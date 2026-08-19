package pipeline

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArtifactReceiptIsTypedDeterministicAndVerifiable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated rights on Windows")
	}
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "bin", "App.app", "Contents"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "bin", "app"), []byte("binary"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "bin", "App.app", "Contents", "app"), []byte("bundle"), 0o755))
	require.NoError(t, os.Symlink("Contents/app", filepath.Join(root, "bin", "App.app", "current")))
	artifacts := []ArtifactReference{
		{Key: "package", Path: "bin/App.app", Identity: ArtifactIdentity{Kind: ArtifactBundle, Target: Target{OS: "darwin", Arch: "arm64"}, Format: "app", Signed: true, Notarized: true}},
		{Key: "binary", Path: "bin/app", Identity: ArtifactIdentity{Kind: ArtifactBinary, Target: Target{OS: "linux", Arch: "amd64"}}},
	}

	first, err := WriteArtifactReceipt(root, ".wails/artifacts/receipt.json", artifacts)
	require.NoError(t, err)
	firstData, err := os.ReadFile(filepath.Join(root, ".wails", "artifacts", "receipt.json"))
	require.NoError(t, err)
	second, err := WriteArtifactReceipt(root, ".wails/artifacts/receipt.json", artifacts)
	require.NoError(t, err)
	secondData, err := os.ReadFile(filepath.Join(root, ".wails", "artifacts", "receipt.json"))
	require.NoError(t, err)
	assert.Equal(t, first, second)
	assert.Equal(t, firstData, secondData)
	require.Len(t, first.Artifacts, 2)
	assert.Equal(t, NodeKey("binary"), first.Artifacts[0].Producer)
	assert.Equal(t, "linux/amd64", first.Artifacts[0].Target)
	assert.Equal(t, int64(6), first.Artifacts[0].Size)
	assert.Regexp(t, `^sha256:[0-9a-f]{64}$`, first.Artifacts[0].Digest)
	assert.Equal(t, NodeKey("package"), first.Artifacts[1].Producer)
	assert.Equal(t, int64(len("bundle")+len("Contents/app")), first.Artifacts[1].Size)
	assert.True(t, first.Artifacts[1].Signed)
	assert.True(t, first.Artifacts[1].Notarized)
	var decoded ArtifactReceipt
	require.NoError(t, json.Unmarshal(firstData, &decoded))
	assert.Equal(t, first, decoded)
}

func TestArtifactReceiptRejectsEscapesMissingInputsAndUnsupportedTypes(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated rights on Windows")
	}
	root := t.TempDir()
	_, err := WriteArtifactReceipt(root, "../receipt.json", nil)
	assert.ErrorContains(t, err, "escapes the project")
	_, err = WriteArtifactReceipt(root, ".wails/receipt.json", []ArtifactReference{{Key: "missing", Path: "bin/missing"}})
	assert.ErrorContains(t, err, "artifact missing")
	require.NoError(t, os.MkdirAll(filepath.Join(root, "bin"), 0o755))
	require.NoError(t, os.Symlink("missing", filepath.Join(root, "bin", "link")))
	_, err = WriteArtifactReceipt(root, ".wails/receipt.json", []ArtifactReference{{Key: "link", Path: "bin/link"}})
	assert.ErrorContains(t, err, "unsupported artifact type")
}

func TestArtifactIdentityDisplayKind(t *testing.T) {
	assert.True(t, (ArtifactIdentity{}).Empty())
	assert.Equal(t, "binary", (ArtifactIdentity{Kind: ArtifactBinary}).DisplayKind())
	assert.Equal(t, "deb", (ArtifactIdentity{Kind: ArtifactPackage, Format: "deb"}).DisplayKind())
}

func TestArtifactReceiptAtomicWriteFaultsNeverPublishPartialData(t *testing.T) {
	want := errors.New("injected receipt failure")
	newFile := func() *receiptFaultFile { return &receiptFaultFile{name: "/output/.receipt-temp"} }
	newOperations := func(file *receiptFaultFile) (receiptWriteOperations, *receiptWriteRecorder) {
		recorder := &receiptWriteRecorder{}
		return receiptWriteOperations{
			mkdirAll:   func(string, fs.FileMode) error { return nil },
			createTemp: func(string, string) (receiptTemporaryFile, error) { return file, nil },
			remove:     func(path string) error { recorder.removed = append(recorder.removed, path); return nil },
			rename: func(source, destination string) error {
				recorder.renamed = [2]string{source, destination}
				return nil
			},
		}, recorder
	}

	operations, _ := newOperations(newFile())
	operations.mkdirAll = func(string, fs.FileMode) error { return want }
	require.ErrorIs(t, writeArtifactReceiptAtomicWithOperations("/output/receipt.json", []byte("data"), operations), want)

	operations, _ = newOperations(newFile())
	operations.createTemp = func(string, string) (receiptTemporaryFile, error) { return nil, want }
	require.ErrorIs(t, writeArtifactReceiptAtomicWithOperations("/output/receipt.json", []byte("data"), operations), want)

	for name, configure := range map[string]func(*receiptFaultFile){
		"chmod":       func(file *receiptFaultFile) { file.chmodErr = want },
		"write":       func(file *receiptFaultFile) { file.writeErr = want },
		"short write": func(file *receiptFaultFile) { file.shortWrite = true },
		"sync":        func(file *receiptFaultFile) { file.syncErr = want },
		"close":       func(file *receiptFaultFile) { file.closeErr = want },
	} {
		t.Run(name, func(t *testing.T) {
			file := newFile()
			configure(file)
			operations, recorder := newOperations(file)
			err := writeArtifactReceiptAtomicWithOperations("/output/receipt.json", []byte("data"), operations)
			require.Error(t, err)
			if name == "short write" {
				require.ErrorIs(t, err, io.ErrShortWrite)
			} else {
				require.ErrorIs(t, err, want)
			}
			assert.Equal(t, []string{file.name}, recorder.removed)
			assert.Equal(t, [2]string{}, recorder.renamed)
		})
	}

	file := newFile()
	operations, recorder := newOperations(file)
	operations.rename = func(string, string) error { return want }
	require.ErrorIs(t, writeArtifactReceiptAtomicWithOperations("/output/receipt.json", []byte("data"), operations), want)
	assert.Equal(t, []string{file.name}, recorder.removed)

	file = newFile()
	operations, recorder = newOperations(file)
	require.NoError(t, writeArtifactReceiptAtomicWithOperations("/output/receipt.json", []byte("data"), operations))
	assert.Equal(t, [2]string{file.name, "/output/receipt.json"}, recorder.renamed)
	assert.Empty(t, recorder.removed)
}

func TestArtifactReceiptCompositionPropagatesEveryFailure(t *testing.T) {
	want := errors.New("injected receipt composition failure")
	base := artifactReceiptOperations{
		digest:  func(string) (string, int64, error) { return "sha256:digest", 4, nil },
		marshal: json.MarshalIndent,
		write:   func(string, []byte) error { return nil },
	}
	artifact := ArtifactReference{Key: "binary", Path: "bin/app", Identity: ArtifactIdentity{Kind: ArtifactBinary, Target: Target{OS: "linux", Arch: "amd64"}}}

	_, err := writeArtifactReceiptWithOperations("/project", "receipt.json", []ArtifactReference{{Key: "escape", Path: "../outside"}}, base)
	require.ErrorContains(t, err, "artifact escape")

	operations := base
	operations.digest = func(string) (string, int64, error) { return "", 0, want }
	_, err = writeArtifactReceiptWithOperations("/project", "receipt.json", []ArtifactReference{artifact}, operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.marshal = func(any, string, string) ([]byte, error) { return nil, want }
	_, err = writeArtifactReceiptWithOperations("/project", "receipt.json", []ArtifactReference{artifact}, operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.write = func(string, []byte) error { return want }
	_, err = writeArtifactReceiptWithOperations("/project", "receipt.json", []ArtifactReference{artifact}, operations)
	require.ErrorIs(t, err, want)

	_, err = writeArtifactReceiptWithOperations("/project", "../receipt.json", []ArtifactReference{artifact}, base)
	require.ErrorContains(t, err, "artifact receipt")

	_, err = artifactPath("/project", "")
	require.ErrorContains(t, err, "empty")
}

func TestArtifactDigestPropagatesFilesystemAndStreamFaults(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated rights on Windows")
	}
	want := errors.New("injected artifact digest failure")
	root := t.TempDir()
	directoryInfo, err := os.Stat(root)
	require.NoError(t, err)
	regularPath := filepath.Join(root, "regular")
	require.NoError(t, os.WriteFile(regularPath, []byte("data"), 0o644))
	regularInfo, err := os.Stat(regularPath)
	require.NoError(t, err)
	symlinkPath := filepath.Join(root, "link")
	require.NoError(t, os.Symlink("regular", symlinkPath))
	symlinkInfo, err := os.Lstat(symlinkPath)
	require.NoError(t, err)

	base := artifactDigestOperations{
		lstat: func(string) (fs.FileInfo, error) { return directoryInfo, nil },
		open: func(string) (io.ReadCloser, error) {
			return &receiptReadCloser{Reader: bytes.NewReader([]byte("data"))}, nil
		},
		walkDir:  func(string, fs.WalkDirFunc) error { return nil },
		rel:      filepath.Rel,
		readlink: func(string) (string, error) { return "regular", nil },
	}

	operations := base
	operations.lstat = func(string) (fs.FileInfo, error) { return nil, want }
	_, _, err = digestArtifactWithOperations(root, operations)
	require.ErrorIs(t, err, want)

	for name, reader := range map[string]*receiptReadCloser{
		"copy":  {Reader: receiptErrorReader{err: want}},
		"close": {Reader: bytes.NewReader([]byte("data")), closeErr: want},
	} {
		t.Run("regular-"+name, func(t *testing.T) {
			operations := base
			operations.lstat = func(string) (fs.FileInfo, error) { return regularInfo, nil }
			operations.open = func(string) (io.ReadCloser, error) { return reader, nil }
			_, _, err := digestArtifactWithOperations(regularPath, operations)
			require.ErrorIs(t, err, want)
		})
	}
	operations = base
	operations.lstat = func(string) (fs.FileInfo, error) { return regularInfo, nil }
	operations.open = func(string) (io.ReadCloser, error) { return nil, want }
	_, _, err = digestArtifactWithOperations(regularPath, operations)
	require.ErrorIs(t, err, want)

	operations = base
	operations.walkDir = func(string, fs.WalkDirFunc) error { return want }
	_, _, err = digestArtifactWithOperations(root, operations)
	require.ErrorIs(t, err, want)

	for name, invoke := range map[string]func(fs.WalkDirFunc) error{
		"walk callback": func(callback fs.WalkDirFunc) error { return callback(root, nil, want) },
		"relative path": func(callback fs.WalkDirFunc) error {
			return callback(regularPath, fs.FileInfoToDirEntry(regularInfo), nil)
		},
		"entry info": func(callback fs.WalkDirFunc) error {
			return callback(regularPath, receiptDirEntry{infoErr: want}, nil)
		},
		"readlink": func(callback fs.WalkDirFunc) error {
			return callback(symlinkPath, fs.FileInfoToDirEntry(symlinkInfo), nil)
		},
		"open": func(callback fs.WalkDirFunc) error {
			return callback(regularPath, fs.FileInfoToDirEntry(regularInfo), nil)
		},
		"copy": func(callback fs.WalkDirFunc) error {
			return callback(regularPath, fs.FileInfoToDirEntry(regularInfo), nil)
		},
		"close": func(callback fs.WalkDirFunc) error {
			return callback(regularPath, fs.FileInfoToDirEntry(regularInfo), nil)
		},
	} {
		t.Run("directory-"+name, func(t *testing.T) {
			operations := base
			operations.walkDir = func(_ string, callback fs.WalkDirFunc) error { return invoke(callback) }
			switch name {
			case "relative path":
				operations.rel = func(string, string) (string, error) { return "", want }
			case "readlink":
				operations.readlink = func(string) (string, error) { return "", want }
			case "open":
				operations.open = func(string) (io.ReadCloser, error) { return nil, want }
			case "copy":
				operations.open = func(string) (io.ReadCloser, error) {
					return &receiptReadCloser{Reader: receiptErrorReader{err: want}}, nil
				}
			case "close":
				operations.open = func(string) (io.ReadCloser, error) {
					return &receiptReadCloser{Reader: bytes.NewReader(nil), closeErr: want}, nil
				}
			}
			_, _, err := digestArtifactWithOperations(root, operations)
			require.ErrorIs(t, err, want)
		})
	}
}

type receiptFaultFile struct {
	name                                  string
	chmodErr, writeErr, syncErr, closeErr error
	shortWrite                            bool
}

func (f *receiptFaultFile) Name() string            { return f.name }
func (f *receiptFaultFile) Chmod(fs.FileMode) error { return f.chmodErr }
func (f *receiptFaultFile) Write(data []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	if f.shortWrite {
		return len(data) - 1, nil
	}
	return len(data), nil
}
func (f *receiptFaultFile) Sync() error  { return f.syncErr }
func (f *receiptFaultFile) Close() error { return f.closeErr }

type receiptWriteRecorder struct {
	removed []string
	renamed [2]string
}

type receiptReadCloser struct {
	io.Reader
	closeErr error
}

func (r *receiptReadCloser) Close() error { return r.closeErr }

type receiptErrorReader struct{ err error }

func (r receiptErrorReader) Read([]byte) (int, error) { return 0, r.err }

type receiptDirEntry struct{ infoErr error }

func (receiptDirEntry) Name() string                     { return "entry" }
func (receiptDirEntry) IsDir() bool                      { return false }
func (receiptDirEntry) Type() fs.FileMode                { return 0 }
func (entry receiptDirEntry) Info() (fs.FileInfo, error) { return nil, entry.infoErr }

func BenchmarkWriteArtifactReceipt(b *testing.B) {
	root := b.TempDir()
	require.NoError(b, os.MkdirAll(filepath.Join(root, "bin"), 0o755))
	artifacts := make([]ArtifactReference, 100)
	for index := range artifacts {
		path := filepath.ToSlash(filepath.Join("bin", fmt.Sprintf("artifact-%03d", index)))
		require.NoError(b, os.WriteFile(filepath.Join(root, filepath.FromSlash(path)), bytes.Repeat([]byte{byte(index)}, 4096), 0o644))
		artifacts[index] = ArtifactReference{Key: NodeKey(fmt.Sprintf("target:%03d", index)), Path: path, Identity: ArtifactIdentity{Kind: ArtifactBinary, Target: Target{OS: "linux", Arch: "amd64"}}}
	}
	b.ReportAllocs()
	b.SetBytes(100 * 4096)
	b.ResetTimer()
	for range b.N {
		if _, err := WriteArtifactReceipt(root, ".wails/artifacts/receipt.json", artifacts); err != nil {
			b.Fatal(err)
		}
	}
}
