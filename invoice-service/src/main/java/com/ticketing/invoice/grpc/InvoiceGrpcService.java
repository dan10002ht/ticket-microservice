package com.ticketing.invoice.grpc;

import com.google.protobuf.ByteString;
import com.ticketing.invoice.domain.Invoice;
import com.ticketing.invoice.domain.InvoiceItem;
import com.ticketing.invoice.grpc.proto.GetInvoicePdfRequest;
import com.ticketing.invoice.grpc.proto.GetInvoicePdfResponse;
import com.ticketing.invoice.grpc.proto.GetInvoiceRequest;
import com.ticketing.invoice.grpc.proto.GetInvoiceResponse;
import com.ticketing.invoice.grpc.proto.HealthRequest;
import com.ticketing.invoice.grpc.proto.HealthResponse;
import com.ticketing.invoice.grpc.proto.InvoiceServiceGrpc;
import com.ticketing.invoice.grpc.proto.ListInvoicesRequest;
import com.ticketing.invoice.grpc.proto.ListInvoicesResponse;
import com.ticketing.invoice.service.InvoiceService;
import com.ticketing.invoice.service.PdfGenerationService;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.server.service.GrpcService;
import org.springframework.data.domain.Page;

import java.util.Optional;
import java.util.UUID;

@Slf4j
@GrpcService
@RequiredArgsConstructor
public class InvoiceGrpcService extends InvoiceServiceGrpc.InvoiceServiceImplBase {

    private final InvoiceService invoiceService;
    private final PdfGenerationService pdfGenerationService;

    @Override
    public void getInvoice(GetInvoiceRequest req,
                           StreamObserver<GetInvoiceResponse> obs) {
        try {
            Optional<Invoice> invoice = req.getInvoiceId().isEmpty()
                    ? invoiceService.getByBookingId(UUID.fromString(req.getBookingId()))
                    : invoiceService.getByInvoiceId(UUID.fromString(req.getInvoiceId()));

            if (invoice.isEmpty()) {
                obs.onError(Status.NOT_FOUND
                        .withDescription("Invoice not found").asRuntimeException());
                return;
            }
            obs.onNext(GetInvoiceResponse.newBuilder()
                    .setSuccess(true)
                    .setInvoice(toProto(invoice.get()))
                    .build());
            obs.onCompleted();
        } catch (Exception e) {
            log.error("getInvoice failed", e);
            obs.onError(Status.INTERNAL.withDescription(e.getMessage()).asRuntimeException());
        }
    }

    @Override
    public void listInvoices(ListInvoicesRequest req,
                             StreamObserver<ListInvoicesResponse> obs) {
        try {
            int page = req.getPage() <= 0 ? 0 : req.getPage() - 1; // proto is 1-based
            int limit = req.getLimit() <= 0 ? 20 : Math.min(req.getLimit(), 200);

            Page<Invoice> result = invoiceService.listByUser(
                    UUID.fromString(req.getUserId()),
                    req.getStatus().isEmpty() ? null : req.getStatus(),
                    page,
                    limit);

            ListInvoicesResponse.Builder builder = ListInvoicesResponse.newBuilder()
                    .setTotal((int) result.getTotalElements())
                    .setPage(req.getPage() <= 0 ? 1 : req.getPage())
                    .setLimit(limit)
                    .setHasMore(result.hasNext());
            result.getContent().forEach(inv -> builder.addInvoices(toProto(inv)));

            obs.onNext(builder.build());
            obs.onCompleted();
        } catch (Exception e) {
            log.error("listInvoices failed", e);
            obs.onError(Status.INTERNAL.withDescription(e.getMessage()).asRuntimeException());
        }
    }

    @Override
    public void getInvoicePdf(GetInvoicePdfRequest req,
                              StreamObserver<GetInvoicePdfResponse> obs) {
        try {
            Optional<Invoice> invoice = invoiceService.getByInvoiceId(
                    UUID.fromString(req.getInvoiceId()));

            if (invoice.isEmpty()) {
                obs.onError(Status.NOT_FOUND
                        .withDescription("Invoice not found").asRuntimeException());
                return;
            }
            byte[] pdf = pdfGenerationService.generate(invoice.get());
            obs.onNext(GetInvoicePdfResponse.newBuilder()
                    .setSuccess(true)
                    .setPdfBytes(ByteString.copyFrom(pdf))
                    .setContentType("application/pdf")
                    .setFilename(invoice.get().getInvoiceNumber() + ".pdf")
                    .build());
            obs.onCompleted();
        } catch (Exception e) {
            log.error("getInvoicePdf failed", e);
            obs.onError(Status.INTERNAL.withDescription(e.getMessage()).asRuntimeException());
        }
    }

    @Override
    public void health(HealthRequest req, StreamObserver<HealthResponse> obs) {
        obs.onNext(HealthResponse.newBuilder()
                .setStatus("UP")
                .setMessage("invoice-service is healthy")
                .build());
        obs.onCompleted();
    }

    // -----------------------------------------------------------------------
    // Mapping helpers
    // -----------------------------------------------------------------------

    private com.ticketing.invoice.grpc.proto.Invoice toProto(Invoice inv) {
        var builder = com.ticketing.invoice.grpc.proto.Invoice.newBuilder()
                .setId(inv.getId().toString())
                .setInvoiceNumber(inv.getInvoiceNumber())
                .setBookingId(inv.getBookingId().toString())
                .setUserId(inv.getUserId().toString())
                .setSubtotal(inv.getSubtotal().doubleValue())
                .setTaxAmount(inv.getTaxAmount().doubleValue())
                .setTotalAmount(inv.getTotalAmount().doubleValue())
                .setCurrency(inv.getCurrency())
                .setStatus(inv.getStatus().name());

        if (inv.getPaymentId() != null) {
            builder.setPaymentId(inv.getPaymentId().toString());
        }
        if (inv.getEventId() != null) {
            builder.setEventId(inv.getEventId().toString());
        }
        if (inv.getIssuedAt() != null) {
            builder.setIssuedAt(inv.getIssuedAt().toEpochMilli());
        }
        if (inv.getCreatedAt() != null) {
            builder.setCreatedAt(inv.getCreatedAt().toEpochMilli());
        }
        for (InvoiceItem item : inv.getItems()) {
            builder.addItems(toItemProto(item));
        }
        return builder.build();
    }

    private com.ticketing.invoice.grpc.proto.InvoiceItem toItemProto(InvoiceItem item) {
        return com.ticketing.invoice.grpc.proto.InvoiceItem.newBuilder()
                .setId(item.getId().toString())
                .setDescription(item.getDescription())
                .setQuantity(item.getQuantity())
                .setUnitPrice(item.getUnitPrice().doubleValue())
                .setTotalPrice(item.getTotalPrice().doubleValue())
                .build();
    }
}
