// mac-feedback demonstrates the macOS feedback APIs in Wails v3: a status
// bar item drawn from an SF Symbol with a tooltip that the user may remove,
// trackpad haptics, system sounds, text to speech and microphone speech
// recognition.
//
// Run it with `go run .` from this directory. Speech recognition needs the
// app to be a bundle whose Info.plist declares
// NSSpeechRecognitionUsageDescription and NSMicrophoneUsageDescription; see
// README.md.
package main

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

const recogniseDuration = 5 * time.Second

func main() {
	app := application.New(application.Options{
		Name:        "Feedback",
		Description: "macOS haptics, sounds, speech and status item demo",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Feedback",
		Width:  640,
		Height: 720,
		URL:    "/",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
			InvisibleTitleBarHeight: 40,
		},
	})

	// Status item: SF Symbol, tooltip, user removal with a stable autosave
	// name so macOS remembers when the user drags it out of the menu bar.
	tray := app.SystemTray.New()
	tray.SetSymbol("waveform.circle").SetSymbolConfiguration(0, application.MacSymbolWeightMedium)
	tray.SetTooltip("Wails feedback demo. Command-drag to remove.")
	tray.SetRemovable(true, "io.wails.examples.mac-feedback.tray")
	tray.OnVisibilityChange(func(visible bool) {
		app.Logger.Info("tray visibility changed", "visible", visible)
		app.Event.Emit("tray:visible", visible)
	})
	trayMenu := app.Menu.New()
	trayMenu.Add("Show Window").OnClick(func(*application.Context) {
		window.Show().Focus()
	})
	trayMenu.Add("Beep").OnClick(func(*application.Context) {
		app.Sound.Beep()
	})
	trayMenu.AddSeparator()
	trayMenu.Add("Quit").OnClick(func(*application.Context) {
		app.Quit()
	})
	tray.SetMenu(trayMenu)

	// Tray visibility from the page.
	app.Event.On("tray:show", func(*application.CustomEvent) {
		tray.Show()
		app.Event.Emit("tray:visible", tray.IsVisible())
	})
	app.Event.On("tray:hide", func(*application.CustomEvent) {
		tray.Hide()
		app.Event.Emit("tray:visible", tray.IsVisible())
	})

	// Haptics: only felt on a Force Touch trackpad while the app is active.
	app.Event.On("haptic", func(event *application.CustomEvent) {
		kind := application.HapticGeneric
		switch event.Data {
		case "alignment":
			kind = application.HapticAlignment
		case "levelChange":
			kind = application.HapticLevelChange
		}
		app.Haptics.Perform(kind)
		app.Event.Emit("status", fmt.Sprintf("haptic %s performed (supported=%v)", kind, app.Haptics.IsSupported()))
	})

	// Sounds.
	app.Event.On("sound:beep", func(*application.CustomEvent) {
		app.Sound.Beep()
		app.Event.Emit("status", "beep")
	})
	app.Event.On("sound:play", func(event *application.CustomEvent) {
		name, _ := event.Data.(string)
		if err := app.Sound.Play(name); err != nil {
			app.Event.Emit("status", "sound error: "+err.Error())
			return
		}
		app.Event.Emit("status", "playing "+name)
	})
	app.Event.On("sound:list", func(*application.CustomEvent) {
		app.Event.Emit("sound:names", app.Sound.SystemSounds())
	})

	// Speech synthesis.
	app.Event.On("speech:voices", func(*application.CustomEvent) {
		app.Event.Emit("speech:voices:result", app.Speech.Voices())
	})
	app.Event.On("speech:speak", func(event *application.CustomEvent) {
		params, _ := event.Data.(map[string]any)
		text, _ := params["text"].(string)
		voice, _ := params["voice"].(string)
		rate, _ := params["rate"].(float64)
		utterance, err := app.Speech.Speak(text, application.SpeechOptions{Voice: voice, Rate: rate})
		if err != nil {
			app.Event.Emit("status", "speech error: "+err.Error())
			return
		}
		app.Event.Emit("status", fmt.Sprintf("speaking utterance %d", utterance.ID()))
		utterance.OnFinished(func() {
			if utterance.WasStopped() {
				app.Event.Emit("status", fmt.Sprintf("utterance %d stopped", utterance.ID()))
				return
			}
			app.Event.Emit("status", fmt.Sprintf("utterance %d finished", utterance.ID()))
		})
	})
	app.Event.On("speech:stop", func(*application.CustomEvent) {
		app.Speech.StopAll()
	})

	// Speech recognition: capture the microphone for five seconds and report
	// partial and final transcripts.
	app.Event.On("speech:recognize", func(*application.CustomEvent) {
		go recognise(app)
	})

	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		tray.Run()
		app.Event.Emit("tray:visible", tray.IsVisible())
		if runtime.GOOS != "darwin" {
			app.Event.Emit("status", "this example targets macOS; most buttons are no-ops here")
		}
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func recognise(app *application.App) {
	session, err := app.Speech.Recognize(application.RecognitionOptions{
		OnPartial: func(text string) {
			app.Event.Emit("speech:partial", text)
		},
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrSpeechRecognitionUsageDescription):
			app.Event.Emit("speech:error", "Speech recognition needs an app bundle with NSSpeechRecognitionUsageDescription and NSMicrophoneUsageDescription in Info.plist (see README.md).")
		case errors.Is(err, application.ErrSpeechRecognitionDenied):
			app.Event.Emit("speech:error", "Permission denied. Allow Speech Recognition and Microphone for this app in System Settings > Privacy & Security.")
		default:
			app.Event.Emit("speech:error", err.Error())
		}
		return
	}
	app.Event.Emit("status", "listening for five seconds")
	time.Sleep(recogniseDuration)
	text, err := session.Stop()
	if err != nil {
		app.Event.Emit("speech:error", err.Error())
		return
	}
	app.Event.Emit("speech:final", text)
}
