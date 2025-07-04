package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/akatsuki/internal/domain/anime/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/errors/stack"
)

// @summary Get anime list.
// @tags Anime
// @produce json
// @param title query string false "title"
// @param nsfw query string false "nsfw" enums(true,false)
// @param type query string false "type" enums(TV,OVA,ONA,MOVIE,SPECIAL,MUSIC,CM,PV,TV_SPECIAL)
// @param status query string false "status" enums(FINISHED,RELEASING,NOT_YET)
// @param season query string false "season" enums(WINTER,SPRING,SUMMER,FALL)
// @param season_year query integer false "season year"
// @param start_mean query number false "start mean"
// @param end_mean query number false "end mean"
// @param start_airing_year query number false "start airing year"
// @param end_airing_year query number false "end airing year"
// @param genre_id query integer false "genre id"
// @param studio_id query integer false "studio id"
// @param sort query string false "sort" enums(ID,-ID,TITLE,-TITLE,START_DATE,-START_DATE,MEAN,-MEAN,RANK,-RANK,POPULARITY,-POPULARITY,MEMBER,-MEMBER,VOTER,-VOTER) default(RANK)
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.Anime,meta=service.Pagination}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /anime [get]
func (api *API) HandleGetAnime(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	nsfw := r.URL.Query().Get("nsfw")
	_type := r.URL.Query().Get("type")
	status := r.URL.Query().Get("status")
	season := r.URL.Query().Get("season")
	seasonYear, _ := strconv.Atoi(r.URL.Query().Get("season_year"))
	startAiringYear, _ := strconv.Atoi(r.URL.Query().Get("start_airing_year"))
	endAiringYear, _ := strconv.Atoi(r.URL.Query().Get("end_airing_year"))
	genreID, _ := strconv.ParseInt(r.URL.Query().Get("genre_id"), 10, 64)
	studioID, _ := strconv.ParseInt(r.URL.Query().Get("studio_id"), 10, 64)
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	var startMean, endMean float64
	if tmp := r.URL.Query().Get("start_mean"); tmp != "" {
		tmp2, err := strconv.ParseFloat(tmp, 64)
		if err != nil {
			utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidFormat("start_mean")))
			return
		}
		startMean = tmp2
	}
	if tmp := r.URL.Query().Get("end_mean"); tmp != "" {
		tmp2, err := strconv.ParseFloat(tmp, 64)
		if err != nil {
			utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidFormat("end_mean")))
			return
		}
		endMean = tmp2
	}

	anime, pagination, code, err := api.service.GetAnime(r.Context(), service.GetAnimeRequest{
		Title:           title,
		NSFW:            utils.ParseToBoolPtr(nsfw),
		Type:            entity.Type(_type),
		Status:          entity.Status(status),
		Season:          entity.Season(season),
		SeasonYear:      seasonYear,
		StartMean:       startMean,
		EndMean:         endMean,
		StartAiringYear: startAiringYear,
		EndAiringYear:   endAiringYear,
		GenreID:         genreID,
		StudioID:        studioID,
		Sort:            entity.Sort(sort),
		Page:            page,
		Limit:           limit,
	})

	utils.ResponseWithJSON(w, code, anime, stack.Wrap(r.Context(), err), pagination)
}

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
func (api *API) HandleGetAnimeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "animeID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidAnimeID))
		return
	}

	anime, code, err := api.service.GetAnimeByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, anime, stack.Wrap(r.Context(), err))
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
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidAnimeID))
		return
	}

	code, err := api.service.UpdateAnimeByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
}

// @summary Get anime stats histories by id.
// @tags Anime
// @produce json
// @param animeID path integer true "anime id"
// @param start_date query string false "start date (yyyy-mm-dd)"
// @param end_date query string false "end date (yyyy-mm-dd)"
// @param group query string false "group" enums(WEEKLY,MONTHLY,YEARLY) default(MONTHLY)
// @success 200 {object} utils.Response{data=[]service.AnimeHistory}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /anime/{animeID}/history [get]
func (api *API) handleGetAnimeHistoriesByID(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	group := r.URL.Query().Get("group")

	id, err := strconv.ParseInt(chi.URLParam(r, "animeID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidAnimeID))
		return
	}

	histories, code, err := api.service.GetAnimeHistoriesByID(r.Context(), service.GetAnimeHistoriesRequest{
		ID:        id,
		StartDate: startDate,
		EndDate:   endDate,
		Group:     entity.HistoryGroup(group),
	})

	utils.ResponseWithJSON(w, code, histories, stack.Wrap(r.Context(), err))
}
