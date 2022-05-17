# Copyright 2021 Layotto Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Layotto commands 👀: 

# A fast and efficient cloud native application runtime 🚀.
# Commands below are used in Development 💻 and GitHub workflow 🌊.

# Usage: make <COMMANDS> <ARGS> ...

# COMMANDS:
#   build               Build layotto for host platform.
#   build.multiarch     Build layotto for multiple platforms. See option PLATFORMS.
#   image               Build docker images for host arch.
#   image.multiarch     Build docker images for multiple platforms. See option PLATFORMS.
#   push                Push docker images to registry.
#   push.multiarch      Push docker images for multiple platforms to registry.
#   app                 Build app docker images for host arch. [`/docker/app` contains apps dockerfiles]
#   app.multiarch       Build app docker images for multiple platforms. See option PLATFORMS.
#   wasm                Build layotto wasm for linux arm64 platform.
#   wasm.multiarch      Build layotto wasm for multiple platform.
#   wasm.image          Build layotto wasm image for multiple platform.
#   wasm.image.push     Push layotto wasm image for multiple platform.
#   check               Run all go checks of code sources.
#   check.lint          Run go syntax and styling of go sources.
#   check.unit          Run go unit test.
#   check.style         Run go style test.
#   style.coverage      Run coverage analysis.
#   style.deadlink      Run deadlink check test.
#   style.quickstart    Run quickstart check test.
#   integrate.wasm      Run integration test with wasm.
#   integrate.runtime   Run integration test with runtime.
#   format              Format layotto go codes style with gofmt and goimports.
#   clean               Remove all files that are created by building.
#   all                 Run format codes, check codes, build Layotto codes for host platform with one command
#   help                Show this help info.

# ARGS:
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

# ==============================================================================
# ROOT Options:
# ==============================================================================

ROOT_PACKAGE=mosn.io/layotto

# ==============================================================================
# Includes:
# ==============================================================================

include make/common.mk
include make/golang.mk
include make/image.mk
include make/wasm.mk
include make/ci.mk

# ==============================================================================
# Targets:
# ==============================================================================

# ==============================================================================
## build: Build layotto for host platform.
# ==============================================================================
.PHONY: build
build:
	@$(MAKE) go.build

# ==============================================================================
## build.multiarch: Build layotto for multiple platforms. See option PLATFORMS.
# ==============================================================================
.PHONY: build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch

# ==============================================================================
## image: Build docker images for host arch.
# ==============================================================================
.PHONY: image
image:
	@$(MAKE) image.build

# ==============================================================================
## image.multiarch: Build docker images for multiple platforms. See option PLATFORMS.
# ==============================================================================
.PHONY: image.multiarch
image.multiarch:
	@$(MAKE) image.build.multiarch

# ==============================================================================
## push: Push docker images to registry.
# ==============================================================================
.PHONY: push
push:
	@$(MAKE) image.push

# ==============================================================================
## push.multiarch: Push docker images for multiple platforms to registry.
# ==============================================================================
.PHONY: push.multiarch
push.multiarch:
	@$(MAKE) image.push.multiarch

# ==============================================================================
## app: Build app docker images for host arch. [`/docker/app` contains apps dockerfiles]
# ==============================================================================
.PHONY: app
app:
	@$(MAKE) app.image

# ==============================================================================
## app.multiarch: Build app docker images for multiple platforms. See option PLATFORMS.
# ==============================================================================
.PHONY: app.multiarch
app.multiarch:
	@$(MAKE) app.image.multiarch

# ==============================================================================
## wasm: Build layotto wasm for linux arm64 platform.
# ==============================================================================
.PHONY: wasm
wasm:
	@$(MAKE) go.wasm

# ==============================================================================
## wasm.multiarch: Build layotto wasm for multiple platform.
# ==============================================================================
.PHONY: wasm.multiarch
wasm.multiarch:
	@$(MAKE) go.wasm.multiarch

# ==============================================================================
## wasm.image: Build layotto wasm image for multiple platform.
# ==============================================================================
.PHONY: wasm.image
wasm.image:
	@$(MAKE) go.wasm.image

# ==============================================================================
## wasm.image.push: Push layotto wasm image for multiple platform.
# ==============================================================================
.PHONY: wasm.image.push
wasm.image.push:
	@$(MAKE) go.wasm.image.push

# ==============================================================================
## check: Run all go checks of code sources.
# ==============================================================================
.PHONY: check
check: check.style check.unit check.lint

# ==============================================================================
## check.lint: Run go syntax and styling of go sources.
# ==============================================================================
.PHONY: check.lint
check.lint:
	@$(MAKE) go.lint

# ==============================================================================
## check.unit: Run go unit test.
# ==============================================================================
.PHONY: check.unit
check.unit:
	@$(MAKE) go.test

# ==============================================================================
## check.style: Run go style test.
# ==============================================================================
.PHONY: check.style
check.style:
	@$(MAKE) go.style

# ==============================================================================
## style.coverage: Run coverage analysis.
# ==============================================================================
.PHONY: style.coverage
style.coverage:
	@$(MAKE) checker.coverage

# ==============================================================================
## style.deadlink: Run deadlink check test.
# ==============================================================================
.PHONY: style.deadlink
style.deadlink:
	@$(MAKE) checker.deadlink

# ==============================================================================
## style.quickstart: Run quickstart check test.
# ==============================================================================
.PHONY: style.quickstart
style.quickstart:
	@$(MAKE) checker.quickstart

# ==============================================================================
## integrate.wasm: Run integration test with wasm.
# ==============================================================================
.PHONY: integrate.wasm
integrate.wasm:
	@$(MAKE) integration.wasm

# ==============================================================================
## integrate.runtime: Run integration test with runtime.
# ==============================================================================
.PHONY: integrate.runtime
integrate.runtime:
	@$(MAKE) integration.runtime

# ==============================================================================
## format: Format layotto go codes style with gofmt and goimports.
# ==============================================================================
.PHONY: format
format: go.format

# ==============================================================================
## clean: Remove all files that are created by building.
# ==============================================================================
.PHONY: clean
clean:
	@$(MAKE) go.clean

# ==============================================================================
## all: Run format codes, check codes, build Layotto codes for host platform with one command
# ==============================================================================
.PHONY: all
all: clean format check style.quickstart clean

# ==============================================================================
# Usage
# ==============================================================================

define USAGE_OPTIONS

ARGS:
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
# Help
# ==============================================================================

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "Layotto commands 👀: \n\nA fast and efficient cloud native application runtime 🚀."
	@echo -e "Commands below are used in Development 💻 and GitHub workflow 🌊.\n"
	@echo -e "Usage: make <COMMANDS> <ARGS> ...\n\nCOMMANDS:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"