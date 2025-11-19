package com.ticketing.payment;

import java.util.List;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.data.jpa.repository.config.EnableJpaAuditing;
import org.springframework.scheduling.annotation.EnableAsync;
import org.springframework.scheduling.annotation.EnableScheduling;

import com.ticketing.payment.adapter.PaymentGatewayAdapter;
import com.ticketing.payment.adapter.PaymentGatewayRegistry;
import com.ticketing.payment.adapter.stripe.StripeProperties;
import com.ticketing.payment.webhook.PaymentWebhookHandler;
import com.ticketing.payment.webhook.WebhookHandlerRegistry;

/**
 * Payment Service Application
 * 
 * Handles all payment processing, refunds, and webhook integrations
 * for the ticket booking system.
 */
@SpringBootApplication
@EnableJpaAuditing
@EnableAsync
@EnableScheduling
@EnableConfigurationProperties({ StripeProperties.class, com.ticketing.payment.webhook.stripe.StripeProperties.class })
public class PaymentServiceApplication {

    public static void main(String[] args) {
        SpringApplication.run(PaymentServiceApplication.class, args);
    }

    @Bean
    public PaymentGatewayRegistry paymentGatewayRegistry(List<PaymentGatewayAdapter> adapters) {
        PaymentGatewayRegistry registry = new PaymentGatewayRegistry();
        adapters.forEach(adapter -> registry.register(adapter.getGateway(), adapter));
        return registry;
    }

    @Bean
    public WebhookHandlerRegistry webhookHandlerRegistry(List<PaymentWebhookHandler> handlers) {
        WebhookHandlerRegistry registry = new WebhookHandlerRegistry();
        handlers.forEach(registry::register);
        return registry;
    }
}
