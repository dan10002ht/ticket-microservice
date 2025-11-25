package com.ticketing.booking.grpc;

import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;

import org.lognet.springboot.grpc.GRpcService;

import com.ticketing.booking.grpc.BookingProto.BookingServiceGrpc.BookingServiceImplBase;
import com.ticketing.booking.service.BookingService;
import com.ticketing.booking.service.dto.BookingCreateCommand;
import com.ticketing.booking.service.mapper.BookingMapper;

import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;

@GRpcService
@RequiredArgsConstructor
public class BookingGrpcService extends BookingServiceImplBase {

    private final BookingService bookingService;
    private final BookingMapper bookingMapper;

    @Override
    public void createBooking(BookingProto.CreateBookingRequest request,
            StreamObserver<BookingProto.BookingResponse> responseObserver) {
        BookingCreateCommand command = BookingCreateCommand.builder()
                .userId(request.getUserId())
                .eventId(request.getEventId())
                .seatNumbers(request.getSeatNumbersList())
                .seatCount(request.getTicketQuantity())
                .currency("USD")
                .totalAmount(BigDecimal.valueOf(Math.max(request.getTicketQuantity(), 1)))
                .metadata(request.getMetadataMap())
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
    }

    @Override
    public void getBooking(BookingProto.GetBookingRequest request,
            StreamObserver<BookingProto.BookingResponse> responseObserver) {
        var booking = bookingService.getBooking(UUID.fromString(request.getBookingId()));
        responseObserver.onNext(BookingProto.BookingResponse.newBuilder()
                .setSuccess(true)
                .setBooking(BookingMapperUtil.toProto(bookingMapper, booking))
                .build());
        responseObserver.onCompleted();
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
}
package com.ticketing.booking.grpc;


