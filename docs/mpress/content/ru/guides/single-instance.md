---
title: "Единственный экземпляр"
description: "Ограничение приложения одним запущенным экземпляром"
slug: "guides/single-instance"
sourcePath: "guides/single-instance.md"
---

Блокировка единственного экземпляра — это механизм, который предотвращает одновременный запуск нескольких экземпляров приложения. Он полезен для приложений, предназначенных для открытия файлов из командной строки или файлового менеджера ОС.

## Использование

Чтобы включить поддержку единственного экземпляра в приложении, передайте структуру `SingleInstanceOptions` при его создании:

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

Структура `SingleInstanceOptions` содержит следующие поля:

- `UniqueID`: уникальный идентификатор приложения. Рекомендуется использовать уникальную строку, обычно в формате обратной доменной записи (например, "com.company.appname").
- `EncryptionKey`: необязательный массив из 32 байт для шифрования передаваемых между экземплярами данных с помощью AES-256-GCM. Если передан массив, в котором хотя бы один байт отличен от нуля, весь обмен данными между экземплярами будет зашифрован.
- `OnSecondInstanceLaunch`: функция обратного вызова, которая вызывается при запуске второго экземпляра приложения. Она получает структуру `SecondInstanceData`, содержащую:
  - `Args`: аргументы командной строки, переданные второму экземпляру
  - `WorkingDir`: рабочий каталог второго экземпляра
  - `AdditionalData`: дополнительные данные, переданные из второго экземпляра (если они указаны)

- `AdditionalData`: необязательное отображение строковых пар «ключ — значение», которое будет передаваться первому экземпляру при запуске последующих экземпляров

@note{type="danger" title="Предупреждение"}
Функция единственного экземпляра реализует необязательный протокол шифрования с помощью AES-256-GCM. Если шифрование не включено, данные, передаваемые между экземплярами, не защищены. При использовании функции единственного экземпляра без шифрования приложению следует считать недоверенными любые данные, переданные ему через функцию обратного вызова второго экземпляра. Следует проверить, что полученные аргументы допустимы и не содержат вредоносных данных.

@end

### Защищённый обмен данными

Чтобы включить защищённый обмен данными между экземплярами, укажите ключ шифрования длиной 32 байт. Этот ключ должен быть одинаковым для всех экземпляров приложения:

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

@note{type="tip" title="Рекомендации по безопасности"}
- Используйте уникальный ключ для своего приложения
- Если ключ загружается из конфигурации, храните его безопасным образом
- Не используйте приведённый выше пример ключа — создайте собственный!

@end

### Управление окном

При обработке запуска второго экземпляра часто требуется вывести окно приложения на передний план. Для этого можно использовать метод окна `Focus()`. Если окно свёрнуто, возможно, сначала потребуется его восстановить:

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

## Принцип работы

@tabs{sync-key="platform"}
[Mac]
Для блокировки единственного экземпляра используется именованный мьютекс. Имя мьютекса формируется на основе указанного вами уникального идентификатора. Данные передаются первому экземпляру через [NSDistributedNotificationCenter](https://developer.apple.com/documentation/foundation/nsdistributednotificationcenter)

[Windows]
Для блокировки единственного экземпляра используется именованный мьютекс. Имя мьютекса формируется на основе указанного вами уникального идентификатора. Данные передаются первому экземпляру через общее окно с помощью [SendMessage](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage)

[Linux]
Для блокировки единственного экземпляра используется [dbus](https://www.freedesktop.org/wiki/Software/dbus/). Имя dbus формируется на основе указанного вами уникального идентификатора. Данные передаются первому экземпляру через [dbus](https://www.freedesktop.org/wiki/Software/dbus/)

@end
