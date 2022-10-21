package mal

import (
	"context"
	"net/http"
)

// UpdateMyAnimeListStatus to update my anime list status.
//
// Need oauth2.
func (c *Client) UpdateMyAnimeListStatus(param UpdateMyAnimeListStatusParam) (*MyAnimeListStatus, int, error) {
	return c.UpdateMyAnimeListStatusWithContext(context.Background(), param)
}

// UpdateMyAnimeListStatusWithContext to update my anime list status with context.
//
// Need oauth2.
func (c *Client) UpdateMyAnimeListStatusWithContext(ctx context.Context, param UpdateMyAnimeListStatusParam) (*MyAnimeListStatus, int, error) {
	url := c.generateURL(nil, "anime", param.ID, "my_list_status")

	var status MyAnimeListStatus
	if code, err := c.patch(ctx, url, param.encode(), &status); err != nil {
		return nil, code, err
	}

	return &status, http.StatusOK, nil
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
	url := c.generateURL(nil, "anime", id, "my_list_status")

	if code, err := c.delete(ctx, url); err != nil {
		return code, err
	}

	return http.StatusOK, nil
}

// GetUserAnimeList to get user anime list.
//
// Need oauth2.
func (c *Client) GetUserAnimeList(param GetUserAnimeListParam, fields ...string) (*UserAnimePaging, int, error) {
	return c.GetUserAnimeListWithContext(context.Background(), param, fields...)
}

// GetUserAnimeListWithContext to get user anime list with context.
//
// Need oauth2.
func (c *Client) GetUserAnimeListWithContext(ctx context.Context, param GetUserAnimeListParam, fields ...string) (*UserAnimePaging, int, error) {
	url := c.generateURL(map[string]interface{}{
		"status": param.Status,
		"sort":   param.Sort,
		"nsfw":   param.Nsfw,
		"limit":  param.Limit,
		"offset": param.Offset,
		"fields": fields,
	}, "users", param.Username, "animelist")

	var userAnime UserAnimePaging
	if code, err := c.get(ctx, url, &userAnime); err != nil {
		return nil, code, err
	}

	return &userAnime, http.StatusOK, nil
}
