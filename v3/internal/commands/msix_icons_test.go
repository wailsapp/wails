package commands

import (
	"github.com/stretchr/testify/require"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestMSIXAssetsUseApplicationIconAtNativeDimensions(t *testing.T) {
	root := t.TempDir()
	icon := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			icon.Set(x, y, color.RGBA{R: 231, G: 27, B: 59, A: 255})
		}
	}
	path := filepath.Join(root, "icon.png")
	f, err := os.Create(path)
	require.NoError(t, err)
	require.NoError(t, png.Encode(f, icon))
	require.NoError(t, f.Close())
	exe := filepath.Join(root, "test.exe")
	require.NoError(t, os.WriteFile(exe, []byte("test"), 0644))
	manifest := filepath.Join(root, "manifest.xml")
	require.NoError(t, os.WriteFile(manifest, []byte("<Package/>"), 0644))
	output := filepath.Join(root, "out")
	require.NoError(t, createMSIXPackageStructure(&MSIXOptions{ExecutablePath: exe, AppxManifest: manifest, IconPath: path}, output))
	for name, size := range map[string]image.Point{"Square44x44Logo.png": {44, 44}, "Square150x150Logo.png": {150, 150}, "Wide310x150Logo.png": {310, 150}, "SplashScreen.png": {620, 300}, "StoreLogo.png": {50, 50}} {
		f, err := os.Open(filepath.Join(output, "Assets", name))
		require.NoError(t, err)
		got, err := png.Decode(f)
		require.NoError(t, err)
		require.NoError(t, f.Close())
		require.Equal(t, size, got.Bounds().Size())
		require.Equal(t, color.RGBA{R: 231, G: 27, B: 59, A: 255}, color.RGBAModel.Convert(got.At(size.X/2, size.Y/2)))
	}
}
