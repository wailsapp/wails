package signing

import (
	"bytes"
	"context"
	"crypto"
	"fmt"
	"os"

	"github.com/sassoftware/relic/v8/lib/authenticode"
	"github.com/sassoftware/relic/v8/lib/certloader"
)

// staticPasswordGetter implements passprompt.PasswordGetter with a static password
type staticPasswordGetter struct {
	password string
}

func (s *staticPasswordGetter) GetPasswd(prompt string) (string, error) {
	return s.password, nil
}

// WindowsRelicSigner uses the relic library for cross-platform Windows signing
type WindowsRelicSigner struct{}

// NewWindowsRelicSigner creates a new relic-based Windows signer
func NewWindowsRelicSigner() *WindowsRelicSigner {
	return &WindowsRelicSigner{}
}

// Backend returns the signer backend type
func (s *WindowsRelicSigner) Backend() SignerBackend {
	return BackendRelic
}

// Available returns true - relic is always available as it's a Go library
func (s *WindowsRelicSigner) Available() bool {
	return true
}

// Sign signs a Windows binary using the relic library
func (s *WindowsRelicSigner) Sign(ctx context.Context, req SignRequest) (*SignResult, error) {
	// Load the certificate
	if req.Certificate.PKCS12Path == "" {
		return nil, fmt.Errorf("PKCS12Path is required for cross-platform signing")
	}

	// Read the PKCS12 file
	p12Data, err := os.ReadFile(req.Certificate.PKCS12Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	// Create a password getter
	prompter := &staticPasswordGetter{password: req.Certificate.PKCS12Password}

	cert, err := certloader.ParsePKCS12(p12Data, prompter)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Determine output path
	outputPath := req.OutputPath
	if outputPath == "" {
		outputPath = req.InputPath
	}

	// Detect file type and sign accordingly
	ext := getExtension(req.InputPath)
	switch ext {
	case ".exe", ".dll", ".sys":
		err = s.signPE(ctx, req.InputPath, outputPath, cert, req)
	case ".ps1", ".psm1", ".psd1":
		err = s.signPowerShell(ctx, req.InputPath, outputPath, cert, req)
	default:
		// Try PE format as default
		err = s.signPE(ctx, req.InputPath, outputPath, cert, req)
	}

	if err != nil {
		return nil, err
	}

	return &SignResult{
		Backend:    BackendRelic,
		OutputPath: outputPath,
	}, nil
}

// signPE signs a PE (EXE/DLL) file
func (s *WindowsRelicSigner) signPE(ctx context.Context, inputPath, outputPath string, cert *certloader.Certificate, req SignRequest) error {
	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Calculate the digest
	hash := crypto.SHA256
	digest, err := authenticode.DigestPE(inputFile, hash, false)
	if err != nil {
		return fmt.Errorf("failed to digest PE file: %w", err)
	}

	// Create opus params for signature metadata
	var opusParams *authenticode.OpusParams
	if req.Description != "" || req.URL != "" {
		opusParams = &authenticode.OpusParams{
			Description: req.Description,
			URL:         req.URL,
		}
	}

	// Create the signature
	patchSet, _, err := digest.Sign(ctx, cert, opusParams)
	if err != nil {
		return fmt.Errorf("failed to sign PE file: %w", err)
	}

	// We need to re-open the file for patching
	inputFile.Close()
	inputFile, err = os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to reopen input file: %w", err)
	}
	defer inputFile.Close()

	// Apply the signature patch
	if err := patchSet.Apply(inputFile, outputPath); err != nil {
		return fmt.Errorf("failed to apply signature: %w", err)
	}

	return nil
}

// signPowerShell signs a PowerShell script
func (s *WindowsRelicSigner) signPowerShell(ctx context.Context, inputPath, outputPath string, cert *certloader.Certificate, req SignRequest) error {
	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Detect signature style
	style, ok := authenticode.GetSigStyle(inputPath)
	if !ok {
		style = authenticode.SigStyleHash
	}

	hash := crypto.SHA256
	digest, err := authenticode.DigestPowershell(inputFile, style, hash)
	if err != nil {
		return fmt.Errorf("failed to digest PowerShell script: %w", err)
	}

	// Create opus params
	var opusParams *authenticode.OpusParams
	if req.Description != "" || req.URL != "" {
		opusParams = &authenticode.OpusParams{
			Description: req.Description,
			URL:         req.URL,
		}
	}

	patchSet, _, err := digest.Sign(ctx, cert, opusParams)
	if err != nil {
		return fmt.Errorf("failed to sign PowerShell script: %w", err)
	}

	// Re-open input for patching
	inputFile.Close()
	inputFile, err = os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to reopen input file: %w", err)
	}
	defer inputFile.Close()

	// Apply the signature patch
	if err := patchSet.Apply(inputFile, outputPath); err != nil {
		return fmt.Errorf("failed to apply signature: %w", err)
	}

	return nil
}

// Verify verifies the signature on a Windows binary using relic
func (s *WindowsRelicSigner) Verify(ctx context.Context, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	ext := getExtension(path)
	switch ext {
	case ".exe", ".dll", ".sys":
		sigs, err := authenticode.VerifyPE(f, false)
		if err != nil {
			return fmt.Errorf("signature verification failed: %w", err)
		}
		if len(sigs) == 0 {
			return fmt.Errorf("no signatures found")
		}
		return nil

	case ".ps1", ".psm1", ".psd1":
		style, ok := authenticode.GetSigStyle(path)
		if !ok {
			style = authenticode.SigStyleHash
		}
		sig, err := authenticode.VerifyPowershell(f, style, false)
		if err != nil {
			return fmt.Errorf("PowerShell signature verification failed: %w", err)
		}
		if sig == nil {
			return fmt.Errorf("no signature found")
		}
		return nil

	default:
		// Try as PE
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sigs, err := authenticode.VerifyPE(bytes.NewReader(data), false)
		if err != nil {
			return fmt.Errorf("signature verification failed: %w", err)
		}
		if len(sigs) == 0 {
			return fmt.Errorf("no signatures found")
		}
		return nil
	}
}

// getExtension returns the lowercase file extension
func getExtension(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			ext := path[i:]
			// Convert to lowercase
			result := make([]byte, len(ext))
			for j := 0; j < len(ext); j++ {
				c := ext[j]
				if c >= 'A' && c <= 'Z' {
					c += 'a' - 'A'
				}
				result[j] = c
			}
			return string(result)
		}
		if path[i] == '/' || path[i] == '\\' {
			break
		}
	}
	return ""
}
