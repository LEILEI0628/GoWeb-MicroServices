SERVICE_NAME := interactive

GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
PROJECT_ROOT := $(shell cd ../../ && pwd)
BIN_DIR := $(PROJECT_ROOT)/bin
GOHOSTOS:=$(shell go env GOHOSTOS)

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	#go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
		   --proto_path=$(PROJECT_ROOT)/third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: wire
# generate wire
wire:
	cd cmd/server && wire

.PHONY: build
build:
	@echo "Building $(SERVICE_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(SERVICE_NAME) ./cmd/server

# 单独生成该服务的 API 代码
.PHONY: generate
generate:
	@protoc \
		-I$(PROJECT_ROOT)/api \
		-I$(PROJECT_ROOT)/third_party \
		--go_out=paths=source_relative:$(PROJECT_ROOT)/api \
		--go-http_out=paths=source_relative:$(PROJECT_ROOT)/api \
		--go-grpc_out=paths=source_relative:$(PROJECT_ROOT)/api \
		--validate_out=paths=source_relative,lang=go:$(PROJECT_ROOT)/api \
		$(PROJECT_ROOT)/api/$(SERVICE_NAME)/v1/*.proto

.PHONY: clean
clean:
	@rm -f $(BIN_DIR)/$(SERVICE_NAME)
	@find . -name '*.pb.go' -delete

.PHONY: all
# generate all
all:
	make generate;
	make wire;
	make config;
	make build;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
