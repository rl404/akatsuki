package nagato

import (
	"context"
	"net/http"

	"github.com/rl404/nagato/mal"
)

// GetMangaDetails to get manga details.
func (c *Client) GetMangaDetails(id int, fields ...MangaField) (*Manga, int, error) {
	return c.GetMangaDetailsWithContext(context.Background(), id, fields...)
}

// GetMangaDetailsWithContext to get manga details with context.
func (c *Client) GetMangaDetailsWithContext(ctx context.Context, id int, fields ...MangaField) (*Manga, int, error) {
	if err := c.validate(&idParam{ID: id}); err != nil {
		return nil, http.StatusBadRequest, err
	}

	manga, code, err := c.mal.GetMangaDetailsWithContext(ctx, id, c.mangaFieldsToStrs(fields...)...)
	if err != nil {
		return nil, code, err
	}

	return c.mangaToManga(manga), http.StatusOK, nil
}

// GetMangaList to get manga list.
func (c *Client) GetMangaList(param GetMangaListParam, fields ...MangaField) ([]Manga, int, error) {
	return c.GetMangaListWithContext(context.Background(), param, fields...)
}

// GetMangaListWithContext to get manga list with context.
func (c *Client) GetMangaListWithContext(ctx context.Context, param GetMangaListParam, fields ...MangaField) ([]Manga, int, error) {
	if err := c.validate(&param); err != nil {
		return nil, http.StatusBadRequest, err
	}

	manga, code, err := c.mal.GetMangaListWithContext(ctx, mal.GetMangaListParam{
		Query:  param.Query,
		Nsfw:   param.NSFW,
		Limit:  param.Limit,
		Offset: param.Offset,
	}, c.mangaFieldsToStrs(fields...)...)
	if err != nil {
		return nil, code, err
	}

	return c.mangaPagingToMangaList(manga), http.StatusOK, nil
}

// GetMangaRanking to get manga ranking.
func (c *Client) GetMangaRanking(param GetMangaRankingParam, fields ...MangaField) ([]Manga, int, error) {
	return c.GetMangaRankingWithContext(context.Background(), param, fields...)
}

// GetMangaRankingWithContext to get manga ranking with context.
func (c *Client) GetMangaRankingWithContext(ctx context.Context, param GetMangaRankingParam, fields ...MangaField) ([]Manga, int, error) {
	if err := c.validate(&param); err != nil {
		return nil, http.StatusBadRequest, err
	}

	manga, code, err := c.mal.GetMangaRankingWithContext(ctx, mal.GetMangaRankingParam{
		RankingType: string(param.RankingType),
		Nsfw:        param.NSFW,
		Limit:       param.Limit,
		Offset:      param.Offset,
	}, c.mangaFieldsToStrs(fields...)...)
	if err != nil {
		return nil, code, err
	}

	return c.mangaRankingPagingToMangaList(manga), http.StatusOK, nil
}
