package com.ticketing.booking.metrics;

import org.springframework.stereotype.Service;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import io.micrometer.core.instrument.Tags;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

/**
 * Service for tracking booking-related metrics.
 * 
 * Exposes metrics to Prometheus for monitoring:
 * - Booking lifecycle events (created, confirmed, cancelled, failed)
 * - Saga orchestration metrics (duration, step counts)
 * - gRPC call metrics (success/failure counts)
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class BookingMetricsService {

    private final MeterRegistry meterRegistry;

    /**
     * Record booking created event
     */
    public void recordBookingCreated(String eventId) {
        Counter.builder("booking.created")
                .description("Total number of bookings created")
                .tags(Tags.of("service", "booking-service"))
                .register(meterRegistry)
                .increment();
        log.debug("Recorded booking.created metric for event: {}", eventId);
    }

    /**
     * Record booking confirmed event
     */
    public void recordBookingConfirmed(String eventId) {
        Counter.builder("booking.confirmed")
                .description("Total number of bookings confirmed")
                .tags(Tags.of("service", "booking-service"))
                .register(meterRegistry)
                .increment();
        log.debug("Recorded booking.confirmed metric for event: {}", eventId);
    }

    /**
     * Record booking cancelled event
     */
    public void recordBookingCancelled(String eventId, String reason) {
        Counter.builder("booking.cancelled")
                .description("Total number of bookings cancelled")
                .tags(Tags.of("service", "booking-service", "reason", reason != null ? reason : "unknown"))
                .register(meterRegistry)
                .increment();
        log.debug("Recorded booking.cancelled metric for event: {}, reason: {}", eventId, reason);
    }

    /**
     * Record booking failed event
     */
    public void recordBookingFailed(String eventId, String errorType) {
        Counter.builder("booking.failed")
                .description("Total number of bookings that failed")
                .tags(Tags.of("service", "booking-service", "error_type", errorType != null ? errorType : "unknown"))
                .register(meterRegistry)
                .increment();
        log.debug("Recorded booking.failed metric for event: {}, error: {}", eventId, errorType);
    }

    /**
     * Record saga execution duration
     */
    public Timer.Sample startSagaExecution() {
        return Timer.start(meterRegistry);
    }

    /**
     * Stop saga execution timer and record duration
     */
    public void recordSagaExecutionDuration(Timer.Sample sample, String status) {
        if (sample != null) {
            Timer timer = Timer.builder("saga.execution.duration")
                    .description("Duration of saga execution in seconds")
                    .tags(Tags.of("service", "booking-service", "status", status))
                    .register(meterRegistry);
            sample.stop(timer);
            log.debug("Recorded saga.execution.duration metric with status: {}", status);
        }
    }

    /**
     * Record saga step execution
     */
    public void recordSagaStep(String stepName, String status) {
        Counter.builder("saga.step.executed")
                .description("Total number of saga steps executed")
                .tags(Tags.of("service", "booking-service", "step", stepName, "status", status))
                .register(meterRegistry)
                .increment();
        log.debug("Recorded saga.step.executed metric: step={}, status={}", stepName, status);
    }

    /**
     * Record saga compensation triggered
     */
    public void recordSagaCompensation(String stepName, String reason) {
        Counter.builder("saga.compensation.triggered")
                .description("Total number of saga compensations triggered")
                .tags(Tags.of("service", "booking-service", "step", stepName, "reason",
                        reason != null ? reason : "unknown"))
                .register(meterRegistry)
                .increment();
        log.debug("Recorded saga.compensation.triggered metric: step={}, reason={}", stepName, reason);
    }

    /**
     * Record gRPC call
     */
    public void recordGrpcCall(String service, String method, String status) {
        Counter.builder("grpc.calls.total")
                .description("Total number of gRPC calls")
                .tags(Tags.of("service", "booking-service", "target_service", service, "method", method, "status",
                        status))
                .register(meterRegistry)
                .increment();
        log.debug("Recorded grpc.calls.total metric: service={}, method={}, status={}", service, method, status);
    }

    /**
     * Record gRPC error
     */
    public void recordGrpcError(String service, String method, String errorCode) {
        Counter.builder("grpc.calls.errors")
                .description("Total number of gRPC call errors")
                .tags(Tags.of("service", "booking-service", "target_service", service, "method", method, "error_code",
                        errorCode != null ? errorCode : "unknown"))
                .register(meterRegistry)
                .increment();
        log.debug("Recorded grpc.calls.errors metric: service={}, method={}, error_code={}", service, method, errorCode);
    }
}

