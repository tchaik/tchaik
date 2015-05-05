GO := go
GOBUILD := $(GO) build
GOFMT := $(GO) fmt
GOINSTALL := $(GO) install
GOTEST := $(GO) test

PACKAGE = github.com/dhowden/tchaik
COMMANDS = cmd/tchaik cmd/tchimport cmd/tchstore
LIBRARIES = index store

ALL_LIST = $(COMMANDS) $(LIBRARIES)
BUILD_LIST = $(foreach cmd, $(COMMANDS), $(cmd)_build)
INSTALL_LIST = $(foreach cmd, $(COMMANDS), $(cmd)_install)
FMT_LIST = $(foreach path, $(ALL_LIST), $(path)_fmt)
TEST_LIST = $(foreach path, $(LIBRARIES), $(path)_test)

.PHONY: $(BUILD_LIST) $(FMT_LIST) ui gotest test build install fmt

$(BUILD_LIST): %_build:
	$(GOBUILD) ./$*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) ./$*
$(FMT_LIST): %_fmt:
	$(GOFMT) ./$*
$(TEST_LIST): %_test:
	$(GOTEST) ./$*
gotest: $(TEST_LIST)
golint:
	./scripts/verify-gofmt.sh $(ALL_LIST)

build: $(BUILD_LIST)
install: $(INSTALL_LIST)
fmt: $(FMT_LIST)

uideps:
	cd cmd/tchaik/ui; npm install
ui: uideps
	cd cmd/tchaik/ui; gulp

uilint: uideps
	cd cmd/tchaik/ui; gulp lint

test: gotest
lint: golint uilint

all: build
