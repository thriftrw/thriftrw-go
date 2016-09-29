#!/bin/bash

set -e
set -x

# This is a workaround to get mockgen to generate mocks for internal
# interfaces. It is based on the script posted in
# https://github.com/golang/mock/issues/29
#
# If you see an error like,
#
#   2016/08/26 15:13:45 Loading input failed: gob: name not registered for interface: "go.uber.org/thriftrw/vendor/github.com/golang/mock/mockgen/model.PredeclaredType"
#
# Make sure that the mockgen you're using is the one vendored in your copy of
# this repo. In the example above, you'll have to do,
#
#   go install go.uber.org/thriftrw/vendor/github.com/golang/mock/mockgen

PACKAGE=go.uber.org/thriftrw/internal/plugin
INTERFACES=Handle,ServiceGenerator
DESTINATION=handletest/mock.go
PACKAGENAME=handletest

mkdir -p _mockgen
mockgen -prog_only "$PACKAGE" "$INTERFACES" > _mockgen/main.go
go build -o _mockgen/gen _mockgen/main.go
mockgen -self_package "$PACKAGENAME" -package "$PACKAGENAME" -destination "$DESTINATION" -exec_only _mockgen/gen "$PACKAGE" "$INTERFACES"
rm -r _mockgen
