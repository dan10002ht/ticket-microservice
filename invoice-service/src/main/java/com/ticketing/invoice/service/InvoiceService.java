package com.ticketing.invoice.service;

import com.ticketing.invoice.domain.Invoice;
import com.ticketing.invoice.domain.InvoiceItem;
import com.ticketing.invoice.domain.InvoiceStatus;
import com.ticketing.invoice.repository.InvoiceRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.time.Instant;
import java.time.YearMonth;
import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicInteger;

@Slf4j
@Service
@RequiredArgsConstructor
public class InvoiceService {

    private final InvoiceRepository invoiceRepository;
    private final PdfGenerationService pdfGenerationService;

    /**
     * Generates an invoice for a completed payment.
     * Idempotent: returns the existing invoice if one already exists for the booking.
     */
    @Transactional
    public Invoice generateInvoice(UUID bookingId, UUID paymentId, UUID userId, UUID eventId,
                                   BigDecimal totalAmount, String currency,
                                   List<InvoiceItemRequest> lineItems) {

        // Idempotency check
        Optional<Invoice> existing = invoiceRepository.findByBookingId(bookingId);
        if (existing.isPresent()) {
            log.info("Invoice already exists for booking {}: {}", bookingId, existing.get().getInvoiceNumber());
            return existing.get();
        }

        String invoiceNumber = generateInvoiceNumber();
        BigDecimal taxAmount = totalAmount.multiply(new BigDecimal("0.10")); // 10% VAT
        BigDecimal subtotal  = totalAmount.subtract(taxAmount);

        Invoice invoice = Invoice.builder()
                .invoiceNumber(invoiceNumber)
                .bookingId(bookingId)
                .paymentId(paymentId)
                .userId(userId)
                .eventId(eventId)
                .subtotal(subtotal)
                .taxAmount(taxAmount)
                .totalAmount(totalAmount)
                .currency(currency != null ? currency : "VND")
                .status(InvoiceStatus.PENDING)
                .build();

        // Add line items
        if (lineItems != null) {
            for (InvoiceItemRequest req : lineItems) {
                InvoiceItem item = InvoiceItem.builder()
                        .invoice(invoice)
                        .description(req.description())
                        .quantity(req.quantity())
                        .unitPrice(req.unitPrice())
                        .totalPrice(req.unitPrice().multiply(BigDecimal.valueOf(req.quantity())))
                        .build();
                invoice.getItems().add(item);
            }
        } else {
            // Default single line item
            InvoiceItem defaultItem = InvoiceItem.builder()
                    .invoice(invoice)
                    .description("Event Ticket")
                    .quantity(1)
                    .unitPrice(subtotal)
                    .totalPrice(subtotal)
                    .build();
            invoice.getItems().add(defaultItem);
        }

        Invoice saved = invoiceRepository.save(invoice);

        // Mark as generated
        saved.setStatus(InvoiceStatus.GENERATED);
        saved.setIssuedAt(Instant.now());
        saved = invoiceRepository.save(saved);

        log.info("Invoice {} generated for booking {}", invoiceNumber, bookingId);
        return saved;
    }

    @Transactional(readOnly = true)
    public Optional<Invoice> getByInvoiceId(UUID invoiceId) {
        return invoiceRepository.findById(invoiceId);
    }

    @Transactional(readOnly = true)
    public Optional<Invoice> getByBookingId(UUID bookingId) {
        return invoiceRepository.findByBookingId(bookingId);
    }

    @Transactional(readOnly = true)
    public Page<Invoice> listByUser(UUID userId, String status, int page, int limit) {
        PageRequest pageRequest = PageRequest.of(
                Math.max(0, page - 1), limit,
                Sort.by(Sort.Direction.DESC, "createdAt")
        );
        if (status != null && !status.isBlank()) {
            InvoiceStatus invoiceStatus = InvoiceStatus.valueOf(status.toUpperCase());
            return invoiceRepository.findByUserIdAndStatus(userId, invoiceStatus, pageRequest);
        }
        return invoiceRepository.findByUserId(userId, pageRequest);
    }

    public byte[] generatePdf(Invoice invoice) {
        return pdfGenerationService.generate(invoice);
    }

    private String generateInvoiceNumber() {
        YearMonth ym = YearMonth.now();
        long count = invoiceRepository.count() + 1;
        return String.format("INV-%d%02d-%06d", ym.getYear(), ym.getMonthValue(), count);
    }

    // -------------------------------------------------------------------------
    // Request record (inline DTO to avoid extra class file)
    // -------------------------------------------------------------------------
    public record InvoiceItemRequest(String description, int quantity, BigDecimal unitPrice) {}
}
