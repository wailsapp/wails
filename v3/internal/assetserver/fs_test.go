package assetserver

import (
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestAssetFileServerIgnoresIndexLookalikes(t *testing.T) {
	tests := []struct {
		name   string
		assets fstest.MapFS
	}{
		{"suffix filename", fstest.MapFS{
			"not-index.html":         {},
			"dist/assets/index.html": {Data: []byte("real index")},
		}},
		{"root directory", fstest.MapFS{
			"index.html/child.txt": {},
			"site/index.html":      {Data: []byte("real index")},
		}},
		{"nested directory", fstest.MapFS{
			"x/index.html/child.txt": {},
			"site/index.html":        {Data: []byte("real index")},
		}},
		{"real root index", fstest.MapFS{
			"index.html":      {Data: []byte("real index")},
			"not-index.html":  {},
			"site/index.html": {Data: []byte("other index")},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			newAssetFileServerFS(tt.assets).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))
			if recorder.Code != http.StatusOK || recorder.Body.String() != "real index" {
				t.Errorf("response = %d %q, want 200 and real index", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestFindPathToFileRejectsIndexLookalikes(t *testing.T) {
	_, err := findPathToFile(fstest.MapFS{
		"not-index.html":       {},
		"index.html/child.txt": {},
	}, indexHTML)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("error = %v, want fs.ErrNotExist", err)
	}
}
