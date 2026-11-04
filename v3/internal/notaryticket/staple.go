package notaryticket

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/applesign"
	"github.com/wailsapp/wails/v3/internal/stagefile"
)

// Digests reads the signed CodeDirectories of a canonical app or signed DMG.
// It does not independently verify the app's signature or resource seal.
func Digests(path string) ([]Digest, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		path, err = applesign.BundleExecutable(path)
		if err != nil {
			return nil, err
		}
	} else if !info.Mode().IsRegular() {
		return nil, errors.New("ticket target must be a regular DMG or an app directory")
	}
	identities, err := applesign.Inspect(path)
	if err != nil {
		return nil, err
	}
	digests := make([]Digest, len(identities))
	for i, id := range identities {
		digests[i] = Digest{Algorithm: id.HashType, Value: id.CDHash}
	}
	return digests, validateDigests(digests)
}

// Staple creates a new output, preserving the input. Output must not exist;
// the caller can publish the completed output transactionally. Tickets must
// cover every architecture's CodeDirectory in the copied target.
func Staple(ctx context.Context, input, output string, ticket *Ticket) (err error) {
	if err = ctx.Err(); err != nil {
		return err
	}
	if ticket == nil || len(ticket.data) == 0 {
		return errors.New("missing verified notarization ticket")
	}
	info, err := os.Lstat(input)
	if err != nil {
		return err
	}
	if info.Mode().IsRegular() {
		return stapleDMG(ctx, input, output, ticket)
	}
	if info.IsDir() {
		err = applesign.CopyBundle(ctx, input, output)
	} else {
		return errors.New("ticket target must be a regular DMG or an app directory")
	}
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(output)
		}
	}()
	hashes, err := Digests(output)
	if err != nil {
		return err
	}
	if err = ticket.includes(hashes); err != nil {
		return err
	}
	path := filepath.Join(output, "Contents", "CodeResources")
	if existing, e := os.Lstat(path); e == nil && !existing.Mode().IsRegular() {
		return errors.New("existing bundle ticket is not a regular file")
	} else if e != nil && !os.IsNotExist(e) {
		return e
	}
	if err = os.WriteFile(path, ticket.data, 0644); err != nil {
		return err
	}
	return ctx.Err()
}

func stapleDMG(ctx context.Context, input, output string, ticket *Ticket) error {
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(filepath.Dir(output), ".notary-ticket-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	path := filepath.Join(stage, "image.dmg")
	if err := stagefile.Copy(ctx, input, path, 0, make([]byte, 64<<10)); err != nil {
		return err
	}
	hashes, err := Digests(path)
	if err != nil {
		return err
	}
	if err := ticket.includes(hashes); err != nil {
		return err
	}
	if err := attachDMG(path, ticket.data); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	// A hard link publishes the completed same-volume file without overwriting
	// an existing output. It needs no Windows symlink privilege.
	return os.Link(path, output)
}
