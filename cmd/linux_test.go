package cmd

import "testing"

func TestUbuntuDetection(t *testing.T) {
	osrelease := `
NAME="Ubuntu"
VERSION="18.04.2 LTS (Bionic Beaver)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 18.04.2 LTS"
VERSION_ID="18.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
VERSION_CODENAME=bionic
UBUNTU_CODENAME=bionic
`

	result := parseOsRelease(osrelease)
	if result.Distribution != Ubuntu {
		t.Errorf("expected 'Ubuntu' ID but got '%d'", result.Distribution)
	}
}

func TestTumbleweedDetection(t *testing.T) {
	osrelease := `
NAME="openSUSE Tumbleweed"
# VERSION="20200414"
ID="opensuse-tumbleweed"
ID_LIKE="opensuse suse"
VERSION_ID="20200414"
PRETTY_NAME="openSUSE Tumbleweed"
ANSI_COLOR="0;32"
CPE_NAME="cpe:/o:opensuse:tumbleweed:20200414"
BUG_REPORT_URL="https://bugs.opensuse.org"
HOME_URL="https://www.opensuse.org/"
LOGO="distributor-logo"
`

	result := parseOsRelease(osrelease)
	if result.Distribution != Tumbleweed {
		t.Errorf("expected 'Tumbleweed' ID but got '%d'", result.Distribution)
	}
}
