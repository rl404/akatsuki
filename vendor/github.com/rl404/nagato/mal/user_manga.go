package mal

import (
	"context"
	"net/http"
)

// UpdateMyMangaListStatus to update my manga list status.
//
// Need oauth2.
func (c *Client) UpdateMyMangaListStatus(param UpdateMyMangaListStatusParam) (*MyMangaListStatus, int, error) {
	return c.UpdateMyMangaListStatusWithContext(context.Background(), param)
}

// UpdateMyMangaListStatusWithContext to update my manga list status with context.
//
// Need oauth2.
func (c *Client) UpdateMyMangaListStatusWithContext(ctx context.Context, param UpdateMyMangaListStatusParam) (*MyMangaListStatus, int, error) {
	url := c.generateURL(nil, "manga", param.ID, "my_list_status")

	var status MyMangaListStatus
	if code, err := c.patch(ctx, url, param.encode(), &status); err != nil {
		return nil, code, err
	}

	return &status, http.StatusOK, nil
}

// DeleteMyMangaListStatus to delete my manga list status.
//
// Need oauth2.
func (c *Client) DeleteMyMangaListStatus(id int) (int, error) {
	return c.DeleteMyMangaListStatusWithContext(context.Background(), id)
}

// DeleteMyMangaListStatusWithContext to delete my manga list status with context.
//
// Need oauth2.
func (c *Client) DeleteMyMangaListStatusWithContext(ctx context.Context, id int) (int, error) {
	url := c.generateURL(nil, "manga", id, "my_list_status")

	if code, err := c.delete(ctx, url); err != nil {
		return code, err
	}

	return http.StatusOK, nil
}

// GetUserMangaList to get user manga list.
//
// Need oauth2.
func (c *Client) GetUserMangaList(param GetUserMangaListParam, fields ...string) (*UserMangaPaging, int, error) {
	return c.GetUserMangaListWithContext(context.Background(), param, fields...)
}

// GetUserMangaListWithContext to get user manga list with context.
//
// Need oauth2.
func (c *Client) GetUserMangaListWithContext(ctx context.Context, param GetUserMangaListParam, fields ...string) (*UserMangaPaging, int, error) {
	url := c.generateURL(map[string]interface{}{
		"status": param.Status,
		"sort":   param.Sort,
		"nsfw":   param.Nsfw,
		"limit":  param.Limit,
		"offset": param.Offset,
		"fields": fields,
	}, "users", param.Username, "mangalist")

	var userManga UserMangaPaging
	if code, err := c.get(ctx, url, &userManga); err != nil {
		return nil, code, err
	}

	return &userManga, http.StatusOK, nil
}
