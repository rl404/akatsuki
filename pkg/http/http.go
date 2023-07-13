package http

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Server serves HTTP endpoints.
type Server interface {
	Run() chan error
	Router() chi.Router
	Close() error
}

type server struct {
	server *http.Server
	router chi.Router
	cfg    Config
}

// Config is basic HTTP server config.
type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	GracefulTimeout time.Duration
}

// New to create new web server.
func New(cfg Config) Server {
	return &server{
		router: chi.NewRouter(),
		cfg:    cfg,
	}
}

// Router returns server router.
func (s *server) Router() chi.Router {
	return s.router
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

	s.server = &http.Server{
		Handler:      s.router,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	ch <- s.server.Serve(listener)
}

// Close to stop the server gracefully.
func (s *server) Close() error {
	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.GracefulTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
