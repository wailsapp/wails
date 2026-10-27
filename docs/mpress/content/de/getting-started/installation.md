---
title: "Installation"
description: "Wails installieren und die Entwicklungsumgebung einrichten"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## Unterstützte Plattformen

- Windows AMD64/ARM64
- macOS AMD64/ARM64. Die generierten Build-Tasks sind auf macOS 12.0 oder neuer ausgerichtet. Dies ist das Bereitstellungsziel des Standard-Builds und bedeutet nicht, dass jede benutzerdefinierte Build-Konfiguration ältere Systeme unterstützt.
- Ubuntu 24.04 AMD64/ARM64 (andere Linux-Distributionen funktionieren möglicherweise ebenfalls!)

## Abhängigkeiten

Wails hat mehrere allgemeine Abhängigkeiten, die vor der Installation erforderlich sind.

@note{type="tip"}
Nach der Installation der Wails CLI kannst du `wails3 setup` ausführen, um diese Abhängigkeiten automatisch zu prüfen und Hilfe bei ihrer Installation zu erhalten.

@end

@tabs
[Go (mindestens 1.26)]
Wails v3 wählt eine Go-Toolchain ab Version 1.26 aus. Ein älterer Go-Launcher kann sie mit `GOTOOLCHAIN=auto` automatisch herunterladen; andernfalls musst du sie ausdrücklich installieren. Die niedrigere Go-Versionsangabe des Startmoduls setzt die Wails-Abhängigkeit nicht außer Kraft.

Lade Go von der [Go-Downloadseite](https://go.dev/dl/) herunter.

Befolge unbedingt die offiziellen [Installationsanweisungen für Go](https://go.dev/doc/install). Stelle außerdem sicher, dass deine Umgebungsvariable `PATH` auch den Pfad zu deinem Verzeichnis `~/go/bin` enthält. Starte dein Terminal neu und führe die folgenden Prüfungen durch:

- Prüfe, ob Go korrekt installiert ist: `go version`
- Prüfe, ob `~/go/bin` in deiner PATH-Variable enthalten ist
  - Mac/Linux: `echo $PATH | grep go/bin`
  - Windows: `$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm (optional)]
Wails setzt npm zwar nicht voraus, die meisten mitgelieferten Vorlagen benötigen es jedoch.

Lade das neueste Node-Installationsprogramm von der [Node-Downloadseite](https://nodejs.org/en/download/) herunter. Verwende möglichst die neueste Version, da wir im Allgemeinen mit dieser testen.

Führe zur Überprüfung `npm --version` aus.

@note{type="info"}
Wenn du einen anderen Paketmanager als npm bevorzugst, kannst du ihn verwenden. Du musst die Taskfiles des Projekts entsprechend anpassen.

@end

@end

## Plattformspezifische Abhängigkeiten

Du musst außerdem plattformspezifische Abhängigkeiten installieren:

@tabs{sync-key="platform"}
[Mac]
Wails setzt die Installation der Xcode-Befehlszeilentools voraus. Führe dazu Folgendes aus:

```sh
xcode-select --install
```

[Windows]
Wails setzt die Installation der [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) voraus. Auf fast allen Windows-Installationen ist sie bereits vorhanden. Du kannst dies mit dem Befehl `wails doctor` prüfen.

[Linux]
Unter Linux sind die üblichen `gcc`-Build-Werkzeuge sowie `gtk4` und `webkitgtk-6.0` erforderlich. Führe nach der Installation <code>wails3 doctor</code> aus, um Anweisungen zur Installation der Abhängigkeiten zu erhalten. Der veraltete GTK3-/WebKit2GTK-Stack 4.1 ist bis v3.1 weiterhin über `-tags gtk3` verfügbar (siehe [Linux-Paketierung – Unterstützung für das veraltete GTK3](/guides/build/linux/#legacy-gtk3-support)). Wenn deine Distribution oder dein Paketmanager nicht unterstützt wird, teile uns dies bitte auf Discord mit.

@end

## Installation

Führe die folgenden Befehle aus, um die Wails CLI mithilfe von Go Modules zu installieren:

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

Führe die folgenden Befehle aus, wenn du die neueste Entwicklungsversion installieren möchtest:

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

Bei Verwendung der Entwicklungsversion nutzen alle generierten Projekte die Go-Direktive [replace](https://go.dev/ref/mod#go-mod-file-replace), damit sie die Entwicklungsversion von Wails verwenden.

## Nächste Schritte

Führe nach der Installation der CLI den Einrichtungsassistenten aus, um deine Entwicklungsumgebung zu konfigurieren:

```shell
wails3 setup
```

@note{type="caution" title="Experimentell"}
Der Einrichtungsassistent ist neu und wurde hauptsächlich unter Linux getestet. Wenn Probleme auftreten, [melde sie](https://github.com/wailsapp/wails/issues/4904) bitte und führe die folgenden Schritte zur manuellen Installation der Abhängigkeiten aus.

@end

Der Einrichtungsassistent führt folgende Aufgaben aus:

- Plattformabhängigkeiten prüfen und bei ihrer Installation helfen
- Projektstandards konfigurieren (Autorenangaben, Präfix der Bundle-ID)
- Optional Docker für plattformübergreifende Builds einrichten
- Codesignierung konfigurieren (falls erforderlich)

Weitere Einzelheiten findest du im [Einrichtungsleitfaden](/getting-started/setup/).

## Manuelle Installation der Abhängigkeiten

Wenn du die Abhängigkeiten lieber manuell installieren möchtest oder der Einrichtungsassistent auf deinem System nicht funktioniert, befolge die obigen plattformspezifischen Anweisungen und führe anschließend Folgendes aus:

```shell
wails3 doctor
```

Dadurch wird geprüft, ob die richtigen Abhängigkeiten installiert sind, und du erhältst Hinweise zu fehlenden Komponenten.

## Der Befehl `wails3` scheint zu fehlen?

Wenn dein System meldet, dass der Befehl `wails3` fehlt, prüfe Folgendes:

- Stelle sicher, dass du die obige **Go-Installationsanleitung** korrekt befolgt hast und das Verzeichnis `go/bin` in der Umgebungsvariable `PATH` enthalten ist.
- Schließe die aktuellen Terminals und öffne sie erneut, damit die neue Variable `PATH` übernommen wird.
