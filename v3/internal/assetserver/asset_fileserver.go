package assetserver

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	iofs "io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

const (
	indexHTML               = "index.html"
	assetFileCopyBufferSize = 32 * 1024
)

var assetFileCopyBuffers = sync.Pool{
	New: func() any { return new([assetFileCopyBufferSize]byte) },
}

type assetFileServer struct {
	fs  iofs.FS
	err error
}

func newAssetFileServerFS(vfs iofs.FS) http.Handler {
	subDir, err := findPathToFile(vfs, indexHTML)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			msg := "no `index.html` could be found in your Assets fs.FS"
			if embedFs, isEmbedFs := vfs.(embed.FS); isEmbedFs {
				rootFolder, _ := findEmbedRootPath(embedFs)
				msg += fmt.Sprintf(", please make sure the embedded directory '%s' is correct and contains your assets", rootFolder)
			}

			err = errors.New(msg)
		}
	} else {
		vfs, err = iofs.Sub(vfs, path.Clean(subDir))
	}

	return &assetFileServer{fs: vfs, err: err}
}

func (d *assetFileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	url := req.URL.Path

	err := d.err
	if err == nil {
		filename := path.Clean(strings.TrimLeft(url, "/"))
		d.logInfo(ctx, "Handling request", "url", url, "file", filename)
		err = d.serveFSFile(rw, req, filename)
		if os.IsNotExist(err) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}

	if err != nil {
		d.logError(ctx, "Unable to handle request", "url", url, "err", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// serveFile will try to load the file from the fs.FS and write it to the response
func (d *assetFileServer) serveFSFile(rw http.ResponseWriter, req *http.Request, filename string) error {
	if d.fs == nil {
		return os.ErrNotExist
	}

	file, err := d.fs.Open(filename)
	if err != nil {
		if s := path.Ext(filename); s == "" {
			filename = filename + ".html"
			file, err = d.fs.Open(filename)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer file.Close()

	statInfo, err := file.Stat()
	if err != nil {
		return err
	}

	url := req.URL.Path
	isDirectoryPath := url == "" || url[len(url)-1] == '/'
	if statInfo.IsDir() {
		if !isDirectoryPath {
			// If the URL doesn't end in a slash normally a http.redirect should be done, but that currently doesn't work on
			// WebKit WebViews (macOS/Linux).
			// So we handle this as a specific error
			return fmt.Errorf("a directory has been requested without a trailing slash, please add a trailing slash to your request")
		}

		filename = path.Join(filename, indexHTML)

		file, err = d.fs.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		statInfo, err = file.Stat()
		if err != nil {
			return err
		}
	} else if isDirectoryPath {
		return fmt.Errorf("a file has been requested with a trailing slash, please remove the trailing slash from your request")
	}

	var prefix []byte
	if _, haveType := rw.Header()[HeaderContentType]; !haveType {
		contentType := mimetypeByExtension(filename)
		if contentType == "" {
			// Only unknown extensions require reading a MIME sniffing prefix.
			prefix = make([]byte, 512)
			n, err := io.ReadFull(file, prefix)
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return err
			}
			prefix = prefix[:n]
			contentType = http.DetectContentType(prefix)
		}
		rw.Header().Set(HeaderContentType, contentType)
	}

	if fileSeeker, _ := file.(io.ReadSeeker); fileSeeker != nil {
		if _, err := fileSeeker.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("seeker can't seek")
		}

		http.ServeContent(rw, req, statInfo.Name(), statInfo.ModTime(), fileSeeker)
		return nil
	}

	rw.Header().Set(HeaderContentLength, fmt.Sprintf("%d", statInfo.Size()))
	if req.Method == http.MethodHead {
		rw.WriteHeader(http.StatusOK)
		return nil
	}

	// Write the first 512 bytes used for MimeType sniffing
	if len(prefix) > 0 {
		written, err := rw.Write(prefix)
		if err != nil {
			return err
		}
		if written != len(prefix) {
			return io.ErrShortWrite
		}
	}

	// Copy the remaining content of the file
	buffer := assetFileCopyBuffers.Get().(*[assetFileCopyBufferSize]byte)
	defer assetFileCopyBuffers.Put(buffer)
	_, err = io.CopyBuffer(rw, file, buffer[:])
	if err == nil {
		rw.WriteHeader(http.StatusOK)
	}
	return err
}

func (d *assetFileServer) logInfo(ctx context.Context, message string, args ...interface{}) {
	logInfo(ctx, "[AssetFileServerFS] "+message, args...)
}

func (d *assetFileServer) logError(ctx context.Context, message string, args ...interface{}) {
	logError(ctx, "[AssetFileServerFS] "+message, args...)
}
