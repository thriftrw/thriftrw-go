export GO15VENDOREXPERIMENT=1

PACKAGES := $(shell glide novendor)

.PHONY: build
build:
	go build

.PHONY: test
test: build
	go test $(PACKAGES)

.PHONY: cover
cover:
	./scripts/cover.sh $(shell go list $(PACKAGES))
	go tool cover -html=cover.out -o cover.html

.PHONY: clean
clean:
	go clean
	rm -rf cover/cover*.out cover.html cover.out

.PHONY: install
install:
	glide --version || go get github.com/Masterminds/glide
	glide install

##############################################################################
# CI

.PHONY: install_ci
install_ci: install
	go get github.com/wadey/gocovmerge
	go get github.com/mattn/goveralls
	go get golang.org/x/tools/cmd/cover

build_ci: build

# Tests don't need to be run separately because goveralls takes care of
# running them.

.PHONY: test_ci
test_ci: build_ci
	./scripts/cover.sh $(shell go list $(PACKAGES))
	goveralls -coverprofile=cover.out -service=travis-ci
