# Variables
SRC_DIR = "."
MAIN = $(SRC_DIR)/main.go
CLIENT_DIR = ./client
GO_SERVER = go run $(MAIN)
REACT_SERVER = npm start --prefix $(CLIENT_DIR)

# Default task
.PHONY: all
all: run

# Run both the Go server and React frontend concurrently
.PHONY: run
run:
	@echo "Running both Go server and React frontend..."
	@$(GO_SERVER) & $(REACT_SERVER)

# Clean cache
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@go clean -cache
	@echo "Clean completed."

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@cd $(CLIENT_DIR) && npm install
	@echo "Dependencies installed."
