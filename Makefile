# Set version if unset
ifeq ($(shell git symbolic-ref -q --short HEAD),)
VERSION ?= $(shell git describe)
else
VERSION ?= 0.0.0-$(shell git symbolic-ref -q --short HEAD)-$(shell git rev-parse HEAD)
endif

BUILD_FLAGS := -trimpath -mod=readonly -ldflags "-s -w -X tm/m/v2/version.Version=$(VERSION)"

build:
	go build $(BUILD_FLAGS) -o build/tm

install:
	go install $(BUILD_FLAGS)

.PHONY: build install
