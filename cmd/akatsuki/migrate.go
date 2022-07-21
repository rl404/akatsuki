package main

import (
	animeSQL "github.com/rl404/akatsuki/internal/domain/anime/repository/sql"
	emptyIDSQL "github.com/rl404/akatsuki/internal/domain/empty_id/repository/sql"
	genreSQL "github.com/rl404/akatsuki/internal/domain/genre/repository/sql"
	studioSQL "github.com/rl404/akatsuki/internal/domain/studio/repository/sql"
	userAnimeSQL "github.com/rl404/akatsuki/internal/domain/user_anime/repository/sql"
	"github.com/rl404/akatsuki/internal/utils"
)

func migrate() error {
	// Get config.
	cfg, err := getConfig()
	if err != nil {
		return err
	}
	utils.Info("config initialized")

	// Init db.
	db, err := newDB(cfg.DB)
	if err != nil {
		return err
	}
	utils.Info("database initialized")
	tmp, _ := db.DB()
	defer tmp.Close()

	// Migrate.
	utils.Info("migrating...")
	if err := db.AutoMigrate(
		animeSQL.Anime{},
		animeSQL.AnimeGenre{},
		animeSQL.AnimePicture{},
		animeSQL.AnimeRelated{},
		animeSQL.AnimeStudio{},
		animeSQL.AnimeStatsHistory{},
		genreSQL.Genre{},
		studioSQL.Studio{},
		userAnimeSQL.UserAnime{},
		emptyIDSQL.EmptyID{},
	); err != nil {
		return err
	}

	utils.Info("done")
	return nil
}
