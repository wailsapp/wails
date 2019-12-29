package wails

import (
	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails/runtime"
)

// AppConfig is the configuration structure used when creating a Wails App object
type AppConfig struct {
	Width, Height    int
	Title            string
	defaultHTML      string
	HTML             string
	JS               string
	CSS              string
	Colour           string
	Resizable        bool
	DisableInspector bool
}

// GetWidth returns the desired width
func (a *AppConfig) GetWidth() int {
	return a.Width
}

// GetHeight returns the desired height
func (a *AppConfig) GetHeight() int {
	return a.Height
}

// GetTitle returns the desired window title
func (a *AppConfig) GetTitle() string {
	return a.Title
}

// GetDefaultHTML returns the default HTML
func (a *AppConfig) GetDefaultHTML() string {
	return a.defaultHTML
}

// GetResizable returns true if the window should be resizable
func (a *AppConfig) GetResizable() bool {
	return a.Resizable
}

// GetDisableInspector returns true if the inspector should be disabled
func (a *AppConfig) GetDisableInspector() bool {
	return a.DisableInspector
}

// GetColour returns the colour
func (a *AppConfig) GetColour() string {
	return a.Colour
}

// GetCSS returns the user CSS
func (a *AppConfig) GetCSS() string {
	return a.CSS
}

// GetJS returns the user Javascript
func (a *AppConfig) GetJS() string {
	return a.JS
}

func (a *AppConfig) merge(in *AppConfig) error {
	if in.CSS != "" {
		a.CSS = in.CSS
	}
	if in.Title != "" {
		a.Title = runtime.ProcessEncoding(in.Title)
	}

	if in.Colour != "" {
		a.Colour = in.Colour
	}

	if in.JS != "" {
		a.JS = in.JS
	}

	if in.Width != 0 {
		a.Width = in.Width
	}
	if in.Height != 0 {
		a.Height = in.Height
	}
	a.Resizable = in.Resizable
	a.DisableInspector = in.DisableInspector

	return nil
}

// Creates the default configuration
func newConfig(userConfig *AppConfig) (*AppConfig, error) {
	result := &AppConfig{
		Width:     800,
		Height:    600,
		Resizable: true,
		Title:     "My Wails App",
		Colour:    "#FFF", // White by default
		HTML:      mewn.String("./runtime/assets/default.html"),
	}

	if userConfig != nil {
		err := result.merge(userConfig)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
