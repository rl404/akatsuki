package nagato

import (
	"context"
	"net/http"

	"github.com/rl404/nagato/mal"
)

// GetUserAnimeList to get user anime list.
func (c *Client) GetUserAnimeList(param GetUserAnimeListParam, fields ...AnimeField) ([]UserAnime, int, error) {
	return c.GetUserAnimeListWithContext(context.Background(), param, fields...)
}

// GetUserAnimeListWithContext to get user anime list with context.
func (c *Client) GetUserAnimeListWithContext(ctx context.Context, param GetUserAnimeListParam, fields ...AnimeField) ([]UserAnime, int, error) {
	if err := c.validate(&param); err != nil {
		return nil, http.StatusBadRequest, err
	}

	anime, code, err := c.mal.GetUserAnimeListWithContext(ctx, mal.GetUserAnimeListParam{
		Username: param.Username,
		Status:   string(param.Status),
		Nsfw:     param.NSFW,
		Sort:     string(param.Sort),
		Limit:    param.Limit,
		Offset:   param.Offset,
	}, c.animeFieldsToStrs(fields...)...)
	if err != nil {
		return nil, code, err
	}

	return c.userAnimePagingToUserAnimeList(anime), http.StatusOK, nil
}
