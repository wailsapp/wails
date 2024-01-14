//go:build !production

package assetserver

import (
	"os"
)

func GetDevServerURL() string {
	return os.Getenv("FRONTEND_DEVSERVER_URL")
}
