FROM golang:alpine AS builder
WORKDIR /src
COPY . .
RUN apk add --update gcc musl-dev
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -ldflags "-s -w -linkmode=external -extldflags='-static'" -o bin/urlshortener ./cmd/urlshortener

FROM alpine:latest
COPY --from=builder /src/bin/urlshortener /urlshortener
USER nobody
CMD ["/urlshortener"]
