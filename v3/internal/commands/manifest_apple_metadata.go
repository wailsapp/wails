package commands

import (
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
	"howett.net/plist"
)

func (h *manifestHandler) applyMicrophoneMetadata(staged string, spec pipeline.AssetsSpec) error {
	if spec.TargetOS != "darwin" || spec.MicrophoneUsageDescription == "" {
		return nil
	}
	path := filepath.Join(staged, "darwin", "Info.plist")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	info := map[string]any{}
	format, err := plist.Unmarshal(data, &info)
	if err != nil {
		return err
	}
	info["NSMicrophoneUsageDescription"] = spec.MicrophoneUsageDescription
	data, err = plist.Marshal(info, format)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	if h.config.Signing.Darwin.Entitlements != "" {
		return nil
	}
	data, err = plist.Marshal(map[string]any{"com.apple.security.device.audio-input": true}, plist.XMLFormat)
	if err != nil {
		return err
	}
	destination := filepath.Join(staged, "signing", "entitlements.plist")
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return err
	}
	return os.WriteFile(destination, data, 0644)
}
