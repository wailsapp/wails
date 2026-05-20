package updater

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"math/big"
)

// verifier authenticates a downloaded artifact. The interface is unexported
// in v1: callers configure verification by populating Release.Verification
// (which the Updater's built-in registry maps to the concrete verifier).
// The interface stays unexported until there is a real third-party need for
// a custom algorithm — at which point we promote it to a public
// RegisterVerifier without breaking existing callers.
type verifier interface {
	// verify checks digest+/or signature against publicKey. The Updater
	// always computes the digest in a streaming pass during download and
	// passes it here, so verifiers do not re-hash. publicKey is nil when
	// the caller did not configure one.
	verify(digest []byte, v *Verification, publicKey []byte) error
}

// digestHasher returns a fresh hash.Hash for the digest algorithm named by
// algo. Unknown algorithms return a nil hash and an error.
func digestHasher(algo string) (hash.Hash, error) {
	switch algo {
	case "", "sha256":
		return sha256.New(), nil
	case "sha512":
		return sha512.New(), nil
	}
	return nil, fmt.Errorf("updater: unknown digest algorithm %q", algo)
}

// verifierFor returns the registered verifier for v.SignatureAlgo, or nil if
// v carries no signature (in which case only digest comparison applies).
func verifierFor(algo string) (verifier, error) {
	if algo == "" {
		return nil, nil
	}
	if vf, ok := verifierRegistry[algo]; ok {
		return vf, nil
	}
	return nil, fmt.Errorf("updater: unsupported signature algorithm %q", algo)
}

var verifierRegistry = map[string]verifier{
	"ed25519":    ed25519Verifier{},
	"ed25519ph":  ed25519phVerifier{},
	"ecdsa-p256": ecdsaP256Verifier{},
}

// runVerification is the single entry point used by the Updater. It enforces
// the contract: if Verification has a digest, it must match; if it has a
// signature, the signature must verify under configKey. Returns nil only when
// every present check passed.
//
// Signature verification uses configKey (Config.PublicKey) and nothing else.
// The release source does not get to choose its own trust anchor — that would
// defeat the purpose of pinning a key out-of-band at build time. Releases that
// carry a Signature without a configured key fail closed.
func runVerification(computedDigest []byte, v *Verification, configKey []byte) error {
	if v == nil {
		return nil // nothing to check
	}

	if len(v.Digest) > 0 {
		if !constantTimeEqual(computedDigest, v.Digest) {
			return errors.New("updater: digest mismatch")
		}
	}

	if v.SignatureAlgo == "" || len(v.Signature) == 0 {
		return nil // digest-only mode
	}

	vf, err := verifierFor(v.SignatureAlgo)
	if err != nil {
		return err
	}
	if len(configKey) == 0 {
		return errors.New("updater: signature requires a public key but none configured")
	}
	return vf.verify(computedDigest, v, configKey)
}

func constantTimeEqual(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

// --- ed25519 (raw, payload-signing) ---

type ed25519Verifier struct{}

func (ed25519Verifier) verify(digest []byte, v *Verification, publicKey []byte) error {
	pub, err := parseEd25519Public(publicKey)
	if err != nil {
		return err
	}
	// Raw Ed25519 signs the message (here, the digest). Callers that want to
	// sign full-file payloads should use ed25519ph instead and let the
	// Updater compute the digest for them.
	if !ed25519.Verify(pub, digest, v.Signature) {
		return errors.New("updater: ed25519 signature did not verify")
	}
	return nil
}

func parseEd25519Public(raw []byte) (ed25519.PublicKey, error) {
	if len(raw) == ed25519.PublicKeySize {
		return ed25519.PublicKey(raw), nil
	}
	if block, _ := pem.Decode(raw); block != nil {
		raw = block.Bytes
	}
	pubAny, err := x509.ParsePKIXPublicKey(raw)
	if err != nil {
		return nil, fmt.Errorf("updater: ed25519 public key parse: %w", err)
	}
	pub, ok := pubAny.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("updater: ed25519 public key has wrong type %T", pubAny)
	}
	return pub, nil
}

// --- ed25519ph (pre-hash) ---

type ed25519phVerifier struct{}

func (ed25519phVerifier) verify(digest []byte, v *Verification, publicKey []byte) error {
	pub, err := parseEd25519Public(publicKey)
	if err != nil {
		return err
	}
	// Ed25519ph signs SHA-512(message). The Updater always streams a SHA-512
	// hash in parallel when verification is configured with this algo.
	if len(digest) != sha512.Size {
		return fmt.Errorf("updater: ed25519ph requires sha512 digest, got %d bytes", len(digest))
	}
	if err := ed25519.VerifyWithOptions(pub, digest, v.Signature, &ed25519.Options{Hash: crypto.SHA512}); err != nil {
		return fmt.Errorf("updater: ed25519ph signature did not verify: %w", err)
	}
	return nil
}

// --- ecdsa P-256 over SHA-256 ---

type ecdsaP256Verifier struct{}

func (ecdsaP256Verifier) verify(digest []byte, v *Verification, publicKey []byte) error {
	pub, err := parseECDSAPublic(publicKey)
	if err != nil {
		return err
	}
	if pub.Curve != elliptic.P256() {
		return fmt.Errorf("updater: ecdsa-p256 requires P-256 key, got %s", pub.Curve.Params().Name)
	}
	// Accept either raw r||s (64 bytes for P-256) or ASN.1 DER.
	r, s, err := splitECDSASig(v.Signature)
	if err != nil {
		return err
	}
	if !ecdsa.Verify(pub, digest, r, s) {
		return errors.New("updater: ecdsa-p256 signature did not verify")
	}
	return nil
}

func parseECDSAPublic(raw []byte) (*ecdsa.PublicKey, error) {
	if block, _ := pem.Decode(raw); block != nil {
		raw = block.Bytes
	}
	pubAny, err := x509.ParsePKIXPublicKey(raw)
	if err != nil {
		return nil, fmt.Errorf("updater: ecdsa public key parse: %w", err)
	}
	pub, ok := pubAny.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("updater: ecdsa public key has wrong type %T", pubAny)
	}
	return pub, nil
}

func splitECDSASig(sig []byte) (r, s *big.Int, err error) {
	// Raw r||s, 64 bytes for P-256.
	if len(sig) == 64 {
		r = new(big.Int).SetBytes(sig[:32])
		s = new(big.Int).SetBytes(sig[32:])
		return r, s, nil
	}
	// ASN.1 DER fallback. Reject signatures with trailing data — accepting
	// them is a path to signature-malleability bugs and they are never
	// produced by conforming signers.
	var seq struct{ R, S *big.Int }
	rest, err := asn1.Unmarshal(sig, &seq)
	if err != nil {
		return nil, nil, fmt.Errorf("updater: ecdsa signature format unrecognised: %w", err)
	}
	if len(rest) != 0 {
		return nil, nil, errors.New("updater: ecdsa signature has trailing data")
	}
	if seq.R == nil || seq.S == nil {
		return nil, nil, errors.New("updater: ecdsa signature missing r/s")
	}
	return seq.R, seq.S, nil
}
