default: format proto_gen swagger_gen tidy

format:
	gofmt -w .

proto_gen:
	@for file in $$(find src -type f -name '*.proto'); do \
		echo $$file; \
		protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative $$file; \
	done
	
swagger_gen:
	swag fmt && swag init -o src/api_service/pkg/docs -g src/api_service/pkg/api/api.go

tidy:
	cd src/api_service && go mod tidy && cd ../..
	cd src/urls_service && go mod tidy && cd ../..

test:
	go test -v ./src/api_service/...
	go test -v ./src/urls_service/...

run_redis:
	docker run -d -p 6379:6379 --name redis_cache redis || docker start redis_cache || echo "redis is running"

run_api:
	env $$(cat .env | grep -v ^# | sed -E -e "s/=.+_service/=localhost/g") go run api_service

run_urls:
	env $$(cat .env | grep -v ^# | sed -E -e "s/=.+_service/=localhost/g") go run urls_service

	
# docker exec -it [container] bash
# docker compose run --rm [container] sh -c "[command]"
# docker compose up --build -V
# docker-compose down --remove-orphans
