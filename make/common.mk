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

# This is a wrapper to set common variables
#
# All make targets related to common variables are defined in this file.

# ====================================================================================================
# Configure Make itself:
# ====================================================================================================

# Turn off .INTERMEDIATE file removal by marking all files as
# .SECONDARY.  .INTERMEDIATE file removal is a space-saving hack from
# a time when drives were small; on modern computers with plenty of
# storage, it causes nothing but headaches.
#
# https://news.ycombinator.com/item?id=16486331

.SECONDARY:

# ==============================================================================
# Common Variables:
# ==============================================================================

SHELL := /bin/bash

ROOT_PACKAGE=mosn.io/layotto

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

PROJECT_NAME := mosn.io/layotto

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/.. && pwd -P))
endif
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif
ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
$(shell mkdir -p $(TMP_DIR))
endif
ifeq ($(origin DOCS_DIR),undefined)
DOCS_DIR := $(ROOT_DIR)/docs
endif
ifeq ($(origin TEST_WASM_DIR),undefined)
TEST_WASM_DIR := $(ROOT_DIR)/test
endif
ifeq ($(origin TEST_RUNTIME_DIR),undefined)
TEST_RUNTIME_DIR := $(ROOT_DIR)/sdk/go-sdk/test
endif
ifeq ($(origin DEMO_DIR),undefined)
DEMO_DIR := $(ROOT_DIR)/demo
endif
ifeq ($(origin CONFIG_DIR),undefined)
CONFIG_DIR := $(ROOT_DIR)/configs
DEFAULT_CONFIG_FILE := $(CONFIG_DIR)/runtime_config.json
endif
ifeq ($(origin SCRIPT_DIR),undefined)
SCRIPT_DIR := $(ROOT_DIR)/etc/script
endif
ifeq ($(origin SUPERVISOR_DIR),undefined)
SUPERVISOR_DIR := $(ROOT_DIR)/etc/supervisor
endif
ifeq ($(origin DEPLOY_DIR),undefined)
DEPLOY_DIR := $(ROOT_DIR)/deploy
endif
ifeq ($(origin K8S_DIR),undefined)
K8S_DIR := $(DEPLOY_DIR)/k8s
endif

# set the version number. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --abbrev=0 --dirty --always --tags | sed 's/-/./g')
endif
# Check if the tree is dirty.  default to dirty
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

PLATFORMS ?= darwin_amd64 darwin_arm64 linux_amd64 linux_arm64 

# Set a specific PLATFORM
ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
	# Use linux as the default OS when building images
	IMAGE_PLAT := linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif

# ==============================================================================
# Colors: globel colors to share.
# ==============================================================================

NO_COLOR := \033[0m
BOLD_COLOR := \n\033[1m
RED_COLOR := \033[0;31m
GREEN_COLOR := \033[0;32m
YELLOW_COLOR := \033[0;33m
BLUE_COLOR := \033[36m

# ==============================================================================
# Includes:
# ==============================================================================

include make/golang.mk
include make/image.mk
include make/proto.mk
include make/deploy.mk
include make/wasm.mk
include make/ci.mk
include make/help.mk
