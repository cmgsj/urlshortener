default: format
	
format:
	gofmt -w .

proto_gen:
	@for file in $$(find ./pkg -type f -name '*.proto'); do \
		echo $$file; \
		protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $$file; \
	done

swagger_gen:
	swag fmt 
	swag init -g pkg/server/server.go

build: format proto_gen swagger_gen
	go build -o bin/server ./cmd/server/main.go
	go build -o bin/cache ./cmd/cache/main.go
	go build -o bin/urls ./cmd/urls/main.go

clean:
	rm -rf bin/*

# docker exec -it urls_service bash
# docker compose run --rm urls_service sh -c "ls -lh"
# docker compose up
# docker compose up --build -V
# docker-compose down
# docker-compose down --remove-orphans
