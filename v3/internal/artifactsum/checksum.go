// Package artifactsum writes SHA256 sidecars for completed distribution files.
package artifactsum

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Write publishes artifact.sha256 in sha256sum format. Hashing uses bounded
// memory and observes cancellation; errors preserve an existing sidecar.
// Callers must finish writing the artifact before invoking Write.
func Write(ctx context.Context, artifact string) error {
	digest, err := checksum(ctx, artifact)
	if err != nil {
		return err
	}
	name := filepath.Base(artifact)
	prefix := ""
	if strings.ContainsAny(name, "\\\n\r") {
		prefix = "\\"
		name = strings.NewReplacer("\\", "\\\\", "\n", "\\n", "\r", "\\r").Replace(name)
	}
	output, err := os.CreateTemp(filepath.Dir(artifact), ".artifact-sha256-*")
	if err != nil {
		return err
	}
	defer os.Remove(output.Name())
	_, writeErr := fmt.Fprintf(output, "%s%x  %s\n", prefix, digest, name)
	if writeErr == nil {
		writeErr = output.Chmod(0644)
	}
	closeErr := output.Close()
	if writeErr != nil {
		return writeErr
	}
	if closeErr != nil {
		return closeErr
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return os.Rename(output.Name(), artifact+".sha256")
}

func checksum(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("checksum input is not a regular file: %s", path)
	}
	input, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer input.Close()
	hash := sha256.New()
	buffer := make([]byte, 128*1024)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		n, err := input.Read(buffer)
		_, _ = hash.Write(buffer[:n])
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return hash.Sum(nil), nil
}
