---
title: "Wake"
description: "Ein experimenteller, auf Wails abgestimmter Build-Runner, der Ihre vorhandenen Taskfiles mit schnelleren inkrementellen Builds, strukturierter Ausgabe und standardmäßig paralleler Ausführung ausführt."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

@note{type="caution" title="Experimentelle Funktion"}
Wake muss über `WAILS_USE_WAKE=true` explizit aktiviert werden und ist **nicht** der standardmäßige Runner. Ist die Variable nicht gesetzt, verhält sich `wails3 build / package / sign / task` genau wie bisher. Funktionsumfang und Verhalten können sich zwischen Releases ändern.

@end

Wake ist ein **experimenteller alternativer Build-Runner** für `wails3`. Er liest dasselbe `Taskfile.yml`, das Ihr Projekt bereits enthält – mit derselben Syntax für Tasks, Abhängigkeiten, Variablen, Vorlagen, Includes und Plattform-Namespaces – und führt es mit einem auf Wails abgestimmten Executor statt mit der universell einsetzbaren [Task](https://taskfile.dev)- Runtime aus.

Das Ziel besteht nicht darin, Task zu ersetzen. Wake soll einen Runner bereitstellen, der speziell auf die tatsächlichen Build-Abläufe von Wails-Projekten zugeschnitten ist und dessen Semantik, Ausgabe und Standardwerte zum übrigen `wails3`-CLI passen. **Wenn Sie ausschließlich Wake verwenden, ändern sich Ihre Taskfiles nicht.**

## Warum es Wake gibt

Sowohl Wake als auch die Task-Runtime sind in `wails3` einkompiliert – keines von beiden erfordert die Installation einer separaten Binärdatei. Der Unterschied besteht darin, dass Wake **die Domäne versteht**. Ein universell einsetzbarer Runner führt die in einem Taskfile aufgeführten Schritte in der vorgegebenen Reihenfolge aus. Wake weiß, was ein Wails-Build tatsächlich *ist*: Das Frontend-Bundle wird in die Binärdatei eingebettet, die Binärdatei wird in plattformspezifische Artefakte verpackt und parallel dazu werden Symbole und Bindings erzeugt. Dieses Wissen nutzt Wake, um den Build auf Arten zu optimieren, die einem generischen Runner nicht möglich sind.

- **Wake führt nur die für einen Build tatsächlich erforderlichen Arbeiten aus.** Wake verfolgt selbst die tatsächlichen Ein- und Ausgaben jedes Schritts. Bei einem Go-Build sind das der Modulgraph und die Ausgaben der Schritte, von denen er abhängt. Wenn sich nichts Relevantes geändert hat, überspringt Wake daher Compiler und Linker vollständig, statt sie erneut auszuführen. Ein universell einsetzbarer Runner kann einen Schritt nur überspringen, wenn im Taskfile vorab genau angegeben wurde, welche Dateien überwacht werden sollen. Wake leitet dies aus seinem vorhandenen Wissen über den Build ab. Bei einem erneuten Build ohne Änderungen sind das ungefähr **~20 ms (Wake) gegenüber ~316 ms (Task)**. Kalte Builds benötigen dieselbe Gesamtzeit – sie werden von `npm install`, Vite und dem Go-Compiler dominiert.

- **Wake weiß, was gleichzeitig ausgeführt werden kann.** Da Wake erkennt, welche Schritte unabhängig voneinander sind, führt es sie standardmäßig parallel aus. Die Ergebniszeile weist den dadurch erzielten Geschwindigkeitsgewinn aus. Deaktivieren Sie dies mit `WAKE_SERIAL=true`, wenn die verschachtelte Ausgabe nebenläufiger Schritte eine Untersuchung erschweren würde.

- **Von wails3 gesteuerte strukturierte Ausgabe.** Wake verwendet den eigenen Reporter von wails3: eine Zeile pro geplantem Schritt, Live-Status, eine farbcodierte Aufschlüsselung der Phasen am Ende und anklickbare `file:line`-Links in Fehlerbereichen. Bei `NO_COLOR` und in Umgebungen ohne TTY (CI-Protokollen) wird die Darstellung kontrolliert vereinfacht.

- **Integriert, damit Wake gemeinsam mit Wails wachsen kann.** Da Wake Teil von `wails3` und kein Drittanbieterwerkzeug ist, können neue Build-Funktionen direkt ergänzt werden – es muss nicht darauf gewartet werden, dass ein separates Projekt sie implementiert. Dadurch eröffnet sich außerdem die Möglichkeit, plattformübergreifende Skripte und Werkzeuge nativ auszuführen, wo ein Taskfile derzeit die `wails3`-Binärdatei über die Shell aufruft und für jeden Aufruf einen Prozess startet. Werden diese Arbeiten in den Prozess integriert, sinkt der Overhead und weitere Geschwindigkeitssteigerungen werden möglich.

## Wake aktivieren

Wake wird vollständig über die Umgebungsvariable `WAILS_USE_WAKE=true` gesteuert. Ist sie nicht gesetzt oder auf einen anderen Wert als `true` gesetzt, verwendet jeder `wails3`-Befehl wie bisher genau die eingebettete Task-Runtime.

```bash
# Default: Task runtime, no Wake involvement
wails3 build

# Opt in: Wake drives build / package / sign / task <name>
WAILS_USE_WAKE=true wails3 build
WAILS_USE_WAKE=true wails3 package
WAILS_USE_WAKE=true wails3 task <some-task-name>
```

Das Flag gilt für `wails3 build`, `wails3 package`, `wails3 sign` und `wails3 task <name>`. `wails3 dev` ist davon derzeit **nicht** betroffen – der Entwicklungs-Watcher verwendet weiterhin seine eigene Pipeline.

@note{type="tip" title="Wake kann gefahrlos aktiviert werden"}
Trifft Wake auf eine nicht implementierte Taskfile-Funktion, übergibt es den gesamten Durchlauf innerhalb desselben Prozesses an die eingebettete Task-Runtime – es muss keine externe `task`- Binärdatei installiert werden. Im ungünstigsten Fall erhalten Sie genau das Verhalten, das Sie auch ohne das Flag erhalten hätten.

@end

## Gestaffelte lokale Überschreibungen

Wake unterstützt ein **Basis-Taskfile mit lokalen Überschreibungen**. Legen Sie neben Ihrem `Taskfile.yml` eine Datei ab; deren Definitionen haben Vorrang:

| Datei | Zweck | Priorität |
| --- | --- | --- |
| `Taskfile.yml` | Basis, eingecheckt | niedrigste |
| `Taskfile.override.yml` / `.yaml` | eingecheckte, teamweite Überschreibungen | mittlere |
| `Taskfile.local.yml` / `.yaml` | persönlich, normalerweise von Git ignoriert | höchste |

**Zusammenführungssemantik (lokale Definition hat Vorrang):**

- Ein Task mit **demselben Namen** überschreibt den Basis-Task. Listenfelder (`cmds`, `deps`, `sources`, `generates`, `platforms`, `status`, `preconditions`, `aliases`) **ersetzen** die entsprechenden Felder der Basis, wenn die Überschreibung sie angibt. Felder, die in der Überschreibung fehlen, werden aus der Basis übernommen.
- `env` und `vars` werden **schlüsselweise zusammengeführt**; bei Konflikten hat die Überschreibung Vorrang.
- Ein Task, der **nur** in einer Überschreibungsdatei vorhanden ist, wird **hinzugefügt**.

Wenn beispielsweise das eingecheckte `Taskfile.yml` mit Entwicklungs-Flags baut, auf Ihrem Rechner aber immer ein Produktions-Build erstellt werden soll:

```yaml
# Taskfile.local.yml (git-ignored, yours)
tasks:
  build:
    cmds:
      - go build -tags production -o bin/app .
  smoke:
    cmds:
      - ./bin/app --selftest
```

Nun führt `build` Ihren Produktionsbefehl aus und `smoke` ist verfügbar, ohne dass das eingecheckte Taskfile geändert wird.

@note{type="note" title="Vertrauensmodell"}
Überschreibungsdateien werden automatisch erkannt und ohne Rückfrage angewendet. Dadurch entstehen keine neuen Möglichkeiten: Ein Taskfile kann bereits beliebige Shell-Befehle ausführen. Eine Überschreibung kann daher nichts bewirken, was nicht auch durch Bearbeiten von `Taskfile.yml` möglich wäre. Das eingecheckte `Taskfile.override.*` erscheint in PR-Diffs; `Taskfile.local.*` wird auf Ihrem eigenen Rechner erstellt. Eine fehlerhafte Überschreibung bricht den Durchlauf ab, statt kommentarlos übersprungen zu werden. Setzen Sie `WAILS_NO_OVERRIDES=true`, um die Erkennung von Überschreibungen für deterministische CI-Builds vollständig zu deaktivieren.

@end

## Automatischer Rückgriff

Trifft Wake auf eine nicht implementierte Taskfile-Funktion, übergibt es den gesamten Durchlauf an die eingebettete Task-Runtime. Folgende Funktionen lösen diesen Rückgriff derzeit aus:

- `dotenv` auf Taskfile-Ebene
- andere `output`-Modi als `interleaved`
- ein `requires`-Block
- `interval` (auf Taskfile- oder Task-Ebene)
- andere `run`-Modi als `always`
- `short` in einem Task
- `defer` in einer Task

## Umgebungsvariablen

| Variable | Wirkung |
| --- | --- |
| `WAILS_USE_WAKE` | `true` aktiviert Wake für die routingfähigen `wails3`-Verben; bei jedem anderen Wert wird die Task-Laufzeit verwendet |
| `WAILS_NO_OVERRIDES` | `true` überspringt die Erkennung von `Taskfile.local.*` / `.override.*` (deterministische Builds) |
| `WAKE_VERBOSE` | stdout/stderr von Unterprozessen live streamen, statt sie nur zur Anzeige bei Fehlern zu erfassen |
| `WAKE_SILENT` | Task-Ausgabe vollständig unterdrücken |
| `WAKE_SERIAL` | `true` deaktiviert den parallelen `deps:`-Fan-out (parallel ist die Standardeinstellung) |
| `WAKE_FORCE` | `true` umgeht für einen vollständig sauberen Neu-Build alle Caches |
| `WAKE_DEBUG` | Interna des Resolvers protokollieren (DAG, Abhängigkeiten, Variablenreferenzen, Ausführungsrouting) |
| `WAKE_NOTICE` | `off`, um den Hinweis „wake (experimental)“ bei jeder Ausführung zu unterdrücken |

Der Build-Cache befindet sich in `.wake/cache.json` (Task verwendet `.task/`).

## Feedback

Wake ist ein Experiment, und Ihr Feedback entscheidet über seine weitere Entwicklung. Wenn Sie es ausprobieren, möchten wir gern erfahren, ob es schneller und verständlicher war und ob etwas nicht mehr funktionierte — die hilfreichsten Berichte nennen, was Sie ausgeführt haben, was Sie erwartet haben und was tatsächlich passiert ist. Schreiben Sie uns in der [Diskussion zum Wake-Feedback](https://github.com/wailsapp/wails/discussions/5679).
