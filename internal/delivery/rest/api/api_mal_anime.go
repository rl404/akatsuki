package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
)

// @summary Get mal anime by id.
// @tags MAL Anime
// @produce json
// @param animeID path string true "anime id"
// @success 200 {object} utils.Response{data=mal.Anime}
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 403 {object} utils.Response
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
