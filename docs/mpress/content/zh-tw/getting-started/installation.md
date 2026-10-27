---
title: "安裝"
description: "安裝 Wails 並設定開發環境"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## 支援的平台

- Windows AMD64/ARM64
- macOS AMD64/ARM64。產生的建置任務以 macOS 12.0 或更新版本為目標；這是標準建置的部署目標，並不表示每一種自訂建置設定都支援較舊的系統。
- Ubuntu 24.04 AMD64/ARM64（其他 Linux 也可能可以運作！）

## 相依套件

安裝 Wails 前，必須先安裝一些通用相依套件。

@note{type="tip"}
安裝 Wails CLI 後，可以執行 `wails3 setup`，自動檢查這些相依套件並協助安裝。

@end

@tabs
[Go（至少為 1.26）]
Wails v3 會選用 Go 1.26 或更新版本的工具鏈。較舊的 Go 啟動程式可能會透過 `GOTOOLCHAIN=auto` 自動下載該工具鏈；否則請明確安裝。入門模組中較低版本的 Go 宣告不會覆寫 Wails 相依套件的需求。

請從 [Go 下載頁面](https://go.dev/dl/)下載 Go。

請務必遵循官方的 [Go 安裝說明](https://go.dev/doc/install)。還必須確認 `PATH` 環境變數包含 `~/go/bin` 目錄的路徑。重新啟動終端機，然後執行下列檢查：

- 確認 Go 已正確安裝：`go version`
- 確認 `~/go/bin` 位於 PATH 變數中
  - Mac / Linux：`echo $PATH | grep go/bin`
  - Windows：`$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm（選用）]
雖然 Wails 不要求安裝 npm，但大多數內附範本都需要使用 npm。

請從 [Node 下載頁面](https://nodejs.org/en/download/)下載最新的 Node 安裝程式。建議使用最新版本，因為我們通常會以該版本進行測試。

執行 `npm --version` 以進行驗證。

@note{type="info"}
如果偏好使用 npm 以外的套件管理員，可以自行選用，但需要更新專案的 Taskfile，讓它們改用該套件管理員。

@end

@end

## 平台專用相依套件

還需要安裝平台專用的相依套件：

@tabs{sync-key="platform"}
[Mac]
Wails 要求安裝 Xcode 命令列工具。可以執行以下命令進行安裝：

```sh
xcode-select --install
```

[Windows]
Wails 要求安裝 [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)。幾乎所有 Windows 系統都已安裝此執行階段。可以使用 `wails doctor` 命令進行檢查。

[Linux]
Linux 除了需要標準的 `gcc` 建置工具外，還需要 `gtk4` 和 `webkitgtk-6.0`。安裝後執行 <code>wails3 doctor</code>，即可查看相依套件的安裝方式。在 v3.1 之前，仍可透過 `-tags gtk3` 使用舊版 GTK3 / WebKit2GTK 4.1 技術堆疊（請參閱 [Linux 封裝－舊版 GTK3 支援](/guides/build/linux/#legacy-gtk3-support)）。如果目前使用的發行版或套件管理員不受支援，請在 Discord 上告知我們。

@end

## 安裝

若要使用 Go Modules 安裝 Wails CLI，請執行下列命令：

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

若要安裝最新的開發版本，請執行下列命令：

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

使用開發版本時，所有產生的專案都會使用 Go 的 [replace](https://go.dev/ref/mod#go-mod-file-replace) 指令，確保專案使用 Wails 的開發版本。

## 後續步驟

安裝 CLI 後，請執行設定精靈來設定開發環境：

```shell
wails3 setup
```

@note{type="caution" title="實驗性功能"}
設定精靈是新功能，目前主要在 Linux 上經過測試。如果遇到問題，請[回報問題](https://github.com/wailsapp/wails/issues/4904)，並改用下方的手動相依套件安裝步驟。

@end

設定精靈會：

- 檢查平台相依套件並協助安裝
- 設定專案預設值（作者資訊、套件組合識別碼前置字串）
- 選擇性設定 Docker，以進行跨平台建置
- 設定程式碼簽署（如有需要）

如需更多詳細資訊，請參閱[設定指南](/getting-started/setup/)。

## 手動安裝相依套件

如果偏好手動安裝相依套件，或設定精靈無法在系統上運作，請按照上方的平台專用說明操作，然後執行：

```shell
wails3 doctor
```

這會檢查是否已安裝正確的相依套件，並告知尚缺少哪些項目。

## 找不到 `wails3` 命令？

如果系統回報找不到 `wails3` 命令，請檢查下列項目：

- 確認已正確遵循上方的 **Go 安裝指南**，且 `go/bin` 目錄已加入 `PATH` 環境變數。
- 關閉後重新開啟目前的終端機，以載入新的 `PATH` 變數。
