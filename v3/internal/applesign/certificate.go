package applesign

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/binary"
	"errors"
	"time"

	"github.com/wailsapp/wails/v3/internal/codesignature"
	"howett.net/plist"
)

var developerIDApplication = asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 6, 1, 13}

// SigningOptions constructs certificate-backed signing options, sharing the
// authenticated PFX/CMS/timestamp implementation with Windows signing. Set the
// returned Identifier for files or DMGs; bundle identifiers come from Info.plist.
// Timestamping uses only the explicitly configured server and trust roots.
// Notarized Developer ID distribution requires a compatible secure timestamp;
// callers must enforce that policy before submitting an artifact.
func SigningOptions(identity codesignature.Identity, timestamp codesignature.Options) (Options, error) {
	if identity.Certificate == nil || identity.Signer == nil {
		return Options{}, errors.New("Apple signing requires a certificate and private signer")
	}
	o := Options{}
	developerID := false
	for _, extension := range identity.Certificate.Extensions {
		if extension.Id.Equal(developerIDApplication) {
			developerID = true
		}
	}
	if len(identity.Certificate.Subject.OrganizationalUnit) == 1 {
		o.TeamID = identity.Certificate.Subject.OrganizationalUnit[0]
	}
	if developerID && o.TeamID == "" {
		return Options{}, errors.New("Developer ID certificate must identify one signing team")
	}
	// SHA-1 is required by the certificate-hash requirement operand. It is not
	// used for executable pages, resource seals, or the CMS signature.
	leafHash := sha1.Sum(identity.Certificate.Raw)
	o.requirement = func(identifier string) ([]byte, error) {
		return designatedRequirement(identifier, o.TeamID, developerID, leafHash[:])
	}
	o.SignCMS = func(ctx context.Context, directory []byte) ([]byte, error) {
		sum := sha256.Sum256(directory)
		legacy, err := plist.Marshal(map[string]any{"cdhashes": [][]byte{sum[:20]}}, plist.XMLFormat)
		if err != nil {
			return nil, err
		}
		v1, err := asn1.Marshal(legacy)
		if err != nil {
			return nil, err
		}
		v2, err := asn1.Marshal(struct {
			Algorithm asn1.ObjectIdentifier
			Digest    []byte
		}{asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}, sum[:]})
		if err != nil {
			return nil, err
		}
		now, err := asn1.Marshal(time.Now().UTC())
		if err != nil {
			return nil, err
		}
		attrs := []codesignature.Attribute{
			{Type: asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 9, 1}, Values: []asn1.RawValue{{FullBytes: v1}}},
			{Type: asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 9, 2}, Values: []asn1.RawValue{{FullBytes: v2}}},
			{Type: asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 5}, Values: []asn1.RawValue{{FullBytes: now}}},
		}
		return codesignature.SignDetached(ctx, identity, directory, attrs, timestamp)
	}
	return o, nil
}

func word(value uint32) []byte { return binary.BigEndian.AppendUint32(nil, value) }
func requirementData(value []byte) []byte {
	result := binary.BigEndian.AppendUint32(nil, uint32(len(value)))
	result = append(result, value...)
	return append(result, make([]byte, (-len(value))&3)...)
}
func conjunction(left, right []byte) []byte { return append(append(word(6), left...), right...) }

// Requirements are prefix expressions with network-order words and padded
// byte strings, as specified by Apple's public requirement format declarations.
func designatedRequirement(identifier, team string, developerID bool, leafHash []byte) ([]byte, error) {
	if err := validateOptions(Options{Identifier: identifier, TeamID: team}); err != nil {
		return nil, err
	}
	if !developerID && len(leafHash) != sha1.Size {
		return nil, errors.New("certificate requirement needs a SHA-1 certificate hash")
	}
	expr := append(word(2), requirementData([]byte(identifier))...)
	if developerID {
		expr = conjunction(expr, word(15)) // generic Apple trust anchor
		for _, cert := range []struct {
			index uint32
			oid   asn1.ObjectIdentifier
		}{{1, asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 6, 2, 6}}, {0, developerIDApplication}} {
			encoded, _ := asn1.Marshal(cert.oid)
			var raw asn1.RawValue
			_, _ = asn1.Unmarshal(encoded, &raw)
			operand := append(word(14), word(cert.index)...)
			operand = append(operand, requirementData(raw.Bytes)...)
			operand = append(operand, word(0)...)
			expr = conjunction(expr, operand)
		}
		field := append(word(11), word(0)...)
		field = append(field, requirementData([]byte("subject.OU"))...)
		field = append(field, word(1)...)
		field = append(field, requirementData([]byte(team))...)
		expr = conjunction(expr, field)
	} else {
		anchor := append(word(4), word(0)...)
		anchor = append(anchor, requirementData(leafHash)...)
		expr = conjunction(expr, anchor)
	}
	requirement, err := wrap(0xfade0c00, append(word(1), expr...))
	if err != nil {
		return nil, err
	}
	return superblob(0xfade0c01, []component{{3, requirement}})
}
