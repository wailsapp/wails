---
title: "iOS"
description: "Сборка и запуск приложений Wails на iOS: настройка, симулятор, сборки для устройств, конфигурация и нативные возможности"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="Экспериментальная возможность"}
Поддержка iOS является экспериментальной и может измениться в будущих выпусках.

@end

@note{type="tip"}
Впервые разрабатываете мобильное приложение с Wails? Начните с пошагового руководства [«Ваше первое мобильное приложение» →](/guides/mobile/first-mobile-app/), а затем вернитесь сюда за полной справочной информацией.

@end

Приложения Wails v3 работают на iOS как полностью нативные приложения — и, что особенно важно, они работают *точно* так же, как настольная версия. Тот же бэкенд на Go, тот же фронтенд и те же `@wailsio/runtime`: привязки сервисов, события, диалоговые окна и буфер обмена работают одинаково, и для мобильной платформы не требуется **никаких** специальных переделок. Не нужны ни отдельная кодовая база для мобильной версии, ни слой портирования, ни специальный API, который пришлось бы осваивать: существующее приложение Wails просто запускается на iOS. Портирование действительно не требует дополнительных усилий: перенесите приложение как есть и выпускайте его.

Один и тот же файл `main.go` компилируется как для настольных платформ, так и для iOS; особенности iOS настраиваются через `application.Options.IOS`.

## Требования

- macOS с установленным **полным пакетом Xcode** (одних инструментов командной строки недостаточно) — `wails3 doctor` показывает найденные пакеты SDK для iOS
- Go 1.25+ и npm

## Симулятор

Из каталога проекта:

```bash
wails3 task ios:run
```

Эта команда собирает приложение, запускает симулятор, если он ещё не запущен, и открывает в нём приложение.

Полезные сопутствующие команды:

```bash
wails3 task ios:logs:dev    # stream the app's logs from the simulator
wails3 task ios:xcode       # open the generated Xcode project
```

В отладочных сборках содержимое WebView можно исследовать через меню Develop в Safari.

## Создание пакетов

```bash
wails3 task ios:package             # production .app for the simulator
wails3 task ios:deploy-simulator    # install + launch it
```

Это оптимизированные производственные сборки без отладочной информации.

## Сборки для устройств

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device` выполняет сборку для физического устройства. Права приложения берутся из `build/ios/entitlements.plist` и применяются только к сборкам для устройств — добавьте ключи возможностей, необходимые вашему приложению.

@note{type="tip"}
Чтобы автоматически управлять подписью, профилями подготовки и архивами App Store, откройте созданный проект Xcode с помощью `wails3 task ios:xcode` и выполните сборку в Xcode.

@end

## Конфигурация

`build/config.yml`:

```yaml
ios:
  bundleID: com.example.myapp
  displayName: My App
  version: 1.0.0
  minIOSVersion: "15.0"
```

Параметры запуска (`application.Options.IOS`) включают `DisableScroll`, `DisableBounce`, `DisableScrollIndicators`, `DisableInputAccessoryView`, `EnableBackForwardNavigationGestures`, `DisableLinkPreview`, `EnableInlineMediaPlayback`, `EnableAutoplayWithoutUserAction`, `DisableInspectable`, `UserAgent`, `ApplicationNameForUserAgent`, `BackgroundColour`, а также нативные нижние вкладки, задаваемые через `EnableNativeTabs` и `NativeTabsItems`.

## Нативные возможности

Возможности, специфичные для iOS, доступны через `application.IOS`. Вызывайте их из Go в файле с ограничением сборки `//go:build ios`, чтобы общий код не зависел от платформы. В Android тот же набор доступен через `application.Android`.

Вызовы однократных действий немедленно возвращают управление:

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

Вспомогательные функции запросов возвращают результаты в формате JSON: `SafeAreaJSON()`, `AppInfoJSON()`, `PowerJSON()`, `NetworkJSON()`, `StorageJSON()`, `GetOrientation()`, `GetBrightness()`. `StoragePath()` возвращает абсолютный путь к каталогу Application Support приложения — удобному месту для баз данных и других постоянных файлов (аналогу `getFilesDir()` из Android в iOS). Каталог создаётся при первом обращении; если создать его не удаётся, `StoragePath()` возвращает пустую строку, поэтому перед использованием проверьте значение на `""`.

### События

Всё, что завершается позднее, — запрос разрешения, поток данных датчика или съёмка камерой — передаёт результат как **событие**, а не как возвращаемое значение. Такое событие можно отслеживать в Go или во фронтенде. Имена событий для возможностей, общих с Android, имеют префикс `common:`, а для возможностей, доступных только в iOS, — префикс `ios:`.

```go
// Go
app.Event.On("common:location", func(e *application.CustomEvent) {
    // e.Data -> {"lat":..,"lng":..,"accuracy":..} or {"error":..}
})
```

```js
// frontend
import { Events } from "@wailsio/runtime";
Events.On("common:notification", (e) => { /* {ok, scheduled, presented, tapped, error} */ });
```

| Событие | Причина срабатывания | Полезная нагрузка |
| --- | --- | --- |
| `common:biometric` | `BiometricAuthenticate(reason)` | `{ok, error}` |
| `common:location` | `GetLocation()` | `{lat, lng, accuracy}` / `{error}` |
| `common:motion` | `SetMotion(true)` | `{x, y, z}` |
| `common:proximity` | `SetProximity(true)` | `{near}` |
| `common:keyboard` | `SetKeyboardWatch(true)` | `{visible, height}` |
| `common:torch` | `SetTorch(bool)` | `{on, available}` |
| `common:notification` | `PostNotification(json)` | `{ok, scheduled, presented, tapped, error}` |
| `common:capture` | `CapturePhoto()` / `CaptureVideo()` | `{type, path, size, thumb}` |
| `common:screenCapture` | `SetScreenProtect(true)` | `{screenshot, recording}` |
| `ios:backgroundTask` | `BeginBackgroundTask(seconds)` | `{message, granted}` |

Комплексный пример в `v3/examples/mobile` объединяет все перечисленные выше возможности в полностью работающий сценарий.

## Управление WebView

Некоторые параметры поведения WebView также можно изменять во время выполнения из Go:

```go
application.IOS.SetScrollEnabled(false)
application.IOS.SetBounceEnabled(false)
application.IOS.SetScrollIndicatorsEnabled(false)
application.IOS.SetBackForwardGesturesEnabled(true)
application.IOS.SetLinkPreviewEnabled(false)
application.IOS.SetInspectableEnabled(true)
application.IOS.SetCustomUserAgent("MyApp/1.0")
```

Встроенный `@wailsio/runtime` также предоставляет во **фронтенде** небольшое пространство имён iOS:

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

Выбор вкладки на нативной нижней панели поступает как событие `nativeTabSelected` в `window`.

## Статус поддержки

| Область | Статус |
| --- | --- |
| Отрисовка фронтенда и ресурсы | ✅ |
| Привязки сервисов, события (в обоих направлениях) | ✅ |
| Диалоговые окна сообщений | ✅ |
| Диалоговые окна открытия файла, нескольких файлов или каталога | ✅ Импортируются как копии в песочнице |
| Диалоговые окна сохранения файла | ❌ Вместо этого записывайте данные в песочницу приложения |
| Буфер обмена | ✅ |
| API экранов | ✅ Включает рабочую область с учётом безопасной зоны |
| События жизненного цикла | ✅ |
| Геометрия окна, меню, системный трей | В iOS ничего не делают |
| Несколько окон | Отображается только первое окно |

## Примечания по переносу

- Код для настольных систем собирается для iOS без изменений — вызовы, относящиеся к окнам, меню и системному трею, просто ничего не делают.
- Замените диалоговые окна сохранения файлов записью в песочницу приложения с последующей отправкой файла.
- Используйте адаптивный дизайн фронтенда; безопасные зоны учитываются автоматически.
