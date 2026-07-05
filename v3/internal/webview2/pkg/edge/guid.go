//go:build windows

package edge

// This code has been adapted from: https://github.com/go-ole/go-ole

/*

The MIT License (MIT)

Copyright © 2013-2017 Yasuhiro Matsumoto, <mattn.jp@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the “Software”), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

const hextable = "0123456789ABCDEF"
const emptyGUID = "{00000000-0000-0000-0000-000000000000}"

// GUID is Windows API specific GUID type.
//
// This exists to match Windows GUID type for direct passing for COM.
// Format is in xxxxxxxx-xxxx-xxxx-xxxxxxxxxxxxxxxx.
type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

// NewGUID converts the given string into a globally unique identifier that is
// compliant with the Windows API.
//
// The supplied string may be in any of these formats:
//
//	XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
//	XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
//	{XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}
//
// The conversion of the supplied string is not case-sensitive.
func NewGUID(guid string) *GUID {
	d := []byte(guid)
	var d1, d2, d3, d4a, d4b []byte

	switch len(d) {
	case 38:
		if d[0] != '{' || d[37] != '}' {
			return nil
		}
		d = d[1:37]
		fallthrough
	case 36:
		if d[8] != '-' || d[13] != '-' || d[18] != '-' || d[23] != '-' {
			return nil
		}
		d1 = d[0:8]
		d2 = d[9:13]
		d3 = d[14:18]
		d4a = d[19:23]
		d4b = d[24:36]
	case 32:
		d1 = d[0:8]
		d2 = d[8:12]
		d3 = d[12:16]
		d4a = d[16:20]
		d4b = d[20:32]
	default:
		return nil
	}

	var g GUID
	var ok1, ok2, ok3, ok4 bool
	g.Data1, ok1 = decodeHexUint32(d1)
	g.Data2, ok2 = decodeHexUint16(d2)
	g.Data3, ok3 = decodeHexUint16(d3)
	g.Data4, ok4 = decodeHexByte64(d4a, d4b)
	if ok1 && ok2 && ok3 && ok4 {
		return &g
	}
	return nil
}

func decodeHexUint32(src []byte) (value uint32, ok bool) {
	var b1, b2, b3, b4 byte
	var ok1, ok2, ok3, ok4 bool
	b1, ok1 = decodeHexByte(src[0], src[1])
	b2, ok2 = decodeHexByte(src[2], src[3])
	b3, ok3 = decodeHexByte(src[4], src[5])
	b4, ok4 = decodeHexByte(src[6], src[7])
	value = (uint32(b1) << 24) | (uint32(b2) << 16) | (uint32(b3) << 8) | uint32(b4)
	ok = ok1 && ok2 && ok3 && ok4
	return
}

func decodeHexUint16(src []byte) (value uint16, ok bool) {
	var b1, b2 byte
	var ok1, ok2 bool
	b1, ok1 = decodeHexByte(src[0], src[1])
	b2, ok2 = decodeHexByte(src[2], src[3])
	value = (uint16(b1) << 8) | uint16(b2)
	ok = ok1 && ok2
	return
}

func decodeHexByte64(s1 []byte, s2 []byte) (value [8]byte, ok bool) {
	var ok1, ok2, ok3, ok4, ok5, ok6, ok7, ok8 bool
	value[0], ok1 = decodeHexByte(s1[0], s1[1])
	value[1], ok2 = decodeHexByte(s1[2], s1[3])
	value[2], ok3 = decodeHexByte(s2[0], s2[1])
	value[3], ok4 = decodeHexByte(s2[2], s2[3])
	value[4], ok5 = decodeHexByte(s2[4], s2[5])
	value[5], ok6 = decodeHexByte(s2[6], s2[7])
	value[6], ok7 = decodeHexByte(s2[8], s2[9])
	value[7], ok8 = decodeHexByte(s2[10], s2[11])
	ok = ok1 && ok2 && ok3 && ok4 && ok5 && ok6 && ok7 && ok8
	return
}

func decodeHexByte(c1, c2 byte) (value byte, ok bool) {
	var n1, n2 byte
	var ok1, ok2 bool
	n1, ok1 = decodeHexChar(c1)
	n2, ok2 = decodeHexChar(c2)
	value = (n1 << 4) | n2
	ok = ok1 && ok2
	return
}

func decodeHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}

// String converts the GUID to string form. It will adhere to this pattern:
//
//	{XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}
//
// If the GUID is nil, the string representation of an empty GUID is returned:
//
//	{00000000-0000-0000-0000-000000000000}
func (guid *GUID) String() string {
	if guid == nil {
		return emptyGUID
	}

	var c [38]byte
	c[0] = '{'
	putUint32Hex(c[1:9], guid.Data1)
	c[9] = '-'
	putUint16Hex(c[10:14], guid.Data2)
	c[14] = '-'
	putUint16Hex(c[15:19], guid.Data3)
	c[19] = '-'
	putByteHex(c[20:24], guid.Data4[0:2])
	c[24] = '-'
	putByteHex(c[25:37], guid.Data4[2:8])
	c[37] = '}'
	return string(c[:])
}

func putUint32Hex(b []byte, v uint32) {
	b[0] = hextable[byte(v>>24)>>4]
	b[1] = hextable[byte(v>>24)&0x0f]
	b[2] = hextable[byte(v>>16)>>4]
	b[3] = hextable[byte(v>>16)&0x0f]
	b[4] = hextable[byte(v>>8)>>4]
	b[5] = hextable[byte(v>>8)&0x0f]
	b[6] = hextable[byte(v)>>4]
	b[7] = hextable[byte(v)&0x0f]
}

func putUint16Hex(b []byte, v uint16) {
	b[0] = hextable[byte(v>>8)>>4]
	b[1] = hextable[byte(v>>8)&0x0f]
	b[2] = hextable[byte(v)>>4]
	b[3] = hextable[byte(v)&0x0f]
}

func putByteHex(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		dst[i*2] = hextable[src[i]>>4]
		dst[i*2+1] = hextable[src[i]&0x0f]
	}
}

// IsEqualGUID compares two GUID.
//
// Not constant time comparison.
func IsEqualGUID(guid1 *GUID, guid2 *GUID) bool {
	return guid1.Data1 == guid2.Data1 &&
		guid1.Data2 == guid2.Data2 &&
		guid1.Data3 == guid2.Data3 &&
		guid1.Data4[0] == guid2.Data4[0] &&
		guid1.Data4[1] == guid2.Data4[1] &&
		guid1.Data4[2] == guid2.Data4[2] &&
		guid1.Data4[3] == guid2.Data4[3] &&
		guid1.Data4[4] == guid2.Data4[4] &&
		guid1.Data4[5] == guid2.Data4[5] &&
		guid1.Data4[6] == guid2.Data4[6] &&
		guid1.Data4[7] == guid2.Data4[7]
}
