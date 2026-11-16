package windowsign

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"context"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf16"
)

const signatureName = "AppxSignature.p7x"

var oidSIPInfo = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 30}

// SignMSIX signs an individual MSIX/APPX package, replacing an existing
// signature. Bundles and ZIP entries requiring 64-bit sizes are not supported.
// Payload compression and the block map are preserved.
func SignMSIX(ctx context.Context, input, output string, identity Identity, options Options) error {
	if identity.Certificate == nil || identity.Signer == nil {
		return errors.New("MSIX signing requires a certificate and private key")
	}
	source, info, err := snapshotInput(ctx, input, output)
	if err != nil {
		return err
	}
	defer os.Remove(source.Name())
	defer source.Close()
	archive, err := zip.NewReader(source, info.Size())
	if err != nil {
		return err
	}
	if len(archive.File) > 65534 || archive.Comment != "" {
		return errors.New("MSIX signing requires fewer than 65535 entries and no ZIP comment")
	}
	entries := make(map[string]*zip.File, len(archive.File))
	seen := map[string]bool{}
	for _, file := range archive.File {
		logicalName, err := url.PathUnescape(file.Name)
		if err != nil {
			return err
		}
		key := strings.ToLower(logicalName)
		if seen[key] || strings.Contains(file.Name, "\\") || !filepath.IsLocal(filepath.FromSlash(logicalName)) || strings.ContainsAny(logicalName, "\\:\x00") || path.Clean(logicalName) != logicalName || file.FileInfo().IsDir() {
			return fmt.Errorf("invalid or duplicate MSIX entry %q", file.Name)
		}
		if file.UncompressedSize64 >= 1<<32 || file.CompressedSize64 >= 1<<32 || (file.Method != zip.Store && file.Method != zip.Deflate) || file.Flags&1 != 0 {
			return fmt.Errorf("unsupported MSIX ZIP entry %q", file.Name)
		}
		seen[key] = true
		if (key == "appxmanifest.xml" || key == "appxblockmap.xml" || key == "[content_types].xml" || key == "appxsignature.p7x") && file.Name != logicalName {
			return errors.New("MSIX footprint names must not be escaped")
		}
		entries[logicalName] = file
	}
	if entries["AppxMetadata/AppxBundleManifest.xml"] != nil {
		return errors.New("MSIX bundles are not supported by portable signing")
	}
	manifest, err := readZIPMetadata(entries["AppxManifest.xml"])
	if err != nil {
		return err
	}
	var app struct {
		XMLName  xml.Name
		Identity struct {
			Publisher string `xml:"Publisher,attr"`
		}
	}
	if err := decodeXML(manifest, &app); err != nil {
		return err
	}
	matches, err := publisherMatches(app.Identity.Publisher, identity.Certificate)
	if err != nil {
		return err
	}
	if app.XMLName.Local != "Package" || !matches {
		return fmt.Errorf("MSIX publisher must match signing certificate subject %q", identity.Certificate.Subject.String())
	}
	blockMap, err := readZIPMetadata(entries["AppxBlockMap.xml"])
	if err != nil {
		return err
	}
	records, err := packageRecords(source, info.Size())
	if err != nil {
		return err
	}
	headers := map[string]int{}
	for _, record := range records {
		name, err := url.PathUnescape(record.name)
		if err != nil {
			return err
		}
		headers[name] = record.header
	}
	if err := verifyBlockMap(ctx, blockMap, entries, headers); err != nil {
		return err
	}
	contentTypes, err := readZIPMetadata(entries["[Content_Types].xml"])
	if err != nil {
		return err
	}
	originalContentTypes := bytes.Clone(contentTypes)
	contentTypes, err = signedContentTypes(contentTypes)
	if err != nil {
		return err
	}
	// Building one unsigned ZIP first gives the precise pre-signature file-record
	// and central-directory bytes required by the APPX SIP digest format.
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	stage, err := os.CreateTemp(filepath.Dir(output), ".wake-msix-sign-*")
	if err != nil {
		return err
	}
	defer os.Remove(stage.Name())
	defer stage.Close()
	directory, central, err := unsignedMSIX(ctx, stage, source, records, contentTypes, !bytes.Equal(contentTypes, originalContentTypes))
	if err != nil {
		return err
	}
	size, err := stage.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	bag := []byte("APPX")
	for _, part := range []struct {
		name   string
		reader io.Reader
	}{
		{"AXPC", io.NewSectionReader(stage, 0, directory)},
		{"AXCD", io.NewSectionReader(stage, directory, size-directory)},
		{"AXCT", bytes.NewReader(contentTypes)},
		{"AXBM", bytes.NewReader(blockMap)},
	} {
		digest, err := hashReader(ctx, part.reader)
		if err != nil {
			return err
		}
		bag = append(bag, part.name...)
		bag = append(bag, digest...)
	}
	if catalog := entries["AppxMetadata/CodeIntegrity.cat"]; catalog != nil {
		reader, err := catalog.Open()
		if err != nil {
			return err
		}
		digest, hashErr := hashReader(ctx, reader)
		closeErr := reader.Close()
		if hashErr != nil {
			return hashErr
		}
		if closeErr != nil {
			return closeErr
		}
		bag = append(bag, "AXCI"...)
		bag = append(bag, digest...)
	}
	sip, err := rawDER(struct {
		Version                                               int
		ID                                                    []byte
		Reserved1, Reserved2, Reserved3, Reserved4, Reserved5 int
	}{0x01010000, []byte{0x4b, 0xdf, 0xc5, 0x0a, 0x07, 0xce, 0xe2, 0x4d, 0xb7, 0x6e, 0x23, 0xc8, 0x39, 0xa0, 0x9f, 0xd1}, 0, 0, 0, 0, 0})
	if err != nil {
		return err
	}
	signature, err := signedContent(ctx, identity, oidSIPInfo, sip, bag, options)
	if err != nil {
		return err
	}
	local, signatureDirectory, err := footprintRecord(signatureName, append([]byte("PKCX"), signature...))
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint32(signatureDirectory[42:], uint32(directory))
	finalDirectory := append(bytes.Clone(central), signatureDirectory...)
	count := len(archive.File) + 1
	if entries[signatureName] != nil {
		count--
	}
	return publishSigned(ctx, output, info.Mode().Perm(), func(target *os.File) error {
		if _, err := io.Copy(target, contextReader{ctx, io.NewSectionReader(stage, 0, directory)}); err != nil {
			return err
		}
		if _, err := target.Write(local); err != nil {
			return err
		}
		return finishMSIX(target, directory+int64(len(local)), finalDirectory, count)
	})
}

func hashReader(ctx context.Context, reader io.Reader) ([]byte, error) {
	h := sha256.New()
	_, err := io.CopyBuffer(h, contextReader{ctx, reader}, make([]byte, 128<<10))
	return h.Sum(nil), err
}

func zipDirectory(file *os.File, size int64) (int64, error) {
	if size < 22 {
		return 0, errors.New("truncated MSIX ZIP")
	}
	end := make([]byte, 22)
	if _, err := file.ReadAt(end, size-22); err != nil {
		return 0, err
	}
	if binary.LittleEndian.Uint32(end) != 0x06054b50 || binary.LittleEndian.Uint16(end[20:]) != 0 {
		return 0, errors.New("invalid MSIX ZIP footer")
	}
	if binary.LittleEndian.Uint32(end[16:]) == 0xffffffff {
		header := make([]byte, 56)
		if size < 98 {
			return 0, errors.New("truncated ZIP64 footer")
		}
		if _, err := file.ReadAt(header, size-98); err != nil {
			return 0, err
		}
		if binary.LittleEndian.Uint32(header) != 0x06064b50 {
			return 0, errors.New("invalid ZIP64 footer")
		}
		offset := int64(binary.LittleEndian.Uint64(header[48:]))
		length := int64(binary.LittleEndian.Uint64(header[40:]))
		if offset < 0 || length < 0 || offset+length != size-98 {
			return 0, errors.New("invalid ZIP64 directory")
		}
		return offset, nil
	}
	offset := int64(binary.LittleEndian.Uint32(end[16:]))
	length := int64(binary.LittleEndian.Uint32(end[12:]))
	if offset+length != size-22 {
		return 0, errors.New("portable signing requires a non-ZIP64 output directory")
	}
	return offset, nil
}

func readZIPMetadata(file *zip.File) ([]byte, error) {
	if file == nil {
		return nil, errors.New("MSIX is missing required signing metadata")
	}
	if file.UncompressedSize64 > 16<<20 {
		return nil, errors.New("MSIX metadata exceeds 16 MiB")
	}
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	data, err := io.ReadAll(io.LimitReader(reader, (16<<20)+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 16<<20 {
		return nil, errors.New("MSIX metadata exceeds 16 MiB")
	}
	return data, nil
}

func xmlText(data []byte) ([]byte, error) {
	if len(data) >= 2 && (bytes.Equal(data[:2], []byte{0xff, 0xfe}) || bytes.Equal(data[:2], []byte{0xfe, 0xff})) {
		if len(data)%2 != 0 {
			return nil, errors.New("invalid UTF-16 XML")
		}
		var order binary.ByteOrder = binary.LittleEndian
		if data[0] == 0xfe {
			order = binary.BigEndian
		}
		units := make([]uint16, (len(data)-2)/2)
		for i := range units {
			units[i] = order.Uint16(data[2+i*2:])
		}
		return []byte(string(utf16.Decode(units))), nil
	}
	return bytes.TrimPrefix(data, []byte{0xef, 0xbb, 0xbf}), nil
}
func xmlDecoder(data []byte) (*xml.Decoder, error) {
	data, err := xmlText(data)
	if err != nil {
		return nil, err
	}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.EqualFold(charset, "UTF-16") {
			return input, nil
		}
		return nil, fmt.Errorf("unsupported XML encoding %q", charset)
	}
	return decoder, nil
}
func decodeXML(data []byte, value any) error {
	decoder, err := xmlDecoder(data)
	if err != nil {
		return err
	}
	if err := decoder.Decode(value); err != nil {
		return err
	}
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if text, ok := token.(xml.CharData); !ok || len(bytes.TrimSpace(text)) != 0 {
			return errors.New("trailing XML content")
		}
	}
}

func signedContentTypes(data []byte) ([]byte, error) {
	var types struct {
		XMLName  xml.Name
		Defaults []struct {
			Extension string `xml:"Extension,attr"`
			Type      string `xml:"ContentType,attr"`
		} `xml:"Default"`
		Overrides []struct {
			Name string `xml:"PartName,attr"`
			Type string `xml:"ContentType,attr"`
		} `xml:"Override"`
	}
	if err := decodeXML(data, &types); err != nil {
		return nil, err
	}
	if types.XMLName.Local != "Types" || types.XMLName.Space != "http://schemas.openxmlformats.org/package/2006/content-types" {
		return nil, errors.New("invalid MSIX content types")
	}
	found := false
	for i := range types.Overrides {
		if types.Overrides[i].Name == "/"+signatureName {
			if found {
				return nil, errors.New("duplicate MSIX signature content type")
			}
			if types.Overrides[i].Type != "application/vnd.ms-appx.signature" {
				return nil, errors.New("incorrect MSIX signature content type")
			}
			found = true
		}
	}
	if found {
		return data, nil
	}
	{
		types.Overrides = append(types.Overrides, struct {
			Name string `xml:"PartName,attr"`
			Type string `xml:"ContentType,attr"`
		}{"/" + signatureName, "application/vnd.ms-appx.signature"})
	}
	body, err := xml.Marshal(types)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), body...), nil
}

func verifyBlockMap(ctx context.Context, data []byte, entries map[string]*zip.File, headers map[string]int) error {
	var blockMap struct {
		XMLName xml.Name
		Hash    string       `xml:"HashMethod,attr"`
		Files   []mappedFile `xml:"File"`
	}

	if err := decodeXML(data, &blockMap); err != nil {
		return err
	}
	if blockMap.XMLName.Local != "BlockMap" || blockMap.Hash != "http://www.w3.org/2001/04/xmlenc#sha256" {
		return errors.New("MSIX signing requires a SHA256 block map")
	}
	seen := map[string]bool{}
	buffer := make([]byte, 65536)
	for _, entry := range blockMap.Files {
		name := strings.ReplaceAll(entry.Name, "\\", "/")
		file := entries[name]
		if file == nil || seen[name] || file.UncompressedSize64 != entry.Size || entry.Header != headers[name] || uint64(len(entry.Blocks)) != (entry.Size+65535)/65536 {
			return fmt.Errorf("invalid MSIX block map entry %q", name)
		}

		var compressedSize uint64
		for _, block := range entry.Blocks {
			if file.Method == zip.Deflate {
				if block.Size == 0 || block.Size > file.CompressedSize64-compressedSize {
					return fmt.Errorf("invalid MSIX compressed block size for %q", name)
				}
				compressedSize += block.Size
			} else if block.Size != 0 {
				return fmt.Errorf("stored MSIX block has a compressed size for %q", name)
			}
		}
		if file.Method == zip.Deflate && compressedSize != file.CompressedSize64 && compressedSize+2 != file.CompressedSize64 {
			return fmt.Errorf("MSIX compressed blocks do not cover %q", name)
		}
		seen[name] = true
		reader, err := file.Open()
		if err != nil {
			return err
		}
		wholeHash := sha256.New()
		payload := io.TeeReader(reader, wholeHash)
		verifyErr := func() error {
			for _, block := range entry.Blocks {
				n, err := io.ReadFull(contextReader{ctx, payload}, buffer)
				if err != nil && err != io.ErrUnexpectedEOF {
					return fmt.Errorf("read MSIX block: %w", err)
				}
				expected, err := base64.StdEncoding.DecodeString(block.Hash)
				if err != nil {
					return err
				}
				actual := sha256.Sum256(buffer[:n])
				if !bytes.Equal(expected, actual[:]) {
					return fmt.Errorf("MSIX payload hash mismatch for %q", name)
				}
			}
			var tail [1]byte
			n, err := payload.Read(tail[:])
			if n != 0 || err != io.EOF {
				return fmt.Errorf("MSIX block map does not cover payload %q", name)
			}
			return nil
		}()
		closeErr := reader.Close()
		if verifyErr != nil {
			return verifyErr
		}
		if entry.FileHash != nil {
			expected, err := base64.StdEncoding.DecodeString(entry.FileHash.Hash)
			if err != nil || !bytes.Equal(expected, wholeHash.Sum(nil)) {
				return fmt.Errorf("MSIX whole-file hash mismatch for %q", name)
			}
		}
		if closeErr != nil {
			return closeErr
		}
		if file.Method == zip.Deflate {
			if err := verifyCompressedBlocks(ctx, file, entry.Blocks, buffer); err != nil {
				return fmt.Errorf("MSIX compressed blocks for %q: %w", name, err)
			}
		}
	}
	for name := range entries {
		if !seen[name] && name != "[Content_Types].xml" && name != "AppxBlockMap.xml" && name != signatureName && name != "AppxMetadata/CodeIntegrity.cat" {
			return fmt.Errorf("MSIX block map omits %q", name)
		}
	}
	return nil
}

type mappedBlock struct {
	Hash string `xml:"Hash,attr"`
	Size uint64 `xml:"Size,attr"`
}
type mappedFile struct {
	Name     string        `xml:"Name,attr"`
	Size     uint64        `xml:"Size,attr"`
	Header   int           `xml:"LfhSize,attr"`
	Blocks   []mappedBlock `xml:"Block"`
	FileHash *struct {
		Hash string `xml:"Hash,attr"`
	} `xml:"http://schemas.microsoft.com/appx/2021/blockmap FileHash"`
}

func verifyCompressedBlocks(ctx context.Context, file *zip.File, blocks []mappedBlock, buffer []byte) error {
	raw, err := file.OpenRaw()
	if err != nil {
		return err
	}
	var compressedSize uint64
	for _, block := range blocks {
		compressedSize += block.Size
	}
	trailer := file.CompressedSize64 - compressedSize
	for i, block := range blocks {
		limited := &io.LimitedReader{R: contextReader{ctx, raw}, N: int64(block.Size)}
		reader := flate.NewReader(limited)
		size := len(buffer)
		if remaining := file.UncompressedSize64 - uint64(i)*65536; remaining < uint64(size) {
			size = int(remaining)
		}
		_, readErr := io.ReadFull(reader, buffer[:size])
		var tail [1]byte
		extra, endErr := reader.Read(tail[:])
		closeErr := reader.Close()
		if readErr != nil {
			return readErr
		}
		if extra != 0 || (endErr != io.EOF && (endErr != io.ErrUnexpectedEOF || (i == len(blocks)-1 && trailer != 2))) || limited.N != 0 {
			return errors.New("compressed block boundaries do not match the block map")
		}
		if closeErr != nil && !((i < len(blocks)-1 || trailer == 2) && closeErr == io.ErrUnexpectedEOF) {
			return closeErr
		}
		expected, err := base64.StdEncoding.DecodeString(block.Hash)
		if err != nil {
			return err
		}
		actual := sha256.Sum256(buffer[:size])
		if !bytes.Equal(expected, actual[:]) {
			return errors.New("compressed block content does not match its hash")
		}
	}
	if trailer != 0 {
		var end [2]byte
		if trailer != 2 {
			return errors.New("unsupported DEFLATE trailer")
		}
		if _, err := io.ReadFull(raw, end[:]); err != nil {
			return err
		}
		if end != [2]byte{3, 0} {
			return errors.New("invalid DEFLATE final block")
		}
	}
	return nil
}
