package main

import (
	"io/fs"
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct{}

// FileInfo holds file metadata including permission bits.
type FileInfo struct {
	Name    string
	Mode    os.FileMode
	FsMode  fs.FileMode
	Perm    os.FileMode
	Size    int64
}

func (*Service) GetFileInfo(path string) (_ FileInfo) {
	return
}

func (*Service) GetMode(path string) (_ os.FileMode) {
	return
}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			application.NewService(&Service{}),
		},
	})

	app.Window.New()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
