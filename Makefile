BINARY_NAME := cloudflare-ddns
GOOS := linux
GOARCH := amd64

build: clean
	@echo "Build executable file"
	@mkdir -p bin
	@GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/${BINARY_NAME} cmd/cloudflare-ddns/main.go
	
run:
	@echo Run app
	go run cmd/cloudflare-ddns/main.go

clean:
	@echo Cleanup
	@go clean
	@rm -f bin/${BINARY_NAME}