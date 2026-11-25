//go:build openharmony && cgo && !server

package application

import "encoding/json"

// openharmonyApp mirrors androidApp: the OHOS lifecycle is driven by the ArkTS
// ability, which owns the JS thread and the ArkWeb component. Go provides the
// runtime, bindings, events and bridge; run() forwards to platformRun.
type openharmonyApp struct {
	parent *App
}

func newPlatformApp(app *App) *openharmonyApp {
	return &openharmonyApp{parent: app}
}

func (a *App) logPlatformInfo() { a.info("Platform Info:", "platform", "openharmony") }

func (a *App) platformEnvironment() map[string]any {
	return map[string]any{"platform": "openharmony"}
}

// The host owns process startup; there is no platform fatal-error hook to install.
func fatalHandler(_ func(error)) {}

func (a *openharmonyApp) run() error {
	a.setupCommonEvents()
	a.parent.platformRun()
	return nil
}

func (a *openharmonyApp) destroy() {}

func (a *openharmonyApp) setApplicationMenu(_ *Menu) {}

func (a *openharmonyApp) name() string { return a.parent.options.Name }

func (a *openharmonyApp) GetFlags(_ Options) map[string]any { return nil }

func (a *openharmonyApp) getAccentColor() string { return "" }

func (a *openharmonyApp) getCurrentWindowID() uint { return 0 }

func (a *openharmonyApp) hide() { ohBridgeVoid("hide", "{}") }

func (a *openharmonyApp) show() { ohBridgeVoid("show", "{}") }

func (a *openharmonyApp) isDarkMode() bool {
	res, ok := ohBridgeString("isDarkMode", "{}")
	if !ok || res == "" {
		return false
	}
	var result struct {
		Dark bool `json:"dark"`
	}
	if err := json.Unmarshal([]byte(res), &result); err != nil {
		return false
	}
	return result.Dark
}

func (a *openharmonyApp) isWindows() bool { return false }

func (a *openharmonyApp) on(eventID uint) { registerOpenharmonyListener(eventID) }

func (a *openharmonyApp) setIcon(_ []byte) {}

func (a *openharmonyApp) showAboutDialog(name, description string, _ []byte) {
	payload, _ := json.Marshal(struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{Name: name, Description: description})
	ohBridgeVoid("showAboutDialog", string(payload))
}

func (a *openharmonyApp) getPrimaryScreen() (*Screen, error) {
	screens, err := a.getScreens()
	if err != nil {
		return nil, err
	}
	if len(screens) == 0 {
		return nil, nil
	}
	return screens[0], nil
}

func (a *openharmonyApp) getScreens() ([]*Screen, error) { return getScreens() }

func (a *openharmonyApp) setupCommonEvents() {
	for sourceEvent, targetEvent := range commonApplicationEventMap {
		sourceEvent := sourceEvent
		targetEvent := targetEvent
		a.parent.Event.OnApplicationEvent(sourceEvent, func(event *ApplicationEvent) {
			applicationEvents <- &ApplicationEvent{
				Id:  uint(targetEvent),
				ctx: event.ctx,
			}
		})
	}
}
