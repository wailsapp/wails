package selfupdate

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// MaxVerifySize is the maximum size of data that can be verified (500MB).
const MaxVerifySize = 500 * 1024 * 1024

// Verifier handles cryptographic verification of update artifacts.
type Verifier struct {
	publicKey ed25519.PublicKey
}

// NewVerifier creates a new Verifier with the given Ed25519 public key.
// The key should be base64-encoded (standard encoding).
//
// Returns an error if the public key is empty or invalid.
func NewVerifier(publicKeyBase64 string) (*Verifier, error) {
	if publicKeyBase64 == "" {
		return nil, fmt.Errorf("public key is required")
	}

	keyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	if len(keyBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key size: expected %d bytes, got %d",
			ed25519.PublicKeySize, len(keyBytes))
	}

	return &Verifier{
		publicKey: ed25519.PublicKey(keyBytes),
	}, nil
}

// VerifySignature verifies that the signature is valid for the given data.
// The signature should be base64-encoded.
//
// Returns an error if verification fails or if the signature is empty.
func (v *Verifier) VerifySignature(data []byte, signatureBase64 string) error {
	if v == nil || v.publicKey == nil {
		return fmt.Errorf("verifier not initialized: public key is required")
	}

	if signatureBase64 == "" {
		return fmt.Errorf("signature is required when public key is configured")
	}

	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	if len(signature) != ed25519.SignatureSize {
		return fmt.Errorf("invalid signature size: expected %d bytes, got %d",
			ed25519.SignatureSize, len(signature))
	}

	if !ed25519.Verify(v.publicKey, data, signature) {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

// VerifySignatureReader verifies a signature against data from a reader.
// This reads the content into memory (up to MaxVerifySize) to compute the signature.
func (v *Verifier) VerifySignatureReader(r io.Reader, signatureBase64 string) error {
	if v == nil || v.publicKey == nil {
		return fmt.Errorf("verifier not initialized: public key is required")
	}

	data, err := io.ReadAll(io.LimitReader(r, MaxVerifySize+1))
	if err != nil {
		return fmt.Errorf("failed to read data for verification: %w", err)
	}
	if int64(len(data)) > MaxVerifySize {
		return fmt.Errorf("data exceeds maximum size of %d bytes", MaxVerifySize)
	}

	return v.VerifySignature(data, signatureBase64)
}

// VerifyChecksum verifies that the data matches the expected SHA256 checksum.
// The checksum should be hex-encoded.
func VerifyChecksum(data []byte, expectedChecksumHex string) error {
	if expectedChecksumHex == "" {
		return nil // Checksum verification disabled
	}

	expectedChecksum, err := hex.DecodeString(expectedChecksumHex)
	if err != nil {
		return fmt.Errorf("failed to decode checksum: %w", err)
	}

	actualChecksum := sha256.Sum256(data)

	if !compareSlices(actualChecksum[:], expectedChecksum) {
		return fmt.Errorf("checksum mismatch: expected %s, got %s",
			expectedChecksumHex, hex.EncodeToString(actualChecksum[:]))
	}

	return nil
}

// VerifyChecksumReader verifies a checksum against data from a reader.
// Reads up to MaxVerifySize bytes.
func VerifyChecksumReader(r io.Reader, expectedChecksumHex string) error {
	if expectedChecksumHex == "" {
		return nil // Checksum verification disabled
	}

	data, err := io.ReadAll(io.LimitReader(r, MaxVerifySize+1))
	if err != nil {
		return fmt.Errorf("failed to read data for checksum: %w", err)
	}
	if int64(len(data)) > MaxVerifySize {
		return fmt.Errorf("data exceeds maximum size of %d bytes", MaxVerifySize)
	}

	return VerifyChecksum(data, expectedChecksumHex)
}

// ComputeChecksum computes the SHA256 checksum of the data.
// Returns the checksum as a hex-encoded string.
func ComputeChecksum(data []byte) string {
	checksum := sha256.Sum256(data)
	return hex.EncodeToString(checksum[:])
}

// ComputeChecksumReader computes the SHA256 checksum from a reader.
func ComputeChecksumReader(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", fmt.Errorf("failed to compute checksum: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// compareSlices compares two byte slices in constant time.
func compareSlices(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := range a {
		result |= a[i] ^ b[i]
	}
	return result == 0
}

// GenerateKeyPair generates a new Ed25519 key pair for signing updates.
// Returns the public and private keys as base64-encoded strings.
//
// Note: This is typically used in CI/CD pipelines, not at runtime.
func GenerateKeyPair() (publicKeyBase64, privateKeyBase64 string, err error) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate key pair: %w", err)
	}

	publicKeyBase64 = base64.StdEncoding.EncodeToString(publicKey)
	privateKeyBase64 = base64.StdEncoding.EncodeToString(privateKey)

	return publicKeyBase64, privateKeyBase64, nil
}

// SignData signs data with the given Ed25519 private key.
// The private key should be base64-encoded.
// Returns the signature as a base64-encoded string.
//
// Note: This is typically used in CI/CD pipelines, not at runtime.
func SignData(data []byte, privateKeyBase64 string) (string, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode private key: %w", err)
	}

	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("invalid private key size: expected %d bytes, got %d",
			ed25519.PrivateKeySize, len(privateKeyBytes))
	}

	privateKey := ed25519.PrivateKey(privateKeyBytes)
	signature := ed25519.Sign(privateKey, data)

	return base64.StdEncoding.EncodeToString(signature), nil
}
