//go:build !ios && !android

package application

// mobileStub is the desktop implementation of MobileManager. Every action is a
// no-op and every query returns its zero value, so cross-platform code that
// calls application.Mobile.* compiles and runs off-device without effect.
type mobileStub struct{}

// Mobile is the cross-platform mobile manager. Off-device it is a no-op stub.
var Mobile MobileManager = mobileStub{}

func (mobileStub) Share(string)                 {}
func (mobileStub) OpenURL(string)               {}
func (mobileStub) SetKeepAwake(bool)            {}
func (mobileStub) SetTorch(bool)                {}
func (mobileStub) SafeAreaJSON() string         { return "" }
func (mobileStub) AppInfoJSON() string          { return "" }
func (mobileStub) SetOrientation(string)        {}
func (mobileStub) SetStatusBar(string)          {}
func (mobileStub) StorageJSON() string          { return "" }
func (mobileStub) StoragePath() string          { return "" }
func (mobileStub) PowerJSON() string            { return "" }
func (mobileStub) NetworkJSON() string          { return "" }
func (mobileStub) BiometricAuthenticate(string) {}
func (mobileStub) SecureGet(string) string      { return "" }
func (mobileStub) SecureDelete(string)          {}
func (mobileStub) GetLocation()                 {}
func (mobileStub) Haptic(string)                {}
func (mobileStub) SetMotion(bool)               {}
func (mobileStub) SetProximity(bool)            {}
func (mobileStub) Speak(string)                 {}
func (mobileStub) StopSpeak()                   {}
func (mobileStub) SetKeyboardWatch(bool)        {}
func (mobileStub) SetScreenProtect(bool)        {}
func (mobileStub) CapturePhoto()                {}
func (mobileStub) CaptureVideo()                {}
