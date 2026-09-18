---
title: "Notificações"
description: "Exiba notificações nativas do sistema com botões de ação e entrada de texto"
slug: "features/notifications/overview"
sourcePath: "features/notifications/overview.md"
---

## Introdução

O Wails oferece um sistema completo de notificações multiplataforma para aplicações desktop. Esse serviço permite exibir notificações nativas do sistema, com suporte a:

- Notificações básicas com título, subtítulo e corpo
- Notificações interativas com botões de ação e respostas de texto
- [Categorias de notificação](#notificaes-interativas) reutilizáveis para ações
- [Sons](#som-personalizado) personalizados (padrão, silencioso ou nomeado)
- [Anexos](#anexos) (imagens em todas as plataformas; áudio/vídeo no macOS)
- [Agrupamento de notificações relacionadas](#encadeamento-e-agrupamento) por `ThreadID`
- [Prioridade](#nvel-de-interrupo) por meio de `InterruptionLevel` (`passive` / `active` / `timeSensitive` / `critical`)
- [Entrega agendada](#entrega-agendada) (nativa no macOS; temporizador no processo no Windows e no Linux)
- [Atualização de uma notificação em andamento](#atualizao-de-notificaes) por ID

Cada novo campo opcional tem uma degradação harmoniosa quando uma plataforma não consegue respeitá-lo; consulte [Considerações sobre plataformas](#consideraes-sobre-as-plataformas) para ver a matriz de suporte de cada recurso.

## Uso básico

### Criação do serviço

Primeiro, inicialize o serviço de notificações:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/notifications"

// Create a new notification service
notifier := notifications.New()

//Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(notifier),
    },
})
```

## Autorização de notificações

As notificações no macOS exigem autorização do usuário. Solicite e verifique a autorização:

```go
authorized, err := notifier.CheckNotificationAuthorization()
if err != nil {
    // Handle authorization error
}
if authorized {
    // Send notifications
} else {
    // Request authorization
    authorized, err = notifier.RequestNotificationAuthorization()
}
```

No Windows e no Linux, isso sempre retorna `true`.

## Tipos de notificação

### Notificações básicas

Envie aos usuários uma notificação básica com ID exclusivo, título, subtítulo opcional (macOS e Linux) e texto no corpo:

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
})

```

### Notificações interativas

Envie uma notificação com botões de ação e campos de entrada de texto. Essas notificações exigem que uma categoria de notificação seja registrada primeiro:

```go
// Define a unique category id
categoryID := "unique-category-id"

// Define a category with actions
category := notifications.NotificationCategory{
    ID: categoryID,
    Actions: []notifications.NotificationAction{
        {
            ID:    "OPEN", 
            Title: "Open",
        },
        {
            ID:          "ARCHIVE", 
            Title:       "Archive", 
            Destructive: true,  /* macOS-specific */
        },
    },
    HasReplyField:    true,
    ReplyPlaceholder: "message...",
    ReplyButtonTitle: "Reply",
}

// Register the category
notifier.RegisterNotificationCategory(category)

// Send an interactive notification with the actions registered in the provided category
notifier.SendNotificationWithActions(notifications.NotificationOptions{
    ID:         "unique-id",
    Title:      "New Message",
    Subtitle:   "From: Jane Doe",
    Body:       "Are you able to make it?",
    CategoryID: categoryID,
})
```

## Respostas às notificações

Processe as interações dos usuários com as notificações:

```go
notifier.OnNotificationResponse(func(result notifications.NotificationResult) {
    response := result.Response
    fmt.Printf("Notification %s was actioned with: %s\n", response.ID, response.ActionIdentifier)

    if response.ActionIdentifier == "TEXT_REPLY" {
        fmt.Printf("User replied: %s\n", response.UserText)
    }

    if data, ok := response.UserInfo["sender"].(string); ok {
        fmt.Printf("Original sender: %s\n", data)
    }

    // Emit an event to the frontend
    app.Event.Emit("notification", result.Response)
})
```

## Personalização de notificações

### Metadados personalizados

Notificações básicas e interativas podem incluir dados personalizados:

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
    Data: map[string]interface{}{
        "sender": "jane.doe@example.com",
        "timestamp": "2025-03-10T15:30:00Z",
    }
})
```

### Som personalizado

Use `Sound` para controlar o áudio reproduzido quando uma notificação é entregue. Deixá-lo como `nil` reproduz o som padrão da plataforma; defina `Silent: true` para suprimir o som; defina `Name` para reproduzir um som nomeado ou incluído no pacote.

```go
// Silent
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "silent-id",
    Title: "Background sync complete",
    Sound: &notifications.NotificationSound{Silent: true},
})

// Named sound
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "named-id",
    Title: "New message",
    Body:  "Tap to read",
    Sound: &notifications.NotificationSound{Name: "Ping"},
})
```

Como `Name` é resolvido em cada plataforma:

- **macOS** — `Name` é passado para `[UNNotificationSound soundNamed:]`; o arquivo de áudio deve estar no diretório `Library/Sounds` do pacote da aplicação.
- **Windows** — se `Name` já começar com `ms-winsoundevent:` ou `ms-appx:`, será usado como está; caso contrário, será encapsulado em `ms-winsoundevent:` para uso como nome de evento de notificação toast integrado (consulte a documentação da Microsoft sobre o esquema de `<audio>` das notificações toast).
- **Linux** — encaminhado como a dica `sound-name` do freedesktop; a reprodução depende do daemon de notificações e do tema de sons ativos.

### Anexos

`Attachments` adiciona arquivos de mídia a uma notificação. O macOS aceita vários anexos de qualquer tipo de mídia; o Windows e o Linux usam o primeiro anexo do tipo imagem.

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "image-id",
    Title: "Photo uploaded",
    Body:  "Tap to view",
    Attachments: []notifications.NotificationAttachment{
        {
            ID:   "preview",
            Path: "/absolute/path/to/image.png",
            // Type hint:
            //   macOS:   UTI like "public.png" / "public.audio" (often inferred)
            //   Windows: "hero" | "appLogoOverride" | "inline" (default "inline")
            //   Linux:   ignored
            Type: "inline",
        },
    },
})
```

`Path` deve ser um caminho absoluto do sistema de arquivos. O macOS também aceita URLs `file://`.

#### Como anexar um arquivo distribuído com a aplicação

O sistema operacional lê o anexo do disco quando a notificação é entregue, portanto `Path` precisa corresponder a um arquivo real na máquina do usuário final. Para um recurso incluído no pacote da aplicação (um ícone ou uma imagem incorporado com `go:embed`), não há um caminho absoluto fixo que possa ser definido diretamente no código, pois o arquivo fica dentro do binário, e não em um local conhecido no disco. Grave-o uma vez em um diretório gravável durante a inicialização e passe esse caminho:

```go
import (
    _ "embed"
    "os"
    "path/filepath"
)

//go:embed assets/preview.png
var previewPNG []byte

// Materialise the embedded asset to a stable path the OS can read.
previewPath := filepath.Join(os.TempDir(), "myapp-preview.png")
if err := os.WriteFile(previewPath, previewPNG, 0o644); err != nil {
    // handle error
}

notifier.SendNotification(notifications.NotificationOptions{
    ID:    "image-id",
    Title: "Photo uploaded",
    Body:  "Tap to view",
    Attachments: []notifications.NotificationAttachment{
        {ID: "preview", Path: previewPath, Type: "inline"},
    },
})
```

Um arquivo fornecido pelo usuário ou baixado já possui um caminho real no disco; portanto, você pode passá-lo diretamente para `Path` sem essa etapa. O suporte para passar os bytes do anexo na memória poderá ser adicionado em uma versão futura.

### Encadeamento e agrupamento

`ThreadID` agrupa notificações relacionadas para que o sistema operacional possa recolhê-las na Central de Notificações/Central de Ações.

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "msg-42",
    Title:    "Jane Doe",
    Body:     "Are you free for lunch?",
    ThreadID: "chat:jane.doe",
})
```

### Nível de interrupção

`InterruptionLevel` controla a prioridade das notificações. Use uma das constantes exportadas:

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:                "alert-id",
    Title:             "Server down",
    Body:              "Investigate immediately",
    InterruptionLevel: notifications.InterruptionLevelTimeSensitive,
})
```

| Constante | Valor | Significado |
| --- | --- | --- |
| `InterruptionLevelPassive` | `"passive"` | Entrega silenciosa; não acende a tela nem reproduz o som padrão |
| `InterruptionLevelActive` | `"active"` | Nível padrão |
| `InterruptionLevelTimeSensitive` | `"timeSensitive"` | Ignora o Foco/Não Perturbe quando permitido |
| `InterruptionLevelCritical` | `"critical"` | Ignora o Foco e o modo silencioso; o macOS exige o direito de Alerta Crítico (sem ele, ocorre uma degradação silenciosa) |

Mapeamento por plataforma:

- **macOS** — define `UNNotificationContent.interruptionLevel`. `critical` exige o macOS 12+ e o direito de Alerta Crítico.
- **Windows** — é mapeado para o atributo `<toast scenario="...">` da notificação toast.
- **Linux** — é mapeado para a dica `urgency` do freedesktop.

### Entrega agendada

`Schedule` adia a entrega. Defina exatamente um entre `DelaySeconds` (segundos a partir de agora) e `At` (segundos Unix, UTC).

```go
// Deliver in 60 seconds
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "reminder-id",
    Title:    "Stand up and stretch",
    Schedule: &notifications.NotificationSchedule{DelaySeconds: 60},
})

// Deliver at an absolute time
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "alarm-id",
    Title:    "Meeting in 5 minutes",
    Schedule: &notifications.NotificationSchedule{At: time.Now().Add(time.Hour).Unix()},
})
```

@note{type="caution" title="Persistência"}
No **macOS**, as notificações agendadas usam um gatilho nativo e persistem após a reinicialização do aplicativo. No **Windows** e no **Linux**, o agendamento recorre a um temporizador `time.AfterFunc` no processo e é **perdido se o aplicativo for encerrado antes da entrega** — nem `wintoast` nem a especificação freedesktop expõem um mecanismo primitivo de entrega adiada.

@end

### Atualização de notificações

`UpdateNotification` substitui uma notificação em andamento que tenha o mesmo `ID`:

```go
notifier.UpdateNotification(notifications.NotificationOptions{
    ID:    "download-id",
    Title: "Download complete",
    Body:  "report.pdf is ready",
})
```

Comportamento por plataforma:

- **macOS** — `UNUserNotificationCenter` elimina duplicatas automaticamente pelo identificador, portanto a notificação existente é atualizada no próprio local.
- **Linux** — usa o parâmetro D-Bus `replaces_id` para substituir a notificação anterior.
- **Windows** — atualmente, entrega novamente como uma nova notificação. A substituição real no próprio local exige que `wintoast` ofereça suporte a `tag` / `group` upstream.

## Considerações sobre as plataformas

@tabs
[macOS]
No macOS, as notificações:

- Exigem autorização do usuário
- Exigem que o aplicativo seja empacotado e assinado (e notarizado para distribuição)
- Usam a aparência de notificação padrão do sistema
- Oferecem suporte a `Subtitle`
- Oferecem suporte à entrada de texto do usuário (respostas)
- Oferecem suporte à opção de ação `Destructive`
- Oferecem suporte a vários `Attachments` de qualquer tipo de mídia (imagens, áudio e vídeo)
- Oferecem suporte a `ThreadID` para agrupamento na Central de Notificações
- Oferecem suporte a todos os valores de `InterruptionLevel` (`critical` exige o direito de Alerta Crítico)
- Oferecem suporte à entrega agendada nativa que persiste após a reinicialização do aplicativo
- Eliminam automaticamente as chamadas duplicadas de `UpdateNotification` com base em `ID`
- Gerenciam automaticamente os modos escuro e claro

[Windows]
No Windows, as notificações:

- Usam os estilos de notificação toast do sistema Windows por meio do backend `wintoast`
- Adaptam-se às configurações de tema do Windows
- Oferecem suporte à entrada de texto do usuário (respostas)
- Oferecem suporte a telas com alta densidade de pixels (DPI)
- Não oferecem suporte a `Subtitle`
- Oferecem suporte a um único `Attachment` de imagem com a dica de posicionamento `hero`, `appLogoOverride` ou `inline` (o padrão é `inline`)
- Oferecem suporte a `ThreadID` para agrupamento na Central de Ações
- Oferecem suporte a `InterruptionLevel` por meio do atributo `scenario` da notificação toast
- Oferecem suporte à entrega agendada por meio de um temporizador no processo — **as notificações agendadas são perdidas se o aplicativo for encerrado antes da entrega**
- `UpdateNotification` atualmente entrega novamente como uma nova notificação (a substituição real no próprio local aguarda o suporte upstream de `wintoast` a `tag`/`group`)

[Linux]
No Linux, as notificações usam a interface D-Bus `org.freedesktop.Notifications`. Um daemon de notificações compatível **deve estar em execução** para que as notificações funcionem.

@note{type="caution" title="Requisito do sistema: daemon de notificações"}
Um daemon de notificações compatível com freedesktop deve estar instalado e em execução. Opções comuns:

- **dunst** — leve e altamente configurável (`apt install dunst` / `dnf install dunst`)
- **mako** — nativo do Wayland (`apt install mako-notifier`)
- **GNOME Shell** — registra automaticamente a interface no GNOME 43+. No Ubuntu 24.04 (GNOME Shell 46), talvez a interface não se registre automaticamente no início da sessão; instale `dunst` como alternativa se as notificações não aparecerem.
- **xfce4-notifyd** — incluído nos ambientes de desktop XFCE

Se nenhum daemon estiver em execução, `SendNotification` retornará um erro D-Bus: `The name org.freedesktop.Notifications was not provided by any .service files`. Trate esse erro no seu aplicativo e informe ao usuário que ele deve instalar um daemon de notificações.

@end

No Linux, as notificações:

- Seguem o tema do ambiente de desktop
- São posicionadas de acordo com as regras do ambiente de desktop
- Oferecem suporte a `Subtitle` (concatenado ao corpo para daemons que não o renderizam separadamente)
- Não oferecem suporte à entrada de texto do usuário (não faz parte da especificação freedesktop)
- Oferecem suporte a um único `Attachment` de imagem por meio da dica `image-path`
- Compatível com `ThreadID` (tratado pelo daemon quando houver suporte)
- `Sound.Name` é encaminhado como a dica `sound-name`; a reprodução depende do daemon ativo e do tema de sons
- Mapeiam `InterruptionLevel` para a dica `urgency` do freedesktop
- Compatível com entrega agendada por meio de um temporizador no processo — **as notificações agendadas serão perdidas se o aplicativo for encerrado antes da entrega**
- `UpdateNotification` usa o parâmetro `replaces_id` do D-Bus para substituir a notificação anterior no mesmo lugar

@end

## Práticas recomendadas

1. Verifique e solicite autorização:
  - O macOS exige autorização do usuário


2. Forneça notificações claras e concisas:
  - Use títulos, subtítulos, textos e títulos de ações descritivos


3. Trate adequadamente as respostas às notificações:
  - Verifique se há erros nas respostas às notificações
  - Forneça feedback para as ações do usuário


4. Considere as convenções da plataforma:
  - Siga os padrões de notificação específicos da plataforma
  - Respeite as configurações do sistema


5. No Linux, trate a dependência do daemon:
  - Verifique o erro retornado por `SendNotification` — a ausência de um daemon produz um erro do D-Bus
  - A documentação do pacote ou o README do aplicativo deve informar que é necessário um daemon de notificações compatível com freedesktop


## Exemplos

Explore este exemplo:

- [Notificações](https://github.com/wailsapp/wails/tree/master/v3/examples/notifications)

## Referência da API

### Gerenciamento do serviço

| Método | Descrição |
| --- | --- |
| `New()` | Cria um novo serviço de notificações |

### Autorização de notificações

| Método | Descrição |
| --- | --- |
| `RequestNotificationAuthorization()` | Solicita permissão para exibir notificações (macOS) |
| `CheckNotificationAuthorization()` | Verifica o status atual da autorização de notificações (macOS) |

### Envio de notificações

| Método | Descrição |
| --- | --- |
| `SendNotification(options NotificationOptions)` | Envia uma notificação básica |
| `SendNotificationWithActions(options NotificationOptions)` | Envia uma notificação interativa com ações |
| `UpdateNotification(options NotificationOptions)` | Atualiza uma notificação em andamento por `ID` (consulte [Atualização de notificações](#atualizao-de-notificaes)) |

### Categorias de notificação

| Método | Descrição |
| --- | --- |
| `RegisterNotificationCategory(category NotificationCategory)` | Registra uma categoria de notificação reutilizável |
| `RemoveNotificationCategory(categoryID string)` | Remove uma categoria registrada anteriormente |

### Gerenciamento de notificações

| Método | Descrição |
| --- | --- |
| `RemoveAllPendingNotifications()` | Remove todas as notificações pendentes (somente macOS e Linux) |
| `RemovePendingNotification(identifier string)` | Remove uma notificação pendente específica (somente macOS e Linux) |
| `RemoveAllDeliveredNotifications()` | Remove todas as notificações entregues (somente macOS e Linux) |
| `RemoveDeliveredNotification(identifier string)` | Remove uma notificação entregue específica (somente macOS e Linux) |
| `RemoveNotification(identifier string)` | Remove uma notificação (específico do Linux) |

### Tratamento de eventos

| Método | Descrição |
| --- | --- |
| `OnNotificationResponse(callback func(result NotificationResult))` | Registra uma função de callback para respostas às notificações |

### Structs e tipos

#### NotificationOptions

```go
type NotificationOptions struct {
    ID         string                 // Unique identifier for the notification (required)
    Title      string                 // Main notification title (required)
    Subtitle   string                 // Subtitle text (macOS and Linux only)
    Body       string                 // Main notification content
    CategoryID string                 // Category identifier for interactive notifications
    Data       map[string]interface{} // Custom data to associate with the notification

    // Sound controls the sound played on delivery. nil = platform default.
    Sound *NotificationSound

    // Attachments are media files shown alongside the notification.
    // macOS supports multiple attachments of any media type;
    // Windows and Linux honour the first image-typed attachment.
    Attachments []NotificationAttachment

    // ThreadID groups related notifications together in
    // Notification Center / Action Center / the Linux notification daemon.
    ThreadID string

    // InterruptionLevel controls priority. One of "passive",
    // "active" (default), "timeSensitive", "critical".
    InterruptionLevel string

    // Schedule defers delivery. macOS uses a native trigger that survives
    // restarts; Windows and Linux use an in-process timer that does NOT.
    Schedule *NotificationSchedule
}
```

#### NotificationSound

```go
type NotificationSound struct {
    Silent bool   // If true, no sound is played
    Name   string // Named/bundled sound (see "Custom Sound" above)
}
```

#### NotificationAttachment

```go
type NotificationAttachment struct {
    ID   string // Optional identifier
    Path string // Absolute filesystem path (macOS also accepts file:// URLs)
    // Type is an optional placement/UTI hint:
    //   macOS:   UTI like "public.png" / "public.audio" (often inferred)
    //   Windows: "hero" | "appLogoOverride" | "inline" (default "inline")
    //   Linux:   ignored (always image-path hint)
    Type string
}
```

#### NotificationSchedule

```go
// Exactly one of DelaySeconds or At must be set. At is Unix seconds (UTC).
type NotificationSchedule struct {
    DelaySeconds int   // Seconds from now until delivery
    At           int64 // Absolute Unix timestamp (seconds, UTC)
}
```

#### Constantes de InterruptionLevel

```go
const (
    InterruptionLevelPassive       = "passive"
    InterruptionLevelActive        = "active" // default
    InterruptionLevelTimeSensitive = "timeSensitive"
    InterruptionLevelCritical      = "critical"
)
```

#### NotificationCategory

```go
type NotificationCategory struct {
    ID               string                // Unique identifier for the category
    Actions          []NotificationAction  // Button actions for the notification
    HasReplyField    bool                  // Whether to include a text input field
    ReplyPlaceholder string                // Placeholder text for the input field
    ReplyButtonTitle string                // Text for the reply button
}
```

#### NotificationAction

```go
type NotificationAction struct {
    ID          string  // Unique identifier for the action
    Title       string  // Button text
    Destructive bool    // Whether the action is destructive (macOS-specific)
}
```

#### NotificationResponse

```go
type NotificationResponse struct {
    ID               string                  // Notification identifier
    ActionIdentifier string                  // Action that was triggered
    CategoryID       string                  // Category of the notification
    Title            string                  // Title of the notification
    Subtitle         string                  // Subtitle of the notification
    Body             string                  // Body text of the notification
    UserText         string                  // Text entered by the user
    UserInfo         map[string]interface{}  // Custom data from the notification
}
```

#### NotificationResult

```go
type NotificationResult struct {
    Response NotificationResponse  // Response data
    Error    error                 // Any error that occurred
}
```
