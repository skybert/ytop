
run: build
	./build/ytop

fmt:
	gofmt -s -w ./cmd ./pkg

.phony: version
version ?= $(shell git describe --tags --always --dirty)

.phony: build
build: fmt
	mkdir -p build
	go build \
	  -o build/ytop \
	  -ldflags "-X main.Version=$(version)" \
	  cmd/*.go

.phony: lint
lint: fmt
	staticcheck ./...

.phony: vuln
vuln:
	govulncheck ./...

install: build
	mkdir -p ~/.local/bin
	cp build/ytop ~/.local/bin/.
