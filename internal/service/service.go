package service

import (
	"context"

	animeRepository "github.com/rl404/akatsuki/internal/domain/anime/repository"
	emptyIDRepository "github.com/rl404/akatsuki/internal/domain/empty_id/repository"
	genreRepository "github.com/rl404/akatsuki/internal/domain/genre/repository"
	malRepository "github.com/rl404/akatsuki/internal/domain/mal/repository"
	"github.com/rl404/akatsuki/internal/domain/publisher/entity"
	publisherRepository "github.com/rl404/akatsuki/internal/domain/publisher/repository"
	studioRepository "github.com/rl404/akatsuki/internal/domain/studio/repository"
	userAnimeRepository "github.com/rl404/akatsuki/internal/domain/user_anime/repository"
)

// Service contains functions for service.
type Service interface {
	GetAnime(ctx context.Context, data GetAnimeRequest) ([]Anime, *Pagination, int, error)
	GetAnimeByID(ctx context.Context, id int64) (*Anime, int, error)
	GetAnimeHistoriesByID(ctx context.Context, data GetAnimeHistoriesRequest) ([]AnimeHistory, int, error)
	UpdateAnimeByID(ctx context.Context, id int64) (int, error)

	GetGenres(ctx context.Context, data GetGenresRequest) ([]Genre, *Pagination, int, error)
	GetGenreByID(ctx context.Context, id int64) (*Genre, int, error)
	GetGenreHistoriesByID(ctx context.Context, data GetGenreHistoriesRequest) ([]GenreHistory, int, error)

	GetStudios(ctx context.Context, data GetStudiosRequest) ([]Studio, *Pagination, int, error)
	GetStudioByID(ctx context.Context, id int64) (*Studio, int, error)
	GetStudioHistoriesByID(ctx context.Context, data GetStudioHistoriesRequest) ([]StudioHistory, int, error)

	GetUserAnime(ctx context.Context, data GetUserAnimeRequest) ([]UserAnime, *Pagination, int, error)
	GetUserAnimeRelations(ctx context.Context, username string) (*UserAnimeRelation, int, error)
	UpdateUserAnime(ctx context.Context, username string) (int, error)

	ConsumeMessage(ctx context.Context, msg entity.Message) error

	QueueOldReleasingAnime(ctx context.Context, limit int) (int, int, error)
	QueueOldFinishedAnime(ctx context.Context, limit int) (int, int, error)
	QueueOldNotYetAnime(ctx context.Context, limit int) (int, int, error)
	QueueMissingAnime(ctx context.Context, limit int) (int, int, error)
	QueueOldUserAnime(ctx context.Context, limit int) (int, int, error)
}

type service struct {
	anime     animeRepository.Repository
	genre     genreRepository.Repository
	studio    studioRepository.Repository
	userAnime userAnimeRepository.Repository
	emptyID   emptyIDRepository.Repository
	publisher publisherRepository.Repository
	mal       malRepository.Repository
}

// New to create new service.
func New(
	anime animeRepository.Repository,
	genre genreRepository.Repository,
	studio studioRepository.Repository,
	userAnime userAnimeRepository.Repository,
	emptyID emptyIDRepository.Repository,
	publisher publisherRepository.Repository,
	mal malRepository.Repository,
) Service {
	return &service{
		anime:     anime,
		genre:     genre,
		studio:    studio,
		userAnime: userAnime,
		emptyID:   emptyID,
		publisher: publisher,
		mal:       mal,
	}
}
