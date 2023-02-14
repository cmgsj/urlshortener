default: 
	@echo "unspecifed target" && exit 1

minikube: docker_build
	minikube start --driver=docker
	kubectl apply -f k8s
	minikube service api-service --url

docker_build: build
	eval $$(minikube -p minikube docker-env)
	cd src/api_service && docker build --no-cache -t cmg/api-svc . && cd ../..
	cd src/url_service && docker build --no-cache -t cmg/url-svc . && cd ../..

build: proto_gen swagger_gen
	cd src/api_service && GOOS=linux \
		go build -o bin ./cmd/api_service && cd ../..
	cd src/url_service && CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 \
		go build -ldflags "-linkmode external -extldflags -static" -o bin ./cmd/url_service && cd ../..

proto_gen:
	@for file in $$(find src/proto -type f -name '*.proto'); do \
		echo $$file; \
		protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative $$file; \
	done
	
swagger_gen:
	swag fmt && swag init -o src/api_service/pkg/docs -g src/api_service/pkg/api/service.go
