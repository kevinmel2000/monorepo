FROM envoyproxy/envoy:latest

COPY svc1 /usr/local/bin
COPY envoy /envoy
COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh", "svc1"]