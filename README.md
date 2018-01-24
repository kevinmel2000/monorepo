# Example of Go Microservice

This is an example and experiment of Go microservice and how we use `circleci`/`drone.io` to handle multiple Go service in single repository.

The case is there are more than 30 engineers that is working in a big team/tribe and this type of repository and we have various products/services to do and maintain.

Open to discussions and all critics + feedbacks is much appreciated.

| CI | Status |
| ---- | ----------- |
| [CircleCI](https://circleci.com/gh/lab46/example) | [![CircleCI](https://circleci.com/gh/lab46/example.svg?style=svg)](https://circleci.com/gh/lab46/example) |
| [drone.io](http://droneio.albertwidi.com/lab46/example) | [![drone.io Build Status](http://droneio.albertwidi.com/api/badges/lab46/example/status.svg?branch=master)](http://droneio.albertwidi.com/lab46/example) |

## Go Test & Build

Test and build is configured around `unix` system environment. If you are using Windows or other systems, there might be some parts that not works for you.

All test is structured and being build around `git` version control environment.

1. Make sure that you have Go in your system. Check the `go` command by type `go` or `go version` in the Terminal:
```shell
❯ go version
go version go1.9.2 darwin/amd64
```
2. Make sure you have git in your system. Check by type `git` or `git version` in the Terminal:
```shell
❯ git version
git version 2.11.0 (Apple Git-81)   
```

3. Make sure you have docker in your system. Check by type `docker` or `docker version` in the Terminal:
```shell
❯ docker version
Client:
 Version:      17.09.0-ce
 API version:  1.32
 Go version:   go1.8.3
 Git commit:   afdb6d4
 Built:        Tue Sep 26 22:40:09 2017
 OS/Arch:      darwin/amd64

Server:
 Version:      17.09.0-ce
 API version:  1.32 (minimum version 1.12)
 Go version:   go1.8.3
 Git commit:   afdb6d4
 Built:        Tue Sep 26 22:45:38 2017
 OS/Arch:      linux/amd64
 Experimental: true
```

4. Start the dependencies first by using `docker-compose` command: `docker-compose up -d`.
```shell
❯ docker-compose up -d
Creating network "example_default" with the default driver
Creating example_redis_1 ...
Creating example_postgres_1 ...
Creating example_redis_1
Creating example_redis_1 ... done
```

All `go test` and `go build` command exist in `GoTest.sh`. The `bash-script` will detect all changed files in one commit. And will only test and build affected packages.

This repo is using several way to Go test and build:
1. `make test` command will trigger `@./GoTest.sh diff` command and will run test based on `git` changed/untrack files.
2. `make test.diffmaster` will diff your current & committed branch against master and test it.
3. `make test.circleci`/`make test.droneio` is used for continous integration and the test is running by executing `./GoTest.sh ${COMMIT_HASH}` or for example: `./GoTest.sh 8e876439933c60badd1f2828655dffe2c34512c8`.

### Unit Test & Integration Test

This project use docker and CI to provide database instance(postgresql). Instead mock the database conection, we are test the schema directly to database by using `sqlimporter`

## Dependencies

Note that all dependencies and configurations made for the sake of this example project and may not suitable for your needs.

### Webserver

Standard Go webserver with several builtin endpoints:
- `metrics` endpoint to expose metrics for Prometheus
- `healthcheck` endpoint to check the health of service
- `status` endpoint to expose current service env and configuration

Internally `webserver` use `route` package and `gorilla/mux` to automatically expose `http-metrics` and add timeout mechanism

### Configuration Loader

Configurations depends on `EXMPLENV` environment variable and all configuration file must be in `*.yaml` format.

If `EXMPLENV` not exist, the default value is `dev`.

### Testutil - Sqlimporter 

Sqlimporter is used to create a random database or schema and import *.sql files schema.

The database/schema can be dropped directly after being used for a test.

## Bookapp and Rentapp

There are two services called bookapp and rentapp. Bookapp is an app to serve list of books and Rentapp is an app to rent a book.

### How to run the service in test mode

1. Run `make docker.compose-test.up`, this will build all go binaries needed and docker image then run docker-compose in daemon/background mode.
2. Run `make docker.compose-test.down` will take down all running container from docker-compose test.

A more old and manual way: 

1. You need to run `make build.all` first to build all docker image required.
2. Run `make docker.run.all` to run all docker container.
3. Run `make docker.stop.all` to stop all docker container and remove its dependencies.

### Service to Service Communication

When you run the service via container, applications will talk to each others via `envoy`. It is configured this way to provide better timeout, circuitbreaking, retry and fallback mechanism to the microservice environment. However a full service mesh is not yet introduced to the ecosystem.

Envoy configuration reference: https://www.envoyproxy.io/docs/envoy/v1.5.0/intro/deployment_types/service_to_service