---
title: "Roteiro"
description: "Status do projeto Wails v3, funcionalidades planejadas e como contribuir"
slug: "status"
sourcePath: "status.md"
---

## Wails v3.0.0 — Estável

Consulte o [registro de alterações](/changelog/) para ver o status mais recente.

Esta documentação abrange o Wails v3.0.0. Após o v3.0.0, correções e melhorias são disponibilizadas em versões semanais, com prévias nightly datadas entre elas e hotfixes para problemas urgentes. Consulte a [política de lançamentos](/contributing/release-policy/).

## Compatibilidade com desktops

O contrato de compatibilidade do v3.0.0 abrange aplicativos para desktop:

| Plataforma | Destinos compatíveis | Requisitos e observações |
| --- | --- | --- |
| Windows | amd64 e arm64 | Runtime do WebView2 |
| macOS | Intel e Apple Silicon | As versões do macOS e do WebKit documentadas no guia de instalação |
| Linux | amd64 e arm64 | GTK4 + WebKitGTK 6.0 por padrão; GTK3 + WebKit2GTK 4.1 continua sendo uma opção legada `-tags gtk3` durante a série v3.0.x e é removida no v3.1 |

Todos os destinos exigem Go 1.26 ou posterior para desenvolvimento. O suporte a Android e iOS continua experimental e está fora do contrato de compatibilidade estável para desktops.

## Como você pode contribuir

- Teste o Wails v3.0.0 e relate bugs reproduzíveis
- Contribua com a documentação e os exemplos
- Participe das discussões e dê feedback sobre os rascunhos de WEPs
- Envie pull requests para correções de bugs, documentação ou WEPs aceitas

Contribuições da comunidade são bem-vindas. Se quiser ajudar com esses objetivos, participe das discussões da comunidade. Propostas de novas funcionalidades devem ser apresentadas em um pull request de WEP, não em uma issue de solicitação de funcionalidade.

## Feedback e atualizações

Este roteiro está sujeito a alterações com base no feedback da comunidade e nas prioridades do projeto. Vamos atualizá-lo regularmente para refletir o progresso e as mudanças de direção. Relate problemas reproduzíveis como issues; proponha novas funcionalidades por meio de um pull request de [WEP (Proposta de Aprimoramento do Wails)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md).
