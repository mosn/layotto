SHELL := /bin/bash

.PHONY: all
all: lint test go.build

# ==============================================================================
# ROOT Options

ROOT_PACKAGE=mosn.io/layotto

# ==============================================================================
# Includes

include make/common.mk
include make/golang.mk
include make/image.mk

# ==============================================================================
# Targets

## go.build: Build layotto for host platform.
.PHONY: go.build
build:
	@$(MAKE) go.build

## go.build.multiarch: Build layotto for multiple platforms. See option PLATFORMS.
.PHONY: go.build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch

## go.wasm: Build layotto wasm for host platform.
.PHONY: go.wasm
go.wasm:
	@$(MAKE) build.wasm

## go.wasm.multiarch: Build layotto wasm for multiple platform.
.PHONY: go.wasm.multiarch
go.wasm.multiarch:
	@$(MAKE) build.wasm.multiarch

## app: Build app docker images for host arch.
.PHONY: app
app:
	@$(MAKE) app.image

## app.multiarch: Build app docker images for multiple platforms. See option PLATFORMS.
.PHONY: app.multiarch
app.multiarch:
	@$(MAKE) app.image.multiarch

## image: Build docker images for host arch.
.PHONY: image
image:
	@$(MAKE) image.build

## image.multiarch: Build docker images for multiple platforms. See option PLATFORMS.
.PHONY: image.multiarch
image.multiarch:
	@$(MAKE) image.build.multiarch

## push: Build docker images for host arch and push images to registry.
.PHONY: push
push:
	@$(MAKE) image.push

## push.multiarch: Build docker images for multiple platforms and push images to registry.
.PHONY: push.multiarch
push.multiarch:
	@$(MAKE) image.push.multiarch

## lint: Run go syntax and styling of go sources.
.PHONY: lint
lint:
	@$(MAKE) go.lint

## test: Run go unit test.
.PHONY: test
test:
	@$(MAKE) go.test

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@$(MAKE) go.clean

## check.coverage: Run coverage analysis.
.PHONY: check.coverage
check.coverage:
	@$(MAKE) checker.coverage

## check.deadlink: Run deadlink check test.
.PHONY: check.deadlink
check.deadlink:
	@$(MAKE) checker.deadlink

## check.quickstart: Run quickstart check test.
.PHONY: check.quickstart
check.quickstart:
	@$(MAKE) checker.quickstart

## integrate.wasm: Run integration test with wasm.
.PHONY: integrate.wasm
integrate.wasm:
	@$(MAKE) integration.wasm

## integrate.runtime: Run integration test with runtime.
.PHONY: integrate.runtime
integrate.runtime:
	@$(MAKE) integration.runtime

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  DEBUG        Whether to generate debug symbols. Default is 0.
  BINS         The binaries to build. Default is all of cmd.
               This option is available when using: make build/build.multiarch
               Example: make build BINS="layotto_multiple_api layotto"
  IMAGES       Backend images to make. Default is all of cmds.
               This option is available when using: make image/image.multiarch/push/push.multiarch
               Example: make image.multiarch IMAGES="layotto_multiple_api layotto"
  PLATFORMS    The multiple platforms to build. Default is linux_amd64 and linux_arm64.
               This option is available when using: make build.multiarch/image.multiarch/push.multiarch
               Example: make image.multiarch IMAGES="layotto_multiple_api layotto" PLATFORMS="linux_amd64 linux_arm64"
               Supported Platforms: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64
endef
export USAGE_OPTIONS

# ==============================================================================

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "Layotto: \nA fast and efficient cloud native application runtime\n"
	@echo -e "Usage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"

build-linux-wasm-layotto:
	docker build --rm -t ${BUILD_IMAGE} build/contrib/builder/image/faas
	docker run --rm -v $(shell pwd):/go/src/${PROJECT_NAME} -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} go build -tags wasmer -o layotto /go/src/${PROJECT_NAME}/cmd/layotto