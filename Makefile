VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build test lint clean install

build:
	go build $(LDFLAGS) -o bin/m ./cmd/m
	ln -sf m bin/mm

test:
	go test ./... -v

lint:
	go vet ./...

clean:
	rm -rf bin/

install: build
	cp bin/m $(GOPATH)/bin/m 2>/dev/null || cp bin/m $(HOME)/go/bin/m
	ln -sf m $(GOPATH)/bin/mm 2>/dev/null || ln -sf m $(HOME)/go/bin/mm
