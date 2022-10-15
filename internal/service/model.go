package service

import "github.com/rl404/akatsuki/internal/domain/anime/entity"

// Pagination is pagination model.
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type alternativeTitle struct {
	Synonyms []string `json:"synonyms"`
	English  string   `json:"english"`
	Japanese string   `json:"japanese"`
}

type date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type episode struct {
	Count    int `json:"count"`
	Duration int `json:"duration"`
}

type season struct {
	Season entity.Season `json:"season" swaggertype:"string"`
	Year   int           `json:"year"`
}

type broadcast struct {
	Day  entity.Day `json:"day" swaggertype:"string"`
	Time string     `json:"time"`
}

type stats struct {
	Status statsStatus `json:"status"`
}

type statsStatus struct {
	Watching  int `json:"watching"`
	Completed int `json:"completed"`
	OnHold    int `json:"on_hold"`
	Dropped   int `json:"dropped"`
	Planned   int `json:"planned"`
}

type genre struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type related struct {
	ID       int64           `json:"id"`
	Title    string          `json:"title"`
	Picture  string          `json:"picture"`
	Relation entity.Relation `json:"relation" swaggertype:"string"`
}

type studio struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (s *service) animeFromEntity(animeDB *entity.Anime) Anime {
	var anime Anime
	anime.ID = animeDB.ID
	anime.Title = animeDB.Title
	anime.AlternativeTitles = alternativeTitle{
		Synonyms: animeDB.AlternativeTitle.Synonyms,
		English:  animeDB.AlternativeTitle.English,
		Japanese: animeDB.AlternativeTitle.Japanese,
	}
	anime.Picture = animeDB.Picture
	anime.StartDate = date{
		Year:  animeDB.StartDate.Year,
		Month: animeDB.StartDate.Month,
		Day:   animeDB.StartDate.Day,
	}
	anime.EndDate = date{
		Year:  animeDB.EndDate.Year,
		Month: animeDB.EndDate.Month,
		Day:   animeDB.EndDate.Day,
	}
	anime.Synopsis = animeDB.Synopsis
	anime.Background = animeDB.Background
	anime.NSFW = animeDB.NSFW
	anime.Type = animeDB.Type
	anime.Status = animeDB.Status
	anime.Episode = episode{
		Count:    animeDB.Episode.Count,
		Duration: animeDB.Episode.Duration,
	}
	if animeDB.Season.Season != "" {
		anime.Season = &season{
			Season: animeDB.Season.Season,
			Year:   animeDB.Season.Year,
		}
	}
	if animeDB.Broadcast.Day != "" {
		anime.Broadcast = &broadcast{
			Day:  animeDB.Broadcast.Day,
			Time: animeDB.Broadcast.Time,
		}
	}
	anime.Source = animeDB.Source
	anime.Rating = animeDB.Rating
	anime.Mean = animeDB.Mean
	anime.Rank = animeDB.Rank
	anime.Popularity = animeDB.Popularity
	anime.Member = animeDB.Member
	anime.Voter = animeDB.Voter
	anime.Stats = stats{
		Status: statsStatus{
			Watching:  animeDB.Stats.Status.Watching,
			Completed: animeDB.Stats.Status.Completed,
			OnHold:    animeDB.Stats.Status.OnHold,
			Dropped:   animeDB.Stats.Status.Dropped,
			Planned:   animeDB.Stats.Status.Planned,
		},
	}
	anime.Pictures = animeDB.Pictures
	anime.Genres = []genre{}
	anime.Related = []related{}
	anime.Studios = []studio{}
	anime.UpdatedAt = animeDB.UpdatedAt
	return anime
}
