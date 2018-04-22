prepare:
	@cd svc1 && GOOS=linux go build
	@cd svc2 && GOOS=linux go build

run:
	@docker-compose up -d

stop:
	@docker-compose down
	-docker image rm -f grpc_svc1
	-docker image rm -f grpc_svc2-1
	-docker image rm -f grpc_svc2-2