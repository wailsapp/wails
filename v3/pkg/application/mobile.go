package application

// MobileManager is the cross-platform surface common to the iOS and Android
// native-feature managers. The package-level Mobile singleton implements it,
// dispatching to the IOS or Android manager on the respective platform and to a
// no-op stub on desktop builds. This lets cross-platform mobile code call
// application.Mobile.* without splitting into //go:build ios / //go:build
// android files of its own.
//
// Only capabilities whose signature is identical on both platforms live here.
// Features that differ in shape between iOS and Android — brightness
// (float64 vs int), the orientation/brightness getters, notifications
// (PostNotification vs Notify), SecureSet (key+value vs JSON) and background
// execution (background task vs foreground service) — are intentionally left
// off; reach for application.IOS / application.Android directly for those.
//
// Because both managers must satisfy this interface, adding a method here also
// guarantees the two platform managers keep that method signature-identical: if
// they ever drift, the platform build fails to compile.
type MobileManager interface {
	// One-way actions
	Share(jsonPayload string)
	OpenURL(url string)
	SetKeepAwake(enabled bool)
	SetTorch(enabled bool)

	// State / query (JSON or path results)
	SafeAreaJSON() string
	AppInfoJSON() string
	SetOrientation(mode string)
	SetStatusBar(jsonPayload string)
	StorageJSON() string
	StoragePath() string
	PowerJSON() string
	NetworkJSON() string

	// Permissions / async results (delivered as common:* events)
	BiometricAuthenticate(reason string)
	SecureGet(key string) string
	SecureDelete(key string)
	GetLocation()

	// Sensors & hardware
	Haptic(hapticType string)
	SetMotion(enabled bool)
	SetProximity(enabled bool)
	Speak(text string)
	StopSpeak()
	SetKeyboardWatch(enabled bool)
	SetScreenProtect(enabled bool)

	// Camera
	CapturePhoto()
	CaptureVideo()
}
