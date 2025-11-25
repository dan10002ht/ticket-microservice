package com.ticketing.payment.grpc;

import com.google.protobuf.Timestamp;
import com.ticketing.payment.entity.Payment;
import com.ticketing.payment.entity.Refund;
import com.ticketing.payment.entity.enums.PaymentGateway;
import com.ticketing.payment.entity.enums.PaymentMethod;
import com.ticketing.payment.entity.enums.PaymentStatus;
import com.ticketing.payment.entity.enums.RefundStatus;
import com.ticketing.payment.entity.enums.RefundType;
import com.ticketing.payment.repository.PaymentRepository;
import com.ticketing.payment.repository.RefundRepository;
import com.ticketing.payment.service.PaymentService;
import com.ticketing.payment.service.RefundService;
import com.ticketing.payment.service.dto.IdempotencyKeyContext;
import com.ticketing.payment.service.dto.PaymentCreateCommand;
import com.ticketing.payment.service.dto.RefundCreateCommand;
import com.ticketing.payment.webhook.PaymentWebhookHandler;
import com.ticketing.payment.webhook.WebhookEvent;
import com.ticketing.payment.webhook.WebhookHandlerRegistry;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.server.service.GrpcService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;

import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.UUID;

/**
 * gRPC adapter exposing payment/refund features to internal services.
 */
@GrpcService
@RequiredArgsConstructor
@Slf4j
public class PaymentGrpcService extends PaymentServiceGrpc.PaymentServiceImplBase {

    private static final String PATH_CREATE_PAYMENT = "/grpc/payment.PaymentService/CreatePayment";
    private static final String PATH_CREATE_REFUND = "/grpc/payment.PaymentService/CreateRefund";

    private final PaymentService paymentService;
    private final RefundService refundService;
    private final PaymentRepository paymentRepository;
    private final RefundRepository refundRepository;
    private final WebhookHandlerRegistry webhookHandlerRegistry;

    @Override
    public void createPayment(CreatePaymentRequest request, StreamObserver<PaymentResponse> responseObserver) {
        try {
            String idempotencyKey = requireNonBlank(request.getIdempotencyKey(), "idempotency_key");

            PaymentCreateCommand command = PaymentCreateCommand.builder()
                    .bookingId(requireNonBlank(request.getBookingId(), "booking_id"))
                    .ticketId(emptyToNull(request.getTicketId()))
                    .userId(requireNonBlank(request.getUserId(), "user_id"))
                    .amount(toBigDecimal(request.getAmount()))
                    .currency(requireNonBlank(request.getCurrency(), "currency"))
                    .paymentMethod(
                            PaymentMethod.valueOf(
                                    requireNonBlank(request.getPaymentMethod(), "payment_method").toUpperCase()))
                    .gatewayProvider(
                            PaymentGateway.valueOf(
                                    requireNonBlank(request.getGatewayProvider(), "gateway_provider").toUpperCase()))
                    .idempotencyKey(idempotencyKey)
                    .metadata(toObjectMap(request.getMetadataMap()))
                    .build();

            IdempotencyKeyContext context = IdempotencyKeyContext.builder()
                    .key(idempotencyKey)
                    .requestPath(PATH_CREATE_PAYMENT)
                    .requestMethod("GRPC")
                    .userId(request.getUserId())
                    .requestBody(toObjectMap(request.getMetadataMap()))
                    .build();

            com.ticketing.payment.service.dto.PaymentResponse serviceResponse = paymentService.createPayment(command,
                    context);
            responseObserver.onNext(PaymentResponse.newBuilder()
                    .setPayment(toProtoPayment(serviceResponse))
                    .build());
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (Exception e) {
            log.error("Failed to create payment", e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void getPayment(GetPaymentRequest request, StreamObserver<PaymentResponse> responseObserver) {
        try {
            UUID paymentId = parseUuid(requireNonBlank(request.getPaymentId(), "payment_id"), "payment_id");
            com.ticketing.payment.service.dto.PaymentResponse serviceResponse = paymentService.getPayment(paymentId);
            responseObserver.onNext(PaymentResponse.newBuilder()
                    .setPayment(toProtoPayment(serviceResponse))
                    .build());
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (Exception e) {
            log.error("Failed to fetch payment {}", request.getPaymentId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void listPayments(ListPaymentsRequest request, StreamObserver<ListPaymentsResponse> responseObserver) {
        try {
            PaymentStatus status = PaymentStatus.valueOf(requireNonBlank(request.getStatus(), "status").toUpperCase());
            int page = Math.max(request.getPage(), 0);
            int size = request.getSize() > 0 ? request.getSize() : 20;
            Page<com.ticketing.payment.service.dto.PaymentResponse> payments = paymentService.listPayments(
                    request.getUserId(),
                    status,
                    PageRequest.of(page, size));
            ListPaymentsResponse.Builder builder = ListPaymentsResponse.newBuilder()
                    .setTotal(payments.getTotalElements());
            payments.forEach(payment -> builder.addPayments(toProtoPayment(payment)));
            responseObserver.onNext(builder.build());
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (Exception e) {
            log.error("Failed to list payments for user {}", request.getUserId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void markPaymentSuccess(MarkPaymentSuccessRequest request,
            StreamObserver<PaymentResponse> responseObserver) {
        try {
            UUID paymentId = parseUuid(request.getPaymentId(), "payment_id");
            com.ticketing.payment.service.dto.PaymentResponse serviceResponse = paymentService.markPaymentSuccess(
                    paymentId,
                    emptyToNull(request.getExternalReference()));
            responseObserver.onNext(PaymentResponse.newBuilder()
                    .setPayment(toProtoPayment(serviceResponse))
                    .build());
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (Exception e) {
            log.error("Failed to mark payment success {}", request.getPaymentId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void markPaymentFailed(MarkPaymentFailedRequest request, StreamObserver<PaymentResponse> responseObserver) {
        try {
            UUID paymentId = parseUuid(request.getPaymentId(), "payment_id");
            com.ticketing.payment.service.dto.PaymentResponse serviceResponse = paymentService.markPaymentFailed(
                    paymentId,
                    request.getFailureReason());
            responseObserver.onNext(PaymentResponse.newBuilder()
                    .setPayment(toProtoPayment(serviceResponse))
                    .build());
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (Exception e) {
            log.error("Failed to mark payment failed {}", request.getPaymentId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void capturePayment(CapturePaymentRequest request, StreamObserver<PaymentResponse> responseObserver) {
        try {
            UUID paymentId = parseUuid(requireNonBlank(request.getPaymentId(), "payment_id"), "payment_id");
            com.ticketing.payment.service.dto.PaymentResponse serviceResponse = paymentService
                    .capturePayment(paymentId);
            responseObserver.onNext(PaymentResponse.newBuilder()
                    .setPayment(toProtoPayment(serviceResponse))
                    .build());
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (Exception e) {
            log.error("Failed to capture payment {}", request.getPaymentId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void cancelPayment(CancelPaymentRequest request, StreamObserver<PaymentResponse> responseObserver) {
        try {
            UUID paymentId = parseUuid(requireNonBlank(request.getPaymentId(), "payment_id"), "payment_id");
            String reason = emptyToNull(request.getReason());
            com.ticketing.payment.service.dto.PaymentResponse serviceResponse = paymentService.cancelPayment(
                    paymentId,
                    reason != null ? reason : "Payment cancelled via gRPC");
            responseObserver.onNext(PaymentResponse.newBuilder()
                    .setPayment(toProtoPayment(serviceResponse))
                    .build());
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (Exception e) {
            log.error("Failed to cancel payment {}", request.getPaymentId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void createRefund(CreateRefundRequest request, StreamObserver<RefundResponse> responseObserver) {
        try {
            String idempotencyKey = requireNonBlank(request.getIdempotencyKey(), "idempotency_key");

            RefundCreateCommand command = RefundCreateCommand.builder()
                    .paymentId(parseUuid(requireNonBlank(request.getPaymentId(), "payment_id"), "payment_id"))
                    .amount(toBigDecimal(request.getAmount()))
                    .reason(requireNonBlank(request.getReason(), "reason"))
                    .description(emptyToNull(request.getDescription()))
                    .refundType(
                            RefundType.valueOf(requireNonBlank(request.getRefundType(), "refund_type").toUpperCase()))
                    .metadata(toObjectMap(request.getMetadataMap()))
                    .build();

            IdempotencyKeyContext context = IdempotencyKeyContext.builder()
                    .key(idempotencyKey)
                    .requestPath(PATH_CREATE_REFUND)
                    .requestMethod("GRPC")
                    .requestBody(toObjectMap(request.getMetadataMap()))
                    .build();

            RefundResponse serviceResponse = refundService.createRefund(command, context);
            responseObserver.onNext(RefundResponse.newBuilder()
                    .setRefund(toProtoRefund(serviceResponse))
                    .build());
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (Exception e) {
            log.error("Failed to create refund", e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void listRefunds(ListRefundsRequest request, StreamObserver<ListRefundsResponse> responseObserver) {
        try {
            UUID paymentId = parseUuid(request.getPaymentId(), "payment_id");
            ListRefundsResponse.Builder builder = ListRefundsResponse.newBuilder();
            refundService.listRefunds(paymentId).forEach(refund -> builder.addRefunds(toProtoRefund(refund)));
            responseObserver.onNext(builder.build());
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (Exception e) {
            log.error("Failed to list refunds for payment {}", request.getPaymentId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void updateRefundStatus(UpdateRefundStatusRequest request, StreamObserver<RefundResponse> responseObserver) {
        try {
            UUID refundId = parseUuid(requireNonBlank(request.getRefundId(), "refund_id"), "refund_id");
            RefundStatus status = RefundStatus.valueOf(requireNonBlank(request.getStatus(), "status").toUpperCase());
            RefundResponse serviceResponse = refundService.markRefundStatus(
                    refundId,
                    status,
                    emptyToNull(request.getExternalReference()),
                    emptyToNull(request.getFailureReason()));
            responseObserver.onNext(RefundResponse.newBuilder()
                    .setRefund(toProtoRefund(serviceResponse))
                    .build());
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (Exception e) {
            log.error("Failed to update refund {}", request.getRefundId(), e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    @Override
    public void processWebhook(ProcessWebhookRequest request, StreamObserver<ProcessWebhookResponse> responseObserver) {
        try {
            PaymentGateway gateway = PaymentGateway
                    .valueOf(requireNonBlank(request.getGateway(), "gateway").toUpperCase());
            PaymentWebhookHandler handler = webhookHandlerRegistry.resolve(gateway);
            WebhookEvent event = handler.handle(request.getPayload(), request.getHeadersMap());

            if (event.getType() != null) {
                switch (event.getType()) {
                    case PAYMENT_SUCCEEDED -> handlePaymentSuccess(event);
                    case PAYMENT_FAILED -> handlePaymentFailure(event);
                    case REFUND_SUCCEEDED -> handleRefund(event, RefundStatus.SUCCESS, null);
                    case REFUND_FAILED -> handleRefund(event, RefundStatus.FAILED, "Gateway failure");
                    default -> log.info("Unhandled webhook event {}", event);
                }
            }

            responseObserver.onNext(ProcessWebhookResponse.newBuilder()
                    .setSuccess(true)
                    .setMessage("Webhook processed")
                    .build());
            responseObserver.onCompleted();
        } catch (Exception e) {
            log.error("Failed to process webhook", e);
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }

    private void handlePaymentSuccess(WebhookEvent event) {
        Payment payment = resolvePayment(event);
        if (payment == null) {
            log.warn("Unable to resolve payment for success event {}", event);
            return;
        }
        paymentService.markPaymentSuccess(payment.getPaymentId(), event.getProviderReference());
    }

    private void handlePaymentFailure(WebhookEvent event) {
        Payment payment = resolvePayment(event);
        if (payment == null) {
            log.warn("Unable to resolve payment for failure event {}", event);
            return;
        }
        paymentService.markPaymentFailed(payment.getPaymentId(), "Gateway failure");
    }

    private void handleRefund(WebhookEvent event, RefundStatus status, String failureReason) {
        Refund refund = resolveRefund(event);
        if (refund == null) {
            log.warn("Unable to resolve refund for event {}", event);
            return;
        }
        refundService.markRefundStatus(refund.getRefundId(), status, event.getProviderReference(), failureReason);
    }

    private Payment resolvePayment(WebhookEvent event) {
        if (StringUtils.isNotBlank(event.getPaymentId())) {
            return paymentRepository.findByPaymentId(UUID.fromString(event.getPaymentId())).orElse(null);
        }
        if (StringUtils.isNotBlank(event.getProviderReference())) {
            return paymentRepository.findByProviderReference(event.getProviderReference()).orElse(null);
        }
        return null;
    }

    private Refund resolveRefund(WebhookEvent event) {
        if (StringUtils.isNotBlank(event.getRefundId())) {
            return refundRepository.findByRefundId(UUID.fromString(event.getRefundId())).orElse(null);
        }
        if (StringUtils.isNotBlank(event.getProviderReference())) {
            return refundRepository.findByProviderReference(event.getProviderReference()).orElse(null);
        }
        return null;
    }

    @Override
    public void health(HealthRequest request, StreamObserver<HealthResponse> responseObserver) {
        responseObserver.onNext(HealthResponse.newBuilder()
                .setStatus("SERVING")
                .setMessage("Payment Service is healthy")
                .build());
        responseObserver.onCompleted();
    }

    private com.ticketing.payment.grpc.Payment toProtoPayment(
            com.ticketing.payment.service.dto.PaymentResponse payment) {
        com.ticketing.payment.grpc.Payment.Builder builder = com.ticketing.payment.grpc.Payment.newBuilder()
                .setId(Optional.ofNullable(payment.getId()).map(String::valueOf).orElse(""))
                .setPaymentId(optionalToString(payment.getPaymentId()))
                .setBookingId(optionalToString(payment.getBookingId()))
                .setTicketId(optionalToString(payment.getTicketId()))
                .setUserId(optionalToString(payment.getUserId()))
                .setAmount(payment.getAmount().doubleValue())
                .setCurrency(optionalToString(payment.getCurrency()))
                .setPaymentMethod(payment.getPaymentMethod().name())
                .setStatus(payment.getStatus().name())
                .setGatewayProvider(Optional.ofNullable(payment.getGatewayProvider()).map(Enum::name).orElse(""))
                .setExternalReference(optionalToString(payment.getExternalReference()))
                .setFailureReason(optionalToString(payment.getFailureReason()))
                .setDescription(optionalToString(payment.getDescription()))
                .putAllMetadata(convertMetadata(payment.getMetadata()))
                .setIdempotencyKey(optionalToString(payment.getIdempotencyKey()));

        Optional.ofNullable(toTimestamp(payment.getCreatedAt())).ifPresent(builder::setCreatedAt);
        Optional.ofNullable(toTimestamp(payment.getUpdatedAt())).ifPresent(builder::setUpdatedAt);
        Optional.ofNullable(toTimestamp(payment.getPaidAt())).ifPresent(builder::setPaidAt);
        Optional.ofNullable(toTimestamp(payment.getCancelledAt())).ifPresent(builder::setCancelledAt);
        return builder.build();
    }

    private com.ticketing.payment.grpc.Refund toProtoRefund(
            com.ticketing.payment.service.dto.RefundResponse refund) {
        com.ticketing.payment.grpc.Refund.Builder builder = com.ticketing.payment.grpc.Refund.newBuilder()
                .setId(Optional.ofNullable(refund.getId()).map(String::valueOf).orElse(""))
                .setRefundId(optionalToString(refund.getRefundId()))
                .setPaymentId(optionalToString(refund.getPaymentId()))
                .setAmount(refund.getAmount().doubleValue())
                .setCurrency(optionalToString(refund.getCurrency()))
                .setRefundType(refund.getRefundType().name())
                .setStatus(refund.getStatus().name())
                .setGatewayProvider(Optional.ofNullable(refund.getGatewayProvider()).map(Enum::name).orElse(""))
                .setReason(optionalToString(refund.getReason()))
                .setDescription(optionalToString(refund.getDescription()))
                .setFailureReason(optionalToString(refund.getFailureReason()))
                .putAllMetadata(convertMetadata(refund.getMetadata()));

        Optional.ofNullable(toTimestamp(refund.getCreatedAt())).ifPresent(builder::setCreatedAt);
        Optional.ofNullable(toTimestamp(refund.getUpdatedAt())).ifPresent(builder::setUpdatedAt);
        Optional.ofNullable(toTimestamp(refund.getRefundedAt())).ifPresent(builder::setRefundedAt);
        Optional.ofNullable(toTimestamp(refund.getCancelledAt())).ifPresent(builder::setCancelledAt);
        return builder.build();
    }

    private UUID parseUuid(String value, String fieldName) {
        try {
            return UUID.fromString(value);
        } catch (Exception ex) {
            throw invalidArgument("Invalid " + fieldName, ex);
        }
    }

    private java.math.BigDecimal toBigDecimal(double amount) {
        return java.math.BigDecimal.valueOf(amount);
    }

    private Timestamp toTimestamp(LocalDateTime time) {
        if (time == null) {
            return null;
        }
        Instant instant = time.toInstant(ZoneOffset.UTC);
        return Timestamp.newBuilder()
                .setSeconds(instant.getEpochSecond())
                .setNanos(instant.getNano())
                .build();
    }

    private Map<String, String> convertMetadata(Map<String, Object> metadata) {
        Map<String, String> result = new HashMap<>();
        if (metadata == null) {
            return result;
        }
        metadata.forEach((key, value) -> result.put(key, Optional.ofNullable(value).map(Object::toString).orElse("")));
        return result;
    }

    private Map<String, Object> toObjectMap(Map<String, String> metadata) {
        Map<String, Object> result = new HashMap<>();
        if (metadata == null) {
            return result;
        }
        result.putAll(metadata);
        return result;
    }

    private String requireNonBlank(String value, String field) {
        if (StringUtils.isBlank(value)) {
            throw invalidArgument(field + " is required");
        }
        return value;
    }

    private StatusRuntimeException invalidArgument(String message) {
        return Status.INVALID_ARGUMENT.withDescription(message).asRuntimeException();
    }

    private StatusRuntimeException invalidArgument(String message, Throwable cause) {
        return Status.INVALID_ARGUMENT.withDescription(message).withCause(cause).asRuntimeException();
    }

    private String optionalToString(Object value) {
        return value == null ? "" : value.toString();
    }

    private String emptyToNull(String value) {
        return StringUtils.isBlank(value) ? null : value;
    }
}
