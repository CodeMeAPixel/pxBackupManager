# Configuration Guide

## Table of Contents

1. [Basic Configuration](#basic-configuration)
2. [S3 Backup Configuration](#s3-backup-configuration)
3. [Discord Notifications](#discord-notifications)
4. [Environment Variables](#environment-variables)
5. [Configuration Files](#configuration-files)
6. [Advanced Examples](#advanced-examples)

---

## Basic Configuration

### Minimal Setup

```bash
./pxBackupManager -db-name "fivem"
```

This backs up:
- FiveM server from `/opt/fivem` (Linux/Mac) or `C:\FiveM\server` (Windows)
- MariaDB database "fivem"
- Backups stored in `./backups`
- Compressed with retention of 30 days

### Custom Directory Backup

```bash
./pxBackupManager \
  -fivem /opt/fivem-production \
  -backup-dir /mnt/backups/prod \
  -db-name "production_db"
```

### Multiple Database Backup

Run separate instances for each database:

```bash
# Database 1
./pxBackupManager -db-name "fivem_main"

# Database 2
./pxBackupManager -db-name "fivem_logs"

# Database 3
./pxBackupManager -db-name "website"
```

---

## S3 Backup Configuration

### Amazon AWS S3

```bash
./pxBackupManager \
  -db-name "fivem" \
  -s3-enabled \
  -s3-bucket "my-company-backups" \
  -s3-region "us-west-2" \
  -s3-access-key "AKIAIOSFODNN7EXAMPLE" \
  -s3-secret-key "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" \
  -s3-upload-fivem \
  -s3-upload-mariadb
```

### Backblaze B2

```bash
./pxBackupManager \
  -db-name "fivem" \
  -s3-enabled \
  -s3-endpoint "https://s3.us-west-000.backblazeb2.com" \
  -s3-bucket "fivem-backups" \
  -s3-region "us-west-000" \
  -s3-access-key "YOUR_S3_KEY_ID" \
  -s3-secret-key "YOUR_S3_APP_KEY" \
  -s3-upload-fivem \
  -s3-upload-mariadb
```

### DigitalOcean Spaces

```bash
./pxBackupManager \
  -db-name "fivem" \
  -s3-enabled \
  -s3-endpoint "https://nyc3.digitaloceanspaces.com" \
  -s3-bucket "fivem-backups" \
  -s3-region "nyc3" \
  -s3-access-key "YOUR_ACCESS_KEY" \
  -s3-secret-key "YOUR_SECRET_KEY" \
  -s3-upload-fivem \
  -s3-upload-mariadb
```

### MinIO (Self-Hosted)

```bash
./pxBackupManager \
  -db-name "fivem" \
  -s3-enabled \
  -s3-endpoint "https://minio.example.com:9000" \
  -s3-bucket "fivem-backups" \
  -s3-region "us-east-1" \
  -s3-access-key "minioadmin" \
  -s3-secret-key "minioadmin" \
  -s3-upload-fivem \
  -s3-upload-mariadb
```

### Conditional S3 Upload

Upload only FiveM to S3 (save costs):

```bash
./pxBackupManager \
  -db-name "fivem" \
  -s3-enabled \
  -s3-bucket "fivem-backups" \
  -s3-access-key "YOUR_KEY" \
  -s3-secret-key "YOUR_SECRET" \
  -s3-upload-fivem
```

Upload only database to S3:

```bash
./pxBackupManager \
  -db-name "fivem" \
  -s3-enabled \
  -s3-bucket "fivem-backups" \
  -s3-access-key "YOUR_KEY" \
  -s3-secret-key "YOUR_SECRET" \
  -s3-upload-mariadb
```

---

## Discord Notifications

### Setup Discord Webhook

1. Go to your Discord server
2. Server Settings → Integrations → Webhooks
3. Create New Webhook
4. Copy the full webhook URL

**Example URL:**
```
https://discord.com/api/webhooks/1234567890/abcdefghijklmnop
```

### Basic Notifications

Notify on success:
```bash
./pxBackupManager \
  -db-name "fivem" \
  -discord-webhook "https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN" \
  -notify-success
```

Notify on failure:
```bash
./pxBackupManager \
  -db-name "fivem" \
  -discord-webhook "https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN" \
  -notify-failure
```

Notify on both:
```bash
./pxBackupManager \
  -db-name "fivem" \
  -discord-webhook "https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN" \
  -notify-success \
  -notify-failure
```

### Complete Setup with All Features

```bash
./pxBackupManager \
  -fivem /opt/fivem \
  -db-name "fivem" \
  -db-host "localhost" \
  -db-user "backup" \
  -db-pass "your_password" \
  -backup-dir /backups/fivem \
  -retention 30 \
  -compress \
  -s3-enabled \
  -s3-bucket "fivem-backups" \
  -s3-access-key "YOUR_KEY" \
  -s3-secret-key "YOUR_SECRET" \
  -s3-upload-fivem \
  -s3-upload-mariadb \
  -discord-webhook "https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN" \
  -notify-success \
  -notify-failure
```

### Notification Format

The Discord notification includes:
- Backup status (Success/Failed)
- Timestamp
- Duration per service
- File size per backup
- S3 upload URLs (if applicable)
- Error messages (if failed)

---

## Environment Variables

Load configuration from environment:

```bash
export DB_HOST="localhost"
export DB_PORT="3306"
export DB_USER="backup"
export DB_PASS="your_password"
export DB_NAME="fivem"
export BACKUP_DIR="/backups/fivem"
export RETENTION="30"
export S3_ACCESS_KEY="YOUR_KEY"
export S3_SECRET_KEY="YOUR_SECRET"
export DISCORD_WEBHOOK="https://..."

./pxBackupManager
```

Or use `.env` file:

```bash
# .env file
DB_HOST=localhost
DB_PORT=3306
DB_USER=backup
DB_PASS=your_password
DB_NAME=fivem
BACKUP_DIR=/backups/fivem
RETENTION=30
S3_ACCESS_KEY=YOUR_KEY
S3_SECRET_KEY=YOUR_SECRET
DISCORD_WEBHOOK=https://...
```

---

## Configuration Files

### Linux/Mac Configuration

Create `/etc/pxBackupManager/config.sh`:

```bash
#!/bin/bash

# Database Configuration
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="backup"
DB_PASS="your_password"
DB_NAME="fivem"

# Backup Configuration
BACKUP_DIR="/backups/fivem"
FIVEM_PATH="/opt/fivem"
RETENTION="30"
COMPRESS="true"

# S3 Configuration
S3_ENABLED="true"
S3_ENDPOINT="https://s3.amazonaws.com"
S3_BUCKET="my-backups"
S3_REGION="us-west-2"
S3_ACCESS_KEY="YOUR_KEY"
S3_SECRET_KEY="YOUR_SECRET"

# Discord Configuration
DISCORD_WEBHOOK="https://discord.com/api/webhooks/..."
NOTIFY_SUCCESS="true"
NOTIFY_FAILURE="true"

# Build command
export PXBM_CMD="pxBackupManager \
  -fivem \"$FIVEM_PATH\" \
  -db-host \"$DB_HOST\" \
  -db-port \"$DB_PORT\" \
  -db-user \"$DB_USER\" \
  -db-pass \"$DB_PASS\" \
  -db-name \"$DB_NAME\" \
  -backup-dir \"$BACKUP_DIR\" \
  -retention \"$RETENTION\" \
  -compress=\"$COMPRESS\" \
  -s3-enabled=\"$S3_ENABLED\" \
  -s3-endpoint \"$S3_ENDPOINT\" \
  -s3-bucket \"$S3_BUCKET\" \
  -s3-region \"$S3_REGION\" \
  -s3-access-key \"$S3_ACCESS_KEY\" \
  -s3-secret-key \"$S3_SECRET_KEY\" \
  -discord-webhook \"$DISCORD_WEBHOOK\" \
  -notify-success=\"$NOTIFY_SUCCESS\" \
  -notify-failure=\"$NOTIFY_FAILURE\""
```

Then run from cron:

```bash
#!/bin/bash
source /etc/pxBackupManager/config.sh
eval "$PXBM_CMD"
```

### Windows Configuration

Create `C:\Scripts\backup-config.ps1`:

```powershell
# Database Configuration
$DbHost = "localhost"
$DbPort = 3306
$DbUser = "backup"
$DbPass = "your_password"
$DbName = "fivem"

# Backup Configuration
$BackupDir = "D:\Backups\fivem"
$FiveMLoc = "C:\FiveM\server"
$Retention = 30

# S3 Configuration
$S3Enabled = $true
$S3Bucket = "my-backups"
$S3Region = "us-west-2"
$S3AccessKey = "YOUR_KEY"
$S3SecretKey = "YOUR_SECRET"

# Discord Configuration
$DiscordWebhook = "https://discord.com/api/webhooks/..."
$NotifySuccess = $true
$NotifyFailure = $true

# Build command
$BackupCmd = @(
    "C:\pxBackupManager\pxBackupManager.exe"
    "-fivem `"$FiveMLoc`""
    "-db-host `"$DbHost`""
    "-db-user `"$DbUser`""
    "-db-pass `"$DbPass`""
    "-db-name `"$DbName`""
    "-backup-dir `"$BackupDir`""
    "-retention $Retention"
    "-compress"
    "-s3-enabled"
    "-s3-bucket `"$S3Bucket`""
    "-s3-access-key `"$S3AccessKey`""
    "-s3-secret-key `"$S3SecretKey`""
    "-discord-webhook `"$DiscordWebhook`""
    "-notify-success"
    "-notify-failure"
) -join ' '

# Execute
Invoke-Expression $BackupCmd
```

---

## Advanced Examples

### Multi-Server Backup Setup

Backup from multiple servers to central storage:

```bash
#!/bin/bash

# Server 1
SERVER1_IP="192.168.1.100"
./pxBackupManager \
  -db-host "$SERVER1_IP" \
  -db-name "fivem" \
  -backup-dir "/backups/server1" \
  -s3-bucket "company-backups-server1"

# Server 2
SERVER2_IP="192.168.1.101"
./pxBackupManager \
  -db-host "$SERVER2_IP" \
  -db-name "fivem" \
  -backup-dir "/backups/server2" \
  -s3-bucket "company-backups-server2"

# Server 3
SERVER3_IP="192.168.1.102"
./pxBackupManager \
  -db-host "$SERVER3_IP" \
  -db-name "fivem" \
  -backup-dir "/backups/server3" \
  -s3-bucket "company-backups-server3"
```

### Tiered Backup Strategy

```bash
# Daily backups (kept for 7 days)
0 2 * * * /usr/local/bin/pxBackupManager -db-name "fivem" -backup-dir /backups/daily -retention 7

# Weekly backups (kept for 30 days)
0 3 * * 0 /usr/local/bin/pxBackupManager -db-name "fivem" -backup-dir /backups/weekly -retention 30

# Monthly backups (kept for 365 days)
0 4 1 * * /usr/local/bin/pxBackupManager -db-name "fivem" -backup-dir /backups/monthly -retention 365

# All to S3 for long-term archival
0 5 * * * /usr/local/bin/pxBackupManager -db-name "fivem" -s3-enabled -s3-bucket "longterm-archive" -s3-upload-fivem -s3-upload-mariadb
```

### High Availability Setup

Backup to local storage and S3, with notifications:

```bash
./pxBackupManager \
  -db-name "fivem" \
  -backup-dir /backups/local \
  -retention 14 \
  -compress \
  -s3-enabled \
  -s3-bucket "ha-backups" \
  -s3-access-key "$S3_KEY" \
  -s3-secret-key "$S3_SECRET" \
  -s3-upload-fivem \
  -s3-upload-mariadb \
  -discord-webhook "$WEBHOOK_1" \
  -notify-success \
  -notify-failure
```

---

## Next Steps

- See [Usage Guide](usage.md) for all command-line options
- See [Deployment Guide](deployment.md) for scheduling
- See [Troubleshooting Guide](troubleshooting.md) for common issues
