---
title: "Abstürze debuggen"
description: "Diagnoseberichte und Minidumps aus einer Wails-Anwendung erfassen"
slug: "guides/debugging-crashes"
sourcePath: "guides/debugging-crashes.md"
---

Wenn in der Produktion etwas schiefläuft, reichen Protokolle selten aus. Dieser Leitfaden behandelt das Paket `github.com/wailsapp/wails/v3/pkg/debug`, das strukturierte Diagnoseberichte und unter Windows vollständige Minidumps erzeugt, die sich in WinDbg oder Visual Studio öffnen lassen.

Das Paket ist bewusst klein und bietet zwei Einstiegspunkte:

- `debug.Dump(...)` — einen Windows-Minidump des aktuellen Prozesses schreiben.
- `debug.Report(...)` — eine ausführliche Diagnoseaufnahme mit optionalem Minidump erfassen.

Beide laufen unter dem Token des aufrufenden Benutzers. Keine Rechteerhöhung, kein `SeDebugPrivilege`, kein `OpenProcess`: Für den Dump des aktuellen Prozesses wird das Pseudohandle `GetCurrentProcess` verwendet, das ACL-Prüfungen beim Zugriff auf den eigenen Prozess umgeht.

## Kurzreferenz

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

`debug.Report` gibt immer einen `*CrashReport` zurück, auch bei teilweisem Fehlschlag. Aufrufer erhalten alle vor dem Fehler erfolgreich erfassten Daten; der Bericht ist deshalb niemals `nil`, wenn `err != nil` gilt.

## Die Struktur `CrashReport`

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

Serialisiere mit `encoding/json` für Absturzmeldedienste oder durchlaufe die Felder direkt, wenn dich nur ein Teil davon interessiert.

## Integration mit `PanicHandler`

Am hilfreichsten ist es, beim Abfangen einer Panic durch Wails einen Bericht und optional einen Minidump zu erfassen und die kombinierten Daten anschließend an die eigene Meldepipeline zu übergeben:

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

Wails ruft `debug.Report` **nicht** automatisch auf. Absturzberichte erfordern häufig die Zustimmung des Benutzers oder das Entfernen personenbezogener Daten; deshalb gehört die Entscheidung in deinen Handler. Im [Leitfaden zur Panic-Behandlung](/guides/panic-handling/) erfährst du, wann Wails `PanicHandler` aufruft und wie du auch selbst gestartete Goroutinen abdeckst.

## Manuelles Debugging

`debug.Dump` und `debug.Report` sind auch ohne Panic nützlich:

- Hänger erkennen: Wenn ein Watchdog auslöst, erstelle einen Prozessdump und öffne ihn in WinDbg, um genau zu sehen, worauf jeder Thread blockiert war.
- Auffällige Zustände melden: Steigt etwa die Zahl der Goroutinen unbegrenzt, erfasse einen Bericht, lasse die Anwendung weiterlaufen und untersuche ihn später.
- Supportpakete auf Anfrage: Verbinde einen Menüpunkt „Fehler melden“ mit `debug.Report(debug.WithDump())` und hänge das Ergebnis an das Supportticket des Benutzers an.

## Windows-Minidumps in der Praxis

Minidumps verwenden standardmäßig einen detaillierten, aber kompakten Satz von Flags:

- `MiniDumpNormal`
- `MiniDumpWithThreadInfo`
- `MiniDumpWithHandleData`
- `MiniDumpWithUnloadedModules`

Dadurch entstehen Dumps von etwa 5–50 MB, die alle Threadstacks, Handle-Tabellen und Modullisten erhalten. Damit lässt sich rekonstruieren, wo sich jede Goroutine zum Zeitpunkt des Dumps befand.

`debug.WithFullMemory()` beziehungsweise `debug.WithDumpFullMemory()` bei `Report` ergänzt `MiniDumpWithFullMemory` und erfasst damit den gesamten Adressraum des Prozesses. Das hilft bei der Untersuchung beschädigter Pointer, erzeugt aber Dateien von Hunderten MB bis zu mehreren GB. Verwende es für eingehende Analysen.

Öffne die erzeugte `.dmp`-Datei mit `windbg -z C:\path\to\your.dmp` in WinDbg oder über Datei → Absturzabbild öffnen in Visual Studio. In Release-Builds werden Symbole üblicherweise entfernt. Verwende für hilfreiche Stacktraces die `.pdb` desselben Builds oder die Go-Binärdatei selbst, wenn sie mit Debug-Informationen gebaut wurde.

## Andere Plattformen als Windows

`debug.Dump` liefert außerhalb von Windows einen Fehler wegen fehlender Implementierung. Linux und macOS haben andere Werkzeuge zur nachträglichen Analyse, etwa Core-Dumps mit `GOTRACEBACK=crash`, `rr` oder plattformspezifische Absturzmelder, die sich nicht direkt auf die API `MiniDumpWriteDump` abbilden lassen. `debug.Report` funktioniert weiterhin: Du erhältst alles außer einer Dump-Datei.

Der Fehler von `Dump` ist prüfbar, damit du kontrolliert mit eingeschränktem Funktionsumfang fortfahren kannst:

```go
path, err := debug.Dump()
if err != nil {
    log.Printf("minidump unavailable on %s: %v", runtime.GOOS, err)
    // Continue with just debug.Report() or pd.StackTrace from wails.
}
```

## Sicherheitshinweise

- **Keine Rechteerhöhung erforderlich.** Das Paket vermeidet bewusst `SeDebugPrivilege` und den Zugriff über `OpenProcess` / fremde Prozess-IDs, wie ihn Werkzeuge zum Abgreifen von Zugangsdaten verwenden. Es erstellt nur Dumps des aufrufenden Prozesses.
- **Dump-Dateien enthalten den gesamten Speicherinhalt.** Dazu gehören gerade verwendete Zugangsdaten, Authentifizierungstoken, Benutzerdaten und Verschlüsselungsschlüssel. Behandle einen Dump wie ein vollständiges Speicherabbild: Schütze ihn bei der Übertragung, bereinige gespeicherte Daten und erwäge das Entfernen sensibler Inhalte vor dem Hochladen an Dritte.
- **Der Standardpfad ist `os.TempDir()`.** Unter Windows entspricht das `%LOCALAPPDATA%\Temp`: benutzerspezifisch und nicht für alle lesbar, aber dauerhaft gespeichert. Verschiebe Dumps an einen sicheren Ort, wenn du sie aufbewahrst.
