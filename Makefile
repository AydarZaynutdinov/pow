.PHONY: test
test:
	go test -v ./...

.PHONY: build_server
build_server:
	docker build -t server:latest -f build/server.Dockerfile .

.PHONY: build_client
build_client:
	docker build -t client:latest -f build/client.Dockerfile .

.PHONY: build_all
build_all: build_server build_client

.PHONY: up
up:
	docker-compose up --build -d

.PHONY: down
down:
	docker-compose down -v