package entity_test

import (
	"context"
	"testing"
	"time"

	"github.com/rl404/akatsuki/internal/domain/anime/entity"
	"github.com/rl404/nagato"
	"github.com/stretchr/testify/assert"
)

func TestAnimeFromMal(t *testing.T) {
	anime := nagato.Anime{
		ID:    1,
		Title: "title",
		AlternativeTitles: nagato.AlternativeTitles{
			Synonyms: []string{"synonym"},
			English:  "english",
			Japanese: "japanese",
		},
		MainPicture: nagato.Picture{
			Medium: "main-pic-medium",
		},
		Pictures: []nagato.Picture{{
			Medium: "pic-medium",
		}},
		StartDate: nagato.Date{
			Year:  2024,
			Month: 2,
			Day:   1,
		},
		EndDate: nagato.Date{
			Year:  2024,
			Month: 3,
			Day:   2,
		},
		Synopsis:               "synopsis",
		NSFW:                   nagato.NsfwWhite,
		MediaType:              nagato.MediaONA,
		Status:                 nagato.StatusCurrentlyAiring,
		NumEpisodes:            5,
		AverageEpisodeDuration: 10 * time.Minute,
		StartSeason: nagato.Season{
			Season: nagato.SeasonFall,
			Year:   2024,
		},
		Broadcast: nagato.Broadcast{
			DayOfTheWeek: nagato.DayFriday,
			StartTime:    "22:00",
		},
		Source:          nagato.Source4KomaManga,
		Rating:          nagato.RatingG,
		Background:      "background",
		Mean:            5.6,
		Rank:            200,
		Popularity:      300,
		NumListUsers:    140,
		NumScoringUsers: 80,
		Statistics: nagato.Statistic{
			Status: nagato.StatisticStatus{
				Watching:    1,
				Completed:   2,
				OnHold:      3,
				Dropped:     4,
				PlanToWatch: 5,
			},
		},
		Genres: []nagato.Genre{{
			ID:   1,
			Name: "genre-name",
		}},
		RelatedAnime: []nagato.RelatedAnime{{
			Anime: nagato.Anime{
				ID: 2,
			},
			RelationType: nagato.RelationAlternativeSetting,
		}},
		Studios: []nagato.Studio{{
			ID: 3,
		}},
	}

	res := entity.AnimeFromMal(context.Background(), &anime)
	assert.Equal(t, int64(anime.ID), res.ID)
	assert.Equal(t, anime.Title, res.Title)
	assert.Equal(t, entity.AlternativeTitle{
		Synonyms: anime.AlternativeTitles.Synonyms,
		English:  anime.AlternativeTitles.English,
		Japanese: anime.AlternativeTitles.Japanese,
	}, res.AlternativeTitle)
	assert.Equal(t, anime.MainPicture.Medium, res.Picture)
	assert.Equal(t, entity.Date{
		Day:   anime.StartDate.Day,
		Month: anime.StartDate.Month,
		Year:  anime.StartDate.Year,
	}, res.StartDate)
	assert.Equal(t, entity.Date{
		Day:   anime.EndDate.Day,
		Month: anime.EndDate.Month,
		Year:  anime.EndDate.Year,
	}, res.EndDate)
	assert.Equal(t, anime.Synopsis, res.Synopsis)
	assert.Equal(t, false, res.NSFW)
	assert.Equal(t, entity.TypeONA, res.Type)
	assert.Equal(t, entity.StatusReleasing, res.Status)
	assert.Equal(t, entity.Episode{
		Count:    anime.NumEpisodes,
		Duration: int(anime.AverageEpisodeDuration / time.Second),
	}, res.Episode)
	assert.Equal(t, entity.SeasonYear{
		Season: entity.SeasonFall,
		Year:   anime.StartSeason.Year,
	}, res.Season)
	assert.Equal(t, entity.Broadcast{
		Day:  entity.DayFriday,
		Time: anime.Broadcast.StartTime,
	}, res.Broadcast)
	assert.Equal(t, entity.Source4Koma, res.Source)
	assert.Equal(t, entity.RatingG, res.Rating)
	assert.Equal(t, anime.Background, res.Background)
	assert.Equal(t, anime.Mean, res.Mean)
	assert.Equal(t, anime.Rank, res.Rank)
	assert.Equal(t, anime.Popularity, res.Popularity)
	assert.Equal(t, anime.NumListUsers, res.Member)
	assert.Equal(t, anime.NumScoringUsers, res.Voter)
	assert.Equal(t, entity.Stats{
		Status: entity.StatsStatus{
			Watching:  int(anime.Statistics.Status.Watching),
			Completed: int(anime.Statistics.Status.Completed),
			OnHold:    int(anime.Statistics.Status.OnHold),
			Dropped:   int(anime.Statistics.Status.Dropped),
			Planned:   int(anime.Statistics.Status.PlanToWatch),
		},
	}, res.Stats)
	assert.Equal(t, []int64{1}, res.GenreIDs)
	assert.Equal(t, []int64{3}, res.StudioIDs)
	assert.Equal(t, []entity.Related{{
		ID:       int64(anime.RelatedAnime[0].Anime.ID),
		Relation: entity.RelationAlternativeSetting,
	}}, res.Related)
}
