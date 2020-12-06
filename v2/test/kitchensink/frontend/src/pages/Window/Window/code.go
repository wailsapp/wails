package main

import (
	"image"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/internal/runtime"
)

type MyStruct struct {
	runtime *wails.Runtime
	image   *runtime.Store
}

func (n *Notepad) WailsInit(runtime *wails.Runtime) error {
	n.runtime = runtime
	n.image = runtime.Store.New("mainimage")
	return nil
}

func (n *MyStruct) LoadImage(filename string) error {

	// Load filedata
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	img, fmtName, err := image.DecodeConfig(f)
	if err != nil {
		return err
	}

	// Sync the image data with the frontend
	n.image.Set(img)

	// Get the size of the image
	n.runtime.Window.SetSize(img.Width, img.Height)

	// Place window in center
	n.runtime.Window.Center()
}
