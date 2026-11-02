package windowsign

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"github.com/wailsapp/wails/v3/internal/codesignature"
	"github.com/wailsapp/wails/v3/internal/windowspackage"
	"hash/crc32"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func signingIdentity(t testing.TB) Identity {
	t.Helper()
	data, err := os.ReadFile("../codesignature/testdata/modern.pfx")
	if err != nil {
		t.Fatal(err)
	}
	id, err := DecodePFX(data, "fixture-only")
	if err != nil {
		t.Fatal(err)
	}
	return id
}
func peFixture(t testing.TB) string {
	t.Helper()
	data := make([]byte, 1024)
	copy(data, "MZ")
	binary.LittleEndian.PutUint32(data[60:], 128)
	copy(data[128:], "PE\x00\x00")
	binary.LittleEndian.PutUint16(data[132:], 0x8664)
	binary.LittleEndian.PutUint16(data[134:], 1)
	binary.LittleEndian.PutUint16(data[148:], 240)
	binary.LittleEndian.PutUint16(data[150:], 2)
	o := data[152:]
	binary.LittleEndian.PutUint16(o, 0x20b)
	binary.LittleEndian.PutUint32(o[32:], 4096)
	binary.LittleEndian.PutUint32(o[36:], 512)
	binary.LittleEndian.PutUint32(o[56:], 8192)
	binary.LittleEndian.PutUint32(o[60:], 512)
	binary.LittleEndian.PutUint32(o[108:], 16)
	section := data[392:]
	copy(section, ".text")
	binary.LittleEndian.PutUint32(section[8:], 512)
	binary.LittleEndian.PutUint32(section[12:], 4096)
	binary.LittleEndian.PutUint32(section[16:], 512)
	binary.LittleEndian.PutUint32(section[20:], 512)
	for i := 512; i < len(data); i++ {
		data[i] = byte(i * 19)
	}
	path := filepath.Join(t.TempDir(), "app.exe")
	if err := os.WriteFile(path, data, 0755); err != nil {
		t.Fatal(err)
	}
	return path
}
func checkCMS(t *testing.T, encoded []byte, id Identity) []byte {
	t.Helper()
	var outer struct {
		Type asn1.ObjectIdentifier
		Data asn1.RawValue `asn1:"explicit,tag:0"`
	}
	if _, err := asn1.Unmarshal(encoded, &outer); err != nil {
		t.Fatal(err)
	}
	var data struct {
		Version int
		Digests []pkix.AlgorithmIdentifier `asn1:"set"`
		Content struct {
			Type  asn1.ObjectIdentifier
			Value asn1.RawValue `asn1:"explicit,tag:0"`
		}
		Certs   asn1.RawValue `asn1:"tag:0"`
		Signers []struct {
			Version int
			Issuer  struct {
				Issuer asn1.RawValue
				Serial *big.Int
			}
			Digest    pkix.AlgorithmIdentifier
			Attrs     []codesignature.Attribute `asn1:"tag:0,set"`
			Algorithm pkix.AlgorithmIdentifier
			Signature []byte
			Unsigned  asn1.RawValue `asn1:"optional,tag:1"`
		} `asn1:"set"`
	}
	if _, err := asn1.Unmarshal(outer.Data.Bytes, &data); err != nil {
		t.Fatal(err)
	}
	if len(data.Signers) != 1 {
		t.Fatal("expected one CMS signer")
	}
	signer := data.Signers[0]
	attrs, err := asn1.MarshalWithParams(signer.Attrs, "set")
	if err != nil {
		t.Fatal(err)
	}
	if err := id.Certificate.CheckSignature(x509.SHA256WithRSA, attrs, signer.Signature); err != nil {
		t.Fatal(err)
	}
	var indirect asn1.RawValue
	if _, err := asn1.Unmarshal(data.Content.Value.Bytes, &indirect); err != nil {
		t.Fatal(err)
	}
	expected := sha256.Sum256(indirect.Bytes)
	matched := false
	for _, attr := range signer.Attrs {
		if attr.Type.String() == "1.2.840.113549.1.9.4" {
			var digest []byte
			if _, err := asn1.Unmarshal(attr.Values[0].FullBytes, &digest); err != nil {
				t.Fatal(err)
			}
			matched = bytes.Equal(digest, expected[:])
		}
	}
	if !matched {
		t.Fatal("CMS content digest mismatch")
	}
	var content struct {
		Format asn1.RawValue
		Digest struct {
			Algorithm pkix.AlgorithmIdentifier
			Value     []byte
		}
	}
	if _, err := asn1.Unmarshal(indirect.FullBytes, &content); err != nil {
		t.Fatal(err)
	}
	return content.Digest.Value
}
func TestSignPEAndReplaceSignature(t *testing.T) {
	id := signingIdentity(t)
	path := peFixture(t)
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := SignPE(context.Background(), path, path, id, Options{}); err != nil {
			t.Fatal(err)
		}
		file, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		info, _ := file.Stat()
		layout, err := inspectPE(file, info.Size())
		if err != nil {
			t.Fatal(err)
		}
		if layout.end != int64(len(original)) {
			t.Fatalf("signature table grew into the image: %d", layout.end)
		}
		table := make([]byte, info.Size()-layout.end)
		if _, err := file.ReadAt(table, layout.end); err != nil {
			t.Fatal(err)
		}
		file.Close()
		digest := checkCMS(t, table[8:], id)
		hash := sha256.New()
		hash.Write(original[:216])
		hash.Write(original[220:296])
		hash.Write(original[304:])
		if !bytes.Equal(digest, hash.Sum(nil)) {
			t.Fatal("PE signature covers incorrect bytes")
		}
	}
}

type mutationSigner struct {
	crypto.Signer
	mutate func()
}

func (s mutationSigner) Sign(random io.Reader, digest []byte, options crypto.SignerOpts) ([]byte, error) {
	s.mutate()
	return s.Signer.Sign(random, digest, options)
}
func TestPESigningUsesImmutableSnapshot(t *testing.T) {
	id := signingIdentity(t)
	input := peFixture(t)
	expected, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	id.Signer = mutationSigner{id.Signer, func() {
		changed := bytes.Clone(expected)
		changed[600] ^= 1
		if err := os.WriteFile(input, changed, 0755); err != nil {
			t.Error(err)
		}
	}}
	output := filepath.Join(t.TempDir(), "signed.exe")
	if err := SignPE(context.Background(), input, output, id, Options{}); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual[512:1024], expected[512:]) {
		t.Fatal("published bytes changed after hashing")
	}
}
func msixFixture(t testing.TB) string {
	t.Helper()
	root := t.TempDir()
	manifest := `<Package xmlns="http://schemas.microsoft.com/appx/manifest/foundation/windows10"><Identity Publisher="CN=Wails Portable Signing Test Only"/></Package>`
	if err := os.WriteFile(filepath.Join(root, "AppxManifest.xml"), []byte(manifest), 0600); err != nil {
		t.Fatal(err)
	}
	payload := make([]byte, 180000)
	for i := range payload {
		payload[i] = byte(i * 19)
	}
	if err := os.WriteFile(filepath.Join(root, "payload.bin"), payload, 0600); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"space name.txt", "français.txt", "empty.txt"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	output := filepath.Join(t.TempDir(), "app.msix")
	if err := windowspackage.Create(context.Background(), windowspackage.Options{Directory: root, Output: output}); err != nil {
		t.Fatal(err)
	}
	return output
}
func TestSignMSIXRawRecordsAndResign(t *testing.T) {
	id := signingIdentity(t)
	path := msixFixture(t)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for iteration := 0; iteration < 2; iteration++ {
		if err := SignMSIX(context.Background(), path, path, id, Options{}); err != nil {
			t.Fatal(err)
		}
		archive, err := zip.OpenReader(path)
		if err != nil {
			t.Fatal(err)
		}
		var signature *zip.File
		for _, entry := range archive.File {
			if entry.Name == signatureName {
				if signature != nil {
					t.Fatal("multiple signatures")
				}
				signature = entry
			}
		}
		if signature == nil || signature.CreatorVersion != 45 {
			t.Fatal("signature creator must match ZIP64 footer")
		}
		reader, err := signature.Open()
		if err != nil {
			t.Fatal(err)
		}
		sig, err := io.ReadAll(reader)
		reader.Close()
		archive.Close()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.HasPrefix(sig, []byte("PKCX")) {
			t.Fatal("incorrect P7X prefix")
		}
		bag := checkCMS(t, sig[4:], id)
		if string(bag[:8]) != "APPXAXPC" {
			t.Fatal("invalid digest bag")
		}
		file, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		info, err := file.Stat()
		if err != nil {
			t.Fatal(err)
		}
		records, err := packageRecords(file, info.Size())
		if err != nil {
			t.Fatal(err)
		}
		for _, record := range records {
			if record.name == signatureName {
				hash, err := hashReader(context.Background(), io.NewSectionReader(file, 0, record.offset))
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(hash, bag[8:40]) {
					t.Fatal("package record digest mismatch")
				}
			}
		}
		file.Close()
		after, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		original, err := zip.NewReader(bytes.NewReader(before), int64(len(before)))
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range original.File {
			if entry.Name == "[Content_Types].xml" {
				offset, err := entry.DataOffset()
				if err != nil {
					t.Fatal(err)
				}
				prefix := offset - int64(30+len(entry.Name)+len(entry.Extra))
				if !bytes.Equal(before[:prefix], after[:prefix]) {
					t.Fatal("existing local records changed")
				}
			}
		}
	}
}
func TestSigningFailurePreservesFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "fake timestamp") }))
	defer server.Close()
	id := signingIdentity(t)
	for _, test := range []struct {
		name    string
		fixture func(testing.TB) string
		sign    func(context.Context, string, string, Identity, Options) error
	}{{"PE", peFixture, SignPE}, {"MSIX", msixFixture, SignMSIX}} {
		t.Run(test.name, func(t *testing.T) {
			path := test.fixture(t)
			before, _ := os.ReadFile(path)
			if err := test.sign(context.Background(), path, path, id, Options{TimestampURL: server.URL}); err == nil {
				t.Fatal("fake timestamp accepted")
			}
			after, _ := os.ReadFile(path)
			if !bytes.Equal(before, after) {
				t.Fatal("failed signing modified input")
			}
			files, err := filepath.Glob(filepath.Join(filepath.Dir(path), ".wake-*"))
			if err != nil || len(files) != 0 {
				t.Fatalf("staged files leaked: %v %v", files, err)
			}
		})
	}
}
func TestMSIXInvalidIdentityAndContentType(t *testing.T) {
	if err := SignMSIX(context.Background(), "unused", "unused", Identity{}, Options{}); err == nil {
		t.Fatal("empty identity accepted")
	}
	data := []byte(`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Override PartName="/AppxSignature.p7x" ContentType="wrong"/></Types>`)
	if _, err := signedContentTypes(data); err == nil {
		t.Fatal("wrong signature content type accepted")
	}
}
func TestPublisherWindowsNames(t *testing.T) {
	subject := pkix.Name{CommonName: "Example, Inc.", Organization: []string{"Example"}, Locality: []string{"Redmond"}, Province: []string{"Washington"}, Country: []string{"US"}}
	raw, err := asn1.Marshal(subject.ToRDNSequence())
	if err != nil {
		t.Fatal(err)
	}
	cert := &x509.Certificate{RawSubject: raw}
	valid := `CN="Example, Inc.", O=Example, L=Redmond, S=Washington, C=US`
	for _, publisher := range []string{valid, `CN=Example\, Inc.,O=Example,L=Redmond,ST=Washington,C=US`} {
		match, err := publisherMatches(publisher, cert)
		if err != nil || !match {
			t.Fatalf("valid publisher %q: %v", publisher, err)
		}
	}
	for _, publisher := range []string{strings.Replace(valid, "O=Example", "O=Different", 1), `CN="unterminated`, strings.Replace(valid, ", O=Example", ", O=Example, O=Example", 1), `C=US,S=Washington,L=Redmond,O=Example,CN="Example, Inc."`} {
		match, err := publisherMatches(publisher, cert)
		if err == nil && match {
			t.Fatalf("incorrect publisher accepted: %q", publisher)
		}
	}
}

func TestMSIXBlockMapValidation(t *testing.T) {
	path := msixFixture(t)
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	archive, err := zip.NewReader(file, info.Size())
	if err != nil {
		t.Fatal(err)
	}
	entries := map[string]*zip.File{}
	for _, entry := range archive.File {
		name, err := url.PathUnescape(entry.Name)
		if err != nil {
			t.Fatal(err)
		}
		entries[name] = entry
	}
	records, err := packageRecords(file, info.Size())
	if err != nil {
		t.Fatal(err)
	}
	headers := map[string]int{}
	for _, record := range records {
		name, _ := url.PathUnescape(record.name)
		headers[name] = record.header
	}
	original, err := readZIPMetadata(entries["AppxBlockMap.xml"])
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		change func(*mappedFile)
	}{
		{"wrong header length", func(f *mappedFile) { f.Header++ }},
		{"wrong block boundary", func(f *mappedFile) { f.Blocks[0].Size++; f.Blocks[1].Size-- }},
		{"wrong payload digest", func(f *mappedFile) { f.Blocks[0].Hash = base64.StdEncoding.EncodeToString(make([]byte, 32)) }},
	} {
		t.Run(test.name, func(t *testing.T) {
			var blockMap struct {
				XMLName xml.Name
				Hash    string       `xml:"HashMethod,attr"`
				Files   []mappedFile `xml:"File"`
			}
			if err := decodeXML(original, &blockMap); err != nil {
				t.Fatal(err)
			}
			for i := range blockMap.Files {
				if blockMap.Files[i].Name == "payload.bin" {
					test.change(&blockMap.Files[i])
				}
			}
			encoded, err := xml.Marshal(blockMap)
			if err != nil {
				t.Fatal(err)
			}
			if err := verifyBlockMap(context.Background(), encoded, entries, headers); err == nil {
				t.Fatal("invalid block map accepted")
			}
		})
	}
}
func BenchmarkSignPE(b *testing.B) {
	path := peFixture(b)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	if _, err := file.Write(make([]byte, 4<<20)); err != nil {
		b.Fatal(err)
	}
	file.Close()
	id := signingIdentity(b)
	output := filepath.Join(b.TempDir(), "signed.exe")
	b.SetBytes(4 << 20)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := SignPE(context.Background(), path, output, id, Options{}); err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkSignMSIX(b *testing.B) {
	path := msixFixture(b)
	id := signingIdentity(b)
	output := filepath.Join(b.TempDir(), "signed.msix")
	b.SetBytes(180000)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := SignMSIX(context.Background(), path, output, id, Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func TestSignMSIXWithNativeZIP64Metadata(t *testing.T) {
	input := msixFixture(t)
	source, err := os.Open(input)
	if err != nil {
		t.Fatal(err)
	}
	defer source.Close()
	info, _ := source.Stat()
	records, err := packageRecords(source, info.Size())
	if err != nil {
		t.Fatal(err)
	}
	start, err := zipDirectory(source, info.Size())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "zip64.msix")
	output, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(output, io.NewSectionReader(source, 0, start)); err != nil {
		t.Fatal(err)
	}
	var directory []byte
	for _, record := range records {
		entry := bytes.Clone(record.directory)
		nameLength := int(binary.LittleEndian.Uint16(entry[28:]))
		extraLength := int(binary.LittleEndian.Uint16(entry[30:]))
		position := 46 + nameLength
		extra := make([]byte, 28)
		binary.LittleEndian.PutUint16(extra, 1)
		binary.LittleEndian.PutUint16(extra[2:], 24)
		binary.LittleEndian.PutUint64(extra[4:], uint64(binary.LittleEndian.Uint32(entry[24:])))
		binary.LittleEndian.PutUint64(extra[12:], uint64(binary.LittleEndian.Uint32(entry[20:])))
		binary.LittleEndian.PutUint64(extra[20:], uint64(record.offset))
		entry = append(append(append([]byte{}, entry[:position]...), extra...), entry[position:]...)
		binary.LittleEndian.PutUint16(entry[30:], uint16(extraLength+28))
		for _, offset := range []int{20, 24, 42} {
			binary.LittleEndian.PutUint32(entry[offset:], 0xffffffff)
		}
		directory = append(directory, entry...)
	}
	if err := finishMSIX(output, start, directory, len(records)); err != nil {
		t.Fatal(err)
	}
	if err := output.Close(); err != nil {
		t.Fatal(err)
	}
	id := signingIdentity(t)
	for i := 0; i < 2; i++ {
		if err := SignMSIX(context.Background(), path, path, id, Options{}); err != nil {
			t.Fatal(err)
		}
	}
}
func TestNativeDeflateTrailer(t *testing.T) {
	payload := bytes.Repeat([]byte("native-style block "), 500)
	var compressed bytes.Buffer
	writer, err := flate.NewWriter(&compressed, flate.DefaultCompression)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write(payload); err != nil {
		t.Fatal(err)
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	// Flush leaves a byte-aligned stream without a final block. Snapshot it
	// before Close so the fixture does not depend on Go's choice of terminator.
	data := bytes.Clone(compressed.Bytes())
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	data = append(data, 3, 0)
	var archive bytes.Buffer
	zipWriter := zip.NewWriter(&archive)
	entry, err := zipWriter.CreateRaw(&zip.FileHeader{Name: "payload", Method: zip.Deflate, CompressedSize64: uint64(len(data)), UncompressedSize64: uint64(len(payload)), CRC32: crc32.ChecksumIEEE(payload)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := zipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	zipReader, err := zip.NewReader(bytes.NewReader(archive.Bytes()), int64(archive.Len()))
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(payload)
	blocks := []mappedBlock{{Hash: base64.StdEncoding.EncodeToString(digest[:]), Size: uint64(len(data) - 2)}}
	if err := verifyCompressedBlocks(context.Background(), zipReader.File[0], blocks, make([]byte, 65536)); err != nil {
		t.Fatal(err)
	}
}
