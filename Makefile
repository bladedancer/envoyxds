PROJECT_NAME := envoyxds
PKG_LIST := $(shell go list ./... | grep -v /vendor/)

.PHONY: clean

all: clean lint build docker-build push ## Build everything

lint: ## Lint the files
	@golint	-set_exit_status	${PKG_LIST}

build: ## Build the binary for linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build	-o ./bin/$(PROJECT_NAME)	./cmd/envoyxds

docker-build: ## Build docker image
	docker build -t docker.pkg.github.com/bladedancer/$(PROJECT_NAME)/$(PROJECT_NAME):latest_dev	.

push: ## Push docker image
	docker push docker.pkg.github.com/bladedancer/$(PROJECT_NAME)/$(PROJECT_NAME):latest_dev

clean: ## Clean out dir
	rm -rf ./bin && \
    docker rmi -f docker.pkg.github.com/bladedancer/$(PROJECT_NAME)/$(PROJECT_NAME):latest_dev

help: ## Display this help screen
	@grep	-h	-E	'^[a-zA-Z_-]+:.*?## .*$$'	$(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
