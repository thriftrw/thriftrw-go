ROOT = ../..
THRIFTRW = $(ROOT)/thriftrw-go
GO_FILES = $(patsubst %, %.go, $(wildcard *.thrift))

.PHONY: generate
generate: $(GO_FILES)

.PHONY: clean
clean:
	make -C $(ROOT) clean

$(THRIFTRW):
	make -C $(ROOT) build

%.go: % $(THRIFTRW)
	$(THRIFTRW) -o $@ $<
