package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
)

// @summary Get studio list.
// @tags Studio
// @produce json
// @param name query string false "name"
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.Studio,meta=service.Pagination}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /studios [get]
func (api *API) handleGetStudios(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	studios, pagination, code, err := api.service.GetStudios(r.Context(), service.GetStudiosRequest{
		Name:  name,
		Page:  page,
		Limit: limit,
	})

	utils.ResponseWithJSON(w, code, studios, errors.Wrap(r.Context(), err), pagination)
}

// @summary Get studio by id.
// @tags Studio
// @produce json
// @param studioID path integer true "studio id"
// @success 200 {object} utils.Response{data=service.Studio}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /studios/{studioID} [get]
func (api *API) handleGetStudioByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "studioID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidStudioID, err))
		return
	}

	studio, code, err := api.service.GetStudioByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, studio, errors.Wrap(r.Context(), err))
}
