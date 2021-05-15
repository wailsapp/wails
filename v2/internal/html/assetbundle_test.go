package html

import (
	"testing"
)

func TestNewAssetBundle(t *testing.T) {
	tests := []struct {
		name       string
		pathToHTML string
		wantAssets []string
		wantErr    bool
	}{
		{
			name:       "basic html",
			pathToHTML: "testdata/basic.html",
			wantAssets: []string{
				AssetTypes.HTML,
				AssetTypes.FAVICON,
				AssetTypes.JS,
				AssetTypes.CSS,
			},
			wantErr: false,
		},
		{
			name:       "self closing tags",
			pathToHTML: "testdata/self_closing.html",
			wantAssets: []string{
				AssetTypes.HTML,
				AssetTypes.FAVICON,
				AssetTypes.JS,
				AssetTypes.CSS,
			},
			wantErr: false,
		},
		{
			name:       "multi-line tags",
			pathToHTML: "testdata/self_closing.html",
			wantAssets: []string{
				AssetTypes.HTML,
				AssetTypes.FAVICON,
				AssetTypes.JS,
				AssetTypes.CSS,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAssetBundle(tt.pathToHTML)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAssetBundle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.assets) != len(tt.wantAssets) {
				t.Errorf("NewAssetBundle() len(assets) = %d, want %d",
					len(got.assets), len(tt.wantAssets))
			}

			for i := range tt.wantAssets {
				if i >= len(got.assets) {
					t.Errorf("NewAssetBundle() missing assets[%d].Type = %s",
						i, tt.wantAssets[i])
				} else {
					if got.assets[i].Type != tt.wantAssets[i] {
						t.Errorf("NewAssetBundle() assets[%d].Type = %s, want %s",
							i, got.assets[i].Type, tt.wantAssets[i])
					}
				}
			}
		})
	}
}
