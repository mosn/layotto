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

.PHONY: proto.gen.doc
proto.gen.doc:
	$(DOCKER) run --rm \
    -v  $(ROOT_DIR)/docs/en/api_reference:/out \
    -v  $(ROOT_DIR)/spec/proto/runtime/v1:/protos \
    pseudomuto/protoc-gen-doc  --doc_opt=/protos/template.tmpl,runtime_v1.md runtime.proto
	$(DOCKER) run --rm \
    -v  $(ROOT_DIR)/docs/en/api_reference:/out \
    -v  $(ROOT_DIR)/spec/proto/runtime/v1:/protos \
    pseudomuto/protoc-gen-doc  --doc_opt=/protos/template.tmpl,appcallback_v1.md appcallback.proto

.PHONY: proto.gen.init
proto.gen.init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.PHONY: proto.gen.code
proto.gen.code:
	$(DOCKER) build -t layotto/protoc $(ROOT_DIR)/docker/proto && \
	$(DOCKER) run --rm \
		-v  $(ROOT_DIR)/spec/proto/runtime/v1:/api/proto \
		layotto/protoc

.PHONY: proto.comments
proto.comments:
	curl -fsSL \
		"https://github.com/bufbuild/buf/releases/download/v1.6.0/buf-$$(uname -s)-$$(uname -m)" \
		-o "$(OUTPUT_DIR)/buf"
	buf lint $(ROOT_DIR)