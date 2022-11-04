package notification

import (
	"os"
	"path"
	"sync"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/internal/fs"
)

var (
	notificatonUUID     string
	notificatonUUIDOnce sync.Once
	tmpDir              string
)

func notificationOnceFunc() {
	notificatonUUID = uuid.New().String()

}

func checkTmpDir() error {
	tmpDir = path.Join(os.TempDir(), "wails-notifications")
	if !fs.DirExists(tmpDir) {
		if err := fs.Mkdir(tmpDir); err != nil {
			return err
		}
	}
	return nil
}

func AppIconPath(data []byte) (string, error) {
	notificatonUUIDOnce.Do(notificationOnceFunc)

	if err := checkTmpDir(); err != nil {
		return "", err
	}

	tmpFilePath := path.Join(tmpDir, notificatonUUID+"-icon")
	return tmpFilePath, os.WriteFile(tmpFilePath, data, 0755)
}

func SoundPath(data []byte) (string, error) {
	notificatonUUIDOnce.Do(notificationOnceFunc)

	if err := checkTmpDir(); err != nil {
		return "", err
	}

	tmpFilePath := path.Join(tmpDir, notificatonUUID+"-sound")
	return tmpFilePath, os.WriteFile(tmpFilePath, data, 0755)
}

func ContentImagePath(data []byte) (string, error) {
	notificatonUUIDOnce.Do(notificationOnceFunc)

	if err := checkTmpDir(); err != nil {
		return "", err
	}

	tmpFilePath := path.Join(tmpDir, notificatonUUID+"-content-image")
	return tmpFilePath, os.WriteFile(tmpFilePath, data, 0755)
}
