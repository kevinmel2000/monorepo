#!bin/sh

# get last commit SHA
LASTCOMMIT = $(shell git log -n 1 --pretty=%H)
CURRENT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

# test

test:
	make test.diff

test.diff:
	@./GoTest.sh diff

test.diffmaster:
	@./GoTest.sh ${CURRENT_BRANCH} branch

test.localcommit:
	@./GoTest.sh ${LASTCOMMIT}

test.dir: 
	@echo ">>> go test based on directory"
	@if [ -d "./$(dir)" ]; then \
		go test -v ./$(dir)...; \
	fi

# test continous integration

test.droneio:
	@./GoTest.sh ${DRONE_COMMIT_SHA}

test.circleci:
	@./GoTest.sh ${CIRCLE_SHA1}

test.circleold:
	@if [ "${CIRCLECI_RETRY}" = "" ]; then \
		echo "test CIRCLE_SHA1"; \
		./GoTest.sh ${CIRCLE_SHA1}; \
	else \
		echo "test CIRCLECI_RETRY_SHA1"; \
		./GoTest.sh ${CIRCLECI_RETRY_SHA1}; \
	fi	

# go build & run

go.build.service:
	@go build -o $(name) service/$(name)/*.go

go.build.tools:
	@go build -o $(name) tools/$(name)/*.go

go.run.service:
	make go.build.service $(name)
	@EXMPLENV=$(env) ./$(name) -log_level=debug -config_dir=service/$(name)/files/config


## docker specific

# docker.go.build function. param1: dockerimage, param2: appname
define docker.go.build
	-docker image rm -f $(1)
	@echo ">>> compiling $(2)"
	@cd files && GOOS=linux go build -o $(2) ../$(2)/*.go
	@echo ">>> docker build"
	@cd files && docker build . -f Dockerfile.$(2) -t $(1)
	@echo ">>> removing $(2) binary"
	@cd files && rm $(2)
endef

## bookapp

docker.build.bookapp:
	$(eval dockerimage = "bookapp:v0.10")
	$(eval appname = "bookapp")
	$(call docker.go.build,$(dockerimage),$(appname))

# /envoy/envoy_linux -c /envoy/envoy_config.json
docker.run.bookapp:
	$(eval dockerimage = "bookapp:v0.10")
	$(eval appname = "bookapp")
	@docker run --name $(appname) -d $(dockerimage) tail -f /dev/null


docker.stop.bookapp:
	-docker stop bookapp
	-docker rm -f bookapp

## rentapp

docker.build.rentapp:
	$(eval dockerimage = "rentapp:v0.10")
	$(eval appname = "rentapp")
	$(call docker.go.build,$(dockerimage),$(appname))

docker.run.rentapp:
	$(eval dockerimage = "rentapp:v0.10")
	$(eval appname = "rentapp")
	@docker run --name $(appname) -d $(dockerimage)

docker.stop.rentapp:
	-docker stop rentapp

## need a better solution for run and stop all

docker.compose-test.up:
	@cd files && GOOS=linux go build -v -o bookapp ../bookapp/*.go
	@cd files && GOOS=linux go build -v -o rentapp ../rentapp/*.go
	@docker-compose -f Docker-compose.test.yaml up -d
	@cd files && rm bookapp
	@cd files && rm rentapp

docker.compose-test.down:
	@docker-compose -f Docker-compose.test.yaml down

docker.build.all:
	make docker.build.bookapp
	make docker.build.rentapp

docker.run.all:
	-docker network create localnet
	@docker run --name bookapp -d --net localnet bookapp:v0.10
	@docker run --name rentapp -d --net localnet rentapp:v0.10

docker.stop.all:
	-docker stop bookapp
	-docker rm -f bookapp
	-docker stop rentapp
	-docker rm -f rentapp
	-docker network rm localnet	