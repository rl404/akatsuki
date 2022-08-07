package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	animeEntity "github.com/rl404/akatsuki/internal/domain/anime/entity"
	publisherEntity "github.com/rl404/akatsuki/internal/domain/publisher/entity"
	"github.com/rl404/akatsuki/internal/domain/user_anime/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
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
		return nil, nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	userAnime, cnt, code, err := s.userAnime.Get(ctx, entity.GetUserAnimeRequest{
		Username: data.Username,
		Page:     data.Page,
		Limit:    data.Limit,
	})
	if err != nil {
		return nil, nil, code, errors.Wrap(ctx, err)
	}

	if cnt == 0 {
		// Queue to parse.
		if err := s.publisher.PublishParseUserAnime(ctx, publisherEntity.ParseUserAnimeRequest{Username: data.Username}); err != nil {
			return nil, nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
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
	AnimeID         int64              `json:"anime_id"`
	Title           string             `json:"title"`
	Status          animeEntity.Status `json:"status" swaggertype:"string"`
	Score           float64            `json:"score"`
	Type            animeEntity.Type   `json:"type" swaggertype:"string"`
	UserAnimeStatus entity.Status      `json:"user_anime_status" swaggertype:"string"`
	UserAnimeScore  int                `json:"user_anime_score"`
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
		return nil, code, errors.Wrap(ctx, err)
	}

	if len(userAnime) == 0 {
		// Queue to parse.
		if err := s.publisher.PublishParseUserAnime(ctx, publisherEntity.ParseUserAnimeRequest{Username: username}); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
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
		return nil, code, errors.Wrap(ctx, err)
	}

	anime, code, err := s.anime.GetByIDs(ctx, animeIDs)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	for _, a := range anime {
		var status entity.Status
		var score int
		if userAnimeMap[a.ID] != nil {
			status = userAnimeMap[a.ID].Status
			score = userAnimeMap[a.ID].Score
		}
		nodes = append(nodes, userAnimeRelationNode{
			AnimeID:         a.ID,
			Title:           a.Title,
			Status:          a.Status,
			Score:           a.Mean,
			Type:            a.Type,
			UserAnimeStatus: status,
			UserAnimeScore:  score,
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
		return code, errors.Wrap(ctx, err)
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
