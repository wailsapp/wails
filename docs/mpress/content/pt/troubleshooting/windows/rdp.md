---
title: "WebView2 trava via Área de Trabalho Remota (RDP)"
description: "Corrija travamentos de vários segundos na interface do WebView2 quando um aplicativo Wails é executado em uma sessão RDP que altera o DPI do monitor durante a sessão."
slug: "troubleshooting/windows/rdp"
sourcePath: "troubleshooting/windows/rdp.md"
---

## Problema

Quando um aplicativo Wails é usado em uma sessão de Área de Trabalho Remota (RDP), a interface pode travar por vários segundos durante interações comuns:

- A abertura de uma janela pop-up leva cerca de 4 a 8 segundos entre o clique e a exibição do conteúdo.
- O fechamento de uma janela bloqueia a janela pai por cerca de 2 segundos.
- O estado de lentidão persiste após reconexões e só desaparece depois que a máquina host é reiniciada.

Isso é observado com mais frequência no cliente Microsoft Remote Desktop para iOS, que provisiona um monitor virtual otimizado para Retina durante a sessão. Qualquer cliente RDP que introduza, durante a sessão, um monitor com um contexto de DPI diferente pode provocar o mesmo comportamento.

## Por que isso acontece

Por padrão, o WebView2 usa hospedagem em janela, na qual sua superfície do compositor reside em uma janela filha. Quando um cliente RDP introduz um monitor cujo contexto de DPI é diferente do contexto da sessão, cada chamada ao controlador do WebView2 (`PutIsVisible`, `MoveFocus`, primeira renderização e liberação da superfície) força um remarshal síncrono do DirectComposition. Cada remarshal bloqueia a thread da interface por cerca de 2 segundos, razão pela qual os travamentos se acumulam em aplicativos que usam muitos pop-ups.

Um aplicativo nativo Win32 com WebView2 na mesma máquina não é afetado, pois usa hospedagem visual. Isso indica que a causa é o modo de hospedagem, e não um problema geral do WebView2 ou do compositor do Windows.

## Solução

Ative a hospedagem visual definindo `UseVisualHosting` nas opções do Windows. Com a hospedagem visual, a superfície do compositor do WebView2 passa a pertencer a um visual do DirectComposition controlado pelo host; assim, as alterações no contexto de DPI deixam de provocar o remarshal síncrono.

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Windows: application.WindowsOptions{
            UseVisualHosting: true,
        },
    })

    // ... create your windows, then:
    app.Run()
}
```

Após a ativação, os pop-ups abrem dentro do tempo normal de navegação (aproximadamente 150 a 500 ms), e o fechamento de uma janela deixa de bloquear a janela pai.

@note{type="caution"}
`UseVisualHosting` deve ser definido antes de `app.Run()`. O Wails lê essa opção durante a inicialização do aplicativo e define a variável de ambiente `COREWEBVIEW2_FORCED_HOSTING_MODE` como `COREWEBVIEW2_HOSTING_MODE_WINDOW_TO_VISUAL` antes que o ambiente do WebView2 seja inicializado. Defini-la posteriormente não produz efeito.

@end

O valor padrão da opção é `false`, portanto, a hospedagem em janela continua sendo o padrão. Os aplicativos existentes não apresentam nenhuma mudança de comportamento, a menos que habilitem explicitamente a opção.

## Quando ativar essa opção

Defina `UseVisualHosting: true` se o aplicativo for usado regularmente via RDP, especialmente com o cliente Microsoft Remote Desktop para iOS, e você observar travamentos de vários segundos ao abrir ou fechar janelas. Se o aplicativo não for executado via RDP, essa opção não é necessária e você pode manter o valor padrão.

## Referências

- [WebView2: hospedagem em janela versus hospedagem visual](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/windowed-vs-visual-hosting)
- [Issue nº 5248](https://github.com/MicrosoftEdge/WebView2Feedback/issues/5248) do WebView2Feedback
- [Issue nº 4485](https://github.com/MicrosoftEdge/WebView2Feedback/issues/4485) do WebView2Feedback
