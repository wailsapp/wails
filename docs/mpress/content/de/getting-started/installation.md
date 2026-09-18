---
title: "Installation"
description: "Wails installieren und die Entwicklungsumgebung einrichten"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## Unterstützte Plattformen

- Windows AMD64/ARM64
- macOS 10.15+ AMD64 (Bereitstellung auf macOS 10.13+ möglich)
- macOS 11.0+ ARM64
- Ubuntu 24.04 AMD64/ARM64 (andere Linux-Distributionen funktionieren möglicherweise ebenfalls!)

## Abhängigkeiten

Wails hat mehrere allgemeine Abhängigkeiten, die vor der Installation vorhanden sein müssen.

@note{type="tip"}
Nach der Installation der Wails-CLI können Sie `wails3 setup` ausführen, um diese Abhängigkeiten automatisch zu prüfen und Unterstützung bei ihrer Installation zu erhalten.

@end

@tabs
[Go (mindestens 1.24)]
Laden Sie Go von der [Go-Downloadseite](https://go.dev/dl/) herunter.

Befolgen Sie unbedingt die offiziellen [Installationsanweisungen für Go](https://go.dev/doc/install). Stellen Sie außerdem sicher, dass die Umgebungsvariable `PATH` auch den Pfad zu Ihrem Verzeichnis `~/go/bin` enthält. Starten Sie Ihr Terminal neu und führen Sie die folgenden Prüfungen durch:

- Prüfen Sie, ob Go ordnungsgemäß installiert ist: `go version`
- Prüfen Sie, ob `~/go/bin` in Ihrer PATH-Variable enthalten ist
  - Mac/Linux: `echo $PATH | grep go/bin`
  - Windows: `$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm (optional)]
Wails setzt keine Installation von npm voraus, die meisten mitgelieferten Vorlagen benötigen es jedoch.

Laden Sie das neueste Node-Installationsprogramm von der [Node-Downloadseite](https://nodejs.org/en/download/) herunter. Verwenden Sie möglichst die neueste Version, da wir unsere Tests in der Regel damit durchführen.

Führen Sie zur Überprüfung `npm --version` aus.

@note{type="info"}
Wenn Sie einen anderen Paketmanager als npm bevorzugen, können Sie diesen verwenden. Sie müssen die Taskfiles des Projekts entsprechend anpassen.

@end

@end

## Plattformspezifische Abhängigkeiten

Sie müssen außerdem plattformspezifische Abhängigkeiten installieren:

@tabs{sync-key="platform"}
[Mac]
Für Wails müssen die Xcode-Befehlszeilentools installiert sein. Führen Sie dazu Folgendes aus:

```sh
xcode-select --install
```

[Windows]
Für Wails muss die [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) installiert sein. Bei fast allen Windows-Installationen ist sie bereits vorhanden. Sie können dies mit dem Befehl `wails doctor` prüfen.

[Linux]
Unter Linux sind die standardmäßigen `gcc`-Build-Werkzeuge sowie `gtk4` und `webkitgtk-6.0` erforderlich. Führen Sie nach der Installation <code>wails3 doctor</code> aus, um Anweisungen zur Installation der Abhängigkeiten zu erhalten. Der veraltete Stack aus GTK3 und WebKit2GTK 4.1 ist bis v3.1 weiterhin über `-tags gtk3` verfügbar (siehe [Linux-Paketierung – Unterstützung für das veraltete GTK3](/guides/build/linux/#legacy-gtk3-support)). Falls Ihre Distribution oder Ihr Paketmanager nicht unterstützt wird, teilen Sie uns dies bitte auf Discord mit.

@end

## Installation

Führen Sie die folgenden Befehle aus, um die Wails-CLI mit Go Modules zu installieren:

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Führen Sie die folgenden Befehle aus, wenn Sie die neueste Entwicklungsversion installieren möchten:

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

Bei Verwendung der Entwicklungsversion nutzen alle generierten Projekte die Go-Direktive [replace](https://go.dev/ref/mod#go-mod-file-replace), damit sie die Entwicklungsversion von Wails verwenden.

## Nächste Schritte

Führen Sie nach der Installation der CLI den Einrichtungsassistenten aus, um Ihre Entwicklungsumgebung zu konfigurieren:

```shell
wails3 setup
```

@note{type="caution" title="Experimentell"}
Der Einrichtungsassistent ist neu und wurde hauptsächlich unter Linux getestet. Falls Probleme auftreten, [melden Sie diese bitte](https://github.com/wailsapp/wails/issues/4904) und führen Sie die unten beschriebenen Schritte zur manuellen Installation der Abhängigkeiten aus.

@end

Der Einrichtungsassistent führt folgende Aufgaben aus:

- Plattformabhängigkeiten prüfen und bei ihrer Installation unterstützen
- Projektstandardwerte konfigurieren (Autoreninformationen, Präfix der Bundle-ID)
- Optional Docker für plattformübergreifende Builds einrichten
- Codesignierung konfigurieren (falls erforderlich)

Weitere Informationen finden Sie im [Einrichtungsleitfaden](/getting-started/setup/).

## Manuelle Installation der Abhängigkeiten

Wenn Sie die Abhängigkeiten lieber manuell installieren möchten oder der Einrichtungsassistent auf Ihrem System nicht funktioniert, befolgen Sie die obigen plattformspezifischen Anweisungen und führen Sie anschließend Folgendes aus:

```shell
wails3 doctor
```

Dadurch wird geprüft, ob die richtigen Abhängigkeiten installiert sind, und Sie erhalten Hinweise zu fehlenden Abhängigkeiten.

## Der Befehl `wails3` scheint zu fehlen?

Falls Ihr System meldet, dass der Befehl `wails3` fehlt, prüfen Sie Folgendes:

- Stellen Sie sicher, dass Sie die obige **Go-Installationsanleitung** korrekt befolgt haben und dass das Verzeichnis `go/bin` in der Umgebungsvariable `PATH` enthalten ist.
- Schließen und öffnen Sie die aktuellen Terminals erneut, damit die neue Variable `PATH` übernommen wird.
