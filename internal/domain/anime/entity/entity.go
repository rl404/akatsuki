package entity

import "time"

// Anime is entity for anime.
type Anime struct {
	ID               int64
	Title            string
	AlternativeTitle AlternativeTitle
	Picture          string
	StartDate        Date
	EndDate          Date
	Synopsis         string
	NSFW             bool
	Type             Type
	Status           Status
	Episode          Episode
	Season           SeasonYear
	Broadcast        Broadcast
	Source           Source
	Rating           Rating
	Background       string
	Mean             float64
	Rank             int
	Popularity       int
	Member           int
	Voter            int
	Stats            Stats

	// Relation.
	GenreIDs  []int64
	Pictures  []string
	Related   []Related
	StudioIDs []int64

	UpdatedAt time.Time
}

// AlternativeTitle is entity for alternative title.
type AlternativeTitle struct {
	Synonyms []string
	English  string
	Japanese string
}

// Date is entity for date.
type Date struct {
	Day   int
	Month int
	Year  int
}

// SeasonYear is entity for season and year.
type SeasonYear struct {
	Season Season
	Year   int
}

// Episode is entity for episode.
type Episode struct {
	Count    int
	Duration int
}

// Broadcast is entity for broadcast.
type Broadcast struct {
	Day  Day
	Time string
}

// Stats is entity for stats.
// Will contain score in the future?
type Stats struct {
	Status StatsStatus
}

// StatsStatus is entity for stats status.
type StatsStatus struct {
	Watching  int
	Completed int
	OnHold    int
	Dropped   int
	Planned   int
}

// Related is entity for related anime.
type Related struct {
	ID       int64
	Relation Relation
}

// AnimeRelated is entity for related anime.
type AnimeRelated struct {
	AnimeID1 int64
	AnimeID2 int64
	Relation Relation
}

// History is entity for anime history.
type History struct {
	Year          int
	Month         int
	Week          int
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
}

// GetHistoriesRequest is get histories request model.
type GetHistoriesRequest struct {
	AnimeID   int64
	StartDate *time.Time
	EndDate   *time.Time
	Group     HistoryGroup
}

// GetRequest is get request model.
type GetRequest struct {
	Title      string
	NSFW       *bool
	Type       Type
	Status     Status
	Season     Season
	SeasonYear int
	StartMean  float64
	EndMean    float64
	GenreID    int64
	StudioID   int64
	Sort       Sort
	Page       int
	Limit      int
}
