BINARY_NAME=ping-1.0.0

.DEFAULT_GOAL := run

build:
	go build -o ./ping/ping.exe ping/main.go

build-all:
	GOARCH=amd64 GOOS=windows go build -o ./builds/${BINARY_NAME}-windows.exe ping/main.go
	GOARCH=amd64 GOOS=linux go build -o ./builds/${BINARY_NAME}-linux ping/main.go
	GOARCH=amd64 GOOS=darwin go build -o ./builds/${BINARY_NAME}-macOS ping/main.go

run: build
	./ping/ping

clean:
	go clean
	rm ./ping/ping

clean-all:
	go clean
	rm ./builds/${BINARY_NAME}-windows 
	rm ./builds/${BINARY_NAME}-linux 
	rm ./builds/${BINARY_NAME}-macOS
	