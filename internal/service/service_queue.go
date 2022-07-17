package service

import (
	"context"
	"net/http"

	"github.com/rl404/akatsuki/internal/domain/publisher/entity"
	"github.com/rl404/akatsuki/internal/errors"
)

// QueueMissingAnime to queue missing anime.
func (s *service) QueueMissingAnime(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	// Get max id.
	maxID, code, err := s.anime.GetMaxID(ctx)
	if err != nil {
		return cnt, code, errors.Wrap(ctx, err)
	}

	// Get all existing anime id.
	animeIDs, code, err := s.anime.GetIDs(ctx)
	if err != nil {
		return cnt, code, errors.Wrap(ctx, err)
	}

	// Get all empty anime id,
	emptyIDs, code, err := s.emptyID.GetIDs(ctx)
	if err != nil {
		return cnt, code, errors.Wrap(ctx, err)
	}

	idMap := make(map[int64]bool)
	for _, id := range animeIDs {
		idMap[id] = true
	}
	for _, id := range emptyIDs {
		idMap[id] = true
	}

	// Loop until max id.
	for id := int64(1); id <= maxID && cnt < limit; id++ {
		if idMap[id] {
			continue
		}

		if err := s.publisher.PublishParseAnime(ctx, entity.ParseAnimeRequest{ID: id}); err != nil {
			return cnt, http.StatusInternalServerError, errors.Wrap(ctx, err)
		}

		cnt++
	}

	return cnt, http.StatusOK, nil
}
