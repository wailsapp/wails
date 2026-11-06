---
title: "Triagem por mantenedores"
description: "Faça a triagem de bugs, relatos sobre a documentação e WEPs de forma consistente"
slug: "contributing/maintainer-triage"
sourcePath: "contributing/maintainer-triage.md"
---

## Finalidade

As issues registram bugs reproduzíveis e problemas na documentação. Um pull request de WEP (Wails Enhancement Proposal) registra uma funcionalidade proposta ou uma alteração no comportamento público.

## Triagem de issues

- Confirme se o relato de bug inclui uma versão lançada, a plataforma, etapas de reprodução, o comportamento esperado, o comportamento real e a saída de `wails3 doctor`.
- Atribua o rótulo `Bug` aos relatos confirmados e aplique os rótulos relevantes de versão e plataforma. Quando necessário, solicite uma reprodução mínima.
- Mantenha abertos os relatos sobre a documentação quando identificarem um defeito concreto nela; incentive a criação de um PR quando a pessoa que fez o relato puder realizar a alteração.
- Redirecione as solicitações de funcionalidades para o guia de WEP e, em seguida, feche-as. O fluxo de trabalho automatizado de redirecionamento processa issues de melhoria que receberam esse rótulo recentemente; use a mesma redação em issues mais antigas.
- Encaminhe perguntas e solicitações de suporte para o GitHub Discussions ou o Discord.

## Triagem de WEPs

1. Verifique se o PR é um rascunho intitulado `[WEP] <title>` e contém apenas o WEP e o material de apoio.
2. Verifique se ele usa o modelo de WEP, identifica uma pessoa responsável pela implementação e aborda compatibilidade, plataformas, testes, manutenção e segurança/privacidade.
3. Mantenha a discussão técnica no PR do WEP. As discussões fornecem contexto útil, mas não constituem o registro da decisão.
4. Registre a decisão dos mantenedores em um comentário no PR — aceito, rejeitado ou retirado — com uma breve justificativa.
5. Para um WEP aceito, atribua seu número, atualize o índice de WEPs, faça o merge do PR do WEP e exija que os PRs de implementação incluam um link para ele.

## Issues de melhoria existentes

Não exclua silenciosamente issues históricas de melhoria. Para cada solicitação ainda relevante, deixe o comentário de redirecionamento e feche-a; colaboradores interessados podem abrir um WEP. Feche duplicatas com um link para o WEP ou a decisão existente.
