---
title: "Codesignierung"
description: "Leitfaden zum Signieren Ihrer Wails-Anwendungen auf allen Plattformen"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## Codesignierung Ihrer Anwendung

Dieser Leitfaden beschreibt, wie Sie Ihre Wails-Anwendungen für macOS, Windows und Linux signieren. Wails v3 bietet integrierte CLI-Werkzeuge für Codesignierung, Notarisierung und die Verwaltung von PGP-Schlüsseln.

- **macOS** – Signieren und notarisieren Sie Ihre macOS-Anwendungen
- **Windows** – Signieren Sie Ihre ausführbaren Windows-Dateien und Pakete
- **Linux** – Signieren Sie DEB- und RPM-Pakete mit PGP-Schlüsseln

## Plattformübergreifende Signierungsmatrix

Diese Matrix zeigt, welche Formate Sie auf den einzelnen Quellplattformen signieren können:

| Zielformat | Unter Windows | Unter macOS | Unter Linux |
| --- | :---: | :---: | :---: |
| Windows EXE/DLL/MSIX/APPX, exportierte PFX-Datei | ✅ | ✅ | ✅ |
| Windows MSI oder Identität aus Zertifikatspeicher/Anbieter | ✅ | ❌ | ❌ |
| macOS .app / Mach-O / DMG, exportierte Identität | ✅ | ✅ | ✅ |
| macOS-Notarisierung, Team-API-Anmeldedaten | ✅ | ✅ | ✅ |
| macOS-Schlüsselbundidentität oder Notarisierungsprofil | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
Exportierte Apple-Identitäten und Team-API-Anmeldedaten ermöglichen portables Signieren und Notarisieren. macOS-Schlüsselbundidentitäten und Notarisierungsprofile erfordern einen Mac. Formatspezifische Einschränkungen werden weiter unten beschrieben.

@end

### Signierungs-Backends

Wails wählt automatisch das beste verfügbare Signierungs-Backend aus:

| Plattform | Natives Backend | Plattformübergreifendes Backend |
| --- | --- | --- |
| Windows | `signtool.exe` für Identitäten aus Zertifikatspeicher/Anbieter, MSI oder ältere RC2-PFX-Dateien | Integriert für exportierte RSA-PFX-Dateien und PE/MSIX |
| macOS | `codesign` für Schlüsselbundidentitäten | Integriert für exportierte PKCS#12-Identitäten oder Ad-hoc-Signierung |
| Linux | Nicht verfügbar | Integriert |

Der Anmeldedatentyp bestimmt auf jedem Host das Apple-Signierungs-Backend. Exportierte Identitäten verwenden Go; macOS-Schlüsselbundidentitäten verwenden `codesign`. Abruf, Überprüfung und Anheften von Tickets erfolgen mit Go.

## Schnellstart

Am schnellsten konfigurieren Sie die Signierung mit dem Einrichtungsassistenten:

```bash
wails3 setup
```

Dadurch wird eine gemeinsame Signierungskonfiguration in `~/.config/wails/defaults.yaml` geschrieben, die **beim Signieren auf jeder Plattform berücksichtigt wird** (siehe [Konfigurationsrangfolge](#konfigurationsrangfolge)). Der Signierungsschritt:

- Zeigt Einrichtungswerkzeuge wie OpenSSL und GnuPG zum Erzeugen oder Exportieren von Zertifikaten und Schlüsseln an. Exportierte Apple- und Windows-Identitäten verwenden den integrierten Signierer; Windows-Identitäten aus Zertifikatspeicher/Anbieter und MSI-Pakete erfordern Windows SignTool.
- Listet unter macOS die Developer-ID-Zertifikate aus Ihrem Schlüsselbund auf.
- Listet unter Linux Ihre GPG-Schlüssel auf und kann einen neuen Schlüssel **erzeugen und exportieren**.
- Kann unter Windows über OpenSSL ein **selbstsigniertes Zertifikat erzeugen** (für Tests).
- **Speichert Passwörter sicher im Schlüsselbund Ihres Systems** (nicht in Taskfiles).

@note{type="tip"}
Passwörter werden im nativen Anmeldedatenspeicher Ihres Systems gespeichert (macOS-Schlüsselbund, Windows-Anmeldeinformationsverwaltung oder Linux Secret Service). Dadurch ist Ihre Signierungskonfiguration geschützt und funktioniert in allen Ihren Wails-Projekten.

@end

## Konfigurationsrangfolge

Wenn Sie eine Signierungsaufgabe ausführen, wird jede Signierungsoption in dieser Reihenfolge aufgelöst (der erste Treffer gilt):

1. Ein explizit an `wails3 tool sign` übergebenes Flag (z. B. `--pgp-key`, `--certificate`, `--identity`).
2. Die passende **Taskfile-Variable des Projekts** (`PGP_KEY`, `SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`, `SIGN_IDENTITY`, …).
3. Die **globale Konfiguration** in `~/.config/wails/defaults.yaml` (geschrieben von `wails3 setup`).

Das bedeutet, dass die Taskfile-Variablen **optionale Überschreibungen** sind: Ist eine Variable nicht gesetzt, wird der global konfigurierte Schlüssel, das Zertifikat oder die Identität verwendet. Stellt keine der drei Quellen einen Wert bereit, meldet der Signierungsbefehl einen eindeutigen Fehler mit einem Hinweis zur Konfiguration.

## Projektbezogene Konfiguration

Um die Signierung nur für ein einzelnes Projekt statt global zu konfigurieren, führen Sie den projektbezogenen Assistenten im Projekt aus. Er schreibt `vars` in die `build/<platform>/Taskfile.yml`-Dateien dieses Projekts:

```bash
wails3 signing                                   # all detected platforms
wails3 signing --platform windows --platform linux
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

Das Passwort wird aus dem Systemschlüsselbund abgerufen (führen Sie zur Konfiguration `wails3 signing` aus).

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

Das Passwort wird aus dem Systemschlüsselbund abgerufen (führen Sie zur Konfiguration `wails3 signing` aus).

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
wails3 signing
```

## macOS-Codesignierung

### Voraussetzungen

- Apple-Developer-Konto ($99/Jahr)
- Developer ID Application-Zertifikat
- Für eine macOS-Schlüsselbundidentität: installierte Xcode Command Line Tools
- Für portables Signieren: eine exportierte RSA-PKCS#12-Identität; Apple-Signierungswerkzeuge sind nicht erforderlich

### Portables Signieren mit Wake

Wake liest `wails.yaml` direkt. Taskfiles werden nicht ausgeführt. Konfigurieren Sie eine exportierte Developer-ID-Identität und einen Team-API-Schlüssel für App Store Connect:

```yaml
darwin:
  signing:
    certificate: signing/developer-id.p12
    credential: APPLE_CERT_PASSWORD
    timestamp_server: https://timestamp.digicert.com
  notarization:
    api_key: signing/AuthKey.p8
    key_id: ABCDEFGHIJ
    issuer_id: 00000000-0000-0000-0000-000000000000
```

`credential` bezeichnet die Umgebungsvariable, die das Zertifikatspasswort enthält. Bewahren Sie beide privaten Schlüsseldateien außerhalb der Versionsverwaltung auf. Aktivieren Sie Signierung und Notarisierung im Release-Profil, das im [YAML-Build-Leitfaden](/experimental/yaml-builds/) beschrieben wird.

Der portable Signierer unterstützt Intel- und Apple-Silicon-Mach-O-Dateien mit 64 Bit, Universal Binaries, kanonische App-Bundles mit verschachteltem Code, Frameworks sowie UDIF-DMGs unter 4 GiB. Nicht unterstützte Layouts werden abgelehnt, bevor eine Ausgabe veröffentlicht wird. Ad-hoc-Signierung verwendet `identity: "-"` und kann nicht notarisiert werden. PKCS#12-Identitäten können AES/PBKDF2 oder die veraltete 3DES-Verschlüsselung verwenden; RC2-Exporte müssen mit einer unterstützten Chiffre erneut exportiert werden.

Portable notarisierte Signaturen erfordern ausdrücklich einen SHA-2-RFC-3161-Zeitstempeldienst. Veraltete SHA-1-Zeitstempelsignaturen werden abgelehnt. Die Übermittlung über die Team-API unterstützt einzelne Uploads bis zu 5 GiB; individuelle API-Schlüssel und mehrteilige Uploads werden nicht unterstützt. Ticketprüfung und Anheften verwenden Apples beobachtetes Dienstprotokoll mit einem fest vorgegebenen Apple Root G3, ohne `xcrun stapler`. Gatekeeper bleibt die maßgebliche Instanz für das Vertrauen in die Distribution und für Sperrungen.

Die folgenden Abschnitte beschreiben den vorhandenen macOS-Schlüsselbund-Workflow, den ältere Projekte und Wake-Konfigurationen verwenden, die eine Schlüsselbundidentität auswählen.

### Signieridentitäten

Prüfen Sie die verfügbaren Signieridentitäten:

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
Für die Distribution außerhalb des App Store benötigen Sie ein **Developer ID Application**-Zertifikat.

@end

### Konfiguration

Bearbeiten Sie `build/darwin/Taskfile.yml` und legen Sie die Signiervariablen fest:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| Variable | Erforderlich | Beschreibung |
| --- | --- | --- |
| `SIGN_IDENTITY` | Ja | Ihre Developer ID (z. B. „Developer ID Application: Ihr Unternehmen (TEAMID)“) |
| `KEYCHAIN_PROFILE` | Für die Notarisierung | Name des Schlüsselbundprofils mit gespeicherten Anmeldedaten |
| `ENTITLEMENTS` | Nein | Pfad zur Entitlements-Datei |

Führen Sie anschließend Folgendes aus:

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### Entitlements

Entitlements steuern, auf welche Funktionen Ihre App zugreifen kann. Wails-Apps benötigen für Entwicklung und Produktion üblicherweise unterschiedliche Entitlements:

- **Entwicklung**: Erfordert Entitlements für JIT, nicht signierten Speicher und Debugging
- **Produktion**: Minimale Entitlements (nur Netzwerkzugriff)

Erzeugen Sie beide Dateien mit dem interaktiven Einrichtungsassistenten:

```bash
wails3 entitlements
```

Dadurch werden folgende Dateien erstellt:

- `build/darwin/entitlements.dev.plist` – Für Entwicklungs-Builds
- `build/darwin/entitlements.plist` – Für Produktions-/signierte Builds

**Verfügbare Voreinstellungen:**\

| Voreinstellung | Beschreibung |
| --- | --- |
| Entwicklung | JIT, nicht signierter Speicher, Debugging, Netzwerk |
| Produktion | Nur Netzwerk (minimal, am sichersten) |
| Beide | Erstellt sowohl die Entwicklungs- als auch die Produktionsdatei (empfohlen) |
| App Store | Sandbox mit Netzwerk- und Dateizugriff aktiviert |
| Benutzerdefiniert | Einzelne Entitlements auswählen |

@note{type="note"}
Der Task `run` im darwin-Taskfile verwendet automatisch `entitlements.dev.plist`. Die `sign`-Tasks verwenden für Produktions-Builds `entitlements.plist`.

@end

Legen Sie anschließend in den Variablen Ihres Taskfiles `ENTITLEMENTS` so fest, dass es auf die passende Datei verweist.

### Notarisierung

Apple verlangt, dass alle distribuierten Apps notarisiert werden.

@steps
### **Speichern Sie Ihre Anmeldedaten im Schlüsselbund** (einmalige Einrichtung). Führen Sie entweder `wails3 signing` aus (der Befehl fragt diese Werte ab und ruft intern `notarytool` auf), oder rufen Sie `notarytool` direkt auf:
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **Legen Sie KEYCHAIN_PROFILE in Ihrem Taskfile** so fest, dass es dem obigen Profilnamen entspricht.
### **Signieren und notarisieren Sie Ihre App**:
```bash
wails3 task darwin:sign:notarize
```

### **Überprüfen Sie die Notarisierung**:
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
Die Notarisierung dauert üblicherweise 1-2 Minuten. Das Ticket wird automatisch an Ihre App angeheftet.

@end

## Codesignierung unter Windows

### Voraussetzungen

- Codesignierungszertifikat (von DigiCert, Sectigo usw.)
- Für die portable Signierung: eine exportierte RSA-PFX-Identität, die mit AES/PBKDF2 oder dem veralteten 3DES verschlüsselt ist; externe Signierwerkzeuge sind nicht erforderlich
- Für Identitäten aus dem Zertifikatspeicher bzw. von einem Anbieter, MSI-Pakete oder ältere RC2-PFX-Dateien: Windows mit den Signierwerkzeugen des Windows SDK unter `PATH`

### Selbstsigniertes Zertifikat erzeugen (Tests)

Wenn Sie lediglich die Signierpipeline testen müssen, führen Sie `wails3 setup` aus, öffnen Sie im Signierschritt die Registerkarte **Windows** und wählen Sie **Selbstsigniertes Zertifikat erzeugen**. Dabei wird mit OpenSSL ein Codesignierungs-`.pfx` erstellt und dessen Pfad in der globalen Konfiguration gespeichert.

@note{type="caution"}
Selbstsignierte Zertifikate sind **ausschließlich für Tests und die interne Distribution vorgesehen** – bei Endbenutzern lösen sie SmartScreen-Warnungen aus. Öffentliche Releases benötigen ein Zertifikat von einer vertrauenswürdigen Zertifizierungsstelle.

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
| `SIGN_CERTIFICATE` | Überschreibung | Pfad zur .pfx/.p12-Zertifikatsdatei (greift auf die globale Konfiguration zurück, wenn nicht festgelegt) |
| `SIGN_THUMBPRINT` | Überschreibung | Fingerabdruck des Zertifikats im Windows-Zertifikatspeicher (Alternative zu `SIGN_CERTIFICATE`) |
| `TIMESTAMP_SERVER` | Nein | URL des Zeitstempelservers (Standard: http://timestamp.digicert.com) |

@note{type="note"}
Diese Variablen sind optionale **Überschreibungen**. Sind sie nicht festgelegt, wird das über `wails3 setup` global konfigurierte Zertifikat verwendet. Siehe [Konfigurationsrangfolge](#konfigurationsrangfolge).

@end

@note{type="note"}
Das Zertifikatspasswort wird im Schlüsselbund Ihres Systems gespeichert, nicht im Taskfile. Führen Sie zur Konfiguration `wails3 signing` aus oder legen Sie in CI die Umgebungsvariable `WAILS_WINDOWS_CERT_PASSWORD` fest.

@end

Führen Sie anschließend Folgendes aus:

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### Plattformübergreifende Signierung

Exportierte RSA-PFX-Identitäten signieren EXE-/DLL-Dateien, ausführbare NSIS-Installationsprogramme und einzelne MSIX-/APPX-Pakete mit Go unter Windows, macOS und Linux. Vorhandene Signaturen können ersetzt werden. Beim Signieren wird die Ausgabe zunächst bereitgestellt und erst veröffentlicht, nachdem die Signatur- und Zeitstempelprüfung erfolgreich war. Bei einem Fehler bleibt das Original erhalten. Verwenden Sie `--output`, um eine separate signierte Datei zu schreiben.

Wake liest diese Einstellungen direkt aus `wails.yaml`:

```yaml
windows:
  signing:
    certificate: signing/windows.pfx
    credential: WINDOWS_SIGNING_PASSWORD
    timestamp_server: http://timestamp.digicert.com
```

`credential` gibt den Namen der Umgebungsvariable für das Passwort an, niemals das Passwort selbst. Aktivieren Sie `sign: true` für das Release-Ziel; siehe den [YAML-Build-Leitfaden](/experimental/yaml-builds/). Wake ruft keine Taskfiles auf. Ältere Projekte können weiterhin die obige Taskfile-Konfiguration verwenden.

Das portable Profil verwendet RSA/SHA-256 und validiert RFC-3161-Zeitstempel einschließlich der vertrauenswürdigen Zertifikatskette des Zeitstempelsignierers. Standardmäßig wird der Dienst von DigiCert verwendet. Der MSIX-Herausgeber muss mit dem Antragsteller des Signaturzertifikats übereinstimmen. Paket-Bundles und Archive mit mindestens 4 GiB werden nicht unterstützt. PFX-Dateien sind auf 16 MiB begrenzt. Andere Zertifikatsalgorithmen und RC2-Verschlüsselung werden abgelehnt. Exportieren Sie eine RC2-PFX-Datei für die portable Signierung erneut mit AES oder 3DES. Unter Windows verwendet eine authentifizierte RC2-PFX-Datei stattdessen SignTool. Ein falsches Passwort oder eine beschädigte PFX-Datei löst diesen Rückgriff niemals aus.

Schlüssel aus Zertifikatspeichern, Hardware-Token und anbietergebundene Schlüssel verwenden weiterhin den nativen Windows-SignTool-Workflow mit `thumbprint`. Auch die MSI-Signierung erfordert SignTool unter Windows. Diese Workflows werden durch die Konfiguration eines exportierten Zertifikats nicht portabel.

### Unterstützte Windows-Formate

| Format | Erweiterung | Hinweise |
| --- | --- | --- |
| Ausführbare Dateien und Bibliotheken | .exe, .dll | Portable PE-Signierung einschließlich NSIS-Installationsprogrammen |
| Installationsprogramme | .msi | Windows SignTool erforderlich |
| App-Pakete | .msix, .appx | Portable Signierung einzelner Pakete; keine Bundles |

## Signieren von Linux-Paketen

Linux-Pakete (DEB und RPM) werden mit PGP-/GPG-Schlüsseln signiert. Anders als die Codesignierung unter Windows und macOS weist die Signierung von Linux-Paketen nach, dass das Paket aus einer vertrauenswürdigen Quelle stammt, nicht dass das Betriebssystem dem Code vertraut.

Eine exportierte, ASCII-armierte oder binäre Datei mit einem privaten Schlüssel verwendet unter Windows, macOS und Linux den nativen Go-Signierer. Beim Signieren eines vorhandenen DEB- oder RPM-Pakets bleibt dessen Paketinhalt ohne Entpacken oder erneute Komprimierung erhalten; eine Installation von GnuPG, `dpkg-sig` oder `rpmsign` ist nicht erforderlich. Arch-Pakete verwenden denselben Schlüssellader, um eine abgetrennte Signatur zu schreiben. GnuPG bleibt zum Erzeugen und Exportieren von Schlüsseln sowie zur unabhängigen Signaturprüfung nützlich. Ältere DEB-/RPM-Konfigurationen, die statt einer Schlüsseldatei eine Schlüsselbundidentität angeben, verwenden weiterhin ihre externen Signierwerkzeuge.

### Voraussetzungen

- PGP-Schlüsselpaar (kann mit Wails erzeugt werden)

### PGP-Schlüssel erzeugen

Am einfachsten geht es mit dem Einrichtungsassistenten: Führen Sie `wails3 setup` aus, öffnen Sie im Signierungsschritt die Registerkarte **Linux** und wählen Sie **Neuen GPG-Schlüssel erstellen**. Der Assistent:

- erzeugt einen RSA-4096-Schlüssel in Ihrem GPG-Schlüsselbund (lassen Sie die Passphrase für einen unbeaufsichtigten, CI-tauglichen Schlüssel leer),
- **exportiert ihn nach `~/.wails/signing/<keyid>.asc`** (die Datei, mit der ein Build signiert) und
- speichert sowohl die Schlüssel-ID als auch den Exportpfad in `~/.config/wails/defaults.yaml`, sodass der Schlüssel beim Signieren automatisch verwendet wird.

@note{type="note"}
Wenn ein Schlüssel nur über die **ID** konfiguriert wurde (beispielsweise vor Einführung des automatischen Schlüsselexports durch den Assistenten), exportiert ihn die Linux-Registerkarte beim nächsten Öffnen in eine Datei und trägt den Pfad ein. Manuelle Schritte sind nicht erforderlich. Builds signieren mit einer Schlüssel*datei*, daher ist der Pfad entscheidend.

@end

Sie können dies auch manuell mit `gpg` erledigen:

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

Empfohlene Einstellungen: RSA mit 4096 Bit, Ablauf nach 1 Jahr und Schutz durch ein starkes Passwort.

@note{type="caution"}
Bewahren Sie Ihren privaten Schlüssel sicher auf! Speichern Sie ihn verschlüsselt und sichern Sie ihn zuverlässig.

@end

### Konfiguration

Bearbeiten Sie `build/linux/Taskfile.yml` und legen Sie die Signierungsvariablen fest:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| Variable | Erforderlich | Beschreibung |
| --- | --- | --- |
| `PGP_KEY` | Überschreibung | Pfad zur exportierten Datei mit dem privaten PGP-Schlüssel (wenn nicht festgelegt, wird der über `wails3 setup` global konfigurierte Schlüssel verwendet) |
| `SIGN_ROLE` | Nein | DEB-Signierungsrolle (Standard: builder) |

@note{type="note"}
Das Passwort des PGP-Schlüssels wird im Schlüsselbund Ihres Systems gespeichert, nicht im Taskfile. Führen Sie `wails3 signing` aus, um es zu konfigurieren, oder legen Sie in CI die Umgebungsvariable `WAILS_PGP_PASSWORD` fest.

@end

Führen Sie anschließend Folgendes aus:

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### DEB-Signierungsrollen

Für DEB-Pakete können Sie die Signierungsrolle über `SIGN_ROLE` angeben:

- `origin`: Signatur der Paketquelle
- `maint`: Signatur des Paketbetreuers
- `archive`: Signatur des Archivbetreuers
- `builder`: Signatur des Paketerstellers (Standard)

### Plattformübergreifende Signierung

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

In CI-Umgebungen werden Passwörter anstelle des Systemschlüsselbunds über Umgebungsvariablen bereitgestellt:

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

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
Ein Linux-Runner für die plattformübergreifende Signierung vereinfacht CI/CD, da keine separaten Windows-Runner erforderlich sind. Apple-Schlüsselbundidentitäten und Notarisierungsprofile erfordern einen macOS-Runner. Exportierte Identitäten und Team-API-Anmeldedaten unterstützen den oben beschriebenen portablen Wake-Workflow.

@end

## CLI-Referenz

### wails3 signing

Interaktiver Assistent zum Konfigurieren der Signierung für Ihr Projekt.

```bash
wails3 signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

Der Assistent führt Sie durch:

- **macOS**: Auswählen eines Developer-ID-Zertifikats und Konfigurieren der Anmeldedaten für die Beglaubigung (ruft `xcrun notarytool store-credentials` auf).
- **Windows**: Auswählen zwischen Zertifikatsdatei und Fingerabdruck sowie Festlegen des Passworts und des Zeitstempelservers.
- **Linux**: Verwenden eines vorhandenen oder Erzeugen eines neuen PGP-Schlüssels (ruft `gpg` auf) sowie Konfigurieren der Signaturrolle.

### wails3 entitlements

Interaktiver Assistent zum Konfigurieren von macOS-Berechtigungen.

```bash
wails3 entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**Voreinstellungen:**

- **Entwicklung**: Erstellt `entitlements.dev.plist` mit JIT, Debugging und Netzwerkzugriff
- **Produktion**: Erstellt `entitlements.plist` mit minimalen Berechtigungen
- **Beide**: Erstellt beide Dateien (empfohlen)
- **App Store**: Erstellt Sandbox-Berechtigungen für den Mac App Store
- **Benutzerdefiniert**: Einzelne Berechtigungen und Zieldatei auswählen

### wails3 sign

Signiert Binärdateien und Pakete für die aktuelle oder angegebene Plattform. Dieser Wrapper ruft die passende plattformspezifische Signaturaufgabe auf.

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

Dadurch wird die entsprechende Aufgabe `<platform>:sign` ausgeführt, die die Signaturkonfiguration aus Ihrem Taskfile verwendet.

### wails3 tool sign

Low-Level-Befehl zum direkten Signieren einer bestimmten Datei. Wird intern von den Taskfiles verwendet.

```bash
wails3 tool sign [flags]
```

**Allgemeine Optionen:**\

| Option | Beschreibung |
| --- | --- |
| `--input` | Pfad zur zu signierenden Datei |
| `--output` | Ausgabepfad (optional; standardmäßig wird die Datei direkt geändert) |
| `--verbose` | Ausführliche Ausgabe aktivieren |

**Optionen für Windows/macOS:**\

| Option | Beschreibung |
| --- | --- |
| `--certificate` | Pfad zum PKCS#12-Zertifikat (.pfx/.p12) |
| `--password` | Zertifikatspasswort |
| `--timestamp` | URL des Zeitstempelservers |

**macOS-spezifische Optionen:**\

| Option | Beschreibung |
| --- | --- |
| `--identity` | Signaturidentität ('-' für Ad-hoc-Signierung verwenden) |
| `--entitlements` | Pfad zur Berechtigungs-plist |
| `--hardened-runtime` | Hardened Runtime aktivieren (Standard: true) |
| `--notarize` | Zur Beglaubigung einreichen |
| `--keychain-profile` | Schlüsselbundprofil für die Beglaubigung |
| `--notary-key` | `.p8`-Datei mit dem privaten Team-API-Schlüssel für App Store Connect; verwendet den Go-Notary-API-Client |
| `--notary-key-id` | Team-API-Schlüssel-ID; zusammen mit `--notary-key` erforderlich |
| `--notary-issuer-id` | Team-Aussteller-ID; zusammen mit `--notary-key` erforderlich |

**Windows-spezifische Optionen:**\

| Option | Beschreibung |
| --- | --- |
| `--thumbprint` | Zertifikatfingerabdruck im Windows-Zertifikatspeicher |

**Linux-spezifische Optionen:**\

| Option | Beschreibung |
| --- | --- |
| `--pgp-key` | Pfad zum privaten PGP-Schlüssel |
| `--pgp-password` | Passwort des PGP-Schlüssels |
| `--role` | DEB-Signaturrolle (origin/maint/archive/builder) |

### Signaturstatus prüfen (native Werkzeuge)

`wails3 signing` konfiguriert die Anmeldedaten für die Signierung. Verwenden Sie zum Prüfen des Signaturstatus direkt die nativen Werkzeuge:

| Aufgabe | Befehl |
| --- | --- |
| macOS-Identitäten für die Codesignierung auflisten | `security find-identity -v -p codesigning` |
| Anmeldedaten für die Beglaubigung speichern | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| PGP-Schlüsseldatei prüfen | `gpg --show-keys <key.asc>` |
| PGP-Schlüsselpaar erzeugen | `gpg --full-generate-key` |
| Öffentlichen Schlüssel exportieren | `gpg --armor --export <email>` |

## Fehlerbehebung

### macOS-Probleme

**„Kein Developer-ID-Zertifikat gefunden“**

- Stellen Sie sicher, dass Ihr Zertifikat im Schlüsselbund installiert ist
- Prüfen Sie mit `security find-identity -v -p codesigning`, ob es abgelaufen ist
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

- Prüfen Sie, ob der Zertifikatspfad korrekt ist
- Prüfen Sie das Zertifikatspasswort
- Stellen Sie sicher, dass das Zertifikat gültig ist (nicht abgelaufen oder widerrufen).

**„Fehler des Zeitstempelservers“**

- Versuchen Sie es mit einem anderen Zeitstempelserver:
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Linux-Probleme

**„Ungültiger PGP-Schlüssel“**

- Stellen Sie sicher, dass die Schlüsseldatei im ASCII-armored-Format vorliegt.
- Prüfen Sie mit `gpg --show-keys <key.asc>`, ob der Schlüssel abgelaufen ist.
- Überprüfen Sie, ob das Passwort korrekt ist.

**„Signaturprüfung fehlgeschlagen“**

- Stellen Sie sicher, dass der öffentliche Schlüssel ordnungsgemäß importiert wurde.
- Prüfen Sie, ob das Paket nach dem Signieren verändert wurde.

## Weitere Ressourcen

### Offizielle Dokumentation

- [Apple-Leitfaden zur Codesignierung](https://developer.apple.com/support/code-signing/)
- [Apple-Dokumentation zur Notarisierung](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Microsoft-Codesignierung](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Signieren von Debian-Paketen](https://wiki.debian.org/SecureApt)
- [Signieren von RPM-Paketen](https://rpm-software-management.github.io/rpm/manual/signatures.html)
