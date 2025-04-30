package entity

import (
	"context"

	"github.com/rl404/nagato"
)

// UserAnimeFromMal to convert mal to user anime.
func UserAnimeFromMal(ctx context.Context, username string, anime nagato.UserAnime) UserAnime {
	return UserAnime{
		Username:     username,
		AnimeID:      int64(anime.Anime.ID),
		Status:       malToStatus(anime.Status.Status),
		Score:        anime.Status.Score,
		Episode:      anime.Status.NumEpisodesWatched,
		StartDay:     anime.Status.StartDate.Day,
		StartMonth:   anime.Status.StartDate.Month,
		StartYear:    anime.Status.StartDate.Year,
		EndDay:       anime.Status.FinishDate.Day,
		EndMonth:     anime.Status.FinishDate.Month,
		EndYear:      anime.Status.FinishDate.Year,
		Priority:     malToPriority(anime.Status.Priority),
		IsRewatching: anime.Status.IsRewatching,
		RewatchCount: anime.Status.NumTimesRewatched,
		RewatchValue: malToRewatchValue(anime.Status.RewatchValue),
		Tags:         anime.Status.Tags,
		Comment:      anime.Status.Comments,
	}
}

func malToStatus(s nagato.UserAnimeStatusType) Status {
	return map[nagato.UserAnimeStatusType]Status{
		nagato.UserAnimeStatusWatching:    StatusWatching,
		nagato.UserAnimeStatusCompleted:   StatusCompleted,
		nagato.UserAnimeStatusOnHold:      StatusOnHold,
		nagato.UserAnimeStatusDropped:     StatusDropped,
		nagato.UserAnimeStatusPlanToWatch: StatusPlanned,
	}[s]
}

func malToPriority(p nagato.PriorityType) Priority {
	return map[nagato.PriorityType]Priority{
		nagato.PriorityLow:    PriorityLow,
		nagato.PriorityMedium: PriorityMedium,
		nagato.PriorityHigh:   PriorityHigh,
	}[p]
}

func malToRewatchValue(v nagato.RewatchValueType) RewatchValue {
	return map[nagato.RewatchValueType]RewatchValue{
		nagato.RewatchValueVeryLow:  RewatchValueVeryLow,
		nagato.RewatchValueLow:      RewatchValueLow,
		nagato.RewatchValueMedium:   RewatchValueMedium,
		nagato.RewatchValueHigh:     RewatchValueHigh,
		nagato.RewatchValueVeryHigh: RewatchValueVeryHigh,
	}[v]
}

// StrToStatus to convert string to status.
func StrToStatus(status string) Status {
	return map[string]Status{
		string(nagato.UserAnimeStatusWatching):    StatusWatching,
		string(nagato.UserAnimeStatusCompleted):   StatusCompleted,
		string(nagato.UserAnimeStatusOnHold):      StatusOnHold,
		string(nagato.UserAnimeStatusDropped):     StatusDropped,
		string(nagato.UserAnimeStatusPlanToWatch): StatusPlanned,
	}[status]
}
