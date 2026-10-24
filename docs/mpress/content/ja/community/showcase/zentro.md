---
title: "Zentro"
description: "PostgreSQL、MySQL、SQL Server、SQLite 向けの、プロジェクト中心のデスクトップ SQL IDE"
slug: "community/showcase/zentro"
sourcePath: "community/showcase/zentro.md"
---

![スキーマ対応のオートコンプリート、スキーマエクスプローラー、レコードパネルを備えた Zentro の SQL エディター](/assets/showcase-images/zentro.png)

**[Zentro](https://zentro-ide.com)** は、PostgreSQL、MySQL、SQL Server、SQLite 向けの **プロジェクト中心のデスクトップ SQL IDE** で、 **Go、Wails v2、React**で構築されています。すべてのプロジェクトは自分がどの環境を指しているかを把握しており、取り消せない操作の前には Zentro が確認するため、ステージング向けのクエリが本番に流れることはありません。

## 主な特長

- **環境の認識** - プロジェクトが環境（開発、ステージング、本番）を持ち、本番を明確に示す表示、文のリスク分類、破壊的な書き込み前の影響行のプレビューを備えています
- **SQL エディター** - スキーマ対応のオートコンプリート、インライン診断、クイックフィックス
- **スキーマエクスプローラーと ERD** - テーブル、列、キー、インデックスを閲覧し、リレーションを可視化します
- **Git ベースの追跡** - 保存したスクリプトとスキーマのスナップショットをコードのようにバージョン管理し、履歴と差分を確認できます
- **SSH トンネルと内蔵ターミナル** - トンネル経由でデータベースに接続し、ローカルまたは SSH のターミナルをクエリの隣に置けます
- **自分のキーで使う AI アシスタント** - 設定したモデルのエンドポイントに接続し、Zentro が運用するサーバーは経由しません。外部エージェント向けの MCP サーバーも備えています
- **キーボード中心** - エクスプローラーと結果グリッドのフォーカス移動、コマンドパレット、連続ショートカット
- **デフォルトで非公開** - データベースの認証情報は OS のキーチェーンに保存されます

Zentro は Windows、macOS、Linux 向けに、プラットフォームごとに 1 つのネイティブバイナリとして提供されます。コアは GPL-3.0 のオープンソースで、任意の有料 Pro プランが本番向けの保護機能などを追加します（30 日間の無料トライアル付き）。Community Edition のソースコードは公開されています。

[ウェブサイト](https://zentro-ide.com) |
[GitHub](https://github.com/ZentroIDE/Zentro) |
[ダウンロード](https://zentro-ide.com/download)
