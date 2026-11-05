---
title: "Dock e barra de tarefas"
description: "Gerencie a visibilidade do ícone no Dock e exiba indicadores no macOS e no Windows"
slug: "features/platform/dock"
sourcePath: "features/platform/dock.md"
---

## Introdução

O Wails fornece um serviço Dock multiplataforma para aplicativos desktop. Esse serviço permite:

- Ocultar e mostrar o ícone do aplicativo no Dock do macOS
- Exibir indicadores no bloco do aplicativo ou no ícone do Dock/barra de tarefas (macOS e Windows)

## Uso básico

### Como criar o serviço

Primeiro, inicialize o serviço Dock:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"

// Create a new Dock service
dockService := dock.New()

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

### Como criar o serviço com opções de indicador personalizadas (somente Windows)

No Windows, você pode personalizar a aparência do indicador com várias opções:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"
import "image/color"

// Create a dock service with custom badge options
options := dock.BadgeOptions{
    TextColour:       color.RGBA{255, 255, 255, 255}, // White text
    BackgroundColour: color.RGBA{0, 0, 255, 255},     // Blue background
    FontName:         "consolab.ttf",                 // Bold Consolas font
    FontSize:         20,                             // Font size for single character
    SmallFontSize:    14,                             // Font size for multiple characters
}

dockService := dock.NewWithOptions(options)

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

## Operações do Dock

### Como ocultar o ícone do aplicativo no Dock

Oculte o ícone do aplicativo do Dock do macOS:

```go
// Hide the app icon
dockService.HideAppIcon()
```

### Como mostrar o ícone do aplicativo no Dock

Mostre o ícone do aplicativo no Dock do macOS:

```go
// Show the app icon
dockService.ShowAppIcon()
```

## Operações com indicadores

### Como definir um indicador

Defina um indicador no bloco/ícone do Dock do aplicativo:

```go
// Set a default badge
dockService.SetBadge("")

// Set a numeric badge
dockService.SetBadge("3")

// Set a text badge
dockService.SetBadge("New")
```

### Como definir um indicador personalizado (somente Windows)

Defina um indicador aplicando opções apenas a essa chamada:

```go
options := dock.BadgeOptions{
    BackgroundColour: color.RGBA{0, 255, 255, 255},
    FontName:         "arialb.ttf", // System font
    FontSize:         16,
    SmallFontSize:    10,
    TextColour:       color.RGBA{0, 0, 0, 255},
}

// Set a default badge
dockService.SetCustomBadge("", options)

// Set a numeric badge
dockService.SetCustomBadge("3", options)

// Set a text badge
dockService.SetCustomBadge("New", options)
```

### Como remover um indicador

Remova o indicador do ícone do aplicativo:

```go
dockService.RemoveBadge()
```

### Como obter o indicador definido

```go
dockService.GetBadge()
```

## Considerações sobre as plataformas

@tabs
[macOS]
No macOS:

- O ícone do Dock pode ser **ocultado** e **mostrado**
- Os indicadores são exibidos diretamente no ícone do Dock
- As opções do indicador **não são personalizáveis** (todas as opções passadas para `NewWithOptions`/`SetCustomBadge` são ignoradas)
- O estilo padrão dos indicadores do Dock do macOS é usado e se adapta automaticamente à aparência
- O sistema trata o estouro do rótulo
- Fornecer um rótulo vazio exibe o indicador padrão "●"

[Windows]
No Windows:

- Atualmente, este serviço não permite ocultar nem mostrar o ícone da barra de tarefas
- Os indicadores são exibidos como um ícone de sobreposição na barra de tarefas
- Os indicadores aceitam valores de texto
- A aparência do indicador pode ser personalizada por meio de `BadgeOptions`
- O aplicativo precisa ter uma janela para que os indicadores sejam exibidos
- Um tamanho de fonte menor é usado automaticamente em rótulos com vários caracteres
- O estouro do rótulo não é tratado
- Opções de personalização:
  - **TextColour**: cor do texto (padrão: branco)
  - **BackgroundColour**: cor de fundo do indicador (padrão: vermelho)
  - **FontName**: nome do arquivo da fonte (padrão: "segoeuib.ttf")
  - **FontSize**: tamanho da fonte para um único caractere (padrão: 18)
  - **SmallFontSize**: tamanho da fonte para vários caracteres (padrão: 14)


[Linux]
No Linux:

- Os recursos de visibilidade do ícone do Dock e de indicadores não estão disponíveis

@end

## Práticas recomendadas

1. **Ao ocultar o ícone do Dock (macOS):**
  - Garanta que os usuários ainda possam acessar o aplicativo (por exemplo, pela [área de notificação do sistema](/features/menus/systray/))
  - Inclua uma opção "Sair" na interface alternativa
  - O aplicativo não aparecerá no alternador Command+Tab
  - As janelas abertas permanecem visíveis e funcionais
  - Fechar todas as janelas pode não encerrar o aplicativo (o comportamento varia no macOS)
  - Os usuários perdem a forma padrão de encerrar o aplicativo clicando com o botão direito no ícone do Dock


2. **Use indicadores com moderação:**
  - Atualizações excessivas do indicador podem distrair os usuários
  - Reserve os indicadores para notificações importantes


3. **Mantenha o texto do indicador curto:**
  - Indicadores numéricos são mais eficazes
  - No macOS, os indicadores de texto devem ser breves


4. **Ao personalizar indicadores no Windows:**
  - Garanta um alto contraste entre as cores do texto e do fundo
  - Teste com diferentes comprimentos de texto, pois o tamanho da fonte diminui à medida que o texto aumenta
  - Use fontes comuns do sistema para garantir a disponibilidade


## Referência da API

### Gerenciamento do serviço

| Método | Descrição |
| --- | --- |
| `New()` | Cria um novo serviço de Dock |
| `NewWithOptions(options BadgeOptions)` | Cria um novo serviço de Dock com opções personalizadas de selo (somente no Windows; as opções são ignoradas no macOS e no Linux) |

### Operações do Dock

| Método | Descrição |
| --- | --- |
| `HideAppIcon()` | Oculta o ícone do aplicativo no Dock do macOS (somente no macOS) |
| `ShowAppIcon()` | Exibe o ícone do aplicativo no Dock do macOS (somente no macOS) |

### Operações de selo

| Método | Descrição |
| --- | --- |
| `SetBadge(label string) error` | Define um selo com o rótulo especificado |
| `SetCustomBadge(label string, options BadgeOptions) error` | Define um selo com o rótulo especificado e opções de estilo personalizadas (somente no Windows) |
| `RemoveBadge() error` | Remove o selo do ícone do aplicativo |
| `GetBadge() *string` | Obtém o selo atual |

### Estruturas e tipos

```go
// Options for customizing badge appearance (Windows only)
type BadgeOptions struct {
    TextColour       color.RGBA  // Color of the badge text
    BackgroundColour color.RGBA  // Color of the badge background
    FontName         string      // Font file name (e.g., "segoeuib.ttf")
    FontSize         int         // Font size for single character
    SmallFontSize    int         // Font size for multiple characters
}
```
