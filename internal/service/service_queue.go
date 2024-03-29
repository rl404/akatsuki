package service

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
)

// QueueOldReleasingAnime to queue old releasing anime data.
func (s *service) QueueOldReleasingAnime(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	ids, code, err := s.anime.GetOldReleasingIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	for i := 0; i < len(ids) && cnt < limit; i, cnt = i+1, cnt+1 {
		if err := s.publisher.PublishParseAnime(ctx, ids[i], false); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
	}

	return cnt, http.StatusOK, nil
}

// QueueOldFinishedAnime to queue old finished anime data.
func (s *service) QueueOldFinishedAnime(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	ids, code, err := s.anime.GetOldFinishedIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	for i := 0; i < len(ids) && cnt < limit; i, cnt = i+1, cnt+1 {
		if err := s.publisher.PublishParseAnime(ctx, ids[i], false); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
	}

	return cnt, http.StatusOK, nil
}

// QueueOldNotYetAnime to queue old not yet released anime data.
func (s *service) QueueOldNotYetAnime(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	ids, code, err := s.anime.GetOldFinishedIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	for i := 0; i < len(ids) && cnt < limit; i, cnt = i+1, cnt+1 {
		if err := s.publisher.PublishParseAnime(ctx, ids[i], false); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
	}

	return cnt, http.StatusOK, nil
}

// QueueMissingAnime to queue missing anime.
func (s *service) QueueMissingAnime(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	// Get max id.
	maxID, code, err := s.anime.GetMaxID(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	// Get all existing anime id.
	animeIDs, code, err := s.anime.GetIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	// Get all empty anime id,
	emptyIDs, code, err := s.emptyID.GetIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
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

		if err := s.publisher.PublishParseAnime(ctx, id, false); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}

		cnt++
	}

	return cnt, http.StatusOK, nil
}

// QueueOldUserAnime to queue old user anime.
func (s *service) QueueOldUserAnime(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	usernames, code, err := s.userAnime.GetOldUsernames(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	for i := 0; i < len(usernames) && cnt < limit; i, cnt = i+1, cnt+1 {
		if err := s.publisher.PublishParseUserAnime(ctx, usernames[i], "", false); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
	}

	return cnt, http.StatusOK, nil
}
