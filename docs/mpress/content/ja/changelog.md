---
title: "変更履歴"
description: "Wails v3のバージョン履歴とリリースノート"
slug: "changelog"
sourcePath: "changelog.md"
---

凡例：

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- このプロジェクトの特筆すべき変更は、すべてこのファイルに記録されます。

この形式は[Keep a Changelog](https://keepachangelog.com/en/1.0.0/)に基づいており、 このプロジェクトは[セマンティック バージョニング](https://semver.org/spec/v2.0.0.html)に準拠しています。

- 新機能には`Added`。
- 既存機能の変更には`Changed`。
- 近く削除される機能には`Deprecated`。
- 削除済みの機能には`Removed`。
- バグ修正には`Fixed`。
- 脆弱性がある場合には`Security`。

_/

/_   * このファイルを更新しないでください *   更新内容は`v3/UNRELEASED_CHANGELOG.md`に追加してください   よろしくお願いします！ _/

## [未リリース]

## v3.0.0-beta.21 - 2026-09-13

## 追加

- [PR](https://github.com/wailsapp/wails/pull/6116) で、Wails の v3 ドキュメントを M-Press で配信（@leaanthony）

## 修正

- [PR](https://github.com/wailsapp/wails/pull/6118) で、変更履歴の生成時に MPD フロントマターから JSON の slug 値を解析するように修正（@leaanthony）
- [PR](https://github.com/wailsapp/wails/pull/6080) で、アップデーターがヘルパーの環境変数を消去し、バックアップに失敗した後に元のターゲットを再起動するように修正（@cnmax）
- [PR](https://github.com/wailsapp/wails/pull/6098) で、App.Run の実行中にデフォルトのシグナルハンドラーを起動するように修正（@leaanthony）
- [PR](https://github.com/wailsapp/wails/pull/6112) で、Windows メニューが nil メニューを処理し、置き換えられたリソースを解放して、メニューバーを再描画するように修正（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6115) で、共有 YAML 設定を使用する新規プロジェクトの MSIX パッケージ化を復元（@leaanthony）
- ジェネリックモデルのクリエーターが後続のヘルパー宣言を参照すると、生成された JavaScript および TypeScript バインディングの読み込みに失敗する問題を修正し、相互依存するジェネリックモデルの作成時にスタックオーバーフローが発生しないように修正（#6062）

## v3.0.0-beta.20 - 2026-09-10

## 変更

- Claveのショーケースリンクを現在のWebサイトとリポジトリに更新（@01xR4inによる[PR](https://github.com/wailsapp/wails/pull/6082)）

## 修正

- Windowsで中止されたアセットリクエスト（ワーカーからのリクエストを含む）をキャンセルしつつ、ナビゲーションをまたいでkeepaliveハンドラーを維持。Appleプラットフォームでは、ネイティブリクエストのコンテキストをアプリケーションラッパー経由で転送。（#5963、#5969）
- 競合するプッシュが発生しても、再試行によって変更履歴のエントリを保持（@leaanthonyによる[PR](https://github.com/wailsapp/wails/pull/6094)）
- モジュールに一度も同梱されていないバイナリを参照する埋め込みを削除し、すべてのプラットフォームで`pattern arm64/WebView2Loader.dll: no matching files found`により`go mod vendor`が失敗する問題を修正。これにより、[#5782](https://github.com/wailsapp/wails/issues/5782)および[#5376](https://github.com/wailsapp/wails/issues/5376)も修正（@Grantmartin2002による[PR](https://github.com/wailsapp/wails/pull/6031)）

## 削除

- 純Goローダーに置き換えられたネイティブWebView2ローダーのサポートを削除。これにより、埋め込まれていた`WebView2Loader.dll`バイナリと`github.com/jchv/go-winloader`依存関係を廃止。`native_webview2loader`ビルドタグは引き続き受け付けられ、エラーも発生しなくなりましたが、v3ビルドには影響しません（@Grantmartin2002による[PR](https://github.com/wailsapp/wails/pull/6031)）。
- macOS APIガイドから未使用のビルドタグとFPSオプションを削除（@leaanthonyによる[PR](https://github.com/wailsapp/wails/pull/6097)）

## v3.0.0-beta.19 - 2026-09-09

## 追加

- オプトインで使用できるよう、macOSのプライベートAPIをビルドタグで制限。[ドキュメント](https://v3.wails.io/features/browser/integration)、[ドキュメント](https://v3.wails.io/features/environment/info)、[ドキュメント](https://v3.wails.io/features/windows/basics)、[ドキュメント](https://v3.wails.io/features/windows/frameless)、[ドキュメント](https://v3.wails.io/features/windows/notch-windows)、[ドキュメント](https://v3.wails.io/features/windows/options)、[ドキュメント](https://v3.wails.io/guides/build/macos)、[ドキュメント](https://v3.wails.io/guides/build/private-macos-apis)、[ドキュメント](https://v3.wails.io/reference/overview)を参照（@leaanthonyによる[PR](https://github.com/wailsapp/wails/pull/6087)）。

## 修正

- 64 MiBを超えるランタイムリクエストをHTTP 413で拒否（@leaanthonyによる[PR](https://github.com/wailsapp/wails/pull/6091)）

## セキュリティ

- トークン認証によりMCPのオリジンとリモートアクセスを強化（@leaanthonyによる[PR](https://github.com/wailsapp/wails/pull/6092)）

## v3.0.0-beta.18 - 2026-09-08

## 修正

- ポインタレシーバーを使用することで、LinuxおよびDarwinでのCallocのメモリリークを修正（@4RH1T3CT0R7による[PR](https://github.com/wailsapp/wails/pull/6083)）

## v3.0.0-beta.17 - 2026-09-06

## 修正

- Windows：WebResourceRequestedハンドラー内の`GetRequest`が失敗するかnilでも、プロセスが終了（`log.Fatal`／nil参照パニック）しないように変更。代わりにリクエストを破棄してログに記録（@midagedevによる[PR](https://github.com/wailsapp/wails/pull/6006)）。

## v3.0.0-beta.16 - 2026-08-29

## 変更

- 新しいターミナルウィンドウで公証用パスワードの入力を要求（@leaanthonyによる[PR](https://github.com/wailsapp/wails/pull/6029)）

## 修正

- macOSでシステムトレイのクリック種別を正しく処理（@ChewbaccaCookieによる[PR](https://github.com/wailsapp/wails/pull/5919)）
- CIで更新前に未使用のMicrosoft aptリポジトリを削除（@Grantmartin2002による[PR](https://github.com/wailsapp/wails/pull/6041)）

## v3.0.0-beta.15 - 2026-08-27

## 修正

- WebView2の埋め込みタイムアウトを60秒に延長（@Grantmartin2002による[PR](https://github.com/wailsapp/wails/pull/6043)）

## v3.0.0-beta.14 - 2026-08-26

## 修正

- macOSでControlキーと文字キーの組み合わせによる押下に正しい名前を付与（@taliesin-aiによる[PR](https://github.com/wailsapp/wails/pull/6032)）
- WindowsでICO形式のトレイアイコンを修正し、タスクバーのテーマに追従するよう変更（@nik9playによる[PR](https://github.com/wailsapp/wails/pull/6016)）

## v3.0.0-beta.13 - 2026-08-25

## 修正

- macOSでモーダルループの実行中もメインスレッドの処理を継続（@leaanthonyによる[PR](https://github.com/wailsapp/wails/pull/6026)）
- モバイルのセキュアストレージが失敗をエラーとして返せるようにし、失敗時はアクセスを拒否するよう変更（@mortenolsrudによる[PR](https://github.com/wailsapp/wails/pull/5923)）
- リスナーが登録されていない場合でも、アプリケーションのイベントフックを実行（@archy-rock3t-cloudによる[PR](https://github.com/wailsapp/wails/pull/5999)）
- コメントとローカライズ済みドキュメントの誤字を修正（@haoku123による[PR](https://github.com/wailsapp/wails/pull/6023)）
- `v3/examples`以下にコミットされていたコンパイル済みmacOSバイナリを削除（@4RH1T3CT0R7による[PR](https://github.com/wailsapp/wails/pull/6025)）

## v3.0.0-beta.12 - 2026-08-21

## 追加

- ライフサイクルとテレメトリの例を備えたmacOSノッチ通知ウィンドウを追加。[ドキュメント](https://v3.wails.io/features/windows/notch-windows)を参照（@leaanthonyによる[PR](https://github.com/wailsapp/wails/pull/6010)）。
- 新しいオプションとネイティブ統合を備えた macOS NSPanel ウィンドウのサポートを追加 — @leaanthony による [PR](https://github.com/wailsapp/wails/pull/6008)（[ドキュメント](https://v3.wails.io/features/windows/options)を参照）

## 修正

- 並行呼び出し時に SQLite Prepare がハングする問題を防止 — @archy-rock3t-cloud による [PR](https://github.com/wailsapp/wails/pull/5998)

## v3.0.0-beta.11 - 2026-08-20

## 削除

- ドキュメントから廃止済みの実装トラッカーを削除 — @leaanthony による [PR](https://github.com/wailsapp/wails/pull/6005)

## v3.0.0-beta.10 - 2026-08-19

## 修正

- GTK4 Linux ホストがカスタムプロトコルおよびファイル関連付けによる起動時の引数を破棄する問題を修正 — @midagedev による [PR](https://github.com/wailsapp/wails/pull/6000)
- 変更履歴の検証で、削除された行と同一ソース内の修正を正しく処理 — @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5993)

## v3.0.0-beta.9 - 2026-08-16

## 追加

- エージェント支援によるプロジェクト管理に対応する、安全な wails3 mcp サーバーを追加 — @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5896)
- バインディング内のモデルに関するドキュメントを追加 — @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5988)（[ドキュメント](https://v3.wails.io/features/bindings/models)を参照）
- Atomic Linux システムでの rpm-ostree によるインストールをサポート — @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5987)
- 週次のスター履歴チャートをネイティブに生成して公開する機能を追加 — @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5986)（[ドキュメント](https://v3.wails.io/credits)、[ドキュメント](https://v3.wails.io/de/credits)、[ドキュメント](https://v3.wails.io/fr/credits)、[ドキュメント](https://v3.wails.io/id/credits)、[ドキュメント](https://v3.wails.io/ja/credits)、[ドキュメント](https://v3.wails.io/ko/credits)、[ドキュメント](https://v3.wails.io/pt/credits)、[ドキュメント](https://v3.wails.io/ru/credits)、[ドキュメント](https://v3.wails.io/zh-cn/credits)、[ドキュメント](https://v3.wails.io/zh-tw/credits)を参照）
- アプリケーションバンドルのリソースを解決する Darwin 専用の mac パッケージを追加 — @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5965)（[ドキュメント](https://v3.wails.io/guides/build/macos)を参照）
- Condui のショーケースページとインデックス項目を追加 — @mgueregath による [PR](https://github.com/wailsapp/wails/pull/5962)（[ドキュメント](https://v3.wails.io/community/showcase/condui)および[ドキュメント](https://v3.wails.io/community/showcase)を参照）
- スクリーンショットとプロジェクトへのリンクを掲載した Redis Viewer のショーケースページを追加 — @redisviewer による [PR](https://github.com/wailsapp/wails/pull/5984)（[ドキュメント](https://v3.wails.io/community/showcase)および[ドキュメント](https://v3.wails.io/community/showcase/redisviewer)を参照）

## 変更

- Linux の GTK アプリケーションフラグを G<em>APPLICATION</em>NON_UNIQUE に更新 — @overlordtm による [PR](https://github.com/wailsapp/wails/pull/5971)
- 存在しないウィンドウのイベントを警告レベルではなくデバッグレベルでログに記録 — @julianstorer による [PR](https://github.com/wailsapp/wails/pull/5914)

## 修正

- 登録済みの macOS キーボードショートカットが webview より優先されるように修正 — @julianstorer による [PR](https://github.com/wailsapp/wails/pull/5902)
- ドキュメントのサイドバーにある壊れたリンクを修復 — @northes による [PR](https://github.com/wailsapp/wails/pull/5937)
- WebKit が対応するカスタムスキームのタスクを中止した際に、macOS および iOS のアセットリクエストコンテキストをキャンセル（#5963）
- WindowSetFullscreenButtonEnabled メッセージを処理 — @archy-rock3t-cloud による [PR](https://github.com/wailsapp/wails/pull/5976)
- preact-ts テンプレートに Fragment をインポートし、ビルドエラーを解消 — @haoku123 による [PR](https://github.com/wailsapp/wails/pull/5979)
- アクティブなウィンドウまたはディスプレイが利用可能になる前に画面検出が実行された場合に、従来の GTK3 サービス専用アプリケーションがクラッシュする問題を防止（#5966）
- 未リリースの変更履歴が空でも、バージョンを明示したリリース処理を続行できるように修正（#5977）

## セキュリティ

- セキュリティアドバイザリに対処するため、Web サイトの nanoid ロックファイルを修正済みの 3.3.18 に更新 — @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5985)

## v3.0.0-beta.8 - 2026-08-12

## 追加

- 変更履歴の自動生成項目にドキュメント URL の生成機能を追加 — @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5957)
- Streams を追加：WebSocket のプログラミングモデルを使用し、リッスンソケットを必要としない、Go と JavaScript 間の双方向バイトストリーム。Go では `app.HandleStream(name, handler)` を使用してストリームを宣言し、フロントエンドから `Stream(name)` を使用して接続します。これは `WebSocket` と同じ構造のオブジェクトを返します。Go→JS の通信はアセットサーバー経由でウィンドウごとに保持される単一のポーリングにより行われ、JS→Go の通信には通常の POST が使用されます。TCP ポートへのバインドは一切行われず、`evaluateJavaScript` を経由する通信もありません。サーバービルド（`-tags server`）では、代わりに同じハンドラーが実際の WebSocket 経由で提供されるため、ビルドが異なってもアプリケーションコードは同一です。@leaanthony
- mailbox の変更履歴項目を「未リリース」に移動 — @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5935)

## 変更

- ドキュメントのサイドバー自動生成とブログ著者型の導出を更新 — @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5938)

## 修正

- WebView2 の初期化で期限とメッセージポンプを使用 — @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5952)
- WebView2 の Cookie テストを、明示的に有効化しない限り CI ではスキップし、実行を現在の OS スレッドに固定 — @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5951)
- Windows のメニュービルダーで、サブメニューの親項目のコマンド ID を復元 — @gilad-ch による [PR](https://github.com/wailsapp/wails/pull/5944)
- 公式のクロスコンパイル用イメージを GTK 4.14+ Linux のサポート基準に合わせる（#5928）
- 継承されたリンカーフラグを維持し、-ObjC を追加するように iOS Xcode プロジェクトを構成 — @mortenolsrud による [PR](https://github.com/wailsapp/wails/pull/5915)
- 大規模なフロントエンドで `wails3 dev` アセットプロキシの TCP 接続が過剰に切り替わる問題を修正。この問題により、ホストのエフェメラルポートが枯渇し、無関係なプロセスで `EADDRNOTAVAIL` エラーが発生する可能性がありました
- 順序どおりのディスパッチとバックプレッシャーのため、ウィンドウごとのイベント JavaScript をキューに格納 — @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5934)

## 削除

- デスクトップバイナリのリリースパイプラインを削除。v3 のリリースはタグのみとなり、`wails3` CLI は `go install` を使用してインストールします。`release-v3.yml` と、それを起動していた nightly ステップを削除 — @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5946)

## v3.0.0-beta.7 - 2026-08-11

## 追加

- [PR](https://github.com/wailsapp/wails/pull/5512)（@Eyalm321）で、メディア再生にユーザー操作を不要とするmacOSの自動再生設定を追加
- [PR](https://github.com/wailsapp/wails/pull/5935)（@leaanthony）で、メールボックスの変更履歴エントリを「未リリース」に移動

## 変更

- [PR](https://github.com/wailsapp/wails/pull/5945)（@savely-krasovsky）で、macOSのズームアニメーションにCADisplayLinkまたはNSTimerを使用し、動作をより滑らかに変更

## 修正

- [PR](https://github.com/wailsapp/wails/pull/5915)（@mortenolsrud）で、継承されたリンカーフラグを保持し、-ObjCを追加するようにiOSのXcodeプロジェクトを構成
- 大規模なフロントエンドで`wails3 dev`アセットプロキシがTCP接続を過剰に切り替える問題を修正。この問題により、ホストのエフェメラルポートが枯渇し、無関係なプロセスが`EADDRNOTAVAIL`で失敗する可能性があった
- [PR](https://github.com/wailsapp/wails/pull/5934)（@leaanthony）で、順序どおりのディスパッチとバックプレッシャーのため、ウィンドウごとのイベントJavaScriptをキューに格納

### 追加

- [PR](https://github.com/wailsapp/wails/pull/5851)（@savely-krasovsky、@DevLumuz）で、イベントを順序どおりに配信するための汎用非同期FIFOメールボックスを実装

## v3.0.0-beta.6 - 2026-08-09

## 追加

- [PR](https://github.com/wailsapp/wails/pull/5930)（@leaanthony）で、サイズ超過イベント向けの容量制限付きホスト側ストレージと、順序どおりのJavaScript配信を実装
- [PR](https://github.com/wailsapp/wails/pull/5921)（@julianstorer）で、ウィンドウの点滅機能をmacOSではDockアイコンのバウンドで実現

## 修正

- [PR](https://github.com/wailsapp/wails/pull/5931)（@leaanthony）で、アセットサーバーがフラッシュ時にContent-Typeスニッファーのエラーと未書き込みのプレフィックスを保持するように修正
- Wailsのコールバックからアプリケーションメニューを置き換えたときに、macOSアプリケーションがクラッシュする問題を修正
- Windows 10 1809／Windows Server 2019（ビルド17763）でネイティブメニューが読めなくなる問題を修正。ダークモード用のuxthemeエクスポートはビルド18334を条件としていたため、これらのホストではアプリレベルのダークモードのオプトインが実行されていなかった。その結果、メニューの背景は暗く描画された一方、Windowsはライトテーマのメニューテキストを描画し続け、暗い背景に暗い文字が表示されていた。序数は17763以降に存在するため、条件をそれに合わせて修正
- `w32.GetStockObject`が`GetStockObject`ではなく`GetDeviceCaps`を呼び出し、すべてのストックオブジェクトに対して0を返していた問題を修正
- [PR](https://github.com/wailsapp/wails/pull/5924)（@jannskiee）で、WebView2ブートストラッパーのダウンロードエラー処理と報告を改善

## v3.0.0-beta.5 - 2026-08-07

## 修正

- [PR](https://github.com/wailsapp/wails/pull/5897)（@julianstorer）で、macOSアプリのアクティベーションが通常のアプリに限ってアクティベーションポリシーに従うように修正
- [PR](https://github.com/wailsapp/wails/pull/5898)（@julianstorer）で、Linuxビルドにおける未初期化のGTKウィンドウを保護
- [PR](https://github.com/wailsapp/wails/pull/5899)（@julianstorer）で、URLを読み込む前にLinuxのWebKitウィンドウへ明示的な不透明の背景色を設定

## v3.0.0-beta.4 - 2026-08-05

## 変更

- [PR](https://github.com/wailsapp/wails/pull/5890)（@mortenolsrud）で、Androidのビルドタスクのデフォルトをarm64に変更し、deploy-emulatorがホストアーキテクチャを選択するように変更

## 修正

- [PR](https://github.com/wailsapp/wails/pull/5900)（@leaanthony）で、ドラッグ中にmacOSウィンドウのズーム状態を保持し、動きを抑制
- `webview_window_windows_nonclient.go`のビルド制約に`!server`を追加し、Windowsのサーバーモードビルドを修正

## v3.0.0-beta.3 - 2026-08-03

## 追加

- [PR](https://github.com/wailsapp/wails/pull/5881)（@leaanthony）で、フェーズ10のベータ検証完了を実装詳細に記載

## 修正

- [PR](https://github.com/wailsapp/wails/pull/5877)（@leaanthony）で、WindowsのダークモードAPIにウィンドウハンドルを渡し、引数を検証
- [PR](https://github.com/wailsapp/wails/pull/5870)（@taliesin-ai）で、macOSのフレームレスウィンドウにおけるタイトルバーボタンの状態解決を一元化
- Windowsアプリケーションがダークモードを要求している一方で、Windowsのアプリテーマがライトの場合に、ネイティブメニューの文字が読めなくなる問題を修正。Windowsがダークモード用のメニューテキストを描画できるようになるまでは、メニューに対応する明るいネイティブ背景を使用するように変更
- Windows 10 1809／Windows Server 2019（ビルド17763）でネイティブメニューが読めなくなる問題を修正。ダークモード用のuxthemeエクスポートはビルド18334を条件としていたため、これらのホストではアプリレベルのダークモードのオプトインが実行されていなかった。その結果、メニューの背景は暗く描画された一方、Windowsはライトテーマのメニューテキストを描画し続け、暗い背景に暗い文字が表示されていた。序数は17763以降に存在するため、条件をそれに合わせて修正

## v3.0.0-beta.2 - 2026-08-02

## 変更

- v3をアルファ版からベータ版に昇格
- システムトレイのスマートデフォルトとポップアップの自動非表示動作を文書化し、クリックハンドラーの選択に対する回帰テストを追加（#5840）
- [PR](https://github.com/wailsapp/wails/pull/5861)（@leaanthony）で、GitHubアップデーターがデフォルトでWindowsインストーラーのアセットを除外するように変更
- [PR](https://github.com/wailsapp/wails/pull/5866)（@leaanthony）で、macOSのフレームレスウィンドウに、角丸、角型、カスタム半径のコーナーを追加

## 修正

- GTK4ウィンドウの現在のサイズを報告し、構成されたサーフェスからサイズ変更、最大化、最小化、フルスクリーンの状態イベントを発行するように修正（#5830）
- [PR](https://github.com/wailsapp/wails/pull/5854)（@taliesin-ai）で、fetchリクエストにBlobまたはFormDataを指定して送信するとLinuxのWebKitがクラッシュする問題を修正
- [PR](https://github.com/wailsapp/wails/pull/5865)（@leaanthony）で、Blob/FormDataヘッダーがない場合にfetch shimがundefinedを渡すように修正

## v3.0.0-alpha2.122 - 2026-08-01

## 追加

## 変更

- [PR](https://github.com/wailsapp/wails/pull/5866)（@leaanthony）で、macOSのフレームレスウィンドウに、角丸、角型、カスタム半径のコーナーを追加

## 修正

- [PR](https://github.com/wailsapp/wails/pull/5865)（@leaanthony）で、Blob/FormDataヘッダーがない場合にfetch shimがundefinedを渡すように修正

## v3.0.0-alpha2.121 - 2026-07-31

## 追加

- macOS向けDMGパッケージ作成のサポートを、新しいオプションおよびビルドタスクとともに追加（@leaanthonyによる[PR](https://github.com/wailsapp/wails/pull/5857)）

## 変更

- GitHubアップデーターがデフォルトでWindowsインストーラーのアセットを除外するように変更（@leaanthonyによる[PR](https://github.com/wailsapp/wails/pull/5861)）

## 修正

- fetchリクエストでBlobまたはFormDataを送信した際にLinuxのWebKitがクラッシュする問題を修正（@taliesin-aiによる[PR](https://github.com/wailsapp/wails/pull/5854)）

## v3.0.0-alpha2.120 - 2026-07-31

## 追加

- macOSのタイトルバーをダブルクリックした際に最大化または最小化するアクションを実装（@taliesin-aiによる[PR](https://github.com/wailsapp/wails/pull/5853)）

## 変更

- QRサービスのチュートリアルをNewServiceWithOptionsを使用するよう更新し、スペースを追加（@jeongkyuによる[PR](https://github.com/wailsapp/wails/pull/5849)）

## 修正

- macOSでのズーム中もWKWebViewが応答し続けるよう修正（@leaanthonyによる[PR](https://github.com/wailsapp/wails/pull/5856)）
- GTK4のウィンドウサイズ取得を修正し、設定された`GdkSurface`から、サイズ変更、最大化、最小化、およびフルスクリーンの状態イベントを発行するよう修正。

## v3.0.0-alpha2.119 - 2026-07-27

## 修正

- 複数言語のドキュメントを更新し、アーキテクチャ図を追加（@taliesin-aiによる[PR](https://github.com/wailsapp/wails/pull/5833)）

## v3.0.0-alpha2.118 - 2026-07-26

## 追加

- アイコン生成の入力および出力にデフォルトパスを追加（@taliesin-aiによる[PR](https://github.com/wailsapp/wails/pull/5825)）
- ランタイムのpackage.jsonのsideEffectsにソースエントリーモジュールを追加（@savely-krasovskyによる[PR](https://github.com/wailsapp/wails/pull/5797)）
- コントリビューションガイドにライセンスおよび来歴のセクションを追加（@taliesin-aiによる[PR](https://github.com/wailsapp/wails/pull/5816)）

## 修正

- GTK4のフレームなしウィンドウにスコープを限定したCSSを適用し、角丸を削除（@savely-krasovskyによる[PR](https://github.com/wailsapp/wails/pull/5800)）
- Windowsでポップアップメニューの表示や画面列挙を行う際、カーソル位置の取得に失敗しても適切に処理するよう修正（@wayneforrestによる[PR](https://github.com/wailsapp/wails/pull/5789)）
- macOSのファイルを開くダイアログで拡張子を正しくフィルタリングし、許可されたファイルを接尾辞で検証するよう修正（@phergulによる[PR](https://github.com/wailsapp/wails/pull/5678)）
- Windowsのダークモード初期化でnilのAPIが呼び出されないようガードを追加（@roachadamによる[PR](https://github.com/wailsapp/wails/pull/5793)）
- アップデーターの32ビットビルド失敗を修正。`GOARCH=386`で`fmt.Errorf`に渡した`maxArchiveTotalSize`定数（2 GiB）が、プラットフォームの`int`をオーバーフローしていました。この定数を明示的に`int64`型としました。
- Windows 10 1809 / Windows Server 2019（ビルド17763）など、ダークモードのuxtheme APIを読み込まないWindowsビルドで、ウィンドウがDark（またはシステムのダーク設定）のタイトルバーを使用すると、起動時にnilポインターによるパニックが発生する問題を修正。ウィンドウテーマの設定にある`AllowDarkModeForWindow`呼び出しを、`w32.SetMenuTheme`ですでに使用しているガードと同様にnilチェックで保護するようにしました。

## v3.0.0-alpha2.117 - 2026-07-08

## 追加

- Windowsの非クライアント領域向けにカスタムヒットテストロジックを実装（@savely-krasovskyによる[PR](https://github.com/wailsapp/wails/pull/5462)）

## 変更

- UseVisualHostingに基づいてWebView2のモニタースケール検出を構成（@wayneforrestによる[PR](https://github.com/wailsapp/wails/pull/5761)）

## v3.0.0-alpha2.116 - 2026-07-07

## 追加

- FAQドキュメントを更新し、Wails v3の機能とガイダンスに重点を置くよう変更（@taliesin-aiによる[PR](https://github.com/wailsapp/wails/pull/5763)）

## v3.0.0-alpha2.115 - 2026-07-06

## 修正

- GTK4 Linuxで`Menu.Update()`がネイティブメニューを再構築しない問題を修正（#5659。@puneetdixit200も#5539で独自に診断し修正）
- ディスプレイ変更時にmacOSの画面を列挙するとクラッシュする問題を、画面ID／名前の文字列をコピーし、件数のスナップショットを取得することで修正（#5565。@x-haoseも#5584で独自に診断し修正）
- 最小化／復元の遷移中に`GetClientRect`がnilを返し、`WM_ERASEBKGND`が単色の背景を描画するとWindowsでクラッシュする問題を修正（ガードは@sinspiredが#5636で報告）
- フロントエンドバインディングのエラーが常にテキストとして解析される問題を修正（@mbaklorによる#5690）
- `server`ビルドタグを使用した際にWindowsのビルドが失敗する問題を修正。原因は、macOSおよびLinuxの同等ファイルにはすでに存在する`!server`ビルド制約が、Windows GUIファイルにはなかったことです（#5680）。

## v3.0.0-alpha2.114 - 2026-07-05

## 追加

- Update Manifestプロトコルおよびエンドポイントプロバイダーを実装（@taliesin-aiによる[PR](https://github.com/wailsapp/wails/pull/5720)）

## 変更

- `webview2`バインディングを`v3/internal/webview2`としてv3モジュールに統合し、独立したモジュール、そのナイトリーリリース／同期ワークフロー、およびgo.modのバージョン調整を削除（利用するのはv3のみ）（@taliesin-aiによる[PR](https://github.com/wailsapp/wails/pull/5711)）

## 修正

- WebView2のモニタースケール検出およびDPI変更時のホスト再同期の修正を「未リリース」セクションへ移動（@taliesin-aiによる[PR](https://github.com/wailsapp/wails/pull/5750)）
- float64およびBOOLパラメーターに対するWebView2のCOMマーシャリングを更新（@wayneforrestによる[PR](https://github.com/wailsapp/wails/pull/5741)）
- Windowsのシステムトレイアイコンの更新および破棄時に発生するパニックとnil参照を防止（@wayneforrestによる[PR](https://github.com/wailsapp/wails/pull/5703)）
- Windowsで非表示ウィンドウが正しく再度非表示にならない問題を修正（@wayneforrestによる[PR](https://github.com/wailsapp/wails/pull/5743)）
- WebView2コントローラーの表示状態を、ウィンドウの最小化／最大化／復元と同期（@wayneforrestによる[PR](https://github.com/wailsapp/wails/pull/5742)）

### 修正

- WebView2のモニタースケール検出を再び有効にし、DPI変更時にのみホストを再同期するよう制御（@taliesin-aiによる[PR](https://github.com/wailsapp/wails/pull/5734)）。@randalmurphalが検証した修正を基に、@eleclinが根本原因を確認し、@qq540491950が実機テストを実施

## v3.0.0-alpha2.113 - 2026-07-04

## 追加

- `ANDROID_KEYSTORE_FILE`を設定せずにリリース用 AAB をビルドした場合の警告を追加し（Google Play はデバッグ署名されたバンドルを拒否します）、@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5730)で App Bundle のパッケージ化と署名について文書化
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5739)で Why Wails ドキュメントを複数の言語で追加
- @fbbdev による[PR](https://github.com/wailsapp/wails/pull/5398)で、バインディングにおける Go の time.Time から JS の Date または string へのマッピングをサポート
- Play Store への提出用に Android App Bundle（AAB）のパッケージ化タスク（`bundle`、`bundle:fat`、`assemble:aab`、`assemble:aab:release`）を追加。ローカル／エミュレーターでのテスト用 APK タスクは引き続き利用可能。@mortenolsrud による[PR](https://github.com/wailsapp/wails/pull/5728)（[#5726](https://github.com/wailsapp/wails/issues/5726)を修正）
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5735)で、Android 実機向けタスクターゲットと、カメラ／位置情報の権限の再開処理を追加

## 変更

- `webview2`を v1.0.28 に更新（[リリースノート](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)）。
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5730)で、Android テンプレートの`compileSdk`／`targetSdk`を34から35へ更新。これは Google Play で新しいアプリを提出するために必要

## 修正

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5745)で、sponsorkit に焼き込まれたアバターマスクを修正
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5730)で、バージョンを辞書順で並べ替えることにより、Android AVD の自動作成時に誤ったシステムイメージまたは cmdline-tools のバージョンが選択される問題を修正
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5730)で、セットアップウィザードが古い Android NDK バージョンを提案する問題を修正（現在は文書化された要件と一致する26.3.11579264）
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5744)で、SvelteKit とオプションに関するフランス語ドキュメントを更新
- @flofreud による[PR](https://github.com/wailsapp/wails/pull/5516)で、ディスプレイ変更時に macOS の画面列挙処理で発生する SIGSEGV を修正

## v3.0.0-alpha2.112 - 2026-07-03

## 追加

- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5724)で、Go ベースのコントリビューター SVG ジェネレーターを追加し、ドキュメント／Web サイトのクレジットページを更新
- @fbbdev による[PR](https://github.com/wailsapp/wails/pull/5398)で、バインディングにおける Go の time.Time から JS の Date または string へのマッピングをサポート

## 変更

- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5719)で、Node ベースのスポンサー画像パイプラインを Go ジェネレーターに置き換え

## 修正

- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5729)で、Android ビルドアセットの依存関係インストールスクリプトを修正
- `ValidateAndSanitizeURL`で U+0085（NEXT LINE）制御文字を拒否するようにし、URL バリデーターの空白文字に関する網羅性を確保
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/4785)で、フレームレスウィンドウの DPI 変更時に DWM フレームを再計算
- @yulesxoxo による[PR](https://github.com/wailsapp/wails/pull/4632)で、Windows の表示スケールが100% 以外の場合に DnD ドロップゾーンの検出が失敗する問題を修正
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5714)で、Darwin のダイアログ、メニュー、トレイ、通知全体の Cocoa オブジェクトに明示的な Objective-C メモリ管理を追加
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5718)で、Linux CGO バックエンドのバグとシステムトレイの問題を修正

## v3.0.0-alpha2.111 - 2026-07-01

## 追加

- @Aliuyanfeng による[PR](https://github.com/wailsapp/wails/pull/5061)で、コミュニティショーケースに HappyTools を追加
- @triadmoko による[PR](https://github.com/wailsapp/wails/pull/5643)で、インドネシア語ロケールのサポートと包括的なドキュメントを追加
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/4813)で、WindowsWindow に DisableMenu オプションを追加

## 変更

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5617)で、GOOS と ARCH を使用してビルド／パッケージタスクを振り分けるように Taskfile テンプレートと CLI を更新

## 修正

- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5708)で、macOS ウィンドウのタブ化に関する問題を修正

## 削除

- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5702)で、コントリビューション、機能、ガイドからドイツ語翻訳済み MDX ファイルを削除

## v3.0.0-alpha2.110 - 2026-06-30

## 追加

- @wayneforrest による[PR](https://github.com/wailsapp/wails/pull/5129)で、macOS WebView の再読み込みと強制再読み込みを実装し、WebContent プロセス終了時の復旧処理を追加
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5396)で、コントリビューション、機能、ガイドに関する包括的なドイツ語ドキュメントを追加
- @popaprozac による[PR](https://github.com/wailsapp/wails/pull/5333)で、通知にサウンド、添付ファイル、スケジュール設定、通知更新用 API を追加

## 修正

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/4785)で、フレームレスウィンドウの DPI 変更時に DWM フレームを再計算
- @yulesxoxo による[PR](https://github.com/wailsapp/wails/pull/4632)で、Windows の表示スケールが100% 以外の場合に DnD ドロップゾーンの検出が失敗する問題を修正

## v3.0.0-alpha2.109 - 2026-06-29

## 追加

- @iamhabbeboy による[PR](https://github.com/wailsapp/wails/pull/5026)で、EventsEmit ドキュメントにコードサンプルを追加
- @MerIijn による[PR](https://github.com/wailsapp/wails/pull/5380)で、Windows WebView2 のビジュアルホスティングオプションを追加
- @SametKUM による[PR](https://github.com/wailsapp/wails/pull/5536)で、コミュニティショーケースのドキュメントに Klustr を追加
- @thiennguyen93 による[PR](https://github.com/wailsapp/wails/pull/5685)で、Kira をコミュニティショーケースに追加し、新しいページと変更履歴の項目も追加
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5694)で、MCP サービスガイドにフィードバックセクションを追加

## 変更

- サーバーモードで、デスクトップのビルドタスクと一貫した正式な本番ビルドを利用できるようになりました（#5693）。`task build:server` はデフォルトで本番用バイナリ（`-tags server,production`、`-trimpath`、シンボル削除済み）をビルドし、`DEV=true`（開発サーバー）、`OBFUSCATED=true`（garble）、`EXTRA_TAGS`を受け付けます。`task run:server` は開発サーバーを実行します。`Dockerfile.server` / `task build:docker` は、まず本番サーバー（`-tags server,production`）と本番フロントエンドをビルドします。イメージはデフォルトで distroless/static 上の Pure Go 静的ビルドとなり、CGO アプリ向けに `CGO_ENABLED`、`GO_IMAGE`、`RUNTIME_IMAGE`を上書き可能なビルド引数として公開します。

## 修正

- 保留中の非同期呼び出しがあるウィンドウを閉じた際のクラッシュを防止（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/4435)）
- Windows で非表示のアプリを開いた際にウィンドウがアクティブになることを防止（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5249)）
- WebKit のリクエストメタデータ、レスポンス完了処理、ボディストリーム処理が GTK メインスレッドで実行されることを保証（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5668)）
- GTK4 Linux で `Menu.Update()` がネイティブメニューを再構築しない問題を修正（#5659、@puneetdixit200 も #5539 で独自に診断して修正）
- ディスプレイ変更時に macOS の画面を列挙するとクラッシュする問題を、画面の ID／名前文字列をコピーし、画面数のスナップショットを取得することで修正（#5565、@x-haose も #5584 で独自に診断して修正）
- Windows で異なる DPI のモニター間にウィンドウをドラッグした後、WebView2 のコンテンツが縮小してから消える問題を、最小化解除時の DPI 再同期と同様に `WM_DPICHANGED` ハンドラー内でコントローラーの境界を再設定することで修正（#5677）

## v3.0.0-alpha2.108 - 2026-06-28

## 追加

- `app.GlobalShortcut`（`Register`、`Unregister`、`UnregisterAll`、`IsRegistered`、`GetAll`）によるグローバル（システム全体）キーボードショートカットを追加。アプリケーションにフォーカスがない場合でもショートカットが発火します。サードパーティ依存関係を使用せず、プラットフォームごとにネイティブ実装されています。macOS では Carbon ホットキー、Windows では `RegisterHotKey`、X11 では `XGrabKey`、Wayland では XDG Desktop Portal のグローバルショートカットインターフェースを使用します。
- 組み込み MCP サーバーを追加。これは `mcp` タグを指定してアプリケーションをビルドすると自動的に起動する Model Context Protocol サーバーで、LLM エージェントが実行中の Wails アプリケーションをテストおよび制御できます。ウィンドウ制御、DOM インスペクション、JavaScript の評価、バインドされたメソッドの呼び出し、イベント、および画面上のアニメーションカーソルで表示されるマウス／キーボード入力のシミュレーションに対応します。ユーザーコードは不要です。`WAILS_MCP=1` が設定されている場合、`wails3 build`/`wails3 dev` によって `mcp` タグが自動的に追加されます。設定はすべて環境変数（`WAILS_MCP_HOST`、`WAILS_MCP_PORT`、`WAILS_MCP_TIMEOUT`、`WAILS_MCP_HIDE_CURSOR`）で行います。

## 修正

- GTK4 Linux で `Menu.Update()` がネイティブメニューを再構築しない問題を修正（#5659、@puneetdixit200 も #5539 で独自に診断して修正）
- ディスプレイ変更時に macOS の画面を列挙するとクラッシュする問題を、画面の ID／名前文字列をコピーし、画面数のスナップショットを取得することで修正（#5565、@x-haose も #5584 で独自に診断して修正）

## v3.0.0-alpha2.107 - 2026-06-27

## 追加

- サイドバーナビゲーションを備えた実験的な Wake ドキュメントを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5613)）

## v3.0.0-alpha2.106 - 2026-06-24

## 変更

- `webview2` を v1.0.27 に更新。
  - ci(webview2): リリースビルドを修正（Windows のクロスコンパイルと完全な go.sum）（#5671）\

  **完全な差分：** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- webview2 のリリースワークフローのクロスコンパイルから go vet を削除（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5672)）
- auto-changelog の OpenRouter モデルを google/gemini-2.5-flash-lite に更新（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5670)）
- `webview2` を v1.0.26 に更新。

### 修正

- **一時的なランタイム COM エラーで終了せず、復旧するように変更**（#5658、#5580）。以前は `Chromium.errorCallback` が<em>あらゆる</em> COM エラーに対して `os.Exit(1)` を呼び出していたため、起動後に復旧可能な一時的不具合が発生しただけでも、アプリケーション全体が終了していました。ランタイムの処理経路（`Resize`/`GetClientRect`、`Navigate`/`NavigateToString`、`Init`、`MessageReceived`、`PutZoomFactor`、`OpenDevToolsWindow`）では、エラーをログに記録して復旧するようになりました。特に、`MessageReceived` 内の不正な形式または信頼できない Web メッセージは、プロセスを停止させず破棄するようになりました。これにより、異なる DPI のモニター間を移動した際に発生する一連のクラッシュ（#5544、#5650）に対処します。環境およびコントローラーの作成処理経路では、引き続き致命的エラーとなります。\

**完全な差分：** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## 修正

- release-webview2 ワークフローが go.sum ファイルを正しく処理するように修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5671)）
- ネイティブメニューをクリアして再構築することで Linux GTK4 のメニュー更新を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5659)）

## v3.0.0-alpha2.105 - 2026-06-21

## 追加

- 共有コードから実行時にプラットフォームを検出するための `application.System` を追加。`System.IsMobile()`（iOS/Android）、`System.IsDesktop()`（macOS/Windows/Linux）、`System.IsServer()`（`server` ビルドタグ）、および単一のターゲットを直接テストするための `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` を提供します。すべてのターゲットでコンパイルできるため、ビルドタグを使わずに分岐できます。対応するフロントエンドヘルパー（`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`）は `@wailsio/runtime` で利用できます。
- 独自の Vite プロジェクトを `frontend/` に組み込む方法を示す「その他のフロントエンドフレームワークの使用」ガイドを追加（Solid、Preact、Lit、SvelteKit、Qwik、Angular などを対象）
- `wails3 setup` ウィザードで、モバイル（iOS/Android）ツールチェーン（Xcode と iOS Simulator ランタイム、JDK、Android SDK/NDK、エミュレーター）を確認するようになりました。該当する場合は、ワンクリックインストールとコピー可能なシェル設定の修正方法も提供します。
- 生成されるプロジェクトに、7 日間の `minimum-release-age` を設定する `frontend/.npmrc` を同梱し、公開直後の（侵害されている可能性がある）パッケージにさらされるリスクを低減します（pnpm と bun では適用され、npm では問題なく無視されます）。

## 変更

- すべての組み込みスターターテンプレートを、新しいネオンマウンテンのヒーロービジュアル（Web、iOS、Android）で再設計
- **TypeScript がスターターテンプレートのデフォルトとなり、修飾子のないテンプレート名を使用するようになりました。** `wails3 init`（`-t` なし）は TypeScript プロジェクトを生成します。`-t vanilla`、`-t react`、`-t vue`、`-t svelte` は TypeScript 版で、JavaScript 版はそれぞれ `-t vanilla-js`、`-t react-js`、`-t vue-js`、`-t svelte-js` です。組み込みテンプレートは `template.yaml` 内の `typescript:` で言語を宣言します。`-ts` サフィックスを使用するコミュニティテンプレートは、フォールバックとして引き続き機能します。
- `wails3 setup` ウィザードを、ネオン調の「digital Wails」テーマ（山を背景にしたすりガラス風の鮮やかな表現）で再設計

## 修正

- WebView2 が一時停止するか、そのレンダリング/GPU プロセスが再生成されるほど長時間アプリを最小化した後、復元すると Windows でクラッシュする問題を修正。最小化/復元時の DPI 再同期（#5544）では、ウィンドウの DPI が実際に変更された場合にのみ WebView2 コントローラーを操作するようになり、一般的な同一 DPI での復元時に、一時停止中のコントローラーに対して致命的な COM 呼び出しが行われることを回避（#5605）
- アセットやメディアを頻繁に読み込む、長時間稼働中の Linux アプリで、ネイティブの `SIGABRT`/`SIGSEGV` クラッシュ（通常は GTK メインループ実行中の `g_object_unref` 内）が繰り返し発生する問題を修正。アセットサーバーがワーカー goroutine から `WebKitURISchemeRequest` を完了させ、スレッドセーフでない WebKit2GTK 関数を GTK メインスレッド以外から呼び出していたため、完了処理（`webkit_uri_scheme_request_finish_with_response`/`finish_error`）をメインスレッドで実行するように変更。#5566 の部分的な修正を完了。GTK3 ビルドと GTK4/WebKitGTK 6.0 ビルドの両方に影響（#5631、#5557）
- Linux/GTK3 の `setupSignalHandlers` で断続的に発生する `fatal error: invalid pointer found on stack` を修正。シグナルの `user_data` として渡されるウィンドウ ID が Go のローカルな `unsafe.Pointer` に保持されていたため、スタックコピー中にガベージコレクターがその（ポインターではない）値をスキャンすると異常終了していた。Go 側で ID を整数型（`uintptr_t`）のまま保持するようにし、GTK4 パスに #4958 で適用されたものと同じ修正を、従来の GTK3 パスにもバックポート。GTK4 パスでは、`-race`/checkptr エラーを解消するため、C のシグナル関数を `uintptr_t` に変更していた（#5631）

## 削除

- `react-swc`、`preact`、`lit`、`solid`、`qwik`、`sveltekit` の各スターターテンプレート（およびそれらの `-ts` バリアント）を削除。サポート対象の組み込みテンプレートは `vanilla`、`react`、`vue`、`svelte` となり、いずれもデフォルトでは TypeScript を使用し、JavaScript バリアントは `-js` となる。その他のフレームワークも、引き続き[独自のフロントエンドを用意する](https://v3.wails.io/guides/dev/frontend-frameworks)か、カスタムテンプレートを使用して利用可能

## v3.0.0-alpha2.104 - 2026-06-18

## 修正

- バインドされた Go サービスメソッドが空文字列を返した際に発生する iOS のクラッシュ（SIGABRT）を修正。iOS のアセットレスポンスライターが本文の長さではなく、`buf != nil` で本文ポインターをチェックしていたため、長さゼロの本文で `&buf[0]` がパニックしていた。デスクトップ版のライターと同様に、長さをチェックするように変更

## v3.0.0-alpha2.103 - 2026-06-15

## 変更

- iOS と Android のネイティブ機能をプラットフォームマネージャーへ移行。従来の `application.IOS*`/`application.Android*` 自由関数ではなく、`application.IOS.*` と `application.Android.*`（例：`application.IOS.Haptic("medium")`、`application.Android.Share(payload)`）を介して呼び出すように変更（#5602）
- モバイルブリッジイベントの名前を変更。クロスプラットフォームイベントでは `common:*` プレフィックス（例：`common:haptic`、`common:location`）を使用し、プラットフォーム固有のイベントでは `ios:*` / `android:*`（例：`ios:backgroundTask`、`android:foregroundService`）を使用するようになった。`native:*` プレフィックスは廃止（#5602）

## v3.0.0-alpha.102 - 2026-06-14

## 追加

- 対話形式のプロジェクトセットアップと依存関係チェックを行う、実験的な `wails3 setup` ウィザードを追加
- 機械可読な出力を生成する `--json` フラグを `wails3 doctor` に追加
- `wails3 doctor` コマンドに署名ステータスのセクションを追加

## 修正

- Linux での npm 検出を修正し、パッケージマネージャーに加えて PATH も確認するように変更

## v3.0.0-alpha.101 - 2026-06-13

## 追加

- iOS：ネイティブメッセージダイアログ（UIAlertController）と、ファイル、複数ファイル、ディレクトリを開くダイアログ（UIDocumentPickerViewController）を追加。保存ダイアログは明示的なエラーを返す
- iOS：UIPasteboard によるクリップボード対応を追加
- iOS：UIScreen による実際の画面メトリクス（ポイント、ピクセル、スケール、セーフエリア内の作業領域）を追加
- iOS：実機向けビルド（`IOS_PLATFORM=device`）、コード署名 ID／プロビジョニングプロファイル／entitlements のサポート、`.ipa` パッケージ化、devicectl を介した `deploy-device` を追加
- iOS：設定可能な最小 iOS バージョン（build/config.yml の `ios.minIOSVersion`）を追加
- iOS：macOS 上で `wails3 doctor` が Xcode と iOS SDK の利用可否を報告するように変更
- iOS：バッテリー、ネットワーク、テーマ、画面ロック、メモリ不足の各システムイベントを、`events.IOS.*` およびプラットフォーム非依存の `events.Common.*` アプリケーションイベントとして公開
- iOS：ネイティブモバイル機能ブリッジ（エクスポートされる `application.IOS*`）を追加。共有シート、URL を開く機能、スリープ防止、トーチ、セーフエリアのインセット、明るさ、アプリ情報、画面方向のロック、ステータスバー、生体認証（Face ID/Touch ID）、ローカル通知、Keychain によるセキュアストレージに対応
- iOS：センサーとハードウェア機能として、触覚フィードバック、単発の位置情報取得、加速度計、近接センサー、テキスト読み上げ、ストレージ情報、電源／バッテリー状態、ネットワーク状態、キーボードのインセット、画面キャプチャ検出を追加
- iOS：ドキュメント（IOS.md およびドキュメントサイトのガイド）を追加
- Android：ネイティブメッセージダイアログ（AlertDialog）と、ファイルおよび複数ファイルを開くダイアログ（Storage Access Framework。キャッシュのコピーとしてインポート）を追加。ディレクトリを開くダイアログと保存ダイアログは明示的なエラーを返す
- Android：ClipboardManager によるクリップボード対応を追加
- Android：WindowMetrics/DisplayMetrics による実際の画面メトリクス（dp、ピクセル、スケール、システムバーを除いた作業領域）を追加
- Android：触覚フィードバック（`Android.Haptics.Vibrate`）、デバイス情報（`Android.Device.Info`）、トースト（`Android.Toast.Show`）のランタイムメソッドを追加
- Android：型付きライフサイクルイベント（`events.Android.*`、events.txt から生成）を追加し、`ActivityCreated` を `Common.ApplicationStarted` にマッピング
- Android：ビルドパイプラインで、インストール可能なデバッグ APK とリリース APK（`android:run`、`android:package`、`android:package:fat`）を生成。リリース署名には、デフォルトでデバッグ用キーストアを使用するか、`ANDROID_KEYSTORE_*` 環境変数を介して実際のキーストアを使用可能
- Android：`wails3 doctor` が Android SDK、NDK、JDK を報告するように変更
- Android：バッテリー、ネットワーク、テーマ、画面ロック、メモリ不足の各システムイベントを、`events.Android.*` およびプラットフォーム非依存の `events.Common.*` アプリケーションイベントとして公開
- Android：ネイティブモバイル機能ブリッジ（エクスポートされる `application.Android*`）を追加。共有、URL を開く機能、スリープ防止、トーチ、セーフエリアのインセット、明るさ、アプリ情報、画面方向のロック、ステータスバー、生体認証（BiometricPrompt）、ローカル通知、EncryptedSharedPreferences によるセキュアストレージに対応
- Android：センサーとハードウェア機能として、触覚フィードバック、単発の位置情報取得、加速度計、近接センサー、テキスト読み上げ、ストレージ情報、電源／バッテリー状態、ネットワーク状態、キーボードのインセット、FLAG_SECURE による画面キャプチャのブロックを追加
- Android：ドキュメント（ANDROID.md およびドキュメントサイトのガイド）を追加
- サンプル：`mobile` の各種機能を網羅したデモに Mobile タブと Hardware タブを追加し、iOS と Android のネイティブ機能ブリッジを実演（ピル型タブは複数行に折り返し）
- モバイル：バッテリー消費を抑えるため、アプリがバックグラウンドに移行すると、加速度計、近接センサー、トーチ、サンプルの定期クロックを一時停止し、復帰時に再開（Android ではバックグラウンドでもプロセスが実行され続け、iOS ではトーチが維持されるハードウェア状態である）。また、Android のシステムイベントレシーバーは、アプリがフォアグラウンドにある間だけ登録
- iOS：カメラ撮影 — `application.IOSCapturePhoto`/`IOSCaptureVideo`（UIImagePickerController → base64 サムネイルを含む `native:capture` イベント）
- iOS：バックグラウンド実行 — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask`（UIApplication のバックグラウンドタスク実行時間枠）、および生成される Info.plist に `UIBackgroundModes` をテンプレート展開する、設定可能な `ios.backgroundModes`（build/config.yml）
- Android：カメラ撮影 — `application.AndroidCapturePhoto`/`AndroidCaptureVideo`（FileProvider 経由のシステムカメラ → `native:capture` イベント）
- Android：フォアグラウンドサービス — `application.AndroidStartForegroundService`/`AndroidStopForegroundService`（継続的な通知を伴う `WailsForegroundService` により、長時間実行されるバックグラウンド処理中もプロセスを維持）
- 例：写真／動画撮影とバックグラウンド実行を実演する Camera タブ（Android ではフォアグラウンドサービス、iOS ではバックグラウンドタスク実行時間枠）

## 修正

- Linux で `getUserMedia` が常に `NotAllowedError` で失敗する問題を修正しました。WebKitGTK は処理するハンドラーがないパーミッション要求を拒否しますが、`permission-request` シグナルが接続されていませんでした。カメラ／マイクは、新しいクロスプラットフォームの `WebviewWindowOptions.Permissions` マップ（`map[PermissionType]Permission`）に基づいて処理されるようになり、Linux（WebKitGTK）と Windows（WebView2）の両方で適用されます。ネイティブプロンプトがない Linux では、カメラ／マイクはデフォルトで許可され（`getUserMedia` の動作を復元）、`PermissionDeny` で無効にできます（#5552）
- iOS：`GOOS=ios` が再びコンパイル可能になりました（`events.IOS` のエクスポート、モバイル用メソッド名スタブ）。また、production タグ付きビルドもコンパイル可能になりました（pkg/application および複数のサービスのビルドタグを修正）
- iOS：Go→JS イベントと ExecJS が動作するようになりました。起動時にページが二重に読み込まれなくなり、`wails:runtime:ready` ハンドシェイクが失われることもなくなりました
- iOS：`ApplicationDidFinishLaunching`/`ApplicationStarted` がアプリの起動と競合しなくなりました。固定の 2 秒の起動待機を削除しました
- iOS：Go→JS の JavaScript 実行ごとに発生していた C 文字列のリークを修正しました
- iOS：`hasListeners` が実際のリスナー登録状態を反映するようになりました
- iOS：本番ビルドではフレームワークのデバッグログがコンパイル時に除外されるようになりました
- Android：`GOOS=android` が再びコンパイル可能になりました。`events.Android` を定義し、範囲外アクセスを起こす `events_android.go` リスナー配列を削除し、モバイル用メソッド名スタブを追加したほか、デスクトップ Linux 用ファイル（`linux_cgo.*`、`events_linux.*`、`environment_linux.go`）が Android ビルドに混入しないようにしました
- Android：JS→Go バインディングが動作するようになりました。WebView は `fetch()` の POST ボディを `shouldInterceptRequest` に渡せないため、ランタイム呼び出しは JavascriptInterface トランスポート（`nativeHandleRuntimeCall`）経由でルーティングされ、nil のリクエストボディによるクラッシュを回避します
- Android：`Screens.*` ランタイム呼び出しが実際のデータを返すようになりました。起動時に ScreenManager が設定されるようになりました（これまでは接続されていなかったため、`GetAll` が nil を返していました）
- Android：本番ビルドではフレームワークのデバッグログがコンパイル時に除外され、デバッグビルドでは `Wails` タグで logcat に出力されるようになりました
- Android：実際に機能する `hasListeners` レジストリ、JNI の参照／例外処理、および単一読み込みのページライフサイクル（二重ナビゲーションなし）を実装しました
- Vite 開発サーバーの実行中に Windows で `wails3 generate bindings` が「Access is denied」で失敗する問題を修正しました。出力ディレクトリを置き換えるように名前変更するのではなく、生成ファイルを出力ディレクトリ内に同期するようにしました（#5515）
- ディスプレイ変更後に画面情報を読み取ると、macOS で断続的に致命的なクラッシュが発生する問題を修正しました。画面 ID と名前が、Go によるコピー前に解放される可能性のある autorelease 済みの `UTF8String` バッファへのポインターを保持していました（解放後使用）。文字列は `strdup` され、変換後に解放されるようになりました。また、画面の列挙を明示的な autorelease pool 内で実行することで、Go の goroutine から呼び出した場合もリークしなくなりました（#5556）
- assetserver が `WebKitURISchemeRequest` を閉じる際、Linux で断続的に SIGSEGV が発生する問題を修正しました。最後の `g_object_unref` が assetserver の goroutine 上で実行され、GTK メインスレッド外で WebKit GObject をファイナライズしていました。unref は `g_main_context_invoke` を介して GTK メインコンテキスト上にマーシャリングされるようになりました（#5557）

## v3.0.0-alpha.100 - 2026-06-13

## 追加

- `MacWebviewPreferences` に WKWebView の追加設定オプション `EnableAutoplayWithoutUserAction`、`AllowsAirPlayForMediaPlayback`、`AllowsMagnification`、`JavaScriptCanOpenWindowsAutomatically`、`MinimumFontSize`、`ApplicationNameForUserAgent` を追加しました（#5549）

## 修正

- Vite 開発サーバーの実行中に Windows で `wails3 generate bindings` が「Access is denied」で失敗する問題を修正しました。出力ディレクトリを置き換えるように名前変更するのではなく、生成ファイルを出力ディレクトリ内に同期するようにしました（#5561）
- Linux のフレームレスウィンドウで JS のリサイズイベントが発火しない問題と、フレームレスウィンドウのスクロールバー端検出を修正しました（#5368）
- 一時ディレクトリとインストールディレクトリが異なるボリュームにある場合、Windows のアップデーターが「invalid cross-device link」で失敗する問題を修正しました（#5560）

## v3.0.0-alpha.99 - 2026-06-10

## 修正

- Vite 開発サーバーの実行中に Windows で `wails3 generate bindings` が「Access is denied」で失敗する問題を修正しました。出力ディレクトリを置き換えるように名前変更するのではなく、生成ファイルを出力ディレクトリ内に同期するようにしました（#5515）

## v3.0.0-alpha.98 - 2026-06-03

## 修正

- アイドル時（インスペクターを開いている場合など）に Linux で WebKit UI がフリーズする問題を修正しました。JavaScriptCore の GC スレッド同期を破壊していた `SIGUSR1` への `SA_ONSTACK` の強制を廃止しました（#5527）

## v3.0.0-alpha.97 - 2026-05-31

## 追加

- デバッグページと `runtime/trace` の操作方法を追加しました

## 変更

- 不要な `_ "embed"` インポートをいくつか削除し、コードを少し整理しました

## 修正

- Windows でウィンドウの最大化を解除した後、最小幅／高さの制約が適用されない問題を修正しました（#4593）
- Frameless + Transparent ウィンドウオプションを使用したフルスクリーンモードで、マウスクリックが背後に透過する問題を修正しました（#4408）

## v3.0.0-alpha.96 - 2026-05-25

## 追加

- Garble 難読化のサポート（[#4563](https://github.com/wailsapp/wails/issues/4563)）を追加しました。安定したバインディングメソッド ID、ビルド／Taskfile の連携（`build --obfuscated --garbleargs`、`generate bindings -obfuscated`）、およびランタイム向けのすべてのペイロードへの JSON 構造体タグ（`EnvironmentInfo`、`OSInfo`、`Screen`、`Rect`、`Point`、`Size`、`Capabilities`）により、Garble がエクスポート済みフィールドの名前を変更してもワイヤーフォーマットが維持されます。

## v3.0.0-alpha.95 - 2026-05-20

## 追加

- 不足していたプロジェクト構成ページを追加しました

## 変更

- ドキュメント：アーキテクチャページの図をいくつか変更し、よりすっきり表示されるようにシーケンス図を使用しました
- ドキュメント：実行の前提条件として D2 のインストールが必要であることを示す注記を追加

## 修正

- GTK4 がデフォルトの場合の `wails3 generate appimage` を修正。バンドラーはランタイムファイルを検索する前にバイナリから GTK スタックを検出するようになったため、GTK4 ビルドでは `libwebkitgtkinjectedbundle.so`（`webkitgtk-6.0/` 配下）を、`-tags gtk3` ビルドでは `libwebkit2gtkinjectedbundle.so`（`webkit2gtk-4.1/` 配下）を選択します。`.relr.dyn` のプローブでも `libgtk-4.so.1` を確認するようになり、スタックにかかわらず、最新のツールチェーンではストリッピングが正しく無効化されます。（#5475）
- 相対 `-builddir` を指定して呼び出すと `wails3 generate appimage` が失敗する問題を修正。バンドラーは `-binary`、`-icon`、`-desktopfile`、`-builddir`、`-outputdir` を処理の開始時に絶対パスへ解決するようになったため、処理途中の `s.CD` によって AppRun のダウンロード用 goroutine やコピー後の `ldd` プローブが機能しなくなることはありません。
- デスクトップの `Name=` フィールドがバイナリのベース名と一致しない場合に、`wails3 generate appimage` が最終的な AppImage を `-outputdir` へ移動できない問題を修正。バンドラーは、linuxdeploy の appimage プラグインに対し、デスクトップファイルから導出した名前ではなく `<binary>-<arch>.AppImage` へ AppImage を書き込むよう、`OUTPUT` 環境変数を介して強制するようになりました。
- alpha.93 で GTK4 + WebKitGTK 6.0 スタックがデフォルトになった後、Linux で `events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`、`Common.SystemDidWake` が発生しない問題を修正。新しいデフォルトの `application_linux.go` `run()` は、`Linux.*` イベントを対応する `Common.*` へ転送する `setupCommonEvents()` も、`monitorPowerEvents()` も呼び出していませんでした。DBus 電源監視ヘルパーは、`application_linux_dbus.go` を介して GTK3 と GTK4 のビルドパス間で共有されるようになりました。（#5474）

## v3.0.0-alpha.94 - 2026-05-19

## 修正

- alpha.93 で GTK4 + WebKitGTK 6.0 スタックがデフォルトになった後、Linux で `events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`、`Common.SystemDidWake` が発生しない問題を修正。新しいデフォルトの `application_linux.go` `run()` は、`Linux.*` イベントを対応する `Common.*` へ転送する `setupCommonEvents()` も、`monitorPowerEvents()` も呼び出していませんでした。DBus 電源監視ヘルパーは、`application_linux_dbus.go` を介して GTK3 と GTK4 のビルドパス間で共有されるようになりました。（#5474）

## v3.0.0-alpha.93 - 2026-05-17

## 追加

- Linux の `wails3 doctor` 出力に `XDG_SESSION_TYPE` を追加（@leaanthony）

## 修正

- appmenu-gtk-module が実体化されていないウィンドウへアクセスすることによって発生する、Wayland 上のウィンドウメニューのクラッシュを修正（#4769、@leaanthony）
- アプリ名に無効な文字（空白、丸括弧など）が含まれる場合に GTK アプリケーションがクラッシュする問題を修正（@leaanthony）
- Windows でドラッグ＆ドロップを初期化する際に発生する「メモリが不足しています」エラーを修正（#4701、@overlordtm）
- mainthread のコールバックストアで、マップからの削除に誤って RLock を使用していたことによる競合状態を修正（Linux、macOS、iOS）（#4424、@leaanthony）
- タスクへコマンドライン引数を渡す際の変数処理を修正しました。KEY=VALUE ペアで指定した CLI 変数が正しく初期化され、タスク実行全体へ伝播されるようになりました。
- macOS での NSWindowZoomButton の競合を修正。起動時と実行時の両方で、`MaximiseButtonState` と `FullscreenButtonState` は、より制限の厳しい状態を適用するようになりました。どちらのセッターも、もう一方の設定を暗黙的に上書きすることはできません（#5319）
- #5463 で CodeRabbit により明らかになった、従来の GTK3 ビルドパス（`-tags gtk3`）に以前から存在していた一連のバグを修正。ファイル関連付けによる起動でスタートアップハンドラーがスキップされなくなりました。`getTheme` は境界および型に対して安全になりました。`appName` は GLib が所有するメモリを解放しなくなりました。`clipboardGet` は GTK が返した `gchar*` をリークしなくなりました。`Calloc` はポインタレシーバーを使用するようになり（また、`NewCalloc` は `*Calloc` を返すようになり）、プールが割り当てを実際に追跡するようになりました。`zoomOut` は、1.0 にクランプされる負の乗数ではなく、`zoomInFactor` の逆数を使用するようになりました。`execJS` は、呼び出しごとに `C.CString("")` をリークせず、事前割り当て済みの空のワールド名を再利用するようになりました。開発用の `fmt.Println` を `menuItem.setAccelerator` から削除しました。#5465 を解決しました。
- デフォルトの GTK4 ビルドパス（`linux_cgo.go`）にある、同じ `Calloc` の値レシーバーによるリークを修正。ポインタレシーバーと `NewCalloc() *Calloc` により、ウィンドウごとの `c.String(...)` の割り当てが実際に追跡され、解放されるようになりました。

## v3.0.0-alpha.92 - 2026-05-15

## 追加

- `PACKAGE_MANAGER` オプションを介して、使用するフロントエンドパッケージマネージャーを制御できるように Taskfile を変更
- Taskfile テンプレートをより予測しやすく記述できるように、テンプレートデータへ `{{.Opn}}` と `{{.Cls}}` を追加

## 変更

- 既存の Taskfile のうち数個を、`{{.Opn}} and {{.Cls}}` を使用するように変更

## 修正

- パネルがトレイメニューを読み取っている間にメニューが更新されると、`linuxSystemTray` で `concurrent map read and map write` ランタイム致命的エラーが発生する問題を修正。
- WebView2 のエラーおよびスタックトレースの出力に `fmt` ではなく `log` を使用するように変更し、Windows でコンソールを接続せずにアプリを実行した場合でもメッセージが失われないようにしました。

## v3.0.0-alpha.91 - 2026-05-12

## 変更

- スポンサー SVG を更新（`@github-actions[bot]` による [PR](https://github.com/wailsapp/wails/pull/5414)）
- **破壊的変更（macOS）：** macOS の座標系を正規化し、`GetScreens`、`Position`、`SetPosition` がすべて同じ座標空間（論理ポイント、Y 軸は下向き、プライマリ画面の左上を `(0,0)` とする）を使用するようにしました。これは Windows、GTK、および Electron と Web の公開 API と一致します。プライマリ画面より物理的に上にある画面では、`Bounds.Y` が正の値ではなく負の値として報告されるようになりました。また、`Position()`/`SetPosition()` の値は `points × primaryScale` ではなく論理ポイント単位になりました。`Position()` → `SetPosition()` の往復変換は維持されます。以前の alpha ビルドで記録した絶対値や、手動で計算した回避策（例：`primaryScale` を掛ける、または画面の高さを基準に Y を反転する）は更新が必要です。[#5117](https://github.com/wailsapp/wails/issues/5117) を解決しました。

## 修正

- [PR](https://github.com/wailsapp/wails/pull/5416) で DBus シグナル名とボディの長さを防御的に検証し、パニックを防止（@leaanthony）
- [PR](https://github.com/wailsapp/wails/pull/5363) で Linux の GTK メニュー処理におけるメモリ安全性の問題を修正（@leaanthony）
- [PR](https://github.com/wailsapp/wails/pull/5295) で NVIDIA GPU を検出し、Linux 上の DMA-BUF レンダラーを無効化（@leaanthony）
- macOS での `SetPosition` の画面間 Y 座標変換を修正。[#5117](https://github.com/wailsapp/wails/issues/5117) でプライマリ画面の高さをグローバルな基準として使用し、プライマリディスプレイから垂直方向にずれたモニターでもウィンドウが正しい位置に配置されるようにしました。
- [PR](https://github.com/wailsapp/wails/pull/5109) で、正しいフィードバック URL を指すように git PR テンプレートを修正（@wayneforrest）
- 壊れた `DestroyMenu` システムコールが1つではなく4つの引数を渡していたため、すべての呼び出しが FALSE を返し、何も解放されなかったことに起因する、一連の Windows システムトレイ `SetMenu` クラッシュを修正。また、メニューの再構築時に HMENU および HBITMAP ハンドル（`MenuItem.SetBitmap` によって実行時に割り当てられたものを含む）を解放し、`Win32Menu.Update` 内の古いチェックボックス／ラジオマップをリセットし、割り当てを倍増させていた `systemtray.updateMenu` 内の冗長な `Update()` 呼び出しを削除。長時間稼働するシステムトレイアプリで、メニューを再構築するたびに GDI/USER オブジェクトがリークしなくなりました。

## v3.0.0-alpha.90 - 2026-05-11

## 追加

- macOS の WKWebView User-Agent に設定可能なアプリケーション名を追加（@vinhvoit225 による [PR](https://github.com/wailsapp/wails/pull/5261)）
- gin-service サンプルに間接依存関係 github.com/coder/websocket を追加（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5400)）
- ビルドアセットのテストにディープイコール比較のサポートを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5402)）

## 変更

- ビルド出力を assets ディレクトリに集約（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5401)）
- スポンサー SVG を更新（`@github-actions[bot]` による [PR](https://github.com/wailsapp/wails/pull/5399)）

## 修正

- macOS のシングルインスタンスメッセージに通知オブジェクトを使用（@overlordtm による [PR](https://github.com/wailsapp/wails/pull/5289)）
- 高負荷時の Promise 消失を防ぐため、Windows コールバックをバッチ処理（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5383)）

## v3.0.0-alpha.89 - 2026-05-10

## 追加

- Go テスト結果を集約する go<em>test</em>results ジョブを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5316)）

## 変更

- 大きな RPC ペイロードを条件に応じてチャンク化された POST リクエストに分割（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5369)）
- すべてのフロントエンドテンプレートで Vite を 5.x.x から 8.0.0 にアップグレード（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5386)）
- Vite 開発サーバーのポート設定を環境変数へ移行（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5365)）
- すべてのテンプレートで Vite 開発サーバーが 127.0.0.1 にバインドするよう設定（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5361)）
- スポンサー SVG を更新（`@github-actions[bot]` による [PR](https://github.com/wailsapp/wails/pull/5384)）

## 修正

- build-assets の更新時に Info.plist テンプレートスタブをサニタイズ（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5312)）
- メニュー項目のミューテーター（`setMenuItemChecked()`、`setMenuItemLabel()`、`setMenuItemDisabled()`、`setMenuItemHidden()`、`setMenuItemTooltip()`）をメインスレッドで同期的に適用することで、macOS メニューの古い状態を修正。すばやく再表示した際にメニューが直前の状態で描画される原因となっていた `dispatch_async` の競合を解消しました（#5002）
- 不要な再ビルドを防ぐため、開発モードで `*_test.go` ファイルを無視（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5203)）
- アプリが実行中でない場合に Menu.Update() でセグメンテーションフォルトが発生する問題を修正（@wucm667 による [PR](https://github.com/wailsapp/wails/pull/5291)）
- Windows で lastSizeWParam を使用してメニューバーの再描画を制御（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5382)）

## v3.0.0-alpha.88 - 2026-05-09

## 変更

- HiddenOnTaskbar で WS<em>EX</em>TOOLWINDOW を使用するよう変更（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5371)）
- 依存関係を並べ替え、go.mod から webview2 の replace ディレクティブを削除（@atterpac による [PR](https://github.com/wailsapp/wails/pull/5370)）
- スポンサー SVG を更新（`@github-actions[bot]` による [PR](https://github.com/wailsapp/wails/pull/5358)）

## 修正

- ジェネリックな間接参照エイリアスを削除し、マップのキー型を統合（@fbbdev による [PR](https://github.com/wailsapp/wails/pull/5331)）

## 削除

- PR-master ワークフローを削除し、ドキュメント、Go テスト、スキップテストを廃止（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5377)）

## v3.0.0-alpha.87 - 2026-05-07

## 追加

- Wails v3 の韓国語ドキュメントを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5352)）
- インストールとクイックスタートのフランス語ドキュメントを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5354)）
- クイックスタート、概念、コミュニティのポルトガル語ドキュメントを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5355)）

## v3.0.0-alpha.86 - 2026-05-06

## 追加

- フランス語のドキュメントローカライズを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5328)）
- ドキュメントサイトにドイツ語ロケールを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5343)）

## 変更

- 翻訳済みの全 8 ロケールをドキュメント設定に登録（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5347)）
- WebView2 向けに Windows 関連の各種ファイルを更新（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5317)）

## 修正

- Linux のダイアログディスパッチを GTK3 と GTK4 に分割（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5340)）
- ダイアログのコールバックが GTK スレッドで確実に実行されるようにし、セグメンテーションフォルトを修正（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5339)）

## v3.0.0-alpha.85 - 2026-05-05

## 追加

- リポジトリに PR テンプレートの URL を追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5179)）
- Wails v3 のドイツ語ドキュメントを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5330)）

## v3.0.0-alpha.84 - 2026-05-03

## 追加

- macOS で Escape キーによるフルスクリーン終了を無効にするオプションを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5307)）
- macOS で Escape キーによるフルスクリーン終了を無効にするオプションを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5310)）
- @yuseferi により、Pausa コミュニティショーケースのドキュメントを[PR](https://github.com/wailsapp/wails/pull/5288)で追加

## 変更

- スポンサー SVG を更新（`@github-actions[bot]` による [PR](https://github.com/wailsapp/wails/pull/5308)）
- @leaanthony により、未対応プラットフォームを処理できるようアイコン生成コマンドを[PR](https://github.com/wailsapp/wails/pull/5309)で更新
- @leaanthony により、真偽値型のフルスクリーン API を三状態の ButtonState に置き換え、プラットフォームバインディングを[PR](https://github.com/wailsapp/wails/pull/5224)で実装

## 修正

- @leaanthony により、コントローラーが nil の状態でも安全になるよう WebView2 のフォーカス操作を[PR](https://github.com/wailsapp/wails/pull/5315)で保護
- @leaanthony により、PR のベースブランチを正しく参照するよう GitHub Actions ワークフローを[PR](https://github.com/wailsapp/wails/pull/5313)で更新
- @leaanthony により、不要な再ビルドを防ぐため、開発モードで`*_test.go`ファイルを無視するよう[PR](https://github.com/wailsapp/wails/pull/5203)で修正
- @wucm667 により、アプリが実行されていないときに Menu.Update() でセグメンテーションフォールトが発生しないよう[PR](https://github.com/wailsapp/wails/pull/5291)で修正

## v3.0.0-alpha.83 - 2026-05-02

## 追加

- @symball により、マシン単位またはユーザー単位のインストールに対応する InstallScope フラグとビルドオプションを[PR](https://github.com/wailsapp/wails/pull/5094)で追加
- @leaanthony により、Window インターフェースを満たすため、何も行わない SetScreen メソッドを BrowserWindow に[PR](https://github.com/wailsapp/wails/pull/5294)で追加

## 修正

- @leaanthony により、Linux で NVIDIA GPU を検出し、DMA-BUF レンダラーを無効にする処理を[PR](https://github.com/wailsapp/wails/pull/5295)で追加
- @wayneforrest により、正しいフィードバック URL を参照するよう Git の PR テンプレートを[PR](https://github.com/wailsapp/wails/pull/5109)で修正
- 1 個ではなく 4 個の引数を渡していた不正な`DestroyMenu`システムコールが原因で、すべての呼び出しが FALSE を返し、何も解放されなかったことにより発生していた、一連の Windows システムトレイ`SetMenu`クラッシュを修正。また、メニューの再構築時に HMENU および HBITMAP ハンドル（実行時に`MenuItem.SetBitmap`を介して割り当てられたものを含む）を解放し、`Win32Menu.Update`内の古いチェックボックス／ラジオボタンのマップをリセットし、割り当てを倍増させていた`systemtray.updateMenu`内の冗長な`Update()`呼び出しを削除。長時間実行されるシステムトレイアプリで、メニューを再構築するたびに GDI/USER オブジェクトがリークしなくなりました。

## v3.0.0-alpha.82 - 2026-05-01

## 修正

- @leaanthony により、デスクトップ名を正しく処理するようデスクトップファイル生成を[PR](https://github.com/wailsapp/wails/pull/5232)で修正

## v3.0.0-alpha.81 - 2026-04-30

## 変更

- @leaanthony により、ナイトリーリリースのスケジュールを15:00 UTC に[PR](https://github.com/wailsapp/wails/pull/5286)で調整

## 修正

- Retina Mac で Screen Bounds、WorkArea、Size の値が半分になる問題を修正 -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## 変更

- @leaanthony により、ドキュメントの依存関係とコンテンツコレクションローダーを[PR](https://github.com/wailsapp/wails/pull/5285)で更新

## v3.0.0-alpha.79 - 2026-04-29

## 追加

- @leaanthony により、trigger-release ジョブに actions: write 権限を[PR](https://github.com/wailsapp/wails/pull/5270)で付与

## 変更

- @leaanthony により、リリースタスクのデフォルトを master ブランチに変更し、変更履歴の文言を[PR](https://github.com/wailsapp/wails/pull/5283)で更新
- @leaanthony により、最新バージョンを使用するよう自動変更履歴ワークフローを[PR](https://github.com/wailsapp/wails/pull/5282)で更新
- @leaanthony により、パスフィルターを追加し、使用されていないワークフローを削除して、ワークフローの効率を[PR](https://github.com/wailsapp/wails/pull/5280)で改善
- @leaanthony により、サンプルのリンクで master ブランチを参照するようドキュメントを[PR](https://github.com/wailsapp/wails/pull/5274)で更新
- @leaanthony により、v3 向けのドキュメントとサンプルを[PR](https://github.com/wailsapp/wails/pull/5272)で更新

## 修正

- @AkagiYui により、開発環境向けの再試行ロジックと IPv4 強制を追加して、リバースプロキシを[PR](https://github.com/wailsapp/wails/pull/5265)で強化
- @leaanthony により、未リリース変更履歴のトリガーワークフローを[PR](https://github.com/wailsapp/wails/pull/5281)で書き直し

## 削除

- @leaanthony により、さまざまなテスト用途のシェルテストスクリプトを[PR](https://github.com/wailsapp/wails/pull/5267)で削除
- @leaanthony により、v3-alpha ドキュメントのデプロイワークフローと CNAME レコードを[PR](https://github.com/wailsapp/wails/pull/5266)で削除

### 追加

- @leaanthony により、サイドバーナビゲーションに「フロントエンドルーティング」の項目を[PR](https://github.com/wailsapp/wails/pull/5196)で追加
- @leaanthony により、フレームワークごとの推奨事項を含むフロントエンドルーティングガイドを[PR](https://github.com/wailsapp/wails/pull/5185)で追加
- モーダルシートのサポートを追加（macOS）
- @leaanthony により、Apple デバイスのサポート向上のため ghw のバージョンを更新 (#4977)
- dock サービスに`GetBadge`メソッドを追加
- カスタム Go ビルドタグ（例：`wails3 build -tags gtk4`）を渡すための`-tags`フラグを`wails3 build`コマンドに追加 (#4957)
- 専用の Enums ページとサイドバーナビゲーションを含む、バインディングジェネレーターでの enum 自動生成に関するドキュメントを追加 (#4972)
- カスタム Go ビルドタグ（例：`wails3 build -tags gtk4`）を渡すための`-tags`フラグを`wails3 build`コマンドに追加 (#4957)
- Storage（localStorage、sessionStorage、IndexedDB、Cache API）、Network（Fetch、WebSocket、XMLHttpRequest、EventSource、Beacon）、Media（Canvas、WebGL、Web Audio、MediaDevices、MediaRecorder、Speech Synthesis）、Device（Geolocation、Clipboard、Fullscreen、Device Orientation、Vibration、Gamepad）、Performance（Performance API、Mutation Observer、Intersection/Resize Observer）、UI（Web Components、Pointer Events、Selection、Dialog、Drag and Drop）などの41ブラウザー API を示す Web API サンプルを`v3/examples/web-apis/`に追加
- プラットフォーム間で200個以上のブラウザー API をテストする WebView API 互換性チェッカーのサンプル（`v3/examples/webview-api-check/`）を追加
- 並列検索、キャッシュ、および Flatpak/Snap/Nix のサポートを備え、Linux 上でネイティブライブラリのパスを検索する`internal/libpath`パッケージを追加
- **WIP:** Linux 向けの実験的な WebKitGTK 6.0 / GTK4 サポートを追加。`-tags gtk4`で利用可能（GTK3/WebKit2GTK 4.1は引き続きデフォルト）
- 注：タイル型ウィンドウマネージャー（Hyprland、Sway など）では、ウィンドウの配置とサイズを WM が制御するため、最小化／最大化の操作が期待どおりに機能しない場合があります
- @AbdelhadiSeddar により、<strong>JavaScriptでのイベントのリスニング</strong>のドキュメントに<strong>一度だけ実行するハンドラー</strong>の実装方法を追加
- @leaanthony により、Windows/Linux 上のウィンドウが`app.Menu.Set()`で設定されたアプリケーションメニューを継承できるようにする`UseApplicationMenu`オプションを`WebviewWindowOptions`に追加
- @wimaha により、Liquid Glass アイコンおよびアセットカタログ（macOS）の生成に`.icon`ファイル（Apple Icon Composer 形式）を使用するためのサポートを追加（#4934）
- ヘッドレス／Web デプロイ向けの実験的なサーバーモード（`-tags server`）を追加。ネイティブ GUI の依存関係なしで Wails アプリを HTTP サーバーとして実行できます。`wails3 task build:server`でビルドしてください。詳細は`examples/server`を参照してください。
- 並列検索、キャッシュ、および Flatpak/Snap/Nix のサポートを備え、Linux 上でネイティブライブラリのパスを検索するための`internal/libpath`パッケージを追加
- @leaanthony により、macOS の Spaces およびフルスクリーン全体でのウィンドウの動作を制御する`CollectionBehavior`オプションを`MacWindow`に追加（#4756）
- @leaanthony により pkg/application の単体テストを追加
- @leaanthony により MSIX パッケージングにカスタムプロトコルのサポートを追加
- Linux にデスクトップ環境の検出機能を追加 [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- @leaanthony により、フロントエンドから印刷ダイアログを表示するための`Window.Print()`メソッドを JavaScript ランタイムに追加（#4290）
- @leaanthony により、Linux 上の`wails3 doctor`出力に`XDG_SESSION_TYPE`を追加
- @leaanthony により、Linux 向けに WebKit2 の追加の読み込み状態変更イベント`WindowLoadStarted`、`WindowLoadRedirected`、`WindowLoadCommitted`、`WindowLoadFinished`を追加（#3896）
- @leaanthony により、Linux 上の`wails3 doctor`出力に`XDG_SESSION_TYPE`を追加
- Linux で、パッケージング時だけでなくビルド時にも`.desktop`ファイルを生成（#4575）
- @leaanthony により、ディストリビューション固有のパッケージ名と nfpm のパッケージング例を含む Linux ランタイム依存関係のドキュメントを追加（#4339）
- @leaanthony により、Linux 上の`wails3 doctor`出力に NVIDIA ドライバーのバージョン情報を追加
- [PR](https://github.com/wailsapp/wails/pull/4710)で、@APshenkin により raw メッセージハンドラーにオリジンを追加
- [PR](https://github.com/wailsapp/wails/pull/4712)で、@APshenkin により macOS のユニバーサルリンク対応を追加
- [PR](https://github.com/wailsapp/wails/pull/4702)で、@APshenkin によりバインディングのトランスポート層をリファクタリング
- [PR](https://github.com/wailsapp/wails/pull/4760)で、@chinenual により helloworld テンプレートに aria-label 識別子を追加し、Appium テストクライアントからサンプルアプリを容易にテストできるように変更
- [PR](https://github.com/wailsapp/wails/pull/4710)で、@APshenkin により raw メッセージハンドラーにオリジンを追加
- [PR](https://github.com/wailsapp/wails/pull/4712)で、@APshenkin により macOS のユニバーサルリンク対応を追加
- [PR](https://github.com/wailsapp/wails/pull/4702)で、@APshenkin によりバインディングのトランスポート層をリファクタリング
- [#4633](https://github.com/wailsapp/wails/pull/4633)で、@fbbdev と @ianvs により型付きイベントを追加
- ツールチップがリアルタイムに更新されるヘッドレストレイを示す`systray-clock`サンプルを追加（#4653）。
- @Tolfx により、#4510で Windows 向けの NSIS Protocol テンプレートを追加
- @Tolfx により、#4510で build-assets のテストを追加
- macOS：@nidib により、[#4588](https://github.com/wailsapp/wails/pull/4588)でメニューバーにネイティブのウィンドウコントロールを表示
- [PR](https://github.com/wailsapp/wails/pull/4451)で、@popaprozac により Dock 内のアプリアイコンを非表示／表示にする macOS Dock サービスを追加
- [PR](https://github.com/wailsapp/wails/pull/4451)で、@popaprozac により Dock 内のアプリアイコンを非表示／表示にする macOS Dock サービスを追加
- [#4534](https://github.com/wailsapp/wails/pull/4534)で、@leaanthony により、NSGlassEffectView（macOS 15.0以降）と NSVisualEffectView へのフォールバックを使用した macOS ネイティブの Liquid Glass エフェクト対応を追加。包括的なマテリアルのカスタマイズオプションを含む
- [#4500](https://github.dev/wailsapp/wails/pull/4500)で、@leaanthony によりブラウザー URL のサニタイズ機能を追加。@APShenkin による[#4484](https://github.com/wailsapp/wails/pull/4484)を基にしています。
- この[PR](https://github.com/wailsapp/wails/pull/4241)で、[@Taiterbase](https://github.com/Taiterbase)による元の実装を基に、[@leaanthony](https://github.com/leaanthony)が Windows/Mac のコンテンツ保護機能を追加
- [PR](https://github.com/wailsapp/wails/pull/4488)で、@leaanthony により、`wails3 build`および`wails3 package`エイリアスを介して CLI 変数を Task コマンドに渡す機能を追加（#4422）
- [#4318](https://github.com/wailsapp/wails/pull/4318)で、[@atterpac](https://github.com/atterpac)により、ドロップされた要素のデータをイベントで提供するドロップゾーンをサポート
- WebView2 ブラウザーに追加のコマンドライン引数を渡せるように、`WindowsWindow`のオプションへ`AdditionalLaunchArgs`を追加。[PR](https://github.com/wailsapp/wails/pull/4467)
- wails init の実行後に go mod tidy を自動実行する機能を追加。[PR](https://github.com/wailsapp/wails/pull/4286)で[@triadmoko](https://github.com/triadmoko)が実装
- [PR](https://github.dev/wailsapp/wails/pull/4463)で、@leaanthony により Windows Snapassist 機能を追加
- WebView2 ブラウザーに追加のコマンドライン引数を渡せるように、`WindowsWindow`のオプションへ`AdditionalLaunchArgs`を追加。[PR](https://github.com/wailsapp/wails/pull/4467)
- wails init の実行後に go mod tidy を自動実行する機能を追加。[PR](https://github.com/wailsapp/wails/pull/4286)で[@triadmoko](https://github.com/triadmoko)が実装
- [PR](https://github.dev/wailsapp/wails/pull/4463)で、@leaanthony により Windows Snapassist 機能を追加
- [PR](https://github.com/wailsapp/wails/pull/4427)で、[@almas-x](https://github.com/almas-x)により Windows 向けの`getAccentColor`実装を追加
- [PR](https://github.com/wailsapp/wails/pull/4427)で、[@almas-x](https://github.com/almas-x)により Windows 向けの`getAccentColor`実装を追加
- Windows のダークテーマ対応メニューおよびメニューバー。[a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)で @leaanthony が実装
- [PR](https://github.com/wailsapp/wails/pull/4405)で、@popaprozac により、JS/TS バインディングをより明確にするため組み込みサービスの名前を変更
- ユーザーのシステムのアクセントカラーを取得する`app.Env.GetAccentColor`。MacOS で動作します。[@etesam913](https://github.com/etesam913)が実装
- [#4137](https://github.com/wailsapp/wails/pull/4137)で、[@atterpac](https://github.com/atterpac)により`window.ToggleFrameless()` API を追加
- @leaanthony により、Linux 向けにディストリビューション固有のビルド依存関係を追加（[PR](https://github.com/wailsapp/wails/pull/4345)）
- @atterpac により、バインディングガイドを追加（[PR](https://github.com/wailsapp/wails/pull/4404)）
- **テストインフラストラクチャの整理**：Docker のテストファイルを専用の `test/docker/` ディレクトリへ移動し、イメージを最適化してビルドの信頼性を向上。[@leaanthony](https://github.com/leaanthony) による変更（[#4359](https://github.com/wailsapp/wails/pull/4359)）
- **リソース管理パターンの改善**：サンプルに、イベントハンドラーの適切なクリーンアップとコンテキストを考慮した goroutine 管理を追加。[@leaanthony](https://github.com/leaanthony) による変更（[#4359](https://github.com/wailsapp/wails/pull/4359)）
- [@AkshayKalose](https://github.com/AkshayKalose) により、aarch64 AppImage のビルドをサポート（[#3981](https://github.com/wailsapp/wails/pull/3981)）
- [@leaanthony](https://github.com/leaanthony) により、`wails doctor` に診断セクションを追加
- [@leaanthony](https://github.com/leaanthony) により、サービスメソッドの呼び出し時にウィンドウをコンテキストへ追加
- [@leaanthony](https://github.com/leaanthony) により、どのウィンドウがサービスを呼び出しているかを確認する方法を示す `window-call` サンプルを追加
- [@leaanthony](https://github.com/leaanthony) による新しい Menu ガイド
- [@leaanthony](https://github.com/leaanthony) により、panic 処理を改善
- [@leaanthony](https://github.com/leaanthony) による新しい Menu ガイド
- [@fbbdev](https://github.com/fbbdev) により、Service API のドキュメントコメントを追加（[#4024](https://github.com/wailsapp/wails/pull/4024)）
- [@leaanthony](https://github.com/leaanthony) により、追加設定を使用してサービスを初期化する関数 `application.NewServiceWithOptions` を追加（[#4024](https://github.com/wailsapp/wails/pull/4024)）
- [@FalcoG](https://github.com/FalcoG) と [@leaanthony](https://github.com/leaanthony) により、メニュー制御を改善（[#4031](https://github.com/wailsapp/wails/pull/4031)）
- [@leaanthony](https://github.com/leaanthony) により、ドキュメントをさらに追加
- [@leaanthony](https://github.com/leaanthony) により、標準イベントリスナーでのイベントキャンセルをサポート
- [@leaanthony](https://github.com/leaanthony) により、Systray の `Hide`、`Show`、`Destroy` をサポート
- [@leaanthony](https://github.com/leaanthony) により、Systray の `SetTooltip` をサポート。原案は [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- [@fbbdev](https://github.com/fbbdev) により、未対応の型に関するバインディングジェネレーターの警告でパッケージパスを報告（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- [@fbbdev](https://github.com/fbbdev) により、バインディングジェネレーターでジェネリックエイリアスをサポート（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- [@fbbdev](https://github.com/fbbdev) により、バインディングジェネレーターで `omitzero` JSON フラグをサポート（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- [@fbbdev](https://github.com/fbbdev) により、選択したサービスメソッドのバインディング生成を防止する `//wails:ignore` ディレクティブを追加（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- [@fbbdev](https://github.com/fbbdev) により、Go ではエクスポートされるが JS/TS ではエクスポートされない型を許可する `//wails:internal` ディレクティブをサービスとモデルに追加（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- [@fbbdev](https://github.com/fbbdev) により、弱い型付けの列挙型を使用できるよう、バインディングジェネレーターでエイリアス型の定数をサポート（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- [@fbbdev](https://github.com/fbbdev) により、Go 1.24 の機能に対するバインディングジェネレーターのテストを追加（[#4068](https://github.com/wailsapp/wails/pull/4068)）
- OS バージョン検出を改善するため、`OSInfo.Branding` に macOS 15「Sequoia」のサポートを追加（[#4065](https://github.com/wailsapp/wails/pull/4065)）
- [@fbbdev](https://github.com/fbbdev) により、シャットダウン処理の完了後にカスタムコードを実行する `PostShutdown` フックを追加（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- [@fbbdev](https://github.com/fbbdev) により、カスタムエラーハンドラーで致命的エラーを検出できるようにする `FatalError` 構造体を追加（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- [@fbbdev](https://github.com/fbbdev) により、サービスの起動順序とシャットダウン順序を標準化し、文書化（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- [@fbbdev](https://github.com/fbbdev) により、アプリケーションの起動／シャットダウンシーケンス用テストハーネスと、サービスの起動／シャットダウンテストを追加（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- [@fbbdev](https://github.com/fbbdev) により、アプリケーション作成後にサービスを登録する `RegisterService` メソッドを追加（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- [@fbbdev](https://github.com/fbbdev) により、バインディング呼び出しのカスタムエラー処理用 `MarshalError` フィールドをアプリケーションとサービスのオプションに追加（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- [@fbbdev](https://github.com/fbbdev) により、キャンセル要求を Promise チェーン全体へ伝播するキャンセル可能な Promise ラッパーを追加（[#4100](https://github.com/wailsapp/wails/pull/4100)）
- [@fbbdev](https://github.com/fbbdev) により、バインディング呼び出しのキャンセルを `AbortSignal` に関連付ける機能を追加（[#4100](https://github.com/wailsapp/wails/pull/4100)）
- [@leaanthony](https://github.com/leaanthony) により、WML で通常の `wml-*` 属性に加えて `data-wml-*` 属性をサポート
- [@fbbdev](https://github.com/fbbdev) により、後からの設定／動的な再設定に使用する `Configure` メソッドをすべてのサービスに追加（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- `fileserver` サービスが未設定の場合、503 Service Unavailable レスポンスを返すように変更。[@fbbdev](https://github.com/fbbdev) による変更（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- `kvstore` サービスが未設定の場合、デフォルトでインメモリのキー値ストアを提供するように変更。[@fbbdev](https://github.com/fbbdev) による変更（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) により、設定変更後にファイルからデータを再読み込みする `Load` メソッドを `kvstore` サービスに追加（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) により、すべてのキーを削除する `Clear` メソッドを `kvstore` サービスに追加（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) により、JS 側にログレベル定数を提供する型 `Level` を `log` サービスに追加（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) により、ログレベルを動的に指定する `Log` メソッドを `log` サービスに追加（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- `sqlite` サービスが未設定の場合、デフォルトでインメモリ DB を提供するように変更。[@fbbdev](https://github.com/fbbdev) による変更（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- `sqlite` サービスに、DB を手動で閉じるための `Close` メソッドを追加。[#4067](https://github.com/wailsapp/wails/pull/4067) で [@fbbdev](https://github.com/fbbdev) が実装
- `sqlite` サービスのクエリメソッドにキャンセル対応を追加。[#4067](https://github.com/wailsapp/wails/pull/4067) で [@fbbdev](https://github.com/fbbdev) が実装
- `sqlite` サービスに、JS バインディングを備えたプリペアドステートメント対応を追加。[#4067](https://github.com/wailsapp/wails/pull/4067) で [@fbbdev](https://github.com/fbbdev) が実装
- Gin に対応。[@AnalogJ](https://github.com/AnalogJ) によるこちらの [PR](https://github.com/wailsapp/wails/pull/3537) の元の成果を基に、[PR](https://github.com/wailsapp/wails/pull/3537) で [Lea Anthony](https://github.com/leaanthony) が実装
- 自動保存とパスワードの自動保存が常に有効になる問題を修正。[#4134](https://github.com/wailsapp/wails/pull/4134) で [@oSethoum](https://github.com/osethoum) が実装
- ウィンドウにメニューを設定できるよう、ウィンドウに `SetMenu()` を追加。[@leaanthony](https://github.com/leaanthony) が実装
- 通知に対応。[#4098](https://github.com/wailsapp/wails/pull/4098) で [@popaprozac](https://github.com/popaprozac) が実装
-  mac のファイル関連付けに対応。[#4177](https://github.com/wailsapp/wails/pull/4177) で [@wimaha](https://github.com/wimaha) が実装
- セマンティックバージョンを更新するための `wails3 tool version` を追加。[@leaanthony](https://github.com/leaanthony) が実装
- macOS と Windows のバッジ表示に対応。[#](https://github.com/wailsapp/wails/pull/4234) で [@popaprozac](https://github.com/popaprozac) が実装
- 登録済みイベント／厳密に型付けされたイベントに対応。[#4161](https://github.com/wailsapp/wails/pull/4161) で [@fbbdev](https://github.com/fbbdev) と [@IanVS](https://github.com/IanVS) が実装
- カスタムイベント用のフックを登録できる機能を追加。[#4161](https://github.com/wailsapp/wails/pull/4161) で [@fbbdev](https://github.com/fbbdev) と [@IanVS](https://github.com/IanVS) が実装
- `path` のパスをシステムのファイルマネージャーで開くための `app.OpenFileManager(path string, selectFile bool)` を追加。`selectFile` による任意の強調表示にも対応。[@Krzysztofz01](https://github.com/Krzysztofz01)、[@rcalixte](https://github.com/rcalixte) が実装
- `wails3 init` コマンドに新しい `-git` フラグを追加。[@leaanthony](https://github.com/leaanthony) が実装
- 新しい `wails3 generate webview2bootstrapper` コマンドを追加。[@leaanthony](https://github.com/leaanthony) が実装
- ランタイムを手動で初期化できるよう、ランタイムに `init()` メソッドを追加。[@leaanthony](https://github.com/leaanthony) が実装
- Window の WindowOptions に `WindowDidMoveDebounceMS` オプションを追加。[@leaanthony](https://github.com/leaanthony) が実装
- シングルインスタンス機能を追加。[@leaanthony](https://github.com/leaanthony) が実装。@APshenkin による [v2 PR](https://github.com/wailsapp/wails/pull/2951) が基になっています。
- `wails3 generate template` コマンドを追加。[@leaanthony](https://github.com/leaanthony) が実装
- `wails3 releasenotes` コマンドを追加。[@leaanthony](https://github.com/leaanthony) が実装
- `wails3 update cli` コマンドを追加。[@leaanthony](https://github.com/leaanthony) が実装
- `wails3 generate bindings` コマンドに `-clean` オプションを追加。[@leaanthony](https://github.com/leaanthony) が実装
- aarch64（arm64）向け AppImage の Linux ビルドに対応。[#3981](https://github.com/wailsapp/wails/pull/3981) で [@AkshayKalose](https://github.com/AkshayKalose) が実装
- スポンサーへのハイパーリンクを追加。[#3958](https://github.com/wailsapp/wails/pull/3958) で @ansxuman が実装
- deb、rpm、Arch Linux 形式の Linux パッケージのビルドに対応。実装者：
- Darwin のユニバーサルビルドおよびパッケージに対応。実装者：
- Web サイトにイベントのドキュメントを追加。実装者：
- 非 SSR 開発用に設定された sveltekit および sveltekit-ts のテンプレート
- 新しい `wails3 update build-assets` コマンドを使用してビルドアセットを更新。実装者：
- HTML Drag and Drop API をテストするサンプル。実装者：
- ファイル関連付けに対応。実装者は [leaanthony](https://github.com/leaanthony)、対象：
- 新しい `wails3 generate runtime` コマンド。実装者：
- 新しい `InitialPosition` オプションで、ウィンドウを中央に配置するか、または
- `application` パッケージに `Path` および `Paths` メソッドを追加。実装者：
- Windows の `GeneralAutofillEnabled` および `PasswordAutosaveEnabled` オプションを追加
- サービスメソッドを呼び出したウィンドウを取得できる機能を追加。実装者：
- WebView2 に `EnabledFeatures` および `DisabledFeatures` オプションを追加。実装者：
- ⊞ 強化された高 DPI モニター対応のための新しい DIP システム。実装者：
- ⊞ ウィンドウクラス名のオプション。実装者は [windom](https://github.com/windom/)、対象：
- プラグイン機能を提供するようにサービスを拡張。実装者：
- 🐧 WindowDidMove / WindowDidResize イベント。対象：
- ⊞ WindowDidResize イベント。対象：
-  Dock を処理できるように Event ApplicationShouldHandleReopen を追加
-  実装に getPrimaryScreen/getScreens を追加。@tmclane が次で実装：
-  macOS のフルスクリーンモードでツールバーを表示するためのオプションを追加。実装者：
- 🐧 Linux のキー入力をアクセラレーターに変換する onKeyPress ロジックを追加
- 🐧 タスク `run:linux` を追加。実装者：
- `SetIcon` メソッドをエクスポート。[@almas-x](https://github.com/almas-x) が次で実装：
- `OnShutdown` を改善。[@almas-x](https://github.com/almas-x) が次で実装：
- `Window` インターフェースの `ToggleMaximise` メソッドを復元。実装者：
- `Environment()` に詳細情報を追加。@leaanthony が次で実装：
- `Window` インターフェースで `WebviewWindow.IsFocused` メソッドを公開。実装者：
- WML システムで、スペース区切りの複数のトリガーイベントに対応。実装者：
- バンドル済み JS ランタイムスクリプトに ESM エクスポートを追加。実装者：
- 代わりにバンドル済み JS ランタイムスクリプトを使用するためのバインディングジェネレーターフラグを追加
- Linux に `setIcon` を実装。[@abichinger](https://github.com/abichinger) が実装
- dev コマンドに `-port` フラグを追加し、環境変数に対応
- バインドされたメソッド呼び出しのテストを追加。担当：
- ⊞ 作成済みのウィンドウに `SetIgnoreMouseEvents` を追加。担当：
-  ウィンドウの重なりレベル（順序）を設定する機能を追加。担当：

### 修正

- Retina Mac で `Screen.Bounds`、`WorkArea`、`Size` が半分の値になる問題を修正。NSScreen のポイント値を `Physical*` フィールド内でデバイスピクセルに変換し、最上位の `Screen.X`/`Y` に値を設定することで、マルチモニター間の接触判定と作業領域への配置が正しくなるようにした（[PR](https://github.com/wailsapp/wails/pull/5168)、@wayneforrest）
- ディスプレイ構成の変更時（例：スリープ／復帰中の外部モニターのホットプラグ）に WebKit DisplayLink のデッドロックを引き起こす ScreenManager のデータ競合を修正
- Assets.car が存在する場合、CFBundleIconName を appicon に直接設定（[PR](https://github.com/wailsapp/wails/pull/5154)、@symball）
- Fedora、openSUSE、Arch、NixOS で `wails3 doctor` が誤った WebKitGTK パッケージを報告する問題を修正。v3 ではコンパイル時に 4.1 API が必要なため、4.0 のフォールバックエントリを削除（#5071）
- openSUSE の webkit2gtk に対して doctor が示すパッケージ名を修正（`webkit2gtk4_1-devel` → 正しい openSUSE パッケージ名である `webkit2gtk3-devel`）（#5071）
- デスクトップ開発モードで `/wails/custom.js` がない場合に発生する `Unexpected token '<'` エラーを修正。HTML SPA のフォールバックが JavaScript として挿入されないように、`/wails/custom.js` 用の明示的な 404 ハンドラーと、大文字と小文字を区別しない `Content-Type` 検証を `loadOptionalScript` に追加（[#5068](https://github.com/wailsapp/wails/issues/5068)）
- macOS のシステムトレイメニューのハイライト状態を修正。メニューが開いているときに、アイコンが選択状態で表示されるようにした（#4910）
- macOS でシステムトレイに関連付けられたウィンドウがほかのウィンドウの背後に表示される問題を修正。適切なポップアップウィンドウレベルを使用するようにした（#4910）
- ドキュメント全体にある誤った `@wailsio/runtime` インポート例を修正（#4989）
- darwin でフレームレスウィンドウを最小化できない問題を修正（#4294）
- `node_modules/` を go-task の最新状態チェックから除外し、`wails3 build` と `wails3 dev` の実行中に 20-30 分間応答しなくなる問題を修正。従来は `sources: "**/*"` の glob パターンにより、go-task が `node_modules/` 内の全ファイル（MUI のような大規模な依存関係がある場合は 50000-100000 件以上）を列挙してチェックサムを計算していたため、特に Windows/NTFS で低速だった（#4939）
- C の `Screen` typedef と X11 の Xlib.h の衝突による GTK4 のビルド失敗を修正（#4957）
- macOS の Dock バッジメソッド間の一貫性を修正
- `InvisibleTitleBarHeight` がフレームレスウィンドウまたはタイトルバーが透明なウィンドウだけでなく、すべての macOS ウィンドウに適用される問題を修正（#4960）
- `InvisibleTitleBarHeight` が有効な状態で上隅からサイズを変更するとウィンドウが揺れる問題を、ウィンドウ端付近ではドラッグを開始しないことで修正（#4960）
- JS/TS バインディングで enum キーを持つマップ型が正しく生成されない問題を修正（#4437、@fbbdev）
- Windows でディスプレイのスケーリングが 100% 以外の場合に、ファイルのドラッグ＆ドロップが機能しない問題を修正
- Windows でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- Windows でファイルのドロップ座標が誤ったピクセル空間（物理ピクセルと CSS ピクセル）になる問題を修正
- Linux でホバー効果を伴うファイルのドラッグ＆ドロップが安定して機能しない問題を修正
- Linux でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- `gtk_window_present()` を使用し、Linux/GTK4 でウィンドウの表示／非表示を切り替えた際に、最小化された状態へ復元されることがある問題を修正（#4957）
- `XTranslateCoordinates`/`XMoveWindow` を介した X11 条件付きサポートを追加し、Linux/GTK4 でウィンドウ位置の取得／設定が常に 0,0 を返す問題を修正（#4957）
- 削除された `gtk_window_set_geometry_hints` の代わりに、シグナルベースでサイズを上限に制限する処理を追加し、Linux/GTK4 でウィンドウの最大サイズが適用されない問題を修正（#4957）
- PhysicalBounds の適切な計算と、`gdk_monitor_get_scale` を介した小数スケーリングのサポートを実装し、Linux/GTK4 の DPI スケーリングを修正（GTK 4.14 以降）
- Linux/GTK4 で新しいウィンドウを作成するとメニュー項目が重複する問題を修正
- JS/TS バインディングで enum キーを持つマップ型が正しく生成されない問題を修正（#4437、@fbbdev）
- Windows でディスプレイのスケーリングが 100% 以外の場合に、ファイルのドラッグ＆ドロップが機能しない問題を修正
- Windows でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- Windows でファイルのドロップ座標が誤ったピクセル空間（物理ピクセルと CSS ピクセル）になる問題を修正
- Linux でホバー効果を伴うファイルのドラッグ＆ドロップが安定して機能しない問題を修正
- Linux でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- PhysicalBounds の適切な計算と、`gdk_monitor_get_scale` を介した小数スケーリングのサポートを実装し、Linux/GTK4 の DPI スケーリングを修正（GTK 4.14 以降）
- Linux/GTK4 で新しいウィンドウを作成するとメニュー項目が重複する問題を修正
- JS/TS バインディングで enum キーを持つマップ型が正しく生成されない問題を修正（#4437、@fbbdev）
- App.Window.Current() で Main Thread から AppKit API にアクセスしていなかったために macOS で発生する「ゴーストウィンドウ」の問題を修正（#4947、@wimaha）
- WKUIDelegate runOpenPanelWithParameters を実装し、macOS で HTML の `<input type="file">` が機能しない問題を修正（#4862）
- macOS/Linux で `@wailsio/runtime` npm モジュールを使用している場合に、ネイティブのファイルドラッグ＆ドロップが機能しない問題を修正（#4953、@leaanthony）
- パッケージをまたぐ型エイリアスのバインディング生成を修正（#4578、@fbbdev）
- GTK のスレッドセーフティ違反によって Linux で OpenFileDialog がクラッシュする問題を修正（#3683、@ddmoney420）
- 非表示または破棄済みのウィンドウで `Focus()` を呼び出した際に SIGSEGV でクラッシュする問題を修正（#4890、@ddmoney420）
- Linux で空のアイコンまたはビットマップを設定した際に panic が発生する可能性がある問題を修正（#4923、@ddmoney420）
- macOS でサービスバインディングから ErrorDialog を呼び出した際にクラッシュする問題を修正（#3631、@leaanthony）
- `v3\examples\dialogs` で Windows OS 上にメニューが表示されるように変更（@ndianabasi）
- ページの再読み込み中に TypeError を引き起こす競合状態を修正（#4872、@ddmoney420）
- `Collector.IsVoidAlias()` メソッドからグローバル状態を削除し、バインディングジェネレーターのテストで誤った出力が生成される問題を修正（#4941、@fbbdev）
- macOS で `<input type="file">` ファイルピッカーが機能しない問題を修正（#4862、@leaanthony）
- macOSで`Position()`と`SetPosition()`が異なる座標系を使用し、状態の保存／復元時にウィンドウ位置がずれる問題を修正（#4816、@leaanthony）
- アプリケーションマニフェストですでにDPI認識が設定されている場合に発生するSetProcessDpiAwarenessContextの「Access is denied」エラーを修正（#4803）
- キーボードショートカットのドキュメントページを更新し、`KeyBinding.Add`のコールバックパラメーターの型を修正（@ndianabasi）
- カスタムバインディングの生成に関するドキュメントを修正。`-o String`ではなく`-d String`を使用する必要があります
- `menu.Update()`でメニューの子項目がクリアされない問題を修正
- ドキュメント内の古いManager API参照を修正（31ファイルを更新し、`app.Window.New()`、`app.Event.Emit()`などの新しいパターンを使用）（@leaanthony）
- WebKitによるシグナルハンドラーの上書きが原因で、JSにバインドされたGoメソッドのpanic時にLinuxでクラッシュする問題を修正（#3965、@leaanthony）
- LinuxでSaveFileDialog.SetFilename()が機能しない問題を修正（#4841、@samstanier）
- ドラッグ＆ドロップの例でドロップ座標がundefinedと表示される問題を修正
- APP_NAMEに空白が含まれる場合、ブレース展開の問題によりmacOSアプリバンドルの作成に失敗する問題を修正
- サービスメソッドの呼び出し時にWindowsで発生するインデックス範囲外のpanicを修正（goccy/go-jsonを元に戻しました）
- Windowsでディスプレイの拡大率が100%以外の場合にファイルのドラッグ＆ドロップが機能しない問題を修正
- Windowsでファイルドロップを有効にすると、HTML5の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- Windowsでファイルのドロップ座標に誤ったピクセル空間（物理ピクセルとCSSピクセル）が使用される問題を修正
- Linuxでホバー効果を使用すると、ファイルのドラッグ＆ドロップが安定して機能しない問題を修正
- Linuxでファイルドロップを有効にすると、HTML5の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- `APP_NAME`などの変数に含まれる空白に対応するため、すべてのオペレーティングシステム向けのTaskfile.ymlファイル内の全コマンドを更新（@ndianabasi）
- Linuxで'build:universal:lipo:go'タスクを実行する際のコマンド引数エラーを修正（@wux1an）
- Linuxで'wails3 build GOOS=darwin GOARCH=arm64'を実行した際に発生するDockerエラー「undefined symbol: **<em>ubsan</em>handle_xxxxxxx」を修正（@wux1an）
- カスタムプロトコルのドキュメントを統合し、Universal Linksのセクションを追加（@leaanthony）
- TrackPopupMenuExの同時呼び出しを防ぐガードを追加し、アイコンを繰り返しクリックした際にWindowsのシステムトレイメニューがクラッシュする問題を修正（#4151、@leaanthony）
- app.Run()より前にsystray.Run()を呼び出した際にアプリがクラッシュするのを防止（@leaanthony）
- ApplicationShouldTerminateAfterLastWindowClosedが有効な状態でHide()/Show()によりウィンドウの表示／非表示を切り替えると、macOSでクラッシュする問題を修正（#4389、@leaanthony）
- macOSおよびWindowsでメニューを繰り返し開いた際に発生するコンテキストメニューのメモリリークを修正（#4012、@leaanthony）
- macOSでコンテキストメニューのネイティブリソースが再利用されず、表示のたびに新しいメニューが作成される問題を修正（#4012、@leaanthony）
- `Hidden: true`を指定してアプリを起動した場合、macOSのDockアイコンをクリックしても非表示のウィンドウが表示されない問題を修正（#4583、@leaanthony）
- CGO呼び出しでウィンドウポインターの型が誤っていたため、macOSで印刷ダイアログが開かない問題を修正（#4290、@leaanthony）
- appmenu-gtk-moduleが未実体化のウィンドウにアクセスすることで、Wayland上のウィンドウメニューがクラッシュする問題を修正（#4769、@leaanthony）
- アプリ名に無効な文字（空白、丸括弧など）が含まれる場合にGTKアプリケーションがクラッシュする問題を修正（@leaanthony）
- Windowsでドラッグ＆ドロップの初期化時に発生する「not enough memory」エラーを修正（#4701、@overlordtm）
- URIのエスケープが誤っていたため、Linuxでファイルエクスプローラーが別のディレクトリを開く問題を修正（#4397、@leaanthony）
- `.relr.dyn` ELFセクションを自動検出してストリップを無効にすることで、最新のLinuxディストリビューション（Arch、Fedora 39以降、Ubuntu 24.04以降）でAppImageのビルドに失敗する問題を修正（#4642、@leaanthony）
- Fedora/DNFベースのシステムで、`wails doctor`がwebkitパッケージをインストール済みと誤って報告する問題を修正（#4457、@leaanthony）
- デフォルトの`config.yml`がプロダクションビルドで`wails3 dev`を実行してしまう問題を修正（@mbaklor）
- 存在しないパッケージのインポートによりiOSサービススタブでビルドが失敗する問題を修正（@leaanthony）
- debug/infoメソッドの構造化ロギングで「no formatting directives」エラーが発生する問題を修正（@leaanthony）
- モバイルプラットフォームのマージから誤って含まれた一時的なデバッグ用print文を削除（@leaanthony）
- DMA-BUFレンダラーを自動的に無効化し、NVIDIA GPUを使用するWayland環境でWebKitGTKがクラッシュする問題（Error 71 Protocol error）を修正（@leaanthony）
- Linuxで`application.WebviewWindowOptions.BackgroundColour`のアルファ値が無視される問題を解決（[#4722](https://github.com/wailsapp/wails/pull/4722)、@BradHacker）
- カスタムアイコンを指定していない場合、Windowsのシステムトレイアイコンにアプリケーションアイコンがデフォルトで使用されない問題を修正（#4704）
- `HICON`の所有権を追跡し、ユーザーが作成したハンドルのみを破棄することで、Explorerの再起動時のクラッシュを防止（#4653）。
- 破棄時にWindowsのシステムテーマリスナーと保持中のトレイアイコンを解放し、goroutineとデバイスコンテキストのリークを解消（#4653）。
- サロゲートペアやマルチバイトグリフの破損を防ぐため、トレイのツールチップを127 UTF-16単位で切り詰めるように変更（#4653）。
- Windowsのパッケージタスクが失敗する問題を修正（#4667）
- LinuxのtaskfileにあるLinux AppImageのappicon変数を修正（[PR #4644](https://github.com/wailsapp/wails/pull/4644)）
- go-webview2 v1.0.22のシグネチャ変更により発生するWindowsのビルドエラーを修正（#4513、#4645）
- LinuxのtaskfileにあるLinux AppImageのappicon変数を修正（[PR #4644](https://github.com/wailsapp/wails/pull/4644)）
- `<.Info.Protocol>`を削除して`<.Protocol>`に変更し、Linuxのdesktop.tmplにあるプロトコルの反復処理を修正（@Tolfx、#4510）
- liquid glassデモの再定義エラーを修正（[#4542](https://github.com/wailsapp/wails/pull/4542)、@Etesam913）
- Linuxでのシステムトレイメニューの更新を修正（[#4604](https://github.com/wailsapp/wails/issues/4604)、[@JackDoan](https://github.com/JackDoan)）
- Windowsで非表示ウィンドウを作成した際に白いウィンドウが表示される問題を修正（@leaanthony、[#4612](https://github.com/wailsapp/wails/pull/4612)）
- ドキュメント内のnotificationsパッケージのインポートパスを修正（@rxliuli、[#4617](https://github.com/wailsapp/wails/pull/4617)）
- npm パッケージ @wailsio/runtime の使用時にドラッグ＆ドロップが機能しない問題を修正（#4489）。@leaanthony が #4616 で対応
- Windows：起動時にウィンドウがちらつく問題と、非表示ウィンドウが誤って表示される問題を修正。@leaanthony による [PR](https://github.com/wailsapp/wails/pull/4600)
- Wayland でウィンドウを最大化した際のサイズ問題（https://github.com/wailsapp/wails/issues/4429）を修正。[@samstanier](https://github.com/samstanier) が対応
- Wayland でウィンドウを最大化した際のサイズ問題（https://github.com/wailsapp/wails/issues/4429）を修正。[@samstanier](https://github.com/samstanier) が対応
- リキッドガラスのデモで発生する再定義エラーを修正。@Etesam913 による [#4542](https://github.com/wailsapp/wails/pull/4542)
- MacOS で AssetServer がクラッシュする可能性がある問題を修正。@jghiloni による [#4576](https://github.com/wailsapp/wails/pull/4576)
- NextJs でのビルド時に発生するコンパイル問題を修正。@rev42 が [#4585](https://github.com/wailsapp/wails/pull/4585) で対応
- ナイトリーリリース用パイプラインを修正。@riadafridishibly による [#4597](https://github.com/wailsapp/wails/pull/4597)
- リキッドガラスのデモで発生する再定義エラーを修正。@Etesam913 による [#4542](https://github.com/wailsapp/wails/pull/4542)
- MacOS で AssetServer がクラッシュする可能性がある問題を修正。@jghiloni による [#4576](https://github.com/wailsapp/wails/pull/4576)
- NextJs でのビルド時に発生するコンパイル問題を修正。@rev42 が [#4585](https://github.com/wailsapp/wails/pull/4585) で対応
- ナイトリーリリース用パイプラインを修正。@riadafridishibly による [#4597](https://github.com/wailsapp/wails/pull/4597)
- リキッドガラスのデモで発生する再定義エラーを修正。@Etesam913 による [#4542](https://github.com/wailsapp/wails/pull/4542)
- Windows の SetBackgroundColour を修正。@PPTGamer による [PR](https://github.com/wailsapp/wails/pull/4492)
- Manager API のリファクタリングによる変更を反映するようドキュメントを更新。@yulesxoxo による [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Linux の Taskfile にある .desktop ファイル用 appicon 変数を修正。[PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Manager API のリファクタリングによる変更を反映するようドキュメントを更新。@yulesxoxo による [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- [#4456](https://github.com/wailsapp/wails/issues/4456) で報告された Windows の nil ポインター参照バグを修正。@leaanthony による [#4460](https://github.com/wailsapp/wails/pull/4460)
- 2 本指のスワイプによるナビゲーションジェスチャーを有効にするため、macOS の WKWebView で `allowsBackForwardNavigationGestures` をサポート（#1857）
- 初期状態で無効に設定されたメニュー項目で onClick が機能しない問題を修正。@leaanthony による [PR #4469](https://github.com/wailsapp/wails/pull/4469)。初期調査を行った @IanVS に感謝します。
- ビルド失敗時に Vite サーバーがクリーンアップされない問題を修正（#4403）
- Windows で `SaveFileDialog` を閉じるかキャンセルした際に発生するパニックを修正。@hkhere が [PR](https://github.com/wailsapp/wails/pull/4284) で対応
- Windows で HTML レベルのドラッグ＆ドロップを修正。[@mbaklor](https://github.com/mbaklor) による [#4259](https://github.com/wailsapp/wails/pull/4259)
- 2 本指のスワイプによるナビゲーションジェスチャーを有効にするため、macOS の WKWebView で `allowsBackForwardNavigationGestures` をサポート（#1857）
- 初期状態で無効に設定されたメニュー項目で onClick が機能しない問題を修正。@leaanthony による [PR #4469](https://github.com/wailsapp/wails/pull/4469)。初期調査を行った @IanVS に感謝します。
- ビルド失敗時に Vite サーバーがクリーンアップされない問題を修正（#4403）
- Windows での通知の解析処理を修正。@popaprozac による [PR](https://github.com/wailsapp/wails/pull/4450)
- Windows SDK の依存関係を確認するよう doctor コマンドを修正。[@kodumulo](https://github.com/kodumulo) による [#4390](https://github.com/wailsapp/wails/issues/4390)
- Mac の processURLRequest における nil ポインター参照を修正。[@etesam913](https://github.com/etesam913) による [#4366](https://github.com/wailsapp/wails/pull/4366)
- Linux でフィルター付きダイアログを使用できないバグを修正。[@bh90210](https://github.com/bh90210) による [#4287](https://github.com/wailsapp/wails/pull/4287)
- Windows および Linux の編集メニューに関する問題を修正。[@leaanthony](https://github.com/leaanthony) による [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- macOS の .plist ファイルにある最小システムバージョンを 10.13.0 から 10.15.0 に更新。[@AkshayKalose](https://github.com/AkshayKalose) による [#3981](https://github.com/wailsapp/wails/pull/3981)
- ウィンドウ ID が飛ばされる問題を修正。[@leaanthony](https://github.com/leaanthony) が対応
- RegisterContextMenu の呼び出し時にメニューが nil になる問題を修正。[@leaanthony](https://github.com/leaanthony) が対応
- バインディングジェネレーターの出力における依存関係の循環を修正。[@fbbdev](https://github.com/fbbdev) による [#4001](https://github.com/wailsapp/wails/pull/4001)
- バインディングジェネレーターの出力における定義前使用エラーを修正。[@fbbdev](https://github.com/fbbdev) による [#4001](https://github.com/wailsapp/wails/pull/4001)
- ビルドフラグをバインディングジェネレーターに渡すよう変更。[@fbbdev](https://github.com/fbbdev) による [#4023](https://github.com/wailsapp/wails/pull/4023)
- Windows 以外のプラットフォームでも動作するよう、windows Taskfile 内のパスをスラッシュ区切りに変更。[@leaanthony](https://github.com/leaanthony) が対応
- Mac および Mac の JS イベントを修正。[@leaanthony](https://github.com/leaanthony) が対応
- macOS でのイベントのデッドロックを修正。[@leaanthony](https://github.com/leaanthony) が対応
- Windows で HTML は指定されているものの JS が指定されていない場合に、Window の初期化で発生する `Parameter incorrect` エラーを修正。[@leaanthony](https://github.com/leaanthony) が対応
- アセットサーバーでコンテンツタイプのスニッフィングに使用するレスポンスのプレフィックスサイズを修正。[@fbbdev](https://github.com/fbbdev) による [#4049](https://github.com/wailsapp/wails/pull/4049)
- アセットサーバーのルートインデックスパスで、404 以外のレスポンスの処理を修正。[@fbbdev](https://github.com/fbbdev) による [#4049](https://github.com/wailsapp/wails/pull/4049)
- バインディングジェネレーターでジェネリック型のプロパティを検査する際の未定義動作を修正。[@fbbdev](https://github.com/fbbdev) による [#4045](https://github.com/wailsapp/wails/pull/4045)
- 基になる型のプロパティが名前付きラッパーと同一でない場合の、モデルに対するバインディングジェネレーターの出力を修正。[@fbbdev](https://github.com/fbbdev) による [#4045](https://github.com/wailsapp/wails/pull/4045)
- マップのキー型および前処理に対するバインディングジェネレーターの出力を修正。[@fbbdev](https://github.com/fbbdev) による [#4045](https://github.com/wailsapp/wails/pull/4045)
- マーシャラーインターフェースを実装する構造体に対するバインディングジェネレーターの出力を修正。[@fbbdev](https://github.com/fbbdev) による [#4045](https://github.com/wailsapp/wails/pull/4045)
- バインディングジェネレーターで、ジェネリック型を含む型の循環検出を修正。[@fbbdev](https://github.com/fbbdev) による [#4045](https://github.com/wailsapp/wails/pull/4045)
- バインディングジェネレーターの出力に含まれる、エクスポートされていないモデルへの無効な参照を修正（[@fbbdev](https://github.com/fbbdev)、[#4045](https://github.com/wailsapp/wails/pull/4045)）
- 注入されるコードをサービスファイルの末尾に移動（[@fbbdev](https://github.com/fbbdev)、[#4045](https://github.com/wailsapp/wails/pull/4045)）
- バインディングジェネレーターで、ファイルを閉じる処理から返されるエラーの扱いを修正（[@fbbdev](https://github.com/fbbdev)、[#4045](https://github.com/wailsapp/wails/pull/4045)）
- ライフサイクルメソッドまたは http メソッドを定義しているものの、それ以外にバインドされたメソッドがないサービスに対する警告を抑制（[@fbbdev](https://github.com/fbbdev)、[#4045](https://github.com/wailsapp/wails/pull/4045)）
- 非 React テンプレートで、システムのカラースキームがライトの場合に Hello World フッターが表示されない問題を修正（[@marcus-crane](https://github.com/marcus-crane)、[#4056](https://github.com/wailsapp/wails/pull/4056)）
- macOS で非表示のメニュー項目に関する問題を修正（[@leaanthony](https://github.com/leaanthony)）
- メッセージプロセッサーでのエラー処理と書式設定を修正（[@fbbdev](https://github.com/fbbdev)、[#4066](https://github.com/wailsapp/wails/pull/4066)）
-  アプリケーション終了時にサービスのシャットダウンが省略される問題を修正（[@fbbdev](https://github.com/fbbdev)、[#4066](https://github.com/wailsapp/wails/pull/4066)）
-  メニューの更新がメインスレッドで行われるように変更（[@leaanthony](https://github.com/leaanthony)）
- ドラッグとサイズ変更の仕組みをより堅牢にし、プラットフォームで期待される動作にさらに近づけた（[@fbbdev](https://github.com/fbbdev)、[#4100](https://github.com/wailsapp/wails/pull/4100)）
- [#4097](https://github.com/wailsapp/wails/issues/4097) Webpack/angular がランタイム初期化コードを破棄する問題を修正（[@fbbdev](https://github.com/fbbdev)、[#4100](https://github.com/wailsapp/wails/pull/4100)）
- 初期状態で非表示のメニュー項目に関する問題を修正（[@IanVS](https://github.com/IanVS)、[#4116](https://github.com/wailsapp/wails/pull/4116)）
- 拡張子のないリクエストで、`[request]`が存在せず`[request].html`が存在する場合に、assetFileServer が`.html`ファイルを配信しない問題を修正
- アイコン生成パスを修正（[@robin-samuel](https://github.com/robin-samuel)、[#4125](https://github.com/wailsapp/wails/pull/4125)）
- `fullscreen`、`unfullscreen`、`unminimise`、`unmaximise`イベントが発行されない問題を修正（[@oSethoum](https://github.com/osethoum)、[#4130](https://github.com/wailsapp/wails/pull/4130)）
- 設定内のデフォルトバージョンに誤ったプレフィックスが付いていたために発生する NSIS エラーを修正（[@robin-samuel](https://github.com/robin-samuel)、[#4126](https://github.com/wailsapp/wails/pull/4126)）
- Windows で Dialogs ランタイム関数がエスケープされたパスを返す問題を修正（[TheGB0077](https://github.com/TheGB0077)、[#4188](https://github.com/wailsapp/wails/pull/4188)）
- HKCU 内の Webview2 検出パスを修正（[@leaanthony](https://github.com/leaanthony)）。
- macOS での入力の問題を修正（[@leaanthony](https://github.com/leaanthony)）。
- Windows アイコン生成タスクのファイル名を修正（[@yulesxoxo](https://github.com/yulesxoxo)、[#4219](https://github.com/wailsapp/wails/pull/4219)）。
- @kron の成果を基に、フレームレスウィンドウの透過に関する問題を修正（[@leaanthony](https://github.com/leaanthony)）。
- @kron の成果を基に、ウィンドウが無効または最小化されている場合のフォーカス呼び出しを修正（[@leaanthony](https://github.com/leaanthony)）。
- @kron の成果を基に、タスクバーの再起動後にシステムトレイが表示されない問題を修正（[@leaanthony](https://github.com/leaanthony)）。
- fallbackResponseWriter が Flush() を実装していない問題を修正（[#4245](https://github.com/wailsapp/wails/pull/4245)）
- fallbackResponseWriter が Flush() を実装していない問題を修正（[@superDingda]、[#4236](https://github.com/wailsapp/wails/issues/4236)）
- 非同期の Go バインド関数呼び出しが保留中のときに macOS でウィンドウを閉じるとクラッシュする問題を修正（[@joshhardy](https://github.com/joshhardy)、[#4354](https://github.com/wailsapp/wails/pull/4354)）
- Windows の効率モード起動時の競合状態を修正（[@leaanthony](https://github.com/leaanthony)）
- Windows のアイコンハンドルのクリーンアップを修正（[@leaanthony](https://github.com/leaanthony)）。
- Windows での`OpenFileManager`を修正（[@PPTGamer](https://github.com/PPTGamer)、[#4375](https://github.com/wailsapp/wails/pull/4375)）。
- Linux の最小／最大幅オプションを修正（@atterpac、[#3979](https://github.com/wailsapp/wails/pull/3979)）
- npm のバージョン更新により、TypeScript テンプレートの型定義を修正（@atterpac、[#3966](https://github.com/wailsapp/wails/pull/3966)）
- Sveltekit テンプレートの CSS 参照を修正（@atterpac、[#3945](https://github.com/wailsapp/wails/pull/3945)）
- window run() 内の主要なコールバックがメインスレッドで呼び出されるように変更（[@leaanthony](https://github.com/leaanthony)）
- ダイアログのディレクトリ選択の例を修正（[@leaanthony](https://github.com/leaanthony)）
- index.html が存在しない場合に表示する新しい中国語エラーページを作成（[@leaanthony](https://github.com/leaanthony)）
-  `windowDidBecomeKey`コールバックがメインスレッドで実行されるように変更（[@leaanthony](https://github.com/leaanthony)）
-  フレームレスウィンドウのフルスクリーンをサポート（[@leaanthony](https://github.com/leaanthony)）
-  ウィンドウ破棄ロジックを改善（[@leaanthony](https://github.com/leaanthony)）
-  システムトレイに関連付けられている場合のウィンドウ位置ロジックを修正（[@leaanthony](https://github.com/leaanthony)）
-  フレームレスウィンドウのフルスクリーンをサポート（[@leaanthony](https://github.com/leaanthony)）
- イベント処理を修正（[@leaanthony](https://github.com/leaanthony)）
- ウィンドウのシャットダウンロジックを修正（[@leaanthony](https://github.com/leaanthony)）
- 共通の taskfile で、TypeScript テンプレート用にデフォルトで TypeScript バインディングを生成するように変更（[@leaanthony](https://github.com/leaanthony)）
- ウィンドウが開いていない場合、またはシステムトレイのみの場合に、WM_CLOSE メッセージでアプリケーションを終了する処理を修正（[@mmalcek](https://github.com/mmalcek)、[#3990](https://github.com/wailsapp/wails/pull/3990)）
- garble ビルドを修正（@5aaee9、[#3192](https://github.com/wailsapp/wails/pull/3192)）
- Windows の NSIS ビルドを修正（[@leaanthony](https://github.com/leaanthony)）
- 閉じられていないことが原因で発生する、複数選択用 Linux ダイアログのデッドロックを修正
- Windows ビルド時の .syso ファイルに対するクロスプラットフォームのクリーンアップを修正（
- amd64 AppImage のコンパイルを修正（@atterpac、
- ビルドアセットの更新を修正（@ansxuman、
- @atterpac により Linux のシステムトレイの `OnClick` および `OnRightClick` の実装を修正
- Mac で `AlwaysOnTop` が動作しない問題を修正：
-  `application.NewEditMenu` に重複が含まれる問題を修正
- 🐧 aarch64 でのコンパイルを修正
- ⊞ ラジオグループのメニュー項目を修正：
- MacOS で実行可能な .app をビルドする際のエラーを修正。その発生条件は 'name' と 'outputfilename' が
- ドラッグ＆ドロップのサンプルで customEventProcessor を使用した際のバグを修正：
- 🐧 IgnoreMouseEvents の追加によって発生した Linux のコンパイルエラーを修正：
- ⊞ syso アイコンファイル生成のバグを修正：
- 🐧 Wayland でネイティブに動作させるための修正を次から取り込み：
- 次で内部サービスメソッドをバインドしないように変更：
- ⊞ 次でシステムトレイの起動時パニックを修正：
- 次で内部サービスメソッドをバインドしないように変更：
- ⊞ 次でシステムトレイの起動時パニックを修正：
- メニュー項目を大幅にリファクタリングし、イベント処理を改善。現時点では主に macOS が改善されます。担当：
- 次でプラグインとイベントのリファクタリング後にテストを修正：
- ⊞ `Failed to unregister class Chrome_WidgetWin_0` の警告を修正。担当：
- モジュールの問題
- [atterpac](https://github.com/atterpac) により、次でサイズ変更イベントのメッセージングを修正：
- 🐧 NixOS でのテーマ処理エラーを修正：
- Windows でボリュームをまたぐプロジェクトのインストールを修正：
- フッターが表示されるように React テンプレートの CSS を修正：
- refresh を最新版に更新し、開発モードで作業中にゾンビプロセスが発生する問題を修正
- [Atterpac](https://github.com/atterpac) により AppImage の WebKit ファイル取得を修正
- [Atterpac](https://github.com/Atterpac) により、次で Doctor の apt パッケージ検証を修正：
- @5aaee9 により、次で終了時にアプリケーションがフリーズする問題（Darwin）を修正：
- Windows 上のサンプルの背景色を修正：
- [mmghv](https://github.com/mmghv) により、次でデフォルトのコンテキストメニューを修正：
- Darwin での矢印キーの16進値を修正：
- Windows のドラッグ＆ドロップを動作するように修正。追加者：
- 適切なドライバーがない場合に Doctor で発生する Linux のバグを修正
- 起動時の DPI スケーリング（Windows）を修正。[@almas-x](https://github.com/almas-x) により次で変更：
- 相対パスを使用するように `go.mod` の置換行を修正。次を含む Windows パスの問題を修正：
- ウィンドウが関連付けられていない場合の MacOS のシステムトレイクリック処理を修正：
- 不明なオプションによって Windows のビルドが失敗する問題を修正：
- 次がない状態で Windows のシステムトレイアイコンを左クリックするとクラッシュする問題を修正：
- @5aaee9 により、ウィンドウを2回開いた際に baseURL が誤る問題を PR で修正
- `WebviewWindow.Restore` メソッド内の if 分岐の順序を修正：
- 次の場合に、複数回の `GetStartURL` 呼び出しにわたって `startURL` を正しく計算：
- `Screen` 構造体の JS 型を、対応する Go の型と一致するように修正：
- 登録済みイベントが適切にクリーンアップされるように `WML.Reload` メソッドを修正
- Linux でカスタムコンテキストメニューが即座に閉じる問題を修正：
- バインディングによって生成されるモデルファイルの出力パスと拡張子を修正
- バインディングによって生成される JS コード内のモデルファイルのインポートパスを修正
- 一部の Linux ディストリビューションでドラッグ＆ドロップが動作しない問題を修正：
- `wails3 task dev` を使用する際に macOS 用タスクが欠落する問題を修正：
- イベントの登録によって nil マップへの代入が発生する問題を修正：
- バインドされたメソッドのパラメーターのアンマーシャリングを修正：
- バインドされたメソッドからの複数の戻り値の処理を修正：
- システムパッケージマネージャーでインストールされていない npm を Doctor が検出できない問題を修正
- MicrosoftEdgeWebview2Setup.exe が欠落する問題を修正。協力：
- @leaanthony により、ウィンドウ ID の処理が原因で Linux 上でランダムにクラッシュする問題を修正。次に基づく：
- Linux で systemTray.setIcon がクラッシュする問題を修正：
- `setFrameless` 関数の初回呼び出し時にウィンドウフレームが確実に適用されるように修正：

### 変更

- **破壊的変更**：生成された JS/TS バインディングのマップキーが、Go のマップのセマンティクスを正確に反映するため、オプションとしてマークされるようになりました。TypeScript でマップの値にアクセスすると、`T` ではなく `T | undefined` が返されるようになり、null チェックまたはアサーションが必要です（#4943）。担当：`@fbbdev`
- `@wailsio/runtime` の変更に従い、`Event` の使用を `Events` に変更し、`Features/Events/Event System` のドキュメント内の関数呼び出しも適切に変更：@AbdelhadiSeddar
- `EnabledFeatures`、`DisabledFeatures`、`AdditionalBrowserArgs` をウィンドウ単位のオプションからアプリケーションレベルの `Options.Windows` に移動（#4559）：@leaanthony
- `Drag N Drop` サンプルの README を更新し、このサンプルで `Internal Drag and Drop` を実演していることを強調：@ndianabasi
- 各種デバッグログのレベルを Info から Debug に変更（@mbaklor）
- **破壊的変更：** ウィンドウオプションの `EnableDragAndDrop` を `EnableFileDrop` に改名
- **破壊的変更：** イベントコンテキストの `DropZoneDetails` を `DropTargetDetails` に改名
- **破壊的変更：** `WindowEventContext` の `DropZoneDetails()` メソッドを `DropTargetDetails()` に改名
- **破壊的変更：** `WindowDropZoneFilesDropped` イベントを削除し、代わりに `WindowFilesDropped` を使用
- **破壊的変更:** HTML 属性を `data-wails-dropzone` から `data-file-drop-target` に変更
- **破壊的変更:** CSS のホバークラスを `wails-dropzone-hover` から `file-drop-target-active` に変更
- **破壊的変更:** Windows から `DragEffect`、`OnEnterEffect`、`OnOverEffect` オプションを削除（削除された IDropTarget の一部）
- ランタイムのすべての JSON 処理（メソッドバインディング、イベント、WebView リクエスト、通知、kvstore）を goccy/go-json に切り替え、パフォーマンスを 21-63% 向上し、メモリ割り当てを 40-60% 削減
- BoundMethod 構造体のレイアウトを最適化し、isVariadic フラグをキャッシュして呼び出しごとのオーバーヘッドを削減
- 引数が `<=8` 個のメソッドではスタックに割り当てた引数バッファを使用し、ヒープ割り当てを回避
- メソッド呼び出し時の結果収集を最適化し、戻り値が 1 つの場合のスライス割り当てを回避
- MIME タイプキャッシュに sync.Map を使用し、並行処理性能を向上
- HTTP トランスポートのリクエストボディ読み取りにバッファプールを使用
- コンテンツタイプスニファーで CloseNotify チャネルを遅延割り当てし、リクエストごとの割り当てを削減
- アセットサーバーから CSS のデバッグログを削除
- MIME タイプの拡張子マップを拡充し、一般的な Web 形式（フォント、音声、動画など）を 50 種類以上サポート
- Window の `X/Y` オプションに関するドキュメントを更新（@ruhuang2001）
- フロントエンドバインディングを生成するためのオプションを追加し、`Frontend Runtime` のドキュメントを更新（@ndianabasi）
- Wails v3 Asset Server のドキュメントページを更新（@ndianabasi）
- **破壊的変更**: パッケージレベルのダイアログ関数（`application.InfoDialog()`、`application.QuestionDialog()` など）を削除。代わりに `app.Dialog` マネージャーの `app.Dialog.Info()`、`app.Dialog.Question()`、`app.Dialog.Warning()`、`app.Dialog.Error()`、`app.Dialog.OpenFile()`、`app.Dialog.SaveFile()` を使用
- 実際の API に合わせてダイアログのドキュメントを更新：`app.Dialog.*`、コールバック付きの `AddButton()`（`SetButtons()` ではない）、`SetDefaultButton(*Button)`（文字列ではない）、`AddFilter()`（`SetFilters()` ではない）、`SetFilename()`（`SetDefaultFilename()` ではない）を使用し、フォルダーの選択には `app.Dialog.OpenFile().CanChooseDirectories(true)` を使用
- **破壊的変更**：プロダクションビルドがデフォルトになりました。開発ビルドを作成するには、Taskfile で `DEV=true` を設定してください。設定例を確認するには、新しいプロジェクトを生成してください（@leaanthony）
- データ引数が 0 個または 1 個のカスタムイベントを発行する場合、データ値はスライスでラップされず、Data フィールドに直接割り当てられるようになりました（[@fbbdev](https://github.com/fbbdev)、[#4633](https://github.com/wailsapp/wails/pull/4633)）
- Windows のトレイで `NIS_HIDDEN` を切り替えることにより、`SystemTray.Show()`/`Hide()` が反映されるようになり、アプリを完全に非表示にしてから再表示できるようになりました（#4653）。
- トレイの登録で解決済みのアイコンを再利用し、`NOTIFYICON_VERSION_4` を一度だけ設定するとともに `NIF_SHOWTIP` を有効にすることで、Explorer の再起動後にツールチップが復元されるようになりました（#4653）。
- macOS：メニューバーと Dock の領域を除外してウィンドウを中央に配置するため、`frame` の代わりに `visibleFrame` を使用
- macOS：メニューバーと Dock の領域を除外してウィンドウを中央に配置するため、`frame` の代わりに `visibleFrame` を使用
- `-config` パラメーターを指定して `wails3 update build-assets` を実行した場合、`-product*` パラメーターで設定した値は
- `window.NativeWindowHandle()` → `window.NativeWindow()`（@leaanthony、[#4471](https://github.com/wailsapp/wails/pull/4471)）
- 内部のウィンドウ処理をリファクタリング（@leaanthony、[#4471](https://github.com/wailsapp/wails/pull/4471)）
- `application.WindowIDKey` と `application.WindowNameKey` を削除（`application.WindowKey` に置換）（[@leaanthony](https://github.com/leaanthony)）
- ContextMenuData が any ではなく文字列を返すように変更（[@leaanthony](https://github.com/leaanthony)）
- JS/TS バインディングで、固定長配列型のクラスフィールドが空ではなく、想定される長さで初期化されるようになりました（[@fbbdev](https://github.com/fbbdev)、[#4001](https://github.com/wailsapp/wails/pull/4001)）
- ContextMenuData が any ではなく文字列を返すように変更（[@leaanthony](https://github.com/leaanthony)）
- `application.NewService` はオプションを省略可能なパラメーターとして受け付けなくなりました（代わりに `application.NewServiceWithOptions` を使用）（[@leaanthony](https://github.com/leaanthony)、[#4024](https://github.com/wailsapp/wails/pull/4024)）
- `nanoid` への依存を削除（[@leaanthony](https://github.com/leaanthony)）
- mica/acrylic/tabbed ウィンドウスタイルに対応するよう Window のサンプルを更新（[@leaanthony](https://github.com/leaanthony)）
- JS/TS バインディングから `internal.js/ts` モデルファイルを削除し、すべてのモデルを `models.js/ts` で参照できるようになりました（[@fbbdev](https://github.com/fbbdev)、[#4045](https://github.com/wailsapp/wails/pull/4045)）
- JS/TS バインディングで、名前付き型が他の名前付き型のエイリアスとしてレンダリングされることはなくなり、従来の動作はエイリアスのみに限定されるようになりました（[@fbbdev](https://github.com/fbbdev)、[#4045](https://github.com/wailsapp/wails/pull/4045)）
- JS/TS バインディングのクラスモードで、型パラメーターを型とする構造体フィールドは省略可能としてマークされ、自動的に初期化されなくなりました（[@fbbdev](https://github.com/fbbdev)、[#4045](https://github.com/wailsapp/wails/pull/4045)）
- テンプレートから ESLint を削除（[@IanVS](https://github.com/IanVS)、[#4059](https://github.com/wailsapp/wails/pull/4059)）
- 著作権の日付を 2025 に更新（[@IanVS](https://github.com/IanVS)、[#4037](https://github.com/wailsapp/wails/pull/4037)）
- event.Sender のドキュメントを追加（[@IanVS](https://github.com/IanVS)、[#4075](https://github.com/wailsapp/wails/pull/4075)）
- Go 1.24 をサポート（[@leaanthony](https://github.com/leaanthony)）
- `ServiceStartup` フックは `application.New` 内ではなく、`App.Run` の呼び出し時に実行されるようになりました（[@fbbdev](https://github.com/fbbdev)、[#4066](https://github.com/wailsapp/wails/pull/4066)）
- `ServiceStartup` のエラーはプロセスを終了させず、`App.Run` から返されるようになりました（[@fbbdev](https://github.com/fbbdev)、[#4066](https://github.com/wailsapp/wails/pull/4066)）
- JS からのバインディングおよびダイアログ呼び出しは、文字列ではなくエラーオブジェクトで reject されるようになりました（[@fbbdev](https://github.com/fbbdev)、[#4066](https://github.com/wailsapp/wails/pull/4066)）
- Windows でのシステムトレイメニューの配置を改善（[@leaanthony](https://github.com/leaanthony)）
- JS ランタイムを TypeScript に移植（[@fbbdev](https://github.com/fbbdev)、[#4100](https://github.com/wailsapp/wails/pull/4100)）
- runtime はインポートされるとすぐに初期化されるため、ウィンドウの読み込みを待つ必要がなくなりました。[#4100](https://github.com/wailsapp/wails/pull/4100) の [@fbbdev](https://github.com/fbbdev) による変更です
- runtime は init メソッドをエクスポートしなくなりました。副作用のためのインポートを使用して初期化できます。[#4100](https://github.com/wailsapp/wails/pull/4100) の [@fbbdev](https://github.com/fbbdev) による変更です
- バインドされたメソッドは、キャンセルされた場合に `CancelError` で reject される `CancellablePromise` を返すようになりました。実際の呼び出し結果は破棄されます。[#4100](https://github.com/wailsapp/wails/pull/4100) の [@fbbdev](https://github.com/fbbdev) による変更です
- 組み込みサービス型の名称を `Service` に統一しました。[#4067](https://github.com/wailsapp/wails/pull/4067) の [@fbbdev](https://github.com/fbbdev) による変更です
- オプションを受け取る組み込みサービス作成関数の名称を `NewWithConfig` に統一しました。[#4067](https://github.com/wailsapp/wails/pull/4067) の [@fbbdev](https://github.com/fbbdev) による変更です
- Go API との一貫性を保つため、`sqlite` サービスの `Select` メソッドを `Query` に改名しました。[#4067](https://github.com/wailsapp/wails/pull/4067) の [@fbbdev](https://github.com/fbbdev) による変更です
- テンプレート：runtime を「dependencies」に移し、package.json ファイルを整理しました。[#4133](https://github.com/wailsapp/wails/pull/4133) の [@IanVS](https://github.com/IanVS) による変更です
- 特定の macOS API を利用できるように、開発時にアプリバンドルを作成してアドホック署名するようにしました。[#4171](https://github.com/wailsapp/wails/pull/4171) の [@popaprozac](https://github.com/popaprozac) による変更です
- ビルドアセットをプラットフォーム固有のディレクトリに移動しました。[@leaanthony](https://github.com/leaanthony) による変更です
- Taskfile をプラットフォーム固有のディレクトリに移動し、名前を変更しました。[@leaanthony](https://github.com/leaanthony) による変更です
- `index.html` が見つからない場合のエクスペリエンスを大幅に改善しました。[@leaanthony](https://github.com/leaanthony) による変更です
- [Windows] 最小化と復元のパフォーマンスを改善しました。[@leaanthony](https://github.com/leaanthony) による変更で、[562589540](https://github.com/562589540) による元の [PR](https://github.com/wailsapp/wails/pull/3955) を基にしています
- `ShouldClose` オプションを削除しました（代わりに events.Common.WindowClosing のフックを登録してください）。[@leaanthony](https://github.com/leaanthony) による変更です
- [Windows] ウィンドウを開く際のちらつきを軽減しました。[@leaanthony](https://github.com/leaanthony) による変更です
- 内部関数として意図されていたため、`Window.Destroy` を削除しました。[@leaanthony](https://github.com/leaanthony) による変更です
- `WindowClose` イベントを `WindowClosing` に改名しました。[@leaanthony](https://github.com/leaanthony) による変更です
- フロントエンドのビルドで、ビルド種別に応じて vite 環境「development」または「production」を使用するようになりました。[@leaanthony](https://github.com/leaanthony) による変更です
- go-webview2 v1.19 に更新しました。[@leaanthony](https://github.com/leaanthony) による変更です
- taskfile のフォークが確実に使用されるようにしました。@leaanthony による変更です
- 次の方法でインストールする際のバージョン問題を修正するため、Taskfile のフォークを更新
- 次の方法でインストールする際のバージョン問題を修正するため、Taskfile のフォークを使用
- `service.OnStartup` は、エラー発生時にアプリケーションを終了し、次を実行するようになりました
- ユーザー操作により適切に対応するよう、システムトレイのクリックメッセージ処理をリファクタリング。変更者：
- 次を生成するフレームワークに対応するため、埋め込みアセットに `all:frontend/dist` を追加
- Taskfile のリファクタリング。[leaanthony](https://github.com/leaanthony) により、次で実施：
- `go-webview2` v1.0.16 にアップグレード。変更者：
- `Screen` 型を修正し、`Id` ではなく `ID` を含めるようにしました。変更者：
- `application.ServiceOptions` に対応するため、`go.mod.tmpl` の Wails バージョンを更新。変更者：
- サービス名の判定を修正しました。[windom](https://github.com/windom/) により、次で実施：
- mkdocs serve が docker を使用するようになりました。[leaanthony](https://github.com/leaanthony) による変更です
- 開発用設定を `config.yml` に統合。変更者：
- システムトレイのダイアログで、利用可能な場合はアプリケーションアイコンをデフォルトで使用するようにしました（Windows）。変更者：
- macOS での GPU とメモリのレポートを改善。変更者：
- `WebviewGpuIsDisabled` と `EnableFraudulentWebsiteWarnings` を削除
- Events API の変更：`On`/`Emit` -> ユーザーイベント、`OnApplicationEvent` ->
- Linux の Events API を修正しました。[TheGB0077](https://github.com/TheGB0077) により、次で実施：
- [CI] Actions を改善し、フォークでも Actions を実行できるようにしたほか、
- `AbsolutePosition()` を `Position()` に改名。変更者：
- Linux の WebKit 依存関係を webkitgtk2-4.0 から webkit2gtk-4.1 に更新し、
- 同梱の JS runtime スクリプトが ESM モジュールになりました。このスクリプトをインポートする script タグは
- `@wailsio/runtime` パッケージは、その API を `window.wails` に公開しません
- Window API モジュール `@wailsio/runtime/src/window` が新たに公開するのは、格納元の
- JS の Window API を、現在の Go `WebviewWindow` と一致するように更新しました
- バインディングジェネレーターは、デフォルトで ID による呼び出しを使用するようになりました。`-id` CLI オプションは
- 新しいバインディングコードのレイアウト：以前、出力ファイルはフォルダー内に整理されていました
- 構造体フィールド `application.Options.Bind` の名前を次のように変更しました：
- バインディングサービスの新しい構文：サービスインスタンスを次のものでラップする必要があります
- 非ターミナル環境または CI 環境ではスピナーを無効化。変更者：

### 削除

- **BREAKING**：ウィンドウごとの `WindowsWindow` オプションから `EnabledFeatures`、`DisabledFeatures`、`AdditionalLaunchArgs` を削除しました。代わりに、アプリケーションレベルの `Options.Windows.EnabledFeatures`、`Options.Windows.DisabledFeatures`、`Options.Windows.AdditionalBrowserArgs` を使用してください。これらのフラグは、共有 WebView2 環境全体に適用されます（#4559）。@leaanthony による変更です
- Windows のネイティブ `IDropTarget` 実装を削除し、JavaScript ベースの手法に置き換えました（v2 の動作と同じです）
- github.com/wailsapp/mimetype 依存関係を削除し、拡張した拡張子マップと標準ライブラリの http.DetectContentType に置き換え、バイナリサイズを約 1.2MB 削減しました
- Linux ファイルエクスプローラー用の最小限の .desktop ファイルパーサーを実装して gopkg.in/ini.v1 依存関係を削除し、約 45KB 削減しました
- Go 1.21以降の標準ライブラリの slices パッケージと最小限の内部ヘルパーを使用して、ランタイムコードから samber/lo を削除し、約310KBを削減
- Darwin の URL スキームハンドラーからデバッグ用の printf 文を削除（#4834）
- **破壊的変更**：`linux:WindowLoadChanged` イベントを削除。WebView の読み込み完了を検出するには、代わりに `linux:WindowLoadFinished` を使用（#3896、@leaanthony）

### 破壊的変更

- **Manager API のリファクタリング**：コード構成と見つけやすさを改善するため、application API をフラットな構造から体系化されたマネージャー構造へ再編。[@leaanthony](https://github.com/leaanthony)による変更（[#4359](https://github.com/wailsapp/wails/pull/4359)）
- `app.NewWebviewWindow()` → `app.Window.New()`
- `app.CurrentWindow()` → `app.Window.Current()`
- `app.GetAllWindows()` → `app.Window.GetAll()`
- `app.WindowByName()` → `app.Window.GetByName()`
- `app.EmitEvent()` → `app.Event.Emit()`
- `app.OnApplicationEvent()` → `app.Event.OnApplicationEvent()`
- `app.OnWindowEvent()` → `app.Event.OnWindowEvent()`
- `app.SetApplicationMenu()` → `app.Menu.SetApplicationMenu()`
- `app.OpenFileDialog()` → `app.Dialog.OpenFile()`
- `app.SaveFileDialog()` → `app.Dialog.SaveFile()`
- `app.MessageDialog()` → `app.Dialog.Message()`
- `app.InfoDialog()` → `app.Dialog.Info()`
- `app.WarningDialog()` → `app.Dialog.Warning()`
- `app.ErrorDialog()` → `app.Dialog.Error()`
- `app.QuestionDialog()` → `app.Dialog.Question()`
- `app.NewSystemTray()` → `app.SystemTray.New()`
- `app.GetSystemTray()` → `app.SystemTray.Get()`
- `app.ShowContextMenu()` → `app.ContextMenu.Show()`
- `app.RegisterKeybinding()` → `app.KeyBinding.Register()`
- `app.UnregisterKeybinding()` → `app.KeyBinding.Unregister()`
- `app.GetPrimaryScreen()` → `app.Screen.GetPrimary()`
- `app.GetAllScreens()` → `app.Screen.GetAll()`
- `app.BrowserOpenURL()` → `app.Browser.OpenURL()`
- `app.Environment()` → `app.Env.GetAll()`
- `app.ClipboardGetText()` → `app.Clipboard.Text()`
- `app.ClipboardSetText()` → `app.Clipboard.SetText()`
- Service メソッドの名前を変更：`Name` -> `ServiceName`、`OnStartup` -> `ServiceStartup`、`OnShutdown` -> `ServiceShutdown`。[@leaanthony](https://github.com/leaanthony)による変更
- `Path` メソッドと `Paths` メソッドを `application` パッケージへ移動。[@leaanthony](https://github.com/leaanthony)による変更
- アプリケーションメニューを macOS 専用に変更。[@leaanthony](https://github.com/leaanthony)による変更

## v3.0.0-alpha.78 - 2026-04-21

## 追加

## 修正

## v3.0.0-alpha.77 - 2026-04-18

## 修正

## v3.0.0-alpha.76 - 2026-04-17

## 修正

## v3.0.0-alpha.75 - 2026-04-16

## 修正

## v3.0.0-alpha.74 - 2026-03-01

## 追加

## 修正

## v3.0.0-alpha.73 - 2026-02-27

## 修正

## v3.0.0-alpha.72 - 2026-02-16

## 修正

## v3.0.0-alpha.71 - 2026-02-10

## 追加

## 修正

## v3.0.0-alpha.70 - 2026-02-09

## 追加

## 修正

## v3.0.0-alpha.69 - 2026-02-08

## 追加

## 修正

## v3.0.0-alpha.68 - 2026-02-07

## 追加

## 変更

## 修正

## v3.0.0-alpha.67 - 2026-02-04

## 追加

## 変更

## 修正

## v3.0.0-alpha.66 - 2026-02-03

## 追加

## 変更

## 修正

## 削除

## v3.0.0-alpha.65 - 2026-02-01

## 追加

## v3.0.0-alpha.64 - 2026-01-26

## 追加

## v3.0.0-alpha.63 - 2026-01-25

## 修正

## v3.0.0-alpha.62 - 2026-01-22

## 修正

## v3.0.0-alpha.61 - 2026-01-20

## 修正

## v3.0.0-alpha.60 - 2026-01-14

## 修正

## v3.0.0-alpha.59 - 2026-01-11

## 変更

## v3.0.0-alpha.58 - 2026-01-09

## 修正

## v3.0.0-alpha.57 - 2026-01-05

## 変更

## 修正

## v3.0.0-alpha.56 - 2026-01-04

## 追加

## 変更

## 修正

## 削除

## v3.0.0-alpha.55 - 2026-01-02

## 変更

## 修正

## 削除

## v3.0.0-alpha.54 - 2025-12-29

## 追加

## 修正

## 削除

## v3.0.0-alpha.53 - 2025-12-27

## 追加

## 修正

## v3.0.0-alpha.52 - 2025-12-26

## 修正

## v3.0.0-alpha.51 - 2025-12-23

## 修正

## v3.0.0-alpha.50 - 2025-12-21

## 変更

## v3.0.0-alpha.49 - 2025-12-18

## 変更

## v3.0.0-alpha.48 - 2025-12-16

## 追加

## 変更

## 修正

## v3.0.0-alpha.47 - 2025-12-15

## 追加

## 修正

## v3.0.0-alpha.46 - 2025-12-14

## 追加

## 削除

## v3.0.0-alpha.45 - 2025-12-13

## 追加

## 修正

## v3.0.0-alpha.44 - 2025-12-12

## 追加

## 変更

## 修正

## v3.0.0-alpha.43 - 2025-12-11

## 追加

## v3.0.0-alpha.42 - 2025-12-10

## 追加

## v3.0.0-alpha.41 - 2025-11-23

## 修正

## v3.0.0-alpha.40 - 2025-11-13

## 修正

## v3.0.0-alpha.39 - 2025-11-12

## 追加

## 変更

## v3.0.0-alpha.38 - 2025-11-04

## 追加

## 変更

## 修正

## v3.0.0-alpha.37 - 2025-11-02

## 修正

## v3.0.0-alpha.36 - 2025-10-15

## 修正

## v3.0.0-alpha.35 - 2025-10-14

## 修正

## v3.0.0-alpha.34 - 2025-10-06

## 追加

## 修正

## v3.0.0-alpha.33 - 2025-10-04

## 修正

## v3.0.0-alpha.32 - 2025-10-02

## 修正

## v3.0.0-alpha.31 - 2025-09-27

## 修正

## v3.0.0-alpha.30 - 2025-09-26

## 修正

## v3.0.0-alpha.29 - 2025-09-25

## 追加

## 変更

## 修正

## v3.0.0-alpha.29 - 2025-09-25

## 追加

## 変更

## 修正

## v3.0.0-alpha.27 - 2025-09-07

## 修正

## v3.0.0-alpha.26 - 2025-08-24

## 追加

## v3.0.0-alpha.25 - 2025-08-16

## 変更

無視されなくなり、設定値を上書きするようになりました。

## v3.0.0-alpha.24 - 2025-08-13

## 追加

## v3.0.0-alpha.23 - 2025-08-11

## 修正

## v3.0.0-alpha.22 - 2025-08-10

## 追加

## 変更

+ Linuxパッケージの過剰に広範な依存関係と、古くなったRPM依存関係を修正。

## v3.0.0-alpha.21 - 2025-08-07

## 修正

## v3.0.0-alpha.20 - 2025-08-06

## 修正

## v3.0.0-alpha.19 - 2025-08-05

## 追加

## 修正

## v3.0.0-alpha.18 - 2025-08-03

## 追加

## 修正

## v3.0.0-alpha.17 - 2025-07-31

## 修正

## v3.0.0-alpha.16 - 2025-07-25

## 追加

## v3.0.0-alpha.15 - 2025-07-25

## 追加

## v3.0.0-alpha.14 - 2025-07-25

## 追加

## v3.0.0-alpha.12 - 2025-07-15

### 追加

### 修正

## v3.0.0-alpha.11 - 2025-07-12

## 追加

## v3.0.0-alpha.10 - 2025-07-06

### 破壊的変更

### 追加

### 修正

### 変更

## v3.0.0-alpha.9 - 2025-01-13

### 追加

### 修正

### 変更

## v3.0.0-alpha.8.3 - 2024-12-07

### 変更

## v3.0.0-alpha.8.2 - 2024-12-07

### 変更

`go install`（@leaanthonyによる変更）

## v3.0.0-alpha.8.1 - 2024-12-07

### 変更

`go install`（@leaanthonyによる変更）

## v3.0.0-alpha.8 - 2024-12-06

### 追加

@atterpac が [#3909](https://github.com/wailsapp/wails/3909) で対応   [ansxuman](https://github.com/ansxuman) が   [#3902](https://github.com/wailsapp/wails/pull/3902) で対応   [atterpac](https://github.com/atterpac) が   [#3867](https://github.com/wailsapp/wails/pull/3867) で対応   [atterpac](https://github.com/atterpac) が   [#3829](https://github.com/wailsapp/wails/pull/3829) で対応   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) が   [#3856](https://github.com/wailsapp/wails/pull/3856) で対応   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   指定された X/Y 座標に配置。   [leaanthony](https://github.com/leaanthony) が   [#3885](https://github.com/wailsapp/wails/pull/3885) で対応   [ansxuman](https://github.com/ansxuman) と   [leaanthony](https://github.com/leaanthony) が   [#3823](https://github.com/wailsapp/wails/pull/3823) で対応   [leaanthony](https://github.com/leaanthony) が   [#3766](https://github.com/wailsapp/wails/pull/3766) で対応   [leaanthony](https://github.com/leaanthony) が   [#3888](https://github.com/wailsapp/wails/pull/3888) で対応   [leaanthony](https://github.com/leaanthony)。 -

### 変更

起動済みのすべてのサービスを対象とする`service.OnShutdown`について、@atterpac が   [#3920](https://github.com/wailsapp/wails/pull/3920) で対応   @atterpac が [#3907](https://github.com/wailsapp/wails/pull/3907) で対応   サブフォルダーについて @atterpac が   [#3887](https://github.com/wailsapp/wails/pull/3887) で対応   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) が   [#3778](https://github.com/wailsapp/wails/pull/3778) で対応   [northes](https://github.com/northes) が   [#3836](https://github.com/wailsapp/wails/pull/3836) で対応   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   （`EnabledFeatures` オプションと `DisabledFeatures` オプションに置き換え）   [leaanthony](https://github.com/leaanthony) が対応

### 修正

channel 変数を修正。@michael-freling が   [#3925](https://github.com/wailsapp/wails/pull/3925) で対応   [ansxuman](https://github.com/ansxuman) が   [#3924](https://github.com/wailsapp/wails/pull/3924) で対応   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   [#3886](https://github.com/wailsapp/wails/pull/3886) で対応   [leaanthony](https://github.com/leaanthony) が   [#3841](https://github.com/wailsapp/wails/pull/3841) で対応   Darwin の編集メニューにおける `PasteAndMatchStyle` ロールを修正。   [johnmccabe](https://github.com/johnmccabe) が   [#3839](https://github.com/wailsapp/wails/pull/3839) で対応   [#3840](https://github.com/wailsapp/wails/issues/3840) を   [#3854](https://github.com/wailsapp/wails/pull/3854) で   [kodflow](https://github.com/kodflow) が対応   [@leaanthony](https://github.com/leaanthony)   が異なる問題を修正。@nickisworking が   [#3789](https://github.com/wailsapp/wails/pull/3789) で対応

## v3.0.0-alpha.7 - 2024-09-18

### 追加

[mmghv](https://github.com/mmghv) が   [#3665](https://github.com/wailsapp/wails/pull/3665) で対応   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) と   [leaanthony](https://github.com/leaanthony) が   [#3570](https://github.com/wailsapp/wails/pull/3570) で対応

### 変更

アプリケーションイベント `OnWindowEvent` をウィンドウイベントに変更。   [leaanthony](https://github.com/leaanthony) が   [#3734](https://github.com/wailsapp/wails/pull/3734) で対応   `v3/` または `v3-` というプレフィックスが付いたブランチ。   [stendler](https://github.com/stendler) が   [#3747](https://github.com/wailsapp/wails/pull/3747) で対応

### 修正

[etesam913](https://github.com/etesam913) が   [#3742](https://github.com/wailsapp/wails/pull/3742) で対応   [atterpac](https://github.com/atterpac) が   [#3721](https://github.com/wailsapp/wails/pull/3721) で対応   [atterpac](https://github.com/atterpac) が   [#3675](https://github.com/wailsapp/wails/pull/3675) で対応   [#1811](https://github.com/wailsapp/wails/pull/1811) を   [#3614](https://github.com/wailsapp/wails/pull/3614) で   [@stendler](https://github.com/stendler) が対応   [#3720](https://github.com/wailsapp/wails/pull/3720) を   [leaanthony](https://github.com/leaanthony) が対応   [#3693](https://github.com/wailsapp/wails/issues/3693) を   [@DeltaLaboratory](https://github.com/DeltaLaboratory) が対応   [#3720](https://github.com/wailsapp/wails/pull/3720) を   [leaanthony](https://github.com/leaanthony) が対応   [#3693](https://github.com/wailsapp/wails/issues/3693) を   [@DeltaLaboratory](https://github.com/DeltaLaboratory) が対応   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746) を   [@stendler](https://github.com/stendler) が対応   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### 修正

## v3.0.0-alpha.5 - 2024-07-30

### 追加

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   アイコンのクリック。@5aaee9 が [#2991](https://github.com/wailsapp/wails/pull/2991) で対応   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) が   [#3282](https://github.com/wailsapp/wails/pull/3282) で対応   @[Atterpac](https://github.com/Atterpac) が   [#3022](https://github.com/wailsapp/wails/pull/3022]) で対応   [@marcus-crane](https://github.com/marcus-crane) が   [#3146](https://github.com/wailsapp/wails/pull/3146) で対応   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) が   [#3281](https://github.com/wailsapp/wails/pull/3281) で対応   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   @Mai-Lapyst による [PR](https://github.com/wailsapp/wails/pull/2044) に基づく   [@fbbdev](https://github.com/fbbdev) が   [#3295](https://github.com/wailsapp/wails/pull/3295) で対応   [@fbbdev](https://github.com/fbbdev) が   [#3295](https://github.com/wailsapp/wails/pull/3295) で対応   [@fbbdev](https://github.com/fbbdev) が   [#3295](https://github.com/wailsapp/wails/pull/3295) で対応   npm パッケージ。[@fbbdev](https://github.com/fbbdev) が   [#3334](https://github.com/wailsapp/wails/pull/3334) で対応   [#3354](https://github.com/wailsapp/wails/pull/3354) で対応   `WAILS_VITE_PORT`。[@abichinger](https://github.com/abichinger) が   [#3429](https://github.com/wailsapp/wails/pull/3429) で対応   [@abichinger](https://github.com/abichinger) が   [#3431](https://github.com/wailsapp/wails/pull/3431) で対応   [@bruxaodev](https://github.com/bruxaodev) が   [#3667](https://github.com/wailsapp/wails/pull/3667) で対応   [@OlegGulevskyy](https://github.com/OlegGulevskyy) が   [#3674](https://github.com/wailsapp/wails/pull/3674) で対応

### 修正

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane)により   [#3515](https://github.com/wailsapp/wails/pull/3515)で対応   [atterpac](https://github.com/atterac)により   [#3512](https://github.com/wailsapp/wails/pull/3512)で対応   [atterpac](https://github.com/atterpac)により   [#3477](https://github.com/wailsapp/wails/pull/3477)で対応   [Atterpac](https://github.com/atterpac)により   [#3320](https://github.com/wailsapp/wails/pull/3320)で対応。   [#3306](https://github.com/wailsapp/wails/pull/3306)で対応。   [#2972](https://github.com/wailsapp/wails/pull/2972)。   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv)により   [#2750](https://github.com/wailsapp/wails/pull/2750)で対応。   [#2753](https://github.com/wailsapp/wails/pull/2753)。   [jaybeecave](https://github.com/jaybeecave)により   [#3052](https://github.com/wailsapp/wails/pull/3052)で対応。   [@pylotlight](https://github.com/pylotlight)により   [PR](https://github.com/wailsapp/wails/pull/3039)で対応   インストール済み。[@pylotlight](https://github.com/pylotlight)により   [PR](https://github.com/wailsapp/wails/pull/3032)で追加   [PR](https://github.com/wailsapp/wails/pull/3145)   スペース（@leaanthony）。   [thomas-senechal](https://github.com/thomas-senechal)により PR   [#3207](https://github.com/wailsapp/wails/pull/3207)で対応   [thomas-senechal](https://github.com/thomas-senechal)により PR   [#3208](https://github.com/wailsapp/wails/pull/3208)で対応   アタッチされたウィンドウ。[tw1nk](https://github.com/tw1nk)により PR   [#3271](https://github.com/wailsapp/wails/pull/3271)で対応   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev)により   [#3279](https://github.com/wailsapp/wails/pull/3279)で対応   `FRONTEND_DEVSERVER_URL`が存在する。   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev)により   [#3295](https://github.com/wailsapp/wails/pull/3295)で対応   リスナー。[@fbbdev](https://github.com/fbbdev)により   [#3295](https://github.com/wailsapp/wails/pull/3295)で対応   [@abichinger](https://github.com/abichinger)により   [#3330](https://github.com/wailsapp/wails/pull/3330)で対応   ジェネレーター。[@fbbdev](https://github.com/fbbdev)により   [#3334](https://github.com/wailsapp/wails/pull/3334)で対応   ジェネレーター。[@fbbdev](https://github.com/fbbdev)により   [#3334](https://github.com/wailsapp/wails/pull/3334)で対応   [@abichinger](https://github.com/abichinger)により   [#3346](https://github.com/wailsapp/wails/pull/3346)で対応   [@hfoxy](https://github.com/hfoxy)により   [#3417](https://github.com/wailsapp/wails/pull/3417)で対応   [@hfoxy](https://github.com/hfoxy)により   [#3426](https://github.com/wailsapp/wails/pull/3426)で対応   [@fbbdev](https://github.com/fbbdev)により   [#3431](https://github.com/wailsapp/wails/pull/3431)で対応   [@fbbdev](https://github.com/fbbdev)により   [#3431](https://github.com/wailsapp/wails/pull/3431)で対応   [@pekim](https://github.com/pekim)により   [#3458](https://github.com/wailsapp/wails/pull/3458)で対応   [@robin-samuel](https://github.com/robin-samuel)。   [@5aaee9](https://github.com/5aaee9)による PR [#3466](https://github.com/wailsapp/wails/pull/3622)。   [@windom](https://github.com/windom/)により   [#3636](https://github.com/wailsapp/wails/pull/3636)で対応。   Windows。[@bruxaodev](https://github.com/bruxaodev/)により   [#3691](https://github.com/wailsapp/wails/pull/3691)で対応。

### 変更

[mmghv](https://github.com/mmghv)により   [#3611](https://github.com/wailsapp/wails/pull/3611)で対応   Ubuntu 24.04 LTS をサポート。[atterpac](https://github.com/atterpac)により   [#3461](https://github.com/wailsapp/wails/pull/3461)で対応   には `type="module"` 属性が必要。   [@fbbdev](https://github.com/fbbdev)により   [#3295](https://github.com/wailsapp/wails/pull/3295)で対応   オブジェクトであり、WML システムを起動しない。この変更は、カプセル化を改善するために行われた。必要に応じて、新しい `WML.Enable` メソッドを呼び出すことで WML システムを手動で起動できる。バンドルされた JS ランタイムスクリプトは、引き続き両方の操作を自動的に実行する。[@fbbdev](https://github.com/fbbdev)により   [#3295](https://github.com/wailsapp/wails/pull/3295)で対応   ウィンドウオブジェクトをデフォルトエクスポートとして提供する。ESM の名前付きインポート構文または名前空間インポート構文を使用して個別のメソッドをインポートすることは、今後できない。   API。一部のメソッドでは名前またはプロトタイプが変更された。具体的には、`Screen` は `GetScreen` に、`GetZoomLevel`/`SetZoomLevel` は `GetZoom`/`SetZoom` になった。`GetZoom`、`Width`、`Height` は、値をオブジェクト内にラップせず、直接返すようになった。[@fbbdev](https://github.com/fbbdev)により   [#3295](https://github.com/wailsapp/wails/pull/3295)で対応   は削除された。名前による呼び出しに戻すには、`-names` CLI オプションを使用する。   [@fbbdev](https://github.com/fbbdev)により   [#3468](https://github.com/wailsapp/wails/pull/3468)で対応   には従来、それらを含むパッケージに基づく名前が付けられていたが、現在はモジュールパスを含む完全な Go インポートパスが使用される。[@fbbdev](https://github.com/fbbdev)により   [#3468](https://github.com/wailsapp/wails/pull/3468)で対応   `application.Options.Services`。[@fbbdev](https://github.com/fbbdev)により   [#3468](https://github.com/wailsapp/wails/pull/3468)で対応   `application.NewService` の呼び出し。[@fbbdev](https://github.com/fbbdev)により   [#3468](https://github.com/wailsapp/wails/pull/3468)で対応   [@DeltaLaboratory](https://github.com/DeltaLaboratory)により   [#3574](https://github.com/wailsapp/wails/pull/3574)で対応
