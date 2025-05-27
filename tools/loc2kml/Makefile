.PHONY: all clean build-linux build-macos build-macos-arm build-windows

# Binary names
BIN_REPUB := loc2kml

# Output folder
OUTPUT_DIR := ./bin

# Setup the -ldflags option for go build here, for example to set the version variable
LDFLAG's := -ldflags="-X main.version=1.0.0"

all: build-linux build-macos build-macos-arm build-windows

build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/linux/$(BIN_REPUB).linux ./src/loc2kml.go

build-macos:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/macos/$(BIN_REPUB) ./src/loc2kml.go

build-macos-arm:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/macos-arm/$(BIN_REPUB) ./src/loc2kml.go

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/windows/$(BIN_REPUB).exe ./src/loc2kml.go

clean:
	rm -rf $(OUTPUT_DIR)
