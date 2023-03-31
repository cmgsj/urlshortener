default: 
	@echo unspecifed target && exit 1

minikube: 
	minikube start --driver=docker

kube_apply: docker_build
	kubectl apply -f k8s
	kubens urlshortener

kube_port_forward_web:
	kubectl port-forward service/web-service 8080:8080

kube_port_forward_url:
	kubectl port-forward service/url-service 8080:8080

kube_delete:
	kubens default
	kubectl delete -f k8s

docker_build: build
	eval $$(minikube -p minikube docker-env)
	docker build -t cmg/web-svc -f cmd/websvc/Dockerfile .
	docker build -t cmg/url-svc -f cmd/urlsvc/Dockerfile .

build: gen
	GOARCH=amd64 GOOS=linux go build -o bin ./cmd/websvc
	CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o bin ./cmd/urlsvc

gen:
	swag fmt pkg/websvc
	swag init -o pkg/websvc/docs -g pkg/websvc/service.go
	sqlc generate -f url.v1.sqlc.yaml
	buf generate
	cp pkg/gen/proto/url/v1/url.swagger.json swagger.json

install_tools:
	grep _ pkg/tools/tools.go | awk -F'"' '{print $$2}' | xargs -tI % go install %

clean:
	rm -f bin/*
