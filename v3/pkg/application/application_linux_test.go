//go:build linux && cgo && !android && !server

package application

import "testing"

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
			if got := applicationID(tt.options); got != tt.want {
				t.Errorf("applicationID() = %q, want %q", got, tt.want)
			}
		})
	}
}

// The option is opt-in, so an application that does not set it has to keep the
// id it had before the option existed.
func TestApplicationIDUnsetIsBackwardCompatible(t *testing.T) {
	options := Options{Name: "My App"}

	if got, want := applicationID(options), "org.wails."+sanitizeAppName(options.Name); got != want {
		t.Errorf("applicationID() = %q, want the derived id %q", got, want)
	}
}
