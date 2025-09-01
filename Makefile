APP := icignore
BIN_DIR := bin
PKG := ./cmd/icignore
VERSION ?= 0.1.0

.PHONY: build install uninstall clean test fmt vet release

build:
	@mkdir -p $(BIN_DIR)
	GO111MODULE=on go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BIN_DIR)/$(APP) $(PKG)

install: build
	@if [ -d "/opt/homebrew/bin" ]; then \
	  echo "Installing to /opt/homebrew/bin"; \
	  install -m 0755 $(BIN_DIR)/$(APP) /opt/homebrew/bin/$(APP); \
	else \
	  echo "Installing to /usr/local/bin"; \
	  install -m 0755 $(BIN_DIR)/$(APP) /usr/local/bin/$(APP); \
	fi
	@echo "Installed $(APP)."

uninstall:
	rm -f /usr/local/bin/$(APP) /opt/homebrew/bin/$(APP) 2>/dev/null || true

clean:
	rm -rf $(BIN_DIR)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

release:
	goreleaser release --clean
