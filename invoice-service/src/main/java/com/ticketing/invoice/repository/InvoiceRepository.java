package com.ticketing.invoice.repository;

import com.ticketing.invoice.domain.Invoice;
import com.ticketing.invoice.domain.InvoiceStatus;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface InvoiceRepository extends JpaRepository<Invoice, UUID> {

    Optional<Invoice> findByBookingId(UUID bookingId);

    Optional<Invoice> findByInvoiceNumber(String invoiceNumber);

    Page<Invoice> findByUserId(UUID userId, Pageable pageable);

    Page<Invoice> findByUserIdAndStatus(UUID userId, InvoiceStatus status, Pageable pageable);

    boolean existsByBookingId(UUID bookingId);
}
