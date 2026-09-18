---
title: "Von einer v3-Alpha aktualisieren"
description: "Ein bestehendes Wails-v3-Alpha-Projekt auf eine festgelegte Beta-Version aktualisieren"
slug: "migration/alpha-to-beta"
sourcePath: "migration/alpha-to-beta.md"
---

Diese Anleitung gilt für bestehende v3-Alpha-Projekte. Für Wails v2 verwenden Sie die [Anleitung von v2 zu v3](/migration/v2-to-v3/).

## Vor dem Upgrade

Committen oder sichern Sie Ihr Projekt. Lesen Sie das [Änderungsprotokoll](/changelog/) zwischen Ihrer Alpha-Version und der gewählten Beta: Quellcode, APIs oder Build-Konfiguration müssen möglicherweise angepasst werden. Prüfen Sie die [Desktop-Kompatibilitätsrichtlinie](/status/) und die Anforderungen Ihrer Plattform.

Die folgenden Befehle verwenden die veröffentlichte Version `v3.0.0-beta.23` als Beispiel für eine genaue Versionsangabe, nicht als Empfehlung, stets die neueste Version zu verwenden. Prüfen Sie bei einer anderen Version die zugehörigen CLI-, Go-Modul- und npm-Runtime-Versionen und passen Sie die Befehle gemeinsam an. In diesem Beispiel entspricht die npm-Version der Go-Version ohne das Präfix `v`.

## 1. CLI aktualisieren

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.23
wails3 version
```

Prüfen Sie, ob `wails3 version` die installierte Version meldet. Eine ältere Binärdatei weiter vorne im `PATH` kann die neue CLI verdecken.

## 2. Go-Modul aktualisieren

Führen Sie die Befehle im Projektstamm aus. Prüfen Sie die Änderungen an den Abhängigkeiten; aktualisieren Sie nicht pauschal unbeteiligte Module.

```sh
go get github.com/wailsapp/wails/v3@v3.0.0-beta.23
go mod tidy
```

## 3. Frontend-Runtime aktualisieren

Für Projekte mit npm und einem Verzeichnis `frontend`:

```sh
cd frontend
npm install --save-exact @wailsio/runtime@3.0.0-beta.23
cd ..
```

Behalten Sie die Sperrdatei und prüfen Sie ihre Änderungen. Passen Sie diesen Schritt bei einem anderen Paketmanager oder Verzeichnis an und behalten Sie eine genaue Runtime-Version bei.

## 4. Neu generieren, bauen und testen

Generieren Sie im Projektstamm die Bindings aus Ihren Go-Services neu und bauen Sie die Anwendung:

```sh
wails3 generate bindings
wails3 build
```

Starten Sie die gebaute Anwendung und testen Sie Ihre Abläufe auf jeder unterstützten Plattform, für die Sie ausliefern. Prüfen und committen Sie Quellcode, generierte Bindings, Moduldateien und Änderungen an der Frontend-Sperrdatei gemeinsam.

## Wenn das Upgrade fehlschlägt

Prüfen Sie die CLI im `PATH`, die Modulversion mit `go list -m github.com/wailsapp/wails/v3` und die installierte Runtime mit `npm --prefix frontend ls @wailsio/runtime`. Generieren Sie nach dem Beheben von Versionskonflikten die Bindings neu. Gehen Sie nicht davon aus, dass jede Alpha ohne Codeänderungen aktualisiert werden kann.

Wenn das Problem bestehen bleibt, [melden Sie ein reproduzierbares Issue](https://github.com/wailsapp/wails/issues/new/choose) mit der alten und neuen Version, dem genauen Fehler und der Ausgabe von `wails3 doctor`. Beachten Sie bei Sicherheitslücken die [Sicherheitsrichtlinie](https://github.com/wailsapp/wails/blob/master/SECURITY.md).
