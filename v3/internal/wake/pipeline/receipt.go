package pipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const ArtifactReceiptVersion = 1

type ArtifactReceipt struct {
	Version   int              `json:"version"`
	Artifacts []ArtifactRecord `json:"artifacts"`
}

type ArtifactRecord struct {
	Producer  NodeKey      `json:"producer"`
	Path      string       `json:"path"`
	Digest    string       `json:"digest"`
	Size      int64        `json:"size"`
	Kind      ArtifactKind `json:"kind"`
	Target    string       `json:"target"`
	Format    string       `json:"format,omitempty"`
	Signed    bool         `json:"signed"`
	Notarized bool         `json:"notarized"`
}

func WriteArtifactReceipt(root, output string, artifacts []ArtifactReference) (ArtifactReceipt, error) {
	return writeArtifactReceiptWithOperations(root, output, artifacts, artifactReceiptOperations{
		digest:  digestArtifact,
		marshal: json.MarshalIndent,
		write:   writeArtifactReceiptAtomic,
	})
}

type artifactReceiptOperations struct {
	digest  func(string) (string, int64, error)
	marshal func(any, string, string) ([]byte, error)
	write   func(string, []byte) error
}

func writeArtifactReceiptWithOperations(root, output string, artifacts []ArtifactReference, operations artifactReceiptOperations) (ArtifactReceipt, error) {
	receipt := ArtifactReceipt{Version: ArtifactReceiptVersion, Artifacts: make([]ArtifactRecord, 0, len(artifacts))}
	for _, artifact := range artifacts {
		path, err := artifactPath(root, artifact.Path)
		if err != nil {
			return ArtifactReceipt{}, fmt.Errorf("artifact %s: %w", artifact.Key, err)
		}
		digest, size, err := operations.digest(path)
		if err != nil {
			return ArtifactReceipt{}, fmt.Errorf("artifact %s: %w", artifact.Key, err)
		}
		receipt.Artifacts = append(receipt.Artifacts, ArtifactRecord{
			Producer: artifact.Key, Path: filepath.ToSlash(filepath.Clean(artifact.Path)), Digest: digest, Size: size,
			Kind: artifact.Identity.Kind, Target: artifact.Identity.Target.OS + "/" + artifact.Identity.Target.Arch,
			Format: artifact.Identity.Format, Signed: artifact.Identity.Signed, Notarized: artifact.Identity.Notarized,
		})
	}
	sort.Slice(receipt.Artifacts, func(i, j int) bool { return receipt.Artifacts[i].Producer < receipt.Artifacts[j].Producer })
	data, err := operations.marshal(receipt, "", "  ")
	if err != nil {
		return ArtifactReceipt{}, err
	}
	data = append(data, '\n')
	destination, err := artifactPath(root, output)
	if err != nil {
		return ArtifactReceipt{}, fmt.Errorf("artifact receipt: %w", err)
	}
	if err := operations.write(destination, data); err != nil {
		return ArtifactReceipt{}, fmt.Errorf("artifact receipt: %w", err)
	}
	return receipt, nil
}

func artifactPath(root, relative string) (string, error) {
	if relative == "" {
		return "", fmt.Errorf("path is empty")
	}
	clean := filepath.Clean(filepath.FromSlash(relative))
	if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes the project", relative)
	}
	return filepath.Join(root, clean), nil
}

func digestArtifact(path string) (string, int64, error) {
	return digestArtifactWithOperations(path, artifactDigestOperations{
		lstat:    os.Lstat,
		open:     func(path string) (io.ReadCloser, error) { return os.Open(path) },
		walkDir:  filepath.WalkDir,
		rel:      filepath.Rel,
		readlink: os.Readlink,
	})
}

type artifactDigestOperations struct {
	lstat    func(string) (fs.FileInfo, error)
	open     func(string) (io.ReadCloser, error)
	walkDir  func(string, fs.WalkDirFunc) error
	rel      func(string, string) (string, error)
	readlink func(string) (string, error)
}

func digestArtifactWithOperations(path string, operations artifactDigestOperations) (string, int64, error) {
	info, err := operations.lstat(path)
	if err != nil {
		return "", 0, err
	}
	if info.Mode().IsRegular() {
		file, err := operations.open(path)
		if err != nil {
			return "", 0, err
		}
		hash := sha256.New()
		_, _ = fmt.Fprintf(hash, "%o\x00", info.Mode().Perm())
		if _, err := io.Copy(hash, file); err != nil {
			_ = file.Close()
			return "", 0, err
		}
		if err := file.Close(); err != nil {
			return "", 0, err
		}
		return "sha256:" + hex.EncodeToString(hash.Sum(nil)), info.Size(), nil
	}
	if !info.IsDir() {
		return "", 0, fmt.Errorf("unsupported artifact type %s", info.Mode().Type())
	}
	hash := sha256.New()
	var size int64
	err = operations.walkDir(path, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := operations.rel(path, current)
		if err != nil {
			return err
		}
		if relative == "." {
			return nil
		}
		relative = filepath.ToSlash(relative)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(hash, "%s\x00%s\x00%o\x00", relative, info.Mode().Type(), info.Mode().Perm())
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := operations.readlink(current)
			if err != nil {
				return err
			}
			size += int64(len(target))
			_, _ = io.WriteString(hash, target)
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		file, err := operations.open(current)
		if err != nil {
			return err
		}
		written, copyErr := io.Copy(hash, file)
		closeErr := file.Close()
		size += written
		if copyErr != nil {
			return copyErr
		}
		return closeErr
	})
	if err != nil {
		return "", 0, err
	}
	return "sha256:" + hex.EncodeToString(hash.Sum(nil)), size, nil
}

func writeArtifactReceiptAtomic(destination string, data []byte) error {
	return writeArtifactReceiptAtomicWithOperations(destination, data, receiptWriteOperations{
		mkdirAll: os.MkdirAll,
		createTemp: func(directory, pattern string) (receiptTemporaryFile, error) {
			return os.CreateTemp(directory, pattern)
		},
		remove: os.Remove,
		rename: os.Rename,
	})
}

type receiptTemporaryFile interface {
	Name() string
	Chmod(fs.FileMode) error
	Write([]byte) (int, error)
	Sync() error
	Close() error
}

type receiptWriteOperations struct {
	mkdirAll   func(string, fs.FileMode) error
	createTemp func(string, string) (receiptTemporaryFile, error)
	remove     func(string) error
	rename     func(string, string) error
}

func writeArtifactReceiptAtomicWithOperations(destination string, data []byte, operations receiptWriteOperations) error {
	directory := filepath.Dir(destination)
	if err := operations.mkdirAll(directory, 0o755); err != nil {
		return err
	}
	temporary, err := operations.createTemp(directory, ".receipt-*")
	if err != nil {
		return err
	}
	temporaryName := temporary.Name()
	committed := false
	defer func() {
		_ = temporary.Close()
		if !committed {
			_ = operations.remove(temporaryName)
		}
	}()
	if err := temporary.Chmod(0o644); err != nil {
		return err
	}
	if written, err := temporary.Write(data); err != nil {
		return err
	} else if written != len(data) {
		return io.ErrShortWrite
	}
	if err := temporary.Sync(); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := operations.rename(temporaryName, destination); err != nil {
		return err
	}
	committed = true
	return nil
}
