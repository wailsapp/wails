//go:build linux
// +build linux

package linux

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Calloc handles alloc/dealloc of C data
type Calloc struct {
	pool []unsafe.Pointer
}

// NewCalloc creates a new allocator
func NewCalloc() Calloc {
	return Calloc{}
}

// String creates a new C string and retains a reference to it
func (c Calloc) String(in string) *C.char {
	result := C.CString(in)
	c.pool = append(c.pool, unsafe.Pointer(result))
	return result
}

// Free frees all allocated C memory
func (c Calloc) Free() {
	for _, str := range c.pool {
		C.free(str)
	}
	c.pool = []unsafe.Pointer{}
}
