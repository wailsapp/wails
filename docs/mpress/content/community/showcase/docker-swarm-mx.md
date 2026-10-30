---
title: "Docker Swarm MX"
description: "Windows desktop app for running a Docker Swarm from day to day"
slug: "community/showcase/docker-swarm-mx"
sourcePath: "community/showcase/docker-swarm-mx.md"
---

![Docker Swarm MX overview page with the nodes, services and load of a swarm](/assets/showcase-images/docker-swarm-mx.png)

**[Docker Swarm MX](https://docker-swarm-mx.eu)** is a **Windows desktop app for running a Docker Swarm** from day to
day, built with **Go, Wails v3 and Svelte 5**. It connects through your Docker CLI contexts and talks to the Docker
Engine API directly, so it installs nothing on the swarm.

## Key Highlights

- **One view across nodes** — load, disk, volumes and images of every node, each reached through its own Docker context
- **Services and stacks** — deploy, roll out a new image pinned to its digest, roll back, scale and rebalance; deploy
  stacks from compose files with a preview of what changes
- **Logs and events** — follow the logs of a service or the whole swarm across nodes, and keep the swarm's events
  longer than the 256 a manager remembers
- **A preview before risky changes** — drain shows which tasks move or stop, promote and demote warn about the quorum
- **Backups and notifications** — back up the swarm and its volumes to a folder, SFTP or S3, and get a message on
  Gotify, ntfy, Slack, Teams, Discord and more when a node goes down
- **A root shell on any node** — straight from the app, without a separate SSH session
- **Protect production** — per context normal, protected or read-only, with an audit log of every change the app makes

Wails keeps it a single native Windows app: the Go backend does the work against the Engine API and keeps watching
the swarm from the notification area, while the frontend stays a normal Svelte app.

Available for Windows, in English, Dutch, French, Spanish, Italian, Portuguese and Simplified Chinese. Open source
under the MIT License.

[Website](https://docker-swarm-mx.eu) |
[Codeberg](https://codeberg.org/MartijnSchimmel/docker-swarm-mx) |
[Download](https://codeberg.org/MartijnSchimmel/docker-swarm-mx/releases/latest)
