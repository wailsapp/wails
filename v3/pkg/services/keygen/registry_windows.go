//go:build windows

package keygen

// SetRegistryValue stores a value in the Windows Registry
// This method is only available on Windows
func (s *Service) SetRegistryValue(name, value string) error {
	if platform, ok := s.platform.(*platformKeygenWindows); ok {
		return platform.SetRegistryValue(name, value)
	}
	return nil
}

// GetRegistryValue retrieves a value from the Windows Registry
// This method is only available on Windows
func (s *Service) GetRegistryValue(name string) (string, error) {
	if platform, ok := s.platform.(*platformKeygenWindows); ok {
		return platform.GetRegistryValue(name)
	}
	return "", nil
}

// DeleteRegistryValue removes a value from the Windows Registry
// This method is only available on Windows
func (s *Service) DeleteRegistryValue(name string) error {
	if platform, ok := s.platform.(*platformKeygenWindows); ok {
		return platform.DeleteRegistryValue(name)
	}
	return nil
}
