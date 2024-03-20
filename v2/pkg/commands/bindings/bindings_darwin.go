//go:build darwin

package bindings

import "github.com/wailsapp/wails/v2/internal/xattr"

func init() {
	fixupXattrs = xattr.RemoveXAttr
}
