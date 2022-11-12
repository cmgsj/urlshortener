FROM golang

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .
COPY pkg pkg
COPY cmd/server/main.go .

RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag fmt 
RUN swag init -g pkg/server/server.go

RUN go install github.com/githubnemo/CompileDaemon@latest

CMD CompileDaemon --build="go build -o app ." --command=./app

# RUN go build -o app .

# CMD ./app
