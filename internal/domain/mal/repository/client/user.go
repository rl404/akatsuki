package client

import (
	"context"
	"net/http"

	"github.com/rl404/akatsuki/internal/domain/mal/entity"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/nagato"
)

// GetUserAnime to get user anime.
func (c *Client) GetUserAnime(ctx context.Context, data entity.GetUserAnimeRequest) ([]nagato.UserAnime, int, error) {
	anime, code, err := c.client.GetUserAnimeListWithContext(ctx, nagato.GetUserAnimeListParam{
		Username: data.Username,
		Status:   nagato.UserAnimeStatusType(data.Status),
		NSFW:     true,
		Limit:    data.Limit,
		Offset:   data.Offset,
	},
		nagato.AnimeFieldUserStatus(
			nagato.UserAnimeNumTimesRewatched,
			nagato.UserAnimeRewatchValue,
			nagato.UserAnimeTags,
			nagato.UserAnimeComments,
		),
	)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	return anime, http.StatusOK, nil
}
