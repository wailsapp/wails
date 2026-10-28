package notifications

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNullCategoriesThenRegister(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	const appName = "wails-category-test"
	dir := filepath.Join(root, appName)
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(dir, "notification-categories.json")
	if err := os.WriteFile(file, []byte("null"), 0600); err != nil {
		t.Fatal(err)
	}
	notifier := &linuxNotifier{appName: appName, categories: map[string]NotificationCategory{"old": {ID: "old"}}}
	if err := notifier.loadCategories(); err != nil {
		t.Fatal(err)
	}
	if len(notifier.categories) != 0 {
		t.Fatalf("old categories retained: %v", notifier.categories)
	}
	if err := notifier.RegisterNotificationCategory(NotificationCategory{ID: "new"}); err != nil {
		t.Fatal(err)
	}
	reloaded := &linuxNotifier{appName: appName}
	if err := reloaded.loadCategories(); err != nil {
		t.Fatal(err)
	}
	if got := reloaded.categories["new"].ID; got != "new" {
		t.Errorf("reloaded category ID=%q, want new", got)
	}
}
