###########################################
.EXPORT_ALL_VARIABLES:
VERSION_PKG := github.com/easysoft/qcadmin/common
ROOT_DIR := $(CURDIR)
BUILD_DIR := $(ROOT_DIR)/_output
BIN_DIR := $(BUILD_DIR)/bin
GO111MODULE = on
GOPROXY = https://goproxy.cn,direct
GOSUMDB = sum.golang.google.cn

BUILD_RELEASE   ?= $(shell cat VERSION || echo "0.0.1")
BUILD_DATE := $(shell date "+%F %T")
GIT_COMMIT := $(shell git rev-parse --short HEAD || echo "abcdefgh")
APP_VERSION := ${BUILD_RELEASE}-${BUILD_DATE}-${GIT_COMMIT}

LDFLAGS := "-w -s \
	-X '$(VERSION_PKG).Version=$(BUILD_RELEASE)' \
	-X '$(VERSION_PKG).BuildDate=$(BUILD_DATE)' \
	-X '$(VERSION_PKG).GitCommitHash=$(GIT_COMMIT)' \
	-X 'k8s.io/client-go/pkg/version.gitVersion=${BUILD_RELEASE}' \
  -X 'k8s.io/client-go/pkg/version.gitCommit=${GIT_COMMIT}' \
  -X 'k8s.io/client-go/pkg/version.gitTreeState=dirty' \
  -X 'k8s.io/client-go/pkg/version.buildDate=${BUILD_DATE}' \
	-X 'k8s.io/client-go/pkg/version.gitMajor=1' \
	-X 'k8s.io/client-go/pkg/version.gitMinor=23' \
  -X 'k8s.io/component-base/version.gitVersion=${BUILD_RELEASE}' \
  -X 'k8s.io/component-base/version.gitCommit=${GIT_COMMIT}' \
  -X 'k8s.io/component-base/version.gitTreeState=dirty' \
	-X 'k8s.io/component-base/version.gitMajor=1' \
	-X 'k8s.io/component-base/version.gitMinor=23' \
  -X 'k8s.io/component-base/version.buildDate=${BUILD_DATE}'"

##########################################################################

help: ## this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

gencopyright: ## add copyright
	@bash hack/scripts/gencopyright.sh

doc: ## gen docs
	go run ./docs/docs.go

fmt: ## fmt code
	gofmt -s -w .
	goimports -w .
	@echo gofmt -l
	@OUTPUT=`gofmt -l . 2>&1`; \
	if [ "$$OUTPUT" ]; then \
		echo "gofmt must be run on the following files:"; \
        echo "$$OUTPUT"; \
        exit 1; \
    fi

lint: ## lint code
	@echo golangci-lint run --skip-files \".*test.go\" -v ./...
	@OUTPUT=`command -v golangci-lint >/dev/null 2>&1 && golangci-lint run --skip-files ".*test.go"  -v ./... 2>&1`; \
	if [ "$$OUTPUT" ]; then \
		echo "go lint errors:"; \
		echo "$$OUTPUT"; \
	fi

default: gencopyright fmt lint ## fmt code

build: ## build binary
	@echo "build bin ${GIT_VERSION} $(GIT_COMMIT) $(GIT_BRANCH) $(BUILD_DATE) $(GIT_TREE_STATE)"
	@GO_ENABLED=1 gox -osarch="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64" \
        -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" \
    		-ldflags ${LDFLAGS}

generate: ## generate
	go generate ./...

dev: ## dev test
	GO_ENABLED=1 gox -osarch="linux/amd64" \
        -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" \
    		-ldflags ${LDFLAGS}
	upx dist/qcadmin_linux_amd64

local: ## dev test
	GO_ENABLED=1 gox -os="darwin" -arch="amd64" -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" -ldflags ${LDFLAGS}

clean: ## clean
	rm -rf dist

upx: ## upx
	upx dist/*

fixme: ## fixme check
	grep -rnw "FIXME" internal

todo: ## todo check
	grep -rnw "TODO" internal

# Legacy code should be removed by the time of release
legacy: # legacy code check
	grep -rnw "\(LEGACY\|Deprecated\)" internal

.PHONY : build prod-docker dev-push clean


snapshot: ## local test goreleaser
	goreleaser release --snapshot --rm-dist --skip-publish
