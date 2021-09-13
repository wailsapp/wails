package runtime

import (
	"context"
)

// BrowserOpenURL uses the system default browser to open the url
func BrowserOpenURL(ctx context.Context, url string) {
	appFrontend := getFrontend(ctx)
	appFrontend.BrowserOpenURL(url)
}
