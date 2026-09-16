package application

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ErrSoundNotSupported is returned by SoundManager methods on platforms that
// cannot play the requested sound.
var ErrSoundNotSupported = errors.New("sound playback is not supported on this platform")

// soundFileExtensions are the audio file types accepted by SystemSounds and
// listed without their extension.
var soundFileExtensions = map[string]bool{
	".aiff": true,
	".aif":  true,
	".wav":  true,
	".caf":  true,
	".m4a":  true,
	".mp3":  true,
}

// SoundManager plays system sounds and alert tones.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type SoundManager struct {
	app *App
}

func newSoundManager(app *App) *SoundManager {
	return &SoundManager{app: app}
}

// Beep plays the system alert sound.
//
// macOS: NSBeep. Windows: MessageBeep. Other platforms: no-op.
func (s *SoundManager) Beep() {
	soundBeep()
}

// Play plays a named system sound or an audio file and returns immediately.
//
// name is either a system sound name such as "Glass", "Ping" or "Submarine"
// (see SystemSounds), or an absolute path to an audio file. Relative paths
// are rejected.
//
// macOS: NSSound, so any format Core Audio can decode is accepted. Windows:
// absolute paths play through PlaySound (WAV files only); names are treated
// as registry sound aliases such as "SystemAsterisk". Other platforms return
// ErrSoundNotSupported.
func (s *SoundManager) Play(name string) error {
	isPath, err := validateSoundName(name)
	if err != nil {
		return err
	}
	if isPath {
		if _, statErr := os.Stat(name); statErr != nil {
			return fmt.Errorf("sound file %q: %w", name, statErr)
		}
		return soundPlayFile(name)
	}
	return soundPlayNamed(name)
}

// PlayData plays audio from memory and returns immediately. The data must be
// a complete audio file (for example WAV or AIFF bytes).
//
// macOS only; other platforms return ErrSoundNotSupported.
func (s *SoundManager) PlayData(data []byte) error {
	if len(data) == 0 {
		return errors.New("sound data is empty")
	}
	return soundPlayData(data)
}

// SystemSounds returns the names accepted by Play for the built-in system
// sounds, sorted alphabetically.
//
// macOS: the contents of /System/Library/Sounds ("Basso", "Blow", "Bottle",
// "Frog", "Funk", "Glass", "Hero", "Morse", "Ping", "Pop", "Purr", "Sosumi",
// "Submarine", "Tink"). Other platforms return nil.
func (s *SoundManager) SystemSounds() []string {
	dir := systemSoundsDir()
	if dir == "" {
		return nil
	}
	return listSoundsInDir(dir)
}

// validateSoundName checks a Play argument. It returns true when name is an
// absolute file path rather than a system sound name.
func validateSoundName(name string) (isPath bool, err error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return false, errors.New("sound name is empty")
	}
	if filepath.IsAbs(name) {
		return true, nil
	}
	if strings.ContainsAny(name, `/\`) || strings.Contains(name, "..") {
		return false, fmt.Errorf("sound %q: relative paths are not supported, use a system sound name or an absolute path", name)
	}
	return false, nil
}

// listSoundsInDir returns the sorted, extension-less names of the audio
// files directly inside dir. Missing or unreadable directories yield nil.
func listSoundsInDir(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	seen := make(map[string]bool)
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if !soundFileExtensions[ext] {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
