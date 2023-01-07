BINARY_NAME := cloudflare-ddns
UNAME := $(shell uname)

ifeq ($(UNAME), Linux)
	GOOS := linux
	GOARCH := amd64
endif
ifeq ($(UNAME), Darwin)
	GOOS := darwin
	GOARCH := amd64
endif



build: clean
	@echo "Build executable file"
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BINARY_NAME} *.go
	
run:
	@echo Run app
	go run *.go

clean:
	@echo Cleanup
	go clean
	rm -f ${BINARY_NAME}