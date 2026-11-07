//go:build linux && !android && !server

package application

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestSingleInstanceNames(t *testing.T) {
	tests := []struct {
		name     string
		uniqueID string
		busName  string
		iface    string
		path     string
	}{
		{
			name:     "reverse dns id is kept verbatim in the bus name",
			uniqueID: "com.myapp.myapplication",
			busName:  "com.myapp.myapplication.SingleInstance",
			iface:    "com.myapp.myapplication.SingleInstance",
			path:     "/com/myapp/myapplication/SingleInstance",
		},
		{
			// Hyphens are legal in a bus name but not in an interface name or
			// an object path, so only the bus name keeps them.
			name:     "hyphens survive in the bus name only",
			uniqueID: "net.my-company.my-app",
			busName:  "net.my-company.my-app.SingleInstance",
			iface:    "net.my_company.my_app.SingleInstance",
			path:     "/net/my_company/my_app/SingleInstance",
		},
		{
			name:     "single element id still yields two name elements",
			uniqueID: "myapp",
			busName:  "myapp.SingleInstance",
			iface:    "myapp.SingleInstance",
			path:     "/myapp/SingleInstance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			busName, interfaceName, path, err := singleInstanceNames(tt.uniqueID)
			if err != nil {
				t.Fatalf("singleInstanceNames(%q) returned error: %v", tt.uniqueID, err)
			}
			if busName != tt.busName {
				t.Errorf("bus name = %q, want %q", busName, tt.busName)
			}
			if interfaceName != tt.iface {
				t.Errorf("interface name = %q, want %q", interfaceName, tt.iface)
			}
			if path != tt.path {
				t.Errorf("object path = %q, want %q", path, tt.path)
			}
		})
	}
}

// The bus name has to stay under the app id, which is what lets a sandboxed
// build claim it without being granted the whole session bus.
func TestSingleInstanceNamesArePrefixedByUniqueID(t *testing.T) {
	const uniqueID = "net.koofr.stage"

	busName, _, _, err := singleInstanceNames(uniqueID)
	if err != nil {
		t.Fatalf("singleInstanceNames(%q) returned error: %v", uniqueID, err)
	}
	if got, want := busName[:len(uniqueID)], uniqueID; got != want {
		t.Errorf("bus name %q is not prefixed by the UniqueID %q", busName, want)
	}
}

func TestSingleInstanceNamesRejectsInvalidIDs(t *testing.T) {
	tests := []struct {
		name     string
		uniqueID string
	}{
		{"empty element", "com..myapp"},
		{"trailing dot", "com.myapp."},
		{"element starting with a digit", "com.1myapp"},
		{"character D-Bus does not allow", "com.my app"},
		{"slash", "com/myapp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, _, err := singleInstanceNames(tt.uniqueID); err == nil {
				t.Errorf("singleInstanceNames(%q) succeeded, want an error", tt.uniqueID)
			}
		})
	}
}

func TestSameExecutableTellsThisBinaryFromAnother(t *testing.T) {
	if same, known := sameExecutable("/proc", os.Getpid()); !known || !same {
		t.Fatalf("this process: same=%v known=%v", same, known)
	}
	sleep, err := exec.LookPath("sleep")
	if err != nil {
		t.Skip(err)
	}
	cmd := exec.Command(sleep, "30")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
	if same, known := sameExecutable("/proc", cmd.Process.Pid); !known || same {
		t.Fatalf("another executable: same=%v known=%v", same, known)
	}
	if _, known := sameExecutable("/proc", 1<<30); known {
		t.Fatal("a missing process was judged")
	}
}

func TestCheckCallerRefusesOnlyAKnownOtherExecutable(t *testing.T) {
	if err := checkCaller("/proc", os.Getpid()); err != nil {
		t.Fatalf("this process refused: %v", err)
	}
	if err := checkCaller("/proc", 1<<30); err != nil {
		t.Fatalf("an unknown caller refused: %v", err)
	}
	sleep, err := exec.LookPath("sleep")
	if err != nil {
		t.Skip(err)
	}
	cmd := exec.Command(sleep, "30")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
	if err := checkCaller("/proc", cmd.Process.Pid); err == nil {
		t.Fatal("another executable was accepted")
	}
}

func TestSameBuildAcceptsTwoMountsOfOneAppImage(t *testing.T) {
	stamp := time.Date(2026, 11, 7, 12, 0, 0, 0, time.UTC)
	copyOf := func(dir, body string, at time.Time) os.FileInfo {
		p := filepath.Join(t.TempDir(), dir, "app")
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(p, at, at); err != nil {
			t.Fatal(err)
		}
		fi, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		return fi
	}
	a := copyOf(".mount_appA", "build-1", stamp)
	if !sameBuild(a, copyOf(".mount_appB", "build-1", stamp)) {
		t.Fatal("two mounts of one image were judged different")
	}
	if sameBuild(a, copyOf(".mount_appC", "build-2", stamp.Add(time.Hour))) {
		t.Fatal("another build was judged the same")
	}
}
