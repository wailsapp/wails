//go:build !production

package assetserver

import (
	"os"
)

func GetDevServerURL() string {
	return os.Getenv("WAILS_DEVSERVER_URL")
}
