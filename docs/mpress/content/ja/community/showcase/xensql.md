---
title: "XenSQL"
description: "複数データベースに対応し、高度なクエリツールを備えたローカルファーストのSQLワークベンチ"
slug: "community/showcase/xensql"
sourcePath: "community/showcase/xensql.md"
---

![インライン編集に対応したXenSQLテーブルエディター](/assets/showcase-images/xensql-1.png) ![テーブル名と列名の候補を提示するXenSQLクエリエディター](/assets/showcase-images/xensql-2.png) ![XenSQLでの複数トランザクションクエリ](/assets/showcase-images/xensql-3.png)

<strong>[XenSQL](https://github.com/Bare7a/XenSQL)</strong>は、<strong>Go、Wails、React</strong>で構築された<strong>高速なローカルファーストのSQLデスクトップワークベンチ</strong>です。PostgreSQL、MySQL/MariaDB、SQLiteを、すっきりとしたネイティブアプリのような単一のインターフェースに統合しています。クラウドの利用も、テレメトリも、アカウントも一切ありません。

## 主な特長

- **高機能なSQLエディター** — Monacoをベースに、スキーマを認識するスマートなオートコンプリート、複数ステートメントの実行、結果のストリーミング、ステートメントごとの結果タブを備えています
- **高度なデータビューアー** — インタラクティブなJSONインスペクター、構文を認識するセルエディター（JSON、XML、HTML、テキスト）、インライン編集、レコード全体の詳細表示を備えています
- **シームレスなデータ編集** — テーブルの閲覧、変更のインラインでのステージング、一括操作、および`RETURNING`対応の安全な`INSERT`/`UPDATE`/`DELETE`を利用できます
- **生産性向上機能** — スキーマエクスプローラー、保存済みクエリ、クエリ履歴、クイック検索（`Ctrl+P`）、キーボード中心のワークフローを備えています
- **エクスポート形式** — CSV、JSON、Markdown、SQL INSERT文

完全にオフラインで動作し、持ち運びも可能です。すべてのデータは、アプリと一緒に移動できる単一の`XenSQL-data/`フォルダーにローカル保存されます。

**対応データベース**：PostgreSQL、MySQL、MariaDB、SQLite（読み取り専用モードとセキュアな通信オプションに対応）。

従来のSQLツールにありがちな肥大化を避け、速度、明快さ、制御性を求める開発者向けに設計されています。
