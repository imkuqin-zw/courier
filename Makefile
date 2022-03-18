# New Makefile for multi-architecture
.PHONY: all

IMAGE_REPOSITORY = ccr.ccs.tencentyun.com/courier
VERSION ?= $(shell date +'%Y%m%d%H%M%S')
VERSION := $(VERSION)
TARGETS = $(shell ls cmd)


ALL_IMAGES=$(TARGETS:=.image)
ALL_PUSH=$(TARGETS:=.push)

all: $(ALL_IMAGES:=.amd64) $(ALL_IMAGES:=.arm64)
all.push: $(ALL_PUSH:=.amd64) $(ALL_PUSH:=.arm64)

all.amd64: $(ALL_IMAGES:=.amd64)
all.arm64: $(ALL_IMAGES:=.arm64)
allpush.amd64: $(ALL_PUSH:=.amd64)
allpush.arm64: $(ALL_PUSH:=.arm64)

%.push.amd64: MAKE_IMAGE ?= ${IMAGE_REPOSITORY}/$*:${VERSION}-amd64
%.push.amd64: %.image.amd64
	@docker push ${MAKE_IMAGE}

%.image.amd64: MAKE_IMAGE ?= ${IMAGE_REPOSITORY}/$*:${VERSION}-amd64
%.image.amd64:
	@mkdir -p dist/amd64 dist/conf
	@cp build/service/$*.Dockerfile dist/amd64/$*.Dockerfile
	@cp -r config/$*/*.yaml dist/conf
	@sed -i -e '/#alpine-git.Dockerfile/ {' -e 'r build/base/amd64/alpine-git.Dockerfile' -e 'd' -e '}' dist/amd64/$*.Dockerfile
	@sed -i -e '/#alpine.Dockerfile/ {' -e 'r build/base/amd64/alpine.Dockerfile' -e 'd' -e '}' dist/amd64/$*.Dockerfile
	@sed -i -e '/#ubuntu-xenial.Dockerfile/ {' -e 'r build/base/amd64/ubuntu-xenial.Dockerfile' -e 'd' -e '}' dist/amd64/$*.Dockerfile
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn,direct go build -v -o dist/$* cmd/$*/main.go
	@docker build -f dist/amd64/$*.Dockerfile --tag ${MAKE_IMAGE} .

%.push.arm64: MAKE_IMAGE ?= ${IMAGE_REPOSITORY}/$*:${VERSION}-amd64
%.push.arm64: %.image.amd64
	docker push ${MAKE_IMAGE}

%.image.arm64: MAKE_IMAGE ?= ${IMAGE_REPOSITORY}/$*:${VERSION}-arm64
%.image.arm64:
	@mkdir -p dist/arm64 dist/conf
	@cp build/service/$*.Dockerfile dist/arm64/$*.Dockerfile
	@cp -r config/$*/*.yaml dist/conf
	@sed -i -e '/#alpine-git.Dockerfile/ {' -e 'r build/base/arm64/alpine-git.Dockerfile' -e 'd' -e '}' dist/arm64/$*.Dockerfile
	@sed -i -e '/#alpine.Dockerfile/ {' -e 'r build/base/arm64/alpine.Dockerfile' -e 'd' -e '}' dist/arm64/$*.Dockerfile
	@sed -i -e '/#ubuntu-xenial.Dockerfile/ {' -e 'r build/base/arm64/ubuntu-xenial.Dockerfile' -e 'd' -e '}' dist/arm64/$*.Dockerfile
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOPROXY=https://goproxy.cn,direct go build -v -o dist/$* cmd/$*/main.go
	@docker build -f dist/arm64/$*.Dockerfile --tag ${MAKE_IMAGE} .

.PHONY: clean
clean:
	@rm -rf dist