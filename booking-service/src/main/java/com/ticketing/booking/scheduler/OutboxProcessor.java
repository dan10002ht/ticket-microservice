package com.ticketing.booking.scheduler;

import java.time.OffsetDateTime;
import java.util.List;
import java.util.concurrent.TimeUnit;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import com.ticketing.booking.entity.OutboxEvent;
import com.ticketing.booking.entity.OutboxEvent.OutboxStatus;
import com.ticketing.booking.repository.OutboxEventRepository;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.Gauge;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import lombok.extern.slf4j.Slf4j;

/**
 * Background processor for the Transactional Outbox Pattern.
 *
 * Polls the outbox_events table and publishes events to Kafka.
 * Uses SELECT FOR UPDATE SKIP LOCKED for concurrent processing support.
 *
 * Features:
 * - Batch processing for efficiency
 * - Automatic retry of failed events
 * - Cleanup of old published events
 * - Metrics for monitoring
 */
@Component
@Slf4j
public class OutboxProcessor {

    private final OutboxEventRepository outboxEventRepository;
    private final KafkaTemplate<String, String> kafkaTemplate;
    private final Counter publishedCounter;
    private final Counter failedCounter;
    private final Timer processingTimer;

    @Value("${outbox.processor.batch-size:100}")
    private int batchSize;

    @Value("${outbox.processor.retention-days:7}")
    private int retentionDays;

    @Value("${kafka.topics.booking-events:booking-events}")
    private String bookingEventsTopic;

    public OutboxProcessor(
            OutboxEventRepository outboxEventRepository,
            KafkaTemplate<String, String> kafkaTemplate,
            MeterRegistry meterRegistry) {
        this.outboxEventRepository = outboxEventRepository;
        this.kafkaTemplate = kafkaTemplate;

        // Metrics
        this.publishedCounter = meterRegistry.counter("outbox.events.published");
        this.failedCounter = meterRegistry.counter("outbox.events.failed");
        this.processingTimer = meterRegistry.timer("outbox.processing.duration");

        // Gauge for pending events
        Gauge.builder("outbox.events.pending", outboxEventRepository,
                repo -> repo.countByStatus(OutboxStatus.PENDING))
                .register(meterRegistry);

        // Gauge for failed events
        Gauge.builder("outbox.events.failed.total", outboxEventRepository,
                repo -> repo.countByStatus(OutboxStatus.FAILED))
                .register(meterRegistry);
    }

    /**
     * Main processing loop - runs every 100ms
     * Polls pending events and publishes to Kafka
     */
    @Scheduled(fixedDelayString = "${outbox.processor.poll-interval-ms:100}")
    @Transactional
    public void processOutboxEvents() {
        Timer.Sample sample = Timer.start();

        try {
            List<OutboxEvent> events = outboxEventRepository.findPendingEventsForProcessing(batchSize);

            if (events.isEmpty()) {
                return;
            }

            log.debug("Processing {} outbox events", events.size());

            for (OutboxEvent event : events) {
                processEvent(event);
            }

        } catch (Exception e) {
            log.error("Error in outbox processor", e);
        } finally {
            sample.stop(processingTimer);
        }
    }

    /**
     * Process a single outbox event
     */
    private void processEvent(OutboxEvent event) {
        try {
            String topic = resolveTopicForEvent(event);
            String key = event.getAggregateId();
            String payload = event.getPayload();

            // Synchronous send to ensure delivery before marking as published
            kafkaTemplate.send(topic, key, payload).get(5, TimeUnit.SECONDS);

            // Mark as published
            event.markPublished();
            outboxEventRepository.save(event);
            publishedCounter.increment();

            log.debug("Published outbox event: id={}, type={}, topic={}",
                    event.getId(), event.getEventType(), topic);

        } catch (Exception e) {
            log.error("Failed to publish outbox event: id={}, type={}, error={}",
                    event.getId(), event.getEventType(), e.getMessage());

            // Mark as failed
            event.markFailed(e.getMessage());
            outboxEventRepository.save(event);
            failedCounter.increment();
        }
    }

    /**
     * Resolve Kafka topic based on event type
     */
    private String resolveTopicForEvent(OutboxEvent event) {
        // For now, all booking events go to booking-events topic
        // Can be extended for different aggregate types
        return switch (event.getAggregateType()) {
            case BOOKING -> bookingEventsTopic;
            case PAYMENT -> "payment-events";
            case TICKET -> "ticket-events";
        };
    }

    /**
     * Retry failed events - runs every 30 seconds
     * Resets failed events with retry count < 5 back to PENDING
     */
    @Scheduled(fixedDelayString = "${outbox.processor.retry-interval-ms:30000}")
    @Transactional
    public void retryFailedEvents() {
        try {
            int resetCount = outboxEventRepository.resetFailedEventsForRetry();
            if (resetCount > 0) {
                log.info("Reset {} failed outbox events for retry", resetCount);
            }
        } catch (Exception e) {
            log.error("Error retrying failed outbox events", e);
        }
    }

    /**
     * Cleanup old published events - runs every hour
     * Deletes events older than retention period
     */
    @Scheduled(fixedDelayString = "${outbox.processor.cleanup-interval-ms:3600000}")
    @Transactional
    public void cleanupOldEvents() {
        try {
            OffsetDateTime cutoff = OffsetDateTime.now().minusDays(retentionDays);
            int deletedCount = outboxEventRepository.deletePublishedEventsBefore(cutoff);

            if (deletedCount > 0) {
                log.info("Cleaned up {} old published outbox events (older than {} days)",
                        deletedCount, retentionDays);
            }
        } catch (Exception e) {
            log.error("Error cleaning up old outbox events", e);
        }
    }

    /**
     * Health check - returns true if outbox is healthy
     * Unhealthy if oldest pending event is older than 1 minute
     */
    public boolean isHealthy() {
        OffsetDateTime oldestPending = outboxEventRepository.findOldestPendingEventTime();
        if (oldestPending == null) {
            return true; // No pending events
        }
        // Unhealthy if events are stuck for more than 1 minute
        return oldestPending.isAfter(OffsetDateTime.now().minusMinutes(1));
    }
}
