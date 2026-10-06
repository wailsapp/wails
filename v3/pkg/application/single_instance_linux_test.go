//go:build linux && !android && !server

package application

import (
	"fmt"
	"os"
	"testing"

	"github.com/godbus/dbus/v5"
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

// A quitting app must give up the bus name before its shutdown work, or an
// instance started meanwhile hands its launch to an app that is exiting.
func TestLinuxLockReleaseFreesTheBusName(t *testing.T) {
	other, err := dbus.ConnectSessionBus()
	if err != nil {
		t.Skipf("no session bus: %v", err)
	}
	defer other.Close()

	uniqueID := fmt.Sprintf("org.wails.test.release.p%d", os.Getpid())
	lock := &linuxLock{manager: &singleInstanceManager{options: &SingleInstanceOptions{}}}
	if err := lock.acquire(uniqueID); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	busName := uniqueID + ".SingleInstance"

	reply, err := other.RequestName(busName, dbus.NameFlagDoNotQueue)
	if err != nil || reply != dbus.RequestNameReplyExists {
		t.Fatalf("while held: reply %d, err %v; want the name taken", reply, err)
	}

	lock.release()

	reply, err = other.RequestName(busName, dbus.NameFlagDoNotQueue)
	if err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		t.Fatalf("after release: reply %d, err %v; want the name free", reply, err)
	}
}
