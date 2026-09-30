---
title: "Triage durch Maintainer"
description: "Fehlerberichte, Dokumentationsmeldungen und WEPs einheitlich triagieren"
slug: "contributing/maintainer-triage"
sourcePath: "contributing/maintainer-triage.md"
---

## Zweck

Issues erfassen reproduzierbare Fehler und Dokumentationsprobleme. Ein Pull Request für einen WEP (Wails Enhancement Proposal) erfasst eine vorgeschlagene Funktion oder Änderung des öffentlichen Verhaltens.

## Issue-Triage

- Stelle sicher, dass ein Fehlerbericht eine veröffentlichte Version, die Plattform, reproduzierbare Schritte, das erwartete Verhalten, das tatsächliche Verhalten und die Ausgabe von `wails3 doctor` enthält.
- Kennzeichne bestätigte Meldungen mit `Bug` und füge die relevanten Versions- und Plattform-Labels hinzu. Bitte bei Bedarf um ein minimales Reproduktionsbeispiel.
- Lasse Dokumentationsmeldungen offen, wenn sie einen konkreten Dokumentationsfehler benennen. Ermutige die meldende Person zu einem PR, wenn sie die Änderung selbst vornehmen kann.
- Verweise Funktionsanfragen auf den WEP-Leitfaden und schließe sie anschließend. Der automatisierte Weiterleitungs-Workflow verarbeitet neu als Enhancement gekennzeichnete Issues; verwende für ältere Issues denselben Wortlaut.
- Verschiebe Fragen und Supportanfragen in GitHub Discussions oder nach Discord.

## WEP-Triage

1. Prüfe, ob der PR ein Entwurf mit dem Titel `[WEP] <title>` ist und ausschließlich den WEP und unterstützendes Material enthält.
2. Prüfe, ob er die WEP-Vorlage verwendet, eine für die Implementierung verantwortliche Person benennt und Kompatibilität, Plattformen, Tests, Wartung sowie Sicherheit und Datenschutz behandelt.
3. Führe die technische Diskussion im WEP-PR. Diskussionen liefern nützlichen Kontext, dokumentieren aber nicht die Entscheidung.
4. Halte die Entscheidung der Maintainer in einem PR-Kommentar fest: angenommen, abgelehnt oder zurückgezogen, jeweils mit einer kurzen Begründung.
5. Weise einem angenommenen WEP seine Nummer zu, aktualisiere den WEP-Index, merge den WEP-PR und verlange, dass Implementierungs-PRs auf ihn zurückverweisen.

## Bestehende Enhancement-Issues

Lösche historische Enhancement-Issues nicht stillschweigend. Hinterlasse für jede weiterhin relevante Anfrage den Weiterleitungskommentar und schließe sie; interessierte Mitwirkende können einen WEP eröffnen. Schließe Duplikate mit einem Link zum bestehenden WEP oder zur bestehenden Entscheidung.
