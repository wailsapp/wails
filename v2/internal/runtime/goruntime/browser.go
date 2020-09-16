package goruntime

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Browser defines all browser related operations
type Browser interface {
	Open(url string) error
}

type browser struct{}

// Open a url / file using the system default application
// Credit: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func (b *browser) Open(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func newBrowser() *browser {
	return &browser{}
}
