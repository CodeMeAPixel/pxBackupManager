# Recommended Setup Guide

This guide explains the recommended directory structure and user setup for running pxBackupManager in production on Linux.

## Overview

The recommended setup uses:
- A dedicated non-root `backup` user for running the service
- Specific directories for organization and security
- systemd for scheduling and logging
- Proper file permissions to prevent unauthorized access

## Directory Structure

```
/opt/pxBackupManager/          # Binary location
  └─ pxBackupManager           # Executable

/backups/pxBackupManager/      # Backup storage
  ├─ fivem-backup-*.tar.gz
  └─ mariadb-backup-*.tar.gz

/var/log/pxBackupManager/      # Log files (via systemd journal)
```

## Why This Structure?

**Why not /home?**
- Backup user doesn't need a login shell or home directory
- Keeps system binaries separate from user files

**Why /opt?**
- Standard location for optional/third-party software
- Separates from system binaries in /usr/bin
- Organized and discoverable

**Why /backups?**
- Dedicated partition for backups (can be on separate disk)
- Easy to manage storage independently
- Standard convention across Linux systems

**Why /var/log?**
- systemd logs go to journalctl by default
- Keeps audit trails with system logs
- Easy to find and analyze backup history

## Step-by-Step Setup

### Quick Setup (Automated Script)

If you have already built or downloaded the binary, you can use the automated setup script:

```bash
# Clone the repository (if you haven't already)
git clone https://github.com/CodeMeAPixel/pxBackupManager.git
cd pxBackupManager

# Build the binary
make build

# Run the setup script as root
sudo bash examples/setup-systemd.sh
```

The script will:
- Create the `pxbackup` system user
- Create `/opt/pxBackupManager/` and `/backups/pxBackupManager/` directories
- Copy the binary to the correct location
- Install systemd service and timer files
- Set proper permissions

After the script completes, edit the service file to add your database credentials:

```bash
sudo nano /etc/systemd/system/pxBackupManager.service
```

Replace `YourPasswordHere` with your actual database password, then:

```bash
sudo systemctl daemon-reload
sudo systemctl start pxBackupManager.timer
```

### Manual Setup (Step-by-Step)

### 1. Create the Backup User

```bash
sudo useradd -r -s /bin/false -d /backups pxbackup
```

**Explanation:**
- `-r`: System account (UID < 1000)
- `-s /bin/false`: No login shell (can't SSH as this user)
- `-d /backups`: Home directory set to backups folder

Verify:
```bash
id pxbackup
```

### 2. Create Directory Structure

```bash
# Create all directories
sudo mkdir -p /opt/pxBackupManager
sudo mkdir -p /backups/pxBackupManager

# Set ownership to backup user
sudo chown -R pxbackup:pxbackup /backups/pxBackupManager
sudo chown -R pxbackup:pxbackup /opt/pxBackupManager

# Set restrictive permissions
sudo chmod 755 /opt/pxBackupManager
sudo chmod 700 /backups/pxBackupManager
```

**Why chmod 700 for /backups/pxBackupManager?**
- Only the `pxbackup` user can read/write/execute
- Other users cannot list or access backup files
- Sensitive data protection

### 3. Install the Binary

Build and copy:

```bash
# Build locally
make build

# Copy to production location
sudo cp pxBackupManager /opt/pxBackupManager/pxBackupManager
sudo chown pxbackup:pxbackup /opt/pxBackupManager/pxBackupManager
sudo chmod 755 /opt/pxBackupManager/pxBackupManager
```

Or install from release:

```bash
# Download from GitHub releases
wget https://github.com/CodeMeAPixel/pxBackupManager/releases/download/v2.0.0/pxBackupManager-linux-x64
sudo cp pxBackupManager /opt/pxBackupManager/
sudo chmod 755 /opt/pxBackupManager/pxBackupManager
```

### 4. Setup Systemd Service

Create `/etc/systemd/system/pxBackupManager.service`:

```bash
sudo tee /etc/systemd/system/pxBackupManager.service > /dev/null << 'EOF'
[Unit]
Description=pxBackupManager - Backup Service
After=network.target mysql.service

[Service]
Type=oneshot
User=pxbackup
Group=pxbackup
ExecStart=/opt/pxBackupManager/pxBackupManager \
  -fivem /opt/fivem \
  -db-name "your_database" \
  -db-user "backup" \
  -db-pass "your_password" \
  -backup-dir /backups/pxBackupManager \
  -retention 30 \
  -compress

StandardOutput=journal
StandardError=journal
SyslogIdentifier=pxBackupManager

[Install]
WantedBy=multi-user.target
EOF
```

**Important:** Replace the following in ExecStart:
- `your_database` - Your MariaDB database name
- `your_password` - Your database password
- Adjust `-fivem` path to your FiveM server location
- Adjust `-db-user` if not "backup"

### 5. Setup Systemd Timer

Create `/etc/systemd/system/pxBackupManager.timer`:

```bash
sudo tee /etc/systemd/system/pxBackupManager.timer > /dev/null << 'EOF'
[Unit]
Description=pxBackupManager - Daily Backup Timer
Requires=pxBackupManager.service

[Timer]
# Run daily at 2:00 AM
OnCalendar=*-*-* 02:00:00

# Run on boot if missed
Persistent=true

# Small random delay to prevent simultaneous system load
RandomizedDelaySec=30s

[Install]
WantedBy=timers.target
EOF
```

### 6. Enable and Start

```bash
# Reload systemd configuration
sudo systemctl daemon-reload

# Enable timer to start at boot
sudo systemctl enable pxBackupManager.timer

# Start the timer now
sudo systemctl start pxBackupManager.timer
```

### 7. Verify Setup

Check timer status:
```bash
sudo systemctl status pxBackupManager.timer
```

List scheduled runs:
```bash
sudo systemctl list-timers pxBackupManager.timer
```

View service logs:
```bash
sudo journalctl -u pxBackupManager.service -n 50
```

Watch live logs:
```bash
sudo journalctl -u pxBackupManager.service -f
```

Verify permissions:
```bash
ls -la /backups/pxBackupManager/
ls -la /opt/pxBackupManager/
```

## Configuration Examples

### Basic Setup

```bash
ExecStart=/opt/pxBackupManager/pxBackupManager \
  -fivem /opt/fivem \
  -db-name "fivem" \
  -db-user "backup" \
  -db-pass "SecurePassword123"
```

### With S3 Backup

```bash
ExecStart=/opt/pxBackupManager/pxBackupManager \
  -fivem /opt/fivem \
  -db-name "fivem" \
  -db-user "backup" \
  -db-pass "SecurePassword123" \
  -s3-enabled \
  -s3-bucket "my-backups" \
  -s3-access-key "YOUR_KEY" \
  -s3-secret-key "YOUR_SECRET" \
  -s3-upload-fivem \
  -s3-upload-mariadb
```

### With Discord Notifications

```bash
ExecStart=/opt/pxBackupManager/pxBackupManager \
  -fivem /opt/fivem \
  -db-name "fivem" \
  -db-user "backup" \
  -db-pass "SecurePassword123" \
  -discord-webhook "https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN" \
  -notify-success \
  -notify-failure
```

### Full Setup with All Options

```bash
ExecStart=/opt/pxBackupManager/pxBackupManager \
  -fivem /opt/fivem \
  -db-name "fivem" \
  -db-user "backup" \
  -db-pass "SecurePassword123" \
  -backup-dir /backups/pxBackupManager \
  -retention 30 \
  -compress \
  -s3-enabled \
  -s3-bucket "my-backups" \
  -s3-access-key "YOUR_KEY" \
  -s3-secret-key "YOUR_SECRET" \
  -s3-upload-fivem \
  -s3-upload-mariadb \
  -discord-webhook "https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN" \
  -notify-success \
  -notify-failure
```

## Scheduling Variations

### Every 6 Hours

```ini
[Timer]
OnCalendar=*-*-* 00,06,12,18:00:00
Persistent=true
```

### Twice Daily (2 AM and 2 PM)

```ini
[Timer]
OnCalendar=*-*-* 02,14:00:00
Persistent=true
```

### Weekly (Every Sunday at 2 AM)

```ini
[Timer]
OnCalendar=Sun *-*-* 02:00:00
Persistent=true
```

After changes, reload with:
```bash
sudo systemctl daemon-reload
sudo systemctl restart pxBackupManager.timer
```

## Security Best Practices

### 1. Secure the Configuration

The service file contains sensitive data:

```bash
# Already restrictive by default (644 or 640)
sudo chmod 640 /etc/systemd/system/pxBackupManager.service
sudo chown root:root /etc/systemd/system/pxBackupManager.service
```

### 2. Restrict Backup Directory

```bash
# Only pxbackup user can access
sudo chmod 700 /backups/pxBackupManager
sudo chown pxbackup:pxbackup /backups/pxBackupManager
```

### 3. Use Environment Variables (Optional)

Instead of hardcoding passwords, use environment variables:

Create `/etc/default/pxBackupManager`:

```bash
DB_USER=backup
DB_PASS=SecurePassword123
DB_NAME=fivem
S3_ACCESS_KEY=YOUR_KEY
S3_SECRET_KEY=YOUR_SECRET
DISCORD_WEBHOOK=https://discord.com/api/webhooks/...
```

Set proper permissions:
```bash
sudo chmod 600 /etc/default/pxBackupManager
sudo chown root:root /etc/default/pxBackupManager
```

Then source in service file:

```ini
EnvironmentFile=/etc/default/pxBackupManager
ExecStart=/opt/pxBackupManager/pxBackupManager \
  -fivem /opt/fivem \
  -db-name $DB_NAME \
  -db-user $DB_USER \
  -db-pass "$DB_PASS" \
  ...
```

### 4. Monitor Logs

Regular log review helps catch issues:

```bash
# Check for errors
sudo journalctl -u pxBackupManager.service -p err

# Check last 24 hours
sudo journalctl -u pxBackupManager.service --since "24 hours ago"

# Export logs to file
sudo journalctl -u pxBackupManager.service > pxbackup-logs.txt
```

## Troubleshooting

### Timer doesn't run

Check if enabled:
```bash
sudo systemctl is-enabled pxBackupManager.timer
# Should output: enabled
```

Check status:
```bash
sudo systemctl status pxBackupManager.timer
```

Check next run:
```bash
sudo systemctl list-timers pxBackupManager.timer
```

### Service fails with permission error

Check backup directory permissions:
```bash
sudo -u pxbackup ls -la /backups/pxBackupManager/
```

If it fails, fix permissions:
```bash
sudo chown -R pxbackup:pxbackup /backups/pxBackupManager
sudo chmod 700 /backups/pxBackupManager
```

### No logs showing up

Check if journalctl is working:
```bash
sudo journalctl -u pxBackupManager.service
```

If empty, try running the service manually:
```bash
sudo systemctl start pxBackupManager.service
sudo journalctl -u pxBackupManager.service -n 50
```

### Database connection refused

Verify credentials:
```bash
sudo -u pxbackup /opt/pxBackupManager/pxBackupManager \
  -db-name "fivem" \
  -db-user "backup" \
  -db-pass "YourPassword" \
  -only-mariadb
```

Check database user has backup privileges:
```bash
mysql -u backup -p
mysql> SHOW GRANTS;
```

### S3 upload fails

Test S3 connectivity:
```bash
sudo -u pxbackup /opt/pxBackupManager/pxBackupManager \
  -s3-enabled \
  -s3-bucket "test" \
  -s3-access-key "key" \
  -s3-secret-key "secret" \
  -only-fivem
```

Check credentials are correct and bucket exists.

## Uninstall

To remove the setup:

```bash
# Stop and disable timer
sudo systemctl stop pxBackupManager.timer
sudo systemctl disable pxBackupManager.timer

# Remove systemd files
sudo rm /etc/systemd/system/pxBackupManager.service
sudo rm /etc/systemd/system/pxBackupManager.timer
sudo systemctl daemon-reload

# Remove binary (optional)
sudo rm -rf /opt/pxBackupManager

# Keep backups (optional - decide based on your needs)
# sudo rm -rf /backups/pxBackupManager
```

## Next Steps

- See [Deployment Guide](deployment.md) for other scheduling options (cron, Docker, Windows)
- Read [Configuration Guide](configuration.md) for advanced setup
- Check [Troubleshooting Guide](troubleshooting.md) for common issues
- Review [Security Guide](../SECURITY.md) for security best practices
