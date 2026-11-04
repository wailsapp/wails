---
title: "Installation"
description: "Wails installieren und für die Anwendungsentwicklung vorbereiten"
slug: "quick-start/installation"
sourcePath: "quick-start/installation.md"
---

## Schnellinstallation (5 Minuten)

@note{type="tip" title="Kurzfassung – für erfahrene Entwickler"}
```bash
# Install Go 1.26+, then:
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
wails3 setup   # Interactive setup wizard (experimental)
```

Oder überprüfen Sie die Installation manuell mit `wails3 doctor`. [Weiter zur ersten App →](/quick-start/first-app/)

@end

## Schrittweise Installation

@steps
### Go installieren (erforderlich)
Wails v3 erfordert eine Go-Build-Toolchain ab Version 1.26. Installieren Sie eine aktuelle Patch-Version von Go 1.26. Ein älterer Go-Launcher kann bei `GOTOOLCHAIN=auto` automatisch eine geeignete Toolchain herunterladen. Sind automatische Downloads deaktiviert, müssen Sie die erforderliche Toolchain selbst installieren. Die `go.mod`-Deklaration des Starterprojekts senkt die Anforderung der Wails-Abhängigkeit nicht.

@tabs{sync-key="os"}
[Windows]
Laden Sie das Windows-Installationsprogramm von **[go.dev/dl](https://go.dev/dl/)** herunter und führen Sie es aus.

**Installation überprüfen:**

```powershell
go version  # Should show 1.26 or later
```

**PATH überprüfen:**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

Falls die Ausgabe leer ist, fügen Sie `C:\Users\YourName\go\bin` zu PATH hinzu.

[macOS]
**Option 1: Offizielles Installationsprogramm**

Laden Sie das macOS-Installationsprogramm (.pkg-Datei) von **[go.dev/dl](https://go.dev/dl/)** herunter und führen Sie es aus.

**Option 2: Homebrew**

```bash
brew install go
```

**Installation überprüfen:**

```bash
go version  # Should show 1.26 or later
echo $PATH | grep go/bin  # Should show ~/go/bin
```

Falls sich `~/go/bin` nicht in PATH befindet, fügen Sie es zu `~/.zshrc` oder `~/.bash_profile` hinzu:

```bash
export PATH=$PATH:~/go/bin
```

[Linux]
**Option 1: Offizielles Tarball**

Laden Sie das Linux-Tarball von **[go.dev/dl](https://go.dev/dl/)** herunter und führen Sie anschließend Folgendes aus:

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.26.0.linux-amd64.tar.gz
```

**Option 2: Paketverwaltung**

```bash
# Ubuntu/Debian
sudo apt install golang-go

# Fedora
sudo dnf install golang

# Arch
sudo pacman -S go
```

**Zu PATH hinzufügen** (zu `~/.bashrc` oder `~/.zshrc` hinzufügen):

```bash
export PATH=$PATH:/usr/local/go/bin:~/go/bin
source ~/.bashrc  # Reload
```

**Überprüfen:**

```bash
go version
echo $PATH | grep go/bin
```

@end

### Plattformabhängigkeiten installieren
@tabs{sync-key="os"}
[Windows]
**WebView2 Runtime** (normalerweise vorinstalliert)

Windows 10/11 enthält WebView2 standardmäßig. Falls es fehlt:

- Von [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/) herunterladen
- Oder führen Sie später `wails3 doctor` aus – der Assistent führt Sie durch die Installation

**Das ist alles!** Es sind keine weiteren Abhängigkeiten erforderlich.

@note{type="tip" title="Leistungstipp für Windows 11"}
Erwägen Sie, Ihre Projekte auf einem [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) zu speichern. Dev Drives sind für Entwicklungsworkloads optimiert und können die Build-Zeiten sowie die Geschwindigkeit von Festplattenzugriffen um bis zu 30 % deutlich verbessern.

@end

[macOS]
**Xcode Command Line Tools** (erforderlich)

```bash
xcode-select --install
```

Klicken Sie im angezeigten Dialogfeld auf „Installieren“.

**Überprüfen:**

```bash
xcode-select -p  # Should show /Library/Developer/CommandLineTools
```

**Das ist alles!** macOS enthält WebKit standardmäßig.

[Linux]
**Build-Werkzeuge und WebKit**

@note{type="caution" title="Mindestversionen der Distributionen"}
Wails v3 erfordert standardmäßig **WebKitGTK 6.0**. Distributionen, die nur WebKit2GTK 4.1 bereitstellen – Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39 und RHEL 9.x –, müssen mit der Legacy-Option `-tags gtk3` gebaut werden. Ältere Versionen, die nur WebKit2GTK 4.0 bereitstellen (Ubuntu 20.04, Debian 11 und RHEL 8), werden nicht unterstützt.

@end

@tabs{sync-key="distro"}
[Ubuntu/Debian]
Für den standardmäßigen GTK4-Stack ist mindestens Ubuntu 24.04 oder Debian 13 erforderlich.

```bash
sudo apt update
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
```

[Fedora]
```bash
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
```

[Arch]
```bash
sudo pacman -S base-devel gtk4 webkitgtk-6.0
```

[openSUSE]
```bash
sudo zypper install gcc pkg-config gtk4-devel webkitgtk-6_0-devel
```

[Gentoo]
```bash
sudo emerge --ask net-libs/webkit-gtk:6
```

[NixOS]
Fügen Sie Folgendes zu `shell.nix` oder `devShell` hinzu:

```nix
buildInputs = with pkgs; [ webkitgtk_6_0 gtk4 pkg-config gcc ];
```

[Sonstige]
Führen Sie nach der Installation von Wails `wails3 doctor` aus. Der Befehl zeigt die genauen Pakete an, die für Ihre Distribution erforderlich sind.

@end

@note{type="info" title="Legacy-GTK3-Stack"}
Falls Ihre Ziel-Distribution WebKitGTK 6.0 noch nicht bereitstellt (z. B. Ubuntu 22.04 LTS oder Debian 12), installieren Sie stattdessen die Entwicklungsbibliotheken für GTK3 und WebKit2GTK 4.1 (`libgtk-3-dev libwebkit2gtk-4.1-dev` unter Debian/Ubuntu; entsprechende Pakete bei anderen Distributionen) und bauen Sie mit `wails3 build -tags gtk3`. Der Legacy-Pfad wird bis einschließlich der Versionsreihe v3.0.x unterstützt und mit v3.1 entfernt. Weitere Informationen finden Sie unter [Linux-Paketierung – Legacy-Unterstützung für GTK3](/guides/build/linux/#legacy-gtk3-support).

@end

@end

### Wails CLI installieren
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

Dadurch wird der Befehl `wails3` in `~/go/bin` (unter Windows in `%USERPROFILE%\go\bin`) installiert.

### Einrichtungsassistenten ausführen (empfohlen)
```bash
wails3 setup
```

Der Einrichtungsassistent überprüft Ihre Abhängigkeiten, hilft bei der Installation fehlender Abhängigkeiten und konfiguriert die Projektstandards.

@note{type="caution" title="Experimentell"}
Der Einrichtungsassistent ist neu und wurde hauptsächlich unter Linux getestet. Falls Probleme auftreten, [melden Sie diese](https://github.com/wailsapp/wails/issues/4904) und verwenden Sie stattdessen `wails3 doctor`.

@end

### Installation überprüfen
```bash
wails3 doctor
```

**Erwartete Ausgabe (oder ähnlich):**

```
Wails (v3.0.0)  Wails Doctor

# System

┌──────────────────────────────────────────────────┐
| Name          | MacOS                            |
| Version       | 26.0                             |
| ID            | 25A354                           |
| Branding      | MacOS 26.0                       |
| Platform      | darwin                           |
| Architecture  | arm64                            |
| Apple Silicon | true                             |
| CPU           | Apple M2 Pro                     |
| CPU 1         | Apple M2 Pro                     |
| CPU 2         | Apple M2 Pro                     |
| GPU           | 16 cores, Metal Support: Metal 4 |
| Memory        | 16 GB                            |
└──────────────────────────────────────────────────┘

# Build Environment

┌─────────────┬─────────────────┐
| Wails CLI   | v3.0.0          |
| Go Version  | go1.26.0        |
└─────────────┴─────────────────┘

# Dependencies

┌─────────────────┬─────────────────────────────────────────────────┐
| npm             | 11.6.2                                          |
| *NSIS           | Not Installed. Install with `brew install...`.  |
| Xcode cli tools | 2412                                            |
└─────────────────┴─────────────────────────────────────────────────┘

# Checking for issues

SUCCESS No issues found

# Diagnosis

SUCCESS Your system is ready for Wails development!
```

@note{type="info" title="Falls der Befehl `wails3` nicht gefunden wird"}
Ihr `~/go/bin` befindet sich nicht in PATH. Beheben Sie dies wie oben in Schritt 1 beschrieben und starten Sie anschließend Ihr Terminal neu.

@end

### npm installieren (optional, aber empfohlen)
Die meisten Wails-Vorlagen verwenden npm für die Frontend-Werkzeuge.

@tabs{sync-key="os"}
[Windows]
Laden Sie das Installationsprogramm von [nodejs.org](https://nodejs.org/) herunter und führen Sie es aus.

**Überprüfen:**

```powershell
npm --version
```

[macOS]
**Option 1: Offizielles Installationsprogramm** Von [nodejs.org](https://nodejs.org/) herunterladen

**Option 2: Homebrew**

```bash
brew install node
```

**Überprüfen:**

```bash
npm --version
```

[Linux]
**Option 1: NodeSource**

```bash
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs  # Ubuntu/Debian
```

**Option 2: Paketverwaltung**

```bash
sudo dnf install nodejs  # Fedora
sudo pacman -S nodejs npm  # Arch
```

**Überprüfen:**

```bash
npm --version
```

@end

@note{type="tip" title="Alternative Paketmanager"}
Bevorzugen Sie `pnpm`, `yarn` oder `bun`? Kein Problem! Passen Sie einfach die `Taskfile.yml` in Ihrem Projekt an, um Ihr bevorzugtes Werkzeug zu verwenden.

@end

@end

## Fehlerbehebung

### Befehl `wails3` nicht gefunden

**Ursache:** `~/go/bin` (oder `%USERPROFILE%\go\bin`) ist nicht in Ihrem PATH enthalten.

**Lösung:**

@tabs{sync-key="os"}
[Windows]
1. Öffnen Sie „Umgebungsvariablen“ (über das Startmenü suchen).
2. Suchen Sie unter „Benutzervariablen“ nach `Path`.
3. Klicken Sie auf „Bearbeiten“ → „Neu“.
4. Fügen Sie Folgendes hinzu: `C:\Users\YourName\go\bin` (ersetzen Sie `YourName`).
5. Klicken Sie in allen Dialogfeldern auf „OK“.
6. **Starten Sie Ihr Terminal neu.**

**Überprüfen:**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

[macOS/Linux]
Fügen Sie Folgendes zu `~/.zshrc` (macOS) oder `~/.bashrc` (Linux) hinzu:

```bash
export PATH=$PATH:~/go/bin
```

Neu laden:

```bash
source ~/.zshrc  # or ~/.bashrc
```

**Überprüfen:**

```bash
echo $PATH | grep go/bin
wails3 version
```

@end

---

#### `wails3 doctor` meldet fehlende Abhängigkeiten

**Linux:** Die Ausgabe gibt genau an, welche Pakete Sie installieren müssen. Beispiel:

```
❌ webkit2gtk not found
   Install with: sudo apt install libwebkit2gtk-4.1-dev
```

**Windows:** Wenn WebView2 fehlt:

- Laden Sie es von [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/) herunter.
- Alternativ wird es beim Ausführen Ihrer ersten App automatisch installiert.

**macOS:** Wenn die Xcode-Werkzeuge fehlen:

```bash
xcode-select --install
```

---

#### Go-Version zu alt

Wails v3 erfordert Go 1.26 oder neuer. Wenn Sie eine ältere Version verwenden:

@tabs{sync-key="os"}
[Windows/macOS]
Laden Sie die neueste Version von [go.dev/dl](https://go.dev/dl/) herunter und installieren Sie sie erneut.

[Linux]
Laden Sie das neueste Tarball von [go.dev/dl](https://go.dev/dl/) herunter und führen Sie anschließend Folgendes aus:

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.26.0.linux-amd64.tar.gz
```

@end

## Entwicklungsversion (Bleeding Edge)

Möchten Sie den neuesten Code aus dem Hauptentwicklungszweig verwenden? Damit erhalten Sie vor der Veröffentlichung Zugriff auf neue Funktionen und Fehlerbehebungen, müssen jedoch mit möglichen Fehlern und inkompatiblen Änderungen rechnen. Dies wird nur Mitwirkenden oder Personen empfohlen, die kommende Funktionen testen müssen.

```bash
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout master
cd v3/cmd/wails3
go install
```

@note{type="caution" title="Entwicklungsversion"}
- Kann Fehler oder inkompatible Änderungen enthalten
- Erstellte Projekte verwenden die `replace`-Direktive, um auf die lokale Wails-Version zu verweisen
- Nur für Mitwirkende oder zum Testen neuer Funktionen empfohlen

@end

## Nächste Schritte

**Installation abgeschlossen!** Ihr System ist für die Entwicklung mit Wails bereit.

@cards{cols="1"}
🚀 Erstellen Sie Ihre erste App
Erstellen Sie in 10 Minuten eine funktionsfähige Anwendung.

[Tutorial: Erste App →](/quick-start/first-app/)

@end

@cards{cols="1"}
📖 Vorlagen erkunden
Sehen Sie, was sofort verfügbar ist.

```bash
wails3 init -l  # List templates
```

@end

---

**Haben Sie Probleme?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder [eröffnen Sie ein Issue](https://github.com/wailsapp/wails/issues).
