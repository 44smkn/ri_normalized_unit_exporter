DOCKER_ARCHS ?= amd64 arm64
DOCKER_IMAGE_NAME ?= ri-normalized-unit-exporter

all:: vet test-e2e common-all

include Makefile.common

PROMTOOL_DOCKER_IMAGE ?= $(shell docker pull -q quay.io/prometheus/prometheus:latest || echo quay.io/prometheus/prometheus:latest)
PROMTOOL ?= docker run -i --rm -w "$(PWD)" -v "$(PWD):$(PWD)" --entrypoint promtool $(PROMTOOL_DOCKER_IMAGE)

.PHONY: test-e2e
test-e2e: build
	@echo ">> running end-to-end tests"
	./e2e-test.sh
