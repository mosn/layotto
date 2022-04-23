# Layotto Commands: 
# A fast and efficient cloud native application runtime

# Usage: make <TARGETS> <OPTIONS> ...

# Targets:
#   go.build             Build layotto for host platform.
#   go.build.multiarch   Build layotto for multiple platforms. See option PLATFORMS.
#   go.wasm              Build layotto wasm for linux arm64 platform.
#   go.wasm.multiarch    Build layotto wasm for multiple platform.
#   go.wasm.image        Build layotto wasm image for multiple platform.
#   go.wasm.image.push   Push layotto wasm image for multiple platform.
#   go.check             Run all go checks of code sources.
#   go.check.lint        Run go syntax and styling of go sources.
#   go.check.unit        Run go unit test.
#   go.check.style       Run go style test.
#   image                Build docker images for host arch.
#   image.multiarch      Build docker images for multiple platforms. See option PLATFORMS.
#   push                 Push docker images to registry.
#   push.multiarch       Push docker images for multiple platforms to registry.
#   app                  Build app docker images for host arch. [`/docker/app` contains apps dockerfiles]
#   app.multiarch        Build app docker images for multiple platforms. See option PLATFORMS.
#   style.coverage       Run coverage analysis.
#   style.deadlink       Run deadlink check test.
#   style.quickstart     Run quickstart check test.
#   integrate.wasm       Run integration test with wasm.
#   integrate.runtime    Run integration test with runtime.
#   clean                Remove all files that are created by building.
#   help                 Show this help info.

# Options:
#   DEBUG        Whether to generate debug symbols. Default is 0.
#   BINS         The binaries to build. Default is all of cmd.
#                This option is available when using: make build/build.multiarch
#                Example: make build BINS="layotto_multiple_api layotto"
#   IMAGES       Backend images to make. Default is all of cmds.
#                This option is available when using: make image/image.multiarch/push/push.multiarch
#                Example: make image.multiarch IMAGES="layotto_multiple_api layotto"
#   PLATFORMS    The multiple platforms to build. Default is linux_amd64 and linux_arm64.
#                This option is available when using: make build.multiarch/image.multiarch/push.multiarch
#                Example: make image.multiarch IMAGES="layotto" PLATFORMS="linux_amd64 linux_arm64"
#                Supported Platforms: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64

SHELL := /bin/bash

.PHONY: all
all: go.lint go.test go.build


.PHONY: check
check: go.check.style go.check.unit go.check.lint

# ==============================================================================
# ROOT Options

ROOT_PACKAGE=mosn.io/layotto

# ==============================================================================
# Includes

include make/common.mk
include make/golang.mk
include make/image.mk
include make/wasm.mk
include make/ci.mk

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

## go.wasm: Build layotto wasm for linux arm64 platform.
.PHONY: go.wasm
go.wasm:
	@$(MAKE) wasm

## go.wasm.multiarch: Build layotto wasm for multiple platform.
.PHONY: go.wasm.multiarch
go.wasm.multiarch:
	@$(MAKE) wasm.multiarch

## go.wasm.image: Build layotto wasm image for multiple platform.
.PHONY: go.wasm.image
go.wasm.image:
	@$(MAKE) wasm.image

## go.wasm.image.push: Push layotto wasm image for multiple platform.
.PHONY: go.wasm.image.push
go.wasm.image.push:
	@$(MAKE) wasm.image.push

## go.check: Run all go checks of code sources.
.PHONY: go.check
go.check: go.check.style go.check.unit go.check.lint

## go.check.lint: Run go syntax and styling of go sources.
.PHONY: go.check.lint
go.check.lint:
	@$(MAKE) go.lint

## go.check.unit: Run go unit test.
.PHONY: go.check.unit
go.check.unit:
	@$(MAKE) go.test

## go.check.style: Run go style test.
.PHONY: go.check.style
go.check.style:
	@$(MAKE) go.style

## image: Build docker images for host arch.
.PHONY: image
image:
	@$(MAKE) image.build

## image.multiarch: Build docker images for multiple platforms. See option PLATFORMS.
.PHONY: image.multiarch
image.multiarch:
	@$(MAKE) image.build.multiarch

## push: Push docker images to registry.
.PHONY: push
push:
	@$(MAKE) image.push

## push.multiarch: Push docker images for multiple platforms to registry.
.PHONY: push.multiarch
push.multiarch:
	@$(MAKE) image.push.multiarch

## app: Build app docker images for host arch. [`/docker/app` contains apps dockerfiles]
.PHONY: app
app:
	@$(MAKE) app.image

## app.multiarch: Build app docker images for multiple platforms. See option PLATFORMS.
.PHONY: app.multiarch
app.multiarch:
	@$(MAKE) app.image.multiarch

## style.coverage: Run coverage analysis.
.PHONY: style.coverage
style.coverage:
	@$(MAKE) checker.coverage

## style.deadlink: Run deadlink check test.
.PHONY: style.deadlink
check.deadlink:
	@$(MAKE) checker.deadlink

## style.quickstart: Run quickstart check test.
.PHONY: style.quickstart
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

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@$(MAKE) go.clean

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
               Example: make image.multiarch IMAGES="layotto" PLATFORMS="linux_amd64 linux_arm64"
               Supported Platforms: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64
endef
export USAGE_OPTIONS

# ==============================================================================

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "Layotto Commands: \nA fast and efficient cloud native application runtime\n"
	@echo -e "Usage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
