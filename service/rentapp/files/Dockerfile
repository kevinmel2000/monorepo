FROM envoyproxy/envoy:latest

COPY rentapp /usr/local/bin
COPY ./config /etc/rentapp
COPY ./envoy /envoy
COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh", "rentapp", "test"]