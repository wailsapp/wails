---
title: Single Instance
description: Limiting your app to a single running instance
sidebar:
   order: 40
---

import { Tabs, TabItem } from "@astrojs/starlight/components";

Single instance locking is a mechanism that prevents multiple instances of your app from running at the same time.
It is useful for apps that are designed to open files from the command line or from the OS file explorer.


## Usage

To enable single instance functionality in your app, provide a `SingleInstanceOptions` struct when creating your application:

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

The `SingleInstanceOptions` struct has the following fields:

- `UniqueID`: A unique identifier for your application. This should be a unique string, typically in reverse domain notation (e.g., "com.company.appname").
- `EncryptionKey`: Optional 32-byte array for encrypting data passed between instances using AES-256-GCM. If provided as a non-zero array, all communication between instances will be encrypted.
- `OnSecondInstanceLaunch`: A callback function that is called when a second instance of your app is launched. The callback receives a `SecondInstanceData` struct containing:
    - `Args`: The command line arguments passed to the second instance
    - `WorkingDir`: The working directory of the second instance
    - `AdditionalData`: Any additional data passed from the second instance (if provided)
- `AdditionalData`: Optional map of string key-value pairs that will be passed to the first instance when subsequent instances are launched

:::danger[Warning]
The Single Instance feature implements an optional encryption protocol using AES-256-GCM. Without encryption enabled,
data passed between instances is not secure. When using the single instance feature without encryption,
your app should treat any data passed to it from second instance callback as untrusted.
You should verify that args that you receive are valid and don't contain any malicious data.
:::

### Secure Communication

To enable secure communication between instances, provide a 32-byte encryption key. This key must be the same for all instances of your application:

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

:::tip[Security Best Practices]
- Use a unique key for your application
- Store the key securely if loading it from configuration
- Do not use the example key shown above - create your own!
:::

### Window Management

When handling second instance launches, you'll often want to bring your application window to the front. You can do this using the window's `Focus()` method. If your window is minimized, you may need to restore it first:

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

## How it works

<Tabs syncKey="platform">
    <TabItem label="Mac" icon="apple">

        Single instance lock using a named mutex. The mutex name is generated from the unique id that you provide. Data is passed to the first instance via [NSDistributedNotificationCenter](https://developer.apple.com/documentation/foundation/nsdistributednotificationcenter)

    </TabItem>
    <TabItem label="Windows" icon="seti:windows">

        Single instance lock using a named mutex. The mutex name is generated from the unique id that you provide. Data is passed to the first instance via a shared window using [SendMessage](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage)

    </TabItem>
    <TabItem label="Linux" icon="linux">

        Single instance lock using [dbus](https://www.freedesktop.org/wiki/Software/dbus/). The dbus name is generated from the unique id that you provide. Data is passed to the first instance via [dbus](https://www.freedesktop.org/wiki/Software/dbus/)

    </TabItem>
</Tabs>
