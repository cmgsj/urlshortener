FROM alpine:latest
COPY bin/urlshortener /usr/local/bin/urlshortener
USER nobody
ENTRYPOINT ["urlshortener"]
