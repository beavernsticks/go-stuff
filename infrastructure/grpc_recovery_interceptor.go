package bsgostuff_infrastructure

import (
	"context"
	"log/slog"
	"runtime/debug"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryUnaryInterceptor возвращает unary-интерцептор с обработкой паник
func RecoveryUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			logPanic(ctx, r, info.FullMethod)
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()

	return handler(ctx, req)
}

// RecoveryStreamInterceptor возвращает stream-интерцептор с обработкой паник
func RecoveryStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logPanic(ss.Context(), r, info.FullMethod)
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()

	return handler(srv, ss)
}

// logPanic логирует информацию о панике с использованием slog
func logPanic(ctx context.Context, panicValue interface{}, method string) {
	stack := string(debug.Stack())
	stackLines := strings.Split(stack, "\n")

	slog.ErrorContext(ctx, "gRPC panic recovered",
		slog.String("method", method),
		slog.Any("panic", panicValue),
		slog.String("stack", strings.Join(stackLines[:min(10, len(stackLines))], "\n")), // Берем первые 10 строк stack trace
	)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
