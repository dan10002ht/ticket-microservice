package com.ticketing.invoice.consumer;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.ticketing.invoice.service.InvoiceService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.support.Acknowledgment;
import org.springframework.kafka.support.KafkaHeaders;
import org.springframework.messaging.handler.annotation.Header;
import org.springframework.stereotype.Component;

import java.math.BigDecimal;
import java.util.UUID;

/**
 * Listens to the payment-events topic and generates invoices when a payment
 * is captured (PAYMENT_CAPTURED event).
 */
@Slf4j
@Component
@RequiredArgsConstructor
public class PaymentEventConsumer {

    private static final String EVENT_PAYMENT_CAPTURED = "PAYMENT_CAPTURED";

    private final InvoiceService invoiceService;
    private final ObjectMapper objectMapper;

    @KafkaListener(
            topics = "${kafka.topics.payment-events:payment-events}",
            groupId = "${spring.kafka.consumer.group-id:invoice-service}",
            containerFactory = "kafkaListenerContainerFactory"
    )
    public void onPaymentEvent(
            String message,
            @Header(KafkaHeaders.RECEIVED_PARTITION) int partition,
            @Header(KafkaHeaders.OFFSET) long offset,
            Acknowledgment ack
    ) {
        log.debug("Received payment event: partition={}, offset={}", partition, offset);

        try {
            JsonNode event = objectMapper.readTree(message);
            String eventType = event.path("eventType").asText();

            if (!EVENT_PAYMENT_CAPTURED.equals(eventType)) {
                ack.acknowledge();
                return;
            }

            processPaymentCaptured(event);
            ack.acknowledge();

        } catch (Exception e) {
            log.error("Failed to process payment event at offset {}: {}", offset, e.getMessage(), e);
            // Acknowledge to avoid infinite retry — dead-letter handling should be added
            // via Spring Kafka's DeadLetterPublishingRecoverer if needed.
            ack.acknowledge();
        }
    }

    private void processPaymentCaptured(JsonNode event) {
        UUID bookingId = UUID.fromString(event.path("bookingId").asText());
        UUID paymentId = parseUUIDOrNull(event.path("paymentId").asText(""));
        UUID userId    = parseUUIDOrNull(event.path("userId").asText(""));
        UUID eventId   = parseUUIDOrNull(event.path("eventId").asText(""));

        BigDecimal totalAmount = new BigDecimal(event.path("amount").asText("0"));
        String currency = event.path("currency").asText("VND");

        log.info("Generating invoice for booking {}", bookingId);
        invoiceService.generateInvoice(bookingId, paymentId, userId, eventId,
                totalAmount, currency, null);
    }

    private UUID parseUUIDOrNull(String value) {
        if (value == null || value.isBlank()) return null;
        try {
            return UUID.fromString(value);
        } catch (IllegalArgumentException e) {
            return null;
        }
    }
}
