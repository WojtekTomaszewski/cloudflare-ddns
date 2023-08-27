BINARY_NAME := cloudflare-ddns
GOOS := linux
GOARCH := amd64

build: clean
	@echo "Building binary..."
	@mkdir -p bin
	@GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/${BINARY_NAME} cmd/cloudflare-ddns/*.go
	
run:
	@echo Run app...
	go run cmd/cloudflare-ddns/main.go

clean:
	@echo Cleanup...
	@go clean
	@rm -f bin/${BINARY_NAME}

build-image:
	@echo "Bulding docker image..."
	docker build . -t docker.io/wojciecht/cloudflare-ddns:latest