---
title: "Protocolos de URL personalizados"
description: "Registre esquemas de URL personalizados para iniciar seu aplicativo por meio de links"
slug: "guides/distribution/custom-protocols"
sourcePath: "guides/distribution/custom-protocols.md"
---

Os protocolos de URL personalizados (também chamados de esquemas de URL) permitem iniciar seu aplicativo quando os usuários clicam em links que usam seu protocolo personalizado, como `myapp://action` ou `myapp://open/document`.

## Visão geral

Os protocolos personalizados permitem:

- **Links profundos**: iniciar seu aplicativo com dados específicos
- **Integração com navegadores**: processar links de páginas da Web
- **Links de e-mail**: abrir seu aplicativo a partir de clientes de e-mail
- **Comunicação entre aplicativos**: iniciar seu aplicativo a partir de outros aplicativos

**Exemplo**: `myapp://open/document?id=123` inicia seu aplicativo e abre o documento 123.

## Configuração

Defina os protocolos personalizados nas opções do aplicativo:

Os protocolos personalizados são declarados em `build/config.yml` (que os empacotadores de cada plataforma — macros do NSIS no Windows, manifesto MSIX, `CFBundleURLTypes` no macOS e `.desktop`/`xdg-mime` no Linux — utilizam durante o empacotamento). Não existe um tipo `application.Protocol` nem um campo `Protocols` em `application.Options`.

```yaml
# build/config.yml
protocols:
  - scheme: myapp
    description: "My Application Protocol"
```

No código Go, escute o evento `ApplicationLaunchedWithUrl` para detectar a inicialização por URL:

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
        Description: "My awesome application",
    })

    app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(e *application.ApplicationEvent) {
        handleCustomURL(e.Context().URL())
    })

    app.Run()
}

func handleCustomURL(url string) {
    // Parse and handle the custom URL
    // Example: myapp://open/document?id=123
    println("Received URL:", url)
}
```

## Manipulador de protocolo

Escute os eventos de protocolo para processar URLs recebidas:

```go
app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(e *application.ApplicationEvent) {
    url := e.Context().URL()

    // Parse the URL
    parsedURL, err := parseCustomURL(url)
    if err != nil {
        app.Logger.Error("Failed to parse URL:", err)
        return
    }

    // Handle different actions
    switch parsedURL.Action {
    case "open":
        openDocument(parsedURL.DocumentID)
    case "settings":
        showSettings()
    case "user":
        showUser Profile(parsedURL.UserID)
    default:
        app.Logger.Warn("Unknown action:", parsedURL.Action)
    }
})
```

## Estrutura da URL

Crie estruturas de URL claras e hierárquicas:

```
myapp://action/resource?param=value

Examples:
myapp://open/document?id=123
myapp://settings/theme?mode=dark
myapp://user/profile?username=john
```

**Práticas recomendadas:**

- Use nomes de esquema em letras minúsculas
- Use esquemas curtos e fáceis de lembrar
- Use caminhos hierárquicos para os recursos
- Inclua parâmetros de consulta para dados opcionais
- Codifique caracteres especiais para uso em URLs

## Registro por plataforma

Os protocolos personalizados são registrados de forma diferente em cada plataforma.

@tabs{sync-key="platform"}
[Windows]
### Instalador NSIS para Windows

**O Wails v3 registra automaticamente os protocolos personalizados** ao usar instaladores NSIS.

#### Registro automático

Quando você compila seu aplicativo com `wails3 build`, o instalador NSIS:

1. Registra automaticamente todos os protocolos declarados em `build/config.yml` na chave `protocols:`
2. Associa os protocolos ao executável do seu aplicativo
3. Cria as entradas apropriadas no Registro
4. Remove as associações de protocolo durante a desinstalação

**Nenhuma configuração adicional é necessária!**

#### Como funciona

O modelo do NSIS inclui macros integradas:

- `wails.associateCustomProtocols` - Registra os protocolos durante a instalação
- `wails.unassociateCustomProtocols` - Remove os protocolos durante a desinstalação

Essas macros são chamadas automaticamente de acordo com sua configuração em `Protocols`.

#### Registro manual (avançado)

Se você precisar fazer o registro manualmente (fora do NSIS):

```batch
@echo off
REM Register custom protocol
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /ve /d "URL:My Application Protocol" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp" /v "URL Protocol" /t REG_SZ /d "" /f
REG ADD "HKEY_CURRENT_USER\SOFTWARE\Classes\myapp\shell\open\command" /ve /d "\"%1\"" /f
```

#### Teste

Teste o registro do protocolo:

```powershell
# Open protocol URL from PowerShell
Start-Process "myapp://test/action"

# Or from command prompt
start myapp://test/action
```

### Pacote MSIX para Windows

Os protocolos personalizados também são registrados automaticamente quando se usa o empacotamento MSIX.

#### Registro automático

Quando você compila seu aplicativo com o MSIX, o manifesto inclui automaticamente os registros de protocolo definidos na configuração de protocolos em `build/config.yml`.

O manifesto gerado inclui:

```xml
<uap:Extension Category="windows.protocol">
  <uap:Protocol Name="myapp">
    <uap:DisplayName>My Application Protocol</uap:DisplayName>
  </uap:Protocol>
</uap:Extension>
```

#### Links universais (vinculação da Web ao aplicativo)

O Windows oferece suporte à **vinculação da Web ao aplicativo**, que funciona de forma semelhante aos Links Universais do macOS. Ao implantar seu aplicativo como um pacote MSIX, você pode habilitar links HTTPS para iniciá-lo diretamente.

@note{type="note"}
A vinculação da Web ao aplicativo exige a configuração manual do manifesto. Os esquemas de protocolo personalizados são configurados automaticamente com base em `build/config.yml`, mas os domínios associados devem ser adicionados manualmente ao manifesto MSIX.

@end

Para habilitar a vinculação da Web ao aplicativo, siga o [guia da Microsoft sobre vinculação da Web ao aplicativo](https://learn.microsoft.com/en-us/windows/apps/develop/launch/web-to-app-linking). Você precisará:

1. **Adicionar manualmente o App URI Handler ao manifesto MSIX** (`build/windows/msix/app_manifest.xml`):
  ```xml
  <uap3:Extension Category="windows.appUriHandler">
    <uap3:AppUriHandler>
      <uap3:Host Name="myawesomeapp.com"/>
    </uap3:AppUriHandler>
  </uap3:Extension>
  ```


2. **Configure `windows-app-web-link` em seu site:** hospede um arquivo `windows-app-web-link` em `https://myawesomeapp.com/.well-known/windows-app-web-link`. Esse arquivo deve conter as informações do pacote do aplicativo e os caminhos que ele processa.

Quando um link da Web para o aplicativo iniciar seu aplicativo, você receberá o mesmo evento `ApplicationLaunchedWithUrl` usado com esquemas de protocolo personalizados.

[macOS]
### Configuração do Info.plist

No macOS, os protocolos são registrados por meio do arquivo `Info.plist`.

#### Configuração automática

O Wails gera automaticamente o `Info.plist` com seus protocolos quando você compila com `wails3 build`.

Os protocolos declarados em `build/config.yml` são adicionados a:

```xml
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleURLName</key>
        <string>My Application Protocol</string>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>myapp</string>
        </array>
        <key>CFBundleTypeRole</key>
        <string>Editor</string>
    </dict>
</array>
```

#### Testes

```bash
# Open protocol URL from terminal
open "myapp://test/action"

# Check registered handlers
/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -dump | grep myapp
```

#### Links universais

Além de esquemas de protocolo personalizados, o macOS também oferece suporte a **links universais**, que permitem iniciar seu aplicativo por meio de links HTTPS comuns (por exemplo, `https://myawesomeapp.com/path`). Os links universais proporcionam uma experiência integrada entre seus aplicativos web e desktop.

@note{type="caution"}
Os links universais exigem que seu aplicativo para macOS seja **assinado digitalmente** com um certificado Apple Developer e um perfil de provisionamento válidos. Compilações sem assinatura ou com assinatura ad hoc não poderão abrir links universais. Verifique se o aplicativo está devidamente assinado antes de testá-lo.

@end

Para habilitar os links universais, siga o [guia da Apple sobre como oferecer suporte a links universais em seu aplicativo](https://developer.apple.com/documentation/xcode/supporting-universal-links-in-your-app). Você precisará:

1. **Adicionar direitos** ao seu `entitlements.plist`:
  ```xml
  <key>com.apple.developer.associated-domains</key>
  <array>
    <string>applinks:myawesomeapp.com</string>
  </array>
  ```


2. **Adicionar NSUserActivityTypes ao Info.plist**:
  ```xml
  <key>NSUserActivityTypes</key>
  <array>
    <string>NSUserActivityTypeBrowsingWeb</string>
  </array>
  ```


3. **Configurar `apple-app-site-association` em seu site:** hospede um arquivo `apple-app-site-association` em `https://myawesomeapp.com/.well-known/apple-app-site-association`.

Quando um link universal aciona seu aplicativo, você recebe o mesmo evento `ApplicationLaunchedWithUrl`, portanto o código de tratamento é idêntico ao usado para esquemas de protocolo personalizados.

[Linux]
### Entrada de desktop

No Linux, os protocolos são registrados por meio de arquivos `.desktop`.

#### Configuração automática

O Wails gera um arquivo de entrada de desktop com manipuladores de protocolo quando você compila com `wails3 build`.

**Corrigido na v3**: agora o modelo de desktop do Linux inclui corretamente o tratamento de protocolos.

O arquivo de desktop gerado inclui:

```ini
[Desktop Entry]
Type=Application
Name=My Application
Exec=/usr/bin/myapp %u
MimeType=x-scheme-handler/myapp;
```

#### Registro manual

Se necessário, instale manualmente o arquivo de desktop:

```bash
# Copy desktop file
cp myapp.desktop ~/.local/share/applications/

# Update desktop database
update-desktop-database ~/.local/share/applications/

# Register protocol handler
xdg-mime default myapp.desktop x-scheme-handler/myapp
```

#### Testes

```bash
# Open protocol URL
xdg-open "myapp://test/action"

# Check registered handler
xdg-mime query default x-scheme-handler/myapp
```

@end

## Exemplo completo

Veja um exemplo completo que trata várias ações de protocolo:

```go
package main

import (
    "fmt"
    "net/url"
    "strings"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type App struct {
    app    *application.App
    window *application.WebviewWindow
}

func main() {
    // Protocol registration lives in build/config.yml (the platform packagers
    // consume it); the application code just listens for the launch event.
    app := application.New(application.Options{
        Name:        "DeepLink Demo",
        Description: "Custom protocol demonstration",
    })

    myApp := &App{app: app}
    myApp.setup()

    app.Run()
}

func (a *App) setup() {
    // Create window
    a.window = a.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "DeepLink Demo",
        Width:  800,
        Height: 600,
        URL:    "http://wails.localhost/",
    })

    // Handle custom protocol URLs
    a.app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(e *application.ApplicationEvent) {
        a.handleDeepLink(e.Context().URL())
    })
}

func (a *App) handleDeepLink(rawURL string) {
    // Parse URL
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        a.app.Logger.Error("Failed to parse URL:", err)
        return
    }

    // Bring window to front
    a.window.Show()
    a.window.Focus()

    // Extract path and query
    path := strings.Trim(parsedURL.Path, "/")
    query := parsedURL.Query()

    // Handle different actions
    parts := strings.Split(path, "/")
    if len(parts) == 0 {
        return
    }

    action := parts[0]

    switch action {
    case "open":
        if len(parts) >= 2 {
            resource := parts[1]
            id := query.Get("id")
            a.openResource(resource, id)
        }

    case "settings":
        section := ""
        if len(parts) >= 2 {
            section = parts[1]
        }
        a.openSettings(section)

    case "user":
        if len(parts) >= 2 {
            username := parts[1]
            a.openUserProfile(username)
        }

    default:
        a.app.Logger.Warn("Unknown action:", action)
    }
}

func (a *App) openResource(resourceType, id string) {
    fmt.Printf("Opening %s with ID: %s\n", resourceType, id)
    // Emit event to frontend
    a.app.Event.Emit("navigate", map[string]string{
        "type": resourceType,
        "id":   id,
    })
}

func (a *App) openSettings(section string) {
    fmt.Printf("Opening settings section: %s\n", section)
    a.app.Event.Emit("navigate", map[string]string{
        "page":    "settings",
        "section": section,
    })
}

func (a *App) openUserProfile(username string) {
    fmt.Printf("Opening user profile: %s\n", username)
    a.app.Event.Emit("navigate", map[string]string{
        "page": "user",
        "user": username,
    })
}
```

## Integração com o frontend

Trate os eventos de navegação no frontend:

```javascript
import { Events } from '@wailsio/runtime'

// Listen for navigation events from protocol handler
Events.On('navigate', (event) => {
    const { type, id, page, section, user } = event.data

    if (type === 'document') {
        // Open document with ID
        router.push(`/document/${id}`)
    } else if (page === 'settings') {
        // Open settings
        router.push(`/settings/${section}`)
    } else if (page === 'user') {
        // Open user profile
        router.push(`/user/${user}`)
    }
})
```

## Considerações de segurança

### Valide todas as entradas

Sempre valide e higienize URLs provenientes de fontes externas:

```go
func (a *App) handleDeepLink(rawURL string) {
    // Parse URL
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        a.app.Logger.Error("Invalid URL:", err)
        return
    }

    // Validate scheme
    if parsedURL.Scheme != "myapp" {
        a.app.Logger.Warn("Invalid scheme:", parsedURL.Scheme)
        return
    }

    // Validate path
    path := strings.Trim(parsedURL.Path, "/")
    if !isValidPath(path) {
        a.app.Logger.Warn("Invalid path:", path)
        return
    }

    // Sanitize parameters
    params := sanitizeQueryParams(parsedURL.Query())

    // Process validated URL
    a.processDeepLink(path, params)
}

func isValidPath(path string) bool {
    // Only allow alphanumeric and forward slashes
    validPath := regexp.MustCompile(`^[a-zA-Z0-9/]+$`)
    return validPath.MatchString(path)
}

func sanitizeQueryParams(query url.Values) map[string]string {
    sanitized := make(map[string]string)
    for key, values := range query {
        if len(values) > 0 {
            // Take first value and sanitize
            sanitized[key] = sanitizeString(values[0])
        }
    }
    return sanitized
}
```

### Evite ataques de injeção

Nunca execute URLs diretamente como código ou SQL:

```go
// ❌ DON'T: Execute URL content
func badHandler(url string) {
    exec.Command("sh", "-c", url).Run() // DANGEROUS!
}

// ✅ DO: Parse and validate
func goodHandler(url string) {
    parsed, _ := url.Parse(url)
    action := parsed.Query().Get("action")

    // Whitelist allowed actions
    allowed := map[string]bool{
        "open":     true,
        "settings": true,
        "help":     true,
    }

    if allowed[action] {
        handleAction(action)
    }
}
```

## Testes

### Testes manuais

Teste os manipuladores de protocolo durante o desenvolvimento:

**Windows:**

```powershell
Start-Process "myapp://test/action?id=123"
```

**macOS:**

```bash
open "myapp://test/action?id=123"
```

**Linux:**

```bash
xdg-open "myapp://test/action?id=123"
```

### Testes com HTML

Crie uma página HTML de teste:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Protocol Test</title>
</head>
<body>
    <h1>Custom Protocol Test Links</h1>

    <ul>
        <li><a href="myapp://open/document?id=123">Open Document 123</a></li>
        <li><a href="myapp://settings/theme?mode=dark">Dark Mode Settings</a></li>
        <li><a href="myapp://user/profile?username=john">User Profile</a></li>
    </ul>
</body>
</html>
```

## Solução de problemas

### Protocolo não registrado

**Windows:**

- Verifique o Registro: `HKEY_CURRENT_USER\SOFTWARE\Classes\<scheme>`
- Reinstale com o instalador NSIS
- Verifique se o instalador foi executado com as permissões adequadas

**macOS:**

- Recompile o aplicativo com `wails3 build`
- Verifique `Info.plist` no pacote do aplicativo: `MyApp.app/Contents/Info.plist`
- Redefina o Launch Services: `/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill`

**Linux:**

- Verifique o arquivo de desktop: `~/.local/share/applications/myapp.desktop`
- Atualize o banco de dados: `update-desktop-database ~/.local/share/applications/`
- Verifique o manipulador: `xdg-mime query default x-scheme-handler/myapp`

### O aplicativo não inicia

**Verifique os logs:**

```go
app := application.New(application.Options{
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // ...
})
```

**Problemas comuns:**

- O aplicativo não está instalado no local esperado
- O caminho do executável no registro não corresponde ao local real
- Problemas de permissão

## Práticas recomendadas

### ✅ Faça

- **Use nomes de esquema descritivos** — `mycompany-myapp` em vez de `mca`
- **Valide todas as entradas** — nunca confie em URLs provenientes de fontes externas
- **Trate os erros de forma adequada** — registre URLs inválidas nos logs; não permita que o aplicativo falhe
- **Forneça feedback ao usuário** — Mostre qual ação foi acionada
- **Teste em todas as plataformas** — O tratamento de protocolos varia
- **Documente a estrutura das suas URLs** — Ajude usuários e integradores

### ❌ Não faça

- **Não use nomes de esquema comuns** — Evite `http`, `file`, `app` etc.
- **Não execute URLs como código** — Isso representa um enorme risco à segurança
- **Não exponha operações sensíveis** — Exija confirmação para ações destrutivas
- **Não presuma que os protocolos funcionam em todos os ambientes** — Disponibilize mecanismos alternativos
- **Não se esqueça da codificação de URLs** — Trate corretamente os caracteres especiais

## Próximas etapas

- [Empacotamento para Windows](/guides/build/windows/) — Conheça as opções do instalador NSIS
- [Associações de arquivos](/guides/file-associations/) — Abra arquivos com seu aplicativo
- [Instância única](/guides/single-instance/) — Evite várias instâncias do aplicativo

---

**Tem alguma dúvida?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos](https://github.com/wailsapp/wails/tree/master/v3/examples).
