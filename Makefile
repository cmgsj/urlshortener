all: build

docker_build: build
	eval $$(minikube -p minikube docker-env)
	docker build -t cmg/urlshortener:latest -f ./cmd/urlshortener/Dockerfile .

build: gen
	CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -mod=vendor -trimpath -ldflags "-linkmode external -extldflags -static" -o bin ./cmd/urlshortener

gen:
	sqlc generate -f sqlc.yaml
	buf generate --exclude-path vendor

install_tools:
	grep _ pkg/tools/tools.go | awk -F'"' '{print $$2}' | xargs -tI % go install %

clean:
	rm -f bin/*

# minikube start --driver=docker
# minikube addons enable ingress
# minikube -n urlshortener service urlshortener --url
# minikube stop