---
title: "程式碼簽署"
description: "在所有平台上簽署 Wails 應用程式的指南"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## 簽署應用程式的程式碼

本指南說明如何為 macOS、Windows 和 Linux 簽署 Wails 應用程式。Wails v3 提供內建 CLI 工具，可用於程式碼簽署、公證及 PGP 金鑰管理。

- **macOS**－簽署並公證 macOS 應用程式
- **Windows**－簽署 Windows 可執行檔和套件
- **Linux**－使用 PGP 金鑰簽署 DEB 和 RPM 套件

## 跨平台簽署矩陣

此矩陣顯示可從各來源平台簽署的項目：

| 目標格式 | 從 Windows | 從 macOS | 從 Linux |
| --- | :---: | :---: | :---: |
| Windows EXE/MSI | ✅ | ✅ | ✅ |
| macOS .app 套件組合 | ❌ | ✅ | ❌ |
| macOS 公證 | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
Windows 和 Linux 套件可以從<strong>任何平台</strong>簽署。由於 Apple 工具的要求，macOS 簽署必須使用 Mac。

@end

### 簽署後端

Wails 會自動選擇最佳的可用簽署後端：

| 平台 | 原生後端 | 跨平台後端 |
| --- | --- | --- |
| Windows | `signtool.exe`（Windows SDK） | 內建 |
| macOS | `codesign`（Xcode） | 不可用 |
| Linux | 不適用 | 內建 |

在原生平台上執行時，Wails 會使用原生工具以達到最大的相容性。進行交叉編譯時，則會使用內建的簽署支援。

## 快速開始

設定簽署最快的方法是使用設定精靈：

```bash
wails3 setup
```

這會將共用簽署設定寫入`~/.config/wails/defaults.yaml`，而且在每個平台執行簽署時都會 **採用此設定**（請參閱[設定優先順序](#heading-4)）。 其簽署步驟會：

- 偵測針對您要建置的目標平台已安裝哪些簽署工具（可<strong>從任何主機</strong>執行），並顯示適用於您作業系統的安裝命令（例如`brew install gnupg`、`sudo apt install osslsigncode`、`winget install GnuPG.Gpg4win`）。
- 在 macOS 上，列出鑰匙圈中的 Developer ID 憑證。
- 針對 Linux，列出您的 GPG 金鑰，並可<strong>產生及匯出</strong>新金鑰。
- 針對 Windows，可透過 OpenSSL<strong>產生自我簽署憑證</strong>（供測試使用）。
- **將密碼安全地儲存在系統鑰匙圈中**（而非 Taskfile 中）。

@note{type="tip"}
密碼會儲存在系統的原生認證資料存放區中（macOS 鑰匙圈、Windows 認證管理員或 Linux Secret Service）。這表示您的簽署設定既安全，又能用於所有 Wails 專案。

@end

## 設定優先順序

執行簽署工作時，會依下列順序解析每個簽署選項（以第一個 相符項目為準）：

1. 明確傳遞給`wails3 tool sign`的旗標（例如`--pgp-key`、`--certificate`、`--identity`）。
2. 相符的<strong>專案 Taskfile 變數</strong>（`PGP_KEY`、`SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`、`SIGN_IDENTITY`……）。
3. `~/.config/wails/defaults.yaml`中的<strong>全域設定</strong>（由`wails3 setup`寫入）。

這表示 Taskfile 變數是<strong>選用的覆寫值</strong>：若未設定某個變數，便會使用 全域設定的金鑰、憑證或身分識別。如果這三者皆未提供 值，簽署命令會回報明確的錯誤，告知您如何進行設定。

## 個別專案設定

若要只為單一專案設定簽署，而非進行全域設定，請在專案內執行個別專案 精靈；它會將`vars`寫入該專案的 `build/<platform>/Taskfile.yml`檔案：

```bash
wails3 setup signing                                   # all detected platforms
wails3 setup signing --platform windows --platform linux
```

### 手動設定

或者，您可以手動編輯各平台專用的 Taskfile。請編輯每個檔案頂端的`vars`區段：

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

密碼會從系統鑰匙圈取得（請執行`wails3 setup signing`進行設定）。

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

密碼會從系統鑰匙圈擷取（執行`wails3 setup signing`進行設定）。

然後執行：

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

您也可以直接從系統檢查簽署狀態：

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

若要以互動方式設定所有平台的簽署組態，請使用精靈：

```bash
wails3 setup signing
```

## macOS 程式碼簽署

### 必要條件

- Apple Developer 帳號（每年 $99）
- Developer ID Application 憑證
- 已安裝 Xcode Command Line Tools

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
若要在 App Store 以外散布，您需要<strong>Developer ID Application</strong>憑證。

@end

### 組態

編輯`build/darwin/Taskfile.yml`並設定簽署變數：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| 變數 | 必要 | 說明 |
| --- | --- | --- |
| `SIGN_IDENTITY` | 是 | 您的 Developer ID（例如「Developer ID Application: Your Company (TEAMID)」） |
| `KEYCHAIN_PROFILE` | 公證時需要 | 已儲存認證資訊的鑰匙圈描述檔名稱 |
| `ENTITLEMENTS` | 否 | 權利檔案的路徑 |

然後執行：

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### 權利

權利會控制您的應用程式可使用哪些功能。Wails 應用程式的開發版與正式版通常需要不同的權利：

- **開發版**：需要 JIT、未簽署記憶體及偵錯權利
- **正式版**：最少的權限（僅網路存取權）

使用互動式設定精靈產生這兩個檔案：

```bash
wails3 setup entitlements
```

這會建立：

- `build/darwin/entitlements.dev.plist`－用於開發組建
- `build/darwin/entitlements.plist`－用於正式版／已簽署組建

**可用的預設組態：**\

| 預設組態 | 說明 |
| --- | --- |
| 開發版 | JIT、未簽署記憶體、偵錯、網路 |
| 正式版 | 僅限網路（最少且最安全） |
| 兩者 | 同時建立開發版與正式版檔案（建議） |
| App Store | 啟用沙箱，並允許網路與檔案存取 |
| 自訂 | 個別選擇權限 |

@note{type="note"}
darwin Taskfile 中的`run`工作會自動使用`entitlements.dev.plist`。`sign`工作則會對正式版組建使用`entitlements.plist`。

@end

然後在 Taskfile 的 vars 中設定`ENTITLEMENTS`，使其指向適當的檔案。

### 公證

Apple 要求所有散布的應用程式都必須經過公證。

@steps
### **將認證資訊儲存在鑰匙圈中**（一次性設定）。您可以執行 `wails3 setup signing`（它會提示您輸入這些值，並在內部呼叫 `notarytool`），或直接呼叫 `notarytool`：
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **在 Taskfile 中設定 KEYCHAIN_PROFILE**，使其符合上述描述檔名稱。
### **簽署並公證您的應用程式**：
```bash
wails3 task darwin:sign:notarize
```

### **驗證公證結果**：
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
公證通常需要1-2分鐘。票證會自動裝訂至您的應用程式。

@end

## Windows 程式碼簽署

### 必要條件

- 程式碼簽署憑證（由 DigiCert、Sectigo 等機構核發）
- 若要在 Windows 上進行原生簽署：必須安裝 Windows SDK（以使用`signtool.exe`）
- 若要從 macOS/Linux 進行跨平台簽署：[`osslsigncode`](https://github.com/mtrojnar/osslsigncode)（`wails3 setup`中的簽署步驟會顯示適用於該主機的安裝命令）

### 產生自我簽署憑證（測試用）

如果只需要測試簽署流程，請執行 `wails3 setup`，開啟簽署步驟中的 **Windows** 分頁，然後選擇<strong>產生自我簽署憑證</strong>。這會使用 OpenSSL 建立程式碼簽署用的 `.pfx`，並將其路徑記錄在全域設定中。

@note{type="caution"}
自我簽署憑證<strong>僅適用於測試和內部散布</strong>，因為它們會向終端使用者觸發 SmartScreen 警告。公開發行版本需要由受信任的 CA 核發憑證。

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
| `SIGN_CERTIFICATE` | 覆寫值 | .pfx/.p12 憑證檔案的路徑（若未設定，則使用全域設定） |
| `SIGN_THUMBPRINT` | 覆寫值 | Windows 憑證存放區中的憑證指紋（`SIGN_CERTIFICATE` 的替代方式） |
| `TIMESTAMP_SERVER` | 否 | 時間戳記伺服器 URL（預設：http://timestamp.digicert.com） |

@note{type="note"}
這些變數是選用的<strong>覆寫值</strong>。若未設定，則使用透過 `wails3 setup` 設定的全域憑證。請參閱[設定優先順序](#heading-4)。

@end

@note{type="note"}
憑證密碼會儲存在系統鑰匙圈中，而不是 Taskfile 內。請執行 `wails3 setup signing` 進行設定，或在 CI 中設定 `WAILS_WINDOWS_CERT_PASSWORD` 環境變數。

@end

然後執行：

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### 跨平台簽署

您可以從任何平台簽署 Windows 可執行檔。同一套 Taskfile 設定和命令也適用於 macOS 與 Linux。

### 支援的 Windows 格式

| 格式 | 副檔名 | 備註 |
| --- | --- | --- |
| 可執行檔 | .exe | 標準 PE 簽署 |
| 安裝程式 | .msi | Windows Installer 套件 |
| 應用程式套件 | .msix, .appx | 現代 Windows 應用程式 |

## Linux 套件簽署

Linux 套件（DEB 和 RPM）使用 PGP/GPG 金鑰簽署。與 Windows 和 macOS 程式碼簽署不同，Linux 套件簽署證明套件來自受信任的來源，而不是證明作業系統信任其程式碼。

### 必要條件

- PGP 金鑰組（可使用 Wails 產生）

### 產生 PGP 金鑰

最簡單的方式是使用設定精靈：執行 `wails3 setup`，開啟簽署步驟中的 **Linux** 分頁，然後選擇<strong>建立新的 GPG 金鑰</strong>。精靈會：

- 在您的 GPG 金鑰圈中產生 RSA 4096 金鑰（若要讓金鑰適合無人值守的 CI，請將複雜密碼留空），
- **將其匯出至 `~/.wails/signing/<keyid>.asc`**（建置簽署時使用的檔案），並且
- 在 `~/.config/wails/defaults.yaml` 中記錄金鑰 ID 和匯出路徑，以便簽署時自動使用。

@note{type="note"}
如果設定金鑰時<strong>只設定了 ID</strong>（例如在精靈尚未自動匯出金鑰前建立），下次開啟 Linux 分頁時，系統會將金鑰匯出至檔案並填入路徑，不需要任何手動步驟。建置會使用金鑰<em>檔案</em>進行簽署，因此真正重要的是路徑。

@end

您也可以使用 `gpg` 手動完成：

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

建議設定：RSA 4096 位元、1 年後到期，並使用高強度密碼保護。

@note{type="caution"}
請妥善保護您的私密金鑰！請加密儲存，並安全地備份。

@end

### 設定

編輯 `build/linux/Taskfile.yml` 並設定簽署變數：

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| 變數 | 必要性 | 說明 |
| --- | --- | --- |
| `PGP_KEY` | 覆寫 | 匯出的 PGP 私密金鑰檔案路徑（若未設定，則使用透過`wails3 setup`全域設定的金鑰） |
| `SIGN_ROLE` | 否 | DEB 簽署角色（預設：builder） |

@note{type="note"}
PGP 金鑰密碼儲存在系統鑰匙圈中，而非 Taskfile 中。請執行`wails3 setup signing`進行設定，或在 CI 中設定`WAILS_PGP_PASSWORD`環境變數。

@end

接著執行：

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### DEB 簽署角色

對於 DEB 套件，可以透過`SIGN_ROLE`指定簽署角色：

- `origin`：來自套件來源的簽章
- `maint`：來自套件維護者的簽章
- `archive`：來自封存庫維護者的簽章
- `builder`：來自套件建置者的簽章（預設）

### 跨平台簽署

可以從任何平台簽署 Linux 套件。同一套 Taskfile 設定與命令也適用於 Windows 和 macOS。

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

在 CI 環境中，密碼會透過環境變數提供，而不使用系統鑰匙圈：

| 環境變數 | 說明 |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Windows 憑證密碼 |
| `WAILS_PGP_PASSWORD` | Linux 套件的 PGP 金鑰密碼 |

也可以直接傳入 Taskfile 變數：

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

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
使用 Linux 執行器進行跨平台簽署，不需要另外使用 Windows 執行器，因此能簡化 CI/CD。由於 Apple 原生工具的要求，macOS 簽署仍需要 macOS 執行器。

@end

## CLI 參考資料

### wails3 setup signing

用於設定專案簽署的互動式精靈。

```bash
wails3 setup signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

此精靈將引導您完成：

- **macOS**：選取 Developer ID 憑證、設定公證認證資訊（呼叫`xcrun notarytool store-credentials`）。
- **Windows**：選擇使用憑證檔案或指紋，並設定密碼和時間戳記伺服器。
- **Linux**：使用現有的 PGP 金鑰或產生新金鑰（呼叫`gpg`），並設定簽署角色。

### wails3 setup entitlements

用於設定 macOS 權利的互動式精靈。

```bash
wails3 setup entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**預設集：**

- **開發**：建立含有 JIT、偵錯和網路權利的`entitlements.dev.plist`
- **正式環境**：建立僅含最低限度權利的`entitlements.plist`
- **兩者**：建立這兩個檔案（建議）
- **App Store**：建立適用於 Mac App Store 的沙箱化權利
- **自訂**：選擇個別權利和目標檔案

### wails3 sign

為目前或指定的平台簽署二進位檔和套件。這是一個包裝命令，會呼叫適用於該平台的簽署工作。

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

這會執行對應的`<platform>:sign`工作，而該工作會使用 Taskfile 中的簽署設定。

### wails3 tool sign

直接簽署特定檔案的低階命令。供 Taskfile 內部使用。

```bash
wails3 tool sign [flags]
```

**通用旗標：**\

| 旗標 | 說明 |
| --- | --- |
| `--input` | 要簽署的檔案路徑 |
| `--output` | 輸出路徑（選填，預設為原地寫入） |
| `--verbose` | 啟用詳細輸出 |

**Windows/macOS 旗標：**\

| 旗標 | 說明 |
| --- | --- |
| `--certificate` | PKCS#12 憑證（.pfx/.p12）的路徑 |
| `--password` | 憑證密碼 |
| `--timestamp` | 時間戳記伺服器 URL |

**macOS 專用旗標：**\

| 旗標 | 說明 |
| --- | --- |
| `--identity` | 簽署身分（臨時簽署請使用 '-'） |
| `--entitlements` | 權利 plist 的路徑 |
| `--hardened-runtime` | 啟用強化執行階段（預設：true） |
| `--notarize` | 提交以進行公證 |
| `--keychain-profile` | 用於公證的鑰匙圈描述檔 |

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

v3 中<strong>沒有</strong> `wails3 signing` 命令。若要檢查簽署狀態，請直接使用原生工具：

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

- 確認您的憑證已安裝至「鑰匙圈」
- 使用`security find-identity -v -p codesigning`檢查憑證是否尚未到期
- 確認您擁有「Developer ID Application」憑證（而不只有「Apple Development」憑證）

**「公證失敗」**

- 檢查公證記錄：`xcrun notarytool log <submission-id> --keychain-profile <profile>`
- 確認已啟用強化執行階段
- 確認您的應用程式未包含未簽署的二進位檔

**「Codesign 失敗」**

- 請確認鑰匙圈已解鎖：`security unlock-keychain`
- 檢查應用程式套件的檔案權限

### Windows 問題

**「找不到憑證」**

- 確認憑證路徑正確
- 檢查憑證密碼
- 確保憑證有效（未過期且未遭撤銷）

**「時間戳記伺服器錯誤」**

- 嘗試其他時間戳記伺服器：
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Linux 問題

**「無效的 PGP 金鑰」**

- 確保金鑰檔案採用 ASCII Armor 格式
- 使用`gpg --show-keys <key.asc>`檢查金鑰是否已過期
- 確認密碼正確

**「簽章驗證失敗」**

- 確保公開金鑰已正確匯入
- 檢查套件在簽署後是否遭到修改

## 其他資源

### 官方文件

- [Apple 程式碼簽署指南](https://developer.apple.com/support/code-signing/)
- [Apple 公證文件](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Microsoft 程式碼簽署](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Debian 套件簽署](https://wiki.debian.org/SecureApt)
- [RPM 套件簽署](https://rpm-software-management.github.io/rpm/manual/signatures.html)
