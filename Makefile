include .env
export
all: build run

deps:
	go mod tidy

build:
	go build -ldflags "-s -w" -o sca *.go
	strip ./sca

run:
	./sca