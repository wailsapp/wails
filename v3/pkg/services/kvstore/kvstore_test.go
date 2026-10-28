package kvstore

import (
	"encoding/json"
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

func TestGetStoreSnapshot(t *testing.T) {
	store := New()
	if err := store.Set("key", "original"); err != nil {
		t.Fatal(err)
	}
	snapshot := store.Get("").(map[string]any)
	if err := store.Set("key", "updated"); err != nil {
		t.Fatal(err)
	}
	if got := snapshot["key"]; got != "original" {
		t.Errorf("snapshot changed after Set: %v", got)
	}
	snapshot["extra"] = "outside"
	delete(snapshot, "key")
	if got := store.Get("key"); got != "updated" {
		t.Errorf("snapshot deletion changed store: %v", got)
	}
	if got := store.Get("extra"); got != nil {
		t.Errorf("snapshot insertion changed store: %v", got)
	}
}

func TestGetStoreSnapshotConcurrentSerialization(t *testing.T) {
	store := New()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := range 1000 {
			if err := store.Set("key", i); err != nil {
				t.Error(err)
				return
			}
		}
	}()
	for range 1000 {
		if _, err := json.Marshal(store.Get("")); err != nil {
			t.Error(err)
			break
		}
	}
	<-done
}
