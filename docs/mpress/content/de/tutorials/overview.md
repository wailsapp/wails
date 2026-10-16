---
title: "Tutorials"
description: "Wails durch die Entwicklung von Anwendungen kennenlernen"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

Schritt-für-Schritt-Tutorials vermitteln Wails-Konzepte, indem vollständige Anwendungen entwickelt werden. Jedes Tutorial enthält funktionsfähigen Code, Erklärungen und praxisnahe Muster.

@note{type="tip" title="Neu bei Go?"}
Absolvieren Sie die [Go Tour](https://go.dev/tour/), bevor Sie mit den Tutorials beginnen.

@end

## QR-Code-Dienst

![QR-Code-Beispiel](/assets/qr1.png)

Lernen Sie die Grundlagen von Wails-Diensten kennen, indem Sie einen QR-Code-Generator entwickeln. Dieses Tutorial führt Sie in die zentralen Konzepte zur Strukturierung Ihrer Anwendungslogik in wiederverwendbare Dienste ein.

**Lerninhalte:**

- Erstellen und Strukturieren eines Wails-Dienstes
- Verwalten externer Go-Abhängigkeiten
- Binden von Go-Methoden an Ihr Frontend
- Übertragen von Daten zwischen Go und JavaScript
- Strukturieren von Code für eine bessere Wartbarkeit

**Ideal für:** Neue Wails-Benutzer, die die Dienstarchitektur verstehen möchten

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>Loslegen</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### TODO-Liste

![TODO-Listen-Anwendung](/assets/todo-app.png)

Entwickeln Sie eine vollständige TODO-Listen-Anwendung mit einer ansprechenden, modernen Benutzeroberfläche. Dieses praxisorientierte Tutorial vermittelt Ihnen zentrale Wails-Muster anhand einer realistischen Anwendung mit Vanilla JavaScript.

**Lerninhalte:**

- Dienstbasierte Architektur mit threadsicherer Zustandsverwaltung
- CRUD-Operationen (Erstellen, Lesen, Aktualisieren, Löschen)
- Typsichere Bindungen zwischen Go und JavaScript
- Entwickeln moderner Benutzeroberflächen ohne die Komplexität eines Frameworks
- Korrekte Fehlerbehandlung und Validierungsmuster

**Bearbeitungszeit:** ca. 20 Minuten

**Ideal für:** Ihre erste vollständige Wails-Anwendung – perfekt, um die Grundlagen zu verstehen, bevor Sie die Komplexität eines Frameworks hinzufügen

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>Loslegen</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Notizen

![Notizen-Anwendung](/assets/notes-app.png)

Entwickeln Sie eine Anwendung im Stil von Apple Notes mit nativen Dateidialogen und automatischer Speicherfunktion. Dieses Tutorial veranschaulicht Desktop-spezifische Funktionen wie Dateioperationen, native Dialoge und professionelle UI-Muster.

**Lerninhalte:**

- Native Dateidialoge (Speichern, Öffnen, Info)
- JSON-basierte Datenpersistenz
- Entprellte Muster für automatisches Speichern
- Professionelle zweispaltige Desktop-Layouts
- Arbeiten mit Dateisystemoperationen in Go

**Bearbeitungszeit:** ca. 30 Minuten

**Ideal für:** Das Erlernen Desktop-spezifischer Funktionen wie Dateioperationen und nativer Betriebssystemdialoge

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>Loslegen</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Wails-App mit automatischen Updates

Fügen Sie einer Wails-Anwendung ausgehend von einem neuen `wails3 init` anwendungsinterne automatische Updates hinzu – einschließlich der Überprüfung signierter Releases und des Austauschs der Binärdatei im Hilfsmodus. Als Update-Quelle dienen GitHub Releases.

**Lerninhalte:**

- Einbinden von `app.Updater` in eine Wails-App
- Konfigurieren des GitHub-Releases-Providers
- Veröffentlichen von Releases mit `SHA256SUMS` zur Digest-Überprüfung
- Hinzufügen einer Ed25519-Signatur zum Schutz vor Manipulationen
- Anpassen des Standardfensters über CSS, benutzerdefiniertes HTML oder BYO
- Regelmäßige Hintergrundprüfungen mit `CheckInterval`

**Bearbeitungszeit:** ca. 25 Minuten

**Ideal für:** Die Auslieferung einer aktualisierbaren Desktop-App – behandelt die vollständige Release-Pipeline, nicht nur die API

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>Loslegen</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
