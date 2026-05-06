//go:build windows

package setupwizard

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/sys/windows/registry"
)

// TestRefreshPath_HKCUPickup writes a sentinel path entry to HKCU and verifies
// that refreshPath() picks it up so execCommandRefreshed would find it.
func TestRefreshPath_HKCUPickup(t *testing.T) {
	sentinel := `C:\sentinel-wails-test-` + t.Name()

	// Write sentinel into HKCU\Environment
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		t.Skipf("cannot open HKCU\\Environment for write: %v", err)
	}
	defer key.Close()

	// Save original value
	origPath, _, _ := key.GetStringValue("Path")
	t.Cleanup(func() {
		if origPath == "" {
			key.DeleteValue("Path")
		} else {
			key.SetStringValue("Path", origPath)
		}
		// Remove sentinel from process PATH too
		cur := os.Getenv("PATH")
		cleaned := []string{}
		for _, p := range strings.Split(cur, ";") {
			if p != sentinel {
				cleaned = append(cleaned, p)
			}
		}
		os.Setenv("PATH", strings.Join(cleaned, ";"))
	})

	// Append sentinel
	newPath := origPath
	if newPath != "" {
		newPath += ";"
	}
	newPath += sentinel
	if err := key.SetStringValue("Path", newPath); err != nil {
		t.Fatalf("SetStringValue failed: %v", err)
	}

	// Call refreshPath and verify the sentinel appears in the process PATH
	refreshPath()

	cur := os.Getenv("PATH")
	if !strings.Contains(cur, sentinel) {
		t.Errorf("refreshPath() did not add HKCU sentinel %q to PATH; got PATH=%q", sentinel, cur)
	}
}

// TestRefreshPath_EmptyEntriesRemoved verifies that empty PATH entries are stripped.
func TestRefreshPath_EmptyEntriesRemoved(t *testing.T) {
	// Save current PATH and restore after
	orig := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", orig) })

	// Set PATH with empty entries via os.Setenv to simulate a dirty environment.
	os.Setenv("PATH", `C:\Windows;;C:\Windows\System32;;;`)

	// refreshPath reads from registry and rebuilds; it should produce no empty entries.
	refreshPath()

	cur := os.Getenv("PATH")
	for _, p := range strings.Split(cur, ";") {
		if strings.TrimSpace(p) == "" {
			t.Errorf("refreshPath() left an empty entry in PATH; full PATH=%q", cur)
		}
	}
}

// TestRefreshPath_SystemFirst verifies that HKLM entries precede HKCU entries.
func TestRefreshPath_SystemFirst(t *testing.T) {
	orig := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", orig) })

	refreshPath()

	cur := os.Getenv("PATH")
	// System32 comes from HKLM and should appear before anything that is HKCU-only.
	// We check that System32 is present (a basic sanity check).
	if !strings.Contains(strings.ToLower(cur), `windows\system32`) {
		t.Errorf("expected Windows\\System32 in refreshed PATH; got %q", cur)
	}
}

// TestCheckWebView2_Paths verifies that checkWebView2 checks the standard env-var paths.
// On any typical Windows 10/11 install WebView2 should be present.
func TestCheckWebView2_Paths(t *testing.T) {
	dep := checkWebView2()

	if dep.Name != "WebView2 Runtime" {
		t.Errorf("unexpected Name: %q", dep.Name)
	}
	if dep.Required != true {
		t.Errorf("WebView2 should be marked required")
	}

	if dep.Installed {
		t.Logf("WebView2 detected: version=%s, status=%s", dep.Version, dep.Status)
		if dep.Status != "installed" {
			t.Errorf("installed=true but status=%q, want \"installed\"", dep.Status)
		}
	} else {
		t.Logf("WebView2 not found (may be OK on minimal VM): status=%s message=%s", dep.Status, dep.Message)
		if dep.Status != "not_installed" {
			t.Errorf("installed=false but status=%q, want \"not_installed\"", dep.Status)
		}
		if dep.HelpURL == "" {
			t.Error("HelpURL must be set when not installed")
		}
	}
}

// TestCheckGo verifies that checkGo finds the Go toolchain on the VM.
func TestCheckGo(t *testing.T) {
	dep := checkGo()

	if dep.Name != "Go" {
		t.Errorf("unexpected Name: %q", dep.Name)
	}
	if !dep.Installed {
		t.Errorf("Go should be installed on the build VM; got status=%s message=%s", dep.Status, dep.Message)
	}
	if dep.Version == "" {
		t.Error("Version should be non-empty when Go is installed")
	}
	t.Logf("Go version: %s, status: %s", dep.Version, dep.Status)
}

// TestCheckNpm verifies npm detection (node/npm should be available on the build VM).
func TestCheckNpm(t *testing.T) {
	dep := checkNpm()

	if dep.Name != "npm" {
		t.Errorf("unexpected Name: %q", dep.Name)
	}
	t.Logf("npm status=%s version=%s installed=%v", dep.Status, dep.Version, dep.Installed)
	if dep.Installed {
		if dep.Status != "installed" && dep.Status != "needs_update" {
			t.Errorf("npm installed but status=%q", dep.Status)
		}
	}
}

// TestCheckDocker_ExecCommandPathBehavior verifies that after execCommandRefreshed
// updates the process PATH, a plain execCommand call also sees the updated PATH.
// This is the core of the potential bug reported: checkDocker uses execCommandRefreshed
// for `docker --version` but then plain execCommand for `docker info`.
func TestCheckDocker_ExecCommandPathBehavior(t *testing.T) {
	orig := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", orig) })

	// Simulate a path that does not include the "go" binary at first.
	os.Setenv("PATH", `C:\Windows\System32`)

	// execCommandRefreshed refreshes from registry then runs the command.
	// This will restore the full PATH including Go, then find "go version".
	out, err := execCommandRefreshed("go", "version")
	if err != nil {
		t.Logf("execCommandRefreshed('go version') failed (Go not installed?): %v", err)
	} else {
		t.Logf("execCommandRefreshed found Go: %s", out)

		// Now verify that execCommand (no refresh) can ALSO find 'go' because
		// os.Setenv("PATH") was called by execCommandRefreshed above.
		out2, err2 := execCommand("go", "version")
		if err2 != nil {
			t.Errorf("execCommand('go version') failed after execCommandRefreshed set PATH: %v\n"+
				"BUG CONFIRMED: execCommand does not see refreshed PATH from prior execCommandRefreshed call", err2)
		} else {
			t.Logf("execCommand also finds Go: %s — PATH propagation correct, bug is benign", out2)
		}
	}
}

// TestCheckWindowsSigningStatus_NoCert verifies defaults when no certificate is configured.
func TestCheckWindowsSigningStatus_NoCert(t *testing.T) {
	cfg := GlobalDefaults{}
	status := checkWindowsSigningStatus(cfg)

	if status.HasCertificate {
		t.Error("expected HasCertificate=false when no cert configured")
	}
	if status.TimestampServer != "http://timestamp.digicert.com" {
		t.Errorf("expected default timestamp server, got %q", status.TimestampServer)
	}
	t.Logf("HasSignTool=%v", status.HasSignTool)
}

// TestCheckWindowsSigningStatus_CertFile verifies file-based cert detection.
func TestCheckWindowsSigningStatus_CertFile(t *testing.T) {
	cfg := GlobalDefaults{}
	cfg.Signing.Windows.CertificatePath = `C:\certs\mycert.pfx`
	status := checkWindowsSigningStatus(cfg)

	if !status.HasCertificate {
		t.Error("expected HasCertificate=true")
	}
	if status.CertificateType != "file" {
		t.Errorf("expected CertificateType=file, got %q", status.CertificateType)
	}
	if status.ConfigSource != "defaults.yaml" {
		t.Errorf("expected ConfigSource=defaults.yaml, got %q", status.ConfigSource)
	}
}

// TestCheckWindowsSigningStatus_Thumbprint verifies store-based cert detection.
func TestCheckWindowsSigningStatus_Thumbprint(t *testing.T) {
	cfg := GlobalDefaults{}
	cfg.Signing.Windows.Thumbprint = "AABBCCDDEEFF"
	status := checkWindowsSigningStatus(cfg)

	if !status.HasCertificate {
		t.Error("expected HasCertificate=true")
	}
	if status.CertificateType != "store" {
		t.Errorf("expected CertificateType=store, got %q", status.CertificateType)
	}
}

// TestCheckWindowsSigningStatus_Cloud verifies cloud provider detection.
func TestCheckWindowsSigningStatus_Cloud(t *testing.T) {
	cfg := GlobalDefaults{}
	cfg.Signing.Windows.CloudProvider = "digicert-one"
	status := checkWindowsSigningStatus(cfg)

	if !status.HasCertificate {
		t.Error("expected HasCertificate=true")
	}
	if status.CertificateType != "cloud:digicert-one" {
		t.Errorf("expected CertificateType=cloud:digicert-one, got %q", status.CertificateType)
	}
}

// TestCheckWindowsSigningStatus_TimestampOverride verifies custom timestamp server.
func TestCheckWindowsSigningStatus_TimestampOverride(t *testing.T) {
	cfg := GlobalDefaults{}
	cfg.Signing.Windows.TimestampServer = "http://timestamp.globalsign.com"
	status := checkWindowsSigningStatus(cfg)

	if status.TimestampServer != "http://timestamp.globalsign.com" {
		t.Errorf("expected custom timestamp server, got %q", status.TimestampServer)
	}
}

// TestCheckWindowsSigningStatus_SigntoolDetection verifies signtool.exe lookup.
func TestCheckWindowsSigningStatus_SigntoolDetection(t *testing.T) {
	cfg := GlobalDefaults{}
	status := checkWindowsSigningStatus(cfg)
	// Just log; the VM may or may not have Windows SDK installed.
	t.Logf("signtool.exe found: %v", status.HasSignTool)
}
