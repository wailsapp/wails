---
title: "Instance unique"
description: "Limiter votre application à une seule instance en cours d’exécution"
slug: "guides/single-instance"
sourcePath: "guides/single-instance.md"
---

Le verrouillage d’instance unique est un mécanisme qui empêche l’exécution simultanée de plusieurs instances de votre application. Il est utile pour les applications conçues pour ouvrir des fichiers depuis la ligne de commande ou l’explorateur de fichiers du système d’exploitation.

## Utilisation

Pour activer la fonctionnalité d’instance unique dans votre application, fournissez une structure `SingleInstanceOptions` lors de sa création :

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

La structure `SingleInstanceOptions` comporte les champs suivants :

- `UniqueID` : identifiant unique de votre application. Il devrait s’agir d’une chaîne unique, généralement au format de nom de domaine inversé (par exemple, "com.company.appname").
- `EncryptionKey` : tableau facultatif de 32 octets permettant de chiffrer avec AES-256-GCM les données transmises entre les instances. Si vous fournissez un tableau non nul, toutes les communications entre les instances seront chiffrées.
- `OnSecondInstanceLaunch` : fonction de rappel appelée lorsqu’une deuxième instance de votre application est lancée. Cette fonction reçoit une structure `SecondInstanceData` contenant :
  - `Args` : les arguments de ligne de commande transmis à la deuxième instance
  - `WorkingDir` : le répertoire de travail de la deuxième instance
  - `AdditionalData` : toute donnée supplémentaire transmise par la deuxième instance (si elle est fournie)

- `AdditionalData` : table de paires clé-valeur de chaînes facultative, transmise à la première instance lors du lancement d’instances ultérieures

@note{type="danger" title="Avertissement"}
La fonctionnalité d’instance unique met en œuvre un protocole de chiffrement facultatif utilisant AES-256-GCM. Si le chiffrement n’est pas activé, les données transmises entre les instances ne sont pas sécurisées. Lorsque vous utilisez la fonctionnalité d’instance unique sans chiffrement, votre application devrait considérer comme non fiables toutes les données qui lui sont transmises par la fonction de rappel de la deuxième instance. Vous devriez vérifier que les arguments reçus sont valides et qu’ils ne contiennent aucune donnée malveillante.

@end

### Communication sécurisée

Pour sécuriser la communication entre les instances, fournissez une clé de chiffrement de 32 octets. Cette clé doit être identique pour toutes les instances de votre application :

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

@note{type="tip" title="Bonnes pratiques de sécurité"}
- Utilisez une clé propre à votre application
- Stockez la clé de manière sécurisée si vous la chargez depuis la configuration
- N’utilisez pas la clé d’exemple indiquée ci-dessus : créez la vôtre !

@end

### Gestion des fenêtres

Lorsque vous gérez le lancement d’une deuxième instance, vous souhaiterez souvent ramener la fenêtre de votre application au premier plan. Pour ce faire, utilisez la méthode `Focus()` de la fenêtre. Si celle-ci est réduite, vous devrez peut-être d’abord la restaurer :

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

## Fonctionnement

@tabs{sync-key="platform"}
[Mac]
Le verrouillage d’instance unique utilise un mutex nommé. Le nom du mutex est généré à partir de l’identifiant unique que vous fournissez. Les données sont transmises à la première instance via [NSDistributedNotificationCenter](https://developer.apple.com/documentation/foundation/nsdistributednotificationcenter).

[Windows]
Le verrouillage d’instance unique utilise un mutex nommé. Le nom du mutex est généré à partir de l’identifiant unique que vous fournissez. Les données sont transmises à la première instance via une fenêtre partagée à l’aide de [SendMessage](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage).

[Linux]
Le verrouillage d’instance unique utilise [dbus](https://www.freedesktop.org/wiki/Software/dbus/). Le nom dbus est généré à partir de l’identifiant unique que vous fournissez. Les données sont transmises à la première instance via [dbus](https://www.freedesktop.org/wiki/Software/dbus/).

@end
