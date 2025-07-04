package entity

import (
	"context"
	"time"

	"github.com/rl404/nagato"
)

// AnimeFromMal to convert mal to anime.
func AnimeFromMal(ctx context.Context, anime *nagato.Anime) Anime {
	picture := anime.MainPicture.Large
	if picture == "" {
		picture = anime.MainPicture.Medium
	}

	genreIDs := make([]int64, len(anime.Genres))
	for i, g := range anime.Genres {
		genreIDs[i] = int64(g.ID)
	}

	pictures := make([]string, len(anime.Pictures))
	for i, p := range anime.Pictures {
		pictures[i] = p.Large
		if pictures[i] == "" {
			pictures[i] = p.Medium
		}
	}

	related := make([]Related, len(anime.RelatedAnime))
	for i, r := range anime.RelatedAnime {
		related[i] = Related{
			ID:       int64(r.Anime.ID),
			Relation: malToRelation(r.RelationType),
		}
	}

	studioIDs := make([]int64, len(anime.Studios))
	for i, s := range anime.Studios {
		studioIDs[i] = int64(s.ID)
	}

	return Anime{
		ID:    int64(anime.ID),
		Title: anime.Title,
		AlternativeTitle: AlternativeTitle{
			Synonyms: anime.AlternativeTitles.Synonyms,
			English:  anime.AlternativeTitles.English,
			Japanese: anime.AlternativeTitles.Japanese,
		},
		Picture: picture,
		StartDate: Date{
			Day:   anime.StartDate.Day,
			Month: anime.StartDate.Month,
			Year:  anime.StartDate.Year,
		},
		EndDate: Date{
			Day:   anime.EndDate.Day,
			Month: anime.EndDate.Month,
			Year:  anime.EndDate.Year,
		},
		Synopsis: anime.Synopsis,
		NSFW:     anime.NSFW != "white",
		Type:     malToType(anime.MediaType),
		Status:   malToStatus(anime.Status),
		Episode: Episode{
			Count:    anime.NumEpisodes,
			Duration: int(anime.AverageEpisodeDuration / time.Second),
		},
		Season: SeasonYear{
			Season: malToSeason(anime.StartSeason.Season),
			Year:   anime.StartSeason.Year,
		},
		Broadcast: Broadcast{
			Day:  malToDay(anime.Broadcast.DayOfTheWeek),
			Time: anime.Broadcast.StartTime,
		},
		Source:     malToSource(anime.Source),
		Rating:     malToRating(anime.Rating),
		Background: anime.Background,
		Mean:       anime.Mean,
		Rank:       anime.Rank,
		Popularity: anime.Popularity,
		Member:     anime.NumListUsers,
		Voter:      anime.NumScoringUsers,
		Stats: Stats{
			Status: StatsStatus{
				Watching:  int(anime.Statistics.Status.Watching),
				Completed: int(anime.Statistics.Status.Completed),
				OnHold:    int(anime.Statistics.Status.OnHold),
				Dropped:   int(anime.Statistics.Status.Dropped),
				Planned:   int(anime.Statistics.Status.PlanToWatch),
			},
		},
		GenreIDs:  genreIDs,
		Pictures:  pictures,
		Related:   related,
		StudioIDs: studioIDs,
	}
}

func malToType(t nagato.MediaType) Type {
	return map[nagato.MediaType]Type{
		"":                    TypeUnknown,
		nagato.MediaTV:        TypeTV,
		nagato.MediaOVA:       TypeOVA,
		nagato.MediaMovie:     TypeMovie,
		nagato.MediaSpecial:   TypeSpecial,
		nagato.MediaONA:       TypeONA,
		nagato.MediaMusic:     TypeMusic,
		nagato.MediaCM:        TypeCM,
		nagato.MediaPV:        TypePV,
		nagato.MediaTVSpecial: TypeTVSpecial,
	}[t]
}

func malToStatus(s nagato.StatusType) Status {
	return map[nagato.StatusType]Status{
		nagato.StatusFinishedAiring:  StatusFinished,
		nagato.StatusCurrentlyAiring: StatusReleasing,
		nagato.StatusNotYetAired:     StatusNotYet,
	}[s]
}

func malToSeason(s nagato.SeasonType) Season {
	return map[nagato.SeasonType]Season{
		nagato.SeasonWinter: SeasonWinter,
		nagato.SeasonSpring: SeasonSpring,
		nagato.SeasonSummer: SeasonSummer,
		nagato.SeasonFall:   SeasonFall,
	}[s]
}

func malToDay(d nagato.DayType) Day {
	return map[nagato.DayType]Day{
		nagato.DayMonday:    DayMonday,
		nagato.DayTuesday:   DayTuesday,
		nagato.DayWednesday: DayWednesday,
		nagato.DayThursday:  DayThursday,
		nagato.DayFriday:    DayFriday,
		nagato.DaySaturday:  DaySaturday,
		nagato.DaySunday:    DaySunday,
		nagato.DayOther:     DayOther,
	}[d]
}

func malToSource(s nagato.SourceType) Source {
	return map[nagato.SourceType]Source{
		nagato.SourceOriginal:     SourceOriginal,
		nagato.SourceManga:        SourceManga,
		nagato.Source4KomaManga:   Source4Koma,
		nagato.SourceWebManga:     SourceWebManga,
		nagato.SourceDigitalManga: SourceDigitalManga,
		nagato.SourceNovel:        SourceNovel,
		nagato.SourceLightNovel:   SourceLightNovel,
		nagato.SourceVisualNovel:  SourceVisualNovel,
		nagato.SourceGame:         SourceGame,
		nagato.SourceCardGame:     SourceCardGame,
		nagato.SourceBook:         SourceBook,
		nagato.SourcePictureBook:  SourcePictureBook,
		nagato.SourceRadio:        SourceRadio,
		nagato.SourceMusic:        SourceMusic,
		nagato.SourceWebNovel:     SourceWebNovel,
		nagato.SourceMixedMedia:   SourceMixedMedia,
		nagato.SourceOther:        SourceOther,
	}[s]
}

func malToRating(r nagato.RatingType) Rating {
	return map[nagato.RatingType]Rating{
		nagato.RatingG:     RatingG,
		nagato.RatingPG:    RatingPG,
		nagato.RatingPG13:  RatingPG13,
		nagato.RatingR:     RatingR,
		nagato.RatingRPlus: RatingRPlus,
		nagato.RatingRX:    RatingRX,
	}[r]
}

func malToRelation(r nagato.RelationType) Relation {
	return map[nagato.RelationType]Relation{
		nagato.RelationSequel:             RelationSequel,
		nagato.RelationPrequel:            RelationPrequel,
		nagato.RelationAlternativeSetting: RelationAlternativeSetting,
		nagato.RelationAlternativeVersion: RelationAlternativeVersion,
		nagato.RelationSideStory:          RelationSideStory,
		nagato.RelationParentStory:        RelationParentStory,
		nagato.RelationSummary:            RelationSummary,
		nagato.RelationFullStory:          RelationFullStory,
		nagato.RelationSpinOff:            RelationSpinOff,
		nagato.RelationOther:              RelationOther,
		nagato.RelationCharacter:          RelationCharacter,
	}[r]
}
