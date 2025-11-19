package com.ticketing.payment.webhook.stripe;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.validation.annotation.Validated;

import jakarta.validation.constraints.NotBlank;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@Validated
@ConfigurationProperties(prefix = "payment.stripe")
public class StripeProperties {

    @NotBlank
    private String apiKey;

    private String webhookSecret;
    private String publicKey;
}
