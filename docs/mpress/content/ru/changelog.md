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

/_-- Все значимые изменения этого проекта будут задокументированы в этом файле.

Формат основан на [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), а проект следует правилам [семантического версионирования](https://semver.org/spec/v2.0.0.html).

- `Added` — новые возможности.
- `Changed` — изменения существующей функциональности.
- `Deprecated` — возможности, которые вскоре будут удалены.
- `Removed` — удалённые возможности.
- `Fixed` — исправления ошибок.
- `Security` — исправления уязвимостей.

_/

/_   * ОБНОВЛЯЕТСЯ АВТОМАТИЧЕСКИ *   Примечания к выпускам формируются из объединённых PR. Не редактируйте этот файл в обычных PR с изменениями и не создавайте заново удалённый `v3/UNRELEASED_CHANGELOG.md`. Инструкции по ведению журнала изменений см. в `AGENTS.md`. _/

## [Не выпущено]

## v3.2.0-nightly.20261103 — 2026-11-03

## Изменено
- [v3] обновлять меню системного трея при Menu.Update — [PR](https://github.com/wailsapp/wails/pull/6314) от @johnmaguire
- [v3/Windows] реализовать SetIcon для окон — [PR](https://github.com/wailsapp/wails/pull/6502) от @hypherionmc

## v3.2.0-nightly.20261102 — 2026-11-02

## Изменено
- DisableOverscroll вместо DisableBounce — [PR](https://github.com/wailsapp/wails/pull/6077) от @lethusang

## v3.2.0-nightly.20261101 — 2026-11-01

## Добавлено
- добавлены экспериментальные серверные части для macOS и Linux, не требующие CGO, — [PR](https://github.com/wailsapp/wails/pull/6521) от @taliesin-ai
- нативные сборки из YAML, переносимая упаковка и подписание — [PR](https://github.com/wailsapp/wails/pull/6531) от @taliesin-ai

## Исправлено
- исправлены различные проблемы ночных сборок v3 и упорядочены зависимости — [PR](https://github.com/wailsapp/wails/pull/6519) от @taliesin-ai
- исправлены проверка выпусков v3 и классификация фикстур — [PR](https://github.com/wailsapp/wails/pull/6533) от @taliesin-ai
- устранены межплатформенные сбои проверки выпусков — [PR](https://github.com/wailsapp/wails/pull/6538) от @taliesin-ai

## v3.1.0 — 2026-11-01

## Добавлено
- добавлены экспериментальные серверные части для macOS и Linux, не требующие CGO, — [PR](https://github.com/wailsapp/wails/pull/6521) от @taliesin-ai
- нативные сборки из YAML, переносимая упаковка и подписание — [PR](https://github.com/wailsapp/wails/pull/6531) от @taliesin-ai

## Изменено
- исключено копирование отбрасываемых метаданных рабочего стола — [PR](https://github.com/wailsapp/wails/pull/6391) от @taliesin-ai
- длины содержимого форматируются без fmt — [PR](https://github.com/wailsapp/wails/pull/6393) от @taliesin-ai
- корни индекса выбираются без сбора кандидатов — [PR](https://github.com/wailsapp/wails/pull/6404) от @taliesin-ai
- буферы копирования файлов без поддержки поиска используются повторно — [PR](https://github.com/wailsapp/wails/pull/6408) от @taliesin-ai
- исключено выделение памяти под средство чтения для определения типа по префиксам — [PR](https://github.com/wailsapp/wails/pull/6410) от @taliesin-ai
- для известных расширений пропускается чтение с целью определения MIME-типа — [PR](https://github.com/wailsapp/wails/pull/6414) от @taliesin-ai
- предотвращены рост очереди и колебания её сжатия — [PR](https://github.com/wailsapp/wails/pull/6415) от @taliesin-ai
- исключено выделение памяти под отбрасываемые контексты запросов — [PR](https://github.com/wailsapp/wails/pull/6419) от @taliesin-ai
- пропускается подготовка отключённого отладочного журналирования — [PR](https://github.com/wailsapp/wails/pull/6420) от @taliesin-ai
- буферы для определения MIME-типа выделяются отложенно — [PR](https://github.com/wailsapp/wails/pull/6428) от @taliesin-ai
- сокращено число выделений памяти при форматировании гиперссылок — [PR](https://github.com/wailsapp/wails/pull/6443) от @taliesin-ai
- устранены временные выделения памяти при форматировании маркеров списка — [PR](https://github.com/wailsapp/wails/pull/6449) от @taliesin-ai
- преобразование UUID в строку сокращено до одного выделения памяти — [PR](https://github.com/wailsapp/wails/pull/6453) от @taliesin-ai
- сокращено число выделений памяти при форматировании дескрипторов инструкций — [PR](https://github.com/wailsapp/wails/pull/6454) от @taliesin-ai
- устранено временное выделение памяти под хеш UUID на основе имени — [PR](https://github.com/wailsapp/wails/pull/6455) от @taliesin-ai
- исключено копирование незавершённых строк вывода команд — [PR](https://github.com/wailsapp/wails/pull/6460) от @taliesin-ai
- исполняемые файлы копируются потоково при подготовке пакета — [PR](https://github.com/wailsapp/wails/pull/6472) от @taliesin-ai
- исключено повторное сканирование незавершённых строк вывода — [PR](https://github.com/wailsapp/wails/pull/6473) от @taliesin-ai
- скомпилированные шаблоны маршрутизации команд используются повторно — [PR](https://github.com/wailsapp/wails/pull/6475) от @taliesin-ai
- исключены обращения к метаданным при обнаружении среды выполнения — [PR](https://github.com/wailsapp/wails/pull/6482) от @taliesin-ai
- сокращено число выделений памяти для локализованной резервной страницы — [PR](https://github.com/wailsapp/wails/pull/6485) от @taliesin-ai
- исключена повторная компиляция шаблонов для каждого обработчика — [PR](https://github.com/wailsapp/wails/pull/6490) от @taliesin-ai

## Исправлено
- реализован SetTooltip для системного трея — [PR](https://github.com/wailsapp/wails/pull/6017) от @fan711
- настройка тёмного режима считывается из портала настроек — [PR](https://github.com/wailsapp/wails/pull/6072) от @TechHutTV
- Window.SetZoom теперь может устанавливать значение ниже 1.0 — [PR](https://github.com/wailsapp/wails/pull/6166) от @APshenkin
- события доставляются каждому слушателю Go в порядке вызовов emit — [PR](https://github.com/wailsapp/wails/pull/6275) от @DevLumuz
- в Emit обработчики событий считываются под блокировкой hookLock — [PR](https://github.com/wailsapp/wails/pull/6276) от @DevLumuz
- свойство Menu системного трея теперь объявляется всегда, чтобы системный трей, созданный после Run, не отбрасывался, — [PR](https://github.com/wailsapp/wails/pull/6284) от @AlbinoGeek
- работа для главного потока, отправленная до Run, удерживается до запуска Run — [PR](https://github.com/wailsapp/wails/pull/6287) от @AlbinoGeek
- в провайдере обновлений GitHub используется правильный URL ресурса — [PR](https://github.com/wailsapp/wails/pull/6304) от @gavidroselj
- тело запроса сохраняется, когда WebView2 не передаёт Content-Length, — [PR](https://github.com/wailsapp/wails/pull/6326) от @strraaa
- экраны macOS и GTK теперь представлены в координатах DIP (границы при смешанном DPI) — [PR](https://github.com/wailsapp/wails/pull/6346) от @julianstorer
- закрытие уведомления больше не регистрируется как щелчок — [PR](https://github.com/wailsapp/wails/pull/6367) от @bercedev
- процессы разработки останавливаются при закрытии терминала — [PR](https://github.com/wailsapp/wails/pull/6371) от @xmasdev
- при проверке подтверждений запуски формирования журнала изменений сопоставляются с коммитами squash-слияния — [PR](https://github.com/wailsapp/wails/pull/6373) от @taliesin-ai
- выдерживать секундную задержку устранения дребезга при перезапуске панели задач в [PR](https://github.com/wailsapp/wails/pull/6386), автор @taliesin-ai
- сохранять ограничения размера при переходе в полноэкранный режим, начатом с этими ограничениями, в [PR](https://github.com/wailsapp/wails/pull/6387), автор @taliesin-ai
- не изменять размер полноэкранного или развёрнутого окна для соответствия новому ограничению в [PR](https://github.com/wailsapp/wails/pull/6388), автор @taliesin-ai
- не открывать диалоговое окно до Run в примере liquid-glass в [PR](https://github.com/wailsapp/wails/pull/6389), автор @taliesin-ai
- диалоговое окно, запрошенное до Run, теперь ничего не делает вместо бесконечного ожидания в [PR](https://github.com/wailsapp/wails/pull/6390), автор @taliesin-ai
- использовать резервный вариант, если в записи рабочего стола нет команды, в [PR](https://github.com/wailsapp/wails/pull/6392), автор @taliesin-ai
- определять тип текущего содержимого, не сохраняя имена файлов, в [PR](https://github.com/wailsapp/wails/pull/6394), автор @taliesin-ai
- распознавать расширения ресурсов с буквами в разных регистрах в [PR](https://github.com/wailsapp/wails/pull/6395), автор @taliesin-ai
- освобождать использованные буферы анализатора содержимого в [PR](https://github.com/wailsapp/wails/pull/6396), автор @taliesin-ai
- освобождать удалённые ссылки ringqueue в [PR](https://github.com/wailsapp/wails/pull/6397), автор @taliesin-ai
- учитывать удаление заголовков до выбора резервного варианта в [PR](https://github.com/wailsapp/wails/pull/6398), автор @taliesin-ai
- фиксировать заголовки резервного ответа перед сбросом буфера в [PR](https://github.com/wailsapp/wails/pull/6399), автор @taliesin-ai
- обнаруживать только файлы с точными именами индексных файлов в [PR](https://github.com/wailsapp/wails/pull/6400), автор @taliesin-ai
- не формировать тело и не выполнять полное чтение для запроса HEAD к потоку без поддержки позиционирования в [PR](https://github.com/wailsapp/wails/pull/6401), автор @taliesin-ai
- заполнять префикс для определения MIME-типа за несколько коротких чтений в [PR](https://github.com/wailsapp/wails/pull/6402), автор @taliesin-ai
- фиксировать успешные пустые ответы без поддержки позиционирования в [PR](https://github.com/wailsapp/wails/pull/6403), автор @taliesin-ai
- не дублировать цели URI после подстановки полей в [PR](https://github.com/wailsapp/wails/pull/6405), автор @taliesin-ai
- сохранять пути при подстановке полей в [PR](https://github.com/wailsapp/wails/pull/6406), автор @taliesin-ai
- учитывать каталоги поиска desktop-файлов XDG в [PR](https://github.com/wailsapp/wails/pull/6407), автор @taliesin-ai
- фиксировать заголовки анализатора содержимого при пустом сбросе буфера в [PR](https://github.com/wailsapp/wails/pull/6409), автор @taliesin-ai
- пропускать каталоги, затеняющие записи рабочего стола, в [PR](https://github.com/wailsapp/wails/pull/6411), автор @taliesin-ai
- сохранять цели, найденные внутри путей к исполняемым файлам, в [PR](https://github.com/wailsapp/wails/pull/6412), автор @taliesin-ai
- освобождать захваченные обратным вызовом данные после срабатывания таймера в [PR](https://github.com/wailsapp/wails/pull/6413), автор @taliesin-ai
- исключить ложные совпадения библиотек ldconfig в [PR](https://github.com/wailsapp/wails/pull/6416), автор @taliesin-ai
- ограничить совпадения при резервном поиске библиотек в файловой системе в [PR](https://github.com/wailsapp/wails/pull/6417), автор @taliesin-ai
- устранить гонку данных при инвалидации кеша в [PR](https://github.com/wailsapp/wails/pull/6418), автор @taliesin-ai
- учитывать необязательную конфигурацию журналирования в [PR](https://github.com/wailsapp/wails/pull/6421), автор @taliesin-ai
- исключить изменение общего начального URL и использование устаревших портов в [PR](https://github.com/wailsapp/wails/pull/6422), автор @taliesin-ai
- прекращать повторные попытки прокси-сервера разработки при отмене запроса в [PR](https://github.com/wailsapp/wails/pull/6423), автор @taliesin-ai
- возвращать 404 для отсутствующих встроенных ресурсов в [PR](https://github.com/wailsapp/wails/pull/6424), автор @taliesin-ai
- обрабатывать запросы HEAD к встроенной среде выполнения в [PR](https://github.com/wailsapp/wails/pull/6425), автор @taliesin-ai
- сохранять экранированные пути при маршрутизации сервисов в [PR](https://github.com/wailsapp/wails/pull/6426), автор @taliesin-ai
- правильно упорядочивать числовые идентификаторы предварительных выпусков, превышающие допустимый размер, в [PR](https://github.com/wailsapp/wails/pull/6427), автор @taliesin-ai
- отклонять ответы с ошибками HTTP до сохранения загружаемых инструментов в [PR](https://github.com/wailsapp/wails/pull/6429), автор @taliesin-ai
- соблюдать крайний срок запуска при поиске записей рабочего стола в [PR](https://github.com/wailsapp/wails/pull/6430), автор @taliesin-ai
- сохранять тела запросов между повторными попытками прокси-сервера разработки в [PR](https://github.com/wailsapp/wails/pull/6432), автор @taliesin-ai
- закрывать тела HTTP-ответов во всех ветвях возврата в [PR](https://github.com/wailsapp/wails/pull/6434), автор @taliesin-ai
- сохранять ошибки отмены при параллельном поиске библиотек в [PR](https://github.com/wailsapp/wails/pull/6435), автор @taliesin-ai
- отклонять неуспешные ответы GitHub с тегами до декодирования в [PR](https://github.com/wailsapp/wails/pull/6436), автор @taliesin-ai
- сохранять буквальные пути шаблонов с квадратными скобками в [PR](https://github.com/wailsapp/wails/pull/6437), автор @taliesin-ai
- устанавливать wails3 один раз и сохранять ошибки установщика в [PR](https://github.com/wailsapp/wails/pull/6438), автор @taliesin-ai
- игнорировать скрытые элементы и заголовки в блоках комментариев в [PR](https://github.com/wailsapp/wails/pull/6439), автор @taliesin-ai
- отображать ошибки поиска стабильного выпуска в [PR](https://github.com/wailsapp/wails/pull/6440), автор @taliesin-ai
- выводить обычные сообщения без интерпретации форматирования в [PR](https://github.com/wailsapp/wails/pull/6441), автор @taliesin-ai
- предотвращать панику при записи после загрузки JSON null в [PR](https://github.com/wailsapp/wails/pull/6442), автор @taliesin-ai
- возвращать снимки состояния при чтении всего хранилища в [PR](https://github.com/wailsapp/wails/pull/6444), автор @taliesin-ai
- передавать ошибки после частичного перебора строк в [PR](https://github.com/wailsapp/wails/pull/6446), автор @taliesin-ai
- сохранять пути назначения распаковки буквально в [PR](https://github.com/wailsapp/wails/pull/6447), автор @taliesin-ai
- освобождать дескрипторы подготовленных выражений при закрытии и повторном открытии в [PR](https://github.com/wailsapp/wails/pull/6448), автор @taliesin-ai
- обрабатывать гонку подготовки с закрытием или повторным открытием в [PR](https://github.com/wailsapp/wails/pull/6450), автор @taliesin-ai
- правильно обрабатывать запросы HEAD к резервной индексной странице в [PR](https://github.com/wailsapp/wails/pull/6451), автор @taliesin-ai
- требовать полные шестнадцатеричные дескрипторы выражений в [PR](https://github.com/wailsapp/wails/pull/6452), автор @taliesin-ai
- сохранять стандартную базу данных в памяти при конкурентном доступе в [PR](https://github.com/wailsapp/wails/pull/6456), автор @taliesin-ai
- отделять stderr от анализируемого вывода HEAD в [PR](https://github.com/wailsapp/wails/pull/6457), автор @taliesin-ai
- создавать ресурсы среды выполнения из встроенных пакетов в [PR](https://github.com/wailsapp/wails/pull/6458), автор @taliesin-ai
- обнаруживать заполнители TODO, FIXME и TBD в [PR](https://github.com/wailsapp/wails/pull/6459), автор @taliesin-ai
- проверять ответы API с примечаниями к выпуску в [PR](https://github.com/wailsapp/wails/pull/6461), автор @taliesin-ai
- учитывать каталоги задач при проверке состояния в [PR](https://github.com/wailsapp/wails/pull/6462), автор @taliesin-ai
- предотвращать потерю данных при копировании файла в самого себя в [PR](https://github.com/wailsapp/wails/pull/6463), автор @taliesin-ai
- включать каталог задачи в идентификатор кеша в [PR](https://github.com/wailsapp/wails/pull/6464), автор @taliesin-ai
- проверять порты один раз за проверку в [PR](https://github.com/wailsapp/wails/pull/6465), автор @taliesin-ai
- предотвращать панику при декодировании записей кеша как null в [PR](https://github.com/wailsapp/wails/pull/6466), автор @taliesin-ai
- увеличивать числовые идентификаторы предварительных выпусков без ограничений целочисленного типа в [PR](https://github.com/wailsapp/wails/pull/6467), автор @taliesin-ai
- не допускать, чтобы устаревшие проверки скрывали промахи кеша, в [PR](https://github.com/wailsapp/wails/pull/6468) от @taliesin-ai
- выполнять задачи по умолчанию в базовом каталоге исполнителя в [PR](https://github.com/wailsapp/wails/pull/6469) от @taliesin-ai
- сохранять литеральные знаки процента в командах проводника в [PR](https://github.com/wailsapp/wails/pull/6470) от @taliesin-ai
- согласовать каталоги предварительных условий и команд задач в [PR](https://github.com/wailsapp/wails/pull/6471) от @taliesin-ai
- заменить символические ссылки сгенерированными обычными файлами в [PR](https://github.com/wailsapp/wails/pull/6474) от @taliesin-ai
- сохранять существующие DMG при сбое сборки в [PR](https://github.com/wailsapp/wails/pull/6476) от @taliesin-ai
- сохранять обратные косые черты в заключённых в кавычки регулярных выражениях и путях в [PR](https://github.com/wailsapp/wails/pull/6477) от @taliesin-ai
- проверять размеры генерируемых значков перед выделением памяти под изображение в [PR](https://github.com/wailsapp/wails/pull/6478) от @taliesin-ai
- передавать цепочки команд и перенаправления оболочке в [PR](https://github.com/wailsapp/wails/pull/6479) от @taliesin-ai
- сохранять предыдущие пакеты Linux при сбое сборки в [PR](https://github.com/wailsapp/wails/pull/6480) от @taliesin-ai
- использовать резервное выполнение через оболочку с учётом платформы в [PR](https://github.com/wailsapp/wails/pull/6481) от @taliesin-ai
- обслуживать ресурсы с повторяющимися начальными косыми чертами в пути в [PR](https://github.com/wailsapp/wails/pull/6483) от @taliesin-ai
- распознавать стабильные версии с метаданными сборки в [PR](https://github.com/wailsapp/wails/pull/6484) от @taliesin-ai
- точно сообщать о сбоях при генерации ресурсов в [PR](https://github.com/wailsapp/wails/pull/6486) от @taliesin-ai
- предотвращать состояния гонки при замене активного средства отчётности в [PR](https://github.com/wailsapp/wails/pull/6487) от @taliesin-ai
- сериализовать одновременное сохранение файлов кеша в [PR](https://github.com/wailsapp/wails/pull/6488) от @taliesin-ai
- поддерживать несколько кандидатов на равенство в шаблоне в [PR](https://github.com/wailsapp/wails/pull/6489) от @taliesin-ai
- нормализовать сохранённые карты категорий со значением null в [PR](https://github.com/wailsapp/wails/pull/6491) от @taliesin-ai
- показывать диагностику неудачно завершившихся переменных команд оболочки в [PR](https://github.com/wailsapp/wails/pull/6492) от @taliesin-ai
- не помещать цепочки команд оболочки в неявный кеш Go в [PR](https://github.com/wailsapp/wails/pull/6493) от @taliesin-ai
- сериализовать сохранение категорий в Linux в [PR](https://github.com/wailsapp/wails/pull/6494) от @taliesin-ai
- кешировать выходные файлы Go с заключёнными в кавычки или экранированными пробелами в [PR](https://github.com/wailsapp/wails/pull/6495) от @taliesin-ai
- сохранять значки Mac при сбое кодирования в [PR](https://github.com/wailsapp/wails/pull/6496) от @taliesin-ai
- не отбрасывать недопустимые флаги количества запусков тестов Go в [PR](https://github.com/wailsapp/wails/pull/6497) от @taliesin-ai
- сохранять права доступа исходного файла при создании его копии в [PR](https://github.com/wailsapp/wails/pull/6498) от @taliesin-ai
- Исправить различные проблемы ночной сборки v3 и очистить зависимости в [PR](https://github.com/wailsapp/wails/pull/6519) от @taliesin-ai
- исправить проверку выпуска v3 и классификацию фикстур в [PR](https://github.com/wailsapp/wails/pull/6533) от @taliesin-ai
- устранить межплатформенные сбои проверки выпуска в [PR](https://github.com/wailsapp/wails/pull/6538) от @taliesin-ai

## v3.1.0-nightly.20261028 — 2026-10-28

## Изменено
- Ускорить создание журнала изменений за счёт повторного использования скомпилированных шаблонов регулярных выражений в [PR](https://github.com/wailsapp/wails/pull/6490) от @taliesin-ai
- Сократить количество выделений памяти при обслуживании локализованных резервных страниц в [PR](https://github.com/wailsapp/wails/pull/6485) от @taliesin-ai
- Ускорить обнаружение среды выполнения AppImage, исключив ненужные обращения к метаданным, в [PR](https://github.com/wailsapp/wails/pull/6482) от @taliesin-ai
- Повысить производительность wake за счёт повторного использования скомпилированных шаблонов маршрутизации команд в [PR](https://github.com/wailsapp/wails/pull/6475) от @taliesin-ai
- Повысить производительность, исключив повторное сканирование незавершённых строк вывода, в [PR](https://github.com/wailsapp/wails/pull/6473) от @taliesin-ai
- Копировать исполняемые файлы потоково во время подготовки пакета MSIX, чтобы сократить потребление памяти и ускорить сборку, в [PR](https://github.com/wailsapp/wails/pull/6472) от @taliesin-ai
- Не копировать незавершённые строки вывода команд для повышения производительности в [PR](https://github.com/wailsapp/wails/pull/6460) от @taliesin-ai
- Сократить количество выделений памяти при генерации UUID на основе имени в [PR](https://github.com/wailsapp/wails/pull/6455) от @taliesin-ai
- Сократить количество выделений памяти при форматировании дескрипторов операторов SQLite в [PR](https://github.com/wailsapp/wails/pull/6454) от @taliesin-ai
- Сократить преобразование UUID в строку до одного выделения памяти в [PR](https://github.com/wailsapp/wails/pull/6453) от @taliesin-ai
- Сократить количество выделений памяти при форматировании пунктов журнала изменений в [PR](https://github.com/wailsapp/wails/pull/6449) от @taliesin-ai
- Сократить количество выделений памяти при форматировании гиперссылок терминала в [PR](https://github.com/wailsapp/wails/pull/6443) от @taliesin-ai
- Сократить количество выделений памяти при обслуживании ресурсов за счёт отложенного анализа MIME-буферов в [PR](https://github.com/wailsapp/wails/pull/6428) от @taliesin-ai
- Не подготавливать отладочные сообщения в сервере ресурсов, когда отладочное журналирование отключено, в [PR](https://github.com/wailsapp/wails/pull/6420) от @taliesin-ai
- Сократить количество выделений памяти при обработке запросов webview, исключив создание отбрасываемых контекстов запросов, в [PR](https://github.com/wailsapp/wails/pull/6419) от @taliesin-ai
- Ускорить обслуживание ресурсов, не определяя MIME для известных расширений файлов, в [PR](https://github.com/wailsapp/wails/pull/6414) от @taliesin-ai
- Сократить количество выделений памяти при обслуживании ресурсов, исключив использование средства чтения для определения типа содержимого, в [PR](https://github.com/wailsapp/wails/pull/6410) от @taliesin-ai
- Сократить количество выделений памяти при обслуживании файлов ресурсов без поддержки позиционирования в [PR](https://github.com/wailsapp/wails/pull/6408) от @taliesin-ai
- Ускорить выбор корня индекса сервером ресурсов, исключив сбор кандидатов, в [PR](https://github.com/wailsapp/wails/pull/6404) от @taliesin-ai
- Повысить производительность сервера ресурсов, форматируя длину содержимого без fmt, в [PR](https://github.com/wailsapp/wails/pull/6393) от @taliesin-ai
- Сократить копирование данных в памяти при чтении метаданных проводника файлов рабочего стола в [PR](https://github.com/wailsapp/wails/pull/6391) от @taliesin-ai

## Исправлено
- Сохранять права доступа исходного файла при копировании в новый файл назначения в [PR](https://github.com/wailsapp/wails/pull/6498) от @taliesin-ai
- Не изменять размер полноэкранных и развёрнутых окон при изменении ограничений минимального или максимального размера в [PR](https://github.com/wailsapp/wails/pull/6388) от @taliesin-ai
- Публиковать запланированные обновления истории количества звёзд для репозитория wailsapp/wails в [PR](https://github.com/wailsapp/wails/pull/6500) от @taliesin-ai
- Не отбрасывать недопустимые флаги количества запусков тестов Go в wake в [PR](https://github.com/wailsapp/wails/pull/6497) от @taliesin-ai
- Сохранять значки Mac при сбое кодирования в [PR](https://github.com/wailsapp/wails/pull/6496) от @taliesin-ai
- Исправить промахи кеша сборки для выходных файлов Go с заключёнными в кавычки или экранированными пробелами в [PR](https://github.com/wailsapp/wails/pull/6495) от @taliesin-ai
- Сериализовать сохранение категорий уведомлений в Linux, чтобы предотвратить состояния гонки при одновременной записи, в [PR](https://github.com/wailsapp/wails/pull/6494) от @taliesin-ai
- Не помещать цепочки команд оболочки в неявный кеш сборки Go в [PR](https://github.com/wailsapp/wails/pull/6493) от @taliesin-ai
- Исправлена ошибка загрузки карт категорий уведомлений, если сохранённые данные имеют значение null, в [PR](https://github.com/wailsapp/wails/pull/6491) от @taliesin-ai
- Добавлен вывод диагностических сообщений при сбое команд оболочки с переменными в [PR](https://github.com/wailsapp/wails/pull/6492) от @taliesin-ai
- Исправлено сопоставление шаблонов wake при наличии нескольких кандидатов на равенство в [PR](https://github.com/wailsapp/wails/pull/6489) от @taliesin-ai
- Добавлена сериализация параллельных сохранений файлов кэша для предотвращения повреждения записываемого кэша в [PR](https://github.com/wailsapp/wails/pull/6488) от @taliesin-ai
- Предотвращены состояния гонки при замене активного средства формирования отчётов в [PR](https://github.com/wailsapp/wails/pull/6487) от @taliesin-ai
- CLI теперь корректно сообщает об ошибках генерации ресурсов в [PR](https://github.com/wailsapp/wails/pull/6486) от @taliesin-ai
- CLI теперь распознаёт стабильные версии, содержащие метаданные сборки, в [PR](https://github.com/wailsapp/wails/pull/6484) от @taliesin-ai
- Добавлена раздача ресурсов, пути запросов к которым содержат повторяющиеся начальные косые черты, в [PR](https://github.com/wailsapp/wails/pull/6483) от @taliesin-ai
- Исправлено выполнение резервной команды оболочки для wake на каждой платформе в [PR](https://github.com/wailsapp/wails/pull/6481) от @taliesin-ai
- Обеспечено сохранение предыдущих пакетов Linux при сбое сборки в [PR](https://github.com/wailsapp/wails/pull/6480) от @taliesin-ai
- Цепочки команд и перенаправления теперь выполняются через оболочку в [PR](https://github.com/wailsapp/wails/pull/6479) от @taliesin-ai
- Обеспечено сохранение заключённых в кавычки регулярных выражений и обратных косых черт в путях при разборе команд wake в [PR](https://github.com/wailsapp/wails/pull/6477) от @taliesin-ai
- Обеспечено сохранение существующих файлов DMG при сбое упаковки сборок в [PR](https://github.com/wailsapp/wails/pull/6476) от @taliesin-ai
- Символические ссылки в выходных данных привязок теперь заменяются сгенерированными обычными файлами в [PR](https://github.com/wailsapp/wails/pull/6474) от @taliesin-ai
- Обеспечено сохранение литеральных знаков процента в командах файлового менеджера Linux в [PR](https://github.com/wailsapp/wails/pull/6470) от @taliesin-ai
- Каталоги команд предварительных условий и задач согласованы для единообразного разрешения в [PR](https://github.com/wailsapp/wails/pull/6471) от @taliesin-ai
- Задачи по умолчанию теперь выполняются в базовом каталоге исполнителя в [PR](https://github.com/wailsapp/wails/pull/6469) от @taliesin-ai
- Устаревшие проверки больше не скрывают промахи кэша в [PR](https://github.com/wailsapp/wails/pull/6468) от @taliesin-ai
- Исправлена нумерация предварительных выпусков в CLI: теперь она увеличивается без ограничений, связанных с переполнением целых чисел, в [PR](https://github.com/wailsapp/wails/pull/6467) от @taliesin-ai
- Предотвращена паника при декодировании записей кэша со значением null в [PR](https://github.com/wailsapp/wails/pull/6466) от @taliesin-ai
- CLI теперь проверяет порты один раз за каждую проверку в [PR](https://github.com/wailsapp/wails/pull/6465) от @taliesin-ai
- Каталог задачи включён в идентификатор кэша, чтобы избежать устаревших результатов wake, в [PR](https://github.com/wailsapp/wails/pull/6464) от @taliesin-ai
- Предотвращена потеря данных при копировании файла поверх самого себя во время упаковки в [PR](https://github.com/wailsapp/wails/pull/6463) от @taliesin-ai
- Проверки состояния теперь учитывают каталоги задач в [PR](https://github.com/wailsapp/wails/pull/6462) от @taliesin-ai
- Добавлено обнаружение заполнителей TODO, FIXME и TBD в записях журнала изменений в [PR](https://github.com/wailsapp/wails/pull/6459) от @taliesin-ai
- CLI теперь генерирует ресурсы среды выполнения из встроенных пакетов в [PR](https://github.com/wailsapp/wails/pull/6458) от @taliesin-ai
- Поток stderr отделён от разбираемого вывода HEAD в командах git в [PR](https://github.com/wailsapp/wails/pull/6457) от @taliesin-ai
- Обеспечено сохранение стандартной базы данных SQLite в памяти при параллельном доступе в [PR](https://github.com/wailsapp/wails/pull/6456) от @taliesin-ai
- Для запросов HEAD к резервной индексной странице теперь возвращаются правильные ответы в [PR](https://github.com/wailsapp/wails/pull/6451) от @taliesin-ai
- Обработана гонка между подготовкой операторов SQLite и закрытием или повторным открытием базы данных в [PR](https://github.com/wailsapp/wails/pull/6450) от @taliesin-ai
- Дескрипторы подготовленных операторов теперь освобождаются при закрытии или повторном открытии базы данных SQLite в [PR](https://github.com/wailsapp/wails/pull/6448) от @taliesin-ai
- При использовании gosod пути назначения для извлечения сохраняются в литеральном виде в [PR](https://github.com/wailsapp/wails/pull/6447) от @taliesin-ai
- Ошибки запросов SQLite теперь передаются после частичной итерации по строкам в [PR](https://github.com/wailsapp/wails/pull/6446) от @taliesin-ai
- При чтении всего хранилища в kvstore теперь возвращаются снимки данных в [PR](https://github.com/wailsapp/wails/pull/6444) от @taliesin-ai
- Предотвращены паники при записи в kvstore после загрузки значения JSON null в [PR](https://github.com/wailsapp/wails/pull/6442) от @taliesin-ai
- Обычные сообщения терминала теперь выводятся без интерпретации спецификаторов форматирования в [PR](https://github.com/wailsapp/wails/pull/6441) от @taliesin-ai
- Средство обновления теперь сообщает об ошибках поиска последнего стабильного выпуска в [PR](https://github.com/wailsapp/wails/pull/6440) от @taliesin-ai
- Скрытые записи и заголовки в блоках комментариев журнала изменений теперь игнорируются в [PR](https://github.com/wailsapp/wails/pull/6439) от @taliesin-ai
- Средство обновления теперь устанавливает wails3 один раз и сообщает об ошибках установщика в [PR](https://github.com/wailsapp/wails/pull/6438) от @taliesin-ai
- В gosod сохраняются литеральные пути шаблонов с квадратными скобками в [PR](https://github.com/wailsapp/wails/pull/6437) от @taliesin-ai
- Средство обновления теперь отклоняет ошибочные ответы с тегами GitHub до их декодирования в [PR](https://github.com/wailsapp/wails/pull/6436) от @taliesin-ai
- Обеспечено сохранение ошибок отмены при параллельном поиске библиотек в [PR](https://github.com/wailsapp/wails/pull/6435) от @taliesin-ai
- Тела HTTP-ответов теперь закрываются во всех ветвях возврата для предотвращения утечек ресурсов в [PR](https://github.com/wailsapp/wails/pull/6434) от @taliesin-ai
- Обеспечено сохранение тел запросов при повторных попытках прокси-сервера разработки в [PR](https://github.com/wailsapp/wails/pull/6432) от @taliesin-ai
- При поиске файлового менеджера рабочего стола теперь соблюдается крайний срок запуска в [PR](https://github.com/wailsapp/wails/pull/6430) от @taliesin-ai
- HTTP-ошибки теперь отклоняются до сохранения загрузок инструментов AppImage в [PR](https://github.com/wailsapp/wails/pull/6429) от @taliesin-ai
- Слишком большие числовые идентификаторы предварительных выпусков теперь правильно упорядочиваются при сравнении semver в [PR](https://github.com/wailsapp/wails/pull/6427) от @taliesin-ai
- При маршрутизации сервисов на сервере ресурсов сохраняются экранированные пути в [PR](https://github.com/wailsapp/wails/pull/6426) от @taliesin-ai
- На сервере ресурсов добавлена обработка запросов HEAD к встроенной среде выполнения в [PR](https://github.com/wailsapp/wails/pull/6425) от @taliesin-ai
- Для отсутствующих встроенных ресурсов теперь возвращается 404 вместо некорректного ответа в [PR](https://github.com/wailsapp/wails/pull/6424) от @taliesin-ai
- Повторные попытки прокси-сервера разработки теперь прекращаются при отмене запроса в [PR](https://github.com/wailsapp/wails/pull/6423) от @taliesin-ai
- Устранены совместное изменение начального URL и устаревшие порты на сервере ресурсов в [PR](https://github.com/wailsapp/wails/pull/6422) от @taliesin-ai
- Сервер ресурсов теперь учитывает необязательную конфигурацию регистратора в [PR](https://github.com/wailsapp/wails/pull/6421) от @taliesin-ai
- Устранена гонка данных при аннулировании кэша путей к библиотекам в [PR](https://github.com/wailsapp/wails/pull/6418) от @taliesin-ai
- Устранены ложные совпадения библиотек при поиске в выводе ldconfig в [PR](https://github.com/wailsapp/wails/pull/6416) от @taliesin-ai
- Предотвращены рост очереди сервера ресурсов и колебания её размера под нагрузкой в [PR](https://github.com/wailsapp/wails/pull/6415) от @taliesin-ai
- Захваченные обратным вызовом значения теперь освобождаются после срабатывания таймера устранения дребезга, чтобы не удерживать устаревшие ссылки, в [PR](https://github.com/wailsapp/wails/pull/6413) от @taliesin-ai
- Обеспечено сохранение целей файлового менеджера, найденных внутри путей к исполняемым файлам, в [PR](https://github.com/wailsapp/wails/pull/6412) от @taliesin-ai
- Пропускать каталоги, которые перекрывают desktop-файлы в файловом менеджере, — [PR](https://github.com/wailsapp/wails/pull/6411) от @taliesin-ai
- Исправить отсутствие заголовков сервера ресурсов при пустом теле ответа — [PR](https://github.com/wailsapp/wails/pull/6409) от @taliesin-ai
- Учитывать каталоги поиска desktop-файлов XDG при поиске значков приложений — [PR](https://github.com/wailsapp/wails/pull/6407) от @taliesin-ai
- Сохранять пути к файлам при подстановке полей в файловом менеджере — [PR](https://github.com/wailsapp/wails/pull/6406) от @taliesin-ai
- Не дублировать цели URI после подстановки в файловом менеджере — [PR](https://github.com/wailsapp/wails/pull/6405) от @taliesin-ai
- Завершать успешные пустые ответы ресурсов без поддержки перемещения по потоку, а не оставлять их в ожидании, — [PR](https://github.com/wailsapp/wails/pull/6403) от @taliesin-ai
- Правильно заполнять префикс для определения MIME-типа при нескольких коротких чтениях на сервере ресурсов — [PR](https://github.com/wailsapp/wails/pull/6402) от @taliesin-ai
- Не считывать тела ответов целиком для HEAD-запросов без поддержки перемещения по потоку на сервере ресурсов — [PR](https://github.com/wailsapp/wails/pull/6401) от @taliesin-ai
- При обработке запросов к каталогам на сервере ресурсов возвращать только индексные файлы с точным совпадением имени — [PR](https://github.com/wailsapp/wails/pull/6400) от @taliesin-ai
- Исправить ошибку, из-за которой отклонённые уведомления Linux регистрировались как нажатия на серверах без поддержки действий, — [PR](https://github.com/wailsapp/wails/pull/6367) от @bercedev
- Останавливать процессы разработки при закрытии терминала — [PR](https://github.com/wailsapp/wails/pull/6371) от @xmasdev
- После перезапуска панели задач Windows выдерживать задержку устранения дребезга в одну секунду перед сбросом флага — [PR](https://github.com/wailsapp/wails/pull/6386) от @taliesin-ai
- Исправить потерю заголовков резервного ответа на сервере ресурсов — [PR](https://github.com/wailsapp/wails/pull/6399) от @taliesin-ai
- Учитывать удаление заголовков перед выбором резервного ответа на сервере ресурсов — [PR](https://github.com/wailsapp/wails/pull/6398) от @taliesin-ai
- Освобождать удалённые ссылки кольцевой очереди на сервере ресурсов, чтобы избежать утечек памяти, — [PR](https://github.com/wailsapp/wails/pull/6397) от @taliesin-ai
- Освобождать использованные буферы определения типа содержимого на сервере ресурсов, чтобы снизить потребление памяти, — [PR](https://github.com/wailsapp/wails/pull/6396) от @taliesin-ai
- Распознавать расширения файлов ресурсов со смешанным регистром при их выдаче — [PR](https://github.com/wailsapp/wails/pull/6395) от @taliesin-ai
- Правильно определять тип содержимого на сервере ресурсов, не сохраняя имена файлов, — [PR](https://github.com/wailsapp/wails/pull/6394) от @taliesin-ai
- Использовать значение по умолчанию, если в desktop-файле отсутствует команда, — [PR](https://github.com/wailsapp/wails/pull/6392) от @taliesin-ai
- Не показывать в примере liquid-glass диалоговое окно до запуска приложения на платформах, отличных от macOS, — [PR](https://github.com/wailsapp/wails/pull/6389) от @taliesin-ai
- Сохранять ограничения размера окна в Windows, если полноэкранный режим включается при размере, заданном с ограничениями, — [PR](https://github.com/wailsapp/wails/pull/6387) от @taliesin-ai
- Исправить бесконечную блокировку диалоговых окон, вызванных до Run, вместо возврата из вызова — [PR](https://github.com/wailsapp/wails/pull/6390) от @taliesin-ai

## Безопасность
- Проверять размеры создаваемых значков до выделения памяти под изображение, чтобы предотвратить чрезмерное потребление памяти, — [PR](https://github.com/wailsapp/wails/pull/6478) от @taliesin-ai
- Проверять ответы API примечаний к выпуску перед использованием — [PR](https://github.com/wailsapp/wails/pull/6461) от @taliesin-ai
- Требовать полные шестнадцатеричные дескрипторы выражений в плагине SQLite — [PR](https://github.com/wailsapp/wails/pull/6452) от @taliesin-ai
- Ограничить резервный поиск библиотек в файловой системе, чтобы предотвратить загрузку непредусмотренных библиотек, — [PR](https://github.com/wailsapp/wails/pull/6417) от @taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## Изменено
- Связать страницу 404 с соответствующими страницами документации Wails v2 на v2.wails.io — [PR](https://github.com/wailsapp/wails/pull/6366) от @taliesin-ai
- Перенести сайт документации v2 на v2.wails.io — [PR](https://github.com/wailsapp/wails/pull/6365) от @taliesin-ai
- Публиковать документацию v3 на wails.io, а документацию v2 — на v2.wails.io; [PR](https://github.com/wailsapp/wails/pull/6364) от @taliesin-ai
- Обновить примечания к выпуску v3.0.0 и исправить утверждение в бюллетене безопасности — [PR](https://github.com/wailsapp/wails/pull/6363) от @taliesin-ai

## Исправлено
- Последовательно передавать обратные вызовы обработчиков событий Go в порядке их отправки — [PR](https://github.com/wailsapp/wails/pull/6275) от @DevLumuz
- Всегда объявлять свойство Menu для значка в области уведомлений, чтобы в Linux не отбрасывались значки, созданные после Run, — [PR](https://github.com/wailsapp/wails/pull/6284) от @AlbinoGeek
- Удерживать работу главного потока, отправленную до Run, пока приложение не запустится, — [PR](https://github.com/wailsapp/wails/pull/6287) от @AlbinoGeek
- Использовать правильный URL API ресурсов при скачивании ресурсов выпуска для средства обновления GitHub — [PR](https://github.com/wailsapp/wails/pull/6304) от @gavidroselj
- Исправить границы и координаты экранов в конфигурациях macOS и Linux со смешанными значениями DPI — [PR](https://github.com/wailsapp/wails/pull/6346) от @julianstorer
- Сохранять тела запросов в Windows, когда WebView2 не передаёт заголовок Content-Length, — [PR](https://github.com/wailsapp/wails/pull/6326) от @strraaa
- Устранить гонку данных при отправке перехватчикам событий, создавая снимок перехватчиков под блокировкой, — [PR](https://github.com/wailsapp/wails/pull/6276) от @DevLumuz
- Исправить бездействие `SystemTray.SetTooltip()` в Linux: теперь задаётся описание всплывающей подсказки StatusNotifierItem (#6017)
- Запускать проверку подтверждающих данных журнала изменений в еженедельном процессе выпуска — [PR](https://github.com/wailsapp/wails/pull/6374) от @taliesin-ai
- Исправить сопоставление запусков с коммитами squash-слияния при проверке подтверждающих данных журнала изменений выпуска — [PR](https://github.com/wailsapp/wails/pull/6373) от @taliesin-ai
- Исправить определение тёмного режима в Linux, которое возвращало светлый режим на тёмном рабочем столе до изменения цветовой схемы, а также проблему, из-за которой `linux:SystemThemeChanged` никогда не срабатывал в бэкенде GTK4. Теперь оба бэкенда используют одну реализацию на основе пространства имён портала `org.freedesktop.appearance` (#6072)
- На окнах с рамкой (Windows, macOS, Linux) `Window.SetZoom` больше не повышает коэффициенты ниже 1.0 до 1.0: `ZoomFactor` WebView2, `magnification` WKWebView и уровень масштабирования WebKitGTK принимают любое положительное значение, а параметр окна `Zoom` уже передавал такие значения без изменений. Для окон без рамки сохраняется нижний предел 1.0, поскольку их области перетаскивания и изменения размера не корректируются с учётом масштаба ([#4590](https://github.com/wailsapp/wails/issues/4590)); теперь неположительные коэффициенты, NaN и Inf игнорируются, а не повышаются до 1.0 — от @APshenkin

## v3.0.0 - 2026-10-26

## Изменено
- Предварительные ночные выпуски публикуются как `vX.Y.0-nightly.YYYYMMDD` и позволяют ознакомиться со следующим еженедельным выпуском, а еженедельные выпуски увеличивают младшую версию (`vX.Y.0`), поэтому команда Go `@latest` всегда разрешается в еженедельный выпуск или выпуск с экстренным исправлением; версия npm-пакета `@wailsio/runtime` всегда совпадает с версией выпуска Wails, а ночные выпуски публикуются с dist-tag `nightly`
- Шаблоны проектов закрепляют `@wailsio/runtime` на версии CLI, с помощью которой они были созданы
- Для Wails v3 требуется Go 1.26 или новее (Go 1.26.10 либо 1.27.3 и более поздние версии содержат исправление безопасности `crypto/x509`)

## Исправлено
- Исправить ошибку, из-за которой `wails3 update cli` не находил выпуски v3: он учитывал только теги v2 и устанавливал их из неверного пути
- Исправлен сбой приложений Windows с `panic: SetTextColor failed` при отрисовке строки меню приложения в некоторых сборках Windows: теперь строка меню отрисовывается в допустимом контексте устройства либо используется стандартная отрисовка, `UAHMENUITEM.Position` соответствует C-структуре `int`, а отрисовка меню больше не вызывает панику (#6311)
- Обеспечена сохранность подготовленного обновления в ожидании перезапуска: загрузка или повторный перезапуск, начатые во время завершения работы, не могут удалить полезную нагрузку, которую ожидает установить вспомогательный процесс замены, или вступить за неё в конфликт, — @julianstorer
- Исправлен сбой приложений Windows в `applyDarkMode`, когда инициализация `pkg/w32` выгружала `uxtheme.dll`, пока указатели на её функции по порядковым номерам ещё использовались: теперь модуль остаётся загруженным на протяжении всего времени жизни процесса (#6356)
- Добавлена проверка подготовленной полезной нагрузки перед заменой и установленного приложения перед удалением резервной копии, чтобы при повреждении обновления выполнялся откат, а не замена рабочего приложения, — @julianstorer

## Удалено
- Удалена интеграция с Crowdin после окончательной синхронизации переводов в [PR](https://github.com/wailsapp/wails/pull/5736) — @github-actions[bot]

## Безопасность
- Зарезервированные сообщения `wails:` из дочерних фреймов другого источника в Linux (GTK3 и GTK4) блокируются путём внедрения защитного скрипта во все фреймы, что снижает риск GHSA-2p5f-2gx8-mh8m
- `github.com/klauspost/compress` обновлён до 1.18.7 (GHSA-259r-337f-4rfw), а `golang.org/x/image` — до 0.45.0; сборка выполняется с Go 1.26.10, чтобы включить исправление `crypto/x509` (GO-2026-6809)
- `vite` обновлён до 8.2.4 или более поздней версии в пакете среды выполнения, мастере настройки, примерах и шаблонах проектов (GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv)

### Известные ограничения
- **Linux:** защита дочерних фреймов не может помешать дочернему фрейму *того же источника* (например, iframe `srcdoc` или `about:blank`) обращаться к нативному обработчику сообщений через родительский фрейм, поскольку WebKitGTK не предоставляет обработчику сведения о фрейме. Фреймы другого источника и изолированные фреймы без `allow-same-origin` блокируются. Не встраивайте в окно Wails ненадёжное содержимое того же источника.
- **Android (экспериментальная поддержка):** JavaScript-мост `window.wails` доступен каждому фрейму в WebView без проверки главного фрейма или источника. Содержимое встроенного iframe, в том числе содержимое другого источника, может вызывать привязанные методы и отправлять зарезервированные сообщения `wails:`. Не встраивайте ненадёжное стороннее содержимое в WebView приложения. Это отличается от macOS и iOS, где сообщения `wails:` из дочерних фреймов игнорируются. Исправление запланировано для 3.0.1.
- **Все платформы:** вызовы привязанных методов из дочерних фреймов передаются по HTTP в `/wails/runtime`, а не через нативный обработчик сообщений, поэтому защита их не ограничивает. Считайте, что любой встроенный iframe имеет доступ к вашим привязанным методам Go.

## v3.0.0-beta.34 - 2026-10-20

## Исправлено
- Исправлено восстановление фокуса веб-представлением Windows при активации окна, а также удержание клавиши Tab внутри веб-представления вместо перехода к нативному окну и обратно (#3783)
- Исправлены структуры таблиц виртуальных методов COM `internal/webview2/pkg/webview2` (`ICoreWebView2_2`..`_27` и другие производные интерфейсы), в которых отсутствовала цепочка унаследованных методов, из-за чего вызовы попадали в неверную позицию таблицы (#6129)
- Исправлены диалоги с вопросом в Windows, которые игнорировали пользовательские подписи кнопок и никогда не вызывали их обратные вызовы; теперь пользовательские кнопки отображаются с помощью TaskDialog, а если comctl32 v6 недоступен — с помощью MessageBox с позиционным сопоставлением обратных вызовов (#5794)
- Исправлено зависание `wails3 dev` на 60 с Vite: автоматическая проверка готовности фронтенда теперь использует TCP-зонд вместо HTTP GET, который отправлялся без заголовка `Accept` и получал ответ 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Меню системного трея продолжают работать при обновлении в открытом состоянии — [PR](https://github.com/wailsapp/wails/pull/6301), @taliesin-ai
- Добавлена корректная обработка временных ошибок COM WebView2 в операциях среды выполнения — [PR](https://github.com/wailsapp/wails/pull/6302), @leaanthony

## v3.0.0-beta.33 - 2026-10-19

## Добавлено
- Генератор привязок теперь создаёт модуль `bindings/events`, который экспортирует типизированную константу имени для каждого обнаруживаемого события, зарегистрированного с помощью `application.RegisterEvent` с именем, заданным константой времени компиляции (например, `UserLogin` для `user:login`). Её можно использовать с `Events.On`/`Emit` без плагина Vite (#5645)

## Исправлено
- Исправлены диалоги с вопросом в Windows, которые игнорировали пользовательские подписи кнопок и никогда не вызывали их обратные вызовы; теперь пользовательские кнопки отображаются с помощью TaskDialog, а если comctl32 v6 недоступен — с помощью MessageBox с позиционным сопоставлением обратных вызовов (#5794)
- Исправлено зависание `wails3 dev` на 60 с Vite: автоматическая проверка готовности фронтенда теперь использует TCP-зонд вместо HTTP GET, который отправлялся без заголовка `Accept` и получал ответ 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Меню системного трея продолжают работать при обновлении в открытом состоянии — [PR](https://github.com/wailsapp/wails/pull/6301), @taliesin-ai
- Добавлена корректная обработка временных ошибок COM WebView2 в операциях среды выполнения — [PR](https://github.com/wailsapp/wails/pull/6302), @leaanthony

## v3.0.0-beta.32 - 2026-10-18

## Добавлено
- Добавлены нативные заголовки разделов меню macOS для macOS 14+ — см. [сайт документации](https://v3.wails.io) (обновлено страниц: 20) — [PR](https://github.com/wailsapp/wails/pull/6297), @leaanthony

## Исправлено
- На macOS позиция SetRelativePosition теперь задаётся относительно рабочей области экрана — [PR](https://github.com/wailsapp/wails/pull/6299), @leaanthony
- После истечения времени ожидания родительского процесса или неудачной замены вспомогательный процесс обновления повторно запускает неизменённую целевую версию — [PR](https://github.com/wailsapp/wails/pull/6298), @taliesin-ai
- Для версионированной документации закреплена версия pnpm 10 — [PR](https://github.com/wailsapp/wails/pull/6293), @taliesin-ai

## v3.0.0-beta.31 - 2026-10-18

## Добавлено
- Добавлен экспорт подпути `@wailsio/runtime/cancellable`, чтобы `CancellablePromise` можно было импортировать без выполнения корневой точки входа среды выполнения (#6153)

## Изменено
- Теперь требуется Go 1.26 или более поздней версии: модуль v3 теперь объявляет `go 1.26.0`, и проекты, созданные с помощью `wails3 init` и `wails3 service init`, объявляют ту же минимальную версию

## Исправлено
- Исправлено изменение вывода генератора привязок при сборке с Go 1.27, где типы `encoding/json`, такие как `Marshaler`, `Unmarshaler` и `RawMessage`, становились псевдонимами типов `encoding/json/v2` и `encoding/json/jsontext`: теперь они снова отображаются как `any`, а модели для `encoding/json/v2` не создаются (#6292)
- Добавлена проверка URL удалённых шаблонов, переданных в `wails3 init -t`, перед клонированием: принимаются только URL `https://`, `http://`, `ssh://`, `git://` и `user@host:path`; значения передаются git после `--`, а при клонировании git разрешено использовать только эти протоколы передачи
- Значения заголовка `x-wails-window-id` вне допустимого диапазона в HTTP-транспорте среды выполнения теперь отклоняются, а не молча усекаются до идентификатора другого окна
- Усилена защита распаковки архивов обновления от zip-slip: теперь пути записей проверяются с помощью `filepath.IsLocal` и дополнительной проверки нахождения объединённого пути внутри целевого каталога; добавлены тесты для вложенных `..`, обратных косых черт, абсолютных путей и записей со сквозной записью через символические ссылки в архивах zip и tar.gz
- Исправлен сбой AppImage, созданных с помощью `wails3 generate appimage`, при запуске вне Debian и Ubuntu с ошибкой «Unable to spawn a new child process: … WebKitNetworkProcess». Теперь встроенный WebKit запускает вспомогательные процессы, поставляемые в AppImage. Песочница bubblewrap в WebKit остаётся включённой
- Исправлена ошибка компиляции сборок `-tags server` для Android и iOS; исправлено [#6271](https://github.com/wailsapp/wails/issues/6271)
- Исправлена ошибка замены файлов средством обновления в Linux, когда временный каталог находится в другой файловой системе (например, в tmpfs `/tmp`): теперь файлы сначала размещаются рядом с установленным исполняемым файлом и копируются между файловыми системами
- Исправлена попытка средства обновления заменить доступную только для чтения точку монтирования AppImage вместо файла `.AppImage`
- Исправлен приём зарезервированных управляющих сообщений окна `wails:` от iframe на macOS и iOS. Теперь они принимаются только из главного фрейма окна; остальные сообщения из iframe по-прежнему поступают в `RawMessageHandler` с информацией об их источнике
- Прекращена отправка `WindowDidUpdate` и `WindowWillUpdate` в macOS: `NSWindowDidUpdateNotification` срабатывает при каждом проходе цикла обработки событий, поэтому страница получала около 120 пустых вызовов `evaluateJavaScript:` в секунду, а в macOS 27.0 каждый из них поглощал кратковременную активацию от щелчка пользователя — `navigator.clipboard.writeText()` и другие API, требующие активации, отклоняли реальные щелчки с ошибкой `NotAllowedError` ([#6170](https://github.com/wailsapp/wails/issues/6170)), автор — @APshenkin
- Исправлен используемый по умолчанию язык Windows в info.json: теперь это en-US (0409), благодаря чему Проводник отображает сведения о версии (#6210)
- Исправлено отключение средства визуализации DMA-BUF в Linux на ноутбуках с дискретным GPU NVIDIA и встроенным GPU AMD или Intel: теперь в Wayland оно отключается только тогда, когда WebKit выполняет визуализацию на NVIDIA (приложение перенесено на дискретный GPU или на нём работает композитор), а не при любой загрузке модуля nvidia, из-за которой визуализация на встроенном GPU становилась крайне медленной; сеансы X11 сохраняют прежнее поведение, поскольку GPU визуализации Xorg задаётся его конфигурацией
- Добавлен резервный вариант с поэтапным копированием для обработки ошибок переименования между устройствами в средстве обновления Linux; [PR](https://github.com/wailsapp/wails/pull/6135), автор — @john-okeefe
- Исправлено отсутствие перерисовки строки меню функцией `WebviewWindow.SetMenu` в Windows, если тема строки меню не была настроена; исправлено [#6103](https://github.com/wailsapp/wails/issues/6103)
- Исправлена паника `WebviewWindow.SetMenu(nil)` в Windows: теперь операция, как и в macOS, ничего не делает; исправлено [#6104](https://github.com/wailsapp/wails/issues/6104)
- Добавлено сохранение записей журнала изменений при конкурирующих отправках изменений с помощью повторных попыток; [PR](https://github.com/wailsapp/wails/pull/6094), автор — @leaanthony
- Исправлен сбой `go mod vendor` с ошибкой `pattern arm64/WebView2Loader.dll: no matching files found` на всех платформах: удалены встраивания, ссылавшиеся на двоичные файлы, которые никогда не поставлялись в составе модуля; исправлены [#5782](https://github.com/wailsapp/wails/issues/5782) и [#5376](https://github.com/wailsapp/wails/issues/5376); [PR](https://github.com/wailsapp/wails/pull/6031), автор — @Grantmartin2002
- Предотвращено повторное уничтожение окна в macOS; [PR](https://github.com/wailsapp/wails/pull/6159), автор — @charlesbalderson
- Добавлена обработка WM_CONTEXTMENU, чтобы меню системного трея открывалось щелчком правой кнопки мыши; [PR](https://github.com/wailsapp/wails/pull/6162), автор — @wzk2025
- Исправлена бесконечная блокировка вызывающей стороны во всех пяти вариантах `InvokeSync*` при панике в переданном обратном вызове; исправлено [#6107](https://github.com/wailsapp/wails/issues/6107); [PR](https://github.com/wailsapp/wails/pull/6108), автор — @RALIST
- Добавлено безопасное повторное создание уничтоженных подменю GTK3; [PR](https://github.com/wailsapp/wails/pull/6236), автор — @leaanthony
- Реализована прослойка перенаправления рабочих потоков Windows для вложенных рабочих потоков; [PR](https://github.com/wailsapp/wails/pull/6261), автор — @taliesin-ai
- Исправлено игнорирование `PromptForMultipleSelection()` функцией `OpenFile()` в Windows при `CanChooseDirectories(true)`: теперь пользователи могут выбирать сразу несколько папок, как в macOS и Linux (#6263)
- Исправлена передача по значению аргументов `double` методам COM WebView2 в windows/arm64 — проблема затрагивала `ZoomFactor`, `RasterizationScale`, `Expires` файлов cookie, а также масштаб, размер страницы и поля в параметрах печати (ранее сеттеры параметров печати возвращали `ErrDoubleArgUnsupported`): теперь затронутые вызовы проходят через трамплин `runtime.cgocall`, загружающий `d0`, чего не может сделать путь системных вызовов Go (golang.org/issue/62583); автор — @APshenkin

## Безопасность
- Теперь серверный режим (`-tags server`) принимает вызовы `/wails/runtime` только как запросы `POST` из собственного источника сервера или источника, указанного в `ServerOptions.WebSocketOriginPatterns`; для других методов возвращается HTTP 405, а для других источников браузера — HTTP 403. При развёртывании за обратным прокси, переписывающим заголовок `Host`, добавьте публичный источник в `WebSocketOriginPatterns`
- Теперь серверный режим (`-tags server`), привязанный к адресу обратной петли (по умолчанию `localhost`, `127.0.0.0/8` или `::1`), отклоняет запросы с кодом HTTP 403, если заголовок `Host` не равен `localhost`, не является IP-адресом в литеральной форме, узлом привязки или узлом, соответствующим `ServerOptions.WebSocketOriginPatterns`. Это блокирует атаки с перепривязкой DNS на `/wails/runtime` и конечные точки WebSocket
- Уведомления Windows больше не используют PowerShell как резервный вариант при сбое API всплывающих уведомлений WinRT; вместо этого возвращается ошибка

## v3.0.0-beta.30 - 2026-10-15

## Исправлено
- Исправлена попытка средства обновления заменить доступную только для чтения точку монтирования AppImage вместо файла `.AppImage`
- Прекращена отправка `WindowDidUpdate` и `WindowWillUpdate` в macOS: `NSWindowDidUpdateNotification` срабатывает при каждом проходе цикла обработки событий, поэтому страница получала около 120 пустых вызовов `evaluateJavaScript:` в секунду, а в macOS 27.0 каждый из них поглощал кратковременную активацию от щелчка пользователя — `navigator.clipboard.writeText()` и другие API, требующие активации, отклоняли реальные щелчки с ошибкой `NotAllowedError` ([#6170](https://github.com/wailsapp/wails/issues/6170)), автор — @APshenkin
- Исправлен используемый по умолчанию язык Windows в info.json: теперь это en-US (0409), благодаря чему Проводник отображает сведения о версии (#6210)
- Добавлен резервный вариант с поэтапным копированием для обработки ошибок переименования между устройствами в средстве обновления Linux; [PR](https://github.com/wailsapp/wails/pull/6135), автор — @john-okeefe
- Исправлено отсутствие перерисовки строки меню функцией `WebviewWindow.SetMenu` в Windows, если тема строки меню не была настроена; исправлено [#6103](https://github.com/wailsapp/wails/issues/6103)
- Исправлена паника `WebviewWindow.SetMenu(nil)` в Windows: теперь операция, как и в macOS, ничего не делает; исправлено [#6104](https://github.com/wailsapp/wails/issues/6104)
- Добавлено сохранение записей журнала изменений при конкурирующих отправках изменений с помощью повторных попыток; [PR](https://github.com/wailsapp/wails/pull/6094), автор — @leaanthony
- Исправлен сбой `go mod vendor` с ошибкой `pattern arm64/WebView2Loader.dll: no matching files found` на всех платформах: удалены встраивания, ссылавшиеся на двоичные файлы, которые никогда не поставлялись в составе модуля; исправлены [#5782](https://github.com/wailsapp/wails/issues/5782) и [#5376](https://github.com/wailsapp/wails/issues/5376); [PR](https://github.com/wailsapp/wails/pull/6031), автор — @Grantmartin2002
- Предотвращено повторное уничтожение окна в macOS; [PR](https://github.com/wailsapp/wails/pull/6159), автор — @charlesbalderson
- Добавлена обработка WM_CONTEXTMENU, чтобы меню системного трея открывалось щелчком правой кнопки мыши; [PR](https://github.com/wailsapp/wails/pull/6162), автор — @wzk2025
- Исправлена бесконечная блокировка вызывающей стороны во всех пяти вариантах `InvokeSync*` при панике в переданном обратном вызове; исправлено [#6107](https://github.com/wailsapp/wails/issues/6107); [PR](https://github.com/wailsapp/wails/pull/6108), автор — @RALIST

## v3.0.0-beta.29 - 2026-10-13

## Добавлено
- Теперь глобальные сочетания клавиш в macOS учитывают раскладку клавиатуры: в раскладках, отличных от QWERTY (AZERTY, QWERTZ, Dvorak, ...), буквенные клавиши быстрого доступа привязываются через `UCKeyTranslate` к клавише с соответствующим обозначением, а не к фиксированной физической позиции QWERTY.

## Исправлено
- Нативное меню теперь инициализируется перед его установкой в Linux; [PR](https://github.com/wailsapp/wails/pull/6233), автор — @abhicodes-007
- Устранены регрессии готовности macOS к общедоступному выпуску в обработке вывода и мастере настройки; [PR](https://github.com/wailsapp/wails/pull/6237), автор — @taliesin-ai
- Исправлены обёртки `AddRef`/`Release` COM WebView2 в Windows, которые сообщали об устаревшей последней ошибке потока как о сбое. Теперь они возвращают счётчик ссылок, поэтому освобождение больше не может аварийно завершить приложение или привести к сбою запроса ресурса (#3931)
- Исправлен вызов слота `AddRef` из `ICoreWebView2Deferral.Release`, из-за которого на Windows не освобождалось отложенное выполнение каждого отложенного запроса ресурсов (#3931)

## v3.0.0-beta.28 - 2026-10-05

## Исправлено
- Вывод приложения macOS перенесён за пределы каталогов файловых провайдеров в [PR](https://github.com/wailsapp/wails/pull/6225), автор — @taliesin-ai
- Завершена работа над мастером настройки и добавлен вывод сведений о системе в [PR](https://github.com/wailsapp/wails/pull/6224), автор — @taliesin-ai
- Улучшены обработка ошибок восстановления WebView2 и очистка ресурсов на Windows в [PR](https://github.com/wailsapp/wails/pull/6222), автор — @leaanthony
- Улучшено восстановление WebView2 на Windows, чтобы предотвратить аварийное завершение приложений и зависание контроллеров, в [PR](https://github.com/wailsapp/wails/pull/6220), автор — @taliesin-ai
- Исправлены рекомендации по авторизации доступа к мультимедиа в macOS в [PR](https://github.com/wailsapp/wails/pull/6221), автор — @taliesin-ai
- Исправлена ошибка, из-за которой щелчок правой кнопкой мыши по значку в области уведомлений Linux не только открывал меню, но и запускал обработчик щелчка, поэтому связанное окно переключалось при каждом щелчке правой кнопкой мыши (#6018)
- Обработчики обратных вызовов WebView2 на Windows сохраняются до освобождения нативной стороной в [PR](https://github.com/wailsapp/wails/pull/6184), автор — @taliesin-ai

## v3.0.0-beta.27 - 2026-10-01

## Исправлено
- Исправлена ошибка, из-за которой параметр `Permissions` игнорировался в macOS 12 и более поздних версиях: теперь реализован метод захвата мультимедиа `WKUIDelegate`, поэтому `PermissionAllow` и `PermissionDeny` применяются к запросам доступа к камере и микрофону так же, как в Linux и Windows. По-прежнему необходимы `NSCameraUsageDescription` / `NSMicrophoneUsageDescription`, а при использовании песочницы — соответствующие разрешения устройств

## v3.0.0-beta.26 - 2026-09-25

## Добавлено
- Средство наблюдения применяет проверки готовности по HTTP к серверу разработки фронтенда в [PR](https://github.com/wailsapp/wails/pull/6150), автор — @atterpac

## Исправлено
- Добавлено восстановление после сбоев процессов WebView2 вместо отображения пустого окна в [PR](https://github.com/wailsapp/wails/pull/6002), автор — @taliesin-ai

## v3.0.0-beta.25 - 2026-09-22

## Исправлено
- `wails3 dev` и его фоновые процессы теперь завершаются при выходе основного приложения (#6048)

## v3.0.0-beta.24 - 2026-09-20

## Безопасность
- Разрешения GITHUB_TOKEN в рабочих процессах сборки и публикации ограничены в соответствии с принципом наименьших привилегий
- Устранены обход путей и выходы через символические ссылки в ПО промежуточного слоя ресурсов примера экрана
- Сохранена работоспособность установки зависимостей мастером настройки при её ограничении утверждёнными командами менеджеров пакетов

## v3.0.0-beta.23 - 2026-09-16

## Добавлено

- Добавлена возможность переопределять идентификатор приложения GTK в Linux в [PR](https://github.com/wailsapp/wails/pull/5972), автор — @overlordtm

## Исправлено

- Механизм единственного экземпляра в Linux регистрирует имя шины под собственным идентификатором приложения в [PR](https://github.com/wailsapp/wails/pull/5973), автор — @overlordtm
- Прерванные запросы пользовательских схем в Linux теперь отменяются в [PR](https://github.com/wailsapp/wails/pull/6126), автор — @leaanthony

## v3.0.0-beta.22 - 2026-09-14

## Добавлено

- Завершены переводы и перенос материалов сообщества — см. [сайт документации](https://v3.wails.io) (обновлено страниц: 1812) — в [PR](https://github.com/wailsapp/wails/pull/6122), автор — @leaanthony
- Добавлена документация по встраиваемым кнопкам Wails различных стилей и размеров в [PR](https://github.com/wailsapp/wails/pull/6120), автор — @leaanthony

## v3.0.0-beta.21 - 2026-09-13

## Добавлено

- Документация Wails v3 теперь обслуживается с помощью M-Press в [PR](https://github.com/wailsapp/wails/pull/6116), автор — @leaanthony

## Исправлено

- При создании журнала изменений значения slug в формате JSON теперь разбираются из метаданных MPD в [PR](https://github.com/wailsapp/wails/pull/6118), автор — @leaanthony
- После сбоев резервного копирования средство обновления очищает вспомогательные переменные окружения и повторно запускает исходную цель в [PR](https://github.com/wailsapp/wails/pull/6080), автор — @cnmax
- Обработчик сигналов по умолчанию теперь запускается во время App.Run в [PR](https://github.com/wailsapp/wails/pull/6098), автор — @leaanthony
- Меню Windows обрабатывает меню со значением nil, освобождает заменённые ресурсы и перерисовывает строку меню в [PR](https://github.com/wailsapp/wails/pull/6112), автор — @taliesin-ai
- Восстановлена упаковка MSIX для новых проектов, использующих общую конфигурацию YAML, в [PR](https://github.com/wailsapp/wails/pull/6115), автор — @leaanthony
- Исправлена ошибка загрузки сгенерированных привязок JavaScript и TypeScript, возникавшая, когда функции создания обобщённых моделей ссылались на объявленные позднее вспомогательные функции; также предотвращено переполнение стека при создании взаимозависимых обобщённых моделей (#6062)

## v3.0.0-beta.20 - 2026-09-10

## Изменено

- Ссылки демонстрационного проекта Clave обновлены и теперь указывают на актуальные сайт и репозиторий в [PR](https://github.com/wailsapp/wails/pull/6082), автор — @01xR4in

## Исправлено

- Прерванные запросы ресурсов Windows, включая запросы рабочих процессов, теперь отменяются, а обработчики keepalive сохраняются при навигации. На платформах Apple контексты нативных запросов передаются через оболочку приложения. (#5963, #5969)
- Записи журнала изменений сохраняются при конкурирующих отправках изменений благодаря повторным попыткам в [PR](https://github.com/wailsapp/wails/pull/6094), автор — @leaanthony
- Исправлен сбой `go mod vendor` с ошибкой `pattern arm64/WebView2Loader.dll: no matching files found` на всех платформах: удалены встраивания, ссылавшиеся на двоичные файлы, которые никогда не поставлялись в модуле; тем самым исправлены [#5782](https://github.com/wailsapp/wails/issues/5782) и [#5376](https://github.com/wailsapp/wails/issues/5376), в [PR](https://github.com/wailsapp/wails/pull/6031), автор — @Grantmartin2002

## Удалено

- Удалена поддержка нативного загрузчика WebView2, заменённого загрузчиком на чистом Go. Из встроенных ресурсов удалены двоичные файлы `WebView2Loader.dll`, а также зависимость `github.com/jchv/go-winloader`. Тег сборки `native_webview2loader` по-прежнему принимается и больше не вызывает ошибку, но не влияет на сборки v3, в [PR](https://github.com/wailsapp/wails/pull/6031), автор — @Grantmartin2002
- Из руководства по API macOS удалены неиспользуемые теги сборки и параметр FPS в [PR](https://github.com/wailsapp/wails/pull/6097), автор — @leaanthony

## v3.0.0-beta.19 - 2026-09-09

## Добавлено

- Доступ к закрытым API macOS ограничен тегами сборки и предоставляется только при явном включении — см. [документацию](https://v3.wails.io/features/browser/integration), [документацию](https://v3.wails.io/features/environment/info), [документацию](https://v3.wails.io/features/windows/basics), [документацию](https://v3.wails.io/features/windows/frameless), [документацию](https://v3.wails.io/features/windows/notch-windows), [документацию](https://v3.wails.io/features/windows/options), [документацию](https://v3.wails.io/guides/build/macos), [документацию](https://v3.wails.io/guides/build/private-macos-apis) и [документацию](https://v3.wails.io/reference/overview), в [PR](https://github.com/wailsapp/wails/pull/6087), автор — @leaanthony

## Исправлено

- Запросы среды выполнения размером более 64 МиБ теперь отклоняются с кодом HTTP 413 в [PR](https://github.com/wailsapp/wails/pull/6091), автор — @leaanthony

## Безопасность

- Усилена защита источников MCP и удалённого доступа с помощью аутентификации по токену в [PR](https://github.com/wailsapp/wails/pull/6092), автор — @leaanthony

## v3.0.0-beta.18 - 2026-09-08

## Исправлено

- Устранена утечка памяти Calloc в Linux и Darwin благодаря использованию получателей-указателей в [PR](https://github.com/wailsapp/wails/pull/6083) от @4RH1T3CT0R7

## v3.0.0-beta.17 - 2026-09-06

## Исправлено

- Windows: сбой или значение nil у `GetRequest` в обработчике WebResourceRequested больше не завершает процесс (паника `log.Fatal` из-за разыменования nil) — вместо этого запрос отбрасывается, а событие записывается в журнал; [PR](https://github.com/wailsapp/wails/pull/6006) от @midagedev

## v3.0.0-beta.16 - 2026-08-29

## Изменено

- Запрос пароля для нотаризации теперь отображается в новом окне терминала; [PR](https://github.com/wailsapp/wails/pull/6029) от @leaanthony

## Исправлено

- Исправлена обработка типов щелчков по значку в системном трее в macOS; [PR](https://github.com/wailsapp/wails/pull/5919) от @ChewbaccaCookie
- Перед обновлением CI удаляет неиспользуемые репозитории Microsoft apt; [PR](https://github.com/wailsapp/wails/pull/6041) от @Grantmartin2002

## v3.0.0-beta.15 - 2026-08-27

## Исправлено

- Тайм-аут встраивания WebView2 увеличен до 60 секунд; [PR](https://github.com/wailsapp/wails/pull/6043) от @Grantmartin2002

## v3.0.0-beta.14 - 2026-08-26

## Исправлено

- Исправлены названия нажатий сочетаний Control с буквенными клавишами в macOS; [PR](https://github.com/wailsapp/wails/pull/6032) от @taliesin-ai
- Исправлены значки формата ICO в системном трее и добавлено следование теме панели задач в Windows; [PR](https://github.com/wailsapp/wails/pull/6016) от @nik9play

## v3.0.0-beta.13 - 2026-08-25

## Исправлено

- В macOS обработка работы главного потока продолжается во время выполнения модального цикла; [PR](https://github.com/wailsapp/wails/pull/6026) от @leaanthony
- Операции с защищённым хранилищем на мобильных платформах теперь могут завершаться ошибкой и по умолчанию запрещают доступ при сбое; [PR](https://github.com/wailsapp/wails/pull/5923) от @mortenolsrud
- Обработчики событий приложения теперь запускаются, даже если не зарегистрирован ни один слушатель; [PR](https://github.com/wailsapp/wails/pull/5999) от @archy-rock3t-cloud
- Исправлены опечатки в комментариях и локализованной документации; [PR](https://github.com/wailsapp/wails/pull/6023) от @haoku123
- Удалены предварительно скомпилированные двоичные файлы macOS, зафиксированные в `v3/examples`; [PR](https://github.com/wailsapp/wails/pull/6025) от @4RH1T3CT0R7

## v3.0.0-beta.12 - 2026-08-21

## Добавлено

- Добавлены окна уведомлений для выреза экрана macOS с примером управления жизненным циклом и телеметрией — см. [документацию](https://v3.wails.io/features/windows/notch-windows); [PR](https://github.com/wailsapp/wails/pull/6010) от @leaanthony
- Добавлена поддержка окон NSPanel в macOS с новыми параметрами и нативной интеграцией — см. [документацию](https://v3.wails.io/features/windows/options); [PR](https://github.com/wailsapp/wails/pull/6008) от @leaanthony

## Исправлено

- Предотвращено зависание SQLite Prepare при параллельных вызовах; [PR](https://github.com/wailsapp/wails/pull/5998) от @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Удалено

- Из документации удалён устаревший трекер реализации; [PR](https://github.com/wailsapp/wails/pull/6005) от @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Исправлено

- Исправлена потеря аргументов запуска для пользовательских протоколов и ассоциаций файлов в хосте GTK4 для Linux; [PR](https://github.com/wailsapp/wails/pull/6000) от @midagedev
- В проверке журнала изменений исправлена обработка удалённых строк и исправлений из того же источника; [PR](https://github.com/wailsapp/wails/pull/5993) от @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Добавлено

- Добавлен защищённый сервер wails3 mcp для управления проектами с помощью агентов; [PR](https://github.com/wailsapp/wails/pull/5896) от @leaanthony
- Добавлена документация по моделям в привязках — см. [документацию](https://v3.wails.io/features/bindings/models); [PR](https://github.com/wailsapp/wails/pull/5988) от @taliesin-ai
- Добавлена поддержка установки через rpm-ostree в атомарных системах Linux; [PR](https://github.com/wailsapp/wails/pull/5987) от @leaanthony
- Добавлены нативное еженедельное создание и публикация графика истории звёзд — см. [документацию](https://v3.wails.io/credits), [документацию](https://v3.wails.io/de/credits), [документацию](https://v3.wails.io/fr/credits), [документацию](https://v3.wails.io/id/credits), [документацию](https://v3.wails.io/ja/credits), [документацию](https://v3.wails.io/ko/credits), [документацию](https://v3.wails.io/pt/credits), [документацию](https://v3.wails.io/ru/credits), [документацию](https://v3.wails.io/zh-cn/credits) и [документацию](https://v3.wails.io/zh-tw/credits); [PR](https://github.com/wailsapp/wails/pull/5986) от @leaanthony
- Добавлен доступный только в Darwin пакет mac для разрешения ресурсов пакета приложения — см. [документацию](https://v3.wails.io/guides/build/macos); [PR](https://github.com/wailsapp/wails/pull/5965) от @leaanthony
- Добавлены страница проекта Condui в галерее и запись в её указателе — см. [документацию](https://v3.wails.io/community/showcase/condui) и [документацию](https://v3.wails.io/community/showcase); [PR](https://github.com/wailsapp/wails/pull/5962) от @mgueregath
- Добавлена страница Redis Viewer в галерее со снимками экрана и ссылкой на проект — см. [документацию](https://v3.wails.io/community/showcase) и [документацию](https://v3.wails.io/community/showcase/redisviewer); [PR](https://github.com/wailsapp/wails/pull/5984) от @redisviewer

## Изменено

- Флаги приложения GTK для Linux обновлены до G<em>APPLICATION</em>NON_UNIQUE; [PR](https://github.com/wailsapp/wails/pull/5971) от @overlordtm
- События отсутствующих окон теперь записываются в журнал на уровне отладки, а не предупреждения; [PR](https://github.com/wailsapp/wails/pull/5914) от @julianstorer

## Исправлено

- Зарегистрированные сочетания клавиш macOS теперь имеют приоритет над webview; [PR](https://github.com/wailsapp/wails/pull/5902) от @julianstorer
- Исправлены неработающие ссылки на боковой панели документации; [PR](https://github.com/wailsapp/wails/pull/5937) от @northes
- Контексты запросов ресурсов в macOS и iOS теперь отменяются, когда WebKit прерывает соответствующую задачу пользовательской схемы (#5963)
- Добавлена обработка сообщения WindowSetFullscreenButtonEnabled; [PR](https://github.com/wailsapp/wails/pull/5976) от @archy-rock3t-cloud
- В шаблон preact-ts импортирован Fragment для устранения ошибки сборки; [PR](https://github.com/wailsapp/wails/pull/5979) от @haoku123
- Предотвращён сбой устаревших служебных приложений GTK3, когда обнаружение экранов выполняется до появления активного окна или доступного дисплея (#5966)
- Запуски выпуска с явно указанной версией теперь могут продолжаться, если раздел невыпущенных изменений пуст (#5977)

## Безопасность

- Файлы блокировки nanoid для веб-сайта обновлены до исправленной версии 3.3.18, чтобы устранить предупреждения безопасности; [PR](https://github.com/wailsapp/wails/pull/5985) от @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Добавлено

- В автоматические записи журнала изменений добавлено создание URL-адресов документации; [PR](https://github.com/wailsapp/wails/pull/5957) от @taliesin-ai
- Добавлены Streams: двунаправленные потоки байтов между Go и JavaScript с моделью программирования WebSocket, не требующие прослушивающего сокета. Объявите поток в Go с помощью `app.HandleStream(name, handler)` и подключитесь к нему из фронтенда с помощью `Stream(name)`, который возвращает объект, совместимый по структуре с `WebSocket`. Передача Go→JS выполняется через сервер ресурсов посредством одного удерживаемого опроса на окно, а JS→Go — посредством обычного POST-запроса; привязка к TCP-порту не выполняется, и данные не проходят через `evaluateJavaScript`. В серверных сборках (`-tags server`) тот же обработчик вместо этого обслуживается через настоящий WebSocket, поэтому код приложения одинаков для всех типов сборки. Автор: @leaanthony
- Запись об очереди сообщений перенесена в раздел «Не выпущено» в [PR](https://github.com/wailsapp/wails/pull/5935). Автор: @leaanthony

## Изменено

- Обновлены автоматическое формирование боковой панели документации и определение типа автора блога в [PR](https://github.com/wailsapp/wails/pull/5938). Автор: @leaanthony

## Исправлено

- При инициализации WebView2 теперь используются крайний срок и цикл обработки сообщений, см. [PR](https://github.com/wailsapp/wails/pull/5952). Автор: @leaanthony
- Тест файлов cookie WebView2 пропускается в CI, если он не включён явно, а его выполнение закрепляется за текущим потоком ОС, см. [PR](https://github.com/wailsapp/wails/pull/5951). Автор: @leaanthony
- Построители меню Windows восстанавливают идентификаторы команд для родительских элементов подменю, см. [PR](https://github.com/wailsapp/wails/pull/5944). Автор: @gilad-ch
- Официальный образ для кросс-компиляции приведён в соответствие с базовым уровнем поддержки GTK 4.14+ в Linux (#5928)
- Проект Xcode для iOS настроен так, чтобы сохранять унаследованные флаги компоновщика и добавлять -ObjC, см. [PR](https://github.com/wailsapp/wails/pull/5915). Автор: @mortenolsrud
- Устранено чрезмерно частое пересоздание TCP-соединений в прокси ресурсов `wails3 dev` на крупных фронтендах, которое могло исчерпать динамические порты хоста и привести к сбою посторонних процессов с ошибкой `EADDRNOTAVAIL`
- JavaScript событий каждого окна теперь ставится в очередь для упорядоченной доставки и обратного давления, см. [PR](https://github.com/wailsapp/wails/pull/5934). Автор: @leaanthony

## Удалено

- Удалён конвейер выпуска двоичных файлов для настольных систем: выпуски v3 теперь содержат только теги, а CLI `wails3` устанавливается с помощью `go install`. Удалены `release-v3.yml` и запускавший его этап ночной сборки, см. [PR](https://github.com/wailsapp/wails/pull/5946). Автор: @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Добавлено

- Добавлена настройка автовоспроизведения macOS, позволяющая отключить требование действия пользователя для воспроизведения мультимедиа, см. [PR](https://github.com/wailsapp/wails/pull/5512). Автор: @Eyalm321
- Запись об очереди сообщений перенесена в раздел «Не выпущено» в [PR](https://github.com/wailsapp/wails/pull/5935). Автор: @leaanthony

## Изменено

- Для более плавной анимации масштабирования в macOS теперь используется CADisplayLink или NSTimer, см. [PR](https://github.com/wailsapp/wails/pull/5945). Автор: @savely-krasovsky

## Исправлено

- Проект Xcode для iOS настроен так, чтобы сохранять унаследованные флаги компоновщика и добавлять -ObjC, см. [PR](https://github.com/wailsapp/wails/pull/5915). Автор: @mortenolsrud
- Устранено чрезмерно частое пересоздание TCP-соединений в прокси ресурсов `wails3 dev` на крупных фронтендах, которое могло исчерпать динамические порты хоста и привести к сбою посторонних процессов с ошибкой `EADDRNOTAVAIL`
- JavaScript событий каждого окна теперь ставится в очередь для упорядоченной доставки и обратного давления, см. [PR](https://github.com/wailsapp/wails/pull/5934). Автор: @leaanthony

### Добавлено

- Реализована универсальная асинхронная очередь сообщений FIFO для упорядоченной доставки событий, см. [PR](https://github.com/wailsapp/wails/pull/5851). Авторы: @savely-krasovsky и @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## Добавлено

- Реализованы ограниченное хранилище на стороне хоста для событий чрезмерного размера и их упорядоченная доставка в JavaScript, см. [PR](https://github.com/wailsapp/wails/pull/5930). Автор: @leaanthony
- Реализована анимация значка приложения в Dock macOS при мигании окна, см. [PR](https://github.com/wailsapp/wails/pull/5921). Автор: @julianstorer

## Исправлено

- При сбросе данных сервер ресурсов теперь сохраняет ошибки анализатора типа содержимого и ещё не записанные префиксы, см. [PR](https://github.com/wailsapp/wails/pull/5931). Автор: @leaanthony
- Предотвращён сбой приложений macOS при замене меню приложения из обратного вызова Wails
- Исправлены нечитаемые нативные меню в Windows 10 1809 / Windows Server 2019 (сборка 17763). Экспортируемые функции uxtheme для тёмного режима были доступны только начиная со сборки 18334, поэтому включение тёмного режима на уровне приложения на этих хостах никогда не выполнялось: фон меню отрисовывался тёмным, но Windows продолжала рисовать текст меню в светлой теме, из-за чего тёмный текст оказывался на тёмном фоне. Порядковые номера существуют начиная с 17763, поэтому теперь условие соответствует этому уровню.
- Исправлен вызов `GetDeviceCaps` вместо `GetStockObject` из `w32.GetStockObject`, из-за которого для каждого стандартного объекта возвращалось 0.
- Улучшены обработка ошибок и сообщения об ошибках при скачивании загрузчика WebView2, см. [PR](https://github.com/wailsapp/wails/pull/5924). Автор: @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## Исправлено

- Активация приложения в macOS теперь учитывает политику активации только для обычных приложений, см. [PR](https://github.com/wailsapp/wails/pull/5897). Автор: @julianstorer
- В сборках для Linux добавлена защита от обращения к неинициализированным окнам GTK, см. [PR](https://github.com/wailsapp/wails/pull/5898). Автор: @julianstorer
- Перед загрузкой URL для окон WebKit в Linux теперь явно задаётся непрозрачный цвет фона, см. [PR](https://github.com/wailsapp/wails/pull/5899). Автор: @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## Изменено

- Задачи сборки Android теперь по умолчанию используют arm64, а deploy-emulator выбирает архитектуру хоста, см. [PR](https://github.com/wailsapp/wails/pull/5890). Автор: @mortenolsrud

## Исправлено

- При перетаскивании теперь сохраняется состояние масштабирования окна macOS и уменьшается объём анимации, см. [PR](https://github.com/wailsapp/wails/pull/5900). Автор: @leaanthony
- Исправлена сборка для Windows в серверном режиме: в ограничение сборки `webview_window_windows_nonclient.go` добавлено `!server`

## v3.0.0-beta.3 - 2026-08-03

## Добавлено

- В сведениях о реализации документировано завершение бета-проверки этапа 10, см. [PR](https://github.com/wailsapp/wails/pull/5881). Автор: @leaanthony

## Исправлено

- В API тёмного режима Windows теперь передаётся дескриптор окна и проверяются аргументы, см. [PR](https://github.com/wailsapp/wails/pull/5877). Автор: @leaanthony
- Определение состояния кнопок строки заголовка для безрамочных окон macOS централизовано, см. [PR](https://github.com/wailsapp/wails/pull/5870). Автор: @taliesin-ai
- Предотвращено появление нечитаемого текста в нативном меню, когда приложение Windows запрашивает тёмный режим, а для приложений Windows выбрана светлая тема. Теперь меню использует соответствующий светлый нативный фон, пока Windows не сможет отрисовать светлый текст меню для тёмного режима.
- Исправлены нечитаемые нативные меню в Windows 10 1809 / Windows Server 2019 (сборка 17763). Экспортируемые функции uxtheme для тёмного режима были доступны только начиная со сборки 18334, поэтому включение тёмного режима на уровне приложения на этих хостах никогда не выполнялось: фон меню отрисовывался тёмным, но Windows продолжала рисовать текст меню в светлой теме, из-за чего тёмный текст оказывался на тёмном фоне. Порядковые номера существуют начиная с 17763, поэтому теперь условие соответствует этому уровню.

## v3.0.0-beta.2 - 2026-08-02

## Изменено

- v3 переведён из альфа-версии в бета-версию
- Документированы интеллектуальные настройки по умолчанию для системного трея и поведение автоматического скрытия всплывающего меню; добавлено регрессионное покрытие выбора обработчика щелчка (#5840).
- Средство обновления через GitHub теперь по умолчанию исключает ресурсы установщика Windows, см. [PR](https://github.com/wailsapp/wails/pull/5861). Автор: @leaanthony
- Добавлена поддержка скруглённых, прямых и настраиваемых углов для безрамочных окон macOS в [PR](https://github.com/wailsapp/wails/pull/5866) от @leaanthony

## Исправлено

- Добавлено получение актуальных размеров окон GTK4 и отправка событий изменения размера, разворачивания, сворачивания и перехода в полноэкранный режим с настроенной поверхности (#5830).
- Исправлен сбой Linux WebKit при отправке Blob или FormData в запросах fetch в [PR](https://github.com/wailsapp/wails/pull/5854) от @taliesin-ai
- Обёртка fetch передаёт undefined при отсутствии заголовков Blob/FormData в [PR](https://github.com/wailsapp/wails/pull/5865) от @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## Добавлено

## Изменено

- Добавлена поддержка скруглённых, прямых и настраиваемых углов для безрамочных окон macOS в [PR](https://github.com/wailsapp/wails/pull/5866) от @leaanthony

## Исправлено

- Обёртка fetch передаёт undefined при отсутствии заголовков Blob/FormData в [PR](https://github.com/wailsapp/wails/pull/5865) от @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## Добавлено

- Добавлена поддержка упаковки macOS-приложений в DMG с новыми параметрами и задачами сборки в [PR](https://github.com/wailsapp/wails/pull/5857) от @leaanthony

## Изменено

- Средство обновления через GitHub теперь по умолчанию исключает артефакты установщика Windows в [PR](https://github.com/wailsapp/wails/pull/5861) от @leaanthony

## Исправлено

- Исправлен сбой Linux WebKit при отправке Blob или FormData в запросах fetch в [PR](https://github.com/wailsapp/wails/pull/5854) от @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## Добавлено

- Реализованы действия разворачивания и сворачивания окна по двойному щелчку на строке заголовка macOS в [PR](https://github.com/wailsapp/wails/pull/5853) от @taliesin-ai

## Изменено

- Руководство по службе QR обновлено для использования NewServiceWithOptions; также добавлены отступы в [PR](https://github.com/wailsapp/wails/pull/5849) от @jeongkyu

## Исправлено

- Сохранена отзывчивость WKWebView при масштабировании в macOS в [PR](https://github.com/wailsapp/wails/pull/5856) от @leaanthony
- Исправлены запросы размеров окон GTK4 и добавлена отправка событий изменения размера, разворачивания, сворачивания и перехода в полноэкранный режим из настроенного `GdkSurface`.

## v3.0.0-alpha2.119 - 2026-07-27

## Исправлено

- Документация на нескольких языках дополнена архитектурными диаграммами в [PR](https://github.com/wailsapp/wails/pull/5833) от @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## Добавлено

- Добавлены пути по умолчанию для входных и выходных файлов генератора значков в [PR](https://github.com/wailsapp/wails/pull/5825) от @taliesin-ai
- Исходные модули точек входа добавлены в sideEffects файла package.json пакета среды выполнения в [PR](https://github.com/wailsapp/wails/pull/5797) от @savely-krasovsky
- В руководство для участников добавлен раздел о лицензии и происхождении кода в [PR](https://github.com/wailsapp/wails/pull/5816) от @taliesin-ai

## Исправлено

- Для удаления радиуса скругления границ применены локализованные стили CSS безрамочных окон GTK4 в [PR](https://github.com/wailsapp/wails/pull/5800) от @savely-krasovsky
- Добавлена корректная обработка сбоев получения позиции курсора в Windows для всплывающих меню и перечисления экранов в [PR](https://github.com/wailsapp/wails/pull/5789) от @wayneforrest
- Диалог открытия файлов macOS теперь правильно фильтрует расширения и проверяет допустимость файлов по суффиксу в [PR](https://github.com/wailsapp/wails/pull/5678) от @phergul
- Инициализация тёмного режима Windows защищена от вызовов nil API в [PR](https://github.com/wailsapp/wails/pull/5793) от @roachadam
- Исправлен сбой сборки средства обновления для 32-разрядной платформы: константа `maxArchiveTotalSize` (2 ГиБ) переполняла платформенный тип `int` при передаче в `fmt.Errorf` на `GOARCH=386`. Теперь ей явно назначен тип `int64`.
- Исправлена паника из-за nil-указателя при запуске, когда окно использует тёмную строку заголовка (или системную тёмную тему) в сборках Windows, которые не загружают API тёмного режима uxtheme, например в Windows 10 1809 / Windows Server 2019 (сборка 17763). Вызовы `AllowDarkModeForWindow` при настройке темы окна теперь защищены от nil аналогично уже используемой защите в `w32.SetMenuTheme`.

## v3.0.0-alpha2.117 - 2026-07-08

## Добавлено

- Реализована пользовательская логика hit test для неклиентских областей в Windows в [PR](https://github.com/wailsapp/wails/pull/5462) от @savely-krasovsky

## Изменено

- Обнаружение масштаба монитора WebView2 настроено в зависимости от UseVisualHosting в [PR](https://github.com/wailsapp/wails/pull/5761) от @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## Добавлено

- Документация FAQ обновлена с акцентом на возможности и рекомендации Wails v3 в [PR](https://github.com/wailsapp/wails/pull/5763) от @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## Исправлено

- Исправлена ошибка, из-за которой `Menu.Update()` не перестраивал нативное меню в GTK4 Linux (#5659; независимо диагностировано и исправлено @puneetdixit200 в #5539)
- Исправлен сбой при перечислении экранов macOS после изменения конфигурации дисплеев: строки идентификаторов и имён экранов теперь копируются, а их количество фиксируется в снимке состояния (#5565; независимо диагностировано и исправлено @x-haose в #5584)
- Исправлен сбой в Windows, когда `WM_ERASEBKGND` закрашивает фон сплошным цветом во время перехода между свёрнутым и восстановленным состояниями, при котором `GetClientRect` возвращает nil (о необходимости проверки сообщил @sinspired в #5636)
- Исправлена ошибка привязок фронтенда, из-за которой данные всегда разбирались как текст, от @mbaklor в #5690
- Исправлен сбой сборки Windows при использовании тега сборки `server`, вызванный отсутствием в файлах графического интерфейса Windows ограничения сборки `!server`, которое уже есть в эквивалентных файлах macOS и Linux (#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## Добавлено

- Реализованы протокол манифеста обновлений и поставщик конечной точки в [PR](https://github.com/wailsapp/wails/pull/5720) от @taliesin-ai

## Изменено

- Привязка `webview2` включена в модуль v3 под именем `v3/internal/webview2`; отдельный модуль, его процессы ночного выпуска и синхронизации, а также манипуляции с версией go.mod удалены (его единственным потребителем был v3) в [PR](https://github.com/wailsapp/wails/pull/5711) от @taliesin-ai

## Исправлено

- Исправление обнаружения масштаба монитора WebView2 и повторной синхронизации хоста при изменении DPI перенесено в раздел «Не выпущено» в [PR](https://github.com/wailsapp/wails/pull/5750) от @taliesin-ai
- Обновлён COM-маршалинг WebView2 для параметров float64 и BOOL в [PR](https://github.com/wailsapp/wails/pull/5741) от @wayneforrest
- Предотвращены паника и разыменование nil при обновлении и удалении значков в области уведомлений Windows в [PR](https://github.com/wailsapp/wails/pull/5703) от @wayneforrest
- Исправлена ошибка, из-за которой скрытые окна в Windows не скрывались повторно должным образом, в [PR](https://github.com/wailsapp/wails/pull/5743) от @wayneforrest
- Видимость контроллера WebView2 синхронизирована со сворачиванием, разворачиванием и восстановлением окна в [PR](https://github.com/wailsapp/wails/pull/5742) от @wayneforrest

### Исправлено

- Повторно включить определение масштаба монитора в WebView2 и выполнять повторную синхронизацию хоста при изменении DPI в [PR](https://github.com/wailsapp/wails/pull/5734) от @taliesin-ai на основе исправления, проверенного @randalmurphal, с подтверждением первопричины от @eleclin и тестированием на оборудовании от @qq540491950

## v3.0.0-alpha2.113 - 2026-07-04

## Добавлено

- Добавить предупреждение при сборке релизного AAB без заданного `ANDROID_KEYSTORE_FILE` (Google Play отклоняет пакеты с отладочной подписью) и документировать упаковку и подписание App Bundle в [PR](https://github.com/wailsapp/wails/pull/5730) от @taliesin-ai
- Добавить документацию «Почему Wails» на нескольких языках в [PR](https://github.com/wailsapp/wails/pull/5739) от @taliesin-ai
- Добавить в привязках сопоставление Go time.Time с JS Date или string в [PR](https://github.com/wailsapp/wails/pull/5398) от @fbbdev
- Добавить задачи упаковки Android App Bundle (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) для публикации в Play Store; задачи APK остаются для локального тестирования и тестирования в эмуляторе — в [PR](https://github.com/wailsapp/wails/pull/5728) от @mortenolsrud (исправляет [#5726](https://github.com/wailsapp/wails/issues/5726))
- Добавить цели задач для физических устройств Android и возобновить запрос разрешений на камеру и местоположение в [PR](https://github.com/wailsapp/wails/pull/5735) от @taliesin-ai

## Изменено

- Обновить `webview2` до v1.0.28 ([примечания к выпуску](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- Обновить `compileSdk`/`targetSdk` в шаблоне Android с 34 до 35, как того требует Google Play для публикации новых приложений, в [PR](https://github.com/wailsapp/wails/pull/5730) от @taliesin-ai

## Исправлено

- Исправить встроенные маски аватаров в sponsorkit в [PR](https://github.com/wailsapp/wails/pull/5745) от @leaanthony
- Исправить выбор неправильного образа системы или версии cmdline-tools при автоматическом создании Android AVD из-за лексикографической сортировки версий в [PR](https://github.com/wailsapp/wails/pull/5730) от @taliesin-ai
- Исправить устаревшую версию Android NDK, предлагаемую мастером настройки (теперь 26.3.11579264, в соответствии с документированным требованием), в [PR](https://github.com/wailsapp/wails/pull/5730) от @taliesin-ai
- Обновить французскую документацию по SvelteKit и параметрам в [PR](https://github.com/wailsapp/wails/pull/5744) от @leaanthony
- Исправить SIGSEGV при перечислении экранов macOS во время изменения конфигурации дисплеев в [PR](https://github.com/wailsapp/wails/pull/5516) от @flofreud

## v3.0.0-alpha2.112 - 2026-07-03

## Добавлено

- Добавить генератор SVG со списком участников на Go и обновить страницы благодарностей в документации и на сайте в [PR](https://github.com/wailsapp/wails/pull/5724) от @taliesin-ai
- Добавить в привязках сопоставление Go time.Time с JS Date или string в [PR](https://github.com/wailsapp/wails/pull/5398) от @fbbdev

## Изменено

- Заменить конвейер создания изображений спонсоров на базе Node генератором на Go в [PR](https://github.com/wailsapp/wails/pull/5719) от @taliesin-ai

## Исправлено

- Исправить скрипт установки зависимостей ресурсов сборки Android в [PR](https://github.com/wailsapp/wails/pull/5729) от @taliesin-ai
- Отклонять управляющий символ U+0085 (NEXT LINE) в `ValidateAndSanitizeURL`, дополнив охват пробельных символов валидатором URL
- Пересчитывать рамку DWM при изменении DPI для окон без рамки в [PR](https://github.com/wailsapp/wails/pull/4785) от @leaanthony
- Исправить сбой определения зоны перетаскивания DnD при масштабе Windows, отличном от 100%, в [PR](https://github.com/wailsapp/wails/pull/4632) от @yulesxoxo
- Добавить явное управление памятью Objective-C для объектов Cocoa в диалогах, меню, системном трее и уведомлениях Darwin в [PR](https://github.com/wailsapp/wails/pull/5714) от @taliesin-ai
- Исправить ошибки бэкенда CGO в Linux и проблемы с системным треем в [PR](https://github.com/wailsapp/wails/pull/5718) от @taliesin-ai

## v3.0.0-alpha2.111 - 2026-07-01

## Добавлено

- Добавить HappyTools в витрину сообщества в [PR](https://github.com/wailsapp/wails/pull/5061) от @Aliuyanfeng
- Добавить поддержку индонезийской локали и полную документацию в [PR](https://github.com/wailsapp/wails/pull/5643) от @triadmoko
- Добавить параметр DisableMenu в WindowsWindow в [PR](https://github.com/wailsapp/wails/pull/4813) от @leaanthony

## Изменено

- Обновить шаблон Taskfile и CLI для диспетчеризации задач сборки и упаковки с GOOS и ARCH в [PR](https://github.com/wailsapp/wails/pull/5617) от @leaanthony

## Исправлено

- Исправить проблему со вкладками окон macOS в [PR](https://github.com/wailsapp/wails/pull/5708) от @taliesin-ai

## Удалено

- Удалить переведённые на немецкий файлы MDX из разделов об участии в разработке, возможностях и руководствах в [PR](https://github.com/wailsapp/wails/pull/5702) от @taliesin-ai

## v3.0.0-alpha2.110 - 2026-06-30

## Добавлено

- Реализовать обычную и принудительную перезагрузку WebView в macOS, а также добавить восстановление после завершения процесса WebContent в [PR](https://github.com/wailsapp/wails/pull/5129) от @wayneforrest
- Добавить полную немецкую документацию по участию в разработке, возможностям и руководствам в [PR](https://github.com/wailsapp/wails/pull/5396) от @leaanthony
- Расширить уведомления поддержкой звука, вложений, планирования и API обновления в [PR](https://github.com/wailsapp/wails/pull/5333) от @popaprozac

## Исправлено

- Пересчитывать рамку DWM при изменении DPI для окон без рамки в [PR](https://github.com/wailsapp/wails/pull/4785) от @leaanthony
- Исправить сбой определения зоны перетаскивания DnD при масштабе Windows, отличном от 100%, в [PR](https://github.com/wailsapp/wails/pull/4632) от @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Добавлено

- Добавить примеры кода в документацию EventsEmit в [PR](https://github.com/wailsapp/wails/pull/5026) от @iamhabbeboy
- Добавить в Windows параметр визуального хостинга WebView2 в [PR](https://github.com/wailsapp/wails/pull/5380) от @MerIijn
- Добавить Klustr в документацию витрины сообщества в [PR](https://github.com/wailsapp/wails/pull/5536) от @SametKUM
- Добавить Kira в витрину сообщества вместе с новыми страницами и записью в журнале изменений в [PR](https://github.com/wailsapp/wails/pull/5685) от @thiennguyen93
- Добавить раздел обратной связи в руководство по службе MCP в [PR](https://github.com/wailsapp/wails/pull/5694) от @taliesin-ai

## Изменено

- Для серверного режима теперь предусмотрена полноценная производственная сборка, согласованная с задачами сборки настольного приложения (#5693). По умолчанию `task build:server` собирает производственный двоичный файл (`-tags server,production`, `-trimpath`, без отладочных символов) и принимает `DEV=true` (сервер разработки), `OBFUSCATED=true` (garble) и `EXTRA_TAGS`. `task run:server` запускает сервер разработки. `Dockerfile.server` / `task build:docker` сначала собирают производственный сервер (`-tags server,production`) и производственный фронтенд; по умолчанию образ содержит статическую сборку на чистом Go в distroless/static, а `CGO_ENABLED`, `GO_IMAGE` и `RUNTIME_IMAGE` доступны как переопределяемые аргументы сборки для приложений CGO.

## Исправлено

- Предотвратить сбой при закрытии окна с ожидающими асинхронными вызовами в [PR](https://github.com/wailsapp/wails/pull/4435) от @leaanthony
- Предотвратить активацию окна при открытии скрытых приложений в Windows в [PR](https://github.com/wailsapp/wails/pull/5249) от @leaanthony
- Обеспечить выполнение обработки метаданных запросов WebKit, завершения ответов и потоков тела в главном потоке GTK в [PR](https://github.com/wailsapp/wails/pull/5668) от @taliesin-ai
- Исправлена ошибка, из-за которой `Menu.Update()` не перестраивал нативное меню в GTK4 на Linux (#5659; независимо диагностировано и исправлено @puneetdixit200 в #5539)
- Исправлен сбой при перечислении экранов macOS после изменения конфигурации дисплеев: строки с идентификаторами и именами экранов теперь копируются, а их количество фиксируется в снимке (#5565; независимо диагностировано и исправлено @x-haose в #5584)
- Исправлена ошибка, из-за которой содержимое WebView2 сжималось, а затем исчезало после перетаскивания окна между мониторами с разным DPI в Windows: в обработчике `WM_DPICHANGED` теперь повторно задаются границы контроллера аналогично повторной синхронизации DPI при выходе из свёрнутого состояния (#5677)

## v3.0.0-alpha2.108 — 2026-06-28

## Добавлено

- Добавлены глобальные (общесистемные) сочетания клавиш через `app.GlobalShortcut` (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Сочетания срабатывают, даже когда приложение не находится в фокусе. Для каждой платформы реализовано нативно, без сторонних зависимостей: горячие клавиши Carbon в macOS, `RegisterHotKey` в Windows, `XGrabKey` в X11 и интерфейс глобальных сочетаний клавиш XDG Desktop Portal в Wayland.
- Добавлен встроенный сервер MCP — сервер Model Context Protocol, который автоматически запускается, если приложение собрано с тегом `mcp`. Он позволяет LLM-агентам тестировать работающие приложения Wails и управлять ими: управлять окнами, исследовать DOM, выполнять JavaScript, вызывать привязанные методы, работать с событиями и имитировать ввод с помощью мыши и клавиатуры с отображением анимированного экранного курсора. Пользовательский код не требуется: `mcp` автоматически добавляет тег `wails3 build`/`wails3 dev`, когда задана переменная `WAILS_MCP=1`. Сервер полностью настраивается через переменные окружения (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Исправлено

- Исправлена ошибка, из-за которой `Menu.Update()` не перестраивал нативное меню в GTK4 на Linux (#5659; независимо диагностировано и исправлено @puneetdixit200 в #5539)
- Исправлен сбой при перечислении экранов macOS после изменения конфигурации дисплеев: строки с идентификаторами и именами экранов теперь копируются, а их количество фиксируется в снимке (#5565; независимо диагностировано и исправлено @x-haose в #5584)

## v3.0.0-alpha2.107 — 2026-06-27

## Добавлено

- Добавлена экспериментальная документация по Wake с навигацией на боковой панели в [PR](https://github.com/wailsapp/wails/pull/5613) от @leaanthony

## v3.0.0-alpha2.106 — 2026-06-24

## Изменено

- Версия `webview2` повышена до v1.0.27.
  - ci(webview2): исправлена релизная сборка (кросс-компиляция для Windows и полный go.sum) (#5671)\

  **Полный список изменений:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- Из рабочего процесса выпуска webview2 удалён go vet при кросс-компиляции в [PR](https://github.com/wailsapp/wails/pull/5672) от @taliesin-ai
- Модель OpenRouter для автоматического формирования журнала изменений обновлена до google/gemini-2.5-flash-lite в [PR](https://github.com/wailsapp/wails/pull/5670) от @taliesin-ai
- Версия `webview2` повышена до v1.0.26.

### Исправления

- **Вместо завершения работы добавлено восстановление после временных ошибок COM во время выполнения** (#5658, #5580). Ранее `Chromium.errorCallback` вызывал `os.Exit(1)` при *любой* ошибке COM, поэтому устранимый временный сбой после запуска завершал работу всего приложения. Теперь пути выполнения (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) регистрируют ошибку в журнале и восстанавливаются. В частности, некорректное или недоверенное веб-сообщение в `MessageReceived` теперь отбрасывается, а не завершает процесс. Это устраняет класс сбоев при перемещении между мониторами с разным DPI (#5544, #5650). Ошибки в путях создания среды и контроллера по-прежнему являются фатальными.\

**Полный список изменений:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Исправлено

- Исправлен рабочий процесс release-webview2, чтобы он корректно обрабатывал файлы go.sum, в [PR](https://github.com/wailsapp/wails/pull/5671) от @taliesin-ai
- Исправлено обновление меню GTK4 в Linux: нативное меню теперь очищается и перестраивается, в [PR](https://github.com/wailsapp/wails/pull/5659) от @taliesin-ai

## v3.0.0-alpha2.105 — 2026-06-21

## Добавлено

- Добавлен `application.System` для определения платформы во время выполнения из общего кода: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (тег сборки `server`), а также `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` для непосредственного тестирования одной целевой платформы. Код компилируется для всех целевых платформ, поэтому ветвление можно выполнять без тегов сборки. Соответствующие вспомогательные средства для фронтенда (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) доступны в `@wailsio/runtime`
- Добавлено руководство «Использование других фронтенд-фреймворков», в котором показано, как поместить собственный проект Vite в `frontend/` (рассматриваются Solid, Preact, Lit, SvelteKit, Qwik, Angular и другие)
- Мастер `wails3 setup` теперь проверяет инструментарий для мобильных платформ (iOS/Android): Xcode и среду выполнения iOS Simulator, JDK, Android SDK/NDK и эмулятор. Там, где это применимо, доступны установка одним нажатием и исправления конфигурации оболочки, которые можно скопировать
- Создаваемые проекты теперь включают `frontend/.npmrc`, задающий `minimum-release-age` длительностью 7 дн., чтобы снизить риск использования недавно опубликованных (и потенциально скомпрометированных) пакетов (учитывается pnpm и bun; npm без каких-либо последствий игнорирует эту настройку)

## Изменено

- Все встроенные стартовые шаблоны получили новый дизайн с главным блоком в стиле неоновых гор (для веб-платформы, iOS и Android)
- **Теперь TypeScript используется в стартовых шаблонах по умолчанию и ему соответствует имя шаблона без суффикса.** `wails3 init` (без `-t`) создаёт проект TypeScript; шаблоны `-t vanilla`, `-t react`, `-t vue` и `-t svelte` используют TypeScript, а их варианты на JavaScript доступны как `-t vanilla-js`, `-t react-js`, `-t vue-js` и `-t svelte-js`. Встроенные шаблоны указывают язык с помощью `typescript:` в `template.yaml`; для совместимости шаблоны сообщества с суффиксом `-ts` продолжают работать как резервный вариант
- Мастер `wails3 setup` получил новый дизайн в неоновой теме «digital Wails» (эффект матового стекла поверх горного фона)

## Исправлено

- Исправлен сбой в Windows при восстановлении приложения, которое оставалось свёрнутым достаточно долго, чтобы WebView2 перешёл в состояние приостановки или его процесс рендеринга/GPU был перезапущен. Теперь повторная синхронизация DPI при сворачивании и восстановлении (#5544) обращается к контроллеру WebView2 только в том случае, если DPI окна действительно изменился. Это предотвращает фатальные вызовы COM к приостановленному контроллеру при обычном восстановлении без изменения DPI (#5605)
- Исправлены повторяющиеся нативные сбои `SIGABRT`/`SIGSEGV` (обычно внутри `g_object_unref` во время основного цикла GTK) в длительно работающих приложениях Linux при частой загрузке ресурсов и мультимедиа. Сервер ресурсов завершал `WebKitURISchemeRequest` из рабочих горутин, вызывая непотокобезопасные функции WebKit2GTK вне главного потока GTK; теперь завершение (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) выполняется в главном потоке. Это дополняет частичное исправление из #5566. Проблема затрагивала сборки 6.0 как с GTK3, так и с GTK4/WebKitGTK (#5631, #5557)
- Исправлен периодически возникавший `fatal error: invalid pointer found on stack` в `setupSignalHandlers` на Linux/GTK3. Идентификаторы окон, передаваемые в качестве сигнальных `user_data`, хранились в локальной переменной Go типа `unsafe.Pointer`, поэтому сборщик мусора аварийно завершал работу, когда при копировании стека сканировал это значение, не являющееся указателем. Теперь на стороне Go идентификатор хранится как целое число (`uintptr_t`). Таким образом, в устаревший путь GTK3 перенесено исправление из #4958, ранее применённое к пути GTK4 (где функции сигналов C были переведены на `uintptr_t` для устранения ошибок `-race`/checkptr) (#5631)

## Удалено

- Удалены стартовые шаблоны `react-swc`, `preact`, `lit`, `solid`, `qwik` и `sveltekit`, а также их варианты `-ts`. Теперь поддерживаемый набор встроенных шаблонов включает `vanilla`, `react`, `vue` и `svelte`: каждый из них по умолчанию использует TypeScript и имеет вариант на JavaScript с `-js`. Любой другой фреймворк по-прежнему можно использовать, [подключив собственный фронтенд](https://v3.wails.io/guides/dev/frontend-frameworks) или пользовательский шаблон

## v3.0.0-alpha2.104 — 2026-06-18

## Исправлено

- Исправлен сбой iOS (SIGABRT), возникавший, когда привязанный метод сервиса Go возвращал пустую строку. В обработчике ответа для ресурсов iOS указатель на тело проверялся с помощью `buf != nil` вместо проверки его длины, поэтому тело нулевой длины вызывало панику в `&buf[0]`; теперь проверяется длина, как и в обработчиках для настольных платформ

## v3.0.0-alpha2.103 — 2026-06-15

## Изменено

- Нативные функции iOS и Android перенесены в менеджеры платформ: теперь вместо прежних свободных функций `application.IOS*`/`application.Android*` они вызываются через `application.IOS.*` и `application.Android.*` (например, `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`) (#5602)
- Переименованы события мобильного моста: кроссплатформенные события теперь используют префикс `common:*` (например, `common:haptic`, `common:location`), а события, относящиеся только к конкретной платформе, — `ios:*` / `android:*` (например, `ios:backgroundTask`, `android:foregroundService`); префикс `native:*` больше не используется (#5602)

## v3.0.0-alpha.102 — 2026-06-14

## Добавлено

- Добавлен экспериментальный мастер `wails3 setup` для интерактивной настройки проекта и проверки зависимостей
- В `wails3 doctor` добавлен флаг `--json` для вывода в машиночитаемом формате
- В команду `wails3 doctor` добавлен раздел со статусом подписывания

## Исправлено

- Исправлено обнаружение npm в Linux: теперь помимо менеджера пакетов проверяется PATH

## v3.0.0-alpha.101 — 2026-06-13

## Добавлено

- iOS: нативные диалоги сообщений (UIAlertController) и диалоги открытия файла, файлов или каталога (UIDocumentPickerViewController); диалоги сохранения возвращают явную ошибку
- iOS: поддержка буфера обмена через UIPasteboard
- iOS: реальные параметры экрана через UIScreen (точки, пиксели, масштаб и рабочая область с учётом безопасной зоны)
- iOS: сборки для устройств (`IOS_PLATFORM=device`), поддержка идентификатора подписи кода, профиля подготовки и разрешений, упаковка `.ipa` и `deploy-device` через devicectl
- iOS: настраиваемая минимальная версия iOS (`ios.minIOSVersion` в build/config.yml)
- iOS: `wails3 doctor` сообщает о доступности Xcode и iOS SDK в macOS
- iOS: системные события — состояние батареи, сети, темы, блокировки экрана и нехватки памяти предоставляются как события приложения `events.IOS.*` и платформонезависимые события `events.Common.*`
- iOS: нативный мост мобильных функций (экспортируемый `application.IOS*`) — системное меню «Поделиться», открытие URL, предотвращение перехода в спящий режим, фонарик, отступы безопасной зоны, яркость, сведения о приложении, блокировка ориентации, строка состояния, биометрия (Face ID/Touch ID), локальные уведомления и защищённое хранилище Keychain
- iOS: датчики и оборудование — тактильная обратная связь, однократное определение геопозиции, акселерометр, датчик приближения, преобразование текста в речь, сведения о хранилище, состояние питания и батареи, состояние сети, отступы для клавиатуры и обнаружение захвата экрана
- iOS: документация (IOS.md и руководство на сайте документации)
- Android: нативные диалоги сообщений (AlertDialog) и диалоги открытия файла или файлов (Storage Access Framework, импорт с созданием копий в кэше); диалоги открытия каталога и сохранения возвращают явную ошибку
- Android: поддержка буфера обмена через ClipboardManager
- Android: реальные параметры экрана через WindowMetrics/DisplayMetrics (dp, пиксели, масштаб и рабочая область с учётом системных панелей)
- Android: методы среды выполнения для тактильной обратной связи (`Android.Haptics.Vibrate`), сведений об устройстве (`Android.Device.Info`) и всплывающих уведомлений (`Android.Toast.Show`)
- Android: типизированные события жизненного цикла (`events.Android.*`, создаваемые из events.txt), где `ActivityCreated` сопоставляется с `Common.ApplicationStarted`
- Android: конвейер сборки создаёт устанавливаемые отладочные и релизные APK (`android:run`, `android:package`, `android:package:fat`); по умолчанию релиз подписывается отладочным хранилищем ключей, а настоящее хранилище ключей можно указать через переменные среды `ANDROID_KEYSTORE_*`
- Android: `wails3 doctor` сообщает сведения об Android SDK, NDK и JDK
- Android: системные события — состояние батареи, сети, темы, блокировки экрана и нехватки памяти предоставляются как события приложения `events.Android.*` и платформонезависимые события `events.Common.*`
- Android: нативный мост мобильных функций (экспортируемый `application.Android*`) — отправка данных, открытие URL, предотвращение перехода в спящий режим, фонарик, отступы безопасной зоны, яркость, сведения о приложении, блокировка ориентации, строка состояния, биометрия (BiometricPrompt), локальные уведомления и защищённое хранилище EncryptedSharedPreferences
- Android: датчики и оборудование — тактильная обратная связь, однократное определение геопозиции, акселерометр, датчик приближения, преобразование текста в речь, сведения о хранилище, состояние питания и батареи, состояние сети, отступы для клавиатуры и блокировка захвата экрана с помощью FLAG_SECURE
- Android: документация (ANDROID.md и руководство на сайте документации)
- Пример: в демонстрационном приложении `mobile` появились вкладки Mobile и Hardware, показывающие работу нативного моста функций в iOS и Android (вкладки-капсулы переносятся на несколько строк)
- Мобильные платформы: экономия заряда батареи — акселерометр, датчик приближения, фонарик и периодические часы из примера приостанавливаются при переходе приложения в фоновый режим и возобновляются после возвращения (Android оставляет процесс работающим в фоне, а фонарик представляет собой состояние оборудования, сохраняющееся в iOS); приёмники системных событий Android регистрируются только тогда, когда приложение находится на переднем плане
- iOS: съёмка камерой — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → событие `native:capture` с миниатюрой в формате base64)
- iOS: фоновое выполнение — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (временное окно фоновой задачи UIApplication) и настраиваемый параметр `ios.backgroundModes` (build/config.yml), подставляющий `UIBackgroundModes` в создаваемый Info.plist
- Android: съёмка камерой — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (системная камера через FileProvider → событие `native:capture`)
- Android: служба переднего плана — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (`WailsForegroundService` с постоянным уведомлением сохраняет процесс активным для длительной фоновой работы)
- Пример: вкладка Camera, демонстрирующая фото- и видеосъёмку, а также фоновое выполнение (служба переднего плана в Android и временное окно фоновой задачи в iOS)

## Исправлено

- Исправлена постоянная ошибка `NotAllowedError` при вызове `getUserMedia` в Linux: WebKitGTK отклоняет запросы разрешений, если их никто не обрабатывает, а сигнал `permission-request` не был подключён. Теперь доступ к камере и микрофону управляется новой кроссплатформенной картой `WebviewWindowOptions.Permissions` (`map[PermissionType]Permission`), которая учитывается как в Linux (WebKitGTK), так и в Windows (WebView2). В Linux, где нет нативного запроса разрешения, доступ к камере и микрофону по умолчанию разрешён (что восстанавливает работу `getUserMedia`), но его можно отключить с помощью `PermissionDeny` (#5552)
- iOS: `GOOS=ios` снова компилируется (экспортирован `events.IOS`, добавлены заглушки имён мобильных методов), также компилируются сборки с тегом production (исправлены теги сборки в pkg/application и нескольких сервисах)
- iOS: события Go→JS и ExecJS теперь работают — страница больше не загружается дважды при запуске, а рукопожатие `wails:runtime:ready` больше не может быть потеряно
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` больше не конкурируют с запуском приложения; удалена фиксированная 2-секундная задержка при запуске
- iOS: устранена утечка C-строки при каждом выполнении JavaScript из Go→JS
- iOS: `hasListeners` теперь отражает фактическую регистрацию слушателей
- iOS: отладочное журналирование фреймворка больше не включается в производственные сборки
- Android: `GOOS=android` снова компилируется — определён `events.Android`, удалён выходящий за границы массив слушателей `events_android.go`, добавлена заглушка имени мобильного метода, а файлы настольной версии Linux (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) больше не попадают в сборки Android
- Android: привязки JS→Go теперь работают — WebView не может передавать тела POST-запросов `fetch()` в `shouldInterceptRequest`, поэтому вызовы среды выполнения теперь направляются через транспорт JavascriptInterface (`nativeHandleRuntimeCall`), а не завершаются сбоем из-за nil в теле запроса
- Android: вызовы среды выполнения `Screens.*` возвращают реальные данные — ScreenManager теперь заполняется при запуске (ранее он вообще не был подключён, поэтому `GetAll` возвращал nil)
- Android: в производственных сборках отладочное журналирование фреймворка исключается при компиляции, а в отладочных сборках направляется в logcat с тегом `Wails`
- Android: полноценный реестр `hasListeners`, обработка ссылок и исключений JNI, а также жизненный цикл страницы с однократной загрузкой (без двойной навигации)
- Исправлен сбой `wails3 generate bindings` с ошибкой «Access is denied» в Windows при запущенном сервере разработки Vite: теперь сгенерированные файлы синхронизируются с выходным каталогом, а не переименовываются поверх него (#5515)
- Исправлен периодически возникавший фатальный сбой в macOS при чтении сведений об экранах после изменения конфигурации дисплеев: идентификатор и имя экрана хранили указатели на автоматически освобождаемые буферы `UTF8String`, которые могли быть освобождены до их копирования в Go (use-after-free). Теперь строки копируются с помощью `strdup` и освобождаются после преобразования, а перечисление экранов выполняется в явном пуле автоматического освобождения, поэтому при вызове из горутин Go утечек больше нет (#5556)
- Исправлена периодически возникавшая ошибка SIGSEGV в Linux при закрытии сервером ресурсов объекта `WebKitURISchemeRequest`: завершающий вызов `g_object_unref` выполнялся в горутине сервера ресурсов, из-за чего GObject WebKit финализировался вне главного потока GTK. Теперь операция unref передаётся в главный контекст GTK через `g_main_context_invoke` (#5557)

## v3.0.0-alpha.100 — 2026-06-13

## Добавлено

- В `MacWebviewPreferences` добавлены дополнительные параметры конфигурации WKWebView: `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` и `ApplicationNameForUserAgent` (#5549)

## Исправлено

- Исправлен сбой `wails3 generate bindings` с ошибкой «Access is denied» в Windows при запущенном сервере разработки Vite: теперь сгенерированные файлы синхронизируются с выходным каталогом, а не переименовываются поверх него (#5561)
- Исправлено отсутствие событий изменения размера JS для окон без рамки в Linux; исправлено определение границы полосы прокрутки для окон без рамки (#5368)
- Исправлен сбой средства обновления в Windows с ошибкой «invalid cross-device link», когда временный каталог находится на другом томе, нежели каталог установки (#5560)

## v3.0.0-alpha.99 — 2026-06-10

## Исправлено

- Исправлен сбой `wails3 generate bindings` с ошибкой «Access is denied» в Windows при запущенном сервере разработки Vite: теперь сгенерированные файлы синхронизируются с выходным каталогом, а не переименовываются поверх него (#5515)

## v3.0.0-alpha.98 — 2026-06-03

## Исправлено

- Исправлено зависание интерфейса WebKit в Linux в режиме простоя (например, при открытом инспекторе): `SA_ONSTACK` больше не устанавливается принудительно для `SIGUSR1`, поскольку это нарушало синхронизацию потока сборщика мусора JavaScriptCore (#5527)

## v3.0.0-alpha.97 — 2026-05-31

## Добавлено

- Добавлена страница об отладке и работе с `runtime/trace`

## Изменено

- Удалены некоторые ненужные импорты `_ "embed"`, что немного упорядочило код

## Исправлено

- Исправлено несоблюдение ограничений минимальной ширины и высоты после выхода окна из развёрнутого состояния в Windows (#4593)
- Исправлено прохождение щелчков мыши сквозь полноэкранное окно с параметрами Frameless + Transparent (#4408)

## v3.0.0-alpha.96 — 2026-05-25

## Добавлено

- Добавлена поддержка обфускации Garble ([#4563](https://github.com/wailsapp/wails/issues/4563)): стабильные идентификаторы методов привязки, интеграция со сборкой и Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`), а также теги структур JSON во всех передаваемых среде выполнения данных (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`), благодаря чему формат передачи сохраняется при переименовании экспортируемых полей инструментом Garble.

## v3.0.0-alpha.95 — 2026-05-20

## Добавлено

- Добавлена отсутствовавшая страница о структуре проекта

## Изменено

- Документация: несколько схем на странице архитектуры заменены диаграммами последовательности для более наглядного отображения
- Документация: добавлено примечание об установке D2 как необходимого условия для запуска

## Исправлено

- Исправлена работа `wails3 generate appimage` со стандартным GTK4: теперь упаковщик определяет стек GTK по исполняемому файлу до поиска файлов среды выполнения, поэтому для сборок GTK4 выбирается `libwebkitgtkinjectedbundle.so` (в `webkitgtk-6.0/`), а для сборок `-tags gtk3` — `libwebkit2gtkinjectedbundle.so` (в `webkit2gtk-4.1/`). Проверка `.relr.dyn` теперь также учитывает `libgtk-4.so.1`, поэтому в современных наборах инструментов удаление символов корректно отключается независимо от стека. (#5475)
- Исправлен сбой `wails3 generate appimage` при вызове с относительным путём `-builddir`: теперь упаковщик заранее преобразует `-binary`, `-icon`, `-desktopfile`, `-builddir` и `-outputdir` в абсолютные пути, чтобы выполняемый в середине процесса `s.CD` не нарушал работу горутины загрузки AppRun или выполняемую после копирования проверку `ldd`.
- Исправлена ошибка `wails3 generate appimage` при перемещении итогового AppImage в `-outputdir`, когда поле `Name=` файла desktop не совпадает с базовым именем исполняемого файла: теперь упаковщик через переменную окружения `OUTPUT` принудительно указывает плагину appimage из linuxdeploy записывать AppImage в `<binary>-<arch>.AppImage`, а не под именем, полученным из файла desktop.
- Исправлено отсутствие событий `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` и `Common.SystemDidWake` в Linux после того, как стек GTK4 + WebKitGTK 6.0 стал стандартным в alpha.93. Новый стандартный `application_linux.go` `run()` не вызывал `setupCommonEvents()` (который перенаправляет события `Linux.*` соответствующим обработчикам `Common.*`) и `monitorPowerEvents()`. Теперь вспомогательный модуль отслеживания питания через DBus используется совместно путями сборки GTK3 и GTK4 посредством `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.94 — 2026-05-19

## Исправлено

- Исправлено отсутствие событий `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` и `Common.SystemDidWake` в Linux после того, как стек GTK4 + WebKitGTK 6.0 стал стандартным в alpha.93. Новый стандартный `application_linux.go` `run()` не вызывал `setupCommonEvents()` (который перенаправляет события `Linux.*` соответствующим обработчикам `Common.*`) и `monitorPowerEvents()`. Теперь вспомогательный модуль отслеживания питания через DBus используется совместно путями сборки GTK3 и GTK4 посредством `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.93 — 2026-05-17

## Добавлено

- В вывод `wails3 doctor` в Linux добавлен `XDG_SESSION_TYPE`, автор — @leaanthony

## Исправлено

- Исправлен вызванный обращением appmenu-gtk-module к нереализованному окну сбой меню окна в Wayland (#4769), автор — @leaanthony
- Исправлен сбой приложения GTK, когда имя приложения содержит недопустимые символы (пробелы, круглые скобки и т. п.), автор — @leaanthony
- Исправлена ошибка «not enough memory» при инициализации перетаскивания в Windows (#4701), автор — @overlordtm
- Исправлено состояние гонки в хранилище обратных вызовов главного потока, вызванное ошибочным использованием RLock при удалении из map (Linux, macOS, iOS) (#4424), автор — @leaanthony
- Исправлена обработка переменных при передаче задачам аргументов командной строки. Теперь переменные CLI, заданные парами KEY=VALUE, правильно инициализируются и передаются на всех этапах выполнения задачи.
- Устранён конфликт NSWindowZoomButton в macOS: теперь `MaximiseButtonState` и `FullscreenButtonState` применяют более строгое состояние как при запуске, так и во время выполнения; ни один из сеттеров больше не может незаметно переопределить другой (#5319)
- Исправлена группа ранее существовавших ошибок в устаревшем пути сборки GTK3 (`-tags gtk3`), обнаруженных CodeRabbit в #5463: при запуске через ассоциацию файлов обработчики запуска больше не пропускаются; `getTheme` теперь безопасен с точки зрения границ и типов; `appName` больше не освобождает память, принадлежащую GLib; `clipboardGet` больше не допускает утечку `gchar*`, возвращаемого GTK; `Calloc` теперь использует получатели-указатели (а `NewCalloc` возвращает `*Calloc`), поэтому пул действительно отслеживает выделение памяти; `zoomOut` использует величину, обратную `zoomInFactor`, вместо отрицательного множителя, из-за которого значение ограничивалось до 1.0; `execJS` повторно использует заранее выделенное пустое имя мира вместо утечки одного `C.CString("")` при каждом вызове; из `menuItem.setAccelerator` удалён отладочный `fmt.Println`. Устраняет #5465.
- Исправлена такая же утечка из-за получателя по значению `Calloc` в стандартном пути сборки GTK4 (`linux_cgo.go`): благодаря получателям-указателям и `NewCalloc() *Calloc` выделения `c.String(...)` для каждого окна теперь действительно отслеживаются и освобождаются.

## v3.0.0-alpha.92 — 2026-05-15

## Добавлено

- Taskfile изменены, чтобы используемым менеджером пакетов фронтенда можно было управлять с помощью параметра `PACKAGE_MANAGER`
- Данные шаблонов дополнены `{{.Opn}}` и `{{.Cls}}`, чтобы поведение шаблонов Taskfile при написании было более предсказуемым

## Изменено

- Несколько существующих Taskfile изменены для использования `{{.Opn}} and {{.Cls}}`

## Исправлено

- Исправлена фатальная ошибка среды выполнения `concurrent map read and map write` в `linuxSystemTray`, возникавшая при обновлении меню в области уведомлений во время его чтения панелью.
- Для вывода ошибок и трассировки стека WebView2 теперь используется `log` вместо `fmt`, поэтому сообщения не теряются, когда приложение работает в Windows без подключённой консоли.

## v3.0.0-alpha.91 — 2026-05-12

## Изменено

- SVG-файл спонсоров обновлён в [PR](https://github.com/wailsapp/wails/pull/5414) автором `@github-actions[bot]`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ (macOS):** система координат macOS унифицирована, и теперь `GetScreens`, `Position` и `SetPosition` используют единое пространство: логические точки, ось Y направлена вниз, а `(0,0)` находится в левом верхнем углу основного экрана. Это соответствует Windows, GTK, а также общедоступным API Electron и веб-платформы. Для экранов, физически расположенных выше основного, `Bounds.Y` теперь имеет отрицательное значение (ранее положительное), а значения `Position()`/`SetPosition()` теперь выражаются в логических точках вместо `points × primaryScale`. Преобразование туда и обратно `Position()` → `SetPosition()` сохранено; абсолютные значения, записанные в журналах предыдущих альфа-сборок, и обходные решения с вычислениями вручную (например, умножение на `primaryScale` или инвертирование Y относительно высоты экрана) потребуется обновить. Устраняет [#5117](https://github.com/wailsapp/wails/issues/5117).

## Исправлено

- Добавлена защитная проверка имени сигнала DBus и длины его тела для предотвращения паник в [PR](https://github.com/wailsapp/wails/pull/5416) от @leaanthony
- Исправлена проблема безопасности памяти при обработке меню GTK в Linux в [PR](https://github.com/wailsapp/wails/pull/5363) от @leaanthony
- Добавлено обнаружение графических процессоров NVIDIA и отключение средства визуализации DMA-BUF в Linux в [PR](https://github.com/wailsapp/wails/pull/5295) от @leaanthony
- Исправлено преобразование Y между экранами в `SetPosition` на macOS: высота основного экрана используется как глобальная точка отсчёта, чтобы окна размещались в правильной позиции на мониторах, смещённых по вертикали относительно основного дисплея, в [#5117](https://github.com/wailsapp/wails/issues/5117)
- Исправлен шаблон PR для Git: теперь он указывает на правильный URL для обратной связи, в [PR](https://github.com/wailsapp/wails/pull/5109) от @wayneforrest
- Исправлена группа аварийных завершений `SetMenu` в области уведомлений Windows, вызванных некорректным системным вызовом `DestroyMenu`, который передавал четыре аргумента вместо одного, поэтому каждый вызов возвращал FALSE и ничего не освобождал. Кроме того, при перестроении меню теперь освобождаются дескрипторы HMENU и HBITMAP, включая выделенные во время выполнения через `MenuItem.SetBitmap`; в `Win32Menu.Update` сбрасываются устаревшие карты флажков и переключателей; из `systemtray.updateMenu` удалён избыточный вызов `Update()`, удваивавший количество выделений памяти. Долго работающие приложения с областью уведомлений больше не допускают утечку объектов GDI/USER при каждом перестроении меню.

## v3.0.0-alpha.90 — 2026-05-11

## Добавлено

- Добавлена возможность настройки имени приложения в User-Agent WKWebView на macOS в [PR](https://github.com/wailsapp/wails/pull/5261) от @vinhvoit225
- В пример gin-service добавлена косвенная зависимость github.com/coder/websocket в [PR](https://github.com/wailsapp/wails/pull/5400) от @taliesin-ai
- В тесты ресурсов сборки добавлена поддержка глубокого сравнения на равенство в [PR](https://github.com/wailsapp/wails/pull/5402) от @leaanthony

## Изменено

- Выходные данные сборки объединены в каталоге assets в [PR](https://github.com/wailsapp/wails/pull/5401) от @taliesin-ai
- SVG-файл спонсоров обновлён в [PR](https://github.com/wailsapp/wails/pull/5399) автором `@github-actions[bot]`

## Исправлено

- Для сообщения одного экземпляра приложения в macOS теперь используется объект уведомления в [PR](https://github.com/wailsapp/wails/pull/5289) от @overlordtm
- Обратные вызовы Windows объединены в пакеты, чтобы предотвратить потерю промисов при высокой нагрузке, в [PR](https://github.com/wailsapp/wails/pull/5383) от @taliesin-ai

## v3.0.0-alpha.89 — 2026-05-10

## Добавлено

- Добавлено задание go<em>test</em>results для объединения результатов тестов Go в [PR](https://github.com/wailsapp/wails/pull/5316) от @leaanthony

## Изменено

- Крупные полезные данные RPC при необходимости разделяются на фрагментированные POST-запросы в [PR](https://github.com/wailsapp/wails/pull/5369) от @leaanthony
- Vite обновлён с 5.x.x до 8.0.0 во всех шаблонах фронтенда в [PR](https://github.com/wailsapp/wails/pull/5386) от @leaanthony
- Конфигурация порта сервера разработки Vite перенесена в переменные среды в [PR](https://github.com/wailsapp/wails/pull/5365) от @leaanthony
- Во всех шаблонах сервер разработки Vite настроен для привязки к 127.0.0.1 в [PR](https://github.com/wailsapp/wails/pull/5361) от @leaanthony
- SVG-файл спонсоров обновлён в [PR](https://github.com/wailsapp/wails/pull/5384) автором `@github-actions[bot]`

## Исправлено

- Заготовки шаблона Info.plist очищаются при обновлении build-assets в [PR](https://github.com/wailsapp/wails/pull/5312) от @leaanthony
- Исправлено устаревшее состояние меню macOS: изменяющие операции элементов меню (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) теперь синхронно выполняются в главном потоке, что устраняет состояние гонки `dispatch_async`, из-за которого при быстром повторном открытии меню отображали предыдущее состояние (#5002)
- В режиме разработки файлы `*_test.go` игнорируются для предотвращения ненужных пересборок в [PR](https://github.com/wailsapp/wails/pull/5203) от @leaanthony
- Предотвращена ошибка сегментации в Menu.Update(), когда приложение не запущено, в [PR](https://github.com/wailsapp/wails/pull/5291) от @wucm667
- Для ограничения перерисовки строки меню в Windows теперь используется lastSizeWParam в [PR](https://github.com/wailsapp/wails/pull/5382) от @taliesin-ai

## v3.0.0-alpha.88 — 2026-05-09

## Изменено

- HiddenOnTaskbar переведён на использование WS<em>EX</em>TOOLWINDOW в [PR](https://github.com/wailsapp/wails/pull/5371) от @leaanthony
- Зависимости переупорядочены, а директива replace для webview2 удалена из go.mod в [PR](https://github.com/wailsapp/wails/pull/5370) от @atterpac
- SVG-файл спонсоров обновлён в [PR](https://github.com/wailsapp/wails/pull/5358) автором `@github-actions[bot]`

## Исправлено

- Удалены универсальные псевдонимы косвенности и унифицированы типы ключей карт в [PR](https://github.com/wailsapp/wails/pull/5331) от @fbbdev

## Удалено

- Удалён рабочий процесс PR-master вместе с документацией, тестами Go и пропуском тестов в [PR](https://github.com/wailsapp/wails/pull/5377), автор: @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## Добавлено

- Добавлена документация Wails v3 на корейском языке в [PR](https://github.com/wailsapp/wails/pull/5352), автор: @leaanthony
- Добавлена документация на французском языке по установке и быстрому старту в [PR](https://github.com/wailsapp/wails/pull/5354), автор: @leaanthony
- Добавлена документация на португальском языке по быстрому старту, основным понятиям и сообществу в [PR](https://github.com/wailsapp/wails/pull/5355), автор: @leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## Добавлено

- Добавлена локализация документации на французский язык в [PR](https://github.com/wailsapp/wails/pull/5328), автор: @leaanthony
- На сайт документации добавлена немецкая локаль в [PR](https://github.com/wailsapp/wails/pull/5343), автор: @leaanthony

## Изменено

- Все 8 локализованные версии зарегистрированы в конфигурации документации в [PR](https://github.com/wailsapp/wails/pull/5347), автор: @leaanthony
- Обновлены различные связанные с Windows файлы для WebView2 в [PR](https://github.com/wailsapp/wails/pull/5317), автор: @leaanthony

## Исправлено

- Диспетчеризация диалогов в Linux разделена между GTK3 и GTK4 в [PR](https://github.com/wailsapp/wails/pull/5340), автор: @leaanthony
- Обеспечено выполнение обратных вызовов диалогов в потоке GTK, что устраняет ошибки сегментации, в [PR](https://github.com/wailsapp/wails/pull/5339), автор: @leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## Добавлено

- В репозиторий добавлен URL шаблона PR в [PR](https://github.com/wailsapp/wails/pull/5179), автор: @leaanthony
- Добавлена документация Wails v3 на немецком языке в [PR](https://github.com/wailsapp/wails/pull/5330), автор: @leaanthony

## v3.0.0-alpha.84 - 2026-05-03

## Добавлено

- Добавлена возможность запретить выход из полноэкранного режима по клавише Escape в macOS в [PR](https://github.com/wailsapp/wails/pull/5307), автор: @leaanthony
- Добавлена возможность запретить выход из полноэкранного режима по клавише Escape в macOS в [PR](https://github.com/wailsapp/wails/pull/5310), автор: @leaanthony
- Добавлена документация о проекте Pausa в разделе работ сообщества в [PR](https://github.com/wailsapp/wails/pull/5288), автор: @yuseferi

## Изменено

- Обновлён SVG-файл спонсоров в [PR](https://github.com/wailsapp/wails/pull/5308), автор: `@github-actions[bot]`
- Команда создания значков обновлена для обработки неподдерживаемых платформ в [PR](https://github.com/wailsapp/wails/pull/5309), автор: @leaanthony
- Логический API полноэкранного режима заменён трёхсостояниевым ButtonState; реализованы привязки для платформ в [PR](https://github.com/wailsapp/wails/pull/5224), автор: @leaanthony

## Исправлено

- Операции фокусировки WebView2 защищены от состояния, в котором контроллер равен nil, в [PR](https://github.com/wailsapp/wails/pull/5315), автор: @leaanthony
- Рабочий процесс GitHub Actions обновлён так, чтобы правильно ссылаться на базовую ветку PR, в [PR](https://github.com/wailsapp/wails/pull/5313), автор: @leaanthony
- Файлы `*_test.go` игнорируются в режиме разработки во избежание ненужных пересборок в [PR](https://github.com/wailsapp/wails/pull/5203), автор: @leaanthony
- Предотвращена ошибка сегментации в Menu.Update(), когда приложение не запущено, в [PR](https://github.com/wailsapp/wails/pull/5291), автор: @wucm667

## v3.0.0-alpha.83 - 2026-05-02

## Добавлено

- Добавлены флаг InstallScope и параметр сборки для установки на уровне компьютера или пользователя в [PR](https://github.com/wailsapp/wails/pull/5094), автор: @symball
- В BrowserWindow добавлен метод-заглушка SetScreen для соответствия интерфейсу Window в [PR](https://github.com/wailsapp/wails/pull/5294), автор: @leaanthony

## Исправлено

- Добавлено обнаружение графических процессоров NVIDIA и отключение средства визуализации DMA-BUF в Linux в [PR](https://github.com/wailsapp/wails/pull/5295), автор: @leaanthony
- Исправлен шаблон PR Git: теперь он указывает правильный URL для обратной связи, в [PR](https://github.com/wailsapp/wails/pull/5109), автор: @wayneforrest
- Исправлена группа сбоев `SetMenu` системного трея Windows, вызванных некорректным системным вызовом `DestroyMenu`, который передавал четыре аргумента вместо одного, поэтому каждый вызов возвращал FALSE и ничего не освобождал. Кроме того, при пересоздании меню теперь освобождаются дескрипторы HMENU и HBITMAP, включая выделенные во время выполнения через `MenuItem.SetBitmap`; в `Win32Menu.Update` сбрасываются устаревшие карты флажков и переключателей; из `systemtray.updateMenu` удалён избыточный вызов `Update()`, удваивавший количество выделений памяти. Длительно работающие приложения с системным треем больше не допускают утечки объектов GDI/USER при каждом пересоздании меню.

## v3.0.0-alpha.82 - 2026-05-01

## Исправлено

- Исправлено создание desktop-файла: теперь имя рабочего стола обрабатывается правильно, в [PR](https://github.com/wailsapp/wails/pull/5232), автор: @leaanthony

## v3.0.0-alpha.81 - 2026-04-30

## Изменено

- Расписание ночных выпусков перенесено на 15:00 UTC в [PR](https://github.com/wailsapp/wails/pull/5286), автор: @leaanthony

## Исправлено

- Исправлены вдвое уменьшенные значения Screen Bounds, WorkArea и Size на компьютерах Mac с дисплеями Retina —  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## Изменено

- Обновлены зависимости документации и загрузчики коллекций содержимого в [PR](https://github.com/wailsapp/wails/pull/5285), автор: @leaanthony

## v3.0.0-alpha.79 - 2026-04-29

## Добавлено

- Заданию trigger-release предоставлено разрешение actions: write в [PR](https://github.com/wailsapp/wails/pull/5270), автор: @leaanthony

## Изменено

- Задача выпуска теперь по умолчанию использует ветку master; также обновлена формулировка журнала изменений в [PR](https://github.com/wailsapp/wails/pull/5283), автор: @leaanthony
- Рабочий процесс auto-changelog обновлён для использования последней версии в [PR](https://github.com/wailsapp/wails/pull/5282), автор: @leaanthony
- Повышена эффективность рабочих процессов благодаря добавлению фильтров путей и удалению неиспользуемых рабочих процессов в [PR](https://github.com/wailsapp/wails/pull/5280), автор: @leaanthony
- Документация обновлена: ссылки на примеры теперь указывают на ветку master, в [PR](https://github.com/wailsapp/wails/pull/5274), автор: @leaanthony
- Обновлены документация и примеры для v3 в [PR](https://github.com/wailsapp/wails/pull/5272), автор: @leaanthony

## Исправлено

- В обратный прокси-сервер добавлены повторные попытки и принудительное использование IPv4 при разработке в [PR](https://github.com/wailsapp/wails/pull/5265), автор: @AkagiYui
- Переработан рабочий процесс запуска создания журнала невыпущенных изменений в [PR](https://github.com/wailsapp/wails/pull/5281), автор: @leaanthony

## Удалено

- Удалены сценарии тестирования оболочки для различных целей тестирования в [PR](https://github.com/wailsapp/wails/pull/5267), автор: @leaanthony
- Удалены рабочий процесс развёртывания документации v3-alpha и запись CNAME в [PR](https://github.com/wailsapp/wails/pull/5266), автор: @leaanthony

### Добавлено

- В боковую панель навигации добавлен пункт «Маршрутизация фронтенда» в [PR](https://github.com/wailsapp/wails/pull/5196), автор: @leaanthony
- Добавлено руководство по маршрутизации фронтенда с рекомендациями для конкретных фреймворков в [PR](https://github.com/wailsapp/wails/pull/5185), автор: @leaanthony
- Добавлена поддержка модальных листов (macOS)
- Версия ghw обновлена для улучшения поддержки устройств Apple; автор: @leaanthony (#4977)
- В сервис Dock добавлен метод `GetBadge`
- В команду `wails3 build` добавлен флаг `-tags` для передачи пользовательских тегов сборки Go (например, `wails3 build -tags gtk4`) (#4957)
- Добавлена документация об автоматическом создании перечислений генератором привязок, включая отдельную страницу Enums и пункт в боковой панели навигации (#4972)
- В команду `wails3 build` добавлен флаг `-tags` для передачи пользовательских тегов сборки Go (например, `wails3 build -tags gtk4`) (#4957)
- В `v3/examples/web-apis/` добавлены примеры использования Web API, демонстрирующие браузерные API 41, включая хранилища (localStorage, sessionStorage, IndexedDB, Cache API), сетевые API (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), мультимедийные API (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), API устройств (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), API производительности (Performance API, Mutation Observer, Intersection/Resize Observer), API пользовательского интерфейса (Web Components, Pointer Events, Selection, Dialog, Drag and Drop) и другие
- Добавлен пример средства проверки совместимости WebView API (`v3/examples/webview-api-check/`), тестирующего более 200 браузерных API на разных платформах
- Добавлен пакет `internal/libpath` для поиска путей к нативным библиотекам в Linux с параллельным поиском, кешированием и поддержкой Flatpak/Snap/Nix
- **В разработке:** добавлена экспериментальная поддержка WebKitGTK 6.0 / GTK4 для Linux, доступная через `-tags gtk4` (по умолчанию по-прежнему используется GTK3/WebKit2GTK 4.1)
- Примечание. В мозаичных оконных менеджерах (например, Hyprland и Sway) операции сворачивания и разворачивания могут работать не так, как ожидается, поскольку геометрией окон управляет оконный менеджер
- В документацию по **прослушиванию событий в JavaScript** добавлено описание **одноразовых обработчиков**; автор: @AbdelhadiSeddar
- В `WebviewWindowOptions` добавлен параметр `UseApplicationMenu`, позволяющий окнам в Windows/Linux наследовать меню приложения, заданное через `app.Menu.Set()`; автор: @leaanthony
- Добавлена поддержка файлов `.icon` (формат Apple Icon Composer) для создания значков Liquid Glass и каталогов ресурсов (macOS) (#4934); автор: @wimaha
- Добавлен экспериментальный серверный режим для развёртываний без графического интерфейса и веб-развёртываний (`-tags server`). Он позволяет запускать приложения Wails как HTTP-серверы без зависимостей от нативного графического интерфейса. Выполните сборку с `wails3 task build:server`. Подробности см. в `examples/server`.
- Добавлен пакет `internal/libpath` для поиска путей к нативным библиотекам в Linux с параллельным поиском, кешированием и поддержкой Flatpak/Snap/Nix
- В `MacWindow` добавлен параметр `CollectionBehavior` для управления поведением окна в пространствах macOS Spaces и полноэкранном режиме (#4756); автор: @leaanthony
- Добавлены модульные тесты для pkg/application; автор: @leaanthony
- В упаковку MSIX добавлена поддержка пользовательских протоколов; автор: @leaanthony
- Добавлено определение среды рабочего стола в Linux, [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- В среду выполнения JavaScript добавлен метод `Window.Print()` для вызова диалогового окна печати из фронтенда (#4290); автор: @leaanthony
- В вывод `wails3 doctor` в Linux добавлен `XDG_SESSION_TYPE`; автор: @leaanthony
- Добавлены дополнительные события изменения состояния загрузки WebKit2 для Linux: `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` (#3896); автор: @leaanthony
- В вывод `wails3 doctor` в Linux добавлен `XDG_SESSION_TYPE`; автор: @leaanthony
- Файл `.desktop` теперь создаётся не только при упаковке, но и во время сборки для Linux (#4575)
- Добавлена документация о зависимостях среды выполнения для Linux с именами пакетов для конкретных дистрибутивов и примерами упаковки с помощью nfpm (#4339); автор: @leaanthony
- В вывод `wails3 doctor` в Linux добавлены сведения о версии драйвера NVIDIA; автор: @leaanthony
- В обработчик необработанных сообщений добавлен источник; автор: @APshenkin, [PR](https://github.com/wailsapp/wails/pull/4710)
- Добавлена поддержка универсальных ссылок для macOS; автор: @APshenkin, [PR](https://github.com/wailsapp/wails/pull/4712)
- Переработан транспортный слой привязок; автор: @APshenkin, [PR](https://github.com/wailsapp/wails/pull/4702)
- В шаблоны helloworld добавлены идентификаторы aria-label, чтобы клиенты тестирования Appium могли легко тестировать пример приложения; автор: @chinenual, [PR](https://github.com/wailsapp/wails/pull/4760)
- В обработчик необработанных сообщений добавлен источник; автор: @APshenkin, [PR](https://github.com/wailsapp/wails/pull/4710)
- Добавлена поддержка универсальных ссылок для macOS; автор: @APshenkin, [PR](https://github.com/wailsapp/wails/pull/4712)
- Переработан транспортный слой привязок; автор: @APshenkin, [PR](https://github.com/wailsapp/wails/pull/4702)
- Типизированные события; авторы: @fbbdev и @ianvs, [#4633](https://github.com/wailsapp/wails/pull/4633)
- Добавлен пример `systray-clock`, демонстрирующий панельный значок без основного окна с динамическим обновлением всплывающей подсказки (#4653).
- Добавлен шаблон протокола NSIS для Windows; автор: @Tolfx, #4510
- Добавлены тесты для build-assets; автор: @Tolfx, #4510
- macOS: нативные элементы управления окном теперь отображаются в строке меню; [#4588](https://github.com/wailsapp/wails/pull/4588), автор: @nidib
- Добавлен сервис Dock для macOS, позволяющий скрывать и отображать значок приложения в Dock; автор: @popaprozac, [PR](https://github.com/wailsapp/wails/pull/4451)
- Добавлен сервис Dock для macOS, позволяющий скрывать и отображать значок приложения в Dock; автор: @popaprozac, [PR](https://github.com/wailsapp/wails/pull/4451)
- Добавлена поддержка нативного эффекта Liquid Glass для macOS с использованием NSGlassEffectView (macOS 15.0+) и резервным использованием NSVisualEffectView, включая широкие возможности настройки материала; автор: @leaanthony, [#4534](https://github.com/wailsapp/wails/pull/4534)
- Санитизация URL браузера; автор: @leaanthony, [#4500](https://github.dev/wailsapp/wails/pull/4500). Основано на [#4484](https://github.com/wailsapp/wails/pull/4484), автор: @APShenkin.
- Добавлена защита содержимого в Windows/Mac; автор: [@leaanthony](https://github.com/leaanthony), на основе первоначальной работы [@Taiterbase](https://github.com/Taiterbase) в этом [PR](https://github.com/wailsapp/wails/pull/4241)
- Добавлена поддержка передачи переменных CLI командам Task через псевдонимы `wails3 build` и `wails3 package` (#4422); автор: @leaanthony, [PR](https://github.com/wailsapp/wails/pull/4488)
- Добавлена поддержка зон перетаскивания, в которых данные о перетаскиваемом элементе поступают из события; автор: [@atterpac](https://github.com/atterpac), [#4318](https://github.com/wailsapp/wails/pull/4318)
- В параметры `WindowsWindow` добавлен `AdditionalLaunchArgs`, позволяющий передавать браузеру WebView2 дополнительные аргументы командной строки; [PR](https://github.com/wailsapp/wails/pull/4467)
- Добавлен автоматический запуск go mod tidy после wails init; автор: [@triadmoko](https://github.com/triadmoko), [PR](https://github.com/wailsapp/wails/pull/4286)
- Функция Windows Snap Assist; автор: @leaanthony, [PR](https://github.dev/wailsapp/wails/pull/4463)
- В параметры `WindowsWindow` добавлен `AdditionalLaunchArgs`, позволяющий передавать браузеру WebView2 дополнительные аргументы командной строки; [PR](https://github.com/wailsapp/wails/pull/4467)
- Добавлен автоматический запуск go mod tidy после wails init; автор: [@triadmoko](https://github.com/triadmoko), [PR](https://github.com/wailsapp/wails/pull/4286)
- Функция Windows Snap Assist; автор: @leaanthony, [PR](https://github.dev/wailsapp/wails/pull/4463)
- Добавлена реализация `getAccentColor` для Windows; автор: [@almas-x](https://github.com/almas-x), [PR](https://github.com/wailsapp/wails/pull/4427)
- Добавлена реализация `getAccentColor` для Windows; автор: [@almas-x](https://github.com/almas-x), [PR](https://github.com/wailsapp/wails/pull/4427)
- Меню и строка меню с тёмной темой в Windows. Автор: @leaanthony, [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Переименование встроенных сервисов для более понятных привязок JS/TS — @popaprozac, [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` для получения системного акцентного цвета пользователя. Работает в MacOS. Автор: [@etesam913](https://github.com/etesam913)
- Добавление API `window.ToggleFrameless()` — [@atterpac](https://github.com/atterpac), [#4137](https://github.com/wailsapp/wails/pull/4137)
- Добавление зависимостей сборки для конкретных дистрибутивов Linux — @leaanthony, [PR](https://github.com/wailsapp/wails/pull/4345)
- Добавлено руководство по привязкам — @atterpac, [PR](https://github.com/wailsapp/wails/pull/4404)
- **Упорядоченная инфраструктура тестирования**: тестовые файлы Docker перенесены в отдельный каталог `test/docker/`, образы оптимизированы, а надёжность сборки повышена — [@leaanthony](https://github.com/leaanthony), [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Улучшенные шаблоны управления ресурсами**: в примеры добавлены корректное удаление обработчиков событий и управление горутинами с учётом контекста — [@leaanthony](https://github.com/leaanthony), [#4359](https://github.com/wailsapp/wails/pull/4359)
- Поддержка сборки AppImage для aarch64 — [@AkshayKalose](https://github.com/AkshayKalose), [#3981](https://github.com/wailsapp/wails/pull/3981)
- Добавление раздела диагностики в `wails doctor` — [@leaanthony](https://github.com/leaanthony)
- Добавление окна в контекст при вызове метода сервиса — [@leaanthony](https://github.com/leaanthony)
- Добавление примера `window-call`, демонстрирующего, как определить, какое окно вызывает сервис, — [@leaanthony](https://github.com/leaanthony)
- Новое руководство по меню — [@leaanthony](https://github.com/leaanthony)
- Улучшенная обработка паник — [@leaanthony](https://github.com/leaanthony)
- Новое руководство по меню — [@leaanthony](https://github.com/leaanthony)
- Добавление документирующих комментариев для Service API — [@fbbdev](https://github.com/fbbdev), [#4024](https://github.com/wailsapp/wails/pull/4024)
- Добавление функции `application.NewServiceWithOptions` для инициализации сервисов с дополнительной конфигурацией — [@leaanthony](https://github.com/leaanthony), [#4024](https://github.com/wailsapp/wails/pull/4024)
- Улучшенное управление меню — [@FalcoG](https://github.com/FalcoG) и [@leaanthony](https://github.com/leaanthony), [#4031](https://github.com/wailsapp/wails/pull/4031)
- Дополнительная документация — [@leaanthony](https://github.com/leaanthony)
- Поддержка отмены событий в стандартных обработчиках событий — [@leaanthony](https://github.com/leaanthony)
- Поддержка `Hide`, `Show` и `Destroy` для системного трея — [@leaanthony](https://github.com/leaanthony)
- Поддержка `SetTooltip` для системного трея — [@leaanthony](https://github.com/leaanthony). Автор исходной идеи: [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- Вывод пути пакета в предупреждениях генератора привязок о неподдерживаемых типах — [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- Добавление в генератор привязок поддержки обобщённых псевдонимов — [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- Добавление в генератор привязок поддержки флага JSON `omitzero` — [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- Добавление директивы `//wails:ignore` для предотвращения генерации привязок для выбранных методов сервисов — [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- Добавление директивы `//wails:internal` для сервисов и моделей, позволяющей использовать типы, которые экспортируются в Go, но не в JS/TS, — [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- Добавление в генератор привязок поддержки констант с типом-псевдонимом для создания слабо типизированных перечислений — [@fbbdev](https://github.com/fbbdev), [#4045](https://github.com/wailsapp/wails/pull/4045)
- Добавление тестов генератора привязок для возможностей Go 1.24 — [@fbbdev](https://github.com/fbbdev), [#4068](https://github.com/wailsapp/wails/pull/4068)
- Добавление поддержки macOS 15 «Sequoia» в `OSInfo.Branding` для более точного определения версии ОС — [#4065](https://github.com/wailsapp/wails/pull/4065)
- Добавление хука `PostShutdown` для выполнения пользовательского кода после завершения процесса остановки — [@fbbdev](https://github.com/fbbdev), [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавление структуры `FatalError` для обнаружения фатальных ошибок в пользовательских обработчиках ошибок — [@fbbdev](https://github.com/fbbdev), [#4066](https://github.com/wailsapp/wails/pull/4066)
- Стандартизация и документирование порядка запуска и остановки сервисов — [@fbbdev](https://github.com/fbbdev), [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавление тестовой среды для последовательности запуска и остановки приложения, а также тестов запуска и остановки сервисов — [@fbbdev](https://github.com/fbbdev), [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавление метода `RegisterService` для регистрации сервисов после создания приложения — [@fbbdev](https://github.com/fbbdev), [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавление поля `MarshalError` в параметры приложения и сервисов для пользовательской обработки ошибок при вызовах привязок — [@fbbdev](https://github.com/fbbdev), [#4066](https://github.com/wailsapp/wails/pull/4066)
- Добавление отменяемой обёртки промиса, передающей запросы на отмену по цепочкам промисов, — [@fbbdev](https://github.com/fbbdev), [#4100](https://github.com/wailsapp/wails/pull/4100)
- Добавление возможности связать отмену вызова привязки с `AbortSignal` — [@fbbdev](https://github.com/fbbdev), [#4100](https://github.com/wailsapp/wails/pull/4100)
- Поддержка атрибутов `data-wml-*` для WML наряду с обычными атрибутами `wml-*` — [@leaanthony](https://github.com/leaanthony)
- Добавление метода `Configure` во все сервисы для поздней или динамической перенастройки — [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ненастроенный сервис `fileserver` отправляет ответ 503 Service Unavailable — [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- Ненастроенный сервис `kvstore` по умолчанию предоставляет хранилище пар «ключ — значение» в памяти — [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- Добавление в сервис `kvstore` метода `Load` для повторной загрузки данных из файла после изменения конфигурации — [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- Добавление в сервис `kvstore` метода `Clear` для удаления всех ключей — [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- Добавление в сервис `log` типа `Level`, предоставляющего константы уровней журналирования на стороне JS, — [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- Добавление в сервис `log` метода `Log` для динамического задания уровня журналирования — [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- Сервис `sqlite` без настройки по умолчанию предоставляет базу данных в памяти — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Добавлен метод `Close` в сервис `sqlite` для ручного закрытия базы данных — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Добавлена поддержка отмены для методов запросов сервиса `sqlite` — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- В сервис `sqlite` добавлена поддержка подготовленных выражений с привязками JS — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Поддержка Gin от [Lea Anthony](https://github.com/leaanthony) в [PR](https://github.com/wailsapp/wails/pull/3537) на основе исходной работы [@AnalogJ](https://github.com/AnalogJ) в этом [PR](https://github.com/wailsapp/wails/pull/3537)
- Исправлена постоянно включённая функция автоматического сохранения и автоматического сохранения паролей — [@oSethoum](https://github.com/osethoum) в [#4134](https://github.com/wailsapp/wails/pull/4134)
- Для окна добавлен `SetMenu()`, позволяющий назначать ему меню, — [@leaanthony](https://github.com/leaanthony)
- Добавлена поддержка уведомлений — [@popaprozac](https://github.com/popaprozac) в [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Добавлена поддержка ассоциаций файлов для mac — [@wimaha](https://github.com/wimaha) в [#4177](https://github.com/wailsapp/wails/pull/4177)
- Добавлен `wails3 tool version` для повышения семантической версии — [@leaanthony](https://github.com/leaanthony)
- Добавлена поддержка значков-индикаторов для macOS и Windows — [@popaprozac](https://github.com/popaprozac) в [#](https://github.com/wailsapp/wails/pull/4234)
- Добавлена поддержка зарегистрированных событий со строгой типизацией — [@fbbdev](https://github.com/fbbdev) и [@IanVS](https://github.com/IanVS) в [#4161](https://github.com/wailsapp/wails/pull/4161)
- Добавлена возможность регистрировать обработчики-перехватчики пользовательских событий — [@fbbdev](https://github.com/fbbdev) и [@IanVS](https://github.com/IanVS) в [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` открывает системный файловый менеджер по пути `path` с необязательным выделением через `selectFile` — [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Новый флаг `-git` для команды `wails3 init` — [@leaanthony](https://github.com/leaanthony)
- Новая команда `wails3 generate webview2bootstrapper` — [@leaanthony](https://github.com/leaanthony)
- В среду выполнения добавлен метод `init()` для её ручной инициализации — [@leaanthony](https://github.com/leaanthony)
- В WindowOptions окна добавлен параметр `WindowDidMoveDebounceMS` — [@leaanthony](https://github.com/leaanthony)
- Добавлена функция единственного экземпляра приложения — [@leaanthony](https://github.com/leaanthony). Основана на [v2 PR](https://github.com/wailsapp/wails/pull/2951) от @APshenkin.
- Команда `wails3 generate template` — [@leaanthony](https://github.com/leaanthony)
- Команда `wails3 releasenotes` — [@leaanthony](https://github.com/leaanthony)
- Команда `wails3 update cli` — [@leaanthony](https://github.com/leaanthony)
- Параметр `-clean` для команды `wails3 generate bindings` — [@leaanthony](https://github.com/leaanthony)
- Разрешена сборка AppImage для Linux на aarch64 (arm64) — [@AkshayKalose](https://github.com/AkshayKalose) в [#3981](https://github.com/wailsapp/wails/pull/3981)
- Добавлена гиперссылка для спонсора — @ansxuman в [#3958](https://github.com/wailsapp/wails/pull/3958)
- Поддержка создания пакетов deb, rpm и Arch Linux в Linux —
- Добавлена поддержка универсальных сборок и пакетов для Darwin —
- Документация по событиям добавлена на сайт —
- Шаблоны для sveltekit и sveltekit-ts, настроенные для разработки без SSR
- Ресурсы сборки обновлены с помощью новой команды `wails3 update build-assets` —
- Пример для тестирования HTML Drag and Drop API —
- Поддержка ассоциаций файлов — [leaanthony](https://github.com/leaanthony) в
- Новая команда `wails3 generate runtime` —
- Новый параметр `InitialPosition`, определяющий, следует ли центрировать окно, —
- В пакет `application` добавлены методы `Path` и `Paths` —
- Добавлены параметры Windows `GeneralAutofillEnabled` и `PasswordAutosaveEnabled`
- Добавлена возможность получить окно, из которого вызван метод сервиса, —
- Добавлены параметры `EnabledFeatures` и `DisabledFeatures` для Webview2 —
- ⊞ Новая система DIP для улучшенной поддержки мониторов с высоким DPI —
- ⊞ Параметр имени класса окна — [windom](https://github.com/windom/) в
- Сервисы расширены для обеспечения функциональности плагинов. Автор:
- 🐧 События WindowDidMove / WindowDidResize в
- ⊞ Событие WindowDidResize в
-  добавлено событие ApplicationShouldHandleReopen для обработки значка в Dock
-  в реализацию добавлены getPrimaryScreen/getScreens — @tmclane в
-  добавлен параметр для отображения панели инструментов в полноэкранном режиме на macOS —
- 🐧 добавлена логика onKeyPress для преобразования нажатия клавиши Linux в сочетание клавиш
- 🐧 добавлена задача `run:linux` —
- Экспортирован метод `SetIcon` — [@almas-x](https://github.com/almas-x) в
- Улучшен `OnShutdown` — [@almas-x](https://github.com/almas-x) в
- В интерфейсе `Window` восстановлен метод `ToggleMaximise` —
- В `Environment()` добавлена дополнительная информация. Автор: @leaanthony, в
- Метод `WebviewWindow.IsFocused` открыт через интерфейс `Window` —
- В системе WML добавлена поддержка нескольких событий-триггеров, разделённых пробелами, —
- В комплектный скрипт среды выполнения JS добавлены экспорты ESM —
- В генератор привязок добавлен флаг для использования комплектного скрипта среды выполнения JS вместо
- `setIcon` реализован в Linux — [@abichinger](https://github.com/abichinger)
- В команду dev добавлен флаг `-port` и поддержка переменной окружения
- Добавлены тесты вызовов связанных методов —
- ⊞ добавлен `SetIgnoreMouseEvents` для уже созданного окна —
-  Добавлена возможность задавать уровень расположения окна (порядок) —

### Исправлено

- Исправлены уменьшенные вдвое `Screen.Bounds`, `WorkArea` и `Size` на компьютерах Mac с Retina: значения точек NSScreen преобразуются в пиксели устройства в полях `Physical*`, а значения верхнего уровня `Screen.X`/`Y` заполняются, чтобы обнаружение соприкосновения нескольких мониторов и размещение в рабочей области выполнялись правильно, — в [PR](https://github.com/wailsapp/wails/pull/5168) от @wayneforrest
- Исправлена гонка данных в ScreenManager, вызывавшая взаимную блокировку WebKit DisplayLink при изменении конфигурации дисплеев (например, при горячем подключении внешнего монитора во время перехода в спящий режим или выхода из него)
- Непосредственно задаёт для CFBundleIconName значение appicon, если существует Assets.car, в [PR](https://github.com/wailsapp/wails/pull/5154) от @symball
- Исправлено указание неверных пакетов WebKitGTK командой `wails3 doctor` в Fedora, openSUSE, Arch и NixOS: резервные записи 4.0 удалены, поскольку v3 требует API 4.1 во время компиляции (#5071)
- Исправлено имя пакета webkit2gtk для проверки doctor в openSUSE (`webkit2gtk4_1-devel` → `webkit2gtk3-devel` — правильное имя пакета openSUSE) (#5071)
- Исправлена ошибка `Unexpected token '<'` при отсутствии `/wails/custom.js` в режиме разработки настольного приложения. Добавлен явный обработчик 404 для `/wails/custom.js` и регистронезависимая проверка `Content-Type` в `loadOptionalScript`, чтобы предотвратить внедрение резервных HTML-страниц SPA в качестве JavaScript. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Исправлено состояние подсветки меню в области состояния системы на macOS: теперь при открытом меню значок отображается выбранным (#4910)
- Исправлено появление прикреплённого к области состояния системы окна позади других окон на macOS: теперь используется правильный уровень всплывающего окна (#4910)
- Исправлены неверные примеры импорта `@wailsio/runtime` в документации (#4989)
- Исправлена невозможность свернуть безрамочное окно на Darwin (#4294)
- Устранены многоминутные зависания 20-30 во время `wails3 build` и `wails3 dev` путём исключения `node_modules/` из проверки актуальности go-task. Ранее glob-шаблон `sources: "**/*"` заставлял go-task перечислять и вычислять контрольную сумму каждого файла в `node_modules/` (50000-100000+ файлов с тяжёлыми зависимостями, такими как MUI), что было особенно медленно в Windows/NTFS (#4939)
- Исправлен сбой сборки GTK4, вызванный конфликтом typedef `Screen` языка C с X11 Xlib.h (#4957)
- Унифицировано поведение методов значка в Dock на macOS
- Исправлено применение `InvisibleTitleBarHeight` ко всем окнам macOS вместо только безрамочных окон и окон с прозрачной строкой заголовка (#4960)
- Устранены дрожание и рывки окна при изменении размера за верхние углы с включённым `InvisibleTitleBarHeight`: запуск перетаскивания вблизи краёв окна теперь пропускается (#4960)
- Исправлена генерация отображаемых типов с ключами-перечислениями в привязках JS/TS (#4437), автор — @fbbdev
- Исправлено перетаскивание файлов в Windows, не работавшее при масштабе экрана, отличном от 100%
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом приёме файлов в Windows
- Исправлено использование неверного пространства пикселей для координат сбрасывания файлов в Windows (физические пиксели вместо CSS-пикселей)
- Исправлена ненадёжная работа перетаскивания файлов в Linux при наличии эффектов наведения
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом приёме файлов в Linux
- Исправлено восстановление окна в свёрнутом состоянии, иногда происходившее при его показе или скрытии в Linux/GTK4, посредством `gtk_window_present()` (#4957)
- Исправлено получение и задание позиции окна в Linux/GTK4, всегда возвращавшее 0,0: добавлена условная поддержка X11 через `XTranslateCoordinates`/`XMoveWindow` (#4957)
- Исправлено несоблюдение максимального размера окна в Linux/GTK4: удалённый `gtk_window_set_geometry_hints` заменён ограничением размера на основе сигналов (#4957)
- Исправлено масштабирование DPI в Linux/GTK4: реализованы правильное вычисление PhysicalBounds и поддержка дробного масштабирования через `gdk_monitor_get_scale` (GTK 4.14+)
- Исправлено дублирование пунктов меню при создании новых окон в Linux/GTK4
- Исправлена генерация отображаемых типов с ключами-перечислениями в привязках JS/TS (#4437), автор — @fbbdev
- Исправлено перетаскивание файлов в Windows, не работавшее при масштабе экрана, отличном от 100%
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом приёме файлов в Windows
- Исправлено использование неверного пространства пикселей для координат сбрасывания файлов в Windows (физические пиксели вместо CSS-пикселей)
- Исправлена ненадёжная работа перетаскивания файлов в Linux при наличии эффектов наведения
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом приёме файлов в Linux
- Исправлено масштабирование DPI в Linux/GTK4: реализованы правильное вычисление PhysicalBounds и поддержка дробного масштабирования через `gdk_monitor_get_scale` (GTK 4.14+)
- Исправлено дублирование пунктов меню при создании новых окон в Linux/GTK4
- Исправлена генерация отображаемых типов с ключами-перечислениями в привязках JS/TS (#4437), автор — @fbbdev
- Исправлена проблема «окон-призраков» в macOS, вызванная обращением к API AppKit не из главного потока в App.Window.Current() (#4947), автор — @wimaha
- Исправлена неработающая функция HTML `<input type="file">` в macOS посредством реализации WKUIDelegate runOpenPanelWithParameters (#4862)
- Исправлено неработающее нативное перетаскивание файлов при использовании npm-модуля `@wailsio/runtime` в macOS/Linux (#4953), автор — @leaanthony
- Исправлена генерация привязок для псевдонимов типов из других пакетов (#4578), автор — @fbbdev
- Исправлен аварийный сбой OpenFileDialog в Linux из-за нарушения потокобезопасности GTK (#3683), автор — @ddmoney420
- Исправлен аварийный сбой SIGSEGV при вызове `Focus()` для скрытого или уничтоженного окна (#4890), автор — @ddmoney420
- Исправлена потенциальная паника при задании пустого значка или растрового изображения в Linux (#4923), автор — @ddmoney420
- Исправлен аварийный сбой ErrorDialog при вызове из привязки сервиса в macOS (#3631), автор — @leaanthony
- Добавлено отображение меню в ОС Windows в `v3\examples\dialogs`, автор — @ndianabasi
- Исправлено состояние гонки, вызывавшее TypeError при перезагрузке страницы (#4872), автор — @ddmoney420
- Исправлен неверный вывод тестов генератора привязок путём удаления глобального состояния из метода `Collector.IsVoidAlias()` (#4941), автор — @fbbdev
- Исправлено неработающее средство выбора файлов `<input type="file">` в macOS (#4862), автор — @leaanthony
- Исправлено использование несовместимых систем координат в `Position()` и `SetPosition()` на macOS, приводившее к смещению позиции окна при сохранении и восстановлении состояния (#4816), автор — @leaanthony
- Исправлена ошибка SetProcessDpiAwarenessContext «Access is denied», возникавшая, когда осведомлённость о DPI уже задана в манифесте приложения (#4803)
- Обновлена страница документации о сочетаниях клавиш и исправлен тип параметра функции обратного вызова для `KeyBinding.Add`, автор — @ndianabasi
- Исправлена документация по генерации пользовательской привязки: необходимо использовать `-d String` вместо `-o String`
- Исправлено меню, не удалявшее дочерние элементы при `menu.Update()`
- Исправлены устаревшие ссылки на API Manager в документации (файлы 31 обновлены для использования нового шаблона, например `app.Window.New()`, `app.Event.Emit()` и т. д.), автор — @leaanthony
- Исправлен аварийный сбой Linux при панике в привязанных к JS методах Go из-за переопределения обработчиков сигналов WebKit (#3965), автор — @leaanthony
- Исправлен метод SaveFileDialog.SetFilename(), не имевший эффекта в Linux (#4841), автор — @samstanier
- Исправлено отображение координат сбрасывания как undefined в примере drag-n-drop
- Исправлен сбой создания пакета приложения macOS, когда APP_NAME содержит пробелы (проблема раскрытия фигурных скобок)
- Исправлена паника из-за выхода индекса за границы в Windows при вызове методов сервиса (выполнен откат goccy/go-json)
- Исправлено перетаскивание файлов в Windows, не работавшее при масштабе экрана, отличном от 100%
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом приёме файлов в Windows
- Исправлено использование неверного пространства пикселей для координат сбрасывания файлов в Windows (физические пиксели вместо CSS-пикселей)
- Исправлена ненадёжная работа перетаскивания файлов в Linux при использовании эффектов наведения
- Исправлена поломка внутреннего перетаскивания HTML5 при включённом сбрасывании файлов в Linux
- Обновлены все команды в файлах Taskfile.yml для всех операционных систем, чтобы корректно обрабатывать пробелы в переменных, таких как `APP_NAME`, — @ndianabasi
- Исправлена ошибка аргумента команды при выполнении задачи 'build:universal:lipo:go' в Linux — @wux1an
- Исправлена ошибка Docker «undefined symbol: **<em>ubsan</em>handle_xxxxxxx» при выполнении 'wails3 build GOOS=darwin GOARCH=arm64' в Linux — @wux1an
- Объединена документация по пользовательским протоколам и добавлены разделы об Universal Links — @leaanthony
- Исправлен сбой меню в системном трее Windows при многократном нажатии на значок: добавлена защита от одновременных вызовов TrackPopupMenuEx (#4151) — @leaanthony
- Предотвращён сбой приложения при вызове systray.Run() до app.Run() — @leaanthony
- Исправлен сбой в macOS при переключении видимости окна через Hide()/Show(), когда включён параметр ApplicationShouldTerminateAfterLastWindowClosed (#4389), — @leaanthony
- Исправлена утечка памяти в контекстных меню macOS и Windows при их многократном открытии (#4012) — @leaanthony
- Исправлено отсутствие повторного использования нативных ресурсов контекстного меню в macOS, из-за которого при каждом показе создавалось новое меню (#4012), — @leaanthony
- Исправлена ошибка, из-за которой нажатие на значок приложения в Dock macOS не показывало скрытые окна, если приложение было запущено с `Hidden: true` (#4583), — @leaanthony
- Исправлена ошибка, из-за которой диалог печати в macOS не открывался вследствие неправильного типа указателя на окно в вызове CGO (#4290), — @leaanthony
- Исправлен сбой меню окна в Wayland, вызванный обращением appmenu-gtk-module к ещё не созданному окну (#4769), — @leaanthony
- Исправлен сбой приложения GTK, если имя приложения содержит недопустимые символы (пробелы, круглые скобки и т. д.), — @leaanthony
- Исправлена ошибка «not enough memory» при инициализации перетаскивания в Windows (#4701) — @overlordtm
- Исправлена ошибка, из-за которой файловый менеджер в Linux открывал неправильный каталог вследствие некорректного экранирования URI (#4397), — @leaanthony
- Исправлен сбой сборки AppImage в современных дистрибутивах Linux (Arch, Fedora 39+, Ubuntu 24.04+): добавлено автоматическое обнаружение разделов ELF `.relr.dyn` и отключение удаления отладочных символов (#4642) — @leaanthony
- Исправлена ошибка, из-за которой `wails doctor` ошибочно сообщал об установленных пакетах webkit в Fedora и системах на основе DNF (#4457), — @leaanthony
- Исправлена ошибка, из-за которой стандартный `config.yml` запускал `wails3 dev` со сборкой для рабочей среды, — @mbaklor
- Исправлены заглушки служб iOS, вызывавшие сбои сборки из-за импорта несуществующего пакета, — @leaanthony
- Исправлена ошибка структурированного журналирования в методах debug/info, приводившая к сообщениям «no formatting directives», — @leaanthony
- Удалены временные отладочные инструкции печати, случайно попавшие в код при слиянии мобильной платформы, — @leaanthony
- Исправлен сбой WebKitGTK в Wayland с графическими процессорами NVIDIA (ошибка 71 Protocol error): добавлено автоматическое отключение средства визуализации DMA-BUF — @leaanthony
- Устранено игнорирование альфа-значения в `application.WebviewWindowOptions.BackgroundColour` в Linux ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Исправлена ошибка, из-за которой значок в системном трее Windows по умолчанию не наследовал значок приложения, если пользовательский значок не задан (#4704)
- Добавлено отслеживание владения `HICON`, чтобы уничтожались только созданные пользователем дескрипторы, что предотвращает сбои при перезапуске Explorer (#4653).
- При уничтожении теперь освобождаются слушатель системной темы Windows и сохранённые значки трея, что устраняет утечки горутин и контекстов устройств (#4653).
- Подсказки в трее обрезаются до 127 единиц UTF-16, чтобы избежать повреждения суррогатных пар и многобайтовых глифов (#4653).
- Исправлен сбой задачи упаковки для Windows (#4667)
- Исправлена переменная appicon для AppImage Linux в Linux Taskfile [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Исправлена ошибка сборки для Windows, вызванная изменением сигнатуры go-webview2 v1.0.22 (#4513, #4645)
- Исправлена переменная appicon для AppImage Linux в Linux Taskfile [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Исправлен диапазон протоколов в linux desktop.tmpl: удалено `<.Info.Protocol>` до `<.Protocol>` — @Tolfx в #4510
- Исправлена ошибка повторного определения в демонстрации liquid glass в [#4542](https://github.com/wailsapp/wails/pull/4542) — @Etesam913
- Исправлено обновление меню в системном трее Linux [#4604](https://github.com/wailsapp/wails/issues/4604) — [@JackDoan](https://github.com/JackDoan)
- Исправлено появление белого окна в Windows при создании скрытого окна — @leaanthony в [#4612](https://github.com/wailsapp/wails/pull/4612)
- Исправлен путь импорта пакета notifications в документации — @rxliuli в [#4617](https://github.com/wailsapp/wails/pull/4617)
- Исправлена ошибка, из-за которой перетаскивание не работало при использовании npm-пакета @wailsio/runtime (#4489), — @leaanthony в #4616
- Windows: устранены мерцание окна при запуске и некорректный показ скрытых окон в [PR](https://github.com/wailsapp/wails/pull/4600) — @leaanthony.
- Исправлены проблемы с размером окна при разворачивании в Wayland (https://github.com/wailsapp/wails/issues/4429) — [@samstanier](https://github.com/samstanier)
- Исправлены проблемы с размером окна при разворачивании в Wayland (https://github.com/wailsapp/wails/issues/4429) — [@samstanier](https://github.com/samstanier)
- Исправлена ошибка повторного определения в демонстрации liquid glass в [#4542](https://github.com/wailsapp/wails/pull/4542) — @Etesam913
- Исправлена ошибка, из-за которой AssetServer мог аварийно завершить работу в macOS, в [#4576](https://github.com/wailsapp/wails/pull/4576) — @jghiloni
- Исправлена проблема компиляции при сборке с NextJs. Исправление внесено в [#4585](https://github.com/wailsapp/wails/pull/4585) — @rev42
- Исправлены конвейеры ночного выпуска в [#4597](https://github.com/wailsapp/wails/pull/4597) — @riadafridishibly
- Исправлена ошибка повторного определения в демонстрации liquid glass в [#4542](https://github.com/wailsapp/wails/pull/4542) — @Etesam913
- Исправлена ошибка, из-за которой AssetServer мог аварийно завершить работу в macOS, в [#4576](https://github.com/wailsapp/wails/pull/4576) — @jghiloni
- Исправлена проблема компиляции при сборке с NextJs. Исправление внесено в [#4585](https://github.com/wailsapp/wails/pull/4585) — @rev42
- Исправлены конвейеры ночного выпуска в [#4597](https://github.com/wailsapp/wails/pull/4597) — @riadafridishibly
- Исправлена ошибка повторного определения в демонстрации liquid glass в [#4542](https://github.com/wailsapp/wails/pull/4542) — @Etesam913
- Исправлен SetBackgroundColour в Windows — @PPTGamer в [PR](https://github.com/wailsapp/wails/pull/4492)
- Документация обновлена с учётом изменений после рефакторинга Manager API — @yulesxoxo в [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Исправлена переменная appicon файла .desktop для Linux в Linux Taskfile [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Документация обновлена с учётом изменений после рефакторинга Manager API — @yulesxoxo в [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Исправлена ошибка разыменования нулевого указателя в Windows, о которой сообщалось в [#4456](https://github.com/wailsapp/wails/issues/4456), — @leaanthony в [#4460](https://github.com/wailsapp/wails/pull/4460)
- В WKWebView для macOS добавлена поддержка `allowsBackForwardNavigationGestures`, обеспечивающая навигационные жесты смахивания двумя пальцами (#1857)
- Исправлена ошибка, из-за которой onClick не работал для пунктов меню, изначально отключённых, — @leaanthony в [PR #4469](https://github.com/wailsapp/wails/pull/4469). Спасибо @IanVS за первоначальное исследование.
- Исправлено отсутствие очистки сервера Vite при сбое сборки (#4403)
- Исправлена паника при закрытии или отмене `SaveFileDialog` в Windows. Исправление внесено @hkhere в [PR](https://github.com/wailsapp/wails/pull/4284)
- Исправлено перетаскивание на уровне HTML в Windows — [@mbaklor](https://github.com/mbaklor) в [#4259](https://github.com/wailsapp/wails/pull/4259)
- В macOS WKWebView добавлена поддержка `allowsBackForwardNavigationGestures` для навигационных жестов смахивания двумя пальцами (#1857)
- Исправлена ошибка, из-за которой onClick не работал для пунктов меню, изначально отключённых, — @leaanthony в [PR #4469](https://github.com/wailsapp/wails/pull/4469). Спасибо @IanVS за первоначальное исследование.
- Исправлено отсутствие очистки сервера Vite при сбое сборки (#4403)
- Исправлен разбор уведомлений в Windows — @popaprozac в [PR](https://github.com/wailsapp/wails/pull/4450)
- В команде doctor исправлена проверка зависимостей Windows SDK — [@kodumulo](https://github.com/kodumulo) в [#4390](https://github.com/wailsapp/wails/issues/4390)
- Исправлено разыменование nil-указателя в processURLRequest на Mac — [@etesam913](https://github.com/etesam913) в [#4366](https://github.com/wailsapp/wails/pull/4366)
- Исправлена ошибка Linux, препятствовавшая работе диалогов с фильтрами, — [@bh90210](https://github.com/bh90210) в [#4287](https://github.com/wailsapp/wails/pull/4287)
- Исправлены проблемы меню Edit в Windows и Linux — [@leaanthony](https://github.com/leaanthony) в [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Минимальная версия системы в файлах .plist для macOS обновлена с 10.13.0 до 10.15.0 — [@AkshayKalose](https://github.com/AkshayKalose) в [#3981](https://github.com/wailsapp/wails/pull/3981)
- Исправлена проблема с пропуском идентификатора окна — [@leaanthony](https://github.com/leaanthony)
- Исправлена проблема с nil-меню при вызове RegisterContextMenu — [@leaanthony](https://github.com/leaanthony)
- Исправлены циклические зависимости в выводе генератора привязок — [@fbbdev](https://github.com/fbbdev) в [#4001](https://github.com/wailsapp/wails/pull/4001)
- Исправлены ошибки использования до объявления в выводе генератора привязок — [@fbbdev](https://github.com/fbbdev) в [#4001](https://github.com/wailsapp/wails/pull/4001)
- В генератор привязок передаются флаги сборки — [@fbbdev](https://github.com/fbbdev) в [#4023](https://github.com/wailsapp/wails/pull/4023)
- Чтобы Taskfile для Windows работал на платформах, отличных от Windows, разделители в путях заменены на прямые косые черты — [@leaanthony](https://github.com/leaanthony)
- Исправлены события Mac и Mac JS — [@leaanthony](https://github.com/leaanthony)
- Исправлена взаимная блокировка событий в macOS — [@leaanthony](https://github.com/leaanthony)
- Исправлена ошибка `Parameter incorrect` при инициализации Window в Windows, когда предоставлен HTML, но отсутствует JS, — [@leaanthony](https://github.com/leaanthony)
- Исправлен размер префикса ответа, используемого сервером ресурсов для определения типа содержимого, — [@fbbdev](https://github.com/fbbdev) в [#4049](https://github.com/wailsapp/wails/pull/4049)
- Исправлена обработка ответов, отличных от 404, для корневого индексного пути на сервере ресурсов — [@fbbdev](https://github.com/fbbdev) в [#4049](https://github.com/wailsapp/wails/pull/4049)
- Исправлено неопределённое поведение генератора привязок при проверке свойств обобщённых типов — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлен вывод генератора привязок для моделей, базовый тип которых имеет свойства, отличные от свойств именованной обёртки, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлены типы ключей map и предварительная обработка в выводе генератора привязок — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлен вывод генератора привязок для структур, реализующих интерфейсы маршалинга, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлено обнаружение циклов типов с участием обобщённых типов в генераторе привязок — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлены недопустимые ссылки на неэкспортируемые модели в выводе генератора привязок — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Внедряемый код перемещён в конец файлов служб — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлена обработка ошибок операций закрытия файлов в генераторе привязок — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Отключены предупреждения для служб, которые определяют методы жизненного цикла или HTTP, но не имеют других привязанных методов, — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- Исправлено отсутствие нижнего колонтитула Hello World в шаблонах, не использующих React, при светлой системной цветовой схеме — [@marcus-crane](https://github.com/marcus-crane) в [#4056](https://github.com/wailsapp/wails/pull/4056)
- Исправлены скрытые пункты меню в macOS — [@leaanthony](https://github.com/leaanthony)
- Исправлены обработка и форматирование ошибок в обработчиках сообщений — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Исправлен пропуск завершения работы службы при выходе из приложения — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Обновления меню теперь гарантированно выполняются в основном потоке — [@leaanthony](https://github.com/leaanthony)
- Механизм перетаскивания и изменения размера стал надёжнее и теперь точнее соответствует ожидаемому поведению платформы — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Исправлена ошибка [#4097](https://github.com/wailsapp/wails/issues/4097), из-за которой Webpack/angular удалял код инициализации среды выполнения, — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Исправлены пункты меню, изначально скрытые, — [@IanVS](https://github.com/IanVS) в [#4116](https://github.com/wailsapp/wails/pull/4116)
- Исправлена ошибка, из-за которой assetFileServer не отдавал файлы `.html` для запроса без расширения, когда `[request]` не существует, а `[request].html` существует
- Исправлены пути генерации значков — [@robin-samuel](https://github.com/robin-samuel) в [#4125](https://github.com/wailsapp/wails/pull/4125)
- Исправлено отсутствие отправки событий `fullscreen`, `unfullscreen`, `unminimise` и `unmaximise` — [@oSethoum](https://github.com/osethoum) в [#4130](https://github.com/wailsapp/wails/pull/4130)
- Исправлена ошибка NSIS, возникавшая из-за неверного префикса версии по умолчанию в конфигурации, — [@robin-samuel](https://github.com/robin-samuel) в [#4126](https://github.com/wailsapp/wails/pull/4126)
- Исправлена функция среды выполнения Dialogs, возвращавшая в Windows пути с экранированными символами, — [TheGB0077](https://github.com/TheGB0077) в [#4188](https://github.com/wailsapp/wails/pull/4188)
- Исправлен путь обнаружения Webview2 в HKCU — [@leaanthony](https://github.com/leaanthony).
- Исправлена проблема с вводом в macOS — [@leaanthony](https://github.com/leaanthony).
- Исправлено имя файла задачи генерации значка Windows; автор: [@yulesxoxo](https://github.com/yulesxoxo), [#4219](https://github.com/wailsapp/wails/pull/4219).
- Исправлена проблема с прозрачностью безрамочных окон; автор: [@leaanthony](https://github.com/leaanthony), на основе работы @kron.
- Исправлены вызовы фокусировки для отключённого или свёрнутого окна; автор: [@leaanthony](https://github.com/leaanthony), на основе работы @kron.
- Исправлено отсутствие системных треев после перезапуска панели задач; автор: [@leaanthony](https://github.com/leaanthony), на основе работы @kron.
- Исправлено отсутствие реализации Flush() в fallbackResponseWriter, [#4245](https://github.com/wailsapp/wails/pull/4245)
- Исправлено отсутствие реализации Flush() в fallbackResponseWriter; автор: [@superDingda], [#4236](https://github.com/wailsapp/wails/issues/4236)
- Исправлены сбои при закрытии окна macOS во время незавершённого асинхронного вызова привязанной функции Go; автор: [@joshhardy](https://github.com/joshhardy), [#4354](https://github.com/wailsapp/wails/pull/4354)
- Исправлено состояние гонки при запуске режима эффективности Windows; автор: [@leaanthony](https://github.com/leaanthony)
- Исправлена очистка дескрипторов значков Windows; автор: [@leaanthony](https://github.com/leaanthony).
- Исправлено `OpenFileManager` в Windows; автор: [@PPTGamer](https://github.com/PPTGamer), [#4375](https://github.com/wailsapp/wails/pull/4375).
- Исправлены параметры минимальной и максимальной ширины в Linux; автор: @atterpac, [#3979](https://github.com/wailsapp/wails/pull/3979)
- Определения типов шаблонов TypeScript обновлены посредством повышения версии npm; автор: @atterpac, [#3966](https://github.com/wailsapp/wails/pull/3966)
- Исправлена ссылка на CSS в шаблоне SvelteKit; автор: @atterpac, [#3945](https://github.com/wailsapp/wails/pull/3945)
- Гарантирован вызов ключевых обратных вызовов в window run() в главном потоке; автор: [@leaanthony](https://github.com/leaanthony)
- Исправлены примеры диалогового окна выбора каталога; автор: [@leaanthony](https://github.com/leaanthony)
- Создана новая страница ошибки на китайском языке для случая, когда отсутствует index.html; автор: [@leaanthony](https://github.com/leaanthony)
-  Гарантировано выполнение обратного вызова `windowDidBecomeKey` в главном потоке; автор: [@leaanthony](https://github.com/leaanthony)
-  Добавлена поддержка полноэкранного режима для безрамочных окон; автор: [@leaanthony](https://github.com/leaanthony)
-  Улучшена логика уничтожения окон; автор: [@leaanthony](https://github.com/leaanthony)
-  Исправлена логика позиционирования окна, прикреплённого к системному трею; автор: [@leaanthony](https://github.com/leaanthony)
-  Добавлена поддержка полноэкранного режима для безрамочных окон; автор: [@leaanthony](https://github.com/leaanthony)
- Исправлена обработка событий; автор: [@leaanthony](https://github.com/leaanthony)
- Исправлена логика завершения работы окна; автор: [@leaanthony](https://github.com/leaanthony)
- Общий Taskfile теперь по умолчанию генерирует привязки TypeScript для шаблонов TypeScript; автор: [@leaanthony](https://github.com/leaanthony)
- Исправлено закрытие приложения при сообщении WM_CLOSE, когда нет открытых окон либо используется только системный трей; автор: [@mmalcek](https://github.com/mmalcek), [#3990](https://github.com/wailsapp/wails/pull/3990)
- Исправлена сборка с garble; автор: @5aaee9, [#3192](https://github.com/wailsapp/wails/pull/3192)
- Исправлены сборки NSIS для Windows; автор: [@leaanthony](https://github.com/leaanthony)
- Исправлена взаимная блокировка в диалоговом окне Linux при множественном выборе, вызванная незакрытым
- Исправлена кроссплатформенная очистка файлов .syso при сборке для Windows; автор:
- Исправлена компиляция AppImage для amd64; автор: @atterpac, в
- Исправлено обновление ресурсов сборки; автор: @ansxuman, в
- Исправлена реализация `OnClick` и `OnRightClick` для системного трея Linux; автор: @atterpac
- Исправлена неработоспособность `AlwaysOnTop` на Mac; автор:
-  Исправлено добавление дубликата в `application.NewEditMenu`
- 🐧 Исправлена компиляция для aarch64
- ⊞ Исправлены пункты меню с группами переключателей; автор:
- Исправлена ошибка при сборке запускаемого пакета .app в macOS, когда 'name' и 'outputfilename'
- Исправлена ошибка использования customEventProcessor в примере перетаскивания; автор:
- 🐧 Исправлена ошибка компиляции в Linux, появившаяся после добавления IgnoreMouseEvents; автор:
- ⊞ Исправлена ошибка генерации файла значка syso; автор:
- 🐧 Интегрировано исправление для нативной работы в Wayland из
- Не создавать привязки для внутренних методов сервисов в
- ⊞ Исправлена паника при запуске системного трея в
- Не создавать привязки для внутренних методов сервисов в
- ⊞ Исправлена паника при запуске системного трея в
- Выполнен значительный рефакторинг пунктов меню и обработки событий. Пока главным образом улучшает macOS. Автор:
- Исправлены тесты после рефакторинга плагинов и событий в
- ⊞ Исправлено предупреждение `Failed to unregister class Chrome_WidgetWin_0`. Автор:
- Проблемы с модулями
- Исправлена передача сообщений о событиях изменения размера; автор: [atterpac](https://github.com/atterpac), в
- 🐧 Исправлена ошибка обработки темы в NixOS; автор:
- Исправлена установка проекта между разными томами в Windows; автор:
- Исправлен CSS шаблона React для отображения нижнего колонтитула; автор:
- Исправлено появление процессов-зомби при работе в режиме разработки путём обновления refresh до последней версии
- Исправлено получение файлов WebKit для AppImage; автор: [Atterpac](https://github.com/atterpac)
- Исправлена проверка пакетов apt в Doctor; автор: [Atterpac](https://github.com/Atterpac), в
- Исправлено зависание приложения при выходе в Darwin; автор: @5aaee9, в
- Исправлены цвета фона примеров в Windows; автор:
- Исправлены контекстные меню по умолчанию; автор: [mmghv](https://github.com/mmghv), в
- Исправлены шестнадцатеричные значения клавиш со стрелками в Darwin; автор:
- Перетаскивание в Windows приведено в рабочее состояние. Добавил:
- Исправлена ошибка в Doctor для Linux, возникавшая, когда у пользователя не установлены подходящие драйверы
- Исправлено масштабирование DPI при запуске в Windows. Изменил: [@almas-x](https://github.com/almas-x), в
- Исправлена строка замены в `go.mod`: теперь используются относительные пути. Это исправляет пути Windows с
- Исправлена обработка щелчка по системному трею macOS при отсутствии прикреплённого окна; автор:
- Исправлена ошибка сборки для Windows из-за неизвестного параметра; автор:
- Исправлен сбой в Windows при щелчке левой кнопкой мыши по значку системного трея при отсутствии
- Исправлен неверный baseURL при двукратном открытии окна; автор: @5aaee9, PR
- Исправлен порядок ветвей if в методе `WebviewWindow.Restore`; автор:
- Обеспечлено корректное вычисление `startURL` при нескольких вызовах `GetStartURL`, когда
- Тип структуры `Screen` в JS приведён в соответствие с её аналогом в Go; автор:
- Исправлен метод `WML.Reload`, чтобы обеспечить надлежащую очистку зарегистрированных событий
- Исправлено немедленное закрытие пользовательского контекстного меню в Linux; автор:
- Исправлен путь вывода и расширение файлов моделей, создаваемых генератором привязок
- Исправлены пути импорта файлов моделей в JS-коде, создаваемом генератором привязок
- Исправлено перетаскивание в некоторых дистрибутивах Linux —
- Исправлено отсутствие задачи для macOS при использовании `wails3 task dev` —
- Исправлена ошибка присваивания в nil-карту при регистрации событий —
- Исправлена десериализация параметров привязанных методов —
- Исправлена обработка нескольких возвращаемых значений привязанных методов —
- Исправлено обнаружение командой doctor установки npm не через системный менеджер пакетов
- Исправлено отсутствие MicrosoftEdgeWebview2Setup.exe. Спасибо
- Исправлен случайный сбой в Linux из-за обработки идентификаторов окон — @leaanthony. На основе
- Исправлен сбой systemTray.setIcon в Linux —
- Исправлено применение рамки окна при первом вызове функции `setFrameless` в

### Изменено

- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ**: ключи карт в сгенерированных привязках JS/TS теперь помечаются как необязательные, чтобы точно отражать семантику карт Go. Доступ к значению карты в TypeScript теперь возвращает `T | undefined` вместо `T`, поэтому требуются проверки на null или утверждения типов (#4943) — `@fbbdev`
- Использование `Event` заменено на `Events` в соответствии с изменениями в `@wailsio/runtime`; соответствующие вызовы функций в документации в `Features/Events/Event System` также обновлены — @AbdelhadiSeddar
- Параметры `EnabledFeatures`, `DisabledFeatures` и `AdditionalBrowserArgs` перенесены из параметров отдельных окон в `Options.Windows` уровня приложения (#4559) — @leaanthony
- Обновлён README примера `Drag N Drop` и подчёркнуто, что пример демонстрирует `Internal Drag and Drop` — @ndianabasi
- Уровень различных отладочных журналов изменён с Info на Debug (@mbaklor)
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** параметр `EnableDragAndDrop` в параметрах окна переименован в `EnableFileDrop`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** параметр `DropZoneDetails` в контексте события переименован в `DropTargetDetails`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** метод `DropZoneDetails()` объекта `WindowEventContext` переименован в `DropTargetDetails()`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** событие `WindowDropZoneFilesDropped` удалено; вместо него используйте `WindowFilesDropped`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** HTML-атрибут `data-wails-dropzone` заменён на `data-file-drop-target`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** CSS-класс наведения `wails-dropzone-hover` заменён на `file-drop-target-active`
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ:** из Windows удалены параметры `DragEffect`, `OnEnterEffect` и `OnOverEffect` (они относились к удалённому IDropTarget)
- Для всей обработки JSON во время выполнения (привязки методов, события, запросы webview, уведомления, kvstore) теперь используется goccy/go-json, что повышает производительность на 21-63% и сокращает количество выделений памяти на 40-60%
- Оптимизирована структура BoundMethod, а флаг isVariadic кэшируется для сокращения накладных расходов каждого вызова
- Для методов с `<=8` аргументами используется размещаемый в стеке буфер аргументов, что позволяет избежать выделений памяти в куче
- Оптимизирован сбор результатов вызовов методов, чтобы не выделять срез для единственного возвращаемого значения
- Для кэша MIME-типов используется sync.Map, что повышает производительность при конкурентном доступе
- Для чтения тела запроса HTTP-транспорта используется пул буферов
- Канал CloseNotify в определителе типа содержимого теперь выделяется отложенно, что сокращает количество выделений памяти на запрос
- Из сервера ресурсов удалено отладочное журналирование CSS
- Карта расширений MIME-типов дополнена поддержкой более 50 распространённых веб-форматов (шрифтов, аудио, видео и т. д.)
- Обновлена документация параметров `X/Y` для Window — @ruhuang2001
- Обновлена документация `Frontend Runtime`: добавлены дополнительные параметры генерации привязок для фронтенда — @ndianabasi
- Обновлена страница документации сервера ресурсов Wails v3 — @ndianabasi
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ**: функции диалогов уровня пакета (`application.InfoDialog()`, `application.QuestionDialog()` и т. д.) удалены. Вместо них используйте менеджер `app.Dialog`: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Документация по диалогам приведена в соответствие с фактическим API: используются `app.Dialog.*`, `AddButton()` с обратными вызовами (не `SetButtons()`), `SetDefaultButton(*Button)` (не строка), `AddFilter()` (не `SetFilters()`), `SetFilename()` (не `SetDefaultFilename()`) и `app.Dialog.OpenFile().CanChooseDirectories(true)` для выбора папки
- **КРИТИЧЕСКОЕ ИЗМЕНЕНИЕ**: теперь по умолчанию создаются производственные сборки. Чтобы создавать сборки для разработки, задайте `DEV=true` в своих Taskfile. Для примеров создайте новый проект — @leaanthony
- При отправке пользовательского события без аргументов данных или с одним таким аргументом значение данных теперь присваивается непосредственно полю Data без оборачивания в срез — [@fbbdev](https://github.com/fbbdev) в [#4633](https://github.com/wailsapp/wails/pull/4633)
- Системные лотки Windows теперь учитывают `SystemTray.Show()`/`Hide()`, переключая `NIS_HIDDEN`, поэтому приложения могут полностью исчезать и возвращаться (#4653).
- При регистрации системного лотка повторно используются уже разрешённые значки, `NOTIFYICON_VERSION_4` задаётся один раз, а `NIF_SHOWTIP` включается, чтобы всплывающие подсказки восстанавливались после перезапуска Explorer (#4653).
- macOS: для центрирования окна используется `visibleFrame` вместо `frame`, чтобы исключить области строки меню и Dock
- macOS: для центрирования окна используется `visibleFrame` вместо `frame`, чтобы исключить области строки меню и Dock
- При запуске `wails3 update build-assets` с параметром `-config` значения, заданные через параметры `-product*`,
- `window.NativeWindowHandle()` -> `window.NativeWindow()` — @leaanthony в [#4471](https://github.com/wailsapp/wails/pull/4471)
- Переработана внутренняя обработка окон — @leaanthony в [#4471](https://github.com/wailsapp/wails/pull/4471)
- Удалены `application.WindowIDKey` и `application.WindowNameKey` (заменены на `application.WindowKey`) — [@leaanthony](https://github.com/leaanthony)
- ContextMenuData теперь возвращает string вместо any — [@leaanthony](https://github.com/leaanthony)
- В привязках JS/TS поля классов с типами массивов фиксированной длины теперь инициализируются с ожидаемой длиной, а не остаются пустыми — [@fbbdev](https://github.com/fbbdev) в [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData теперь возвращает string вместо any — [@leaanthony](https://github.com/leaanthony)
- `application.NewService` больше не принимает параметры в виде необязательного аргумента (вместо этого используйте `application.NewServiceWithOptions`) — [@leaanthony](https://github.com/leaanthony) в [#4024](https://github.com/wailsapp/wails/pull/4024)
- Удалена зависимость `nanoid` — [@leaanthony](https://github.com/leaanthony)
- Обновлён пример Window для стилей окон mica/acrylic/tabbed — [@leaanthony](https://github.com/leaanthony)
- Из привязок JS/TS удалены файлы моделей `internal.js/ts`; теперь все модели находятся в `models.js/ts` — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- В привязках JS/TS именованные типы больше никогда не отображаются как псевдонимы других именованных типов; прежнее поведение теперь применяется только к псевдонимам — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- В привязках JS/TS в режиме классов поля структур, тип которых является параметром типа, помечаются как необязательные и никогда не инициализируются автоматически — [@fbbdev](https://github.com/fbbdev) в [#4045](https://github.com/wailsapp/wails/pull/4045)
- ESLint удалён из шаблонов — [@IanVS](https://github.com/IanVS) в [#4059](https://github.com/wailsapp/wails/pull/4059)
- Дата авторского права обновлена до 2025 — [@IanVS](https://github.com/IanVS) в [#4037](https://github.com/wailsapp/wails/pull/4037)
- Добавлена документация для event.Sender — [@IanVS](https://github.com/IanVS) в [#4075](https://github.com/wailsapp/wails/pull/4075)
- Поддержка Go 1.24 — [@leaanthony](https://github.com/leaanthony)
- Хуки `ServiceStartup` теперь вызываются при вызове `App.Run`, а не в `application.New` — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Ошибки `ServiceStartup` теперь возвращаются из `App.Run`, а не завершают процесс — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Вызовы привязок и диалогов из JS теперь отклоняются с объектами ошибок вместо строк — [@fbbdev](https://github.com/fbbdev) в [#4066](https://github.com/wailsapp/wails/pull/4066)
- Улучшено позиционирование меню в области уведомлений Windows — [@leaanthony](https://github.com/leaanthony)
- Среда выполнения JS перенесена на TypeScript — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Среда выполнения инициализируется сразу после импорта; ждать загрузки окна больше не нужно — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Среда выполнения больше не экспортирует метод init. Для её инициализации можно использовать импорт ради побочных эффектов — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Связанные методы теперь возвращают `CancellablePromise`, который при отмене отклоняется с `CancelError`. Фактический результат вызова отбрасывается — [@fbbdev](https://github.com/fbbdev) в [#4100](https://github.com/wailsapp/wails/pull/4100)
- Типы встроенных сервисов теперь единообразно называются `Service` — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Функции создания встроенных сервисов с параметрами теперь единообразно называются `NewWithConfig` — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Метод `Select` сервиса `sqlite` теперь называется `Query` для согласованности с API Go — [@fbbdev](https://github.com/fbbdev) в [#4067](https://github.com/wailsapp/wails/pull/4067)
- Шаблоны: среда выполнения перенесена в "dependencies", файлы package.json упорядочены — [@IanVS](https://github.com/IanVS) в [#4133](https://github.com/wailsapp/wails/pull/4133)
- При разработке создаёт пакеты приложений и подписывает их специальной подписью, чтобы задействовать некоторые API macOS — [@popaprozac](https://github.com/popaprozac) в [#4171](https://github.com/wailsapp/wails/pull/4171)
- Ресурсы сборки перемещены в каталоги для соответствующих платформ — [@leaanthony](https://github.com/leaanthony)
- Файлы Taskfile перемещены и переименованы в каталоги для соответствующих платформ — [@leaanthony](https://github.com/leaanthony)
- Значительно улучшено взаимодействие при отсутствии `index.html` — [@leaanthony](https://github.com/leaanthony)
- [Windows] Повышена производительность сворачивания и восстановления окна — [@leaanthony](https://github.com/leaanthony). Основано на исходном [PR](https://github.com/wailsapp/wails/pull/3955) от [562589540](https://github.com/562589540)
- Удалён параметр `ShouldClose` (вместо него зарегистрируйте хук для events.Common.WindowClosing) — [@leaanthony](https://github.com/leaanthony)
- [Windows] Уменьшено мерцание при открытии окна — [@leaanthony](https://github.com/leaanthony)
- Удалена `Window.Destroy`, поскольку эта функция предназначалась для внутреннего использования — [@leaanthony](https://github.com/leaanthony)
- События `WindowClose` переименованы в `WindowClosing` — [@leaanthony](https://github.com/leaanthony)
- Сборки фронтенда теперь используют окружение Vite "development" или "production" в зависимости от типа сборки — [@leaanthony](https://github.com/leaanthony)
- go-webview2 обновлён до v1.19 — [@leaanthony](https://github.com/leaanthony)
- Обеспечено использование форка taskfile — @leaanthony
- Форк Taskfile обновлён для устранения проблем с версией при установке с помощью
- Использование форка Taskfile для устранения проблем с версией при установке с помощью
- `service.OnStartup` теперь завершает работу приложения при ошибке и выполняет
- Обмен сообщениями о щелчках по значку в области уведомлений переработан для лучшего соответствия действиям пользователя —
- Во встраиваемые ресурсы включён `all:frontend/dist` для поддержки фреймворков, которые создают
- Рефакторинг Taskfile — [leaanthony](https://github.com/leaanthony) в
- Обновление до `go-webview2` v1.0.16 —
- Тип `Screen` исправлен: теперь он включает `ID`, а не `Id` —
- Версия Wails в `go.mod.tmpl` обновлена для поддержки `application.ServiceOptions` —
- Исправлено определение имени сервиса — [windom](https://github.com/windom/) в
- mkdocs serve теперь использует Docker — [leaanthony](https://github.com/leaanthony)
- Конфигурация для разработки объединена в `config.yml` —
- Диалог в области уведомлений теперь по умолчанию использует значок приложения, если он доступен (Windows) —
- Улучшено отображение сведений о GPU и памяти в macOS —
- Удалены `WebviewGpuIsDisabled` и `EnableFraudulentWebsiteWarnings`
- Изменение API событий: `On`/`Emit` -> пользовательские события, `OnApplicationEvent` ->
- Исправлен API событий в Linux — [TheGB0077](https://github.com/TheGB0077) в
- [CI] улучшены workflows GitHub Actions и разрешён их запуск также в форках и
- `AbsolutePosition()` переименован в `Position()` —
- Зависимость Linux от WebKit обновлена с webkitgtk2-4.0 до webkit2gtk-4.1, чтобы
- Встроенный скрипт среды выполнения JS теперь является модулем ESM: теги script, которые его импортируют
- Пакет `@wailsio/runtime` не публикует свой API в `window.wails`
- Модуль API окна `@wailsio/runtime/src/window` теперь предоставляет содержащий его
- API окна JS обновлён в соответствии с текущим Go `WebviewWindow`
- Генератор привязок теперь по умолчанию использует вызовы по идентификатору. Параметр CLI `-id`
- Новая структура кода привязок: ранее выходные файлы были распределены по каталогам
- Поле структуры `application.Options.Bind` переименовано в
- Новый синтаксис привязки сервисов: теперь экземпляры сервисов необходимо оборачивать в
- Индикатор выполнения отключён в нетерминальной среде или среде CI —

### Удалено

- **НЕСОВМЕСТИМОЕ ИЗМЕНЕНИЕ**: параметры `EnabledFeatures`, `DisabledFeatures` и `AdditionalLaunchArgs` удалены из настроек `WindowsWindow` отдельных окон. Вместо них используйте параметры уровня приложения `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures` и `Options.Windows.AdditionalBrowserArgs`. Эти флаги применяются глобально к общей среде WebView2 (#4559), автор — @leaanthony
- Нативная реализация `IDropTarget` в Windows удалена в пользу подхода на основе JavaScript (соответствует поведению v2)
- Зависимость github.com/wailsapp/mimetype удалена в пользу расширенной карты расширений и функции стандартной библиотеки http.DetectContentType, что уменьшило размер двоичного файла примерно на 1.2 МБ
- Зависимость gopkg.in/ini.v1 удалена благодаря реализации минимального парсера файлов .desktop для файлового менеджера Linux, что позволило сэкономить примерно 45 КБ
- Зависимость samber/lo удалена из кода среды выполнения благодаря использованию пакета slices из стандартной библиотеки Go 1.21+ и минимальных внутренних вспомогательных функций, что позволило сэкономить примерно 310 КБ
- Из обработчика схем URL в Darwin удалены отладочные вызовы printf (#4834)
- **НЕСОВМЕСТИМОЕ ИЗМЕНЕНИЕ**: событие `linux:WindowLoadChanged` удалено. Чтобы определить, когда WebView завершил загрузку, используйте вместо него `linux:WindowLoadFinished` (#3896), автор — @leaanthony

### Несовместимые изменения

- **Рефакторинг API менеджеров**: плоская структура API приложения преобразована в упорядоченный набор менеджеров для улучшения организации кода и упрощения поиска возможностей; автор — [@leaanthony](https://github.com/leaanthony), [#4359](https://github.com/wailsapp/wails/pull/4359)
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
- Методы `Path` и `Paths` перемещены в пакет `application`; автор — [@leaanthony](https://github.com/leaanthony)
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

## v3.0.0-alpha.11 — 2025-07-12

## Добавлено

## v3.0.0-alpha.10 — 2025-07-06

### Критические изменения

### Добавлено

### Исправлено

### Изменено

## v3.0.0-alpha.9 — 2025-01-13

### Добавлено

### Исправлено

### Изменено

## v3.0.0-alpha.8.3 — 2024-12-07

### Изменено

## v3.0.0-alpha.8.2 — 2024-12-07

### Изменено

`go install` от @leaanthony

## v3.0.0-alpha.8.1 — 2024-12-07

### Изменено

`go install` от @leaanthony

## v3.0.0-alpha.8 — 2024-12-06

### Добавлено

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### Изменено

`service.OnShutdown`для всех ранее запущенных сервисов, от @atterpac в   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac в [#3907](https://github.com/wailsapp/wails/pull/3907)   подпапки, от @atterpac в   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) в   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) в   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (заменено параметрами `EnabledFeatures` и `DisabledFeatures`), от   [leaanthony](https://github.com/leaanthony)

### Исправлено

переменная канала, от @michael-freling в   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) в   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   в [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) в   [#3841](https://github.com/wailsapp/wails/pull/3841)   роль `PasteAndMatchStyle` в меню редактирования на Darwin, от   [johnmccabe](https://github.com/johnmccabe) в   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) в   [#3854](https://github.com/wailsapp/wails/pull/3854), от   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   различаются, от @nickisworking в   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 — 2024-09-18

### Добавлено

[mmghv](https://github.com/mmghv) в   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) и   [leaanthony](https://github.com/leaanthony) в   [#3570](https://github.com/wailsapp/wails/pull/3570)

### Изменено

События приложения `OnWindowEvent` -> события окна, от   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   ветки с префиксом `v3/` или `v3-`, от   [stendler](https://github.com/stendler) в   [#3747](https://github.com/wailsapp/wails/pull/3747)

### Исправлено

[etesam913](https://github.com/etesam913) в   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) в   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) в   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) в   [#3614](https://github.com/wailsapp/wails/pull/3614), от   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720), от   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693), от   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720), от   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693), от   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746), от   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 — 2024-07-30

### Исправлено

## v3.0.0-alpha.5 — 2024-07-30

### Добавлено

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   щелчок по значку, от @5aaee9 в [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) в   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   в [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) в   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) в   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   на основе [PR](https://github.com/wailsapp/wails/pull/2044) от @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   пакет npm, от [@fbbdev](https://github.com/fbbdev) в   [#3334](https://github.com/wailsapp/wails/pull/3334)   в [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` от [@abichinger](https://github.com/abichinger) в   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) в   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) в   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) в   [#3674](https://github.com/wailsapp/wails/pull/3674)

### Исправлено

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) в   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) в   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) в   [#3477](https://github.com/wailsapp/wails/pull/3477)   автор — [Atterpac](https://github.com/atterpac), в   [#3320](https://github.com/wailsapp/wails/pull/3320).   в [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) в   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) в   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) в   [PR](https://github.com/wailsapp/wails/pull/3039)   установлено. Добавлено пользователем [@pylotlight](https://github.com/pylotlight) в   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   пробелы — @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) в PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) в PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   прикреплённое окно [tw1nk](https://github.com/tw1nk) в PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) в   [#3279](https://github.com/wailsapp/wails/pull/3279)   присутствует `FRONTEND_DEVSERVER_URL`.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) в   [#3295](https://github.com/wailsapp/wails/pull/3295)   обработчики событий — автор [@fbbdev](https://github.com/fbbdev), в   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) в   [#3330](https://github.com/wailsapp/wails/pull/3330)   генератор — автор [@fbbdev](https://github.com/fbbdev), в   [#3334](https://github.com/wailsapp/wails/pull/3334)   генератор — автор [@fbbdev](https://github.com/fbbdev), в   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) в   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) в   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) в   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) в   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) в   [#3431](https://github.com/wailsapp/wails/pull/3431)   автор — [@pekim](https://github.com/pekim), в   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) от   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) в   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows — автор [@bruxaodev](https://github.com/bruxaodev/), в   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Изменено

[mmghv](https://github.com/mmghv) в   [#3611](https://github.com/wailsapp/wails/pull/3611)   поддержка Ubuntu 24.04 LTS — автор [atterpac](https://github.com/atterpac), в   [#3461](https://github.com/wailsapp/wails/pull/3461)   должен иметь атрибут `type="module"`. Автор —   [@fbbdev](https://github.com/fbbdev), в   [#3295](https://github.com/wailsapp/wails/pull/3295)   объект и не запускает систему WML. Это сделано для улучшения   инкапсуляции. При необходимости систему WML можно запустить вручную, вызвав   новый метод `WML.Enable`. Встроенный сценарий среды выполнения JS по-прежнему выполняет обе   операции автоматически. Автор — [@fbbdev](https://github.com/fbbdev), в   [#3295](https://github.com/wailsapp/wails/pull/3295)   объект окна в качестве экспорта по умолчанию. Импортировать   отдельные методы с помощью именованного импорта ESM или импорта пространства имён больше нельзя.   API. У некоторых методов изменились имена или прототипы, а именно: `Screen`   заменён на `GetScreen`; `GetZoomLevel`/`SetZoomLevel` заменены на `GetZoom`/`SetZoom`;   `GetZoom`, `Width` и `Height` теперь возвращают значения напрямую, а не оборачивают   их в объекты. Автор — [@fbbdev](https://github.com/fbbdev), в   [#3295](https://github.com/wailsapp/wails/pull/3295)   удалён. Чтобы вернуться к вызовам по имени, используйте параметр CLI `-names`.   Автор — [@fbbdev](https://github.com/fbbdev), в   [#3468](https://github.com/wailsapp/wails/pull/3468)   получали имена по содержащим их пакетам; теперь используются полные пути импорта Go,   включая путь модуля. Автор — [@fbbdev](https://github.com/fbbdev), в   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Автор — [@fbbdev](https://github.com/fbbdev), в   [#3468](https://github.com/wailsapp/wails/pull/3468)   вызов `application.NewService`. Автор — [@fbbdev](https://github.com/fbbdev), в   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) в   [#3574](https://github.com/wailsapp/wails/pull/3574)
