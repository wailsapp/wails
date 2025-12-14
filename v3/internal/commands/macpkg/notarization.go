package macpkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// NotarizationService handles Apple notarization workflow
type NotarizationService struct {
	AppleID     string
	AppPassword string
	TeamID      string
}

// NotarizationStatus represents the status of a notarization request
type NotarizationStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Name   string `json:"name"`
}

// NotarizationResult contains the full result of notarization
type NotarizationResult struct {
	Status        string `json:"status"`
	StatusSummary string `json:"statusSummary"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	CreatedDate   string `json:"createdDate"`
}

// NewNotarizationService creates a new notarization service instance
func NewNotarizationService(appleID, appPassword, teamID string) *NotarizationService {
	return &NotarizationService{
		AppleID:     appleID,
		AppPassword: appPassword,
		TeamID:      teamID,
	}
}

// NotarizePackage submits a package for notarization and waits for completion
func (n *NotarizationService) NotarizePackage(pkgPath string) error {
	if err := n.validateCredentials(); err != nil {
		return fmt.Errorf("credential validation failed: %w", err)
	}
	
	// Step 1: Submit for notarization
	submissionID, err := n.submitForNotarization(pkgPath)
	if err != nil {
		return fmt.Errorf("submission failed: %w", err)
	}
	
	fmt.Printf("Package submitted for notarization with ID: %s\n", submissionID)
	
	// Step 2: Wait for notarization to complete
	if err := n.waitForNotarization(submissionID); err != nil {
		return fmt.Errorf("notarization failed: %w", err)
	}
	
	// Step 3: Staple the notarization ticket
	if err := n.staplePackage(pkgPath); err != nil {
		return fmt.Errorf("stapling failed: %w", err)
	}
	
	fmt.Println("Package successfully notarized and stapled!")
	return nil
}

// submitForNotarization submits the package to Apple's notarization service
func (n *NotarizationService) submitForNotarization(pkgPath string) (string, error) {
	cmd := exec.Command("xcrun", "notarytool", "submit",
		pkgPath,
		"--output-format", "json",
	)

	// Pass credentials via environment variables to avoid exposure in process list
	cmd.Env = append(os.Environ(),
		"NOTARYTOOL_APPLE_ID="+n.AppleID,
		"NOTARYTOOL_PASSWORD="+n.AppPassword,
		"NOTARYTOOL_TEAM_ID="+n.TeamID,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("notarytool submit failed: %w\nStderr: %s", err, stderr.String())
	}
	
	// Parse the JSON response to get submission ID
	var result map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return "", fmt.Errorf("failed to parse notarytool response: %w", err)
	}
	
	submissionID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("submission ID not found in response")
	}
	
	return submissionID, nil
}

// waitForNotarization polls the notarization status until completion
func (n *NotarizationService) waitForNotarization(submissionID string) error {
	maxRetries := 60 // Maximum 30 minutes (30 seconds * 60)
	retryInterval := 30 * time.Second
	
	for i := 0; i < maxRetries; i++ {
		status, err := n.getNotarizationStatus(submissionID)
		if err != nil {
			return fmt.Errorf("failed to get notarization status: %w", err)
		}
		
		fmt.Printf("Notarization status: %s\n", status.Status)
		
		switch status.Status {
		case "Accepted":
			return nil
		case "Invalid":
			// Get detailed log for debugging
			if err := n.getNotarizationLog(submissionID); err != nil {
				fmt.Printf("Warning: Could not retrieve notarization log: %v\n", err)
			}
			return fmt.Errorf("notarization was rejected")
		case "In Progress":
			time.Sleep(retryInterval)
			continue
		default:
			return fmt.Errorf("unexpected notarization status: %s", status.Status)
		}
	}
	
	return fmt.Errorf("notarization timed out after %d minutes", maxRetries/2)
}

// getNotarizationStatus checks the current status of a notarization request
func (n *NotarizationService) getNotarizationStatus(submissionID string) (*NotarizationStatus, error) {
	cmd := exec.Command("xcrun", "notarytool", "info",
		submissionID,
		"--output-format", "json",
	)

	// Pass credentials via environment variables to avoid exposure in process list
	cmd.Env = append(os.Environ(),
		"NOTARYTOOL_APPLE_ID="+n.AppleID,
		"NOTARYTOOL_PASSWORD="+n.AppPassword,
		"NOTARYTOOL_TEAM_ID="+n.TeamID,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("notarytool info failed: %w\nStderr: %s", err, stderr.String())
	}
	
	var status NotarizationStatus
	if err := json.Unmarshal(stdout.Bytes(), &status); err != nil {
		return nil, fmt.Errorf("failed to parse status response: %w", err)
	}
	
	return &status, nil
}

// getNotarizationLog retrieves the notarization log for debugging
func (n *NotarizationService) getNotarizationLog(submissionID string) error {
	cmd := exec.Command("xcrun", "notarytool", "log", submissionID)

	// Pass credentials via environment variables to avoid exposure in process list
	cmd.Env = append(os.Environ(),
		"NOTARYTOOL_APPLE_ID="+n.AppleID,
		"NOTARYTOOL_PASSWORD="+n.AppPassword,
		"NOTARYTOOL_TEAM_ID="+n.TeamID,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("notarytool log failed: %w\nStderr: %s", err, stderr.String())
	}
	
	fmt.Println("Notarization log:")
	fmt.Println(stdout.String())
	
	return nil
}

// staplePackage attaches the notarization ticket to the package
func (n *NotarizationService) staplePackage(pkgPath string) error {
	args := []string{
		"xcrun", "stapler", "staple",
		pkgPath,
	}
	
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("stapler failed: %w", err)
	}
	
	return nil
}

// validateCredentials ensures all required credentials are provided
func (n *NotarizationService) validateCredentials() error {
	if n.AppleID == "" {
		return fmt.Errorf("apple_id is required for notarization")
	}
	if n.AppPassword == "" {
		return fmt.Errorf("app_password is required for notarization")
	}
	if n.TeamID == "" {
		return fmt.Errorf("team_id is required for notarization")
	}
	return nil
}

// CheckNotarizationDependencies verifies that required tools are available
func CheckNotarizationDependencies() error {
	// Check for xcrun
	if _, err := exec.LookPath("xcrun"); err != nil {
		return fmt.Errorf("xcrun not found - Xcode command line tools required")
	}
	
	// Check for notarytool specifically
	cmd := exec.Command("xcrun", "notarytool", "--help")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("notarytool not available - requires Xcode 13.0 or later")
	}
	
	// Check for stapler
	cmd = exec.Command("xcrun", "stapler", "--help")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("stapler not available - required for stapling notarization tickets")
	}
	
	return nil
}

// ValidatePackageSignature verifies that a package is properly signed
func ValidatePackageSignature(pkgPath string) error {
	args := []string{
		"spctl", "--assess", "--verbose", "--type", "install", pkgPath,
	}
	
	cmd := exec.Command(args[0], args[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		output := strings.TrimSpace(stderr.String())
		if output != "" {
			return fmt.Errorf("package signature validation failed: %s", output)
		}
		return fmt.Errorf("package signature validation failed: %w", err)
	}
	
	return nil
}