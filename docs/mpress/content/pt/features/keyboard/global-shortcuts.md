---
title: "Atalhos globais"
description: "Registre atalhos de teclado disponíveis em todo o sistema, que são acionados mesmo quando seu aplicativo não está em foco"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

Atalhos globais são atalhos de teclado disponíveis em todo o sistema que são acionados independentemente de qual aplicativo esteja em foco, enquanto seu aplicativo Wails estiver em execução. Eles são ideais para teclas de atalho de exibir/ocultar, ferramentas de captura rápida, controles de mídia e outros recursos que os usuários esperam acessar de qualquer lugar.

@note{type="info" title="Atalhos globais versus associações de teclas"}
[As associações de teclas](/features/keyboard/shortcuts/) (`app.KeyBinding`) só são acionadas enquanto uma das janelas do seu aplicativo está em foco. Os atalhos globais (`app.GlobalShortcut`) são acionados em todo o sistema, mesmo quando seu aplicativo está em segundo plano. Use a opção que atenda à sua necessidade.

@end

Os atalhos globais usam diretamente os recursos nativos de cada plataforma e não adicionam dependências de terceiros.

## Como acessar o gerenciador de atalhos globais

O gerenciador está disponível na propriedade `GlobalShortcut` da instância do seu aplicativo:

```go
app := application.New(application.Options{
    Name: "Global Shortcuts Demo",
})

globalShortcuts := app.GlobalShortcut
```

## Como registrar um atalho

`Register` recebe um acelerador e um callback. O callback é executado em sua própria goroutine sempre que o atalho é pressionado.

```go
err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
    // Runs even when another application is focused.
    window.Show()
    window.Focus()
})
if err != nil {
    app.Logger.Error("could not register shortcut", "error", err)
}
```

Você pode registrar atalhos antes de chamar `app.Run`. A associação com o sistema operacional é feita automaticamente quando o aplicativo é iniciado.

@note{type="tip" title="Como acessar a interface do usuário em um callback"}
Os callbacks são executados fora da thread principal. Se o callback precisar interagir com janelas ou outros elementos da interface do usuário, os métodos de janela cuidam disso para você; para executar uma tarefa personalizada na thread principal, porém, encapsule-a com `application.InvokeSync`.

@end

### Formato do acelerador

Os atalhos globais usam o mesmo formato de acelerador dos aceleradores de menu e das associações de teclas:

```go
"CmdOrCtrl+Shift+G"  // Command on macOS, Control elsewhere
"Ctrl+Alt+K"         // Control + Alt + K
"Cmd+Option+Space"   // Command + Option + Space (macOS)
"Super+D"            // Super / Windows / Logo key + D
"Ctrl+Shift+F5"      // Function keys are supported
```

`CmdOrCtrl` corresponde a Command no macOS e a Control no Windows e no Linux, o que o torna conveniente para atalhos multiplataforma.

## Como gerenciar atalhos

```go
// Check whether a shortcut is registered (modifier order does not matter).
registered := app.GlobalShortcut.IsRegistered("Ctrl+Shift+G")

// List every shortcut this application has registered.
for _, accelerator := range app.GlobalShortcut.GetAll() {
    app.Logger.Info("global shortcut", "accelerator", accelerator)
}

// Release a single shortcut.
app.GlobalShortcut.Unregister("Ctrl+Shift+G")

// Release everything (also done automatically on shutdown).
app.GlobalShortcut.UnregisterAll()
```

Todos os atalhos registrados são liberados automaticamente quando o aplicativo é encerrado, portanto não é necessário removê-los manualmente.

## O que acontece quando o mesmo atalho é registrado duas vezes

Há dois casos distintos, e o Wails os trata de maneiras diferentes.

### O mesmo aplicativo registra um atalho duas vezes

Esse caso é resolvido pelo próprio Wails e se comporta da mesma forma em todas as plataformas. A segunda chamada a `Register` retorna um erro, e a associação original é mantida ("erro e preservação"). Isso mantém o comportamento previsível e evidencia o engano, em vez de substituir silenciosamente um atalho que funciona.

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

Para alterar o callback de um atalho, primeiro use `Unregister` para cancelar seu registro e depois use `Register` para registrá-lo novamente.

### Outro aplicativo já detém o atalho

Isso é decidido pelo sistema operacional, portanto o resultado varia conforme a plataforma:

| Plataforma | Comportamento quando outro aplicativo detém o atalho |
| --- | --- |
| **macOS** | O registro é bem-sucedido. O macOS permite que vários aplicativos registrem a mesma tecla de atalho; portanto, seu callback é adicionado junto ao do aplicativo que já a detém, em vez de ser rejeitado. |
| **Windows** | O registro falha, e `Register` retorna um erro. O aplicativo que registrou o atalho primeiro o mantém. |
| **Linux (X11)** | O registro falha, e `Register` retorna um erro, porque o servidor X recusa uma segunda captura da mesma combinação. |
| **Linux (Wayland)** | O compositor arbitra. Normalmente, o usuário precisa aprovar ou escolher a associação na caixa de diálogo de atalhos globais do ambiente de desktop. |

Devido a essas diferenças, sempre verifique o erro retornado por `Register` e forneça um atalho alternativo ou informe o usuário quando não for possível reservar um atalho.

## Considerações específicas de cada plataforma

@tabs
[macOS]
Os atalhos globais usam a API de teclas de atalho do Carbon Event Manager. Esse é o mecanismo padrão para teclas de atalho disponíveis em todo o sistema no macOS e não requer permissão de Acessibilidade.

Os atalhos com letras levam em conta o layout. A API Carbon faz a associação com posições físicas das teclas; por isso, o Wails traduz uma letra para a tecla que a produz no layout de teclado ativo (por meio de `UCKeyTranslate`). Assim, em layouts AZERTY, QWERTZ, Dvorak e semelhantes, `Cmd+A` corresponde à tecla identificada como **A**, e não à posição dela no QWERTY. Dígitos, sinais de pontuação e teclas nomeadas (setas, teclas de função e assim por diante) usam suas posições físicas fixas. A associação é determinada quando o atalho é registrado; se você mudar o layout do teclado enquanto o aplicativo estiver em execução, os atalhos de letras existentes manterão a tecla física à qual foram associados originalmente.

@note{type="caution" title="Atalhos para ocultar e `ApplicationShouldTerminateAfterLastWindowClosed`"}
`window.Hide()` usa `orderOut:` no macOS, o que torna a janela não visível. O AppKit considera fechada a última janela não visível; portanto, se você definir `Mac.ApplicationShouldTerminateAfterLastWindowClosed: true` e usar um atalho global para ocultar sua única janela, o aplicativo será encerrado em vez de permanecer em segundo plano. Quando depender de uma tecla de atalho para ocultar/exibir, deixe essa opção sem definir (o padrão), para que a janela possa ser ocultada e reaberta posteriormente.

@end

[Windows]
Os atalhos globais usam a API Win32 `RegisterHotKey`. A repetição automática é suprimida; portanto, manter as teclas pressionadas aciona o callback uma única vez, e não repetidamente.

O registro falha se outro aplicativo já detiver a combinação; portanto, prefira combinações menos comuns como padrão.

[Linux]
Em sessões **X11**, o Wails captura o atalho diretamente do servidor X; portanto, o acelerador solicitado é associado exatamente como especificado.

Em sessões **Wayland**, por projeto, não há como um aplicativo capturar teclas diretamente. Em vez disso, o Wails usa a interface `org.freedesktop.portal.GlobalShortcuts` do XDG Desktop Portal. Com o portal, o acelerador informado é um acionador *preferencial*, e o compositor (e, em última instância, o usuário) decide a combinação final de teclas. Seu callback ainda é executado quando o atalho é ativado, mas não há garantia de que as teclas exatas correspondam à solicitação, e `IsRegistered`/`GetAll` informam o que você solicitou, não o que o compositor associou.

O portal requer um ambiente de desktop que implemente o portal de atalhos globais (por exemplo, versões recentes do GNOME ou do KDE Plasma).

@end

## Exemplo completo

```go
package main

import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Global Shortcuts Demo",
    })

    window := app.Window.New()

    // Bring the window to the front from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
        window.Show()
        window.Focus()
    }); err != nil {
        log.Printf("could not register show shortcut: %v", err)
    }

    // Hide the window from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+H", func() {
        window.Hide()
    }); err != nil {
        log.Printf("could not register hide shortcut: %v", err)
    }

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

@note{type="danger" title="Evite atalhos críticos do sistema"}
Algumas combinações são reservadas pelo sistema operacional ou pelo ambiente de desktop e não podem ser usadas por aplicativos. Escolha padrões com pouca probabilidade de conflito e sempre trate o erro de `Register`.

@end
