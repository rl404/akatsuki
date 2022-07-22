package service

import (
	"context"
	"net/http"
	"time"

	publisherEntity "github.com/rl404/akatsuki/internal/domain/publisher/entity"
	"github.com/rl404/akatsuki/internal/domain/user_anime/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
)

// UserAnime is user anime model.
type UserAnime struct {
	AnimeID   int64         `json:"anime_id"`
	Status    entity.Status `json:"status"`
	Score     int           `json:"score"`
	Episode   int           `json:"episode"`
	Tags      []string      `json:"tags"`
	Comment   string        `json:"comment"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// GetUserAnimeRequest is get user anime request model.
type GetUserAnimeRequest struct {
	Username string `validate:"required" mod:"trim"`
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
