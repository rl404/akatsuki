package entity

import (
	"context"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
)

// UserAnimeFromMal to convert mal to user anime.
func UserAnimeFromMal(ctx context.Context, username string, anime mal.UserAnime) (*UserAnime, error) {
	startY, startM, startD, err := utils.SplitDate(anime.Status.StartDate)
	if err != nil {
		return nil, errors.Wrap(ctx, err)
	}

	endY, endM, endD, err := utils.SplitDate(anime.Status.FinishDate)
	if err != nil {
		return nil, errors.Wrap(ctx, err)
	}

	return &UserAnime{
		Username:     username,
		AnimeID:      int64(anime.Anime.ID),
		Status:       malToStatus(anime.Status.Status),
		Score:        anime.Status.Score,
		Episode:      anime.Status.NumEpisodesWatched,
		StartDay:     startD,
		StartMonth:   startM,
		StartYear:    startY,
		EndDay:       endD,
		EndMonth:     endM,
		EndYear:      endY,
		Priority:     malToPriority(anime.Status.Priority),
		IsRewatching: anime.Status.IsRewatching,
		RewatchCount: anime.Status.NumTimesRewatched,
		RewatchValue: malToRewatchValue(anime.Status.RewatchValue),
		Tags:         anime.Status.Tags,
		Comment:      anime.Status.Comments,
	}, nil
}

func malToStatus(s mal.AnimeStatus) Status {
	return map[mal.AnimeStatus]Status{
		mal.AnimeStatusWatching:    StatusWatching,
		mal.AnimeStatusCompleted:   StatusCompleted,
		mal.AnimeStatusOnHold:      StatusOnHold,
		mal.AnimeStatusDropped:     StatusDropped,
		mal.AnimeStatusPlanToWatch: StatusPlanned,
	}[s]
}

func malToPriority(p int) Priority {
	return map[int]Priority{
		0: PriorityLow,
		1: PriorityMedium,
		2: PriorityHigh,
	}[p]
}

func malToRewatchValue(v int) RewatchValue {
	return map[int]RewatchValue{
		1: RewatchValueVeryLow,
		2: RewatchValueLow,
		3: RewatchValueMedium,
		4: RewatchValueHigh,
		5: RewatchValueVeryHigh,
	}[v]
}
