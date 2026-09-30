---
title: "Linux-Paketierung"
description: "Ihre Wails-Anwendung für die Linux-Distribution paketieren"
slug: "guides/build/linux"
sourcePath: "guides/build/linux.md"
---

## Paketformate

Paketieren Sie Ihre Anwendung für die Linux-Distribution:

```bash
wails3 package GOOS=linux
```

Dadurch werden im Verzeichnis `bin/` mehrere Formate erstellt:

- **AppImage**: Portabel, läuft auf jeder Linux-Distribution
- **DEB**: Für Debian, Ubuntu und deren Derivate
- **RPM**: Für Fedora, RHEL und deren Derivate
- **Arch**: Für Arch Linux und dessen Derivate

### Einzelne Formate

Erstellen Sie bestimmte Formate:

```bash
wails3 task linux:create:appimage
wails3 task linux:create:deb
wails3 task linux:create:rpm
wails3 task linux:create:aur
```

## Pakete anpassen

### Desktop-Eintrag

Die Datei `.desktop` legt fest, wie Ihre Anwendung in Anwendungsmenüs angezeigt wird. Sie wird aus den Werten in `build/linux/Taskfile.yml` generiert:

```yaml
vars:
  APP_NAME: 'MyApp'
  EXEC: 'MyApp'
  ICON: 'MyApp'
  CATEGORIES: 'Development;'
```

### Paketmetadaten

Bearbeiten Sie `build/linux/nfpm/nfpm.yaml`, um DEB- und RPM-Pakete anzupassen:

```yaml
name: myapp
version: 1.0.0
maintainer: Your Name <you@example.com>
description: My awesome Wails application
homepage: https://example.com
license: MIT
```

### AppImage

Die AppImage-Konfiguration befindet sich in `build/linux/appimage/`. Das Anwendungssymbol stammt aus `build/appicon.png`.

## Pakete signieren

Signieren Sie DEB- und RPM-Pakete mit einem PGP-Schlüssel:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=linux

# Or using tasks directly
wails3 task linux:sign:deb
wails3 task linux:sign:rpm
wails3 task linux:sign:packages  # Both
```

Konfigurieren Sie die Signierung in `build/linux/Taskfile.yml`:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  SIGN_ROLE: "builder"  # origin, maint, archive, or builder
```

Speichern Sie das Passwort Ihres Schlüssels:

```bash
wails3 setup signing
```

Weitere Informationen finden Sie unter [Anwendungen signieren](/guides/build/signing/).

## Für ARM bauen

```bash
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

@note{type="note"}
ARM64-Builds auf x86_64-Hosts verwenden Docker für die CGO-Cross-Kompilierung.

@end

## Unterstützung für das veraltete GTK3

Wails v3 baut standardmäßig auf **GTK4 mit WebKitGTK 6.0**. Für Distributionen, die WebKitGTK 6.0 noch nicht bereitstellen (Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x), ist weiterhin ein veralteter GTK3-/WebKit2GTK-4.1-Pfad verfügbar. Dieser veraltete Pfad muss über ein Build-Tag explizit aktiviert werden und soll in v3.1 entfernt werden.

@note{type="caution" title="Veralteter Pfad"}
Der GTK3-/WebKit2GTK-4.1-Pfad wird während der gesamten v3.0.x-Reihe unterstützt. Planen Sie die Migration zu GTK4 in Abstimmung mit der Verfügbarkeit von GTK4/WebKitGTK 6.0 in Ihrer Ziel-Distribution – `-tags gtk3` wird in v3.1 entfernt.

@end

### Abhängigkeiten

Installieren Sie die Entwicklungsbibliotheken für GTK3 und WebKit2GTK 4.1:

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

Die erforderlichen pkg-config-Pakete sind `gtk+-3.0` und `webkit2gtk-4.1`.

### Mit GTK3 bauen

Verwenden Sie das Flag `-tags gtk3`:

```bash
wails3 build -tags gtk3
```

Oder direkt mit Go:

```bash
go build -tags gtk3 -o myapp .
```

### Bekannte Unterschiede zu GTK4

- **Dateidialoge**: GTK4 verwendet standardmäßig `xdg-desktop-portal` für Dateidialoge. Daher verhalten sich einige Dialogoptionen, etwa das Standardverzeichnis und die Anzeige benutzerdefinierter Filter, anders als unter GTK3. Weitere Informationen finden Sie unter [Dialogreferenz – Verhalten von Dialogen unter Linux](/reference/dialogs/#linux-dialog-behavior).
- **Menüstil**: GTK4 unterstützt eine `LinuxMenuStylePrimaryMenu`-Option, die gemäß den GNOME HIG eine Hamburger-Schaltfläche (☰) in der Kopfleiste anzeigt. Diese Option hat bei `-tags gtk3`-Builds keine Wirkung. Siehe [Window API – Linux MenuStyle](/reference/window/#linux).
- **DPI-Skalierung**: GTK4 verwendet `gdk_monitor_get_scale` (GTK 4.14+) zur Unterstützung einer fraktionalen Skalierung.

### Build überprüfen

Führen Sie `wails3 doctor` aus, um Ihre Konfiguration zu überprüfen. Ohne Flags wird nach GTK4/WebKitGTK 6.0 gesucht (Standardeinstellung). Die veralteten GTK3-/WebKit2GTK-4.1-Pakete werden als optional aufgeführt.

## Fehlerbehebung

### AppImage wird nicht ausgeführt

Machen Sie die Datei ausführbar:

```bash
chmod +x MyApp-x86_64.AppImage
```

### Fehlende Abhängigkeiten

Wenn die Anwendung nicht startet, prüfen Sie, ob WebKit-Abhängigkeiten fehlen:

```bash
# Debian/Ubuntu
sudo apt install libwebkit2gtk-4.1-0

# Fedora
sudo dnf install webkit2gtk4.1

# Arch
sudo pacman -S webkit2gtk-4.1
```

### Kein C-Compiler gefunden

Das Build-System benötigt GCC oder Clang für CGO:

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

Alternativ können Sie `wails3 task setup:docker` ausführen; das Build-System verwendet dann automatisch Docker.

### Leeres oder weißes Fenster mit NVIDIA-GPU

Unter Linux mit proprietären NVIDIA-Treibern zeigen Wails-Anwendungen beim Start möglicherweise ein leeres oder weißes Fenster. Ursache ist ein Fehler in WebKitGTK, durch den der DMA-BUF-Renderer mit `gbm_bo_map()` und dem proprietären NVIDIA-Treiber fehlschlägt (betrifft X11 und Wayland, die Treiberversionen 377–580+ sowie GPUs der 10-Serie und ältere GT-710-Modelle).

**Wails wendet `WEBKIT_DISABLE_DMABUF_RENDERER=1` automatisch an**, wenn das NVIDIA-Kernelmodul (`/sys/module/nvidia`) erkannt wird. Daher müssen die meisten Benutzer nichts unternehmen.

Wenn weiterhin ein leeres Fenster angezeigt wird, beispielsweise in einem Container, in dem der Modulpfad nicht sichtbar ist, setzen Sie die Umgebungsvariable vor dem Start Ihrer Anwendung manuell:

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

Zugehörige Upstream-Fehler: [WebKit #262607](https://bugs.webkit.org/show_bug.cgi?id=262607), [WebKit #180739](https://bugs.webkit.org/show_bug.cgi?id=180739).

### AppImage-Kompatibilität beim Strippen

Auf modernen Linux-Distributionen (Arch Linux, Fedora 39+, Ubuntu 24.04+) werden Systembibliotheken mit `.relr.dyn`-ELF-Abschnitten kompiliert, um Relokationen effizienter auszuführen. Das zum Erstellen von AppImages verwendete Werkzeug `linuxdeploy` enthält eine ältere `strip`-Binärdatei, die diese modernen Abschnitte nicht verarbeiten kann.

Wails erkennt diese Situation automatisch, indem es vor dem Erstellen des AppImage die GTK-Systembibliotheken prüft. Wird sie erkannt, wird das Strippen deaktiviert (`NO_STRIP=1`), um die Kompatibilität sicherzustellen.

**Das bedeutet:**

- AppImages sind auf betroffenen Systemen etwas größer (~20-40 %).
- Die Funktionalität der Anwendung wird nicht beeinträchtigt.
- Dies wird automatisch gehandhabt – es ist keine Aktion erforderlich.

Wenn Sie auf modernen Systemen kleinere AppImages benötigen, können Sie eine neuere `strip`-Binärdatei installieren und `linuxdeploy` so konfigurieren, dass diese anstelle der mitgelieferten Version verwendet wird.
