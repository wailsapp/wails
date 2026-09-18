---
title: "Personalização de janelas no Wails"
description: "Personalize a aparência e o comportamento das janelas em seus aplicativos Wails"
slug: "guides/customising-windows"
sourcePath: "guides/customising-windows.md"
---

Plataformas relevantes: <span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

O Wails fornece uma API para controlar a aparência e a funcionalidade dos controles de uma janela. Essa funcionalidade está disponível no Windows e no macOS, mas não no Linux.

## Definição dos estados dos botões da janela

Os estados dos botões são definidos pela enumeração `ButtonState`:

```go
type ButtonState int

const (
    ButtonEnabled   ButtonState = 0
    ButtonDisabled  ButtonState = 1
    ButtonHidden    ButtonState = 2
)
```

- `ButtonEnabled`: o botão está habilitado e visível.
- `ButtonDisabled`: o botão está visível, mas desabilitado (acinzentado).
- `ButtonHidden`: o botão fica oculto na barra de título.

Os estados dos botões podem ser definidos durante a criação da janela ou em tempo de execução.

### Definição dos estados dos botões durante a criação da janela

Ao criar uma nova janela, você pode definir o estado inicial dos botões usando a estrutura `WebviewWindowOptions`:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        MinimiseButtonState:   application.ButtonHidden,
        MaximiseButtonState:   application.ButtonDisabled,
        CloseButtonState:      application.ButtonEnabled,
        FullscreenButtonState: application.ButtonEnabled,
    })

    app.Run()
}
```

No exemplo acima, o botão de minimizar está oculto, o botão de maximizar está inativo (acinzentado) e o botão de fechar está ativo.

### Definição dos estados dos botões em tempo de execução

Você também pode alterar os estados dos botões em tempo de execução usando os seguintes métodos da interface `Window`:

```go
window.SetMinimiseButtonState(wails.ButtonHidden)
window.SetMaximiseButtonState(wails.ButtonEnabled)
window.SetCloseButtonState(wails.ButtonDisabled)
window.SetFullscreenButtonState(wails.ButtonEnabled)
```

### macOS: MaximiseButtonState e FullscreenButtonState compartilham um botão

No macOS, o botão verde do semáforo (`NSWindowZoomButton`) é o mesmo controle físico para maximização e tela cheia. Se `MaximiseButtonState` e `FullscreenButtonState` fossem definidos com valores diferentes durante a criação da janela, isso resultaria em uma substituição silenciosa na qual prevaleceria o último valor definido.

Para evitar isso, durante a inicialização, o Wails aplica o estado **mais restritivo** dos dois, na ordem `ButtonEnabled` < `ButtonDisabled` < `ButtonHidden`.

| `MaximiseButtonState` | `FullscreenButtonState` | Estado efetivo no macOS |
| --- | --- | --- |
| `ButtonEnabled` | `ButtonEnabled` | `ButtonEnabled` |
| `ButtonDisabled` | `ButtonEnabled` | `ButtonDisabled` |
| `ButtonEnabled` | `ButtonHidden` | `ButtonHidden` |
| `ButtonDisabled` | `ButtonHidden` | `ButtonHidden` |

Em tempo de execução, `SetMaximiseButtonState` e `SetFullscreenButtonState` controlam `NSWindowZoomButton` no macOS; portanto, prevalece a última chamada.

### Diferenças entre plataformas

A funcionalidade de estado dos botões se comporta de maneira ligeiramente diferente no Windows e no macOS:

|  | Windows | Mac |
| --- | --- | --- |
| Desabilitar Minimizar/Maximizar/Fechar | Desabilita Minimizar/Maximizar/Fechar | Desabilita Minimizar/Maximizar/Fechar |
| Ocultar Minimizar | Desabilita Minimizar | Oculta o botão Minimizar |
| Ocultar Maximizar | Desabilita Maximizar | Oculta o botão Maximizar |
| Ocultar Fechar | Oculta todos os controles | Oculta Fechar |
| `FullscreenButtonState` | Sem efeito | Controla o botão de zoom (verde) |

Observação: no Windows, não é possível ocultar os botões de minimizar e maximizar individualmente. No entanto, desabilitar ambos ocultará os dois controles e exibirá apenas o botão Fechar. O Windows não tem um botão dedicado para tela cheia na barra de título padrão, portanto `FullscreenButtonState` não tem efeito nessa plataforma.

### Como controlar o estilo da janela (Windows)

Para controlar o estilo da barra de título no Windows, use o campo `ExStyle` na struct `WebviewWindowOptions`:

Exemplo:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/w32"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Windows: application.WindowsWindow{
            ExStyle: w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOPMOST,
        },
    })

    app.Run()
}
```

Esta configuração substituirá outras opções que afetam o estilo estendido de uma janela:

- HiddenOnTaskbar
- AlwaysOnTop
- IgnoreMouseEvents
- BackgroundType
