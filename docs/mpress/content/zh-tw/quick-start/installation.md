---
title: "安裝"
description: "安裝 Wails，準備開始建置應用程式"
slug: "quick-start/installation"
sourcePath: "quick-start/installation.md"
---

## 快速安裝（5 分鐘）

@note{type="tip" title="摘要：適合有經驗的開發人員"}
```bash
# Install Go 1.26+, then:
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
wails3 setup   # Interactive setup wizard (experimental)
```

或使用`wails3 doctor`手動驗證。[跳至第一個應用程式 →](/quick-start/first-app/)

@end

## 逐步安裝

@steps
### 安裝 Go（必要）
Wails v3需要 Go 1.26或更新版本的建置工具鏈。請安裝目前的 Go 1.26修補版本。執行`GOTOOLCHAIN=auto`時，較舊的 Go 啟動程式可能會自動下載適用的工具鏈；若已停用自動下載，請自行安裝所需的工具鏈。起始專案本身的`go.mod`宣告不會降低 Wails 相依套件的版本要求。

@tabs{sync-key="os"}
[Windows]
從<strong>[go.dev/dl](https://go.dev/dl/)</strong>下載 Windows 安裝程式並執行。

**驗證安裝：**

```powershell
go version  # Should show 1.26 or later
```

**檢查 PATH：**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

若輸出為空，請將`C:\Users\YourName\go\bin`加入 PATH。

[macOS]
**選項 1：官方安裝程式**

從<strong>[go.dev/dl](https://go.dev/dl/)</strong>下載 macOS 安裝程式（.pkg 檔案）並執行。

**選項 2：Homebrew**

```bash
brew install go
```

**驗證安裝：**

```bash
go version  # Should show 1.26 or later
echo $PATH | grep go/bin  # Should show ~/go/bin
```

若 PATH 中沒有`~/go/bin`，請將它加入`~/.zshrc`或`~/.bash_profile`：

```bash
export PATH=$PATH:~/go/bin
```

[Linux]
**選項 1：官方 Tarball**

從<strong>[go.dev/dl](https://go.dev/dl/)</strong>下載 Linux tarball，然後執行：

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.26.0.linux-amd64.tar.gz
```

**選項 2：套件管理員**

```bash
# Ubuntu/Debian
sudo apt install golang-go

# Fedora
sudo dnf install golang

# Arch
sudo pacman -S go
```

**加入 PATH**（加入`~/.bashrc`或`~/.zshrc`）：

```bash
export PATH=$PATH:/usr/local/go/bin:~/go/bin
source ~/.bashrc  # Reload
```

**驗證：**

```bash
go version
echo $PATH | grep go/bin
```

@end

### 安裝平台相依套件
@tabs{sync-key="os"}
[Windows]
**WebView2 Runtime**（通常已預先安裝）

Windows 10/11預設包含 WebView2。若系統缺少此元件：

- 從[Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)下載
- 或稍後執行`wails3 doctor`，它會引導您完成設定

<strong>這樣就完成了！</strong>不需要其他相依套件。

@note{type="tip" title="Windows 11效能提示"}
建議使用[Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/)儲存專案。Dev Drive 針對開發人員工作負載進行最佳化，最多可將建置時間和磁碟存取速度大幅提升30%。

@end

[macOS]
**Xcode Command Line Tools**（必要）

```bash
xcode-select --install
```

在顯示的對話方塊中按一下「安裝」。

**驗證：**

```bash
xcode-select -p  # Should show /Library/Developer/CommandLineTools
```

<strong>這樣就完成了！</strong>macOS 預設包含 WebKit。

[Linux]
**建置工具與 WebKit**

@note{type="caution" title="最低發行版版本"}
Wails v3預設需要**WebKitGTK 6.0**。僅隨附 WebKit2GTK 4.1的發行版（Ubuntu 22.04 LTS、Debian 12、Fedora ≤ 39、RHEL 9.x）必須選擇使用舊版`-tags gtk3`進行建置。不支援僅隨附 WebKit2GTK 4.0的更舊版本（Ubuntu 20.04、Debian 11、RHEL 8）。

@end

@tabs{sync-key="distro"}
[Ubuntu/Debian]
預設 GTK4 技術堆疊需要 Ubuntu 24.04+ 或 Debian 13+。

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
加入您的`shell.nix`或`devShell`：

```nix
buildInputs = with pkgs; [ webkitgtk_6_0 gtk4 pkg-config gcc ];
```

[其他]
安裝 Wails 後執行`wails3 doctor`，它會顯示您的發行版所需的確切套件。

@end

@note{type="info" title="舊版 GTK3 技術堆疊"}
若目標發行版尚未隨附 WebKitGTK 6.0（例如 Ubuntu 22.04 LTS、Debian 12），請改為安裝 GTK3 + WebKit2GTK 4.1開發程式庫（Debian/Ubuntu 使用`libgtk-3-dev libwebkit2gtk-4.1-dev`；其他發行版則使用對應套件），並以`wails3 build -tags gtk3`進行建置。舊版路徑支援至v3.0.x 系列，並將於v3.1移除。詳情請參閱[Linux 封裝—舊版 GTK3 支援](/guides/build/linux/#legacy-gtk3-support)。

@end

@end

### 安裝 Wails CLI
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

這會將`wails3`命令安裝至`~/go/bin`（Windows 上則為`%USERPROFILE%\go\bin`）。

### 執行設定精靈（建議）
```bash
wails3 setup
```

設定精靈會檢查相依套件、協助安裝缺少的套件，並設定專案預設值。

@note{type="caution" title="實驗性功能"}
設定精靈是新功能，目前主要在 Linux 上測試。若遇到問題，請[回報問題](https://github.com/wailsapp/wails/issues/4904)，並改用`wails3 doctor`。

@end

### 驗證安裝
```bash
wails3 doctor
```

**預期輸出（或類似內容）：**

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

@note{type="info" title="若找不到 `wails3` 命令"}
您的`~/go/bin`不在 PATH 中。請參閱上方步驟1進行修正，然後重新啟動終端機。

@end

### 安裝 npm（選用但建議）
大多數 Wails 範本使用 npm 作為前端工具。

@tabs{sync-key="os"}
[Windows]
從[nodejs.org](https://nodejs.org/)下載並執行安裝程式。

**驗證：**

```powershell
npm --version
```

[macOS]
**選項 1：官方安裝程式** 從[nodejs.org](https://nodejs.org/)下載

**選項 2：Homebrew**

```bash
brew install node
```

**驗證：**

```bash
npm --version
```

[Linux]
**選項 1：NodeSource**

```bash
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs  # Ubuntu/Debian
```

**選項 2：套件管理員**

```bash
sudo dnf install nodejs  # Fedora
sudo pacman -S nodejs npm  # Arch
```

**驗證：**

```bash
npm --version
```

@end

@note{type="tip" title="其他套件管理工具"}
偏好使用`pnpm`、`yarn`或`bun`？沒問題！只要更新專案中的`Taskfile.yml`，改用您偏好的工具即可。

@end

@end

## 疑難排解

### 找不到`wails3`命令

**原因：**`~/go/bin`（或`%USERPROFILE%\go\bin`）不在您的 PATH 中。

**解決方法：**

@tabs{sync-key="os"}
[Windows]
1. 開啟「環境變數」（在「開始」功能表中搜尋）
2. 在「使用者變數」下方找到`Path`
3. 按一下「編輯」→「新增」
4. 新增：`C:\Users\YourName\go\bin`（取代`YourName`）
5. 在所有對話方塊中按一下「確定」
6. **重新啟動終端機**

**驗證：**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

[macOS/Linux]
新增至`~/.zshrc`（macOS）或`~/.bashrc`（Linux）：

```bash
export PATH=$PATH:~/go/bin
```

重新載入：

```bash
source ~/.zshrc  # or ~/.bashrc
```

**驗證：**

```bash
echo $PATH | grep go/bin
wails3 version
```

@end

---

#### `wails3 doctor`回報缺少相依套件

<strong>Linux：</strong>輸出會明確指出需要安裝哪些套件。例如：

```
❌ webkit2gtk not found
   Install with: sudo apt install libwebkit2gtk-4.1-dev
```

<strong>Windows：</strong>如果缺少 WebView2：

- 從[Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)下載
- 或者，在您執行第一個應用程式時會自動安裝

<strong>macOS：</strong>如果缺少 Xcode 工具：

```bash
xcode-select --install
```

---

#### Go 版本過舊

Wails v3需要 Go 1.26以上版本。如果您使用的是較舊版本：

@tabs{sync-key="os"}
[Windows/macOS]
從[go.dev/dl](https://go.dev/dl/)下載最新版本並重新安裝。

[Linux]
從[go.dev/dl](https://go.dev/dl/)下載最新的 tarball，然後執行：

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.26.0.linux-amd64.tar.gz
```

@end

## 開發版本（最前沿版本）

想使用主要開發分支上的最新程式碼嗎？這可讓您在新功能與修正正式發布前搶先使用，但也可能遇到錯誤或破壞性變更。僅建議貢獻者或需要測試即將推出之功能的人員使用。

```bash
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout master
cd v3/cmd/wails3
go install
```

@note{type="caution" title="開發版本"}
- 可能含有錯誤或破壞性變更
- 建立的專案將使用`replace`指令指向本機的 Wails
- 僅建議貢獻者或測試新功能時使用

@end

## 後續步驟

<strong>安裝完成！</strong>您的系統已準備好進行 Wails 開發。

@cards{cols="1"}
🚀 建立您的第一個應用程式
在10分鐘內建立可運作的應用程式。

[第一個應用程式教學 →](/quick-start/first-app/)

@end

@cards{cols="1"}
📖 探索範本
查看有哪些可立即使用的內容。

```bash
wails3 init -l  # List templates
```

@end

---

<strong>遇到問題嗎？</strong>請到[Discord](https://discord.gg/JDdSxwjhGf)提問，或[建立問題](https://github.com/wailsapp/wails/issues)。
