package main

import (
	_ "embed"
	"log"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService struct{}

// FileInfo holds file metadata with stdlib type fields.
// This tests that os.FileMode and other stdlib named types
// are rendered as their underlying primitive types (number, string, etc.)
// rather than causing a crash or emitting invalid TypeScript.
type FileInfo struct {
	Name    string
	Mode    os.FileMode
	ModTime time.Time
	Size    int64
	IsDir   bool
}

// GetFileInfo returns file info for path.
func (*GreetService) GetFileInfo(path string) FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		return FileInfo{}
	}
	return FileInfo{
		Name:    info.Name(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		Size:    info.Size(),
		IsDir:   info.IsDir(),
	}
}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
	})

	app.Window.New()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
