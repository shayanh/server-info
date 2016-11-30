GIT ?= git
GO ?= go
COMMIT := $(shell $(GIT) log -1 --format="%H")
VERSION ?= $(shell $(GIT) describe --tags --exact-match --abbrev=0 --tags ${COMMIT} 2> /dev/null || echo "$(shell $(GIT) log -1 --format="%h")")
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
TARGET := github.com/shayanh/server-info
LD_FLAGS := -X $(TARGET)/common.Version=$(VERSION) -X $(TARGET)/common.BuildTime=$(BUILD_TIME)
FORMAT := '{{ join .Deps " " }}'
DOCKER_IMAGE := "quay.io/server-info:$(VERSION)"

.PHONY: help clean build

clean:
	rm talk 1>/dev/null 2>/dev/null || exit 0

build: main.go clean
	$(GO) build -o="server-info" -ldflags="$(LD_FLAGS)" $(TARGET)

docker: Dockerfile
	docker build -t $(DOCKER_IMAGE) .

push:
	docker push $(DOCKER_IMAGE)
