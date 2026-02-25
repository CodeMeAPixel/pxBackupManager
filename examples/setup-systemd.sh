#!/bin/bash

# Setup script for systemd service and timer
# Run as root: sudo bash setup-systemd.sh

set -e

echo "====== Flux Backup Systemd Setup ======"

# Create backup user
echo "Creating backup user..."
if ! id -u backup > /dev/null 2>&1; then
    useradd -r -s /bin/bash -d /backups backup
    echo "Created user 'backup'"
else
    echo "User 'backup' already exists"
fi

# Create directories
echo "Creating directories..."
mkdir -p /opt/flux-backup
mkdir -p /backups/flux
mkdir -p /var/log/flux-backup

# Set permissions
chown -R backup:backup /backups
chown -R backup:backup /var/log/flux-backup
chmod 755 /opt/flux-backup

# Copy binary
echo "Installing flux-backup binary..."
if [ -f "./flux-backup" ]; then
    cp ./flux-backup /opt/flux-backup/flux-backup
    chmod +x /opt/flux-backup/flux-backup
    echo "Binary installed to /opt/flux-backup/flux-backup"
else
    echo "ERROR: flux-backup binary not found in current directory"
    exit 1
fi

# Copy systemd files
echo "Installing systemd service and timer..."
cp systemd/flux-backup.service /etc/systemd/system/
cp systemd/flux-backup.timer /etc/systemd/system/

# Set proper permissions
chmod 644 /etc/systemd/system/flux-backup.service
chmod 644 /etc/systemd/system/flux-backup.timer

# Reload systemd daemon
echo "Reloading systemd daemon..."
systemctl daemon-reload

# Enable services
echo "Enabling flux-backup timer..."
systemctl enable flux-backup.timer

echo ""
echo "====== Setup Complete ======"
echo ""
echo "Next steps:"
echo "1. Edit /etc/systemd/system/flux-backup.service to set database credentials:"
echo "   sudo nano /etc/systemd/system/flux-backup.service"
echo ""
echo "2. The timer runs daily at 2 AM. To modify the schedule, edit:"
echo "   sudo nano /etc/systemd/system/flux-backup.timer"
echo ""
echo "3. Start the timer:"
echo "   sudo systemctl start flux-backup.timer"
echo ""
echo "4. Check status:"
echo "   sudo systemctl status flux-backup.timer"
echo "   sudo systemctl list-timers flux-backup.timer"
echo ""
echo "5. View logs:"
echo "   sudo journalctl -u flux-backup.service -f"
echo ""
