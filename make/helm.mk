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

##@ Helm Development

OCI_REGISTRY ?= oci://docker.io/layotto
CHART_NAME ?= injector-helm
CHART_VERSION ?= ${VERSION}
APP_VERSION ?= ${VERSION}

# ==============================================================================
# Public Commands:
# ==============================================================================

.PHONY: helm-package
helm-package: ## Package layotto injector helm chart.
helm-package: helm.package

.PHONY: helm-push
helm-push:	## Push layotto injector helm chart to OCI registry.
helm-push: helm.push

.PHONY: helm-install
helm-install: ## Install layotto injector helm chart from OCI registry.
helm-install: helm.install

.PHONY: helm-uninstall
helm-uninstall: ## Uninstall layotto injector helm chart.
helm-uninstall: helm.uninstall

# ==============================================================================
# Private Commands:
# ==============================================================================

.PHONY: helm.package
helm.package:
	@echo "===========> Package layotto injector helm chart $(CHART_VERSION)"
	helm package $(CHART_DIR)/${CHART_NAME} --app-version ${APP_VERSION} --version ${CHART_VERSION} --destination ${OUTPUT_DIR}/charts/

.PHONY: helm.push
helm.push:
	@echo "===========> Push layotto injector helm chart $(CHART_VERSION) to OCI registry"
	helm push ${OUTPUT_DIR}/charts/${CHART_NAME}-${CHART_VERSION}.tgz ${OCI_REGISTRY}

.PHONY: helm.install
helm.install:
	@echo "===========> Install layotto injector helm chart $(CHART_VERSION) from OCI registry"
	helm install injector ${OCI_REGISTRY}/${CHART_NAME} --version ${CHART_VERSION} -n layotto-system --create-namespace

.PHONY: helm.uninstall
helm.uninstall:
	@echo "===========> Uninstall layotto injector helm chart $(CHART_VERSION)"
	helm uninstall injector -n layotto-system



