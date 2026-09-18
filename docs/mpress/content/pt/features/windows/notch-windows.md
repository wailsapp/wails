---
title: "Janelas do notch"
description: "Crie janelas nativas do macOS anexadas ao compartimento da câmera."
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow` cria um painel do macOS com formato personalizado e que não ativa o aplicativo, anexado ao compartimento da câmera. O Wails controla o posicionamento nativo, as extensões pretas de encaixe, a área externa transparente, o nível da janela, o comportamento nos Spaces e as animações opcionais para mostrar e ocultar. Seu conteúdo web ocupa apenas o retângulo interno solicitado. Ao mover o ponteiro para dentro da janela, sua webview se torna imediatamente a principal receptora de entrada, sem ativar o aplicativo.

@note{type="caution" title="A transparência da webview usa uma API privada"}
`NewNotchWindow` requer `-tags private_mac_apis` para que o conteúdo web transparente revele o formato nativo do notch. Sem a tag, o painel, o posicionamento e as animações continuam funcionando, mas a webview permanece opaca. Consulte [APIs privadas do macOS](/guides/build/private-macos-apis/#webview-transparency-and-background).

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="Uma notificação de notch do Wails deslizando para baixo a partir do compartimento da câmera do MacBook, exibindo métricas do sistema em tempo real e ocultando-se novamente sob o notch"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    Uma notificação de notch animada que usa uma webview persistente com transições nativas para mostrar
    e ocultar.
  </figcaption>
</figure>

```go
alert := app.Window.NewNotchWindow(application.NotchWindowOptions{
    Width:    660,
    Height:   92,
    Animated: true,
    WindowOptions: application.WebviewWindowOptions{
        Name: "alert",
        URL:  "/alert",
    },
})

alert.Show()
alert.Hide()
visible := alert.Visibility()
alert.Close()
```

## Opções

| Campo | Tipo | Padrão | Descrição |
| --- | --- | --- | --- |
| `Width` | `int` | `660` | Largura útil da webview dentro das bordas com formato nativo. |
| `Height` | `int` | `92` | Altura útil da webview dentro das bordas com formato nativo. |
| `Animated` | `bool` | `false` | Desliza a janela para baixo em `Show` e para cima em `Hide`. |
| `AnimationSpeed` | `time.Duration` | `420ms` | Duração da exibição. A ocultação usa dois terços desse valor, `280ms` por padrão. |
| `Screen` | `*Screen` | tela principal | Define uma tela específica como destino. Caso contrário, o Wails usa a tela principal. |
| `WindowOptions` | `WebviewWindowOptions` | valores padrão | Fornece o nome, a URL ou o HTML, o CSS, o JavaScript, os atalhos de teclado e outros comportamentos da webview. |

`NewNotchWindow` controla o tamanho externo, a posição, a moldura, a transparência, a política de redimensionamento, a classe do painel nativo, o nível da janela e o comportamento de coleção. Os valores desses campos em `WindowOptions` são substituídos intencionalmente. Os demais campos são preservados. O arraste nativo pelo plano de fundo e as regiões de arraste CSS são desativados para que a janela permaneça anexada ao compartimento da câmera.

O `NotchWindow` retornado expõe intencionalmente apenas `Show`, `Hide`, `Visibility` e `Close`; a geometria nativa não pode ser alterada por meio da referência de alto nível.

@note{type="note"}
As janelas de notch requerem macOS. Em um Mac sem compartimento para a câmera, o Wails posiciona a janela no centro da parte superior, abaixo da barra de menus. Em plataformas não compatíveis, `NewNotchWindow` retorna uma referência inerte cujos métodos de ciclo de vida são operações sem efeito seguras e cujo `Visibility` é sempre falso.

@end

## Ciclo de vida

- `Show` revela a janela nativa existente. Com a animação ativada, ela desliza para baixo a partir de uma posição acima da tela.
- `Hide` mantém a janela ativa para reutilização. Com a animação ativada, ela desliza novamente para uma posição acima da tela antes de ser removida da exibição. A webview, o estado do JavaScript, os bindings e os listeners de eventos permanecem carregados, mas a janela oculta não tem uma área que responda à passagem do ponteiro; o aplicativo deve chamar `Show` para revelá-la novamente.
- `Visibility` informa a visibilidade nativa atual.
- `Close` destrói permanentemente a janela nativa. Crie outra antes de mostrar essa notificação novamente.
- A entrada do ponteiro traz essa janela de notch para a frente e direciona o foco à sua webview para permitir a interação imediata pelo teclado, preservando o comportamento do painel que não ativa o aplicativo.

Cada chamada a `NewNotchWindow` cria uma janela independente com conteúdo, visibilidade e estado de animação próprios. As janelas usam o mesmo nível nativo; portanto, a instância mostrada mais recentemente ou na qual o ponteiro entrou por último aparece à frente e pode se sobrepor às instâncias anteriores. Quando as janelas têm telas diferentes como destino, cada uma é centralizada no compartimento da câmera da respectiva tela. O macOS não coordena janelas de notch pertencentes a aplicativos diferentes; se aplicativos distintos mostrarem janelas no mesmo local e nível, a janela ordenada mais recentemente aparecerá à frente.

Em cargas de trabalho de notificações, os aplicativos normalmente reutilizam uma janela oculta ou mantêm a própria fila em uma única webview. O Wails não impõe enfileiramento, substituição, dispensa automática nem uma política de janela única.

@note{type="note"}
Superfícies recolhidas persistentes que reabrem ao passar o ponteiro são intencionalmente separadas da ocultação de notificações e são acompanhadas em [#6009](https://github.com/wailsapp/wails/issues/6009).

@end

Consulte o exemplo [notch-notification](https://github.com/wailsapp/wails/tree/master/v3/examples/notch-notification) de um monitor de sistema compacto cujo estado ativo do JavaScript persiste durante ciclos repetidos de exibição e ocultação de notificações.
