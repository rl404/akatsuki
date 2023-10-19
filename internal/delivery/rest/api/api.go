package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/monitoring/newrelic/middleware"
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
func (api *API) Register(r chi.Router, nrApp *newrelic.Application) {
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.NewHTTP(nrApp))
		r.Use(log.HTTPMiddlewareWithLog(utils.GetLogger(0), log.APIMiddlewareConfig{Error: true}))
		r.Use(log.HTTPMiddlewareWithLog(utils.GetLogger(1), log.APIMiddlewareConfig{
			RequestHeader:  true,
			RequestBody:    true,
			ResponseHeader: true,
			ResponseBody:   true,
			RawPath:        true,
			Error:          true,
		}))
		r.Use(utils.Recoverer)

		r.Get("/anime", api.handleGetAnime)
		r.Get("/anime/{animeID}", api.handleGetAnimeByID)
		r.Post("/anime/{animeID}/update", api.handleUpdateAnimeByID)
		r.Get("/anime/{animeID}/history", api.handleGetAnimeHistoriesByID)

		r.Get("/genres", api.handleGetGenres)
		r.Get("/genres/{genreID}", api.handleGetGenreByID)
		r.Get("/genres/{genreID}/history", api.handleGetGenreHistoriesByID)

		r.Get("/studios", api.handleGetStudios)
		r.Get("/studios/{studioID}", api.handleGetStudioByID)
		r.Get("/studios/{studioID}/history", api.handleGetStudioHistoriesByID)

		r.Get("/user/{username}/anime", api.handleGetUserAnime)
		r.Get("/user/{username}/anime/relations", api.handleGetUserAnimeRelations)
		r.Post("/user/{username}/update", api.handleUpdateUserAnime)
	})
}
