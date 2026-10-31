//go:build windows || darwin || linux || freebsd || openbsd || netbsd || dragonfly

package credentialstore

import (
	"errors"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestNativeCredentialRoundTrip(t *testing.T) {
	if os.Getenv("WAILS_TEST_CREDENTIALSTORE") != "1" {
		t.Skip("set WAILS_TEST_CREDENTIALSTORE=1 to exercise the native user vault")
	}
	service := "wails-test-" + strconv.FormatInt(time.Now().UnixNano(), 36)
	const user = "credential-roundtrip"
	t.Cleanup(func() {
		if err := Delete(service, user); err != nil && !errors.Is(err, ErrNotFound) {
			t.Errorf("native credential cleanup failed: %v", err)
		}
	})
	for _, value := range []string{"", "  sample password\n🔒", "replacement"} {
		if err := Set(service, user, value); err != nil {
			t.Fatal(err)
		}
		if got, err := Get(service, user); err != nil || got != value {
			t.Fatalf("native credential roundtrip failed: %v", err)
		}
	}
	if err := Delete(service, user); err != nil {
		t.Fatal(err)
	}
	if _, err := Get(service, user); !errors.Is(err, ErrNotFound) {
		t.Fatalf("deleted credential remains available: %v", err)
	}
}
