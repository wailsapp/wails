package pcinames

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const downloadURL = "https://pci-ids.ucw.cz/v2.2/pci.ids.gz"
const maxDownloadBytes = 32 * 1024 * 1024

type config struct {
	root, path, cache  string
	cacheOnly, network bool
}

func Load() (*Database, error) {
	client := *http.DefaultClient
	if client.Timeout <= 0 || client.Timeout > 15*time.Second {
		client.Timeout = 15 * time.Second
	}
	return load(environment(), &client, downloadURL)
}

func environment() config {
	cache := os.Getenv("PCIDB_CACHE_PATH")
	if cache == "" {
		if home, err := os.UserHomeDir(); err == nil {
			cache = filepath.Join(home, ".cache", "pci.ids")
		}
	}
	return config{
		root: os.Getenv("PCIDB_CHROOT"), path: os.Getenv("PCIDB_PATH"), cache: cache,
		cacheOnly: enabled("PCIDB_CACHE_ONLY"), network: enabled("PCIDB_ENABLE_NETWORK_FETCH"),
	}
}

func enabled(key string) bool {
	value, _ := strconv.ParseBool(os.Getenv(key))
	return value
}

func (c config) paths() []string {
	paths := []string{c.path, c.cache}
	if c.cacheOnly {
		return paths
	}
	for _, path := range []string{
		"usr/share/hwdata/pci.ids", "usr/share/misc/pci.ids",
		"usr/share/hwdata/pci.ids.gz", "usr/share/misc/pci.ids.gz",
	} {
		paths = append(paths, filepath.Join(c.root, "/", path))
	}
	return paths
}

func load(c config, client *http.Client, endpoint string) (*Database, error) {
	for _, path := range c.paths() {
		if path == "" {
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		db, err := read(f)
		f.Close()
		return db, err
	}
	if !c.network {
		return nil, errors.New("no local PCI names database found")
	}
	return fetch(client, endpoint, c.cache)
}

func fetch(client *http.Client, endpoint, cache string) (*Database, error) {
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Wails doctor PCI names")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch PCI names: HTTP %d", response.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, maxDownloadBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxDownloadBytes {
		return nil, errors.New("PCI names download exceeds size limit")
	}
	db, err := read(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	// A usable download can still serve this report if the cache is read-only.
	if cache != "" {
		_ = store(cache, data)
	}
	return db, nil
}

func store(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".pci-names-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	_, writeErr := f.Write(data)
	closeErr := f.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
