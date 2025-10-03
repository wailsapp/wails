package conv

import (
	"unsafe"
)

// BytesToString converts []byte to string without memory allocation.
// WARNING: The returned string must be used as read-only!
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
