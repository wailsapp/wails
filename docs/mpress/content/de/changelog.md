---
title: "Änderungsprotokoll"
description: "Versionsverlauf und Versionshinweise für Wails v3"
slug: "changelog"
sourcePath: "changelog.md"
---

Legende:

-  – macOS
- ⊞ – Windows
- 🐧 – Linux

/_-- Alle nennenswerten Änderungen an diesem Projekt werden in dieser Datei dokumentiert.

Das Format basiert auf [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), und dieses Projekt folgt der [semantischen Versionierung](https://semver.org/spec/v2.0.0.html).

- `Added` für neue Funktionen.
- `Changed` für Änderungen an vorhandener Funktionalität.
- `Deprecated` für Funktionen, die demnächst entfernt werden.
- `Removed` für inzwischen entfernte Funktionen.
- `Fixed` für Fehlerbehebungen aller Art.
- `Security` bei Sicherheitslücken.

_/

/_   * AUTOMATISCH AKTUALISIERT *   Versionshinweise werden aus zusammengeführten PRs generiert. Bearbeiten Sie diese Datei nicht in regulären Beitrags-PRs und erstellen Sie die entfernte Datei `v3/UNRELEASED_CHANGELOG.md` nicht neu. Informationen zur Pflege des Änderungsprotokolls finden Sie unter `AGENTS.md`. _/

## [Unveröffentlicht]

## v3.2.0-nightly.20261103 – 2026-11-03

## Geändert
- [v3] Systray-Menüs bei Menu.Update aktualisieren, in [PR](https://github.com/wailsapp/wails/pull/6314) von @johnmaguire
- [v3/Windows] SetIcon für Fenster implementieren, in [PR](https://github.com/wailsapp/wails/pull/6502) von @hypherionmc

## v3.2.0-nightly.20261102 – 2026-11-02

## Geändert
- DisableOverscroll statt DisableBounce, in [PR](https://github.com/wailsapp/wails/pull/6077) von @lethusang

## v3.2.0-nightly.20261101 – 2026-11-01

## Hinzugefügt
- Experimentelle CGO-freie Backends für macOS und Linux hinzufügen, in [PR](https://github.com/wailsapp/wails/pull/6521) von @taliesin-ai
- Native YAML-Builds, portable Paketierung und Signierung, in [PR](https://github.com/wailsapp/wails/pull/6531) von @taliesin-ai

## Behoben
- Verschiedene Probleme mit v3-Nightly-Builds beheben und Abhängigkeiten bereinigen, in [PR](https://github.com/wailsapp/wails/pull/6519) von @taliesin-ai
- Validierung und Fixture-Klassifizierung für v3-Releases reparieren, in [PR](https://github.com/wailsapp/wails/pull/6533) von @taliesin-ai
- Plattformübergreifende Fehler bei der Release-Validierung beheben, in [PR](https://github.com/wailsapp/wails/pull/6538) von @taliesin-ai

## v3.1.0 – 2026-11-01

## Hinzugefügt
- Experimentelle CGO-freie Backends für macOS und Linux hinzufügen, in [PR](https://github.com/wailsapp/wails/pull/6521) von @taliesin-ai
- Native YAML-Builds, portable Paketierung und Signierung, in [PR](https://github.com/wailsapp/wails/pull/6531) von @taliesin-ai

## Geändert
- Verworfene Desktop-Metadaten nicht kopieren, in [PR](https://github.com/wailsapp/wails/pull/6391) von @taliesin-ai
- Inhaltslängen ohne fmt formatieren, in [PR](https://github.com/wailsapp/wails/pull/6393) von @taliesin-ai
- Indexwurzeln auswählen, ohne Kandidaten zu sammeln, in [PR](https://github.com/wailsapp/wails/pull/6404) von @taliesin-ai
- Kopierpuffer für nicht durchsuchbare Dateien wiederverwenden, in [PR](https://github.com/wailsapp/wails/pull/6408) von @taliesin-ai
- Keinen Reader zum Ermitteln von Präfixen allozieren, in [PR](https://github.com/wailsapp/wails/pull/6410) von @taliesin-ai
- Lesevorgänge zur MIME-Erkennung bei bekannten Erweiterungen überspringen, in [PR](https://github.com/wailsapp/wails/pull/6414) von @taliesin-ai
- Wachstum der Warteschlange und Oszillation beim Verkleinern verhindern, in [PR](https://github.com/wailsapp/wails/pull/6415) von @taliesin-ai
- Verworfene Anfragekontexte nicht allozieren, in [PR](https://github.com/wailsapp/wails/pull/6419) von @taliesin-ai
- Vorbereitung deaktivierter Debug-Protokolle überspringen, in [PR](https://github.com/wailsapp/wails/pull/6420) von @taliesin-ai
- Puffer zur MIME-Erkennung verzögert allozieren, in [PR](https://github.com/wailsapp/wails/pull/6428) von @taliesin-ai
- Allokationen bei der Formatierung von Hyperlinks reduzieren, in [PR](https://github.com/wailsapp/wails/pull/6443) von @taliesin-ai
- Temporäre Allokationen für die Formatierung von Aufzählungszeichen entfernen, in [PR](https://github.com/wailsapp/wails/pull/6449) von @taliesin-ai
- UUID-Stringkonvertierung auf eine Allokation reduzieren, in [PR](https://github.com/wailsapp/wails/pull/6453) von @taliesin-ai
- Allokationen bei der Formatierung von Statement-Handles reduzieren, in [PR](https://github.com/wailsapp/wails/pull/6454) von @taliesin-ai
- Temporäre Digest-Allokation für namensbasierte UUIDs beseitigen, in [PR](https://github.com/wailsapp/wails/pull/6455) von @taliesin-ai
- Unvollständige Zeilen der Befehlsausgabe nicht kopieren, in [PR](https://github.com/wailsapp/wails/pull/6460) von @taliesin-ai
- Kopien ausführbarer Dateien während des Paket-Stagings streamen, in [PR](https://github.com/wailsapp/wails/pull/6472) von @taliesin-ai
- Unvollständige Ausgabezeilen nicht erneut durchsuchen, in [PR](https://github.com/wailsapp/wails/pull/6473) von @taliesin-ai
- Kompilierte Muster für das Befehlsrouting wiederverwenden, in [PR](https://github.com/wailsapp/wails/pull/6475) von @taliesin-ai
- Metadatenabfragen während der Runtime-Erkennung vermeiden, in [PR](https://github.com/wailsapp/wails/pull/6482) von @taliesin-ai
- Allokationen lokalisierter Fallback-Seiten reduzieren, in [PR](https://github.com/wailsapp/wails/pull/6485) von @taliesin-ai
- Muster nicht für jeden Prozessor neu kompilieren, in [PR](https://github.com/wailsapp/wails/pull/6490) von @taliesin-ai

## Behoben
- SetTooltip für den System-Tray implementieren, in [PR](https://github.com/wailsapp/wails/pull/6017) von @fan711
- Dunkelmodus aus dem Einstellungsportal auslesen, in [PR](https://github.com/wailsapp/wails/pull/6072) von @TechHutTV
- Window.SetZoom Werte unter 1.0 erlauben, in [PR](https://github.com/wailsapp/wails/pull/6166) von @APshenkin
- In Emit-Reihenfolge an jeden Go-Listener zustellen, in [PR](https://github.com/wailsapp/wails/pull/6275) von @DevLumuz
- Event-Hooks in Emit unter hookLock lesen, in [PR](https://github.com/wailsapp/wails/pull/6276) von @DevLumuz
- Die Tray-Eigenschaft Menu immer bekannt geben, damit ein nach Run erstelltes Tray nicht verworfen wird, in [PR](https://github.com/wailsapp/wails/pull/6284) von @AlbinoGeek
- Vor Run gesendete Arbeit für den Hauptthread zurückhalten, bis Run startet, in [PR](https://github.com/wailsapp/wails/pull/6287) von @AlbinoGeek
- Korrekte Asset-URL im GitHub-Updater-Provider verwenden, in [PR](https://github.com/wailsapp/wails/pull/6304) von @gavidroselj
- Anfragetext beibehalten, wenn WebView2 Content-Length auslässt, in [PR](https://github.com/wailsapp/wails/pull/6326) von @strraaa
- macOS- und GTK-Bildschirme im DIP-Raum melden (Begrenzungen bei gemischten DPI-Werten), in [PR](https://github.com/wailsapp/wails/pull/6346) von @julianstorer
- Eine verworfene Benachrichtigung nicht als Klick melden, in [PR](https://github.com/wailsapp/wails/pull/6367) von @bercedev
- Entwicklungsprozesse beim Schließen des Terminals beenden, in [PR](https://github.com/wailsapp/wails/pull/6371) von @xmasdev
- Änderungsprotokoll-Läufe bei der Nachweisprüfung Squash-Merge-Commits zuordnen, in [PR](https://github.com/wailsapp/wails/pull/6373) von @taliesin-ai
- Entprellung für Taskleisten-Neustarts eine Sekunde lang aufrechterhalten in [PR](https://github.com/wailsapp/wails/pull/6386) von @taliesin-ai
- Größenbeschränkungen über einen dort gestarteten Vollbildmodus hinweg beibehalten in [PR](https://github.com/wailsapp/wails/pull/6387) von @taliesin-ai
- Vollbild- oder maximierte Fenster nicht an neue Beschränkungen anpassen in [PR](https://github.com/wailsapp/wails/pull/6388) von @taliesin-ai
- im Liquid-Glass-Beispiel vor Run keinen Dialog öffnen in [PR](https://github.com/wailsapp/wails/pull/6389) von @taliesin-ai
- vor Run angeforderte Dialoge wirkungslos lassen, statt endlos zu warten in [PR](https://github.com/wailsapp/wails/pull/6390) von @taliesin-ai
- Fallback verwenden, wenn ein Desktop-Eintrag keinen Befehl enthält in [PR](https://github.com/wailsapp/wails/pull/6392) von @taliesin-ai
- aktuellen Inhalt ohne Beibehaltung von Dateinamen erkennen in [PR](https://github.com/wailsapp/wails/pull/6394) von @taliesin-ai
- Asset-Erweiterungen mit gemischter Groß- und Kleinschreibung erkennen in [PR](https://github.com/wailsapp/wails/pull/6395) von @taliesin-ai
- verbrauchte Puffer der Inhaltserkennung freigeben in [PR](https://github.com/wailsapp/wails/pull/6396) von @taliesin-ai
- entfernte ringqueue-Referenzen freigeben in [PR](https://github.com/wailsapp/wails/pull/6397) von @taliesin-ai
- gelöschte Header vor der Fallback-Auswahl berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6398) von @taliesin-ai
- Fallback-Antwortheader vor dem Flush festschreiben in [PR](https://github.com/wailsapp/wails/pull/6399) von @taliesin-ai
- nur exakt übereinstimmende Indexdateien erkennen in [PR](https://github.com/wailsapp/wails/pull/6400) von @taliesin-ai
- bei nicht durchsuchbaren HEAD-Antworten Body und vollständiges Lesen vermeiden in [PR](https://github.com/wailsapp/wails/pull/6401) von @taliesin-ai
- Präfix für die MIME-Erkennung über kurze Lesevorgänge hinweg füllen in [PR](https://github.com/wailsapp/wails/pull/6402) von @taliesin-ai
- erfolgreiche leere, nicht durchsuchbare Antworten festschreiben in [PR](https://github.com/wailsapp/wails/pull/6403) von @taliesin-ai
- Duplizieren erweiterter URI-Ziele vermeiden in [PR](https://github.com/wailsapp/wails/pull/6405) von @taliesin-ai
- Pfade bei der Felderweiterung beibehalten in [PR](https://github.com/wailsapp/wails/pull/6406) von @taliesin-ai
- XDG-Suchverzeichnisse für Desktop-Dateien berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6407) von @taliesin-ai
- Header der Inhaltserkennung bei einem leeren Flush festschreiben in [PR](https://github.com/wailsapp/wails/pull/6409) von @taliesin-ai
- Verzeichnisse überspringen, die Desktop-Einträge verdecken in [PR](https://github.com/wailsapp/wails/pull/6411) von @taliesin-ai
- in Pfaden ausführbarer Dateien gefundene Ziele beibehalten in [PR](https://github.com/wailsapp/wails/pull/6412) von @taliesin-ai
- Callback-Captures nach Auslösen des Timers freigeben in [PR](https://github.com/wailsapp/wails/pull/6413) von @taliesin-ai
- falsche ldconfig-Bibliothekstreffer vermeiden in [PR](https://github.com/wailsapp/wails/pull/6416) von @taliesin-ai
- Fallback-Treffer für Bibliotheken im Dateisystem einschränken in [PR](https://github.com/wailsapp/wails/pull/6417) von @taliesin-ai
- Datenrennen bei der Cache-Invalidierung beseitigen in [PR](https://github.com/wailsapp/wails/pull/6418) von @taliesin-ai
- optionale Logger-Konfiguration berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6421) von @taliesin-ai
- gemeinsame Mutation der Start-URL und veraltete Ports vermeiden in [PR](https://github.com/wailsapp/wails/pull/6422) von @taliesin-ai
- Wiederholungsversuche des Entwicklungs-Proxys beim Abbruch einer Anfrage beenden in [PR](https://github.com/wailsapp/wails/pull/6423) von @taliesin-ai
- bei fehlenden gebündelten Assets 404 zurückgeben in [PR](https://github.com/wailsapp/wails/pull/6424) von @taliesin-ai
- HEAD-Anfragen an die gebündelte Runtime verarbeiten in [PR](https://github.com/wailsapp/wails/pull/6425) von @taliesin-ai
- maskierte Pfade beim Routing von Diensten beibehalten in [PR](https://github.com/wailsapp/wails/pull/6426) von @taliesin-ai
- überlange numerische Vorabversionsbezeichner korrekt sortieren in [PR](https://github.com/wailsapp/wails/pull/6427) von @taliesin-ai
- HTTP-Fehler vor dem Speichern heruntergeladener Werkzeuge zurückweisen in [PR](https://github.com/wailsapp/wails/pull/6429) von @taliesin-ai
- Startfrist bei der Desktop-Erkennung berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6430) von @taliesin-ai
- Anfrage-Bodys bei Wiederholungsversuchen des Entwicklungs-Proxys beibehalten in [PR](https://github.com/wailsapp/wails/pull/6432) von @taliesin-ai
- HTTP-Antwort-Bodys auf jedem Rückgabepfad schließen in [PR](https://github.com/wailsapp/wails/pull/6434) von @taliesin-ai
- Abbruchfehler bei parallelen Bibliothekssuchen beibehalten in [PR](https://github.com/wailsapp/wails/pull/6435) von @taliesin-ai
- fehlgeschlagene Antworten auf GitHub-Tag-Anfragen vor dem Dekodieren zurückweisen in [PR](https://github.com/wailsapp/wails/pull/6436) von @taliesin-ai
- literale, in Klammern eingeschlossene Vorlagenpfade beibehalten in [PR](https://github.com/wailsapp/wails/pull/6437) von @taliesin-ai
- wails3 einmal installieren und Installationsfehler beibehalten in [PR](https://github.com/wailsapp/wails/pull/6438) von @taliesin-ai
- versteckte Einträge und Überschriften in Kommentarblöcken ignorieren in [PR](https://github.com/wailsapp/wails/pull/6439) von @taliesin-ai
- Fehler bei der Suche nach stabilen Releases sichtbar machen in [PR](https://github.com/wailsapp/wails/pull/6440) von @taliesin-ai
- reine Textnachrichten ohne Formatinterpretation ausgeben in [PR](https://github.com/wailsapp/wails/pull/6441) von @taliesin-ai
- Schreib-Panics nach dem Laden von JSON-null vermeiden in [PR](https://github.com/wailsapp/wails/pull/6442) von @taliesin-ai
- Snapshots für Lesevorgänge des gesamten Stores zurückgeben in [PR](https://github.com/wailsapp/wails/pull/6444) von @taliesin-ai
- Fehler nach teilweiser Zeileniteration weitergeben in [PR](https://github.com/wailsapp/wails/pull/6446) von @taliesin-ai
- Zielpfade beim Extrahieren literal beibehalten in [PR](https://github.com/wailsapp/wails/pull/6447) von @taliesin-ai
- Handles vorbereiteter Anweisungen beim Schließen und erneuten Öffnen freigeben in [PR](https://github.com/wailsapp/wails/pull/6448) von @taliesin-ai
- Datenrennen zwischen Vorbereitung und Schließen oder erneutem Öffnen behandeln in [PR](https://github.com/wailsapp/wails/pull/6450) von @taliesin-ai
- HEAD-Anfragen für die Fallback-Indexseite berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6451) von @taliesin-ai
- vollständige hexadezimale Anweisungs-Handles verlangen in [PR](https://github.com/wailsapp/wails/pull/6452) von @taliesin-ai
- standardmäßige In-Memory-Datenbank bei Nebenläufigkeit beibehalten in [PR](https://github.com/wailsapp/wails/pull/6456) von @taliesin-ai
- stderr von der analysierten HEAD-Ausgabe trennen in [PR](https://github.com/wailsapp/wails/pull/6457) von @taliesin-ai
- Runtime-Assets aus eingebetteten Bundles generieren in [PR](https://github.com/wailsapp/wails/pull/6458) von @taliesin-ai
- TODO-, FIXME- und TBD-Platzhalter erkennen in [PR](https://github.com/wailsapp/wails/pull/6459) von @taliesin-ai
- API-Antworten für Release Notes validieren in [PR](https://github.com/wailsapp/wails/pull/6461) von @taliesin-ai
- Task-Verzeichnisse bei Statusprüfungen berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6462) von @taliesin-ai
- Datenverlust beim Kopieren einer Datei auf sich selbst verhindern in [PR](https://github.com/wailsapp/wails/pull/6463) von @taliesin-ai
- Task-Verzeichnis in die Cache-Identität einbeziehen in [PR](https://github.com/wailsapp/wails/pull/6464) von @taliesin-ai
- Ports pro Prüfung nur einmal abfragen in [PR](https://github.com/wailsapp/wails/pull/6465) von @taliesin-ai
- Panic vermeiden, wenn Cache-Einträge als null dekodiert werden in [PR](https://github.com/wailsapp/wails/pull/6466) von @taliesin-ai
- numerische Vorabversionen ohne Ganzzahlgrenzen erhöhen in [PR](https://github.com/wailsapp/wails/pull/6467) von @taliesin-ai
- verhindern, dass Legacy-Prüfungen Cache-Fehlschläge verbergen, in [PR](https://github.com/wailsapp/wails/pull/6468) von @taliesin-ai
- Standardaufgaben im Basisverzeichnis des Executors ausführen in [PR](https://github.com/wailsapp/wails/pull/6469) von @taliesin-ai
- literale Prozentzeichen in Explorer-Befehlen beibehalten in [PR](https://github.com/wailsapp/wails/pull/6470) von @taliesin-ai
- Verzeichnisse für Vorbedingungen und Aufgabenbefehle angleichen in [PR](https://github.com/wailsapp/wails/pull/6471) von @taliesin-ai
- symbolische Links durch generierte reguläre Dateien ersetzen in [PR](https://github.com/wailsapp/wails/pull/6474) von @taliesin-ai
- vorhandene DMGs bei fehlgeschlagenen Builds beibehalten in [PR](https://github.com/wailsapp/wails/pull/6476) von @taliesin-ai
- Backslashes in Pfaden und in regulären Ausdrücken in Anführungszeichen beibehalten in [PR](https://github.com/wailsapp/wails/pull/6477) von @taliesin-ai
- Größen generierter Symbole vor der Bildallokation validieren in [PR](https://github.com/wailsapp/wails/pull/6478) von @taliesin-ai
- Befehlsketten und Umleitungen über die Shell ausführen in [PR](https://github.com/wailsapp/wails/pull/6479) von @taliesin-ai
- vorherige Linux-Pakete bei fehlgeschlagenen Builds beibehalten in [PR](https://github.com/wailsapp/wails/pull/6480) von @taliesin-ai
- plattformspezifische Shell-Ausführung als Rückfalloption verwenden in [PR](https://github.com/wailsapp/wails/pull/6481) von @taliesin-ai
- Assets mit wiederholten führenden Schrägstrichen im Pfad ausliefern in [PR](https://github.com/wailsapp/wails/pull/6483) von @taliesin-ai
- stabile Versionen mit Build-Metadaten erkennen in [PR](https://github.com/wailsapp/wails/pull/6484) von @taliesin-ai
- Fehler bei der Ressourcengenerierung korrekt melden in [PR](https://github.com/wailsapp/wails/pull/6486) von @taliesin-ai
- Race Conditions beim Ersetzen des aktiven Reporters verhindern in [PR](https://github.com/wailsapp/wails/pull/6487) von @taliesin-ai
- gleichzeitige Speichervorgänge für Cache-Dateien serialisieren in [PR](https://github.com/wailsapp/wails/pull/6488) von @taliesin-ai
- mehrere Gleichheitskandidaten für Vorlagen unterstützen in [PR](https://github.com/wailsapp/wails/pull/6489) von @taliesin-ai
- gespeicherte Kategorietabellen mit dem Wert null normalisieren in [PR](https://github.com/wailsapp/wails/pull/6491) von @taliesin-ai
- Diagnosen fehlgeschlagener variabler Shell-Befehle anzeigen in [PR](https://github.com/wailsapp/wails/pull/6492) von @taliesin-ai
- Shell-Befehlsketten aus dem impliziten Go-Cache ausschließen in [PR](https://github.com/wailsapp/wails/pull/6493) von @taliesin-ai
- Persistierung von Linux-Kategorien serialisieren in [PR](https://github.com/wailsapp/wails/pull/6494) von @taliesin-ai
- Go-Ausgaben mit Leerzeichen in Anführungszeichen oder maskierten Leerzeichen zwischenspeichern in [PR](https://github.com/wailsapp/wails/pull/6495) von @taliesin-ai
- Mac-Symbole bei fehlgeschlagener Codierung beibehalten in [PR](https://github.com/wailsapp/wails/pull/6496) von @taliesin-ai
- ungültige Go-Test-Flags für die Ausführungsanzahl nicht mehr verwerfen in [PR](https://github.com/wailsapp/wails/pull/6497) von @taliesin-ai
- Quellberechtigungen bei neuen Dateikopien beibehalten in [PR](https://github.com/wailsapp/wails/pull/6498) von @taliesin-ai
- verschiedene Probleme mit v3 Nightly-Builds beheben und Abhängigkeiten bereinigen in [PR](https://github.com/wailsapp/wails/pull/6519) von @taliesin-ai
- Validierung von v3 Releases und Klassifizierung von Fixtures korrigieren in [PR](https://github.com/wailsapp/wails/pull/6533) von @taliesin-ai
- Fehler bei der plattformübergreifenden Release-Validierung beheben in [PR](https://github.com/wailsapp/wails/pull/6538) von @taliesin-ai

## v3.1.0-nightly.20261028 - 2026-10-28

## Geändert
- Changelog-Generierung durch Wiederverwendung kompilierter Muster für reguläre Ausdrücke beschleunigen in [PR](https://github.com/wailsapp/wails/pull/6490) von @taliesin-ai
- Allokationen beim Ausliefern lokalisierter Rückfallseiten reduzieren in [PR](https://github.com/wailsapp/wails/pull/6485) von @taliesin-ai
- Ermittlung der AppImage-Laufzeit durch Vermeidung unnötiger Metadatenabfragen beschleunigen in [PR](https://github.com/wailsapp/wails/pull/6482) von @taliesin-ai
- kompilierte Muster für das Befehlsrouting wiederverwenden, um die Wake-Leistung zu verbessern, in [PR](https://github.com/wailsapp/wails/pull/6475) von @taliesin-ai
- Leistung verbessern, indem unvollständige Ausgabezeilen nicht erneut durchsucht werden, in [PR](https://github.com/wailsapp/wails/pull/6473) von @taliesin-ai
- Kopien ausführbarer Dateien während des Stagings von MSIX-Paketen streamen, um den Speicherverbrauch zu reduzieren und Builds zu beschleunigen, in [PR](https://github.com/wailsapp/wails/pull/6472) von @taliesin-ai
- Kopieren unvollständiger Befehlsausgabezeilen vermeiden, um die Leistung zu verbessern, in [PR](https://github.com/wailsapp/wails/pull/6460) von @taliesin-ai
- Allokationen beim Generieren namensbasierter UUIDs reduzieren in [PR](https://github.com/wailsapp/wails/pull/6455) von @taliesin-ai
- Allokationen beim Formatieren von SQLite-Anweisungshandles reduzieren in [PR](https://github.com/wailsapp/wails/pull/6454) von @taliesin-ai
- UUID-Zeichenfolgen mit nur einer Allokation konvertieren in [PR](https://github.com/wailsapp/wails/pull/6453) von @taliesin-ai
- Allokationen beim Formatieren von Changelog-Aufzählungseinträgen reduzieren in [PR](https://github.com/wailsapp/wails/pull/6449) von @taliesin-ai
- Allokationen beim Formatieren von Terminal-Hyperlinks reduzieren in [PR](https://github.com/wailsapp/wails/pull/6443) von @taliesin-ai
- Allokationen beim Ausliefern von Assets durch verzögerte MIME-Erkennung in Puffern reduzieren in [PR](https://github.com/wailsapp/wails/pull/6428) von @taliesin-ai
- Vorbereitung von Debug-Protokollen im Asset-Server überspringen, wenn die Debug-Protokollierung deaktiviert ist, in [PR](https://github.com/wailsapp/wails/pull/6420) von @taliesin-ai
- Allokationen bei der Verarbeitung von WebView-Anfragen reduzieren, indem verworfene Anfragekontexte vermieden werden, in [PR](https://github.com/wailsapp/wails/pull/6419) von @taliesin-ai
- Auslieferung von Assets beschleunigen, indem die MIME-Erkennung bei bekannten Dateierweiterungen übersprungen wird, in [PR](https://github.com/wailsapp/wails/pull/6414) von @taliesin-ai
- Allokationen beim Ausliefern von Assets reduzieren, indem kein Reader für die Inhaltserkennung verwendet wird, in [PR](https://github.com/wailsapp/wails/pull/6410) von @taliesin-ai
- Allokationen beim Ausliefern nicht durchsuchbarer Asset-Dateien reduzieren in [PR](https://github.com/wailsapp/wails/pull/6408) von @taliesin-ai
- Auswahl des Index-Stammverzeichnisses im Asset-Server durch Vermeidung einer Kandidatensammlung beschleunigen in [PR](https://github.com/wailsapp/wails/pull/6404) von @taliesin-ai
- Leistung des Asset-Servers verbessern, indem Inhaltslängen ohne fmt formatiert werden, in [PR](https://github.com/wailsapp/wails/pull/6393) von @taliesin-ai
- Speicherkopiervorgänge beim Lesen von Metadaten des Desktop-Datei-Explorers reduzieren in [PR](https://github.com/wailsapp/wails/pull/6391) von @taliesin-ai

## Behoben
- Berechtigungen der Quelldatei beim Kopieren in neue Zieldateien beibehalten in [PR](https://github.com/wailsapp/wails/pull/6498) von @taliesin-ai
- Größenänderung von Vollbild- und maximierten Fenstern verhindern, wenn sich die Beschränkungen für die minimale oder maximale Größe ändern, in [PR](https://github.com/wailsapp/wails/pull/6388) von @taliesin-ai
- geplante Aktualisierungen des Sternverlaufs für das Repository wailsapp/wails veröffentlichen in [PR](https://github.com/wailsapp/wails/pull/6500) von @taliesin-ai
- ungültige Go-Test-Flags für die Ausführungsanzahl in Wake nicht mehr verwerfen in [PR](https://github.com/wailsapp/wails/pull/6497) von @taliesin-ai
- Mac-Symbole bei fehlgeschlagener Codierung beibehalten in [PR](https://github.com/wailsapp/wails/pull/6496) von @taliesin-ai
- Build-Cache-Fehlschläge bei Go-Ausgaben mit Leerzeichen in Anführungszeichen oder maskierten Leerzeichen beheben in [PR](https://github.com/wailsapp/wails/pull/6495) von @taliesin-ai
- Persistierung von Benachrichtigungskategorien unter Linux serialisieren, um Race Conditions bei gleichzeitigen Schreibvorgängen zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/6494) von @taliesin-ai
- Shell-Befehlsketten aus dem impliziten Go-Build-Cache ausschließen in [PR](https://github.com/wailsapp/wails/pull/6493) von @taliesin-ai
- Fehler beim Laden von Benachrichtigungskategorie-Zuordnungen beheben, wenn die gespeicherten Daten null sind, in [PR](https://github.com/wailsapp/wails/pull/6491) von @taliesin-ai
- Diagnosen fehlgeschlagener variabler Shell-Befehle anzeigen in [PR](https://github.com/wailsapp/wails/pull/6492) von @taliesin-ai
- Wake-Vorlagenabgleich bei mehreren Gleichheitskandidaten korrigieren in [PR](https://github.com/wailsapp/wails/pull/6489) von @taliesin-ai
- Gleichzeitiges Speichern von Cache-Dateien serialisieren, um beschädigte Cache-Schreibvorgänge zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/6488) von @taliesin-ai
- Race Conditions beim Ersetzen des aktiven Reporters verhindern in [PR](https://github.com/wailsapp/wails/pull/6487) von @taliesin-ai
- Fehler bei der Ressourcengenerierung in der CLI korrekt melden in [PR](https://github.com/wailsapp/wails/pull/6486) von @taliesin-ai
- Stabile Versionen mit Build-Metadaten in der CLI erkennen in [PR](https://github.com/wailsapp/wails/pull/6484) von @taliesin-ai
- Assets ausliefern, deren Anfragepfade wiederholte führende Schrägstriche enthalten, in [PR](https://github.com/wailsapp/wails/pull/6483) von @taliesin-ai
- Shell-Fallback-Ausführung für Wake auf jeder Plattform korrigieren in [PR](https://github.com/wailsapp/wails/pull/6481) von @taliesin-ai
- Vorhandene Linux-Pakete bei einem fehlgeschlagenen Build beibehalten in [PR](https://github.com/wailsapp/wails/pull/6480) von @taliesin-ai
- Befehlsketten und Umleitungen über die Shell ausführen in [PR](https://github.com/wailsapp/wails/pull/6479) von @taliesin-ai
- Backslashes in Anführungszeichen gesetzter regulärer Ausdrücke und Pfade beim Parsen von Wake-Befehlen beibehalten in [PR](https://github.com/wailsapp/wails/pull/6477) von @taliesin-ai
- Vorhandene DMG-Dateien bei fehlgeschlagenen Paketierungs-Builds beibehalten in [PR](https://github.com/wailsapp/wails/pull/6476) von @taliesin-ai
- Symbolische Links in der Bindings-Ausgabe durch generierte reguläre Dateien ersetzen in [PR](https://github.com/wailsapp/wails/pull/6474) von @taliesin-ai
- Literale Prozentzeichen in Linux-Dateiexplorer-Befehlen beibehalten in [PR](https://github.com/wailsapp/wails/pull/6470) von @taliesin-ai
- Verzeichnisse von Vorbedingungen und Task-Befehlen angleichen, damit sie konsistent aufgelöst werden, in [PR](https://github.com/wailsapp/wails/pull/6471) von @taliesin-ai
- Standard-Tasks im Basisverzeichnis des Executors ausführen in [PR](https://github.com/wailsapp/wails/pull/6469) von @taliesin-ai
- Verhindern, dass Legacy-Prüfungen Cache-Fehlschläge verbergen, in [PR](https://github.com/wailsapp/wails/pull/6468) von @taliesin-ai
- Nummerierung von Vorabversionen in der CLI korrigieren, sodass sie ohne Begrenzung durch Ganzzahlüberläufe erhöht wird, in [PR](https://github.com/wailsapp/wails/pull/6467) von @taliesin-ai
- Panic vermeiden, wenn Cache-Einträge als null dekodiert werden, in [PR](https://github.com/wailsapp/wails/pull/6466) von @taliesin-ai
- Ports in der CLI nur einmal pro Prüfung testen in [PR](https://github.com/wailsapp/wails/pull/6465) von @taliesin-ai
- Task-Verzeichnis in die Cache-Identität aufnehmen, um veraltete Wake-Ergebnisse zu vermeiden, in [PR](https://github.com/wailsapp/wails/pull/6464) von @taliesin-ai
- Datenverlust verhindern, wenn bei der Paketierung eine Datei auf sich selbst kopiert wird, in [PR](https://github.com/wailsapp/wails/pull/6463) von @taliesin-ai
- Task-Verzeichnisse bei Statusprüfungen berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6462) von @taliesin-ai
- Platzhalter TODO, FIXME und TBD in Changelog-Einträgen erkennen in [PR](https://github.com/wailsapp/wails/pull/6459) von @taliesin-ai
- Runtime-Assets in der CLI aus eingebetteten Bundles generieren in [PR](https://github.com/wailsapp/wails/pull/6458) von @taliesin-ai
- stderr in Git-Befehlen von der geparsten HEAD-Ausgabe trennen in [PR](https://github.com/wailsapp/wails/pull/6457) von @taliesin-ai
- Standardmäßige SQLite-In-Memory-Datenbank bei gleichzeitigem Zugriff beibehalten in [PR](https://github.com/wailsapp/wails/pull/6456) von @taliesin-ai
- Korrekte Antworten auf HEAD-Anfragen für die Fallback-Indexseite zurückgeben in [PR](https://github.com/wailsapp/wails/pull/6451) von @taliesin-ai
- Race Condition zwischen dem Vorbereiten von SQLite-Anweisungen und dem Schließen oder erneuten Öffnen der Datenbank behandeln in [PR](https://github.com/wailsapp/wails/pull/6450) von @taliesin-ai
- Handles vorbereiteter Anweisungen freigeben, wenn die SQLite-Datenbank geschlossen oder erneut geöffnet wird, in [PR](https://github.com/wailsapp/wails/pull/6448) von @taliesin-ai
- Zielpfade für die Extraktion bei Verwendung von gosod literal beibehalten in [PR](https://github.com/wailsapp/wails/pull/6447) von @taliesin-ai
- Fehler aus SQLite-Abfragen nach teilweiser Zeileniteration weitergeben in [PR](https://github.com/wailsapp/wails/pull/6446) von @taliesin-ai
- Bei Lesevorgängen des gesamten Speichers in kvstore Snapshots zurückgeben in [PR](https://github.com/wailsapp/wails/pull/6444) von @taliesin-ai
- Schreib-Panics in kvstore nach dem Laden eines JSON-null-Werts verhindern in [PR](https://github.com/wailsapp/wails/pull/6442) von @taliesin-ai
- Einfache Terminalmeldungen ausgeben, ohne Formatierungsverben zu interpretieren, in [PR](https://github.com/wailsapp/wails/pull/6441) von @taliesin-ai
- Fehler beim Ermitteln der neuesten stabilen Version im Updater sichtbar machen in [PR](https://github.com/wailsapp/wails/pull/6440) von @taliesin-ai
- Ausgeblendete Einträge und Überschriften in Changelog-Kommentarblöcken ignorieren in [PR](https://github.com/wailsapp/wails/pull/6439) von @taliesin-ai
- wails3 nur einmal installieren und Installationsfehler im Updater sichtbar machen in [PR](https://github.com/wailsapp/wails/pull/6438) von @taliesin-ai
- Literale Vorlagenpfade in eckigen Klammern in gosod beibehalten in [PR](https://github.com/wailsapp/wails/pull/6437) von @taliesin-ai
- Fehlgeschlagene Antworten auf GitHub-Tags im Updater vor dem Dekodieren ablehnen in [PR](https://github.com/wailsapp/wails/pull/6436) von @taliesin-ai
- Abbruchfehler bei parallelen Bibliothekssuchen beibehalten in [PR](https://github.com/wailsapp/wails/pull/6435) von @taliesin-ai
- HTTP-Antwortkörper in jedem Rückgabepfad schließen, um Ressourcenlecks zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/6434) von @taliesin-ai
- Anfragekörper beibehalten, wenn der Entwicklungsproxy Anfragen erneut versucht, in [PR](https://github.com/wailsapp/wails/pull/6432) von @taliesin-ai
- Startfrist bei der Ermittlung des Desktop-Dateiexplorers einhalten in [PR](https://github.com/wailsapp/wails/pull/6430) von @taliesin-ai
- HTTP-Fehler vor dem Speichern heruntergeladener AppImage-Werkzeuge ablehnen in [PR](https://github.com/wailsapp/wails/pull/6429) von @taliesin-ai
- Überlange numerische Vorabversionsbezeichner in SemVer-Vergleichen korrekt sortieren in [PR](https://github.com/wailsapp/wails/pull/6427) von @taliesin-ai
- Maskierte Pfade beim Routing von Diensten im Asset-Server beibehalten in [PR](https://github.com/wailsapp/wails/pull/6426) von @taliesin-ai
- HEAD-Anfragen für die gebündelte Runtime im Asset-Server verarbeiten in [PR](https://github.com/wailsapp/wails/pull/6425) von @taliesin-ai
- Für fehlende gebündelte Assets 404 statt einer falschen Antwort zurückgeben in [PR](https://github.com/wailsapp/wails/pull/6424) von @taliesin-ai
- Wiederholungsversuche des Entwicklungsproxys beenden, wenn eine Anfrage abgebrochen wird, in [PR](https://github.com/wailsapp/wails/pull/6423) von @taliesin-ai
- Gemeinsame Mutation der Start-URL und veraltete Ports im Asset-Server vermeiden in [PR](https://github.com/wailsapp/wails/pull/6422) von @taliesin-ai
- Optionale Logger-Konfiguration im Asset-Server berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6421) von @taliesin-ai
- Data Race bei der Invalidierung des Bibliothekspfad-Caches beseitigen in [PR](https://github.com/wailsapp/wails/pull/6418) von @taliesin-ai
- Falsche Bibliothekstreffer beim Durchsuchen der ldconfig-Ausgabe vermeiden in [PR](https://github.com/wailsapp/wails/pull/6416) von @taliesin-ai
- Wachstum der Asset-Server-Warteschlange und Oszillation beim Verkleinern unter Last verhindern in [PR](https://github.com/wailsapp/wails/pull/6415) von @taliesin-ai
- Vom Callback erfasste Werte nach Auslösen des Debounce-Timers freigeben, damit veraltete Referenzen nicht beibehalten werden, in [PR](https://github.com/wailsapp/wails/pull/6413) von @taliesin-ai
- In Pfaden ausführbarer Dateien gefundene Dateiexplorer-Ziele beibehalten in [PR](https://github.com/wailsapp/wails/pull/6412) von @taliesin-ai
- Verzeichnisse überspringen, die Desktop-Einträge im Datei-Explorer verdecken, in [PR](https://github.com/wailsapp/wails/pull/6411) von @taliesin-ai
- Fehler behoben, durch den Asset-Server-Header bei leerem Antworttext nicht gesendet wurden, in [PR](https://github.com/wailsapp/wails/pull/6409) von @taliesin-ai
- XDG-Suchverzeichnisse für Desktop-Dateien bei der Suche nach Anwendungssymbolen berücksichtigen, in [PR](https://github.com/wailsapp/wails/pull/6407) von @taliesin-ai
- Dateipfade bei der Felderweiterung im Datei-Explorer beibehalten, in [PR](https://github.com/wailsapp/wails/pull/6406) von @taliesin-ai
- Duplizierte erweiterte URI-Ziele im Datei-Explorer vermeiden, in [PR](https://github.com/wailsapp/wails/pull/6405) von @taliesin-ai
- Erfolgreiche leere, nicht durchsuchbare Asset-Antworten abschließen, statt sie ausstehend zu lassen, in [PR](https://github.com/wailsapp/wails/pull/6403) von @taliesin-ai
- Präfix für die MIME-Erkennung im Asset-Server auch bei kurzen Lesevorgängen korrekt füllen, in [PR](https://github.com/wailsapp/wails/pull/6402) von @taliesin-ai
- Bei nicht durchsuchbaren HEAD-Anfragen im Asset-Server nicht den gesamten Antworttext lesen, in [PR](https://github.com/wailsapp/wails/pull/6401) von @taliesin-ai
- Beim Auflösen von Verzeichnisanfragen im Asset-Server nur exakt passende Indexdateien ausliefern, in [PR](https://github.com/wailsapp/wails/pull/6400) von @taliesin-ai
- Fehler behoben, durch den verworfene Linux-Benachrichtigungen auf Servern ohne Aktionsunterstützung als Klicks gemeldet wurden, in [PR](https://github.com/wailsapp/wails/pull/6367) von @bercedev
- Entwicklungsprozesse beim Schließen des Terminals beenden, in [PR](https://github.com/wailsapp/wails/pull/6371) von @xmasdev
- Debounce für den Neustart der Windows-Taskleiste eine Sekunde lang aufrechterhalten, bevor das Flag gelöscht wird, in [PR](https://github.com/wailsapp/wails/pull/6386) von @taliesin-ai
- Fehler behoben, durch den Header von Fallback-Antworten im Asset-Server verloren gingen, in [PR](https://github.com/wailsapp/wails/pull/6399) von @taliesin-ai
- Gelöschte Header vor der Auswahl eines Fallbacks im Asset-Server berücksichtigen, in [PR](https://github.com/wailsapp/wails/pull/6398) von @taliesin-ai
- Referenzen auf entfernte Ringpuffer-Einträge im Asset-Server freigeben, um Speicherlecks zu vermeiden, in [PR](https://github.com/wailsapp/wails/pull/6397) von @taliesin-ai
- Verbrauchte Puffer der Inhaltserkennung im Asset-Server freigeben, um den Speicherverbrauch zu senken, in [PR](https://github.com/wailsapp/wails/pull/6396) von @taliesin-ai
- Asset-Dateierweiterungen mit gemischter Groß- und Kleinschreibung bei der Asset-Auslieferung erkennen, in [PR](https://github.com/wailsapp/wails/pull/6395) von @taliesin-ai
- Inhaltstyp im Asset-Server korrekt erkennen, ohne Dateinamen im Speicher zu behalten, in [PR](https://github.com/wailsapp/wails/pull/6394) von @taliesin-ai
- Auf einen Standardwert zurückgreifen, wenn ein Desktop-Eintrag keinen Befehl enthält, in [PR](https://github.com/wailsapp/wails/pull/6392) von @taliesin-ai
- Verhindern, dass das Liquid-Glass-Beispiel auf anderen Plattformen als macOS vor dem Anwendungsstart ein Dialogfeld anzeigt, in [PR](https://github.com/wailsapp/wails/pull/6389) von @taliesin-ai
- Fenstergrößenbeschränkungen unter Windows beibehalten, wenn der Vollbildmodus mit einer beschränkten Größe gestartet wird, in [PR](https://github.com/wailsapp/wails/pull/6387) von @taliesin-ai
- Fehler behoben, durch den vor Run aufgerufene Dialogfelder dauerhaft blockierten, statt zurückzukehren, in [PR](https://github.com/wailsapp/wails/pull/6390) von @taliesin-ai

## Sicherheit
- Größen generierter Symbole vor der Bildspeicherzuweisung validieren, um übermäßigen Speicherverbrauch zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/6478) von @taliesin-ai
- Antworten der Release-Notes-API vor der Verwendung validieren, in [PR](https://github.com/wailsapp/wails/pull/6461) von @taliesin-ai
- Vollständige hexadezimale Statement-Handles im SQLite-Plugin voraussetzen, in [PR](https://github.com/wailsapp/wails/pull/6452) von @taliesin-ai
- Fallback-Treffer für Dateisystembibliotheken einschränken, um das Laden unbeabsichtigter Bibliotheken zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/6417) von @taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## Geändert
- Die Seite 404 mit den entsprechenden Wails-Dokumentationsseiten unter v2.wails.io verknüpfen, in [PR](https://github.com/wailsapp/wails/pull/6366) von @taliesin-ai v2
- Die Dokumentationswebsite v2 nach v2.wails.io verschieben, in [PR](https://github.com/wailsapp/wails/pull/6365) von @taliesin-ai
- Die Dokumentation v3 über wails.io und die Dokumentation v2 über v2.wails.io bereitstellen, in [PR](https://github.com/wailsapp/wails/pull/6364) von @taliesin-ai
- Release Notes für v3.0.0 aktualisieren und die Aussage in der Sicherheitswarnung korrigieren, in [PR](https://github.com/wailsapp/wails/pull/6363) von @taliesin-ai

## Behoben
- Callbacks von Go-Ereignis-Listenern seriell in Emit-Reihenfolge ausführen, in [PR](https://github.com/wailsapp/wails/pull/6275) von @DevLumuz
- Die Tray-Eigenschaft Menu immer bekannt geben, damit nach Run erstellte Tray-Symbole unter Linux nicht verworfen werden, in [PR](https://github.com/wailsapp/wails/pull/6284) von @AlbinoGeek
- Vor Run eingereichte Arbeiten für den Hauptthread bis zum Anwendungsstart zurückhalten, in [PR](https://github.com/wailsapp/wails/pull/6287) von @AlbinoGeek
- Beim Herunterladen von Release-Assets des GitHub-Updaters die korrekte Asset-API-URL verwenden, in [PR](https://github.com/wailsapp/wails/pull/6304) von @gavidroselj
- Bildschirmgrenzen und -positionen in macOS- und Linux-Konfigurationen mit unterschiedlichen DPI-Werten korrigieren, in [PR](https://github.com/wailsapp/wails/pull/6346) von @julianstorer
- Anfragetexte unter Windows beibehalten, wenn WebView2 den Content-Length-Header auslässt, in [PR](https://github.com/wailsapp/wails/pull/6326) von @strraaa
- Data Race beim Auslösen von Ereignis-Hooks beheben, indem unter einer Sperre ein Snapshot der Hooks erstellt wird, in [PR](https://github.com/wailsapp/wails/pull/6276) von @DevLumuz
- Fehler behoben, durch den `SystemTray.SetTooltip()` unter Linux wirkungslos blieb: Die Tooltip-Beschreibung des StatusNotifierItem wird jetzt gesetzt (#6017)
- Prüfung der Changelog-Nachweise im wöchentlichen Release-Workflow ausführen, in [PR](https://github.com/wailsapp/wails/pull/6374) von @taliesin-ai
- Prüfung der Release-Changelog-Nachweise korrigieren, damit Ausführungen Squash-Merge-Commits zugeordnet werden, in [PR](https://github.com/wailsapp/wails/pull/6373) von @taliesin-ai
- Fehler behoben, durch den die Linux-Dunkelmoduserkennung auf einem dunklen Desktop bis zur Änderung des Farbschemas den hellen Modus meldete und `linux:SystemThemeChanged` im GTK4-Backend nie ausgelöst wurde. Beide Backends verwenden jetzt dieselbe Implementierung, die auf dem Portal-Namespace `org.freedesktop.appearance` basiert (#6072)
- `Window.SetZoom` hebt Faktoren unter 1.0 bei Fenstern mit Rahmen (Windows, macOS, Linux) nicht mehr auf 1.0 an: WebView2 `ZoomFactor`, WKWebView `magnification` und die Zoomstufe von WebKitGTK akzeptieren jeden positiven Wert, und die Fensteroption `Zoom` leitete solche Werte bereits weiter. Rahmenlose Fenster behalten die Untergrenze 1.0 bei, da ihre Zieh- und Größenänderungsbereiche nicht an den Zoom angepasst werden ([#4590](https://github.com/wailsapp/wails/issues/4590)); nicht positive, NaN- und Inf-Faktoren werden jetzt ignoriert, statt auf 1.0 angehoben zu werden — von @APshenkin

## v3.0.0 - 2026-10-26

## Geändert
- Nächtliche Vorabversionen werden als `vX.Y.0-nightly.YYYYMMDD` veröffentlicht und geben eine Vorschau auf das nächste wöchentliche Release; wöchentliche Releases erhöhen die Minor-Version (`vX.Y.0`), sodass Gos `@latest` immer ein wöchentliches oder Hotfix-Release auflöst. Die Version des npm-Pakets `@wailsio/runtime` entspricht immer der Wails-Release-Version; Nightlies verwenden den Dist-Tag `nightly`
- Projektvorlagen fixieren `@wailsio/runtime` auf die Version der CLI, mit der sie generiert wurden
- Wails v3 erfordert Go 1.26 oder neuer (Go 1.26.10 oder 1.27.3 und neuer enthalten die Sicherheitskorrektur `crypto/x509`)

## Behoben
- Fehler behoben, durch den `wails3 update cli` keine v3-Releases fand: Es berücksichtigte nur v2-Tags und installierte aus dem falschen Pfad
- Behebt Abstürze von Windows-Apps mit `panic: SetTextColor failed` beim Zeichnen der Anwendungsmenüleiste in einigen Windows-Builds: Die Menüleiste zeichnet in einen gültigen Gerätekontext oder greift auf die Standarddarstellung zurück, `UAHMENUITEM.Position` entspricht jetzt dem C-`int`, und das Zeichnen des Menüs löst niemals eine Panic aus (#6311)
- Schützt das bereitgestellte Update, solange ein Neustart aussteht: Ein Download oder ein zweiter Neustart, der während des Herunterfahrens gestartet wird, kann die Nutzdaten, auf deren Installation der Austausch-Helfer wartet, weder löschen noch mit ihnen in Konflikt geraten, von @julianstorer
- Behebt Abstürze von Windows-Apps in `applyDarkMode`, wenn die Initialisierung von `pkg/w32` das Modul `uxtheme.dll` entladen hat, während dessen ordinale Funktionszeiger noch verwendet wurden: Das Modul bleibt jetzt während der gesamten Prozesslaufzeit geladen (#6356)
- Validiert die bereitgestellten Nutzdaten vor dem Austausch und die installierte Anwendung, bevor die Sicherung verworfen wird, sodass ein beschädigtes Update zurückgesetzt wird, statt eine funktionierende Anwendung zu ersetzen, von @julianstorer

## Entfernt
- Entfernt nach der abschließenden Übersetzungssynchronisierung in [PR](https://github.com/wailsapp/wails/pull/5736) die Crowdin-Integration, von @github-actions[bot]

## Sicherheit
- Blockiert reservierte `wails:`-Nachrichten aus ursprungsfremden untergeordneten Frames unter Linux (GTK3 und GTK4), indem ein Schutzskript in alle Frames injiziert wird; dadurch wird GHSA-2p5f-2gx8-mh8m entschärft
- Aktualisiert `github.com/klauspost/compress` auf 1.18.7 (GHSA-259r-337f-4rfw) und `golang.org/x/image` auf 0.45.0 und erstellt den Build mit Go 1.26.10, um die Korrektur für `crypto/x509` zu übernehmen (GO-2026-6809)
- Aktualisiert `vite` im Runtime-Paket, im Einrichtungsassistenten, in Beispielen und in Projektvorlagen auf 8.2.4 oder neuer (GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv)

### Bekannte Einschränkungen
- **Linux:** Der Schutz für untergeordnete Frames kann einen *gleichursprünglichen* untergeordneten Frame (beispielsweise einen `srcdoc`- oder `about:blank`-iframe) nicht daran hindern, über seinen übergeordneten Frame auf den nativen Nachrichtenhandler zuzugreifen, da WebKitGTK dem Handler keine Frame-Informationen bereitstellt. Ursprungsfremde Frames und Sandbox-Frames ohne `allow-same-origin` werden blockiert. Betten Sie keine nicht vertrauenswürdigen gleichursprünglichen Inhalte in ein Wails-Fenster ein.
- **Android (experimentell):** Die `window.wails`-JavaScript-Bridge ist für jeden Frame in der WebView verfügbar, ohne Prüfung des Hauptframes oder Ursprungs. Inhalte in einem eingebetteten iframe, einschließlich ursprungsfremder Inhalte, können gebundene Methoden aufrufen und reservierte `wails:`-Nachrichten senden. Betten Sie keine nicht vertrauenswürdigen Inhalte von Drittanbietern in die WebView der App ein. Dies unterscheidet sich von macOS und iOS, die `wails:`-Nachrichten aus untergeordneten Frames ignorieren. Eine Korrektur ist für 3.0.1 geplant.
- **Alle Plattformen:** Aufrufe gebundener Methoden aus untergeordneten Frames werden über HTTP an `/wails/runtime` und nicht an den nativen Nachrichtenhandler übertragen, weshalb der Schutz sie nicht einschränkt. Behandeln Sie jeden eingebetteten iframe so, als hätte er Zugriff auf Ihre gebundenen Go-Methoden.

## v3.0.0-beta.34 - 2026-10-20

## Behoben
- Behebt, dass die Windows-WebView bei der Aktivierung des Fensters den Fokus nicht wiedererlangt und Tab innerhalb der WebView festgehalten wird, statt zum nativen Fenster und zurück zu wechseln (#3783)
- Behebt, dass in den COM-vtable-Strukturen von `internal/webview2/pkg/webview2` (`ICoreWebView2_2`..`_27` und weiteren abgeleiteten Schnittstellen) die geerbte Methodenkette fehlte, wodurch Aufrufe im falschen vtable-Slot landeten (#6129)
- Behebt, dass Windows-Fragedialoge benutzerdefinierte Schaltflächenbeschriftungen ignorierten und Schaltflächen-Callbacks nie auslösten; benutzerdefinierte Schaltflächen werden jetzt mit einem TaskDialog angezeigt (mit Rückgriff auf MessageBox und positionsbasierter Callback-Zuordnung, wenn comctl32 v6 nicht verfügbar ist) (#5794)
- Behebt, dass `wails3 dev` bei 60s mit Vite hängen blieb: Die automatische Bereitschaftsprüfung des Frontends verwendet jetzt eine TCP-Sonde statt eines HTTP-GET, der keinen `Accept`-Header sendete und eine 404-Antwort erhielt ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Sorgt dafür, dass Systemleistenmenüs weiterhin funktionieren, wenn sie im geöffneten Zustand aktualisiert werden, in [PR](https://github.com/wailsapp/wails/pull/6301) von @taliesin-ai
- Behandelt vorübergehende WebView2-COM-Fehler bei Runtime-Operationen ordnungsgemäß, in [PR](https://github.com/wailsapp/wails/pull/6302) von @leaanthony

## v3.0.0-beta.33 - 2026-10-19

## Hinzugefügt
- Der Binding-Generator erzeugt jetzt ein `bindings/events`-Modul, das für jedes auffindbare Ereignis, das mit einem zur Kompilierzeit konstanten Namen über `application.RegisterEvent` registriert wurde, eine typisierte Namenskonstante exportiert (z. B. `UserLogin` für `user:login`); diese kann ohne das Vite-Plug-in mit `Events.On`/`Emit` verwendet werden (#5645)

## Behoben
- Behebt, dass Windows-Fragedialoge benutzerdefinierte Schaltflächenbeschriftungen ignorierten und Schaltflächen-Callbacks nie auslösten; benutzerdefinierte Schaltflächen werden jetzt mit einem TaskDialog angezeigt (mit Rückgriff auf MessageBox und positionsbasierter Callback-Zuordnung, wenn comctl32 v6 nicht verfügbar ist) (#5794)
- Behebt, dass `wails3 dev` bei 60s mit Vite hängen blieb: Die automatische Bereitschaftsprüfung des Frontends verwendet jetzt eine TCP-Sonde statt eines HTTP-GET, der keinen `Accept`-Header sendete und eine 404-Antwort erhielt ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Sorgt dafür, dass Systemleistenmenüs weiterhin funktionieren, wenn sie im geöffneten Zustand aktualisiert werden, in [PR](https://github.com/wailsapp/wails/pull/6301) von @taliesin-ai
- Behandelt vorübergehende WebView2-COM-Fehler bei Runtime-Operationen ordnungsgemäß, in [PR](https://github.com/wailsapp/wails/pull/6302) von @leaanthony

## v3.0.0-beta.32 - 2026-10-18

## Hinzugefügt
- Fügt native macOS-Menüabschnittsüberschriften für macOS 14+ hinzu – siehe die [Dokumentationswebsite](https://v3.wails.io) (20 Seiten aktualisiert), in [PR](https://github.com/wailsapp/wails/pull/6297) von @leaanthony

## Behoben
- Sorgt dafür, dass SetRelativePosition unter macOS relativ zum Arbeitsbereich des Bildschirms arbeitet, in [PR](https://github.com/wailsapp/wails/pull/6299) von @leaanthony
- Der Updater-Helfer startet nach einer Zeitüberschreitung des übergeordneten Prozesses oder einem fehlgeschlagenen Austausch das unveränderte Ziel erneut, in [PR](https://github.com/wailsapp/wails/pull/6298) von @taliesin-ai
- Fixiert pnpm für versionierte Dokumentation auf Version 10, in [PR](https://github.com/wailsapp/wails/pull/6293) von @taliesin-ai

## v3.0.0-beta.31 - 2026-10-18

## Hinzugefügt
- Stellt den Subpfadexport `@wailsio/runtime/cancellable` bereit, sodass `CancellablePromise` importiert werden kann, ohne den Root-Einstiegspunkt der Runtime auszuwerten (#6153)

## Geändert
- Erfordert Go 1.26 oder neuer: Das Modul v3 deklariert jetzt `go 1.26.0`, und von `wails3 init` und `wails3 service init` erzeugte Projekte deklarieren dieselbe Mindestversion

## Behoben
- Behebt, dass sich die Ausgabe des Binding-Generators beim Erstellen mit Go 1.27 änderte, wobei `encoding/json`-Typen wie `Marshaler`, `Unmarshaler` und `RawMessage` zu Aliasen von `encoding/json/v2`- und `encoding/json/jsontext`-Typen wurden: Sie werden wieder als `any` dargestellt, und für `encoding/json/v2` werden keine Modelle erzeugt (#6292)
- Validiert an `wails3 init -t` übergebene URLs entfernter Vorlagen vor dem Klonen: Akzeptiert werden nur `https://`-, `http://`-, `ssh://`-, `git://`- und `user@host:path`-URLs, Werte werden nach `--` an git übergeben, und git wird während des Klonens auf diese Transportprotokolle beschränkt
- Lehnt außerhalb des gültigen Bereichs liegende Werte des `x-wails-window-id`-Headers im HTTP-Runtime-Transport ab, statt sie unbemerkt auf die ID eines anderen Fensters zu kürzen
- Härtet die Archivextraktion des Updaters gegen Zip-Slip ab: Eintragspfade werden jetzt mit `filepath.IsLocal` sowie durch eine Prüfung geprüft, ob der zusammengesetzte Pfad innerhalb des Zielverzeichnisses liegt; Tests decken verschachtelte `..`-Einträge, Einträge mit Rückwärtsschrägstrichen und absoluten Pfaden sowie durch symbolische Links hindurch schreibende Einträge in zip- und tar.gz-Archiven ab
- Absturz von AppImages aus `wails3 generate appimage` beim Start außerhalb von Debian und Ubuntu mit „Unable to spawn a new child process: … WebKitNetworkProcess“ behoben. Das gebündelte WebKit führt nun die im AppImage enthaltenen Hilfsprozesse aus. Die bubblewrap-Sandbox von WebKit bleibt aktiviert
- Fehler behoben, durch den `-tags server`-Builds für Android und iOS nicht kompiliert werden konnten; damit ist [#6271](https://github.com/wailsapp/wails/issues/6271) behoben
- Fehler behoben, durch den der Updater unter Linux den Austausch nicht durchführen konnte, wenn sich das temporäre Verzeichnis auf einem anderen Dateisystem befand (etwa einem tmpfs unter `/tmp`). Dazu werden die Dateien neben der installierten Binärdatei bereitgestellt und dateisystemübergreifend kopiert
- Fehler behoben, durch den der Updater versuchte, den schreibgeschützten Einhängepunkt eines AppImage statt der Datei `.AppImage` zu ersetzen
- Fehler behoben, durch den reservierte `wails:`-Nachrichten zur Fenstersteuerung unter macOS und iOS aus iframes akzeptiert wurden. Sie werden jetzt nur noch vom Hauptframe des Fensters akzeptiert; andere Nachrichten aus iframes erreichen weiterhin `RawMessageHandler` einschließlich ihrer Herkunftsinformationen
- Senden von `WindowDidUpdate` und `WindowWillUpdate` unter macOS eingestellt: `NSWindowDidUpdateNotification` wird bei jedem Durchlauf der Run Loop ausgelöst, sodass die Seite pro Sekunde etwa 120 leere `evaluateJavaScript:`-Aufrufe empfing. Unter macOS verbrauchte 27.0 dabei jeweils die vorübergehende Aktivierung durch den Klick des Benutzers. Dadurch wurden `navigator.clipboard.writeText()` und andere aktivierungsgebundene APIs bei echten Klicks mit `NotAllowedError` abgelehnt ([#6170](https://github.com/wailsapp/wails/issues/6170)), von @APshenkin
- Standardsprache von Windows-info.json auf en-US (0409) korrigiert, damit der Explorer Versionsdetails anzeigt (#6210)
- Fehler behoben, durch den der Linux-DMA-BUF-Renderer auf Laptops mit einer NVIDIA-dGPU neben einer AMD- oder Intel-iGPU deaktiviert wurde: Unter Wayland wird er jetzt nur deaktiviert, wenn WebKit auf NVIDIA rendert (weil die App auf die dGPU ausgelagert wurde oder der Compositor darauf läuft), und nicht bereits, wenn das nvidia-Modul geladen ist, was das Rendering auf der iGPU extrem verlangsamte. X11-Sitzungen behalten das bisherige Verhalten bei, da die Rendering-GPU von Xorg durch dessen Konfiguration festgelegt wird
- Fehler beim dateisystemübergreifenden Umbenennen im Linux-Updater durch eine Ausweichlösung mit bereitgestellter Kopie in [PR](https://github.com/wailsapp/wails/pull/6135) behandelt, von @john-okeefe
- Fehler behoben, durch den `WebviewWindow.SetMenu` unter Windows die Menüleiste nur neu zeichnete, wenn ein Menüleisten-Theme konfiguriert war; damit ist [#6103](https://github.com/wailsapp/wails/issues/6103) behoben
- Fehler behoben, durch den `WebviewWindow.SetMenu(nil)` unter Windows eine Panic auslöste, statt wie unter macOS keine Wirkung zu haben; damit ist [#6104](https://github.com/wailsapp/wails/issues/6104) behoben
- Changelog-Einträge werden bei konkurrierenden Pushes durch Wiederholungsversuche in [PR](https://github.com/wailsapp/wails/pull/6094) beibehalten, von @leaanthony
- Fehler behoben, durch den `go mod vendor` auf jeder Plattform mit `pattern arm64/WebView2Loader.dll: no matching files found` fehlschlug. Dazu wurden Einbettungen entfernt, die auf nie mit dem Modul ausgelieferte Binärdateien verwiesen. Damit sind [#5782](https://github.com/wailsapp/wails/issues/5782) und [#5376](https://github.com/wailsapp/wails/issues/5376) behoben, in [PR](https://github.com/wailsapp/wails/pull/6031) von @Grantmartin2002
- Doppelte Fensterbereinigung unter macOS in [PR](https://github.com/wailsapp/wails/pull/6159) verhindert, von @charlesbalderson
- WM_CONTEXTMENU für das Taskleistenmenü behandelt, damit es sich per Rechtsklick öffnet, in [PR](https://github.com/wailsapp/wails/pull/6162) von @wzk2025
- Fehler in allen fünf `InvokeSync*`-Varianten behoben, durch den sie ihren Aufrufer dauerhaft blockierten, wenn der weitergeleitete Callback eine Panic auslöste. Damit ist [#6107](https://github.com/wailsapp/wails/issues/6107) behoben, in [PR](https://github.com/wailsapp/wails/pull/6108) von @RALIST
- Zerstörte GTK3-Untermenüs werden in [PR](https://github.com/wailsapp/wails/pull/6236) sicher neu erstellt, von @leaanthony
- Windows-Weiterleitungs-Shim für verschachtelte Worker in [PR](https://github.com/wailsapp/wails/pull/6261) implementiert, von @taliesin-ai
- Fehler behoben, durch den `OpenFile()` unter Windows `PromptForMultipleSelection()` ignorierte, wenn `CanChooseDirectories(true)`: Benutzer können jetzt wie unter macOS und Linux mehrere Ordner gleichzeitig auswählen (#6263)
- Fehler bei als Wert übergebenen `double`-Argumenten für WebView2-COM-Methoden unter windows/arm64 behoben. Betroffen waren dort `ZoomFactor`, `RasterizationScale`, Cookie-`Expires` sowie Skalierung, Seitengröße und Ränder der Druckeinstellungen (die Setter der Druckeinstellungen gaben zuvor `ErrDoubleArgUnsupported` zurück). Betroffene Aufrufe durchlaufen jetzt ein `runtime.cgocall`-Trampolin, das `d0` lädt, wozu der Syscall-Pfad von Go nicht in der Lage ist (golang.org/issue/62583), von @APshenkin

## Sicherheit
- Der Servermodus (`-tags server`) akzeptiert `/wails/runtime`-Aufrufe jetzt nur als `POST`-Anfragen vom eigenen Ursprung des Servers oder von einem in `ServerOptions.WebSocketOriginPatterns` aufgeführten Ursprung. Andere Methoden erhalten HTTP 405 und andere Browserursprünge HTTP 403. Bei Bereitstellungen hinter einem Reverse Proxy, der den `Host`-Header umschreibt, muss der öffentliche Ursprung zu `WebSocketOriginPatterns` hinzugefügt werden
- Der an eine Loopback-Adresse gebundene Servermodus (`-tags server`; standardmäßig `localhost`, `127.0.0.0/8` oder `::1`) lehnt Anfragen jetzt mit HTTP 403 ab, sofern der `Host`-Header nicht `localhost`, ein IP-Literal, der Bind-Host oder ein mit `ServerOptions.WebSocketOriginPatterns` übereinstimmender Host ist. Dadurch werden DNS-Rebinding-Angriffe auf `/wails/runtime` und die WebSocket-Endpunkte blockiert
- Windows-Benachrichtigungen greifen bei einem Fehlschlag der WinRT-Toast-API nicht mehr auf PowerShell zurück; stattdessen wird der Fehler zurückgegeben

## v3.0.0-beta.30 - 2026-10-15

## Behoben
- Fehler behoben, durch den der Updater versuchte, den schreibgeschützten Einhängepunkt eines AppImage statt der Datei `.AppImage` zu ersetzen
- Senden von `WindowDidUpdate` und `WindowWillUpdate` unter macOS eingestellt: `NSWindowDidUpdateNotification` wird bei jedem Durchlauf der Run Loop ausgelöst, sodass die Seite pro Sekunde etwa 120 leere `evaluateJavaScript:`-Aufrufe empfing. Unter macOS verbrauchte 27.0 dabei jeweils die vorübergehende Aktivierung durch den Klick des Benutzers. Dadurch wurden `navigator.clipboard.writeText()` und andere aktivierungsgebundene APIs bei echten Klicks mit `NotAllowedError` abgelehnt ([#6170](https://github.com/wailsapp/wails/issues/6170)), von @APshenkin
- Standardsprache von Windows-info.json auf en-US (0409) korrigiert, damit der Explorer Versionsdetails anzeigt (#6210)
- Fehler beim dateisystemübergreifenden Umbenennen im Linux-Updater durch eine Ausweichlösung mit bereitgestellter Kopie in [PR](https://github.com/wailsapp/wails/pull/6135) behandelt, von @john-okeefe
- Fehler behoben, durch den `WebviewWindow.SetMenu` unter Windows die Menüleiste nur neu zeichnete, wenn ein Menüleisten-Theme konfiguriert war; damit ist [#6103](https://github.com/wailsapp/wails/issues/6103) behoben
- Fehler behoben, durch den `WebviewWindow.SetMenu(nil)` unter Windows eine Panic auslöste, statt wie unter macOS keine Wirkung zu haben; damit ist [#6104](https://github.com/wailsapp/wails/issues/6104) behoben
- Changelog-Einträge werden bei konkurrierenden Pushes durch Wiederholungsversuche in [PR](https://github.com/wailsapp/wails/pull/6094) beibehalten, von @leaanthony
- Fehler behoben, durch den `go mod vendor` auf jeder Plattform mit `pattern arm64/WebView2Loader.dll: no matching files found` fehlschlug. Dazu wurden Einbettungen entfernt, die auf nie mit dem Modul ausgelieferte Binärdateien verwiesen. Damit sind [#5782](https://github.com/wailsapp/wails/issues/5782) und [#5376](https://github.com/wailsapp/wails/issues/5376) behoben, in [PR](https://github.com/wailsapp/wails/pull/6031) von @Grantmartin2002
- Doppelte Fensterbereinigung unter macOS in [PR](https://github.com/wailsapp/wails/pull/6159) verhindert, von @charlesbalderson
- WM_CONTEXTMENU für das Taskleistenmenü behandelt, damit es sich per Rechtsklick öffnet, in [PR](https://github.com/wailsapp/wails/pull/6162) von @wzk2025
- Fehler in allen fünf `InvokeSync*`-Varianten behoben, durch den sie ihren Aufrufer dauerhaft blockierten, wenn der weitergeleitete Callback eine Panic auslöste. Damit ist [#6107](https://github.com/wailsapp/wails/issues/6107) behoben, in [PR](https://github.com/wailsapp/wails/pull/6108) von @RALIST

## v3.0.0-beta.29 - 2026-10-13

## Hinzugefügt
- Globale Tastenkürzel unter macOS berücksichtigen jetzt das Tastaturlayout: Buchstaben-Tastenkürzel werden auf Nicht-QWERTY-Layouts (AZERTY, QWERTZ, Dvorak, ...) über `UCKeyTranslate` an die beschriftete Taste statt an die feste physische QWERTY-Position gebunden.

## Behoben
- Natives Menü unter Linux vor dem Festlegen initialisiert, in [PR](https://github.com/wailsapp/wails/pull/6233) von @abhicodes-007
- Regressionen der GA-Bereitschaft unter macOS bei der Ausgabeverarbeitung und im Einrichtungsassistenten behoben, in [PR](https://github.com/wailsapp/wails/pull/6237) von @taliesin-ai
- Fehler behoben, durch den die WebView2-COM-Wrapper `AddRef`/`Release` unter Windows einen veralteten Thread-Last-Error als Fehlschlag meldeten. Sie geben jetzt den Referenzzähler zurück, sodass eine Freigabe weder die Anwendung abbrechen noch eine Asset-Anfrage fehlschlagen lassen kann (#3931)
- Behebt, dass `ICoreWebView2Deferral.Release` den `AddRef`-Slot aufruft, wodurch unter Windows die Zurückstellung jeder verzögerten Asset-Anfrage nicht aufgehoben wurde (#3931)

## v3.0.0-beta.28 - 2026-10-05

## Behoben
- Hält die macOS-App-Ausgabe außerhalb von Dateianbietern in [PR](https://github.com/wailsapp/wails/pull/6225) von @taliesin-ai
- Vervollständigt den Einrichtungsassistenten und zeigt Systemdetails an in [PR](https://github.com/wailsapp/wails/pull/6224) von @taliesin-ai
- Verbessert unter Windows die Fehlerbehandlung und Ressourcenbereinigung bei der WebView2-Wiederherstellung in [PR](https://github.com/wailsapp/wails/pull/6222) von @leaanthony
- Verbessert unter Windows die WebView2-Wiederherstellung, um App-Abstürze und verwaiste Controller zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/6220) von @taliesin-ai
- Korrigiert die Anleitung zur Medienautorisierung unter macOS in [PR](https://github.com/wailsapp/wails/pull/6221) von @taliesin-ai
- Behebt, dass ein Rechtsklick auf ein Linux-Symbol im Infobereich neben dem Öffnen des Menüs auch den Klick-Handler ausführte, sodass ein verknüpftes Fenster bei jedem Rechtsklick umgeschaltet wurde (#6018)
- Hält WebView2-Callback-Handler unter Windows bis zur nativen Freigabe aktiv in [PR](https://github.com/wailsapp/wails/pull/6184) von @taliesin-ai

## v3.0.0-beta.27 - 2026-10-01

## Behoben
- Behebt, dass die Option `Permissions` unter macOS 12 und neuer ignoriert wurde: Die Medienaufnahmemethode `WKUIDelegate` ist jetzt implementiert, sodass `PermissionAllow` und `PermissionDeny` für Kamera- und Mikrofonanfragen ebenso gelten wie unter Linux und Windows. `NSCameraUsageDescription` / `NSMicrophoneUsageDescription` sowie bei Verwendung einer Sandbox die entsprechenden Geräteberechtigungen sind weiterhin erforderlich

## v3.0.0-beta.26 - 2026-09-25

## Hinzugefügt
- Der Watcher führt HTTP-Bereitschaftsprüfungen für den Frontend-Entwicklungsserver durch in [PR](https://github.com/wailsapp/wails/pull/6150) von @atterpac

## Behoben
- Stellt den Betrieb nach WebView2-Prozessfehlern wieder her, statt ein leeres Fenster zurückzulassen, in [PR](https://github.com/wailsapp/wails/pull/6002) von @taliesin-ai

## v3.0.0-beta.25 - 2026-09-22

## Behoben
- Beendet `wails3 dev` und dessen Hintergrundprozesse, wenn die primäre Anwendung beendet wird (#6048)

## v3.0.0-beta.24 - 2026-09-20

## Sicherheit
- Beschränkt die GITHUB_TOKEN-Berechtigungen in Build- und Veröffentlichungsworkflows gemäß dem Prinzip der geringsten Rechte
- Behebt Pfadtraversierung und das Ausbrechen über symbolische Links in der Asset-Middleware des Bildschirmbeispiels
- Hält die Installation von Abhängigkeiten durch den Einrichtungsassistenten funktionsfähig und beschränkt sie zugleich auf genehmigte Paketmanagerbefehle

## v3.0.0-beta.23 - 2026-09-16

## Hinzugefügt

- Ermöglicht das Überschreiben der GTK-Anwendungs-ID unter Linux in [PR](https://github.com/wailsapp/wails/pull/5972) von @overlordtm

## Behoben

- Die Linux-Einzelinstanz beansprucht den Busnamen unter der eigenen ID der App in [PR](https://github.com/wailsapp/wails/pull/5973) von @overlordtm
- Bricht abgebrochene Anfragen für benutzerdefinierte Schemas unter Linux ab in [PR](https://github.com/wailsapp/wails/pull/6126) von @leaanthony

## v3.0.0-beta.22 - 2026-09-14

## Hinzugefügt

- Vervollständigt die Übersetzungen und portiert Community-Medien – siehe die [Dokumentationswebsite](https://v3.wails.io) (1812 Seiten aktualisiert) – in [PR](https://github.com/wailsapp/wails/pull/6122) von @leaanthony
- Dokumentiert einbettbare Wails-Schaltflächen in verschiedenen Stilen und Größen in [PR](https://github.com/wailsapp/wails/pull/6120) von @leaanthony

## v3.0.0-beta.21 - 2026-09-13

## Hinzugefügt

- Stellt die Wails-v3-Dokumentation mit M-Press bereit in [PR](https://github.com/wailsapp/wails/pull/6116) von @leaanthony

## Behoben

- Parst JSON-Slug-Werte aus dem MPD-Frontmatter für die Changelog-Generierung in [PR](https://github.com/wailsapp/wails/pull/6118) von @leaanthony
- Der Updater entfernt Hilfsumgebungsvariablen und startet nach fehlgeschlagenen Sicherungen das ursprüngliche Ziel neu in [PR](https://github.com/wailsapp/wails/pull/6080) von @cnmax
- Startet den standardmäßigen Signal-Handler während App.Run in [PR](https://github.com/wailsapp/wails/pull/6098) von @leaanthony
- Das Windows-Menü verarbeitet nil-Menüs, gibt ersetzte Ressourcen frei und zeichnet die Menüleiste neu in [PR](https://github.com/wailsapp/wails/pull/6112) von @taliesin-ai
- Stellt die MSIX-Paketierung für neue Projekte mithilfe einer gemeinsamen YAML-Konfiguration wieder her in [PR](https://github.com/wailsapp/wails/pull/6115) von @leaanthony
- Behebt, dass generierte JavaScript- und TypeScript-Bindings nicht geladen werden konnten, wenn Erzeuger generischer Modelle auf später deklarierte Hilfsfunktionen verwiesen, und verhindert Stacküberläufe beim Erstellen gegenseitig abhängiger generischer Modelle (#6062)

## v3.0.0-beta.20 - 2026-09-10

## Geändert

- Aktualisiert die Clave-Showcase-Links auf die aktuelle Website und das aktuelle Repository in [PR](https://github.com/wailsapp/wails/pull/6082) von @01xR4in

## Behoben

- Bricht abgebrochene Windows-Asset-Anfragen einschließlich Worker-Anfragen ab und erhält dabei Keepalive-Handler über Navigationen hinweg. Leitet auf Apple-Plattformen native Anfragekontexte durch den Anwendungs-Wrapper weiter. (#5963, #5969)
- Bewahrt Changelog-Einträge bei konkurrierenden Pushes durch Wiederholungsversuche in [PR](https://github.com/wailsapp/wails/pull/6094) von @leaanthony
- Behebt, dass `go mod vendor` auf jeder Plattform mit `pattern arm64/WebView2Loader.dll: no matching files found` fehlschlug, indem die Einbettungen entfernt wurden, die auf nie mit dem Modul ausgelieferte Binärdateien verwiesen. Dadurch werden [#5782](https://github.com/wailsapp/wails/issues/5782) und [#5376](https://github.com/wailsapp/wails/issues/5376) behoben, in [PR](https://github.com/wailsapp/wails/pull/6031) von @Grantmartin2002

## Entfernt

- Entfernt die Unterstützung für den nativen WebView2-Loader, der durch den reinen Go-Loader ersetzt wurde. Dadurch entfallen die eingebetteten `WebView2Loader.dll`-Binärdateien und die Abhängigkeit `github.com/jchv/go-winloader`. Das Build-Tag `native_webview2loader` wird weiterhin akzeptiert und verursacht keinen Fehler mehr, hat jedoch keine Auswirkung auf v3-Builds, in [PR](https://github.com/wailsapp/wails/pull/6031) von @Grantmartin2002
- Entfernt ungenutzte Build-Tags und die FPS-Option aus dem macOS-API-Leitfaden in [PR](https://github.com/wailsapp/wails/pull/6097) von @leaanthony

## v3.0.0-beta.19 - 2026-09-09

## Hinzugefügt

- Schützt private macOS-APIs mit Build-Tags für die explizite Aktivierung – siehe [Dokumentation](https://v3.wails.io/features/browser/integration) und [Dokumentation](https://v3.wails.io/features/environment/info) und [Dokumentation](https://v3.wails.io/features/windows/basics) und [Dokumentation](https://v3.wails.io/features/windows/frameless) und [Dokumentation](https://v3.wails.io/features/windows/notch-windows) und [Dokumentation](https://v3.wails.io/features/windows/options) und [Dokumentation](https://v3.wails.io/guides/build/macos) und [Dokumentation](https://v3.wails.io/guides/build/private-macos-apis) und [Dokumentation](https://v3.wails.io/reference/overview) in [PR](https://github.com/wailsapp/wails/pull/6087) von @leaanthony

## Behoben

- Lehnt Laufzeitanfragen über 64 MiB mit HTTP 413 ab in [PR](https://github.com/wailsapp/wails/pull/6091) von @leaanthony

## Sicherheit

- Sichert MCP-Ursprünge und Remotezugriffe durch Token-Authentifizierung ab in [PR](https://github.com/wailsapp/wails/pull/6092) von @leaanthony

## v3.0.0-beta.18 - 2026-09-08

## Behoben

- Malloc-Speicherleck unter Linux und Darwin durch die Verwendung von Pointer-Receivern in [PR](https://github.com/wailsapp/wails/pull/6083) von @4RH1T3CT0R7 behoben

## v3.0.0-beta.17 - 2026-09-06

## Behoben

- Windows: Ein fehlgeschlagenes oder nil-`GetRequest` im WebResourceRequested-Handler beendet den Prozess nicht mehr (`log.Fatal` / Panic durch nil-Dereferenzierung). Stattdessen wird die Anfrage verworfen und protokolliert, in [PR](https://github.com/wailsapp/wails/pull/6006) von @midagedev

## v3.0.0-beta.16 - 2026-08-29

## Geändert

- Notarisierungskennwort in einem neuen Terminalfenster abfragen, in [PR](https://github.com/wailsapp/wails/pull/6029) von @leaanthony

## Behoben

- Klicktypen im Infobereich unter macOS korrekt verarbeiten, in [PR](https://github.com/wailsapp/wails/pull/5919) von @ChewbaccaCookie
- CI entfernt vor der Aktualisierung nicht verwendete Microsoft-apt-Repositorys, in [PR](https://github.com/wailsapp/wails/pull/6041) von @Grantmartin2002

## v3.0.0-beta.15 - 2026-08-27

## Behoben

- Zeitüberschreitung für das Einbetten von WebView2 auf 60 Sekunden erhöht, in [PR](https://github.com/wailsapp/wails/pull/6043) von @Grantmartin2002

## v3.0.0-beta.14 - 2026-08-26

## Behoben

- Tastendrücke mit Control und Buchstaben unter macOS korrekt benannt, in [PR](https://github.com/wailsapp/wails/pull/6032) von @taliesin-ai
- ICO-Symbole im Infobereich korrigiert und unter Windows an das Taskleisten-Theme angepasst, in [PR](https://github.com/wailsapp/wails/pull/6016) von @nik9play

## v3.0.0-beta.13 - 2026-08-25

## Behoben

- Aufgaben des Hauptthreads unter macOS auch während einer modalen Schleife weiter ausführen, in [PR](https://github.com/wailsapp/wails/pull/6026) von @leaanthony
- Fehler für den sicheren Speicher auf Mobilgeräten ermöglichen und bei Fehlern den Zugriff verweigern, in [PR](https://github.com/wailsapp/wails/pull/5923) von @mortenolsrud
- Ereignis-Hooks der Anwendung auch ausführen, wenn kein Listener registriert ist, in [PR](https://github.com/wailsapp/wails/pull/5999) von @archy-rock3t-cloud
- Tippfehler in Kommentaren und lokalisierter Dokumentation korrigiert, in [PR](https://github.com/wailsapp/wails/pull/6023) von @haoku123
- Die unter `v3/examples` eingecheckten vorkompilierten macOS-Binärdateien entfernt, in [PR](https://github.com/wailsapp/wails/pull/6025) von @4RH1T3CT0R7

## v3.0.0-beta.12 - 2026-08-21

## Hinzugefügt

- Benachrichtigungsfenster für die macOS-Notch mit Lebenszyklus- und Telemetriebeispiel hinzugefügt – siehe [Dokumentation](https://v3.wails.io/features/windows/notch-windows), in [PR](https://github.com/wailsapp/wails/pull/6010) von @leaanthony
- Unterstützung für macOS-NSPanel-Fenster mit neuen Optionen und nativer Integration hinzugefügt – siehe [Dokumentation](https://v3.wails.io/features/windows/options), in [PR](https://github.com/wailsapp/wails/pull/6008) von @leaanthony

## Behoben

- Verhindert, dass SQLite Prepare bei gleichzeitigen Aufrufen hängen bleibt, in [PR](https://github.com/wailsapp/wails/pull/5998) von @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Entfernt

- Veralteten Implementierungs-Tracker aus der Dokumentation entfernt, in [PR](https://github.com/wailsapp/wails/pull/6005) von @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Behoben

- Behoben, dass der GTK4-Linux-Host Startargumente für benutzerdefinierte Protokolle und Dateizuordnungen verwirft, in [PR](https://github.com/wailsapp/wails/pull/6000) von @midagedev
- Gelöschte Zeilen und Korrekturen aus derselben Quelle bei der Changelog-Validierung korrekt verarbeitet, in [PR](https://github.com/wailsapp/wails/pull/5993) von @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Hinzugefügt

- Sicheren wails3-mcp-Server für die agentengestützte Projektverwaltung hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5896) von @leaanthony
- Dokumentation für Modelle in Bindings hinzugefügt – siehe [Dokumentation](https://v3.wails.io/features/bindings/models), in [PR](https://github.com/wailsapp/wails/pull/5988) von @taliesin-ai
- rpm-ostree-Installation für atomare Linux-Systeme unterstützt, in [PR](https://github.com/wailsapp/wails/pull/5987) von @leaanthony
- Native wöchentliche Erzeugung und Veröffentlichung von Diagrammen zum Sternverlauf hinzugefügt – siehe [Dokumentation](https://v3.wails.io/credits) und [Dokumentation](https://v3.wails.io/de/credits) und [Dokumentation](https://v3.wails.io/fr/credits) und [Dokumentation](https://v3.wails.io/id/credits) und [Dokumentation](https://v3.wails.io/ja/credits) und [Dokumentation](https://v3.wails.io/ko/credits) und [Dokumentation](https://v3.wails.io/pt/credits) und [Dokumentation](https://v3.wails.io/ru/credits) und [Dokumentation](https://v3.wails.io/zh-cn/credits) und [Dokumentation](https://v3.wails.io/zh-tw/credits), in [PR](https://github.com/wailsapp/wails/pull/5986) von @leaanthony
- Nur für Darwin verfügbares mac-Paket zum Auflösen von Ressourcen in Anwendungspaketen hinzugefügt – siehe [Dokumentation](https://v3.wails.io/guides/build/macos), in [PR](https://github.com/wailsapp/wails/pull/5965) von @leaanthony
- Showcase-Seite und Indexeintrag für Condui hinzugefügt – siehe [Dokumentation](https://v3.wails.io/community/showcase/condui) und [Dokumentation](https://v3.wails.io/community/showcase), in [PR](https://github.com/wailsapp/wails/pull/5962) von @mgueregath
- Showcase-Seite für Redis Viewer mit Screenshots und Projektlink hinzugefügt – siehe [Dokumentation](https://v3.wails.io/community/showcase) und [Dokumentation](https://v3.wails.io/community/showcase/redisviewer), in [PR](https://github.com/wailsapp/wails/pull/5984) von @redisviewer

## Geändert

- GTK-Anwendungsflags für Linux auf G<em>APPLICATION</em>NON_UNIQUE aktualisiert, in [PR](https://github.com/wailsapp/wails/pull/5971) von @overlordtm
- Fehlende Fensterereignisse auf Debug- statt auf Warnstufe protokollieren, in [PR](https://github.com/wailsapp/wails/pull/5914) von @julianstorer

## Behoben

- Registrierten macOS-Tastenkombinationen Vorrang vor der Webview gewährt, in [PR](https://github.com/wailsapp/wails/pull/5902) von @julianstorer
- Defekte Links in der Seitenleiste der Dokumentation repariert, in [PR](https://github.com/wailsapp/wails/pull/5937) von @northes
- Kontexte von Ressourcenanfragen unter macOS und iOS abbrechen, wenn WebKit die zugehörige Aufgabe für das benutzerdefinierte Schema abbricht (#5963)
- WindowSetFullscreenButtonEnabled-Nachricht verarbeitet, in [PR](https://github.com/wailsapp/wails/pull/5976) von @archy-rock3t-cloud
- Fragment in die preact-ts-Vorlage importiert, um den Build-Fehler zu beheben, in [PR](https://github.com/wailsapp/wails/pull/5979) von @haoku123
- Absturz älterer reiner GTK3-Dienstanwendungen verhindert, wenn die Bildschirmerkennung ausgeführt wird, bevor ein aktives Fenster oder Display verfügbar ist (#5966)
- Release-Läufe mit expliziter Version können fortfahren, wenn der Changelog für unveröffentlichte Änderungen leer ist (#5977)

## Sicherheit

- nanoid-Lockfiles der Website auf die gepatchte Version 3.3.18 aktualisiert, um Sicherheitshinweise zu beheben, in [PR](https://github.com/wailsapp/wails/pull/5985) von @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Hinzugefügt

- Generierung von Dokumentations-URLs zu automatischen Changelog-Einträgen hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5957) von @taliesin-ai
- Streams hinzufügen: bidirektionale Byte-Streams zwischen Go und JavaScript mit dem WebSocket-Programmiermodell, jedoch ohne lauschenden Socket. Ein Stream wird in Go mit `app.HandleStream(name, handler)` deklariert und vom Frontend mit `Stream(name)` verbunden; dies gibt ein Objekt in der Form von `WebSocket` zurück. Go→JS wird über den Asset-Server durch eine gehaltene Abfrage pro Fenster übertragen, JS→Go durch einen normalen POST; dabei wird weder ein TCP-Port gebunden noch etwas über `evaluateJavaScript` geleitet. In Server-Builds (`-tags server`) wird derselbe Handler stattdessen über einen echten WebSocket bereitgestellt, sodass der Anwendungscode in allen Builds identisch ist. Von @leaanthony
- Mailbox-Changelog-Eintrag in [PR](https://github.com/wailsapp/wails/pull/5935) nach „Unveröffentlicht“ verschieben, von @leaanthony

## Geändert

- Automatische Generierung der Dokumentations-Seitenleiste und Ableitung des Blog-Autorentyps in [PR](https://github.com/wailsapp/wails/pull/5938) aktualisieren, von @leaanthony

## Behoben

- Die WebView2-Initialisierung verwendet in [PR](https://github.com/wailsapp/wails/pull/5952) eine Frist und eine Nachrichtenschleife, von @leaanthony
- Der WebView2-Cookie-Test wird in CI übersprungen, sofern er nicht explizit aktiviert wurde, und seine Ausführung wird in [PR](https://github.com/wailsapp/wails/pull/5951) an den aktuellen Betriebssystem-Thread gebunden, von @leaanthony
- Windows-Menü-Builder stellen in [PR](https://github.com/wailsapp/wails/pull/5944) Befehls-IDs für übergeordnete Untermenüeinträge wieder her, von @gilad-ch
- Offizielles Cross-Compilation-Image an die Mindestversion für GTK-4.14+- und Linux-Unterstützung angleichen (#5928)
- Das iOS-Xcode-Projekt so konfigurieren, dass geerbte Linker-Flags erhalten bleiben, und in [PR](https://github.com/wailsapp/wails/pull/5915) -ObjC hinzufügen, von @mortenolsrud
- Übermäßigen Wechsel von TCP-Verbindungen im `wails3 dev`-Asset-Proxy bei großen Frontends beheben, der die kurzlebigen Ports des Hosts erschöpfen und dazu führen konnte, dass unabhängige Prozesse mit `EADDRNOTAVAIL` fehlschlugen
- Fensterspezifisches Ereignis-JavaScript für geordnete Auslieferung und Rückstaukontrolle in [PR](https://github.com/wailsapp/wails/pull/5934) in eine Warteschlange stellen, von @leaanthony

## Entfernt

- Release-Pipeline für Desktop-Binärdateien entfernen: v3-Releases bestehen nur aus Tags, und die `wails3`-CLI wird mit `go install` installiert. Entfernt `release-v3.yml` und den Nightly-Schritt, der sie ausgelöst hat, in [PR](https://github.com/wailsapp/wails/pull/5946), von @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Hinzugefügt

- macOS-Einstellung für automatische Wiedergabe hinzufügen, mit der sich die Anforderung einer Benutzeraktion für die Medienwiedergabe deaktivieren lässt, in [PR](https://github.com/wailsapp/wails/pull/5512), von @Eyalm321
- Mailbox-Changelog-Eintrag in [PR](https://github.com/wailsapp/wails/pull/5935) nach „Unveröffentlicht“ verschieben, von @leaanthony

## Geändert

- Die macOS-Zoomanimation verwendet in [PR](https://github.com/wailsapp/wails/pull/5945) CADisplayLink oder NSTimer für eine flüssigere Darstellung, von @savely-krasovsky

## Behoben

- Das iOS-Xcode-Projekt so konfigurieren, dass geerbte Linker-Flags erhalten bleiben, und in [PR](https://github.com/wailsapp/wails/pull/5915) -ObjC hinzufügen, von @mortenolsrud
- Übermäßigen Wechsel von TCP-Verbindungen im `wails3 dev`-Asset-Proxy bei großen Frontends beheben, der die kurzlebigen Ports des Hosts erschöpfen und dazu führen konnte, dass unabhängige Prozesse mit `EADDRNOTAVAIL` fehlschlugen
- Fensterspezifisches Ereignis-JavaScript für geordnete Auslieferung und Rückstaukontrolle in [PR](https://github.com/wailsapp/wails/pull/5934) in eine Warteschlange stellen, von @leaanthony

### Hinzugefügt

- Generische asynchrone FIFO-Mailbox für die geordnete Ereigniszustellung in [PR](https://github.com/wailsapp/wails/pull/5851) implementieren, von @savely-krasovsky und @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## Hinzugefügt

- Begrenzten hostseitigen Speicher für übergroße Ereignisse und geordnete JavaScript-Auslieferung in [PR](https://github.com/wailsapp/wails/pull/5930) implementieren, von @leaanthony
- Springen des macOS-Dock-Symbols für das Aufblinken von Fenstern in [PR](https://github.com/wailsapp/wails/pull/5921) implementieren, von @julianstorer

## Behoben

- Der Asset-Server bewahrt in [PR](https://github.com/wailsapp/wails/pull/5931) beim Leeren Fehler der Inhaltstyperkennung und noch nicht geschriebene Präfixe, von @leaanthony
- Abstürze von macOS-Anwendungen verhindern, wenn das Anwendungsmenü aus einem Wails-Callback heraus ersetzt wird
- Unlesbare native Menüs unter Windows 10 1809 / Windows Server 2019 (Build 17763) beheben. Die uxtheme-Exporte für den Dunkelmodus waren auf Build 18334 beschränkt, sodass die Aktivierung des Dunkelmodus auf Anwendungsebene auf diesen Hosts nie ausgeführt wurde: Der Menühintergrund wurde dunkel gezeichnet, Windows zeichnete den Menütext jedoch weiterhin im hellen Design, sodass dunkler Text auf dunklem Hintergrund erschien. Die Ordinalwerte sind ab 17763 vorhanden; die Versionsprüfung stimmt nun damit überein.
- Beheben, dass `w32.GetStockObject` statt `GetStockObject` die Funktion `GetDeviceCaps` aufrief, wodurch für jedes vordefinierte Objekt 0 zurückgegeben wurde.
- Fehlerbehandlung und -meldung beim Download des WebView2-Bootstrappers in [PR](https://github.com/wailsapp/wails/pull/5924) verbessern, von @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## Behoben

- Die Aktivierung von macOS-Apps berücksichtigt in [PR](https://github.com/wailsapp/wails/pull/5897) die Aktivierungsrichtlinie nun nur für reguläre Apps, von @julianstorer
- Nicht initialisierte GTK-Fenster in Linux-Builds in [PR](https://github.com/wailsapp/wails/pull/5898) absichern, von @julianstorer
- Vor dem Laden der URL eine explizite deckende Hintergrundfarbe für Linux-WebKit-Fenster in [PR](https://github.com/wailsapp/wails/pull/5899) festlegen, von @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## Geändert

- Android-Build-Aufgaben verwenden standardmäßig arm64, und deploy-emulator wählt in [PR](https://github.com/wailsapp/wails/pull/5890) die Hostarchitektur aus, von @mortenolsrud

## Behoben

- Zoomzustand von macOS-Fenstern beim Ziehen beibehalten und Bewegung in [PR](https://github.com/wailsapp/wails/pull/5900) reduzieren, von @leaanthony
- Build im Windows-Servermodus durch Hinzufügen von `!server` zur Build-Bedingung `webview_window_windows_nonclient.go` beheben

## v3.0.0-beta.3 - 2026-08-03

## Hinzugefügt

- Abschluss der Beta-Verifizierung von Phase 10 in den Implementierungsdetails in [PR](https://github.com/wailsapp/wails/pull/5881) dokumentieren, von @leaanthony

## Behoben

- Fensterhandle an die Windows-Dunkelmodus-API übergeben und Argumente in [PR](https://github.com/wailsapp/wails/pull/5877) validieren, von @leaanthony
- Ermittlung des Zustands der macOS-Titelleistenschaltflächen für rahmenlose Fenster in [PR](https://github.com/wailsapp/wails/pull/5870) zentralisieren, von @taliesin-ai
- Unlesbaren Text in nativen Menüs verhindern, wenn eine Windows-Anwendung den Dunkelmodus anfordert, während das Windows-App-Design hell ist. Das Menü verwendet nun den passenden hellen nativen Hintergrund, bis Windows dunklen Menütext darstellen kann.
- Unlesbare native Menüs unter Windows 10 1809 / Windows Server 2019 (Build 17763) beheben. Die uxtheme-Exporte für den Dunkelmodus waren auf Build 18334 beschränkt, sodass die Aktivierung des Dunkelmodus auf Anwendungsebene auf diesen Hosts nie ausgeführt wurde: Der Menühintergrund wurde dunkel gezeichnet, Windows zeichnete den Menütext jedoch weiterhin im hellen Design, sodass dunkler Text auf dunklem Hintergrund erschien. Die Ordinalwerte sind ab 17763 vorhanden; die Versionsprüfung stimmt nun damit überein.

## v3.0.0-beta.2 - 2026-08-02

## Geändert

- v3 von Alpha auf Beta hochstufen
- Intelligente Standardwerte und automatisches Ausblenden von Pop-ups im System-Tray dokumentieren, einschließlich Regressionstests für die Auswahl des Klick-Handlers (#5840).
- Der GitHub-Updater schließt Windows-Installer-Assets in [PR](https://github.com/wailsapp/wails/pull/5861) standardmäßig aus, von @leaanthony
- Unterstützung für abgerundete, eckige und benutzerdefinierte Eckradien bei rahmenlosen macOS-Fenstern hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5866) von @leaanthony

## Behoben

- Aktuelle GTK4-Fenstergrößen melden und vom konfigurierten Surface Ereignisse für Größenänderung, Maximierung, Minimierung und Vollbildstatus ausgeben (#5830).
- Linux-WebKit-Absturz beim Senden von Blob oder FormData in fetch-Anfragen behoben, in [PR](https://github.com/wailsapp/wails/pull/5854) von @taliesin-ai
- Der Fetch-Shim übergibt bei fehlenden Blob-/FormData-Headern undefined, in [PR](https://github.com/wailsapp/wails/pull/5865) von @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## Hinzugefügt

## Geändert

- Unterstützung für abgerundete, eckige und benutzerdefinierte Eckradien bei rahmenlosen macOS-Fenstern hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5866) von @leaanthony

## Behoben

- Der Fetch-Shim übergibt bei fehlenden Blob-/FormData-Headern undefined, in [PR](https://github.com/wailsapp/wails/pull/5865) von @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## Hinzugefügt

- Unterstützung für die macOS-DMG-Paketierung mit neuen Optionen und Build-Tasks implementiert, in [PR](https://github.com/wailsapp/wails/pull/5857) von @leaanthony

## Geändert

- Der GitHub-Updater schließt Windows-Installer-Artefakte standardmäßig aus, in [PR](https://github.com/wailsapp/wails/pull/5861) von @leaanthony

## Behoben

- Linux-WebKit-Absturz beim Senden von Blob oder FormData in fetch-Anfragen behoben, in [PR](https://github.com/wailsapp/wails/pull/5854) von @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## Hinzugefügt

- Doppelklick-Aktionen auf die macOS-Titelleiste zum Maximieren oder Minimieren implementiert, in [PR](https://github.com/wailsapp/wails/pull/5853) von @taliesin-ai

## Geändert

- Tutorial für den QR-Dienst auf NewServiceWithOptions umgestellt und Abstände hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5849) von @jeongkyu

## Behoben

- WKWebView bleibt beim Zoomen unter macOS reaktionsfähig, in [PR](https://github.com/wailsapp/wails/pull/5856) von @leaanthony
- Abfragen der GTK4-Fenstergröße korrigiert und vom konfigurierten `GdkSurface` Ereignisse für Größenänderung, Maximierung, Minimierung und Vollbildstatus ausgegeben.

## v3.0.0-alpha2.119 - 2026-07-27

## Behoben

- Mehrsprachige Dokumentation um Architekturdiagramme ergänzt, in [PR](https://github.com/wailsapp/wails/pull/5833) von @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## Hinzugefügt

- Standardpfade für Ein- und Ausgaben der Symbolgenerierung bereitgestellt, in [PR](https://github.com/wailsapp/wails/pull/5825) von @taliesin-ai
- Quell-Einstiegsmodule zu sideEffects in der runtime-package.json hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5797) von @savely-krasovsky
- Abschnitt zu Lizenz und Herkunft zum Leitfaden für Mitwirkende hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5816) von @taliesin-ai

## Behoben

- Auf GTK4 beschränktes CSS für rahmenlose Fenster angewendet, um den Eckradius zu entfernen, in [PR](https://github.com/wailsapp/wails/pull/5800) von @savely-krasovsky
- Fehler beim Ermitteln der Cursorposition unter Windows für Popup-Menüs und die Bildschirmauflistung werden kontrolliert behandelt, in [PR](https://github.com/wailsapp/wails/pull/5789) von @wayneforrest
- Der macOS-Dialog zum Öffnen von Dateien filtert Erweiterungen korrekt und validiert zulässige Dateien anhand ihres Suffixes, in [PR](https://github.com/wailsapp/wails/pull/5678) von @phergul
- Initialisierung des Windows-Dunkelmodus gegen nil-API-Aufrufe abgesichert, in [PR](https://github.com/wailsapp/wails/pull/5793) von @roachadam
- 32-Bit-Buildfehler im Updater behoben: Die Konstante `maxArchiveTotalSize` (2 GiB) überschritt beim Übergeben an `fmt.Errorf` auf `GOARCH=386` den Wertebereich des plattformspezifischen Typs `int`. Sie ist nun explizit als `int64` typisiert.
- nil-Zeiger-Panic beim Start behoben, wenn ein Fenster eine dunkle (oder systembedingt dunkle) Titelleiste auf Windows-Builds verwendet, welche die uxtheme-APIs für den Dunkelmodus nicht laden, etwa Windows 10 1809 / Windows Server 2019 (Build 17763). Die `AllowDarkModeForWindow`-Aufrufe beim Einrichten des Fensterdesigns sind nun gegen nil abgesichert, entsprechend der bereits in `w32.SetMenuTheme` verwendeten Absicherung.

## v3.0.0-alpha2.117 - 2026-07-08

## Hinzugefügt

- Benutzerdefinierte Hit-Test-Logik für Nicht-Client-Bereiche unter Windows implementiert, in [PR](https://github.com/wailsapp/wails/pull/5462) von @savely-krasovsky

## Geändert

- WebView2-Erkennung der Monitorskalierung anhand von UseVisualHosting konfiguriert, in [PR](https://github.com/wailsapp/wails/pull/5761) von @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## Hinzugefügt

- FAQ-Dokumentation auf Funktionen und Anleitungen für Wails v3 ausgerichtet, in [PR](https://github.com/wailsapp/wails/pull/5763) von @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## Behoben

- Behoben, dass `Menu.Update()` das native Menü unter GTK4 Linux nicht neu erstellt (#5659, unabhängig von @puneetdixit200 in #5539 diagnostiziert und behoben)
- Absturz beim Auflisten der macOS-Bildschirme nach einer Anzeigeänderung behoben, indem die Zeichenfolgen für Bildschirm-ID und -Name kopiert und die Anzahl als Momentaufnahme gespeichert werden (#5565, unabhängig von @x-haose in #5584 diagnostiziert und behoben)
- Windows-Absturz behoben, wenn `WM_ERASEBKGND` während eines Übergangs beim Minimieren/Wiederherstellen einen einfarbigen Hintergrund zeichnet und `GetClientRect` nil zurückgibt (Absicherung von @sinspired in #5636 gemeldet)
- Behoben, dass Fehler von Frontend-Bindings immer als Text geparst werden, von @mbaklor in #5690
- Windows-Buildfehler bei Verwendung des Build-Tags `server` behoben. Ursache war, dass in den Windows-GUI-Dateien die Build-Einschränkung `!server` fehlte, die in den entsprechenden macOS- und Linux-Dateien bereits vorhanden ist (#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## Hinzugefügt

- Update-Manifest-Protokoll und Endpoint-Provider implementiert, in [PR](https://github.com/wailsapp/wails/pull/5720) von @taliesin-ai

## Geändert

- Das Binding `webview2` als `v3/internal/webview2` in das Modul v3 integriert und dabei das eigenständige Modul, dessen Workflows für Nightly-Releases und Synchronisierung sowie den Versionsabgleich in go.mod entfernt (v3 ist der einzige Verbraucher), in [PR](https://github.com/wailsapp/wails/pull/5711) von @taliesin-ai

## Behoben

- WebView2-Erkennung der Monitorskalierung und Fehlerbehebung für die Host-Neusynchronisierung bei DPI-Änderungen in den Abschnitt „Unveröffentlicht“ verschoben, in [PR](https://github.com/wailsapp/wails/pull/5750) von @taliesin-ai
- WebView2-COM-Marshalling für float64- und BOOL-Parameter aktualisiert, in [PR](https://github.com/wailsapp/wails/pull/5741) von @wayneforrest
- Panic und nil-Dereferenzierung beim Aktualisieren und Zerstören von Windows-Infobereichssymbolen verhindert, in [PR](https://github.com/wailsapp/wails/pull/5703) von @wayneforrest
- Behoben, dass ausgeblendete Fenster unter Windows nicht korrekt erneut ausgeblendet werden, in [PR](https://github.com/wailsapp/wails/pull/5743) von @wayneforrest
- Sichtbarkeit des WebView2-Controllers mit dem Minimieren, Maximieren und Wiederherstellen des Fensters synchronisiert, in [PR](https://github.com/wailsapp/wails/pull/5742) von @wayneforrest

### Behoben

- WebView2-Erkennung der Monitorskalierung wieder aktivieren und die Host-Neusynchronisierung bei DPI-Änderungen absichern in [PR](https://github.com/wailsapp/wails/pull/5734) von @taliesin-ai, basierend auf der von @randalmurphal validierten Korrektur, mit Überprüfung der Grundursache durch @eleclin und Hardwaretests durch @qq540491950

## v3.0.0-alpha2.113 - 2026-07-04

## Hinzugefügt

- Warnung beim Erstellen eines Release-AAB ohne gesetztes `ANDROID_KEYSTORE_FILE` hinzufügen (Google Play lehnt mit Debug-Signatur signierte Bundles ab) und die Paketierung und Signierung von App Bundles dokumentieren in [PR](https://github.com/wailsapp/wails/pull/5730) von @taliesin-ai
- Mehrsprachige „Warum Wails?“-Dokumentation hinzufügen in [PR](https://github.com/wailsapp/wails/pull/5739) von @taliesin-ai
- Abbildung von Go time.Time auf JS Date oder string in Bindings unterstützen in [PR](https://github.com/wailsapp/wails/pull/5398) von @fbbdev
- Aufgaben zur Paketierung von Android App Bundles (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) für die Einreichung im Play Store hinzufügen — APK-Aufgaben bleiben für lokale Tests und Emulatortests erhalten in [PR](https://github.com/wailsapp/wails/pull/5728) von @mortenolsrud (behebt [#5726](https://github.com/wailsapp/wails/issues/5726))
- Aufgabenziele für physische Android-Geräte hinzufügen und Kamera-/Standortberechtigungen fortsetzen in [PR](https://github.com/wailsapp/wails/pull/5735) von @taliesin-ai

## Geändert

- `webview2` auf v1.0.28 aktualisieren ([Versionshinweise](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- `compileSdk`/`targetSdk` der Android-Vorlage von 34 auf 35 aktualisieren, was Google Play für neue App-Einreichungen voraussetzt, in [PR](https://github.com/wailsapp/wails/pull/5730) von @taliesin-ai

## Behoben

- Eingebrannte Avatar-Masken in sponsorkit korrigieren in [PR](https://github.com/wailsapp/wails/pull/5745) von @leaanthony
- Fehler beheben, durch den die automatische Erstellung einer Android-AVD wegen lexikografischer Versionssortierung das falsche Systemabbild oder die falsche cmdline-tools-Version auswählte, in [PR](https://github.com/wailsapp/wails/pull/5730) von @taliesin-ai
- Fehler beheben, durch den der Einrichtungsassistent eine veraltete Android-NDK-Version vorschlug (jetzt 26.3.11579264, entsprechend der dokumentierten Anforderung), in [PR](https://github.com/wailsapp/wails/pull/5730) von @taliesin-ai
- Französische Dokumentation für SvelteKit und Optionen aktualisieren in [PR](https://github.com/wailsapp/wails/pull/5744) von @leaanthony
- SIGSEGV bei der Auflistung von macOS-Bildschirmen während Anzeigeänderungen beheben in [PR](https://github.com/wailsapp/wails/pull/5516) von @flofreud

## v3.0.0-alpha2.112 - 2026-07-03

## Hinzugefügt

- Go-basierten SVG-Generator für Mitwirkende hinzufügen und Danksagungsseiten der Dokumentation/Website aktualisieren in [PR](https://github.com/wailsapp/wails/pull/5724) von @taliesin-ai
- Abbildung von Go time.Time auf JS Date oder string in Bindings unterstützen in [PR](https://github.com/wailsapp/wails/pull/5398) von @fbbdev

## Geändert

- Node-basierte Pipeline für Sponsorenbilder durch Go-Generator ersetzen in [PR](https://github.com/wailsapp/wails/pull/5719) von @taliesin-ai

## Behoben

- Installationsskript für Abhängigkeiten der Android-Build-Assets korrigieren in [PR](https://github.com/wailsapp/wails/pull/5729) von @taliesin-ai
- Steuerzeichen U+0085 (NEXT LINE) in `ValidateAndSanitizeURL` ablehnen und damit die Abdeckung von Leerraumzeichen durch den URL-Validator vervollständigen
- DWM-Rahmen bei DPI-Änderungen für rahmenlose Fenster neu berechnen in [PR](https://github.com/wailsapp/wails/pull/4785) von @leaanthony
- Fehlschlagende Erkennung von DnD-Ablagezonen bei einer Windows-Skalierung ungleich 100 % beheben in [PR](https://github.com/wailsapp/wails/pull/4632) von @yulesxoxo
- Explizite Objective-C-Speicherverwaltung für Cocoa-Objekte in Darwin-Dialogen, -Menüs, -Taskleistensymbolen und -Benachrichtigungen hinzufügen in [PR](https://github.com/wailsapp/wails/pull/5714) von @taliesin-ai
- Fehler im Linux-CGO-Backend und Probleme mit dem Taskleistensymbol beheben in [PR](https://github.com/wailsapp/wails/pull/5718) von @taliesin-ai

## v3.0.0-alpha2.111 - 2026-07-01

## Hinzugefügt

- HappyTools zur Community-Galerie hinzufügen in [PR](https://github.com/wailsapp/wails/pull/5061) von @Aliuyanfeng
- Unterstützung für die indonesische Sprache und umfassende Dokumentation hinzufügen in [PR](https://github.com/wailsapp/wails/pull/5643) von @triadmoko
- Option DisableMenu zu WindowsWindow hinzufügen in [PR](https://github.com/wailsapp/wails/pull/4813) von @leaanthony

## Geändert

- Taskfile-Vorlage und CLI aktualisieren, um Build-/Paketierungsaufgaben mit GOOS und ARCH auszuführen, in [PR](https://github.com/wailsapp/wails/pull/5617) von @leaanthony

## Behoben

- Problem mit Fenster-Tabs auf dem Mac beheben in [PR](https://github.com/wailsapp/wails/pull/5708) von @taliesin-ai

## Entfernt

- Deutsche MDX-Übersetzungen aus den Bereichen „Mitwirken“, „Funktionen“ und „Anleitungen“ entfernen in [PR](https://github.com/wailsapp/wails/pull/5702) von @taliesin-ai

## v3.0.0-alpha2.110 - 2026-06-30

## Hinzugefügt

- Neuladen und erzwungenes Neuladen der WebView unter macOS implementieren und Wiederherstellung nach Beendigung des WebContent-Prozesses hinzufügen in [PR](https://github.com/wailsapp/wails/pull/5129) von @wayneforrest
- Umfassende deutsche Dokumentation für die Bereiche „Mitwirken“, „Funktionen“ und „Anleitungen“ hinzufügen in [PR](https://github.com/wailsapp/wails/pull/5396) von @leaanthony
- Benachrichtigungen um Ton, Anhänge, Zeitplanung und eine Aktualisierungs-API erweitern in [PR](https://github.com/wailsapp/wails/pull/5333) von @popaprozac

## Behoben

- DWM-Rahmen bei DPI-Änderungen für rahmenlose Fenster neu berechnen in [PR](https://github.com/wailsapp/wails/pull/4785) von @leaanthony
- Fehlschlagende Erkennung von DnD-Ablagezonen bei einer Windows-Skalierung ungleich 100 % beheben in [PR](https://github.com/wailsapp/wails/pull/4632) von @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Hinzugefügt

- Codebeispiele zur EventsEmit-Dokumentation hinzufügen in [PR](https://github.com/wailsapp/wails/pull/5026) von @iamhabbeboy
- Option für visuelles Windows-WebView2-Hosting hinzufügen in [PR](https://github.com/wailsapp/wails/pull/5380) von @MerIijn
- Klustr zur Dokumentation der Community-Galerie hinzufügen in [PR](https://github.com/wailsapp/wails/pull/5536) von @SametKUM
- Kira mit neuen Seiten und einem Changelog-Eintrag zur Community-Galerie hinzufügen in [PR](https://github.com/wailsapp/wails/pull/5685) von @thiennguyen93
- Feedback-Abschnitt zum Leitfaden für den MCP-Dienst hinzufügen in [PR](https://github.com/wailsapp/wails/pull/5694) von @taliesin-ai

## Geändert

- Der Servermodus verfügt jetzt über einen vollwertigen Produktions-Build, der den Desktop-Build-Aufgaben entspricht (#5693). `task build:server` erstellt standardmäßig eine Produktionsbinärdatei (`-tags server,production`, `-trimpath`, bereinigt) und akzeptiert `DEV=true` (Entwicklungsserver), `OBFUSCATED=true` (garble) und `EXTRA_TAGS`. `task run:server` führt einen Entwicklungsserver aus. `Dockerfile.server` / `task build:docker` erstellen zuerst den Produktionsserver (`-tags server,production`) und das Produktions-Frontend; das Image verwendet standardmäßig einen statischen Pure-Go-Build auf distroless/static, wobei `CGO_ENABLED`, `GO_IMAGE` und `RUNTIME_IMAGE` als überschreibbare Build-Argumente für CGO-Apps bereitgestellt werden.

## Behoben

- Absturz beim Schließen eines Fensters mit ausstehenden asynchronen Aufrufen verhindern in [PR](https://github.com/wailsapp/wails/pull/4435) von @leaanthony
- Fensteraktivierung beim Öffnen ausgeblendeter Apps unter Windows verhindern in [PR](https://github.com/wailsapp/wails/pull/5249) von @leaanthony
- Sicherstellen, dass die Verarbeitung von WebKit-Anfragemetadaten, Antwortabschluss und Body-Stream im GTK-Hauptthread erfolgt, in [PR](https://github.com/wailsapp/wails/pull/5668) von @taliesin-ai
- Behebt, dass `Menu.Update()` unter GTK4 Linux das native Menü nicht neu aufgebaut hat (#5659, unabhängig von @puneetdixit200 in #5539 diagnostiziert und behoben)
- Behebt einen Absturz beim Auflisten der macOS-Bildschirme nach einer Anzeigeänderung, indem die Zeichenfolgen für Bildschirm-ID und -Name kopiert und die Anzahl als Momentaufnahme erfasst werden (#5565, unabhängig von @x-haose in #5584 diagnostiziert und behoben)
- Behebt, dass WebView2-Inhalte nach dem Ziehen eines Fensters über Windows-Monitore mit unterschiedlichen DPI-Werten zunächst schrumpften und dann verschwanden. Dazu werden die Controller-Grenzen im `WM_DPICHANGED`-Handler erneut gesetzt, analog zur DPI-Neusynchronisierung beim Aufheben der Minimierung (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Hinzugefügt

- Fügt über `app.GlobalShortcut` globale (systemweite) Tastenkürzel hinzu (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Die Tastenkürzel werden auch ausgelöst, wenn die Anwendung nicht fokussiert ist. Sie sind ohne Drittanbieterabhängigkeiten nativ für jede Plattform implementiert: Carbon-Hotkeys unter macOS, `RegisterHotKey` unter Windows, `XGrabKey` unter X11 und die Schnittstelle für globale Tastenkürzel des XDG Desktop Portal unter Wayland.
- Fügt einen integrierten MCP-Server hinzu: einen Model-Context-Protocol-Server, der automatisch startet, wenn die Anwendung mit dem Tag `mcp` gebaut wird. Dadurch können LLM-Agenten eine laufende Wails-Anwendung testen und steuern – einschließlich Fenstersteuerung, DOM-Inspektion, JavaScript-Auswertung, Aufrufen gebundener Methoden, Ereignissen sowie simulierten Maus- und Tastatureingaben, die mit einem animierten Bildschirmcursor dargestellt werden. Kein Benutzercode erforderlich: Wenn `WAILS_MCP=1` gesetzt ist, wird das Tag `mcp` von `wails3 build`/`wails3 dev` automatisch hinzugefügt. Die Konfiguration erfolgt vollständig über Umgebungsvariablen (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Behoben

- Behebt, dass `Menu.Update()` unter GTK4 Linux das native Menü nicht neu aufgebaut hat (#5659, unabhängig von @puneetdixit200 in #5539 diagnostiziert und behoben)
- Behebt einen Absturz beim Auflisten der macOS-Bildschirme nach einer Anzeigeänderung, indem die Zeichenfolgen für Bildschirm-ID und -Name kopiert und die Anzahl als Momentaufnahme erfasst werden (#5565, unabhängig von @x-haose in #5584 diagnostiziert und behoben)

## v3.0.0-alpha2.107 - 2026-06-27

## Hinzugefügt

- Fügt in [PR](https://github.com/wailsapp/wails/pull/5613) von @leaanthony eine experimentelle Wake-Dokumentation mit Seitenleistennavigation hinzu

## v3.0.0-alpha2.106 - 2026-06-24

## Geändert

- Aktualisiert `webview2` auf v1.0.27.
  - ci(webview2): Release-Build korrigiert (Windows-Cross-Kompilierung + vollständige go.sum) (#5671)\

  **Vollständiger Diff:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- Entfernt go vet aus der Cross-Kompilierung des webview2-Release-Workflows in [PR](https://github.com/wailsapp/wails/pull/5672) von @taliesin-ai
- Aktualisiert das OpenRouter-Modell des automatischen Changelogs in [PR](https://github.com/wailsapp/wails/pull/5670) von @taliesin-ai auf google/gemini-2.5-flash-lite
- Aktualisiert `webview2` auf v1.0.26.

### Korrekturen

- **Stellt nach vorübergehenden COM-Laufzeitfehlern den Betrieb wieder her, statt die Anwendung zu beenden** (#5658, #5580). Zuvor rief `Chromium.errorCallback` bei *jedem* COM-Fehler `os.Exit(1)` auf, sodass eine behebbare Störung nach dem Start die gesamte Anwendung beendete. Laufzeitpfade (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) protokollieren den Fehler nun und stellen den Betrieb wieder her. Insbesondere wird eine fehlerhafte oder nicht vertrauenswürdige Webnachricht in `MessageReceived` nun verworfen, statt den Prozess zu beenden. Dies behebt die Klasse von Abstürzen beim Wechsel zwischen Monitoren mit unterschiedlichen DPI-Werten (#5544, #5650). Fehler in den Pfaden zur Erstellung der Umgebung oder des Controllers bleiben fatal.\

**Vollständiger Diff:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Behoben

- Korrigiert in [PR](https://github.com/wailsapp/wails/pull/5671) von @taliesin-ai den release-webview2-Workflow, sodass go.sum-Dateien korrekt verarbeitet werden
- Behebt in [PR](https://github.com/wailsapp/wails/pull/5659) von @taliesin-ai die Aktualisierung von GTK4-Menüs unter Linux, indem das native Menü geleert und neu aufgebaut wird

## v3.0.0-alpha2.105 - 2026-06-21

## Hinzugefügt

- Fügt `application.System` zur Laufzeiterkennung der Plattform aus gemeinsam genutztem Code hinzu: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (das Build-Tag `server`) und `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)`, um ein einzelnes Ziel direkt zu testen. Es lässt sich für jedes Ziel kompilieren, sodass Verzweigungen ohne Build-Tags möglich sind. Entsprechende Frontend-Hilfsfunktionen (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) sind in `@wailsio/runtime` verfügbar
- Fügt einen Leitfaden „Andere Frontend-Frameworks verwenden“ hinzu, der zeigt, wie ein eigenes Vite-Projekt in `frontend/` eingebunden wird (behandelt Solid, Preact, Lit, SvelteKit, Qwik, Angular usw.)
- Der Assistent `wails3 setup` prüft nun die Toolchain für Mobilplattformen (iOS/Android) – Xcode und die iOS-Simulator-Laufzeit, JDK, Android SDK/NDK und Emulator – und bietet, sofern zutreffend, eine Installation mit einem Klick sowie kopierbare Korrekturen für die Shell-Konfiguration
- Generierte Projekte enthalten eine `frontend/.npmrc`, die eine `minimum-release-age` von 7 Tagen festlegt, um die Gefährdung durch kürzlich veröffentlichte und möglicherweise kompromittierte Pakete zu verringern (wird von pnpm und bun berücksichtigt und von npm ohne Auswirkungen ignoriert)

## Geändert

- Gestaltet alle integrierten Startvorlagen mit einer neuen Hero-Optik mit Neonbergen neu (Web, iOS und Android)
- **TypeScript ist nun der Standard für Startvorlagen und verwendet den unveränderten Vorlagennamen.** `wails3 init` (ohne `-t`) erzeugt ein TypeScript-Projekt; `-t vanilla`, `-t react`, `-t vue` und `-t svelte` verwenden TypeScript. JavaScript-Varianten sind unter `-t vanilla-js`, `-t react-js`, `-t vue-js` und `-t svelte-js` verfügbar. Integrierte Vorlagen deklarieren ihre Sprache mit `typescript:` in `template.yaml`; Community-Vorlagen mit dem Suffix `-ts` funktionieren weiterhin als Ausweichlösung
- Gestaltet den Assistenten `wails3 setup` mit dem Neon-Design „digital Wails“ neu (Milchglas-Vibranz vor einer Bergkulisse)

## Behoben

- Behebt einen Absturz unter Windows beim Wiederherstellen einer Anwendung, die lange genug minimiert war, dass WebView2 angehalten oder dessen Rendering-/GPU-Prozess neu erstellt wurde. Die DPI-Neusynchronisierung beim Minimieren und Wiederherstellen (#5544) greift nun nur auf den WebView2-Controller zu, wenn sich die DPI des Fensters tatsächlich geändert haben. Dadurch werden beim häufigen Wiederherstellen mit unveränderten DPI fatale COM-Aufrufe an einen angehaltenen Controller vermieden (#5605)
- Behebt wiederholte native `SIGABRT`/`SIGSEGV`-Abstürze bei lange laufenden Linux-Anwendungen unter häufiger Asset-/Medienlast (typischerweise innerhalb von `g_object_unref` während der GTK-Hauptschleife). Der Asset-Server schloss `WebKitURISchemeRequest`s aus Worker-Goroutinen ab und rief dadurch nicht threadsichere WebKit2GTK-Funktionen außerhalb des GTK-Hauptthreads auf; der Abschluss (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) erfolgt nun im Hauptthread. Vervollständigt die teilweise Korrektur aus #5566. Betrifft sowohl GTK3- als auch GTK4/WebKitGTK-6.0-Builds (#5631, #5557)
- Behebt sporadisches `fatal error: invalid pointer found on stack` in `setupSignalHandlers` unter Linux/GTK3. Als Signal-`user_data` übergebene Fenster-IDs wurden in einer lokalen Go-Variable vom Typ `unsafe.Pointer` gespeichert. Daher brach der Garbage Collector ab, als er den Wert, der kein Zeiger ist, während einer Stack-Kopie prüfte. Auf der Go-Seite wird die ID nun als Ganzzahl (`uintptr_t`) gespeichert. Damit wird dieselbe Korrektur, die #4958 auf den GTK4-Pfad anwendete, auf den älteren GTK3-Pfad zurückportiert. Dort wurden die C-Signalfunktionen auf `uintptr_t` umgestellt, um `-race`-/checkptr-Fehler zu beseitigen (#5631)

## Entfernt

- Entfernt die Startvorlagen `react-swc`, `preact`, `lit`, `solid`, `qwik` und `sveltekit` sowie deren `-ts`-Varianten. Die unterstützte integrierte Auswahl umfasst nun `vanilla`, `react`, `vue` und `svelte`. Jede verwendet standardmäßig TypeScript und bietet `-js`-JavaScript-Varianten. Jedes andere Framework kann weiterhin verwendet werden, indem Sie [Ihr eigenes Frontend einbinden](https://v3.wails.io/guides/dev/frontend-frameworks) oder eine benutzerdefinierte Vorlage nutzen

## v3.0.0-alpha2.104 - 2026-06-18

## Behoben

- Behebt einen iOS-Absturz (SIGABRT), wenn eine gebundene Go-Servicemethode eine leere Zeichenfolge zurückgibt. Der iOS-Writer für Asset-Antworten prüfte den Body-Zeiger mit `buf != nil` statt dessen Länge, sodass ein Body der Länge null in `&buf[0]` eine Panic auslöste. Nun wird wie bei den Desktop-Writern die Länge geprüft

## v3.0.0-alpha2.103 - 2026-06-15

## Geändert

- Verschiebt native iOS- und Android-Funktionen in Plattformmanager: Ruft sie über `application.IOS.*` und `application.Android.*` auf (z. B. `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`), statt über die bisherigen freien Funktionen `application.IOS*`/`application.Android*` (#5602)
- Benennt Ereignisse der mobilen Bridge um: Plattformübergreifende Ereignisse verwenden jetzt das Präfix `common:*` (z. B. `common:haptic`, `common:location`), plattformspezifische Ereignisse verwenden `ios:*` bzw. `android:*` (z. B. `ios:backgroundTask`, `android:foregroundService`). Das Präfix `native:*` wird nicht mehr verwendet (#5602)

## v3.0.0-alpha.102 - 2026-06-14

## Hinzugefügt

- Fügt einen experimentellen Assistenten `wails3 setup` zur interaktiven Projekteinrichtung und Abhängigkeitsprüfung hinzu
- Fügt `wails3 doctor` das Flag `--json` für maschinenlesbare Ausgaben hinzu
- Fügt dem Befehl `wails3 doctor` einen Abschnitt zum Signierungsstatus hinzu

## Behoben

- Behebt die npm-Erkennung unter Linux, sodass zusätzlich zum Paketmanager auch PATH geprüft wird

## v3.0.0-alpha.101 - 2026-06-13

## Hinzugefügt

- iOS: Native Meldungsdialoge (UIAlertController) und Dialoge zum Öffnen von Dateien, mehreren Dateien und Verzeichnissen (UIDocumentPickerViewController); Speicherdialoge geben einen expliziten Fehler zurück
- iOS: Zwischenablageunterstützung über UIPasteboard
- iOS: Tatsächliche Bildschirmmetriken über UIScreen (Punkte, Pixel, Skalierung und um sichere Bereiche bereinigter Arbeitsbereich)
- iOS: Geräte-Builds (`IOS_PLATFORM=device`), Unterstützung für Codesignierungsidentität, Bereitstellungsprofil und Berechtigungen, Paketierung als `.ipa` sowie `deploy-device` über devicectl
- iOS: Konfigurierbare iOS-Mindestversion (`ios.minIOSVersion` in build/config.yml)
- iOS: `wails3 doctor` meldet unter macOS die Verfügbarkeit von Xcode und dem iOS SDK
- iOS: Systemereignisse für Akku, Netzwerk, Theme, Bildschirmsperre und knappen Arbeitsspeicher werden als `events.IOS.*` und plattformneutrale Anwendungsereignisse vom Typ `events.Common.*` bereitgestellt
- iOS: Native Bridge für Mobilfunktionen (exportiertes `application.IOS*`) – Freigabedialog, Öffnen von URLs, Aktivhalten, Taschenlampe, Abstände sicherer Bereiche, Helligkeit, App-Informationen, Ausrichtungssperre, Statusleiste, Biometrie (Face ID/Touch ID), lokale Benachrichtigungen und sichere Speicherung im Keychain
- iOS: Sensoren und Hardware – Haptik, einmalige Geolokalisierung, Beschleunigungssensor, Näherungssensor, Text-to-Speech, Speicherinformationen, Stromversorgungs-/Akkustatus, Netzwerkstatus, Tastaturabstände und Erkennung von Bildschirmaufnahmen
- iOS: Dokumentation (IOS.md und ein Leitfaden auf der Dokumentationswebsite)
- Android: Native Meldungsdialoge (AlertDialog) und Dialoge zum Öffnen von Dateien bzw. mehreren Dateien (Storage Access Framework, als Kopien in den Cache importiert); Dialoge zum Öffnen von Verzeichnissen und zum Speichern geben einen expliziten Fehler zurück
- Android: Zwischenablageunterstützung über ClipboardManager
- Android: Tatsächliche Bildschirmmetriken über WindowMetrics/DisplayMetrics (dp, Pixel, Skalierung und um Systemleisten bereinigter Arbeitsbereich)
- Android: Laufzeitmethoden für Haptik (`Android.Haptics.Vibrate`), Geräteinformationen (`Android.Device.Info`) und Toast-Meldungen (`Android.Toast.Show`)
- Android: Typisierte Lebenszyklusereignisse (`events.Android.*`, aus events.txt generiert), wobei `ActivityCreated` auf `Common.ApplicationStarted` abgebildet wird
- Android: Die Build-Pipeline erzeugt installierbare Debug- und Release-APKs (`android:run`, `android:package`, `android:package:fat`); Release-Signierung erfolgt standardmäßig mit dem Debug-Keystore oder über Umgebungsvariablen vom Typ `ANDROID_KEYSTORE_*` mit einem echten Keystore
- Android: `wails3 doctor` meldet Android SDK, NDK und JDK
- Android: Systemereignisse für Akku, Netzwerk, Theme, Bildschirmsperre und knappen Arbeitsspeicher werden als `events.Android.*` und plattformneutrale Anwendungsereignisse vom Typ `events.Common.*` bereitgestellt
- Android: Native Bridge für Mobilfunktionen (exportiertes `application.Android*`) – Teilen, Öffnen von URLs, Aktivhalten, Taschenlampe, Abstände sicherer Bereiche, Helligkeit, App-Informationen, Ausrichtungssperre, Statusleiste, Biometrie (BiometricPrompt), lokale Benachrichtigungen und sichere Speicherung mit EncryptedSharedPreferences
- Android: Sensoren und Hardware – Haptik, einmalige Geolokalisierung, Beschleunigungssensor, Näherungssensor, Text-to-Speech, Speicherinformationen, Stromversorgungs-/Akkustatus, Netzwerkstatus, Tastaturabstände und Blockierung von Bildschirmaufnahmen mit FLAG_SECURE
- Android: Dokumentation (ANDROID.md und ein Leitfaden auf der Dokumentationswebsite)
- Beispiel: Das Kitchen-Sink-Beispiel `mobile` erhält die Tabs „Mobil“ und „Hardware“, die die native Funktions-Bridge unter iOS und Android demonstrieren (Pill-Tabs werden auf mehrere Zeilen umgebrochen)
- Mobil: Akku – Beschleunigungs- und Näherungssensor, Taschenlampe sowie die periodische Uhr des Beispiels werden pausiert, wenn die App in den Hintergrund wechselt, und bei der Rückkehr wiederhergestellt (Android lässt den Prozess im Hintergrund weiterlaufen, und die Taschenlampe ist unter iOS ein dauerhaft bestehender Hardwarezustand). Empfänger für Android-Systemereignisse werden nur registriert, solange sich die App im Vordergrund befindet
- iOS: Kameraaufnahme – `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → ein Ereignis vom Typ `native:capture` mit einem Base64-Vorschaubild)
- iOS: Hintergrundausführung – `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (ein UIApplication-Hintergrundaufgabenfenster) und ein konfigurierbares `ios.backgroundModes` (build/config.yml), das `UIBackgroundModes` in die generierte Info.plist einfügt
- Android: Kameraaufnahme – `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (Systemkamera über FileProvider → ein Ereignis vom Typ `native:capture`)
- Android: Vordergrunddienst – `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (ein `WailsForegroundService` mit einer dauerhaften Benachrichtigung hält den Prozess für lang laufende Hintergrundarbeiten aktiv)
- Beispiel: Ein Kamera-Tab demonstriert Foto-/Videoaufnahmen und die Hintergrundausführung (Vordergrunddienst unter Android, Hintergrundaufgabenfenster unter iOS)

## Behoben

- Behebt, dass `getUserMedia` unter Linux immer mit `NotAllowedError` fehlschlug: WebKitGTK verweigert Berechtigungsanfragen, die niemand verarbeitet, und das Signal `permission-request` war nicht verbunden. Kamera und Mikrofon werden jetzt gemäß einer neuen plattformübergreifenden Zuordnung `WebviewWindowOptions.Permissions` (`map[PermissionType]Permission`) behandelt, die sowohl unter Linux (WebKitGTK) als auch unter Windows (WebView2) berücksichtigt wird. Unter Linux, wo es keine native Eingabeaufforderung gibt, sind Kamera und Mikrofon standardmäßig zugelassen (wodurch `getUserMedia` wieder funktioniert) und können mit `PermissionDeny` deaktiviert werden (#5552)
- iOS: `GOOS=ios` lässt sich wieder kompilieren (exportiertes `events.IOS`, Stubs für Namen mobiler Methoden), und Builds mit Produktions-Tag lassen sich kompilieren (Korrekturen an Build-Tags in pkg/application und mehreren Diensten)
- iOS: Go→JS-Ereignisse und ExecJS funktionieren jetzt – die Seite wird beim Start nicht mehr zweimal geladen, und der Handshake `wails:runtime:ready` kann nicht mehr verloren gehen
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` konkurrieren nicht mehr mit dem App-Start; die feste Startverzögerung von 2 Sekunden wurde entfernt
- iOS: Behebt ein C-String-Speicherleck bei jeder Ausführung von JavaScript über Go→JS
- iOS: `hasListeners` gibt jetzt die tatsächliche Listener-Registrierung wieder
- iOS: Framework-Debug-Protokollierung wird aus Produktions-Builds herauskompiliert
- Android: `GOOS=android` lässt sich wieder kompilieren – `events.Android` wurde definiert, das außerhalb der Array-Grenzen zugreifende Listener-Array `events_android.go` entfernt, der Stub für Namen mobiler Methoden hinzugefügt und verhindert, dass Desktop-Linux-Dateien (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) in Android-Builds einfließen
- Android: JS→Go-Bindings funktionieren jetzt – die WebView kann POST-Bodys vom Typ `fetch()` nicht an `shouldInterceptRequest` übermitteln. Daher werden Laufzeitaufrufe über einen JavascriptInterface-Transport (`nativeHandleRuntimeCall`) geleitet, statt wegen eines nil-Anfragebodys abzustürzen
- Android: Laufzeitaufrufe von `Screens.*` geben echte Daten zurück – der ScreenManager wird jetzt beim Start befüllt (er war nie eingebunden, sodass `GetAll` nil zurückgab)
- Android: Das Debug-Logging des Frameworks wird nicht in Produktions-Builds kompiliert und in Debug-Builds unter dem Tag `Wails` an logcat weitergeleitet
- Android: echte `hasListeners`-Registry, Behandlung von JNI-Referenzen und -Ausnahmen sowie Seitenlebenszyklus mit einmaligem Laden (keine doppelte Navigation)
- Behebt den Fehler „Access is denied“ bei `wails3 generate bindings` unter Windows, wenn der Vite-Entwicklungsserver ausgeführt wird, indem generierte Dateien mit dem Ausgabeverzeichnis synchronisiert werden, statt sie darüber umzubenennen (#5515)
- Behebt einen sporadischen schwerwiegenden Absturz unter macOS beim Lesen von Bildschirminformationen nach einer Anzeigeänderung: Die Bildschirm-ID und der Bildschirmname enthielten Zeiger auf automatisch freigegebene `UTF8String`-Puffer, die freigegeben werden konnten, bevor Go sie kopierte (Use-after-free). Die Zeichenfolgen werden nun mit `strdup` verarbeitet und nach der Konvertierung freigegeben. Außerdem läuft die Bildschirmaufzählung in einem expliziten Autorelease-Pool, sodass bei Aufrufen aus Go-Goroutinen kein Speicher mehr verloren geht (#5556)
- Behebt einen sporadischen SIGSEGV unter Linux, wenn der Assetserver ein `WebKitURISchemeRequest` schließt: Das abschließende `g_object_unref` lief in der Assetserver-Goroutine und finalisierte dadurch ein WebKit-GObject außerhalb des GTK-Hauptthreads. Die Unref-Operation wird nun über `g_main_context_invoke` an den GTK-Hauptkontext übergeben (#5557)

## v3.0.0-alpha.100 – 2026-06-13

## Hinzugefügt

- Erweitert `MacWebviewPreferences` um zusätzliche WKWebView-Konfigurationsoptionen: `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` und `ApplicationNameForUserAgent` (#5549)

## Behoben

- Behebt den Fehler „Access is denied“ bei `wails3 generate bindings` unter Windows, wenn der Vite-Entwicklungsserver ausgeführt wird, indem generierte Dateien mit dem Ausgabeverzeichnis synchronisiert werden, statt sie darüber umzubenennen (#5561)
- Behebt, dass JS-Größenänderungsereignisse für rahmenlose Fenster unter Linux nicht ausgelöst werden, und korrigiert die Erkennung von Scrollleistenrändern für rahmenlose Fenster (#5368)
- Behebt den Fehler „invalid cross-device link“ des Updaters unter Windows, wenn sich das temporäre Verzeichnis auf einem anderen Volume als das Installationsverzeichnis befindet (#5560)

## v3.0.0-alpha.99 – 2026-06-10

## Behoben

- Behebt den Fehler „Access is denied“ bei `wails3 generate bindings` unter Windows, wenn der Vite-Entwicklungsserver ausgeführt wird, indem generierte Dateien mit dem Ausgabeverzeichnis synchronisiert werden, statt sie darüber umzubenennen (#5515)

## v3.0.0-alpha.98 – 2026-06-03

## Behoben

- Behebt das Einfrieren der WebKit-Benutzeroberfläche unter Linux im Leerlauf (z. B. bei geöffnetem Inspektor), indem `SA_ONSTACK` nicht mehr für `SIGUSR1` erzwungen wird, da dies die Synchronisierung des GC-Threads von JavaScriptCore beeinträchtigte (#5527)

## v3.0.0-alpha.97 – 2026-05-31

## Hinzugefügt

- Debugging-Seite und Arbeiten mit `runtime/trace` hinzugefügt

## Geändert

- Einige unnötige `_ "embed"`-Importe entfernt und den Code dadurch etwas aufgeräumt

## Behoben

- Behebt, dass Einschränkungen für die minimale Breite und Höhe nach dem Aufheben der Fenstermaximierung unter Windows nicht durchgesetzt werden (#4593)
- Behebt das Durchreichen von Mausklicks im Vollbildmodus bei Verwendung der Fensteroptionen Frameless und Transparent (#4408)

## v3.0.0-alpha.96 – 2026-05-25

## Hinzugefügt

- Unterstützung für Garble-Verschleierung hinzugefügt ([#4563](https://github.com/wailsapp/wails/issues/4563)): stabile Methoden-IDs für Bindings, Integration in Build und Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`) sowie JSON-Struct-Tags für alle zur Laufzeit verwendeten Payloads (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`), damit das Übertragungsformat die Umbenennung exportierter Felder durch Garble übersteht.

## v3.0.0-alpha.95 – 2026-05-20

## Hinzugefügt

- Fehlende Seite zur Projektstruktur hinzugefügt

## Geändert

- Dokumentation: Einige Diagramme auf der Architekturseite wurden für eine übersichtlichere Darstellung in Sequenzdiagramme geändert
- Dokumentation: Hinweis zur erforderlichen Installation von D2 vor der Ausführung hinzugefügt

## Behoben

- Behebt `wails3 generate appimage` mit dem GTK4-Standard: Der Bundler erkennt nun den GTK-Stack aus der Binärdatei, bevor er nach Laufzeitdateien sucht. Dadurch wählt er für GTK4-Builds `libwebkitgtkinjectedbundle.so` (unter `webkitgtk-6.0/`) und für `-tags gtk3`-Builds `libwebkit2gtkinjectedbundle.so` (unter `webkit2gtk-4.1/`). Die `.relr.dyn`-Prüfung berücksichtigt außerdem `libgtk-4.so.1`, sodass das Stripping bei modernen Toolchains unabhängig vom Stack korrekt deaktiviert wird. (#5475)
- Behebt, dass `wails3 generate appimage` beim Aufruf mit einem relativen `-builddir` fehlschlägt: Der Bundler löst `-binary`, `-icon`, `-desktopfile`, `-builddir` und `-outputdir` nun vorab in absolute Pfade auf, damit das während des Ablaufs ausgeführte `s.CD` weder die AppRun-Download-Goroutine noch die nach dem Kopieren ausgeführte `ldd`-Prüfung beeinträchtigt.
- Behebt, dass `wails3 generate appimage` das fertige AppImage nicht nach `-outputdir` verschieben kann, wenn das Feld `Name=` der Desktop-Datei nicht mit dem Basisnamen der Binärdatei übereinstimmt: Der Bundler erzwingt nun über die Umgebungsvariable `OUTPUT`, dass das AppImage-Plug-in von linuxdeploy das AppImage als `<binary>-<arch>.AppImage` schreibt und nicht unter dem aus der Desktop-Datei abgeleiteten Namen.
- Behebt, dass `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` und `Common.SystemDidWake` unter Linux nicht ausgelöst werden, nachdem der GTK4- und WebKitGTK-6.0-Stack in alpha.93 zum Standard erhoben wurde. Das neue standardmäßige `application_linux.go` `run()` rief weder `setupCommonEvents()` auf, das `Linux.*`-Ereignisse an ihre `Common.*`-Entsprechungen weiterleitet, noch `monitorPowerEvents()`. Der DBus-Helfer für die Energieüberwachung wird nun über `application_linux_dbus.go` von den GTK3- und GTK4-Buildpfaden gemeinsam verwendet. (#5474)

## v3.0.0-alpha.94 – 2026-05-19

## Behoben

- Behebt, dass `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` und `Common.SystemDidWake` unter Linux nicht ausgelöst werden, nachdem der GTK4- und WebKitGTK-6.0-Stack in alpha.93 zum Standard erhoben wurde. Das neue standardmäßige `application_linux.go` `run()` rief weder `setupCommonEvents()` auf, das `Linux.*`-Ereignisse an ihre `Common.*`-Entsprechungen weiterleitet, noch `monitorPowerEvents()`. Der DBus-Helfer für die Energieüberwachung wird nun über `application_linux_dbus.go` von den GTK3- und GTK4-Buildpfaden gemeinsam verwendet. (#5474)

## v3.0.0-alpha.93 – 2026-05-17

## Hinzugefügt

- `XDG_SESSION_TYPE` von @leaanthony zur `wails3 doctor`-Ausgabe unter Linux hinzugefügt

## Behoben

- Behebt einen durch den Zugriff von appmenu-gtk-module auf ein noch nicht realisiertes Fenster verursachten Absturz des Fenstermenüs unter Wayland (#4769) von @leaanthony
- Behebt einen Absturz der GTK-Anwendung, wenn der Anwendungsname ungültige Zeichen wie Leerzeichen oder Klammern enthält, von @leaanthony
- Behebt den Fehler „not enough memory“ beim Initialisieren von Drag-and-drop unter Windows (#4701) von @overlordtm
- Behebt eine Race-Condition im Callback-Speicher des Hauptthreads, die durch die falsche Verwendung von RLock beim Löschen aus einer Map verursacht wurde (Linux, macOS, iOS) (#4424) von @leaanthony
- Korrigiert die Variablenverarbeitung bei der Übergabe von Befehlszeilenargumenten an Tasks. Als KEY=VALUE-Paare angegebene CLI-Variablen werden nun ordnungsgemäß initialisiert und während der gesamten Taskausführung weitergegeben.
- Behebt einen NSWindowZoomButton-Konflikt unter macOS: `MaximiseButtonState` und `FullscreenButtonState` wenden nun sowohl beim Start als auch zur Laufzeit den restriktiveren Zustand an; keiner der beiden Setter kann den anderen unbemerkt überschreiben (#5319)
- Behebt mehrere bereits vorhandene Fehler im veralteten GTK3-Build-Pfad (`-tags gtk3`), die CodeRabbit bei #5463 aufgedeckt hat: Starts über Dateizuordnungen überspringen die Startup-Handler nicht mehr; `getTheme` ist nun bereichs- und typsicher; `appName` gibt keinen von GLib verwalteten Speicher mehr frei; `clipboardGet` verursacht kein Speicherleck mehr für das von GTK zurückgegebene `gchar*`; `Calloc` verwendet nun Pointer-Receiver (und `NewCalloc` gibt `*Calloc` zurück), sodass der Pool Allokationen tatsächlich nachverfolgt; `zoomOut` verwendet den Kehrwert von `zoomInFactor` statt eines negativen Multiplikators, der auf 1.0 begrenzt wurde; `execJS` verwendet den vorab zugewiesenen leeren Weltnamen wieder, statt bei jedem Aufruf ein `C.CString("")` zu verlieren; ein Entwicklungs-`fmt.Println` wurde aus `menuItem.setAccelerator` entfernt. Behebt #5465.
- Behebt dasselbe durch einen Value-Receiver verursachte Speicherleck in `Calloc` im standardmäßigen GTK4-Build-Pfad (`linux_cgo.go`): Pointer-Receiver und `NewCalloc() *Calloc` sorgen dafür, dass die `c.String(...)`-Allokationen jedes Fensters tatsächlich nachverfolgt und freigegeben werden.

## v3.0.0-alpha.92 - 2026-05-15

## Hinzugefügt

- Taskfiles geändert, damit sich der verwendete Frontend-Paketmanager über die Option `PACKAGE_MANAGER` steuern lässt
- Vorlagendaten um `{{.Opn}}` und `{{.Cls}}` erweitert, um das Erstellen von Taskfile-Vorlagen vorhersehbarer zu machen

## Geändert

- Einige vorhandene Taskfiles geändert, damit sie `{{.Opn}} and {{.Cls}}` verwenden

## Behoben

- Fatalen `concurrent map read and map write`-Laufzeitfehler in `linuxSystemTray` behoben, der auftrat, wenn das Tray-Menü aktualisiert wurde, während das Panel es las.
- Für WebView2-Fehler und Stacktraces wird nun `log` statt `fmt` verwendet, damit Meldungen nicht verloren gehen, wenn die App unter Windows ohne verbundene Konsole ausgeführt wird.

## v3.0.0-alpha.91 - 2026-05-12

## Geändert

- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5414) von `@github-actions[bot]` aktualisiert
- **BREAKING (macOS):** Das macOS-Koordinatensystem wurde vereinheitlicht, sodass `GetScreens`, `Position` und `SetPosition` denselben Koordinatenraum verwenden: logische Punkte, nach unten zunehmende Y-Werte und `(0,0)` oben links auf dem primären Bildschirm. Dies entspricht Windows, GTK sowie den öffentlichen APIs von Electron und dem Web. Bildschirme, die sich physisch oberhalb des primären Bildschirms befinden, melden nun negative `Bounds.Y`-Werte (zuvor positive), und die Werte von `Position()`/`SetPosition()` werden nun in logischen Punkten statt in `points × primaryScale` angegeben. Die verlustfreie Hin- und Rückkonvertierung `Position()` → `SetPosition()` bleibt erhalten. Absolute Werte, die von früheren Alpha-Builds protokolliert wurden, sowie manuell berechnete Workarounds (z. B. Multiplikation mit `primaryScale` oder Spiegelung von Y anhand einer Bildschirmhöhe) müssen aktualisiert werden. Behebt [#5117](https://github.com/wailsapp/wails/issues/5117).

## Behoben

- DBus-Signalnamen und Textlänge defensiv validiert, um Panics zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/5416) von @leaanthony
- Problem mit der Speichersicherheit bei der GTK-Menüverarbeitung unter Linux behoben, in [PR](https://github.com/wailsapp/wails/pull/5363) von @leaanthony
- Erkennung von NVIDIA-GPUs und Deaktivierung des DMA-BUF-Renderers unter Linux, in [PR](https://github.com/wailsapp/wails/pull/5295) von @leaanthony
- Bildschirmübergreifende Y-Konvertierung von `SetPosition` unter macOS korrigiert: Die Höhe des primären Bildschirms dient nun als globale Referenz, sodass Fenster auf Monitoren, die gegenüber dem primären Bildschirm vertikal versetzt sind, an der richtigen Position erscheinen, in [#5117](https://github.com/wailsapp/wails/issues/5117)
- Git-PR-Vorlage korrigiert, sodass sie auf die richtige Feedback-URL verweist, in [PR](https://github.com/wailsapp/wails/pull/5109) von @wayneforrest
- Mehrere `SetMenu`-Abstürze des Windows-Infobereichs behoben, die durch einen fehlerhaften `DestroyMenu`-Systemaufruf verursacht wurden: Er übergab vier Argumente statt eines, sodass jeder Aufruf FALSE zurückgab und nichts freigab. Außerdem werden HMENU- und HBITMAP-Handles – einschließlich der zur Laufzeit über `MenuItem.SetBitmap` zugewiesenen – beim Neuaufbau des Menüs freigegeben, veraltete Checkbox-/Optionsfeld-Maps in `Win32Menu.Update` zurückgesetzt und ein redundanter `Update()`-Aufruf in `systemtray.updateMenu` entfernt, der die Allokationen verdoppelte. Lang laufende Infobereich-Apps verlieren nun nicht mehr bei jedem Neuaufbau des Menüs GDI-/USER-Objekte.

## v3.0.0-alpha.90 - 2026-05-11

## Hinzugefügt

- Konfigurierbarer Anwendungsname für den WKWebView-User-Agent unter macOS hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5261) von @vinhvoit225
- Indirekte Abhängigkeit github.com/coder/websocket zum gin-service-Beispiel hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5400) von @taliesin-ai
- Unterstützung für tiefgehende Gleichheitsvergleiche zu den Build-Asset-Tests hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5402) von @leaanthony

## Geändert

- Build-Ausgabe im Assets-Verzeichnis zusammengeführt, in [PR](https://github.com/wailsapp/wails/pull/5401) von @taliesin-ai
- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5399) von `@github-actions[bot]` aktualisiert

## Behoben

- Benachrichtigungsobjekt für die macOS-Single-Instance-Nachricht verwendet, in [PR](https://github.com/wailsapp/wails/pull/5289) von @overlordtm
- Windows-Callbacks gebündelt, um den Verlust von Promises unter hoher Last zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/5383) von @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## Hinzugefügt

- Job go<em>test</em>results zum Aggregieren von Go-Testergebnissen hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5316) von @leaanthony

## Geändert

- Große RPC-Nutzlasten bedingt in mehrere segmentierte POST-Anfragen aufgeteilt, in [PR](https://github.com/wailsapp/wails/pull/5369) von @leaanthony
- Vite in allen Frontend-Vorlagen von 5.x.x auf 8.0.0 aktualisiert, in [PR](https://github.com/wailsapp/wails/pull/5386) von @leaanthony
- Portkonfiguration des Vite-Entwicklungsservers auf Umgebungsvariablen umgestellt, in [PR](https://github.com/wailsapp/wails/pull/5365) von @leaanthony
- Vite-Entwicklungsserver in allen Vorlagen so konfiguriert, dass er an 127.0.0.1 gebunden wird, in [PR](https://github.com/wailsapp/wails/pull/5361) von @leaanthony
- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5384) von `@github-actions[bot]` aktualisiert

## Behoben

- Info.plist-Vorlagen-Stubs beim Aktualisieren der Build-Assets bereinigt, in [PR](https://github.com/wailsapp/wails/pull/5312) von @leaanthony
- Veralteten Zustand in macOS-Menüs behoben, indem die Menüelement-Mutatoren (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) synchron im Hauptthread angewendet werden. Dadurch entfällt das `dispatch_async`-Race, das beim schnellen erneuten Öffnen dazu führte, dass Menüs den vorherigen Zustand darstellten (#5002).
- `*_test.go`-Dateien im Entwicklungsmodus ignoriert, um unnötige Neubuilds zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/5203) von @leaanthony
- Segfault bei Menu.Update() verhindert, wenn die App nicht ausgeführt wird, in [PR](https://github.com/wailsapp/wails/pull/5291) von @wucm667
- lastSizeWParam zur Steuerung der Neudarstellung der Menüleiste unter Windows verwendet, in [PR](https://github.com/wailsapp/wails/pull/5382) von @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## Geändert

- HiddenOnTaskbar auf die Verwendung von WS<em>EX</em>TOOLWINDOW umgestellt, in [PR](https://github.com/wailsapp/wails/pull/5371) von @leaanthony
- Abhängigkeiten neu angeordnet und die webview2-Replace-Direktive aus go.mod entfernt, in [PR](https://github.com/wailsapp/wails/pull/5370) von @atterpac
- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5358) von `@github-actions[bot]` aktualisiert

## Behoben

- Generische Indirektionsaliase entfernt und Map-Schlüsseltypen zusammengeführt, in [PR](https://github.com/wailsapp/wails/pull/5331) von @fbbdev

## Entfernt

- PR-master-Workflow löschen und dabei Dokumentation, Go-Tests und übersprungene Tests in [PR](https://github.com/wailsapp/wails/pull/5377) von @leaanthony entfernen

## v3.0.0-alpha.87 - 2026-05-07

## Hinzugefügt

- Koreanische Dokumentation für Wails v3 in [PR](https://github.com/wailsapp/wails/pull/5352) von @leaanthony hinzufügen
- Französische Dokumentation zur Installation und zum Schnellstart in [PR](https://github.com/wailsapp/wails/pull/5354) von @leaanthony hinzufügen
- Portugiesische Dokumentation zu Schnellstart, Konzepten und Community in [PR](https://github.com/wailsapp/wails/pull/5355) von @leaanthony hinzufügen

## v3.0.0-alpha.86 - 2026-05-06

## Hinzugefügt

- Französische Dokumentationslokalisierung in [PR](https://github.com/wailsapp/wails/pull/5328) von @leaanthony hinzufügen
- Deutsche Locale zur Dokumentationswebsite in [PR](https://github.com/wailsapp/wails/pull/5343) von @leaanthony hinzufügen

## Geändert

- Alle 8 übersetzten Locales in der Dokumentationskonfiguration in [PR](https://github.com/wailsapp/wails/pull/5347) von @leaanthony registrieren
- Verschiedene Windows-bezogene Dateien für WebView2 in [PR](https://github.com/wailsapp/wails/pull/5317) von @leaanthony aktualisieren

## Behoben

- Dialog-Dispatch unter Linux in [PR](https://github.com/wailsapp/wails/pull/5340) von @leaanthony zwischen GTK3 und GTK4 aufteilen
- Sicherstellen, dass Dialog-Callbacks im GTK-Thread ausgeführt werden, und dadurch Segmentierungsfehler in [PR](https://github.com/wailsapp/wails/pull/5339) von @leaanthony beheben

## v3.0.0-alpha.85 - 2026-05-05

## Hinzugefügt

- URL der PR-Vorlage in [PR](https://github.com/wailsapp/wails/pull/5179) von @leaanthony zum Repository hinzufügen
- Deutsche Dokumentation für Wails v3 in [PR](https://github.com/wailsapp/wails/pull/5330) von @leaanthony hinzufügen

## v3.0.0-alpha.84 - 2026-05-03

## Hinzugefügt

- Option zum Deaktivieren des Beendens des Vollbildmodus mit der Escape-Taste unter macOS in [PR](https://github.com/wailsapp/wails/pull/5307) von @leaanthony hinzufügen
- Option zum Deaktivieren des Beendens des Vollbildmodus mit der Escape-Taste unter macOS in [PR](https://github.com/wailsapp/wails/pull/5310) von @leaanthony hinzufügen
- Dokumentation für die Community-Präsentation von Pausa in [PR](https://github.com/wailsapp/wails/pull/5288) von @yuseferi hinzufügen

## Geändert

- SVG der Sponsoren in [PR](https://github.com/wailsapp/wails/pull/5308) von `@github-actions[bot]` aktualisieren
- Befehl zur Symbolgenerierung in [PR](https://github.com/wailsapp/wails/pull/5309) von @leaanthony aktualisieren, damit nicht unterstützte Plattformen verarbeitet werden
- Boolesche Vollbild-API durch den dreistufigen ButtonState ersetzen und Plattformbindungen in [PR](https://github.com/wailsapp/wails/pull/5224) von @leaanthony implementieren

## Behoben

- WebView2-Fokusoperationen in [PR](https://github.com/wailsapp/wails/pull/5315) von @leaanthony gegen einen nil-Controller-Zustand absichern
- GitHub-Actions-Workflow in [PR](https://github.com/wailsapp/wails/pull/5313) von @leaanthony aktualisieren, damit er den Basis-Branch des PR korrekt referenziert
- `*_test.go`-Dateien im Entwicklungsmodus ignorieren, um unnötige Neubuilds in [PR](https://github.com/wailsapp/wails/pull/5203) von @leaanthony zu vermeiden
- Segmentierungsfehler durch Menu.Update() bei nicht laufender App in [PR](https://github.com/wailsapp/wails/pull/5291) von @wucm667 verhindern

## v3.0.0-alpha.83 - 2026-05-02

## Hinzugefügt

- InstallScope-Flag und Build-Option für die systemweite bzw. benutzerspezifische Installation in [PR](https://github.com/wailsapp/wails/pull/5094) von @symball hinzufügen
- No-Op-Methode SetScreen zu BrowserWindow hinzufügen, um das Window-Interface in [PR](https://github.com/wailsapp/wails/pull/5294) von @leaanthony zu erfüllen

## Behoben

- NVIDIA-GPUs erkennen und den DMA-BUF-Renderer unter Linux in [PR](https://github.com/wailsapp/wails/pull/5295) von @leaanthony deaktivieren
- Git-PR-Vorlage in [PR](https://github.com/wailsapp/wails/pull/5109) von @wayneforrest korrigieren, sodass sie auf die richtige Feedback-URL verweist
- Eine Reihe von Abstürzen der Windows-Taskleisten-App `SetMenu` beheben, die durch einen fehlerhaften `DestroyMenu`-Systemaufruf verursacht wurden, der vier Argumente statt eines übergab, sodass jeder Aufruf FALSE zurückgab und nichts freigab. Außerdem bei der Neuerstellung von Menüs HMENU- und HBITMAP-Handles freigeben, einschließlich der zur Laufzeit über `MenuItem.SetBitmap` zugewiesenen Handles, veraltete Checkbox-/Optionsfeld-Maps in `Win32Menu.Update` zurücksetzen und einen redundanten `Update()`-Aufruf in `systemtray.updateMenu` entfernen, der die Zuweisungen verdoppelte. Lang laufende Taskleisten-Apps verlieren nun nicht mehr bei jeder Neuerstellung des Menüs GDI-/USER-Objekte.

## v3.0.0-alpha.82 - 2026-05-01

## Behoben

- Generierung von Desktop-Dateien in [PR](https://github.com/wailsapp/wails/pull/5232) von @leaanthony korrigieren, damit der Desktop-Name richtig verarbeitet wird

## v3.0.0-alpha.81 - 2026-04-30

## Geändert

- Zeitplan für Nightly-Releases in [PR](https://github.com/wailsapp/wails/pull/5286) von @leaanthony auf 15:00 UTC anpassen

## Behoben

- Halbierte Werte für Screen Bounds, WorkArea und Size auf Retina-Macs beheben –  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## Geändert

- Dokumentationsabhängigkeiten und Loader für Inhaltssammlungen in [PR](https://github.com/wailsapp/wails/pull/5285) von @leaanthony aktualisieren

## v3.0.0-alpha.79 - 2026-04-29

## Hinzugefügt

- Dem trigger-release-Job in [PR](https://github.com/wailsapp/wails/pull/5270) von @leaanthony die Berechtigung actions: write erteilen

## Geändert

- Release-Task verwendet standardmäßig den master-Branch und aktualisiert die Formulierung des Changelogs in [PR](https://github.com/wailsapp/wails/pull/5283) von @leaanthony
- auto-changelog-Workflow in [PR](https://github.com/wailsapp/wails/pull/5282) von @leaanthony auf die neueste Version aktualisieren
- Workflow-Effizienz durch Pfadfilter und das Entfernen nicht mehr verwendeter Workflows in [PR](https://github.com/wailsapp/wails/pull/5280) von @leaanthony verbessern
- Dokumentation in [PR](https://github.com/wailsapp/wails/pull/5274) von @leaanthony aktualisieren, sodass Beispiellinks auf den master-Branch verweisen
- Dokumentation und Beispiele für v3 in [PR](https://github.com/wailsapp/wails/pull/5272) von @leaanthony aktualisieren

## Behoben

- Reverse-Proxy in [PR](https://github.com/wailsapp/wails/pull/5265) von @AkagiYui um Wiederholungslogik und erzwungenes IPv4 für die Entwicklung erweitern
- Trigger-Workflow für den unveröffentlichten Changelog in [PR](https://github.com/wailsapp/wails/pull/5281) von @leaanthony neu schreiben

## Entfernt

- Shell-Testskripte für verschiedene Testzwecke in [PR](https://github.com/wailsapp/wails/pull/5267) von @leaanthony entfernen
- Workflow zur Bereitstellung der v3-alpha-Dokumentation und CNAME-Eintrag in [PR](https://github.com/wailsapp/wails/pull/5266) von @leaanthony löschen

### Hinzugefügt

- Eintrag „Frontend-Routing“ in [PR](https://github.com/wailsapp/wails/pull/5196) von @leaanthony zur Seitenleistennavigation hinzufügen
- Leitfaden zum Frontend-Routing mit Framework-spezifischen Empfehlungen in [PR](https://github.com/wailsapp/wails/pull/5185) von @leaanthony hinzufügen
- Unterstützung für modale Sheets hinzufügen (macOS)
- ghw-Version für bessere Unterstützung von Apple-Geräten aktualisiert, von @leaanthony (#4977)
- Methode `GetBadge` zum Dock-Dienst hinzufügen
- Flag `-tags` zum Befehl `wails3 build` hinzufügen, um benutzerdefinierte Go-Build-Tags zu übergeben (z. B. `wails3 build -tags gtk4`) (#4957)
- Dokumentation zur automatischen Enum-Generierung im Binding-Generator hinzufügen, einschließlich einer eigenen Enums-Seite und der Navigation in der Seitenleiste (#4972)
- Flag `-tags` zum Befehl `wails3 build` hinzufügen, um benutzerdefinierte Go-Build-Tags zu übergeben (z. B. `wails3 build -tags gtk4`) (#4957)
- Web-API-Beispiele in `v3/examples/web-apis/` hinzufügen, die 41 Browser-APIs demonstrieren, darunter Speicher (localStorage, sessionStorage, IndexedDB, Cache API), Netzwerk (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), Medien (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), Gerät (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), Leistung (Performance API, Mutation Observer, Intersection/Resize Observer), Benutzeroberfläche (Web Components, Pointer Events, Selection, Dialog, Drag and Drop) und weitere
- Beispiel für eine WebView-API-Kompatibilitätsprüfung (`v3/examples/webview-api-check/`) hinzufügen, die plattformübergreifend mehr als 200 Browser-APIs testet
- Paket `internal/libpath` zum Ermitteln nativer Bibliothekspfade unter Linux hinzufügen, mit paralleler Suche, Caching und Unterstützung für Flatpak/Snap/Nix
- **WIP:** Experimentelle Unterstützung für WebKitGTK 6.0 / GTK4 unter Linux hinzufügen, verfügbar über `-tags gtk4` (GTK3/WebKit2GTK 4.1 bleibt die Standardeinstellung)
- Hinweis: Bei Tiling-Window-Managern (z. B. Hyprland, Sway) funktionieren Minimieren und Maximieren möglicherweise nicht wie erwartet, da der Window-Manager die Fenstergeometrie steuert
- Dokumentation zu **einmalig ausgeführten Handlern** unter **Ereignisse in JavaScript überwachen** hinzugefügt, von @AbdelhadiSeddar
- Option `UseApplicationMenu` zu `WebviewWindowOptions` hinzufügen, mit der Fenster unter Windows/Linux das über `app.Menu.Set()` festgelegte Anwendungsmenü übernehmen können, von @leaanthony
- Unterstützung für `.icon`-Dateien (Apple Icon Composer-Format) zur Generierung von Liquid-Glass-Symbolen und Asset-Katalogen (macOS) hinzufügen (#4934), von @wimaha
- Experimentellen Servermodus für Headless-/Web-Bereitstellungen hinzufügen (`-tags server`). Damit können Wails-Apps ohne native GUI-Abhängigkeiten als HTTP-Server ausgeführt werden. Mit `wails3 task build:server` bauen. Einzelheiten finden Sie unter `examples/server`.
- Paket `internal/libpath` zum Ermitteln nativer Bibliothekspfade unter Linux hinzufügen, mit paralleler Suche, Caching und Unterstützung für Flatpak/Snap/Nix
- Option `CollectionBehavior` zu `MacWindow` hinzufügen, um das Fensterverhalten über macOS Spaces und den Vollbildmodus hinweg zu steuern (#4756), von @leaanthony
- Unit-Tests für pkg/application hinzufügen, von @leaanthony
- Unterstützung benutzerdefinierter Protokolle zur MSIX-Paketierung hinzufügen, von @leaanthony
- Erkennung der Desktop-Umgebung unter Linux hinzufügen, [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- Methode `Window.Print()` zur JavaScript-Runtime hinzufügen, um den Druckdialog vom Frontend aus zu öffnen (#4290), von @leaanthony
- `XDG_SESSION_TYPE` zur Ausgabe von `wails3 doctor` unter Linux hinzufügen, von @leaanthony
- Weitere WebKit2-Ereignisse für Änderungen des Ladezustands unter Linux hinzufügen: `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` (#3896), von @leaanthony
- `XDG_SESSION_TYPE` zur Ausgabe von `wails3 doctor` unter Linux hinzufügen, von @leaanthony
- Datei `.desktop` beim Linux-Build generieren, nicht nur bei der Paketierung (#4575)
- Dokumentation der Linux-Runtime-Abhängigkeiten mit distributionsspezifischen Paketnamen und nfpm-Paketierungsbeispielen hinzufügen (#4339), von @leaanthony
- Informationen zur NVIDIA-Treiberversion zur Ausgabe von `wails3 doctor` unter Linux hinzufügen, von @leaanthony
- Ursprung zum Handler für Rohdaten-Nachrichten hinzufügen, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4710)
- Unterstützung für Universal Links unter macOS hinzufügen, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4712)
- Binding-Transportschicht refaktorieren, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4702)
- aria-label-Kennungen zu den helloworld-Vorlagen hinzufügen, damit die Beispiel-App einfach mit Appium-Testclients getestet werden kann, von @chinenual in [PR](https://github.com/wailsapp/wails/pull/4760)
- Ursprung zum Handler für Rohdaten-Nachrichten hinzufügen, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4710)
- Unterstützung für Universal Links unter macOS hinzufügen, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4712)
- Binding-Transportschicht refaktorieren, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4702)
- Typisierte Ereignisse, von @fbbdev und @ianvs in [#4633](https://github.com/wailsapp/wails/pull/4633)
- Beispiel `systray-clock` hinzufügen, das ein Headless-Taskleistensymbol mit Live-Aktualisierungen des Tooltips zeigt (#4653).
- NSIS-Protokollvorlage für Windows hinzugefügt, von @Tolfx in #4510
- Tests für build-assets hinzugefügt, von @Tolfx in #4510
- macOS: Zeigt native Fenstersteuerelemente in der Menüleiste, in [#4588](https://github.com/wailsapp/wails/pull/4588) von @nidib
- macOS-Dock-Dienst hinzufügen, um das App-Symbol im Dock ein- oder auszublenden, von @popaprozac in [PR](https://github.com/wailsapp/wails/pull/4451)
- macOS-Dock-Dienst hinzufügen, um das App-Symbol im Dock ein- oder auszublenden, von @popaprozac in [PR](https://github.com/wailsapp/wails/pull/4451)
- Native Unterstützung für den Liquid-Glass-Effekt unter macOS mit NSGlassEffectView (macOS 15.0+) und NSVisualEffectView als Fallback hinzufügen, einschließlich umfassender Optionen zur Anpassung des Materials, von @leaanthony in [#4534](https://github.com/wailsapp/wails/pull/4534)
- Bereinigung von Browser-URLs, von @leaanthony in [#4500](https://github.dev/wailsapp/wails/pull/4500). Basiert auf [#4484](https://github.com/wailsapp/wails/pull/4484) von @APShenkin.
- Inhaltsschutz unter Windows/Mac hinzufügen, von [@leaanthony](https://github.com/leaanthony), basierend auf der ursprünglichen Arbeit von [@Taiterbase](https://github.com/Taiterbase) in diesem [PR](https://github.com/wailsapp/wails/pull/4241)
- Unterstützung für die Übergabe von CLI-Variablen an Task-Befehle über die Aliasse `wails3 build` und `wails3 package` hinzufügen (#4422), von @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4488)
- Unterstützung für Dropzones, bei denen das Ereignis die Daten des abgelegten Elements als Quelle enthält, von [@atterpac](https://github.com/atterpac) in [#4318](https://github.com/wailsapp/wails/pull/4318)
- `AdditionalLaunchArgs` zu den Optionen von `WindowsWindow` hinzugefügt, damit zusätzliche Befehlszeilenargumente an den WebView2-Browser übergeben werden können, in [PR](https://github.com/wailsapp/wails/pull/4467)
- Automatische Ausführung von go mod tidy nach wails init hinzugefügt, von [@triadmoko](https://github.com/triadmoko) in [PR](https://github.com/wailsapp/wails/pull/4286)
- Windows-Snapassist-Funktion, von @leaanthony in [PR](https://github.dev/wailsapp/wails/pull/4463)
- `AdditionalLaunchArgs` zu den Optionen von `WindowsWindow` hinzugefügt, damit zusätzliche Befehlszeilenargumente an den WebView2-Browser übergeben werden können, in [PR](https://github.com/wailsapp/wails/pull/4467)
- Automatische Ausführung von go mod tidy nach wails init hinzugefügt, von [@triadmoko](https://github.com/triadmoko) in [PR](https://github.com/wailsapp/wails/pull/4286)
- Windows-Snapassist-Funktion, von @leaanthony in [PR](https://github.dev/wailsapp/wails/pull/4463)
- Windows-Implementierung von `getAccentColor` hinzufügen, von [@almas-x](https://github.com/almas-x) in [PR](https://github.com/wailsapp/wails/pull/4427)
- Windows-Implementierung von `getAccentColor` hinzufügen, von [@almas-x](https://github.com/almas-x) in [PR](https://github.com/wailsapp/wails/pull/4427)
- Menüs und Menüleiste im dunklen Windows-Design. Von @leaanthony in [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Integrierte Dienste für verständlichere JS-/TS-Bindings umbenannt, von @popaprozac in [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` zum Abrufen der Akzentfarbe des Systems eines Benutzers. Funktioniert unter macOS. Von [@etesam913](https://github.com/etesam913)
- `window.ToggleFrameless()`-API hinzugefügt, von [@atterpac](https://github.com/atterpac) in [#4137](https://github.com/wailsapp/wails/pull/4137)
- Distributionsspezifische Build-Abhängigkeiten für Linux hinzugefügt, von @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4345)
- Binding-Leitfaden hinzugefügt, von @atterpac in [PR](https://github.com/wailsapp/wails/pull/4404)
- **Organisierte Testinfrastruktur**: Docker-Testdateien in ein eigenes `test/docker/`-Verzeichnis mit optimierten Images verschoben und die Build-Zuverlässigkeit verbessert, von [@leaanthony](https://github.com/leaanthony) in [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Verbesserte Muster für die Ressourcenverwaltung**: Beispiele um die ordnungsgemäße Bereinigung von Event-Handlern und eine kontextbezogene Goroutine-Verwaltung ergänzt, von [@leaanthony](https://github.com/leaanthony) in [#4359](https://github.com/wailsapp/wails/pull/4359)
- Unterstützung für aarch64-AppImage-Builds hinzugefügt, von [@AkshayKalose](https://github.com/AkshayKalose) in [#3981](https://github.com/wailsapp/wails/pull/3981)
- Diagnoseabschnitt zu `wails doctor` hinzugefügt, von [@leaanthony](https://github.com/leaanthony)
- Beim Aufruf einer Dienstmethode das Fenster zum Kontext hinzugefügt, von [@leaanthony](https://github.com/leaanthony)
- `window-call`-Beispiel hinzugefügt, das zeigt, wie sich ermitteln lässt, welches Fenster einen Dienst aufruft, von [@leaanthony](https://github.com/leaanthony)
- Neuer Menüleitfaden von [@leaanthony](https://github.com/leaanthony)
- Verbesserte Behandlung von Panics, von [@leaanthony](https://github.com/leaanthony)
- Neuer Menüleitfaden von [@leaanthony](https://github.com/leaanthony)
- Dokumentationskommentare für die Service-API hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4024](https://github.com/wailsapp/wails/pull/4024)
- Funktion `application.NewServiceWithOptions` hinzugefügt, um Dienste mit zusätzlicher Konfiguration zu initialisieren, von [@leaanthony](https://github.com/leaanthony) in [#4024](https://github.com/wailsapp/wails/pull/4024)
- Verbesserte Menüsteuerung von [@FalcoG](https://github.com/FalcoG) und [@leaanthony](https://github.com/leaanthony) in [#4031](https://github.com/wailsapp/wails/pull/4031)
- Weitere Dokumentation von [@leaanthony](https://github.com/leaanthony)
- Unterstützung für das Abbrechen von Events in Standard-Event-Listenern hinzugefügt, von [@leaanthony](https://github.com/leaanthony)
- Unterstützung für Systray-`Hide`, -`Show` und -`Destroy`, von [@leaanthony](https://github.com/leaanthony)
- Unterstützung für Systray-`SetTooltip`, von [@leaanthony](https://github.com/leaanthony). Ursprüngliche Idee von [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- Paketpfad in Warnungen des Binding-Generators zu nicht unterstützten Typen angegeben, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Unterstützung für generische Aliasse zum Binding-Generator hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Unterstützung für das JSON-Flag `omitzero` zum Binding-Generator hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Direktive `//wails:ignore` hinzugefügt, um die Binding-Generierung für ausgewählte Dienstmethoden zu verhindern, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Direktive `//wails:internal` für Dienste und Modelle hinzugefügt, um Typen zu ermöglichen, die in Go, aber nicht in JS/TS exportiert werden, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Unterstützung für Konstanten eines Aliastyps zum Binding-Generator hinzugefügt, um schwach typisierte Enums zu ermöglichen, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Tests des Binding-Generators für Go-1.24-Funktionen hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4068](https://github.com/wailsapp/wails/pull/4068)
- Unterstützung für macOS 15 „Sequoia“ zu `OSInfo.Branding` hinzugefügt, um die Erkennung der Betriebssystemversion zu verbessern, in [#4065](https://github.com/wailsapp/wails/pull/4065)
- Hook `PostShutdown` hinzugefügt, um benutzerdefinierten Code nach Abschluss des Herunterfahrens auszuführen, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Struct `FatalError` hinzugefügt, um die Erkennung schwerwiegender Fehler in benutzerdefinierten Fehler-Handlern zu unterstützen, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Reihenfolge für das Starten und Herunterfahren von Diensten vereinheitlicht und dokumentiert, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Test-Harness für die Start-/Herunterfahrsequenz der Anwendung sowie Tests für das Starten und Herunterfahren von Diensten hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Methode `RegisterService` zum Registrieren von Diensten nach dem Erstellen der Anwendung hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Feld `MarshalError` zu den Anwendungs- und Dienstoptionen hinzugefügt, um eine benutzerdefinierte Fehlerbehandlung bei Binding-Aufrufen zu ermöglichen, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Abbrechbaren Promise-Wrapper hinzugefügt, der Abbruchanforderungen durch Promise-Ketten weitergibt, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Möglichkeit hinzugefügt, den Abbruch eines Binding-Aufrufs an ein `AbortSignal` zu koppeln, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Unterstützung für `data-wml-*`-Attribute in WML zusätzlich zu den üblichen `wml-*`-Attributen hinzugefügt, von [@leaanthony](https://github.com/leaanthony)
- Methode `Configure` für alle Dienste zur nachträglichen Konfiguration und dynamischen Neukonfiguration hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Der Dienst `fileserver` sendet im nicht konfigurierten Zustand die Antwort 503 Service Unavailable, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Der Dienst `kvstore` stellt im nicht konfigurierten Zustand standardmäßig einen In-Memory-Schlüssel-Wert-Speicher bereit, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Methode `Load` zum Dienst `kvstore` hinzugefügt, um Daten nach Konfigurationsänderungen erneut aus einer Datei zu laden, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Methode `Clear` zum Dienst `kvstore` hinzugefügt, um alle Schlüssel zu löschen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Typ `Level` zum Dienst `log` hinzugefügt, um JS-seitige Konstanten für Protokollierungsstufen bereitzustellen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Methode `Log` zum Dienst `log` hinzugefügt, um die Protokollierungsstufe dynamisch festzulegen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Der Dienst `sqlite` stellt standardmäßig eine In-Memory-Datenbank bereit, wenn er nicht konfiguriert ist, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Methode `Close` zum manuellen Schließen der Datenbank zum Dienst `sqlite` hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Unterstützung für den Abbruch von Abfragemethoden des Dienstes `sqlite` hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Unterstützung für vorbereitete Anweisungen mit JS-Bindings zum Dienst `sqlite` hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Gin-Unterstützung von [Lea Anthony](https://github.com/leaanthony) in [PR](https://github.com/wailsapp/wails/pull/3537), basierend auf der ursprünglichen Arbeit von [@AnalogJ](https://github.com/AnalogJ) in diesem [PR](https://github.com/wailsapp/wails/pull/3537)
- Behoben, dass automatisches Speichern und automatisches Speichern von Passwörtern stets aktiviert waren, von [@oSethoum](https://github.com/osethoum) in [#4134](https://github.com/wailsapp/wails/pull/4134)
- `SetMenu()` zu Fenstern hinzugefügt, um einem Fenster ein Menü zuweisen zu können, von [@leaanthony](https://github.com/leaanthony)
- Unterstützung für Benachrichtigungen hinzugefügt, von [@popaprozac](https://github.com/popaprozac) in [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Unterstützung für Dateizuordnungen unter macOS hinzugefügt, von [@wimaha](https://github.com/wimaha) in [#4177](https://github.com/wailsapp/wails/pull/4177)
- `wails3 tool version` zum Erhöhen semantischer Versionsnummern hinzugefügt, von [@leaanthony](https://github.com/leaanthony)
- Badge-Unterstützung für macOS und Windows hinzugefügt, von [@popaprozac](https://github.com/popaprozac) in [#](https://github.com/wailsapp/wails/pull/4234)
- Unterstützung für registrierte/streng typisierte Ereignisse hinzugefügt, von [@fbbdev](https://github.com/fbbdev) und [@IanVS](https://github.com/IanVS) in [#4161](https://github.com/wailsapp/wails/pull/4161)
- Möglichkeit zur Registrierung von Hooks für benutzerdefinierte Ereignisse hinzugefügt, von [@fbbdev](https://github.com/fbbdev) und [@IanVS](https://github.com/IanVS) in [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)`, um den System-Dateimanager am Pfad `path` zu öffnen, mit optionaler Hervorhebung über `selectFile`, von [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Neues Flag `-git` für den Befehl `wails3 init`, von [@leaanthony](https://github.com/leaanthony)
- Neuer Befehl `wails3 generate webview2bootstrapper` von [@leaanthony](https://github.com/leaanthony)
- Methode `init()` zur manuellen Initialisierung der Runtime zur Runtime hinzugefügt, von [@leaanthony](https://github.com/leaanthony)
- Option `WindowDidMoveDebounceMS` zu den WindowOptions eines Fensters hinzugefügt, von [@leaanthony](https://github.com/leaanthony)
- Single-Instance-Funktion hinzugefügt, von [@leaanthony](https://github.com/leaanthony). Basiert auf dem [v2-PR](https://github.com/wailsapp/wails/pull/2951) von @APshenkin.
- Befehl `wails3 generate template` von [@leaanthony](https://github.com/leaanthony)
- Befehl `wails3 releasenotes` von [@leaanthony](https://github.com/leaanthony)
- Befehl `wails3 update cli` von [@leaanthony](https://github.com/leaanthony)
- Option `-clean` für den Befehl `wails3 generate bindings`, von [@leaanthony](https://github.com/leaanthony)
- AppImage-Linux-Builds für aarch64 (arm64) ermöglicht, von [@AkshayKalose](https://github.com/AkshayKalose) in [#3981](https://github.com/wailsapp/wails/pull/3981)
- Hyperlink für Sponsor hinzugefügt, von @ansxuman in [#3958](https://github.com/wailsapp/wails/pull/3958)
- Unterstützung für die Linux-Paketierung von deb-, rpm- und Arch-Linux-Paket-Builds von
- Unterstützung für universelle Darwin-Builds und -Pakete hinzugefügt von
- Ereignisdokumentation für die Website von
- Vorlagen für sveltekit und sveltekit-ts, die für die Entwicklung ohne SSR eingerichtet sind
- Build-Assets mit dem neuen Befehl `wails3 update build-assets` aktualisiert von
- Beispiel zum Testen der HTML Drag and Drop API von
- Unterstützung für Dateizuordnungen von [leaanthony](https://github.com/leaanthony) in
- Neuer Befehl `wails3 generate runtime` von
- Neue Option `InitialPosition`, um festzulegen, ob das Fenster zentriert werden soll oder
- Methoden `Path` und `Paths` zum Paket `application` hinzugefügt von
- Windows-Optionen `GeneralAutofillEnabled` und `PasswordAutosaveEnabled` hinzugefügt
- Möglichkeit hinzugefügt, das Fenster abzurufen, das eine Dienstmethode aufruft, von
- Optionen `EnabledFeatures` und `DisabledFeatures` für WebView2 hinzugefügt von
- ⊞ Neues DIP-System für verbesserte Unterstützung von High-DPI-Monitoren von
- ⊞ Option für den Namen der Fensterklasse von [windom](https://github.com/windom/) in
- Dienste wurden um Plugin-Funktionalität erweitert. Von
- 🐧 Ereignisse WindowDidMove / WindowDidResize in
- ⊞ Ereignis WindowDidResize in
-  Ereignis ApplicationShouldHandleReopen hinzugefügt, um das Dock verarbeiten zu können
-  getPrimaryScreen/getScreens zur Implementierung hinzugefügt, von @tmclane in
-  Option zum Anzeigen der Symbolleiste im Vollbildmodus unter macOS hinzugefügt von
- 🐧 onKeyPress-Logik hinzugefügt, um Linux-Tastendrücke in einen Accelerator umzuwandeln
- 🐧 Task `run:linux` hinzugefügt von
- Methode `SetIcon` exportiert von [@almas-x](https://github.com/almas-x) in
- `OnShutdown` verbessert von [@almas-x](https://github.com/almas-x) in
- Methode `ToggleMaximise` in der Schnittstelle `Window` wiederhergestellt von
- Weitere Informationen zu `Environment()` hinzugefügt. Von @leaanthony in
- Methode `WebviewWindow.IsFocused` über die Schnittstelle `Window` verfügbar gemacht von
- Unterstützung für mehrere durch Leerzeichen getrennte Trigger-Ereignisse im WML-System von
- ESM-Exporte aus dem gebündelten JS-Runtime-Skript hinzugefügt von
- Flag für den Binding-Generator hinzugefügt, um das gebündelte JS-Runtime-Skript anstelle von
- `setIcon` unter Linux implementiert von [@abichinger](https://github.com/abichinger)
- Flag `-port` zum dev-Befehl hinzugefügt und Umgebungsvariable unterstützt
- Tests für Aufrufe gebundener Methoden hinzugefügt von
- ⊞ `SetIgnoreMouseEvents` für bereits erstellte Fenster hinzugefügt von
-  Möglichkeit zum Festlegen der Stapelreihenfolge eines Fensters hinzugefügt von

### Behoben

- Halbierte `Screen.Bounds`, `WorkArea` und `Size` auf Retina-Macs behoben, indem NSScreen-Punktwerte in den Feldern `Physical*` in Gerätepixel umgerechnet und `Screen.X`/`Y` auf oberster Ebene befüllt werden, sodass die Erkennung sich berührender Monitore und die Platzierung im Arbeitsbereich korrekt sind, in [PR](https://github.com/wailsapp/wails/pull/5168) von @wayneforrest
- Datenrennen in ScreenManager behoben, das bei einer Änderung der Anzeigekonfiguration einen WebKit-DisplayLink-Deadlock verursacht (z. B. beim Anschließen eines externen Monitors während des Ruhezustands/Aufwachens)
- Setzt CFBundleIconName direkt auf appicon, wenn Assets.car vorhanden ist, in [PR](https://github.com/wailsapp/wails/pull/5154) von @symball
- Behebt, dass `wails3 doctor` unter Fedora, openSUSE, Arch und NixOS falsche WebKitGTK-Pakete meldet. Die 4.0-Fallback-Einträge wurden entfernt, da v3 die 4.1-API zur Kompilierzeit benötigt (#5071)
- Korrigiert den webkit2gtk-Paketnamen des Doctors für openSUSE (`webkit2gtk4_1-devel` → `webkit2gtk3-devel`, der korrekte openSUSE-Paketname) (#5071)
- Behebt den `Unexpected token '<'`-Fehler, wenn `/wails/custom.js` im Desktop-Entwicklungsmodus fehlt. Ein expliziter 404-Handler für `/wails/custom.js` und eine Validierung von `Content-Type` ohne Berücksichtigung der Groß-/Kleinschreibung in `loadOptionalScript` wurden hinzugefügt, um zu verhindern, dass HTML-SPA-Fallbacks als JavaScript eingeschleust werden. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Behebt den Hervorhebungsstatus des Systemleistenmenüs unter macOS – das Symbol zeigt nun den ausgewählten Zustand an, wenn das Menü geöffnet ist (#4910)
- Behebt, dass das an die Systemleiste angeheftete Fenster unter macOS hinter anderen Fenstern erscheint – es verwendet nun die korrekte Popup-Fensterebene (#4910)
- Korrigiert fehlerhafte `@wailsio/runtime`-Importbeispiele in der gesamten Dokumentation (#4989)
- Behebt, dass rahmenlose Fenster unter darwin nicht minimiert werden können (#4294)
- Behebt minutenlange Hänger während 20-30 und `wails3 build`, indem `wails3 dev` von der Aktualitätsprüfung von go-task ausgeschlossen wird. Zuvor führte das `node_modules/`-Glob dazu, dass go-task jede Datei in `sources: "**/*"` aufzählte und ihre Prüfsumme berechnete (über `node_modules/` Dateien mit umfangreichen Abhängigkeiten wie MUI), was insbesondere unter Windows/NTFS langsam war (50000-100000#4939)
- Behebt einen GTK4-Buildfehler, der durch eine Kollision des C-`Screen`-typedefs mit X11 Xlib.h verursacht wurde (#4957)
- Vereinheitlicht die Dock-Badge-Methoden unter macOS
- Behebt, dass `InvisibleTitleBarHeight` auf alle macOS-Fenster statt nur auf rahmenlose Fenster oder Fenster mit transparenter Titelleiste angewendet wurde (#4960)
- Behebt das Ruckeln/Zittern des Fensters beim Ändern der Größe über die oberen Ecken mit aktiviertem `InvisibleTitleBarHeight`, indem die Einleitung des Ziehvorgangs nahe den Fensterrändern übersprungen wird (#4960)
- Behebt die Generierung abgebildeter Typen mit Enum-Schlüsseln in JS/TS-Bindings (#4437) von @fbbdev
- Behebt, dass Datei-Drag-and-drop unter Windows bei einer Anzeigeskalierung ungleich 100 % nicht funktioniert
- Behebt, dass internes HTML5-Drag-and-drop unter Windows nicht funktioniert, wenn das Ablegen von Dateien aktiviert ist
- Behebt, dass Dateiablagekoordinaten unter Windows im falschen Pixelraum vorliegen (physische statt CSS-Pixel)
- Behebt, dass Datei-Drag-and-drop unter Linux mit Hover-Effekten nicht zuverlässig funktioniert
- Behebt, dass internes HTML5-Drag-and-drop unter Linux nicht funktioniert, wenn das Ablegen von Dateien aktiviert ist
- Behebt, dass beim Ein-/Ausblenden eines Fensters unter Linux/GTK4 gelegentlich der minimierte Zustand wiederhergestellt wird, indem `gtk_window_present()` verwendet wird (#4957)
- Behebt, dass das Abrufen/Festlegen der Fensterposition unter Linux/GTK4 stets 0,0 zurückgibt, indem bedingte X11-Unterstützung über `XTranslateCoordinates`/`XMoveWindow` hinzugefügt wird (#4957)
- Behebt, dass die maximale Fenstergröße unter Linux/GTK4 nicht erzwungen wird, indem eine signalbasierte Größenbegrenzung als Ersatz für das entfernte `gtk_window_set_geometry_hints` hinzugefügt wird (#4957)
- Korrigiert die DPI-Skalierung unter Linux/GTK4 durch eine korrekte Berechnung von PhysicalBounds und Unterstützung für fraktionale Skalierung über `gdk_monitor_get_scale` (GTK 4.14+)
- Behebt die Duplizierung von Menüeinträgen beim Erstellen neuer Fenster unter Linux/GTK4
- Behebt die Generierung abgebildeter Typen mit Enum-Schlüsseln in JS/TS-Bindings (#4437) von @fbbdev
- Behebt, dass Datei-Drag-and-drop unter Windows bei einer Anzeigeskalierung ungleich 100 % nicht funktioniert
- Behebt, dass internes HTML5-Drag-and-drop unter Windows nicht funktioniert, wenn das Ablegen von Dateien aktiviert ist
- Behebt, dass Dateiablagekoordinaten unter Windows im falschen Pixelraum vorliegen (physische statt CSS-Pixel)
- Behebt, dass Datei-Drag-and-drop unter Linux mit Hover-Effekten nicht zuverlässig funktioniert
- Behebt, dass internes HTML5-Drag-and-drop unter Linux nicht funktioniert, wenn das Ablegen von Dateien aktiviert ist
- Korrigiert die DPI-Skalierung unter Linux/GTK4 durch eine korrekte Berechnung von PhysicalBounds und Unterstützung für fraktionale Skalierung über `gdk_monitor_get_scale` (GTK 4.14+)
- Behebt die Duplizierung von Menüeinträgen beim Erstellen neuer Fenster unter Linux/GTK4
- Behebt die Generierung abgebildeter Typen mit Enum-Schlüsseln in JS/TS-Bindings (#4437) von @fbbdev
- Behebt das Problem mit „Geisterfenstern“ unter macOS, das dadurch verursacht wurde, dass AppKit-APIs in App.Window.Current() nicht vom Main Thread aus aufgerufen wurden (#4947) von @wimaha
- Behebt, dass HTML-`<input type="file">` unter macOS nicht funktioniert, durch Implementierung von WKUIDelegate runOpenPanelWithParameters (#4862)
- Behebt, dass natives Datei-Drag-and-drop bei Verwendung des npm-Moduls `@wailsio/runtime` unter macOS/Linux nicht funktioniert (#4953) von @leaanthony
- Behebt die Binding-Generierung für paketübergreifende Typaliase (#4578) von @fbbdev
- Behebt einen Absturz von OpenFileDialog unter Linux aufgrund einer Verletzung der GTK-Threadsicherheit (#3683) von @ddmoney420
- Behebt einen SIGSEGV-Absturz beim Aufruf von `Focus()` für ein ausgeblendetes oder zerstörtes Fenster (#4890) von @ddmoney420
- Behebt eine mögliche Panic beim Festlegen eines leeren Symbols oder Bitmaps unter Linux (#4923) von @ddmoney420
- Behebt einen Absturz von ErrorDialog bei einem Aufruf aus einem Service-Binding unter macOS (#3631) von @leaanthony
- Ermöglicht die Anzeige von Menüs unter Windows in `v3\examples\dialogs` von @ndianabasi
- Behebt eine Race Condition, die beim Neuladen der Seite einen TypeError verursacht (#4872) von @ddmoney420
- Behebt fehlerhafte Ausgaben der Binding-Generator-Tests, indem globaler Zustand aus der Methode `Collector.IsVoidAlias()` entfernt wird (#4941) von @fbbdev
- Behebt, dass die `<input type="file">`-Dateiauswahl unter macOS nicht funktioniert (#4862) von @leaanthony
- Behebt, dass `Position()` und `SetPosition()` unter macOS unterschiedliche Koordinatensysteme verwenden, wodurch sich die Fensterposition beim Speichern/Wiederherstellen des Zustands verschiebt (#4816) von @leaanthony
- Behebt den SetProcessDpiAwarenessContext-Fehler „Access is denied“, wenn die DPI-Awareness bereits über das Anwendungsmanifest festgelegt wurde (#4803)
- Aktualisiert die Dokumentationsseite zu Tastenkombinationen und korrigiert den Typ des Callback-Parameters für `KeyBinding.Add` von @ndianabasi
- Korrigiert die Dokumentation zur Generierung benutzerdefinierter Bindings: Statt `-o String` muss `-d String` verwendet werden
- Behebt, dass das Menü unter `menu.Update()` seine untergeordneten Elemente nicht entfernt
- Korrigiert veraltete Verweise auf die Manager-API in der Dokumentation (31 Dateien wurden auf das neue Muster mit beispielsweise `app.Window.New()` und `app.Event.Emit()` aktualisiert) von @leaanthony
- Behebt einen Linux-Absturz bei einer Panic in an JS gebundenen Go-Methoden, der durch das Überschreiben von Signal-Handlern durch WebKit verursacht wird (#3965) von @leaanthony
- Behebt, dass SaveFileDialog.SetFilename() unter Linux keine Wirkung hat (#4841) von @samstanier
- Behebt, dass Ablagekoordinaten im Drag-and-drop-Beispiel als undefined angezeigt werden
- Behebt das Fehlschlagen der Erstellung von macOS-App-Bundles, wenn APP_NAME Leerzeichen enthält (Problem mit der Klammererweiterung)
- Behebt eine Index-out-of-bounds-Panic unter Windows beim Aufruf von Service-Methoden (goccy/go-json zurückgesetzt)
- Behebt, dass Datei-Drag-and-drop unter Windows bei einer Anzeigeskalierung ungleich 100 % nicht funktioniert
- Behebt, dass internes HTML5-Drag-and-drop unter Windows nicht funktioniert, wenn das Ablegen von Dateien aktiviert ist
- Behebt, dass Dateiablagekoordinaten unter Windows im falschen Pixelraum vorliegen (physische statt CSS-Pixel)
- Behebt, dass Datei-Drag-and-Drop unter Linux mit Hover-Effekten nicht zuverlässig funktioniert
- Behebt defektes internes HTML5-Drag-and-Drop, wenn Datei-Drop unter Linux aktiviert war
- Aktualisiert alle Befehle in Taskfile.yml-Dateien für alle Betriebssysteme, damit Leerzeichen in Variablen wie `APP_NAME` unterstützt werden, von @ndianabasi
- Behebt den Fehler bei Befehlsargumenten beim Ausführen der Task 'build:universal:lipo:go' unter Linux, von @wux1an
- Behebt den Docker-Fehler „undefined symbol: **<em>ubsan</em>handle_xxxxxxx“ beim Ausführen von 'wails3 build GOOS=darwin GOARCH=arm64' unter Linux, von @wux1an
- Fasst die Dokumentation zu benutzerdefinierten Protokollen zusammen und ergänzt Abschnitte zu Universal Links, von @leaanthony
- Behebt einen Absturz des Windows-Menüs im Infobereich beim wiederholten Klicken auf das Symbol durch eine Schutzprüfung gegen gleichzeitige TrackPopupMenuEx-Aufrufe (#4151), von @leaanthony
- Verhindert einen Anwendungsabsturz, wenn systray.Run() vor app.Run() aufgerufen wird, von @leaanthony
- Behebt einen Absturz unter macOS beim Umschalten der Fenstersichtbarkeit über Hide()/Show(), wenn ApplicationShouldTerminateAfterLastWindowClosed aktiviert ist (#4389), von @leaanthony
- Behebt ein Speicherleck in Kontextmenüs unter macOS und Windows beim wiederholten Öffnen von Menüs (#4012), von @leaanthony
- Behebt, dass native Ressourcen von Kontextmenüs unter macOS nicht wiederverwendet wurden und das Menü bei jeder Anzeige neu erstellt wurde (#4012), von @leaanthony
- Behebt, dass ein Klick auf das macOS-Dock-Symbol ausgeblendete Fenster nicht anzeigte, wenn die Anwendung mit `Hidden: true` gestartet wurde (#4583), von @leaanthony
- Behebt, dass der Druckdialog unter macOS aufgrund eines falschen Fensterzeigertyps im CGO-Aufruf nicht geöffnet wurde (#4290), von @leaanthony
- Behebt einen Absturz des Fenstermenüs unter Wayland, der durch den Zugriff von appmenu-gtk-module auf ein noch nicht realisiertes Fenster verursacht wurde (#4769), von @leaanthony
- Behebt einen Absturz der GTK-Anwendung, wenn der Anwendungsname ungültige Zeichen enthält (Leerzeichen, Klammern usw.), von @leaanthony
- Behebt den Fehler „not enough memory“ beim Initialisieren von Drag-and-Drop unter Windows (#4701), von @overlordtm
- Behebt, dass der Dateimanager unter Linux aufgrund einer fehlerhaften URI-Escapierung das falsche Verzeichnis öffnete (#4397), von @leaanthony
- Behebt AppImage-Buildfehler auf modernen Linux-Distributionen (Arch, Fedora 39+, Ubuntu 24.04+) durch automatische Erkennung von `.relr.dyn`-ELF-Abschnitten und Deaktivierung des Strippings (#4642), von @leaanthony
- Behebt, dass `wails doctor` auf Fedora- bzw. DNF-basierten Systemen WebKit-Pakete fälschlicherweise als installiert meldete (#4457), von @leaanthony
- Behebt, dass die Standardeinstellung `config.yml` mit einem Produktions-Build `wails3 dev` ausführte, von @mbaklor
- Behebt Buildfehler durch iOS-Service-Stubs, die ein nicht vorhandenes Paket importierten, von @leaanthony
- Behebt Fehler des Typs „no formatting directives“, die durch strukturiertes Logging in Debug-/Info-Methoden verursacht wurden, von @leaanthony
- Entfernt temporäre Debug-Ausgaben, die versehentlich aus der Zusammenführung der mobilen Plattform übernommen wurden, von @leaanthony
- Behebt einen WebKitGTK-Absturz unter Wayland mit NVIDIA-GPUs (Fehler 71 Protocol error), indem der DMA-BUF-Renderer automatisch deaktiviert wird, von @leaanthony
- Behebt, dass der Alphawert in `application.WebviewWindowOptions.BackgroundColour` unter Linux ignoriert wurde ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Behebt, dass das Windows-Symbol im Infobereich nicht standardmäßig das Anwendungssymbol verwendete, wenn kein benutzerdefiniertes Symbol angegeben war (#4704)
- Verfolgt die Eigentümerschaft von `HICON`, sodass nur vom Benutzer erstellte Handles zerstört werden, und verhindert dadurch Abstürze beim Neustart des Explorers (#4653).
- Gibt beim Zerstören den Windows-Listener für das Systemdesign und beibehaltene Infobereichssymbole frei, damit keine Goroutinen und Gerätekontexte mehr lecken (#4653).
- Kürzt Tooltips im Infobereich auf 127 UTF-16-Einheiten, um die Beschädigung von Surrogatpaaren und Mehrbyte-Glyphen zu vermeiden (#4653).
- Behebt einen Fehler der Windows-Paket-Task (#4667)
- Korrigiert die appicon-Variable für Linux-AppImage im Linux-Taskfile [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Behebt einen Windows-Buildfehler, der durch eine Signaturänderung in go-webview2 v1.0.22 verursacht wurde (#4513, #4645)
- Korrigiert die appicon-Variable für Linux-AppImage im Linux-Taskfile [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Korrigiert den Protokollbereich in linux desktop.tmpl durch Entfernen von `<.Info.Protocol>` bis `<.Protocol>`, von @Tolfx in #4510
- Behebt einen Neudefinitionsfehler in der Liquid-Glass-Demo in [#4542](https://github.com/wailsapp/wails/pull/4542), von @Etesam913
- Behebt Aktualisierungen des Infobereichsmenüs unter Linux in [#4604](https://github.com/wailsapp/wails/issues/4604), von [@JackDoan](https://github.com/JackDoan)
- Behebt das weiße Fenster, das unter Windows beim Erstellen eines ausgeblendeten Fensters erschien, von @leaanthony in [#4612](https://github.com/wailsapp/wails/pull/4612)
- Korrigiert den Importpfad des Benachrichtigungspakets in der Dokumentation, von @rxliuli in [#4617](https://github.com/wailsapp/wails/pull/4617)
- Behebt, dass Drag-and-Drop bei Verwendung des npm-Pakets @wailsio/runtime nicht funktionierte (#4489), von @leaanthony in #4616
- Windows: Behebt das Flackern des Fensters beim Start und die fehlerhafte Anzeige ausgeblendeter Fenster in [PR](https://github.com/wailsapp/wails/pull/4600), von @leaanthony.
- Probleme bei der Maximierung der Fenstergröße unter Wayland behoben (https://github.com/wailsapp/wails/issues/4429) von [@samstanier](https://github.com/samstanier)
- Probleme beim Maximieren der Fenstergröße unter Wayland behoben (https://github.com/wailsapp/wails/issues/4429) von [@samstanier](https://github.com/samstanier)
- Behebt einen Neudefinitionsfehler in der Liquid-Glass-Demo in [#4542](https://github.com/wailsapp/wails/pull/4542), von @Etesam913
- Behebt ein Problem, durch das AssetServer unter macOS abstürzen konnte, in [#4576](https://github.com/wailsapp/wails/pull/4576), von @jghiloni
- Behebt ein Kompilierungsproblem beim Build mit NextJs. Behoben in [#4585](https://github.com/wailsapp/wails/pull/4585) von @rev42
- Korrigiert die Pipelines für Nightly-Releases in [#4597](https://github.com/wailsapp/wails/pull/4597), von @riadafridishibly
- Behebt einen Neudefinitionsfehler in der Liquid-Glass-Demo in [#4542](https://github.com/wailsapp/wails/pull/4542), von @Etesam913
- Behebt ein Problem, durch das AssetServer unter macOS abstürzen konnte, in [#4576](https://github.com/wailsapp/wails/pull/4576), von @jghiloni
- Behebt ein Kompilierungsproblem beim Build mit NextJs. Behoben in [#4585](https://github.com/wailsapp/wails/pull/4585) von @rev42
- Korrigiert die Pipelines für Nightly-Releases in [#4597](https://github.com/wailsapp/wails/pull/4597), von @riadafridishibly
- Behebt einen Neudefinitionsfehler in der Liquid-Glass-Demo in [#4542](https://github.com/wailsapp/wails/pull/4542), von @Etesam913
- Korrigiert SetBackgroundColour unter Windows, von @PPTGamer in [PR](https://github.com/wailsapp/wails/pull/4492)
- Aktualisiert die Dokumentation entsprechend den Änderungen durch das Refactoring der Manager API, von @yulesxoxo in [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Korrigiert die appicon-Variable der Linux-.desktop-Datei im Linux-Taskfile [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Aktualisiert die Dokumentation entsprechend den Änderungen durch das Refactoring der Manager API, von @yulesxoxo in [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Behebt einen unter [#4456](https://github.com/wailsapp/wails/issues/4456) gemeldeten Fehler durch Dereferenzierung eines nil-Zeigers unter Windows, von @leaanthony in [#4460](https://github.com/wailsapp/wails/pull/4460)
- Fügt Unterstützung für `allowsBackForwardNavigationGestures` in WKWebView unter macOS hinzu, um Navigationsgesten durch Wischen mit zwei Fingern zu ermöglichen (#1857)
- Behebt ein Problem, durch das onClick bei Menüeinträgen nicht funktionierte, die anfänglich als deaktiviert festgelegt waren; behoben von @leaanthony in [PR #4469](https://github.com/wailsapp/wails/pull/4469). Vielen Dank an @IanVS für die erste Untersuchung.
- Behebt ein Problem, durch das der Vite-Server nach einem fehlgeschlagenen Build nicht bereinigt wurde (#4403)
- Panic beim Schließen oder Abbrechen eines `SaveFileDialog` unter Windows behoben. Behoben von @hkhere in [PR](https://github.com/wailsapp/wails/pull/4284)
- Drag-and-drop auf HTML-Ebene unter Windows behoben von [@mbaklor](https://github.com/mbaklor) in [#4259](https://github.com/wailsapp/wails/pull/4259)
- Unterstützung für `allowsBackForwardNavigationGestures` in WKWebView unter macOS hinzugefügt, um Navigationsgesten durch Wischen mit zwei Fingern zu ermöglichen (#1857)
- Behebt ein Problem, durch das onClick bei Menüeinträgen nicht funktionierte, die anfänglich als deaktiviert festgelegt waren; behoben von @leaanthony in [PR #4469](https://github.com/wailsapp/wails/pull/4469). Vielen Dank an @IanVS für die erste Untersuchung.
- Behebt ein Problem, durch das der Vite-Server nach einem fehlgeschlagenen Build nicht bereinigt wurde (#4403)
- Parsen von Benachrichtigungen unter Windows behoben von @popaprozac in [PR](https://github.com/wailsapp/wails/pull/4450)
- Prüfung auf Abhängigkeiten des Windows SDK im doctor-Befehl behoben von [@kodumulo](https://github.com/kodumulo) in [#4390](https://github.com/wailsapp/wails/issues/4390)
- Dereferenzierung eines nil-Zeigers in processURLRequest auf dem Mac behoben von [@etesam913](https://github.com/etesam913) in [#4366](https://github.com/wailsapp/wails/pull/4366)
- Linux-Fehler behoben, der gefilterte Dialoge verhinderte, von [@bh90210](https://github.com/bh90210) in [#4287](https://github.com/wailsapp/wails/pull/4287)
- Probleme mit dem Bearbeiten-Menü unter Windows und Linux behoben von [@leaanthony](https://github.com/leaanthony) in [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Minimale Systemversion in macOS-.plist-Dateien von 10.13.0 auf 10.15.0 aktualisiert von [@AkshayKalose](https://github.com/AkshayKalose) in [#3981](https://github.com/wailsapp/wails/pull/3981)
- Problem mit übersprungenen Fenster-IDs behoben von [@leaanthony](https://github.com/leaanthony)
- Problem mit nil-Menüs beim Aufruf von RegisterContextMenu behoben von [@leaanthony](https://github.com/leaanthony)
- Abhängigkeitszyklen in der Ausgabe des Binding-Generators behoben von [@fbbdev](https://github.com/fbbdev) in [#4001](https://github.com/wailsapp/wails/pull/4001)
- Use-before-define-Fehler in der Ausgabe des Binding-Generators behoben von [@fbbdev](https://github.com/fbbdev) in [#4001](https://github.com/wailsapp/wails/pull/4001)
- Build-Flags werden an den Binding-Generator übergeben, von [@fbbdev](https://github.com/fbbdev) in [#4023](https://github.com/wailsapp/wails/pull/4023)
- Pfade im Windows-Taskfile auf Schrägstriche umgestellt, damit es auf Nicht-Windows-Plattformen funktioniert, von [@leaanthony](https://github.com/leaanthony)
- Mac- und Mac-JS-Ereignisse behoben von [@leaanthony](https://github.com/leaanthony)
- Event-Deadlock unter macOS behoben von [@leaanthony](https://github.com/leaanthony)
- `Parameter incorrect`-Fehler bei der Fensterinitialisierung unter Windows behoben, wenn HTML, aber kein JS bereitgestellt wurde, von [@leaanthony](https://github.com/leaanthony)
- Größe des Antwortpräfixes korrigiert, das der Asset-Server zur Erkennung des Inhaltstyps verwendet, von [@fbbdev](https://github.com/fbbdev) in [#4049](https://github.com/wailsapp/wails/pull/4049)
- Verarbeitung von Nicht-404-Antworten auf dem Root-Index-Pfad im Asset-Server behoben von [@fbbdev](https://github.com/fbbdev) in [#4049](https://github.com/wailsapp/wails/pull/4049)
- Undefiniertes Verhalten im Binding-Generator beim Prüfen von Eigenschaften generischer Typen behoben von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ausgabe des Binding-Generators für Modelle korrigiert, wenn der zugrunde liegende Typ nicht dieselben Eigenschaften wie der benannte Wrapper besitzt, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ausgabe des Binding-Generators für Map-Schlüsseltypen und die Vorverarbeitung korrigiert von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ausgabe des Binding-Generators für Structs korrigiert, die Marshaler-Schnittstellen implementieren, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Erkennung von Typzyklen mit generischen Typen im Binding-Generator behoben von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ungültige Verweise auf nicht exportierte Modelle in der Ausgabe des Binding-Generators behoben von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Injizierten Code an das Ende der Service-Dateien verschoben von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Verarbeitung von Fehlern beim Schließen von Dateien im Binding-Generator behoben von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Warnungen für Services unterdrückt, die Lebenszyklus- oder HTTP-Methoden, aber keine anderen gebundenen Methoden definieren, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Problem behoben, durch das Nicht-React-Vorlagen bei Verwendung eines hellen Systemfarbschemas die Hello-World-Fußzeile nicht anzeigten, von [@marcus-crane](https://github.com/marcus-crane) in [#4056](https://github.com/wailsapp/wails/pull/4056)
- Ausgeblendete Menüeinträge unter macOS behoben von [@leaanthony](https://github.com/leaanthony)
- Verarbeitung und Formatierung von Fehlern in Nachrichtenprozessoren behoben von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Übersprungenes Herunterfahren von Services beim Beenden der Anwendung behoben von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Sichergestellt, dass Menüaktualisierungen im Hauptthread erfolgen, von [@leaanthony](https://github.com/leaanthony)
- Der Mechanismus zum Ziehen und Ändern der Größe ist jetzt robuster und entspricht dem erwarteten Plattformverhalten genauer, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- [#4097](https://github.com/wailsapp/wails/issues/4097) behoben: Webpack/Angular verwirft den Initialisierungscode der Runtime; behoben von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Anfänglich ausgeblendete Menüeinträge behoben von [@IanVS](https://github.com/IanVS) in [#4116](https://github.com/wailsapp/wails/pull/4116)
- Problem behoben, durch das assetFileServer bei einer Anfrage ohne Dateierweiterung keine `.html`-Dateien auslieferte, wenn `[request]` nicht vorhanden war, `[request].html` jedoch schon
- Pfade für die Symbolgenerierung korrigiert von [@robin-samuel](https://github.com/robin-samuel) in [#4125](https://github.com/wailsapp/wails/pull/4125)
- Problem behoben, durch das die Ereignisse `fullscreen`, `unfullscreen`, `unminimise` und `unmaximise` nicht ausgelöst wurden, von [@oSethoum](https://github.com/osethoum) in [#4130](https://github.com/wailsapp/wails/pull/4130)
- NSIS-Fehler aufgrund eines falschen Präfixes bei der Standardversion in der Konfiguration behoben von [@robin-samuel](https://github.com/robin-samuel) in [#4126](https://github.com/wailsapp/wails/pull/4126)
- Problem behoben, durch das die Dialogs-Runtime-Funktion unter Windows maskierte Pfade zurückgab, von [TheGB0077](https://github.com/TheGB0077) in [#4188](https://github.com/wailsapp/wails/pull/4188)
- Pfad zur Webview2-Erkennung in HKCU korrigiert von [@leaanthony](https://github.com/leaanthony).
- Eingabeproblem unter macOS behoben von [@leaanthony](https://github.com/leaanthony).
- Dateinamen der Task-Datei für die Windows-Symbolgenerierung korrigiert von [@yulesxoxo](https://github.com/yulesxoxo) in [#4219](https://github.com/wailsapp/wails/pull/4219).
- Transparenzproblem bei rahmenlosen Fenstern korrigiert von [@leaanthony](https://github.com/leaanthony), basierend auf der Arbeit von @kron.
- Fokusaufrufe bei deaktiviertem oder minimiertem Fenster korrigiert von [@leaanthony](https://github.com/leaanthony), basierend auf der Arbeit von @kron.
- Problem behoben, durch das Taskleistensymbole nach einem Neustart der Taskleiste nicht angezeigt wurden, von [@leaanthony](https://github.com/leaanthony), basierend auf der Arbeit von @kron.
- Behoben, dass fallbackResponseWriter Flush() nicht implementierte, in [#4245](https://github.com/wailsapp/wails/pull/4245)
- Behoben, dass fallbackResponseWriter Flush() nicht implementierte, von [@superDingda] in [#4236](https://github.com/wailsapp/wails/issues/4236)
- Abstürze beim Schließen eines macOS-Fensters mit einem ausstehenden asynchronen Aufruf einer Go-gebundenen Funktion behoben von [@joshhardy](https://github.com/joshhardy) in [#4354](https://github.com/wailsapp/wails/pull/4354)
- Race-Condition beim Start des Windows-Effizienzmodus behoben von [@leaanthony](https://github.com/leaanthony)
- Bereinigung von Windows-Symbol-Handles korrigiert von [@leaanthony](https://github.com/leaanthony).
- `OpenFileManager` unter Windows korrigiert von [@PPTGamer](https://github.com/PPTGamer) in [#4375](https://github.com/wailsapp/wails/pull/4375).
- Optionen für minimale/maximale Breite unter Linux korrigiert von @atterpac in [#3979](https://github.com/wailsapp/wails/pull/3979)
- Typdefinitionen für TypeScript-Vorlagen durch Anheben der npm-Version von @atterpac in [#3966](https://github.com/wailsapp/wails/pull/3966)
- CSS-Referenz der SvelteKit-Vorlage korrigiert von @atterpac in [#3945](https://github.com/wailsapp/wails/pull/3945)
- Sichergestellt, dass wichtige Callbacks in window run() im Hauptthread aufgerufen werden, von [@leaanthony](https://github.com/leaanthony)
- Beispiele für die Verzeichnisauswahl in Dialogen korrigiert von [@leaanthony](https://github.com/leaanthony)
- Neue chinesische Fehlerseite für den Fall erstellt, dass index.html fehlt, von [@leaanthony](https://github.com/leaanthony)
-  Sichergestellt, dass der Callback `windowDidBecomeKey` im Hauptthread ausgeführt wird, von [@leaanthony](https://github.com/leaanthony)
-  Vollbildmodus für rahmenlose Fenster unterstützt von [@leaanthony](https://github.com/leaanthony)
-  Logik zum Zerstören von Fenstern verbessert von [@leaanthony](https://github.com/leaanthony)
-  Logik für die Fensterposition bei Verknüpfung mit Taskleistensymbolen korrigiert von [@leaanthony](https://github.com/leaanthony)
-  Vollbildmodus für rahmenlose Fenster unterstützt von [@leaanthony](https://github.com/leaanthony)
- Ereignisbehandlung korrigiert von [@leaanthony](https://github.com/leaanthony)
- Logik zum Herunterfahren von Fenstern korrigiert von [@leaanthony](https://github.com/leaanthony)
- Die gemeinsame Task-Datei generiert für TypeScript-Vorlagen nun standardmäßig TypeScript-Bindings, von [@leaanthony](https://github.com/leaanthony)
- Schließen der Anwendung bei einer WM_CLOSE-Nachricht korrigiert, wenn keine Fenster geöffnet sind bzw. nur ein Taskleistensymbol vorhanden ist, von [@mmalcek](https://github.com/mmalcek) in [#3990](https://github.com/wailsapp/wails/pull/3990)
- Garble-Build korrigiert von @5aaee9 in [#3192](https://github.com/wailsapp/wails/pull/3192)
- Windows-NSIS-Builds korrigiert von [@leaanthony](https://github.com/leaanthony)
- Deadlock im Linux-Dialog bei Mehrfachauswahl behoben, der durch ein nicht geschlossenes
- Plattformübergreifende Bereinigung von .syso-Dateien während des Windows-Builds korrigiert von
- Kompilierung des amd64-AppImage korrigiert von @atterpac in
- Aktualisierung der Build-Assets korrigiert von @ansxuman in
- Implementierung von `OnClick` und `OnRightClick` für das Linux-Taskleistensymbol korrigiert von @atterpac
- Behoben, dass `AlwaysOnTop` auf dem Mac nicht funktionierte, von
-  Behoben, dass `application.NewEditMenu` ein Duplikat enthielt
- 🐧 aarch64-Kompilierung korrigiert
- ⊞ Menüelemente für Optionsgruppen korrigiert von
- Fehler beim Erstellen einer ausführbaren .app unter macOS behoben, wenn 'name' und 'outputfilename'
- Fehler bei der Verwendung von customEventProcessor im Drag-and-drop-Beispiel behoben von
- 🐧 Durch die Ergänzung von IgnoreMouseEvents verursachten Linux-Kompilierungsfehler behoben von
- ⊞ Fehler bei der Generierung der syso-Symboldatei behoben von
- 🐧 Korrektur für die native Ausführung unter Wayland übernommen aus
- Interne Dienstmethoden nicht binden in
- ⊞ Panic beim Start des Taskleistensymbols behoben in
- Interne Dienstmethoden nicht binden in
- ⊞ Panic beim Start des Taskleistensymbols behoben in
- Umfassende Überarbeitung der Menüelemente und der Ereignisbehandlung. Verbessert vorerst hauptsächlich macOS. Von
- Tests nach der Überarbeitung der Plugins und Ereignisse korrigiert in
- ⊞ Warnung zu `Failed to unregister class Chrome_WidgetWin_0` behoben. Von
- Modulprobleme
- Nachrichtenübermittlung für Größenänderungsereignisse korrigiert von [atterpac](https://github.com/atterpac) in
- 🐧 Fehler bei der Theme-Verarbeitung unter NixOS behoben von
- Projektinstallation über Laufwerksgrenzen hinweg unter Windows korrigiert von
- CSS der React-Vorlage korrigiert, damit die Fußzeile angezeigt wird, von
- Zombie-Prozesse beim Arbeiten im Entwicklungsmodus durch Aktualisierung auf die neueste refresh-Version behoben
- Einbindung der WebKit-Datei im AppImage korrigiert von [Atterpac](https://github.com/atterpac)
- Prüfung von apt-Paketen in Doctor korrigiert von [Atterpac](https://github.com/Atterpac) in
- Einfrieren der Anwendung beim Beenden unter Darwin behoben von @5aaee9 in
- Hintergrundfarben der Beispiele unter Windows korrigiert von
- Standardkontextmenüs korrigiert von [mmghv](https://github.com/mmghv) in
- Hexadezimalwerte für Pfeiltasten unter Darwin korrigiert von
- Drag-and-drop unter Windows funktionsfähig gemacht. Hinzugefügt von
- Linux-Fehler in Doctor für den Fall behoben, dass der Benutzer nicht über die richtigen Treiber verfügt
- DPI-Skalierung beim Start unter Windows korrigiert. Geändert von [@almas-x](https://github.com/almas-x) in
- Ersetzungszeile in `go.mod` zur Verwendung relativer Pfade korrigiert. Behebt Windows-Pfade mit
- Klickbehandlung für das macOS-Taskleistensymbol ohne verknüpftes Fenster korrigiert von
- Fehlschlagenden Windows-Build aufgrund einer unbekannten Option korrigiert von
- Absturz unter Windows beim Linksklick auf das Taskleistensymbol behoben, wenn kein
- Falsche baseURL beim zweimaligen Öffnen eines Fensters korrigiert von @5aaee9 in PR
- Reihenfolge der if-Zweige in der Methode `WebviewWindow.Restore` korrigiert von
- `startURL` über mehrere Aufrufe von `GetStartURL` hinweg korrekt berechnet, wenn
- JS-Typ der Struktur `Screen` an ihr Go-Gegenstück angepasst von
- Methode `WML.Reload` korrigiert, um die ordnungsgemäße Bereinigung registrierter Ereignisse sicherzustellen
- Sofortiges Schließen benutzerdefinierter Kontextmenüs unter Linux behoben von
- Ausgabepfad und Dateiendung der vom Binding erzeugten Modelldateien korrigiert
- Importpfade der vom Binding im JS-Code erzeugten Modelldateien korrigiert
- Drag-and-drop auf einigen Linux-Distributionen korrigiert von
- Fehlenden Task für macOS bei Verwendung von `wails3 task dev` korrigiert von
- Zuweisung an eine nil-Map bei der Ereignisregistrierung korrigiert von
- Unmarshalling gebundener Methodenparameter korrigiert von
- Verarbeitung mehrerer Rückgabewerte gebundener Methoden korrigiert von
- npm-Erkennung durch Doctor korrigiert, wenn npm nicht mit dem Systempaketmanager installiert wurde
- Fehlende MicrosoftEdgeWebview2Setup.exe ergänzt. Vielen Dank an
- Zufälligen Absturz unter Linux aufgrund der Verarbeitung von Fenster-IDs korrigiert von @leaanthony. Basierend auf
- Absturz von systemTray.setIcon unter Linux korrigiert von
- Sichergestellt, dass der Fensterrahmen beim ersten Aufruf der Funktion `setFrameless` angewendet wird auf

### Geändert

- **INKOMPATIBLE ÄNDERUNG**: Map-Schlüssel in generierten JS/TS-Bindings sind jetzt als optional gekennzeichnet, um die Semantik von Go-Maps korrekt abzubilden. Der Zugriff auf Map-Werte in TypeScript gibt jetzt `T | undefined` statt `T` zurück, sodass Nullprüfungen oder Assertions erforderlich sind (#4943) von `@fbbdev`
- Verwendung von `Event` gemäß den Änderungen in `@wailsio/runtime` zu `Events` geändert und die entsprechenden Funktionsaufrufe in der Dokumentation unter `Features/Events/Event System` angepasst von @AbdelhadiSeddar
- `EnabledFeatures`, `DisabledFeatures` und `AdditionalBrowserArgs` aus den fensterspezifischen Optionen in `Options.Windows` auf Anwendungsebene verschoben (#4559) von @leaanthony
- README für das Beispiel `Drag N Drop` aktualisiert und hervorgehoben, dass `Internal Drag and Drop` im Beispiel demonstriert wird, von @ndianabasi
- Verschiedene Debug-Protokollmeldungen von Info auf Debug umgestellt (von @mbaklor)
- **INKOMPATIBLE ÄNDERUNG:** `EnableDragAndDrop` in den Fensteroptionen in `EnableFileDrop` umbenannt
- **INKOMPATIBLE ÄNDERUNG:** `DropZoneDetails` im Ereigniskontext in `DropTargetDetails` umbenannt
- **INKOMPATIBLE ÄNDERUNG:** Methode `DropZoneDetails()` auf `WindowEventContext` in `DropTargetDetails()` umbenannt
- **INKOMPATIBLE ÄNDERUNG:** Ereignis `WindowDropZoneFilesDropped` entfernt; stattdessen `WindowFilesDropped` verwenden
- **INKOMPATIBLE ÄNDERUNG:** HTML-Attribut von `data-wails-dropzone` in `data-file-drop-target` geändert
- **INKOMPATIBLE ÄNDERUNG:** CSS-Hover-Klasse von `wails-dropzone-hover` in `file-drop-target-active` geändert
- **INKOMPATIBLE ÄNDERUNG:** Optionen `DragEffect`, `OnEnterEffect` und `OnOverEffect` aus Windows entfernt (sie gehörten zum entfernten IDropTarget)
- Für die gesamte JSON-Verarbeitung zur Laufzeit (Methoden-Bindings, Ereignisse, WebView-Anfragen, Benachrichtigungen, kvstore) auf goccy/go-json umgestellt, wodurch die Leistung um 21-63 % steigt und die Speicherallokationen um 40-60 % sinken
- Layout der BoundMethod-Struktur optimiert und isVariadic-Flag zwischengespeichert, um den Overhead pro Aufruf zu reduzieren
- Für Methoden mit `<=8` Argumenten einen auf dem Stack allozierten Argumentpuffer verwendet, um Heap-Allokationen zu vermeiden
- Erfassung der Ergebnisse von Methodenaufrufen optimiert, um bei einzelnen Rückgabewerten die Slice-Allokation zu vermeiden
- sync.Map für den MIME-Typ-Cache verwendet, um die Leistung bei gleichzeitigem Zugriff zu verbessern
- Pufferpool zum Lesen des Anfragekörpers beim HTTP-Transport verwendet
- CloseNotify-Kanal bei der Inhaltstyperkennung verzögert alloziert, um die Allokationen pro Anfrage zu reduzieren
- Debug-CSS-Protokollierung aus dem Asset-Server entfernt
- Zuordnung von Dateiendungen zu MIME-Typen um mehr als 50 gängige Webformate erweitert (Schriftarten, Audio, Video usw.)
- Dokumentation für die Window-Optionen `X/Y` aktualisiert von @ruhuang2001
- Dokumentation zu `Frontend Runtime` um weitere Optionen zum Generieren von Frontend-Bindings ergänzt von @ndianabasi
- Dokumentationsseite für den Wails-Asset-Server v3 aktualisiert von @ndianabasi
- **INKOMPATIBLE ÄNDERUNG**: Dialogfunktionen auf Paketebene entfernt (`application.InfoDialog()`, `application.QuestionDialog()` usw.). Stattdessen den Manager `app.Dialog` verwenden: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Dialogdokumentation an die tatsächliche API angepasst: `app.Dialog.*` und `AddButton()` mit Callbacks verwenden (nicht `SetButtons()`), `SetDefaultButton(*Button)` (kein String), `AddFilter()` (nicht `SetFilters()`), `SetFilename()` (nicht `SetDefaultFilename()`) sowie `app.Dialog.OpenFile().CanChooseDirectories(true)` zur Ordnerauswahl
- **INKOMPATIBLE ÄNDERUNG**: Produktions-Builds sind jetzt die Standardeinstellung. Für Entwicklungs-Builds `DEV=true` in den Taskfiles festlegen. Neues Projekt für Beispiele generiert von @leaanthony
- Beim Auslösen eines benutzerdefinierten Ereignisses mit keinem oder einem Datenargument wird der Datenwert direkt dem Feld Data zugewiesen, ohne ihn in ein Slice einzuschließen, von [@fbbdev](https://github.com/fbbdev) in [#4633](https://github.com/wailsapp/wails/pull/4633)
- Windows-Tray-Icons berücksichtigen jetzt `SystemTray.Show()`/`Hide()`, indem sie `NIS_HIDDEN` umschalten. Dadurch können Anwendungen tatsächlich verschwinden und wieder eingeblendet werden (#4653).
- Die Tray-Registrierung verwendet aufgelöste Icons erneut, setzt `NOTIFYICON_VERSION_4` einmalig und aktiviert `NIF_SHOWTIP`, damit Tooltips nach einem Neustart von Explorer wiederhergestellt werden (#4653).
- macOS: Zum Zentrieren von Fenstern `visibleFrame` statt `frame` verwendet, um die Bereiche von Menüleiste und Dock auszuschließen
- macOS: Zum Zentrieren von Fenstern `visibleFrame` statt `frame` verwendet, um die Bereiche von Menüleiste und Dock auszuschließen
- Beim Ausführen von `wails3 update build-assets` mit dem Parameter `-config` werden über die Parameter `-product*` festgelegte Werte
- `window.NativeWindowHandle()` -> `window.NativeWindow()` von @leaanthony in [#4471](https://github.com/wailsapp/wails/pull/4471)
- Interne Fensterverarbeitung refaktoriert von @leaanthony in [#4471](https://github.com/wailsapp/wails/pull/4471)
- `application.WindowIDKey` und `application.WindowNameKey` entfernt (durch `application.WindowKey` ersetzt) von [@leaanthony](https://github.com/leaanthony)
- ContextMenuData gibt jetzt einen String statt any zurück, von [@leaanthony](https://github.com/leaanthony)
- In JS/TS-Bindings werden Klassenfelder mit Array-Typen fester Länge jetzt mit der erwarteten Länge initialisiert, statt leer zu sein, von [@fbbdev](https://github.com/fbbdev) in [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData gibt jetzt einen String statt any zurück, von [@leaanthony](https://github.com/leaanthony)
- `application.NewService` akzeptiert keine Optionen mehr als optionalen Parameter (stattdessen `application.NewServiceWithOptions` verwenden), von [@leaanthony](https://github.com/leaanthony) in [#4024](https://github.com/wailsapp/wails/pull/4024)
- Abhängigkeit von `nanoid` entfernt von [@leaanthony](https://github.com/leaanthony)
- Window-Beispiel für Mica-, Acryl- und Registerkarten-Fensterstile aktualisiert von [@leaanthony](https://github.com/leaanthony)
- In JS/TS-Bindings wurden die Modelldateien `internal.js/ts` entfernt; alle Modelle befinden sich jetzt in `models.js/ts`, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- In JS/TS-Bindings werden benannte Typen niemals als Aliase für andere benannte Typen dargestellt; das bisherige Verhalten ist jetzt auf Aliase beschränkt, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- In JS/TS-Bindings werden im Klassenmodus Struct-Felder, deren Typ ein Typparameter ist, als optional markiert und nie automatisch initialisiert, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- ESLint aus Vorlagen entfernt, von [@IanVS](https://github.com/IanVS) in [#4059](https://github.com/wailsapp/wails/pull/4059)
- Copyright-Jahr auf 2025 aktualisiert, von [@IanVS](https://github.com/IanVS) in [#4037](https://github.com/wailsapp/wails/pull/4037)
- Dokumentation für event.Sender hinzugefügt, von [@IanVS](https://github.com/IanVS) in [#4075](https://github.com/wailsapp/wails/pull/4075)
- Unterstützung für Go 1.24, von [@leaanthony](https://github.com/leaanthony)
- `ServiceStartup`-Hooks werden jetzt beim Aufruf von `App.Run` statt in `application.New` ausgeführt, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- `ServiceStartup`-Fehler werden jetzt von `App.Run` zurückgegeben, statt den Prozess zu beenden, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Binding- und Dialogaufrufe aus JS werden jetzt mit Fehlerobjekten statt mit Zeichenketten abgelehnt, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Positionierung des Systray-Menüs unter Windows verbessert, von [@leaanthony](https://github.com/leaanthony)
- Die JS-Runtime wurde nach TypeScript portiert, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Die Runtime wird unmittelbar beim Import initialisiert; es muss nicht mehr auf das Laden des Fensters gewartet werden, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Die Runtime exportiert keine init-Methode mehr. Sie kann mit einem Import für Seiteneffekte initialisiert werden, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Gebundene Methoden geben jetzt ein `CancellablePromise` zurück, das bei Abbruch mit einem `CancelError` abgelehnt wird. Das eigentliche Ergebnis des Aufrufs wird verworfen, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Integrierte Diensttypen heißen jetzt einheitlich `Service`, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Funktionen zum Erstellen integrierter Dienste mit Optionen heißen jetzt einheitlich `NewWithConfig`, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Die Methode `Select` des Dienstes `sqlite` heißt jetzt `Query`, um mit den Go-APIs konsistent zu sein, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Vorlagen: Runtime nach "dependencies" verschoben und package.json-Dateien geordnet, von [@IanVS](https://github.com/IanVS) in [#4133](https://github.com/wailsapp/wails/pull/4133)
- Erstellt und signiert App-Bundles während der Entwicklung ad hoc, um bestimmte macOS-APIs zu aktivieren, von [@popaprozac](https://github.com/popaprozac) in [#4171](https://github.com/wailsapp/wails/pull/4171)
- Build-Assets in plattformspezifische Verzeichnisse verschoben, von [@leaanthony](https://github.com/leaanthony)
- Taskfiles verschoben, umbenannt und in plattformspezifischen Verzeichnissen abgelegt, von [@leaanthony](https://github.com/leaanthony)
- Deutlich bessere Benutzerführung geschaffen, wenn `index.html` fehlt, von [@leaanthony](https://github.com/leaanthony)
- [Windows] Leistung beim Minimieren und Wiederherstellen verbessert, von [@leaanthony](https://github.com/leaanthony). Basierend auf dem ursprünglichen [PR](https://github.com/wailsapp/wails/pull/3955) von [562589540](https://github.com/562589540)
- Option `ShouldClose` entfernt (registriere stattdessen einen Hook für events.Common.WindowClosing), von [@leaanthony](https://github.com/leaanthony)
- [Windows] Flackern beim Öffnen eines Fensters reduziert, von [@leaanthony](https://github.com/leaanthony)
- `Window.Destroy` entfernt, da dies als interne Funktion vorgesehen war, von [@leaanthony](https://github.com/leaanthony)
- `WindowClose`-Ereignisse in `WindowClosing` umbenannt, von [@leaanthony](https://github.com/leaanthony)
- Frontend-Builds verwenden jetzt abhängig vom Build-Typ die Vite-Umgebung "development" oder "production", von [@leaanthony](https://github.com/leaanthony)
- Auf go-webview2 v1.19 aktualisiert, von [@leaanthony](https://github.com/leaanthony)
- Sichergestellt, dass der Taskfile-Fork verwendet wird, von @leaanthony
- Taskfile-Fork aktualisiert, um Versionsprobleme bei der Installation mit folgendem Verfahren zu beheben:
- Taskfile-Fork verwendet, um Versionsprobleme bei der Installation mit folgendem Verfahren zu beheben:
- `service.OnStartup` beendet die Anwendung jetzt bei einem Fehler und führt Folgendes aus:
- Nachrichtenübermittlung bei Systray-Klicks überarbeitet, um Benutzerinteraktionen besser abzubilden, von
- `all:frontend/dist` in eingebettete Assets aufgenommen, um Frameworks zu unterstützen, die Folgendes generieren:
- Taskfile-Überarbeitung von [leaanthony](https://github.com/leaanthony) in
- Upgrade auf `go-webview2` v1.0.16, von
- Typ `Screen` korrigiert, sodass er `ID` statt `Id` enthält, von
- Wails-Version von `go.mod.tmpl` aktualisiert, um `application.ServiceOptions` zu unterstützen, von
- Ermittlung des Dienstnamens korrigiert, von [windom](https://github.com/windom/) in
- mkdocs serve verwendet jetzt Docker, von [leaanthony](https://github.com/leaanthony)
- Entwicklungskonfiguration in `config.yml` zusammengeführt, von
- Der Systray-Dialog verwendet jetzt standardmäßig das Anwendungssymbol, sofern verfügbar (Windows), von
- Verbesserte Erfassung von GPU und Arbeitsspeicher unter macOS, von
- `WebviewGpuIsDisabled` und `EnableFraudulentWebsiteWarnings` entfernt
- Änderung der Ereignis-API: `On`/`Emit` -> Benutzerereignisse, `OnApplicationEvent` ->
- Korrektur der Ereignis-API unter Linux, von [TheGB0077](https://github.com/TheGB0077) in
- [CI] Actions verbessert und ihre Ausführung auch in Forks und
- `AbsolutePosition()` in `Position()` umbenannt, von
- Linux-WebKit-Abhängigkeit von webkitgtk2-4.0 auf webkit2gtk-4.1 aktualisiert, um
- Das gebündelte JS-Runtime-Skript ist jetzt ein ESM-Modul: Skript-Tags, die es importieren,
- Das Paket `@wailsio/runtime` veröffentlicht seine API nicht im `window.wails`
- Das Window-API-Modul `@wailsio/runtime/src/window` stellt jetzt das enthaltende Element bereit
- Die JS-Fenster-API wurde an die aktuelle Go-API `WebviewWindow` angepasst
- Der Binding-Generator verwendet jetzt standardmäßig Aufrufe per ID. Die CLI-Option `-id`
- Neue Struktur des Binding-Codes: Ausgabedateien waren zuvor in Ordnern organisiert
- Das Struct-Feld `application.Options.Bind` wurde umbenannt in
- Neue Syntax für Binding-Dienste: Dienstinstanzen müssen jetzt in Folgendes eingeschlossen werden:
- Spinner in Nicht-Terminal- oder CI-Umgebungen deaktiviert, von

### Entfernt

- **NICHT ABWÄRTSKOMPATIBEL**: `EnabledFeatures`, `DisabledFeatures` und `AdditionalLaunchArgs` aus den fensterspezifischen `WindowsWindow`-Optionen entfernt. Verwenden Sie stattdessen die anwendungsweiten Optionen `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures` und `Options.Windows.AdditionalBrowserArgs`. Diese Flags gelten global für die gemeinsam genutzte WebView2-Umgebung (#4559) von @leaanthony
- Native `IDropTarget`-Implementierung unter Windows zugunsten eines JavaScript-basierten Ansatzes entfernt (entspricht dem Verhalten von v2)
- Abhängigkeit github.com/wailsapp/mimetype zugunsten einer erweiterten Erweiterungszuordnung und der Standardbibliotheksfunktion http.DetectContentType entfernt, wodurch sich die Binärgröße um etwa 1.2 MB verringert
- Abhängigkeit gopkg.in/ini.v1 durch Implementierung eines minimalen Parsers für .desktop-Dateien im Linux-Dateiexplorer entfernt; spart etwa 45 KB
- samber/lo durch Verwendung des slices-Pakets der Go-1.21+-Standardbibliothek und minimaler interner Hilfsfunktionen aus dem Laufzeitcode entfernt; spart etwa 310 KB
- Debug-printf-Anweisungen aus dem Darwin-Handler für URL-Schemata entfernt (#4834)
- **NICHT ABWÄRTSKOMPATIBEL**: Ereignis `linux:WindowLoadChanged` entfernt. Verwenden Sie stattdessen `linux:WindowLoadFinished`, um zu erkennen, wann das Laden der WebView abgeschlossen ist (#3896) von @leaanthony

### Inkompatible Änderungen

- **Überarbeitung der Manager-API**: Die Anwendungs-API wurde von einer flachen Struktur in organisierte Manager umstrukturiert, um die Codeorganisation und Auffindbarkeit zu verbessern, von [@leaanthony](https://github.com/leaanthony) in [#4359](https://github.com/wailsapp/wails/pull/4359)
- `app.NewWebviewWindow()` → `app.Window.New()`
- `app.CurrentWindow()` → `app.Window.Current()`
- `app.GetAllWindows()` → `app.Window.GetAll()`
- `app.WindowByName()` → `app.Window.GetByName()`
- `app.EmitEvent()` → `app.Event.Emit()`
- `app.OnApplicationEvent()` → `app.Event.OnApplicationEvent()`
- `app.OnWindowEvent()` → `app.Event.OnWindowEvent()`
- `app.SetApplicationMenu()` → `app.Menu.SetApplicationMenu()`
- `app.OpenFileDialog()` → `app.Dialog.OpenFile()`
- `app.SaveFileDialog()` → `app.Dialog.SaveFile()`
- `app.MessageDialog()` → `app.Dialog.Message()`
- `app.InfoDialog()` → `app.Dialog.Info()`
- `app.WarningDialog()` → `app.Dialog.Warning()`
- `app.ErrorDialog()` → `app.Dialog.Error()`
- `app.QuestionDialog()` → `app.Dialog.Question()`
- `app.NewSystemTray()` → `app.SystemTray.New()`
- `app.GetSystemTray()` → `app.SystemTray.Get()`
- `app.ShowContextMenu()` → `app.ContextMenu.Show()`
- `app.RegisterKeybinding()` → `app.KeyBinding.Register()`
- `app.UnregisterKeybinding()` → `app.KeyBinding.Unregister()`
- `app.GetPrimaryScreen()` → `app.Screen.GetPrimary()`
- `app.GetAllScreens()` → `app.Screen.GetAll()`
- `app.BrowserOpenURL()` → `app.Browser.OpenURL()`
- `app.Environment()` → `app.Env.GetAll()`
- `app.ClipboardGetText()` → `app.Clipboard.Text()`
- `app.ClipboardSetText()` → `app.Clipboard.SetText()`
- Service-Methoden umbenannt: `Name` -> `ServiceName`, `OnStartup` -> `ServiceStartup`, `OnShutdown` -> `ServiceShutdown` von [@leaanthony](https://github.com/leaanthony)
- Methoden `Path` und `Paths` in das Paket `application` verschoben von [@leaanthony](https://github.com/leaanthony)
- Das Anwendungsmenü ist jetzt ausschließlich unter macOS verfügbar, von [@leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.78 - 2026-04-21

## Hinzugefügt

## Behoben

## v3.0.0-alpha.77 - 2026-04-18

## Behoben

## v3.0.0-alpha.76 - 2026-04-17

## Behoben

## v3.0.0-alpha.75 - 2026-04-16

## Behoben

## v3.0.0-alpha.74 - 2026-03-01

## Hinzugefügt

## Behoben

## v3.0.0-alpha.73 - 2026-02-27

## Behoben

## v3.0.0-alpha.72 - 2026-02-16

## Behoben

## v3.0.0-alpha.71 - 2026-02-10

## Hinzugefügt

## Behoben

## v3.0.0-alpha.70 - 2026-02-09

## Hinzugefügt

## Behoben

## v3.0.0-alpha.69 - 2026-02-08

## Hinzugefügt

## Behoben

## v3.0.0-alpha.68 - 2026-02-07

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.67 - 2026-02-04

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.66 - 2026-02-03

## Hinzugefügt

## Geändert

## Behoben

## Entfernt

## v3.0.0-alpha.65 - 2026-02-01

## Hinzugefügt

## v3.0.0-alpha.64 - 2026-01-26

## Hinzugefügt

## v3.0.0-alpha.63 - 2026-01-25

## Behoben

## v3.0.0-alpha.62 - 2026-01-22

## Behoben

## v3.0.0-alpha.61 - 2026-01-20

## Behoben

## v3.0.0-alpha.60 - 2026-01-14

## Behoben

## v3.0.0-alpha.59 - 2026-01-11

## Geändert

## v3.0.0-alpha.58 - 2026-01-09

## Behoben

## v3.0.0-alpha.57 - 2026-01-05

## Geändert

## Behoben

## v3.0.0-alpha.56 - 2026-01-04

## Hinzugefügt

## Geändert

## Behoben

## Entfernt

## v3.0.0-alpha.55 - 2026-01-02

## Geändert

## Behoben

## Entfernt

## v3.0.0-alpha.54 - 2025-12-29

## Hinzugefügt

## Behoben

## Entfernt

## v3.0.0-alpha.53 - 2025-12-27

## Hinzugefügt

## Behoben

## v3.0.0-alpha.52 - 2025-12-26

## Behoben

## v3.0.0-alpha.51 - 2025-12-23

## Behoben

## v3.0.0-alpha.50 - 2025-12-21

## Geändert

## v3.0.0-alpha.49 - 2025-12-18

## Geändert

## v3.0.0-alpha.48 - 2025-12-16

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.47 - 2025-12-15

## Hinzugefügt

## Behoben

## v3.0.0-alpha.46 - 2025-12-14

## Hinzugefügt

## Entfernt

## v3.0.0-alpha.45 - 2025-12-13

## Hinzugefügt

## Behoben

## v3.0.0-alpha.44 - 2025-12-12

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.43 - 2025-12-11

## Hinzugefügt

## v3.0.0-alpha.42 - 2025-12-10

## Hinzugefügt

## v3.0.0-alpha.41 - 2025-11-23

## Behoben

## v3.0.0-alpha.40 - 2025-11-13

## Behoben

## v3.0.0-alpha.39 - 2025-11-12

## Hinzugefügt

## Geändert

## v3.0.0-alpha.38 - 2025-11-04

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.37 - 2025-11-02

## Behoben

## v3.0.0-alpha.36 - 2025-10-15

## Behoben

## v3.0.0-alpha.35 - 2025-10-14

## Behoben

## v3.0.0-alpha.34 - 2025-10-06

## Hinzugefügt

## Behoben

## v3.0.0-alpha.33 - 2025-10-04

## Behoben

## v3.0.0-alpha.32 - 2025-10-02

## Behoben

## v3.0.0-alpha.31 - 2025-09-27

## Behoben

## v3.0.0-alpha.30 - 2025-09-26

## Behoben

## v3.0.0-alpha.29 - 2025-09-25

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.29 - 2025-09-25

## Hinzugefügt

## Geändert

## Behoben

## v3.0.0-alpha.27 - 2025-09-07

## Behoben

## v3.0.0-alpha.26 - 2025-08-24

## Hinzugefügt

## v3.0.0-alpha.25 - 2025-08-16

## Geändert

werden nicht mehr ignoriert und überschreiben den Konfigurationswert.

## v3.0.0-alpha.24 - 2025-08-13

## Hinzugefügt

## v3.0.0-alpha.23 - 2025-08-11

## Behoben

## v3.0.0-alpha.22 - 2025-08-10

## Hinzugefügt

## Geändert

+ Zu weit gefasste Abhängigkeiten von Linux-Paketen und veraltete RPM-Abhängigkeiten behoben.

## v3.0.0-alpha.21 - 2025-08-07

## Behoben

## v3.0.0-alpha.20 - 2025-08-06

## Behoben

## v3.0.0-alpha.19 - 2025-08-05

## Hinzugefügt

## Behoben

## v3.0.0-alpha.18 - 2025-08-03

## Hinzugefügt

## Behoben

## v3.0.0-alpha.17 - 2025-07-31

## Behoben

## v3.0.0-alpha.16 - 2025-07-25

## Hinzugefügt

## v3.0.0-alpha.15 - 2025-07-25

## Hinzugefügt

## v3.0.0-alpha.14 - 2025-07-25

## Hinzugefügt

## v3.0.0-alpha.12 - 2025-07-15

### Hinzugefügt

### Behoben

## v3.0.0-alpha.11 - 2025-07-12

## Hinzugefügt

## v3.0.0-alpha.10 - 2025-07-06

### Inkompatible Änderungen

### Hinzugefügt

### Behoben

### Geändert

## v3.0.0-alpha.9 - 2025-01-13

### Hinzugefügt

### Behoben

### Geändert

## v3.0.0-alpha.8.3 - 2024-12-07

### Geändert

## v3.0.0-alpha.8.2 - 2024-12-07

### Geändert

`go install` von @leaanthony

## v3.0.0-alpha.8.1 - 2024-12-07

### Geändert

`go install` von @leaanthony

## v3.0.0-alpha.8 - 2024-12-06

### Hinzugefügt

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### Geändert

`service.OnShutdown`für alle zuvor gestarteten Dienste von @atterpac in   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac in [#3907](https://github.com/wailsapp/wails/pull/3907)   Unterordner von @atterpac in   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) in   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) in   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (durch die Optionen `EnabledFeatures` und `DisabledFeatures` ersetzt) von   [leaanthony](https://github.com/leaanthony)

### Behoben

Kanalvariable von @michael-freling in   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) in   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   in [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) in   [#3841](https://github.com/wailsapp/wails/pull/3841)   `PasteAndMatchStyle`-Rolle im Bearbeitungsmenü unter Darwin von   [johnmccabe](https://github.com/johnmccabe) in   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) in   [#3854](https://github.com/wailsapp/wails/pull/3854) von   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   unterscheiden sich. Von @nickisworking in   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### Hinzugefügt

[mmghv](https://github.com/mmghv) in   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) und   [leaanthony](https://github.com/leaanthony) in   [#3570](https://github.com/wailsapp/wails/pull/3570)

### Geändert

Anwendungsereignisse `OnWindowEvent` -> Fensterereignisse, von   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   Branches mit dem Präfix `v3/` oder `v3-` von   [stendler](https://github.com/stendler) in   [#3747](https://github.com/wailsapp/wails/pull/3747)

### Behoben

[etesam913](https://github.com/etesam913) in   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) in   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) in   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) in   [#3614](https://github.com/wailsapp/wails/pull/3614) von   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720) von   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) von   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720) von   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) von   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746) von   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### Behoben

## v3.0.0-alpha.5 - 2024-07-30

### Hinzugefügt

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   Symbolklick von @5aaee9 in [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) in   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   in [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) in   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) in   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   basierend auf [PR](https://github.com/wailsapp/wails/pull/2044) von @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   das npm-Paket von [@fbbdev](https://github.com/fbbdev) in   [#3334](https://github.com/wailsapp/wails/pull/3334)   in [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` von [@abichinger](https://github.com/abichinger) in   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) in   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) in   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) in   [#3674](https://github.com/wailsapp/wails/pull/3674)

### Behoben

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) in   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) in   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) in   [#3477](https://github.com/wailsapp/wails/pull/3477)   von [Atterpac](https://github.com/atterpac) in   [#3320](https://github.com/wailsapp/wails/pull/3320).   in [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) in   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) in   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) in   [PR](https://github.com/wailsapp/wails/pull/3039)   installiert. Hinzugefügt von [@pylotlight](https://github.com/pylotlight) in   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   Leerzeichen – @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) in PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) in PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   angehängtes Fenster [tw1nk](https://github.com/tw1nk) in PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) in   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` ist vorhanden.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   Listener von [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) in   [#3330](https://github.com/wailsapp/wails/pull/3330)   Generator von [@fbbdev](https://github.com/fbbdev) in   [#3334](https://github.com/wailsapp/wails/pull/3334)   Generator von [@fbbdev](https://github.com/fbbdev) in   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) in   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) in   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) in   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) in   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) in   [#3431](https://github.com/wailsapp/wails/pull/3431)   von [@pekim](https://github.com/pekim) in   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) von   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) in   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows von [@bruxaodev](https://github.com/bruxaodev/) in   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Geändert

[mmghv](https://github.com/mmghv) in   [#3611](https://github.com/wailsapp/wails/pull/3611)   Unterstützung für Ubuntu 24.04 LTS von [atterpac](https://github.com/atterpac) in   [#3461](https://github.com/wailsapp/wails/pull/3461)   muss das Attribut `type="module"` besitzen. Von   [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   Objekt und startet das WML-System nicht. Dies verbessert die Kapselung. Das WML-System kann bei Bedarf manuell durch Aufruf der neuen Methode `WML.Enable` gestartet werden. Das gebündelte JS-Laufzeitskript führt weiterhin beide Vorgänge automatisch aus. Von [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   Fensterobjekt als Standardexport. Einzelne Methoden können nicht mehr über benannte ESM-Importe oder die Namespace-Importsyntax importiert werden.   API. Einige Methoden wurden umbenannt oder ihre Prototypen geändert, insbesondere: Aus `Screen` wird `GetScreen`; aus `GetZoomLevel`/`SetZoomLevel` wird `GetZoom`/`SetZoom`; `GetZoom`, `Width` und `Height` geben Werte jetzt direkt zurück, statt sie in Objekte einzuschließen. Von [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   wurde entfernt. Verwenden Sie die CLI-Option `-names`, um wieder Aufrufe anhand des Namens zu verwenden. Von [@fbbdev](https://github.com/fbbdev) in   [#3468](https://github.com/wailsapp/wails/pull/3468)   wurden nach ihrem enthaltenden Paket benannt; jetzt werden vollständige Go-Importpfade einschließlich des Modulpfads verwendet. Von [@fbbdev](https://github.com/fbbdev) in   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Von [@fbbdev](https://github.com/fbbdev) in   [#3468](https://github.com/wailsapp/wails/pull/3468)   Aufruf von `application.NewService`. Von [@fbbdev](https://github.com/fbbdev) in   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) in   [#3574](https://github.com/wailsapp/wails/pull/3574)
