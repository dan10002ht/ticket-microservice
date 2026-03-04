package com.ticketing.invoice.domain;

public enum InvoiceStatus {
    PENDING,      // event received, invoice not yet generated
    GENERATED,    // PDF created and invoice is ready
    CANCELLED     // booking was refunded / cancelled after invoice was generated
}
