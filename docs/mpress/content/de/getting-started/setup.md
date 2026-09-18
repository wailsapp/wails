---
title: "Einrichtung"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Experimentell"}
Der Einrichtungsassistent ist neu und wurde hauptsächlich unter Linux getestet. Wenn Probleme auftreten, [melden Sie sie](https://github.com/wailsapp/wails/issues/4904) und führen Sie stattdessen die [Schritte zur manuellen Installation](/getting-started/installation/#platform-specific-dependencies) aus.

@end

## Schnellstart

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard
wails3 setup
```

Der Assistent wird in Ihrem Browser geöffnet und führt Sie durch die Prüfung der Abhängigkeiten, die Projektstandardeinstellungen und die optionale Einrichtung plattformübergreifender Builds.

Danach können Sie Ihr erstes Projekt erstellen:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## Funktionsumfang

- **Prüft Abhängigkeiten** – Überprüft Go, npm und plattformspezifische Werkzeuge
- **Konfiguriert Standardeinstellungen** – Autoreninformationen, Bundle-ID-Präfix und bevorzugte Vorlagen
- **Plattformübergreifende Builds** – Optionale Docker-Einrichtung für Builds auf jedem Hostsystem
- **Codesignierung** – Optionale Einrichtung für macOS, Windows und Linux

Die Konfiguration wird unter `~/.config/wails/config.yaml` gespeichert und von `wails3 init` verwendet.

## Unterbefehle

```bash
wails3 setup signing      # Configure code signing
wails3 setup entitlements # Configure macOS entitlements
```

## Probleme?

1. Führen Sie `wails3 doctor` aus, um Probleme zu diagnostizieren
2. Führen Sie die [Schritte zur manuellen Installation](/getting-started/installation/#platform-specific-dependencies) aus
3. [Melden Sie das Problem](https://github.com/wailsapp/wails/issues/4904) und fügen Sie die Ausgabe von `wails3 doctor` bei
