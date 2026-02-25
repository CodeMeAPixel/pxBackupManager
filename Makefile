all: build

BINARY_DIR := binaries
PLATFORMS := linux-amd64 linux-arm64 windows-amd64 darwin-amd64 darwin-arm64

# Detect OS for cross-platform commands
ifeq ($(OS),Windows_NT)
    MKDIR = mkdir
    RMDIR = powershell -Command if(Test-Path \\'$(BINARY_DIR)\\') { Remove-Item -Path \\'$(BINARY_DIR)\\' -Recurse -Force }
    RM = del /q
else
    MKDIR = mkdir -p
    RMDIR = rm -rf $(BINARY_DIR)
    RM = rm -f
endif

build:
	@echo "Building pxBackupManager for current platform, please wait..."
	@go build -v -o pxBackupManager

build-all: clean-binaries
	@echo "Building pxBackupManager for all platforms..."
	@$(MKDIR) $(BINARY_DIR)
	@GOOS=linux GOARCH=amd64 go build -v -o $(BINARY_DIR)/pxBackupManager-linux-amd64
	@GOOS=linux GOARCH=arm64 go build -v -o $(BINARY_DIR)/pxBackupManager-linux-arm64
	@GOOS=windows GOARCH=amd64 go build -v -o $(BINARY_DIR)/pxBackupManager-windows-amd64.exe
	@GOOS=darwin GOARCH=amd64 go build -v -o $(BINARY_DIR)/pxBackupManager-darwin-amd64
	@GOOS=darwin GOARCH=arm64 go build -v -o $(BINARY_DIR)/pxBackupManager-darwin-arm64
	@echo "All platform builds complete!"

run: build
	@echo "Running pxBackupManager..."
	@./pxBackupManager

dev:
	@echo "Running pxBackupManager in development mode..."
	@go run main.go

backup-fivem:
	@go run main.go -fivem /opt/fivem -backup-dir ./backups

backup-mariadb:
	@go run main.go -only-mariadb -db-name "your_database_name" -backup-dir ./backups

backup-all:
	@go run main.go -fivem /opt/fivem -db-name "your_database_name" -backup-dir ./backups

clean:
	@echo "Cleaning build artifacts..."
	@$(RM) pxBackupManager pxBackupManager.exe

clean-binaries:
	@echo "Cleaning binary directory..."
	@$(RMDIR)

test:
	@echo "Running tests..."
	@go test -v ./...

.PHONY: all build build-all run dev backup-fivem backup-mariadb backup-all clean clean-binaries test
