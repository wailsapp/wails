package pipeline

import "fmt"

func validateWindowsSigningHost(spec SignSpec, host HostCapabilities) error {
	signing := spec.Config
	if signing.Certificate == "" && signing.Thumbprint == "" {
		return hostConfigurationField("certificate", fmt.Errorf("windows signing requires certificate or thumbprint"))
	}
	if signing.Credential != "" && !host.hasCredential(signing.Credential) {
		return hostConfigurationField("credential", fmt.Errorf("Windows signing credential %q is unavailable", signing.Credential))
	}
	if signing.Certificate != "" && (spec.Format == "binary" || spec.Format == "exe" || spec.Format == "dll" || spec.Format == "msix" || spec.Format == "appx" || spec.Format == "nsis") {
		return nil
	}
	if host.os != "windows" {
		return hostConfigurationField("sign", fmt.Errorf("Windows certificate-store/provider signing and MSI signing require a windows host with signtool.exe"))
	}
	return hostConfigurationField("sign", requireHostTools(host, "Windows native signing", "signtool.exe"))
}
