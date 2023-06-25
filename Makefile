BUILD_TAG = $(shell git describe --tags $(git describe --tags --abbrev=0 --always) --always)
BUILD_TIME = $(shell date +%s)
GO_VERSION = $(shell go version | cut -c 14-)

OS=linux
ARCH=amd64
GO_BUILD_PARAMS = -ldflags="-X 'main.AppVersion=$(BUILD_TAG)' -X 'main.GoVersion=$(GO_VERSION)' -X 'main.Time=$(BUILD_TIME)'"

ifeq ($(shell uname -s), Darwin)
	OS=darwin
endif

ifeq ($(shell uname -m), arm64)
	ARCH=arm64
endif

.PHONY: install-tools
install-tools:
	make -f tools/Makefile install

.PHONY: install-godeps
install-godeps:
	go mod tidy

.PHONY: generate
generate:
	go generate ./pkg/...

.PHONY: deps
deps: install-godeps install-tools generate

.PHONY: leverage-test-deps
leverage-test-deps:
	@docker compose -f deploy/development/docker-compose.tests.yaml up -d

.PHONY: test
test: leverage-test-deps
	go test ./... -v -race
	@docker compose -f deploy/development/docker-compose.tests.yaml down

.PHONY: test-cov
test-cov:
	go test ./... -v -race --cover

test-cov-docker:leverage-test-deps
	go test ./... -v -race --cover
	@docker compose -f deploy/development/docker-compose.tests.yaml down

.PHONY: lint
lint:
	golangci-lint run ./... -c .golangci.yaml

.PHONY: lint-fix
lint-fix:
	golangci-lint run ./... -c .golangci.yaml --fix

.PHONY: swagger
swagger:
	 swag init --parseDependency -o api -g cmd/app/main.go

.PHONY: migrate
migrate:
	dbmate --url '${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?${DB_OPTIONS}' up

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build $(GO_BUILD_PARAMS) -o bin/app ./cmd/app

build-docker: build
	docker build -t pedroyremolo/identity-service:latest -t pedroyremolo/identity-service:$(BUILD_TAG) .