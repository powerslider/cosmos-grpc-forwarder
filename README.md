# Cosmos gRPC Forwarder

## Overview

## Getting Started

### Local Setup

**Step 0.** Install [pre-commit](https://pre-commit.com/):

```shell
pip install pre-commit

# For macOS users.
brew install pre-commit
```

Then run `pre-commit install` to setup git hook scripts.
Used hooks can be found [here](.pre-commit-config.yaml).

______________________________________________________________________

NOTE

> `pre-commit` aids in running checks (end of file fixing,
> markdown linting, go linting, runs go tests, json validation, etc.)
> before you perform your git commits.

______________________________________________________________________

**Step 1.** Install external tooling (golangci-lint, etc.):

```shell script
make install
```

**Step 2.** Setup project for local testing (code lint, runs tests, builds all needed binaries):

```shell script
make all
```

______________________________________________________________________

NOTE

> All binaries can be found in `<project_root>/bin` directory.
> Use `make clean` to delete old binaries.

______________________________________________________________________

**Step 3.** Run server:

```shell
make run-server
```

______________________________________________________________________

NOTE

> Check [Makefile](Makefile) for other useful commands.

______________________________________________________________________

### Docker-compose Setup

**Step 1.** Run `docker-compose` to build and run the application as a Docker container:

```shell script
docker-compose up -d
```

## License

[MIT](LICENSE)
