---
title: "Windows-Paketierung"
description: "Paketieren Sie Ihre Wails-Anwendung für die Verteilung unter Windows"
slug: "guides/build/windows"
sourcePath: "guides/build/windows.md"
---

## NSIS-Installationsprogramm

Das standardmäßige Paketformat erstellt ein NSIS-Installationsprogramm:

```bash
wails3 package GOOS=windows
```

Dadurch wird `wails3 task windows:package` ausgeführt. Der Vorgang:

1. Erstellt die Anwendung
2. Generiert den WebView2-Bootstrapper
3. Erstellt ein NSIS-Installationsprogramm

Ausgabe: `build/windows/nsis/<AppName>-installer.exe`

### MSIX-Paket

Für die Verteilung über den Microsoft Store oder die moderne Windows-Bereitstellung:

```bash
wails3 package GOOS=windows FORMAT=msix
```

Ausgabe: `bin/<AppName>-<arch>.msix`

@note{type="note"}
MSIX erfordert entweder `makeappx.exe` (Windows SDK) oder die eigenständigen MSIX-Werkzeuge. Das Windows-Taskfile stellt die Installationsaufgabe als `wails3 task install:msix:tools` bereit.

@end

## Installationsprogramm anpassen

Die NSIS-Konfiguration befindet sich in `build/windows/nsis/project.nsi`. Bearbeiten Sie diese Datei, um Folgendes anzupassen:

- Benutzeroberfläche und Branding des Installationsprogramms
- Installationsverzeichnis
- Verknüpfungen im Startmenü und auf dem Desktop
- Dateizuordnungen
- Lizenzvereinbarung

Die Anwendungsmetadaten stammen aus `build/windows/info.json`:

```json
{
  "fixed": {
    "file_version": "1.0.0"
  },
  "info": {
    "0000": {
      "ProductVersion": "1.0.0",
      "CompanyName": "My Company",
      "FileDescription": "My Application",
      "ProductName": "MyApp"
    }
  }
}
```

## Codesignierung

Signieren Sie Ihre ausführbare Datei und das Installationsprogramm, um SmartScreen-Warnungen zu vermeiden:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=windows

# Or using tasks directly
wails3 task windows:sign
wails3 task windows:sign:installer
```

Konfigurieren Sie die Signierung in `build/windows/Taskfile.yml`:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint for certificates in Windows store
  SIGN_THUMBPRINT: "certificate-thumbprint"
  TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

Bewahren Sie das Kennwort für Ihr Zertifikat sicher auf:

```bash
wails3 setup signing
```

Weitere Informationen finden Sie unter [Anwendungen signieren](/guides/build/signing/).

## Für ARM erstellen

```bash
wails3 build GOOS=windows GOARCH=arm64
wails3 package GOOS=windows GOARCH=arm64
```

## Fehlerbehebung

### makensis nicht gefunden

Installieren Sie NSIS:

```bash
# Windows
winget install NSIS.NSIS

# Or download from https://nsis.sourceforge.io/
```

### SmartScreen-Warnung

Ihre ausführbare Datei ist nicht signiert. Lesen Sie den obigen Abschnitt [Codesignierung](#codesignierung).

### WebView2 fehlt

Das Installationsprogramm enthält einen WebView2-Bootstrapper, der die Laufzeitumgebung bei Bedarf herunterlädt. Wenn Sie eine Offlineinstallation benötigen, laden Sie den Evergreen Standalone Installer von Microsoft herunter.
