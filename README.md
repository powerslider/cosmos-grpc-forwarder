# Cosmos SDK gRPC Forwarder

## Overview

This projects aims to serve as a forwarder for Cosmos SDK enabled gRPC endpoints.

## Getting Started

**Step 1 (OPTIONAL).** Install external tooling (golangci-lint, etc.):

```shell script
make install
```

**Step 2 (OPTIONAL).** Setup project for local testing (code lint, runs tests, builds all needed binaries):

```shell script
make all
```

**Step 3.** Run gRPC server:

```shell script
make run-server
```

**Step 4.** Run gRPC client (in a separate terminal session):

```shell script
make run-client
```

______________________________________________________________________

NOTE

> Instead of having a file called `test_grpc.go` there is a file called
> `service_handler_test.go` where the comparison between the local server
> and the public Cosmos SDK gRPC endpoint. This is because having it this
> way gets it picked up by `go test` command.

______________________________________________________________________

**Step 5.** Change gRPC endpoint:

- At `<project-root>/.env.dist` and `<projec-root>/.env.test.dist` you can find the following environment variables:

```shell script
SERVER_NAME=cosmos-grpc-forwarder
SERVER_HOST=localhost
SERVER_PORT=8080
LOG_LEVEL=debug
LOG_FORMAT=json
COSMOS_SDK_GRPC_ENDPOINT=grpc.osmosis.zone:9090
```

- `COSMOS_SDK_GRPC_ENDPOINT` could be easily swapped for another Cosmos SDK enabled mainnet or testnet endpoint and
  it should work just fine.
- After modification run the tests again with `make test` to verify compatibility.

## Development Setup

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

**Step 3.** Other commands:

```shell script
# Re-generate proto stubs.
make proto-gen

# Lint proto files.
makke proto-lint
```

______________________________________________________________________

NOTE

> All binaries can be found in `<project_root>/bin` directory.
> Use `make clean` to delete old binaries.

______________________________________________________________________

______________________________________________________________________

NOTE

> Check [Makefile](Makefile) for other useful commands.

______________________________________________________________________

## License

[MIT](LICENSE)
