---
title: "更新資訊清單協定"
description: "Wails 應用程式用來探索並驗證自我更新的開放式 JSON 協定，可由任何靜態檔案主機或動態更新伺服器提供。"
slug: "reference/update-manifest"
sourcePath: "reference/update-manifest.md"
---

Wails Update Manifest 協定是 Wails 應用程式與更新來源之間一份精簡且開放的 JSON 約定。任何可透過 HTTPS 提供 JSON 檔案的服務都能提供 Wails 更新，例如 S3 儲存貯體、GitHub Pages、CDN，或依授權決定是否提供版本的動態更新伺服器。

用戶端以 `endpoint` 提供者（`github.com/wailsapp/wails/v3/pkg/updater/providers/endpoint`）的形式隨框架提供。本頁是供伺服器端實作者參考的線路格式規格。

## 設計目標

1. <strong>適合靜態託管。</strong>每個通道只需一份資訊清單檔案，列出各平台的成品，即構成完整實作。不需要伺服器端程式碼。
2. <strong>適合動態伺服器。</strong>用戶端每次檢查都會傳送 `platform`、`arch`、`version` 和 `channel`，因此伺服器可以只回應一個成品、套用授權規則，或在呼叫端已是最新版本時傳回 `204 No Content`。
3. <strong>優先驗證。</strong>資訊清單會為每個成品提供總和檢查碼與簽章，而 Wails 更新程式會使用建置時固定在應用程式二進位檔中的公開金鑰加以驗證。更新來源絕不會自行選擇信任根。

## 請求

用戶端會向設定的資訊清單 URL 發出 `GET`，並附帶 `Accept: application/json` 以及應用程式設定的任何標頭（例如 `Authorization: License <key>`）。

URL 可內嵌預留位置，用戶端會在每次檢查時替換這些預留位置：

| 預留位置 | 替換為 |
| --- | --- |
| `{{platform}}` | 執行中的作業系統，以 Go `GOOS` 值表示（`darwin`、`windows`、`linux`） |
| `{{arch}}` | 執行中的架構，以 Go `GOARCH` 值表示（`amd64`、`arm64`……） |
| `{{version}}` | 目前安裝的版本 |
| `{{channel}}` | 已設定的發行通道（若有設定） |

四個值中，任何未由預留位置使用的值，都會以同名查詢參數附加至 URL（`channel` 僅在已設定時附加）。因此，以下兩種設定都有效且等效：

```text
# Dynamic server: reads query parameters
https://updates.example.com/check
  -> GET /check?platform=darwin&arch=arm64&version=1.0.0&channel=stable

# Static host: one manifest per platform/arch/channel path
https://cdn.example.com/updates/{{platform}}/{{arch}}/{{channel}}.json
  -> GET /updates/darwin/arm64/stable.json?version=1.0.0
```

靜態主機只需忽略收到的查詢參數。

## 回應

| 狀態 | 含義 |
| --- | --- |
| `200 OK` | 後接資訊清單。用戶端自行判斷這是否為升級。 |
| `204 No Content` | 伺服器已比較版本，且呼叫端目前是最新版本。 |
| `404 Not Found` | 尚未發佈任何內容（處理方式與已是最新版本相同）。 |
| 任何其他狀態 | 發生錯誤。更新程式會接著嘗試下一個已設定的提供者。 |

`200` 回應主體是一份資訊清單文件：

```json
{
  "schemaVersion": 1,
  "version": "2.1.0",
  "channel": "stable",
  "name": "Summer Release",
  "notes": "## What's new\n\n- Faster startup\n- New themes",
  "publishedAt": "2026-07-03T10:00:00Z",
  "artifacts": [
    {
      "url": "MyApp-2.1.0-darwin-arm64.zip",
      "platform": "darwin",
      "arch": "arm64",
      "filetype": "zip",
      "size": 8388608,
      "digestAlgo": "sha512",
      "digest": "base64-encoded digest bytes",
      "signatureAlgo": "ed25519ph",
      "signature": "base64-encoded signature bytes"
    },
    {
      "url": "MyApp-2.1.0-windows-amd64.zip",
      "platform": "windows",
      "arch": "amd64",
      "filetype": "zip",
      "size": 9437184,
      "digestAlgo": "sha512",
      "digest": "...",
      "signatureAlgo": "ed25519ph",
      "signature": "..."
    }
  ]
}
```

### 頂層欄位

| 欄位 | 型別 | 必填 | 備註 |
| --- | --- | --- | --- |
| `schemaVersion` | int | 否 | 協定版本。省略時表示 `1`。用戶端會拒絕比其可理解版本更新的值。 |
| `version` | string | **是** | SemVer 2.0.0，可帶或不帶前置的 `v`。 |
| `channel` | string | 否 | 僅供參考。若用戶端設定為其他通道，會將此資訊清單視為沒有更新。 |
| `name` | string | 否 | 供人閱讀的版本標題，顯示在更新視窗中。 |
| `notes` | string | 否 | 以 Markdown 撰寫的版本資訊，會在更新視窗中呈現。 |
| `publishedAt` | string | 否 | RFC 3339 時間戳記。 |
| `artifacts` | array | **是** | 每個可下載成品各有一筆項目。順序表示發佈者的偏好。 |
| `metadata` | object | 否 | 自由格式的鍵/值資料，會原樣傳遞給應用程式。 |

用戶端會忽略未知欄位，因此伺服器可加入自訂欄位，而不會破壞任何用戶端。伺服器專用的擴充內容應放在`metadata`中。

### 成品欄位

| 欄位 | 型別 | 必填 | 備註 |
| --- | --- | --- | --- |
| `url` | 字串 | **是** | 絕對 URL，或相對於資訊清單 URL 的相對 URL。僅限`http(s)`。 |
| `platform` | 字串 | 否 | Go 的`GOOS`值。接受常見別名（`macos`、`win`等）。空值會比對所有平台。 |
| `arch` | 字串 | 否 | Go 的`GOARCH`值。接受常見別名（`x86_64`、`aarch64`等）。空值會比對所有架構。 |
| `filename` | 字串 | 否 | 預設為`url`的最後一個路徑區段。 |
| `filetype` | 字串 | 否 | 預設為副檔名。 |
| `size` | 整數 | 否 | 位元組數，用於顯示下載進度。 |
| `digestAlgo` / `digest` | 字串 / base64 | 否 | `sha256`或`sha512`。 |
| `signatureAlgo` / `signature` | 字串 / base64 | 否 | `ed25519`、`ed25519ph`或`ecdsa-p256`。只要有`signature`，就必須提供`signatureAlgo`。請參閱[更新程式指南](/guides/updater/#cryptographic-verification)，瞭解各演算法簽署的內容。 |

用戶端會選取<strong>第一個</strong>其`platform`和`arch`與執行中系統相符的成品。Base64 值無論是否含有填補字元皆可接受。

### 版本比較

資訊清單是否代表升級，一律由用戶端依據 SemVer 2.0.0優先順序判定：資訊清單中的`version`必須嚴格地比已安裝版本更新。這讓靜態託管能輕易維持正確（資訊清單始終描述最新版本，而已是最新版本的用戶端不會執行任何動作），同時動態伺服器仍可使用`204`來節省頻寬。

## 驗證與信任

總和檢查碼和簽章會隨資訊清單提供，但信任根不會；簽章是使用應用程式在建置時透過`updater.Config.PublicKey`固定的公開金鑰進行驗證。遭入侵或替換的更新來源無法提供自己的金鑰。若應用程式未固定金鑰，帶有簽章的成品就會遭到拒絕；若簽章未宣告`signatureAlgo`或無法解碼，成品也會遭到拒絕：用戶端絕不會在未提示的情況下退回僅驗證摘要。

僅含摘要的成品會在摘要檢查通過後安裝；這可防止資料損毀，但其防竄改能力仰賴 TLS 及主機本身的完整性。任何涉及安全性的內容都應提供簽章。

使用框架的`ed25519ph`機制簽署成品，只需要幾行 Go 程式碼：

```go
digest := sha512.Sum512(artifactBytes)
sig, _ := privateKey.Sign(nil, digest[:], &ed25519.Options{Hash: crypto.SHA512})
manifest.Artifacts[i].DigestAlgo = "sha512"
manifest.Artifacts[i].Digest = base64.StdEncoding.EncodeToString(digest[:])
manifest.Artifacts[i].SignatureAlgo = "ed25519ph"
manifest.Artifacts[i].Signature = base64.StdEncoding.EncodeToString(sig)
```

實務上很少需要自行撰寫這段程式碼：CLI 會代為處理。

## 使用 wails3 CLI 發佈

`wails3 updater`命令群組涵蓋完整的發佈管線。發佈一個版本只需三個命令：

```bash
# Once per application: create the signing keypair.
wails3 updater genkey
# updater.key      keep secret (CI secret store), signs every release
# updater.key.pub  embed in the app and pass as updater.Config.PublicKey

# Per release: digest, sign and describe every artifact in one manifest.
wails3 updater manifest -version 2.1.0 -channel stable \
    -key updater.key -notes-file notes.md \
    -url-prefix "https://cdn.example.com/myapp/2.1.0" \
    bin/updates/

# Before uploading: re-verify the files exactly as a shipped app would.
wails3 updater verify -manifest manifest.json -publickey updater.key.pub
```

`manifest`可接受檔案或目錄（會自動略過金鑰材料、`.json`、總和檢查碼及備註附屬檔案），以串流方式對每個成品執行 SHA-512，在提供`-key`時使用 Ed25519ph 簽署摘要，並從`MyApp-2.1.0-darwin-arm64.zip`等慣用檔名推斷`platform`和`arch`（可識別`macOS`、`win64`、`x86_64`及`aarch64`等常見別名；無法推斷時會顯示警告，該成品的匹配範圍將不受平台限制）。省略`-url-prefix`即可輸出相對 URL，並將資訊清單上傳至與成品相同的目錄。

只要有任何不相符，`verify`就會以非零狀態碼結束，因此很適合作為建置與發佈之間的 CI 閘門。對於自行組合資訊清單的伺服器，`wails3 updater sign -key updater.key <files...>`會將每個檔案的`digest`/`signature`欄位輸出為 JSON，以便直接合併至您自己的文件。

## 驗證身分

驗證身分是伺服器的責任；通訊協定只負責傳遞標頭。用戶端會在每次請求資訊清單時重新傳送已設定的標頭。下載成品時，只有在成品 URL 與資訊清單位於相同主機，且未從`https`降級至`http`的情況下，才會傳送`Authorization`標頭；發生任何跨來源或降級重新導向時都會移除該標頭，因此憑證絕不會洩漏給 CDN 或物件儲存服務，也絕不會以明文傳輸。

受授權限制的範例，可自然地與託管式授權服務搭配使用：

```go
ep, _ := endpoint.New(endpoint.Config{
    URL:     "https://updates.example.com/check",
    Headers: map[string]string{"Authorization": "License " + licenseKey},
})
```

## 用戶端設定

請參閱[更新程式指南](/guides/updater/#providers)，以取得完整的`endpoint.Config`參考資料，並瞭解如何將此提供者與 GitHub、keygen.sh 及 AppCast 提供者一同納入後援鏈。
