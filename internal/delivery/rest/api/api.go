package api

import (
	"github.com/go-chi/chi"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/monitoring/prometheus/middleware"
)

// API contains all functions for api endpoints.
type API struct {
	service service.Service
}

// New to create new api endpoints.
func New(service service.Service) *API {
	return &API{
		service: service,
	}
}

// Register to register api routes.
func (api *API) Register(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.NewHTTP())
		r.Use(log.MiddlewareWithLog(utils.GetLogger(0), log.MiddlewareConfig{Error: true}))
		r.Use(utils.Recoverer)

		r.Get("/anime/{animeID}", api.handleGetAnimeByID)
		r.Post("/anime/{animeID}/update", api.handleUpdateAnimeByID)

		r.Get("/user/{username}/anime", api.handleGetUserAnime)
		r.Get("/user/{username}/anime/relations", api.handleGetUserAnimeRelations)
		r.Post("/user/{username}/update", api.handleUpdateUserAnime)

		r.Get("/mal/anime/{animeID}", api.handleGetMalAnimeByID)
		r.Get("/mal/users/{username}/animelist", api.handleGetMalUserAnime)
	})
}
