---
title: "Guia de eventos"
description: "Um guia prático para usar eventos no Wails v3 na comunicação e no gerenciamento do ciclo de vida do aplicativo"
slug: "guides/events-reference"
sourcePath: "guides/events-reference.md"
---

**OBSERVAÇÃO: Este guia está em elaboração**

## Guia de eventos

Os eventos são o principal mecanismo de comunicação nos aplicativos Wails. Eles permitem que diferentes partes do aplicativo se comuniquem sem um acoplamento rígido. Este guia apresenta tudo o que você precisa saber para usar eventos com eficiência no seu aplicativo Wails.

## Entenda os eventos do Wails

Pense nos eventos como mensagens transmitidas por todo o aplicativo. Qualquer parte do aplicativo pode escutar essas mensagens e reagir de acordo. Isso é particularmente útil para:

- **Responder a alterações da janela**: saiba quando a janela é minimizada, maximizada ou movida
- **Tratar eventos do sistema**: reaja a alterações de tema ou a eventos de energia
- **Lógica personalizada do aplicativo**: crie seus próprios eventos para recursos como atualizações de dados ou ações do usuário
- **Comunicação entre componentes**: permita que diferentes partes do aplicativo se comuniquem sem dependências diretas

## Convenção de nomenclatura de eventos

Todos os eventos do Wails seguem um padrão de namespace para indicar claramente sua origem:

- `common:` — Eventos multiplataforma que funcionam no Windows, macOS e Linux
- `windows:` — Eventos específicos do Windows
- `mac:` — Eventos específicos do macOS\
- `linux:` — Eventos específicos do Linux

Por exemplo:

- `common:WindowFocus` — A janela recebeu foco (funciona em todas as plataformas)
- `windows:APMSuspend` — O sistema está entrando em suspensão (somente no Windows)
- `mac:ApplicationDidBecomeActive` — O aplicativo tornou-se ativo (somente no macOS)

## Primeiros passos com eventos

### Escutar eventos (frontend)

O caso de uso mais comum é escutar eventos no código do frontend:

```javascript
import { Events } from '@wailsio/runtime';

// Listen for when the window gains focus
Events.On('common:WindowFocus', () => {
    console.log('Window is now focused!');
    // Maybe refresh some data or resume animations
});

// Listen for theme changes
Events.On('common:ThemeChanged', (event) => {
    console.log('Theme changed:', event.data);
    // Update your app's theme accordingly
});

// Listen for custom events from your Go backend
Events.On('my-app:data-updated', (event) => {
    console.log('Data updated:', event.data);
    // Update your UI with the new data
});
```

### Emitir eventos (backend)

No código Go, você pode emitir eventos que o frontend pode escutar:

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "time"
)

func (s *Service) UpdateData() {
    // Do some data processing...

    // Notify the frontend
    app := application.Get()
    app.Event.Emit("my-app:data-updated",
        map[string]interface{}{
            "timestamp": time.Now(),
            "count": 42,
        },
	)
}
```

### Emitir eventos (frontend)

Embora seja menos comum, você também pode emitir eventos no frontend que o código Go pode escutar:

```javascript
import { Events } from '@wailsio/runtime';

// Event without data
Events.Emit('myapp:close-window')

// Event with data
Events.Emit('myapp:disconnect-requested', 'id-123')
```

Se você usa TypeScript no frontend e [registra eventos tipados](#eventos-tipados-com-segurana-de-tipos) no código Go, terá preenchimento automático e verificação dos nomes de eventos, além da verificação dos tipos de dados.

### Remover listeners de eventos

Sempre remova os listeners de eventos quando eles não forem mais necessários:

```javascript
import { Events } from '@wailsio/runtime';

// Store the handler reference
const focusHandler = () => {
    console.log('Window focused');
};

// Add the listener — capture the unsubscribe function it returns
const unsubscribe = Events.On('common:WindowFocus', focusHandler);

// Later, remove this specific listener via the returned unsubscribe
unsubscribe();

// Or remove ALL listeners for one (or more) event names — Events.Off takes only event-name strings
Events.Off('common:WindowFocus');
// Events.Off('common:WindowFocus', 'common:WindowLostFocus'); // variadic
// Events.OffAll(); // remove every listener for every event (no args)
```

## Casos de uso comuns

### 1. Pausar e retomar conforme o foco da janela

Muitos aplicativos precisam pausar determinadas atividades quando a janela perde o foco:

```javascript
import { Events } from '@wailsio/runtime';

let animationRunning = true;

Events.On('common:WindowLostFocus', () => {
    animationRunning = false;
    pauseBackgroundTasks();
});

Events.On('common:WindowFocus', () => {
    animationRunning = true;
    resumeBackgroundTasks();
});
```

### 2. Responder a alterações de tema

Mantenha o aplicativo sincronizado com o tema do sistema:

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:ThemeChanged', (event) => {
    const isDarkMode = event.data.isDarkMode;

    if (isDarkMode) {
        document.body.classList.add('dark-theme');
        document.body.classList.remove('light-theme');
    } else {
        document.body.classList.add('light-theme');
        document.body.classList.remove('dark-theme');
    }
});
```

### 3. Tratar arquivos arrastados e soltos

Permita que o aplicativo aceite arquivos arrastados:

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:WindowFilesDropped', (event) => {
    const files = event.data.files;

    files.forEach(file => {
        console.log('File dropped:', file);
        // Process the dropped files
        handleFileUpload(file);
    });
});
```

### 4. Gerenciar o ciclo de vida da janela

Responda a alterações no estado da janela:

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:WindowClosing', () => {
    // Save user data before closing
    saveApplicationState();

    // You could also prevent closing by returning false
    // from a registered window close handler
});

Events.On('common:WindowMaximise', () => {
    // Adjust UI for maximized view
    adjustLayoutForMaximized();
});

Events.On('common:WindowRestore', () => {
    // Return UI to normal state
    adjustLayoutForNormal();
});
```

### 5. Recursos específicos da plataforma

Trate eventos específicos da plataforma quando necessário:

```javascript
import { Events } from '@wailsio/runtime';

// Windows-specific power management
Events.On('windows:APMSuspend', () => {
    console.log('System is going to sleep');
    saveState();
});

Events.On('windows:APMResumeSuspend', () => {
    console.log('System woke up');
    refreshData();
});

// macOS-specific app lifecycle
Events.On('mac:ApplicationWillTerminate', () => {
    console.log('App is about to quit');
    performCleanup();
});
```

## Criar eventos personalizados

Você pode criar seus próprios eventos para atender às necessidades específicas do aplicativo.

### Backend (Go)

```go
// Emit a custom event when data changes

func (s *Service) ProcessUserData(userData UserData) error {
    // Process the data...

    app := application.Get()
    // Notify all listeners
    app.Event.Emit("user:data-processed",
        map[string]interface{}{
            "userId": userData.ID,
            "status": "completed",
            "timestamp": time.Now(),
        },
    )
    return nil
}

// Emit periodic updates
func (s *Service) StartMonitoring() {
    app := application.Get()
    ticker := time.NewTicker(5 * time.Second)
    go func() {
        for range ticker.C {
            stats := s.collectStats()
            app.Event.Emit("monitor:stats-updated", stats)
        }
    }()
}
```

### Frontend (JavaScript)

```javascript
import { Events } from '@wailsio/runtime';

// Listen for your custom events
Events.On('user:data-processed', (event) => {
    const { userId, status, timestamp } = event.data;

    showNotification(`User ${userId} processing ${status}`);
    updateUIWithNewData();
});

Events.On('monitor:stats-updated', (event) => {
    updateDashboard(event.data);
});
```

## Eventos tipados com segurança de tipos

O Wails v3 oferece suporte a eventos tipados com segurança de tipos completa do TypeScript por meio do registro de eventos e da geração automática de bindings.

### Registrar eventos personalizados

Chame `application.RegisterEvent` durante a inicialização para registrar nomes de eventos personalizados com seus tipos de dados:

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

type UserLoginData struct {
    UserID   string
    Username string
    LoginTime string
}

type MonitorStats struct {
    CPUUsage    float64
    MemoryUsage float64
}

func init() {
    // Register events with their data types
    application.RegisterEvent[UserLoginData]("user:login")
    application.RegisterEvent[MonitorStats]("monitor:stats")

    // Register events without data (void events)
    application.RegisterEvent[application.Void]("app:ready")
}
```

@note{type="caution"}
`RegisterEvent` deve ser chamado durante a inicialização e causará um panic se:

- Os argumentos não forem válidos
- O mesmo nome de evento for registrado duas vezes com tipos de dados diferentes

@end

@note{type="info"}
É seguro registrar o mesmo evento várias vezes, desde que o tipo de dados seja sempre o mesmo. Isso pode ser útil para garantir que um evento seja registrado quando qualquer um entre vários pacotes for carregado.

@end

### Benefícios do registro de eventos

Após o registro, os argumentos de dados passados para `Event.Emit` são verificados em relação ao tipo especificado. Em caso de incompatibilidade:

- Um erro é emitido e registrado no log (ou passado ao manipulador de erros registrado)
- O evento que causou o erro não será propagado
- Isso garante que o campo de dados dos eventos registrados sempre possa ser atribuído ao tipo declarado

### Modo estrito

Use a tag de build `strictevents` para habilitar avisos sobre eventos não registrados durante o desenvolvimento:

```bash
go build -tags strictevents
```

Com o modo estrito habilitado, o runtime emite no máximo um aviso por nome de evento não registrado para evitar o excesso de mensagens nos logs.

### Geração de bindings TypeScript

O gerador de bindings produz definições TypeScript e código de integração para oferecer suporte transparente a eventos tipados no frontend.

#### 1. Configure o plugin do Vite

No seu `vite.config.ts`:

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

#### 2. Gere os bindings

Execute o gerador de bindings:

```bash
wails3 generate bindings
```

Isso cria arquivos TypeScript no diretório do frontend, com criadores de eventos tipados e interfaces de dados.

#### 3. Use eventos tipados no frontend

```typescript
import { Events } from '@wailsio/runtime'
import { UserLogin, MonitorStats } from './bindings/events'

// Type-safe event emission with autocomplete
Events.Emit(UserLogin({
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
}))

// Type-safe event listening
Events.On(UserLogin, (event) => {
    // event.data is typed as UserLoginData
    console.log(`User ${event.data.Username} logged in`)
})

Events.On(MonitorStats, (event) => {
    // event.data is typed as MonitorStats
    updateDashboard({
        cpu: event.data.CPUUsage,
        memory: event.data.MemoryUsage
    })
})
```

Os eventos tipados oferecem:

- **Preenchimento automático** de nomes de eventos
- **Verificação de tipos** dos dados dos eventos
- **Erros em tempo de compilação** para tipos de dados incompatíveis
- Documentação via **IntelliSense**

## Referência de eventos

### Eventos comuns (multiplataforma)

Estes eventos funcionam em todas as plataformas:

| Evento | Descrição | Quando usar |
| --- | --- | --- |
| `common:ApplicationStarted` | O aplicativo foi totalmente iniciado | Inicializar o aplicativo e carregar o estado salvo |
| `common:WindowRuntimeReady` | O runtime do Wails está pronto | Começar a fazer chamadas à API do Wails |
| `common:ThemeChanged` | O tema do sistema foi alterado | Atualizar a aparência do aplicativo |
| `common:SystemWillSleep` | O sistema está prestes a entrar em suspensão | Persistir o estado e fechar os sockets |
| `common:SystemDidWake` | O sistema retomou a operação após a suspensão | Reconectar e atualizar dados desatualizados |
| `common:WindowFocus` | A janela recebeu foco | Retomar as atividades e atualizar os dados |
| `common:WindowLostFocus` | A janela perdeu o foco | Pausar as atividades e salvar o estado |
| `common:WindowMinimise` | A janela foi minimizada | Pausar a renderização e reduzir o uso de recursos |
| `common:WindowMaximise` | A janela foi maximizada | Ajustar o layout para tela cheia |
| `common:WindowRestore` | A janela foi restaurada após ser minimizada ou maximizada | Retornar ao layout normal |
| `common:WindowClosing` | A janela está prestes a ser fechada | Salvar os dados e liberar os recursos |
| `common:WindowFilesDropped` | Arquivos foram soltos na janela | Processar importações de arquivos |
| `common:WindowDidResize` | A janela foi redimensionada | Ajustar o layout e renderizar os gráficos novamente |
| `common:WindowDidMove` | A janela foi movida | Atualizar funcionalidades que dependem da posição |

### Eventos específicos da plataforma

#### Eventos do Windows

Principais eventos para aplicativos Windows:

| Evento | Descrição | Caso de uso |
| --- | --- | --- |
| `windows:SystemThemeChanged` | O tema do Windows foi alterado | Atualizar as cores do aplicativo |
| `windows:APMSuspend` | O sistema está sendo suspenso | Salvar o estado e pausar as operações |
| `windows:APMResumeAutomatic` | O sistema retomou a atividade (sempre disparado na retomada) | Restaurar o estado e atualizar os dados |
| `windows:APMResumeSuspend` | O sistema retomou a atividade por uma entrada do usuário (após `APMResumeAutomatic`) | Distinguir a reativação iniciada pelo usuário |
| `windows:APMPowerStatusChange` | O status de energia foi alterado | Ajustar as configurações de desempenho |

#### Eventos do macOS

Eventos importantes de aplicativos macOS:

| Evento | Descrição | Caso de uso |
| --- | --- | --- |
| `mac:ApplicationDidBecomeActive` | O aplicativo tornou-se ativo | Retomar as operações |
| `mac:ApplicationDidResignActive` | O aplicativo tornou-se inativo | Pausar as operações |
| `mac:ApplicationWillTerminate` | O aplicativo será encerrado | Realizar a limpeza final |
| `mac:ApplicationWillSleep` | O sistema está prestes a ser suspenso | Salvar o estado e fechar os sockets |
| `mac:ApplicationDidWake` | O sistema retomou a atividade | Reconectar e atualizar |
| `mac:ApplicationScreensDidSleep` | As telas entraram em repouso | Pausar a renderização (diferente da suspensão do sistema) |
| `mac:ApplicationScreensDidWake` | As telas foram reativadas | Retomar a renderização |
| `mac:WindowDidEnterFullScreen` | Entrou no modo de tela cheia | Adaptar a interface para o modo de tela cheia |
| `mac:WindowDidExitFullScreen` | Saiu do modo de tela cheia | Restaurar a interface normal |

#### Eventos do Linux

Principais eventos de janela do Linux:

| Evento | Descrição | Caso de uso |
| --- | --- | --- |
| `linux:SystemThemeChanged` | O tema da área de trabalho foi alterado | Atualizar o tema do aplicativo |
| `linux:SystemWillSleep` | O sistema está prestes a ser suspenso (logind) | Salvar o estado |
| `linux:SystemDidWake` | O sistema retomou a atividade (logind) | Reconectar e atualizar |
| `linux:WindowFocusIn` | A janela recebeu o foco | Retomar as atividades |
| `linux:WindowFocusOut` | A janela perdeu o foco | Pausar atividades |
| `linux:WindowLoadStarted` | O WebView começou a carregar | Exibir o indicador de carregamento |
| `linux:WindowLoadRedirected` | O WebView foi redirecionado | Rastrear redirecionamentos de navegação |
| `linux:WindowLoadCommitted` | O WebView confirmou o carregamento | O conteúdo está sendo recebido |
| `linux:WindowLoadFinished` | O WebView concluiu o carregamento | Ocultar o indicador de carregamento e injetar JS/CSS |

## Boas práticas

### 1. Use namespaces de eventos

Ao criar eventos personalizados, use namespaces para evitar conflitos:

```javascript
import { Events } from '@wailsio/runtime';

// Good - namespaced events
Events.Emit('myapp:user:login');
Events.Emit('myapp:data:updated');
Events.Emit('myapp:network:connected');

// Avoid - generic names that might conflict
Events.Emit('login');
Events.Emit('update');
```

### 2. Remova os listeners

Sempre remova os listeners de eventos quando os componentes forem desmontados:

```javascript
import { Events } from '@wailsio/runtime';

// React example
useEffect(() => {
    const handler = (event) => {
        // Handle event
    };

    const off = Events.On('common:WindowDidResize', handler);

    // Cleanup — call the unsubscribe returned by Events.On
    return () => {
        off();
    };
}, []);
```

### 3. Considere as diferenças entre plataformas

Ao usar eventos específicos de uma plataforma, verifique a disponibilidade deles nela:

```javascript
import { Events } from '@wailsio/runtime';

// Platform-specific events can be registered unconditionally;
// they will simply never fire on unsupported platforms.
Events.On('windows:APMSuspend', handleSuspend);
Events.On('mac:ApplicationWillTerminate', handleTerminate);
```

### 4. Não use eventos em excesso

Embora os eventos sejam poderosos, não os use para tudo:

- ✅ Use eventos para: notificações do sistema, mudanças no ciclo de vida e atualizações transmitidas a vários destinatários
- ❌ Evite eventos para: retornos diretos de funções, atualizações de um único componente e operações síncronas

## Depuração de eventos

Para depurar problemas com eventos:

```javascript
import { Events } from '@wailsio/runtime';

// Log all events (development only)
if (isDevelopment) {
    const originalOn = Events.On;
    Events.On = function(eventName, handler) {
        console.log(`[Event Registered] ${eventName}`);
        return originalOn.call(this, eventName, function(event) {
            console.log(`[Event Fired] ${eventName}`, event);
            return handler(event);
        });
    };
}
```

## Fonte oficial

A lista completa dos eventos disponíveis pode ser encontrada no código-fonte do Wails:

- Eventos do frontend: [`v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts`](https://github.com/wailsapp/wails/blob/main/v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts)
- Eventos do backend: [`v3/pkg/events/events.go`](https://github.com/wailsapp/wails/blob/main/v3/pkg/events/events.go)

Sempre consulte esses arquivos para obter os nomes e a disponibilidade mais atualizados dos eventos.

## Resumo

Os eventos no Wails oferecem uma maneira poderosa e desacoplada de gerenciar a comunicação no seu aplicativo. Seguindo os padrões e as práticas deste guia, você pode criar aplicativos responsivos e cientes da plataforma, que reagem de forma fluida às mudanças do sistema e às interações do usuário.

Lembre-se: comece com eventos comuns para garantir a compatibilidade entre plataformas, adicione eventos específicos de cada plataforma quando necessário e sempre remova os listeners de eventos para evitar vazamentos de memória.
