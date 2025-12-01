package com.ticketing.booking.config;

import org.springframework.context.annotation.Configuration;
import org.springframework.retry.annotation.EnableRetry;

/**
 * Configuration for gRPC retry logic.
 * 
 * Enables Spring Retry for handling transient failures in gRPC calls.
 * Retry policies are configured via @Retryable annotations on gRPC client methods.
 */
@Configuration
@EnableRetry
public class GrpcRetryConfig {
    // Retry configuration is done via @Retryable annotations
    // and application.yml properties
}

