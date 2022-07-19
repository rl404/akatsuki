package api

import (
	"github.com/rl404/akatsuki/internal/delivery/grpc/schema"
	"github.com/rl404/akatsuki/internal/service"
)

// API contains all functions for api endpoints.
type API struct {
	service service.Service
	schema.UnimplementedAPIServer
}

// New to create new api endpoints.
func New(service service.Service) *API {
	return &API{
		service: service,
	}
}
