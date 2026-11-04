package appimage

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"slices"
	"strconv"
)

// RuntimeIdentity fingerprints the exact selected deployment, including data
// added since a previous build. Call it after executable inputs materialize.
func (options Options) RuntimeIdentity(ctx context.Context) (string, error) {
	if len(options.Binaries) == 0 {
		return "", fmt.Errorf("runtime identity requires executable inputs")
	}
	deployment, err := options.collectRuntime(ctx)
	if err != nil {
		return "", err
	}
	runtimePath, err := options.PrepareRuntime(ctx)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	frame := func(value string) { writeRuntimePart(hash, []byte(value)) }
	for _, value := range []string{"native-appimage-v1", deployment.root, options.Arch, strconv.Itoa(deployment.gtk)} {
		frame(value)
	}
	digests := map[string][]byte{}
	fingerprint := func(filename string) error {
		frame(filename)
		digest, ok := digests[filename]
		if !ok {
			file, err := os.Open(filename)
			if err != nil {
				return err
			}
			info, err := file.Stat()
			if err != nil {
				file.Close()
				return err
			}
			if !info.Mode().IsRegular() {
				file.Close()
				return &os.PathError{Op: "hash runtime", Path: filename, Err: os.ErrInvalid}
			}
			content := sha256.New()
			writeRuntimePart(content, []byte(info.Mode().String()))
			_, err = io.Copy(contextWriter{ctx, content}, file)
			closeErr := file.Close()
			if err != nil {
				return err
			}
			if closeErr != nil {
				return closeErr
			}
			digest = content.Sum(nil)
			digests[filename] = digest
		}
		writeRuntimePart(hash, digest)
		return nil
	}
	if err := fingerprint(runtimePath); err != nil {
		return "", err
	}
	if options.Tools.AppRun != "" {
		if err := fingerprint(options.path(options.Tools.AppRun)); err != nil {
			return "", err
		}
	} else {
		writeRuntimePart(hash, nativeLauncher(options.Binaries[0].Name))
	}
	destinations := make([]string, 0, len(deployment.files))
	for destination := range deployment.files {
		destinations = append(destinations, destination)
	}
	slices.Sort(destinations)
	for _, destination := range destinations {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		frame(destination)
		if err := fingerprint(deployment.files[destination]); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func writeRuntimePart(writer hash.Hash, value []byte) {
	var length [8]byte
	binary.BigEndian.PutUint64(length[:], uint64(len(value)))
	writer.Write(length[:])
	writer.Write(value)
}
