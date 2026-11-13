---
title: "程式碼簽署"
description: "在所有平台上簽署 Wails 應用程式的指南"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## 簽署應用程式的程式碼

本指南說明如何在 macOS、Windows 和 Linux 上簽署 Wails 應用程式。Wails v3提供內建的 CLI 工具，用於程式碼簽署、公證及 PGP 金鑰管理。

- **macOS** — 簽署並公證 macOS 應用程式
- **Windows** — 簽署 Windows 可執行檔與套件
- **Linux** — 使用 PGP 金鑰簽署 DEB 與 RPM 套件

## 跨平台簽署矩陣

此矩陣顯示可從各來源平台簽署的格式：

| 目標格式 | 從 Windows | 從 macOS | 從 Linux |
| --- | :---: | :---: | :---: |
| Windows EXE/DLL/MSIX/APPX、匯出的 PFX | ✅ | ✅ | ✅ |
| Windows MSI 或憑證存放區／提供者身分 | ✅ | ❌ | ❌ |
| macOS .app／Mach-O／DMG、匯出的身分 | ✅ | ✅ | ✅ |
| macOS 公證、團隊 API 認證 | ✅ | ✅ | ✅ |
| macOS 鑰匙圈身分或公證設定檔 | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
匯出的 Apple 身分與團隊 API 認證可用於可攜式簽署和公證。macOS 鑰匙圈身分與公證設定檔需要 Mac。各格式的特定限制如下所述。

@end

### 簽署後端

Wails 會自動選取最佳的可用簽署後端：

| 平台 | 原生後端 | 跨平台後端 |
| --- | --- | --- |
| Windows | `signtool.exe`，用於存放區／提供者身分、MSI 或舊版 RC2 PFX | 內建支援匯出的 RSA PFX 及 PE/MSIX |
| macOS | `codesign`，用於鑰匙圈身分 | 內建支援匯出的 PKCS#12 身分或臨時簽署 |
| Linux | 不適用 | 內建 |

在每個主機上，Apple 簽署後端皆由認證類型決定。匯出的身分使用 Go；macOS 鑰匙圈身分使用`codesign`。票證擷取、驗證及裝訂均使用 Go。

## 快速開始

設定簽署最快的方法是使用設定精靈：

```bash
wails3 setup
```

這會將共用簽署設定寫入`~/.config/wails/defaults.yaml`，且**在每個平台簽署時都會採用此設定**（請參閱[設定優先順序](#heading-4)）。其簽署步驟會：

- 顯示用於產生或匯出憑證和金鑰的設定工具，例如 OpenSSL 和 GnuPG。匯出的 Apple 與 Windows 身分使用內建簽署器；Windows 憑證存放區／提供者身分及 MSI 套件則需要 Windows SignTool。
- 在 macOS 上，列出鑰匙圈中的 Developer ID 憑證。
- 在 Linux 上，列出 GPG 金鑰，並可**產生及匯出**新的金鑰。
- 在 Windows 上，可透過 OpenSSL **產生自我簽署憑證**（用於測試）。
- **將密碼安全地儲存在系統鑰匙圈中**（而非 Taskfile 中）。

@note{type="tip"}
密碼會儲存在系統的原生認證存放區中（macOS 鑰匙圈、Windows 認證管理員或 Linux Secret Service）。因此，簽署設定既安全，也適用於所有 Wails 專案。

@end

## 設定優先順序

執行簽署工作時，每個簽署選項會依下列順序解析（第一個相符項目優先）：

1. 明確傳遞給`wails3 tool sign`的旗標（例如`--pgp-key`、`--certificate`、`--identity`）。
2. 相符的**專案 Taskfile 變數**（`PGP_KEY`、`SIGN_CERTIFICATE`／`SIGN_THUMBPRINT`、`SIGN_IDENTITY`……）。
3. `~/.config/wails/defaults.yaml`中的**全域設定**（由`wails3 setup`寫入）。

這表示 Taskfile 變數是**選用的覆寫值**：若未設定變數，便會使用全域設定的金鑰／憑證／身分。若這三者都未提供值，簽署命令會回報清楚的錯誤，告知如何進行設定。

## 個別專案設定

若要只為單一專案設定簽署，而非進行全域設定，請在專案內執行個別專案精靈；它會將`vars`寫入該專案的`build/<platform>/Taskfile.yml`檔案：

```bash
wails3 signing                                   # all detected platforms
wails3 signing --platform windows --platform linux
```

### 手動設定

或者，也可以手動編輯各平台專用的 Taskfile。請編輯每個檔案頂端的`vars`區段：

@tabs
[macOS]
編輯`build/darwin/Taskfile.yml`：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  # ENTITLEMENTS: "build/darwin/entitlements.plist"
```

然後執行：

```bash
wails3 task darwin:sign           # Sign only
wails3 task darwin:sign:notarize  # Sign and notarize
```

[Windows]
編輯`build/windows/Taskfile.yml`：

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

密碼會從系統鑰匙圈擷取（請執行`wails3 signing`進行設定）。

然後執行：

```bash
wails3 task windows:sign           # Sign executable
wails3 task windows:sign:installer # Sign NSIS installer
```

[Linux]
編輯`build/linux/Taskfile.yml`：

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

密碼會從系統鑰匙圈擷取（請執行`wails3 signing`進行設定）。

然後執行：

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

也可以直接從系統檢查簽署狀態：

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

若要以互動方式設定每個平台的簽署組態，請使用精靈：

```bash
wails3 signing
```

## macOS 程式碼簽署

### 先決條件

- Apple Developer 帳號（每年 $99）
- Developer ID Application 憑證
- 若使用 macOS 鑰匙圈身分：須已安裝 Xcode Command Line Tools
- 若使用可攜式簽署：匯出的 RSA PKCS#12 身分；不需要 Apple 簽署工具

### 使用 Wake 進行可攜式簽署

Wake 會直接讀取`wails.yaml`，不會執行 Taskfile。請設定匯出的 Developer ID 身分及 App Store Connect 團隊 API 金鑰：

```yaml
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

`credential`指定用來存放憑證密碼的環境變數。請將兩個私密金鑰檔案都排除在版本控制之外。在[YAML 建置指南](/experimental/yaml-builds/)所述的發行設定檔中啟用簽署與公證。

可攜式簽署工具支援64 位元 Intel 與 Apple Silicon Mach-O 檔案、通用二進位檔、包含巢狀程式碼與框架的標準 App 套件，以及小於 4 GiB 的 UDIF DMG。它會在發佈輸出前拒絕不支援的配置。臨時簽署使用 `identity: "-"`，且無法進行公證。PKCS#12 身分可使用 AES/PBKDF2 或舊版 3DES 加密；使用 RC2 的匯出檔必須以支援的加密演算法重新匯出。

可攜式公證簽章需要明確指定 SHA-2 RFC 3161 時間戳記服務。不接受舊版 SHA-1 時間戳記簽章。團隊 API 提交支援最大 5 GiB 的單一上傳；不支援個人 API 金鑰與多部分上傳。票證驗證與裝訂使用觀察所得的 Apple 服務通訊協定，並固定採用 Apple Root G3，無需 `xcrun stapler`。Gatekeeper 仍是判定發佈信任與撤銷狀態的權威。

以下各節說明舊版專案，以及選用鑰匙圈身分的 Wake 設定所使用的現有 macOS 鑰匙圈工作流程。

### 簽署身分

檢查可用的簽署身分：

```bash
security find-identity -v -p codesigning
```

輸出：

```
Found 2 signing identities:

  Developer ID Application: Your Company (ABCD1234) [valid]
    Hash: ABC123DEF456...

  Apple Development: your@email.com (XYZ789) [valid]
    Hash: DEF789ABC123...
```

@note{type="tip"}
若要在 App Store 以外發佈，您需要 **Developer ID Application** 憑證。

@end

### 設定

編輯 `build/darwin/Taskfile.yml` 並設定簽署變數：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| 變數 | 必要性 | 說明 |
| --- | --- | --- |
| `SIGN_IDENTITY` | 是 | 您的 Developer ID（例如「Developer ID Application: Your Company (TEAMID)」） |
| `KEYCHAIN_PROFILE` | 進行公證時 | 已儲存認證資訊的鑰匙圈描述檔名稱 |
| `ENTITLEMENTS` | 否 | 權利設定檔的路徑 |

然後執行：

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### 權利設定

權利設定控制您的 App 可存取哪些功能。Wails App 的開發版本與正式版本通常需要不同的權利設定：

- **開發**：需要 JIT、未簽署記憶體及偵錯權利
- **正式環境**：最少權利（僅限網路存取）

使用互動式設定精靈產生這兩個檔案：

```bash
wails3 entitlements
```

這會建立：

- `build/darwin/entitlements.dev.plist`－用於開發組建
- `build/darwin/entitlements.plist`－用於正式／已簽署組建

**可用的預設組態：**\

| 預設組態 | 說明 |
| --- | --- |
| 開發 | JIT、未簽署記憶體、偵錯、網路 |
| 正式環境 | 僅限網路（最精簡、最安全） |
| 兩者 | 同時建立開發與正式環境檔案（建議） |
| App Store | 啟用沙箱，並允許網路與檔案存取 |
| 自訂 | 選擇個別權利 |

@note{type="note"}
darwin Taskfile 中的 `run` 任務會自動使用 `entitlements.dev.plist`。`sign` 任務會在正式組建中使用 `entitlements.plist`。

@end

接著在 Taskfile 變數中設定 `ENTITLEMENTS`，使其指向適當的檔案。

### 公證

Apple 要求所有發佈的 App 都必須經過公證。

@steps
### **將認證資訊儲存在鑰匙圈中**（一次性設定）。您可以執行 `wails3 signing`（它會提示輸入這些值，並在底層呼叫 `notarytool`），也可以直接呼叫 `notarytool`：
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **在 Taskfile 中設定 KEYCHAIN_PROFILE**，使其與上述描述檔名稱相符。
### **簽署並公證您的 App**：
```bash
wails3 task darwin:sign:notarize
```

### **驗證公證**：
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
公證通常需要 1-2 分鐘。票證會自動裝訂至您的 App。

@end

## Windows 程式碼簽署

### 必要條件

- 程式碼簽署憑證（由 DigiCert、Sectigo 等機構核發）
- 可攜式簽署：使用 AES/PBKDF2 或舊版 3DES 加密的匯出 RSA PFX 身分；不需要外部簽署工具
- 憑證存放區／提供者身分、MSI 套件或舊版 RC2 PFX：在 `PATH` 上使用 Windows SDK 簽署工具的 Windows

### 產生自我簽署憑證（測試）

如果只需要測試簽署管線，請執行 `wails3 setup`，開啟簽署步驟的 **Windows** 分頁，然後選擇 **產生自我簽署憑證**。這會使用 OpenSSL 建立程式碼簽署 `.pfx`，並將其路徑記錄在全域設定中。

@note{type="caution"}
自我簽署憑證**僅供測試與內部發佈使用**，因為它們會對終端使用者觸發 SmartScreen 警告。公開發行需要由受信任的 CA 核發憑證。

@end

### 設定

編輯 `build/windows/Taskfile.yml` 並設定簽署變數：

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| 變數 | 必要性 | 說明 |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | 覆寫 | .pfx/.p12 憑證檔案的路徑（若未設定，則改用全域設定） |
| `SIGN_THUMBPRINT` | 覆寫 | Windows 憑證存放區中的憑證指紋（`SIGN_CERTIFICATE` 的替代選項） |
| `TIMESTAMP_SERVER` | 否 | 時間戳記伺服器 URL（預設：http://timestamp.digicert.com） |

@note{type="note"}
這些變數是選用的**覆寫值**。若未設定，則使用透過 `wails3 setup` 全域設定的憑證。請參閱[設定優先順序](#heading-4)。

@end

@note{type="note"}
憑證密碼儲存在系統鑰匙圈中，而非 Taskfile 中。請執行 `wails3 signing` 進行設定，或在 CI 中設定 `WAILS_WINDOWS_CERT_PASSWORD` 環境變數。

@end

然後執行：

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### 跨平台簽署

匯出的 RSA PFX 身分可在 Windows、macOS 和 Linux 上使用 Go 簽署 EXE/DLL 檔案、NSIS 安裝程式執行檔，以及個別的 MSIX/APPX 套件。可取代現有簽章。簽署時會先暫存輸出，只有在簽署及時間戳記驗證成功後才發布；若失敗，則保留原始檔案。使用`--output`寫入另一個已簽署檔案。

Wake 會直接從`wails.yaml`讀取這些設定：

```yaml
windows:
  signing:
    certificate: signing/windows.pfx
    credential: WINDOWS_SIGNING_PASSWORD
    timestamp_server: http://timestamp.digicert.com
```

`credential`指定密碼環境變數的名稱，絕不是密碼本身。請為發布目標啟用`sign: true`；請參閱[YAML 建置指南](/experimental/yaml-builds/)。Wake 不會叫用 Taskfile。舊版專案可繼續使用上述 Taskfile 設定。

可攜式設定檔使用 RSA/SHA-256，並驗證 RFC 3161時間戳記，包括時間戳記簽署者受信任的憑證鏈。預設服務為 DigiCert。MSIX Publisher 必須與簽署憑證的主體相符；不支援套件組合包及大小為4 GiB 以上的封存檔。PFX 檔案上限為16 MiB。其他憑證演算法及 RC2 加密會遭拒絕。若要進行可攜式簽署，請使用 AES 或3DES 重新匯出 RC2 PFX；在 Windows 上，已驗證的 RC2 PFX 會改用 SignTool。密碼錯誤或 PFX 損毀絕不會觸發此備援機制。

憑證存放區、硬體權杖及由提供者支援的金鑰，會保留使用`thumbprint`的原生 Windows SignTool 工作流程。在 Windows 上簽署 MSI 也需要 SignTool。設定匯出的憑證並不會讓這些工作流程變得可攜。

### 支援的 Windows 格式

| 格式 | 副檔名 | 備註 |
| --- | --- | --- |
| 執行檔與程式庫 | .exe、.dll | 可攜式 PE 簽署，包括 NSIS 安裝程式 |
| 安裝程式 | .msi | 需要 Windows SignTool |
| 應用程式套件 | .msix、.appx | 可攜式個別套件簽署；不支援組合包 |

## Linux 套件簽署

Linux 套件（DEB 和 RPM）使用 PGP/GPG 金鑰簽署。不同於 Windows 和 macOS 程式碼簽署，Linux 套件簽署證明套件來自受信任的來源，而非證明作業系統信任其程式碼。

匯出的 ASCII Armor 或二進位私密金鑰檔案，會在 Windows、macOS 和 Linux 上使用原生 Go 簽署器。簽署現有 DEB 或 RPM 時，無須解壓縮或重新壓縮即可保留其套件內容，也不需要安裝 GnuPG、`dpkg-sig`或`rpmsign`。Arch 套件使用相同的金鑰載入器寫入分離式簽章。GnuPG 仍適合用於產生／匯出金鑰，以及獨立驗證簽章。以金鑰環身分而非金鑰檔案進行設定的舊版 DEB/RPM 組態，會保留其外部簽署工具。

### 必要條件

- PGP 金鑰組（可使用 Wails 產生）

### 產生 PGP 金鑰

最簡單的方法是使用設定精靈：執行`wails3 setup`，開啟簽署步驟中的**Linux**分頁，然後選擇**建立新的 GPG 金鑰**。精靈會：

- 在您的 GPG 金鑰環中產生 RSA 4096金鑰（若要建立適合無人值守及 CI 的金鑰，請將複雜密碼留空），
- **將其匯出至`~/.wails/signing/<keyid>.asc`**（建置用來簽署的檔案），以及
- 在`~/.config/wails/defaults.yaml`中記錄金鑰 ID 和匯出路徑，以便簽署時自動使用。

@note{type="note"}
若金鑰僅以**ID**設定（例如在精靈開始自動匯出金鑰之前建立），下次開啟 Linux 分頁時，精靈會將其匯出至檔案並填入路徑，無須任何手動步驟。建置會使用金鑰*檔案*進行簽署，因此路徑才是關鍵。

@end

您也可以使用`gpg`手動完成：

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

建議設定：RSA 4096位元、1年後到期，並以高強度密碼保護。

@note{type="caution"}
請妥善保護您的私密金鑰！以加密方式儲存，並安全地備份。

@end

### 設定

編輯`build/linux/Taskfile.yml`並設定簽署變數：

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| 變數 | 必要性 | 說明 |
| --- | --- | --- |
| `PGP_KEY` | 覆寫 | 匯出的 PGP 私密金鑰檔案路徑（若未設定，則退回使用透過`wails3 setup`全域設定的金鑰） |
| `SIGN_ROLE` | 否 | DEB 簽署角色（預設：builder） |

@note{type="note"}
PGP 金鑰密碼儲存在系統鑰匙圈中，而非 Taskfile。執行`wails3 signing`進行設定，或在 CI 中設定`WAILS_PGP_PASSWORD`環境變數。

@end

接著執行：

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### DEB 簽署角色

對於 DEB 套件，您可以透過`SIGN_ROLE`指定簽署角色：

- `origin`：來自套件來源的簽章
- `maint`：來自套件維護者的簽章
- `archive`：來自封存庫維護者的簽章
- `builder`：來自套件建置者的簽章（預設）

### 跨平台簽署

可從任何平台簽署 Linux 套件。相同的 Taskfile 設定及命令也適用於 Windows 和 macOS。

### 檢視金鑰資訊

```bash
gpg --show-keys signing-key.asc
```

輸出：

```
pub   rsa4096 2024-01-15 [SC] [expires: 2025-01-15]
      1234 5678 90AB CDEF 1234 5678 90AB CDEF 1234 5678
uid                      Your Name <your@email.com>
```

### 驗證 Linux 套件

```bash
# Verify DEB signature
dpkg-sig --verify myapp_1.0.0_amd64.deb

# Verify RPM signature
rpm --checksig myapp-1.0.0.x86_64.rpm
```

### 散布您的公開金鑰

使用者需要您的公開金鑰才能驗證套件：

```bash
# Export public key for distribution
gpg --armor --export "your@email.com" > myapp-signing.pub.asc

# Users can import it:
# For DEB (apt):
sudo apt-key add myapp-signing.pub.asc
# Or for modern apt:
sudo cp myapp-signing.pub.asc /etc/apt/trusted.gpg.d/

# For RPM:
sudo rpm --import myapp-signing.pub.asc
```

## GitHub Actions 整合

在 CI 環境中，密碼會透過環境變數提供，而非使用系統鑰匙圈：

| 環境變數 | 說明 |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Windows 憑證密碼 |
| `WAILS_PGP_PASSWORD` | Linux 套件的 PGP 金鑰密碼 |

您也可以直接傳入 Taskfile 變數：

```bash
wails3 task darwin:sign SIGN_IDENTITY="$SIGN_IDENTITY" KEYCHAIN_PROFILE="$KEYCHAIN_PROFILE"
```

### macOS 工作流程

```yaml
name: Build and Sign macOS

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.MACOS_CERTIFICATE }}
          CERTIFICATE_PASSWORD: ${{ secrets.MACOS_CERTIFICATE_PASSWORD }}
        run: |
          echo $CERTIFICATE_BASE64 | base64 --decode > certificate.p12
          security create-keychain -p "" build.keychain
          security default-keychain -s build.keychain
          security unlock-keychain -p "" build.keychain
          security import certificate.p12 -k build.keychain -P "$CERTIFICATE_PASSWORD" -T /usr/bin/codesign
          security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "" build.keychain

      - name: Store Notarization Credentials
        env:
          APPLE_ID: ${{ secrets.APPLE_ID }}
          APPLE_TEAM_ID: ${{ secrets.APPLE_TEAM_ID }}
          APPLE_APP_PASSWORD: ${{ secrets.APPLE_APP_PASSWORD }}
        run: |
          xcrun notarytool store-credentials "notarize-profile" \
            --apple-id "$APPLE_ID" \
            --team-id "$APPLE_TEAM_ID" \
            --password "$APPLE_APP_PASSWORD"

      - name: Build, Sign, and Notarize
        env:
          SIGN_IDENTITY: ${{ secrets.MACOS_SIGN_IDENTITY }}
        run: |
          wails3 task darwin:sign:notarize \
            SIGN_IDENTITY="$SIGN_IDENTITY" \
            KEYCHAIN_PROFILE="notarize-profile"

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-macOS
          path: bin/*.app
```

### Windows 工作流程

```yaml
name: Build and Sign Windows

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
        run: |
          $certBytes = [Convert]::FromBase64String($env:CERTIFICATE_BASE64)
          [IO.File]::WriteAllBytes("certificate.pfx", $certBytes)

      - name: Build and Sign
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-Windows
          path: bin/*.exe
```

### 跨平台工作流程（Linux 執行器）

從單一 Linux 執行器簽署 Windows 和 Linux 套件：

```yaml
name: Build and Sign (Cross-Platform)

on:
  push:
    tags: ['v*']

jobs:
  build-and-sign:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Install Build Dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y nsis rpm

      # Import certificates
      - name: Import Certificates
        env:
          WINDOWS_CERT_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
          PGP_KEY_BASE64: ${{ secrets.PGP_PRIVATE_KEY }}
        run: |
          echo "$WINDOWS_CERT_BASE64" | base64 -d > certificate.pfx
          echo "$PGP_KEY_BASE64" | base64 -d > signing-key.asc

      # Build and sign Windows
      - name: Build and Sign Windows
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      # Build and sign Linux packages
      - name: Build and Sign Linux Packages
        env:
          WAILS_PGP_PASSWORD: ${{ secrets.PGP_PASSWORD }}
        run: |
          wails3 task linux:sign:packages PGP_KEY=signing-key.asc

      # Cleanup secrets
      - name: Cleanup
        if: always()
        run: rm -f certificate.pfx signing-key.asc

      - name: Upload Artifacts
        uses: actions/upload-artifact@v4
        with:
          name: signed-binaries
          path: |
            bin/*.exe
            bin/*.deb
            bin/*.rpm
```

@note{type="note"}
使用 Linux 執行器進行跨平台簽署，可免除另外使用 Windows 執行器的需求，從而簡化 CI/CD。Apple 鑰匙圈身分和公證設定檔需要 macOS 執行器。匯出的身分和團隊 API 認證支援上述可攜式 Wake 工作流程。

@end

## CLI 參考

### wails3 signing

用於設定專案簽署的互動式精靈。

```bash
wails3 signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

精靈會引導您完成：

- **macOS**：選取 Developer ID 憑證、設定公證憑證（在內部呼叫`xcrun notarytool store-credentials`）。
- **Windows**：選擇使用憑證檔案或指紋，並設定密碼與時間戳記伺服器。
- **Linux**：使用現有 PGP 金鑰或產生新金鑰（呼叫`gpg`），並設定簽署角色。

### wails3 entitlements

用於設定 macOS 權利的互動式精靈。

```bash
wails3 entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**預設集：**

- **開發**：建立包含 JIT、偵錯及網路權利的`entitlements.dev.plist`
- **正式環境**：建立僅含最低限度權利的`entitlements.plist`
- **兩者**：建立兩個檔案（建議）
- **App Store**：建立適用於 Mac App Store 的沙箱化權利
- **自訂**：選擇個別權利與目標檔案

### wails3 sign

簽署目前平台或指定平台的二進位檔與套件。這是一層包裝，會呼叫適合該平台的簽署工作。

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

這會執行對應的`<platform>:sign`工作，該工作會使用 Taskfile 中的簽署設定。

### wails3 tool sign

直接簽署特定檔案的低階命令。Taskfile 會在內部使用此命令。

```bash
wails3 tool sign [flags]
```

**通用旗標：**\

| 旗標 | 說明 |
| --- | --- |
| `--input` | 要簽署的檔案路徑 |
| `--output` | 輸出路徑（選用，預設為就地寫入） |
| `--verbose` | 啟用詳細輸出 |

**Windows/macOS 旗標：**\

| 旗標 | 說明 |
| --- | --- |
| `--certificate` | PKCS#12（.pfx/.p12）憑證的路徑 |
| `--password` | 憑證密碼 |
| `--timestamp` | 時間戳記伺服器 URL |

**macOS 專用旗標：**\

| 旗標 | 說明 |
| --- | --- |
| `--identity` | 簽署身分（使用「-」進行臨時簽署） |
| `--entitlements` | 權利 plist 的路徑 |
| `--hardened-runtime` | 啟用強化執行階段（預設值：true） |
| `--notarize` | 送交公證 |
| `--keychain-profile` | 用於公證的鑰匙圈設定檔 |
| `--notary-key` | App Store Connect 團隊 API 私密金鑰`.p8`檔案；使用 Go Notary API 用戶端 |
| `--notary-key-id` | 團隊 API 金鑰 ID；與`--notary-key`搭配使用時為必填 |
| `--notary-issuer-id` | 團隊簽發者 ID；與`--notary-key`搭配使用時為必填 |

**Windows 專用旗標：**\

| 旗標 | 說明 |
| --- | --- |
| `--thumbprint` | Windows 憑證存放區中的憑證指紋 |

**Linux 專用旗標：**\

| 旗標 | 說明 |
| --- | --- |
| `--pgp-key` | PGP 私密金鑰的路徑 |
| `--pgp-password` | PGP 金鑰密碼 |
| `--role` | DEB 簽署角色（origin/maint/archive/builder） |

### 檢查簽署狀態（原生工具）

`wails3 signing`會設定簽署憑證。若要檢查簽署狀態，請直接使用原生工具：

| 工作 | 命令 |
| --- | --- |
| 列出 macOS 程式碼簽署身分 | `security find-identity -v -p codesigning` |
| 儲存公證憑證 | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| 檢查 PGP 金鑰檔案 | `gpg --show-keys <key.asc>` |
| 產生 PGP 金鑰組 | `gpg --full-generate-key` |
| 匯出公開金鑰 | `gpg --armor --export <email>` |

## 疑難排解

### macOS 問題

**「找不到 Developer ID 憑證」**

- 確認憑證已安裝於鑰匙圈中
- 使用`security find-identity -v -p codesigning`檢查憑證是否尚未過期
- 確認您擁有「Developer ID Application」憑證，而不只是「Apple Development」憑證

**「公證失敗」**

- 檢查公證記錄：`xcrun notarytool log <submission-id> --keychain-profile <profile>`
- 確認已啟用強化執行階段
- 確認您的應用程式不包含未簽署的二進位檔

**「程式碼簽署失敗」**

- 確認鑰匙圈已解鎖：`security unlock-keychain`
- 檢查應用程式套件組合的檔案權限

### Windows 問題

**「找不到憑證」**

- 確認憑證路徑正確
- 檢查憑證密碼
- 確保憑證有效（未過期或遭撤銷）

**「時間戳記伺服器錯誤」**

- 請嘗試其他時間戳記伺服器：
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Linux 問題

**「無效的 PGP 金鑰」**

- 確保金鑰檔案採用 ASCII 裝甲格式
- 使用`gpg --show-keys <key.asc>`檢查金鑰是否尚未過期
- 確認密碼正確

**「簽章驗證失敗」**

- 確保已正確匯入公開金鑰
- 檢查套件在簽署後是否未遭修改

## 其他資源

### 官方文件

- [Apple 程式碼簽署指南](https://developer.apple.com/support/code-signing/)
- [Apple 公證文件](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Microsoft 程式碼簽署](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Debian 套件簽署](https://wiki.debian.org/SecureApt)
- [RPM 套件簽署](https://rpm-software-management.github.io/rpm/manual/signatures.html)
