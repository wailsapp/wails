package main

import (
	"fmt"
	"io"
	"net/http"
)

func downloadBytes(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("download %s: HTTP %s", url, response.Status)
	}
	return io.ReadAll(response.Body)
}
