SHELL := /bin/bash

.PHONY: default
default: fmt gen

.PHONY: fmt
fmt:
	@go fmt ./...
	@goimports -w -local github.com/cmgsj/urlshortener $$(find . -type f -name "*.go" ! -path "./vendor/*")

.PHONY: gen
gen:
	@sqlc generate --file sqlc.yaml
	@buf format --write proto && buf generate --template proto/buf.gen.yaml proto

.PHONY: test
test:
	@go test -v ./...

.PHONY: docker
docker:
	@eval $$(minikube -p minikube docker-env) && docker build -t cmg/urlshortener:latest .

# minikube start --driver=docker
# minikube addons enable ingress
# minikube addons enable ingress-dns
# kubectl delete -f k8s
# kubectl apply -f k8s
# minikube -n urlshortener service urlshortener
# minikube stop