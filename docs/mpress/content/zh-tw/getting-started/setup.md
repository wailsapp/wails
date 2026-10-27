---
title: "設定"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="實驗性功能"}
設定精靈是新功能，目前主要在 Linux 上測試。如果遇到問題，請[回報問題](https://github.com/wailsapp/wails/issues/4904)，並改為依照[手動安裝步驟](/getting-started/installation/#platform-specific-dependencies)操作。

@end

## 快速開始

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

精靈會在瀏覽器中開啟，引導你完成相依性檢查、專案預設值設定，以及選用的跨平台建置設定。

接著就可以開始建立第一個專案：

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 功能說明

- **檢查相依性** - 驗證 Go、npm 和平台工具
- **設定預設值** - 作者資訊、套件組合 ID 前綴、偏好的範本
- **跨平台建置** - 選用的 Docker 設定，可從任何主機進行建置
- **程式碼簽署** - 適用於 macOS、Windows 和 Linux 的選用設定

設定會儲存至`$XDG_CONFIG_HOME/wails/defaults.yaml`，若無法使用則改用`~/.config/wails/defaults.yaml`，並由`wails3 init`使用。這與專案的`build/config.yml`分開。

## 子命令

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## 遇到問題？

1. 執行`wails3 doctor`以診斷問題
2. 依照[手動安裝步驟](/getting-started/installation/#platform-specific-dependencies)操作
3. 附上`wails3 doctor`的輸出並[回報問題](https://github.com/wailsapp/wails/issues/4904)
