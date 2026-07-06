package notifications

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateNotificationOptions_RequiredFields(t *testing.T) {
	cases := []struct {
		name    string
		opts    NotificationOptions
		wantErr string
	}{
		{
			name:    "empty id",
			opts:    NotificationOptions{Title: "ok"},
			wantErr: "notification ID cannot be empty",
		},
		{
			name:    "empty title",
			opts:    NotificationOptions{ID: "n1"},
			wantErr: "notification title cannot be empty",
		},
		{
			name: "ok",
			opts: NotificationOptions{ID: "n1", Title: "ok"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateNotificationOptions(tc.opts)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("got %v, want error containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestValidateNotificationOptions_InterruptionLevel(t *testing.T) {
	for _, level := range []string{
		"", InterruptionLevelPassive, InterruptionLevelActive,
		InterruptionLevelTimeSensitive, InterruptionLevelCritical,
	} {
		opts := NotificationOptions{ID: "n", Title: "t", InterruptionLevel: level}
		if err := validateNotificationOptions(opts); err != nil {
			t.Errorf("level %q should be valid, got %v", level, err)
		}
	}

	bad := NotificationOptions{ID: "n", Title: "t", InterruptionLevel: "shouty"}
	err := validateNotificationOptions(bad)
	if err == nil || !strings.Contains(err.Error(), "invalid interruption level") {
		t.Fatalf("expected invalid interruption level error, got %v", err)
	}
}

func TestValidateNotificationOptions_Schedule(t *testing.T) {
	cases := []struct {
		name    string
		sched   *NotificationSchedule
		wantErr string
	}{
		{name: "nil schedule"},
		{name: "delay only", sched: &NotificationSchedule{DelaySeconds: 30}},
		{name: "at only", sched: &NotificationSchedule{At: 1717181600}},
		{
			name:    "both set",
			sched:   &NotificationSchedule{DelaySeconds: 30, At: 1717181600},
			wantErr: "mutually exclusive",
		},
		{
			name:    "neither set",
			sched:   &NotificationSchedule{},
			wantErr: "must set either delaySeconds or at",
		},
		{
			name:    "negative delay",
			sched:   &NotificationSchedule{DelaySeconds: -5},
			wantErr: "cannot be negative",
		},
		{
			name:    "negative at",
			sched:   &NotificationSchedule{At: -1},
			wantErr: "cannot be negative",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := NotificationOptions{ID: "n", Title: "t", Schedule: tc.sched}
			err := validateNotificationOptions(opts)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("got %v, want error containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestValidateNotificationOptions_Attachments(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "image.png")
	if err := os.WriteFile(good, []byte("not really a png"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Run("present file passes", func(t *testing.T) {
		opts := NotificationOptions{ID: "n", Title: "t", Attachments: []NotificationAttachment{{Path: good}}}
		if err := validateNotificationOptions(opts); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("missing file fails", func(t *testing.T) {
		opts := NotificationOptions{ID: "n", Title: "t", Attachments: []NotificationAttachment{{Path: filepath.Join(dir, "absent.png")}}}
		err := validateNotificationOptions(opts)
		if err == nil || !strings.Contains(err.Error(), "not accessible") {
			t.Fatalf("got %v, want not-accessible error", err)
		}
		// Error must not expose raw OS error text (no filesystem oracle).
		if strings.Contains(err.Error(), "no such file") || strings.Contains(err.Error(), "cannot find") {
			t.Fatalf("error leaks OS-level details: %v", err)
		}
	})

	t.Run("empty path fails", func(t *testing.T) {
		opts := NotificationOptions{ID: "n", Title: "t", Attachments: []NotificationAttachment{{Path: ""}}}
		err := validateNotificationOptions(opts)
		if err == nil || !strings.Contains(err.Error(), "cannot be empty") {
			t.Fatalf("got %v, want empty-path error", err)
		}
	})

	t.Run("relative path fails", func(t *testing.T) {
		opts := NotificationOptions{ID: "n", Title: "t", Attachments: []NotificationAttachment{{Path: "relative/image.png"}}}
		err := validateNotificationOptions(opts)
		if err == nil || !strings.Contains(err.Error(), "absolute") {
			t.Fatalf("got %v, want absolute-path error", err)
		}
	})

	t.Run("relative file:// URL fails", func(t *testing.T) {
		opts := NotificationOptions{ID: "n", Title: "t", Attachments: []NotificationAttachment{{Path: "file://relative/image.png"}}}
		err := validateNotificationOptions(opts)
		if err == nil || !strings.Contains(err.Error(), "absolute") {
			t.Fatalf("got %v, want absolute-path error", err)
		}
	})

	// macOS UNNotificationAttachment accepts file:// URLs, and the package
	// godoc on NotificationAttachment.Path documents them. The validator
	// must not reject the URL form by trying to os.Stat it as a literal
	// path.
	t.Run("file:// URL passes", func(t *testing.T) {
		opts := NotificationOptions{
			ID:    "n",
			Title: "t",
			Attachments: []NotificationAttachment{
				{Path: "file://" + good},
			},
		}
		if err := validateNotificationOptions(opts); err != nil {
			t.Fatalf("unexpected error for file:// URL: %v", err)
		}
	})
}
