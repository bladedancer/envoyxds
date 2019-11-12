PROJECT_NAME ?= envoyxds
PKG_LIST := $(shell go list ./... | grep -v /vendor/)
REGISTRY ?= bladedancer

PKG := bladedancer/$(PROJECT_NAME)

.PHONY: clean

all: clean lint protoc build docker-build push ## Build everything

lint: ## Lint the files
	@golint	-set_exit_status	${PKG_LIST}

build: ## Build the binary for linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build	-o ./bin/$(PROJECT_NAME)	./cmd/$(PROJECT_NAME)

docker-build: ## Build docker image
	docker build -f ./cmd/$(PROJECT_NAME)/Dockerfile -t $(REGISTRY)/$(PROJECT_NAME):latest	.

push: ## Push docker image
	docker push $(REGISTRY)/$(PROJECT_NAME):latest

clean: ## Clean out dir
	rm -rf ./bin && \
    docker rmi -f $(REGISTRY)/$(PROJECT_NAME):latest

help: ## Display this help screen
	@grep	-h	-E	'^[a-zA-Z_-]+:.*?## .*$$'	$(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

PROTODIRS := pkg
WORKSPACE ?= $$(pwd)
# standard protobuf files
PROTOFILES := $(shell find $(PROTODIRS) -type f -name '*.proto')
PROTOTARGETS := $(PROTOFILES:.proto=.pb.go)

PROTOOPTS := -I/go/src/ \
	--go_out=plugins=grpc:/go/src/

PROTOALLTARGETS := $(PROTOTARGETS)

%.pb.go %.pb.gw.go %.swagger.json: %.proto
	@echo $<
#	@protoc $(PROTOOPTS) $(GOPATH)/src/$(REPO)/$<

	@docker run -i --rm  -v "$(WORKSPACE):/go/src/$(PKG)" \
	-u $$(id -u):$$(id -g)                    \
	chrisccoy/golang-gw:1.0.0 	protoc \
	-I /go/src -I/go/src/$(PKG)/vendor \
	--go_out=plugins=grpc:/go/src  \
	/go/src/$(PKG)/$<

protoc: $(PROTOALLTARGETS)

