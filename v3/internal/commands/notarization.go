package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/bundlearchive"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/notaryservice"
	"github.com/wailsapp/wails/v3/internal/notaryticket"
)

type notaryClient interface {
	Submit(context.Context, string) (notaryservice.Submission, error)
	Wait(context.Context, string, time.Duration) (notaryservice.Submission, error)
	Log(context.Context, string) ([]byte, error)
}

var newNotaryClient = func(credentials notaryservice.Credentials) (notaryClient, error) {
	return notaryservice.New(credentials, nil)
}

var stapleNotarization = staplePortableNotarization

func staplePortableNotarization(ctx context.Context, input string) error {
	digests, err := notaryticket.Digests(input)
	if err != nil {
		return fmt.Errorf("inspect notarization ticket target: %w", err)
	}
	ticket, err := (notaryticket.Client{}).Fetch(ctx, digests)
	if err != nil {
		return err
	}
	return produceNamedPathTransactional(input, ".notary-staple-*", filepath.Base(input), func(staged string) error {
		return notaryticket.Staple(ctx, input, staged, ticket)
	})
}

func usesNotaryAPI(options *flags.Sign) bool {
	return options.NotaryKey != "" || options.NotaryKeyID != "" || options.NotaryIssuerID != ""
}

func notarizeMacOSWithAPI(ctx context.Context, options *flags.Sign) error {
	key, err := readNotaryKey(options.NotaryKey)
	if err != nil {
		return err
	}
	defer clear(key)
	client, err := newNotaryClient(notaryservice.Credentials{KeyID: options.NotaryKeyID, IssuerID: options.NotaryIssuerID, PrivateKeyPEM: key})
	if err != nil {
		return err
	}
	archive := options.Input
	if strings.EqualFold(filepath.Ext(options.Input), ".app") {
		directory, err := os.MkdirTemp("", "wails-notary-upload-*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(directory)
		archive = filepath.Join(directory, "application.zip")
		if err := bundlearchive.Create(ctx, archive, options.Input); err != nil {
			return err
		}
	}
	submission, err := client.Submit(ctx, archive)
	if err != nil {
		return fmt.Errorf("notarization submission %s: %w", submission.ID, err)
	}
	result, err := client.Wait(ctx, submission.ID, 5*time.Second)
	if err != nil {
		return fmt.Errorf("notarization submission %s: %w", submission.ID, err)
	}
	if result.Status != "Accepted" {
		return rejectedNotarization(ctx, client, submission.ID, result.Status)
	}
	if err := stapleNotarization(ctx, options.Input); err != nil {
		return fmt.Errorf("Apple accepted submission %s, but %w", submission.ID, err)
	}
	pterm.Success.Printfln("Notarization accepted and stapled: %s", submission.ID)
	return nil
}

func readNotaryKey(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read notary API key: %w", err)
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, (64<<10)+1))
	if err != nil || len(data) > 64<<10 {
		clear(data)
		return nil, fmt.Errorf("notary API key must be a readable PEM file of at most 64 KiB")
	}
	return data, nil
}

func rejectedNotarization(ctx context.Context, client notaryClient, id, status string) error {
	log, err := client.Log(ctx, id)
	if err != nil {
		return fmt.Errorf("notarization submission %s was %s; cannot retrieve log: %w", id, status, err)
	}
	file, err := os.CreateTemp("", "wails-notary-rejection-*.json")
	if err != nil {
		return fmt.Errorf("notarization submission %s was %s; cannot save log: %w", id, status, err)
	}
	_, writeErr := file.Write(log)
	closeErr := file.Close()
	if writeErr != nil || closeErr != nil {
		_ = os.Remove(file.Name())
		return fmt.Errorf("notarization submission %s was %s; cannot save complete log", id, status)
	}
	return fmt.Errorf("notarization submission %s was %s; log saved to %s", id, status, file.Name())
}

func submitNotarizationProfile(ctx context.Context, path, profile string) error {
	output, err := runManifestCommand(ctx, filepath.Dir(path), nil, "xcrun", "notarytool", "submit", path, "--keychain-profile", profile, "--wait")
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil {
		return fmt.Errorf("notarization failed: %w: %s", err, output)
	}
	return nil
}
