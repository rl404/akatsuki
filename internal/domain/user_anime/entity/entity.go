package entity

import "time"

// UserAnime is user anime entity.
type UserAnime struct {
	ID           int64
	Username     string
	AnimeID      int64
	Status       Status
	Score        int
	Episode      int
	StartDay     int
	StartMonth   int
	StartYear    int
	EndDay       int
	EndMonth     int
	EndYear      int
	Priority     Priority
	IsRewatching bool
	RewatchCount int
	RewatchValue RewatchValue
	Tags         []string
	Comment      string
	UpdatedAt    time.Time
}

// GetUserAnimeRequest is get user anime request model.
type GetUserAnimeRequest struct {
	Username string
	Page     int
	Limit    int
}
