package utils_test

import (
	"fmt"
	"strings"

	animeSQL "github.com/rl404/akatsuki/internal/domain/anime/repository/sql"
	emptyIDSQL "github.com/rl404/akatsuki/internal/domain/empty_id/repository/sql"
	genreSQL "github.com/rl404/akatsuki/internal/domain/genre/repository/sql"
	studioSQL "github.com/rl404/akatsuki/internal/domain/studio/repository/sql"
	userAnimeSQL "github.com/rl404/akatsuki/internal/domain/user_anime/repository/sql"
	"github.com/rl404/akatsuki/internal/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func GetDB(cfg *config) (*gorm.DB, error) {
	// Split host and port.
	split := strings.Split(cfg.DB.Address, ":")
	if len(split) != 2 {
		return nil, errors.ErrInvalidDBFormat
	}

	var dialector gorm.Dialector
	switch cfg.DB.Dialect {
	case "mysql":
		dialector = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local", cfg.DB.User, cfg.DB.Password, cfg.DB.Address, cfg.DB.Name))
	case "postgresql":
		dialector = postgres.Open(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", split[0], split[1], cfg.DB.User, cfg.DB.Password, cfg.DB.Name))
	default:
		return nil, errors.ErrOneOfField("dialect", "mysql postgresql")
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}

	if err := TruncateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
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
	)
}

func TruncateDB(db *gorm.DB) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if err := tx.Unscoped().Delete(&animeSQL.Anime{}).Error; err != nil {
		return err
	}

	if err := tx.Unscoped().Delete(&animeSQL.AnimeGenre{}).Error; err != nil {
		return err
	}

	if err := tx.Unscoped().Delete(&animeSQL.AnimePicture{}).Error; err != nil {
		return err
	}

	if err := tx.Unscoped().Delete(&animeSQL.AnimeRelated{}).Error; err != nil {
		return err
	}

	if err := tx.Unscoped().Delete(&animeSQL.AnimeStudio{}).Error; err != nil {
		return err
	}

	if err := tx.Unscoped().Delete(&animeSQL.AnimeStatsHistory{}).Error; err != nil {
		return err
	}

	if err := tx.Unscoped().Delete(&genreSQL.Genre{}).Error; err != nil {
		return err
	}

	if err := tx.Unscoped().Delete(&studioSQL.Studio{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}
