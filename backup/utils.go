package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// GetBackupFilename generates a backup filename with timestamp
func GetBackupFilename(service string, compress bool) string {
	timestamp := time.Now().Format("2006-01-02-150405")
	ext := ".tar"
	if compress {
		ext = ".tar.gz"
	}
	return fmt.Sprintf("%s-backup-%s%s", service, timestamp, ext)
}

// EnsureBackupDir creates backup directory if it doesn't exist
func EnsureBackupDir(backupDir string) error {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}
	return nil
}

// CleanupOldBackups removes backups older than retentionDays
func CleanupOldBackups(backupDir string, retentionDays int) error {
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoffTime) {
			filePath := filepath.Join(backupDir, entry.Name())
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("Warning: failed to remove old backup %s: %v\n", filePath, err)
			} else {
				fmt.Printf("Removed old backup: %s\n", filePath)
			}
		}
	}

	return nil
}

// GetDirectorySize returns the total size of a directory in bytes
func GetDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// getMySQLDumpCommand creates a mysqldump command compatible with cross-platform systems
func getMySQLDumpCommand(args []string) (*exec.Cmd, error) {
	// Try direct mysqldump first (works on Linux/Mac and Windows if in PATH)
	cmd := exec.Command("mysqldump", args...)
	if _, err := exec.LookPath("mysqldump"); err == nil {
		return cmd, nil
	}

	// On Windows, try common installation paths
	if runtime.GOOS == "windows" {
		commonPaths := []string{
			"C:\\Program Files\\MySQL\\MySQL Server 8.0\\bin\\mysqldump.exe",
			"C:\\Program Files\\MySQL\\MySQL Server 5.7\\bin\\mysqldump.exe",
			"C:\\Program Files (x86)\\MySQL\\MySQL Server 8.0\\bin\\mysqldump.exe",
			"C:\\Program Files (x86)\\MySQL\\MySQL Server 5.7\\bin\\mysqldump.exe",
		}

		for _, path := range commonPaths {
			if _, err := os.Stat(path); err == nil {
				return exec.Command(path, args...), nil
			}
		}
	}

	// If not found, return error
	return nil, fmt.Errorf("mysqldump not found in PATH or standard installation directories")
}
