---
title: "從 v3 Alpha 升級"
description: "將現有 Wails v3 Alpha 專案升級至固定的 Beta 版本"
slug: "migration/alpha-to-beta"
sourcePath: "migration/alpha-to-beta.md"
---

本指南適用於現有的 v3 Alpha 專案。Wails v2 使用者請參閱 [v2 升級至 v3 指南](/migration/v2-to-v3/)。

## 升級之前

提交或備份您的專案。閱讀從目前 Alpha 版本到所選 Beta 版本之間的[變更記錄](/changelog/)：可能需要修改原始碼、API 或建置設定。請確認[桌面相容性政策](/status/)與平台需求。

以下命令使用已發布的 `v3.0.0-beta.23` 作為固定版本的範例，並非建議持續追蹤最新版本。如果選擇其他版本，請確認其 CLI、Go 模組與 npm 執行階段版本，並一併更新命令。在此範例中，npm 版本與移除 `v` 前綴的 Go 版本相同。

## 1. 更新 CLI

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.23
wails3 version
```

確認 `wails3 version` 顯示您安裝的版本。`PATH` 中較前面的舊執行檔可能會遮蔽新的 CLI。

## 2. 更新 Go 模組

在專案根目錄執行。檢查相依套件變更，不要大量升級無關的模組。

```sh
go get github.com/wailsapp/wails/v3@v3.0.0-beta.23
go mod tidy
```

## 3. 更新前端執行階段

對於使用 npm 與 `frontend` 目錄的專案：

```sh
cd frontend
npm install --save-exact @wailsio/runtime@3.0.0-beta.23
cd ..
```

保留鎖定檔並檢查其變更。如果前端使用其他套件管理工具或目錄，請調整此步驟，同時維持明確固定的執行階段版本。

## 4. 重新產生、建置與測試

在專案根目錄，從 Go 服務重新產生繫結並建置：

```sh
wails3 generate bindings
wails3 build
```

執行建置後的應用程式，並在您發布的每個受支援平台上測試工作流程。檢查並一併提交原始碼、產生的繫結、模組檔案與前端鎖定檔的變更。

## 升級失敗時

檢查 `PATH` 中的 CLI，使用 `go list -m github.com/wailsapp/wails/v3` 確認模組版本，並使用 `npm --prefix frontend ls @wailsio/runtime` 確認已安裝的執行階段。解決版本不符的問題後，重新產生繫結。不要假設每個 Alpha 版本都能在不修改程式碼的情況下升級。

若問題仍然存在，請[回報可重現的問題](https://github.com/wailsapp/wails/issues/new/choose)，附上新舊版本、完整錯誤訊息與 `wails3 doctor` 輸出。漏洞回報請遵循[安全性政策](https://github.com/wailsapp/wails/blob/master/SECURITY.md)。
