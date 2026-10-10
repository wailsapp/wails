//go:build linux

package system

import (
	"testing"

	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
)

func TestCheckLocallyInstalled(t *testing.T) {
	found := func() *packagemanager.Dependency {
		return &packagemanager.Dependency{PackageName: "N/A", Installed: true, Version: "11.16.0"}
	}
	notFound := func() *packagemanager.Dependency {
		return &packagemanager.Dependency{PackageName: "N/A"}
	}

	tests := []struct {
		name    string
		dep     packagemanager.Dependency
		checker func() *packagemanager.Dependency
		want    packagemanager.Dependency
	}{
		{
			name:    "unknown to the package manager, found locally",
			dep:     packagemanager.Dependency{Name: "npm"},
			checker: found,
			want:    packagemanager.Dependency{Name: "npm", PackageName: "N/A", Installed: true, Version: "11.16.0"},
		},
		{
			name:    "available from the package manager, found locally",
			dep:     packagemanager.Dependency{Name: "npm", PackageName: "npm", Version: "8.5.1~ds-1"},
			checker: found,
			want:    packagemanager.Dependency{Name: "npm", PackageName: "npm", Installed: true, Version: "11.16.0"},
		},
		{
			name:    "installed by the package manager",
			dep:     packagemanager.Dependency{Name: "npm", PackageName: "npm", Installed: true, Version: "8.5.1~ds-1"},
			checker: found,
			want:    packagemanager.Dependency{Name: "npm", PackageName: "npm", Installed: true, Version: "8.5.1~ds-1"},
		},
		{
			name:    "not found anywhere",
			dep:     packagemanager.Dependency{Name: "npm"},
			checker: notFound,
			want:    packagemanager.Dependency{Name: "npm"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := tt.dep
			checkLocallyInstalled(tt.checker, &dep)
			if dep != tt.want {
				t.Errorf("got %+v, want %+v", dep, tt.want)
			}
		})
	}
}

func TestApplyLocalCheckOnInstalledDependency(t *testing.T) {
	dep := packagemanager.Dependency{Name: "nsis", PackageName: "nsis", Installed: true, Version: "3.08-2"}
	applyLocalCheck(&dep, &packagemanager.Dependency{PackageName: "N/A", Installed: true, Version: "v3.08-2"})

	want := packagemanager.Dependency{Name: "nsis", PackageName: "nsis", Installed: true, Version: "v3.08-2"}
	if dep != want {
		t.Errorf("got %+v, want %+v", dep, want)
	}
}
