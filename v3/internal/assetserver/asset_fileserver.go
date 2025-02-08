package assetserver

import (
	"bytes"
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
)

const (
	indexHTML = "index.html"
)

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
		filename := path.Clean(strings.TrimPrefix(url, "/"))
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
		return err
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

	var buf [512]byte
	var n int
	if _, haveType := rw.Header()[HeaderContentType]; !haveType {
		// Detect MimeType by sniffing the first 512 bytes
		n, err = file.Read(buf[:])
		if err != nil && err != io.EOF {
			return err
		}

		// Do the custom MimeType sniffing even though http.ServeContent would do it in case
		// of an io.ReadSeeker. We would like to have a consistent behaviour in both cases.
		if contentType := GetMimetype(filename, buf[:n]); contentType != "" {
			rw.Header().Set(HeaderContentType, contentType)
		}
	}

	if fileSeeker, _ := file.(io.ReadSeeker); fileSeeker != nil {
		if _, err := fileSeeker.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("seeker can't seek")
		}

		http.ServeContent(rw, req, statInfo.Name(), statInfo.ModTime(), fileSeeker)
		return nil
	}

	rw.Header().Set(HeaderContentLength, fmt.Sprintf("%d", statInfo.Size()))

	// Write the first 512 bytes used for MimeType sniffing
	_, err = io.Copy(rw, bytes.NewReader(buf[:n]))
	if err != nil {
		return err
	}

	// Copy the remaining content of the file
	_, err = io.Copy(rw, file)
	return err
}

func (d *assetFileServer) logInfo(ctx context.Context, message string, args ...interface{}) {
	logInfo(ctx, "[AssetFileServerFS] "+message, args...)
}

func (d *assetFileServer) logError(ctx context.Context, message string, args ...interface{}) {
	logError(ctx, "[AssetFileServerFS] "+message, args...)
}
