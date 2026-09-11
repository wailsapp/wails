package commands

import (
	"encoding/xml"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
	"gopkg.in/yaml.v3"
)

func TestMSIXPackageStructure(t *testing.T) {
	for _, config := range []string{
		"info:\n  companyName: Example\n  productName: Example App\n  productIdentifier: com.example.app\n  version: 1.2.3\n  description: Example description\n",
		`{"info":{"companyName":"Example","productName":"Example App","productIdentifier":"com.example.app","version":"1.2.3","description":"Example description"}}`,
	} {
		for _, arch := range []string{"amd64", "arm64"} {
			t.Run(arch+"/"+config[:1], func(t *testing.T) {
				var options MSIXOptions
				if err := yaml.Unmarshal([]byte(config), &options.WailsConfig); err != nil {
					t.Fatal(err)
				}
				options.ExecutablePath = filepath.Join(t.TempDir(), "source.exe")
				options.ExecutableName = "example.exe"
				options.ProcessorArchitecture = arch
				if err := os.WriteFile(options.ExecutablePath, []byte("executable payload"), 0600); err != nil {
					t.Fatal(err)
				}
				if err := validateMSIXOptions(&options); err != nil {
					t.Fatal(err)
				}
				dir := t.TempDir()
				if err := createMSIXPackageStructure(&options, dir); err != nil {
					t.Fatal(err)
				}
				data, err := os.ReadFile(filepath.Join(dir, "AppxManifest.xml"))
				if err != nil {
					t.Fatal(err)
				}
				var manifest struct {
					Identity struct {
						Name      string `xml:"Name,attr"`
						Version   string `xml:"Version,attr"`
						Arch      string `xml:"ProcessorArchitecture,attr"`
						Publisher string `xml:"Publisher,attr"`
					} `xml:"Identity"`
					Application struct {
						Executable string `xml:"Executable,attr"`
					} `xml:"Applications>Application"`
				}
				if err := xml.Unmarshal(data, &manifest); err != nil {
					t.Fatal(err)
				}
				if manifest.Identity.Name != "com.example.app" || manifest.Identity.Version != "1.2.3.0" || manifest.Identity.Arch != archToMSIX(arch) || manifest.Identity.Publisher != "CN=Example" || manifest.Application.Executable != "example.exe" {
					t.Fatalf("incorrect manifest: %+v", manifest)
				}
				payload, err := os.ReadFile(filepath.Join(dir, manifest.Application.Executable))
				if err != nil || string(payload) != "executable payload" {
					t.Fatalf("packaged executable: %q, %v", payload, err)
				}
				for name, size := range map[string][2]int{"Square150x150Logo.png": {150, 150}, "Square44x44Logo.png": {44, 44}, "Wide310x150Logo.png": {310, 150}, "SplashScreen.png": {620, 300}, "StoreLogo.png": {50, 50}} {
					f, err := os.Open(filepath.Join(dir, "Assets", name))
					if err != nil {
						t.Fatal(err)
					}
					img, err := png.Decode(f)
					f.Close()
					if err != nil {
						t.Fatal(err)
					}
					if img.Bounds().Dx() != size[0] || img.Bounds().Dy() != size[1] {
						t.Fatalf("wrong dimensions for %s: %v", name, img.Bounds())
					}
				}
			})
		}
	}
}

func TestFindWindowsSDKTool(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows SDK discovery")
	}
	root := t.TempDir()
	t.Setenv("ProgramFiles(x86)", root)
	t.Setenv("PATH", t.TempDir())
	// The SDK uses a lowercase filename; discovery must be case-insensitive.
	dir := filepath.Join(root, "Windows Kits", "10", "bin", "10.0.26100.0", archToMSIX(runtime.GOARCH))
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "makeappx.exe")
	if err := os.WriteFile(path, nil, 0600); err != nil {
		t.Fatal(err)
	}
	got, err := findWindowsSDKTool("MakeAppx.exe")
	if err != nil || !strings.EqualFold(got, path) {
		t.Fatalf("SDK lookup = %q, %v; want %q", got, err, path)
	}
	pathDir := t.TempDir()
	pathTool := filepath.Join(pathDir, "makeappx.exe")
	if err := os.WriteFile(pathTool, nil, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", pathDir)
	got, err = findWindowsSDKTool("MakeAppx.exe")
	if err != nil || !strings.EqualFold(got, pathTool) {
		t.Fatalf("PATH precedence = %q, %v", got, err)
	}
	if _, err := findWindowsSDKTool("missing-sdk-tool.exe"); err == nil {
		t.Fatal("missing SDK tool should fail")
	}
}

func TestMSIXPackagingToolDoesNotRequireSDK(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows packaging tools")
	}
	dir := t.TempDir()
	t.Setenv("PATH", dir)
	t.Setenv("ProgramFiles(x86)", t.TempDir())
	if err := os.WriteFile(filepath.Join(dir, "MsixPackagingTool.exe"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	if err := checkMSIXTools(&flags.ToolMSIX{UseMsixPackagingTool: true, CertificatePath: "certificate.pfx"}); err != nil {
		t.Fatalf("Packaging Tool should not require the SDK: %v", err)
	}
	if err := checkMSIXTools(&flags.ToolMSIX{UseMakeAppx: true, CertificatePath: "certificate.pfx"}); err == nil {
		t.Fatal("MakeAppx should require the SDK")
	}
}
