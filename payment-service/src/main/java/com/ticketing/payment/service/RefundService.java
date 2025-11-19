package com.ticketing.payment.service;

import com.ticketing.payment.entity.IdempotencyKey;
import com.ticketing.payment.adapter.GatewayRequestContext;
import com.ticketing.payment.adapter.GatewayResponse;
import com.ticketing.payment.adapter.PaymentGatewayRegistry;
import com.ticketing.payment.entity.Payment;
import com.ticketing.payment.entity.Refund;
import com.ticketing.payment.entity.enums.RefundStatus;
import com.ticketing.payment.entity.enums.TransactionStatus;
import com.ticketing.payment.entity.enums.TransactionType;
import com.ticketing.payment.repository.PaymentRepository;
import com.ticketing.payment.repository.RefundRepository;
import com.ticketing.payment.service.dto.IdempotencyKeyContext;
import com.ticketing.payment.service.dto.RefundCreateCommand;
import com.ticketing.payment.service.dto.RefundResponse;
import com.ticketing.payment.service.mapper.RefundMapper;
import jakarta.transaction.Transactional;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

/**
 * Service handling refund lifecycle.
 */
@Service
@RequiredArgsConstructor
@Slf4j
public class RefundService {

    private final PaymentRepository paymentRepository;
    private final RefundRepository refundRepository;
    private final RefundMapper refundMapper;
    private final TransactionLogService transactionLogService;
    private final IdempotencyKeyService idempotencyKeyService;
    private final PaymentGatewayRegistry paymentGatewayRegistry;

    /**
     * Initiate refund request for a payment.
     */
    @Transactional
    public RefundResponse createRefund(@Valid RefundCreateCommand command, IdempotencyKeyContext context) {
        IdempotencyKey key = idempotencyKeyService.registerOrRetrieve(context)
                .orElseThrow(() -> new IllegalArgumentException("Idempotency key is required"));

        if (key.isCompleted()) {
            throw new IllegalStateException("Refund already processed for this idempotency key");
        }

        Payment payment = paymentRepository.findByPaymentId(command.getPaymentId())
                .orElseThrow(() -> new IllegalArgumentException("Payment not found"));

        validateRefundAmount(payment, command.getAmount());

        Refund refund = refundMapper.toEntity(payment, command);
        refund.setGatewayProvider(payment.getGatewayProvider());
        refund = refundRepository.save(refund);
        key.setPayment(payment);
        key.setPaymentUuid(payment.getPaymentId());
        key.setRefund(refund);
        key.setRefundUuid(refund.getRefundId());

        transactionLogService.logRefundEvent(
                TransactionType.REFUND_INITIATED,
                TransactionStatus.SUCCESS,
                refund,
                "refund.initiated",
                null,
                null,
                null);

        if (payment.getProviderReference() == null) {
            throw new IllegalStateException("Payment missing provider reference, cannot refund");
        }

        GatewayRequestContext gatewayContext = GatewayRequestContext.builder()
                .refundCommand(command)
                .providerReference(payment.getProviderReference())
                .build();

        GatewayResponse gatewayResponse = paymentGatewayRegistry.resolve(payment.getGatewayProvider())
                .refund(command, gatewayContext);

        if (!gatewayResponse.isSuccessful()) {
            refund.markAsFailed(gatewayResponse.getErrorMessage());
            refundRepository.save(refund);
            transactionLogService.logRefundEvent(
                    TransactionType.REFUND_FAILED,
                    TransactionStatus.FAILED,
                    refund,
                    "refund.failed",
                    gatewayResponse.getProviderReference(),
                    gatewayResponse.getErrorCode(),
                    gatewayResponse.getErrorMessage());

            Map<String, Object> failedPayload = new HashMap<>();
            failedPayload.put("error", gatewayResponse.getErrorMessage());
            idempotencyKeyService.markFailed(key, 422, failedPayload);
            throw new IllegalStateException("Gateway refund failed: " + gatewayResponse.getErrorMessage());
        }

        refund.setProviderReference(gatewayResponse.getProviderReference());
        refundRepository.save(refund);

        transactionLogService.logRefundEvent(
                TransactionType.REFUND_SUCCESS,
                TransactionStatus.SUCCESS,
                refund,
                "refund.success",
                gatewayResponse.getProviderReference(),
                null,
                null);

        log.info("Refund {} created for payment {}", refund.getRefundId(), refund.getPaymentUuid());
        RefundResponse response = refundMapper.toResponse(refund);
        Map<String, Object> payload = new HashMap<>();
        payload.put("refundId", refund.getRefundId().toString());
        payload.put("providerReference", gatewayResponse.getProviderReference());
        idempotencyKeyService.markCompleted(key, 202, payload);
        return response;
    }

    /**
     * List refunds for a payment.
     */
    public List<RefundResponse> listRefunds(UUID paymentId) {
        return refundRepository.findByPayment_PaymentIdOrderByCreatedAtDesc(paymentId).stream()
                .map(refundMapper::toResponse)
                .toList();
    }

    /**
     * Mark refund success/failure (later, triggered by gateway webhook).
     */
    @Transactional
    public RefundResponse markRefundStatus(UUID refundId, RefundStatus status, String externalReference,
            String failureReason) {
        Refund refund = refundRepository.findByRefundId(refundId)
                .orElseThrow(() -> new IllegalArgumentException("Refund not found"));

        switch (status) {
            case SUCCESS -> {
                refund.markAsSuccess(LocalDateTime.now());
                refund.setExternalReference(externalReference);
                transactionLogService.logRefundEvent(
                        TransactionType.REFUND_SUCCESS,
                        TransactionStatus.SUCCESS,
                        refund,
                        "refund.success",
                        externalReference,
                        null,
                        null);
            }
            case FAILED -> {
                refund.markAsFailed(failureReason);
                transactionLogService.logRefundEvent(
                        TransactionType.REFUND_FAILED,
                        TransactionStatus.FAILED,
                        refund,
                        "refund.failed",
                        null,
                        "REFUND_FAILED",
                        failureReason);
            }
            case CANCELLED -> {
                refund.markAsCancelled("system", failureReason);
                transactionLogService.logRefundEvent(
                        TransactionType.REFUND_FAILED,
                        TransactionStatus.FAILED,
                        refund,
                        "refund.cancelled",
                        null,
                        "REFUND_CANCELLED",
                        failureReason);
            }
            default -> {
                // For pending/processing just update status
                refund.setStatus(status);
                transactionLogService.logRefundEvent(
                        TransactionType.REFUND_INITIATED,
                        TransactionStatus.PENDING,
                        refund,
                        "refund." + status.name().toLowerCase(),
                        null,
                        null,
                        null);
            }
        }

        refundRepository.save(refund);
        // TODO: update Payment aggregate status based on total refunded amount
        // TODO: log transaction event
        return refundMapper.toResponse(refund);
    }

    private void validateRefundAmount(Payment payment, java.math.BigDecimal amount) {
        if (amount.compareTo(payment.getRemainingRefundableAmount()) > 0) {
            throw new IllegalArgumentException("Refund amount exceeds remaining balance");
        }
    }
}
