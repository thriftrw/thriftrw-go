ROOT = ../..
THRIFTRW = $(ROOT)/thriftrw
THRIFT_FILES = $(wildcard thrift/*.thrift)
PACKAGES = $(patsubst %.thrift, %, $(notdir $(THRIFT_FILES)))

.PHONY: all
all: $(PACKAGES)

.PHONY: clean
clean:
	make -C $(ROOT) clean
	rm -rf $(PACKAGES)

$(THRIFTRW):
	make -C $(ROOT) build

%: thrift/%.thrift $(THRIFTRW)
	$(THRIFTRW) --no-recurse $<
