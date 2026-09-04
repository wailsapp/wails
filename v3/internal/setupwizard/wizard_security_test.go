package setupwizard

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestGetSafeCommand(t *testing.T) {
	tests := []struct {
		name            string
		parts           []string
		wantSafeCmd     string
		wantElevatedCmd string
		wantArgs        []string
		wantOK          bool
	}{
		{"apt install", []string{"apt", "install", "libgtk-4-dev"}, "apt", "", []string{"install", "libgtk-4-dev"}, true},
		{"sudo apt install", []string{"sudo", "apt", "install", "libgtk-4-dev"}, "sudo", "apt", []string{"install", "libgtk-4-dev"}, true},
		{"sudo apt-get with confirmation", []string{"sudo", "apt-get", "install", "-y", "openjdk-21-jdk"}, "sudo", "apt-get", []string{"install", "-y", "openjdk-21-jdk"}, true},
		{"sudo pacman with confirmation", []string{"sudo", "pacman", "-S", "--noconfirm", "jdk-openjdk"}, "sudo", "pacman", []string{"-S", "--noconfirm", "jdk-openjdk"}, true},
		{"brew cask", []string{"brew", "install", "--cask", "android-commandlinetools"}, "brew", "", []string{"install", "--cask", "android-commandlinetools"}, true},
		{"sdkmanager package with semicolon", []string{"sdkmanager", "--install", "ndk;26.3.11579264"}, "sdkmanager", "", []string{"--install", "ndk;26.3.11579264"}, true},
		{"xcode runtime", []string{"xcodebuild", "-downloadPlatform", "iOS"}, "xcodebuild", "", []string{"-downloadPlatform", "iOS"}, true},
		{"sudo option bypass", []string{"sudo", "-u", "apt", "bash"}, "", "", nil, false},
		{"sudo shell", []string{"sudo", "bash", "-c", "id"}, "", "", nil, false},
		{"apt option injection", []string{"apt", "install", "--option=Dir::Etc::sourcelist=/tmp/evil", "pkg"}, "", "", nil, false},
		{"shell command", []string{"sh", "-c", "id"}, "", "", nil, false},
		{"empty", nil, "", "", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			safeCmd, elevatedCmd, args, ok := getSafeCommand(tt.parts)
			if ok != tt.wantOK || safeCmd != tt.wantSafeCmd || elevatedCmd != tt.wantElevatedCmd || !reflect.DeepEqual(args, tt.wantArgs) {
				t.Fatalf("getSafeCommand(%v) = (%q, %q, %v, %v), want (%q, %q, %v, %v)", tt.parts, safeCmd, elevatedCmd, args, ok, tt.wantSafeCmd, tt.wantElevatedCmd, tt.wantArgs, tt.wantOK)
			}
		})
	}
}

func TestInstallHandlerRejectsCrossOriginRequests(t *testing.T) {
	w := New()
	w.origin = "http://127.0.0.1:4321"
	req := httptest.NewRequest(http.MethodPost, "/api/dependencies/install", strings.NewReader(`{"command":"apt install pkg"}`))
	req.Header.Set("Origin", "https://attacker.example")
	rr := httptest.NewRecorder()

	w.handleInstallDependency(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusForbidden)
	}
}

func TestInstallHandlerAllowsSameOriginValidation(t *testing.T) {
	w := New()
	w.origin = "http://127.0.0.1:4321"
	req := httptest.NewRequest(http.MethodPost, "/api/dependencies/install", strings.NewReader(`{"command":"bash -c id"}`))
	req.Header.Set("Origin", w.origin)
	rr := httptest.NewRecorder()

	w.handleInstallDependency(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "approved dependency installer") {
		t.Fatalf("response did not reject the command: %s", rr.Body.String())
	}
}
