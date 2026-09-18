---
title: "Roteiro"
description: "Status do projeto Wails v3, funcionalidades planejadas e como contribuir"
slug: "status"
sourcePath: "status.md"
---

## Status atual: Beta

Consulte o [Registro de alterações](/changelog/) para ver o status mais recente.

Nosso objetivo é chegar a uma versão v3.0 estável. Este roteiro descreve as principais funcionalidades e melhorias que precisamos implementar antes da versão final. Este é um documento dinâmico e pode ser atualizado à medida que as prioridades mudarem ou surgirem novas informações.

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

Este roteiro está sujeito a alterações com base no feedback da comunidade e nas prioridades do projeto. Vamos atualizá-lo regularmente para refletir o progresso e as mudanças de direção. Relate problemas reproduzíveis como issues; proponha novas funcionalidades por meio de um PR de [WEP (Proposta de Aprimoramento do Wails)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md).
