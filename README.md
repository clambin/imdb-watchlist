# imdb-watchlist
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/clambin/imdb-watchlist?color=green&label=Release&style=plastic)
![Build)](https://github.com/clambin/imdb-watchlist/workflows/Build/badge.svg)
![Codecov](https://img.shields.io/codecov/c/gh/clambin/imdb-watchlist?style=plastic)
![Go Report Card](https://goreportcard.com/badge/github.com/clambin/imdb-watchlist)
![GitHub](https://img.shields.io/github/license/clambin/imdb-watchlist?style=plastic)

Import IMDB Watchlist in Sonarr

## Installation

Binaries are available on the [release](https://github.com/clambin/imdb-watchlist/releases) page. Docker images are available on [docker hub](https://hub.docker.com/r/clambin/imdb-watchlist).


## Running
### Command-line options

```
usage: imdb-watchlist --list=LIST [<flags>]

imdb-watchlist

Flags:
  -h, --help           Show context-sensitive help (also try --help-long and --help-man).
  -v, --version        Show application version.
      --debug          Log debug messages
      --port=8080      API listener port
      --list=LIST      IMDB Watchlist ID
      --apikey=APIKEY  API Key

```

If '--apikey' is not specified, imdb-watchlist will generate one automatically and list in the logfile.

### Configuring Sonarr

Add an Import List in Sonarr, specifying the API KEY described above.

## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.