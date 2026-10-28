package kvstore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEmptyStoreRemainsWritable(t *testing.T) {
	for _, content := range []string{"null", " \nnull\t", "{}", ""} {
		t.Run(content, func(t *testing.T) {
			filename := filepath.Join(t.TempDir(), "store.json")
			if err := os.WriteFile(filename, []byte(content), 0600); err != nil {
				t.Fatal(err)
			}
			store := New()
			if err := store.Set("old", "old value"); err != nil {
				t.Fatal(err)
			}
			store.Configure(&Config{Filename: filename, AutoSave: true})
			if err := store.Load(); err != nil {
				t.Fatal(err)
			}
			if store.Get("old") != nil {
				t.Fatal("load retained the previous store contents")
			}
			if err := store.Set("new", "new value"); err != nil {
				t.Fatal(err)
			}
			reloaded := NewWithConfig(&Config{Filename: filename})
			if err := reloaded.Load(); err != nil {
				t.Fatal(err)
			}
			if got := reloaded.Get("new"); got != "new value" {
				t.Errorf("saved value=%v, want new value", got)
			}
		})
	}
}
