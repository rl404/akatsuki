package client

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/nagato"
)

// GetAnimeByID to get anime by id.
func (c *Client) GetAnimeByID(ctx context.Context, id int) (*nagato.Anime, int, error) {
	anime, code, err := c.client.GetAnimeDetailsWithContext(ctx, id,
		nagato.AnimeFieldAlternativeTitles,
		nagato.AnimeFieldStartDate,
		nagato.AnimeFieldEndDate,
		nagato.AnimeFieldSynopsis,
		nagato.AnimeFieldMean,
		nagato.AnimeFieldRank,
		nagato.AnimeFieldPopularity,
		nagato.AnimeFieldNumListUsers,
		nagato.AnimeFieldNumScoringUsers,
		nagato.AnimeFieldNSFW,
		nagato.AnimeFieldGenres,
		nagato.AnimeFieldMediaType,
		nagato.AnimeFieldStatus,
		nagato.AnimeFieldNumEpisodes,
		nagato.AnimeFieldStartSeason,
		nagato.AnimeFieldBroadcast,
		nagato.AnimeFieldSource,
		nagato.AnimeFieldAverageEpisodeDuration,
		nagato.AnimeFieldRating,
		nagato.AnimeFieldStudios,
		nagato.AnimeFieldPictures,
		nagato.AnimeFieldBackground,
		nagato.AnimeFieldStatistics,
		nagato.AnimeFieldNumFavorites,
		nagato.AnimeFieldOpeningThemes,
		nagato.AnimeFieldEndingThemes,
		nagato.AnimeFieldVideos,
		nagato.AnimeFieldRelatedAnime(),
	)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	return anime, http.StatusOK, nil
}
