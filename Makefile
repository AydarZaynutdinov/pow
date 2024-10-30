.PHONY: deps
deps:
	@echo 'Install dependencies'
	go mod tidy -v

.PHONY: build-pow
build-pow:
	@echo 'Building pow app'
	go build -buildvcs=false -o ./bin/pow ./cmd/pow-server

.PHONY: build
build: deps build-pow

.PHONY: run-pow
run-pow:
	./bin/pow -config ./test/config/config.yaml

.PHONY: clean-testcache-environment
clean-testcache-environment:
	go clean -testcache

.PHONY: unit-test
unit-test: clean-testcache-environment
	@echo 'Running unit tests...'
	go test -race -short ./...

.PHONY: integration-test
integration-test: clean-testcache-environment
	@echo 'Running integration tests...'
	go test -race -run Integration ./test

.PHONY: all-test
all-test: unit-test integration-test

.PHONY: all-test-one-command
all-test-one-command: build
	docker-compose -f ./scripts/test/docker-compose.yml up --build -d
	- make run-pow >/dev/null 2>&1 &
	go clean -testcache
	make all-test;\
	EXIT_CODE=$$?;\
	ps aux | grep -i ./bin/pow | grep -v grep | awk {'print $$2'} | xargs kill -9;\
	docker-compose -f ./scripts/test/docker-compose.yml down -v;\
	exit $$EXIT_CODE