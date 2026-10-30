---
title: "Android"
description: "Сборка и запуск приложений Wails на Android: настройка инструментов, эмулятор, подпись APK, упаковка для Play Store и справочник по API"
slug: "guides/mobile/android"
sourcePath: "guides/mobile/android.md"
---

@note{type="caution" title="Экспериментальная функция"}
Поддержка Android является экспериментальной и может измениться в будущих выпусках.

@end

@note{type="tip"}
Впервые разрабатываете мобильные приложения с Wails? Начните с пошагового руководства [«Ваше первое мобильное приложение» →](/guides/mobile/first-mobile-app/), а затем вернитесь сюда за полной справочной информацией.

@end

Приложения Wails v3 работают на Android как нативные приложения: Android `WebView` отображает фронтенд, ресурсы обслуживаются **внутри процесса** через `WebViewAssetLoader` на базе сервера ресурсов Go (без локального сервера и открытых портов), а стандартный `@wailsio/runtime` работает без изменений — привязки сервисов, события, диалоговые окна и буфер обмена взаимодействуют через обработчик сообщений Go.

Один и тот же `main.go` собирается для настольных систем и Android. Код Go компилируется в динамическую библиотеку C (`libwails.so`, `GOOS=android` + набор инструментов NDK) и загружается небольшим хостом на Java. Специфичное для Android поведение реализовано в отдельных файлах Go для этой платформы, защищённых условием `//go:build android`.

## Требования

- **Android SDK** с platform-tools, платформой SDK (API 35), build-tools и **NDK** (26.3.x) — `wails3 doctor` показывает обнаруженные компоненты
- **JDK** (например, OpenJDK 21) для Gradle; задайте `JAVA_HOME`, если `java` отсутствует в `PATH`
- Go 1.26+ и npm
- `ANDROID_HOME` (или `ANDROID_SDK_ROOT`), указывающая на SDK

Установите компоненты SDK с помощью инструментов командной строки:

```bash
sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" \
           "ndk;26.3.11579264" "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails \
           --package "system-images;android-35;google_apis;arm64-v8a" \
           --device pixel_7
```

## Запуск в эмуляторе

Из каталога проекта:

```bash
wails3 task android:run
```

Эта команда запускает эмулятор, если он ещё не запущен, генерирует привязки, собирает фронтенд, компилирует код Go в `libwails.so` для ABI эмулятора, создаёт отладочный APK с помощью Gradle, а затем устанавливает и запускает его.

Полезные сопутствующие команды:

```bash
wails3 task android:logs    # stream the app's logcat output
```

В отладочных сборках WebView можно инспектировать из Chrome по адресу `chrome://inspect`.

## Упаковка

```bash
wails3 task android:package             # production release APK
wails3 task android:deploy-emulator     # install + launch it
wails3 task android:bundle              # production release AAB (Android App Bundle)
wails3 task android:bundle:fat          # release AAB containing all ABIs
wails3 task android:run:device          # debug install + launch on a physical device
wails3 task android:deploy-device       # install + launch on a physical device
DEVICE_ID=<serial> wails3 task android:run:device
DEVICE_ID=<serial> wails3 task android:deploy-device
```

Рабочие сборки используют `-tags production,android`, очищаются от отладочной информации и компилируются без внутренней диагностики фреймворка. `wails3 task android:package:fat` включает `arm64-v8a` и `x86_64` в один APK.

Google Play требует формат Android App Bundle (`.aab`) для публикации новых приложений, причём новые приложения должны быть ориентированы на Android 15 (API 35) или более позднюю версию; шаблон проекта задаёт для `compileSdk` и `targetSdk` значение 35 в `build/android/app/build.gradle`. Команда `wails3 task android:bundle:fat` создаёт `bin/<AppName>.aab`, включающий обе ABI; Google Play генерирует из него оптимизированные APK для каждого устройства, поэтому универсальный пакет является подходящим артефактом для загрузки в магазин. Для локального тестирования и тестирования в эмуляторе быстрее всего использовать APK, поскольку `.aab` нельзя установить напрямую с помощью `adb`.

Задачи `android:run` и `android:deploy-emulator` предназначены для эмулятора. Для физического устройства Android используйте `android:run:device`, чтобы получить отладочный APK, или `android:deploy-device`, чтобы получить релизный APK. Обе задачи выполняют сборку для `arm64`, выбирают первое подключённое физическое устройство из вывода `adb devices`, устанавливают приложение и запускают `com.wails.app.MainActivity`. Чтобы выбрать конкретное устройство, передайте `DEVICE_ID=<serial>`.

## Подпись и релизные сборки

Если хранилище ключей не указано, релизные сборки подписываются **отладочным** хранилищем ключей Android, чтобы их можно было установить для тестирования. Чтобы подписать сборку собственным хранилищем ключей, задайте:

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=... \
ANDROID_KEY_ALIAS=... \
ANDROID_KEY_PASSWORD=... \
  wails3 task android:package
```

Те же переменные используются для подписи App Bundle: задайте их и запустите `wails3 task android:bundle:fat`, чтобы создать готовый для Play пакет `.aab`. Без этих переменных пакет будет подписан отладочным хранилищем ключей, и Google Play отклонит его, поэтому задача выводит предупреждение.

@note{type="tip"}
При использовании [Play App Signing](https://support.google.com/googleplay/android-developer/answer/9842756) хранилище ключей, которым вы подписываете приложение локально, содержит ваш **ключ загрузки**: Google проверяет им загруженный пакет, а затем повторно подписывает приложение управляемым ключом подписи приложения. Учтите также, что для каждой загрузки Google Play требует более высокое значение `versionCode`; увеличивайте его в `build/android/app/build.gradle`.

@end

## Конфигурация

Во время выполнения фронтенд управляет функциями Android через объект среды выполнения `Android`: `Android.Haptics.Vibrate(durationMs)`, `Android.Device.Info()`, `Android.Toast.Show(message)`. Имя пакета задаётся параметром `APP_ID` в задачах сборки.

## Что работает, а что нет

| Область | Состояние |
| --- | --- |
| WebView и ресурсы внутри процесса (`WebViewAssetLoader`) | ✅ |
| Привязки сервисов и события (в обоих направлениях) | ✅ |
| Диалоговые окна сообщений | ✅ AlertDialog с обратными вызовами кнопок |
| Диалоговые окна открытия одного или нескольких файлов | ✅ Storage Access Framework (файлы импортируются как копии в кэше) |
| Диалоговые окна выбора каталога и сохранения файла | ❌ Возвращают ошибку — вместо этого записывайте данные в песочнице приложения |
| Буфер обмена | ✅ ClipboardManager |
| API экранов | ✅ WindowMetrics, включая рабочую область без системных панелей |
| События жизненного цикла (`events.Android.*`) | ✅ |
| Тактильная обратная связь, сведения об устройстве и всплывающие уведомления | ✅ API среды выполнения `Android.*` |
| Сборки для эмулятора и физического устройства | ✅ `android:run`, `android:run:device`, `android:deploy-emulator`, `android:deploy-device` |
| Геометрия окна, меню и системный трей | Преднамеренно не выполняют никаких действий |
| Несколько окон | Отображается только первое окно |

## Примечания по переносу

- Код для настольных систем компилируется без изменений при `GOOS=android`; вызовы для управления геометрией, меню и системным треем не выполняют никаких действий, поскольку приложения Android работают в полноэкранном режиме.
- `android` **подразумевает тег сборки `linux`** (Android использует ядро Linux): для файлов, предназначенных только для настольного Linux, требуется `//go:build linux && !android`, а во время выполнения `runtime.GOOS` имеет значение `"android"`.
- Вместо диалоговых окон сохранения файла и выбора каталога записывайте данные в песочницу приложения, а затем передавайте их через системный механизм общего доступа. Диалоговые окна открытия файлов работают и импортируют выбранные документы как копии в каталоге кэша, поэтому вы получаете реальные пути в файловой системе.
- Настоящее приложение всегда собирается с `CGO_ENABLED=1` и NDK; путь без cgo существует только для того, чтобы такие инструменты, как `wails3 generate bindings`, могли загрузить пакет.
- Создавайте адаптивный фронтенд; рабочая область `Screens` не включает строку состояния и панель навигации.
