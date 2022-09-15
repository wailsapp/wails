package gomod

import (
	"reflect"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/matryer/is"
)

const basic string = `module changeme

go 1.17

require github.com/wailsapp/wails/v2 v2.0.0-beta.7

require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fasthttp/websocket v0.0.0-20200320073529-1554a54587ab // indirect
	github.com/wailsapp/mimetype v1.4.1-beta.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gofiber/fiber/v2 v2.17.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2 v0.0.0-20211007092718-65d2f028ef2d // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leaanthony/typescriptify-golang-structs v0.1.7 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc v0.0.0-20210921073452-54963136bf18 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20200117113501-90175b0fbe3f // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

//replace github.com/wailsapp/wails/v2 v2.0.0-beta.7 => C:\Users\leaan\Documents\wails-v2-beta\wails\v2
`

func TestGetWailsVersion(t *testing.T) {
	tests := []struct {
		name      string
		goModText []byte
		want      *semver.Version
		wantErr   bool
	}{
		{"basic", []byte(basic), semver.MustParse("v2.0.0-beta.7"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetWailsVersionFromModFile(tt.goModText)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWailsVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWailsVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

const basicUpdated string = `module changeme

go 1.17

require github.com/wailsapp/wails/v2 v2.0.0-beta.20

require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fasthttp/websocket v0.0.0-20200320073529-1554a54587ab // indirect
	github.com/wailsapp/mimetype v1.3.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gofiber/fiber/v2 v2.17.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/wailsapp/wails/v2/internal/go-common-file-dialog v1.0.3 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2 v0.0.0-20211007092718-65d2f028ef2d // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leaanthony/typescriptify-golang-structs v0.1.7 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc v0.0.0-20210921073452-54963136bf18 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20200117113501-90175b0fbe3f // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

//replace github.com/wailsapp/wails/v2 v2.0.0-beta.7 => C:\Users\leaan\Documents\wails-v2-beta\wails\v2
`

const multilineRequire = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.7
)
require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fasthttp/websocket v0.0.0-20200320073529-1554a54587ab // indirect
	github.com/wailsapp/mimetype v1.3.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gofiber/fiber/v2 v2.17.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/wailsapp/wails/v2/internal/go-common-file-dialog v1.0.3 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2 v0.0.0-20211007092718-65d2f028ef2d // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leaanthony/typescriptify-golang-structs v0.1.7 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc v0.0.0-20210921073452-54963136bf18 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20200117113501-90175b0fbe3f // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

//replace github.com/wailsapp/wails/v2 v2.0.0-beta.7 => C:\Users\leaan\Documents\wails-v2-beta\wails\v2
`
const multilineReplace = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.7
)
require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fasthttp/websocket v0.0.0-20200320073529-1554a54587ab // indirect
	github.com/wailsapp/mimetype v1.3.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gofiber/fiber/v2 v2.17.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/wailsapp/wails/v2/internal/go-common-file-dialog v1.0.3 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2 v0.0.0-20211007092718-65d2f028ef2d // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leaanthony/typescriptify-golang-structs v0.1.7 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc v0.0.0-20210921073452-54963136bf18 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20200117113501-90175b0fbe3f // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

replace github.com/wailsapp/wails/v2 v2.0.0-beta.7 => C:\Users\leaan\Documents\wails-v2-beta\wails\v2
`

const multilineReplaceNoVersion = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.7
)
require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fasthttp/websocket v0.0.0-20200320073529-1554a54587ab // indirect
	github.com/wailsapp/mimetype v1.3.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gofiber/fiber/v2 v2.17.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/wailsapp/wails/v2/internal/go-common-file-dialog v1.0.3 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2 v0.0.0-20211007092718-65d2f028ef2d // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leaanthony/typescriptify-golang-structs v0.1.7 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc v0.0.0-20210921073452-54963136bf18 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20200117113501-90175b0fbe3f // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

replace github.com/wailsapp/wails/v2 => C:\Users\leaan\Documents\wails-v2-beta\wails\v2
`

const multilineReplaceNoVersionBlock = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.7
)
require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fasthttp/websocket v0.0.0-20200320073529-1554a54587ab // indirect
	github.com/wailsapp/mimetype v1.3.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gofiber/fiber/v2 v2.17.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/wailsapp/wails/v2/internal/go-common-file-dialog v1.0.3 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2 v0.0.0-20211007092718-65d2f028ef2d // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leaanthony/typescriptify-golang-structs v0.1.7 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc v0.0.0-20210921073452-54963136bf18 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20200117113501-90175b0fbe3f // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

replace (
	github.com/wailsapp/wails/v2 => C:\Users\leaan\Documents\wails-v2-beta\wails\v2
)
`

const multilineReplaceBlock = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.7
)
require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fasthttp/websocket v0.0.0-20200320073529-1554a54587ab // indirect
	github.com/wailsapp/mimetype v1.3.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gofiber/fiber/v2 v2.17.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/wailsapp/wails/v2/internal/go-common-file-dialog v1.0.3 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2 v0.0.0-20211007092718-65d2f028ef2d // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leaanthony/typescriptify-golang-structs v0.1.7 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc v0.0.0-20210921073452-54963136bf18 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20200117113501-90175b0fbe3f // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

replace (
	github.com/wailsapp/wails/v2 v2.0.0-beta.7 => C:\Users\leaan\Documents\wails-v2-beta\wails\v2
)
`

const multilineRequireUpdated = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.20
)

require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fasthttp/websocket v0.0.0-20200320073529-1554a54587ab // indirect
	github.com/wailsapp/mimetype v1.3.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gofiber/fiber/v2 v2.17.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/wailsapp/wails/v2/internal/go-common-file-dialog v1.0.3 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2 v0.0.0-20211007092718-65d2f028ef2d // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leaanthony/typescriptify-golang-structs v0.1.7 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc v0.0.0-20210921073452-54963136bf18 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20200117113501-90175b0fbe3f // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

//replace github.com/wailsapp/wails/v2 v2.0.0-beta.7 => C:\Users\leaan\Documents\wails-v2-beta\wails\v2
`

func TestUpdateGoModVersion(t *testing.T) {
	is2 := is.New(t)

	type args struct {
		goModText      []byte
		currentVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"basic", args{[]byte(basic), "v2.0.0-beta.20"}, []byte(basicUpdated), false},
		{"basicmultiline", args{[]byte(multilineRequire), "v2.0.0-beta.20"}, []byte(multilineRequireUpdated), false},
		{"basicmultilinereplace", args{[]byte(multilineReplace), "v2.0.0-beta.20"}, []byte(multilineReplaceUpdated), false},
		{"basicmultilinereplaceblock", args{[]byte(multilineReplaceBlock), "v2.0.0-beta.20"}, []byte(multilineReplaceBlockUpdated), false},
		{"basicmultilinereplacenoversion", args{[]byte(multilineReplaceNoVersion), "v2.0.0-beta.20"}, []byte(multilineReplaceNoVersionUpdated), false},
		{"basicmultilinereplacenoversionblock", args{[]byte(multilineReplaceNoVersionBlock), "v2.0.0-beta.20"}, []byte(multilineReplaceNoVersionBlockUpdated), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateGoModVersion(tt.args.goModText, tt.args.currentVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateGoModVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			is2.Equal(got, tt.want)
		})
	}
}

func TestGoModOutOfSync(t *testing.T) {
	is2 := is.New(t)

	type args struct {
		goModData      []byte
		currentVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"basic", args{[]byte(basic), "v2.0.0-beta.20"}, true, false},
		{"basicmultiline", args{[]byte(multilineRequire), "v2.0.0-beta.20"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoModOutOfSync(tt.args.goModData, tt.args.currentVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoModOutOfSync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			is2.Equal(got, tt.want)
		})
	}
}

const multilineReplaceUpdated = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.20
)

require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fasthttp/websocket v0.0.0-20200320073529-1554a54587ab // indirect
	github.com/wailsapp/mimetype v1.3.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gofiber/fiber/v2 v2.17.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/wailsapp/wails/v2/internal/go-common-file-dialog v1.0.3 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2 v0.0.0-20211007092718-65d2f028ef2d // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leaanthony/typescriptify-golang-structs v0.1.7 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc v0.0.0-20210921073452-54963136bf18 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20200117113501-90175b0fbe3f // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

replace github.com/wailsapp/wails/v2 v2.0.0-beta.20 => C:\Users\leaan\Documents\wails-v2-beta\wails\v2
`
const multilineReplaceNoVersionUpdated = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.20
)

require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fasthttp/websocket v0.0.0-20200320073529-1554a54587ab // indirect
	github.com/wailsapp/mimetype v1.3.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gofiber/fiber/v2 v2.17.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/wailsapp/wails/v2/internal/go-common-file-dialog v1.0.3 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2 v0.0.0-20211007092718-65d2f028ef2d // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leaanthony/typescriptify-golang-structs v0.1.7 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc v0.0.0-20210921073452-54963136bf18 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20200117113501-90175b0fbe3f // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

replace github.com/wailsapp/wails/v2 => C:\Users\leaan\Documents\wails-v2-beta\wails\v2
`
const multilineReplaceNoVersionBlockUpdated = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.20
)

require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fasthttp/websocket v0.0.0-20200320073529-1554a54587ab // indirect
	github.com/wailsapp/mimetype v1.3.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gofiber/fiber/v2 v2.17.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/wailsapp/wails/v2/internal/go-common-file-dialog v1.0.3 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2 v0.0.0-20211007092718-65d2f028ef2d // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leaanthony/typescriptify-golang-structs v0.1.7 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc v0.0.0-20210921073452-54963136bf18 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20200117113501-90175b0fbe3f // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

replace (
	github.com/wailsapp/wails/v2 => C:\Users\leaan\Documents\wails-v2-beta\wails\v2
)
`

const multilineReplaceBlockUpdated = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.20
)

require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fasthttp/websocket v0.0.0-20200320073529-1554a54587ab // indirect
	github.com/wailsapp/mimetype v1.3.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gofiber/fiber/v2 v2.17.0 // indirect
	github.com/gofiber/websocket/v2 v2.0.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/wailsapp/wails/v2/internal/go-common-file-dialog v1.0.3 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2 v0.0.0-20211007092718-65d2f028ef2d // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leaanthony/typescriptify-golang-structs v0.1.7 // indirect
	github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc v0.0.0-20210921073452-54963136bf18 // indirect
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20200117113501-90175b0fbe3f // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.28.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

replace (
	github.com/wailsapp/wails/v2 v2.0.0-beta.20 => C:\Users\leaan\Documents\wails-v2-beta\wails\v2
)
`

const basicGo117 string = `module changeme

go 1.17

require github.com/wailsapp/wails/v2 v2.0.0-beta.7

`

const basicGo118 string = `module changeme

go 1.18

require github.com/wailsapp/wails/v2 v2.0.0-beta.7
`

const basicGo119 string = `module changeme

go 1.19

require github.com/wailsapp/wails/v2 v2.0.0-beta.7
`

func TestUpdateGoModGoVersion(t *testing.T) {
	is2 := is.New(t)

	type args struct {
		goModText      []byte
		currentVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		updated bool
	}{
		{"basic1.17", args{[]byte(basicGo117), "1.18"}, []byte(basicGo118), true},
		{"basic1.18", args{[]byte(basicGo118), "1.18"}, []byte(basicGo118), false},
		{"basic1.19", args{[]byte(basicGo119), "1.17"}, []byte(basicGo119), false},
		{"basic1.19", args{[]byte(basicGo119), "1.18"}, []byte(basicGo119), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, updated, err := SyncGoVersion(tt.args.goModText, tt.args.currentVersion)
			if err != nil {
				t.Errorf("UpdateGoModVersion() error = %v", err)
				return
			}
			if updated != tt.updated {
				t.Errorf("UpdateGoModVersion() updated = %t, want = %t", updated, tt.updated)
				return
			}
			is2.Equal(got, tt.want)
		})
	}
}
