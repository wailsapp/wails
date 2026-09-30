---
title: "安裝"
description: "安裝 Wails 並設定開發環境"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## 支援的平台

- Windows AMD64/ARM64
- macOS 10.15+ AMD64（可部署至 macOS 10.13+）
- macOS 11.0+ ARM64
- Ubuntu 24.04 AMD64/ARM64（其他 Linux 發行版也可能適用！）

## 相依套件

安裝 Wails 前，必須先安裝一些共用相依套件。

@note{type="tip"}
安裝 Wails CLI 後，可以執行`wails3 setup`，自動檢查這些相依套件並協助安裝。

@end

@tabs
[Go（至少為1.24）]
請從[Go 下載頁面](https://go.dev/dl/)下載 Go。

請務必遵循官方的[Go 安裝說明](https://go.dev/doc/install)。另請確保`PATH`環境變數中也包含`~/go/bin`目錄的路徑。重新啟動終端機，然後進行下列檢查：

- 檢查 Go 是否已正確安裝：`go version`
- 檢查`~/go/bin`是否位於 PATH 變數中
  - Mac / Linux：`echo $PATH | grep go/bin`
  - Windows：`$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm（選用）]
雖然 Wails 不要求安裝 npm，但大多數隨附範本都需要使用 npm。

請從[Node 下載頁面](https://nodejs.org/en/download/)下載最新的 Node 安裝程式。建議使用最新版本，因為我們通常以該版本進行測試。

執行`npm --version`以進行驗證。

@note{type="info"}
如果偏好使用 npm 以外的套件管理工具，也可以自行選用。你需要更新專案的 Taskfile，以改用該工具。

@end

@end

## 平台特定相依套件

你也需要安裝平台特定的相依套件：

@tabs{sync-key="platform"}
[Mac]
Wails 要求安裝 Xcode 命令列工具。可以執行以下命令進行安裝：

```sh
xcode-select --install
```

[Windows]
Wails 要求安裝[WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)。幾乎所有 Windows 系統都已安裝此執行階段。你可以使用`wails doctor`命令進行檢查。

[Linux]
Linux 需要標準的`gcc`建置工具，以及`gtk4`和`webkitgtk-6.0`。安裝後執行<code>wails3 doctor</code>，即可查看如何安裝這些相依套件。舊版 GTK3 / WebKit2GTK 4.1技術堆疊仍可透過`-tags gtk3`使用（請參閱[Linux 封裝－舊版 GTK3 支援](/guides/build/linux/#legacy-gtk3-support)），直到 v3.1為止。如果你的發行版／套件管理工具不受支援，請在 Discord 上告訴我們。

@end

## 安裝

若要使用 Go Modules 安裝 Wails CLI，請執行下列命令：

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest
```

若要安裝最新的開發版本，請執行下列命令：

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

使用開發版本時，所有產生的專案都會使用 Go 的[replace](https://go.dev/ref/mod#go-mod-file-replace)指令，以確保專案使用 Wails 的開發版本。

## 後續步驟

安裝 CLI 後，請執行設定精靈來設定開發環境：

```shell
wails3 setup
```

@note{type="caution" title="實驗性功能"}
設定精靈是新功能，目前主要在 Linux 上進行測試。如果遇到問題，請[回報問題](https://github.com/wailsapp/wails/issues/4904)，並改用下方的手動安裝相依套件步驟。

@end

設定精靈將會：

- 檢查平台相依套件並協助安裝
- 設定專案預設值（作者資訊、套件組合識別碼前置字串）
- 選擇性設定 Docker 以進行跨平台建置
- 設定程式碼簽署（如有需要）

如需更多詳細資訊，請參閱[設定指南](/getting-started/setup/)。

## 手動安裝相依套件

如果偏好手動安裝相依套件，或設定精靈無法在你的系統上運作，請依照上方的平台特定說明操作，然後執行：

```shell
wails3 doctor
```

這會檢查是否已安裝正確的相依套件，並告知你缺少哪些項目。

## 找不到`wails3`命令？

如果系統回報找不到`wails3`命令，請檢查下列項目：

- 請確認已正確依照上述<strong>Go安裝指南</strong>操作，且`go/bin`目錄位於`PATH`環境變數中。
- 關閉並重新開啟目前的終端機，使新的`PATH`變數生效。
