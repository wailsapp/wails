package runtime

import "context"

func IsNotificationAvailable(ctx context.Context) bool {
	frontend := getFrontend(ctx)
	return frontend.IsNotificationAvailable()
}

func RequestNotificationAuthorization(ctx context.Context) (bool, error) {
	frontend := getFrontend(ctx)
	return frontend.RequestNotificationAuthorization()
}

func CheckNotificationAuthorization(ctx context.Context) (bool, error) {
	frontend := getFrontend(ctx)
	return frontend.CheckNotificationAuthorization()
}
