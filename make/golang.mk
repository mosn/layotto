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

GO := go
GO_FMT := gofmt
GO_IMPORTS := goimports

GO_MODULE := mosn.io/layotto
VERSION_PACKAGE := main
BINARY_PREFIX := layotto

GO_LDFLAGS += -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	# -X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	# -X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	# -X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \

ifeq ($(ROOT_PACKAGE),)
	$(error the variable ROOT_PACKAGE must be set prior to including golang.mk)
endif

GOPATH := $(shell go env GOPATH | cut -d ':' -f 1)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif

COMMANDS ?= $(filter-out %.md, $(wildcard ${ROOT_DIR}/cmd/*))
BINS ?= $(foreach cmd,${COMMANDS},$(notdir ${cmd}))

ifeq (${COMMANDS},)
  $(error Could not determine COMMANDS, set ROOT_DIR or run in source dir)
endif
ifeq (${BINS},)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif

UNITTEST_OUT := $(OUTPUT_DIR)/unittest.out

##@ Golang Development

# ==============================================================================
# Public Commands:
# ==============================================================================

.PHONY: build
build: ## Build layotto for host platform.
build: go.build

.PHONY: multiarch
multiarch: ## Build layotto for multiple platforms.
multiarch: go.build.multiarch

.PHONY: clean
clean: ## clean all unused generated files.
clean: go.clean

.PHONY: lint
lint: ## Run go syntax and styling of go sources.
lint: go.lint

.PHONY: tests
test: ## Run golang unit test in target paths.
test: go.test

.PHONY: workspace
workspace: ## check if workspace is clean and committed.
workspace: go.style

.PHONY: format
format: ## Format codes style with gofmt and goimports.
format: go.format

# ==============================================================================
# Private Commands:
# ==============================================================================


.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)"
	@mkdir -p $(OUTPUT_DIR)/$(OS)/$(ARCH)
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build -o $(OUTPUT_DIR)/$(OS)/$(ARCH)/$(COMMAND) -ldflags "$(GO_LDFLAGS)" $(ROOT_PACKAGE)/cmd/$(COMMAND)

.PHONY: go.build
go.build:  $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS)))

.PHONY: go.build.multiarch
go.build.multiarch:  $(foreach p,$(PLATFORMS),$(addprefix go.build., $(addprefix $(p)., $(BINS))))

.PHONY: go.clean
go.clean:
	@echo "===========> Cleaning all build output"
	@rm -rf $(OUTPUT_DIR)
	@rm -rf $(ROOT_DIR)/cover.out
	@rm -f cmd/layotto/layotto
	@rm -f cmd/layotto/nohup.out
	@rm -f cmd/layotto_multiple_api/layotto
	@rm -f cmd/layotto_multiple_api/nohup.out
	@rm -rf default.etcd/
	@rm -f demo/configuration/common/client
	@rm -f demo/file/client
	@rm -f demo/flowcontrol/client
	@rm -f demo/lock/redis/client
	@rm -f demo/pubsub/redis/client/publisher
	@rm -f demo/pubsub/redis/server/nohup.out
	@rm -f demo/pubsub/redis/server/subscriber
	@rm -f demo/sequencer/common/client
	@rm -f demo/state/common/client
	@rm -f etc/script/mdx
	@rm -f etcd
	@rm -f layotto_wasmtime
	@rm -f nohup.out

.PHONY: go.lint.verify
go.lint.verify:
ifeq (,$(shell which golangci-lint))
	@echo "===========> Installing golangci lint"
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(${GO} env GOPATH)/bin
endif

.PHONY: go.lint
go.lint: go.lint.verify
	@echo "===========> Run golangci to lint source codes"
	@golangci-lint run -v

.PHONY: go.test
go.unittest:
	@echo "===========> Run unit test in diagnostics" > $(UNITTEST_OUT) && \
	$(GO) test -count=1 -timeout=10m -short -v `go list ./diagnostics/...` >> $(UNITTEST_OUT) && \
	echo "===========> Run unit test in sdk/go-sdk" >> $(UNITTEST_OUT) && \
	cd sdk/go-sdk && $(GO) test -count=1 -timeout=10m -short -v `go list ./... | grep -v runtime` >> $(UNITTEST_OUT) && \
	echo "===========> Run unit test in components" >> $(UNITTEST_OUT) && \
	cd ../../components && $(GO) test -count=1 -timeout=10m -short -v `go list ./...` >> $(UNITTEST_OUT) && \
	echo "===========> Run unit test in pkg" >> $(UNITTEST_OUT) && \
	cd .. && $(GO) test -count=1 -timeout=10m -short -v `go list ./pkg/...` >> $(UNITTEST_OUT) || true

go.test: go.unittest
	@cat $(UNITTEST_OUT)
	@if grep -q "FAIL" $(UNITTEST_OUT); then \
    		grep "FAIL" $(UNITTEST_OUT); \
    		rm -f $(UNITTEST_OUT); \
    		exit 1; \
	fi;
	@rm -f $(UNITTEST_OUT)

.PHONY: go.style
go.style:  
	@echo "===========> Running go style check"
	@$(MAKE) format && git status && [[ -z `git status -s` ]] || echo -e "\n${RED}Error: there are uncommitted changes after formatting all the code. ${GREEN}\nHow to fix it:${NO_COLOR} please 'make format' and then use git to commit all those changed files. "

.PHONY: go.format.verify
go.format.verify:  
ifeq ($(shell which goimports), )
	@echo "===========> Installing missing goimports"
	@mkdir -p $(GOPATH)/src/github.com/golang
	@mkdir -p $(GOPATH)/src/golang.org/x
ifeq ($(shell if [ -d $(GOPATH)/src/github.com/golang/tools ]; then echo "exist"; else echo ""; fi;), )
	@git clone https://github.com/golang/tools.git $(GOPATH)/src/github.com/golang/tools
endif 
ifeq ($(shell if [ -d $(GOPATH)/src/golang.org/x/tools ]; then echo "exist"; else echo ""; fi;), )
	@ln -s $(GOPATH)/src/github.com/golang/tools $(GOPATH)/src/golang.org/x/tools
endif

ifeq ($(shell if [ -d $(GOPATH)/src/github.com/golang/mod ]; then echo "exist"; else echo ""; fi;), )
	@git clone https://github.com/golang/mod.git $(GOPATH)/src/github.com/golang/mod
endif 
ifeq ($(shell if [ -d $(GOPATH)/src/golang.org/x/mod ]; then echo "exist"; else echo ""; fi;), )
	@ln -s $(GOPATH)/src/github.com/golang/mod $(GOPATH)/src/golang.org/x/mod
endif

ifeq ($(shell if [ -d $(GOPATH)/src/github.com/golang/sys ]; then echo "exist"; else echo ""; fi;), )
	@git clone https://github.com/golang/sys.git $(GOPATH)/src/github.com/golang/sys
endif 
ifeq ($(shell if [ -d $(GOPATH)/src/golang.org/x/sys ]; then echo "exist"; else echo ""; fi;), )
	@ln -s $(GOPATH)/src/github.com/golang/sys $(GOPATH)/src/golang.org/x/sys
endif
	@GO111MODULE=off $(GO) build $(GOPATH)/src/golang.org/x/tools/cmd/goimports
	@GO111MODULE=off $(GO) install $(GOPATH)/src/golang.org/x/tools/cmd/goimports
endif

.PHONY: go.format
go.format: go.format.verify
	@echo "===========> Running go codes format"
	$(GO_FMT) -s -w .
	$(GOPATH)/bin/$(GO_IMPORTS) -w -local $(GO_MODULE) .
	$(GO) mod tidy
	cd components && $(GO) mod tidy
	cd demo && $(GO) mod tidy
	cd sdk/go-sdk && $(GO) mod tidy
	cd spec && $(GO) mod tidy
