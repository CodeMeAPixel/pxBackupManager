#!/bin/bash

# Setup script for systemd service and timer
# Run as root: sudo bash setup-systemd.sh

set -e

echo "====== pxBackupManager Systemd Setup ======"

# Create backup user
echo "Creating backup user..."
if ! id -u pxbackup > /dev/null 2>&1; then
    useradd -r -s /bin/false -d /backups pxbackup
    echo "Created user 'pxbackup'"
else
    echo "User 'pxbackup' already exists"
fi

# Create directories
echo "Creating directories..."
mkdir -p /opt/pxBackupManager
mkdir -p /backups/pxBackupManager

# Set permissions
chown -R pxbackup:pxbackup /backups/pxBackupManager
chown -R pxbackup:pxbackup /opt/pxBackupManager
chmod 755 /opt/pxBackupManager
chmod 700 /backups/pxBackupManager

# Copy binary
echo "Installing pxBackupManager binary..."
if [ -f "./pxBackupManager" ]; then
    cp ./pxBackupManager /opt/pxBackupManager/pxBackupManager
    chmod +x /opt/pxBackupManager/pxBackupManager
    echo "Binary installed to /opt/pxBackupManager/pxBackupManager"
else
    echo "ERROR: pxBackupManager binary not found in current directory"
    exit 1
fi

# Copy systemd files
echo "Installing systemd service and timer..."
cp systemd/pxBackupManager.service /etc/systemd/system/
cp systemd/pxBackupManager.timer /etc/systemd/system/

# Set proper permissions
chmod 644 /etc/systemd/system/pxBackupManager.service
chmod 644 /etc/systemd/system/pxBackupManager.timer

# Reload systemd daemon
echo "Reloading systemd daemon..."
systemctl daemon-reload

# Enable services
echo "Enabling pxBackupManager timer..."
systemctl enable pxBackupManager.timer

echo ""
echo "====== Setup Complete ======"
echo ""
echo "Next steps:"
echo "1. Edit /etc/systemd/system/pxBackupManager.service to set database credentials:"
echo "   sudo nano /etc/systemd/system/pxBackupManager.service"
echo ""
echo "2. The timer runs daily at 2 AM. To modify the schedule, edit:"
echo "   sudo nano /etc/systemd/system/pxBackupManager.timer"
echo ""
echo "3. Start the timer:"
echo "   sudo systemctl start pxBackupManager.timer"
echo ""
echo "4. Check status:"
echo "   sudo systemctl status pxBackupManager.timer"
echo "   sudo systemctl list-timers pxBackupManager.timer"
echo ""
echo "5. View logs:"
echo "   sudo journalctl -u pxBackupManager.service -f"
echo ""
