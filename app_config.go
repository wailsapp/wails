package wails

import (
	"strings"

	"github.com/dchest/htmlmin"
	"github.com/leaanthony/mewn"
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
	isHTMLFragment   bool
}

func (a *AppConfig) merge(in *AppConfig) error {
	if in.CSS != "" {
		a.CSS = in.CSS
	}
	if in.Title != "" {
		a.Title = in.Title
	}
	if in.HTML != "" {
		minified, err := htmlmin.Minify([]byte(in.HTML), &htmlmin.Options{
			MinifyScripts: true,
		})
		if err != nil {
			return err
		}
		inlineHTML := string(minified)
		inlineHTML = strings.Replace(inlineHTML, "'", "\\'", -1)
		inlineHTML = strings.Replace(inlineHTML, "\n", " ", -1)
		a.HTML = strings.TrimSpace(inlineHTML)

		// Deduce whether this is a full html page or a fragment
		// The document is determined to be a fragment if an HMTL
		// tag exists and is located before the first div tag
		HTMLTagIndex := strings.Index(a.HTML, "<html")
		DivTagIndex := strings.Index(a.HTML, "<div")

		if HTMLTagIndex == -1 {
			a.isHTMLFragment = true
		} else {
			if DivTagIndex < HTMLTagIndex {
				a.isHTMLFragment = true
			}
		}
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
func newAppConfig(userConfig *AppConfig) (*AppConfig, error) {
	result := &AppConfig{
		Width:     800,
		Height:    600,
		Resizable: true,
		Title:     "My Wails App",
		Colour:    "#FFF", // White by default
		HTML:      mewn.String("./wailsruntimeassets/default/default.html"),
	}

	if userConfig != nil {
		err := result.merge(userConfig)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
