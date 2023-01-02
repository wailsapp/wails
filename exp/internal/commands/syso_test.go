package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGenerateSyso(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *SysoOptions
		wantErr bool
		test    func() error
	}{
		{
			name: "should generate example info and manifest files when using the `example` flag",
			setup: func() *SysoOptions {
				return &SysoOptions{
					Example: true,
				}
			},
			wantErr: false,
			test: func() error {
				// the file `info.json` should be created in the current directory
				// check for the existence of the file
				f, err := os.Stat("info.json")
				if err != nil {
					return err
				}
				m, err := os.Stat("wails.exe.manifest")
				if err != nil {
					return err
				}
				defer func() {
					err := os.Remove("info.json")
					err2 := os.Remove("wails.exe.manifest")
					if err != nil {
						panic(err)
					}
					if err2 != nil {
						panic(err2)
					}
				}()
				if f.IsDir() {
					return fmt.Errorf("info.json is a directory")
				}
				if f.Size() == 0 {
					return fmt.Errorf("info.json is empty")
				}
				if m.IsDir() {
					return fmt.Errorf("wails.exe.manifest is a directory")
				}
				if m.Size() == 0 {
					return fmt.Errorf("wails.exe.manifest is empty")
				}
				return nil
			},
		},
		{
			name: "should error if manifest filename is not provided",
			setup: func() *SysoOptions {
				return &SysoOptions{
					Manifest: "",
				}
			},
			wantErr: true,
		},
		{
			name: "should error if icon filename is not provided",
			setup: func() *SysoOptions {
				return &SysoOptions{
					Manifest: "test.manifest",
					Icon:     "",
				}
			},
			wantErr: true,
		},
		{
			name: "should error if icon filename does not exist",
			setup: func() *SysoOptions {
				return &SysoOptions{
					Manifest: "test.manifest",
					Icon:     "icon.ico",
				}
			},
			wantErr: true,
		},
		{
			name: "should error if icon is wrong format",
			setup: func() *SysoOptions {
				_, thisFile, _, _ := runtime.Caller(1)
				return &SysoOptions{
					Manifest: "test.manifest",
					Icon:     thisFile,
				}
			},
			wantErr: true,
		},
		{
			name: "should error if manifest filename does not exist",
			setup: func() *SysoOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				// Get the path to the example icon
				exampleIcon := filepath.Join(localDir, "examples", "icon.ico")
				return &SysoOptions{
					Manifest: "test.manifest",
					Icon:     exampleIcon,
				}
			},
			wantErr: true,
		},
		{
			name: "should error if manifest is wrong format",
			setup: func() *SysoOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				// Get the path to the example icon
				exampleIcon := filepath.Join(localDir, "examples", "icon.ico")
				return &SysoOptions{
					Manifest: exampleIcon,
					Icon:     exampleIcon,
				}
			},
			wantErr: true,
		},
		{
			name: "should error if info file does not exist",
			setup: func() *SysoOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				// Get the path to the example icon
				exampleIcon := filepath.Join(localDir, "examples", "icon.ico")
				// Get the path to the example manifest
				exampleManifest := filepath.Join(localDir, "examples", "wails.exe.manifest")
				return &SysoOptions{
					Manifest: exampleManifest,
					Icon:     exampleIcon,
					Info:     "doesnotexist.json",
				}
			},
			wantErr: true,
		},
		{
			name: "should error if info file is wrong format",
			setup: func() *SysoOptions {
				// Get the directory of this file
				_, thisFile, _, _ := runtime.Caller(1)
				localDir := filepath.Dir(thisFile)
				// Get the path to the example icon
				exampleIcon := filepath.Join(localDir, "examples", "icon.ico")
				// Get the path to the example manifest
				exampleManifest := filepath.Join(localDir, "examples", "wails.exe.manifest")
				return &SysoOptions{
					Manifest: exampleManifest,
					Icon:     exampleIcon,
					Info:     thisFile,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := tt.setup()
			err := GenerateSyso(options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSyso() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && tt.wantErr {
				println(err.Error())
				return
			}
			if tt.test != nil {
				if err := tt.test(); err != nil {
					t.Errorf("GenerateSyso() test error = %v", err)
				}
			}
		})
	}
}
