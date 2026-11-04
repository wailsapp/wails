---
title: "iOS"
description: "Сборка и запуск приложений Wails на iOS: настройка, симулятор, сборки для устройств, конфигурация и нативные функции"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="Экспериментальная функция"}
Поддержка iOS является экспериментальной и может измениться в будущих выпусках.

@end

@note{type="tip"}
Впервые разрабатываете мобильное приложение с Wails? Начните с пошагового руководства [«Ваше первое мобильное приложение» →](/guides/mobile/first-mobile-app/), а затем вернитесь сюда за полной справочной информацией.

@end

Приложения Wails v3 выполняются в iOS как полностью нативные приложения — и самое приятное, что они работают *точно* так же, как настольная версия. Тот же бэкенд на Go, тот же фронтенд, те же `@wailsio/runtime`: привязки сервисов, события, диалоговые окна и буфер обмена ведут себя одинаково, и для мобильной платформы не требуется **никакой** специальной перенастройки. Не нужны отдельная кодовая база для мобильной версии, слой портирования или специальный API, который пришлось бы осваивать: существующее приложение Wails просто запускается на iOS. Портирование действительно проходит без затруднений: перенесите приложение как есть и выпускайте его.

Одна и та же `main.go` используется для сборки как настольной версии, так и версии для iOS; параметры, относящиеся к iOS, настраиваются через `application.Options.IOS`.

## Требования

- macOS с установленной **полной версией Xcode** (одних инструментов командной строки недостаточно) — команда `wails3 doctor` показывает доступные ей SDK для iOS
- Go 1.26+ и npm

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

В отладочных сборках WebView можно инспектировать через меню Develop в Safari.

## Упаковка

```bash
wails3 task ios:package             # production .app for the simulator
wails3 task ios:deploy-simulator    # install + launch it
```

Это оптимизированные производственные сборки с удалённой отладочной информацией.

## Сборки для устройств

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device` создаёт сборку для физического устройства. Разрешения берутся из `build/ios/entitlements.plist` и применяются только к сборкам для устройств — добавьте ключи возможностей, необходимые приложению.

@note{type="tip"}
Чтобы автоматически управлять подписью и профилями подготовки, а также создавать архивы для App Store, откройте созданный проект Xcode с помощью `wails3 task ios:xcode` и выполните сборку в Xcode.

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

К параметрам запуска (`application.Options.IOS`) относятся `DisableScroll`, `DisableBounce`, `DisableScrollIndicators`, `DisableInputAccessoryView`, `EnableBackForwardNavigationGestures`, `DisableLinkPreview`, `EnableInlineMediaPlayback`, `EnableAutoplayWithoutUserAction`, `DisableInspectable`, `UserAgent`, `ApplicationNameForUserAgent`, `BackgroundColour`, а также нативные нижние вкладки, задаваемые через `EnableNativeTabs` + `NativeTabsItems`.

## Нативные функции

Специальные возможности iOS доступны через `application.IOS`. Вызывайте их из Go в файле `//go:build ios`, чтобы общий код оставался независимым от платформы. В Android тот же набор доступен через `application.Android`.

Однократные действия завершаются немедленно:

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

Вспомогательные функции запросов возвращают результаты в формате JSON: `SafeAreaJSON()`, `AppInfoJSON()`, `PowerJSON()`, `NetworkJSON()`, `StorageJSON()`, `GetOrientation()`, `GetBrightness()`. `StoragePath()` возвращает абсолютный путь к каталогу Application Support приложения — подходящему месту для баз данных и других постоянных файлов (аналогу `getFilesDir()` из Android в iOS). Каталог создаётся при первом обращении; если создать его не удаётся, `StoragePath()` возвращает пустую строку, поэтому перед использованием проверьте значение на `""`.

### События

Всё, что завершается позднее — запрос разрешения, поток данных датчика или съёмка камерой, — передаёт результат не как возвращаемое значение, а как **событие**, которое можно прослушивать в Go или во фронтенде. Имена возможностей, общих с Android, имеют префикс `common:`, а доступных только в iOS — `ios:`.

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

| Событие | Инициируется | Полезная нагрузка |
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

Комплексный пример в `v3/examples/mobile` обеспечивает сквозную интеграцию всех перечисленных выше функций.

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

Встроенный `@wailsio/runtime` также предоставляет небольшое пространство имён iOS для **фронтенда**:

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

Выбор нативной нижней вкладки поступает как событие `nativeTabSelected` в `window`.

## Состояние поддержки

| Область | Состояние |
| --- | --- |
| Отрисовка фронтенда и ресурсы | ✅ |
| Привязки сервисов, события (в обоих направлениях) | ✅ |
| Диалоговые окна сообщений | ✅ |
| Диалоговые окна открытия файла, файлов или каталога | ✅ Импортируются как копии в песочнице |
| Диалоговые окна сохранения файла | ❌ Вместо этого записывайте данные внутри песочницы приложения |
| Буфер обмена | ✅ |
| API экранов | ✅ Рабочая область учитывает безопасную зону |
| События жизненного цикла | ✅ |
| Геометрия окна, меню, системный трей | В iOS не выполняют никаких действий |
| Несколько окон | Отображается только первое окно |

## Примечания по портированию

- Код для настольных систем собирается для iOS без изменений — вызовы для окон, меню и системного трея просто ничего не делают.
- Замените диалоговые окна сохранения файлов записью в песочницу приложения с последующим предоставлением общего доступа.
- Создавайте адаптивный фронтенд; безопасные области обрабатываются автоматически.
