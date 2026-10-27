---
title: "Mitwirken"
description: "Zu Wails beitragen"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## Willkommen, Mitwirkende!

Wir freuen uns über Beiträge zu Wails! Ob du Fehler behebst, Funktionen hinzufügst oder die Dokumentation verbesserst – deine Hilfe ist willkommen.

## Möglichkeiten zur Mitwirkung

### 1. Probleme melden

Du hast einen Fehler gefunden? [Erstelle ein Issue](https://github.com/wailsapp/wails/issues/new) mit folgenden Angaben:

- Klare Beschreibung
- Schritte zum Reproduzieren
- Erwartetes und tatsächliches Verhalten
- Systeminformationen
- Codebeispiele

### 2. Dokumentation verbessern

Korrektur-PRs sind auch ohne vorheriges Issue oder fehlgeschlagenen Codetest willkommen.  
Informationen zu Downloads, zur Installation unter macOS/Linux/Windows, zur Live-Vorschau und zum Einreichen deiner Änderung findest du unter [Dokumentation mit M-Press schreiben](/contributing/documentation/).

Verbesserungen an der Dokumentation sind jederzeit willkommen:

- Tippfehler und Fehler beheben
- Beispiele hinzufügen
- Erklärungen präzisieren
- Inhalte übersetzen

### 3. Code einreichen

Reiche Code über Pull Requests ein:

- Fehlerbehebungen
- Neue Funktionen
- Leistungsverbesserungen
- Tests

### 4. Eine Erweiterung vorschlagen (WEP)

Für neue Funktionen und Änderungen am öffentlichen Verhalten gilt der Prozess für Wails Enhancement Proposals (WEPs). Er macht die Funktionsentwicklung transparent und stellt sicher, dass jeder angenommene Vorschlag von jemandem implementiert wird. Erstelle kein Issue für einen Funktionswunsch.

1. Du kannst deine Idee optional in der Kategorie [Ideas](https://github.com/wailsapp/wails/discussions/categories/ideas) in GitHub Discussions oder auf [Discord](https://discord.gg/JDdSxwjhGf) vorstellen, um das Interesse daran einzuschätzen.
2. Kopiere [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/WEP_TEMPLATE.md) nach `v3/wep/proposals/<proposal name>/proposal.md` und fülle jeden Abschnitt aus.
3. Erstelle einen als Entwurf markierten Pull Request mit dem Titel `[WEP] <title>`, der ausschließlich den Vorschlag enthält. Der PR ist der offizielle Ort für dessen Diskussion.
4. Hole Feedback und Unterstützung ein (Kommentare und „Daumen hoch“-Reaktionen im PR). Plane mindestens zwei Wochen für die Diskussion ein und einigt euch darauf, wer den Vorschlag implementiert.
5. Markiere den PR als bereit zur Überprüfung. Die Maintainer treffen die endgültige Entscheidung: Angenommene Vorschläge erhalten eine WEP-Nummer und werden zusammengeführt.

Der vollständige Prozess ist in [`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md) dokumentiert.

## Erste Schritte

### Fork erstellen und klonen

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git
```

### Aus dem Quellcode erstellen

```bash
# Build the v3 CLI. Go downloads any required modules automatically.
cd v3
go build -o ../wails3 ./cmd/wails3

# Confirm the built CLI runs and reports its version.
../wails3 version
```

### Tests ausführen

Tests sind Bestandteil der Änderung und keine abschließende oberflächliche Prüfung. Platziere Unit-Tests neben dem getesteten Code und bevorzuge tabellengesteuerte Tests, wenn ein Verhalten mit mehreren Eingaben oder Grenzfällen geprüft wird. Benenne jeden Testfall so, dass ein Fehlschlag das jeweilige Szenario erkennen lässt.

Neue und geänderte Logik sollte vollständig abgedeckt sein. Strebe für den von deinem PR hinzugefügten oder geänderten Code eine Go-Anweisungsabdeckung von 100 % an. Betrachte einen repositoryweiten Prozentsatz nicht als Ersatz für das Testen der Änderung. Eine Lücke kann gerechtfertigt sein – beispielsweise bei einem betriebssystemspezifischen Fehlerpfad oder einer Bedingung, die sich ohne echte Hardware nicht praktikabel reproduzieren lässt. Erläutere dann aber in der PR-Beschreibung die Lücke und warum sie sich nicht mit vertretbarem Aufwand testen lässt.

```bash
# Run all v3 tests
cd v3
go test ./...

# Run specific package tests
go test ./pkg/application

# Inspect coverage for the packages you changed
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
cd ..
```

Informationen zu Integrations-Testreihen, Race Detection und den vollständigen CI-äquivalenten Befehlen findest du unter [Tests und Continuous Integration](/contributing/testing-ci/).

## Änderungen vornehmen

### Branch erstellen

```bash
# Update master
git checkout master
git pull upstream master

# Create feature branch
git checkout -b feature/my-feature
```

### Änderungen vornehmen

1. **Schreibe Code** gemäß den Go-Konventionen
2. **Füge Tests** für neue Funktionen hinzu
3. **Aktualisiere bei Bedarf die Dokumentation**
4. **Führe Tests aus**, um sicherzustellen, dass nichts beeinträchtigt wird
5. **Committe Änderungen** mit klaren Meldungen

### Commit-Richtlinien

```bash
# Good commit messages
git commit -m "fix: resolve window focus issue on macOS"
git commit -m "feat: add support for custom window chrome"
git commit -m "docs: improve bindings documentation"

# Use conventional commits:
# - feat: New feature
# - fix: Bug fix
# - docs: Documentation
# - test: Tests
# - refactor: Code refactoring
# - chore: Maintenance
```

### Pull Request einreichen

```bash
# Push to your fork
git push origin feature/my-feature

# Open pull request on GitHub
# Provide clear description
# Reference related issues
```

## Richtlinien für Pull Requests

### Gute PR-Beschreibung

```markdown
## Description
Brief description of changes

## Changes
- Added feature X
- Fixed bug Y
- Updated documentation

## Testing
- Tested on macOS 14
- Tested on Windows 11
- All tests passing

## Related Issues
Fixes #123
```

### PR-Checkliste

- [ ] Code entspricht den Go-Konventionen
- [ ] Tests hinzugefügt/aktualisiert
- [ ] Dokumentation aktualisiert
- [ ] Alle Tests erfolgreich
- [ ] Keine inkompatiblen Änderungen (oder diese sind dokumentiert)
- [ ] Klare Commit-Meldungen

## Code-Richtlinien

### Go-Codestil

```go
// ✅ Good: Clear, documented, tested
// ProcessData processes the input data and returns the result.
// It returns an error if the data is invalid.
func ProcessData(data string) (string, error) {
    if data == "" {
        return "", errors.New("data cannot be empty")
    }
    
    result := process(data)
    return result, nil
}

// ❌ Bad: No docs, no error handling
func ProcessData(data string) string {
    return process(data)
}
```

### Tests

```go
func TestProcessData(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "processed", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessData(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessData() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ProcessData() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Dokumentation

### Dokumentation schreiben

Den vollständigen Dokumentationsworkflow findest du unter [Dokumentation mit M-Press schreiben](/contributing/documentation/). Er behandelt die festgelegte Version, Installationsoptionen für die ausführbare Datei und Go, bearbeitbare Quelldateien, Navigation, Übersetzungsprüfungen und das Einreichen eines PRs für `master`.

Für reine Dokumentationsarbeiten müssen weder eine Wails-Anwendung erstellt noch Plattformabhängigkeiten installiert werden. M-Press zeigt vom Stammverzeichnis des Repositorys aus mit `mpress dev` eine Vorschau der bestehenden Website an.

### Dokumentationsstil

- Verwende internationale englische Rechtschreibung
- Beginne mit dem Problem
- Stelle funktionsfähige Beispiele bereit
- Füge Hinweise zur Fehlerbehebung hinzu
- Verweise auf verwandte Inhalte

## Community

### Hilfe erhalten

- **Discord:** [Tritt unserer Community bei](https://discord.gg/JDdSxwjhGf)
- **GitHub Discussions:** Stelle Fragen
- **GitHub Issues:** Melde Fehler

### Verhaltenskodex

Sei respektvoll, inklusiv und professionell. Wir sind alle hier, um gemeinsam großartige Software zu entwickeln. Einzelheiten findest du im [Verhaltenskodex](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md).

## Anerkennung

Mitwirkende werden hier gewürdigt:

- Versionshinweise
- Liste der Mitwirkenden
- GitHub Insights

Vielen Dank für deinen Beitrag zu Wails! 🎉

## Nächste Schritte

@cards{cols="2"}
◆ GitHub-Repository
Besuchen Sie das Wails-Repository.

[Auf GitHub ansehen →](https://github.com/wailsapp/wails)

---
◆ Discord-Community
Treten Sie der Community bei.

[Discord beitreten →](https://discord.gg/JDdSxwjhGf)

---
📖 Dokumentation
Lesen Sie die Dokumentation.

[Dokumentation durchsuchen →](/quick-start/why-wails/)

@end
