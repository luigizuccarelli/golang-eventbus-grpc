.PHONY: all test build clean

all: clean test build

build: 
	mkdir -p build
	go build -o build -tags real ./...

build-arm64:
	mkdir -p buildarm
	GOOS=linux GOARCH=arm64 go build -o buildarm -ldflags="-s -w" -tags real ./...

test:
	go test -v -coverprofile=tests/results/cover.out -tags fake ./pkg/...

cover:
	go tool cover -html=tests/results/cover.out -o tests/results/cover.html

clean:
	rm -rf build/*
	go clean ./...

container:
	podman build -t  tfld-docker-prd-local.repo.14west.io/golang-simple-oc4service:1.14.2 .

push:
	podman push tfld-docker-prd-local.repo.14west.io/golang-simple-oc4service:1.14.2 
