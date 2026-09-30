---
title: "CLI 參考資料"
description: "Wails CLI 指令的完整參考資料"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

## 概觀

Wails CLI（`wails3`）是用來建立、開發、建置、簽署、封裝及檢查 Wails 3應用程式的命令列進入點。大部分的建置協調工作會委派給各專案的 Taskfile（位於專案的`build/`下）——許多`wails3`指令只是呼叫特定工作項目的薄層包裝。

若要取得任何指令的最新說明，請執行：

```bash
wails3 --help
wails3 <command> --help
```

## 專案生命週期

| 指令 | 說明 |
| --- | --- |
| `wails3 init` | 從範本建立新專案。旗標：`-n`（專案名稱）、`-t`（範本，預設為`vanilla`）、`-p`（Go 套件名稱，預設為`main`）、`-d`（專案目錄，預設為`.`）、`-q`（安靜模式）、`-l`（列出範本）、`-mod`（Go 模組路徑）、`--git`（Git 儲存庫 URL）、`--skipgomodtidy`、`-s`（略過遠端範本警告）、`--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`。 |
| `wails3 dev` | 以開發模式執行應用程式，並啟用前端熱重新載入。旗標：`--config`（預設為`./build/config.yml`）、`--port`（Vite 開發伺服器連接埠）、`-s`（啟用 HTTPS）。 |
| `wails3 build` | 建置專案。這是 Taskfile `build`工作項目的薄層包裝。旗標：`--tags`（以`EXTRA_TAGS=`轉送）、`--obfuscated`（使用 Garble 建置；請參閱[混淆建置](/guides/build/obfuscation/)）、`--garbleargs`（在`build`子指令之前，將額外旗標轉送至`garble`）。 |
| `wails3 package` | 執行平台專用的`package` Taskfile 工作項目。 |
| `wails3 task [name]` | 執行任何 Taskfile 工作項目；未指定名稱時，`--list`會顯示所有已註冊的工作項目。 |
| `wails3 mcp` | 啟動專案 MCP 伺服器。由代理程式啟動的程序會自動使用 stdio；在互動式終端機中使用時，則會自動使用迴路 Streamable HTTP。 |
| `wails3 doctor` | 輸出環境診斷報告。 |
| `wails3 doctor-ng` | `doctor`的新版 TUI 變體。 |
| `wails3 version` | 輸出 CLI 版本。 |
| `wails3 releasenotes` | 輸出近期版本資訊。 |
| `wails3 docs` | 在瀏覽器中開啟文件網站。 |
| `wails3 sponsor` | 開啟贊助頁面。 |

## 產生

`wails3 generate <subcommand>`：

| 子指令 | 說明 |
| --- | --- |
| `generate bindings` | 產生從 Go 到前端的繫結。旗標：`-d`（輸出目錄）、`-models`、`-index`、`-ts`、`-i`（介面）、`-b`（套件組合）、`-names`（輸出`Call.ByName`）、`-noevents`、`-noindex`、`-dry`、`-silent`、`-v`、`-clean`（預設為`true`）、`-f`、`-obfuscated`（產生具有穩定繫結 ID 的`wails_obfuscated.gen.go`，供 Garble 建置使用；請參閱[混淆建置](/guides/build/obfuscation/)）、`-obfuscated-output`（產生檔案的目錄；預設為 main 套件目錄）。接受套件模式（例如`./...`）；若未提供，則改用目前目錄。 |
| `generate icons` | 將來源 PNG 轉換為各平台的圖示格式。旗標：`-input`、`-windowsfilename`、`-macfilename`、`-iconcomposerinput`、`-macassetdir`。 |
| `generate build-assets` | 根據`build/config.yml`產生`build/`目錄的內容（Taskfile 程式碼片段、NSIS 檔案、`Info.plist`、`.desktop`範本等）。 |
| `generate runtime` | 重新產生隨附並提供給 webview 的預先建置`/wails/runtime.js`。 |
| `generate syso` | 產生 Windows `.syso`資源檔案（圖示 + 資訊清單 + 版本資訊）。 |
| `generate webview2bootstrapper` | 產生適用於 Windows 的 WebView2 啟動安裝程式。 |
| `generate constants` | 根據 Go 事件型別產生 JS 事件名稱常數。 |
| `generate template` | 建立新專案範本的基本架構。 |
| `generate .desktop` | 產生 Linux `.desktop`檔案（供 AppImage/DEB/RPM 使用）。 |
| `generate appimage` | 產生 AppImage 建置目錄。 |

## 更新

`wails3 update <subcommand>`：

| 子指令 | 說明 |
| --- | --- |
| `update build-assets` | 從`build/config.yml`重新整理`build/`目錄（盡可能保留使用者的編輯）。 |
| `update cli` | 自行更新`wails3`二進位檔。 |

## 程式碼簽署與封裝

| 指令 | 說明 |
| --- | --- |
| `wails3 setup signing` | 互動式精靈，可為其在`build/`中偵測到的平台設定簽署。旗標：`--platform`（可重複指定；預設從建置目錄自動偵測）。 |
| `wails3 setup entitlements` | 用於設定 macOS 權利的互動式精靈。旗標：`--output`（路徑；預設為`build/darwin/entitlements.plist`）。 |
| `wails3 sign [GOOS=…]` | 包裝器，會執行目前作業系統專屬的`*:sign` Taskfile 工作（或透過`GOOS`指定的平台工作）。 |
| `wails3 tool sign` | 底層的直接簽署進入點。旗標：`--input`、`--output`、`--verbose`、`--certificate`、`--password`、`--thumbprint`、`--timestamp`、`--identity`、`--entitlements`、`--hardened-runtime`、`--notarize`、`--keychain-profile`、`--pgp-key`、`--pgp-password`、`--role`。 |

並<strong>沒有</strong>`wails3 signing`子命令——鑰匙圈憑證請直接使用`xcrun notarytool store-credentials`，PGP 金鑰則直接使用`gpg`（`wails3 setup signing`精靈會自動執行這兩項操作）。

## 工具

`wails3 tool <subcommand>`：

| 子命令 | 說明 |
| --- | --- |
| `tool checkport` | 檢查 TCP 連接埠是否開啟（適合用於等待 Vite）。 |
| `tool watcher` | 每當監看中的檔案變更時執行命令。 |
| `tool cp` | 跨平台複製檔案。 |
| `tool buildinfo` | 輸出二進位檔中嵌入的 Go 建置資訊。 |
| `tool package` | 從`build/linux/nfpm`建置 Linux 套件（`deb`、`rpm`、`archlinux`）。 |
| `tool version` | 遞增專案的語意化版本。 |
| `tool lipo` | 將多個 macOS 架構的二進位檔合併成通用二進位檔。 |
| `tool capabilities` | 探測系統是否提供 GTK3/GTK4／WebKit。 |
| `tool sign` | （請參閱[程式碼簽署與封裝](#heading-4)。） |

## 服務

`wails3 service <subcommand>`：

| 子命令 | 說明 |
| --- | --- |
| `service init` | 建立新的服務套件骨架。 |

## iOS

`wails3 ios <subcommand>`：

| 子命令 | 說明 |
| --- | --- |
| `ios overlay:gen` | 產生 iOS 橋接相容層的 Go overlay。 |
| `ios xcode:gen` | 在輸出目錄中產生 Xcode 專案。 |

## 建置輸出路徑

- 原生二進位檔會輸出至`bin/<APP_NAME>`（在 Windows 上則為`bin/<APP_NAME>.exe`）。不存在`build/bin/`。
- 封裝後的輸出（`.app`、`.dmg`、NSIS 安裝程式、MSIX、DEB/RPM/AppImage）同樣會輸出至`bin/`（或由相應 Taskfile 工作建立的各平台專屬子目錄中）。

## 全域旗標

| 旗標 | 適用範圍 | 說明 |
| --- | --- | --- |
| `--no-colour` | 所有命令 | 停用 CLI 輸出中的 ANSI 色彩。 |

---

<strong>有問題嗎？</strong>請在[Discord](https://discord.gg/JDdSxwjhGf)中提問，或查看[範例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
