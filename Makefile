REGISTRY:=quay.io
PACKAGE:=jeffersongirao/hello-k8s

VERSION:=0.1.0
SHELL:=$(shell which bash)
DOCKER:=$(shell command -v docker)
IMAGE:=$(REGISTRY)/$(PACKAGE)

deps-development:
ifndef DOCKER
	@echo "Docker is not available. Please install docker"
	@exit 1
endif

image: deps-development
	docker build \
	-t $(IMAGE):latest \
	-t $(IMAGE):$(VERSION) \
	-f docker/Dockerfile \
	.

publish: deps-development
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest