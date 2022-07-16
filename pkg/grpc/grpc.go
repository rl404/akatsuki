package grpc

import (
	"net"
	"time"

	"google.golang.org/grpc"
)

// Server serves GRPC endpoints.
type Server interface {
	Run() chan error
	Server() *grpc.Server
	Close() error
}

type server struct {
	server *grpc.Server
	cfg    Config
}

// Config is basic GRPC server config.
type Config struct {
	Port              string
	Timeout           time.Duration
	UnaryInterceptors []grpc.UnaryServerInterceptor
	StreamInterceptos []grpc.StreamServerInterceptor
	opts              []grpc.ServerOption
}

// New to create new grpc server.
func New(cfg Config) Server {
	return &server{
		server: grpc.NewServer(append(cfg.opts,
			grpc.ConnectionTimeout(cfg.Timeout),
			grpc.UnaryInterceptor(chainUnaryServer(cfg.UnaryInterceptors...)),
			grpc.StreamInterceptor(chainStreamServer(cfg.StreamInterceptos...)))...),
		cfg: cfg,
	}
}

// Run to start serving HTTP.
func (s *server) Run() chan error {
	var ch = make(chan error)
	go s.run(ch)
	return ch
}

func (s *server) run(ch chan error) {
	listener, err := net.Listen("tcp", ":"+s.cfg.Port)
	if err != nil {
		ch <- err
		return
	}
	ch <- s.server.Serve(listener)
}

// Server returns server router.
func (s *server) Server() *grpc.Server {
	return s.server
}

// Close to stop the server gracefully.
func (s *server) Close() error {
	if s.server == nil {
		return nil
	}
	s.server.GracefulStop()
	return nil
}
