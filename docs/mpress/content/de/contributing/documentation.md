---
title: "Dokumentation mit M-Press verfassen"
description: "M-Press herunterladen und installieren, eine Vorschau der Wails-Dokumentation anzeigen und einen Dokumentationsbeitrag einreichen."
sourcePath: "contributing/documentation.md"
---

Wails v3 verwendet **M-Press**, um aus Markdown diese Dokumentationswebsite zu erstellen. Sie können einen Tippfehler beheben, eine Erklärung verbessern, ein Beispiel oder eine Anleitung hinzufügen oder eine Übersetzung korrigieren. Für reine Dokumentationskorrekturen sind weder ein vorheriges Issue noch ein fehlgeschlagener Codetest erforderlich.

Sie benötigen Git, um einen Beitrag vorzubereiten, und ein GitHub-Konto, um einen Pull Request einzureichen. Die herunterladbare ausführbare M-Press-Datei enthält den Dokumentationsgenerator und den Diagramm-Renderer. Für die Vorschau der Dokumentation benötigen Sie weder Go noch Node.js, die Wails CLI oder Abhängigkeiten für Desktop-Builds. Python 3 ist für die zusätzlichen Validierungsskripte des Repositorys erforderlich.

## Von einer Dokumentationsseite aus beginnen

Öffnen Sie in der Seitenleiste **Mitwirken → Dokumentation**, um diese Anleitung aufzurufen. Um an einer bestimmten Seite zu arbeiten, wählen Sie in deren Kopfzeile das **Bearbeitungssymbol** und anschließend **Diese Seite lokal bearbeiten** aus. Kopieren Sie den Befehl für Ihr Betriebssystem und führen Sie ihn in einem Terminal aus.

Der Befehl verwendet M-Press, sofern es installiert ist. Andernfalls lädt das Installationsprogramm ein offizielles Release herunter und überprüft dessen Prüfsumme. Es richtet eine lokale Arbeitskopie für den Beitrag ein und öffnet die ausgewählte Seite in M-Press. Prüfen Sie `mpress version` anhand der unten festgelegten Version, bevor Sie die Repository-Prüfungen ausführen: Das automatische Installationsprogramm kann ein neueres Release herunterladen.

Sie können die Quelldatei in Ihrem eigenen Editor bearbeiten, während die lokale Vorschau neu erstellt wird. Der Assistent bietet außerdem **Änderungen überprüfen**, **Änderungen committen**, **Beitrags-Branch pushen** und **Pull Request als Entwurf öffnen** als separate Aktionen an. Befolgen Sie vor dem Einreichen die nachstehenden Hinweise zu Prüfungen und Pull Requests. Das Generieren von Übersetzungen ist optional und erfordert möglicherweise einen konfigurierten Anbieter; für das Bearbeiten und Anzeigen einer Vorschau von Markdown ist dies nicht erforderlich.

Führen Sie für eine reproduzierbare manuelle Einrichtung die folgenden Schritte aus.

## M-Press herunterladen und installieren

Das Build-Skript des Repositorys lädt die **neueste stabile M-Press-Version** herunter. Verwenden Sie lokal dieselbe Version: Wählen Sie unten eine Installationsmethode aus, führen Sie `mpress version` aus und vergleichen Sie die Ausgabe mit der auf der Release-Seite angegebenen Version.

### Eine ausführbare Datei herunterladen

Öffnen Sie die Seite der [neuesten stabilen M-Press-Version](https://github.com/leaanthony/mpress/releases/latest) oder laden Sie unten das passende Archiv herunter. ARM64 umfasst Apple Silicon; AMD64 ist die übliche Architektur für PCs mit Intel- oder AMD-Prozessor.

| Computer | AMD64 / Intel | ARM64 / Apple Silicon |
| --- | --- | --- |
| macOS | [mpress-darwin-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-amd64.tar.gz) | [mpress-darwin-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-arm64.tar.gz) |
| Linux | [mpress-linux-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-amd64.tar.gz) | [mpress-linux-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-arm64.tar.gz) |
| Windows | [mpress-windows-amd64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-amd64.zip) | [mpress-windows-arm64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-arm64.zip) |

Laden Sie [checksums.txt](https://github.com/leaanthony/mpress/releases/latest/download/checksums.txt) aus demselben Release herunter. Berechnen Sie in Ihrem Downloadverzeichnis den SHA-256-Hash des Archivs und vergleichen Sie ihn mit der Zeile für diesen Dateinamen in `checksums.txt`:

@tabs
[macOS]
```sh
shasum -a 256 mpress-darwin-arm64.tar.gz
```

[Linux]
```sh
sha256sum mpress-linux-amd64.tar.gz
```

[Windows PowerShell]
```powershell
Get-FileHash .\mpress-windows-amd64.zip -Algorithm SHA256
```
@end

Verwenden Sie den tatsächlichen Namen Ihrer heruntergeladenen Datei, falls Sie eine andere Architektur gewählt haben. Wenn die Hashes übereinstimmen, entpacken und installieren Sie die ausführbare Datei:

@tabs
[macOS]
Aus dem Downloadverzeichnis mit dem Apple-Silicon-Archiv:

```sh
mkdir -p mpress-download
tar -xzf mpress-darwin-arm64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

Ersetzen Sie für Intel-Macs `mpress-darwin-amd64.tar.gz`. Fügen Sie die Zeile `export PATH` zu `~/.zshrc` hinzu, damit der Befehl auch in neuen Terminals verfügbar ist.

[Linux]
Aus dem Downloadverzeichnis:

```sh
mkdir -p mpress-download
tar -xzf mpress-linux-amd64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

Ersetzen Sie für ARM64 `mpress-linux-arm64.tar.gz`. Fügen Sie die Zeile `export PATH` zur Startdatei Ihrer Shell hinzu, beispielsweise zu `~/.bashrc`, damit der Befehl auch in neuen Terminals verfügbar ist.

[Windows PowerShell]
Aus dem Downloadverzeichnis:

```powershell
Expand-Archive .\mpress-windows-amd64.zip -DestinationPath .\mpress-download
New-Item -ItemType Directory -Force "$HOME\bin"
Copy-Item .\mpress-download\mpress.exe "$HOME\bin\mpress.exe"
$env:Path = "$HOME\bin;$env:Path"
mpress version
```

Ersetzen Sie für ARM64 `mpress-windows-arm64.zip`. Fügen Sie Ihr benutzerspezifisches Verzeichnis `bin` unter Windows in den **Umgebungsvariablen** zum benutzerspezifischen **Path** hinzu und öffnen Sie das Terminal erneut, damit der Befehl auch in zukünftigen Sitzungen verfügbar ist.
@end

### Mit Go installieren

Wenn Sie Go bereits verwenden, können Sie stattdessen die neueste stabile Version erstellen. Prüfen Sie in `go.mod`, welche Toolchain erforderlich ist; die automatische Toolchain-Auswahl von Go lädt sie gegebenenfalls herunter. Die für diesen Dokumentationsdurchlauf verwendete Version erfordert Go 1.26.10 oder neuer.

```sh
go install github.com/leaanthony/mpress/cmd/mpress@latest
mpress version
```

Go installiert den Befehl in `GOBIN` oder, falls `GOBIN` nicht gesetzt ist, in das Verzeichnis `bin` unter `go env GOPATH`. Fügen Sie dieses Verzeichnis zu Ihrem `PATH` hinzu und öffnen Sie das Terminal erneut, falls `mpress` nicht gefunden wird.

## Dokumentationsquellen abrufen

[Forken Sie wailsapp/wails](https://github.com/wailsapp/wails/fork) und führen Sie anschließend diese Befehle aus. Ersetzen Sie dabei `YOUR_USERNAME` durch Ihren GitHub-Benutzernamen:

```sh
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
git fetch upstream master
git switch -c docs/my-change upstream/master
```

Damit beginnt Ihr Dokumentations-Branch beim Upstream-Branch **master**. Verwenden Sie anstelle von `docs/my-change` einen eigenen aussagekräftigen Branch-Namen.

## Live-Vorschau ausführen

Führen Sie den Befehl im Stammverzeichnis des Wails-Repositorys aus, in dem sich `mpress.yaml` befindet:

```sh
mpress dev
```

Öffnen Sie [http://127.0.0.1:3000](http://127.0.0.1:3000) in Ihrem Browser. Rufen Sie für diese Anleitung [http://127.0.0.1:3000/contributing/documentation/](http://127.0.0.1:3000/contributing/documentation/) auf. Lassen Sie den Befehl während der Bearbeitung weiterlaufen; gespeicherte Änderungen an den Quelldateien erstellen die Vorschau neu. Drücken Sie **Ctrl+C**, um den Prozess zu beenden, bevor Sie die nachstehenden Validierungsbefehle ausführen.

Das Repository enthält bereits die M-Press-Konfiguration. Es muss keine neue Website initialisiert werden. Um einen anderen Port zu verwenden, führen Sie `mpress dev --port 3001` aus und öffnen Sie die entsprechende localhost-Adresse.

## Eine Seite bearbeiten oder hinzufügen

Die englischen Quellseiten befinden sich in `docs/mpress/content/`. Übersetzungen befinden sich in Sprachordnern wie `fr/` und `id/`. Die über das Bearbeitungssymbol aufgerufenen Anweisungen für Beiträge geben die Quelldatei der aktuellen Seite an.

Bearbeiten Sie bei einer vorhandenen Seite deren Markdown und behalten Sie die YAML-Metadaten am Anfang bei. Verwenden Sie Überschriften, Absätze, Listen, Links und mit Begrenzungszeichen versehene Codeblöcke. Behalten Sie zusammengehörige M-Press-Komponenten wie `@note` / `@end` und `@tabs` / `@end` bei. Die [Dokumentation zur Inhaltserstellung mit M-Press](https://m-press.me/authoring/) beschreibt die unterstützte Syntax.

Erstellen Sie für eine neue Anleitung eine Datei wie `docs/mpress/content/guides/my-guide.md`:

```markdown
---
title: "My guide"
description: "The task this guide helps someone complete."
slug: "guides/my-guide"
sourcePath: "guides/my-guide.md"
---

Explain the task, give working steps, and show the expected result.
```

Fügen Sie die Seite dem passenden Abschnitt von `docs/mpress/content/_nav.yaml` hinzu, sofern dieser Abschnitt nicht bereits `autogenerate` verwendet. Prüfen Sie, ob Leser die Seite über die Navigation und über Links auf verwandten Seiten erreichen können. Für neue englische Seiten sind außerdem vollständige Kopien für alle Sprachversionen erforderlich, bevor die Abdeckungsprüfungen erfolgreich sein können; stimmen Sie diese mit den Maintainern oder Übersetzungsmitwirkenden ab.

Die generierte Ausgabe befindet sich in `docs/mpress/site/`. Committen Sie die Quellseiten sowie alle beabsichtigten Änderungen an Navigation oder statischen Assets; generierte Websitedateien und lokale Caches gehören nicht in den Pull Request.

## Beitrag prüfen

Aus dem Stammverzeichnis des Repositorys:

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Verwenden Sie unter Windows `python`, falls Python 3 auf diese Weise installiert ist. Wenn [Task](https://taskfile.dev/) bereits installiert ist, führt `task docs:build` dieselben Prüfungen mit dem `mpress` in Ihrem PATH aus.

Die Übersetzungsprüfung schlägt bei fehlenden, leeren oder verwaisten Übersetzungsdateien fehl. Veraltete oder beschädigte übersetzte Passagen werden standardmäßig als Warnungen gemeldet; mit `--strict` führen auch diese Befunde zum Fehlschlagen. Eine Änderung an bestehendem englischem Text kann eingereicht werden, während die Übersetzungen synchronisiert werden. Bei der Korrektur einer einzelnen Übersetzung kann nur die betreffende Sprache geändert werden. Befehle, API-Namen, Links, Codebeispiele und Diagrammverbindungen müssen unverändert bleiben; verwenden Sie keine englischen Platzhalter als übersetzten Inhalt.

Lesen Sie die geänderte Seite im Browser und führen Sie alle von Ihnen geänderten Codebeispiele aus. Wenn eine Prüfung bei einer unveränderten Seite fehlschlägt, fügen Sie die Ausgabe und den relevanten Kontext Ihrem PR hinzu, statt die Prüfung abzuschwächen.

## Pull Request einreichen

Prüfen Sie Ihren Diff, committen Sie Ihre Änderungen am Quelltext und pushen Sie Ihren Beitrags-Branch in Ihren Fork. Öffnen Sie in `wailsapp/wails` einen PR gegen **master**. Beschreiben Sie die Änderungen, verlinken Sie die betroffene Seite und führen Sie die ausgeführten Prüfungen auf. Fügen Sie bei Layoutänderungen Vorher-/Nachher-Screenshots und bei geänderten Codebeispielen Angaben zu Plattform und Version hinzu.

Die Produktionswebsite ist statisch; für die lokale Vorschau und Validierung sind weder Anmeldedaten für die Bereitstellung noch ein Übersetzungsanbieter erforderlich. Die Maintainer übernehmen die Veröffentlichung nach der Überprüfung.

Informationen zu Codebeiträgen und Funktionsvorschlägen finden Sie unter [Beiträge zu Wails](/contributing/). Informationen zu den Interna des Frameworks finden Sie in der [Technischen Übersicht](/contributing/overview/).
