package nagato

import (
	"context"
	"net/http"
	"strings"

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

// UpdateMyAnimeListStatus to update my anime list status.
//
// Need oauth2.
func (c *Client) UpdateMyAnimeListStatus(param UpdateMyAnimeListStatusParam) (*UserAnimeListStatus, int, error) {
	return c.UpdateMyAnimeListStatusWithContext(context.Background(), param)
}

// UpdateMyAnimeListStatusWithContext to update my anime list status with context.
//
// Need oauth2.
func (c *Client) UpdateMyAnimeListStatusWithContext(ctx context.Context, param UpdateMyAnimeListStatusParam) (*UserAnimeListStatus, int, error) {
	if err := c.validate(&param); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := c.validateDate(param.StartDate); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := c.validateDate(param.FinishDate); err != nil {
		return nil, http.StatusBadRequest, err
	}

	status, code, err := c.mal.UpdateMyAnimeListStatusWithContext(ctx, mal.UpdateMyAnimeListStatusParam{
		ID:                 param.ID,
		Status:             string(param.Status),
		IsRewatching:       param.IsRewatching,
		Score:              param.Score,
		NumWatchedEpisodes: param.Episode,
		Priority:           int(param.Priority),
		NumTimesRewatched:  param.RewatchedTimes,
		RewatchValue:       int(param.RewatchValue),
		Tags:               strings.Join(param.Tags, ","),
		Comments:           param.Comment,
		StartDate:          c.dateToStr(param.StartDate),
		FinishDate:         c.dateToStr(param.FinishDate),
	})
	if err != nil {
		return nil, code, err
	}

	return c.listStatusToUserAnimeListStatus(*status), http.StatusOK, nil
}

// DeleteMyAnimeListStatus to delete my anime list status.
//
// Need oauth2.
func (c *Client) DeleteMyAnimeListStatus(id int) (int, error) {
	return c.DeleteMyAnimeListStatusWithContext(context.Background(), id)
}

// DeleteMyAnimeListStatusWithContext to delete my anime list status with context.
//
// Need oauth2.
func (c *Client) DeleteMyAnimeListStatusWithContext(ctx context.Context, id int) (int, error) {
	if err := c.validate(&idParam{ID: id}); err != nil {
		return http.StatusBadRequest, err
	}

	return c.mal.DeleteMyAnimeListStatusWithContext(ctx, id)
}
