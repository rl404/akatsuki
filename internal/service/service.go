package service

import (
	"context"

	"github.com/nstratos/go-myanimelist/mal"
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
	GetAnimeByID(ctx context.Context, id int64) (*Anime, int, error)

	GetUserAnime(ctx context.Context, data GetUserAnimeRequest) ([]UserAnime, *Pagination, int, error)

	GetMalAnimeByID(ctx context.Context, id int) (*mal.Anime, int, error)
	GetMalUserAnime(ctx context.Context, data GetMalUserAnimeRequest) ([]mal.UserAnime, int, error)

	ConsumeMessage(ctx context.Context, msg entity.Message) error

	QueueOldReleasingAnime(ctx context.Context, limit int) (int, int, error)
	QueueOldFinishedAnime(ctx context.Context, limit int) (int, int, error)
	QueueOldNotYetAnime(ctx context.Context, limit int) (int, int, error)
	QueueMissingAnime(ctx context.Context, limit int) (int, int, error)
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
