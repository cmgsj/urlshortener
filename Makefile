default: 
	@echo unspecifed target && exit 1

minikube: 
	minikube start --driver=docker
	make docker_build
	kubectl apply -f k8s
	minikube service web-service --url

docker_build: build
	eval $$(minikube -p minikube docker-env)
	docker build -t cmg/web-svc -f cmd/websvc/Dockerfile .
	docker build -t cmg/url-svc -f cmd/urlsvc/Dockerfile .

build: proto_gen swagger_gen
	GOOS=linux go build -o bin ./cmd/websvc
	CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o bin ./cmd/urlsvc

proto_gen:
	@for file in $$(find pkg/proto -type f -name '*.proto'); do \
		echo $$file; \
		protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative $$file; \
	done
	
swagger_gen:
	swag fmt pkg/websvc
	swag init -o pkg/websvc/docs -g pkg/websvc/service.go

sqlc_gen:
	sqlc generate -f pkg/urlsvc/sqlc.yaml

clean:
	rm -f bin/*