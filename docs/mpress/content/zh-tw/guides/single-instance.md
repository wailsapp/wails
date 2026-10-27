---
title: "單一執行個體"
description: "限制應用程式只能執行單一執行個體"
slug: "guides/single-instance"
sourcePath: "guides/single-instance.md"
---

單一執行個體鎖定是一種防止應用程式的多個執行個體同時執行的機制。 這對於設計為從命令列或作業系統檔案總管開啟檔案的應用程式很實用。

## 使用方式

若要在應用程式中啟用單一執行個體功能，請在建立應用程式時提供一個 `SingleInstanceOptions` 結構：

```go
app := application.New(application.Options{
    // ... other options ...
    SingleInstance: &application.SingleInstanceOptions{
        UniqueID: "com.myapp.unique-id",
        OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
            log.Printf("Second instance launched with args: %v", data.Args)
            log.Printf("Working directory: %s", data.WorkingDir)
            log.Printf("Additional data: %v", data.AdditionalData)
        },
        // Optional: Pass additional data to second instance
        AdditionalData: map[string]string{
            "launchtime": time.Now().String(),
        },
    },
})
```

`SingleInstanceOptions` 結構包含下列欄位：

- `UniqueID`：應用程式的唯一識別碼。此值應為唯一字串，通常採用反向網域名稱表示法（例如 "com.company.appname"）。
- `EncryptionKey`：選用的 32 位元組陣列，用於透過 AES-256-GCM 加密執行個體之間傳遞的資料。若提供非全零陣列，執行個體之間的所有通訊都會加密。
- `OnSecondInstanceLaunch`：啟動應用程式的第二個執行個體時呼叫的回呼函式。此回呼會收到一個 `SecondInstanceData` 結構，其中包含：
  - `Args`：傳遞給第二個執行個體的命令列引數
  - `WorkingDir`：第二個執行個體的工作目錄
  - `AdditionalData`：從第二個執行個體傳遞的任何額外資料（若有提供）

- `AdditionalData`：選用的字串鍵值配對映射，啟動後續執行個體時會將其傳遞給第一個執行個體

@note{type="danger" title="警告"}
單一執行個體功能使用 AES-256-GCM 實作選用的加密通訊協定。若未啟用加密， 執行個體之間傳遞的資料並不安全。在未加密的情況下使用單一執行個體功能時， 應用程式應將第二個執行個體回呼所傳遞的任何資料視為不受信任。 您應驗證收到的引數有效，且不含任何惡意資料。

@end

### 安全通訊

若要啟用執行個體之間的安全通訊，請提供一個 32 位元組的加密金鑰。應用程式的所有執行個體都必須使用相同的金鑰：

```go
// Define your encryption key (must be exactly 32 bytes)
var encryptionKey = [32]byte{
    0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
    0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
    0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
    0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
}

// Use the key in SingleInstanceOptions
SingleInstance: &application.SingleInstanceOptions{
    UniqueID: "com.myapp.unique-id",
    // Enable encryption for instance communication
    EncryptionKey: encryptionKey,
    // ... other options ...
}
```

@note{type="tip" title="安全性最佳實務"}
- 為應用程式使用專屬金鑰
- 若從組態載入金鑰，請安全地儲存金鑰
- 請勿使用上方所示的範例金鑰，請自行建立金鑰！

@end

### 視窗管理

處理第二個執行個體的啟動時，通常需要將應用程式視窗移至最前方。您可以使用該視窗的 `Focus()` 方法來執行此操作。若視窗已最小化，可能需要先將其還原：

```go

    var mainWindow *application.WebviewWindow

    SingleInstance: &application.SingleInstanceOptions{
        // Other options...
        OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
            // Focus the window if needed
            if mainWindow != nil {
                mainWindow.Restore()
                mainWindow.Focus()
            }
        },
    }
```

## 運作方式

@tabs{sync-key="platform"}
[Mac]
使用具名互斥鎖進行單一執行個體鎖定。互斥鎖名稱會根據您提供的唯一識別碼產生。資料透過 [NSDistributedNotificationCenter](https://developer.apple.com/documentation/foundation/nsdistributednotificationcenter) 傳遞給第一個執行個體。

[Windows]
使用具名互斥鎖進行單一執行個體鎖定。互斥鎖名稱會根據您提供的唯一識別碼產生。資料透過使用 [SendMessage](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage) 的共用視窗傳遞給第一個執行個體。

[Linux]
使用 [dbus](https://www.freedesktop.org/wiki/Software/dbus/) 進行單一執行個體鎖定。dbus 名稱會根據您提供的唯一識別碼產生。資料透過 [dbus](https://www.freedesktop.org/wiki/Software/dbus/) 傳遞給第一個執行個體。

@end
