---
title: "Status do projeto"
description: "Compatibilidade da versão beta do Wails v3, suporte de segurança e orientações de atualização"
slug: "status"
sourcePath: "status.md"
---

## Status atual: Beta

Consulte o [Registro de alterações](/changelog/) para ver o status mais recente.

Nosso objetivo é uma versão v3.0 estável. Wails v2 continua sendo a versão estável atual e segue recebendo correções. Teste as versões beta com seu aplicativo antes da implantação.

## Compromisso de compatibilidade da versão Beta

O contrato da versão Beta da v3 abrange aplicativos para desktop:

| Plataforma | Arquiteturas compatíveis | Requisitos e observações |
| --- | --- | --- |
| Windows | amd64 e arm64 | Runtime do WebView2 |
| macOS | Intel e Apple Silicon | As versões do macOS e do WebKit documentadas no guia de instalação |
| Linux | amd64 e arm64 | GTK4 + WebKitGTK 6.0 por padrão; GTK3 + WebKit2GTK 4.1 continua sendo uma opção legada com `-tags gtk3` durante toda a série v3.0.x e é removida na v3.1 |

Todos os destinos exigem Go 1.25 ou posterior para o desenvolvimento. O suporte a Android e iOS é experimental e não impede o lançamento da versão Beta para desktop. As APIs Beta buscam estabilidade, mas defeitos de pré-lançamento e alterações anunciadas explicitamente ainda podem ser corrigidos antes da v3.0.0.

## Como você pode contribuir

- Teste a versão Beta mais recente e relate bugs reproduzíveis
- Contribua para a documentação e os exemplos
- Participe das discussões e dê feedback sobre os rascunhos de WEPs
- Envie pull requests para correções de bugs, documentação ou WEPs aprovadas

Agradecemos as contribuições da comunidade. Se quiser ajudar com esses objetivos, participe das discussões da comunidade. Propostas de novas funcionalidades devem ser apresentadas em um PR de WEP, não em uma issue de solicitação de funcionalidade.

## Feedback e atualizações

Relate problemas reproduzíveis como issues; proponha novas funcionalidades por meio de um PR de [WEP (Proposta de Aprimoramento do Wails)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md).

## Usar a versão beta

Fixe versões exatas da CLI, do módulo Go e do runtime frontend em vez de usar `latest`. Para um projeto alpha existente, siga o [guia de atualização de alpha para beta](/migration/alpha-to-beta/).

A [política de segurança](https://github.com/wailsapp/wails/blob/master/SECURITY.md) lista as versões beta da v3 como suportadas e as versões alpha como não suportadas. Relate vulnerabilidades pelo [canal privado de relato de vulnerabilidades](https://github.com/wailsapp/wails/security/advisories/new), não em issues públicas.

## Trabalho acompanhado

- [Bugs abertos com o rótulo v3](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3ABug+label%3Av3)
- [Issues v3 abertas com o rótulo P0 ou P1](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3Av3+label%3AP0%2CP1)
- [Marcos de lançamento](https://github.com/wailsapp/wails/milestones)

Essas consultas em tempo real dependem dos rótulos das issues; não são uma lista completa nem uma promessa de data ou escopo de lançamento. Leia as issues para avaliar o impacto no seu projeto.
