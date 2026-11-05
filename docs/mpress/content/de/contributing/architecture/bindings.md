---
title: "Bindungssystem"
description: "Wie das Bindungssystem JavaScript-/TypeScript-Code erfasst, verarbeitet und generiert"
slug: "contributing/architecture/bindings"
sourcePath: "contributing/architecture/bindings.md"
---

Dieser Leitfaden erläutert die interne Funktionsweise des Wails-Bindungssystems und vermittelt Entwicklern Einblicke in die Mechanismen der automatischen Codegenerierung.

## Architekturübersicht

Das Wails-Bindungssystem besteht aus drei Hauptkomponenten:

1. **Erfassung**: Analysiert Go-Code, um Informationen über Services, Modelle und andere Deklarationen zu extrahieren
2. **Konfiguration**: Verwaltet Einstellungen und Optionen für die Bindungsgenerierung
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

- **Collector**: Verwaltet Paketinformationen und speichert erfasste Daten zwischen
- **Package**: Repräsentiert ein analysiertes Go-Paket und speichert die erfassten Services, Modelle und Direktiven
- **Service**: Erfasst Informationen über Service-Typen und deren Methoden
- **Model**: Erfasst detaillierte Informationen über Modelltypen, einschließlich Feldern, Werten und Typparametern
- **Directive**: Analysiert und interpretiert `//wails:`-Direktiven im Go-Quellcode

### Erfassungsablauf

1. Der Collector durchsucht die im Projekt angegebenen Go-Pakete
2. Er identifiziert Service-Typen (Structs mit Methoden, die dem Frontend zugänglich gemacht werden)
3. Für jeden Service erfasst er Informationen über dessen Methoden
4. Er identifiziert Modelltypen (Structs, die in Service-Methoden als Parameter oder Rückgabewerte verwendet werden)
5. Für jedes Modell erfasst er Informationen über dessen Felder und Typparameter
6. Er verarbeitet alle im Code gefundenen `//wails:`-Direktiven

## Renderingprozess

Der Renderingprozess generiert anhand der erfassten Informationen JavaScript-/TypeScript-Code. Dies übernimmt das Paket `render`.

### Hauptkomponenten

- **Renderer**: Steuert das Rendering von Service-, Modell- und Indexdateien
- **Module**: Repräsentiert ein einzelnes generiertes JavaScript-/TypeScript-Modul
- **Templates**: Für die Codegenerierung verwendete Textvorlagen

### Renderingablauf

1. Für jeden Service generiert der Renderer eine JavaScript-/TypeScript-Datei mit Funktionen, die den Service-Methoden entsprechen
2. Für jedes Modell generiert der Renderer eine JavaScript-/TypeScript-Klasse, die dem Modell-Struct entspricht
3. Der Renderer generiert Indexdateien, die alle Services und Modelle erneut exportieren
4. Der Renderer wendet alle durch `//wails:inject`-Direktiven angegebenen benutzerdefinierten Code-Injektionen an

## Typzuordnung

Einer der wichtigsten Aspekte des Bindungssystems ist die Zuordnung von Go-Typen zu JavaScript-/TypeScript-Typen. Die folgende Übersicht fasst diese Zuordnung zusammen:

| Go-Typ | JavaScript-Typ | TypeScript-Typ |
| --- | --- | --- |
| `bool` | `boolean` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64` | `number` | `number` |
| `string` | `string` | `string` |
| `[]byte` | `Uint8Array` | `Uint8Array` |
| `[]T` | `Array<T>` | `T[]` |
| `map[string]V` | `Object` | `{ [_: string]: V }` |
| `map[K]V` (`K` ohne Zeichenfolgen) | `Object` | `{ [_ in K]?: V }` |
| `struct` | `Object` | Benutzerdefinierte Klasse |
| `interface{}` | `any` | `any` |
| `*T` | `T \| null` | `T \| null` |
| `func` | Nicht unterstützt | Nicht unterstützt |
| `chan` | Nicht unterstützt | Nicht unterstützt |

## Direktivensystem

Das Bindungssystem unterstützt mehrere Direktiven, mit denen sich der generierte Code anpassen lässt. Fügen Sie diese Direktiven als Kommentare in Ihren Go-Code ein.

### Verfügbare Direktiven

- `//wails:inject`: Fügt benutzerdefinierten JavaScript-/TypeScript-Code in die generierten Bindungen ein
- `//wails:include`: Nimmt zusätzliche Dateien in die generierten Bindungen auf
- `//wails:internal`: Unterdrückt Exporte im generierten SDK; verhindert jedoch nicht den Laufzeitaufruf exportierter Go-Methoden
- `//wails:ignore`: Ignoriert eine Methode bei der Bindungsgenerierung vollständig
- `//wails:id`: Legt eine benutzerdefinierte ID für eine Methode fest und überschreibt damit die standardmäßige hashbasierte ID

Diese Direktiven steuern den generierten Code. Bei normalen Builds berücksichtigt die Laufzeitregistrierung weiterhin exportierte Dienstmethoden. Verwenden Sie für private Operationen nicht exportierte Methoden oder nicht gebundene Typen und erzwingen Sie die Autorisierung in offengelegten Methoden.

### Verarbeitung von Direktiven

1. Während der Erfassungsphase erkennt und analysiert der Collector Direktiven im Go-Code
2. Die Direktiven werden zusammen mit den zugehörigen Deklarationen (Dienste, Methoden, Modelle usw.) gespeichert
3. Während der Rendering-Phase wendet der Renderer die Direktiven an, um den generierten Code anzupassen

## Erweiterte Funktionen

### Bedingte Codegenerierung

Das Bindungssystem unterstützt die bedingte Codegenerierung mit einem zweistelligen Bedingungspräfix für die Direktiven `include` und `inject`:

```
<language><style>:<content>
```

Dabei gilt:

- `<language>` kann folgende Werte annehmen:
  - `*` – JavaScript und TypeScript
  - `j` – nur JavaScript
  - `t` – nur TypeScript


- `<style>` kann folgende Werte annehmen:
  - `*` – Klassen und Schnittstellen
  - `c` – nur Klassen
  - `i` – nur Schnittstellen


Beispiel:

```go
//wails:inject j*:console.log("JavaScript only");
//wails:inject t*:console.log("TypeScript only");
```

### Benutzerdefinierte Methoden-IDs

Standardmäßig werden Methoden durch eine Hash-basierte ID identifiziert. Sie können jedoch mit der Direktive `//wails:id` eine benutzerdefinierte ID angeben:

```go
//wails:id 42
func (s *Service) CustomIDMethod() {}
```

Dies kann hilfreich sein, um bei der Umstrukturierung von Code die Kompatibilität aufrechtzuerhalten.

## Leistungsaspekte

Der Bindungsgenerator ist auf Effizienz ausgelegt. Beachten Sie jedoch Folgendes:

1. Der erste Durchlauf ist langsamer, da dabei ein Cache der zu durchsuchenden Pakete aufgebaut wird
2. Nachfolgende Durchläufe sind schneller, da sie die zwischengespeicherten Informationen verwenden
3. Der Generator verarbeitet alle Pakete im Projekt, was bei großen Projekten zeitaufwendig sein kann
4. Mit dem Flag `-clean` können Sie das Ausgabeverzeichnis vor der Generierung leeren

## Fehlerbehebung

Wenn bei der Bindungsgenerierung Probleme auftreten, können Sie mit dem Flag `-v` die Debug-Ausgabe aktivieren:

```bash
wails3 generate bindings -v
```

Dadurch erhalten Sie detaillierte Informationen zum Erfassungs- und Rendering-Prozess, die helfen können, die Ursache des Problems zu ermitteln.
