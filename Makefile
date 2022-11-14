default: build

build: format proto_gen swagger_gen
	go build -o bin/api_service cmd/api/main.go
	go build -o bin/urls_service cmd/urls/main.go
	go build -o bin/cache_service cmd/cache/main.go
	
format:
	gofmt -w .

proto_gen:
	@for file in $$(find ./pkg -type f -name '*.proto'); do \
		echo $$file; \
		protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $$file; \
	done

swagger_gen:
	swag fmt 
	swag init -g pkg/services/api/api.go

docker_swagger_update: swagger_gen
	docker compose run --rm api_service sh -c "swag fmt; swag init -g pkg/services/api/api.go"

run_redis:
	docker run -d -p 6379:6379 --name redi_cache_local redis

run_cache:
# ./bin/cache_service -redis_addr=localhost:6379
	go run cmd/cache/main.go -redis_addr=localhost:6379

run_urls:
# ./bin/urls_service
	go run cmd/urls/main.go

run_api:
# ./bin/api_service -urls_addr=localhost:8081 -cache_addr=localhost:8082
	go run cmd/api/main.go -urls_addr=localhost:8081 -cache_addr=localhost:8082

clean:
	rm -f bin/*

# docker exec -it [container] bash
# docker compose run --rm [container] sh -c "[command]"
# docker compose up --build -V
# docker-compose down --remove-orphans
