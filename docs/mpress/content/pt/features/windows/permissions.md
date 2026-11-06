---
title: "Permissões"
description: "Controle solicitações de acesso à câmera, ao microfone, à geolocalização e a outros recursos feitas pelo conteúdo web"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

O conteúdo web que chama `navigator.mediaDevices.getUserMedia()`, a API de Geolocalização ou a API de Notificações precisa que o aplicativo host conceda ou negue essas solicitações. O Wails disponibiliza um mapa multiplataforma `Permissions` em `WebviewWindowOptions`, que permite controlar isso de forma declarativa, sem exigir código específico para cada plataforma.

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

As solicitações de acesso à câmera e ao microfone feitas pelo conteúdo web dessa janela são concedidas sem exibir uma solicitação de permissão do navegador.

## Tipos de permissão

`PermissionType` (uint8) identifica um recurso que o conteúdo web pode solicitar.

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

`PermissionDefault` é o valor zero, portanto as entradas não definidas no mapa adotam o comportamento padrão.

## Comportamento nas plataformas

Cada plataforma trata `PermissionDefault` de forma diferente porque suas webviews subjacentes têm comportamentos nativos distintos.

### Linux (WebKitGTK)

O WebKitGTK **não tem uma solicitação de permissão nativa**. Sem um manipulador associado, ele nega silenciosamente todas as solicitações — por isso `getUserMedia` sempre retornava `NotAllowedError` antes da adição deste recurso.

Atualmente, no Linux, o Wails trata solicitações de acesso à **câmera e ao microfone**. A geolocalização, as notificações e a leitura da área de transferência ainda não estão integradas e continuam sendo negadas, independentemente da política definida.

| Política | Câmera/microfone | Geolocalização, notificações e área de transferência |
| --- | --- | --- |
| `PermissionDefault` | **Permitido** (restaura getUserMedia) | Sempre negado |
| `PermissionAllow` | Permitido | Sempre negado (ainda não implementado) |
| `PermissionDeny` | Negado | Sempre negado |

### Windows (WebView2)

O WebView2 tem uma solicitação de permissão nativa e uma API de permissões por tipo. Os cinco tipos de recurso têm suporte completo.

| Política | Comportamento |
| --- | --- |
| `PermissionDefault` | O WebView2 exibe sua solicitação de permissão nativa do sistema operacional/navegador |
| `PermissionAllow` | Concedido silenciosamente |
| `PermissionDeny` | Negado silenciosamente |

**Importante:** antes da existência deste recurso, o Wails chamava `SetGlobalPermission(Allow)` incondicionalmente, concedendo silenciosamente todos os recursos. Agora, quando há qualquer entrada em `Permissions`, essa concessão geral **não** é definida. Para os recursos não definidos, o WebView2 exibe sua solicitação nativa em vez de concedê-los automaticamente.

Isso significa que, se você configurar `Permissions` de qualquer forma no Windows, qualquer recurso que não esteja explicitamente listado exibirá uma solicitação de permissão, em vez de ser permitido silenciosamente. Defina explicitamente os recursos necessários.

### macOS (WKWebView + TCC)

No macOS, duas camadas precisam estar de acordo. O WKWebView consulta o aplicativo antes de iniciar uma sessão de captura, e o mapa `Permissions` responde a essa solicitação. Abaixo dele, a estrutura de privacidade do sistema (TCC) controla o acesso ao próprio dispositivo: a solicitação do sistema operacional aparece na primeira vez que o aplicativo realmente tenta acessar a câmera ou o microfone, e a escolha do usuário é memorizada por aplicativo em Ajustes do Sistema → Privacidade e Segurança.

O Wails trata solicitações de acesso à **câmera e ao microfone** no macOS 12 e versões posteriores. A geolocalização, as notificações e a leitura da área de transferência não têm um equivalente a `WKUIDelegate`, portanto não estão integradas e ficam exclusivamente a cargo do TCC — assim como no Linux, a política definida para elas não tem efeito.

| Política | Câmera/microfone | Geolocalização, notificações e área de transferência |
| --- | --- | --- |
| `PermissionDefault` | O WebKit exibe sua própria solicitação de permissão | Somente TCC |
| `PermissionAllow` | A solicitação do WebKit é ignorada — **o TCC ainda se aplica** | Somente TCC |
| `PermissionDeny` | Negado pelo WebKit; a autorização do sistema é separada | Somente TCC |

`PermissionAllow` concede a solicitação da webview, não o acesso ao dispositivo. A autorização do sistema continua sendo aplicada de forma independente, e um aplicativo cujo acesso foi negado nos Ajustes do Sistema permanece sem acesso. Somente a solicitação do próprio WebKit é ignorada.

Em versões anteriores ao macOS 12, o método delegado não existe; portanto, o mapa é ignorado e todas as solicitações recorrem à solicitação de permissão do WebKit.

Verifique se o seu `Info.plist` inclui as chaves de descrição de uso apropriadas:

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="As descrições de uso são obrigatórias"}

O WebKit pode solicitar autorização do sistema antes de chamar o delegado de permissões do Wails. `PermissionDeny` bloqueia a solicitação do WebKit, mas não garante que nenhuma solicitação do TCC ou de autorização do AVFoundation seja exibida.

Inclua as descrições de uso apropriadas para os dispositivos solicitados pelo aplicativo e, quando ele for executado em sandbox, os direitos de acesso correspondentes aos dispositivos. Não use `PermissionDeny` como substituto dessas declarações.

@end

## Padrões comuns

### Aplicativo de captura de mídia

Conceda acesso à câmera e ao microfone em todas as plataformas:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

No **Linux**, isso permite explicitamente os dois dispositivos; os demais recursos permanecem negados. No **Windows**, isso permite ambos; qualquer outro recurso não listado exibirá uma solicitação nativa. No **macOS**, isso permite ambos na camada do WebKit, portanto nenhuma solicitação do navegador é exibida; o TCC ainda solicita acesso aos próprios dispositivos no primeiro uso, e as duas chaves de descrição de uso devem estar presentes em `Info.plist`.

### Negar captura de mídia no Linux

Por padrão, o Linux permite o acesso à câmera e ao microfone. Para desativá-lo:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### Permitir todos os recursos no Windows

Para conceder todos os recursos sem exibir solicitações:

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

O campo `Windows.Permissions` por janela (`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`) continua funcionando e pode substituir recursos individuais depois que o mapa multiplataforma é aplicado. Use-o quando precisar acessar tipos de permissão do WebView2 que não tenham equivalente multiplataforma (por exemplo, `CoreWebView2PermissionKindOtherSensors`).

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

A ordem de avaliação no Windows é:

1. Mapa multiplataforma `Permissions` (define o estado por tipo por meio de `SetPermission`)
2. Mapa `Windows.Permissions` (substitui tipos individuais)
3. Para qualquer tipo não abrangido por nenhum dos mapas: solicitação nativa do WebView2 (quando uma política está configurada) ou permissão automática (quando nenhuma política está configurada — o comportamento legado)

## Matriz de compatibilidade por plataforma

| Recurso | Linux | Windows | macOS |
| --- | --- | --- | --- |
| Microfone | ✅ | ✅ | ✅ (macOS 12+) |
| Câmera | ✅ | ✅ | ✅ (macOS 12+) |
| Geolocalização | ❌ ainda não | ✅ | ❌ ainda não |
| Notificações | ❌ ainda não | ✅ | ❌ ainda não |
| Leitura da área de transferência | ❌ ainda não | ✅ | ❌ ainda não |

Nos casos marcados com ✅ no macOS, a política responde à solicitação do WebKit; além disso, o TCC ainda controla o acesso ao dispositivo. Nos casos marcados com ❌, o recurso fica sujeito somente ao TCC.

## Solução de problemas

**`getUserMedia` continua falhando no Linux após a atualização**

Verifique se você não definiu explicitamente `PermissionMicrophone: PermissionDeny` nem `PermissionCamera: PermissionDeny`. O padrão (não definido) permite a captura de mídia no Linux.

**O Windows está solicitando permissões que não configurei**

Assim que houver qualquer entrada em `Permissions`, o Wails deixará de definir a concessão geral `Allow`. Os recursos não listados exibirão a solicitação nativa do WebView2. Adicione entradas `PermissionAllow` explícitas para todos os recursos usados pelo seu aplicativo.

**As permissões do macOS não estão funcionando**

`Permissions` abrange a câmera e o microfone no macOS 12 e posteriores; a geolocalização, as notificações e a leitura da área de transferência ainda não estão integradas e ignoram a política. Quando ela é respeitada, o TCC ainda controla adicionalmente o acesso ao dispositivo: `PermissionAllow` remove a solicitação do WebKit, mas não a do sistema. Verifique se seu `Info.plist` inclui as chaves corretas de descrição de uso (`NSMicrophoneUsageDescription`, `NSCameraUsageDescription`) e se o usuário concedeu acesso em Ajustes do Sistema → Privacidade e Segurança.

**Meu aplicativo para macOS é encerrado quando o conteúdo da Web solicita acesso à câmera ou ao microfone**

Verifique se `Info.plist` inclui `NSCameraUsageDescription` ou `NSMicrophoneUsageDescription` para cada dispositivo solicitado pelo aplicativo. A ausência das descrições pode fazer com que a solicitação seja rejeitada; acessar o AVFoundation sem uma descrição obrigatória pode encerrar o aplicativo. `PermissionDeny` controla a decisão do WebKit e não substitui essas chaves.

**Geolocalização, notificações e área de transferência não têm efeito no Linux**

Atualmente, apenas a câmera e o microfone são tratados no Linux. A compatibilidade com outros tipos de recurso ainda não foi implementada — eles permanecem negados independentemente da política definida.
