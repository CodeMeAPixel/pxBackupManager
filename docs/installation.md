# Installation Guide

## Requirements

### System Requirements
- **Operating System**: Linux, macOS, or Windows
- **Go**: Version 1.23 or later
- **Disk Space**: Sufficient space for backups (recommend separate disk/filesystem)

### Dependencies

#### For FiveM Backups
- Built into Go (uses `archive/tar` and `compress/gzip`)
- No external utilities required (cross-platform compatible)

#### For MariaDB/MySQL Backups
- `mysqldump` utility (usually included with MySQL/MariaDB client)
  - Linux: `sudo apt-get install mysql-client` (Ubuntu/Debian) or `mariadb-client` (MariaDB)
  - macOS: `brew install mysql-client`
  - Windows: Include MySQL client in PATH or specify full path

## Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/CodeMeAPixel/pxBackupManager.git
cd pxBackupManager

# Build the project
make build

# Verify installation
./pxBackupManager -version
```

This creates an executable named `pxBackupManager` (or `pxBackupManager.exe` on Windows).

### Option 2: Download Pre-built Binary

Pre-built binaries are available on the [GitHub Releases](https://github.com/CodeMeAPixel/pxBackupManager/releases) page.

```bash
# Download (replace VERSION with actual release)
wget https://github.com/CodeMeAPixel/pxBackupManager/releases/download/vVERSION/pxBackupManager-linux-amd64

# Make executable
chmod +x pxBackupManager-linux-amd64

# Run
./pxBackupManager-linux-amd64 -version
```

### Option 3: Install via Makefile

```bash
# Build and install to /usr/local/bin
sudo make install

# Verify
pxBackupManager -version
```

## Verifying Installation

```bash
# Check version
pxBackupManager -version

# Display help
pxBackupManager -help

# Verify database connectivity
pxBackupManager -db-host localhost -db-user root -db-pass "password" -only-mariadb -db-name test
```

## Platform-Specific Setup

### Linux (Systemd)

```bash
# Create backup directory
sudo mkdir -p /backups
sudo chmod 700 /backups

# Create dedicated backup user
sudo useradd -r -s /bin/false backup

# Set permissions
sudo chown backup:backup /backups
```

### macOS

```bash
# Install MySQL client if needed
brew install mysql-client

# Create backup directory
mkdir -p ~/backups
chmod 700 ~/backups
```

### Windows

```powershell
# Create backup directory
New-Item -ItemType Directory -Path "C:\Backups" -Force

# Add to PATH if using -mysqldump from installation location
$env:Path += ";C:\Program Files\MySQL\MySQL Server 8.0\bin"
```

## Next Steps

- See [Usage Guide](usage.md) for command-line flags and basic usage
- See [Deployment Guide](deployment.md) for scheduling backups
- See [Configuration Guide](configuration.md) for advanced examples
