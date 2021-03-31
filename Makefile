GOLANG_VERSION := 1.15.6
ALPINE_VERSION := 3.13

APPNAME := whatismyip

SERVER_ADDR ?= 80

.DEFAULT_GOAL := help
.PHONY: build remove docker-build docker-run docker-stop

help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ \
	{ printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build:	## Build bin
	CGO_ENABLED=1 \
	go build \
	-v \
	-o $(APPNAME) .

remove:	## Remove bin
	rm -rf $(APPNAME)

docker-build:	## Build docker image
	docker build \
	--pull \
	--build-arg GOLANG_VERSION="$(GOLANG_VERSION)" \
	--build-arg ALPINE_VERSION="$(ALPINE_VERSION)" \
	--build-arg APPNAME="$(APPNAME)" \
	--tag="$(APPNAME):latest" \
	.

docker-run:	## Run docker
	docker run --rm -d -p 8080:8080 --name $(APPNAME) $(APPNAME):latest && \
	docker ps

docker-stop:	## Stop docker
	docker stop $(APPNAME)