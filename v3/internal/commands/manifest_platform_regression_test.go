package commands

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestManifestRegressionGeneratedMobileSettings(t *testing.T) {
	for _, target := range []string{"android", "ios"} {
		t.Run(target, func(t *testing.T) {
			root := t.TempDir()
			t.Chdir(root)
			hcl := `version = 3
project {
 name = "reviewapp"
 product_name = "Review App"
 identifier = "com.example.review"
 version = "1.0.0"
}
android {
 version_name = "2.4.1"
 version_code = 42
 minimum_sdk = 28
 target_sdk = 35
}
ios {
 background_modes = ["audio", "remote-notification"]
}
`
			require.NoError(t, os.WriteFile(filepath.Join(root, "wails.hcl"), []byte(hcl), 0644))
			loaded, err := manifest.Load(root, "")
			require.NoError(t, err)
			plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{TargetOS: target, TargetArch: "arm64"})
			require.NoError(t, err)
			spec := plan.Nodes[pipeline.NodeKey("target:"+target+"/arm64:assets")].Spec.(pipeline.AssetsSpec)
			handler := &manifestHandler{root: root, config: loaded.Config}
			_, err = handler.assets(spec)
			require.NoError(t, err)
			if target == "android" {
				data, err := os.ReadFile(filepath.Join(root, spec.Directory, "android/app/build.gradle"))
				require.NoError(t, err)

				assert.Contains(t, string(data), "versionCode = 42")
				assert.Contains(t, string(data), `versionName = "2.4.1"`)
				assert.Contains(t, string(data), "minSdk = 28")
				assert.Contains(t, string(data), "targetSdk = 35")
				assert.NotContains(t, string(data), "signingConfig = hasKeystore ? signingConfigs.release : signingConfigs.debug", "production Android release must not fall back to debug signing")
			} else {
				data, err := os.ReadFile(filepath.Join(root, spec.Directory, "ios/xcode/main/Info.plist"))
				require.NoError(t, err)

				assert.Contains(t, string(data), "<key>UIBackgroundModes</key>")
				assert.Contains(t, string(data), "<string>audio</string>")
				assert.Contains(t, string(data), "<string>remote-notification</string>")
			}
		})
	}
}

func TestManifestRegressionDevPreservesProductionBundle(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{TargetOS: "darwin", TargetArch: "arm64", Development: true})
	require.NoError(t, err)
	spec := plan.Nodes["assemble:darwin/arm64"].Spec.(pipeline.PackageSpec)
	t.Logf("Development assembly output: %s", spec.Output)
	sentinel := filepath.Join(root, "bin/app.app/release-sentinel")
	require.NoError(t, os.MkdirAll(filepath.Dir(sentinel), 0755))
	require.NoError(t, os.WriteFile(sentinel, []byte("production"), 0644))
	binary := filepath.Join(root, spec.Binary)
	require.NoError(t, os.MkdirAll(filepath.Dir(binary), 0755))
	require.NoError(t, os.WriteFile(binary, []byte("development executable"), 0755))
	_, err = (&manifestHandler{root: root, config: loaded.Config}).packageApp(spec)
	require.NoError(t, err)
	_, err = os.Stat(sentinel)
	assert.NoError(t, err, "development assembly must preserve production bundle")
	assert.Contains(t, spec.Output, ".wails/dev/")
}

func TestManifestIOSAssetsUseSelectedIconAndPreserveReplacement(t *testing.T) {
	for _, platformIcon := range []bool{false, true} {
		t.Run(map[bool]string{false: "project-icon", true: "platform-icon"}[platformIcon], func(t *testing.T) {
			root := t.TempDir()
			t.Chdir(root)
			require.NoError(t, os.MkdirAll(filepath.Join(root, "build"), 0755))
			writeIcon := func(name string, colour color.RGBA) {
				img := image.NewRGBA(image.Rect(0, 0, 32, 32))
				draw.Draw(img, img.Bounds(), image.NewUniform(colour), image.Point{}, draw.Src)
				var data bytes.Buffer
				require.NoError(t, png.Encode(&data, img))
				require.NoError(t, os.WriteFile(filepath.Join(root, name), data.Bytes(), 0644))
			}
			writeIcon("project.png", color.RGBA{R: 255, A: 255})
			writeIcon("platform.png", color.RGBA{G: 255, A: 255})
			writeIcon("build/appicon.png", color.RGBA{B: 255, A: 255})
			replacement := []byte("<?xml version=\"1.0\"?><plist version=\"1.0\"><dict><key>CFBundleExecutable</key><string>custom</string><key>CFBundleVersion</key><string>99</string></dict></plist>\n")
			require.NoError(t, os.WriteFile(filepath.Join(root, "custom.plist"), replacement, 0644))
			doc := manifest.NewDocument(manifest.Project{Name: "icons", ProductName: "Icons", Identifier: "com.example.icons", Version: "1.0.0", Icon: "project.png"})
			doc.Targets.IOS.InfoPlist = "custom.plist"
			if platformIcon {
				doc.Targets.IOS.Icon = "platform.png"
			}
			require.NoError(t, manifest.WriteMigrationDraft(root, doc))
			loaded, err := manifest.LoadFile(root, filepath.Join(root, manifest.MigratedFilename), "")
			require.NoError(t, err)
			plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{TargetOS: "ios", TargetArch: "arm64"})
			require.NoError(t, err)
			spec := plan.Nodes["target:ios/arm64:assets"].Spec.(pipeline.AssetsSpec)
			_, err = (&manifestHandler{root: root, config: loaded.Config}).assets(spec)
			require.NoError(t, err)
			generated := filepath.Join(root, spec.Directory, "ios", "xcode", "main")
			actual, err := os.ReadFile(filepath.Join(generated, "Info.plist"))
			require.NoError(t, err)
			require.Equal(t, replacement, actual)
			icon, err := os.Open(filepath.Join(generated, "Assets.xcassets", "AppIcon.appiconset", "icon-20.png"))
			require.NoError(t, err)
			defer icon.Close()
			img, err := png.Decode(icon)
			require.NoError(t, err)
			want := color.RGBA{R: 255, A: 255}
			if platformIcon {
				want = color.RGBA{G: 255, A: 255}
			}
			require.Equal(t, want, color.RGBAModel.Convert(img.At(0, 0)))
		})
	}
}
