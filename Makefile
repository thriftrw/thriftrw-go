export GO15VENDOREXPERIMENT=1

BUILD := ./build
PACKAGES := $(shell glide novendor)
TEST_ARGS ?= -race -v

.PHONY: clean
clean:
	go clean
	rm -rf $(BUILD)

.PHONY: setup
setup:
	mkdir -p $(BUILD)

.PHONY: build
build: setup
	go build -o $(BUILD)/thriftrw

.PHONY: install
install:
	glide --version || go get github.com/Masterminds/glide
	glide install

.PHONY: test
test:
	go test $(PACKAGES) $(TEST_ARGS)
