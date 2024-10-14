SHELL := /bin/bash

MODULE := $$(go list -m)

.PHONY: default
default: fmt gen

.PHONY: fmt
fmt:
	@find . -type f -name "*.go" ! -path "./pkg/gen/*" ! -path "./vendor/*" | while read -r file; do \
		go fmt "$${file}" 2>&1 | grep -v "is a program, not an importable package"; \
		goimports -w -local $(MODULE) "$${file}"; \
	done

.PHONY: gen
gen:
	@sqlc generate
	@buf format --write && buf generate

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