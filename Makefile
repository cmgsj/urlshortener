default: format proto_gen swagger_gen

format:
	gofmt -w .

proto_gen:
	@for file in $$(find src/proto/pkg -type f -name '*.proto'); do \
		echo $$file; \
		protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative $$file; \
	done
	
swagger_gen:
	swag fmt 
	swag init -o src/api_service/pkg/docs -g src/api_service/pkg/api/api.go

test:
	go test -v ./src/api_service/...
	go test -v ./src/auth_service/...
	go test -v ./src/cache_service/...
	go test -v ./src/urls_service/...

run_api:
	env $$(cat .env | grep -v ^# | sed -E -e "s/=.+_service/=localhost/g") go run api_service

run_auth:
	env $$(cat .env | grep -v ^# | sed -E -e "s/=.+_service/=localhost/g") go run auth_service

run_cache:
	env $$(cat .env | grep -v ^# | sed -E -e "s/=.+_service/=localhost/g") go run cache_service

run_urls:
	env $$(cat .env | grep -v ^# | sed -E -e "s/=.+_service/=localhost/g") go run urls_service

run_redis:
	docker run -d -p 6379:6379 --name redis_local redis || docker start redis_local || echo "redis is running"

# docker exec -it [container] bash
# docker compose run --rm [container] sh -c "[command]"
# docker compose up --build -V
# docker-compose down --remove-orphans
