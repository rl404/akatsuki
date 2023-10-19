package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	animeEntity "github.com/rl404/akatsuki/internal/domain/anime/entity"
	"github.com/rl404/akatsuki/internal/domain/user_anime/entity"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/errors/stack"
)

// UserAnime is user anime model.
type UserAnime struct {
	AnimeID   int64         `json:"anime_id"`
	Status    entity.Status `json:"status" swaggertype:"string"`
	Score     int           `json:"score"`
	Episode   int           `json:"episode"`
	Tags      []string      `json:"tags"`
	Comment   string        `json:"comment"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// GetUserAnimeRequest is get user anime request model.
type GetUserAnimeRequest struct {
	Username string `validate:"required" mod:"trim,lcase"`
	Page     int    `validate:"required,gte=1" mod:"default=1"`
	Limit    int    `validate:"required,gte=-1" mod:"default=20"`
}

// GetUserAnime to get user anime.
func (s *service) GetUserAnime(ctx context.Context, data GetUserAnimeRequest) ([]UserAnime, *Pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	userAnime, cnt, code, err := s.userAnime.Get(ctx, entity.GetUserAnimeRequest{
		Username: data.Username,
		Page:     data.Page,
		Limit:    data.Limit,
	})
	if err != nil {
		return nil, nil, code, stack.Wrap(ctx, err)
	}

	if cnt == 0 {
		// Queue to parse.
		if err := s.publisher.PublishParseUserAnime(ctx, data.Username, "", false); err != nil {
			return nil, nil, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
		return nil, nil, http.StatusAccepted, nil
	}

	res := make([]UserAnime, len(userAnime))
	for i, ua := range userAnime {
		res[i] = UserAnime{
			AnimeID:   ua.AnimeID,
			Status:    ua.Status,
			Score:     ua.Score,
			Episode:   ua.Episode,
			Tags:      ua.Tags,
			Comment:   ua.Comment,
			UpdatedAt: ua.UpdatedAt,
		}
	}

	return res, &Pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: cnt,
	}, http.StatusOK, nil
}

// UserAnimeRelation is user anime relation model.
type UserAnimeRelation struct {
	Nodes []userAnimeRelationNode `json:"nodes"`
	Links []userAnimeRelationLink `json:"links"`
}

type userAnimeRelationNode struct {
	AnimeID          int64              `json:"anime_id"`
	Title            string             `json:"title"`
	Status           animeEntity.Status `json:"status" swaggertype:"string"`
	Score            float64            `json:"score"`
	Type             animeEntity.Type   `json:"type" swaggertype:"string"`
	Source           animeEntity.Source `json:"source" swaggertype:"string"`
	EpisodeCount     int                `json:"episode_count"`
	EpisodeDuration  int                `json:"episode_duration"`
	StartYear        int                `json:"start_year"`
	Season           animeEntity.Season `json:"season" swaggertype:"string"`
	SeasonYear       int                `json:"season_year"`
	UserAnimeStatus  entity.Status      `json:"user_anime_status" swaggertype:"string"`
	UserAnimeScore   int                `json:"user_anime_score"`
	UserEpisodeCount int                `json:"user_episode_count"`
}

type userAnimeRelationLink struct {
	AnimeID1 int64                `json:"anime_id1"`
	AnimeID2 int64                `json:"anime_id2"`
	Relation animeEntity.Relation `json:"relation" swaggertype:"string"`
}

// GetUserAnimeRelations to get user anime relation.
func (s *service) GetUserAnimeRelations(ctx context.Context, username string) (*UserAnimeRelation, int, error) {
	username = strings.ToLower(username)

	userAnime, _, code, err := s.userAnime.Get(ctx, entity.GetUserAnimeRequest{
		Username: username,
		Page:     1,
		Limit:    -1,
	})
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	if len(userAnime) == 0 {
		// Queue to parse.
		if err := s.publisher.PublishParseUserAnime(ctx, username, "", false); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
		return nil, http.StatusAccepted, nil
	}

	var animeIDs []int64
	nodes := []userAnimeRelationNode{}
	links := []userAnimeRelationLink{}
	nodeMap := make(map[int64]bool)
	userAnimeMap := make(map[int64]*entity.UserAnime)

	for _, ua := range userAnime {
		animeIDs = append(animeIDs, ua.AnimeID)
		userAnimeMap[ua.AnimeID] = ua
	}

	if code, err := s.getUserAnimeRelations(ctx, &animeIDs, nodeMap, &links); err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	anime, code, err := s.anime.GetByIDs(ctx, animeIDs)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	for _, a := range anime {
		var status entity.Status
		var score int
		var userEpisode int
		if userAnimeMap[a.ID] != nil {
			status = userAnimeMap[a.ID].Status
			score = userAnimeMap[a.ID].Score
			userEpisode = userAnimeMap[a.ID].Episode
		}
		nodes = append(nodes, userAnimeRelationNode{
			AnimeID:          a.ID,
			Title:            a.Title,
			Status:           a.Status,
			Score:            a.Mean,
			Type:             a.Type,
			Source:           a.Source,
			StartYear:        a.StartDate.Year,
			EpisodeCount:     a.Episode.Count,
			EpisodeDuration:  a.Episode.Duration,
			Season:           a.Season.Season,
			SeasonYear:       a.Season.Year,
			UserAnimeStatus:  status,
			UserAnimeScore:   score,
			UserEpisodeCount: userEpisode,
		})
	}

	return &UserAnimeRelation{
		Nodes: nodes,
		Links: links,
	}, http.StatusOK, nil
}

func (s *service) getUserAnimeRelations(ctx context.Context, animeIDs *[]int64, nodeMap map[int64]bool, links *[]userAnimeRelationLink) (int, error) {
	var ids []int64

	for _, id := range *animeIDs {
		if !nodeMap[id] {
			ids = append(ids, id)
			nodeMap[id] = true
		}
	}

	if len(ids) == 0 {
		return http.StatusOK, nil
	}

	related, code, err := s.anime.GetRelatedByIDs(ctx, ids)
	if err != nil {
		return code, stack.Wrap(ctx, err)
	}

	for _, r := range related {
		if !nodeMap[r.AnimeID2] {
			*animeIDs = append(*animeIDs, r.AnimeID2)
		}

		*links = append(*links, userAnimeRelationLink{
			AnimeID1: r.AnimeID1,
			AnimeID2: r.AnimeID2,
			Relation: r.Relation,
		})
	}

	return s.getUserAnimeRelations(ctx, animeIDs, nodeMap, links)
}
