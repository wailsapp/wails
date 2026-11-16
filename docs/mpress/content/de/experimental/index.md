---
title: "Experimentelle Funktionen"
description: "Ein Ort für laufende Experimente in Wails v3 – was sie sind, warum es sie gibt und wie Sie Feedback geben können."
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="Hier gibt es Experimente"}
Alles in diesem Abschnitt ist per Definition ein Experiment. Diese Funktionen müssen ausdrücklich aktiviert werden, sind standardmäßig deaktiviert und können zwischen Releases ihre Form oder ihren Namen ändern oder vollständig entfernt werden. Verwenden Sie sie nur für tragende Komponenten, wenn Sie auf häufige Änderungen vorbereitet sind.

@end

## Was „experimentell“ bedeutet

Wails entwickelt Experimente öffentlich. Ein Experiment ist eine Idee, die wir für vielversprechend genug halten, um sie Ihnen frühzeitig zur Verfügung zu stellen – zu der wir uns aber noch nicht vollständig bekannt haben. Wir veröffentlichen sie, *weil* wir aus dem praktischen Einsatz lernen möchten, bevor wir entscheiden, ob sie ein dauerhafter, unterstützter Bestandteil von Wails wird.

Daraus ergeben sich für alles in diesem Abschnitt einige Punkte:

- **Die Aktivierung ist optional.** Experimente ändern niemals das Standardverhalten von `wails3`. Sie aktivieren sie bewusst – üblicherweise mit einer Umgebungsvariablen oder einem Build-Flag. Bleiben sie deaktiviert, ändert sich nichts an Ihrem bestehenden Arbeitsablauf.
- **Das Experiment bleibt möglicherweise nicht erhalten.** Einige Experimente werden zu stabilen Funktionen. Andere werden bis zur Unkenntlichkeit überarbeitet oder eingestellt. Wir probieren Dinge lieber öffentlich aus und lernen schnell daraus, statt nur das zu veröffentlichen, dessen wir uns bereits sicher sind.
- **Die API ist nicht eingefroren.** Während ein Experiment seine endgültige Form findet, können sich Namen, Flags, Standardwerte und Verhalten zwischen Releases ändern. Die Release Notes weisen auf Änderungen hin, aber erwarten Sie nicht die Stabilitätsgarantien stabiler Funktionen.

## Wir möchten Ihr Feedback

Dies ist der entscheidende Punkt. Ob Experimente fortgeführt oder eingestellt werden, hängt von den Rückmeldungen der Personen ab, die sie tatsächlich verwenden. Wenn Sie eines dieser Experimente ausprobieren, möchten wir wirklich Folgendes wissen:

- Hat es für Ihr Projekt funktioniert? Wo ist es gescheitert?
- War es schneller, verständlicher oder angenehmer – oder den Wechsel nicht wert?
- Welche Voraussetzungen müssten erfüllt sein, damit Sie es standardmäßig verwenden?

Am hilfreichsten ist konkretes Feedback: was Sie ausgeführt haben, was Sie erwartet haben und was tatsächlich passiert ist. Für jedes Experiment gibt es in den GitHub Discussions einen eigenen Thread in der Kategorie **Experimente**:

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Diskussionen zu Experimenten" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="Suchen Sie den Thread zu Ihrem Experiment und teilen Sie uns mit, wie es funktioniert hat, was nicht funktioniert oder was fehlt."}
@end

## Aktuelle Experimente

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="Probieren Sie das native YAML-Build-System aus: Erstellen und konfigurieren Sie ein Projekt in einem einzigen Assistenten, bauen und paketieren Sie es oder prüfen Sie eine unidirektionale Migration von Taskfiles."}
@linkcard{title="CGO-freie Builds" href="/experimental/#cgo-free-builds-purego" description="Bauen Sie macOS- und Linux-Apps mit CGO_ENABLED=0: keine C-Toolchain, Cross-Kompilierung ohne Docker. Aktivierung mit WAILS_PUREGO=1 oder dem Build-Tag purego."}
@linkcard{title="LLM-Steuerung (MCP)" href="/guides/mcp-service/" description="Ein integrierter Model-Context-Protocol-Server, mit dem LLM-Agenten eine laufende Wails-App untersuchen, testen und steuern können – ohne erforderlichen Benutzercode und aktiviert über ein Build-Tag."}
@end

## CGO-freie Builds (purego)

@note{type="caution" title="Experimentelle Funktion"}
Das CGO-freie Backend muss über das Build-Tag `purego` (oder `WAILS_PUREGO=1`) ausdrücklich aktiviert werden und ist **nicht** die Standardeinstellung. Ohne diese Aktivierung verwenden macOS- und Linux-Builds wie bisher das cgo-Backend. Das Verhalten kann sich zwischen Releases ändern.

@end

### Funktionsweise

Unter macOS und Linux kommuniziert Wails normalerweise über cgo mit den UI-Bibliotheken des Systems: AppKit und WKWebView unter macOS sowie GTK und WebKitGTK unter Linux. Das CGO-freie Backend führt dieselben Aufrufe aus reinem Go mithilfe von [purego](https://github.com/ebitengine/purego) aus, das die Systembibliotheken zur Laufzeit lädt.

Die öffentliche Anwendungs-API bleibt unverändert. Wählen Sie das Backend beim Bauen aus; der Anwendungscode muss nicht geändert werden. Testen Sie Ihre App vor der Verteilung auf dem Zielsystem.

Ihre Vorteile:

- **Keine C-Toolchain für die Go-Kompilierung.** Der Go-Build benötigt weder die Xcode-Kommandozeilenwerkzeuge noch `gcc` oder Linux-Entwicklungsheader. Frontend-Builds benötigen weiterhin ihre üblichen JavaScript-Werkzeuge.
- **Cross-Kompilierung ohne Docker.** Bauen Sie mit einfachem `go build` auf jedem Betriebssystem eine Binärdatei für macOS oder Linux.

Windows wird bereits ohne cgo gebaut und ist nicht betroffen.

### Vor der Aktivierung

Verwenden Sie eine CLI- und Framework-Version, die beide dieses Experiment enthalten. Wenn Sie nur die CLI aktualisieren, wird das Backend nicht zu einer älteren Framework-Abhängigkeit in der `go.mod` Ihrer Anwendung hinzugefügt. Führen Sie zum Testen der Entwicklungsversion diese Befehle im Verzeichnis Ihrer Anwendung aus:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@master
go get github.com/wailsapp/wails/v3@master
```

Beide Befehle folgen dem sich laufend ändernden Branch `master`. Verwenden Sie für reproduzierbare Builds denselben getesteten Commit für CLI und Framework. Befolgen Sie zur Einrichtung der Go-Toolchain und der CLI die [Installationsanleitung](/getting-started/installation/).

Die experimentellen Backends für macOS und Linux unterstützen `amd64` und `arm64`. Bei der Cross-Kompilierung wird eine ausführbare Datei für das angeforderte Ziel erstellt; sie wird nicht auf dem Build-Rechner ausgeführt. Testen Sie sie auf dem Zielbetriebssystem und der Zielarchitektur.

### Verwendung

Wählen Sie eine der folgenden Methoden aus, um das Backend festzulegen.

#### Mit der Wails CLI

In einem mit dieser CLI-Version erzeugten Projekt verwenden `build`, `package` und `dev` das experimentelle Backend, wenn `WAILS_PUREGO` aktiviert ist. In einer POSIX-Shell:

```sh
WAILS_PUREGO=1 wails3 dev
WAILS_PUREGO=1 wails3 build
WAILS_PUREGO=1 wails3 package
```

In PowerShell:

```powershell
$env:WAILS_PUREGO = "1"
wails3 build
```

Die CLI erkennt `1`, `true`, `on` und `yes` unabhängig von Groß- und Kleinschreibung sowie umgebenden Leerzeichen. Ein leerer Wert aktiviert dieses Experiment nicht.

Alternativ können Sie das Tag an einen einzelnen Build- oder Paketierungsbefehl übergeben:

```sh
wails3 build -tags purego
wails3 package -tags purego
```

Mit gesetztem Tag bauen die aktualisierten Standard-Taskfiles mit `CGO_ENABLED=0` und verwenden dabei die Go-Toolchain des Hostsystems. So bauen Sie beispielsweise unter macOS oder Windows eine ausführbare Linux-ARM64-Datei:

```sh
wails3 build -tags purego GOOS=linux ARCH=arm64
```

Ältere Projekte behalten ihre vorhandenen Taskfiles. Vergleichen Sie benutzerdefinierte Build-Tasks mit den aktualisierten Vorlagen oder testen Sie das Backend mit einfachem Go, ohne diese Tasks zu ändern. Bei Befehlen der Wails CLI wählt die Umgebungsvariable das Tag aus; bei einfachem `go build` muss das Tag ausdrücklich angegeben werden.

#### Mit einfachem Go

Bauen Sie zuerst Ihr Frontend, damit alle von `go:embed` referenzierten Assets vorhanden sind, und kompilieren Sie anschließend die Anwendung:

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

`CGO_ENABLED=1 go build -tags purego` funktioniert ebenfalls: Das Backend bleibt dasselbe; cgo steht lediglich weiterhin für den Rest Ihres Programms zur Verfügung.

Legen Sie für einen Cross-Build zusätzlich `GOOS` und `GOARCH` fest:

```sh
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags purego,production -o myapp .
```

Wählen Sie unter Linux den GTK3-Stack mit `-tags purego,gtk3,production` aus.

Um zum Standard-Backend zurückzukehren, entfernen Sie das Tag `purego` und heben Sie die Definition von `WAILS_PUREGO` auf (mit `unset WAILS_PUREGO` in einer POSIX-Shell oder `Remove-Item Env:WAILS_PUREGO` in PowerShell). Entfernen Sie außerdem eine explizite Einstellung für `CGO_ENABLED=0`, wenn Sie zum cgo-Backend für macOS oder Linux zurückkehren.

Wenn Sie mit `CGO_ENABLED=0` bauen, aber das Tag vergessen, wird der Build mit einer einzelnen Fehlermeldung abgebrochen, die die erforderliche Korrektur nennt:

```
undefined: requires_CGO_ENABLED_1_or_build_tag_purego
```

### Laufzeitanforderungen

Die App lädt die Systembibliotheken beim Start. Daher müssen sie wie bei einem cgo-Build auf dem Zielrechner installiert sein:

| Plattform | Bibliotheken | Älteste unterstützte Version |
|---|---|---|
| macOS | AppKit, WebKit (Bestandteil des Betriebssystems) | Wie beim cgo-Backend |
| Linux (Standard) | GTK 4, WebKitGTK 6.0 | GTK 4.8, WebKitGTK 2.40 (Debian 12) |
| Linux (`gtk3`) | GTK 3, WebKit2GTK 4.1 | GTK 3.24, WebKit2GTK 2.36 (Ubuntu 22.04) |

Funktionen, die neuer als diese Mindestversionen sind, werden beim Start der App gesucht und nur verwendet, wenn die installierte Bibliothek sie bereitstellt. Dadurch läuft dasselbe Binärprogramm auf alten und neuen Systemen.

Die App verwendet weiterhin native UI-Bibliotheken und kann Systembibliotheken wie die libc unter Linux dynamisch einbinden. `CGO_ENABLED=0` bedeutet nicht, dass die ausführbare Datei vollständig statisch ist. Fehlt eine erforderliche Laufzeitbibliothek, installiere den passenden GTK/WebKit-Stack auf dem Zielsystem. Die Installation von C-Headern auf dem Build-System behebt keine fehlende Laufzeitbibliothek.

### Paketierung und Signierung

Das Backend beseitigt die Abhängigkeit von einer C-Toolchain beim Kompilieren. Jedes Paketformat behält seine eigenen Anforderungen. Insbesondere macht die Auswahl von `purego` die vorhandenen AppImage-Bereitstellungstools nicht plattformübergreifend nutzbar. Siehe [Linux-Paketierung](/guides/build/linux/).

Für andere Plattformen erstellte macOS-Anwendungen müssen vor der Verteilung weiterhin ordnungsgemäß signiert und notarisiert werden. Windows-Installationsprogramme benötigen weiterhin ihre Installationswerkzeuge. Siehe [macOS-Paketierung](/guides/build/macos/) und [Windows-Paketierung](/guides/build/windows/). Bei Abhängigkeiten, die cgo verwenden, erfordert `CGO_ENABLED=1` außerdem einen C-Compiler sowie sämtliche Header oder Bibliotheken für das Zielsystem. Die Auswahl des Pure-Go-Backends von Wails beseitigt die Anforderungen dieser Abhängigkeiten nicht.

### Validierung

Prüfungen der generierten Bindings und `puregocheck` untersuchen native Symbole, Schutzmechanismen für optionale Versionen und die Backend-Abdeckung. Der Linux-Paritätsworkflow testet GTK4 und GTK3 mit beiden Backends, einschließlich Asset-Anfragen, binären Anfragekörpern und nativer Weiterleitung; außerdem kompiliert er Linux-ARM64-Ziele für andere Plattformen. Cross-Compilation-Prüfungen stellen fest, ob sich der Code erstellen lässt, während Laufzeitprüfungen die auf dem Testsystem installierten nativen Bibliotheken testen.

Der [Backend-Testleitfaden](https://github.com/wailsapp/wails/blob/master/v3/TESTING.md) beschreibt die Paritätsprüfungen und End-to-End-Eingabetests. Unterscheide beim Melden eines Problems zwischen einem Kompilierungsfehler und einem Laufzeit- oder Paketierungsfehler.

### Bekannte Unterschiede

- **Menü-Tastenkürzel unter Linux werden auf Anwendungsebene registriert.** Weise verschiedenen Menüs keine kollidierenden Tastenkürzel zu. Wird das Menü eines Fensters mit `SetMenu(nil)` ausdrücklich entfernt, können die Tastenkürzel des entfernten Menüs in diesem Fenster ebenfalls nicht mehr ausgelöst werden; andere Fenster behalten ihre Menüs und Tastenkürzel.
- **Dialogfelder zum Öffnen und Speichern unter macOS** werden wie bei cgo in einem separaten Systemprozess ausgeführt; die Tests prüfen, ob sie sich öffnen und ordnungsgemäß abbrechen lassen.

Etwas anderes gefunden? Bitte teile es uns mit.

### Feedback

Berichte uns im Thread **CGO-freie Builds** unter [Experimente](https://github.com/wailsapp/wails/discussions/categories/experiments), wie es gelaufen ist: welches Betriebssystem und welche Version du verwendest, wie du den Build erstellt hast und was funktioniert hat oder nicht. Ein reproduzierbarer Fehlerfall mit der Ausgabe von `wails3 doctor` ist für uns der hilfreichste Bericht.
