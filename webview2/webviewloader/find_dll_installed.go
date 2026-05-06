//go:build windows && !native_webview2loader

package webviewloader

import (
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const (
	kNumChannels              = 4
	kInstallKeyPath           = "Software\\Microsoft\\EdgeUpdate\\ClientState\\"
	kMinimumCompatibleVersion = "86.0.616.0"
)

var (
	kChannelName = [kNumChannels]string{
		"", "beta", "dev", "canary", // "internal"
	}

	kChannelUuid = [kNumChannels]string{
		"{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}",
		"{2CD8A007-E189-409D-A2C8-9AF4EF3C72AA}",
		"{0D50BFEC-CD6A-4F9A-964C-C7416E3ACB10}",
		"{65C35B14-6C1D-4122-AC46-7148CC9D6497}",
		//"{BE59E8FD-089A-411B-A3B0-051D9E417818}",
	}

	minimumCompatibleVersion, _ = parseVersion(kMinimumCompatibleVersion)
)

func findInstalledClientDll(preferCanary bool) (clientPath string, version *version, err error) {
	for i := 0; i < kNumChannels; i++ {
		channel := i
		if preferCanary {
			channel = (kNumChannels - 1) - i
		}

		key := kInstallKeyPath + kChannelUuid[channel]
		for _, checkSystem := range []bool{true, false} {
			clientPath, version, err := findInstalledClientDllForChannel(key, checkSystem)
			if err == errNoClientDLLFound {
				continue
			}
			if err != nil {
				return "", nil, err
			}

			version.channel = kChannelName[channel]
			return clientPath, version, nil
		}
	}
	return "", nil, errNoClientDLLFound
}

func findInstalledClientDllForChannel(subKey string, system bool) (clientPath string, clientVersion *version, err error) {
	key := registry.LOCAL_MACHINE
	if !system {
		key = registry.CURRENT_USER
	}

	regKey, err := registry.OpenKey(key, subKey, registry.READ|registry.WOW64_32KEY)
	if err != nil {
		return "", nil, mapFindErr(err)
	}
	defer regKey.Close()

	embeddedEdgeSubFolder, _, err := regKey.GetStringValue("EBWebView")
	if err != nil {
		return "", nil, mapFindErr(err)
	}

	if embeddedEdgeSubFolder == "" {
		return "", nil, errNoClientDLLFound
	}

	versionString := filepath.Base(embeddedEdgeSubFolder)
	version, err := parseVersion(versionString)
	if err != nil {
		return "", nil, errNoClientDLLFound
	}

	if version.compare(minimumCompatibleVersion) < 0 {
		return "", nil, errNoClientDLLFound
	}

	dllPath, err := findEmbeddedClientDll(embeddedEdgeSubFolder)
	if err != nil {
		return "", nil, mapFindErr(err)
	}

	return dllPath, &version, nil
}
