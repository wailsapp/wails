//go:build !darwin

package commands

func macOSDirectoryHasFileProviderMetadata(string) (bool, error) {
	return false, nil
}

func lockMacOSOutput(string) (func(), error) {
	return func() {}, nil
}

func sameMacOSFilesystem(string, string) (bool, error) {
	return true, nil
}
