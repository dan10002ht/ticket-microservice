package interceptors

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const correlationIDKey = "correlation-id"

type correlationIDCtxKey struct{}

// CorrelationServerInterceptor extracts correlation-id from incoming gRPC metadata,
// stores it in context, and logs each call.
func CorrelationServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		correlationID := "unknown"
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get(correlationIDKey); len(vals) > 0 {
				correlationID = vals[0]
			}
		}
		logger.Debug("grpc request",
			zap.String("method", info.FullMethod),
			zap.String("correlation_id", correlationID),
		)
		ctx = context.WithValue(ctx, correlationIDCtxKey{}, correlationID)
		return handler(ctx, req)
	}
}

// CorrelationClientInterceptor propagates correlation-id from context to outgoing gRPC calls.
func CorrelationClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if correlationID, ok := ctx.Value(correlationIDCtxKey{}).(string); ok && correlationID != "" {
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(correlationIDKey, correlationID))
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// FromContext returns the correlation ID stored in ctx, or "unknown".
func FromContext(ctx context.Context) string {
	if v, ok := ctx.Value(correlationIDCtxKey{}).(string); ok {
		return v
	}
	return "unknown"
}
