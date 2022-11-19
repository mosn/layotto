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

WASM_PLATFORM ?= linux_amd64
WASM_PLATFORMS ?= linux_amd64
WASM_BUILD ?= faas

##@ WebAssembly Development

# ==============================================================================
# Public Commands:
# ==============================================================================

.PHONY: wasm-build
wasm-build: ## Build layotto wasm for linux arm64 platform.
wasm-build: go.wasm

.PHONY: wasm-image
wasm-image: ## Build layotto wasm image for multiple platform.
wasm-image: go.wasm.image

.PHONY: wasm-push
wasm-push: ## Push layotto wasm image for multiple platform.
wasm-push: go.wasm.image.push

# ==============================================================================
# Private Commands:
# ==============================================================================

.PHONY: go.wasm
go.wasm:  $(addprefix go.wasm., $(addprefix $(WASM_PLATFORM)., $(WASM_BUILD)))

.PHONY: go.wasm.multiarch
go.wasm.multiarch:  $(foreach p,$(WASM_PLATFORMS),$(addprefix go.wasm., $(addprefix $(p)., $(WASM_BUILD))))

.PHONY: go.wasm.%
go.wasm.%:
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
	$(eval ACTION := $(GO) build -o $(OUTPUT_PATH) -tags wasmcomm,wasmtime -ldflags "$(GO_LDFLAGS)" $(ROOT_PACKAGE)/cmd/layotto)
	$(DOCKER) run --rm -v $(ROOT_DIR):/go/src/${PROJECT_NAME} -e GOOS=$(OS) -e GOARCH=$(ARCH) -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} ${ACTION}

.PHONY: go.wasm.image
go.wasm.image: go.wasm
	$(eval IMAGE := layotto)
	$(eval IMAGE_PLAT := $(subst _,/,$(WASM_PLATFORM)))
	$(eval ARCH := $(word 2,$(subst _, ,$(WASM_PLATFORM))))
	@echo "===========> Building docker image $(IMAGE) $(VERSION) for $(IMAGE_PLAT)"
	@mkdir -p $(TMP_DIR)/$(IMAGE)
	@cat $(ROOT_DIR)/docker/$(IMAGE)/Dockerfile\
		>$(TMP_DIR)/$(IMAGE)/Dockerfile
	@cp $(OUTPUT_DIR)/$(IMAGE_PLAT)/layotto $(TMP_DIR)/$(IMAGE)/
	$(eval BUILD_SUFFIX := $(_DOCKER_BUILD_EXTRA_ARGS) --pull -t $(REGISTRY_PREFIX)/$(IMAGE).wasm.$(ARCH):$(VERSION) $(TMP_DIR)/$(IMAGE))
	$(DOCKER) buildx build --platform $(IMAGE_PLAT) $(BUILD_SUFFIX)

.PHONY: go.wasm.image.push
go.wasm.image.push:
	$(eval IMAGE := layotto)
	$(eval ARCH := $(word 2,$(subst _, ,$(WASM_PLATFORM))))
	$(eval IMAGE_PLAT := $(subst _,/,$(WASM_PLATFORM)))
	@echo "===========> Pushing image $(IMAGE) $(VERSION) to $(REGISTRY_PREFIX)"
	$(DOCKER) push $(REGISTRY_PREFIX)/$(IMAGE).wasm.$(ARCH):$(VERSION)
