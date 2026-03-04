package com.ticketing.invoice;

import com.ticketing.invoice.domain.Invoice;
import com.ticketing.invoice.domain.InvoiceStatus;
import com.ticketing.invoice.repository.InvoiceRepository;
import com.ticketing.invoice.service.InvoiceService;
import com.ticketing.invoice.service.PdfGenerationService;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.math.BigDecimal;
import java.util.Optional;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class InvoiceServiceTest {

    @Mock
    private InvoiceRepository invoiceRepository;

    @Mock
    private PdfGenerationService pdfGenerationService;

    @InjectMocks
    private InvoiceService invoiceService;

    @Test
    void generateInvoice_idempotent_returnExistingIfAlreadyPresent() {
        UUID bookingId = UUID.randomUUID();
        Invoice existing = Invoice.builder()
                .id(UUID.randomUUID())
                .invoiceNumber("INV-202412-000001")
                .bookingId(bookingId)
                .status(InvoiceStatus.GENERATED)
                .build();

        when(invoiceRepository.findByBookingId(bookingId)).thenReturn(Optional.of(existing));

        Invoice result = invoiceService.generateInvoice(
                bookingId, null, UUID.randomUUID(), null,
                new BigDecimal("100000"), "VND", null
        );

        assertThat(result.getInvoiceNumber()).isEqualTo("INV-202412-000001");
        verify(invoiceRepository, never()).save(any());
    }

    @Test
    void generateInvoice_createsNewInvoiceWhenNotExists() {
        UUID bookingId = UUID.randomUUID();
        when(invoiceRepository.findByBookingId(bookingId)).thenReturn(Optional.empty());
        when(invoiceRepository.count()).thenReturn(0L);

        Invoice saved = Invoice.builder()
                .id(UUID.randomUUID())
                .invoiceNumber("INV-202412-000001")
                .bookingId(bookingId)
                .status(InvoiceStatus.GENERATED)
                .subtotal(new BigDecimal("90000"))
                .taxAmount(new BigDecimal("10000"))
                .totalAmount(new BigDecimal("100000"))
                .build();

        when(invoiceRepository.save(any(Invoice.class))).thenReturn(saved);

        Invoice result = invoiceService.generateInvoice(
                bookingId, null, UUID.randomUUID(), null,
                new BigDecimal("100000"), "VND", null
        );

        assertThat(result).isNotNull();
        verify(invoiceRepository, atLeastOnce()).save(any(Invoice.class));
    }
}
