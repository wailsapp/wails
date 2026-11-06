---
title: "Redis Viewer"
description: "Wails で構築されたデスクトップ向け Redis GUI"
slug: "community/showcase/redisviewer"
sourcePath: "community/showcase/redisviewer.md"
---

![RedisViewer のスクリーンショット](/assets/showcase-images/redisviewer-overview1.webp)

![RedisViewer のスクリーンショット](/assets/showcase-images/redisviewer-overview2.webp)

[RedisViewer](https://redisviewer.com/) は、Wails で構築されたモダンなデスクトップ向け Redis GUI です。操作性を損なうことなく、複雑な値の調査、コマンドの実行、Redis のパフォーマンス分析を行えます。

Wails の WebView + Go アーキテクチャを中心に設計されており、大規模なキースペースやサイズの大きいペイロードをすべてフロントエンドへ送るのではなく、Go バックエンドに保持します。これにより、Go のメモリモデルを活用し、JS を多用するクライアントで一般的なヒープ負荷やメモリリークのリスクを回避します。UI は画面上に表示される内容だけをレンダリングするよう最適化されているため、膨大なデータセットを閲覧しても高い応答性が保たれ、ネイティブのデスクトップアプリに近い操作感を得られます。

[プロジェクトの Web サイトを見る](https://redisviewer.com/)
