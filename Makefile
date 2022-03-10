SHELL = /bin/bash
export GO111MODULE=on

MAJOR_VERSION    = $(shell cat VERSION)
TARGET           = layotto
ARM_TARGET       = layotto.aarch64
PROJECT_NAME     = mosn.io/layotto
CONFIG_FILE      = runtime_config.json
BUILD_IMAGE      = godep-builder
IMAGE_NAME       = layotto
GIT_VERSION      = $(shell git log -1 --pretty=format:%h)
REPOSITORY       = layotto/${IMAGE_NAME}

SCRIPT_DIR       = $(shell pwd)/etc/script
IMAGE_BUILD_DIR  = IMAGEBUILD

build-local:
	@rm -rf build/bundles/${MAJOR_VERSION}/binary
	CGO_ENABLED=1 go build \
		-ldflags "-B 0x$(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n') -X main.Version=${MAJOR_VERSION}(${GIT_VERSION}) -X ${PROJECT_NAME}/pkg/types.IstioVersion=${ISTIO_VERSION}" \
		-v -o ${TARGET} \
		${PROJECT_NAME}/cmd/layotto
	mkdir -p build/bundles/${MAJOR_VERSION}/binary
	mv ${TARGET} build/bundles/${MAJOR_VERSION}/binary
	@cd build/bundles/${MAJOR_VERSION}/binary && $(shell which md5sum) -b ${TARGET} | cut -d' ' -f1  > ${TARGET}.md5
	cp configs/${CONFIG_FILE} build/bundles/${MAJOR_VERSION}/binary
	@cd build/bundles/${MAJOR_VERSION}/binary

build-arm64:
	@rm -rf build/bundles/${MAJOR_VERSION}/binary
	GOOS=linux GOARCH=arm64 go build\
		-a -ldflags "-B 0x$(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')" \
		-v -o ${ARM_TARGET} \
		${PROJECT_NAME}/cmd/layotto
	mkdir -p build/bundles/${MAJOR_VERSION}/binary
	mv ${ARM_TARGET} build/bundles/${MAJOR_VERSION}/binary
	@cd build/bundles/${MAJOR_VERSION}/binary && $(shell which md5sum) -b ${ARM_TARGET} | cut -d' ' -f1  > ${ARM_TARGET}.md5

image: build-local
	@rm -rf ${IMAGE_BUILD_DIR}
	cp -r build/contrib/builder/image ${IMAGE_BUILD_DIR} && cp build/bundles/${MAJOR_VERSION}/binary/${TARGET} ${IMAGE_BUILD_DIR} && cp -r configs ${IMAGE_BUILD_DIR} && cp -r etc ${IMAGE_BUILD_DIR}
	docker build --rm -t ${REPOSITORY}:${MAJOR_VERSION}-${GIT_VERSION} ${IMAGE_BUILD_DIR}
	rm -rf ${IMAGE_BUILD_DIR}

wasm-integrate-ci:
	docker build --rm -t ${BUILD_IMAGE} build/contrib/builder/image/faas
	docker run --rm -v $(shell pwd):/go/src/${PROJECT_NAME} -v $(shell pwd)/test/test.sh:/go/src/${PROJECT_NAME}/test.sh -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} sh ./test.sh

coverage:
	sh ${SCRIPT_DIR}/report.sh

build-linux-wasm-layotto:
	go build -tags wasmer -o layotto $(shell pwd)/cmd/layotto
