package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/internal/defaults"
)

type SigningStatus struct {
	Darwin  DarwinSigningStatus  `json:"darwin"`
	Windows WindowsSigningStatus `json:"windows"`
	Linux   LinuxSigningStatus   `json:"linux"`
}

type DarwinSigningStatus struct {
	HasIdentity     bool     `json:"hasIdentity"`
	Identity        string   `json:"identity,omitempty"`
	Identities      []string `json:"identities,omitempty"`
	HasNotarization bool     `json:"hasNotarization"`
	TeamID          string   `json:"teamID,omitempty"`
	ConfigSource    string   `json:"configSource,omitempty"`
}

type WindowsSigningStatus struct {
	HasCertificate  bool   `json:"hasCertificate"`
	CertificateType string `json:"certificateType,omitempty"`
	HasSignTool     bool   `json:"hasSignTool"`
	TimestampServer string `json:"timestampServer,omitempty"`
	ConfigSource    string `json:"configSource,omitempty"`
}

type LinuxSigningStatus struct {
	HasGPGKey    bool   `json:"hasGpgKey"`
	GPGKeyID     string `json:"gpgKeyID,omitempty"`
	ConfigSource string `json:"configSource,omitempty"`
}

func CheckSigning() SigningStatus {
	globalDefaults, err := defaults.Load()
	if err != nil {
		// Log warning but continue with empty defaults
		// This allows the check to proceed even if config is invalid
		fmt.Fprintf(os.Stderr, "Warning: could not load global defaults: %v\n", err)
	}

	return SigningStatus{
		Darwin:  checkDarwinSigning(globalDefaults),
		Windows: checkWindowsSigning(globalDefaults),
		Linux:   checkLinuxSigning(globalDefaults),
	}
}

func checkDarwinSigning(cfg defaults.GlobalDefaults) DarwinSigningStatus {
	status := DarwinSigningStatus{}

	if cfg.Signing.Darwin.Identity != "" {
		status.HasIdentity = true
		status.Identity = cfg.Signing.Darwin.Identity
		status.TeamID = cfg.Signing.Darwin.TeamID
		status.ConfigSource = "defaults.yaml"
	}

	if cfg.Signing.Darwin.KeychainProfile != "" || cfg.Signing.Darwin.APIKeyID != "" {
		status.HasNotarization = true
	}

	if runtime.GOOS == "darwin" {
		identities := getMacOSSigningIdentities()
		status.Identities = identities
		if len(identities) > 0 && !status.HasIdentity {
			status.HasIdentity = true
			status.Identity = identities[0]
			status.ConfigSource = "keychain"
		}
	}

	return status
}

func checkWindowsSigning(cfg defaults.GlobalDefaults) WindowsSigningStatus {
	status := WindowsSigningStatus{}

	if cfg.Signing.Windows.CertificatePath != "" {
		status.HasCertificate = true
		status.CertificateType = "file"
		status.ConfigSource = "defaults.yaml"
	} else if cfg.Signing.Windows.Thumbprint != "" {
		status.HasCertificate = true
		status.CertificateType = "store"
		status.ConfigSource = "defaults.yaml"
	} else if cfg.Signing.Windows.CloudProvider != "" {
		status.HasCertificate = true
		status.CertificateType = "cloud:" + cfg.Signing.Windows.CloudProvider
		status.ConfigSource = "defaults.yaml"
	}

	status.TimestampServer = cfg.Signing.Windows.TimestampServer
	if status.TimestampServer == "" {
		status.TimestampServer = "http://timestamp.digicert.com"
	}

	if runtime.GOOS == "windows" {
		_, err := exec.LookPath("signtool.exe")
		status.HasSignTool = err == nil
	}

	return status
}

func checkLinuxSigning(cfg defaults.GlobalDefaults) LinuxSigningStatus {
	status := LinuxSigningStatus{}

	if cfg.Signing.Linux.GPGKeyPath != "" {
		status.HasGPGKey = true
		status.GPGKeyID = cfg.Signing.Linux.GPGKeyID
		status.ConfigSource = "defaults.yaml"
	}

	if !status.HasGPGKey {
		keyID := getDefaultGPGKey()
		if keyID != "" {
			status.HasGPGKey = true
			status.GPGKeyID = keyID
			status.ConfigSource = "gpg"
		}
	}

	return status
}

func getMacOSSigningIdentities() []string {
	if runtime.GOOS != "darwin" {
		return nil
	}

	cmd := exec.Command("security", "find-identity", "-v", "-p", "codesigning")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var identities []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "\"") && strings.Contains(line, "Developer ID") {
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start != -1 && end > start {
				identity := line[start+1 : end]
				identities = append(identities, identity)
			}
		}
	}

	return identities
}

func getDefaultGPGKey() string {
	cmd := exec.Command("gpg", "--list-secret-keys", "--keyid-format", "long")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "sec") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.Contains(part, "/") {
					keyParts := strings.Split(part, "/")
					if len(keyParts) > 1 {
						return keyParts[1]
					}
				}
			}
		}
	}

	return ""
}

func formatSigningStatus(status SigningStatus) map[string]string {
	result := make(map[string]string)

	if status.Darwin.HasIdentity {
		identity := status.Darwin.Identity
		if len(identity) > 50 {
			identity = identity[:47] + "..."
		}
		value := identity
		if status.Darwin.HasNotarization {
			value += " (notarization configured)"
		}
		result["macOS Signing"] = value
	} else {
		result["macOS Signing"] = "Not configured"
	}

	if status.Windows.HasCertificate {
		value := "Configured (" + status.Windows.CertificateType + ")"
		if status.Windows.HasSignTool {
			value += " - signtool available"
		}
		result["Windows Signing"] = value
	} else {
		result["Windows Signing"] = "Not configured"
	}

	if status.Linux.HasGPGKey {
		result["Linux Signing"] = "GPG key: " + status.Linux.GPGKeyID
	} else {
		result["Linux Signing"] = "Not configured (GPG)"
	}

	return result
}
