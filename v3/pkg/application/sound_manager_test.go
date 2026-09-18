package application

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestValidateSoundName(t *testing.T) {
	absolute := filepath.Join(t.TempDir(), "click.wav")
	tests := []struct {
		name     string
		input    string
		wantPath bool
		wantErr  bool
	}{
		{name: "system sound", input: "Glass"},
		{name: "system sound with spaces trimmed", input: "  Ping  "},
		{name: "absolute path", input: absolute, wantPath: true},
		{name: "empty", input: "", wantErr: true},
		{name: "whitespace", input: "   ", wantErr: true},
		{name: "relative path", input: "sounds/click.wav", wantErr: true},
		{name: "parent traversal", input: "../Glass", wantErr: true},
		{name: "backslash", input: `sounds\click.wav`, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isPath, err := validateSoundName(test.input)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateSoundName(%q) error = %v, wantErr %v", test.input, err, test.wantErr)
			}
			if isPath != test.wantPath {
				t.Errorf("validateSoundName(%q) isPath = %v, want %v", test.input, isPath, test.wantPath)
			}
		})
	}
}

func TestSoundManagerPlayRejectsInvalidNames(t *testing.T) {
	manager := newSoundManager(nil)
	if err := manager.Play(""); err == nil {
		t.Error("Play(\"\") should fail")
	}
	if err := manager.Play("../etc/passwd"); err == nil {
		t.Error("Play with traversal should fail")
	}
	missing := filepath.Join(t.TempDir(), "missing.wav")
	if err := manager.Play(missing); err == nil {
		t.Error("Play with a missing file should fail before touching the platform")
	}
	if err := manager.PlayData(nil); err == nil {
		t.Error("PlayData(nil) should fail")
	}
}

func TestListSoundsInDir(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"Ping.aiff", "glass.WAV", "Basso.aiff", "notes.txt", "Ping.wav", "noext"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "Sub.aiff"), 0o755); err != nil {
		t.Fatal(err)
	}
	got := listSoundsInDir(dir)
	want := []string{"Basso", "Ping", "glass"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("listSoundsInDir = %v, want %v", got, want)
	}
	if got := listSoundsInDir(filepath.Join(dir, "does-not-exist")); got != nil {
		t.Errorf("missing dir should yield nil, got %v", got)
	}
}

func TestSoundManagerSystemSoundsPlatform(t *testing.T) {
	sounds := newSoundManager(nil).SystemSounds()
	if runtime.GOOS != "darwin" {
		if sounds != nil {
			t.Errorf("SystemSounds should be nil off macOS, got %v", sounds)
		}
		return
	}
	if len(sounds) == 0 {
		t.Skip("no system sounds directory on this build (server or ios tags)")
	}
	found := false
	for _, name := range sounds {
		if name == "Glass" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected Glass in system sounds, got %v", sounds)
	}
}
