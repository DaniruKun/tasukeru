BINARY_NAME=tasukeru

all: build

build:
	go build -o bin/${BINARY_NAME} .

compile:
	@echo "Compiling for every OS and Platform"
	GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/${BINARY_NAME}-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build -o bin/${BINARY_NAME}-windows-arm64.exe .
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o bin/${BINARY_NAME}-linux-arm64 .

clean:
	go clean
	rm ./bin/*

run:
	go run main.go

format:
	@echo "Formatting the entire project"
	go fmt
