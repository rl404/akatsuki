package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/akatsuki/pkg/cache"
	"github.com/rl404/akatsuki/pkg/pubsub"
	"github.com/rl404/fairy/monitoring/newrelic/database"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type config struct {
	App      appConfig      `envconfig:"APP"`
	HTTP     httpConfig     `envconfig:"HTTP"`
	GRPC     grpcConfig     `envconfig:"GRPC"`
	Cache    cacheConfig    `envconfig:"CACHE"`
	DB       dbConfig       `envconfig:"DB"`
	PubSub   pubsubConfig   `envconfig:"PUBSUB"`
	Mal      malConfig      `envconfig:"MAL"`
	Cron     cronConfig     `envconfig:"CRON"`
	Log      logConfig      `envconfig:"LOG"`
	Newrelic newrelicConfig `envconfig:"NEWRELIC"`
}

type appConfig struct {
	Env string `envconfig:"ENV" validate:"required,oneof=dev prod" mod:"default=dev,no_space,lcase"`
}

type httpConfig struct {
	Port            string        `envconfig:"PORT" validate:"required" mod:"default=45001,no_space"`
	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" validate:"required,gt=0" mod:"default=5s"`
	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" validate:"required,gt=0" mod:"default=5s"`
	GracefulTimeout time.Duration `envconfig:"GRACEFUL_TIMEOUT" validate:"required,gt=0" mod:"default=10s"`
}

type grpcConfig struct {
	Port    string        `envconfig:"PORT" validate:"required" mod:"default=46001,no_space"`
	Timeout time.Duration `envconfig:"TIMEOUT" validate:"required,gt=0" mod:"default=10s"`
}

type cacheConfig struct {
	Dialect  string        `envconfig:"DIALECT" validate:"required,oneof=nocache redis inmemory" mod:"default=inmemory,no_space,lcase"`
	Address  string        `envconfig:"ADDRESS"`
	Password string        `envconfig:"PASSWORD"`
	Time     time.Duration `envconfig:"TIME" default:"24h" validate:"required,gt=0"`
}

type dbConfig struct {
	Dialect         string        `envconfig:"DIALECT" validate:"required,oneof=mysql postgresql" mod:"default=mysql,no_space,lcase"`
	Address         string        `envconfig:"ADDRESS" validate:"required" mod:"default=localhost:3306,no_space"`
	Name            string        `envconfig:"NAME" validate:"required" mod:"default=akatsuki"`
	User            string        `envconfig:"USER" validate:"required" mod:"default=root"`
	Password        string        `envconfig:"PASSWORD"`
	MaxConnOpen     int           `envconfig:"MAX_CONN_OPEN" validate:"required,gt=0" mod:"default=10"`
	MaxConnIdle     int           `envconfig:"MAX_CONN_IDLE" validate:"required,gt=0" mod:"default=10"`
	MaxConnLifetime time.Duration `envconfig:"MAX_CONN_LIFETIME" validate:"required,gt=0" mod:"default=1m"`
}

type pubsubConfig struct {
	Dialect  string `envconfig:"DIALECT" validate:"required,oneof=rabbitmq redis google" mod:"default=rabbitmq,no_space,lcase"`
	Address  string `envconfig:"ADDRESS" validate:"required"`
	Password string `envconfig:"PASSWORD"`
}

type malConfig struct {
	ClientID string `envconfig:"CLIENT_ID" validate:"required" mod:"no_space"`
}

type cronConfig struct {
	UpdateLimit  int `envconfig:"UPDATE_LIMIT" validate:"required,gte=0" mod:"default=10"`
	FillLimit    int `envconfig:"FILL_LIMIT" validate:"required,gte=0" mod:"default=30"`
	ReleasingAge int `envconfig:"RELEASING_AGE" validate:"required,gt=0" mod:"default=1"`  // days
	FinishedAge  int `envconfig:"FINISHED_AGE" validate:"required,gt=0" mod:"default=30"`  // days
	NotYetAge    int `envconfig:"NOT_YET_AGE" validate:"required,gt=0" mod:"default=7"`    // days
	UserAnimeAge int `envconfig:"USER_ANIME_AGE" validate:"required,gt=0" mod:"default=7"` // days
}

type logConfig struct {
	Level utils.LogLevel `envconfig:"LEVEL" default:"-1"`
	JSON  bool           `envconfig:"JSON" default:"false"`
	Color bool           `envconfig:"COLOR" default:"true"`
}

type newrelicConfig struct {
	Name       string `envconfig:"NAME" default:"akatsuki"`
	LicenseKey string `envconfig:"LICENSE_KEY"`
}

const envPath = "../../.env"
const envPrefix = "AKATSUKI"
const pubsubTopic = "akatsuki-pubsub"

var cacheType = map[string]cache.CacheType{
	"nocache":  cache.NOP,
	"redis":    cache.Redis,
	"inmemory": cache.InMemory,
}

var pubsubType = map[string]pubsub.PubsubType{
	"rabbitmq": pubsub.RabbitMQ,
	"redis":    pubsub.Redis,
	"google":   pubsub.Google,
}

func getConfig() (*config, error) {
	var cfg config

	// Load .env file.
	_ = godotenv.Load(envPath)

	// Convert env to struct.
	if err := envconfig.Process(envPrefix, &cfg); err != nil {
		return nil, err
	}

	// Override PORT env.
	if port := os.Getenv("PORT"); port != "" {
		cfg.HTTP.Port = port
	}

	// Handle google pubsub credential.
	if cfg.PubSub.Dialect == "google" {
		credFilename, err := generateGoogleServiceAccountJSON("gcp-service-account.json", cfg.PubSub.Password)
		if err != nil {
			return nil, err
		}
		cfg.PubSub.Password = credFilename
	}

	// Validate.
	if err := utils.Validate(&cfg); err != nil {
		return nil, err
	}

	// Init global log.
	if err := utils.InitLog(cfg.Log.Level, cfg.Log.JSON, cfg.Log.Color); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func newDB(cfg dbConfig) (*gorm.DB, error) {
	// Split host and port.
	split := strings.Split(cfg.Address, ":")
	if len(split) != 2 {
		return nil, errors.ErrInvalidDBFormat
	}

	var dialector gorm.Dialector
	switch cfg.Dialect {
	case "mysql":
		dialector = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Address, cfg.Name))
	case "postgresql":
		dialector = postgres.Open(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", split[0], split[1], cfg.User, cfg.Password, cfg.Name))
	default:
		return nil, errors.ErrOneOfField("dialect", "mysql postgresql sqlite sqlserver clickhouse")
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	tmp, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set basic config.
	tmp.SetMaxIdleConns(cfg.MaxConnIdle)
	tmp.SetMaxOpenConns(cfg.MaxConnOpen)
	tmp.SetConnMaxLifetime(time.Duration(cfg.MaxConnLifetime) * time.Second)

	db.Use(database.NewGORM(cfg.Address, cfg.Name))

	return db, nil
}

func generateGoogleServiceAccountJSON(filename, value string) (string, error) {
	if err := os.WriteFile(filename, []byte(value), 0644); err != nil {
		return "", err
	}
	return filename, nil
}
