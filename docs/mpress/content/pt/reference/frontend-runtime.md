---
title: "Runtime do frontend"
description: "O pacote de runtime JavaScript do Wails para integração com o frontend"
slug: "reference/frontend-runtime"
sourcePath: "reference/frontend-runtime.md"
---

O runtime do frontend do Wails é a biblioteca padrão para aplicações Wails. Ele oferece diversos recursos que podem ser usados em suas aplicações, incluindo:

- Gerenciamento de janelas
- Caixas de diálogo
- Integração com o navegador
- Área de transferência
- Menus
- Informações do sistema
- Eventos
- Menus de contexto
- Telas
- WML (Linguagem de Marcação do Wails)

O runtime é necessário para a integração entre o Go e o frontend. Há 2 maneiras de integrar o runtime:

- Usando o pacote `@wailsio/runtime`
- Usando um bundle pré-compilado

## Usando o pacote npm

O pacote `@wailsio/runtime` é um pacote JavaScript que fornece acesso ao runtime do Wails pelo frontend. Ele é usado por todos os modelos padrão e é a maneira recomendada de integrar o runtime à sua aplicação. Ao usar o pacote `@wailsio/runtime`, você incluirá apenas as partes do runtime que utiliza.

O pacote está disponível no npm e pode ser instalado usando:

```shell
npm install --save @wailsio/runtime
```

## Usando um bundle pré-compilado

Alguns projetos não usam um empacotador JavaScript e podem preferir uma versão pré-compilada e empacotada do runtime. Essa versão pode ser gerada localmente usando o seguinte comando:

```shell
wails3 generate runtime
```

O comando gerará um arquivo `runtime.js` (e `runtime-debug.js`) no diretório atual. Esse arquivo é um módulo ES que pode ser importado pelos scripts da sua aplicação da mesma forma que o pacote npm, mas a API também é exportada para o objeto global da janela. Assim, em aplicações mais simples, você pode usá-lo da seguinte maneira:

```html
<html>
    <head>
        <script type="module" src="./runtime.js"></script>
        <script>
            window.onload = function () {
                wails.Window.SetTitle("A new window title");
            }
        </script>
    </head>
    <!--- ... -->
</html>
```

@note{type="caution"}
É importante incluir o atributo `type="module"` na tag `<script>` que carrega o runtime e aguardar o carregamento completo da página antes de chamar a API, pois scripts com o atributo `type="module"` são executados de forma assíncrona.

@end

## Inicialização

Além das funções da API, o runtime oferece suporte a menus de contexto e ao arraste de janelas. Esses recursos funcionarão como esperado somente depois que o runtime for inicializado. Mesmo que você não use a API, inclua uma instrução de importação por efeito colateral em algum lugar do código do frontend:

```javascript
import "@wailsio/runtime";
```

Seu empacotador deve detectar a presença de efeitos colaterais e incluir no build todo o código de inicialização necessário.

@note{type="info"}
Se preferir o bundle pré-compilado, basta adicionar uma tag de script conforme mostrado acima.

@end

## Plugin do Vite para eventos tipados

O runtime inclui um plugin do Vite que habilita o suporte a HMR (Hot Module Replacement) para eventos tipados durante o desenvolvimento.

### Configuração

Adicione o plugin ao seu `vite.config.ts`:

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails("./bindings")],
})
```

O argumento é o diretório dos bindings gerados, relativo ao projeto Vite. Gere os bindings antes de compilar o frontend.

### Benefícios

- **Recarregamento automático**: os bindings de eventos são gerados novamente e recarregados automaticamente quando você executa `wails3 generate bindings`
- **Modo de desenvolvimento**: funciona perfeitamente com `wails3 dev` para oferecer atualizações instantâneas
- **Segurança de tipos**: suporte completo ao TypeScript, com preenchimento automático e verificação de tipos

### Uso com registro de eventos

Registre seus eventos no Go:

```go
type UserData struct {
    ID   string
    Name string
}

func init() {
    application.RegisterEvent[UserData]("user-updated")
}
```

Gere os bindings:

```bash
wails3 generate bindings

# Or, to include TypeScript definitions
wails3 generate bindings -ts

# For more options, see:
wails3 generate bindings -help
```

Use eventos tipados no frontend:

```typescript
import { Events } from '@wailsio/runtime'
import { UserUpdated } from './bindings/events'

// Type-safe event with autocomplete
Events.Emit(UserUpdated, {
    ID: "123",
    Name: "John Doe"
})
```

## Referência da API

O runtime é organizado em módulos, cada um oferecendo funcionalidades específicas. Importe apenas o que for necessário:

```javascript
import { Events, Window, Clipboard } from '@wailsio/runtime'
```

### Eventos

Sistema de eventos para comunicação entre Go e JavaScript. As exportações de eventos geradas, como `UserUpdated`, são constantes de string, não fábricas de payload. O gerador amplia `CustomEvents` para associar os nomes dos eventos registrados aos tipos de payload. Os tipos abaixo (`WailsEventName`, `WailsEventCallback`, `WailsEventData`) são exportados pelo módulo Events.

#### On()

Registra uma função de callback para um evento.

```typescript
function On<E extends WailsEventName = WailsEventName>(eventName: E, callback: WailsEventCallback<E>): () => void
```

**Retorno:** função para cancelar a inscrição

**Exemplo:**

```javascript
import { Events } from '@wailsio/runtime'

// Basic event listening
const unsubscribe = Events.On('user-logged-in', (event) => {
    console.log('User:', event.data.username)
})

// With typed events (TypeScript)
import { UserLogin } from './bindings/events'

Events.On(UserLogin, (event) => {
    // event.data is typed as UserLoginData
    console.log('User:', event.data.username)
})

// Later: unsubscribe()
```

#### Once()

Registra uma função de callback que é executada apenas uma vez.

```typescript
function Once<E extends WailsEventName = WailsEventName>(eventName: E, callback: WailsEventCallback<E>): () => void
```

**Exemplo:**

```javascript
import { Events } from '@wailsio/runtime'

Events.Once('app-ready', () => {
    console.log('App initialized')
})
```

#### Emit()

Emite um evento para o backend em Go ou para outras janelas.

```typescript
function Emit<E extends WailsEventName = WailsEventName>(name: E, data: WailsEventData<E>): Promise<boolean>
function Emit<E extends WailsEventName = WailsEventName>(name: WailsEventData<E> extends null | void ? E : never): Promise<boolean>
```

**Retorno:** uma Promise que é resolvida como `true` se o evento tiver sido cancelado; caso contrário, como `false`

**Exemplo:**

```javascript
import { Events } from '@wailsio/runtime'

// Basic event emission
const wasCancelled = await Events.Emit('button-clicked', { buttonId: 'submit' })

// With typed events (TypeScript)
import { UserLogin } from './bindings/events'

const cancelled = await Events.Emit(UserLogin, {
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
})

if (cancelled) {
    console.log('Login was cancelled by a hook')
}
```

@note{type="info"}
O valor retornado indica se o evento foi cancelado por um hook. A maioria dos eventos não pode ser cancelada e sempre retornará `false`.

@end

#### Off()

Remove listeners de eventos.

```typescript
function Off(...eventNames: string[]): void
```

**Exemplo:**

```javascript
import { Events } from '@wailsio/runtime'

Events.Off('user-logged-in', 'user-logged-out')
```

#### OffAll()

Remove todos os listeners de eventos.

```typescript
function OffAll(): void
```

### Janela

Métodos de gerenciamento de janelas. A exportação padrão é a janela atual.

```javascript
import { Window } from '@wailsio/runtime'

// Current window
await Window.SetTitle('New Title')
await Window.Center()

// Get another window
const otherWindow = Window.Get('secondary')
await otherWindow.Show()
```

#### Visibilidade

**Show()** - Exibe a janela

```typescript
function Show(): Promise<void>
```

**Hide()** - Oculta a janela

```typescript
function Hide(): Promise<void>
```

**Close()** - Fecha a janela

```typescript
function Close(): Promise<void>
```

#### Tamanho e posição

**SetSize(width, height)** - Define o tamanho da janela

```typescript
function SetSize(width: number, height: number): Promise<void>
```

**Size()** - Obtém o tamanho da janela

```typescript
function Size(): Promise<{ width: number, height: number }>
```

**SetPosition(x, y)** - Define a posição absoluta

```typescript
function SetPosition(x: number, y: number): Promise<void>
```

**Position()** - Obtém a posição absoluta

```typescript
function Position(): Promise<{ x: number, y: number }>
```

**Center()** - Centraliza a janela

```typescript
function Center(): Promise<void>
```

**Exemplo:**

```javascript
import { Window } from '@wailsio/runtime'

// Resize and center
await Window.SetSize(800, 600)
await Window.Center()

// Get current size
const { width, height } = await Window.Size()
```

#### Estado da janela

**Minimise()** - Minimiza a janela

```typescript
function Minimise(): Promise<void>
```

**Maximise()** - Maximiza a janela

```typescript
function Maximise(): Promise<void>
```

**Fullscreen()** - Entra no modo de tela cheia

```typescript
function Fullscreen(): Promise<void>
```

**Restore()** - Restaura a janela do estado minimizado, maximizado ou de tela cheia

```typescript
function Restore(): Promise<void>
```

**IsMinimised()** - Verifica se a janela está minimizada

```typescript
function IsMinimised(): Promise<boolean>
```

**IsMaximised()** - Verifica se a janela está maximizada

```typescript
function IsMaximised(): Promise<boolean>
```

**IsFullscreen()** - Verifica se a janela está em tela cheia

```typescript
function IsFullscreen(): Promise<boolean>
```

#### Propriedades da janela

**SetTitle(title)** - Define o título da janela

```typescript
function SetTitle(title: string): Promise<void>
```

**Name()** - Obtém o nome da janela

```typescript
function Name(): Promise<string>
```

**SetBackgroundColour(r, g, b, a)** - Define a cor de fundo

```typescript
function SetBackgroundColour(r: number, g: number, b: number, a: number): Promise<void>
```

**SetAlwaysOnTop(alwaysOnTop)** - Mantém a janela sempre visível sobre as demais

```typescript
function SetAlwaysOnTop(alwaysOnTop: boolean): Promise<void>
```

**SetResizable(resizable)** - Permite redimensionar a janela

```typescript
function SetResizable(resizable: boolean): Promise<void>
```

#### Foco e tela

**Focus()** - Coloca o foco na janela

```typescript
function Focus(): Promise<void>
```

**IsFocused()** - Verifica se a janela está em foco

```typescript
function IsFocused(): Promise<boolean>
```

**GetScreen()** - Obtém a tela em que a janela está

```typescript
function GetScreen(): Promise<Screen>
```

#### Conteúdo

**Reload()** - Recarrega a página

```typescript
function Reload(): Promise<void>
```

**ForceReload()** - Força o recarregamento da página (limpa o cache)

```typescript
function ForceReload(): Promise<void>
```

#### Zoom

**SetZoom(level)** - Define o nível de zoom

```typescript
function SetZoom(level: number): Promise<void>
```

**GetZoom()** - Obtém o nível de zoom

```typescript
function GetZoom(): Promise<number>
```

**ZoomIn()** - Aumenta o zoom

```typescript
function ZoomIn(): Promise<void>
```

**ZoomOut()** - Diminui o zoom

```typescript
function ZoomOut(): Promise<void>
```

**ZoomReset()** - Redefine o zoom para 100%

```typescript
function ZoomReset(): Promise<void>
```

#### Impressão

**Print()** - Abre a caixa de diálogo nativa de impressão

```typescript
function Print(): Promise<void>
```

**Exemplo:**

```javascript
import { Window } from '@wailsio/runtime'

// Open print dialog for current window
await Window.Print()
```

**Observação:** Isso abre a caixa de diálogo nativa de impressão do sistema operacional, permitindo que o usuário selecione as configurações da impressora e imprima o conteúdo atual da janela. Ao contrário de `window.print()`, que pode não funcionar em webviews, esse método usa a API nativa de impressão da plataforma.

### Área de transferência

Operações da área de transferência.

#### SetText()

Define o texto da área de transferência.

```typescript
function SetText(text: string): Promise<void>
```

**Exemplo:**

```javascript
import { Clipboard } from '@wailsio/runtime'

await Clipboard.SetText('Hello from Wails!')
```

#### Text()

Obtém o texto da área de transferência.

```typescript
function Text(): Promise<string>
```

**Exemplo:**

```javascript
import { Clipboard } from '@wailsio/runtime'

const clipboardText = await Clipboard.Text()
console.log('Clipboard:', clipboardText)
```

### Sistema

Métodos de baixo nível para comunicação direta com o backend.

#### invoke()

Envia uma mensagem bruta diretamente ao backend. Isso ignora o sistema padrão de vinculação e é tratado pelo `RawMessageHandler` nas opções do aplicativo.

```typescript
function invoke(message: any): void
```

**Exemplo:**

```javascript
import { System } from '@wailsio/runtime'

// Send a raw message to the backend
System.invoke('my-custom-message')

// Send structured data as JSON
System.invoke(JSON.stringify({ action: 'update', value: 42 }))
```

@note{type="caution"}
Essa função apenas envia a mensagem, sem aguardar resposta, e não retorna nenhum valor. Use eventos para receber respostas do backend.

@end

Para obter mais detalhes, consulte o [Guia de mensagens brutas](/guides/raw-messages/).

### Aplicativo

Métodos no nível do aplicativo.

#### Show()

Exibe todas as janelas do aplicativo.

```typescript
function Show(): Promise<void>
```

#### Hide()

Oculta todas as janelas do aplicativo.

```typescript
function Hide(): Promise<void>
```

#### Quit()

Encerra o aplicativo.

```typescript
function Quit(): Promise<void>
```

**Exemplo:**

```javascript
import { Application } from '@wailsio/runtime'

// Add quit button
document.getElementById('quit-btn').addEventListener('click', async () => {
    await Application.Quit()
})
```

### Navegador

Abra URLs no navegador padrão.

#### OpenURL()

Abre uma URL no navegador do sistema.

```typescript
function OpenURL(url: string | URL): Promise<void>
```

**Exemplo:**

```javascript
import { Browser } from '@wailsio/runtime'

await Browser.OpenURL('https://wails.io')
```

### Telas

Informações e gerenciamento de telas.

#### GetAll()

Obtém todas as telas.

```typescript
function GetAll(): Promise<Screen[]>
```

#### GetPrimary()

Obtém a tela principal.

```typescript
function GetPrimary(): Promise<Screen>
```

#### GetCurrent()

Obtém a tela ativa atual.

```typescript
function GetCurrent(): Promise<Screen>
```

**Interface Screen:**

```typescript
interface Screen {
    ID: string
    Name: string
    ScaleFactor: number
    X: number
    Y: number
    Size: { Width: number, Height: number }
    Bounds: { X: number, Y: number, Width: number, Height: number }
    WorkArea: { X: number, Y: number, Width: number, Height: number }
    IsPrimary: boolean
    Rotation: number
}
```

**Exemplo:**

```javascript
import { Screens } from '@wailsio/runtime'

// List all screens
const screens = await Screens.GetAll()
screens.forEach(screen => {
    console.log(`${screen.Name}: ${screen.Size.Width}x${screen.Size.Height}`)
})

// Get primary screen
const primary = await Screens.GetPrimary()
console.log('Primary screen:', primary.Name)
```

### Caixas de diálogo

Caixas de diálogo nativas do sistema operacional abertas pelo JavaScript.

#### Info()

Exibe uma caixa de diálogo de informações.

```typescript
function Info(options: MessageDialogOptions): Promise<string>
```

**Exemplo:**

```javascript
import { Dialogs } from '@wailsio/runtime'

await Dialogs.Info({
    Title: 'Success',
    Message: 'Operation completed successfully!'
})
```

#### Error()

Exibe uma caixa de diálogo de erro.

```typescript
function Error(options: MessageDialogOptions): Promise<string>
```

#### Warning()

Exibe uma caixa de diálogo de aviso.

```typescript
function Warning(options: MessageDialogOptions): Promise<string>
```

#### Question()

Exibe uma caixa de diálogo de pergunta com botões personalizados.

```typescript
function Question(options: MessageDialogOptions): Promise<string>
```

**Exemplo:**

```javascript
import { Dialogs } from '@wailsio/runtime'

const result = await Dialogs.Question({
    Title: 'Confirm Delete',
    Message: 'Are you sure you want to delete this file?',
    Buttons: [
        { Label: 'Delete', IsDefault: false },
        { Label: 'Cancel', IsDefault: true }
    ]
})

if (result === 'Delete') {
    // Delete the file
}
```

#### OpenFile()

Exibe uma caixa de diálogo para abrir arquivos.

```typescript
function OpenFile(options: OpenFileDialogOptions): Promise<string | string[]>
```

**Exemplo:**

```javascript
import { Dialogs } from '@wailsio/runtime'

const file = await Dialogs.OpenFile({
    Title: 'Select Image',
    Filters: [
        { DisplayName: 'Images', Pattern: '*.png;*.jpg;*.jpeg' },
        { DisplayName: 'All Files', Pattern: '*.*' }
    ]
})

if (file) {
    console.log('Selected:', file)
}
```

#### SaveFile()

Exibe uma caixa de diálogo para salvar arquivos.

```typescript
function SaveFile(options: SaveFileDialogOptions): Promise<string>
```

### WML (Wails Markup Language)

A WML fornece atributos declarativos para ações comuns. Adicione atributos aos elementos HTML:

#### Atributos

**wml-event** - Emite um evento ao receber um clique

```html
<button wml-event="save-clicked">Save</button>
```

**wml-window** - Chama um método da janela

```html
<button wml-window="Close">Close Window</button>
<button wml-window="Minimise">Minimize</button>
```

**wml-target-window** - Especifica a janela de destino de wml-window

```html
<button wml-window="Show" wml-target-window="settings">
    Show Settings
</button>
```

**wml-openurl** - Abre uma URL no navegador

```html
<a href="#" wml-openurl="https://wails.io">Visit Wails</a>
```

**wml-confirm** - Exibe uma caixa de diálogo de confirmação antes da ação

```html
<button wml-window="Close" wml-confirm="Are you sure you want to close?">
    Close
</button>
```

**Exemplo:**

```html
<div>
    <button wml-event="save-clicked">Save</button>
    <button wml-window="Minimise">Minimize</button>
    <button wml-window="Close" wml-confirm="Close window?">Close</button>
    <a href="#" wml-openurl="https://github.com/wailsapp/wails">GitHub</a>
</div>
```

## Exemplo completo

```javascript
import { Events, Window, Clipboard, Dialogs, Screens } from '@wailsio/runtime'

// Listen for events from Go
Events.On('data-updated', (event) => {
    console.log('Data:', event.data)
    updateUI(event.data)
})

// Window management
document.getElementById('center-btn').addEventListener('click', async () => {
    await Window.Center()
})

document.getElementById('fullscreen-btn').addEventListener('click', async () => {
    const isFullscreen = await Window.IsFullscreen()
    if (isFullscreen) {
        await Window.UnFullscreen()
    } else {
        await Window.Fullscreen()
    }
})

// Clipboard operations
document.getElementById('copy-btn').addEventListener('click', async () => {
    await Clipboard.SetText('Copied from Wails!')
})

// Dialog with confirmation
document.getElementById('delete-btn').addEventListener('click', async () => {
    const result = await Dialogs.Question({
        Title: 'Confirm',
        Message: 'Delete this item?',
        Buttons: [
            { Label: 'Delete' },
            { Label: 'Cancel', IsDefault: true }
        ]
    })

    if (result === 'Delete') {
        await Events.Emit('delete-item', { id: currentItemId })
    }
})

// Screen information
const screens = await Screens.GetAll()
console.log(`Detected ${screens.length} screen(s)`)
screens.forEach(screen => {
    console.log(`- ${screen.Name}: ${screen.Size.Width}x${screen.Size.Height}`)
})
```

## Práticas recomendadas

### ✅ Faça

- **Importe seletivamente** - Importe apenas o necessário
- **Trate as promises** - Todos os métodos retornam promises
- **Use WML para ações simples** — Mais simples do que JavaScript
- **Verifique os valores retornados** — Especialmente em caixas de diálogo
- **Cancele a inscrição em eventos** — Faça a limpeza ao terminar

### ❌ Não faça

- **Não se esqueça de usar await** — A maioria dos métodos é assíncrona
- **Não bloqueie a interface** — Use async/await corretamente
- **Não ignore erros** — Sempre trate as rejeições

## Suporte a TypeScript

O runtime inclui definições completas de TypeScript:

```typescript
import { Events, Window } from '@wailsio/runtime'

Events.On('custom-event', (event) => {
    // TypeScript knows event.data, event.name, event.sender
    console.log(event.data)
})

// All methods are fully typed
const size: { width: number, height: number } = await Window.Size()
```
