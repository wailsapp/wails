---
title: "모바일 API"
description: "크로스 플랫폼 application.Mobile 관리자 — iOS와 Android가 공유하는 네이티브 모바일 기능을 위한 빌드 가드 적용 단일 진입점"
slug: "guides/mobile/mobile-api"
sourcePath: "guides/mobile/mobile-api.md"
---

@note{type="caution" title="실험적 기능"}
모바일 지원은 실험적 기능이며 향후 릴리스에서 변경될 수 있습니다.

@end

네이티브 모바일 기능은 다음 두 가지 방식으로 제공됩니다.

- **플랫폼별 관리자** — `//go:build ios` 파일의 `application.IOS`와 `//go:build android` 파일의 `application.Android`입니다. 플랫폼별 기능에는 이 관리자들을 사용하세요. 플랫폼별 API 전체는 [iOS](/guides/mobile/ios/) 및 [Android](/guides/mobile/android/) 참조 문서를 확인하세요.
- **`application.Mobile`** — 두 플랫폼에서 동일하게 동작하는 기능의 일부를 제공하는, 빌드 가드가 적용된 단일 관리자입니다. 모든 플랫폼에서 컴파일되고 실행되는 하나의 코드 경로가 필요할 때 사용하세요.

## `application.Mobile`

`application.Mobile`는 iOS에서는 `IOS`로, Android에서는 `Android`로 디스패치하고, 데스크톱에서는 아무 작업도 하지 않는 스텁으로 디스패치합니다. 빌드 제약 조건이 없으므로 별도의 `//go:build` 파일 없이 일반적인 플랫폼 독립적 Go 코드에서 호출할 수 있습니다.

```go
// Works in any file, on any target.
// On desktop this returns "" (no-op); on device it returns the real path.
dbDir := application.Mobile.StoragePath()
if dbDir == "" {
    // Off-device, or the directory could not be created — handle accordingly.
    return
}
db, _ := sql.Open("sqlite", filepath.Join(dbDir, "app.db"))
```

`StoragePath()`는 앱의 비공개 파일 디렉터리 절대 경로를 반환합니다. Android에서는 `getFilesDir()`, iOS에서는 Application Support 디렉터리이며, 데이터베이스와 기타 영구 파일을 저장할 권장 위치입니다. 데스크톱에서는 빈 문자열을 반환하며, 기기에서도 디렉터리를 사용할 수 없으면(iOS에서는 디렉터리를 생성할 수 없으면) 빈 문자열을 반환하므로 사용하기 전에 `""`인지 확인하세요.

@note{type="note"}
기기 외부 환경(데스크톱 빌드)에서는 모든 `Mobile` 메서드가 아무 작업도 하지 않으며, 모든 쿼리는 해당 형식의 제로 값(문자열의 경우 `""`)을 반환합니다. 따라서 크로스 플랫폼 코드에서 조건 없이 `application.Mobile.*`를 호출할 수 있습니다. 데스크톱에서도 실제 경로가 필요하면 플랫폼에 따라 분기하여 `os.UserConfigDir()` 또는 이와 유사한 방식으로 대체하세요.

@end

## 기능

`Mobile` 관리자는 iOS와 Android에서 시그니처가 동일한 다음 기능을 제공합니다.

| 기능 | API | 참고 |
| --- | --- | --- |
| 공유 시트 | `Mobile.Share(json)` | `{text, url}` |
| 외부에서 URL 열기 | `Mobile.OpenURL(url)` | 시스템 브라우저 |
| 화면 켜짐 유지 | `Mobile.SetKeepAwake(bool)` |  |
| 토치 / 손전등 | `Mobile.SetTorch(bool)` | → `common:torch` |
| 안전 영역 인셋 | `Mobile.SafeAreaJSON()` | `{top,bottom,left,right}` |
| 앱 정보 | `Mobile.AppInfoJSON()` | `{name,version,build,bundleId}` |
| 화면 방향 잠금 | `Mobile.SetOrientation(mode)` | `portrait` / `landscape` / `auto` |
| 상태 표시줄 | `Mobile.SetStatusBar(json)` | 스타일 + 표시 여부 |
| 저장 공간 정보 | `Mobile.StorageJSON()` | `{free,total}`바이트 |
| 저장소 경로 | `Mobile.StoragePath()` | 앱 비공개 파일 디렉터리 |
| 전원 / 배터리 | `Mobile.PowerJSON()` | `{level,charging,lowPower}` |
| 네트워크 상태 | `Mobile.NetworkJSON()` | `{connected,type}` |
| 생체 인식 | `Mobile.BiometricAuthenticate(reason)` | → `common:biometric` |
| 보안 저장소 | `Mobile.SecureGet(key)` / `Mobile.SecureDelete(key)` | Keychain / `EncryptedSharedPreferences` |
| 위치 정보 | `Mobile.GetLocation()` | 일회성 → `common:location` |
| 햅틱 | `Mobile.Haptic(type)` | 충격 / 알림 / 선택 |
| 가속도계 | `Mobile.SetMotion(bool)` | → `common:motion` |
| 근접 센서 | `Mobile.SetProximity(bool)` | → `common:proximity` |
| 텍스트 음성 변환 | `Mobile.Speak(text)` / `Mobile.StopSpeak()` |  |
| 키보드 인셋 | `Mobile.SetKeyboardWatch(bool)` | → `common:keyboard` |
| 화면 캡처 | `Mobile.SetScreenProtect(bool)` | → `common:screenCapture` |
| 카메라 | `Mobile.CapturePhoto()` / `Mobile.CaptureVideo()` | → `common:capture` |

비동기 결과는 플랫폼별 관리자와 정확히 동일하게 `common:*` 이벤트로 전달됩니다. 페이로드에 대해서는 [이벤트](/guides/mobile/ios/#events)를 참조하세요.

## 플랫폼별로 유지되는 항목

iOS와 Android에서 형태가 다른 기능은 `Mobile`에 **포함되지 않습니다**. 빌드 태그가 지정된 파일에서 `application.IOS` / `application.Android`를 통해 호출하세요.

| 항목 | iOS | Android |
| --- | --- | --- |
| 밝기(설정) | `IOS.SetBrightness(0.0-1.0)` | `Android.SetBrightness(0-100)` |
| 밝기 / 방향(가져오기) | `IOS.GetBrightness()` / `IOS.GetOrientation()` | `Android.BrightnessJSON()` / `Android.OrientationJSON()` |
| 로컬 알림 | `IOS.PostNotification(json)` | `Android.Notify(json)` |
| 보안 저장소(쓰기) | `IOS.SecureSet(key, value)` | `Android.SecureSet(json)` |
| 백그라운드 실행 | `IOS.BeginBackgroundTask` / `EndBackgroundTask` | `Android.StartForegroundService` / `StopForegroundService` |

@note{type="tip"}
`MobileManager` 인터페이스는 `application.Mobile`의 기반이 되는 계약입니다. 두 플랫폼 관리자가 모두 이 인터페이스를 충족해야 하므로, 위에 나열된 모든 메서드는 iOS와 Android에서 동일한 시그니처를 유지합니다. 두 시그니처가 조금이라도 달라지면 해당 플랫폼 빌드의 컴파일이 실패합니다.

@end
