package com.ticketing.booking.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import com.ticketing.booking.grpcclient.PaymentServiceClient;
import com.ticketing.booking.grpcclient.TicketServiceClient;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Configuration
public class GrpcClientConfig {

    @Bean
    public ManagedChannel paymentServiceChannel(
            @Value("${grpc.payment-service.host:localhost}") String host,
            @Value("${grpc.payment-service.port:9090}") int port) {
        log.info("Creating gRPC channel to payment-service: {}:{}", host, port);
        return ManagedChannelBuilder.forAddress(host, port)
                .usePlaintext() // TODO: Use TLS in production
                .build();
    }

    @Bean
    public ManagedChannel ticketServiceChannel(
            @Value("${grpc.ticket-service.host:localhost}") String host,
            @Value("${grpc.ticket-service.port:50054}") int port) {
        log.info("Creating gRPC channel to ticket-service: {}:{}", host, port);
        return ManagedChannelBuilder.forAddress(host, port)
                .usePlaintext() // TODO: Use TLS in production
                .build();
    }

    @Bean
    public PaymentServiceClient paymentServiceClient(ManagedChannel paymentServiceChannel) {
        return new PaymentServiceClient(paymentServiceChannel);
    }

    @Bean
    public TicketServiceClient ticketServiceClient(ManagedChannel ticketServiceChannel) {
        return new TicketServiceClient(ticketServiceChannel);
    }
}
