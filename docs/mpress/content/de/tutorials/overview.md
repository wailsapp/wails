---
title: "Tutorials"
description: "Wails durch das Erstellen von Anwendungen kennenlernen"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

Wähle eine geführte Lektion, wenn du Wails kennenlernen möchtest, indem du etwas erstellst. Jede Lektion führt ein Projekt zu einem konkreten Ergebnis.

Wenn du zum ersten Mal hier bist, schließe zuerst [Installation](/quick-start/installation/) und [Deine erste App](/quick-start/first-app/) ab.

## Tutorial auswählen

Arbeite nach deiner ersten App die Desktop-Lektionen in dieser Reihenfolge durch:

| Lektion | Was du üben wirst |
| --- | --- |
| [QR-Code-Dienst](/tutorials/01-creating-a-service/) | Go-Logik über einen Dienst bereitstellen und aus dem Frontend aufrufen |
| [TODO-Liste](/tutorials/02-todo-vanilla/) | Anwendungsdaten erstellen, aktualisieren und löschen |
| [Notizen-App](/tutorials/03-notes-vanilla/) | Native Dialoge verwenden und Daten in Dateien speichern |
| [App mit automatischen Updates](/tutorials/04-self-update-a-wails-app/) | Updates aus GitHub Releases suchen und anwenden |

Arbeite für eine optionale Mobilgeräte-Lektion [Deine erste mobile App](/guides/mobile/first-mobile-app/) durch. Die Unterstützung für Mobilgeräte ist experimentell.

Wenn du an einer bestehenden App arbeitest und eine bestimmte Aufgabe erledigen musst, wähle stattdessen eine [Schritt-für-Schritt-Anleitung](/quick-start/next-steps/#features). Hintergrundinformationen zum Design findest du unter [Kernkonzepte](/concepts/architecture/).

@note{type="tip" title="Go ist neu für dich?"}
Absolviere die [Go Tour](https://go.dev/tour/), bevor du mit den Tutorials beginnst.

@end

## QR-Code-Dienst

![QR-Code-Beispiel](/assets/qr1.png)

Lerne die Grundlagen von Wails-Diensten kennen, indem du einen QR-Code-Generator erstellst. Dieses Tutorial führt dich in die Kernkonzepte zur Organisation deiner Anwendungslogik in wiederverwendbaren Diensten ein.

**Das wirst du lernen:**

- Einen Wails-Dienst erstellen und strukturieren
- Externe Go-Abhängigkeiten verwalten
- Go-Methoden an dein Frontend binden
- Daten zwischen Go und JavaScript übertragen
- Code wartungsfreundlich organisieren

**Am besten geeignet für:** Personen, die Wails zum ersten Mal verwenden und die Dienstarchitektur verstehen möchten

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>Loslegen</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### TODO-Liste

![TODO-Listen-Anwendung](/assets/todo-app.png)

Erstelle eine vollständige TODO-Listen-Anwendung mit einer ansprechenden, modernen Oberfläche. In diesem praxisorientierten Tutorial lernst du zentrale Wails-Muster anhand einer realitätsnahen Anwendung mit Vanilla JavaScript kennen.

**Das wirst du lernen:**

- Dienstbasierte Architektur mit threadsicherer Zustandsverwaltung
- CRUD-Operationen (Erstellen, Lesen, Aktualisieren, Löschen)
- Typsichere Bindungen zwischen Go und JavaScript
- Moderne Benutzeroberflächen ohne die Komplexität eines Frameworks erstellen
- Geeignete Muster für Fehlerbehandlung und Validierung

**Benötigte Zeit:** ca. 20 Minuten

**Am besten geeignet für:** Deine erste vollständige Wails-Anwendung – ideal, um die Grundlagen zu verstehen, bevor du die Komplexität eines Frameworks hinzufügst

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>Loslegen</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Notizen

![Notizen-Anwendung](/assets/notes-app.png)

Erstelle eine Anwendung im Stil von Apple Notes mit nativen Dateidialogen und automatischer Speicherung. Dieses Tutorial veranschaulicht Desktop-spezifische Funktionen wie Dateioperationen, native Dialoge und professionelle UI-Muster.

**Das wirst du lernen:**

- Native Dateidialoge (Speichern, Öffnen, Info)
- JSON-basierte Datenpersistenz
- Entprellte Muster für automatisches Speichern
- Professionelle zweispaltige Desktop-Layouts
- Mit Dateisystemoperationen in Go arbeiten

**Benötigte Zeit:** ca. 30 Minuten

**Am besten geeignet für:** Desktop-spezifische Funktionen wie Dateioperationen und native Betriebssystemdialoge kennenlernen

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>Loslegen</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Wails-App mit automatischen Updates

Füge einer Wails-Anwendung ausgehend von einer neuen `wails3 init` Selbstaktualisierungen innerhalb der App hinzu – einschließlich der Überprüfung signierter Releases und des Binärdateiaustauschs im Hilfsmodus. Als Updatequelle dienen GitHub Releases.

**Das wirst du lernen:**

- Wie `app.Updater` in eine Wails-App eingebunden wird
- Den GitHub-Releases-Provider konfigurieren
- Releases mit `SHA256SUMS` zur Digest-Verifizierung veröffentlichen
- Ed25519-Signaturen zum Schutz vor Manipulationen hinzufügen
- Das Standardfenster über CSS, benutzerdefiniertes HTML oder BYO anpassen
- Regelmäßige Hintergrundprüfungen mit `CheckInterval`

**Benötigte Zeit:** ca. 25 Minuten

**Am besten geeignet für:** Die Auslieferung einer aktualisierbaren Desktop-App – behandelt die vollständige Release-Pipeline, nicht nur die API

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>Loslegen</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
