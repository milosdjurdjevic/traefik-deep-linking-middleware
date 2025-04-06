# Makefile

.PHONY: build test clean

build:
	go build -o traefik_deep_linking_middleware .

clean:
	go clean