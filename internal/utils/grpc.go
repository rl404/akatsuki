package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ResponseWithGRPC to convert http response code and error message
// to grpc error format.
func ResponseWithGRPC(code int, err error) error {
	if err == nil {
		return err
	}
	switch code {
	case http.StatusBadRequest:
		return status.Error(codes.InvalidArgument, err.Error())
	case http.StatusNotFound:
		return status.Error(codes.NotFound, err.Error())
	case http.StatusUnauthorized:
		return status.Error(codes.Unauthenticated, err.Error())
	case http.StatusInternalServerError:
		return status.Error(codes.Internal, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

// RecovererGRPC is custom recoverer middleware.
func RecovererGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Panic: %+v\n", r)
			debug.PrintStack()

			err = status.Error(codes.Internal, "panic")
		}
	}()
	return handler(ctx, req)
}
