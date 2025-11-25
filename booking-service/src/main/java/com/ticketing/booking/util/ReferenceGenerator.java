package com.ticketing.booking.util;

import java.time.OffsetDateTime;
import java.time.temporal.ChronoUnit;
import java.util.Locale;
import java.util.UUID;

public final class ReferenceGenerator {

    private ReferenceGenerator() {
    }

    public static String bookingReference() {
        return "BKG-" + UUID.randomUUID().toString()
                .replace("-", "")
                .substring(0, 12)
                .toUpperCase(Locale.ROOT);
    }

    public static OffsetDateTime bookingExpiry() {
        return OffsetDateTime.now().plus(15, ChronoUnit.MINUTES);
    }
}

