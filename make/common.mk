SHELL := /bin/bash

# include the common make file
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
PROJECT_NAME := mosn.io/layotto
BINARY_PREFIX := layotto

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
ifeq ($(origin TEST_DIR),undefined)
TEST_DIR := $(ROOT_DIR)/test
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

COMMA := ,
SPACE :=
SPACE +=
