---
title: "Binding-System"
description: "Wie das Binding-System JavaScript-/TypeScript-Code erfasst, verarbeitet und generiert"
slug: "contributing/architecture/bindings"
sourcePath: "contributing/architecture/bindings.md"
---

Dieser Leitfaden erläutert die interne Funktionsweise des Wails-Binding-Systems und bietet Entwicklern Einblicke in die Mechanismen der automatischen Codegenerierung.

## Architekturüberblick

Das Wails-Binding-System besteht aus drei Hauptkomponenten:

1. **Erfassung**: Analysiert Go-Code, um Informationen über Services, Modelle und andere Deklarationen zu extrahieren
2. **Konfiguration**: Verwaltet Einstellungen und Optionen für die Binding-Generierung
3. **Rendering**: Generiert anhand der erfassten Informationen JavaScript-/TypeScript-Code

@filetree
- internal/generator/
  - collect/     # Paketanalyse und Informationsextraktion
  - config/      # Konfigurationsstrukturen und Schnittstellen
  - render/      # Codegenerierung für JS/TS
@end

## Erfassungsprozess

Der Erfassungsprozess analysiert Go-Pakete und extrahiert Informationen über Services, Modelle und andere Deklarationen. Dies übernimmt das Paket `collect`.

### Hauptkomponenten

- **Collector**: Verwaltet Paketinformationen und speichert erfasste Daten im Cache
- **Package**: Repräsentiert ein analysiertes Go-Paket und speichert erfasste Services, Modelle und Direktiven
- **Service**: Erfasst Informationen über Service-Typen und deren Methoden
- **Model**: Erfasst detaillierte Informationen über Modelltypen, einschließlich Feldern, Werten und Typparametern
- **Directive**: Analysiert und interpretiert `//wails:`-Direktiven im Go-Quellcode

### Erfassungsablauf

1. Der Collector durchsucht die im Projekt angegebenen Go-Pakete
2. Er identifiziert Service-Typen (Structs mit Methoden, die dem Frontend bereitgestellt werden)
3. Für jeden Service erfasst er Informationen über dessen Methoden
4. Er identifiziert Modelltypen (Structs, die in Service-Methoden als Parameter oder Rückgabewerte verwendet werden)
5. Für jedes Modell erfasst er Informationen über dessen Felder und Typparameter
6. Er verarbeitet alle im Code gefundenen `//wails:`-Direktiven

## Renderingprozess

Der Renderingprozess generiert anhand der erfassten Informationen JavaScript-/TypeScript-Code. Dies übernimmt das Paket `render`.

### Hauptkomponenten

- **Renderer**: Koordiniert das Rendering von Service-, Modell- und Indexdateien
- **Module**: Repräsentiert ein einzelnes generiertes JavaScript-/TypeScript-Modul
- **Templates**: Für die Codegenerierung verwendete Textvorlagen

### Renderingablauf

1. Für jeden Service generiert der Renderer eine JavaScript-/TypeScript-Datei mit Funktionen, die den Service-Methoden entsprechen
2. Für jedes Modell generiert der Renderer eine JavaScript-/TypeScript-Klasse, die dem Modell-Struct entspricht
3. Der Renderer generiert Indexdateien, die alle Services und Modelle erneut exportieren
4. Der Renderer wendet alle durch `//wails:inject`-Direktiven angegebenen benutzerdefinierten Code-Injektionen an

## Typzuordnung

Einer der wichtigsten Aspekte des Binding-Systems ist die Zuordnung von Go-Typen zu JavaScript-/TypeScript-Typen. Die folgende Übersicht fasst die Zuordnungen zusammen:

| Go-Typ | JavaScript-Typ | TypeScript-Typ |
| --- | --- | --- |
| `bool` | `boolean` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64` | `number` | `number` |
| `string` | `string` | `string` |
| `[]byte` | `Uint8Array` | `Uint8Array` |
| `[]T` | `Array<T>` | `T[]` |
| `map[string]V` | `Object` | `{ [_: string]: V }` |
| `map[K]V` (nicht als Zeichenfolge vorliegendes `K`) | `Object` | `{ [_ in K]?: V }` |
| `struct` | `Object` | Benutzerdefinierte Klasse |
| `interface{}` | `any` | `any` |
| `*T` | `T \| null` | `T \| null` |
| `func` | Nicht unterstützt | Nicht unterstützt |
| `chan` | Nicht unterstützt | Nicht unterstützt |

## Direktivensystem

Das Bindungssystem unterstützt mehrere Direktiven, mit denen Sie den generierten Code anpassen können. Diese Direktiven werden Ihrem Go-Code als Kommentare hinzugefügt.

### Verfügbare Direktiven

- `//wails:inject`: Fügt benutzerdefinierten JavaScript-/TypeScript-Code in die generierten Bindungen ein
- `//wails:include`: Fügt den generierten Bindungen zusätzliche Dateien hinzu
- `//wails:internal`: Kennzeichnet einen Typ oder eine Methode als intern und verhindert so den Export an das Frontend
- `//wails:ignore`: Ignoriert eine Methode bei der Bindungsgenerierung vollständig
- `//wails:id`: Legt eine benutzerdefinierte ID für eine Methode fest und überschreibt damit die standardmäßige hashbasierte ID

### Verarbeitung von Direktiven

1. Während der Erfassungsphase identifiziert und analysiert der Collector die Direktiven im Go-Code
2. Die Direktiven werden zusammen mit den zugehörigen Deklarationen (Diensten, Methoden, Modellen usw.) gespeichert
3. Während der Renderingphase wendet der Renderer die Direktiven an, um den generierten Code anzupassen

## Erweiterte Funktionen

### Bedingte Codegenerierung

Das Bindungssystem unterstützt die bedingte Codegenerierung mit einem zweistelligen Bedingungspräfix für die Direktiven `include` und `inject`:

```
<language><style>:<content>
```

Dabei gilt:

- `<language>` kann folgende Werte haben:
  - `*` – JavaScript und TypeScript
  - `j` – nur JavaScript
  - `t` – nur TypeScript


- `<style>` kann folgende Werte haben:
  - `*` – Klassen und Schnittstellen
  - `c` – nur Klassen
  - `i` – nur Schnittstellen


Beispiel:

```go
//wails:inject j*:console.log("JavaScript only");
//wails:inject t*:console.log("TypeScript only");
```

### Benutzerdefinierte Methoden-IDs

Standardmäßig werden Methoden durch eine hashbasierte ID identifiziert. Sie können jedoch mit der Direktive `//wails:id` eine benutzerdefinierte ID festlegen:

```go
//wails:id 42
func (s *Service) CustomIDMethod() {}
```

Dies kann nützlich sein, um beim Refactoring von Code die Kompatibilität zu erhalten.

## Überlegungen zur Leistung

Der Bindungsgenerator ist auf Effizienz ausgelegt. Beachten Sie dennoch Folgendes:

1. Der erste Durchlauf ist langsamer, da zunächst ein Cache der zu durchsuchenden Pakete aufgebaut wird
2. Nachfolgende Durchläufe sind schneller, da sie die zwischengespeicherten Informationen verwenden
3. Der Generator verarbeitet alle Pakete im Projekt, was bei großen Projekten zeitaufwendig sein kann
4. Mit dem Flag `-clean` können Sie das Ausgabeverzeichnis vor der Generierung bereinigen

## Fehlersuche

Wenn bei der Bindungsgenerierung Probleme auftreten, können Sie mit dem Flag `-v` die Debug-Ausgabe aktivieren:

```bash
wails3 generate bindings -v
```

Dadurch erhalten Sie detaillierte Informationen zum Erfassungs- und Renderingprozess, mit denen sich die Ursache des Problems ermitteln lässt.
