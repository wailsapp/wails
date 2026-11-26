module openharmony-hellowails

go 1.26.0

require github.com/wailsapp/wails/v3 v3.1.0

require (
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/term v0.46.0 // indirect
)

replace github.com/wailsapp/wails/v3 => ../../..
