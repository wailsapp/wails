---
title: "Журнал изменений"
description: "История версий и примечания к выпускам Wails v3"
slug: "changelog"
sourcePath: "changelog.md"
---

Условные обозначения:

-  — macOS
- ⊞ — Windows
- 🐧 — Linux

/_-- Все значимые изменения в этом проекте будут задокументированы в этом файле.

Формат основан на [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), а проект следует принципам [семантического версионирования](https://semver.org/spec/v2.0.0.html).

- `Added` — новые возможности.
- `Changed` — изменения существующей функциональности.
- `Deprecated` — возможности, которые вскоре будут удалены.
- `Removed` — удалённые возможности.
- `Fixed` — исправления ошибок.
- `Security` — устранение уязвимостей.

_/

/_   * ПОЖАЛУЙСТА, НЕ ОБНОВЛЯЙТЕ ЭТОТ ФАЙЛ *   Добавляйте обновления в `v3/UNRELEASED_CHANGELOG.md`   Спасибо! _/

## [Не выпущено]

## v3.0.0-beta.21 — 2026-09-13

## Добавлено

- Документация Wails v3 теперь предоставляется через M-Press в [PR](https://github.com/wailsapp/wails/pull/6116) от @leaanthony

## Исправлено

- Добавлен разбор JSON-значений slug из метаданных MPD для формирования журнала изменений в [PR](https://github.com/wailsapp/wails/pull/6118) от @leaanthony
- Средство обновления очищает вспомогательные переменные окружения и повторно запускает исходный целевой файл после сбоев резервного копирования в [PR](https://github.com/wailsapp/wails/pull/6080) от @cnmax
- Обработчик сигналов по умолчанию запускается во время App.Run в [PR](https://github.com/wailsapp/wails/pull/6098) от @leaanthony
- Меню Windows корректно обрабатывает меню со значением nil, освобождает заменённые ресурсы и перерисовывает строку меню в [PR](https://github.com/wailsapp/wails/pull/6112) от @taliesin-ai
- Восстановлена упаковка MSIX для новых проектов, использующих общую конфигурацию YAML, в [PR](https://github.com/wailsapp/wails/pull/6115) от @leaanthony
- Исправлена ошибка загрузки сгенерированных привязок JavaScript и TypeScript, возникавшая, когда функции создания обобщённых моделей ссылались на вспомогательные объявления, расположенные ниже; также предотвращено переполнение стека при создании взаимозависимых обобщённых моделей (#6062)

## v3.0.0-beta.20 — 2026-09-10

## Изменено

- Ссылки на демонстрационный проект Clave обновлены и теперь ведут на текущий сайт и репозиторий — [PR](https://github.com/wailsapp/wails/pull/6082) от @01xR4in

## Исправлено

- Отменяются прерванные запросы ресурсов в Windows, включая запросы Worker, при этом обработчики keepalive сохраняются между переходами. На платформах Apple контексты нативных запросов передаются через обёртку приложения. (#5963, #5969)
- Записи журнала изменений сохраняются при конкурирующих отправках изменений благодаря повторным попыткам — [PR](https://github.com/wailsapp/wails/pull/6094) от @leaanthony
- Исправлен сбой `go mod vendor` с ошибкой `pattern arm64/WebView2Loader.dll: no matching files found` на всех платформах: удалены встраивания, ссылавшиеся на двоичные файлы, которые никогда не поставлялись в составе модуля. Это исправляет [#5782](https://github.com/wailsapp/wails/issues/5782) и [#5376](https://github.com/wailsapp/wails/issues/5376) — [PR](https://github.com/wailsapp/wails/pull/6031) от @Grantmartin2002

## Удалено

- Удалена поддержка нативного загрузчика WebView2, заменённого загрузчиком на чистом Go. Удалены встроенные двоичные файлы `WebView2Loader.dll` и зависимость `github.com/jchv/go-winloader`. Тег сборки `native_webview2loader` по-прежнему принимается и больше не вызывает ошибку, но не влияет на сборки v3 — [PR](https://github.com/wailsapp/wails/pull/6031) от @Grantmartin2002
- Из руководства по API macOS удалены неиспользуемые теги сборки и параметр FPS — [PR](https://github.com/wailsapp/wails/pull/6097) от @leaanthony

## v3.0.0-beta.19 — 2026-09-09

## Добавлено

- Доступ к закрытым API macOS теперь предоставляется только при явном включении с помощью тегов сборки — см. [документацию](https://v3.wails.io/features/browser/integration), [документацию](https://v3.wails.io/features/environment/info), [документацию](https://v3.wails.io/features/windows/basics), [документацию](https://v3.wails.io/features/windows/frameless), [документацию](https://v3.wails.io/features/windows/notch-windows), [документацию](https://v3.wails.io/features/windows/options), [документацию](https://v3.wails.io/guides/build/macos), [документацию](https://v3.wails.io/guides/build/private-macos-apis) и [документацию](https://v3.wails.io/reference/overview) — [PR](https://github.com/wailsapp/wails/pull/6087) от @leaanthony

## Исправлено

- Запросы среды выполнения размером более 64 МиБ отклоняются с кодом HTTP 413 — [PR](https://github.com/wailsapp/wails/pull/6091) от @leaanthony

## Безопасность

- Усилена защита источников MCP и удалённого доступа с помощью аутентификации по токену — [PR](https://github.com/wailsapp/wails/pull/6092) от @leaanthony

## v3.0.0-beta.18 — 2026-09-08

## Исправлено

- Устранена утечка памяти Calloc в Linux и Darwin за счёт использования получателей-указателей — [PR](https://github.com/wailsapp/wails/pull/6083) от @4RH1T3CT0R7

## v3.0.0-beta.17 — 2026-09-06

## Исправлено

- Windows: ошибка или результат nil при вызове `GetRequest` в обработчике WebResourceRequested больше не приводят к завершению процесса (через `log.Fatal` или панику при разыменовании nil). Вместо этого запрос отбрасывается, а событие записывается в журнал — [PR](https://github.com/wailsapp/wails/pull/6006) от @midagedev

## v3.0.0-beta.16 — 2026-08-29

## Изменено

- Пароль для нотаризации запрашивается в новом окне терминала — [PR](https://github.com/wailsapp/wails/pull/6029) от @leaanthony

## Исправлено

- Типы щелчков по значку в области уведомлений macOS теперь обрабатываются правильно — [PR](https://github.com/wailsapp/wails/pull/5919) от @ChewbaccaCookie
- Перед обновлением CI удаляет неиспользуемые репозитории apt от Microsoft — [PR](https://github.com/wailsapp/wails/pull/6041) от @Grantmartin2002

## v3.0.0-beta.15 — 2026-08-27

## Исправлено

- Время ожидания встраивания WebView2 увеличено до 60 секунд — [PR](https://github.com/wailsapp/wails/pull/6043) от @Grantmartin2002

## v3.0.0-beta.14 — 2026-08-26

## Исправлено

- Нажатия сочетаний Control с буквенными клавишами в macOS теперь именуются правильно — [PR](https://github.com/wailsapp/wails/pull/6032) от @taliesin-ai
- Исправлены значки формата ICO в области уведомлений и добавлено следование теме панели задач Windows — [PR](https://github.com/wailsapp/wails/pull/6016) от @nik9play

## v3.0.0-beta.13 — 2026-08-25

## Исправлено

- В macOS обработка задач главного потока продолжается во время работы модального цикла — [PR](https://github.com/wailsapp/wails/pull/6026) от @leaanthony
- Операции с защищённым хранилищем на мобильных платформах теперь могут завершаться ошибкой и безопасно блокируют доступ при сбое — [PR](https://github.com/wailsapp/wails/pull/5923) от @mortenolsrud
- Обработчики событий приложения запускаются, даже если не зарегистрировано ни одного слушателя — [PR](https://github.com/wailsapp/wails/pull/5999) от @archy-rock3t-cloud
- Исправлены опечатки в комментариях и локализованной документации — [PR](https://github.com/wailsapp/wails/pull/6023) от @haoku123
- Удалены предварительно скомпилированные двоичные файлы macOS, добавленные в репозиторий в каталоге `v3/examples` — [PR](https://github.com/wailsapp/wails/pull/6025) от @4RH1T3CT0R7

## v3.0.0-beta.12 — 2026-08-21

## Добавлено

- Добавлены окна уведомлений для области выреза экрана macOS, а также пример управления их жизненным циклом и телеметрией — см. [документацию](https://v3.wails.io/features/windows/notch-windows) в [PR](https://github.com/wailsapp/wails/pull/6010) от @leaanthony
- Добавлена поддержка окон NSPanel в macOS с новыми параметрами и нативной интеграцией — см. [документацию](https://v3.wails.io/features/windows/options) в [PR](https://github.com/wailsapp/wails/pull/6008) от @leaanthony

## Исправлено

- Предотвращено зависание SQLite Prepare при параллельных вызовах в [PR](https://github.com/wailsapp/wails/pull/5998) от @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Удалено

- Из документации удалён устаревший трекер реализации в [PR](https://github.com/wailsapp/wails/pull/6005) от @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Исправлено

- Исправлена потеря аргументов запуска для пользовательских протоколов и ассоциаций файлов в хосте GTK4 для Linux в [PR](https://github.com/wailsapp/wails/pull/6000) от @midagedev
- При проверке журнала изменений теперь корректно обрабатываются удалённые строки и исправления из того же источника в [PR](https://github.com/wailsapp/wails/pull/5993) от @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Добавлено

- Добавлен защищённый MCP-сервер wails3 для управления проектами с помощью агентов в [PR](https://github.com/wailsapp/wails/pull/5896) от @leaanthony
- Добавлена документация по моделям в привязках — см. [документацию](https://v3.wails.io/features/bindings/models) в [PR](https://github.com/wailsapp/wails/pull/5988) от @taliesin-ai
- Добавлена поддержка установки через rpm-ostree в атомарных системах Linux в [PR](https://github.com/wailsapp/wails/pull/5987) от @leaanthony
- Добавлены встроенное еженедельное создание и публикация графика истории звёзд — см. [документацию](https://v3.wails.io/credits), [документацию](https://v3.wails.io/de/credits), [документацию](https://v3.wails.io/fr/credits), [документацию](https://v3.wails.io/id/credits), [документацию](https://v3.wails.io/ja/credits), [документацию](https://v3.wails.io/ko/credits), [документацию](https://v3.wails.io/pt/credits), [документацию](https://v3.wails.io/ru/credits), [документацию](https://v3.wails.io/zh-cn/credits) и [документацию](https://v3.wails.io/zh-tw/credits) в [PR](https://github.com/wailsapp/wails/pull/5986) от @leaanthony
- Добавлен пакет mac только для Darwin, предназначенный для поиска ресурсов пакета приложения, — см. [документацию](https://v3.wails.io/guides/build/macos) в [PR](https://github.com/wailsapp/wails/pull/5965) от @leaanthony
- Добавлены демонстрационная страница Condui и запись в указателе — см. [документацию](https://v3.wails.io/community/showcase/condui) и [документацию](https://v3.wails.io/community/showcase) в [PR](https://github.com/wailsapp/wails/pull/5962) от @mgueregath
- Добавлена демонстрационная страница Redis Viewer со снимками экрана и ссылкой на проект — см. [документацию](https://v3.wails.io/community/showcase) и [документацию](https://v3.wails.io/community/showcase/redisviewer) в [PR](https://github.com/wailsapp/wails/pull/5984) от @redisviewer

## Изменено

- Флаги приложения GTK для Linux изменены на G<em>APPLICATION</em>NON_UNIQUE в [PR](https://github.com/wailsapp/wails/pull/5971) от @overlordtm
- Отсутствующие события окон теперь регистрируются на уровне отладки, а не предупреждения, в [PR](https://github.com/wailsapp/wails/pull/5914) от @julianstorer

## Исправлено

- Зарегистрированным сочетаниям клавиш macOS разрешено иметь приоритет над webview в [PR](https://github.com/wailsapp/wails/pull/5902) от @julianstorer
- Исправлены неработающие ссылки на боковой панели документации в [PR](https://github.com/wailsapp/wails/pull/5937) от @northes
- Контексты запросов ресурсов в macOS и iOS теперь отменяются, когда WebKit прерывает соответствующую задачу пользовательской схемы (#5963)
- Добавлена обработка сообщения WindowSetFullscreenButtonEnabled в [PR](https://github.com/wailsapp/wails/pull/5976) от @archy-rock3t-cloud
- В шаблон preact-ts импортирован Fragment для устранения ошибки сборки в [PR](https://github.com/wailsapp/wails/pull/5979) от @haoku123
- Предотвращён сбой устаревших приложений GTK3, выполняющих только сервисы, когда обнаружение экранов выполняется до появления активного окна или доступного дисплея (#5966)
- Запуски выпуска с явно указанной версией теперь могут продолжаться, даже если раздел невыпущенных изменений в журнале пуст (#5977)

## Безопасность

- Файлы блокировок nanoid для веб-сайта обновлены до исправленной версии 3.3.18 для устранения уведомлений о проблемах безопасности в [PR](https://github.com/wailsapp/wails/pull/5985) от @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Добавлено

- В автоматические записи журнала изменений добавлено создание URL-адресов документации в [PR](https://github.com/wailsapp/wails/pull/5957) от @taliesin-ai
- Добавлены Streams: двунаправленные потоки байтов между Go и JavaScript с программной моделью WebSocket и без прослушивающего сокета. Объявите поток в Go с помощью `app.HandleStream(name, handler)` и подключитесь к нему из фронтенда с помощью `Stream(name)`, который возвращает объект структуры `WebSocket`. Данные Go→JS передаются через сервер ресурсов посредством одного удерживаемого опроса на окно, а JS→Go — обычным запросом POST; при этом порт TCP не привязывается, а данные не проходят через `evaluateJavaScript`. В серверных сборках (`-tags server`) тот же обработчик вместо этого обслуживается через настоящий WebSocket, поэтому код приложения одинаков во всех сборках. Автор: @leaanthony
- Запись журнала изменений о почтовом ящике перенесена в раздел невыпущенных изменений в [PR](https://github.com/wailsapp/wails/pull/5935) от @leaanthony

## Изменено

- Обновлены автоматическое создание боковой панели документации и определение типа автора блога в [PR](https://github.com/wailsapp/wails/pull/5938) от @leaanthony

## Исправлено

- При инициализации WebView2 теперь используются крайний срок и цикл обработки сообщений в [PR](https://github.com/wailsapp/wails/pull/5952) от @leaanthony
- Тест файлов cookie WebView2 пропускается в CI, если он не включён явно, а его выполнение закрепляется за текущим потоком ОС в [PR](https://github.com/wailsapp/wails/pull/5951) от @leaanthony
- Построители меню Windows восстанавливают идентификаторы команд для родительских элементов подменю в [PR](https://github.com/wailsapp/wails/pull/5944) от @gilad-ch
- Официальный образ для кросс-компиляции приведён в соответствие с базовым уровнем поддержки GTK 4.14+ в Linux (#5928)
- Проект Xcode для iOS настроен так, чтобы сохранять унаследованные флаги компоновщика и добавлять -ObjC, в [PR](https://github.com/wailsapp/wails/pull/5915) от @mortenolsrud
- Устранено чрезмерное создание и закрытие TCP-соединений в прокси ресурсов `wails3 dev` для крупных фронтендов, которое могло исчерпать эфемерные порты хоста и привести к сбою несвязанных процессов с ошибкой `EADDRNOTAVAIL`
- JavaScript событий для каждого окна теперь помещается в очередь для упорядоченной доставки и управления обратным давлением в [PR](https://github.com/wailsapp/wails/pull/5934) от @leaanthony

## Удалено

- Удалён конвейер выпуска исполняемых файлов для настольных систем: выпуски v3 теперь представлены только тегами, а CLI `wails3` устанавливается с помощью `go install`. Удалены `release-v3.yml` и запускавший его этап ночной сборки в [PR](https://github.com/wailsapp/wails/pull/5946) от @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Добавлено

- Добавлена настройка автовоспроизведения в macOS, позволяющая отключить требование пользовательского действия для воспроизведения мультимедиа, в [PR](https://github.com/wailsapp/wails/pull/5512) от @Eyalm321
- Запись об изменении почтового ящика перенесена в раздел «Невыпущенные изменения» в [PR](https://github.com/wailsapp/wails/pull/5935) от @leaanthony

## Изменено

- Для более плавной анимации масштабирования в macOS теперь используются CADisplayLink или NSTimer в [PR](https://github.com/wailsapp/wails/pull/5945) от @savely-krasovsky

## Исправлено

- Проект Xcode для iOS настроен так, чтобы сохранять унаследованные флаги компоновщика и добавлять -ObjC, в [PR](https://github.com/wailsapp/wails/pull/5915) от @mortenolsrud
- Устранено чрезмерное создание и закрытие TCP-соединений в прокси ресурсов `wails3 dev` для крупных фронтендов, которое могло исчерпать эфемерные порты хоста и привести к сбоям несвязанных процессов с ошибкой `EADDRNOTAVAIL`
- JavaScript-код событий каждого окна теперь ставится в очередь для упорядоченной отправки и обеспечения обратного давления в [PR](https://github.com/wailsapp/wails/pull/5934) от @leaanthony

### Добавлено

- Реализован универсальный асинхронный почтовый ящик FIFO для упорядоченной доставки событий в [PR](https://github.com/wailsapp/wails/pull/5851) от @savely-krasovsky и @DevLumuz

## v3.0.0-beta.6 — 2026-08-09

## Добавлено

- Реализовано ограниченное по объёму хранилище на стороне хоста для слишком крупных событий и их упорядоченная доставка в JavaScript в [PR](https://github.com/wailsapp/wails/pull/5930) от @leaanthony
- Реализовано подпрыгивание значка в Dock macOS при привлечении внимания к окну в [PR](https://github.com/wailsapp/wails/pull/5921) от @julianstorer

## Исправлено

- При сбросе буфера сервер ресурсов теперь сохраняет ошибки анализатора типа содержимого и незаписанные префиксы в [PR](https://github.com/wailsapp/wails/pull/5931) от @leaanthony
- Предотвращено аварийное завершение приложений macOS при замене меню приложения из обратного вызова Wails
- Исправлены нечитаемые нативные меню в Windows 10 1809 / Windows Server 2019 (сборка 17763). Проверка версии разрешала использование экспортируемых функций uxtheme для тёмного режима только начиная со сборки 18334, поэтому на этих хостах не выполнялось включение тёмного режима на уровне приложения: фон меню отрисовывался тёмным, но Windows продолжала отображать текст меню в светлой теме, из-за чего тёмный текст оказывался на тёмном фоне. Функции с этими порядковыми номерами существуют начиная с 17763, поэтому порог проверки теперь соответствует этой сборке.
- Исправлен вызов `GetDeviceCaps` вместо `GetStockObject` в `w32.GetStockObject`, из-за которого для каждого стандартного объекта возвращалось 0.
- Улучшены обработка ошибок и сообщения об ошибках при скачивании загрузчика WebView2 в [PR](https://github.com/wailsapp/wails/pull/5924) от @jannskiee

## v3.0.0-beta.5 — 2026-08-07

## Исправлено

- Активация приложения в macOS теперь учитывает политику активации только для обычных приложений в [PR](https://github.com/wailsapp/wails/pull/5897) от @julianstorer
- Добавлена защита от обращения к неинициализированным окнам GTK в сборках для Linux в [PR](https://github.com/wailsapp/wails/pull/5898) от @julianstorer
- Для окон WebKit в Linux перед загрузкой URL теперь явно задаётся непрозрачный цвет фона в [PR](https://github.com/wailsapp/wails/pull/5899) от @julianstorer

## v3.0.0-beta.4 — 2026-08-05

## Изменено

- Задачи сборки Android теперь по умолчанию используют arm64, а deploy-emulator выбирает архитектуру хоста, в [PR](https://github.com/wailsapp/wails/pull/5890) от @mortenolsrud

## Исправлено

- Состояние масштабирования окна macOS теперь сохраняется при перетаскивании, а объём анимации уменьшен в [PR](https://github.com/wailsapp/wails/pull/5900) от @leaanthony
- Исправлена сборка Windows в серверном режиме: в ограничение сборки `webview_window_windows_nonclient.go` добавлено `!server`

## v3.0.0-beta.3 — 2026-08-03

## Добавлено

- В подробностях реализации задокументировано завершение бета-проверки этапа 10 в [PR](https://github.com/wailsapp/wails/pull/5881) от @leaanthony

## Исправлено

- В API тёмного режима Windows теперь передаётся дескриптор окна и проверяются аргументы в [PR](https://github.com/wailsapp/wails/pull/5877) от @leaanthony
- Централизовано определение состояния кнопок строки заголовка для безрамочных окон macOS в [PR](https://github.com/wailsapp/wails/pull/5870) от @taliesin-ai
- Предотвращено появление нечитаемого текста в нативных меню, когда приложение Windows запрашивает тёмный режим, а тема приложений Windows остаётся светлой. Теперь меню использует соответствующий светлый нативный фон, пока Windows не сможет отображать текст меню в тёмном режиме.
- Исправлены нечитаемые нативные меню в Windows 10 1809 / Windows Server 2019 (сборка 17763). Проверка версии разрешала использование экспортируемых функций uxtheme для тёмного режима только начиная со сборки 18334, поэтому на этих хостах не выполнялось включение тёмного режима на уровне приложения: фон меню отрисовывался тёмным, но Windows продолжала отображать текст меню в светлой теме, из-за чего тёмный текст оказывался на тёмном фоне. Функции с этими порядковыми номерами существуют начиная с 17763, поэтому порог проверки теперь соответствует этой сборке.

## v3.0.0-beta.2 — 2026-08-02

## Изменено

- v3 переведена из альфа-версии в бета-версию
- Задокументированы интеллектуальные значения по умолчанию для системного трея и автоматическое скрытие всплывающего меню; добавлены регрессионные тесты выбора обработчика щелчка (#5840).
- Средство обновления через GitHub теперь по умолчанию исключает ресурсы установщика Windows в [PR](https://github.com/wailsapp/wails/pull/5861) от @leaanthony
- Добавлена поддержка скруглённых и прямых углов, а также углов с настраиваемым радиусом для безрамочных окон macOS в [PR](https://github.com/wailsapp/wails/pull/5866) от @leaanthony

## Исправлено

- Теперь сообщаются текущие размеры окон GTK4, а настроенная поверхность отправляет события изменения размера, разворачивания, сворачивания и изменения состояния полноэкранного режима (#5830).
- Исправлено аварийное завершение WebKit в Linux при отправке Blob или FormData в запросах fetch в [PR](https://github.com/wailsapp/wails/pull/5854) от @taliesin-ai
- При отсутствии заголовков Blob/FormData прослойка fetch теперь передаёт undefined в [PR](https://github.com/wailsapp/wails/pull/5865) от @leaanthony

## v3.0.0-alpha2.122 — 2026-08-01

## Добавлено

## Изменено

- Добавлена поддержка скруглённых и прямых углов, а также углов с настраиваемым радиусом для безрамочных окон macOS в [PR](https://github.com/wailsapp/wails/pull/5866) от @leaanthony

## Исправлено

- При отсутствии заголовков Blob/FormData прослойка fetch теперь передаёт undefined в [PR](https://github.com/wailsapp/wails/pull/5865) от @leaanthony

## v3.0.0-alpha2.121 — 2026-07-31

## Добавлено

- Добавлена поддержка упаковки в DMG для macOS с новыми параметрами и задачами сборки в [PR](https://github.com/wailsapp/wails/pull/5857) от @leaanthony

## Изменено

- Средство обновления через GitHub теперь по умолчанию исключает ресурсы установщика Windows в [PR](https://github.com/wailsapp/wails/pull/5861) от @leaanthony

## Исправлено

- Исправлен сбой WebKit в Linux при отправке Blob или FormData в запросах fetch в [PR](https://github.com/wailsapp/wails/pull/5854) от @taliesin-ai

## v3.0.0-alpha2.120 — 2026-07-31

## Добавлено

- Реализованы действия разворачивания или сворачивания окна по двойному щелчку на строке заголовка в macOS в [PR](https://github.com/wailsapp/wails/pull/5853) от @taliesin-ai

## Изменено

- Руководство по сервису QR обновлено для использования NewServiceWithOptions; также добавлены отступы в [PR](https://github.com/wailsapp/wails/pull/5849) от @jeongkyu

## Исправлено

- Обеспечена отзывчивость WKWebView во время масштабирования в macOS в [PR](https://github.com/wailsapp/wails/pull/5856) от @leaanthony
- Исправлены запросы размера окна GTK4 и отправка событий изменения размера, разворачивания, сворачивания и изменения полноэкранного состояния из настроенного `GdkSurface`.

## v3.0.0-alpha2.119 — 2026-07-27

## Исправлено

- Документация на нескольких языках дополнена архитектурными диаграммами в [PR](https://github.com/wailsapp/wails/pull/5833) от @taliesin-ai

## v3.0.0-alpha2.118 — 2026-07-26

## Добавлено

- Добавлены пути по умолчанию для входных и выходных файлов генерации значков в [PR](https://github.com/wailsapp/wails/pull/5825) от @taliesin-ai
- Исходные модули точки входа добавлены в sideEffects файла package.json пакета среды выполнения в [PR](https://github.com/wailsapp/wails/pull/5797) от @savely-krasovsky
- В руководство для участников добавлен раздел о лицензии и происхождении кода в [PR](https://github.com/wailsapp/wails/pull/5816) от @taliesin-ai

## Исправлено

- Применены ограниченные по области действия стили CSS для безрамочных окон GTK4, устраняющие скругление углов, в [PR](https://github.com/wailsapp/wails/pull/5800) от @savely-krasovsky
- Добавлена корректная обработка ошибок определения положения курсора в Windows для всплывающих меню и перечисления экранов в [PR](https://github.com/wailsapp/wails/pull/5789) от @wayneforrest
- Диалог открытия файла в macOS теперь правильно фильтрует расширения и проверяет допустимость файлов по суффиксу в [PR](https://github.com/wailsapp/wails/pull/5678) от @phergul
- Инициализация тёмного режима в Windows защищена от вызовов nil API в [PR](https://github.com/wailsapp/wails/pull/5793) от @roachadam
- Исправлена ошибка сборки средства обновления для 32-разрядной платформы: при передаче в `fmt.Errorf` на `GOARCH=386` константа `maxArchiveTotalSize` (2 ГиБ) переполняла платформенный тип `int`. Теперь для неё явно указан тип `int64`.
- Исправлена паника из-за nil-указателя при запуске, когда окно использует тёмную строку заголовка (или системную тёмную тему) в сборках Windows, где не загружаются API тёмного режима uxtheme, например в Windows 10 1809 / Windows Server 2019 (сборка 17763). Вызовы `AllowDarkModeForWindow` при настройке темы окна теперь защищены от nil аналогично уже применяемой защите в `w32.SetMenuTheme`.

## v3.0.0-alpha2.117 — 2026-07-08

## Добавлено

- Реализована пользовательская логика hit-test для неклиентских областей в Windows в [PR](https://github.com/wailsapp/wails/pull/5462) от @savely-krasovsky

## Изменено

- Определение масштаба монитора в WebView2 настроено в зависимости от UseVisualHosting в [PR](https://github.com/wailsapp/wails/pull/5761) от @wayneforrest

## v3.0.0-alpha2.116 — 2026-07-07

## Добавлено

- Документация FAQ обновлена с акцентом на возможности и рекомендации Wails v3 в [PR](https://github.com/wailsapp/wails/pull/5763) от @taliesin-ai

## v3.0.0-alpha2.115 — 2026-07-06

## Исправлено

- Исправлена ошибка, из-за которой `Menu.Update()` не перестраивал нативное меню в GTK4 Linux (#5659; независимо диагностирована и исправлена @puneetdixit200 в #5539)
- Исправлен сбой при перечислении экранов macOS после изменения конфигурации дисплеев: строки идентификаторов и имён экранов теперь копируются, а количество фиксируется в снимке (#5565; проблема независимо диагностирована и исправлена @x-haose в #5584)
- Исправлен сбой в Windows, когда `WM_ERASEBKGND` отрисовывает сплошной фон во время перехода между свёрнутым и восстановленным состояниями, при котором `GetClientRect` возвращает nil (о необходимости проверки сообщил @sinspired в #5636)
- Исправлена ошибка, из-за которой ошибка привязок фронтенда всегда разбиралась как текст, — @mbaklor в #5690
- Исправлена ошибка сборки для Windows при использовании тега сборки `server`, вызванная отсутствием в файлах графического интерфейса Windows ограничения сборки `!server`, которое уже присутствует в эквивалентных файлах macOS и Linux (#5680)

## v3.0.0-alpha2.114 — 2026-07-05

## Добавлено

- Реализованы протокол Update Manifest и поставщик конечной точки в [PR](https://github.com/wailsapp/wails/pull/5720) от @taliesin-ai

## Изменено

- Привязка `webview2` включена в модуль v3 как `v3/internal/webview2`; отдельный модуль, его процессы ночных выпусков и синхронизации, а также манипуляции с версией go.mod удалены, поскольку v3 — его единственный потребитель, в [PR](https://github.com/wailsapp/wails/pull/5711) от @taliesin-ai

## Исправлено

- Исправления определения масштаба монитора в WebView2 и повторной синхронизации хоста при изменении DPI перенесены в раздел «Невыпущенные изменения» в [PR](https://github.com/wailsapp/wails/pull/5750) от @taliesin-ai
- Обновлён маршалинг COM в WebView2 для параметров float64 и BOOL в [PR](https://github.com/wailsapp/wails/pull/5741) от @wayneforrest
- Предотвращены паника и разыменование nil при обновлении и уничтожении значка в системном трее Windows в [PR](https://github.com/wailsapp/wails/pull/5703) от @wayneforrest
- Исправлена ошибка, из-за которой скрытые окна в Windows не скрывались повторно должным образом, в [PR](https://github.com/wailsapp/wails/pull/5743) от @wayneforrest
- Видимость контроллера WebView2 синхронизирована со сворачиванием, разворачиванием и восстановлением окна в [PR](https://github.com/wailsapp/wails/pull/5742) от @wayneforrest

### Исправлено

- Повторно включено определение масштаба монитора в WebView2, а повторная синхронизация хоста теперь выполняется только при изменении DPI — [PR](https://github.com/wailsapp/wails/pull/5734) от @taliesin-ai на основе исправления, проверенного @randalmurphal, с подтверждением первопричины от @eleclin и тестированием на оборудовании от @qq540491950

## v3.0.0-alpha2.113 — 2026-07-04

## Добавлено

- Добавлено предупреждение при сборке релизного AAB без заданного `ANDROID_KEYSTORE_FILE` (Google Play отклоняет пакеты с отладочной подписью), а также документация по упаковке и подписыванию App Bundle в [PR](https://github.com/wailsapp/wails/pull/5730) от @taliesin-ai
- Добавлена документация Why Wails на нескольких языках в [PR](https://github.com/wailsapp/wails/pull/5739) от @taliesin-ai
- Добавлена поддержка преобразования Go time.Time в JS Date или строку в привязках в [PR](https://github.com/wailsapp/wails/pull/5398) от @fbbdev
- Добавлены задачи для упаковки Android App Bundle (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) для отправки в Play Store; задачи для APK сохранены для локального тестирования и тестирования в эмуляторе — [PR](https://github.com/wailsapp/wails/pull/5728) от @mortenolsrud (исправляет [#5726](https://github.com/wailsapp/wails/issues/5726))
- Добавлены задачи для физических устройств Android и возобновление запроса разрешений на использование камеры и геолокации в [PR](https://github.com/wailsapp/wails/pull/5735) от @taliesin-ai

## Изменено

- Версия `webview2` повышена до v1.0.28 ([примечания к выпуску](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- В шаблоне Android версия `compileSdk`/`targetSdk` повышена с 34 до 35, как того требует Google Play для публикации новых приложений, в [PR](https://github.com/wailsapp/wails/pull/5730) от @taliesin-ai

## Исправлено

- Исправлены встроенные маски аватаров в sponsorkit в [PR](https://github.com/wailsapp/wails/pull/5745) от @leaanthony
- Исправлен выбор неверного образа системы или версии cmdline-tools при автоматическом создании Android AVD из-за лексикографической сортировки версий в [PR](https://github.com/wailsapp/wails/pull/5730) от @taliesin-ai
- Исправлена устаревшая версия Android NDK, предлагаемая мастером настройки (теперь 26.3.11579264 в соответствии с документированным требованием), в [PR](https://github.com/wailsapp/wails/pull/5730) от @taliesin-ai
- Обновлена французская документация по SvelteKit и параметрам в [PR](https://github.com/wailsapp/wails/pull/5744) от @leaanthony
- Исправлен SIGSEGV при перечислении экранов в macOS во время изменения конфигурации дисплеев в [PR](https://github.com/wailsapp/wails/pull/5516) от @flofreud

## v3.0.0-alpha2.112 — 2026-07-03

## Добавлено

- Добавлен написанный на Go генератор SVG со списком участников проекта и обновлены страницы благодарностей в документации и на сайте в [PR](https://github.com/wailsapp/wails/pull/5724) от @taliesin-ai
- Добавлена поддержка преобразования Go time.Time в JS Date или строку в привязках в [PR](https://github.com/wailsapp/wails/pull/5398) от @fbbdev

## Изменено

- Конвейер создания изображений спонсоров на базе Node заменён генератором на Go в [PR](https://github.com/wailsapp/wails/pull/5719) от @taliesin-ai

## Исправлено

- Исправлен скрипт установки зависимостей ресурсов сборки Android в [PR](https://github.com/wailsapp/wails/pull/5729) от @taliesin-ai
- В `ValidateAndSanitizeURL` запрещён управляющий символ U+0085 (NEXT LINE), что обеспечивает полную проверку пробельных символов валидатором URL
- Добавлен повторный расчёт рамки DWM при изменении DPI для окон без системной рамки в [PR](https://github.com/wailsapp/wails/pull/4785) от @leaanthony
- Исправлена ошибка определения зоны перетаскивания DnD в Windows при масштабировании, отличном от 100%, в [PR](https://github.com/wailsapp/wails/pull/4632) от @yulesxoxo
- Добавлено явное управление памятью Objective-C для объектов Cocoa в диалоговых окнах, меню, системном трее и уведомлениях Darwin в [PR](https://github.com/wailsapp/wails/pull/5714) от @taliesin-ai
- Исправлены ошибки бэкенда Linux CGO и проблемы с системным треем в [PR](https://github.com/wailsapp/wails/pull/5718) от @taliesin-ai

## v3.0.0-alpha2.111 — 2026-07-01

## Добавлено

- HappyTools добавлен в витрину проектов сообщества в [PR](https://github.com/wailsapp/wails/pull/5061) от @Aliuyanfeng
- Добавлены поддержка индонезийской локали и подробная документация в [PR](https://github.com/wailsapp/wails/pull/5643) от @triadmoko
- В WindowsWindow добавлен параметр DisableMenu в [PR](https://github.com/wailsapp/wails/pull/4813) от @leaanthony

## Изменено

- Шаблон Taskfile и CLI обновлены для запуска задач сборки и упаковки с GOOS и ARCH в [PR](https://github.com/wailsapp/wails/pull/5617) от @leaanthony

## Исправлено

- Исправлена проблема с объединением окон macOS во вкладки в [PR](https://github.com/wailsapp/wails/pull/5708) от @taliesin-ai

## Удалено

- Удалены немецкие переводы файлов MDX из разделов об участии в проекте, возможностях и руководствах в [PR](https://github.com/wailsapp/wails/pull/5702) от @taliesin-ai

## v3.0.0-alpha2.110 — 2026-06-30

## Добавлено

- Реализованы обычная и принудительная перезагрузка WebView в macOS, а также восстановление после завершения процесса WebContent в [PR](https://github.com/wailsapp/wails/pull/5129) от @wayneforrest
- Добавлена подробная документация на немецком языке по участию в проекте, возможностям и руководствам в [PR](https://github.com/wailsapp/wails/pull/5396) от @leaanthony
- Уведомления дополнены звуком, вложениями, планированием и API для их обновления в [PR](https://github.com/wailsapp/wails/pull/5333) от @popaprozac

## Исправлено

- Добавлен повторный расчёт рамки DWM при изменении DPI для окон без системной рамки в [PR](https://github.com/wailsapp/wails/pull/4785) от @leaanthony
- Исправлена ошибка определения зоны перетаскивания DnD в Windows при масштабировании, отличном от 100%, в [PR](https://github.com/wailsapp/wails/pull/4632) от @yulesxoxo

## v3.0.0-alpha2.109 — 2026-06-29

## Добавлено

- В документацию EventsEmit добавлены примеры кода в [PR](https://github.com/wailsapp/wails/pull/5026) от @iamhabbeboy
- Добавлен параметр визуального хостинга WebView2 в Windows в [PR](https://github.com/wailsapp/wails/pull/5380) от @MerIijn
- Klustr добавлен в документацию витрины проектов сообщества в [PR](https://github.com/wailsapp/wails/pull/5536) от @SametKUM
- Kira добавлен в витрину проектов сообщества вместе с новыми страницами и записью в журнале изменений в [PR](https://github.com/wailsapp/wails/pull/5685) от @thiennguyen93
- В руководство по сервису MCP добавлен раздел обратной связи в [PR](https://github.com/wailsapp/wails/pull/5694) от @taliesin-ai

## Изменено

- В серверном режиме теперь предусмотрена полноценная производственная сборка, согласованная с задачами сборки настольного приложения (#5693). По умолчанию `task build:server` создаёт производственный исполняемый файл (`-tags server,production`, `-trimpath`, с удалёнными отладочными символами) и принимает параметры `DEV=true` (сервер разработки), `OBFUSCATED=true` (garble) и `EXTRA_TAGS`. `task run:server` запускает сервер разработки. `Dockerfile.server` / `task build:docker` сначала собирают производственный сервер (`-tags server,production`) и производственную версию фронтенда; по умолчанию образ содержит статическую сборку на чистом Go на базе distroless/static, а `CGO_ENABLED`, `GO_IMAGE` и `RUNTIME_IMAGE` доступны как переопределяемые аргументы сборки для приложений с CGO.

## Исправлено

- Предотвращён сбой при закрытии окна с ожидающими асинхронными вызовами в [PR](https://github.com/wailsapp/wails/pull/4435) от @leaanthony
- Предотвращена активация окна при открытии скрытых приложений в Windows в [PR](https://github.com/wailsapp/wails/pull/5249) от @leaanthony
- Обеспечено выполнение обработки метаданных запросов WebKit, завершения ответов и потоков тела в главном потоке GTK в [PR](https://github.com/wailsapp/wails/pull/5668) от @taliesin-ai
- Исправлена ошибка, из-за которой `Menu.Update()` не пересобирал нативное меню в GTK4 для Linux (#5659; независимо диагностировано и исправлено @puneetdixit200 в #5539)
- Исправлен сбой при перечислении экранов macOS после изменения конфигурации дисплеев: строки идентификаторов и имён экранов теперь копируются, а количество фиксируется в снимке состояния (#5565; независимо диагностировано и исправлено @x-haose в #5584)
- Исправлена ошибка, из-за которой содержимое WebView2 сжималось, а затем исчезало после перетаскивания окна между мониторами с разным DPI в Windows: в обработчике `WM_DPICHANGED` теперь повторно задаются границы контроллера, аналогично повторной синхронизации DPI при восстановлении свёрнутого окна (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Добавлено

- Добавлены глобальные (общесистемные) сочетания клавиш через `app.GlobalShortcut` (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Сочетания срабатывают, даже когда приложение не находится в фокусе. Для каждой платформы реализована нативная поддержка без сторонних зависимостей: горячие клавиши Carbon в macOS, `RegisterHotKey` в Windows, `XGrabKey` в X11 и интерфейс глобальных сочетаний клавиш XDG Desktop Portal в Wayland.
- Добавлен встроенный сервер MCP — сервер Model Context Protocol, который запускается автоматически, если приложение собрано с тегом `mcp`. Он позволяет LLM-агентам тестировать работающие приложения Wails и управлять ими: управлять окнами, инспектировать DOM, выполнять JavaScript, вызывать привязанные методы, работать с событиями и имитировать ввод с помощью мыши и клавиатуры, отображая анимированный экранный курсор. Пользовательский код не требуется: тег `mcp` автоматически добавляется командой `wails3 build`/`wails3 dev`, когда задано `WAILS_MCP=1`. Настройка полностью выполняется с помощью переменных окружения (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Исправлено

- Исправлена ошибка, из-за которой `Menu.Update()` не пересобирал нативное меню в GTK4 для Linux (#5659; независимо диагностировано и исправлено @puneetdixit200 в #5539)
- Исправлен сбой при перечислении экранов macOS после изменения конфигурации дисплеев: строки идентификаторов и имён экранов теперь копируются, а количество фиксируется в снимке состояния (#5565; независимо диагностировано и исправлено @x-haose в #5584)

## v3.0.0-alpha2.107 - 2026-06-27

## Добавлено

- Добавлена экспериментальная документация по Wake с навигацией на боковой панели в [PR](https://github.com/wailsapp/wails/pull/5613) от @leaanthony

## v3.0.0-alpha2.106 - 2026-06-24

## Изменено

- Версия `webview2` повышена до v1.0.27.
  - ci(webview2): исправлена сборка выпуска (кросс-компиляция для Windows и полный go.sum) (#5671)\

  **Полный список различий:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- Команда go vet удалена из кросс-компиляции в процессе выпуска webview2 в [PR](https://github.com/wailsapp/wails/pull/5672) от @taliesin-ai
- Модель OpenRouter для автоматического создания журнала изменений обновлена до google/gemini-2.5-flash-lite в [PR](https://github.com/wailsapp/wails/pull/5670) от @taliesin-ai
- Версия `webview2` повышена до v1.0.26.

### Исправления

- **Восстановление после временных ошибок COM во время выполнения вместо завершения работы** (#5658, #5580). Ранее `Chromium.errorCallback` вызывал `os.Exit(1)` при *любой* ошибке COM, поэтому устранимый кратковременный сбой после запуска завершал работу всего приложения. Пути выполнения (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) теперь записывают ошибку в журнал и восстанавливают работу. В частности, некорректное веб-сообщение или веб-сообщение из недоверенного источника в `MessageReceived` теперь отбрасывается, а не завершает процесс. Это устраняет класс сбоев при перемещении окна между мониторами с разным DPI (#5544, #5650). Ошибки на путях создания среды и контроллера по-прежнему являются фатальными.\

**Полный список различий:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Исправлено

- Исправлен процесс release-webview2 для корректной обработки файлов go.sum в [PR](https://github.com/wailsapp/wails/pull/5671) от @taliesin-ai
- Исправлено обновление меню GTK4 в Linux: нативное меню теперь очищается и пересобирается в [PR](https://github.com/wailsapp/wails/pull/5659) от @taliesin-ai

## v3.0.0-alpha2.105 - 2026-06-21

## Добавлено

- Добавлено `application.System` для определения платформы во время выполнения из общего кода: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (тег сборки `server`) и `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` для непосредственной проверки одной целевой платформы. Код компилируется для каждой целевой платформы, поэтому ветвление возможно без тегов сборки. Соответствующие вспомогательные функции фронтенда (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) доступны в `@wailsio/runtime`
- Добавлено руководство «Использование других фронтенд-фреймворков», показывающее, как поместить собственный проект Vite в `frontend/` (рассматриваются Solid, Preact, Lit, SvelteKit, Qwik, Angular и другие)
- Мастер `wails3 setup` теперь проверяет инструментарий для мобильных платформ (iOS/Android): Xcode и среду выполнения iOS Simulator, JDK, Android SDK/NDK и эмулятор. Там, где это применимо, доступны установка одним щелчком и готовые к копированию исправления конфигурации оболочки
- Создаваемые проекты включают `frontend/.npmrc`, где задаётся `minimum-release-age` длительностью 7 дней, чтобы снизить риск, связанный с недавно опубликованными (и потенциально скомпрометированными) пакетами (учитывается pnpm и bun; npm без каких-либо последствий игнорирует эту настройку)

## Изменено

- Все встроенные начальные шаблоны получили новый дизайн главного блока с неоновыми горами (для веба, iOS и Android)
- **Теперь TypeScript используется по умолчанию для начальных шаблонов и занимает имя шаблона без уточняющего суффикса.** `wails3 init` (без `-t`) создаёт каркас проекта TypeScript; `-t vanilla`, `-t react`, `-t vue` и `-t svelte` используют TypeScript, а их варианты на JavaScript доступны как `-t vanilla-js`, `-t react-js`, `-t vue-js` и `-t svelte-js`. Встроенные шаблоны указывают свой язык с помощью `typescript:` в `template.yaml`; шаблоны сообщества с суффиксом `-ts` продолжают работать как резервный вариант
- Мастер `wails3 setup` получил новый дизайн в неоновой теме «цифровой Wails» (эффект матового стекла на фоне гор)

## Исправлено

- Исправлен сбой в Windows при восстановлении приложения, которое было свёрнуто достаточно долго, чтобы WebView2 приостановил работу или перезапустил процесс рендеринга/GPU. Повторная синхронизация DPI при сворачивании и восстановлении (#5544) теперь обращается к контроллеру WebView2 только при фактическом изменении DPI окна, что позволяет избежать фатальных вызовов COM к приостановленному контроллеру в обычном случае восстановления с прежним DPI (#5605)
- Исправлены повторяющиеся нативные сбои `SIGABRT`/`SIGSEGV` (обычно внутри `g_object_unref` во время основного цикла GTK) в длительно работающих приложениях Linux при частой загрузке ресурсов и медиафайлов. Сервер ресурсов завершал обработку `WebKitURISchemeRequest` из рабочих горутин, вызывая не потокобезопасные функции WebKit2GTK вне главного потока GTK; теперь завершение (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) выполняется в главном потоке. Это завершает частичное исправление из #5566. Затронуты сборки как с GTK3, так и с GTK4/WebKitGTK 6.0 (#5631, #5557)
- Исправлена периодическая ошибка `fatal error: invalid pointer found on stack` в `setupSignalHandlers` в Linux/GTK3. Идентификаторы окон, передаваемые как `user_data` сигнала, хранились в локальной переменной Go типа `unsafe.Pointer`, поэтому сборщик мусора аварийно завершал работу, когда при копировании стека сканировал это значение, не являющееся указателем. Теперь на стороне Go идентификатор сохраняется как целое число (`uintptr_t`): в устаревший код GTK3 перенесено то же исправление, которое в #4958 было применено к коду GTK4 (там функции сигналов C переведены на `uintptr_t` для устранения ошибок `-race`/checkptr) (#5631)

## Удалено

- Удалены стартовые шаблоны `react-swc`, `preact`, `lit`, `solid`, `qwik` и `sveltekit` (включая их варианты `-ts`). Теперь поддерживаемый встроенный набор включает `vanilla`, `react`, `vue` и `svelte` — по умолчанию каждый использует TypeScript, также доступны варианты JavaScript `-js`. Любой другой фреймворк по-прежнему можно использовать, [подключив собственный фронтенд](https://v3.wails.io/guides/dev/frontend-frameworks) или пользовательский шаблон

## v3.0.0-alpha2.104 — 2026-06-18

## Исправлено

- Исправлен сбой iOS (SIGABRT), возникавший, когда метод привязанного сервиса Go возвращал пустую строку. Обработчик ответа с ресурсом в iOS проверял указатель на тело с помощью `buf != nil` вместо проверки его длины, поэтому тело нулевой длины вызывало панику в `&buf[0]`; теперь проверяется длина, как и в обработчиках для настольных платформ

## v3.0.0-alpha2.103 — 2026-06-15

## Изменено

- Нативные возможности iOS и Android перенесены в менеджеры платформ: теперь вызывайте их через `application.IOS.*` и `application.Android.*` (например, `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`) вместо прежних свободных функций `application.IOS*`/`application.Android*` (#5602)
- Переименованы события мобильного моста: кроссплатформенные события теперь используют префикс `common:*` (например, `common:haptic`, `common:location`), а события только для определённой платформы — `ios:*` / `android:*` (например, `ios:backgroundTask`, `android:foregroundService`); префикс `native:*` больше не используется (#5602)

## v3.0.0-alpha.102 — 2026-06-14

## Добавлено

- Добавлен экспериментальный мастер `wails3 setup` для интерактивной настройки проекта и проверки зависимостей
- В `wails3 doctor` добавлен флаг `--json` для вывода в машиночитаемом формате
- В команду `wails3 doctor` добавлен раздел со статусом подписывания

## Исправлено

- Исправлено обнаружение npm в Linux: теперь помимо менеджера пакетов проверяется PATH

## v3.0.0-alpha.101 — 2026-06-13

## Добавлено

- iOS: нативные диалоговые окна сообщений (UIAlertController) и диалоговые окна открытия файла, файлов и каталога (UIDocumentPickerViewController); диалоговые окна сохранения возвращают явную ошибку
- iOS: поддержка буфера обмена через UIPasteboard
- iOS: реальные параметры экрана через UIScreen (точки, пиксели, масштаб и рабочая область с учётом безопасной зоны)
- iOS: сборки для устройств (`IOS_PLATFORM=device`), поддержка идентификатора для подписывания кода, профиля подготовки и прав доступа, упаковка `.ipa` и `deploy-device` через devicectl
- iOS: настраиваемая минимальная версия iOS (`ios.minIOSVersion` в build/config.yml)
- iOS: `wails3 doctor` сообщает о доступности Xcode и iOS SDK в macOS
- iOS: системные события — события батареи, сети, темы, блокировки экрана и нехватки памяти предоставляются как события приложения `events.IOS.*` и не зависящие от платформы `events.Common.*`
- iOS: нативный мост мобильных возможностей (экспортируемый `application.IOS*`) — меню «Поделиться», открытие URL, предотвращение перехода в спящий режим, фонарик, отступы безопасной зоны, яркость, сведения о приложении, блокировка ориентации, строка состояния, биометрия (Face ID/Touch ID), локальные уведомления и защищённое хранилище Keychain
- iOS: датчики и аппаратные возможности — тактильная обратная связь, однократное определение геопозиции, акселерометр, датчик приближения, синтез речи, сведения о хранилище, состояние питания и батареи, состояние сети, отступы для клавиатуры и обнаружение захвата экрана
- iOS: документация (IOS.md и руководство на сайте документации)
- Android: нативные диалоговые окна сообщений (AlertDialog) и диалоговые окна открытия файла и файлов (Storage Access Framework, файлы импортируются в виде копий в кэше); диалоговые окна открытия каталога и сохранения возвращают явную ошибку
- Android: поддержка буфера обмена через ClipboardManager
- Android: реальные параметры экрана через WindowMetrics/DisplayMetrics (dp, пиксели, масштаб и рабочая область с учётом системных панелей)
- Android: методы среды выполнения для тактильной обратной связи (`Android.Haptics.Vibrate`), сведений об устройстве (`Android.Device.Info`) и всплывающих уведомлений toast (`Android.Toast.Show`)
- Android: типизированные события жизненного цикла (`events.Android.*`, создаваемые из events.txt), при этом `ActivityCreated` сопоставляется с `Common.ApplicationStarted`
- Android: конвейер сборки создаёт устанавливаемые отладочные и релизные APK (`android:run`, `android:package`, `android:package:fat`); по умолчанию релиз подписывается отладочным хранилищем ключей, а настоящее хранилище ключей можно указать через переменные среды `ANDROID_KEYSTORE_*`
- Android: `wails3 doctor` сообщает об Android SDK, NDK и JDK
- Android: системные события — события батареи, сети, темы, блокировки экрана и нехватки памяти предоставляются как события приложения `events.Android.*` и не зависящие от платформы `events.Common.*`
- Android: нативный мост мобильных возможностей (экспортируемый `application.Android*`) — общий доступ, открытие URL, предотвращение перехода в спящий режим, фонарик, отступы безопасной зоны, яркость, сведения о приложении, блокировка ориентации, строка состояния, биометрия (BiometricPrompt), локальные уведомления и защищённое хранилище EncryptedSharedPreferences
- Android: датчики и аппаратные возможности — тактильная обратная связь, однократное определение геопозиции, акселерометр, датчик приближения, синтез речи, сведения о хранилище, состояние питания и батареи, состояние сети, отступы для клавиатуры и блокировка захвата экрана с помощью FLAG_SECURE
- Android: документация (ANDROID.md и руководство на сайте документации)
- Пример: в демонстрационном приложении `mobile` добавлены вкладки «Мобильные устройства» и «Оборудование», показывающие работу моста нативных возможностей в iOS и Android (вкладки в форме пилюль переносятся на несколько строк)
- Мобильные платформы: экономия заряда батареи — акселерометр, датчик приближения, фонарик и периодические часы в примере приостанавливаются, когда приложение переходит в фоновый режим, и возобновляют работу после возвращения (Android сохраняет процесс работающим в фоне, а состояние фонарика в iOS сохраняется на уровне оборудования); получатели системных событий Android регистрируются только тогда, когда приложение находится на переднем плане
- iOS: съёмка камерой — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → событие `native:capture` с миниатюрой в формате base64)
- iOS: фоновое выполнение — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (ограниченный интервал времени для фоновой задачи UIApplication) и настраиваемый параметр `ios.backgroundModes` (build/config.yml), который подставляет `UIBackgroundModes` в создаваемый файл Info.plist
- Android: съёмка камерой — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (системная камера через FileProvider → событие `native:capture`)
- Android: служба переднего плана — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (`WailsForegroundService` с постоянным уведомлением поддерживает процесс активным для длительной работы в фоновом режиме)
- Пример: вкладка «Камера», демонстрирующая съёмку фото и видео, а также фоновое выполнение (служба переднего плана в Android, ограниченный интервал времени для фоновой задачи в iOS)

## Исправлено

- Исправлен постоянный сбой `getUserMedia` с ошибкой `NotAllowedError` в Linux: WebKitGTK отклоняет запросы разрешений, если их никто не обрабатывает, а сигнал `permission-request` не был подключён. Теперь доступ к камере и микрофону обрабатывается с помощью новой кроссплатформенной карты `WebviewWindowOptions.Permissions` (`map[PermissionType]Permission`), которая учитывается как в Linux (WebKitGTK), так и в Windows (WebView2). В Linux, где нет нативного запроса разрешения, доступ к камере и микрофону по умолчанию разрешён (что восстанавливает `getUserMedia`), но его можно отключить с помощью `PermissionDeny` (#5552)
- iOS: `GOOS=ios` снова компилируется (экспортирован `events.IOS`, добавлены заглушки имён методов для мобильных платформ); также теперь компилируются сборки с тегом production (исправлены теги сборки в pkg/application и нескольких сервисах)
- iOS: события Go→JS и ExecJS теперь работают — страница больше не загружается дважды при запуске, а рукопожатие `wails:runtime:ready` больше не может быть потеряно
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` больше не вступают в состояние гонки с запуском приложения; удалена фиксированная задержка запуска длительностью 2 с
- iOS: устранена утечка C-строки при каждом выполнении JavaScript из Go
- iOS: `hasListeners` теперь отражает фактическую регистрацию слушателей
- iOS: отладочное журналирование фреймворка исключается при компиляции production-сборок
- Android: `GOOS=android` снова компилируется — определён `events.Android`, удалён выходящий за границы массив слушателей `events_android.go`, добавлена заглушка имени метода для мобильных платформ, а файлы для настольной версии Linux (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) больше не попадают в сборки Android
- Android: привязки JS→Go теперь работают — WebView не может передавать тела POST-запросов `fetch()` в `shouldInterceptRequest`, поэтому вызовы среды выполнения направляются через транспорт JavascriptInterface (`nativeHandleRuntimeCall`), а не завершаются аварийно из-за тела запроса, равного nil
- Android: вызовы среды выполнения `Screens.*` возвращают реальные данные — теперь ScreenManager заполняется при запуске (ранее он вообще не был подключён, поэтому `GetAll` возвращал nil)
- Android: отладочное журналирование фреймворка исключается при компиляции production-сборок, а в отладочных сборках направляется в logcat с тегом `Wails`
- Android: полноценный реестр `hasListeners`, обработка ссылок и исключений JNI, а также жизненный цикл страницы с однократной загрузкой (без двойной навигации)
- Исправлен сбой `wails3 generate bindings` с ошибкой «Access is denied» в Windows при запущенном сервере разработки Vite: теперь созданные файлы синхронизируются с выходным каталогом, а не переименовываются поверх него (#5515)
- Исправлен периодический критический сбой в macOS при чтении сведений об экране после изменения конфигурации дисплеев: идентификатор и имя экрана хранили указатели на автоматически освобождаемые буферы `UTF8String`, которые могли быть освобождены до их копирования в Go (use-after-free). Теперь для строк вызывается `strdup`, а после преобразования они освобождаются; перечисление экранов выполняется в явном пуле автоматического освобождения, поэтому при вызове из горутин Go утечек больше нет (#5556)
- Исправлен периодический SIGSEGV в Linux при закрытии `WebKitURISchemeRequest` сервером ресурсов: финальный `g_object_unref` выполнялся в горутине сервера ресурсов, из-за чего финализация GObject WebKit происходила вне главного потока GTK. Теперь операция unref передаётся в главный контекст GTK через `g_main_context_invoke` (#5557)

## v3.0.0-alpha.100 — 2026-06-13

## Добавлено

- В `MacWebviewPreferences` добавлены дополнительные параметры конфигурации WKWebView: `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` и `ApplicationNameForUserAgent` (#5549)

## Исправлено

- Исправлен сбой `wails3 generate bindings` с ошибкой «Access is denied» в Windows при запущенном сервере разработки Vite: теперь созданные файлы синхронизируются с выходным каталогом, а не переименовываются поверх него (#5561)
- Исправлено отсутствие событий JS при изменении размера безрамочных окон в Linux; исправлено определение края полосы прокрутки для безрамочных окон (#5368)
- Исправлен сбой средства обновления в Windows с ошибкой «invalid cross-device link», возникавший, когда временный каталог находился на другом томе, нежели каталог установки (#5560)

## v3.0.0-alpha.99 — 2026-06-10

## Исправлено

- Исправлен сбой `wails3 generate bindings` с ошибкой «Access is denied» в Windows при запущенном сервере разработки Vite: теперь созданные файлы синхронизируются с выходным каталогом, а не переименовываются поверх него (#5515)

## v3.0.0-alpha.98 — 2026-06-03

## Исправлено

- Устранено зависание пользовательского интерфейса WebKit в Linux в режиме ожидания (например, при открытом инспекторе): `SA_ONSTACK` больше не устанавливается принудительно для `SIGUSR1`, поскольку это нарушало синхронизацию потока сборщика мусора JavaScriptCore (#5527)

## v3.0.0-alpha.97 — 2026-05-31

## Добавлено

- Добавлена страница отладки и работа с `runtime/trace`

## Изменено

- Удалены некоторые ненужные импорты `_ "embed"`, что позволило немного упорядочить код

## Исправлено

- Исправлено несоблюдение ограничений минимальной ширины и высоты после выхода окна из развёрнутого состояния в Windows (#4593)
- Исправлено прохождение щелчков мыши сквозь окно в полноэкранном режиме при использовании параметров окна Frameless + Transparent (#4408)

## v3.0.0-alpha.96 — 2026-05-25

## Добавлено

- Добавлена поддержка обфускации с помощью Garble ([#4563](https://github.com/wailsapp/wails/issues/4563)): стабильные идентификаторы методов привязок, интеграция со сборкой и Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`), а также теги структур JSON для всех данных, доступных среде выполнения (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`), благодаря чему формат обмена сохраняется при переименовании экспортируемых полей средствами Garble.

## v3.0.0-alpha.95 — 2026-05-20

## Добавлено

- Добавлена отсутствовавшая страница о структуре проекта

## Изменено

- Документация: несколько диаграмм на странице архитектуры заменены диаграммами последовательности для более наглядного отображения
- Документация: добавлено примечание об установке D2 как необходимого условия для запуска

## Исправлено

- Исправлена работа `wails3 generate appimage` со стеком GTK4 по умолчанию: теперь упаковщик определяет стек GTK по исполняемому файлу до поиска файлов среды выполнения, поэтому для сборок GTK4 он выбирает `libwebkitgtkinjectedbundle.so` (в `webkitgtk-6.0/`), а для сборок `-tags gtk3` — `libwebkit2gtkinjectedbundle.so` (в `webkit2gtk-4.1/`). Проверка `.relr.dyn` также учитывает `libgtk-4.so.1`, поэтому удаление отладочной информации корректно отключается в современных наборах инструментов независимо от стека. (#5475)
- Исправлен сбой `wails3 generate appimage` при вызове с относительным `-builddir`: теперь упаковщик заранее преобразует `-binary`, `-icon`, `-desktopfile`, `-builddir` и `-outputdir` в абсолютные пути, чтобы выполняемый в середине процесса `s.CD` не нарушал работу горутины загрузки AppRun и проверки `ldd` после копирования.
- Исправлена ошибка, из-за которой `wails3 generate appimage` не удавалось переместить итоговый AppImage в `-outputdir`, если поле `Name=` desktop-файла не совпадало с базовым именем исполняемого файла: теперь упаковщик через переменную среды `OUTPUT` предписывает плагину appimage из linuxdeploy записывать AppImage в `<binary>-<arch>.AppImage`, а не использовать имя, полученное из desktop-файла.
- Исправлена ошибка, из-за которой `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` и `Common.SystemDidWake` не срабатывали в Linux после того, как стек GTK4 + WebKitGTK 6.0 стал использоваться по умолчанию в alpha.93. Новый используемый по умолчанию `application_linux.go` `run()` не вызывал `setupCommonEvents()` (который перенаправляет события `Linux.*` в соответствующие события `Common.*`) и `monitorPowerEvents()`. Теперь вспомогательный компонент отслеживания питания через DBus используется совместно путями сборки GTK3 и GTK4 посредством `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.94 — 2026-05-19

## Исправлено

- Исправлена ошибка, из-за которой `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` и `Common.SystemDidWake` не срабатывали в Linux после того, как стек GTK4 + WebKitGTK 6.0 стал использоваться по умолчанию в alpha.93. Новый используемый по умолчанию `application_linux.go` `run()` не вызывал `setupCommonEvents()` (который перенаправляет события `Linux.*` в соответствующие события `Common.*`) и `monitorPowerEvents()`. Теперь вспомогательный компонент отслеживания питания через DBus используется совместно путями сборки GTK3 и GTK4 посредством `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.93 — 2026-05-17

## Добавлено

- В вывод `wails3 doctor` в Linux добавлен `XDG_SESSION_TYPE`, автор — @leaanthony

## Исправлено

- Исправлен сбой меню окна в Wayland, вызванный обращением appmenu-gtk-module к окну, нативные ресурсы которого ещё не созданы (#4769), автор — @leaanthony
- Исправлен сбой приложения GTK, если имя приложения содержит недопустимые символы (пробелы, круглые скобки и т. п.), автор — @leaanthony
- Исправлена ошибка «недостаточно памяти» при инициализации перетаскивания в Windows (#4701), автор — @overlordtm
- Устранено состояние гонки в хранилище обратных вызовов основного потока, где для удаления из map ошибочно использовалась RLock (Linux, macOS, iOS) (#4424), автор — @leaanthony
- Исправлена обработка переменных при передаче задачам аргументов командной строки. Теперь переменные CLI, заданные парами KEY=VALUE, корректно инициализируются и передаются на всех этапах выполнения задачи.
- Устранён конфликт NSWindowZoomButton в macOS: теперь `MaximiseButtonState` и `FullscreenButtonState` применяют более ограничивающее состояние как при запуске, так и во время выполнения; ни один из сеттеров больше не может незаметно переопределить другой (#5319)
- Исправлена группа ранее существовавших ошибок в устаревшем пути сборки GTK3 (`-tags gtk3`), выявленных CodeRabbit в #5463: при запуске через файловые ассоциации обработчики запуска больше не пропускаются; `getTheme` теперь безопасен с точки зрения границ и типов; `appName` больше не освобождает память, принадлежащую GLib; `clipboardGet` больше не приводит к утечке возвращаемого GTK объекта `gchar*`; теперь `Calloc` использует получатели-указатели (а `NewCalloc` возвращает `*Calloc`), поэтому пул действительно отслеживает выделение памяти; `zoomOut` использует величину, обратную `zoomInFactor`, вместо отрицательного множителя, из-за которого значение ограничивалось до 1.0; `execJS` повторно использует заранее выделенное пустое имя мира вместо утечки одного `C.CString("")` при каждом вызове; из `menuItem.setAccelerator` удалён отладочный `fmt.Println`. Решает #5465.
- Исправлена такая же утечка из-за получателя-значения `Calloc` в используемом по умолчанию пути сборки GTK4 (`linux_cgo.go`): получатели-указатели и `NewCalloc() *Calloc` обеспечивают фактическое отслеживание и освобождение выделений `c.String(...)` для каждого окна.

## v3.0.0-alpha.92 — 2026-05-15

## Добавлено

- Файлы Taskfile изменены так, чтобы используемым менеджером пакетов фронтенда можно было управлять с помощью параметра `PACKAGE_MANAGER`
- Данные шаблонов дополнены `{{.Opn}}` и `{{.Cls}}`, чтобы сделать создание шаблонов Taskfile более предсказуемым

## Изменено

- Несколько существующих файлов Taskfile изменены для использования `{{.Opn}} and {{.Cls}}`

## Исправлено

- Устранена фатальная ошибка среды выполнения `concurrent map read and map write` в `linuxSystemTray`, возникавшая, когда меню в области уведомлений обновлялось во время его чтения панелью.
- Для вывода ошибки и трассировки стека WebView2 теперь используется `log` вместо `fmt`, чтобы сообщения не терялись при работе приложения в Windows без подключённой консоли.

## v3.0.0-alpha.91 — 2026-05-12

## Изменено

- Обновлён SVG-файл спонсоров в [PR](https://github.com/wailsapp/wails/pull/5414), автор — `@github-actions[bot]`
- **ИЗМЕНЕНИЕ, НАРУШАЮЩЕЕ ОБРАТНУЮ СОВМЕСТИМОСТЬ (macOS):** система координат macOS нормализована, чтобы `GetScreens`, `Position` и `SetPosition` использовали единое пространство: логические точки, ось Y направлена вниз, а `(0,0)` находится в левом верхнем углу основного экрана. Это соответствует Windows, GTK, а также общедоступным API Electron и веб-платформы. Для экранов, физически расположенных выше основного, теперь возвращается отрицательное значение `Bounds.Y` (ранее положительное), а значения `Position()`/`SetPosition()` теперь выражаются в логических точках вместо `points × primaryScale`. Корректность преобразования туда и обратно `Position()` → `SetPosition()` сохранена; абсолютные значения, записанные в журналы предыдущими alpha-сборками, и обходные решения с вычислениями вручную (например, умножением на `primaryScale` или инвертированием Y относительно высоты экрана) потребуется обновить. Решает [#5117](https://github.com/wailsapp/wails/issues/5117).

## Исправлено

- Добавлена защитная проверка имени сигнала DBus и длины его тела для предотвращения паник в [PR](https://github.com/wailsapp/wails/pull/5416), автор — @leaanthony
- Исправлена проблема с безопасностью памяти при обработке меню GTK в Linux в [PR](https://github.com/wailsapp/wails/pull/5363), автор — @leaanthony
- Добавлено обнаружение графических процессоров NVIDIA и отключение средства визуализации DMA-BUF в Linux в [PR](https://github.com/wailsapp/wails/pull/5295), автор — @leaanthony
- Исправлено преобразование координаты Y между экранами в `SetPosition` в macOS: высота основного экрана теперь используется как глобальная точка отсчёта, поэтому окна размещаются в правильной позиции на мониторах, смещённых по вертикали относительно основного экрана — [#5117](https://github.com/wailsapp/wails/issues/5117)
- Исправлен шаблон git для PR: теперь он указывает на правильный URL для обратной связи в [PR](https://github.com/wailsapp/wails/pull/5109), автор — @wayneforrest
- Исправлено семейство сбоев `SetMenu` системного трея Windows, вызванных некорректным системным вызовом `DestroyMenu`, которому передавалось четыре аргумента вместо одного, поэтому каждый вызов возвращал FALSE и ничего не освобождал. Кроме того, при перестроении меню теперь освобождаются дескрипторы HMENU и HBITMAP (включая выделенные во время выполнения через `MenuItem.SetBitmap`), в `Win32Menu.Update` сбрасываются устаревшие карты флажков и переключателей, а из `systemtray.updateMenu` удалён избыточный вызов `Update()`, удваивавший количество выделений памяти. Длительно работающие приложения с системным треем больше не допускают утечки объектов GDI/USER при каждом перестроении меню.

## v3.0.0-alpha.90 — 2026-05-11

## Добавлено

- Добавлена возможность настраивать имя приложения в User-Agent WKWebView на macOS в [PR](https://github.com/wailsapp/wails/pull/5261), автор: @vinhvoit225
- В пример gin-service добавлена косвенная зависимость github.com/coder/websocket в [PR](https://github.com/wailsapp/wails/pull/5400), автор: @taliesin-ai
- В тесты ресурсов сборки добавлена поддержка глубокого сравнения на равенство в [PR](https://github.com/wailsapp/wails/pull/5402), автор: @leaanthony

## Изменено

- Выходные данные сборки объединены в каталоге assets в [PR](https://github.com/wailsapp/wails/pull/5401), автор: @taliesin-ai
- Обновлён SVG-файл спонсоров в [PR](https://github.com/wailsapp/wails/pull/5399), автор — `@github-actions[bot]`

## Исправлено

- Для сообщения единственного экземпляра приложения на macOS теперь используется объект уведомления в [PR](https://github.com/wailsapp/wails/pull/5289), автор: @overlordtm
- Обратные вызовы Windows теперь объединяются в пакеты, чтобы предотвратить потерю промисов при высокой нагрузке, в [PR](https://github.com/wailsapp/wails/pull/5383), автор: @taliesin-ai

## v3.0.0-alpha.89 — 2026-05-10

## Добавлено

- Добавлено задание go<em>test</em>results для агрегирования результатов тестов Go в [PR](https://github.com/wailsapp/wails/pull/5316), автор: @leaanthony

## Изменено

- Добавлено условное разделение крупных полезных данных RPC на фрагментированные POST-запросы в [PR](https://github.com/wailsapp/wails/pull/5369), автор: @leaanthony
- Во всех шаблонах фронтенда Vite обновлён с версии 5.x.x до 8.0.0 в [PR](https://github.com/wailsapp/wails/pull/5386), автор: @leaanthony
- Конфигурация порта сервера разработки Vite перенесена в переменные окружения в [PR](https://github.com/wailsapp/wails/pull/5365), автор: @leaanthony
- Во всех шаблонах сервер разработки Vite настроен на привязку к 127.0.0.1 в [PR](https://github.com/wailsapp/wails/pull/5361), автор: @leaanthony
- Обновлён SVG-файл спонсоров в [PR](https://github.com/wailsapp/wails/pull/5384), автор — `@github-actions[bot]`

## Исправлено

- При обновлении build-assets теперь очищаются заготовки шаблона Info.plist в [PR](https://github.com/wailsapp/wails/pull/5312), автор: @leaanthony
- Исправлено устаревшее состояние меню macOS: методы изменения пунктов меню (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) теперь применяются синхронно в основном потоке. Это устраняет состояние гонки `dispatch_async`, из-за которого при быстром повторном открытии меню отображалось его предыдущее состояние (#5002)
- В режиме разработки файлы `*_test.go` теперь игнорируются, чтобы предотвратить ненужные повторные сборки, в [PR](https://github.com/wailsapp/wails/pull/5203), автор: @leaanthony
- Предотвращена ошибка сегментации в Menu.Update(), когда приложение не запущено, в [PR](https://github.com/wailsapp/wails/pull/5291), автор: @wucm667
- В Windows перерисовка строки меню теперь выполняется в зависимости от lastSizeWParam в [PR](https://github.com/wailsapp/wails/pull/5382), автор: @taliesin-ai

## v3.0.0-alpha.88 — 2026-05-09

## Изменено

- HiddenOnTaskbar переведён на использование WS<em>EX</em>TOOLWINDOW в [PR](https://github.com/wailsapp/wails/pull/5371), автор: @leaanthony
- Изменён порядок зависимостей и удалена директива replace для webview2 из go.mod в [PR](https://github.com/wailsapp/wails/pull/5370), автор: @atterpac
- Обновлён SVG-файл спонсоров в [PR](https://github.com/wailsapp/wails/pull/5358), автор — `@github-actions[bot]`

## Исправлено

- Удалены обобщённые псевдонимы косвенного обращения, а типы ключей карт унифицированы в [PR](https://github.com/wailsapp/wails/pull/5331), автор: @fbbdev

## Удалено

- Удалён рабочий процесс PR-master вместе с документацией, тестами Go и механизмом пропуска тестов в [PR](https://github.com/wailsapp/wails/pull/5377), автор: @leaanthony

## v3.0.0-alpha.87 — 2026-05-07

## Добавлено

- Добавлена документация Wails v3 на корейском языке в [PR](https://github.com/wailsapp/wails/pull/5352), автор: @leaanthony
- Добавлена документация на французском языке по установке и быстрому началу работы в [PR](https://github.com/wailsapp/wails/pull/5354), автор: @leaanthony
- Добавлена документация на португальском языке по быстрому началу работы, концепциям и сообществу в [PR](https://github.com/wailsapp/wails/pull/5355), автор: @leaanthony

## v3.0.0-alpha.86 — 2026-05-06

## Добавлено

- Добавлена французская локализация документации в [PR](https://github.com/wailsapp/wails/pull/5328), автор: @leaanthony
- На сайт документации добавлена немецкая локаль в [PR](https://github.com/wailsapp/wails/pull/5343), автор: @leaanthony

## Изменено

- Все 8 переведённых локалей зарегистрированы в конфигурации документации в [PR](https://github.com/wailsapp/wails/pull/5347), автор: @leaanthony
- Обновлены различные связанные с Windows файлы для WebView2 в [PR](https://github.com/wailsapp/wails/pull/5317), автор: @leaanthony

## Исправлено

- Диспетчеризация диалогов в Linux разделена между GTK3 и GTK4 в [PR](https://github.com/wailsapp/wails/pull/5340), автор: @leaanthony
- Выполнение обратных вызовов диалогов в потоке GTK теперь гарантировано, что устраняет ошибки сегментации, в [PR](https://github.com/wailsapp/wails/pull/5339), автор: @leaanthony

## v3.0.0-alpha.85 — 2026-05-05

## Добавлено

- В репозиторий добавлен URL-адрес шаблона PR в [PR](https://github.com/wailsapp/wails/pull/5179), автор: @leaanthony
- Добавлена документация Wails v3 на немецком языке в [PR](https://github.com/wailsapp/wails/pull/5330), автор: @leaanthony

## v3.0.0-alpha.84 — 2026-05-03

## Добавлено

- Добавлена возможность отключить выход из полноэкранного режима по клавише Escape на macOS в [PR](https://github.com/wailsapp/wails/pull/5307), автор: @leaanthony
- Добавлена возможность отключить выход из полноэкранного режима по клавише Escape на macOS в [PR](https://github.com/wailsapp/wails/pull/5310), автор: @leaanthony
- Добавлена документация о проекте Pausa в галерее сообщества в [PR](https://github.com/wailsapp/wails/pull/5288) от @yuseferi

## Изменено

- Обновлён SVG-файл со спонсорами в [PR](https://github.com/wailsapp/wails/pull/5308), автор — `@github-actions[bot]`
- Обновлена команда создания значков для обработки неподдерживаемых платформ в [PR](https://github.com/wailsapp/wails/pull/5309) от @leaanthony
- Логический API полноэкранного режима заменён на трёхпозиционный ButtonState, реализованы привязки для платформ в [PR](https://github.com/wailsapp/wails/pull/5224) от @leaanthony

## Исправлено

- Операции фокусировки WebView2 защищены от состояния, при котором контроллер равен nil, в [PR](https://github.com/wailsapp/wails/pull/5315) от @leaanthony
- Рабочий процесс GitHub Actions обновлён для правильного обращения к базовой ветке PR в [PR](https://github.com/wailsapp/wails/pull/5313) от @leaanthony
- Файлы `*_test.go` игнорируются в режиме разработки во избежание лишних пересборок в [PR](https://github.com/wailsapp/wails/pull/5203) от @leaanthony
- Предотвращена ошибка сегментации в Menu.Update(), когда приложение не запущено, в [PR](https://github.com/wailsapp/wails/pull/5291) от @wucm667

## v3.0.0-alpha.83 — 2026-05-02

## Добавлено

- Добавлены флаг InstallScope и параметр сборки для установки на уровне компьютера или пользователя в [PR](https://github.com/wailsapp/wails/pull/5094) от @symball
- В BrowserWindow добавлен ничего не выполняющий метод SetScreen для соответствия интерфейсу Window в [PR](https://github.com/wailsapp/wails/pull/5294) от @leaanthony

## Исправлено

- Добавлено обнаружение графических процессоров NVIDIA и отключение средства визуализации DMA-BUF в Linux в [PR](https://github.com/wailsapp/wails/pull/5295) от @leaanthony
- Исправлен шаблон PR для Git: теперь он указывает правильный URL для обратной связи, в [PR](https://github.com/wailsapp/wails/pull/5109) от @wayneforrest
- Исправлена группа сбоев `SetMenu` в области уведомлений Windows, вызванных неисправным системным вызовом `DestroyMenu`, который передавал четыре аргумента вместо одного, поэтому каждый вызов возвращал FALSE и ничего не освобождал. Кроме того, при перестроении меню теперь освобождаются дескрипторы HMENU и HBITMAP, включая выделенные во время выполнения через `MenuItem.SetBitmap`; в `Win32Menu.Update` сбрасываются устаревшие таблицы флажков и переключателей; а из `systemtray.updateMenu` удалён избыточный вызов `Update()`, удваивавший количество выделений памяти. Длительно работающие приложения с областью уведомлений больше не вызывают утечку объектов GDI/USER при каждом перестроении меню.

## v3.0.0-alpha.82 — 2026-05-01

## Исправлено

- Исправлено создание desktop-файла: теперь имя записи приложения обрабатывается правильно, в [PR](https://github.com/wailsapp/wails/pull/5232) от @leaanthony

## v3.0.0-alpha.81 — 2026-04-30

## Изменено

- Расписание ночных выпусков перенесено на 15:00 UTC в [PR](https://github.com/wailsapp/wails/pull/5286) от @leaanthony

## Исправлено

- Исправлены уменьшенные вдвое значения Screen Bounds, WorkArea и Size на компьютерах Mac с дисплеем Retina —  (#5168)

## v3.0.0-alpha.80 — 2026-04-29

## Изменено

- Обновлены зависимости документации и загрузчики коллекций содержимого в [PR](https://github.com/wailsapp/wails/pull/5285) от @leaanthony

## v3.0.0-alpha.79 — 2026-04-29

## Добавлено

- Заданию trigger-release предоставлено разрешение actions: write в [PR](https://github.com/wailsapp/wails/pull/5270) от @leaanthony

## Изменено

- Для задачи выпуска веткой по умолчанию назначена master и обновлена формулировка журнала изменений в [PR](https://github.com/wailsapp/wails/pull/5283) от @leaanthony
- Рабочий процесс автоматического создания журнала изменений обновлён для использования последней версии в [PR](https://github.com/wailsapp/wails/pull/5282) от @leaanthony
- Повышена эффективность рабочих процессов за счёт добавления фильтров путей и удаления неиспользуемых рабочих процессов в [PR](https://github.com/wailsapp/wails/pull/5280) от @leaanthony
- В документации для ссылок на примеры теперь указана ветка master в [PR](https://github.com/wailsapp/wails/pull/5274) от @leaanthony
- Обновлены документация и примеры для v3 в [PR](https://github.com/wailsapp/wails/pull/5272) от @leaanthony

## Исправлено

- В обратный прокси-сервер добавлены повторные попытки и принудительное использование IPv4 при разработке в [PR](https://github.com/wailsapp/wails/pull/5265) от @AkagiYui
- Переработан рабочий процесс запуска обновления раздела «Невыпущенные изменения» в журнале изменений — [PR](https://github.com/wailsapp/wails/pull/5281) от @leaanthony

## Удалено

- Удалены сценарии оболочки для различных задач тестирования в [PR](https://github.com/wailsapp/wails/pull/5267) от @leaanthony
- Удалены рабочий процесс развёртывания документации v3-alpha и запись CNAME в [PR](https://github.com/wailsapp/wails/pull/5266) от @leaanthony

### Добавлено

- В боковую панель навигации добавлен пункт «Маршрутизация во фронтенде» в [PR](https://github.com/wailsapp/wails/pull/5196) от @leaanthony
- Добавлено руководство по маршрутизации во фронтенде с рекомендациями для конкретных фреймворков в [PR](https://github.com/wailsapp/wails/pull/5185) от @leaanthony
- Добавлена поддержка модальных листов (macOS)
- Версия ghw повышена для улучшения поддержки устройств Apple; автор — @leaanthony (#4977)
- В службу Dock добавлен метод `GetBadge`
- В команду `wails3 build` добавлен флаг `-tags` для передачи пользовательских тегов сборки Go (например, `wails3 build -tags gtk4`) (#4957)
- Добавлена документация по автоматическому созданию перечислений в генераторе привязок, включая отдельную страницу «Перечисления» и пункт в боковой панели навигации (#4972)
- В команду `wails3 build` добавлен флаг `-tags` для передачи пользовательских тегов сборки Go (например, `wails3 build -tags gtk4`) (#4957)
- В `v3/examples/web-apis/` добавлены примеры Web API, демонстрирующие 41 браузерный API, включая хранилища (localStorage, sessionStorage, IndexedDB, Cache API), сетевые возможности (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), мультимедиа (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), возможности устройств (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), производительность (Performance API, Mutation Observer, Intersection/Resize Observer), пользовательский интерфейс (Web Components, Pointer Events, Selection, Dialog, Drag and Drop) и многое другое
- Добавлен пример средства проверки совместимости API WebView (`v3/examples/webview-api-check/`), которое тестирует более 200 API браузера на разных платформах
- Добавлен пакет `internal/libpath` для поиска путей к нативным библиотекам в Linux с параллельным поиском, кэшированием и поддержкой Flatpak, Snap и Nix
- **В разработке:** добавлена экспериментальная поддержка WebKitGTK 6.0 / GTK4 для Linux, доступная через `-tags gtk4` (GTK3/WebKit2GTK 4.1 по-прежнему используется по умолчанию)
- Примечание. В мозаичных оконных менеджерах (например, Hyprland и Sway) операции сворачивания и разворачивания могут работать не так, как ожидается, поскольку геометрией окон управляет оконный менеджер
- В раздел документации «**Прослушивание событий в JavaScript**» добавлены инструкции по использованию **одноразовых обработчиков** от @AbdelhadiSeddar
- В `WebviewWindowOptions` добавлен параметр `UseApplicationMenu`, позволяющий окнам в Windows/Linux наследовать меню приложения, заданное через `app.Menu.Set()`, от @leaanthony
- Добавлена поддержка файлов `.icon` (формат Apple Icon Composer) для создания значков Liquid Glass и каталогов ресурсов (macOS) (#4934) от @wimaha
- Добавлен экспериментальный серверный режим для безголовых и веб-развёртываний (`-tags server`). Он позволяет запускать приложения Wails как HTTP-серверы без зависимостей от нативного графического интерфейса. Выполняйте сборку с `wails3 task build:server`. Подробности см. в `examples/server`.
- Добавлен пакет `internal/libpath` для поиска путей к нативным библиотекам в Linux с параллельным поиском, кешированием и поддержкой Flatpak/Snap/Nix
- В `MacWindow` добавлен параметр `CollectionBehavior` для управления поведением окна в пространствах macOS Spaces и полноэкранном режиме (#4756) от @leaanthony
- Добавлены модульные тесты для pkg/application от @leaanthony
- В упаковку MSIX добавлена поддержка пользовательских протоколов от @leaanthony
- Добавлено определение окружения рабочего стола в Linux [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- В среду выполнения JavaScript добавлен метод `Window.Print()` для открытия диалогового окна печати из фронтенда (#4290) от @leaanthony
- В вывод `wails3 doctor` в Linux добавлен `XDG_SESSION_TYPE` от @leaanthony
- Добавлены дополнительные события изменения состояния загрузки WebKit2 для Linux: `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` (#3896) от @leaanthony
- В вывод `wails3 doctor` в Linux добавлен `XDG_SESSION_TYPE` от @leaanthony
- Файл `.desktop` теперь создаётся во время сборки для Linux, а не только при упаковке (#4575)
- Добавлена документация по зависимостям среды выполнения Linux с названиями пакетов для отдельных дистрибутивов и примерами упаковки с помощью nfpm (#4339) от @leaanthony
- В вывод `wails3 doctor` в Linux добавлена информация о версии драйвера NVIDIA от @leaanthony
- В обработчик необработанных сообщений добавлен источник сообщения от @APshenkin в [PR](https://github.com/wailsapp/wails/pull/4710)
- Добавлена поддержка универсальных ссылок для macOS от @APshenkin в [PR](https://github.com/wailsapp/wails/pull/4712)
- Рефакторинг транспортного слоя привязок от @APshenkin в [PR](https://github.com/wailsapp/wails/pull/4702)
- В шаблоны helloworld добавлены идентификаторы aria-label, чтобы пример приложения можно было легко тестировать с помощью клиентов Appium, от @chinenual в [PR](https://github.com/wailsapp/wails/pull/4760)
- В обработчик необработанных сообщений добавлен источник сообщения от @APshenkin в [PR](https://github.com/wailsapp/wails/pull/4710)
- Добавлена поддержка универсальных ссылок для macOS от @APshenkin в [PR](https://github.com/wailsapp/wails/pull/4712)
- Рефакторинг транспортного слоя привязок от @APshenkin в [PR](https://github.com/wailsapp/wails/pull/4702)
- Типизированные события от @fbbdev и @ianvs в [#4633](https://github.com/wailsapp/wails/pull/4633)
- Добавлен пример `systray-clock`, демонстрирующий работу приложения в области уведомлений без основного окна с динамическим обновлением всплывающей подсказки (#4653).
- Добавлен шаблон протокола NSIS для Windows от @Tolfx в #4510
- Добавлены тесты для build-assets от @Tolfx в #4510
- macOS: нативные элементы управления окном отображаются в строке меню в [#4588](https://github.com/wailsapp/wails/pull/4588) от @nidib
- Добавлена служба Dock для macOS, позволяющая скрывать и отображать значок приложения в Dock, от @popaprozac в [PR](https://github.com/wailsapp/wails/pull/4451)
- Добавлена служба Dock для macOS, позволяющая скрывать и отображать значок приложения в Dock, от @popaprozac в [PR](https://github.com/wailsapp/wails/pull/4451)
- Добавлена нативная поддержка эффекта Liquid Glass для macOS с NSGlassEffectView (macOS 15.0+) и резервным использованием NSVisualEffectView, включая широкие возможности настройки материала, от @leaanthony в [#4534](https://github.com/wailsapp/wails/pull/4534)
- Санитизация URL браузера от @leaanthony в [#4500](https://github.dev/wailsapp/wails/pull/4500). Основано на [#4484](https://github.com/wailsapp/wails/pull/4484) от @APShenkin.
- Добавлена защита содержимого в Windows/Mac от [@leaanthony](https://github.com/leaanthony) на основе исходной работы [@Taiterbase](https://github.com/Taiterbase) в этом [PR](https://github.com/wailsapp/wails/pull/4241)
- Добавлена поддержка передачи переменных CLI командам Task через псевдонимы `wails3 build` и `wails3 package` (#4422) от @leaanthony в [PR](https://github.com/wailsapp/wails/pull/4488)
- Поддержка зон перетаскивания, в которых данные о перетаскиваемом элементе поступают из события, от [@atterpac](https://github.com/atterpac) в [#4318](https://github.com/wailsapp/wails/pull/4318)
- В параметры `WindowsWindow` добавлен `AdditionalLaunchArgs`, позволяющий передавать браузеру WebView2 дополнительные аргументы командной строки, в [PR](https://github.com/wailsapp/wails/pull/4467)
- Добавлен автоматический запуск go mod tidy после wails init от [@triadmoko](https://github.com/triadmoko) в [PR](https://github.com/wailsapp/wails/pull/4286)
- Функция Windows Snap Assist от @leaanthony в [PR](https://github.dev/wailsapp/wails/pull/4463)
- В параметры `WindowsWindow` добавлен `AdditionalLaunchArgs`, позволяющий передавать браузеру WebView2 дополнительные аргументы командной строки, в [PR](https://github.com/wailsapp/wails/pull/4467)
- Добавлен автоматический запуск go mod tidy после wails init от [@triadmoko](https://github.com/triadmoko) в [PR](https://github.com/wailsapp/wails/pull/4286)
- Функция Windows Snap Assist от @leaanthony в [PR](https://github.dev/wailsapp/wails/pull/4463)
- Добавлена реализация `getAccentColor` для Windows от [@almas-x](https://github.com/almas-x) в [PR](https://github.com/wailsapp/wails/pull/4427)
- Добавлена реализация `getAccentColor` для Windows от [@almas-x](https://github.com/almas-x) в [PR](https://github.com/wailsapp/wails/pull/4427)
- Тёмная тема меню и строки меню в Windows. Реализовано @leaanthony в [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Встроенные службы переименованы для создания более понятных привязок JS/TS, от @popaprozac в [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` для получения системного акцентного цвета пользователя. Работает в macOS. От [@etesam913](https://github.com/etesam913)
- Добавлен API `window.ToggleFrameless()` от [@atterpac](https://github.com/atterpac) в [#4137](https://github.com/wailsapp/wails/pull/4137)
- Добавлены зависимости сборки для Linux с учётом конкретного дистрибутива — @leaanthony в [PR](https://github.com/wailsapp/wails/pull/4345)
- Добавлено руководство по привязкам — @atterpac в [PR](https://github.com/wailsapp/wails/pull/4404)
- **Упорядочена инфраструктура тестирования**: тестовые файлы Docker перемещены в отдельный каталог `test/docker/`, образы оптимизированы, а надёжность сборки повышена — [@leaanthony](https://github.com/leaanthony) в [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Улучшены шаблоны управления ресурсами**: в примеры добавлены корректная очистка обработчиков событий и управление горутинами с учётом контекста — [@leaanthony](https://github.com/leaanthony) в [#4359](https://github.com/wailsapp/wails/pull/4359)
- Добавлена поддержка сборки AppImage для aarch64 — [@AkshayKalose](https://github.com/AkshayKalose) в [#3981](https://github.com/wailsapp/wails/pull/3981)
- В `wails doctor` добавлен раздел диагностики — [@leaanthony](https://github.com/leaanthony)
- При вызове метода сервиса окно добавляется в контекст — [@leaanthony](https://github.com/leaanthony)
- Добавлен пример `window-call`, демонстрирующий, как определить, какое окно вызывает сервис, — [@leaanthony](https://github.com/leaanthony)
- Новое руководство по меню — [@leaanthony](https://github.com/leaanthony)
- Улучшена обработка паник — [@leaanthony](https://github.com/leaanthony)
- Новое руководство по меню — [@leaanthony](https://github.com/leaanthony)
- Добавлены документирующие комментарии для Service API — [@fbbdev](https://github.com/fbbdev) в [#4024](https://github.com/wailsapp/wails/pull/4024)
- Добавлена функция `application.NewServiceWithOptions` для инициализации сервисов с дополнительной конфигурацией — [@leaanthony](https://github.com/leaanthony) в [#4024](https://github.com/wailsapp/wails/pull/4024)
- Улучшено управление меню — [@FalcoG](https://github.com/FalcoG) и [@leaanthony](https://github.com/leaanthony) в [#4031](https://github.com/wailsapp/wails/pull/4031)
- Добавлена дополнительная документация — [@leaanthony](https://github.com/leaanthony)
- Добавлена поддержка отмены событий в стандартных слушателях событий — [@leaanthony](https://github.com/leaanthony)
- Добавлена поддержка `Hide`, `Show` и `Destroy` для значка в области уведомлений — [@leaanthony](https://github.com/leaanthony)
- Добавлена поддержка `SetTooltip` для значка в области уведомлений — [@leaanthony](https://github.com/leaanthony). Автор исходной идеи — [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- В предупреждениях генератора привязок о неподдерживаемых типах теперь указывается путь к пакету — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- В генератор привязок добавлена поддержка обобщённых псевдонимов — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- В генератор привязок добавлена поддержка флага JSON `omitzero` — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Добавлена директива `//wails:ignore`, предотвращающая генерацию привязок для выбранных методов сервиса, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Для сервисов и моделей добавлена директива `//wails:internal`, позволяющая использовать типы, которые экспортируются в Go, но не в JS/TS, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- В генератор привязок добавлена поддержка констант типа-псевдонима, позволяющая создавать перечисления со слабой типизацией, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Добавлены тесты генератора привязок для возможностей Go 1.24 — [@fbbdev](https://github.com/fbbdev) в [#4068](https://github.com/wailsapp/wails/pull/4068)
- В `OSInfo.Branding` добавлена поддержка macOS 15 «Sequoia» для более точного определения версии ОС — [#4065](https://github.com/wailsapp/wails/pull/4065)
- Добавлен хук `PostShutdown` для выполнения пользовательского кода после завершения процесса остановки — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавлена структура `FatalError` для поддержки обнаружения фатальных ошибок в пользовательских обработчиках ошибок — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Порядок запуска и остановки сервисов стандартизирован и задокументирован — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавлены тестовая инфраструктура для последовательности запуска и остановки приложения, а также тесты запуска и остановки сервисов — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавлен метод `RegisterService` для регистрации сервисов после создания приложения — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- В параметры приложения и сервисов добавлено поле `MarshalError` для пользовательской обработки ошибок при вызовах привязок — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавлена отменяемая обёртка промиса, которая передаёт запросы на отмену по цепочкам промисов, — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Добавлена возможность привязать отмену вызова привязки к `AbortSignal` — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- В WML наряду с обычными атрибутами `wml-*` поддерживаются атрибуты `data-wml-*` — [@leaanthony](https://github.com/leaanthony)
- Во все сервисы добавлен метод `Configure` для поздней настройки и динамической перенастройки — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Если сервис `fileserver` не настроен, он отправляет ответ 503 Service Unavailable — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Если сервис `kvstore` не настроен, по умолчанию он предоставляет хранилище ключей и значений в памяти — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- В сервис `kvstore` добавлен метод `Load` для повторной загрузки данных из файла после изменения конфигурации — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- В сервис `kvstore` добавлен метод `Clear` для удаления всех ключей — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- В сервис `log` добавлен тип `Level`, предоставляющий константы уровней журналирования на стороне JS, — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- В сервис `log` добавлен метод `Log` для динамического задания уровня журналирования — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Если сервис `sqlite` не настроен, по умолчанию он предоставляет базу данных в памяти — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Добавлен метод `Close` в сервис `sqlite` для ручного закрытия БД; автор — [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- Добавлена поддержка отмены для методов выполнения запросов в сервисе `sqlite`; автор — [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- В сервис `sqlite` добавлена поддержка подготовленных выражений с привязками JS; автор — [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- Поддержка Gin от [Lea Anthony](https://github.com/leaanthony) в [PR](https://github.com/wailsapp/wails/pull/3537) на основе исходной работы [@AnalogJ](https://github.com/AnalogJ) в этом [PR](https://github.com/wailsapp/wails/pull/3537)
- Исправлена проблема, из-за которой автосохранение и автоматическое сохранение паролей всегда были включены; автор — [@oSethoum](https://github.com/osethoum), [#4134](https://github.com/wailsapp/wails/pull/4134)
- В окно добавлен `SetMenu()`, позволяющий назначить окну меню; автор — [@leaanthony](https://github.com/leaanthony)
- Добавлена поддержка уведомлений; автор — [@popaprozac](https://github.com/popaprozac), [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Добавлена поддержка ассоциаций файлов для mac; автор — [@wimaha](https://github.com/wimaha), [#4177](https://github.com/wailsapp/wails/pull/4177)
- Добавлен `wails3 tool version` для повышения семантической версии; автор — [@leaanthony](https://github.com/leaanthony)
- Добавлена поддержка меток на значке приложения для macOS и Windows; автор — [@popaprozac](https://github.com/popaprozac), [#](https://github.com/wailsapp/wails/pull/4234)
- Добавлена поддержка зарегистрированных событий со строгой типизацией; авторы — [@fbbdev](https://github.com/fbbdev) и [@IanVS](https://github.com/IanVS), [#4161](https://github.com/wailsapp/wails/pull/4161)
- Добавлена возможность регистрировать хуки для пользовательских событий; авторы — [@fbbdev](https://github.com/fbbdev) и [@IanVS](https://github.com/IanVS), [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` для открытия системного файлового менеджера по пути `path` с необязательным выделением через `selectFile`; авторы — [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Новый флаг `-git` для команды `wails3 init`; автор — [@leaanthony](https://github.com/leaanthony)
- Новая команда `wails3 generate webview2bootstrapper`; автор — [@leaanthony](https://github.com/leaanthony)
- В среду выполнения добавлен метод `init()` для её ручной инициализации; автор — [@leaanthony](https://github.com/leaanthony)
- В WindowOptions окна добавлен параметр `WindowDidMoveDebounceMS`; автор — [@leaanthony](https://github.com/leaanthony)
- Добавлена функция единственного экземпляра; автор — [@leaanthony](https://github.com/leaanthony). Основано на [PR для v2](https://github.com/wailsapp/wails/pull/2951) от @APshenkin.
- Команда `wails3 generate template`; автор — [@leaanthony](https://github.com/leaanthony)
- Команда `wails3 releasenotes`; автор — [@leaanthony](https://github.com/leaanthony)
- Команда `wails3 update cli`; автор — [@leaanthony](https://github.com/leaanthony)
- Параметр `-clean` для команды `wails3 generate bindings`; автор — [@leaanthony](https://github.com/leaanthony)
- Добавлена возможность собирать AppImage для Linux на aarch64 (arm64); автор — [@AkshayKalose](https://github.com/AkshayKalose), [#3981](https://github.com/wailsapp/wails/pull/3981)
- Добавлена гиперссылка для спонсора; автор — @ansxuman, [#3958](https://github.com/wailsapp/wails/pull/3958)
- Поддержка сборки пакетов deb, rpm и Arch Linux в Linux от
- Добавлена поддержка универсальных сборок и пакетов для darwin от
- Документация по событиям для веб-сайта от
- Шаблоны для sveltekit и sveltekit-ts, настроенные для разработки без SSR
- Обновление ресурсов сборки с помощью новой команды `wails3 update build-assets` от
- Пример для тестирования HTML Drag and Drop API от
- Поддержка ассоциаций файлов от [leaanthony](https://github.com/leaanthony) в
- Новая команда `wails3 generate runtime` от
- Новый параметр `InitialPosition`, определяющий, следует ли центрировать окно или
- Добавлены методы `Path` и `Paths` в пакет `application` от
- Добавлены параметры Windows `GeneralAutofillEnabled` и `PasswordAutosaveEnabled`
- Добавлена возможность получить окно, вызывающее метод сервиса, от
- Добавлены параметры `EnabledFeatures` и `DisabledFeatures` для Webview2 от
- ⊞ Новая система DIP для улучшенной поддержки мониторов с высоким DPI от
- ⊞ Параметр имени класса окна от [windom](https://github.com/windom/) в
- Сервисы расширены для поддержки функциональности плагинов. Автор —
- 🐧 События WindowDidMove / WindowDidResize в
- ⊞ Событие WindowDidResize в
-  Добавлено событие ApplicationShouldHandleReopen для обработки значка в Dock
-  В реализацию добавлены getPrimaryScreen/getScreens; автор — @tmclane, в
-  Добавлен параметр для отображения панели инструментов в полноэкранном режиме в macOS от
- 🐧 Добавлена логика onKeyPress для преобразования нажатия клавиши в Linux в акселератор
- 🐧 Добавлена задача `run:linux` от
- Экспортирован метод `SetIcon`; автор — [@almas-x](https://github.com/almas-x), в
- Улучшен `OnShutdown`; автор — [@almas-x](https://github.com/almas-x), в
- В интерфейсе `Window` восстановлен метод `ToggleMaximise` от
- В `Environment()` добавлена дополнительная информация. Автор — @leaanthony, в
- Метод `WebviewWindow.IsFocused` предоставлен через интерфейс `Window` от
- Поддержка нескольких разделённых пробелами событий-триггеров в системе WML от
- Добавлены экспорты ESM из встроенного JS-скрипта среды выполнения от
- Добавлен флаг генератора привязок для использования встроенного JS-скрипта среды выполнения вместо
- Реализован `setIcon` в linux; автор — [@abichinger](https://github.com/abichinger)
- В команду dev добавлен флаг `-port` и поддержка переменной окружения
- Добавлены тесты вызовов привязанных методов, автор:
- ⊞ Добавлен `SetIgnoreMouseEvents` для уже созданного окна, автор:
-  Добавлена возможность задавать уровень расположения окна (порядок), автор:

### Исправлено

- Исправлены уменьшенные вдвое `Screen.Bounds`, `WorkArea` и `Size` на компьютерах Mac с дисплеем Retina: значения точек NSScreen преобразуются в физические пиксели в полях `Physical*`, а поля верхнего уровня `Screen.X`/`Y` заполняются, чтобы корректно определять соприкосновение нескольких мониторов и размещать окна в рабочей области; [PR](https://github.com/wailsapp/wails/pull/5168), автор: @wayneforrest
- Исправлена гонка данных в ScreenManager, вызывавшая взаимную блокировку WebKit DisplayLink при изменении конфигурации дисплеев (например, при горячем подключении внешнего монитора во время перехода в спящий режим или выхода из него)
- Если Assets.car существует, CFBundleIconName теперь напрямую получает значение appicon; [PR](https://github.com/wailsapp/wails/pull/5154), автор: @symball
- Исправлена ошибка, из-за которой `wails3 doctor` сообщал о неверных пакетах WebKitGTK в Fedora, openSUSE, Arch и NixOS. Резервные варианты 4.0 удалены, поскольку v3 требует API 4.1 во время компиляции (#5071)
- Исправлено имя пакета webkit2gtk, указываемое командой doctor для openSUSE (`webkit2gtk4_1-devel` → `webkit2gtk3-devel` — правильное имя пакета openSUSE) (#5071)
- Исправлена ошибка `Unexpected token '<'` при отсутствии `/wails/custom.js` в режиме разработки настольного приложения. Добавлены явный обработчик 404 для `/wails/custom.js` и регистронезависимая проверка `Content-Type` в `loadOptionalScript`, чтобы предотвратить внедрение резервных HTML-страниц SPA в качестве JavaScript. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Исправлено состояние выделения меню в области уведомлений macOS: когда меню открыто, значок теперь отображается как выбранный (#4910)
- Исправлено появление привязанного к области уведомлений окна позади других окон в macOS: теперь используется правильный уровень всплывающего окна (#4910)
- Исправлены неверные примеры импорта `@wailsio/runtime` в документации (#4989)
- Исправлена невозможность свернуть безрамочное окно в Darwin (#4294)
- Устранены зависания на 20-30 минут во время `wails3 build` и `wails3 dev`: `node_modules/` исключён из проверки актуальности go-task. Ранее glob-шаблон `sources: "**/*"` заставлял go-task перечислять все файлы в `node_modules/` и вычислять их контрольные суммы (50000-100000+ файлов при использовании тяжёлых зависимостей, таких как MUI), что было особенно медленно в Windows/NTFS (#4939)
- Исправлен сбой сборки GTK4, вызванный конфликтом C-typedef `Screen` с X11 Xlib.h (#4957)
- Унифицировано поведение методов управления меткой на значке приложения в Dock в macOS
- Исправлена ошибка, из-за которой `InvisibleTitleBarHeight` применялся ко всем окнам macOS, а не только к безрамочным окнам и окнам с прозрачной строкой заголовка (#4960)
- Устранены дрожание и рывки окна при изменении размера за верхние углы с включённым `InvisibleTitleBarHeight`: запуск перетаскивания рядом с краями окна теперь пропускается (#4960)
- Исправлена генерация отображённых типов с ключами-перечислениями в привязках JS/TS (#4437), автор: @fbbdev
- Исправлена неработающая функция перетаскивания файлов в Windows при масштабе экрана, отличном от 100%
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом сбрасывании файлов в Windows
- Исправлены координаты сбрасывания файлов в Windows, использовавшие неверное пространство пикселей (физические пиксели и CSS-пиксели)
- Исправлена ненадёжная работа перетаскивания файлов с эффектами наведения в Linux
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом сбрасывании файлов в Linux
- Исправлена ошибка, из-за которой после скрытия и показа окна в Linux/GTK4 оно иногда восстанавливалось свёрнутым: теперь используется `gtk_window_present()` (#4957)
- Исправлено получение и задание положения окна в Linux/GTK4, всегда возвращавшее 0,0: добавлена условная поддержка X11 через `XTranslateCoordinates`/`XMoveWindow` (#4957)
- Исправлена ошибка, из-за которой максимальный размер окна не соблюдался в Linux/GTK4: удалённый `gtk_window_set_geometry_hints` заменён ограничением размера на основе сигналов (#4957)
- Исправлено масштабирование DPI в Linux/GTK4: реализованы правильный расчёт PhysicalBounds и поддержка дробного масштабирования через `gdk_monitor_get_scale` (GTK 4.14+)
- Исправлено дублирование пунктов меню при создании новых окон в Linux/GTK4
- Исправлена генерация отображаемых типов с ключами-перечислениями в привязках JS/TS (#4437), автор: @fbbdev
- Исправлена неработающая функция перетаскивания файлов в Windows при масштабе экрана, отличном от 100%
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом сбрасывании файлов в Windows
- Исправлены координаты сбрасывания файлов в Windows, использовавшие неверное пространство пикселей (физические пиксели и CSS-пиксели)
- Исправлена ненадёжная работа перетаскивания файлов с эффектами наведения в Linux
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом сбрасывании файлов в Linux
- Исправлено масштабирование DPI в Linux/GTK4: реализованы правильный расчёт PhysicalBounds и поддержка дробного масштабирования через `gdk_monitor_get_scale` (GTK 4.14+)
- Исправлено дублирование пунктов меню при создании новых окон в Linux/GTK4
- Исправлена генерация отображаемых типов с ключами-перечислениями в привязках JS/TS (#4437), автор: @fbbdev
- Исправлена проблема с «окнами-призраками» в macOS, вызванная тем, что AppKit API вызывались не из главного потока в App.Window.Current() (#4947), автор: @wimaha
- Исправлена неработающая функция HTML `<input type="file">` в macOS: реализован WKUIDelegate runOpenPanelWithParameters (#4862)
- Исправлена неработающая функция нативного перетаскивания файлов при использовании npm-модуля `@wailsio/runtime` в macOS/Linux (#4953), автор: @leaanthony
- Исправлена генерация привязок для псевдонимов типов из других пакетов (#4578), автор: @fbbdev
- Исправлен сбой OpenFileDialog в Linux из-за нарушения потокобезопасности GTK (#3683), автор: @ddmoney420
- Исправлен сбой SIGSEGV при вызове `Focus()` для скрытого или уничтоженного окна (#4890), автор: @ddmoney420
- Устранена возможная паника при задании пустого значка или растрового изображения в Linux (#4923), автор: @ddmoney420
- Исправлен сбой ErrorDialog при вызове из привязки службы в macOS (#3631), автор: @leaanthony
- Обеспечено отображение меню в ОС Windows в `v3\examples\dialogs`, автор: @ndianabasi
- Исправлено состояние гонки, вызывавшее TypeError при перезагрузке страницы (#4872), автор: @ddmoney420
- Исправлен неверный вывод тестов генератора привязок: из метода `Collector.IsVoidAlias()` удалено глобальное состояние (#4941), автор: @fbbdev
- Исправлено неработающее диалоговое окно выбора файлов `<input type="file">` в macOS (#4862), автор: @leaanthony
- Исправлено использование `Position()` и `SetPosition()` несогласованных систем координат в macOS, из-за которого положение окна смещалось при сохранении и восстановлении состояния (#4816), автор — @leaanthony
- Исправлена ошибка SetProcessDpiAwarenessContext «Доступ запрещён», возникавшая, когда режим учёта DPI уже был задан в манифесте приложения (#4803)
- Обновлена страница документации о сочетаниях клавиш и исправлен тип параметра, принимающего функцию обратного вызова, у `KeyBinding.Add`, автор — @ndianabasi
- Исправлена документация по созданию пользовательских привязок: необходимо использовать `-d String` вместо `-o String`
- Исправлена ошибка, из-за которой при вызове `menu.Update()` не удалялись дочерние элементы меню
- Исправлены устаревшие ссылки на Manager API в документации (31 файл обновлён для использования нового шаблона, например `app.Window.New()`, `app.Event.Emit()` и т. д.), автор — @leaanthony
- Исправлен аварийный сбой в Linux при панике в методах Go, привязанных к JS, вызванный переопределением обработчиков сигналов в WebKit (#3965), автор — @leaanthony
- Исправлена ошибка, из-за которой SaveFileDialog.SetFilename() не действовал в Linux (#4841), автор — @samstanier
- Исправлена ошибка, из-за которой в примере перетаскивания координаты места сброса отображались как undefined
- Исправлена ошибка создания пакета приложения для macOS, возникавшая, когда APP_NAME содержал пробелы (проблема с раскрытием фигурных скобок)
- Исправлена паника из-за выхода индекса за границы массива в Windows при вызове методов служб (отменён переход на goccy/go-json)
- Исправлено неработающее перетаскивание файлов в Windows при масштабе экрана, отличном от 100%
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом сбросе файлов в Windows
- Исправлено вычисление координат сброса файлов в неверном пространстве пикселей в Windows (физические пиксели вместо CSS-пикселей)
- Исправлена ненадёжная работа перетаскивания файлов с эффектами наведения в Linux
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом сбросе файлов в Linux
- Обновлены все команды в файлах Taskfile.yml для всех операционных систем, чтобы корректно обрабатывать пробелы в таких переменных, как `APP_NAME`, автор — @ndianabasi
- Исправлена ошибка аргумента команды при выполнении задачи 'build:universal:lipo:go' в Linux, автор — @wux1an
- Исправлена ошибка Docker «undefined symbol: **<em>ubsan</em>handle_xxxxxxx» при выполнении 'wails3 build GOOS=darwin GOARCH=arm64' в Linux, автор — @wux1an
- Объединена документация по пользовательским протоколам и добавлены разделы об Universal Links, автор — @leaanthony
- Исправлен аварийный сбой меню системного трея Windows при многократном нажатии на значок: добавлена защита от одновременных вызовов TrackPopupMenuEx (#4151), автор — @leaanthony
- Предотвращён аварийный сбой приложения при вызове systray.Run() до app.Run(), автор — @leaanthony
- Исправлен аварийный сбой в macOS при переключении видимости окна через Hide()/Show(), когда включён параметр ApplicationShouldTerminateAfterLastWindowClosed (#4389), автор — @leaanthony
- Исправлена утечка памяти в контекстных меню macOS и Windows при их многократном открытии (#4012), автор — @leaanthony
- Исправлена ошибка, из-за которой нативные ресурсы контекстного меню в macOS не использовались повторно и при каждом показе создавалось новое меню (#4012), автор — @leaanthony
- Исправлена ошибка, из-за которой нажатие на значок в Dock macOS не отображало скрытые окна, если приложение было запущено с `Hidden: true` (#4583), автор — @leaanthony
- Исправлена ошибка, из-за которой диалог печати macOS не открывался вследствие неверного типа указателя на окно в вызове CGO (#4290), автор — @leaanthony
- Исправлен аварийный сбой меню окна в Wayland, вызванный обращением appmenu-gtk-module к окну, для которого ещё не созданы нативные ресурсы (#4769), автор — @leaanthony
- Исправлен аварийный сбой приложения GTK, когда имя приложения содержало недопустимые символы (пробелы, круглые скобки и т. д.), автор — @leaanthony
- Исправлена ошибка «недостаточно памяти» при инициализации перетаскивания в Windows (#4701), автор — @overlordtm
- Исправлена ошибка, из-за которой файловый менеджер открывал неверный каталог в Linux вследствие неправильного экранирования URI (#4397), автор — @leaanthony
- Исправлен сбой сборки AppImage в современных дистрибутивах Linux (Arch, Fedora 39+, Ubuntu 24.04+): добавлено автоматическое обнаружение секций ELF `.relr.dyn` и отключение удаления символов и служебных данных из бинарного файла (#4642), автор — @leaanthony
- Исправлена ошибка, из-за которой `wails doctor` ошибочно сообщал об установленных пакетах webkit в Fedora и системах на основе DNF (#4457), автор — @leaanthony
- Исправлена ошибка, из-за которой стандартный `config.yml` запускал `wails3 dev` со сборкой для рабочей среды, автор — @mbaklor
- Исправлены заглушки служб iOS, вызывавшие сбои сборки из-за импорта несуществующего пакета, автор — @leaanthony
- Исправлено структурированное журналирование в методах debug/info, вызывавшее ошибки «нет директив форматирования», автор — @leaanthony
- Удалены временные отладочные операторы печати, случайно попавшие в код при слиянии изменений для мобильных платформ, автор — @leaanthony
- Исправлен аварийный сбой WebKitGTK в Wayland на графических процессорах NVIDIA (ошибка 71 Protocol error): добавлено автоматическое отключение средства визуализации DMA-BUF, автор — @leaanthony
- Исправлено игнорирование значения альфа-канала в `application.WebviewWindowOptions.BackgroundColour` в Linux ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Исправлена ошибка, из-за которой значком в системном трее Windows по умолчанию не становился значок приложения, если пользовательский значок не был задан (#4704)
- Добавлено отслеживание владельца `HICON`, чтобы уничтожались только созданные пользователем дескрипторы, что предотвращает сбои при перезапуске Explorer (#4653).
- При уничтожении теперь освобождаются обработчик изменений системной темы Windows и сохранённые значки трея, что устраняет утечки горутин и контекстов устройств (#4653).
- Всплывающие подсказки значков трея теперь обрезаются до 127 единиц UTF-16, чтобы избежать повреждения суррогатных пар и многобайтовых глифов (#4653).
- Исправлен сбой задачи упаковки для Windows (#4667)
- Исправлена переменная appicon для Linux AppImage в файле задач Linux, [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Исправлена ошибка сборки для Windows, вызванная изменением сигнатуры в go-webview2 v1.0.22 (#4513, #4645)
- Исправлена переменная appicon для Linux AppImage в файле задач Linux, [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Исправлен перебор протоколов в desktop.tmpl для Linux путём замены `<.Info.Protocol>` на `<.Protocol>`, автор — @Tolfx в #4510
- Исправлена ошибка повторного определения в демонстрации liquid glass в [#4542](https://github.com/wailsapp/wails/pull/4542), автор — @Etesam913
- Исправлено обновление меню системного трея в Linux, [#4604](https://github.com/wailsapp/wails/issues/4604), автор — [@JackDoan](https://github.com/JackDoan)
- Исправлено появление белого окна в Windows при создании скрытого окна, автор — @leaanthony в [#4612](https://github.com/wailsapp/wails/pull/4612)
- Исправлен путь импорта пакета уведомлений в документации, автор — @rxliuli в [#4617](https://github.com/wailsapp/wails/pull/4617)
- Исправлена неработающая функция перетаскивания при использовании npm-пакета @wailsio/runtime (#4489), автор — @leaanthony, #4616
- Windows: устранены мерцание окна при запуске и некорректное отображение скрытых окон в [PR](https://github.com/wailsapp/wails/pull/4600), автор — @leaanthony.
- Исправлены проблемы с размером окна при разворачивании в Wayland (https://github.com/wailsapp/wails/issues/4429), автор — [@samstanier](https://github.com/samstanier)
- Исправлены проблемы с размером окна при разворачивании в Wayland (https://github.com/wailsapp/wails/issues/4429), автор — [@samstanier](https://github.com/samstanier)
- Исправлена ошибка повторного определения в демонстрационном примере liquid glass в [#4542](https://github.com/wailsapp/wails/pull/4542), автор — @Etesam913
- Исправлена проблема, из-за которой AssetServer мог аварийно завершить работу в MacOS, в [#4576](https://github.com/wailsapp/wails/pull/4576), автор — @jghiloni
- Исправлена проблема компиляции при сборке с NextJs. Исправление внесено в [#4585](https://github.com/wailsapp/wails/pull/4585), автор — @rev42
- Исправлены конвейеры ночных выпусков в [#4597](https://github.com/wailsapp/wails/pull/4597), автор — @riadafridishibly
- Исправлена ошибка повторного определения в демонстрационном примере liquid glass в [#4542](https://github.com/wailsapp/wails/pull/4542), автор — @Etesam913
- Исправлена проблема, из-за которой AssetServer мог аварийно завершить работу в MacOS, в [#4576](https://github.com/wailsapp/wails/pull/4576), автор — @jghiloni
- Исправлена проблема компиляции при сборке с NextJs. Исправление внесено в [#4585](https://github.com/wailsapp/wails/pull/4585), автор — @rev42
- Исправлены конвейеры ночных выпусков в [#4597](https://github.com/wailsapp/wails/pull/4597), автор — @riadafridishibly
- Исправлена ошибка повторного определения в демонстрационном примере liquid glass в [#4542](https://github.com/wailsapp/wails/pull/4542), автор — @Etesam913
- Исправлена работа SetBackgroundColour в Windows, автор — @PPTGamer, [PR](https://github.com/wailsapp/wails/pull/4492)
- Документация обновлена с учётом изменений после рефакторинга Manager API, автор — @yulesxoxo, [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Исправлена переменная appicon для файла Linux .desktop в Linux Taskfile, [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Документация обновлена с учётом изменений после рефакторинга Manager API, автор — @yulesxoxo, [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Исправлена ошибка разыменования nil-указателя в Windows, о которой сообщалось в [#4456](https://github.com/wailsapp/wails/issues/4456); исправление внесено @leaanthony в [#4460](https://github.com/wailsapp/wails/pull/4460)
- В macOS WKWebView добавлена поддержка `allowsBackForwardNavigationGestures`, включающая жесты навигации смахиванием двумя пальцами (#1857)
- Исправлена проблема, из-за которой onClick не работал для пунктов меню, изначально заданных как отключённые, автор — @leaanthony, [PR #4469](https://github.com/wailsapp/wails/pull/4469). Благодарим @IanVS за первоначальное исследование.
- Исправлена проблема, из-за которой сервер Vite не завершал работу после сбоя сборки (#4403)
- Устранена паника при закрытии или отмене `SaveFileDialog` в windows. Исправление внесено в [PR](https://github.com/wailsapp/wails/pull/4284), автор — @hkhere
- Исправлено перетаскивание на уровне HTML в Windows, автор — [@mbaklor](https://github.com/mbaklor), в [#4259](https://github.com/wailsapp/wails/pull/4259)
- В macOS WKWebView добавлена поддержка `allowsBackForwardNavigationGestures`, включающая жесты навигации смахиванием двумя пальцами (#1857)
- Исправлена проблема, из-за которой onClick не работал для пунктов меню, изначально заданных как отключённые, автор — @leaanthony, [PR #4469](https://github.com/wailsapp/wails/pull/4469). Благодарим @IanVS за первоначальное исследование.
- Исправлена проблема, из-за которой сервер Vite не завершал работу после сбоя сборки (#4403)
- Исправлен разбор уведомлений в Windows, автор — @popaprozac, [PR](https://github.com/wailsapp/wails/pull/4450)
- Команда doctor исправлена для проверки зависимостей Windows SDK, автор — [@kodumulo](https://github.com/kodumulo), в [#4390](https://github.com/wailsapp/wails/issues/4390)
- Исправлено разыменование nil-указателя в processURLRequest на Mac, автор — [@etesam913](https://github.com/etesam913), в [#4366](https://github.com/wailsapp/wails/pull/4366)
- Исправлена ошибка в linux, препятствовавшая работе диалоговых окон с фильтрами, автор — [@bh90210](https://github.com/bh90210), в [#4287](https://github.com/wailsapp/wails/pull/4287)
- Исправлены проблемы меню «Правка» в Windows и Linux, автор — [@leaanthony](https://github.com/leaanthony), в [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Минимальная версия системы в файлах macOS .plist обновлена с 10.13.0 до 10.15.0, автор — [@AkshayKalose](https://github.com/AkshayKalose), в [#3981](https://github.com/wailsapp/wails/pull/3981)
- Исправлена проблема с пропуском идентификаторов окон, автор — [@leaanthony](https://github.com/leaanthony)
- Исправлена проблема с nil-меню при вызове RegisterContextMenu, автор — [@leaanthony](https://github.com/leaanthony)
- Исправлены циклические зависимости в выходных данных генератора привязок, автор — [@fbbdev](https://github.com/fbbdev), в [#4001](https://github.com/wailsapp/wails/pull/4001)
- Исправлены ошибки использования до объявления в выходных данных генератора привязок, автор — [@fbbdev](https://github.com/fbbdev), в [#4001](https://github.com/wailsapp/wails/pull/4001)
- В генератор привязок передаются флаги сборки, автор — [@fbbdev](https://github.com/fbbdev), в [#4023](https://github.com/wailsapp/wails/pull/4023)
- В путях в Windows Taskfile обратные косые черты заменены на прямые, чтобы обеспечить работу на платформах, отличных от Windows, автор — [@leaanthony](https://github.com/leaanthony)
- Исправлены события Mac и Mac JS, автор — [@leaanthony](https://github.com/leaanthony)
- Исправлена взаимная блокировка событий в macOS, автор — [@leaanthony](https://github.com/leaanthony)
- Исправлена ошибка `Parameter incorrect` при инициализации Window в Windows, когда предоставлен HTML, но отсутствует JS, автор — [@leaanthony](https://github.com/leaanthony)
- Исправлен размер префикса ответа, используемого сервером ресурсов для определения типа содержимого, автор — [@fbbdev](https://github.com/fbbdev), в [#4049](https://github.com/wailsapp/wails/pull/4049)
- Исправлена обработка ответов, отличных от 404, для корневого пути индекса на сервере ресурсов, автор — [@fbbdev](https://github.com/fbbdev), в [#4049](https://github.com/wailsapp/wails/pull/4049)
- Исправлено неопределённое поведение генератора привязок при проверке свойств обобщённых типов, автор — [@fbbdev](https://github.com/fbbdev), в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлены выходные данные генератора привязок для моделей, у которых базовый тип имеет свойства, отличные от свойств именованной обёртки, автор — [@fbbdev](https://github.com/fbbdev), в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлены выходные данные генератора привязок для типов ключей отображений и предварительная обработка, автор — [@fbbdev](https://github.com/fbbdev), в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлены выходные данные генератора привязок для структур, реализующих интерфейсы маршалинга, автор — [@fbbdev](https://github.com/fbbdev), в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлено обнаружение циклов типов с участием обобщённых типов в генераторе привязок, автор — [@fbbdev](https://github.com/fbbdev), в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлены недопустимые ссылки на неэкспортируемые модели в выходных данных генератора привязок — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Внедряемый код перемещён в конец файлов сервисов — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлена обработка ошибок операций закрытия файлов в генераторе привязок — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Отключены предупреждения для сервисов, которые определяют методы жизненного цикла или HTTP, но не имеют других привязанных методов, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлена ошибка, из-за которой шаблоны без React не отображали нижний колонтитул Hello World при использовании светлой системной цветовой схемы, — [@marcus-crane](https://github.com/marcus-crane) в [#4056](https://github.com/wailsapp/wails/pull/4056)
- Исправлены скрытые пункты меню в macOS — [@leaanthony](https://github.com/leaanthony)
- Исправлены обработка и форматирование ошибок в обработчиках сообщений — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Исправлен пропуск завершения работы сервиса при выходе из приложения — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Обновления меню теперь выполняются в главном потоке — [@leaanthony](https://github.com/leaanthony)
- Механизм перетаскивания и изменения размера стал надёжнее и теперь точнее соответствует ожидаемому поведению платформы — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Исправлена ошибка [#4097](https://github.com/wailsapp/wails/issues/4097), из-за которой Webpack/Angular отбрасывал код инициализации среды выполнения, — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Исправлены пункты меню, изначально находящиеся в скрытом состоянии, — [@IanVS](https://github.com/IanVS) в [#4116](https://github.com/wailsapp/wails/pull/4116)
- Исправлена ошибка, из-за которой assetFileServer не отдавал файлы `.html` при запросе без расширения, когда `[request]` не существует, а `[request].html` существует
- Исправлены пути генерации значков — [@robin-samuel](https://github.com/robin-samuel) в [#4125](https://github.com/wailsapp/wails/pull/4125)
- Исправлена ошибка, из-за которой события `fullscreen`, `unfullscreen`, `unminimise` и `unmaximise` не генерировались, — [@oSethoum](https://github.com/osethoum) в [#4130](https://github.com/wailsapp/wails/pull/4130)
- Исправлена ошибка NSIS, возникавшая из-за неверного префикса версии по умолчанию в конфигурации, — [@robin-samuel](https://github.com/robin-samuel) в [#4126](https://github.com/wailsapp/wails/pull/4126)
- Исправлена функция среды выполнения Dialogs, возвращавшая экранированные пути в Windows, — [TheGB0077](https://github.com/TheGB0077) в [#4188](https://github.com/wailsapp/wails/pull/4188)
- Исправлен путь обнаружения WebView2 в HKCU — [@leaanthony](https://github.com/leaanthony).
- Исправлена проблема с вводом в macOS — [@leaanthony](https://github.com/leaanthony).
- Исправлено имя файла задачи генерации значков Windows — [@yulesxoxo](https://github.com/yulesxoxo) в [#4219](https://github.com/wailsapp/wails/pull/4219).
- Исправлена проблема с прозрачностью безрамочных окон — [@leaanthony](https://github.com/leaanthony) на основе работы @kron.
- Исправлены вызовы установки фокуса, когда окно отключено или свёрнуто, — [@leaanthony](https://github.com/leaanthony) на основе работы @kron.
- Исправлена ошибка, из-за которой значки в системной области не появлялись после перезапуска панели задач, — [@leaanthony](https://github.com/leaanthony) на основе работы @kron.
- Исправлена отсутствующая реализация Flush() в fallbackResponseWriter в [#4245](https://github.com/wailsapp/wails/pull/4245)
- Исправлена отсутствующая реализация Flush() в fallbackResponseWriter — [@superDingda] в [#4236](https://github.com/wailsapp/wails/issues/4236)
- Исправлены сбои при закрытии окна macOS, когда ожидалось завершение асинхронного вызова функции, привязанной к Go, — [@joshhardy](https://github.com/joshhardy) в [#4354](https://github.com/wailsapp/wails/pull/4354)
- Исправлено состояние гонки при запуске режима эффективности Windows — [@leaanthony](https://github.com/leaanthony)
- Исправлено освобождение дескриптора значка Windows — [@leaanthony](https://github.com/leaanthony).
- Исправлено `OpenFileManager` в Windows — [@PPTGamer](https://github.com/PPTGamer) в [#4375](https://github.com/wailsapp/wails/pull/4375).
- Исправлены параметры минимальной и максимальной ширины для Linux — @atterpac в [#3979](https://github.com/wailsapp/wails/pull/3979)
- Определения типов для шаблонов TypeScript обновлены посредством повышения версии npm — @atterpac в [#3966](https://github.com/wailsapp/wails/pull/3966)
- Исправлена ссылка на CSS в шаблоне SvelteKit — @atterpac в [#3945](https://github.com/wailsapp/wails/pull/3945)
- Ключевые функции обратного вызова в run() окна теперь гарантированно вызываются в главном потоке — [@leaanthony](https://github.com/leaanthony)
- Исправлены примеры выбора каталога в диалоговом окне — [@leaanthony](https://github.com/leaanthony)
- Создана новая страница ошибки на китайском языке, отображаемая при отсутствии index.html, — [@leaanthony](https://github.com/leaanthony)
-  Функция обратного вызова `windowDidBecomeKey` теперь гарантированно выполняется в главном потоке — [@leaanthony](https://github.com/leaanthony)
-  Добавлена поддержка полноэкранного режима для безрамочных окон — [@leaanthony](https://github.com/leaanthony)
-  Улучшена логика уничтожения окон — [@leaanthony](https://github.com/leaanthony)
-  Исправлена логика позиционирования окон, прикреплённых к значкам в системной области, — [@leaanthony](https://github.com/leaanthony)
-  Добавлена поддержка полноэкранного режима для безрамочных окон — [@leaanthony](https://github.com/leaanthony)
- Исправлена обработка событий — [@leaanthony](https://github.com/leaanthony)
- Исправлена логика завершения работы окна — [@leaanthony](https://github.com/leaanthony)
- Общий taskfile теперь по умолчанию генерирует привязки TypeScript для шаблонов TypeScript — [@leaanthony](https://github.com/leaanthony)
- Исправлено закрытие приложения при получении сообщения WM_CLOSE, когда нет открытых окон и используется только значок в системной области, — [@mmalcek](https://github.com/mmalcek) в [#3990](https://github.com/wailsapp/wails/pull/3990)
- Исправлена сборка с garble — @5aaee9 в [#3192](https://github.com/wailsapp/wails/pull/3192)
- Исправлены сборки NSIS для Windows — [@leaanthony](https://github.com/leaanthony)
- Исправлена взаимная блокировка в диалоговом окне Linux при выборе нескольких элементов, вызванная незакрытым
- Исправлена кроссплатформенная очистка файлов .syso во время сборки для Windows —
- Исправлена компиляция AppImage для amd64 — @atterpac в
- Исправлено обновление ресурсов сборки — @ansxuman в
- Исправлена реализация `OnClick` и `OnRightClick` для области уведомлений Linux; автор — @atterpac
- Исправлена неработоспособность `AlwaysOnTop` на Mac; автор —
-  Исправлено дублирование в `application.NewEditMenu`
- 🐧 Исправлена компиляция для aarch64
- ⊞ Исправлены пункты меню с переключателями; автор —
- Исправлена ошибка при сборке запускаемого пакета .app на MacOS, когда 'name' и 'outputfilename'
- Исправлена ошибка использования customEventProcessor в примере drag-n-drop; автор —
- 🐧 Исправлена ошибка компиляции в Linux, появившаяся после добавления IgnoreMouseEvents; автор —
- ⊞ Исправлена ошибка создания файла значка syso; автор —
- 🐧 Интегрировано исправление для нативного запуска в Wayland из
- Не привязывать внутренние методы сервисов в
- ⊞ Исправлена паника при запуске области уведомлений в
- Не привязывать внутренние методы сервисов в
- ⊞ Исправлена паника при запуске области уведомлений в
- Значительно переработаны пункты меню и обработка событий. Пока что в основном улучшена macOS. Автор —
- Исправлены тесты после переработки плагинов и событий в
- ⊞ Исправлено предупреждение `Failed to unregister class Chrome_WidgetWin_0`. Автор —
- Проблемы с модулями
- Исправлена передача сообщений о событии изменения размера; автор — [atterpac](https://github.com/atterpac), в
- 🐧 Исправлена ошибка обработки темы в NixOS; автор —
- Исправлена установка проекта между разными томами в Windows; автор —
- Исправлены стили CSS шаблона React, чтобы отображался нижний колонтитул; автор —
- Устранены процессы-зомби при работе в режиме разработки благодаря обновлению refresh до последней версии
- Исправлен поиск файла WebKit для AppImage; автор — [Atterpac](https://github.com/atterpac)
- Исправлена проверка пакетов apt в Doctor; автор — [Atterpac](https://github.com/Atterpac), в
- Исправлено зависание приложения при завершении работы (Darwin); автор — @5aaee9, в
- Исправлены цвета фона примеров в Windows; автор —
- Исправлены контекстные меню по умолчанию; автор — [mmghv](https://github.com/mmghv), в
- Исправлены шестнадцатеричные значения клавиш со стрелками в Darwin; автор —
- В Windows реализована работа drag-n-drop. Добавил
- Исправлена ошибка в Doctor для Linux, возникавшая при отсутствии у пользователя подходящих драйверов
- Исправлено масштабирование DPI при запуске (Windows). Изменил [@almas-x](https://github.com/almas-x), в
- Исправлена строка замены в `go.mod`: теперь используются относительные пути. Это исправляет пути Windows с
- Исправлена обработка щелчков по значку в области уведомлений MacOS при отсутствии связанного окна; автор —
- Исправлен сбой сборки для Windows из-за неизвестного параметра; автор —
- Исправлен сбой в Windows при щелчке левой кнопкой по значку в области уведомлений, когда отсутствует
- Исправлен неверный baseURL при повторном открытии окна; автор — @5aaee9, в PR
- Исправлен порядок ветвей if в методе `WebviewWindow.Restore`; автор —
- Реализовано правильное вычисление `startURL` при нескольких вызовах `GetStartURL`, когда
- Тип JS структуры `Screen` приведён в соответствие с её аналогом в Go; автор —
- Исправлен метод `WML.Reload`, чтобы обеспечить надлежащую очистку зарегистрированного события
- Исправлено немедленное закрытие пользовательского контекстного меню в Linux; автор —
- Исправлены выходной путь и расширение файлов моделей, создаваемых генератором привязок
- Исправлены пути импорта файлов моделей в JS-коде, создаваемом генератором привязок
- Исправлена работа drag-n-drop в некоторых дистрибутивах Linux; автор —
- Добавлена отсутствовавшая задача для macOS при использовании `wails3 task dev`; автор —
- Исправлена ошибка присваивания в nil-карту при регистрации событий; автор —
- Исправлена десериализация параметров привязанных методов; автор —
- Исправлена обработка нескольких возвращаемых значений привязанных методов; автор —
- Исправлено обнаружение в Doctor пакета npm, установленного не через системный менеджер пакетов
- Добавлен отсутствовавший MicrosoftEdgeWebview2Setup.exe. Благодарим
- Исправлен случайный сбой в Linux из-за обработки идентификатора окна; автор — @leaanthony. На основе
- Исправлен сбой systemTray.setIcon в Linux; автор —
- Исправлено применение рамки окна при первом вызове функции `setFrameless` в

### Изменено

- **ИЗМЕНЕНИЕ, НАРУШАЮЩЕЕ ОБРАТНУЮ СОВМЕСТИМОСТЬ**: ключи карт в сгенерированных привязках JS/TS теперь помечаются как необязательные, чтобы точно отражать семантику карт Go. Теперь доступ к значению карты в TypeScript возвращает `T | undefined` вместо `T`, поэтому требуются проверки на null или утверждения типов (#4943); автор — `@fbbdev`
- Использование `Event` заменено на `Events` в соответствии с изменениями в `@wailsio/runtime`; соответствующие вызовы функций в документации обновлены в `Features/Events/Event System`; автор — @AbdelhadiSeddar
- `EnabledFeatures`, `DisabledFeatures` и `AdditionalBrowserArgs` перенесены из параметров отдельных окон в `Options.Windows` уровня приложения (#4559); автор — @leaanthony
- Обновлён README для примера `Drag N Drop`; подчёркнуто, что пример демонстрирует `Internal Drag and Drop`; автор — @ndianabasi
- Уровень различных отладочных сообщений изменён с Info на Debug (автор — @mbaklor)
- **ИЗМЕНЕНИЕ, НАРУШАЮЩЕЕ ОБРАТНУЮ СОВМЕСТИМОСТЬ:** `EnableDragAndDrop` переименован в `EnableFileDrop` в параметрах окна
- **ИЗМЕНЕНИЕ, НАРУШАЮЩЕЕ ОБРАТНУЮ СОВМЕСТИМОСТЬ:** `DropZoneDetails` переименован в `DropTargetDetails` в контексте события
- **ИЗМЕНЕНИЕ, НАРУШАЮЩЕЕ ОБРАТНУЮ СОВМЕСТИМОСТЬ:** метод `DropZoneDetails()` объекта `WindowEventContext` переименован в `DropTargetDetails()`
- **ИЗМЕНЕНИЕ, НАРУШАЮЩЕЕ ОБРАТНУЮ СОВМЕСТИМОСТЬ:** событие `WindowDropZoneFilesDropped` удалено; вместо него используйте `WindowFilesDropped`
- **ИЗМЕНЕНИЕ, НАРУШАЮЩЕЕ ОБРАТНУЮ СОВМЕСТИМОСТЬ:** HTML-атрибут `data-wails-dropzone` заменён на `data-file-drop-target`
- **ИЗМЕНЕНИЕ, НАРУШАЮЩЕЕ ОБРАТНУЮ СОВМЕСТИМОСТЬ:** CSS-класс наведения `wails-dropzone-hover` заменён на `file-drop-target-active`
- **ИЗМЕНЕНИЕ, НАРУШАЮЩЕЕ ОБРАТНУЮ СОВМЕСТИМОСТЬ:** из Windows удалены параметры `DragEffect`, `OnEnterEffect` и `OnOverEffect` (они относились к удалённому IDropTarget)
- Для всей обработки JSON во время выполнения (привязок методов, событий, запросов WebView, уведомлений и kvstore) теперь используется goccy/go-json, что повышает производительность на 21-63% и сокращает количество выделений памяти на 40-60%
- Оптимизирована компоновка структуры BoundMethod и кэшируется флаг isVariadic, чтобы сократить накладные расходы при каждом вызове
- Для методов с `<=8` аргументами используется выделяемый в стеке буфер аргументов, чтобы избежать выделений памяти в куче
- Оптимизирован сбор результатов вызовов методов, чтобы не выделять срез для единственного возвращаемого значения
- Для кэша MIME-типов используется sync.Map, что повышает производительность при конкурентном доступе
- Для чтения тела запросов HTTP-транспорта используется пул буферов
- Канал CloseNotify в определителе типа содержимого теперь выделяется отложенно, чтобы сократить количество выделений памяти на каждый запрос
- С сервера ресурсов удалено отладочное журналирование CSS
- Карта расширений MIME-типов дополнена поддержкой более 50 распространённых веб-форматов (шрифтов, аудио, видео и т. д.)
- Обновлена документация по параметрам Window `X/Y` — @ruhuang2001
- В документацию по `Frontend Runtime` добавлены дополнительные параметры генерации привязок для фронтенда — @ndianabasi
- Обновлена страница документации по серверу ресурсов Wails v3 — @ndianabasi
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ**: удалены функции диалогов уровня пакета (`application.InfoDialog()`, `application.QuestionDialog()` и т. д.). Вместо них используйте менеджер `app.Dialog`: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Документация по диалогам приведена в соответствие с фактическим API: используются `app.Dialog.*`, `AddButton()` с функциями обратного вызова (не `SetButtons()`), `SetDefaultButton(*Button)` (не строка), `AddFilter()` (не `SetFilters()`), `SetFilename()` (не `SetDefaultFilename()`) и `app.Dialog.OpenFile().CanChooseDirectories(true)` для выбора папки
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ**: теперь по умолчанию создаются производственные сборки. Чтобы создавать сборки для разработки, укажите `DEV=true` в файлах Taskfile. Для примеров создайте новый проект — @leaanthony
- При отправке пользовательского события без аргументов данных или с одним таким аргументом значение данных присваивается непосредственно полю Data без оборачивания в срез — [@fbbdev](https://github.com/fbbdev) в [#4633](https://github.com/wailsapp/wails/pull/4633)
- Значки приложений в области уведомлений Windows теперь учитывают `SystemTray.Show()`/`Hide()`, переключая `NIS_HIDDEN`, благодаря чему приложения могут полностью исчезать и возвращаться (#4653).
- При регистрации значка в области уведомлений повторно используются уже полученные значки, `NOTIFYICON_VERSION_4` задаётся один раз, а `NIF_SHOWTIP` включается, чтобы всплывающие подсказки восстанавливались после перезапуска Explorer (#4653).
- macOS: для центрирования окна используется `visibleFrame` вместо `frame`, чтобы исключить области строки меню и Dock
- macOS: для центрирования окна используется `visibleFrame` вместо `frame`, чтобы исключить области строки меню и Dock
- При запуске `wails3 update build-assets` с параметром `-config` значения, заданные через параметры `-product*`,
- `window.NativeWindowHandle()` → `window.NativeWindow()` — @leaanthony в [#4471](https://github.com/wailsapp/wails/pull/4471)
- Переработана внутренняя обработка окон — @leaanthony в [#4471](https://github.com/wailsapp/wails/pull/4471)
- Удалены `application.WindowIDKey` и `application.WindowNameKey` (заменены на `application.WindowKey`) — [@leaanthony](https://github.com/leaanthony)
- ContextMenuData теперь возвращает строку вместо any — [@leaanthony](https://github.com/leaanthony)
- В привязках JS/TS поля классов с типами массивов фиксированной длины теперь инициализируются с ожидаемой длиной, а не остаются пустыми — [@fbbdev](https://github.com/fbbdev) в [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData теперь возвращает строку вместо any — [@leaanthony](https://github.com/leaanthony)
- `application.NewService` больше не принимает параметры в качестве необязательного аргумента (вместо этого используйте `application.NewServiceWithOptions`) — [@leaanthony](https://github.com/leaanthony) в [#4024](https://github.com/wailsapp/wails/pull/4024)
- Удалена зависимость `nanoid` — [@leaanthony](https://github.com/leaanthony)
- Обновлён пример Window для стилей окон mica/acrylic/tabbed — [@leaanthony](https://github.com/leaanthony)
- Из привязок JS/TS удалены файлы моделей `internal.js/ts`; теперь все модели находятся в `models.js/ts` — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- В привязках JS/TS именованные типы больше никогда не отображаются как псевдонимы других именованных типов; прежнее поведение теперь применяется только к псевдонимам — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- В привязках JS/TS в режиме классов поля структур, типом которых является параметр типа, помечаются как необязательные и никогда не инициализируются автоматически — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- ESLint удалён из шаблонов — [@IanVS](https://github.com/IanVS) в [#4059](https://github.com/wailsapp/wails/pull/4059)
- Год в уведомлении об авторских правах обновлён до 2025 — [@IanVS](https://github.com/IanVS) в [#4037](https://github.com/wailsapp/wails/pull/4037)
- Добавлена документация по event.Sender — [@IanVS](https://github.com/IanVS) в [#4075](https://github.com/wailsapp/wails/pull/4075)
- Поддержка Go 1.24 — [@leaanthony](https://github.com/leaanthony)
- Обработчики `ServiceStartup` теперь вызываются при вызове `App.Run`, а не в `application.New` — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ошибки `ServiceStartup` теперь возвращаются из `App.Run` вместо завершения процесса — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Вызовы привязок и диалогов из JS теперь отклоняются с объектами ошибок вместо строк — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Улучшено позиционирование меню области уведомлений в Windows — [@leaanthony](https://github.com/leaanthony)
- Среда выполнения JS перенесена на TypeScript — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Среда выполнения инициализируется сразу после импорта, поэтому ждать загрузки окна не требуется; автор: [@fbbdev](https://github.com/fbbdev), см. [#4100](https://github.com/wailsapp/wails/pull/4100)
- Среда выполнения больше не экспортирует метод init. Для её инициализации можно использовать импорт только ради побочных эффектов; автор: [@fbbdev](https://github.com/fbbdev), см. [#4100](https://github.com/wailsapp/wails/pull/4100)
- Связанные методы теперь возвращают `CancellablePromise`, который при отмене отклоняется с `CancelError`. Фактический результат вызова отбрасывается; автор: [@fbbdev](https://github.com/fbbdev), см. [#4100](https://github.com/wailsapp/wails/pull/4100)
- Типы встроенных сервисов теперь единообразно называются `Service`; автор: [@fbbdev](https://github.com/fbbdev), см. [#4067](https://github.com/wailsapp/wails/pull/4067)
- Функции создания встроенных сервисов с параметрами теперь единообразно называются `NewWithConfig`; автор: [@fbbdev](https://github.com/fbbdev), см. [#4067](https://github.com/wailsapp/wails/pull/4067)
- Метод `Select` сервиса `sqlite` теперь называется `Query` для согласованности с API Go; автор: [@fbbdev](https://github.com/fbbdev), см. [#4067](https://github.com/wailsapp/wails/pull/4067)
- Шаблоны: среда выполнения перенесена в "dependencies", файлы package.json упорядочены; автор: [@IanVS](https://github.com/IanVS), см. [#4133](https://github.com/wailsapp/wails/pull/4133)
- В режиме разработки создаёт пакеты приложений и выполняет их ad-hoc-подпись, чтобы сделать доступными некоторые API macOS; автор: [@popaprozac](https://github.com/popaprozac), см. [#4171](https://github.com/wailsapp/wails/pull/4171)
- Ресурсы сборки перенесены в каталоги для соответствующих платформ; автор: [@leaanthony](https://github.com/leaanthony)
- Файлы Taskfile перемещены, переименованы и помещены в каталоги для соответствующих платформ; автор: [@leaanthony](https://github.com/leaanthony)
- Значительно улучшено взаимодействие при отсутствии `index.html`; автор: [@leaanthony](https://github.com/leaanthony)
- [Windows] Повышена производительность сворачивания и восстановления окон; автор: [@leaanthony](https://github.com/leaanthony). На основе исходного [PR](https://github.com/wailsapp/wails/pull/3955) от [562589540](https://github.com/562589540)
- Удалён параметр `ShouldClose` (вместо него зарегистрируйте обработчик для events.Common.WindowClosing); автор: [@leaanthony](https://github.com/leaanthony)
- [Windows] Уменьшено мерцание при открытии окна; автор: [@leaanthony](https://github.com/leaanthony)
- Удалена `Window.Destroy`, поскольку эта функция предназначалась для внутреннего использования; автор: [@leaanthony](https://github.com/leaanthony)
- События `WindowClose` переименованы в `WindowClosing`; автор: [@leaanthony](https://github.com/leaanthony)
- Сборки фронтенда теперь используют окружение vite "development" или "production" в зависимости от типа сборки; автор: [@leaanthony](https://github.com/leaanthony)
- Обновление до go-webview2 v1.19; автор: [@leaanthony](https://github.com/leaanthony)
- Обеспечено использование форка taskfile; автор: @leaanthony
- Форк Taskfile обновлён для устранения проблем с версией при установке с помощью
- Использование форка Taskfile для устранения проблем с версией при установке с помощью
- `service.OnStartup` теперь завершает работу приложения при ошибке и запускает
- Обмен сообщениями о нажатиях на значок в области уведомлений переработан для лучшего соответствия действиям пользователя; автор:
- Во встраиваемые ресурсы добавлен `all:frontend/dist` для поддержки фреймворков, которые создают
- Рефакторинг Taskfile; автор: [leaanthony](https://github.com/leaanthony), см.
- Обновление до `go-webview2` v1.0.16; автор:
- Тип `Screen` исправлен: теперь он включает `ID`, а не `Id`; автор:
- Версия Wails в `go.mod.tmpl` обновлена для поддержки `application.ServiceOptions`; автор:
- Исправлено определение имени сервиса; автор: [windom](https://github.com/windom/), см.
- mkdocs serve теперь использует Docker; автор: [leaanthony](https://github.com/leaanthony)
- Конфигурация разработки объединена в `config.yml`; автор:
- Диалоговое окно значка в области уведомлений теперь по умолчанию использует значок приложения, если он доступен (Windows); автор:
- Улучшено отображение сведений о GPU и памяти в macOS; автор:
- Удалены `WebviewGpuIsDisabled` и `EnableFraudulentWebsiteWarnings`
- Изменение API событий: `On`/`Emit` -> пользовательские события, `OnApplicationEvent` ->
- Исправление API событий в Linux; автор: [TheGB0077](https://github.com/TheGB0077), см.
- [CI] улучшены рабочие процессы и добавлена возможность запускать их также в форках и
- `AbsolutePosition()` переименован в `Position()`; автор:
- Зависимость Linux WebKit обновлена с webkitgtk2-4.0 до webkit2gtk-4.1, чтобы
- Встроенный скрипт среды выполнения JS теперь является модулем ESM: теги script, импортирующие его
- Пакет `@wailsio/runtime` не публикует свой API в `window.wails`
- Модуль API окон `@wailsio/runtime/src/window` теперь предоставляет содержащий его
- API окон JS обновлён в соответствии с текущим `WebviewWindow` Go
- Генератор привязок теперь по умолчанию использует вызовы по идентификатору. Параметр CLI `-id`
- Новая структура кода привязок: ранее выходные файлы были организованы по каталогам
- Поле структуры `application.Options.Bind` переименовано в
- Новый синтаксис привязки сервисов: экземпляры сервисов теперь необходимо оборачивать в
- Отключение индикатора выполнения в нетерминальной среде или среде CI; автор:

### Удалено

- **НЕСОВМЕСТИМОЕ ИЗМЕНЕНИЕ**: из параметров `WindowsWindow` для отдельных окон удалены `EnabledFeatures`, `DisabledFeatures` и `AdditionalLaunchArgs`. Вместо них используйте параметры уровня приложения `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures` и `Options.Windows.AdditionalBrowserArgs`. Эти флаги глобально применяются к общей среде WebView2 (#4559); автор: @leaanthony
- Удалена нативная реализация `IDropTarget` в Windows в пользу подхода на основе JavaScript (как в v2)
- Удалена зависимость github.com/wailsapp/mimetype в пользу расширенной карты расширений и http.DetectContentType из стандартной библиотеки, что уменьшает размер двоичного файла примерно на 1.2 МБ
- Удалена зависимость gopkg.in/ini.v1: для файлового менеджера Linux реализован минимальный синтаксический анализатор файлов .desktop, что экономит около 45 КБ
- Удалить samber/lo из кода среды выполнения, используя пакет slices стандартной библиотеки Go 1.21+ и минимальные внутренние вспомогательные функции, что сократит размер примерно на 310 КБ
- Удалить отладочные вызовы printf из обработчика схем URL для Darwin (#4834)
- **ИЗМЕНЕНИЕ, НАРУШАЮЩЕЕ ОБРАТНУЮ СОВМЕСТИМОСТЬ**: удалено событие `linux:WindowLoadChanged`; вместо него используйте `linux:WindowLoadFinished`, чтобы определять завершение загрузки WebView (#3896), автор — @leaanthony

### Изменения, нарушающие обратную совместимость

- **Рефакторинг API менеджеров**: плоская структура API приложения преобразована в упорядоченный набор менеджеров для лучшей организации кода и упрощения поиска возможностей; автор — [@leaanthony](https://github.com/leaanthony), изменение внесено в [#4359](https://github.com/wailsapp/wails/pull/4359)
- `app.NewWebviewWindow()` → `app.Window.New()`
- `app.CurrentWindow()` → `app.Window.Current()`
- `app.GetAllWindows()` → `app.Window.GetAll()`
- `app.WindowByName()` → `app.Window.GetByName()`
- `app.EmitEvent()` → `app.Event.Emit()`
- `app.OnApplicationEvent()` → `app.Event.OnApplicationEvent()`
- `app.OnWindowEvent()` → `app.Event.OnWindowEvent()`
- `app.SetApplicationMenu()` → `app.Menu.SetApplicationMenu()`
- `app.OpenFileDialog()` → `app.Dialog.OpenFile()`
- `app.SaveFileDialog()` → `app.Dialog.SaveFile()`
- `app.MessageDialog()` → `app.Dialog.Message()`
- `app.InfoDialog()` → `app.Dialog.Info()`
- `app.WarningDialog()` → `app.Dialog.Warning()`
- `app.ErrorDialog()` → `app.Dialog.Error()`
- `app.QuestionDialog()` → `app.Dialog.Question()`
- `app.NewSystemTray()` → `app.SystemTray.New()`
- `app.GetSystemTray()` → `app.SystemTray.Get()`
- `app.ShowContextMenu()` → `app.ContextMenu.Show()`
- `app.RegisterKeybinding()` → `app.KeyBinding.Register()`
- `app.UnregisterKeybinding()` → `app.KeyBinding.Unregister()`
- `app.GetPrimaryScreen()` → `app.Screen.GetPrimary()`
- `app.GetAllScreens()` → `app.Screen.GetAll()`
- `app.BrowserOpenURL()` → `app.Browser.OpenURL()`
- `app.Environment()` → `app.Env.GetAll()`
- `app.ClipboardGetText()` → `app.Clipboard.Text()`
- `app.ClipboardSetText()` → `app.Clipboard.SetText()`
- Методы Service переименованы: `Name` -> `ServiceName`, `OnStartup` -> `ServiceStartup`, `OnShutdown` -> `ServiceShutdown`; автор — [@leaanthony](https://github.com/leaanthony)
- Методы `Path` и `Paths` перенесены в пакет `application`; автор — [@leaanthony](https://github.com/leaanthony)
- Меню приложения теперь доступно только в macOS; автор — [@leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.78 - 2026-04-21

## Добавлено

## Исправлено

## v3.0.0-alpha.77 - 2026-04-18

## Исправлено

## v3.0.0-alpha.76 - 2026-04-17

## Исправлено

## v3.0.0-alpha.75 - 2026-04-16

## Исправлено

## v3.0.0-alpha.74 - 2026-03-01

## Добавлено

## Исправлено

## v3.0.0-alpha.73 - 2026-02-27

## Исправлено

## v3.0.0-alpha.72 - 2026-02-16

## Исправлено

## v3.0.0-alpha.71 - 2026-02-10

## Добавлено

## Исправлено

## v3.0.0-alpha.70 - 2026-02-09

## Добавлено

## Исправлено

## v3.0.0-alpha.69 - 2026-02-08

## Добавлено

## Исправлено

## v3.0.0-alpha.68 - 2026-02-07

## Добавлено

## Изменено

## Исправлено

## v3.0.0-alpha.67 - 2026-02-04

## Добавлено

## Изменено

## Исправлено

## v3.0.0-alpha.66 - 2026-02-03

## Добавлено

## Изменено

## Исправлено

## Удалено

## v3.0.0-alpha.65 - 2026-02-01

## Добавлено

## v3.0.0-alpha.64 - 2026-01-26

## Добавлено

## v3.0.0-alpha.63 - 2026-01-25

## Исправлено

## v3.0.0-alpha.62 - 2026-01-22

## Исправлено

## v3.0.0-alpha.61 - 2026-01-20

## Исправлено

## v3.0.0-alpha.60 - 2026-01-14

## Исправлено

## v3.0.0-alpha.59 - 2026-01-11

## Изменено

## v3.0.0-alpha.58 - 2026-01-09

## Исправлено

## v3.0.0-alpha.57 - 2026-01-05

## Изменено

## Исправлено

## v3.0.0-alpha.56 - 2026-01-04

## Добавлено

## Изменено

## Исправлено

## Удалено

## v3.0.0-alpha.55 - 2026-01-02

## Изменено

## Исправлено

## Удалено

## v3.0.0-alpha.54 - 2025-12-29

## Добавлено

## Исправлено

## Удалено

## v3.0.0-alpha.53 - 2025-12-27

## Добавлено

## Исправлено

## v3.0.0-alpha.52 - 2025-12-26

## Исправлено

## v3.0.0-alpha.51 - 2025-12-23

## Исправлено

## v3.0.0-alpha.50 - 2025-12-21

## Изменено

## v3.0.0-alpha.49 - 2025-12-18

## Изменено

## v3.0.0-alpha.48 - 2025-12-16

## Добавлено

## Изменено

## Исправлено

## v3.0.0-alpha.47 - 2025-12-15

## Добавлено

## Исправлено

## v3.0.0-alpha.46 - 2025-12-14

## Добавлено

## Удалено

## v3.0.0-alpha.45 - 2025-12-13

## Добавлено

## Исправлено

## v3.0.0-alpha.44 - 2025-12-12

## Добавлено

## Изменено

## Исправлено

## v3.0.0-alpha.43 - 2025-12-11

## Добавлено

## v3.0.0-alpha.42 - 2025-12-10

## Добавлено

## v3.0.0-alpha.41 - 2025-11-23

## Исправлено

## v3.0.0-alpha.40 - 2025-11-13

## Исправлено

## v3.0.0-alpha.39 - 2025-11-12

## Добавлено

## Изменено

## v3.0.0-alpha.38 - 2025-11-04

## Добавлено

## Изменено

## Исправлено

## v3.0.0-alpha.37 - 2025-11-02

## Исправлено

## v3.0.0-alpha.36 - 2025-10-15

## Исправлено

## v3.0.0-alpha.35 - 2025-10-14

## Исправлено

## v3.0.0-alpha.34 - 2025-10-06

## Добавлено

## Исправлено

## v3.0.0-alpha.33 - 2025-10-04

## Исправлено

## v3.0.0-alpha.32 - 2025-10-02

## Исправлено

## v3.0.0-alpha.31 - 2025-09-27

## Исправлено

## v3.0.0-alpha.30 - 2025-09-26

## Исправлено

## v3.0.0-alpha.29 - 2025-09-25

## Добавлено

## Изменено

## Исправлено

## v3.0.0-alpha.29 - 2025-09-25

## Добавлено

## Изменено

## Исправлено

## v3.0.0-alpha.27 - 2025-09-07

## Исправлено

## v3.0.0-alpha.26 - 2025-08-24

## Добавлено

## v3.0.0-alpha.25 - 2025-08-16

## Изменено

больше не игнорируются и переопределяют значение конфигурации.

## v3.0.0-alpha.24 - 2025-08-13

## Добавлено

## v3.0.0-alpha.23 - 2025-08-11

## Исправлено

## v3.0.0-alpha.22 - 2025-08-10

## Добавлено

## Изменено

+ Исправлены чрезмерно широкие зависимости пакетов Linux и устаревшие зависимости RPM.

## v3.0.0-alpha.21 - 2025-08-07

## Исправлено

## v3.0.0-alpha.20 - 2025-08-06

## Исправлено

## v3.0.0-alpha.19 - 2025-08-05

## Добавлено

## Исправлено

## v3.0.0-alpha.18 - 2025-08-03

## Добавлено

## Исправлено

## v3.0.0-alpha.17 - 2025-07-31

## Исправлено

## v3.0.0-alpha.16 - 2025-07-25

## Добавлено

## v3.0.0-alpha.15 - 2025-07-25

## Добавлено

## v3.0.0-alpha.14 - 2025-07-25

## Добавлено

## v3.0.0-alpha.12 - 2025-07-15

### Добавлено

### Исправлено

## v3.0.0-alpha.11 - 2025-07-12

## Добавлено

## v3.0.0-alpha.10 - 2025-07-06

### Несовместимые изменения

### Добавлено

### Исправлено

### Изменено

## v3.0.0-alpha.9 - 2025-01-13

### Добавлено

### Исправлено

### Изменено

## v3.0.0-alpha.8.3 - 2024-12-07

### Изменено

## v3.0.0-alpha.8.2 - 2024-12-07

### Изменено

`go install`, автор: @leaanthony

## v3.0.0-alpha.8.1 - 2024-12-07

### Изменено

`go install`, автор: @leaanthony

## v3.0.0-alpha.8 - 2024-12-06

### Добавлено

@atterpac в [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) в   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) в   [#3867](https://github.com/wailsapp/wails/pull/3867)   — [atterpac](https://github.com/atterpac) в   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) в   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   размещено в заданных координатах X/Y —   [leaanthony](https://github.com/leaanthony) в   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) и   [leaanthony](https://github.com/leaanthony) в   [#3823](https://github.com/wailsapp/wails/pull/3823)   — [leaanthony](https://github.com/leaanthony) в   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) в   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### Изменено

`service.OnShutdown`для всех ранее запущенных сервисов — @atterpac в   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac в [#3907](https://github.com/wailsapp/wails/pull/3907)   вложенные папки — @atterpac в   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) в   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) в   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (заменено параметрами `EnabledFeatures` и `DisabledFeatures`) —   [leaanthony](https://github.com/leaanthony)

### Исправлено

переменная канала — @michael-freling в   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) в   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   в [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) в   [#3841](https://github.com/wailsapp/wails/pull/3841)   `PasteAndMatchStyle` роль в меню редактирования на Darwin —   [johnmccabe](https://github.com/johnmccabe) в   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) в   [#3854](https://github.com/wailsapp/wails/pull/3854) —   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   различаются — @nickisworking в   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 — 2024-09-18

### Добавлено

[mmghv](https://github.com/mmghv) в   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) и   [leaanthony](https://github.com/leaanthony) в   [#3570](https://github.com/wailsapp/wails/pull/3570)

### Изменено

События приложения `OnWindowEvent` → события окна —   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   ветки с префиксом `v3/` или `v3-` —   [stendler](https://github.com/stendler) в   [#3747](https://github.com/wailsapp/wails/pull/3747)

### Исправлено

[etesam913](https://github.com/etesam913) в   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) в   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) в   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) в   [#3614](https://github.com/wailsapp/wails/pull/3614) —   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720) —   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) —   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720) —   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) —   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746) —   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 — 2024-07-30

### Исправлено

## v3.0.0-alpha.5 — 2024-07-30

### Добавлено

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   щелчок по значку — @5aaee9 в [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) в   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   в [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) в   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) в   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   на основе [PR](https://github.com/wailsapp/wails/pull/2044) от @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   пакет npm — [@fbbdev](https://github.com/fbbdev) в   [#3334](https://github.com/wailsapp/wails/pull/3334)   в [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` — [@abichinger](https://github.com/abichinger) в   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) в   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) в   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) в   [#3674](https://github.com/wailsapp/wails/pull/3674)

### Исправлено

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) в   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) в   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) в   [#3477](https://github.com/wailsapp/wails/pull/3477)   от [Atterpac](https://github.com/atterpac) в   [#3320](https://github.com/wailsapp/wails/pull/3320).   в [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) в   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) в   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) в   [PR](https://github.com/wailsapp/wails/pull/3039)   установлен. Добавлено [@pylotlight](https://github.com/pylotlight) в   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   пробелы — @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) в PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) в PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   прикреплённое окно [tw1nk](https://github.com/tw1nk) в PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) в   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` присутствует.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   обработчики событий от [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) в   [#3330](https://github.com/wailsapp/wails/pull/3330)   генератор от [@fbbdev](https://github.com/fbbdev) в   [#3334](https://github.com/wailsapp/wails/pull/3334)   генератор от [@fbbdev](https://github.com/fbbdev) в   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) в   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) в   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) в   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) в   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) в   [#3431](https://github.com/wailsapp/wails/pull/3431)   от [@pekim](https://github.com/pekim) в   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) от   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) в   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows от [@bruxaodev](https://github.com/bruxaodev/) в   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Изменено

[mmghv](https://github.com/mmghv) в   [#3611](https://github.com/wailsapp/wails/pull/3611)   поддержка Ubuntu 24.04 LTS от [atterpac](https://github.com/atterpac) в   [#3461](https://github.com/wailsapp/wails/pull/3461)   должен иметь атрибут `type="module"`. Автор:   [@fbbdev](https://github.com/fbbdev), в   [#3295](https://github.com/wailsapp/wails/pull/3295)   объект и не запускает систему WML. Это сделано для улучшения   инкапсуляции. При необходимости систему WML можно запустить вручную, вызвав   новый метод `WML.Enable`. Встроенный скрипт среды выполнения JS по-прежнему автоматически выполняет обе   операции. Автор: [@fbbdev](https://github.com/fbbdev), в   [#3295](https://github.com/wailsapp/wails/pull/3295)   объект окна в качестве экспорта по умолчанию. Импортировать   отдельные методы с помощью синтаксиса именованного импорта или импорта пространства имён ESM больше нельзя.   API. У некоторых методов изменились имя или прототип, а именно: `Screen`   заменён на `GetScreen`; `GetZoomLevel`/`SetZoomLevel` заменены на `GetZoom`/`SetZoom`;   `GetZoom`, `Width` и `Height` теперь возвращают значения напрямую, а не оборачивают   их в объекты. Автор: [@fbbdev](https://github.com/fbbdev), в   [#3295](https://github.com/wailsapp/wails/pull/3295)   удалён. Чтобы вернуться к вызовам по имени, используйте параметр CLI `-names`.   Автор: [@fbbdev](https://github.com/fbbdev), в   [#3468](https://github.com/wailsapp/wails/pull/3468)   именовались по содержащим их пакетам; теперь используются полные пути импорта Go,   включая путь модуля. Автор: [@fbbdev](https://github.com/fbbdev), в   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Автор: [@fbbdev](https://github.com/fbbdev), в   [#3468](https://github.com/wailsapp/wails/pull/3468)   вызов `application.NewService`. Автор: [@fbbdev](https://github.com/fbbdev), в   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) в   [#3574](https://github.com/wailsapp/wails/pull/3574)
