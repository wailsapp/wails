package wails

import (
	"strings"

	"github.com/dchest/htmlmin"
	// "github.com/leaanthony/mewn"
)

// Config is the configuration structure used when creating a Wails App object
type Config struct {
	Width, Height    int
	Title            string
	defaultHTML      string
	HTML             string
	JS               string
	CSS              string
	Colour           string
	Resizable        bool
	DisableInspector bool
	// isHTMLFragment   bool
}

// GetWidth returns the desired width
func (a *Config) GetWidth() int {
	return a.Width
}

// GetHeight returns the desired height
func (a *Config) GetHeight() int {
	return a.Height
}

// GetTitle returns the desired window title
func (a *Config) GetTitle() string {
	return a.Title
}

// GetDefaultHTML returns the desired window title
func (a *Config) GetDefaultHTML() string {
	return a.defaultHTML
}

// GetResizable returns true if the window should be resizable
func (a *Config) GetResizable() bool {
	return a.Resizable
}

// GetDisableInspector returns true if the inspector should be disabled
func (a *Config) GetDisableInspector() bool {
	return a.DisableInspector
}

// GetColour returns the colour
func (a *Config) GetColour() string {
	return a.Colour
}

// GetCSS returns the user CSS
func (a *Config) GetCSS() string {
	return a.CSS
}

// GetJS returns the user Javascript
func (a *Config) GetJS() string {
	return a.JS
}

func (a *Config) merge(in *Config) error {
	if in.CSS != "" {
		a.CSS = in.CSS
	}
	if in.Title != "" {
		a.Title = in.Title
	}
	// if in.HTML != "" {
	// 	minified, err := htmlmin.Minify([]byte(in.HTML), &htmlmin.Options{
	// 		MinifyScripts: true,
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	inlineHTML := string(minified)
	// 	inlineHTML = strings.Replace(inlineHTML, "'", "\\'", -1)
	// 	inlineHTML = strings.Replace(inlineHTML, "\n", " ", -1)
	// 	a.HTML = strings.TrimSpace(inlineHTML)

	// 	// Deduce whether this is a full html page or a fragment
	// 	// The document is determined to be a fragment if an HTML
	// 	// tag exists and is located before the first div tag
	// 	HTMLTagIndex := strings.Index(a.HTML, "<html")
	// 	DivTagIndex := strings.Index(a.HTML, "<div")

	// 	if HTMLTagIndex == -1 {
	// 		a.isHTMLFragment = true
	// 	} else {
	// 		if DivTagIndex < HTMLTagIndex {
	// 			a.isHTMLFragment = true
	// 		}
	// 	}
	// }

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
func newConfig(userConfig *Config) (*Config, error) {
	result := &Config{
		Width:     800,
		Height:    600,
		Resizable: true,
		Title:     "My Wails App",
		Colour:    "#FFF", // White by default
		// HTML:      mewn.String("./runtime/assets/default.html"),
	}

	if userConfig != nil {
		err := result.merge(userConfig)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
