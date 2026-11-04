package keychain

import "testing"

func TestEnvironmentCredentials(t *testing.T) {
	for _, key := range []string{KeyWindowsCertPassword, KeyPGPPassword} {
		t.Run(key, func(t *testing.T) {
			value := "  test password\n🔒"
			t.Setenv("WAILS_"+toEnvName(key), value)
			got, err := Get(key)
			if err != nil || got != value {
				t.Fatalf("environment credential was not preserved: %v", err)
			}
			if !Exists(key) {
				t.Fatal("environment credential was not found")
			}
		})
	}
}
