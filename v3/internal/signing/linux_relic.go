package signing

import (
	"bytes"
	"context"
	"crypto"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	rpmutils "github.com/sassoftware/go-rpmutils"
	"github.com/sassoftware/relic/v8/lib/binpatch"
	"github.com/sassoftware/relic/v8/lib/signdeb"
)

// LinuxRelicSigner uses the relic library for cross-platform Linux package signing
type LinuxRelicSigner struct{}

// NewLinuxRelicSigner creates a new relic-based Linux signer
func NewLinuxRelicSigner() *LinuxRelicSigner {
	return &LinuxRelicSigner{}
}

// Backend returns the signer backend type
func (s *LinuxRelicSigner) Backend() SignerBackend {
	return BackendRelic
}

// Available returns true - relic is always available as it's a Go library
func (s *LinuxRelicSigner) Available() bool {
	return true
}

// Sign signs a Linux package (DEB or RPM) using the relic library
func (s *LinuxRelicSigner) Sign(ctx context.Context, req SignRequest) (*SignResult, error) {
	// Load the PGP key
	if req.Certificate.PGPKeyPath == "" {
		return nil, fmt.Errorf("PGPKeyPath is required for Linux package signing")
	}

	entity, err := loadPGPKey(req.Certificate.PGPKeyPath, req.Certificate.PGPKeyPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to load PGP key: %w", err)
	}

	// Determine output path
	outputPath := req.OutputPath
	if outputPath == "" {
		outputPath = req.InputPath
	}

	// Detect package type and sign accordingly
	ext := getExtension(req.InputPath)
	switch ext {
	case ".deb":
		err = s.signDEB(ctx, req.InputPath, outputPath, entity, req)
	case ".rpm":
		err = s.signRPM(ctx, req.InputPath, outputPath, entity, req)
	default:
		return nil, fmt.Errorf("unsupported Linux package format: %s (supported: .deb, .rpm)", ext)
	}

	if err != nil {
		return nil, err
	}

	return &SignResult{
		Backend:    BackendRelic,
		OutputPath: outputPath,
	}, nil
}

// signDEB signs a Debian package
func (s *LinuxRelicSigner) signDEB(ctx context.Context, inputPath, outputPath string, entity *openpgp.Entity, req SignRequest) error {
	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Determine role (default: builder)
	role := "builder"
	if req.Description != "" {
		// Use description as role if it's a valid role
		switch req.Description {
		case "builder", "origin", "maint", "archive":
			role = req.Description
		}
	}

	// Sign the DEB
	sig, err := signdeb.Sign(inputFile, entity, crypto.SHA256, role)
	if err != nil {
		return fmt.Errorf("failed to sign DEB: %w", err)
	}

	// Re-open input for patching
	inputFile.Close()
	inputFile, err = os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to reopen input file: %w", err)
	}
	defer inputFile.Close()

	// Apply the signature patch
	if err := applyPatch(inputFile, outputPath, sig.PatchSet); err != nil {
		return fmt.Errorf("failed to apply signature: %w", err)
	}

	return nil
}

// signRPM signs an RPM package
func (s *LinuxRelicSigner) signRPM(ctx context.Context, inputPath, outputPath string, entity *openpgp.Entity, req SignRequest) error {
	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Sign the RPM
	config := &rpmutils.SignatureOptions{
		Hash: crypto.SHA256,
	}

	header, err := rpmutils.SignRpmStream(inputFile, entity.PrivateKey, config)
	if err != nil {
		return fmt.Errorf("failed to sign RPM: %w", err)
	}

	// Get the signature header blob
	blob, err := header.DumpSignatureHeader(true)
	if err != nil {
		return fmt.Errorf("failed to dump signature header: %w", err)
	}

	// Create patch to replace the signature header
	patch := binpatch.New()
	patch.Add(0, int64(header.OriginalSignatureHeaderSize()), blob)

	// Re-open input for patching
	inputFile.Close()
	inputFile, err = os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to reopen input file: %w", err)
	}
	defer inputFile.Close()

	// Apply the signature patch
	if err := applyPatch(inputFile, outputPath, patch); err != nil {
		return fmt.Errorf("failed to apply signature: %w", err)
	}

	return nil
}

// Verify verifies the signature on a Linux package
func (s *LinuxRelicSigner) Verify(ctx context.Context, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	ext := getExtension(path)
	switch ext {
	case ".deb":
		sigmap, err := signdeb.Verify(f, nil, false)
		if err != nil {
			return fmt.Errorf("DEB signature verification failed: %w", err)
		}
		if len(sigmap) == 0 {
			return fmt.Errorf("no signatures found in DEB")
		}
		return nil

	case ".rpm":
		_, sigs, err := rpmutils.Verify(f, nil)
		if err != nil {
			return fmt.Errorf("RPM signature verification failed: %w", err)
		}
		if len(sigs) == 0 {
			return fmt.Errorf("no signatures found in RPM")
		}
		return nil

	default:
		return fmt.Errorf("unsupported package format for verification: %s", ext)
	}
}

// loadPGPKey loads a PGP private key from a file
func loadPGPKey(keyPath, password string) (*openpgp.Entity, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	var reader io.Reader = bytes.NewReader(keyData)

	// Check if it's ASCII armored
	if bytes.HasPrefix(keyData, []byte("-----BEGIN PGP")) {
		block, err := armor.Decode(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to decode armored key: %w", err)
		}
		reader = block.Body
	}

	// Read the entity
	entity, err := openpgp.ReadEntity(packet.NewReader(reader))
	if err != nil {
		return nil, fmt.Errorf("failed to read PGP entity: %w", err)
	}

	// Decrypt the private key if encrypted
	if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
		if password == "" {
			return nil, fmt.Errorf("PGP key is encrypted but no password provided")
		}
		if err := entity.PrivateKey.Decrypt([]byte(password)); err != nil {
			return nil, fmt.Errorf("failed to decrypt private key: %w", err)
		}
	}

	// Also decrypt any subkeys
	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted {
			if password == "" {
				return nil, fmt.Errorf("PGP subkey is encrypted but no password provided")
			}
			if err := subkey.PrivateKey.Decrypt([]byte(password)); err != nil {
				return nil, fmt.Errorf("failed to decrypt subkey: %w", err)
			}
		}
	}

	return entity, nil
}

// applyPatch applies a binary patch to create the signed output
func applyPatch(input *os.File, outputPath string, patch *binpatch.PatchSet) error {
	// Get file info for permissions
	info, err := input.Stat()
	if err != nil {
		return err
	}

	// Create output file
	output, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer output.Close()

	// Apply the patch
	return patch.Apply(input, outputPath)
}

// isLinuxPackage checks if a file is a Linux package based on extension
func isLinuxPackage(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".deb") || strings.HasSuffix(lower, ".rpm")
}
