---
title: "API REST"
description: "Exposer des API RESTful depuis votre application Wails v3"
slug: "guides/patterns/rest-api"
sourcePath: "guides/patterns/rest-api.md"
---

@note{type="info"}
Cette page est provisoire. La documentation complète sur les API REST sera bientôt disponible.

@end

Vous pouvez exposer des API RESTful depuis votre application Wails v3 à l’aide de gestionnaires Go `net/http` standard ou d’un framework tel que [Gin](/guides/gin-routing/). Comme Wails sert votre frontend par l’intermédiaire d’un gestionnaire de ressources, tout gestionnaire HTTP peut être monté et rendu accessible depuis votre interface utilisateur ou d’autres clients.
