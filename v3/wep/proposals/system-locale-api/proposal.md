# Wails Enhancement Proposal (WEP)

## System locale API

**WEP Number**: (assigned on acceptance)  
**Status**: Draft  
**Author**: Morten Olsrud (@mortenolsrud)  
**Created**: 2026-11-06  
**Discussion**: [#5894](https://github.com/wailsapp/wails/pull/5894)  
**Implementor**: Morten Olsrud  
**Target**: Wails v3

## Summary

Add `application.PreferredLanguages() []string`: the languages the user has chosen for applications' user interfaces, most preferred first, as canonical BCP-47 tags, on Windows, macOS, Linux, iOS, Android and server builds. `application.PreferredLanguage()` returns the first of them, or `""` when there are none.

Both answer one question: which language should this app show? They do not report the regional format for dates, numbers and currency, which every platform keeps as a separate setting.

## Motivation

A Wails app has no portable way to learn the user's language from Go, and the obvious workarounds give wrong answers.

- **The environment is empty where it matters.** A macOS app launched from Finder, a Windows GUI process, and every iOS and Android app start with no `LANG` or `LC_*`. Code that reads them finds nothing on four of the six targets.
- **Android needs the bridge.** The language is only available from Java (`LocaleList`), and only Wails owns the JNI bridge. An app cannot add this on its own.
- **The "current locale" is a different setting.** `GetUserDefaultLocaleName` on Windows and `NSLocale.currentLocale` on Apple platforms report the regional *format*. Someone running English Windows in Norway gets `nb-NO`, and an app that picks its translation from that shows Norwegian to a user who chose English.
- **One value loses the fallback.** A user in Norway may list Northern Sami (`se-NO`) first and Norwegian Bokmål (`nb-NO`) second. An app with no Sami translation should show Bokmål, not its built-in default; only the list says so.
- **The webview is not enough.** `navigator.languages` depends on each webview engine's configuration, and it is not available to Go before a window exists, for menus, the tray, dialogs or notifications.

## Detailed Design

Two package-level functions in `v3/pkg/application`:

```go
// PreferredLanguages returns the languages the user has chosen for
// application user interfaces, most preferred first, as canonical BCP-47
// tags such as "nb-NO", "en-US" or "zh-Hans-CN". It returns an empty
// slice when the platform reports none; the app then applies its own default.
func PreferredLanguages() []string

// PreferredLanguage returns the first entry of PreferredLanguages, or ""
// when the list is empty.
func PreferredLanguage() string
```

**Rules**

1. **Meaning.** The display language from the operating system's language settings, never the regional format.
2. **Format.** Canonical BCP-47: hyphens, language in lower case, script in title case, region in upper case.
   - POSIX codesets (`.UTF-8`) and modifiers (`@euro`) are removed, except script modifiers: `sr_RS@latin` becomes `sr-Latn-RS`.
   - Entries that don't parse are dropped, and a duplicate keeps its first position.
3. **Unknown is empty.** There is no `"en"` fallback. An app whose own default is not English must be able to tell "unknown" from "English".
4. **No caching.** Each call asks the platform. Whether a running process sees a settings change depends on the platform, and the docs say so.
5. **One entry point everywhere.** `application.PreferredLanguages()` also works on iOS and Android, so `MobileManager` gets no new method.

**Implementation steps**

1. **The shared rules.** A file with no build tags holds the POSIX environment rules, normalisation and de-duplication, with unit tests that run on any host. Linux and server share this one copy.
2. **One raw function per platform.** Each platform file provides `platformPreferredLanguages() []string`, the raw strings the OS returns and nothing more, so all the logic stays testable. macOS has a cgo version and one for `purego` builds.
3. **Android.** The `WailsBridge.java` template gets `getPreferredLanguages()`. An app that keeps its own copy of the template from before this change must get an empty list, not a crash. The existing bridge call already does this: it clears the `NoSuchMethodError` and reports no result.
4. **Docs and an example.** A section on the Environment docs page, and a line in the environment example that prints the list.

**Normalisation** uses `golang.org/x/text/language`. v3 already requires v0.42.0 indirectly, and this makes it a direct dependency. A small hand-written canonicaliser could replace it if the maintainers prefer no new direct dependency.

**Impact on existing functionality:** none. This adds two functions and one bridge method, and changes nothing that exists.

## Non-Goals

- **The regional format** (dates, numbers, currency, first day of the week). It is a separate setting on every platform and deserves its own proposal if anyone needs it.
- **Translation.** Wails does not choose a translation or load message catalogues. Apps match the list against their own languages, for example with `language.NewMatcher`.
- **Change notifications.** No event fires when the user changes language while the app runs.
- **Setting a language.** No per-app language override is written or read beyond what the platform itself applies.
- **A frontend runtime call.** This WEP adds the Go API only. Apps can expose it through a bound service today, and a `@wailsio/runtime` call can follow separately.

## Platform Considerations

Every target is supported. Each reads the list the user edits in the OS's language settings, never the regional-format setting beside it.

| Platform | Source | Available from | Notes |
| --- | --- | --- | --- |
| Windows | `GetUserPreferredUILanguages(MUI_LANGUAGE_NAME)` | Windows Vista | Display languages, not "Regional format". Wrapped by `golang.org/x/sys/windows`, already a dependency. |
| macOS | `CFLocaleCopyPreferredLanguages` | macOS 10.5 | System Settings › General › Language & Region; the same list as `NSLocale.preferredLanguages`, which `purego` builds read through the Objective-C runtime. |
| iOS | `CFLocaleCopyPreferredLanguages` | iOS 2.0 | Settings › General › Language & Region. |
| Android | `LocaleList.getDefault().toLanguageTags()` | API 24 | Through the new bridge method. On API 21–23 (the template's `minSdk` is 21) one entry, from `Locale.getDefault().toLanguageTag()`. |
| Linux | `LANGUAGE`, `LC_ALL`, `LC_MESSAGES`, `LANG` | — | The rules below. Desktop sessions export the user's choice into the environment that launched apps inherit. |
| Server | As Linux | — | The *host's* environment, not the remote viewer's. An app serving browsers should use the request's `Accept-Language`. |

**Linux and server rules**, following GNU gettext:

1. The effective message locale is `LC_ALL`, else `LC_MESSAGES`, else `LANG`.
2. If that is unset, `C` or `POSIX` (with any codeset, so `C.UTF-8` too), the list is empty: the user asked for no localisation, and gettext ignores `LANGUAGE` in that case too.
3. Otherwise the entries of `LANGUAGE` (a colon-separated priority list such as `nb:no:en`) come first, in order, followed by the effective locale.

## Pros/Cons

**Pros**

- One call gives the same answer on six targets, including the two where an app cannot get it alone (Android) or easily (iOS).
- It answers the question apps actually ask, so translations follow the user's display language rather than their region.
- The list keeps the user's fallbacks, which a single value drops.
- No permissions, and no new dependency beyond making `x/text` direct.

**Cons**

- A new method in `WailsBridge.java`, a template many apps copy and own. The missing-method handling limits the damage, but apps must update their copy to get the list.
- On Linux the answer is only as good as the session environment. A process started without it, from a cron job or a bare `ssh` session for instance, reports an empty list.

## Alternatives Considered

- **A single string, as #5894 first proposed.** It drops the user's fallbacks: the `se-NO`, `nb-NO` user above gets the app's default instead of Bokmål. `PreferredLanguage()` keeps the one-value convenience on top of the list.
- **The current locale** (`GetUserDefaultLocaleName`, `NSLocale.currentLocale`, `Locale.getDefault()`), as #5894 first implemented. It answers "how should dates look", not "which language", and gives `nb-NO` for English Windows in Norway.
- **The name `SystemLocale()`**, as in #5894's first version and this proposal's title. On Apple platforms and Windows, "locale" names the regional-format setting this API deliberately does not return, so the name invites the wrong expectation. `PreferredLanguages` matches the platform APIs it wraps, `NSLocale.preferredLanguages` and `GetUserPreferredUILanguages`.
- **`navigator.languages` in the frontend.** It depends on each webview engine, and Go code that runs before a window exists cannot use it.
- **Leave it to apps.** Each app would write its own cgo and Windows calls, and none can reach Android's list without changing the bridge Wails owns.
- **Return the best match for the app's languages.** Matching is the app's policy; `language.NewMatcher` already does it well on top of the list.

## Backwards Compatibility

Nothing breaks. The change only adds package-level functions; no exported interface changes, so code that implements `MobileManager` elsewhere is unaffected. #5894's first version added `SystemLocale()` to `MobileManager`; that method is gone, and no release carried it.

An app with its own copy of `WailsBridge.java` from before this change keeps working and gets an empty list, under the missing-method rule in the design.

## Security and Privacy

No permission is needed on any platform, and nothing leaves the device. The language list can help fingerprint a user, but the app already runs with the user's rights, and its webview exposes the same list as `navigator.languages`. There are no other effects.

## Test Plan

**Unit tests, on any host.** The shared file takes `getenv` and the raw platform strings as parameters, so the tests never touch the process environment.

| Input | Expected |
| --- | --- |
| `LANG=nb_NO.UTF-8` | `["nb-NO"]` |
| `LANG=sr_RS@latin` | `["sr-Latn-RS"]` |
| `LANG=de_DE@euro` | `["de-DE"]` |
| `LANG` unset, `C`, `POSIX` or `C.UTF-8` | `[]` |
| `LANGUAGE=nb:en`, `LANG=en_US.UTF-8` | `["nb", "en", "en-US"]` |
| `LC_ALL=C`, `LANGUAGE=nb` | `[]` |
| `LC_ALL=en_GB`, `LC_MESSAGES=nb_NO`, `LANG=de_DE` | `["en-GB"]` |
| Platform strings `zh-hans-cn`, `nb-NO`, `nb-NO`, `!!` | `["zh-Hans-CN", "nb-NO"]` |

**Manual checks per platform,** with two languages set in a known order, on the oldest supported OS version as well as the current one. #5894's first version called an Apple API newer than its availability check, and only an old OS shows that.

- **Windows 11:** English display language with Norwegian regional format must give `en-*` first.
- **macOS 13 and current:** app launched from Finder, so with no `LANG`; both the cgo and the `purego` build.
- **iOS 16 and current.**
- **Android API 23 and current:** one entry on 23. An app whose `WailsBridge.java` lacks the method must get `[]` and not crash.
- **Ubuntu with GNOME:** the order set in Settings.
- **Server build:** `LANG` set on the host.

## Reference Implementation

[#5894](https://github.com/wailsapp/wails/pull/5894) implements this proposal:

- `preferred_languages.go`: the two public functions, the POSIX rules and the normalisation, with no build tags. `preferred_languages_test.go` covers the table above and more, on any host.
- One `platformPreferredLanguages()` per platform: Windows (`x/sys/windows`), macOS and iOS (CoreFoundation through cgo), macOS `purego` (`NSLocale` through the Objective-C runtime), Android (the new bridge method), and Linux, server builds and other targets (the POSIX rules).
- `getPreferredLanguages()` in the `WailsBridge.java` template.
- A "Preferred Languages" section on the Environment docs page, and a line in the environment example.

Its description lists what has been tested so far.

**Prior art:** the desktop half of this design runs in a Wails v3 dictionary app from a leading Norwegian publishing house, currently in beta. It reads the preferred-language lists on Windows and macOS, and keeps the POSIX rules in a shared, host-tested file.

## Maintenance Plan

The code is small: one thin function per platform and one shared file. The OS APIs it calls are stable and old, from Windows Vista, macOS 10.5, iOS 2.0 and Android 7. Unit tests on any host pin the behaviour, so regressions show up in normal CI. The author will maintain it and answer its issues.

## Conclusion

Apps need to know which language to show, and today each must find out per platform, or cannot at all on Android. A list of preferred display languages in BCP-47, with a one-value convenience, answers that the same way on every target. It adds no breaking change and needs no permission.
