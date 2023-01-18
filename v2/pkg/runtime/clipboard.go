package runtime

import "context"

func ClipboardGetText(ctx context.Context) (string, error) {
	appFrontend := getFrontend(ctx)
	return appFrontend.ClipboardGetText()
}

func ClipboardSetText(ctx context.Context, text string) error {
	appFrontend := getFrontend(ctx)
	return appFrontend.ClipboardSetText(text)
}
