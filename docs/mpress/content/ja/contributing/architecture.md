---
title: "Wails v3 のアーキテクチャ"
description: "Wails v3 を構成するすべての要素を詳しく解説する図と説明"
slug: "contributing/architecture"
sourcePath: "contributing/architecture.md"
---

Wails v3 は、Go ランタイム、JavaScript ブリッジ、タスク駆動型ツールチェーン、各種テンプレートで構成される<strong>フルスタックのデスクトップフレームワーク</strong>です。 最新の Web 技術を活用したネイティブアプリケーションをリリースできます。

このページでは、次の4つの図で<em>全体像</em>を示します。

1. **全体アーキテクチャ** – 各サブシステムがどのように接続されているか\
2. **ランタイムフロー** – JS から Go を呼び出す場合と、その逆の場合に何が起こるか\
3. **開発環境と本番環境** – アセットサーバーの2つのモード\
4. **プラットフォーム別実装** – OS 固有のコードが配置されている場所\

---

## 1 · 全体アーキテクチャ

**Wails v3 – 上位レベルのスタック**

**［上位レベルのスタック図のプレースホルダー］**

---

## 2 · ランタイム呼び出しフロー

**ランタイム – JavaScript ⇄ Go の呼び出し経路**

**［ランタイム呼び出しフロー図のプレースホルダー］**

要点：

- **HTTP／IPC は不使用** – ブリッジはネイティブ WebView のインメモリチャネルを使用します\
- **メソッド ID** – 決定論的な FNV ハッシュにより、Go で O(1) の検索が可能です\
- **Promise** – エラーはスタックとコードを伴う reject として伝播します

---

## 3 · 開発環境と本番環境のアセットフロー

**開発環境 ↔ 本番環境のアセットサーバー**

**［アセットフロー図のプレースホルダー］**

- <strong>開発環境</strong>では、サーバーは未知のパスをフレームワークのライブリロードサーバーへプロキシし、静的アセットをディスクから配信します。
- <strong>本番環境</strong>では、同じ API が `go:embed` を基盤として動作し、依存関係のないバイナリを生成します。

---

## 4 · プラットフォーム固有ランタイムの分割

**OS 別ランタイムファイル**

**［プラットフォーム分割図のプレースホルダー］**

すべての機能は次のパターンに従います。

1. `pkg/application` 内の<strong>共通インターフェース</strong>\
2. `pkg/application/messageprocessor_*.go` 内の<strong>メッセージプロセッサ</strong>のエントリ\
3. `pkg/application/*_{darwin,linux,windows}.go` 内の<strong>OS 別実装</strong>（例：`webview_window_darwin.go`、`clipboard_linux.go`、`dialogs_windows.go`、`systemtray_*.go`、`mainthread_*.go`）。ビルドタグによって制御されます。Linux ではさらに、`linux_cgo.go`／`linux_cgo_gtk4.{go,c,h}` に cgo ブリッジがあります。

`internal/runtime/` に含まれるのは、小規模な`runtime{_darwin,_linux,_windows,_android,_dev,_prod}.go`ビルドタグ用グルーコードと、`internal/runtime/desktop/` 配下に埋め込まれた JS ランタイムだけです。

`internal/capabilities/` はプラットフォーム別の機能セットを宣言するために存在しますが、 `ErrCapability` センチネルはありません。機能の制御には通常のビルドタグと、 プラットフォーム固有のスタブ戻り値（`nil` や機能固有のエラーなど）を使用します。

---

## まとめ

これらの図は、**コードがどこにあるか**、**データがどのように移動するか**、 そして<strong>各レイヤーがどの責務を担うか</strong>を示しています。 後続の詳細ページを読み進める際には、これらの図を手元に置いてください。Wails v3 のソースツリーを案内する地図として役立ちます。
