package shared_infrastructure

import (
	"context"
	"log/slog"

	bsgostuff_domain "github.com/beavernsticks/go-stuff/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	slog.Error("gRPC error", slog.String("method", info.FullMethod), slog.Any("error", err))

	// Обрабатываем стандартные ошибки
	switch err {
	case context.DeadlineExceeded:
		return nil, status.Error(codes.DeadlineExceeded, "request timed out")
	case context.Canceled:
		return nil, status.Error(codes.Canceled, "request canceled")
	case bsgostuff_domain.ErrForbidden:
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	case bsgostuff_domain.ErrInvalidArgument:
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	case bsgostuff_domain.ErrNotFound:
		return nil, status.Error(codes.NotFound, "not found")
	case bsgostuff_domain.ErrDuplicate:
		return nil, status.Error(codes.AlreadyExists, "already exists")
	default:
		return nil, status.Error(codes.Internal, "internal server error")
	}
}

func ErrorStreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, stream)
	if err == nil {
		return nil
	}

	slog.Error("gRPC error", slog.String("method", info.FullMethod), slog.Any("error", err))

	// Обрабатываем стандартные ошибки
	switch err {
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, "request timed out")
	case context.Canceled:
		return status.Error(codes.Canceled, "request canceled")
	case bsgostuff_domain.ErrForbidden:
		return status.Error(codes.PermissionDenied, "forbidden")
	case bsgostuff_domain.ErrInvalidArgument:
		return status.Error(codes.InvalidArgument, "invalid argument")
	case bsgostuff_domain.ErrNotFound:
		return status.Error(codes.NotFound, "not found")
	case bsgostuff_domain.ErrDuplicate:
		return status.Error(codes.AlreadyExists, "already exists")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

func ErrorClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		slog.Error("gRPC error", slog.String("method", method), slog.Any("error", err))

		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.DeadlineExceeded:
				return context.DeadlineExceeded
			case codes.Canceled:
				return context.Canceled
			case codes.PermissionDenied:
				return bsgostuff_domain.ErrForbidden
			case codes.InvalidArgument:
				return bsgostuff_domain.ErrInvalidArgument
			case codes.NotFound:
				return bsgostuff_domain.ErrNotFound
			case codes.AlreadyExists:
				return bsgostuff_domain.ErrDuplicate
			default:
				return bsgostuff_domain.ErrInternal
			}
		}

		return bsgostuff_domain.ErrInternal
	}

	return nil
}

func ErrorStreamClientInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	// Создаем обертку для ClientStream
	clientStream, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		slog.Error("gRPC error", slog.String("method", method), slog.Any("error", err))

		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.DeadlineExceeded:
				return nil, context.DeadlineExceeded
			case codes.Canceled:
				return nil, context.Canceled
			case codes.PermissionDenied:
				return nil, bsgostuff_domain.ErrForbidden
			case codes.InvalidArgument:
				return nil, bsgostuff_domain.ErrInvalidArgument
			case codes.NotFound:
				return nil, bsgostuff_domain.ErrNotFound
			case codes.AlreadyExists:
				return nil, bsgostuff_domain.ErrDuplicate
			default:
				return nil, bsgostuff_domain.ErrInternal
			}
		}

		return nil, bsgostuff_domain.ErrInternal
	}

	return clientStream, nil
}
