---
title: "Мобильный API"
description: "Кроссплатформенный менеджер application.Mobile — единая точка входа с защитой на уровне сборки для нативных мобильных возможностей, общих для iOS и Android"
slug: "guides/mobile/mobile-api"
sourcePath: "guides/mobile/mobile-api.md"
---

@note{type="caution" title="Экспериментальная функция"}
Поддержка мобильных платформ является экспериментальной и может измениться в будущих выпусках.

@end

Нативные мобильные возможности доступны двумя способами:

- **Менеджеры для отдельных платформ** — `application.IOS` (в файлах `//go:build ios`) и `application.Android` (в файлах `//go:build android`). Используйте их для всего, что зависит от конкретной платформы. Полный набор API для каждой платформы приведён в справочниках по [iOS](/guides/mobile/ios/) и [Android](/guides/mobile/android/).
- **`application.Mobile`** — единый менеджер с защитой на уровне сборки, охватывающий подмножество возможностей, которые одинаково работают на обеих платформах. Используйте его, если нужен единый путь выполнения кода, который компилируется и работает везде.

## `application.Mobile`

`application.Mobile` перенаправляет вызовы к `IOS` в iOS, к `Android` в Android и к заглушке, которая ничего не делает, в настольных системах. Поскольку ограничений сборки у него нет, его можно вызывать из обычного, не зависящего от платформы кода Go — собственные файлы `//go:build` не нужны:

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

`StoragePath()` возвращает абсолютный путь к каталогу личных файлов приложения — `getFilesDir()` в Android или каталогу Application Support в iOS. Это рекомендуемое место для баз данных и других постоянных файлов. Метод возвращает пустую строку в настольных системах, а на устройстве — если каталог недоступен (в iOS — если его не удаётся создать), поэтому перед использованием проверьте результат на `""`.

@note{type="note"}
Вне устройства (в настольных сборках) каждый метод `Mobile` ничего не делает, а каждый запрос возвращает нулевое значение своего типа (`""` для строк). Благодаря этому кроссплатформенный код может безусловно вызывать `application.Mobile.*`. Если реальный путь нужен и в настольной системе, добавьте ветвление по платформе и используйте в качестве запасного варианта `os.UserConfigDir()` или аналогичный API.

@end

## Возможности

Менеджер `Mobile` предоставляет возможности с одинаковыми сигнатурами в iOS и Android:

| Возможность | API | Примечания |
| --- | --- | --- |
| Системное меню «Поделиться» | `Mobile.Share(json)` | `{text, url}` |
| Открытие URL во внешнем приложении | `Mobile.OpenURL(url)` | Системный браузер |
| Предотвращение отключения экрана | `Mobile.SetKeepAwake(bool)` |  |
| Фонарик | `Mobile.SetTorch(bool)` | → `common:torch` |
| Отступы безопасной области | `Mobile.SafeAreaJSON()` | `{top,bottom,left,right}` |
| Сведения о приложении | `Mobile.AppInfoJSON()` | `{name,version,build,bundleId}` |
| Блокировка ориентации | `Mobile.SetOrientation(mode)` | `portrait` / `landscape` / `auto` |
| Строка состояния | `Mobile.SetStatusBar(json)` | стиль и видимость |
| Сведения о хранилище | `Mobile.StorageJSON()` | `{free,total}` байт |
| Путь к хранилищу | `Mobile.StoragePath()` | Каталог файлов, доступный только приложению |
| Питание и аккумулятор | `Mobile.PowerJSON()` | `{level,charging,lowPower}` |
| Состояние сети | `Mobile.NetworkJSON()` | `{connected,type}` |
| Биометрия | `Mobile.BiometricAuthenticate(reason)` | → `common:biometric` |
| Защищённое хранилище | `Mobile.SecureGet(key)` / `Mobile.SecureDelete(key)` | Keychain / `EncryptedSharedPreferences` |
| Геолокация | `Mobile.GetLocation()` | однократный запрос → `common:location` |
| Тактильная обратная связь | `Mobile.Haptic(type)` | воздействие / уведомление / выбор |
| Акселерометр | `Mobile.SetMotion(bool)` | → `common:motion` |
| Датчик приближения | `Mobile.SetProximity(bool)` | → `common:proximity` |
| Синтез речи | `Mobile.Speak(text)` / `Mobile.StopSpeak()` |  |
| Отступы экранной клавиатуры | `Mobile.SetKeyboardWatch(bool)` | → `common:keyboard` |
| Захват экрана | `Mobile.SetScreenProtect(bool)` | → `common:screenCapture` |
| Камера | `Mobile.CapturePhoto()` / `Mobile.CaptureVideo()` | → `common:capture` |

Асинхронные результаты поступают в виде событий `common:*`, точно так же, как при использовании менеджеров отдельных платформ. Описание полезной нагрузки см. в разделе [События](/guides/mobile/ios/#events).

## Что остаётся платформозависимым

Возможности, форма которых различается в iOS и Android, **не** представлены в `Mobile`; вызывайте их через `application.IOS` / `application.Android` из файла с тегом сборки:

| Аспект | iOS | Android |
| --- | --- | --- |
| Яркость (установка) | `IOS.SetBrightness(0.0-1.0)` | `Android.SetBrightness(0-100)` |
| Яркость / ориентация (получение) | `IOS.GetBrightness()` / `IOS.GetOrientation()` | `Android.BrightnessJSON()` / `Android.OrientationJSON()` |
| Локальное уведомление | `IOS.PostNotification(json)` | `Android.Notify(json)` |
| Защищённое хранилище (запись) | `IOS.SecureSet(key, value)` | `Android.SecureSet(json)` |
| Выполнение в фоновом режиме | `IOS.BeginBackgroundTask` / `EndBackgroundTask` | `Android.StartForegroundService` / `StopForegroundService` |

@note{type="tip"}
Интерфейс `MobileManager` определяет контракт, лежащий в основе `application.Mobile`. Поскольку оба платформенных менеджера должны ему соответствовать, любой из перечисленных выше методов гарантированно сохраняет одинаковую сигнатуру в iOS и Android. Если сигнатуры когда-либо разойдутся, сборка для соответствующей платформы завершится ошибкой компиляции.

@end
