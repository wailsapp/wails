package common

import (
	"fmt"
	"net/url"
)

var ErrUnexpectedScheme = fmt.Errorf("unexpected scheme")
var ErrUnexpectedHost = fmt.Errorf("unexpected host")

func translateUriToFile(uri string, expectedScheme string, expectedHosts ...string) (file string, err error) {
	url, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	if url.Scheme != expectedScheme {
		return "", ErrUnexpectedScheme
	}

	for _, expectedHost := range expectedHosts {
		if url.Host != expectedHost {
			continue
		}

		filePath := url.Path
		if filePath == "" {
			filePath = "/"
		}
		return filePath, nil
	}

	return "", ErrUnexpectedHost
}
