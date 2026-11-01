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

Der Assistent wird in Ihrem Browser geöffnet und führt Sie durch die Prüfung der Abhängigkeiten, die Projektstandardwerte und die optionale Einrichtung plattformübergreifender Builds. Die obigen Anweisungen verwenden den standardmäßigen Einrichtungsablauf; der experimentelle Ablauf für Wake-Projekte wird nachfolgend beschrieben.

Danach können Sie Ihr erstes Projekt erstellen:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## Funktionsweise

- **Prüft Abhängigkeiten** – Überprüft Go, npm und die Plattformwerkzeuge
- **Konfiguriert Standardwerte** – Autoreninformationen, Bundle-ID-Präfix und bevorzugte Vorlagen
- **Plattformübergreifende Builds** – Optionale Docker-Einrichtung für Builds auf jedem Hostsystem
- **Codesignierung** – Optionale Einrichtung für macOS, Windows und Linux

Die Konfiguration wird in `$XDG_CONFIG_HOME/wails/defaults.yaml` gespeichert, ersatzweise in `~/.config/wails/defaults.yaml`, und von `wails3 init` verwendet. Sie ist von der Datei `build/config.yml` eines Projekts getrennt.

## Wake-Projekt konfigurieren

Führen Sie die Einrichtung in Ihrem Projektverzeichnis mit aktiviertem [Wake-Experiment](/experimental/yaml-builds/#enable-yaml-builds) aus:

```bash
export WAILS_EXP_USE_WAKE=1
wails3 setup
```

In PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
wails3 setup
```

Dadurch wird ein projektbezogener Ablauf geöffnet:

1. **Projekt** — Geben Sie Projektmetadaten ein, wenn kein natives Manifest vorhanden ist.
2. **Anwendung** — Wählen Sie GUI, CLI oder Server sowie ein Build-Backend oder bereits erstellte ausführbare Dateien aus. Reines Go ist eine ausdrücklich experimentelle Option.
3. **Ziele** — wählen Sie Betriebssysteme und Architekturen aus. iOS und Android bleiben experimentell; konfigurieren Sie deren Signierungseinstellungen in YAML.
4. **Pakete** — Wählen Sie Paketformate aus und legen Sie fest, ob sie signiert werden sollen.
5. **Werkzeuge** — Prüfen Sie die Anforderungen für den ausgewählten Build- und Paketierungsweg.
6. **Signierung** — Konfigurieren Sie bei ausgewählter Signierung die vorhandenen Signierungsoptionen.
7. **Speichern** — Schreiben Sie die Projektkonfiguration und zeigen Sie den auszuführenden Befehl an.

Die Werkzeuganforderungen hängen von Ihrer Auswahl ab. Für die Desktop-Kompilierung mit reinem Go sind weder Docker noch ein C-Compiler erforderlich. Plattformübergreifende CGo-Builds benötigen einen kompatiblen Compiler und die Zielbibliotheken oder eine unterstützte Container-Toolchain. NSIS benötigt weiterhin `makensis`; GUI-AppImages, die auf einem anderen Hostsystem erstellt werden, benötigen einen bereitgestellten Linux-Laufzeitbaum. Exportierte Signierungsidentitäten können die Go-Signierer von Wails verwenden, während Identitäten aus nativen Schlüsselbund- oder Zertifikatspeichern weiterhin Plattformwerkzeuge benötigen. Siehe [Anforderungen für Paketierung und Signierung](/experimental/yaml-builds/#package-and-sign-for-distribution).

Bei bereits erstellten ausführbaren Dateien werden die Go-Kompilierung, die Installation und Builds des Frontends, Bindings und UPX übersprungen. Die ausgewählten Paketierungs- und Signierungswerkzeuge sowie gültige Eingabebinärdateien sind weiterhin erforderlich. Siehe [Pfade für vorgefertigte Eingabedateien](/experimental/yaml-builds/#use-prebuilt-desktop-executables). Diese Auswahl gilt für den nach dem Speichern angezeigten Befehl; sie ist keine dauerhafte YAML-Buildeinstellung.

Mit „Speichern“ wird `wails.yaml` geschrieben, die Anwendungsart, die Frontend-Aktivierung und die ausgewählten Backend-Einstellungen werden aktualisiert und die separate Auswahl von Ziel, Paket und Signierung für `profiles.setup` wird ersetzt. Andere Profile, Frontend-Befehle, Begleitprogramme, Hooks und nicht zugehörige Projekteinstellungen bleiben erhalten. Backend-Einstellungen gelten projektweit, sodass andere Profile, die sie erben, ebenfalls das ausgewählte Backend verwenden. Neue Projektmetadaten verbleiben bis zum Speichern im Assistenten. Fehlende Tools werden gemeldet, damit Sie sie vor dem Build installieren können; beim Speichern einer gültigen Konfiguration werden weder Tools installiert noch ein Build ausgeführt.

Signierte Builds übernehmen fehlende Signierungseinstellungen aus Ihren Standardwerten auf Benutzerebene, einschließlich Zertifikatsdateien außerhalb des Projekts. Explizite Projekteinstellungen haben Vorrang. Die Einrichtung kopiert keine privaten Schlüssel, Passwörter oder globalen Anmeldeinformationspfade nach `wails.yaml`. Stellen Sie für global konfigurierte Windows- und Linux-Schlüssel das Passwort über den vorhandenen Anmeldeinformationsspeicher oder `WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD` bereit. Setzen Sie für eine verschlüsselte Apple-Identität `darwin.signing.credential` auf den Namen einer Umgebungsvariable, die das Passwort enthält.

Das gespeicherte Profil wird mit `wails3 build setup` ausgeführt; die Option für bereits erstellte Dateien fügt `--skip-build` hinzu.

Lassen Sie Wake aktiviert, wenn Sie den angezeigten Befehl ausführen. Ein aktives `wails.yaml` greift bei nicht gesetztem Flag niemals auf Taskfiles zurück. Die Einrichtung konvertiert weder vorhandene Taskfiles noch ältere Metadaten: Verwenden Sie den [geprüften Migrationsablauf](/experimental/yaml-builds/#migrate-an-existing-project), bevor Sie diese Projekte mit Wake konfigurieren.

## Unterbefehle

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Probleme?

1. Führen Sie `wails3 doctor` aus, um Probleme zu diagnostizieren
2. Führen Sie die [manuellen Installationsschritte](/getting-started/installation/#platform-specific-dependencies) aus
3. [Melden Sie das Problem](https://github.com/wailsapp/wails/issues/4904) zusammen mit Ihrer `wails3 doctor`-Ausgabe
