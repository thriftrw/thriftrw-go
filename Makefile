export GO15VENDOREXPERIMENT=1

BUILD := ./build
PACKAGES := $(shell glide novendor)

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
	go test $(PACKAGES) -v

.PHONY: cover
cover:
	./scripts/cover.sh $(shell go list $(PACKAGES))
	go tool cover -html=cover.out -o cover.html

##############################################################################
# CI

.PHONY: install_ci
install_ci: install
	go get github.com/wadey/gocovmerge
	go get github.com/mattn/goveralls
	go get golang.org/x/tools/cmd/cover

# Tests don't need to be run separately because goveralls takes care of
# running them.

.PHONY: test_ci
test_ci:
	./scripts/cover.sh $(shell go list $(PACKAGES))
	goveralls -coverprofile=cover.out -service=travis-ci
