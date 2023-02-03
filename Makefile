default: build

build: proto_gen swagger_gen
	cd src/api_service && GOOS=linux \
		go build -o bin ./cmd/api_service && cd ../..
	cd src/urls_service && CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 \
		go build -ldflags "-linkmode external -extldflags -static" -o bin ./cmd/urls_service && cd ../..

proto_gen:
	@for file in $$(find src/proto -type f -name '*.proto'); do \
		echo $$file; \
		protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative $$file; \
	done
	
swagger_gen:
	swag fmt && swag init -o src/api_service/pkg/docs -g src/api_service/pkg/api/service.go

run_redis:
	docker run -d -p 6379:6379 --name redis_cache redis || docker start redis_cache || echo "redis is running"

run_api:
	env $$(cat .env | grep -v ^# | sed -E -e "s/=.+_service/=localhost/g") go run src/api_service/cmd/api_service/main.go

run_urls:
	env $$(cat .env | grep -v ^# | sed -E -e "s/=.+_service/=localhost/g") go run src/urls_service/cmd/urls_service/main.go
	
# docker exec -it [container] bash
# docker compose run --rm [container] sh -c "[command]"
# docker compose up --build -V
# docker-compose down --remove-orphans
