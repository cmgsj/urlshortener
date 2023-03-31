FROM alpine:latest
COPY bin/urlsvc /usr/local/bin/app
COPY swagger-ui /usr/local/bin/swagger-ui
COPY swagger.json /usr/local/bin/swagger.json
ENTRYPOINT ["app"]