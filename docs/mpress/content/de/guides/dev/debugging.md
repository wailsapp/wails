---
title: "Debugging"
description: "Probleme untersuchen und die App-Leistung analysieren"
slug: "guides/dev/debugging"
sourcePath: "guides/dev/debugging.md"
---

Dieser Leitfaden stellt verschiedene Werkzeuge vor, mit denen Sie mögliche Leistungsprobleme in Ihrer Wails-App untersuchen und analysieren können. Dabei verwenden Sie

- [`runtime/trace`](https://pkg.go.dev/runtime/trace), um Leistungsdiagramme zu erstellen und sie in einem Browser zu untersuchen

## Leistungstraces erstellen

@steps
### App für das Tracing vorbereiten
Stellen Sie in der Nähe des Einstiegspunkts Ihres Programms sicher, dass Code nach folgendem Muster vorhanden ist

```go
// Create the file to store our trace data within
traceFile, err := os.Create("trace.out")
if err != nil {
  log.Fatalf("trace.out could not be created: %v", err)
}

// Start the trace
if err := trace.Start(traceFile); err != nil {
  _ = traceFile.Close()
  log.Fatalf("trace.start could not start: %v", err)
}

// Trace cleanup on exit. Alternatively,
defer func() {
  trace.Stop()
  _ = traceFile.Close()
}()

...Start your wails app here...
```

Die Ausgabe wird unter `trace.out` in Ihrem aktuellen Arbeitsverzeichnis gespeichert und enthält Messwerte für die gesamte Laufzeit Ihrer App. Der Standard-Trace enthält weitgehend unverarbeitete Daten. Lesen Sie daher die weiterführende Dokumentation zum Tracing, um zusätzliche Kontextinformationen hinzuzufügen und die tatsächlich aufgezeichneten Daten einzuschränken. Beispiel: `WithRegion, NewTask, Log`

### App ausführen und Trace erstellen
Während Ihre App ausgeführt wird, wird bis zum Beenden fortlaufend ein Trace erstellt. Führen Sie daher einige Aktionen in Ihrer App aus und beenden Sie sie anschließend

### Visualisierungswerkzeuge installieren
Für einige Ansichten ist [Graphviz](https://graphviz.org/) erforderlich. Führen Sie `dot -V` aus, um zu prüfen, ob es installiert ist

```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz

# Arch
sudo pacman -S graphviz
```

### Trace visualisieren
Sobald die Trace-Daten verfügbar sind, können Sie die Weboberfläche starten

```bash
go tool trace trace.out
```

Dadurch sollte sich Ihr Standardbrowser mit der Startseite des Trace-Ereignis-Viewers öffnen. Empfohlen wird ein Chromium-basierter Browser.

Für Erstnutzer ist die Ansicht `Syscall profile` wahrscheinlich am nützlichsten. Sie zeigt mit Zeitangaben detailliert, was das Programm ausführt

@end
