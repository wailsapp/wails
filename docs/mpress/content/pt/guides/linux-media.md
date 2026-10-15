---
title: "Reproduzir áudio e vídeo locais"
description: "Reproduza mídia incluída no aplicativo Linux usando URLs de blob com tamanho limitado e libere-a ao terminar."
slug: "guides/linux-media"
sourcePath: "guides/linux-media.md"
---

Use `Media.SetSource` para reproduzir clipes locais curtos em um aplicativo Wails v3. No Linux, WebKitGTK delega a reprodução de mídia ao GStreamer, que não pode carregar URLs `wails://` diretamente. A função recebe o clipe por um fluxo Wails e atribui uma URL de blob ao reprodutor.

As transferências no desktop usam o transporte de recursos existente do Wails sem abrir um socket de escuta. A API funciona com elementos de áudio e vídeo. Mídias HTTP/HTTPS comuns podem usar diretamente o `src` do reprodutor.

## Registrar os arquivos de mídia

Exponha um sistema de arquivos contendo os clipes que o frontend pode reproduzir e registre o manipulador de mídia em um fluxo nomeado:

```go
import (
    "embed"
    "io/fs"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/services/media"
)

//go:embed clips
var clips embed.FS

func registerMedia(app *application.App) {
    mediaFiles, err := fs.Sub(clips, "clips")
    if err != nil {
        log.Fatal(err)
    }
    handler, err := media.NewHandler(mediaFiles, 32 << 20) // 32 MiB per file
    if err != nil {
        log.Fatal(err)
    }
    app.HandleStream("media", handler)
}
```

Chame `registerMedia(app)` depois de criar o aplicativo e antes de chamar `app.Run()`.

Para arquivos em disco, use `os.OpenRoot(directory)` e passe `root.FS()` para `media.NewHandler`. Mantenha a raiz aberta até `app.Run()` retornar e depois feche-a. Escolha um diretório contendo apenas arquivos que o frontend pode ler; `os.Root` restringe o acesso mesmo se um link simbólico apontar para fora desse diretório.

## Carregar um clipe

Crie um reprodutor:

```html
<video id="player" controls></video>
```

Para um frontend que usa o runtime npm:

```javascript
import { Media } from '@wailsio/runtime';

const player = document.getElementById('player');

try {
    await Media.SetSource(player, 'media', 'welcome.mp4');
} catch (error) {
    if (error.name !== 'AbortError') {
        console.error('Could not load the clip:', error);
    }
}
```

Para um aplicativo que usa o runtime incluído, altere a importação para:

```javascript
import { Media } from '/wails/runtime.js';
```

O segundo argumento é o nome do fluxo registrado. O terceiro é um caminho de arquivo separado por barras, relativo à raiz do sistema de arquivos, como `welcome.mp4` ou `tutorials/intro.mp4`. Não é uma URL nem um caminho do sistema operacional.

A promessa é resolvida quando a fonte é atribuída. Em seguida, o reprodutor a decodifica. Trate o evento `error` do reprodutor para detectar codecs não suportados e chame `player.play()` a partir de uma interação do usuário se precisar iniciar a reprodução por conta própria.

## Substituir ou liberar um clipe

Chame `Media.SetSource` novamente para trocar o clipe. Isso cancela um carregamento anterior pendente para esse reprodutor, impedindo que uma resposta lenta sobrescreva a seleção mais recente. O clipe anterior permanece disponível até que seu substituto seja carregado com sucesso. Sua URL de blob é então revogada.

Chame `Media.ClearSource(player)` ao fechar um reprodutor ou desmontar seu componente:

```javascript
Media.ClearSource(player);
```

Isso cancela um carregamento pendente, redefine o reprodutor e libera a URL de blob. Chame antes de remover o elemento do documento. Um componente de framework deve chamar a função no hook de desmontagem ou descarte. Remover apenas o elemento não libera a URL de blob. Use essas funções de forma consistente para a fonte do reprodutor.

Para a marcação `<video><source ...></video>`, passe o nome do arquivo escolhido para `Media.SetSource(video, "media", name)`. A função define o `src` do reprodutor pai, que tem precedência sobre os filhos `<source>`. Limpar esse atributo permite que o navegador considere esses filhos novamente; use um reprodutor vazio ao gerenciar todas as fontes por esta API.

Objetos de áudio fora do documento também funcionam:

```javascript
const sound = new Audio();
await Media.SetSource(sound, 'media', 'notification.mp3');
sound.addEventListener('ended', () => Media.ClearSource(sound), { once: true });
// Call sound.play() from an appropriate user interaction.
// Also clear it if playback is cancelled or the owning component is disposed.
```

## Limitar downloads e cancelar o carregamento

O limite padrão é de **32 MiB por fonte**. Você pode escolher um limite inteiro positivo menor ou maior, em bytes, dentro do limite configurado em Go:

```javascript
const controller = new AbortController();
const loading = Media.SetSource(player, 'media', 'welcome.mp4', {
    maxBytes: 8 * 1024 * 1024,
    signal: controller.signal,
});

// Call controller.abort() to cancel this load.
await loading;
```

Um arquivo muito grande causa rejeição com `RangeError`. O manipulador Go verifica o tamanho antes de ler o conteúdo e transfere no máximo o menor valor entre seu limite configurado e o limite do frontend. O frontend também verifica o tamanho recebido e rejeita transferências incompletas. O cancelamento causa rejeição com `AbortError` ou com o motivo passado para `AbortController.abort(reason)`.

Os arquivos são enviados em quadros de 64 KiB. O transporte de fluxos do Wails no desktop limita as respostas de consulta a 1 MiB, de modo que o armazenamento de respostas completas em buffer pelo WebView2 não armazena todo o arquivo de mídia em uma única resposta. As filas de fluxos têm seu próprio controle limitado de contrapressão. Esses limites de transporte não alteram o comportamento de blob completo da função do reprodutor.

**O arquivo inteiro é baixado antes da reprodução.** O limite de bytes é por arquivo, não um limite de memória total do aplicativo. Vários reprodutores, o clipe anterior durante a substituição, a construção do blob e a mídia decodificada podem usar memória adicional. A função transfere ao ser chamada, independentemente da configuração `preload` do reprodutor. Chame quando o usuário escolher carregar um clipe.

Para arquivos locais grandes, essa função não é uma solução de streaming. Aumentar o limite também aumenta o uso de memória. Mídias já hospedadas por HTTP/HTTPS devem usar o carregamento nativo para que o navegador possa transmitir e buscar posições com requisições de intervalo.

## Solucionar problemas de reprodução no Linux

- Se a reprodução local direta informar **No URI handler implemented for "wails"**, carregue o clipe com `Media.SetSource`. Isso afeta tanto a pilha GTK4 padrão quanto a pilha legada `-tags gtk3`.
- Se o carregamento funcionar, mas a decodificação falhar, verifique os codecs GStreamer instalados no sistema de destino. MP4 geralmente exige suporte a vídeo H.264 e áudio AAC; MP3 precisa de um decodificador MP3. Teste os formatos distribuídos nas distribuições suportadas.
- Se a transferência falhar, verifique o nome do fluxo registrado, o nome relativo do arquivo e as permissões do sistema de arquivos. Alterar um arquivo durante a transferência pode torná-la incompleta; tente novamente após a conclusão da gravação.
- Se o aplicativo definir uma Content Security Policy, permita a origem dos recursos Wails em `connect-src` e `blob:` em `media-src`. Para uma política exclusivamente local, essas diretivas podem ser `connect-src 'self'; media-src 'self' blob:`. Preserve as outras diretivas.
- Se um arquivo ultrapassar o limite, escolha um clipe mais curto ou menor ou defina um limite explícito adequado ao orçamento de memória do aplicativo.

Execute o [exemplo audio-video](https://github.com/wailsapp/wails/tree/master/v3/examples/audio-video) com `go run .` para testar as amostras MP3 e MP4 incluídas na sua máquina.
