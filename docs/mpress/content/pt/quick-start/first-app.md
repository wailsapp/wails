---
title: "Seu primeiro aplicativo"
description: "Crie um aplicativo Wails funcional em 10 minutos"
slug: "quick-start/first-app"
sourcePath: "quick-start/first-app.md"
---

Vamos criar um aplicativo simples de saudação que demonstra os principais conceitos do Wails:

- Backend em Go gerenciando a lógica
- Frontend chamando funções Go
- Bindings com segurança de tipos
- Hot reload durante o desenvolvimento

**Tempo para concluir:** 10 minutos

@note{type="tip" title="Dica de desempenho para usuários do Windows 11"}
Considere usar uma [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) para armazenar seus projetos. As Dev Drives são otimizadas para cargas de trabalho de desenvolvimento e podem melhorar significativamente os tempos de compilação e as velocidades de acesso ao disco em até 30% em comparação com unidades NTFS comuns.

@end

## Crie seu projeto

@steps
### Gere o projeto
```bash
wails3 init -n myapp
cd myapp
```

Isso cria um novo projeto com o template padrão Vanilla + Vite (HTML/CSS/TypeScript com o empacotador Vite).

@note{type="tip" title="Outros templates"}
Experimente `-t react`, `-t vue` ou `-t svelte` para usar seu framework preferido. Por padrão, esses templates usam TypeScript; para usar JavaScript puro, use `-t vanilla-js` ou `-t react-js`. Execute `wails3 init -l` para ver todos os templates disponíveis ou [use seu próprio framework de frontend](/guides/dev/frontend-frameworks/).

@end

### Entenda a estrutura do projeto
```
myapp/
├── main.go              # Application entry point
├── greetservice.go      # Greet service
├── frontend/            # Your UI code
│   ├── index.html       # HTML entry point
│   ├── src/
│   │   └── main.ts      # Frontend TypeScript
│   ├── public/
│   │   └── style.css    # Styles
│   ├── package.json     # Frontend dependencies
│   ├── tsconfig.json    # TypeScript configuration
│   └── vite.config.ts   # Vite bundler config
├── build/               # Build configuration
└── Taskfile.yml         # Build tasks
```

### Execute o aplicativo
```bash
wails3 dev
```

@note{type="info" title="Primeira execução"}
A primeira execução pode demorar mais do que o esperado, pois instala as dependências do frontend, gera os bindings etc. As execuções seguintes são muito mais rápidas.

@end

O aplicativo abre e exibe uma interface de saudação. Digite seu nome e clique em "Saudar": o backend em Go processa a entrada e retorna uma saudação.

@end

## Como funciona

Vamos entender o código que faz isso funcionar.

### O backend em Go

Abra `greetservice.go`:

```go {title="greetservice.go"}
package main

import (
	"fmt"
)

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
```

**Principais conceitos:**

1. **Serviço** — Uma struct Go com métodos exportados
2. **Método exportado** — `Greet` começa com letra maiúscula, o que o disponibiliza ao frontend
3. **Lógica simples** — Recebe um nome e retorna uma saudação
4. **Segurança de tipos** — Os tipos de entrada e saída são definidos

@note{type="tip" title="Entenda serviços e bindings"}
**Serviços** são módulos Go autocontidos que expõem funcionalidades ao frontend. Eles são simplesmente structs Go comuns com métodos exportados, registradas no campo `Services` da configuração do aplicativo.

**Bindings** são o SDK TypeScript/JavaScript gerado automaticamente que permite ao frontend chamar esses serviços. Quando você executa `wails3 dev` ou `wails3 build`, o Wails analisa os serviços registrados e gera bindings com segurança de tipos em `frontend/bindings/`.

Pense nos serviços como a API do seu backend e nos bindings como a biblioteca cliente que se comunica com ela.

@end

### Registro do serviço

Abra `main.go` e localize o registro do serviço:

```go {title="main.go" highlight="4-6"}
err := application.New(application.Options{
    Name: "myapp",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    // ... other options
})
```

Isso registra seu `GreetService` no Wails, disponibilizando todos os seus métodos exportados ao frontend.

### O frontend

Abra `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
import {GreetService} from "../bindings/changeme";

window.greet = async () => {
    const nameElement = document.getElementById('name');
    const resultElement = document.getElementById('result');

    const name = nameElement.value;
    if (!name) {
        return;
    }

    try {
        const result = await GreetService.Greet(name);
        resultElement.innerText = result;
    } catch (err) {
        console.error(err);
    }
};
```

**Principais conceitos:**

1. **Bindings gerados automaticamente** — `GreetService` é importado do código gerado
2. **Chamadas com segurança de tipos** — Os nomes e as assinaturas dos métodos correspondem ao seu código Go
3. **Assíncrono por padrão** — Todas as chamadas Go retornam Promises
4. **Tratamento de erros** — Os erros do Go são capturados em try/catch

@note{type="info" title="Onde estão os bindings?"}
Os bindings gerados ficam em `frontend/bindings/`. Eles são criados automaticamente quando você executa `wails3 dev` ou `wails3 build`.

**Nunca edite esses arquivos manualmente** — eles são gerados novamente a cada compilação.

@end

## Personalize seu aplicativo

Vamos adicionar um novo recurso para entender o fluxo de trabalho.

### Adicione o recurso "Saudar vários"

@steps
### Adicione o método a GreetService
Adicione o seguinte a `greetservice.go`:

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = fmt.Sprintf("Hello %s!", name)
    }
    return greetings
}
```

### O aplicativo será recompilado automaticamente
Salve o arquivo, e `wails3 dev` recompilará automaticamente seu código Go e reiniciará o aplicativo.

@note{type="info" title="Recompilação automática"}
Alterações no código Go acionam automaticamente a recompilação e a reinicialização. Alterações no frontend são recarregadas a quente sem reinicialização.

@end

### Use-o no frontend
Adicione o seguinte a `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

Abra o console do navegador e chame `greetMany()` — você verá o array de saudações.

@end

## Compile para produção

Quando estiver pronto para distribuir seu aplicativo:

```bash
wails3 build
```

**O que isso faz:**

- Compila o código Go com otimizações
- Compila o frontend para produção (minificado)
- Cria um executável nativo em `bin/`

@tabs{sync-key="os"}
[Windows]
**Saída:** `bin/myapp.exe`

Clique duas vezes para executar. Não é necessário instalar dependências (o WebView2 faz parte do Windows).

[macOS]
**Saída:** `bin/myapp.app`

Arraste para a pasta Aplicativos ou clique duas vezes para executar.

[Linux]
**Saída:** `bin/myapp`

Execute com `./bin/myapp` ou crie um arquivo `.desktop` para o inicializador.

@end

@note{type="tip" title="Compilações multiplataforma"}
Quer compilar para outras plataformas? Consulte [Compilações multiplataforma →](/guides/build/cross-platform/)

@end

## O que aprendemos

**Estrutura do projeto**

- `main.go` para o backend em Go
- `frontend/` para o código da interface do usuário
- `Taskfile.yml` para tarefas de compilação

**Serviços**

- Crie structs em Go com métodos exportados
- Registre com `application.NewService()`
- Os métodos ficam automaticamente disponíveis no frontend

**Bindings**

- Definições TypeScript geradas automaticamente
- Chamadas de função com segurança de tipos
- Assíncrono por padrão (Promises)

**Fluxo de trabalho de desenvolvimento**

- `wails3 dev` para recarregamento automático
- Alterações no código Go recompilam e reiniciam o aplicativo automaticamente
- Alterações no frontend são recarregadas instantaneamente

---

**Tem alguma dúvida?** Entre no [Discord](https://discord.gg/JDdSxwjhGf) e pergunte à comunidade.
