package nagato

import (
	"context"
	"net/http"
	"strings"

	"github.com/rl404/nagato/mal"
)

// GetUserMangaList to get user manga list.
func (c *Client) GetUserMangaList(param GetUserMangaListParam, fields ...MangaField) ([]UserManga, int, error) {
	return c.GetUserMangaListWithContext(context.Background(), param, fields...)
}

// GetUserMangaListWithContext to get user manga list with context.
func (c *Client) GetUserMangaListWithContext(ctx context.Context, param GetUserMangaListParam, fields ...MangaField) ([]UserManga, int, error) {
	if err := c.validate(&param); err != nil {
		return nil, http.StatusBadRequest, err
	}

	manga, code, err := c.mal.GetUserMangaListWithContext(ctx, mal.GetUserMangaListParam{
		Username: param.Username,
		Status:   string(param.Status),
		Nsfw:     param.NSFW,
		Sort:     string(param.Sort),
		Limit:    param.Limit,
		Offset:   param.Offset,
	}, c.mangaFieldsToStrs(fields...)...)
	if err != nil {
		return nil, code, err
	}

	return c.userMangaPagingToUserMangaList(manga), http.StatusOK, nil
}

// UpdateMyMangaListStatus to update my manga list status.
//
// Need oauth2.
func (c *Client) UpdateMyMangaListStatus(param UpdateMyMangaListStatusParam) (*UserMangaListStatus, int, error) {
	return c.UpdateMyMangaListStatusWithContext(context.Background(), param)
}

// UpdateMyMangaListStatusWithContext to update my manga list status with context.
//
// Need oauth2.
func (c *Client) UpdateMyMangaListStatusWithContext(ctx context.Context, param UpdateMyMangaListStatusParam) (*UserMangaListStatus, int, error) {
	if err := c.validate(&param); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := c.validateDate(param.StartDate); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := c.validateDate(param.FinishDate); err != nil {
		return nil, http.StatusBadRequest, err
	}

	status, code, err := c.mal.UpdateMyMangaListStatusWithContext(ctx, mal.UpdateMyMangaListStatusParam{
		ID:              param.ID,
		Status:          string(param.Status),
		IsRereading:     param.IsRereading,
		Score:           param.Score,
		NumChaptersRead: param.Chapter,
		NumVolumesRead:  param.Volume,
		Priority:        int(param.Priority),
		NumTimesReread:  param.RereadTimes,
		RereadValue:     int(param.RereadValue),
		Tags:            strings.Join(param.Tags, ","),
		Comments:        param.Comment,
		StartDate:       c.dateToStr(param.StartDate),
		FinishDate:      c.dateToStr(param.FinishDate),
	})
	if err != nil {
		return nil, code, err
	}

	return c.listStatusToUserMangaListStatus(*status), http.StatusOK, nil
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
	if err := c.validate(&idParam{ID: id}); err != nil {
		return http.StatusBadRequest, err
	}

	return c.mal.DeleteMyMangaListStatusWithContext(ctx, id)
}
