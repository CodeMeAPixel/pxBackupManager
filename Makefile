all: build

build:
	@echo "Building pxBackupManager, please wait..."
	@go build -v -o pxBackupManager

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
	@rm -f flux-backup

test:
	@echo "Running tests..."
	@go test -v ./...

.PHONY: all build run dev backup-fivem backup-mariadb backup-all clean test
