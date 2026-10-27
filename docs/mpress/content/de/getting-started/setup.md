---
title: "Einrichtung"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Experimentell"}
Der Einrichtungsassistent ist neu und wurde hauptsächlich unter Linux getestet. Wenn Probleme auftreten, [melden Sie sie](https://github.com/wailsapp/wails/issues/4904) bitte und führen Sie stattdessen die Schritte zur [manuellen Installation](/getting-started/installation/#platform-specific-dependencies) aus.

@end

## Schnellstart

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

Der Assistent wird in Ihrem Browser geöffnet und führt Sie durch die Prüfung der Abhängigkeiten, die Projektstandardeinstellungen und die optionale Einrichtung plattformübergreifender Builds.

Danach können Sie Ihr erstes Projekt erstellen:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## Funktionsweise

- **Prüft Abhängigkeiten** – Überprüft Go, npm und die plattformspezifischen Werkzeuge
- **Konfiguriert Standardeinstellungen** – Angaben zum Autor, Bundle-ID-Präfix und bevorzugte Vorlagen
- **Plattformübergreifende Builds** – Optionale Docker-Einrichtung für Builds auf jedem Hostsystem
- **Codesignierung** – Optionale Einrichtung für macOS, Windows und Linux

Die Konfiguration wird in `$XDG_CONFIG_HOME/wails/defaults.yaml` gespeichert, wobei ersatzweise `~/.config/wails/defaults.yaml` verwendet wird, und von `wails3 init` genutzt. Sie ist von der Datei `build/config.yml` eines Projekts getrennt.

## Unterbefehle

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Probleme?

1. Führen Sie `wails3 doctor` aus, um Probleme zu diagnostizieren
2. Führen Sie die Schritte zur [manuellen Installation](/getting-started/installation/#platform-specific-dependencies) aus
3. [Melden Sie das Problem](https://github.com/wailsapp/wails/issues/4904) zusammen mit der Ausgabe von `wails3 doctor`
