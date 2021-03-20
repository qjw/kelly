REPO = github.com/qjw/kelly
# GIT_COMMIT := $(shell git rev-parse --short HEAD)
# GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
# GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "*" || true)
# BUILD_VERSION := $(shell git describe --abbrev=0 --tags)
# BUILD_TIME := $(shell date +%FT%T%z)

# LDFLAGS := " \
# -X ${REPO}/pkg/version.buildGitCommit=${GIT_COMMIT}${GIT_DIRTY} \
# -X ${REPO}/pkg/version.buildGitBranch=${GIT_BRANCH} \
# -X ${REPO}/pkg/version.buildVersion=${BUILD_VERSION} \
# -X ${REPO}/pkg/version.buildTime=${BUILD_TIME} \
# "

# BINARIES:=./cmd/gateway


.PHONY: build
build:
	go build -mod=vendor

.PHONY: mod
mod:
	go mod tidy && go mod vendor
