.PHONY: build test clean install

BIN_NAME=omc-tui
BIN_DIR=bin
INSTALL_DIR=$(HOME)/.local/bin
GO=$(HOME)/.local/go/bin/go

build:
	@echo "Building $(BIN_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BIN_NAME) ./cmd/omc-tui
	@echo "Build complete: $(BIN_DIR)/$(BIN_NAME)"

test:
	@echo "Running tests..."
	$(GO) test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	@echo "Clean complete"

install: build
	@echo "Installing $(BIN_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	cp $(BIN_DIR)/$(BIN_NAME) $(INSTALL_DIR)/$(BIN_NAME)
	@echo "Install complete: $(INSTALL_DIR)/$(BIN_NAME)"
	@echo "Ensure $(INSTALL_DIR) is in your PATH"
