//go:build ios

package dock

import "github.com/wailsapp/wails/v3/pkg/services"

type iosBadgeService struct{}

func NewService() services.Service {
	return &iosBadgeService{}
}

func (s *iosBadgeService) Name() string {
	return "badge"
}

func (s *iosBadgeService) Route() string {
	return "/badge"
}

func (s *iosBadgeService) Shutdown() {}

func setBadgeLabel(_ string) error {
	// iOS badge implementation would go here
	return nil
}

func clearBadgeLabel() error {
	// iOS badge implementation would go here
	return nil
}