---
title: "Atalhos globais"
description: "Registre atalhos de teclado disponíveis em todo o sistema, acionados mesmo quando seu aplicativo não está em foco"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

Atalhos globais são atalhos de teclado disponíveis em todo o sistema e acionados independentemente do aplicativo que está em foco no momento, enquanto seu aplicativo Wails estiver em execução. Eles são ideais para teclas de atalho de exibir/ocultar, ferramentas de captura rápida, controles de mídia e outros recursos que os usuários esperam poder acessar de qualquer lugar.

@note{type="info" title="Atalhos globais versus associações de teclas"}
[Associações de teclas](/features/keyboard/shortcuts/) (`app.KeyBinding`) só são acionadas enquanto uma das janelas do seu aplicativo está em foco. Atalhos globais (`app.GlobalShortcut`) são acionados em todo o sistema, mesmo quando seu aplicativo está em segundo plano. Use a opção adequada à sua necessidade.

@end

Os atalhos globais são implementados diretamente com os recursos nativos de cada plataforma e não adicionam dependências de terceiros.

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

Você pode registrar atalhos antes de chamar `app.Run`. A associação com o sistema operacional é então realizada automaticamente quando o aplicativo é iniciado.

@note{type="tip" title="Como interagir com a interface do usuário em um callback"}
Os callbacks são executados fora da thread principal. Se o callback precisar interagir com janelas ou outros elementos da interface do usuário, os métodos da janela cuidarão disso para você; porém, para executar trabalho personalizado na thread principal, envolva-o com `application.InvokeSync`.

@end

### Formato do acelerador

Os atalhos globais usam o mesmo formato de acelerador que os aceleradores de menu e as associações de teclas:

```go
"CmdOrCtrl+Shift+G"  // Command on macOS, Control elsewhere
"Ctrl+Alt+K"         // Control + Alt + K
"Cmd+Option+Space"   // Command + Option + Space (macOS)
"Super+D"            // Super / Windows / Logo key + D
"Ctrl+Shift+F5"      // Function keys are supported
```

`CmdOrCtrl` é resolvido como Command no macOS e Control no Windows e no Linux, o que facilita a criação de atalhos multiplataforma.

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

Há dois casos distintos, e o Wails trata cada um de forma diferente.

### O mesmo aplicativo registra um atalho duas vezes

Esse caso é resolvido pelo próprio Wails e apresenta o mesmo comportamento em todas as plataformas. A segunda chamada a `Register` retorna um erro, e a associação original é mantida ("erro e preservação"). Isso mantém o comportamento previsível e revela o engano, em vez de substituir silenciosamente um atalho funcional.

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

Se quiser alterar o callback de um atalho, primeiro chame `Unregister` para removê-lo e depois chame `Register` para registrá-lo novamente.

### Outro aplicativo já é o proprietário do atalho

Esse caso é decidido pelo sistema operacional, portanto o resultado varia conforme a plataforma:

| Plataforma | Comportamento quando outro aplicativo é o proprietário do atalho |
| --- | --- |
| **macOS** | O registro é bem-sucedido. O macOS permite que vários aplicativos registrem a mesma tecla de atalho, portanto seu callback é adicionado junto ao proprietário existente, em vez de ser rejeitado. |
| **Windows** | O registro falha, e `Register` retorna um erro. O aplicativo que registrou o atalho primeiro continua sendo seu proprietário. |
| **Linux (X11)** | O registro falha, e `Register` retorna um erro, pois o servidor X recusa uma segunda captura da mesma combinação. |
| **Linux (Wayland)** | O compositor arbitra o conflito. Normalmente, o usuário precisa aprovar ou escolher a associação na caixa de diálogo de atalhos globais do ambiente de desktop. |

Devido a essas diferenças, sempre verifique o erro retornado por `Register` e forneça um atalho alternativo ou uma mensagem ao usuário quando não for possível obter a combinação.

## Considerações específicas de cada plataforma

@tabs
[macOS]
Os atalhos globais usam a API de teclas de atalho do Carbon Event Manager. Esse é o mecanismo padrão para teclas de atalho disponíveis em todo o sistema no macOS e não requer permissão de Acessibilidade.

As teclas de atalho são associadas às posições físicas das teclas; portanto, em layouts que não sejam QWERTY, um atalho corresponde à tecla situada na posição padrão ANSI/QWERTY.

@note{type="caution" title="Atalhos para ocultar e `ApplicationShouldTerminateAfterLastWindowClosed`"}
`window.Hide()` usa `orderOut:` no macOS, o que torna a janela não visível. O AppKit considera fechada a última janela não visível; portanto, se você definir `Mac.ApplicationShouldTerminateAfterLastWindowClosed: true` e usar um atalho global para ocultar sua única janela, o aplicativo será encerrado em vez de permanecer em segundo plano. Quando depender de uma tecla de atalho para ocultar/exibir, deixe essa opção sem definir (o padrão), para que a janela possa ser ocultada e reaberta posteriormente.

@end

[Windows]
Os atalhos globais usam a API Win32 `RegisterHotKey`. A repetição automática é suprimida; portanto, manter as teclas pressionadas aciona o callback uma única vez, em vez de acioná-lo repetidamente.

O registro falha se outro aplicativo já for o proprietário da combinação; portanto, prefira combinações menos comuns como padrão.

[Linux]
Em sessões **X11**, o Wails captura o atalho diretamente do servidor X, portanto o acelerador solicitado é associado exatamente como especificado.

Em sessões **Wayland**, por definição, não há como um aplicativo capturar teclas diretamente. Em vez disso, o Wails usa a interface `org.freedesktop.portal.GlobalShortcuts` do XDG Desktop Portal. Com o portal, o acelerador fornecido é um acionador *preferencial*, e o compositor (e, em última instância, o usuário) decide a combinação de teclas final. Seu callback ainda é executado quando o atalho é ativado, mas não há garantia de que as teclas exatas correspondam à sua solicitação, e `IsRegistered`/`GetAll` informam o que você solicitou, não o que o compositor associou.

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
Algumas combinações são reservadas pelo sistema operacional ou pelo ambiente de desktop e não podem ser capturadas por aplicativos. Escolha combinações padrão com pouca probabilidade de conflito e sempre trate o erro retornado por `Register`.

@end
