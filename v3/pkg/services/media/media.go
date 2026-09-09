// Package media transfers local clips over Wails streams for playback with
// runtime Media.SetSource. Desktop transfers use the existing asset transport
// without opening a listening socket.
package media

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"mime"
	"path"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const chunkBytes = 64 << 10

// NewHandler creates a stream handler for regular files in files. Register it
// with app.HandleStream. maxBytes must be positive and is the application-side
// limit on each complete clip; the frontend can request a smaller limit.
//
// Only expose a filesystem containing media the frontend may read. For disk
// files, os.Root.FS can confine access even when the directory contains symlinks.
// The caller owns the filesystem and must keep it open for the handler's lifetime.
// Transfers are chunked, but the frontend retains the complete clip as a blob.
func NewHandler(files fs.FS, maxBytes int64) (application.StreamHandler, error) {
	if files == nil || maxBytes <= 0 {
		return nil, errors.New("media requires a filesystem and a positive byte limit")
	}
	return func(c *application.StreamConn) { transfer(c, files, maxBytes) }, nil
}

type connection interface {
	Context() context.Context
	Receive() ([]byte, error)
	Send([]byte) error
}

type request struct {
	Version  int    `json:"version"`
	Name     string `json:"name"`
	MaxBytes int64  `json:"maxBytes"`
}

type metadata struct {
	Version int    `json:"version"`
	Size    int64  `json:"size"`
	Type    string `json:"type"`
	Error   string `json:"error,omitempty"`
	Limit   bool   `json:"limit,omitempty"`
}

func sendMetadata(c connection, value metadata) error {
	value.Version = 1
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Send(data)
}

func transfer(c connection, files fs.FS, maxBytes int64) {
	frame, err := c.Receive()
	if err != nil {
		return
	}
	var req request
	if len(frame) > 4096 || json.Unmarshal(frame, &req) != nil || req.Version != 1 || req.MaxBytes <= 0 ||
		!fs.ValidPath(req.Name) || req.Name == "." || strings.ContainsAny(req.Name, `\:`) {
		_ = sendMetadata(c, metadata{Error: "Invalid media request"})
		return
	}
	file, err := files.Open(req.Name)
	if err != nil {
		_ = sendMetadata(c, metadata{Error: "Media file could not be opened"})
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() < 0 {
		_ = sendMetadata(c, metadata{Error: "Media must be a regular file with a known size"})
		return
	}
	limit := min(maxBytes, req.MaxBytes)
	if info.Size() > limit {
		_ = sendMetadata(c, metadata{Error: "Media exceeds the byte limit", Limit: true})
		return
	}
	contentType := mime.TypeByExtension(path.Ext(req.Name))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if sendMetadata(c, metadata{Size: info.Size(), Type: contentType}) != nil {
		return
	}
	// Each Send hands ownership of its slice to the transport. Allocate each
	// chunk separately; reusing a scratch buffer would corrupt queued frames.
	for remaining := info.Size(); remaining > 0; {
		if c.Context().Err() != nil {
			return
		}
		chunk := make([]byte, min(int64(chunkBytes), remaining))
		if _, err := io.ReadFull(file, chunk); err != nil {
			// A short transfer is rejected by the frontend using the declared size.
			return
		}
		if c.Send(chunk) != nil {
			return
		}
		remaining -= int64(len(chunk))
	}
}
