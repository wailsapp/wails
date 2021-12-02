package common

import "net/url"

func TranslateUriToFile(uri string, expectedScheme string, expectedHost string) (file string, match bool, err error) {
	url, err := url.Parse(uri)
	if err != nil {
		return "", false, err
	}

	if url.Scheme != expectedScheme || url.Host != expectedHost {
		return "", false, nil
	}

	filePath := url.Path
	if filePath == "" {
		filePath = "/"
	}
	return filePath, true, nil
}
