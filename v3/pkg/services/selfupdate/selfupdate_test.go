package selfupdate

import (
	"bytes"
	"context"
	"testing"
)

func TestNewService(t *testing.T) {
	t.Run("nil config creates default", func(t *testing.T) {
		s := New(nil)
		if s == nil {
			t.Fatal("expected service, got nil")
		}
		if s.config == nil {
			t.Fatal("expected config, got nil")
		}
	})

	t.Run("config is stored", func(t *testing.T) {
		cfg := &Config{
			CurrentVersion: "1.0.0",
			Provider:       "github",
		}
		s := New(cfg)
		if s.config.CurrentVersion != "1.0.0" {
			t.Errorf("expected CurrentVersion=1.0.0, got %s", s.config.CurrentVersion)
		}
	})

	t.Run("default max download size", func(t *testing.T) {
		s := New(nil)
		if s.config.MaxDownloadSize != DefaultMaxDownloadSize {
			t.Errorf("expected default MaxDownloadSize=%d, got %d",
				DefaultMaxDownloadSize, s.config.MaxDownloadSize)
		}
	})

	t.Run("default check timeout", func(t *testing.T) {
		s := New(nil)
		if s.config.CheckTimeout != DefaultCheckTimeout {
			t.Errorf("expected default CheckTimeout=%v, got %v",
				DefaultCheckTimeout, s.config.CheckTimeout)
		}
	})
}

func TestServiceName(t *testing.T) {
	s := New(nil)
	name := s.ServiceName()
	if name != "github.com/wailsapp/wails/v3/pkg/services/selfupdate" {
		t.Errorf("unexpected service name: %s", name)
	}
}

func TestGetCurrentVersion(t *testing.T) {
	s := New(&Config{CurrentVersion: "2.0.0"})
	v := s.GetCurrentVersion()
	if v != "2.0.0" {
		t.Errorf("expected 2.0.0, got %s", v)
	}
}

func TestProviderRegistry(t *testing.T) {
	t.Run("github provider is registered", func(t *testing.T) {
		if !HasProvider("github") {
			t.Error("github provider should be registered")
		}
	})

	t.Run("unknown provider returns error", func(t *testing.T) {
		_, err := GetProvider("unknown")
		if err == nil {
			t.Error("expected error for unknown provider")
		}
	})

	t.Run("available providers includes github", func(t *testing.T) {
		providers := AvailableProviders()
		found := false
		for _, p := range providers {
			if p == "github" {
				found = true
				break
			}
		}
		if !found {
			t.Error("github should be in available providers")
		}
	})
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1, v2 string
		want   int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.0.1", -1},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"1.1.0", "1.0.0", 1},
		{"1.0.0", "1.1.0", -1},
		{"1.0", "1.0.0", 0},
		{"1", "1.0.0", 0},
		{"10.0.0", "9.0.0", 1},
		{"0.0.1", "0.0.2", -1},
		// Pre-release tests (M1)
		{"1.0.0-alpha", "1.0.0", -1},
		{"1.0.0", "1.0.0-alpha", 1},
		{"1.0.0-alpha", "1.0.0-beta", -1},
		{"1.0.0-beta", "1.0.0-alpha", 1},
		{"1.0.0-alpha.1", "1.0.0-alpha.2", -1},
		{"1.0.0-rc.1", "1.0.0", -1},
	}

	for _, tt := range tests {
		t.Run(tt.v1+" vs "+tt.v2, func(t *testing.T) {
			got := compareVersions(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

func TestSplitPrerelease(t *testing.T) {
	tests := []struct {
		input      string
		wantBase   string
		wantPre    string
	}{
		{"1.0.0", "1.0.0", ""},
		{"1.0.0-alpha", "1.0.0", "alpha"},
		{"1.0.0-beta.1", "1.0.0", "beta.1"},
		{"1.0.0-rc.1", "1.0.0", "rc.1"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			base, pre := splitPrerelease(tt.input)
			if base != tt.wantBase {
				t.Errorf("base: expected %q, got %q", tt.wantBase, base)
			}
			if pre != tt.wantPre {
				t.Errorf("pre: expected %q, got %q", tt.wantPre, pre)
			}
		})
	}
}

func TestPatternVariables(t *testing.T) {
	vars := PatternVariables{
		Name:    "myapp",
		Version: "1.0.0",
		GOOS:    "darwin",
		GOARCH:  "arm64",
	}

	pattern := "{name}_{goos}_{goarch}"
	result := ResolveAssetPattern(pattern, vars)
	expected := "myapp_darwin_arm64"

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestFindMatchingAsset(t *testing.T) {
	assets := []string{
		"myapp_darwin_amd64.tar.gz",
		"myapp_darwin_arm64.tar.gz",
		"myapp_linux_amd64.tar.gz",
		"myapp_windows_amd64.zip",
		"checksums.txt",
	}

	vars := PatternVariables{
		Name:   "myapp",
		GOOS:   "darwin",
		GOARCH: "arm64",
	}

	match, _ := FindMatchingAsset(assets, DefaultAssetPattern, vars)
	if match != "myapp_darwin_arm64.tar.gz" {
		t.Errorf("expected myapp_darwin_arm64.tar.gz, got %s", match)
	}
}

func TestExtractPatternVariables(t *testing.T) {
	t.Run("detects linux amd64", func(t *testing.T) {
		vars, ok := ExtractPatternVariables("myapp_linux_amd64.tar.gz")
		if !ok {
			t.Error("expected extraction to succeed")
		}
		if vars.GOOS != "linux" {
			t.Errorf("expected linux, got %s", vars.GOOS)
		}
		if vars.GOARCH != "amd64" {
			t.Errorf("expected amd64, got %s", vars.GOARCH)
		}
	})

	t.Run("detects darwin arm64", func(t *testing.T) {
		vars, ok := ExtractPatternVariables("myapp_darwin_arm64.dmg")
		if !ok {
			t.Error("expected extraction to succeed")
		}
		if vars.GOOS != "darwin" {
			t.Errorf("expected darwin, got %s", vars.GOOS)
		}
		if vars.GOARCH != "arm64" {
			t.Errorf("expected arm64, got %s", vars.GOARCH)
		}
	})

	t.Run("returns false for unrecognizable", func(t *testing.T) {
		_, ok := ExtractPatternVariables("checksums.txt")
		if ok {
			t.Error("expected extraction to fail for checksums.txt")
		}
	})
}

func TestVerifyChecksum(t *testing.T) {
	data := []byte("hello world")
	// SHA256 of "hello world"
	validChecksum := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	t.Run("valid checksum passes", func(t *testing.T) {
		err := VerifyChecksum(data, validChecksum)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid checksum fails", func(t *testing.T) {
		err := VerifyChecksum(data, "invalid")
		if err == nil {
			t.Error("expected error for invalid checksum")
		}
	})

	t.Run("empty checksum is skipped", func(t *testing.T) {
		err := VerifyChecksum(data, "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestComputeChecksum(t *testing.T) {
	data := []byte("hello world")
	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	checksum := ComputeChecksum(data)
	if checksum != expected {
		t.Errorf("expected %s, got %s", expected, checksum)
	}
}

func TestNewVerifier(t *testing.T) {
	t.Run("empty key returns error", func(t *testing.T) {
		_, err := NewVerifier("")
		if err == nil {
			t.Error("expected error for empty public key")
		}
	})

	t.Run("invalid base64 returns error", func(t *testing.T) {
		_, err := NewVerifier("not-valid-base64!!!")
		if err == nil {
			t.Error("expected error for invalid base64")
		}
	})
}

func TestGenerateAndVerifySignature(t *testing.T) {
	// Generate key pair
	pubKey, privKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	data := []byte("test data to sign")

	// Sign the data
	signature, err := SignData(data, privKey)
	if err != nil {
		t.Fatalf("failed to sign data: %v", err)
	}

	// Create verifier
	verifier, err := NewVerifier(pubKey)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	// Verify the signature
	err = verifier.VerifySignature(data, signature)
	if err != nil {
		t.Errorf("signature verification failed: %v", err)
	}

	// Verify with wrong data should fail
	err = verifier.VerifySignature([]byte("wrong data"), signature)
	if err == nil {
		t.Error("expected verification to fail with wrong data")
	}

	// Verify with empty signature should fail (C1)
	err = verifier.VerifySignature(data, "")
	if err == nil {
		t.Error("expected error for empty signature when verifier has public key")
	}
}

func TestGitHubProvider(t *testing.T) {
	t.Run("configure requires owner", func(t *testing.T) {
		p := NewGitHubProvider()
		err := p.Configure(context.Background(), &ProviderConfig{
			Settings: map[string]any{
				"repo": "myrepo",
			},
		})
		if err == nil {
			t.Error("expected error for missing owner")
		}
	})

	t.Run("configure requires repo", func(t *testing.T) {
		p := NewGitHubProvider()
		err := p.Configure(context.Background(), &ProviderConfig{
			Settings: map[string]any{
				"owner": "myorg",
			},
		})
		if err == nil {
			t.Error("expected error for missing repo")
		}
	})

	t.Run("configure validates owner format", func(t *testing.T) {
		p := NewGitHubProvider()
		err := p.Configure(context.Background(), &ProviderConfig{
			Settings: map[string]any{
				"owner": "../traversal",
				"repo":  "myrepo",
			},
		})
		if err == nil {
			t.Error("expected error for invalid owner format")
		}
	})

	t.Run("configure validates repo format", func(t *testing.T) {
		p := NewGitHubProvider()
		err := p.Configure(context.Background(), &ProviderConfig{
			Settings: map[string]any{
				"owner": "myorg",
				"repo":  "my/repo",
			},
		})
		if err == nil {
			t.Error("expected error for invalid repo format")
		}
	})

	t.Run("configure rejects HTTP baseURL", func(t *testing.T) {
		p := NewGitHubProvider()
		err := p.Configure(context.Background(), &ProviderConfig{
			Settings: map[string]any{
				"owner":   "myorg",
				"repo":    "myrepo",
				"baseURL": "http://insecure.example.com",
			},
		})
		if err == nil {
			t.Error("expected error for HTTP baseURL")
		}
	})

	t.Run("configure with valid settings", func(t *testing.T) {
		p := NewGitHubProvider()
		err := p.Configure(context.Background(), &ProviderConfig{
			CurrentVersion: "1.0.0",
			Settings: map[string]any{
				"owner": "myorg",
				"repo":  "myrepo",
			},
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("name returns github", func(t *testing.T) {
		p := NewGitHubProvider()
		if p.Name() != "github" {
			t.Errorf("expected github, got %s", p.Name())
		}
	})

	t.Run("verify requires signature when public key set", func(t *testing.T) {
		pubKey, _, _ := GenerateKeyPair()
		p := NewGitHubProvider()
		_ = p.Configure(context.Background(), &ProviderConfig{
			CurrentVersion: "1.0.0",
			PublicKey:       pubKey,
			Settings: map[string]any{
				"owner": "myorg",
				"repo":  "myrepo",
			},
		})

		// Missing signature should error when public key is set (C1)
		err := p.VerifyUpdate(context.Background(), &UpdateResult{
			Signature: "", // No signature
		}, bytes.NewReader([]byte("data")))
		if err == nil {
			t.Error("expected error when signature is missing but public key is configured")
		}
	})

	t.Run("validate download URL rejects non-github", func(t *testing.T) {
		p := NewGitHubProvider()
		p.baseURL = "https://api.github.com"

		err := p.validateDownloadURL("https://evil.example.com/malware.exe")
		if err == nil {
			t.Error("expected error for non-GitHub download URL")
		}
	})

	t.Run("validate download URL accepts github domains", func(t *testing.T) {
		p := NewGitHubProvider()
		p.baseURL = "https://api.github.com"

		// github.com is allowed
		err := p.validateDownloadURL("https://github.com/myorg/myrepo/releases/download/v1.0.0/app.tar.gz")
		if err != nil {
			t.Errorf("unexpected error for github.com URL: %v", err)
		}

		// objects.githubusercontent.com is allowed (actual download host)
		err = p.validateDownloadURL("https://objects.githubusercontent.com/github-production-release-asset/12345/app.tar.gz")
		if err != nil {
			t.Errorf("unexpected error for objects.githubusercontent.com URL: %v", err)
		}

		// github-releases.githubusercontent.com is allowed
		err = p.validateDownloadURL("https://github-releases.githubusercontent.com/12345/app.tar.gz")
		if err != nil {
			t.Errorf("unexpected error for github-releases.githubusercontent.com URL: %v", err)
		}
	})

	t.Run("validate download URL rejects HTTP", func(t *testing.T) {
		p := NewGitHubProvider()
		p.baseURL = "https://api.github.com"

		err := p.validateDownloadURL("http://github.com/download/v1.0.0/app.tar.gz")
		if err == nil {
			t.Error("expected error for HTTP download URL")
		}
	})
}
