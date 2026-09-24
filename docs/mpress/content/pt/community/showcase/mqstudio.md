---
title: "MQ Studio"
description: "Cliente de desktop com prioridade ao armazenamento local para RocketMQ, RabbitMQ, Kafka e outros"
slug: "community/showcase/mqstudio"
sourcePath: "community/showcase/mqstudio.md"
---

![Janela de nova conexão do MQ Studio mostrando os drivers de RocketMQ, Kafka e RabbitMQ](/assets/showcase-images/mqstudio.webp)

**[MQ Studio](https://mq-studio.amigoer.com)** é um **cliente de desktop para filas de mensagens que prioriza o armazenamento local**, desenvolvido com **Go, Wails e React**. Cada broker tem seu próprio console: o RocketMQ tem um, o Kafka tem outro e o RabbitMQ oferece um plugin de gerenciamento. São interfaces e vocabulários diferentes, cada qual com um serviço para implantar e manter. O MQ Studio substitui todos por um único aplicativo: cada broker é acessado por um driver por trás da mesma interface, mantendo as páginas e o fluxo de trabalho iguais, independentemente do sistema conectado.

## Principais recursos

- **Uma interface para todos os brokers** — RocketMQ, RabbitMQ e Kafka atualmente, com Pulsar, NATS, MQTT e SQS no planejamento.
- **Tópicos, filas e mensagens** — inspecione tópicos, filas, exchanges e vínculos; consulte e rastreie mensagens, acompanhe logs, produza mensagens com chaves e cabeçalhos, reenvie e trate mensagens não entregues.
- **Consumidores e atraso** — grupos, clientes, assinaturas e atraso por partição, com redefinição de offsets e tratamento de novas tentativas e filas de mensagens não entregues (DLQ).
- **Cluster e alertas** — integridade dos brokers, métricas de execução, vazão, uso de disco e notificações nativas do desktop.
- **Capacidades transparentes** — cada driver declara o que seu endpoint realmente pode fazer, e a interface oferece apenas as operações compatíveis com o broker.
- **Privacidade por padrão** — a configuração fica no seu dispositivo e as credenciais são criptografadas em repouso.

Não há componente de servidor para implantar, console web para manter nem telemetria. O Wails torna isso possível: os clientes de administração desses brokers são bibliotecas Go. A camada de drivers conversa diretamente com elas no mesmo processo, enquanto a interface continua sendo um aplicativo React comum — um binário por plataforma em vez de um serviço que alguém precisa operar.

Disponível para macOS, Windows e Linux, em inglês e chinês.

[Site](https://mq-studio.amigoer.com) |
[GitHub](https://github.com/amigoer/mq-studio) |
[Baixar](https://github.com/amigoer/mq-studio/releases/latest)
