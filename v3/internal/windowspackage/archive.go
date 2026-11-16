// Package windowspackage writes unsigned MSIX containers using Go's ZIP,
// DEFLATE, XML and SHA256 implementations. It does not invoke packaging tools.
package windowspackage

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"
)

const blockSize = 65536

type Options struct {
	Directory   string
	Output      string
	Compression int // 0 selects level 4; 1..9 selects DEFLATE effort.
	Prepare     func(context.Context, string) error
}

type blockMap struct {
	XMLName xml.Name   `xml:"http://schemas.microsoft.com/appx/2010/blockmap BlockMap"`
	Hash    string     `xml:"HashMethod,attr"`
	Files   []fileHash `xml:"File"`
}

type fileHash struct {
	Name   string      `xml:"Name,attr"`
	Size   uint64      `xml:"Size,attr"`
	Header int         `xml:"LfhSize,attr"`
	Blocks []blockHash `xml:"Block"`
}

type blockHash struct {
	Hash string `xml:"Hash,attr"`
	Size uint64 `xml:"Size,attr,omitempty"`
}

// Create publishes only after every input, block map and ZIP directory has
// been written successfully. Compression is bounded by a block buffer and a
// temporary file, rather than holding each executable in memory.
func Create(ctx context.Context, options Options) error {
	if options.Compression < 0 || options.Compression > 9 {
		return fmt.Errorf("MSIX compression must be between 0 and 9")
	}
	root, err := filepath.Abs(options.Directory)
	if err != nil {
		return err
	}
	output, err := filepath.Abs(options.Output)
	if err != nil {
		return err
	}
	if relative, err := filepath.Rel(root, output); err == nil && filepath.IsLocal(relative) {
		return fmt.Errorf("MSIX output must be outside the input directory")
	}
	files, err := packageFiles(ctx, root)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(filepath.Dir(output), ".wake-msix-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	path := filepath.Join(stage, "package.msix")
	if err = writePackage(ctx, options, root, stage, path, files); err != nil {
		return err
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	if options.Prepare != nil {
		if err = options.Prepare(ctx, path); err != nil {
			return err
		}
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	return os.Rename(path, output)
}

func packageFiles(ctx context.Context, root string) ([]string, error) {
	var files []string
	seen := make(map[string]bool)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if err = ctx.Err(); err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("MSIX input is not a regular file: %s", path)
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(relative)
		if err = validateName(name); err != nil {
			return err
		}
		key := strings.ToLower(name)
		if key == "appxmanifest.xml" && name != "AppxManifest.xml" {
			return fmt.Errorf("MSIX manifest must be named AppxManifest.xml")
		}
		if seen[key] || key == "appxblockmap.xml" || key == "[content_types].xml" || key == "appxsignature.p7x" {
			return fmt.Errorf("duplicate or reserved MSIX entry %q", name)
		}
		seen[key] = true
		files = append(files, name)
		return nil
	})
	if err == nil && !seen["appxmanifest.xml"] {
		err = fmt.Errorf("MSIX requires AppxManifest.xml")
	}
	if err == nil {
		err = validateManifest(ctx, filepath.Join(root, "AppxManifest.xml"))
	}
	return files, err
}

func validateManifest(ctx context.Context, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := xml.NewDecoder(contextReader{ctx, file})
	depth, roots := 0, 0
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			if roots != 1 {
				return fmt.Errorf("MSIX manifest requires one Package root")
			}
			return nil
		}
		if err != nil {
			return fmt.Errorf("read MSIX manifest: %w", err)
		}
		switch token := token.(type) {
		case xml.StartElement:
			if depth == 0 {
				roots++
				if roots != 1 || token.Name.Local != "Package" {
					return fmt.Errorf("MSIX manifest requires one Package root")
				}
			}
			depth++
		case xml.EndElement:
			depth--
		case xml.CharData:
			if depth == 0 && strings.TrimSpace(string(token)) != "" {
				return fmt.Errorf("text outside MSIX manifest root")
			}
		}
	}
}

func validateName(name string) error {
	if len(utf16.Encode([]rune(name))) > 260 || strings.ContainsAny(name, "\\<>:\"|?*\x00\r\n") {
		return fmt.Errorf("invalid MSIX entry name %q", name)
	}
	for _, component := range strings.Split(name, "/") {
		if component == "" || component == "." || component == ".." || strings.HasSuffix(component, ".") || strings.HasSuffix(component, " ") {
			return fmt.Errorf("invalid MSIX entry name %q", name)
		}
	}
	return nil
}

func writePackage(ctx context.Context, options Options, root, stage, path string, files []string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, file.Close()) }()
	archive := zip.NewWriter(file)
	defer func() { err = errors.Join(err, archive.Close()) }()
	hashes := blockMap{Hash: "http://www.w3.org/2001/04/xmlenc#sha256"}
	buffer := make([]byte, blockSize)
	for _, name := range files {
		hash, err := addPayload(ctx, archive, filepath.Join(root, filepath.FromSlash(name)), name, stage, options.Compression, buffer)
		if err != nil {
			return err
		}
		hashes.Files = append(hashes.Files, hash)
	}
	data, err := xml.Marshal(hashes)
	if err != nil {
		return err
	}
	if err = writeFootprint(archive, "AppxBlockMap.xml", data, options.Compression); err != nil {
		return err
	}
	return writeFootprint(archive, "[Content_Types].xml", contentTypes(files), options.Compression)
}

type contentType struct {
	Part string `xml:"PartName,attr"`
	Type string `xml:"ContentType,attr"`
}

func contentTypes(files []string) []byte {
	content := struct {
		XMLName xml.Name      `xml:"http://schemas.openxmlformats.org/package/2006/content-types Types"`
		Parts   []contentType `xml:"Override"`
	}{}
	for _, name := range append(files, "AppxBlockMap.xml") {
		kind := "application/octet-stream"
		switch name {
		case "AppxManifest.xml":
			kind = "application/vnd.ms-appx.manifest+xml"
		case "AppxBlockMap.xml":
			kind = "application/vnd.ms-appx.blockmap+xml"
		}
		uri := (&url.URL{Path: "/" + name}).EscapedPath()
		content.Parts = append(content.Parts, contentType{uri, kind})
	}
	data, _ := xml.Marshal(content) // Only strings and XML names; no failing field types.
	return data
}

func writeFootprint(archive *zip.Writer, name string, data []byte, level int) error {
	header, compressed, err := EncodeFootprint(name, data, level)
	if err != nil {
		return err
	}
	writer, err := archive.CreateRaw(header)
	if err != nil {
		return err
	}
	_, err = writer.Write(compressed)
	return err
}

// EncodeFootprint prepares a small MSIX metadata record without a data descriptor.
// Signing also uses it, adjusting the ZIP creator version for the signature.
func EncodeFootprint(name string, data []byte, level int) (*zip.FileHeader, []byte, error) {
	if level == 0 {
		level = 4
	}
	var compressed bytes.Buffer
	compressor, err := flate.NewWriter(&compressed, level)
	if err != nil {
		return nil, nil, err
	}
	if _, err := compressor.Write(data); err != nil {
		return nil, nil, err
	}
	if err := compressor.Close(); err != nil {
		return nil, nil, err
	}
	header := &zip.FileHeader{Name: name, Method: zip.Deflate, Flags: 0x800, ReaderVersion: 20, CreatorVersion: 20, ModifiedDate: 33, CRC32: crc32.ChecksumIEEE(data), CompressedSize64: uint64(compressed.Len()), UncompressedSize64: uint64(len(data))}
	return header, compressed.Bytes(), nil
}

type contextReader struct {
	ctx context.Context
	io.Reader
}

func (r contextReader) Read(data []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.Reader.Read(data)
}
