package commands

import (
	"fmt"
	"os"

	"github.com/wailsapp/wails/v3/internal/keychain"
	"github.com/wailsapp/wails/v3/internal/signingconfig"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func manifestSigningPassword(config manifest.SigningPlatform, platform string, environment []string) (string, error) {
	if config.Credential != "" {
		password := environmentValues(mergeEnvironment(os.Environ(), environment))[config.Credential]
		if password == "" {
			return "", fmt.Errorf("%s signing password environment variable %s is empty", platform, config.Credential)
		}
		return password, nil
	}
	if !signingconfig.IsGlobalPath(config, "certificate", config.Certificate) {
		return "", nil
	}
	key := keychain.KeyWindowsCertPassword
	if platform == "linux" {
		key = keychain.KeyPGPPassword
	}
	password, _ := keychain.Get(key)
	return password, nil
}
