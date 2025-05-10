# Define the application name
APP_NAME = json-key-checker

# Define the path to the main package source
BUILD_PATH = ./cmd/$(APP_NAME)

# Define the installation path
INSTALL_PATH = $(GOBIN)

.PHONY: build test install clean all help

all: build

build:
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) $(BUILD_PATH)
	@echo "Build complete. Executable created: ./bin/$(APP_NAME)"

test:
	@echo "Running tests..."
	go test ./... -v
	@echo "Tests complete."

install: build
	@echo "Installing $(APP_NAME) from ./bin/ to $(INSTALL_PATH)..."
	install bin/$(APP_NAME) $(INSTALL_PATH)
	@echo "$(APP_NAME) installed successfully."

clean:
	@echo "Cleaning build artifacts..."
	go clean $(BUILD_PATH)
	rm -f bin/$(APP_NAME)
	@echo "Clean complete."

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build   : Builds the $(APP_NAME) executable in the ./bin/ directory."
	@echo "  test    : Runs the tests for the application."
	@echo "  install : Installs the built $(APP_NAME) executable from ./bin/ to the defined INSTALL_PATH."
	@echo "  clean   : Removes build artifacts and the executable from the ./bin/ directory."
	@echo "  all     : Builds the application (default)."
	@echo "  help    : Displays this help message."