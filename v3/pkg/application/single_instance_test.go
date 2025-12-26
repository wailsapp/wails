package application

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := [32]byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}

	plaintext := []byte("Hello, World! This is a test message.")

	encrypted, err := encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if encrypted == "" {
		t.Error("encrypted should not be empty")
	}

	if encrypted == string(plaintext) {
		t.Error("encrypted should be different from plaintext")
	}

	decrypted, err := decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("decrypted = %q, want %q", string(decrypted), string(plaintext))
	}
}

func TestEncryptDecrypt_EmptyData(t *testing.T) {
	key := [32]byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}

	plaintext := []byte{}

	encrypted, err := encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	decrypted, err := decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if len(decrypted) != 0 {
		t.Errorf("decrypted should be empty, got %d bytes", len(decrypted))
	}
}

func TestDecrypt_InvalidData(t *testing.T) {
	key := [32]byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}

	tests := []struct {
		name string
		data string
	}{
		{"invalid base64", "not-valid-base64!!!"},
		{"too short", "YWJj"}, // "abc" base64 encoded (3 bytes)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decrypt(key, tt.data)
			if err == nil {
				t.Error("decrypt should return error for invalid data")
			}
		})
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	key1 := [32]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}

	key2 := [32]byte{0xff, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}

	plaintext := []byte("Secret message")

	encrypted, err := encrypt(key1, plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	_, err = decrypt(key2, encrypted)
	if err == nil {
		t.Error("decrypt with wrong key should return error")
	}
}

func TestGetLockPath(t *testing.T) {
	uniqueID := "com.example.myapp"
	path := getLockPath(uniqueID)

	if path == "" {
		t.Error("getLockPath should return non-empty path")
	}

	expectedFileName := uniqueID + ".lock"
	actualFileName := filepath.Base(path)
	if actualFileName != expectedFileName {
		t.Errorf("filename = %q, want %q", actualFileName, expectedFileName)
	}

	// Path should be in temp directory
	// Use filepath.Clean to normalize paths (os.TempDir may have trailing slash on macOS)
	tmpDir := filepath.Clean(os.TempDir())
	if filepath.Dir(path) != tmpDir {
		t.Errorf("path should be in temp directory %q, got %q", tmpDir, filepath.Dir(path))
	}
}

func TestGetCurrentWorkingDir(t *testing.T) {
	dir := getCurrentWorkingDir()

	// Should return a non-empty path
	if dir == "" {
		t.Error("getCurrentWorkingDir should return non-empty path")
	}

	// Should match os.Getwd()
	expected, err := os.Getwd()
	if err != nil {
		t.Skipf("os.Getwd failed: %v", err)
	}

	if dir != expected {
		t.Errorf("getCurrentWorkingDir() = %q, want %q", dir, expected)
	}
}

func TestSecondInstanceData_Fields(t *testing.T) {
	data := SecondInstanceData{
		Args:       []string{"arg1", "arg2"},
		WorkingDir: "/home/user",
		AdditionalData: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	if len(data.Args) != 2 {
		t.Error("Args not set correctly")
	}
	if data.WorkingDir != "/home/user" {
		t.Error("WorkingDir not set correctly")
	}
	if len(data.AdditionalData) != 2 {
		t.Error("AdditionalData not set correctly")
	}
}

func TestSingleInstanceOptions_Defaults(t *testing.T) {
	opts := SingleInstanceOptions{}

	if opts.UniqueID != "" {
		t.Error("UniqueID should default to empty string")
	}
	if opts.OnSecondInstanceLaunch != nil {
		t.Error("OnSecondInstanceLaunch should default to nil")
	}
	if opts.AdditionalData != nil {
		t.Error("AdditionalData should default to nil")
	}
	if opts.ExitCode != 0 {
		t.Error("ExitCode should default to 0")
	}
	var zeroKey [32]byte
	if opts.EncryptionKey != zeroKey {
		t.Error("EncryptionKey should default to zero array")
	}
}

func TestSingleInstanceManager_Cleanup_Nil(t *testing.T) {
	// Calling cleanup on nil manager should not panic
	var m *singleInstanceManager
	m.cleanup() // Should not panic
}

func TestSingleInstanceManager_Cleanup_NilLock(t *testing.T) {
	// Calling cleanup with nil lock should not panic
	m := &singleInstanceManager{}
	m.cleanup() // Should not panic
}

func TestNewSingleInstanceManager_NilOptions(t *testing.T) {
	manager, err := newSingleInstanceManager(nil, nil)
	if err != nil {
		t.Errorf("newSingleInstanceManager(nil, nil) should not return error: %v", err)
	}
	if manager != nil {
		t.Error("newSingleInstanceManager(nil, nil) should return nil manager")
	}
}

func TestAlreadyRunningError(t *testing.T) {
	if alreadyRunningError == nil {
		t.Error("alreadyRunningError should not be nil")
	}
	if alreadyRunningError.Error() != "application is already running" {
		t.Errorf("alreadyRunningError.Error() = %q", alreadyRunningError.Error())
	}
}
