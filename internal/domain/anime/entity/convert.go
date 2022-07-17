package entity

import (
	"strconv"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rl404/akatsuki/internal/utils"
)

// AnimeFromMal to convert mal to anime.
func AnimeFromMal(anime *mal.Anime) (*Anime, error) {
	picture := anime.MainPicture.Large
	if picture == "" {
		picture = anime.MainPicture.Medium
	}

	startY, startM, startD, err := utils.SplitDate(anime.StartDate)
	if err != nil {
		return nil, err
	}

	endY, endM, endD, err := utils.SplitDate(anime.EndDate)
	if err != nil {
		return nil, err
	}

	watching, err := strconv.Atoi(anime.Statistics.Status.Watching)
	if err != nil {
		return nil, err
	}

	completed, err := strconv.Atoi(anime.Statistics.Status.Completed)
	if err != nil {
		return nil, err
	}

	onHold, err := strconv.Atoi(anime.Statistics.Status.OnHold)
	if err != nil {
		return nil, err
	}

	dropped, err := strconv.Atoi(anime.Statistics.Status.Dropped)
	if err != nil {
		return nil, err
	}

	planned, err := strconv.Atoi(anime.Statistics.Status.PlanToWatch)
	if err != nil {
		return nil, err
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
			ID:       int64(r.Node.ID),
			Relation: malToRelation(r.RelationType),
		}
	}

	studioIDs := make([]int64, len(anime.Studios))
	for i, s := range anime.Studios {
		studioIDs[i] = int64(s.ID)
	}

	return &Anime{
		ID:    int64(anime.ID),
		Title: anime.Title,
		AlternativeTitle: AlternativeTitle{
			Synonyms: anime.AlternativeTitles.Synonyms,
			English:  anime.AlternativeTitles.En,
			Japanese: anime.AlternativeTitles.Ja,
		},
		Picture: picture,
		StartDate: Date{
			Day:   startD,
			Month: startM,
			Year:  startY,
		},
		EndDate: Date{
			Day:   endD,
			Month: endM,
			Year:  endY,
		},
		Synopsis: anime.Synopsis,
		NSFW:     anime.NSFW != "white",
		Type:     malToType(anime.MediaType),
		Status:   malToStatus(anime.Status),
		Episode: Episode{
			Count:    anime.NumEpisodes,
			Duration: anime.AverageEpisodeDuration,
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
				Watching:  watching,
				Completed: completed,
				OnHold:    onHold,
				Dropped:   dropped,
				Planned:   planned,
			},
		},
		GenreIDs:  genreIDs,
		Pictures:  pictures,
		Related:   related,
		StudioIDs: studioIDs,
	}, nil
}

func malToType(t string) Type {
	return map[string]Type{
		"":        TypeUnknown,
		"tv":      TypeTV,
		"ova":     TypeOVA,
		"ona":     TypeONA,
		"movie":   TypeMovie,
		"special": TypeSpecial,
		"music":   TypeMusic,
	}[t]
}

func malToStatus(s string) Status {
	return map[string]Status{
		"finished_airing":  StatusFinished,
		"currently_airing": StatusReleasing,
		"not_yet_aired":    StatusNotYet,
	}[s]
}

func malToSeason(s string) Season {
	return map[string]Season{
		"winter": SeasonWinter,
		"spring": SeasonSpring,
		"summer": SeasonSummer,
		"fall":   SeasonFall,
	}[s]
}

func malToDay(d string) Day {
	return map[string]Day{
		"monday":    DayMonday,
		"tuesday":   DayTuesday,
		"wednesday": DayWednesday,
		"thurday":   DayThurday,
		"friday":    DayFriday,
		"saturday":  DaySaturday,
		"sunday":    DaySunday,
		"other":     DayOther,
	}[d]
}

func malToSource(s string) Source {
	return map[string]Source{
		"original":      SourceOriginal,
		"manga":         SourceManga,
		"4_koma_manga":  Source4Koma,
		"web_manga":     SourceWebManga,
		"digital_manga": SourceDigitalManga,
		"novel":         SourceNovel,
		"light_novel":   SourceLightNovel,
		"visual_novel":  SourceVisualNovel,
		"game":          SourceGame,
		"card_game":     SourceCardGame,
		"book":          SourceBook,
		"picture_book":  SourcePictureBook,
		"radio":         SourceRadio,
		"music":         SourceMusic,
		"other":         SourceOther,
		"web_novel":     SourceWebNovel,
		"mixed_media":   SourceMixedMedia,
	}[s]
}

func malToRating(r string) Rating {
	return map[string]Rating{
		"g":     RatingG,
		"pg":    RatingPG,
		"pg_13": RatingPG13,
		"r":     RatingR,
		"r+":    RatingRPlus,
		"rx":    RatingRX,
	}[r]
}

func malToRelation(r string) Relation {
	return map[string]Relation{
		"sequel":              RelationSequel,
		"prequel":             RelationPrequel,
		"alternative_setting": RelationAlternativeSetting,
		"alternative_version": RelationAlternativeVersion,
		"side_story":          RelationSideStory,
		"parent_story":        RelationParentStory,
		"summary":             RelationSummary,
		"full_story":          RelationFullStory,
	}[r]
}
