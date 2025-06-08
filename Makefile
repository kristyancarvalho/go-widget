APP_NAME := gowidget
BUILD_DIR := bin
MAIN_WEB := ./cmd/widget/main.go
MAIN_DESKTOP := ./cmd/desktop/main.go

.PHONY: all build build-web build-desktop run run-web run-desktop test clean

all: build

deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

build-web: deps
	@echo "🔧 Building web version..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME)-web $(MAIN_WEB)

build-desktop: deps
	@echo "🔧 Building desktop version..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME)-desktop.exe $(MAIN_DESKTOP)

build: build-web build-desktop

run-web: build-web
	@echo "🚀 Running web version..."
	@echo "📊 Dashboard: http://localhost:8080"
	./$(BUILD_DIR)/$(APP_NAME)-web

run-desktop: build-desktop
	@echo "🚀 Running desktop version..."
	./$(BUILD_DIR)/$(APP_NAME)-desktop.exe

run: run-desktop

test:
	@echo "🧪 Running tests..."
	go test ./...

clean:
	@echo "🧹 Cleaning up..."
	rm -rf $(BUILD_DIR)

build-windows: deps
	@echo "🔧 Building for Windows..."
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows.exe $(MAIN_DESKTOP)

package-windows: build-windows
	@echo "📦 Creating Windows package..."
	mkdir -p $(BUILD_DIR)/windows-package
	cp $(BUILD_DIR)/$(APP_NAME)-windows.exe $(BUILD_DIR)/windows-package/
	cp README.md $(BUILD_DIR)/windows-package/ 2>/dev/null || true
	cd $(BUILD_DIR) && zip -r github-widget-windows.zip windows-package/

help:
	@echo "Available commands:"
	@echo "  make build-web      - Build web version"
	@echo "  make build-desktop  - Build desktop version"
	@echo "  make run-web        - Run web version"
	@echo "  make run-desktop    - Run desktop version"
	@echo "  make build-windows  - Cross-compile for Windows"
	@echo "  make package-windows- Create Windows distribution package"