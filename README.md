# Example of Go Microservice

This is an example of Go microservice and how we use `circleci`/`drone.io` to handle multiple Go service in one repository.

All critics and feedbacks is much appreciated.

| CI | Status |
| ---- | ----------- |
| [CircleCI](https://circleci.com/gh/lab46/example) | [![CircleCI](https://circleci.com/gh/lab46/example.svg?style=svg)](https://circleci.com/gh/lab46/example) |
| [drone.io](http://droneio.albertwidi.com/lab46/example) | [![drone.io Build Status](http://droneio.albertwidi.com/api/badges/lab46/example/status.svg?branch=master)](http://droneio.albertwidi.com/lab46/example) |

## Go Test & Build

Test and build is configured around `unix` system environment. If you are using Windows or other systems, there might be some part that not works for you.

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

All `go test` and `go build` command exist in `GoTest.sh`. The `bash-cript` will detect all changed files in one commit. And will only test and build affected packages.

This repo is using several way to Go test and build:
1. `make test` command will trigger `@./GoTest.sh diff` command and will run test based on `git` changed/untrack files.
2. `make test.diffmaster` will diff your current & committed branch against master and test it.
3. `make test.circleci`/`make test.droneio` is used for continous integration and the test is running by executing `./GoTest.sh ${COMMIT_HASH}` or for example: `./GoTest.sh 8e876439933c60badd1f2828655dffe2c34512c8`.

## Dependencies

Note that all dependencies and configurations made for the sake of this example project and may not suitable for your needs.

### Webserver

Standard Go webserver with several builtin endpoints:
- `metrics` endpoint to expose metrics for Prometheus
- `healthcheck` endpoint to check the health of service

Internally `webserver` use `route` package and `gorilla/mux` to automatically expose `http-metrics` and add timeout mechanism

### Configuration Loader

Configurations depends on `EXMPLENV` environment variable and all configuration file must be in `*.yaml` format.

If `EXMPLENV` not exist, the default value is `dev`.

## Bookapp and Rentapp

There are two services called bookapp and rentapp. Bookapp is an app to serve list of books and Rentapp is an app to rent a book.

### How to run the test

1. Run `make docker.compose-test.up`, this will build all go binaries needed and docker image then run docker-compose in daemon/background mode.
2. Run `make docker.compose-test.down` will take down all running container from docker-compose test.

A more old and manual way: 

1. You need to run `make build.all` first to build all docker image required.
2. Run `make docker.run.all` to run all docker container.
3. Run `make docker.stop.all` to stop all docker container and remove its dependencies.

### Service to Service Communication

When you run the service via container, applications will talk to each others via `envoy`. It is configured this way to provide better timeout, circuitbreaking, retry and fallback mechanism to the microservice environment. However a full service mesh is not yet introduced to the ecosystem.

Envoy configuration reference: https://www.envoyproxy.io/docs/envoy/v1.5.0/intro/deployment_types/service_to_service