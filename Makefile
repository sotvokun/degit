VERSION_VERSION := $(shell git describe --tags --always)
VERSION_DATETIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
VERSION_COMMIT := $(shell git rev-parse --short HEAD)
ifeq ($(OS), Windows_NT)
	VERSION_DATETIME := $(shell sh.exe -c "date -u +%Y-%m-%dT%H:%M:%SZ")
endif


BUILD_OS ?= $(shell go env GOOS)
BUILD_ARCH ?= $(shell go env GOARCH)
BUILD_BINARY_NAME := emit-$(BUILD_OS)-$(BUILD_ARCH)-$(VERSION_VERSION)
ifeq ($(BUILD_OS), windows)
	BUILD_BINARY_NAME := $(BUILD_BINARY_NAME).exe
endif
BUILD_OUTPUT_DIR ?= dist


build:
	GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) go build -o $(BUILD_OUTPUT_DIR)/$(BUILD_BINARY_NAME) \
		-ldflags "-X github.com/sotvokun/emit/internal/service/version.version=$(VERSION_VERSION) \
		-X github.com/sotvokun/emit/internal/service/version.commit=$(VERSION_COMMIT) \
		-X github.com/sotvokun/emit/internal/service/version.date=$(VERSION_DATETIME)" \
		github.com/sotvokun/emit/cmd/emit
