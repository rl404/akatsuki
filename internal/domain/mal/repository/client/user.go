package client

import (
	"context"
	_errors "errors"
	"net/http"
	"strings"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rl404/akatsuki/internal/domain/mal/entity"
	"github.com/rl404/akatsuki/internal/errors"
)

// GetUserAnime to get user anime.
func (c *Client) GetUserAnime(ctx context.Context, data entity.GetUserAnimeRequest) ([]mal.UserAnime, int, error) {
	c.limiter.Take()

	anime, resp, err := c.client.User.AnimeList(ctx, data.Username,
		data.Status,
		data.Sort,
		mal.Limit(data.Limit),
		mal.Offset(data.Offset),
		mal.Fields{"list_status{num_times_rewatched,priority,rewatch_value,tags,comments}"},
		mal.NSFW(true),
	)
	if err != nil {
		if resp != nil {
			return nil, resp.StatusCode, errors.Wrap(ctx, _errors.New(strings.ToLower(http.StatusText(resp.StatusCode))), err)
		}
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	return anime, http.StatusOK, nil
}
