---
title: "Foundation"
description: "本番利用に対応した Wails v3 デスクトップアプリのスキャフォールド"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![9 種類のレイアウトプリセットを切り替えられる Foundation のレイアウト設定画面](/assets/showcase-images/foundation-layouts.webp)

![ダークテーマの Foundation ホーム画面](/assets/showcase-images/foundation-home-dark.webp)

![Obsidian テーマで表示した Foundation の X-Pro チャート](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** は、**Go、React 19、TypeScript、Mantine 9** で構築された、本番利用に対応した **Wails v3 向けデスクトップアプリのスキャフォールド**です。階層化された Go バックエンド、プラットフォームごとのウィンドウ外観、永続化、テーマ、国際化、エンドツーエンドテストなど、通常は構築に数週間かかる部分をあらかじめ用意しているため、すぐにビジネスロジックの開発に取りかかれます。

## 主な特徴

- **各プラットフォームでネイティブなウィンドウ** — Windows と Linux ではフロントエンドが描画するフレームレスのタイトルバー、macOS ではネイティブの信号機ボタンを残した没入型タイトルバーを採用しています。
- **切り替え可能な 9 種類のレイアウト** — アイコンレール、サイドバー、トップナビゲーション、タイトルバー統合、右レール、下部 Dock、ドラッグ可能なフローティング Dock、インセットカード、フォーカスモードを、すべて CSS Grid のシェル上で提供します。独自のプリセットも 1 行で登録できます。
- **テーマと国際化** — ライト、ダーク、Obsidian のプリセットとカスタムパレットエディターを備え、英語と中国語に標準対応しています。メインウィンドウ、子ウィンドウ、トレイパネルの間でテーマと言語がリアルタイムに同期されます。
- **SQLite による永続化** — GORM の AutoMigrate によりマイグレーションファイルは不要です。データベースの保存場所はユーザーが選択でき、テーブルごとの使用量もグラフで確認できます。
- **デスクトップ統合** — システムトレイパネル、子ウィンドウ、ネイティブのファイル／メッセージダイアログに加え、ネイティブの「名前を付けて保存」ダイアログから CSV / XLSX をエクスポートできます。
- **堅牢なユーティリティ** — HTTP クライアント、Windows Job Object に紐づくサブプロセス、AES-GCM 暗号化、ローテーションするログ、アトミックな書き込みを備えています。
- **エンドツーエンドでテスト済み** — Playwright が Wails のサーバーモードビルド（実際の Go バックエンドと実際のフロントエンド）を操作し、同じテストスイートで README のスクリーンショットも生成します。

Wails により、すべてがプラットフォームごとに 1 つのバイナリにまとまります。Go のサービスは生成されたバインディングから直接呼び出され、フロントエンドは機能が増えても保守しやすい厳格な MVVM 構成に従います。

Windows、macOS、Linux に対応しています。

[GitHub](https://github.com/FxRayHughes/foundation) |
[リリース](https://github.com/FxRayHughes/foundation/releases/latest)
