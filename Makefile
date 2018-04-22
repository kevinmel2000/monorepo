#!bin/sh

# get last commit SHA
LASTCOMMIT = $(shell git log -n 1 --pretty=%H)
CURRENT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

# repo

repo.setup:
	@go get -v -u github.com/golang/dep/cmd/dep
	@dep ensure -v --vendor-only
	make tools.update

# test

test:
	make test.diff

test.diff:
	@set -e
	@./GoTest.sh diff
	@git-test diff

test.commit:
	@./GoTest.sh $(commit)
	@git-test commit $(commit)

test.env:
	@echo "Build Parameters:"
	@echo "MONOREPO_SERVICE=$$MONOREPO_SERVICE"

# test.diffmaster:
# 	@./GoTest.sh ${CURRENT_BRANCH} branch

test.commit.local:
	make test commit=${LASTCOMMIT}

test.dir: 
	@echo ">>> go test based on directory"
	@if [[ -d "./$(dir)" ]]; then \
		go test -v ./$(dir)...; \
	fi

test.service:
	@git-test service $(service)

# test continous integration	

test.droneio:
	@./GoTest.sh ${DRONE_COMMIT_SHA}
	@git-test commit ${CIRCLE_SHA1}

test.circleci:
	@.circleci/circleci.sh

# go build & run

go.ensure.run:
	make go.dep.ensure
	make go.run.service

go.build.service:
	@echo ">>> building $(service)"
	@go build -v -o $(service) svc/$(service)/*.go

go.build.tools:
	@go build -o $(name) tools/$(name)/*.go

go.run.service:
	make go.build.service $(service)
	@echo ">>> running $(service) with env $(env)"
	@TKPENV=$(env) ./$(service) -log_level=debug -config_dir=svc/$(service)/files/config

go.clean.bin:
	- rm bookapp
	- rm rentapp
	- rm ongkirapp
	- rm togel

go.dep.init:
	@echo ">>> dep init $(service)"
	@cd svc/$(service) && dep init -v

go.dep.ensure:
	@echo ">>> dep ensure $(service)"
	@cd svc/$(service) && dep ensure -v --vendor-only

go.dep.update:
	@echo ">>> dep update $(service)"
	@cd svc/$(service) && dep ensure -v -update

# tools

tools.update:
	@echo ">>> updating git-test"
	@go install -v github.com/lab46/monorepo/tools/git-test
	@echo ">>> updating sqlimporter"
	@go install -v github.com/lab46/monorepo/tools/sqlimporter
	@echo ">>> tools update complete"

## docker specific

# docker.go.build function. param1: dockerimage, param2: appname
define docker.go.build
	-docker image rm -f $(1)
	@echo ">>> compiling $(2)"
	@cd svc/$(2) && GOOS=linux go build -o $(2)
	@echo ">>> docker build"
	@cd svc/$(2) && docker build --build-arg service=$(2) . -t $(1)
	@echo ">>> removing $(2) binary"
	@cd svc/$(2) && rm $(2)
endef

docker.build:
	$(eval dockerimage = "$(service):v0.10")
	$(call docker.go.build,$(dockerimage),$(service))


docker.run:
	$(eval dockerimage = "$(service):v0.10")
	@docker run --name $(service) -d $(dockerimage)

docker.stop:
	-docker stop $(service)
	-docker rm -f $(service)