---
title: "チュートリアル"
description: "アプリケーションを構築しながら Wails を学ぶ"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

実際に何かを構築しながら Wails を学びたい場合は、ガイド付きレッスンを選んでください。各レッスンでは、1 つのプロジェクトを進めて具体的な成果物を完成させます。

初めてアクセスした場合は、まず[インストール](/quick-start/installation/)と[最初のアプリ](/quick-start/first-app/)を完了してください。

## チュートリアルを選ぶ

最初のアプリを完成させたら、次の順序でデスクトップ向けレッスンを進めてください。

| レッスン | 学習内容 |
| --- | --- |
| [QR コードサービス](/tutorials/01-creating-a-service/) | Go のロジックをサービスとして公開し、フロントエンドから呼び出す |
| [TODO リスト](/tutorials/02-todo-vanilla/) | アプリケーションデータを作成、更新、削除する |
| [メモアプリ](/tutorials/03-notes-vanilla/) | ネイティブダイアログを使用し、データをファイルに保存する |
| [自己更新アプリ](/tutorials/04-self-update-a-wails-app/) | GitHub Releases から更新を確認して適用する |

任意のモバイル向けレッスンとして、[最初のモバイルアプリ](/guides/mobile/first-mobile-app/)を進めることもできます。モバイル対応は試験的な機能です。

既存のアプリで特定のタスクを完了する必要がある場合は、代わりに[ハウツーガイド](/quick-start/next-steps/#features)を選んでください。設計の背景については、[コアコンセプト](/concepts/architecture/)を参照してください。

@note{type="tip" title="Go を初めて使う場合"}
チュートリアルを始める前に、[Go Tour](https://go.dev/tour/)を完了してください。

@end

## QR コードサービス

![QR コードのサンプル](/assets/qr1.png)

QR コードジェネレーターを構築しながら、Wails サービスの基礎を学びます。このチュートリアルでは、アプリケーションロジックを再利用可能なサービスとして整理するためのコアコンセプトを紹介します。

**学習内容：**

- Wails サービスを作成して構成する方法
- 外部 Go 依存関係の管理
- Go メソッドをフロントエンドにバインドする方法
- Go と JavaScript の間でデータを受け渡す方法
- 保守性を高めるコード構成

<strong>対象：</strong>サービスアーキテクチャを理解したい、Wails を初めて使う方

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>始める</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### TODO リスト

![TODO リストアプリケーション](/assets/todo-app.png)

美しくモダンなインターフェースを備えた、完全な TODO リストアプリケーションを構築します。この実践的なチュートリアルでは、バニラ JavaScript を使用した実用的なアプリケーションを通じて、Wails の主要なパターンを学びます。

**学習内容：**

- スレッドセーフな状態管理を備えたサービスベースのアーキテクチャ
- CRUD 操作（作成、読み取り、更新、削除）
- Go と JavaScript 間の型安全なバインディング
- フレームワークによる複雑さを増やさずにモダンな UI を構築する方法
- 適切なエラー処理と検証のパターン

<strong>所要時間：</strong>約20分

<strong>対象：</strong>初めて完全な Wails アプリケーションを構築する方。フレームワークによる複雑さを加える前に、基礎を理解するのに最適です

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>始める</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### メモ

![メモアプリケーション](/assets/notes-app.png)

ネイティブのファイルダイアログと自動保存機能を備えた、Apple Notes 風のアプリケーションを構築します。このチュートリアルでは、ファイル操作、ネイティブダイアログ、プロフェッショナルな UI パターンなど、デスクトップ固有の機能を紹介します。

**学習内容：**

- ネイティブのファイルダイアログ（保存、開く、情報）
- JSON ベースのデータ永続化
- デバウンスを使用した自動保存パターン
- プロフェッショナルな 2 カラムのデスクトップレイアウト
- Go でファイルシステム操作を扱う方法

<strong>所要時間：</strong>約30分

<strong>対象：</strong>ファイル操作や OS ネイティブのダイアログなど、デスクトップ固有の機能を学びたい方

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>始める</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 自己更新する Wails アプリ

新しい`wails3 init`から始めて、署名付きリリースの検証とヘルパーモードによるバイナリの入れ替えを実装し、Wails アプリケーションにアプリ内自己更新機能を追加します。更新元には GitHub Releases を使用します。

**学習内容：**

- `app.Updater`を Wails アプリに組み込む方法
- GitHub Releases プロバイダーの設定
- ダイジェスト検証用の`SHA256SUMS`を使用してリリースを公開する方法
- 改ざん耐性を高める Ed25519 署名の追加
- CSS、カスタム HTML、または BYO を使用してデフォルトウィンドウをカスタマイズする方法
- `CheckInterval`を使用した定期的なバックグラウンド確認

<strong>所要時間：</strong>約25分

<strong>対象：</strong>更新可能なデスクトップアプリをリリースしたい方。API だけでなく、リリースパイプライン全体を扱います

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>始める</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
