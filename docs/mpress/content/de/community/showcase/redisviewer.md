---
title: "Redis-Viewer"
description: "Eine mit Wails entwickelte Redis-GUI für den Desktop"
slug: "community/showcase/redisviewer"
sourcePath: "community/showcase/redisviewer.md"
---

![Screenshot von RedisViewer](/assets/showcase-images/redisviewer-overview1.webp)

![Screenshot von RedisViewer](/assets/showcase-images/redisviewer-overview2.webp)

[RedisViewer](https://redisviewer.com/) ist eine moderne, mit Wails entwickelte Redis-GUI für den Desktop. Damit können Sie komplexe Werte untersuchen, Befehle ausführen und die Redis-Performance analysieren, ohne Abstriche bei der Interaktionsqualität zu machen.

Die Anwendung basiert auf der WebView- und Go-Architektur von Wails und hält große Schlüsselräume und umfangreiche Nutzdaten im Go-Backend, statt alles an das Frontend zu übertragen. So nutzt sie das Speichermodell von Go und vermeidet den Heap-Druck und die Risiken von Speicherlecks, die bei stark auf JS basierenden Clients häufig auftreten. Die Benutzeroberfläche ist darauf optimiert, nur die aktuell sichtbaren Inhalte zu rendern. Dadurch bleibt das Durchsuchen riesiger Datensätze reaktionsschnell und fühlt sich nahezu wie bei einer nativen Desktopanwendung an.

[Projektwebsite besuchen](https://redisviewer.com/)
