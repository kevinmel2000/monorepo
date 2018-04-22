FROM envoyproxy/envoy:latest

COPY svc2 /usr/local/bin
COPY envoy /envoy
COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh", "svc2"]