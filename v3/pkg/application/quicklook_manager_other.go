//go:build !darwin || ios || server

package application

func quickLookPreview([]string) error {
	return ErrQuickLookNotSupported
}

func quickLookClosePreview() {}

func quickLookIsPreviewOpen() bool {
	return false
}

func quickLookThumbnail(string, ThumbnailOptions) ([]byte, error) {
	return nil, ErrQuickLookNotSupported
}
