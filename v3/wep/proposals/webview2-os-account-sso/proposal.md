# Wails Enhancement Proposal (WEP)

## WebView2 single sign-on with the Windows account

**WEP Number**: (leave blank, assigned on acceptance)  
**Status**: Draft  
**Author**: @Grantmartin2002  
**Created**: 2026-11-19  
**Discussion**: none  
**Implementor**: @Grantmartin2002  
**Target**: Wails v3

## Summary

Add `WindowsOptions.AllowSingleSignOnUsingOSPrimaryAccount bool`. When it is true, Wails sets WebView2's
[`AllowSingleSignOnUsingOSPrimaryAccount`](https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environmentoptions#get_allowsinglesignonusingosprimaryaccount)
environment option, so pages in the webview can sign in to Microsoft Entra ID and Microsoft accounts with the
account the user signed in to Windows with. It defaults to false, which is today's behaviour.

## Motivation

On an Entra-joined or hybrid-joined PC, Edge signs users in to Entra ID without a prompt by using the Primary
Refresh Token (PRT) that Windows holds for the signed-in account. The PRT also carries the device identity, so
Conditional Access policies that require a managed or compliant device pass.

WebView2 does none of this unless the host turns the environment option on. Today a Wails app cannot: the
option is already bound in `internal/webview2/webviewloader` (`WithAllowSingleSignOnUsingOSPrimaryAccount`),
but nothing passes it, and `createCoreWebView2EnvironmentWithOptions` only forwards the browser folder, data
folder and browser arguments. There is no browser argument that does the same thing.

The result is that a Wails app loading an Entra ID login page asks the user for credentials they already gave
Windows, and under device-based Conditional Access it can be blocked outright. Two common cases:

- Apps that embed an internal web app behind Entra ID.
- Apps that authenticate with MSAL by loading the authorize URL in a (possibly hidden) webview with
  `prompt=none`. On macOS the Enterprise SSO extension already completes this inside WKWebView, so the same
  code works there and fails on Windows.

## Detailed Design

A new field on `application.WindowsOptions`, next to the other WebView2 environment settings:

```go
// AllowSingleSignOnUsingOSPrimaryAccount lets pages sign in to Microsoft Entra ID and Microsoft
// accounts with the account signed in to Windows. Applies to all windows.
AllowSingleSignOnUsingOSPrimaryAccount bool
```

Implementation:

1. `edge.Chromium` gets an `AllowSingleSignOnUsingOSPrimaryAccount bool` setting.
2. `createCoreWebView2EnvironmentWithOptions` takes it and passes
   `webviewloader.WithAllowSingleSignOnUsingOSPrimaryAccount(...)`.
3. `webview_window_windows.go` copies the value from the application-level `WindowsOptions`, the same way it
   copies `WebviewUserDataPath` and `WebviewBrowserPath`.

It is application-level, not per-window, on purpose. WebView2 requires every environment that shares a user
data folder to be created with the same options, so a per-window value would make the second window fail to
create.

## Non-Goals

- Calling the Windows Web Account Manager (WAM) broker from Go, or any token API in Wails. This only lets
  WebView2 do what Edge already does.
- Exposing the other unbound environment options (`ExclusiveUserDataFolderAccess`, `Language`, etc.). Each
  can be its own small change if someone needs it.
- macOS and Linux. See below.

## Platform Considerations

- **Windows**: SSO works when the user is signed in to Windows with an Entra ID account on an Entra-joined or
  hybrid-joined device, or with a Microsoft account. On other machines the option has no effect and the normal
  login page is shown. The property has been in the WebView2 SDK since 0.9.622 (2020), so every
  runtime Wails supports has it.
- **macOS**: no option needed. WKWebView requests are already handled by the Enterprise SSO extension when the
  device is configured for it and the app is covered by its configuration. The field lives on
  `WindowsOptions`, so it has no effect on other platforms.
- **Linux**: WebKitGTK has no equivalent integration. Unsupported.

## Pros/Cons

Pros:

- Small change that reuses an existing binding. Opt-in and off by default.
- Brings WebView2 to parity with Edge and with WKWebView on managed Macs.

Cons:

- One more Windows option to document.
- When on, any page loaded in the app can obtain Entra ID or Microsoft account sign-ins for the Windows user
  without a prompt. Apps that load untrusted content should leave it off.

## Alternatives Considered

- **Browser arguments or feature flags**: there is no Chromium switch or WebView2 feature flag that enables
  this. It is only exposed through `ICoreWebView2EnvironmentOptions`.
- **A generic hook to set environment options**: more flexible, but it would make the internal COM types part
  of the public API. A plain bool matches how the other WebView2 settings are exposed.
- **Opening the system browser instead**: works, but it is a visible round trip, and it is what apps that want
  silent sign-in are trying to avoid.

## Backwards Compatibility

None. The field defaults to false, which is what Wails passes today (the option is never set).

## Security and Privacy

When enabled, WebView2 attaches the Windows account's credentials to requests to Microsoft identity endpoints,
the same way Edge does. A page can then complete a Microsoft sign-in without user interaction. This is the
reason it is opt-in and documented with a caution not to enable it for apps that load untrusted origins.
Nothing is sent anywhere it would not be in Edge, and Wails does not read or store any credential.

## Test Plan

- On an Entra-joined Windows 11 machine (`dsregcmd /status` shows `AzureAdJoined : YES` and
  `AzureAdPrt : YES`), load `https://login.microsoftonline.com/<tenant>/oauth2/v2.0/authorize?...&prompt=none`
  in a window with the option on and confirm it redirects with a code. With the option off, confirm it returns
  `interaction_required`.
- Open two windows with the option on and confirm both are created (no options mismatch).
- On a machine that is not joined, confirm the option is harmless and the login page shows.

## Reference Implementation

#6612

## Maintenance Plan

The change is a pass-through to a stable WebView2 property, so it should need no upkeep beyond the existing
WebView2 code. The implementor will handle issues raised against it.

## Conclusion

Wails already carries the binding for this. Exposing it as one opt-in Windows option lets Wails apps get the
same silent Entra ID sign-in that Edge, and WKWebView on managed Macs, already get.
