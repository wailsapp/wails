---
title: "設定"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="實驗性功能"}
設定精靈是新功能，目前主要在 Linux 上測試。如果遇到問題，請[回報問題](https://github.com/wailsapp/wails/issues/4904)，並改為按照[手動安裝步驟](/getting-started/installation/#platform-specific-dependencies)操作。

@end

## 快速開始

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

精靈會在瀏覽器中開啟，引導您檢查相依套件、設定全域專案預設值、選用的跨平台工具與簽署認證。這些偏好設定可供不同專案重複使用。

接著就可以開始建立第一個專案：

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 功能說明

- **檢查相依性** - 驗證 Go、npm 和平台工具
- **設定預設值** - 作者資訊、套件組合 ID 前置字串、偏好的範本
- **跨平台建置** - 選用的 Docker 設定，可從任何主機進行建置
- **程式碼簽署** - 適用於 macOS、Windows 和 Linux 的選用設定

設定會儲存至`$XDG_CONFIG_HOME/wails/defaults.yaml`；若無法使用，則改用`~/.config/wails/defaults.yaml`，並由`wails3 init`使用。這與專案的`build/config.yml`分開。

## 建立 Wake 專案

啟用 Wake 實驗性功能後，`wails3 init -ui`會收集專案專屬的建置選項，並在成功建立專案時將其儲存至`wails.yaml`。建立桌面專案時，請按照[Wake 入門教學](/experimental/wake/)選擇建置後端、目標與套件格式。`wails3 setup`仍負責全域偏好設定。

現有 Wake 專案是在`wails.yaml`中設定。請使用[YAML 建置指南](/experimental/yaml-builds/)變更其設定檔與封裝設定。現有 Taskfile 專案在使用 Wake 前，必須先完成[經審查的單向遷移](/experimental/wake/#migrate-an-existing-project)；設定程序不會轉換或執行其 Taskfile。

## 重複使用簽署偏好設定

已簽署的 Wake 建置若缺少簽署設定，會重複使用使用者層級預設值中的設定，包括專案外部的憑證檔案。明確指定的專案設定優先。私密金鑰、密碼和全域憑證路徑不會複製到`wails.yaml`。

對於全域設定的 Windows 和 Linux 金鑰，請透過現有的憑證儲存區或`WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`提供密碼。若使用加密的 Apple 身分，請將`darwin.signing.credential`設為包含其密碼的環境變數名稱。

## 子命令

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## 遇到問題？

1. 執行`wails3 doctor`以診斷問題
2. 按照[手動安裝步驟](/getting-started/installation/#platform-specific-dependencies)操作
3. 附上`wails3 doctor`的輸出並[回報問題](https://github.com/wailsapp/wails/issues/4904)
