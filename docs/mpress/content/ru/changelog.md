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

/_-- Все значимые изменения этого проекта будут документироваться в этом файле.

Формат основан на [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), а проект придерживается [семантического версионирования](https://semver.org/spec/v2.0.0.html).

- `Added` для новых возможностей.
- `Changed` для изменений существующей функциональности.
- `Deprecated` для возможностей, которые вскоре будут удалены.
- `Removed` для уже удалённых возможностей.
- `Fixed` для любых исправлений ошибок.
- `Security` для уязвимостей.

_/

/_   * ОБНОВЛЯЕТСЯ АВТОМАТИЧЕСКИ *   Примечания к выпускам формируются из объединённых PR. Не редактируйте этот файл в обычных PR от участников и не создавайте заново удалённый `v3/UNRELEASED_CHANGELOG.md`. Инструкции по ведению журнала изменений см. в `AGENTS.md`. _/

## [Не выпущено]

## v3.2.0-nightly.20261102 — 2026-11-02

## Изменено
- DisableOverscroll вместо DisableBounce в [PR](https://github.com/wailsapp/wails/pull/6077) от @lethusang

## v3.2.0-nightly.20261101 — 2026-11-01

## Добавлено
- добавлены экспериментальные бэкенды macOS и Linux без CGO в [PR](https://github.com/wailsapp/wails/pull/6521) от @taliesin-ai
- нативные сборки из YAML, переносимая упаковка и подписание в [PR](https://github.com/wailsapp/wails/pull/6531) от @taliesin-ai

## Исправлено
- исправлены различные проблемы ночной сборки v3 и очищены зависимости в [PR](https://github.com/wailsapp/wails/pull/6519) от @taliesin-ai
- исправлены проверка выпуска v3 и классификация фикстур в [PR](https://github.com/wailsapp/wails/pull/6533) от @taliesin-ai
- устранены сбои межплатформенной проверки выпусков в [PR](https://github.com/wailsapp/wails/pull/6538) от @taliesin-ai

## v3.1.0 — 2026-11-01

## Добавлено
- добавлены экспериментальные бэкенды macOS и Linux без CGO в [PR](https://github.com/wailsapp/wails/pull/6521) от @taliesin-ai
- нативные сборки из YAML, переносимая упаковка и подписание в [PR](https://github.com/wailsapp/wails/pull/6531) от @taliesin-ai

## Изменено
- исключено копирование отброшенных метаданных рабочего стола в [PR](https://github.com/wailsapp/wails/pull/6391) от @taliesin-ai
- длина содержимого форматируется без fmt в [PR](https://github.com/wailsapp/wails/pull/6393) от @taliesin-ai
- корни индекса выбираются без сбора кандидатов в [PR](https://github.com/wailsapp/wails/pull/6404) от @taliesin-ai
- буферы копирования файлов без произвольного доступа используются повторно в [PR](https://github.com/wailsapp/wails/pull/6408) от @taliesin-ai
- исключено выделение памяти под reader для анализа префиксов в [PR](https://github.com/wailsapp/wails/pull/6410) от @taliesin-ai
- для известных расширений пропускается чтение с целью определения MIME-типа в [PR](https://github.com/wailsapp/wails/pull/6414) от @taliesin-ai
- предотвращены рост очереди и колебания её сокращения в [PR](https://github.com/wailsapp/wails/pull/6415) от @taliesin-ai
- исключено выделение памяти под отбрасываемые контексты запросов в [PR](https://github.com/wailsapp/wails/pull/6419) от @taliesin-ai
- пропускается подготовка отключённых отладочных журналов в [PR](https://github.com/wailsapp/wails/pull/6420) от @taliesin-ai
- буферы для определения MIME-типа выделяются отложенно в [PR](https://github.com/wailsapp/wails/pull/6428) от @taliesin-ai
- сокращено число выделений памяти при форматировании гиперссылок в [PR](https://github.com/wailsapp/wails/pull/6443) от @taliesin-ai
- устранены временные выделения памяти при форматировании маркеров списка в [PR](https://github.com/wailsapp/wails/pull/6449) от @taliesin-ai
- преобразование UUID в строку сведено к одному выделению памяти в [PR](https://github.com/wailsapp/wails/pull/6453) от @taliesin-ai
- сокращено число выделений памяти при форматировании дескрипторов выражений в [PR](https://github.com/wailsapp/wails/pull/6454) от @taliesin-ai
- устранено временное выделение памяти под дайджест UUID на основе имени в [PR](https://github.com/wailsapp/wails/pull/6455) от @taliesin-ai
- исключено копирование незавершённых строк вывода команд в [PR](https://github.com/wailsapp/wails/pull/6460) от @taliesin-ai
- исполняемые файлы копируются потоково при подготовке пакета в [PR](https://github.com/wailsapp/wails/pull/6472) от @taliesin-ai
- исключено повторное сканирование незавершённых строк вывода в [PR](https://github.com/wailsapp/wails/pull/6473) от @taliesin-ai
- скомпилированные шаблоны маршрутизации команд используются повторно в [PR](https://github.com/wailsapp/wails/pull/6475) от @taliesin-ai
- исключён поиск метаданных при обнаружении среды выполнения в [PR](https://github.com/wailsapp/wails/pull/6482) от @taliesin-ai
- сокращено число выделений памяти для локализованной резервной страницы в [PR](https://github.com/wailsapp/wails/pull/6485) от @taliesin-ai
- исключена повторная компиляция шаблонов для каждого обработчика в [PR](https://github.com/wailsapp/wails/pull/6490) от @taliesin-ai

## Исправлено
- реализован SetTooltip для системного трея в [PR](https://github.com/wailsapp/wails/pull/6017) от @fan711
- тёмный режим считывается из портала настроек в [PR](https://github.com/wailsapp/wails/pull/6072) от @TechHutTV
- Window.SetZoom теперь может устанавливать значение ниже 1.0 в [PR](https://github.com/wailsapp/wails/pull/6166) от @APshenkin
- события доставляются каждому слушателю Go в порядке вызова emit в [PR](https://github.com/wailsapp/wails/pull/6275) от @DevLumuz
- хуки событий в Emit считываются под блокировкой hookLock в [PR](https://github.com/wailsapp/wails/pull/6276) от @DevLumuz
- свойство Menu трея теперь объявляется всегда, чтобы трей, созданный после Run, не отбрасывался, в [PR](https://github.com/wailsapp/wails/pull/6284) от @AlbinoGeek
- работа для главного потока, отправленная до Run, удерживается до запуска Run в [PR](https://github.com/wailsapp/wails/pull/6287) от @AlbinoGeek
- в провайдере обновлений GitHub используется правильный URL ресурса в [PR](https://github.com/wailsapp/wails/pull/6304) от @gavidroselj
- тело запроса сохраняется, когда WebView2 не передаёт Content-Length, в [PR](https://github.com/wailsapp/wails/pull/6326) от @strraaa
- экраны macOS и GTK представляются в пространстве DIP (границы при смешанном DPI) в [PR](https://github.com/wailsapp/wails/pull/6346) от @julianstorer
- отклонённое уведомление не регистрируется как щелчок в [PR](https://github.com/wailsapp/wails/pull/6367) от @bercedev
- процессы разработки останавливаются при закрытии терминала в [PR](https://github.com/wailsapp/wails/pull/6371) от @xmasdev
- при проверке подтверждений запуски журнала изменений сопоставляются с коммитами squash-слияния в [PR](https://github.com/wailsapp/wails/pull/6373) от @taliesin-ai
- устранение дребезга при перезапуске панели задач удерживается в течение секунды в [PR](https://github.com/wailsapp/wails/pull/6386) от @taliesin-ai
- ограничения размера сохраняются при переходе в полноэкранный режим из начального полноэкранного состояния в [PR](https://github.com/wailsapp/wails/pull/6387) от @taliesin-ai
- полноэкранное или развёрнутое окно не изменяет размер для соответствия новому ограничению в [PR](https://github.com/wailsapp/wails/pull/6388) от @taliesin-ai
- не показывать диалоговое окно до вызова Run в примере liquid-glass в [PR](https://github.com/wailsapp/wails/pull/6389) от @taliesin-ai
- диалоговое окно, запрошенное до вызова Run, ничего не делает вместо бесконечного ожидания в [PR](https://github.com/wailsapp/wails/pull/6390) от @taliesin-ai
- использовать резервный вариант, если в desktop entry отсутствует команда, в [PR](https://github.com/wailsapp/wails/pull/6392) от @taliesin-ai
- определять тип текущего содержимого без сохранения имён файлов в [PR](https://github.com/wailsapp/wails/pull/6394) от @taliesin-ai
- распознавать расширения ресурсов с буквами в смешанном регистре в [PR](https://github.com/wailsapp/wails/pull/6395) от @taliesin-ai
- освобождать использованные буферы определения типа содержимого в [PR](https://github.com/wailsapp/wails/pull/6396) от @taliesin-ai
- освобождать удалённые ссылки ringqueue в [PR](https://github.com/wailsapp/wails/pull/6397) от @taliesin-ai
- учитывать удаление заголовков до выбора резервного варианта в [PR](https://github.com/wailsapp/wails/pull/6398) от @taliesin-ai
- фиксировать заголовки резервного ответа перед сбросом буфера в [PR](https://github.com/wailsapp/wails/pull/6399) от @taliesin-ai
- обнаруживать только файлы с точным именем index в [PR](https://github.com/wailsapp/wails/pull/6400) от @taliesin-ai
- не формировать тела и не выполнять полное чтение для незапросопозиционируемых данных при HEAD-запросах в [PR](https://github.com/wailsapp/wails/pull/6401) от @taliesin-ai
- заполнять префикс для определения MIME-типа при нескольких коротких чтениях в [PR](https://github.com/wailsapp/wails/pull/6402) от @taliesin-ai
- фиксировать успешные пустые ответы из незапросопозиционируемых источников в [PR](https://github.com/wailsapp/wails/pull/6403) от @taliesin-ai
- избегать дублирования целей URI после подстановки в [PR](https://github.com/wailsapp/wails/pull/6405) от @taliesin-ai
- сохранять пути при подстановке полей в [PR](https://github.com/wailsapp/wails/pull/6406) от @taliesin-ai
- учитывать каталоги поиска desktop-файлов XDG в [PR](https://github.com/wailsapp/wails/pull/6407) от @taliesin-ai
- фиксировать заголовки определения типа содержимого при сбросе пустого буфера в [PR](https://github.com/wailsapp/wails/pull/6409) от @taliesin-ai
- пропускать каталоги, затеняющие desktop entries, в [PR](https://github.com/wailsapp/wails/pull/6411) от @taliesin-ai
- сохранять цели, найденные внутри путей к исполняемым файлам, в [PR](https://github.com/wailsapp/wails/pull/6412) от @taliesin-ai
- освобождать захваченные обратным вызовом данные после срабатывания таймера в [PR](https://github.com/wailsapp/wails/pull/6413) от @taliesin-ai
- исключить ложные совпадения библиотек ldconfig в [PR](https://github.com/wailsapp/wails/pull/6416) от @taliesin-ai
- ограничить совпадения при резервном поиске библиотек в файловой системе в [PR](https://github.com/wailsapp/wails/pull/6417) от @taliesin-ai
- устранить гонку данных при инвалидации кеша в [PR](https://github.com/wailsapp/wails/pull/6418) от @taliesin-ai
- учитывать необязательную конфигурацию журналирования в [PR](https://github.com/wailsapp/wails/pull/6421) от @taliesin-ai
- избегать изменения общего начального URL и использования устаревших портов в [PR](https://github.com/wailsapp/wails/pull/6422) от @taliesin-ai
- прекращать повторные попытки прокси-сервера разработки при отмене запроса в [PR](https://github.com/wailsapp/wails/pull/6423) от @taliesin-ai
- возвращать 404 при отсутствии встроенных ресурсов в [PR](https://github.com/wailsapp/wails/pull/6424) от @taliesin-ai
- обрабатывать HEAD-запросы к встроенной среде выполнения в [PR](https://github.com/wailsapp/wails/pull/6425) от @taliesin-ai
- сохранять экранированные пути при маршрутизации сервисов в [PR](https://github.com/wailsapp/wails/pull/6426) от @taliesin-ai
- правильно упорядочивать чрезмерно большие числовые идентификаторы предварительных выпусков в [PR](https://github.com/wailsapp/wails/pull/6427) от @taliesin-ai
- отклонять ответы с ошибками HTTP до сохранения загружаемых инструментов в [PR](https://github.com/wailsapp/wails/pull/6429) от @taliesin-ai
- соблюдать предельный срок запуска при обнаружении desktop entries в [PR](https://github.com/wailsapp/wails/pull/6430) от @taliesin-ai
- сохранять тела запросов при повторных попытках прокси-сервера разработки в [PR](https://github.com/wailsapp/wails/pull/6432) от @taliesin-ai
- закрывать тела HTTP-ответов при любом пути возврата в [PR](https://github.com/wailsapp/wails/pull/6434) от @taliesin-ai
- сохранять ошибки отмены при параллельном поиске библиотек в [PR](https://github.com/wailsapp/wails/pull/6435) от @taliesin-ai
- отклонять неуспешные ответы на запрос тегов GitHub до декодирования в [PR](https://github.com/wailsapp/wails/pull/6436) от @taliesin-ai
- сохранять буквальные пути шаблонов с квадратными скобками в [PR](https://github.com/wailsapp/wails/pull/6437) от @taliesin-ai
- устанавливать wails3 один раз и сохранять ошибки установщика в [PR](https://github.com/wailsapp/wails/pull/6438) от @taliesin-ai
- игнорировать скрытые записи и заголовки в блоках комментариев в [PR](https://github.com/wailsapp/wails/pull/6439) от @taliesin-ai
- выводить ошибки поиска стабильного выпуска в [PR](https://github.com/wailsapp/wails/pull/6440) от @taliesin-ai
- печатать обычные сообщения без интерпретации форматирования в [PR](https://github.com/wailsapp/wails/pull/6441) от @taliesin-ai
- предотвращать панику при записи после загрузки JSON null в [PR](https://github.com/wailsapp/wails/pull/6442) от @taliesin-ai
- возвращать снимки состояния при чтении всего хранилища в [PR](https://github.com/wailsapp/wails/pull/6444) от @taliesin-ai
- передавать ошибки после частичного перебора строк в [PR](https://github.com/wailsapp/wails/pull/6446) от @taliesin-ai
- сохранять пути назначения распаковки буквально в [PR](https://github.com/wailsapp/wails/pull/6447) от @taliesin-ai
- освобождать дескрипторы подготовленных выражений при закрытии и повторном открытии в [PR](https://github.com/wailsapp/wails/pull/6448) от @taliesin-ai
- обрабатывать гонку подготовки с закрытием или повторным открытием в [PR](https://github.com/wailsapp/wails/pull/6450) от @taliesin-ai
- корректно обрабатывать HEAD-запросы к резервной индексной странице в [PR](https://github.com/wailsapp/wails/pull/6451) от @taliesin-ai
- требовать полные шестнадцатеричные дескрипторы выражений в [PR](https://github.com/wailsapp/wails/pull/6452) от @taliesin-ai
- сохранять стандартную базу данных в памяти при конкурентном доступе в [PR](https://github.com/wailsapp/wails/pull/6456) от @taliesin-ai
- отделять stderr от разбираемого вывода HEAD в [PR](https://github.com/wailsapp/wails/pull/6457) от @taliesin-ai
- создавать ресурсы среды выполнения из встроенных пакетов в [PR](https://github.com/wailsapp/wails/pull/6458) от @taliesin-ai
- обнаруживать заполнители TODO, FIXME и TBD в [PR](https://github.com/wailsapp/wails/pull/6459) от @taliesin-ai
- проверять ответы API примечаний к выпуску в [PR](https://github.com/wailsapp/wails/pull/6461) от @taliesin-ai
- учитывать каталоги задач при проверке состояния в [PR](https://github.com/wailsapp/wails/pull/6462) от @taliesin-ai
- предотвращать потерю данных при копировании файла в самого себя в [PR](https://github.com/wailsapp/wails/pull/6463) от @taliesin-ai
- включать каталог задачи в идентификатор кеша в [PR](https://github.com/wailsapp/wails/pull/6464) от @taliesin-ai
- проверять порты один раз за проверку в [PR](https://github.com/wailsapp/wails/pull/6465) от @taliesin-ai
- избегать паники, когда записи кеша декодируются как null, в [PR](https://github.com/wailsapp/wails/pull/6466) от @taliesin-ai
- увеличивать числовые номера предварительных выпусков без ограничений целочисленного типа в [PR](https://github.com/wailsapp/wails/pull/6467) от @taliesin-ai
- не позволять устаревшим проверкам скрывать промахи кеша в [PR](https://github.com/wailsapp/wails/pull/6468) от @taliesin-ai
- выполнять стандартные задачи в базовом каталоге исполнителя в [PR](https://github.com/wailsapp/wails/pull/6469) от @taliesin-ai
- сохранять буквальные знаки процента в командах проводника в [PR](https://github.com/wailsapp/wails/pull/6470) от @taliesin-ai
- согласовать каталоги команд для предварительных условий и задач в [PR](https://github.com/wailsapp/wails/pull/6471) от @taliesin-ai
- заменить символические ссылки сгенерированными обычными файлами в [PR](https://github.com/wailsapp/wails/pull/6474) от @taliesin-ai
- сохранять существующие DMG при сбое сборки в [PR](https://github.com/wailsapp/wails/pull/6476) от @taliesin-ai
- сохранять обратные косые черты в заключённых в кавычки регулярных выражениях и путях в [PR](https://github.com/wailsapp/wails/pull/6477) от @taliesin-ai
- проверять размеры сгенерированных значков до выделения памяти под изображения в [PR](https://github.com/wailsapp/wails/pull/6478) от @taliesin-ai
- выполнять цепочки команд и перенаправления через оболочку в [PR](https://github.com/wailsapp/wails/pull/6479) от @taliesin-ai
- сохранять предыдущие пакеты Linux при сбое сборки в [PR](https://github.com/wailsapp/wails/pull/6480) от @taliesin-ai
- использовать резервное выполнение через оболочку с учётом платформы в [PR](https://github.com/wailsapp/wails/pull/6481) от @taliesin-ai
- обслуживать ресурсы с повторяющимися косыми чертами в начале пути в [PR](https://github.com/wailsapp/wails/pull/6483) от @taliesin-ai
- распознавать стабильные версии с метаданными сборки в [PR](https://github.com/wailsapp/wails/pull/6484) от @taliesin-ai
- точно сообщать об ошибках создания ресурсов в [PR](https://github.com/wailsapp/wails/pull/6486) от @taliesin-ai
- предотвращать состояния гонки при замене активного репортера в [PR](https://github.com/wailsapp/wails/pull/6487) от @taliesin-ai
- сериализовать одновременные сохранения файлов кэша в [PR](https://github.com/wailsapp/wails/pull/6488) от @taliesin-ai
- поддерживать несколько кандидатов на равенство в шаблонах в [PR](https://github.com/wailsapp/wails/pull/6489) от @taliesin-ai
- нормализовать сохранённые карты категорий со значением null в [PR](https://github.com/wailsapp/wails/pull/6491) от @taliesin-ai
- показывать диагностические данные при сбое команд оболочки с переменными в [PR](https://github.com/wailsapp/wails/pull/6492) от @taliesin-ai
- исключить цепочки команд оболочки из неявного кэша Go в [PR](https://github.com/wailsapp/wails/pull/6493) от @taliesin-ai
- сериализовать сохранение категорий в Linux в [PR](https://github.com/wailsapp/wails/pull/6494) от @taliesin-ai
- кэшировать выходные файлы Go с пробелами в кавычках или с экранированными пробелами в [PR](https://github.com/wailsapp/wails/pull/6495) от @taliesin-ai
- сохранять значки Mac при сбое кодирования в [PR](https://github.com/wailsapp/wails/pull/6496) от @taliesin-ai
- не отбрасывать недопустимые флаги количества запусков тестов Go в [PR](https://github.com/wailsapp/wails/pull/6497) от @taliesin-ai
- сохранять разрешения исходного файла при создании его копии в [PR](https://github.com/wailsapp/wails/pull/6498) от @taliesin-ai
- исправить различные проблемы ночной сборки v3 и очистить зависимости в [PR](https://github.com/wailsapp/wails/pull/6519) от @taliesin-ai
- исправить проверку выпуска v3 и классификацию тестовых данных в [PR](https://github.com/wailsapp/wails/pull/6533) от @taliesin-ai
- устранить межплатформенные сбои проверки выпуска в [PR](https://github.com/wailsapp/wails/pull/6538) от @taliesin-ai

## v3.1.0-nightly.20261028 — 2026-10-28

## Изменено
- ускорить создание журнала изменений за счёт повторного использования скомпилированных регулярных выражений в [PR](https://github.com/wailsapp/wails/pull/6490) от @taliesin-ai
- сократить количество выделений памяти при обслуживании локализованных резервных страниц в [PR](https://github.com/wailsapp/wails/pull/6485) от @taliesin-ai
- ускорить обнаружение среды выполнения AppImage, исключив ненужные обращения к метаданным, в [PR](https://github.com/wailsapp/wails/pull/6482) от @taliesin-ai
- повторно использовать скомпилированные шаблоны маршрутизации команд для повышения производительности wake в [PR](https://github.com/wailsapp/wails/pull/6475) от @taliesin-ai
- повысить производительность, исключив повторное сканирование незавершённых строк вывода, в [PR](https://github.com/wailsapp/wails/pull/6473) от @taliesin-ai
- копировать исполняемые файлы потоково при подготовке пакета MSIX, чтобы сократить использование памяти и ускорить сборку, в [PR](https://github.com/wailsapp/wails/pull/6472) от @taliesin-ai
- не копировать незавершённые строки вывода команд для повышения производительности в [PR](https://github.com/wailsapp/wails/pull/6460) от @taliesin-ai
- сократить количество выделений памяти при создании UUID на основе имени в [PR](https://github.com/wailsapp/wails/pull/6455) от @taliesin-ai
- сократить количество выделений памяти при форматировании дескрипторов операторов SQLite в [PR](https://github.com/wailsapp/wails/pull/6454) от @taliesin-ai
- свести преобразование UUID в строку к одному выделению памяти в [PR](https://github.com/wailsapp/wails/pull/6453) от @taliesin-ai
- сократить количество выделений памяти при форматировании пунктов журнала изменений в [PR](https://github.com/wailsapp/wails/pull/6449) от @taliesin-ai
- сократить количество выделений памяти при форматировании гиперссылок терминала в [PR](https://github.com/wailsapp/wails/pull/6443) от @taliesin-ai
- сократить количество выделений памяти при обслуживании ресурсов за счёт отложенного анализа MIME-буферов в [PR](https://github.com/wailsapp/wails/pull/6428) от @taliesin-ai
- пропускать подготовку отладочных сообщений на сервере ресурсов, если отладочное журналирование отключено, в [PR](https://github.com/wailsapp/wails/pull/6420) от @taliesin-ai
- сократить количество выделений памяти при обработке запросов webview, исключив отбрасываемые контексты запросов, в [PR](https://github.com/wailsapp/wails/pull/6419) от @taliesin-ai
- ускорить обслуживание ресурсов, пропуская определение MIME для известных расширений файлов, в [PR](https://github.com/wailsapp/wails/pull/6414) от @taliesin-ai
- сократить количество выделений памяти при обслуживании ресурсов, отказавшись от средства чтения для определения типа содержимого, в [PR](https://github.com/wailsapp/wails/pull/6410) от @taliesin-ai
- сократить количество выделений памяти при обслуживании файлов ресурсов без поддержки позиционирования в [PR](https://github.com/wailsapp/wails/pull/6408) от @taliesin-ai
- ускорить выбор корневого каталога индекса сервером ресурсов, исключив сбор кандидатов, в [PR](https://github.com/wailsapp/wails/pull/6404) от @taliesin-ai
- повысить производительность сервера ресурсов, форматируя длину содержимого без fmt, в [PR](https://github.com/wailsapp/wails/pull/6393) от @taliesin-ai
- сократить копирование памяти при чтении метаданных проводника файлов рабочего стола в [PR](https://github.com/wailsapp/wails/pull/6391) от @taliesin-ai

## Исправлено
- сохранять разрешения исходного файла при копировании в новые файлы назначения в [PR](https://github.com/wailsapp/wails/pull/6498) от @taliesin-ai
- не изменять размер полноэкранных и развёрнутых окон при изменении ограничений минимального или максимального размера в [PR](https://github.com/wailsapp/wails/pull/6388) от @taliesin-ai
- публиковать запланированные обновления истории количества звёзд репозитория wailsapp/wails в [PR](https://github.com/wailsapp/wails/pull/6500) от @taliesin-ai
- не отбрасывать недопустимые флаги количества запусков тестов Go в wake в [PR](https://github.com/wailsapp/wails/pull/6497) от @taliesin-ai
- сохранять значки Mac при сбое кодирования в [PR](https://github.com/wailsapp/wails/pull/6496) от @taliesin-ai
- исправить промахи кэша сборки для выходных файлов Go с пробелами в кавычках или с экранированными пробелами в [PR](https://github.com/wailsapp/wails/pull/6495) от @taliesin-ai
- сериализовать сохранение категорий уведомлений в Linux, чтобы предотвратить состояния гонки при одновременной записи, в [PR](https://github.com/wailsapp/wails/pull/6494) от @taliesin-ai
- исключить цепочки команд оболочки из неявного кэша сборки Go в [PR](https://github.com/wailsapp/wails/pull/6493) от @taliesin-ai
- исправить ошибку загрузки карт категорий уведомлений, когда сохранённые данные имеют значение null, в [PR](https://github.com/wailsapp/wails/pull/6491) от @taliesin-ai
- показывать диагностические данные при сбое команд оболочки с переменными в [PR](https://github.com/wailsapp/wails/pull/6492) от @taliesin-ai
- исправить сопоставление шаблонов wake при наличии нескольких кандидатов на равенство в [PR](https://github.com/wailsapp/wails/pull/6489) от @taliesin-ai
- Сериализовать параллельные сохранения файлов кеша, чтобы предотвратить повреждение данных при записи кеша, в [PR](https://github.com/wailsapp/wails/pull/6488) от @taliesin-ai
- Предотвратить состояния гонки при замене активного средства формирования отчётов в [PR](https://github.com/wailsapp/wails/pull/6487) от @taliesin-ai
- Обеспечить точное сообщение об ошибках создания ресурсов в CLI в [PR](https://github.com/wailsapp/wails/pull/6486) от @taliesin-ai
- Распознавать в CLI стабильные версии, содержащие метаданные сборки, в [PR](https://github.com/wailsapp/wails/pull/6484) от @taliesin-ai
- Обслуживать ресурсы, пути запросов к которым содержат несколько начальных косых черт, в [PR](https://github.com/wailsapp/wails/pull/6483) от @taliesin-ai
- Исправить выполнение резервных команд оболочки для wake на каждой платформе в [PR](https://github.com/wailsapp/wails/pull/6481) от @taliesin-ai
- Сохранять предыдущие пакеты Linux при сбое сборки в [PR](https://github.com/wailsapp/wails/pull/6480) от @taliesin-ai
- Выполнять цепочки команд и перенаправления через оболочку в [PR](https://github.com/wailsapp/wails/pull/6479) от @taliesin-ai
- Сохранять обратные косые черты в заключённых в кавычки регулярных выражениях и путях при разборе команд wake в [PR](https://github.com/wailsapp/wails/pull/6477) от @taliesin-ai
- Сохранять существующие файлы DMG при сбоях упаковки сборок в [PR](https://github.com/wailsapp/wails/pull/6476) от @taliesin-ai
- Заменять символические ссылки созданными обычными файлами в выходных данных привязок в [PR](https://github.com/wailsapp/wails/pull/6474) от @taliesin-ai
- Сохранять литеральные знаки процента в командах файлового менеджера Linux в [PR](https://github.com/wailsapp/wails/pull/6470) от @taliesin-ai
- Согласовать каталоги предварительных условий и команд задач, чтобы они разрешались единообразно, в [PR](https://github.com/wailsapp/wails/pull/6471) от @taliesin-ai
- Выполнять задачи по умолчанию в базовом каталоге исполнителя в [PR](https://github.com/wailsapp/wails/pull/6469) от @taliesin-ai
- Не допускать, чтобы устаревшие проверки скрывали промахи кеша, в [PR](https://github.com/wailsapp/wails/pull/6468) от @taliesin-ai
- Исправить нумерацию предварительных выпусков в CLI, чтобы она увеличивалась без ограничений, связанных с переполнением целых чисел, в [PR](https://github.com/wailsapp/wails/pull/6467) от @taliesin-ai
- Предотвратить панику при декодировании записей кеша как null в [PR](https://github.com/wailsapp/wails/pull/6466) от @taliesin-ai
- Проверять порты в CLI один раз за каждую проверку в [PR](https://github.com/wailsapp/wails/pull/6465) от @taliesin-ai
- Включить каталог задачи в идентификатор кеша, чтобы избежать устаревших результатов wake, в [PR](https://github.com/wailsapp/wails/pull/6464) от @taliesin-ai
- Предотвратить потерю данных, когда при упаковке файл копируется сам в себя, в [PR](https://github.com/wailsapp/wails/pull/6463) от @taliesin-ai
- Учитывать каталоги задач при проверках состояния в [PR](https://github.com/wailsapp/wails/pull/6462) от @taliesin-ai
- Обнаруживать заполнители TODO, FIXME и TBD в записях журнала изменений в [PR](https://github.com/wailsapp/wails/pull/6459) от @taliesin-ai
- Создавать ресурсы среды выполнения из встроенных пакетов в CLI в [PR](https://github.com/wailsapp/wails/pull/6458) от @taliesin-ai
- Отделять stderr от разбираемого вывода HEAD в командах git в [PR](https://github.com/wailsapp/wails/pull/6457) от @taliesin-ai
- Сохранять стандартную базу данных SQLite в памяти при параллельном доступе в [PR](https://github.com/wailsapp/wails/pull/6456) от @taliesin-ai
- Возвращать правильные ответы на запросы HEAD для резервной индексной страницы в [PR](https://github.com/wailsapp/wails/pull/6451) от @taliesin-ai
- Обрабатывать состояние гонки между подготовкой инструкций SQLite и закрытием или повторным открытием базы данных в [PR](https://github.com/wailsapp/wails/pull/6450) от @taliesin-ai
- Освобождать дескрипторы подготовленных инструкций при закрытии или повторном открытии базы данных SQLite в [PR](https://github.com/wailsapp/wails/pull/6448) от @taliesin-ai
- Сохранять пути назначения при извлечении в литеральном виде при использовании gosod в [PR](https://github.com/wailsapp/wails/pull/6447) от @taliesin-ai
- Передавать ошибки запросов SQLite после частичного перебора строк в [PR](https://github.com/wailsapp/wails/pull/6446) от @taliesin-ai
- Возвращать снимки при чтении всего хранилища в kvstore в [PR](https://github.com/wailsapp/wails/pull/6444) от @taliesin-ai
- Предотвратить паники при записи в kvstore после загрузки значения JSON null в [PR](https://github.com/wailsapp/wails/pull/6442) от @taliesin-ai
- Выводить обычные сообщения терминала без интерпретации спецификаторов форматирования в [PR](https://github.com/wailsapp/wails/pull/6441) от @taliesin-ai
- Отображать ошибки при поиске последнего стабильного выпуска в средстве обновления в [PR](https://github.com/wailsapp/wails/pull/6440) от @taliesin-ai
- Игнорировать скрытые записи и заголовки в блоках комментариев журнала изменений в [PR](https://github.com/wailsapp/wails/pull/6439) от @taliesin-ai
- Устанавливать wails3 один раз и отображать ошибки установщика в средстве обновления в [PR](https://github.com/wailsapp/wails/pull/6438) от @taliesin-ai
- Сохранять литеральные пути шаблонов в квадратных скобках в gosod в [PR](https://github.com/wailsapp/wails/pull/6437) от @taliesin-ai
- Отклонять неуспешные ответы на запросы тегов GitHub до их декодирования в средстве обновления в [PR](https://github.com/wailsapp/wails/pull/6436) от @taliesin-ai
- Сохранять ошибки отмены при параллельном поиске библиотек в [PR](https://github.com/wailsapp/wails/pull/6435) от @taliesin-ai
- Закрывать тела HTTP-ответов на каждом пути возврата, чтобы предотвратить утечки ресурсов, в [PR](https://github.com/wailsapp/wails/pull/6434) от @taliesin-ai
- Сохранять тела запросов при их повторной отправке прокси-сервером разработки в [PR](https://github.com/wailsapp/wails/pull/6432) от @taliesin-ai
- Соблюдать предельный срок запуска при обнаружении файлового менеджера рабочего стола в [PR](https://github.com/wailsapp/wails/pull/6430) от @taliesin-ai
- Отклонять HTTP-ошибки до сохранения загрузок инструментария AppImage в [PR](https://github.com/wailsapp/wails/pull/6429) от @taliesin-ai
- Правильно упорядочивать чрезмерно большие числовые идентификаторы предварительных выпусков при сравнении версий semver в [PR](https://github.com/wailsapp/wails/pull/6427) от @taliesin-ai
- Сохранять экранированные пути при маршрутизации служб на сервере ресурсов в [PR](https://github.com/wailsapp/wails/pull/6426) от @taliesin-ai
- Обрабатывать запросы HEAD к встроенной среде выполнения на сервере ресурсов в [PR](https://github.com/wailsapp/wails/pull/6425) от @taliesin-ai
- Возвращать 404 для отсутствующих встроенных ресурсов вместо неправильного ответа в [PR](https://github.com/wailsapp/wails/pull/6424) от @taliesin-ai
- Прекращать повторные попытки прокси-сервера разработки при отмене запроса в [PR](https://github.com/wailsapp/wails/pull/6423) от @taliesin-ai
- Избегать общей мутации начального URL и устаревших портов на сервере ресурсов в [PR](https://github.com/wailsapp/wails/pull/6422) от @taliesin-ai
- Учитывать необязательную конфигурацию средства журналирования на сервере ресурсов в [PR](https://github.com/wailsapp/wails/pull/6421) от @taliesin-ai
- Устранить гонку данных при аннулировании кеша путей к библиотекам в [PR](https://github.com/wailsapp/wails/pull/6418) от @taliesin-ai
- Избегать ложных совпадений библиотек при поиске в выводе ldconfig в [PR](https://github.com/wailsapp/wails/pull/6416) от @taliesin-ai
- Предотвратить рост очереди сервера ресурсов и колебания её размера под нагрузкой в [PR](https://github.com/wailsapp/wails/pull/6415) от @taliesin-ai
- Освобождать захваченные обратным вызовом данные после срабатывания таймера устранения дребезга, чтобы не удерживать устаревшие ссылки, в [PR](https://github.com/wailsapp/wails/pull/6413) от @taliesin-ai
- Сохранять цели файлового менеджера, найденные внутри путей к исполняемым файлам, в [PR](https://github.com/wailsapp/wails/pull/6412) от @taliesin-ai
- Пропускать каталоги, затеняющие записи рабочего стола в файловом менеджере, в [PR](https://github.com/wailsapp/wails/pull/6411) от @taliesin-ai
- Исправить отсутствие отправки заголовков сервера ресурсов при пустом теле ответа в [PR](https://github.com/wailsapp/wails/pull/6409) от @taliesin-ai
- Учитывать каталоги поиска desktop-файлов XDG при определении расположения значков приложений в [PR](https://github.com/wailsapp/wails/pull/6407) от @taliesin-ai
- Сохранять пути к файлам при раскрытии полей в файловом проводнике, [PR](https://github.com/wailsapp/wails/pull/6406) от @taliesin-ai
- Не дублировать раскрытые целевые URI в файловом проводнике, [PR](https://github.com/wailsapp/wails/pull/6405) от @taliesin-ai
- Завершать успешные пустые ответы для неперематываемых ресурсов, а не оставлять их в ожидании, [PR](https://github.com/wailsapp/wails/pull/6403) от @taliesin-ai
- Правильно заполнять префикс для определения MIME-типа при нескольких коротких чтениях на сервере ресурсов, [PR](https://github.com/wailsapp/wails/pull/6402) от @taliesin-ai
- Не читать тела ответов целиком для неперематываемых HEAD-запросов на сервере ресурсов, [PR](https://github.com/wailsapp/wails/pull/6401) от @taliesin-ai
- При разрешении запросов к каталогам на сервере ресурсов отдавать только точно совпадающие индексные файлы, [PR](https://github.com/wailsapp/wails/pull/6400) от @taliesin-ai
- Исправить ошибку, из-за которой отклонённые уведомления Linux регистрировались как нажатия на серверах без поддержки действий, [PR](https://github.com/wailsapp/wails/pull/6367) от @bercedev
- Останавливать процессы разработки при закрытии терминала, [PR](https://github.com/wailsapp/wails/pull/6371) от @xmasdev
- После перезапуска панели задач Windows выдерживать задержку устранения дребезга в одну секунду перед сбросом флага, [PR](https://github.com/wailsapp/wails/pull/6386) от @taliesin-ai
- Исправить потерю заголовков резервного ответа на сервере ресурсов, [PR](https://github.com/wailsapp/wails/pull/6399) от @taliesin-ai
- Учитывать удаление заголовков до выбора резервного ответа на сервере ресурсов, [PR](https://github.com/wailsapp/wails/pull/6398) от @taliesin-ai
- Освобождать удалённые ссылки кольцевой очереди на сервере ресурсов, чтобы избежать утечек памяти, [PR](https://github.com/wailsapp/wails/pull/6397) от @taliesin-ai
- Освобождать использованные буферы определения типа содержимого на сервере ресурсов, чтобы сократить потребление памяти, [PR](https://github.com/wailsapp/wails/pull/6396) от @taliesin-ai
- Распознавать расширения файлов ресурсов со смешанным регистром при их выдаче, [PR](https://github.com/wailsapp/wails/pull/6395) от @taliesin-ai
- Правильно определять тип содержимого на сервере ресурсов, не сохраняя имена файлов, [PR](https://github.com/wailsapp/wails/pull/6394) от @taliesin-ai
- Использовать значение по умолчанию, если в записи рабочего стола не указана команда, [PR](https://github.com/wailsapp/wails/pull/6392) от @taliesin-ai
- Не показывать диалог в примере liquid-glass до запуска приложения на системах, отличных от macOS, [PR](https://github.com/wailsapp/wails/pull/6389) от @taliesin-ai
- Сохранять ограничения размера окна при переходе в полноэкранный режим из ограниченного размера в Windows, [PR](https://github.com/wailsapp/wails/pull/6387) от @taliesin-ai
- Исправить бесконечную блокировку диалогов, вызванных до Run, вместо возврата из вызова, [PR](https://github.com/wailsapp/wails/pull/6390) от @taliesin-ai

## Безопасность
- Проверять размеры создаваемых значков до выделения памяти под изображение, чтобы предотвратить чрезмерное потребление памяти, [PR](https://github.com/wailsapp/wails/pull/6478) от @taliesin-ai
- Проверять ответы API примечаний к выпуску перед использованием, [PR](https://github.com/wailsapp/wails/pull/6461) от @taliesin-ai
- Требовать полные шестнадцатеричные дескрипторы выражений в плагине SQLite, [PR](https://github.com/wailsapp/wails/pull/6452) от @taliesin-ai
- Ограничить резервный поиск библиотек в файловой системе, чтобы предотвратить загрузку непредусмотренных библиотек, [PR](https://github.com/wailsapp/wails/pull/6417) от @taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## Изменено
- Добавить на страницу 404 ссылки на соответствующие страницы документации Wails v2 на v2.wails.io, [PR](https://github.com/wailsapp/wails/pull/6366) от @taliesin-ai
- Перенести сайт документации v2 на v2.wails.io, [PR](https://github.com/wailsapp/wails/pull/6365) от @taliesin-ai
- Размещать документацию v3 на wails.io, а документацию v2 — на v2.wails.io, [PR](https://github.com/wailsapp/wails/pull/6364) от @taliesin-ai
- Обновить примечания к выпуску v3.0.0 и исправить утверждение в бюллетене безопасности, [PR](https://github.com/wailsapp/wails/pull/6363) от @taliesin-ai

## Исправлено
- Последовательно передавать обратные вызовы обработчиков событий Go в порядке их отправки, [PR](https://github.com/wailsapp/wails/pull/6275) от @DevLumuz
- Всегда объявлять свойство Menu системного трея, чтобы в Linux не отбрасывались элементы трея, созданные после Run, [PR](https://github.com/wailsapp/wails/pull/6284) от @AlbinoGeek
- Удерживать работу основного потока, отправленную до Run, пока приложение не запустится, [PR](https://github.com/wailsapp/wails/pull/6287) от @AlbinoGeek
- Использовать правильный URL API ресурсов при скачивании ресурсов выпуска GitHub средством обновления, [PR](https://github.com/wailsapp/wails/pull/6304) от @gavidroselj
- Исправить границы и позиции экранов в конфигурациях macOS и Linux со смешанными значениями DPI, [PR](https://github.com/wailsapp/wails/pull/6346) от @julianstorer
- Сохранять тела запросов в Windows, когда WebView2 не передаёт заголовок Content-Length, [PR](https://github.com/wailsapp/wails/pull/6326) от @strraaa
- Устранить гонку данных при отправке событий, создавая снимок перехватчиков под блокировкой, [PR](https://github.com/wailsapp/wails/pull/6276) от @DevLumuz
- Исправить бездействие `SystemTray.SetTooltip()` в Linux: теперь задаётся описание всплывающей подсказки StatusNotifierItem (#6017)
- Запускать проверку подтверждающих данных журнала изменений в еженедельном процессе выпуска, [PR](https://github.com/wailsapp/wails/pull/6374) от @taliesin-ai
- Исправить проверку подтверждающих данных журнала изменений выпуска, чтобы она сопоставляла запуски с коммитами squash-слияния, [PR](https://github.com/wailsapp/wails/pull/6373) от @taliesin-ai
- Исправить определение тёмного режима в Linux: ранее на тёмном рабочем столе возвращался светлый режим до изменения цветовой схемы, а `linux:SystemThemeChanged` никогда не срабатывал в бэкенде GTK4. Теперь оба бэкенда используют одну реализацию на основе пространства имён портала `org.freedesktop.appearance` (#6072)
- `Window.SetZoom` больше не повышает коэффициенты ниже 1.0 до 1.0 для окон с рамкой (Windows, macOS, Linux): WebView2 `ZoomFactor`, WKWebView `magnification` и уровень масштабирования WebKitGTK принимают любое положительное значение, а параметр окна `Zoom` уже передавал такие значения без изменений. Для безрамочных окон сохраняется нижний предел 1.0, поскольку их области перетаскивания и изменения размера не корректируются с учётом масштаба ([#4590](https://github.com/wailsapp/wails/issues/4590)); теперь неположительные коэффициенты, NaN и Inf игнорируются, а не повышаются до 1.0 — от @APshenkin

## v3.0.0 - 2026-10-26

## Изменено
- Предварительные ночные выпуски публикуются как `vX.Y.0-nightly.YYYYMMDD` и служат предварительным представлением следующего еженедельного выпуска, а еженедельные выпуски увеличивают минорную версию (`vX.Y.0`), поэтому `@latest` в Go всегда разрешается в еженедельный выпуск или выпуск с экстренным исправлением; версия npm-пакета `@wailsio/runtime` всегда совпадает с версией выпуска Wails, а ночные сборки публикуются с dist-тегом `nightly`
- Шаблоны проектов фиксируют `@wailsio/runtime` на версии CLI, с помощью которой они были созданы
- Для Wails v3 требуется Go 1.26 или новее (Go 1.26.10 либо 1.27.3 и более поздние версии содержат исправление безопасности `crypto/x509`)

## Исправлено
- Исправить ошибку, из-за которой `wails3 update cli` не находил выпуски v3: он учитывал только теги v2 и устанавливал из неправильного пути
- Исправить сбой приложений Windows с `panic: SetTextColor failed` при отрисовке строки меню приложения в некоторых сборках Windows: теперь строка меню отрисовывается в допустимом контексте устройства или используется стандартная отрисовка, `UAHMENUITEM.Position` соответствует C `int`, а отрисовка меню никогда не вызывает панику (#6311)
- Обеспечить безопасность подготовленного обновления в ожидании перезапуска: скачивание или повторный перезапуск, начатые во время завершения работы, не могут удалить полезную нагрузку, которую ожидает установить вспомогательный процесс замены, или вступить с ним в конфликт за неё — от @julianstorer
- Исправить сбой приложений Windows в `applyDarkMode`, когда инициализация `pkg/w32` выгружала `uxtheme.dll`, пока указатели на его функции по порядковым номерам ещё использовались: теперь модуль остаётся загруженным на протяжении всего времени работы процесса (#6356)
- Проверять подготовленную полезную нагрузку перед заменой, а установленное приложение — перед удалением резервной копии, чтобы при повреждённом обновлении выполнялся откат, а не замена рабочего приложения, от @julianstorer

## Удалено
- Удалена интеграция с Crowdin после окончательной синхронизации переводов в [PR](https://github.com/wailsapp/wails/pull/5736), от @github-actions[bot]

## Безопасность
- Заблокированы зарезервированные сообщения `wails:` из дочерних фреймов другого источника в Linux (GTK3 и GTK4) путём внедрения защитного скрипта во все фреймы, что снижает риск уязвимости GHSA-2p5f-2gx8-mh8m
- Обновлены `github.com/klauspost/compress` до 1.18.7 (GHSA-259r-337f-4rfw) и `golang.org/x/image` до 0.45.0; сборка теперь выполняется с Go 1.26.10, чтобы включить исправление `crypto/x509` (GO-2026-6809)
- Обновлён `vite` до 8.2.4 или более поздней версии в пакете среды выполнения, мастере установки, примерах и шаблонах проектов (GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv)

### Известные ограничения
- **Linux:** защита дочерних фреймов не может помешать дочернему фрейму *того же источника* (например, iframe `srcdoc` или `about:blank`) обращаться к нативному обработчику сообщений через родительский фрейм, поскольку WebKitGTK не предоставляет обработчику сведения о фрейме. Фреймы другого источника и изолированные фреймы без `allow-same-origin` блокируются. Не встраивайте ненадёжное содержимое того же источника в окно Wails.
- **Android (экспериментальная поддержка):** мост JavaScript `window.wails` доступен каждому фрейму в WebView без проверки главного фрейма или источника. Содержимое встроенного iframe, в том числе содержимое другого источника, может вызывать связанные методы и отправлять зарезервированные сообщения `wails:`. Не встраивайте ненадёжное стороннее содержимое в WebView приложения. Поведение отличается от macOS и iOS, где сообщения `wails:` из дочерних фреймов игнорируются. Исправление запланировано для 3.0.1.
- **Все платформы:** вызовы связанных методов из дочерних фреймов передаются по HTTP в `/wails/runtime`, а не через нативный обработчик сообщений, поэтому защита их не ограничивает. Считайте, что любой встроенный iframe имеет доступ к связанным методам Go.

## v3.0.0-beta.34 - 2026-10-20

## Исправлено
- Исправлено: WebView в Windows не возвращал фокус при активации окна, а клавиша Tab оставалась внутри WebView вместо перехода к нативному окну и обратно (#3783)
- Исправлены структуры таблиц виртуальных методов COM `internal/webview2/pkg/webview2` (`ICoreWebView2_2`..`_27` и другие производные интерфейсы), в которых отсутствовала цепочка унаследованных методов, из-за чего вызовы попадали в неверные позиции таблицы (#6129)
- Исправлено: диалоги с вопросом в Windows игнорировали пользовательские подписи кнопок и никогда не вызывали обработчики кнопок. Теперь пользовательские кнопки отображаются через TaskDialog, а если comctl32 v6 недоступен — через MessageBox с позиционным сопоставлением обработчиков (#5794)
- Исправлено зависание `wails3 dev` на 60 с Vite: автоматическая проверка готовности фронтенда теперь использует TCP-запрос вместо HTTP GET, который отправлялся без заголовка `Accept` и получал ответ 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Меню системного трея продолжают работать при обновлении в открытом состоянии — [PR](https://github.com/wailsapp/wails/pull/6301), от @taliesin-ai
- Добавлена корректная обработка временных ошибок COM WebView2 во время операций среды выполнения — [PR](https://github.com/wailsapp/wails/pull/6302), от @leaanthony

## v3.0.0-beta.33 - 2026-10-19

## Добавлено
- Генератор привязок теперь создаёт модуль `bindings/events`, экспортирующий типизированную константу имени для каждого доступного для обнаружения события, зарегистрированного через `application.RegisterEvent` с константным именем, известным во время компиляции (например, `UserLogin` для `user:login`). Константу можно использовать с `Events.On`/`Emit` без плагина Vite (#5645)

## Исправлено
- Исправлено: диалоги с вопросом в Windows игнорировали пользовательские подписи кнопок и никогда не вызывали обработчики кнопок. Теперь пользовательские кнопки отображаются через TaskDialog, а если comctl32 v6 недоступен — через MessageBox с позиционным сопоставлением обработчиков (#5794)
- Исправлено зависание `wails3 dev` на 60 с Vite: автоматическая проверка готовности фронтенда теперь использует TCP-запрос вместо HTTP GET, который отправлялся без заголовка `Accept` и получал ответ 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Меню системного трея продолжают работать при обновлении в открытом состоянии — [PR](https://github.com/wailsapp/wails/pull/6301), от @taliesin-ai
- Добавлена корректная обработка временных ошибок COM WebView2 во время операций среды выполнения — [PR](https://github.com/wailsapp/wails/pull/6302), от @leaanthony

## v3.0.0-beta.32 - 2026-10-18

## Добавлено
- Добавлены нативные заголовки разделов меню macOS для macOS 14+ — см. [сайт документации](https://v3.wails.io) (обновлены страницы 20) в [PR](https://github.com/wailsapp/wails/pull/6297), от @leaanthony

## Исправлено
- Метод SetRelativePosition теперь позиционирует относительно рабочей области экрана в macOS — [PR](https://github.com/wailsapp/wails/pull/6299), от @leaanthony
- Вспомогательный процесс обновления повторно запускает неизменённое целевое приложение после истечения времени ожидания родительского процесса или неудачной замены — [PR](https://github.com/wailsapp/wails/pull/6298), от @taliesin-ai
- Для версионированной документации закреплена версия pnpm 10 — [PR](https://github.com/wailsapp/wails/pull/6293), от @taliesin-ai

## v3.0.0-beta.31 - 2026-10-18

## Добавлено
- Добавлен экспорт подпути `@wailsio/runtime/cancellable`, чтобы `CancellablePromise` можно было импортировать без выполнения корневой точки входа среды выполнения (#6153)

## Изменено
- Теперь требуется Go 1.26 или более поздней версии: модуль v3 теперь объявляет `go 1.26.0`, и проекты, созданные с помощью `wails3 init` и `wails3 service init`, объявляют тот же минимум

## Исправлено
- Исправлено изменение вывода генератора привязок при сборке с Go 1.27, где типы `encoding/json`, такие как `Marshaler`, `Unmarshaler` и `RawMessage`, становились псевдонимами типов `encoding/json/v2` и `encoding/json/jsontext`: теперь они снова отображаются как `any`, а модели для `encoding/json/v2` не создаются (#6292)
- Добавлена проверка URL удалённых шаблонов, переданных в `wails3 init -t`, перед клонированием: принимаются только URL `https://`, `http://`, `ssh://`, `git://` и `user@host:path`; значения передаются git после `--`, а при клонировании git разрешено использовать только эти протоколы передачи
- Значения заголовка `x-wails-window-id` вне допустимого диапазона теперь отклоняются транспортом среды выполнения HTTP, а не молча усекаются до идентификатора другого окна
- Распаковка архивов средством обновления защищена от zip-slip: пути записей теперь проверяются с помощью `filepath.IsLocal` и дополнительной проверки нахождения объединённого пути внутри целевого каталога; добавлены тесты для вложенных `..`, обратных косых черт, абсолютных путей и записей, выполняющих запись через символьные ссылки, в архивах zip и tar.gz
- Исправлен сбой при запуске AppImage из `wails3 generate appimage` вне Debian и Ubuntu с сообщением «Unable to spawn a new child process: … WebKitNetworkProcess». Встроенный WebKit теперь запускает вспомогательные процессы, поставляемые в составе AppImage. Песочница bubblewrap WebKit остаётся включённой
- Исправлены ошибки компиляции сборок `-tags server` для Android и iOS, что исправляет [#6271](https://github.com/wailsapp/wails/issues/6271)
- Исправлен сбой замены при обновлении в Linux, когда временный каталог находится в другой файловой системе (например, в tmpfs `/tmp`): подготовка теперь выполняется рядом с установленным исполняемым файлом, а между файловыми системами данные копируются
- Исправлено: средство обновления пыталось заменить доступный только для чтения смонтированный образ AppImage вместо файла `.AppImage`
- Исправлен приём зарезервированных сообщений управления окном `wails:` из iframe на macOS и iOS. Теперь они принимаются только из главного фрейма окна; остальные сообщения из iframe по-прежнему поступают в `RawMessageHandler` вместе с информацией об их источнике
- Прекращена отправка `WindowDidUpdate` и `WindowWillUpdate` в macOS: `NSWindowDidUpdateNotification` срабатывает при каждом проходе цикла выполнения, поэтому страница получала около 120 пустых вызовов `evaluateJavaScript:` в секунду, а в macOS 27.0 каждый из них поглощал кратковременную активацию от щелчка пользователя — `navigator.clipboard.writeText()` и другие API, требующие активации, отклоняли реальные щелчки с ошибкой `NotAllowedError` ([#6170](https://github.com/wailsapp/wails/issues/6170)), автор @APshenkin
- Исправлен язык по умолчанию в Windows info.json на en-US (0409), чтобы Проводник показывал сведения о версии (#6210)
- Исправлено отключение средства визуализации Linux DMA-BUF на ноутбуках с дискретным GPU NVIDIA наряду со встроенным GPU AMD или Intel: теперь в Wayland оно отключается, только когда WebKit выполняет визуализацию на NVIDIA (приложение было перенесено на дискретный GPU либо на нём работает композитор), а не при любой загрузке модуля nvidia, из-за чего визуализация на встроенном GPU замедлялась до предела; в сеансах X11 сохранено прежнее поведение, поскольку используемый Xorg GPU для визуализации задаётся его конфигурацией
- Добавлена обработка сбоев переименования между устройствами в средстве обновления Linux с резервным поэтапным копированием в [PR](https://github.com/wailsapp/wails/pull/6135), автор @john-okeefe
- Исправлена ошибка, из-за которой `WebviewWindow.SetMenu` в Windows не перерисовывал строку меню, если для неё не была настроена тема; исправлено [#6103](https://github.com/wailsapp/wails/issues/6103)
- Исправлена паника `WebviewWindow.SetMenu(nil)` в Windows: теперь, как и в macOS, операция ничего не делает; исправлено [#6104](https://github.com/wailsapp/wails/issues/6104)
- Добавлено сохранение записей журнала изменений при конкурирующих отправках изменений благодаря повторным попыткам в [PR](https://github.com/wailsapp/wails/pull/6094), автор @leaanthony
- Исправлен сбой `go mod vendor` с ошибкой `pattern arm64/WebView2Loader.dll: no matching files found` на всех платформах: удалены встраивания, ссылавшиеся на двоичные файлы, которые никогда не поставлялись в составе модуля; исправлены [#5782](https://github.com/wailsapp/wails/issues/5782) и [#5376](https://github.com/wailsapp/wails/issues/5376) в [PR](https://github.com/wailsapp/wails/pull/6031), автор @Grantmartin2002
- Предотвращено повторное уничтожение окна в macOS в [PR](https://github.com/wailsapp/wails/pull/6159), автор @charlesbalderson
- Добавлена обработка WM_CONTEXTMENU, чтобы меню системного трея открывалось по щелчку правой кнопкой мыши, в [PR](https://github.com/wailsapp/wails/pull/6162), автор @wzk2025
- Исправлены все пять вариантов `InvokeSync*`, навсегда блокировавшие вызывающий код при панике отправленного обратного вызова; исправлено [#6107](https://github.com/wailsapp/wails/issues/6107) в [PR](https://github.com/wailsapp/wails/pull/6108), автор @RALIST
- Добавлено безопасное повторное создание уничтоженных подменю GTK3 в [PR](https://github.com/wailsapp/wails/pull/6236), автор @leaanthony
- Реализована прослойка перенаправления рабочих процессов Windows для вложенных рабочих процессов в [PR](https://github.com/wailsapp/wails/pull/6261), автор @taliesin-ai
- Исправлена ошибка, из-за которой `OpenFile()` в Windows игнорировал `PromptForMultipleSelection()` при `CanChooseDirectories(true)`: теперь пользователи могут выбрать несколько папок одновременно, как в macOS и Linux (#6263)
- Исправлена передача аргументов `double` по значению в методы COM WebView2 на windows/arm64 — проблема затрагивала `ZoomFactor`, `RasterizationScale`, `Expires` файлов cookie, а также масштаб, размер страницы и поля в параметрах печати (ранее сеттеры параметров печати возвращали `ErrDoubleArgUnsupported`): теперь затронутые вызовы проходят через трамплин `runtime.cgocall`, загружающий `d0`, чего не может сделать путь системных вызовов Go (golang.org/issue/62583), автор @APshenkin

## Безопасность
- Теперь серверный режим (`-tags server`) принимает вызовы `/wails/runtime` только как запросы `POST` с собственного источника сервера или источника, указанного в `ServerOptions.WebSocketOriginPatterns`; для других методов возвращается HTTP 405, а для других источников браузера — HTTP 403. При развёртывании за обратным прокси, который перезаписывает заголовок `Host`, добавьте публичный источник в `WebSocketOriginPatterns`
- Теперь серверный режим (`-tags server`), привязанный к адресу обратной петли (по умолчанию `localhost`, `127.0.0.0/8` или `::1`), отклоняет запросы с кодом HTTP 403, если заголовок `Host` не равен `localhost`, не является IP-литералом, хостом привязки или хостом, соответствующим `ServerOptions.WebSocketOriginPatterns`. Это блокирует атаки с перепривязкой DNS на `/wails/runtime` и конечные точки WebSocket
- Уведомления Windows больше не используют PowerShell как резервный вариант при сбое API всплывающих уведомлений WinRT; вместо этого возвращается ошибка

## v3.0.0-beta.30 — 2026-10-15

## Исправлено
- Исправлена попытка средства обновления заменить доступную только для чтения точку монтирования AppImage вместо файла `.AppImage`
- Прекращена отправка `WindowDidUpdate` и `WindowWillUpdate` в macOS: `NSWindowDidUpdateNotification` срабатывает при каждом проходе цикла выполнения, поэтому страница получала около 120 пустых вызовов `evaluateJavaScript:` в секунду, а в macOS 27.0 каждый из них поглощал кратковременную активацию от щелчка пользователя — `navigator.clipboard.writeText()` и другие API, требующие активации, отклоняли реальные щелчки с ошибкой `NotAllowedError` ([#6170](https://github.com/wailsapp/wails/issues/6170)), автор @APshenkin
- Исправлен язык по умолчанию в Windows info.json на en-US (0409), чтобы Проводник показывал сведения о версии (#6210)
- Добавлена обработка сбоев переименования между устройствами в средстве обновления Linux с резервным поэтапным копированием в [PR](https://github.com/wailsapp/wails/pull/6135), автор @john-okeefe
- Исправлена ошибка, из-за которой `WebviewWindow.SetMenu` в Windows не перерисовывал строку меню, если для неё не была настроена тема; исправлено [#6103](https://github.com/wailsapp/wails/issues/6103)
- Исправлена паника `WebviewWindow.SetMenu(nil)` в Windows: теперь, как и в macOS, операция ничего не делает; исправлено [#6104](https://github.com/wailsapp/wails/issues/6104)
- Добавлено сохранение записей журнала изменений при конкурирующих отправках изменений благодаря повторным попыткам в [PR](https://github.com/wailsapp/wails/pull/6094), автор @leaanthony
- Исправлен сбой `go mod vendor` с ошибкой `pattern arm64/WebView2Loader.dll: no matching files found` на всех платформах: удалены встраивания, ссылавшиеся на двоичные файлы, которые никогда не поставлялись в составе модуля; исправлены [#5782](https://github.com/wailsapp/wails/issues/5782) и [#5376](https://github.com/wailsapp/wails/issues/5376) в [PR](https://github.com/wailsapp/wails/pull/6031), автор @Grantmartin2002
- Предотвращено повторное уничтожение окна в macOS в [PR](https://github.com/wailsapp/wails/pull/6159), автор @charlesbalderson
- Добавлена обработка WM_CONTEXTMENU, чтобы меню системного трея открывалось по щелчку правой кнопкой мыши, в [PR](https://github.com/wailsapp/wails/pull/6162), автор @wzk2025
- Исправлены все пять вариантов `InvokeSync*`, навсегда блокировавшие вызывающий код при панике отправленного обратного вызова; исправлено [#6107](https://github.com/wailsapp/wails/issues/6107) в [PR](https://github.com/wailsapp/wails/pull/6108), автор @RALIST

## v3.0.0-beta.29 — 2026-10-13

## Добавлено
- Теперь глобальные сочетания клавиш в macOS учитывают раскладку клавиатуры: в раскладках, отличных от QWERTY (AZERTY, QWERTZ, Dvorak, ...), буквенные акселераторы через `UCKeyTranslate` привязываются к клавише с соответствующим обозначением, а не к фиксированной физической позиции QWERTY.

## Исправлено
- Добавлена инициализация нативного меню перед его установкой в Linux в [PR](https://github.com/wailsapp/wails/pull/6233), автор @abhicodes-007
- Устранены регрессии готовности macOS к общедоступному выпуску в обработке выходных данных и мастере настройки в [PR](https://github.com/wailsapp/wails/pull/6237), автор @taliesin-ai
- Исправлены обёртки COM WebView2 `AddRef`/`Release` в Windows, которые сообщали об устаревшей последней ошибке потока как о сбое. Теперь они возвращают счётчик ссылок, поэтому освобождение больше не может аварийно завершить приложение или привести к сбою запроса ресурса (#3931)
- Исправлен вызов слота `AddRef` из `ICoreWebView2Deferral.Release`, из-за которого в Windows утекала отсрочка каждого отложенного запроса ресурса (#3931)

## v3.0.0-beta.28 — 2026-10-05

## Исправлено
- Выходные данные приложения macOS сохранены вне поставщиков файлов в [PR](https://github.com/wailsapp/wails/pull/6225), автор @taliesin-ai
- Завершена реализация мастера настройки и добавлен показ сведений о системе в [PR](https://github.com/wailsapp/wails/pull/6224), автор @taliesin-ai
- Улучшены обработка ошибок восстановления WebView2 и очистка ресурсов в Windows в [PR](https://github.com/wailsapp/wails/pull/6222), автор @leaanthony
- Улучшено восстановление WebView2 в Windows, чтобы предотвратить аварийное завершение приложения и зависание контроллеров, в [PR](https://github.com/wailsapp/wails/pull/6220) от @taliesin-ai
- Исправлены рекомендации по авторизации доступа к мультимедиа в macOS в [PR](https://github.com/wailsapp/wails/pull/6221) от @taliesin-ai
- Исправлена ошибка, из-за которой щелчок правой кнопкой мыши по значку в системном трее Linux не только открывал меню, но и запускал обработчик щелчка, поэтому привязанное окно переключалось при каждом таком щелчке (#6018)
- Обработчики обратных вызовов WebView2 сохраняются до их освобождения нативным кодом в Windows в [PR](https://github.com/wailsapp/wails/pull/6184) от @taliesin-ai

## v3.0.0-beta.27 — 2026-10-01

## Исправлено
- Исправлено игнорирование параметра `Permissions` в macOS 12 и более поздних версиях: теперь реализован метод захвата мультимедиа `WKUIDelegate`, поэтому `PermissionAllow` и `PermissionDeny` применяются к запросам доступа к камере и микрофону так же, как в Linux и Windows. По-прежнему требуются `NSCameraUsageDescription` / `NSMicrophoneUsageDescription`, а при работе в песочнице — соответствующие разрешения на устройства

## v3.0.0-beta.26 — 2026-09-25

## Добавлено
- Средство наблюдения выполняет HTTP-проверки готовности сервера разработки фронтенда в [PR](https://github.com/wailsapp/wails/pull/6150) от @atterpac

## Исправлено
- Добавлено восстановление после сбоев процессов WebView2 вместо отображения пустого окна в [PR](https://github.com/wailsapp/wails/pull/6002) от @taliesin-ai

## v3.0.0-beta.25 — 2026-09-22

## Исправлено
- При завершении основного приложения теперь останавливаются `wails3 dev` и его фоновые процессы (#6048)

## v3.0.0-beta.24 — 2026-09-20

## Безопасность
- Разрешения GITHUB_TOKEN в рабочих процессах сборки и публикации ограничены в соответствии с принципом наименьших привилегий
- Устранены обход путей и выходы через символические ссылки в промежуточном ПО для ресурсов примера screen
- Установка зависимостей мастером настройки сохранена в рабочем состоянии, но ограничена одобренными командами менеджеров пакетов

## v3.0.0-beta.23 — 2026-09-16

## Добавлено

- Добавлена возможность переопределять идентификатор приложения GTK в Linux в [PR](https://github.com/wailsapp/wails/pull/5972) от @overlordtm

## Исправлено

- Механизм единственного экземпляра в Linux теперь регистрирует имя в шине под собственным идентификатором приложения в [PR](https://github.com/wailsapp/wails/pull/5973) от @overlordtm
- Прерванные запросы пользовательских схем в Linux теперь отменяются в [PR](https://github.com/wailsapp/wails/pull/6126) от @leaanthony

## v3.0.0-beta.22 — 2026-09-14

## Добавлено

- Завершены переводы и перенесены материалы сообщества — см. [сайт документации](https://v3.wails.io) (обновлено страниц: 1812) — в [PR](https://github.com/wailsapp/wails/pull/6122) от @leaanthony
- Документированы встраиваемые кнопки Wails различных стилей и размеров в [PR](https://github.com/wailsapp/wails/pull/6120) от @leaanthony

## v3.0.0-beta.21 — 2026-09-13

## Добавлено

- Документация Wails v3 теперь обслуживается с помощью M-Press в [PR](https://github.com/wailsapp/wails/pull/6116) от @leaanthony

## Исправлено

- Добавлен разбор JSON-значений slug из метаданных MPD для создания журнала изменений в [PR](https://github.com/wailsapp/wails/pull/6118) от @leaanthony
- После сбоев резервного копирования средство обновления очищает вспомогательные переменные окружения и повторно запускает исходную целевую программу в [PR](https://github.com/wailsapp/wails/pull/6080) от @cnmax
- Обработчик сигналов по умолчанию теперь запускается во время App.Run в [PR](https://github.com/wailsapp/wails/pull/6098) от @leaanthony
- Меню Windows теперь обрабатывает nil-меню, освобождает заменённые ресурсы и перерисовывает строку меню в [PR](https://github.com/wailsapp/wails/pull/6112) от @taliesin-ai
- Восстановлена упаковка MSIX для новых проектов с общей конфигурацией YAML в [PR](https://github.com/wailsapp/wails/pull/6115) от @leaanthony
- Исправлен сбой загрузки сгенерированных привязок JavaScript и TypeScript, когда создатели обобщённых моделей ссылались на вспомогательные объявления, расположенные ниже; также предотвращено переполнение стека при создании взаимозависимых обобщённых моделей (#6062)

## v3.0.0-beta.20 — 2026-09-10

## Изменено

- Ссылки на демонстрационный проект Clave обновлены и теперь указывают на текущие сайт и репозиторий в [PR](https://github.com/wailsapp/wails/pull/6082) от @01xR4in

## Исправлено

- Прерванные запросы ресурсов в Windows, включая запросы рабочих процессов, теперь отменяются, а обработчики keepalive сохраняются при навигации. На платформах Apple контексты нативных запросов передаются через оболочку приложения. (#5963, #5969)
- Записи журнала изменений сохраняются при конкурирующих отправках изменений благодаря повторным попыткам в [PR](https://github.com/wailsapp/wails/pull/6094) от @leaanthony
- Исправлен сбой `go mod vendor` с ошибкой `pattern arm64/WebView2Loader.dll: no matching files found` на всех платформах: удалены встраивания, ссылавшиеся на двоичные файлы, которые никогда не поставлялись в составе модуля; тем самым исправлены [#5782](https://github.com/wailsapp/wails/issues/5782) и [#5376](https://github.com/wailsapp/wails/issues/5376), в [PR](https://github.com/wailsapp/wails/pull/6031) от @Grantmartin2002

## Удалено

- Удалена поддержка нативного загрузчика WebView2, заменённого загрузчиком на чистом Go. Удалены встроенные двоичные файлы `WebView2Loader.dll` и зависимость `github.com/jchv/go-winloader`. Тег сборки `native_webview2loader` по-прежнему принимается и больше не вызывает ошибку, но не влияет на сборки v3, в [PR](https://github.com/wailsapp/wails/pull/6031) от @Grantmartin2002
- Из руководства по API macOS удалены неиспользуемые теги сборки и параметр FPS в [PR](https://github.com/wailsapp/wails/pull/6097) от @leaanthony

## v3.0.0-beta.19 — 2026-09-09

## Добавлено

- Доступ к закрытым API macOS теперь предоставляется только при явном включении с помощью тегов сборки — см. [документацию](https://v3.wails.io/features/browser/integration), [документацию](https://v3.wails.io/features/environment/info), [документацию](https://v3.wails.io/features/windows/basics), [документацию](https://v3.wails.io/features/windows/frameless), [документацию](https://v3.wails.io/features/windows/notch-windows), [документацию](https://v3.wails.io/features/windows/options), [документацию](https://v3.wails.io/guides/build/macos), [документацию](https://v3.wails.io/guides/build/private-macos-apis) и [документацию](https://v3.wails.io/reference/overview) — в [PR](https://github.com/wailsapp/wails/pull/6087) от @leaanthony

## Исправлено

- Запросы к среде выполнения размером более 64 МиБ теперь отклоняются с кодом HTTP 413 в [PR](https://github.com/wailsapp/wails/pull/6091) от @leaanthony

## Безопасность

- Усилена защита источников MCP и удалённого доступа с помощью аутентификации по токену в [PR](https://github.com/wailsapp/wails/pull/6092) от @leaanthony

## v3.0.0-beta.18 — 2026-09-08

## Исправлено

- Устранена утечка памяти Calloc в Linux и Darwin за счёт использования получателей-указателей в [PR](https://github.com/wailsapp/wails/pull/6083) от @4RH1T3CT0R7

## v3.0.0-beta.17 — 2026-09-06

## Исправлено

- Windows: ошибочный или nil-результат `GetRequest` в обработчике WebResourceRequested больше не завершает процесс (`log.Fatal` / паника из-за разыменования nil) — вместо этого запрос отбрасывается, а событие записывается в журнал, в [PR](https://github.com/wailsapp/wails/pull/6006) от @midagedev

## v3.0.0-beta.16 — 2026-08-29

## Изменено

- Запрашивать пароль для нотариального заверения в новом окне терминала — [PR](https://github.com/wailsapp/wails/pull/6029) от @leaanthony

## Исправлено

- Правильно обрабатывать типы нажатий на значок в области уведомлений macOS — [PR](https://github.com/wailsapp/wails/pull/5919) от @ChewbaccaCookie
- Перед обновлением CI удаляет неиспользуемые репозитории Microsoft apt — [PR](https://github.com/wailsapp/wails/pull/6041) от @Grantmartin2002

## v3.0.0-beta.15 — 2026-08-27

## Исправлено

- Увеличить время ожидания встраивания WebView2 до 60 секунд — [PR](https://github.com/wailsapp/wails/pull/6043) от @Grantmartin2002

## v3.0.0-beta.14 — 2026-08-26

## Исправлено

- Правильно именовать сочетания клавиш Control с буквами в macOS — [PR](https://github.com/wailsapp/wails/pull/6032) от @taliesin-ai
- Исправить значки ICO в области уведомлений и учитывать тему панели задач Windows — [PR](https://github.com/wailsapp/wails/pull/6016) от @nik9play

## v3.0.0-beta.13 — 2026-08-25

## Исправлено

- Продолжать обрабатывать работу главного потока в macOS во время выполнения модального цикла — [PR](https://github.com/wailsapp/wails/pull/6026) от @leaanthony
- Сделать операции безопасного хранилища на мобильных платформах допускающими ошибку и безопасно запрещающими доступ при сбое — [PR](https://github.com/wailsapp/wails/pull/5923) от @mortenolsrud
- Выполнять обработчики событий приложения, даже если не зарегистрирован ни один слушатель — [PR](https://github.com/wailsapp/wails/pull/5999) от @archy-rock3t-cloud
- Исправить опечатки в комментариях и локализованной документации — [PR](https://github.com/wailsapp/wails/pull/6023) от @haoku123
- Удалить предварительно скомпилированные двоичные файлы macOS, зафиксированные в `v3/examples`, — [PR](https://github.com/wailsapp/wails/pull/6025) от @4RH1T3CT0R7

## v3.0.0-beta.12 — 2026-08-21

## Добавлено

- Добавить окна уведомлений для выреза экрана macOS с примером управления жизненным циклом и телеметрии — см. [документацию](https://v3.wails.io/features/windows/notch-windows); [PR](https://github.com/wailsapp/wails/pull/6010) от @leaanthony
- Добавить поддержку окон NSPanel в macOS с новыми параметрами и нативной интеграцией — см. [документацию](https://v3.wails.io/features/windows/options); [PR](https://github.com/wailsapp/wails/pull/6008) от @leaanthony

## Исправлено

- Предотвратить зависание SQLite Prepare при параллельных вызовах — [PR](https://github.com/wailsapp/wails/pull/5998) от @archy-rock3t-cloud

## v3.0.0-beta.11 — 2026-08-20

## Удалено

- Удалить из документации устаревший трекер реализации — [PR](https://github.com/wailsapp/wails/pull/6005) от @leaanthony

## v3.0.0-beta.10 — 2026-08-19

## Исправлено

- Исправить потерю аргументов запуска для пользовательских протоколов и ассоциаций файлов хостом GTK4 в Linux — [PR](https://github.com/wailsapp/wails/pull/6000) от @midagedev
- Правильно обрабатывать удалённые строки и исправления в том же исходном файле при проверке журнала изменений — [PR](https://github.com/wailsapp/wails/pull/5993) от @taliesin-ai

## v3.0.0-beta.9 — 2026-08-16

## Добавлено

- Добавить защищённый сервер wails3 mcp для управления проектами с помощью агентов — [PR](https://github.com/wailsapp/wails/pull/5896) от @leaanthony
- Добавить документацию по моделям в привязках — см. [документацию](https://v3.wails.io/features/bindings/models); [PR](https://github.com/wailsapp/wails/pull/5988) от @taliesin-ai
- Добавить поддержку установки через rpm-ostree в атомарных системах Linux — [PR](https://github.com/wailsapp/wails/pull/5987) от @leaanthony
- Добавить нативное еженедельное создание и публикацию графика истории звёзд — см. [документацию](https://v3.wails.io/credits), [документацию](https://v3.wails.io/de/credits), [документацию](https://v3.wails.io/fr/credits), [документацию](https://v3.wails.io/id/credits), [документацию](https://v3.wails.io/ja/credits), [документацию](https://v3.wails.io/ko/credits), [документацию](https://v3.wails.io/pt/credits), [документацию](https://v3.wails.io/ru/credits), [документацию](https://v3.wails.io/zh-cn/credits) и [документацию](https://v3.wails.io/zh-tw/credits); [PR](https://github.com/wailsapp/wails/pull/5986) от @leaanthony
- Добавить пакет mac только для Darwin для разрешения ресурсов пакета приложения — см. [документацию](https://v3.wails.io/guides/build/macos); [PR](https://github.com/wailsapp/wails/pull/5965) от @leaanthony
- Добавить демонстрационную страницу Condui и запись в указателе — см. [документацию](https://v3.wails.io/community/showcase/condui) и [документацию](https://v3.wails.io/community/showcase); [PR](https://github.com/wailsapp/wails/pull/5962) от @mgueregath
- Добавить демонстрационную страницу Redis Viewer со снимками экрана и ссылкой на проект — см. [документацию](https://v3.wails.io/community/showcase) и [документацию](https://v3.wails.io/community/showcase/redisviewer); [PR](https://github.com/wailsapp/wails/pull/5984) от @redisviewer

## Изменено

- Обновить флаги приложения GTK для Linux до G<em>APPLICATION</em>NON_UNIQUE — [PR](https://github.com/wailsapp/wails/pull/5971) от @overlordtm
- Записывать события отсутствующих окон на уровне отладки, а не предупреждения — [PR](https://github.com/wailsapp/wails/pull/5914) от @julianstorer

## Исправлено

- Разрешить зарегистрированным сочетаниям клавиш macOS иметь приоритет над веб-представлением — [PR](https://github.com/wailsapp/wails/pull/5902) от @julianstorer
- Исправить неработающие ссылки на боковой панели документации — [PR](https://github.com/wailsapp/wails/pull/5937) от @northes
- Отменять контексты запросов ресурсов в macOS и iOS, когда WebKit прерывает соответствующую задачу пользовательской схемы (#5963)
- Обрабатывать сообщение WindowSetFullscreenButtonEnabled — [PR](https://github.com/wailsapp/wails/pull/5976) от @archy-rock3t-cloud
- Импортировать Fragment в шаблоне preact-ts, чтобы устранить сбой сборки, — [PR](https://github.com/wailsapp/wails/pull/5979) от @haoku123
- Предотвратить сбой устаревших служебных приложений GTK3, когда поиск экранов выполняется до появления активного окна или доступного дисплея (#5966)
- Разрешить выпускам с явно заданной версией продолжать выполнение, когда раздел невыпущенных изменений пуст (#5977)

## Безопасность

- Обновить lock-файлы nanoid для веб-сайта до исправленной версии 3.3.18, чтобы устранить предупреждения безопасности, — [PR](https://github.com/wailsapp/wails/pull/5985) от @taliesin-ai

## v3.0.0-beta.8 — 2026-08-12

## Добавлено

- Добавить создание URL-адресов документации в автоматические записи журнала изменений — [PR](https://github.com/wailsapp/wails/pull/5957) от @taliesin-ai
- Добавить Streams — двунаправленные потоки байтов между Go и JavaScript с программной моделью WebSocket и без прослушивающего сокета. Объявите поток в Go с помощью `app.HandleStream(name, handler)` и подключитесь к нему из фронтенда с помощью `Stream(name)`, который возвращает объект в форме `WebSocket`. Данные Go→JS передаются через сервер ресурсов с помощью одного удерживаемого опроса на окно, а JS→Go — обычным запросом POST; TCP-порт не привязывается, и через `evaluateJavaScript` ничего не проходит. В серверных сборках (`-tags server`) тот же обработчик вместо этого обслуживается через настоящий WebSocket, поэтому код приложения одинаков во всех вариантах сборки. Автор: @leaanthony
- Переместить запись журнала изменений о почтовом ящике в раздел невыпущенных изменений — [PR](https://github.com/wailsapp/wails/pull/5935) от @leaanthony

## Изменено

- Обновить автоматическое создание боковой панели документации и выведение типа автора блога — [PR](https://github.com/wailsapp/wails/pull/5938) от @leaanthony

## Исправлено

- При инициализации WebView2 используются крайний срок и цикл обработки сообщений — [PR](https://github.com/wailsapp/wails/pull/5952) от @leaanthony
- Тест файлов cookie WebView2 пропускается в CI, если он не включён явно, а его выполнение привязывается к текущему потоку ОС, в [PR](https://github.com/wailsapp/wails/pull/5951) от @leaanthony
- Построители меню Windows восстанавливают идентификаторы команд для родительских элементов подменю в [PR](https://github.com/wailsapp/wails/pull/5944) от @gilad-ch
- Официальный образ для кросс-компиляции приведён в соответствие с базовым уровнем поддержки GTK 4.14+ в Linux (#5928)
- Проект Xcode для iOS настроен на сохранение унаследованных флагов компоновщика и добавление -ObjC в [PR](https://github.com/wailsapp/wails/pull/5915) от @mortenolsrud
- Устранено чрезмерное пересоздание TCP-соединений в прокси ресурсов `wails3 dev` для крупных фронтендов, которое могло исчерпать эфемерные порты хоста и привести к сбою посторонних процессов с ошибкой `EADDRNOTAVAIL`
- JavaScript событий для каждого окна помещается в очередь для упорядоченной отправки и обратного давления в [PR](https://github.com/wailsapp/wails/pull/5934) от @leaanthony

## Удалено

- Удалён конвейер выпуска двоичных файлов для настольных платформ: выпуски v3 теперь содержат только теги, а CLI `wails3` устанавливается с помощью `go install`. Удалены `release-v3.yml` и запускавший его этап ночной сборки в [PR](https://github.com/wailsapp/wails/pull/5946) от @leaanthony

## v3.0.0-beta.7 — 2026-08-11

## Добавлено

- Добавлена настройка автовоспроизведения в macOS, позволяющая отключить требование действия пользователя для воспроизведения мультимедиа, в [PR](https://github.com/wailsapp/wails/pull/5512) от @Eyalm321
- Запись об изменении почтового ящика перенесена в раздел «Не выпущено» в [PR](https://github.com/wailsapp/wails/pull/5935) от @leaanthony

## Изменено

- Для более плавной анимации масштабирования в macOS теперь используется CADisplayLink или NSTimer в [PR](https://github.com/wailsapp/wails/pull/5945) от @savely-krasovsky

## Исправлено

- Проект Xcode для iOS настроен на сохранение унаследованных флагов компоновщика и добавление -ObjC в [PR](https://github.com/wailsapp/wails/pull/5915) от @mortenolsrud
- Устранено чрезмерное пересоздание TCP-соединений в прокси ресурсов `wails3 dev` для крупных фронтендов, которое могло исчерпать эфемерные порты хоста и привести к сбою посторонних процессов с ошибкой `EADDRNOTAVAIL`
- JavaScript событий для каждого окна помещается в очередь для упорядоченной отправки и обратного давления в [PR](https://github.com/wailsapp/wails/pull/5934) от @leaanthony

### Добавлено

- Реализован универсальный асинхронный почтовый ящик FIFO для упорядоченной доставки событий в [PR](https://github.com/wailsapp/wails/pull/5851) от @savely-krasovsky и @DevLumuz

## v3.0.0-beta.6 — 2026-08-09

## Добавлено

- Реализованы ограниченное хранилище на стороне хоста для событий чрезмерного размера и их упорядоченная доставка в JavaScript в [PR](https://github.com/wailsapp/wails/pull/5930) от @leaanthony
- Реализовано подпрыгивание значка в Dock macOS при мигании окна в [PR](https://github.com/wailsapp/wails/pull/5921) от @julianstorer

## Исправлено

- Сервер ресурсов сохраняет ошибки определения типа содержимого и незаписанные префиксы при сбросе буфера в [PR](https://github.com/wailsapp/wails/pull/5931) от @leaanthony
- Предотвращён аварийный сбой приложений macOS при замене меню приложения из обратного вызова Wails
- Исправлены нечитаемые нативные меню в Windows 10 1809 / Windows Server 2019 (сборка 17763). Экспортируемые функции uxtheme для тёмного режима были доступны только начиная со сборки 18334, поэтому явное включение тёмного режима на уровне приложения никогда не выполнялось на этих хостах: фон меню окрашивался в тёмный цвет, но Windows продолжала отрисовывать текст меню в светлой теме, из-за чего тёмный текст оказывался на тёмном фоне. Порядковые номера существуют начиная с 17763, поэтому теперь условие соответствует этому номеру сборки.
- Исправлен вызов `GetDeviceCaps` вместо `GetStockObject` в `w32.GetStockObject`, из-за которого для каждого стандартного объекта возвращалось 0.
- Улучшены обработка ошибок и отчётность при загрузке установщика WebView2 в [PR](https://github.com/wailsapp/wails/pull/5924) от @jannskiee

## v3.0.0-beta.5 — 2026-08-07

## Исправлено

- Активация приложения macOS теперь учитывает политику активации только для обычных приложений в [PR](https://github.com/wailsapp/wails/pull/5897) от @julianstorer
- Добавлена защита от неинициализированных окон GTK в сборках Linux в [PR](https://github.com/wailsapp/wails/pull/5898) от @julianstorer
- Для окон WebKit в Linux перед загрузкой URL явно задаётся непрозрачный цвет фона в [PR](https://github.com/wailsapp/wails/pull/5899) от @julianstorer

## v3.0.0-beta.4 — 2026-08-05

## Изменено

- Задачи сборки Android теперь по умолчанию используют arm64, а deploy-emulator выбирает архитектуру хоста, в [PR](https://github.com/wailsapp/wails/pull/5890) от @mortenolsrud

## Исправлено

- Состояние масштабирования окна macOS сохраняется при перетаскивании, а объём анимации уменьшен в [PR](https://github.com/wailsapp/wails/pull/5900) от @leaanthony
- Исправлена сборка Windows в серверном режиме путём добавления `!server` в ограничение сборки `webview_window_windows_nonclient.go`

## v3.0.0-beta.3 — 2026-08-03

## Добавлено

- Завершение бета-проверки этапа 10 задокументировано в подробностях реализации в [PR](https://github.com/wailsapp/wails/pull/5881) от @leaanthony

## Исправлено

- Дескриптор окна передаётся в API тёмного режима Windows, а аргументы проверяются, в [PR](https://github.com/wailsapp/wails/pull/5877) от @leaanthony
- Определение состояния кнопок строки заголовка для окон macOS без рамки централизовано в [PR](https://github.com/wailsapp/wails/pull/5870) от @taliesin-ai
- Предотвращено отображение нечитаемого текста нативного меню, когда приложение Windows запрашивает тёмный режим, а тема приложений Windows остаётся светлой. Теперь меню использует соответствующий светлый нативный фон, пока Windows не сможет отрисовать светлый текст меню для тёмного режима.
- Исправлены нечитаемые нативные меню в Windows 10 1809 / Windows Server 2019 (сборка 17763). Экспортируемые функции uxtheme для тёмного режима были доступны только начиная со сборки 18334, поэтому явное включение тёмного режима на уровне приложения никогда не выполнялось на этих хостах: фон меню окрашивался в тёмный цвет, но Windows продолжала отрисовывать текст меню в светлой теме, из-за чего тёмный текст оказывался на тёмном фоне. Порядковые номера существуют начиная с 17763, поэтому теперь условие соответствует этому номеру сборки.

## v3.0.0-beta.2 — 2026-08-02

## Изменено

- v3 переведён из альфа-версии в бета-версию
- Задокументированы интеллектуальные значения по умолчанию для системного трея и автоматическое скрытие всплывающего меню; добавлено регрессионное покрытие выбора обработчика щелчка (#5840).
- Средство обновления через GitHub по умолчанию исключает ресурсы установщика Windows в [PR](https://github.com/wailsapp/wails/pull/5861) от @leaanthony
- Добавлена поддержка скруглённых и прямых углов, а также углов с настраиваемым радиусом для окон macOS без рамки в [PR](https://github.com/wailsapp/wails/pull/5866) от @leaanthony

## Исправлено

- Добавлено получение актуальных размеров окон GTK4 и отправка событий изменения размера, развёртывания, свёртывания и перехода в полноэкранный режим с настроенной поверхности (#5830).
- Исправлен аварийный сбой WebKit в Linux при отправке Blob или FormData в запросах fetch в [PR](https://github.com/wailsapp/wails/pull/5854) от @taliesin-ai
- При отсутствии заголовков Blob/FormData прослойка fetch передаёт undefined в [PR](https://github.com/wailsapp/wails/pull/5865) от @leaanthony

## v3.0.0-alpha2.122 — 2026-08-01

## Добавлено

## Изменено

- Добавлена поддержка скруглённых и прямых углов, а также углов с настраиваемым радиусом для окон macOS без рамки в [PR](https://github.com/wailsapp/wails/pull/5866) от @leaanthony

## Исправлено

- Fetch-шим передаёт undefined для отсутствующих заголовков Blob/FormData в [PR](https://github.com/wailsapp/wails/pull/5865) от @leaanthony

## v3.0.0-alpha2.121 — 2026-07-31

## Добавлено

- Добавлена поддержка упаковки macOS-приложений в DMG с новыми параметрами и задачами сборки в [PR](https://github.com/wailsapp/wails/pull/5857) от @leaanthony

## Изменено

- Средство обновления через GitHub теперь по умолчанию исключает артефакты установщика Windows в [PR](https://github.com/wailsapp/wails/pull/5861) от @leaanthony

## Исправлено

- Исправлен сбой Linux WebKit при отправке Blob или FormData в запросах fetch в [PR](https://github.com/wailsapp/wails/pull/5854) от @taliesin-ai

## v3.0.0-alpha2.120 — 2026-07-31

## Добавлено

- Реализованы действия разворачивания или сворачивания окна по двойному щелчку по строке заголовка в macOS в [PR](https://github.com/wailsapp/wails/pull/5853) от @taliesin-ai

## Изменено

- Руководство по сервису QR обновлено для использования NewServiceWithOptions, также добавлены отступы в [PR](https://github.com/wailsapp/wails/pull/5849) от @jeongkyu

## Исправлено

- Обеспечена отзывчивость WKWebView при масштабировании в macOS в [PR](https://github.com/wailsapp/wails/pull/5856) от @leaanthony
- Исправлены запросы размера окна в GTK4 и отправка событий изменения размера, разворачивания, сворачивания и перехода в полноэкранный режим из настроенного `GdkSurface`.

## v3.0.0-alpha2.119 — 2026-07-27

## Исправлено

- В документацию на нескольких языках добавлены архитектурные схемы в [PR](https://github.com/wailsapp/wails/pull/5833) от @taliesin-ai

## v3.0.0-alpha2.118 — 2026-07-26

## Добавлено

- Добавлены пути по умолчанию для входных и выходных файлов генерации значков в [PR](https://github.com/wailsapp/wails/pull/5825) от @taliesin-ai
- Исходные модули-точки входа добавлены в sideEffects файла package.json пакета среды выполнения в [PR](https://github.com/wailsapp/wails/pull/5797) от @savely-krasovsky
- В руководство для участников проекта добавлен раздел о лицензии и происхождении материалов в [PR](https://github.com/wailsapp/wails/pull/5816) от @taliesin-ai

## Исправлено

- Для удаления радиуса скругления границ применён ограниченный областью действия CSS для безрамочных окон GTK4 в [PR](https://github.com/wailsapp/wails/pull/5800) от @savely-krasovsky
- В Windows обеспечена корректная обработка сбоев получения позиции курсора для всплывающих меню и перечисления экранов в [PR](https://github.com/wailsapp/wails/pull/5789) от @wayneforrest
- Диалог открытия файла в macOS теперь правильно фильтрует расширения и проверяет допустимость файлов по суффиксу в [PR](https://github.com/wailsapp/wails/pull/5678) от @phergul
- Инициализация тёмного режима в Windows защищена от вызовов nil API в [PR](https://github.com/wailsapp/wails/pull/5793) от @roachadam
- Исправлена ошибка сборки средства обновления для 32-разрядной платформы: константа `maxArchiveTotalSize` (2 ГиБ) переполняла платформенный тип `int` при передаче в `fmt.Errorf` на `GOARCH=386`. Теперь ей явно назначен тип `int64`.
- Исправлена паника из-за nil-указателя при запуске, когда окно использует тёмную (или системную тёмную) строку заголовка в сборках для Windows, которые не загружают API тёмного режима uxtheme, например Windows 10 1809 / Windows Server 2019 (сборка 17763). Вызовы `AllowDarkModeForWindow` при настройке темы окна теперь защищены от nil аналогично уже используемой защите в `w32.SetMenuTheme`.

## v3.0.0-alpha2.117 — 2026-07-08

## Добавлено

- Реализована пользовательская логика hit-test для неклиентских областей в Windows в [PR](https://github.com/wailsapp/wails/pull/5462) от @savely-krasovsky

## Изменено

- Обнаружение масштаба монитора в WebView2 настроено с учётом UseVisualHosting в [PR](https://github.com/wailsapp/wails/pull/5761) от @wayneforrest

## v3.0.0-alpha2.116 — 2026-07-07

## Добавлено

- Документация FAQ обновлена с акцентом на функции и рекомендации Wails v3 в [PR](https://github.com/wailsapp/wails/pull/5763) от @taliesin-ai

## v3.0.0-alpha2.115 — 2026-07-06

## Исправлено

- Исправлена ошибка, из-за которой `Menu.Update()` не перестраивал нативное меню в GTK4 Linux (#5659; независимо диагностировано и исправлено @puneetdixit200 в #5539)
- Исправлен сбой при перечислении экранов macOS после изменения конфигурации дисплеев: строки идентификаторов и имён экранов теперь копируются, а количество фиксируется в снимке (#5565; независимо диагностировано и исправлено @x-haose в #5584)
- Исправлен сбой в Windows, когда `WM_ERASEBKGND` закрашивает фон сплошным цветом во время перехода сворачивания/восстановления, при котором `GetClientRect` возвращает nil (о необходимости проверки сообщил @sinspired в #5636)
- Исправлена ошибка, из-за которой ошибка привязок фронтенда всегда разбиралась как текст, от @mbaklor в #5690
- Исправлена ошибка сборки для Windows при использовании тега сборки `server`, вызванная отсутствием в файлах GUI для Windows ограничения сборки `!server`, которое уже есть в эквивалентных файлах для macOS и Linux (#5680)

## v3.0.0-alpha2.114 — 2026-07-05

## Добавлено

- Реализованы протокол Update Manifest и поставщик конечной точки в [PR](https://github.com/wailsapp/wails/pull/5720) от @taliesin-ai

## Изменено

- Привязка `webview2` включена в модуль v3 под именем `v3/internal/webview2`; удалены отдельный модуль, его процессы ночного выпуска и синхронизации, а также манипуляции с версией go.mod (его единственный потребитель — v3) в [PR](https://github.com/wailsapp/wails/pull/5711) от @taliesin-ai

## Исправлено

- Исправление обнаружения масштаба монитора в WebView2 и повторной синхронизации хоста при изменении DPI перенесено в раздел Unreleased в [PR](https://github.com/wailsapp/wails/pull/5750) от @taliesin-ai
- Обновлён маршалинг COM в WebView2 для параметров float64 и BOOL в [PR](https://github.com/wailsapp/wails/pull/5741) от @wayneforrest
- Предотвращены паника и разыменование nil при обновлении и уничтожении значка в системном трее Windows в [PR](https://github.com/wailsapp/wails/pull/5703) от @wayneforrest
- Исправлена ошибка, из-за которой скрытые окна в Windows не скрывались повторно должным образом, в [PR](https://github.com/wailsapp/wails/pull/5743) от @wayneforrest
- Видимость контроллера WebView2 синхронизирована со сворачиванием, разворачиванием и восстановлением окна в [PR](https://github.com/wailsapp/wails/pull/5742) от @wayneforrest

### Исправлено

- Повторно включено обнаружение масштаба монитора в WebView2, а повторная синхронизация хоста теперь выполняется только при изменении DPI, в [PR](https://github.com/wailsapp/wails/pull/5734) от @taliesin-ai на основе исправления, проверенного @randalmurphal, с подтверждением первопричины от @eleclin и тестированием на оборудовании от @qq540491950

## v3.0.0-alpha2.113 — 2026-07-04

## Добавлено

- Добавлено предупреждение при сборке релизного AAB без заданного `ANDROID_KEYSTORE_FILE` (Google Play отклоняет пакеты с отладочной подписью), а также документация по упаковке и подписыванию App Bundle в [PR](https://github.com/wailsapp/wails/pull/5730) от @taliesin-ai
- Документация «Почему Wails» добавлена на нескольких языках в [PR](https://github.com/wailsapp/wails/pull/5739) от @taliesin-ai
- Добавлена поддержка сопоставления Go time.Time с JS Date или string в привязках в [PR](https://github.com/wailsapp/wails/pull/5398) от @fbbdev
- Добавлены задачи упаковки Android App Bundle (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) для публикации в Play Store; задачи APK сохранены для локального тестирования и тестирования в эмуляторе — [PR](https://github.com/wailsapp/wails/pull/5728) от @mortenolsrud (исправляет [#5726](https://github.com/wailsapp/wails/issues/5726))
- Добавлены цели задач для физических устройств Android и возобновление разрешений на использование камеры и местоположения — [PR](https://github.com/wailsapp/wails/pull/5735) от @taliesin-ai

## Изменено

- `webview2` обновлён до v1.0.28 ([примечания к выпуску](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- В шаблоне Android `compileSdk`/`targetSdk` обновлён с 34 до 35, как того требует Google Play для публикации новых приложений, — [PR](https://github.com/wailsapp/wails/pull/5730) от @taliesin-ai

## Исправлено

- Исправлены встроенные маски аватаров в sponsorkit — [PR](https://github.com/wailsapp/wails/pull/5745) от @leaanthony
- Исправлен выбор неверного образа системы или версии cmdline-tools при автоматическом создании Android AVD из-за лексикографической сортировки версий — [PR](https://github.com/wailsapp/wails/pull/5730) от @taliesin-ai
- Исправлена устаревшая версия Android NDK, которую предлагал мастер настройки (теперь 26.3.11579264 в соответствии с документированным требованием), — [PR](https://github.com/wailsapp/wails/pull/5730) от @taliesin-ai
- Обновлена французская документация по SvelteKit и параметрам — [PR](https://github.com/wailsapp/wails/pull/5744) от @leaanthony
- Исправлена ошибка SIGSEGV при перечислении экранов в macOS во время изменения конфигурации дисплеев — [PR](https://github.com/wailsapp/wails/pull/5516) от @flofreud

## v3.0.0-alpha2.112 - 2026-07-03

## Добавлено

- Добавлен написанный на Go генератор SVG со списком участников проекта и обновлены страницы благодарностей в документации и на сайте — [PR](https://github.com/wailsapp/wails/pull/5724) от @taliesin-ai
- Добавлена поддержка преобразования Go time.Time в JS Date или строку в привязках — [PR](https://github.com/wailsapp/wails/pull/5398) от @fbbdev

## Изменено

- Конвейер изображений спонсоров на основе Node заменён генератором на Go — [PR](https://github.com/wailsapp/wails/pull/5719) от @taliesin-ai

## Исправлено

- Исправлен скрипт установки зависимостей ресурсов сборки Android — [PR](https://github.com/wailsapp/wails/pull/5729) от @taliesin-ai
- В `ValidateAndSanitizeURL` запрещён управляющий символ U+0085 (NEXT LINE), что обеспечивает полный охват пробельных символов валидатором URL
- Добавлен повторный расчёт рамки DWM при изменении DPI для окон без рамки — [PR](https://github.com/wailsapp/wails/pull/4785) от @leaanthony
- Исправлено обнаружение зоны перетаскивания DnD, которое не работало в Windows при масштабе, отличном от 100 %, — [PR](https://github.com/wailsapp/wails/pull/4632) от @yulesxoxo
- Добавлено явное управление памятью Objective-C для объектов Cocoa в диалогах, меню, системном трее и уведомлениях Darwin — [PR](https://github.com/wailsapp/wails/pull/5714) от @taliesin-ai
- Исправлены ошибки бэкенда CGO для Linux и проблемы с системным треем — [PR](https://github.com/wailsapp/wails/pull/5718) от @taliesin-ai

## v3.0.0-alpha2.111 - 2026-07-01

## Добавлено

- HappyTools добавлен в витрину проектов сообщества — [PR](https://github.com/wailsapp/wails/pull/5061) от @Aliuyanfeng
- Добавлены поддержка индонезийской локали и подробная документация — [PR](https://github.com/wailsapp/wails/pull/5643) от @triadmoko
- В WindowsWindow добавлен параметр DisableMenu — [PR](https://github.com/wailsapp/wails/pull/4813) от @leaanthony

## Изменено

- Шаблон Taskfile и CLI обновлены для запуска задач сборки и упаковки с GOOS и ARCH — [PR](https://github.com/wailsapp/wails/pull/5617) от @leaanthony

## Исправлено

- Исправлена проблема с объединением окон во вкладки в macOS — [PR](https://github.com/wailsapp/wails/pull/5708) от @taliesin-ai

## Удалено

- Из разделов об участии в разработке, возможностях и руководствах удалены переведённые на немецкий файлы MDX — [PR](https://github.com/wailsapp/wails/pull/5702) от @taliesin-ai

## v3.0.0-alpha2.110 - 2026-06-30

## Добавлено

- Реализованы обычная и принудительная перезагрузка WebView в macOS, а также восстановление после завершения процесса WebContent — [PR](https://github.com/wailsapp/wails/pull/5129) от @wayneforrest
- Добавлена подробная документация на немецком языке по участию в разработке, возможностям и руководствам — [PR](https://github.com/wailsapp/wails/pull/5396) от @leaanthony
- Уведомления дополнены звуком, вложениями, планированием и API обновления — [PR](https://github.com/wailsapp/wails/pull/5333) от @popaprozac

## Исправлено

- Добавлен повторный расчёт рамки DWM при изменении DPI для окон без рамки — [PR](https://github.com/wailsapp/wails/pull/4785) от @leaanthony
- Исправлено обнаружение зоны перетаскивания DnD, которое не работало в Windows при масштабе, отличном от 100 %, — [PR](https://github.com/wailsapp/wails/pull/4632) от @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Добавлено

- В документацию EventsEmit добавлены примеры кода — [PR](https://github.com/wailsapp/wails/pull/5026) от @iamhabbeboy
- Добавлен параметр визуального хостинга WebView2 для Windows — [PR](https://github.com/wailsapp/wails/pull/5380) от @MerIijn
- Klustr добавлен в документацию витрины проектов сообщества — [PR](https://github.com/wailsapp/wails/pull/5536) от @SametKUM
- Kira добавлен в витрину проектов сообщества вместе с новыми страницами и записью в журнале изменений — [PR](https://github.com/wailsapp/wails/pull/5685) от @thiennguyen93
- В руководство по службе MCP добавлен раздел обратной связи — [PR](https://github.com/wailsapp/wails/pull/5694) от @taliesin-ai

## Изменено

- Для серверного режима теперь предусмотрена полноценная производственная сборка, согласованная с задачами сборки настольного приложения (#5693). По умолчанию `task build:server` собирает производственный исполняемый файл (`-tags server,production`, `-trimpath`, с удалёнными символами) и принимает `DEV=true` (сервер разработки), `OBFUSCATED=true` (garble) и `EXTRA_TAGS`. `task run:server` запускает сервер разработки. `Dockerfile.server` / `task build:docker` сначала собирают производственный сервер (`-tags server,production`) и производственный фронтенд; по умолчанию образ содержит статическую сборку на чистом Go поверх distroless/static, а `CGO_ENABLED`, `GO_IMAGE` и `RUNTIME_IMAGE` доступны как переопределяемые аргументы сборки для приложений CGO.

## Исправлено

- Предотвращён сбой при закрытии окна с незавершёнными асинхронными вызовами — [PR](https://github.com/wailsapp/wails/pull/4435) от @leaanthony
- Предотвращена активация окна при открытии скрытых приложений в Windows — [PR](https://github.com/wailsapp/wails/pull/5249) от @leaanthony
- Обеспечено выполнение обработки метаданных запросов WebKit, завершения ответов и потоков тела в основном потоке GTK — [PR](https://github.com/wailsapp/wails/pull/5668) от @taliesin-ai
- Исправлена ошибка, из-за которой `Menu.Update()` не перестраивал нативное меню в GTK4 Linux (#5659; независимо диагностировано и исправлено @puneetdixit200 в #5539)
- Исправлен сбой при перечислении экранов macOS во время изменения конфигурации дисплеев: строки идентификаторов и имён экранов теперь копируются, а количество фиксируется в снимке (#5565; независимо диагностировано и исправлено @x-haose в #5584)
- Исправлено сжатие с последующим исчезновением содержимого WebView2 после перетаскивания окна между мониторами с разным DPI в Windows: границы контроллера повторно задаются в обработчике `WM_DPICHANGED` аналогично повторной синхронизации DPI при восстановлении свёрнутого окна (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Добавлено

- Добавлена поддержка глобальных (общесистемных) сочетаний клавиш через `app.GlobalShortcut` (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Сочетания срабатывают, даже когда приложение не находится в фокусе. Поддержка реализована нативно для каждой платформы без сторонних зависимостей: горячие клавиши Carbon в macOS, `RegisterHotKey` в Windows, `XGrabKey` в X11 и интерфейс глобальных сочетаний клавиш XDG Desktop Portal в Wayland.
- Добавлен встроенный сервер MCP — сервер Model Context Protocol, который автоматически запускается, если приложение собрано с тегом `mcp`. Он позволяет LLM-агентам тестировать работающие приложения Wails и управлять ими: управлять окнами, исследовать DOM, выполнять JavaScript, вызывать привязанные методы, работать с событиями и имитировать ввод с помощью мыши и клавиатуры с отображением анимированного экранного курсора. Пользовательский код не требуется: если задано `WAILS_MCP=1`, тег `mcp` автоматически добавляется командами `wails3 build`/`wails3 dev`. Настройка полностью выполняется через переменные окружения (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Исправлено

- Исправлена ошибка, из-за которой `Menu.Update()` не пересобирал нативное меню в Linux с GTK4 (#5659; независимо диагностировано и исправлено @puneetdixit200 в #5539)
- Устранён сбой при перечислении экранов macOS после изменения конфигурации дисплеев: строки с идентификаторами и именами экранов теперь копируются, а их количество сохраняется в виде снимка (#5565; независимо диагностировано и исправлено @x-haose в #5584)

## v3.0.0-alpha2.107 - 2026-06-27

## Добавлено

- Добавлена экспериментальная документация по Wake с навигацией на боковой панели в [PR](https://github.com/wailsapp/wails/pull/5613) от @leaanthony

## v3.0.0-alpha2.106 - 2026-06-24

## Изменено

- Версия `webview2` повышена до v1.0.27.
  - ci(webview2): исправлена релизная сборка (кросс-компиляция для Windows и полный go.sum) (#5671)\

  **Полный список изменений:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- Из процесса выпуска webview2 при кросс-компиляции удалён go vet в [PR](https://github.com/wailsapp/wails/pull/5672) от @taliesin-ai
- Модель OpenRouter для автоматического журнала изменений обновлена до google/gemini-2.5-flash-lite в [PR](https://github.com/wailsapp/wails/pull/5670) от @taliesin-ai
- Версия `webview2` повышена до v1.0.26.

### Исправления

- **Приложение теперь восстанавливается после временных ошибок COM среды выполнения, а не завершает работу** (#5658, #5580). Ранее `Chromium.errorCallback` вызывал `os.Exit(1)` при *любой* ошибке COM, поэтому устранимый кратковременный сбой после запуска завершал работу всего приложения. Теперь пути среды выполнения (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) регистрируют ошибку в журнале и восстанавливают работу. В частности, некорректное или недоверенное веб-сообщение в `MessageReceived` теперь отбрасывается, а не завершает процесс. Это устраняет класс сбоев при перемещении между мониторами с разным DPI (#5544, #5650). Ошибки при создании среды или контроллера по-прежнему являются фатальными.\

**Полный список изменений:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Исправлено

- Исправлен процесс release-webview2, чтобы он корректно обрабатывал файлы go.sum, в [PR](https://github.com/wailsapp/wails/pull/5671) от @taliesin-ai
- Исправлено обновление меню GTK4 в Linux: нативное меню теперь очищается и пересобирается, в [PR](https://github.com/wailsapp/wails/pull/5659) от @taliesin-ai

## v3.0.0-alpha2.105 - 2026-06-21

## Добавлено

- Добавлен `application.System` для определения платформы во время выполнения из общего кода: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (тег сборки `server`) и `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` для непосредственного тестирования одной целевой платформы. Он компилируется для всех целевых платформ, поэтому ветвление можно выполнять без тегов сборки. Соответствующие вспомогательные функции для фронтенда (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) доступны в `@wailsio/runtime`
- Добавлено руководство «Использование других фронтенд-фреймворков», показывающее, как поместить собственный проект Vite в `frontend/` (рассматриваются Solid, Preact, Lit, SvelteKit, Qwik, Angular и другие)
- Мастер `wails3 setup` теперь проверяет инструментарий для мобильных платформ (iOS/Android): Xcode и среду выполнения iOS Simulator, JDK, Android SDK/NDK и эмулятор. Где применимо, доступны установка одним щелчком и готовые для копирования исправления конфигурации оболочки
- Создаваемые проекты включают настройку `frontend/.npmrc`, которая задаёт `minimum-release-age` длительностью 7 дн., чтобы снизить риск использования недавно опубликованных (и потенциально скомпрометированных) пакетов (поддерживается pnpm и bun; npm безопасно игнорирует эту настройку)

## Изменено

- Все встроенные начальные шаблоны получили новый дизайн с главным блоком в стиле неоновых гор (веб, iOS и Android)
- **TypeScript теперь используется в начальных шаблонах по умолчанию и занимает имя шаблона без суффикса.** Команда `wails3 init` (без `-t`) создаёт структуру проекта TypeScript; `-t vanilla`, `-t react`, `-t vue` и `-t svelte` используют TypeScript, а варианты JavaScript доступны как `-t vanilla-js`, `-t react-js`, `-t vue-js` и `-t svelte-js`. Встроенные шаблоны объявляют свой язык с помощью `typescript:` в `template.yaml`; для совместимости шаблоны сообщества с суффиксом `-ts` продолжают работать как резервный вариант
- Мастер `wails3 setup` получил новый дизайн в неоновой теме «цифровой Wails» (эффект матового стекла поверх горного фона)

## Исправлено

- Устранён сбой в Windows при восстановлении приложения, которое оставалось свёрнутым достаточно долго, чтобы WebView2 приостановил работу или перезапустил процесс рендеринга/GPU. Теперь повторная синхронизация DPI при сворачивании и восстановлении (#5544) обращается к контроллеру WebView2 только при фактическом изменении DPI окна, что предотвращает фатальные вызовы COM к приостановленному контроллеру в обычном случае восстановления с тем же DPI (#5605)
- Устранены повторяющиеся нативные сбои `SIGABRT`/`SIGSEGV` (обычно внутри `g_object_unref` во время главного цикла GTK) в длительно работающих приложениях Linux при частой загрузке ресурсов и мультимедиа. Сервер ресурсов завершал `WebKitURISchemeRequest` из рабочих горутин, вызывая потоконебезопасные функции WebKit2GTK вне главного потока GTK; теперь завершение (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) выполняется в главном потоке. Это дополняет частичное исправление в #5566. Проблема затрагивала сборки 6.0 как с GTK3, так и с GTK4/WebKitGTK (#5631, #5557)
- Устранена периодическая ошибка `fatal error: invalid pointer found on stack` в `setupSignalHandlers` в Linux/GTK3. Идентификаторы окон, передаваемые как `user_data` сигнала, хранились в локальной переменной Go типа `unsafe.Pointer`, поэтому сборщик мусора аварийно завершал работу, когда при копировании стека сканировал это значение, не являющееся указателем. Теперь на стороне Go идентификатор сохраняет целочисленный тип (`uintptr_t`): то же исправление, которое #4958 применило к пути GTK4 (где функции сигналов C были переведены на `uintptr_t` для устранения ошибок `-race`/checkptr), перенесено в устаревший путь GTK3 (#5631)

## Удалено

- Удалены начальные шаблоны `react-swc`, `preact`, `lit`, `solid`, `qwik` и `sveltekit`, а также их варианты `-ts`. Теперь поддерживаемый встроенный набор состоит из `vanilla`, `react`, `vue` и `svelte` — каждый по умолчанию использует TypeScript, а варианты JavaScript имеют суффикс `-js`. Любой другой фреймворк по-прежнему можно использовать, [добавив собственный фронтенд](https://v3.wails.io/guides/dev/frontend-frameworks) или пользовательский шаблон

## v3.0.0-alpha2.104 - 2026-06-18

## Исправлено

- Устранён сбой iOS (SIGABRT), возникавший, когда привязанный метод службы Go возвращал пустую строку. Средство записи ответа ресурсов iOS проверяло указатель на тело с помощью `buf != nil` вместо его длины, поэтому тело нулевой длины вызывало панику `&buf[0]`; теперь проверяется длина, как и в средствах записи для настольных платформ

## v3.0.0-alpha2.103 - 2026-06-15

## Изменено

- Нативные возможности iOS и Android перенесены в менеджеры платформ: теперь вызывайте их через `application.IOS.*` и `application.Android.*` (например, `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`), а не через прежние свободные функции `application.IOS*`/`application.Android*` (#5602)
- Переименованы события мобильного моста: кроссплатформенные события теперь используют префикс `common:*` (например, `common:haptic`, `common:location`), а события только для конкретной платформы — `ios:*` / `android:*` (например, `ios:backgroundTask`, `android:foregroundService`); префикс `native:*` больше не используется (#5602)

## v3.0.0-alpha.102 — 2026-06-14

## Добавлено

- Добавлен экспериментальный мастер `wails3 setup` для интерактивной настройки проекта и проверки зависимостей
- В `wails3 doctor` добавлен флаг `--json` для вывода в машиночитаемом формате
- В команду `wails3 doctor` добавлен раздел со сведениями о состоянии подписи

## Исправлено

- Исправлено обнаружение npm в Linux: теперь помимо диспетчера пакетов проверяется PATH

## v3.0.0-alpha.101 — 2026-06-13

## Добавлено

- iOS: нативные диалоговые окна сообщений (UIAlertController) и диалоговые окна открытия файла, файлов или каталога (UIDocumentPickerViewController); диалоговые окна сохранения возвращают явную ошибку
- iOS: поддержка буфера обмена через UIPasteboard
- iOS: реальные параметры экрана через UIScreen (точки, пиксели, масштаб и рабочая область с учётом безопасной зоны)
- iOS: сборки для устройств (`IOS_PLATFORM=device`), поддержка удостоверения для подписи кода, профиля подготовки и прав, упаковка `.ipa` и `deploy-device` через devicectl
- iOS: настраиваемая минимальная версия iOS (`ios.minIOSVersion` в build/config.yml)
- iOS: `wails3 doctor` сообщает о доступности Xcode и iOS SDK в macOS
- iOS: системные события — состояние батареи и сети, тема, блокировка экрана и нехватка памяти — предоставляются как события приложения `events.IOS.*` и платформонезависимые события `events.Common.*`
- iOS: нативный мост мобильных функций (экспортируемый `application.IOS*`) — системное меню общего доступа, открытие URL, предотвращение перехода в спящий режим, фонарик, отступы безопасной зоны, яркость, сведения о приложении, блокировка ориентации, строка состояния, биометрия (Face ID/Touch ID), локальные уведомления и защищённое хранилище Keychain
- iOS: датчики и оборудование — тактильная отдача, однократное определение геопозиции, акселерометр, датчик приближения, преобразование текста в речь, сведения о хранилище, состояние питания и батареи, состояние сети, отступы для клавиатуры и обнаружение захвата экрана
- iOS: документация (IOS.md и руководство на сайте документации)
- Android: нативные диалоговые окна сообщений (AlertDialog) и диалоговые окна открытия файла или файлов (Storage Access Framework; файлы импортируются как копии в кэше); диалоговые окна открытия каталога и сохранения возвращают явную ошибку
- Android: поддержка буфера обмена через ClipboardManager
- Android: реальные параметры экрана через WindowMetrics/DisplayMetrics (dp, пиксели, масштаб и рабочая область с учётом системных панелей)
- Android: методы среды выполнения для тактильной отдачи (`Android.Haptics.Vibrate`), сведений об устройстве (`Android.Device.Info`) и всплывающих уведомлений (`Android.Toast.Show`)
- Android: типизированные события жизненного цикла (`events.Android.*`, генерируемые из events.txt), где `ActivityCreated` сопоставляется с `Common.ApplicationStarted`
- Android: конвейер сборки создаёт устанавливаемые отладочные и релизные APK (`android:run`, `android:package`, `android:package:fat`); по умолчанию релиз подписывается отладочным хранилищем ключей, а настоящее хранилище можно указать с помощью переменных среды `ANDROID_KEYSTORE_*`
- Android: `wails3 doctor` сообщает сведения об Android SDK, NDK и JDK
- Android: системные события — состояние батареи и сети, тема, блокировка экрана и нехватка памяти — предоставляются как события приложения `events.Android.*` и платформонезависимые события `events.Common.*`
- Android: нативный мост мобильных функций (экспортируемый `application.Android*`) — общий доступ, открытие URL, предотвращение перехода в спящий режим, фонарик, отступы безопасной зоны, яркость, сведения о приложении, блокировка ориентации, строка состояния, биометрия (BiometricPrompt), локальные уведомления и защищённое хранилище EncryptedSharedPreferences
- Android: датчики и оборудование — тактильная отдача, однократное определение геопозиции, акселерометр, датчик приближения, преобразование текста в речь, сведения о хранилище, состояние питания и батареи, состояние сети, отступы для клавиатуры и блокировка захвата экрана с помощью FLAG_SECURE
- Android: документация (ANDROID.md и руководство на сайте документации)
- Пример: в демонстрационное приложение `mobile` добавлены вкладки Mobile и Hardware, демонстрирующие нативный мост функций в iOS и Android (вкладки-капсулы переносятся на несколько строк)
- Мобильные платформы: энергопотребление — акселерометр, датчик приближения, фонарик и периодические часы из примера приостанавливаются при переходе приложения в фоновый режим и возобновляются после возвращения (Android продолжает выполнять процесс в фоне, а на iOS состояние фонарика сохраняется аппаратно); приёмники системных событий Android регистрируются только тогда, когда приложение находится на переднем плане
- iOS: съёмка камерой — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → событие `native:capture` с миниатюрой в формате base64)
- iOS: выполнение в фоновом режиме — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (окно фоновой задачи UIApplication) и настраиваемый параметр `ios.backgroundModes` (build/config.yml), по шаблону добавляющий `UIBackgroundModes` в создаваемый Info.plist
- Android: съёмка камерой — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (системная камера через FileProvider → событие `native:capture`)
- Android: служба переднего плана — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (`WailsForegroundService` с постоянным уведомлением поддерживает процесс активным для длительной фоновой работы)
- Пример: вкладка Camera, демонстрирующая фото- и видеосъёмку, а также выполнение в фоновом режиме (служба переднего плана в Android и окно фоновой задачи в iOS)

## Исправлено

- Исправлена неизменная ошибка `NotAllowedError` при вызове `getUserMedia` в Linux: WebKitGTK отклоняет запросы разрешений, которые никто не обрабатывает, а сигнал `permission-request` не был подключён. Разрешения для камеры и микрофона теперь обрабатываются с помощью новой кроссплатформенной карты `WebviewWindowOptions.Permissions` (`map[PermissionType]Permission`), которая учитывается как в Linux (WebKitGTK), так и в Windows (WebView2). В Linux нет нативного запроса разрешения, поэтому доступ к камере и микрофону разрешён по умолчанию (что восстанавливает работу `getUserMedia`), но его можно отключить с помощью `PermissionDeny` (#5552)
- iOS: `GOOS=ios` снова компилируется (экспортирован `events.IOS`, добавлены заглушки имён мобильных методов); также компилируются сборки с тегом production (исправлены теги сборки в pkg/application и нескольких службах)
- iOS: события Go→JS и ExecJS теперь работают — страница больше не загружается дважды при запуске, а рукопожатие `wails:runtime:ready` больше не может быть потеряно
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` больше не вступают в состояние гонки с запуском приложения; удалена фиксированная 2-секундная задержка запуска
- iOS: устранена утечка C-строки при каждом выполнении JavaScript из Go в JS
- iOS: `hasListeners` теперь отражает фактическую регистрацию обработчиков событий
- iOS: отладочное журналирование фреймворка исключается при компиляции производственных сборок
- Android: `GOOS=android` снова компилируется — определён `events.Android`, удалён выходящий за границы массив обработчиков `events_android.go`, добавлена заглушка имени мобильного метода, а файлы для настольной версии Linux (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) больше не попадают в сборки Android
- Android: привязки JS→Go теперь работают — WebView не может передавать тела POST-запросов `fetch()` в `shouldInterceptRequest`, поэтому вызовы среды выполнения направляются через транспорт JavascriptInterface (`nativeHandleRuntimeCall`), а не завершаются аварийно из-за пустого тела запроса
- Android: вызовы среды выполнения `Screens.*` возвращают реальные данные — ScreenManager теперь заполняется при запуске (ранее он вообще не был подключён, поэтому `GetAll` возвращал nil)
- Android: отладочное журналирование фреймворка исключается при компиляции производственных сборок, а в отладочных сборках направляется в logcat с тегом `Wails`
- Android: настоящий реестр `hasListeners`, обработка ссылок и исключений JNI, а также жизненный цикл с однократной загрузкой страницы (без двойной навигации)
- Исправлена ошибка «Access is denied» при выполнении `wails3 generate bindings` в Windows во время работы сервера разработки Vite: создаваемые файлы теперь синхронизируются с выходным каталогом вместо переименования поверх него (#5515)
- Исправлен периодический фатальный сбой в macOS при чтении сведений об экране после смены дисплея: идентификатор и имя экрана хранили указатели на автоматически освобождаемые буферы `UTF8String`, которые могли быть освобождены до того, как Go их скопирует (use-after-free). Теперь для строк вызывается `strdup`, после преобразования они освобождаются, а перечисление экранов выполняется в явном пуле автоматического освобождения, поэтому при вызове из горутин Go утечки больше не возникает (#5556)
- Исправлен периодический SIGSEGV в Linux при закрытии assetserver объекта `WebKitURISchemeRequest`: финальный `g_object_unref` выполнялся в горутине assetserver, из-за чего финализация GObject WebKit происходила вне главного потока GTK. Теперь unref передаётся в главный контекст GTK через `g_main_context_invoke` (#5557)

## v3.0.0-alpha.100 — 2026-06-13

## Добавлено

- В `MacWebviewPreferences` добавлены дополнительные параметры конфигурации WKWebView: `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` и `ApplicationNameForUserAgent` (#5549)

## Исправлено

- Исправлена ошибка «Access is denied» в `wails3 generate bindings` на Windows при работающем сервере разработки Vite: теперь созданные файлы синхронизируются с выходным каталогом, а не переименовываются поверх него (#5561)
- Исправлены отсутствие событий изменения размера JS для окон без рамки в Linux и определение границ полосы прокрутки в таких окнах (#5368)
- Исправлена ошибка «invalid cross-device link» в средстве обновления на Windows, возникавшая, когда временный каталог и каталог установки находились на разных томах (#5560)

## v3.0.0-alpha.99 — 2026-06-10

## Исправлено

- Исправлена ошибка «Access is denied» в `wails3 generate bindings` на Windows при работающем сервере разработки Vite: теперь созданные файлы синхронизируются с выходным каталогом, а не переименовываются поверх него (#5515)

## v3.0.0-alpha.98 — 2026-06-03

## Исправлено

- Устранено зависание интерфейса WebKit в Linux во время простоя (например, при открытом инспекторе): `SA_ONSTACK` больше не применяется принудительно к `SIGUSR1`, поскольку это нарушало синхронизацию потока сборщика мусора JavaScriptCore (#5527)

## v3.0.0-alpha.97 — 2026-05-31

## Добавлено

- Добавлена страница об отладке и работе с `runtime/trace`

## Изменено

- Удалены некоторые ненужные импорты `_ "embed"`, что немного упростило код

## Исправлено

- Исправлено несоблюдение ограничений минимальной ширины и высоты после выхода окна из развёрнутого состояния в Windows (#4593)
- Исправлено прохождение щелчков мыши сквозь полноэкранное окно с параметрами Frameless + Transparent (#4408)

## v3.0.0-alpha.96 — 2026-05-25

## Добавлено

- Добавлена поддержка обфускации Garble ([#4563](https://github.com/wailsapp/wails/issues/4563)): стабильные идентификаторы методов привязки, интеграция со сборкой и Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`), а также теги структур JSON во всех передаваемых среде выполнения данных (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`), благодаря чему формат обмена сохраняется при переименовании экспортируемых полей инструментом Garble.

## v3.0.0-alpha.95 — 2026-05-20

## Добавлено

- Добавлена отсутствовавшая страница о структуре проекта

## Изменено

- Документация: несколько диаграмм на странице архитектуры заменены диаграммами последовательности для более наглядного отображения
- Документация: добавлено примечание об установке D2 как необходимом условии для запуска

## Исправлено

- Исправлен `wails3 generate appimage` при использовании GTK4 по умолчанию: теперь упаковщик определяет стек GTK по бинарному файлу до поиска файлов среды выполнения, поэтому для сборок GTK4 он выбирает `libwebkitgtkinjectedbundle.so` (в `webkitgtk-6.0/`), а для сборок `-tags gtk3` — `libwebkit2gtkinjectedbundle.so` (в `webkit2gtk-4.1/`). Проверка `.relr.dyn` теперь также учитывает `libgtk-4.so.1`, поэтому в современных наборах инструментов удаление символов корректно отключается независимо от стека. (#5475)
- Исправлен сбой `wails3 generate appimage` при вызове с относительным `-builddir`: теперь упаковщик заранее преобразует `-binary`, `-icon`, `-desktopfile`, `-builddir` и `-outputdir` в абсолютные пути, чтобы выполняемый в середине процесса `s.CD` не нарушал работу горутины загрузки AppRun и проверки `ldd` после копирования.
- Исправлена ошибка, из-за которой `wails3 generate appimage` не мог переместить итоговый AppImage в `-outputdir`, когда поле `Name=` desktop-файла не совпадало с базовым именем бинарного файла: теперь упаковщик через переменную окружения `OUTPUT` принудительно указывает плагину appimage инструмента linuxdeploy записывать AppImage в `<binary>-<arch>.AppImage`, а не под именем, полученным из desktop-файла.
- Исправлено отсутствие срабатывания `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` и `Common.SystemDidWake` в Linux после того, как стек GTK4 + WebKitGTK 6.0 стал используемым по умолчанию в alpha.93. Новый используемый по умолчанию `application_linux.go` `run()` не вызывал `setupCommonEvents()`, который перенаправляет события `Linux.*` соответствующим обработчикам `Common.*`, или `monitorPowerEvents()`. Теперь вспомогательный компонент монитора питания DBus используется совместно путями сборки GTK3 и GTK4 через `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.94 — 2026-05-19

## Исправлено

- Исправлено отсутствие срабатывания `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` и `Common.SystemDidWake` в Linux после того, как стек GTK4 + WebKitGTK 6.0 стал используемым по умолчанию в alpha.93. Новый используемый по умолчанию `application_linux.go` `run()` не вызывал `setupCommonEvents()`, который перенаправляет события `Linux.*` соответствующим обработчикам `Common.*`, или `monitorPowerEvents()`. Теперь вспомогательный компонент монитора питания DBus используется совместно путями сборки GTK3 и GTK4 через `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.93 — 2026-05-17

## Добавлено

- В вывод `wails3 doctor` в Linux добавлен `XDG_SESSION_TYPE`, автор: @leaanthony

## Исправлено

- Исправлен сбой меню окна в Wayland, вызванный обращением appmenu-gtk-module к ещё не реализованному окну (#4769), автор: @leaanthony
- Исправлен сбой приложения GTK, когда имя приложения содержит недопустимые символы (пробелы, круглые скобки и т. п.), автор: @leaanthony
- Исправлена ошибка «not enough memory» при инициализации перетаскивания в Windows (#4701), автор: @overlordtm
- Исправлена гонка данных в хранилище обратных вызовов главного потока, вызванная ошибочным использованием RLock при удалении из map (Linux, macOS, iOS) (#4424), автор: @leaanthony
- Исправлена обработка переменных при передаче аргументов командной строки задачам. Теперь переменные CLI, заданные парами KEY=VALUE, правильно инициализируются и передаются на всех этапах выполнения задачи.
- Устранён конфликт NSWindowZoomButton в macOS: теперь `MaximiseButtonState` и `FullscreenButtonState` применяют более строгое состояние как при запуске, так и во время выполнения; ни один из сеттеров больше не может незаметно переопределить другой (#5319)
- Исправлена группа ранее существовавших ошибок в устаревшем пути сборки GTK3 (`-tags gtk3`), обнаруженных CodeRabbit в #5463: запуск через файловые ассоциации больше не пропускает обработчики запуска; `getTheme` теперь безопасен с точки зрения границ и типов; `appName` больше не освобождает память, принадлежащую GLib; `clipboardGet` больше не приводит к утечке `gchar*`, возвращённого GTK; теперь `Calloc` использует получатели-указатели, а `NewCalloc` возвращает `*Calloc`, поэтому пул действительно отслеживает выделения памяти; `zoomOut` использует величину, обратную `zoomInFactor`, вместо отрицательного множителя, из-за которого значение ограничивалось до 1.0; `execJS` повторно использует заранее выделенное пустое имя мира вместо утечки одного `C.CString("")` при каждом вызове; из `menuItem.setAccelerator` удалён отладочный `fmt.Println`. Решает #5465.
- Исправлена такая же утечка из-за получателя-значения `Calloc` в стандартной сборке GTK4 (`linux_cgo.go`): теперь используются получатели-указатели и `NewCalloc() *Calloc`, поэтому выделения `c.String(...)` для каждого окна действительно отслеживаются и освобождаются.

## v3.0.0-alpha.92 - 2026-05-15

## Добавлено

- Taskfile изменены, чтобы используемым менеджером пакетов фронтенда можно было управлять с помощью параметра `PACKAGE_MANAGER`
- Данные шаблонов дополнены `{{.Opn}}` и `{{.Cls}}`, чтобы сделать создание шаблонов Taskfile более предсказуемым

## Изменено

- Несколько существующих Taskfile изменены для использования `{{.Opn}} and {{.Cls}}`

## Исправлено

- Исправлена фатальная ошибка среды выполнения `concurrent map read and map write` в `linuxSystemTray`, возникавшая при обновлении меню в области уведомлений во время его чтения панелью.
- Для вывода ошибок и трассировок стека WebView2 вместо `fmt` теперь используется `log`, чтобы сообщения не терялись при запуске приложения в Windows без подключённой консоли.

## v3.0.0-alpha.91 - 2026-05-12

## Изменено

- Обновлён SVG-файл спонсоров в [PR](https://github.com/wailsapp/wails/pull/5414) от `@github-actions[bot]`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ (macOS):** система координат macOS унифицирована: `GetScreens`, `Position` и `SetPosition` теперь используют единое пространство — логические точки, ось Y направлена вниз, а `(0,0)` находится в левом верхнем углу основного экрана. Это соответствует Windows, GTK и публичным API Electron и веб-платформы. Для экранов, физически расположенных выше основного, `Bounds.Y` теперь имеет отрицательное значение (ранее положительное), а значения `Position()`/`SetPosition()` теперь выражаются в логических точках вместо `points × primaryScale`. Обратимое преобразование `Position()` → `SetPosition()` сохранено; абсолютные значения из журналов прежних альфа-сборок и обходные решения с вычислениями вручную (например, умножение на `primaryScale` или инверсия Y относительно высоты экрана) потребуется обновить. Решает [#5117](https://github.com/wailsapp/wails/issues/5117).

## Исправлено

- Добавлена защитная проверка имени сигнала DBus и длины его тела для предотвращения паник в [PR](https://github.com/wailsapp/wails/pull/5416) от @leaanthony
- Исправлена проблема безопасности памяти при обработке меню GTK в Linux в [PR](https://github.com/wailsapp/wails/pull/5363) от @leaanthony
- Добавлено обнаружение графических процессоров NVIDIA и отключение средства визуализации DMA-BUF в Linux в [PR](https://github.com/wailsapp/wails/pull/5295) от @leaanthony
- Исправлено преобразование Y между экранами в `SetPosition` на macOS: высота основного экрана теперь используется как глобальная точка отсчёта, поэтому окна располагаются правильно на мониторах, смещённых по вертикали относительно основного дисплея, в [#5117](https://github.com/wailsapp/wails/issues/5117)
- В шаблоне PR для Git исправлена ссылка на правильный URL для обратной связи в [PR](https://github.com/wailsapp/wails/pull/5109) от @wayneforrest
- Исправлена группа сбоев `SetMenu` в области уведомлений Windows, вызванных некорректным системным вызовом `DestroyMenu`: он передавал четыре аргумента вместо одного, поэтому каждый вызов возвращал FALSE и ничего не освобождал. Кроме того, при перестроении меню теперь освобождаются дескрипторы HMENU и HBITMAP (включая выделенные во время выполнения через `MenuItem.SetBitmap`), в `Win32Menu.Update` сбрасываются устаревшие карты флажков и переключателей, а из `systemtray.updateMenu` удалён избыточный вызов `Update()`, удваивавший число выделений памяти. Долго работающие приложения со значком в области уведомлений больше не теряют объекты GDI/USER при каждом перестроении меню.

## v3.0.0-alpha.90 - 2026-05-11

## Добавлено

- Добавлена возможность настраивать имя приложения в User-Agent WKWebView на macOS в [PR](https://github.com/wailsapp/wails/pull/5261) от @vinhvoit225
- В пример gin-service добавлена косвенная зависимость github.com/coder/websocket в [PR](https://github.com/wailsapp/wails/pull/5400) от @taliesin-ai
- В тесты ресурсов сборки добавлена поддержка глубокого сравнения на равенство в [PR](https://github.com/wailsapp/wails/pull/5402) от @leaanthony

## Изменено

- Выходные данные сборки объединены в каталоге ресурсов в [PR](https://github.com/wailsapp/wails/pull/5401) от @taliesin-ai
- Обновлён SVG-файл спонсоров в [PR](https://github.com/wailsapp/wails/pull/5399) от `@github-actions[bot]`

## Исправлено

- Для сообщения единственного экземпляра приложения в macOS использован объект уведомления в [PR](https://github.com/wailsapp/wails/pull/5289) от @overlordtm
- Обратные вызовы Windows объединены в пакеты для предотвращения потери промисов при высокой нагрузке в [PR](https://github.com/wailsapp/wails/pull/5383) от @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## Добавлено

- Добавлено задание go<em>test</em>results для объединения результатов тестов Go в [PR](https://github.com/wailsapp/wails/pull/5316) от @leaanthony

## Изменено

- Добавлено условное разбиение крупных полезных данных RPC на фрагментированные POST-запросы в [PR](https://github.com/wailsapp/wails/pull/5369) от @leaanthony
- Vite обновлён с 5.x.x до 8.0.0 во всех шаблонах фронтенда в [PR](https://github.com/wailsapp/wails/pull/5386) от @leaanthony
- Конфигурация порта сервера разработки Vite перенесена в переменные окружения в [PR](https://github.com/wailsapp/wails/pull/5365) от @leaanthony
- Во всех шаблонах сервер разработки Vite настроен на привязку к 127.0.0.1 в [PR](https://github.com/wailsapp/wails/pull/5361) от @leaanthony
- Обновлён SVG-файл спонсоров в [PR](https://github.com/wailsapp/wails/pull/5384) от `@github-actions[bot]`

## Исправлено

- Добавлена очистка заготовок шаблона Info.plist при обновлении ресурсов сборки в [PR](https://github.com/wailsapp/wails/pull/5312) от @leaanthony
- Исправлено устаревшее состояние меню macOS: методы изменения пунктов меню (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) теперь применяются синхронно в главном потоке, что устраняет состояние гонки `dispatch_async`, из-за которого при быстром повторном открытии меню отображали предыдущее состояние (#5002)
- Файлы `*_test.go` игнорируются в режиме разработки для предотвращения ненужных повторных сборок в [PR](https://github.com/wailsapp/wails/pull/5203) от @leaanthony
- Предотвращена ошибка сегментации в Menu.Update(), когда приложение не запущено, в [PR](https://github.com/wailsapp/wails/pull/5291) от @wucm667
- Для управления перерисовкой строки меню в Windows теперь используется lastSizeWParam в [PR](https://github.com/wailsapp/wails/pull/5382) от @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## Изменено

- HiddenOnTaskbar переведён на использование WS<em>EX</em>TOOLWINDOW в [PR](https://github.com/wailsapp/wails/pull/5371) от @leaanthony
- Зависимости переупорядочены, а директива replace для webview2 удалена из go.mod в [PR](https://github.com/wailsapp/wails/pull/5370) от @atterpac
- Обновлён SVG-файл спонсоров в [PR](https://github.com/wailsapp/wails/pull/5358) от `@github-actions[bot]`

## Исправлено

- Удалены универсальные псевдонимы косвенного обращения и унифицированы типы ключей карт в [PR](https://github.com/wailsapp/wails/pull/5331) от @fbbdev

## Удалено

- Удалён рабочий процесс PR-master вместе с документацией, тестами Go и тестами пропуска в [PR](https://github.com/wailsapp/wails/pull/5377) от @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## Добавлено

- Добавлена документация Wails v3 на корейском языке в [PR](https://github.com/wailsapp/wails/pull/5352) от @leaanthony
- Добавлена документация на французском языке по установке и быстрому началу работы в [PR](https://github.com/wailsapp/wails/pull/5354) от @leaanthony
- Добавлена документация на португальском языке по быстрому началу работы, основным концепциям и сообществу в [PR](https://github.com/wailsapp/wails/pull/5355) от @leaanthony

## v3.0.0-alpha.86 — 2026-05-06

## Добавлено

- Добавлена французская локализация документации в [PR](https://github.com/wailsapp/wails/pull/5328) от @leaanthony
- На сайт документации добавлена немецкая локаль в [PR](https://github.com/wailsapp/wails/pull/5343) от @leaanthony

## Изменено

- Все 8 переведённые локали зарегистрированы в конфигурации документации в [PR](https://github.com/wailsapp/wails/pull/5347) от @leaanthony
- Обновлены различные связанные с Windows файлы для WebView2 в [PR](https://github.com/wailsapp/wails/pull/5317) от @leaanthony

## Исправлено

- Обработка диалогов в Linux разделена между GTK3 и GTK4 в [PR](https://github.com/wailsapp/wails/pull/5340) от @leaanthony
- Выполнение обратных вызовов диалогов перенесено в поток GTK, что устраняет ошибки сегментации, в [PR](https://github.com/wailsapp/wails/pull/5339) от @leaanthony

## v3.0.0-alpha.85 — 2026-05-05

## Добавлено

- В репозиторий добавлен URL шаблона PR в [PR](https://github.com/wailsapp/wails/pull/5179) от @leaanthony
- Добавлена немецкая документация для Wails v3 в [PR](https://github.com/wailsapp/wails/pull/5330) от @leaanthony

## v3.0.0-alpha.84 — 2026-05-03

## Добавлено

- Добавлена возможность запретить выход из полноэкранного режима по клавише Escape в macOS в [PR](https://github.com/wailsapp/wails/pull/5307) от @leaanthony
- Добавлена возможность запретить выход из полноэкранного режима по клавише Escape в macOS в [PR](https://github.com/wailsapp/wails/pull/5310) от @leaanthony
- Добавлена документация о проекте Pausa в галерее сообщества в [PR](https://github.com/wailsapp/wails/pull/5288) от @yuseferi

## Изменено

- Обновлён SVG-файл спонсоров в [PR](https://github.com/wailsapp/wails/pull/5308) от `@github-actions[bot]`
- Команда генерации значков обновлена для обработки неподдерживаемых платформ в [PR](https://github.com/wailsapp/wails/pull/5309) от @leaanthony
- Булев API полноэкранного режима заменён на трёхсостоянийный ButtonState, реализованы привязки для платформ в [PR](https://github.com/wailsapp/wails/pull/5224) от @leaanthony

## Исправлено

- Операции фокусировки WebView2 защищены от состояния nil контроллера в [PR](https://github.com/wailsapp/wails/pull/5315) от @leaanthony
- Рабочий процесс GitHub Actions обновлён для правильного обращения к базовой ветке PR в [PR](https://github.com/wailsapp/wails/pull/5313) от @leaanthony
- Файлы `*_test.go` игнорируются в режиме разработки во избежание ненужных повторных сборок в [PR](https://github.com/wailsapp/wails/pull/5203) от @leaanthony
- Предотвращена ошибка сегментации в Menu.Update(), когда приложение не запущено, в [PR](https://github.com/wailsapp/wails/pull/5291) от @wucm667

## v3.0.0-alpha.83 — 2026-05-02

## Добавлено

- Добавлены флаг InstallScope и параметр сборки для установки на уровне компьютера или пользователя в [PR](https://github.com/wailsapp/wails/pull/5094) от @symball
- В BrowserWindow добавлен ничего не выполняющий метод SetScreen для соответствия интерфейсу Window в [PR](https://github.com/wailsapp/wails/pull/5294) от @leaanthony

## Исправлено

- Добавлено обнаружение GPU NVIDIA и отключение средства визуализации DMA-BUF в Linux в [PR](https://github.com/wailsapp/wails/pull/5295) от @leaanthony
- В шаблоне git PR исправлен URL для обратной связи в [PR](https://github.com/wailsapp/wails/pull/5109) от @wayneforrest
- Исправлена группа сбоев `SetMenu` в области уведомлений Windows, вызванных неисправным системным вызовом `DestroyMenu`, который передавал четыре аргумента вместо одного, из-за чего каждый вызов возвращал FALSE и ничего не освобождал. Кроме того, при перестроении меню теперь освобождаются дескрипторы HMENU и HBITMAP (включая выделенные во время выполнения через `MenuItem.SetBitmap`), в `Win32Menu.Update` сбрасываются устаревшие карты флажков и переключателей, а из `systemtray.updateMenu` удалён лишний вызов `Update()`, удваивавший количество выделений памяти. Длительно работающие приложения с областью уведомлений больше не допускают утечки объектов GDI/USER при каждом перестроении меню.

## v3.0.0-alpha.82 — 2026-05-01

## Исправлено

- Исправлена генерация desktop-файла для правильной обработки имени рабочего стола в [PR](https://github.com/wailsapp/wails/pull/5232) от @leaanthony

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

- В качестве ветки по умолчанию для задачи выпуска выбрана master и обновлена формулировка журнала изменений в [PR](https://github.com/wailsapp/wails/pull/5283) от @leaanthony
- Рабочий процесс автоматического создания журнала изменений обновлён до последней версии в [PR](https://github.com/wailsapp/wails/pull/5282) от @leaanthony
- Повышена эффективность рабочих процессов благодаря добавлению фильтров путей и удалению неиспользуемых рабочих процессов в [PR](https://github.com/wailsapp/wails/pull/5280) от @leaanthony
- В документации ссылки на примеры обновлены для использования ветки master в [PR](https://github.com/wailsapp/wails/pull/5274) от @leaanthony
- Обновлены документация и примеры для v3 в [PR](https://github.com/wailsapp/wails/pull/5272) от @leaanthony

## Исправлено

- В обратный прокси-сервер добавлены повторные попытки и принудительное использование IPv4 при разработке в [PR](https://github.com/wailsapp/wails/pull/5265) от @AkagiYui
- Переработан рабочий процесс запуска создания журнала невыпущенных изменений в [PR](https://github.com/wailsapp/wails/pull/5281) от @leaanthony

## Удалено

- Удалены сценарии тестирования оболочки для различных целей тестирования в [PR](https://github.com/wailsapp/wails/pull/5267) от @leaanthony
- Удалены рабочий процесс развёртывания документации v3-alpha и запись CNAME в [PR](https://github.com/wailsapp/wails/pull/5266) от @leaanthony

### Добавлено

- В боковую панель навигации добавлен пункт «Маршрутизация фронтенда» в [PR](https://github.com/wailsapp/wails/pull/5196) от @leaanthony
- Добавлено руководство по маршрутизации фронтенда с рекомендациями для конкретных фреймворков в [PR](https://github.com/wailsapp/wails/pull/5185) от @leaanthony
- Добавлена поддержка модальных листов (macOS)
- Версия ghw повышена для улучшения поддержки устройств Apple, автор — @leaanthony (#4977)
- В службу Dock добавлен метод `GetBadge`
- В команду `wails3 build` добавлен флаг `-tags` для передачи пользовательских тегов сборки Go (например, `wails3 build -tags gtk4`) (#4957)
- Добавлена документация по автоматической генерации перечислений в генераторе привязок, включая отдельную страницу «Перечисления» и пункт в боковой панели навигации (#4972)
- В команду `wails3 build` добавлен флаг `-tags` для передачи пользовательских тегов сборки Go (например, `wails3 build -tags gtk4`) (#4957)
- Добавить в `v3/examples/web-apis/` примеры Web API, демонстрирующие браузерные API 41, включая хранилища (localStorage, sessionStorage, IndexedDB, Cache API), сетевые API (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), мультимедийные API (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), API устройств (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), API производительности (Performance API, Mutation Observer, Intersection/Resize Observer), API пользовательского интерфейса (Web Components, Pointer Events, Selection, Dialog, Drag and Drop) и многое другое
- Добавить пример средства проверки совместимости API WebView (`v3/examples/webview-api-check/`), тестирующий более 200 браузерных API на разных платформах
- Добавить пакет `internal/libpath` для поиска путей к нативным библиотекам в Linux с параллельным поиском, кэшированием и поддержкой Flatpak/Snap/Nix
- **В разработке:** добавить экспериментальную поддержку WebKitGTK 6.0 / GTK4 для Linux, доступную через `-tags gtk4` (по умолчанию по-прежнему используется GTK3/WebKit2GTK 4.1)
- Примечание: в тайловых оконных менеджерах (например, Hyprland и Sway) операции сворачивания и разворачивания могут работать не так, как ожидается, поскольку геометрией окна управляет оконный менеджер
- В документацию по теме **«Прослушивание событий в JavaScript»** добавлено описание создания **одноразовых обработчиков** от @AbdelhadiSeddar
- Добавить в `WebviewWindowOptions` параметр `UseApplicationMenu`, позволяющий окнам в Windows/Linux наследовать меню приложения, заданное через `app.Menu.Set()`, от @leaanthony
- Добавить поддержку файлов `.icon` (формат Apple Icon Composer) для создания значков Liquid Glass и каталогов ресурсов (macOS) (#4934) от @wimaha
- Добавить экспериментальный серверный режим для безголового и веб-развёртывания (`-tags server`). Он позволяет запускать приложения Wails как HTTP-серверы без зависимостей от нативного графического интерфейса. Выполняйте сборку с `wails3 task build:server`. Подробности см. в `examples/server`.
- Добавить пакет `internal/libpath` для поиска путей к нативным библиотекам в Linux с параллельным поиском, кэшированием и поддержкой Flatpak/Snap/Nix
- Добавить в `MacWindow` параметр `CollectionBehavior` для управления поведением окна в пространствах macOS Spaces и полноэкранном режиме (#4756) от @leaanthony
- Добавить модульные тесты для pkg/application от @leaanthony
- Добавить поддержку пользовательских протоколов в упаковку MSIX от @leaanthony
- Добавить обнаружение среды рабочего стола в Linux, [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- Добавить в среду выполнения JavaScript метод `Window.Print()` для вызова диалогового окна печати из фронтенда (#4290) от @leaanthony
- Добавить `XDG_SESSION_TYPE` в вывод `wails3 doctor` в Linux от @leaanthony
- Добавить дополнительные события изменения состояния загрузки WebKit2 для Linux: `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` (#3896) от @leaanthony
- Добавить `XDG_SESSION_TYPE` в вывод `wails3 doctor` в Linux от @leaanthony
- Создавать файл `.desktop` во время сборки для Linux, а не только при упаковке (#4575)
- Добавить документацию по зависимостям среды выполнения Linux с названиями пакетов для отдельных дистрибутивов и примерами упаковки с помощью nfpm (#4339) от @leaanthony
- Добавить сведения о версии драйвера NVIDIA в вывод `wails3 doctor` в Linux от @leaanthony
- Добавить origin в обработчик необработанных сообщений от @APshenkin в [PR](https://github.com/wailsapp/wails/pull/4710)
- Добавить поддержку универсальных ссылок для macOS от @APshenkin в [PR](https://github.com/wailsapp/wails/pull/4712)
- Выполнить рефакторинг транспортного уровня привязок от @APshenkin в [PR](https://github.com/wailsapp/wails/pull/4702)
- Добавить идентификаторы aria-label в шаблоны helloworld, чтобы пример приложения можно было легко тестировать с помощью клиентов Appium, от @chinenual в [PR](https://github.com/wailsapp/wails/pull/4760)
- Добавить origin в обработчик необработанных сообщений от @APshenkin в [PR](https://github.com/wailsapp/wails/pull/4710)
- Добавить поддержку универсальных ссылок для macOS от @APshenkin в [PR](https://github.com/wailsapp/wails/pull/4712)
- Выполнить рефакторинг транспортного уровня привязок от @APshenkin в [PR](https://github.com/wailsapp/wails/pull/4702)
- Типизированные события от @fbbdev и @ianvs в [#4633](https://github.com/wailsapp/wails/pull/4633)
- Добавить пример `systray-clock`, демонстрирующий трей без основного окна с динамически обновляемой всплывающей подсказкой (#4653).
- Добавлен шаблон протокола NSIS для Windows от @Tolfx в #4510
- Добавлены тесты для build-assets от @Tolfx в #4510
- macOS: отображение нативных элементов управления окном в строке меню в [#4588](https://github.com/wailsapp/wails/pull/4588) от @nidib
- Добавить службу Dock для macOS, позволяющую скрывать и отображать значок приложения в Dock, от @popaprozac в [PR](https://github.com/wailsapp/wails/pull/4451)
- Добавить службу Dock для macOS, позволяющую скрывать и отображать значок приложения в Dock, от @popaprozac в [PR](https://github.com/wailsapp/wails/pull/4451)
- Добавить поддержку нативного эффекта Liquid Glass для macOS с NSGlassEffectView (macOS 15.0+) и резервным использованием NSVisualEffectView, включая обширные возможности настройки материала, от @leaanthony в [#4534](https://github.com/wailsapp/wails/pull/4534)
- Санитизация URL браузера от @leaanthony в [#4500](https://github.dev/wailsapp/wails/pull/4500). Основано на [#4484](https://github.com/wailsapp/wails/pull/4484) от @APShenkin.
- Добавить защиту содержимого в Windows/Mac от [@leaanthony](https://github.com/leaanthony) на основе исходной работы [@Taiterbase](https://github.com/Taiterbase) в этом [PR](https://github.com/wailsapp/wails/pull/4241)
- Добавить поддержку передачи переменных CLI командам Task через псевдонимы `wails3 build` и `wails3 package` (#4422) от @leaanthony в [PR](https://github.com/wailsapp/wails/pull/4488)
- Поддержка зон перетаскивания, в которых событие предоставляет данные о перетащенном элементе, от [@atterpac](https://github.com/atterpac) в [#4318](https://github.com/wailsapp/wails/pull/4318)
- В параметры `WindowsWindow` добавлен `AdditionalLaunchArgs`, позволяющий передавать браузеру WebView2 дополнительные аргументы командной строки, в [PR](https://github.com/wailsapp/wails/pull/4467)
- Добавлен автоматический запуск go mod tidy после wails init от [@triadmoko](https://github.com/triadmoko) в [PR](https://github.com/wailsapp/wails/pull/4286)
- Функция Windows Snap Assist от @leaanthony в [PR](https://github.dev/wailsapp/wails/pull/4463)
- В параметры `WindowsWindow` добавлен `AdditionalLaunchArgs`, позволяющий передавать браузеру WebView2 дополнительные аргументы командной строки, в [PR](https://github.com/wailsapp/wails/pull/4467)
- Добавлен автоматический запуск go mod tidy после wails init от [@triadmoko](https://github.com/triadmoko) в [PR](https://github.com/wailsapp/wails/pull/4286)
- Функция Windows Snap Assist от @leaanthony в [PR](https://github.dev/wailsapp/wails/pull/4463)
- Добавить реализацию `getAccentColor` для Windows от [@almas-x](https://github.com/almas-x) в [PR](https://github.com/wailsapp/wails/pull/4427)
- Добавить реализацию `getAccentColor` для Windows от [@almas-x](https://github.com/almas-x) в [PR](https://github.com/wailsapp/wails/pull/4427)
- Меню и строка меню с тёмной темой в Windows. От @leaanthony в [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Переименовать встроенные службы, чтобы сделать привязки JS/TS понятнее, от @popaprozac в [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` для получения акцентного цвета системы пользователя. Работает в macOS. От [@etesam913](https://github.com/etesam913)
- Добавить API `window.ToggleFrameless()` от [@atterpac](https://github.com/atterpac) в [#4137](https://github.com/wailsapp/wails/pull/4137)
- Добавить зависимости сборки для отдельных дистрибутивов Linux от @leaanthony в [PR](https://github.com/wailsapp/wails/pull/4345)
- Руководство по привязкам добавлено @atterpac в [PR](https://github.com/wailsapp/wails/pull/4404)
- **Упорядочена инфраструктура тестирования**: файлы тестов Docker перенесены в отдельный каталог `test/docker/`, образы оптимизированы, а надёжность сборки повышена — [@leaanthony](https://github.com/leaanthony) в [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Улучшены шаблоны управления ресурсами**: в примеры добавлены корректная очистка обработчиков событий и управление горутинами с учётом контекста — [@leaanthony](https://github.com/leaanthony) в [#4359](https://github.com/wailsapp/wails/pull/4359)
- Добавлена поддержка сборок AppImage для aarch64 — [@AkshayKalose](https://github.com/AkshayKalose) в [#3981](https://github.com/wailsapp/wails/pull/3981)
- В `wails doctor` добавлен раздел диагностики — [@leaanthony](https://github.com/leaanthony)
- При вызове метода сервиса окно добавлено в контекст — [@leaanthony](https://github.com/leaanthony)
- Добавлен пример `window-call`, показывающий, как определить, какое окно вызывает сервис, — [@leaanthony](https://github.com/leaanthony)
- Новое руководство по меню — [@leaanthony](https://github.com/leaanthony)
- Улучшена обработка паник — [@leaanthony](https://github.com/leaanthony)
- Новое руководство по меню — [@leaanthony](https://github.com/leaanthony)
- Добавлены документирующие комментарии для Service API — [@fbbdev](https://github.com/fbbdev) в [#4024](https://github.com/wailsapp/wails/pull/4024)
- Добавлена функция `application.NewServiceWithOptions` для инициализации сервисов с дополнительной конфигурацией — [@leaanthony](https://github.com/leaanthony) в [#4024](https://github.com/wailsapp/wails/pull/4024)
- Улучшено управление меню — [@FalcoG](https://github.com/FalcoG) и [@leaanthony](https://github.com/leaanthony) в [#4031](https://github.com/wailsapp/wails/pull/4031)
- Расширена документация — [@leaanthony](https://github.com/leaanthony)
- Добавлена поддержка отмены событий в стандартных обработчиках событий — [@leaanthony](https://github.com/leaanthony)
- Добавлена поддержка `Hide`, `Show` и `Destroy` для системного трея — [@leaanthony](https://github.com/leaanthony)
- Добавлена поддержка `SetTooltip` для системного трея — [@leaanthony](https://github.com/leaanthony). Автор исходной идеи — [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- В предупреждениях генератора привязок о неподдерживаемых типах теперь указывается путь к пакету — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- В генератор привязок добавлена поддержка обобщённых псевдонимов — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- В генератор привязок добавлена поддержка флага JSON `omitzero` — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Добавлена директива `//wails:ignore`, предотвращающая создание привязок для выбранных методов сервисов, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Для сервисов и моделей добавлена директива `//wails:internal`, позволяющая использовать типы, экспортируемые в Go, но не в JS/TS, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- В генератор привязок добавлена поддержка констант типов-псевдонимов, позволяющая создавать перечисления со слабой типизацией, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Добавлены тесты генератора привязок для возможностей Go 1.24 — [@fbbdev](https://github.com/fbbdev) в [#4068](https://github.com/wailsapp/wails/pull/4068)
- В `OSInfo.Branding` добавлена поддержка macOS 15 «Sequoia» для улучшенного определения версии ОС — [#4065](https://github.com/wailsapp/wails/pull/4065)
- Добавлен хук `PostShutdown` для выполнения пользовательского кода после завершения процесса остановки — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавлена структура `FatalError` для обнаружения фатальных ошибок в пользовательских обработчиках ошибок — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Порядок запуска и остановки сервисов стандартизирован и задокументирован — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавлен тестовый стенд для последовательности запуска и остановки приложения, а также тесты запуска и остановки сервисов — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавлен метод `RegisterService` для регистрации сервисов после создания приложения — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- В параметры приложения и сервисов добавлено поле `MarshalError` для пользовательской обработки ошибок при вызовах привязок — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавлена отменяемая обёртка Promise, передающая запросы на отмену по цепочкам Promise, — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Добавлена возможность связать отмену вызова привязки с `AbortSignal` — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Для WML наряду с обычными атрибутами `wml-*` добавлена поддержка атрибутов `data-wml-*` — [@leaanthony](https://github.com/leaanthony)
- Во все сервисы добавлен метод `Configure` для поздней или динамической повторной настройки — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Если сервис `fileserver` не настроен, он отправляет ответ 503 Service Unavailable — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Если сервис `kvstore` не настроен, по умолчанию он предоставляет хранилище пар «ключ — значение» в памяти — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- В сервис `kvstore` добавлен метод `Load` для повторной загрузки данных из файла после изменения конфигурации — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- В сервис `kvstore` добавлен метод `Clear` для удаления всех ключей — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- В сервис `log` добавлен тип `Level`, предоставляющий на стороне JS константы уровней журналирования, — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- В сервис `log` добавлен метод `Log` для динамического задания уровня журналирования — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Если сервис `sqlite` не настроен, по умолчанию он предоставляет базу данных в памяти — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- В сервис `sqlite` добавлен метод `Close` для ручного закрытия базы данных — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Для методов запросов сервиса `sqlite` добавлена поддержка отмены — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- В сервис `sqlite` с привязками JS добавлена поддержка подготовленных выражений — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Поддержка Gin от [Lea Anthony](https://github.com/leaanthony) в [PR](https://github.com/wailsapp/wails/pull/3537) на основе исходной работы [@AnalogJ](https://github.com/AnalogJ) в этом [PR](https://github.com/wailsapp/wails/pull/3537)
- Исправлена ошибка, из-за которой автосохранение и автоматическое сохранение паролей всегда были включены, — [@oSethoum](https://github.com/osethoum) в [#4134](https://github.com/wailsapp/wails/pull/4134)
- В окно добавлен `SetMenu()`, позволяющий задать для него меню, — [@leaanthony](https://github.com/leaanthony)
- Добавлена поддержка уведомлений от [@popaprozac](https://github.com/popaprozac) в [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Добавлена поддержка ассоциаций файлов для mac от [@wimaha](https://github.com/wimaha) в [#4177](https://github.com/wailsapp/wails/pull/4177)
- Добавлен `wails3 tool version` для повышения семантической версии — [@leaanthony](https://github.com/leaanthony)
- Добавлена поддержка значков-индикаторов для macOS и Windows от [@popaprozac](https://github.com/popaprozac) в [#](https://github.com/wailsapp/wails/pull/4234)
- Добавлена поддержка зарегистрированных событий со строгой типизацией от [@fbbdev](https://github.com/fbbdev) и [@IanVS](https://github.com/IanVS) в [#4161](https://github.com/wailsapp/wails/pull/4161)
- Добавлена возможность регистрировать обработчики-перехватчики пользовательских событий от [@fbbdev](https://github.com/fbbdev) и [@IanVS](https://github.com/IanVS) в [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` для открытия системного файлового менеджера по пути `path` с необязательным выделением через `selectFile` — [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Новый флаг `-git` для команды `wails3 init` от [@leaanthony](https://github.com/leaanthony)
- Новая команда `wails3 generate webview2bootstrapper` от [@leaanthony](https://github.com/leaanthony)
- В среду выполнения добавлен метод `init()`, позволяющий инициализировать её вручную, — [@leaanthony](https://github.com/leaanthony)
- В WindowOptions окна добавлен параметр `WindowDidMoveDebounceMS` — [@leaanthony](https://github.com/leaanthony)
- Добавлена функция единственного экземпляра от [@leaanthony](https://github.com/leaanthony). Основано на [v2 PR](https://github.com/wailsapp/wails/pull/2951) от @APshenkin.
- Команда `wails3 generate template` от [@leaanthony](https://github.com/leaanthony)
- Команда `wails3 releasenotes` от [@leaanthony](https://github.com/leaanthony)
- Команда `wails3 update cli` от [@leaanthony](https://github.com/leaanthony)
- Параметр `-clean` для команды `wails3 generate bindings` от [@leaanthony](https://github.com/leaanthony)
- Добавлена возможность создавать сборки Linux AppImage для aarch64 (arm64) от [@AkshayKalose](https://github.com/AkshayKalose) в [#3981](https://github.com/wailsapp/wails/pull/3981)
- Добавлена гиперссылка для спонсора от @ansxuman в [#3958](https://github.com/wailsapp/wails/pull/3958)
- Поддержка сборок пакетов deb, rpm и Arch Linux для Linux от
- Добавлена поддержка универсальных сборок и пакетов для Darwin от
- Документация по событиям на веб-сайте от
- Шаблоны для sveltekit и sveltekit-ts, настроенные для разработки без SSR
- Обновление ресурсов сборки с помощью новой команды `wails3 update build-assets` от
- Пример для тестирования HTML Drag and Drop API от
- Поддержка ассоциаций файлов от [leaanthony](https://github.com/leaanthony) в
- Новая команда `wails3 generate runtime` от
- Новый параметр `InitialPosition`, определяющий, должно ли окно быть отцентрировано или
- В пакет `application` добавлены методы `Path` и `Paths` от
- Добавлены параметры Windows `GeneralAutofillEnabled` и `PasswordAutosaveEnabled`
- Добавлена возможность получить окно, вызывающее метод сервиса, от
- Добавлены параметры `EnabledFeatures` и `DisabledFeatures` для WebView2 от
- ⊞ Новая система DIP для улучшенной поддержки мониторов с высоким DPI от
- ⊞ Параметр имени класса окна от [windom](https://github.com/windom/) в
- Возможности сервисов расширены для поддержки плагинов. Автор:
- 🐧 События WindowDidMove / WindowDidResize в
- ⊞ Событие WindowDidResize в
-  добавлено событие ApplicationShouldHandleReopen для обработки значка в Dock
-  в реализацию добавлены getPrimaryScreen/getScreens от @tmclane в
-  добавлен параметр для отображения панели инструментов в полноэкранном режиме macOS от
- 🐧 добавлена логика onKeyPress для преобразования нажатия клавиши Linux в акселератор
- 🐧 добавлена задача `run:linux` от
- Экспортирован метод `SetIcon` от [@almas-x](https://github.com/almas-x) в
- Улучшен `OnShutdown` — [@almas-x](https://github.com/almas-x) в
- В интерфейсе `Window` восстановлен метод `ToggleMaximise` от
- В `Environment()` добавлена дополнительная информация. Автор: @leaanthony в
- Метод `WebviewWindow.IsFocused` предоставлен через интерфейс `Window` —
- Поддержка нескольких разделённых пробелами событий-триггеров в системе WML от
- Добавлен экспорт ESM из комплектного скрипта среды выполнения JS от
- Добавлен флаг генератора привязок для использования комплектного скрипта среды выполнения JS вместо
- `setIcon` реализован в Linux — [@abichinger](https://github.com/abichinger)
- В команду dev добавлен флаг `-port` и поддержка переменной окружения
- Добавлены тесты вызовов привязанных методов от
- ⊞ для уже созданного окна добавлен `SetIgnoreMouseEvents` от
-  Добавлена возможность задавать уровень расположения окна в стеке (порядок) от

### Исправлено

- Исправлены уменьшенные вдвое `Screen.Bounds`, `WorkArea` и `Size` на компьютерах Mac с Retina: значения точек NSScreen преобразуются в пиксели устройства в полях `Physical*`, а значения верхнего уровня `Screen.X`/`Y` заполняются, чтобы обнаружение соприкосновения нескольких мониторов и размещение в рабочей области выполнялись правильно; [PR](https://github.com/wailsapp/wails/pull/5168) от @wayneforrest
- Исправлена гонка данных в ScreenManager, которая приводила к взаимной блокировке WebKit DisplayLink при изменении конфигурации дисплеев (например, при горячем подключении внешнего монитора во время перехода в режим сна или выхода из него)
- Если существует Assets.car, CFBundleIconName напрямую задаётся как appicon; [PR](https://github.com/wailsapp/wails/pull/5154) от @symball
- Исправлено отображение неверных пакетов WebKitGTK командой `wails3 doctor` в Fedora, openSUSE, Arch и NixOS — резервные записи 4.0 удалены, поскольку v3 во время компиляции требует API 4.1 (#5071)
- Исправлено имя пакета webkit2gtk для openSUSE в doctor (`webkit2gtk4_1-devel` → `webkit2gtk3-devel` — правильное имя пакета openSUSE) (#5071)
- Исправлена ошибка `Unexpected token '<'` при отсутствии `/wails/custom.js` в режиме разработки настольного приложения. Добавлены явный обработчик 404 для `/wails/custom.js` и проверка `Content-Type` без учёта регистра в `loadOptionalScript`, чтобы предотвратить внедрение резервных HTML-страниц SPA в качестве JavaScript. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Исправлено состояние подсветки меню в области уведомлений macOS: теперь значок отображается как выбранный, когда меню открыто (#4910)
- Исправлено появление привязанного к области уведомлений окна позади других окон в macOS: теперь используется корректный уровень всплывающего окна (#4910)
- Исправлены неверные примеры импорта `@wailsio/runtime` в документации (#4989)
- Исправлена невозможность свернуть безрамочное окно в darwin (#4294)
- Устранены зависания на 20-30 минут во время `wails3 build` и `wails3 dev` путём исключения `node_modules/` из проверки актуальности в go-task. Ранее glob-шаблон `sources: "**/*"` заставлял go-task перечислять и вычислять контрольную сумму каждого файла в `node_modules/` (более 50000-100000 файлов с тяжёлыми зависимостями, такими как MUI), что было особенно медленно в Windows/NTFS (#4939)
- Исправлена ошибка сборки GTK4, вызванная конфликтом C-объявления `Screen` typedef с X11 Xlib.h (#4957)
- Унифицировано поведение методов значка-счётчика в Dock на macOS
- Исправлено применение `InvisibleTitleBarHeight` ко всем окнам macOS вместо применения только к безрамочным окнам или окнам с прозрачной строкой заголовка (#4960)
- Устранено дрожание окна при изменении размера за верхние углы с включённым `InvisibleTitleBarHeight`: запуск перетаскивания теперь пропускается вблизи краёв окна (#4960)
- Исправлена генерация отображаемых типов с ключами-перечислениями в привязках JS/TS (#4437), автор — @fbbdev
- Исправлено перетаскивание файлов в Windows, не работавшее при масштабе экрана, отличном от 100%
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом приёме файлов в Windows
- Исправлено использование неверного пространства пикселей для координат приёма файлов в Windows (физические пиксели вместо CSS-пикселей)
- Исправлена ненадёжная работа перетаскивания файлов в Linux при наличии эффектов наведения
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом приёме файлов в Linux
- Исправлено восстановление окна в свёрнутом состоянии, иногда происходившее при его отображении или скрытии в Linux/GTK4: теперь используется `gtk_window_present()` (#4957)
- Исправлены получение и задание положения окна в Linux/GTK4, всегда возвращавшие 0,0: добавлена условная поддержка X11 через `XTranslateCoordinates`/`XMoveWindow` (#4957)
- Исправлено несоблюдение максимального размера окна в Linux/GTK4: вместо удалённого `gtk_window_set_geometry_hints` добавлено ограничение размера на основе сигналов (#4957)
- Исправлено масштабирование DPI в Linux/GTK4: реализованы корректное вычисление PhysicalBounds и поддержка дробного масштабирования через `gdk_monitor_get_scale` (GTK 4.14+)
- Исправлено дублирование пунктов меню при создании новых окон в Linux/GTK4
- Исправлена генерация отображаемых типов с ключами-перечислениями в привязках JS/TS (#4437), автор — @fbbdev
- Исправлено перетаскивание файлов в Windows, не работавшее при масштабе экрана, отличном от 100%
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом приёме файлов в Windows
- Исправлено использование неверного пространства пикселей для координат приёма файлов в Windows (физические пиксели вместо CSS-пикселей)
- Исправлена ненадёжная работа перетаскивания файлов в Linux при наличии эффектов наведения
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом приёме файлов в Linux
- Исправлено масштабирование DPI в Linux/GTK4: реализованы корректное вычисление PhysicalBounds и поддержка дробного масштабирования через `gdk_monitor_get_scale` (GTK 4.14+)
- Исправлено дублирование пунктов меню при создании новых окон в Linux/GTK4
- Исправлена генерация отображаемых типов с ключами-перечислениями в привязках JS/TS (#4437), автор — @fbbdev
- Исправлена проблема «окон-призраков» в macOS, вызванная обращением к API AppKit не из главного потока в App.Window.Current() (#4947), автор — @wimaha
- Исправлена неработоспособность HTML-элемента `<input type="file">` в macOS путём реализации WKUIDelegate runOpenPanelWithParameters (#4862)
- Исправлена неработоспособность нативного перетаскивания файлов при использовании npm-модуля `@wailsio/runtime` в macOS/Linux (#4953), автор — @leaanthony
- Исправлена генерация привязок для псевдонимов типов из других пакетов (#4578), автор — @fbbdev
- Исправлен аварийный сбой OpenFileDialog в Linux из-за нарушения потокобезопасности GTK (#3683), автор — @ddmoney420
- Исправлен аварийный сбой SIGSEGV при вызове `Focus()` для скрытого или уничтоженного окна (#4890), автор — @ddmoney420
- Устранена потенциальная паника при задании пустого значка или растрового изображения в Linux (#4923), автор — @ddmoney420
- Исправлен аварийный сбой ErrorDialog при вызове из привязки службы в macOS (#3631), автор — @leaanthony
- Обеспечено отображение меню в ОС Windows в `v3\examples\dialogs`, автор — @ndianabasi
- Исправлено состояние гонки, вызывавшее TypeError при перезагрузке страницы (#4872), автор — @ddmoney420
- Исправлен неверный результат тестов генератора привязок путём удаления глобального состояния из метода `Collector.IsVoidAlias()` (#4941), автор — @fbbdev
- Исправлена неработоспособность средства выбора файлов `<input type="file">` в macOS (#4862), автор — @leaanthony
- Исправлено использование несовместимых систем координат в `Position()` и `SetPosition()` в macOS, из-за которого положение окна смещалось при сохранении и восстановлении состояния (#4816), автор — @leaanthony
- Исправлена ошибка SetProcessDpiAwarenessContext «Access is denied», возникавшая, когда поддержка DPI уже была настроена через манифест приложения (#4803)
- Обновлена страница документации о сочетаниях клавиш и исправлен тип параметра функции обратного вызова для `KeyBinding.Add`, автор — @ndianabasi
- Исправлена документация по генерации пользовательской привязки: необходимо использовать `-d String` вместо `-o String`
- Исправлена ошибка, из-за которой дочерние элементы меню не очищались при `menu.Update()`
- Исправлены устаревшие ссылки на API Manager в документации (файлы 31 обновлены для использования нового шаблона, например `app.Window.New()`, `app.Event.Emit()` и т. д.), автор — @leaanthony
- Исправлен аварийный сбой Linux при панике в привязанных к JS методах Go, вызванный переопределением обработчиков сигналов в WebKit (#3965), автор — @leaanthony
- Исправлена ошибка, из-за которой SaveFileDialog.SetFilename() не действовал в Linux (#4841), автор — @samstanier
- Исправлено отображение координат приёма как undefined в примере drag-and-drop
- Исправлена ошибка создания пакета приложения macOS, когда APP_NAME содержит пробелы (проблема с раскрытием фигурных скобок)
- Исправлена паника из-за выхода индекса за границы массива в Windows при вызове методов служб (откачен goccy/go-json)
- Исправлено перетаскивание файлов в Windows, не работавшее при масштабе экрана, отличном от 100%
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом приёме файлов в Windows
- Исправлено использование неверного пространства пикселей для координат приёма файлов в Windows (физические пиксели вместо CSS-пикселей)
- Исправлена ненадёжная работа перетаскивания файлов в Linux при наличии эффектов наведения
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом приёме файлов в Linux
- Обновлены все команды в файлах Taskfile.yml для всех операционных систем, чтобы учитывать пробелы в таких переменных, как `APP_NAME`, автор — @ndianabasi
- Исправлена ошибка аргумента команды при выполнении задачи 'build:universal:lipo:go' в Linux, автор — @wux1an
- Исправлена ошибка Docker «undefined symbol: **<em>ubsan</em>handle_xxxxxxx» при выполнении 'wails3 build GOOS=darwin GOARCH=arm64' в Linux — @wux1an
- Объединена документация по пользовательским протоколам и добавлены разделы об Universal Links — @leaanthony
- Исправлен сбой меню в области уведомлений Windows при многократном нажатии на значок: добавлена защита от одновременных вызовов TrackPopupMenuEx (#4151) — @leaanthony
- Предотвращён сбой приложения при вызове systray.Run() до app.Run() — @leaanthony
- Исправлен сбой в macOS при переключении видимости окна через Hide()/Show(), когда включён параметр ApplicationShouldTerminateAfterLastWindowClosed (#4389) — @leaanthony
- Исправлена утечка памяти в контекстных меню macOS и Windows при их многократном открытии (#4012) — @leaanthony
- Исправлена проблема, из-за которой нативные ресурсы контекстного меню не использовались повторно в macOS и при каждом отображении создавалось новое меню (#4012) — @leaanthony
- Исправлена проблема в macOS, из-за которой нажатие на значок в Dock не отображало скрытые окна, если приложение запускалось с `Hidden: true` (#4583) — @leaanthony
- Исправлена проблема в macOS, из-за которой диалог печати не открывался вследствие неверного типа указателя окна в вызове CGO (#4290) — @leaanthony
- Исправлен сбой меню окна в Wayland, вызванный обращением appmenu-gtk-module к ещё не созданному окну (#4769) — @leaanthony
- Исправлен сбой приложения GTK, если имя приложения содержит недопустимые символы (пробелы, круглые скобки и т. д.), — @leaanthony
- Исправлена ошибка «not enough memory» при инициализации перетаскивания в Windows (#4701) — @overlordtm
- Исправлена проблема, из-за которой файловый менеджер открывал неверный каталог в Linux вследствие неправильного экранирования URI (#4397) — @leaanthony
- Исправлена ошибка сборки AppImage в современных дистрибутивах Linux (Arch, Fedora 39+, Ubuntu 24.04+): добавлено автоматическое обнаружение разделов ELF `.relr.dyn` и отключение удаления символов (#4642) — @leaanthony
- Исправлена проблема, из-за которой `wails doctor` ошибочно сообщал об установленных пакетах webkit в Fedora и системах на основе DNF (#4457) — @leaanthony
- Исправлено поведение по умолчанию, при котором `config.yml` запускал `wails3 dev` со сборкой для производственной среды, — @mbaklor
- Исправлены заглушки сервисов iOS, приводившие к сбоям сборки из-за импорта несуществующего пакета, — @leaanthony
- Исправлено структурированное журналирование в методах debug/info, вызывавшее ошибки «no formatting directives», — @leaanthony
- Удалены временные отладочные операторы печати, случайно добавленные при слиянии кода мобильной платформы, — @leaanthony
- Исправлен сбой WebKitGTK в Wayland на графических процессорах NVIDIA (Error 71 Protocol error): рендерер DMA-BUF теперь автоматически отключается — @leaanthony
- Исправлена проблема, из-за которой альфа-канал игнорировался в `application.WebviewWindowOptions.BackgroundColour` в Linux ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Исправлена проблема, из-за которой при отсутствии пользовательского значка в области уведомлений Windows по умолчанию не использовался значок приложения (#4704)
- Добавлено отслеживание владельца `HICON`, чтобы уничтожались только созданные пользователем дескрипторы, что предотвращает сбои при перезапуске Explorer (#4653).
- При уничтожении теперь освобождаются обработчик системной темы Windows и сохранённые значки области уведомлений, что устраняет утечки горутин и контекстов устройств (#4653).
- Подсказки значков области уведомлений теперь обрезаются до 127 единиц UTF-16, чтобы не повреждать суррогатные пары и многобайтовые глифы (#4653).
- Исправлена ошибка задачи упаковки для Windows (#4667)
- Исправлена переменная appicon для AppImage в Linux в файле задач Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Исправлена ошибка сборки для Windows, вызванная изменением сигнатуры go-webview2 v1.0.22 (#4513, #4645)
- Исправлена переменная appicon для AppImage в Linux в файле задач Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Исправлен диапазон протоколов в linux desktop.tmpl: удалено от `<.Info.Protocol>` до `<.Protocol>` — @Tolfx в #4510
- Исправлена ошибка повторного определения в демонстрационном примере liquid glass в [#4542](https://github.com/wailsapp/wails/pull/4542) — @Etesam913
- Исправлено обновление меню в области уведомлений Linux [#4604](https://github.com/wailsapp/wails/issues/4604) — [@JackDoan](https://github.com/JackDoan)
- Исправлено появление белого окна в Windows при создании скрытого окна — @leaanthony в [#4612](https://github.com/wailsapp/wails/pull/4612)
- Исправлен путь импорта пакета уведомлений в документации — @rxliuli в [#4617](https://github.com/wailsapp/wails/pull/4617)
- Исправлена проблема, из-за которой перетаскивание не работало при использовании npm-пакета @wailsio/runtime (#4489), — @leaanthony в #4616
- Windows: устранены мерцание окна при запуске и неправильное отображение скрытых окон в [PR](https://github.com/wailsapp/wails/pull/4600) — @leaanthony.
- Исправлены проблемы с размером окна при разворачивании в Wayland (https://github.com/wailsapp/wails/issues/4429) — [@samstanier](https://github.com/samstanier)
- Исправлены проблемы с размером окна при разворачивании в Wayland (https://github.com/wailsapp/wails/issues/4429) — [@samstanier](https://github.com/samstanier)
- Исправлена ошибка повторного определения в демонстрационном примере liquid glass в [#4542](https://github.com/wailsapp/wails/pull/4542) — @Etesam913
- Исправлена проблема, из-за которой AssetServer мог аварийно завершить работу в macOS, в [#4576](https://github.com/wailsapp/wails/pull/4576) — @jghiloni
- Исправлена ошибка компиляции при сборке с NextJs. Исправление внесено в [#4585](https://github.com/wailsapp/wails/pull/4585) — @rev42
- Исправлены конвейеры ночного выпуска в [#4597](https://github.com/wailsapp/wails/pull/4597) — @riadafridishibly
- Исправлена ошибка повторного определения в демонстрационном примере liquid glass в [#4542](https://github.com/wailsapp/wails/pull/4542) — @Etesam913
- Исправлена проблема, из-за которой AssetServer мог аварийно завершить работу в macOS, в [#4576](https://github.com/wailsapp/wails/pull/4576) — @jghiloni
- Исправлена ошибка компиляции при сборке с NextJs. Исправление внесено в [#4585](https://github.com/wailsapp/wails/pull/4585) — @rev42
- Исправлены конвейеры ночного выпуска в [#4597](https://github.com/wailsapp/wails/pull/4597) — @riadafridishibly
- Исправлена ошибка повторного определения в демонстрационном примере liquid glass в [#4542](https://github.com/wailsapp/wails/pull/4542) — @Etesam913
- Исправлен SetBackgroundColour в Windows — @PPTGamer в [PR](https://github.com/wailsapp/wails/pull/4492)
- Документация обновлена с учётом изменений после рефакторинга Manager API — @yulesxoxo в [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Исправлена переменная appicon файла .desktop в Linux в файле задач Linux [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Документация обновлена с учётом изменений после рефакторинга Manager API — @yulesxoxo в [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Исправлена ошибка разыменования нулевого указателя в Windows, о которой сообщалось в [#4456](https://github.com/wailsapp/wails/issues/4456), — @leaanthony в [#4460](https://github.com/wailsapp/wails/pull/4460)
- В WKWebView для macOS добавлена поддержка `allowsBackForwardNavigationGestures`, позволяющая использовать жесты навигации двумя пальцами (#1857)
- Исправлена проблема, из-за которой onClick не работал для пунктов меню, изначально помеченных как отключённые, — @leaanthony в [PR #4469](https://github.com/wailsapp/wails/pull/4469). Благодарим @IanVS за первоначальное исследование.
- Исправлена проблема, из-за которой сервер Vite не завершал работу при сбое сборки (#4403)
- Исправлена паника при закрытии или отмене `SaveFileDialog` в Windows. Исправление внесено в [PR](https://github.com/wailsapp/wails/pull/4284) — @hkhere
- Исправлено перетаскивание на уровне HTML в Windows — [@mbaklor](https://github.com/mbaklor) в [#4259](https://github.com/wailsapp/wails/pull/4259)
- Добавлена поддержка `allowsBackForwardNavigationGestures` в WKWebView для macOS, чтобы включить жесты навигации смахиванием двумя пальцами (#1857)
- Исправлена ошибка, из-за которой onClick не работал для пунктов меню, изначально заданных как отключённые, — @leaanthony в [PR #4469](https://github.com/wailsapp/wails/pull/4469). Спасибо @IanVS за первоначальное исследование.
- Исправлена ошибка, из-за которой сервер Vite не завершал работу при сбое сборки (#4403)
- Исправлен разбор уведомлений в Windows — @popaprozac в [PR](https://github.com/wailsapp/wails/pull/4450)
- В команде doctor исправлена проверка зависимостей Windows SDK — [@kodumulo](https://github.com/kodumulo) в [#4390](https://github.com/wailsapp/wails/issues/4390)
- Исправлено разыменование nil-указателя в processURLRequest на Mac — [@etesam913](https://github.com/etesam913) в [#4366](https://github.com/wailsapp/wails/pull/4366)
- Исправлена ошибка Linux, препятствовавшая работе диалоговых окон с фильтрами, — [@bh90210](https://github.com/bh90210) в [#4287](https://github.com/wailsapp/wails/pull/4287)
- Исправлены проблемы с меню Edit в Windows и Linux — [@leaanthony](https://github.com/leaanthony) в [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Минимальная версия системы в файлах .plist для macOS обновлена с 10.13.0 до 10.15.0 — [@AkshayKalose](https://github.com/AkshayKalose) в [#3981](https://github.com/wailsapp/wails/pull/3981)
- Исправлена проблема с пропуском идентификаторов окон — [@leaanthony](https://github.com/leaanthony)
- Исправлена проблема с nil-меню при вызове RegisterContextMenu — [@leaanthony](https://github.com/leaanthony)
- Исправлены циклические зависимости в выходных данных генератора привязок — [@fbbdev](https://github.com/fbbdev) в [#4001](https://github.com/wailsapp/wails/pull/4001)
- Исправлены ошибки использования до объявления в выходных данных генератора привязок — [@fbbdev](https://github.com/fbbdev) в [#4001](https://github.com/wailsapp/wails/pull/4001)
- Флаги сборки теперь передаются генератору привязок — [@fbbdev](https://github.com/fbbdev) в [#4023](https://github.com/wailsapp/wails/pull/4023)
- В путях файла Taskfile для Windows обратные косые черты заменены прямыми, чтобы обеспечить работу на платформах, отличных от Windows, — [@leaanthony](https://github.com/leaanthony)
- Исправлены события Mac и Mac JS — [@leaanthony](https://github.com/leaanthony)
- Исправлена взаимная блокировка событий в macOS — [@leaanthony](https://github.com/leaanthony)
- Исправлена ошибка `Parameter incorrect` при инициализации Window в Windows, когда предоставлен HTML, но отсутствует JS, — [@leaanthony](https://github.com/leaanthony)
- Исправлен размер префикса ответа, используемого сервером ресурсов для определения типа содержимого, — [@fbbdev](https://github.com/fbbdev) в [#4049](https://github.com/wailsapp/wails/pull/4049)
- Исправлена обработка ответов, отличных от 404, для корневого пути индекса на сервере ресурсов — [@fbbdev](https://github.com/fbbdev) в [#4049](https://github.com/wailsapp/wails/pull/4049)
- Исправлено неопределённое поведение генератора привязок при проверке свойств обобщённых типов — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлены выходные данные генератора привязок для моделей, у которых базовый тип имеет свойства, отличные от свойств именованной обёртки, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлены выходные данные генератора привязок для типов ключей карт и предварительная обработка — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлены выходные данные генератора привязок для структур, реализующих интерфейсы маршалинга, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлено обнаружение циклов типов с участием обобщённых типов в генераторе привязок — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлены недопустимые ссылки на неэкспортируемые модели в выходных данных генератора привязок — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Внедряемый код перемещён в конец файлов сервисов — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлена обработка ошибок операций закрытия файлов в генераторе привязок — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Подавлены предупреждения для сервисов, которые определяют методы жизненного цикла или HTTP, но не имеют других привязанных методов, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлена ошибка, из-за которой шаблоны без React не отображали нижний колонтитул Hello World при использовании светлой системной цветовой схемы, — [@marcus-crane](https://github.com/marcus-crane) в [#4056](https://github.com/wailsapp/wails/pull/4056)
- Исправлены скрытые пункты меню в macOS — [@leaanthony](https://github.com/leaanthony)
- Исправлены обработка и форматирование ошибок в обработчиках сообщений — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Исправлен пропуск завершения работы сервисов при выходе из приложения — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Обновления меню теперь гарантированно выполняются в главном потоке — [@leaanthony](https://github.com/leaanthony)
- Механизм перетаскивания и изменения размера стал надёжнее и теперь точнее соответствует ожидаемому поведению платформы — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Исправлена ошибка [#4097](https://github.com/wailsapp/wails/issues/4097): Webpack/angular удаляет код инициализации среды выполнения — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Исправлены изначально скрытые пункты меню — [@IanVS](https://github.com/IanVS) в [#4116](https://github.com/wailsapp/wails/pull/4116)
- Исправлена ошибка, из-за которой assetFileServer не обслуживал файлы `.html` при запросе без расширения, если `[request]` не существует, но `[request].html` существует
- Исправлены пути генерации значков — [@robin-samuel](https://github.com/robin-samuel) в [#4125](https://github.com/wailsapp/wails/pull/4125)
- Исправлена ошибка, из-за которой события `fullscreen`, `unfullscreen`, `unminimise` и `unmaximise` не генерировались, — [@oSethoum](https://github.com/osethoum) в [#4130](https://github.com/wailsapp/wails/pull/4130)
- Исправлена ошибка NSIS, вызванная неправильным префиксом версии по умолчанию в конфигурации, — [@robin-samuel](https://github.com/robin-samuel) в [#4126](https://github.com/wailsapp/wails/pull/4126)
- Исправлена ошибка, из-за которой функция среды выполнения Dialogs возвращала экранированные пути в Windows, — [TheGB0077](https://github.com/TheGB0077) в [#4188](https://github.com/wailsapp/wails/pull/4188)
- Исправлен путь обнаружения Webview2 в HKCU — [@leaanthony](https://github.com/leaanthony).
- Исправлена проблема с вводом в macOS — [@leaanthony](https://github.com/leaanthony).
- Исправлено имя файла задачи генерации значков Windows — [@yulesxoxo](https://github.com/yulesxoxo) в [#4219](https://github.com/wailsapp/wails/pull/4219).
- Исправлена проблема с прозрачностью окон без рамки — [@leaanthony](https://github.com/leaanthony) на основе работы @kron.
- Исправлены вызовы фокусировки, когда окно отключено или свёрнуто, — [@leaanthony](https://github.com/leaanthony) на основе работы @kron.
- Исправлена ошибка, из-за которой значки в области уведомлений не появлялись после перезапуска панели задач, — [@leaanthony](https://github.com/leaanthony) на основе работы @kron.
- Исправлено отсутствие реализации Flush() в fallbackResponseWriter в [#4245](https://github.com/wailsapp/wails/pull/4245)
- Исправлено отсутствие реализации Flush() в fallbackResponseWriter; автор — [@superDingda], [#4236](https://github.com/wailsapp/wails/issues/4236)
- Исправлены сбои в macOS при закрытии окна во время ожидающего асинхронного вызова привязанной функции Go; автор — [@joshhardy](https://github.com/joshhardy), [#4354](https://github.com/wailsapp/wails/pull/4354)
- Исправлена гонка при запуске режима эффективности Windows; автор — [@leaanthony](https://github.com/leaanthony)
- Исправлено освобождение дескрипторов значков Windows; автор — [@leaanthony](https://github.com/leaanthony).
- Исправлено `OpenFileManager` в Windows; автор — [@PPTGamer](https://github.com/PPTGamer), [#4375](https://github.com/wailsapp/wails/pull/4375).
- Исправлены параметры минимальной и максимальной ширины в Linux; автор — @atterpac, [#3979](https://github.com/wailsapp/wails/pull/3979)
- Определения типов шаблонов TypeScript обновлены посредством повышения версии npm; автор — @atterpac, [#3966](https://github.com/wailsapp/wails/pull/3966)
- Исправлена ссылка на CSS в шаблоне SvelteKit; автор — @atterpac, [#3945](https://github.com/wailsapp/wails/pull/3945)
- Обеспечен вызов ключевых функций обратного вызова в window run() в основном потоке; автор — [@leaanthony](https://github.com/leaanthony)
- Исправлены примеры диалогового окна выбора каталога; автор — [@leaanthony](https://github.com/leaanthony)
- Создана новая страница ошибки на китайском языке для случая отсутствия index.html; автор — [@leaanthony](https://github.com/leaanthony)
-  Обеспечено выполнение функции обратного вызова `windowDidBecomeKey` в основном потоке; автор — [@leaanthony](https://github.com/leaanthony)
-  Добавлена поддержка полноэкранного режима для окон без рамки; автор — [@leaanthony](https://github.com/leaanthony)
-  Улучшена логика уничтожения окон; автор — [@leaanthony](https://github.com/leaanthony)
-  Исправлена логика позиционирования окна, прикреплённого к значку в системной области уведомлений; автор — [@leaanthony](https://github.com/leaanthony)
-  Добавлена поддержка полноэкранного режима для окон без рамки; автор — [@leaanthony](https://github.com/leaanthony)
- Исправлена обработка событий; автор — [@leaanthony](https://github.com/leaanthony)
- Исправлена логика завершения работы окон; автор — [@leaanthony](https://github.com/leaanthony)
- Общий Taskfile теперь по умолчанию создаёт привязки TypeScript для шаблонов TypeScript; автор — [@leaanthony](https://github.com/leaanthony)
- Исправлено закрытие приложения при получении сообщения WM_CLOSE, когда нет открытых окон или используется только системная область уведомлений; автор — [@mmalcek](https://github.com/mmalcek), [#3990](https://github.com/wailsapp/wails/pull/3990)
- Исправлена сборка с garble; автор — @5aaee9, [#3192](https://github.com/wailsapp/wails/pull/3192)
- Исправлены сборки NSIS для Windows; автор — [@leaanthony](https://github.com/leaanthony)
- Исправлена взаимная блокировка в диалоговом окне Linux при выборе нескольких элементов, вызванная незакрытым
- Исправлена кроссплатформенная очистка файлов .syso во время сборки для Windows; автор —
- Исправлена компиляция AppImage для amd64; автор — @atterpac, в
- Исправлено обновление ресурсов сборки; автор — @ansxuman, в
- Исправлена реализация `OnClick` и `OnRightClick` для системной области уведомлений Linux; автор — @atterpac
- Исправлено неработающее `AlwaysOnTop` на Mac; автор —
-  Исправлено включение дубликата в `application.NewEditMenu`
- 🐧 Исправлена компиляция для aarch64
- ⊞ Исправлены пункты меню с переключателями; автор —
- Исправлена ошибка при сборке запускаемого файла .app в macOS, когда значения 'name' и 'outputfilename'
- Исправлена ошибка использования customEventProcessor в примере перетаскивания; автор —
- 🐧 Исправлена ошибка компиляции в Linux, появившаяся после добавления IgnoreMouseEvents; автор —
- ⊞ Исправлена ошибка создания файла значка syso; автор —
- 🐧 Интегрировано исправление для нативной работы в Wayland из
- Не привязывать внутренние методы служб в
- ⊞ Исправлена паника при запуске системной области уведомлений в
- Не привязывать внутренние методы служб в
- ⊞ Исправлена паника при запуске системной области уведомлений в
- Выполнен крупный рефакторинг пунктов меню и обработки событий. Пока в основном улучшена macOS. Автор —
- Исправлены тесты после рефакторинга плагинов и событий в
- ⊞ Исправлено предупреждение `Failed to unregister class Chrome_WidgetWin_0`. Автор —
- Проблемы с модулями
- Исправлена передача сообщений о событиях изменения размера; автор — [atterpac](https://github.com/atterpac), в
- 🐧 Исправлена ошибка обработки темы в NixOS; автор —
- Исправлена установка проекта между томами в Windows; автор —
- Исправлен CSS шаблона React для отображения нижнего колонтитула; автор —
- Исправлены процессы-зомби при работе в режиме разработки посредством обновления refresh до последней версии
- Исправлено получение файла WebKit для AppImage; автор — [Atterpac](https://github.com/atterpac)
- Исправлена проверка пакета apt в Doctor; автор — [Atterpac](https://github.com/Atterpac), в
- Исправлено зависание приложения при выходе в Darwin; автор — @5aaee9, в
- Исправлены цвета фона примеров в Windows; автор —
- Исправлены контекстные меню по умолчанию; автор — [mmghv](https://github.com/mmghv), в
- Исправлены шестнадцатеричные значения клавиш со стрелками в Darwin; автор —
- Перетаскивание в Windows приведено в рабочее состояние. Добавил
- Исправлена ошибка в Doctor для Linux, возникавшая при отсутствии у пользователя необходимых драйверов
- Исправлено масштабирование DPI при запуске в Windows. Автор изменения — [@almas-x](https://github.com/almas-x), в
- Исправлена строка замены в `go.mod`: теперь используются относительные пути. Это исправляет пути Windows с
- Исправлена обработка щелчка по значку в системной области уведомлений macOS при отсутствии прикреплённого окна; автор —
- Исправлен сбой сборки для Windows из-за неизвестного параметра; автор —
- Исправлен сбой в Windows при щелчке левой кнопкой мыши по значку в системной области уведомлений при отсутствии
- Исправлен неверный baseURL при повторном открытии окна; автор — @5aaee9, в PR
- Исправлен порядок ветвей if в методе `WebviewWindow.Restore`; автор —
- Обеспечено правильное вычисление `startURL` при нескольких вызовах `GetStartURL`, когда
- Исправлен JS-тип структуры `Screen`, чтобы он соответствовал аналогу в Go; автор —
- Исправлен метод `WML.Reload`, чтобы обеспечить надлежащую очистку зарегистрированных обработчиков событий
- Исправлено немедленное закрытие пользовательского контекстного меню в Linux; автор —
- Исправлены выходной путь и расширение файлов моделей, создаваемых генератором привязок
- Исправлены пути импорта файлов моделей в JS-коде, создаваемом генератором привязок
- Исправлено перетаскивание в некоторых дистрибутивах Linux; автор —
- Исправлено отсутствие задачи для macOS при использовании `wails3 task dev`; автор —
- Исправлена ошибка присваивания nil-карте при регистрации событий; автор —
- Исправлена десериализация параметров привязанных методов; автор —
- Исправлена обработка нескольких возвращаемых значений привязанных методов, автор:
- Исправлено обнаружение в doctor установки npm не через системный менеджер пакетов
- Исправлено отсутствие MicrosoftEdgeWebview2Setup.exe. Спасибо
- Исправлен случайный сбой в Linux из-за обработки идентификатора окна, автор: @leaanthony. На основе
- Исправлен сбой systemTray.setIcon в Linux, автор:
- Исправлено применение рамки окна при первом вызове функции `setFrameless` в

### Изменено

- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ**: ключи map в сгенерированных привязках JS/TS теперь помечаются как необязательные, чтобы точно отражать семантику map в Go. Доступ к значениям map в TypeScript теперь возвращает `T | undefined` вместо `T`, поэтому требуется проверка на null или утверждение типа (#4943), автор: `@fbbdev`
- Использование `Event` заменено на `Events` в соответствии с изменениями в `@wailsio/runtime`; соответствующие вызовы функций в документации также обновлены в `Features/Events/Event System`, автор: @AbdelhadiSeddar
- Параметры `EnabledFeatures`, `DisabledFeatures` и `AdditionalBrowserArgs` перенесены из параметров отдельных окон в `Options.Windows` уровня приложения (#4559), автор: @leaanthony
- Обновлён README для примера `Drag N Drop` и подчёркнуто, что в нём демонстрируется `Internal Drag and Drop`, автор: @ndianabasi
- Уровень различных отладочных сообщений изменён с Info на Debug (автор: @mbaklor)
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** `EnableDragAndDrop` переименован в `EnableFileDrop` в параметрах окна
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** `DropZoneDetails` переименован в `DropTargetDetails` в контексте события
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** метод `DropZoneDetails()` объекта `WindowEventContext` переименован в `DropTargetDetails()`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** событие `WindowDropZoneFilesDropped` удалено; вместо него используйте `WindowFilesDropped`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** HTML-атрибут изменён с `data-wails-dropzone` на `data-file-drop-target`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** CSS-класс наведения изменён с `wails-dropzone-hover` на `file-drop-target-active`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** из Windows удалены параметры `DragEffect`, `OnEnterEffect` и `OnOverEffect` (они относились к удалённому IDropTarget)
- Вся обработка JSON во время выполнения (привязки методов, события, запросы webview, уведомления, kvstore) переведена на goccy/go-json, что повысило производительность на 21-63% и сократило выделение памяти на 40-60%
- Оптимизирована компоновка структуры BoundMethod и добавлено кеширование флага isVariadic для снижения накладных расходов каждого вызова
- Для методов с `<=8` аргументами используется буфер аргументов, выделяемый в стеке, чтобы избежать выделений в куче
- Оптимизирован сбор результатов вызовов методов, чтобы не выделять срез для одного возвращаемого значения
- Для кеша MIME-типов используется sync.Map, что повышает производительность при конкурентном доступе
- Для чтения тела запроса транспорта HTTP используется пул буферов
- Канал CloseNotify в определителе типа содержимого теперь выделяется отложенно, что сокращает число выделений памяти на каждый запрос
- Из сервера ресурсов удалено отладочное журналирование CSS
- Карта расширений MIME-типов дополнена и теперь охватывает более 50 распространённых веб-форматов (шрифты, аудио, видео и т. д.)
- Обновлена документация по параметрам окна `X/Y`, автор: @ruhuang2001
- Документация `Frontend Runtime` дополнена параметрами генерации привязок для фронтенда, автор: @ndianabasi
- Обновлена страница документации сервера ресурсов Wails v3, автор: @ndianabasi
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ**: функции диалогов уровня пакета (`application.InfoDialog()`, `application.QuestionDialog()` и т. д.) удалены. Вместо них используйте диспетчер `app.Dialog`: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Документация по диалогам приведена в соответствие с фактическим API: используются `app.Dialog.*` и `AddButton()` с функциями обратного вызова (не `SetButtons()`), `SetDefaultButton(*Button)` (не строка), `AddFilter()` (не `SetFilters()`), `SetFilename()` (не `SetDefaultFilename()`), а для выбора папки — `app.Dialog.OpenFile().CanChooseDirectories(true)`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ**: теперь по умолчанию создаются производственные сборки. Чтобы создавать сборки для разработки, задайте `DEV=true` в своих Taskfile. Для получения примеров создайте новый проект, автор: @leaanthony
- При отправке пользовательского события без аргументов данных или с одним таким аргументом значение данных теперь записывается непосредственно в поле Data без оборачивания в срез, автор: [@fbbdev](https://github.com/fbbdev) в [#4633](https://github.com/wailsapp/wails/pull/4633)
- Системные треи Windows теперь учитывают `SystemTray.Show()`/`Hide()`, переключая `NIS_HIDDEN`, поэтому приложения действительно могут исчезать и возвращаться (#4653).
- При регистрации системного трея повторно используются уже разрешённые значки, `NOTIFYICON_VERSION_4` задаётся один раз, а `NIF_SHOWTIP` включается, чтобы всплывающие подсказки восстанавливались после перезапуска Explorer (#4653).
- macOS: для центрирования окна используется `visibleFrame` вместо `frame`, чтобы исключить области строки меню и Dock
- macOS: для центрирования окна используется `visibleFrame` вместо `frame`, чтобы исключить области строки меню и Dock
- При запуске `wails3 update build-assets` с параметром `-config` значения, заданные через параметры `-product*`,
- `window.NativeWindowHandle()` -> `window.NativeWindow()`, автор: @leaanthony в [#4471](https://github.com/wailsapp/wails/pull/4471)
- Переработана внутренняя обработка окон, автор: @leaanthony в [#4471](https://github.com/wailsapp/wails/pull/4471)
- Удалены `application.WindowIDKey` и `application.WindowNameKey` (заменены на `application.WindowKey`), автор: [@leaanthony](https://github.com/leaanthony)
- ContextMenuData теперь возвращает строку вместо any, автор: [@leaanthony](https://github.com/leaanthony)
- В привязках JS/TS поля классов с типами массивов фиксированной длины теперь инициализируются с ожидаемой длиной, а не остаются пустыми, автор: [@fbbdev](https://github.com/fbbdev) в [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData теперь возвращает строку вместо any, автор: [@leaanthony](https://github.com/leaanthony)
- `application.NewService` больше не принимает параметры в качестве необязательного аргумента (вместо этого используйте `application.NewServiceWithOptions`), автор: [@leaanthony](https://github.com/leaanthony) в [#4024](https://github.com/wailsapp/wails/pull/4024)
- Удалена зависимость `nanoid`, автор: [@leaanthony](https://github.com/leaanthony)
- Обновлён пример Window для стилей окон mica/acrylic/tabbed, автор: [@leaanthony](https://github.com/leaanthony)
- Из привязок JS/TS удалены файлы моделей `internal.js/ts`; теперь все модели находятся в `models.js/ts`, автор: [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- В привязках JS/TS именованные типы больше никогда не отображаются как псевдонимы других именованных типов; прежнее поведение теперь применяется только к псевдонимам, автор: [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- В привязках JS/TS в режиме классов поля структур, типом которых является параметр типа, помечаются как необязательные и никогда не инициализируются автоматически, автор: [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- ESLint удалён из шаблонов, автор: [@IanVS](https://github.com/IanVS) в [#4059](https://github.com/wailsapp/wails/pull/4059)
- Дата авторского права обновлена до 2025, автор: [@IanVS](https://github.com/IanVS) в [#4037](https://github.com/wailsapp/wails/pull/4037)
- Добавлена документация по event.Sender — [@IanVS](https://github.com/IanVS) в [#4075](https://github.com/wailsapp/wails/pull/4075)
- Поддержка Go 1.24 — [@leaanthony](https://github.com/leaanthony)
- Хуки `ServiceStartup` теперь вызываются при вызове `App.Run`, а не в `application.New` — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ошибки `ServiceStartup` теперь возвращаются из `App.Run`, а не завершают процесс — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Вызовы привязок и диалогов из JS теперь отклоняются с объектами ошибок вместо строк — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Улучшено позиционирование меню в области уведомлений Windows — [@leaanthony](https://github.com/leaanthony)
- Среда выполнения JS перенесена на TypeScript — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Среда выполнения инициализируется сразу после импорта; ждать загрузки окна больше не нужно — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Среда выполнения больше не экспортирует метод init. Для её инициализации можно использовать импорт только ради побочных эффектов — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Привязанные методы теперь возвращают `CancellablePromise`, который при отмене отклоняется с `CancelError`. Фактический результат вызова отбрасывается — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Типы встроенных сервисов теперь единообразно называются `Service` — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Функции создания встроенных сервисов с параметрами теперь единообразно называются `NewWithConfig` — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Метод `Select` сервиса `sqlite` теперь называется `Query` для согласованности с API Go — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Шаблоны: среда выполнения перемещена в "dependencies", файлы package.json упорядочены — [@IanVS](https://github.com/IanVS) в [#4133](https://github.com/wailsapp/wails/pull/4133)
- В режиме разработки создаёт пакеты приложений и добавляет к ним специальную подпись, чтобы сделать доступными некоторые API macOS — [@popaprozac](https://github.com/popaprozac) в [#4171](https://github.com/wailsapp/wails/pull/4171)
- Ресурсы сборки перемещены в каталоги для отдельных платформ — [@leaanthony](https://github.com/leaanthony)
- Файлы Taskfile перемещены и переименованы в каталоги для отдельных платформ — [@leaanthony](https://github.com/leaanthony)
- Значительно улучшена работа при отсутствии `index.html` — [@leaanthony](https://github.com/leaanthony)
- [Windows] Повышена производительность сворачивания и восстановления окна — [@leaanthony](https://github.com/leaanthony). Основано на исходном [PR](https://github.com/wailsapp/wails/pull/3955) от [562589540](https://github.com/562589540)
- Удалён параметр `ShouldClose` (вместо него зарегистрируйте хук для events.Common.WindowClosing) — [@leaanthony](https://github.com/leaanthony)
- [Windows] Уменьшено мерцание при открытии окна — [@leaanthony](https://github.com/leaanthony)
- Удалена `Window.Destroy`, поскольку эта функция предназначалась для внутреннего использования — [@leaanthony](https://github.com/leaanthony)
- События `WindowClose` переименованы в `WindowClosing` — [@leaanthony](https://github.com/leaanthony)
- Сборки фронтенда теперь используют окружение vite "development" или "production" в зависимости от типа сборки — [@leaanthony](https://github.com/leaanthony)
- Обновление до go-webview2 v1.19 — [@leaanthony](https://github.com/leaanthony)
- Обеспечено использование форка taskfile — @leaanthony
- Обновлён форк Taskfile для исправления проблем с версиями при установке с помощью
- Используется форк Taskfile для исправления проблем с версиями при установке с помощью
- `service.OnStartup` теперь завершает работу приложения при ошибке и запускает
- Переработана передача сообщений о щелчках в области уведомлений для лучшего соответствия действиям пользователя —
- Во встраиваемые ресурсы включён `all:frontend/dist` для поддержки фреймворков, которые создают
- Рефакторинг Taskfile — [leaanthony](https://github.com/leaanthony) в
- Обновление до `go-webview2` v1.0.16 —
- Исправлен тип `Screen`: теперь он включает `ID`, а не `Id` —
- Версия Wails в `go.mod.tmpl` обновлена для поддержки `application.ServiceOptions` —
- Исправлено определение имени сервиса — [windom](https://github.com/windom/) в
- mkdocs serve теперь использует docker — [leaanthony](https://github.com/leaanthony)
- Конфигурация разработки объединена в `config.yml` —
- Диалог из области уведомлений теперь по умолчанию использует значок приложения, если он доступен (Windows) —
- Улучшено отображение сведений о GPU и памяти в macOS —
- Удалены `WebviewGpuIsDisabled` и `EnableFraudulentWebsiteWarnings`
- Изменение API событий: `On`/`Emit` -> пользовательские события, `OnApplicationEvent` ->
- Исправление API событий в Linux — [TheGB0077](https://github.com/TheGB0077) в
- [CI] улучшены actions и добавлена возможность запускать actions также в форках и
- `AbsolutePosition()` переименован в `Position()` —
- Зависимость Linux WebKit изменена с webkitgtk2-4.0 на webkit2gtk-4.1, чтобы
- Встроенный сценарий среды выполнения JS теперь является модулем ESM: теги script, импортирующие его,
- Пакет `@wailsio/runtime` не публикует свой API в `window.wails`
- Модуль API окна `@wailsio/runtime/src/window` теперь предоставляет содержащий его
- API окна JS обновлён в соответствии с текущим `WebviewWindow` Go
- Генератор привязок теперь по умолчанию использует вызовы по идентификатору. Параметр CLI `-id`
- Новая структура кода привязок: ранее выходные файлы распределялись по каталогам
- Поле структуры `application.Options.Bind` переименовано в
- Новый синтаксис привязки сервисов: экземпляры сервисов теперь необходимо оборачивать в
- Отключение индикатора выполнения в нетерминальной среде или среде CI —

### Удалено

- **BREAKING**: из параметров `WindowsWindow` отдельных окон удалены `EnabledFeatures`, `DisabledFeatures` и `AdditionalLaunchArgs`. Вместо них используйте `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures` и `Options.Windows.AdditionalBrowserArgs` уровня приложения. Эти флаги глобально применяются к общей среде WebView2 (#4559) — @leaanthony
- Нативная реализация `IDropTarget` в Windows удалена в пользу подхода на основе JavaScript (соответствует поведению v2)
- Зависимость github.com/wailsapp/mimetype удалена в пользу расширенной карты расширений и http.DetectContentType из стандартной библиотеки, что уменьшает размер двоичного файла примерно на 1.2 МБ
- Удалена зависимость gopkg.in/ini.v1: вместо неё реализован минимальный парсер файлов .desktop для файлового менеджера Linux, что позволило сэкономить около 45 КБ
- Из кода среды выполнения удалена зависимость samber/lo: вместо неё используются пакет slices из стандартной библиотеки Go 1.21+ и минимальные внутренние вспомогательные функции, что позволило сэкономить около 310 КБ
- Из обработчика схем URL для Darwin удалены отладочные вызовы printf (#4834)
- **НЕСОВМЕСТИМОЕ ИЗМЕНЕНИЕ**: удалено событие `linux:WindowLoadChanged`. Чтобы определить, когда WebView завершил загрузку, используйте вместо него `linux:WindowLoadFinished` (#3896), автор — @leaanthony

### Несовместимые изменения

- **Рефакторинг API менеджеров**: плоская структура API приложения преобразована в упорядоченный набор менеджеров, чтобы улучшить организацию кода и упростить поиск нужных возможностей. Изменение внесено [@leaanthony](https://github.com/leaanthony) в [#4359](https://github.com/wailsapp/wails/pull/4359)
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
- Методы Service переименованы: `Name` -> `ServiceName`, `OnStartup` -> `ServiceStartup`, `OnShutdown` -> `ServiceShutdown`. Автор — [@leaanthony](https://github.com/leaanthony)
- Методы `Path` и `Paths` перенесены в пакет `application`. Автор — [@leaanthony](https://github.com/leaanthony)
- Меню приложения теперь доступно только в macOS. Автор — [@leaanthony](https://github.com/leaanthony)

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

### Критические изменения

### Добавлено

### Исправлено

### Изменено

## v3.0.0-alpha.9 - 2025-01-13

### Добавлено

### Исправлено

### Изменено

## v3.0.0-alpha.8.3 — 2024-12-07

### Изменено

## v3.0.0-alpha.8.2 — 2024-12-07

### Изменено

`go install`, автор: @leaanthony

## v3.0.0-alpha.8.1 — 2024-12-07

### Изменено

`go install`, автор: @leaanthony

## v3.0.0-alpha.8 — 2024-12-06

### Добавлено

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### Изменено

`service.OnShutdown`для всех ранее запущенных служб, автор: @atterpac, в [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac в [#3907](https://github.com/wailsapp/wails/pull/3907)   подпапках, автор: @atterpac, в [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) в [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) в [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (заменено параметрами `EnabledFeatures` и `DisabledFeatures`), автор: [leaanthony](https://github.com/leaanthony)

### Исправлено

переменная канала, автор: @michael-freling, в [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) в [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   в [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) в [#3841](https://github.com/wailsapp/wails/pull/3841)   роль `PasteAndMatchStyle` в меню редактирования на Darwin, автор: [johnmccabe](https://github.com/johnmccabe), в [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) в [#3854](https://github.com/wailsapp/wails/pull/3854), автор: [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   различаются, автор: @nickisworking, в [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 — 2024-09-18

### Добавлено

[mmghv](https://github.com/mmghv) в [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) и [leaanthony](https://github.com/leaanthony) в [#3570](https://github.com/wailsapp/wails/pull/3570)

### Изменено

События приложения `OnWindowEvent` -> События окна, автор: [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   ветки с префиксом `v3/` или `v3-`, автор: [stendler](https://github.com/stendler), в [#3747](https://github.com/wailsapp/wails/pull/3747)

### Исправлено

[etesam913](https://github.com/etesam913) в [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) в [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) в [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) в [#3614](https://github.com/wailsapp/wails/pull/3614), автор: [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720), автор: [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693), автор: [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720), автор: [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693), автор: [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746), автор: [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 — 2024-07-30

### Исправлено

## v3.0.0-alpha.5 — 2024-07-30

### Добавлено

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   нажатие значка, автор: @5aaee9, в [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) в [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac) в [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) в [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) в [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   на основе [PR](https://github.com/wailsapp/wails/pull/2044), автор: @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) в [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) в [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) в [#3295](https://github.com/wailsapp/wails/pull/3295)   пакет npm, автор: [@fbbdev](https://github.com/fbbdev), в [#3334](https://github.com/wailsapp/wails/pull/3334)   в [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT`, автор: [@abichinger](https://github.com/abichinger), в [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) в [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) в [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) в [#3674](https://github.com/wailsapp/wails/pull/3674)

### Исправлено

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) в   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) в   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) в   [#3477](https://github.com/wailsapp/wails/pull/3477)   от [Atterpac](https://github.com/atterpac) в   [#3320](https://github.com/wailsapp/wails/pull/3320).   в [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) в   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) в   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) в   [PR](https://github.com/wailsapp/wails/pull/3039)   установлен. Добавлено [@pylotlight](https://github.com/pylotlight) в   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   пробелы — @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) в PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) в PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   прикреплённое окно [tw1nk](https://github.com/tw1nk) в PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) в   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` присутствует.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   обработчики событий от [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) в   [#3330](https://github.com/wailsapp/wails/pull/3330)   генератор от [@fbbdev](https://github.com/fbbdev) в   [#3334](https://github.com/wailsapp/wails/pull/3334)   генератор от [@fbbdev](https://github.com/fbbdev) в   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) в   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) в   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) в   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) в   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) в   [#3431](https://github.com/wailsapp/wails/pull/3431)   от [@pekim](https://github.com/pekim) в   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) от   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) в   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows от [@bruxaodev](https://github.com/bruxaodev/) в   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Изменено

[mmghv](https://github.com/mmghv) в   [#3611](https://github.com/wailsapp/wails/pull/3611)   поддержка Ubuntu 24.04 LTS от [atterpac](https://github.com/atterpac) в   [#3461](https://github.com/wailsapp/wails/pull/3461)   должен иметь атрибут `type="module"`. От   [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   объект и не запускает систему WML. Это сделано для улучшения инкапсуляции. При необходимости систему WML можно запустить вручную, вызвав новый метод `WML.Enable`. Встроенный сценарий среды выполнения JS по-прежнему автоматически выполняет обе операции. От [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   объект окна как экспорт по умолчанию. Импортировать отдельные методы с помощью именованного импорта ESM или импорта пространства имён теперь невозможно.   API. У некоторых методов изменились имя или прототип, а именно: `Screen`   заменён на `GetScreen`; `GetZoomLevel`/`SetZoomLevel` заменены на `GetZoom`/`SetZoom`;   `GetZoom`, `Width` и `Height` теперь возвращают значения напрямую, а не оборачивают их в объекты. От [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   удалён. Чтобы вернуться к вызовам по имени, используйте параметр CLI `-names`.   От [@fbbdev](https://github.com/fbbdev) в   [#3468](https://github.com/wailsapp/wails/pull/3468)   именовались по содержащему их пакету; теперь используются полные пути импорта Go, включая путь модуля. От [@fbbdev](https://github.com/fbbdev) в   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. От [@fbbdev](https://github.com/fbbdev) в   [#3468](https://github.com/wailsapp/wails/pull/3468)   вызов `application.NewService`. От [@fbbdev](https://github.com/fbbdev) в   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) в   [#3574](https://github.com/wailsapp/wails/pull/3574)
