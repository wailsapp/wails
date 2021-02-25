package runtime

import (
	"bytes"
	"runtime"

	"github.com/abadojack/whatlanggo"
	"github.com/wailsapp/wails/lib/interfaces"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func detectEncoding(text string) (encoding.Encoding, string) {
	// korean
	var enc encoding.Encoding
	info := whatlanggo.Detect(text)
	//fmt.Println("Language:", info.Lang.String(), " Script:", whatlanggo.Scripts[info.Script], " Confidence: ", info.Confidence)
	switch info.Lang.String() {
	case "Korean":
		enc = korean.EUCKR
	case "Mandarin":
		enc = simplifiedchinese.GBK
	case "Japanese":
		enc = japanese.EUCJP
	}
	return enc, info.Lang.String()
}

// ProcessEncoding attempts to convert CKJ strings to UTF-8
func ProcessEncoding(text string) string {
	if runtime.GOOS != "windows" {
		return text
	}

	encoding, _ := detectEncoding(text)
	if encoding != nil {
		var bufs bytes.Buffer
		wr := transform.NewWriter(&bufs, encoding.NewEncoder())
		_, err := wr.Write([]byte(text))
		defer wr.Close()
		if err != nil {
			return ""
		}

		return bufs.String()
	}
	return text
}

// Window exposes an interface for manipulating the window
type Window struct {
	renderer interfaces.Renderer
}

// NewWindow creates a new Window struct
func NewWindow(renderer interfaces.Renderer) *Window {
	return &Window{
		renderer: renderer,
	}
}

// SetColour sets the the window colour
func (r *Window) SetColour(colour string) error {
	return r.renderer.SetColour(colour)
}

// SetMinSize sets the minimum size of a resizable window
func (r *Window) SetMinSize(width int, height int) {
	r.renderer.SetMinSize(width, height)
}

// SetMaxSize sets the maximum size of a resizable window
func (r *Window) SetMaxSize(width int, height int) {
	r.renderer.SetMaxSize(width, height)
}


// Fullscreen makes the window fullscreen
func (r *Window) Fullscreen() {
	r.renderer.Fullscreen()
}

// UnFullscreen attempts to restore the window to the size/position before fullscreen
func (r *Window) UnFullscreen() {
	r.renderer.UnFullscreen()
}

// SetTitle sets the the window title
func (r *Window) SetTitle(title string) {
	title = ProcessEncoding(title)
	r.renderer.SetTitle(title)
}

// Close shuts down the window and therefore the app
func (r *Window) Close() {
	r.renderer.Close()
}
