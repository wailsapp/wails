---
title: "Tutoriais"
description: "Aprenda Wails criando aplicativos"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

Tutoriais passo a passo que ensinam os conceitos do Wails por meio da criação de aplicativos completos. Cada tutorial inclui código funcional, explicações e padrões práticos.

@note{type="tip" title="Ainda não conhece Go?"}
Conclua o [Tour de Go](https://go.dev/tour/) antes de iniciar os tutoriais.

@end

## Serviço de códigos QR

![Exemplo de código QR](/assets/qr1.png)

Aprenda os fundamentos dos serviços do Wails criando um gerador de códigos QR. Este tutorial apresenta os principais conceitos para organizar a lógica do seu aplicativo em serviços reutilizáveis.

**O que você aprenderá:**

- Como criar e estruturar um serviço do Wails
- Gerenciamento de dependências externas do Go
- Vinculação de métodos Go ao frontend
- Transferência de dados entre Go e JavaScript
- Organização do código para facilitar a manutenção

**Ideal para:** quem está usando o Wails pela primeira vez e deseja entender a arquitetura de serviços

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>Começar</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Lista de tarefas

![Aplicativo de lista de tarefas](/assets/todo-app.png)

Crie um aplicativo completo de lista de tarefas com uma interface moderna e elegante. Neste tutorial prático, você aprenderá os principais padrões do Wails criando um aplicativo útil para situações reais com JavaScript puro.

**O que você aprenderá:**

- Arquitetura baseada em serviços com gerenciamento de estado seguro para threads
- Operações CRUD (criar, ler, atualizar e excluir)
- Vinculações com segurança de tipos entre Go e JavaScript
- Criação de interfaces modernas sem a complexidade de frameworks
- Padrões adequados de tratamento de erros e validação

**Tempo para concluir:** cerca de 20 minutos

**Ideal para:** seu primeiro aplicativo Wails completo — perfeito para compreender os fundamentos antes de adicionar a complexidade de um framework

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>Começar</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Notas

![Aplicativo de notas](/assets/notes-app.png)

Crie um aplicativo no estilo do Apple Notes com caixas de diálogo nativas para arquivos e salvamento automático. Este tutorial demonstra recursos específicos de desktop, como operações de arquivo, caixas de diálogo nativas e padrões profissionais de interface.

**O que você aprenderá:**

- Caixas de diálogo nativas para arquivos (Salvar, Abrir e Informações)
- Persistência de dados baseada em JSON
- Padrões de salvamento automático com debounce
- Layouts profissionais de desktop com duas colunas
- Como trabalhar com operações do sistema de arquivos em Go

**Tempo para concluir:** cerca de 30 minutos

**Ideal para:** aprender recursos específicos de desktop, como operações de arquivo e caixas de diálogo nativas do sistema operacional

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>Começar</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### Aplicativo Wails com atualização automática

Adicione atualizações automáticas dentro do próprio aplicativo a um aplicativo Wails, começando com um `wails3 init` novo e avançando até a verificação de versões assinadas e a substituição do binário no modo auxiliar. Usa o GitHub Releases como fonte das atualizações.

**O que você aprenderá:**

- Como `app.Updater` se integra a um aplicativo Wails
- Configuração do provedor do GitHub Releases
- Publicação de versões com `SHA256SUMS` para verificação de resumo criptográfico
- Adição de assinaturas Ed25519 para oferecer resistência a adulterações
- Personalização da janela padrão por meio de CSS, HTML personalizado ou BYO
- Verificações periódicas em segundo plano com `CheckInterval`

**Tempo para concluir:** cerca de 25 minutos

**Ideal para:** distribuir um aplicativo de desktop atualizável — abrange todo o pipeline de lançamento, não apenas a API

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>Começar</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
