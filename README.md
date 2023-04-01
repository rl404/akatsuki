# Akatsuki

[![Go Report Card](https://goreportcard.com/badge/github.com/rl404/akatsuki)](https://goreportcard.com/report/github.com/rl404/akatsuki)
![License: MIT](https://img.shields.io/github/license/rl404/akatsuki)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/rl404/akatsuki)
[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/rl404/akatsuki)](https://hub.docker.com/r/rl404/akatsuki)
[![publish & deploy](https://github.com/rl404/akatsuki/actions/workflows/publish-deploy.yml/badge.svg)](https://github.com/rl404/akatsuki/actions/workflows/publish-deploy.yml)

Akatsuki is [MyAnimeList](https://myanimelist.net/) anime database dump and REST API.

Powered by my [nagato](https://github.com/rl404/nagato) library and [MyAnimeList API](https://myanimelist.net/apiconfig/references/api/v2) as reference.

## Features

- Save anime details
    - Anime data
    - Anime genres
    - Anime pictures
    - Anime relation (with other anime)
    - Anime studios
- Save anime stats history
- Save user anime list
- Get all anime related in user anime list
- Handle empty anime id
- Auto update anime & user data (cron)
- Interchangeable database
    - [MySQL](https://www.mysql.com/)
    - [PostgreSQL](https://www.postgresql.org/)
    - [SQLite](https://www.sqlite.org/)
    - [SQL server](https://www.microsoft.com/en-us/sql-server/sql-server-downloads)
    - [ClickHouse](https://clickhouse.com/)
- Interchangeable cache
    - no cache
    - inmemory
    - [Redis](https://redis.io/)
    - [Memcache](https://memcached.org/)
- Interchangeable pubsub
    - [NSQ](https://nsq.io/)
    - [RabbitMQ](https://www.rabbitmq.com/)
    - [Redis](https://redis.io/)
    - [Google PubSub](https://cloud.google.com/pubsub)
- [Swagger](https://github.com/swaggo/swag)
- [Docker](https://www.docker.com/)
- [Newrelic](https://newrelic.com/) monitoring
    - HTTP
    - Cron
    - Database
    - Cache
    - Pubsub
    - External API

*More will be coming soon...*

## Requirement

- [Go](https://go.dev/)
- [MyAnimeList](https://myanimelist.net/) [client id](https://myanimelist.net/apiconfig)
- Database ([MySQL](https://www.mysql.com/)/[PostgreSQL](https://www.postgresql.org/)/[SQLite](https://www.sqlite.org/)/[SQL server](https://www.microsoft.com/en-us/sql-server/sql-server-downloads)/[ClickHouse](https://clickhouse.com/))
- PubSub ([NSQ](https://nsq.io/)/[RabbitMQ](https://www.rabbitmq.com/)/[Redis](https://redis.io/)/[Google PubSub](https://cloud.google.com/pubsub))
- (optional) Cache ([Redis](https://redis.io/)/[Memcache](https://memcached.org/))
- (optional) [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- (optional) [Newrelic](https://newrelic.com/) license key

## Installation

### Without [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)

1. Clone the repository.
```sh
git clone github.com/rl404/akatsuki
```
2. Rename `.env.sample` to `.env` and modify the values according to your setup.
3. Create the database according to your `.env`.
4. Migrate the tables.
```sh
make migrate
```
5. Run. You need at least 2 consoles/terminals.
```sh
# Run the API.
make

# Run the consumer.
make consumer
```
6. [localhost:45001](http://localhost:45001) is ready (port may varies depend on your `.env`).

#### Other commands

```sh
# Update old anime data.
make cron-update

# Fill missing anime data.
make cron-fill
```

### With [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)

1. Clone the repository.
```sh
git clone github.com/rl404/akatsuki
```
2. Rename `.env.sample` to `.env` and modify the values according to your setup.
3. Create the database according to your `.env`.
4. Get docker image.
```sh
# Pull existing image.
docker pull rl404/akatsuki

# Or build your own.
make docker-build
```
5. Migrate the tables.
```sh
make docker-migrate
```

6. Run the container. You need at least 2 consoles/terminals.
```sh
# Run the API.
make docker-api

# Run the consumer.
make docker-consumer
```
7. [localhost:45001](http://localhost:45001) is ready (port may varies depend on your `.env`).

#### Other commands

```sh
# Update old anime data.
make docker-cron-update

# Fill missing anime data.
make docker-cron-fill

# Stop running containers.
make docker-stop
```

## Environment Variables

Env | Default | Description
--- | :---: | ---
`AKATSUKI_APP_ENV` | `dev` | Environment type (`dev`/`prod`).
`AKATSUKI_HTTP_PORT` | `45001` | HTTP server port.
`AKATSUKI_HTTP_READ_TIMEOUT` | `5s` | HTTP read timeout.
`AKATSUKI_HTTP_WRITE_TIMEOUT` | `5s` | HTTP write timeout.
`AKATSUKI_HTTP_GRACEFUL_TIMEOUT` | `10s` | HTTP gracefull timeout.
`AKATSUKI_GRPC_PORT` | `46001` | GRPC server port.
`AKATSUKI_GRPC_TIMEOUT` | `10s` | GRPC timeout.
`AKATSUKI_CACHE_DIALECT` | `inmemory` | Cache type (`nocache`/`redis`/`inmemory`/`memcache`)
`AKATSUKI_CACHE_ADDRESS` | | Cache address.
`AKATSUKI_CACHE_PASSWORD` | | Cache password.
`AKATSUKI_CACHE_TIME` | `24h` | Cache time.
`AKATSUKI_DB_DIALECT` | `mysql` | Database type (`mysql`/`postgresql`/`sqlite`/`sqlserver`/`clickhouse`)
`AKATSUKI_DB_ADDRESS` | `localhost:3306` | Database address with port.
`AKATSUKI_DB_NAME` | `akatsuki` | Database name.
`AKATSUKI_DB_USER` | | Database username.
`AKATSUKI_DB_PASSWORD` | | Database password.
`AKATSUKI_DB_MAX_CONN_OPEN` | `10` | Max open database connection.
`AKATSUKI_DB_MAX_CONN_IDLE` | `10` | Max idle database connection.
`AKATSUKI_DB_MAX_CONN_LIFETIME` | `1m` | Max database connection lifetime.
`AKATSUKI_PUBSUB_DIALECT` | `rabbitmq` | Pubsub type (`nsq`/`rabbitmq`/`redis`/`google`)
`AKATSUKI_PUBSUB_ADDRESS` | | Pubsub address (if you are using `google`, this will be your google project id).
`AKATSUKI_PUBSUB_PASSWORD` | | Pubsub password (if you are using `google`, this will be the content of your google service account json).
`AKATSUKI_MAL_CLIENT_ID` | | MyAnimeList client id.
`AKATSUKI_CRON_UPDATE_LIMIT` | `10` | Anime count limit when updating old data.
`AKATSUKI_CRON_FILL_LIMIT` | `30` | Anime count limit when filling missing anime data.
`AKATSUKI_CRON_RELEASING_AGE` | `1` | Age of old releasing/airing anime data (in days).
`AKATSUKI_CRON_FINISHED_AGE` | `30` | Age of old finished anime data (in days).
`AKATSUKI_CRON_NOT_YET_AGE` | `7` | Age of old not yet released/aired anime (in days).
`AKATSUKI_CRON_USER_ANIME_AGE` | `7` | Age of old user anime list (in days).
`AKATSUKI_NEWRELIC_NAME` | `akatsuki` | Newrelic application name.
`AKATSUKI_NEWRELIC_LICENSE_KEY` | | Newrelic license key.

## Trivia

[Akatsuki](https://en.wikipedia.org/wiki/Japanese_destroyer_Akatsuki_(1932))'s name is taken from japanese destroyer with her sisters (Inazuma, Hibiki, Ikazuchi). Also, [exists](https://en.kancollewiki.net/Akatsuki) in Kantai Collection games and anime.

## Disclaimer

Akatsuki is meant for educational purpose and personal usage only. Please use it responsibly according to MyAnimeList [API License and Developer Agreement](https://myanimelist.net/static/apiagreement.html).

All data belong to their respective copyrights owners, akatsuki does not have any affiliation with content providers.

## License

MIT License

Copyright (c) 2022 Axel
