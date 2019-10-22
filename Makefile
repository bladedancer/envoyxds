GO_VERSION=1.11
GLIDE_DOCKER_VERSION=0.12.3-go1.8

PROJECT_NAME := envoyxds
PKG_LIST := $(shell go list ./... | grep -v /vendor/)

all: lint build docker ## Build everything

lint: ## Lint the files
	@golint	-set_exit_status	${PKG_LIST}

build: ## Build the binary for linux
	CGO_ENABLED=0 GOOS=linux go build	-o bin/$(PROJECT_NAME)	./cmd/envoyxds

docker:
	APP=$(PROJECT_NAME) docker build -t docker.pkg.github.com/bladedancer/$(PROJECT_NAME):latest_dev	.

push:
	docker push docker.pkg.github.com/bladedancer/$(PROJECT_NAME):latest_dev

clean:
	@rm	-rf	${OUT_DIR}

help: ## Display this help screen
	@grep	-h	-E	'^[a-zA-Z_-]+:.*?## .*$$'	$(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
