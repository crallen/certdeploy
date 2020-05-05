all: build

build:
	go build -o build/certdeploy .

test:
	go test ./...

clean:
	rm -rf build

.PHONY: all build test clean