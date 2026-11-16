package appbundle

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"howett.net/plist"
)

func TestMetadataOverridesOnlyStagedPlistBeforePrepare(t *testing.T) {
	for _, format := range []int{plist.XMLFormat, plist.BinaryFormat} {
		options := fixture(t)
		original, err := plist.Marshal(map[string]any{"CFBundleExecutable": "app", "CFBundleShortVersionString": "old", "Custom": []string{"preserved"}}, format)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(options.InfoPlist, original, 0400); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(options.InfoPlist, 0400); err != nil {
			t.Fatal(err)
		}
		options.InfoValues = map[string]string{"CFBundleShortVersionString": "1.2.3"}
		prepared := false
		options.Prepare = func(ctx context.Context, path string) error {
			data, err := os.ReadFile(filepath.Join(path, "Contents/Info.plist"))
			if err != nil {
				return err
			}
			var info map[string]any
			gotFormat, err := plist.Unmarshal(data, &info)
			if err != nil {
				return err
			}
			if gotFormat != format || info["CFBundleShortVersionString"] != "1.2.3" || info["Custom"] == nil {
				t.Errorf("format=%d info=%v", gotFormat, info)
			}
			prepared = true
			return nil
		}
		if err := Create(context.Background(), options); err != nil {
			t.Fatal(err)
		}
		if !prepared {
			t.Fatal("prepare skipped")
		}
		data, err := os.ReadFile(options.InfoPlist)
		if err != nil || string(data) != string(original) {
			t.Fatal("source plist changed")
		}
	}
}
