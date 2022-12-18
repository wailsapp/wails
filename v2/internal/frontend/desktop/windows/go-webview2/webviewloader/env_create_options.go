//go:build windows && !native_webview2loader

package webviewloader

import (
	"unicode/utf16"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2/pkg/combridge"
	"golang.org/x/sys/windows"
)

// WithBrowserExecutableFolder to specify whether WebView2 controls use a fixed or installed version
// of the WebView2 Runtime that exists on a user machine.
//
// To use a fixed version of the WebView2 Runtime,
// pass the folder path that contains the fixed version of the WebView2 Runtime.
// BrowserExecutableFolder supports both relative (to the application's executable) and absolute files paths.
// To create WebView2 controls that use the installed version of the WebView2 Runtime that exists on user
// machines, pass a empty string to WithBrowserExecutableFolder. In this scenario, the API tries to find a
// compatible version of the WebView2 Runtime that is installed on the user machine (first at the machine level,
// and then per user) using the selected channel preference. The path of fixed version of the WebView2 Runtime
// should not contain \Edge\Application\. When such a path is used, the API fails with HRESULT_FROM_WIN32(ERROR_NOT_SUPPORTED).
func WithBrowserExecutableFolder(folder string) option {
	return func(wvep *environmentOptions) {
		wvep.browserExecutableFolder = folder
	}
}

// WithUserDataFolder specifies to user data folder location for WebView2
//
// You may specify the userDataFolder to change the default user data folder location for WebView2.
// The path is either an absolute file path or a relative file path that is interpreted as relative
// to the compiled code for the current process.
// Dhe default user data ({Executable File Name}.WebView2) folder is created in the same directory
// next to the compiled code for the app. WebView2 creation fails if the compiled code is running
// in a directory in which the process does not have permission to create a new directory.
// The app is responsible to clean up the associated user data folder when it is done.
func WithUserDataFolder(folder string) option {
	return func(wvep *environmentOptions) {
		wvep.userDataFolder = folder
	}
}

// WithAdditionalBrowserArguments changes the behavior of the WebView.
//
// The arguments are passed to the
// browser process as part of the command.  For more information about
// using command-line switches with Chromium browser processes, navigate to
// [Run Chromium with Flags][ChromiumDevelopersHowTosRunWithFlags].
// The value appended to a switch is appended to the browser process, for
// example, in `--edge-webview-switches=xxx` the value is `xxx`.  If you
// specify a switch that is important to WebView functionality, it is
// ignored, for example, `--user-data-dir`.  Specific features are disabled
// internally and blocked from being enabled. If a switch is specified
// multiple times, only the last instance is used.
//
// \> [!NOTE]\n\> A merge of the different values of the same switch is not attempted,
// except for disabled and enabled features. The features specified by
// `--enable-features` and `--disable-features` are merged with simple
// logic.\n\> *   The features is the union of the specified features
// and built-in features.  If a feature is disabled, it is removed from the
// enabled features list.
//
// If you specify command-line switches and use the
// `additionalBrowserArguments` parameter, the `--edge-webview-switches`
// value takes precedence and is processed last.  If a switch fails to
// parse, the switch is ignored.  The default state for the operation is
// to run the browser process with no extra flags.
//
// [ChromiumDevelopersHowTosRunWithFlags]: https://www.chromium.org/developers/how-tos/run-chromium-with-flags "Run Chromium with flags | The Chromium Projects"
func WithAdditionalBrowserArguments(args string) option {
	return func(wvep *environmentOptions) {
		wvep.additionalBrowserArguments = args
	}
}

// WithLanguage sets the default display language for WebView.
//
// It applies to browser UI such as
// context menu and dialogs.  It also applies to the `accept-languages` HTTP
// header that WebView sends to websites.  It is in the format of
//
// `language[-country]` where `language` is the 2-letter code from
// [ISO 639][ISO639LanguageCodesHtml]
// and `country` is the
// 2-letter code from
// [ISO 3166][ISOStandard72482Html].
//
// [ISO639LanguageCodesHtml]: https://www.iso.org/iso-639-language-codes.html "ISO 639 | ISO"
// [ISOStandard72482Html]: https://www.iso.org/standard/72482.html "ISO 3166-1:2020 | ISO"
func WithLanguage(lang string) option {
	return func(wvep *environmentOptions) {
		wvep.language = lang
	}
}

// WithTargetCompatibleBrowserVersion secifies the version of the WebView2 Runtime binaries required to be
// compatible with your app.
//
// This defaults to the WebView2 Runtime version
// that corresponds with the version of the SDK the app is using.  The
// format of this value is the same as the format of the
// `BrowserVersionString` property and other `BrowserVersion` values.  Only
// the version part of the `BrowserVersion` value is respected.  The channel
// suffix, if it exists, is ignored.  The version of the WebView2 Runtime
// binaries actually used may be different from the specified
// `TargetCompatibleBrowserVersion`. The binaries are only guaranteed to be
// compatible. Verify the actual version on the `BrowserVersionString`
// property on the `ICoreWebView2Environment`.
func WithTargetCompatibleBrowserVersion(version string) option {
	return func(wvep *environmentOptions) {
		wvep.targetCompatibleBrowserVersion = version
	}
}

// WithAllowSingleSignOnUsingOSPrimaryAccount is used to enable
// single sign on with Azure Active Directory (AAD) and personal Microsoft
// Account (MSA) resources inside WebView. All AAD accounts, connected to
// Windows and shared for all apps, are supported. For MSA, SSO is only enabled
// for the account associated for Windows account login, if any.
// Default is disabled. Universal Windows Platform apps must also declare
// `enterpriseCloudSSO`
// [Restricted capabilities][WindowsUwpPackagingAppCapabilityDeclarationsRestrictedCapabilities]
// for the single sign on (SSO) to work.
//
// [WindowsUwpPackagingAppCapabilityDeclarationsRestrictedCapabilities]: /windows/uwp/packaging/app-capability-declarations\#restricted-capabilities "Restricted capabilities - App capability declarations | Microsoft Docs"
func WithAllowSingleSignOnUsingOSPrimaryAccount(allow bool) option {
	return func(wvep *environmentOptions) {
		wvep.allowSingleSignOnUsingOSPrimaryAccount = allow
	}
}

// WithExclusiveUserDataFolderAccess specifies that the WebView environment
// obtains exclusive access to the user data folder.
//
// If the user data folder is already being used by another WebView environment with a
// different value for `ExclusiveUserDataFolderAccess` property, the creation of a WebView2Controller
// using the environment object will fail with `HRESULT_FROM_WIN32(ERROR_INVALID_STATE)`.
// When set as TRUE, no other WebView can be created from other processes using WebView2Environment
// objects with the same UserDataFolder. This prevents other processes from creating WebViews
// which share the same browser process instance, since sharing is performed among
// WebViews that have the same UserDataFolder. When another process tries to create a
// WebView2Controller from an WebView2Environment object created with the same user data folder,
// it will fail with `HRESULT_FROM_WIN32(ERROR_INVALID_STATE)`.
func WithExclusiveUserDataFolderAccess(exclusive bool) option {
	return func(wvep *environmentOptions) {
		wvep.exclusiveUserDataFolderAccess = exclusive
	}
}

type option func(*environmentOptions)

var _ iCoreWebView2EnvironmentOptions = &environmentOptions{}
var _ iCoreWebView2EnvironmentOptions2 = &environmentOptions{}

type environmentOptions struct {
	browserExecutableFolder string
	userDataFolder          string
	preferCanary            bool

	additionalBrowserArguments             string
	language                               string
	targetCompatibleBrowserVersion         string
	allowSingleSignOnUsingOSPrimaryAccount bool
	exclusiveUserDataFolderAccess          bool
}

func (o *environmentOptions) AdditionalBrowserArguments() string {
	return o.additionalBrowserArguments
}

func (o *environmentOptions) Language() string {
	return o.language
}

func (o *environmentOptions) TargetCompatibleBrowserVersion() string {
	v := o.targetCompatibleBrowserVersion
	if v == "" {
		v = kMinimumCompatibleVersion
	}
	return v
}

func (o *environmentOptions) AllowSingleSignOnUsingOSPrimaryAccount() bool {
	return o.allowSingleSignOnUsingOSPrimaryAccount
}

func (o *environmentOptions) ExclusiveUserDataFolderAccess() bool {
	return o.exclusiveUserDataFolderAccess
}

type iCoreWebView2EnvironmentOptions interface {
	combridge.IUnknown

	AdditionalBrowserArguments() string
	Language() string
	TargetCompatibleBrowserVersion() string
	AllowSingleSignOnUsingOSPrimaryAccount() bool
}

type iCoreWebView2EnvironmentOptions2 interface {
	combridge.IUnknown

	ExclusiveUserDataFolderAccess() bool
}

func init() {
	combridge.RegisterVTable[combridge.IUnknown, iCoreWebView2EnvironmentOptions](
		"{2fde08a8-1e9a-4766-8c05-95a9ceb9d1c5}",
		_iCoreWebView2EnvironmentOptionsAdditionalBrowserArguments,
		_iCoreWebView2EnvironmentOptionsNOP,
		_iCoreWebView2EnvironmentOptionsLanguage,
		_iCoreWebView2EnvironmentOptionsNOP,
		_iCoreWebView2EnvironmentTargetCompatibleBrowserVersion,
		_iCoreWebView2EnvironmentOptionsNOP,
		_iCoreWebView2EnvironmentOptionsAllowSingleSignOnUsingOSPrimaryAccount,
		_iCoreWebView2EnvironmentOptionsNOP,
	)

	combridge.RegisterVTable[combridge.IUnknown, iCoreWebView2EnvironmentOptions2](
		"{ff85c98a-1ba7-4a6b-90c8-2b752c89e9e2}",
		_iCoreWebView2EnvironmentOptions2ExclusiveUserDataFolderAccess,
		_iCoreWebView2EnvironmentOptionsNOP,
	)
}
func _iCoreWebView2EnvironmentOptionsNOP(this uintptr) uintptr {
	return uintptr(windows.S_FALSE)
}

func _iCoreWebView2EnvironmentOptionsAdditionalBrowserArguments(this uintptr, value **uint16) uintptr {
	v := combridge.Resolve[iCoreWebView2EnvironmentOptions](this).AdditionalBrowserArguments()
	*value = stringToOleString(v)
	return uintptr(windows.S_OK)
}

func _iCoreWebView2EnvironmentOptionsLanguage(this uintptr, value **uint16) uintptr {
	args := combridge.Resolve[iCoreWebView2EnvironmentOptions](this).Language()
	*value = stringToOleString(args)
	return uintptr(windows.S_OK)
}

func _iCoreWebView2EnvironmentTargetCompatibleBrowserVersion(this uintptr, value **uint16) uintptr {
	args := combridge.Resolve[iCoreWebView2EnvironmentOptions](this).TargetCompatibleBrowserVersion()
	*value = stringToOleString(args)
	return uintptr(windows.S_OK)
}

func _iCoreWebView2EnvironmentOptionsAllowSingleSignOnUsingOSPrimaryAccount(this uintptr, value *int32) uintptr {
	v := combridge.Resolve[iCoreWebView2EnvironmentOptions](this).AllowSingleSignOnUsingOSPrimaryAccount()
	*value = boolToInt(v)
	return uintptr(windows.S_OK)
}

func _iCoreWebView2EnvironmentOptions2ExclusiveUserDataFolderAccess(this uintptr, value *int32) uintptr {
	v := combridge.Resolve[iCoreWebView2EnvironmentOptions2](this).ExclusiveUserDataFolderAccess()
	*value = boolToInt(v)
	return uintptr(windows.S_OK)
}

func stringToOleString(v string) *uint16 {
	wstr := utf16.Encode([]rune(v + "\x00"))
	lwstr := len(wstr)
	ptr := (*uint16)(coTaskMemAlloc(2 * lwstr))

	copy(unsafe.Slice(ptr, lwstr), wstr)

	return ptr
}

func boolToInt(v bool) int32 {
	if v {
		return 1
	}
	return 0
}
