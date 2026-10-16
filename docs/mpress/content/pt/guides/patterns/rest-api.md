---
title: "APIs REST"
description: "Como expor APIs RESTful no seu aplicativo Wails v3"
slug: "guides/patterns/rest-api"
sourcePath: "guides/patterns/rest-api.md"
---

@note{type="info"}
Esta página é provisória. O conteúdo completo sobre APIs REST estará disponível em breve.

@end

Você pode expor APIs RESTful no seu aplicativo Wails v3 usando manipuladores padrão do Go `net/http` ou um framework como o [Gin](/guides/gin-routing/). Como o Wails disponibiliza o frontend por meio de um manipulador de ativos, qualquer manipulador HTTP pode ser montado e acessado pela sua interface ou por outros clientes.
