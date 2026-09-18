---
title: "Ваше первое мобильное приложение"
description: "Запустите приложение Wails в iOS Simulator или Android Emulator за считаные минуты"
slug: "guides/mobile/first-mobile-app"
sourcePath: "guides/mobile/first-mobile-app.md"
---

В этом руководстве показано, как запустить обычное настольное приложение Wails в iOS Simulator или Android Emulator. **Изменять код Go не нужно.** Один и тот же `main.go` используется для сборки под все целевые платформы.

**Время выполнения:** 15–30 минут (бо́льшая часть этого времени при первом запуске уходит на установку инструментов)

## Начните с настольного проекта

Если у вас ещё нет проекта, создайте новый:

```bash
wails3 init -n mymobileapp
cd mymobileapp
```

Сначала убедитесь, что настольное приложение работает:

```bash
wails3 dev
```

Когда приложение откроется, закройте его и продолжайте. Всё, что работает на настольной платформе, работает и на мобильной — в рамках этого руководства вам не придётся изменять `main.go` или какой-либо код Go.

---

## Выберите платформу

@tabs{sync-key="mobile-platform"}
[iOS Simulator]
### Требования

- **macOS** (сборки для iOS доступны только в macOS)
- **Полная версия Xcode** — одних инструментов командной строки недостаточно. Установите её из App Store, затем выполните:
  ```bash
  sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
  sudo xcodebuild -license accept
  ```


- **Go 1.25+** и **npm** (они уже установлены, если вы запускали `wails3 init`)

Для проверки выполните `wails3 doctor` — команда выведет список обнаруженных SDK для iOS.

### Запустите в симуляторе

@steps
### Запустите приложение
```bash
wails3 task ios:run
```

Готово: эта команда соберёт приложение, запустит симулятор, если он ещё не запущен, и откроет в нём приложение.

@note{type="tip"}
Первый запуск занимает несколько минут, поскольку выполняются компиляция и кэширование фреймворка Wails для iOS. Все последующие запуски проходят гораздо быстрее.

@end

После запуска ваше настольное приложение без каких-либо изменений будет работать в iOS Simulator — с тем же `main.go` и тем же фронтендом:

![Стандартное приложение Wails, запущенное в iOS Simulator](/assets/ios-simulator-first-app.png)

### Просматривайте журналы в реальном времени
В отдельном терминале выполните:

```bash
wails3 task ios:logs:dev
```

Команда непрерывно выводит журнал симулятора, отфильтрованный по вашему приложению. Здесь отображаются выходные данные `fmt.Println` и `log.Println`.

### Проверьте WebView
В Safari выберите **Разработка → Симулятор → ваше приложение**. Доступны все возможности веб-инспектора: консоль, отладчик, панель сети и всё остальное.

### Внесите изменение
Измените любой файл фронтенда (`frontend/src/main.js`, `index.html` и т. д.) и снова выполните `wails3 task ios:run`. Wails повторно соберёт фронтенд и перезапустит приложение.

После изменений в Go также снова выполните `wails3 task ios:run`. Go выполняет инкрементную перекомпиляцию, поэтому повторно собираются только изменённые пакеты.

@end

### Откройте в Xcode (необязательно)

```bash
wails3 task ios:xcode
```

Команда открывает `build/ios/` в Xcode. Xcode можно использовать для развёртывания на устройствах, расширенного профилирования и управления профилями подготовки. Wails заново создаёт проект Xcode при каждой сборке, поэтому не изменяйте созданные файлы напрямую.

[Android Emulator]
### Требования

Вам потребуются **Android SDK**, **NDK** и **JDK**. Проще всего установить их с помощью Android Studio, но можно использовать и инструменты командной строки:

@steps
### Установите инструменты командной строки Android
Скачайте архив с [developer.android.com/studio#command-line-tools-only](https://developer.android.com/studio#command-line-tools-only) и распакуйте его в `~/android-sdk/cmdline-tools/latest/`.

### Установите компоненты SDK
```bash
sdkmanager "platform-tools" \
           "platforms;android-35" \
           "build-tools;35.0.0" \
           "ndk;26.3.11579264" \
           "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
```

### Создайте эмулятор
```bash
avdmanager create avd \
  --name wails \
  --package "system-images;android-35;google_apis;arm64-v8a" \
  --device pixel_7
```

### Задайте переменные окружения
Добавьте в `~/.zshrc` или `~/.bashrc`:

```bash
export ANDROID_HOME=~/android-sdk
export ANDROID_SDK_ROOT=~/android-sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/cmdline-tools/latest/bin
```

Перезагрузите конфигурацию: `source ~/.zshrc`

### Установите JDK
```bash
# macOS
brew install openjdk@21
export JAVA_HOME=$(brew --prefix openjdk@21)

# Ubuntu/Debian
sudo apt install openjdk-21-jdk
export JAVA_HOME=/usr/lib/jvm/java-21-openjdk-amd64

# Windows (scoop)
scoop install openjdk21
```

@end

Выполните `wails3 doctor`, чтобы убедиться, что все компоненты обнаружены.

### Запустите в эмуляторе

@steps
### Запустите приложение
```bash
wails3 task android:run
```

При первом запуске команда выполняет следующие действия:

- Запускает эмулятор, если ни один эмулятор ещё не запущен
- Создаёт привязки и собирает фронтенд
- Компилирует код Go в `libwails.so` с помощью кросс-компилятора NDK
- Собирает отладочный APK с помощью Gradle
- Устанавливает и запускает его в эмуляторе

@note{type="tip"}
При первой сборке загружается Gradle и компилируется набор инструментов NDK — это займёт около 5–10 минут. Последующие сборки выполняются инкрементно и занимают менее минуты.

@end

### Просматривайте журналы в реальном времени
В отдельном терминале выполните:

```bash
wails3 task android:logs
```

Команда запускает `adb logcat` с фильтрацией по вашему приложению. Здесь отображаются выходные данные `fmt.Println`.

### Проверьте WebView
Откройте Chrome и перейдите по адресу `chrome://inspect`. WebView вашего приложения появится в разделе **Удалённая цель** — нажмите **проверить**, чтобы открыть DevTools.

### Внесите изменение
Измените любой файл и снова выполните `wails3 task android:run`. Благодаря инкрементной сборке Gradle перекомпилируется только изменённый код.

@end

@end

---

## Что произошло

Ваш `main.go` совсем не изменился. Wails выполнил всё необходимое:

- **Система сборки** — файл `Taskfile.yml` в вашем проекте содержит задачи `ios:*` и `android:*`, которые запускают набор инструментов для соответствующей платформы.
- **Кросс-компиляция Go** — с помощью `GOOS=ios` или `GOOS=android` с подходящими значениями `GOARCH` и sysroot.
- **Нативное хост-приложение** — сгенерированный проект Xcode (iOS) или Gradle (Android), который включает скомпилированный код Go и размещает WebView.
- **Предоставление ресурсов** — содержимое `frontend/dist/` встраивается в двоичный файл Go и предоставляется внутри процесса. Локальный сервер не требуется.

---

## Адаптируйте приложение для мобильных устройств

Приложение уже работает, но на экране телефона выглядит как настольное. Несколько небольших изменений заметно улучшат результат.

### Адаптивный CSS

Экраны мобильных устройств уже и используют другие способы ввода. В файле `frontend/public/style.css` (или его аналоге):

```css
/* Prevent horizontal scrolling */
body {
  overflow-x: hidden;
}

/* Touch-friendly tap targets */
button {
  min-height: 44px;
  min-width: 44px;
}

/* Respect the iOS safe area (notch, home indicator) */
body {
  padding-top: env(safe-area-inset-top);
  padding-bottom: env(safe-area-inset-bottom);
  padding-left: env(safe-area-inset-left);
  padding-right: env(safe-area-inset-right);
}
```

### Определение платформы в Go

Используйте теги сборки, чтобы добавить поведение для конкретной платформы, не загромождая общий код.

Создайте `mobile_ios.go` для кода, предназначенного только для iOS:

```go {title="mobile_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.IOSOptions {
    return application.IOSOptions{
        DisableBounce: true,
    }
}
```

Создайте `mobile_android.go` для кода, предназначенного только для Android:

```go {title="mobile_android.go"}
//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.AndroidOptions {
    return application.AndroidOptions{}
}
```

Создайте `mobile_desktop.go` как заглушку, чтобы общий код компилировался и для настольных платформ:

```go {title="mobile_desktop.go"}
//go:build !ios && !android

package main

type mobileOptions struct{}

func platformOptions() mobileOptions { return mobileOptions{} }
```

### Определение платформы в JavaScript и отображение интерфейса только для мобильных платформ

Объекты среды выполнения `IOS.*` и `Android.*` существуют только на соответствующих платформах. Их вызов на настольной платформе приводит к исключению. Правильный подход, используемый в Kitchen Sink, — один раз определить платформу и полностью скрыть элементы управления, предназначенные только для мобильных платформ:

```javascript
// Detect platform from the bridge the host injects into the WebView
const platform = (() => {
  if (typeof window.wails?.platform === 'function') return window.wails.platform(); // Android
  if (window.webkit?.messageHandlers?.external) return 'ios';
  return 'desktop';
})();

const isIOS     = platform === 'ios';
const isAndroid = platform === 'android';
const isMobile  = isIOS || isAndroid;

// Hide any element marked as mobile-only
document.querySelectorAll('.mobile-only').forEach(el => {
  el.style.display = isMobile ? '' : 'none';
});
```

Затем добавьте в HTML:

```html
<section class="mobile-only">
  <button id="btnHaptic">Haptic feedback</button>
</section>
```

При таком подходе кнопки, предназначенные только для мобильных платформ, никогда не отрисовываются на настольных платформах, и вам не приходится защищать каждый отдельный вызов проверкой `if (isMobile)`.

На стороне Go дополните это заглушкой с тегом сборки, чтобы обработчики событий регистрировались только на тех платформах, где они нужны:

```go {title="native_desktop.go"}
//go:build !ios && !android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

// No-op on desktop — mobile tabs are hidden in the frontend so these
// events are never emitted.
func registerNativeFeatures(app *application.App) {}
```

```go {title="native_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func registerNativeFeatures(app *application.App) {
    app.Event.On("common:haptic", func(e *application.CustomEvent) {
        // only compiled and called on iOS
        application.IOS.Haptic("medium")
    })
    // ... other handlers
}
```

Именно этот подход используется в примере [Kitchen Sink](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile) — см. `native_features_stub.go`, `native_features_ios.go` и `native_features_android.go`.

@note{type="note" title="API нативных функций и именование событий"}
Полезно знать о двух соглашениях:

- **Нативные функции на стороне Go используют менеджеры платформ.** Вызывайте их через синглтоны `application.IOS.*` и `application.Android.*` — например, `application.IOS.Haptic("medium")` или `application.Android.Share(payload)`. Каждый менеджер существует только на своей платформе, поэтому его вызовы находятся в файлах `//go:build ios` / `//go:build android`.
- **Пространства имён событий зависят от охвата.** Для всего, что поддерживают обе платформы, используется префикс `common:*` (`common:haptic`, `common:location`, …); для событий, которые может создавать или обрабатывать только одна платформа, используется `ios:*` или `android:*` (например, `ios:backgroundTask`, `android:foregroundService`). Поскольку почти все мобильные функции являются общими, во фронтенде достаточно одного обработчика для каждого события в пространстве имён `common:*`.

@end

### Добавление тактильной обратной связи (iOS)

```javascript
import { IOS } from '@wailsio/runtime';

async function onButtonTap() {
  if (isIOS) {
    await IOS.Haptics.Impact({ style: 'medium' });
  }
  // ... rest of your handler
}
```

### Добавление вибрации (Android)

```javascript
import { Android } from '@wailsio/runtime';

async function onButtonTap() {
  if (isAndroid) {
    await Android.Haptics.Vibrate(50); // 50ms
  }
}
```

---

## Сборка для промышленной эксплуатации

@tabs{sync-key="mobile-platform"}
[iOS]
**Сборка для симулятора** (для тестирования в симуляторе; подпись не требуется):

```bash
wails3 task ios:package
wails3 task ios:deploy-simulator
```

**Сборка для устройства** (требуются идентификатор подписи и профиль подготовки):

```bash
wails3 task ios:package \
  IOS_PLATFORM=device \
  CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
  PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device   # installs via xcrun devicectl
```

**Дистрибутив IPA** (для App Store или TestFlight):

```bash
wails3 task ios:package:ipa IOS_PLATFORM=device \
  CODESIGN_IDENTITY="..." \
  PROVISIONING_PROFILE=path/to/distribution.mobileprovision
```

@note{type="tip"}
Для загрузки в App Store Connect используйте `wails3 task ios:xcode` и позвольте Xcode управлять подписью и архивацией — он автоматически выполнит сложную работу с сертификатами, профилями и нотариальным заверением.

@end

[Android]
**Отладочный APK** (подписан отладочным хранилищем ключей Android, устанавливается напрямую):

```bash
wails3 task android:package
wails3 task android:deploy-emulator
```

**Релизный APK** (подписан вашим хранилищем ключей):

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=yourpassword \
ANDROID_KEY_ALIAS=youralias \
ANDROID_KEY_PASSWORD=yourkeypassword \
  wails3 task android:package
```

**Универсальный APK** (arm64 и x86_64 в одном файле):

```bash
wails3 task android:package:fat
```

@note{type="tip"}
Для загрузки в Play Store создайте `.aab` (Android App Bundle) вместо APK: откройте `build/android/` в Android Studio и выберите **Build → Generate Signed Bundle / APK**.

@end

@end

---

## Устранение неполадок

### Ошибка «no iOS SDKs found» при выполнении `wails3 task ios:run`

Необходимо установить полную версию Xcode и выбрать её:

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
xcode-select -p  # should print the Xcode path
```

#### Ошибка «SDK not found» при выполнении `wails3 task android:run`

Убедитесь, что переменная `ANDROID_HOME` задана и экспортирована. Проверьте это с помощью:

```bash
echo $ANDROID_HOME
ls $ANDROID_HOME/platform-tools/adb
```

#### Симулятор не загружается

Выведите список доступных симуляторов и вручную запустите один из них:

```bash
xcrun simctl list devices available
xcrun simctl boot "iPhone 16"
```

#### В `chrome://inspect` не отображаются объекты для отладки

WebView должен работать в режиме отладки (этот режим по умолчанию используется для `android:run`). Убедитесь, что запущена отладочная, а не промышленная сборка. Также проверьте, что `adb devices` показывает эмулятор как подключённый.

#### Отступы безопасной области не применяются

Убедитесь, что HTML содержит метатег viewport:

```html
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
```

---

## Изучите Kitchen Sink

Когда ваше первое приложение заработает, пример **Kitchen Sink** позволит быстрее всего узнать о других возможностях. Это полноценное приложение Wails, которое запускается на iOS, Android и настольных платформах из единой кодовой базы и демонстрирует тактильную обратную связь, геолокацию, биометрию, локальные уведомления, защищённое хранилище и многое другое:

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

Просмотрите исходный код по адресу [`v3/examples/mobile`](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile) — файлы `native_features_ios.go` и `native_features_android.go` особенно удобны как заготовки для копирования и вставки при реализации функций, специфичных для отдельных платформ.

## Что дальше

@cards{cols="2"}
Руководство по iOS
Полный справочник: параметры конфигурации, нативные вкладки, переключатели WKWebView, сборки для устройств и подписание.

[Руководство по iOS →](/guides/mobile/ios/)

---
Руководство по Android
Полный справочник: конфигурация, всплывающие уведомления, создание пакетов для Play Store и сведения о NDK.

[Руководство по Android →](/guides/mobile/android/)

---
📖 Исходный код Kitchen Sink
Тактильная обратная связь, геолокация, биометрия, уведомления и защищённое хранилище — всё в одном готовом к запуску приложении.

[Посмотреть на GitHub →](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile)

@end
