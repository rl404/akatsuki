package client

import (
	"context"
	_errors "errors"
	"net/http"
	"strings"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rl404/akatsuki/internal/errors"
)

// GetAnimeByID to get anime by id.
func (c *Client) GetAnimeByID(ctx context.Context, id int) (*mal.Anime, int, error) {
	c.limiter.Take()

	anime, resp, err := c.client.Anime.Details(ctx, id, mal.Fields{
		"alternative_titles",
		"start_date",
		"end_date",
		"synopsis",
		"mean",
		"rank",
		"popularity",
		"num_list_users",
		"num_scoring_users",
		"nsfw",
		"genres",
		"media_type",
		"status",
		"num_episodes",
		"start_season",
		"broadcast",
		"source",
		"average_episode_duration",
		"rating",
		"studios",
		"pictures",
		"background",
		"related_anime",
		"related_manga",
		"recommendations",
		"statistics",
		"opening_themes", // undocumented
		"ending_themes",  // undocumented
	})
	if err != nil {
		if resp != nil {
			return nil, resp.StatusCode, errors.Wrap(ctx, _errors.New(strings.ToLower(http.StatusText(resp.StatusCode))), err)
		}
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	return anime, http.StatusOK, nil
}
