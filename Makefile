PKG = github.com/deis/minibroker
DOCKER_IMG = minibroker-build

USE_DOCKER ?= true

build:

ifeq ($(USE_DOCKER),true)
  DO = docker run --rm -it -v $$HOME/.kube:/root/.kube -v $$HOME/.minikube:$$HOME/.minikube -v $$(pwd):/go/src/$(PKG) $(DOCKER_IMG)
else
  DO =
endif

default: build

.PHONY: buildimage build run create-cluster init test clean deploy

buildimage:
	docker build -t $(DOCKER_IMG) ./build

build: buildimage
	$(DO) ./build/build.sh

run: buildimage
	$(DO) ./build/run.sh
	$(DO) svcat get brokers

create-cluster:
	./build/create-cluster.sh

init: buildimage
	$(DO) ./build/init.sh

test:
	$(DO) ./build/test.sh

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build --ldflags="-s" -i bin/linux/amd64/minibroker ./cmd/minibroker

image: linux
	cp minibroker image/
	docker build image/ -t minibroker

clean:
	rm -f bin
