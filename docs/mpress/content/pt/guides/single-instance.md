---
title: "Instância única"
description: "Limitação do aplicativo a uma única instância em execução"
slug: "guides/single-instance"
sourcePath: "guides/single-instance.md"
---

O bloqueio de instância única é um mecanismo que impede a execução simultânea de várias instâncias do aplicativo. Ele é útil para aplicativos projetados para abrir arquivos pela linha de comando ou pelo explorador de arquivos do sistema operacional.

## Uso

Para habilitar a funcionalidade de instância única no aplicativo, forneça uma struct `SingleInstanceOptions` ao criar o aplicativo:

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

A struct `SingleInstanceOptions` tem os seguintes campos:

- `UniqueID`: um identificador exclusivo para o aplicativo. Deve ser uma string exclusiva, normalmente em notação de domínio reverso (por exemplo, "com.empresa.nomeaplicativo").
- `EncryptionKey`: array opcional de 32 bytes para criptografar, com AES-256-GCM, os dados transmitidos entre as instâncias. Se for fornecido como um array não zerado, toda a comunicação entre as instâncias será criptografada.
- `OnSecondInstanceLaunch`: uma função de callback chamada quando uma segunda instância do aplicativo é iniciada. O callback recebe uma struct `SecondInstanceData` que contém:
  - `Args`: os argumentos de linha de comando passados para a segunda instância
  - `WorkingDir`: o diretório de trabalho da segunda instância
  - `AdditionalData`: quaisquer dados adicionais passados pela segunda instância (se fornecidos)

- `AdditionalData`: mapa opcional de pares chave-valor de strings que será passado para a primeira instância quando instâncias subsequentes forem iniciadas

@note{type="danger" title="Aviso"}
O recurso de instância única implementa um protocolo de criptografia opcional que usa AES-256-GCM. Sem a criptografia habilitada, os dados transmitidos entre as instâncias não são seguros. Ao usar o recurso de instância única sem criptografia, o aplicativo deve tratar como não confiáveis todos os dados que receber pelo callback da segunda instância. Verifique se os argumentos recebidos são válidos e não contêm dados maliciosos.

@end

### Comunicação segura

Para habilitar a comunicação segura entre as instâncias, forneça uma chave de criptografia de 32 bytes. Essa chave deve ser a mesma para todas as instâncias do aplicativo:

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

@note{type="tip" title="Práticas recomendadas de segurança"}
- Use uma chave exclusiva para o aplicativo
- Armazene a chave com segurança caso ela seja carregada de uma configuração
- Não use a chave de exemplo mostrada acima — crie a sua própria!

@end

### Gerenciamento de janelas

Ao processar a inicialização de uma segunda instância, você frequentemente vai querer trazer a janela do aplicativo para a frente. Para isso, você pode usar o método `Focus()` da janela. Se a janela estiver minimizada, talvez seja necessário restaurá-la primeiro:

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

## Como funciona

@tabs{sync-key="platform"}
[Mac]
O bloqueio de instância única usa um mutex nomeado. O nome do mutex é gerado a partir do identificador exclusivo fornecido. Os dados são passados para a primeira instância por meio de [NSDistributedNotificationCenter](https://developer.apple.com/documentation/foundation/nsdistributednotificationcenter)

[Windows]
O bloqueio de instância única usa um mutex nomeado. O nome do mutex é gerado a partir do identificador exclusivo fornecido. Os dados são passados para a primeira instância por meio de uma janela compartilhada usando [SendMessage](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage)

[Linux]
O bloqueio de instância única usa [dbus](https://www.freedesktop.org/wiki/Software/dbus/). O nome do dbus é gerado a partir do identificador exclusivo fornecido. Os dados são passados para a primeira instância por meio de [dbus](https://www.freedesktop.org/wiki/Software/dbus/)

@end
