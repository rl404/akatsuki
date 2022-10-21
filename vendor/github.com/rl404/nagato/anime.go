package nagato

import (
	"context"
	"net/http"

	"github.com/rl404/nagato/mal"
)

// GetAnimeDetails to get anime details.
func (c *Client) GetAnimeDetails(id int, fields ...AnimeField) (*Anime, int, error) {
	return c.GetAnimeDetailsWithContext(context.Background(), id, fields...)
}

// GetAnimeDetailsWithContext to get anime details with context.
func (c *Client) GetAnimeDetailsWithContext(ctx context.Context, id int, fields ...AnimeField) (*Anime, int, error) {
	if err := c.validate(&idParam{ID: id}); err != nil {
		return nil, http.StatusBadRequest, err
	}

	anime, code, err := c.mal.GetAnimeDetailsWithContext(ctx, id, c.animeFieldsToStrs(fields...)...)
	if err != nil {
		return nil, code, err
	}

	return c.animeToAnime(anime), http.StatusOK, nil
}

// GetAnimeList to get anime list.
func (c *Client) GetAnimeList(param GetAnimeListParam, fields ...AnimeField) ([]Anime, int, error) {
	return c.GetAnimeListWithContext(context.Background(), param, fields...)
}

// GetAnimeListWithContext to get anime list with context.
func (c *Client) GetAnimeListWithContext(ctx context.Context, param GetAnimeListParam, fields ...AnimeField) ([]Anime, int, error) {
	if err := c.validate(&param); err != nil {
		return nil, http.StatusBadRequest, err
	}

	anime, code, err := c.mal.GetAnimeListWithContext(ctx, mal.GetAnimeListParam{
		Query:  param.Query,
		Nsfw:   param.NSFW,
		Limit:  param.Limit,
		Offset: param.Offset,
	}, c.animeFieldsToStrs(fields...)...)
	if err != nil {
		return nil, code, err
	}

	return c.animePagingToAnimeList(anime), http.StatusOK, nil
}

// GetAnimeRanking to get anime ranking list.
func (c *Client) GetAnimeRanking(param GetAnimeRankingParam, fields ...AnimeField) ([]Anime, int, error) {
	return c.GetAnimeRankingWithContext(context.Background(), param, fields...)
}

// GetAnimeRankingWithContext to get anime ranking list with context.
func (c *Client) GetAnimeRankingWithContext(ctx context.Context, param GetAnimeRankingParam, fields ...AnimeField) ([]Anime, int, error) {
	if err := c.validate(&param); err != nil {
		return nil, http.StatusBadRequest, err
	}

	anime, code, err := c.mal.GetAnimeRankingWithContext(ctx, mal.GetAnimeRankingParam{
		RankingType: string(param.RankingType),
		Nsfw:        param.NSFW,
		Limit:       param.Limit,
		Offset:      param.Offset,
	}, c.animeFieldsToStrs(fields...)...)
	if err != nil {
		return nil, code, err
	}

	return c.animeRankingPagingToAnimeList(anime), http.StatusOK, nil
}

// GetSeasonalAnime to get seasonal anime list.
func (c *Client) GetSeasonalAnime(param GetSeasonalAnimeParam, fields ...AnimeField) ([]Anime, int, error) {
	return c.GetSeasonalAnimeWithContext(context.Background(), param, fields...)
}

// GetSeasonalAnimeWithContext to get seasonal anime list with context.
func (c *Client) GetSeasonalAnimeWithContext(ctx context.Context, param GetSeasonalAnimeParam, fields ...AnimeField) ([]Anime, int, error) {
	if err := c.validate(&param); err != nil {
		return nil, http.StatusBadRequest, err
	}

	anime, code, err := c.mal.GetSeasonalAnimeWithContext(ctx, mal.GetSeasonalAnimeParam{
		Year:   param.Year,
		Season: string(param.Season),
		Nsfw:   param.NSFW,
		Sort:   string(param.Sort),
		Limit:  param.Limit,
		Offset: param.Offset,
	}, c.animeFieldsToStrs(fields...)...)
	if err != nil {
		return nil, code, err
	}

	return c.seasonalAnimePagingToAnimeList(anime), http.StatusOK, nil
}

// GetSuggestedAnime to get suggested anime list.
//
// Need oauth2.
func (c *Client) GetSuggestedAnime(param GetSuggestedAnimeParam, fields ...AnimeField) ([]Anime, int, error) {
	return c.GetSuggestedAnimeWithContext(context.Background(), param, fields...)
}

// GetSuggestedAnimeWithContext to get suggested anime list with context.
//
// Need oauth2.
func (c *Client) GetSuggestedAnimeWithContext(ctx context.Context, param GetSuggestedAnimeParam, fields ...AnimeField) ([]Anime, int, error) {
	if err := c.validate(&param); err != nil {
		return nil, http.StatusBadRequest, err
	}

	anime, code, err := c.mal.GetSuggestedAnimeWithContext(ctx, mal.GetSuggestedAnimeParam{
		Nsfw:   param.NSFW,
		Limit:  param.Limit,
		Offset: param.Offset,
	}, c.animeFieldsToStrs(fields...)...)
	if err != nil {
		return nil, code, err
	}

	return c.suggestedAnimePagingToAnimeList(anime), http.StatusOK, nil
}
