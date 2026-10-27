---
title: "Linux 封裝"
description: "封裝 Wails 應用程式以發佈至 Linux"
slug: "guides/build/linux"
sourcePath: "guides/build/linux.md"
---

## 套件格式

封裝應用程式以發佈至 Linux：

```bash
wails3 package GOOS=linux
```

這會在`bin/`目錄中建立多種格式：

- **AppImage**：可攜式，可在任何 Linux 發行版上執行
- **DEB**：適用於 Debian、Ubuntu 及其衍生發行版
- **RPM**：適用於 Fedora、RHEL 及其衍生發行版
- **Arch**：適用於 Arch Linux 及其衍生發行版

### 個別格式

建置特定格式：

```bash
wails3 task linux:create:appimage
wails3 task linux:create:deb
wails3 task linux:create:rpm
wails3 task linux:create:aur
```

## 自訂套件

### 桌面項目

`.desktop`檔案控制應用程式在應用程式選單中的顯示方式。此檔案會根據`build/linux/Taskfile.yml`中的值產生：

```yaml
vars:
  APP_NAME: 'MyApp'
  EXEC: 'MyApp'
  ICON: 'MyApp'
  CATEGORIES: 'Development;'
```

### 套件中繼資料

編輯`build/linux/nfpm/nfpm.yaml`以自訂 DEB 和 RPM 套件：

```yaml
name: myapp
version: 1.0.0
maintainer: Your Name <you@example.com>
description: My awesome Wails application
homepage: https://example.com
license: MIT
```

### AppImage

AppImage 的設定位於`build/linux/appimage/`。應用程式圖示來自`build/appicon.png`。

AppImage 會將 WebKit 及其輔助程序封裝在內，並修補封裝的程式庫，使其從 AppImage 內啟動這些程序。它的`AppRun`會從 AppImage 內的`usr/lib`啟動應用程式；修補後的路徑在 WebKit 的 bubblewrap 沙箱外部與內部都能正確解析，因此沙箱會維持啟用。只有在您自己的`AppRun`能保留該工作目錄時，才應加以提供。

## 簽署套件

使用 PGP 金鑰簽署 DEB 和 RPM 套件：

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=linux

# Or using tasks directly
wails3 task linux:sign:deb
wails3 task linux:sign:rpm
wails3 task linux:sign:packages  # Both
```

在`build/linux/Taskfile.yml`中設定簽署：

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  SIGN_ROLE: "builder"  # origin, maint, archive, or builder
```

儲存金鑰密碼：

```bash
wails3 signing
```

詳情請參閱[簽署應用程式](/guides/build/signing/)。

## 針對 ARM 建置

```bash
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

@note{type="note"}
從 x86_64 主機進行 ARM64 建置時，會使用 Docker 執行 CGO 交叉編譯。

@end

## 舊版 GTK3 支援

Wails v3預設以**GTK4 搭配 WebKitGTK 6.0**進行建置。對於尚未提供 WebKitGTK 6.0的發行版（Ubuntu 22.04 LTS、Debian 12、Fedora ≤ 39、RHEL 9.x），仍可使用舊版 GTK3／WebKit2GTK 4.1建置途徑。舊版途徑必須透過建置標籤選擇啟用，並預計於v3.1移除。

@note{type="caution" title="舊版途徑"}
GTK3／WebKit2GTK 4.1途徑會支援至v3.0.x 版本線。請配合目標發行版提供 GTK4／WebKitGTK 6.0的時程規劃遷移至 GTK4，因為`-tags gtk3`將於v3.1移除。

@end

### 相依套件

安裝 GTK3 和 WebKit2GTK 4.1開發程式庫：

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

必要的 pkg-config 套件為`gtk+-3.0`和`webkit2gtk-4.1`。

### 使用 GTK3 建置

使用`-tags gtk3`旗標：

```bash
wails3 build -tags gtk3
```

或直接使用 Go：

```bash
go build -tags gtk3 -o myapp .
```

### 與 GTK4 的已知差異

- **檔案對話方塊**：GTK4 預設使用`xdg-desktop-portal`顯示檔案對話方塊，因此某些對話方塊選項（例如預設目錄、自訂篩選器的顯示方式）的行為與 GTK3 不同。詳情請參閱[對話方塊參考資料－Linux 對話方塊行為](/reference/dialogs/#linux-dialog-behavior)。
- **選單樣式**：GTK4 支援`LinuxMenuStylePrimaryMenu`選項，會依循 GNOME HIG 在標題列顯示漢堡選單按鈕（☰）。此選項不會影響`-tags gtk3`建置。請參閱[Window API－Linux MenuStyle](/reference/window/#linux)。
- **DPI 縮放**：GTK4 使用`gdk_monitor_get_scale`（GTK 4.14+）來支援非整數縮放。

### 檢查建置環境

執行`wails3 doctor`以驗證設定。未指定旗標時，它會檢查 GTK4／WebKitGTK 6.0（預設值）。舊版 GTK3／WebKit2GTK 4.1套件則會列為選用項目。

## 疑難排解

### AppImage 無法執行

將其設為可執行：

```bash
chmod +x MyApp-x86_64.AppImage
```

### 缺少相依套件

請安裝與二進位檔建置變體相符的程式庫。對於預設的 GTK4／WebKitGTK 6.0建置：

```bash
# Debian/Ubuntu versions that provide WebKitGTK 6.0
sudo apt install libgtk-4-1 libwebkitgtk-6.0-4
# Fedora
sudo dnf install gtk4 webkitgtk6.0
# Arch
sudo pacman -S gtk4 webkitgtk-6.0
```

對於明確使用`-tags gtk3`建置的二進位檔，請改為安裝舊版 GTK3／WebKit2GTK 4.1程式庫：

```bash
# Debian/Ubuntu
sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0
# Fedora
sudo dnf install gtk3 webkit2gtk4.1
# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

套件是否可用及其名稱取決於發行版版本。請使用具備必要程式庫的受支援版本，或使用文件記載的交叉建置／容器途徑；安裝舊版程式庫並不能滿足預設建置的需求。AppImage 的可攜性取決於其封裝實際包含哪些項目，以及仍需使用哪些主機程式庫。

### 找不到 C 編譯器

建置系統需要 GCC 或 Clang 才能使用 CGO：

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

或者，執行`wails3 task setup:docker`，建置系統便會自動使用 Docker。

### 在 NVIDIA GPU 上顯示空白或全白視窗

在使用 NVIDIA 專有驅動程式的 Linux 上，Wails 應用程式啟動時可能顯示空白或全白視窗。這是 WebKitGTK 的錯誤所致：DMA-BUF 轉譯器搭配 NVIDIA 專有驅動程式時會因`gbm_bo_map()`而失敗（影響 X11 和 Wayland、377–580+ 版驅動程式、10系列及更舊的 GT 710 GPU）。

**Wails 偵測到 NVIDIA 核心模組（`/sys/module/nvidia`）時，會`WEBKIT_DISABLE_DMABUF_RENDERER=1`自動套用**，因此大多數使用者不需要採取任何動作。

如果仍看到空白視窗（例如在無法看見模組路徑的容器中），請先手動設定環境變數，再啟動應用程式：

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

相關的上游錯誤：[WebKit #262607](https://bugs.webkit.org/show_bug.cgi?id=262607)、[WebKit #180739](https://bugs.webkit.org/show_bug.cgi?id=180739)。

### AppImage strip 相容性

在現代 Linux 發行版（Arch Linux、Fedora 39+、Ubuntu 24.04+）上，系統程式庫會使用`.relr.dyn` ELF 區段編譯，以提高重定位效率。用於建立 AppImage 的`linuxdeploy`工具內含較舊的`strip`二進位檔，無法處理這些現代區段。

Wails 會在建置 AppImage 前檢查系統 GTK 程式庫，自動偵測此情況。偵測到時，會停用符號剝除（`NO_STRIP=1`）以確保相容性。

**這表示：**

- 在受影響的系統上，AppImage 會稍大一些（約增加20-40%）
- 應用程式功能不受影響
- 此情況會自動處理，無須採取任何動作

若需要在現代系統上產生較小的 AppImage，可以安裝較新的`strip`二進位檔，並設定`linuxdeploy`使用該檔案，而非其內含版本。
