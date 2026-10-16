---
title: "Dokumentation korrigieren"
description: "Reiche mit M-Press einen Korrektur-PR für die Wails-v3-Dokumentation ein."
sourcePath: "contributing/documentation.md"
---

Korrektur-PRs sind willkommen. Korrigiere Tippfehler, defekte Links, veraltete Beispiele, unklare Erklärungen oder Übersetzungen. Für eine reine Dokumentationskorrektur benötigst du weder ein Issue noch einen fehlgeschlagenen Code-Test.

## Lokale Vorschau

Forke [wailsapp/wails](https://github.com/wailsapp/wails/fork), klone deinen Fork und erstelle aus `master` einen Branch.

Installiere den festgelegten Dokumentationsgenerator:

```sh
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
```

Sie können auch eine verifizierte Binärdatei aus dem [M-Press-Release v1.0.17](https://github.com/leaanthony/mpress/releases/tag/v1.0.17) herunterladen.

Führe im Stammverzeichnis des Wails-Repositorys Folgendes aus:

```sh
mpress version
mpress dev
```

Bearbeite die `.md`-Quelldateien in `docs/mpress/content/`. Englisch ist die Standardsprache und befindet sich direkt in diesem Verzeichnis. Vorhandene Übersetzungen befinden sich in Sprachordnern wie `fr/` und `id/`. Die Vorschau wird beim Speichern neu erstellt.

Behalte den Metadatenblock am Anfang jeder Seite sowie zusammengehörige `@...`- und `@end`-Komponenten bei. Gewöhnliche Absätze, Überschriften, Listen und abgegrenzte Codeblöcke kannst du als Text bearbeiten. Bearbeite keine generierten Dateien in `docs/mpress/site/`.

## Korrektur prüfen

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Prüfe die geänderte Seite im Browser und führe alle Codebeispiele aus, die du änderst. Vergleiche bei Übersetzungskorrekturen die vollständige geänderte Passage mit dem englischen Text. Behalte Befehle, API-Namen, Links, Codebeispiele und Diagrammverbindungen unverändert bei. Verwende natürliche technische Sprache, erhalte Anforderungen und Vorbehalte und übersetze neben dem Fließtext auch sichtbare Diagrammbeschriftungen, Navigationsbeschriftungen und Bildbeschreibungen.

Für jede veröffentlichte Sprache muss jede englische Seite vollständig übersetzt sein. Verwende keine englischen Platzhalter oder Ausweichseiten. Bei der Korrektur einer einzelnen Übersetzung kannst du nur diese Sprache ändern. Wenn du die Bedeutung des englischen Textes änderst, aktualisiere die entsprechenden Seiten in den anderen veröffentlichten Sprachen; nicht betroffene Seiten müssen nicht neu generiert werden.

## Pull Request einreichen

Eröffne einen PR gegen `master`. Beschreibe, was falsch war, erläutere deine Korrektur und führe die vorgenommenen Prüfungen auf. Füge bei sichtbaren Layoutänderungen Screenshots und bei geänderten Codebeispielen Angaben zu Plattform und Version hinzu.

Cloudflare-Zugangsdaten und private Dienste sind nicht erforderlich. Die öffentlichen PR-Prüfungen erstellen und validieren die statische Website ohne Zugangsdaten für die Bereitstellung.

Informationen zu Codeänderungen und Funktionsvorschlägen findest du unter [Zu Wails beitragen](/contributing/). Informationen zu den Interna findest du in der [Technischen Übersicht](/contributing/overview/).
