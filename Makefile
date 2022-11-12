BINARY_NAME := tasukeru

.PHONY: all build compile-cli compile-windows compile-mac clean run format

.DEFAULT: build

build:
	@echo "Building tasukeru build for the current platform"
	go build -o bin/${BINARY_NAME} -ldflags '-s -w' .

compile-cli:
	@echo "Compiling simple CLI for every OS and Platform"
	GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}-darwin-amd64 .
	fyne-cross windows -ldflags '-s -w' -arch amd64 -console -name tasukeru-cli.exe

compile-windows: FyneApp.toml
	@echo "Building native Windows cross compiled build"
	fyne-cross windows -ldflags '-s -w'

compile-mac:
	@echo "Building native Mac app"
	fyne package -os darwin -icon Tasukeru.png --release

compile-linux: FyneApp.toml
	@echo "Building native Linux app"
	fyne-cross linux -ldflags '-s -w'	

all: compile-cli compile-windows compile-mac compile-linux

clean:
	go clean
	rm ./bin/*

run: build
	@echo "Running dev build of Tasukeru"
	./bin/${BINARY_NAME}

format:
	@echo "Formatting the entire project"
	go fmt
