---
title: "Android"
description: "Compile e execute aplicativos Wails no Android — configuração da cadeia de ferramentas, emulador, assinatura de APK, empacotamento para a Play Store e referência da API"
slug: "guides/mobile/android"
sourcePath: "guides/mobile/android.md"
---

@note{type="caution" title="Recurso experimental"}
O suporte ao Android é experimental e pode mudar em versões futuras.

@end

@note{type="tip"}
Está começando no desenvolvimento para dispositivos móveis com Wails? Comece por [Seu primeiro aplicativo móvel →](/guides/mobile/first-mobile-app/) para ver um passo a passo e depois volte aqui para consultar a referência completa.

@end

Os aplicativos Wails v3 são executados no Android como aplicativos nativos: uma `WebView` do Android renderiza o frontend, os ativos são servidos **no processo** por meio de uma `WebViewAssetLoader`, com suporte do servidor de ativos em Go (sem servidor localhost nem portas abertas), e a `@wailsio/runtime` padrão funciona sem alterações — vinculações de serviços, eventos, caixas de diálogo e a área de transferência são encaminhados pelo processador de mensagens em Go.

O mesmo `main.go` é compilado para desktop e Android. O código Go é compilado como uma biblioteca compartilhada C (`libwails.so`, `GOOS=android` + a cadeia de ferramentas do NDK) e carregado por um pequeno host Java. O comportamento específico do Android fica em arquivos Go próprios da plataforma, protegidos por `//go:build android`.

## Requisitos

- O **SDK do Android** com platform-tools, uma plataforma do SDK (API 35), build-tools e o **NDK** (26.3.x) — `wails3 doctor` mostra o que foi encontrado
- Um **JDK** (por exemplo, OpenJDK 21) para o Gradle; defina `JAVA_HOME` se `java` não estiver no seu `PATH`
- Go 1.26+ e npm
- `ANDROID_HOME` (ou `ANDROID_SDK_ROOT`) apontando para o SDK

Instale os componentes do SDK com as ferramentas de linha de comando:

```bash
sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" \
           "ndk;26.3.11579264" "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails \
           --package "system-images;android-35;google_apis;arm64-v8a" \
           --device pixel_7
```

## Execução no emulador

No diretório do seu projeto:

```bash
wails3 task android:run
```

Isso inicializa um emulador caso nenhum esteja em execução, gera as vinculações, compila o frontend, compila seu código Go como `libwails.so` para a ABI do emulador, monta um APK de depuração com o Gradle e, em seguida, o instala e inicia.

Comandos complementares úteis:

```bash
wails3 task android:logs    # stream the app's logcat output
```

Em compilações de depuração, a WebView pode ser inspecionada pelo Chrome em `chrome://inspect`.

## Empacotamento

```bash
wails3 task android:package             # production release APK
wails3 task android:deploy-emulator     # install + launch it
wails3 task android:bundle              # production release AAB (Android App Bundle)
wails3 task android:bundle:fat          # release AAB containing all ABIs
wails3 task android:run:device          # debug install + launch on a physical device
wails3 task android:deploy-device       # install + launch on a physical device
DEVICE_ID=<serial> wails3 task android:run:device
DEVICE_ID=<serial> wails3 task android:deploy-device
```

As compilações de produção usam `-tags production,android`, têm os símbolos removidos e excluem da compilação os diagnósticos internos do framework. `wails3 task android:package:fat` compila `arm64-v8a` e `x86_64` em um único APK.

O Google Play exige o formato Android App Bundle (`.aab`) para o envio de novos aplicativos, que devem ter como destino o Android 15 (API 35) ou posterior; o modelo do projeto define `compileSdk` e `targetSdk` como 35 em `build/android/app/build.gradle`. `wails3 task android:bundle:fat` gera `bin/<AppName>.aab` com as duas ABIs incluídas; o Google Play gera a partir dele APKs otimizados para cada dispositivo, portanto o pacote universal é o artefato correto para enviar à loja. Os APKs continuam sendo a opção mais rápida para testes locais e no emulador, pois um `.aab` não pode ser instalado diretamente com `adb`.

`android:run` e `android:deploy-emulator` são tarefas voltadas ao emulador. Para um dispositivo Android físico, use `android:run:device` para gerar um APK de depuração ou `android:deploy-device` para gerar um APK de lançamento. Ambas compilam para `arm64`, selecionam a primeira entrada conectada que não seja um emulador em `adb devices`, instalam o APK e iniciam `com.wails.app.MainActivity`. Passe `DEVICE_ID=<serial>` para selecionar um dispositivo específico.

## Assinatura e compilações de lançamento

Sem um keystore, as compilações de lançamento são assinadas com o keystore de **depuração** do Android para que possam ser instaladas para testes. Para assinar com seu próprio keystore, defina:

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=... \
ANDROID_KEY_ALIAS=... \
ANDROID_KEY_PASSWORD=... \
  wails3 task android:package
```

As mesmas variáveis assinam App Bundles: execute `wails3 task android:bundle:fat` com elas definidas para gerar um `.aab` pronto para a Play Store. Sem elas, o pacote é assinado com o keystore de depuração e o Google Play o rejeitará; por isso, a tarefa exibe um aviso.

@note{type="tip"}
Com o [Play App Signing](https://support.google.com/googleplay/android-developer/answer/9842756), o keystore usado para assinar localmente contém sua **chave de upload**: o Google a utiliza para verificar seu upload e depois assina novamente o aplicativo com a chave de assinatura do aplicativo que ele gerencia. Observe também que o Google Play exige um `versionCode` maior a cada upload; incremente-o em `build/android/app/build.gradle`.

@end

## Configuração

O frontend controla os recursos do Android em tempo de execução por meio do objeto de runtime `Android`: `Android.Haptics.Vibrate(durationMs)`, `Android.Device.Info()`, `Android.Toast.Show(message)`. O nome do pacote é controlado por `APP_ID` nas tarefas de compilação.

## O que funciona e o que não funciona

| Área | Status |
| --- | --- |
| WebView + ativos no processo (`WebViewAssetLoader`) | ✅ |
| Vinculações de serviços e eventos (nos dois sentidos) | ✅ |
| Caixas de diálogo de mensagens | ✅ AlertDialog com callbacks de botões |
| Caixas de diálogo para abrir arquivo/arquivos | ✅ Storage Access Framework (arquivos importados como cópias no cache) |
| Caixas de diálogo para abrir diretório/salvar arquivo | ❌ Retornam um erro — grave no sandbox do aplicativo |
| Área de transferência | ✅ ClipboardManager |
| API de telas | ✅ WindowMetrics, incluindo a área de trabalho delimitada pelas barras do sistema |
| Eventos de ciclo de vida (`events.Android.*`) | ✅ |
| Resposta tátil, informações do dispositivo e toast | ✅ API de runtime `Android.*` |
| Compilações para emulador e dispositivo físico | ✅ `android:run`, `android:run:device`, `android:deploy-emulator`, `android:deploy-device` |
| Geometria da janela, menus e bandeja do sistema | Operações inativas intencionalmente |
| Várias janelas | Somente a primeira janela é exibida |

## Observações sobre portabilidade

- O código para desktop é compilado sem alterações com `GOOS=android`; as chamadas de geometria, menu e bandeja tornam-se operações inativas porque os aplicativos Android são executados em tela cheia.
- `android` **implica a tag de compilação `linux`** (o Android usa um kernel Linux): arquivos exclusivos do Linux para desktop precisam de `//go:build linux && !android` e, em tempo de execução, `runtime.GOOS` é `"android"`.
- Substitua as caixas de diálogo para salvar arquivos e escolher diretórios por gravações no sandbox do aplicativo, acompanhadas de um fluxo de compartilhamento por intent. As caixas de diálogo para abrir arquivos funcionam e importam os documentos escolhidos como cópias no diretório de cache, fornecendo caminhos reais do sistema de arquivos.
- Um aplicativo real sempre é compilado com `CGO_ENABLED=1` e o NDK; o caminho sem cgo existe apenas para que ferramentas como `wails3 generate bindings` possam carregar o pacote.
- Crie um frontend responsivo; a área de trabalho de `Screens` exclui as barras de status e de navegação.
