---
title: "Tutoriais"
description: "Aprenda Wails criando aplicações"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

Escolha uma aula guiada quando quiser aprender Wails criando algo. Cada aula acompanha um projeto até chegar a um resultado concreto.

Se esta é sua primeira visita, conclua primeiro [Instalação](/quick-start/installation/) e [Sua primeira aplicação](/quick-start/first-app/).

## Escolha um tutorial

Depois de criar sua primeira aplicação, siga as aulas para desktop nesta ordem:

| Aula | O que você praticará |
| --- | --- |
| [Serviço de código QR](/tutorials/01-creating-a-service/) | Exposição da lógica em Go por meio de um serviço e chamada desse serviço pelo frontend |
| [Lista de tarefas](/tutorials/02-todo-vanilla/) | Criação, atualização e exclusão de dados da aplicação |
| [Aplicação de notas](/tutorials/03-notes-vanilla/) | Uso de caixas de diálogo nativas e salvamento de dados em arquivos |
| [Aplicação com atualização automática](/tutorials/04-self-update-a-wails-app/) | Verificação e aplicação de atualizações do GitHub Releases |

Para uma aula opcional sobre dispositivos móveis, siga [Sua primeira aplicação móvel](/guides/mobile/first-mobile-app/). O suporte a dispositivos móveis é experimental.

Se você estiver trabalhando em uma aplicação existente e precisar concluir uma tarefa específica, escolha um [guia prático](/quick-start/next-steps/#features). Para conhecer os fundamentos de design, consulte [Conceitos fundamentais](/concepts/architecture/).

@note{type="tip" title="Ainda não conhece Go?"}
Conclua o [Tour de Go](https://go.dev/tour/) antes de começar os tutoriais.

@end

## Serviço de código QR

![Exemplo de código QR](/assets/qr1.png)

Aprenda os fundamentos dos serviços do Wails criando um gerador de códigos QR. Este tutorial apresenta os conceitos essenciais para organizar a lógica da sua aplicação em serviços reutilizáveis.

**O que você aprenderá:**

- Como criar e estruturar um serviço do Wails
- Gerenciamento de dependências externas do Go
- Vinculação de métodos Go ao frontend
- Transferência de dados entre Go e JavaScript
- Organização do código para facilitar a manutenção

**Ideal para:** quem está começando a usar Wails e quer entender a arquitetura de serviços

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>Começar</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Lista de tarefas

![Aplicação de lista de tarefas](/assets/todo-app.png)

Crie uma aplicação completa de lista de tarefas com uma interface bonita e moderna. Neste tutorial prático, você aprenderá os principais padrões do Wails por meio de uma aplicação real e funcional que usa JavaScript puro.

**O que você aprenderá:**

- Arquitetura baseada em serviços com gerenciamento de estado seguro para concorrência
- Operações CRUD (Criar, Ler, Atualizar e Excluir)
- Vinculações com segurança de tipos entre Go e JavaScript
- Criação de interfaces modernas sem a complexidade de frameworks
- Padrões adequados de tratamento de erros e validação

**Tempo para concluir:** cerca de 20 minutos

**Ideal para:** sua primeira aplicação Wails completa — perfeita para entender os fundamentos antes de adicionar a complexidade de um framework

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>Começar</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Notas

![Aplicação de notas](/assets/notes-app.png)

Crie uma aplicação no estilo do Apple Notes, com caixas de diálogo nativas para arquivos e funcionalidade de salvamento automático. Este tutorial demonstra recursos específicos de desktop, como operações com arquivos, caixas de diálogo nativas e padrões profissionais de interface.

**O que você aprenderá:**

- Caixas de diálogo nativas para arquivos (Salvar, Abrir e Informações)
- Persistência de dados baseada em JSON
- Padrões de salvamento automático com debounce
- Layouts profissionais de desktop com duas colunas
- Uso de operações do sistema de arquivos em Go

**Tempo para concluir:** cerca de 30 minutos

**Ideal para:** aprender recursos específicos de desktop, como operações com arquivos e caixas de diálogo nativas do sistema operacional

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>Começar</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Aplicação Wails com atualização automática

Adicione atualizações automáticas dentro da aplicação a uma aplicação Wails, começando com um `wails3 init` novo, passando pela verificação assinada de versões e chegando à substituição do binário no modo auxiliar. Usa o GitHub Releases como fonte das atualizações.

**O que você aprenderá:**

- Como `app.Updater` se integra a uma aplicação Wails
- Configuração do provedor do GitHub Releases
- Publicação de versões com `SHA256SUMS` para verificação de resumo criptográfico
- Adição de assinatura Ed25519 para oferecer resistência contra adulterações
- Personalização da janela padrão por meio de CSS, HTML personalizado ou BYO
- Verificações periódicas em segundo plano com `CheckInterval`

**Tempo para concluir:** cerca de 25 minutos

**Ideal para:** distribuir uma aplicação para desktop que possa ser atualizada — abrange todo o pipeline de lançamento, não apenas a API

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>Começar</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
