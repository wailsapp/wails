---
title: "단일 인스턴스"
description: "앱의 실행 인스턴스를 하나로 제한하기"
slug: "guides/single-instance"
sourcePath: "guides/single-instance.md"
---

단일 인스턴스 잠금은 앱의 여러 인스턴스가 동시에 실행되지 않도록 하는 메커니즘입니다. 명령줄이나 OS 파일 탐색기에서 파일을 열도록 설계된 앱에 유용합니다.

## 사용법

앱에서 단일 인스턴스 기능을 활성화하려면 애플리케이션을 생성할 때 `SingleInstanceOptions` 구조체를 제공하세요:

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

`SingleInstanceOptions` 구조체에는 다음 필드가 있습니다:

- `UniqueID`: 애플리케이션의 고유 식별자입니다. 일반적으로 역방향 도메인 표기법(예: "com.company.appname")을 사용하는 고유 문자열을 권장합니다.
- `EncryptionKey`: AES-256-GCM을 사용해 인스턴스 간에 전달되는 데이터를 암호화하기 위한 선택적 32바이트 배열입니다. 0이 아닌 배열을 제공하면 인스턴스 간의 모든 통신이 암호화됩니다.
- `OnSecondInstanceLaunch`: 앱의 두 번째 인스턴스가 실행될 때 호출되는 콜백 함수입니다. 콜백은 다음 항목이 포함된 `SecondInstanceData` 구조체를 받습니다:
  - `Args`: 두 번째 인스턴스에 전달된 명령줄 인수
  - `WorkingDir`: 두 번째 인스턴스의 작업 디렉터리
  - `AdditionalData`: 두 번째 인스턴스에서 전달된 추가 데이터(제공된 경우)

- `AdditionalData`: 후속 인스턴스가 실행될 때 첫 번째 인스턴스에 전달할 선택적 문자열 키-값 쌍 맵

@note{type="danger" title="경고"}
단일 인스턴스 기능은 AES-256-GCM을 사용하는 선택적 암호화 프로토콜을 구현합니다. 암호화를 활성화하지 않으면 인스턴스 간에 전달되는 데이터는 안전하지 않습니다. 암호화 없이 단일 인스턴스 기능을 사용할 때는 두 번째 인스턴스의 콜백에서 앱으로 전달되는 모든 데이터를 신뢰할 수 없는 데이터로 취급하는 것이 좋습니다. 수신한 인수가 유효하고 악성 데이터를 포함하지 않는지 확인하는 것이 좋습니다.

@end

### 보안 통신

인스턴스 간 보안 통신을 활성화하려면 32바이트 암호화 키를 제공하세요. 이 키는 애플리케이션의 모든 인스턴스에서 동일해야 합니다:

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

@note{type="tip" title="보안 모범 사례"}
- 애플리케이션에 고유한 키를 사용하세요
- 구성에서 키를 불러오는 경우 안전하게 저장하세요
- 위에 표시된 예제 키를 사용하지 말고 직접 생성하세요!

@end

### 창 관리

두 번째 인스턴스의 실행을 처리할 때는 애플리케이션 창을 맨 앞으로 가져와야 하는 경우가 많습니다. 창의 `Focus()` 메서드를 사용하면 됩니다. 창이 최소화되어 있다면 먼저 복원해야 할 수 있습니다:

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

## 작동 방식

@tabs{sync-key="platform"}
[Mac]
명명된 뮤텍스를 사용해 단일 인스턴스를 잠급니다. 뮤텍스 이름은 사용자가 제공한 고유 ID로부터 생성됩니다. 데이터는 [NSDistributedNotificationCenter](https://developer.apple.com/documentation/foundation/nsdistributednotificationcenter)를 통해 첫 번째 인스턴스에 전달됩니다.

[Windows]
명명된 뮤텍스를 사용해 단일 인스턴스를 잠급니다. 뮤텍스 이름은 사용자가 제공한 고유 ID로부터 생성됩니다. 데이터는 [SendMessage](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage)를 사용하는 공유 창을 통해 첫 번째 인스턴스에 전달됩니다.

[Linux]
[dbus](https://www.freedesktop.org/wiki/Software/dbus/)를 사용해 단일 인스턴스를 잠급니다. dbus 이름은 사용자가 제공한 고유 ID로부터 생성됩니다. 데이터는 [dbus](https://www.freedesktop.org/wiki/Software/dbus/)를 통해 첫 번째 인스턴스에 전달됩니다.

@end
