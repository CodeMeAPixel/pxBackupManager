package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"pxBackupManager/backup"
	"pxBackupManager/types"
)

func main() {
	// Determine platform-specific defaults
	defaultFiveMLoc := "/opt/fivem" // Linux/Mac default
	if runtime.GOOS == "windows" {
		defaultFiveMLoc = "C:\\FiveM\\server" // Windows default
	}

	// Define command-line flags
	fiveMLoc := flag.String("fivem", defaultFiveMLoc, "Path to FiveM server directory")
	mariaDBHost := flag.String("db-host", "localhost", "MariaDB host")
	mariaDBPort := flag.Int("db-port", 3306, "MariaDB port")
	mariaDBUser := flag.String("db-user", "root", "MariaDB user")
	mariaDBPass := flag.String("db-pass", "", "MariaDB password")
	mariaDBName := flag.String("db-name", "", "MariaDB database name")
	backupDir := flag.String("backup-dir", "./backups", "Backup destination directory")
	retentionDays := flag.Int("retention", 30, "Backup retention period in days (0 = no cleanup)")
	compress := flag.Bool("compress", true, "Compress backups")
	onlyFiveM := flag.Bool("only-fivem", false, "Only backup FiveM server")
	onlyMariaDB := flag.Bool("only-mariadb", false, "Only backup MariaDB database")
	cleanup := flag.Bool("cleanup", true, "Cleanup old backups")

	// S3 backup flags
	s3Enabled := flag.Bool("s3-enabled", false, "Enable S3 backup uploads")
	s3Endpoint := flag.String("s3-endpoint", "", "S3 endpoint URL (for S3-compatible services like Backblaze)")
	s3Bucket := flag.String("s3-bucket", "", "S3 bucket name")
	s3Region := flag.String("s3-region", "us-east-1", "S3 region")
	s3AccessKey := flag.String("s3-access-key", "", "S3 access key")
	s3SecretKey := flag.String("s3-secret-key", "", "S3 secret key")
	s3UploadFiveM := flag.Bool("s3-upload-fivem", false, "Upload FiveM backup to S3")
	s3UploadMariaDB := flag.Bool("s3-upload-mariadb", false, "Upload MariaDB backup to S3")

	// Discord webhook flags
	discordWebhook := flag.String("discord-webhook", "", "Discord webhook URL for notifications")
	notifyOnSuccess := flag.Bool("notify-success", false, "Send Discord notification on success")
	notifyOnFailure := flag.Bool("notify-failure", false, "Send Discord notification on failure")

	version := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *version {
		fmt.Println("pxBackupManager v2.0.0")
		os.Exit(0)
	}

	// Resolve backup directory to absolute path
	absBackupDir, err := filepath.Abs(*backupDir)
	if err != nil {
		log.Fatalf("Failed to resolve backup directory: %v\n", err)
	}

	// Ensure backup directory exists
	if err := backup.EnsureBackupDir(absBackupDir); err != nil {
		log.Fatalf("Failed to setup backup directory: %v\n", err)
	}

	// Determine which backups to run
	enableFiveM := !*onlyMariaDB
	enableMariaDB := !*onlyFiveM

	// Validate configuration
	if enableMariaDB && *mariaDBName == "" {
		log.Fatal("MariaDB database name must be specified with -db-name flag")
	}

	if *s3Enabled {
		if *s3Bucket == "" || *s3AccessKey == "" || *s3SecretKey == "" {
			log.Fatal("S3 requires -s3-bucket, -s3-access-key, and -s3-secret-key flags")
		}
	}

	// Create backup configuration
	config := types.BackupConfig{
		FiveMLoc:          *fiveMLoc,
		MariaDBHost:       *mariaDBHost,
		MariaDBPort:       *mariaDBPort,
		MariaDBUser:       *mariaDBUser,
		MariaDBPassword:   *mariaDBPass,
		MariaDBDatabase:   *mariaDBName,
		BackupDir:         absBackupDir,
		RetentionDays:     *retentionDays,
		Compress:          *compress,
		EnableFiveM:       enableFiveM,
		EnableMariaDB:     enableMariaDB,
		S3Enabled:         *s3Enabled,
		S3Endpoint:        *s3Endpoint,
		S3Bucket:          *s3Bucket,
		S3Region:          *s3Region,
		S3AccessKey:       *s3AccessKey,
		S3SecretKey:       *s3SecretKey,
		S3UploadFiveM:     *s3UploadFiveM,
		S3UploadMariaDB:   *s3UploadMariaDB,
		DiscordWebhookURL: *discordWebhook,
		NotifyOnSuccess:   *notifyOnSuccess,
		NotifyOnFailure:   *notifyOnFailure,
	}

	// Print configuration
	fmt.Println("===== PX BACKUP MANAGER =====")
	fmt.Printf("Backup Directory: %s\n", config.BackupDir)
	fmt.Printf("Compression: %v\n", config.Compress)
	fmt.Printf("Retention: %d days\n", config.RetentionDays)
	fmt.Printf("FiveM Backup: %v\n", config.EnableFiveM)
	fmt.Printf("MariaDB Backup: %v\n", config.EnableMariaDB)
	if config.S3Enabled {
		fmt.Printf("S3 Upload: enabled (bucket: %s)\n", config.S3Bucket)
	}
	if config.DiscordWebhookURL != "" {
		fmt.Printf("Discord Notifications: enabled\n")
	}
	fmt.Println("=============================")

	startTime := time.Now()
	results := make([]types.BackupResult, 0)

	// Run FiveM backup
	if config.EnableFiveM {
		start := time.Now()
		filename, err := backup.BackupFiveM(config.FiveMLoc, config.BackupDir, config.Compress)
		if err != nil {
			fmt.Printf("FiveM backup failed: %v\n", err)
			results = append(results, types.BackupResult{
				Service:   "fivem",
				Success:   false,
				Message:   err.Error(),
				Duration:  time.Since(start).Milliseconds(),
				Timestamp: start.Format(time.RFC3339),
			})
		} else {
			// Get file info
			backupPath := filepath.Join(config.BackupDir, filename)
			fileInfo, _ := os.Stat(backupPath)

			result := types.BackupResult{
				Service:   "fivem",
				Success:   true,
				Message:   "Successfully backed up FiveM server",
				Filename:  filename,
				Size:      fileInfo.Size(),
				Duration:  time.Since(start).Milliseconds(),
				Timestamp: start.Format(time.RFC3339),
			}

			// Handle S3 upload if enabled
			if config.S3Enabled && config.S3UploadFiveM {
				fmt.Printf("Uploading FiveM backup to S3...\n")
				s3URL, err := backup.UploadToS3(backupPath, config.S3Bucket, config.S3Region, config.S3Endpoint, config.S3AccessKey, config.S3SecretKey)
				if err != nil {
					fmt.Printf("S3 upload failed: %v\n", err)
				} else {
					result.S3URL = s3URL
				}
			}

			results = append(results, result)
		}
		fmt.Println()
	}

	// Run MariaDB backup
	if config.EnableMariaDB {
		start := time.Now()
		filename, err := backup.BackupMariaDB(
			config.MariaDBHost,
			config.MariaDBPort,
			config.MariaDBUser,
			config.MariaDBPassword,
			config.MariaDBDatabase,
			config.BackupDir,
			config.Compress,
		)
		if err != nil {
			fmt.Printf("MariaDB backup failed: %v\n", err)
			results = append(results, types.BackupResult{
				Service:   "mariadb",
				Success:   false,
				Message:   err.Error(),
				Duration:  time.Since(start).Milliseconds(),
				Timestamp: start.Format(time.RFC3339),
			})
		} else {
			// Get file info
			backupPath := filepath.Join(config.BackupDir, filename)
			fileInfo, _ := os.Stat(backupPath)

			result := types.BackupResult{
				Service:   "mariadb",
				Success:   true,
				Message:   "Successfully backed up MariaDB database",
				Filename:  filename,
				Size:      fileInfo.Size(),
				Duration:  time.Since(start).Milliseconds(),
				Timestamp: start.Format(time.RFC3339),
			}

			// Handle S3 upload if enabled
			if config.S3Enabled && config.S3UploadMariaDB {
				fmt.Printf("Uploading MariaDB backup to S3...\n")
				s3URL, err := backup.UploadToS3(backupPath, config.S3Bucket, config.S3Region, config.S3Endpoint, config.S3AccessKey, config.S3SecretKey)
				if err != nil {
					fmt.Printf("S3 upload failed: %v\n", err)
				} else {
					result.S3URL = s3URL
				}
			}

			results = append(results, result)
		}
		fmt.Println()
	}

	// Cleanup old backups
	if *cleanup && config.RetentionDays > 0 {
		fmt.Println("Cleaning up old backups...")
		if err := backup.CleanupOldBackups(config.BackupDir, config.RetentionDays); err != nil {
			fmt.Printf("Warning: cleanup failed: %v\n", err)
		}
		fmt.Println()
	}

	// Print summary
	fmt.Println("===== BACKUP SUMMARY =====")
	totalSuccess := 0
	totalFailed := 0
	for _, result := range results {
		status := "✓ SUCCESS"
		if !result.Success {
			status = "✗ FAILED"
			totalFailed++
		} else {
			totalSuccess++
		}
		fmt.Printf("[%s] %s - %dms\n", result.Service, status, result.Duration)
		if result.Message != "" && !result.Success {
			fmt.Printf("    Error: %s\n", result.Message)
		}
		if result.Size > 0 {
			fmt.Printf("    Size: %.2f MB\n", float64(result.Size)/1024/1024)
		}
		if result.S3URL != "" {
			fmt.Printf("    S3: %s\n", result.S3URL)
		}
	}
	fmt.Printf("\nTotal: %d succeeded, %d failed\n", totalSuccess, totalFailed)
	fmt.Printf("Elapsed time: %v\n", time.Since(startTime))
	fmt.Println("==========================")

	// Send Discord notification if configured
	if config.DiscordWebhookURL != "" {
		shouldNotify := (totalFailed == 0 && config.NotifyOnSuccess) || (totalFailed > 0 && config.NotifyOnFailure)
		if shouldNotify {
			summary := fmt.Sprintf("Backup completed: %d succeeded, %d failed", totalSuccess, totalFailed)
			if err := backup.SendDiscordNotification(config.DiscordWebhookURL, results, summary); err != nil {
				fmt.Printf("Warning: failed to send Discord notification: %v\n", err)
			}
		}
	}

	// Exit with error if any backup failed
	if totalFailed > 0 {
		os.Exit(1)
	}
}
