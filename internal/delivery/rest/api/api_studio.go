package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/akatsuki/internal/domain/studio/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/errors/stack"
)

// @summary Get studio list.
// @tags Studio
// @produce json
// @param name query string false "name"
// @param sort query string false "sort" enums(NAME,-NAME,COUNT,-COUNT,MEAN,-MEAN,MEMBER,-MEMBER) default(NAME)
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.Studio,meta=service.Pagination}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /studios [get]
func (api *API) handleGetStudios(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	studios, pagination, code, err := api.service.GetStudios(r.Context(), service.GetStudiosRequest{
		Name:  name,
		Sort:  entity.Sort(sort),
		Page:  page,
		Limit: limit,
	})

	utils.ResponseWithJSON(w, code, studios, stack.Wrap(r.Context(), err), pagination)
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
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidStudioID))
		return
	}

	studio, code, err := api.service.GetStudioByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, studio, stack.Wrap(r.Context(), err))
}

// @summary Get studio stats histories by id.
// @tags Studio
// @produce json
// @param studioID path integer true "studio id"
// @param start_year query integer false "start year"
// @param end_date query integer false "end year"
// @param group query string false "group" enums(MONTHLY,YEARLY) default(MONTHLY)
// @success 200 {object} utils.Response{data=[]service.StudioHistory}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /studios/{studioID}/history [get]
func (api *API) handleGetStudioHistoriesByID(w http.ResponseWriter, r *http.Request) {
	startYear, _ := strconv.Atoi(r.URL.Query().Get("start_year"))
	endYear, _ := strconv.Atoi(r.URL.Query().Get("end_year"))
	group := r.URL.Query().Get("group")

	id, err := strconv.ParseInt(chi.URLParam(r, "studioID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidStudioID))
		return
	}

	histories, code, err := api.service.GetStudioHistoriesByID(r.Context(), service.GetStudioHistoriesRequest{
		ID:        id,
		StartYear: startYear,
		EndYear:   endYear,
		Group:     entity.HistoryGroup(group),
	})

	utils.ResponseWithJSON(w, code, histories, stack.Wrap(r.Context(), err))
}
