---
title: "建置與封裝流程"
description: "說明執行`wails3 build`時的內部運作方式、如何產生跨平台二進位檔，以及如何為各作業系統產生安裝程式。"
slug: "contributing/build-packaging"
sourcePath: "contributing/build-packaging.md"
---

`wails3 build`刻意設計得<strong>很精簡</strong>：它是 Taskfile 包裝器，會將額外的建置標籤轉送至主專案的`build`工作。實際的繁重工作位於專案本身的`build/Taskfile.yml`（由`wails3 init`產生）、`internal/commands/build-assets.go`（管理烘焙階段的資產），以及`internal/packager`（Linux nfpm 封裝）和`internal/commands/appimage.go`、`internal/commands/msix.go`、`internal/commands/dmg/dmg.go`、`internal/commands/dot_desktop.go`（各平台的安裝程式）。

本頁涵蓋：

1. 實際的 CLI 進入點
2. 由 Taskfile 驅動的建置流程
3. 資產烘焙與建置資訊注入
4. 各平台的封裝後端
5. 自訂流程
6. 疑難排解

---

## 1. 實際的 CLI 進入點

```
wails3 build       → internal/commands.Build       (in task_wrapper.go)
wails3 package     → internal/commands.Package     (in task_wrapper.go)
wails3 generate build-assets → GenerateBuildAssets (in build-assets.go)
wails3 update build-assets   → UpdateBuildAssets   (in build-assets.go)
wails3 tool buildinfo        → BuildInfoOptions    (in tool_buildinfo.go)
wails3 tool package          → internal/packager   (nfpm wrapper)
wails3 generate .desktop     → in dot_desktop.go
```

`internal/commands/task_wrapper.go`：

```go
func Build(buildFlags *flags.Build, otherArgs []string) error {
    // forwards --tags / EXTRA_TAGS, then defers to a Taskfile target
    return wrapTask("build", otherArgs)
}
```

`flags.Build`只提供<strong>一個</strong>旗標——`--tags`（轉送為`EXTRA_TAGS=`）。`wails3 build`上<strong>沒有</strong>`-platform`、`-o`、`-skipbindings`、`-skip-package`、`-package`、`-ldflags`、`-verbose`、`-debug`、`-devbuild`、`-icon`或`-clean`旗標。交叉編譯、輸出路徑、圖示等項目是在<strong>`Taskfile.yml`</strong>、<strong>`build/config.yml`</strong>以及輔助的`wails3 generate icons`／`wails3 generate build-assets`命令中設定。

`build/build.json`<strong>不</strong>屬於 v3——其設定由`Taskfile.yml`加上`build/config.yml`構成。

---

## 2. 由 Taskfile 驅動的建置流程

剛初始化的專案隨附一個`build/Taskfile.yml`，其命名空間大致如下：

| 命名空間 | 工作（節選） |
| --- | --- |
| `darwin:` | `build`、`build:universal`、`package`、`run`、`dev` |
| `windows:` | `build`、`package`、`run`、`dev` |
| `linux:` | `build`、`package`、`run`、`dev` |
| `common:` | `update:build-assets`、`generate:icons`、`generate:syso` |

`wails3 build`預設會叫用主機作業系統的`build`命名空間；專案的 Taskfile 接著會使用主機專用旗標，透過 shell 執行`go build`。若要為其他作業系統建置，必須直接執行其工作（例如`wails3 task darwin:build:universal`），而不是將旗標傳給`wails3 build`。

預設輸出目錄為<strong>`bin/<APP_NAME>`</strong>（沒有`build/bin/`前綴）。

---

## 3. 烘焙階段的資產與建置資訊

| 關注項目 | 檔案 |
| --- | --- |
| 建置資產的產生／更新 | `internal/commands/build-assets.go` |
| 建置資訊輸出工具（CLI：`wails3 tool buildinfo`） | `internal/commands/tool_buildinfo.go`——輸出資訊；它<strong>不是</strong>`ldflags`注入器 |
| 正式環境存根 | `internal/assetserver/build_production.go`——`//go:build production` |
| 前端打包資產 | 透過應用程式自身套件中的`//go:embed`嵌入（例如放在`main.go`旁） |
| Windows 資源（`.syso`） | `internal/commands/syso.go`——產生`rsrc_windows_<arch>.syso` |
| Windows MSIX | `internal/commands/msix.go` + `internal/commands/webview2/` |
| macOS DMG 輸入項目 | `internal/commands/dmg/` |
| Linux `.desktop` | `internal/commands/dot_desktop.go` |

CLI 不會自動為應用程式烘焙`bundled_assetserver.go`——`internal/assetserver/bundled_assetserver.go`是<strong>手寫的</strong>，並包裝`bundledassets/`下嵌入的 JS 執行階段。

---

## 4. 封裝後端

### Linux

Linux 封裝由<strong>nfpm</strong>驅動（不是`fpm`）：

- `internal/packager/packager.go`會包裝`github.com/goreleaser/nfpm/v2`，並提供`CreatePackageFromConfig(pkgType, configPath, output)`／`CreatePackageFromConfigWriter(...)`。
- 產生的專案會在`internal/commands/`下隨附 nfpm 樣式的`myapp.DEB`、`myapp.RPM`和`myapp.ARCHLINUX`設定（由`wails3 tool package`使用）。
- AppImage 的產生邏輯位於`internal/commands/appimage.go`，它會呼叫`linuxdeploy` + `linuxdeploy-plugin-gtk`（外掛程式隨附於`internal/commands/linuxdeploy-plugin-gtk.sh`）。

`wails3 build`上<strong>沒有</strong>`-package deb`／`rpm`旗標。請使用`wails3 tool package`或平台專用的 Taskfile 目標。

### macOS

- 專案 Taskfile 中的`darwin:package`會產生`.app`套件組合。
- DMG 資源位於 `internal/commands/dmg/`；`darwin:package`完成後，專案可以使用 `hdiutil` 將應用程式套件封裝成 DMG（較新範本中的 Taskfile 包含 `dmg` 輔助工具）。
- CFBundle 識別碼、版本和著作權資訊，來自執行 `wails3 init` 時的 `-product*` 旗標以及 `build/config.yml`。

### Windows

- Windows 封裝的目標格式是 **MSIX**（而非 WiX/MSI）。完整工作流程請參閱 `internal/commands/msix.go` 和 `internal/commands/webview2/`。
- **沒有** `internal/commands/packager.go`，也<strong>沒有</strong> `internal/commands/windows_resources/` 目錄。
- 可選的執行檔程式碼簽署透過 `wails3 tool sign`（Authenticode）執行，詳見 `internal/commands/sign.go`。

---

## 5. 自訂管線

| 需求 | 做法 |
| --- | --- |
| 額外的建置標籤 | `wails3 build --tags myFeature,otherTag` |
| Linter／建置前步驟 | 在 `build/Taskfile.yml` 中新增工作，並讓作業系統專用的 `build` 工作依賴該工作 |
| 交叉編譯 | 執行相應的作業系統工作（例如 `wails3 task linux:build`）；沒有 `-platform` 旗標 |
| 略過封裝 | 只執行 `build` 工作；`package` 是獨立的工作 |
| 自訂封裝工具 | 將設定檔放在 `internal/commands/myapp.*` 下，並使用 `-config <file>` 呼叫 `wails3 tool package` |
| 移除符號 | 編輯 `darwin:/windows:/linux:` 的 `build` 工作，將 `-ldflags "-s -w"` 直接傳給 `go build`；`wails3 build` 本身沒有 `-ldflags` 旗標 |

所有 Taskfile 目標都會遵循 Wails 發布的環境變數（`APP_NAME`、`WAILS_VITE_PORT`、`FRONTEND_DEVSERVER_URL`……），因此自訂工作可以依賴這些變數。

---

## 6. 疑難排解

| 症狀 | 可能原因 | 修正方式 |
| --- | --- | --- |
| **`ld: framework not found WebKit`（mac）** | 缺少 Xcode CLI 工具 | `xcode-select --install` |
| **正式版建置中出現空白視窗** | 前端建置失敗或 SPA 路由問題 | 確認 `frontend/dist/index.html` 存在，且資源處理常式會回退至該檔案 |
| **缺少 MSIX 封裝工具** | 尚未安裝 `WebView2` SDK／MSIX 工具 | 執行 `wails3 task install:msix:tools` |
| **`linuxdeploy` 找不到** | PATH 中缺少外掛程式 | 安裝 `linuxdeploy`，並透過 CLI 的自動安裝步驟執行 `internal/commands/linuxdeploy-plugin-gtk.sh` |

`wails3 build` 沒有 `-verbose` 旗標。設定 `TASK_X_VERBOSE=1`（Taskfile），或直接檢查工作目標以查看正在執行的命令。

---

## 7. 主要原始碼對照表

| 相關功能 | 檔案 |
| --- | --- |
| 建置包裝器 | `internal/commands/task_wrapper.go`（`Build`、`Package`、`SignWrapper`、`wrapTask`） |
| 建置資源產生 | `internal/commands/build-assets.go`（`GenerateBuildAssets`、`UpdateBuildAssets`） |
| 建置資訊輸出器 | `internal/commands/tool_buildinfo.go` |
| AppImage 建置器 | `internal/commands/appimage.go` |
| Linux 封裝（nfpm） | `internal/packager/packager.go`、`internal/commands/myapp.{DEB,RPM,ARCHLINUX}` |
| Windows MSIX | `internal/commands/msix.go`、`internal/commands/webview2/` |
| Windows 資源產生器 | `internal/commands/syso.go` |
| macOS DMG 資源 | `internal/commands/dmg/` |
| `.desktop` 產生器 | `internal/commands/dot_desktop.go` |
| 版本常數 | `internal/version/version.go` |

追查建置失敗時，請隨時參考此表。

---

現在，從<strong>原始碼</strong>到<strong>安裝程式</strong>，你已掌握完整全貌。簡而言之，`wails3 build` 本身只是一層薄包裝；幾乎所有自訂都在專案的 `Taskfile.yml`／`build/config.yml` 中進行，或透過明確的 `wails3 generate …`／`wails3 tool …` 子命令完成。祝發布順利！
