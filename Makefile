VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build test lint clean install

build:
	go build $(LDFLAGS) -o bin/manifestor ./cmd/manifestor

test:
	go test ./... -v

lint:
	go vet ./...

clean:
	rm -rf bin/

install: build
	cp bin/manifestor $(GOPATH)/bin/manifestor 2>/dev/null || \
	cp bin/manifestor $(HOME)/go/bin/manifestor
