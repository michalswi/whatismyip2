GOLANG_VERSION := 1.15.6
ALPINE_VERSION := 3.13

DOCKER_REPO := michalsw
APPNAME := whatismyip

SERVER_PORT ?= 8080

AZ_LOCATION ?= westeurope
AZ_RANDOM ?=$(shell head /dev/urandom | tr -dc a-z0-9 | head -c 7)
AZ_DNS_LABEL ?= $(APPNAME)-$(AZ_RANDOM)

.DEFAULT_GOAL := help
.PHONY: build docker-build docker-run docker-stop

help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ \
	{ printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build bin
	CGO_ENABLED=1 \
	go build \
	-v \
	-o $(APPNAME) .

docker-build: ## Build docker image
	docker build \
	--pull \
	--build-arg GOLANG_VERSION="$(GOLANG_VERSION)" \
	--build-arg ALPINE_VERSION="$(ALPINE_VERSION)" \
	--build-arg APPNAME="$(APPNAME)" \
	--tag="$(DOCKER_REPO)/$(APPNAME):latest" \
	.

docker-run: ## Run docker
	docker run --rm -d \
	--name $(APPNAME) \
	-p $(SERVER_PORT):$(SERVER_PORT) \
	$(DOCKER_REPO)/$(APPNAME):latest &&\
	docker ps

docker-stop: ## Stop docker
	docker stop $(APPNAME)

azure-rg:	## Create the Azure Resource Group
	az group create --name $(APPNAME)-rg --location $(AZ_LOCATION)

azure-rg-del:	## Delete the Azure Resource Group
	az group delete --name $(APPNAME)-rg

azure-aci:	## Run Azure Container Instance
	az container create \
	--resource-group $(APPNAME)-rg \
	--name $(APPNAME) \
	--image $(DOCKER_REPO)/$(APPNAME) \
	--restart-policy Always \
	--ports $(SERVER_PORT) \
	--dns-name-label $(AZ_DNS_LABEL) \
	--location $(AZ_LOCATION) \
	--environment-variables \
		SERVER_PORT=$(SERVER_PORT)
