---
title: "MSIX-Paketierung"
description: "Ihre Wails-v3-Anwendung als MSIX-Paket paketieren"
slug: "guides/build/msix"
sourcePath: "guides/build/msix.md"
---

MSIX ist das moderne Paketformat für Windows-Anwendungen. Wails kann im Rahmen Ihres Windows-Builds ein MSIX-Paket erstellen.

Anweisungen zur MSIX-Paketierung finden Sie im Leitfaden [Windows-Paketierung](/guides/build/windows/#msix-package).

## Pakete direkt über die CLI erstellen

Führe das MSIX-Werkzeug unter Windows aus, auch in CI. Das Standard-Backend verwendet `MakeAppx.exe`; zum Signieren wird zusätzlich `signtool.exe` benötigt. Beide sind im Windows SDK enthalten. Der Installationshelfer öffnet den Microsoft Store und bei Bedarf die SDK-Downloadseite; schließe die Installation vor dem Paketieren ab.

Lege die Identität in `build/config.yml` fest und erstelle anschließend die ausführbare Datei und das Paket:

```yaml
info:
  companyName: "Example Corp"
  productName: "MyApp"
  productIdentifier: "com.example.myapp"
  description: "MyApp"
  version: "1.0.0"
```

```powershell
wails3 tool msix-install-tools
wails3 build GOOS=windows
wails3 tool msix --executable bin/myapp.exe --name myapp.exe
```

Der direkte Befehl schreibt `MyApp.msix` in das aktuelle Verzeichnis. Die [Windows-Paketierungsaufgabe](/guides/build/windows/#msix-package) gibt hingegen einen eigenen Ausgabepfad vor.

## CLI-Optionen

| Option | Bedeutung |
| --- | --- |
| `--config` | Konfigurationsdatei; Standard: `build/config.yml`. |
| `--executable`, `--name` | Vorhandene ausführbare Datei und ihr Dateiname im Paket; beide sind erforderlich. |
| `--out` | Ausgabedatei; Standard: `<ProductName>.msix`. |
| `--arch` | Paketarchitektur: `x64` (Standard), `x86`, `arm`, `arm64`, `x86a64` oder `neutral`. Die Go-Aliase `amd64` und `386` werden akzeptiert. Verwende die Architektur der ausführbaren Datei. |
| `--publisher` | Herausgeberidentität; Standard: `CN=<companyName>`. |
| `--cert`, `--cert-password` | Pfad zum PFX-Zertifikat und Passwort zum Signieren. |
| `--use-makeappx` | Den Standard-Paketierer des Windows SDK verwenden. |
| `--use-msix-tool` | `MsixPackagingTool.exe` ausdrücklich auswählen; es muss in `PATH` liegen. |

## Signierung und CI

Signiere für die Verteilung außerhalb des Stores mit einem Zertifikat, dem der Zielrechner vertraut. Das Subject des Zertifikats muss exakt mit `--publisher` übereinstimmen. Das MakeAppx-Backend ruft SignTool mit SHA256 auf, wenn `--cert` angegeben wird. Siehe [Microsofts Signierungsanleitung](https://learn.microsoft.com/en-us/windows/msix/package/sign-msix-package-guide).

Dieser Windows-Workflow-Schritt setzt voraus, dass Wails und das SDK installiert sind und ein vorheriger Schritt die PFX-Datei sicher unter `CERT_PATH` bereitgestellt hat. Ein Secret mit einem Pfad lädt allein kein Zertifikat hoch:

```yaml
- name: MSIX
  if: runner.os == 'Windows'
  shell: pwsh
  run: |
    wails3 build GOOS=windows
    wails3 tool msix --executable bin/myapp.exe --name myapp.exe --publisher "$env:MSIX_PUBLISHER" --cert "$env:CERT_PATH" --cert-password "$env:CERT_PASSWORD"
  env:
    MSIX_PUBLISHER: ${{ vars.MSIX_PUBLISHER }}
    CERT_PATH: ${{ secrets.WINDOWS_CERT_PATH }}
    CERT_PASSWORD: ${{ secrets.WINDOWS_CERT_PASSWORD }}
```

## Dateizuordnungen und Ressourcen

Füge Erweiterungen ohne führenden Punkt in `build/config.yml` ein; das generierte Manifest ergänzt den Punkt:

```yaml
fileAssociations:
  - ext: myext
    name: MyApp Document
    description: MyApp Document
    iconName: fileicon
```

Behandle das Öffnen von Dateien zur Laufzeit wie unter [Dateizuordnungen](/guides/file-associations/) beschrieben. Das MakeAppx-Backend kopiert derzeit nur die ausführbare Datei und erzeugt transparente Platzhalterbilder. Es importiert keine `Assets/`-Dateien des Projekts und konvertiert keine `iconName`-Symbole. Verwende für eigene Grafiken oder zusätzliche DLLs einen angepassten Paketierungsablauf.

| Generierte Ressource | Größe (Pixel) |
| --- | --- |
| `Square150x150Logo.png` | 150×150 |
| `Square44x44Logo.png` | 44×44 |
| `Wide310x150Logo.png` | 310×150 |
| `StoreLogo.png` | 50×50 |
| `SplashScreen.png` | 620×300 |
| `FileIcon.png` | 44×44 |

`FileIcon.png` wird nur erzeugt, wenn Dateizuordnungen konfiguriert sind. Diese Dateien befinden sich im Verzeichnis `Assets/` des Pakets.

## Store-Einreichung und Fehlerbehebung

Reserviere deine Anwendung im [Partner Center für Entwickler](https://partner.microsoft.com/dashboard) und verwende dessen Paketidentität und Herausgeber bei der Vorbereitung der Einreichung. Der Store signiert MSIX-Pakete bei der Einreichung; dafür brauchst du kein Signierungszertifikat zu kaufen. Siehe [Microsofts Paketanforderungen](https://learn.microsoft.com/en-us/windows/apps/publish/publish-your-app/msix/app-package-requirements).

Wenn `MakeAppx.exe` oder `signtool.exe` nicht gefunden wird, installiere oder repariere das Windows SDK. Wails durchsucht `PATH` und die üblichen SDK-Verzeichnisse. Prüfe bei Signierungsfehlern das Subject des Zertifikats, seine Gültigkeit und das Vertrauen auf dem Zielrechner; siehe [MSIX-Fehlerbehebung](https://learn.microsoft.com/en-us/windows/msix/msix-troubleshooting-guide).
