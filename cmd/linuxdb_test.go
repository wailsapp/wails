package cmd

import "testing"

func TestNewLinuxDB(t *testing.T) {
	_ = NewLinuxDB()
}

func TestKnownDistro(t *testing.T) {
	var linuxDB = NewLinuxDB()
	result := linuxDB.GetDistro("ubuntu")
	if result == nil {
		t.Error("Cannot get distro 'ubuntu'")
	}
}

func TestUnknownDistro(t *testing.T) {
	var linuxDB = NewLinuxDB()
	result := linuxDB.GetDistro("unknown")
	if result != nil {
		t.Error("Should get nil for distribution 'unknown'")
	}
}

func TestDefaultRelease(t *testing.T) {
	var linuxDB = NewLinuxDB()
	result := linuxDB.GetDistro("ubuntu")
	if result == nil {
		t.Error("Cannot get distro 'ubuntu'")
	}

	release := result.GetRelease("default")
	if release == nil {
		t.Error("Cannot get release 'default' for distro 'ubuntu'")
	}
}

func TestUnknownRelease(t *testing.T) {
	var linuxDB = NewLinuxDB()
	result := linuxDB.GetDistro("ubuntu")
	if result == nil {
		t.Error("Cannot get distro 'ubuntu'")
	}

	release := result.GetRelease("16.04")
	if release == nil {
		t.Error("Failed to get release 'default' for unknown release version '16.04'")
	}

	if release.Version != "default" {
		t.Errorf("Got version '%s' instead of 'default' for unknown release version '16.04'", result.ID)
	}
}

func TestGetPrerequisites(t *testing.T) {
	var linuxDB = NewLinuxDB()
	result := linuxDB.GetDistro("debian")
	if result == nil {
		t.Error("Cannot get distro 'debian'")
	}

	release := result.GetRelease("default")
	if release == nil {
		t.Error("Failed to get release 'default' for unknown release version '16.04'")
	}

	if release.Version != "default" {
		t.Errorf("Got version '%s' instead of 'default' for unknown release version '16.04'", result.ID)
	}

	if release.Name != "Debian" {
		t.Errorf("Got Release Name '%s' instead of 'debian' for unknown release version '16.04'", release.Name)
	}

	if len(release.Programs) != 3 {
		t.Errorf("Expected %d programs for unknown release version '16.04'", len(release.Programs))
	}
	if len(release.Libraries) != 2 {
		t.Errorf("Expected %d libraries for unknown release version '16.04'", len(release.Libraries))
	}
}
