#!/bin/bash

# Setup script for cron jobs
# Run as the user that will run backups: bash setup-cron.sh

set -e

echo "====== Flux Backup Cron Setup ======"

# Ensure backup log directory exists
mkdir -p ~/.backup-logs

# Generate crontab entry
CRON_JOB="0 2 * * * /opt/flux-backup/flux-backup -fivem /opt/fivem -db-name fivem -db-user backup -db-pass YourPassword -backup-dir /backups/flux >> ~/.backup-logs/flux-backup.log 2>&1"

# Check if cron entry already exists
if crontab -l 2>/dev/null | grep -q "flux-backup"; then
    echo "Backup cron job already exists in crontab"
    echo "Current cron jobs:"
    crontab -l | grep flux-backup
else
    echo "Adding cron job to crontab..."
    (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -
    echo "Cron job added successfully!"
    echo ""
    echo "New cron job:"
    crontab -l | grep flux-backup
fi

echo ""
echo "====== Setup Complete ======"
echo ""
echo "The backup will run daily at 2 AM"
echo ""
echo "To view logs:"
echo "  tail -f ~/.backup-logs/flux-backup.log"
echo ""
echo "To edit cron schedule:"
echo "  crontab -e"
echo ""
echo "To remove the backup job:"
echo "  crontab -e"
echo "  (and delete the flux-backup line)"
echo ""
