SNAP_DIR  := test_data
STAGE_DIR := $(SNAP_DIR)/build

VERSION  := $(shell git describe --tags --always --dirty)
LDFLAGS  := -s -w -X 'main.version=$(VERSION)'

# CGO_ENABLED=0 produces a statically linked binary
GO_BUILD := CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)"

.PHONY: all build test stage snap install clean

all: build

build:
	$(GO_BUILD) -o modelctl .

test:
	go test -count 1 -failfast ./...

stage:
	rm -rf $(STAGE_DIR)
	mkdir -p $(STAGE_DIR)/cli/bin $(STAGE_DIR)/completer
	$(GO_BUILD) -o $(STAGE_DIR)/cli/bin/modelctl .
	cp completer/snap.bash $(STAGE_DIR)/completer/
	git describe --always --dirty > $(STAGE_DIR)/version

snap: stage test_data/dev/build.sh
	cd $(SNAP_DIR) && dev/build.sh

install: test_data/dev/build.sh
	cd $(SNAP_DIR) && dev/install.sh

clean:
	cd $(SNAP_DIR) && snapcraft clean
	rm -f modelctl
	rm -rf $(STAGE_DIR) $(SNAP_DIR)/*.snap $(SNAP_DIR)/*.comp

test_data/dev/build.sh:
	git submodule update --init test_data/dev
