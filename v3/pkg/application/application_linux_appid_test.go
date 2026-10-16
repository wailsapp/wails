//go:build linux && cgo && !android && !server

package application

import (
	"strings"
	"testing"
)

// Both backends derive the id the same way, so this covers the GTK4 and GTK3
// builds alike.
func TestApplicationID(t *testing.T) {
	tests := []struct {
		name    string
		options Options
		want    string
	}{
		{
			name:    "derived from Name when unset",
			options: Options{Name: "My App"},
			want:    "org.wails.my_app",
		},
		{
			name:    "hyphens are kept when derived",
			options: Options{Name: "koofr-stage"},
			want:    "org.wails.koofr-stage",
		},
		{
			name:    "empty Name falls back to wailsapp",
			options: Options{},
			want:    "org.wails.wailsapp",
		},
		{
			name: "used verbatim when set",
			options: Options{
				Name:  "My App",
				Linux: LinuxOptions{ApplicationID: "com.myapp.myapplication"},
			},
			want: "com.myapp.myapplication",
		},
		{
			// Sandboxed builds depend on this: the id has to be exactly the one
			// the packaging declares, with nothing derived from Name mixed in.
			name: "set id wins over Name entirely",
			options: Options{
				Name:  "something else",
				Linux: LinuxOptions{ApplicationID: "com.example.WailsFlatpakAppId"},
			},
			want: "com.example.WailsFlatpakAppId",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applicationID(tt.options)
			if err != nil {
				t.Fatalf("applicationID() returned an unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("applicationID() = %q, want %q", got, tt.want)
			}
		})
	}
}

// The option is opt-in, so an application that does not set it has to keep the
// id it had before the option existed.
func TestApplicationIDUnsetIsBackwardCompatible(t *testing.T) {
	options := Options{Name: "My App"}

	got, err := applicationID(options)
	if err != nil {
		t.Fatalf("applicationID() returned an unexpected error: %v", err)
	}
	if want := "org.wails." + sanitizeAppName(options.Name); got != want {
		t.Errorf("applicationID() = %q, want the derived id %q", got, want)
	}
}

// An id GTK would reject has to be reported rather than handed to
// gtk_application_new(), and the application still needs an id to start with.
func TestApplicationIDInvalidFallsBackToDerived(t *testing.T) {
	options := Options{
		Name:  "My App",
		Linux: LinuxOptions{ApplicationID: "yeehaw"},
	}

	got, err := applicationID(options)
	if err == nil {
		t.Fatal("applicationID() accepted an id without a '.' separator")
	}
	if want := "org.wails.my_app"; got != want {
		t.Errorf("applicationID() = %q, want the derived id %q", got, want)
	}
}

// The derived id is never routed through validateApplicationID at runtime, so
// make sure sanitizeAppName cannot produce one GTK would refuse.
func TestDerivedApplicationIDIsAlwaysValid(t *testing.T) {
	names := []string{
		"My App",
		"koofr-stage",
		"",
		"1Password",
		"...",
		"app.with.dots",
		"ünïcodé",
		"__leading__and__trailing__",
		"9",
		"-",
		strings.Repeat("long", 200),
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			id, err := applicationID(Options{Name: name})
			if err != nil {
				t.Fatalf("applicationID() returned an unexpected error: %v", err)
			}
			if err := validateApplicationID(id); err != nil {
				t.Errorf("derived id %q is not a valid GTK application id: %v", id, err)
			}
		})
	}
}

// On Wayland the surface app_id comes from g_get_prgname(), so setting only
// ApplicationID has to be enough to have windows match their .desktop file.
func TestProgramName(t *testing.T) {
	tests := []struct {
		name    string
		options Options
		want    string
	}{
		{
			name:    "left alone when neither option is set",
			options: Options{Name: "My App"},
			want:    "",
		},
		{
			name: "inherits the application id",
			options: Options{
				Name:  "My App",
				Linux: LinuxOptions{ApplicationID: "com.example.MyApp"},
			},
			want: "com.example.MyApp",
		},
		{
			name: "an explicit program name wins",
			options: Options{
				Name: "My App",
				Linux: LinuxOptions{
					ApplicationID: "com.example.MyApp",
					ProgramName:   "myapp",
				},
			},
			want: "myapp",
		},
		{
			name: "kept without an application id",
			options: Options{
				Name:  "My App",
				Linux: LinuxOptions{ProgramName: "myapp"},
			},
			want: "myapp",
		},
		{
			// The id GTK was given, not the one it would have rejected.
			name: "inherits the fallback when the id is invalid",
			options: Options{
				Name:  "My App",
				Linux: LinuxOptions{ApplicationID: "yeehaw"},
			},
			want: "org.wails.my_app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appID, _ := applicationID(tt.options)
			if got := programName(tt.options, appID); got != tt.want {
				t.Errorf("programName() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Mirrors the contract of g_application_id_is_valid().
// See: https://docs.gtk.org/gio/type_func.Application.id_is_valid.html
func TestValidateApplicationID(t *testing.T) {
	tests := []struct {
		name  string
		id    string
		valid bool
	}{
		{name: "reverse dns", id: "com.example.MyApp", valid: true},
		{name: "two elements", id: "com.example", valid: true},
		{name: "underscores", id: "com.example.my_app", valid: true},
		{name: "hyphens are discouraged but legal", id: "org.wails.koofr-stage", valid: true},
		{name: "digits inside an element", id: "com.example.App2", valid: true},
		{name: "255 characters", id: "com." + strings.Repeat("a", 251), valid: true},

		{name: "empty", id: "", valid: false},
		{name: "single element", id: "yeehaw", valid: false},
		{name: "leading dot", id: ".com.example", valid: false},
		{name: "trailing dot", id: "com.example.", valid: false},
		{name: "consecutive dots", id: "com..example", valid: false},
		{name: "element starting with a digit", id: "com.example.2ndApp", valid: false},
		{name: "first element starting with a digit", id: "2com.example", valid: false},
		{name: "slash", id: "com.example/MyApp", valid: false},
		{name: "space", id: "com.example.My App", valid: false},
		{name: "colon prefixed unique name", id: ":1.42", valid: false},
		{name: "non ascii", id: "com.example.Mü", valid: false},
		{name: "256 characters", id: "com." + strings.Repeat("a", 252), valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateApplicationID(tt.id)
			if tt.valid && err != nil {
				t.Errorf("validateApplicationID(%q) = %v, want no error", tt.id, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("validateApplicationID(%q) = nil, want an error", tt.id)
			}
		})
	}
}
