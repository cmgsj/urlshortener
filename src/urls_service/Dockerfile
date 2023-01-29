FROM golang

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .
COPY pkg pkg
COPY main.go .

RUN go mod download

RUN go install github.com/githubnemo/CompileDaemon@latest
CMD CompileDaemon --build="go build -o ./bin/app ." --command=./bin/app

# RUN go build -o ./bin/app .
# CMD ./bin/app