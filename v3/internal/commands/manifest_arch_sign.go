package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func (h *manifestHandler) signArch(ctx context.Context, spec pipeline.SignSpec) (pipeline.RunResult, error) {
	keyPath := spec.Config.Certificate
	if !filepath.IsAbs(keyPath) {
		keyPath = filepath.Join(h.root, keyPath)
	}
	password := environmentValues(mergeEnvironment(os.Environ(), h.environment))[spec.Config.Credential]
	if spec.Config.Credential != "" && password == "" {
		return pipeline.RunResult{}, fmt.Errorf("Arch signing credential %s is empty", spec.Config.Credential)
	}
	input, err := existingArtifactPathInsideProject(h.root, spec.Input)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	err = signArchPackage(ctx, input, keyPath, password)
	return pipeline.RunResult{Output: spec.Input + ".sig"}, err
}

func signArchPackage(ctx context.Context, input, keyPath, password string) error {
	if keyPath == "" {
		return fmt.Errorf("Arch signing requires an exported armored private-key file")
	}
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("read Arch private key: %w", err)
	}
	defer clear(key)
	entities, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(key))
	if err != nil {
		return fmt.Errorf("parse Arch armored private key: %w", err)
	}
	signer, err := archSigner(entities, []byte(password))
	if err != nil {
		return err
	}
	return producePathTransactional(input+".sig", ".arch-sign-*", func(staged string) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		message, err := os.Open(input)
		if err != nil {
			return err
		}
		defer message.Close()
		signature, err := os.OpenFile(staged, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		signErr := openpgp.DetachSign(signature, signer, &signingReader{ctx: ctx, file: message}, nil)
		closeErr := signature.Close()
		if signErr != nil {
			return fmt.Errorf("sign Arch package: %w", signErr)
		}
		if closeErr != nil {
			return closeErr
		}
		return ctx.Err()
	})
}

type signingReader struct {
	ctx  context.Context
	file *os.File
}

func (r *signingReader) Read(data []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.file.Read(data)
}

func archSigner(entities openpgp.EntityList, password []byte) (*openpgp.Entity, error) {
	defer clear(password)
	for _, entity := range entities {
		keys := []*packet.PrivateKey{entity.PrivateKey}
		for _, subkey := range entity.Subkeys {
			keys = append(keys, subkey.PrivateKey)
		}
		for _, key := range keys {
			if key != nil && key.Encrypted {
				if err := key.Decrypt(password); err != nil {
					return nil, fmt.Errorf("decrypt Arch private key: %w", err)
				}
			}
		}
		if key, ok := entity.SigningKey(time.Now()); ok && key.PrivateKey != nil {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("Arch signing requires an armored private key with a signing-capable key")
}
