---
title: "REST API"
description: "Wails v3 アプリケーションから RESTful API を公開する"
slug: "guides/patterns/rest-api"
sourcePath: "guides/patterns/rest-api.md"
---

@note{type="info"}
このページはプレースホルダーです。REST API に関する詳細な内容は近日公開予定です。

@end

標準の Go `net/http` ハンドラー、または [Gin](/guides/gin-routing/) などのフレームワークを使用して、Wails v3 アプリケーションから RESTful API を公開できます。Wails はアセットハンドラーを介してフロントエンドを配信するため、任意の HTTP ハンドラーをマウントし、UI やその他のクライアントからアクセスできます。
