---
title: "Projektstatus"
description: "Beta-Kompatibilität, Sicherheitsunterstützung und Upgrade-Anleitung für Wails v3"
slug: "status"
sourcePath: "status.md"
---

## Aktueller Status: Beta

Den aktuellen Status finden Sie im [Änderungsprotokoll](/changelog/).

Unser Ziel ist eine stabile Version v3.0. Wails v2 bleibt die aktuelle stabile Version und erhält weiterhin Fehlerbehebungen. Testen Sie Beta-Versionen mit Ihrer Anwendung vor der Bereitstellung.

## Kompatibilitätsversprechen für die Beta

Die Beta-Kompatibilitätszusage für v3 gilt für Desktop-Anwendungen:

| Plattform | Unterstützte Zielsysteme | Anforderungen und Hinweise |
| --- | --- | --- |
| Windows | amd64 und arm64 | WebView2-Laufzeitumgebung |
| macOS | Intel und Apple Silicon | Die in der Installationsanleitung dokumentierten macOS- und WebKit-Versionen |
| Linux | amd64 und arm64 | Standardmäßig GTK4 + WebKitGTK 6.0; GTK3 + WebKit2GTK 4.1 bleibt bis einschließlich v3.0.x als `-tags gtk3` Legacy-Option verfügbar und wird in v3.1 entfernt |

Für die Entwicklung auf allen Zielsystemen ist Go 1.25 oder neuer erforderlich. Die Unterstützung für Android und iOS ist experimentell und steht der Desktop-Beta nicht im Weg. Die Beta-APIs sollen stabil bleiben, doch Fehler in Vorabversionen und ausdrücklich angekündigte Änderungen können noch vor v3.0.0 korrigiert werden.

## So können Sie mitwirken

- Testen Sie die neueste Beta-Version und melden Sie reproduzierbare Fehler
- Wirken Sie an der Dokumentation und an Beispielen mit
- Beteiligen Sie sich an Diskussionen und geben Sie Feedback zu WEP-Entwürfen
- Reichen Sie Pull Requests für Fehlerbehebungen, Dokumentation oder angenommene WEPs ein

Beiträge aus der Community sind willkommen. Wenn Sie uns bei diesen Zielen unterstützen möchten, beteiligen Sie sich an den Diskussionen der Community. Vorschläge für neue Funktionen gehören in einen WEP-PR und nicht in ein Issue für eine Funktionsanfrage.

## Feedback und Aktualisierungen

Melden Sie reproduzierbare Probleme als Issues; schlagen Sie neue Funktionen über einen PR für ein [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md) vor.

## Die Beta verwenden

Legen Sie genaue Versionen für CLI, Go-Modul und Frontend-Runtime fest, statt `latest` zu verwenden. Folgen Sie bei einem bestehenden Alpha-Projekt der [Upgrade-Anleitung von Alpha auf Beta](/migration/alpha-to-beta/).

Die [Sicherheitsrichtlinie](https://github.com/wailsapp/wails/blob/master/SECURITY.md) führt v3-Beta-Versionen als unterstützt und Alpha-Versionen als nicht unterstützt auf. Melden Sie Sicherheitslücken über die [private Schwachstellenmeldung](https://github.com/wailsapp/wails/security/advisories/new), nicht in öffentlichen Issues.

## Erfasste Arbeiten

- [Offene Fehler mit dem Label v3](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3ABug+label%3Av3)
- [Offene v3-Issues mit dem Label P0 oder P1](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3Av3+label%3AP0%2CP1)
- [Release-Meilensteine](https://github.com/wailsapp/wails/milestones)

Diese aktuellen Abfragen hängen von den Issue-Labels ab. Sie sind weder eine vollständige Liste noch eine Zusage zu einem Veröffentlichungsdatum oder Umfang. Lesen Sie die Issues, um ihre Auswirkungen auf Ihr Projekt einzuschätzen.
