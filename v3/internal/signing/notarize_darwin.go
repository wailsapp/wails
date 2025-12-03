//go:build darwin

package signing

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// NotarizeOptions defines options for notarization
type NotarizeOptions struct {
	// AppPath is the path to the signed .app bundle
	AppPath string
	// AppleID is the Apple ID for notarization
	AppleID string
	// TeamID is the Apple Developer Team ID
	TeamID string
	// KeychainProfile is the name of the stored credentials profile
	KeychainProfile string
	// AppSpecificPassword is an app-specific password (alternative to keychain profile)
	AppSpecificPassword string
	// Wait determines whether to wait for notarization to complete
	Wait bool
	// Verbose enables verbose output
	Verbose bool
}

// NotarizationStatus represents the status of a notarization submission
type NotarizationStatus struct {
	ID        string
	Status    string
	Message   string
	CreatedAt time.Time
}

// notarytoolOutput represents the JSON output from notarytool
type notarytoolOutput struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Notarize submits an application for notarization
func Notarize(options NotarizeOptions) (*NotarizationStatus, error) {
	if options.AppPath == "" {
		return nil, fmt.Errorf("app path is required")
	}

	if !strings.HasSuffix(options.AppPath, ".app") {
		return nil, fmt.Errorf("path must be an .app bundle")
	}

	// Create a temporary zip file for submission
	zipPath, err := createZipForNotarization(options.AppPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create zip: %w", err)
	}
	defer os.Remove(zipPath)

	// Build notarytool command
	args := []string{"notarytool", "submit", zipPath, "--output-format", "json"}

	if options.KeychainProfile != "" {
		args = append(args, "--keychain-profile", options.KeychainProfile)
	} else {
		if options.AppleID == "" {
			return nil, fmt.Errorf("apple_id is required when not using keychain_profile")
		}
		if options.TeamID == "" {
			return nil, fmt.Errorf("team_id is required when not using keychain_profile")
		}
		if options.AppSpecificPassword == "" {
			return nil, fmt.Errorf("app_specific_password is required when not using keychain_profile")
		}
		args = append(args,
			"--apple-id", expandEnvVar(options.AppleID),
			"--team-id", expandEnvVar(options.TeamID),
			"--password", expandEnvVar(options.AppSpecificPassword),
		)
	}

	if options.Wait {
		args = append(args, "--wait")
	}

	cmd := exec.Command("xcrun", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if options.Verbose {
		fmt.Printf("Running: xcrun %s\n", strings.Join(args, " "))
	}

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("notarytool submit failed: %s\n%s", stderr.String(), err)
	}

	// Parse output
	var output notarytoolOutput
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return nil, fmt.Errorf("failed to parse notarytool output: %w", err)
	}

	return &NotarizationStatus{
		ID:        output.ID,
		Status:    output.Status,
		Message:   output.Message,
		CreatedAt: time.Now(),
	}, nil
}

// WaitForNotarization waits for a notarization submission to complete
func WaitForNotarization(submissionID string, options NotarizeOptions) (*NotarizationStatus, error) {
	args := []string{"notarytool", "wait", submissionID, "--output-format", "json"}

	if options.KeychainProfile != "" {
		args = append(args, "--keychain-profile", options.KeychainProfile)
	} else {
		args = append(args,
			"--apple-id", expandEnvVar(options.AppleID),
			"--team-id", expandEnvVar(options.TeamID),
			"--password", expandEnvVar(options.AppSpecificPassword),
		)
	}

	cmd := exec.Command("xcrun", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("notarytool wait failed: %s\n%s", stderr.String(), err)
	}

	var output notarytoolOutput
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return nil, fmt.Errorf("failed to parse notarytool output: %w", err)
	}

	return &NotarizationStatus{
		ID:      output.ID,
		Status:  output.Status,
		Message: output.Message,
	}, nil
}

// GetNotarizationLog retrieves the log for a notarization submission
func GetNotarizationLog(submissionID string, options NotarizeOptions) (string, error) {
	args := []string{"notarytool", "log", submissionID}

	if options.KeychainProfile != "" {
		args = append(args, "--keychain-profile", options.KeychainProfile)
	} else {
		args = append(args,
			"--apple-id", expandEnvVar(options.AppleID),
			"--team-id", expandEnvVar(options.TeamID),
			"--password", expandEnvVar(options.AppSpecificPassword),
		)
	}

	cmd := exec.Command("xcrun", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("notarytool log failed: %s\n%s", stderr.String(), err)
	}

	return stdout.String(), nil
}

// Staple staples the notarization ticket to the application
func Staple(appPath string) error {
	cmd := exec.Command("xcrun", "stapler", "staple", appPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("stapler failed: %s", stderr.String())
	}

	return nil
}

// NotarizeAndStaple performs the complete notarization workflow
func NotarizeAndStaple(options NotarizeOptions) error {
	options.Wait = true

	fmt.Printf("Submitting %s for notarization...\n", filepath.Base(options.AppPath))

	status, err := Notarize(options)
	if err != nil {
		return err
	}

	fmt.Printf("Submission ID: %s\n", status.ID)
	fmt.Printf("Status: %s\n", status.Status)

	if status.Status != "Accepted" {
		log, _ := GetNotarizationLog(status.ID, options)
		return fmt.Errorf("notarization failed with status '%s': %s\nLog:\n%s", status.Status, status.Message, log)
	}

	fmt.Println("Stapling notarization ticket...")
	if err := Staple(options.AppPath); err != nil {
		return fmt.Errorf("failed to staple: %w", err)
	}

	fmt.Println("Notarization complete!")
	return nil
}

// createZipForNotarization creates a zip file of the app bundle for notarization
func createZipForNotarization(appPath string) (string, error) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "notarize-*.zip")
	if err != nil {
		return "", err
	}
	tmpPath := tmpFile.Name()

	zipWriter := zip.NewWriter(tmpFile)

	baseDir := filepath.Dir(appPath)
	appName := filepath.Base(appPath)

	err = filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from app directory
		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		// Create zip header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		// Handle symlinks
		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			header.SetMode(info.Mode())
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}
			_, err = writer.Write([]byte(link))
			return err
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		}

		return nil
	})

	if err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", err
	}

	if err := zipWriter.Close(); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", err
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return "", err
	}

	_ = appName // silence unused warning
	return tmpPath, nil
}

// expandEnvVar expands environment variable references like ${VAR_NAME}
func expandEnvVar(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		varName := value[2 : len(value)-1]
		return os.Getenv(varName)
	}
	return value
}

// StoreCredentials stores notarization credentials in the keychain
// This is a helper to run: xcrun notarytool store-credentials
func StoreCredentials(profileName, appleID, teamID, password string) error {
	cmd := exec.Command("xcrun", "notarytool", "store-credentials", profileName,
		"--apple-id", appleID,
		"--team-id", teamID,
		"--password", password,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
