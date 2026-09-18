---
title: "Visão geral para dispositivos móveis"
description: "Crie aplicativos para iOS e Android usando a mesma base de código Go do seu aplicativo para desktop"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

O Wails v3 é executado no **iOS e no Android** usando o mesmo `main.go` e o mesmo frontend que você já desenvolve para desktop. Não há um projeto separado para dispositivos móveis, uma ponte para compartilhamento de código nem a necessidade de reescrever o código: o binário Go é compilado para a plataforma móvel de destino, e uma WebView nativa renderiza o frontend existente.

@cards{cols="2"}
iOS
Host WKWebView + UIKit. Os recursos são servidos por meio de um esquema `wails://` personalizado — sem portas abertas. Requer **macOS** com a instalação completa do Xcode.

[Guia do iOS →](/guides/mobile/ios/)

---
Android
Android WebView + `WebViewAssetLoader`. O código Go é compilado como `libwails.so` por meio do NDK. Funciona no macOS, Linux e Windows.

[Guia do Android →](/guides/mobile/android/)

@end

## Veja em execução: exemplo Kitchen Sink

A melhor maneira de entender o que é possível fazer é conferir o **Kitchen Sink** — um único aplicativo Wails que, a partir de uma só base de código, é executado de forma idêntica no iOS, Android e desktop:

@linkcard{title="Mobile Kitchen Sink — GitHub" href="https://github.com/wailsapp/wails/tree/master/v3/examples/mobile" description="Bindings · Eventos · Caixas de diálogo · Resposta tátil · Geolocalização · Biometria · Notificações · Armazenamento seguro · e muito mais — tudo a partir de um único main.go"}
Ele demonstra todas as principais APIs para dispositivos móveis em 7 abas — e também é executado no desktop. As abas **Dispositivos móveis** e **Hardware** ficam ocultas no desktop por meio de uma verificação de plataforma no frontend; quando compilado para desktop, o código Go não registra nenhum handler para eventos móveis `common:*`. Esse é o padrão recomendado para distribuir uma única base de código em todas as plataformas.

| Aba | Plataformas | O que demonstra |
| --- | --- | --- |
| **Bindings** | todas | Chamadas de serviços JS → Go que retornam valores, structs e erros |
| **Eventos** | todas | Relógio Go → JS, ping/pong JS → Go → JS e eventos do sistema operacional (bateria, rede e tema) |
| **Caixas de diálogo** | todas | Caixas de diálogo de mensagem nativas em cada plataforma |
| **Sistema** | todas | Área de transferência, métricas da tela e informações do dispositivo |
| **Dispositivos móveis** | iOS + Android | Menu de compartilhamento, manutenção da tela ativa, lanterna, brilho, biometria, notificações locais e armazenamento seguro |
| **Hardware** | iOS + Android | Resposta tátil, geolocalização, acelerômetro, proximidade e conversão de texto em fala |
| **Nativo** | iOS + Android | iOS: resposta tátil + controles para ativar ou desativar opções da WKWebView · Android: vibração + toast |

Para executá-lo por conta própria:

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator (macOS + Xcode required)
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

## Como funciona

O mesmo modelo de aplicativo se aplica a todas as plataformas:

1. **Backend Go** — seus serviços, handlers de eventos e lógica do aplicativo são compilados sem alterações para `GOOS=ios` e `GOOS=android`.
2. **Frontend** — exatamente o mesmo HTML/JS/CSS. O pacote `@wailsio/runtime` funciona de forma idêntica; bindings de serviços, eventos, caixas de diálogo e a área de transferência são todos encaminhados pelo mesmo transporte dentro do processo.
3. **Host da WebView** — no iOS, uma `WKWebView` dentro de um `UIViewController`; no Android, uma `WebView` dentro de uma `Activity`. O Wails configura automaticamente a ponte de mensagens.
4. **Fornecimento de recursos dentro do processo** — os recursos são servidos diretamente da memória do Go, e não por um servidor localhost. Sem portas abertas, sem loopback e sem latência adicional.

O comportamento específico de cada plataforma fica em arquivos protegidos por `//go:build ios` ou `//go:build android`, mantendo limpo o código compartilhado.

## Visão geral dos pré-requisitos

| Requisito | iOS | Android |
| --- | --- | --- |
| Sistema operacional | Somente macOS | macOS, Linux, Windows |
| Conjunto de ferramentas | Instalação completa do Xcode (não apenas as ferramentas de linha de comando) | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| Verifique com | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
Execute `wails3 doctor` depois de configurar seu conjunto de ferramentas — o comando mostra exatamente o que foi encontrado e o que está faltando para cada plataforma.

@end

## O que é compatível

As duas plataformas compartilham o mesmo conjunto principal de recursos:

| Recurso | iOS | Android |
| --- | --- | --- |
| Vinculações de serviços (JS → Go) | ✅ | ✅ |
| Eventos (em ambas as direções) | ✅ | ✅ |
| Caixas de diálogo de mensagem | ✅ UIAlertController | ✅ AlertDialog |
| Caixas de diálogo para abrir arquivos | ✅ UIDocumentPicker | ✅ Storage Access Framework |
| Caixas de diálogo para salvar arquivos | ❌ grave na sandbox em vez disso | ❌ grave na sandbox em vez disso |
| Área de transferência | ✅ UIPasteboard | ✅ ClipboardManager |
| Métricas de tela/área segura | ✅ | ✅ |
| Eventos do ciclo de vida | ✅ `events.IOS.*` | ✅ `events.Android.*` |
| Resposta tátil | ✅ `IOS.Haptics.*` | ✅ `Android.Haptics.Vibrate` |
| Informações do dispositivo | ✅ `IOS.Device.Info()` | ✅ `Android.Device.Info()` |
| Abas nativas (iOS) | ✅ UITabBar | — |
| Mensagens toast (Android) | — | ✅ `Android.Toast.Show` |
| Várias janelas | ❌ apenas a primeira janela | ❌ apenas a primeira janela |
| Geometria da janela/menus/bandeja do sistema | operações nulas intencionais | operações nulas intencionais |

## Regras de tags de compilação

Duas regras importantes para escrever código condicional por plataforma:

- **`ios` implica `darwin`** — um arquivo com a tag `//go:build darwin` também será compilado para iOS. Para direcionar apenas ao macOS, use `//go:build darwin && !ios`.
- **`android` implica `linux`** — um arquivo com a tag `//go:build linux` também será compilado para Android. Para direcionar apenas ao Linux para desktop, use `//go:build linux && !android`.

Em tempo de execução, `runtime.GOOS` retorna `"ios"` e `"android"`, respectivamente.

## Detecção da plataforma em tempo de execução

As tags de compilação destinam-se a código que só pode ser *compilado* em determinada plataforma. Para ramificações comuns em código compartilhado, use `application.System` — ele está disponível em todas as compilações (sem necessidade de tags de compilação), portanto o mesmo arquivo funciona em qualquer plataforma:

```go
import "github.com/wailsapp/wails/v3/pkg/application"

if application.System.IsMobile() {
    // iOS or Android
} else if application.System.IsDesktop() {
    // macOS, Windows or Linux
}

// Or test a single target directly:
if application.System.IsPlatform(application.PlatformIOS) {
    // iOS only
}
```

Disponíveis: `IsMobile()`, `IsDesktop()`, `IsServer()` (a tag de compilação `server`) e `IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)`.

O frontend tem as funções auxiliares correspondentes em `@wailsio/runtime`:

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## Próximas etapas

@cards{cols="2"}
🚀 Seu primeiro aplicativo móvel
Pegue um aplicativo Wails para desktop e execute-o no iOS Simulator ou no Android Emulator em poucos minutos.

[Comece agora →](/guides/mobile/first-mobile-app/)

---
Guia do iOS
Configuração completa da cadeia de ferramentas do iOS, simulador, compilações para dispositivos, assinatura, configuração e referência da API.

[Guia do iOS →](/guides/mobile/ios/)

---
Guia do Android
Configuração completa do SDK/NDK do Android, emulador, assinatura de APKs, empacotamento para a Play Store e referência da API.

[Guia do Android →](/guides/mobile/android/)

@end
