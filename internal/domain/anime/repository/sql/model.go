package sql

import (
	"encoding/json"
	"time"

	"github.com/rl404/akatsuki/internal/domain/anime/entity"
	"gorm.io/gorm"
)

// Anime is anime database model.
type Anime struct {
	ID              int64 `gorm:"primaryKey"`
	Title           string
	TitleSynonym    string
	TitleEnglish    string
	TitleJapanese   string
	Picture         string
	StartDay        int
	StartMonth      int
	StartYear       int
	EndDay          int
	EndMonth        int
	EndYear         int
	Synopsis        string
	NSFW            bool
	Type            entity.Type
	Status          entity.Status
	Episode         int
	EpisodeDuration int // in seconds
	Season          entity.Season
	SeasonYear      int
	BroadcastDay    entity.Day
	BroadcastTime   string
	Source          entity.Source
	Rating          entity.Rating
	Background      string

	// Stats.
	Mean          float64
	Rank          int
	Popularity    int
	Member        int
	Voter         int
	UserWatching  int
	UserCompleted int
	UserOnHold    int
	UserDropped   int
	UserPlanned   int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// AnimeGenre is anime_genre database model.
type AnimeGenre struct {
	AnimeID int64 `gorm:"primaryKey"`
	GenreID int64 `gorm:"primaryKey"`
}

// AnimePicture is anime_picture database model.
type AnimePicture struct {
	AnimeID int64 `gorm:"index"`
	URL     string
}

// AnimeRelated is anime_related database model.
type AnimeRelated struct {
	AnimeID1 int64           `gorm:"primaryKey"`
	AnimeID2 int64           `gorm:"primaryKey"`
	Relation entity.Relation `gorm:"primaryKey"`
}

// AnimeStudio is anime_studio database model.
type AnimeStudio struct {
	AnimeID  int64 `gorm:"primaryKey"`
	StudioID int64 `gorm:"primaryKey"`
}

// AnimeStatsHistory is anime_stats_history database model.
type AnimeStatsHistory struct {
	ID            int64
	AnimeID       int64 `gorm:"index"`
	Mean          float64
	Rank          int
	Popularity    int
	Member        int
	Voter         int
	UserWatching  int
	UserCompleted int
	UserOnHold    int
	UserDropped   int
	UserPlanned   int
	CreatedAt     time.Time
}

func (a *Anime) toEntity() *entity.Anime {
	var synonyms []string
	_ = json.Unmarshal([]byte(a.TitleSynonym), &synonyms)

	return &entity.Anime{
		ID:    a.ID,
		Title: a.Title,
		AlternativeTitle: entity.AlternativeTitle{
			Synonyms: synonyms,
			English:  a.TitleEnglish,
			Japanese: a.TitleJapanese,
		},
		Picture: a.Picture,
		StartDate: entity.Date{
			Year:  a.StartYear,
			Month: a.StartMonth,
			Day:   a.StartDay,
		},
		EndDate: entity.Date{
			Year:  a.EndYear,
			Month: a.EndMonth,
			Day:   a.EndDay,
		},
		Synopsis: a.Synopsis,
		NSFW:     a.NSFW,
		Type:     a.Type,
		Status:   a.Status,
		Episode: entity.Episode{
			Count:    a.Episode,
			Duration: a.EpisodeDuration,
		},
		Season: entity.SeasonYear{
			Season: a.Season,
			Year:   a.SeasonYear,
		},
		Broadcast: entity.Broadcast{
			Day:  a.BroadcastDay,
			Time: a.BroadcastTime,
		},
		Source:     a.Source,
		Rating:     a.Rating,
		Background: a.Background,
		Mean:       a.Mean,
		Rank:       a.Rank,
		Popularity: a.Popularity,
		Member:     a.Member,
		Voter:      a.Voter,
		Stats: entity.Stats{
			Status: entity.StatsStatus{
				Watching:  a.UserWatching,
				Completed: a.UserCompleted,
				OnHold:    a.UserOnHold,
				Dropped:   a.UserDropped,
				Planned:   a.UserPlanned,
			},
		},
		UpdatedAt: a.UpdatedAt,
	}
}

func (sql *SQL) animeToEntities(data []Anime) []*entity.Anime {
	a := make([]*entity.Anime, len(data))
	for i, aa := range data {
		a[i] = aa.toEntity()
	}
	return a
}

func (sql *SQL) animeFromEntity(anime entity.Anime) *Anime {
	synonym, _ := json.Marshal(anime.AlternativeTitle.Synonyms)

	return &Anime{
		ID:              anime.ID,
		Title:           anime.Title,
		TitleSynonym:    string(synonym),
		TitleEnglish:    anime.AlternativeTitle.English,
		TitleJapanese:   anime.AlternativeTitle.Japanese,
		Picture:         anime.Picture,
		StartDay:        anime.StartDate.Day,
		StartMonth:      anime.StartDate.Month,
		StartYear:       anime.StartDate.Year,
		EndDay:          anime.EndDate.Day,
		EndMonth:        anime.EndDate.Month,
		EndYear:         anime.EndDate.Year,
		Synopsis:        anime.Synopsis,
		NSFW:            anime.NSFW,
		Type:            anime.Type,
		Status:          anime.Status,
		Episode:         anime.Episode.Count,
		EpisodeDuration: anime.Episode.Duration,
		Season:          anime.Season.Season,
		SeasonYear:      anime.Season.Year,
		BroadcastDay:    anime.Broadcast.Day,
		BroadcastTime:   anime.Broadcast.Time,
		Source:          anime.Source,
		Rating:          anime.Rating,
		Background:      anime.Background,
		Mean:            anime.Mean,
		Rank:            anime.Rank,
		Popularity:      anime.Popularity,
		Member:          anime.Member,
		Voter:           anime.Voter,
		UserWatching:    anime.Stats.Status.Watching,
		UserCompleted:   anime.Stats.Status.Completed,
		UserOnHold:      anime.Stats.Status.OnHold,
		UserDropped:     anime.Stats.Status.Dropped,
		UserPlanned:     anime.Stats.Status.Planned,
	}
}

func (sql *SQL) animeGenreFromEntity(anime entity.Anime) []AnimeGenre {
	ag := make([]AnimeGenre, len(anime.GenreIDs))
	for i, id := range anime.GenreIDs {
		ag[i] = AnimeGenre{
			AnimeID: anime.ID,
			GenreID: id,
		}
	}
	return ag
}

func (sql *SQL) animePictureFromEntity(anime entity.Anime) []AnimePicture {
	ap := make([]AnimePicture, len(anime.Pictures))
	for i, p := range anime.Pictures {
		ap[i] = AnimePicture{
			AnimeID: anime.ID,
			URL:     p,
		}
	}
	return ap
}

func (sql *SQL) animeRelatedFromEntity(anime entity.Anime) []AnimeRelated {
	ar := make([]AnimeRelated, len(anime.Related))
	for i, r := range anime.Related {
		ar[i] = AnimeRelated{
			AnimeID1: anime.ID,
			AnimeID2: r.ID,
			Relation: r.Relation,
		}
	}
	return ar
}

func (sql *SQL) animeStudioFromEntity(anime entity.Anime) []AnimeStudio {
	as := make([]AnimeStudio, len(anime.StudioIDs))
	for i, s := range anime.StudioIDs {
		as[i] = AnimeStudio{
			AnimeID:  anime.ID,
			StudioID: s,
		}
	}
	return as
}

func (sql *SQL) animeStatsFromEntity(anime entity.Anime) *AnimeStatsHistory {
	return &AnimeStatsHistory{
		AnimeID:       anime.ID,
		Mean:          anime.Mean,
		Rank:          anime.Rank,
		Popularity:    anime.Popularity,
		Member:        anime.Member,
		Voter:         anime.Voter,
		UserWatching:  anime.Stats.Status.Watching,
		UserCompleted: anime.Stats.Status.Completed,
		UserOnHold:    anime.Stats.Status.OnHold,
		UserDropped:   anime.Stats.Status.Dropped,
		UserPlanned:   anime.Stats.Status.Planned,
	}
}

func (ar *AnimeRelated) toEntity() *entity.AnimeRelated {
	return &entity.AnimeRelated{
		AnimeID1: ar.AnimeID1,
		AnimeID2: ar.AnimeID2,
		Relation: ar.Relation,
	}
}

func (sql *SQL) animeRelatedToEntities(data []AnimeRelated) []*entity.AnimeRelated {
	ar := make([]*entity.AnimeRelated, len(data))
	for i, aa := range data {
		ar[i] = aa.toEntity()
	}
	return ar
}
