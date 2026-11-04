package appbundle

import (
	"os"

	"howett.net/plist"
)

func overrideInfo(path string, overrides map[string]string) error {
	if len(overrides) == 0 {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	info := map[string]any{}
	format, err := plist.Unmarshal(data, &info)
	if err != nil {
		return err
	}
	for key, value := range overrides {
		info[key] = value
	}
	data, err = plist.Marshal(info, format)
	if err != nil {
		return err
	}
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	// Replace only the staged copy, including read-only source plist files.
	if err := os.Remove(path); err != nil {
		return err
	}
	return os.WriteFile(path, data, stat.Mode().Perm())
}
