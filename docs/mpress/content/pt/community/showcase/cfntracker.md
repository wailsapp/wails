---
title: "CFN Tracker"
description: "Um aplicativo para desktop criado com Wails"
slug: "community/showcase/cfntracker"
sourcePath: "community/showcase/cfntracker.md"
---

![CFN Tracker](/assets/showcase-images/cfntracker.webp)

[CFN Tracker](https://github.com/williamsjokvist/cfn-tracker) — Acompanhe ao vivo as partidas de qualquer perfil do Street Fighter 6 ou V no CFN. Acesse [o site](https://cfn.williamsjokvist.se/) para começar.

## Recursos

- Acompanhamento de partidas em tempo real
- Armazenamento de registros e estatísticas das partidas
- Suporte à exibição de estatísticas ao vivo no OBS por meio de uma Fonte de navegador
- Suporte a SF6 e SFV
- Possibilidade de os usuários criarem seus próprios temas de navegador para o OBS com CSS

### Principais tecnologias usadas com o Wails

- [Task](https://github.com/go-task/task) — encapsula a CLI do Wails para facilitar o uso de comandos comuns
- [React](https://github.com/facebook/react) — escolhido por seu vasto ecossistema (radix, framer-motion)
- [Bun](https://github.com/oven-sh/bun) — usado por sua rápida resolução de dependências e seu baixo tempo de compilação
- [Rod](https://github.com/go-rod/rod) — automação de navegador sem interface gráfica para autenticação e sondagem de alterações
- [SQLite](https://github.com/mattn/go-sqlite3) — usado para armazenar partidas, sessões e perfis
- [Eventos enviados pelo servidor](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events) — um fluxo HTTP para enviar atualizações de acompanhamento às fontes de navegador do OBS
- [i18next](https://github.com/i18next/) — com um conector de back-end para fornecer objetos de localização provenientes da camada Go
- [xstate](https://github.com/statelyai/xstate) — máquinas de estados para o processo de autenticação e o acompanhamento
