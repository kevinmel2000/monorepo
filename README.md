# Example of Go Webservice

This is an example of Go microservice and how we use `drone.io` to handle multiple Go service in one repository

## Go Test & Build

All `go test` and `go build` command exist in `GoTest.sh`. The `bash-cript` will detect all changed files in one commit. And will only test and build affected packages.

To run the test and build: `./GoTest.sh ${COMMIT_HASH}` or for example: `./GoTest.sh 8e876439933c60badd1f2828655dffe2c34512c8`

## Dependencies

Note that all dependencies and configurations made for the sake of this example project and may not suitable for your needs.

### Webserver

Standard Go webserver with several endpoints builtin:
- `metrics` endpoint to expose metrics for Prometheus
- `healthcheck` endpoint to check the health of service

Internally `webserver` use `route` package to automatically expose `http-metrics` and timeout mechanism

### Configuration Loader

Configurations depends on `EXMPLENV` environment variable and all configuration file must be in `*.yaml` format.

If `EXMPLENV` not exist, the default value is `dev`.

## Bookapp and Rentapp

There are two services called bookapp and rentapp. Bookapp is an app to serve list of books and Rentapp is an app to rent a book.

Will provide more documentation later

## Continous Integration

CI Integration is running via drone.io/Drone CI