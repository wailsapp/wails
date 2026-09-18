---
title: "Einzelinstanz"
description: "Beschränken der App auf eine einzige laufende Instanz"
slug: "guides/single-instance"
sourcePath: "guides/single-instance.md"
---

Die Einzelinstanzsperre verhindert, dass mehrere Instanzen Ihrer App gleichzeitig ausgeführt werden. Sie ist für Apps nützlich, die Dateien über die Befehlszeile oder den Dateimanager des Betriebssystems öffnen sollen.

## Verwendung

Um die Einzelinstanzfunktion in Ihrer App zu aktivieren, übergeben Sie beim Erstellen Ihrer Anwendung eine `SingleInstanceOptions`-Struktur:

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

Die `SingleInstanceOptions`-Struktur enthält die folgenden Felder:

- `UniqueID`: Eine eindeutige Kennung für Ihre Anwendung. Dies sollte eine eindeutige Zeichenfolge sein, üblicherweise in umgekehrter Domainnotation (z. B. „com.company.appname“).
- `EncryptionKey`: Optionales 32-Byte-Array zum Verschlüsseln der zwischen Instanzen übertragenen Daten mit AES-256-GCM. Wenn Sie ein Array angeben, das nicht ausschließlich Nullwerte enthält, wird die gesamte Kommunikation zwischen den Instanzen verschlüsselt.
- `OnSecondInstanceLaunch`: Eine Callback-Funktion, die beim Starten einer zweiten Instanz Ihrer App aufgerufen wird. Der Callback erhält eine `SecondInstanceData`-Struktur mit folgenden Daten:
  - `Args`: Die an die zweite Instanz übergebenen Befehlszeilenargumente
  - `WorkingDir`: Das Arbeitsverzeichnis der zweiten Instanz
  - `AdditionalData`: Alle zusätzlichen Daten, die von der zweiten Instanz übergeben wurden (sofern vorhanden)

- `AdditionalData`: Optionale Map mit Schlüssel-Wert-Paaren aus Zeichenfolgen, die beim Starten weiterer Instanzen an die erste Instanz übergeben wird

@note{type="danger" title="Warnung"}
Die Einzelinstanzfunktion implementiert ein optionales Verschlüsselungsprotokoll mit AES-256-GCM. Ohne aktivierte Verschlüsselung sind die zwischen den Instanzen übertragenen Daten nicht sicher. Wenn Sie die Einzelinstanzfunktion ohne Verschlüsselung verwenden, sollte Ihre App alle Daten, die ihr über den Callback der zweiten Instanz übergeben werden, als nicht vertrauenswürdig behandeln. Prüfen Sie, ob die empfangenen Argumente gültig sind und keine schädlichen Daten enthalten.

@end

### Sichere Kommunikation

Um die sichere Kommunikation zwischen Instanzen zu aktivieren, geben Sie einen 32 Byte langen Verschlüsselungsschlüssel an. Dieser Schlüssel muss für alle Instanzen Ihrer Anwendung identisch sein:

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

@note{type="tip" title="Bewährte Sicherheitspraktiken"}
- Verwenden Sie einen eindeutigen Schlüssel für Ihre Anwendung
- Speichern Sie den Schlüssel sicher, wenn Sie ihn aus der Konfiguration laden
- Verwenden Sie nicht den oben gezeigten Beispielschlüssel – erstellen Sie einen eigenen!

@end

### Fensterverwaltung

Beim Verarbeiten des Starts einer zweiten Instanz möchten Sie das Anwendungsfenster häufig in den Vordergrund holen. Verwenden Sie dazu die Methode `Focus()` des Fensters. Wenn Ihr Fenster minimiert ist, müssen Sie es möglicherweise zuerst wiederherstellen:

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

## Funktionsweise

@tabs{sync-key="platform"}
[Mac]
Die Einzelinstanzsperre verwendet einen benannten Mutex. Der Mutexname wird aus der von Ihnen angegebenen eindeutigen ID erzeugt. Die Daten werden über [NSDistributedNotificationCenter](https://developer.apple.com/documentation/foundation/nsdistributednotificationcenter) an die erste Instanz übergeben.

[Windows]
Die Einzelinstanzsperre verwendet einen benannten Mutex. Der Mutexname wird aus der von Ihnen angegebenen eindeutigen ID erzeugt. Die Daten werden über ein gemeinsam verwendetes Fenster mithilfe von [SendMessage](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage) an die erste Instanz übergeben.

[Linux]
Die Einzelinstanzsperre verwendet [dbus](https://www.freedesktop.org/wiki/Software/dbus/). Der dbus-Name wird aus der von Ihnen angegebenen eindeutigen ID erzeugt. Die Daten werden über [dbus](https://www.freedesktop.org/wiki/Software/dbus/) an die erste Instanz übergeben.

@end
