---
title: "チュートリアル"
description: "アプリケーションを構築しながら Wails を学ぶ"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

完全なアプリケーションを構築しながら Wails の概念を学べる、ステップ形式のチュートリアルです。各チュートリアルには、動作するコード、解説、実践的なパターンが含まれています。

@note{type="tip" title="Go は初めてですか？"}
チュートリアルを始める前に、[Go Tour](https://go.dev/tour/) を完了してください。

@end

## QR コードサービス

![QR コードの例](/assets/qr1.png)

QR コードジェネレーターを構築しながら、Wails サービスの基礎を学びます。このチュートリアルでは、アプリケーションロジックを再利用可能なサービスとして整理するための中核的な概念を紹介します。

**学習内容：**

- Wails サービスの作成方法と構成方法
- 外部 Go 依存関係の管理
- Go メソッドをフロントエンドにバインドする方法
- Go と JavaScript の間でデータを受け渡す方法
- 保守しやすいコードの整理方法

<strong>対象者：</strong>サービスアーキテクチャを理解したい、Wails を初めて使う方

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>始める</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### TODO リスト

![TODO リストアプリケーション](/assets/todo-app.png)

美しくモダンなインターフェースを備えた、完全な TODO リストアプリケーションを構築します。この実践的なチュートリアルでは、素の JavaScript を使用した実用的なアプリケーションを通して、Wails の中核的なパターンを学びます。

**学習内容：**

- スレッドセーフな状態管理を備えたサービスベースのアーキテクチャ
- CRUD 操作（作成、読み取り、更新、削除）
- Go と JavaScript 間の型安全なバインディング
- フレームワークの複雑さを伴わないモダンな UI の構築
- 適切なエラー処理と検証のパターン

<strong>所要時間：</strong>約20分

<strong>対象者：</strong>初めて完全な Wails アプリケーションを構築する方。フレームワークによる複雑さを加える前に、基礎を理解するのに最適です

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>始める</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### メモ

![メモアプリケーション](/assets/notes-app.png)

ネイティブのファイルダイアログと自動保存機能を備えた、Apple Notes 風のアプリケーションを構築します。このチュートリアルでは、ファイル操作、ネイティブダイアログ、プロフェッショナルな UI パターンなど、デスクトップ固有の機能を説明します。

**学習内容：**

- ネイティブのファイルダイアログ（保存、開く、情報）
- JSON ベースのデータ永続化
- デバウンスを使用した自動保存パターン
- プロフェッショナルな 2 カラムのデスクトップレイアウト
- Go でファイルシステム操作を扱う方法

<strong>所要時間：</strong>約30分

<strong>対象者：</strong>ファイル操作やネイティブ OS ダイアログなど、デスクトップ固有の機能を学びたい方

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>始める</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 自己更新対応 Wails アプリ

新規の `wails3 init` から始め、署名済みリリースの検証とヘルパーモードによるバイナリの入れ替えまで、Wails アプリケーションにアプリ内自己更新機能を追加します。更新元として GitHub Releases を使用します。

**学習内容：**

- `app.Updater` を Wails アプリに組み込む方法
- GitHub Releases プロバイダーの設定
- ダイジェスト検証用の `SHA256SUMS` を使用してリリースを公開する方法
- 改ざん耐性を確保するための Ed25519 署名の追加
- CSS、カスタム HTML、または BYO を使用したデフォルトウィンドウのカスタマイズ
- `CheckInterval` を使用した定期的なバックグラウンドチェック

<strong>所要時間：</strong>約25分

<strong>対象者：</strong>更新可能なデスクトップアプリをリリースする方。API だけでなく、リリースパイプライン全体を扱います

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>始める</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
