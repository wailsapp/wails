package edge

import (
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/webview2/internal/w32"
	"golang.org/x/sys/windows"
)

func TestCookieManager(t *testing.T) {
	// Initialize COM
	err := windows.CoInitializeEx(0, windows.COINIT_APARTMENTTHREADED)
	if err != nil {
		t.Fatalf("Failed to initialize COM: %v", err)
	}
	defer windows.CoUninitialize()

	// Create a temporary window for WebView2
	var hinstance windows.Handle
	_ = windows.GetModuleHandleEx(0, nil, &hinstance)

	// Load default icon
	icow, _, _ := w32.User32GetSystemMetrics.Call(w32.SystemMetricsCxIcon)
	icoh, _, _ := w32.User32GetSystemMetrics.Call(w32.SystemMetricsCyIcon)
	icon, _, _ := w32.User32LoadImageW.Call(uintptr(hinstance), 32512, icow, icoh, 0)

	className, _ := windows.UTF16PtrFromString("WebView2Test")
	wc := w32.WndClassExW{
		CbSize:        uint32(unsafe.Sizeof(w32.WndClassExW{})),
		HInstance:     hinstance,
		LpszClassName: className,
		HIcon:         windows.Handle(icon),
		HIconSm:       windows.Handle(icon),
		LpfnWndProc:   windows.NewCallback(w32.DefWindowProc),
	}
	_, _, _ = w32.User32RegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	windowName, _ := windows.UTF16PtrFromString("WebView2 Test Window")
	hwnd, _, _ := w32.User32CreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		0xCF0000, // WS_OVERLAPPEDWINDOW
		uintptr(w32.CW_USEDEFAULT),
		uintptr(w32.CW_USEDEFAULT),
		640,
		480,
		0,
		0,
		uintptr(hinstance),
		0,
	)
	if hwnd == 0 {
		t.Fatal("Failed to create window")
	}
	defer w32.DestroyWindow(hwnd)

	_, _, _ = w32.User32ShowWindow.Call(hwnd, w32.SWShow)
	_, _, _ = w32.User32UpdateWindow.Call(hwnd)
	_, _, _ = w32.User32SetFocus.Call(hwnd)

	// Create a new Chromium instance
	chromium := NewChromium()
	require.NotNil(t, chromium, "Chromium instance should not be nil")

	// Initialize WebView2
	success := chromium.Embed(uintptr(hwnd))
	require.True(t, success, "WebView2 initialization should succeed")

	// Get the cookie manager
	cookieManager, err := chromium.GetCookieManager()
	require.NoError(t, err, "Should get cookie manager without error")
	require.NotNil(t, cookieManager, "Cookie manager should not be nil")
	defer cookieManager.Release()

	// Delete all cookies to start with a clean slate
	err = cookieManager.DeleteAllCookies()
	require.NoError(t, err, "Should delete all cookies without error")

	t.Run("Test Cookie Creation and Properties", func(t *testing.T) {
		// Create a new cookie
		cookie, err := cookieManager.CreateCookie("testCookie", "testValue", "example.com", "/test")
		require.NoError(t, err, "Should create cookie without error")
		require.NotNil(t, cookie, "Cookie should not be nil")
		defer cookie.Release()

		// Test GetName
		name, err := cookie.GetName()
		assert.NoError(t, err, "Should get name without error")
		assert.Equal(t, "testCookie", name, "Cookie name should match")

		// Test GetValue/PutValue
		value, err := cookie.GetValue()
		assert.NoError(t, err, "Should get value without error")
		assert.Equal(t, "testValue", value, "Cookie value should match")

		err = cookie.PutValue("newValue")
		assert.NoError(t, err, "Should put value without error")
		value, err = cookie.GetValue()
		assert.NoError(t, err, "Should get updated value without error")
		assert.Equal(t, "newValue", value, "Cookie value should be updated")

		// Test GetDomain
		domain, err := cookie.GetDomain()
		assert.NoError(t, err, "Should get domain without error")
		assert.Equal(t, "example.com", domain, "Cookie domain should match")

		// Test GetPath
		path, err := cookie.GetPath()
		assert.NoError(t, err, "Should get path without error")
		assert.Equal(t, "/test", path, "Cookie path should match")

		// Test Expires
		now := time.Now().Add(24 * time.Hour) // 24 hours from now
		comTime := float64(now.Unix())
		err = cookie.PutExpires(comTime)
		assert.NoError(t, err, "Should set expiration without error")
		expires, err := cookie.GetExpires()
		assert.NoError(t, err, "Should get expiration without error")
		assert.Equal(t, comTime, expires, "Cookie expiration should match")

		// Test IsHttpOnly
		err = cookie.PutIsHttpOnly(true)
		assert.NoError(t, err, "Should set HttpOnly without error")
		isHttpOnly, err := cookie.GetIsHttpOnly()
		assert.NoError(t, err, "Should get HttpOnly without error")
		assert.True(t, isHttpOnly, "Cookie should be HttpOnly")

		// Test IsSecure
		err = cookie.PutIsSecure(true)
		assert.NoError(t, err, "Should set Secure without error")
		isSecure, err := cookie.GetIsSecure()
		assert.NoError(t, err, "Should get Secure without error")
		assert.True(t, isSecure, "Cookie should be Secure")

		// Test SameSite
		err = cookie.PutSameSite(2) // 2 = Lax
		assert.NoError(t, err, "Should set SameSite without error")
		sameSite, err := cookie.GetSameSite()
		assert.NoError(t, err, "Should get SameSite without error")
		assert.Equal(t, int32(2), sameSite, "Cookie SameSite should be Lax")
	})

	t.Run("Test Cookie Management", func(t *testing.T) {
		// Create and add a cookie
		cookie, err := cookieManager.CreateCookie("managedCookie", "testValue", "example.com", "/test")
		require.NoError(t, err, "Should create cookie without error")
		defer cookie.Release()

		err = cookieManager.AddOrUpdateCookie(cookie)
		assert.NoError(t, err, "Should add cookie without error")

		// Delete the cookie
		err = cookieManager.DeleteCookie(cookie)
		assert.NoError(t, err, "Should delete cookie without error")

		// Delete all cookies
		err = cookieManager.DeleteAllCookies()
		assert.NoError(t, err, "Should delete all cookies without error")
	})
}
