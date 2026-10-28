package notifications

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
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

func TestConcurrentCategoryChangesPersist(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	notifier := &linuxNotifier{appName: "wails-concurrent-category-test", categories: make(map[string]NotificationCategory)}
	const workersCount = 8
	const perWorker = 4
	const count = workersCount * perWorker
	start := make(chan struct{})
	var workers sync.WaitGroup
	for worker := range workersCount {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			for offset := range perWorker {
				i := worker*perWorker + offset
				category := NotificationCategory{ID: fmt.Sprintf("category-%d", i), Actions: []NotificationAction{{ID: "action", Title: strings.Repeat("label", 1024)}}}
				if err := notifier.RegisterNotificationCategory(category); err != nil {
					t.Errorf("register %s: %v", category.ID, err)
				}
				if i%2 == 0 {
					if err := notifier.RemoveNotificationCategory(category.ID); err != nil {
						t.Errorf("remove %s: %v", category.ID, err)
					}
				}
			}
		}()
	}
	close(start)
	workers.Wait()
	reloaded := &linuxNotifier{appName: notifier.appName}
	if err := reloaded.loadCategories(); err != nil {
		t.Fatal(err)
	}
	if len(reloaded.categories) != count/2 {
		t.Fatalf("persisted %d categories, want %d", len(reloaded.categories), count/2)
	}
	if !reflect.DeepEqual(reloaded.categories, notifier.categories) {
		t.Error("saved categories differ from the final in-memory categories")
	}
}
