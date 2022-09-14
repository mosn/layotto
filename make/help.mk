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

# ==============================================================================
# Usage
# ==============================================================================

## help: Show this help info.
.PHONY: help
help:
	@echo -e "$(BOLD_COLOR)Layotto$(NO_COLOR) is a fast and efficient cloud native application runtime."
	@echo -e "$(BOLD_COLOR)Usage:$(NO_COLOR)\n  make \033[36m<Target>\033[0m \033[36m<Option>\033[0m\n$(BOLD_COLOR)Targets:$(NO_COLOR)"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo -e "$$USAGE_OPTIONS"

define USAGE_OPTIONS

$(BOLD_COLOR)Options$(NO_COLOR):

  $(BLUE_COLOR)BINS$(NO_COLOR)         The binaries to build. Default is all of cmd.
               This option is available when using: make build/multiarch
               Examples:
               $(GREEN_COLOR)* make multiarch BINS="layotto"$(NO_COLOR)
               $(GREEN_COLOR)* make build BINS="layotto_multiple_api layotto"$(NO_COLOR)
  $(BLUE_COLOR)IMAGES$(NO_COLOR)       Backend images to make. Default is all of cmds.
               This option is available when using: make image/image-multiarch/push/push-multiarch
               Examples: 
               $(GREEN_COLOR)* make image IMAGES="layotto"$(NO_COLOR)
               $(GREEN_COLOR)* make image-multiarch IMAGES="layotto"$(NO_COLOR)
               $(GREEN_COLOR)* make push IMAGES="layotto_multiple_api"$(NO_COLOR)
               $(GREEN_COLOR)* make push-multiarch IMAGES="layotto_multiple_api"$(NO_COLOR)
  $(BLUE_COLOR)NAMESPACE$(NO_COLOR)    The namepace to deploy. Default is `default`.
               This option is available when using: make deploy-k8s/undeploy-k8s
               Examples: 
               $(GREEN_COLOR)* make deploy-k8s NAMESPACE="layotto"$(NO_COLOR)
               $(GREEN_COLOR)* make undeploy-k8s NAMESPACE="default"$(NO_COLOR)
  $(BLUE_COLOR)VERSION$(NO_COLOR)    The image tag version to build. Default is the latest release tag.
               This option is available when using: make image/image-multiarch/push/push-multiarch
               Examples: 
               $(GREEN_COLOR)* make image VERSION="latest"$(NO_COLOR)
               $(GREEN_COLOR)* make image-multiarch VERSION="v1.0.0"$(NO_COLOR)
               $(GREEN_COLOR)* make push-multiarch VERSION="v2.0.0"$(NO_COLOR)
  $(BLUE_COLOR)REGISTRY_PREFIX$(NO_COLOR)    The docker image registry repo name to push. Default is `layotto`.
               This option is available when using: make push/push-multiarch
               Examples: 
               $(GREEN_COLOR)* make push IMAGES="layotto" REGISTRY_PREFIX="mosn"$(NO_COLOR)
               $(GREEN_COLOR)* make push IMAGES="layotto_multiple_api" REGISTRY_PREFIX="mosn"$(NO_COLOR)
               Supported Platforms: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64
  $(BLUE_COLOR)PLATFORMS$(NO_COLOR)    The multiple platforms to build. Default is linux_amd64 and linux_arm64.
               This option is available when using: make multiarch
               Examples: 
               $(GREEN_COLOR)* make multiarch BINS="layotto" PLATFORMS="linux_amd64 linux_arm64"$(NO_COLOR)
               Supported Platforms: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64
endef
export USAGE_OPTIONS
