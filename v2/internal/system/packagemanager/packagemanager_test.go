//go:build linux
// +build linux

package packagemanager

import "testing"

type fakePackageManager struct {
	available map[string]bool
	installed map[string]bool
}

func (f *fakePackageManager) Name() string { return "fake" }

func (f *fakePackageManager) Packages() packagemap {
	return packagemap{
		"libwebkit": []*Package{
			{Name: "webkit2gtk", SystemPackage: true, Library: true},
			{Name: "webkit2gtk-4.1", SystemPackage: true, Library: true, BuildTags: "webkit2_41"},
		},
	}
}

func (f *fakePackageManager) PackageInstalled(pkg *Package) (bool, error) {
	return f.installed[pkg.Name], nil
}

func (f *fakePackageManager) PackageAvailable(pkg *Package) (bool, error) {
	return f.available[pkg.Name], nil
}

func (f *fakePackageManager) InstallCommand(pkg *Package) string {
	return "install " + pkg.Name
}

func TestDependenciesPrefersInstalledAlternative(t *testing.T) {
	tests := []struct {
		name          string
		available     map[string]bool
		installed     map[string]bool
		wantPackage   string
		wantInstalled bool
		wantBuildTags string
	}{
		{
			name:          "both available, only second installed",
			available:     map[string]bool{"webkit2gtk": true, "webkit2gtk-4.1": true},
			installed:     map[string]bool{"webkit2gtk-4.1": true},
			wantPackage:   "webkit2gtk-4.1",
			wantInstalled: true,
			wantBuildTags: "webkit2_41",
		},
		{
			name:        "both available, neither installed",
			available:   map[string]bool{"webkit2gtk": true, "webkit2gtk-4.1": true},
			wantPackage: "webkit2gtk",
		},
		{
			name:          "first unavailable, second installed",
			available:     map[string]bool{"webkit2gtk-4.1": true},
			installed:     map[string]bool{"webkit2gtk-4.1": true},
			wantPackage:   "webkit2gtk-4.1",
			wantInstalled: true,
			wantBuildTags: "webkit2_41",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps, err := Dependencies(&fakePackageManager{available: tt.available, installed: tt.installed})
			if err != nil {
				t.Fatal(err)
			}
			if len(deps) != 1 {
				t.Fatalf("expected 1 dependency, got %d", len(deps))
			}
			got := deps[0]
			if got.PackageName != tt.wantPackage {
				t.Errorf("PackageName: expected %q, got %q", tt.wantPackage, got.PackageName)
			}
			if got.Installed != tt.wantInstalled {
				t.Errorf("Installed: expected %v, got %v", tt.wantInstalled, got.Installed)
			}
			if got.BuildTags != tt.wantBuildTags {
				t.Errorf("BuildTags: expected %q, got %q", tt.wantBuildTags, got.BuildTags)
			}
			if got.InstallCommand != "install "+tt.wantPackage {
				t.Errorf("InstallCommand: expected %q, got %q", "install "+tt.wantPackage, got.InstallCommand)
			}
		})
	}
}
