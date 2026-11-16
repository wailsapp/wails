---
title: "YAML 執行設定"
description: "wails3 run 的本機執行預設值、平台覆寫與啟動行為。"
---

## 可用性

本參考文件說明實驗性的 run 擴充功能，不保證已發布的 CLI 會提供此功能。

`wails3 run`會使用本機資訊清單的預設值建置並啟動應用程式。
由資訊清單驅動的執行需要透過
[`WAILS_EXP_USE_WAKE=1`選擇啟用](/experimental/yaml-builds#enable-yaml-builds)。
若找不到資訊清單，此命令會在目前目錄中執行`go run .`，
且不需要選擇啟用。資訊清單若無效或無法讀取，將視為錯誤；
絕不會因此改用備援方式。若資訊清單存在但 YAML 已停用，
命令會說明如何啟用 YAML。

產生的[`wails.yaml`欄位參考文件](/reference/wails-yaml-fields)列出
欄位型別與預設值。在發布版本宣稱支援原生行動平台與 macOS App 套件之前，
必須先在適當的主機上驗證其可接受性。

## 命令契約

`wails3 run`會尋找最近的`wails.yaml`、解析單一目標的設定、
透過現有的快取管線進行建置，然後啟動結果。
探索程序會從工作目錄開始向上搜尋，並在`go.mod`邊界停止前檢查資訊清單。
資訊清單所在目錄即為專案根目錄，也是桌面應用程式的工作目錄。

預設目標是主機的作業系統與架構。行動平台目標必須明確選取。
交叉編譯不代表桌面可執行檔能在主機上執行；不相容的桌面目標會在建置前失敗。
這也包括使用`CGO_ENABLED=0`建置的可執行檔：僅使用 Go 成功完成交叉建置，
並不會提供模擬器或變更本機啟動契約。
請使用[`wails3 build`進行交叉編譯](/experimental/yaml-builds/#cross-compile-desktop-applications)，
並在目標系統上執行或安裝產生的檔案。
Android 請使用`--target android/arm64 --device <adb-serial>`或`--target android/arm64 --emulator <avd-name>`。
所選架構必須與裝置或模擬器相符。iOS 模擬器請使用`--target ios/arm64 --device <udid>`；若為實體裝置，另加`--destination device`。規劃及執行 iOS 目標需要 macOS 與 Xcode。

使用 YAML 時，建置會採用嵌入式資產與正式環境建置路徑，並自動提供
`production`標籤。`wails3 dev`仍負責即時重新載入；
`dev.tags`不會提供內容給`wails3 run`。

```shell
wails3 run
wails3 run --tags extra_tag
wails3 run -- --example-mode other
wails3 run --plan
wails3 run --skip-build
```

`--plan`會解析並顯示所選目標、實際採用的建置標籤、
可執行或可安裝的成品、引數及工作目錄，而不進行建置、安裝或啟動。
計畫會列出環境覆寫的鍵名，但省略其值。

## CLI 選項

| 選項 | 意義 |
| --- | --- |
| `--tags <tags>` | 其他以逗號分隔的 Go 建置標籤。 |
| `--skip-build` | 從 YAML 輸出目錄匯入現有的桌面可執行檔，然後不經編譯直接啟動。 |
| `--target <os>/<arch>` | 選取一個目標；預設為主機的作業系統與架構。 |
| `--plan` | 顯示解析後的設定，不進行建置或啟動。允許規劃跨平台桌面目標。 |
| `--device <id>` | Android adb 序號或 iOS 模擬器／裝置識別碼。 |
| `--emulator <name>` | 要啟動或重複使用的 Android 虛擬裝置；不得與`--device`併用。 |
| `--destination <kind>` | iOS 目的地：`simulator`（預設）或`device`。 |
| `-- <args...>` | 取代已設定的應用程式引數。單獨使用`--`會將其清除。 |
| `--appargs <string>` | 相容拼法：使用 POSIX 風格的引號規則剖析單一字串。不得重複使用，也不得與`--`之後的引數併用。 |

行動平台選取器需要資訊清單以及相符的行動平台目標。

`--skip-build`需要有效的 YAML 專案和相容的桌面可執行檔。此模式會略過前端工作、繫結、編譯、UPX 和建置掛鉤；已設定的啟動引數和環境仍然適用。匯入的可執行檔經過驗證後，會複製到 Wake 產生的工作區，而不會修改其來源檔案。macOS GUI 應用程式會使用這些輸入和目前的套件資產組裝。請參閱[預先建置輸入的路徑與需求](/experimental/yaml-builds/#use-prebuilt-desktop-executables)。
建置標籤無法變更預先建置的可執行檔，因此`--tags`和已設定的建置標籤在此模式下不會生效。

## 欄位

選用的`run`對應表提供共用啟動預設值。`targets["<os>"]`或`targets["<os>/<arch>"]`內選用的`run`
對應表提供平台或架構專用的預設值。
無論對應表中的順序為何，架構專用設定的優先順序都高於整個作業系統共用的設定。

| 欄位 | 型別 | 預設值 | 意義 |
| --- | --- | --- | --- |
| `run.tags` | `list(string)` | `[]` | 僅供 run 建置使用的其他 Go 標籤。 |
| `run.args` | `list(string)` | `[]` | 傳遞給應用程式的引數。 |
| `run.environment` | `map(string)` | `{}` | 執行階段環境覆寫，與編譯器環境分開。 |
| `target["<os>"].run.tags` | `list(string)` | `[]` | 所選平台的其他 run 標籤。 |
| `target["<os>"].run.args` | `list(string)` | 繼承 | 若明確設定，便會取代共用引數；`[]`會將其清除。 |
| `target["<os>"].run.environment` | `map(string)` | `{}` | 依鍵覆寫共用的執行階段環境值。 |

引數會以引數向量傳遞，不經 shell 展開。依照現有 YAML 僅允許常值的契約，
環境值均為常值字串。

## 平台範例

此片段是一般專案設定的補充。`example_private_api`
是範例實際所需私有 API 標籤的預留位置。

```yaml
build:
  tags: [devtools]
run:
  args: [--example-mode]
targets:
  darwin:
    tags: [example_private_api]
    run:
      environment:
        EXAMPLE_WINDOW_MODE: native
  linux:
    run:
      environment:
        GDK_BACKEND: wayland
  windows:
    run:
      args: [--example-mode, --windows-option]
```

Linux 環境設定僅適用於此範例，且需要 Wayland 工作階段；
它不是 Wails 的預設值。私有 API 標籤僅在 macOS 上選取。
現有的`target.tags`適用於正式環境建置，包括 run 建置。
`target.run.tags`會將另一個標籤限制為僅供 run 建置使用。

架構專用覆寫可以進一步調整整個作業系統共用的預設值：

```yaml
targets:
  linux/arm64:
    run:
      tags: [example_arm64]
      args: []
```

與前述 Linux 對應表搭配使用時，這會保留`GDK_BACKEND`、加入架構專用標籤，
並清除共用的應用程式引數。

## 解析規則

| 設定 | 解析順序（優先順序由低至高） |
| --- | --- |
| 標籤 | 合併`build.tags`、所選`target.tags`、自動加入的`production`、`run.tags`、所選`target.run.tags`及 CLI `--tags`；每個標籤僅保留第一次出現的位置。 |
| 執行階段環境 | 繼承啟動命令的 shell 環境，疊加`run.environment`，再疊加所選`target.run.environment`。空值會設定為空字串。 |
| 應用程式引數 | 先採用共用`run.args`，再由明確設定的目標引數取代，最後由`--`之後明確提供的引數取代。單獨使用`--`會清除已設定的引數。 |

在所選目標內，會先解析整個作業系統共用的設定，再解析架構專用設定：
標籤會累加、環境值會依鍵合併，而明確提供的引數則會取代較不具體的清單。
run 建置也會針對對應的行動平台目標提供`android`或`ios`。

繫結產生程序與 Go 編譯會收到相同的解析後標籤。建置快取鍵包含這些標籤。
若只變更應用程式引數或執行階段環境，不會使編譯快取失效；下次啟動時會使用新值。

平台對應表用於自訂平台；其存在並不表示不支援其他平台。
`project.supported_platforms`會限制建置與`run`接受的平台；省略此清單或使用空清單，則允許所有平台。例如，在`project`內設定`supported_platforms: [darwin]`，會使僅適用於 macOS 的範例在 Linux 上提早失敗。

## 平台啟動行為

| 平台 | 啟動行為 |
| --- | --- |
| Linux 與 Windows | 啟動可執行檔，連接標準輸入／輸出，處理終止作業，並傳遞應用程式結束狀態。 |
| macOS | 組裝 app bundle，並執行其中的`Contents/MacOS`二進位檔，同時連接標準輸入／輸出。向 Launch Services 註冊該 bundle，以支援 URL scheme 和檔案關聯。簽署與原生生命週期驗收需要在 macOS 上進行驗證。 |
| Android | 建置 APK、選取裝置或模擬器，然後安裝並啟動應用程式。分別回報建置、安裝與啟動失敗。 |
| iOS | 建置 app bundle、選取模擬器或裝置、滿足簽署與佈建需求，然後安裝並啟動。裝置版建置需要適當的憑證和 Apple 工具鏈。 |

Android 會在部署前拒絕應用程式引數和執行階段環境覆寫。若同時支援 Android，僅適用於桌面的共用設定應置於桌面目標對應中。Android 執行時會附接應用程式記錄，並在附接結束時停止應用程式。iOS 會將引數傳給所選的啟動器。模擬器透過`SIMCTL_CHILD_`變數接收環境覆寫；實體裝置則透過`devicectl --environment-variables`接收 JSON 對應。行動平台的結束狀態描述的是啟動器或記錄監控的結果，而非可跨平台使用的應用程式結束碼。

## 僅使用 Go 的範例與後援機制

若僅使用 Go 的範例採用既有或嵌入式資產，會停用前端階段：

```yaml
frontend:
  disabled: true
```

這會略過前端相依套件安裝、繫結產生及前端建置。需要產生繫結類別的範例，則可在`frontend.bindings`下設定`interfaces: false`。

無 manifest 時的後援機制會執行`go run [-tags <tags>] . [application arguments]`，且不會加入`production`或`devtools`。其結束狀態行為遵循 Go 本身的規則。`--`之後的引數會傳給應用程式，包括`--help`之類的旗標。`--plan`顯示的執行階段值僅限環境變數鍵名，不會列印值。後援計畫會顯示 Go 命令和工作目錄。

## 開發環境的應用程式引數

使用有效的 YAML 專案時，`wails3 dev`接受的應用程式引數獨立於正式環境的`run`預設值。`dev.args`提供共用預設值；`targets["<os>"]`或`targets["<os>/<arch>"]`中巢狀的`dev`對應會針對該目標取代這些預設值。無論對應中的順序為何，架構設定的優先順序都高於整個作業系統適用的設定。明確的`args: []`會清除繼承的引數。

```yaml
dev:
  args: [--config-path, config/testing.yaml]
targets:
  windows:
    dev:
      args: [--config-path, config/testing-windows.yaml]
```

`run.args`和`run.tags`不會影響開發環境的啟動。既有的`dev.tags`建置原則維持不變。應用程式會以 manifest 目錄為基準解讀相對應用程式路徑。

```shell
wails3 dev
wails3 dev -- --config-path "profile with spaces.yaml"
wails3 dev --appargs='--config-path "profile with spaces.yaml"'
wails3 dev --
wails3 dev --plan
```

`--`之後的引數會取代已設定的引數，即使清單為空亦然。該分隔符號之後的`--help`屬於應用程式。`--appargs=""`也會清除預設值。相容性選項接受一個加上引號的字串；若將它與`--`、位置式應用程式引數或第二個`--appargs`一起使用，便會發生錯誤。系統會在啟動工作階段前拒絕格式錯誤的引號。

相容性剖析器在每個主機上都使用相同的 POSIX 風格引號規則。它不會執行 shell 命令或展開變數；叫用它的 shell 可能在將文字傳給 Wails 前就已完成展開。`--`之後的引數已是個別的值，不會再進行額外的詞元化。這是 Windows 路徑的建議形式，例如在 PowerShell 中：

```powershell
wails3 dev -- --config-path 'C:\Users\Example User\testing.yaml'
```

工作階段會在原始碼重新建置後保留已解析的引數向量。若只變更`dev.args`或所選目標的引數，後端會使用快取的可執行檔重新啟動，而不會重新啟動前端。CLI 覆寫在重新載入 manifest 後仍然有效。若取代失敗，會保留或還原先前的可執行檔及其先前的引數。應用程式引數不是編譯器旗標，也不會納入編譯快取鍵。`dev --plan`會在不啟動應用程式的情況下，連同有限建置計畫顯示實際生效的引數。

Linux、Windows 和 macOS 會將引數向量直接傳給可執行檔。iOS 會透過模擬器或裝置啟動器傳遞該向量。Android 沒有受支援的對等引數通道，因此會在探索裝置或建置前拒絕非空的開發環境應用程式引數。僅適用於桌面的預設值可置於桌面目標對應中，或在`targets.android.dev`下使用`args: []`清除。

這些開發選項需要有效的`wails.yaml`，且必須透過 YAML 明確啟用。它們不會改寫舊版開發路徑中的 Taskfile 或監看器 shell 命令。
