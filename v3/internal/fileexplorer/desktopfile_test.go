package fileexplorer

import (
	"strings"
	"testing"
)

func TestParseDesktopReader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantExec string
		wantErr  bool
	}{
		{
			name: "simple desktop file",
			input: `[Desktop Entry]
Name=Files
Exec=nautilus %U
Icon=org.gnome.Nautilus
`,
			wantExec: "nautilus %U",
		},
		{
			name: "exec with full path",
			input: `[Desktop Entry]
Name=1Password
Exec=/opt/1Password/1password %U
`,
			wantExec: "/opt/1Password/1password %U",
		},
		{
			name: "exec without arguments",
			input: `[Desktop Entry]
Name=Btop
Exec=btop
Terminal=true
`,
			wantExec: "btop",
		},
		{
			name: "exec with spaces in path",
			input: `[Desktop Entry]
Name=My App
Exec="/path/with spaces/myapp" %f
`,
			wantExec: `"/path/with spaces/myapp" %f`,
		},
		{
			name: "comments are ignored",
			input: `# This is a comment
[Desktop Entry]
# Another comment
Name=Files
Exec=nautilus
# Comment after
`,
			wantExec: "nautilus",
		},
		{
			name: "empty lines are ignored",
			input: `

[Desktop Entry]

Name=Files

Exec=nautilus

`,
			wantExec: "nautilus",
		},
		{
			name: "key before section is ignored",
			input: `Exec=ignored
[Desktop Entry]
Exec=nautilus
`,
			wantExec: "nautilus",
		},
		{
			name: "other sections after Desktop Entry are ignored",
			input: `[Desktop Entry]
Exec=nautilus --new-window %U
Icon=nautilus

[Desktop Action new-window]
Name=New Window
Exec=nautilus --new-window
`,
			wantExec: "nautilus --new-window %U",
		},
		{
			name: "section before Desktop Entry is ignored",
			input: `[Some Other Section]
Exec=ignored

[Desktop Entry]
Exec=nautilus
`,
			wantExec: "nautilus",
		},
		{
			name: "case sensitive section name",
			input: `[desktop entry]
Exec=ignored

[Desktop Entry]
Exec=correct
`,
			wantExec: "correct",
		},
		{
			name: "case sensitive key name",
			input: `[Desktop Entry]
exec=ignored
EXEC=also ignored
Exec=correct
`,
			wantExec: "correct",
		},
		{
			name: "value with equals sign",
			input: `[Desktop Entry]
Exec=env VAR=value myapp
`,
			wantExec: "env VAR=value myapp",
		},
		{
			name: "value with multiple equals signs",
			input: `[Desktop Entry]
Exec=env A=1 B=2 C=3 myapp
`,
			wantExec: "env A=1 B=2 C=3 myapp",
		},
		{
			name: "localized keys are separate",
			input: `[Desktop Entry]
Name[en]=Files
Name=Default Files
Exec[en]=ignored
Exec=nautilus
`,
			wantExec: "nautilus",
		},
		{
			name: "whitespace in section header",
			input: `[Desktop Entry]
Exec=nautilus
`,
			wantExec: "nautilus",
		},
		{
			name: "no exec key",
			input: `[Desktop Entry]
Name=Files
Icon=nautilus
`,
			wantExec: "",
		},
		{
			name:     "empty file",
			input:    ``,
			wantExec: "",
		},
		{
			name: "only comments",
			input: `# Comment 1
# Comment 2
`,
			wantExec: "",
		},
		{
			name: "no Desktop Entry section",
			input: `[Other Section]
Exec=ignored
`,
			wantExec: "",
		},
		{
			name: "real nautilus desktop file structure",
			input: `[Desktop Entry]
Name[en_CA]=Files
Name[en_GB]=Files
Name=Files
Comment=Access and organize files
Keywords=folder;manager;explore;disk;filesystem;nautilus;
Exec=nautilus --new-window %U
Icon=org.gnome.Nautilus
Terminal=false
Type=Application
DBusActivatable=true
StartupNotify=true
Categories=GNOME;GTK;Utility;Core;FileManager;
MimeType=inode/directory;application/x-7z-compressed;
X-GNOME-UsesNotifications=true
Actions=new-window;

[Desktop Action new-window]
Name=New Window
Exec=nautilus --new-window
`,
			wantExec: "nautilus --new-window %U",
		},
		{
			name: "thunar style",
			input: `[Desktop Entry]
Version=1.0
Name=Thunar File Manager
Exec=thunar %F
Icon=Thunar
Type=Application
Categories=System;FileTools;FileManager;
`,
			wantExec: "thunar %F",
		},
		{
			name: "dolphin style",
			input: `[Desktop Entry]
Type=Application
Exec=dolphin %u
Icon=system-file-manager
Name=Dolphin
GenericName=File Manager
`,
			wantExec: "dolphin %u",
		},
		{
			name: "pcmanfm style",
			input: `[Desktop Entry]
Type=Application
Name=PCManFM
GenericName=File Manager
Exec=pcmanfm %U
Icon=system-file-manager
`,
			wantExec: "pcmanfm %U",
		},
		{
			name: "exec with environment variable",
			input: `[Desktop Entry]
Exec=env GDK_BACKEND=x11 nautilus %U
`,
			wantExec: "env GDK_BACKEND=x11 nautilus %U",
		},
		{
			name: "trailing whitespace in value preserved",
			input: `[Desktop Entry]
Exec=nautilus   
`,
			wantExec: "nautilus   ",
		},
		{
			name: "leading whitespace in key",
			input: `[Desktop Entry]
  Exec=nautilus
`,
			wantExec: "nautilus",
		},
		{
			name: "space around equals",
			input: `[Desktop Entry]
Exec = nautilus
`,
			wantExec: " nautilus", // We trim the key, value starts after =
		},
		{
			name: "line without equals is ignored",
			input: `[Desktop Entry]
InvalidLine
Exec=nautilus
AnotherInvalidLine
`,
			wantExec: "nautilus",
		},
		{
			name: "UTF-8 in exec path",
			input: `[Desktop Entry]
Exec=/usr/bin/文件管理器 %U
`,
			wantExec: "/usr/bin/文件管理器 %U",
		},
		{
			name: "special characters in exec",
			input: `[Desktop Entry]
Exec=sh -c "echo 'hello world' && nautilus %U"
`,
			wantExec: `sh -c "echo 'hello world' && nautilus %U"`,
		},
		{
			name: "multiple Desktop Entry sections (invalid file, last value wins)",
			input: `[Desktop Entry]
Exec=first

[Desktop Entry]
Exec=second
`,
			wantExec: "second", // Invalid file, but we handle it gracefully
		},
		{
			name: "very long exec line",
			input: `[Desktop Entry]
Exec=` + strings.Repeat("a", 1000) + `
`,
			wantExec: strings.Repeat("a", 1000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := ParseDesktopReader(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDesktopReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if entry.Exec != tt.wantExec {
				t.Errorf("ParseDesktopReader() Exec = %q, want %q", entry.Exec, tt.wantExec)
			}
		})
	}
}

func TestParseDesktopReader_LineScanning(t *testing.T) {
	// Test that we handle lines at the scanner's limit (64KB default)
	// bufio.Scanner returns an error for lines > 64KB, which is acceptable
	// since real .desktop files never have lines that long

	t.Run("line at buffer limit returns error", func(t *testing.T) {
		// Create a line that exceeds the buffer size (64KB)
		longValue := strings.Repeat("x", 65536)
		input := "[Desktop Entry]\nExec=" + longValue + "\n"

		_, err := ParseDesktopReader(strings.NewReader(input))
		if err == nil {
			t.Error("Expected error for line exceeding buffer size, got nil")
		}
	})

	t.Run("line under buffer limit works", func(t *testing.T) {
		// Create a line that's under the limit (should work fine)
		longValue := strings.Repeat("x", 60000)
		input := "[Desktop Entry]\nExec=" + longValue + "\n"

		entry, err := ParseDesktopReader(strings.NewReader(input))
		if err != nil {
			t.Errorf("Unexpected error for long but valid line: %v", err)
			return
		}
		if entry.Exec != longValue {
			t.Errorf("Long line not parsed correctly, got length %d, want %d", len(entry.Exec), len(longValue))
		}
	})
}

func TestParseDesktopReader_RealWorldFiles(t *testing.T) {
	// These are actual .desktop file contents from real systems
	realWorldTests := []struct {
		name     string
		content  string
		wantExec string
	}{
		{
			name: "GNOME Nautilus 43.x",
			content: `[Desktop Entry]
Name=Files
Comment=Access and organize files
Keywords=folder;manager;explore;disk;filesystem;nautilus;
Exec=nautilus --new-window %U
Icon=org.gnome.Nautilus
Terminal=false
Type=Application
DBusActivatable=true
StartupNotify=true
Categories=GNOME;GTK;Utility;Core;FileManager;
MimeType=inode/directory;application/x-7z-compressed;
Actions=new-window;

[Desktop Action new-window]
Name=New Window
Exec=nautilus --new-window`,
			wantExec: "nautilus --new-window %U",
		},
		{
			name: "KDE Dolphin",
			content: `[Desktop Entry]
Type=Application
Exec=dolphin %u
Icon=system-file-manager
Terminal=false
InitialPreference=9
Name=Dolphin
GenericName=File Manager
MimeType=inode/directory;
Categories=Qt;KDE;System;FileTools;FileManager;
Actions=new-window;

[Desktop Action new-window]
Name=Open a New Window
Exec=dolphin %u`,
			wantExec: "dolphin %u",
		},
		{
			name: "Thunar",
			content: `[Desktop Entry]
Version=1.0
Name=Thunar File Manager
GenericName=File Manager
Comment=Browse the filesystem with the file manager
Exec=thunar %F
Icon=Thunar
Terminal=false
StartupNotify=true
Type=Application
Categories=System;FileTools;FileManager;
`,
			wantExec: "thunar %F",
		},
		{
			name: "PCManFM",
			content: `[Desktop Entry]
Type=Application
Name=PCManFM
GenericName=File Manager
Comment=Browse the file system
Exec=pcmanfm %U
Icon=system-file-manager
Terminal=false
StartupNotify=true
Categories=Utility;FileManager;`,
			wantExec: "pcmanfm %U",
		},
		{
			name: "Caja (MATE)",
			content: `[Desktop Entry]
Name=Files
Comment=Access and organize files
Exec=caja %U
Icon=system-file-manager
Terminal=false
Type=Application
Categories=MATE;System;FileManager;
StartupNotify=true`,
			wantExec: "caja %U",
		},
		{
			name: "Nemo (Cinnamon)",
			content: `[Desktop Entry]
Name=Files
Comment=Access and organize files
Exec=nemo %U
Icon=folder
Terminal=false
Type=Application
StartupNotify=true
Categories=GNOME;GTK;Utility;Core;
MimeType=inode/directory;`,
			wantExec: "nemo %U",
		},
	}

	for _, tt := range realWorldTests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := ParseDesktopReader(strings.NewReader(tt.content))
			if err != nil {
				t.Fatalf("ParseDesktopReader() error = %v", err)
			}
			if entry.Exec != tt.wantExec {
				t.Errorf("ParseDesktopReader() Exec = %q, want %q", entry.Exec, tt.wantExec)
			}
		})
	}
}

// BenchmarkParseDesktopReader measures parsing performance
func BenchmarkParseDesktopReader(b *testing.B) {
	// Real Nautilus .desktop file content
	content := `[Desktop Entry]
Name=Files
Comment=Access and organize files
Keywords=folder;manager;explore;disk;filesystem;nautilus;
Exec=nautilus --new-window %U
Icon=org.gnome.Nautilus
Terminal=false
Type=Application
DBusActivatable=true
StartupNotify=true
Categories=GNOME;GTK;Utility;Core;FileManager;
MimeType=inode/directory;application/x-7z-compressed;
Actions=new-window;

[Desktop Action new-window]
Name=New Window
Exec=nautilus --new-window
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseDesktopReader(strings.NewReader(content))
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseDesktopReader_Large tests parsing a file with many localized entries
func BenchmarkParseDesktopReader_Large(b *testing.B) {
	// Simulate a desktop file with many localized Name entries (like Nautilus)
	var sb strings.Builder
	sb.WriteString("[Desktop Entry]\n")
	for i := 0; i < 100; i++ {
		sb.WriteString("Name[lang")
		sb.WriteString(strings.Repeat("x", 5))
		sb.WriteString("]=Localized Name\n")
	}
	sb.WriteString("Exec=nautilus %U\n")
	sb.WriteString("[Desktop Action new-window]\n")
	sb.WriteString("Name=New Window\n")
	sb.WriteString("Exec=nautilus\n")

	content := sb.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseDesktopReader(strings.NewReader(content))
		if err != nil {
			b.Fatal(err)
		}
	}
}
