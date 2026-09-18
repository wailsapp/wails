---
title: "Отладка сбоев"
description: "Сбор диагностических отчётов и минидампов из приложения Wails"
slug: "guides/debugging-crashes"
sourcePath: "guides/debugging-crashes.md"
---

Когда в рабочей среде что-то идёт не так, журналов редко бывает достаточно. Этот раздел описывает пакет `github.com/wailsapp/wails/v3/pkg/debug`, который создаёт структурированные диагностические отчёты, а в Windows — полные минидампы для открытия в WinDbg или Visual Studio.

Пакет намеренно небольшой и имеет две точки входа:

- `debug.Dump(...)` — записать минидамп Windows текущего процесса.
- `debug.Report(...)` — собрать подробный диагностический снимок с необязательным минидампом.

Обе функции выполняются с токеном вызывающего пользователя. Без повышения привилегий, без `SeDebugPrivilege`, без `OpenProcess`: для дампа текущего процесса используется псевдодескриптор `GetCurrentProcess`, обходящий проверки ACL при доступе к собственному процессу.

## Краткий справочник

```go
import "github.com/wailsapp/wails/v3/pkg/debug"

// 1. Just a minidump (Windows only).
path, err := debug.Dump()

// 2. Dump to a specific path.
path, err := debug.Dump(debug.WithPath("C:\\crashes\\my.dmp"))

// 3. Full-memory minidump (large, gigabytes).
path, err := debug.Dump(debug.WithFullMemory())

// 4. Full diagnostic report, no minidump.
r, err := debug.Report()

// 5. Report + minidump.
r, err := debug.Report(debug.WithDump())

// 6. Report + minidump at a specific path, full-memory.
r, err := debug.Report(
    debug.WithDumpPath("C:\\crashes\\my.dmp"),
    debug.WithDumpFullMemory(),
)
```

`debug.Report` всегда возвращает `*CrashReport`, даже при частичном сбое: вызывающий код получает всё, что удалось собрать до ошибки, поэтому отчёт никогда не равен `nil`, когда `err != nil`.

## Структура `CrashReport`

```go
type CrashReport struct {
    Timestamp   time.Time          // when the report was generated
    System      SystemInfo         // OS, hardware, CPU, GPU from doctor
    Build       BuildInfo          // Go version, buildmode, compiler, CGO flag
    Crash       *CrashInfo         // process, memory, modules, env vars
    Diagnostics []DiagnosticResult // doctor's health-check results
    DumpPath    string             // set only if WithDump() was passed
}
```

Сериализуйте данные через `encoding/json` для сервисов отчётов о сбоях или обходите поля напрямую, если нужна лишь часть информации.

## Интеграция с `PanicHandler`

Наиболее полезный подход — собирать отчёт и при необходимости минидамп в момент перехвата паники Wails, а затем передавать объединённые данные в собственный процесс обработки отчётов:

```go
app := application.New(application.Options{
    PanicHandler: func(pd *application.PanicDetails) {
        // Always: capture a diagnostic snapshot plus a minidump.
        report, reportErr := debug.Report(debug.WithDump())
        if reportErr != nil {
            log.Printf("debug.Report: %v", reportErr)
        }

        // Now you have:
        //   pd.Error, pd.StackTrace, pd.FullStackTrace  — from wails
        //   report.DumpPath                              — minidump (Windows)
        //   report.System / report.Crash / report.Build — context
        //
        // Ship it off (Sentry, S3, support ticket, local log...).
        mycrashservice.Upload(pd, report)
    },
})
```

Wails **не** вызывает `debug.Report` автоматически: отчёты о сбоях часто требуют согласия пользователя или удаления персональных данных, поэтому решение принимает ваш обработчик. См. [руководство по обработке паник](/guides/panic-handling/), чтобы узнать, когда Wails вызывает `PanicHandler` и как охватить горутины, запущенные вашим кодом.

## Ручная отладка

`debug.Dump` и `debug.Report` полезны и вне обработки паник:

- Обнаружение зависаний: при срабатывании сторожевого механизма создайте дамп процесса, чтобы открыть его в WinDbg и точно увидеть, на чём заблокирован каждый поток.
- Отчёты об аномальном состоянии: например, если число горутин растёт без ограничений, сохраните отчёт, оставьте приложение работать и исследуйте данные позже.
- Пакеты для поддержки по запросу: свяжите пункт меню «Сообщить об ошибке» с `debug.Report(debug.WithDump())` и приложите результат к обращению пользователя.

## Минидампы Windows на практике

По умолчанию минидампы используют подробный, но компактный набор флагов:

- `MiniDumpNormal`
- `MiniDumpWithThreadInfo`
- `MiniDumpWithHandleData`
- `MiniDumpWithUnloadedModules`

В результате получаются дампы размером примерно 5–50 МБ, сохраняющие стеки всех потоков, таблицы дескрипторов и списки модулей. Этого достаточно, чтобы восстановить положение каждой горутины в момент создания дампа.

`debug.WithFullMemory()` (или `debug.WithDumpFullMemory()` для `Report`) добавляет `MiniDumpWithFullMemory`, захватывающий всё адресное пространство процесса. Это полезно для выяснения причин повреждения конкретного указателя, но размеры файлов достигают сотен МБ или ГБ. Используйте для углублённых исследований.

Откройте полученный `.dmp` в WinDbg командой `windbg -z C:\path\to\your.dmp` либо через Файл → Открыть аварийный дамп в Visual Studio. Из выпускных сборок символы обычно удаляют: для полезных трассировок стека используйте `.pdb` той же сборки или сам исполняемый файл Go, если он собран с отладочной информацией.

## Платформы, отличные от Windows

`debug.Dump` возвращает ошибку отсутствующей реализации вне Windows. В Linux и macOS используются другие средства анализа после сбоя: core dump через `GOTRACEBACK=crash`, `rr` или платформенные средства отчётов. Они не соответствуют напрямую API `MiniDumpWriteDump`. Функция `debug.Report` продолжает работать: вы получите всё, кроме файла дампа.

Ошибка от `Dump` допускает проверку, чтобы приложение могло корректно продолжить работу с ограниченной функциональностью:

```go
path, err := debug.Dump()
if err != nil {
    log.Printf("minidump unavailable on %s: %v", runtime.GOOS, err)
    // Continue with just debug.Report() or pd.StackTrace from wails.
}
```

## Замечания о безопасности

- **Повышение привилегий не требуется.** Пакет намеренно избегает `SeDebugPrivilege` и пути через `OpenProcess` / PID другого процесса, используемого инструментами сбора учётных данных. Он создаёт дамп только вызывающего процесса.
- **Файлы дампа содержат всё находящееся в памяти.** Это включает используемые учётные данные, токены аутентификации, пользовательские данные и ключи шифрования. Обращайтесь с дампом как с полным образом памяти: защищайте при передаче, очищайте сохранённые данные и рассмотрите удаление конфиденциальных сведений перед отправкой третьим лицам.
- **Путь по умолчанию — `os.TempDir()`.** В Windows это `%LOCALAPPDATA%\Temp`: каталог отдельного пользователя, не доступный для чтения всем, но постоянный. Переместите дампы в безопасное место, если собираетесь их хранить.
