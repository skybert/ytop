
run: build
	./ytop

fmt:
	gofmt -s -w ./cmd ./internal

build: fmt
	go build -o ytop cmd/*.go

install: build
	mkdir -p ~/.local/bin
	cp ytop ~/.local/bin/.
