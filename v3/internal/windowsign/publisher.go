package windowsign

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Windows renders distinguished names with reversed RDN order, optional spaces
// after separators, S for state, and quoted values. Compare the original
// certificate RDNs instead of rebuilding/reordering them through pkix.Name.
func publisherMatches(publisher string, certificate *x509.Certificate) (bool, error) {
	if certificate == nil || len(publisher) > 16<<10 || !utf8.ValidString(publisher) {
		return false, errors.New("invalid MSIX publisher")
	}
	var expected pkix.RDNSequence
	rest, err := asn1.Unmarshal(certificate.RawSubject, &expected)
	if err != nil || len(rest) != 0 {
		return false, errors.New("invalid signing certificate subject")
	}
	names, err := splitDN(publisher, ',')
	if err != nil {
		return false, err
	}
	if len(names) != len(expected) {
		return false, nil
	}
	for i, name := range names {
		attributes, err := splitDN(name, '+')
		if err != nil {
			return false, err
		}
		rdn := expected[len(expected)-1-i]
		if len(attributes) != len(rdn) {
			return false, nil
		}
		used := make([]bool, len(rdn))
		for _, attribute := range attributes {
			key, value, ok := strings.Cut(attribute, "=")
			if !ok {
				return false, errors.New("publisher attribute requires a value")
			}
			oid, err := publisherOID(strings.TrimSpace(key))
			if err != nil {
				return false, err
			}
			value, err = publisherValue(strings.TrimSpace(value))
			if err != nil {
				return false, err
			}
			matched := false
			for j, candidate := range rdn {
				if used[j] || !candidate.Type.Equal(oid) {
					continue
				}
				text, ok := candidate.Value.(string)
				if ok && text == value {
					used[j] = true
					matched = true
					break
				}
			}
			if !matched {
				return false, nil
			}
		}
	}
	return true, nil
}
func splitDN(value string, separator byte) ([]string, error) {
	var result []string
	start := 0
	quoted := false
	for i := 0; i < len(value); i++ {
		switch value[i] {
		case '\\':
			i++
			if i >= len(value) {
				return nil, errors.New("truncated publisher escape")
			}
		case '"':
			if quoted && i+1 < len(value) && value[i+1] == '"' {
				i++
				continue
			}
			quoted = !quoted
		default:
			if value[i] == separator && !quoted {
				part := strings.TrimSpace(value[start:i])
				if part == "" {
					return nil, errors.New("empty publisher attribute")
				}
				result = append(result, part)
				start = i + 1
			}
		}
	}
	if quoted {
		return nil, errors.New("unterminated publisher quote")
	}
	part := strings.TrimSpace(value[start:])
	if part == "" {
		return nil, errors.New("empty publisher attribute")
	}
	return append(result, part), nil
}
func publisherValue(value string) (string, error) {
	quoted := strings.HasPrefix(value, "\"")
	if quoted {
		if len(value) < 2 || value[len(value)-1] != '"' {
			return "", errors.New("invalid quoted publisher value")
		}
		value = value[1 : len(value)-1]
	}
	var result strings.Builder
	for i := 0; i < len(value); i++ {
		switch value[i] {
		case '\\':
			i++
			if i >= len(value) {
				return "", errors.New("truncated publisher escape")
			}
			if i+1 < len(value) {
				if b, err := strconv.ParseUint(value[i:i+2], 16, 8); err == nil {
					result.WriteByte(byte(b))
					i++
					continue
				}
			}
			result.WriteByte(value[i])
		case '"':
			if !quoted || i+1 >= len(value) || value[i+1] != '"' {
				return "", errors.New("invalid publisher quote")
			}
			result.WriteByte('"')
			i++
		default:
			result.WriteByte(value[i])
		}
	}
	if !utf8.ValidString(result.String()) {
		return "", errors.New("invalid UTF-8 publisher value")
	}
	return result.String(), nil
}
func publisherOID(key string) (asn1.ObjectIdentifier, error) {
	aliases := map[string]string{"CN": "2.5.4.3", "C": "2.5.4.6", "L": "2.5.4.7", "S": "2.5.4.8", "ST": "2.5.4.8", "STREET": "2.5.4.9", "O": "2.5.4.10", "OU": "2.5.4.11", "T": "2.5.4.12", "TITLE": "2.5.4.12", "SERIALNUMBER": "2.5.4.5", "SN": "2.5.4.4", "SURNAME": "2.5.4.4", "G": "2.5.4.42", "GN": "2.5.4.42", "GIVENNAME": "2.5.4.42", "I": "2.5.4.43", "INITIALS": "2.5.4.43", "E": "1.2.840.113549.1.9.1", "EMAILADDRESS": "1.2.840.113549.1.9.1", "DC": "0.9.2342.19200300.100.1.25", "UID": "0.9.2342.19200300.100.1.1"}
	oid := aliases[strings.ToUpper(key)]
	if oid == "" {
		oid = strings.TrimPrefix(strings.ToUpper(key), "OID.")
	}
	parts := strings.Split(oid, ".")
	if len(parts) < 2 || len(parts) > 32 {
		return nil, errors.New("unknown publisher attribute")
	}
	result := make(asn1.ObjectIdentifier, len(parts))
	for i, p := range parts {
		n, err := strconv.ParseUint(p, 10, 31)
		if err != nil {
			return nil, errors.New("invalid publisher attribute OID")
		}
		result[i] = int(n)
	}
	if result[0] > 2 || (result[0] < 2 && result[1] > 39) {
		return nil, errors.New("invalid publisher attribute OID")
	}
	return result, nil
}
