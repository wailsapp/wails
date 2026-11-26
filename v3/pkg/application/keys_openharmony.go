//go:build openharmony && !server

package application

// OHOS keyboard handling stub (mirrors keys_android.go)
func acceleratorToString(accelerator *accelerator) string {
	return ""
}
