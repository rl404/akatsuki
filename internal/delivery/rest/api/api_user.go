package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
)

// @summary Get user's anime.
// @tags User
// @produce json
// @param username path string true "username"
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.UserAnime,meta=service.Pagination}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /user/{username}/anime [get]
func (api *API) handleGetUserAnime(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	anime, pagination, code, err := api.service.GetUserAnime(r.Context(), service.GetUserAnimeRequest{
		Username: username,
		Page:     page,
		Limit:    limit,
	})

	utils.ResponseWithJSON(w, code, anime, errors.Wrap(r.Context(), err), pagination)
}
