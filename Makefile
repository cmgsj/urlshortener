default: build

build: format proto_gen swagger_gen
	go build -o bin/api_service cmd/api/main.go
	go build -o bin/urls_service cmd/urls/main.go
	go build -o bin/cache_service cmd/cache/main.go
	
format:
	gofmt -w .

proto_gen:
	@for file in $$(find pkg/proto -type f -name '*.proto'); do \
		echo $$file; \
		protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative $$file; \
	done
	
swagger_gen:
	swag fmt 
	swag init -o pkg/services/api/api_docs -g pkg/services/api/api.go

test:
	go test -v ./...

clean:
	rm -f bin/*

docker_swagger_update: swagger_gen
	docker compose run --rm api_service sh -c "swag init -o pkg/services/api/api_docs -g pkg/services/api/api.go"

run_redis:
	docker run -d -p 6379:6379 --name redi_cache_local redis || docker start redi_cache_local || echo "docker daemon not running"" 

run_cache:
	go run cmd/cache/main.go -redis_addr=localhost:6379

run_urls:
	go run cmd/urls/main.go

run_api:
	go run cmd/api/main.go -urls_addr=localhost:8081 -cache_addr=localhost:8082

# grpc-health-probe -addr=localhost:8081 -service=urlspb.UrlsService
# docker exec -it [container] bash
# docker compose run --rm [container] sh -c "[command]"
# docker compose up --build -V
# docker-compose down --remove-orphans
