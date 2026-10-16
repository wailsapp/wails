---
title: "REST API"
description: "Предоставление RESTful API из приложения Wails v3"
slug: "guides/patterns/rest-api"
sourcePath: "guides/patterns/rest-api.md"
---

@note{type="info"}
Эта страница пока не заполнена. Полная документация по REST API появится в ближайшее время.

@end

RESTful API из приложения Wails v3 можно предоставлять с помощью стандартных обработчиков Go `net/http` или фреймворка, например [Gin](/guides/gin-routing/). Поскольку Wails обслуживает фронтенд через обработчик ресурсов, можно подключить любой обработчик HTTP и обращаться к нему из пользовательского интерфейса или других клиентов.
