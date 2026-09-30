---
title: "APIs privadas do macOS"
description: "Todos os recursos e opções do Wails que dependem de APIs privadas do macOS, com comandos de ativação e alternativas para builds públicos."
slug: "guides/build/private-macos-apis"
sourcePath: "guides/build/private-macos-apis.md"
---

Por padrão, o Wails usa APIs públicas do macOS. A única tag de build do Go, `private_mac_apis`, habilita as chamadas privadas do WebKit e do AppKit listadas nesta página. Todas as opções e todos os métodos públicos do Go permanecem disponíveis em ambos os builds. Sem a tag, as operações exclusivas de APIs privadas não realizam nenhuma ação; os recursos com alternativas públicas usam essas alternativas.

@note{type="caution" title="Ativar o comportamento privado do macOS"}
Definir uma opção de janela não habilita APIs privadas. Adicione `private_mac_apis` ao comando de build para habilitá-las. A tag se aplica somente a builds de desktop para macOS, não a builds para iOS, Android, Windows, Linux ou servidor.

@end

## Habilitar APIs privadas

```bash
# Build an application
wails3 build -tags private_mac_apis

# Run with live reload
EXTRA_TAGS=private_mac_apis wails3 dev

# Package an application
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis

# Run a Go-only example from its directory
go run -tags private_mac_apis .
```

Para um build de produção direto, use `go build -tags production,private_mac_apis .`. Para os exemplos de frontend, siga o README correspondente para compilar os bindings e os assets antes da execução. Taskfiles personalizados ou antigos devem encaminhar `EXTRA_TAGS` ao compilador Go.

## Inventário de recursos

| Recurso ou valor | O que `private_mac_apis` habilita | Sem a tag |
| --- | --- | --- |
| `Mac.Backdrop: MacBackdropTransparent` | WKWebView transparente sobre a janela nativa | A janela nativa é configurada, mas a webview permanece opaca |
| `Mac.Backdrop: MacBackdropTranslucent` | WKWebView transparente para que o desfoque nativo fique visível através dela | O desfoque é configurado atrás de uma webview opaca |
| `Mac.Backdrop: MacBackdropLiquidGlass` | WKWebView transparente sobre a camada de vidro, com controle privado do plano de fundo da webview | A camada de vidro é configurada atrás de uma webview opaca; a estilização usa alternativas públicas |
| Limpeza do plano de fundo da webview durante a configuração do Liquid Glass | Controle privado `backgroundColor` do WebKit | `underPageBackgroundColor` público no macOS 12+ ou uma cor de camada em versões anteriores do macOS; não torna a webview transparente |
| `app.Window.NewNotchWindow(...)` | Webview transparente dentro do painel de notch com formato personalizado | O painel continua funcionando, inclusive quanto ao posicionamento e às animações, mas sua webview permanece opaca |
| `Mac.LiquidGlass.Style` | Mapeamento de estilos nativos existente do Wails, incluindo um valor de estilo escuro não documentado | Usa estilos públicos regular/clear e aparências clara/escura; consulte a tabela de valores abaixo |
| `Mac.LiquidGlass.GroupID` | Solicita o agrupamento privado de vidro para um identificador não vazio | Ignorado; nenhum agrupamento é solicitado |
| `Mac.LiquidGlass.GroupSpacing` | Solicita o espaçamento privado de grupo para um valor maior que zero | Ignorado |
| `window.OpenDevTools()` e `Window.OpenDevTools()` do JavaScript | Abre programaticamente o inspetor do WebKit no macOS 12+ | Não realiza nenhuma ação |
| `WebviewWindowOptions.OpenInspectorOnStartup: true` | Solicita a abertura programática do inspetor quando a janela é exibida pela primeira vez, no macOS 12+ | Não realiza nenhuma ação |
| Habilitação legada do inspetor, anterior ao macOS 13.3 | Habilita os recursos extras para desenvolvedores do WebKit quando o suporte ao inspetor está incluído no build | Não realiza nenhuma ação; a inspeção pública pelo Safari requer macOS 13.3+ |

## Transparência e plano de fundo da webview

**Requer APIs privadas:** a transparência da webview usada por `MacBackdropTransparent`, `MacBackdropTranslucent`, `MacBackdropLiquidGlass` e janelas de notch. Internamente, o Wails define a chave privada `drawsBackground` do WebKit. Um plano de fundo HTML ou CSS transparente, por si só, não pode tornar transparente uma WKWebView nativa opaca.

```go
Mac: application.MacWindow{
    // Requires -tags private_mac_apis for the blur to show through the webview.
    Backdrop: application.MacBackdropTranslucent,
},
```

A operação privada de cor do plano de fundo da webview usa a chave `backgroundColor` do WebKit. A configuração do Liquid Glass usa essa operação para limpar o plano de fundo da webview. Sem a tag, essa operação interna usa `underPageBackgroundColor` público no macOS 12+ ou a camada da view em versões anteriores do macOS. Essas alternativas não tornam a webview transparente.

`WebviewWindowOptions.BackgroundColour` e `window.SetBackgroundColour()` definem a cor da **janela nativa** no macOS e, por si só, não requerem APIs privadas. Da mesma forma, `Frameless` e `Mac.TitleBar.AppearsTransparent` usam APIs públicas do AppKit; a dependência privada é a transparência da webview, não a transparência da barra de título. Para aplicar efeitos de pano de fundo no macOS, configure `Mac.Backdrop` em vez de depender apenas de `BackgroundType`.

Consulte [opções de janela](/features/windows/options/#mac-options), [janelas sem moldura](/features/windows/frameless/#with-transparent-background) e [janelas de notch](/features/windows/notch-windows/).

## Valores do Liquid Glass

Quando um `NSGlassEffectView` nativo está disponível (macOS 26+), aplicam-se os mapeamentos a seguir. Somente os valores de estilo nativo `0` (regular) e `1` (clear) estão documentados. As constantes Go mantêm os valores existentes em ambos os builds.

| Valor de `MacLiquidGlassStyle` | Com `private_mac_apis` | Sem a tag |
| --- | --- | --- |
| `LiquidGlassStyleAutomatic` (`0`) | Estilo nativo regular (`0`) | Estilo nativo regular (`0`) |
| `LiquidGlassStyleLight` (`1`) | Mapeamento de estilo nativo existente (`1`, clear) | Estilo nativo regular (`0`) com aparência Aqua |
| `LiquidGlassStyleDark` (`2`) | **Valor de estilo nativo não documentado `2`** | Estilo regular nativo (`0`) com a aparência Dark Aqua |
| `LiquidGlassStyleVibrant` (`3`) | Mapeia para o estilo claro/transparente nativo existente (`1`) | Estilo transparente nativo (`1`) |

Automatic e Vibrant usam valores de estilo nativo documentados, mas um plano de fundo Liquid Glass que ocupe toda a janela ainda precisa da tag para a **transparência da webview**. A aparência de Light difere entre as duas compilações. O mapeamento de estilo privado não garante que uma versão futura do macOS renderize o mesmo efeito.

**Sempre privados:** `GroupID` e `GroupSpacing`. Antes de solicitar o agrupamento, o Wails verifica os seletores privados `setGroupIdentifier:`, `setGroupName:` e `setGroupSpacing:`. Sem a tag, essas operações não têm efeito. Habilitar a tag não garante que a versão do macOS em execução seja compatível com esses seletores.

`MacLiquidGlass.Material`, `CornerRadius` e `TintColor` não exigem APIs privadas por si só. Em versões do macOS sem suporte nativo a Liquid Glass, o Wails configura uma alternativa translúcida; para que ela fique visível através da webview, a tag ainda é necessária.

## Inspetor da Web

**Requer APIs privadas:** chamar `OpenDevTools()` ou definir `OpenInspectorOnStartup: true` para abrir o inspetor do WebKit a partir do aplicativo. O Wails usa o seletor privado `_inspector`. Sem `private_mac_apis`, essas operações não têm efeito e não exibem nenhuma mensagem.

O suporte ao inspetor também precisa ser incluído na compilação. As tags `production` e `devtools` existentes mantêm seus significados:

| Tags de compilação | Abertura programática do inspetor | Inspeção pública pelo Safari no macOS 13.3+ |
| --- | --- | --- |
| Nenhuma | Sem efeito | Habilitada |
| `private_mac_apis` | Habilitada no macOS 12+ | Habilitada |
| `production` | Sem efeito | Desabilitada |
| `production,private_mac_apis` | Sem efeito | Desabilitada |
| `production,devtools` | Sem efeito | Habilitada |
| `production,devtools,private_mac_apis` | Habilitada no macOS 12+ | Habilitada |

No macOS 13.3 ou posterior, o Wails habilita a inspeção pelo Safari usando a API pública `WKWebView.inspectable`; isso não requer APIs privadas. Antes do macOS 13.3, o mecanismo alternativo usa a preferência privada `developerExtrasEnabled` e, portanto, requer `private_mac_apis`, além de o suporte ao inspetor ter sido incluído na compilação.

```bash
# Production build with programmatic inspector support
go build -tags production,devtools,private_mac_apis .
```

## Manutenção desta lista

Todas as chamadas nativas privadas ficam isoladas em `v3/pkg/application/mac_private_api_darwin.go`; as compilações padrão selecionam `mac_public_api_darwin.go`. Este inventário abrange transparência, cor de fundo da webview, estilos de vidro, agrupamento de vidro, abertura do inspetor e habilitação do inspetor legado. Alterações nessas implementações devem atualizar esta página e, ao mesmo tempo, a documentação das opções afetadas.
