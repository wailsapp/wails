package commands

import (
	"fmt"
	"os"

	"howett.net/plist"
)

// Asset catalogs produce both files and bundle metadata. Preserve explicit
// project plist values while adding the generated icon declarations they lack.
func mergeIOSAssetInfo(infoPath, partialPath string) error {
	read := func(path string) (map[string]any, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var result map[string]any
		if _, err = plist.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("read iOS plist %s: %w", path, err)
		}
		return result, nil
	}
	info, err := read(infoPath)
	if err != nil {
		return err
	}
	partial, err := read(partialPath)
	if err != nil {
		return err
	}
	if !mergeIOSPlistDefaults(info, partial) {
		return nil
	}
	data, err := plist.MarshalIndent(info, plist.XMLFormat, "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(infoPath, data, 0644)
}

func mergeIOSPlistDefaults(info, defaults map[string]any) bool {
	changed := false
	for key, value := range defaults {
		existing, present := info[key]
		if !present {
			info[key] = value
			changed = true
			continue
		}
		current, currentOK := existing.(map[string]any)
		generated, generatedOK := value.(map[string]any)
		if currentOK && generatedOK {
			changed = mergeIOSPlistDefaults(current, generated) || changed
		}
	}
	return changed
}
