---
title: "API móvel"
description: "O gerenciador multiplataforma application.Mobile — um único ponto de entrada protegido por restrições de compilação para os recursos móveis nativos compartilhados pelo iOS e pelo Android"
slug: "guides/mobile/mobile-api"
sourcePath: "guides/mobile/mobile-api.md"
---

@note{type="caution" title="Recurso experimental"}
O suporte a dispositivos móveis é experimental e pode mudar em versões futuras.

@end

Os recursos móveis nativos são disponibilizados de duas formas:

- **Gerenciadores específicos de cada plataforma** — `application.IOS` (em arquivos `//go:build ios`) e `application.Android` (em arquivos `//go:build android`). Use-os para qualquer funcionalidade específica da plataforma. Consulte as referências do [iOS](/guides/mobile/ios/) e do [Android](/guides/mobile/android/) para ver a API completa de cada plataforma.
- **`application.Mobile`** — um único gerenciador, protegido por restrições de compilação, que abrange o subconjunto de recursos com comportamento idêntico nas duas plataformas. Use-o quando quiser um único fluxo de código que seja compilado e executado em todos os ambientes.

## `application.Mobile`

`application.Mobile` encaminha as chamadas para `IOS` no iOS, para `Android` no Android e para uma implementação stub que não realiza nenhuma operação no desktop. Como não tem restrições de compilação, você pode chamá-lo em código Go comum e independente de plataforma — sem precisar criar seus próprios arquivos `//go:build`:

```go
// Works in any file, on any target.
// On desktop this returns "" (no-op); on device it returns the real path.
dbDir := application.Mobile.StoragePath()
if dbDir == "" {
    // Off-device, or the directory could not be created — handle accordingly.
    return
}
db, _ := sql.Open("sqlite", filepath.Join(dbDir, "app.db"))
```

`StoragePath()` retorna o caminho absoluto do diretório privado de arquivos do aplicativo — `getFilesDir()` no Android e o diretório Application Support no iOS —, o local recomendado para bancos de dados e outros arquivos persistentes. Retorna uma string vazia no desktop e, no dispositivo, se o diretório não estiver disponível (no iOS, se não puder ser criado); portanto, verifique se o resultado é `""` antes de usá-lo.

@note{type="note"}
Fora do dispositivo (em compilações para desktop), todos os métodos de `Mobile` não realizam nenhuma operação e todas as consultas retornam o valor zero correspondente (`""` para strings). Isso permite que o código multiplataforma chame `application.Mobile.*` incondicionalmente. Quando também precisar de um caminho real no desktop, verifique a plataforma e use `os.UserConfigDir()` ou algo semelhante como alternativa.

@end

## Recursos

O gerenciador `Mobile` disponibiliza os recursos cujas assinaturas são idênticas no iOS e no Android:

| Recurso | API | Observações |
| --- | --- | --- |
| Folha de compartilhamento | `Mobile.Share(json)` | `{text, url}` |
| Abrir URL externamente | `Mobile.OpenURL(url)` | Navegador do sistema |
| Manter a tela ativa | `Mobile.SetKeepAwake(bool)` |  |
| Lanterna | `Mobile.SetTorch(bool)` | → `common:torch` |
| Margens da área segura | `Mobile.SafeAreaJSON()` | `{top,bottom,left,right}` |
| Informações do aplicativo | `Mobile.AppInfoJSON()` | `{name,version,build,bundleId}` |
| Bloqueio de orientação | `Mobile.SetOrientation(mode)` | `portrait` / `landscape` / `auto` |
| Barra de status | `Mobile.SetStatusBar(json)` | estilo + visibilidade |
| Informações de armazenamento | `Mobile.StorageJSON()` | `{free,total}` bytes |
| Caminho de armazenamento | `Mobile.StoragePath()` | Diretório privado de arquivos do aplicativo |
| Energia / bateria | `Mobile.PowerJSON()` | `{level,charging,lowPower}` |
| Status da rede | `Mobile.NetworkJSON()` | `{connected,type}` |
| Biometria | `Mobile.BiometricAuthenticate(reason)` | → `common:biometric` |
| Armazenamento seguro | `Mobile.SecureGet(key)` / `Mobile.SecureDelete(key)` | Keychain / `EncryptedSharedPreferences` |
| Geolocalização | `Mobile.GetLocation()` | consulta única → `common:location` |
| Resposta tátil | `Mobile.Haptic(type)` | impacto / notificação / seleção |
| Acelerômetro | `Mobile.SetMotion(bool)` | → `common:motion` |
| Proximidade | `Mobile.SetProximity(bool)` | → `common:proximity` |
| Conversão de texto em fala | `Mobile.Speak(text)` / `Mobile.StopSpeak()` |  |
| Margens do teclado | `Mobile.SetKeyboardWatch(bool)` | → `common:keyboard` |
| Captura de tela | `Mobile.SetScreenProtect(bool)` | → `common:screenCapture` |
| Câmera | `Mobile.CapturePhoto()` / `Mobile.CaptureVideo()` | → `common:capture` |

Os resultados assíncronos chegam como eventos `common:*`, exatamente como ocorre com os gerenciadores específicos de cada plataforma — consulte [Eventos](/guides/mobile/ios/#events) para ver os payloads.

## O que permanece específico de cada plataforma

Os recursos cujo formato difere entre iOS e Android **não** estão disponíveis em `Mobile`; chame-os por meio de `application.IOS` / `application.Android` em um arquivo com tags de compilação:

| Funcionalidade | iOS | Android |
| --- | --- | --- |
| Brilho (definir) | `IOS.SetBrightness(0.0-1.0)` | `Android.SetBrightness(0-100)` |
| Brilho / orientação (obter) | `IOS.GetBrightness()` / `IOS.GetOrientation()` | `Android.BrightnessJSON()` / `Android.OrientationJSON()` |
| Notificação local | `IOS.PostNotification(json)` | `Android.Notify(json)` |
| Armazenamento seguro (gravar) | `IOS.SecureSet(key, value)` | `Android.SecureSet(json)` |
| Execução em segundo plano | `IOS.BeginBackgroundTask` / `EndBackgroundTask` | `Android.StartForegroundService` / `StopForegroundService` |

@note{type="tip"}
A interface `MobileManager` é o contrato no qual `application.Mobile` se baseia. Como ambos os gerenciadores de plataforma precisam implementá-la, qualquer método listado acima tem a garantia de manter uma assinatura idêntica no iOS e no Android — se elas algum dia divergirem, a compilação da plataforma falhará.

@end
