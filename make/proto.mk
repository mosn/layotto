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

##@ Proto Development

# ==============================================================================
# Public Commands:
# ==============================================================================

.PHONY: proto
proto: ## Generate code and documentation based on the proto files.
proto: proto.gen.all

.PHONY: proto-doc
proto-doc: ## Generate documentation based on the proto files.
proto-doc: proto.gen.doc

.PHONY: proto-code
proto-code: ## Generate code based on the proto files.
proto-code: proto.gen.code

.PHONY: proto-lint
proto-lint: ## Run Protobuffer Linter with Buf Tool
proto-lint: proto.lint


# ==============================================================================
# Private Commands:
# ==============================================================================

.PHONY: proto.gen.doc
proto.gen.doc:
	sh ${SCRIPT_DIR}/generate-doc.sh

.PHONY: proto.gen.init
proto.gen.init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/layotto/protoc-gen-p6@latest

.PHONY: proto.gen.code
proto.gen.code:
	sh ${SCRIPT_DIR}/generate-code.sh
	$(MAKE) format

.PHONY: proto.lint
proto.lint:
ifeq (,$(shell which buf))
	@echo "===========> Installing buf linter"
	@curl -fsSL \
		"https://github.com/bufbuild/buf/releases/download/v1.6.0/buf-$$(uname -s)-$$(uname -m)" \
		-o "$(OUTPUT_DIR)/buf"
	@sudo install -m 0755 $(OUTPUT_DIR)/buf /usr/local/bin/buf
endif
	@echo "===========> Running buf linter"
	buf lint $(ROOT_DIR)

.PHONY: proto.gen.all
proto.gen.all: proto.gen.code proto.gen.doc
