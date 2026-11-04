package appimage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func ensureTool(ctx context.Context, client *http.Client, url, filename string, force bool) (err error) {
	if err = ctx.Err(); err != nil {
		return err
	}
	if !force {
		if info, err := os.Stat(filename); err == nil && info.Mode().IsRegular() && info.Size() > 0 {
			return os.Chmod(filename, 0755)
		}
	}
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Minute}
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: HTTP %s", url, response.Status)
	}
	if err = os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	output, err := os.CreateTemp(filepath.Dir(filename), ".tool-*")
	if err != nil {
		return err
	}
	defer os.Remove(output.Name())
	defer output.Close()
	count, err := io.Copy(contextWriter{ctx, output}, response.Body)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("download %s was empty", url)
	}
	if err = output.Chmod(0755); err != nil {
		return err
	}
	if err = output.Close(); err != nil {
		return err
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	return os.Rename(output.Name(), filename)
}
