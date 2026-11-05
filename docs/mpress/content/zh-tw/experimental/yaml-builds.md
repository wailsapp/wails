---
title: "YAML 建置"
description: "使用單一 wails.yaml 檔案設定並執行實驗性的 Wails v3 原生建置系統。"
sidebar:
  order: 2
  badge:
    text: "實驗性功能"
    variant: caution
---

@note{type="caution" title="實驗性預覽"}
首個可用版本僅支援常值設定。這項實驗也會建立有界生命週期掛鉤的原型，以便在提議納入正式版本前評估其形式。自訂階段、YAML 運算式和具型別的 Wails 工具呼叫仍暫緩實作。現有專案會繼續使用 Taskfile，直到您明確啟用經審查的遷移。
@end

第一次使用 Wake？請先參閱[入門教學](/experimental/wake/)。
該教學涵蓋此實驗的限制、所有專案與建置設定畫面、首次本機建置，以及經過審查的移轉。本頁將更詳細地說明各個建置、封裝與設定工作流程。

## 啟用 YAML 建置

在每個執行實驗性命令的終端機中設定`WAILS_EXP_USE_WAKE`。若使用 Bash 或 Zsh：

```bash
export WAILS_EXP_USE_WAKE=1
```

若使用 PowerShell：

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

此變數是依是否存在來選擇啟用：即使值為空或`false`，也會啟用。若要停用這項實驗，請取消設定此變數（Bash/Zsh 使用`unset WAILS_EXP_USE_WAKE`，PowerShell 使用`Remove-Item Env:WAILS_EXP_USE_WAKE`）。啟用`wails.yaml`後，命令會要求重新啟用 Wake，且絕不會回退至 Taskfile。若未設定此變數，YAML 建置／開發路由和 YAML 專案產生功能都會停用，且無法使用`migrate`、`config`和`eject`命令。尚未遷移的專案會繼續透過獨立的內嵌 Task 執行階段運作。

## 設定專案

若要在瀏覽器中建立及設定專案，請在啟用選擇加入功能的情況下執行`wails3 init -ui`，並依照[Wake 教學課程](/experimental/wake/)操作。專案選項儲存在`wails.yaml`中；若要變更現有專案，請編輯該檔案。`wails3 setup`用於設定共用偏好設定、工具及簽署認證。

啟用選擇加入功能後，請在專案根目錄新增`wails.yaml`。此檔案的存在會選用原生建置系統，並使 Wails 忽略該專案中的所有 Taskfile。執行下列所有命令時，請保持選擇加入功能已啟用。結構描述 modeline 可在支援 YAML 語言伺服器的編輯器中啟用自動完成和行內驗證。

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

執行預設的主機建置或完整的具名發行要求：

```bash
wails3 build
wails3 build release
```

若要執行一次性建置，可選取以逗號分隔的目標和格式，而不必新增設定檔：

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## 交叉編譯桌面應用程式

若要使用實驗性的純 Go 桌面後端，請在建置環境中設定`CGO_ENABLED: "0"`。Wake 會為 macOS 和 Linux GUI 程式選用`purego`標籤，並在 Windows、macOS 和 Linux 主機上直接使用 Go。在此模式下，建置桌面可執行檔不需要 Docker、模擬器、C 編譯器或目標平台 SDK。應用程式相依項目也必須支援`CGO_ENABLED=0`。

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

第二個命令會合併兩種 macOS 架構，並使用 Go 建立未簽署的 DMG。它可在任何受支援的桌面主機上執行。

每個可執行檔會各自選擇解析後的模式。已宣告的`build.environment.CGO_ENABLED`、目標環境或二進位檔環境會覆寫呼叫命令之終端機的值。否則，明確選取的`purego`標籤、CLI／伺服器程式或 Windows 目標預設使用`0`；macOS 和 Linux GUI 程式預設使用`1`。`--plan --json`會報告每個可執行檔解析後的編譯器、CGo 模式和標籤。交叉編譯會建立目標檔案；若要執行不相容的桌面目標，仍需使用適當的主機。

使用`CGO_ENABLED=1`時，Wake 會先使用相容的原生工具鏈或已設定的`CC`交叉編譯器；若進行 Windows 交叉建置且 Zig 可用，接著會使用 Zig；之後則在支援的平台上使用既有的 Docker/Podman 回退方案。明確指定的`CC`必須提供目標編譯器、程式庫和 sysroot；Zig 本身不會提供 Linux GTK/WebKit sysroot。封裝和簽署可能需要其他工具，如下文所述。行動裝置目標需要 CGo 及其平台 SDK。

僅使用 Go 編譯不會封裝目的地的系統程式庫。Linux GUI 應用程式（包括 GUI 隨附可執行檔）仍需要其建置標籤所選定的 GTK/WebKitGTK 執行階段。必須由套件相依項目或適當的部署環境提供該執行階段。預設情況下，僅含 CLI／伺服器的套件不會新增這些 GUI 相依項目；仍可明確覆寫相依項目。

## 建置多個可執行檔

專案會定義主要程式。將`project.package`設為其 Go 套件（預設為`.`），並將`project.kind`設為`gui`、`cli`或`server`（預設為`gui`）。其他程式使用`binaries`下的具名項目：

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

每個程式都有自己的編譯和快取輸入。系統會追蹤匯入的專案套件和內嵌檔案；編輯獨立程式不會重新建置其他程式。桌面平台支援其他可執行檔，並會將其納入預設的應用程式套件組合與安裝程式。由使用者擁有的替代套件檔案會自行控制其內容。

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

具名設定檔是一項完整要求；請在匿名建置中使用這些選取覆寫。`--jobs`會限制執行工作程序的數量。若只選取 CLI／伺服器程式，則會略過前端工作。

若要在不執行應用程式碼的情況下讀取版本，請同時設定`project.version_file`和`project.version_symbol`。該檔案必須包含靜態 Go 字串宣告。之後即可省略常值`project.version`。在專案或具名二進位檔上設定`commit_symbol: main.Commit`，可將儲存庫修訂版注入該可執行檔的字串變數。要求注入提交資訊時，若缺少儲存庫中繼資料，則會發生錯誤。

`upx_level: 1`到`9`可選擇性壓縮發行版 Linux 和 Windows 可執行檔；設為零會停用此功能。必須安裝 UPX。開發建置和 macOS 建置不會使用壓縮。編譯或壓縮失敗時，會保留先前發布的程式。

## 封裝並簽署以供發佈

### 使用預先建置的桌面可執行檔

若其他建置系統或 CI 作業已產生可執行檔，請使用`--skip-build`：

```bash
wails3 build --targets linux/amd64 --formats deb,rpm --skip-build
wails3 run --skip-build -- --example-mode
```

請將每個可執行檔放在`build.output`（預設為`bin`）中，主要可執行檔使用`project.binary_name`，配套可執行檔則使用各個`binaries`項目的名稱。Windows 檔名以`.exe`結尾。一次指定多個目標時，會從各目標的`<os>-<arch>`子目錄讀取檔案。對於 macOS GUI，若`bin/<name>`存在，Wake 會從該處讀取原始可執行檔；否則會從`bin/<name>.app/Contents/MacOS/<name>`讀取。`darwin/universal`輸入必須同時包含 amd64 和 arm64 的可執行檔切片。

Wake 會驗證檔案類型和目標架構、將可執行檔複製到其產生的工作區，並追蹤檔案內容以供套件快取使用。可執行檔輸入若缺失、不相容或為符號連結，作業會直接失敗，而不會重新建置或還原舊輸入。系統不會覆寫提供的可執行檔，也不會對其執行 chmod。任何可執行檔的變更都會使使用它的套件快取失效。

此模式會略過前端安裝／建置、繫結、Go 編譯、UPX，以及`before_build`／`after_build`掛鉤。它仍會產生套件資產、組裝 macOS 應用程式套件組合，並執行已設定的封裝／簽署掛鉤。套件和簽署工具仍須符合下表的限制；不需要 Go 和前端建置工具。套件中繼資料來自目前的資訊清單，包括已設定的`version_file`；系統不會從二進位檔推斷版本。請使用相符的建置和套件設定，包括 GUI 或 CLI 類型。

也支援`wails3 package --skip-build`。使用`wails3 run`時，啟動引數和環境仍會套用，且可執行檔必須能在主機上執行。此模式需要作用中的 YAML 專案，且僅支援正式環境的桌面目標；不適用於開發工作階段或行動裝置封裝。

建立未簽署套件與進行簽署的主機／工具需求不同：

IPA 承載內容和 macOS 公證 ZIP 皆以 Go 寫入，並保留套件路徑、執行權限及相對符號連結。這兩項作業都不需要外部`zip`或`ditto`執行檔。這些封存檔包含應用程式檔案，而非 Finder 中繼資料或資源分支；macOS 簽署本來就會移除這些中繼資料，因為它們會使已簽署的套件失效。

| 輸出 | 建立套件所需的主機和工具 | 簽署 |
| --- | --- | --- |
| Windows MSIX | Windows、macOS 或 Linux；使用 Go 寫入器，不需要 Windows SDK。 | 可在任何桌面主機上使用匯出的 RSA PFX 進行 Go 簽署；憑證存放區／提供者身分需要 Windows SignTool。 |
| Windows NSIS 安裝程式 | Windows、macOS 或 Linux；需要主機原生的`makensis`。 | 可在任何桌面主機上使用匯出的 RSA PFX 進行 Go PE 簽署；憑證存放區／提供者身分需要 Windows SignTool。 |
| macOS 應用程式和 DMG | Windows、macOS 或 Linux；使用 Go 組裝、合併通用二進位檔，並透過可攜式 DMG 寫入器建立。 | 使用匯出的 RSA PKCS#12身分進行 Go 簽署，或進行臨時簽署；鑰匙圈身分仍使用 macOS 工具。API 金鑰公證和票證裝訂以 Go 執行。 |
| Linux DEB、RPM、Arch 和 Void | Windows、macOS 或 Linux；使用原生 Go 套件寫入器。 | 匯出的 OpenPGP 金鑰使用原生 DEB／RPM／Arch 簽署器。舊式 DEB／RPM 金鑰 ID 需要外部工具。Void 輸出不會簽署。 |
| Linux AppImage | Windows、macOS 或 Linux；使用 Go 進行部署並寫入 SquashFS，另需適用於目標的 AppImage 執行階段。跨平台封裝 GUI 時，必須提供 Linux 執行階段根目錄樹。 | 沒有內建的 AppImage 簽章流程。 |

對於已簽署的 DMG，Wake 通常會先簽署應用程式的私有副本，再由 Go 寫入器壓縮映像，然後簽署映像。這可避免使用`hdiutil`掛載、展開並重新壓縮 DMG。每次提出要求時都會執行簽署；已簽署的套件不會從建置快取中重複使用。具有`after_package`或`before_sign`掛鉤的專案仍採用原本先封裝再簽署的流程，包括`hdiutil`，因此這些掛鉤對未簽署映像所做的編輯會予以保留。這條與掛鉤相容的映像編輯路徑需要 macOS 和`hdiutil`。

若要進行可攜式 Apple 簽署，請設定`darwin.signing.identity: "-"`以使用臨時簽章，或使用`darwin.signing.certificate`指定匯出的 RSA PKCS#12身分。`darwin.signing.credential`指定其密碼環境變數的名稱。憑證簽章使用共用的 Go PFX/CMS 實作；系統不會自動從鑰匙圈匯出私密金鑰。使用產生的測試身分時，精簡與通用執行檔、巢狀 app 套件、framework 和 DMG 均已通過 Apple 嚴格的原生簽章檢查，包括拒絕遭竄改的內容。

進行公證時，匯出的 Developer ID 身分還需要相容的 RFC 3161 `darwin.signing.timestamp_server`。共用驗證器接受其支援的 SHA-2時間戳記設定檔，並拒絕不受支援的舊式演算法。臨時簽章無法進行公證。原生憑證格式驗證與真正的 Developer ID 信任及 Apple 是否接受提交內容是彼此獨立的。

Go 裝訂器會驗證 Apple 票證簽章的真實性，確認它涵蓋所有架構後，再將其附加至私有副本。它使用觀察到的票證傳遞通訊協定，而 Apple 並未將此通訊協定發布為受支援的 API。它不能取代 Gatekeeper 目前的撤銷檢查。App 和已簽署的 DMG 輸出均已通過原生`stapler validate`及嚴格的簽章檢查。

Windows MSIX 封存檔直接以 Go 建立；建立未簽署的套件不需要 Windows SDK 或 Docker。匯出的 RSA PFX 身分可在任何桌面主機上直接以 Go 簽署 Windows 可執行檔與個別 MSIX 套件。將 `windows.signing.certificate` 設為專案相對路徑的 PFX，並將 `windows.signing.credential` 設為其密碼環境變數。支援 AES/PBKDF2 與舊式 3DES 匯出；RC2 匯出必須重新匯出，或改用原生 Windows SignTool 備援方案。憑證存放區／提供者身分與 MSI 簽署需要 Windows SignTool。預設時間戳記服務為 DigiCert；發布已簽署的輸出前會先驗證時間戳記。格式限制請參閱[簽署指南](/guides/build/signing/)。
NSIS 使用外部 NSIS 編譯器，並透過 `packages.nsis.defines` 接受常值定義。在 macOS 上，若 NSIS 編譯器繼承的地區設定不是 UTF-8，Wails 僅會向該編譯器子行程提供 UTF-8 地區設定；終端機與應用程式環境不會變更。
當所有封裝的可執行檔都是 CLI 或伺服器時，原廠 NSIS 安裝程式會省略 WebView2 啟動載入程式。若主要程式或附屬程式是 GUI，則會保留 WebView2。自訂 NSIS 指令碼會保留啟動載入程式，以相容於其本身的執行階段安裝邏輯。

Linux 支援 AppImage、DEB、RPM、Arch Linux 和 Void 套件。Wake 直接以 Go 組裝 AppImage，包括 GTK/WebKit 部署、SquashFS 壓縮及 zsync 產生作業。它不會叫用 linuxdeploy、GTK 外掛程式、appimagetool 或 mksquashfs。系統會為目標架構下載版本固定且經總和檢查碼驗證的 AppImage type-2執行階段；也可改用`packages.appimage.runtime_path`選取相對於專案的執行階段 ELF。

在 Linux 上，預設程式庫來源是主機的根檔案系統。若要為其他架構封裝 GUI，或從 Windows／macOS 封裝 GUI，請將`packages.appimage.runtime_root`設為相對於專案的目標 Linux 根目錄樹，其中須包含其程式庫、已編譯的結構描述及模組快取。這必須是包含`usr/`的完整根目錄配置，而非`usr/`目錄本身。請提供來自您預計支援之最舊 Linux 發行版的執行階段。Wake 會讀取 ELF 匯入項目及明確的 GTK/WebKit 設定檔，而不執行目標程式。它會封裝選定的相依閉包，但不包含 glibc 或動態載入器。執行階段更新後，程式庫和資源內容的雜湊會使套件快取失效。

支援純 Go GUI 建置。對於使用`gtk3`建置的 GUI 執行檔，Wake 會選取 GTK 3；其他情況則選取 GTK 4，包括 CLI／伺服器的 GUI 隨附程式。同一映像中的所有 GUI 執行檔都必須使用同一個 GTK 版本。對於預先建置的二進位檔，可透過`packages.appimage.gtk_version`明確選取`3`或`4`；匯入的 CGo 程式庫必須與該選擇一致。

設定`packages.appimage.update_info`以要求產生 zsync 檔案，並設定`packages.appimage.checksum: true`以要求產生 SHA256 檔案。附屬檔案會從套件快取還原，並與 AppImage 一同列在成品收據中。已使用所提供的 GTK/WebKit 目錄樹驗證 Windows 和 macOS 封裝；產生的純 Go GUI 映像已通過 Linux 上的執行階段檢查。在一次條件相符的 CGo GUI 建置中，原生路徑耗時8.77秒，linuxdeploy 則耗時26.19秒。原生 gzip 映像為133 MiB，舊 zstd 路徑則為114 MiB，因此這項量測涵蓋整個封裝路徑，而不只是壓縮。

Arch 套件會在原生套件中繼資料中保留預發行版本，例如`1.2.3-beta.2`會變成`1.2.3beta.2-1`，排序在最終版本之前。Void 相依項目可接受`gtk4`之類的純套件名稱，並將其正規化為`gtk4>=0`；明確的版本比較和完整的`name-version_revision`值則會保留。請使用 Void 套件名稱，而非程式庫檔名。

DEB、RPM 和 Arch 簽署可在任何受支援的桌面主機上使用匯出的 ASCII 裝甲格式或二進位 OpenPGP 私密金鑰，無須 GnuPG、`dpkg-sig`或`rpmsign`。將`linux.signing.certificate`設為金鑰檔案，並將`credential`設為其複雜密碼環境變數的名稱。DEB 和 RPM 會內嵌簽章；Arch 則發布分離式簽章。只有簽署步驟會讀取私密金鑰，而且不會將其複製到套件資產或可重複使用的快取成品中。

包含裸 GnuPG 金鑰識別碼的現有 DEB／RPM 設定會繼續使用其外部簽署器。若要移除此相依性，請匯出私密金鑰，並改為設定其路徑。請使用如`./signing/release.asc`的路徑；若金鑰檔案不存在、無法讀取或無效，作業會失敗，而不會退回使用金鑰圈。DEB 簽章會保留`dpkg-sig`套件簽章格式；這不能取代對 apt 儲存庫的 release 中繼資料進行簽署。

若要建立 Arch 套件和分離式簽章，請使用匯出的 ASCII 裝甲格式 PGP 私密金鑰檔案，並可選擇指定複雜密碼環境變數的名稱：

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

透過本機環境或 CI 祕密儲存區設定指定名稱的環境變數，然後執行`wails3 build release`。若金鑰未加密，請省略`credential`。Go 簽章程式可在所有受支援的桌面主機上執行，且不需要 GPG。它會發布原始`.pkg.tar.zst`及同層的`.pkg.tar.zst.sig`；兩者會分別記錄於成品收據中。金鑰內容不會暫存或儲存在建置快取中。簽章及其發布絕不會從快取重複使用。簽章失敗時會保留先前發布的一對檔案；若安裝其中任一檔案失敗，發布作業會回復這兩個檔案。這不會讓兩次檔案系統重新命名對並行讀取者呈現原子性。

macOS app 套件包含隨附程式。設定`darwin.microphone_usage_description`以新增用途說明文字，並在未提供自訂權利檔案時新增音訊輸入權利。自訂 plist 和權利來源不會變更。匯出的 RSA PKCS#12身分會在任何桌面主機上使用 Go 簽署；`identity: "-"`會選取可攜式臨時簽署。鑰匙圈身分在 macOS 上仍使用`codesign`。票證裝訂以 Go 執行。提交公證時，可使用現有的 macOS 公證鑰匙圈設定檔，或 App Store Connect 團隊 API 認證：

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
darwin:
  signing:
    certificate: signing/developer-id.p12
    credential: APPLE_CERT_PASSWORD
    timestamp_server: https://timestamp.digicert.com
  notarization:
    api_key: signing/AuthKey.p8
    key_id: ABCDEFGHIJ
    issuer_id: 00000000-0000-0000-0000-000000000000
```

API 模式會以 Go 提交、上傳及輪詢，不使用`notarytool`。請設定全部三個欄位並省略`credential`；請勿將私密金鑰納入版本控制。只有在簽署時才會讀取金鑰，而且不會將其複製到產生的資產或快取中。此模式支援團隊 API 金鑰，以及最大5 GiB 的封存檔；尚未實作多部分上傳和個人 API 金鑰。提交遭拒時會保留提交 ID，並將 Apple 的診斷記錄儲存至私有暫存檔案；錯誤訊息會回報該檔案的路徑。已取消或遭拒的提交不會取代先前發布的輸出。接受提交後，Wake 會擷取並驗證 Apple 的票證，再以 Go 將其裝訂。搭配匯出的身分和團隊 API 認證時，簽署、提交及裝訂可在 Linux、Windows 或 macOS 上執行。鑰匙圈設定檔仍需要 macOS 和`notarytool`。

可攜式憑證簽署支援 AES/PBKDF2 或舊式3DES PKCS#12檔案中的 RSA 身分。以 RC2 加密的匯出內容必須使用受支援的密碼重新匯出。已公證的可攜式簽章需要明確的 SHA-2 RFC 3161時間戳記服務；Go 驗證器會拒絕舊式 SHA-1時間戳記簽章。票證服務遵循觀察到的 Apple 通訊協定，並固定信任 Apple Root G3；Gatekeeper 仍是判定散發信任與撤銷的權威。

已簽署的輸出會保留其平台副檔名，例如
`notes.signed.app` 和 `notes.signed.dmg`。

## 開發控制項

執行`wails3 dev`期間，輸入命令後按 Enter：

| 命令 | 動作 |
| --- | --- |
| `r` | 強制重新建置並重新啟動，包括未變更的後端。 |
| `b` | 重新產生繫結並重新建置。 |
| `f` | 重新啟動前端。 |
| `c` | 清除終端機畫面。 |
| `q` | 停止工作階段及其子程序。 |
| `h` | 顯示控制項。 |

建置失敗時，目前的應用程式會繼續執行。控制項需要互動式終端機；以管線傳入的輸入不會改變工作階段行為。

## 執行前檢查工作

檢查計畫時，會解析執行所使用的同一個不可變圖。文字供人員閱讀；JSON 是具有版本的機器可讀合約。

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## 不建置即可驗證設定

請在 CI、編輯器工作中，或提交設定前使用專用的預檢命令。未指定設定檔時，它會驗證基礎資訊清單及每個具名設定檔；指定設定檔時，則只檢查該完整要求：

```bash
wails3 config check
wails3 config check release
```

語法、結構描述、語意、路徑、目標、格式、簽章、規劃器，以及由設定決定的主機能力失敗，均使用與 Build 和 Dev 相同的診斷呈現器。診斷會指出欄位和來源範圍、印出原始碼行與插入號，並在有幫助時提供修正提示。診斷的來源摘錄會遮蔽環境值。若主機失敗僅由命令列選項造成，仍會作為一般命令錯誤處理，而不會被誤報為無效的 YAML。

## 使用專案自有資產

資訊清單中的路徑指向您擁有的檔案，例如圖示、平台資訊清單、套件範本和維護者指令碼。Wails 絕不會編輯這些來源。在工具使用這些檔案之前，Wails 會將其複製到可捨棄的`.wails/`工作區。最終成品會置於`build.output`下，並將已驗證的收據寫入`.wails/artifacts/receipt.json`。

## 部署至 Android

正式環境 Android 設定檔會產生 AAB。進行本機部署時，Wails 會在內部建置 ABI 特定的偵錯 APK、使用 ADB 安裝，然後啟動應用程式：

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

未明確指定選擇器時，互動式終端機會列出已連線的裝置和已設定的 Android 虛擬裝置供您選擇。若自動化環境中有多個可用目標，必須提供`--device`或`--emulator`。

附加應用程式範圍的記錄串流，並監控目標，直到應用程式結束或您按下 Ctrl+C：

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

附加執行會跟隨應用程式重新啟動，並回報目標中斷連線或轉為離線。模擬器預設會保持執行，以便快速反覆運算。`--stop-emulator`必須明確指定，且需要`--logs`；附加工作階段結束時，它只會停止模擬器目標。若 Wails 啟動新的模擬器，但等候開機就緒失敗或遭取消，則會自動清理該程序。

## 執行有界生命週期掛鉤

實驗性的掛鉤介面會在每個已宣告的階段叫用一個專案自有指令碼。掛鉤不是殼層命令：Wails 會直接執行該檔案，且除非設定`directory`，否則會從專案根目錄執行。

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

支援的階段為`before_build`、`after_build`、`before_package`、`after_package`、`before_sign`和`after_sign`。建置掛鉤會視情況針對共用專案執行一次，或針對每個目標各執行一次。封裝和簽章掛鉤會在針對某個目標要求的所有格式外形成屏障，而屏障內的各項格式作業仍會並行執行。

掛鉤會收到一個由 Wails 擁有的環境變數：`WAILS_HOOK_CONTEXT_FILE`。它指定一個受保護的 JSON 檔案，其中包含具有版本的叫用內容：

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

內容包含解析後的路徑，但絕不會複製繼承的程序環境或祕密。Wails 會在`.wails/hooks/context/`下，以僅擁有者可存取的檔案權限，透過原子操作建立該檔案。成功或取消後會移除該檔案。若失敗，則會保留以供診斷，且其路徑會包含在錯誤中。非零結束狀態會使一般 Plan 步驟失敗，而取消作業會終止掛鉤程序群組。

掛鉤預設會在每次叫用時執行。只有當指令碼完整宣告輸入和輸出時，才選擇使用快取：

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

指令碼及其可執行模式是隱含輸入。多個輸出必須共用一個非根目錄，該目錄即為 Wails 記錄和還原的有界成品。輸出根目錄不得包含指令碼或已宣告的輸入，而快取的輸出樹只能包含一般檔案和目錄，絕不能包含符號連結或特殊檔案。每次皆執行的掛鉤會使其下游作業無法快取，因為無法安全地為未宣告的副作用建立指紋。可快取的掛鉤則會透過其已宣告的成品識別資訊，保留下游快取能力。

## 移轉現有專案

建立並檢閱尚未啟用的草稿：

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

保守且不帶引數、會叫用專案自有指令碼的生命週期工作，會轉換為掛鉤。其他可觸及的自訂 Taskfile 行為會回報為阻礙。檢閱
`wails.migrated.yaml`後，請明確啟用它：

```bash
wails3 migrate --activate
```

啟用時會以不可分割的方式建立`wails.yaml`。此操作不會刪除、重新命名或編輯 Taskfile、`build/config.yml`、指令碼或資產。無效的 YAML 會以 YAML 錯誤告終，絕不會退回舊版路徑。

## 具體化所有預設值

若要檢查或複製完整且解析完成的設定，請寫入未啟用的快照：

```bash
wails3 eject
```

這會建立`wails.ejected.yaml`，絕不會重寫作用中的設定。

## 相容性期間

`wails3 package`和`wails3 sign`是同一規劃器已棄用的別名；新的自動化應搭配設定檔或格式使用`wails3 build`。推出時程受版本控制：第一個版本會讓 YAML 成為新專案的使用路徑，下一個版本會針對舊版 Taskfile 路由發出警告，再下一個版本則會移除 Taskfile 執行功能及這些別名。每次轉換都要求該版本的原生主機驗收矩陣全部通過。

如需瞭解每個欄位、型別、預設值及適用性規則，請參閱產生的[`wails.yaml`欄位參考](/reference/wails-yaml-fields/)。

[`wails3 run`設定參考](/reference/wails-yaml-run/)說明本機啟動預設值、平台覆寫，以及沒有資訊清單時的`go run .`後援機制。
