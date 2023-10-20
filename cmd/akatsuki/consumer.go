package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	_consumer "github.com/rl404/akatsuki/internal/delivery/consumer"
	animeRepository "github.com/rl404/akatsuki/internal/domain/anime/repository"
	animeSQL "github.com/rl404/akatsuki/internal/domain/anime/repository/sql"
	emptyIDRepository "github.com/rl404/akatsuki/internal/domain/empty_id/repository"
	emptyIDSQL "github.com/rl404/akatsuki/internal/domain/empty_id/repository/sql"
	genreRepository "github.com/rl404/akatsuki/internal/domain/genre/repository"
	genreSQL "github.com/rl404/akatsuki/internal/domain/genre/repository/sql"
	malRepository "github.com/rl404/akatsuki/internal/domain/mal/repository"
	malClient "github.com/rl404/akatsuki/internal/domain/mal/repository/client"
	publisherRepository "github.com/rl404/akatsuki/internal/domain/publisher/repository"
	publisherPubsub "github.com/rl404/akatsuki/internal/domain/publisher/repository/pubsub"
	studioRepository "github.com/rl404/akatsuki/internal/domain/studio/repository"
	studioSQL "github.com/rl404/akatsuki/internal/domain/studio/repository/sql"
	userAnimeRepository "github.com/rl404/akatsuki/internal/domain/user_anime/repository"
	userAnimeSQL "github.com/rl404/akatsuki/internal/domain/user_anime/repository/sql"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/akatsuki/pkg/pubsub"
	_nr "github.com/rl404/fairy/log/newrelic"
	nrPS "github.com/rl404/fairy/monitoring/newrelic/pubsub"
)

func consumer() error {
	// Get config.
	cfg, err := getConfig()
	if err != nil {
		return err
	}
	utils.Info("config initialized")

	// Init newrelic.
	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.Newrelic.Name),
		newrelic.ConfigLicense(cfg.Newrelic.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		utils.Error(err.Error())
	} else {
		defer nrApp.Shutdown(10 * time.Second)
		utils.AddLog(_nr.NewFromNewrelicApp(nrApp, _nr.ErrorLevel))
		utils.Info("newrelic initialized")
	}

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
	ps = nrPS.New(cfg.PubSub.Dialect, ps, nrApp)
	utils.Info("pubsub initialized")
	defer ps.Close()

	// Init anime.
	var anime animeRepository.Repository = animeSQL.New(db, cfg.Cron.FinishedAge, cfg.Cron.ReleasingAge, cfg.Cron.NotYetAge)
	utils.Info("repository anime initialized")

	// Init genre.
	var genre genreRepository.Repository = genreSQL.New(db)
	utils.Info("repository genre initialized")

	// Init studio.
	var studio studioRepository.Repository = studioSQL.New(db)
	utils.Info("repository studio initialized")

	// Init user anime.
	var userAnime userAnimeRepository.Repository = userAnimeSQL.New(db, cfg.Cron.UserAnimeAge)
	utils.Info("repository user anime initialized")

	// Init empty id.
	var emptyID emptyIDRepository.Repository = emptyIDSQL.New(db)
	utils.Info("repository empty id initialized")

	// Init mal.
	var mal malRepository.Repository = malClient.New(cfg.Mal.ClientID)
	utils.Info("repository mal initialized")

	// Init publisher.
	var publisher publisherRepository.Repository = publisherPubsub.New(ps, pubsubTopic)
	utils.Info("repository publisher initialized")

	// Init service.
	service := service.New(anime, genre, studio, userAnime, emptyID, publisher, mal)
	utils.Info("service initialized")

	// Init consumer.
	consumer := _consumer.New(service, ps, pubsubTopic)
	utils.Info("consumer initialized")
	defer consumer.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Start subscribe.
	if err := consumer.Subscribe(nrApp); err != nil {
		return err
	}

	utils.Info("consumer ready")
	<-sigChan

	return nil
}
