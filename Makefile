default: format proto_gen swagger_gen
	
format:
	gofmt -w .

proto_gen:
	@for file in $$(find ./pkg -type f -name '*.proto'); do \
		echo $$file; \
		protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $$file; \
	done

swagger_gen:
	swag fmt 
	swag init -g pkg/api/api.go

docker_swagger_update: swagger_gen
	docker compose run --rm api_service sh -c "swag fmt; swag init -g pkg/api/api.go"

# docker exec -it [container] bash
# docker compose run --rm [container] sh -c "[command]"
# docker compose up --build -V
# docker-compose down --remove-orphans
