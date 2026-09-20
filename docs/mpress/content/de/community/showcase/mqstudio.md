---
title: "MQ Studio"
description: "Desktop-Client mit lokaler Datenhaltung für RocketMQ, RabbitMQ, Kafka und weitere Systeme"
slug: "community/showcase/mqstudio"
sourcePath: "community/showcase/mqstudio.md"
---

![Dialog für eine neue Verbindung in MQ Studio mit den Treibern für RocketMQ, Kafka und RabbitMQ](/assets/showcase-images/mqstudio.webp)

**[MQ Studio](https://mq-studio.amigoer.com)** ist ein **Desktop-Client für Nachrichtenwarteschlangen mit lokaler Datenhaltung**, entwickelt mit **Go, Wails und React**. Jeder Broker bringt seine eigene Konsole mit: RocketMQ hat eine, Kafka eine andere, und RabbitMQ liefert ein Verwaltungs-Plugin. Unterschiedliche Oberflächen, unterschiedliche Begriffe und jeweils ein Dienst, der bereitgestellt und betrieben werden muss. MQ Studio ersetzt sie durch eine einzige Anwendung: Jeder Broker wird über einen Treiber hinter derselben Oberfläche angesprochen. Seiten und Arbeitsabläufe bleiben dadurch unabhängig vom verbundenen System gleich.

## Wichtige Funktionen

- **Eine Oberfläche für alle Broker** — derzeit RocketMQ, RabbitMQ und Kafka; Pulsar, NATS, MQTT und SQS sind geplant.
- **Topics, Warteschlangen und Nachrichten** — Topics, Warteschlangen, Exchanges und Bindings untersuchen; Nachrichten abfragen und verfolgen, Protokolle mitlesen, Nachrichten mit Schlüsseln und Headern erzeugen, erneut senden und unzustellbare Nachrichten bearbeiten.
- **Consumer und Rückstand** — Gruppen, Clients, Abonnements und Rückstand pro Partition, einschließlich Zurücksetzen von Offsets und Behandlung von Wiederholungsversuchen und Dead-Letter-Queues (DLQ).
- **Cluster und Warnungen** — Broker-Zustand, Laufzeitmetriken, Durchsatz, Speicherplatzbelegung und native Desktop-Benachrichtigungen.
- **Transparente Fähigkeiten** — jeder Treiber gibt an, was sein Endpunkt tatsächlich kann. Die Oberfläche bietet nur Vorgänge an, die der Broker unterstützt.
- **Standardmäßig privat** — die Konfiguration bleibt auf Ihrem Gerät, und Zugangsdaten werden verschlüsselt gespeichert.

Es gibt keine Serverkomponente bereitzustellen, keine Webkonsole zu betreiben und keine Telemetrie. Wails macht dies möglich: Die Verwaltungsclients dieser Broker sind Go-Bibliotheken. Die Treiberschicht spricht sie direkt im selben Prozess an, während die Oberfläche eine normale React-App bleibt — eine Binärdatei pro Plattform statt eines zu betreibenden Dienstes.

Verfügbar für macOS, Windows und Linux, auf Englisch und Chinesisch.

[Website](https://mq-studio.amigoer.com) |
[GitHub](https://github.com/amigoer/mq-studio) |
[Herunterladen](https://github.com/amigoer/mq-studio/releases/latest)
