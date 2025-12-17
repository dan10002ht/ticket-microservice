package com.ticketing.booking.grpc;

import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;

import org.lognet.springboot.grpc.GRpcService;

import com.ticketing.booking.exception.BookingException;
import com.ticketing.booking.exception.BookingLockException;
import com.ticketing.booking.exception.BookingNotFoundException;
import com.ticketing.booking.exception.BookingValidationException;
import com.ticketing.booking.grpc.BookingProto.BookingServiceGrpc.BookingServiceImplBase;
import com.ticketing.booking.service.BookingService;
import com.ticketing.booking.service.dto.BookingCreateCommand;
import com.ticketing.booking.service.mapper.BookingMapper;

import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@GRpcService
@RequiredArgsConstructor
public class BookingGrpcService extends BookingServiceImplBase {

    private final BookingService bookingService;
    private final BookingMapper bookingMapper;

    @Override
    public void createBooking(BookingProto.CreateBookingRequest request,
            StreamObserver<BookingProto.BookingResponse> responseObserver) {
        try {
            // Extract idempotency key from request (for duplicate detection)
            String idempotencyKey = request.getIdempotencyKey();
            if (idempotencyKey != null && idempotencyKey.isEmpty()) {
                idempotencyKey = null;
            }

            BookingCreateCommand command = BookingCreateCommand.builder()
                    .userId(request.getUserId())
                    .eventId(request.getEventId())
                    .seatNumbers(request.getSeatNumbersList())
                    .seatCount(request.getTicketQuantity())
                    .currency("USD")
                    .totalAmount(BigDecimal.valueOf(Math.max(request.getTicketQuantity(), 1)))
                    .metadata(request.getMetadataMap())
                    .idempotencyKey(idempotencyKey) // Pass idempotency key for duplicate prevention
                    .build();

            var result = bookingService.createBooking(command);
            var booking = bookingService.getBooking(result.getBookingId());
            var response = BookingProto.BookingResponse.newBuilder()
                    .setSuccess(true)
                    .setBooking(BookingMapperUtil.toProto(bookingMapper, booking))
                    .setMessage("Booking created")
                    .build();
            responseObserver.onNext(response);
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (BookingValidationException e) {
            log.warn("Validation error in createBooking: {}", e.getMessage());
            responseObserver.onError(Status.INVALID_ARGUMENT
                    .withDescription(e.getMessage())
                    .asRuntimeException());
        } catch (BookingLockException e) {
            log.warn("Lock error in createBooking: {}", e.getMessage());
            responseObserver.onError(Status.RESOURCE_EXHAUSTED
                    .withDescription(e.getMessage())
                    .asRuntimeException());
        } catch (Exception e) {
            log.error("Failed to create booking", e);
            responseObserver.onError(Status.INTERNAL
                    .withDescription("Failed to create booking: " + e.getMessage())
                    .withCause(e)
                    .asRuntimeException());
        }
    }

    @Override
    public void getBooking(BookingProto.GetBookingRequest request,
            StreamObserver<BookingProto.BookingResponse> responseObserver) {
        try {
            UUID bookingId = UUID.fromString(request.getBookingId());
            var booking = bookingService.getBooking(bookingId);
            responseObserver.onNext(BookingProto.BookingResponse.newBuilder()
                    .setSuccess(true)
                    .setBooking(BookingMapperUtil.toProto(bookingMapper, booking))
                    .build());
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (BookingNotFoundException e) {
            log.warn("Booking not found: {}", request.getBookingId());
            responseObserver.onError(Status.NOT_FOUND
                    .withDescription(e.getMessage())
                    .asRuntimeException());
        } catch (IllegalArgumentException e) {
            log.warn("Invalid booking ID: {}", request.getBookingId());
            responseObserver.onError(Status.INVALID_ARGUMENT
                    .withDescription("Invalid booking ID format")
                    .asRuntimeException());
        } catch (Exception e) {
            log.error("Failed to get booking: {}", request.getBookingId(), e);
            responseObserver.onError(Status.INTERNAL
                    .withDescription("Failed to get booking: " + e.getMessage())
                    .withCause(e)
                    .asRuntimeException());
        }
    }

    private static class BookingMapperUtil {
        static BookingProto.Booking toProto(BookingMapper mapper, com.ticketing.booking.entity.Booking booking) {
            var dto = mapper.toResponse(booking);
            double amount = dto.getTotalAmount() == null ? 0D : dto.getTotalAmount().doubleValue();
            return BookingProto.Booking.newBuilder()
                    .setId(String.valueOf(booking.getId()))
                    .setBookingReference(dto.getBookingReference())
                    .setUserId(dto.getUserId())
                    .setEventId(dto.getEventId())
                    .setTicketQuantity(dto.getSeatCount())
                    .addAllSeatNumbers(dto.getSeatNumbers() == null ? List.of() : dto.getSeatNumbers())
                    .setTotalAmount(amount)
                    .setCurrency(dto.getCurrency())
                    .setStatus(dto.getStatus())
                    .setPaymentStatus(dto.getPaymentStatus())
                    .setPaymentReference(dto.getPaymentReference() == null ? "" : dto.getPaymentReference())
                    .putAllMetadata(dto.getMetadata())
                    .build();
        }
    }

    @Override
    public void confirmBooking(BookingProto.ConfirmBookingRequest request,
            StreamObserver<BookingProto.ConfirmBookingResponse> responseObserver) {
        try {
            UUID bookingId = UUID.fromString(request.getBookingId());
            var result = bookingService.confirmBooking(bookingId, request.getPaymentReference());
            var booking = bookingService.getBooking(result.getBookingId());
            var response = BookingProto.ConfirmBookingResponse.newBuilder()
                    .setSuccess(true)
                    .setBooking(BookingMapperUtil.toProto(bookingMapper, booking))
                    .setMessage("Booking confirmed")
                    .build();
            responseObserver.onNext(response);
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (BookingNotFoundException e) {
            log.warn("Booking not found for confirmation: {}", request.getBookingId());
            responseObserver.onError(Status.NOT_FOUND
                    .withDescription(e.getMessage())
                    .asRuntimeException());
        } catch (IllegalArgumentException e) {
            log.warn("Invalid argument in confirmBooking: {}", e.getMessage());
            responseObserver.onError(Status.INVALID_ARGUMENT
                    .withDescription(e.getMessage())
                    .asRuntimeException());
        } catch (Exception e) {
            log.error("Failed to confirm booking: {}", request.getBookingId(), e);
            responseObserver.onError(Status.INTERNAL
                    .withDescription("Failed to confirm booking: " + e.getMessage())
                    .withCause(e)
                    .asRuntimeException());
        }
    }

    @Override
    public void cancelBooking(BookingProto.CancelBookingRequest request,
            StreamObserver<BookingProto.CancelBookingResponse> responseObserver) {
        try {
            bookingService.cancelBooking(
                    UUID.fromString(request.getBookingId()),
                    request.getReason());
            var response = BookingProto.CancelBookingResponse.newBuilder()
                    .setSuccess(true)
                    .setMessage("Booking cancelled")
                    .build();
            responseObserver.onNext(response);
            responseObserver.onCompleted();
        } catch (StatusRuntimeException e) {
            responseObserver.onError(e);
        } catch (BookingNotFoundException e) {
            log.warn("Booking not found for cancellation: {}", request.getBookingId());
            responseObserver.onError(Status.NOT_FOUND
                    .withDescription(e.getMessage())
                    .asRuntimeException());
        } catch (IllegalArgumentException e) {
            log.warn("Invalid argument in cancelBooking: {}", e.getMessage());
            responseObserver.onError(Status.INVALID_ARGUMENT
                    .withDescription(e.getMessage())
                    .asRuntimeException());
        } catch (Exception e) {
            log.error("Failed to cancel booking: {}", request.getBookingId(), e);
            responseObserver.onError(Status.INTERNAL
                    .withDescription("Failed to cancel booking: " + e.getMessage())
                    .withCause(e)
                    .asRuntimeException());
        }
    }

    @Override
    public void health(BookingProto.HealthRequest request,
            StreamObserver<BookingProto.HealthResponse> responseObserver) {
        var response = BookingProto.HealthResponse.newBuilder()
                .setStatus("UP")
                .setMessage("Booking service is healthy")
                .putDetails("service", "booking-service")
                .build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }
}
