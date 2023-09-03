//go:build darwin

package doctor

func getInfo() (map[string]string, bool) {
	result := make(map[string]string)
	return result, true
}
