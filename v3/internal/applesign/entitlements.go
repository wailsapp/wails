package applesign

import (
	"encoding/asn1"
	"errors"
	"fmt"
	"math"
	"slices"

	"howett.net/plist"
)

// Modern Apple entitlement DER uses an application-16 envelope containing a
// version integer and a context-16 dictionary. Dictionaries order their UTF-8
// keys lexically; arrays preserve source order. Native codesign fixtures cover
// these format details independently of this encoder.
func derEntitlements(xml []byte) ([]byte, error) {
	var value map[string]any
	if _, err := plist.Unmarshal(xml, &value); err != nil {
		return nil, fmt.Errorf("parse entitlements: %w", err)
	}
	if value == nil {
		return nil, errors.New("entitlements must be a dictionary")
	}
	nodes := 0
	dictionary, err := entitlementValue(value, 0, &nodes)
	if err != nil {
		return nil, err
	}
	return derContainer(0x70, append([]byte{2, 1, 1}, dictionary...))
}

func derContainer(tag byte, contents []byte) ([]byte, error) {
	return asn1.Marshal(asn1.RawValue{Class: int(tag >> 6), Tag: int(tag & 31), IsCompound: tag&32 != 0, Bytes: contents})
}

func entitlementValue(value any, depth int, nodes *int) ([]byte, error) {
	*nodes++
	if depth > 32 || *nodes > 100000 {
		return nil, errors.New("entitlement structure exceeds limits")
	}
	switch v := value.(type) {
	case string:
		return asn1.MarshalWithParams(v, "utf8")
	case bool, int64:
		return asn1.Marshal(v)
	case uint64:
		if v > math.MaxInt64 {
			return nil, errors.New("entitlement integer exceeds signed 64-bit range")
		}
		return asn1.Marshal(int64(v))
	case []any:
		var body []byte
		for _, child := range v {
			data, err := entitlementValue(child, depth+1, nodes)
			if err != nil {
				return nil, err
			}
			body = append(body, data...)
		}
		return derContainer(0x30, body)
	case map[string]any:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		slices.Sort(keys)
		var body []byte
		for _, key := range keys {
			name, err := asn1.MarshalWithParams(key, "utf8")
			if err != nil {
				return nil, err
			}
			data, err := entitlementValue(v[key], depth+1, nodes)
			if err != nil {
				return nil, err
			}
			entry, err := derContainer(0x30, append(name, data...))
			if err != nil {
				return nil, err
			}
			body = append(body, entry...)
		}
		return derContainer(0xb0, body)
	default:
		return nil, fmt.Errorf("unsupported entitlement value %T", value)
	}
}
