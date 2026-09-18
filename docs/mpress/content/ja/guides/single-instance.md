---
title: "シングルインスタンス"
description: "アプリの実行中のインスタンスを1つに制限する"
slug: "guides/single-instance"
sourcePath: "guides/single-instance.md"
---

シングルインスタンスロックは、アプリの複数のインスタンスが同時に実行されるのを防ぐ仕組みです。 コマンドラインやOSのファイルエクスプローラーからファイルを開くように設計されたアプリに役立ちます。

## 使用方法

アプリでシングルインスタンス機能を有効にするには、アプリケーションの作成時に `SingleInstanceOptions` 構造体を指定します。

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

`SingleInstanceOptions` 構造体には、次のフィールドがあります。

- `UniqueID`：アプリケーションの一意の識別子です。一意の文字列を指定してください。通常は逆ドメイン表記を使用します（例："com.company.appname"）。
- `EncryptionKey`：AES-256-GCMを使用してインスタンス間で渡されるデータを暗号化するための、省略可能な 32 バイト配列です。ゼロ以外の配列を指定すると、インスタンス間のすべての通信が暗号化されます。
- `OnSecondInstanceLaunch`：アプリの2つ目のインスタンスが起動されたときに呼び出されるコールバック関数です。このコールバックは、次の情報を含む `SecondInstanceData` 構造体を受け取ります。
  - `Args`：2つ目のインスタンスに渡されたコマンドライン引数
  - `WorkingDir`：2つ目のインスタンスの作業ディレクトリ
  - `AdditionalData`：2つ目のインスタンスから渡された追加データ（指定されている場合）

- `AdditionalData`：後続のインスタンスが起動されたときに最初のインスタンスへ渡される、文字列のキーと値のペアからなる省略可能なマップ

@note{type="danger" title="警告"}
シングルインスタンス機能には、AES-256-GCMを使用する省略可能な暗号化プロトコルが実装されています。暗号化を有効にしない場合、 インスタンス間で渡されるデータは安全ではありません。暗号化せずにシングルインスタンス機能を使用する場合、 アプリは、2つ目のインスタンスのコールバックから渡されるすべてのデータを信頼できないものとして扱ってください。 受け取った引数が有効であり、悪意のあるデータが含まれていないことを確認してください。

@end

### セキュアな通信

インスタンス間のセキュアな通信を有効にするには、32 バイトの暗号化キーを指定します。このキーは、アプリケーションのすべてのインスタンスで同一でなければなりません。

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

@note{type="tip" title="セキュリティのベストプラクティス"}
- アプリケーション固有のキーを使用してください
- 設定からキーを読み込む場合は、安全に保管してください
- 上記のサンプルキーは使用せず、独自のキーを作成してください！

@end

### ウィンドウ管理

2つ目のインスタンスの起動を処理する際は、アプリケーションウィンドウを最前面に表示したいことがよくあります。これには、ウィンドウの `Focus()` メソッドを使用できます。ウィンドウが最小化されている場合は、先に元の状態に戻す必要があるかもしれません：

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

## 仕組み

@tabs{sync-key="platform"}
[Mac]
名前付きミューテックスを使用してシングルインスタンスをロックします。ミューテックス名は、指定した一意のIDから生成されます。データは [NSDistributedNotificationCenter](https://developer.apple.com/documentation/foundation/nsdistributednotificationcenter) を介して最初のインスタンスに渡されます。

[Windows]
名前付きミューテックスを使用してシングルインスタンスをロックします。ミューテックス名は、指定した一意のIDから生成されます。データは、[SendMessage](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage) を使用する共有ウィンドウを介して最初のインスタンスに渡されます。

[Linux]
[dbus](https://www.freedesktop.org/wiki/Software/dbus/) を使用してシングルインスタンスをロックします。dbus名は、指定した一意のIDから生成されます。データは [dbus](https://www.freedesktop.org/wiki/Software/dbus/) を介して最初のインスタンスに渡されます。

@end
