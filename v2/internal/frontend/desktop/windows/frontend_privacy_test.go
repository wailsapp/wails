//go:build windows

package windows

import (
	"errors"
	"github.com/wailsapp/go-webview2/pkg/edge"
	opts "github.com/wailsapp/wails/v2/pkg/options/windows"
	"testing"
)

type privacyControl uint8

const (
	autofill privacyControl = iota
	password
	reputation
)

type privacyProbe struct {
	set func(privacyControl, bool) error
}

func (p privacyProbe) PutIsGeneralAutofillEnabled(v bool) error     { return p.set(autofill, v) }
func (p privacyProbe) PutIsPasswordAutosaveEnabled(v bool) error    { return p.set(password, v) }
func (p privacyProbe) PutIsReputationCheckingRequired(v bool) error { return p.set(reputation, v) }

func TestLocalShellPrivacyPrecedesEmbedAndNavigation(t *testing.T) {
	refused := errors.New("capability/HRESULT refused")
	for _, tc := range []struct {
		enabled, embed bool
		fail           *privacyControl
		path           string
		wantError      bool
	}{
		{enabled: true, embed: true}, {embed: true}, {enabled: true, embed: false, wantError: true},
		{enabled: true, embed: true, fail: ptrControl(autofill), wantError: true},
		{enabled: true, embed: true, fail: ptrControl(password), wantError: true},
		{enabled: true, embed: true, fail: ptrControl(reputation), wantError: true},
		{enabled: true, embed: true, path: `C:\private trace\net.json`},
		{enabled: true, embed: true, path: `relative.json`, wantError: true},
		{enabled: true, embed: true, path: `C:\trace" --remote-debugging-port=9`, wantError: true},
	} {
		chromium := &edge.Chromium{}
		embedded, navigated, calls := false, false, 0
		probe := privacyProbe{set: func(control privacyControl, value bool) error {
			privacyCheck(t, embedded && !value, embedded, value)
			calls++
			if tc.fail != nil && *tc.fail == control {
				return refused
			}
			return nil
		}}
		err := embedChromium(chromium, &opts.Options{LocalShellPrivacy: tc.enabled, DiagnosticNetLogPath: tc.path}, func() bool {
			privacyCheck(t, tc.enabled == chromium.CustomCrashReportingEnabled, tc.enabled, chromium.CustomCrashReportingEnabled)
			wantArgs := 0
			if tc.path != "" {
				wantArgs = 1
				privacyCheck(t, len(chromium.AdditionalBrowserArgs) == 1 && chromium.AdditionalBrowserArgs[0] == "--log-net-log=\""+tc.path+"\"", chromium.AdditionalBrowserArgs)
			}
			privacyCheck(t, len(chromium.AdditionalBrowserArgs) == wantArgs, chromium.AdditionalBrowserArgs, wantArgs)
			embedded = true
			return tc.embed
		}, probe, func() {
			navigated = true
			if tc.enabled {
				privacyCheck(t, 3 == calls, 3, calls)
			}
		})
		privacyCheck(t, tc.wantError == (err != nil) && navigated == !tc.wantError, err, navigated)
		if tc.fail != nil {
			privacyCheck(t, errors.Is(err, refused), err, refused)
		}
		if !tc.enabled {
			privacyCheck(t, calls == 0, "Zero")
		}
	}
}
func ptrControl(c privacyControl) *privacyControl { return &c }

func privacyCheck(t *testing.T, ok bool, values ...any) {
	t.Helper()
	if !ok {
		t.Fatal(values...)
	}
}
