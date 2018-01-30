FROM envoyproxy/envoy:latest

COPY bookapp /usr/local/bin
COPY ./config /etc/bookapp
COPY ./envoy /envoy
COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh", "bookapp", "test"]