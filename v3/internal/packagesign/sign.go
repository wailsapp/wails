// Package packagesign signs existing Linux package archives without invoking
// external tools. The archive payload is preserved byte for byte.
package packagesign

import (
	"bytes"
	"context"
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

type Options struct {
	Format   string
	KeyFile  string
	Password string
	// KeyID optionally selects a primary key or signing subkey from the export.
	// It accepts a full fingerprint or a 16-digit key ID.
	KeyID string
	// Role is the dpkg-sig role; an empty value means builder.
	Role string
}

// Sign writes a signed archive to output, which must differ from input.
// The caller owns transactional publication and removal of a partial output.
// KeyFile must contain an exported OpenPGP secret key, not a GnuPG key ID.
func Sign(ctx context.Context, input, output string, options Options) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if options.Format != "deb" && options.Format != "rpm" && options.Format != "archlinux" {
		return fmt.Errorf("unsupported package signing format %q", options.Format)
	}
	signer, err := loadSigner(options)
	if err != nil {
		return err
	}
	src, err := os.Open(input)
	if err != nil {
		return err
	}
	defer src.Close()
	info, err := src.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("package input must be a regular file")
	}
	if dst, err := os.Stat(output); err == nil && os.SameFile(info, dst) {
		return fmt.Errorf("package signing output must differ from input")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if options.Format == "deb" {
		return signDEB(ctx, src, output, signer, options.Role)
	}
	if options.Format == "archlinux" {
		return writeOutput(output, func(dst io.Writer) error {
			return openpgp.DetachSign(dst, signer.entity, contextReader{ctx, src}, signer.config)
		})
	}
	return signRPM(ctx, src, output, signer)
}

type signingKey struct {
	entity *openpgp.Entity
	key    *packet.PrivateKey
	config *packet.Config
}

func loadSigner(options Options) (*signingKey, error) {
	f, err := os.Open(options.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("read exported signing key: %w", err)
	}
	defer f.Close()
	const maxKeySize = 16 << 20
	data, err := io.ReadAll(io.LimitReader(f, maxKeySize+1))
	if err != nil {
		return nil, fmt.Errorf("read exported signing key: %w", err)
	}
	defer clear(data)
	if len(data) > maxKeySize {
		return nil, fmt.Errorf("exported signing key exceeds 16 MiB")
	}
	var entities openpgp.EntityList
	if bytes.HasPrefix(bytes.TrimSpace(data), []byte("-----BEGIN PGP")) {
		entities, err = openpgp.ReadArmoredKeyRing(bytes.NewReader(data))
	} else {
		entities, err = openpgp.ReadKeyRing(bytes.NewReader(data))
	}
	if err != nil {
		return nil, fmt.Errorf("parse exported signing key: %w", err)
	}
	selector := strings.TrimPrefix(strings.ToUpper(options.KeyID), "0X")
	if selector != "" {
		if _, err := hex.DecodeString(selector); err != nil || (len(selector) != 16 && len(selector) != 40 && len(selector) != 64) {
			return nil, fmt.Errorf("signing key selector must be a fingerprint or 16-digit key ID")
		}
	}
	now := time.Now()
	var chosen *signingKey
	for _, entity := range entities {
		id := uint64(0)
		if selector != "" && !matchesKey(entity.PrimaryKey, selector) {
			for _, sub := range entity.Subkeys {
				if matchesKey(sub.PublicKey, selector) {
					id = sub.PublicKey.KeyId
				}
			}
			if id == 0 {
				continue
			}
		}
		key, ok := entity.SigningKeyById(now, id)
		if !ok || key.PrivateKey == nil || key.PrivateKey.Dummy() {
			continue
		}
		if chosen != nil {
			return nil, fmt.Errorf("export contains multiple signing identities; select a fingerprint with KeyID")
		}
		chosen = &signingKey{entity: entity, key: key.PrivateKey, config: &packet.Config{DefaultHash: crypto.SHA256, SigningKeyId: key.PublicKey.KeyId, Time: func() time.Time { return now }}}
	}
	if chosen == nil {
		return nil, fmt.Errorf("export contains no valid private signing key matching %q", options.KeyID)
	}
	if chosen.key.Encrypted {
		password := []byte(options.Password)
		err := chosen.key.Decrypt(password)
		clear(password)
		if err != nil {
			return nil, fmt.Errorf("decrypt exported signing key: %w", err)
		}
	}
	return chosen, nil
}

func matchesKey(key *packet.PublicKey, selector string) bool {
	return selector == strings.ToUpper(hex.EncodeToString(key.Fingerprint)) || selector == fmt.Sprintf("%016X", key.KeyId)
}

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (r contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}

func writeOutput(path string, write func(io.Writer) error) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	err = write(f)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
}
