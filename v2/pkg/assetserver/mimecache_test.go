package assetserver

import (
	"testing"
)

func TestGetMimetype(t *testing.T) {
	type args struct {
		filename string
		data     []byte
	}
	bomUTF8 := []byte{0xef, 0xbb, 0xbf}
	var emptyMsg []byte
	css := []byte("body{margin:0;padding:0;background-color:#d579b2}#app{font-family:Avenir,Helvetica,Arial,sans-serif;-webkit-font-smoothing:antialiased;-moz-osx-font-smoothing:grayscale;text-align:center;color:#2c3e50;background-color:#ededed}#nav{padding:30px}#nav a{font-weight:700;color:#2c\n3e50}#nav a.router-link-exact-active{color:#42b983}.hello[data-v-4e26ad49]{margin:10px 0}")
	html := []byte("<!DOCTYPE html><html><head>title</head><body></body></html>")
	bomHtml := append(bomUTF8, html...)
	svg := []byte("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"16\" height=\"16\" viewBox=\"0 0 16 16\"><path d=\"M15.707 14.293l-4.822-4.822a6.019 6.019 0 1 0-1.414 1.414l4.822 4.822a1 1 0 0 0 1.414-1.414zM6 10a4 4 0 1 1 4-4 4 4 0 0 1-4 4z\"></path></svg>")
	svgWithComment := append([]byte("<!-- this is a comment -->"), svg...)
	svgWithCommentAndControlChars := append([]byte("    \r\n "), svgWithComment...)
	svgWithBomCommentAndControlChars := append(bomUTF8, append([]byte("    \r\n "), svgWithComment...)...)

	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"nil data", args{"nil.svg", nil}, "image/svg+xml"},
		{"empty data", args{"empty.html", emptyMsg}, "text/html; charset=utf-8"},
		{"css", args{"test.css", css}, "text/css; charset=utf-8"},
		{"js", args{"test.js", []byte("let foo = 'bar'; console.log(foo);")}, "text/javascript; charset=utf-8"},
		{"mjs", args{"test.mjs", []byte("let foo = 'bar'; console.log(foo);")}, "text/javascript; charset=utf-8"},
		{"html-utf8", args{"test_utf8.html", html}, "text/html; charset=utf-8"},
		{"html-bom-utf8", args{"test_bom_utf8.html", bomHtml}, "text/html; charset=utf-8"},
		{"svg", args{"test.svg", svg}, "image/svg+xml"},
		{"svg-w-comment", args{"test_comment.svg", svgWithComment}, "image/svg+xml"},
		{"svg-w-control-comment", args{"test_control_comment.svg", svgWithCommentAndControlChars}, "image/svg+xml"},
		{"svg-w-bom-control-comment", args{"test_bom_control_comment.svg", svgWithBomCommentAndControlChars}, "image/svg+xml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMimetype(tt.args.filename, tt.args.data); got != tt.want {
				t.Errorf("GetMimetype() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}
