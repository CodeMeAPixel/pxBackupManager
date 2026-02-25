# Features & Capabilities

## Core Features

### 1. FiveM Server Backups

**What it backs up:**
- Complete FiveM server directory with all scripts, resources, and configurations
- Preserves directory structure and file permissions
- Cross-platform compatible (Windows, Linux, macOS)

**How it works:**
- Uses Go's native `archive/tar` for compatibility
- Automatically handles path normalization
- Supports both compressed and uncompressed backups

**Example:**
```bash
./pxBackupManager -only-fivem -fivem /opt/fivem
```

### 2. MariaDB/MySQL Database Backups

**What it backs up:**
- Complete database schema and data
- Triggers, stored procedures, and views
- User privileges and access controls

**How it works:**
- Uses `mysqldump` for consistent backups
- Supports `--single-transaction` for InnoDB tables (live backups without locking)
- Compresses database dumps on-the-fly

**Example:**
```bash
./pxBackupManager -only-mariadb -db-name "fivem" -db-user "backup"
```

### 3. Compression

**What it does:**
- Reduces backup file sizes by 60-80%
- Saves storage costs for S3 and local storage
- Uses gzip compression standard

**Performance:**
- Default: Enabled
- Speed impact: ~20-30% slower backup time
- Space savings: Typically 70-80% reduction

**Example:**
```bash
./pxBackupManager -compress=true    # Enabled by default
./pxBackupManager -compress=false   # Disable for faster backups
```

### 4. Retention Policy

**What it does:**
- Automatically deletes backups older than specified days
- Prevents unlimited growth of backup storage
- Runs after successful backup

**Configuration:**
- Default: 30 days
- Set to 0 to disable cleanup
- Runs after each backup completion

**Example:**
```bash
./pxBackupManager -retention 30      # Keep 30 days (default)
./pxBackupManager -retention 90      # Keep 90 days
./pxBackupManager -retention 0       # Keep forever
./pxBackupManager -cleanup=false     # Disable cleanup
```

---

## Cloud Storage

### S3 Compatible Storage

**What it does:**
- Uploads backups to S3-compatible storage services
- Supports AWS S3, Backblaze B2, DigitalOcean Spaces, MinIO, etc.
- Parallel uploads for faster transfer

**Supported Services:**
- **AWS S3** - Standard cloud storage
- **Backblaze B2** - Cost-effective alternative to S3
- **DigitalOcean Spaces** - Integrated object storage
- **MinIO** - Self-hosted S3-compatible storage
- **Wasabi** - S3 API compatible
- **Oracle Object Storage** - OCI offering
- Any S3-compatible API endpoint

**Features:**
- Configurable custom endpoints
- Separate storage paths for different backup types
- Automatic retry on failure
- Server-side encryption support

**Example - AWS:**
```bash
./pxBackupManager \
  -db-name "fivem" \
  -s3-enabled \
  -s3-bucket "my-backups" \
  -s3-region "us-west-2" \
  -s3-access-key "AKIA..." \
  -s3-secret-key "..." \
  -s3-upload-fivem \
  -s3-upload-mariadb
```

**Example - Backblaze B2:**
```bash
./pxBackupManager \
  -db-name "fivem" \
  -s3-enabled \
  -s3-endpoint "https://s3.us-west-000.backblazeb2.com" \
  -s3-bucket "backups" \
  -s3-region "us-west-000" \
  -s3-access-key "app-key-id" \
  -s3-secret-key "app-key" \
  -s3-upload-fivem
```

---

## Notifications

### Discord Webhooks

**What it does:**
- Sends backup status notifications to Discord
- Rich embeds with detailed information
- Color-coded messages (green for success, red for failure)

**Information Included:**
- Backup status (success/failed)
- Timestamp of backup
- Duration per service
- File sizes
- S3 upload URLs (if applicable)
- Error messages (if failed)

**Notification Types:**
- **Success notifications** - Sent when all backups succeed
- **Failure notifications** - Sent when any backup fails
- Both types can be enabled independently

**Example:**
```bash
./pxBackupManager \
  -db-name "fivem" \
  -discord-webhook "https://discord.com/api/webhooks/ID/TOKEN" \
  -notify-success \
  -notify-failure
```

**Sample Discord Message:**
```
✓ Backup Completed Successfully

[fivem] ✓ Success - 2345ms
Size: 45.23 MB
[mariadb] ✓ Success - 1234ms
Size: 12.45 MB
S3: s3://my-bucket/fivem-backup-2024-01-15.tar.gz

Total: 2 succeeded, 0 failed
```

---

## Backup Naming & Organization

### File Naming Convention

```
{service}-backup-{YYYY-MM-DD-HHMMSS}.{extension}
```

**Examples:**
- `fivem-backup-2024-01-15-143022.tar.gz` - Compressed FiveM backup
- `mariadb-backup-2024-01-15-143022.tar.gz` - Compressed database backup
- `fivem-backup-2024-01-15-143022.tar` - Uncompressed FiveM backup

**Advantages:**
- Timestamp-based sorting
- Clear service identification
- Easy to identify by date

### Directory Organization

Example backup directory structure:
```
/backups/
├── fivem-backup-2024-01-15-020000.tar.gz          # Yesterday
├── fivem-backup-2024-01-14-020000.tar.gz          # 2 days ago
├── mariadb-backup-2024-01-15-020000.tar.gz        # Yesterday
├── mariadb-backup-2024-01-14-020000.tar.gz        # 2 days ago
└── ... (older backups)
```

---

## Cross-Platform Support

### Operating Systems

**Linux:**
- Ubuntu 18.04+
- Debian 10+
- CentOS 7+
- RHEL 8+
- Rocky Linux
- AlmaLinux
- Any systemd-based distribution

**macOS:**
- macOS 10.13+
- macOS 11+ (Big Sur and later)
- M1/M2 native support

**Windows:**
- Windows Server 2016+
- Windows 10+
- Windows 11
- WSL2 with Linux support

### Platform-Specific Features

**Linux:**
- Systemd service and timer support (recommended)
- Cron job scheduling
- journalctl logging integration
- Full permission control

**macOS:**
- launchd support for scheduling
- Cron job support
- Integration with system logs

**Windows:**
- Task Scheduler integration
- PowerShell script support
- Event Viewer logging
- UNC path support for network drives

---

## Performance Features

### Efficient Backup Creation

**Speed Optimizations:**
- Streaming compression (no intermediate files)
- Parallel file processing
- Minimal memory footprint
- Configurable compression levels

**Memory Usage:**
- Fixed memory allocation regardless of backup size
- Suitable for resource-constrained environments
- Stream processing instead of buffering

### Large Backup Support

**Tested with:**
- FiveM servers: 100GB+
- MariaDB databases: 50GB+
- S3 uploads: Multi-gigabyte files

**Limitations:**
- Only limited by available disk space
- Network transfer limited by bandwidth
- Upload timeouts configurable per service

---

## Scheduling & Automation

### Supported Scheduling Methods

1. **Systemd Timers** (Linux, recommended)
   - Flexible scheduling syntax
   - System integration
   - Automatic retry on failure

2. **Cron Jobs** (Unix-like systems)
   - Traditional, widely supported
   - Simple configuration
   - Compatible with older systems

3. **Windows Task Scheduler**
   - Native Windows scheduling
   - Integration with Event Viewer
   - Detailed execution history

4. **Docker**
   - Containerized deployment
   - Easy horizontal scaling
   - Consistent environments

5. **Manual Scripts**
   - Custom logic
   - Integration with other tools
   - Maximum flexibility

### Example Schedules

```bash
# Daily at 2 AM
OnCalendar=*-*-* 02:00:00

# Every 6 hours
OnCalendar=*-*-* 00,06,12,18:00:00

# Weekly on Monday
OnCalendar=Mon *-*-* 02:00:00

# Multiple times per day
OnCalendar=*-*-* 02,14:00:00
```

---

## Monitoring & Logging

### Log Output

**Console Output:**
- Real-time progress display
- Success/failure indicators
- File sizes and durations
- Error messages with details

**File Logging:**
- Compatible with syslog
- Integration with journalctl
- Structured error reporting

### Status Codes

- `0` - All backups successful
- `1` - One or more backups failed

---

## Security Features

### Built-in Security

1. **Credential Handling**
   - Command-line flags (process list visible)
   - Environment variables (more secure)
   - Configuration files (with permission control)

2. **Access Control**
   - File permission preservation (tar backups)
   - User/group settings for scheduled backups
   - Restricted backup directory permissions

3. **Transport Security**
   - SSL/TLS for S3 uploads
   - HTTPS for Discord webhooks
   - Standard cryptographic libraries

### Best Practices

- Use dedicated backup user account
- Restrict backup directory permissions (chmod 700)
- Store credentials in environment variables or secure files
- Rotate S3 access keys regularly
- Use Discord webhook tokens sparingly
- Monitor backup logs for errors

---

## Limitations & Considerations

### Known Limitations

1. **Database Backups**
   - `mysqldump` may lock tables (mitigated with `--single-transaction`)
   - Point-in-time requires additional binary logs
   - Large databases may require additional memory

2. **Large Backups**
   - Very large backups (>100GB) may take hours
   - Network bandwidth affects S3 upload times
   - Disk I/O may be bottlenecked on slower drives

3. **Retention Policy**
   - Only based on modification time
   - May require manual cleanup for corrupted files
   - Disabled backups not cleaned up automatically

### Performance Considerations

- **CPU**: Compression adds ~20-30% CPU overhead
- **Disk I/O**: Read performance depends on source storage speed
- **Network**: S3 upload speed limited by connection bandwidth
- **Memory**: Fixed allocation, minimal impact even for large backups

---

## Next Steps

- See [Installation Guide](installation.md) for setup
- See [Configuration Guide](configuration.md) for examples
- See [Deployment Guide](deployment.md) for scheduling
