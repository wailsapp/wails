---
title: "檔案關聯"
description: "為您的 Wails 應用程式設定檔案關聯"
slug: "guides/file-associations"
sourcePath: "guides/file-associations.md"
---

適用平台：<span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

檔案關聯可讓您的應用程式在使用者開啟特定類型的檔案時進行處理。這對文字編輯器、圖片檢視器或任何使用特定檔案格式的應用程式尤其實用。本指南說明如何在 Wails v3 應用程式中實作檔案關聯。

## 概觀

Wails v3 目前在下列平台支援檔案關聯：

- Windows（NSIS 安裝程式套件）
- macOS（應用程式套件）

## 設定

檔案關聯是在專案的 `build` 目錄內，透過 `config.yml` 檔案進行設定。

### 基本設定

若要設定檔案關聯：

1. 開啟 `build/config.yml`
2. 在 `fileAssociations` 區段下新增檔案關聯
3. 執行 `wails3 update build-assets` 以更新建置資產
4. 在應用程式選項中設定 `FileAssociations` 欄位
5. 使用 `wails3 package` 封裝應用程式

以下是設定範例：

```yaml
fileAssociations:
  - ext: myapp
    name: MyApp Document
    description: MyApp Document File
    iconName: myappFileIcon
    role: Editor
  - ext: custom
    name: Custom Format
    description: Custom File Format
    iconName: customFileIcon
    role: Editor
```

### 設定屬性

| 屬性 | 說明 | 平台 |
| --- | --- | --- |
| ext | 不含開頭句點的副檔名（例如 `txt`） | 全部 |
| name | 檔案類型的顯示名稱 | 全部 |
| description | 顯示於檔案內容中的說明 | Windows |
| iconName | build 資料夾中圖示檔案的名稱（不含副檔名） | 全部 |
| role | 應用程式處理此檔案類型時的角色（例如 `Editor`、`Viewer`） | macOS |
| mimeType | 檔案的 MIME 類型（例如 `image/jpeg`） | macOS |

## 監聽檔案開啟事件

若要在應用程式中處理檔案開啟事件，可以監聽 `events.Common.ApplicationOpenedWithFile` 事件：

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
    })

    // Listen for files being used to open the application
    app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
        associatedFile := event.Context().Filename()
        app.Dialog.Info().SetMessage("Application opened with file: " + associatedFile).Show()
    })

    // Create your window and run the app...
}

```

## 逐步教學

接下來逐步為簡易文字編輯器設定檔案關聯：

@steps
### 建立圖示
- 為您的檔案類型建立圖示（建議尺寸：16x16、32x32、48x48、256x256）
- 將圖示儲存至專案的 `build` 資料夾
- 依照 `iconName` 設定為圖示命名（例如 `textFileIcon.png`）

@note{type="tip"}
您可以使用 `wails3 generate icons` 產生所需的圖示。執行 `wails3 generate icons --help` 以取得更多資訊。

@end

- 在 macOS 上，請在 `create:app:bundle:` 工作中加入類似 `cp build/darwin/documenticon.icns {{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources` 的複製陳述式。

### 設定檔案關聯
編輯 `build/config.yml` 檔案以新增檔案關聯：

```yaml
# build/config.yml
fileAssociations:
  - ext: txt
    name: Text Document
    description: Plain Text Document
    iconName: textFileIcon
    role: Editor
```

### 更新建置資產
執行下列命令以更新建置資產：

```bash
wails3 update build-assets
```

### 在應用程式選項中設定檔案關聯
在 `main.go` 檔案中，設定應用程式選項的 `FileAssociations` 欄位：

```go
app := application.New(application.Options{
  Name: "MyApp",
  FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
})
```

@note{type="tip" title="為什麼應用程式設定和 config.yml 都必須指定副檔名？"}
在 Windows 上，透過檔案關聯開啟檔案時，應用程式會啟動，並將檔案名稱作為應用程式的第一個引數。應用程式無法得知第一個引數是檔案還是命令列引數，因此會使用應用程式選項中的 `FileAssociations` 欄位，判斷第一個引數是否為已建立關聯的檔案。

@end

### 封裝應用程式
使用下列命令封裝應用程式：

```bash
wails3 package
```

封裝後的應用程式會建立於 `bin` 目錄中。接著即可安裝並測試應用程式。

## 其他注意事項

- 圖示應以 PNG 格式放置於 build 資料夾中
- 必須安裝封裝後的應用程式，才能測試檔案關聯

@end
