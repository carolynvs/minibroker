PKG = github.com/deis/minibroker
DOCKER_IMG = minibroker-build

USE_DOCKER ?= true

ifeq ($(USE_DOCKER),true)
  DO = docker run --rm -it -v $$HOME/.kube:/root/.kube -v $$HOME/.minikube:$$HOME/.minikube -v $$(pwd):/go/src/$(PKG) $(DOCKER_IMG)
else
  DO =
endif

default: build

.PHONY: buildimage build run create-cluster init

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
