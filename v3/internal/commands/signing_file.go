package commands

import (
	"fmt"
	"io"
	"os"
)

func readSigningFile(path string, limit int64) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > limit {
		return nil, fmt.Errorf("signing input %q must be a regular file no larger than %d bytes", path, limit)
	}
	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if int64(len(data)) > limit {
		clear(data)
		return nil, fmt.Errorf("signing input exceeds %d bytes", limit)
	}
	return data, err
}
