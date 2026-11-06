---
title: "Foundation"
description: "Wails v3.1、React 19、Mantine 9 によるデスクトップアプリのスキャフォールド"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![9 種類のレイアウトプリセットを切り替えられる Foundation のレイアウト設定画面](/assets/showcase-images/foundation-layouts.webp)

![ダークテーマの Foundation ホーム画面](/assets/showcase-images/foundation-home-dark.webp)

![Obsidian テーマで表示した Foundation の X-Pro チャート](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** は、**Go、React 19、TypeScript、Mantine 9** で構築された **Wails v3.1 向けデスクトップアプリのスキャフォールド**です。階層化された Go バックエンド、プラットフォームごとのウィンドウ外観、SQLite による永続化、テーマ、国際化、エンドツーエンドテストを、新しいプロジェクトの出発点として提供します。

## 機能

- **プラットフォームごとのウィンドウ** — Windows と Linux ではフレームレスのウィンドウとフロントエンドが描画するタイトルバー、macOS ではネイティブの信号機ボタンを残した没入型タイトルバーを使用します。
- **9 種類のレイアウトプリセット** — アイコンレール、サイドバー、トップナビゲーション、タイトルバー統合、右レール、下部 Dock、ドラッグ可能なフローティング Dock、インセットカード、フォーカスモードを CSS Grid のシェル上で提供します。独自のプリセットも登録できます。
- **テーマと国際化** — ライト、ダーク、Obsidian のプリセットとカスタムパレットエディターを備え、UI は英語と中国語に対応しています。メインウィンドウ、子ウィンドウ、トレイパネルの間でテーマと言語が同期されます。
- **SQLite による永続化** — GORM の AutoMigrate によりマイグレーションファイルは不要です。データベースの保存場所はユーザーが選択でき、テーブルごとの使用量をグラフで確認できます。
- **デスクトップ統合** — システムトレイパネル、子ウィンドウ、ネイティブのファイル／メッセージダイアログに加え、ネイティブの「名前を付けて保存」ダイアログから CSV / XLSX をエクスポートできます。
- **ユーティリティ** — HTTP クライアント、Windows Job Object に紐づくサブプロセス、AES-GCM 暗号化、ローテーションするログ、アトミックな書き込みを備えています。
- **エンドツーエンドテスト** — Playwright が Wails のサーバーモードビルド（実際の Go バックエンドとフロントエンド）に対して実行され、同じテストスイートで README のスクリーンショットも生成します。

フロントエンドは Wails が生成したバインディングを通じて Go のサービスを呼び出し、各プラットフォーム向けに単一のバイナリとしてビルドされます。フロントエンドは MVVM 構成を採用しており、README にはスキャフォールドを自分のプロジェクト名に変更する手順が記載されています。

Windows、macOS、Linux で動作します。

[GitHub](https://github.com/FxRayHughes/foundation) |
[リリース](https://github.com/FxRayHughes/foundation/releases/latest)
