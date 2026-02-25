# Deployment & Scheduling Guide

This guide covers multiple ways to schedule and run pxBackupManager automatically on your server.

## Table of Contents

1. [Systemd Service + Timer (Recommended)](#option-1-systemd-service--timer-recommended)
2. [Cron Jobs (Traditional Linux)](#option-2-cron-jobs-traditional-linux)
3. [Docker (Container Deployment)](#option-3-docker-container-deployment)
4. [Windows Task Scheduler](#option-4-windows-task-scheduler)
5. [Manual Recurring Script](#option-5-manual-recurring-script)

---

## Option 1: Systemd Service + Timer (Recommended)

**Best for:** Modern Linux distributions (systemd-based)

**Advantages:**
- Better logging integration with journalctl
- Dependency management
- Automatic restart on failure
- Easy scheduling with timers
- System-wide availability

**Before starting:** Read [Recommended Setup](recommended-setup.md) for production directory structure and user permissions.

### Quick Start

The [Recommended Setup Guide](recommended-setup.md) provides detailed steps for:
- Creating a dedicated `pxbackup` user
- Setting up the directory structure
- Installing the binary
- Creating systemd service and timer files
- Verifying the setup

Follow that guide for a production-ready installation.

### Manual Setup (Advanced)

If you prefer manual configuration:

1. **Build the binary:**
   ```bash
   make build
   ```

2. **Copy binary to standard location:**
   ```bash
   sudo cp pxBackupManager /usr/local/bin/
   sudo chmod +x /usr/local/bin/pxBackupManager
   ```

3. **Create systemd service file:**
   ```bash
   sudo tee /etc/systemd/system/pxBackupManager.service > /dev/null << EOF
   [Unit]
   Description=pxBackupManager - Backup Service
   After=network.target

   [Service]
   Type=oneshot
   User=backup
   Group=backup
   ExecStart=/usr/local/bin/pxBackupManager \
     -fivem /opt/fivem \
     -db-name "fivem" \
     -db-user "backup" \
     -db-pass "your_password" \
     -backup-dir /backups/pxbackup \
     -retention 30 \
     -compress \
     -notify-success \
     -discord-webhook "https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN"

   StandardOutput=journal
   StandardError=journal

   [Install]
   WantedBy=multi-user.target
   EOF
   ```

4. **Create systemd timer file:**
   ```bash
   sudo tee /etc/systemd/system/pxBackupManager.timer > /dev/null << EOF
   [Unit]
   Description=pxBackupManager - Daily Backup Timer
   Requires=pxBackupManager.service

   [Timer]
   OnCalendar=daily
   OnCalendar=*-*-* 02:00:00
   OnBootSec=5min
   Persistent=true

   [Install]
   WantedBy=timers.target
   EOF
   ```

5. **Create dedicated backup user:**
   ```bash
   sudo useradd -r -s /bin/false backup
   ```

6. **Setup backup directory:**
   ```bash
   sudo mkdir -p /backups/pxbackup
   sudo chown -R backup:backup /backups/pxbackup
   sudo chmod 700 /backups/pxbackup
   ```

7. **Enable and start the timer:**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable pxBackupManager.timer
   sudo systemctl start pxBackupManager.timer
   ```

### Managing the Timer

```bash
# Check timer status
sudo systemctl status pxBackupManager.timer

# List all timers
sudo systemctl list-timers pxBackupManager.timer

# View service logs
sudo journalctl -u pxBackupManager.service -n 50 -f

# Manually trigger backup now
sudo systemctl start pxBackupManager.service

# Edit timer schedule
sudo systemctl edit --full pxBackupManager.timer

# Disable timer (keeps it installed)
sudo systemctl stop pxBackupManager.timer
sudo systemctl disable pxBackupManager.timer
```

### Timer Schedule Examples

```bash
# Every day at 2 AM
OnCalendar=*-*-* 02:00:00

# Every 6 hours
OnCalendar=*-*-* 00,06,12,18:00:00

# Twice daily (2 AM and 2 PM)
OnCalendar=*-*-* 02,14:00:00

# Every Monday at 2 AM
OnCalendar=Mon *-*-* 02:00:00

# Every 12 hours
OnCalendar=*-*-* 00:00:00
OnCalendar=*-*-* 12:00:00
```

---

## Option 2: Cron Jobs (Traditional Linux)

**Best for:** Older Linux systems without systemd

**Advantages:**
- Simple and lightweight
- Compatible with all Unix-like systems
- Easy to understand and modify

### Setup Steps

1. **Build the binary:**
   ```bash
   make build
   ```

2. **Install to standard location:**
   ```bash
   sudo mkdir -p /opt/pxBackupManager
   sudo cp pxBackupManager /opt/pxBackupManager/
   sudo chmod +x /opt/pxBackupManager/pxBackupManager
   ```

3. **Create log directory:**
   ```bash
   sudo mkdir -p /var/log/pxBackupManager
   sudo chmod 755 /var/log/pxBackupManager
   ```

4. **Create cron script:**
   ```bash
   sudo tee /opt/pxBackupManager/backup.sh > /dev/null << 'EOF'
   #!/bin/bash
   /opt/pxBackupManager/pxBackupManager \
     -fivem /opt/fivem \
     -db-name "fivem" \
     -db-user "backup" \
     -db-pass "your_password" \
     -backup-dir /backups/pxbackup \
     -retention 30 \
     -compress \
     >> /var/log/pxBackupManager/backup.log 2>&1
   EOF
   
   sudo chmod +x /opt/pxBackupManager/backup.sh
   ```

5. **Add cron job:**
   ```bash
   sudo crontab -e
   ```

   Add this line (daily at 2 AM):
   ```cron
   0 2 * * * /opt/pxBackupManager/backup.sh
   ```

### Cron Schedule Examples

```cron
# Daily at 2 AM
0 2 * * * /opt/pxBackupManager/backup.sh

# Every 6 hours
0 */6 * * * /opt/pxBackupManager/backup.sh

# Twice daily (2 AM and 2 PM)
0 2,14 * * * /opt/pxBackupManager/backup.sh

# Weekly on Monday at 2 AM
0 2 * * 1 /opt/pxBackupManager/backup.sh

# Every 30 minutes
*/30 * * * * /opt/pxBackupManager/backup.sh
```

### Monitoring Cron

```bash
# View cron logs (Ubuntu/Debian)
sudo grep CRON /var/log/syslog | tail -20

# View cron logs (CentOS/RHEL)
sudo grep CRON /var/log/cron | tail -20

# View backup logs
sudo tail -f /var/log/pxBackupManager/backup.log
```

---

## Option 3: Docker (Container Deployment)

**Best for:** Container-based infrastructure

### Create Dockerfile

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o pxBackupManager

FROM alpine:latest

RUN apk add --no-cache mysql-client

WORKDIR /app
COPY --from=builder /app/pxBackupManager .

ENTRYPOINT ["./pxBackupManager"]
```

### Create docker-compose.yml

```yaml
version: '3.8'

services:
  pxBackupManager:
    build: .
    container_name: pxBackupManager
    environment:
      DB_HOST: mariadb
      DB_NAME: fivem
      DB_USER: backup
      DB_PASSWORD: your_password
    volumes:
      - /opt/fivem:/fivem:ro
      - /backups:/backups
    networks:
      - backup-network
    restart: unless-stopped
    # Run daily at 2 AM
    entrypoint: >
      sh -c "
      apk add --no-cache dcron &&
      echo '0 2 * * * ./pxBackupManager -fivem /fivem -db-name fivem -db-user backup -db-pass $$DB_PASSWORD -backup-dir /backups -retention 30 -compress' > /etc/crontabs/root &&
      crond -f
      "

  mariadb:
    image: mariadb:latest
    container_name: mariadb
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: fivem
      MYSQL_USER: backup
      MYSQL_PASSWORD: your_password
    volumes:
      - mariadb_data:/var/lib/mysql
    networks:
      - backup-network
    restart: unless-stopped

volumes:
  mariadb_data:

networks:
  backup-network:
    driver: bridge
```

### Run Docker

```bash
# Build and start
docker-compose up -d

# View logs
docker logs -f pxBackupManager

# Execute manual backup
docker-compose exec pxBackupManager ./pxBackupManager -db-name fivem -db-user backup -db-pass your_password

# Stop services
docker-compose down
```

---

## Option 4: Windows Task Scheduler

**Best for:** Windows Server deployments

### Setup Steps

1. **Build the executable:**
   ```powershell
   make build
   ```

2. **Create backup script:**
   Create `C:\Scripts\backup.ps1`:
   ```powershell
   & "C:\pxBackupManager\pxBackupManager.exe" `
     -fivem "C:\FiveM\server" `
     -db-name "fivem" `
     -db-user "backup" `
     -db-pass "your_password" `
     -backup-dir "D:\Backups\pxBackup" `
     -retention 30 `
     -compress
   ```

3. **Create Task Scheduler task:**
   
   Open Task Scheduler and create a new task:
   
   - **General Tab:**
     - Name: "pxBackupManager Daily Backup"
     - Run with highest privileges
   
   - **Triggers Tab:**
     - New Trigger
     - Begin the task: On a schedule
     - Daily at 02:00 AM
   
   - **Actions Tab:**
     - Program/script: `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`
     - Add arguments: `-NoProfile -ExecutionPolicy Bypass -File "C:\Scripts\backup.ps1"`
     - Start in: `C:\Scripts`
   
   - **Conditions Tab:**
     - Uncheck "Wake the computer to run this task" if desired
   
   - **Settings Tab:**
     - Check "Run task as soon as possible after a scheduled start is missed"
     - Multiple instances: "Do not start a new instance"

4. **Test the task:**
   ```powershell
   # Run the task manually
   Start-ScheduledTask -TaskName "pxBackupManager Daily Backup"
   
   # Check task history
   Get-ScheduledTaskInfo -TaskName "pxBackupManager Daily Backup"
   ```

---

## Option 5: Manual Recurring Script

**Best for:** Simple, script-based deployments

### Create backup script

Unix/Linux/macOS:
```bash
#!/bin/bash

set -e

BACKUP_BIN="/opt/pxBackupManager/pxBackupManager"
LOG_FILE="/var/log/pxBackupManager/backup.log"

echo "[$(date +'%Y-%m-%d %H:%M:%S')] Starting backup..." >> $LOG_FILE

$BACKUP_BIN \
  -fivem /opt/fivem \
  -db-name "fivem" \
  -db-user "backup" \
  -db-pass "your_password" \
  -backup-dir /backups/pxbackup \
  -retention 30 \
  -compress \
  -notify-success \
  -discord-webhook "YOUR_WEBHOOK_URL" \
  >> $LOG_FILE 2>&1

if [ $? -eq 0 ]; then
  echo "[$(date +'%Y-%m-%d %H:%M:%S')] Backup completed successfully" >> $LOG_FILE
else
  echo "[$(date +'%Y-%m-%d %H:%M:%S')] Backup failed!" >> $LOG_FILE
  exit 1
fi
```

Windows:
```powershell
$BackupBin = "C:\pxBackupManager\pxBackupManager.exe"
$LogFile = "C:\Logs\pxBackupManager.log"

Add-Content $LogFile "$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') Starting backup..."

& $BackupBin `
  -fivem "C:\FiveM\server" `
  -db-name "fivem" `
  -db-user "backup" `
  -db-pass "your_password" `
  -backup-dir "D:\Backups\pxBackup" `
  -retention 30 `
  -compress

if ($LASTEXITCODE -eq 0) {
  Add-Content $LogFile "$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') Backup completed successfully"
} else {
  Add-Content $LogFile "$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') Backup failed!"
  exit 1
}
```

---

## Monitoring Backups

### Check Recent Backups

```bash
# List backups by date
ls -lh /backups/pxbackup | grep -E "\.tar\.gz|\.sql"

# Sort by modification time (newest first)
ls -lht /backups/pxbackup | head -10

# Calculate total backup storage
du -sh /backups/pxbackup
```

### Verify Backup Integrity

```bash
# List tar contents
tar -tzf /backups/pxbackup/fivem-backup-*.tar.gz | head -20

# Verify tar file
tar -tzf /backups/pxbackup/fivem-backup-*.tar.gz > /dev/null && echo "OK" || echo "CORRUPTED"

# List SQL dump
head -50 /backups/pxbackup/mariadb-backup-*.tar.gz | gunzip | head -20
```

---

## Backup Recovery

### Restore FiveM from Backup

```bash
# Extract to temporary directory first
mkdir /tmp/fivem-restore
tar -xzf /backups/pxbackup/fivem-backup-2024-01-15-143022.tar.gz -C /tmp/fivem-restore

# Review contents
ls -la /tmp/fivem-restore

# Restore (backup original first!)
cp -r /tmp/fivem-restore/fivem /opt/fivem.backup
cp -r /tmp/fivem-restore/fivem /opt/fivem
```

### Restore MariaDB from Backup

```bash
# Decompress if needed
tar -xzf /backups/pxbackup/mariadb-backup-2024-01-15-143022.tar.gz -O | mysql -u root -p

# Or directly import
mysql -u root -p fivem < /backups/pxbackup/mariadb-backup-2024-01-15-143022.sql
```

---

## Next Steps

- See [Usage Guide](usage.md) for command-line options
- See [Configuration Guide](configuration.md) for advanced setups
- See [Troubleshooting Guide](troubleshooting.md) for common issues
