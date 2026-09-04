package webcontentsview

import "github.com/leaanthony/u"

// WebPreferences closely mirrors Electron's webPreferences for WebContentsView.
type WebPreferences struct {
	// DevTools enables or disables the developer tools. Default is true.
	DevTools u.Bool
	
	// Javascript enables or disables javascript execution. Default is true.
	Javascript u.Bool
	
	// WebSecurity enables or disables web security (CORS, etc.). Default is true.
	WebSecurity u.Bool
	
	// AllowRunningInsecureContent allows an https page to run http code. Default is false.
	AllowRunningInsecureContent u.Bool
	
	// Images enables or disables image loading. Default is true.
	Images u.Bool
	
	// TextAreasAreResizable controls whether text areas can be resized. Default is true.
	TextAreasAreResizable u.Bool
	
	// WebGL enables or disables WebGL. Default is true.
	WebGL u.Bool
	
	// Plugins enables or disables plugins. Default is false.
	Plugins u.Bool
	
	// ZoomFactor sets the default zoom factor of the page. Default is 1.0.
	ZoomFactor float64
	
	// NavigateOnDragDrop controls whether dropping files triggers navigation. Default is false.
	NavigateOnDragDrop u.Bool
	
	// DefaultFontSize sets the default font size. Default is 16.
	DefaultFontSize int
	
	// DefaultMonospaceFontSize sets the default monospace font size. Default is 13.
	DefaultMonospaceFontSize int
	
	// MinimumFontSize sets the minimum font size. Default is 0.
	MinimumFontSize int
	
	// DefaultEncoding sets the default character encoding. Default is "UTF-8".
	DefaultEncoding string
	// UserAgent sets a custom user agent for the webview.
	UserAgent string

}
