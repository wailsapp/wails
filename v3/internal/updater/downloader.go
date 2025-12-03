package updater

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Downloader handles downloading updates
type Downloader struct {
	httpClient *http.Client
	tempDir    string
}

// NewDownloader creates a new update downloader
func NewDownloader() *Downloader {
	return &Downloader{
		httpClient: &http.Client{
			Timeout: 0, // No timeout for downloads
		},
		tempDir: os.TempDir(),
	}
}

// SetTempDir sets the temporary directory for downloads
func (d *Downloader) SetTempDir(dir string) {
	d.tempDir = dir
}

// Download downloads an update file and verifies its checksum
func (d *Downloader) Download(ctx context.Context, info *UpdateInfo, progress ProgressCallback) (string, error) {
	// Determine what to download (prefer patch if available)
	downloadURL := info.DownloadURL
	expectedChecksum := info.Checksum
	expectedSize := info.Size
	isPatch := false

	if info.PatchURL != "" {
		downloadURL = info.PatchURL
		expectedChecksum = info.PatchChecksum
		expectedSize = info.PatchSize
		isPatch = true
	}

	// Create temp file
	ext := ".tar.gz"
	if isPatch {
		ext = ".patch"
	}
	tempFile, err := os.CreateTemp(d.tempDir, fmt.Sprintf("wails-update-*%s", ext))
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()

	// Ensure cleanup on error
	success := false
	defer func() {
		tempFile.Close()
		if !success {
			os.Remove(tempPath)
		}
	}()

	// Start download
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Wails-Updater/1.0")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	// Get content length
	contentLength := resp.ContentLength
	if contentLength < 0 {
		contentLength = expectedSize
	}

	// Create progress writer
	hasher := sha256.New()
	writer := io.MultiWriter(tempFile, hasher)

	var downloaded int64
	lastProgress := time.Now()
	lastBytes := int64(0)

	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := writer.Write(buf[:n]); writeErr != nil {
				return "", fmt.Errorf("write failed: %w", writeErr)
			}
			downloaded += int64(n)

			// Report progress
			if progress != nil {
				now := time.Now()
				elapsed := now.Sub(lastProgress).Seconds()
				if elapsed >= 0.1 { // Update every 100ms
					bytesPerSecond := float64(downloaded-lastBytes) / elapsed
					progress(DownloadProgress{
						Downloaded:     downloaded,
						Total:          contentLength,
						Percentage:     float64(downloaded) / float64(contentLength) * 100,
						BytesPerSecond: bytesPerSecond,
					})
					lastProgress = now
					lastBytes = downloaded
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("read failed: %w", err)
		}
	}

	// Final progress update
	if progress != nil {
		progress(DownloadProgress{
			Downloaded:     downloaded,
			Total:          contentLength,
			Percentage:     100,
			BytesPerSecond: 0,
		})
	}

	// Verify checksum
	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	if actualChecksum != expectedChecksum {
		return "", fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	success = true
	return tempPath, nil
}

// DownloadFull downloads the full update (not patch)
func (d *Downloader) DownloadFull(ctx context.Context, info *UpdateInfo, progress ProgressCallback) (string, error) {
	// Create a copy of info without patch URL to force full download
	fullInfo := *info
	fullInfo.PatchURL = ""
	fullInfo.PatchChecksum = ""
	fullInfo.PatchSize = 0

	return d.Download(ctx, &fullInfo, progress)
}

// VerifyChecksum verifies the SHA256 checksum of a file
func VerifyChecksum(filePath, expectedChecksum string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return err
	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

// GetDownloadPath returns the path where an update would be downloaded
func (d *Downloader) GetDownloadPath(version string, isPatch bool) string {
	ext := ".tar.gz"
	if isPatch {
		ext = ".patch"
	}
	return filepath.Join(d.tempDir, fmt.Sprintf("wails-update-%s%s", version, ext))
}

// CleanupOldDownloads removes old update downloads
func (d *Downloader) CleanupOldDownloads() error {
	pattern := filepath.Join(d.tempDir, "wails-update-*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}

		// Remove files older than 24 hours
		if time.Since(info.ModTime()) > 24*time.Hour {
			os.Remove(match)
		}
	}

	return nil
}
