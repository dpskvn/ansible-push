SHORTNAME				:= sample-connector
IMAGE_NAME				 = tls-protect-${SHORTNAME}
IMAGE_TAG				:= latest
BUILDX_BUILDER_NAME    	?= default
MODULE_NAME     		:= ${SHORTNAME}
BUILDX_EXTRA_ARGS  		:= --build-arg=MODULE_NAME=${MODULE_NAME}
CONTAINER_REPOSITORY 	?= <CONTAINER_REPOSITORY>

.DEFAULT: clean build

.PHONY: prepare_build
prepare_build:
	@mkdir -p output

.PHONY: clean
clean:
	@rm -rf output/

.PHONY: build
build: BUILDX_OUTPUT:=-o type=local,dest=output
build: BUILDX_TARGET:=--target reports
build: prepare_build buildx

.PHONY: image
image: BUILDX_OUTPUT:=--output type=image,name=${CONTAINER_REPOSITORY}/${IMAGE_NAME}:${IMAGE_TAG}
image: BUILDX_TARGET:=--target image
image: BUILDX_EXTRA_ARGS:=--metadata-file=buildx-digest.json
image: buildx

.PHONY: push
push: BUILDX_OUTPUT:=--output type=image,name=${CONTAINER_REPOSITORY}/${IMAGE_NAME}:${IMAGE_TAG},push=true
push: BUILDX_TARGET:=--target image
push: BUILDX_EXTRA_ARGS:=--metadata-file=buildx-digest.json
push: buildx

.PHONY: buildx
buildx:
	docker buildx build ${BUILDX_OUTPUT} ${BUILDX_TARGET} ${BUILDX_EXTRA_ARGS} -f build/Dockerfile --builder ${BUILDX_BUILDER_NAME} .

.PHONY: lint
lint:
	golangci-lint run --config build/golangci.yaml --out-format colored-line-number --issues-exit-code 1 ./...

.PHONY: test
test:
	go test -cover ./...
