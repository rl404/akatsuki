package log

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/rl404/fairy/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// UnaryMiddlewareWithLog is unary grpc middleware that will log the request and response.
func UnaryMiddlewareWithLog(logger Logger, middlewareConfig ...MiddlewareConfig) grpc.UnaryServerInterceptor {
	var cfg MiddlewareConfig
	if len(middlewareConfig) > 0 {
		cfg = middlewareConfig[0]
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if logger == nil {
			return handler(ctx, req)
		}

		// Prepare error stack tracing.
		s := errors.New()
		ctx = s.Init(ctx)
		start := time.Now()

		// Call handler.
		resp, err := handler(ctx, req)

		// Prepare map for logging.
		m := map[string]interface{}{
			"level":    getLevelFromError(err),
			"duration": time.Since(start).String(),
			"code":     getCodeFromErr(err),
			"ip":       getIPFromCtx(ctx),
		}

		m["service"], m["method"] = splitMethodName(info.FullMethod)

		if cfg.RequestHeader {
			meta, _ := metadata.FromIncomingContext(ctx)
			header, _ := json.Marshal(meta)
			m["request_header"] = string(header)
		}

		if cfg.RequestBody {
			body, _ := json.Marshal(req)
			m["request_body"] = string(body)
		}

		if cfg.ResponseBody {
			body, _ := json.Marshal(resp)
			m["response_body"] = string(body)
		}

		if cfg.ResponseHeader {
			// todo: how?
		}

		// Include the error stack if you use it.
		errStack := s.Get(ctx).([]string)
		if len(errStack) > 0 {
			// Reverse the stack order.
			for i, j := 0, len(errStack)-1; i < j; i, j = i+1, j-1 {
				errStack[i], errStack[j] = errStack[j], errStack[i]
			}
			m["error"] = errStack
		}

		logger.Log(m)

		return resp, err
	}
}

// StreamMiddlewareWithLog is stream grpc middleware that will log the request and response.
// Todo: implement logger here.
func StreamMiddlewareWithLog(logger Logger, middlewareConfig ...MiddlewareConfig) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, ss)
	}
}

func getCodeFromErr(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	c, ok := status.FromError(err)
	if !ok {
		return codes.Unknown
	}
	return c.Code()
}

func getLevelFromError(err error) LogLevel {
	if err == nil {
		return InfoLevel
	}
	e, _ := status.FromError(err)
	switch e.Code() {
	case
		codes.Canceled,
		codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.Aborted:
		return WarnLevel
	default:
		return ErrorLevel
	}
}

func getIPFromCtx(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	return p.Addr.String()
}

func splitMethodName(fullMethodName string) (string, string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/")
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		return fullMethodName[:i], fullMethodName[i+1:]
	}
	return "unknown", "unknown"
}
