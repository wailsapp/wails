---
title: "单实例"
description: "将应用限制为仅运行一个实例"
slug: "guides/single-instance"
sourcePath: "guides/single-instance.md"
---

单实例锁是一种防止应用同时运行多个实例的机制。 它适用于设计为从命令行或操作系统文件资源管理器打开文件的应用。

## 用法

要在应用中启用单实例功能，请在创建应用时提供一个`SingleInstanceOptions`结构体：

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

`SingleInstanceOptions`结构体包含以下字段：

- `UniqueID`：应用的唯一标识符。它应为唯一字符串，通常采用反向域名表示法（例如“com.company.appname”）。
- `EncryptionKey`：可选的32字节数组，用于通过AES-256-GCM加密实例之间传递的数据。如果提供的数组不全为零，则实例之间的所有通信都会被加密。
- `OnSecondInstanceLaunch`：启动应用的第二个实例时调用的回调函数。该回调会接收一个`SecondInstanceData`结构体，其中包含：
  - `Args`：传递给第二个实例的命令行参数
  - `WorkingDir`：第二个实例的工作目录
  - `AdditionalData`：从第二个实例传递的任何附加数据（如果提供）

- `AdditionalData`：可选的字符串键值对映射，启动后续实例时会将其传递给第一个实例

@note{type="danger" title="警告"}
单实例功能使用AES-256-GCM实现了一种可选的加密协议。未启用加密时， 实例之间传递的数据并不安全。在不加密的情况下使用单实例功能时， 应用应将第二个实例回调传递给它的任何数据视为不可信数据。 你应验证收到的参数是否有效，并确保其中不含任何恶意数据。

@end

### 安全通信

要启用实例之间的安全通信，请提供一个32字节的加密密钥。应用的所有实例都必须使用相同的密钥：

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

@note{type="tip" title="安全最佳实践"}
- 为应用使用唯一密钥
- 如果从配置中加载密钥，请安全地存储密钥
- 不要使用上面所示的示例密钥——请创建自己的密钥！

@end

### 窗口管理

处理第二个实例的启动时，通常需要将应用窗口置于最前。可以使用窗口的`Focus()`方法执行此操作。如果窗口已最小化，可能需要先将其还原：

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

## 工作原理

@tabs{sync-key="platform"}
[Mac]
使用命名互斥体实现单实例锁。互斥体名称根据你提供的唯一ID生成。数据通过[NSDistributedNotificationCenter](https://developer.apple.com/documentation/foundation/nsdistributednotificationcenter)传递给第一个实例。

[Windows]
使用命名互斥体实现单实例锁。互斥体名称根据你提供的唯一ID生成。数据通过使用[SendMessage](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage)的共享窗口传递给第一个实例。

[Linux]
使用[dbus](https://www.freedesktop.org/wiki/Software/dbus/)实现单实例锁。dbus名称根据你提供的唯一ID生成。数据通过[dbus](https://www.freedesktop.org/wiki/Software/dbus/)传递给第一个实例。

@end
