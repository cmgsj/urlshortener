default: 
	@echo unspecifed target && exit 1

minikube: 
	minikube start --driver=docker

kube: docker_build
	kubectl apply -f k8s
	kubens urlshortener

kube_port_forward:
	kubectl port-forward service/web-service 8080:8080

kube_delete:
	kubens default
	kubectl delete -f k8s

docker_build: gen
	eval $$(minikube -p minikube docker-env)
	docker build -t cmg/web-svc -f cmd/websvc/Dockerfile .
	docker build -t cmg/url-svc -f cmd/urlsvc/Dockerfile .

build: gen
	go build -o bin ./cmd/websvc
	go build -o bin ./cmd/urlsvc

install_tools:
	grep _ pkg/tools/tools.go | awk -F'"' '{print $$2}' | xargs -tI % go install %

gen:
	swag fmt pkg/websvc
	swag init -o pkg/websvc/docs -g pkg/websvc/service.go
	sqlc generate -f url.v1.sqlc.yaml
	buf generate

clean:
	rm -f bin/*
