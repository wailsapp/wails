//go:build !darwin

package commands

func lockMacOSOutput(string) (func(), error) {
	return func() {}, nil
}

func sameMacOSFilesystem(string, string) (bool, error) {
	return true, nil
}
