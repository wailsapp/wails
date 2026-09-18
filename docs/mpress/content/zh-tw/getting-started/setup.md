---
title: "設定"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="實驗性功能"}
設定精靈是新功能，目前主要在 Linux 上進行測試。如果遇到問題，請[回報](https://github.com/wailsapp/wails/issues/4904)，並改為依照[手動安裝步驟](/getting-started/installation/#platform-specific-dependencies)操作。

@end

## 快速開始

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard
wails3 setup
```

精靈會在瀏覽器中開啟，引導您檢查相依套件、設定專案預設值，以及選擇性設定跨平台建置環境。

接著便可開始建立第一個專案：

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## 功能說明

- **檢查相依套件**——驗證 Go、npm 和平台工具
- **設定預設值**——作者資訊、套件組合 ID 前置詞、偏好的範本
- **跨平台建置**——選擇性設定 Docker，以便從任何主機進行建置
- **程式碼簽署**——適用於 macOS、Windows 和 Linux 的選擇性設定

設定會儲存至`~/.config/wails/config.yaml`，並供`wails3 init`使用。

## 子命令

```bash
wails3 setup signing      # Configure code signing
wails3 setup entitlements # Configure macOS entitlements
```

## 遇到問題？

1. 執行`wails3 doctor`以診斷問題
2. 依照[手動安裝步驟](/getting-started/installation/#platform-specific-dependencies)操作
3. 附上`wails3 doctor`的輸出並[回報問題](https://github.com/wailsapp/wails/issues/4904)
