package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
)

// @summary Get mal anime by id.
// @tags MAL
// @produce json
// @param animeID path integer true "anime id"
// @success 200 {object} utils.Response{data=mal.Anime}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /mal/anime/{animeID} [get]
func (api *API) handleGetMalAnimeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "animeID"))
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidAnimeID, err))
		return
	}

	anime, code, err := api.service.GetMalAnimeByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, anime, errors.Wrap(r.Context(), err))
}

// @summary Get mal user anime.
// @tags MAL
// @produce json
// @param username path string true "username"
// @param status query string false "status" enum(watching,completed,on_hold,dropped,plan_to_watch)
// @param sort query string false "sort" enum(list_score,list_updated_at,anime_title,anime_start_date)
// @param limit query integer false "limit" default(100)
// @param offset query integer false "offset" default(0)
// @success 200 {object} utils.Response{data=[]mal.UserAnime}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /mal/users/{username}/animelist [get]
func (api *API) handleGetMalUserAnime(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	status := r.URL.Query().Get("status")
	sort := r.URL.Query().Get("sort")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	anime, code, err := api.service.GetMalUserAnime(r.Context(), service.GetMalUserAnimeRequest{
		UserName: username,
		Status:   status,
		Sort:     sort,
		Limit:    limit,
		Offset:   offset,
	})

	utils.ResponseWithJSON(w, code, anime, errors.Wrap(r.Context(), err))
}
