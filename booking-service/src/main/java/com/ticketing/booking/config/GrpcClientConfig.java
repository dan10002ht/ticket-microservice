package com.ticketing.booking.config;

import java.io.File;

import javax.net.ssl.SSLException;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import com.ticketing.booking.grpcclient.PaymentServiceClient;
import com.ticketing.booking.grpcclient.TicketServiceClient;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.netty.shaded.io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.shaded.io.grpc.netty.NettyChannelBuilder;
import io.grpc.netty.shaded.io.netty.handler.ssl.SslContext;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Configuration
public class GrpcClientConfig {

    @Value("${GRPC_TLS_ENABLED:false}")
    private boolean tlsEnabled;

    @Value("${GRPC_TLS_CERT:/certs/server.crt}")
    private String tlsCert;

    @Value("${GRPC_TLS_KEY:/certs/server.key}")
    private String tlsKey;

    @Value("${GRPC_TLS_CA:/certs/ca.crt}")
    private String tlsCa;

    @Bean
    public ManagedChannel paymentServiceChannel(
            @Value("${grpc.payment-service.host:localhost}") String host,
            @Value("${grpc.payment-service.port:50062}") int port) throws SSLException {
        log.info("Creating gRPC channel to payment-service: {}:{} (tls={})", host, port, tlsEnabled);
        return createChannel(host, port);
    }

    @Bean
    public ManagedChannel ticketServiceChannel(
            @Value("${grpc.ticket-service.host:localhost}") String host,
            @Value("${grpc.ticket-service.port:50053}") int port) throws SSLException {
        log.info("Creating gRPC channel to ticket-service: {}:{} (tls={})", host, port, tlsEnabled);
        return createChannel(host, port);
    }

    @Bean
    public PaymentServiceClient paymentServiceClient(ManagedChannel paymentServiceChannel) {
        return new PaymentServiceClient(paymentServiceChannel);
    }

    @Bean
    public TicketServiceClient ticketServiceClient(ManagedChannel ticketServiceChannel) {
        return new TicketServiceClient(ticketServiceChannel);
    }

    private ManagedChannel createChannel(String host, int port) throws SSLException {
        if (!tlsEnabled) {
            return ManagedChannelBuilder.forAddress(host, port)
                    .usePlaintext()
                    .build();
        }

        SslContext sslContext = GrpcSslContexts.forClient()
                .trustManager(new File(tlsCa))
                .keyManager(new File(tlsCert), new File(tlsKey))
                .build();

        return NettyChannelBuilder.forAddress(host, port)
                .sslContext(sslContext)
                .build();
    }
}
