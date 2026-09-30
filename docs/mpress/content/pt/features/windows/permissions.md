---
title: "Permissões"
description: "Controle solicitações de acesso à câmera, ao microfone, à geolocalização e a outros recursos feitas por conteúdo da Web"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

O conteúdo da Web que chama `navigator.mediaDevices.getUserMedia()`, a API de geolocalização ou a API de notificações precisa que o aplicativo host conceda ou negue essas solicitações. O Wails disponibiliza um mapa multiplataforma `Permissions` em `WebviewWindowOptions`, permitindo controlar isso de forma declarativa, sem a necessidade de código específico para cada plataforma.

## Início rápido

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})
```

As solicitações de acesso à câmera e ao microfone feitas pelo conteúdo da Web dessa janela são concedidas sem que o navegador exiba uma solicitação de permissão.

## Tipos de permissão

`PermissionType` (uint8) identifica um recurso que o conteúdo da Web pode solicitar.

| Constante | Recurso |
| --- | --- |
| `PermissionMicrophone` | `getUserMedia({audio: true})` |
| `PermissionCamera` | `getUserMedia({video: true})` |
| `PermissionGeolocation` | `navigator.geolocation` |
| `PermissionNotifications` | `Notification.requestPermission()` |
| `PermissionClipboardRead` | `navigator.clipboard.readText()` |

## Valores de permissão

`Permission` (uint8) é a política aplicada a um determinado tipo.

| Constante | Valor | Significado |
| --- | --- | --- |
| `PermissionDefault` | 0 | Usar o tratamento nativo da plataforma (veja abaixo) |
| `PermissionAllow` | 1 | Conceder sem solicitar confirmação |
| `PermissionDeny` | 2 | Negar sem solicitar confirmação |

`PermissionDefault` é o valor zero; portanto, as entradas não definidas no mapa usam o comportamento padrão.

## Comportamento nas plataformas

Cada plataforma trata `PermissionDefault` de maneira diferente porque suas webviews subjacentes têm comportamentos nativos distintos.

### Linux (WebKitGTK)

O WebKitGTK **não tem uma solicitação de permissão nativa**. Sem um manipulador associado, ele nega silenciosamente todas as solicitações — por isso, `getUserMedia` sempre retornava `NotAllowedError` antes da adição desse recurso.

Atualmente, no Linux, o Wails trata solicitações de acesso à **câmera e ao microfone**. A geolocalização, as notificações e a leitura da área de transferência ainda não estão integradas e continuam sendo negadas, independentemente da política definida.

| Política | Câmera/microfone | Geolocalização, notificações e área de transferência |
| --- | --- | --- |
| `PermissionDefault` | **Permitido** (restaura getUserMedia) | Sempre negado |
| `PermissionAllow` | Permitido | Sempre negado (ainda não implementado) |
| `PermissionDeny` | Negado | Sempre negado |

### Windows (WebView2)

O WebView2 tem uma solicitação de permissão nativa e uma API de permissões por tipo. Os cinco tipos de recurso são totalmente compatíveis.

| Política | Comportamento |
| --- | --- |
| `PermissionDefault` | O WebView2 exibe a solicitação de permissão nativa do sistema operacional/navegador |
| `PermissionAllow` | Concedido silenciosamente |
| `PermissionDeny` | Negado silenciosamente |

**Importante:** antes da existência desse recurso, o Wails chamava `SetGlobalPermission(Allow)` incondicionalmente, concedendo silenciosamente todos os recursos. Agora, quando há qualquer entrada em `Permissions`, essa concessão geral **não** é definida. Os recursos não definidos passam a exibir a solicitação nativa do WebView2, em vez de serem concedidos automaticamente.

Isso significa que, se você configurar `Permissions` de qualquer forma no Windows, todo recurso que não estiver listado explicitamente exibirá uma solicitação de permissão, em vez de ser permitido silenciosamente. Defina explicitamente os recursos necessários.

### macOS (TCC)

O macOS gerencia o acesso à câmera, ao microfone, à geolocalização e às notificações por meio de sua estrutura de privacidade do sistema. A solicitação do sistema operacional aparece automaticamente na primeira vez em que o conteúdo da Web solicita um recurso, e a escolha do usuário é lembrada para cada aplicativo em Ajustes do Sistema → Privacidade e Segurança.

Isso funciona corretamente sem nenhuma configuração de `Permissions`. Atualmente, o mapa **é ignorado no macOS** — todas as solicitações passam pelo TCC, independentemente do valor definido. Na prática, a limitação é que `PermissionDeny` não tem efeito no macOS: não é possível impedir que uma webview use um recurso que o TCC já tenha concedido no nível do sistema.

Verifique se `Info.plist` inclui as chaves de descrição de uso apropriadas:

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

## Padrões comuns

### Aplicativo de captura de mídia

Conceda acesso à câmera e ao microfone em todas as plataformas:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

No **Linux**, isso permite explicitamente ambos os dispositivos; os demais recursos continuam bloqueados. No **Windows**, isso permite ambos; qualquer outro recurso que você não listar exibirá uma solicitação nativa. No **macOS**, isso não tem efeito; o TCC gerencia tudo.

### Bloquear a captura de mídia no Linux

Por padrão, o Linux permite o acesso à câmera e ao microfone. Para desativá-lo:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### Permitir todos os recursos no Windows

Para permitir todos os recursos sem exibir solicitações:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone:    application.PermissionAllow,
    application.PermissionCamera:        application.PermissionAllow,
    application.PermissionGeolocation:   application.PermissionAllow,
    application.PermissionNotifications: application.PermissionAllow,
    application.PermissionClipboardRead: application.PermissionAllow,
},
```

### Políticas por janela

Janelas diferentes podem ter políticas diferentes:

```go
// Main app window — full media access
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})

// Settings window — no special capabilities
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Settings",
    // No Permissions entry — uses platform defaults
})

// Embedded content window — deny media capture
embeddedWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Embedded",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionDeny,
        application.PermissionCamera:     application.PermissionDeny,
    },
})
```

## Substituição específica do Windows

O campo `Windows.Permissions` por janela (`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`) continua funcionando e pode substituir as configurações de recursos individuais depois que o mapa multiplataforma é aplicado. Use-o quando precisar acessar tipos de permissão do WebView2 que não tenham um equivalente multiplataforma (por exemplo, `CoreWebView2PermissionKindOtherSensors`).

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

A ordem de avaliação no Windows é:

1. Mapa multiplataforma `Permissions` (define o estado de cada tipo por meio de `SetPermission`)
2. Mapa `Windows.Permissions` (substitui tipos individuais)
3. Para qualquer tipo não abrangido por nenhum dos mapas: solicitação nativa do WebView2 (quando uma política está configurada) ou permissão automática (quando nenhuma política está configurada — o comportamento legado)

## Matriz de compatibilidade entre plataformas

| Recurso | Linux | Windows | macOS |
| --- | --- | --- | --- |
| Microfone | ✅ | ✅ | Somente TCC |
| Câmera | ✅ | ✅ | Somente TCC |
| Geolocalização | ❌ ainda não | ✅ | Somente TCC |
| Notificações | ❌ ainda não | ✅ | Somente TCC |
| Leitura da área de transferência | ❌ ainda não | ✅ | Somente TCC |

## Solução de problemas

**`getUserMedia` continua falhando no Linux após a atualização**

Verifique se você não definiu explicitamente `PermissionMicrophone: PermissionDeny` nem `PermissionCamera: PermissionDeny`. O padrão (não definido) permite a captura de mídia no Linux.

**O Windows está solicitando permissões que eu não configurei**

Quando há alguma entrada em `Permissions`, o Wails deixa de definir a permissão geral `Allow`. Os recursos que você não listar exibirão a solicitação nativa do WebView2. Adicione entradas `PermissionAllow` explícitas para todos os recursos usados pelo seu aplicativo.

**As permissões do macOS não estão funcionando**

O mapa `Permissions` não tem efeito no macOS. Verifique se `Info.plist` inclui as chaves corretas de descrição de uso (`NSMicrophoneUsageDescription`, `NSCameraUsageDescription` etc.) e se o usuário concedeu acesso em Ajustes do Sistema → Privacidade e Segurança.

**Geolocalização, notificações e área de transferência não têm efeito no Linux**

Atualmente, somente a câmera e o microfone são tratados no Linux. A compatibilidade com outros tipos de recurso ainda não foi implementada — eles permanecem bloqueados independentemente da política definida.
