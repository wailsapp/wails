---
title: "macOS-Paketierung"
description: "Ihre Wails-Anwendung für die Verteilung unter macOS paketieren"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## Private macOS-APIs

Wails v3 verwendet standardmäßig öffentliche macOS-APIs. Um Funktionen zu aktivieren, die undokumentierte Apple-APIs benötigen, erstellen Sie Ihre App mit dem einzelnen Go-Build-Tag `private_mac_apis`:

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

Verwenden Sie für einen direkten Go-Build `go build -tags private_mac_apis .` (oder `-tags production,private_mac_apis` für die Produktion). Bestehende Anwendungen, die auf privatem Verhalten beruhen, müssen dieses Tag hinzufügen, um es beizubehalten. Das Tag gilt nur für macOS-Desktop-Builds.

Eine vollständige Übersicht der betroffenen Funktionen und Optionswerte, die genauen Fallbacks öffentlicher Builds, die Zuordnungen der Liquid-Glass-Stile und die Build-Kombinationen für den Inspector finden Sie unter [Private macOS-APIs](/guides/build/private-macos-apis/). Ausschließlich private Operationen bleiben ohne das Tag wirkungslos; die öffentliche Go-API bleibt unverändert.

## Anwendungspaket

Paketieren Sie Ihre App als standardmäßiges macOS-`.app`-Bundle:

```bash
wails3 package GOOS=darwin
```

Dadurch wird `bin/<AppName>.app` mit folgendem Inhalt erstellt:

- Die kompilierte Binärdatei in `Contents/MacOS/`
- Das App-Symbol in `Contents/Resources/` (aus `icons.icns` oder, sofern vorhanden, aus dem Asset-Katalog `Assets.car`)
- `Info.plist` mit App-Metadaten

### Projekte in synchronisierten Ordnern

Einige macOS-Dateianbieter fügen `.app`-Verzeichnissen Finder-Metadaten hinzu. Dadurch lehnt `codesign` das Bundle ab oder dessen Signatur wird nach dem Signieren ungültig. Es reicht nicht, nur die Metadaten zu entfernen, wenn der Anbieter sie erneut hinzufügt.

Die generierten macOS-Build-Tasks prüfen auf dieses Verhalten. Wird es erkannt, verschiebt Wails das generierte Ausgabeverzeichnis des Projekts nach `~/Library/Application Support/Wails/build-output/<project>-<hash>` und ersetzt das ursprüngliche Ausgabeverzeichnis durch einen symbolischen Link. Vorhandene Ausgabedateien bleiben erhalten, und `bin/<AppName>.app` bleibt der für Paketierung, Signierung und Entwicklungsmodus verwendete Pfad. Die CLI gibt beim Verschieben den lokalen Speicherort aus.

Legen Sie mit `WAILS_MACOS_OUTPUT_ROOT` einen anderen lokalen, nicht synchronisierten Speicherort fest. Das Stammverzeichnis für die Verlagerung muss sich auf demselben Dateisystem wie das Projekt befinden. Ein explizit konfiguriertes `BIN_DIR` außerhalb des Projekts wird unverändert verwendet. Bestehende Projekte benötigen das aktualisierte `build/darwin/Taskfile.yml`, um diese Behandlung zu nutzen.

Kopieren Sie für die Verteilung `bin/<AppName>.app` oder erstellen Sie ein DMG. Beim alleinigen Kopieren des symbolischen Links `bin` werden die Build-Artefakte nicht einbezogen. Der direkte Befehl `wails3 tool sign` meldet einen Fehler, wenn sein Eingabeordner erneut Metadaten erhält. Verschieben Sie dieses Bundle vor dem Signieren in einen nicht synchronisierten Ordner.

## Bundle-Ressourcen

`Contents/Resources/` ist der Standardspeicherort für schreibgeschützte Dateien, die mit einer macOS-App ausgeliefert werden. Verwenden Sie ihn für größere Vorlagen, Ausgangsdaten, Medien, Sprachpakete oder andere Inhalte, die bei Bedarf geöffnet und nicht mit `embed` in die ausführbare Go-Datei kompiliert werden sollen.

Wails legt das Anwendungssymbol bereits in diesem Verzeichnis ab. Um eigene Dateien hinzuzufügen, legen Sie diese in einem Quellverzeichnis wie `build/resources/` ab und fügen Sie anschließend dem Task `create:app:bundle` in `build/darwin/Taskfile.yml` einen Kopierschritt hinzu:

```yaml
tasks:
  create:app:bundle:
    cmds:
      # Existing bundle creation commands...
      - |
          if [ -d build/resources ]; then
            cp -R build/resources/. "{{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources/"
          fi
```

Wenn Sie den Task `darwin:run` des Taskfiles verwenden, fügen Sie dessen Task `run` den entsprechenden Befehl mit `{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/` als Ziel hinzu.

### Ressourcen aus Go lesen

Importieren Sie das macOS-Plattformpaket:

```go
import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/mac"
)
```

Verwenden Sie für kleine Dateien `LoadResource`:

```go
func loadSplash() ([]byte, error) {
	return mac.LoadResource("images/splash.png")
}
```

Verwenden Sie für größere Dateien `ResourceFS`. Die Funktion gibt ein in `Contents/Resources` verwurzeltes `io/fs.FS` zurück, sodass Aufrufer eine Ressource öffnen und streamen können, ohne sie zuerst vollständig in ein Go-Byte-Slice zu laden:

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

Ressourcennamen sind durch Schrägstriche getrennte Pfade relativ zu `Contents/Resources`. `ResourceFS` und `LoadResource` geben `mac.ErrNotInAppBundle` zurück, sofern die ausführbare Datei nicht aus `.app/Contents/MacOS` ausgeführt wird.

Behandeln Sie Bundle-Ressourcen als unveränderlich. Änderungen an Dateien innerhalb einer signierten Anwendung machen deren Codesignatur ungültig. Speichern Sie heruntergeladene, generierte oder vom Benutzer bearbeitbare Daten stattdessen im Application-Support-Verzeichnis des Benutzers.

### Universelle Binärdatei

Erstellen Sie einen Build für Macs mit Apple Silicon und Intel-Prozessoren:

```bash
wails3 task darwin:package:universal
```

Dadurch wird eine einzelne `.app` erstellt, die auf beiden Architekturen nativ ausgeführt wird. Universelle Binärdateien können auf jeder Plattform erstellt werden – unter Linux und Windows wird `wails3 tool lipo` automatisch verwendet.

## Bundle anpassen

Bearbeiten Sie `build/darwin/Info.plist`, um Folgendes anzupassen:

- Bundle-Kennung (`CFBundleIdentifier`)
- App-Name und -Version
- MacOS-Mindestversion
- Dateizuordnungen
- URL-Schemata

Das App-Symbol wird aus Assets im Verzeichnis `build/` generiert. Verwenden Sie den Task `generate:icons`:

```bash
wails3 task common:generate:icons
```

Dabei erzeugt `build/appicon.png` die Dateien `darwin/icons.icns` und `windows/icon.ico`. Unter macOS können Sie auch `build/appicon.icon` (Icon-Composer-Format) bereitstellen: Der Task übergibt `-iconcomposerinput appicon.icon -macassetdir darwin`, wodurch aus der Datei `.icon` die Dateien `Assets.car` und `darwin/icons.icns` erzeugt werden (auf anderen Plattformen als macOS wird dieser Schritt übersprungen). Wenn `Assets.car` vorhanden ist, führen Sie den Task `update:build-assets` aus, damit `Info.plist` und `CFBundleIconName` entsprechend aktualisiert werden:

```bash
wails3 task common:update:build-assets
```

So führen Sie den Symbolbefehl manuell aus dem Verzeichnis `build/` aus:

```bash
cd build
wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico -iconcomposerinput appicon.icon -macassetdir darwin
```

## Codesignierung

Signieren Sie Ihre App für die Verteilung:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=darwin

# Or using the task directly
wails3 task darwin:sign
```

Konfigurieren Sie die Signierung in `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

### Notarisierung

Für Apps, die außerhalb des Mac App Store verteilt werden, verlangt Apple eine Notarisierung:

```bash
wails3 task darwin:sign:notarize
```

Speichern Sie zunächst Ihre Anmeldedaten. Führen Sie entweder den interaktiven Assistenten (`wails3 signing`) aus oder rufen Sie `notarytool` direkt auf:

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

Konfigurieren Sie dies in `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

Weitere Informationen finden Sie unter [Anwendungen signieren](/guides/build/signing/).

## DMG-Installationsprogramm

Die mit Wails ausgelieferte Vorlage 3 stellt `wails3 task darwin:package:dmg` bereit. Sie erstellt zunächst `.app` und anschließend mithilfe der DMG-Bibliothek ein gestaltetes DMG. Standardmäßig verwendet das DMG einen Verlaufshintergrund im Wails-Branding mit dem roten Drachenlogo und der WAILS-Wortmarke.

```bash
wails3 task darwin:package:dmg
```

Der untergeordnete Task `darwin:create:dmg` erstellt aus einem vorhandenen `.app`-Bundle ein DMG und kann direkt im Taskfile konfiguriert werden:

```yaml
vars:
  # These are the template defaults; override them when needed.
  DMG_BACKGROUND: build/darwin/dmg-background.png
  DMG_VOLUME_ICON: build/darwin/icons.icns
  DMG_FILE_ICON: build/darwin/dmg-file-icon.icns
  DMG_WINDOW_WIDTH: 540
  DMG_WINDOW_HEIGHT: 380
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

### Standardlayout

Das generierte DMG enthält:

- Das Anwendungspaket auf der linken Seite
- Einen `Applications`-Link auf der rechten Seite
- Ein Finder-Fenster mit einer Größe von 540 × 380 Pixeln
- Eine Symbolgröße von 96 Punkt mit Beschriftungen unter jedem Symbol
- Den Hintergrund im Wails-Branding aus `build/darwin/dmg-background.png`

Die Symbole der Anwendung und von `Applications` werden relativ zu den konfigurierten Fensterabmessungen positioniert. Dadurch bleiben die Abstände des standardmäßigen Layouts mit zwei Symbolen proportional, wenn Sie `DMG_WINDOW_WIDTH` oder `DMG_WINDOW_HEIGHT` ändern. Das beste Ergebnis erzielen Sie mit einem Hintergrundbild, dessen Pixelabmessungen denen des Finder-Fensters entsprechen.

### DMG-Assets ersetzen

Die generierten Dateien unter `build/darwin/` sind normale Projekt-Assets und können ersetzt werden:

- `DMG_BACKGROUND` steuert das Bild, das hinter dem Inhalt des Finder-Fensters angezeigt wird.
- `DMG_VOLUME_ICON` steuert das für das eingebundene Volume angezeigte Symbol.
- `DMG_FILE_ICON` steuert das Symbol, das im Finder für die resultierende Datei `.dmg` angezeigt wird.

Das Volume-Symbol und das Symbol der DMG-Datei sind separate Ressourcen. Wenn Sie das Anwendungssymbol ersetzen, werden diese beiden Symbole nicht automatisch ersetzt.

### Zusätzliche Dateien hinzufügen

Verwenden Sie `DMG_FILES`, um Installationsskripte, Versionshinweise, Lizenzdateien oder andere Ressourcen neben der Anwendung einzubinden. Der Wert ist eine kommagetrennte Liste von `name=path`-Paaren:

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

Der Name vor `=` ist der im DMG angezeigte Dateiname. Der Pfad nach `=` verweist auf die Quelldatei im Projekt. Führende und nachgestellte Leerzeichen werden ignoriert.

Jeder angezeigte Name muss eindeutig sein. Zusätzliche Dateien können keine bereits vom Paketierungswerkzeug erstellten Einträge ersetzen, einschließlich des Anwendungspakets oder des Eintrags `Applications`. Bei einem Namenskonflikt schlägt die Paketierung mit einer Fehlermeldung fehl, anstatt ein beschädigtes DMG zu erzeugen.

@note{type="note"}
Das Erstellen von DMGs wird nur unter macOS unterstützt, da dafür die Werkzeuge von macOS für Festplattenabbilder und Finder verwendet werden. Kreuzkompilierte `.app`-Pakete können auch auf anderen Systemen erstellt werden, das endgültige DMG muss jedoch auf einem Mac erzeugt werden.

@end

## Fehlerbehebung

### „App ist beschädigt und kann nicht geöffnet werden“

Die App ist nicht signiert. Signieren Sie sie entweder mit einem Developer-ID-Zertifikat, oder lassen Sie die Benutzer Gatekeeper umgehen:

```bash
xattr -cr /path/to/YourApp.app
```

### Notarisierung schlägt fehl

Häufige Probleme:

- **Ungültige Anmeldedaten**: Führen Sie `xcrun notarytool store-credentials` (oder `wails3 signing`) erneut aus
- **Gehärtete Laufzeit erforderlich**: Stellen Sie sicher, dass die Berechtigungen bei Bedarf `com.apple.security.cs.allow-unsigned-executable-memory` enthalten
- **Fehlender Zeitstempel**: Der Signierungsprozess sollte automatisch einen Zeitstempel hinzufügen

### Kreuzkompilierte App wird nicht ausgeführt

Kreuzkompilierte macOS-Binärdateien sind nicht signiert. Übertragen Sie sie auf einen Mac und signieren Sie sie vor dem Testen:

```bash
codesign --force --deep --sign - YourApp.app
```
