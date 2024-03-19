# imdb-watchlist
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/clambin/imdb-watchlist?color=green&label=Release&style=plastic)
![Build)](https://github.com/clambin/imdb-watchlist/workflows/Build/badge.svg)
![Codecov](https://img.shields.io/codecov/c/gh/clambin/imdb-watchlist?style=plastic)
![Go Report Card](https://goreportcard.com/badge/github.com/clambin/imdb-watchlist)
![GitHub](https://img.shields.io/github/license/clambin/imdb-watchlist?style=plastic)

Import IMDB Watchlist in Sonarr

## Installation

Binaries are available on the [release](https://github.com/clambin/imdb-watchlist/releases) page. Docker images are available on [ghcr.io](https://github.com/clambin/imdb-watchlist/pkgs/container/imdb-watchlist).


## Running
### Command-line options

```
Usage: imdb-watchlist --list=LIST [<flags>]

  -addr string
        Server address (default ":8080")
  -apikey string
        APIKey
  -debug
        Log debug messages
  -list string
        IMDB List ID(s) (required, comma-separated)
  -prometheus string
        Prometheus metrics address (default ":9090")

```

If '-apikey' is not specified, imdb-watchlist generates & logs a key automatically.

### Configuring Sonarr

Add an Import List in Sonarr/Radarr, specifying the API key described above.

## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
