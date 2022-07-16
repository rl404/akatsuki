package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	_consumer "github.com/rl404/akatsuki/internal/delivery/consumer"
	animeRepository "github.com/rl404/akatsuki/internal/domain/anime/repository"
	animeCache "github.com/rl404/akatsuki/internal/domain/anime/repository/cache"
	animeSQL "github.com/rl404/akatsuki/internal/domain/anime/repository/sql"
	emptyIDRepository "github.com/rl404/akatsuki/internal/domain/empty_id/repository"
	emptyIDCache "github.com/rl404/akatsuki/internal/domain/empty_id/repository/cache"
	emptyIDSQL "github.com/rl404/akatsuki/internal/domain/empty_id/repository/sql"
	genreRepository "github.com/rl404/akatsuki/internal/domain/genre/repository"
	genreCache "github.com/rl404/akatsuki/internal/domain/genre/repository/cache"
	genreSQL "github.com/rl404/akatsuki/internal/domain/genre/repository/sql"
	malRepository "github.com/rl404/akatsuki/internal/domain/mal/repository"
	malClient "github.com/rl404/akatsuki/internal/domain/mal/repository/client"
	publisherRepository "github.com/rl404/akatsuki/internal/domain/publisher/repository"
	publisherPubsub "github.com/rl404/akatsuki/internal/domain/publisher/repository/pubsub"
	studioRepository "github.com/rl404/akatsuki/internal/domain/studio/repository"
	studioCache "github.com/rl404/akatsuki/internal/domain/studio/repository/cache"
	studioSQL "github.com/rl404/akatsuki/internal/domain/studio/repository/sql"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/pubsub"
)

func consumer() error {
	// Get config.
	cfg, err := getConfig()
	if err != nil {
		return err
	}
	utils.Info("config initialized")

	// Init cache.
	c, err := cache.New(cacheType[cfg.Cache.Dialect], cfg.Cache.Address, cfg.Cache.Password, cfg.Cache.Time)
	if err != nil {
		return err
	}
	utils.Info("cache initialized")
	defer c.Close()

	// Init db.
	db, err := newDB(cfg.DB)
	if err != nil {
		return err
	}
	utils.Info("database initialized")
	tmp, _ := db.DB()
	defer tmp.Close()

	// Init pubsub.
	ps, err := pubsub.New(pubsubType[cfg.PubSub.Dialect], cfg.PubSub.Address, cfg.PubSub.Password)
	if err != nil {
		return err
	}
	utils.Info("pubsub initialized")
	defer ps.Close()

	// Init anime.
	var anime animeRepository.Repository
	anime = animeSQL.New(db, time.Duration(cfg.App.OldAge)*24*time.Hour)
	anime = animeCache.New(c, anime)
	utils.Info("repository anime initialized")

	// Init genre.
	var genre genreRepository.Repository
	genre = genreSQL.New(db)
	genre = genreCache.New(c, genre)
	utils.Info("repository genre initialized")

	// Init studio.
	var studio studioRepository.Repository
	studio = studioSQL.New(db)
	studio = studioCache.New(c, studio)
	utils.Info("repository studio initialized")

	// Init empty id.
	var emptyID emptyIDRepository.Repository
	emptyID = emptyIDSQL.New(db)
	emptyID = emptyIDCache.New(c, emptyID)
	utils.Info("repository empty id initialized")

	// Init mal.
	var mal malRepository.Repository = malClient.New(cfg.Mal.ClientID)
	utils.Info("repository mal initialized")

	// Init publisher.
	var publisher publisherRepository.Repository = publisherPubsub.New(ps, pubsubTopic)
	utils.Info("repository publisher initialized")

	// Init service.
	service := service.New(anime, genre, studio, emptyID, publisher, mal)
	utils.Info("service initialized")

	// Init consumer.
	consumer, err := _consumer.New(service, ps, pubsubTopic)
	if err != nil {
		return err
	}
	utils.Info("consumer initialized")
	defer consumer.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Start subscribe.
	if err := consumer.Subscribe(); err != nil {
		return err
	}

	utils.Info("consumer ready")
	<-sigChan

	return nil
}
