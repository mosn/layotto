#!/bin/bash
#
# Copyright 2021 Layotto Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Proto buf generation
# https://medium.com/blokur/how-to-implement-a-grpc-client-and-server-in-typescript-fa3ac807855e
PATH_ROOT=$(pwd)
PATH_PROTO_ROOT="${PATH_ROOT}/../../spec/proto/runtime/v1"
PATH_PROTO_OUTPUT="${PATH_ROOT}/proto"

prerequisiteCheckProtobuf() {
    if ! type "protoc" > /dev/null; then
        echo "protoc is not installed, trying to install"
        sudo apt update
        sudo apt install -y protobuf-compiler
        protoc --version

        prerequisiteCheckProtobuf
    else
        echo "protoc ($(protoc --version)) installed"
    fi
}

generateGrpc() {
    PATH_PROTO=$1
    PATH_FILE=$2

    echo "[protoc] Generating RPC for $PATH_PROTO/$PATH_FILE"

    # Tools to be installed by npm (see package.json)
    # npm install grpc-tools --save-dev
    # npm install grpc_tools_node_protoc_ts --save-dev
    PROTOC_GEN_TS_PATH="${PATH_ROOT}/node_modules/.bin/protoc-gen-ts"
    PROTOC_GEN_GRPC_PATH="${PATH_ROOT}/node_modules/.bin/grpc_tools_node_protoc_plugin"

    # commonjs
    protoc \
        --proto_path="${PATH_PROTO}" \
        --plugin="protoc-gen-ts=${PROTOC_GEN_TS_PATH}" \
        --plugin=protoc-gen-grpc=${PROTOC_GEN_GRPC_PATH} \
        --js_out="import_style=commonjs,binary:$PATH_PROTO_OUTPUT" \
        --ts_out="grpc_js:$PATH_PROTO_OUTPUT" \
        --grpc_out="grpc_js:$PATH_PROTO_OUTPUT" \
        "$PATH_PROTO/$PATH_FILE"
}

echo "Checking Dependencies"
prerequisiteCheckProtobuf

echo ""
echo "Removing old Proto Files: ${PATH_PROTO_OUTPUT}"
rm -rf $PATH_PROTO_OUTPUT
mkdir -p $PATH_PROTO_OUTPUT

echo ""
echo "Compiling gRPC files"
generateGrpc $PATH_PROTO_ROOT "runtime.proto"
generateGrpc $PATH_PROTO_ROOT "appcallback.proto"

echo ""
echo "DONE"
