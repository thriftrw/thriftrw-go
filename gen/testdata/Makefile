ROOT = ../..
THRIFTRW = $(ROOT)/thriftrw-go
THRIFT_FILES = $(wildcard thrift/*.thrift)
PACKAGES = $(patsubst %.thrift, %, $(THRIFT_FILES))

.PHONY: all
all: $(PACKAGES)

.PHONY: clean
clean:
	make -C $(ROOT) clean
	rm -rf $(PACKAGES)

$(THRIFTRW):
	make -C $(ROOT) build

%: %.thrift $(THRIFTRW)
	$(THRIFTRW) -no-recurse $<
