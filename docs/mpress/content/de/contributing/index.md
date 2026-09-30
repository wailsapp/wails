---
title: "Mitwirken"
description: "Zu Wails beitragen"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## Willkommen, Mitwirkende!

Wir freuen uns über Beiträge zu Wails! Ob Sie Fehler beheben, Funktionen hinzufügen oder die Dokumentation verbessern – Ihre Hilfe ist willkommen.

## Möglichkeiten zur Mitwirkung

### 1. Probleme melden

Einen Fehler gefunden? [Erstellen Sie einen Issue](https://github.com/wailsapp/wails/issues/new) mit folgenden Angaben:

- Klare Beschreibung
- Schritte zum Reproduzieren
- Erwartetes und tatsächliches Verhalten
- Systeminformationen
- Codebeispiele

### 2. Dokumentation verbessern

PRs mit Korrekturen sind auch ohne vorherigen Issue oder fehlgeschlagenen Codetest willkommen.  
Folgen Sie der Anleitung [Dokumentation korrigieren](/contributing/documentation/), um eine Änderung mit M-Press in der Vorschau anzuzeigen und zu validieren.

Verbesserungen der Dokumentation sind jederzeit willkommen:

- Tippfehler und Fehler beheben
- Beispiele hinzufügen
- Erklärungen verdeutlichen
- Inhalte übersetzen

### 3. Code beitragen

Tragen Sie Code über Pull Requests bei:

- Fehlerbehebungen
- Neue Funktionen
- Leistungsverbesserungen
- Tests

### 4. Verbesserung vorschlagen (WEP)

Neue Funktionen und Änderungen am öffentlichen Verhalten durchlaufen den Prozess für Wails Enhancement Proposals (WEP). Er hält die Funktionsentwicklung transparent und stellt sicher, dass es für jeden angenommenen Vorschlag eine für die Implementierung zuständige Person gibt. Erstellen Sie keinen Issue für einen Funktionswunsch.

1. Optional können Sie Ihre Idee in der Kategorie [Ideas](https://github.com/wailsapp/wails/discussions/categories/ideas) der GitHub Discussions oder auf [Discord](https://discord.gg/JDdSxwjhGf) vorstellen, um das Interesse auszuloten.
2. Kopieren Sie [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/WEP_TEMPLATE.md) nach `v3/wep/proposals/<proposal name>/proposal.md` und füllen Sie jeden Abschnitt aus.
3. Erstellen Sie einen als Entwurf markierten Pull Request mit dem Titel `[WEP] <title>`, der ausschließlich den Vorschlag enthält. Der PR ist der offizielle Ort für dessen Diskussion.
4. Holen Sie Feedback und Unterstützung ein (Kommentare und „Daumen hoch“-Reaktionen im PR). Planen Sie mindestens zwei Wochen für die Diskussion ein und einigen Sie sich darauf, wer den Vorschlag implementiert.
5. Markieren Sie den PR als bereit zur Überprüfung. Die Maintainer treffen die endgültige Entscheidung: Angenommene Vorschläge erhalten eine WEP-Nummer und werden zusammengeführt.

Der vollständige Prozess ist in [`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md) dokumentiert.

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

Tests sind Bestandteil der Änderung und keine abschließende oberflächliche Prüfung. Platzieren Sie Unit-Tests neben dem Code, den sie prüfen, und verwenden Sie bevorzugt tabellengesteuerte Tests, wenn ein Verhalten mit mehreren Eingaben oder Grenzfällen geprüft wird. Benennen Sie jeden Testfall so, dass bei einem Fehlschlag das Szenario erkennbar ist.

Neue und geänderte Logik sollte vollständig abgedeckt sein. Streben Sie für den durch Ihren PR hinzugefügten oder geänderten Code eine Go-Anweisungsabdeckung von 100 % an; betrachten Sie einen prozentualen Wert für das gesamte Repository nicht als Ersatz für das Testen der Änderung. Eine Lücke kann gerechtfertigt sein – beispielsweise bei einem nur unter einem bestimmten Betriebssystem auftretenden Fehlerpfad oder einer Bedingung, die sich ohne echte Hardware nicht praktikabel reproduzieren lässt. Erläutern Sie in der PR-Beschreibung jedoch die Lücke und warum sie sich nicht mit vertretbarem Aufwand testen lässt.

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

Informationen zu Integrations-Test-Suites, Race Detection und den vollständigen CI-äquivalenten Befehlen finden Sie unter [Tests und kontinuierliche Integration](/contributing/testing-ci/).

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

1. **Code schreiben** und dabei die Go-Konventionen einhalten
2. Für neue Funktionen **Tests hinzufügen**
3. Bei Bedarf die **Dokumentation aktualisieren**
4. **Tests ausführen**, um sicherzustellen, dass nichts beschädigt wird
5. **Änderungen committen** und klare Meldungen verwenden

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

## Pull-Request-Richtlinien

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
- [ ] Tests hinzugefügt oder aktualisiert
- [ ] Dokumentation aktualisiert
- [ ] Alle Tests erfolgreich
- [ ] Keine Breaking Changes (oder dokumentiert)
- [ ] Commit-Meldungen sind klar

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

### Dokumentation verfassen

Die Dokumentation verwendet M-Press. Bearbeite die `.md`-Dateien unter  
`docs/mpress/content/` und führe anschließend die Vorschau und Validierung aus dem Stammverzeichnis des Repositorys aus:

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

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

Sei respektvoll, inklusiv und professionell. Wir alle sind hier, um gemeinsam großartige Software zu entwickeln. Weitere Einzelheiten findest du im [Verhaltenskodex](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md).

## Anerkennung

Mitwirkende werden hier gewürdigt:

- Versionshinweise
- Liste der Mitwirkenden
- GitHub Insights

Vielen Dank für deinen Beitrag zu Wails! 🎉

## Nächste Schritte

@cards{cols="2"}
◆ GitHub-Repository
Besuche das Wails-Repository.

[Auf GitHub ansehen →](https://github.com/wailsapp/wails)

---
◆ Discord-Community
Tritt der Community bei.

[Discord beitreten →](https://discord.gg/JDdSxwjhGf)

---
📖 Dokumentation
Lies die Dokumentation.

[Dokumentation durchsuchen →](/quick-start/why-wails/)

@end
