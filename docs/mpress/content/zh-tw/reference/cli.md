---
title: "CLI 指令概覽"
description: "依用途尋找 Wails CLI 指令：專案、建置、簽署、工具與服務"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

從這裡開始依用途尋找指令。如需特定指令的旗標和使用範例，請參閱[CLI 旗標與範例](/guides/cli/)。

## 概覽

Wails CLI（`wails3`）是建立、開發、建置、簽署、封裝及檢查 Wails 3應用程式的命令列進入點。大部分建置協調工作會委派給各專案的 Taskfile（位於專案中的`build/`之下）——許多`wails3`指令只是呼叫特定任務的輕量包裝。

若要取得任何指令的最新說明，請執行：

```bash
wails3 --help
wails3 <command> --help
```

## 專案生命週期

| 指令 | 說明 |
| --- | --- |
| `wails3 init` | 從範本建立新專案。旗標：`-n`（專案名稱）、`-t`（範本，預設為`vanilla`）、`-p`（Go 套件名稱，預設為`main`）、`-d`（專案目錄，預設為`.`）、`-q`（安靜模式）、`-l`（列出範本）、`-mod`（Go 模組路徑）、`--git`（Git 儲存庫 URL）、`--skipgomodtidy`、`-s`（略過遠端範本警告）、`--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`。 |
| `wails3 dev` | 以開發模式執行應用程式，並啟用前端熱重新載入。旗標：`--config`（預設為`./build/config.yml`）、`--port`（Vite 開發連接埠）、`-s`（啟用 HTTPS）。 |
| `wails3 build` | 建置專案。這是 Taskfile `build`任務的輕量包裝。旗標：`--tags`（以`EXTRA_TAGS=`轉送）、`--obfuscated`（使用 Garble 建置；請參閱[混淆建置](/guides/build/obfuscation/)）、`--garbleargs`（在`build`子指令之前轉送給`garble`的額外旗標）。 |
| `wails3 package` | 執行平台專屬的`package` Taskfile 任務。 |
| `wails3 task [name]` | 執行任意 Taskfile 任務；若未指定名稱，`--list`會顯示所有已註冊的任務。 |
| `wails3 mcp` | 啟動專案 MCP 伺服器。由代理程式啟動的處理程序會自動使用標準輸入輸出，互動式終端機則會使用回送介面的 Streamable HTTP。 |
| `wails3 doctor` | 輸出環境診斷報告。 |
| `wails3 doctor-ng` | `doctor`的較新版 TUI 變體。 |
| `wails3 version` | 輸出 CLI 版本。 |
| `wails3 releasenotes` | 輸出近期版本資訊。 |
| `wails3 docs` | 在瀏覽器中開啟文件網站。 |
| `wails3 sponsor` | 開啟贊助頁面。 |

## 產生

`wails3 generate <subcommand>`:

| 子指令 | 說明 |
| --- | --- |
| `generate bindings` | 產生從 Go 到前端的繫結。旗標：`-d`（輸出目錄）、`-models`、`-index`、`-ts`、`-i`（介面）、`-b`（套件組合）、`-names`（發出`Call.ByName`）、`-noevents`、`-noindex`、`-dry`、`-silent`、`-v`、`-clean`（預設為`true`）、`-f`、`-obfuscated`（產生具有穩定繫結 ID 的`wails_obfuscated.gen.go`，供 Garble 建置使用；請參閱[混淆建置](/guides/build/obfuscation/)）、`-obfuscated-output`（產生檔案的目錄；預設為 main 套件目錄）。接受套件模式（例如`./...`）；若未提供，則改用目前目錄。 |
| `generate icons` | 將來源 PNG 轉換為各平台的圖示格式。旗標：`-input`、`-windowsfilename`、`-macfilename`、`-iconcomposerinput`、`-macassetdir`。 |
| `generate build-assets` | 從`build/config.yml`產生`build/`目錄的內容（Taskfile 片段、NSIS 檔案、`Info.plist`、`.desktop`範本等）。 |
| `generate runtime` | 重新產生提供給 WebView 的預先建置`/wails/runtime.js`。 |
| `generate syso` | 產生 Windows `.syso`資源檔案（圖示 + 資訊清單 + 版本資訊）。 |
| `generate webview2bootstrapper` | 產生適用於 Windows 的 WebView2 啟動安裝程式。 |
| `generate constants` | 從 Go 事件型別產生 JS 事件名稱常數。 |
| `generate template` | 建立新專案範本的基本架構。 |
| `generate .desktop` | 產生 Linux `.desktop`檔案（供 AppImage/DEB/RPM 使用）。 |
| `generate appimage` | 產生 AppImage 建置目錄。 |

## 更新

`wails3 update <subcommand>`:

| 子指令 | 說明 |
| --- | --- |
| `update build-assets` | 從`build/config.yml`重新整理`build/`目錄（盡可能保留使用者的編輯內容）。 |
| `update cli` | 自行更新`wails3`二進位檔。 |

## 程式碼簽署與封裝

| 指令 | 說明 |
| --- | --- |
| `wails3 signing` | 互動式精靈，為其在`build/`中偵測到的平台設定簽署。旗標：`--platform`（可重複指定；預設從建置目錄自動偵測）。 |
| `wails3 entitlements` | macOS 權利設定的互動式精靈。旗標：`--output`（路徑；預設為`build/darwin/entitlements.plist`）。 |
| `wails3 sign [GOOS=…]` | 包裝程式，針對目前作業系統（或透過`GOOS`指定的作業系統）執行平台專屬的`*:sign` Taskfile 任務。 |
| `wails3 tool sign` | 低階直接簽署進入點。旗標：`--input`、`--output`、`--verbose`、`--certificate`、`--password`、`--thumbprint`、`--timestamp`、`--identity`、`--entitlements`、`--hardened-runtime`、`--notarize`、`--keychain-profile`、`--pgp-key`、`--pgp-password`、`--role`。 |

在此版本中，CLI 剖析器會在頂層公開`signing`和`entitlements`。請查看`wails3 signing --help`和`wails3 entitlements --help`；`wails3 setup signing`不會選取簽署精靈。若要直接管理認證，請使用`xcrun notarytool store-credentials`或`gpg`。

## 工具

`wails3 tool <subcommand>`:

| 子命令 | 說明 |
| --- | --- |
| `tool checkport` | 檢查 TCP 連接埠是否開啟（適合用來等待 Vite）。 |
| `tool watcher` | 每當監看中的檔案變更時執行命令。 |
| `tool cp` | 跨平台複製檔案。 |
| `tool buildinfo` | 輸出二進位檔中嵌入的 Go 建置資訊。 |
| `tool package` | 從`build/linux/nfpm`建置 Linux 套件（`deb`、`rpm`、`archlinux`）。 |
| `tool version` | 遞增專案的語意化版本。 |
| `tool lipo` | 將多個 macOS 架構的二進位檔合併為通用二進位檔。 |
| `tool capabilities` | 探測系統是否提供 GTK3/GTK4／WebKit。 |
| `tool sign` | （請參閱[程式碼簽署與封裝](#heading-4)。） |

## 服務

`wails3 service <subcommand>`:

| 子命令 | 說明 |
| --- | --- |
| `service init` | 建立新服務套件的基本架構。 |

## iOS

`wails3 ios <subcommand>`:

| 子命令 | 說明 |
| --- | --- |
| `ios overlay:gen` | 為 iOS 橋接相容層產生 Go 疊加層。 |
| `ios xcode:gen` | 在輸出目錄中產生 Xcode 專案。 |

## 建置輸出路徑

- 原生二進位檔會輸出至`bin/<APP_NAME>`（在 Windows 上則為`bin/<APP_NAME>.exe`）。不存在`build/bin/`。
- 封裝後的輸出（`.app`、`.dmg`、NSIS 安裝程式、MSIX、DEB/RPM/AppImage）同樣會輸出至`bin/`（或由相關 Taskfile 任務建立的平台專用子目錄中）。

## 共用旗標

| 旗標 | 適用於 | 說明 |
| --- | --- | --- |
| `-nocolour` | 在`--help`中公開此旗標的命令（例如`init`和`update cli`） | 停用 CLI 輸出中的 ANSI 色彩。並非所有命令都接受此旗標，其中包括`version`。 |

---

<strong>有問題嗎？</strong>請到[Discord](https://discord.gg/JDdSxwjhGf)提問，或查看[範例](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)。
