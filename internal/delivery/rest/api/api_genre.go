package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/akatsuki/internal/domain/genre/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
)

// @summary Get genre list.
// @tags Genre
// @produce json
// @param name query string false "name"
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.Genre,meta=service.Pagination}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /genres [get]
func (api *API) handleGetGenres(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	genres, pagination, code, err := api.service.GetGenres(r.Context(), service.GetGenresRequest{
		Name:  name,
		Page:  page,
		Limit: limit,
	})

	utils.ResponseWithJSON(w, code, genres, errors.Wrap(r.Context(), err), pagination)
}

// @summary Get genre by id.
// @tags Genre
// @produce json
// @param genreID path integer true "genre id"
// @success 200 {object} utils.Response{data=service.Genre}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /genres/{genreID} [get]
func (api *API) handleGetGenreByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "genreID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidGenreID, err))
		return
	}

	genre, code, err := api.service.GetGenreByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, genre, errors.Wrap(r.Context(), err))
}

// @summary Get genre stats histories by id.
// @tags Genre
// @produce json
// @param genreID path integer true "genre id"
// @param start_year query integer false "start year"
// @param end_date query integer false "end year"
// @param group query string false "group" enums(MONTHLY,YEARLY) default(MONTHLY)
// @success 200 {object} utils.Response{data=[]service.GenreHistory}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /genres/{genreID}/history [get]
func (api *API) handleGetGenreHistoriesByID(w http.ResponseWriter, r *http.Request) {
	startYear, _ := strconv.Atoi(r.URL.Query().Get("start_year"))
	endYear, _ := strconv.Atoi(r.URL.Query().Get("end_year"))
	group := r.URL.Query().Get("group")

	id, err := strconv.ParseInt(chi.URLParam(r, "genreID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidGenreID, err))
		return
	}

	histories, code, err := api.service.GetGenreHistoriesByID(r.Context(), service.GetGenreHistoriesRequest{
		ID:        id,
		StartYear: startYear,
		EndYear:   endYear,
		Group:     entity.HistoryGroup(group),
	})

	utils.ResponseWithJSON(w, code, histories, errors.Wrap(r.Context(), err))
}
