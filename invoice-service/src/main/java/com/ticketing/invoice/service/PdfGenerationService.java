package com.ticketing.invoice.service;

import com.itextpdf.io.font.constants.StandardFonts;
import com.itextpdf.kernel.colors.DeviceRgb;
import com.itextpdf.kernel.font.PdfFontFactory;
import com.itextpdf.kernel.pdf.PdfDocument;
import com.itextpdf.kernel.pdf.PdfWriter;
import com.itextpdf.layout.Document;
import com.itextpdf.layout.element.Cell;
import com.itextpdf.layout.element.Paragraph;
import com.itextpdf.layout.element.Table;
import com.itextpdf.layout.properties.TextAlignment;
import com.itextpdf.layout.properties.UnitValue;
import com.ticketing.invoice.domain.Invoice;
import com.ticketing.invoice.domain.InvoiceItem;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.math.BigDecimal;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;

@Slf4j
@Service
public class PdfGenerationService {

    private static final DateTimeFormatter DATE_FMT =
            DateTimeFormatter.ofPattern("dd/MM/yyyy HH:mm").withZone(ZoneId.of("Asia/Ho_Chi_Minh"));

    private static final DeviceRgb HEADER_COLOR  = new DeviceRgb(41, 128, 185);
    private static final DeviceRgb ROW_ALT_COLOR = new DeviceRgb(240, 248, 255);

    /**
     * Generates a PDF invoice and returns the raw bytes.
     */
    public byte[] generate(Invoice invoice) {
        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        try (PdfDocument pdf = new PdfDocument(new PdfWriter(baos));
             Document doc = new Document(pdf)) {

            addTitle(doc, invoice);
            addMetaTable(doc, invoice);
            addItemsTable(doc, invoice);
            addTotals(doc, invoice);
            addFooter(doc);

        } catch (IOException e) {
            log.error("PDF generation failed for invoice {}", invoice.getInvoiceNumber(), e);
            throw new RuntimeException("PDF generation failed: " + e.getMessage(), e);
        }
        return baos.toByteArray();
    }

    private void addTitle(Document doc, Invoice invoice) throws IOException {
        doc.add(new Paragraph("INVOICE")
                .setFont(PdfFontFactory.createFont(StandardFonts.HELVETICA_BOLD))
                .setFontSize(24)
                .setFontColor(HEADER_COLOR)
                .setTextAlignment(TextAlignment.CENTER)
                .setMarginBottom(4));

        doc.add(new Paragraph(invoice.getInvoiceNumber())
                .setFont(PdfFontFactory.createFont(StandardFonts.HELVETICA))
                .setFontSize(12)
                .setTextAlignment(TextAlignment.CENTER)
                .setMarginBottom(20));
    }

    private void addMetaTable(Document doc, Invoice invoice) throws IOException {
        Table meta = new Table(UnitValue.createPercentArray(new float[]{50, 50})).useAllAvailableWidth();

        String issuedAt = invoice.getIssuedAt() != null
                ? DATE_FMT.format(invoice.getIssuedAt())
                : DATE_FMT.format(invoice.getCreatedAt());

        meta.addCell(labelCell("Invoice Number:")).addCell(valueCell(invoice.getInvoiceNumber()));
        meta.addCell(labelCell("Booking ID:")).addCell(valueCell(invoice.getBookingId().toString()));
        meta.addCell(labelCell("User ID:")).addCell(valueCell(invoice.getUserId().toString()));
        meta.addCell(labelCell("Date Issued:")).addCell(valueCell(issuedAt));
        meta.addCell(labelCell("Status:")).addCell(valueCell(invoice.getStatus().name()));
        meta.addCell(labelCell("Currency:")).addCell(valueCell(invoice.getCurrency()));

        doc.add(meta.setMarginBottom(20));
    }

    private void addItemsTable(Document doc, Invoice invoice) throws IOException {
        float[] colWidths = {10f, 50f, 15f, 25f, 25f};
        Table table = new Table(UnitValue.createPercentArray(colWidths)).useAllAvailableWidth();

        // Header row
        for (String h : new String[]{"#", "Description", "Qty", "Unit Price", "Total"}) {
            table.addHeaderCell(
                    new Cell().add(new Paragraph(h)
                            .setFont(PdfFontFactory.createFont(StandardFonts.HELVETICA_BOLD))
                            .setFontColor(new DeviceRgb(255, 255, 255)))
                            .setBackgroundColor(HEADER_COLOR));
        }

        int row = 0;
        for (InvoiceItem item : invoice.getItems()) {
            DeviceRgb bg = (row++ % 2 == 0) ? null : ROW_ALT_COLOR;
            addItemRow(table, row, item, bg);
        }

        doc.add(table.setMarginBottom(10));
    }

    private void addItemRow(Table table, int idx, InvoiceItem item, DeviceRgb bg) throws IOException {
        String[] values = {
                String.valueOf(idx),
                item.getDescription(),
                String.valueOf(item.getQuantity()),
                formatAmount(item.getUnitPrice(), ""),
                formatAmount(item.getTotalPrice(), "")
        };
        for (String v : values) {
            Cell cell = new Cell().add(new Paragraph(v)
                    .setFont(PdfFontFactory.createFont(StandardFonts.HELVETICA)));
            if (bg != null) cell.setBackgroundColor(bg);
            table.addCell(cell);
        }
    }

    private void addTotals(Document doc, Invoice invoice) throws IOException {
        Table totals = new Table(UnitValue.createPercentArray(new float[]{70, 30})).useAllAvailableWidth();

        totals.addCell(new Cell().add(new Paragraph("Subtotal:")
                .setFont(PdfFontFactory.createFont(StandardFonts.HELVETICA))
                .setTextAlignment(TextAlignment.RIGHT)).setBorder(null));
        totals.addCell(new Cell().add(new Paragraph(formatAmount(invoice.getSubtotal(), invoice.getCurrency()))
                .setTextAlignment(TextAlignment.RIGHT)).setBorder(null));

        totals.addCell(new Cell().add(new Paragraph("Tax:")
                .setFont(PdfFontFactory.createFont(StandardFonts.HELVETICA))
                .setTextAlignment(TextAlignment.RIGHT)).setBorder(null));
        totals.addCell(new Cell().add(new Paragraph(formatAmount(invoice.getTaxAmount(), invoice.getCurrency()))
                .setTextAlignment(TextAlignment.RIGHT)).setBorder(null));

        totals.addCell(new Cell().add(new Paragraph("Total:")
                .setFont(PdfFontFactory.createFont(StandardFonts.HELVETICA_BOLD))
                .setFontColor(HEADER_COLOR)
                .setTextAlignment(TextAlignment.RIGHT)).setBorder(null));
        totals.addCell(new Cell().add(new Paragraph(formatAmount(invoice.getTotalAmount(), invoice.getCurrency()))
                .setFont(PdfFontFactory.createFont(StandardFonts.HELVETICA_BOLD))
                .setFontColor(HEADER_COLOR)
                .setTextAlignment(TextAlignment.RIGHT)).setBorder(null));

        doc.add(totals.setMarginBottom(30));
    }

    private void addFooter(Document doc) throws IOException {
        doc.add(new Paragraph("Thank you for your purchase!")
                .setFont(PdfFontFactory.createFont(StandardFonts.HELVETICA_OBLIQUE))
                .setFontSize(10)
                .setTextAlignment(TextAlignment.CENTER)
                .setFontColor(new DeviceRgb(128, 128, 128)));
    }

    private Cell labelCell(String text) throws IOException {
        return new Cell().add(new Paragraph(text)
                .setFont(PdfFontFactory.createFont(StandardFonts.HELVETICA_BOLD)));
    }

    private Cell valueCell(String text) throws IOException {
        return new Cell().add(new Paragraph(text)
                .setFont(PdfFontFactory.createFont(StandardFonts.HELVETICA)));
    }

    private String formatAmount(BigDecimal amount, String currency) {
        if (amount == null) return "0";
        String formatted = String.format("%,.0f", amount);
        return currency.isEmpty() ? formatted : formatted + " " + currency;
    }
}
