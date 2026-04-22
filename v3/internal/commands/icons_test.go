package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGenerateIcon(t *testing.T) {
	tests := []struct {
		name             string
		setup            func(t *testing.T) *IconsOptions
		wantErr          bool
		wantErrContains  string
		requireDarwin    bool
		requireNonDarwin bool
		test             func(t *testing.T, options *IconsOptions) error
	}{
		{
			name: "should generate an icon when using the `example` flag",
			setup: func(t *testing.T) *IconsOptions {
				return &IconsOptions{
					Example: true,
				}
			},
			wantErr: false,
			test: func(t *testing.T, options *IconsOptions) error {
				// the file `appicon.png` should be created in the current directory
				// check for the existence of the file
				f, err := os.Stat("appicon.png")
				if err != nil {
					return err
				}
				defer func() {
					err := os.Remove("appicon.png")
					if err != nil {
						panic(err)
					}
				}()
				if f.IsDir() {
					return fmt.Errorf("appicon.png is a directory")
				}
				if f.Size() == 0 {
					return fmt.Errorf("appicon.png is empty")
				}
				return nil
			},
		},
		{
			name: "should generate a .ico file when using the `input` flag and `windowsfilename` flag",
			setup: func(t *testing.T) *IconsOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				// Get the path to the example icon
				exampleIcon := filepath.Join(localDir, "build_assets", "appicon.png")
				return &IconsOptions{
					Input:           exampleIcon,
					WindowsFilename: "appicon.ico",
				}
			},
			wantErr: false,
			test: func(t *testing.T, options *IconsOptions) error {
				// the file `appicon.ico` should be created in the current directory
				// check for the existence of the file
				f, err := os.Stat("appicon.ico")
				if err != nil {
					return err
				}
				defer func() {
					// Remove the file
					_ = os.Remove("appicon.ico")
				}()
				if f.IsDir() {
					return fmt.Errorf("appicon.ico is a directory")
				}
				if f.Size() == 0 {
					return fmt.Errorf("appicon.ico is empty")
				}
				// Remove the file
				err = os.Remove("appicon.ico")
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "should generate a .icns file when using the `input` flag and `macfilename` flag",
			setup: func(t *testing.T) *IconsOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				// Get the path to the example icon
				exampleIcon := filepath.Join(localDir, "build_assets", "appicon.png")
				return &IconsOptions{
					Input:       exampleIcon,
					MacFilename: "appicon.icns",
				}
			},
			wantErr: false,
			test: func(t *testing.T, options *IconsOptions) error {
				// the file `appicon.icns` should be created in the current directory
				// check for the existence of the file
				f, err := os.Stat("appicon.icns")
				if err != nil {
					return err
				}
				defer func() {
					// Remove the file
					err = os.Remove("appicon.icns")
					if err != nil {
						panic(err)
					}
				}()
				if f.IsDir() {
					return fmt.Errorf("appicon.icns is a directory")
				}
				if f.Size() == 0 {
					return fmt.Errorf("appicon.icns is empty")
				}
				// Remove the file

				return nil
			},
		},

		{
			name:          "should generate a Assets.car and icons.icns file when using the `IconComposerInput` flag and `MacAssetDir` flag",
			requireDarwin: true,
			setup: func(t *testing.T) *IconsOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				// Get the path to the example icon
				exampleIcon := filepath.Join(localDir, "build_assets", "appicon.icon")
				return &IconsOptions{
					IconComposerInput: exampleIcon,
					MacAssetDir:       t.TempDir(),
				}
			},
			wantErr: false,
			test: func(t *testing.T, options *IconsOptions) error {
				carPath := filepath.Join(options.MacAssetDir, "Assets.car")
				icnsPath := filepath.Join(options.MacAssetDir, "icons.icns")
				f, err := os.Stat(carPath)
				if err != nil {
					return err
				}
				if f.IsDir() {
					return fmt.Errorf("Assets.car is a directory")
				}
				if f.Size() == 0 {
					return fmt.Errorf("Assets.car is empty")
				}
				f, err = os.Stat(icnsPath)
				if err != nil {
					return err
				}
				if f.IsDir() {
					return fmt.Errorf("icons.icns is a directory")
				}
				if f.Size() == 0 {
					return fmt.Errorf("icons.icns is empty")
				}
				return nil
			},
		},
		{
			name:             "should return a descriptive error when icon composer assets are unsupported without fallback input",
			requireNonDarwin: true,
			setup: func(t *testing.T) *IconsOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				exampleIcon := filepath.Join(localDir, "build_assets", "appicon.icon")
				return &IconsOptions{
					IconComposerInput: exampleIcon,
					MacAssetDir:       t.TempDir(),
				}
			},
			wantErr:         true,
			wantErrContains: "mac asset generation requires macOS 26 or later",
		},
		{
			name:             "should fall back to image-based mac icon generation when icon composer assets are unsupported",
			requireNonDarwin: true,
			setup: func(t *testing.T) *IconsOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				exampleIcon := filepath.Join(localDir, "build_assets", "appicon.icon")
				examplePNG := filepath.Join(localDir, "build_assets", "appicon.png")
				return &IconsOptions{
					Input:             examplePNG,
					MacFilename:       "appicon.icns",
					IconComposerInput: exampleIcon,
					MacAssetDir:       t.TempDir(),
				}
			},
			wantErr: false,
			test: func(t *testing.T, options *IconsOptions) error {
				f, err := os.Stat(options.MacFilename)
				if err != nil {
					return err
				}
				defer func() {
					err = os.Remove(options.MacFilename)
					if err != nil {
						panic(err)
					}
				}()
				if f.IsDir() {
					return fmt.Errorf("%s is a directory", options.MacFilename)
				}
				if f.Size() == 0 {
					return fmt.Errorf("%s is empty", options.MacFilename)
				}
				if _, err := os.Stat(filepath.Join(options.MacAssetDir, "Assets.car")); !os.IsNotExist(err) {
					return fmt.Errorf("expected no Assets.car fallback artifact, got err=%v", err)
				}
				if _, err := os.Stat(filepath.Join(options.MacAssetDir, "icons.icns")); !os.IsNotExist(err) {
					return fmt.Errorf("expected no icons.icns fallback artifact, got err=%v", err)
				}
				return nil
			},
		},
		{
			name: "should generate a small .ico file when using the `input` flag and `sizes` flag",
			setup: func(t *testing.T) *IconsOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				// Get the path to the example icon
				exampleIcon := filepath.Join(localDir, "build_assets", "appicon.png")
				return &IconsOptions{
					Input:           exampleIcon,
					Sizes:           "16",
					WindowsFilename: "appicon.ico",
				}
			},
			wantErr: false,
			test: func(t *testing.T, options *IconsOptions) error {
				// the file `appicon.ico` should be created in the current directory
				// check for the existence of the file
				f, err := os.Stat("appicon.ico")
				if err != nil {
					return err
				}
				defer func() {
					err := os.Remove("appicon.ico")
					if err != nil {
						panic(err)
					}
				}()
				// The size of the file should be 571 bytes
				if f.Size() != 571 {
					return fmt.Errorf("appicon.ico is not the correct size. Got %d", f.Size())
				}
				if f.IsDir() {
					return fmt.Errorf("appicon.ico is a directory")
				}
				if f.Size() == 0 {
					return fmt.Errorf("appicon.ico is empty")
				}
				return nil
			},
		},
		{
			name: "should error if no input file is provided",
			setup: func(t *testing.T) *IconsOptions {
				return &IconsOptions{}
			},
			wantErr: true,
		},
		{
			name: "should error if neither mac or windows filename is provided",
			setup: func(t *testing.T) *IconsOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				// Get the path to the example icon
				exampleIcon := filepath.Join(localDir, "build_assets", "appicon.png")
				return &IconsOptions{
					Input: exampleIcon,
				}
			},
			wantErr: true,
		},
		{
			name: "should error if bad sizes provided",
			setup: func(t *testing.T) *IconsOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				// Get the path to the example icon
				exampleIcon := filepath.Join(localDir, "build_assets", "appicon.png")
				return &IconsOptions{
					Input:           exampleIcon,
					WindowsFilename: "appicon.ico",
					Sizes:           "bad",
				}
			},
			wantErr: true,
		},
		{
			name: "should ignore 0 size",
			setup: func(t *testing.T) *IconsOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				// Get the path to the example icon
				exampleIcon := filepath.Join(localDir, "build_assets", "appicon.png")
				return &IconsOptions{
					Input:           exampleIcon,
					WindowsFilename: "appicon.ico",
					Sizes:           "0,16",
				}
			},
			wantErr: false,
			test: func(t *testing.T, options *IconsOptions) error {
				// Test the file exists and has 571 bytes
				f, err := os.Stat("appicon.ico")
				if err != nil {
					return err
				}
				defer func() {
					err := os.Remove("appicon.ico")
					if err != nil {
						panic(err)
					}
				}()
				if f.Size() != 571 {
					return fmt.Errorf("appicon.ico is not the correct size. Got %d", f.Size())
				}
				if f.IsDir() {
					return fmt.Errorf("appicon.ico is a directory")
				}
				if f.Size() == 0 {
					return fmt.Errorf("appicon.ico is empty")
				}
				return nil
			},
		},
		{
			name: "should error if the input file does not exist",
			setup: func(t *testing.T) *IconsOptions {
				return &IconsOptions{
					Input:           "doesnotexist.png",
					WindowsFilename: "appicon.ico",
				}
			},
			wantErr: true,
		},
		{
			name: "should error if the input file is not a png",
			setup: func(t *testing.T) *IconsOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				return &IconsOptions{
					Input:           thisFile,
					WindowsFilename: "appicon.ico",
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requireDarwin && (runtime.GOOS != "darwin" || os.Getenv("CI") != "") {
				t.Skip("Assets.car generation is only supported on macOS and not in CI")
			}
			if tt.requireNonDarwin && runtime.GOOS == "darwin" {
				t.Skip("unsupported-platform behavior is only exercised on non-macOS hosts")
			}

			options := tt.setup(t)
			err := GenerateIcons(options)
			if tt.requireDarwin && err != nil {
				var notSupported *macAssetNotSupportedError
				if errors.As(err, &notSupported) {
					t.Skip(notSupported.Error())
				}
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateIcon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErrContains != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErrContains)) {
				t.Errorf("GenerateIcon() error = %v, want error containing %q", err, tt.wantErrContains)
				return
			}
			if tt.test != nil {
				if err := tt.test(t, options); err != nil {
					t.Errorf("GenerateIcon() test error = %v", err)
				}
			}
		})
	}
}
