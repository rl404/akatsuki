# Nagato

[![Go Report Card](https://goreportcard.com/badge/github.com/rl404/nagato)](https://goreportcard.com/report/github.com/rl404/nagato)
![License: MIT](https://img.shields.io/github/license/rl404/nagato.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/rl404/nagato.svg)](https://pkg.go.dev/github.com/rl404/nagato)

**_nagato_** is just another API wrapper library to use [MyAnimeList](https://myanimelist.net/) API.

This library contains 2 types of wrapper library, **_nagato_** and **_mal_**.
**_mal_** is library to call MyAnimeList API without any validation and modification of the request and response. **_nagato_** wraps **_mal_** but will validate the request first before calling the API and modify the response to an easier to use format.

This libary is built following MyAnimeList's [API reference](https://myanimelist.net/apiconfig/references/api/v2) and inspired by [nstratos](https://github.com/nstratos/go-myanimelist)'s library. This library also contains some [undocumented documentation](https://myanimelist.net/forum/?topicid=2006357), so use it at your own risk.

## Features

| Feature | **_mal_** | **_nagato_** |
| --- | :---: | :---: |
| Oauth2 | :heavy_check_mark: | :heavy_check_mark: |
| Get anime list | :heavy_check_mark: | :heavy_check_mark: |
| Get anime details | :heavy_check_mark: | :heavy_check_mark: |
| Get anime ranking | :heavy_check_mark: | :heavy_check_mark: |
| Get seasonal anime | :heavy_check_mark: | :heavy_check_mark: |
| Get suggested anime | :heavy_check_mark: | :heavy_check_mark: |
| Update my anime list | :heavy_check_mark: | :hourglass: |
| Delete my anime list | :heavy_check_mark: | :hourglass: |
| Get user anime list | :heavy_check_mark: | :heavy_check_mark: |
| Get forum boards | :heavy_check_mark: | :hourglass: |
| Get forum topic detail |  :heavy_check_mark: | :hourglass: |
| Get forum topics | :heavy_check_mark: | :hourglass: |
| Get manga list | :heavy_check_mark: | :heavy_check_mark: |
| Get manga details | :heavy_check_mark: | :heavy_check_mark: |
| Get manga ranking | :heavy_check_mark: | :heavy_check_mark: |
| Update my manga list | :heavy_check_mark: | :hourglass: |
| Delete my manga list | :heavy_check_mark: | :hourglass: |
| Get user manga list | :heavy_check_mark: | :heavy_check_mark: |
| Get my user info | :heavy_check_mark: | :hourglass: |
| **Additional Features** |
| Return response code | :heavy_check_mark: | :heavy_check_mark: |
| Rate limit | :heavy_check_mark: | :heavy_check_mark: |
| Validate request | :x: | :heavy_check_mark: |
| Custom data type | :x: | :heavy_check_mark: |
| Prettier returned struct | :x: | :heavy_check_mark: |

## Installation

```sh
# For nagato.
go get github.com/rl404/nagato

# For mal.
go get github.com/rl404/nagato/mal
```

## Requirement

You need **client id** and **client secret** to use the API. You can get them by
registering from [here](https://myanimelist.net/apiconfig).

## Example

Please go to `/example` folder for examples. Also, you can read the [documentation](https://pkg.go.dev/github.com/rl404/nagato) for more details.

## Trivia

[Nagato](https://en.wikipedia.org/wiki/Japanese_battleship_Nagato)'s name is taken from japanese battleship. She is the lead ship of the Nagato-class. Also, [exists](https://en.kancollewiki.net/Nagato) in Kantai Collection games and anime.

## License

MIT License

Copyright (c) 2022 Axel
