package com.ticketing.booking.grpcclient;

import java.util.List;

import org.springframework.retry.annotation.Backoff;
import org.springframework.retry.annotation.Retryable;
import org.springframework.stereotype.Component;

import com.ticketing.booking.grpc.TicketProto;
import com.ticketing.booking.grpc.TicketProto.TicketServiceGrpc;

import io.grpc.ManagedChannel;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Component
@RequiredArgsConstructor
public class TicketServiceClient {

    private final ManagedChannel channel;
    private TicketServiceGrpc.TicketServiceBlockingStub stub;

    private TicketServiceGrpc.TicketServiceBlockingStub getStub() {
        if (stub == null) {
            stub = TicketServiceGrpc.newBlockingStub(channel);
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
    public TicketProto.ReserveTicketsResponse reserveTickets(
            String eventId,
            List<String> seatNumbers,
            String userId,
            int timeoutSeconds) {
        try {
            TicketProto.ReserveTicketsRequest request = TicketProto.ReserveTicketsRequest.newBuilder()
                    .setEventId(eventId)
                    .addAllSeatNumbers(seatNumbers)
                    .setUserId(userId)
                    .setTimeoutSeconds(timeoutSeconds)
                    .build();

            return getStub().reserveTickets(request);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            // Don't retry on INVALID_ARGUMENT, NOT_FOUND, PERMISSION_DENIED
            if (status.getCode() == Status.Code.INVALID_ARGUMENT
                    || status.getCode() == Status.Code.NOT_FOUND
                    || status.getCode() == Status.Code.PERMISSION_DENIED) {
                log.warn("Non-retryable gRPC error: {}", status.getCode(), e);
                throw e;
            }
            log.warn("Retryable gRPC error during ticket reservation: {}", status.getCode(), e);
            throw e;
        } catch (Exception e) {
            log.error("Unexpected error during ticket reservation", e);
            throw new RuntimeException("Ticket reservation failed", e);
        }
    }

    @Retryable(
            retryFor = { StatusRuntimeException.class, RuntimeException.class },
            maxAttemptsExpression = "${booking.grpc.retry.max-attempts:3}",
            backoff = @Backoff(
                    delayExpression = "${booking.grpc.retry.initial-interval-ms:500}",
                    multiplierExpression = "${booking.grpc.retry.multiplier:2.0}",
                    maxDelayExpression = "${booking.grpc.retry.max-interval-ms:5000}"))
    public TicketProto.ReleaseTicketsResponse releaseTickets(String reservationId, List<String> ticketIds) {
        try {
            TicketProto.ReleaseTicketsRequest.Builder builder = TicketProto.ReleaseTicketsRequest.newBuilder()
                    .setReservationId(reservationId);

            if (ticketIds != null && !ticketIds.isEmpty()) {
                builder.addAllTicketIds(ticketIds);
            }

            return getStub().releaseTickets(builder.build());
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            // Don't retry on INVALID_ARGUMENT, NOT_FOUND
            if (status.getCode() == Status.Code.INVALID_ARGUMENT
                    || status.getCode() == Status.Code.NOT_FOUND) {
                log.warn("Non-retryable gRPC error: {}", status.getCode(), e);
                throw e;
            }
            log.warn("Retryable gRPC error during ticket release: {}", status.getCode(), e);
            throw e;
        } catch (Exception e) {
            log.error("Unexpected error during ticket release", e);
            throw new RuntimeException("Ticket release failed", e);
        }
    }

    @Retryable(
            retryFor = { StatusRuntimeException.class, RuntimeException.class },
            maxAttemptsExpression = "${booking.grpc.retry.max-attempts:3}",
            backoff = @Backoff(
                    delayExpression = "${booking.grpc.retry.initial-interval-ms:500}",
                    multiplierExpression = "${booking.grpc.retry.multiplier:2.0}",
                    maxDelayExpression = "${booking.grpc.retry.max-interval-ms:5000}"))
    public TicketProto.CheckAvailabilityResponse checkAvailability(String eventId, List<String> seatNumbers) {
        try {
            TicketProto.CheckAvailabilityRequest request = TicketProto.CheckAvailabilityRequest.newBuilder()
                    .setEventId(eventId)
                    .addAllSeatNumbers(seatNumbers)
                    .build();

            return getStub().checkAvailability(request);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            // Don't retry on INVALID_ARGUMENT, NOT_FOUND
            if (status.getCode() == Status.Code.INVALID_ARGUMENT
                    || status.getCode() == Status.Code.NOT_FOUND) {
                log.warn("Non-retryable gRPC error: {}", status.getCode(), e);
                throw e;
            }
            log.warn("Retryable gRPC error during availability check: {}", status.getCode(), e);
            throw e;
        } catch (Exception e) {
            log.error("Unexpected error during availability check", e);
            throw new RuntimeException("Availability check failed", e);
        }
    }
}
