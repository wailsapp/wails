---
title: "Einrichtung"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Experimentell"}
Der Einrichtungsassistent ist neu und wurde hauptsächlich unter Linux getestet. Wenn Probleme auftreten, [melden Sie sie bitte](https://github.com/wailsapp/wails/issues/4904) und führen Sie stattdessen die [manuellen Installationsschritte](/getting-started/installation/#platform-specific-dependencies) aus.

@end

## Schnellstart

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

Der Assistent wird in Ihrem Browser geöffnet und führt Sie durch die Prüfung der Abhängigkeiten, globale Projektstandards, optionale plattformübergreifende Werkzeuge und Anmeldedaten für die Signierung. Diese Einstellungen können projektübergreifend wiederverwendet werden.

Danach können Sie Ihr erstes Projekt erstellen:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## Funktionsumfang

- **Prüft Abhängigkeiten** – Überprüft Go, npm und plattformspezifische Werkzeuge
- **Konfiguriert Standardeinstellungen** – Angaben zum Autor, Bundle-ID-Präfix und bevorzugte Vorlagen
- **Plattformübergreifende Builds** – Optionale Docker-Einrichtung für Builds auf jedem Hostsystem
- **Codesignierung** – Optionale Einrichtung für macOS, Windows und Linux

Die Konfiguration wird in `$XDG_CONFIG_HOME/wails/defaults.yaml` gespeichert, ersatzweise in `~/.config/wails/defaults.yaml`, und von `wails3 init` verwendet. Sie ist von der `build/config.yml` eines Projekts getrennt.

## Ein Wake-Projekt erstellen

Wenn das Wake-Experiment aktiviert ist, erfasst `wails3 init -ui` Ihre projektspezifischen Build-Einstellungen und speichert sie nach erfolgreicher Erstellung in `wails.yaml`. Folgen Sie dem [Wake-Einstiegstutorial](/experimental/wake/), um beim Erstellen eines Desktopprojekts ein Build-Backend, Zielplattformen und Paketformate auszuwählen. `wails3 setup` behält seine Funktion für globale Einstellungen bei.

Bestehende Wake-Projekte werden in `wails.yaml` konfiguriert. Verwenden Sie den [Leitfaden für YAML-Builds](/experimental/yaml-builds/), um deren Profile und Paketierungseinstellungen zu ändern. Bestehende Taskfile-Projekte müssen vor der Verwendung von Wake die [geprüfte Einwegmigration](/experimental/wake/#migrate-an-existing-project) durchlaufen; die Einrichtung übersetzt deren Taskfiles nicht und führt sie auch nicht aus.

## Signierungseinstellungen wiederverwenden

Signierte Wake-Builds übernehmen fehlende Signierungseinstellungen aus Ihren benutzerspezifischen Standardeinstellungen, einschließlich Zertifikatsdateien außerhalb des Projekts. Explizite Projekteinstellungen haben Vorrang. Private Schlüssel, Passwörter und globale Pfade zu Anmeldedaten werden nicht nach `wails.yaml` kopiert.

Geben Sie das Passwort für global konfigurierte Windows- und Linux-Schlüssel über den vorhandenen Anmeldedatenspeicher oder `WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD` an. Setzen Sie bei einer verschlüsselten Apple-Identität `darwin.signing.credential` auf den Namen einer Umgebungsvariablen, die das zugehörige Passwort enthält.

## Unterbefehle

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Probleme?

1. Führen Sie `wails3 doctor` aus, um Probleme zu diagnostizieren
2. Führen Sie die [manuellen Installationsschritte](/getting-started/installation/#platform-specific-dependencies) aus
3. [Melden Sie das Problem](https://github.com/wailsapp/wails/issues/4904) zusammen mit der Ausgabe von `wails3 doctor`
