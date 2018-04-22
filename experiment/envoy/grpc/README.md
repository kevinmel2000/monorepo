# gRPC Experiment

This experiment is around gRPC environment

## gRPC Load Balancing Via Envoy

To test load balancing gRPC via envoy do this step:

1. `make prepare`to build all go binaries

2. `make run` to run all services
After all services running, try to hit `curl localhost:9090/pingserver`. You will see response from different server everytime you hit.

    `/pingserver` is calling grpc function to call other services via `envoy` and envoy will load-balance the traffic using `round-robin`.

3. `make stop` to stop docker compose and delete the images