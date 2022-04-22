DOCKER := docker
DOCKER_SUPPORTED_API_VERSION ?= 1.32

REGISTRY_PREFIX ?= layotto

EXTRA_ARGS ?=
_DOCKER_BUILD_EXTRA_ARGS :=

ifdef HTTP_PROXY
_DOCKER_BUILD_EXTRA_ARGS += --build-arg HTTP_PROXY=${HTTP_PROXY}
endif

ifneq ($(EXTRA_ARGS), )
_DOCKER_BUILD_EXTRA_ARGS += $(EXTRA_ARGS)
endif

# Determine image files by looking into build/docker/*/Dockerfile
IMAGES_DIR ?= $(wildcard ${ROOT_DIR}/docker/*)
# Determine images names by stripping out the dir names
IMAGES ?= layotto
IMAGE_PLATFORMS ?= darwin_amd64 darwin_arm64 linux_amd64 linux_arm64 

ifeq (${IMAGES},)
  $(error Could not determine IMAGES, set ROOT_DIR or run in source dir)
endif

.PHONY: image.verify
image.verify:
	$(eval API_VERSION := $(shell $(DOCKER) version | grep -E 'API version: {1,6}[0-9]' | head -n1 | awk '{print $$3} END { if (NR==0) print 0}' ))
	$(eval PASS := $(shell echo "$(API_VERSION) > $(DOCKER_SUPPORTED_API_VERSION)" | bc))
	@if [ $(PASS) -ne 1 ]; then \
		$(DOCKER) -v ;\
		echo "Unsupported docker version. Docker API version should be greater than $(DOCKER_SUPPORTED_API_VERSION)"; \
		exit 1; \
	fi

.PHONY: image.build
image.build: image.verify  $(addprefix image.build., $(addprefix $(IMAGE_PLAT)., $(IMAGES)))

.PHONY: image.build.multiarch
image.build.multiarch: image.verify  $(foreach p,$(IMAGE_PLATFORMS),$(addprefix image.build., $(addprefix $(p)., $(IMAGES))))

.PHONY: image.build.%
image.build.%: go.build.%
	$(eval IMAGE := $(COMMAND))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	@echo "===========> Building docker image $(IMAGE) $(VERSION) for $(IMAGE_PLAT)"
	@mkdir -p $(TMP_DIR)/$(IMAGE)
	@mkdir -p $(TMP_DIR)/$(IMAGE)/etc
	@mkdir -p $(TMP_DIR)/$(IMAGE)/configs
	@cat $(ROOT_DIR)/docker/$(IMAGE)/Dockerfile\
		>$(TMP_DIR)/$(IMAGE)/Dockerfile
	@cp -r $(SCRIPT_DIR) $(TMP_DIR)/$(IMAGE)/etc
	@cp -r $(SUPERVISOR_DIR) $(TMP_DIR)/$(IMAGE)/etc
	@cp $(DEFAULT_CONFIG_FILE) $(TMP_DIR)/$(IMAGE)/configs
	@cp $(OUTPUT_DIR)/$(IMAGE_PLAT)/$(IMAGE) $(TMP_DIR)/$(IMAGE)/
	$(eval BUILD_SUFFIX := $(_DOCKER_BUILD_EXTRA_ARGS) --pull -t $(REGISTRY_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION) $(TMP_DIR)/$(IMAGE))
	$(DOCKER) build --platform $(IMAGE_PLAT) $(BUILD_SUFFIX)
	@rm -rf $(TMP_DIR)/$(IMAGE)

APPS ?= faas integrate
APP_PLATFORMS = linux_amd64 linux_arm64 

.PHONY: app.image
app.image: image.verify  $(addprefix app.image., $(addprefix $(IMAGE_PLAT)., $(APPS)))

.PHONY: app.image.multiarch
app.image.multiarch: image.verify  $(foreach p,$(APP_PLATFORMS),$(addprefix app.image., $(addprefix $(p)., $(APPS))))

.PHONY: app.image.%
app.image.%:
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval APP := $(COMMAND))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	@echo "===========> Building docker image $(APP) $(VERSION) for $(IMAGE_PLAT)"
	@mkdir -p $(TMP_DIR)/$(APP)
	@cat $(ROOT_DIR)/docker/app/$(APP)/Dockerfile\
		>$(TMP_DIR)/$(APP)/Dockerfile
	$(eval BUILD_SUFFIX := $(_DOCKER_BUILD_EXTRA_ARGS) --pull -t $(REGISTRY_PREFIX)/$(APP)-$(ARCH):$(VERSION) $(TMP_DIR)/$(APP))
	$(DOCKER) build --platform $(IMAGE_PLAT) $(BUILD_SUFFIX)
	@rm -rf $(TMP_DIR)/$(APP)

.PHONY: image.push
image.push: image.verify $(addprefix image.push., $(addprefix $(IMAGE_PLAT)., $(IMAGES)))

.PHONY: image.push.multiarch
image.push.multiarch: image.verify  $(foreach p,$(IMAGE_PLATFORMS),$(addprefix image.push., $(addprefix $(p)., $(IMAGES)))) 

.PHONY: image.push.%
image.push.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval IMAGE := $(COMMAND))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	@echo "===========> Pushing image $(IMAGE) $(VERSION) to $(REGISTRY_PREFIX)"
	$(DOCKER) push $(REGISTRY_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION)

.PHONY: integration.wasm
integration.wasm: app.image.linux_arm64.faas
	@echo "===========> Integration Test With WASM"
	$(eval ACTION := sh ./wasm_test.sh)
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	$(eval BUILD_IMAGE := $(REGISTRY_PREFIX)/faas-$(ARCH):$(VERSION))
	$(eval WORKDIR := -w /go/src/${PROJECT_NAME} )
	$(eval INTEGRATE_SUFFIX := -v $(ROOT_DIR):/go/src/${PROJECT_NAME} -v ${TEST_DIR}/wasm/wasm_test.sh:/go/src/${PROJECT_NAME}/wasm_test.sh $(WORKDIR))
	$(DOCKER) run --rm $(INTEGRATE_SUFFIX) $(BUILD_IMAGE) $(ACTION)

.PHONY: integration.runtime
integration.runtime: app.image.linux_arm64.integrate
	@echo "===========> Integration Test With Runtime"
	$(eval ACTION := sh ./integrate_test.sh)
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	$(eval BUILD_IMAGE := $(REGISTRY_PREFIX)/integrate-$(ARCH):$(VERSION))
	$(eval WORKDIR := -w /go/src/${PROJECT_NAME} )
	$(eval INTEGRATE_SUFFIX := -v $(ROOT_DIR):/go/src/${PROJECT_NAME} -v ${TEST_DIR}/runtime/integrate_test.sh:/go/src/${PROJECT_NAME}/integrate_test.sh $(WORKDIR))
	$(DOCKER) run --rm $(INTEGRATE_SUFFIX) ${BUILD_IMAGE} $(ACTION)