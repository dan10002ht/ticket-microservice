package com.ticketing.payment.service;

import com.ticketing.payment.adapter.GatewayRequestContext;
import com.ticketing.payment.adapter.PaymentGatewayRegistry;
import com.ticketing.payment.adapter.GatewayResponse;
import com.ticketing.payment.entity.IdempotencyKey;
import com.ticketing.payment.entity.Payment;
import com.ticketing.payment.entity.enums.PaymentStatus;
import com.ticketing.payment.entity.enums.TransactionStatus;
import com.ticketing.payment.entity.enums.TransactionType;
import com.ticketing.payment.repository.PaymentRepository;
import com.ticketing.payment.service.dto.IdempotencyKeyContext;
import com.ticketing.payment.service.dto.PaymentCreateCommand;
import com.ticketing.payment.service.dto.PaymentResponse;
import com.ticketing.payment.service.mapper.PaymentMapper;
import jakarta.transaction.Transactional;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

/**
 * Core payment domain service. This layer encapsulates business rules,
 * validation, idempotency, and orchestration with external providers (Stripe,
 * etc.). For now we scaffold the structure with TODO markers so that future
 * work can focus on gateway adapters without reshaping the service.
 */
@Service
@Validated
@RequiredArgsConstructor
@Slf4j
public class PaymentService {

    private final PaymentRepository paymentRepository;
    private final PaymentMapper paymentMapper;
    private final TransactionLogService transactionLogService;
    private final IdempotencyKeyService idempotencyKeyService;
    private final PaymentGatewayRegistry gatewayRegistry;

    /**
     * Create payment with idempotency guard. Actual gateway capture logic will be
     * plugged in later.
     */
    @Transactional
    public PaymentResponse createPayment(@Valid PaymentCreateCommand command, IdempotencyKeyContext context) {
        IdempotencyKey key = idempotencyKeyService.registerOrRetrieve(context)
                .orElseThrow(() -> new IllegalArgumentException("Idempotency key is required"));

        if (key.isCompleted()) {
            throw new IllegalStateException("Payment already processed for this idempotency key");
        }

        Payment payment = paymentMapper.toEntity(command);
        payment.markAsProcessing();
        payment = paymentRepository.save(payment);
        key.setPayment(payment);
        key.setPaymentUuid(payment.getPaymentId());

        GatewayRequestContext gatewayContext = GatewayRequestContext.builder()
                .paymentCommand(command)
                .build();
        GatewayResponse gatewayResponse = gatewayRegistry.resolve(command.getGatewayProvider())
                .authorize(command, gatewayContext);
        payment.setGatewayResponse(gatewayResponse.getRawPayload());
        payment.setProviderReference(gatewayResponse.getProviderReference());

        if (!gatewayResponse.isSuccessful()) {
            transactionLogService.logPaymentEvent(
                    TransactionType.PAYMENT_FAILED,
                    TransactionStatus.FAILED,
                    payment,
                    "payment.failed",
                    gatewayResponse.getProviderReference(),
                    gatewayResponse.getErrorCode(),
                    gatewayResponse.getErrorMessage());
            payment.markAsFailed(gatewayResponse.getErrorMessage());
            paymentRepository.save(payment);

            Map<String, Object> failedPayload = new HashMap<>();
            failedPayload.put("error", gatewayResponse.getErrorMessage());
            idempotencyKeyService.markFailed(key, 422, failedPayload);
            throw new IllegalStateException("Payment failed at gateway: " + gatewayResponse.getErrorMessage());
        }

        transactionLogService.logPaymentEvent(
                TransactionType.PAYMENT_INITIATED,
                TransactionStatus.SUCCESS,
                payment,
                "payment.initiated",
                gatewayResponse.getProviderReference(),
                null,
                null);
        Map<String, Object> payload = new HashMap<>();
        payload.put("paymentId", payment.getPaymentId().toString());
        payload.put("providerReference", gatewayResponse.getProviderReference());
        idempotencyKeyService.markCompleted(key, 202, payload);

        log.info("Payment {} created for booking {}", payment.getPaymentId(), payment.getBookingId());
        return paymentMapper.toResponse(payment);
    }

    /**
     * Retrieve payment by UUID.
     */
    @Transactional
    public PaymentResponse getPayment(UUID paymentId) {
        Payment payment = paymentRepository.findByPaymentId(paymentId)
                .orElseThrow(() -> new IllegalArgumentException("Payment not found"));
        return paymentMapper.toResponse(payment);
    }

    /**
     * List payments for a user filtered by status.
     */
    @Transactional
    public Page<PaymentResponse> listPayments(String userId, PaymentStatus status, Pageable pageable) {
        return paymentRepository.findByUserIdAndStatusOrderByCreatedAtDesc(userId, status, pageable)
                .map(paymentMapper::toResponse);
    }

    /**
     * Mark payment as successful. Later this will be triggered by gateway
     * callback or internal flow.
     */
    @Transactional
    public PaymentResponse markPaymentSuccess(UUID paymentId, String externalReference) {
        Payment payment = paymentRepository.findByPaymentId(paymentId)
                .orElseThrow(() -> new IllegalArgumentException("Payment not found"));
        payment.markAsSuccess(LocalDateTime.now());
        payment.setExternalReference(externalReference);
        paymentRepository.save(payment);

        transactionLogService.logPaymentEvent(
                TransactionType.PAYMENT_SUCCESS,
                TransactionStatus.SUCCESS,
                payment,
                "payment.success",
                externalReference,
                null,
                null);
        return paymentMapper.toResponse(payment);
    }

    /**
     * Mark payment as failed with reason.
     */
    @Transactional
    public PaymentResponse markPaymentFailed(UUID paymentId, String failureReason) {
        Payment payment = paymentRepository.findByPaymentId(paymentId)
                .orElseThrow(() -> new IllegalArgumentException("Payment not found"));
        payment.markAsFailed(failureReason);
        paymentRepository.save(payment);

        transactionLogService.logPaymentEvent(
                TransactionType.PAYMENT_FAILED,
                TransactionStatus.FAILED,
                payment,
                "payment.failed",
                null,
                "PAYMENT_FAILED",
                failureReason);
        return paymentMapper.toResponse(payment);
    }

    /**
     * Capture a previously authorized payment. This is used for manual capture
     * flows where payment was authorized but not yet captured.
     */
    @Transactional
    public PaymentResponse capturePayment(UUID paymentId) {
        Payment payment = paymentRepository.findByPaymentId(paymentId)
                .orElseThrow(() -> new IllegalArgumentException("Payment not found"));

        if (payment.getProviderReference() == null) {
            throw new IllegalStateException("Payment does not have a provider reference");
        }

        if (payment.getStatus() != PaymentStatus.PROCESSING) {
            throw new IllegalStateException("Payment must be in PROCESSING status to capture");
        }

        GatewayRequestContext gatewayContext = GatewayRequestContext.builder()
                .providerReference(payment.getProviderReference())
                .build();

        PaymentCreateCommand dummyCommand = PaymentCreateCommand.builder()
                .gatewayProvider(payment.getGatewayProvider())
                .build();

        GatewayResponse gatewayResponse = gatewayRegistry.resolve(payment.getGatewayProvider())
                .capture(dummyCommand, gatewayContext);

        payment.setGatewayResponse(gatewayResponse.getRawPayload());

        if (gatewayResponse.isSuccessful()) {
            payment.markAsSuccess(LocalDateTime.now());
            transactionLogService.logPaymentEvent(
                    TransactionType.PAYMENT_CAPTURED,
                    TransactionStatus.SUCCESS,
                    payment,
                    "payment.captured",
                    gatewayResponse.getProviderReference(),
                    null,
                    null);
        } else {
            payment.markAsFailed(gatewayResponse.getErrorMessage());
            transactionLogService.logPaymentEvent(
                    TransactionType.PAYMENT_FAILED,
                    TransactionStatus.FAILED,
                    payment,
                    "payment.capture.failed",
                    gatewayResponse.getProviderReference(),
                    gatewayResponse.getErrorCode(),
                    gatewayResponse.getErrorMessage());
        }

        paymentRepository.save(payment);
        log.info("Payment {} capture attempted, status: {}", paymentId, gatewayResponse.isSuccessful());
        return paymentMapper.toResponse(payment);
    }

    /**
     * Cancel a payment that is in PROCESSING status. This releases any held
     * funds or authorization.
     */
    @Transactional
    public PaymentResponse cancelPayment(UUID paymentId, String reason) {
        Payment payment = paymentRepository.findByPaymentId(paymentId)
                .orElseThrow(() -> new IllegalArgumentException("Payment not found"));

        if (payment.getProviderReference() == null) {
            throw new IllegalStateException("Payment does not have a provider reference");
        }

        if (payment.getStatus() != PaymentStatus.PROCESSING) {
            throw new IllegalStateException("Payment must be in PROCESSING status to cancel");
        }

        GatewayRequestContext gatewayContext = GatewayRequestContext.builder()
                .providerReference(payment.getProviderReference())
                .build();

        GatewayResponse gatewayResponse = gatewayRegistry.resolve(payment.getGatewayProvider())
                .cancelPayment(payment.getProviderReference(), reason, gatewayContext);

        payment.setGatewayResponse(gatewayResponse.getRawPayload());

        if (gatewayResponse.isSuccessful()) {
            payment.markAsCancelled("system", reason != null ? reason : "Payment cancelled via gateway");
            transactionLogService.logPaymentEvent(
                    TransactionType.PAYMENT_CANCELLED,
                    TransactionStatus.SUCCESS,
                    payment,
                    "payment.cancelled",
                    gatewayResponse.getProviderReference(),
                    null,
                    null);
        } else {
            transactionLogService.logPaymentEvent(
                    TransactionType.PAYMENT_FAILED,
                    TransactionStatus.FAILED,
                    payment,
                    "payment.cancel.failed",
                    gatewayResponse.getProviderReference(),
                    gatewayResponse.getErrorCode(),
                    gatewayResponse.getErrorMessage());
            throw new IllegalStateException("Failed to cancel payment: " + gatewayResponse.getErrorMessage());
        }

        paymentRepository.save(payment);
        log.info("Payment {} cancelled", paymentId);
        return paymentMapper.toResponse(payment);
    }

    // Legacy helper kept for future use (e.g., when context not available)
    private void verifyIdempotency(String idempotencyKey) {
        if (idempotencyKey == null) {
            return;
        }
        if (paymentRepository.existsByIdempotencyKey(idempotencyKey)) {
            throw new IllegalStateException("Duplicate request (idempotency violation)");
        }
    }
}
