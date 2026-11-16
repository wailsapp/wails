---
title: "次のステップ"
description: "最初のアプリを構築した後の進め方"
slug: "quick-start/next-steps"
sourcePath: "quick-start/next-steps.md"
---

最初の Wails アプリを構築できました。以下の学習手順に沿って進めるか、作りたいものがすでに決まっている場合は該当するガイドに直接進んでください。

## さらに学ぶ

完全なアプリケーションを一つずつ作りながら、学習を続けましょう。

1. <strong>[QR コードサービスを作成する](/tutorials/01-creating-a-service/)。</strong>サービスを通じて Go のロジックをフロントエンドに公開する方法を練習します。
2. <strong>[TODO リストを作成する](/tutorials/02-todo-vanilla/)。</strong>サービス呼び出しを通じて、アプリケーションデータを追加、変更、削除します。
3. <strong>[メモアプリを作成する](/tutorials/03-notes-vanilla/)。</strong>ファイル、ネイティブダイアログ、保存データを扱います。

リリースについて学ぶ準備ができたら、[自己更新アプリのチュートリアル](/tutorials/04-self-update-a-wails-app/)に沿って進めてください。各レッスンの学習成果については、[すべてのチュートリアル](/tutorials/overview/)を参照してください。

## 機能

**すでに独自のアプリを構築していますか？** 目の前の作業に対応するハウツーガイドを利用してください。

| やりたいこと | 最初に読むページ |
| --- | --- |
| フロントエンドから Go を呼び出す | [JavaScript から Go を呼び出す](/features/bindings/methods/) |
| ウィンドウを開く、または管理する | [ウィンドウの基本](/features/windows/basics/) |
| ネイティブメニューまたはトレイアイコンを追加する | [アプリケーションメニュー](/features/menus/application/)または[システムトレイメニュー](/features/menus/systray/) |
| ユーザーにファイルを開く、または保存してもらう | [ファイルダイアログ](/features/dialogs/file/) |
| コンポーネント間でイベントを送信する | [イベントシステム](/features/events/system/) |
| テキストをコピーする、またはドロップされたファイルを受け取る | [クリップボード](/features/clipboard/basics/)または[ファイルのドラッグ＆ドロップ](/features/drag-and-drop/files/) |
| キーボードショートカットを追加する | [キーボードショートカット](/features/keyboard/shortcuts/) |
| アプリをビルドしてパッケージ化する | [アプリケーションのビルド](/guides/build/building/)を行った後、[Windows](/guides/build/windows/)、[macOS](/guides/build/macos/)、または[Linux](/guides/build/linux/)向けにパッケージ化する |

これらのページは、サイドバーの**ハウツーガイド**配下にあります。セクション全体を順番に読まなくても、該当するトピックを直接開けます。

## Wails の仕組みを理解する

コードの背景にある設計を理解したい場合は、**コア概念**を参照してください。

- [Wails の仕組み](/concepts/architecture/)では、Go バックエンド、フロントエンド、ネイティブ WebView のつながりを説明します。
- [Go とフロントエンドのブリッジ](/concepts/bridge/)では、Go と JavaScript の間で呼び出しがどのように行われるかを説明します。
- [アプリケーションのライフサイクル](/concepts/lifecycle/)では、起動と終了について説明します。
- [ビルドシステム](/concepts/build-system/)では、プロジェクトタスクとフロントエンドツールの役割について説明します。

## API を調べる

レッスンではなく、型、メソッド、オプション、コマンドを調べる必要がある場合は、[API リファレンス](/reference/overview/)を利用してください。コマンドラインでの作業については、まず[CLI コマンドの概要](/reference/cli/)を参照し、各コマンドの詳細には[CLI のフラグと例](/guides/cli/)を利用してください。

## ヘルプを得る

[よくある質問](/faq/)を確認するか、[動作するサンプル](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples)を参照するか、[Discord](https://discord.gg/JDdSxwjhGf)で質問してください。プラットフォーム固有の問題は、サイドバーの**トラブルシューティング**にまとめられています。
