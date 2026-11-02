package iconfile

import "image/png"

// pngBuffers keeps one compression workspace for the sizes of a single icon.
// Each encoding call owns its pool and uses it sequentially.
type pngBuffers struct{ buffer *png.EncoderBuffer }

func (pool *pngBuffers) Get() *png.EncoderBuffer {
	buffer := pool.buffer
	pool.buffer = nil
	return buffer
}

func (pool *pngBuffers) Put(buffer *png.EncoderBuffer) {
	pool.buffer = buffer
}
