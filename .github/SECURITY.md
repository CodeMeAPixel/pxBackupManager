# Security Policy

## Reporting a Vulnerability

pxBackupManager takes security seriously. If you discover a security vulnerability, please report it responsibly by emailing **hey@codemeapixel.dev** instead of using the public issue tracker.

### What to Include in a Security Report

When reporting a security vulnerability, please provide:

1. **Description** - A clear description of the vulnerability
2. **Location** - Which file(s) or component(s) are affected
3. **Severity** - Your assessment of the severity (Critical, High, Medium, Low)
4. **Steps to Reproduce** - Instructions on how to reproduce the vulnerability
5. **Proof of Concept** - Code or detailed steps demonstrating the vulnerability (optional but helpful)
6. **Potential Impact** - Description of the potential harm

### Response Timeline

- **24 hours**: Initial acknowledgment of your report
- **7 days**: Initial assessment and proposed fix timeline
- **30 days**: Target for publishing a security patch (may vary depending on complexity)

## Security Considerations

### Credentials and Sensitive Data

- **Never commit credentials** to version control (database passwords, API keys, S3 secrets)
- Use environment variables or secure configuration files for sensitive data
- Ensure backup directories have appropriate file permissions (e.g., `chmod 700`)
- Store database passwords securely and restrict access to configuration files

### Database Backups

- Store backups on a separate disk/filesystem from the server data
- Implement proper access controls on backup directories
- Consider encrypting backups at rest
- Verify backup integrity periodically

### S3 Configuration

- Use IAM users/roles with minimal required permissions for S3 access
- Never share S3 credentials in logs or configuration files
- Use S3 bucket policies to restrict access
- Enable S3 server-side encryption for backups at rest
- Consider using temporary credentials with STS tokens

### Discord Webhooks

- Discord webhook URLs contain sensitive tokens
- Protect webhook URLs in configuration files
- Rotate webhook URLs periodically if they're exposed
- Use Discord permissions carefully when creating webhooks

### User Access Control

- Create a dedicated backup user account with minimal privileges
- Restrict cron/systemd service permissions
- Use file permissions to control who can read/modify backups
- Implement proper authentication if backing up from multiple servers

## Running Securely

### Linux/Unix Best Practices

```bash
# Create dedicated backup user
sudo useradd -r -s /bin/false backup

# Set restrictive permissions on backup directory
sudo mkdir -p /backups
sudo chown backup:backup /backups
sudo chmod 700 /backups

# Run backup tool with reduced privileges
sudo -u backup /path/to/pxBackupManager -db-name "your_db"
```

### Windows Best Practices

- Run the backup tool with a dedicated service account
- Use Windows Task Scheduler with restricted permissions
- Store configuration in files with restricted NTFS permissions
- Use Windows credential manager for sensitive data

## Known Limitations

- Backups are point-in-time snapshots; they may not include data written between backup intervals
- MariaDB `mysqldump` locks tables during backup (use `--single-transaction` for InnoDB)
- Large backups may consume significant disk space
- Network bandwidth may be a bottleneck when uploading to S3

## Future Security Improvements

- Encryption of backups at rest using AES-256
- Backup integrity verification (SHA-256 checksums)
- Audit logging for backup operations
- Automatic backup deletion with secure wiping
- Integration with key management services

## Security Updates

We recommend:

1. **Keep Go updated** - Use the latest stable version of Go
2. **Monitor dependencies** - Watch for security updates in AWS SDK and other dependencies
3. **Follow releases** - Star the repository to be notified of security updates
4. **Update regularly** - Apply updates as soon as they're released

## Supported Versions

Security updates are provided for:

- Current version (latest release)
- Previous major version (if applicable)

Older versions are not supported and may not receive security patches.

## Additional Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Guidelines](https://golang.org/doc/security)
- [AWS Security Best Practices](https://aws.amazon.com/security/best-practices/)
- [Database Security](https://mariadb.com/kb/en/security-overview/)

## Vulnerability Disclosure

After a security patch is released, we will:

1. Publish a security advisory on GitHub
2. Document the vulnerability and fix in release notes
3. Recommend all users update to the patched version
4. Credit the security researcher (if they wish to be credited)

## Contact

- **Security Email**: hey@codemeapixel.dev
- **GitHub**: https://github.com/CodeMeAPixel/pxBackupManager

Thank you for helping keep pxBackupManager secure!
