VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build test lint clean install

build:
	go build $(LDFLAGS) -o bin/manifestor ./cmd/m
	ln -sf manifestor bin/m
	ln -sf manifestor bin/mm

test:
	go test ./... -v

lint:
	go vet ./...

clean:
	rm -rf bin/

install: build
	install -d $(or $(GOPATH),$(HOME)/go)/bin
	cp bin/manifestor $(or $(GOPATH),$(HOME)/go)/bin/manifestor
	ln -sf manifestor $(or $(GOPATH),$(HOME)/go)/bin/m
	ln -sf manifestor $(or $(GOPATH),$(HOME)/go)/bin/mm
