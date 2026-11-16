package windowsign

import (
	"context"
	"crypto/x509/pkix"
	"encoding/asn1"
	"github.com/wailsapp/wails/v3/internal/codesignature"
)

type Identity = codesignature.Identity
type Options = codesignature.Options

func DecodePFX(data []byte, password string) (Identity, error) {
	return codesignature.DecodePFX(data, password)
}

var (
	oidIndirectData = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 4}
	oidPEImage      = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 15}
	oidSHA256       = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}
)

func rawDER(value any) (asn1.RawValue, error) {
	data, err := asn1.Marshal(value)
	return asn1.RawValue{FullBytes: data}, err
}
func explicitContent(data []byte) asn1.RawValue {
	return asn1.RawValue{Class: 2, Tag: 0, IsCompound: true, Bytes: data}
}
func signedContent(ctx context.Context, identity Identity, format asn1.ObjectIdentifier, formatData asn1.RawValue, digest []byte, options Options) ([]byte, error) {
	indirect := struct {
		Data struct {
			Type  asn1.ObjectIdentifier
			Value asn1.RawValue
		}
		Digest struct {
			Algorithm pkix.AlgorithmIdentifier
			Digest    []byte
		}
	}{}
	indirect.Data.Type = format
	indirect.Data.Value = formatData
	indirect.Digest.Algorithm = pkix.AlgorithmIdentifier{Algorithm: oidSHA256, Parameters: asn1.RawValue{Tag: asn1.TagNull}}
	indirect.Digest.Digest = digest
	content, err := asn1.Marshal(indirect)
	if err != nil {
		return nil, err
	}
	opus, err := rawDER(struct{}{})
	if err != nil {
		return nil, err
	}
	statement, err := rawDER([]asn1.ObjectIdentifier{{1, 3, 6, 1, 4, 1, 311, 2, 1, 21}})
	if err != nil {
		return nil, err
	}
	attrs := []codesignature.Attribute{{Type: asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 12}, Values: []asn1.RawValue{opus}}, {Type: asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 11}, Values: []asn1.RawValue{statement}}}
	options.TimestampAttribute = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 3, 3, 1}
	return codesignature.SignPKCS7(ctx, identity, oidIndirectData, content, attrs, options)
}
func peFormatData() (asn1.RawValue, error) {
	text := []byte{}
	for _, r := range "<<<Obsolete>>>" {
		text = append(text, byte(r>>8), byte(r))
	}
	stringData, err := asn1.Marshal(asn1.RawValue{Class: 2, Tag: 0, Bytes: text})
	if err != nil {
		return asn1.RawValue{}, err
	}
	link, err := asn1.Marshal(asn1.RawValue{Class: 2, Tag: 2, IsCompound: true, Bytes: stringData})
	if err != nil {
		return asn1.RawValue{}, err
	}
	image := struct {
		Flags asn1.BitString
		File  asn1.RawValue
	}{asn1.BitString{}, explicitContent(link)}
	encoded, err := asn1.Marshal(image)
	if err != nil {
		return asn1.RawValue{}, err
	}
	// Microsoft's Authenticode profile encodes the PE image value directly.
	return asn1.RawValue{FullBytes: encoded}, nil
}
