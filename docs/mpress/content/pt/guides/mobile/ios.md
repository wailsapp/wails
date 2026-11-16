---
title: "iOS"
description: "Compile e execute aplicativos Wails no iOS — configuração inicial, simulador, builds para dispositivos, configurações e recursos nativos"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="Recurso experimental"}
O suporte ao iOS é experimental e pode mudar em versões futuras.

@end

@note{type="tip"}
Está começando a desenvolver para dispositivos móveis com o Wails? Comece por [Seu primeiro aplicativo móvel →](/guides/mobile/first-mobile-app/) para ver um passo a passo e depois volte aqui para consultar a referência completa.

@end

Os aplicativos Wails v3 são executados no iOS como aplicativos totalmente nativos — e o melhor é que funcionam *exatamente* como a versão para desktop. O mesmo backend em Go, o mesmo frontend e os mesmos `@wailsio/runtime`: vinculações de serviços, eventos, caixas de diálogo e área de transferência se comportam de forma idêntica, com **zero** adaptações específicas para dispositivos móveis. Não há uma base de código móvel separada, uma camada de portabilidade nem uma API especial para aprender — seu aplicativo Wails existente simplesmente é executado no iOS. A portabilidade é realmente transparente: leve o aplicativo como está e publique-o.

O mesmo `main.go` é compilado para desktop e iOS; os ajustes específicos do iOS são definidos por meio de `application.Options.IOS`.

## Requisitos

- macOS com o **Xcode completo** instalado (somente as ferramentas de linha de comando não são suficientes) — `wails3 doctor` mostra os SDKs do iOS que ele consegue encontrar
- Go 1.26+ e npm

## Simulador

No diretório do projeto:

```bash
wails3 task ios:run
```

Esse comando compila o aplicativo, inicializa um simulador caso nenhum já esteja em execução e abre o aplicativo.

Comandos complementares úteis:

```bash
wails3 task ios:logs:dev    # stream the app's logs from the simulator
wails3 task ios:xcode       # open the generated Xcode project
```

Em builds de depuração, a WebView pode ser inspecionada pelo menu Desenvolvedor do Safari.

## Empacotamento

```bash
wails3 task ios:package             # production .app for the simulator
wails3 task ios:deploy-simulator    # install + launch it
```

Esses são builds de produção otimizados e sem símbolos desnecessários.

## Builds para dispositivos

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device` compila para um dispositivo físico. As autorizações vêm de `build/ios/entitlements.plist` e se aplicam somente a builds para dispositivos — adicione as chaves de recursos necessárias ao aplicativo.

@note{type="tip"}
Para gerenciar automaticamente a assinatura, o provisionamento e os arquivos da App Store, abra o projeto Xcode gerado com `wails3 task ios:xcode` e faça o build pelo Xcode.

@end

## Configuração

`build/config.yml`:

```yaml
ios:
  bundleID: com.example.myapp
  displayName: My App
  version: 1.0.0
  minIOSVersion: "15.0"
```

As opções de inicialização (`application.Options.IOS`) incluem `DisableScroll`, `DisableBounce`, `DisableScrollIndicators`, `DisableInputAccessoryView`, `EnableBackForwardNavigationGestures`, `DisableLinkPreview`, `EnableInlineMediaPlayback`, `EnableAutoplayWithoutUserAction`, `DisableInspectable`, `UserAgent`, `ApplicationNameForUserAgent`, `BackgroundColour` e abas inferiores nativas por meio de `EnableNativeTabs` + `NativeTabsItems`.

## Recursos nativos

Os recursos específicos do iOS estão disponíveis por meio de `application.IOS`, chamado pelo Go dentro de um arquivo `//go:build ios` para que o código compartilhado permaneça independente de plataforma. O Android oferece o mesmo conjunto por meio de `application.Android`.

As ações pontuais retornam imediatamente:

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

As funções auxiliares de consulta retornam os resultados como JSON — `SafeAreaJSON()`, `AppInfoJSON()`, `PowerJSON()`, `NetworkJSON()`, `StorageJSON()`, `GetOrientation()`, `GetBrightness()`. `StoragePath()` retorna o caminho absoluto do diretório Application Support do aplicativo — um bom local para bancos de dados e outros arquivos persistentes (o equivalente no iOS ao `getFilesDir()` do Android). O diretório é criado no primeiro acesso; `StoragePath()` retorna uma string vazia se não for possível criá-lo. Portanto, verifique se o valor é `""` antes de usá-lo.

### Eventos

Tudo que for concluído posteriormente — uma solicitação de permissão, um fluxo de dados de sensor ou uma captura da câmera — entrega o resultado como um **evento**, e não como valor de retorno. Você pode escutá-lo no Go ou no frontend. Os nomes recebem o prefixo `common:` para recursos compartilhados com o Android e `ios:` para recursos exclusivos do iOS.

```go
// Go
app.Event.On("common:location", func(e *application.CustomEvent) {
    // e.Data -> {"lat":..,"lng":..,"accuracy":..} or {"error":..}
})
```

```js
// frontend
import { Events } from "@wailsio/runtime";
Events.On("common:notification", (e) => { /* {ok, scheduled, presented, tapped, error} */ });
```

| Evento | Acionado por | Conteúdo |
| --- | --- | --- |
| `common:biometric` | `BiometricAuthenticate(reason)` | `{ok, error}` |
| `common:location` | `GetLocation()` | `{lat, lng, accuracy}` / `{error}` |
| `common:motion` | `SetMotion(true)` | `{x, y, z}` |
| `common:proximity` | `SetProximity(true)` | `{near}` |
| `common:keyboard` | `SetKeyboardWatch(true)` | `{visible, height}` |
| `common:torch` | `SetTorch(bool)` | `{on, available}` |
| `common:notification` | `PostNotification(json)` | `{ok, scheduled, presented, tapped, error}` |
| `common:capture` | `CapturePhoto()` / `CaptureVideo()` | `{type, path, size, thumb}` |
| `common:screenCapture` | `SetScreenProtect(true)` | `{screenshot, recording}` |
| `ios:backgroundTask` | `BeginBackgroundTask(seconds)` | `{message, granted}` |

O exemplo completo em `v3/examples/mobile` integra todos os recursos acima de ponta a ponta.

## Controles da WebView

Alguns comportamentos da WebView também podem ser alterados em tempo de execução pelo Go:

```go
application.IOS.SetScrollEnabled(false)
application.IOS.SetBounceEnabled(false)
application.IOS.SetScrollIndicatorsEnabled(false)
application.IOS.SetBackForwardGesturesEnabled(true)
application.IOS.SetLinkPreviewEnabled(false)
application.IOS.SetInspectableEnabled(true)
application.IOS.SetCustomUserAgent("MyApp/1.0")
```

O `@wailsio/runtime` incluído também expõe um pequeno namespace do iOS no **frontend**:

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

As seleções de abas inferiores nativas chegam como um evento `nativeTabSelected` em `window`.

## Status do suporte

| Área | Status |
| --- | --- |
| Renderização do frontend e recursos | ✅ |
| Vinculações de serviços e eventos (em ambas as direções) | ✅ |
| Caixas de diálogo de mensagens | ✅ |
| Caixas de diálogo para abrir arquivo, arquivos ou diretório | ✅ Importados como cópias no sandbox |
| Caixas de diálogo para salvar arquivo | ❌ Grave dentro do sandbox do aplicativo |
| Área de transferência | ✅ |
| API de telas | ✅ Inclui a área de trabalho da área segura |
| Eventos do ciclo de vida | ✅ |
| Geometria da janela, menus e bandeja do sistema | Sem efeito no iOS |
| Várias janelas | Somente a primeira janela é exibida |

## Observações sobre portabilidade

- O código para desktop é compilado para iOS sem alterações — as chamadas de janela, menu e bandeja do sistema simplesmente não fazem nada.
- Substitua as caixas de diálogo para salvar arquivos pela gravação no sandbox do aplicativo seguida de um compartilhamento.
- Crie o frontend com um design responsivo; as áreas seguras são tratadas automaticamente.
