---
title: "Roadmap"
description: "Projektstatus von Wails v3, geplante Funktionen und Möglichkeiten zur Mitwirkung"
slug: "status"
sourcePath: "status.md"
---

## Aktueller Status: Beta

Den aktuellen Status finden Sie im [Änderungsprotokoll](/changelog/).

Unser Ziel ist eine stabile Version v3.0. Diese Roadmap beschreibt die wichtigsten Funktionen und Verbesserungen, die wir vor der endgültigen Veröffentlichung implementieren müssen. Beachten Sie, dass dies ein lebendes Dokument ist und aktualisiert werden kann, wenn sich Prioritäten verschieben oder neue Erkenntnisse ergeben.

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

Diese Roadmap kann sich aufgrund des Feedbacks der Community und der Projektprioritäten ändern. Wir aktualisieren sie regelmäßig, um Fortschritte und Richtungsänderungen widerzuspiegeln. Melden Sie reproduzierbare Probleme als Issues; schlagen Sie neue Funktionen über einen PR für ein [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md) vor.
