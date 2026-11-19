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

Add `WindowsOptions.AllowSingleSignOnUsingOSPrimaryAccount bool`, which sets the WebView2 environment option of
the same name. Pages in the webview can then sign in to Microsoft Entra ID and Microsoft accounts with the
account signed in to Windows, without a prompt.

## Motivation

On an Entra-joined or hybrid-joined PC, Edge signs in to Entra ID silently using the Windows account's Primary
Refresh Token, which also satisfies device-based Conditional Access. WebView2 only does this when the host sets
[`AllowSingleSignOnUsingOSPrimaryAccount`](https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environmentoptions#get_allowsinglesignonusingosprimaryaccount).

Wails never sets it. `webviewloader.WithAllowSingleSignOnUsingOSPrimaryAccount` exists, but
`createCoreWebView2EnvironmentWithOptions` doesn't pass it, and no browser argument does the same thing. So a
Wails app that loads an Entra ID login page asks for credentials the user already gave Windows, or is blocked by
Conditional Access. On a managed Mac the same page signs in silently, because the Enterprise SSO extension
handles WKWebView.

## Detailed Design

1. Add `AllowSingleSignOnUsingOSPrimaryAccount bool` to `application.WindowsOptions`, next to
   `AdditionalBrowserArgs`.
2. Add the same field to `edge.Chromium`, set from the application options in `webview_window_windows.go`.
3. Pass it to `webviewloader.WithAllowSingleSignOnUsingOSPrimaryAccount` in
   `createCoreWebView2EnvironmentWithOptions`.

```go
app := application.New(application.Options{
    Windows: application.WindowsOptions{
        AllowSingleSignOnUsingOSPrimaryAccount: true,
    },
})
```

The option is application-level because WebView2 requires all environments sharing a user data folder to be
created with the same options.

## Non-Goals

- Token acquisition or the Web Account Manager (WAM) broker. This only lets WebView2 do what Edge does.
- The other unbound environment options.

## Platform Considerations

- **Windows**: works for an Entra ID account on an Entra-joined or hybrid-joined device, or a Microsoft account
  signed in to Windows. Elsewhere it has no effect. The property has been in WebView2 since SDK 0.9.622.
- **macOS**: not needed. The Enterprise SSO extension already covers WKWebView on managed devices.
- **Linux**: WebKitGTK has no equivalent.

## Pros/Cons

### Pros

- Small, opt-in change that reuses an existing binding.

### Cons

- When enabled, any page loaded in the app can complete a Microsoft sign-in as the Windows user.

## Alternatives Considered

- A generic hook for environment options. More flexible, but it would expose internal COM types.
- Opening the system browser. Works, but it is a visible round trip, which is what silent sign-in avoids.

## Backwards Compatibility

None. The default is false, which matches current behaviour.

## Security and Privacy

WebView2 sends the Windows account's credentials to Microsoft identity endpoints, as Edge does. The docs will
warn against enabling it in apps that load untrusted content. Wails does not read or store credentials.

## Test Plan

On an Entra-joined Windows 11 machine, load a `prompt=none` authorize URL with the option on and off. On should
redirect with a code; off should return `interaction_required`. Also check that two windows with the option on
both open.

## Reference Implementation

#6612

## Maintenance Plan

This feature will be maintained and supported by the implementor and the Wails developers.

## Conclusion

A one-field option that brings WebView2 sign-in in line with Edge and WKWebView on managed devices.
