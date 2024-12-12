package utils_test

import (
	animeRepository "github.com/rl404/akatsuki/internal/domain/anime/repository"
	animeCache "github.com/rl404/akatsuki/internal/domain/anime/repository/cache"
	animeSQL "github.com/rl404/akatsuki/internal/domain/anime/repository/sql"
	emptyIDRepository "github.com/rl404/akatsuki/internal/domain/empty_id/repository"
	emptyIDCache "github.com/rl404/akatsuki/internal/domain/empty_id/repository/cache"
	emptyIDSQL "github.com/rl404/akatsuki/internal/domain/empty_id/repository/sql"
	genreRepository "github.com/rl404/akatsuki/internal/domain/genre/repository"
	genreCache "github.com/rl404/akatsuki/internal/domain/genre/repository/cache"
	genreSQL "github.com/rl404/akatsuki/internal/domain/genre/repository/sql"
	publisherRepository "github.com/rl404/akatsuki/internal/domain/publisher/repository"
	publisherPubsub "github.com/rl404/akatsuki/internal/domain/publisher/repository/pubsub"
	studioRepository "github.com/rl404/akatsuki/internal/domain/studio/repository"
	studioCache "github.com/rl404/akatsuki/internal/domain/studio/repository/cache"
	studioSQL "github.com/rl404/akatsuki/internal/domain/studio/repository/sql"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/pubsub"
	"gorm.io/gorm"
)

func GetService(cfg *config, db *gorm.DB, c cache.Cacher, ps pubsub.PubSub) service.Service {
	// Init anime.
	var anime animeRepository.Repository
	anime = animeSQL.New(db, cfg.Cron.FinishedAge, cfg.Cron.ReleasingAge, cfg.Cron.NotYetAge)
	anime = animeCache.New(c, anime)

	// Init genre.
	var genre genreRepository.Repository
	genre = genreSQL.New(db)
	genre = genreCache.New(c, genre)

	// Init studio.
	var studio studioRepository.Repository
	studio = studioSQL.New(db)
	studio = studioCache.New(c, studio)

	// Init empty id.
	var emptyID emptyIDRepository.Repository
	emptyID = emptyIDSQL.New(db)
	emptyID = emptyIDCache.New(c, emptyID)

	// Init publisher.
	var publisher publisherRepository.Repository = publisherPubsub.New(ps, pubsubTopic)

	return service.New(anime, genre, studio, nil, emptyID, publisher, nil)
}
