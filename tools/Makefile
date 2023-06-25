GOPATHBIN=$(GOPATH)/bin
GOENVPATH=$(shell go env GOPATH)

ifeq ($(GOENVPATH),)
	GOENVPATH=$(GOPATHBIN)
else
	GOENVPATH = $(shell go env GOPATH)/bin
endif

.PHONY: all
all: install

.PHONY: install
install: install-swagger install-golangci install-moq

.PHONY: install-swagger
install-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: install-golangci
install-golangci:
	curl -sSfl https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOENVPATH) v1.52.2

.PHONY: install-moq
install-moq:
	go install github.com/matryer/moq@latest

install-dbmate:
	curl -fsSL -o $(GOENVPATH)/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64