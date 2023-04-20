.DEFAULT_GOAL:=help

.EXPORT_ALL_VARIABLES:

ifndef VERBOSE
.SILENT:
endif

SHELL=/bin/bash -o pipefail -o errexit

export GO111MODULE=on
# export GIN_MODE=release

DISTDIR:=$(shell pwd)/sysbase
T:=$(shell mkdir -p $(DISTDIR)/bin)

REPO_INFO ?= $(shell git config --get remote.origin.url)
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
COMMIT_TIME ?= $(shell git show --pretty="%ci %cr" | head -1 | awk '{print $$1"/"$$2"/"$$3}')
BUILD_TIME ?= $(shell date "+%Y-%m-%d/%H:%M:%S/%z")
VERSION ?= v1.0.0

LDFLAGS=-w -X sysbase/version.Pkg=$(REPO_INFO) -X sysbase/version.Version=$(VERSION) -X sysbase/version.GitCommitSha=$(COMMIT_SHA) -X sysbase/version.GitCommitTime=$(COMMIT_TIME) -X sysbase/version.BuildTime=$(BUILD_TIME)

HOST_ARCH = $(shell which go >/dev/null 2>&1 && go env GOARCH)
HOST_OS = $(shell which go >/dev/null 2>&1 && go env GOOS)
ARCH ?= $(HOST_ARCH)
OS ?= $(HOST_OS)

ifeq ($(ARCH),)
    $(error mandatory variable ARCH is empty, either set it when calling the command or make sure 'go env GOARCH' works)
endif

help:  ## 请查看下面帮助
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: build
build:  ## 编译 sysbase
	@cd backend && \
		rm -rf $(DISTDIR)/bin/sysbase-$(ARCH) && \
		go mod tidy && \
		CGO_ENABLED=0 \
		GOOS=$(OS) \
		GOARCH=$(ARCH) \
		go build -ldflags "$(LDFLAGS)" -o $(DISTDIR)/bin/sysbase-$(VERSION)-$(OS)-$(ARCH) ./

.PHONY: runf
runf:  ## 运行 sysbase 前端
	@cd frontend && \
		yarn serve

.PHONY: runb
runb:  ## 运行 sysbase 后端
	@cd backend && \
		go mod tidy && \
		CGO_ENABLED=0 \
		GOOS=$(OS) \
		GOARCH=$(ARCH) \
		go run -ldflags "$(LDFLAGS)" ./

.PHONY: runbin
runbin:  ## 运行 sysbase 后端
	@cd $(DISTDIR) && \
	./bin/sysbase-$(VERSION)-$(OS)-$(ARCH)

.PHONY: package
package: dist build ## 打包 sysbase
	@/bin/cp -rpf backend/etc $(DISTDIR)/ && \
		/bin/cp -rpf backend/etc/config.js $(DISTDIR)/static

.PHONY: dist
dist:  ## 打包 sysbase 前端
	@cd frontend && \
		yarn build && \
		rm -rf $(DISTDIR)/static && \
		mv dist $(DISTDIR)/static

.PHONY: clean
clean:  ## 清理 sysbase
	@rm -rf $(DISTDIR) && \
		mkdir -p $(DISTDIR)/bin
