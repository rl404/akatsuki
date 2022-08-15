package middleware

import (
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/newrelic/go-agent/v3/newrelic"
	"google.golang.org/grpc"
)

// NewUnaryGRPC to create new grpc unary newrelic middleware.
// Just use the original one.
func NewUnaryGRPC(nrApp *newrelic.Application) grpc.UnaryServerInterceptor {
	return nrgrpc.UnaryServerInterceptor(nrApp)
}

// NewStreamGRPC to create new grpc stream newrelic middleware.
// Just use the original one.
func NewStreamGRPC(nrApp *newrelic.Application) grpc.StreamServerInterceptor {
	return nrgrpc.StreamServerInterceptor(nrApp)
}
