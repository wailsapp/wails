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
