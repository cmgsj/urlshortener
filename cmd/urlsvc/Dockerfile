FROM golang:1.20-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go test -v ./...
RUN go build -o bin/urlsvc ./cmd/urlsvc

FROM alpine:latest
RUN apk add ca-certificates
COPY --from=builder /app/bin/urlsvc /usr/local/bin/urlsvc
USER nobody
ENTRYPOINT ["urlsvc"]