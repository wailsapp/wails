package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGenerateIcon(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *IconsOptions
		wantErr bool
		test    func() error
	}{
		{
			name: "should generate an icon when using the `example` flag",
			setup: func() *IconsOptions {
				return &IconsOptions{
					Example: true,
				}
			},
			wantErr: false,
			test: func() error {
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
			setup: func() *IconsOptions {
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
			test: func() error {
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
			setup: func() *IconsOptions {
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
			test: func() error {
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
			name: "should generate a small .ico file when using the `input` flag and `sizes` flag",
			setup: func() *IconsOptions {
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
			test: func() error {
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
			setup: func() *IconsOptions {
				return &IconsOptions{}
			},
			wantErr: true,
		},
		{
			name: "should error if neither mac or windows filename is provided",
			setup: func() *IconsOptions {
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
			setup: func() *IconsOptions {
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
			setup: func() *IconsOptions {
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
			test: func() error {
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
			setup: func() *IconsOptions {
				return &IconsOptions{
					Input:           "doesnotexist.png",
					WindowsFilename: "appicon.ico",
				}
			},
			wantErr: true,
		},
		{
			name: "should error if the input file is not a png",
			setup: func() *IconsOptions {
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
			options := tt.setup()
			err := GenerateIcons(options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateIcon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.test != nil {
				if err := tt.test(); err != nil {
					t.Errorf("GenerateIcon() test error = %v", err)
				}
			}
		})
	}
}
