---
title: "MSIX 封裝"
description: "將您的 Wails v3 應用程式封裝為 MSIX 套件"
slug: "guides/build/msix"
sourcePath: "guides/build/msix.md"
---

MSIX 是現代化的 Windows 應用程式封裝格式。Wails 可在 Windows 建置過程中產生 MSIX 套件。

MSIX 封裝的操作說明請參閱[Windows 封裝](/guides/build/windows/#msix-package)指南。

## 直接使用 CLI 封裝

請在 Windows 上執行 MSIX 工具，CI 環境亦同。預設後端使用 `MakeAppx.exe`；簽署還需要 `signtool.exe`。兩者皆包含於 Windows SDK。安裝輔助工具會開啟 Microsoft Store，並視需要開啟 SDK 下載頁面；請先完成安裝，再進行封裝。

在 `build/config.yml` 中設定識別資訊，然後建置並封裝執行檔：

```yaml
info:
  companyName: "Example Corp"
  productName: "MyApp"
  productIdentifier: "com.example.myapp"
  description: "MyApp"
  version: "1.0.0"
```

```powershell
wails3 tool msix-install-tools
wails3 build GOOS=windows
wails3 tool msix --executable bin/myapp.exe --name myapp.exe
```

直接執行此命令會產生 `MyApp.msix`，檔案位於目前目錄。[Windows 封裝工作](/guides/build/windows/#msix-package)則會提供自己的輸出路徑。

## CLI 選項

| 選項 | 意義 |
| --- | --- |
| `--config` | 設定檔；預設為 `build/config.yml`。 |
| `--executable`, `--name` | 現有執行檔的路徑及其在套件內的檔名；兩者皆為必要項目。 |
| `--out` | 輸出檔案；預設為 `<ProductName>.msix`。 |
| `--arch` | 套件架構：`x64`（預設）、`x86`、`arm`、`arm64`、`x86a64` 或 `neutral`。也接受 Go 別名 `amd64` 與 `386`。請與執行檔的架構保持一致。 |
| `--publisher` | 發行者識別資訊；預設為 `CN=<companyName>`。 |
| `--cert`, `--cert-password` | 用於簽署的 PFX 憑證路徑與密碼。 |
| `--use-makeappx` | 使用 Windows SDK 的預設封裝工具。 |
| `--use-msix-tool` | 明確選用 `MsixPackagingTool.exe`，此工具必須位於 `PATH` 中。 |

## 簽署與 CI

若要在 Store 之外散發，請使用目標電腦信任的憑證簽署。憑證的 Subject 必須與 `--publisher` 完全相符。提供 `--cert` 後，MakeAppx 後端會呼叫 SignTool 並使用 SHA256。請參閱 [Microsoft 簽署指南](https://learn.microsoft.com/en-us/windows/msix/package/sign-msix-package-guide)。

此 Windows 工作流程步驟假設已安裝 Wails 與 SDK，且先前的步驟已將 PFX 檔案安全地放置於 `CERT_PATH` 指定的位置。只有路徑的祕密並不會自動上傳憑證：

```yaml
- name: MSIX
  if: runner.os == 'Windows'
  shell: pwsh
  run: |
    wails3 build GOOS=windows
    wails3 tool msix --executable bin/myapp.exe --name myapp.exe --publisher "$env:MSIX_PUBLISHER" --cert "$env:CERT_PATH" --cert-password "$env:CERT_PASSWORD"
  env:
    MSIX_PUBLISHER: ${{ vars.MSIX_PUBLISHER }}
    CERT_PATH: ${{ secrets.WINDOWS_CERT_PATH }}
    CERT_PASSWORD: ${{ secrets.WINDOWS_CERT_PASSWORD }}
```

## 檔案關聯與資源

在 `build/config.yml` 中新增不含開頭句點的副檔名；產生的資訊清單會加上句點：

```yaml
fileAssociations:
  - ext: myext
    name: MyApp Document
    description: MyApp Document
    iconName: fileicon
```

請依照[檔案關聯](/guides/file-associations/)的說明處理執行期間的檔案開啟操作。目前 MakeAppx 後端只會複製執行檔並產生透明的預留位置影像，不會匯入專案中的 `Assets/` 檔案，也不會轉換 `iconName` 圖示。若需要品牌影像或額外的 DLL，請使用自訂封裝工作流程。

| 產生的資源 | 尺寸（像素） |
| --- | --- |
| `Square150x150Logo.png` | 150×150 |
| `Square44x44Logo.png` | 44×44 |
| `Wide310x150Logo.png` | 310×150 |
| `StoreLogo.png` | 50×50 |
| `SplashScreen.png` | 620×300 |
| `FileIcon.png` | 44×44 |

只有設定了檔案關聯時，才會產生 `FileIcon.png`。這些檔案位於套件的 `Assets/` 目錄中。

## 提交至 Store 與疑難排解

在 [Partner Center 入口網站](https://partner.microsoft.com/dashboard) 中保留您的應用程式，並在準備提交時使用其中的套件識別資訊與發行者。Store 會在提交過程中簽署 MSIX 套件，因此透過此途徑散發時無須購買簽署憑證。請參閱 [Microsoft 套件需求](https://learn.microsoft.com/en-us/windows/apps/publish/publish-your-app/msix/app-package-requirements)。

若找不到 `MakeAppx.exe` 或 `signtool.exe`，請安裝或修復 Windows SDK。Wails 會搜尋 `PATH` 與標準 SDK 位置。簽署失敗時，請檢查憑證的 Subject、有效期限，以及目標電腦是否信任該憑證；請參閱 [MSIX 疑難排解](https://learn.microsoft.com/en-us/windows/msix/msix-troubleshooting-guide)。
