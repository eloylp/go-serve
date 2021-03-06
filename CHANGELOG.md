
## [v2.0.0](https://github.com/eloylp/go-serve/releases/v2.0.0) - 2021-06-17
Check new [README.md](README.md) for more info.
  
## [v1.3.1](https://github.com/eloylp/go-serve/releases/v1.3.1) - 2020-06-27
### Changed
- Go version to 1.14.4
- Fix HTTP shutdown with shutdown wrapper. Was not waiting for ending
connections.

## [v1.3.0](https://github.com/eloylp/go-serve/releases/v1.3.0) - 2020-04-05
### Changed
- Golang version
- An updated README.md
### Added
- Support for basic auth file.
- Now server stops gracefully after a SIGTERM or SIGINT signals.
### Changed
- Code readability improved 
- Added more tests
- Added middleware tooling

## [v1.2.0](https://github.com/eloylp/go-serve/releases/v1.2.0) - 2019-09-02
### Added
- Support for request logging
- Support for a prefix when serving files
### Changed
- An improved README.md
- Improved CHANGELOG.md
### Fixed
- An outdated link to the latest binary in the README.md

## [v1.1.0](https://github.com/eloylp/go-serve/releases/v1.1.0) - 2019-09-01
### Added
- Better binary version information both in http header and CLI advice

## [v1.0.0](https://github.com/eloylp/go-serve/releases/v1.0.0) - 2019-08-31
### Added
- The main HTTP server
- .goreleaser file for binary releases
- The readme file
- This changelog file