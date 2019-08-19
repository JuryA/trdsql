VERSION=$(shell git describe --tags)

GOCMD=go
TAGS="json1"
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
BUILDFLAG=-tags $(TAGS) -ldflags="-X github.com/noborus/trdsql.Version=$(VERSION)"
GOBUILD=$(GOCMD) build $(BUILDFLAG)
GOTEST=$(GOCMD) test -tags $(TAGS)
GOINSTALL=$(GOCMD) install $(BUILDFLAG)

GOXCMD=gox -cgo $(BUILDFLAG)

TARGET="dist/trdsql-$(VERSION)-{{.OS}}_{{.Arch}}/{{.Dir}}"

BINARY_NAME := trdsql
SRCS := $(shell git ls-files '*.go')

all: test build

.PHONY: test
test: $(SRCS)
	$(GOTEST)

.PHONY: build
build: trdsql

$(BINARY_NAME): $(SRCS)
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/trdsql

.PHONY: pkg
pkg: linux_pkg windows_pkg

.PHONY: linux_pkg
linux_pkg:
	$(GOXCMD) -os "linux" -arch "amd64" -output $(TARGET) ./cmd/trdsql

.PHONY: windows_pkg
windows_pkg:
	CC=x86_64-w64-mingw32-gcc $(GOXCMD) -os "windows" -arch "amd64" -output $(TARGET) ./cmd/trdsql

.PHONY: macOS_pkg
macOS_pkg:
	$(GOXCMD) -os "darwin" -arch "amd64" -output ${TARGET} ./cmd/trdsql

.PHONY: install
install:
	$(GOINSTALL) ./cmd/trdsql

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf dist
