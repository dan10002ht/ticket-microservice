package com.ticketing.booking.service;

import java.time.OffsetDateTime;
import java.util.List;

import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.ticketing.booking.entity.FailedCompensation;
import com.ticketing.booking.entity.FailedCompensation.CompensationStatus;
import com.ticketing.booking.entity.FailedCompensation.CompensationType;
import com.ticketing.booking.grpcclient.PaymentServiceClient;
import com.ticketing.booking.grpcclient.TicketServiceClient;
import com.ticketing.booking.metrics.BookingMetricsService;
import com.ticketing.booking.repository.FailedCompensationRepository;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

/**
 * Service to process failed compensations from the Dead Letter Queue.
 * Runs periodically to retry failed compensation operations.
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class CompensationRetryService {

    private final FailedCompensationRepository failedCompensationRepository;
    private final PaymentServiceClient paymentServiceClient;
    private final TicketServiceClient ticketServiceClient;
    private final BookingMetricsService metricsService;
    private final ObjectMapper objectMapper;

    /**
     * Save a failed compensation to DLQ for retry
     */
    @Transactional
    public void saveFailedCompensation(
            CompensationType type,
            String referenceId,
            Long bookingId,
            String bookingReference,
            String errorMessage,
            Object payload) {

        String payloadJson = null;
        try {
            if (payload != null) {
                payloadJson = objectMapper.writeValueAsString(payload);
            }
        } catch (Exception e) {
            log.warn("Failed to serialize compensation payload", e);
        }

        FailedCompensation fc = FailedCompensation.builder()
                .compensationType(type)
                .referenceId(referenceId)
                .bookingId(bookingId)
                .bookingReference(bookingReference)
                .errorMessage(errorMessage)
                .payload(payloadJson)
                .status(CompensationStatus.PENDING)
                .build();

        failedCompensationRepository.save(fc);

        log.warn("Saved failed compensation to DLQ: type={}, referenceId={}, bookingRef={}",
                type, referenceId, bookingReference);

        metricsService.recordSagaCompensation(type.name().toLowerCase() + "_dlq", "saved");
    }

    /**
     * Scheduled job to process failed compensations
     * Runs every 30 seconds
     */
    @Scheduled(fixedDelay = 30000, initialDelay = 60000)
    @Transactional
    public void processFailedCompensations() {
        List<FailedCompensation> pendingCompensations = failedCompensationRepository
                .findReadyForRetry(CompensationStatus.PENDING, OffsetDateTime.now());

        if (pendingCompensations.isEmpty()) {
            return;
        }

        log.info("Processing {} failed compensations", pendingCompensations.size());

        for (FailedCompensation fc : pendingCompensations) {
            processCompensation(fc);
        }
    }

    /**
     * Process a single failed compensation
     */
    private void processCompensation(FailedCompensation fc) {
        fc.setStatus(CompensationStatus.RETRYING);
        failedCompensationRepository.save(fc);

        try {
            switch (fc.getCompensationType()) {
                case CANCEL_PAYMENT -> retryPaymentCancellation(fc);
                case RELEASE_SEATS -> retrySeatRelease(fc);
                case REFUND_PAYMENT -> retryPaymentRefund(fc);
            }

            fc.markSucceeded();
            failedCompensationRepository.save(fc);

            log.info("Compensation succeeded: type={}, referenceId={}, retryCount={}",
                    fc.getCompensationType(), fc.getReferenceId(), fc.getRetryCount());

            metricsService.recordSagaCompensation(
                    fc.getCompensationType().name().toLowerCase() + "_retry", "success");

        } catch (Exception e) {
            log.error("Compensation retry failed: type={}, referenceId={}, attempt={}",
                    fc.getCompensationType(), fc.getReferenceId(), fc.getRetryCount() + 1, e);

            fc.setErrorMessage(e.getMessage());
            fc.scheduleNextRetry();
            failedCompensationRepository.save(fc);

            metricsService.recordSagaCompensation(
                    fc.getCompensationType().name().toLowerCase() + "_retry", "failed");

            // If max retries reached, alert
            if (fc.getStatus() == CompensationStatus.FAILED) {
                alertManualIntervention(fc);
            }
        }
    }

    private void retryPaymentCancellation(FailedCompensation fc) {
        if (fc.getReferenceId() == null) {
            throw new IllegalStateException("Payment ID is required for cancellation");
        }
        paymentServiceClient.cancelPayment(fc.getReferenceId());
    }

    private void retrySeatRelease(FailedCompensation fc) {
        if (fc.getReferenceId() == null) {
            throw new IllegalStateException("Reservation ID is required for seat release");
        }
        ticketServiceClient.releaseTickets(fc.getReferenceId(), null);
    }

    private void retryPaymentRefund(FailedCompensation fc) {
        if (fc.getReferenceId() == null) {
            throw new IllegalStateException("Payment ID is required for refund");
        }
        // TODO: Implement refund retry when refund service is available
        log.warn("Refund retry not implemented yet: {}", fc.getReferenceId());
        throw new UnsupportedOperationException("Refund retry not implemented");
    }

    /**
     * Alert for manual intervention when all retries exhausted
     */
    private void alertManualIntervention(FailedCompensation fc) {
        log.error("ALERT: Compensation requires manual intervention! " +
                        "type={}, referenceId={}, bookingRef={}, retryCount={}",
                fc.getCompensationType(),
                fc.getReferenceId(),
                fc.getBookingReference(),
                fc.getRetryCount());

        // TODO: Send alert to PagerDuty/Slack/Email
        // alertService.sendCriticalAlert("Compensation failed after max retries", fc);
    }

    /**
     * Get statistics for monitoring
     */
    public CompensationStats getStats() {
        return CompensationStats.builder()
                .pendingCount(failedCompensationRepository.countByStatus(CompensationStatus.PENDING))
                .failedCount(failedCompensationRepository.countByStatus(CompensationStatus.FAILED))
                .manualCount(failedCompensationRepository.countByStatus(CompensationStatus.MANUAL))
                .build();
    }

    @lombok.Builder
    @lombok.Value
    public static class CompensationStats {
        long pendingCount;
        long failedCount;
        long manualCount;
    }
}
