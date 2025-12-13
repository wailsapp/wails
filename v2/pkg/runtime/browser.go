package runtime

import (
	"context"
	"fmt"
	"net/url"
)

// BrowserOpenURL uses the system default browser to open the url
func BrowserOpenURL(ctx context.Context, rawURL string) {
	appFrontend := getFrontend(ctx)

	parsed, err := url.Parse(rawURL)
	if err != nil {
		LogError(ctx, fmt.Sprintf("BrowserOpenURL cannot parse url: %s", err.Error()))
		return
	}

	appFrontend.BrowserOpenURL(parsed.String())
}
