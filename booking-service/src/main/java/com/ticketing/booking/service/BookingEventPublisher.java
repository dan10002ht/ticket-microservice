package com.ticketing.booking.service;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.ticketing.booking.entity.Booking;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Component
@RequiredArgsConstructor
@Slf4j
public class BookingEventPublisher {

    private final KafkaTemplate<String, String> kafkaTemplate;
    private final ObjectMapper objectMapper;

    @Value("${kafka.topics.booking-events:booking-events}")
    private String bookingTopic;

    public void publishBookingCreated(Booking booking) {
        publish("BOOKING_CREATED", booking);
    }

    public void publishBookingConfirmed(Booking booking) {
        publish("BOOKING_CONFIRMED", booking);
    }

    public void publishBookingCancelled(Booking booking) {
        publish("BOOKING_CANCELLED", booking);
    }

    private void publish(String eventType, Booking booking) {
        try {
            BookingEventPayload payload = new BookingEventPayload(eventType,
                    booking.getBookingId().toString(),
                    booking.getUserId(),
                    booking.getEventId());
            kafkaTemplate.send(bookingTopic, booking.getBookingReference(), objectMapper.writeValueAsString(payload));
            log.info("Booking event published {} {}", eventType, booking.getBookingReference());
        } catch (JsonProcessingException ex) {
            log.error("Failed to serialize booking event {}", booking.getBookingReference(), ex);
        }
    }

    private record BookingEventPayload(String type, String bookingId, String userId, String eventId) {
    }
}

