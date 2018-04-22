FROM envoyproxy/envoy:latest

COPY bookapp /usr/local/bin
COPY files/config /etc/bookapp
COPY files/envoy /envoy
COPY files/entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh", "bookapp", "test"]