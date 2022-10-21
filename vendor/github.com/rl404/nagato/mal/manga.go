package mal

import (
	"context"
	"net/http"
)

// GetMangaDetails to get manga details.
func (c *Client) GetMangaDetails(id int, fields ...string) (*Manga, int, error) {
	return c.GetMangaDetailsWithContext(context.Background(), id, fields...)
}

// GetMangaDetailsWithContext to get manga details with context.
func (c *Client) GetMangaDetailsWithContext(ctx context.Context, id int, fields ...string) (*Manga, int, error) {
	url := c.generateURL(map[string]interface{}{
		"fields": fields,
	}, "manga", id)

	var manga Manga
	if code, err := c.get(ctx, url, &manga); err != nil {
		return nil, code, err
	}

	return &manga, http.StatusOK, nil
}

// GetMangaList to get manga list.
func (c *Client) GetMangaList(param GetMangaListParam, fields ...string) (*MangaPaging, int, error) {
	return c.GetMangaListWithContext(context.Background(), param, fields...)
}

//  GetMangaListWithContext to get manga list with context.
func (c *Client) GetMangaListWithContext(ctx context.Context, param GetMangaListParam, fields ...string) (*MangaPaging, int, error) {
	url := c.generateURL(map[string]interface{}{
		"q":      param.Query,
		"nsfw":   param.Nsfw,
		"limit":  param.Limit,
		"offset": param.Offset,
		"fields": fields,
	}, "manga")

	var manga MangaPaging
	if code, err := c.get(ctx, url, &manga); err != nil {
		return nil, code, err
	}

	return &manga, http.StatusOK, nil
}

// GetMangaRanking to get manga ranking list.
func (c *Client) GetMangaRanking(param GetMangaRankingParam, fields ...string) (*MangaRankingPaging, int, error) {
	return c.GetMangaRankingWithContext(context.Background(), param, fields...)
}

// GetMangaRankingWithContext to get manga ranking list with context.
func (c *Client) GetMangaRankingWithContext(ctx context.Context, param GetMangaRankingParam, fields ...string) (*MangaRankingPaging, int, error) {
	url := c.generateURL(map[string]interface{}{
		"ranking_type": param.RankingType,
		"nsfw":         param.Nsfw,
		"limit":        param.Limit,
		"offset":       param.Offset,
		"fields":       fields,
	}, "manga", "ranking")

	var manga MangaRankingPaging
	if code, err := c.get(ctx, url, &manga); err != nil {
		return nil, code, err
	}

	return &manga, http.StatusOK, nil
}
