// Package voidpackage writes XBPS archives for Void Linux. The wire contract is
// documented at https://man.voidlinux.org/xbps-create.1 and exercised against the
// official XBPS tools; no package-manager implementation is embedded here.
package voidpackage

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	"howett.net/plist"
)

type Metadata struct {
	Name, Version, Revision, Arch              string
	Description, Maintainer, Homepage, License string
	SourceVersion                              string
	Depends, Provides, Conflicts, Replaces     []string
}

type Entry struct {
	Source, Destination, Type string // Type is file, config, dir or symlink.
	Mode                      os.FileMode
	Owner, Group              string
	MTime                     time.Time
}

type Options struct {
	Metadata
	Contents        []Entry
	MTime           time.Time
	Install, Remove []byte
}

type identity struct{ filename, version, arch string }

func packageIdentity(m Metadata) (identity, error) {
	if !token(m.Name, "+.-_") || m.Name[0] == '-' || !token(m.Version, ".+~") {
		return identity{}, fmt.Errorf("invalid XBPS package name or version")
	}
	revision := m.Revision
	if revision == "" {
		revision = "1"
	}
	n, err := strconv.ParseUint(revision, 10, 64)
	if err != nil || n == 0 {
		return identity{}, fmt.Errorf("XBPS revision must be a positive integer")
	}
	arch := m.Arch
	switch arch {
	case "amd64":
		arch = "x86_64"
	case "arm64":
		arch = "aarch64"
	case "386":
		arch = "i686"
	case "x86_64", "x86_64-musl", "aarch64", "aarch64-musl", "i686", "armv6l", "armv6l-musl", "armv7l", "armv7l-musl", "noarch":
	default:
		return identity{}, fmt.Errorf("unsupported XBPS architecture %q; configure its Void architecture explicitly", arch)
	}
	version := m.Version + "_" + revision
	return identity{m.Name + "-" + version + "." + arch + ".xbps", version, arch}, nil
}

func token(s, punctuation string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || strings.ContainsRune(punctuation, c) {
			continue
		}
		return false
	}
	return true
}

func FileName(metadata Metadata) (string, error) {
	id, err := packageIdentity(metadata)
	return id.filename, err
}

type payload struct {
	Entry
	size int64
	hash string
}

// Write uses supported gzip compression and bounded payload buffers. Hashes are
// checked again while writing so a changing input cannot publish a false manifest.
func Write(ctx context.Context, output io.Writer, options Options) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	id, err := packageIdentity(options.Metadata)
	if err != nil {
		return err
	}
	options.Metadata, err = normalizeDependencies(options.Metadata)
	if err != nil {
		return err
	}
	description, _, _ := strings.Cut(options.Description, "\n")
	if description == "" || len([]rune(description)) > 80 {
		return fmt.Errorf("XBPS description needs a nonempty first line of at most 80 characters")
	}
	entries, err := prepare(ctx, options)
	if err != nil {
		return err
	}
	props := map[string]any{"architecture": id.arch, "pkgname": options.Name, "version": id.version, "pkgver": options.Name + "-" + id.version, "short_desc": description}
	for key, value := range map[string]string{"maintainer": options.Maintainer, "homepage": options.Homepage, "license": options.License, "long_desc": options.Description, "wails-source-version": options.SourceVersion} {
		if value != "" {
			props[key] = value
		}
	}
	for key, value := range map[string][]string{"run_depends": options.Depends, "provides": options.Provides, "conflicts": options.Conflicts, "replaces": options.Replaces} {
		if len(value) > 0 {
			props[key] = value
		}
	}
	files := map[string][]map[string]any{}
	var installed uint64
	var configFiles []string
	for _, entry := range entries {
		item := map[string]any{"file": entry.Destination}
		kind := "files"
		switch entry.Type {
		case "dir":
			kind = "dirs"
		case "symlink":
			kind = "links"
			target := entry.Source
			if !path.IsAbs(target) {
				target = path.Join(path.Dir(entry.Destination), target)
			}
			item["target"] = path.Clean(target)
		default:
			if entry.Type == "config" {
				kind = "conf_files"
				configFiles = append(configFiles, entry.Destination)
			}
			item["sha256"], item["size"] = entry.hash, uint64(entry.size)
			installed += uint64(entry.size)
		}
		files[kind] = append(files[kind], item)
	}
	props["installed_size"] = installed
	if len(options.Install) > 0 {
		props["install-script"] = options.Install
	}
	if len(options.Remove) > 0 {
		props["remove-script"] = options.Remove
	}
	if len(configFiles) > 0 {
		props["conf_files"] = configFiles
	}
	return writeArchive(ctx, output, options, entries, props, files)
}

func prepare(ctx context.Context, options Options) ([]payload, error) {
	contents := map[string]payload{}
	for _, entry := range options.Contents {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if !path.IsAbs(entry.Destination) || entry.Destination == "/" || path.Clean(entry.Destination) != entry.Destination || strings.ContainsRune(entry.Destination, 0) {
			return nil, fmt.Errorf("invalid XBPS destination %q", entry.Destination)
		}
		switch entry.Destination {
		case "/props.plist", "/files.plist", "/INSTALL", "/REMOVE":
			return nil, fmt.Errorf("reserved XBPS metadata destination %q", entry.Destination)
		}
		if _, ok := contents[entry.Destination]; ok {
			return nil, fmt.Errorf("duplicate XBPS destination %q", entry.Destination)
		}
		if entry.Owner == "" {
			entry.Owner = "root"
		}
		if entry.Group == "" {
			entry.Group = "root"
		}
		if entry.MTime.IsZero() {
			entry.MTime = options.MTime
		}
		if entry.Type == "" {
			entry.Type = "file"
		}
		p := payload{Entry: entry}
		switch entry.Type {
		case "file", "config":
			size, hash, err := fileHash(ctx, entry.Source)
			if err != nil {
				return nil, err
			}
			p.size, p.hash = size, hash
		case "dir":
			if p.Mode == 0 {
				p.Mode = 0755
			}
		case "symlink":
			if entry.Source == "" || strings.ContainsRune(entry.Source, 0) {
				return nil, fmt.Errorf("invalid XBPS link target")
			}
		default:
			return nil, fmt.Errorf("unsupported XBPS content type %q", entry.Type)
		}
		contents[entry.Destination] = p
	}
	for name := range contents {
		for parent := path.Dir(name); parent != "/"; parent = path.Dir(parent) {
			if entry, ok := contents[parent]; ok {
				if entry.Type != "dir" {
					return nil, fmt.Errorf("XBPS destination %q has a non-directory ancestor %q", name, parent)
				}
			} else {
				contents[parent] = payload{Entry: Entry{Destination: parent, Type: "dir", Mode: 0755, Owner: "root", Group: "root", MTime: options.MTime}}
			}
		}
	}
	entries := make([]payload, 0, len(contents))
	for _, entry := range contents {
		entries = append(entries, entry)
	}
	slices.SortFunc(entries, func(a, b payload) int { return strings.Compare(a.Destination, b.Destination) })
	return entries, nil
}

type contextReader struct {
	context.Context
	io.Reader
}

func (r contextReader) Read(data []byte) (int, error) {
	if err := r.Err(); err != nil {
		return 0, err
	}
	return r.Reader.Read(data)
}

func fileHash(ctx context.Context, filename string) (int64, string, error) {
	if err := ctx.Err(); err != nil {
		return 0, "", err
	}
	info, err := os.Stat(filename)
	if err != nil {
		return 0, "", err
	}
	if !info.Mode().IsRegular() {
		return 0, "", fmt.Errorf("XBPS input %q is not a regular file", filename)
	}
	f, err := os.Open(filename)
	if err != nil {
		return 0, "", err
	}
	defer f.Close()
	info, err = f.Stat()
	if err != nil {
		return 0, "", err
	}
	if !info.Mode().IsRegular() {
		return 0, "", fmt.Errorf("XBPS input %q is not a regular file", filename)
	}
	h := sha256.New()
	size, err := io.Copy(h, contextReader{ctx, f})
	return size, hex.EncodeToString(h.Sum(nil)), err
}

func writeArchive(ctx context.Context, output io.Writer, options Options, entries []payload, props map[string]any, files map[string][]map[string]any) (err error) {
	compressed := gzip.NewWriter(output)
	archive := tar.NewWriter(compressed)
	defer func() { err = errors.Join(err, archive.Close(), compressed.Close()) }()
	for _, script := range []struct {
		name string
		data []byte
	}{{"INSTALL", options.Install}, {"REMOVE", options.Remove}} {
		if len(script.data) > 0 {
			if err := writeData(archive, script.name, script.data, 0755, options.MTime); err != nil {
				return err
			}
		}
	}
	for _, metadata := range []struct {
		name  string
		value any
	}{{"props.plist", props}, {"files.plist", files}} {
		data, err := plist.Marshal(metadata.value, plist.XMLFormat)
		if err != nil {
			return err
		}
		if err := writeData(archive, metadata.name, data, 0644, options.MTime); err != nil {
			return err
		}
	}
	buffer := make([]byte, 64*1024)
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return err
		}
		header := tar.Header{Name: "." + entry.Destination, Mode: tarMode(entry.Mode), Uname: entry.Owner, Gname: entry.Group, ModTime: entry.MTime, Size: entry.size, Typeflag: tar.TypeReg, Format: tar.FormatPAX}
		switch entry.Type {
		case "dir":
			header.Typeflag = tar.TypeDir
			header.Name += "/"
		case "symlink":
			header.Typeflag = tar.TypeSymlink
			header.Linkname = entry.Source
		}
		if err := archive.WriteHeader(&header); err != nil {
			return err
		}
		if header.Typeflag == tar.TypeReg {
			if err := copyPayload(ctx, archive, entry, buffer); err != nil {
				return err
			}
		}
	}
	return ctx.Err()
}

func tarMode(mode os.FileMode) int64 {
	bits := int64(mode & 07777)
	if mode&os.ModeSetuid != 0 {
		bits |= 04000
	}
	if mode&os.ModeSetgid != 0 {
		bits |= 02000
	}
	if mode&os.ModeSticky != 0 {
		bits |= 01000
	}
	return bits
}

func copyPayload(ctx context.Context, archive io.Writer, entry payload, buffer []byte) error {
	f, err := os.Open(entry.Source)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	size, err := io.CopyBuffer(io.MultiWriter(archive, h), contextReader{ctx, f}, buffer)
	if err != nil {
		return err
	}
	if size != entry.size || hex.EncodeToString(h.Sum(nil)) != entry.hash {
		return fmt.Errorf("XBPS input %q changed while packaging", entry.Source)
	}
	return nil
}

func writeData(archive *tar.Writer, name string, data []byte, mode int64, mtime time.Time) error {
	if err := archive.WriteHeader(&tar.Header{Name: "./" + name, Size: int64(len(data)), Mode: mode, Uname: "root", Gname: "root", ModTime: mtime, Typeflag: tar.TypeReg, Format: tar.FormatPAX}); err != nil {
		return err
	}
	_, err := archive.Write(data)
	return err
}
