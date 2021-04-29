SHELL = /bin/bash
export GO111MODULE=off

MAJOR_VERSION    = $(shell cat VERSION)
TARGET           = runtime
ARM_TARGET       = runtime.aarch64
PROJECT_NAME     = gitlab.alipay-inc.com/ant-mesh/runtime
CONFIG_FILE     = runtime_config.json
BUILD_IMAGE     = godep-builder
IMAGE_NAME      = runtime
GIT_VERSION     = $(shell git log -1 --pretty=format:%h)
REPOSITORY      = mosnio/${IMAGE_NAME}

build:
	@rm -rf build/bundles/${MAJOR_VERSION}/binary
	CGO_ENABLED=1 go build\
		-ldflags "-B 0x$(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')" \
		-v -o ${TARGET} \
		${PROJECT_NAME}/cmd/layotto
	mkdir -p build/bundles/${MAJOR_VERSION}/binary
	mv ${TARGET} build/bundles/${MAJOR_VERSION}/binary
	@cd build/bundles/${MAJOR_VERSION}/binary

build-local:
	@rm -rf build/bundles/${MAJOR_VERSION}/binary
	GO111MODULE=off CGO_ENABLED=1 go build \
		-ldflags "-B 0x$(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n') -X main.Version=${MAJOR_VERSION}(${GIT_VERSION}) -X ${PROJECT_NAME}/pkg/types.IstioVersion=${ISTIO_VERSION}" \
		-v -o ${TARGET} \
		${PROJECT_NAME}/cmd/layotto
	mkdir -p build/bundles/${MAJOR_VERSION}/binary
	mv ${TARGET} build/bundles/${MAJOR_VERSION}/binary
	@cd build/bundles/${MAJOR_VERSION}/binary && $(shell which md5sum) -b ${TARGET} | cut -d' ' -f1  > ${TARGET}.md5
	cp configs/${CONFIG_FILE} build/bundles/${MAJOR_VERSION}/binary

build-arm64:
	@rm -rf build/bundles/${MAJOR_VERSION}/binary
	GOOS=linux GOARCH=arm64 go build\
		-a -ldflags "-B 0x$(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')" \
		-v -o ${ARM_TARGET} \
		${PROJECT_NAME}/cmd/layotto
	mkdir -p build/bundles/${MAJOR_VERSION}/binary
	mv ${ARM_TARGET} build/bundles/${MAJOR_VERSION}/binary
	@cd build/bundles/${MAJOR_VERSION}/binary && $(shell which md5sum) -b ${ARM_TARGET} | cut -d' ' -f1  > ${ARM_TARGET}.md5

build-image:
	docker build --rm -t ${BUILD_IMAGE} build/contrib/builder/binary
	docker run --rm -v $(shell pwd):/go/src/${PROJECT_NAME} -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} make build-local

image:
	make build-image
	@rm -rf IMAGEBUILD
	cp -r build/contrib/builder/image IMAGEBUILD && cp build/bundles/${MAJOR_VERSION}/binary/${TARGET} IMAGEBUILD && cp -r configs IMAGEBUILD && cp -r etc IMAGEBUILD
	#docker build --no-cache --rm -t ${IMAGE_NAME}:${MAJOR_VERSION}-${GIT_VERSION} IMAGEBUILD
	docker build --rm -t ${IMAGE_NAME}:${MAJOR_VERSION}-${GIT_VERSION} IMAGEBUILD
	docker tag ${IMAGE_NAME}:${MAJOR_VERSION}-${GIT_VERSION} ${REPOSITORY}:${MAJOR_VERSION}-${GIT_VERSION}
	rm -rf IMAGEBUILD

.PHONY: build
