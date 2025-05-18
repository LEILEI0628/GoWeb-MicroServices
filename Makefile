GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

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
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

PROJECT_ROOT := $(shell pwd)
PROTO_DEST_DIR := api
# 新增第三方依赖路径变量
PROTO_INCLUDE := -I$(PROJECT_ROOT)/third_party
.PHONY: all
all: generate build

# 生成所有 API 代码
.PHONY: generate
generate:
	@find ./api -name '*.proto' -print0 | xargs -0 -t -I{} protoc \
		-I./api \
		-I./third_party \
		--go_out=paths=source_relative:$(PROJECT_ROOT)/$(PROTO_DEST_DIR) \
		--go-http_out=paths=source_relative:$(PROJECT_ROOT)/$(PROTO_DEST_DIR) \
		--go-grpc_out=paths=source_relative:$(PROJECT_ROOT)/$(PROTO_DEST_DIR) \
		--validate_out=paths=source_relative,lang=go:$(PROJECT_ROOT)/$(PROTO_DEST_DIR) {}

# 递归执行子目录的 make
.PHONY: build
build:
	@for dir in $(shell find app -name Makefile); do \
		(cd `dirname $$dir` && make build); \
	done

.PHONY: clean
clean:
	@rm -rf bin/*
	@find . -name '*.pb.go' -delete
	@find . -name '*_http.pb.go' -delete
	@find . -name '*_grpc.pb.go' -delete

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
