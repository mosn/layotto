GO := go
GO_LDFLAGS += -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \

ifeq ($(ROOT_PACKAGE),)
	$(error the variable ROOT_PACKAGE must be set prior to including golang.mk)
endif

GOPATH := $(shell go env GOPATH)
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

WASM_PLATFORM ?= linux_amd64
WASM_PLATFORMS ?= linux_amd64
WASM_BUILD ?= faas

.PHONY: wasm
wasm:  $(addprefix wasm., $(addprefix $(WASM_PLATFORM)., $(WASM_BUILD)))

.PHONY: wasm.multiarch
wasm.multiarch:  $(foreach p,$(WASM_PLATFORMS),$(addprefix wasm., $(addprefix $(p)., $(WASM_BUILD))))

.PHONY: wasm.%
wasm.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	$(eval BUILD_IMAGE := $(REGISTRY_PREFIX)/faas-$(ARCH):$(VERSION))
	@mkdir -p $(TMP_DIR)/$(COMMAND)
	@mkdir -p $(OUTPUT_DIR)/$(OS)/$(ARCH)
	@cat $(ROOT_DIR)/docker/app/$(COMMAND)/Dockerfile\
		>$(TMP_DIR)/$(COMMAND)/Dockerfile
	$(eval DOCKER_FILE := $(TMP_DIR)/$(COMMAND)/Dockerfile)
	@echo "===========> Building wasm base image in $(VERSION) for $(OS) $(ARCH)"
	$(DOCKER) build -f ${DOCKER_FILE} -t  ${BUILD_IMAGE} .
	@echo "===========> Building binary wasm in $(VERSION) for $(OS) $(ARCH)"
	$(eval OUTPUT_PATH := ./_output/$(OS)/$(ARCH)/layotto)
	$(eval ACTION := $(GO) build -o $(OUTPUT_PATH) -tags wasmer -ldflags "$(GO_LDFLAGS)" $(ROOT_PACKAGE)/cmd/layotto)
	$(DOCKER) run --rm -v $(ROOT_DIR):/go/src/${PROJECT_NAME} -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} ${ACTION}

.PHONY: go.clean
go.clean:
	@echo "===========> Cleaning all build output"
	@rm -rf $(OUTPUT_DIR)
	@rm -rf $(ROOT_DIR)/cover.out

.PHONY: go.lint.verify
go.lint.verify:
ifeq (,$(shell which golangci-lint))
	@echo "===========> Installing golangci lint"
	@GO111MODULE=off $(GO) get -u github.com/golangci/golangci-lint/cmd/golangci-lint
endif

.PHONY: go.lint
go.lint: go.lint.verify
	@echo "===========> Run golangci to lint source codes"
	@golangci-lint run $(ROOT_DIR)/...

.PHONY: go.test.verify
go.test.verify:  
ifeq ($(shell which go-junit-report), )
	@echo "===========> Installing go-junit-report"
	@GO111MODULE=off $(GO) get -u github.com/jstemmer/go-junit-report
endif

.PHONY: go.test
go.test: go.test.verify
	@echo "===========> Run unit test"
	$(GO) test -count=1 -timeout=10m -short -v `go list ./...|grep -v mosn.io/layotto/test` 2>&1 | tee >(go-junit-report --set-exit-code >$(OUTPUT_DIR)/report.xml)

.PHONY: go.style
go.style:  
	@echo "===========> Running go style check"
	$(GO) fmt ./... && git status && [[ -z `git status -s` ]]

.PHONY: checker.deadlink
checker.deadlink:
	@echo "===========> Checking Dead Links"
	sh ${SCRIPT_DIR}/check-dead-link.sh

.PHONY: checker.quickstart
checker.quickstart:
	@echo "===========> Checking QuickStart Doc"
	curl -o ${ROOT_DIR}/mdsh.sh https://raw.githubusercontent.com/seeflood/mdsh/master/bin/mdsh
	mv ${ROOT_DIR}/mdsh.sh ${SCRIPT_DIR}
	chmod +x  ${SCRIPT_DIR}/mdsh.sh
	sh ${SCRIPT_DIR}/test-quickstart.sh

.PHONY: checker.coverage
checker.coverage:
	@echo "===========> Coverage Analysis"
	sh ${SCRIPT_DIR}/report.sh