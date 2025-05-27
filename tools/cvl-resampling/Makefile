.PHONY: all clean build-linux build-macos build-macos-arm build-windows fmt

# Binary name
BIN := cvl-resampling

# Output folder
OUTPUT_DIR := ./bin

# Setup the -ldflags option for go build here, for example to set the version variable
LDFLAGS := -ldflags="-X main.version=1.0.0"

all: fmt build-linux build-macos build-macos-arm build-windows

# Format the code and remove unused imports
fmt:
	go fmt ./...
	go vet ./...

build-linux: fmt
	mkdir -p $(OUTPUT_DIR)/linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/linux/$(BIN).linux ./cvl-resampling

build-macos: fmt
	mkdir -p $(OUTPUT_DIR)/macos
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/macos/$(BIN) ./cvl-resampling

build-macos-arm: fmt
	mkdir -p $(OUTPUT_DIR)/macos-arm
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/macos-arm/$(BIN) ./cvl-resampling

build-windows: fmt
	mkdir -p $(OUTPUT_DIR)/windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/windows/$(BIN).exe ./cvl-resampling

clean:
	rm -rf $(OUTPUT_DIR)