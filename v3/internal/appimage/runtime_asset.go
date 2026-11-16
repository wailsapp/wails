package appimage

import (
	"context"
	"crypto/sha256"
	"debug/elf"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type runtimeAsset struct{ id, sha256 string }

// Official type2-runtime release, commit fedfad51acc2860b2916a1eef32ba4a500ff31ef.
// IDs identify immutable release assets; the digest is verified before use.
var runtimeAssets = map[string]runtimeAsset{
	"amd64": {"657374641", "617aa77314d57c8b8e4fae9c3f0b35af7b30d5529cb69e98f64f92510ed98a46"},
	"arm64": {"657373827", "9a72a744d221576c9347ec074243acc5c3c7568896b12714e0e82e9f0bf47900"},
}

type runtimeTransport struct{ base http.RoundTripper }

func (t runtimeTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	copy := request.Clone(request.Context())
	copy.Header.Set("Accept", "application/octet-stream")
	return t.base.RoundTrip(copy)
}

func (options Options) PrepareRuntime(ctx context.Context) (string, error) {
	if _, err := architectureLabel(options.Arch); err != nil {
		return "", err
	}
	if options.RuntimePath != "" {
		filename := options.path(options.RuntimePath)
		return filename, verify(filename, options.Arch)
	}
	asset, ok := runtimeAssets[options.Arch]
	if !ok {
		return "", fmt.Errorf("unsupported AppImage runtime architecture %s", options.Arch)
	}
	cache := options.CacheDir
	if cache == "" {
		cache = ".wake/appimage-tools"
	}
	filename := filepath.Join(options.path(cache), options.Arch, "runtime-"+asset.sha256)
	client := http.Client{Timeout: 5 * time.Minute}
	if options.Client != nil {
		client = *options.Client
	}
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	client.Transport = runtimeTransport{transport}
	url := "https://api.github.com/repos/AppImage/type2-runtime/releases/assets/" + asset.id
	if err := ensureTool(ctx, &client, url, filename, options.ForceTools); err != nil {
		return "", err
	}
	input, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	_, err = io.Copy(contextWriter{ctx, hash}, input)
	input.Close()
	if err != nil {
		return "", err
	}
	if fmt.Sprintf("%x", hash.Sum(nil)) != asset.sha256 {
		os.Remove(filename)
		return "", fmt.Errorf("AppImage runtime SHA256 does not match pinned release asset %s", asset.id)
	}
	return filename, verify(filename, options.Arch)
}

func copyRuntime(ctx context.Context, source string, output *os.File, update string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	if _, err = io.Copy(contextWriter{ctx, output}, input); err != nil {
		return err
	}
	if update == "" {
		return nil
	}
	if strings.ContainsAny(update, "\x00\r\n") {
		return fmt.Errorf("invalid AppImage update information")
	}
	image, err := elf.NewFile(input)
	if err != nil {
		return err
	}
	section := image.Section(".upd_info")
	info, err := input.Stat()
	if err != nil {
		return err
	}
	if section == nil || section.Offset > uint64(info.Size()) || section.Size > uint64(info.Size())-section.Offset || uint64(len(update)+1) > section.Size {
		return fmt.Errorf("AppImage runtime has no space for update information")
	}
	value := make([]byte, section.Size)
	copy(value, update)
	_, err = output.WriteAt(value, int64(section.Offset))
	return err
}
