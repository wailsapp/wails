---
title: "YAML 建置"
description: "使用單一 wails.yaml 檔案設定並執行實驗性的 Wails v3 原生建置系統。"
sidebar:
  order: 2
  badge:
    text: "實驗性"
    variant: caution
---

@note{type="caution" title="實驗性預覽"}
第一個可接受的版本僅支援常值設定。這項實驗也會製作有界生命週期掛鉤的原型，以便在提議納入正式版本前評估其形式。自訂階段、YAML 運算式及具型別的 Wails 工具呼叫仍延後實作。現有專案會繼續使用 Taskfile，直到您明確啟用經審查的遷移為止。
@end

## 啟用 YAML 建置

在每個執行實驗性命令的終端機中設定 `WAILS_EXP_USE_WAKE`。若使用 Bash 或 Zsh：

```bash
export WAILS_EXP_USE_WAKE=1
```

若使用 PowerShell：

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

這個變數採用「存在即選用」的機制：即使值為空或設為 `false`，也會啟用此功能。取消設定該變數（在 Bash/Zsh 中使用 `unset WAILS_EXP_USE_WAKE`，或在 PowerShell 中使用 `Remove-Item Env:WAILS_EXP_USE_WAKE`）即可停用這項實驗。啟用 `wails.yaml` 後，命令會要求重新啟用 Wake，且絕不會退回使用 Taskfile。
若未設定該變數，YAML 建置／開發路由及 YAML 專案產生功能都會停用，`migrate`、`config`及 `eject` 命令也無法使用。
尚未遷移的專案會繼續透過獨立的內嵌 Task 執行階段運作。

## 設定專案

啟用選用設定後，請在專案根目錄加入 `wails.yaml`。該檔案的存在會選用原生建置系統，並讓 Wails 忽略專案中的所有 Taskfile。執行下列所有命令時，請維持選用設定。結構描述 modeline 可在支援 YAML 語言伺服器的編輯器中啟用自動完成及行內驗證。

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/wails.v3.json
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
frontend:
  directory: frontend
  install: [pnpm, install, --frozen-lockfile]
  build: [pnpm, run, build]
  dev: [pnpm, run, dev]
  output: frontend/dist
build:
  output: bin
  trim_path: true
  strip: true
profiles:
  release:
    targets:
      linux/amd64:
        formats: [deb, rpm, appimage]
```

執行預設的主機建置，或完整的具名發行要求：

```bash
wails3 build
wails3 build release
```

若要進行一次性建置，請選取以逗號分隔的目標及格式，不必新增設定檔：

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## 交叉編譯桌面應用程式

若要使用實驗性的純 Go 桌面後端，請在建置環境中設定 `CGO_ENABLED: "0"`。Wake 會為 macOS 和 Linux GUI 程式選用 `purego` 標籤，並直接在 Windows、macOS 及 Linux 主機上使用 Go。在此模式下，建置桌面可執行檔不需要 Docker、模擬器、C 編譯器或目標平台 SDK。應用程式的相依套件也必須支援 `CGO_ENABLED=0`。

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  output: bin
  environment:
    CGO_ENABLED: "0"
```

```bash
wails3 build --targets windows/amd64,windows/arm64,linux/amd64,linux/arm64,darwin/amd64,darwin/arm64
wails3 build --targets darwin/universal --formats dmg
```

第二個命令會合併兩種 macOS 架構，並使用 Go 建立未簽署的 DMG。它可在任何受支援的桌面主機上執行。程式碼簽署及 Apple 公證仍是各自獨立的 macOS 作業。

每個可執行檔都會獨立選取解析後的模式。已宣告的 `build.environment.CGO_ENABLED`、目標環境或二進位檔環境，會覆寫呼叫命令之終端機中的值。否則，明確選取的 `purego` 標籤、CLI／伺服器程式或 Windows 目標預設使用 `0`；macOS 和 Linux GUI 程式預設使用 `1`。`--plan --json` 會報告每個可執行檔解析後的編譯器、CGo 模式及標籤。交叉編譯會建立目標檔案；若要執行不相容的桌面目標，仍須使用適當的主機。

使用 `CGO_ENABLED=1` 時，Wake 會先使用相容的原生工具鏈或已設定的 `CC` 交叉編譯器；進行 Windows 交叉建置時，若 Zig 可用則接著使用 Zig；之後在受支援的情況下，才會使用現有的 Docker／Podman 備援方案。明確指定的 `CC` 必須提供目標編譯器、程式庫及 sysroot；僅有 Zig 並不會提供 Linux GTK/WebKit sysroot。封裝及簽署可能需要下文所述的其他工具。行動裝置目標需要 CGo 及其平台 SDK。

僅使用 Go 編譯不會封裝目標平台的系統程式庫。
Linux GUI 應用程式（包括 GUI 隨附可執行檔）仍需要其建置標籤所選取的
GTK/WebKitGTK 執行階段。套件相依項目或適當的部署環境必須提供該執行階段。僅含 CLI／伺服器的套件
預設不會加入這些 GUI 相依項目；仍可明確覆寫相依項目。

## 建置多個可執行檔

專案會定義主要程式。將 `project.package` 設為其 Go 套件（預設為 `.`），並將 `project.kind` 設為 `gui`、`cli` 或 `server`（預設為 `gui`）。其他程式則使用 `binaries` 下的具名項目：

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
binaries:
  notes-tool:
    package: cmd/notes-tool
    kind: cli
    tags: [tools]
    ldflags: [-X, main.Mode=release]
    environment:
      CGO_ENABLED: "0"
profiles:
  release:
    targets:
      linux/amd64: {}
```

每個程式都有各自的編譯及快取輸入。系統會追蹤匯入的專案套件及內嵌檔案；編輯某個獨立程式不會重新建置其他程式。桌面平台支援額外的可執行檔，且預設的應用程式套件組合及安裝程式會納入這些檔案。由使用者自行提供的替代套件檔案，其內容由使用者自行控制。

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

具名設定檔是一項完整要求；這些選取覆寫應搭配匿名建置使用。`--jobs` 會限制執行工作站數量。若只選取 CLI／伺服器程式，便會省略前端工作。

若要在不執行應用程式碼的情況下讀取版本，請同時設定 `project.version_file` 和 `project.version_symbol`。該檔案必須包含靜態 Go 字串宣告。之後即可省略常值 `project.version`。在專案或具名二進位檔上設定 `commit_symbol: main.Commit`，即可將儲存庫修訂版本注入該可執行檔的字串變數。要求注入提交資訊時，若缺少儲存庫中繼資料即會產生錯誤。

`upx_level: 1` 到 `9` 可選擇性壓縮發行版 Linux 及 Windows 可執行檔；設為零則停用。必須安裝 UPX。開發版及 macOS 建置不會壓縮。編譯或壓縮失敗時，會保留先前發佈的程式。

## 封裝並簽署以供散發

建立未簽署套件與簽署作業的主機／工具需求不同：

| 輸出 | 建立作業的主機與工具 | 簽署 |
| --- | --- | --- |
| Windows MSIX | Windows、macOS 或 Linux；使用 Go 寫入器，不需要 Windows SDK。 | Windows 主機及簽署工具。 |
| Windows NSIS 安裝程式 | Windows、macOS 或 Linux；需要主機原生的 `makensis`。 | Windows 主機及簽署工具。 |
| macOS app 與 DMG | Windows、macOS 或 Linux；使用 Go 組裝、通用二進位檔合併及可攜式 DMG 寫入器。 | 具備 Apple 工具及認證資訊的 macOS。 |
| Linux DEB、RPM、Arch 及 Void | Windows、macOS 或 Linux；使用原生 Go 套件寫入器。 | DEB／RPM 使用各自的外部簽署工具；Arch 使用下述 Go 簽署器。Void 輸出未經簽署。 |
| Linux AppImage | 相符的 Linux 主機及目標架構；使用外部 linuxdeploy 及 GTK 部署外掛程式。 | 沒有內建的 AppImage 簽署流程。 |

Windows MSIX 封存檔直接以 Go 建立；建立未簽署的套件
不需要 Windows SDK 或 Docker。Windows 簽署仍需要
Windows 主機及其簽署工具。NSIS 使用外部 NSIS 編譯器，並透過`packages.nsis.defines`接受常值定義。在 macOS 上，僅當 NSIS 編譯器子行程繼承的
地區設定不是 UTF-8 時，Wails 才會只為該子行程
提供 UTF-8 地區設定；終端機與應用程式環境不會變更。
當封裝的每個可執行檔都是 CLI 或伺服器時，內建 NSIS 安裝程式會省略 WebView2 啟動載入器。若主要或隨附可執行檔為 GUI，則會保留 WebView2。
自訂 NSIS 指令碼會保留啟動載入器，以便與其本身的執行階段安裝邏輯相容。

Linux 支援 AppImage、DEB、RPM、Arch Linux 及 Void 套件。AppImage 會保留外部 linuxdeploy 及 GTK 外掛程式的邊界。設定 `packages.appimage.update_info` 以要求產生 zsync 檔案，設定 `packages.appimage.checksum: true` 以要求產生 SHA256 檔案。附屬檔案會從套件快取還原，並與 AppImage 一同列在成品收據中。AppImage 封裝目前需要 Linux，且主機架構必須與目標架構相符。其 GTK 偵測無法找出純 Go GUI 動態載入的程式庫，因此此配接器尚不支援純 Go GUI 部署。該建置模式請改用其中一種 Linux 套件管理器格式。

Arch 套件會在原生套件中繼資料中保留預發行版本，例如 `1.2.3-beta.2` 會變成 `1.2.3beta.2-1`，且排序在最終版本之前。Void 相依套件可接受 `gtk4` 之類的純套件名稱，並正規化為 `gtk4>=0`；明確的版本比較及完整的 `name-version_revision` 值都會保留。請使用 Void 套件名稱，而非程式庫檔名。

若要建立 Arch 套件及分離式簽章，請使用已匯出的 ASCII Armor PGP 私密金鑰檔案，並可選擇指定密碼片語環境變數的名稱：

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  environment:
    CGO_ENABLED: "0"
linux:
  signing:
    certificate: signing/private-key.asc
    credential: ARCH_SIGNING_PASSPHRASE
profiles:
  release:
    targets:
      linux/amd64:
        formats: [archlinux]
        sign: true
```

請透過本機環境或 CI 祕密儲存區設定指定的環境變數，然後執行`wails3 build release`。若金鑰未加密，請省略`credential`。Go 簽署器可在所有支援的桌面主機上執行，且不需要 GPG。它會發布原始的`.pkg.tar.zst`及同層的`.pkg.tar.zst.sig`；兩者會分別記錄在成品收據中。金鑰內容不會暫存或儲存在建置快取中。簽章及其發布結果絕不會從快取重複使用。若簽署失敗，會保留先前發布的檔案組；若安裝任一檔案失敗，發布程序會復原這兩個檔案。這並不會讓並行讀取者所見的兩次檔案系統重新命名成為不可分割的原子操作。

macOS App 套件包含隨附程式。設定`darwin.microphone_usage_description`即可加入用途說明文字；若未提供自訂權利檔案，也會加入音訊輸入權利。自訂 plist 和權利來源檔案會保持不變。簽署與公證會在 macOS 上使用既有的 Apple 工具和憑證；公證需要真實的公證設定檔。簽署後的輸出會保留其平台副檔名，例如`notes.signed.app`和`notes.signed.dmg`。

## 開發控制項

執行`wails3 dev`期間，輸入命令後按 Enter：

| 命令 | 動作 |
| --- | --- |
| `r` | 強制重新建置並重新啟動，包括未變更的後端。 |
| `b` | 重新產生繫結並重新建置。 |
| `f` | 重新啟動前端。 |
| `c` | 清除終端機。 |
| `q` | 停止工作階段及其子處理程序。 |
| `h` | 顯示控制項。 |

建置失敗時，目前的應用程式會繼續執行。控制項需要互動式終端機；透過管線傳入輸入不會改變工作階段的行為。

## 執行前檢查工作

計畫檢查會解析執行作業所使用的同一個不可變圖。文字輸出供人員閱讀；JSON 則是具版本的機器可讀合約。

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## 不建置而驗證設定

請在 CI、編輯器工作中，或提交設定前使用專用的預檢命令。未指定設定檔時，它會驗證基礎資訊清單和每個具名設定檔；指定設定檔時，則只會檢查該完整請求：

```bash
wails3 config check
wails3 config check release
```

語法、結構描述、語意、路徑、目標、格式、簽署、規劃器，以及由設定所決定的主機功能失敗，皆會使用與 Build 和 Dev 相同的診斷轉譯器。診斷會指出欄位與來源範圍、印出來源行和插入號，並在有幫助時提供修正提示。診斷的來源摘錄會遮蔽環境值。僅由命令列選項造成的主機失敗仍是一般命令錯誤，不會被誤報為無效的 YAML。

## 使用專案自有的資產

資訊清單中的路徑會參照您擁有的檔案，例如圖示、平台資訊清單、套件範本和維護者指令碼。Wails 絕不會編輯這些來源檔案。工具使用前，Wails 會將它們複製到可拋棄的`.wails/`工作區。最終成品會放在`build.output`之下，並將經驗證的收據寫入`.wails/artifacts/receipt.json`。

## 部署至 Android

正式環境 Android 設定檔會產生 AAB。進行本機部署時，Wails 會在內部建置特定 ABI 的偵錯 APK、使用 ADB 安裝，然後啟動應用程式：

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

若未明確指定選取器，互動式終端機會列出已連線的裝置和已設定的 Android 虛擬裝置供您選擇。當有多個可用目標時，自動化流程必須提供`--device`或`--emulator`。

附加應用程式範圍的記錄資料流並監控目標，直到應用程式結束或您按下 Ctrl+C：

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

附加執行會跟隨應用程式重新啟動，並回報目標中斷連線或轉為離線的情況。模擬器預設會繼續執行，以便快速反覆測試。`--stop-emulator`必須明確指定，且需要`--logs`；附加工作階段結束時，它只會停止模擬器目標。若 Wails 啟動新的模擬器，但等候開機就緒失敗或遭取消，Wails 會自動清理該處理程序。

## 執行有界限的生命週期掛鉤

實驗性掛鉤介面會在每個已宣告的階段叫用一個專案自有的指令碼。掛鉤不是殼層命令：除非設定了`directory`，否則 Wails 會從專案根目錄直接執行該檔案。

在 Unix 主機上，指令碼必須可執行並以 shebang 開頭。Windows 主機會透過固定的直譯器配接器接受`.cmd`、`.bat`和`.ps1`指令碼；Wails 在任一平台上都不會剖析內嵌殼層命令。

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
  after_package:
    script: scripts/check-packages.sh
profiles:
  release:
    targets:
      linux/amd64: {}
```

支援的階段為`before_build`、`after_build`、`before_package`、`after_package`、`before_sign`和`after_sign`。建置掛鉤會依適用情況，對共用專案執行一次或對每個目標執行一次。套件與簽署掛鉤會在某個目標所要求的所有格式之外形成屏障，而屏障內的格式作業仍會並行執行。

掛鉤會收到一個由 Wails 管理的環境變數：`WAILS_HOOK_CONTEXT_FILE`。它指定一個受保護的 JSON 檔案，其中包含具版本的叫用內容：

```json
{
  "version": 1,
  "phase": "after_build",
  "command": "build",
  "scope": "target",
  "project_dir": "/work/notes",
  "working_directory": "/work/notes",
  "profile": "release",
  "target": { "os": "linux", "arch": "amd64" },
  "output": "/work/notes/bin/notes",
  "declared_outputs": ["/work/notes/generated/version.go"]
}
```

內容包含已解析的路徑，但絕不會複製繼承的處理程序環境或祕密。Wails 會在`.wails/hooks/context/`下，以僅限擁有者存取的檔案權限，以不可分割的方式建立該檔案。成功或取消後，該檔案會被移除。失敗時則會保留以供診斷，且錯誤訊息會包含其路徑。非零結束碼會使一般的 Plan 步驟失敗，取消作業則會終止掛鉤處理程序群組。

掛鉤預設會在每次叫用時執行。只有在指令碼完整宣告輸入與輸出時，才應選擇啟用快取：

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
    cache: true
    inputs: [version.txt]
    outputs: [generated/version.go]
profiles:
  release:
    targets:
      linux/amd64: {}
```

指令碼及其可執行模式是隱含輸入。多個輸出必須共用一個非根目錄，該目錄就是 Wails 記錄並還原的有界成品。輸出根目錄不得包含指令碼或已宣告的輸入，而快取的輸出樹狀結構只能包含一般檔案和目錄，絕不能包含符號連結或特殊檔案。一律執行的掛鉤會使其下游作業無法快取，因為無法安全地為未宣告的副作用建立指紋。可快取的掛鉤則會透過其宣告的成品識別資訊，保留下游快取能力。

## 移轉現有專案

建立並檢閱未啟用的草稿：

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

保守且不帶引數、會叫用專案自有指令碼的生命週期工作，會轉換為掛鉤。其他可到達的自訂 Taskfile 行為會回報為阻礙因素。檢閱`wails.migrated.yaml`後，請明確啟用它：

```bash
wails3 migrate --activate
```

啟用作業會以不可分割的方式建立`wails.yaml`。它不會刪除、重新命名或編輯 Taskfile、`build/config.yml`、指令碼或資產。無效的 YAML 會以 YAML 錯誤而失敗；絕不會退回舊版路由。

## 具體化所有預設值

若要檢查或複製完整的已解析設定，請寫入未啟用的快照：

```bash
wails3 eject
```

這會建立`wails.ejected.yaml`，且絕不會重寫使用中的設定。

## 相容性期間

`wails3 package`和`wails3 sign`是同一個規劃器已淘汰的別名；新的自動化流程應搭配設定檔或格式使用`wails3 build`。推出作業受版本管控：第一個版本會讓 YAML 成為新專案採用的路徑，下一個版本會針對舊版 Taskfile 路由發出警告，再下一個版本則會一併移除 Taskfile 執行功能與這些別名。每次轉換都要求該版本的原生主機驗收矩陣全部通過。

如需瞭解每個欄位、型別、預設值和適用規則，請參閱產生的[`wails.yaml`欄位參考](/reference/wails-yaml-fields/)。

[`wails3 run`設定參考](/reference/wails-yaml-run/)說明本機啟動預設值、平台覆寫，以及沒有資訊清單時的`go run .`備援機制。
