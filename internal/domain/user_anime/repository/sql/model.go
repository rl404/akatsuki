package sql

import (
	"encoding/json"
	"time"

	"github.com/rl404/akatsuki/internal/domain/user_anime/entity"
	"gorm.io/gorm"
)

// UserAnime is user_anime database model.
type UserAnime struct {
	ID           int64
	Username     string `gorm:"index:username_index"`
	AnimeID      int64  `gorm:"index:anime_id_index"`
	Status       entity.Status
	Score        int
	Episode      int
	StartDay     int
	StartMonth   int
	StartYear    int
	EndDay       int
	EndMonth     int
	EndYear      int
	Priority     entity.Priority
	IsRewatching bool
	RewatchCount int
	RewatchValue entity.RewatchValue
	Tags         string
	Comment      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

func (ua *UserAnime) toEntity() *entity.UserAnime {
	var tags []string
	_ = json.Unmarshal([]byte(ua.Tags), &tags)

	return &entity.UserAnime{
		ID:           ua.ID,
		Username:     ua.Username,
		AnimeID:      ua.AnimeID,
		Status:       ua.Status,
		Score:        ua.Score,
		Episode:      ua.Episode,
		StartDay:     ua.StartDay,
		StartMonth:   ua.StartMonth,
		StartYear:    ua.StartYear,
		EndDay:       ua.EndDay,
		EndMonth:     ua.EndMonth,
		EndYear:      ua.EndYear,
		Priority:     ua.Priority,
		IsRewatching: ua.IsRewatching,
		RewatchCount: ua.RewatchCount,
		RewatchValue: ua.RewatchValue,
		Tags:         tags,
		Comment:      ua.Comment,
		UpdatedAt:    ua.UpdatedAt,
	}
}

func (sql *SQL) userAnimeToEntities(data []UserAnime) []*entity.UserAnime {
	a := make([]*entity.UserAnime, len(data))
	for i, aa := range data {
		a[i] = aa.toEntity()
	}
	return a
}

func (sql *SQL) userAnimeFromEntity(anime entity.UserAnime) *UserAnime {
	tags, _ := json.Marshal(anime.Tags)

	return &UserAnime{
		ID:           anime.ID,
		Username:     anime.Username,
		AnimeID:      anime.AnimeID,
		Status:       anime.Status,
		Score:        anime.Score,
		Episode:      anime.Episode,
		StartDay:     anime.StartDay,
		StartMonth:   anime.StartMonth,
		StartYear:    anime.StartYear,
		EndDay:       anime.EndDay,
		EndMonth:     anime.EndMonth,
		EndYear:      anime.EndYear,
		Priority:     anime.Priority,
		IsRewatching: anime.IsRewatching,
		RewatchCount: anime.RewatchCount,
		RewatchValue: anime.RewatchValue,
		Tags:         string(tags),
		Comment:      anime.Comment,
	}
}
