//go:build darwin

package build

import "github.com/wailsapp/wails/v2/internal/xattr"

func init() {
	fixupXattrs = xattr.RemoveXAttr
}
