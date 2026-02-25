# Troubleshooting & Monitoring Guide

## Monitoring

### Systemd Service Monitoring

```bash
# Check service status
sudo systemctl status pxBackupManager.timer
sudo systemctl status pxBackupManager.service

# View last 50 lines
sudo journalctl -u pxBackupManager.service -n 50

# Follow logs in real-time
sudo journalctl -u pxBackupManager.service -f

# View logs for specific time
sudo journalctl -u pxBackupManager.service --since "2024-01-15 02:00:00"

# See all timer runs
sudo journalctl -u pxBackupManager.service | grep -E "failed|success"
```

### Cron Job Monitoring

```bash
# View recent backups
ls -lh /backups/pxbackup | tail -10

# Monitor backup.log
tail -f /var/log/pxBackupManager/backup.log

# Check cron logs
grep pxBackupManager /var/log/cron
```

### Disk Space Monitoring

```bash
# Check backup directory size
du -sh /backups/pxbackup

# List backups by size (largest first)
ls -lhS /backups/pxbackup | head -20

# Find backups older than 30 days
find /backups/pxbackup -type f -mtime +30 -ls

# Calculate disk usage for specific month
du -sh /backups/pxbackup/mariadb-backup-2024-01*
```

---

## Common Issues

### Issue: Database Connection Failure

```
ERROR: Failed to connect to database: (connection refused)
```

**Causes:**
- MariaDB/MySQL service not running
- Wrong host/port
- Wrong credentials

**Solutions:**

```bash
# Check if database is running
sudo systemctl status mariadb
# or
sudo systemctl status mysql

# Start database if stopped
sudo systemctl start mariadb

# Test connection manually
mysql -h localhost -u backup -p"your_password" -e "SELECT 1;"

# Check if port is listening
sudo netstat -tlnp | grep 3306
# or
sudo ss -tlnp | grep 3306

# Verify credentials
mysql -u backup -p"your_password" -e "SHOW DATABASES;"
```

**Verify flags:**
- `-db-host` should match database host
- `-db-port` should match database port
- `-db-user` should have backup privileges
- `-db-pass` should be correct

### Issue: FiveM Directory Not Found

```
ERROR: FiveM location does not exist: /opt/fivem
```

**Causes:**
- Wrong path specified
- Directory doesn't exist
- Permission denied

**Solutions:**

```bash
# Find FiveM directory
find / -type d -name "fivem" 2>/dev/null

# Check if directory exists
ls -la /opt/fivem

# Verify permissions (user running backup should have read access)
ls -ld /opt/fivem

# Fix permissions if needed
sudo chmod 755 /opt/fivem
```

### Issue: Insufficient Disk Space

```
ERROR: Not enough space in backup directory
```

**Causes:**
- Backup directory full
- Retention policy not cleaning up
- Disk drive full

**Solutions:**

```bash
# Check disk space
df -h /backups

# Calculate backup size
du -sh /backups/pxbackup

# Manually clean old backups (example: 30 days)
find /backups/pxbackup -type f -mtime +30 -delete

# Check for large files
du -sh /backups/pxbackup/* | sort -rh | head

# Expand volume (Linux)
sudo lvextend -L +100G /dev/mapper/vg0-backups
sudo resize2fs /dev/mapper/vg0-backups
```

### Issue: mysqldump Command Not Found

```
ERROR: mysqldump command failed: command not found
```

**Causes:**
- MySQL client not installed
- Not in system PATH
- Wrong installation location

**Solutions:**

**Linux:**
```bash
# Install MySQL client
sudo apt-get install mysql-client       # Ubuntu/Debian
sudo yum install mysql                  # CentOS/RHEL
sudo pacman -S mysql-clients            # Arch Linux

# Verify installation
which mysqldump
mysqldump --version
```

**macOS:**
```bash
# Install via Homebrew
brew install mysql-client

# Add to PATH if needed
export PATH="/usr/local/opt/mysql-client/bin:$PATH"
```

**Windows:**
```powershell
# Check if installed
Get-Command mysqldump -ErrorAction SilentlyContinue

# Add to PATH manually in PowerShell
$env:Path += ";C:\Program Files\MySQL\MySQL Server 8.0\bin"

# Verify
mysqldump --version
```

### Issue: S3 Upload Fails

```
ERROR: failed to upload to S3: NoSuchBucket
```

**Causes:**
- Bucket doesn't exist
- Wrong credentials
- Access denied
- Endpoint URL wrong (for S3-compatible services)

**Solutions:**

```bash
# Verify AWS credentials
aws s3 ls --profile default

# List buckets
aws s3 ls

# Check bucket permissions
aws s3api head-bucket --bucket my-bucket

# Test with AWS CLI first
aws s3 cp test.txt s3://my-bucket/test.txt

# For Backblaze B2:
# 1. Verify endpoint URL
# 2. Check access key and app key
# 3. Verify bucket name

# Test S3 connection manually
aws s3 cp /backups/pxbackup/test-backup.tar.gz s3://my-bucket/ \
  --endpoint-url https://s3.us-west-000.backblazeb2.com
```

**Check flags:**
- `-s3-bucket` - Correct bucket name
- `-s3-endpoint` - Correct for S3-compatible services
- `-s3-access-key` - Valid access key
- `-s3-secret-key` - Valid secret key
- `-s3-region` - Matches bucket region

### Issue: Discord Webhook Notifications Not Send

```
WARNING: failed to send Discord notification
```

**Causes:**
- Invalid webhook URL
- Webhook deleted/revoked
- Network connectivity issue
- Discord API rate limiting

**Solutions:**

```bash
# Test webhook manually with curl
curl -X POST "https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"Test notification"}'

# Recreate webhook in Discord:
# 1. Go to Server Settings -> Integrations -> Webhooks
# 2. Create New Webhook
# 3. Copy full URL (includes both ID and token)

# Verify webhook still exists
curl -X GET "https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN"
```

### Issue: Permission Denied

```
ERROR: permission denied: /backups/pxbackup
```

**Causes:**
- User running backup doesn't have permissions
- Wrong ownership
- Restrictive permissions

**Solutions:**

```bash
# Check current permissions
ls -ld /backups/pxbackup

# Fix ownership
sudo chown -R backup:backup /backups/pxbackup

# Fix permissions
sudo chmod 700 /backups/pxbackup

# If running via systemd, verify User in service file
sudo grep "^User=" /etc/systemd/system/pxBackupManager.service

# If running via cron, verify crontab user
sudo crontab -u backup -l
```

### Issue: Backup Hangs/Takes Too Long

**Causes:**
- Large FiveM directory
- Large database
- Slow disk I/O
- Network issues (for S3 uploads)

**Solutions:**

```bash
# Monitor in real-time
watch -n 1 'ls -lh /backups/pxbackup | tail -5'

# Check disk I/O
iostat -x 1 10

# Check system resources
top
htop

# Disable compression to speed up (trades disk space for speed)
# Add -compress=false to backup command

# Increase timeout for systemd service
# Edit service file and add:
# TimeoutStartSec=3600

# Check database size before backup
sudo mysql -u root -p"password" -e "SELECT table_schema, SUM(data_length + index_length) FROM information_schema.tables GROUP BY table_schema;"
```

### Issue: Backups Not Retention Policy Not Working

```
Backup retention period set but old backups not deleted
```

**Causes:**
- Cleanup disabled
- Wrong retention value
- File timestamps not accurate

**Solutions:**

```bash
# Check retention setting in command
# Verify -cleanup is not disabled
# Verify -retention is set correctly

# Manual cleanup (example: 30 days)
find /backups/pxbackup -type f -mtime +30 -delete

# Check file modification times
ls -la --time-style=long-iso /backups/pxbackup/

# Debug cleanup
# Run manually to see output
./pxBackupManager -db-name "test" -cleanup -retention 30
```

### Issue: Service Not Auto-Starting After Reboot

**Causes:**
- Timer not enabled
- Service dependencies not met
- Startup scripts not loaded

**Solutions:**

```bash
# Check if timer is enabled
sudo systemctl is-enabled pxBackupManager.timer

# Enable timer
sudo systemctl enable pxBackupManager.timer

# Check timer status
sudo systemctl list-timers pxBackupManager.timer

# Check service dependencies
sudo systemctl show -p After pxBackupManager.service
sudo systemctl show -p Requires pxBackupManager.service

# Check for startup errors
sudo journalctl -u pxBackupManager -b
```

---

## Performance Optimization

### Reduce Backup Time

1. **Disable compression if disk space permits:**
   ```bash
   -compress=false
   ```

2. **Run during off-hours:**
   Schedule backups at times when server load is lowest

3. **Use faster storage:**
   Store backups on SSD instead of mechanical disk

4. **Optimize database:**
   ```bash
   mysql -u backup -p"password" -e "OPTIMIZE TABLE your_table;"
   ```

5. **Exclude unnecessary files:**
   Consider backing up only essential directories

### Reduce Storage Usage

1. **Enable compression:**
   ```bash
   -compress=true  # default
   ```

2. **Reduce retention period:**
   ```bash
   -retention 15  # Keep only 15 days instead of 30
   ```

3. **Upload to S3 and delete local:**
   Keep only recent backups locally, older ones on S3

4. **Exclude cache/temp:**
   Backup only critical data

### Monitoring Storage

```bash
# Set up daily storage check
0 5 * * * du -sh /backups/pxbackup | mail -s "Backup Storage Report" admin@example.com

# Create alert if above threshold
#!/bin/bash
SIZE=$(du -s /backups/pxbackup | awk '{print $1}')
LIMIT=$((500 * 1024 * 1024))  # 500GB in KB

if [ $SIZE -gt $LIMIT ]; then
  echo "Backup storage exceeds limit!" | mail -s "ALERT: Backup Storage" admin@example.com
fi
```

---

## Next Steps

- See [Installation Guide](installation.md) for setup
- See [Usage Guide](usage.md) for command-line options
- See [Deployment Guide](deployment.md) for scheduling
