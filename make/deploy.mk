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

# This file contains commands to deploy/undeploy layotto into Kubernetes
# Default namespace is `default`

##@ Kubernetes Development

NAMESPACE := default

# ==============================================================================
# Public Commands:
# ==============================================================================

.PHONY: deploy-k8s
deploy-k8s: ## Install Layotto in Kubernetes.
deploy-k8s: deploy.k8s 

.PHONY: undeploy-k8s
undeploy-k8s: ## Uninstall Layotto in Kubernetes.
undeploy-k8s: undeploy.k8s 

# ==============================================================================
# Private Commands:
# ==============================================================================

.PHONY: deploy.k8s
deploy.k8s: deploy.k8s.standalone

.PHONY: deploy.k8s.standalone
deploy.k8s.standalone: 
	@echo "===========> Deploy Layotto to Kubernetes in namespace ${NAMESPACE} in standalone mode"
	@kubectl apply -f $(K8S_DIR)/standalone/default_quickstart.yaml -n ${NAMESPACE}

.PHONY: undeploy.k8s
undeploy.k8s: undeploy.k8s.standalone

.PHONY: undeploy.k8s.standalone
undeploy.k8s.standalone:
	@echo "===========> Clean Layotto to Kubernetes in namespace ${NAMESPACE} in standalone mode"
	@kubectl delete -f $(K8S_DIR)/standalone/default_quickstart.yaml -n ${NAMESPACE}
