---
title: "Anwendungen erstellen"
description: "Wails-Anwendung erstellen und paketieren"
slug: "guides/build/building"
sourcePath: "guides/build/building.md"
---

Wails v3 verwendet [Task](https://taskfile.dev) als Build-System. Die Befehle `wails3 build` und `wails3 package` sind praktische Wrapper für Task.

## Erstellen

Für die aktuelle Plattform erstellen:

```bash
wails3 build
```

Für eine bestimmte Plattform erstellen:

```bash
wails3 build GOOS=windows
wails3 build GOOS=darwin
wails3 build GOOS=linux

# With architecture
wails3 build GOOS=darwin GOARCH=arm64

# Environment variable style works too
GOOS=windows wails3 build
```

Die Ausgabe wird im Verzeichnis `bin/` gespeichert.

@note{type="tip"}
Für die Cross-Kompilierung für macOS oder Linux von einer anderen Plattform aus ist Docker erforderlich. Informationen zur Einrichtung finden Sie unter [Plattformübergreifende Builds](/guides/build/cross-platform/).

@end

## Entwicklung

Anwendung mit Hot Reload ausführen:

```bash
wails3 dev
```

Dadurch wird ein Datei-Watcher gestartet, der Ihre Anwendung bei Änderungen neu erstellt und startet. Der Frontend-Entwicklungsserver wird standardmäßig auf Port 9245 ausgeführt.

```bash
# Custom port
wails3 dev -port 3000

# Enable HTTPS
wails3 dev -s
```

## Paketierung

Anwendung für die Verteilung paketieren:

```bash
wails3 package
wails3 package GOOS=windows
wails3 package GOOS=darwin
wails3 package GOOS=linux
```

Dadurch werden plattformspezifische Pakete erstellt:

- **Windows**: NSIS-Installationsprogramm – siehe [Paketierung für Windows](/guides/build/windows/)
- **macOS**: Anwendungspaket (`.app`) – siehe [Paketierung für macOS](/guides/build/macos/)
- **Linux**: AppImage, deb und rpm – siehe [Paketierung für Linux](/guides/build/linux/)

## Benutzerdefinierte Build-Tags

Benutzerdefinierte Go-Build-Tags mit dem Flag `-tags` übergeben:

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI, CGO-free)
wails3 build -tags server

# Combine multiple tags
wails3 build -tags gtk3,customtag
```

Die Tags werden als `EXTRA_TAGS` an das zugrunde liegende Taskfile weitergeleitet. Weitere Informationen finden Sie unter [Server-Build](/guides/server-build/) und [Linux-Paketierung – Unterstützung für älteres GTK3](/guides/build/linux/#legacy-gtk3-support).

## Task direkt verwenden

Verwenden Sie Task direkt, um mehr Kontrolle zu erhalten:

```bash
# List available tasks
wails3 task --list

# Verbose output
wails3 task build -v

# Dry run
wails3 task --dry

# Force rebuild
wails3 task build -f

# Pass variables
wails3 task darwin:build ARCH=amd64
```

Plattformspezifische Tasks wie `linux:create:deb` oder `darwin:build:universal` sind nur über Task verfügbar.

## Assets generieren

Symbole neu generieren oder die Build-Konfiguration aktualisieren:

```bash
wails3 generate icons -input build/appicon.png
wails3 update build-assets -name "MyApp" -config build/config.yml -dir build
```
