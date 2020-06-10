all: build unit-test
.PHONY: all

build:
	go build -o /dev/null

unit-test:
	go test -race -coverprofile=cover.txt -covermode=atomic
