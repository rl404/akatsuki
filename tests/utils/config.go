package utils_test

import (
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rl404/akatsuki/internal/utils"
)

type config struct {
	Cache  cacheConfig  `envconfig:"CACHE"`
	DB     dbConfig     `envconfig:"DB"`
	PubSub pubsubConfig `envconfig:"PUBSUB"`
	Cron   cronConfig   `envconfig:"CRON"`
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
	Address  string `envconfig:"ADDRESS" validate:"required" mod:"default=amqp://guest:guest@localhost:5672"`
	Password string `envconfig:"PASSWORD"`
}

type cronConfig struct {
	UpdateLimit  int `envconfig:"UPDATE_LIMIT" validate:"required,gte=0" mod:"default=10"`
	FillLimit    int `envconfig:"FILL_LIMIT" validate:"required,gte=0" mod:"default=30"`
	ReleasingAge int `envconfig:"RELEASING_AGE" validate:"required,gt=0" mod:"default=1"`  // days
	FinishedAge  int `envconfig:"FINISHED_AGE" validate:"required,gt=0" mod:"default=30"`  // days
	NotYetAge    int `envconfig:"NOT_YET_AGE" validate:"required,gt=0" mod:"default=7"`    // days
	UserAnimeAge int `envconfig:"USER_ANIME_AGE" validate:"required,gt=0" mod:"default=7"` // days
}

const envPath = ".env"
const envPrefix = "AKATSUKI_TEST"
const pubsubTopic = "akatsuki-pubsub-test"

func getRepoPath() string {
	re := regexp.MustCompile(`^(.*akatsuki)`)
	cwd, _ := os.Getwd()
	return string(re.Find([]byte(cwd)))
}

func GetConfig() (*config, error) {
	var cfg config

	_ = godotenv.Load(getRepoPath() + "/" + envPath)

	if err := envconfig.Process(envPrefix, &cfg); err != nil {
		return nil, err
	}

	if err := utils.Validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
