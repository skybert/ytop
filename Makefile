
run: build
	./ytop

fmt:
	gofmt -s -w ./src

build: fmt
	go build -o ytop src/*.go
