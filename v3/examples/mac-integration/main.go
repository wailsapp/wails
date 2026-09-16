package main

import (
	"embed"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

// activityType must be listed under NSUserActivityTypes in Info.plist for
// Handoff to advertise it.
const activityType = "com.wails.example.mac-integration.editing"

// IntegrationService is bound to the frontend and wraps the ServicesProvider,
// Activity, Browser and Env managers.
type IntegrationService struct {
	app *application.App

	lock      sync.Mutex
	published *application.PublishedActivity
}

// ServiceInfoPlist returns the NSServices XML for the registered services so
// the page can show what packaging has to include.
func (s *IntegrationService) ServiceInfoPlist() string {
	return s.app.ServicesProvider.InfoPlistXML()
}

// PublishActivity mirrors the text field into a Handoff activity: the first
// call publishes, later calls update the live activity's UserInfo.
func (s *IntegrationService) PublishActivity(text string) (application.UserActivity, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	userInfo := map[string]any{"text": text, "updatedAt": time.Now().Format(time.RFC3339)}
	if s.published != nil {
		if err := s.published.Update(userInfo); err != nil {
			return application.UserActivity{}, err
		}
		return s.published.Activity(), nil
	}
	published, err := s.app.Activity.Publish(application.UserActivity{
		Type:               activityType,
		Title:              "Editing in Wails mac-integration",
		UserInfo:           userInfo,
		WebpageURL:         "https://wails.io/",
		EligibleForHandoff: true,
		EligibleForSearch:  true,
		Keywords:           []string{"wails", "handoff"},
	})
	if err != nil {
		return application.UserActivity{}, err
	}
	s.published = published
	return published.Activity(), nil
}

// InvalidateActivity ends the Handoff activity.
func (s *IntegrationService) InvalidateActivity() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.published != nil {
		s.published.Invalidate()
		s.published = nil
	}
}

// ApplicationsForFile lists the applications that can open path.
func (s *IntegrationService) ApplicationsForFile(path string) []application.AppInfo {
	return s.app.Browser.ApplicationsForFile(path)
}

// OpenWith opens path with the given bundle identifier or .app path.
func (s *IntegrationService) OpenWith(path string, app string) error {
	return s.app.Browser.OpenWith(path, app)
}

// ActivateApplication brings a running application to the front.
func (s *IntegrationService) ActivateApplication(bundleID string) error {
	return s.app.Browser.ActivateApplication(bundleID)
}

// DefaultHandler reports which application opens a content type or scheme.
func (s *IntegrationService) DefaultHandler(contentType string, scheme string) (application.AppInfo, error) {
	return s.app.Env.DefaultHandler(application.DefaultHandler{ContentType: contentType, URLScheme: scheme})
}

// summarise keeps the first sentence and appends a word count.
func summarise(text string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return ""
	}
	first := trimmed
	for _, stop := range []string{". ", "! ", "? ", "\n"} {
		if index := strings.Index(first, stop); index >= 0 {
			first = first[:index+1]
		}
	}
	return fmt.Sprintf("%s (%d words)", strings.TrimSpace(first), len(strings.Fields(trimmed)))
}

func main() {
	service := &IntegrationService{}

	app := application.New(application.Options{
		Name:        "Wails Example",
		Description: "Services menu, Handoff, universal links and workspace helpers",
		Services: []application.Service{
			application.NewService(service),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	service.app = app

	// Services menu: "Summarise with Wails Example" replaces the selected
	// text in any app with its summary. The matching NSServices entry lives
	// in build/config.yml under `services` (see README.md).
	err := app.ServicesProvider.Register(application.ServiceDefinition{
		Name:          "SummariseWithWailsExample",
		MenuTitle:     "Summarise with Wails Example",
		SendTypes:     []string{"public.utf8-plain-text"},
		ReturnTypes:   []string{"public.utf8-plain-text"},
		KeyEquivalent: "S",
		Handler: func(ctx *application.Context, req application.ServiceRequest) (application.ServiceResponse, error) {
			summary := summarise(req.Text)
			app.Event.Emit("integration:service", map[string]any{"input": req.Text, "output": summary, "types": req.Types})
			return application.ServiceResponse{Text: summary}, nil
		},
	})
	if err != nil {
		log.Printf("services: %v", err)
	}

	// Handoff and universal links. A universal link arrives both here (as a
	// browsing activity) and through ApplicationLaunchedWithUrl below.
	app.Activity.OnWillContinue(func(activityType string) bool {
		app.Event.Emit("integration:activity", map[string]any{"phase": "will-continue", "type": activityType})
		return true
	})
	app.Activity.OnContinue(func(ctx *application.Context, activity application.UserActivity) bool {
		app.Event.Emit("integration:activity", map[string]any{"phase": "continue", "type": activity.Type, "title": activity.Title, "userInfo": activity.UserInfo, "webpageURL": activity.WebpageURL})
		return activity.Type == activityType || activity.Type == application.UserActivityTypeBrowsingWeb
	})
	app.Activity.OnFailed(func(activityType string, err error) {
		app.Event.Emit("integration:activity", map[string]any{"phase": "failed", "type": activityType, "error": err.Error()})
	})
	app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(event *application.ApplicationEvent) {
		app.Event.Emit("integration:url", event.Context().URL())
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Platform Integration",
		Width:  960,
		Height: 800,
		URL:    "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(fmt.Errorf("run: %w", err))
	}
}
