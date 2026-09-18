---
title: "Kira"
description: "Um cliente de desktop nativo para macOS que reúne AWS — ECS, RDS, S3, DynamoDB e muito mais — em uma única interface operada pelo teclado"
slug: "community/showcase/kira"
sourcePath: "community/showcase/kira.md"
---

**[Kira](https://kira.thiennguyen.dev)** é um **cliente de desktop nativo para macOS voltado à AWS** — *"AWS, sem complicações"*. Desenvolvido com **Go, Wails e React + TypeScript**, ele consolida as operações da AWS em um único aplicativo operado pelo teclado. Autentique-se uma vez pelo AWS SSO e gerencie a infraestrutura de várias contas sem precisar alternar entre o console, a CLI e diversas ferramentas de banco de dados.

Tudo começou com uma necessidade pessoal: alternar entre abas do navegador, comandos `aws` e um cliente SQL separado apenas para disponibilizar uma alteração parecia lento. O Kira reúne esses fluxos de trabalho em uma única janela nativa e rápida — entre, escolha uma conta e acesse tudo o que precisar com apenas uma tecla.

![Visão geral de várias contas no Kira — contas de produção, homologação e desenvolvimento em um só lugar](/assets/showcase-images/kira_screenshot_1.png)

## Principais destaques

- **AWS SSO para várias contas** - Entre uma vez e alterne entre contas e regiões em um só lugar
- **ECS** - Explore clusters, serviços e tarefas; reimplante, escale e reverta serviços; consulte definições de tarefas; monitore métricas dos serviços; e abra um shell interativo pelo ECS Exec
- **Bancos de dados** - Execute SQL no RDS, faça consultas e varreduras no DynamoDB e conecte-se ao PostgreSQL, MySQL e Redshift — com tunelamento SSH seguro e credenciais armazenadas nas Chaves do macOS
- **S3** - Navegue por buckets e prefixos; visualize, envie, baixe, copie, renomeie e exclua objetos; e crie pastas
- **Secrets Manager** - Liste segredos e recupere seus valores para a conta ativa
- **CloudWatch Logs** - Acompanhe e pesquise fluxos de logs
- **Smart Query** - Geração opcional de SQL assistida por IA, disponibilizada pela CLI `claude`
- **Extensões** - Instale pacotes `.kext` personalizados que adicionam botões de ação implementados por pequenos scripts em Go
- **Navegação rápida** - Uma tecla de atalho global para abrir o aplicativo, uma paleta de comandos `Cmd+K` e links diretos `kira://`

## Uma visão mais detalhada

Acompanhe seus serviços do ECS em tempo real — integridade das tarefas, CPU e memória, estado da implantação — e reimplante, escale ou reverta sem sair da lista.

![Kira exibindo serviços do ECS com métricas em tempo real](/assets/showcase-images/kira_screenshot_17.png)

Navegue pelo S3 como em um gerenciador de arquivos. Visualize objetos, inspecione metadados e versões e envie, baixe, renomeie ou exclua itens diretamente na interface.

![Navegador de objetos do S3 do Kira com visualização de objetos e metadados](/assets/showcase-images/kira_screenshot_5.png)

Distribuído como um `.dmg` assinado e autenticado pela Apple para macOS.

[Visite o Kira](https://kira.thiennguyen.dev) | [Leia a documentação](https://docs.kira.thiennguyen.dev)
