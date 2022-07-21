package entity

import "github.com/nstratos/go-myanimelist/mal"

// GetUserAnimeRequest is get user anime request entity.
type GetUserAnimeRequest struct {
	Username string
	Status   mal.AnimeStatus
	Sort     mal.SortAnimeList
	Limit    int
	Offset   int
}
