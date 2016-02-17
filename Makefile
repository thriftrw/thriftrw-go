export GO15VENDOREXPERIMENT=1

PACKAGES := $(shell glide novendor)

.PHONY: clean
clean:
	go clean

.PHONY: build
build:
	go build

.PHONY: install
install:
	glide --version || go get github.com/Masterminds/glide
	glide install

.PHONY: test
test: build
	go test $(PACKAGES)

.PHONY: cover
cover:
	@$(foreach pkg, $(shell go list $(PACKAGES) | cut -d/ -f4-), \
		go test ./$(pkg) -v -cover &&) echo "success"

##############################################################################
# CI

.PHONY: install_ci
install_ci: install
	go get github.com/axw/gocov/gocov
	go get github.com/mattn/goveralls
	go get golang.org/x/tools/cmd/cover

build_ci: build

# Tests don't need to be run separately because goveralls takes care of
# running them.

.PHONY: test_ci
test_ci: build_ci
	goveralls -service=travis-ci -v $(PACKAGES)
