---
title: "REST API"
description: "Wails v3 애플리케이션에서 RESTful API 제공하기"
slug: "guides/patterns/rest-api"
sourcePath: "guides/patterns/rest-api.md"
---

@note{type="info"}
이 페이지는 자리 표시자입니다. 전체 REST API 콘텐츠는 곧 제공될 예정입니다.

@end

표준 Go `net/http` 핸들러나 [Gin](/guides/gin-routing/) 같은 프레임워크를 사용하여 Wails v3 애플리케이션에서 RESTful API를 제공할 수 있습니다. Wails는 애셋 핸들러를 통해 프런트엔드를 제공하므로, 모든 HTTP 핸들러를 마운트하여 UI나 다른 클라이언트에서 접근할 수 있습니다.
