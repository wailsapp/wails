# macOS feedback: status item, haptics, sounds and speech

This example exercises the platform feedback APIs on macOS:

- a status bar item drawn from an SF Symbol (`SystemTray.SetSymbol`,
  `SetSymbolConfiguration`) with a tooltip (`SetTooltip`) that the user can
  Command-drag out of the menu bar (`SetRemovable` with an autosave name);
  `OnVisibilityChange`, `IsVisible` and `SetVisible` track and restore it;
- trackpad haptics through `app.Haptics.Perform` (`HapticGeneric`,
  `HapticAlignment`, `HapticLevelChange`);
- system sounds through `app.Sound.Beep`, `app.Sound.Play` (named sounds
  such as "Glass" or absolute file paths), `app.Sound.PlayData` and
  `app.Sound.SystemSounds`;
- text to speech through `app.Speech.Speak`, `app.Speech.Voices` and
  `app.Speech.StopAll`, with `Utterance.Stop` and `Utterance.OnFinished`;
- microphone speech recognition through `app.Speech.Recognize`, with partial
  transcripts and `RecognitionSession.Stop` returning the final text.

The page talks to Go with custom events only, so there are no bindings to
generate.

## Running

```shell
go run .
```

Everything except speech recognition works from `go run`. Recognition uses
`SFSpeechRecognizer` and the microphone, and macOS only grants those to an app
bundle whose `Info.plist` declares why it needs them. Without the keys
`Recognize` returns `application.ErrSpeechRecognitionUsageDescription` and the
page shows that message.

To try recognition, package the example with `wails3 package` (or build your
own bundle) and add these keys to the bundle's `Info.plist`:

```xml
<key>NSSpeechRecognitionUsageDescription</key>
<string>Transcribes what you say into the text field.</string>
<key>NSMicrophoneUsageDescription</key>
<string>Listens to the microphone while you dictate.</string>
```

The first call to `Recognize` prompts for both permissions. A refusal returns
`application.ErrSpeechRecognitionDenied`; the decision can be changed later in
System Settings > Privacy & Security > Speech Recognition and Microphone.
Recognition happens on device when the locale supports it, otherwise the audio
is sent to Apple's servers.

## Platform notes

- Haptics are only felt on a Force Touch trackpad or Magic Trackpad and only
  while the app is active. On iOS and Android `Perform` is routed to the
  platform haptics engine; on Windows and Linux it is a no-op.
- `Sound.Beep` uses `NSBeep` on macOS and `MessageBeep` on Windows. `Play`
  accepts registry sound aliases such as "SystemAsterisk" or WAV file paths
  on Windows. Linux returns `ErrSoundNotSupported`.
- Speech synthesis needs macOS 10.14 or later and recognition needs 10.15 or
  later; other platforms return `ErrSpeechNotSupported` and
  `ErrSpeechRecognitionNotSupported`.
- `SetSymbol` needs macOS 11 or later and is ignored on older versions and
  other platforms, where `SetIcon` or `SetTemplateIcon` should be used.
