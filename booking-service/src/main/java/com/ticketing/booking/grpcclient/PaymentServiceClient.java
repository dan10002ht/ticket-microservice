package com.ticketing.booking.grpcclient;

import java.math.BigDecimal;
import java.util.Map;

import org.springframework.retry.annotation.Backoff;
import org.springframework.retry.annotation.Retryable;
import org.springframework.stereotype.Component;

import com.ticketing.payment.grpc.PaymentProto;
import com.ticketing.payment.grpc.PaymentProto.PaymentServiceGrpc;

import io.grpc.ManagedChannel;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Component
@RequiredArgsConstructor
public class PaymentServiceClient {

    private final ManagedChannel channel;
    private PaymentServiceGrpc.PaymentServiceBlockingStub stub;

    private PaymentServiceGrpc.PaymentServiceBlockingStub getStub() {
        if (stub == null) {
            stub = PaymentServiceGrpc.newBlockingStub(channel);
        }
        return stub;
    }

    @Retryable(
            retryFor = { StatusRuntimeException.class, RuntimeException.class },
            maxAttemptsExpression = "${booking.grpc.retry.max-attempts:3}",
            backoff = @Backoff(
                    delayExpression = "${booking.grpc.retry.initial-interval-ms:500}",
                    multiplierExpression = "${booking.grpc.retry.multiplier:2.0}",
                    maxDelayExpression = "${booking.grpc.retry.max-interval-ms:5000}"))
    public PaymentProto.Payment createPayment(
            String bookingId,
            String userId,
            BigDecimal amount,
            String currency,
            String paymentMethod,
            String gatewayProvider,
            String idempotencyKey,
            Map<String, String> metadata) {
        try {
            PaymentProto.CreatePaymentRequest request = PaymentProto.CreatePaymentRequest.newBuilder()
                    .setBookingId(bookingId)
                    .setUserId(userId)
                    .setAmount(amount.doubleValue())
                    .setCurrency(currency)
                    .setPaymentMethod(paymentMethod)
                    .setGatewayProvider(gatewayProvider)
                    .setIdempotencyKey(idempotencyKey)
                    .putAllMetadata(metadata != null ? metadata : Map.of())
                    .build();

            PaymentProto.PaymentResponse response = getStub().createPayment(request);
            return response.getPayment();
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            // Don't retry on INVALID_ARGUMENT, NOT_FOUND, PERMISSION_DENIED, FAILED_PRECONDITION
            if (status.getCode() == Status.Code.INVALID_ARGUMENT
                    || status.getCode() == Status.Code.NOT_FOUND
                    || status.getCode() == Status.Code.PERMISSION_DENIED
                    || status.getCode() == Status.Code.FAILED_PRECONDITION) {
                log.warn("Non-retryable gRPC error: {}", status.getCode(), e);
                throw e;
            }
            log.warn("Retryable gRPC error during payment creation: {}", status.getCode(), e);
            throw e;
        } catch (Exception e) {
            log.error("Unexpected error during payment creation", e);
            throw new RuntimeException("Payment service call failed", e);
        }
    }

    @Retryable(
            retryFor = { StatusRuntimeException.class, RuntimeException.class },
            maxAttemptsExpression = "${booking.grpc.retry.max-attempts:3}",
            backoff = @Backoff(
                    delayExpression = "${booking.grpc.retry.initial-interval-ms:500}",
                    multiplierExpression = "${booking.grpc.retry.multiplier:2.0}",
                    maxDelayExpression = "${booking.grpc.retry.max-interval-ms:5000}"))
    public PaymentProto.Payment capturePayment(String paymentId) {
        try {
            PaymentProto.CapturePaymentRequest request = PaymentProto.CapturePaymentRequest.newBuilder()
                    .setPaymentId(paymentId)
                    .build();

            PaymentProto.PaymentResponse response = getStub().capturePayment(request);
            return response.getPayment();
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            // Don't retry on INVALID_ARGUMENT, NOT_FOUND, FAILED_PRECONDITION
            if (status.getCode() == Status.Code.INVALID_ARGUMENT
                    || status.getCode() == Status.Code.NOT_FOUND
                    || status.getCode() == Status.Code.FAILED_PRECONDITION) {
                log.warn("Non-retryable gRPC error: {}", status.getCode(), e);
                throw e;
            }
            log.warn("Retryable gRPC error during payment capture: {}", status.getCode(), e);
            throw e;
        } catch (Exception e) {
            log.error("Unexpected error during payment capture", e);
            throw new RuntimeException("Payment capture failed", e);
        }
    }

    @Retryable(
            retryFor = { StatusRuntimeException.class, RuntimeException.class },
            maxAttemptsExpression = "${booking.grpc.retry.max-attempts:3}",
            backoff = @Backoff(
                    delayExpression = "${booking.grpc.retry.initial-interval-ms:500}",
                    multiplierExpression = "${booking.grpc.retry.multiplier:2.0}",
                    maxDelayExpression = "${booking.grpc.retry.max-interval-ms:5000}"))
    public PaymentProto.Payment cancelPayment(String paymentId) {
        try {
            PaymentProto.CancelPaymentRequest request = PaymentProto.CancelPaymentRequest.newBuilder()
                    .setPaymentId(paymentId)
                    .build();

            PaymentProto.PaymentResponse response = getStub().cancelPayment(request);
            return response.getPayment();
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            // Don't retry on INVALID_ARGUMENT, NOT_FOUND, FAILED_PRECONDITION
            if (status.getCode() == Status.Code.INVALID_ARGUMENT
                    || status.getCode() == Status.Code.NOT_FOUND
                    || status.getCode() == Status.Code.FAILED_PRECONDITION) {
                log.warn("Non-retryable gRPC error: {}", status.getCode(), e);
                throw e;
            }
            log.warn("Retryable gRPC error during payment cancellation: {}", status.getCode(), e);
            throw e;
        } catch (Exception e) {
            log.error("Unexpected error during payment cancellation", e);
            throw new RuntimeException("Payment cancellation failed", e);
        }
    }

    @Retryable(
            retryFor = { StatusRuntimeException.class, RuntimeException.class },
            maxAttemptsExpression = "${booking.grpc.retry.max-attempts:3}",
            backoff = @Backoff(
                    delayExpression = "${booking.grpc.retry.initial-interval-ms:500}",
                    multiplierExpression = "${booking.grpc.retry.multiplier:2.0}",
                    maxDelayExpression = "${booking.grpc.retry.max-interval-ms:5000}"))
    public PaymentProto.Payment getPayment(String paymentId) {
        try {
            PaymentProto.GetPaymentRequest request = PaymentProto.GetPaymentRequest.newBuilder()
                    .setPaymentId(paymentId)
                    .build();

            PaymentProto.PaymentResponse response = getStub().getPayment(request);
            return response.getPayment();
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            // Don't retry on INVALID_ARGUMENT, NOT_FOUND
            if (status.getCode() == Status.Code.INVALID_ARGUMENT
                    || status.getCode() == Status.Code.NOT_FOUND) {
                log.warn("Non-retryable gRPC error: {}", status.getCode(), e);
                throw e;
            }
            log.warn("Retryable gRPC error during get payment: {}", status.getCode(), e);
            throw e;
        } catch (Exception e) {
            log.error("Unexpected error during get payment", e);
            throw new RuntimeException("Get payment failed", e);
        }
    }
}

