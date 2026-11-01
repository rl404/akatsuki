package api

import (
	"github.com/rl404/akatsuki/api"
	"github.com/rl404/akatsuki/internal/service"
)

// API contains all functions for api endpoints.
type API struct {
	api.UnimplementedAPIServer
	service service.Service
}

// New to create new api endpoints.
func New(service service.Service) *API {
	return &API{
		service: service,
	}
}
