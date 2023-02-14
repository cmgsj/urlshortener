default: 
	@echo "unspecifed target" && exit 1

minikube: 
	minikube start --driver=docker
	eval $$(minikube -p minikube docker-env)
	make docker_build
	kubectl apply -f k8s
	minikube service api-service --url

docker_build: build
	docker build -t cmg/api-svc -f cmd/api_service/Dockerfile .
	docker build -t cmg/url-svc -f cmd/url_service/Dockerfile .

build: proto_gen swagger_gen
	GOOS=linux go build -o bin ./cmd/api_service
	CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o bin ./cmd/url_service

proto_gen:
	@for file in $$(find pkg/proto -type f -name '*.proto'); do \
		echo $$file; \
		protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative $$file; \
	done
	
swagger_gen:
	swag fmt pkg/api_service && swag init -o pkg/api_service/docs -g pkg/api_service/service.go

clean:
	rm -f bin/*