all: build

build:
	go build -mod vendor -o build/certdeploy .

clean:
	rm -rf build

.PHONY: all build clean