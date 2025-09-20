//go:build ios

package notifications

import "github.com/wailsapp/wails/v3/pkg/services"

type darwinNotificationService struct{}

func NewService() services.Service {
	return &darwinNotificationService{}
}

func (s *darwinNotificationService) Name() string {
	return "notifications"
}

func (s *darwinNotificationService) Route() string {
	return "/notifications"
}

func (s *darwinNotificationService) Shutdown() {}

func sendNotification(opts SendNotificationOptions) error {
	// iOS notification implementation would go here
	return nil
}