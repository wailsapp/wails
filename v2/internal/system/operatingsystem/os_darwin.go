package operatingsystem

import "github.com/wailsapp/wails/v2/internal/shell"

func platformInfo() (*OS, error) {
	// Default value
	var result OS
	result.ID = "Unknown"
	result.Name = "MacOS"
	result.Version = "Unknown"

	stdout, stderr, err := shell.RunCommand(".", "sysctl", "kern.osrelease")
	println(stdout)
	println(stderr)
	println(err)
// 		cmd := CreateCommand(directory, command, args...)
// 		var stdo, stde bytes.Buffer
// 		cmd.Stdout = &stdo
// 		cmd.Stderr = &stde
// 		err := cmd.Run()
// 		return stdo.String(), stde.String(), err
// 	}
// 	sysctl := shell.NewCommand("sysctl")
// 	kern.ostype: Darwin
// kern.osrelease: 20.1.0
// kern.osrevision: 199506


	return nil, nil
}