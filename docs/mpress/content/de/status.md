---
title: "Roadmap"
description: "Projektstatus von Wails v3, geplante Funktionen und Möglichkeiten zur Mitwirkung"
slug: "status"
sourcePath: "status.md"
---

## Wails v3.0.0 — stabil

Den aktuellen Stand finden Sie im [Changelog](/changelog/).

Diese Dokumentation behandelt Wails v3.0.0. Nach v3.0.0 werden Fehlerbehebungen und Verbesserungen in wöchentlichen Releases veröffentlicht, dazwischen gibt es datierte Nightly-Vorabversionen und bei dringenden Problemen Hotfixes. Weitere Informationen finden Sie in der [Release-Richtlinie](/contributing/release-policy/).

## Desktop-Kompatibilität

Der Kompatibilitätsvertrag für v3.0.0 gilt für Desktop-Anwendungen:

| Plattform | Unterstützte Zielsysteme | Anforderungen und Hinweise |
| --- | --- | --- |
| Windows | amd64 und arm64 | WebView2-Laufzeitumgebung |
| macOS | Intel und Apple Silicon | Die in der Installationsanleitung dokumentierten Versionen von macOS und WebKit |
| Linux | amd64 und arm64 | Standardmäßig GTK4 + WebKitGTK 6.0; GTK3 + WebKit2GTK 4.1 bleibt bis einschließlich v3.0.x eine `-tags gtk3` Legacy-Option und wird in v3.1 entfernt |

Für die Entwicklung auf allen Zielsystemen ist Go 1.26 oder neuer erforderlich. Die Unterstützung für Android und iOS bleibt experimentell und ist nicht Bestandteil des stabilen Kompatibilitätsvertrags für Desktop-Anwendungen.

## So können Sie mitwirken

- Wails v3.0.0 testen und reproduzierbare Fehler melden
- Zur Dokumentation und zu den Beispielen beitragen
- An Diskussionen teilnehmen und Feedback zu WEP-Entwürfen geben
- Pull Requests für Fehlerkorrekturen, Dokumentation oder angenommene WEPs einreichen

Beiträge aus der Community sind willkommen. Wenn Sie an diesen Zielen mitwirken möchten, beteiligen Sie sich an den Community-Diskussionen. Vorschläge für neue Funktionen gehören in einen WEP-PR, nicht in ein Feature-Request-Issue.

## Feedback und Aktualisierungen

Diese Roadmap kann sich aufgrund von Community-Feedback und Projektprioritäten ändern. Wir aktualisieren sie regelmäßig, um Fortschritte und Kursänderungen abzubilden. Melden Sie reproduzierbare Probleme als Issues; schlagen Sie neue Funktionen über einen [WEP-PR (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md) vor.
