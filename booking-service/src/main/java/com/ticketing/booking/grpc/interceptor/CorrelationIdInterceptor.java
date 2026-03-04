package com.ticketing.booking.grpc.interceptor;

import io.grpc.*;
import net.devh.boot.grpc.server.interceptor.GrpcGlobalServerInterceptor;
import org.slf4j.MDC;

/**
 * gRPC server interceptor that extracts the {@code correlation-id} metadata header
 * and binds it to the SLF4J MDC for the duration of the call.
 *
 * Registered globally via @GrpcGlobalServerInterceptor — no per-service wiring needed.
 */
@GrpcGlobalServerInterceptor
public class CorrelationIdInterceptor implements ServerInterceptor {

    private static final Metadata.Key<String> CORRELATION_ID_KEY =
            Metadata.Key.of("correlation-id", Metadata.ASCII_STRING_MARSHALLER);

    private static final String MDC_KEY = "correlationId";

    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
            ServerCall<ReqT, RespT> call,
            Metadata headers,
            ServerCallHandler<ReqT, RespT> next) {

        String correlationId = headers.get(CORRELATION_ID_KEY);
        if (correlationId == null || correlationId.isBlank()) {
            correlationId = "unknown";
        }

        MDC.put(MDC_KEY, correlationId);
        try {
            return new ForwardingServerCallListener.SimpleForwardingServerCallListener<>(
                    next.startCall(call, headers)) {
                @Override
                public void onComplete() {
                    MDC.remove(MDC_KEY);
                    super.onComplete();
                }

                @Override
                public void onCancel() {
                    MDC.remove(MDC_KEY);
                    super.onCancel();
                }
            };
        } finally {
            // MDC is cleared in onComplete/onCancel; no need to clear here
            // but ensure it's not leaked in synchronous error paths
        }
    }
}
