package types

// BackupConfig holds the configuration for backups
type BackupConfig struct {
	// FiveM server path
	FiveMLoc string

	// MariaDB configuration
	MariaDBHost     string
	MariaDBPort     int
	MariaDBUser     string
	MariaDBPassword string
	MariaDBDatabase string

	// Backup destination
	BackupDir string

	// Retention policy (days)
	RetentionDays int

	// Compression
	Compress bool

	// Enable FiveM backup
	EnableFiveM bool

	// Enable MariaDB backup
	EnableMariaDB bool

	// S3 configuration
	S3Enabled       bool
	S3Endpoint      string // For S3-compatible services like Backblaze
	S3Bucket        string
	S3Region        string
	S3AccessKey     string
	S3SecretKey     string
	S3UploadFiveM   bool
	S3UploadMariaDB bool

	// Discord webhook configuration
	DiscordWebhookURL string
	NotifyOnSuccess   bool
	NotifyOnFailure   bool
}

// BackupResult contains information about backup execution
type BackupResult struct {
	Service   string // "fivem" or "mariadb"
	Success   bool
	Message   string
	Filename  string
	Size      int64
	Duration  int64 // milliseconds
	Timestamp string
	S3URL     string // S3 URL if uploaded to S3
}
