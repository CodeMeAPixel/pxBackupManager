package backup

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func BackupMariaDB(host string, port int, user, password, database, backupDir string, compress bool) (string, error) {
	filename := GetBackupFilename("mariadb", compress)
	backupPath := filepath.Join(backupDir, filename)

	fmt.Printf("Starting MariaDB backup for database: %s\n", database)
	fmt.Printf("Backup destination: %s\n", backupPath)

	args := []string{
		"--host=" + host,
		"--port=" + fmt.Sprintf("%d", port),
		"--user=" + user,
		"--password=" + password,
		"--single-transaction",
		"--lock-tables=false",
		database,
	}

	cmd, err := getMySQLDumpCommand(args)
	if err != nil {
		return "", fmt.Errorf("mysqldump not available or path issue: %w", err)
	}

	backupFile, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer backupFile.Close()

	var writer io.WriteCloser
	if compress {
		writer = gzip.NewWriter(backupFile)
	} else {
		writer = backupFile
	}
	defer writer.Close()

	cmd.Stdout = writer
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		os.Remove(backupPath)
		return "", fmt.Errorf("mysqldump command failed: %w", err)
	}

	if info, err := os.Stat(backupPath); err != nil {
		return "", fmt.Errorf("backup file verification failed: %w", err)
	} else {
		fmt.Printf("MariaDB backup completed successfully\n")
		fmt.Printf("File: %s\n", filename)
		fmt.Printf("Size: %.2f MB\n", float64(info.Size())/1024/1024)
	}

	return filename, nil
}

func BackupMariaDBASync(host string, port int, user, password, database, backupDir string, compress bool) <-chan string {
	statusChan := make(chan string)

	go func() {
		defer close(statusChan)
		statusChan <- "MariaDB backup started"

		filename, err := BackupMariaDB(host, port, user, password, database, backupDir, compress)
		if err != nil {
			statusChan <- fmt.Sprintf("ERROR: %v", err)
			return
		}

		backupPath := filepath.Join(backupDir, filename)
		if info, err := os.Stat(backupPath); err != nil {
			statusChan <- fmt.Sprintf("ERROR: backup verification failed: %v", err)
		} else {
			statusChan <- fmt.Sprintf("SUCCESS: MariaDB backup completed - %s (%.2f MB)", filename, float64(info.Size())/1024/1024)
		}
	}()

	return statusChan
}
