---
title: "MQ Studio"
description: "Client de bureau privilégiant le stockage local pour RocketMQ, RabbitMQ, Kafka et plus encore"
slug: "community/showcase/mqstudio"
sourcePath: "community/showcase/mqstudio.md"
---

![Fenêtre de création de connexion de MQ Studio présentant les pilotes RocketMQ, Kafka et RabbitMQ](/assets/showcase-images/mqstudio.webp)

**[MQ Studio](https://mq-studio.amigoer.com)** est un **client de bureau pour les files de messages qui privilégie le stockage local**, développé avec **Go, Wails et React**. Chaque courtier possède sa propre console : RocketMQ en a une, Kafka une autre, et RabbitMQ propose un module de gestion. Les interfaces et le vocabulaire diffèrent, et chacune nécessite un service à déployer et à maintenir. MQ Studio les remplace par une seule application : chaque courtier est accessible par un pilote derrière la même interface, de sorte que les pages et le déroulement des opérations restent identiques quel que soit le système connecté.

## Points forts

- **Une interface pour tous les courtiers** — RocketMQ, RabbitMQ et Kafka aujourd’hui ; Pulsar, NATS, MQTT et SQS sont prévus.
- **Sujets, files et messages** — inspectez les sujets, files, échanges et liaisons ; recherchez et tracez les messages, suivez un journal, produisez des messages avec des clés et des en-têtes, renvoyez-les et traitez les messages en échec.
- **Consommateurs et retard** — groupes, clients, abonnements et retard par partition, avec réinitialisation des offsets et gestion des nouvelles tentatives et des files de lettres mortes (DLQ).
- **Cluster et alertes** — état des courtiers, mesures d’exécution, débit, utilisation du disque et notifications natives du bureau.
- **Des capacités clairement annoncées** — chaque pilote déclare ce que son point de connexion peut réellement faire ; l’interface ne propose que les opérations prises en charge par le courtier.
- **Confidentialité par défaut** — la configuration reste sur votre appareil et les identifiants sont chiffrés au repos.

Aucun composant serveur à déployer, aucune console web à maintenir et aucune télémétrie. Wails rend cela possible : les clients d’administration de ces courtiers sont des bibliothèques Go. Les pilotes les appellent directement dans le même processus, tandis que l’interface reste une application React classique : un binaire par plateforme, sans service à exploiter.

Disponible pour macOS, Windows et Linux, en anglais et en chinois.

[Site web](https://mq-studio.amigoer.com) |
[GitHub](https://github.com/amigoer/mq-studio) |
[Télécharger](https://github.com/amigoer/mq-studio/releases/latest)
