package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	grpcAPI "github.com/rl404/akatsuki/internal/delivery/grpc/api"
	"github.com/rl404/akatsuki/internal/delivery/grpc/schema"
	httpAPI "github.com/rl404/akatsuki/internal/delivery/rest/api"
	"github.com/rl404/akatsuki/internal/delivery/rest/ping"
	"github.com/rl404/akatsuki/internal/delivery/rest/swagger"
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
	"github.com/rl404/akatsuki/pkg/grpc"
	"github.com/rl404/akatsuki/pkg/http"
	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/log"
	promCache "github.com/rl404/fairy/monitoring/prometheus/cache"
	promDB "github.com/rl404/fairy/monitoring/prometheus/database"
	promMW "github.com/rl404/fairy/monitoring/prometheus/middleware"
	"github.com/rl404/fairy/pubsub"
	_grpc "google.golang.org/grpc"
)

func server() error {
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
	c = promCache.New(cfg.Cache.Dialect, c)
	utils.Info("cache initialized")
	defer c.Close()

	// Init in-memory.
	im, err := cache.New(cache.InMemory, "", "", 5*time.Second)
	if err != nil {
		return err
	}
	im = promCache.New("inmemory", im)
	utils.Info("in-memory initialized")
	defer im.Close()

	// Init db.
	db, err := newDB(cfg.DB)
	if err != nil {
		return err
	}
	promDB.RegisterGORM(cfg.DB.Name, db)
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
	anime = animeSQL.New(db, cfg.Cron.FinishedAge, cfg.Cron.ReleasingAge, cfg.Cron.NotYetAge)
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

	// Init web server.
	httpServer := http.New(http.Config{
		Port:            cfg.HTTP.Port,
		ReadTimeout:     cfg.HTTP.ReadTimeout,
		WriteTimeout:    cfg.HTTP.WriteTimeout,
		GracefulTimeout: cfg.HTTP.GracefulTimeout,
	})
	utils.Info("http server initialized")

	r := httpServer.Router()
	r.Use(middleware.RealIP)
	r.Use(utils.Recoverer)
	utils.Info("http server middleware initialized")

	// Register ping route.
	ping.New().Register(r)
	utils.Info("http route ping initialized")

	// Register swagger route.
	swagger.New().Register(r)
	utils.Info("http route swagger initialized")

	// Register api route.
	httpAPI.New(service).Register(r)
	utils.Info("http route api initialized")

	// Run web server.
	httpServerChan := httpServer.Run()
	utils.Info("http server listening at :%s", cfg.HTTP.Port)

	// Init GRPC.
	grpcAPI := grpcAPI.New(service)
	grpcServer := grpc.New(grpc.Config{
		Port:    cfg.GRPC.Port,
		Timeout: cfg.GRPC.Timeout,
		UnaryInterceptors: []_grpc.UnaryServerInterceptor{
			utils.RecovererGRPC,
			log.UnaryMiddlewareWithLog(utils.GetLogger()),
			promMW.NewUnaryGRPC,
		},
	})
	utils.Info("grpc server initialized")

	// Register api route.
	schema.RegisterAPIServer(grpcServer.Server(), grpcAPI)
	utils.Info("grpc route api initialized")

	// Run grpc server.
	grpcServerChan := grpcServer.Run()
	utils.Info("grpc server listening at :%s", cfg.GRPC.Port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case err := <-httpServerChan:
		if err != nil {
			return err
		}
	case err := <-grpcServerChan:
		if err != nil {
			return err
		}
	case <-sigChan:
	}

	return nil
}
