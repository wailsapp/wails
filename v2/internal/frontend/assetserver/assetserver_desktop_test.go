// +build desktop

package assetserver

import (
	"embed"
	"github.com/matryer/is"
	"github.com/wailsapp/wails/v2/internal/frontend/assetserver/testdata"
	"strconv"
	"testing"
)

var runtimeInjection = `<script defer src="/wails/runtime.js"></script>`
var expected = `<html><head><link rel="stylesheet" href="/main.css"></head><body data-wails-drag><div id="logo"></div>` + runtimeInjection + `</body></html>`

//go:embed testdata/subdir
var subdir embed.FS

//go:embed testdata
var multiple embed.FS

func TestAssetServer_Init(t *testing.T) {

	is2 := is.New(t)

	tests := []struct {
		assets  embed.FS
		want    string
		wantErr bool
	}{
		{testdata.TopLevelFS, expected, false},
		{subdir, expected, false},
		{multiple, expected, true},
	}

	for idx, tt := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			server, err := NewAssetServer(tt.assets)
			if tt.wantErr {
				is2.True(err != nil)
			} else {
				is2.NoErr(err)
				is2.Equal(string(server.indexFile), tt.want)
			}
		})
	}

}
