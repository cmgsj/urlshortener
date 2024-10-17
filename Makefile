SHELL := /bin/bash

MODULE := $$(go list -m)
SWAGGER_UI_VERSION :=

.PHONY: default
default: fmt generate

.PHONY: fmt
fmt:
	@find . -type f -name "*.go" ! -path "./pkg/gen/*" ! -path "./vendor/*" | while read -r file; do \
		go fmt "$${file}" 2>&1 | grep -v "is a program, not an importable package"; \
		goimports -w -local $(MODULE) "$${file}"; \
	done

.PHONY: generate
generate: generate/buf generate/swagger generate/sqlc

.PHONY: generate/buf
generate/buf:
	@rm -rf pkg/gen; \
	find swagger/dist -type f -name '*.swagger.json' -delete; \
	buf format --write; \
	buf lint; \
	buf breaking --against "https://$(MODULE).git#branch=main"; \
	buf generate

.PHONY: generate/swagger
generate/swagger:
	@version=$(SWAGGER_UI_VERSION); \
	if [[ -z "$${version}" ]]; then \
		version="$$(curl -sSL https://api.github.com/repos/swagger-api/swagger-ui/releases/latest | jq -r '.tag_name' | sed 's/^v//')"; \
	fi; \
	rm -rf /tmp/swagger-ui.tar.gz; \
	curl -sSLo /tmp/swagger-ui.tar.gz "https://github.com/swagger-api/swagger-ui/archive/refs/tags/v$${version}.tar.gz"; \
	rm -rf /tmp/swagger-ui; \
	mkdir -p /tmp/swagger-ui; \
	tar -xzf /tmp/swagger-ui.tar.gz -C /tmp/swagger-ui; \
	mkdir -p swagger/dist; \
	find swagger/dist -type f -not -name '*.swagger.json' -delete; \
	cp -r /tmp/swagger-ui/swagger-ui-$${version}/dist/ swagger/dist/; \
	urls="    urls: ["; \
	for file in "$$(find swagger/dist -type f -name "*.swagger.json")"; do \
		path="$${file#swagger/dist/}"; \
		urls+="\n      { name: \"$${path}\", url: \"$${path}\" },\n"; \
	done; \
	urls+="    ],"; \
	line="$$(cat swagger/dist/swagger-initializer.js | grep -n "url" | cut -d: -f1)"; \
	before="$$(head -n "$$(($${line} - 1))" swagger/dist/swagger-initializer.js)"; \
	after="$$(tail -n +"$$(($${line} + 1))" swagger/dist/swagger-initializer.js)"; \
	echo -e "$${before}\n$${urls}\n$${after}" >swagger/dist/swagger-initializer.js

.PHONY: generate/sqlc
generate/sqlc:
	@sqlc generate

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