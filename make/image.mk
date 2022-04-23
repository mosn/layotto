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
IMAGE_PLATFORMS ?= linux_amd64 linux_arm64 

ifeq (${IMAGES},)
  $(error Could not determine IMAGES, set ROOT_DIR or run in source dir)
endif

.PHONY: image.daemon.verify
image.daemon.verify:
	$(eval PASS := $(shell $(DOCKER) version | grep -q -E 'Experimental: {1,5}true' && echo 1 || echo 0))
	@if [ $(PASS) -ne 1 ]; then \
		echo "Experimental features of Docker daemon is not enabled. Please add \"experimental\": true in '/etc/docker/daemon.json' and then restart Docker daemon."; \
		exit 1; \
	fi

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
	@cat $(ROOT_DIR)/docker/$(IMAGE)/Dockerfile\
		>$(TMP_DIR)/$(IMAGE)/Dockerfile
	@cp $(OUTPUT_DIR)/$(IMAGE_PLAT)/$(IMAGE) $(TMP_DIR)/$(IMAGE)/
	$(eval BUILD_SUFFIX := $(_DOCKER_BUILD_EXTRA_ARGS) --pull -t $(REGISTRY_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION) $(TMP_DIR)/$(IMAGE))
	$(DOCKER) build --platform $(IMAGE_PLAT) $(BUILD_SUFFIX)

APPS ?= faas integrate
APP_PLATFORMS = linux_amd64 linux_arm64 

.PHONY: app.image
app.image: image.verify  $(addprefix app.image., $(addprefix $(IMAGE_PLAT)., $(APPS)))

.PHONY: app.image.multiarch
app.image.multiarch: image.verify  $(foreach p,$(APP_PLATFORMS),$(addprefix app.image., $(addprefix $(p)., $(APPS))))

.PHONY: app.image.%
app.image.%:
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval APP := $(word 2,$(subst ., ,$*)))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	@echo "===========> Building docker image $(APP) $(VERSION) for $(IMAGE_PLAT)"
	@mkdir -p $(TMP_DIR)/$(APP)
	@cat $(ROOT_DIR)/docker/app/$(APP)/Dockerfile\
		>$(TMP_DIR)/$(APP)/Dockerfile
	$(eval BUILD_SUFFIX := $(_DOCKER_BUILD_EXTRA_ARGS) --pull -t $(REGISTRY_PREFIX)/$(APP)-$(ARCH):$(VERSION) $(TMP_DIR)/$(APP))
	$(DOCKER) build --platform $(IMAGE_PLAT) $(BUILD_SUFFIX)

.PHONY: wasm.image
wasm.image: wasm
	$(eval IMAGE := layotto)
	$(eval IMAGE_PLAT := $(subst _,/,$(WASM_PLATFORM)))
	$(eval ARCH := $(word 2,$(subst _, ,$(WASM_PLATFORM))))
	@echo "===========> Building docker image $(IMAGE) $(VERSION) for $(IMAGE_PLAT)"
	@mkdir -p $(TMP_DIR)/$(IMAGE)
	@cat $(ROOT_DIR)/docker/$(IMAGE)/Dockerfile\
		>$(TMP_DIR)/$(IMAGE)/Dockerfile
	@cp $(OUTPUT_DIR)/$(IMAGE_PLAT)/layotto $(TMP_DIR)/$(IMAGE)/
	$(MAKE) image.daemon.verify
	$(eval BUILD_SUFFIX := $(_DOCKER_BUILD_EXTRA_ARGS) --pull -t $(REGISTRY_PREFIX)/$(IMAGE).wasm.$(ARCH):$(VERSION) $(TMP_DIR)/$(IMAGE))
	$(DOCKER) buildx build --platform $(IMAGE_PLAT) $(BUILD_SUFFIX)

.PHONY: wasm.image.push
wasm.image.push:
	$(eval IMAGE := layotto)
	$(eval ARCH := $(word 2,$(subst _, ,$(WASM_PLATFORM))))
	$(eval IMAGE_PLAT := $(subst _,/,$(WASM_PLATFORM)))
	@echo "===========> Pushing image $(IMAGE) $(VERSION) to $(REGISTRY_PREFIX)"
	$(DOCKER) push $(REGISTRY_PREFIX)/$(IMAGE).wasm.$(ARCH):$(VERSION)

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
integration.wasm: app.image.linux_amd64.faas
	@echo "===========> Integration Test With WASM"
	$(eval ACTION := sh ./wasm_test.sh)
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	$(eval BUILD_IMAGE := $(REGISTRY_PREFIX)/faas-$(ARCH):$(VERSION))
	$(eval WORKDIR := -w /go/src/${PROJECT_NAME} )
	$(eval INTEGRATE_SUFFIX := -v $(ROOT_DIR):/go/src/${PROJECT_NAME} -v ${TEST_DIR}/wasm/wasm_test.sh:/go/src/${PROJECT_NAME}/wasm_test.sh $(WORKDIR))
	$(DOCKER) run --rm $(INTEGRATE_SUFFIX) $(BUILD_IMAGE) $(ACTION)

.PHONY: integration.runtime
integration.runtime: app.image.linux_amd64.integrate
	@echo "===========> Integration Test With Runtime"
	$(eval ACTION := sh ./integrate_test.sh)
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	$(eval BUILD_IMAGE := $(REGISTRY_PREFIX)/integrate-$(ARCH):$(VERSION))
	$(eval WORKDIR := -w /go/src/${PROJECT_NAME} )
	$(eval INTEGRATE_SUFFIX := -v $(ROOT_DIR):/go/src/${PROJECT_NAME} -v ${TEST_DIR}/runtime/integrate_test.sh:/go/src/${PROJECT_NAME}/integrate_test.sh $(WORKDIR))
	$(DOCKER) run --rm $(INTEGRATE_SUFFIX) ${BUILD_IMAGE} $(ACTION)