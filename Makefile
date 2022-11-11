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

run_redis: 
	docker run --name redis_cache -p 6379:6379 -d redis

run_cache:
	go run cmd/cache/main.go

run_urls:
	go run cmd/urls/main.go

run_server:
	go run cmd/server/main.go

build: format proto_gen swagger_gen
	go build -o bin/server ./cmd/server/main.go
	go build -o bin/cache ./cmd/cache/main.go
	go build -o bin/url ./cmd/url/main.go

clean:
	rm -rf bin/*
