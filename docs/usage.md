# Usage Guide

## Basic Commands

### Backup Both FiveM and MariaDB

```bash
./pxBackupManager -db-name "your_database_name"
```

### Backup Only FiveM

```bash
./pxBackupManager -only-fivem
```

### Backup Only MariaDB

```bash
./pxBackupManager -only-mariadb -db-name "your_database_name"
```

### Display Version

```bash
./pxBackupManager -version
```

### Display Help

```bash
./pxBackupManager -help
```

## Command-line Flags

### FiveM Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `-fivem` | `/opt/fivem` (Linux/Mac)<br/>`C:\FiveM\server` (Windows) | Path to FiveM server directory |
| `-only-fivem` | `false` | Only backup FiveM server, skip MariaDB |

### MariaDB Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `-db-host` | `localhost` | MariaDB/MySQL host address |
| `-db-port` | `3306` | MariaDB/MySQL port |
| `-db-user` | `root` | MariaDB/MySQL user |
| `-db-pass` | empty | MariaDB/MySQL password |
| `-db-name` | empty | Database name to backup (required for MariaDB) |
| `-only-mariadb` | `false` | Only backup MariaDB database, skip FiveM |

### Backup Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `-backup-dir` | `./backups` | Backup destination directory |
| `-retention` | `30` | Backup retention period in days (0 = no cleanup) |
| `-compress` | `true` | Compress backups with gzip |
| `-cleanup` | `true` | Automatically clean up old backups |
| `-version` | `false` | Show version information |

### S3 Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `-s3-enabled` | `false` | Enable S3 backup uploads |
| `-s3-endpoint` | empty | S3 endpoint URL (for S3-compatible services) |
| `-s3-bucket` | empty | S3 bucket name |
| `-s3-region` | `us-east-1` | AWS region or S3-compatible region |
| `-s3-access-key` | empty | S3 access key |
| `-s3-secret-key` | empty | S3 secret key |
| `-s3-upload-fivem` | `false` | Upload FiveM backup to S3 |
| `-s3-upload-mariadb` | `false` | Upload MariaDB backup to S3 |

### Discord Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `-discord-webhook` | empty | Discord webhook URL for notifications |
| `-notify-success` | `false` | Send Discord notification on success |
| `-notify-failure` | `false` | Send Discord notification on failure |

## Common Examples

### Full Backup with Details

```bash
./pxBackupManager \
  -fivem /opt/fivem \
  -db-name "fivem_database" \
  -db-pass "your_password" \
  -backup-dir /backups/flux \
  -retention 60 \
  -compress
```

### Backup Without Compression

```bash
./pxBackupManager \
  -db-name "fivem_database" \
  -db-pass "your_password" \
  -backup-dir /mnt/backup/fivem \
  -compress=false
```

### Backup with S3 Upload (AWS)

```bash
./pxBackupManager \
  -db-name "fivem_database" \
  -db-pass "your_password" \
  -s3-enabled \
  -s3-bucket "my-backup-bucket" \
  -s3-region "us-west-2" \
  -s3-access-key "YOUR_ACCESS_KEY" \
  -s3-secret-key "YOUR_SECRET_KEY" \
  -s3-upload-fivem \
  -s3-upload-mariadb
```

### Backup with S3-Compatible Storage (Backblaze B2)

```bash
./pxBackupManager \
  -db-name "fivem_database" \
  -db-pass "your_password" \
  -s3-enabled \
  -s3-endpoint "https://s3.us-west-000.backblazeb2.com" \
  -s3-bucket "my-bucket" \
  -s3-region "us-west-000" \
  -s3-access-key "YOUR_KEY_ID" \
  -s3-secret-key "YOUR_APP_KEY" \
  -s3-upload-fivem
```

### Backup with Discord Notifications

```bash
./pxBackupManager \
  -db-name "fivem_database" \
  -db-pass "your_password" \
  -backup-dir /backups \
  -discord-webhook "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_TOKEN" \
  -notify-success \
  -notify-failure
```

### Disable Automatic Cleanup

```bash
./pxBackupManager \
  -db-name "fivem_database" \
  -db-pass "your_password" \
  -cleanup=false
```

### Custom Retention Period

```bash
./pxBackupManager \
  -db-name "fivem_database" \
  -db-pass "your_password" \
  -retention 90
```

## Output

### Successful Backup

```
===== PX BACKUP MANAGER =====
Backup Directory: /backups
Compression: true
Retention: 30 days
FiveM Backup: true
MariaDB Backup: true
=============================

Starting FiveM backup from /opt/fivem
Backup destination: /backups/fivem-backup-2024-01-15-143022.tar.gz
FiveM backup completed successfully
File: fivem-backup-2024-01-15-143022.tar.gz
Size: 45.23 MB

Starting MariaDB backup for database: fivem
Backup destination: /backups/mariadb-backup-2024-01-15-143022.tar.gz
MariaDB backup completed successfully
File: mariadb-backup-2024-01-15-143022.tar.gz
Size: 12.45 MB

===== BACKUP SUMMARY =====
[fivem] ✓ SUCCESS - 2345ms
Size: 45.23 MB
[mariadb] ✓ SUCCESS - 1234ms
Size: 12.45 MB

Total: 2 succeeded, 0 failed
Elapsed time: 3.579s
==========================
```

### Exit Codes

- `0` - All backups completed successfully
- `1` - One or more backups failed

## Environment Variables

You can pass sensitive credentials via environment variables instead of command-line flags:

```bash
export DB_PASSWORD="your_password"
export S3_ACCESS_KEY="your_key"
export S3_SECRET_KEY="your_secret"
export DISCORD_WEBHOOK="https://..."

./pxBackupManager -db-name "fivem"
```

Then reference in command or update the source to support env var loading.

## Next Steps

- See [Deployment Guide](deployment.md) for scheduling backups
- See [Configuration Guide](configuration.md) for advanced setups
- See [Troubleshooting Guide](troubleshooting.md) for common issues
