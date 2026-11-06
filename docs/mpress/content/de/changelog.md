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
- `Fixed` für sämtliche Fehlerbehebungen.
- `Security` bei Sicherheitslücken.

_/

/_   * AUTOMATISCH AKTUALISIERT *   Die Versionshinweise werden aus zusammengeführten PRs generiert. Bearbeiten Sie diese Datei nicht in regulären Beitrags-PRs und erstellen Sie die entfernte Datei `v3/UNRELEASED_CHANGELOG.md` nicht neu. Informationen zur Pflege des Änderungsprotokolls finden Sie unter `AGENTS.md`. _/

## [Unveröffentlicht]

## v3.2.0-nightly.20261102 – 2026-11-02

## Geändert
- DisableOverscroll statt DisableBounce in [PR](https://github.com/wailsapp/wails/pull/6077) von @lethusang

## v3.2.0-nightly.20261101 – 2026-11-01

## Hinzugefügt
- Experimentelle CGO-freie Backends für macOS und Linux hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/6521) von @taliesin-ai
- Native YAML-Builds, portable Paketierung und Signierung in [PR](https://github.com/wailsapp/wails/pull/6531) von @taliesin-ai

## Behoben
- Verschiedene Probleme bei nächtlichen v3 Builds behoben und Abhängigkeiten bereinigt in [PR](https://github.com/wailsapp/wails/pull/6519) von @taliesin-ai
- Validierung und Fixture-Klassifizierung für v3 Releases repariert in [PR](https://github.com/wailsapp/wails/pull/6533) von @taliesin-ai
- Plattformübergreifende Fehler bei der Release-Validierung behoben in [PR](https://github.com/wailsapp/wails/pull/6538) von @taliesin-ai

## v3.1.0 – 2026-11-01

## Hinzugefügt
- Experimentelle CGO-freie Backends für macOS und Linux hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/6521) von @taliesin-ai
- Native YAML-Builds, portable Paketierung und Signierung in [PR](https://github.com/wailsapp/wails/pull/6531) von @taliesin-ai

## Geändert
- Verworfene Desktop-Metadaten werden nicht mehr kopiert in [PR](https://github.com/wailsapp/wails/pull/6391) von @taliesin-ai
- Inhaltslängen ohne fmt formatiert in [PR](https://github.com/wailsapp/wails/pull/6393) von @taliesin-ai
- Indexwurzeln ausgewählt, ohne Kandidaten zu sammeln, in [PR](https://github.com/wailsapp/wails/pull/6404) von @taliesin-ai
- Kopierpuffer für nicht durchsuchbare Dateien wiederverwendet in [PR](https://github.com/wailsapp/wails/pull/6408) von @taliesin-ai
- Allokation eines Readers zur Erkennung von Präfixen vermieden in [PR](https://github.com/wailsapp/wails/pull/6410) von @taliesin-ai
- Lesevorgänge zur MIME-Erkennung bei bekannten Erweiterungen übersprungen in [PR](https://github.com/wailsapp/wails/pull/6414) von @taliesin-ai
- Anwachsen der Warteschlange und Oszillation beim Verkleinern verhindert in [PR](https://github.com/wailsapp/wails/pull/6415) von @taliesin-ai
- Allokation verworfener Anfragekontexte vermieden in [PR](https://github.com/wailsapp/wails/pull/6419) von @taliesin-ai
- Vorbereitung deaktivierter Debug-Protokolle übersprungen in [PR](https://github.com/wailsapp/wails/pull/6420) von @taliesin-ai
- Puffer zur MIME-Erkennung verzögert alloziert in [PR](https://github.com/wailsapp/wails/pull/6428) von @taliesin-ai
- Allokationen bei der Hyperlink-Formatierung reduziert in [PR](https://github.com/wailsapp/wails/pull/6443) von @taliesin-ai
- Temporäre Allokationen für die Aufzählungsformatierung entfernt in [PR](https://github.com/wailsapp/wails/pull/6449) von @taliesin-ai
- UUID-Zeichenfolgenkonvertierung auf eine Allokation reduziert in [PR](https://github.com/wailsapp/wails/pull/6453) von @taliesin-ai
- Allokationen bei der Formatierung von Statement-Handles reduziert in [PR](https://github.com/wailsapp/wails/pull/6454) von @taliesin-ai
- Temporäre Digest-Allokation für namensbasierte UUIDs beseitigt in [PR](https://github.com/wailsapp/wails/pull/6455) von @taliesin-ai
- Kopieren unvollständiger Zeilen der Befehlsausgabe vermieden in [PR](https://github.com/wailsapp/wails/pull/6460) von @taliesin-ai
- Ausführbare Dateien während des Paket-Stagings als Stream kopiert in [PR](https://github.com/wailsapp/wails/pull/6472) von @taliesin-ai
- Erneutes Scannen unvollständiger Ausgabezeilen vermieden in [PR](https://github.com/wailsapp/wails/pull/6473) von @taliesin-ai
- Kompilierte Muster für das Befehlsrouting wiederverwendet in [PR](https://github.com/wailsapp/wails/pull/6475) von @taliesin-ai
- Metadatenabfragen bei der Runtime-Erkennung vermieden in [PR](https://github.com/wailsapp/wails/pull/6482) von @taliesin-ai
- Allokationen lokalisierter Fallback-Seiten reduziert in [PR](https://github.com/wailsapp/wails/pull/6485) von @taliesin-ai
- Erneutes Kompilieren von Mustern für jeden Prozessor vermieden in [PR](https://github.com/wailsapp/wails/pull/6490) von @taliesin-ai

## Behoben
- SetTooltip für den System-Tray implementiert in [PR](https://github.com/wailsapp/wails/pull/6017) von @fan711
- Dark Mode aus dem Einstellungsportal ausgelesen in [PR](https://github.com/wailsapp/wails/pull/6072) von @TechHutTV
- Window.SetZoom kann jetzt Werte unter 1.0 verwenden in [PR](https://github.com/wailsapp/wails/pull/6166) von @APshenkin
- Ereignisse an jeden Go-Listener in Emit-Reihenfolge zugestellt in [PR](https://github.com/wailsapp/wails/pull/6275) von @DevLumuz
- Event-Hooks in Emit unter hookLock gelesen in [PR](https://github.com/wailsapp/wails/pull/6276) von @DevLumuz
- Die Tray-Eigenschaft Menu wird immer bekannt gegeben, damit ein nach Run erstellter Tray nicht verworfen wird, in [PR](https://github.com/wailsapp/wails/pull/6284) von @AlbinoGeek
- Vor Run gesendete Arbeit für den Haupt-Thread wird bis zum Start von Run zurückgehalten in [PR](https://github.com/wailsapp/wails/pull/6287) von @AlbinoGeek
- Korrekte Asset-URL im GitHub-Updater-Provider verwendet in [PR](https://github.com/wailsapp/wails/pull/6304) von @gavidroselj
- Anfragetext beibehalten, wenn WebView2 Content-Length auslässt, in [PR](https://github.com/wailsapp/wails/pull/6326) von @strraaa
- macOS- und GTK-Bildschirme im DIP-Raum gemeldet (Begrenzungen bei gemischten DPI-Werten) in [PR](https://github.com/wailsapp/wails/pull/6346) von @julianstorer
- Eine verworfene Benachrichtigung wird nicht als Klick gemeldet in [PR](https://github.com/wailsapp/wails/pull/6367) von @bercedev
- Entwicklungsprozesse beim Schließen des Terminals beendet in [PR](https://github.com/wailsapp/wails/pull/6371) von @xmasdev
- Changelog-Läufe bei der Evidenzprüfung Squash-Merge-Commits zugeordnet in [PR](https://github.com/wailsapp/wails/pull/6373) von @taliesin-ai
- Entprellung nach einem Neustart der Taskleiste eine Sekunde lang aufrechterhalten in [PR](https://github.com/wailsapp/wails/pull/6386) von @taliesin-ai
- Größenbeschränkungen über einen dort beginnenden Vollbildmodus hinweg beibehalten in [PR](https://github.com/wailsapp/wails/pull/6387) von @taliesin-ai
- Ein Vollbild- oder maximiertes Fenster wird nicht an eine neue Beschränkung angepasst in [PR](https://github.com/wailsapp/wails/pull/6388) von @taliesin-ai
- im Liquid-Glass-Beispiel vor Run keinen Dialog anzeigen in [PR](https://github.com/wailsapp/wails/pull/6389) von @taliesin-ai
- ein vor Run angeforderter Dialog bleibt wirkungslos, statt endlos zu warten, in [PR](https://github.com/wailsapp/wails/pull/6390) von @taliesin-ai
- Fallback verwenden, wenn der Desktop-Eintrag keinen Befehl enthält, in [PR](https://github.com/wailsapp/wails/pull/6392) von @taliesin-ai
- aktuellen Inhalt ohne Beibehaltung von Dateinamen erkennen in [PR](https://github.com/wailsapp/wails/pull/6394) von @taliesin-ai
- Asset-Erweiterungen mit gemischter Groß-/Kleinschreibung erkennen in [PR](https://github.com/wailsapp/wails/pull/6395) von @taliesin-ai
- verbrauchte Puffer der Inhaltserkennung freigeben in [PR](https://github.com/wailsapp/wails/pull/6396) von @taliesin-ai
- entfernte ringqueue-Referenzen freigeben in [PR](https://github.com/wailsapp/wails/pull/6397) von @taliesin-ai
- gelöschte Header vor der Fallback-Auswahl berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6398) von @taliesin-ai
- Fallback-Antwort-Header vor dem Flush festschreiben in [PR](https://github.com/wailsapp/wails/pull/6399) von @taliesin-ai
- nur exakt übereinstimmende Indexdateien finden in [PR](https://github.com/wailsapp/wails/pull/6400) von @taliesin-ai
- bei nicht durchsuchbaren HEAD-Antworten Bodys und vollständiges Lesen vermeiden in [PR](https://github.com/wailsapp/wails/pull/6401) von @taliesin-ai
- Präfix für die MIME-Erkennung über kurze Lesevorgänge hinweg füllen in [PR](https://github.com/wailsapp/wails/pull/6402) von @taliesin-ai
- erfolgreiche leere, nicht durchsuchbare Antworten festschreiben in [PR](https://github.com/wailsapp/wails/pull/6403) von @taliesin-ai
- Duplizierung expandierter URI-Ziele vermeiden in [PR](https://github.com/wailsapp/wails/pull/6405) von @taliesin-ai
- Pfade bei der Feldexpansion beibehalten in [PR](https://github.com/wailsapp/wails/pull/6406) von @taliesin-ai
- XDG-Suchverzeichnisse für Desktop-Dateien berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6407) von @taliesin-ai
- Header der Inhaltserkennung bei leerem Flush festschreiben in [PR](https://github.com/wailsapp/wails/pull/6409) von @taliesin-ai
- Verzeichnisse überspringen, die Desktop-Einträge verdecken, in [PR](https://github.com/wailsapp/wails/pull/6411) von @taliesin-ai
- in Pfaden ausführbarer Dateien gefundene Ziele beibehalten in [PR](https://github.com/wailsapp/wails/pull/6412) von @taliesin-ai
- Callback-Captures nach Auslösung des Timers freigeben in [PR](https://github.com/wailsapp/wails/pull/6413) von @taliesin-ai
- falsche ldconfig-Bibliothekstreffer vermeiden in [PR](https://github.com/wailsapp/wails/pull/6416) von @taliesin-ai
- Fallback-Treffer für Bibliotheken im Dateisystem einschränken in [PR](https://github.com/wailsapp/wails/pull/6417) von @taliesin-ai
- Datenrennen bei der Cache-Invalidierung beseitigen in [PR](https://github.com/wailsapp/wails/pull/6418) von @taliesin-ai
- optionale Logger-Konfiguration berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6421) von @taliesin-ai
- gemeinsame Mutation der Start-URL und veraltete Ports vermeiden in [PR](https://github.com/wailsapp/wails/pull/6422) von @taliesin-ai
- Wiederholungsversuche des Entwicklungsproxys bei Abbruch der Anfrage beenden in [PR](https://github.com/wailsapp/wails/pull/6423) von @taliesin-ai
- bei fehlenden gebündelten Assets 404 zurückgeben in [PR](https://github.com/wailsapp/wails/pull/6424) von @taliesin-ai
- HEAD-Anfragen für die gebündelte Runtime verarbeiten in [PR](https://github.com/wailsapp/wails/pull/6425) von @taliesin-ai
- maskierte Pfade beim Routing von Diensten beibehalten in [PR](https://github.com/wailsapp/wails/pull/6426) von @taliesin-ai
- übergroße numerische Vorabversionsbezeichner korrekt sortieren in [PR](https://github.com/wailsapp/wails/pull/6427) von @taliesin-ai
- HTTP-Fehler vor dem Speichern von Tooling-Downloads zurückweisen in [PR](https://github.com/wailsapp/wails/pull/6429) von @taliesin-ai
- Startfrist während der Desktop-Erkennung berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6430) von @taliesin-ai
- Anfrage-Bodys bei Wiederholungsversuchen des Entwicklungsproxys beibehalten in [PR](https://github.com/wailsapp/wails/pull/6432) von @taliesin-ai
- HTTP-Antwort-Bodys in jedem Rückgabepfad schließen in [PR](https://github.com/wailsapp/wails/pull/6434) von @taliesin-ai
- Abbruchfehler bei parallelen Bibliothekssuchen beibehalten in [PR](https://github.com/wailsapp/wails/pull/6435) von @taliesin-ai
- fehlgeschlagene Antworten auf GitHub-Tag-Anfragen vor der Dekodierung zurückweisen in [PR](https://github.com/wailsapp/wails/pull/6436) von @taliesin-ai
- literale Vorlagenpfade mit eckigen Klammern beibehalten in [PR](https://github.com/wailsapp/wails/pull/6437) von @taliesin-ai
- wails3 einmal installieren und Installationsfehler beibehalten in [PR](https://github.com/wailsapp/wails/pull/6438) von @taliesin-ai
- verborgene Einträge und Header in Kommentarblöcken ignorieren in [PR](https://github.com/wailsapp/wails/pull/6439) von @taliesin-ai
- Fehler bei der Suche nach stabilen Releases ausgeben in [PR](https://github.com/wailsapp/wails/pull/6440) von @taliesin-ai
- reine Textnachrichten ohne Formatinterpretation ausgeben in [PR](https://github.com/wailsapp/wails/pull/6441) von @taliesin-ai
- Schreib-Panics nach dem Laden von JSON-null vermeiden in [PR](https://github.com/wailsapp/wails/pull/6442) von @taliesin-ai
- bei Lesevorgängen des gesamten Speichers Snapshots zurückgeben in [PR](https://github.com/wailsapp/wails/pull/6444) von @taliesin-ai
- Fehler nach teilweiser Zeileniteration weitergeben in [PR](https://github.com/wailsapp/wails/pull/6446) von @taliesin-ai
- Zielpfade für die Extraktion literal beibehalten in [PR](https://github.com/wailsapp/wails/pull/6447) von @taliesin-ai
- Handles vorbereiteter Anweisungen beim Schließen und erneuten Öffnen freigeben in [PR](https://github.com/wailsapp/wails/pull/6448) von @taliesin-ai
- Datenrennen zwischen Vorbereitung und Schließen oder erneutem Öffnen behandeln in [PR](https://github.com/wailsapp/wails/pull/6450) von @taliesin-ai
- HEAD-Anfragen für die Fallback-Indexseite berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6451) von @taliesin-ai
- vollständige hexadezimale Anweisungs-Handles verlangen in [PR](https://github.com/wailsapp/wails/pull/6452) von @taliesin-ai
- standardmäßige In-Memory-Datenbank bei Nebenläufigkeit beibehalten in [PR](https://github.com/wailsapp/wails/pull/6456) von @taliesin-ai
- stderr von der geparsten HEAD-Ausgabe trennen in [PR](https://github.com/wailsapp/wails/pull/6457) von @taliesin-ai
- Runtime-Assets aus eingebetteten Bundles erzeugen in [PR](https://github.com/wailsapp/wails/pull/6458) von @taliesin-ai
- TODO-, FIXME- und TBD-Platzhalter erkennen in [PR](https://github.com/wailsapp/wails/pull/6459) von @taliesin-ai
- API-Antworten für Release Notes validieren in [PR](https://github.com/wailsapp/wails/pull/6461) von @taliesin-ai
- Task-Verzeichnisse bei Statusprüfungen berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6462) von @taliesin-ai
- Datenverlust beim Kopieren einer Datei auf sich selbst verhindern in [PR](https://github.com/wailsapp/wails/pull/6463) von @taliesin-ai
- Task-Verzeichnis in die Cache-Identität einbeziehen in [PR](https://github.com/wailsapp/wails/pull/6464) von @taliesin-ai
- Ports einmal pro Prüfung testen in [PR](https://github.com/wailsapp/wails/pull/6465) von @taliesin-ai
- Panic vermeiden, wenn Cache-Einträge als null dekodiert werden, in [PR](https://github.com/wailsapp/wails/pull/6466) von @taliesin-ai
- numerische Vorabversionen ohne Ganzzahlgrenzen erhöhen in [PR](https://github.com/wailsapp/wails/pull/6467) von @taliesin-ai
- verhindern, dass Legacy-Prüfungen Cache-Fehlschläge verbergen, in [PR](https://github.com/wailsapp/wails/pull/6468) von @taliesin-ai
- Standard-Tasks im Basisverzeichnis des Executors ausführen in [PR](https://github.com/wailsapp/wails/pull/6469) von @taliesin-ai
- literale Prozentzeichen in Explorer-Befehlen beibehalten in [PR](https://github.com/wailsapp/wails/pull/6470) von @taliesin-ai
- Verzeichnisse für Vorbedingungen und Task-Befehle in [PR](https://github.com/wailsapp/wails/pull/6471) von @taliesin-ai angleichen
- Symlinks durch generierte reguläre Dateien in [PR](https://github.com/wailsapp/wails/pull/6474) von @taliesin-ai ersetzen
- Vorhandene DMGs bei fehlgeschlagenen Builds in [PR](https://github.com/wailsapp/wails/pull/6476) von @taliesin-ai beibehalten
- Backslashes in maskierten regulären Ausdrücken und Pfaden in [PR](https://github.com/wailsapp/wails/pull/6477) von @taliesin-ai beibehalten
- Größen generierter Symbole vor der Bildzuweisung in [PR](https://github.com/wailsapp/wails/pull/6478) von @taliesin-ai validieren
- Befehlsketten und Umleitungen in [PR](https://github.com/wailsapp/wails/pull/6479) von @taliesin-ai über die Shell ausführen
- Vorherige Linux-Pakete bei fehlgeschlagenen Builds in [PR](https://github.com/wailsapp/wails/pull/6480) von @taliesin-ai beibehalten
- Plattformspezifische Shell-Ausführung als Fallback in [PR](https://github.com/wailsapp/wails/pull/6481) von @taliesin-ai verwenden
- Assets bei mehrfachen führenden Schrägstrichen im Pfad in [PR](https://github.com/wailsapp/wails/pull/6483) von @taliesin-ai ausliefern
- Stabile Versionen mit Build-Metadaten in [PR](https://github.com/wailsapp/wails/pull/6484) von @taliesin-ai erkennen
- Fehler bei der Ressourcengenerierung in [PR](https://github.com/wailsapp/wails/pull/6486) von @taliesin-ai korrekt melden
- Race Conditions beim Ersetzen des aktiven Reporters in [PR](https://github.com/wailsapp/wails/pull/6487) von @taliesin-ai verhindern
- Gleichzeitiges Speichern von Cache-Dateien in [PR](https://github.com/wailsapp/wails/pull/6488) von @taliesin-ai serialisieren
- Mehrere Gleichheitskandidaten für Vorlagen in [PR](https://github.com/wailsapp/wails/pull/6489) von @taliesin-ai unterstützen
- Gespeicherte Kategorie-Maps mit dem Wert null in [PR](https://github.com/wailsapp/wails/pull/6491) von @taliesin-ai normalisieren
- Diagnoseausgaben fehlgeschlagener variabler Shell-Befehle in [PR](https://github.com/wailsapp/wails/pull/6492) von @taliesin-ai anzeigen
- Shell-Befehlsketten in [PR](https://github.com/wailsapp/wails/pull/6493) von @taliesin-ai vom impliziten Go-Cache ausschließen
- Persistierung von Linux-Kategorien in [PR](https://github.com/wailsapp/wails/pull/6494) von @taliesin-ai serialisieren
- Go-Ausgaben mit maskierten oder durch Escapezeichen geschützten Leerzeichen in [PR](https://github.com/wailsapp/wails/pull/6495) von @taliesin-ai zwischenspeichern
- Mac-Symbole bei fehlgeschlagener Codierung in [PR](https://github.com/wailsapp/wails/pull/6496) von @taliesin-ai beibehalten
- Ungültige Go-Test-Flags für die Anzahl in [PR](https://github.com/wailsapp/wails/pull/6497) von @taliesin-ai nicht mehr verwerfen
- Quellberechtigungen bei neuen Dateikopien in [PR](https://github.com/wailsapp/wails/pull/6498) von @taliesin-ai beibehalten
- Verschiedene Probleme mit v3-Nightly-Builds beheben und Abhängigkeiten in [PR](https://github.com/wailsapp/wails/pull/6519) von @taliesin-ai bereinigen
- Validierung von v3-Releases und Klassifizierung von Test-Fixtures in [PR](https://github.com/wailsapp/wails/pull/6533) von @taliesin-ai korrigieren
- Plattformübergreifende Fehler bei der Release-Validierung in [PR](https://github.com/wailsapp/wails/pull/6538) von @taliesin-ai beheben

## v3.1.0-nightly.20261028 - 2026-10-28

## Geändert
- Changelog-Generierung durch Wiederverwendung kompilierter Regex-Muster in [PR](https://github.com/wailsapp/wails/pull/6490) von @taliesin-ai beschleunigen
- Allokationen bei der Auslieferung lokalisierter Fallback-Seiten in [PR](https://github.com/wailsapp/wails/pull/6485) von @taliesin-ai reduzieren
- Ermittlung der AppImage-Laufzeit durch Vermeidung unnötiger Metadatenabfragen in [PR](https://github.com/wailsapp/wails/pull/6482) von @taliesin-ai beschleunigen
- Kompilierte Muster für die Befehlsweiterleitung wiederverwenden, um die Leistung von wake in [PR](https://github.com/wailsapp/wails/pull/6475) von @taliesin-ai zu verbessern
- Leistung durch Vermeidung des erneuten Durchsuchens unvollständiger Ausgabezeilen in [PR](https://github.com/wailsapp/wails/pull/6473) von @taliesin-ai verbessern
- Kopien ausführbarer Dateien beim Staging von MSIX-Paketen streamen, um den Speicherverbrauch zu senken und Builds in [PR](https://github.com/wailsapp/wails/pull/6472) von @taliesin-ai zu beschleunigen
- Kopieren unvollständiger Befehlsausgabezeilen vermeiden, um die Leistung in [PR](https://github.com/wailsapp/wails/pull/6460) von @taliesin-ai zu verbessern
- Allokationen beim Generieren namensbasierter UUIDs in [PR](https://github.com/wailsapp/wails/pull/6455) von @taliesin-ai reduzieren
- Allokationen beim Formatieren von SQLite-Anweisungs-Handles in [PR](https://github.com/wailsapp/wails/pull/6454) von @taliesin-ai reduzieren
- UUID-Zeichenkettenkonvertierung in [PR](https://github.com/wailsapp/wails/pull/6453) von @taliesin-ai auf eine einzige Allokation reduzieren
- Allokationen beim Formatieren von Changelog-Aufzählungseinträgen in [PR](https://github.com/wailsapp/wails/pull/6449) von @taliesin-ai reduzieren
- Allokationen beim Formatieren von Terminal-Hyperlinks in [PR](https://github.com/wailsapp/wails/pull/6443) von @taliesin-ai reduzieren
- Allokationen bei der Asset-Auslieferung durch verzögerte MIME-Puffererkennung in [PR](https://github.com/wailsapp/wails/pull/6428) von @taliesin-ai reduzieren
- Vorbereitung von Debug-Protokollen im Asset-Server überspringen, wenn die Debug-Protokollierung in [PR](https://github.com/wailsapp/wails/pull/6420) von @taliesin-ai deaktiviert ist
- Allokationen bei der Verarbeitung von WebView-Anfragen durch Vermeidung verworfener Anfragekontexte in [PR](https://github.com/wailsapp/wails/pull/6419) von @taliesin-ai reduzieren
- Asset-Auslieferung beschleunigen, indem die MIME-Erkennung bei bekannten Dateierweiterungen in [PR](https://github.com/wailsapp/wails/pull/6414) von @taliesin-ai übersprungen wird
- Allokationen bei der Asset-Auslieferung reduzieren, indem in [PR](https://github.com/wailsapp/wails/pull/6410) von @taliesin-ai kein Reader zur Inhaltserkennung verwendet wird
- Allokationen bei der Auslieferung nicht durchsuchbarer Asset-Dateien in [PR](https://github.com/wailsapp/wails/pull/6408) von @taliesin-ai reduzieren
- Auswahl des Index-Stammverzeichnisses des Asset-Servers beschleunigen, indem in [PR](https://github.com/wailsapp/wails/pull/6404) von @taliesin-ai keine Kandidaten gesammelt werden
- Leistung des Asset-Servers verbessern, indem Inhaltslängen in [PR](https://github.com/wailsapp/wails/pull/6393) von @taliesin-ai ohne fmt formatiert werden
- Speicherkopien beim Lesen von Metadaten des Desktop-Datei-Explorers in [PR](https://github.com/wailsapp/wails/pull/6391) von @taliesin-ai reduzieren

## Behoben
- Berechtigungen der Quelldatei beim Kopieren in neue Zieldateien in [PR](https://github.com/wailsapp/wails/pull/6498) von @taliesin-ai beibehalten
- Verhindern, dass Vollbild- und maximierte Fenster ihre Größe ändern, wenn sich Mindest- oder Höchstgrößenbeschränkungen in [PR](https://github.com/wailsapp/wails/pull/6388) von @taliesin-ai ändern
- Geplante Aktualisierungen des Sternverlaufs für das Repository wailsapp/wails in [PR](https://github.com/wailsapp/wails/pull/6500) von @taliesin-ai veröffentlichen
- Ungültige Go-Test-Flags für die Anzahl in wake in [PR](https://github.com/wailsapp/wails/pull/6497) von @taliesin-ai nicht mehr verwerfen
- Mac-Symbole bei fehlgeschlagener Codierung in [PR](https://github.com/wailsapp/wails/pull/6496) von @taliesin-ai beibehalten
- Build-Cache-Fehlschläge bei Go-Ausgaben mit maskierten oder durch Escapezeichen geschützten Leerzeichen in [PR](https://github.com/wailsapp/wails/pull/6495) von @taliesin-ai beheben
- Persistierung von Benachrichtigungskategorien unter Linux serialisieren, um Race Conditions bei gleichzeitigen Schreibvorgängen in [PR](https://github.com/wailsapp/wails/pull/6494) von @taliesin-ai zu verhindern
- Shell-Befehlsketten in [PR](https://github.com/wailsapp/wails/pull/6493) von @taliesin-ai vom impliziten Go-Build-Cache ausschließen
- Fehler beim Laden von Maps für Benachrichtigungskategorien beheben, wenn die gespeicherten Daten in [PR](https://github.com/wailsapp/wails/pull/6491) von @taliesin-ai null sind
- Diagnoseausgaben fehlgeschlagener variabler Shell-Befehle in [PR](https://github.com/wailsapp/wails/pull/6492) von @taliesin-ai anzeigen
- Vorlagenabgleich in wake bei mehreren Gleichheitskandidaten in [PR](https://github.com/wailsapp/wails/pull/6489) von @taliesin-ai korrigieren
- Gleichzeitige Speicherungen von Cache-Dateien serialisieren, um beschädigte Cache-Schreibvorgänge zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/6488) von @taliesin-ai
- Race Conditions beim Ersetzen des aktiven Reporters verhindern in [PR](https://github.com/wailsapp/wails/pull/6487) von @taliesin-ai
- Fehler bei der Ressourcengenerierung in der CLI korrekt melden in [PR](https://github.com/wailsapp/wails/pull/6486) von @taliesin-ai
- Stabile Versionen mit Build-Metadaten in der CLI erkennen in [PR](https://github.com/wailsapp/wails/pull/6484) von @taliesin-ai
- Assets ausliefern, deren Anfragepfade wiederholte führende Schrägstriche enthalten, in [PR](https://github.com/wailsapp/wails/pull/6483) von @taliesin-ai
- Shell-Fallback-Ausführung für wake auf jeder Plattform korrigieren in [PR](https://github.com/wailsapp/wails/pull/6481) von @taliesin-ai
- Vorhandene Linux-Pakete bei einem fehlgeschlagenen Build beibehalten in [PR](https://github.com/wailsapp/wails/pull/6480) von @taliesin-ai
- Befehlsketten und Umleitungen über die Shell ausführen in [PR](https://github.com/wailsapp/wails/pull/6479) von @taliesin-ai
- Maskierte Backslashes in regulären Ausdrücken und Pfaden beim Parsen von wake-Befehlen beibehalten in [PR](https://github.com/wailsapp/wails/pull/6477) von @taliesin-ai
- Vorhandene DMG-Dateien beibehalten, wenn das Paketieren von Builds fehlschlägt, in [PR](https://github.com/wailsapp/wails/pull/6476) von @taliesin-ai
- Symlinks in der Bindings-Ausgabe durch generierte reguläre Dateien ersetzen in [PR](https://github.com/wailsapp/wails/pull/6474) von @taliesin-ai
- Wörtliche Prozentzeichen in Linux-Explorer-Befehlen beibehalten in [PR](https://github.com/wailsapp/wails/pull/6470) von @taliesin-ai
- Verzeichnisse für Vorbedingungen und Task-Befehle angleichen, damit sie konsistent aufgelöst werden, in [PR](https://github.com/wailsapp/wails/pull/6471) von @taliesin-ai
- Standard-Tasks im Basisverzeichnis des Executors ausführen in [PR](https://github.com/wailsapp/wails/pull/6469) von @taliesin-ai
- Verhindern, dass Legacy-Prüfungen Cache-Fehlschläge verbergen, in [PR](https://github.com/wailsapp/wails/pull/6468) von @taliesin-ai
- Nummerierung von Vorabversionen in der CLI korrigieren, sodass sie ohne Begrenzung durch Ganzzahlüberläufe inkrementiert wird, in [PR](https://github.com/wailsapp/wails/pull/6467) von @taliesin-ai
- Panic verhindern, wenn Cache-Einträge als null dekodiert werden, in [PR](https://github.com/wailsapp/wails/pull/6466) von @taliesin-ai
- Ports in der CLI pro Prüfung nur einmal sondieren in [PR](https://github.com/wailsapp/wails/pull/6465) von @taliesin-ai
- Das Task-Verzeichnis in die Cache-Identität aufnehmen, um veraltete wake-Ergebnisse zu vermeiden, in [PR](https://github.com/wailsapp/wails/pull/6464) von @taliesin-ai
- Datenverlust verhindern, wenn beim Paketieren eine Datei auf sich selbst kopiert wird, in [PR](https://github.com/wailsapp/wails/pull/6463) von @taliesin-ai
- Task-Verzeichnisse bei Statusprüfungen berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6462) von @taliesin-ai
- Die Platzhalter TODO, FIXME und TBD in Changelog-Einträgen erkennen in [PR](https://github.com/wailsapp/wails/pull/6459) von @taliesin-ai
- Runtime-Assets in der CLI aus eingebetteten Bundles generieren in [PR](https://github.com/wailsapp/wails/pull/6458) von @taliesin-ai
- stderr in Git-Befehlen von der geparsten HEAD-Ausgabe trennen in [PR](https://github.com/wailsapp/wails/pull/6457) von @taliesin-ai
- Die standardmäßige SQLite-In-Memory-Datenbank bei gleichzeitigem Zugriff erhalten in [PR](https://github.com/wailsapp/wails/pull/6456) von @taliesin-ai
- Korrekte Antworten auf HEAD-Anfragen für die Fallback-Indexseite zurückgeben in [PR](https://github.com/wailsapp/wails/pull/6451) von @taliesin-ai
- Race Conditions zwischen dem Vorbereiten von SQLite-Anweisungen und dem Schließen oder erneuten Öffnen der Datenbank behandeln in [PR](https://github.com/wailsapp/wails/pull/6450) von @taliesin-ai
- Handles vorbereiteter Anweisungen beim Schließen oder erneuten Öffnen der SQLite-Datenbank freigeben in [PR](https://github.com/wailsapp/wails/pull/6448) von @taliesin-ai
- Zielpfade für die Extraktion bei Verwendung von gosod wörtlich beibehalten in [PR](https://github.com/wailsapp/wails/pull/6447) von @taliesin-ai
- Fehler aus SQLite-Abfragen nach teilweiser Zeileniteration weitergeben in [PR](https://github.com/wailsapp/wails/pull/6446) von @taliesin-ai
- Bei Lesevorgängen des gesamten Speichers in kvstore Snapshots zurückgeben in [PR](https://github.com/wailsapp/wails/pull/6444) von @taliesin-ai
- Schreib-Panics in kvstore nach dem Laden eines JSON-null-Werts verhindern in [PR](https://github.com/wailsapp/wails/pull/6442) von @taliesin-ai
- Unformatierte Terminalmeldungen ausgeben, ohne Formatierungsverben zu interpretieren, in [PR](https://github.com/wailsapp/wails/pull/6441) von @taliesin-ai
- Fehler bei der Ermittlung des neuesten stabilen Releases im Updater melden in [PR](https://github.com/wailsapp/wails/pull/6440) von @taliesin-ai
- Ausgeblendete Einträge und Überschriften in Changelog-Kommentarblöcken ignorieren in [PR](https://github.com/wailsapp/wails/pull/6439) von @taliesin-ai
- wails3 nur einmal installieren und Installationsfehler im Updater melden in [PR](https://github.com/wailsapp/wails/pull/6438) von @taliesin-ai
- Wörtliche, in eckige Klammern eingeschlossene Vorlagenpfade in gosod beibehalten in [PR](https://github.com/wailsapp/wails/pull/6437) von @taliesin-ai
- Fehlgeschlagene GitHub-Tag-Antworten im Updater vor dem Dekodieren ablehnen in [PR](https://github.com/wailsapp/wails/pull/6436) von @taliesin-ai
- Abbruchfehler bei parallelen Bibliothekssuchen beibehalten in [PR](https://github.com/wailsapp/wails/pull/6435) von @taliesin-ai
- HTTP-Antwortkörper auf jedem Rückgabepfad schließen, um Ressourcenlecks zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/6434) von @taliesin-ai
- Anfragekörper beibehalten, wenn der Entwicklungsproxy Anfragen erneut versucht, in [PR](https://github.com/wailsapp/wails/pull/6432) von @taliesin-ai
- Startfrist bei der Ermittlung des Desktop-Datei-Explorers einhalten in [PR](https://github.com/wailsapp/wails/pull/6430) von @taliesin-ai
- HTTP-Fehler vor dem Speichern heruntergeladener AppImage-Werkzeuge ablehnen in [PR](https://github.com/wailsapp/wails/pull/6429) von @taliesin-ai
- Übergroße numerische Bezeichner von Vorabversionen bei semver-Vergleichen korrekt sortieren in [PR](https://github.com/wailsapp/wails/pull/6427) von @taliesin-ai
- Maskierte Pfade beim Routing von Diensten im Asset-Server beibehalten in [PR](https://github.com/wailsapp/wails/pull/6426) von @taliesin-ai
- HEAD-Anfragen für die gebündelte Runtime im Asset-Server verarbeiten in [PR](https://github.com/wailsapp/wails/pull/6425) von @taliesin-ai
- Für fehlende gebündelte Assets 404 statt einer falschen Antwort zurückgeben in [PR](https://github.com/wailsapp/wails/pull/6424) von @taliesin-ai
- Wiederholungsversuche des Entwicklungsproxys beenden, wenn eine Anfrage abgebrochen wird, in [PR](https://github.com/wailsapp/wails/pull/6423) von @taliesin-ai
- Gemeinsame Mutation der Start-URL und veraltete Ports im Asset-Server vermeiden in [PR](https://github.com/wailsapp/wails/pull/6422) von @taliesin-ai
- Optionale Logger-Konfiguration im Asset-Server berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6421) von @taliesin-ai
- Data Race bei der Invalidierung des Bibliothekspfad-Caches beseitigen in [PR](https://github.com/wailsapp/wails/pull/6418) von @taliesin-ai
- Falsche Bibliothekstreffer bei der Suche in der ldconfig-Ausgabe vermeiden in [PR](https://github.com/wailsapp/wails/pull/6416) von @taliesin-ai
- Wachstum und Schrumpfoszillation der Asset-Server-Warteschlange unter Last verhindern in [PR](https://github.com/wailsapp/wails/pull/6415) von @taliesin-ai
- Callback-Captures nach Auslösen des Debounce-Timers freigeben, damit veraltete Referenzen nicht erhalten bleiben, in [PR](https://github.com/wailsapp/wails/pull/6413) von @taliesin-ai
- In Pfaden zu ausführbaren Dateien gefundene Datei-Explorer-Ziele beibehalten in [PR](https://github.com/wailsapp/wails/pull/6412) von @taliesin-ai
- Verzeichnisse überspringen, die Desktop-Einträge im Datei-Explorer verdecken, in [PR](https://github.com/wailsapp/wails/pull/6411) von @taliesin-ai
- Korrigieren, dass Asset-Server-Header bei leerem Antwortkörper nicht gesendet werden, in [PR](https://github.com/wailsapp/wails/pull/6409) von @taliesin-ai
- XDG-Suchverzeichnisse für Desktop-Dateien beim Ermitteln von Anwendungssymbolen berücksichtigen in [PR](https://github.com/wailsapp/wails/pull/6407) von @taliesin-ai
- Dateipfade bei der Felderweiterung im Datei-Explorer beibehalten, in [PR](https://github.com/wailsapp/wails/pull/6406) von @taliesin-ai
- Doppelte erweiterte URI-Ziele im Datei-Explorer vermeiden, in [PR](https://github.com/wailsapp/wails/pull/6405) von @taliesin-ai
- Erfolgreiche leere, nicht durchsuchbare Asset-Antworten abschließen, statt sie ausstehend zu lassen, in [PR](https://github.com/wailsapp/wails/pull/6403) von @taliesin-ai
- Präfix für die MIME-Erkennung im Asset-Server auch bei kurzen Lesevorgängen korrekt füllen, in [PR](https://github.com/wailsapp/wails/pull/6402) von @taliesin-ai
- Bei nicht durchsuchbaren HEAD-Anfragen im Asset-Server nicht den vollständigen Antworttext lesen, in [PR](https://github.com/wailsapp/wails/pull/6401) von @taliesin-ai
- Beim Auflösen von Verzeichnisanfragen im Asset-Server nur exakt passende Indexdateien ausliefern, in [PR](https://github.com/wailsapp/wails/pull/6400) von @taliesin-ai
- Behoben, dass verworfene Linux-Benachrichtigungen auf Servern ohne Aktionsunterstützung als Klicks gemeldet wurden, in [PR](https://github.com/wailsapp/wails/pull/6367) von @bercedev
- Entwicklungsprozesse beim Schließen des Terminals beenden, in [PR](https://github.com/wailsapp/wails/pull/6371) von @xmasdev
- Beim Neustart der Windows-Taskleiste vor dem Löschen des Flags eine Sekunde Entprellzeit einhalten, in [PR](https://github.com/wailsapp/wails/pull/6386) von @taliesin-ai
- Verlust von Fallback-Antwortheadern im Asset-Server behoben, in [PR](https://github.com/wailsapp/wails/pull/6399) von @taliesin-ai
- Header-Löschungen im Asset-Server vor der Fallback-Auswahl berücksichtigen, in [PR](https://github.com/wailsapp/wails/pull/6398) von @taliesin-ai
- Entfernte Referenzen der Ringwarteschlange im Asset-Server freigeben, um Speicherlecks zu vermeiden, in [PR](https://github.com/wailsapp/wails/pull/6397) von @taliesin-ai
- Verbrauchte Puffer der Inhaltserkennung im Asset-Server freigeben, um den Speicherverbrauch zu senken, in [PR](https://github.com/wailsapp/wails/pull/6396) von @taliesin-ai
- Asset-Dateierweiterungen mit gemischter Groß-/Kleinschreibung beim Ausliefern von Assets erkennen, in [PR](https://github.com/wailsapp/wails/pull/6395) von @taliesin-ai
- Inhaltstyp im Asset-Server korrekt erkennen, ohne Dateinamen im Speicher zu behalten, in [PR](https://github.com/wailsapp/wails/pull/6394) von @taliesin-ai
- Auf einen Standardwert zurückfallen, wenn ein Desktop-Eintrag keinen Befehl enthält, in [PR](https://github.com/wailsapp/wails/pull/6392) von @taliesin-ai
- Verhindern, dass das liquid-glass-Beispiel unter anderen Systemen als macOS vor dem Start der App einen Dialog anzeigt, in [PR](https://github.com/wailsapp/wails/pull/6389) von @taliesin-ai
- Fenstergrößenbeschränkungen unter Windows beibehalten, wenn der Vollbildmodus von einer beschränkten Größe aus gestartet wird, in [PR](https://github.com/wailsapp/wails/pull/6387) von @taliesin-ai
- Behoben, dass vor Run aufgerufene Dialoge dauerhaft blockierten, statt zurückzukehren, in [PR](https://github.com/wailsapp/wails/pull/6390) von @taliesin-ai

## Sicherheit
- Generierte Symbolgrößen vor der Bildspeicherzuweisung validieren, um übermäßigen Speicherverbrauch zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/6478) von @taliesin-ai
- API-Antworten für Versionshinweise vor der Verwendung validieren, in [PR](https://github.com/wailsapp/wails/pull/6461) von @taliesin-ai
- Vollständige hexadezimale Statement-Handles im SQLite-Plugin verlangen, in [PR](https://github.com/wailsapp/wails/pull/6452) von @taliesin-ai
- Fallback-Treffer der Dateisystembibliothek einschränken, um das Laden unbeabsichtigter Bibliotheken zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/6417) von @taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## Geändert
- Die Seite 404 mit den entsprechenden Dokumentationsseiten für Wails v2 auf v2.wails.io verknüpfen, in [PR](https://github.com/wailsapp/wails/pull/6366) von @taliesin-ai
- Die Dokumentationswebsite für v2 nach v2.wails.io verschieben, in [PR](https://github.com/wailsapp/wails/pull/6365) von @taliesin-ai
- Die Dokumentation für v3 über wails.io und die Dokumentation für v2 über v2.wails.io bereitstellen, in [PR](https://github.com/wailsapp/wails/pull/6364) von @taliesin-ai
- Versionshinweise für v3.0.0 aktualisieren und die Aussage im Sicherheitshinweis korrigieren, in [PR](https://github.com/wailsapp/wails/pull/6363) von @taliesin-ai

## Behoben
- Callbacks von Go-Ereignis-Listenern seriell in der Reihenfolge ihrer Auslösung zustellen, in [PR](https://github.com/wailsapp/wails/pull/6275) von @DevLumuz
- Die Eigenschaft Menu des Tray-Icons immer bekannt geben, damit nach Run erstellte Tray-Icons unter Linux nicht verworfen werden, in [PR](https://github.com/wailsapp/wails/pull/6284) von @AlbinoGeek
- Vor Run eingereichte Arbeiten für den Hauptthread bis zum Start der Anwendung zurückhalten, in [PR](https://github.com/wailsapp/wails/pull/6287) von @AlbinoGeek
- Beim Herunterladen von Release-Assets für den GitHub-Updater die korrekte Asset-API-URL verwenden, in [PR](https://github.com/wailsapp/wails/pull/6304) von @gavidroselj
- Bildschirmgrenzen und -positionen in macOS- und Linux-Konfigurationen mit unterschiedlichen DPI-Werten korrigieren, in [PR](https://github.com/wailsapp/wails/pull/6346) von @julianstorer
- Anfragetexte unter Windows beibehalten, wenn WebView2 den Content-Length-Header auslässt, in [PR](https://github.com/wailsapp/wails/pull/6326) von @strraaa
- Datenrennen bei der Ausgabe von Ereignis-Hooks beheben, indem die Hooks unter Sperre als Snapshot erfasst werden, in [PR](https://github.com/wailsapp/wails/pull/6276) von @DevLumuz
- Behoben, dass `SystemTray.SetTooltip()` unter Linux wirkungslos blieb: Die Tooltip-Beschreibung des StatusNotifierItem wird nun gesetzt (#6017)
- Nachweisprüfung für das Änderungsprotokoll im wöchentlichen Release-Workflow ausführen, in [PR](https://github.com/wailsapp/wails/pull/6374) von @taliesin-ai
- Nachweisprüfung für das Release-Änderungsprotokoll korrigieren, sodass Läufe Squash-Merge-Commits zugeordnet werden, in [PR](https://github.com/wailsapp/wails/pull/6373) von @taliesin-ai
- Behoben, dass die Linux-Dunkelmoduserkennung auf einem dunklen Desktop bis zur Änderung des Farbschemas den hellen Modus zurückgab und `linux:SystemThemeChanged` im GTK4-Backend nie ausgelöst wurde. Beide Backends verwenden jetzt eine gemeinsame Implementierung, die auf dem Portal-Namespace `org.freedesktop.appearance` basiert (#6072)
- `Window.SetZoom` hebt Faktoren unter 1.0 bei Fenstern mit Rahmen (Windows, macOS, Linux) nicht mehr auf 1.0 an: WebView2 `ZoomFactor`, WKWebView `magnification` und die Zoomstufe von WebKitGTK akzeptieren jeden positiven Wert, und die Fensteroption `Zoom` leitete solche Werte bereits unverändert weiter. Rahmenlose Fenster behalten die Untergrenze 1.0 bei, da ihre Zieh- und Größenänderungsbereiche nicht an den Zoom angepasst werden ([#4590](https://github.com/wailsapp/wails/issues/4590)). Nicht positive, NaN- und Inf-Faktoren werden jetzt ignoriert, statt auf 1.0 angehoben zu werden — von @APshenkin

## v3.0.0 - 2026-10-26

## Geändert
- Nächtliche Vorabversionen werden als `vX.Y.0-nightly.YYYYMMDD` veröffentlicht und geben eine Vorschau auf das nächste wöchentliche Release. Wöchentliche Releases erhöhen die Minor-Version (`vX.Y.0`), sodass Go mit `@latest` immer ein wöchentliches oder Hotfix-Release auflöst. Die Version des npm-Pakets `@wailsio/runtime` entspricht immer der Wails-Releaseversion; Nightlies verwenden das dist-tag `nightly`
- Projektvorlagen fixieren `@wailsio/runtime` auf die Version der CLI, mit der sie generiert wurden
- Wails v3 erfordert Go 1.26 oder neuer (Go 1.26.10 beziehungsweise 1.27.3 und neuer enthalten die Sicherheitskorrektur `crypto/x509`)

## Behoben
- Behoben, dass `wails3 update cli` keine Releases von v3 fand: Es berücksichtigte nur v2-Tags und installierte aus dem falschen Pfad
- Abstürze von Windows-Apps mit `panic: SetTextColor failed` beim Zeichnen der Anwendungsmenüleiste in einigen Windows-Builds behoben: Die Menüleiste zeichnet in einen gültigen Gerätekontext oder greift auf das Standardzeichnen zurück, `UAHMENUITEM.Position` entspricht jetzt dem C-`int`, und das Zeichnen des Menüs löst niemals eine Panic aus (#6311)
- Das bereitgestellte Update bleibt sicher, solange ein Neustart aussteht: Ein Download oder ein zweiter Neustart, der während des Herunterfahrens gestartet wird, kann die Nutzlast, auf deren Installation der Austausch-Helfer wartet, weder löschen noch mit ihr in Konflikt geraten — von @julianstorer
- Abstürze von Windows-Apps in `applyDarkMode` behoben, wenn die Initialisierung von `pkg/w32` `uxtheme.dll` entlud, während dessen ordinale Funktionszeiger noch verwendet wurden: Das Modul bleibt jetzt während der gesamten Prozesslaufzeit geladen (#6356)
- Die bereitgestellten Nutzdaten vor dem Austausch und die installierte Anwendung vor dem Verwerfen der Sicherung validieren, damit ein beschädigtes Update zurückgesetzt wird, statt eine funktionierende Anwendung zu ersetzen, von @julianstorer

## Entfernt
- Crowdin-Integration nach der abschließenden Synchronisierung der Übersetzungen in [PR](https://github.com/wailsapp/wails/pull/5736) entfernen, von @github-actions[bot]

## Sicherheit
- Reservierte `wails:`-Nachrichten aus ursprungsübergreifenden untergeordneten Frames unter Linux (GTK3 und GTK4) blockieren, indem ein Schutzskript in alle Frames injiziert wird; dies entschärft GHSA-2p5f-2gx8-mh8m
- `github.com/klauspost/compress` auf 1.18.7 (GHSA-259r-337f-4rfw) und `golang.org/x/image` auf 0.45.0 aktualisieren sowie mit Go 1.26.10 bauen, um die Korrektur für `crypto/x509` zu übernehmen (GO-2026-6809)
- `vite` im Runtime-Paket, Einrichtungsassistenten, in Beispielen und Projektvorlagen auf 8.2.4 oder neuer aktualisieren (GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv)

### Bekannte Einschränkungen
- **Linux:** Der Schutz für untergeordnete Frames kann nicht verhindern, dass ein untergeordneter Frame mit *demselben Ursprung* (beispielsweise ein `srcdoc`- oder `about:blank`-iframe) den nativen Nachrichten-Handler über seinen übergeordneten Frame erreicht, da WebKitGTK dem Handler keine Frame-Informationen bereitstellt. Ursprungsübergreifende Frames und Sandbox-Frames ohne `allow-same-origin` werden blockiert. Betten Sie keine nicht vertrauenswürdigen Inhalte desselben Ursprungs in ein Wails-Fenster ein.
- **Android (experimentell):** Die JavaScript-Bridge `window.wails` ist für jeden Frame im WebView zugänglich, ohne Prüfung des Haupt-Frames oder Ursprungs. Inhalte in einem eingebetteten iframe, einschließlich ursprungsübergreifender Inhalte, können gebundene Methoden aufrufen und reservierte `wails:`-Nachrichten senden. Betten Sie keine nicht vertrauenswürdigen Drittanbieterinhalte in den WebView der App ein. Dies unterscheidet sich von macOS und iOS, die `wails:`-Nachrichten aus untergeordneten Frames ignorieren. Eine Korrektur ist für 3.0.1 geplant.
- **Alle Plattformen:** Aufrufe gebundener Methoden aus untergeordneten Frames werden über HTTP an `/wails/runtime` übertragen, nicht über den nativen Nachrichten-Handler. Der Schutz schränkt sie daher nicht ein. Behandeln Sie jeden eingebetteten iframe so, als hätte er Zugriff auf Ihre gebundenen Go-Methoden.

## v3.0.0-beta.34 - 2026-10-20

## Behoben
- Behoben: Der Windows-WebView erhielt bei der Fensteraktivierung den Fokus nicht zurück und Tab blieb im WebView gefangen, statt zum nativen Fenster und zurück zu wechseln (#3783)
- Behoben: In den COM-vtable-Strukturen von `internal/webview2/pkg/webview2` (`ICoreWebView2_2`..`_27` und weiteren abgeleiteten Schnittstellen) fehlte die Kette geerbter Methoden, wodurch Aufrufe im falschen vtable-Slot landeten (#6129)
- Behoben: Windows-Fragedialoge ignorierten benutzerdefinierte Schaltflächenbeschriftungen und lösten nie Schaltflächen-Callbacks aus. Benutzerdefinierte Schaltflächen werden jetzt mit einem TaskDialog angezeigt; wenn comctl32 v6 nicht verfügbar ist, wird auf eine MessageBox mit positionsbasierter Callback-Zuordnung zurückgegriffen (#5794)
- Behoben: `wails3 dev` blieb bei 60s mit Vite hängen. Die automatische Prüfung der Frontend-Bereitschaft verwendet jetzt eine TCP-Sonde statt eines HTTP-GET, der keinen `Accept`-Header sendete und eine 404-Antwort erhielt ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Systray-Menüs bleiben funktionsfähig, wenn sie im geöffneten Zustand aktualisiert werden, in [PR](https://github.com/wailsapp/wails/pull/6301) von @taliesin-ai
- Vorübergehende WebView2-COM-Fehler bei Runtime-Operationen ordnungsgemäß behandeln, in [PR](https://github.com/wailsapp/wails/pull/6302) von @leaanthony

## v3.0.0-beta.33 - 2026-10-19

## Hinzugefügt
- Der Binding-Generator erzeugt jetzt ein `bindings/events`-Modul, das für jedes auffindbare, mit `application.RegisterEvent` unter Verwendung eines zur Kompilierzeit konstanten Namens registrierte Ereignis eine typisierte Namenskonstante exportiert (z. B. `UserLogin` für `user:login`). Diese kann ohne das Vite-Plug-in mit `Events.On`/`Emit` verwendet werden (#5645)

## Behoben
- Behoben: Windows-Fragedialoge ignorierten benutzerdefinierte Schaltflächenbeschriftungen und lösten nie Schaltflächen-Callbacks aus. Benutzerdefinierte Schaltflächen werden jetzt mit einem TaskDialog angezeigt; wenn comctl32 v6 nicht verfügbar ist, wird auf eine MessageBox mit positionsbasierter Callback-Zuordnung zurückgegriffen (#5794)
- Behoben: `wails3 dev` blieb bei 60s mit Vite hängen. Die automatische Prüfung der Frontend-Bereitschaft verwendet jetzt eine TCP-Sonde statt eines HTTP-GET, der keinen `Accept`-Header sendete und eine 404-Antwort erhielt ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Systray-Menüs bleiben funktionsfähig, wenn sie im geöffneten Zustand aktualisiert werden, in [PR](https://github.com/wailsapp/wails/pull/6301) von @taliesin-ai
- Vorübergehende WebView2-COM-Fehler bei Runtime-Operationen ordnungsgemäß behandeln, in [PR](https://github.com/wailsapp/wails/pull/6302) von @leaanthony

## v3.0.0-beta.32 - 2026-10-18

## Hinzugefügt
- Native macOS-Menüabschnittsüberschriften für macOS 14+ hinzufügen – siehe die [Dokumentationswebsite](https://v3.wails.io) (20 Seiten aktualisiert), in [PR](https://github.com/wailsapp/wails/pull/6297) von @leaanthony

## Behoben
- SetRelativePosition unter macOS relativ zum Arbeitsbereich des Bildschirms positionieren, in [PR](https://github.com/wailsapp/wails/pull/6299) von @leaanthony
- Der Updater-Helfer startet nach einer Zeitüberschreitung des übergeordneten Prozesses oder einem fehlgeschlagenen Austausch das unveränderte Ziel erneut, in [PR](https://github.com/wailsapp/wails/pull/6298) von @taliesin-ai
- pnpm für versionierte Dokumentation auf Version 10 festsetzen, in [PR](https://github.com/wailsapp/wails/pull/6293) von @taliesin-ai

## v3.0.0-beta.31 - 2026-10-18

## Hinzugefügt
- Den Unterpfadexport `@wailsio/runtime/cancellable` bereitstellen, damit `CancellablePromise` importiert werden kann, ohne den Runtime-Haupteinstieg auszuwerten (#6153)

## Geändert
- Go 1.26 oder neuer voraussetzen: Das Modul v3 deklariert jetzt `go 1.26.0`, und von `wails3 init` und `wails3 service init` erzeugte Projekte deklarieren dieselbe Mindestversion

## Behoben
- Behoben: Die Ausgabe des Binding-Generators änderte sich beim Bauen mit Go 1.27, da `encoding/json`-Typen wie `Marshaler`, `Unmarshaler` und `RawMessage` zu Aliasen von `encoding/json/v2`- und `encoding/json/jsontext`-Typen wurden. Sie werden wieder als `any` dargestellt, und für `encoding/json/v2` werden keine Modelle erzeugt (#6292)
- An `wails3 init -t` übergebene URLs entfernter Vorlagen vor dem Klonen validieren: Akzeptiert werden nur `https://`-, `http://`-, `ssh://`-, `git://`- und `user@host:path`-URLs; Werte werden nach `--` an git übergeben, und git wird während des Klonens auf diese Transportprotokolle beschränkt
- Außerhalb des gültigen Bereichs liegende Werte des `x-wails-window-id`-Headers im HTTP-Runtime-Transport zurückweisen, statt sie stillschweigend auf die ID eines anderen Fensters zu kürzen
- Die Archivextraktion des Updaters gegen Zip-Slip absichern: Eintragspfade werden jetzt mit `filepath.IsLocal` sowie einer Einschlussprüfung des zusammengesetzten Pfads geprüft. Tests decken verschachtelte `..`-, Backslash-, Absolutpfad- und durch Symlinks hindurch schreibende Einträge in zip- und tar.gz-Archiven ab
- Behoben: AppImages aus `wails3 generate appimage` stürzten außerhalb von Debian und Ubuntu beim Start mit „Unable to spawn a new child process: … WebKitNetworkProcess“ ab. Das gebündelte WebKit führt jetzt die im AppImage enthaltenen Hilfsprozesse aus. Die bubblewrap-Sandbox von WebKit bleibt aktiviert
- Behoben: `-tags server`-Builds konnten für Android und iOS nicht kompiliert werden; damit wird [#6271](https://github.com/wailsapp/wails/issues/6271) behoben
- Behoben: Der Updater konnte unter Linux keinen Austausch durchführen, wenn sich das temporäre Verzeichnis auf einem anderen Dateisystem befand (beispielsweise einem tmpfs unter `/tmp`). Die Bereitstellung erfolgt jetzt neben der installierten Binärdatei, wobei über Dateisystemgrenzen hinweg kopiert wird
- Behoben: Der Updater versuchte, den schreibgeschützten Mount eines AppImage statt der Datei `.AppImage` zu ersetzen
- Behebt, dass reservierte `wails:`-Nachrichten zur Fenstersteuerung unter macOS und iOS aus iframes akzeptiert wurden. Sie werden jetzt nur noch vom Hauptframe des Fensters akzeptiert; andere Nachrichten aus iframes erreichen `RawMessageHandler` weiterhin einschließlich ihrer Ursprungsinformationen
- Sendet `WindowDidUpdate` und `WindowWillUpdate` unter macOS nicht mehr: `NSWindowDidUpdateNotification` wird bei jedem Durchlauf der Run Loop ausgelöst, sodass die Seite pro Sekunde etwa 120 leere `evaluateJavaScript:`-Aufrufe empfing. Unter macOS verbrauchte 27.0 bei jedem dieser Aufrufe die vorübergehende Aktivierung durch den Klick des Benutzers. Dadurch wurden `navigator.clipboard.writeText()` und andere aktivierungsgebundene APIs bei echten Klicks mit `NotAllowedError` abgelehnt ([#6170](https://github.com/wailsapp/wails/issues/6170)) von @APshenkin
- Setzt die Standardsprache der Windows-Datei info.json auf en-US (0409), damit der Explorer Versionsdetails anzeigt (#6210)
- Behebt, dass der Linux-DMA-BUF-Renderer auf Laptops mit einer NVIDIA-dGPU neben einer AMD- oder Intel-iGPU deaktiviert wurde: Unter Wayland wird er jetzt nur deaktiviert, wenn WebKit auf der NVIDIA-GPU rendert, weil die App auf die dGPU ausgelagert wurde oder der Compositor darauf läuft. Er wird nicht mehr bereits dann deaktiviert, wenn das nvidia-Modul geladen ist, was das Rendering auf der iGPU extrem verlangsamte. X11-Sitzungen behalten das bisherige Verhalten bei, da die rendernde GPU von Xorg durch dessen Konfiguration festgelegt wird
- Behandelt im Linux-Updater fehlgeschlagene geräteübergreifende Umbenennungen durch einen Fallback auf eine stufenweise Kopie in [PR](https://github.com/wailsapp/wails/pull/6135) von @john-okeefe
- Behebt, dass `WebviewWindow.SetMenu` unter Windows die Menüleiste nur neu zeichnete, wenn ein Menüleisten-Theme konfiguriert war; behebt [#6103](https://github.com/wailsapp/wails/issues/6103)
- Behebt, dass `WebviewWindow.SetMenu(nil)` unter Windows eine Panic auslöste, statt wie unter macOS keine Aktion auszuführen; behebt [#6104](https://github.com/wailsapp/wails/issues/6104)
- Bewahrt Changelog-Einträge bei konkurrierenden Pushes durch Wiederholungsversuche in [PR](https://github.com/wailsapp/wails/pull/6094) von @leaanthony
- Behebt, dass `go mod vendor` auf allen Plattformen mit `pattern arm64/WebView2Loader.dll: no matching files found` fehlschlug, indem Einbettungen entfernt wurden, die auf nie mit dem Modul ausgelieferte Binärdateien verwiesen. Behebt [#5782](https://github.com/wailsapp/wails/issues/5782) und [#5376](https://github.com/wailsapp/wails/issues/5376) in [PR](https://github.com/wailsapp/wails/pull/6031) von @Grantmartin2002
- Verhindert den doppelten Abbau von Fenstern unter macOS in [PR](https://github.com/wailsapp/wails/pull/6159) von @charlesbalderson
- Behandelt WM_CONTEXTMENU, damit sich das Systray-Menü per Rechtsklick öffnet, in [PR](https://github.com/wailsapp/wails/pull/6162) von @wzk2025
- Behebt, dass alle fünf `InvokeSync*`-Varianten ihren Aufrufer dauerhaft blockierten, wenn der weitergeleitete Callback eine Panic auslöste. Behebt [#6107](https://github.com/wailsapp/wails/issues/6107) in [PR](https://github.com/wailsapp/wails/pull/6108) von @RALIST
- Erstellt zerstörte GTK3-Untermenüs sicher neu in [PR](https://github.com/wailsapp/wails/pull/6236) von @leaanthony
- Implementiert unter Windows einen Weiterleitungs-Shim für verschachtelte Worker in [PR](https://github.com/wailsapp/wails/pull/6261) von @taliesin-ai
- Behebt, dass `OpenFile()` unter Windows `PromptForMultipleSelection()` bei `CanChooseDirectories(true)` ignorierte: Benutzer können jetzt wie unter macOS und Linux mehrere Ordner gleichzeitig auswählen (#6263)
- Behebt als Wert übergebene `double`-Argumente für WebView2-COM-Methoden unter windows/arm64. Betroffen waren dort `ZoomFactor`, `RasterizationScale`, Cookie-`Expires` sowie Skalierung, Seitengröße und Ränder der Druckeinstellungen; die Setter der Druckeinstellungen gaben zuvor `ErrDoubleArgUnsupported` zurück. Betroffene Aufrufe laufen jetzt über ein `runtime.cgocall`-Trampolin, das `d0` lädt, was über den syscall-Pfad von Go nicht möglich ist (golang.org/issue/62583), von @APshenkin

## Sicherheit
- Der Servermodus (`-tags server`) akzeptiert `/wails/runtime`-Aufrufe jetzt ausschließlich als `POST`-Anfragen vom eigenen Ursprung des Servers oder von einem in `ServerOptions.WebSocketOriginPatterns` aufgeführten Ursprung. Andere Methoden erhalten HTTP 405, andere Browserursprünge HTTP 403. Bereitstellungen hinter einem Reverse Proxy, der den `Host`-Header umschreibt, sollten ihren öffentlichen Ursprung zu `WebSocketOriginPatterns` hinzufügen
- Der an eine Loopback-Adresse gebundene Servermodus (`-tags server`; standardmäßig `localhost`, `127.0.0.0/8` oder `::1`) lehnt Anfragen jetzt mit HTTP 403 ab, sofern der `Host`-Header nicht `localhost`, ein IP-Literal, der Bind-Host oder ein mit `ServerOptions.WebSocketOriginPatterns` übereinstimmender Host ist. Dies blockiert DNS-Rebinding-Angriffe auf `/wails/runtime` und die WebSocket-Endpunkte
- Windows-Benachrichtigungen greifen bei einem Fehler der WinRT-Toast-API nicht mehr auf PowerShell zurück; stattdessen wird der Fehler zurückgegeben

## v3.0.0-beta.30 - 2026-10-15

## Behoben
- Behebt, dass der Updater versuchte, den schreibgeschützten Einhängepunkt eines AppImage statt der `.AppImage`-Datei zu ersetzen
- Sendet `WindowDidUpdate` und `WindowWillUpdate` unter macOS nicht mehr: `NSWindowDidUpdateNotification` wird bei jedem Durchlauf der Run Loop ausgelöst, sodass die Seite pro Sekunde etwa 120 leere `evaluateJavaScript:`-Aufrufe empfing. Unter macOS verbrauchte 27.0 bei jedem dieser Aufrufe die vorübergehende Aktivierung durch den Klick des Benutzers. Dadurch wurden `navigator.clipboard.writeText()` und andere aktivierungsgebundene APIs bei echten Klicks mit `NotAllowedError` abgelehnt ([#6170](https://github.com/wailsapp/wails/issues/6170)) von @APshenkin
- Setzt die Standardsprache der Windows-Datei info.json auf en-US (0409), damit der Explorer Versionsdetails anzeigt (#6210)
- Behandelt im Linux-Updater fehlgeschlagene geräteübergreifende Umbenennungen durch einen Fallback auf eine stufenweise Kopie in [PR](https://github.com/wailsapp/wails/pull/6135) von @john-okeefe
- Behebt, dass `WebviewWindow.SetMenu` unter Windows die Menüleiste nur neu zeichnete, wenn ein Menüleisten-Theme konfiguriert war; behebt [#6103](https://github.com/wailsapp/wails/issues/6103)
- Behebt, dass `WebviewWindow.SetMenu(nil)` unter Windows eine Panic auslöste, statt wie unter macOS keine Aktion auszuführen; behebt [#6104](https://github.com/wailsapp/wails/issues/6104)
- Bewahrt Changelog-Einträge bei konkurrierenden Pushes durch Wiederholungsversuche in [PR](https://github.com/wailsapp/wails/pull/6094) von @leaanthony
- Behebt, dass `go mod vendor` auf allen Plattformen mit `pattern arm64/WebView2Loader.dll: no matching files found` fehlschlug, indem Einbettungen entfernt wurden, die auf nie mit dem Modul ausgelieferte Binärdateien verwiesen. Behebt [#5782](https://github.com/wailsapp/wails/issues/5782) und [#5376](https://github.com/wailsapp/wails/issues/5376) in [PR](https://github.com/wailsapp/wails/pull/6031) von @Grantmartin2002
- Verhindert den doppelten Abbau von Fenstern unter macOS in [PR](https://github.com/wailsapp/wails/pull/6159) von @charlesbalderson
- Behandelt WM_CONTEXTMENU, damit sich das Systray-Menü per Rechtsklick öffnet, in [PR](https://github.com/wailsapp/wails/pull/6162) von @wzk2025
- Behebt, dass alle fünf `InvokeSync*`-Varianten ihren Aufrufer dauerhaft blockierten, wenn der weitergeleitete Callback eine Panic auslöste. Behebt [#6107](https://github.com/wailsapp/wails/issues/6107) in [PR](https://github.com/wailsapp/wails/pull/6108) von @RALIST

## v3.0.0-beta.29 - 2026-10-13

## Hinzugefügt
- Globale Tastenkürzel unter macOS berücksichtigen jetzt das Tastaturlayout: Über `UCKeyTranslate` werden Buchstaben-Tastenkürzel bei Nicht-QWERTY-Layouts (AZERTY, QWERTZ, Dvorak, ...) an die beschriftete Taste statt an die feste physische QWERTY-Position gebunden.

## Behoben
- Initialisiert das native Menü unter Linux, bevor es gesetzt wird, in [PR](https://github.com/wailsapp/wails/pull/6233) von @abhicodes-007
- Behebt macOS-Regressionen bei der GA-Bereitschaft in der Ausgabeverarbeitung und im Einrichtungsassistenten in [PR](https://github.com/wailsapp/wails/pull/6237) von @taliesin-ai
- Behebt, dass die WebView2-COM-Wrapper `AddRef`/`Release` unter Windows einen veralteten Last-Error des Threads als Fehler meldeten. Sie geben jetzt den Referenzzähler zurück, sodass eine Freigabe die Anwendung nicht mehr abbrechen oder eine Asset-Anfrage fehlschlagen lassen kann (#3931)
- Behebt, dass `ICoreWebView2Deferral.Release` den `AddRef`-Slot aufrief, wodurch der Aufschub jeder zurückgestellten Asset-Anfrage unter Windows nicht freigegeben wurde (#3931)

## v3.0.0-beta.28 - 2026-10-05

## Behoben
- Speichert die macOS-App-Ausgabe außerhalb von Dateianbietern in [PR](https://github.com/wailsapp/wails/pull/6225) von @taliesin-ai
- Vervollständigt den Einrichtungsassistenten und zeigt Systemdetails an in [PR](https://github.com/wailsapp/wails/pull/6224) von @taliesin-ai
- Verbessert unter Windows die Fehlerbehandlung und Ressourcenbereinigung bei der WebView2-Wiederherstellung in [PR](https://github.com/wailsapp/wails/pull/6222) von @leaanthony
- Verbessert die Wiederherstellung von WebView2 unter Windows, um App-Abstürze und verwaiste Controller zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/6220) von @taliesin-ai
- Korrigiert die Anleitung zur Medienautorisierung unter macOS in [PR](https://github.com/wailsapp/wails/pull/6221) von @taliesin-ai
- Behebt, dass ein Rechtsklick auf ein Linux-Taskleistensymbol neben dem Öffnen des Menüs auch den Klick-Handler ausführte, sodass ein angehängtes Fenster bei jedem Rechtsklick umgeschaltet wurde (#6018)
- Hält WebView2-Callback-Handler unter Windows bis zur nativen Freigabe aktiv, in [PR](https://github.com/wailsapp/wails/pull/6184) von @taliesin-ai

## v3.0.0-beta.27 - 2026-10-01

## Behoben
- Behebt, dass die Option `Permissions` unter macOS 12 und neuer ignoriert wurde: Die Methode `WKUIDelegate` zur Medienaufnahme ist jetzt implementiert, sodass `PermissionAllow` und `PermissionDeny` wie unter Linux und Windows für Kamera- und Mikrofonanfragen gelten. `NSCameraUsageDescription` / `NSMicrophoneUsageDescription` sowie bei Verwendung der Sandbox die entsprechenden Geräteberechtigungen sind weiterhin erforderlich

## v3.0.0-beta.26 - 2026-09-25

## Hinzugefügt
- Der Watcher führt HTTP-Bereitschaftsprüfungen für den Frontend-Entwicklungsserver durch, in [PR](https://github.com/wailsapp/wails/pull/6150) von @atterpac

## Behoben
- Stellt nach Ausfällen des WebView2-Prozesses die Funktion wieder her, anstatt ein leeres Fenster zurückzulassen, in [PR](https://github.com/wailsapp/wails/pull/6002) von @taliesin-ai

## v3.0.0-beta.25 - 2026-09-22

## Behoben
- Beendet `wails3 dev` und dessen Hintergrundprozesse, wenn die primäre Anwendung beendet wird (#6048)

## v3.0.0-beta.24 - 2026-09-20

## Sicherheit
- Beschränkt die Berechtigungen von GITHUB_TOKEN in Build- und Veröffentlichungs-Workflows gemäß dem Prinzip der geringsten Rechte
- Behebt Path-Traversal- und Symlink-Escape-Schwachstellen in der Asset-Middleware des Bildschirmbeispiels
- Erhält die Funktionsfähigkeit der Abhängigkeitsinstallation im Einrichtungsassistenten und beschränkt sie zugleich auf genehmigte Paketmanagerbefehle

## v3.0.0-beta.23 - 2026-09-16

## Hinzugefügt

- Ermöglicht das Überschreiben der GTK-Anwendungs-ID unter Linux, in [PR](https://github.com/wailsapp/wails/pull/5972) von @overlordtm

## Behoben

- Die Linux-Einzelinstanz beansprucht den Busnamen unter der eigenen ID der App, in [PR](https://github.com/wailsapp/wails/pull/5973) von @overlordtm
- Bricht abgebrochene Anfragen für benutzerdefinierte Schemas unter Linux ab, in [PR](https://github.com/wailsapp/wails/pull/6126) von @leaanthony

## v3.0.0-beta.22 - 2026-09-14

## Hinzugefügt

- Vervollständigt Übersetzungen und portiert Community-Medien – siehe die [Dokumentationswebsite](https://v3.wails.io) (1812 Seiten aktualisiert), in [PR](https://github.com/wailsapp/wails/pull/6122) von @leaanthony
- Dokumentiert einbettbare Wails-Schaltflächen in verschiedenen Stilen und Größen, in [PR](https://github.com/wailsapp/wails/pull/6120) von @leaanthony

## v3.0.0-beta.21 - 2026-09-13

## Hinzugefügt

- Stellt die Wails-v3-Dokumentation mit M-Press bereit, in [PR](https://github.com/wailsapp/wails/pull/6116) von @leaanthony

## Behoben

- Verarbeitet JSON-Slug-Werte aus dem MPD-Frontmatter für die Changelog-Erzeugung, in [PR](https://github.com/wailsapp/wails/pull/6118) von @leaanthony
- Der Updater löscht Hilfsumgebungsvariablen und startet nach fehlgeschlagenen Sicherungen das ursprüngliche Ziel erneut, in [PR](https://github.com/wailsapp/wails/pull/6080) von @cnmax
- Startet den standardmäßigen Signal-Handler während App.Run, in [PR](https://github.com/wailsapp/wails/pull/6098) von @leaanthony
- Das Windows-Menü verarbeitet nil-Menüs, gibt ersetzte Ressourcen frei und zeichnet die Menüleiste neu, in [PR](https://github.com/wailsapp/wails/pull/6112) von @taliesin-ai
- Stellt die MSIX-Paketierung für neue Projekte mit gemeinsamer YAML-Konfiguration wieder her, in [PR](https://github.com/wailsapp/wails/pull/6115) von @leaanthony
- Behebt, dass erzeugte JavaScript- und TypeScript-Bindings nicht geladen werden konnten, wenn Erzeuger generischer Modelle auf später deklarierte Hilfsfunktionen verwiesen, und verhindert Stacküberläufe beim Erstellen voneinander abhängiger generischer Modelle (#6062)

## v3.0.0-beta.20 - 2026-09-10

## Geändert

- Aktualisiert die Links der Clave-Präsentation auf die aktuelle Website und das aktuelle Repository, in [PR](https://github.com/wailsapp/wails/pull/6082) von @01xR4in

## Behoben

- Bricht abgebrochene Windows-Asset-Anfragen einschließlich Worker-Anfragen ab und erhält Keepalive-Handler bei der Navigation. Leitet native Anfragekontexte auf Apple-Plattformen durch den Anwendungs-Wrapper weiter. (#5963, #5969)
- Bewahrt Changelog-Einträge bei konkurrierenden Pushes durch Wiederholungsversuche, in [PR](https://github.com/wailsapp/wails/pull/6094) von @leaanthony
- Behebt, dass `go mod vendor` auf jeder Plattform mit `pattern arm64/WebView2Loader.dll: no matching files found` fehlschlug, indem die Einbettungen entfernt wurden, die auf nie mit dem Modul ausgelieferte Binärdateien verwiesen. Dies behebt [#5782](https://github.com/wailsapp/wails/issues/5782) und [#5376](https://github.com/wailsapp/wails/issues/5376), in [PR](https://github.com/wailsapp/wails/pull/6031) von @Grantmartin2002

## Entfernt

- Entfernt die Unterstützung für den nativen WebView2-Loader, der durch den reinen Go-Loader ersetzt wurde. Dadurch entfallen die eingebetteten `WebView2Loader.dll`-Binärdateien und die Abhängigkeit `github.com/jchv/go-winloader`. Das Build-Tag `native_webview2loader` wird weiterhin akzeptiert und verursacht keinen Fehler mehr, hat bei v3-Builds jedoch keine Wirkung, in [PR](https://github.com/wailsapp/wails/pull/6031) von @Grantmartin2002
- Entfernt ungenutzte Build-Tags und die FPS-Option aus dem macOS-API-Leitfaden, in [PR](https://github.com/wailsapp/wails/pull/6097) von @leaanthony

## v3.0.0-beta.19 - 2026-09-09

## Hinzugefügt

- Schirmt private macOS-APIs durch Build-Tags für die explizite Nutzung ab – siehe [Dokumentation](https://v3.wails.io/features/browser/integration), [Dokumentation](https://v3.wails.io/features/environment/info), [Dokumentation](https://v3.wails.io/features/windows/basics), [Dokumentation](https://v3.wails.io/features/windows/frameless), [Dokumentation](https://v3.wails.io/features/windows/notch-windows), [Dokumentation](https://v3.wails.io/features/windows/options), [Dokumentation](https://v3.wails.io/guides/build/macos), [Dokumentation](https://v3.wails.io/guides/build/private-macos-apis) und [Dokumentation](https://v3.wails.io/reference/overview), in [PR](https://github.com/wailsapp/wails/pull/6087) von @leaanthony

## Behoben

- Lehnt Laufzeitanfragen über 64 MiB mit HTTP 413 ab, in [PR](https://github.com/wailsapp/wails/pull/6091) von @leaanthony

## Sicherheit

- Sichert MCP-Ursprünge und Remotezugriff durch Token-Authentifizierung ab, in [PR](https://github.com/wailsapp/wails/pull/6092) von @leaanthony

## v3.0.0-beta.18 - 2026-09-08

## Behoben

- Behebt ein Calloc-Speicherleck unter Linux und Darwin durch die Verwendung von Zeiger-Receivern, in [PR](https://github.com/wailsapp/wails/pull/6083) von @4RH1T3CT0R7

## v3.0.0-beta.17 - 2026-09-06

## Behoben

- Windows: Ein fehlgeschlagenes oder nil-`GetRequest` im WebResourceRequested-Handler beendet den Prozess nicht mehr (`log.Fatal` / Panic durch nil-Dereferenzierung). Stattdessen wird die Anfrage verworfen und protokolliert, in [PR](https://github.com/wailsapp/wails/pull/6006) von @midagedev

## v3.0.0-beta.16 - 2026-08-29

## Geändert

- In einem neuen Terminalfenster zur Eingabe des Notarisierungspassworts auffordern, in [PR](https://github.com/wailsapp/wails/pull/6029) von @leaanthony

## Behoben

- Systray-Klicktypen unter macOS korrekt verarbeiten, in [PR](https://github.com/wailsapp/wails/pull/5919) von @ChewbaccaCookie
- CI entfernt vor der Aktualisierung nicht verwendete Microsoft-apt-Repositorys, in [PR](https://github.com/wailsapp/wails/pull/6041) von @Grantmartin2002

## v3.0.0-beta.15 - 2026-08-27

## Behoben

- Zeitüberschreitung für das Einbetten von WebView2 auf 60 Sekunden erhöhen, in [PR](https://github.com/wailsapp/wails/pull/6043) von @Grantmartin2002

## v3.0.0-beta.14 - 2026-08-26

## Behoben

- Tastendrücke mit Control-Buchstaben unter macOS korrekt benennen, in [PR](https://github.com/wailsapp/wails/pull/6032) von @taliesin-ai
- ICO-Systray-Symbole korrigieren und unter Windows dem Taskleisten-Design folgen, in [PR](https://github.com/wailsapp/wails/pull/6016) von @nik9play

## v3.0.0-beta.13 - 2026-08-25

## Behoben

- Während einer modalen Schleife unter macOS weiterhin Aufgaben im Hauptthread verarbeiten, in [PR](https://github.com/wailsapp/wails/pull/6026) von @leaanthony
- Fehler im sicheren Speicher auf Mobilgeräten zulassen und dabei nach dem Fail-Closed-Prinzip vorgehen, in [PR](https://github.com/wailsapp/wails/pull/5923) von @mortenolsrud
- Hooks für Anwendungsereignisse auch ausführen, wenn kein Listener registriert ist, in [PR](https://github.com/wailsapp/wails/pull/5999) von @archy-rock3t-cloud
- Tippfehler in Kommentaren und lokalisierter Dokumentation korrigieren, in [PR](https://github.com/wailsapp/wails/pull/6023) von @haoku123
- Die unter `v3/examples` eingecheckten vorkompilierten macOS-Binärdateien entfernen, in [PR](https://github.com/wailsapp/wails/pull/6025) von @4RH1T3CT0R7

## v3.0.0-beta.12 - 2026-08-21

## Hinzugefügt

- macOS-Benachrichtigungsfenster für die Display-Aussparung mit Lebenszyklus- und Telemetriebeispiel hinzufügen – siehe [Dokumentation](https://v3.wails.io/features/windows/notch-windows), in [PR](https://github.com/wailsapp/wails/pull/6010) von @leaanthony
- Unterstützung für macOS-NSPanel-Fenster mit neuen Optionen und nativer Integration hinzufügen – siehe [Dokumentation](https://v3.wails.io/features/windows/options), in [PR](https://github.com/wailsapp/wails/pull/6008) von @leaanthony

## Behoben

- Verhindern, dass SQLite Prepare bei gleichzeitigen Aufrufen hängen bleibt, in [PR](https://github.com/wailsapp/wails/pull/5998) von @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Entfernt

- Veralteten Implementierungs-Tracker aus der Dokumentation entfernen, in [PR](https://github.com/wailsapp/wails/pull/6005) von @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Behoben

- Beheben, dass der GTK4-Linux-Host Startargumente für benutzerdefinierte Protokolle und Dateizuordnungen verwirft, in [PR](https://github.com/wailsapp/wails/pull/6000) von @midagedev
- Gelöschte Zeilen und Korrekturen aus derselben Quelle bei der Changelog-Validierung korrekt verarbeiten, in [PR](https://github.com/wailsapp/wails/pull/5993) von @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Hinzugefügt

- Den sicheren wails3-mcp-Server für agentengestützte Projektverwaltung hinzufügen, in [PR](https://github.com/wailsapp/wails/pull/5896) von @leaanthony
- Dokumentation für Modelle in Bindings hinzufügen – siehe [Dokumentation](https://v3.wails.io/features/bindings/models), in [PR](https://github.com/wailsapp/wails/pull/5988) von @taliesin-ai
- rpm-ostree-Installation für atomare Linux-Systeme unterstützen, in [PR](https://github.com/wailsapp/wails/pull/5987) von @leaanthony
- Native wöchentliche Generierung und Veröffentlichung eines Diagramms zur Sternhistorie hinzufügen – siehe [Dokumentation](https://v3.wails.io/credits) und [Dokumentation](https://v3.wails.io/de/credits) und [Dokumentation](https://v3.wails.io/fr/credits) und [Dokumentation](https://v3.wails.io/id/credits) und [Dokumentation](https://v3.wails.io/ja/credits) und [Dokumentation](https://v3.wails.io/ko/credits) und [Dokumentation](https://v3.wails.io/pt/credits) und [Dokumentation](https://v3.wails.io/ru/credits) und [Dokumentation](https://v3.wails.io/zh-cn/credits) und [Dokumentation](https://v3.wails.io/zh-tw/credits), in [PR](https://github.com/wailsapp/wails/pull/5986) von @leaanthony
- Nur unter Darwin verfügbares mac-Paket zum Auflösen von Ressourcen im Anwendungspaket hinzufügen – siehe [Dokumentation](https://v3.wails.io/guides/build/macos), in [PR](https://github.com/wailsapp/wails/pull/5965) von @leaanthony
- Showcase-Seite und Indexeintrag für Condui hinzufügen – siehe [Dokumentation](https://v3.wails.io/community/showcase/condui) und [Dokumentation](https://v3.wails.io/community/showcase), in [PR](https://github.com/wailsapp/wails/pull/5962) von @mgueregath
- Showcase-Seite für Redis Viewer mit Screenshots und Projektlink hinzufügen – siehe [Dokumentation](https://v3.wails.io/community/showcase) und [Dokumentation](https://v3.wails.io/community/showcase/redisviewer), in [PR](https://github.com/wailsapp/wails/pull/5984) von @redisviewer

## Geändert

- GTK-Anwendungsflags für Linux auf G<em>APPLICATION</em>NON_UNIQUE aktualisieren, in [PR](https://github.com/wailsapp/wails/pull/5971) von @overlordtm
- Fehlende Fensterereignisse auf Debug- statt auf Warnstufe protokollieren, in [PR](https://github.com/wailsapp/wails/pull/5914) von @julianstorer

## Behoben

- Registrierten macOS-Tastenkürzeln Vorrang vor der Webview geben, in [PR](https://github.com/wailsapp/wails/pull/5902) von @julianstorer
- Defekte Links in der Seitenleiste der Dokumentation reparieren, in [PR](https://github.com/wailsapp/wails/pull/5937) von @northes
- Kontexte von Asset-Anfragen unter macOS und iOS abbrechen, wenn WebKit die zugehörige Aufgabe des benutzerdefinierten Schemas abbricht (#5963)
- WindowSetFullscreenButtonEnabled-Nachricht verarbeiten, in [PR](https://github.com/wailsapp/wails/pull/5976) von @archy-rock3t-cloud
- Fragment in die preact-ts-Vorlage importieren, um den Buildfehler zu beheben, in [PR](https://github.com/wailsapp/wails/pull/5979) von @haoku123
- Abstürze älterer reiner GTK3-Dienstanwendungen verhindern, wenn die Bildschirmerkennung ausgeführt wird, bevor ein aktives Fenster oder Display verfügbar ist (#5966)
- Release-Läufe mit expliziter Version auch bei leerem Changelog für unveröffentlichte Änderungen fortsetzen lassen (#5977)

## Sicherheit

- nanoid-Lockfiles der Website auf die korrigierte Version 3.3.18 aktualisieren, um Sicherheitshinweise zu beheben, in [PR](https://github.com/wailsapp/wails/pull/5985) von @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Hinzugefügt

- Generierung von Dokumentations-URLs zu automatischen Changelog-Einträgen hinzufügen, in [PR](https://github.com/wailsapp/wails/pull/5957) von @taliesin-ai
- Streams hinzufügen: bidirektionale Byte-Streams zwischen Go und JavaScript mit dem WebSocket-Programmiermodell und ohne lauschenden Socket. Einen Stream in Go mit `app.HandleStream(name, handler)` deklarieren und vom Frontend mit `Stream(name)` verbinden; dies gibt ein Objekt in der Form von `WebSocket` zurück. Go→JS wird über den Asset-Server durch eine gehaltene Poll-Anfrage pro Fenster übertragen, JS→Go durch einen normalen POST; es wird weder ein TCP-Port gebunden noch etwas über `evaluateJavaScript` geleitet. In Server-Builds (`-tags server`) wird derselbe Handler stattdessen über einen echten WebSocket bereitgestellt, sodass der Anwendungscode in allen Builds identisch ist. Von @leaanthony
- Mailbox-Changelog-Eintrag nach „Unveröffentlicht“ verschieben, in [PR](https://github.com/wailsapp/wails/pull/5935) von @leaanthony

## Geändert

- Automatische Generierung der Dokumentationsseitenleiste und Ableitung des Blogautorentyps aktualisieren, in [PR](https://github.com/wailsapp/wails/pull/5938) von @leaanthony

## Behoben

- WebView2-Initialisierung verwendet eine Deadline und eine Nachrichtenschleife, in [PR](https://github.com/wailsapp/wails/pull/5952) von @leaanthony
- Der WebView2-Cookie-Test wird in CI übersprungen, sofern er nicht ausdrücklich aktiviert wurde, und bindet die Ausführung an den aktuellen Betriebssystem-Thread, in [PR](https://github.com/wailsapp/wails/pull/5951) von @leaanthony
- Windows-Menü-Builder stellen Befehls-IDs für übergeordnete Untermenüeinträge wieder her, in [PR](https://github.com/wailsapp/wails/pull/5944) von @gilad-ch
- Offizielles Cross-Compilation-Image an die Basis für GTK-4.14+-Linux-Unterstützung angleichen (#5928)
- iOS-Xcode-Projekt so konfigurieren, dass geerbte Linker-Flags beibehalten werden, und -ObjC hinzufügen, in [PR](https://github.com/wailsapp/wails/pull/5915) von @mortenolsrud
- Übermäßigen Wechsel von TCP-Verbindungen im `wails3 dev`-Asset-Proxy bei großen Frontends beheben, der die ephemeren Ports des Hosts erschöpfen und dazu führen konnte, dass nicht zugehörige Prozesse mit `EADDRNOTAVAIL` fehlschlugen
- Ereignis-JavaScript pro Fenster für geordnete Zustellung und Rückstaukontrolle in eine Warteschlange stellen, in [PR](https://github.com/wailsapp/wails/pull/5934) von @leaanthony

## Entfernt

- Release-Pipeline für Desktop-Binärdateien entfernen: v3-Releases bestehen nur aus Tags, und die `wails3`-CLI wird mit `go install` installiert. Entfernt `release-v3.yml` und den Nightly-Schritt, der sie ausgelöst hat, in [PR](https://github.com/wailsapp/wails/pull/5946) von @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Hinzugefügt

- macOS-Autoplay-Einstellung zum Deaktivieren der erforderlichen Benutzeraktion für die Medienwiedergabe hinzufügen, in [PR](https://github.com/wailsapp/wails/pull/5512) von @Eyalm321
- Mailbox-Changelog-Eintrag nach „Unveröffentlicht“ verschieben, in [PR](https://github.com/wailsapp/wails/pull/5935) von @leaanthony

## Geändert

- Die macOS-Zoomanimation verwendet CADisplayLink oder NSTimer für eine flüssigere Darstellung, in [PR](https://github.com/wailsapp/wails/pull/5945) von @savely-krasovsky

## Behoben

- iOS-Xcode-Projekt so konfigurieren, dass geerbte Linker-Flags beibehalten werden, und -ObjC hinzufügen, in [PR](https://github.com/wailsapp/wails/pull/5915) von @mortenolsrud
- Übermäßigen Wechsel von TCP-Verbindungen im `wails3 dev`-Asset-Proxy bei großen Frontends beheben, der die ephemeren Ports des Hosts erschöpfen und dazu führen konnte, dass nicht zugehörige Prozesse mit `EADDRNOTAVAIL` fehlschlugen
- Ereignis-JavaScript pro Fenster für geordnete Zustellung und Rückstaukontrolle in eine Warteschlange stellen, in [PR](https://github.com/wailsapp/wails/pull/5934) von @leaanthony

### Hinzugefügt

- Generische asynchrone FIFO-Mailbox für die geordnete Ereigniszustellung implementieren, in [PR](https://github.com/wailsapp/wails/pull/5851) von @savely-krasovsky und @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## Hinzugefügt

- Begrenzten hostseitigen Speicher für übergroße Ereignisse und geordnete JavaScript-Zustellung implementieren, in [PR](https://github.com/wailsapp/wails/pull/5930) von @leaanthony
- Springen des macOS-Dock-Symbols für das Aufblinken von Fenstern implementieren, in [PR](https://github.com/wailsapp/wails/pull/5921) von @julianstorer

## Behoben

- Der Asset-Server bewahrt beim Leeren Fehler der Inhaltstyperkennung und noch nicht geschriebene Präfixe, in [PR](https://github.com/wailsapp/wails/pull/5931) von @leaanthony
- Abstürze von macOS-Anwendungen verhindern, wenn das Anwendungsmenü aus einem Wails-Callback heraus ersetzt wird
- Unlesbare native Menüs unter Windows 10 1809 / Windows Server 2019 (Build 17763) beheben. Die uxtheme-Exporte für den Dunkelmodus waren auf Build 18334 beschränkt, sodass die Aktivierung des Dunkelmodus auf Anwendungsebene auf diesen Hosts nie ausgeführt wurde: Der Menühintergrund wurde dunkel gezeichnet, Windows zeichnete den Menütext jedoch weiterhin im hellen Design, wodurch dunkler Text auf dunklem Hintergrund erschien. Die Ordinalzahlen sind ab 17763 vorhanden; die Einschränkung stimmt nun damit überein.
- Beheben, dass `w32.GetStockObject` statt `GetStockObject` die Funktion `GetDeviceCaps` aufrief, wodurch für jedes Standardobjekt 0 zurückgegeben wurde.
- Fehlerbehandlung und -meldung beim Herunterladen des WebView2-Bootstrappers verbessern, in [PR](https://github.com/wailsapp/wails/pull/5924) von @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## Behoben

- Die Aktivierung von macOS-Apps berücksichtigt die Aktivierungsrichtlinie jetzt nur für reguläre Apps, in [PR](https://github.com/wailsapp/wails/pull/5897) von @julianstorer
- Nicht initialisierte GTK-Fenster in Linux-Builds absichern, in [PR](https://github.com/wailsapp/wails/pull/5898) von @julianstorer
- Vor dem Laden der URL eine explizite deckende Hintergrundfarbe für Linux-WebKit-Fenster festlegen, in [PR](https://github.com/wailsapp/wails/pull/5899) von @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## Geändert

- Android-Build-Tasks verwenden standardmäßig arm64, und deploy-emulator wählt die Hostarchitektur aus, in [PR](https://github.com/wailsapp/wails/pull/5890) von @mortenolsrud

## Behoben

- Zoomzustand von macOS-Fenstern beim Ziehen beibehalten und Bewegung reduzieren, in [PR](https://github.com/wailsapp/wails/pull/5900) von @leaanthony
- Build im Windows-Servermodus durch Hinzufügen von `!server` zur Build-Einschränkung `webview_window_windows_nonclient.go` reparieren

## v3.0.0-beta.3 - 2026-08-03

## Hinzugefügt

- Abschluss der Beta-Verifizierung für Phase 10 in den Implementierungsdetails dokumentieren, in [PR](https://github.com/wailsapp/wails/pull/5881) von @leaanthony

## Behoben

- Fenster-Handle an die Windows-Dunkelmodus-API übergeben und Argumente validieren, in [PR](https://github.com/wailsapp/wails/pull/5877) von @leaanthony
- Ermittlung des Zustands der Titelleistenschaltflächen für rahmenlose macOS-Fenster zentralisieren, in [PR](https://github.com/wailsapp/wails/pull/5870) von @taliesin-ai
- Unlesbaren Text in nativen Menüs verhindern, wenn eine Windows-Anwendung den Dunkelmodus anfordert, während das Windows-App-Design hell ist. Das Menü verwendet nun den passenden hellen nativen Hintergrund, bis Windows dunklen Menütext darstellen kann.
- Unlesbare native Menüs unter Windows 10 1809 / Windows Server 2019 (Build 17763) beheben. Die uxtheme-Exporte für den Dunkelmodus waren auf Build 18334 beschränkt, sodass die Aktivierung des Dunkelmodus auf Anwendungsebene auf diesen Hosts nie ausgeführt wurde: Der Menühintergrund wurde dunkel gezeichnet, Windows zeichnete den Menütext jedoch weiterhin im hellen Design, wodurch dunkler Text auf dunklem Hintergrund erschien. Die Ordinalzahlen sind ab 17763 vorhanden; die Einschränkung stimmt nun damit überein.

## v3.0.0-beta.2 - 2026-08-02

## Geändert

- v3 von Alpha auf Beta hochstufen
- Intelligente Standardwerte und das automatische Ausblenden von Pop-ups im Infobereich dokumentieren, einschließlich Regressionstests für die Auswahl des Klick-Handlers (#5840).
- Der GitHub-Updater schließt Windows-Installer-Assets standardmäßig aus, in [PR](https://github.com/wailsapp/wails/pull/5861) von @leaanthony
- Unterstützung für abgerundete, eckige und benutzerdefinierte Eckenradien bei rahmenlosen macOS-Fenstern hinzufügen, in [PR](https://github.com/wailsapp/wails/pull/5866) von @leaanthony

## Behoben

- Aktuelle Größen von GTK4-Fenstern melden und Ereignisse für Größenänderung, Maximierung, Minimierung und Vollbildstatus von der konfigurierten Oberfläche ausgeben (#5830).
- Absturz von Linux WebKit beim Senden von Blob oder FormData in fetch-Anfragen beheben, in [PR](https://github.com/wailsapp/wails/pull/5854) von @taliesin-ai
- Der Fetch-Shim übergibt bei fehlenden Blob-/FormData-Headern undefined, in [PR](https://github.com/wailsapp/wails/pull/5865) von @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## Hinzugefügt

## Geändert

- Unterstützung für abgerundete, eckige und benutzerdefinierte Eckenradien bei rahmenlosen macOS-Fenstern hinzufügen, in [PR](https://github.com/wailsapp/wails/pull/5866) von @leaanthony

## Behoben

- Der Fetch-Shim übergibt bei fehlenden Blob-/FormData-Headern undefined in [PR](https://github.com/wailsapp/wails/pull/5865) von @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## Hinzugefügt

- Unterstützung für die macOS-DMG-Paketierung mit neuen Optionen und Build-Tasks hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5857) von @leaanthony

## Geändert

- Der GitHub-Updater schließt Windows-Installer-Assets standardmäßig aus in [PR](https://github.com/wailsapp/wails/pull/5861) von @leaanthony

## Behoben

- Linux-WebKit-Absturz beim Senden von Blob oder FormData in Fetch-Anfragen behoben in [PR](https://github.com/wailsapp/wails/pull/5854) von @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## Hinzugefügt

- Doppelklickaktionen auf die macOS-Titelleiste zum Maximieren oder Minimieren implementiert in [PR](https://github.com/wailsapp/wails/pull/5853) von @taliesin-ai

## Geändert

- Tutorial zum QR-Dienst auf NewServiceWithOptions umgestellt und Abstände ergänzt in [PR](https://github.com/wailsapp/wails/pull/5849) von @jeongkyu

## Behoben

- Reaktionsfähigkeit von WKWebView während des Zoomens unter macOS sichergestellt in [PR](https://github.com/wailsapp/wails/pull/5856) von @leaanthony
- Abfragen der GTK4-Fenstergröße korrigiert und Ereignisse für Größenänderung, Maximierung, Minimierung und Vollbildstatus über das konfigurierte `GdkSurface` ausgegeben.

## v3.0.0-alpha2.119 - 2026-07-27

## Behoben

- Mehrsprachige Dokumentation um Architekturdiagramme ergänzt in [PR](https://github.com/wailsapp/wails/pull/5833) von @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## Hinzugefügt

- Standardpfade für Ein- und Ausgaben der Symbolgenerierung bereitgestellt in [PR](https://github.com/wailsapp/wails/pull/5825) von @taliesin-ai
- Quell-Einstiegsmodule zu sideEffects in der runtime-package.json hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5797) von @savely-krasovsky
- Abschnitt zu Lizenz und Herkunftsnachweisen zum Beitragsleitfaden hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5816) von @taliesin-ai

## Behoben

- Bereichsbegrenztes CSS für rahmenlose GTK4-Fenster angewendet, um den Eckenradius zu entfernen, in [PR](https://github.com/wailsapp/wails/pull/5800) von @savely-krasovsky
- Fehler beim Ermitteln der Cursorposition unter Windows für Popup-Menüs und die Bildschirmauflistung zuverlässig behandelt in [PR](https://github.com/wailsapp/wails/pull/5789) von @wayneforrest
- Der macOS-Dialog zum Öffnen von Dateien filtert Erweiterungen korrekt und validiert zulässige Dateien anhand ihrer Endung in [PR](https://github.com/wailsapp/wails/pull/5678) von @phergul
- Initialisierung des Dunkelmodus unter Windows gegen nil-API-Aufrufe abgesichert in [PR](https://github.com/wailsapp/wails/pull/5793) von @roachadam
- Fehler beim 32-Bit-Build des Updaters behoben: Die Konstante `maxArchiveTotalSize` (2 GiB) führte beim Übergeben an `fmt.Errorf` unter `GOARCH=386` zu einem Überlauf des plattformspezifischen `int`. Sie ist jetzt explizit als `int64` typisiert.
- Eine nil-Zeiger-Panik beim Start behoben, wenn ein Fenster eine dunkle (oder systembedingt dunkle) Titelleiste in Windows-Builds verwendet, welche die uxtheme-APIs für den Dunkelmodus nicht laden, etwa Windows 10 1809 / Windows Server 2019 (Build 17763). Die `AllowDarkModeForWindow`-Aufrufe beim Einrichten des Fensterdesigns sind jetzt gegen nil abgesichert, entsprechend der bereits in `w32.SetMenuTheme` verwendeten Absicherung.

## v3.0.0-alpha2.117 - 2026-07-08

## Hinzugefügt

- Benutzerdefinierte Hit-Test-Logik für Nicht-Client-Bereiche unter Windows implementiert in [PR](https://github.com/wailsapp/wails/pull/5462) von @savely-krasovsky

## Geändert

- WebView2-Erkennung der Monitorskalierung anhand von UseVisualHosting konfiguriert in [PR](https://github.com/wailsapp/wails/pull/5761) von @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## Hinzugefügt

- FAQ-Dokumentation auf Funktionen und Anleitungen für Wails v3 ausgerichtet in [PR](https://github.com/wailsapp/wails/pull/5763) von @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## Behoben

- Fehler behoben, durch den `Menu.Update()` das native Menü unter GTK4 Linux nicht neu aufbaute (#5659, unabhängig von @puneetdixit200 in #5539 diagnostiziert und behoben)
- Absturz beim Auflisten von macOS-Bildschirmen nach einer Anzeigeänderung behoben, indem ID-/Namenszeichenfolgen der Bildschirme kopiert und die Anzahl als Momentaufnahme erfasst wird (#5565, unabhängig von @x-haose in #5584 diagnostiziert und behoben)
- Absturz unter Windows behoben, wenn `WM_ERASEBKGND` während eines Minimierungs-/Wiederherstellungsübergangs einen einfarbigen Hintergrund zeichnet und `GetClientRect` dabei nil zurückgibt (Absicherung von @sinspired in #5636 gemeldet)
- Fehler behoben, durch den Fehler aus Frontend-Bindings stets als Text geparst wurden, von @mbaklor in #5690
- Fehler beim Windows-Build mit dem Build-Tag `server` behoben. Ursache waren Windows-GUI-Dateien, denen die bereits bei den macOS- und Linux-Entsprechungen vorhandene Build-Bedingung `!server` fehlte (#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## Hinzugefügt

- Update-Manifest-Protokoll und Endpunktanbieter implementiert in [PR](https://github.com/wailsapp/wails/pull/5720) von @taliesin-ai

## Geändert

- Das Binding `webview2` als `v3/internal/webview2` in das Modul v3 integriert und dabei das eigenständige Modul, dessen Workflows für nächtliche Releases und Synchronisierung sowie den go.mod-Versionstrick entfernt (v3 ist sein einziger Verbraucher) in [PR](https://github.com/wailsapp/wails/pull/5711) von @taliesin-ai

## Behoben

- WebView2-Erkennung der Monitorskalierung und Korrektur der Host-Neusynchronisierung bei DPI-Änderungen in den Abschnitt „Unveröffentlicht“ verschoben in [PR](https://github.com/wailsapp/wails/pull/5750) von @taliesin-ai
- WebView2-COM-Marshalling für float64- und BOOL-Parameter aktualisiert in [PR](https://github.com/wailsapp/wails/pull/5741) von @wayneforrest
- Panik und nil-Dereferenzierung beim Aktualisieren und Zerstören von Taskleistensymbolen unter Windows verhindert in [PR](https://github.com/wailsapp/wails/pull/5703) von @wayneforrest
- Fehler behoben, durch den ausgeblendete Fenster unter Windows nicht korrekt erneut ausgeblendet wurden, in [PR](https://github.com/wailsapp/wails/pull/5743) von @wayneforrest
- Sichtbarkeit des WebView2-Controllers mit dem Minimieren, Maximieren und Wiederherstellen des Fensters synchronisiert in [PR](https://github.com/wailsapp/wails/pull/5742) von @wayneforrest

### Behoben

- WebView2-Erkennung der Monitorskalierung wieder aktiviert und die Host-Neusynchronisierung auf DPI-Änderungen beschränkt in [PR](https://github.com/wailsapp/wails/pull/5734) von @taliesin-ai, basierend auf der von @randalmurphal validierten Korrektur, mit Überprüfung der Grundursache durch @eleclin und Hardwaretests durch @qq540491950

## v3.0.0-alpha2.113 - 2026-07-04

## Hinzugefügt

- Warnung beim Erstellen eines Release-AAB ohne gesetztes `ANDROID_KEYSTORE_FILE` hinzugefügt (Google Play lehnt mit einem Debug-Schlüssel signierte Bundles ab) und die Paketierung sowie Signierung von App Bundles dokumentiert in [PR](https://github.com/wailsapp/wails/pull/5730) von @taliesin-ai
- Mehrsprachige Dokumentation „Warum Wails?“ hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5739) von @taliesin-ai
- Unterstützung für die Zuordnung von Go time.Time zu JS Date oder string in Bindings hinzugefügt in [PR](https://github.com/wailsapp/wails/pull/5398) von @fbbdev
- Fügt Tasks zum Paketieren von Android App Bundles (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) für die Einreichung im Play Store hinzu; die APK-Tasks bleiben für lokale Tests und Emulatortests erhalten, in [PR](https://github.com/wailsapp/wails/pull/5728) von @mortenolsrud (behebt [#5726](https://github.com/wailsapp/wails/issues/5726))
- Fügt Android-Task-Targets für physische Geräte hinzu und stellt Kamera- und Standortberechtigungen wieder her, in [PR](https://github.com/wailsapp/wails/pull/5735) von @taliesin-ai

## Geändert

- Aktualisiert `webview2` auf v1.0.28 ([Versionshinweise](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- Aktualisiert im Android-Template `compileSdk`/`targetSdk` von 34 auf 35, wie von Google Play für die Einreichung neuer Apps vorgeschrieben, in [PR](https://github.com/wailsapp/wails/pull/5730) von @taliesin-ai

## Behoben

- Behebt fest eingebettete Avatar-Masken in sponsorkit, in [PR](https://github.com/wailsapp/wails/pull/5745) von @leaanthony
- Behebt, dass die automatische Erstellung eines Android-AVD aufgrund lexikografischer Versionssortierung das falsche System-Image oder die falsche cmdline-tools-Version auswählt, in [PR](https://github.com/wailsapp/wails/pull/5730) von @taliesin-ai
- Behebt, dass der Einrichtungsassistent eine veraltete Android-NDK-Version vorschlägt (jetzt 26.3.11579264, entsprechend der dokumentierten Anforderung), in [PR](https://github.com/wailsapp/wails/pull/5730) von @taliesin-ai
- Aktualisiert die französische Dokumentation für SvelteKit und Optionen, in [PR](https://github.com/wailsapp/wails/pull/5744) von @leaanthony
- Behebt SIGSEGV bei der Auflistung von macOS-Bildschirmen während Anzeigeänderungen, in [PR](https://github.com/wailsapp/wails/pull/5516) von @flofreud

## v3.0.0-alpha2.112 - 2026-07-03

## Hinzugefügt

- Fügt einen Go-basierten SVG-Generator für Mitwirkende hinzu und aktualisiert die Danksagungsseiten in der Dokumentation und auf der Website, in [PR](https://github.com/wailsapp/wails/pull/5724) von @taliesin-ai
- Unterstützt in Bindings die Abbildung von Go time.Time auf JS Date oder string, in [PR](https://github.com/wailsapp/wails/pull/5398) von @fbbdev

## Geändert

- Ersetzt die Node-basierte Pipeline für Sponsorenbilder durch einen Go-Generator, in [PR](https://github.com/wailsapp/wails/pull/5719) von @taliesin-ai

## Behoben

- Behebt das Installationsskript für Abhängigkeiten der Android-Build-Assets, in [PR](https://github.com/wailsapp/wails/pull/5729) von @taliesin-ai
- Lehnt das Steuerzeichen U+0085 (NEXT LINE) in `ValidateAndSanitizeURL` ab und vervollständigt damit die Abdeckung von Leerraumzeichen im URL-Validator
- Berechnet den DWM-Rahmen bei DPI-Änderungen für rahmenlose Fenster neu, in [PR](https://github.com/wailsapp/wails/pull/4785) von @leaanthony
- Behebt die fehlschlagende Erkennung von DnD-Ablagezonen unter Windows bei einer Skalierung ungleich 100 %, in [PR](https://github.com/wailsapp/wails/pull/4632) von @yulesxoxo
- Fügt explizite Objective-C-Speicherverwaltung für Cocoa-Objekte in Darwin-Dialogen, -Menüs, -Taskleistensymbolen und -Benachrichtigungen hinzu, in [PR](https://github.com/wailsapp/wails/pull/5714) von @taliesin-ai
- Behebt Fehler im Linux-CGO-Backend und Probleme mit dem Taskleistensymbol, in [PR](https://github.com/wailsapp/wails/pull/5718) von @taliesin-ai

## v3.0.0-alpha2.111 - 2026-07-01

## Hinzugefügt

- Fügt HappyTools zur Community-Galerie hinzu, in [PR](https://github.com/wailsapp/wails/pull/5061) von @Aliuyanfeng
- Fügt Unterstützung für die indonesische Sprache und eine umfassende Dokumentation hinzu, in [PR](https://github.com/wailsapp/wails/pull/5643) von @triadmoko
- Fügt WindowsWindow die Option DisableMenu hinzu, in [PR](https://github.com/wailsapp/wails/pull/4813) von @leaanthony

## Geändert

- Aktualisiert das Taskfile-Template und die CLI, um Build- und Paketierungs-Tasks mit GOOS und ARCH auszuführen, in [PR](https://github.com/wailsapp/wails/pull/5617) von @leaanthony

## Behoben

- Behebt ein Problem mit Fenster-Tabs unter macOS, in [PR](https://github.com/wailsapp/wails/pull/5708) von @taliesin-ai

## Entfernt

- Entfernt die übersetzten deutschen MDX-Dateien aus den Bereichen Mitwirken, Funktionen und Anleitungen, in [PR](https://github.com/wailsapp/wails/pull/5702) von @taliesin-ai

## v3.0.0-alpha2.110 - 2026-06-30

## Hinzugefügt

- Implementiert für macOS das Neuladen und erzwungene Neuladen der WebView und fügt eine Wiederherstellung nach Beendigung des WebContent-Prozesses hinzu, in [PR](https://github.com/wailsapp/wails/pull/5129) von @wayneforrest
- Fügt eine umfassende deutsche Dokumentation für die Bereiche Mitwirken, Funktionen und Anleitungen hinzu, in [PR](https://github.com/wailsapp/wails/pull/5396) von @leaanthony
- Erweitert Benachrichtigungen um Töne, Anhänge und Zeitplanung und aktualisiert die API, in [PR](https://github.com/wailsapp/wails/pull/5333) von @popaprozac

## Behoben

- Berechnet den DWM-Rahmen bei DPI-Änderungen für rahmenlose Fenster neu, in [PR](https://github.com/wailsapp/wails/pull/4785) von @leaanthony
- Behebt die fehlschlagende Erkennung von DnD-Ablagezonen unter Windows bei einer Skalierung ungleich 100 %, in [PR](https://github.com/wailsapp/wails/pull/4632) von @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Hinzugefügt

- Fügt der EventsEmit-Dokumentation Codebeispiele hinzu, in [PR](https://github.com/wailsapp/wails/pull/5026) von @iamhabbeboy
- Fügt eine Option für visuelles WebView2-Hosting unter Windows hinzu, in [PR](https://github.com/wailsapp/wails/pull/5380) von @MerIijn
- Fügt Klustr zur Dokumentation der Community-Galerie hinzu, in [PR](https://github.com/wailsapp/wails/pull/5536) von @SametKUM
- Fügt Kira mit neuen Seiten und einem Changelog-Eintrag zur Community-Galerie hinzu, in [PR](https://github.com/wailsapp/wails/pull/5685) von @thiennguyen93
- Fügt dem Leitfaden zum MCP-Dienst einen Feedbackabschnitt hinzu, in [PR](https://github.com/wailsapp/wails/pull/5694) von @taliesin-ai

## Geändert

- Der Servermodus verfügt jetzt über einen vollwertigen Produktions-Build, entsprechend den Desktop-Build-Tasks (#5693). `task build:server` erstellt standardmäßig ein Produktions-Binary (`-tags server,production`, `-trimpath`, bereinigt) und akzeptiert `DEV=true` (Entwicklungsserver), `OBFUSCATED=true` (garble) und `EXTRA_TAGS`. `task run:server` führt einen Entwicklungsserver aus. `Dockerfile.server` / `task build:docker` erstellt zuerst den Produktionsserver (`-tags server,production`) und das Produktions-Frontend. Das Image verwendet standardmäßig einen statischen Pure-Go-Build auf distroless/static; `CGO_ENABLED`, `GO_IMAGE` und `RUNTIME_IMAGE` stehen für CGO-Apps als überschreibbare Build-Argumente zur Verfügung.

## Behoben

- Verhindert einen Absturz beim Schließen eines Fensters mit ausstehenden asynchronen Aufrufen, in [PR](https://github.com/wailsapp/wails/pull/4435) von @leaanthony
- Verhindert unter Windows die Fensteraktivierung beim Öffnen ausgeblendeter Apps, in [PR](https://github.com/wailsapp/wails/pull/5249) von @leaanthony
- Stellt sicher, dass die Verarbeitung von WebKit-Anfragemetadaten, der Abschluss von Antworten und die Verarbeitung des Body-Streams im GTK-Hauptthread erfolgen, in [PR](https://github.com/wailsapp/wails/pull/5668) von @taliesin-ai
- Behebt, dass `Menu.Update()` unter GTK4 Linux das native Menü nicht neu erstellt (#5659; unabhängig von @puneetdixit200 in #5539 diagnostiziert und behoben)
- Behebt einen Absturz beim Auflisten von macOS-Bildschirmen nach einer Anzeigeänderung, indem die Zeichenfolgen für Bildschirm-ID und -Name kopiert und die Anzahl als Momentaufnahme erfasst werden (#5565; unabhängig von @x-haose in #5584 diagnostiziert und behoben)
- Behebt, dass WebView2-Inhalte unter Windows nach dem Ziehen eines Fensters über Monitore mit unterschiedlichen DPI-Einstellungen zunächst schrumpfen und dann verschwinden. Dazu werden die Controller-Grenzen im `WM_DPICHANGED`-Handler erneut festgelegt, entsprechend der DPI-Neusynchronisierung beim Aufheben der Minimierung (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Hinzugefügt

- Globale (systemweite) Tastenkürzel über `app.GlobalShortcut` hinzufügen (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Die Tastenkürzel werden auch ausgelöst, wenn die Anwendung nicht fokussiert ist. Sie sind ohne Drittanbieterabhängigkeiten nativ für jede Plattform implementiert: Carbon-Hotkeys unter macOS, `RegisterHotKey` unter Windows, `XGrabKey` unter X11 und die Schnittstelle für globale Tastenkürzel des XDG Desktop Portal unter Wayland.
- Integrierten MCP-Server hinzufügen: einen Model-Context-Protocol-Server, der automatisch startet, wenn die Anwendung mit dem Tag `mcp` erstellt wurde. Dadurch können LLM-Agenten eine laufende Wails-Anwendung testen und steuern – einschließlich Fenstersteuerung, DOM-Inspektion, JavaScript-Auswertung, Aufrufen gebundener Methoden, Ereignissen und simulierten Maus- und Tastatureingaben, die mit einem animierten Bildschirmcursor dargestellt werden. Kein Benutzercode erforderlich: `wails3 build`/`wails3 dev` fügt das Tag `mcp` automatisch hinzu, wenn `WAILS_MCP=1` gesetzt ist. Die Konfiguration erfolgt vollständig über Umgebungsvariablen (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Behoben

- Beheben, dass `Menu.Update()` unter GTK4/Linux das native Menü nicht neu erstellt (#5659, unabhängig von @puneetdixit200 in #5539 diagnostiziert und behoben)
- Absturz beim Auflisten von macOS-Bildschirmen nach einer Anzeigeänderung beheben, indem die Zeichenfolgen für Bildschirm-ID und -Name kopiert und die Anzahl als Momentaufnahme erfasst wird (#5565, unabhängig von @x-haose in #5584 diagnostiziert und behoben)

## v3.0.0-alpha2.107 - 2026-06-27

## Hinzugefügt

- Experimentelle Wake-Dokumentation mit Seitenleistennavigation in [PR](https://github.com/wailsapp/wails/pull/5613) von @leaanthony hinzufügen

## v3.0.0-alpha2.106 - 2026-06-24

## Geändert

- `webview2` auf v1.0.27 aktualisieren.
  - ci(webview2): Release-Build korrigieren (Windows-Cross-Kompilierung + vollständige go.sum) (#5671)\

  **Vollständiger Diff:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- go vet in [PR](https://github.com/wailsapp/wails/pull/5672) von @taliesin-ai aus der Cross-Kompilierung des webview2-Release-Workflows entfernen
- OpenRouter-Modell des automatischen Changelogs in [PR](https://github.com/wailsapp/wails/pull/5670) von @taliesin-ai auf google/gemini-2.5-flash-lite aktualisieren
- `webview2` auf v1.0.26 aktualisieren.

### Korrekturen

- **Nach vorübergehenden COM-Laufzeitfehlern wiederherstellen, statt die Anwendung zu beenden** (#5658, #5580). Zuvor rief `Chromium.errorCallback` bei *jedem* COM-Fehler `os.Exit(1)` auf, sodass eine behebbare Störung nach dem Start die gesamte Anwendung beendete. Laufzeitpfade (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) protokollieren den Fehler nun und stellen den Betrieb wieder her. Insbesondere wird eine fehlerhafte oder nicht vertrauenswürdige Webnachricht in `MessageReceived` nun verworfen, statt den Prozess zu beenden. Dies behebt die Absturzklasse beim Wechsel zwischen Monitoren mit unterschiedlicher DPI (#5544, #5650). Pfade zum Erstellen der Umgebung bzw. des Controllers bleiben fatal.\

**Vollständiger Diff:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Behoben

- release-webview2-Workflow in [PR](https://github.com/wailsapp/wails/pull/5671) von @taliesin-ai korrigieren, damit go.sum-Dateien ordnungsgemäß verarbeitet werden
- Aktualisierungen des GTK4-Menüs unter Linux in [PR](https://github.com/wailsapp/wails/pull/5659) von @taliesin-ai korrigieren, indem das native Menü geleert und neu erstellt wird

## v3.0.0-alpha2.105 - 2026-06-21

## Hinzugefügt

- `application.System` zur Laufzeiterkennung der Plattform aus gemeinsam genutztem Code hinzufügen: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (das Build-Tag `server`) und `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)`, um ein einzelnes Ziel direkt zu prüfen. Es lässt sich für jedes Ziel kompilieren, sodass Verzweigungen ohne Build-Tags möglich sind. Entsprechende Frontend-Hilfsfunktionen (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) sind in `@wailsio/runtime` verfügbar
- Leitfaden „Andere Frontend-Frameworks verwenden“ hinzufügen, der zeigt, wie ein eigenes Vite-Projekt in `frontend/` eingebunden wird (behandelt Solid, Preact, Lit, SvelteKit, Qwik, Angular usw.)
- Der Assistent `wails3 setup` prüft nun die Toolchain für Mobilgeräte (iOS/Android) – Xcode und die iOS-Simulator-Laufzeit, JDK, Android SDK/NDK und Emulator – und bietet, sofern möglich, eine Installation mit einem Klick sowie kopierbare Korrekturen für die Shell-Konfiguration
- Generierte Projekte enthalten eine `frontend/.npmrc`, die eine `minimum-release-age` von 7 Tagen festlegt, um die Gefährdung durch neu veröffentlichte (und möglicherweise kompromittierte) Pakete zu reduzieren (wird von pnpm und bun berücksichtigt und von npm ohne Auswirkungen ignoriert)

## Geändert

- Alle integrierten Startervorlagen mit einem neuen Hero-Design aus Neonbergen neu gestalten (Web, iOS und Android)
- **TypeScript ist nun die Standardeinstellung für Startervorlagen und übernimmt den unveränderten Vorlagennamen.** `wails3 init` (ohne `-t`) legt ein TypeScript-Projekt an; `-t vanilla`, `-t react`, `-t vue` und `-t svelte` verwenden TypeScript. JavaScript-Varianten stehen unter `-t vanilla-js`, `-t react-js`, `-t vue-js` und `-t svelte-js` bereit. Integrierte Vorlagen deklarieren ihre Sprache mit `typescript:` in `template.yaml`; Community-Vorlagen mit dem Suffix `-ts` funktionieren weiterhin als Fallback
- Den Assistenten `wails3 setup` mit dem Neon-Design „digital Wails“ neu gestalten (Milchglas-Vibranz vor einer Bergkulisse)

## Behoben

- Absturz unter Windows beim Wiederherstellen einer Anwendung beheben, die so lange minimiert war, dass WebView2 angehalten oder dessen Render-/GPU-Prozess wiederverwendet wurde. Die DPI-Neusynchronisierung beim Minimieren und Wiederherstellen (#5544) greift nun nur auf den WebView2-Controller zu, wenn sich die DPI des Fensters tatsächlich geändert hat. Dadurch werden beim üblichen Wiederherstellen mit unveränderter DPI fatale COM-Aufrufe an einen angehaltenen Controller vermieden (#5605)
- Wiederholte native Abstürze in `SIGABRT`/`SIGSEGV` beheben, die bei lange laufenden Linux-Anwendungen unter häufigem Laden von Assets und Medien auftreten (typischerweise während der GTK-Hauptschleife in `g_object_unref`). Der Asset-Server schloss `WebKitURISchemeRequest`s aus Worker-Goroutinen ab und rief dadurch nicht threadsichere WebKit2GTK-Funktionen außerhalb des GTK-Hauptthreads auf; der Abschluss (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) erfolgt nun im Hauptthread. Vervollständigt die teilweise Korrektur aus #5566. Betrifft sowohl GTK3- als auch GTK4/WebKitGTK-Builds mit 6.0 (#5631, #5557)
- Sporadisches `fatal error: invalid pointer found on stack` in `setupSignalHandlers` unter Linux/GTK3 beheben. Als Signal-`user_data` übergebene Fenster-IDs wurden in einer lokalen Go-Variablen vom Typ `unsafe.Pointer` gehalten. Daher brach der Garbage Collector ab, als er den Wert, der kein Zeiger war, während einer Stack-Kopie untersuchte. Auf der Go-Seite wird die ID nun als Ganzzahl (`uintptr_t`) gespeichert. Damit wird dieselbe Korrektur, die #4958 auf den GTK4-Pfad angewendet hatte, auf den alten GTK3-Pfad zurückportiert (dort wurden die C-Signalfunktionen auf `uintptr_t` umgestellt, um `-race`-/checkptr-Fehler zu beseitigen) (#5631)

## Entfernt

- Die Startervorlagen `react-swc`, `preact`, `lit`, `solid`, `qwik` und `sveltekit` sowie deren `-ts`-Varianten entfernen. Die unterstützten integrierten Vorlagen sind nun `vanilla`, `react`, `vue` und `svelte` – jeweils standardmäßig mit TypeScript und mit JavaScript-Varianten unter `-js`. Jedes andere Framework kann weiterhin verwendet werden, indem Sie [Ihr eigenes Frontend mitbringen](https://v3.wails.io/guides/dev/frontend-frameworks) oder eine benutzerdefinierte Vorlage verwenden

## v3.0.0-alpha2.104 - 2026-06-18

## Behoben

- iOS-Absturz (SIGABRT) beheben, wenn eine gebundene Go-Dienstmethode eine leere Zeichenfolge zurückgibt. Der Writer für iOS-Asset-Antworten prüfte den Body-Zeiger mit `buf != nil` statt dessen Länge, sodass ein Body der Länge null in `&buf[0]` eine Panic auslöste. Nun wird die Länge geprüft, entsprechend den Desktop-Writern

## v3.0.0-alpha2.103 - 2026-06-15

## Geändert

- Native iOS- und Android-Funktionen in Plattformmanager verschieben: Statt der bisherigen freien Funktionen `application.IOS*`/`application.Android*` werden sie über `application.IOS.*` und `application.Android.*` aufgerufen (z. B. `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`) (#5602)
- Mobile-Bridge-Ereignisse umbenannt: Plattformübergreifende Ereignisse verwenden jetzt das Präfix `common:*` (z. B. `common:haptic`, `common:location`), plattformspezifische Ereignisse verwenden `ios:*` bzw. `android:*` (z. B. `ios:backgroundTask`, `android:foregroundService`); das Präfix `native:*` wird nicht mehr verwendet (#5602)

## v3.0.0-alpha.102 – 2026-06-14

## Hinzugefügt

- Experimentellen Assistenten `wails3 setup` zur interaktiven Projekteinrichtung und Abhängigkeitsprüfung hinzugefügt
- Flag `--json` für maschinenlesbare Ausgaben zu `wails3 doctor` hinzugefügt
- Abschnitt zum Signierungsstatus zum Befehl `wails3 doctor` hinzugefügt

## Behoben

- npm-Erkennung unter Linux korrigiert, sodass zusätzlich zum Paketmanager auch PATH geprüft wird

## v3.0.0-alpha.101 – 2026-06-13

## Hinzugefügt

- iOS: native Meldungsdialoge (UIAlertController) sowie Dialoge zum Öffnen einer Datei, mehrerer Dateien oder eines Verzeichnisses (UIDocumentPickerViewController); Speicherdialoge geben einen expliziten Fehler zurück
- iOS: Unterstützung der Zwischenablage über UIPasteboard
- iOS: reale Bildschirmmetriken über UIScreen (Punkte, Pixel, Skalierung und um sichere Bereiche bereinigter Arbeitsbereich)
- iOS: Geräte-Builds (`IOS_PLATFORM=device`), Unterstützung für Codesignierungsidentität, Bereitstellungsprofil und Berechtigungen, Paketierung als `.ipa` sowie `deploy-device` über devicectl
- iOS: konfigurierbare iOS-Mindestversion (`ios.minIOSVersion` in build/config.yml)
- iOS: `wails3 doctor` meldet unter macOS die Verfügbarkeit von Xcode und dem iOS SDK
- iOS: Systemereignisse für Akku, Netzwerk, Theme, Bildschirmsperre und wenig Arbeitsspeicher werden als `events.IOS.*`- und plattformneutrale `events.Common.*`-Anwendungsereignisse bereitgestellt
- iOS: native Bridge für mobile Funktionen (exportiertes `application.IOS*`) – Teilen-Dialog, Öffnen von URLs, Verhindern des Ruhezustands, Taschenlampe, Abstände sicherer Bereiche, Helligkeit, App-Informationen, Ausrichtungssperre, Statusleiste, Biometrie (Face ID/Touch ID), lokale Benachrichtigungen und sicherer Speicher im Keychain
- iOS: Sensoren und Hardware – Haptik, einmalige Standortbestimmung, Beschleunigungssensor, Näherungssensor, Text-to-Speech, Speicherinformationen, Stromversorgungs-/Akkustatus, Netzwerkstatus, Tastaturabstände und Erkennung von Bildschirmaufnahmen
- iOS: Dokumentation (IOS.md und ein Leitfaden auf der Dokumentationswebsite)
- Android: native Meldungsdialoge (AlertDialog) sowie Dialoge zum Öffnen einer oder mehrerer Dateien (Storage Access Framework; als Kopien in den Cache importiert); Dialoge zum Öffnen eines Verzeichnisses und zum Speichern geben einen expliziten Fehler zurück
- Android: Unterstützung der Zwischenablage über ClipboardManager
- Android: reale Bildschirmmetriken über WindowMetrics/DisplayMetrics (dp, Pixel, Skalierung und um Systemleisten bereinigter Arbeitsbereich)
- Android: Laufzeitmethoden für Haptik (`Android.Haptics.Vibrate`), Geräteinformationen (`Android.Device.Info`) und Toast-Benachrichtigungen (`Android.Toast.Show`)
- Android: typisierte Lebenszyklusereignisse (`events.Android.*`, aus events.txt generiert), wobei `ActivityCreated` auf `Common.ApplicationStarted` abgebildet wird
- Android: Die Build-Pipeline erzeugt installierbare Debug- und Release-APKs (`android:run`, `android:package`, `android:package:fat`); Release-Signierung standardmäßig mit dem Debug-Keystore oder mit einem echten Keystore über `ANDROID_KEYSTORE_*`-Umgebungsvariablen
- Android: `wails3 doctor` meldet Android SDK, NDK und JDK
- Android: Systemereignisse für Akku, Netzwerk, Theme, Bildschirmsperre und wenig Arbeitsspeicher werden als `events.Android.*`- und plattformneutrale `events.Common.*`-Anwendungsereignisse bereitgestellt
- Android: native Bridge für mobile Funktionen (exportiertes `application.Android*`) – Teilen, Öffnen von URLs, Verhindern des Ruhezustands, Taschenlampe, Abstände sicherer Bereiche, Helligkeit, App-Informationen, Ausrichtungssperre, Statusleiste, Biometrie (BiometricPrompt), lokale Benachrichtigungen und sicherer Speicher in EncryptedSharedPreferences
- Android: Sensoren und Hardware – Haptik, einmalige Standortbestimmung, Beschleunigungssensor, Näherungssensor, Text-to-Speech, Speicherinformationen, Stromversorgungs-/Akkustatus, Netzwerkstatus, Tastaturabstände und Blockierung von Bildschirmaufnahmen durch FLAG_SECURE
- Android: Dokumentation (ANDROID.md und ein Leitfaden auf der Dokumentationswebsite)
- Beispiel: Das Kitchen-Sink-Beispiel `mobile` erhält die Registerkarten „Mobile“ und „Hardware“, die die native Funktions-Bridge für iOS und Android demonstrieren (pill-förmige Registerkarten werden auf mehrere Zeilen umgebrochen)
- Mobile: Akku – Beschleunigungs- und Näherungssensor, Taschenlampe sowie die periodische Uhr des Beispiels werden pausiert, wenn die App in den Hintergrund wechselt, und bei der Rückkehr wiederhergestellt (Android lässt den Prozess im Hintergrund weiterlaufen, und unter iOS bleibt der Hardwarezustand der Taschenlampe erhalten); Android-Empfänger für Systemereignisse werden nur registriert, solange sich die App im Vordergrund befindet
- iOS: Kameraaufnahme – `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → ein `native:capture`-Ereignis mit einem Base64-Vorschaubild)
- iOS: Hintergrundausführung – `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (ein UIApplication-Hintergrundaufgabenfenster) und ein konfigurierbares `ios.backgroundModes` (build/config.yml), das `UIBackgroundModes` als Vorlage in die generierte Info.plist einfügt
- Android: Kameraaufnahme – `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (Systemkamera über FileProvider → ein `native:capture`-Ereignis)
- Android: Vordergrunddienst – `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (ein `WailsForegroundService` mit einer dauerhaften Benachrichtigung hält den Prozess für lang laufende Hintergrundaufgaben aktiv)
- Beispiel: Eine Kamera-Registerkarte demonstriert Foto-/Videoaufnahmen und die Hintergrundausführung (Vordergrunddienst unter Android, Hintergrundaufgabenfenster unter iOS)

## Behoben

- Behoben, dass `getUserMedia` unter Linux immer mit `NotAllowedError` fehlschlug: WebKitGTK lehnt Berechtigungsanfragen ab, die von niemandem verarbeitet werden, und das Signal `permission-request` war nicht verbunden. Kamera und Mikrofon werden jetzt anhand einer neuen plattformübergreifenden `WebviewWindowOptions.Permissions`-Zuordnung (`map[PermissionType]Permission`) verarbeitet, die sowohl unter Linux (WebKitGTK) als auch unter Windows (WebView2) berücksichtigt wird. Unter Linux, wo es keine native Abfrage gibt, sind Kamera und Mikrofon standardmäßig zugelassen (wodurch `getUserMedia` wiederhergestellt wird) und können mit `PermissionDeny` deaktiviert werden (#5552)
- iOS: `GOOS=ios` lässt sich wieder kompilieren (exportiertes `events.IOS`, Stubs für Namen mobiler Methoden), und mit dem production-Tag versehene Builds lassen sich kompilieren (Korrekturen der Build-Tags in pkg/application und mehreren Diensten)
- iOS: Go→JS-Ereignisse und ExecJS funktionieren jetzt – die Seite wird beim Start nicht mehr zweimal geladen, und der `wails:runtime:ready`-Handshake kann nicht mehr verloren gehen
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` konkurrieren nicht mehr mit dem App-Start; die feste Startverzögerung von 2 Sekunden wurde entfernt
- iOS: C-String-Speicherleck bei jeder Go→JS-JavaScript-Ausführung behoben
- iOS: `hasListeners` gibt jetzt die tatsächliche Listener-Registrierung wieder
- iOS: Debug-Protokollierung des Frameworks wird aus Produktions-Builds herauskompiliert
- Android: `GOOS=android` lässt sich wieder kompilieren – `events.Android` wurde definiert, das außerhalb der Grenzen liegende `events_android.go`-Listener-Array entfernt, der Stub für Namen mobiler Methoden hinzugefügt und verhindert, dass Desktop-Linux-Dateien (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) in Android-Builds einfließen
- Android: JS→Go-Bindings funktionieren jetzt – die WebView kann `fetch()`-POST-Bodys nicht an `shouldInterceptRequest` übermitteln, daher werden Laufzeitaufrufe über einen JavascriptInterface-Transport (`nativeHandleRuntimeCall`) geleitet, statt wegen eines nil-Anfragebodys abzustürzen
- Android: `Screens.*`-Laufzeitaufrufe geben echte Daten zurück – der ScreenManager wird jetzt beim Start befüllt (er war nie angebunden, weshalb `GetAll` nil zurückgab)
- Android: Debug-Protokollierung des Frameworks wird aus Produktions-Builds herauskompiliert und in Debug-Builds unter dem Tag `Wails` an logcat weitergeleitet
- Android: echte `hasListeners`-Registry, Behandlung von JNI-Referenzen und -Ausnahmen sowie ein Seitenlebenszyklus mit einmaligem Laden (keine doppelte Navigation)
- Behoben, dass `wails3 generate bindings` unter Windows bei laufendem Vite-Entwicklungsserver mit „Access is denied“ fehlschlug: Generierte Dateien werden jetzt mit dem Ausgabeverzeichnis synchronisiert, statt durch Umbenennen darübergelegt zu werden (#5515)
- Behebt einen sporadischen schwerwiegenden Absturz unter macOS beim Auslesen von Bildschirminformationen nach einer Anzeigeänderung: Die für Bildschirm-ID und -name gespeicherten Zeiger verwiesen auf automatisch freigegebene `UTF8String`-Puffer, die freigegeben werden konnten, bevor Go sie kopierte (Use-after-free). Die Strings werden nun mit `strdup` kopiert und nach der Konvertierung freigegeben. Außerdem läuft die Bildschirmaufzählung in einem expliziten Autorelease-Pool, sodass bei Aufrufen aus Go-Goroutinen kein Speicher mehr verloren geht (#5556)
- Behebt einen sporadischen SIGSEGV unter Linux, wenn der Assetserver ein `WebKitURISchemeRequest` schließt: Das abschließende `g_object_unref` lief in der Assetserver-Goroutine und finalisierte dadurch ein WebKit-GObject außerhalb des GTK-Hauptthreads. Der Unref-Aufruf wird nun über `g_main_context_invoke` in den GTK-Hauptkontext eingereiht (#5557)

## v3.0.0-alpha.100 - 2026-06-13

## Hinzugefügt

- Erweitert `MacWebviewPreferences` um zusätzliche WKWebView-Konfigurationsoptionen: `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` und `ApplicationNameForUserAgent` (#5549)

## Behoben

- Behebt, dass `wails3 generate bindings` unter Windows bei laufendem Vite-Entwicklungsserver mit „Access is denied“ fehlschlug, indem generierte Dateien in das Ausgabeverzeichnis synchronisiert werden, statt vorhandene Dateien durch Umbenennen zu ersetzen (#5561)
- Behebt, dass JS-Größenänderungsereignisse für rahmenlose Fenster unter Linux nicht ausgelöst wurden, sowie die Erkennung des Bildlaufleistenrands bei rahmenlosen Fenstern (#5368)
- Behebt, dass der Updater unter Windows mit „invalid cross-device link“ fehlschlug, wenn sich das temporäre Verzeichnis auf einem anderen Volume als das Installationsverzeichnis befand (#5560)

## v3.0.0-alpha.99 - 2026-06-10

## Behoben

- Behebt, dass `wails3 generate bindings` unter Windows bei laufendem Vite-Entwicklungsserver mit „Access is denied“ fehlschlug, indem generierte Dateien in das Ausgabeverzeichnis synchronisiert werden, statt vorhandene Dateien durch Umbenennen zu ersetzen (#5515)

## v3.0.0-alpha.98 - 2026-06-03

## Behoben

- Behebt das Einfrieren der WebKit-Benutzeroberfläche unter Linux im Leerlauf (z. B. bei geöffnetem Inspektor), indem `SA_ONSTACK` für `SIGUSR1` nicht mehr erzwungen wird, da dies die Synchronisierung des GC-Threads von JavaScriptCore beeinträchtigte (#5527)

## v3.0.0-alpha.97 - 2026-05-31

## Hinzugefügt

- Fügt eine Seite zum Debugging und zur Arbeit mit `runtime/trace` hinzu

## Geändert

- Entfernt einige unnötige `_ "embed"`-Importe und bereinigt dadurch den Code etwas

## Behoben

- Behebt, dass Mindestbreite und -höhe nach dem Aufheben der Fenstermaximierung unter Windows nicht erzwungen wurden (#4593)
- Behebt das Durchreichen von Mausklicks im Vollbildmodus bei den Fensteroptionen Frameless und Transparent (#4408)

## v3.0.0-alpha.96 - 2026-05-25

## Hinzugefügt

- Fügt Unterstützung für Garble-Verschleierung hinzu ([#4563](https://github.com/wailsapp/wails/issues/4563)): stabile IDs für Binding-Methoden, Integration in Build und Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`) sowie JSON-Struct-Tags für alle zur Laufzeit verwendeten Payloads (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`), damit das Übertragungsformat die Umbenennung exportierter Felder durch Garble übersteht.

## v3.0.0-alpha.95 - 2026-05-20

## Hinzugefügt

- Fügt die fehlende Seite zur Projektstruktur hinzu

## Geändert

- Dokumentation: Stellt einige Diagramme auf der Architekturseite für eine übersichtlichere Darstellung auf Sequenzdiagramme um
- Dokumentation: Fügt einen Hinweis hinzu, dass D2 zum Ausführen installiert sein muss

## Behoben

- Behebt `wails3 generate appimage` mit dem GTK4-Standard: Der Bundler erkennt nun vor der Suche nach Laufzeitdateien den GTK-Stack anhand der Binärdatei. Dadurch wählt er für GTK4-Builds `libwebkitgtkinjectedbundle.so` (unter `webkitgtk-6.0/`) und für `-tags gtk3`-Builds `libwebkit2gtkinjectedbundle.so` (unter `webkit2gtk-4.1/`). Die `.relr.dyn`-Prüfung berücksichtigt außerdem `libgtk-4.so.1`, sodass das Stripping mit modernen Toolchains unabhängig vom Stack korrekt deaktiviert wird. (#5475)
- Behebt, dass `wails3 generate appimage` bei Aufruf mit einem relativen `-builddir` fehlschlug: Der Bundler löst `-binary`, `-icon`, `-desktopfile`, `-builddir` und `-outputdir` nun vorab in absolute Pfade auf, damit das während des Ablaufs ausgeführte `s.CD` weder die AppRun-Download-Goroutine noch die nach dem Kopieren ausgeführte `ldd`-Prüfung beeinträchtigt.
- Behebt, dass `wails3 generate appimage` das fertige AppImage nicht nach `-outputdir` verschieben konnte, wenn das `Name=`-Feld der Desktop-Datei nicht mit dem Basisnamen der Binärdatei übereinstimmte: Der Bundler erzwingt nun über die Umgebungsvariable `OUTPUT`, dass das AppImage-Plug-in von linuxdeploy das AppImage nach `<binary>-<arch>.AppImage` schreibt und nicht unter dem aus der Desktop-Datei abgeleiteten Namen.
- Behebt, dass `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` und `Common.SystemDidWake` unter Linux nicht mehr ausgelöst wurden, nachdem der GTK4- und WebKitGTK-6.0-Stack in alpha.93 zum Standard erklärt worden war. Das neue standardmäßige `application_linux.go` `run()` rief weder `setupCommonEvents()` auf, das `Linux.*`-Ereignisse an die entsprechenden `Common.*`-Ereignisse weiterleitet, noch `monitorPowerEvents()`. Der DBus-Helfer für die Energieüberwachung wird nun über `application_linux_dbus.go` von den GTK3- und GTK4-Build-Pfaden gemeinsam verwendet. (#5474)

## v3.0.0-alpha.94 - 2026-05-19

## Behoben

- Behebt, dass `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` und `Common.SystemDidWake` unter Linux nicht mehr ausgelöst wurden, nachdem der GTK4- und WebKitGTK-6.0-Stack in alpha.93 zum Standard erklärt worden war. Das neue standardmäßige `application_linux.go` `run()` rief weder `setupCommonEvents()` auf, das `Linux.*`-Ereignisse an die entsprechenden `Common.*`-Ereignisse weiterleitet, noch `monitorPowerEvents()`. Der DBus-Helfer für die Energieüberwachung wird nun über `application_linux_dbus.go` von den GTK3- und GTK4-Build-Pfaden gemeinsam verwendet. (#5474)

## v3.0.0-alpha.93 - 2026-05-17

## Hinzugefügt

- Fügt `XDG_SESSION_TYPE` zur Ausgabe von `wails3 doctor` unter Linux hinzu, von @leaanthony

## Behoben

- Behebt einen Absturz des Fenstermenüs unter Wayland, der dadurch verursacht wurde, dass appmenu-gtk-module auf ein noch nicht realisiertes Fenster zugriff (#4769), von @leaanthony
- Behebt einen Absturz der GTK-Anwendung, wenn der Anwendungsname ungültige Zeichen enthält, etwa Leerzeichen oder Klammern, von @leaanthony
- Behebt den Fehler „not enough memory“ beim Initialisieren von Drag-and-drop unter Windows (#4701), von @overlordtm
- Behebt eine Race-Condition im Callback-Speicher des Hauptthreads, bei der zum Löschen aus der Map fälschlicherweise RLock verwendet wurde (Linux, macOS, iOS) (#4424), von @leaanthony
- Behebt die Variablenverarbeitung bei der Übergabe von Befehlszeilenargumenten an Tasks. Als KEY=VALUE-Paare angegebene CLI-Variablen werden nun ordnungsgemäß initialisiert und während der gesamten Task-Ausführung weitergegeben.
- Behebt einen Konflikt mit NSWindowZoomButton unter macOS: `MaximiseButtonState` und `FullscreenButtonState` wenden nun sowohl beim Start als auch zur Laufzeit den restriktiveren Zustand an; keiner der beiden Setter kann den anderen unbemerkt überschreiben (#5319)
- Behebt mehrere bereits vorhandene Fehler im veralteten GTK3-Build-Pfad (`-tags gtk3`), die CodeRabbit bei #5463 aufgedeckt hat: Starts über Dateizuordnungen überspringen keine Startup-Handler mehr; `getTheme` ist nun grenzen- und typsicher; `appName` gibt keinen GLib-eigenen Speicher mehr frei; `clipboardGet` verliert das von GTK zurückgegebene `gchar*` nicht mehr; `Calloc` verwendet nun Pointer-Receiver und `NewCalloc` gibt `*Calloc` zurück, sodass der Pool die Allokationen tatsächlich nachverfolgt; `zoomOut` verwendet den Kehrwert von `zoomInFactor` statt eines negativen Multiplikators, der auf 1.0 begrenzt wurde; `execJS` verwendet den vorab allozierten leeren World-Namen wieder, statt bei jedem Aufruf ein `C.CString("")` zu verlieren; ein für die Entwicklung vorgesehenes `fmt.Println` wurde aus `menuItem.setAccelerator` entfernt. Behebt #5465.
- Dasselbe durch einen Wertempfänger verursachte Speicherleck bei `Calloc` im standardmäßigen GTK4-Build-Pfad (`linux_cgo.go`) behoben: Zeigerempfänger und `NewCalloc() *Calloc` stellen sicher, dass die pro Fenster vorgenommenen `c.String(...)`-Allokationen tatsächlich nachverfolgt und freigegeben werden.

## v3.0.0-alpha.92 - 2026-05-15

## Hinzugefügt

- Taskfiles geändert, sodass sich der verwendete Frontend-Paketmanager über die Option `PACKAGE_MANAGER` steuern lässt
- Vorlagendaten um `{{.Opn}}` und `{{.Cls}}` erweitert, damit Taskfile-Vorlagen vorhersehbarer erstellt werden können

## Geändert

- Einige bestehende Taskfiles für die Verwendung von `{{.Opn}} and {{.Cls}}` angepasst

## Behoben

- Schwerwiegenden `concurrent map read and map write`-Laufzeitfehler in `linuxSystemTray` behoben, der auftrat, wenn das Tray-Menü aktualisiert wurde, während das Panel es las.
- Für die Ausgabe von WebView2-Fehlern und Stacktraces `log` statt `fmt` verwendet, damit Meldungen nicht verloren gehen, wenn die Anwendung unter Windows ohne verbundene Konsole ausgeführt wird.

## v3.0.0-alpha.91 - 2026-05-12

## Geändert

- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5414) von `@github-actions[bot]` aktualisiert
- **BREAKING (macOS):** Das macOS-Koordinatensystem wurde vereinheitlicht, sodass `GetScreens`, `Position` und `SetPosition` denselben Koordinatenraum verwenden: logische Punkte, nach unten verlaufende Y-Achse und `(0,0)` oben links auf dem primären Bildschirm. Dies entspricht Windows, GTK sowie den öffentlichen APIs von Electron und dem Web. Bildschirme, die sich physisch oberhalb des primären Bildschirms befinden, melden nun negative `Bounds.Y`-Werte (zuvor positive), und die Werte von `Position()`/`SetPosition()` werden jetzt in logischen Punkten statt in `points × primaryScale` angegeben. Die verlustfreie Umwandlung `Position()` → `SetPosition()` bleibt erhalten. Absolute Werte aus Protokollen früherer Alpha-Builds oder manuell berechnete Behelfslösungen (z. B. Multiplikation mit `primaryScale` oder Umkehrung von Y anhand einer Bildschirmhöhe) müssen aktualisiert werden. Behebt [#5117](https://github.com/wailsapp/wails/issues/5117).

## Behoben

- DBus-Signalnamen und Länge des Inhalts defensiv validiert, um Panics zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/5416) von @leaanthony
- Problem mit der Speichersicherheit bei der GTK-Menüverarbeitung unter Linux behoben, in [PR](https://github.com/wailsapp/wails/pull/5363) von @leaanthony
- Erkennung von NVIDIA-GPUs und Deaktivierung des DMA-BUF-Renderers unter Linux hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5295) von @leaanthony
- Bildschirmübergreifende Y-Konvertierung für `SetPosition` unter macOS korrigiert: Die Höhe des primären Bildschirms dient nun als globale Referenz, damit Fenster auf Monitoren, die gegenüber dem primären Bildschirm vertikal versetzt sind, an der richtigen Position erscheinen, in [#5117](https://github.com/wailsapp/wails/issues/5117)
- Git-PR-Vorlage korrigiert, sodass sie auf die richtige Feedback-URL verweist, in [PR](https://github.com/wailsapp/wails/pull/5109) von @wayneforrest
- Eine Gruppe von `SetMenu`-Abstürzen des Windows-Infobereichs behoben. Ursache war ein fehlerhafter `DestroyMenu`-Systemaufruf, der vier Argumente statt eines übergab, sodass jeder Aufruf FALSE zurückgab und nichts freigab. Außerdem werden HMENU- und HBITMAP-Handles, einschließlich der zur Laufzeit über `MenuItem.SetBitmap` allokierten, beim Neuaufbau des Menüs freigegeben, veraltete Zuordnungen für Kontrollkästchen und Optionsfelder in `Win32Menu.Update` zurückgesetzt und ein redundanter `Update()`-Aufruf in `systemtray.updateMenu` entfernt, der die Allokationen verdoppelte. Lang laufende Anwendungen im Infobereich verlieren nun nicht mehr bei jedem Neuaufbau des Menüs GDI-/USER-Objekte.

## v3.0.0-alpha.90 - 2026-05-11

## Hinzugefügt

- Konfigurierbaren Anwendungsnamen für den WKWebView-User-Agent unter macOS hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5261) von @vinhvoit225
- Indirekte Abhängigkeit github.com/coder/websocket zum gin-service-Beispiel hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5400) von @taliesin-ai
- Unterstützung für tiefgehende Gleichheitsvergleiche zu den Tests der Build-Assets hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5402) von @leaanthony

## Geändert

- Build-Ausgabe im Assets-Verzeichnis zusammengeführt, in [PR](https://github.com/wailsapp/wails/pull/5401) von @taliesin-ai
- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5399) von `@github-actions[bot]` aktualisiert

## Behoben

- Benachrichtigungsobjekt für macOS-Nachrichten bei einer einzelnen Anwendungsinstanz verwendet, in [PR](https://github.com/wailsapp/wails/pull/5289) von @overlordtm
- Windows-Callbacks gebündelt, um den Verlust von Promises bei hoher Last zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/5383) von @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## Hinzugefügt

- Job go<em>test</em>results zum Zusammenfassen der Go-Testergebnisse hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5316) von @leaanthony

## Geändert

- Große RPC-Nutzlasten bei Bedarf auf mehrere POST-Anfragen mit Datenblöcken aufgeteilt, in [PR](https://github.com/wailsapp/wails/pull/5369) von @leaanthony
- Vite in allen Frontend-Vorlagen von 5.x.x auf 8.0.0 aktualisiert, in [PR](https://github.com/wailsapp/wails/pull/5386) von @leaanthony
- Portkonfiguration des Vite-Entwicklungsservers auf Umgebungsvariablen umgestellt, in [PR](https://github.com/wailsapp/wails/pull/5365) von @leaanthony
- Vite-Entwicklungsserver in allen Vorlagen so konfiguriert, dass er an 127.0.0.1 gebunden wird, in [PR](https://github.com/wailsapp/wails/pull/5361) von @leaanthony
- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5384) von `@github-actions[bot]` aktualisiert

## Behoben

- Info.plist-Vorlagenfragmente bei der Aktualisierung der Build-Assets bereinigt, in [PR](https://github.com/wailsapp/wails/pull/5312) von @leaanthony
- Veralteten Zustand in macOS-Menüs behoben, indem die Mutatoren für Menüeinträge (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) synchron im Hauptthread angewendet werden. Dadurch wurde die `dispatch_async`-Race-Condition beseitigt, wegen der Menüs beim schnellen erneuten Öffnen den vorherigen Zustand anzeigten (#5002)
- `*_test.go`-Dateien im Entwicklungsmodus ignoriert, um unnötige Neubuilds zu verhindern, in [PR](https://github.com/wailsapp/wails/pull/5203) von @leaanthony
- Speicherzugriffsfehler bei Menu.Update() verhindert, wenn die Anwendung nicht ausgeführt wird, in [PR](https://github.com/wailsapp/wails/pull/5291) von @wucm667
- lastSizeWParam verwendet, um das Neuzeichnen der Menüleiste unter Windows zu steuern, in [PR](https://github.com/wailsapp/wails/pull/5382) von @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## Geändert

- HiddenOnTaskbar auf die Verwendung von WS<em>EX</em>TOOLWINDOW umgestellt, in [PR](https://github.com/wailsapp/wails/pull/5371) von @leaanthony
- Abhängigkeiten neu angeordnet und die replace-Direktive für webview2 aus go.mod entfernt, in [PR](https://github.com/wailsapp/wails/pull/5370) von @atterpac
- Sponsoren-SVG in [PR](https://github.com/wailsapp/wails/pull/5358) von `@github-actions[bot]` aktualisiert

## Behoben

- Generische Indirektionsaliase entfernt und Typen für Map-Schlüssel vereinheitlicht, in [PR](https://github.com/wailsapp/wails/pull/5331) von @fbbdev

## Entfernt

- PR-master-Workflow einschließlich Dokumentation, Go-Tests und übersprungener Tests gelöscht, in [PR](https://github.com/wailsapp/wails/pull/5377) von @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## Hinzugefügt

- Koreanische Dokumentation für Wails v3 hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5352) von @leaanthony
- Französische Dokumentation für Installation und Schnellstart hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5354) von @leaanthony
- Portugiesische Dokumentation für Schnellstart, Konzepte und Community hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5355) von @leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## Hinzugefügt

- Französische Lokalisierung der Dokumentation in [PR](https://github.com/wailsapp/wails/pull/5328) von @leaanthony hinzugefügt
- Deutsche Locale zur Dokumentationswebsite in [PR](https://github.com/wailsapp/wails/pull/5343) von @leaanthony hinzugefügt

## Geändert

- Alle 8 übersetzten Locales in der Dokumentationskonfiguration in [PR](https://github.com/wailsapp/wails/pull/5347) von @leaanthony registriert
- Verschiedene Windows-bezogene Dateien für WebView2 in [PR](https://github.com/wailsapp/wails/pull/5317) von @leaanthony aktualisiert

## Behoben

- Dialog-Dispatch unter Linux in [PR](https://github.com/wailsapp/wails/pull/5340) von @leaanthony zwischen GTK3 und GTK4 aufgeteilt
- Sichergestellt, dass Dialog-Callbacks im GTK-Thread ausgeführt werden, wodurch Speicherzugriffsfehler behoben werden, in [PR](https://github.com/wailsapp/wails/pull/5339) von @leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## Hinzugefügt

- URL der PR-Vorlage in [PR](https://github.com/wailsapp/wails/pull/5179) von @leaanthony zum Repository hinzugefügt
- Deutsche Dokumentation für Wails v3 in [PR](https://github.com/wailsapp/wails/pull/5330) von @leaanthony hinzugefügt

## v3.0.0-alpha.84 - 2026-05-03

## Hinzugefügt

- Option zum Deaktivieren des Beendens des Vollbildmodus mit der Escape-Taste unter macOS in [PR](https://github.com/wailsapp/wails/pull/5307) von @leaanthony hinzugefügt
- Option zum Deaktivieren des Beendens des Vollbildmodus mit der Escape-Taste unter macOS in [PR](https://github.com/wailsapp/wails/pull/5310) von @leaanthony hinzugefügt
- Dokumentation für Pausa zur Community-Galerie in [PR](https://github.com/wailsapp/wails/pull/5288) von @yuseferi hinzugefügt

## Geändert

- SVG-Datei der Sponsoren in [PR](https://github.com/wailsapp/wails/pull/5308) von `@github-actions[bot]` aktualisiert
- Befehl zur Symbolgenerierung für die Behandlung nicht unterstützter Plattformen in [PR](https://github.com/wailsapp/wails/pull/5309) von @leaanthony aktualisiert
- Boolesche Vollbild-API durch den dreiwertigen ButtonState ersetzt und Plattform-Bindings in [PR](https://github.com/wailsapp/wails/pull/5224) von @leaanthony implementiert

## Behoben

- WebView2-Fokusoperationen in [PR](https://github.com/wailsapp/wails/pull/5315) von @leaanthony gegen einen nil-Controller-Zustand abgesichert
- GitHub-Actions-Workflow in [PR](https://github.com/wailsapp/wails/pull/5313) von @leaanthony aktualisiert, sodass er den Basis-Branch des PR korrekt referenziert
- `*_test.go`-Dateien im Entwicklungsmodus ignoriert, um unnötige Neubuilds zu vermeiden, in [PR](https://github.com/wailsapp/wails/pull/5203) von @leaanthony
- Speicherzugriffsfehler durch Menu.Update() verhindert, wenn die App nicht ausgeführt wird, in [PR](https://github.com/wailsapp/wails/pull/5291) von @wucm667

## v3.0.0-alpha.83 - 2026-05-02

## Hinzugefügt

- InstallScope-Flag und Build-Option für die computer- bzw. benutzerbezogene Installation in [PR](https://github.com/wailsapp/wails/pull/5094) von @symball hinzugefügt
- SetScreen-Methode ohne Funktion zu BrowserWindow hinzugefügt, um das Window-Interface zu erfüllen, in [PR](https://github.com/wailsapp/wails/pull/5294) von @leaanthony

## Behoben

- NVIDIA-GPUs erkannt und DMA-BUF-Renderer unter Linux deaktiviert, in [PR](https://github.com/wailsapp/wails/pull/5295) von @leaanthony
- Git-PR-Vorlage so korrigiert, dass sie auf die richtige Feedback-URL verweist, in [PR](https://github.com/wailsapp/wails/pull/5109) von @wayneforrest
- Eine Reihe von Abstürzen der Windows-Taskleistenanwendung `SetMenu` behoben, die durch einen fehlerhaften `DestroyMenu`-Systemaufruf verursacht wurden. Dieser übergab vier Argumente statt eines, sodass jeder Aufruf FALSE zurückgab und nichts freigab. Außerdem werden HMENU- und HBITMAP-Handles – einschließlich der zur Laufzeit über `MenuItem.SetBitmap` zugewiesenen – bei der Neuerstellung von Menüs freigegeben, veraltete Kontrollkästchen-/Optionsfeld-Zuordnungen in `Win32Menu.Update` zurückgesetzt und ein redundanter `Update()`-Aufruf in `systemtray.updateMenu` entfernt, der die Zuweisungen verdoppelte. Lang laufende Taskleistenanwendungen verlieren dadurch nicht mehr bei jeder Neuerstellung eines Menüs GDI-/USER-Objekte.

## v3.0.0-alpha.82 - 2026-05-01

## Behoben

- Generierung der Desktop-Datei so korrigiert, dass der Desktop-Name richtig verarbeitet wird, in [PR](https://github.com/wailsapp/wails/pull/5232) von @leaanthony

## v3.0.0-alpha.81 - 2026-04-30

## Geändert

- Zeitplan für Nightly-Releases in [PR](https://github.com/wailsapp/wails/pull/5286) von @leaanthony auf 15:00 UTC angepasst

## Behoben

- Halbierte Werte für Screen Bounds, WorkArea und Size auf Retina-Macs behoben –  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## Geändert

- Dokumentationsabhängigkeiten und Loader für Inhaltssammlungen in [PR](https://github.com/wailsapp/wails/pull/5285) von @leaanthony aktualisiert

## v3.0.0-alpha.79 - 2026-04-29

## Hinzugefügt

- Dem trigger-release-Job in [PR](https://github.com/wailsapp/wails/pull/5270) von @leaanthony die Berechtigung actions: write erteilt

## Geändert

- Release-Task verwendet standardmäßig den master-Branch und aktualisiert die Formulierung des Changelogs, in [PR](https://github.com/wailsapp/wails/pull/5283) von @leaanthony
- auto-changelog-Workflow in [PR](https://github.com/wailsapp/wails/pull/5282) von @leaanthony auf die neueste Version aktualisiert
- Workflow-Effizienz durch das Hinzufügen von Pfadfiltern und Entfernen nicht mehr verwendeter Workflows verbessert, in [PR](https://github.com/wailsapp/wails/pull/5280) von @leaanthony
- Dokumentation so aktualisiert, dass Beispiellinks auf den master-Branch verweisen, in [PR](https://github.com/wailsapp/wails/pull/5274) von @leaanthony
- Dokumentation und Beispiele für v3 in [PR](https://github.com/wailsapp/wails/pull/5272) von @leaanthony aktualisiert

## Behoben

- Reverse-Proxy für die Entwicklung um Wiederholungslogik und die Erzwingung von IPv4 erweitert, in [PR](https://github.com/wailsapp/wails/pull/5265) von @AkagiYui
- Trigger-Workflow für den unveröffentlichten Changelog in [PR](https://github.com/wailsapp/wails/pull/5281) von @leaanthony neu geschrieben

## Entfernt

- Shell-Testskripte für verschiedene Testzwecke in [PR](https://github.com/wailsapp/wails/pull/5267) von @leaanthony entfernt
- Workflow zur Bereitstellung der v3-Alpha-Dokumentation und CNAME-Eintrag in [PR](https://github.com/wailsapp/wails/pull/5266) von @leaanthony gelöscht

### Hinzugefügt

- Eintrag „Frontend-Routing“ zur Seitenleistennavigation hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5196) von @leaanthony
- Leitfaden zum Frontend-Routing mit Framework-spezifischen Empfehlungen hinzugefügt, in [PR](https://github.com/wailsapp/wails/pull/5185) von @leaanthony
- Unterstützung für modale Sheets hinzugefügt (macOS)
- ghw-Version für eine bessere Unterstützung von Apple-Geräten durch @leaanthony aktualisiert (#4977)
- Methode `GetBadge` zum Dock-Dienst hinzugefügt
- Flag `-tags` zum Befehl `wails3 build` hinzugefügt, um benutzerdefinierte Go-Build-Tags zu übergeben (z. B. `wails3 build -tags gtk4`) (#4957)
- Dokumentation zur automatischen Enum-Generierung im Binding-Generator hinzugefügt, einschließlich einer eigenen Enums-Seite und Seitenleistennavigation (#4972)
- Flag `-tags` zum Befehl `wails3 build` hinzugefügt, um benutzerdefinierte Go-Build-Tags zu übergeben (z. B. `wails3 build -tags gtk4`) (#4957)
- Web-API-Beispiele in `v3/examples/web-apis/` hinzufügen, die 41 Browser-APIs demonstrieren, darunter Speicher (localStorage, sessionStorage, IndexedDB, Cache API), Netzwerk (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), Medien (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), Geräte (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), Leistung (Performance API, Mutation Observer, Intersection/Resize Observer), Benutzeroberfläche (Web Components, Pointer Events, Selection, Dialog, Drag and Drop) und weitere
- Beispiel für eine WebView-API-Kompatibilitätsprüfung (`v3/examples/webview-api-check/`) hinzufügen, die plattformübergreifend mehr als 200 Browser-APIs testet
- Paket `internal/libpath` zum Ermitteln nativer Bibliothekspfade unter Linux hinzufügen, mit paralleler Suche, Caching und Unterstützung für Flatpak/Snap/Nix
- **In Arbeit:** Experimentelle Unterstützung für WebKitGTK 6.0 / GTK4 unter Linux hinzufügen, verfügbar über `-tags gtk4` (GTK3/WebKit2GTK 4.1 bleibt die Standardeinstellung)
- Hinweis: Bei Tiling-Fenstermanagern (z. B. Hyprland, Sway) funktionieren Minimieren und Maximieren möglicherweise nicht wie erwartet, da der Fenstermanager die Fenstergeometrie steuert
- Die Dokumentation zu **Ereignisse in JavaScript überwachen** wurde von @AbdelhadiSeddar um die Verwendung von **einmaligen Handlern** ergänzt
- Option `UseApplicationMenu` zu `WebviewWindowOptions` hinzufügen, damit Fenster unter Windows/Linux das über `app.Menu.Set()` festgelegte Anwendungsmenü übernehmen können, von @leaanthony
- Unterstützung für `.icon`-Dateien (Apple-Icon-Composer-Format) zum Erzeugen von Liquid-Glass-Symbolen und Asset-Katalogen (macOS) hinzufügen (#4934), von @wimaha
- Experimentellen Servermodus für Headless-/Web-Bereitstellungen hinzufügen (`-tags server`). Damit können Wails-Apps ohne native GUI-Abhängigkeiten als HTTP-Server ausgeführt werden. Mit `wails3 task build:server` bauen. Weitere Informationen finden Sie unter `examples/server`.
- Paket `internal/libpath` zum Ermitteln nativer Bibliothekspfade unter Linux hinzufügen, mit paralleler Suche, Caching und Unterstützung für Flatpak/Snap/Nix
- Option `CollectionBehavior` zu `MacWindow` hinzufügen, um das Fensterverhalten über macOS Spaces und den Vollbildmodus hinweg zu steuern (#4756), von @leaanthony
- Unit-Tests für pkg/application hinzufügen, von @leaanthony
- Unterstützung benutzerdefinierter Protokolle zur MSIX-Paketierung hinzufügen, von @leaanthony
- Erkennung der Desktop-Umgebung unter Linux hinzufügen, [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- Methode `Window.Print()` zur JavaScript-Laufzeit hinzufügen, um den Druckdialog im Frontend zu öffnen (#4290), von @leaanthony
- `XDG_SESSION_TYPE` zur Ausgabe von `wails3 doctor` unter Linux hinzufügen, von @leaanthony
- Weitere WebKit2-Ereignisse für Änderungen des Ladestatus unter Linux hinzufügen: `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` (#3896), von @leaanthony
- `XDG_SESSION_TYPE` zur Ausgabe von `wails3 doctor` unter Linux hinzufügen, von @leaanthony
- Datei `.desktop` bereits beim Linux-Build und nicht erst bei der Paketierung erzeugen (#4575)
- Dokumentation der Linux-Laufzeitabhängigkeiten mit distributionsspezifischen Paketnamen und nfpm-Paketierungsbeispielen hinzufügen (#4339), von @leaanthony
- Informationen zur NVIDIA-Treiberversion zur Ausgabe von `wails3 doctor` unter Linux hinzufügen, von @leaanthony
- Ursprung zum Handler für unformatierte Nachrichten hinzufügen, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4710)
- Unterstützung für Universal Links unter macOS hinzufügen, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4712)
- Transportebene für Bindings überarbeiten, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4702)
- aria-label-Bezeichner zu den helloworld-Vorlagen hinzufügen, damit Appium-Testclients die Beispiel-App problemlos testen können, von @chinenual in [PR](https://github.com/wailsapp/wails/pull/4760)
- Ursprung zum Handler für unformatierte Nachrichten hinzufügen, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4710)
- Unterstützung für Universal Links unter macOS hinzufügen, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4712)
- Transportebene für Bindings überarbeiten, von @APshenkin in [PR](https://github.com/wailsapp/wails/pull/4702)
- Typisierte Ereignisse von @fbbdev und @ianvs in [#4633](https://github.com/wailsapp/wails/pull/4633)
- Beispiel `systray-clock` hinzufügen, das ein Headless-Taskleistensymbol mit live aktualisierten Tooltips zeigt (#4653).
- NSIS-Protokollvorlage für Windows hinzugefügt, von @Tolfx in #4510
- Tests für build-assets hinzugefügt, von @Tolfx in #4510
- macOS: Native Fenstersteuerelemente werden in [#4588](https://github.com/wailsapp/wails/pull/4588) von @nidib in der Menüleiste angezeigt
- macOS-Dock-Dienst zum Ein- und Ausblenden des App-Symbols im Dock hinzufügen, von @popaprozac in [PR](https://github.com/wailsapp/wails/pull/4451)
- macOS-Dock-Dienst zum Ein- und Ausblenden des App-Symbols im Dock hinzufügen, von @popaprozac in [PR](https://github.com/wailsapp/wails/pull/4451)
- Native Unterstützung für den Liquid-Glass-Effekt unter macOS mit NSGlassEffectView (macOS 15.0+) und NSVisualEffectView als Fallback hinzufügen, einschließlich umfassender Optionen zur Materialanpassung, von @leaanthony in [#4534](https://github.com/wailsapp/wails/pull/4534)
- Bereinigung von Browser-URLs, von @leaanthony in [#4500](https://github.dev/wailsapp/wails/pull/4500). Basierend auf [#4484](https://github.com/wailsapp/wails/pull/4484) von @APShenkin.
- Inhaltsschutz unter Windows/Mac hinzufügen, von [@leaanthony](https://github.com/leaanthony), basierend auf der ursprünglichen Arbeit von [@Taiterbase](https://github.com/Taiterbase) in diesem [PR](https://github.com/wailsapp/wails/pull/4241)
- Unterstützung für die Übergabe von CLI-Variablen an Task-Befehle über die Aliasse `wails3 build` und `wails3 package` hinzufügen (#4422), von @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4488)
- Unterstützung für Dropzones, bei denen das Ereignis die Daten des abgelegten Elements bereitstellt, von [@atterpac](https://github.com/atterpac) in [#4318](https://github.com/wailsapp/wails/pull/4318)
- `AdditionalLaunchArgs` zu den Optionen von `WindowsWindow` hinzugefügt, damit zusätzliche Befehlszeilenargumente an den WebView2-Browser übergeben werden können, in [PR](https://github.com/wailsapp/wails/pull/4467)
- Automatische Ausführung von go mod tidy nach wails init hinzugefügt, von [@triadmoko](https://github.com/triadmoko) in [PR](https://github.com/wailsapp/wails/pull/4286)
- Windows-Snap-Assist-Funktion von @leaanthony in [PR](https://github.dev/wailsapp/wails/pull/4463)
- `AdditionalLaunchArgs` zu den Optionen von `WindowsWindow` hinzugefügt, damit zusätzliche Befehlszeilenargumente an den WebView2-Browser übergeben werden können, in [PR](https://github.com/wailsapp/wails/pull/4467)
- Automatische Ausführung von go mod tidy nach wails init hinzugefügt, von [@triadmoko](https://github.com/triadmoko) in [PR](https://github.com/wailsapp/wails/pull/4286)
- Windows-Snap-Assist-Funktion von @leaanthony in [PR](https://github.dev/wailsapp/wails/pull/4463)
- Windows-Implementierung für `getAccentColor` hinzufügen, von [@almas-x](https://github.com/almas-x) in [PR](https://github.com/wailsapp/wails/pull/4427)
- Windows-Implementierung für `getAccentColor` hinzufügen, von [@almas-x](https://github.com/almas-x) in [PR](https://github.com/wailsapp/wails/pull/4427)
- Menüs und Menüleiste im dunklen Windows-Design. Von @leaanthony in [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Integrierte Dienste für verständlichere JS-/TS-Bindings umbenennen, von @popaprozac in [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` zum Abrufen der Akzentfarbe des Benutzersystems. Funktioniert unter macOS. Von [@etesam913](https://github.com/etesam913)
- API `window.ToggleFrameless()` hinzufügen, von [@atterpac](https://github.com/atterpac) in [#4137](https://github.com/wailsapp/wails/pull/4137)
- Distributionsspezifische Build-Abhängigkeiten für Linux hinzufügen, von @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4345)
- Bindings-Leitfaden von @atterpac in [PR](https://github.com/wailsapp/wails/pull/4404) hinzugefügt
- **Testinfrastruktur organisiert**: Docker-Testdateien in ein eigenes Verzeichnis `test/docker/` mit optimierten Images verschoben und die Build-Zuverlässigkeit verbessert, von [@leaanthony](https://github.com/leaanthony) in [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Muster für die Ressourcenverwaltung verbessert**: Beispiele um eine ordnungsgemäße Bereinigung von Event-Handlern und eine kontextabhängige Verwaltung von Goroutinen ergänzt, von [@leaanthony](https://github.com/leaanthony) in [#4359](https://github.com/wailsapp/wails/pull/4359)
- Unterstützung für aarch64-AppImage-Builds von [@AkshayKalose](https://github.com/AkshayKalose) in [#3981](https://github.com/wailsapp/wails/pull/3981)
- Diagnoseabschnitt zu `wails doctor` hinzugefügt, von [@leaanthony](https://github.com/leaanthony)
- Beim Aufruf einer Servicemethode das Fenster zum Kontext hinzugefügt, von [@leaanthony](https://github.com/leaanthony)
- Beispiel `window-call` hinzugefügt, das zeigt, wie sich ermitteln lässt, welches Fenster einen Service aufruft, von [@leaanthony](https://github.com/leaanthony)
- Neuer Menüleitfaden von [@leaanthony](https://github.com/leaanthony)
- Verbesserte Behandlung von Panics durch [@leaanthony](https://github.com/leaanthony)
- Neuer Menüleitfaden von [@leaanthony](https://github.com/leaanthony)
- Dokumentationskommentare für die Service-API hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4024](https://github.com/wailsapp/wails/pull/4024)
- Funktion `application.NewServiceWithOptions` hinzugefügt, um Services mit zusätzlicher Konfiguration zu initialisieren, von [@leaanthony](https://github.com/leaanthony) in [#4024](https://github.com/wailsapp/wails/pull/4024)
- Verbesserte Menüsteuerung von [@FalcoG](https://github.com/FalcoG) und [@leaanthony](https://github.com/leaanthony) in [#4031](https://github.com/wailsapp/wails/pull/4031)
- Weitere Dokumentation von [@leaanthony](https://github.com/leaanthony)
- Abbruch von Events in Standard-Event-Listenern unterstützt, von [@leaanthony](https://github.com/leaanthony)
- Unterstützung für Systray-`Hide`, -`Show` und -`Destroy` von [@leaanthony](https://github.com/leaanthony)
- Unterstützung für Systray-`SetTooltip` von [@leaanthony](https://github.com/leaanthony). Ursprüngliche Idee von [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- Paketpfad in Warnungen des Binding-Generators zu nicht unterstützten Typen ausgegeben, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Unterstützung für generische Aliase zum Binding-Generator hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Unterstützung für das JSON-Flag `omitzero` zum Binding-Generator hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Direktive `//wails:ignore` hinzugefügt, um die Binding-Generierung für ausgewählte Servicemethoden zu verhindern, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Direktive `//wails:internal` für Services und Modelle hinzugefügt, um Typen zu ermöglichen, die in Go, aber nicht in JS/TS exportiert werden, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Unterstützung für Konstanten eines Aliastyps zum Binding-Generator hinzugefügt, um schwach typisierte Enums zu ermöglichen, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Tests des Binding-Generators für Go-1.24-Funktionen hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4068](https://github.com/wailsapp/wails/pull/4068)
- Unterstützung für macOS 15 „Sequoia“ zu `OSInfo.Branding` hinzugefügt, um die Erkennung der Betriebssystemversion zu verbessern, in [#4065](https://github.com/wailsapp/wails/pull/4065)
- Hook `PostShutdown` hinzugefügt, um nach Abschluss des Herunterfahrens benutzerdefinierten Code auszuführen, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Struct `FatalError` hinzugefügt, um die Erkennung schwerwiegender Fehler in benutzerdefinierten Fehler-Handlern zu unterstützen, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Start- und Herunterfahrreihenfolge von Services vereinheitlicht und dokumentiert, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Testgerüst für die Start-/Herunterfahrsequenz der Anwendung sowie Tests für das Starten und Herunterfahren von Services hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Methode `RegisterService` hinzugefügt, um Services nach dem Erstellen der Anwendung zu registrieren, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Feld `MarshalError` zu den Anwendungs- und Serviceoptionen hinzugefügt, um Fehler bei Binding-Aufrufen benutzerdefiniert zu behandeln, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Abbrechbaren Promise-Wrapper hinzugefügt, der Abbruchanforderungen durch Promise-Ketten weitergibt, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Möglichkeit hinzugefügt, den Abbruch eines Binding-Aufrufs an einen `AbortSignal` zu koppeln, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- `data-wml-*`-Attribute neben den üblichen `wml-*`-Attributen für WML unterstützt, von [@leaanthony](https://github.com/leaanthony)
- Methode `Configure` zu allen Services hinzugefügt, um eine späte Konfiguration oder dynamische Neukonfiguration zu ermöglichen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Der Service `fileserver` sendet im unkonfigurierten Zustand eine 503-Service-Unavailable-Antwort, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Der Service `kvstore` stellt im unkonfigurierten Zustand standardmäßig einen In-Memory-Key-Value-Store bereit, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Methode `Load` zum Service `kvstore` hinzugefügt, um Daten nach Konfigurationsänderungen erneut aus der Datei zu laden, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Methode `Clear` zum Service `kvstore` hinzugefügt, um alle Schlüssel zu löschen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Typ `Level` zum Service `log` hinzugefügt, um Konstanten für Protokollierungsstufen auf der JS-Seite bereitzustellen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Methode `Log` zum Service `log` hinzugefügt, um die Protokollierungsstufe dynamisch festzulegen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Der Service `sqlite` stellt im unkonfigurierten Zustand standardmäßig eine In-Memory-Datenbank bereit, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Methode `Close` zum Service `sqlite` hinzugefügt, um die Datenbank manuell zu schließen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Abbruchunterstützung für Abfragemethoden des Services `sqlite` hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Unterstützung für vorbereitete Anweisungen mit JS-Bindings zum Service `sqlite` hinzugefügt, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Gin-Unterstützung von [Lea Anthony](https://github.com/leaanthony) in [PR](https://github.com/wailsapp/wails/pull/3537), basierend auf der ursprünglichen Arbeit von [@AnalogJ](https://github.com/AnalogJ) in diesem [PR](https://github.com/wailsapp/wails/pull/3537)
- Behebt, dass automatisches Speichern und automatisches Speichern von Passwörtern immer aktiviert sind, von [@oSethoum](https://github.com/osethoum) in [#4134](https://github.com/wailsapp/wails/pull/4134)
- Fügt `SetMenu()` für Fenster hinzu, um für ein Fenster ein Menü festlegen zu können, von [@leaanthony](https://github.com/leaanthony)
- Fügt Benachrichtigungsunterstützung hinzu, von [@popaprozac](https://github.com/popaprozac) in [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Fügt Unterstützung für Dateizuordnungen auf dem Mac hinzu, von [@wimaha](https://github.com/wimaha) in [#4177](https://github.com/wailsapp/wails/pull/4177)
- Fügt `wails3 tool version` zum Erhöhen semantischer Versionsnummern hinzu, von [@leaanthony](https://github.com/leaanthony)
- Fügt Badge-Unterstützung für macOS und Windows hinzu, von [@popaprozac](https://github.com/popaprozac) in [#](https://github.com/wailsapp/wails/pull/4234)
- Fügt Unterstützung für registrierte/streng typisierte Ereignisse hinzu, von [@fbbdev](https://github.com/fbbdev) und [@IanVS](https://github.com/IanVS) in [#4161](https://github.com/wailsapp/wails/pull/4161)
- Ermöglicht das Registrieren von Hooks für benutzerdefinierte Ereignisse, von [@fbbdev](https://github.com/fbbdev) und [@IanVS](https://github.com/IanVS) in [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` zum Öffnen des System-Dateimanagers am Pfad `path` mit optionaler Hervorhebung über `selectFile`, von [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Neues Flag `-git` für den Befehl `wails3 init`, von [@leaanthony](https://github.com/leaanthony)
- Neuer Befehl `wails3 generate webview2bootstrapper` von [@leaanthony](https://github.com/leaanthony)
- Fügt dem Runtime-Modul die Methode `init()` hinzu, um die Runtime manuell initialisieren zu können, von [@leaanthony](https://github.com/leaanthony)
- Fügt den WindowOptions eines Fensters die Option `WindowDidMoveDebounceMS` hinzu, von [@leaanthony](https://github.com/leaanthony)
- Fügt die Single-Instance-Funktion hinzu, von [@leaanthony](https://github.com/leaanthony). Basiert auf dem [v2 PR](https://github.com/wailsapp/wails/pull/2951) von @APshenkin.
- Befehl `wails3 generate template` von [@leaanthony](https://github.com/leaanthony)
- Befehl `wails3 releasenotes` von [@leaanthony](https://github.com/leaanthony)
- Befehl `wails3 update cli` von [@leaanthony](https://github.com/leaanthony)
- Option `-clean` für den Befehl `wails3 generate bindings`, von [@leaanthony](https://github.com/leaanthony)
- Ermöglicht AppImage-Linux-Builds für aarch64 (arm64), von [@AkshayKalose](https://github.com/AkshayKalose) in [#3981](https://github.com/wailsapp/wails/pull/3981)
- Fügt einen Hyperlink für den Sponsor hinzu, von @ansxuman in [#3958](https://github.com/wailsapp/wails/pull/3958)
- Unterstützung für die Linux-Paketierung als deb-, rpm- und Arch-Linux-Pakete von
- Unterstützung für universelle Darwin-Builds und -Pakete hinzugefügt von
- Ereignisdokumentation für die Website von
- Vorlagen für sveltekit und sveltekit-ts, die für die Entwicklung ohne SSR konfiguriert sind
- Aktualisiert Build-Assets mit dem neuen Befehl `wails3 update build-assets` von
- Beispiel zum Testen der HTML Drag and Drop API von
- Unterstützung für Dateizuordnungen von [leaanthony](https://github.com/leaanthony) in
- Neuer Befehl `wails3 generate runtime` von
- Neue Option `InitialPosition`, mit der angegeben wird, ob das Fenster zentriert werden soll oder
- Fügt dem Paket `application` die Methoden `Path` und `Paths` hinzu, von
- Fügt die Windows-Optionen `GeneralAutofillEnabled` und `PasswordAutosaveEnabled` hinzu
- Ermöglicht das Abrufen des Fensters, das eine Dienstmethode aufruft, von
- Fügt die Optionen `EnabledFeatures` und `DisabledFeatures` für WebView2 hinzu, von
- ⊞ Neues DIP-System für verbesserte Unterstützung von Monitoren mit hoher DPI von
- ⊞ Option für den Fensterklassennamen von [windom](https://github.com/windom/) in
- Dienste wurden erweitert, um Plugin-Funktionalität bereitzustellen. Von
- 🐧 Ereignisse WindowDidMove/WindowDidResize in
- ⊞ Ereignis WindowDidResize in
-  Fügt das Ereignis ApplicationShouldHandleReopen hinzu, um das Dock verarbeiten zu können
-  Fügt getPrimaryScreen/getScreens zur Implementierung hinzu, von @tmclane in
-  Fügt eine Option zum Anzeigen der Symbolleiste im Vollbildmodus unter macOS hinzu, von
- 🐧 Fügt onKeyPress-Logik hinzu, um einen Linux-Tastendruck in einen Accelerator umzuwandeln
- 🐧 Fügt die Task `run:linux` hinzu, von
- Exportiert die Methode `SetIcon`, von [@almas-x](https://github.com/almas-x) in
- Verbessert `OnShutdown`, von [@almas-x](https://github.com/almas-x) in
- Stellt die Methode `ToggleMaximise` in der Schnittstelle `Window` wieder her, von
- Fügt `Environment()` weitere Informationen hinzu. Von @leaanthony in
- Macht die Methode `WebviewWindow.IsFocused` über die Schnittstelle `Window` verfügbar, von
- Unterstützt im WML-System mehrere durch Leerzeichen getrennte auslösende Ereignisse, von
- Fügt ESM-Exporte aus dem gebündelten JS-Runtime-Skript hinzu, von
- Fügt dem Binding-Generator ein Flag hinzu, um das gebündelte JS-Runtime-Skript zu verwenden anstelle von
- Implementiert `setIcon` unter Linux, von [@abichinger](https://github.com/abichinger)
- Fügt dem dev-Befehl das Flag `-port` hinzu und unterstützt die Umgebungsvariable
- Fügt Tests für Aufrufe gebundener Methoden hinzu, von
- ⊞ Fügt `SetIgnoreMouseEvents` für bereits erstellte Fenster hinzu, von
-  Ermöglicht das Festlegen der Stapelungsebene (Reihenfolge) eines Fensters, von

### Behoben

- Behebt halbierte Werte für `Screen.Bounds`, `WorkArea` und `Size` auf Retina-Macs, indem NSScreen-Punktwerte in den Feldern `Physical*` in Gerätepixel umgerechnet werden. Außerdem werden `Screen.X`/`Y` auf oberster Ebene befüllt, damit in [PR](https://github.com/wailsapp/wails/pull/5168) von @wayneforrest die Erkennung sich berührender Monitore und die Platzierung im Arbeitsbereich korrekt funktionieren
- Behebt einen Data Race im ScreenManager, der bei einer Änderung der Bildschirmkonfiguration zu einem WebKit-DisplayLink-Deadlock führt (z. B. beim Anschließen eines externen Monitors während des Ruhezustands/Aufwachens)
- Setzt CFBundleIconName direkt auf appicon, wenn Assets.car vorhanden ist, in [PR](https://github.com/wailsapp/wails/pull/5154) von @symball
- Behebt, dass `wails3 doctor` unter Fedora, openSUSE, Arch und NixOS falsche WebKitGTK-Pakete meldet — die Fallback-Einträge für 4.0 wurden entfernt, da v3 zur Kompilierzeit die API 4.1 benötigt (#5071)
- Korrigiert den Namen des openSUSE-webkit2gtk-Pakets in doctor (`webkit2gtk4_1-devel` → `webkit2gtk3-devel`, der korrekte openSUSE-Paketname) (#5071)
- Fehler `Unexpected token '<'` behoben, wenn `/wails/custom.js` im Desktop-Entwicklungsmodus fehlt. Ein expliziter 404-Handler für `/wails/custom.js` und eine Validierung von `Content-Type` ohne Beachtung der Groß-/Kleinschreibung in `loadOptionalScript` wurden hinzugefügt, damit HTML-SPA-Fallbacks nicht als JavaScript eingefügt werden. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Hervorhebungsstatus des Systemleistenmenüs unter macOS korrigiert – das Symbol zeigt nun den ausgewählten Zustand, wenn das Menü geöffnet ist (#4910)
- Fehler behoben, durch den das an die Systemleiste angeheftete Fenster unter macOS hinter anderen Fenstern erschien – es verwendet nun die korrekte Fensterebene für Pop-ups (#4910)
- Fehlerhafte `@wailsio/runtime`-Importbeispiele in der gesamten Dokumentation korrigiert (#4989)
- Fehler behoben, durch den sich rahmenlose Fenster unter darwin nicht minimieren ließen (#4294)
- 20-30-minütige Hänger während `wails3 build` und `wails3 dev` behoben, indem `node_modules/` von der Aktualitätsprüfung von go-task ausgeschlossen wurde. Zuvor führte das `sources: "**/*"`-Glob dazu, dass go-task jede Datei in `node_modules/` auflistete und ihre Prüfsumme berechnete (über 50000-100000 Dateien mit umfangreichen Abhängigkeiten wie MUI), was insbesondere unter Windows/NTFS langsam war (#4939)
- GTK4-Buildfehler behoben, der durch eine Kollision des C-`Screen`-typedef mit X11 Xlib.h verursacht wurde (#4957)
- Konsistenz der Dock-Badge-Methoden unter macOS korrigiert
- Fehler behoben, durch den `InvisibleTitleBarHeight` auf alle macOS-Fenster statt nur auf rahmenlose Fenster oder Fenster mit transparenter Titelleiste angewendet wurde (#4960)
- Ruckeln und Zittern des Fensters beim Ändern der Größe an den oberen Ecken mit aktiviertem `InvisibleTitleBarHeight` behoben, indem das Einleiten des Ziehvorgangs nahe den Fensterrändern übersprungen wird (#4960)
- Generierung abgebildeter Typen mit Enum-Schlüsseln in JS/TS-Bindings korrigiert (#4437) von @fbbdev
- Fehler behoben, durch den Datei-Drag-and-drop unter Windows bei einer Anzeigeskalierung ungleich 100 % nicht funktionierte
- Fehler behoben, durch den internes HTML5-Drag-and-drop bei aktiviertem Datei-Drop unter Windows nicht funktionierte
- Fehler behoben, durch den Datei-Drop-Koordinaten unter Windows im falschen Pixelraum vorlagen (physische statt CSS-Pixel)
- Fehler behoben, durch den Datei-Drag-and-drop unter Linux mit Hover-Effekten nicht zuverlässig funktionierte
- Fehler behoben, durch den internes HTML5-Drag-and-drop bei aktiviertem Datei-Drop unter Linux nicht funktionierte
- Fehler behoben, durch den das Ein-/Ausblenden von Fenstern unter Linux/GTK4 gelegentlich den minimierten Zustand wiederherstellte; hierzu wird nun `gtk_window_present()` verwendet (#4957)
- Fehler behoben, durch den das Abrufen/Festlegen der Fensterposition unter Linux/GTK4 stets 0,0 zurückgab; hierzu wurde bedingte X11-Unterstützung über `XTranslateCoordinates`/`XMoveWindow` hinzugefügt (#4957)
- Fehler behoben, durch den die maximale Fenstergröße unter Linux/GTK4 nicht erzwungen wurde; die entfernte Funktion `gtk_window_set_geometry_hints` wurde durch eine signalbasierte Größenbegrenzung ersetzt (#4957)
- DPI-Skalierung unter Linux/GTK4 durch eine korrekte PhysicalBounds-Berechnung und Unterstützung für fraktionale Skalierung über `gdk_monitor_get_scale` korrigiert (GTK 4.14+)
- Duplizierung von Menüeinträgen beim Erstellen neuer Fenster unter Linux/GTK4 behoben
- Generierung abgebildeter Typen mit Enum-Schlüsseln in JS/TS-Bindings korrigiert (#4437) von @fbbdev
- Fehler behoben, durch den Datei-Drag-and-drop unter Windows bei einer Anzeigeskalierung ungleich 100 % nicht funktionierte
- Fehler behoben, durch den internes HTML5-Drag-and-drop bei aktiviertem Datei-Drop unter Windows nicht funktionierte
- Fehler behoben, durch den Datei-Drop-Koordinaten unter Windows im falschen Pixelraum vorlagen (physische statt CSS-Pixel)
- Fehler behoben, durch den Datei-Drag-and-drop unter Linux mit Hover-Effekten nicht zuverlässig funktionierte
- Fehler behoben, durch den internes HTML5-Drag-and-drop bei aktiviertem Datei-Drop unter Linux nicht funktionierte
- DPI-Skalierung unter Linux/GTK4 durch eine korrekte PhysicalBounds-Berechnung und Unterstützung für fraktionale Skalierung über `gdk_monitor_get_scale` korrigiert (GTK 4.14+)
- Duplizierung von Menüeinträgen beim Erstellen neuer Fenster unter Linux/GTK4 behoben
- Generierung abgebildeter Typen mit Enum-Schlüsseln in JS/TS-Bindings korrigiert (#4437) von @fbbdev
- Problem mit „Geisterfenstern“ unter macOS behoben, das dadurch verursacht wurde, dass AppKit-APIs in App.Window.Current() nicht vom Main Thread aus aufgerufen wurden (#4947) von @wimaha
- Fehler behoben, durch den HTML-`<input type="file">` unter macOS nicht funktionierte; hierzu wurde WKUIDelegate runOpenPanelWithParameters implementiert (#4862)
- Fehler behoben, durch den natives Datei-Drag-and-drop bei Verwendung des npm-Moduls `@wailsio/runtime` unter macOS/Linux nicht funktionierte (#4953) von @leaanthony
- Binding-Generierung für paketübergreifende Typaliase korrigiert (#4578) von @fbbdev
- Absturz von OpenFileDialog unter Linux aufgrund einer Verletzung der GTK-Threadsicherheit behoben (#3683) von @ddmoney420
- SIGSEGV-Absturz beim Aufruf von `Focus()` für ein ausgeblendetes oder zerstörtes Fenster behoben (#4890) von @ddmoney420
- Mögliche Panic beim Festlegen eines leeren Symbols oder Bitmaps unter Linux behoben (#4923) von @ddmoney420
- Absturz von ErrorDialog beim Aufruf aus einem Service-Binding unter macOS behoben (#3631) von @leaanthony
- Anzeige von Menüs unter Windows OS in `v3\examples\dialogs` ermöglicht von @ndianabasi
- Race Condition behoben, die beim Neuladen der Seite einen TypeError verursachte (#4872) von @ddmoney420
- Fehlerhafte Ausgabe der Tests des Binding-Generators behoben, indem der globale Zustand aus der Methode `Collector.IsVoidAlias()` entfernt wurde (#4941) von @fbbdev
- Fehler behoben, durch den die `<input type="file">`-Dateiauswahl unter macOS nicht funktionierte (#4862) von @leaanthony
- Fehler behoben, durch den `Position()` und `SetPosition()` unter macOS unterschiedliche Koordinatensysteme verwendeten, sodass sich die Fensterposition beim Speichern/Wiederherstellen des Zustands verschob (#4816) von @leaanthony
- SetProcessDpiAwarenessContext-Fehler „Access is denied“ behoben, wenn die DPI-Sensitivität bereits über das Anwendungsmanifest festgelegt ist (#4803)
- Dokumentationsseite zu Tastenkombinationen aktualisiert und Typ des Callback-Parameters für `KeyBinding.Add` korrigiert von @ndianabasi
- Dokumentation zur Generierung benutzerdefinierter Bindings korrigiert: Es muss `-d String` statt `-o String` verwendet werden
- Fehler behoben, durch den unter `menu.Update()` untergeordnete Menüelemente nicht entfernt wurden
- Veraltete Verweise auf die Manager-API in der Dokumentation korrigiert (31-Dateien wurden auf das neue Muster mit `app.Window.New()`, `app.Event.Emit()` usw. umgestellt) von @leaanthony
- Linux-Absturz bei einer Panic in an JS gebundenen Go-Methoden behoben, der durch das Überschreiben von Signal-Handlern durch WebKit verursacht wurde (#3965) von @leaanthony
- Fehler behoben, durch den SaveFileDialog.SetFilename() unter Linux wirkungslos war (#4841) von @samstanier
- Fehler behoben, durch den Drop-Koordinaten im Drag-and-drop-Beispiel als undefined angezeigt wurden
- Fehler beim Erstellen von macOS-App-Bundles behoben, wenn APP_NAME Leerzeichen enthält (Problem mit Brace Expansion)
- Index-out-of-bounds-Panic unter Windows beim Aufrufen von Service-Methoden behoben (goccy/go-json zurückgesetzt)
- Fehler behoben, durch den Datei-Drag-and-drop unter Windows bei einer Anzeigeskalierung ungleich 100 % nicht funktionierte
- Fehler behoben, durch den internes HTML5-Drag-and-drop bei aktiviertem Datei-Drop unter Windows nicht funktionierte
- Fehler behoben, durch den Datei-Drop-Koordinaten unter Windows im falschen Pixelraum vorlagen (physische statt CSS-Pixel)
- Fehler behoben, durch den Datei-Drag-and-drop unter Linux mit Hover-Effekten nicht zuverlässig funktionierte
- Fehler behoben, durch den internes HTML5-Drag-and-drop bei aktiviertem Datei-Drop unter Linux nicht funktionierte
- Alle Befehle in Taskfile.yml-Dateien für sämtliche Betriebssysteme aktualisiert, damit Leerzeichen in Variablen wie `APP_NAME` unterstützt werden von @ndianabasi
- Fehler bei Befehlsargumenten beim Ausführen des Tasks 'build:universal:lipo:go' unter Linux behoben von @wux1an
- Docker-Fehler „undefined symbol: **<em>ubsan</em>handle_xxxxxxx“ beim Ausführen von 'wails3 build GOOS=darwin GOARCH=arm64' unter Linux behoben von @wux1an
- Dokumentation zu benutzerdefinierten Protokollen zusammengeführt und Abschnitte zu Universal Links hinzugefügt von @leaanthony
- Absturz des Windows-Taskleistenmenüs beim wiederholten Klicken auf das Symbol durch eine Schutzabfrage gegen gleichzeitige TrackPopupMenuEx-Aufrufe behoben (#4151) von @leaanthony
- Absturz der Anwendung beim Aufruf von systray.Run() vor app.Run() verhindert von @leaanthony
- Absturz unter macOS beim Umschalten der Fenstersichtbarkeit über Hide()/Show() mit aktiviertem ApplicationShouldTerminateAfterLastWindowClosed behoben (#4389) von @leaanthony
- Speicherleck in Kontextmenüs unter macOS und Windows beim wiederholten Öffnen der Menüs behoben (#4012) von @leaanthony
- Fehlende Wiederverwendung nativer Kontextmenüressourcen unter macOS behoben, durch die bei jeder Anzeige ein neues Menü erstellt wurde (#4012) von @leaanthony
- Fehler behoben, durch den ein Klick auf das macOS-Dock-Symbol ausgeblendete Fenster nicht anzeigte, wenn die Anwendung mit `Hidden: true` gestartet wurde (#4583) von @leaanthony
- Fehler behoben, durch den sich der Druckdialog unter macOS wegen eines falschen Fensterzeigertyps im CGO-Aufruf nicht öffnete (#4290) von @leaanthony
- Absturz des Fenstermenüs unter Wayland behoben, der durch den Zugriff von appmenu-gtk-module auf ein nicht realisiertes Fenster verursacht wurde (#4769) von @leaanthony
- Absturz der GTK-Anwendung behoben, wenn der Anwendungsname ungültige Zeichen wie Leerzeichen oder Klammern enthält von @leaanthony
- Fehler „not enough memory“ beim Initialisieren von Drag-and-drop unter Windows behoben (#4701) von @overlordtm
- Fehler behoben, durch den der Datei-Explorer unter Linux wegen einer fehlerhaften URI-Escapierung das falsche Verzeichnis öffnete (#4397) von @leaanthony
- Fehlschlagen von AppImage-Builds auf modernen Linux-Distributionen (Arch, Fedora 39+, Ubuntu 24.04+) durch automatische Erkennung von `.relr.dyn`-ELF-Abschnitten und Deaktivierung des Strippings behoben (#4642) von @leaanthony
- Fehler behoben, durch den `wails doctor` auf Fedora- und DNF-basierten Systemen fälschlicherweise WebKit-Pakete als installiert meldete (#4457) von @leaanthony
- Fehler behoben, durch den der Standardwert `config.yml` bei einem Produktions-Build `wails3 dev` ausführte von @mbaklor
- Build-Fehler behoben, die durch den Import eines nicht vorhandenen Pakets in iOS-Service-Stubs verursacht wurden von @leaanthony
- Fehler „no formatting directives“ behoben, die durch strukturiertes Logging in Debug-/Info-Methoden verursacht wurden von @leaanthony
- Temporäre Debug-Ausgaben entfernt, die versehentlich aus dem Merge der Mobilplattform übernommen worden waren von @leaanthony
- WebKitGTK-Absturz unter Wayland mit NVIDIA-GPUs (Fehler 71 Protocol error) durch automatisches Deaktivieren des DMA-BUF-Renderers behoben von @leaanthony
- Fehler behoben, durch den der Alphawert in `application.WebviewWindowOptions.BackgroundColour` unter Linux ignoriert wurde ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Fehler behoben, durch den unter Windows nicht standardmäßig das Anwendungssymbol als Taskleistensymbol verwendet wurde, wenn kein benutzerdefiniertes Symbol angegeben war (#4704)
- Eigentümerschaft von `HICON` wird nachverfolgt, sodass nur vom Benutzer erstellte Handles zerstört werden und Abstürze beim Recycling durch den Explorer verhindert werden (#4653).
- Listener für das Windows-Systemdesign und vorgehaltene Taskleistensymbole beim Zerstören freigegeben, um Lecks von Goroutinen und Gerätekontexten zu verhindern (#4653).
- Taskleisten-Tooltips auf 127 UTF-16-Einheiten gekürzt, um die Beschädigung von Surrogatpaaren und Mehrbyte-Glyphen zu vermeiden (#4653).
- Fehlschlagen der Windows-Paketierungsaufgabe behoben (#4667)
- Appicon-Variable für Linux-AppImage in der Linux-Taskfile korrigiert [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Windows-Build-Fehler behoben, der durch eine Signaturänderung in go-webview2 v1.0.22 verursacht wurde (#4513, #4645)
- Appicon-Variable für Linux-AppImage in der Linux-Taskfile korrigiert [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Protokollbereich in linux desktop.tmpl durch Entfernen von `<.Info.Protocol>` bis `<.Protocol>` korrigiert von @Tolfx in #4510
- Neudefinitionsfehler in der Liquid-Glass-Demo behoben in [#4542](https://github.com/wailsapp/wails/pull/4542) von @Etesam913
- Aktualisierungen des Taskleistenmenüs unter Linux behoben [#4604](https://github.com/wailsapp/wails/issues/4604) von [@JackDoan](https://github.com/JackDoan)
- Erscheinen eines weißen Fensters unter Windows beim Erstellen eines ausgeblendeten Fensters behoben von @leaanthony in [#4612](https://github.com/wailsapp/wails/pull/4612)
- Importpfad des Benachrichtigungspakets in der Dokumentation korrigiert von @rxliuli in [#4617](https://github.com/wailsapp/wails/pull/4617)
- Fehler behoben, durch den Drag-and-drop bei Verwendung des npm-Pakets @wailsio/runtime nicht funktionierte (#4489) von @leaanthony in #4616
- Windows: Flackern des Fensters beim Start und fehlerhafte Anzeige ausgeblendeter Fenster behoben in [PR](https://github.com/wailsapp/wails/pull/4600) von @leaanthony.
- Probleme mit der Fenstergröße beim Maximieren unter Wayland behoben (https://github.com/wailsapp/wails/issues/4429) von [@samstanier](https://github.com/samstanier)
- Probleme mit der Fenstergröße beim Maximieren unter Wayland behoben (https://github.com/wailsapp/wails/issues/4429) von [@samstanier](https://github.com/samstanier)
- Neudefinitionsfehler in der Liquid-Glass-Demo behoben in [#4542](https://github.com/wailsapp/wails/pull/4542) von @Etesam913
- Problem behoben, durch das AssetServer unter macOS abstürzen konnte, in [#4576](https://github.com/wailsapp/wails/pull/4576) von @jghiloni
- Kompilierungsproblem beim Build mit NextJs behoben. Behoben in [#4585](https://github.com/wailsapp/wails/pull/4585) von @rev42
- Pipelines für Nightly-Releases korrigiert in [#4597](https://github.com/wailsapp/wails/pull/4597) von @riadafridishibly
- Neudefinitionsfehler in der Liquid-Glass-Demo behoben in [#4542](https://github.com/wailsapp/wails/pull/4542) von @Etesam913
- Problem behoben, durch das AssetServer unter macOS abstürzen konnte, in [#4576](https://github.com/wailsapp/wails/pull/4576) von @jghiloni
- Kompilierungsproblem beim Build mit NextJs behoben. Behoben in [#4585](https://github.com/wailsapp/wails/pull/4585) von @rev42
- Pipelines für Nightly-Releases korrigiert in [#4597](https://github.com/wailsapp/wails/pull/4597) von @riadafridishibly
- Neudefinitionsfehler in der Liquid-Glass-Demo behoben in [#4542](https://github.com/wailsapp/wails/pull/4542) von @Etesam913
- SetBackgroundColour unter Windows korrigiert von @PPTGamer in [PR](https://github.com/wailsapp/wails/pull/4492)
- Dokumentation an die Änderungen durch das Refactoring der Manager API angepasst von @yulesxoxo in [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Appicon-Variable der Linux-.desktop-Datei in der Linux-Taskfile korrigiert [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Dokumentation an die Änderungen durch das Refactoring der Manager API angepasst von @yulesxoxo in [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Unter Windows auftretenden Fehler durch Dereferenzierung eines nil-Zeigers behoben, gemeldet in [#4456](https://github.com/wailsapp/wails/issues/4456), von @leaanthony in [#4460](https://github.com/wailsapp/wails/pull/4460)
- Unterstützung für `allowsBackForwardNavigationGestures` in WKWebView unter macOS hinzugefügt, um Navigationsgesten durch Wischen mit zwei Fingern zu ermöglichen (#1857)
- Problem behoben, durch das onClick bei Menüeinträgen nicht funktionierte, die anfänglich als deaktiviert festgelegt waren, von @leaanthony in [PR #4469](https://github.com/wailsapp/wails/pull/4469). Vielen Dank an @IanVS für die anfängliche Untersuchung.
- Fehler behoben, durch den der Vite-Server nach einem fehlgeschlagenen Build nicht bereinigt wurde (#4403)
- Panic beim Schließen oder Abbrechen eines `SaveFileDialog` unter Windows behoben. Behoben in [PR](https://github.com/wailsapp/wails/pull/4284) von @hkhere
- Drag-and-drop auf HTML-Ebene unter Windows wurde von [@mbaklor](https://github.com/mbaklor) in [#4259](https://github.com/wailsapp/wails/pull/4259) korrigiert
- Unterstützung für `allowsBackForwardNavigationGestures` in WKWebView unter macOS hinzugefügt, um Navigationsgesten durch Wischen mit zwei Fingern zu ermöglichen (#1857)
- Problem behoben, durch das onClick bei Menüeinträgen nicht funktionierte, die anfänglich als deaktiviert festgelegt waren; behoben von @leaanthony in [PR #4469](https://github.com/wailsapp/wails/pull/4469). Vielen Dank an @IanVS für die anfängliche Untersuchung.
- Bereinigung des Vite-Servers bei einem fehlgeschlagenen Build korrigiert (#4403)
- Parsen von Benachrichtigungen unter Windows korrigiert, von @popaprozac in [PR](https://github.com/wailsapp/wails/pull/4450)
- Prüfung auf Abhängigkeiten des Windows SDK im doctor-Befehl wurde von [@kodumulo](https://github.com/kodumulo) in [#4390](https://github.com/wailsapp/wails/issues/4390) korrigiert
- Nil-Zeiger-Dereferenzierung in processURLRequest auf dem Mac wurde von [@etesam913](https://github.com/etesam913) in [#4366](https://github.com/wailsapp/wails/pull/4366) behoben
- Linux-Fehler behoben, der gefilterte Dialoge verhinderte, von [@bh90210](https://github.com/bh90210) in [#4287](https://github.com/wailsapp/wails/pull/4287)
- Probleme mit dem Bearbeiten-Menü unter Windows und Linux wurden von [@leaanthony](https://github.com/leaanthony) in [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062) behoben
- Minimale Systemversion in macOS-.plist-Dateien wurde von 10.13.0 auf 10.15.0 aktualisiert, von [@AkshayKalose](https://github.com/AkshayKalose) in [#3981](https://github.com/wailsapp/wails/pull/3981)
- Problem mit übersprungenen Fenster-IDs behoben von [@leaanthony](https://github.com/leaanthony)
- Problem mit einem nil-Menü beim Aufruf von RegisterContextMenu behoben von [@leaanthony](https://github.com/leaanthony)
- Abhängigkeitszyklen in der Ausgabe des Binding-Generators wurden von [@fbbdev](https://github.com/fbbdev) in [#4001](https://github.com/wailsapp/wails/pull/4001) behoben
- Use-before-define-Fehler in der Ausgabe des Binding-Generators wurden von [@fbbdev](https://github.com/fbbdev) in [#4001](https://github.com/wailsapp/wails/pull/4001) behoben
- Build-Flags werden an den Binding-Generator übergeben, von [@fbbdev](https://github.com/fbbdev) in [#4023](https://github.com/wailsapp/wails/pull/4023)
- Pfade im Windows-Taskfile wurden auf Schrägstriche umgestellt, damit es auf Nicht-Windows-Plattformen funktioniert, von [@leaanthony](https://github.com/leaanthony)
- Mac- und Mac-JS-Ereignisse wurden von [@leaanthony](https://github.com/leaanthony) korrigiert
- Ereignis-Deadlock unter macOS behoben von [@leaanthony](https://github.com/leaanthony)
- Ein `Parameter incorrect`-Fehler bei der Fensterinitialisierung unter Windows wurde behoben, wenn HTML, aber kein JS bereitgestellt wurde, von [@leaanthony](https://github.com/leaanthony)
- Größe des Antwortpräfixes für die Inhaltstyperkennung im Asset-Server wurde von [@fbbdev](https://github.com/fbbdev) in [#4049](https://github.com/wailsapp/wails/pull/4049) korrigiert
- Verarbeitung von Nicht-404-Antworten auf dem Root-Indexpfad im Asset-Server wurde von [@fbbdev](https://github.com/fbbdev) in [#4049](https://github.com/wailsapp/wails/pull/4049) korrigiert
- Undefiniertes Verhalten im Binding-Generator beim Prüfen von Eigenschaften generischer Typen wurde von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045) behoben
- Ausgabe des Binding-Generators für Modelle korrigiert, wenn der zugrunde liegende Typ nicht dieselben Eigenschaften wie der benannte Wrapper besitzt, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Ausgabe des Binding-Generators für Map-Schlüsseltypen und die Vorverarbeitung wurde von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045) korrigiert
- Ausgabe des Binding-Generators für Structs, die Marshaler-Schnittstellen implementieren, wurde von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045) korrigiert
- Erkennung von Typzyklen mit generischen Typen im Binding-Generator wurde von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045) korrigiert
- Ungültige Verweise auf nicht exportierte Modelle in der Ausgabe des Binding-Generators wurden von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045) behoben
- Injizierter Code wurde von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045) an das Ende der Service-Dateien verschoben
- Verarbeitung von Fehlern beim Schließen von Dateien im Binding-Generator wurde von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045) korrigiert
- Warnungen für Services unterdrückt, die Lifecycle- oder HTTP-Methoden, aber keine anderen gebundenen Methoden definieren, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- Bei Nicht-React-Vorlagen wurde behoben, dass bei einem hellen Systemfarbschema die Hello-World-Fußzeile nicht angezeigt wurde, von [@marcus-crane](https://github.com/marcus-crane) in [#4056](https://github.com/wailsapp/wails/pull/4056)
- Ausgeblendete Menüeinträge unter macOS behoben von [@leaanthony](https://github.com/leaanthony)
- Verarbeitung und Formatierung von Fehlern in Nachrichtenprozessoren wurde von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066) korrigiert
-  Übersprungenes Herunterfahren von Services beim Beenden der Anwendung wurde von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066) behoben
-  Menüaktualisierungen erfolgen nun garantiert im Hauptthread, von [@leaanthony](https://github.com/leaanthony)
- Der Mechanismus zum Ziehen und Ändern der Größe ist jetzt robuster und entspricht dem erwarteten Plattformverhalten genauer, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- [#4097](https://github.com/wailsapp/wails/issues/4097) behoben: Webpack/Angular verwirft Runtime-Initialisierungscode, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Anfänglich ausgeblendete Menüeinträge wurden von [@IanVS](https://github.com/IanVS) in [#4116](https://github.com/wailsapp/wails/pull/4116) korrigiert
- Behoben, dass assetFileServer bei einer Anfrage ohne Erweiterung keine `.html`-Dateien auslieferte, wenn `[request]` nicht existiert, `[request].html` jedoch schon
- Pfade für die Symbolgenerierung wurden von [@robin-samuel](https://github.com/robin-samuel) in [#4125](https://github.com/wailsapp/wails/pull/4125) korrigiert
- Behoben, dass die Ereignisse `fullscreen`, `unfullscreen`, `unminimise` und `unmaximise` nicht ausgelöst wurden, von [@oSethoum](https://github.com/osethoum) in [#4130](https://github.com/wailsapp/wails/pull/4130)
- NSIS-Fehler aufgrund eines falschen Präfixes bei der Standardversion in der Konfiguration wurde von [@robin-samuel](https://github.com/robin-samuel) in [#4126](https://github.com/wailsapp/wails/pull/4126) behoben
- Behoben, dass die Runtime-Funktion Dialogs unter Windows maskierte Pfade zurückgab, von [TheGB0077](https://github.com/TheGB0077) in [#4188](https://github.com/wailsapp/wails/pull/4188)
- WebView2-Erkennungspfad in HKCU wurde von [@leaanthony](https://github.com/leaanthony) korrigiert.
- Eingabeproblem unter macOS behoben von [@leaanthony](https://github.com/leaanthony).
- Name der Taskdatei für die Windows-Symbolgenerierung wurde von [@yulesxoxo](https://github.com/yulesxoxo) in [#4219](https://github.com/wailsapp/wails/pull/4219) korrigiert.
- Transparenzproblem bei rahmenlosen Fenstern behoben von [@leaanthony](https://github.com/leaanthony), basierend auf der Arbeit von @kron.
- Fokusaufrufe bei deaktivierten oder minimierten Fenstern wurden von [@leaanthony](https://github.com/leaanthony) korrigiert, basierend auf der Arbeit von @kron.
- Behoben, dass Infobereichssymbole nach einem Neustart der Taskleiste nicht angezeigt wurden, von [@leaanthony](https://github.com/leaanthony), basierend auf der Arbeit von @kron.
- Behoben: fallbackResponseWriter implementierte Flush() in [#4245](https://github.com/wailsapp/wails/pull/4245) nicht
- Behoben: fallbackResponseWriter implementierte Flush() nicht; von [@superDingda] in [#4236](https://github.com/wailsapp/wails/issues/4236)
- Behoben: Abstürze beim Schließen eines macOS-Fensters während eines ausstehenden asynchronen Aufrufs einer an Go gebundenen Funktion; von [@joshhardy](https://github.com/joshhardy) in [#4354](https://github.com/wailsapp/wails/pull/4354)
- Behoben: Race-Condition beim Start des Windows-Effizienzmodus; von [@leaanthony](https://github.com/leaanthony)
- Behoben: Bereinigung von Windows-Icon-Handles; von [@leaanthony](https://github.com/leaanthony).
- Behoben: `OpenFileManager` unter Windows; von [@PPTGamer](https://github.com/PPTGamer) in [#4375](https://github.com/wailsapp/wails/pull/4375).
- Behoben: Optionen für minimale/maximale Breite unter Linux; von @atterpac in [#3979](https://github.com/wailsapp/wails/pull/3979)
- TypeScript-Typdefinitionen für Vorlagen durch Anhebung der npm-Version; von @atterpac in [#3966](https://github.com/wailsapp/wails/pull/3966)
- Behoben: CSS-Referenz der SvelteKit-Vorlage; von @atterpac in [#3945](https://github.com/wailsapp/wails/pull/3945)
- Sichergestellt, dass wichtige Callbacks in window run() im Hauptthread aufgerufen werden; von [@leaanthony](https://github.com/leaanthony)
- Behoben: Beispiele für die Verzeichnisauswahl in Dialogen; von [@leaanthony](https://github.com/leaanthony)
- Neue chinesische Fehlerseite für eine fehlende index.html erstellt; von [@leaanthony](https://github.com/leaanthony)
-  Sichergestellt, dass der `windowDidBecomeKey`-Callback im Hauptthread ausgeführt wird; von [@leaanthony](https://github.com/leaanthony)
-  Vollbildmodus für rahmenlose Fenster unterstützt; von [@leaanthony](https://github.com/leaanthony)
-  Logik zum Zerstören von Fenstern verbessert; von [@leaanthony](https://github.com/leaanthony)
-  Behoben: Logik der Fensterposition bei Anbindung an Taskleistensymbole; von [@leaanthony](https://github.com/leaanthony)
-  Vollbildmodus für rahmenlose Fenster unterstützt; von [@leaanthony](https://github.com/leaanthony)
- Behoben: Ereignisbehandlung; von [@leaanthony](https://github.com/leaanthony)
- Behoben: Logik zum Schließen von Fenstern; von [@leaanthony](https://github.com/leaanthony)
- Das gemeinsame Taskfile erzeugt für TypeScript-Vorlagen nun standardmäßig TypeScript-Bindings; von [@leaanthony](https://github.com/leaanthony)
- Behoben: Schließen der Anwendung bei einer WM_CLOSE-Nachricht, wenn keine Fenster geöffnet sind bzw. nur ein Taskleistensymbol vorhanden ist; von [@mmalcek](https://github.com/mmalcek) in [#3990](https://github.com/wailsapp/wails/pull/3990)
- Behoben: garble-Build; von @5aaee9 in [#3192](https://github.com/wailsapp/wails/pull/3192)
- Behoben: Windows-NSIS-Builds; von [@leaanthony](https://github.com/leaanthony)
- Behoben: Deadlock im Linux-Dialog für Mehrfachauswahl, verursacht durch nicht geschlossene
- Behoben: Plattformübergreifende Bereinigung von .syso-Dateien während des Windows-Builds; von
- Behoben: Kompilierung des amd64-AppImage; von @atterpac in
- Behoben: Aktualisierung der Build-Assets; von @ansxuman in
- Behoben: Implementierung von `OnClick` und `OnRightClick` für Linux-Taskleistensymbole; von @atterpac
- Behoben: `AlwaysOnTop` funktionierte auf dem Mac nicht; von
-  Behoben: `application.NewEditMenu` enthielt ein Duplikat
- 🐧 Behoben: aarch64-Kompilierung
- ⊞ Behoben: Menüelemente für Optionsgruppen; von
- Behoben: Fehler beim Erstellen einer ausführbaren .app unter macOS, wenn 'name' und 'outputfilename'
- Behoben: Fehler bei der Verwendung von customEventProcessor im Drag-and-drop-Beispiel; von
- 🐧 Behoben: Durch das Hinzufügen von IgnoreMouseEvents verursachter Linux-Kompilierungsfehler; von
- ⊞ Behoben: Fehler bei der Erzeugung der syso-Icon-Datei; von
- 🐧 Korrektur für die native Ausführung unter Wayland übernommen von
- Interne Dienstmethoden nicht binden in
- ⊞ Behoben: Panic beim Start des Taskleistensymbols in
- Interne Dienstmethoden nicht binden in
- ⊞ Behoben: Panic beim Start des Taskleistensymbols in
- Umfassende Überarbeitung der Menüelemente und Ereignisbehandlung. Verbessert vorerst hauptsächlich macOS. Von
- Tests nach der Überarbeitung von Plugins und Ereignissen korrigiert in
- ⊞ Behoben: `Failed to unregister class Chrome_WidgetWin_0`-Warnung. Von
- Modulprobleme
- Behoben: Nachrichtenübermittlung bei Größenänderungsereignissen; von [atterpac](https://github.com/atterpac) in
- 🐧 Behoben: Fehler bei der Theme-Verarbeitung unter NixOS; von
- Behoben: Projektinstallation über Volume-Grenzen hinweg unter Windows; von
- CSS der React-Vorlage korrigiert, damit die Fußzeile angezeigt wird; von
- Behoben: Zombie-Prozesse beim Arbeiten im Entwicklungsmodus durch Aktualisierung auf die neueste refresh-Version
- Behoben: Einbindung der WebKit-Datei für AppImage; von [Atterpac](https://github.com/atterpac)
- Behoben: Prüfung von apt-Paketen durch Doctor; von [Atterpac](https://github.com/Atterpac) in
- Behoben: Anwendung fror beim Beenden ein (Darwin); von @5aaee9 in
- Behoben: Hintergrundfarben der Beispiele unter Windows; von
- Behoben: Standardkontextmenüs; von [mmghv](https://github.com/mmghv) in
- Behoben: Hexadezimalwerte für Pfeiltasten unter Darwin; von
- Drag-and-drop unter Windows funktionsfähig gemacht. Hinzugefügt von
- Behoben: Linux-Fehler in Doctor, wenn der Benutzer nicht über die richtigen Treiber verfügt
- DPI-Skalierung beim Start korrigiert (Windows). Geändert von [@almas-x](https://github.com/almas-x) in
- Ersetzungszeile in `go.mod` auf relative Pfade umgestellt. Korrigiert Windows-Pfade mit
- Behoben: Behandlung von Klicks auf das macOS-Taskleistensymbol ohne angebundenes Fenster; von
- Behoben: Fehlschlagender Windows-Build aufgrund einer unbekannten Option; von
- Behoben: Absturz unter Windows beim Linksklick auf das Taskleistensymbol, wenn kein
- Behoben: Falsche baseURL beim zweimaligen Öffnen eines Fensters; von @5aaee9 in PR
- Reihenfolge der if-Zweige in der Methode `WebviewWindow.Restore` korrigiert; von
- `startURL` über mehrere Aufrufe von `GetStartURL` hinweg korrekt berechnet, wenn
- JS-Typ der Struktur `Screen` an das Go-Gegenstück angeglichen; von
- Methode `WML.Reload` korrigiert, um die ordnungsgemäße Bereinigung registrierter Ereignisse sicherzustellen
- Behoben: Benutzerdefiniertes Kontextmenü schloss sich unter Linux sofort; von
- Ausgabepfad und Erweiterung der vom Binding erzeugten Modelldateien korrigiert
- Importpfade der Modelldateien im vom Binding erzeugten JS-Code korrigiert
- Drag-and-drop auf einigen Linux-Distributionen korrigiert; von
- Fehlenden Task für macOS bei Verwendung von `wails3 task dev` ergänzt; von
- Behoben: Registrierung von Ereignissen verursachte eine Zuweisung zu einer nil-Map; von
- Unmarshalling der Parameter gebundener Methoden korrigiert; von
- Verarbeitung mehrerer Rückgabewerte gebundener Methoden korrigiert von
- Erkennung von npm durch Doctor korrigiert, wenn npm nicht mit dem Systempaketmanager installiert wurde
- Fehlende MicrosoftEdgeWebview2Setup.exe ergänzt. Dank an
- Zufälligen Absturz unter Linux aufgrund der Handhabung von Fenster-IDs korrigiert von @leaanthony. Basierend auf
- Absturz von systemTray.setIcon unter Linux korrigiert von
- Sichergestellt, dass der Fensterrahmen beim ersten Aufruf der Funktion `setFrameless` angewendet wird auf

### Geändert

- **INAKOMPATIBLE ÄNDERUNG**: Map-Schlüssel in generierten JS/TS-Bindings sind jetzt als optional gekennzeichnet, um die Semantik von Go-Maps korrekt abzubilden. Der Zugriff auf Map-Werte in TypeScript gibt jetzt `T | undefined` statt `T` zurück, sodass Nullprüfungen oder Assertions erforderlich sind (#4943) von `@fbbdev`
- Verwendung von `Event` gemäß den Änderungen in `@wailsio/runtime` in `Events` geändert und die entsprechenden Funktionsaufrufe in der Dokumentation unter `Features/Events/Event System` angepasst von @AbdelhadiSeddar
- `EnabledFeatures`, `DisabledFeatures` und `AdditionalBrowserArgs` aus den fensterspezifischen Optionen in `Options.Windows` auf Anwendungsebene verschoben (#4559) von @leaanthony
- README für das Beispiel `Drag N Drop` aktualisiert und hervorgehoben, dass `Internal Drag and Drop` im Beispiel demonstriert wird, von @ndianabasi
- Protokollierungsstufe verschiedener Debug-Meldungen von Info in Debug geändert (von @mbaklor)
- **INAKOMPATIBLE ÄNDERUNG:** `EnableDragAndDrop` in den Fensteroptionen in `EnableFileDrop` umbenannt
- **INAKOMPATIBLE ÄNDERUNG:** `DropZoneDetails` im Ereigniskontext in `DropTargetDetails` umbenannt
- **INAKOMPATIBLE ÄNDERUNG:** Methode `DropZoneDetails()` von `WindowEventContext` in `DropTargetDetails()` umbenannt
- **INAKOMPATIBLE ÄNDERUNG:** Ereignis `WindowDropZoneFilesDropped` entfernt; stattdessen `WindowFilesDropped` verwenden
- **INAKOMPATIBLE ÄNDERUNG:** HTML-Attribut von `data-wails-dropzone` in `data-file-drop-target` geändert
- **INAKOMPATIBLE ÄNDERUNG:** CSS-Hover-Klasse von `wails-dropzone-hover` in `file-drop-target-active` geändert
- **INAKOMPATIBLE ÄNDERUNG:** Optionen `DragEffect`, `OnEnterEffect` und `OnOverEffect` aus Windows entfernt (sie gehörten zum entfernten IDropTarget)
- Für die gesamte JSON-Verarbeitung zur Laufzeit (Methodenbindungen, Ereignisse, WebView-Anfragen, Benachrichtigungen, kvstore) auf goccy/go-json umgestellt; dadurch wurde die Leistung um 21-63 % verbessert und die Anzahl der Speicherallokationen um 40-60 % reduziert
- Layout der BoundMethod-Struktur optimiert und isVariadic-Flag zwischengespeichert, um den Aufwand pro Aufruf zu reduzieren
- Für Methoden mit `<=8` Argumenten einen auf dem Stack allozierten Argumentpuffer verwendet, um Heap-Allokationen zu vermeiden
- Erfassung der Ergebnisse von Methodenaufrufen optimiert, um bei einzelnen Rückgabewerten keine Slice-Allokation vorzunehmen
- sync.Map für den MIME-Typ-Cache verwendet, um die Leistung bei gleichzeitigem Zugriff zu verbessern
- Pufferpool zum Lesen des Request-Bodys beim HTTP-Transport verwendet
- CloseNotify-Kanal im Inhaltstyp-Sniffer verzögert alloziert, um die Anzahl der Allokationen pro Anfrage zu reduzieren
- CSS-Debug-Protokollierung aus dem Asset-Server entfernt
- Zuordnung von MIME-Typen zu Dateierweiterungen auf über 50 gängige Webformate erweitert (Schriftarten, Audio, Video usw.)
- Dokumentation der Window-Optionen `X/Y` aktualisiert von @ruhuang2001
- Dokumentation zu `Frontend Runtime` um weitere Optionen zum Generieren von Frontend-Bindings ergänzt von @ndianabasi
- Dokumentationsseite zum Wails-Asset-Server v3 aktualisiert von @ndianabasi
- **INAKOMPATIBLE ÄNDERUNG**: Dialogfunktionen auf Paketebene (`application.InfoDialog()`, `application.QuestionDialog()` usw.) entfernt. Stattdessen den Manager `app.Dialog` verwenden: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Dialogdokumentation an die tatsächliche API angepasst: `app.Dialog.*` und `AddButton()` mit Callbacks verwenden (nicht `SetButtons()`), außerdem `SetDefaultButton(*Button)` (kein String), `AddFilter()` (nicht `SetFilters()`), `SetFilename()` (nicht `SetDefaultFilename()`) sowie `app.Dialog.OpenFile().CanChooseDirectories(true)` zur Ordnerauswahl
- **INAKOMPATIBLE ÄNDERUNG**: Produktions-Builds sind jetzt die Standardeinstellung. Für Entwicklungs-Builds `DEV=true` in den Taskfiles festlegen. Für Beispiele ein neues Projekt generieren von @leaanthony
- Beim Auslösen eines benutzerdefinierten Ereignisses mit keinem oder einem Datenargument wird der Datenwert direkt dem Feld Data zugewiesen, ohne ihn in ein Slice einzuschließen, von [@fbbdev](https://github.com/fbbdev) in [#4633](https://github.com/wailsapp/wails/pull/4633)
- Windows-Tray-Symbole berücksichtigen jetzt `SystemTray.Show()`/`Hide()`, indem sie `NIS_HIDDEN` umschalten, sodass Anwendungen vollständig verschwinden und wieder angezeigt werden können (#4653).
- Die Tray-Registrierung verwendet aufgelöste Symbole erneut, setzt `NOTIFYICON_VERSION_4` einmalig und aktiviert `NIF_SHOWTIP`, sodass Tooltips nach einem Neustart des Explorers wiederhergestellt werden (#4653).
- macOS: Für die Fensterzentrierung `visibleFrame` statt `frame` verwendet, um die Bereiche von Menüleiste und Dock auszuschließen
- macOS: Für die Fensterzentrierung `visibleFrame` statt `frame` verwendet, um die Bereiche von Menüleiste und Dock auszuschließen
- Beim Ausführen von `wails3 update build-assets` mit dem Parameter `-config` werden die über die Parameter `-product*` festgelegten Werte
- `window.NativeWindowHandle()` -> `window.NativeWindow()` von @leaanthony in [#4471](https://github.com/wailsapp/wails/pull/4471)
- Interne Fensterverwaltung überarbeitet von @leaanthony in [#4471](https://github.com/wailsapp/wails/pull/4471)
- `application.WindowIDKey` und `application.WindowNameKey` entfernt (durch `application.WindowKey` ersetzt) von [@leaanthony](https://github.com/leaanthony)
- ContextMenuData gibt jetzt einen String statt any zurück von [@leaanthony](https://github.com/leaanthony)
- In JS/TS-Bindings werden Klassenfelder mit Array-Typen fester Länge jetzt mit ihrer erwarteten Länge initialisiert, statt leer zu bleiben, von [@fbbdev](https://github.com/fbbdev) in [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData gibt jetzt einen String statt any zurück von [@leaanthony](https://github.com/leaanthony)
- `application.NewService` akzeptiert Optionen nicht mehr als optionalen Parameter (stattdessen `application.NewServiceWithOptions` verwenden) von [@leaanthony](https://github.com/leaanthony) in [#4024](https://github.com/wailsapp/wails/pull/4024)
- Abhängigkeit von `nanoid` entfernt von [@leaanthony](https://github.com/leaanthony)
- Window-Beispiel für Mica-, Acryl- und Tabbed-Fensterstile aktualisiert von [@leaanthony](https://github.com/leaanthony)
- In JS/TS-Bindings wurden die Modelldateien `internal.js/ts` entfernt; alle Modelle befinden sich jetzt in `models.js/ts`, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- In JS/TS-Bindings werden benannte Typen niemals als Aliase für andere benannte Typen dargestellt; das bisherige Verhalten ist jetzt auf Aliase beschränkt, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- In JS/TS-Bindings im Klassenmodus werden Strukturfelder, deren Typ ein Typparameter ist, als optional gekennzeichnet und niemals automatisch initialisiert, von [@fbbdev](https://github.com/fbbdev) in [#4045](https://github.com/wailsapp/wails/pull/4045)
- ESLint aus den Vorlagen entfernt von [@IanVS](https://github.com/IanVS) in [#4059](https://github.com/wailsapp/wails/pull/4059)
- Copyright-Jahr auf 2025 aktualisiert von [@IanVS](https://github.com/IanVS) in [#4037](https://github.com/wailsapp/wails/pull/4037)
- Dokumentation für event.Sender hinzugefügt von [@IanVS](https://github.com/IanVS) in [#4075](https://github.com/wailsapp/wails/pull/4075)
- Unterstützung für Go 1.24 von [@leaanthony](https://github.com/leaanthony)
- `ServiceStartup`-Hooks werden jetzt beim Aufruf von `App.Run` statt in `application.New` ausgeführt, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- `ServiceStartup`-Fehler werden jetzt von `App.Run` zurückgegeben, statt den Prozess zu beenden, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Binding- und Dialogaufrufe aus JS werden jetzt mit Fehlerobjekten statt Zeichenfolgen abgelehnt, von [@fbbdev](https://github.com/fbbdev) in [#4066](https://github.com/wailsapp/wails/pull/4066)
- Positionierung des Taskleistenmenüs unter Windows verbessert von [@leaanthony](https://github.com/leaanthony)
- Die JS-Laufzeit wurde von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100) nach TypeScript portiert
- Die Laufzeit wird direkt beim Import initialisiert; auf das Laden des Fensters muss nicht mehr gewartet werden, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Die Laufzeit exportiert keine init-Methode mehr. Sie kann mit einem Import für Seiteneffekte initialisiert werden, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Gebundene Methoden geben jetzt ein `CancellablePromise` zurück, das bei einem Abbruch mit einem `CancelError` abgelehnt wird. Das eigentliche Ergebnis des Aufrufs wird verworfen, von [@fbbdev](https://github.com/fbbdev) in [#4100](https://github.com/wailsapp/wails/pull/4100)
- Integrierte Diensttypen heißen jetzt einheitlich `Service`, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Funktionen zum Erstellen integrierter Dienste mit Optionen heißen jetzt einheitlich `NewWithConfig`, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Die Methode `Select` des Dienstes `sqlite` heißt jetzt `Query`, um mit den Go-APIs übereinzustimmen, von [@fbbdev](https://github.com/fbbdev) in [#4067](https://github.com/wailsapp/wails/pull/4067)
- Vorlagen: Laufzeit nach „dependencies“ verschoben und package.json-Dateien organisiert von [@IanVS](https://github.com/IanVS) in [#4133](https://github.com/wailsapp/wails/pull/4133)
- Erstellt und signiert App-Bundles während der Entwicklung ad hoc, um bestimmte macOS-APIs zu aktivieren, von [@popaprozac](https://github.com/popaprozac) in [#4171](https://github.com/wailsapp/wails/pull/4171)
- Build-Assets in plattformspezifische Verzeichnisse verschoben von [@leaanthony](https://github.com/leaanthony)
- Taskfiles in plattformspezifische Verzeichnisse verschoben und umbenannt von [@leaanthony](https://github.com/leaanthony)
- Die Benutzererfahrung bei fehlendem `index.html` deutlich verbessert von [@leaanthony](https://github.com/leaanthony)
- [Windows] Leistung beim Minimieren und Wiederherstellen verbessert von [@leaanthony](https://github.com/leaanthony). Basierend auf dem ursprünglichen [PR](https://github.com/wailsapp/wails/pull/3955) von [562589540](https://github.com/562589540)
- Option `ShouldClose` entfernt (registrieren Sie stattdessen einen Hook für events.Common.WindowClosing) von [@leaanthony](https://github.com/leaanthony)
- [Windows] Flackern beim Öffnen eines Fensters reduziert von [@leaanthony](https://github.com/leaanthony)
- `Window.Destroy` entfernt, da dies als interne Funktion vorgesehen war, von [@leaanthony](https://github.com/leaanthony)
- `WindowClose`-Ereignisse in `WindowClosing` umbenannt von [@leaanthony](https://github.com/leaanthony)
- Frontend-Builds verwenden jetzt abhängig vom Build-Typ die Vite-Umgebung „development“ oder „production“, von [@leaanthony](https://github.com/leaanthony)
- Aktualisierung auf go-webview2 v1.19 von [@leaanthony](https://github.com/leaanthony)
- Sichergestellt, dass der Fork von taskfile verwendet wird, von @leaanthony
- Fork von Taskfile aktualisiert, um Versionsprobleme bei der Installation mit zu beheben
- Fork von Taskfile verwendet, um Versionsprobleme bei der Installation mit zu beheben
- `service.OnStartup` beendet die Anwendung jetzt bei einem Fehler und führt aus
- Nachrichtenübermittlung bei Klicks auf das Taskleistensymbol überarbeitet, um besser den Benutzerinteraktionen zu entsprechen, von
- Asset-Einbettung um `all:frontend/dist` erweitert, um Frameworks zu unterstützen, die generieren
- Taskfile überarbeitet von [leaanthony](https://github.com/leaanthony) in
- Upgrade von `go-webview2` auf v1.0.16 von
- Typ `Screen` korrigiert, sodass er `ID` statt `Id` enthält, von
- Wails-Version von `go.mod.tmpl` aktualisiert, um `application.ServiceOptions` zu unterstützen, von
- Ermittlung des Dienstnamens korrigiert von [windom](https://github.com/windom/) in
- mkdocs serve verwendet jetzt Docker, von [leaanthony](https://github.com/leaanthony)
- Entwicklungskonfiguration in `config.yml` zusammengeführt von
- Der Taskleistendialog verwendet jetzt standardmäßig das Anwendungssymbol, sofern verfügbar (Windows), von
- Verbesserte Meldung von GPU und Arbeitsspeicher unter macOS von
- `WebviewGpuIsDisabled` und `EnableFraudulentWebsiteWarnings` entfernt
- Änderung der Events-API: `On`/`Emit` -> Benutzerereignisse, `OnApplicationEvent` ->
- Events-API unter Linux korrigiert von [TheGB0077](https://github.com/TheGB0077) in
- [CI] Actions verbessert und ihre Ausführung auch in Forks aktiviert sowie
- `AbsolutePosition()` in `Position()` umbenannt von
- Linux-WebKit-Abhängigkeit von webkitgtk2-4.0 auf webkit2gtk-4.1 aktualisiert, um
- Das gebündelte JS-Laufzeitskript ist jetzt ein ESM-Modul: script-Tags, die es importieren
- Das Paket `@wailsio/runtime` veröffentlicht seine API nicht im `window.wails`
- Das Window-API-Modul `@wailsio/runtime/src/window` stellt jetzt das enthaltende
- Die JS-Window-API wurde an das aktuelle Go-`WebviewWindow` angepasst
- Der Binding-Generator verwendet jetzt standardmäßig Aufrufe per ID. Die CLI-Option `-id`
- Neues Layout des Binding-Codes: Ausgabedateien waren zuvor in Ordnern organisiert
- Das Struct-Feld `application.Options.Bind` wurde umbenannt in
- Neue Syntax für Binding-Dienste: Dienstinstanzen müssen jetzt eingeschlossen werden in ein
- Spinner in Nicht-Terminal- oder CI-Umgebungen deaktiviert von

### Entfernt

- **BREAKING**: `EnabledFeatures`, `DisabledFeatures` und `AdditionalLaunchArgs` aus den fensterspezifischen `WindowsWindow`-Optionen entfernt. Verwenden Sie stattdessen die anwendungsweiten Optionen `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures` und `Options.Windows.AdditionalBrowserArgs`. Diese Flags gelten global für die gemeinsam genutzte WebView2-Umgebung (#4559), von @leaanthony
- Native `IDropTarget`-Implementierung unter Windows zugunsten eines JavaScript-basierten Ansatzes entfernt (entspricht dem Verhalten von v2)
- Abhängigkeit github.com/wailsapp/mimetype zugunsten einer erweiterten Zuordnung von Dateierweiterungen und http.DetectContentType aus der Standardbibliothek entfernt; dadurch wurde die Binärgröße um etwa 1.2 MB reduziert
- Entfernt die Abhängigkeit gopkg.in/ini.v1 durch Implementierung eines minimalen Parsers für .desktop-Dateien im Linux-Dateiexplorer und spart dadurch ~45 KB
- Entfernt samber/lo aus dem Laufzeitcode durch Verwendung des stdlib-Pakets slices aus Go 1.21+ und minimaler interner Hilfsfunktionen und spart dadurch ~310 KB
- Entfernt Debug-printf-Anweisungen aus dem Darwin-Handler für URL-Schemata (#4834)
- **NICHT ABWÄRTSKOMPATIBEL**: Entfernt das Ereignis `linux:WindowLoadChanged`. Verwenden Sie stattdessen `linux:WindowLoadFinished`, um zu erkennen, wann die WebView vollständig geladen ist (#3896), von @leaanthony

### Nicht abwärtskompatible Änderungen

- **Überarbeitung der Manager-API**: Die Anwendungs-API wurde zur besseren Codeorganisation und leichteren Auffindbarkeit von einer flachen Struktur in übersichtliche Manager umstrukturiert, von [@leaanthony](https://github.com/leaanthony) in [#4359](https://github.com/wailsapp/wails/pull/4359)
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
- Service-Methoden umbenannt: `Name` -> `ServiceName`, `OnStartup` -> `ServiceStartup`, `OnShutdown` -> `ServiceShutdown`, von [@leaanthony](https://github.com/leaanthony)
- Die Methoden `Path` und `Paths` wurden von [@leaanthony](https://github.com/leaanthony) in das Paket `application` verschoben
- Das Anwendungsmenü ist jetzt nur noch unter macOS verfügbar, von [@leaanthony](https://github.com/leaanthony)

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

+ Zu weit gefasste Linux-Paketabhängigkeiten und veraltete RPM-Abhängigkeiten behoben.

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

[mmghv](https://github.com/mmghv) in   [#3611](https://github.com/wailsapp/wails/pull/3611)   unterstützt Ubuntu 24.04 LTS, von [atterpac](https://github.com/atterpac) in   [#3461](https://github.com/wailsapp/wails/pull/3461)   muss das Attribut `type="module"` besitzen. Von   [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   Objekt und startet das WML-System nicht. Dies verbessert die Kapselung. Das WML-System kann bei Bedarf manuell durch Aufruf der neuen Methode `WML.Enable` gestartet werden. Das gebündelte JS-Laufzeitskript führt weiterhin beide Vorgänge automatisch aus. Von [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   Fensterobjekt als Standardexport. Einzelne Methoden können nicht mehr über benannte ESM-Importe oder die Namespace-Importsyntax importiert werden.   API. Einige Methoden wurden umbenannt oder ihre Signatur wurde geändert, insbesondere: `Screen`   wird zu `GetScreen`; `GetZoomLevel`/`SetZoomLevel` werden zu `GetZoom`/`SetZoom`;   `GetZoom`, `Width` und `Height` geben Werte nun direkt zurück, statt sie in Objekte einzuschließen. Von [@fbbdev](https://github.com/fbbdev) in   [#3295](https://github.com/wailsapp/wails/pull/3295)   wurde entfernt. Verwenden Sie die CLI-Option `-names`, um wieder zu Aufrufen nach Namen zu wechseln.   Von [@fbbdev](https://github.com/fbbdev) in   [#3468](https://github.com/wailsapp/wails/pull/3468)   wurden nach ihrem übergeordneten Paket benannt; nun werden vollständige Go-Importpfade einschließlich des Modulpfads verwendet. Von [@fbbdev](https://github.com/fbbdev) in   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Von [@fbbdev](https://github.com/fbbdev) in   [#3468](https://github.com/wailsapp/wails/pull/3468)   Aufruf von `application.NewService`. Von [@fbbdev](https://github.com/fbbdev) in   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) in   [#3574](https://github.com/wailsapp/wails/pull/3574)
