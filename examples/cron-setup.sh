#!/bin/bash

# Cron job setup for pxBackupManager
# Add this to your crontab with: crontab -e

# Daily backup at 2 AM
0 2 * * * /opt/pxBackupManager/pxBackupManager -fivem /opt/fivem -db-name fivem -db-user backup -db-pass YourPassword -backup-dir /backups/pxBackupManager >> /var/log/pxBackupManager.log 2>&1

# Backup every 6 hours
0 */6 * * * /opt/pxBackupManager/pxBackupManager -fivem /opt/fivem -db-name fivem -db-user backup -db-pass YourPassword -backup-dir /backups/pxBackupManager >> /var/log/pxBackupManager.log 2>&1

# Backup twice daily (2 AM and 2 PM)
0 2,14 * * * /opt/pxBackupManager/pxBackupManager -fivem /opt/fivem -db-name fivem -db-user backup -db-pass YourPassword -backup-dir /backups/pxBackupManager >> /var/log/pxBackupManager.log 2>&1
