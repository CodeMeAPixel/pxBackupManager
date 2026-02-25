# pxBackupManager

<div align="center">

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.23%2B-blue)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/CodeMeAPixel/pxBackupManager?include_prereleases)](https://github.com/CodeMeAPixel/pxBackupManager/releases)
[![GitHub Issues](https://img.shields.io/github/issues/CodeMeAPixel/pxBackupManager)](https://github.com/CodeMeAPixel/pxBackupManager/issues)

**A powerful, cross-platform backup manager for FiveM & RedM servers with cloud storage and Discord notifications**

[Features](#-features) • [Quick Start](#-quick-start) • [Documentation](#-documentation) • [Contributing](#-contributing)

</div>

---

## What is pxBackupManager?

pxBackupManager is a backup tool for FiveM & RedM servers and MariaDB databases. It supports S3 storage, automatic cleanup, and Discord notifications.

**Perfect for:**
- FiveM & RedM server administrators
- Businesses running multiple game servers
- DevOps engineers managing cloud infrastructure
- Teams requiring automated backups

---

## Features

### Core Capabilities
- **FiveM & RedM Server Backups** - Complete server directory backup with compression (works with both)
- **MariaDB Backups** - Database dumps with safety
- **Auto Compression** - Reduce backup size by 70-80% with gzip
- **Retention Policy** - Automatically clean up old backups (configurable)
- **Cross-Platform** - Works on Linux, macOS, and Windows

### Cloud Storage
- **S3 Compatible** - AWS S3, Backblaze B2, DigitalOcean Spaces, MinIO, and more
- **Direct Upload** - Backups streamed directly to S3 (no local storage needed)
- **Secure** - SSL/TLS encryption for all transfers

### Notifications & Monitoring
- **Discord Webhooks** - Notifications with backup status and file sizes
- **Structured Logging** - Integration with systemd journalctl
- **Status Codes** - Exit codes for success/failure integration

### Automation
- **Systemd Timers** - Linux scheduling (recommended)
- **Cron Support** - Unix scheduling
- **Docker Ready** - Container deployment
- **Windows Task Scheduler** - Windows scheduling

---

## Quick Start

### Installation

**Via direct build:**
```bash
git clone https://github.com/CodeMeAPixel/pxBackupManager.git
cd pxBackupManager
make build
./pxBackupManager -version
```

**Via releases:**
Download pre-built binaries from [GitHub Releases](https://github.com/CodeMeAPixel/pxBackupManager/releases)

### Basic Usage

```bash
# Backup both FiveM and database
./pxBackupManager -db-name "your_database"

# Backup only FiveM
./pxBackupManager -only-fivem

# With S3 upload and Discord notification
./pxBackupManager -db-name "fivem" \
  -s3-enabled \
  -s3-bucket "my-backups" \
  -s3-access-key "YOUR_KEY" \
  -s3-secret-key "YOUR_SECRET" \
  -discord-webhook "YOUR_WEBHOOK_URL" \
  -notify-success
```

---

## Documentation

Comprehensive documentation is available at [codemeapixel.dev/docs](https://codemeapixel.dev/docs/assets/server-management/pxbackupmanager)

| Document | Purpose |
|----------|---------|
| [Recommended Setup](https://codemeapixel.dev/docs/assets/server-management/pxbackupmanager/recommended-setup) | Production setup with directory structure and user permissions |
| [Installation Guide](https://codemeapixel.dev/docs/assets/server-management/pxbackupmanager/installation) | Setup, requirements, and build options |
| [Usage Guide](https://codemeapixel.dev/docs/assets/server-management/pxbackupmanager/usage) | Command-line flags and practical examples |
| [Configuration Guide](https://codemeapixel.dev/docs/assets/server-management/pxbackupmanager/configuration) | Advanced setup, S3, Discord, and multi-server |
| [Deployment Guide](https://codemeapixel.dev/docs/assets/server-management/pxbackupmanager/deployment) | Scheduling with systemd, cron, Docker, Windows |
| [Features & Capabilities](https://codemeapixel.dev/docs/assets/server-management/pxbackupmanager/features) | Detailed feature documentation and performance |
| [Troubleshooting Guide](https://codemeapixel.dev/docs/assets/server-management/pxbackupmanager/troubleshooting) | Common issues, monitoring, and recovery |

---

## Common Examples

### Daily Backup with Auto-Cleanup
```bash
./pxBackupManager -db-name "fivem" \
  -backup-dir /backups/fivem \
  -retention 30 \
  -compress
```

### Backup to Backblaze B2
```bash
./pxBackupManager -db-name "fivem" \
  -s3-enabled \
  -s3-endpoint "https://s3.us-west-000.backblazeb2.com" \
  -s3-bucket "fivem-backups" \
  -s3-access-key "YOUR_KEY" \
  -s3-secret-key "YOUR_APP_KEY" \
  -s3-upload-fivem \
  -s3-upload-mariadb
```

### Setup Systemd Timer (Linux)
See [Deployment Guide](https://codemeapixel.dev/docs/assets/server-management/pxbackupmanager/deployment) for complete setup instructions.

---

## Security

Read our [SECURITY.md](SECURITY.md) for security guidelines:

- Credential handling best practices
- Backup storage security
- S3 and Discord webhook security
- User access control recommendations
- Vulnerability reporting process

Key recommendations:
- Use dedicated backup user account
- Restrict backup directory permissions (`chmod 700`)
- Store credentials in environment variables or secure files
- Monitor backup logs for errors and failures

---

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for:

- Bug reporting guidelines
- Feature request process
- Development setup
- Code style guidelines
- Pull request requirements

---

## Requirements

### System
- **OS**: Linux, macOS, or Windows
- **Go**: 1.23 or later
- **Disk**: Sufficient space for backups (recommend separate disk)

### Optional Dependencies
- **MySQL Client**: For MariaDB backups (`mysqldump`)

---

## Performance

- **FiveM Backups**: 100GB+ tested
- **Database Backups**: 50GB+ tested
- **Compression**: 70-80% reduction in backup size
- **Memory**: Fixed allocation
- **Speed**: ~20-30% slower with compression enabled

---


## Command-Line Flags Summary

### Core Flags
| Flag | Default | Description |
|------|---------|-------------|
| `-fivem` | `/opt/fivem` | Path to FiveM server |
| `-db-name` | empty | Database name (required for MariaDB) |
| `-backup-dir` | `./backups` | Backup destination |
| `-retention` | `30` | Keep backups for N days |
| `-compress` | `true` | Enable gzip compression |
| `-version` | `false` | Show version |

### S3 Flags
| Flag | Default | Description |
|------|---------|-------------|
| `-s3-enabled` | `false` | Enable S3 uploads |
| `-s3-bucket` | empty | S3 bucket name |
| `-s3-endpoint` | empty | S3 endpoint (for non-AWS) |
| `-s3-access-key` | empty | S3 access key |
| `-s3-secret-key` | empty | S3 secret key |

### Discord Flags
| Flag | Default | Description |
|------|---------|-------------|
| `-discord-webhook` | empty | Discord webhook URL |
| `-notify-success` | `false` | Notify on success |
| `-notify-failure` | `false` | Notify on failure |

See [Usage Guide](https://codemeapixel.dev/docs/assets/server-management/pxbackupmanager/usage) for complete flag reference.

---

## License

This project is licensed under the **GNU Affero General Public License v3.0** (AGPL 3.0).

See [LICENSE](LICENSE) for details.

---

## Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [AWS SDK v2](https://github.com/aws/aws-sdk-go-v2)
- Inspired by the FiveM community

---

## Support

- **Report a bug**: [GitHub Issues](https://github.com/CodeMeAPixel/pxBackupManager/issues)
- **Request a feature**: [GitHub Issues](https://github.com/CodeMeAPixel/pxBackupManager/issues)
- **Report security issue**: [hey@codemeapixel.dev](mailto:hey@codemeapixel.dev)
- **Join the community**: [Discord Server](https://discord.gg/BsEhHBTbXw)
- **View changelog**: [CHANGELOG.md](CHANGELOG.md)
- **Read documentation**: [`docs/`](docs/)

---

<div align="center">

**[⬆ back to top](#pxbackupmanager)**

</div>
