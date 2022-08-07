package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
)

// @summary Get anime by id.
// @tags Anime
// @produce json
// @param animeID path integer true "anime id"
// @success 200 {object} utils.Response{data=service.Anime}
// @failure 202 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /anime/{animeID} [get]
func (api *API) handleGetAnimeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "animeID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidAnimeID, err))
		return
	}

	anime, code, err := api.service.GetAnimeByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, anime, errors.Wrap(r.Context(), err))
}

// @summary Update anime by id.
// @tags Anime
// @produce json
// @param animeID path integer true "anime id"
// @success 202 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /anime/{animeID}/update [post]
func (api *API) handleUpdateAnimeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "animeID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidAnimeID, err))
		return
	}

	code, err := api.service.UpdateAnimeByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, nil, errors.Wrap(r.Context(), err))
}
