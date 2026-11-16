# Read-only application filesystem

This package writes SquashFS v4 images with standard-library zlib compression.
The writer is original Wails code based on the documented disk format; no
third-party implementation was copied or adapted and no module was added.

`Write` starts at the current output position and uses offsets relative to that
position, allowing an AppImage runtime to precede the filesystem. The caller owns
publication and cleanup of incomplete outputs. Input entries must remain stable
while writing; a reader that produces more or less than its declared size fails.

Supported entries are regular files, directories and symbolic links. Paths are
slash-relative, missing parents become directories with mode 0755, and explicit
permissions are preserved independently of the host operating system. Ownership
is root:root. Entries are sorted to produce identical bytes independent of input
ordering. A zero timestamp means the Unix epoch. Regular files are streamed one
block at a time; directory and inode metadata stays in memory.

Data blocks default to 128 KiB and fall back to uncompressed storage when zlib
would make them larger. Metadata uses uncompressed 8 KiB frames. Extended directory
inodes support listings larger than 64 KiB without requiring directory indexes.
Fragments, xattrs, devices, hard-link deduplication and sparse-file accounting are
not emitted. Individual regular files and their starting data offsets must fit
32 bits. The writer rejects more than 100,000 entries, paths over 4096 bytes,
filenames over 256 bytes, and estimated metadata over 64 MiB.

## Verification and measurement

```sh
go test -race ./internal/readonlyfs
go test ./internal/readonlyfs -run '^$' -bench BenchmarkWrite -benchmem
```

When `unsquashfs` is installed, tests independently extract both standalone and
runtime-prefixed images. The 3,206-entry fixture crosses inode/directory metadata
boundaries and exceeds the basic directory listing limit. It checks every file's
bytes and permissions, symlinks, empty files and empty directories. Other tests
cover deterministic ordering, invalid paths/types, changing source sizes and
cancellation. On 2026-11-01 this passed with unsquashfs 4.6.1.

Three 200 ms benchmark samples on Go 1.26.10, Linux/amd64, Ryzen 3950X and
GOMAXPROCS=4 measured a median of 74.669 ms (56.17 MB/s) for a 4 MiB mixed-
compressibility payload. Allocation volume was about 1.21 MB/op; this is not peak
RSS. The CPU profile attributed about 96% to standard-library deflate work.

A separate seven-sample alternating comparison with the already-installed
`klauspost/compress/zlib` measured 71.713 ms for stdlib versus 51.534 ms for that
backend, with output growing from 2,231,376 to 2,328,211 bytes. The writer retains
stdlib to avoid another direct dependency relationship and preserve the smaller
output. Only results of that external comparison were retained; its dependency-
using harness was removed. These are bounded codec/writer measurements, not a
whole-AppImage comparison. Custom assembly was not introduced.

## Format references

- [SquashFS on-disk format documentation](https://github.com/AgentD/squashfs-tools-ng/blob/master/doc/format.adoc)
- [Linux SquashFS filesystem documentation](https://docs.kernel.org/filesystems/squashfs.html)
