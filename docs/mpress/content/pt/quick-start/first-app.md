---
title: "Seu primeiro aplicativo"
description: "Crie um aplicativo Wails funcional em 10 minutos"
slug: "quick-start/first-app"
sourcePath: "quick-start/first-app.md"
---

Vamos criar um aplicativo simples de saudação que demonstra os principais conceitos do Wails:

- Backend em Go responsável pela lógica
- Frontend que chama funções Go
- Bindings com segurança de tipos
- Recarregamento automático durante o desenvolvimento

**Tempo para concluir:** 10 minutos

## Crie seu projeto

@steps
### Gere o projeto
```bash
wails3 init -n myapp -t vanilla-js
cd myapp
```

Esta lição seleciona explicitamente Vanilla + Vite com JavaScript para que os arquivos de origem abaixo correspondam ao projeto gerado. O template `vanilla` básico usa TypeScript.

@note{type="tip" title="Outros templates"}
Experimente `-t react`, `-t vue` ou `-t svelte` com o framework de sua preferência. Por padrão, esses templates usam TypeScript; para usar JavaScript puro, escolha `-t vanilla-js` ou `-t react-js`. Execute `wails3 init -l` para ver todos os templates disponíveis ou [use seu próprio framework de frontend](/guides/dev/frontend-frameworks/).

@end

### Entenda a estrutura do projeto
```
myapp/
├── main.go              # Application entry point
├── greetservice.go      # Greet service
├── frontend/            # Your UI code
│   ├── index.html       # HTML entry point
│   ├── src/
│   │   └── main.js      # Frontend JavaScript
│   ├── public/
│   │   └── style.css    # Styles
│   ├── package.json     # Frontend dependencies
│   ├── tsconfig.json    # Editor and JavaScript tooling configuration
│   └── vite.config.js   # Vite bundler config
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

O aplicativo será aberto com uma interface de saudação. Digite seu nome e clique em "Saudar": o backend em Go processará a entrada e retornará uma saudação.

@end

@note{type="tip" title="Dica de desempenho para usuários do Windows 11"}
Considere usar o [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) para armazenar seus projetos. Os Dev Drives são otimizados para cargas de trabalho de desenvolvimento e podem melhorar significativamente o tempo de compilação e a velocidade de acesso ao disco em até 30% em comparação com unidades NTFS comuns.

@end

## Como funciona

Vamos entender o código que faz tudo isso funcionar.

### O backend em Go

Abra `greetservice.go`:

```go {title="greetservice.go"}
package main

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return "Hello " + name + "!"
}
```

**Principais conceitos:**

1. **Serviço** — Uma struct Go com métodos exportados
2. **Método exportado** — `Greet` começa com letra maiúscula, o que o torna disponível para o frontend
3. **Lógica simples** — Recebe um nome e retorna uma saudação
4. **Segurança de tipos** — Os tipos de entrada e saída são definidos

@note{type="tip" title="Entenda serviços e bindings"}
**Serviços** são módulos Go autônomos que disponibilizam funcionalidades ao frontend. Eles são simplesmente structs Go comuns com métodos exportados, registradas no campo `Services` da configuração do aplicativo.

**Bindings** são o SDK TypeScript/JavaScript gerado automaticamente que permite ao frontend chamar esses serviços. Quando você executa `wails3 dev` ou `wails3 build`, o Wails analisa os serviços registrados e gera bindings com segurança de tipos em `frontend/bindings/`.

Pense nos serviços como a API do backend e nos bindings como a biblioteca cliente que se comunica com ela.

@end

### Registre o serviço

Abra `main.go` e localize o registro do serviço:

```go {title="main.go" highlight="4-6"}
app := application.New(application.Options{
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

window.doGreet = async () => {
    const nameElement = document.getElementById('name');

    const name = nameElement.value;
    if (!name) {
        return;
    }

    try {
        const result = await GreetService.Greet(name);
        showToast(result);
    } catch (err) {
        console.error(err);
    }
};
```

Esse manipulador simplificado substitui o `window.doGreet` do projeto inicial; mantenha o helper `showToast` existente e o restante do código de inicialização. O botão HTML gerado chama `doGreet`, e `showToast` torna o resultado visível. O caminho de importação segue o nome do módulo Go (`changeme`, a menos que tenha sido substituído).

**Principais conceitos:**

1. **Bindings gerados automaticamente** — `GreetService` é importado do código gerado
2. **Chamadas com segurança de tipos** — Os nomes e as assinaturas dos métodos correspondem ao seu código Go
3. **Assíncrono por padrão** — Todas as chamadas a Go retornam Promises
4. **Tratamento de erros** — Os erros do Go são capturados em try/catch

@note{type="info" title="Onde estão os bindings?"}
Os bindings gerados ficam em `frontend/bindings/`. Eles são criados automaticamente quando você executa `wails3 dev` ou `wails3 build`.

**Nunca edite esses arquivos manualmente** — eles são gerados novamente a cada compilação.

@end

## Personalize seu aplicativo

Vamos adicionar uma nova funcionalidade para entender o fluxo de trabalho.

### Adicione a funcionalidade "Saudar vários"

@steps
### Adicione o método a GreetService
Adicione o seguinte a `greetservice.go`:

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = "Hello " + name + "!"
    }
    return greetings
}
```

### O aplicativo será recompilado automaticamente
Salve o arquivo, e `wails3 dev` recompilará automaticamente seu código Go e reiniciará o aplicativo.

@note{type="info" title="Recompilação automática"}
Alterações no código Go acionam automaticamente a recompilação e a reinicialização. Alterações no frontend são recarregadas sem reiniciar o aplicativo.

@end

### Use a funcionalidade no frontend
Adicione o seguinte a `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

Abra o console do navegador e chame `greetMany()`; você verá o array de saudações.

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

Clique duas vezes para executar. Nenhuma dependência é necessária (o WebView2 faz parte do Windows).

[macOS]
**Saída:** `bin/myapp.app`

Arraste para a pasta Aplicativos ou clique duas vezes para executar.

[Linux]
**Saída:** `bin/myapp`

Execute com `./bin/myapp` ou crie um arquivo `.desktop` para seu inicializador.

@end

@note{type="tip" title="Compilações multiplataforma"}
Quer compilar para outras plataformas? Consulte [Compilações multiplataforma →](/guides/build/cross-platform/)

@end

## O que aprendemos

**Estrutura do projeto**

- `main.go` para o backend em Go
- `frontend/` para o código da interface
- `Taskfile.yml` para tarefas de build

**Serviços**

- Crie structs Go com métodos exportados
- Registre com `application.NewService()`
- Os métodos ficam disponíveis automaticamente no frontend

**Bindings**

- Definições TypeScript geradas automaticamente
- Chamadas de função com segurança de tipos
- Assíncronas por padrão (Promises)

**Fluxo de desenvolvimento**

- `wails3 dev` para recarregamento automático
- Alterações no código Go recompilam e reiniciam o aplicativo automaticamente
- Alterações no frontend são recarregadas instantaneamente

---

**Dúvidas?** Entre no [Discord](https://discord.gg/JDdSxwjhGf) e pergunte à comunidade.
