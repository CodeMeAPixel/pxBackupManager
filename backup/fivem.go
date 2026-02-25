package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// BackupFiveM creates a backup of the FiveM server and returns the backup filename
func BackupFiveM(fiveMLoc, backupDir string, compress bool) (string, error) {
	// Verify FiveM location exists
	if _, err := os.Stat(fiveMLoc); os.IsNotExist(err) {
		return "", fmt.Errorf("FiveM location does not exist: %s", fiveMLoc)
	}

	// Create backup filename
	filename := GetBackupFilename("fivem", compress)
	backupPath := filepath.Join(backupDir, filename)

	fmt.Printf("Starting FiveM backup from %s\n", fiveMLoc)
	fmt.Printf("Backup destination: %s\n", backupPath)

	// Create backup file
	backupFile, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer backupFile.Close()

	// Create tar writer (optionally wrapped with gzip)
	var tw *tar.Writer
	if compress {
		gz := gzip.NewWriter(backupFile)
		defer gz.Close()
		tw = tar.NewWriter(gz)
	} else {
		tw = tar.NewWriter(backupFile)
	}
	defer tw.Close()

	// Walk through the FiveM directory and add files to tar
	baseDir := filepath.Dir(fiveMLoc)

	err = filepath.Walk(fiveMLoc, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path for tar archive
		relPath, err := filepath.Rel(baseDir, filePath)
		if err != nil {
			return err
		}

		// Use forward slashes for tar archive
		tarPath := strings.ReplaceAll(relPath, "\\", "/")

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name = tarPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// Write file content if it's a regular file
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		os.Remove(backupPath) // Clean up on failure
		return "", fmt.Errorf("tar creation failed: %w", err)
	}

	// Verify backup was created
	if info, err := os.Stat(backupPath); err != nil {
		return "", fmt.Errorf("backup file verification failed: %w", err)
	} else {
		fmt.Printf("FiveM backup completed successfully\n")
		fmt.Printf("File: %s\n", filename)
		fmt.Printf("Size: %.2f MB\n", float64(info.Size())/1024/1024)
	}

	return filename, nil
}

// BackupFiveMASync creates a backup of FiveM server asynchronously and returns status
func BackupFiveMASync(fiveMLoc, backupDir string, compress bool) <-chan string {
	statusChan := make(chan string)

	go func() {
		defer close(statusChan)
		statusChan <- "FiveM backup started"

		if _, err := os.Stat(fiveMLoc); os.IsNotExist(err) {
			statusChan <- fmt.Sprintf("ERROR: FiveM location does not exist: %s", fiveMLoc)
			return
		}

		filename, err := BackupFiveM(fiveMLoc, backupDir, compress)
		if err != nil {
			statusChan <- fmt.Sprintf("ERROR: %v", err)
			return
		}
		backupPath := filepath.Join(backupDir, filename)

		if info, err := os.Stat(backupPath); err != nil {
			statusChan <- fmt.Sprintf("ERROR: backup verification failed: %v", err)
		} else {
			statusChan <- fmt.Sprintf("SUCCESS: FiveM backup completed - %s (%.2f MB)", filename, float64(info.Size())/1024/1024)
		}
	}()

	return statusChan
}
