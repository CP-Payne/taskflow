# Makefile

# --- Variables ---

# Attempt to automatically determine the Go module path
# Ensure you have initialized your Go module (go mod init your_module_path)
GOMODULE := $(shell go list -m)
ifeq ($(GOMODULE),)
    $(error Please initialize your Go module first: go mod init your_module_path)
endif

# Tools (assumes they are in your PATH)
PROTOC := protoc
GO_PLUGIN := protoc-gen-go
GRPC_PLUGIN := protoc-gen-go-grpc

# Directories
API_DIR := ./api
GEN_DIR := ./pkg/gen

# Find all .proto files within the api/<service>/v1 structure
PROTO_FILES := $(shell find $(API_DIR) -path '*/v1/*.proto' -print)

# --- Targets ---

# Default target (optional)
all: proto

# Generate Go code from proto files
.PHONY: proto
proto: check-tools
	@echo "Generating Go code for $(GOMODULE)..."
	@# Ensure the base generated code directory exists
	@mkdir -p $(GEN_DIR)
	@echo "Found proto files: $(PROTO_FILES)"
	$(PROTOC) \
		--proto_path=$(API_DIR) \
		--go_out=. \
		--go_opt=module=$(GOMODULE) \
		--go-grpc_out=. \
		--go-grpc_opt=module=$(GOMODULE) \
		$(PROTO_FILES)
	@echo "Protobuf Go code generation complete."

# Clean generated code
.PHONY: clean-proto
clean-proto:
	@echo "Cleaning generated protobuf Go files..."
	@rm -rf $(GEN_DIR)/*

# Check if required tools are installed
.PHONY: check-tools
check-tools:
	@command -v $(PROTOC) >/dev/null 2>&1 || { echo >&2 "$(PROTOC) not found. Please install protobuf compiler."; exit 1; }
	@command -v $(GO_PLUGIN) >/dev/null 2>&1 || { echo >&2 "$(GO_PLUGIN) not found. Please run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"; exit 1; }
	@command -v $(GRPC_PLUGIN) >/dev/null 2>&1 || { echo >&2 "$(GRPC_PLUGIN) not found. Please run: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"; exit 1; }

# Help target (optional)
.PHONY: help
help:
	@echo "Makefile targets:"
	@echo "  proto        : Generate Go code from .proto files."
	@echo "  clean-proto  : Remove generated Go protobuf files."
	@echo "  check-tools  : Verify required tools (protoc, plugins) are installed."
	@echo "  help         : Show this help message."


