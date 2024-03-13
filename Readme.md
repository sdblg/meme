
# Meme microservice

This microservice manages the APIs and resources that the Memes for 

## REST APIs

## Usage of Makefile

```bash
make

Usage:
  make <target>

Build
  build                                     build the binary
  docker                                    Build docker
    
Run
  run                                       Build as per tag and run the binary

Utility
  check                                     Check code quality before commit
  cover                                     Run the go test with coverage
  test                                      Check unit test with coverage
  init                                      Run this to install dependencies
  help                                      Display this help.
```

Run `make init check test` for installing all dependencies, checking the code and running the unit test

### Check Code

To check for secrets being committed, linting errors, and openAPI spec errors:

```bash
make check
```

Note: The openAPI spec is generated automatically from the go source code and located at `docs/openapi-spec.yaml` .

## Build and run

Build just the `meme` executable and execute it:

```bash
make run
```

Or build it as a docker image and run that:

```bash
make docker run-docker
```

Note: The default config uses the `postgres` DB driver. This means you need to start a postgres database to run `meme`. You can do this either by installing postgres manually or using the `test-start-postgres` make target to start postgres as a container.

### Publish the docker image to artifactory

```bash
make publish
```


### Environment Variables

| EnvVar Name   | Description  |  Default Value |
| ------------- | ------------- | -------------- |
| `DB_HOST`| Host required for the DB | `localhost`|
| `DB_PORT`| Port required by the DB |`5432`|
| `DB_SSLROOTCERT`| SSL ROOT certificate in case of TLS connection| ""|
| `DB_SSLMODE`| SSL mode enabled | disabled |
| `DB_TIMEZONE`| Timezone of the DB to use | `GMT` |
| `TOKEN_SIGNING_KEY`| HUB signing key used to manage JWT tokens | "" |

### Heath check
```bash
curl -sS http://localhost:8080/v1/ping
```

If successful, the HTTP response code will be `200` and the response body will be the string `pong` .

If successful, the HTTP response code will be `200` and the response body will be the build version.

Note: If you are running these APIs remotely replace `localhost` with the externally accessible IP or hostname of the machine the service is running on.
