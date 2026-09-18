---
title: "Janelas sem moldura"
description: "Crie elementos de janela personalizados com janelas sem moldura"
slug: "features/windows/frameless"
sourcePath: "features/windows/frameless.md"
---

## Janelas sem moldura

O Wails oferece **suporte a janelas sem moldura**, com regiões de arraste baseadas em CSS e comportamento nativo da plataforma. Remova a barra de título nativa da plataforma para ter controle total sobre os elementos da janela, criar designs personalizados e proporcionar experiências de usuário únicas, mantendo funcionalidades essenciais como arrastar, redimensionar e usar os controles do sistema.

![O aplicativo inicial TypeScript padrão do Wails v3 executado como uma janela sem moldura, com cantos nativos do macOS](/assets/screenshots/frameless-v3-native-corners-macos.png)

O exemplo acima é o aplicativo inicial TypeScript padrão do Wails v3 com `Frameless: true` habilitado.

## Início rápido

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Frameless App",
    Width:     800,
    Height:    600,
    Frameless: true,
})
```

**CSS para uma barra de título arrastável:**

```css
.titlebar {
    --wails-draggable: drag;
    height: 40px;
    background: #333;
}

.titlebar button {
    --wails-draggable: no-drag;
}
```

**HTML:**

```html
<div class="titlebar">
    <span>My Application</span>
    <button onclick="window.close()">×</button>
</div>
```

**Pronto!** Agora você tem uma barra de título personalizada.

## Como criar janelas sem moldura

### Raio dos cantos (macOS)

Por padrão, as janelas sem moldura preservam os cantos arredondados padrão do AppKit no macOS. Defina `Mac.CornerRadius` para usar um raio personalizado (em pontos):

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerRadius: 16,
    },
})
```

Defina `Mac.CornerType` como `MacWindowCornerTypeSquare` para obter cantos retos. Essa configuração ignora `CornerRadius`:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerType: application.MacWindowCornerTypeSquare,
    },
})
```

### Janela sem moldura básica

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Width:     800,
    Height:    600,
})
```

**O que você obtém:**

- Sem barra de título
- Sem bordas de janela
- Sem botões do sistema
- Plano de fundo transparente (opcional)

**O que você precisa implementar:**

- Área arrastável
- Botões para fechar, minimizar e maximizar
- Alças de redimensionamento (se a janela for redimensionável)

### Com plano de fundo transparente

**API privada no macOS:** defina `Mac.Backdrop: application.MacBackdropTransparent` e compile com `-tags private_mac_apis` para permitir a transparência da webview. Sem a tag, a webview nativa permanece opaca, mesmo que o plano de fundo em HTML/CSS seja transparente. `Frameless` e `TitleBar.AppearsTransparent` usam APIs públicas. Consulte [APIs privadas do macOS](/guides/build/private-macos-apis/#webview-transparency-and-background).

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

**Casos de uso:**

- Cantos arredondados
- Formas personalizadas
- Janelas de sobreposição
- Telas de abertura

## Regiões de arraste

### Arraste baseado em CSS

Use a propriedade CSS `--wails-draggable`:

```css
/* Draggable area */
.titlebar {
    --wails-draggable: drag;
}

/* Non-draggable elements within draggable area */
.titlebar button {
    --wails-draggable: no-drag;
}
```

**Valores:**

- `drag` — A área pode ser arrastada
- `no-drag` — A área não pode ser arrastada (mesmo que a área pai possa)

### Exemplo completo de barra de título

```html
<div class="titlebar">
    <div class="title">My Application</div>
    <div class="controls">
        <button class="minimize">−</button>
        <button class="maximize">□</button>
        <button class="close">×</button>
    </div>
</div>
```

```css
.titlebar {
    --wails-draggable: drag;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    background: #2c2c2c;
    color: white;
    padding: 0 16px;
}

.title {
    font-size: 14px;
    user-select: none;
}

.controls {
    display: flex;
    gap: 8px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: white;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
}

.controls button:hover {
    background: rgba(255, 255, 255, 0.1);
}

.controls .close:hover {
    background: #e81123;
}
```

**JavaScript para os botões:**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

## Regiões não clientes nativas no Windows

O Windows pode tratar partes de uma barra de título personalizada como áreas não clientes nativas. Assim, você pode desenhar a barra de título e seus botões com qualquer design em HTML/CSS, preservando o comportamento nativo do Windows: a área da barra de título arrasta a janela; o botão de maximização pode exibir os recursos Snap Assist/Snap Layouts do Windows 11; e os botões de minimizar, maximizar e fechar recebem a detecção nativa de regiões sob o ponteiro e o estado nativo do mouse.

O vídeo abaixo mostra uma barra de título personalizada em HTML/CSS que usa a detecção nativa de regiões sob o ponteiro do Windows, incluindo os recursos Snap Assist/Snap Layouts do Windows 11 em um botão de maximização personalizado.

<video src="/assets/windows-native-non-client-regions/wails-app-region.mp4" controls muted playsInline></video>

O Wails oferece suporte a dois mecanismos específicos do Windows:

- `app-region` por meio do suporte nativo do WebView2 a regiões não clientes
- `--wails-non-client-region` por meio do rastreamento do runtime do Wails para botões de legenda personalizados

### Como escolher um modo

@note{type="caution" title="Experimental"}
`WebView2CompositionHosting` altera internamente a forma como a janela hospeda o WebView2 e interage com ele. Em vez do controlador WebView2 padrão hospedado em HWND, o Wails usa a hospedagem por controlador de composição e encaminha explicitamente as entradas. Esse modo pode apresentar problemas de renderização, entrada, foco ou compatibilidade com o WebView2 Runtime. Habilite-o somente quando precisar do comportamento nativo de botões de legenda personalizados e teste cuidadosamente seu aplicativo nas versões do Windows e do WebView2 Runtime às quais você oferece suporte.

@end

Escolha de acordo com o que você precisa do Windows:

- Use `NonClientRegionSupport` para permitir o arraste nativo simples do aplicativo com os recursos `app-region: drag` e `app-region: no-drag` do WebView2.
- Use `WebView2CompositionHosting` quando quiser que seus botões personalizados de minimizar, maximizar e fechar se comportem como botões de legenda nativos do Windows.
- Habilite ambos quando a mesma janela precisar do suporte nativo do WebView2 a `app-region` e de regiões de botões de legenda personalizados gerenciadas pelo Wails.

`NonClientRegionSupport` é a alternativa nativa e leve ao rastreamento de `--wails-draggable` do Wails. Você marca com CSS as áreas arrastáveis e não arrastáveis, o WebView2 determina quais pixels pertencem à barra de título e, durante a detecção de regiões sob o ponteiro, o Wails solicita ao WebView2 a região nativa.

Atualmente, esse é todo o escopo deste modo. Ele não faz com que os botões personalizados de minimizar, maximizar ou fechar se comportem como botões de legenda nativos do Windows, nem habilita os recursos Snap Assist/Snap Layouts do Windows 11 para um botão de maximização personalizado. Use-o quando precisar do arraste nativo simples do aplicativo sem os mecanismos adicionais de `--wails-draggable`.

`WebView2CompositionHosting` destina-se a botões personalizados da barra de título com comportamento nativo. O Wails rastreia retângulos do DOM marcados com `--wails-non-client-region`, mapeia-os para valores de detecção de regiões sob o ponteiro do Windows, como `HTMINBUTTON`, `HTMAXBUTTON` e `HTCLOSE`, e encaminha a entrada do mouse de volta para a superfície do WebView2 hospedada por composição. É isso que permite que um botão personalizado de maximização participe dos recursos Snap Assist/Snap Layouts do Windows 11, mantendo qualquer design visual que você escolher.

Em outras palavras: `NonClientRegionSupport` é o suporte nativo do WebView2 a regiões CSS. `WebView2CompositionHosting` significa que o Wails assume a responsabilidade pela composição controlada pelo host e pela detecção personalizada de regiões não clientes sob o ponteiro.

### app-region do WebView2

Ative o suporte nativo do WebView2 a regiões não cliente para a janela:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport: true,
    },
})
```

Em seguida, marque as áreas arrastáveis com a propriedade CSS `app-region`:

```css
.titlebar {
    app-region: drag;
}

.titlebar button,
.titlebar input,
.titlebar select,
.titlebar textarea {
    app-region: no-drag;
}
```

Use esse modo quando você precisar apenas do arraste nativo pela barra de título e os controles dessa barra forem tratados por cliques normais no frontend.

A limitação desse modo é que ele está restrito ao suporte do próprio WebView2 a regiões não cliente. Nas versões atuais do WebView2, isso significa apenas regiões arrastáveis e não arrastáveis. Ele não se destina a representar botões de legenda totalmente personalizados no frontend, com funções nativas distintas de minimizar, maximizar e fechar.

### Botões de legenda personalizados com comportamento nativo

Para que botões personalizados de minimizar, maximizar e fechar se comportem como botões de legenda do sistema, ative a hospedagem por composição:

@note{type="caution" title="Experimental"}
`WebView2CompositionHosting` usa a hospedagem do controlador de composição do WebView2 com DirectComposition. Consulte [Escolha de um modo](#como-escolher-um-modo) antes de ativá-la.

@end

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        WebView2CompositionHosting: true,
    },
})
```

Em seguida, marque cada região do frontend com `--wails-non-client-region`:

```html
<div class="titlebar">
    <div class="title">My Application</div>
    <div class="window-controls">
        <button class="window-button minimize" aria-label="Minimize"></button>
        <button class="window-button maximize" aria-label="Maximize"></button>
        <button class="window-button close" aria-label="Close"></button>
    </div>
</div>
```

```css
.titlebar {
    --wails-non-client-region: caption;
    height: 40px;
}

.window-controls {
    display: flex;
    height: 100%;
}

.window-button {
    width: 46px;
    border: 0;
    background: transparent;
}

.window-button.minimize {
    --wails-non-client-region: minimize;
}

.window-button.maximize {
    --wails-non-client-region: maximize;
}

.window-button.close {
    --wails-non-client-region: close;
}
```

Valores de `--wails-non-client-region` compatíveis:

- `caption` — área arrastável da barra de título
- `minimize` — alvo de clique do botão nativo de minimizar
- `maximize` — alvo de clique do botão nativo de maximizar, incluindo o comportamento do Snap Assist/Layouts de Ajuste do Windows 11 ao passar o ponteiro
- `close` — alvo de clique do botão nativo de fechar

O runtime do Wails observa alterações no DOM, nos estilos, no tamanho, na rolagem e na viewport e, em seguida, envia instantâneos das regiões para a janela nativa. A geometria das regiões é medida em pixels CSS e convertida em pixels físicos para a detecção de regiões sob o ponteiro do Windows.

O design visual permanece inteiramente sob seu controle. As regiões apenas informam ao Windows o significado de cada retângulo; a forma, o ícone, a cor, o espaçamento, o estilo ao passar o ponteiro e o layout dos botões continuam sendo definidos pelo seu frontend.

### Combinação dos dois modos

Você pode ativar as duas opções quando quiser usar o suporte do WebView2 a `app-region` e regiões de botões da barra de título gerenciadas pelo Wails na mesma janela:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport:    true,
        WebView2CompositionHosting: true,
    },
})
```

## Botões do sistema

### Implementação de fechar, minimizar e maximizar

**No lado do Go:**

```go
type WindowControls struct {
    window *application.WebviewWindow
}

func (wc *WindowControls) Minimise() {
    wc.window.Minimise()
}

func (wc *WindowControls) Maximise() {
    if wc.window.IsMaximised() {
        wc.window.UnMaximise()
    } else {
        wc.window.Maximise()
    }
}

func (wc *WindowControls) Close() {
    wc.window.Close()
}
```

**No lado do JavaScript:**

```javascript
import { Minimise, Maximise, Close } from './bindings/WindowControls'

document.querySelector('.minimize').addEventListener('click', Minimise)
document.querySelector('.maximize').addEventListener('click', Maximise)
document.querySelector('.close').addEventListener('click', Close)
```

**Ou use os métodos do runtime:**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

### Alternância do estado de maximização

Acompanhe o estado de maximização para atualizar o ícone do botão:

```javascript
import { Window } from '@wailsio/runtime'

async function toggleMaximise() {
    const isMaximised = await Window.IsMaximised()

    if (isMaximised) {
        await Window.Restore()
    } else {
        await Window.Maximise()
    }

    updateMaximiseButton()
}

async function updateMaximiseButton() {
    const isMaximised = await Window.IsMaximised()
    const button = document.querySelector('.maximize')
    button.textContent = isMaximised ? '❐' : '□'
}
```

## Alças de redimensionamento

### Redimensionamento baseado em CSS

O Wails fornece alças de redimensionamento automáticas para janelas sem moldura:

```css
/* Enable resize on all edges */
body {
    --wails-resize: all;
}

/* Or specific edges */
.resize-top {
    --wails-resize: top;
}

.resize-bottom {
    --wails-resize: bottom;
}

.resize-left {
    --wails-resize: left;
}

.resize-right {
    --wails-resize: right;
}

/* Corners */
.resize-top-left {
    --wails-resize: top-left;
}

.resize-top-right {
    --wails-resize: top-right;
}

.resize-bottom-left {
    --wails-resize: bottom-left;
}

.resize-bottom-right {
    --wails-resize: bottom-right;
}
```

**Valores:**

- `all` — redimensionar por todas as bordas
- `top`, `bottom`, `left`, `right` — bordas específicas
- `top-left`, `top-right`, `bottom-left`, `bottom-right` — cantos
- `none` — sem redimensionamento

### Exemplo de alça de redimensionamento

```html
<div class="window">
    <div class="titlebar">...</div>
    <div class="content">...</div>
    <div class="resize-handle resize-bottom-right"></div>
</div>
```

```css
.resize-handle {
    position: absolute;
    width: 16px;
    height: 16px;
}

.resize-bottom-right {
    --wails-resize: bottom-right;
    bottom: 0;
    right: 0;
    cursor: nwse-resize;
}
```

## Comportamento específico da plataforma

@tabs{sync-key="platform"}
[Windows]
**Janelas sem moldura no Windows:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        DisableFramelessWindowDecorations: false,
    },
})
```

**Recursos:**

- Sombra projetada automática
- Suporte a Layouts de Ajuste (Windows 11)
- Suporte ao Aero Snap
- Dimensionamento de DPI

**Desativar decorações:**

```go
Windows: application.WindowsWindow{
    DisableFramelessWindowDecorations: true,
},
```

**Snap Assist:**

```go
// Trigger Windows 11 Snap Assist
window.SnapAssist()
```

Isso aciona os Layouts de Ajuste por meio do atalho de teclado do Windows. Para usar um botão HTML personalizado de maximizar com Layouts de Ajuste nativos ao passar o ponteiro, use [Regiões não cliente nativas no Windows](#regies-no-clientes-nativas-no-windows).

**Altura personalizada da barra de título:** O Windows detecta automaticamente as regiões de arraste pelo CSS.

[macOS]
**Janelas sem moldura no macOS:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        TitleBar: application.MacTitleBar{
            AppearsTransparent: true,
        },
        InvisibleTitleBarHeight: 40,
    },
})
```

**Recursos:**

- Suporte nativo a tela cheia
- Botões de semáforo (opcionais)
- Efeitos de vibrância
- Barra de título transparente

**Oculte completamente a barra de título** (use as variantes predefinidas exportadas pelo pacote `application` — não há nenhum campo `TitleBarStyle` nem uma constante `MacTitleBarStyleHidden`):

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBarHidden,
},
```

Outras predefinições incluem `MacTitleBarDefault`, `MacTitleBarHiddenInset` e `MacTitleBarHiddenInsetUnified`.

**Barra de título invisível:** Permite arrastar a janela enquanto oculta a barra de título. Isso só tem efeito quando a janela não tem moldura ou usa `AppearsTransparent`:

```go
Mac: application.MacWindow{
    InvisibleTitleBarHeight: 40,
},
```

[Linux]
**Janelas sem moldura no Linux:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
})
```

**Recursos:**

- Suporte básico a janelas sem moldura
- Regiões de arraste em CSS
- Varia de acordo com o ambiente de desktop

**Observações sobre ambientes de desktop:**

- **GNOME:** Bom suporte
- **KDE Plasma:** Bom suporte
- **XFCE:** Suporte básico
- **Gerenciadores de janelas lado a lado:** Suporte limitado

**Compositor necessário:** A transparência requer um compositor (a maioria dos ambientes de desktop modernos tem um).

@end

## Padrões comuns

### Padrão 1: barra de título moderna

```html
<div class="modern-titlebar">
    <div class="app-icon">
        <img src="/icon.png" alt="App Icon">
    </div>
    <div class="title">My Application</div>
    <div class="controls">
        <button class="minimize">−</button>
        <button class="maximize">□</button>
        <button class="close">×</button>
    </div>
</div>
```

```css
.modern-titlebar {
    --wails-draggable: drag;
    display: flex;
    align-items: center;
    height: 40px;
    background: linear-gradient(to bottom, #3a3a3a, #2c2c2c);
    border-bottom: 1px solid #1a1a1a;
    padding: 0 16px;
}

.app-icon {
    --wails-draggable: no-drag;
    width: 24px;
    height: 24px;
    margin-right: 12px;
}

.title {
    flex: 1;
    font-size: 13px;
    color: #e0e0e0;
    user-select: none;
}

.controls {
    display: flex;
    gap: 1px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 46px;
    height: 32px;
    border: none;
    background: transparent;
    color: #e0e0e0;
    font-size: 14px;
    cursor: pointer;
    transition: background 0.2s;
}

.controls button:hover {
    background: rgba(255, 255, 255, 0.1);
}

.controls .close:hover {
    background: #e81123;
    color: white;
}
```

### Padrão 2: tela de abertura

```go
splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "Loading...",
    Width:          400,
    Height:         300,
    Frameless:      true,
    AlwaysOnTop:    true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
    display: flex;
    justify-content: center;
    align-items: center;
}

.splash {
    background: white;
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    padding: 40px;
    text-align: center;
}
```

### Padrão 3: janela arredondada

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
    margin: 8px;
}

.window {
    background: white;
    border-radius: 16px;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.15);
    overflow: hidden;
    height: calc(100vh - 16px);
}

.titlebar {
    --wails-draggable: drag;
    background: #f5f5f5;
    border-bottom: 1px solid #e0e0e0;
}
```

### Padrão 4: janela de sobreposição

```go
overlay := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    AlwaysOnTop:    true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
}

.overlay {
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(10px);
    border-radius: 8px;
    padding: 20px;
}
```

## Exemplo completo

Veja uma janela sem moldura pronta para produção:

**Go:**

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "Frameless App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:     "Frameless Application",
        Width:     1000,
        Height:    700,
        MinWidth:  800,
        MinHeight: 600,
        Frameless: true,

        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            InvisibleTitleBarHeight: 40,
        },

        Windows: application.WindowsWindow{
            DisableFramelessWindowDecorations: false,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

**HTML:**

```html
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="/style.css">
</head>
<body>
    <div class="window">
        <div class="titlebar">
            <div class="title">Frameless Application</div>
            <div class="controls">
                <button class="minimize" title="Minimise">−</button>
                <button class="maximize" title="Maximise">□</button>
                <button class="close" title="Close">×</button>
            </div>
        </div>
        <div class="content">
            <h1>Hello from Frameless Window!</h1>
        </div>
    </div>
    <script src="/main.js" type="module"></script>
</body>
</html>
```

**CSS:**

```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #f5f5f5;
}

.window {
    height: 100vh;
    display: flex;
    flex-direction: column;
}

.titlebar {
    --wails-draggable: drag;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    background: #ffffff;
    border-bottom: 1px solid #e0e0e0;
    padding: 0 16px;
}

.title {
    font-size: 13px;
    font-weight: 500;
    color: #333;
    user-select: none;
}

.controls {
    display: flex;
    gap: 8px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: #666;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
}

.controls button:hover {
    background: #f0f0f0;
    color: #333;
}

.controls .close:hover {
    background: #e81123;
    color: white;
}

.content {
    flex: 1;
    padding: 40px;
    overflow: auto;
}
```

**JavaScript:**

```javascript
import { Window } from '@wailsio/runtime'

// Minimise button
document.querySelector('.minimize').addEventListener('click', () => {
    Window.Minimise()
})

// Maximise/restore button
const maximiseBtn = document.querySelector('.maximize')
maximiseBtn.addEventListener('click', async () => {
    const isMaximised = await Window.IsMaximised()

    if (isMaximised) {
        await Window.Restore()
    } else {
        await Window.Maximise()
    }

    updateMaximiseButton()
})

// Close button
document.querySelector('.close').addEventListener('click', () => {
    Window.Close()
})

// Update maximise button icon
async function updateMaximiseButton() {
    const isMaximised = await Window.IsMaximised()
    maximiseBtn.textContent = isMaximised ? '❐' : '□'
    maximiseBtn.title = isMaximised ? 'Restore' : 'Maximise'
}

// Initial state
updateMaximiseButton()
```

## Práticas recomendadas

### ✅ Faça

- **Forneça uma área arrastável** — os usuários precisam mover a janela
- **Implemente os botões do sistema** — fechar, minimizar e maximizar
- **Defina um tamanho mínimo** — evite layouts inutilizáveis
- **Teste em todas as plataformas** — o comportamento varia
- **Use CSS para as regiões de arraste** — é flexível e fácil de manter
- **Forneça feedback visual** — estados de foco do ponteiro nos botões

### ❌ Não faça

- **Não se esqueça das alças de redimensionamento** — se a janela for redimensionável
- **Não torne a janela inteira arrastável** — isso impede a interação
- **Não se esqueça de desativar o arraste nos botões** — caso contrário, eles não funcionarão
- **Não use áreas de arraste minúsculas** — são difíceis de agarrar
- **Não se esqueça das diferenças entre plataformas** — teste minuciosamente

## Solução de problemas

### Não é possível arrastar a janela

**Causa:** `--wails-draggable: drag` ausente

**Solução:**

```css
.titlebar {
    --wails-draggable: drag;
}
```

### Os botões não funcionam

**Causa:** os botões estão na área arrastável

**Solução:**

```css
.titlebar button {
    --wails-draggable: no-drag;
}
```

### Não é possível redimensionar a janela

**Causa:** alças de redimensionamento ausentes

**Solução:**

```css
body {
    --wails-resize: all;
}
```

## Próximas etapas

@cards{cols="2"}
▣ Conceitos básicos de janelas
Conheça os fundamentos do gerenciamento de janelas.

[Saiba mais →](/features/windows/basics/)

---
⚙ Opções de janela
Referência completa das opções de janela.

[Saiba mais →](/features/windows/options/)

---
🚀 Eventos de janela
Trate os eventos do ciclo de vida da janela.

[Saiba mais →](/features/windows/events/)

---
◆ Várias janelas
Padrões para aplicativos com várias janelas.

[Saiba mais →](/features/windows/multiple/)

@end

---

**Tem alguma dúvida?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte o [exemplo de janela sem moldura](https://github.com/wailsapp/wails/tree/master/v3/examples/frameless).
