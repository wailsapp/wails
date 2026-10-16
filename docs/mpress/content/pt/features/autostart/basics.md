---
title: "Inicialização automática"
description: "Registre seu aplicativo para ser iniciado quando o usuário fizer login no macOS, Windows e Linux"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## Inicialização automática

`app.Autostart` registra seu aplicativo para ser iniciado automaticamente quando o usuário faz login. Ele seleciona o mecanismo nativo adequado para cada plataforma e resolve caminhos de instalação com links simbólicos (Homebrew, Scoop), para que os registros não deixem de funcionar quando o binário for atualizado.

O registro entra em vigor no **próximo login**, não imediatamente.

## Início rápido

```go
import "github.com/wailsapp/wails/v3/pkg/application"

// Register to launch at login
if err := app.Autostart.Enable(); err != nil {
    app.Logger.Error("autostart enable failed", "error", err)
}

// Stop launching at login
if err := app.Autostart.Disable(); err != nil {
    app.Logger.Error("autostart disable failed", "error", err)
}

// Check status
enabled, err := app.Autostart.IsEnabled()
```

## API

### `Enable`

Registra o aplicativo para ser iniciado no login com as opções padrão.

```go
func (m *AutostartManager) Enable() error
```

É seguro chamar `Enable` repetidamente — o registro é sobrescrito a cada vez, portanto, é correto chamá-lo em toda inicialização se você tiver armazenado a preferência do usuário.

### `EnableWithOptions`

Faz o registro com opções personalizadas.

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions`:**

| Campo | Tipo | Descrição |
| --- | --- | --- |
| `Identifier` | `string` | Substitui o ID de registro derivado automaticamente. Consulte “Identificador” abaixo. |
| `Arguments` | `[]string` | Argumentos adicionais anexados ao caminho do executável quando ele é iniciado no login (por exemplo, `--hidden`). |

### `Disable`

Remove o registro de inicialização automática. Retorna `nil` se o aplicativo não estava registrado — a desativação é idempotente.

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

Informa se existe um registro. É rápido, pois não valida o caminho registrado.

```go
func (m *AutostartManager) IsEnabled() (bool, error)
```

### `Status`

Retorna o estado completo do registro.

```go
func (m *AutostartManager) Status() (AutostartStatus, error)
```

**`AutostartStatus`:**

| Campo | Tipo | Descrição |
| --- | --- | --- |
| `Enabled` | `bool` | Indica se existe um registro. |
| `Path` | `string` | Local em disco do artefato de registro (caminho do plist, caminho de `.desktop` ou subchave do Registro). Fica vazio quando `Enabled` é false. |
| `Strategy` | `AutostartStrategy` | Indica qual mecanismo registrou o aplicativo (consulte [Comportamento por plataforma](#comportamento-por-plataforma)). |

## Comportamento por plataforma

@tabs{sync-key="platform"}
[macOS]
Dois mecanismos são usados, dependendo de como o aplicativo foi empacotado:

- **macOS 13+, `.app`** empacotado: `SMAppService.mainAppService`. Funciona com aplicativos em sandbox e builds da Mac App Store. Não exibe solicitação de automação do TCC (a abordagem histórica com AppleScript exibia uma).
- **macOS anterior ao 13 ou binário não empacotado**: um plist de LaunchAgent gravado em `~/Library/LaunchAgents/<identifier>.plist` com `RunAtLoad=true`.

`Status()` retorna `AutostartStrategySMAppService` ou `AutostartStrategyLaunchAgent` para que o código chamador possa identificar qual caminho foi usado.

Quando o aplicativo é atualizado de não empacotado para empacotado, ambos os caminhos são verificados por `Status()` e limpos por `Disable()`, para que um LaunchAgent órfão não continue iniciando a build antiga.

[Windows]
Um valor do Registro é adicionado em `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, com o nome do valor definido como o `Identifier` da inicialização automática e os dados definidos como o caminho do executável entre aspas, seguido de quaisquer `Arguments`.

A colocação de argumentos entre aspas segue as regras de `CommandLineToArgvW` (as barras invertidas são duplicadas antes das aspas), para que caminhos com espaços ou aspas sejam convertidos e restaurados corretamente.

`Status().Strategy` retorna `AutostartStrategyRegistryRun`.

[Linux]
Uma entrada de inicialização automática XDG é gravada em `$XDG_CONFIG_HOME/autostart/<identifier>.desktop` (com o padrão `~/.config/autostart/`) com:

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

O campo `Exec` é escapado conforme a [especificação Desktop Entry do freedesktop.org](https://specifications.freedesktop.org/desktop-entry-spec/) — caracteres reservados (`"`, `` ` ``, `$`, `\\`) são escapados com uma barra invertida, e o valor é colocado entre aspas duplas quando contém espaços em branco.

`Status().Strategy` retorna `AutostartStrategyXDGAutostart`.

[iOS / Android / servidor]
Não há suporte. Todos os métodos retornam `ErrAutostartNotSupported`. Use `errors.Is(err, application.ErrAutostartNotSupported)` para detectar isso de forma adequada:

```go
if err := app.Autostart.Enable(); err != nil {
    if errors.Is(err, application.ErrAutostartNotSupported) {
        // hide the toggle in the UI
        return
    }
    // real failure — surface it
}
```

@end

## Identificador

Se `Options.Identifier` estiver vazio, um valor padrão será derivado do nome do aplicativo:

| Plataforma | Identificador padrão |
| --- | --- |
| macOS (empacotado) | O identificador do bundle do aplicativo, por exemplo, `com.example.MyApp` |
| macOS (não empacotado) | `wails.autostart.<slug>`, em que `<slug>` é derivado de `application.Options.Name` |
| Windows | Slug de `application.Options.Name` (em minúsculas, com caracteres que não sejam `A-Za-z0-9._-` removidos e espaços substituídos por hífens) |
| Linux | O mesmo slug usado no Windows |

Os identificadores devem corresponder a `^[A-Za-z0-9._-]+$` e ter no máximo 200 caracteres. O formato DNS reverso é recomendado para o macOS (ele corresponde à forma como os Labels do launchd são convencionalmente escritos).

Quando `AutostartOptions.Identifier` é substituído, o mesmo identificador é reutilizado como nome do valor do Registro no Windows e como nome do arquivo `.desktop` no Linux, de modo que uma única string identifica o registro em todas as plataformas.

## Detecção de registros obsoletos

`Disable()` e `Status()` localizam o registro **comparando o caminho do executável registrado com `os.Executable()` (resolvido por meio de quaisquer links simbólicos)**, e não consultando o identificador. Isso significa que:

- **É seguro alterar o identificador entre versões.** O registro antigo ainda pode ser localizado por `Status()` e removido por `Disable()` — desde que o caminho do executável permaneça o mesmo.
- **Uma segunda cópia do aplicativo em outro caminho não sobrescreverá o registro da primeira.** Cada localização do binário é rastreada de forma independente.
- **Instalações por links simbólicos (Homebrew, Scoop) permanecem estáveis.** `filepath.EvalSymlinks` é aplicado a `os.Executable()` antes da comparação, portanto uma atualização do Homebrew que troca o destino não deixa a entrada órfã.

O que isso *não* abrange: se o usuário mover ou renomear o binário para um caminho não relacionado, o registro antigo ficará órfão (apontará para o arquivo que agora não existe). Aplicativos distribuídos em um caminho de instalação estável não precisam se preocupar com isso; aplicativos distribuídos como binários portáteis de arquivo único devem chamar `Disable()` antes de mover a si próprios ou sempre ser iniciados por meio de um link simbólico estável.

## Exemplo

```go
package main

import (
    "errors"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Restore the user's preference on startup
    if userPrefersAutostart() {
        if err := app.Autostart.Enable(); err != nil {
            if !errors.Is(err, application.ErrAutostartNotSupported) {
                app.Logger.Error("autostart", "error", err)
            }
        }
    }

    app.Run()
}
```

Um exemplo completo e executável, com botões de status, ativação e desativação, está em [`examples/autostart/`](https://github.com/wailsapp/wails/tree/master/v3/examples/autostart).
