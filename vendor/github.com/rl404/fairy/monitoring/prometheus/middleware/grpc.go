package middleware

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const grpcReqName = "grpc_requests_total"
const grpcLatencyName = "grpc_request_duration_seconds"
const grpcTypeUnary = "unary"
const grpcTypeStream = "stream"

var gm grpcMiddleware

type grpcMiddleware struct {
	req *prometheus.CounterVec
	lat *prometheus.HistogramVec
}

func init() {
	gm = grpcMiddleware{
		req: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: grpcReqName,
				Help: "How many GRPC requests processed, partitioned by status code, type, service and method.",
			},
			[]string{"code", "type", "service", "method"},
		),
		lat: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: grpcLatencyName,
				Help: "How long it took to process the request, partitioned by status code, type, service and method.",
			},
			[]string{"code", "type", "service", "method"},
		),
	}

	prometheus.MustRegister(gm.req)
	prometheus.MustRegister(gm.lat)
}

// NewUnaryGRPC to create new grpc unary prometheus middleware.
func NewUnaryGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	service, method := gm.splitMethodName(info.FullMethod)
	gm.req.WithLabelValues(strconv.Itoa(gm.getCodeFromErr(err)), grpcTypeUnary, service, method).Inc()
	gm.lat.WithLabelValues(strconv.Itoa(gm.getCodeFromErr(err)), grpcTypeUnary, service, method).Observe(float64(time.Since(start).Seconds()))
	return resp, err
}

// NewStreamGRPC to create new grpc stream prometheus middleware.
// Todo: implement prometheus.
func NewStreamGRPC(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return handler(srv, ss)
}

func (m *grpcMiddleware) getCodeFromErr(err error) int {
	if err == nil {
		return int(codes.OK)
	}
	c, ok := status.FromError(err)
	if !ok {
		return int(codes.Unknown)
	}
	return int(c.Code())
}

func (m *grpcMiddleware) splitMethodName(fullMethodName string) (string, string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/")
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		return fullMethodName[:i], fullMethodName[i+1:]
	}
	return "unknown", "unknown"
}
