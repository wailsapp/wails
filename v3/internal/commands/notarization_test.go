package commands

import (
	"archive/zip"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/notaryservice"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

type fixtureNotaryClient struct {
	t       *testing.T
	status  string
	waitErr error
	archive string
}

func (c *fixtureNotaryClient) Submit(ctx context.Context, path string) (notaryservice.Submission, error) {
	c.archive = path
	archive, err := zip.OpenReader(path)
	if err != nil {
		c.t.Fatal(err)
	}
	defer archive.Close()
	if len(archive.File) != 2 || archive.File[1].Name != "Example.app/binary" {
		c.t.Fatal("notary submission did not preserve bundle root")
	}
	return notaryservice.Submission{ID: "test-submission"}, nil
}

func (c *fixtureNotaryClient) Wait(context.Context, string, time.Duration) (notaryservice.Submission, error) {
	return notaryservice.Submission{ID: "test-submission", Status: c.status}, c.waitErr
}

func (c *fixtureNotaryClient) Log(context.Context, string) ([]byte, error) {
	return []byte(`{"status":"Invalid","issues":[]}`), nil
}

func TestPortableNotarizationArchivesWaitsAndStaplesOnlyAccepted(t *testing.T) {
	for _, status := range []string{"Accepted", "Invalid", "Cancelled"} {
		t.Run(status, func(t *testing.T) {
			root := t.TempDir()
			app := filepath.Join(root, "Example.app")
			if err := os.Mkdir(app, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(app, "binary"), []byte("signed executable"), 0755); err != nil {
				t.Fatal(err)
			}
			key := filepath.Join(root, "AuthKey.p8")
			if err := os.WriteFile(key, []byte("test-only-key"), 0600); err != nil {
				t.Fatal(err)
			}
			client := &fixtureNotaryClient{t: t, status: status}
			if status == "Cancelled" {
				client.waitErr = context.Canceled
			}
			previousClient, previousStaple, previousSign := newNotaryClient, stapleNotarization, manifestSign
			t.Cleanup(func() {
				newNotaryClient, stapleNotarization, manifestSign = previousClient, previousStaple, previousSign
			})
			signed := false
			manifestSign = func(options *flags.Sign) error {
				if options.Notarize || !options.HardenedRuntime || options.Input == app {
					t.Fatal("API route must sign privately with hardened runtime before notarization")
				}
				signed = true
				return nil
			}
			newNotaryClient = func(credentials notaryservice.Credentials) (notaryClient, error) {
				if !signed {
					t.Fatal("notarization started before signing")
				}
				if credentials.KeyID != "API_KEY_ID" || credentials.IssuerID != "TEAM" || string(credentials.PrivateKeyPEM) != "test-only-key" {
					t.Fatal("incorrect API credential routing")
				}
				return client, nil
			}
			stapled := false
			stapleNotarization = func(ctx context.Context, input string) error {
				if input == app || filepath.Ext(input) != ".app" {
					t.Fatal("stapled the upload archive instead of the application")
				}
				stapled = true
				return nil
			}
			_, err := (&manifestHandler{root: root}).sign(t.Context(), pipeline.SignSpec{TargetOS: "darwin", Format: "app", Input: "Example.app", Config: manifest.SigningPlatform{Enabled: true, Identity: "Developer ID", Notarize: true, NotarizationKey: "AuthKey.p8", NotarizationKeyID: "API_KEY_ID", NotarizationIssuerID: "TEAM"}})
			if stapled != (status == "Accepted") || (err == nil) != stapled {
				t.Fatalf("status=%s stapled=%v err=%v", status, stapled, err)
			}
			if _, statErr := os.Stat(client.archive); !os.IsNotExist(statErr) {
				t.Fatal("temporary upload archive was retained")
			}
			if status == "Cancelled" && !errors.Is(err, context.Canceled) {
				t.Fatal("cancellation was lost", err)
			}
			if status == "Invalid" {
				_, log, found := strings.Cut(err.Error(), "log saved to ")
				if !found {
					t.Fatal("rejection did not retain diagnostic log", err)
				}
				defer os.Remove(log)
				if data, readErr := os.ReadFile(log); readErr != nil || !strings.Contains(string(data), "Invalid") {
					t.Fatal("missing rejection log", readErr)
				}
			}
		})
	}
}
