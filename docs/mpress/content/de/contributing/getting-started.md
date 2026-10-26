---
title: "Erste Schritte"
description: "So tragen Sie zu Wails bei v3"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## Willkommen!

Vielen Dank für Ihr Interesse, zu Wails beizutragen! Dieser Leitfaden hilft Ihnen bei Ihrem ersten Beitrag.

## Voraussetzungen

Stellen Sie zunächst sicher, dass Folgendes vorhanden ist:

- **Go 1.26+** ist installiert ([herunterladen](https://go.dev/dl/))
- **Node.js 20+** und **npm** sind installiert ([herunterladen](https://nodejs.org/))
- **Git** ist für Ihr GitHub-Konto konfiguriert
- Grundkenntnisse in Go und JavaScript/TypeScript

### Plattformspezifische Anforderungen

**macOS:**

- Xcode Command Line Tools: `xcode-select --install`

**Windows:**

- MSYS2 oder eine ähnliche Unix-artige Umgebung wird empfohlen
- WebView2-Laufzeit (unter Windows 11 normalerweise vorinstalliert)

**Linux:**

- `gcc`, `pkg-config`, `libgtk-4-dev`, `libwebkitgtk-6.0-dev` (standardmäßiger GTK4-Stack)
- Installation über: `sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev` (Debian/Ubuntu)
- Installieren Sie für den veralteten `-tags gtk3`-Build-Pfad außerdem `libgtk-3-dev` und `libwebkit2gtk-4.1-dev`

## Überblick über den Beitragsprozess

Der typische Arbeitsablauf für Beiträge umfasst folgende Schritte:

1. **Forken und klonen** – Erstellen Sie Ihre eigene Kopie des Wails-Repositorys
2. **Einrichten** – Bauen Sie die Wails-CLI und überprüfen Sie Ihre Umgebung
3. **Branch erstellen** – Erstellen Sie einen Feature-Branch für Ihre Änderungen
4. **Entwickeln** – Nehmen Sie Ihre Änderungen gemäß unseren Programmierstandards vor
5. **Testen** – Führen Sie Tests aus, um sicherzustellen, dass alles funktioniert
6. **Committen** – Committen Sie mit eindeutigen, konventionellen Commit-Nachrichten
7. **Einreichen** – Öffnen Sie einen Pull Request zur Überprüfung
8. **Überarbeiten** – Reagieren Sie auf Feedback und nehmen Sie Anpassungen vor
9. **Zusammenführen** – Nach der Genehmigung werden Ihre Änderungen Bestandteil von Wails!

## Schritt-für-Schritt-Anleitung

Wählen Sie die Art Ihres Beitrags:

@tabs
[Fehlerbehebung]
@steps
### Fehler suchen oder melden
- Prüfen Sie in den [GitHub-Issues](https://github.com/wailsapp/wails/issues), ob der Fehler bereits gemeldet wurde
- Erstellen Sie andernfalls ein neues Issue mit Schritten zum Reproduzieren des Fehlers
- Warten Sie auf eine Bestätigung, bevor Sie mit der Arbeit beginnen

### Forken und klonen
Forken Sie das Repository unter [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)

Klonen Sie Ihren Fork:

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### Bauen und überprüfen
Bauen Sie Wails und überprüfen Sie, ob Sie den Fehler reproduzieren können:

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Reproduce the bug to understand it
```

### Branch für die Fehlerbehebung erstellen
Erstellen Sie einen Branch für Ihre Fehlerbehebung:

```bash
git checkout -b fix/issue-123-window-crash
```

### Fehler beheben
- Nehmen Sie nur die mindestens erforderlichen Änderungen vor, um den Fehler zu beheben
- Refaktorieren Sie keinen nicht damit zusammenhängenden Code
- Fügen Sie Tests hinzu oder aktualisieren Sie diese, um Regressionen zu verhindern

```bash
# Make your changes
# Add tests in *_test.go files
```

### Fehlerbehebung testen
Führen Sie Tests aus, um sicherzustellen, dass die Fehlerbehebung funktioniert:

```bash
go test ./...

# Test the specific package
go test ./pkg/application -v

# Run with race detector
go test ./... -race
```

### Fehlerbehebung committen
Committen Sie mit einer eindeutigen Nachricht:

```bash
git commit -m "fix: prevent window crash when closing during initialization

Fixes #123"
```

### Pull Request einreichen
Pushen Sie Ihre Änderungen und erstellen Sie einen PR:

```bash
git push origin fix/issue-123-window-crash
```

Geben Sie in Ihrer PR-Beschreibung Folgendes an:

- Erläutern Sie den Fehler und seine Ursache
- Beschreiben Sie Ihre Fehlerbehebung
- Verweisen Sie auf das Issue: „Fixes #123“
- Beschreiben Sie das Verhalten vor und nach der Änderung

### Auf Feedback reagieren
Bearbeiten Sie die Review-Kommentare und aktualisieren Sie Ihren PR nach Bedarf.

@end

[WEP (Erweiterung)]
@steps
### WEP verfassen
- Lesen Sie den [Prozess für WEPs (Wails Enhancement Proposals)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)
- Kopieren Sie die WEP-Vorlage nach `v3/wep/proposals/<name>/proposal.md`
- Öffnen Sie einen Entwurfs-PR mit dem Titel `[WEP] <title>`, der ausschließlich das WEP enthält
- Warten Sie vor der Implementierung auf die Entscheidung eines Maintainers

### Forken und klonen
Forken Sie das Repository unter [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)

Klonen Sie Ihren Fork:

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### Entwicklungsumgebung einrichten
Bauen Sie Wails und überprüfen Sie Ihre Umgebung:

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### Feature-Branch erstellen
Erstellen Sie einen aussagekräftig benannten Branch:

```bash
git checkout -b feat/window-transparency-support
```

### Feature implementieren
- Befolgen Sie unsere [Programmierstandards](/contributing/standards/)
- Beschränken Sie die Änderungen auf das Feature
- Schreiben Sie sauberen, dokumentierten Code
- Fügen Sie umfassende Tests hinzu

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### Gründlich testen
Testen Sie Ihr Feature:

```bash
# Unit tests
go test ./pkg/application -v

# Integration test - create a test app
cd ..
./wails3 init -n feature-test
cd feature-test
# Add code using your new feature
../wails3 dev
```

### Feature dokumentieren
- Fügen Sie allen öffentlichen APIs Dokumentationskommentare hinzu
- Aktualisieren Sie die relevante Dokumentation in `/docs/mpress/content/`
- Fügen Sie gegebenenfalls Beispiele hinzu

### Konventionsgemäß committen
Verwenden Sie Conventional Commits:

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### Pull Request einreichen
Pushen Sie Ihre Änderungen und erstellen Sie einen PR:

```bash
git push origin feat/window-transparency-support
```

Geben Sie in Ihrem PR Folgendes an:

- Beschreiben Sie das Feature und seine Anwendungsfälle
- Zeigen Sie Beispiele oder Screenshots
- Führen Sie alle Breaking Changes auf
- Verweisen Sie auf den angenommenen WEP-PR

### Anhand des Reviews überarbeiten
Maintainer können Änderungen anfordern. Seien Sie geduldig und arbeiten Sie konstruktiv mit ihnen zusammen.

@end

[Dokumentation]
Korrektur-PRs sind willkommen, ohne dass zuvor ein Issue eröffnet werden muss. Reine Dokumentationskorrekturen benötigen keinen fehlschlagenden Codetest. Informationen zur Installation von M-Press, zu Quellpfaden, Vorschau, Validierung und den PR-Schritten finden Sie unter [Dokumentation mit M-Press verfassen](/contributing/documentation/).

@end

## Issues finden, die Sie bearbeiten können

- Suchen Sie nach Labels vom Typ [`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue)
- Prüfen Sie Issues vom Typ [`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted)
- Durchsuchen Sie [offene Issues](https://github.com/wailsapp/wails/issues) und bitten Sie darum, Ihnen eines zuzuweisen

## Hilfe erhalten

- **Discord:** Treten Sie dem [Wails-Discord](https://discord.gg/JDdSxwjhGf) bei
- **Diskussionen:** Veröffentlichen Sie einen Beitrag in den [GitHub Discussions](https://github.com/wailsapp/wails/discussions)
- **Issues:** Eröffnen Sie für einen reproduzierbaren Fehler ein Issue; verwenden Sie Discussions für Fragen und einen WEP-PR für Erweiterungen

## Verhaltenskodex

Seien Sie respektvoll, konstruktiv und aufgeschlossen. Wir bauen eine freundliche Community auf, in der wir gemeinsam großartige Software entwickeln.

## Nächste Schritte

- Richten Sie Ihre [Entwicklungsumgebung](/contributing/setup/) ein
- Lesen Sie unsere [Programmierrichtlinien](/contributing/standards/)
- Erkunden Sie die [technische Dokumentation](/contributing/overview/)
