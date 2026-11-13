---
title: "Janelas de notch"
description: "Crie janelas nativas do macOS anexadas ao alojamento da câmera."
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow` cria um painel do macOS com formato personalizado e que não ativa o aplicativo, anexado ao alojamento da câmera. O Wails controla o posicionamento nativo, as extensões laterais pretas, a tela externa transparente, o nível da janela, o comportamento nos Spaces e as animações opcionais de exibição e ocultação. Seu conteúdo web ocupa somente o retângulo interno solicitado. Mover o ponteiro para dentro da janela torna sua webview imediatamente apta a receber entrada do teclado, sem ativar o aplicativo.

@note{type="caution" title="A transparência da webview usa uma API privada"}
`NewNotchWindow` requer `-tags private_mac_apis` para que o conteúdo web transparente revele o formato nativo do notch. Sem a tag, o painel, o posicionamento e as animações continuam funcionando, mas a webview permanece opaca. Consulte [APIs privadas do macOS](/guides/build/private-macos-apis/#webview-transparency-and-background).

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="Uma notificação de notch do Wails deslizando para baixo a partir do alojamento da câmera do MacBook, exibindo métricas do sistema em tempo real e ocultando-se novamente sob o notch"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    Uma notificação de notch animada que usa uma webview persistente com transições nativas
    de exibição e ocultação.
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
| `Width` | `int` | `660` | Largura utilizável da webview dentro das bordas com formato nativo. |
| `Height` | `int` | `92` | Altura utilizável da webview dentro das bordas com formato nativo. |
| `Animated` | `bool` | `false` | Desliza a janela para baixo em `Show` e para cima em `Hide`. |
| `AnimationSpeed` | `time.Duration` | `420ms` | Duração da exibição. A ocultação usa dois terços desse valor, `280ms` por padrão. |
| `Screen` | `*Screen` | tela principal | Define uma tela específica como destino. Caso contrário, o Wails usa a tela principal. |
| `WindowOptions` | `WebviewWindowOptions` | valores padrão | Fornece o nome, a URL ou o HTML, o CSS, o JavaScript, os atalhos de teclado e outros comportamentos da webview. |

`NewNotchWindow` controla o tamanho externo, a posição, o quadro, a transparência, a política de redimensionamento, a classe do painel nativo, o nível da janela e o comportamento de coleção. Os valores desses campos em `WindowOptions` são substituídos intencionalmente. Os outros campos são preservados. O arraste nativo pelo plano de fundo e as regiões de arraste CSS são desativados para que a janela permaneça anexada ao alojamento da câmera.

O `NotchWindow` retornado expõe intencionalmente apenas `Show`, `Hide`, `Visibility` e `Close`; a geometria nativa não pode ser alterada por meio do identificador de alto nível.

@note{type="note"}
As janelas de notch exigem macOS. Em um Mac sem alojamento de câmera, o Wails posiciona a janela na parte superior central, abaixo da barra de menus. Em plataformas sem suporte, `NewNotchWindow` retorna um identificador inerte, cujos métodos de ciclo de vida são operações nulas seguras e cujo `Visibility` é sempre falso.

@end

## Ciclo de vida

- `Show` revela a janela nativa existente. Com a animação ativada, ela desliza para baixo a partir de uma posição acima da tela.
- `Hide` mantém a janela ativa para reutilização. Com a animação ativada, ela desliza de volta para uma posição acima da tela antes de ser retirada da exibição. A webview, o estado do JavaScript, os bindings e os listeners de eventos permanecem carregados, mas a janela oculta não tem uma área de destino para o ponteiro; o aplicativo deve chamar `Show` para revelá-la novamente.
- `Visibility` informa a visibilidade nativa atual.
- `Close` destrói permanentemente a janela nativa. Crie uma nova antes de exibir essa notificação novamente.
- A entrada do ponteiro traz essa janela de notch para a frente e dá foco à sua webview para interação imediata pelo teclado, preservando o comportamento do painel que não ativa o aplicativo.

Cada chamada a `NewNotchWindow` cria uma janela independente, com conteúdo, visibilidade e estado de animação próprios. As janelas usam o mesmo nível nativo; portanto, a instância exibida mais recentemente ou na qual o ponteiro entrou por último aparece à frente e pode se sobrepor às instâncias anteriores. Quando as janelas têm telas diferentes como destino, cada uma é centralizada no alojamento da câmera da respectiva tela. O macOS não coordena janelas de notch pertencentes a aplicativos diferentes; se aplicativos distintos exibirem janelas no mesmo local e nível, a janela ordenada mais recentemente aparecerá à frente.

Para cargas de trabalho de notificação, os aplicativos geralmente reutilizam uma janela oculta ou mantêm sua própria fila em uma única webview. O Wails não impõe enfileiramento, substituição, fechamento automático nem uma política de janela única.

@note{type="note"}
Superfícies recolhidas persistentes que reabrem ao passar o ponteiro são intencionalmente separadas da ocultação de notificações e são acompanhadas em [#6009](https://github.com/wailsapp/wails/issues/6009).

@end

Consulte o [exemplo de notificação de notch](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/notch-notification) para ver um monitor de sistema compacto cujo estado ativo do JavaScript é preservado durante ciclos repetidos de exibição e ocultação da notificação.
