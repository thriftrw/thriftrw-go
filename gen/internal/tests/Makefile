ROOT = ../../..
THRIFTRW = $(ROOT)/build/bin/thriftrw
THRIFT_FILES = $(wildcard thrift/*.thrift)
PACKAGES = $(patsubst %.thrift, %, $(notdir $(THRIFT_FILES)))

THRIFTRW_FLAGS = --pkg-prefix go.uber.org/thriftrw/gen/internal/tests

.PHONY: all
all: $(PACKAGES)

.PHONY: clean
clean:
	make -C $(ROOT) clean
	rm -rf $(PACKAGES)

$(THRIFTRW):
	make -C $(ROOT) build

nozap: thrift/nozap.thrift $(THRIFTRW)
	$(THRIFTRW) $(THRIFTRW_FLAGS) --no-recurse --no-zap $<

enum-text-marshal-strict: thrift/enum-text-marshal-strict.thrift $(THRIFTRW)
	$(THRIFTRW) $(THRIFTRW_FLAGS) --no-recurse --enum-text-marshal-strict $<

%: thrift/%.thrift $(THRIFTRW)
	$(THRIFTRW) $(THRIFTRW_FLAGS) $<
