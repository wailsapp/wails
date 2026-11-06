---
title: "CFN Tracker"
description: "Wailsで構築されたデスクトップアプリケーション"
slug: "community/showcase/cfntracker"
sourcePath: "community/showcase/cfntracker.md"
---

![CFN Tracker](/assets/showcase-images/cfntracker.webp)

[CFN Tracker](https://github.com/williamsjokvist/cfn-tracker) - Street Fighter 6またはVの任意のCFNプロフィールについて、進行中の対戦を追跡できます。利用を開始するには、[Webサイト](https://cfn.williamsjokvist.se/)をご確認ください。

## 機能

- リアルタイムの対戦追跡
- 対戦ログと統計情報の保存
- ブラウザソースを介したOBSへのリアルタイム統計情報の表示に対応
- SF6とSFVの両方に対応
- ユーザーがCSSを使用して独自のOBSブラウザテーマを作成可能

### Wailsと併用している主な技術

- [Task](https://github.com/go-task/task) - 一般的なコマンドを簡単に使用できるようにWails CLIをラップ
- [React](https://github.com/facebook/react) - 豊富なエコシステム（radix、framer-motion）を理由に採用
- [Bun](https://github.com/oven-sh/bun) - 依存関係の解決とビルドが高速なため使用
- [Rod](https://github.com/go-rod/rod) - 認証と変更のポーリングに使用するヘッドレスブラウザ自動化ツール
- [SQLite](https://github.com/mattn/go-sqlite3) - 対戦、セッション、プロフィールの保存に使用
- [Server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events) - 追跡の更新をOBSブラウザソースに送信するためのHTTPストリーム
- [i18next](https://github.com/i18next/) - Goレイヤーからローカライズオブジェクトを提供するバックエンドコネクターと併用
- [xstate](https://github.com/statelyai/xstate) - 認証プロセスと追跡に使用するステートマシン
