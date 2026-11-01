---
title: "設定"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="實驗性功能"}
設定精靈是新功能，目前主要在 Linux 上進行測試。如果遇到問題，請[回報問題](https://github.com/wailsapp/wails/issues/4904)，並改為依照[手動安裝步驟](/getting-started/installation/#platform-specific-dependencies)操作。

@end

## 快速開始

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

精靈會在瀏覽器中開啟，引導你檢查相依項目、設定專案預設值，以及選擇性設定跨平台建置。上述指示使用預設設定流程；實驗性的 Wake 專案流程則說明如下。

接著即可建立第一個專案：

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 功能說明

- **檢查相依項目** — 驗證 Go、npm 和平台工具
- **設定預設值** — 作者資訊、套件組合 ID 前綴、偏好的範本
- **跨平台建置** — 選擇性設定 Docker，以便從任何主機進行建置
- **程式碼簽署** — 選擇性設定 macOS、Windows 和 Linux 的簽署

設定會儲存至`$XDG_CONFIG_HOME/wails/defaults.yaml`，若無法使用則改用`~/.config/wails/defaults.yaml`，並由`wails3 init`使用。此設定與專案的`build/config.yml`分開。

## 設定 Wake 專案

在專案目錄中執行設定，並[啟用 Wake 實驗性功能](/experimental/yaml-builds/#enable-yaml-builds)：

```bash
export WAILS_EXP_USE_WAKE=1
wails3 setup
```

在 PowerShell 中：

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
wails3 setup
```

這會開啟以專案為主的流程：

1. **專案** — 若不存在原生資訊清單，請輸入專案中繼資料。
2. **應用程式** — 選擇 GUI、CLI 或伺服器，以及建置後端，或選用已建置的可執行檔。純 Go 是明確的實驗性選項。
3. **目標** — 選取作業系統和架構。iOS 和 Android 仍為實驗性功能；請在 YAML 中設定其簽署設定。
4. **套件** — 選擇套件格式，以及是否簽署套件。
5. **工具** — 檢視所選建置與封裝路徑的需求。
6. **簽署** — 若已選擇簽署，請使用現有的簽署控制項進行設定。
7. **儲存** — 寫入專案設定並顯示要執行的命令。

工具需求取決於你的選擇。純 Go 桌面編譯不需要 Docker 或 C 編譯器。CGo 跨平台建置需要相容的編譯器和目標程式庫，或受支援的容器工具鏈。NSIS 仍需要`makensis`；跨主機建置 GUI AppImage 時，需要提供 Linux 執行階段目錄樹。匯出的簽署身分可使用 Wails 的 Go 簽署器，而原生鑰匙圈或憑證存放區中的身分則繼續使用平台工具。請參閱[封裝與簽署需求](/experimental/yaml-builds/#package-and-sign-for-distribution)。

選用已建置的可執行檔會略過 Go 編譯、前端安裝與建置、繫結和 UPX。仍需準備所選的封裝／簽署工具及有效的輸入二進位檔。請參閱[預先建置的輸入路徑](/experimental/yaml-builds/#use-prebuilt-desktop-executables)。此選項適用於儲存後顯示的命令，並不是永久性的 YAML 建置設定。

「儲存」會寫入`wails.yaml`、更新應用程式種類、前端啟用狀態與所選的後端設定，並取代專用的`profiles.setup`目標／套件／簽署選項。其他設定檔、前端命令、隨附的可執行檔、掛鉤及不相關的專案設定都會保留。後端設定適用於整個專案，因此其他繼承這些設定的設定檔也會使用所選的後端。新的專案中繼資料會保留在精靈中，直到按下「儲存」。系統會回報缺少的工具，讓你能在建置前安裝；儲存有效的設定不會安裝工具或執行建置。

簽署建置會從使用者層級的預設值沿用缺少的簽署設定，包括專案外部的憑證檔案。明確的專案設定優先。設定程序不會將私密金鑰、密碼或全域認證路徑複製到`wails.yaml`。對於全域設定的 Windows 和 Linux 金鑰，請透過現有的認證存放區或`WAILS_WINDOWS_CERT_PASSWORD`／`WAILS_PGP_PASSWORD`提供密碼。若 Apple 身分已加密，請將`darwin.signing.credential`設為存有其密碼的環境變數名稱。

儲存的設定檔會使用`wails3 build setup`執行；選用已建置的可執行檔時，則會加入`--skip-build`。

執行顯示的命令時，請保持啟用 Wake。若未設定該旗標，作用中的`wails.yaml`絕不會改用 Taskfile。設定程序不會轉換現有的 Taskfile 或舊版中繼資料：使用 Wake 設定這些專案前，請先依照[經審查的移轉工作流程](/experimental/yaml-builds/#migrate-an-existing-project)操作。

## 子命令

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## 遇到問題？

1. 執行`wails3 doctor`以診斷問題
2. 依照[手動安裝步驟](/getting-started/installation/#platform-specific-dependencies)操作
3. 請附上`wails3 doctor`的輸出並[回報問題](https://github.com/wailsapp/wails/issues/4904)
