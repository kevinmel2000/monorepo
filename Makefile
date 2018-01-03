build_bookapp:
	@go build -o book bookapp/*.go

build_rentapp:
	@go build -o rent rentapp/*.go

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