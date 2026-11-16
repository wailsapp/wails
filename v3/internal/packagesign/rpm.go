package packagesign

import (
	"bytes"
	"context"
	"crypto/md5"  // RPM v4's signature header retains its legacy package digest.
	"crypto/sha1" // RPM v4's signature header retains its legacy header digest.
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

var rpmHeaderMagic = []byte{0x8e, 0xad, 0xe8, 1, 0, 0, 0, 0}

type rpmEntry struct {
	tag, kind, count uint32
	data             []byte
}

// RPM v4 envelope layout: https://rpm.org/docs/latest/manual/format_v4.html
// Only the signature header is rewritten; the immutable header and compressed
// payload are neither reconstructed nor recompressed.
func signRPM(ctx context.Context, src *os.File, output string, signer *signingKey) error {
	headerTag, payloadTag, err := rpmSignatureTags(signer.key)
	if err != nil {
		return err
	}
	lead := make([]byte, 96)
	if _, err := io.ReadFull(src, lead); err != nil || !bytes.Equal(lead[:4], []byte{0xed, 0xab, 0xee, 0xdb}) || lead[4] != 3 || binary.BigEndian.Uint16(lead[78:80]) != 5 {
		return fmt.Errorf("unsupported RPM lead")
	}
	oldSignature, entries, err := readRPMHeader(src)
	if err != nil {
		return err
	}
	if _, err := src.Seek(int64((8-len(oldSignature)%8)%8), io.SeekCurrent); err != nil {
		return err
	}
	start, err := src.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	header, metadata, err := readRPMHeader(src)
	if err != nil {
		return err
	}
	if err := validateRPMRegion(header, 63); err != nil {
		return err
	}
	var headerSignature, payloadSignature bytes.Buffer
	if err := openpgp.DetachSign(&headerSignature, signer.entity, bytes.NewReader(header), signer.config); err != nil {
		return err
	}
	original, md, payloadHash := sha256.New(), md5.New(), sha256.New()
	content := io.MultiReader(bytes.NewReader(header), io.TeeReader(contextReader{ctx, src}, payloadHash))
	message := io.TeeReader(content, io.MultiWriter(original, md))
	if err := openpgp.DetachSign(&payloadSignature, signer.entity, message, signer.config); err != nil {
		return err
	}
	sha1Header, sha256Header := sha1.Sum(header), sha256.Sum256(header)
	if err := checkRPMDigests(entries, md.Sum(nil), sha1Header[:], sha256Header[:]); err != nil {
		return err
	}
	if err := checkRPMPayloadDigest(metadata, payloadHash.Sum(nil)); err != nil {
		return err
	}
	retained := make([]rpmEntry, 0, len(entries)+2)
	for _, entry := range entries {
		switch entry.tag {
		case 62, 267, 268, 1002, 1005, 1008:
			// Replace package signatures and their optional reserved space.
		default:
			retained = append(retained, entry)
		}
	}
	retained = append(retained,
		rpmEntry{headerTag, 7, uint32(headerSignature.Len()), headerSignature.Bytes()},
		rpmEntry{payloadTag, 7, uint32(payloadSignature.Len()), payloadSignature.Bytes()})
	signature := encodeRPMHeader(retained)
	if _, err := src.Seek(start, io.SeekStart); err != nil {
		return err
	}
	return writeOutput(output, func(dst io.Writer) error {
		for _, part := range [][]byte{lead, signature, make([]byte, (8-len(signature)%8)%8)} {
			if _, err := dst.Write(part); err != nil {
				return err
			}
		}
		copied := sha256.New()
		if _, err := io.Copy(io.MultiWriter(dst, copied), contextReader{ctx, src}); err != nil {
			return err
		}
		if !bytes.Equal(original.Sum(nil), copied.Sum(nil)) {
			return fmt.Errorf("RPM input changed while signing")
		}
		return nil
	})
}

func rpmSignatureTags(key *packet.PrivateKey) (uint32, uint32, error) {
	if key.Version != 4 {
		return 0, 0, fmt.Errorf("RPM signing requires an OpenPGP v4 key")
	}
	switch key.PubKeyAlgo {
	case packet.PubKeyAlgoRSA, packet.PubKeyAlgoRSASignOnly:
		return 268, 1002, nil
	case packet.PubKeyAlgoDSA, packet.PubKeyAlgoECDSA, packet.PubKeyAlgoEdDSA:
		return 267, 1005, nil
	default:
		return 0, 0, fmt.Errorf("unsupported RPM signing key algorithm %d", key.PubKeyAlgo)
	}
}

func readRPMHeader(src io.Reader) ([]byte, []rpmEntry, error) {
	intro := make([]byte, 16)
	if _, err := io.ReadFull(src, intro); err != nil || !bytes.Equal(intro[:8], rpmHeaderMagic) {
		return nil, nil, fmt.Errorf("invalid RPM header")
	}
	n, size := binary.BigEndian.Uint32(intro[8:]), binary.BigEndian.Uint32(intro[12:])
	// Bound allocations before reading lengths from an untrusted package.
	if n == 0 || n > 65536 || size > 64<<20 {
		return nil, nil, fmt.Errorf("RPM header exceeds supported limits")
	}
	header := append(intro, make([]byte, int(n)*16+int(size))...)
	if _, err := io.ReadFull(src, header[16:]); err != nil {
		return nil, nil, fmt.Errorf("truncated RPM header: %w", err)
	}
	data := header[16+int(n)*16:]
	entries := make([]rpmEntry, 0, n)
	for i := range int(n) {
		index := header[16+i*16 : 32+i*16]
		tag, kind := binary.BigEndian.Uint32(index), binary.BigEndian.Uint32(index[4:])
		offset, count := binary.BigEndian.Uint32(index[8:]), binary.BigEndian.Uint32(index[12:])
		if i > 0 && tag <= entries[i-1].tag {
			return nil, nil, fmt.Errorf("RPM header tags are not strictly ordered")
		}
		value, err := rpmValue(data, kind, offset, count)
		if err != nil {
			return nil, nil, fmt.Errorf("RPM tag %d: %w", tag, err)
		}
		entries = append(entries, rpmEntry{tag, kind, count, value})
	}
	return header, entries, nil
}

func rpmValue(data []byte, kind, offset, count uint32) ([]byte, error) {
	if uint64(offset) > uint64(len(data)) {
		return nil, fmt.Errorf("value offset out of bounds")
	}
	tail := data[offset:]
	if kind == 6 || kind == 8 || kind == 9 {
		if count == 0 || (kind == 6 && count != 1) || uint64(count) > uint64(len(tail)) {
			return nil, fmt.Errorf("invalid string count")
		}
		end := 0
		for range count {
			zero := bytes.IndexByte(tail[end:], 0)
			if zero < 0 {
				return nil, fmt.Errorf("unterminated string")
			}
			end += zero + 1
		}
		return tail[:end], nil
	}
	width := uint64(1)
	switch kind {
	case 1, 2, 7:
	case 3, 4, 5:
		width = 1 << (kind - 2)
	default:
		return nil, fmt.Errorf("unsupported value type %d", kind)
	}
	length := uint64(count) * width
	if uint64(offset)%width != 0 || length > uint64(len(tail)) {
		return nil, fmt.Errorf("invalid value length or alignment")
	}
	return tail[:length], nil
}

func validateRPMRegion(header []byte, tag uint32) error {
	n := binary.BigEndian.Uint32(header[8:])
	first := header[16:32]
	data := header[16+int(n)*16:]
	offset := binary.BigEndian.Uint32(first[8:])
	if binary.BigEndian.Uint32(first) != tag || binary.BigEndian.Uint32(first[4:]) != 7 || binary.BigEndian.Uint32(first[12:]) != 16 || uint64(offset)+16 != uint64(len(data)) {
		return fmt.Errorf("RPM header must contain a complete immutable region")
	}
	trailer := data[offset:]
	if binary.BigEndian.Uint32(trailer) != tag || binary.BigEndian.Uint32(trailer[4:]) != 7 || int32(binary.BigEndian.Uint32(trailer[8:])) != -int32(n*16) || binary.BigEndian.Uint32(trailer[12:]) != 16 {
		return fmt.Errorf("invalid RPM immutable region trailer")
	}
	return nil
}

func checkRPMDigests(entries []rpmEntry, md, sha1sum, sha256sum []byte) error {
	for _, entry := range entries {
		var expected []byte
		switch entry.tag {
		case 1004:
			expected = md
		case 269:
			expected = append([]byte(hex.EncodeToString(sha1sum)), 0)
		case 273:
			expected = append([]byte(hex.EncodeToString(sha256sum)), 0)
		default:
			continue
		}
		if !bytes.Equal(entry.data, expected) {
			return fmt.Errorf("RPM digest mismatch for tag %d", entry.tag)
		}
	}
	return nil
}

func checkRPMPayloadDigest(entries []rpmEntry, digest []byte) error {
	for _, entry := range entries {
		switch entry.tag {
		case 5114:
			if entry.kind != 4 || entry.count != 1 || binary.BigEndian.Uint32(entry.data) != 4 {
				return fmt.Errorf("only the RPM v4 package format is supported")
			}
		case 5092:
			want := append([]byte(hex.EncodeToString(digest)), 0)
			if entry.count != 1 || !bytes.Equal(entry.data, want) {
				return fmt.Errorf("RPM payload SHA256 digest mismatch")
			}
		case 5093:
			if entry.kind != 4 || entry.count != 1 || binary.BigEndian.Uint32(entry.data) != 8 {
				return fmt.Errorf("unsupported RPM payload digest algorithm")
			}
		}
	}
	return nil
}

func encodeRPMHeader(entries []rpmEntry) []byte {
	entries = append(entries, rpmEntry{tag: 62, kind: 7, count: 16})
	sort.Slice(entries, func(i, j int) bool { return entries[i].tag < entries[j].tag })
	header := make([]byte, 16+16*len(entries))
	copy(header, rpmHeaderMagic)
	binary.BigEndian.PutUint32(header[8:], uint32(len(entries)))
	var data []byte
	for i, entry := range entries {
		if entry.tag == 62 {
			continue
		}
		alignment := 1
		if entry.kind >= 3 && entry.kind <= 5 {
			alignment = 1 << (entry.kind - 2)
		}
		data = append(data, make([]byte, (alignment-len(data)%alignment)%alignment)...)
		putRPMIndex(header[16+i*16:], entry, uint32(len(data)))
		data = append(data, entry.data...)
	}
	region := rpmEntry{tag: 62, kind: 7, count: 16}
	putRPMIndex(header[16:], region, uint32(len(data)))
	trailer := make([]byte, 16)
	putRPMIndex(trailer, region, uint32(-int32(len(entries)*16)))
	data = append(data, trailer...)
	binary.BigEndian.PutUint32(header[12:], uint32(len(data)))
	return append(header, data...)
}

func putRPMIndex(dst []byte, entry rpmEntry, offset uint32) {
	binary.BigEndian.PutUint32(dst, entry.tag)
	binary.BigEndian.PutUint32(dst[4:], entry.kind)
	binary.BigEndian.PutUint32(dst[8:], offset)
	binary.BigEndian.PutUint32(dst[12:], entry.count)
}
