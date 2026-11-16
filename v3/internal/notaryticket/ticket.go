// Package notaryticket verifies and attaches Apple's observed notarization
// ticket format. This is not a supported public Apple service contract.
package notaryticket

import (
	"bytes"
	"crypto/x509"
	_ "embed"
	"encoding/asn1"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

const maxTicketSize = 4 << 20

// Apple publishes this public trust anchor at its certificate authority site.
//
//go:embed apple-root-g3.der
var appleRootDER []byte

// Digest identifies a code directory. This implementation accepts SHA256,
// represented by algorithm 2 and its first 20 digest bytes, only.
type Digest struct {
	Algorithm uint8
	Value     [20]byte
}

// Ticket holds an authenticated, positive ticket; its signed bytes are private
// so callers cannot modify them after verification.
type Ticket struct {
	data   []byte
	hashes map[Digest]bool
}

func (t *Ticket) Bytes() []byte {
	if t == nil {
		return nil
	}
	return bytes.Clone(t.data)
}

// Verify establishes Apple ticket-signing trust and association with all the
// expected code directories. It does not assess executable integrity or query
// Gatekeeper's revocation database.
func Verify(data []byte, expected []Digest) (*Ticket, error) {
	if err := validateDigests(expected); err != nil {
		return nil, err
	}
	if len(data) < 16 || len(data) > maxTicketSize || string(data[:4]) != "s8ch" || binary.LittleEndian.Uint32(data[4:]) != 1 {
		return nil, errors.New("unsupported notarization ticket envelope")
	}
	certSize, payloadSize := uint64(binary.LittleEndian.Uint32(data[8:])), uint64(binary.LittleEndian.Uint32(data[12:]))
	end := uint64(16) + certSize + payloadSize
	if certSize == 0 || payloadSize < 24 || end >= uint64(len(data)) || uint64(len(data))-end > 144 {
		return nil, errors.New("invalid notarization ticket bounds")
	}
	payload := data[16+certSize : end]
	hashes, issued, err := ticketHashes(payload)
	if err != nil {
		return nil, err
	}
	if err := verifyAppleSignature(data[16:16+certSize], data[:end], data[end:], issued); err != nil {
		return nil, err
	}
	ticket := &Ticket{data: bytes.Clone(data), hashes: hashes}
	if err := ticket.includes(expected); err != nil {
		return nil, err
	}
	return ticket, nil
}

func validateDigests(digests []Digest) error {
	if len(digests) == 0 || len(digests) > 64 {
		return errors.New("ticket lookup requires 1–64 code directory hashes")
	}
	for _, digest := range digests {
		if digest.Algorithm != 2 || digest.Value == ([20]byte{}) {
			return errors.New("ticket lookup requires nonzero SHA256 code directory hashes")
		}
	}
	return nil
}

func (t *Ticket) includes(expected []Digest) error {
	if err := validateDigests(expected); err != nil {
		return err
	}
	if t == nil {
		return errors.New("missing verified notarization ticket")
	}
	for _, hash := range expected {
		if !t.hashes[hash] {
			return fmt.Errorf("notarization ticket does not cover code directory %x", hash.Value)
		}
	}
	return nil
}

func ticketHashes(payload []byte) (map[Digest]bool, time.Time, error) {
	le := binary.LittleEndian
	if len(payload) < 24 || string(payload[:4]) != "g8tk" || le.Uint16(payload[4:]) != 2 || le.Uint16(payload[6:]) != 20 || le.Uint32(payload[12:]) != 0 {
		return nil, time.Time{}, errors.New("unsupported or non-positive notarization ticket payload")
	}
	count := uint64(le.Uint32(payload[8:]))
	if count == 0 || count > 100_000 || 24+21*count != uint64(len(payload)) {
		return nil, time.Time{}, errors.New("invalid notarization hash list bounds")
	}
	seconds := le.Uint64(payload[16:])
	if seconds < 946684800 || seconds > uint64(time.Now().Add(5*time.Minute).Unix()) {
		return nil, time.Time{}, errors.New("invalid notarization issuance time")
	}
	hashes := make(map[Digest]bool, count)
	for offset := 24; offset < len(payload); offset += 21 {
		entry := Digest{Algorithm: payload[offset]}
		copy(entry.Value[:], payload[offset+1:offset+21])
		if entry.Algorithm != 2 {
			return nil, time.Time{}, errors.New("unsupported notarization hash algorithm")
		}
		hashes[entry] = true
	}
	return hashes, time.Unix(int64(seconds), 0), nil
}

func verifyAppleSignature(encoded, signed, signature []byte, issued time.Time) error {
	// Observed Apple envelopes pad the DER signature to a four-byte boundary.
	// Padding is outside the signature; only zero padding of the exact size is allowed.
	var signatureValue asn1.RawValue
	padding, err := asn1.Unmarshal(signature, &signatureValue)
	if err != nil || len(signature)%4 != 0 || len(padding) > 3 || !bytes.Equal(padding, make([]byte, len(padding))) || signatureValue.Tag != asn1.TagSequence || signatureValue.Class != 0 {
		return errors.New("invalid ticket signature encoding")
	}
	signature = signatureValue.FullBytes
	var sequence asn1.RawValue
	rest, err := asn1.Unmarshal(encoded, &sequence)
	if err != nil || len(rest) != 0 || sequence.Class != 0 || sequence.Tag != asn1.TagSequence || !sequence.IsCompound {
		return errors.New("invalid ticket certificate sequence")
	}
	var certificates []*x509.Certificate
	for data := sequence.Bytes; len(data) > 0; {
		if len(certificates) == 8 {
			return errors.New("too many ticket certificates")
		}
		var value asn1.RawValue
		rest, err := asn1.Unmarshal(data, &value)
		if err != nil {
			return err
		}
		certificate, err := x509.ParseCertificate(value.FullBytes)
		if err != nil {
			return err
		}
		certificates = append(certificates, certificate)
		data = rest
	}
	if len(certificates) == 0 {
		return errors.New("missing ticket signing certificate")
	}
	leaf := certificates[0]
	authorized := false
	for _, extension := range leaf.Extensions {
		if extension.Id.Equal(asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 6, 1, 30}) && bytes.Equal(extension.Value, []byte{5, 0}) {
			authorized = true
		}
	}
	if !authorized || leaf.IsCA || leaf.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
		return errors.New("certificate is not authorized for Apple software tickets")
	}
	root, err := x509.ParseCertificate(appleRootDER)
	if err != nil {
		return err
	}
	roots, intermediates := x509.NewCertPool(), x509.NewCertPool()
	roots.AddCert(root)
	for _, certificate := range certificates[1:] {
		intermediates.AddCert(certificate)
	}
	if _, err := leaf.Verify(x509.VerifyOptions{Roots: roots, Intermediates: intermediates, CurrentTime: issued, KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny}}); err != nil {
		return fmt.Errorf("untrusted Apple ticket certificate: %w", err)
	}
	if err := leaf.CheckSignature(x509.ECDSAWithSHA256, signed, signature); err != nil {
		return fmt.Errorf("invalid notarization ticket signature: %w", err)
	}
	return nil
}
