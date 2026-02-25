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
	var skipCount int

	err = filepath.Walk(fiveMLoc, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Warning: error accessing %s: %v\n", filePath, err)
			return nil // Skip files we can't access
		}

		// Get relative path for tar archive
		relPath, err := filepath.Rel(baseDir, filePath)
		if err != nil {
			fmt.Printf("Warning: failed to get relative path for %s: %v\n", filePath, err)
			skipCount++
			return nil
		}

		// Use forward slashes for tar archive
		tarPath := strings.ReplaceAll(relPath, "\\", "/")

		// Skip files with extremely long paths (> 8000 chars indicates potential issues)
		if len(tarPath) > 8000 {
			fmt.Printf("Warning: skipping file with extremely long path: %s\n", tarPath)
			skipCount++
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			fmt.Printf("Warning: failed to create tar header for %s: %v\n", filePath, err)
			skipCount++
			return nil
		}

		header.Name = tarPath

		if err := tw.WriteHeader(header); err != nil {
			fmt.Printf("Warning: failed to write tar header for %s: %v\n", filePath, err)
			skipCount++
			return nil // Continue with next file instead of failing entirely
		}

		// Write file content if it's a regular file
		if !info.IsDir() && info.Mode().IsRegular() {
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Warning: failed to open file %s: %v\n", filePath, err)
				skipCount++
				return nil
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				fmt.Printf("Warning: failed to write file content for %s: %v\n", filePath, err)
				skipCount++
				return nil
			}
		}

		return nil
	})

	if err != nil {
		os.Remove(backupPath) // Clean up on failure
		return "", fmt.Errorf("tar creation failed: %w", err)
	}

	if skipCount > 0 {
		fmt.Printf("Warning: %d files were skipped during backup\n", skipCount)
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
