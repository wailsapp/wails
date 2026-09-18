---
title: "Arquitetura do Wails v3"
description: "Diagramas e explicações detalhadas de cada componente interno do Wails v3"
slug: "contributing/architecture"
sourcePath: "contributing/architecture.md"
---

O Wails v3 é um **framework full-stack para desktop** composto por um runtime em Go, uma ponte JavaScript, uma cadeia de ferramentas orientada por tarefas e uma coleção de modelos que permitem distribuir aplicações nativas baseadas em tecnologias web modernas.

Esta página apresenta a *visão geral* em quatro diagramas:

1. **Arquitetura geral** – como todos os subsistemas se conectam\
2. **Fluxo do runtime** – o que acontece quando o JS chama o Go e vice-versa\
3. **Desenvolvimento vs. produção** – dois modos do servidor de assets\
4. **Implementações específicas de plataforma** – onde fica o código específico de cada sistema operacional\

---

## 1 · Arquitetura geral

**Wails v3 – Pilha de alto nível**

**[Espaço reservado para o diagrama da pilha de alto nível]**

---

## 2 · Fluxo de chamadas do runtime

**Runtime – Fluxo de chamadas entre JavaScript ⇄ Go**

**[Espaço reservado para o diagrama do fluxo de chamadas do runtime]**

Pontos principais:

- **Sem HTTP/IPC** – a ponte usa o canal em memória nativo da WebView\
- **IDs de métodos** – um hash FNV determinístico permite buscas O(1) em Go\
- **Promises** – os erros são propagados como rejeições com stack e código

---

## 3 · Fluxo de assets em desenvolvimento vs. produção

**Servidor de assets em desenvolvimento ↔ produção**

**[Espaço reservado para o diagrama do fluxo de assets]**

- Em **desenvolvimento**, o servidor encaminha caminhos desconhecidos ao servidor de recarregamento em tempo real do framework e fornece assets estáticos a partir do disco.
- Em **produção**, a mesma API usa `go:embed` como implementação subjacente, gerando um binário sem dependências.

---

## 4 · Divisão do runtime específica de plataforma

**Arquivos de runtime por sistema operacional**

**[Espaço reservado para o diagrama da divisão por plataforma]**

Cada recurso segue este padrão:

1. **Interface comum** em `pkg/application`\
2. Entrada do **processador de mensagens** em `pkg/application/messageprocessor_*.go`\
3. **Implementação por sistema operacional** em `pkg/application/*_{darwin,linux,windows}.go` (por exemplo, `webview_window_darwin.go`, `clipboard_linux.go`, `dialogs_windows.go`, `systemtray_*.go`, `mainthread_*.go`), protegida por build tags. Além disso, o Linux tem a ponte cgo em `linux_cgo.go`/`linux_cgo_gtk4.{go,c,h}`.

`internal/runtime/` contém apenas o pequeno código de integração de `runtime{_darwin,_linux,_windows,_android,_dev,_prod}.go` com build tags, além do runtime JS incorporado em `internal/runtime/desktop/`.

`internal/capabilities/` existe para declarar conjuntos de recursos por plataforma, mas não há nenhum sentinela `ErrCapability` — a disponibilização condicional de recursos é feita por build tags comuns e retornos de stubs específicos de plataforma (por exemplo, `nil` ou erros específicos do recurso).

---

## Resumo

Estes diagramas mostram **onde fica o código**, **como os dados circulam** e **quais camadas são responsáveis por quais tarefas**. Mantenha-os à mão ao explorar as páginas detalhadas a seguir – eles são seu mapa da árvore de código-fonte do Wails v3.
