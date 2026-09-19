---
title: "Codesignierung"
description: "Anleitung zum Signieren Ihrer Wails-Anwendungen auf allen Plattformen"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## Codesignierung Ihrer Anwendung

Diese Anleitung beschreibt, wie Sie Ihre Wails-Anwendungen für macOS, Windows und Linux signieren. Wails v3 bietet integrierte CLI-Werkzeuge für Codesignierung, Notarisierung und die Verwaltung von PGP-Schlüsseln.

- **macOS** – Signieren und notarisieren Sie Ihre macOS-Anwendungen
- **Windows** – Signieren Sie Ihre ausführbaren Windows-Dateien und -Pakete
- **Linux** – Signieren Sie DEB- und RPM-Pakete mit PGP-Schlüsseln

## Plattformübergreifende Signierungsmatrix

Diese Matrix zeigt, was Sie von jeder Quellplattform aus signieren können:

| Zielformat | Unter Windows | Unter macOS | Unter Linux |
| --- | :---: | :---: | :---: |
| Windows EXE/MSI | ✅ | ✅ | ✅ |
| macOS-.app-Bundle | ❌ | ✅ | ❌ |
| macOS-Notarisierung | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
Windows- und Linux-Pakete können von **jeder Plattform** aus signiert werden. Aufgrund der Anforderungen der Apple-Werkzeuge ist für die macOS-Signierung ein Mac erforderlich.

@end

### Signierungs-Backends

Wails wählt automatisch das beste verfügbare Signierungs-Backend aus:

| Plattform | Natives Backend | Plattformübergreifendes Backend |
| --- | --- | --- |
| Windows | `signtool.exe` (Windows SDK) | Integriert |
| macOS | `codesign` (Xcode) | Nicht verfügbar |
| Linux | Nicht zutreffend | Integriert |

Bei der Ausführung auf der nativen Plattform verwendet Wails für maximale Kompatibilität die nativen Werkzeuge. Bei der Cross-Kompilierung verwendet es die integrierte Signierungsunterstützung.

## Schnellstart

Am schnellsten konfigurieren Sie die Signierung mit dem Einrichtungsassistenten:

```bash
wails3 setup
```

Dadurch wird eine gemeinsame Signierungskonfiguration in `~/.config/wails/defaults.yaml` gespeichert, die **beim Signieren auf jeder Plattform berücksichtigt wird** (siehe [Konfigurationsrangfolge](#konfigurationsrangfolge)). Der Signierungsschritt:

- Erkennt – **von jedem Host aus** –, welche Signierungswerkzeuge für die gewünschte Zielplattform installiert sind, und zeigt den Installationsbefehl für Ihr Betriebssystem an (z. B. `brew install gnupg`, `sudo apt install osslsigncode`, `winget install GnuPG.Gpg4win`).
- Listet unter macOS die Developer-ID-Zertifikate aus Ihrem Schlüsselbund auf.
- Listet für Linux Ihre GPG-Schlüssel auf und kann einen neuen Schlüssel **erzeugen und exportieren**.
- Kann für Windows über OpenSSL ein **selbstsigniertes Zertifikat erzeugen** (zu Testzwecken).
- **Speichert Passwörter sicher im Schlüsselbund Ihres Systems** (nicht in Taskfiles).

@note{type="tip"}
Passwörter werden im nativen Anmeldedatenspeicher Ihres Systems gespeichert (macOS-Schlüsselbund, Windows-Anmeldeinformationsverwaltung oder Linux Secret Service). Dadurch ist Ihre Signierungskonfiguration sicher und funktioniert projektübergreifend für alle Ihre Wails-Projekte.

@end

## Konfigurationsrangfolge

Wenn Sie eine Signierungsaufgabe ausführen, wird jede Signierungsoption in dieser Reihenfolge aufgelöst (der erste Treffer gilt):

1. Ein explizit an `wails3 tool sign` übergebenes Flag (z. B. `--pgp-key`, `--certificate`, `--identity`).
2. Die entsprechende **Taskfile-Variable des Projekts** (`PGP_KEY`, `SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`, `SIGN_IDENTITY`, …).
3. Die **globale Konfiguration** in `~/.config/wails/defaults.yaml` (geschrieben von `wails3 setup`).

Das bedeutet, dass die Taskfile-Variablen **optionale Überschreibungen** sind: Ist eine Variable nicht gesetzt, wird der global konfigurierte Schlüssel, das global konfigurierte Zertifikat oder die global konfigurierte Identität verwendet. Liefert keine der drei Quellen einen Wert, meldet der Signierungsbefehl einen eindeutigen Fehler mit Hinweisen zur Konfiguration.

## Projektbezogene Konfiguration

Um die Signierung nur für ein einzelnes Projekt statt global zu konfigurieren, führen Sie den projektbezogenen Assistenten im Projekt aus. Er schreibt `vars` in die `build/<platform>/Taskfile.yml`-Dateien dieses Projekts:

```bash
wails3 setup signing                                   # all detected platforms
wails3 setup signing --platform windows --platform linux
```

### Manuelle Konfiguration

Alternativ können Sie die plattformspezifischen Taskfiles manuell bearbeiten. Bearbeiten Sie den Abschnitt `vars` am Anfang jeder Datei:

@tabs
[macOS]
Bearbeiten Sie `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  # ENTITLEMENTS: "build/darwin/entitlements.plist"
```

Führen Sie anschließend Folgendes aus:

```bash
wails3 task darwin:sign           # Sign only
wails3 task darwin:sign:notarize  # Sign and notarize
```

[Windows]
Bearbeiten Sie `build/windows/Taskfile.yml`:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

Das Passwort wird aus dem Schlüsselbund des Systems abgerufen (führen Sie zur Konfiguration `wails3 setup signing` aus).

Führen Sie anschließend Folgendes aus:

```bash
wails3 task windows:sign           # Sign executable
wails3 task windows:sign:installer # Sign NSIS installer
```

[Linux]
Bearbeiten Sie `build/linux/Taskfile.yml`:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

Das Passwort wird aus dem Schlüsselbund des Systems abgerufen (führen Sie zur Konfiguration `wails3 setup signing` aus).

Führen Sie anschließend Folgendes aus:

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

Sie können den Signierungsstatus auch direkt über das System prüfen:

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

Verwenden Sie für die interaktive Einrichtung der Signierungskonfiguration aller Plattformen den Assistenten:

```bash
wails3 setup signing
```

## Codesignierung unter macOS

### Voraussetzungen

- Apple-Developer-Account ($99/Jahr)
- Developer ID Application-Zertifikat
- Installierte Xcode Command Line Tools

### Signierungsidentitäten

Prüfen Sie die verfügbaren Signierungsidentitäten:

```bash
security find-identity -v -p codesigning
```

Ausgabe:

```
Found 2 signing identities:

  Developer ID Application: Your Company (ABCD1234) [valid]
    Hash: ABC123DEF456...

  Apple Development: your@email.com (XYZ789) [valid]
    Hash: DEF789ABC123...
```

@note{type="tip"}
Für die Verteilung außerhalb des App Store benötigen Sie ein **Developer ID Application**-Zertifikat.

@end

### Konfiguration

Bearbeiten Sie `build/darwin/Taskfile.yml` und legen Sie die Signierungsvariablen fest:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| Variable | Erforderlich | Beschreibung |
| --- | --- | --- |
| `SIGN_IDENTITY` | Ja | Ihre Developer ID (z. B. „Developer ID Application: Ihr Unternehmen (TEAMID)“) |
| `KEYCHAIN_PROFILE` | Für die Beglaubigung | Name des Schlüsselbundprofils mit den gespeicherten Anmeldedaten |
| `ENTITLEMENTS` | Nein | Pfad zur Entitlements-Datei |

Führen Sie anschließend Folgendes aus:

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### Entitlements

Entitlements steuern, auf welche Funktionen Ihre App zugreifen darf. Wails-Apps benötigen für Entwicklung und Produktion in der Regel unterschiedliche Entitlements:

- **Entwicklung**: Erfordert Entitlements für JIT, nicht signierten Speicher und Debugging
- **Produktion**: Minimale Entitlements (nur Netzwerkzugriff)

Verwenden Sie den interaktiven Einrichtungsassistenten, um beide Dateien zu erzeugen:

```bash
wails3 setup entitlements
```

Dadurch werden folgende Dateien erstellt:

- `build/darwin/entitlements.dev.plist` – Für Entwicklungs-Builds
- `build/darwin/entitlements.plist` – Für Produktions-/signierte Builds

**Verfügbare Voreinstellungen:**\

| Voreinstellung | Beschreibung |
| --- | --- |
| Entwicklung | JIT, nicht signierter Speicher, Debugging, Netzwerk |
| Produktion | Nur Netzwerk (minimal, höchste Sicherheit) |
| Beide | Erstellt sowohl Entwicklungs- als auch Produktionsdateien (empfohlen) |
| App Store | Sandbox mit Netzwerk- und Dateizugriff aktiviert |
| Benutzerdefiniert | Einzelne Entitlements auswählen |

@note{type="note"}
Der Task `run` im darwin-Taskfile verwendet automatisch `entitlements.dev.plist`. Die Tasks `sign` verwenden für Produktions-Builds `entitlements.plist`.

@end

Legen Sie anschließend in den Variablen Ihres Taskfiles `ENTITLEMENTS` so fest, dass es auf die entsprechende Datei verweist.

### Beglaubigung

Apple verlangt, dass alle verteilten Apps beglaubigt werden.

@steps
### **Speichern Sie Ihre Anmeldedaten im Schlüsselbund** (einmalige Einrichtung). Führen Sie entweder `wails3 setup signing` aus (der Befehl fragt diese Werte ab und ruft intern `notarytool` auf), oder rufen Sie `notarytool` direkt auf:
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **Legen Sie KEYCHAIN_PROFILE in Ihrem Taskfile** so fest, dass es dem oben angegebenen Profilnamen entspricht.
### **Signieren und beglaubigen Sie Ihre App**:
```bash
wails3 task darwin:sign:notarize
```

### **Überprüfen Sie die Beglaubigung**:
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
Die Beglaubigung dauert in der Regel 1-2 Minuten. Das Ticket wird automatisch an Ihre App angeheftet.

@end

## Codesignierung unter Windows

### Voraussetzungen

- Codesignierungszertifikat (von DigiCert, Sectigo usw.)
- Für die native Signierung unter Windows: installiertes Windows SDK (für `signtool.exe`)
- Für plattformübergreifendes Signieren unter macOS/Linux: [`osslsigncode`](https://github.com/mtrojnar/osslsigncode) (im Signierschritt `wails3 setup` wird der hostspezifische Installationsbefehl angezeigt)

### Selbstsigniertes Zertifikat erzeugen (Testzwecke)

Wenn Sie lediglich die Signierpipeline testen möchten, führen Sie `wails3 setup` aus, öffnen Sie im Signierschritt die Registerkarte **Windows** und wählen Sie **Selbstsigniertes Zertifikat erzeugen**. Dabei wird mit OpenSSL eine Codesignatur-`.pfx` erstellt und ihr Pfad in der globalen Konfiguration gespeichert.

@note{type="caution"}
Selbstsignierte Zertifikate sind **ausschließlich für Tests und die interne Verteilung vorgesehen** — bei Endbenutzern lösen sie SmartScreen-Warnungen aus. Für öffentliche Releases ist ein Zertifikat von einer vertrauenswürdigen Zertifizierungsstelle erforderlich.

@end

### Konfiguration

Bearbeiten Sie `build/windows/Taskfile.yml` und legen Sie die Signiervariablen fest:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| Variable | Erforderlich | Beschreibung |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | Überschreibung | Pfad zur Zertifikatsdatei im Format .pfx/.p12 (greift auf die globale Konfiguration zurück, wenn nicht festgelegt) |
| `SIGN_THUMBPRINT` | Überschreibung | Fingerabdruck des Zertifikats im Windows-Zertifikatspeicher (Alternative zu `SIGN_CERTIFICATE`) |
| `TIMESTAMP_SERVER` | Nein | URL des Zeitstempelservers (Standard: http://timestamp.digicert.com) |

@note{type="note"}
Diese Variablen sind optionale **Überschreibungen**. Sind sie nicht festgelegt, wird das über `wails3 setup` global konfigurierte Zertifikat verwendet. Siehe [Konfigurationsrangfolge](#konfigurationsrangfolge).

@end

@note{type="note"}
Das Zertifikatspasswort wird im Schlüsselbund Ihres Systems und nicht im Taskfile gespeichert. Führen Sie zum Konfigurieren `wails3 setup signing` aus oder setzen Sie in CI die Umgebungsvariable `WAILS_WINDOWS_CERT_PASSWORD`.

@end

Führen Sie anschließend Folgendes aus:

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### Plattformübergreifendes Signieren

Ausführbare Windows-Dateien können auf jeder Plattform signiert werden. Dieselbe Taskfile-Konfiguration und dieselben Befehle funktionieren unter macOS und Linux.

### Unterstützte Windows-Formate

| Format | Erweiterung | Hinweise |
| --- | --- | --- |
| Ausführbare Dateien | .exe | Standardmäßige PE-Signierung |
| Installationsprogramme | .msi | Windows-Installer-Pakete |
| App-Pakete | .msix, .appx | Moderne Windows-Apps |

## Signieren von Linux-Paketen

Linux-Pakete (DEB und RPM) werden mit PGP-/GPG-Schlüsseln signiert. Anders als die Codesignierung unter Windows und macOS weist das Signieren von Linux-Paketen nach, dass das Paket aus einer vertrauenswürdigen Quelle stammt, und nicht, dass das Betriebssystem dem Code vertraut.

### Voraussetzungen

- PGP-Schlüsselpaar (kann mit Wails erzeugt werden)

### PGP-Schlüssel erzeugen

Am einfachsten geht dies mit dem Einrichtungsassistenten: Führen Sie `wails3 setup` aus, öffnen Sie im Signierschritt die Registerkarte **Linux** und wählen Sie **Neuen GPG-Schlüssel erstellen**. Der Assistent:

- erzeugt einen RSA-4096-Schlüssel in Ihrem GPG-Schlüsselbund (lassen Sie die Passphrase für einen unbeaufsichtigt verwendbaren, CI-geeigneten Schlüssel leer),
- **exportiert ihn nach `~/.wails/signing/<keyid>.asc`** (die Datei, mit der ein Build signiert wird) und
- speichert sowohl die Schlüssel-ID als auch den Exportpfad in `~/.config/wails/defaults.yaml`, sodass der Schlüssel beim Signieren automatisch verwendet wird.

@note{type="note"}
Wenn ein Schlüssel **nur anhand seiner ID** konfiguriert wurde (z. B. bevor der Assistent Schlüssel automatisch exportierte), exportiert ihn die Registerkarte „Linux“ beim nächsten Öffnen in eine Datei und trägt den Pfad ein — manuelle Schritte sind nicht erforderlich. Builds signieren mit einer Schlüssel<em>datei</em>, daher ist der Pfad entscheidend.

@end

Sie können dies mit `gpg` auch manuell durchführen:

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

Empfohlene Einstellungen: RSA mit 4096 Bit, 1 Jahr Gültigkeitsdauer und Schutz durch ein starkes Passwort.

@note{type="caution"}
Bewahren Sie Ihren privaten Schlüssel sicher auf! Speichern Sie ihn verschlüsselt und sichern Sie ihn zuverlässig.

@end

### Konfiguration

Bearbeiten Sie `build/linux/Taskfile.yml` und legen Sie die Signiervariablen fest:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| Variable | Erforderlich | Beschreibung |
| --- | --- | --- |
| `PGP_KEY` | Überschreibung | Pfad zur exportierten privaten PGP-Schlüsseldatei (greift auf den global über `wails3 setup` konfigurierten Schlüssel zurück, wenn nicht festgelegt) |
| `SIGN_ROLE` | Nein | DEB-Signaturrolle (Standard: builder) |

@note{type="note"}
Das Passwort des PGP-Schlüssels wird im Schlüsselbund Ihres Systems gespeichert, nicht im Taskfile. Führen Sie `wails3 setup signing` aus, um es zu konfigurieren, oder legen Sie in CI die Umgebungsvariable `WAILS_PGP_PASSWORD` fest.

@end

Führen Sie anschließend Folgendes aus:

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### DEB-Signaturrollen

Für DEB-Pakete können Sie die Signaturrolle über `SIGN_ROLE` angeben:

- `origin`: Signatur des Paketursprungs
- `maint`: Signatur des Paketbetreuers
- `archive`: Signatur des Archivbetreuers
- `builder`: Signatur des Paketerstellers (Standard)

### Plattformübergreifendes Signieren

Linux-Pakete können auf jeder Plattform signiert werden. Dieselbe Taskfile-Konfiguration und dieselben Befehle funktionieren unter Windows und macOS.

### Schlüsselinformationen anzeigen

```bash
gpg --show-keys signing-key.asc
```

Ausgabe:

```
pub   rsa4096 2024-01-15 [SC] [expires: 2025-01-15]
      1234 5678 90AB CDEF 1234 5678 90AB CDEF 1234 5678
uid                      Your Name <your@email.com>
```

### Linux-Pakete überprüfen

```bash
# Verify DEB signature
dpkg-sig --verify myapp_1.0.0_amd64.deb

# Verify RPM signature
rpm --checksig myapp-1.0.0.x86_64.rpm
```

### Öffentlichen Schlüssel verteilen

Benutzer benötigen Ihren öffentlichen Schlüssel, um Pakete zu überprüfen:

```bash
# Export public key for distribution
gpg --armor --export "your@email.com" > myapp-signing.pub.asc

# Users can import it:
# For DEB (apt):
sudo apt-key add myapp-signing.pub.asc
# Or for modern apt:
sudo cp myapp-signing.pub.asc /etc/apt/trusted.gpg.d/

# For RPM:
sudo rpm --import myapp-signing.pub.asc
```

## GitHub-Actions-Integration

In CI-Umgebungen werden Passwörter über Umgebungsvariablen statt über den Schlüsselbund des Systems bereitgestellt:

| Umgebungsvariable | Beschreibung |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Passwort des Windows-Zertifikats |
| `WAILS_PGP_PASSWORD` | Passwort des PGP-Schlüssels für Linux-Pakete |

Sie können Taskfile-Variablen auch direkt übergeben:

```bash
wails3 task darwin:sign SIGN_IDENTITY="$SIGN_IDENTITY" KEYCHAIN_PROFILE="$KEYCHAIN_PROFILE"
```

### macOS-Workflow

```yaml
name: Build and Sign macOS

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.MACOS_CERTIFICATE }}
          CERTIFICATE_PASSWORD: ${{ secrets.MACOS_CERTIFICATE_PASSWORD }}
        run: |
          echo $CERTIFICATE_BASE64 | base64 --decode > certificate.p12
          security create-keychain -p "" build.keychain
          security default-keychain -s build.keychain
          security unlock-keychain -p "" build.keychain
          security import certificate.p12 -k build.keychain -P "$CERTIFICATE_PASSWORD" -T /usr/bin/codesign
          security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "" build.keychain

      - name: Store Notarization Credentials
        env:
          APPLE_ID: ${{ secrets.APPLE_ID }}
          APPLE_TEAM_ID: ${{ secrets.APPLE_TEAM_ID }}
          APPLE_APP_PASSWORD: ${{ secrets.APPLE_APP_PASSWORD }}
        run: |
          xcrun notarytool store-credentials "notarize-profile" \
            --apple-id "$APPLE_ID" \
            --team-id "$APPLE_TEAM_ID" \
            --password "$APPLE_APP_PASSWORD"

      - name: Build, Sign, and Notarize
        env:
          SIGN_IDENTITY: ${{ secrets.MACOS_SIGN_IDENTITY }}
        run: |
          wails3 task darwin:sign:notarize \
            SIGN_IDENTITY="$SIGN_IDENTITY" \
            KEYCHAIN_PROFILE="notarize-profile"

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-macOS
          path: bin/*.app
```

### Windows-Workflow

```yaml
name: Build and Sign Windows

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
        run: |
          $certBytes = [Convert]::FromBase64String($env:CERTIFICATE_BASE64)
          [IO.File]::WriteAllBytes("certificate.pfx", $certBytes)

      - name: Build and Sign
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-Windows
          path: bin/*.exe
```

### Plattformübergreifender Workflow (Linux-Runner)

Signieren Sie Windows- und Linux-Pakete mit einem einzigen Linux-Runner:

```yaml
name: Build and Sign (Cross-Platform)

on:
  push:
    tags: ['v*']

jobs:
  build-and-sign:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

      - name: Install Build Dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y nsis rpm

      # Import certificates
      - name: Import Certificates
        env:
          WINDOWS_CERT_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
          PGP_KEY_BASE64: ${{ secrets.PGP_PRIVATE_KEY }}
        run: |
          echo "$WINDOWS_CERT_BASE64" | base64 -d > certificate.pfx
          echo "$PGP_KEY_BASE64" | base64 -d > signing-key.asc

      # Build and sign Windows
      - name: Build and Sign Windows
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      # Build and sign Linux packages
      - name: Build and Sign Linux Packages
        env:
          WAILS_PGP_PASSWORD: ${{ secrets.PGP_PASSWORD }}
        run: |
          wails3 task linux:sign:packages PGP_KEY=signing-key.asc

      # Cleanup secrets
      - name: Cleanup
        if: always()
        run: rm -f certificate.pfx signing-key.asc

      - name: Upload Artifacts
        uses: actions/upload-artifact@v4
        with:
          name: signed-binaries
          path: |
            bin/*.exe
            bin/*.deb
            bin/*.rpm
```

@note{type="note"}
Ein Linux-Runner für das plattformübergreifende Signieren vereinfacht CI/CD, da keine separaten Windows-Runner erforderlich sind. Für das Signieren unter macOS ist aufgrund der Anforderungen von Apples nativen Werkzeugen weiterhin ein macOS-Runner erforderlich.

@end

## CLI-Referenz

### wails3 setup signing

Interaktiver Assistent zum Konfigurieren der Signierung für Ihr Projekt.

```bash
wails3 setup signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

Der Assistent führt Sie durch folgende Schritte:

- **macOS**: Auswählen eines Developer-ID-Zertifikats und Konfigurieren der Anmeldedaten für die Beglaubigung (ruft `xcrun notarytool store-credentials` auf).
- **Windows**: Auswählen zwischen Zertifikatsdatei und Fingerabdruck sowie Festlegen von Passwort und Zeitstempelserver.
- **Linux**: Verwenden eines vorhandenen PGP-Schlüssels oder Erzeugen eines neuen Schlüssels (ruft `gpg` auf) sowie Konfigurieren der Signaturrolle.

### wails3 setup entitlements

Interaktiver Assistent zum Konfigurieren von macOS-Berechtigungen.

```bash
wails3 setup entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**Voreinstellungen:**

- **Entwicklung**: Erstellt `entitlements.dev.plist` mit Berechtigungen für JIT, Debugging und Netzwerkzugriff
- **Produktion**: Erstellt `entitlements.plist` mit minimalen Berechtigungen
- **Beide**: Erstellt beide Dateien (empfohlen)
- **App Store**: Erstellt Sandbox-Berechtigungen für den Mac App Store
- **Benutzerdefiniert**: Einzelne Berechtigungen und Zieldatei auswählen

### wails3 sign

Signiert Binärdateien und Pakete für die aktuelle oder angegebene Plattform. Dies ist ein Wrapper, der die entsprechende plattformspezifische Signatur-Task aufruft.

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

Dadurch wird die entsprechende Task `<platform>:sign` ausgeführt, die die Signaturkonfiguration aus Ihrem Taskfile verwendet.

### wails3 tool sign

Low-Level-Befehl zum direkten Signieren einer bestimmten Datei. Wird intern von den Taskfiles verwendet.

```bash
wails3 tool sign [flags]
```

**Allgemeine Optionen:**\

| Option | Beschreibung |
| --- | --- |
| `--input` | Pfad zur zu signierenden Datei |
| `--output` | Ausgabepfad (optional; standardmäßig direkte Änderung) |
| `--verbose` | Ausführliche Ausgabe aktivieren |

**Flags für Windows/macOS:**\

| Flag | Beschreibung |
| --- | --- |
| `--certificate` | Pfad zum PKCS#12-Zertifikat (.pfx/.p12) |
| `--password` | Zertifikatspasswort |
| `--timestamp` | URL des Zeitstempelservers |

**macOS-spezifische Flags:**\

| Flag | Beschreibung |
| --- | --- |
| `--identity` | Signierungsidentität ('-' für Ad-hoc verwenden) |
| `--entitlements` | Pfad zur Entitlements-Plist-Datei |
| `--hardened-runtime` | Hardened Runtime aktivieren (Standard: true) |
| `--notarize` | Zur Beglaubigung einreichen |
| `--keychain-profile` | Schlüsselbundprofil für die Beglaubigung |

**Windows-spezifische Flags:**\

| Flag | Beschreibung |
| --- | --- |
| `--thumbprint` | Fingerabdruck des Zertifikats im Windows-Zertifikatspeicher |

**Linux-spezifische Flags:**\

| Flag | Beschreibung |
| --- | --- |
| `--pgp-key` | Pfad zum privaten PGP-Schlüssel |
| `--pgp-password` | Passwort des PGP-Schlüssels |
| `--role` | DEB-Signierungsrolle (origin/maint/archive/builder) |

### Signierungsstatus prüfen (native Werkzeuge)

In v3 gibt es **keinen** `wails3 signing`-Befehl. Verwenden Sie zum Prüfen des Signierungsstatus direkt die nativen Werkzeuge:

| Aufgabe | Befehl |
| --- | --- |
| Identitäten für die macOS-Codesignierung auflisten | `security find-identity -v -p codesigning` |
| Anmeldedaten für die Beglaubigung speichern | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| PGP-Schlüsseldatei prüfen | `gpg --show-keys <key.asc>` |
| PGP-Schlüsselpaar erzeugen | `gpg --full-generate-key` |
| Öffentlichen Schlüssel exportieren | `gpg --armor --export <email>` |

## Fehlerbehebung

### macOS-Probleme

**„Kein Developer-ID-Zertifikat gefunden“**

- Stellen Sie sicher, dass Ihr Zertifikat im Schlüsselbund installiert ist
- Prüfen Sie mit `security find-identity -v -p codesigning`, dass es nicht abgelaufen ist.
- Stellen Sie sicher, dass Sie über ein „Developer ID Application“-Zertifikat verfügen (nicht nur über „Apple Development“)

**„Beglaubigung fehlgeschlagen“**

- Prüfen Sie das Beglaubigungsprotokoll: `xcrun notarytool log <submission-id> --keychain-profile <profile>`
- Stellen Sie sicher, dass die Hardened Runtime aktiviert ist
- Vergewissern Sie sich, dass Ihre App keine unsignierten Binärdateien enthält

**„Codesign fehlgeschlagen“**

- Stellen Sie sicher, dass der Schlüsselbund entsperrt ist: `security unlock-keychain`
- Prüfen Sie die Dateiberechtigungen des App-Bundles

### Windows-Probleme

**„Zertifikat nicht gefunden“**

- Überprüfen Sie, ob der Zertifikatspfad korrekt ist
- Prüfen Sie das Zertifikatspasswort
- Stellen Sie sicher, dass das Zertifikat gültig ist (nicht abgelaufen oder widerrufen)

**„Fehler des Zeitstempelservers“**

- Versuchen Sie es mit einem anderen Zeitstempelserver:
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Linux-Probleme

**„Ungültiger PGP-Schlüssel“**

- Stellen Sie sicher, dass die Schlüsseldatei im ASCII-armierten Format vorliegt
- Stellen Sie mit `gpg --show-keys <key.asc>` sicher, dass der Schlüssel nicht abgelaufen ist
- Überprüfen Sie, ob das Passwort korrekt ist

**„Signaturprüfung fehlgeschlagen“**

- Stellen Sie sicher, dass der öffentliche Schlüssel ordnungsgemäß importiert wurde
- Stellen Sie sicher, dass das Paket nach dem Signieren nicht verändert wurde

## Weitere Ressourcen

### Offizielle Dokumentation

- [Apple-Leitfaden zur Codesignierung](https://developer.apple.com/support/code-signing/)
- [Apple-Dokumentation zur Notarisierung](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Microsoft-Codesignierung](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Signieren von Debian-Paketen](https://wiki.debian.org/SecureApt)
- [Signieren von RPM-Paketen](https://rpm-software-management.github.io/rpm/manual/signatures.html)
