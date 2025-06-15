package assetserver

import (
	"sync"
)

// Buffer sizes for different use cases
const (
	// ContentSnifferBufferSize is the standard size for content type detection
	ContentSnifferBufferSize = 512
	// AssetBufferSize is used for general asset serving operations  
	AssetBufferSize = 4096
	// LargeBufferSize for larger operations like image processing
	LargeBufferSize = 32768
)

// Global buffer pools for different sizes
var (
	contentSnifferPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, ContentSnifferBufferSize)
		},
	}
	
	assetBufferPool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 0, AssetBufferSize)
			return &buf
		},
	}
	
	largeBufferPool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 0, LargeBufferSize)
			return &buf
		},
	}
)

// GetContentSnifferBuffer returns a pooled buffer for content type sniffing
// The returned buffer has zero length but ContentSnifferBufferSize capacity
func GetContentSnifferBuffer() []byte {
	buf := contentSnifferPool.Get().([]byte)
	return buf[:0] // Reset length to 0, keep capacity
}

// PutContentSnifferBuffer returns a buffer to the content sniffer pool
// The buffer is reset to zero length before being returned to the pool
func PutContentSnifferBuffer(buf []byte) {
	if cap(buf) < ContentSnifferBufferSize {
		// Don't pool buffers that are too small
		return
	}
	
	// Reset the buffer and return to pool
	contentSnifferPool.Put(buf[:0])
}

// GetAssetBuffer returns a pooled buffer for general asset operations
func GetAssetBuffer() []byte {
	bufPtr := assetBufferPool.Get().(*[]byte)
	buf := *bufPtr
	return buf[:0]
}

// PutAssetBuffer returns a buffer to the asset buffer pool
func PutAssetBuffer(buf []byte) {
	if cap(buf) != AssetBufferSize {
		return
	}
	
	buf = buf[:0]
	assetBufferPool.Put(&buf)
}

// GetLargeBuffer returns a pooled buffer for large operations
func GetLargeBuffer() []byte {
	bufPtr := largeBufferPool.Get().(*[]byte)
	buf := *bufPtr
	return buf[:0]
}

// PutLargeBuffer returns a buffer to the large buffer pool
func PutLargeBuffer(buf []byte) {
	if cap(buf) != LargeBufferSize {
		return
	}
	
	buf = buf[:0]
	largeBufferPool.Put(&buf)
}

// Pool statistics for monitoring and debugging
type PoolStats struct {
	ContentSnifferHits   int64
	ContentSnifferMisses int64
	AssetBufferHits      int64
	AssetBufferMisses    int64
	LargeBufferHits      int64
	LargeBufferMisses    int64
}

// GetPoolStats returns current pool usage statistics
// Note: This is approximate and for debugging purposes only
func GetPoolStats() PoolStats {
	return PoolStats{
		// sync.Pool doesn't expose hit/miss stats directly
		// This is a placeholder for potential future monitoring
	}
}

// ResetPools clears all buffer pools (useful for testing)
func ResetPools() {
	contentSnifferPool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 0, ContentSnifferBufferSize)
			return &buf
		},
	}
	
	assetBufferPool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 0, AssetBufferSize)
			return &buf
		},
	}
	
	largeBufferPool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 0, LargeBufferSize)
			return &buf
		},
	}
}