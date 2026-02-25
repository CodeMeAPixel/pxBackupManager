# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.2] - 02-25-2026

### Added
- **RedM Support**: Full compatibility with RedM servers alongside FiveM (identical backup mechanism)
- Documentation now hosted at [codemeapixel.dev](https://codemeapixel.dev/docs/assets/server-management/pxbackupmanager)

### Changed
- Updated all documentation to reflect RedM support
- Updated help text and flag descriptions to mention both FiveM and RedM
- Updated systemd service description to include RedM
- Enhanced README to highlight dual FiveM/RedM compatibility

## [0.1.1] - 02-25-2026

### Fixed
- **FiveM Backup Tar Creation**: Fixed "archive/tar: write too long" error by implementing robust error handling that skips problematic files (extremely long paths, inaccessible files, special files) instead of failing the entire backup
- **Backup Resilience**: FiveM backups now continue on file errors, only skipping files that can't be archived and logging warnings for each skipped file
- **File Type Validation**: Improved file type handling to ensure only regular files are copied (skip special files, symlinks with long targets, etc.)

### Changed
- FiveM backup now logs detailed warnings for each skipped file instead of silently failing

## [0.1.0] - 02-24-2026

### Added
- **Windows Support**: Full cross-platform compatibility (Windows, macOS, Linux)
- **FiveM Server Backups**: Complete server directory backup with compression
- **MariaDB Backups**: Database dumps with mysqldump
- **S3 Storage**: Upload backups to AWS S3 or S3-compatible services (Backblaze B2, DigitalOcean Spaces, MinIO)
- **Discord Webhooks**: Rich notifications for backup success/failure with detailed status information
- **Platform Detection**: Automatic detection of OS-specific paths (e.g., C:\FiveM\server on Windows, /opt/fivem on Unix)
- **MySQL Path Discovery**: Automatic detection of mysqldump on Windows common installation directories
- **Retention Policy**: Automatic cleanup of old backups (configurable)
- **Compression**: gzip compression support for backup files
- **Systemd Integration**: Service and timer files for Linux scheduling
- **Cron Support**: Cron job examples for traditional scheduling
- **Docker Support**: Dockerfile and docker-compose examples
- **Windows Task Scheduler**: Support for Windows scheduled tasks
- Comprehensive documentation in `/docs` directory:
  - `recommended-setup.md`: Production setup guide with directory structure and permissions
  - `installation.md`: Platform-specific installation instructions
  - `usage.md`: Command-line flags and practical examples
  - `configuration.md`: Advanced setup, S3, and Discord configuration
  - `deployment.md`: Multiple scheduling methods (systemd, cron, Docker, Windows Task Scheduler)
  - `troubleshooting.md`: Common issues and solutions
  - `features.md`: Detailed feature breakdown
- **Security Policy** (`SECURITY.md`) with vulnerability disclosure and best practices
- **Contributing Guidelines** (`CONTRIBUTING.md`)

### Changed
- Uses Go's archive/tar and compress/gzip packages for cross-platform compatibility
- Backup functions return filenames for better tracking and S3 URL generation

### Fixed
- Cross-platform compatibility for backup creation
- AWS SDK context usage (uses context.Background())
