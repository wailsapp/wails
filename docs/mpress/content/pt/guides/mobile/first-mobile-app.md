---
title: "Seu primeiro aplicativo móvel"
description: "Execute seu aplicativo Wails no Simulador do iOS ou no Emulador do Android em poucos minutos"
slug: "guides/mobile/first-mobile-app"
sourcePath: "guides/mobile/first-mobile-app.md"
---

Este guia parte de um aplicativo Wails padrão para desktop e o executa no Simulador do iOS ou no Emulador do Android. **Você não precisa alterar seu código Go.** O mesmo `main.go` é compilado para todos os destinos.

**Tempo para concluir:** 15–30 minutos (a maior parte desse tempo é gasta na instalação da cadeia de ferramentas na primeira execução)

## Comece com um projeto para desktop

Se você ainda não tiver um, crie um novo projeto:

```bash
wails3 init -n mymobileapp
cd mymobileapp
```

Primeiro, confirme se o aplicativo para desktop funciona:

```bash
wails3 dev
```

Depois que ele abrir, encerre-o e prossiga. Tudo o que é executado no desktop também é executado em dispositivos móveis — neste guia, você não precisará alterar `main.go` nem qualquer código Go.

---

## Escolha sua plataforma

@tabs{sync-key="mobile-platform"}
[Simulador do iOS]
### Requisitos

- **macOS** (as compilações para iOS só podem ser feitas no macOS)
- **Xcode completo** — não apenas as ferramentas de linha de comando. Instale-o pela App Store e depois execute:
  ```bash
  sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
  sudo xcodebuild -license accept
  ```


- **Go 1.25+** e **npm** (já instalados se você executou `wails3 init`)

Execute `wails3 doctor` para verificar — esse comando lista os SDKs do iOS que consegue encontrar.

### Execute no Simulador

@steps
### Inicie o aplicativo
```bash
wails3 task ios:run
```

É só isso — esse comando compila seu aplicativo, inicializa um simulador se nenhum já estiver em execução e abre o aplicativo.

@note{type="tip"}
A primeira execução leva alguns minutos (o framework Wails para iOS está sendo compilado e armazenado em cache). Todas as execuções seguintes são muito mais rápidas.

@end

Quando for iniciado, seu aplicativo para desktop sem modificações estará em execução no Simulador do iOS — o mesmo `main.go` e o mesmo frontend:

![Um aplicativo Wails padrão em execução no Simulador do iOS](/assets/ios-simulator-first-app.png)

### Acompanhe os logs em tempo real
Em outro terminal:

```bash
wails3 task ios:logs:dev
```

Esse comando acompanha o log do simulador em tempo real, filtrado para seu aplicativo. A saída de `fmt.Println` e `log.Println` aparece aqui.

### Inspecione a WebView
No Safari: **Desenvolvimento → Simulador → seu aplicativo**. O Web Inspector completo funciona — console, depurador, painel de rede, tudo.

### Faça uma alteração
Edite qualquer arquivo do frontend (`frontend/src/main.js`, `index.html` etc.) e execute `wails3 task ios:run` novamente. O Wails recompila o frontend e reinicia o aplicativo.

Para alterações em Go, execute também `wails3 task ios:run` novamente. A recompilação do Go é incremental, portanto apenas os pacotes alterados são recompilados.

@end

### Abra no Xcode (opcional)

```bash
wails3 task ios:xcode
```

Esse comando abre `build/ios/` no Xcode. Você pode usar o Xcode para implantação em dispositivos, análise avançada de desempenho ou gerenciamento de perfis de provisionamento. O Wails gera novamente o projeto do Xcode a cada compilação, portanto não modifique diretamente os arquivos gerados.

[Emulador do Android]
### Requisitos

Você precisa do **SDK do Android**, do **NDK** e de um **JDK**. A maneira mais fácil é usar o Android Studio ou as ferramentas de linha de comando:

@steps
### Instale as ferramentas de linha de comando do Android
Baixe-as em [developer.android.com/studio#command-line-tools-only](https://developer.android.com/studio#command-line-tools-only) e descompacte-as em `~/android-sdk/cmdline-tools/latest/`.

### Instale os componentes do SDK
```bash
sdkmanager "platform-tools" \
           "platforms;android-35" \
           "build-tools;35.0.0" \
           "ndk;26.3.11579264" \
           "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
```

### Crie um emulador
```bash
avdmanager create avd \
  --name wails \
  --package "system-images;android-35;google_apis;arm64-v8a" \
  --device pixel_7
```

### Defina as variáveis de ambiente
Adicione a `~/.zshrc` ou `~/.bashrc`:

```bash
export ANDROID_HOME=~/android-sdk
export ANDROID_SDK_ROOT=~/android-sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/cmdline-tools/latest/bin
```

Recarregue: `source ~/.zshrc`

### Instale um JDK
```bash
# macOS
brew install openjdk@21
export JAVA_HOME=$(brew --prefix openjdk@21)

# Ubuntu/Debian
sudo apt install openjdk-21-jdk
export JAVA_HOME=/usr/lib/jvm/java-21-openjdk-amd64

# Windows (scoop)
scoop install openjdk21
```

@end

Execute `wails3 doctor` para confirmar que todos os componentes foram encontrados.

### Execute no Emulador

@steps
### Inicie o aplicativo
```bash
wails3 task android:run
```

Na primeira execução, esse comando:

- Inicializa o emulador se nenhum estiver em execução
- Gera os bindings e compila o frontend
- Compila seu código Go para `libwails.so` usando o compilador cruzado do NDK
- Monta um APK de depuração com o Gradle
- Instala e inicia o APK no emulador

@note{type="tip"}
A primeira compilação baixa o Gradle e compila a cadeia de ferramentas do NDK — espere de 5 a 10 minutos. As compilações seguintes são incrementais e levam menos de um minuto.

@end

### Acompanhe os logs em tempo real
Em outro terminal:

```bash
wails3 task android:logs
```

Esse comando executa `adb logcat` com um filtro para seu aplicativo. A saída de `fmt.Println` aparece aqui.

### Inspecione a WebView
Abra o Chrome e acesse `chrome://inspect`. A WebView do seu aplicativo aparece em **Destino remoto** — clique em **inspecionar** para abrir o DevTools.

### Faça uma alteração
Edite qualquer arquivo e execute `wails3 task android:run` novamente. Com a compilação incremental do Gradle, apenas o código alterado é recompilado.

@end

@end

---

## Entenda o que aconteceu

Seu `main.go` não mudou em nada. O Wails cuidou de tudo:

- **Sistema de compilação** — o `Taskfile.yml` do seu projeto contém as tarefas `ios:*` e `android:*`, que acionam o conjunto de ferramentas específico da plataforma.
- **Compilação cruzada do Go** — `GOOS=ios` ou `GOOS=android` com o `GOARCH` e o sysroot apropriados.
- **Host nativo** — um projeto Xcode (iOS) ou Gradle (Android) gerado que incorpora seu código Go compilado e hospeda a WebView.
- **Disponibilização de recursos** — seu `frontend/dist/` é incorporado ao binário Go e disponibilizado no próprio processo. Não é necessário um servidor localhost.

---

## Adapte seu aplicativo a dispositivos móveis

Seu aplicativo já funciona, mas parece um aplicativo para desktop na tela de um celular. Algumas pequenas alterações fazem uma grande diferença.

### CSS responsivo

As telas de dispositivos móveis são mais estreitas e usam padrões de entrada diferentes. Em `frontend/public/style.css` (ou equivalente):

```css
/* Prevent horizontal scrolling */
body {
  overflow-x: hidden;
}

/* Touch-friendly tap targets */
button {
  min-height: 44px;
  min-width: 44px;
}

/* Respect the iOS safe area (notch, home indicator) */
body {
  padding-top: env(safe-area-inset-top);
  padding-bottom: env(safe-area-inset-bottom);
  padding-left: env(safe-area-inset-left);
  padding-right: env(safe-area-inset-right);
}
```

### Detecte a plataforma em Go

Use tags de compilação para adicionar comportamentos específicos da plataforma sem sobrecarregar o código compartilhado.

Crie `mobile_ios.go` para o código exclusivo do iOS:

```go {title="mobile_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.IOSOptions {
    return application.IOSOptions{
        DisableBounce: true,
    }
}
```

Crie `mobile_android.go` para o código exclusivo do Android:

```go {title="mobile_android.go"}
//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.AndroidOptions {
    return application.AndroidOptions{}
}
```

Crie `mobile_desktop.go` como stub para que o código compartilhado também seja compilado no desktop:

```go {title="mobile_desktop.go"}
//go:build !ios && !android

package main

type mobileOptions struct{}

func platformOptions() mobileOptions { return mobileOptions{} }
```

### Detecte a plataforma em JavaScript e restrinja a interface exclusiva de dispositivos móveis

Os objetos de runtime `IOS.*` e `Android.*` só existem em suas respectivas plataformas. Chamá-los no desktop gera uma exceção. O padrão correto — usado pelo Kitchen Sink — é detectar a plataforma uma única vez e ocultar completamente os controles exclusivos de dispositivos móveis:

```javascript
// Detect platform from the bridge the host injects into the WebView
const platform = (() => {
  if (typeof window.wails?.platform === 'function') return window.wails.platform(); // Android
  if (window.webkit?.messageHandlers?.external) return 'ios';
  return 'desktop';
})();

const isIOS     = platform === 'ios';
const isAndroid = platform === 'android';
const isMobile  = isIOS || isAndroid;

// Hide any element marked as mobile-only
document.querySelectorAll('.mobile-only').forEach(el => {
  el.style.display = isMobile ? '' : 'none';
});
```

Depois, no seu HTML:

```html
<section class="mobile-only">
  <button id="btnHaptic">Haptic feedback</button>
</section>
```

Dessa forma, os botões exclusivos de dispositivos móveis nunca são renderizados no desktop, e você não precisa proteger cada chamada individual com uma verificação `if (isMobile)`.

No lado do Go, combine isso com um stub controlado por tag de compilação para que os manipuladores de eventos sejam registrados somente nas plataformas que precisam deles:

```go {title="native_desktop.go"}
//go:build !ios && !android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

// No-op on desktop — mobile tabs are hidden in the frontend so these
// events are never emitted.
func registerNativeFeatures(app *application.App) {}
```

```go {title="native_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func registerNativeFeatures(app *application.App) {
    app.Event.On("common:haptic", func(e *application.CustomEvent) {
        // only compiled and called on iOS
        application.IOS.Haptic("medium")
    })
    // ... other handlers
}
```

Esse é exatamente o padrão usado pelo [Kitchen Sink](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile) — consulte `native_features_stub.go`, `native_features_ios.go` e `native_features_android.go`.

@note{type="note" title="API de recursos nativos e nomenclatura de eventos"}
Vale a pena conhecer duas convenções:

- **Os recursos nativos no lado do Go usam gerenciadores de plataforma.** Chame-os por meio dos singletons `application.IOS.*` e `application.Android.*` — por exemplo, `application.IOS.Haptic("medium")` ou `application.Android.Share(payload)`. Cada gerenciador existe apenas em sua própria plataforma, portanto suas chamadas ficam em arquivos `//go:build ios` / `//go:build android`.
- **Os eventos são organizados em namespaces conforme seu alcance.** Tudo o que ambas as plataformas entendem usa o prefixo `common:*` (`common:haptic`, `common:location`, …); os eventos que apenas uma plataforma pode produzir ou tratar usam `ios:*` ou `android:*` (por exemplo, `ios:backgroundTask`, `android:foregroundService`). Como quase todos os recursos móveis são compartilhados, seu frontend mantém um único listener por evento em `common:*`.

@end

### Adicione feedback tátil (iOS)

```javascript
import { IOS } from '@wailsio/runtime';

async function onButtonTap() {
  if (isIOS) {
    await IOS.Haptics.Impact({ style: 'medium' });
  }
  // ... rest of your handler
}
```

### Adicione vibração (Android)

```javascript
import { Android } from '@wailsio/runtime';

async function onButtonTap() {
  if (isAndroid) {
    await Android.Haptics.Vibrate(50); // 50ms
  }
}
```

---

## Compile para produção

@tabs{sync-key="mobile-platform"}
[iOS]
**Build para simulador** (para testes no simulador; não é necessário assinar):

```bash
wails3 task ios:package
wails3 task ios:deploy-simulator
```

**Build para dispositivo** (requer uma identidade de assinatura e um perfil de provisionamento):

```bash
wails3 task ios:package \
  IOS_PLATFORM=device \
  CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
  PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device   # installs via xcrun devicectl
```

**IPA de distribuição** (para a App Store ou o TestFlight):

```bash
wails3 task ios:package:ipa IOS_PLATFORM=device \
  CODESIGN_IDENTITY="..." \
  PROVISIONING_PROFILE=path/to/distribution.mobileprovision
```

@note{type="tip"}
Para enviar arquivos ao App Store Connect, use `wails3 task ios:xcode` e deixe o Xcode gerenciar a assinatura e o arquivamento — ele cuida automaticamente da complexidade dos certificados, perfis e notarização.

@end

[Android]
**APK de depuração** (assinado com o keystore de depuração do Android; pode ser instalado diretamente):

```bash
wails3 task android:package
wails3 task android:deploy-emulator
```

**APK de lançamento** (assinado com seu próprio keystore):

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=yourpassword \
ANDROID_KEY_ALIAS=youralias \
ANDROID_KEY_PASSWORD=yourkeypassword \
  wails3 task android:package
```

**APK universal** (arm64 + x86_64 em um único arquivo):

```bash
wails3 task android:package:fat
```

@note{type="tip"}
Para enviar arquivos à Play Store, gere um `.aab` (Android App Bundle) em vez de um APK — abra `build/android/` no Android Studio e use **Build → Generate Signed Bundle / APK**.

@end

@end

---

## Solução de problemas

### `wails3 task ios:run` falha com "nenhum SDK do iOS encontrado"

A versão completa do Xcode deve estar instalada e selecionada:

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
xcode-select -p  # should print the Xcode path
```

#### `wails3 task android:run` falha com "SDK não encontrado"

Verifique se `ANDROID_HOME` está definida e exportada. Confirme com:

```bash
echo $ANDROID_HOME
ls $ANDROID_HOME/platform-tools/adb
```

#### O simulador não inicializa

Liste os simuladores disponíveis e inicialize um manualmente:

```bash
xcrun simctl list devices available
xcrun simctl boot "iPhone 16"
```

#### `chrome://inspect` não mostra nenhum destino

A WebView deve estar no modo de depuração (o padrão de `android:run`). Verifique se você está executando um build de depuração, não um de produção. Confirme também se `adb devices` mostra o emulador como conectado.

#### Os recuos da área segura não são aplicados

Verifique se o HTML inclui a metatag de viewport:

```html
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
```

---

## Explore o Kitchen Sink

Depois que seu primeiro aplicativo estiver em execução, o exemplo **Kitchen Sink** será a maneira mais rápida de descobrir outras possibilidades. Ele é um aplicativo Wails completo que, com uma única base de código, é executado no iOS, Android e desktop e inclui recursos de resposta tátil, geolocalização, biometria, notificações locais, armazenamento seguro e muito mais:

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

Consulte o código-fonte em [`v3/examples/mobile`](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile) — os arquivos `native_features_ios.go` e `native_features_android.go` são especialmente úteis como pontos de partida para copiar e colar ao implementar recursos específicos de cada plataforma.

## Próximos passos

@cards{cols="2"}
Guia do iOS
Referência completa: opções de configuração, abas nativas, opções de ativação do WKWebView, builds para dispositivos e assinatura.

[Guia do iOS →](/guides/mobile/ios/)

---
Guia do Android
Referência completa: configuração, toasts, empacotamento para a Play Store e detalhes do NDK.

[Guia do Android →](/guides/mobile/android/)

---
📖 Código-fonte do Kitchen Sink
Resposta tátil, geolocalização, biometria, notificações e armazenamento seguro — tudo em um único aplicativo executável.

[Ver no GitHub →](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile)

@end
