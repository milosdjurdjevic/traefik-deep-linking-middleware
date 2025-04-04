# Makefile

.PHONY: build test clean

build:
	go build -o deep_linking .

test:
	go test -v ./...

clean:
	go clean
	rm -f deep_linking