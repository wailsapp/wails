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

/_-- このプロジェクトの注目すべき変更はすべて、このファイルに記録されます。

この形式は[Keep a Changelog](https://keepachangelog.com/en/1.0.0/)に基づいており、このプロジェクトは[セマンティック バージョニング](https://semver.org/spec/v2.0.0.html)に準拠しています。

- `Added`：新機能。
- `Changed`：既存機能の変更。
- `Deprecated`：近日中に削除される機能。
- `Removed`：削除済みの機能。
- `Fixed`：あらゆるバグ修正。
- `Security`：脆弱性に関する変更。

_/

/_   * 自動更新 *   リリースノートは、マージ済みの PR から生成されます。通常のコントリビューション PR ではこのファイルを編集せず、削除された`v3/UNRELEASED_CHANGELOG.md`を再作成しないでください。変更履歴のメンテナンスについては、`AGENTS.md`を参照してください。 _/

## [未リリース]

## v3.2.0-nightly.20261103 - 2026-11-03

## 変更
- [v3] Menu.Update 時にシステムトレイのメニューを更新（@johnmaguire による[PR](https://github.com/wailsapp/wails/pull/6314)）
- [v3/Windows] ウィンドウ用の SetIcon を実装（@hypherionmc による[PR](https://github.com/wailsapp/wails/pull/6502)）

## v3.2.0-nightly.20261102 - 2026-11-02

## 変更
- DisableBounce を DisableOverscroll に変更（@lethusang による[PR](https://github.com/wailsapp/wails/pull/6077)）

## v3.2.0-nightly.20261101 - 2026-11-01

## 追加
- CGO を使用しない実験的な macOS および Linux バックエンドを追加（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6521)）
- ネイティブ YAML ビルド、ポータブルなパッケージングおよび署名を追加（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6531)）

## 修正
- v3の各種ナイトリービルドの問題を修正し、依存関係を整理（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6519)）
- v3のリリース検証とフィクスチャ分類を修復（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6533)）
- クロスプラットフォームのリリース検証エラーを解消（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6538)）

## v3.1.0 - 2026-11-01

## 追加
- CGO を使用しない実験的な macOS および Linux バックエンドを追加（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6521)）
- ネイティブ YAML ビルド、ポータブルなパッケージングおよび署名を追加（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6531)）

## 変更
- 破棄されるデスクトップメタデータのコピーを回避（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6391)）
- fmt を使用せずにコンテンツ長をフォーマット（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6393)）
- 候補を収集せずにインデックスルートを選択（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6404)）
- シーク不可能なファイルのコピーバッファーを再利用（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6408)）
- プレフィックス判定用リーダーの割り当てを回避（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6410)）
- 既知の拡張子では MIME 判定のための読み取りを省略（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6414)）
- キューの増大と縮小の繰り返しを防止（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6415)）
- 破棄されるリクエストコンテキストの割り当てを回避（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6419)）
- 無効なデバッグログの準備を省略（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6420)）
- MIME 判定用バッファーを遅延割り当て（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6428)）
- ハイパーリンクのフォーマット処理に伴う割り当てを削減（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6443)）
- 箇条書きのフォーマット処理に伴う一時的な割り当てを削除（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6449)）
- UUID の文字列変換を 1 回の割り当てに削減（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6453)）
- ステートメントハンドルのフォーマット処理に伴う割り当てを削減（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6454)）
- 名前ベース UUID のダイジェストに伴う一時的な割り当てを排除（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6455)）
- 未完了のコマンド出力行のコピーを回避（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6460)）
- パッケージのステージング中に実行可能ファイルをストリーミングコピー（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6472)）
- 未完了の出力行の再スキャンを回避（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6473)）
- コンパイル済みのコマンドルーティングパターンを再利用（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6475)）
- ランタイム検出中のメタデータ検索を回避（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6482)）
- ローカライズされたフォールバックページの割り当てを削減（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6485)）
- プロセッサごとのパターン再コンパイルを回避（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6490)）

## 修正
- システムトレイ用の SetTooltip を実装（@fan711 による[PR](https://github.com/wailsapp/wails/pull/6017)）
- 設定ポータルからダークモードを読み取り（@TechHutTV による[PR](https://github.com/wailsapp/wails/pull/6072)）
- Window.SetZoom で1.0未満を指定可能に変更（@APshenkin による[PR](https://github.com/wailsapp/wails/pull/6166)）
- 各 Go リスナーに emit 順で配信（@DevLumuz による[PR](https://github.com/wailsapp/wails/pull/6275)）
- Emit で hookLock の保持中にイベントフックを読み取り（@DevLumuz による[PR](https://github.com/wailsapp/wails/pull/6276)）
- Run の後に作成されたトレイが破棄されないよう、トレイの Menu プロパティを常に公開（@AlbinoGeek による[PR](https://github.com/wailsapp/wails/pull/6284)）
- Run の前に送信されたメインスレッドの処理を、Run が開始されるまで保留（@AlbinoGeek による[PR](https://github.com/wailsapp/wails/pull/6287)）
- GitHub アップデータープロバイダーで正しいアセット URL を使用（@gavidroselj による[PR](https://github.com/wailsapp/wails/pull/6304)）
- WebView2 が Content-Length を省略した場合でもリクエストボディを保持（@strraaa による[PR](https://github.com/wailsapp/wails/pull/6326)）
- macOS および GTK の画面情報を DIP 空間で報告（DPI が混在する境界）（@julianstorer による[PR](https://github.com/wailsapp/wails/pull/6346)）
- 閉じられた通知をクリックとして報告しないよう修正（@bercedev による[PR](https://github.com/wailsapp/wails/pull/6367)）
- ターミナルを閉じたときに開発プロセスを停止（@xmasdev による[PR](https://github.com/wailsapp/wails/pull/6371)）
- エビデンスチェックで変更履歴の実行を squash merge コミットに対応付け（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6373)）
- タスクバー再起動のデバウンスを1秒間維持する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6386)）
- フルスクリーンで開始した場合もサイズ制約を維持する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6387)）
- 新しい制約に合わせるために、フルスクリーンまたは最大化されたウィンドウのサイズを変更しない（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6388)）
- liquid-glass のサンプルで Run より前にダイアログを表示しない（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6389)）
- Run より前に要求されたダイアログは、永遠に待機せず何もしない（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6390)）
- デスクトップエントリにコマンドがない場合はフォールバックする（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6392)）
- ファイル名を保持せずに現在のコンテンツを判別する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6394)）
- 大文字と小文字が混在するアセット拡張子を認識する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6395)）
- 使用済みのコンテンツ判別用バッファを解放する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6396)）
- ringqueue から削除された参照を解放する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6397)）
- フォールバックの選択前に、ヘッダーの削除を反映する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6398)）
- フラッシュ前にフォールバックレスポンスのヘッダーを確定する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6399)）
- 名前が完全に一致するインデックスファイルのみを検出する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6400)）
- シーク不能なリソースへの HEAD では、ボディと全体読み込みを回避する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6401)）
- 短い読み込みを繰り返して MIME 判別用プレフィックスを満たす（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6402)）
- 成功した空のシーク不能レスポンスを確定する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6403)）
- 展開された URI ターゲットの重複を回避する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6405)）
- フィールド展開時にパスを保持する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6406)）
- XDG のデスクトップファイル検索ディレクトリに従う（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6407)）
- 空のフラッシュでもコンテンツ判別用ヘッダーを確定する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6409)）
- デスクトップエントリを隠しているディレクトリをスキップする（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6411)）
- 実行可能ファイルのパス内で見つかったターゲットを保持する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6412)）
- タイマーの発火後にコールバックのキャプチャを解放する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6413)）
- ldconfig でライブラリが誤一致することを回避する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6416)）
- ファイルシステムでのライブラリのフォールバック照合を限定する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6417)）
- キャッシュ無効化時のデータ競合を解消する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6418)）
- 任意指定のロガー設定を反映する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6421)）
- 共有される開始 URL の変更と古いポートの使用を回避する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6422)）
- リクエストがキャンセルされたら開発プロキシの再試行を停止する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6423)）
- バンドルされたアセットがない場合は 404 を返す（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6424)）
- バンドルされたランタイムへの HEAD リクエストを処理する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6425)）
- サービスのルーティング時にエスケープ済みパスを保持する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6426)）
- サイズが上限を超える数値のプレリリース識別子を正しく並べ替える（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6427)）
- ツールのダウンロードを保存する前に HTTP エラーを拒否する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6429)）
- デスクトップ検出時に起動期限を遵守する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6430)）
- 開発プロキシの再試行時もリクエストボディを保持する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6432)）
- すべてのリターンパスで HTTP レスポンスボディを閉じる（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6434)）
- 並列ライブラリ検索でキャンセルエラーを保持する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6435)）
- デコード前に失敗した GitHub タグレスポンスを拒否する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6436)）
- 角括弧を含むテンプレートパスをリテラルのまま保持する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6437)）
- wails3 を一度だけインストールし、インストーラーのエラーを保持する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6438)）
- コメントブロック内の非表示エントリとヘッダーを無視する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6439)）
- 安定版リリースの検索エラーを通知する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6440)）
- プレーンテキストのメッセージをフォーマットとして解釈せずに出力する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6441)）
- JSON の null を読み込んだ後の書き込み時パニックを回避する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6442)）
- ストア全体の読み込みではスナップショットを返す（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6444)）
- 行の反復処理が一部完了した後のエラーを伝播する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6446)）
- 展開先のパスをリテラルのまま保持する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6447)）
- クローズおよび再オープン時にプリペアドステートメントのハンドルを解放する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6448)）
- ステートメントの準備とクローズまたは再オープンとの競合を処理する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6450)）
- フォールバック用インデックスページへの HEAD リクエストを正しく処理する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6451)）
- 完全な16進数形式のステートメントハンドルを必須にする（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6452)）
- 並行処理下でもデフォルトのインメモリデータベースを維持する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6456)）
- 解析対象の HEAD 出力から stderr を分離する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6457)）
- 埋め込みバンドルからランタイムアセットを生成する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6458)）
- TODO、FIXME、TBD のプレースホルダーを検出する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6459)）
- リリースノート API のレスポンスを検証する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6461)）
- ステータスチェック時にタスクディレクトリを考慮する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6462)）
- 同一ファイルへのコピーによるデータ損失を防止する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6463)）
- キャッシュの識別情報にタスクディレクトリを含める（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6464)）
- 各チェックでポートを一度だけプローブする（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6465)）
- キャッシュエントリが null としてデコードされた場合のパニックを回避する（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6466)）
- 整数の上限に制約されずに数値のプレリリース識別子をインクリメントする（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6467)）
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6468)で、レガシーチェックによってキャッシュミスが隠される問題を防止
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6469)で、デフォルトタスクをエグゼキューターのベースディレクトリで実行
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6470)で、エクスプローラーコマンド内のリテラルのパーセント記号を保持
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6471)で、事前条件とタスクコマンドのディレクトリを一致させるように変更
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6474)で、シンボリックリンクを生成された通常ファイルに置換
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6476)で、ビルド失敗時にも既存の DMG を保持
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6477)で、引用符で囲まれた正規表現とパスのバックスラッシュを保持
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6478)で、画像を割り当てる前に生成アイコンのサイズを検証
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6479)で、コマンドチェーンとリダイレクトをシェル経由で処理
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6480)で、ビルド失敗時にも以前の Linux パッケージを保持
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6481)で、プラットフォーム固有のシェルフォールバック実行を使用
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6483)で、先頭のパス区切りスラッシュが重複している場合もアセットを配信
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6484)で、ビルドメタデータを含む安定版を認識
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6486)で、リソース生成の失敗を正確に報告
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6487)で、アクティブなレポーターの置換時に発生する競合を防止
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6488)で、並行して行われるキャッシュファイルの保存を直列化
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6489)で、テンプレートの複数の等価候補をサポート
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6491)で、保存済みカテゴリーマップが null の場合に正規化
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6492)で、失敗した変数シェルコマンドの診断情報を表示
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6493)で、シェルのコマンドチェーンを暗黙的な Go キャッシュの対象外に変更
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6494)で、Linux におけるカテゴリーの永続化を直列化
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6495)で、引用符またはエスケープされた空白を含む Go 出力をキャッシュ
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6496)で、エンコード失敗時にも Mac アイコンを保持
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6497)で、無効な Go テスト回数フラグを破棄しないように修正
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6498)で、新規ファイルへのコピー時にコピー元の権限を保持
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6519)で、各種v3ナイトリービルドの問題を修正し、依存関係を整理
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6533)で、v3リリース検証とフィクスチャ分類を修正
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6538)で、クロスプラットフォームのリリース検証エラーを解消

## v3.1.0-nightly.20261028 - 2026-10-28

## 変更
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6490)で、コンパイル済み正規表現パターンを再利用して変更履歴の生成を高速化
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6485)で、ローカライズされたフォールバックページの配信時の割り当てを削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6482)で、不要なメタデータ検索を避けて AppImage ランタイムの検出を高速化
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6475)で、コンパイル済みコマンドルーティングパターンを再利用し、wake のパフォーマンスを向上
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6473)で、未完了の出力行の再スキャンを避けてパフォーマンスを向上
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6472)で、MSIX パッケージのステージング中に実行可能ファイルをストリーミングコピーし、メモリ使用量を削減してビルドを高速化
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6460)で、未完了のコマンド出力行のコピーを避けてパフォーマンスを向上
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6455)で、名前ベース UUID の生成時の割り当てを削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6454)で、SQLite ステートメントハンドルのフォーマット時の割り当てを削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6453)で、UUID の文字列変換における割り当てを 1 回に削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6449)で、変更履歴の箇条書き項目をフォーマットする際の割り当てを削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6443)で、ターミナルのハイパーリンクをフォーマットする際の割り当てを削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6428)で、MIME バッファを遅延判定し、アセット配信時の割り当てを削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6420)で、デバッグログが無効な場合、アセットサーバーでのデバッグログ準備を省略
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6419)で、破棄されるリクエストコンテキストを作らず、webview のリクエスト処理における割り当てを削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6414)で、既知のファイル拡張子に対する MIME 判定を省略し、アセット配信を高速化
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6410)で、コンテンツ判定用のリーダーを使用せず、アセット配信時の割り当てを削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6408)で、シーク不能なアセットファイルの配信時の割り当てを削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6404)で、候補の収集を避けてアセットサーバーのインデックスルート選択を高速化
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6393)で、fmt を使わずにコンテンツ長をフォーマットし、アセットサーバーのパフォーマンスを向上
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6391)で、デスクトップファイルエクスプローラーのメタデータ読み取り時のメモリコピーを削減

## 修正
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6498)で、新しいコピー先ファイルへのコピー時にコピー元ファイルの権限を保持
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6388)で、最小または最大サイズの制約が変更されても、全画面表示または最大化されたウィンドウのサイズが変わらないように修正
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6500)で、wailsapp/wails リポジトリのスター履歴をスケジュールに従って更新・公開
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6497)で、wake が無効な Go テスト回数フラグを破棄しないように修正
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6496)で、エンコード失敗時にも Mac アイコンを保持
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6495)で、引用符またはエスケープされた空白を含む Go 出力でビルドキャッシュがミスする問題を修正
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6494)で、Linux における通知カテゴリーの永続化を直列化し、並行書き込みの競合を防止
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6493)で、シェルのコマンドチェーンを暗黙的な Go ビルドキャッシュの対象外に変更
- 保存データが null の場合に通知カテゴリーマップを読み込めない問題を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6491)）
- 変数を含むシェルコマンドが失敗した場合に診断情報を表示（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6492)）
- 等価条件の候補が複数ある場合の wake テンプレートのマッチングを修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6489)）
- キャッシュファイルの同時保存を直列化し、キャッシュ書き込みの破損を防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6488)）
- アクティブなレポーターを置き換える際の競合を防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6487)）
- リソース生成の失敗を CLI で正確に報告（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6486)）
- ビルドメタデータを含む安定版を CLI で認識（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6484)）
- リクエストパスの先頭にスラッシュが連続しているアセットを配信（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6483)）
- 各プラットフォームでの wake のシェルフォールバック実行を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6481)）
- ビルド失敗時に既存の Linux パッケージを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6480)）
- コマンドチェーンとリダイレクトをシェル経由で実行（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6479)）
- wake コマンドの解析時に、引用符で囲まれた正規表現とパスのバックスラッシュを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6477)）
- ビルドのパッケージ化に失敗した場合も既存の DMG ファイルを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6476)）
- バインディング出力内のシンボリックリンクを、生成した通常ファイルに置換（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6474)）
- Linux のファイルエクスプローラーコマンドでリテラルのパーセント記号を保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6470)）
- 前提条件とタスクコマンドのディレクトリを揃え、一貫して解決されるように変更（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6471)）
- デフォルトタスクをエグゼキューターのベースディレクトリで実行（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6469)）
- レガシーチェックによってキャッシュミスが隠れることを防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6468)）
- 整数オーバーフローの上限なく増分できるよう、CLI のプレリリース番号付けを修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6467)）
- キャッシュエントリーが null としてデコードされた場合のパニックを回避（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6466)）
- CLI でチェックごとにポートを一度だけ検査（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6465)）
- 古い wake 結果を回避するため、キャッシュ識別情報にタスクディレクトリを追加（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6464)）
- パッケージ化でファイルを同じファイル自身にコピーする場合のデータ損失を防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6463)）
- ステータスチェックでタスクディレクトリを考慮（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6462)）
- 変更履歴のエントリーにある TODO、FIXME、TBD プレースホルダーを検出（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6459)）
- CLI で埋め込みバンドルからランタイムアセットを生成（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6458)）
- git コマンドで、解析対象の HEAD 出力から stderr を分離（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6457)）
- 同時アクセス時もデフォルトのインメモリー SQLite データベースを維持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6456)）
- フォールバック用インデックスページへの HEAD リクエストに正しいレスポンスを返却（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6451)）
- SQLite ステートメントの準備とデータベースのクローズまたは再オープンとの競合を処理（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6450)）
- SQLite データベースのクローズまたは再オープン時に、準備済みステートメントのハンドルを解放（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6448)）
- gosod の使用時に展開先パスをリテラルのまま保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6447)）
- 行を途中まで反復処理した後に SQLite クエリーで発生したエラーを伝播（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6446)）
- kvstore でストア全体を読み取る際にスナップショットを返却（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6444)）
- JSON の null 値を読み込んだ後に kvstore で発生する書き込み時のパニックを防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6442)）
- 書式指定子を解釈せず、プレーンなターミナルメッセージを出力（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6441)）
- アップデーターで最新の安定版を検索する際のエラーを表示（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6440)）
- 変更履歴のコメントブロック内にある非表示エントリーと見出しを無視（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6439)）
- wails3 を一度だけインストールし、アップデーターでインストーラーのエラーを表示（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6438)）
- gosod で角括弧を含むリテラルのテンプレートパスを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6437)）
- アップデーターで GitHub タグの失敗レスポンスをデコード前に拒否（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6436)）
- ライブラリーの並列検索中にキャンセルエラーを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6435)）
- すべてのリターンパスで HTTP レスポンスボディを閉じ、リソースリークを防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6434)）
- 開発プロキシがリクエストを再試行する際にリクエストボディを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6432)）
- デスクトップのファイルエクスプローラーを検出する際に起動期限を遵守（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6430)）
- AppImage ツールのダウンロードを保存する前に HTTP エラーを拒否（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6429)）
- semver 比較で、過大な数値のプレリリース識別子を正しく順序付け（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6427)）
- アセットサーバーでサービスをルーティングする際にエスケープ済みパスを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6426)）
- アセットサーバーでバンドル済みランタイムへの HEAD リクエストを処理（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6425)）
- バンドル済みアセットが存在しない場合、誤ったレスポンスの代わりに 404 を返却（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6424)）
- リクエストがキャンセルされた場合、開発プロキシの再試行を停止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6423)）
- アセットサーバーで開始 URL の共有ミューテーションと古いポートの使用を回避（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6422)）
- アセットサーバーで任意指定のロガー設定を尊重（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6421)）
- ライブラリーパスのキャッシュ無効化におけるデータ競合を解消（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6418)）
- ldconfig の出力を検索する際の誤ったライブラリー一致を回避（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6416)）
- 負荷時にアセットサーバーのキューが増大すること、および拡大と縮小を繰り返すことを防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6415)）
- デバウンスタイマーの発火後にコールバックのキャプチャーを解放し、古い参照が保持されることを回避（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6413)）
- 実行可能ファイルのパス内で見つかったファイルエクスプローラーのターゲットを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6412)）
- ファイルエクスプローラーでデスクトップエントリを隠してしまうディレクトリをスキップ（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6411)）
- レスポンス本文が空の場合にアセットサーバーのヘッダーが送信されない問題を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6409)）
- アプリケーションアイコンの検索時に、XDG のデスクトップファイル検索ディレクトリを考慮（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6407)）
- ファイルエクスプローラーでのフィールド展開時にファイルパスを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6406)）
- ファイルエクスプローラーで展開された URI ターゲットが重複しないように修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6405)）
- 正常に完了した空のシーク不可アセットレスポンスを保留状態のままにせずコミット（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6403)）
- アセットサーバーで短い読み取りが複数回発生しても、MIME スニッフィング用プレフィックスを正しく充填（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6402)）
- アセットサーバーで、シーク不可の HEAD リクエストに対してレスポンス本文全体を読み込まないように修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6401)）
- アセットサーバーでディレクトリリクエストを解決する際、完全に一致するインデックスファイルのみを配信（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6400)）
- アクションをサポートしないサーバーで、閉じられた Linux 通知がクリックとして報告される問題を修正（@bercedev による [PR](https://github.com/wailsapp/wails/pull/6367)）
- ターミナルが閉じたときに開発プロセスを停止（@xmasdev による [PR](https://github.com/wailsapp/wails/pull/6371)）
- Windows タスクバー再起動のデバウンスフラグをクリアする前に、1 秒間保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6386)）
- アセットサーバーでフォールバックレスポンスのヘッダーが失われる問題を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6399)）
- アセットサーバーでフォールバックを選択する前に、ヘッダーの削除を反映（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6398)）
- メモリリークを防ぐため、アセットサーバーで削除済みのリングキュー参照を解放（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6397)）
- メモリ使用量を削減するため、アセットサーバーで使用済みのコンテンツスニッファーバッファーを解放（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6396)）
- アセット配信時に、大文字と小文字が混在するアセットファイル拡張子を認識（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6395)）
- ファイル名を保持せずに、アセットサーバーでコンテンツタイプを正しく検出（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6394)）
- デスクトップエントリにコマンドがない場合はデフォルトにフォールバック（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6392)）
- macOS 以外で、liquid-glass サンプルがアプリの起動前にダイアログを表示しないように修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6389)）
- Windows で制約付きサイズからフルスクリーンを開始した場合も、ウィンドウサイズの制約を保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6387)）
- Run より前に呼び出したダイアログが復帰せず永久にブロックする問題を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6390)）

## セキュリティ
- 過剰なメモリ使用を防ぐため、画像を割り当てる前に生成アイコンのサイズを検証（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6478)）
- リリースノート API のレスポンスを使用前に検証（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6461)）
- SQLite プラグインで完全な 16 進数のステートメントハンドルを必須化（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6452)）
- 意図しないライブラリの読み込みを防ぐため、ファイルシステム上のライブラリに対するフォールバックの一致範囲を制限（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6417)）

## v3.1.0-nightly.20261027 - 2026-10-27

## 変更
- 404 ページから、v2.wails.io 上の対応する Wails v2 ドキュメントページへリンク（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6366)）
- v2 ドキュメントサイトを v2.wails.io に移行（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6365)）
- v3 ドキュメントを wails.io から、v2 ドキュメントを v2.wails.io から配信（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6364)）
- v3.0.0 のリリースノートを更新し、セキュリティアドバイザリに関する記述を訂正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6363)）

## 修正
- Go のイベントリスナーコールバックを emit 順に直列実行（@DevLumuz による [PR](https://github.com/wailsapp/wails/pull/6275)）
- Run 後に作成したトレイが Linux で破棄されないよう、トレイの Menu プロパティを常に通知（@AlbinoGeek による [PR](https://github.com/wailsapp/wails/pull/6284)）
- Run より前に送信されたメインスレッドの処理を、アプリケーションが起動するまで保留（@AlbinoGeek による [PR](https://github.com/wailsapp/wails/pull/6287)）
- GitHub アップデーターでリリースアセットをダウンロードする際に、正しいアセット API URL を使用（@gavidroselj による [PR](https://github.com/wailsapp/wails/pull/6304)）
- DPI が混在する macOS および Linux 環境で、画面の境界と位置を修正（@julianstorer による [PR](https://github.com/wailsapp/wails/pull/6346)）
- Windows で WebView2 が Content-Length ヘッダーを省略した場合も、リクエスト本文を保持（@strraaa による [PR](https://github.com/wailsapp/wails/pull/6326)）
- ロック中にフックのスナップショットを作成し、イベントフックの発行時に発生するデータ競合を修正（@DevLumuz による [PR](https://github.com/wailsapp/wails/pull/6276)）
- Linux で `SystemTray.SetTooltip()` が何も実行しない問題を修正：StatusNotifierItem のツールチップ説明が設定されるようになりました（#6017）
- 週次リリースワークフローで変更履歴の根拠チェックを実行（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6374)）
- 実行結果を squash merge コミットに対応付けるよう、リリース変更履歴の根拠チェックを修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6373)）
- 暗いデスクトップでもカラースキームが変更されるまで Linux のダークモード検出がライトを返す問題と、GTK4 バックエンドで `linux:SystemThemeChanged` が一度も発火しない問題を修正しました。両バックエンドは、`org.freedesktop.appearance` ポータル名前空間を基盤とする単一の実装を共有するようになりました（#6072）
- フレーム付きウィンドウ（Windows、macOS、Linux）で、`Window.SetZoom` が 1.0 未満の係数を 1.0 に引き上げることはなくなりました。WebView2 の `ZoomFactor`、WKWebView の `magnification`、WebKitGTK のズームレベルは任意の正の値を受け入れ、`Zoom` ウィンドウオプションもすでにそのような値をそのまま渡していました。フレームレスウィンドウでは、ドラッグ領域とサイズ変更領域がズームに合わせて調整されないため、1.0 の下限を維持します（[#4590](https://github.com/wailsapp/wails/issues/4590)）。正でない係数、NaN、Inf は 1.0 に引き上げず、無視するようになりました — @APshenkin

## v3.0.0 - 2026-10-26

## 変更
- ナイトリープレリリースは次回の週次リリースのプレビューとして `vX.Y.0-nightly.YYYYMMDD` の形式で公開され、週次リリースではマイナーバージョンが上がります（`vX.Y.0`）。そのため、Go の `@latest` は常に週次リリースまたはホットフィックスリリースに解決されます。`@wailsio/runtime` npm パッケージのバージョンは常に Wails のリリースバージョンと一致し、ナイトリーには `nightly` dist-tag が付きます
- プロジェクトテンプレートでは、`@wailsio/runtime` をテンプレート生成に使用した CLI のバージョンに固定
- Wails v3 には Go 1.26 以降が必要です（Go 1.26.10、または 1.27.3 以降には `crypto/x509` のセキュリティ修正が含まれます）

## 修正
- `wails3 update cli` が v3 リリースを検出できない問題を修正：v2 タグのみを対象とし、誤ったパスからインストールしていました
- 一部の Windows ビルドでアプリケーションのメニューバーを描画すると、Windows アプリが `panic: SetTextColor failed` でクラッシュする問題を修正：メニューバーは有効なデバイスコンテキストに描画されるか、既定の描画にフォールバックするようになり、`UAHMENUITEM.Position` は C の `int` と一致し、メニューの描画でパニックが発生しなくなりました（#6311）
- 再起動待ちの間もステージング済みアップデートの安全性を維持：シャットダウン中にダウンロードまたは2回目の再起動が開始されても、スワップヘルパーがインストールを待っているペイロードが削除されたり競合したりしないようにしました（@julianstorer）
- `pkg/w32` の初期化処理が、序数で参照する関数ポインターの使用中に `uxtheme.dll` をアンロードしたため、Windows アプリが `applyDarkMode` でクラッシュする問題を修正：モジュールをプロセスの存続期間中ロードしたままにするようにしました（#6356）
- スワップ前にステージング済みペイロードを検証し、バックアップを破棄する前にインストール済みアプリケーションを検証することで、破損したアップデートが正常に動作するアプリケーションを置き換えず、ロールバックされるようにしました（@julianstorer）

## 削除
- [PR](https://github.com/wailsapp/wails/pull/5736) で最後の翻訳同期を行った後、Crowdin 連携を削除（@github-actions[bot]）

## セキュリティ
- すべてのフレームにガードスクリプトを挿入し、Linux（GTK3 および GTK4）上のクロスオリジン子フレームから送信される予約済み `wails:` メッセージをブロックすることで、GHSA-2p5f-2gx8-mh8m を緩和
- `github.com/klauspost/compress` を 1.18.7（GHSA-259r-337f-4rfw）に、`golang.org/x/image` を 0.45.0 に更新し、`crypto/x509` の修正（GO-2026-6809）を取り込むため Go 1.26.10 でビルド
- ランタイムパッケージ、セットアップウィザード、サンプル、プロジェクトテンプレートの `vite` を 8.2.4 以降に更新（GHSA-9jrq-w75r-8gcw、GHSA-rq7h-c2jc-7f22、GHSA-vfpm-58rq-9qcg、GHSA-r825-w8mm-5pxw、GHSA-587h-7wmj-h5pm、GHSA-698v-4r3x-3vcv）

### 既知の制限事項
- **Linux**：WebKitGTK はハンドラーにフレーム情報を提供しないため、子フレームガードでは、*同一オリジン*の子フレーム（たとえば `srcdoc` または `about:blank` の iframe）が親フレームを経由してネイティブメッセージハンドラーに到達するのを防止できません。クロスオリジンのフレームと、`allow-same-origin` が指定されていないサンドボックス化されたフレームはブロックされます。信頼できない同一オリジンのコンテンツを Wails ウィンドウに埋め込まないでください。
- **Android（試験的）：**`window.wails` JavaScript ブリッジは、メインフレームまたはオリジンの確認なしに WebView 内のすべてのフレームへ公開されます。クロスオリジンコンテンツを含む、埋め込まれた iframe 内のコンテンツは、バインド済みメソッドを呼び出し、予約済み `wails:` メッセージを送信できます。信頼できないサードパーティーコンテンツをアプリの WebView に埋め込まないでください。これは、子フレームからの `wails:` メッセージを無視する macOS および iOS とは異なります。3.0.1 での修正を予定しています。
- <strong>すべてのプラットフォーム：</strong>子フレームからのバインド済みメソッド呼び出しは、ネイティブメッセージハンドラーではなく HTTP 経由で `/wails/runtime` に送信されるため、ガードによる制限を受けません。埋め込むすべての iframe が、バインド済み Go メソッドへアクセスできるものとして扱ってください。

## v3.0.0-beta.34 - 2026-10-20

## 修正
- ウィンドウをアクティブにしても Windows の webview がフォーカスを取り戻さない問題と、Tab キーによるフォーカス移動が webview 内に閉じ込められ、ネイティブウィンドウとの間を行き来できない問題を修正（#3783）
- `internal/webview2/pkg/webview2` の COM vtable 構造体（`ICoreWebView2_2`..`_27` およびその他の派生インターフェイス）で継承されたメソッドチェーンが欠落し、呼び出しが誤った vtable スロットに到達する問題を修正（#6129）
- Windows の質問ダイアログでカスタムボタンのラベルが無視され、ボタンのコールバックが一度も実行されない問題を修正：カスタムボタンを TaskDialog で表示するようにし、comctl32 v6 が利用できない場合は、位置に基づくコールバックの対応付けを使用する MessageBox にフォールバックします（#5794）
- Vite に対する 60s で `wails3 dev` がハングする問題を修正：自動フロントエンド準備完了ゲートは、`Accept` ヘッダーを送信せず 404 を受信していた HTTP GET の代わりに、TCP プローブを使用するようになりました（[#6209](https://github.com/wailsapp/wails/issues/6209)）
- [PR](https://github.com/wailsapp/wails/pull/6301) で、システムトレイメニューを開いたまま更新しても動作し続けるように修正（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6302) で、ランタイム操作中の一時的な WebView2 COM エラーを適切に処理（@leaanthony）

## v3.0.0-beta.33 - 2026-10-19

## 追加
- バインディングジェネレーターが、コンパイル時定数の名前を使用して `application.RegisterEvent` に登録された検出可能な各イベントについて、型付きの名前定数をエクスポートする `bindings/events` モジュールを生成するようになりました（例：`user:login` に対する `UserLogin`）。Vite プラグインなしで `Events.On`/`Emit` とともに使用できます（#5645）

## 修正
- Windows の質問ダイアログでカスタムボタンのラベルが無視され、ボタンのコールバックが一度も実行されない問題を修正：カスタムボタンを TaskDialog で表示するようにし、comctl32 v6 が利用できない場合は、位置に基づくコールバックの対応付けを使用する MessageBox にフォールバックします（#5794）
- Vite に対する 60s で `wails3 dev` がハングする問題を修正：自動フロントエンド準備完了ゲートは、`Accept` ヘッダーを送信せず 404 を受信していた HTTP GET の代わりに、TCP プローブを使用するようになりました（[#6209](https://github.com/wailsapp/wails/issues/6209)）
- [PR](https://github.com/wailsapp/wails/pull/6301) で、システムトレイメニューを開いたまま更新しても動作し続けるように修正（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6302) で、ランタイム操作中の一時的な WebView2 COM エラーを適切に処理（@leaanthony）

## v3.0.0-beta.32 - 2026-10-18

## 追加
- macOS 14 以降向けにネイティブの macOS メニューセクション見出しを追加 — [ドキュメントサイト](https://v3.wails.io)を参照（20 ページを更新）。[PR](https://github.com/wailsapp/wails/pull/6297)（@leaanthony）

## 修正
- [PR](https://github.com/wailsapp/wails/pull/6299) で、macOS の SetRelativePosition が画面の作業領域を基準にするように修正（@leaanthony）
- [PR](https://github.com/wailsapp/wails/pull/6298) で、親プロセスのタイムアウトまたは置換の失敗後に、アップデーターヘルパーが変更されていない対象を再起動するように修正（@taliesin-ai）
- 10PR[ で、バージョン別ドキュメント用の pnpm をバージョン ](https://github.com/wailsapp/wails/pull/6293) に固定（@taliesin-ai）

## v3.0.0-beta.31 - 2026-10-18

## 追加
- ルートのランタイムエントリを評価せずに `CancellablePromise` をインポートできるよう、`@wailsio/runtime/cancellable` サブパスのエクスポートを公開（#6153）

## 変更
- Go 1.26 以降を必須化：v3 モジュールが `go 1.26.0` を宣言するようになり、`wails3 init` および `wails3 service init` で生成されるプロジェクトも同じ最小バージョンを宣言します

## 修正
- Go 1.27 でビルドすると、`Marshaler`、`Unmarshaler`、`RawMessage` などの `encoding/json` 型が `encoding/json/v2` 型および `encoding/json/jsontext` 型のエイリアスになり、バインディングジェネレーターの出力が変化する問題を修正：これらを再び `any` としてレンダリングし、`encoding/json/v2` のモデルを生成しないようにしました（#6292）
- `wails3 init -t` に渡されたリモートテンプレート URL をクローン前に検証：`https://`、`http://`、`ssh://`、`git://`、`user@host:path` の URL のみを受け入れ、値は `--` の後に git へ渡し、クローン中に git が使用できるトランスポートをこれらに制限します
- HTTP ランタイムトランスポートで、範囲外の `x-wails-window-id` ヘッダー値を別のウィンドウ ID に暗黙的に切り詰めず、拒否するように変更
- アップデーターのアーカイブ展開を zip-slip に対して強化：エントリパスを `filepath.IsLocal` で確認し、さらに結合後のパスが展開先に含まれることを確認するようにしました。zip および tar.gz アーカイブ内の、入れ子になった `..`、バックスラッシュ、絶対パス、シンボリックリンク経由の書き込みを行うエントリに対するテストも追加しました
- `wails3 generate appimage` から作成した AppImage が Debian および Ubuntu 以外で起動時に「Unable to spawn a new child process: … WebKitNetworkProcess」と表示してクラッシュする問題を修正。バンドルされた WebKit は、AppImage に同梱されているヘルパープロセスを実行するようになりました。WebKit の bubblewrap サンドボックスは引き続き有効です
- `-tags server` ビルドが Android および iOS 向けにコンパイルできない問題を修正し、[#6271](https://github.com/wailsapp/wails/issues/6271) を解決
- 一時ディレクトリが別のファイルシステム上（tmpfs の `/tmp` など）にある場合に Linux でアップデーターが入れ替えに失敗する問題を、インストール済みバイナリの隣でステージングし、ファイルシステム間でコピーすることで修正
- アップデーターが `.AppImage` ファイルではなく、AppImage の読み取り専用マウントを置き換えようとする問題を修正
- macOS および iOS で、予約済みの `wails:` ウィンドウ制御メッセージが iframe から受け入れられる問題を修正。これらはウィンドウのメインフレームからのみ受け入れられるようになりました。iframe からのその他のメッセージは、引き続きオリジン情報とともに `RawMessageHandler` に到達します
- macOS での `WindowDidUpdate` および `WindowWillUpdate` の送信を停止。`NSWindowDidUpdateNotification` は run loop の各処理サイクルで発火するため、ページは毎秒約 120 回の空の `evaluateJavaScript:` 呼び出しを受信していました。また macOS の 27.0 では、そのたびにユーザーのクリックによる一時的なアクティベーションが消費され、実際のクリックから呼び出した `navigator.clipboard.writeText()` やその他のアクティベーション必須 API が `NotAllowedError` で拒否されていました（[#6170](https://github.com/wailsapp/wails/issues/6170)）。@APshenkin による修正
- Explorer にバージョンの詳細が表示されるよう、Windows の info.json のデフォルト言語を en-US（0409）に修正（#6210）
- AMD または Intel の iGPU と NVIDIA dGPU を搭載したノート PC で、Linux の DMA-BUF レンダラーが無効になる問題を修正。Wayland では、WebKit が NVIDIA 上でレンダリングする場合（アプリが dGPU にオフロードされた場合、またはコンポジターがそこで動作している場合）にのみ無効になり、nvidia モジュールがロードされているだけでは無効にならなくなりました。これまでは iGPU 上のレンダリングが極端に低速になっていました。X11 セッションでは、Xorg が使用するレンダリング GPU はその構成で決まるため、従来の動作を維持します
- @john-okeefe による [PR](https://github.com/wailsapp/wails/pull/6135) で、Linux アップデーターのデバイス間 rename 失敗を、ステージングコピーへのフォールバックで処理
- Windows の `WebviewWindow.SetMenu` が、メニューバーのテーマを構成していない場合にメニューバーを再描画しない問題を修正し、[#6103](https://github.com/wailsapp/wails/issues/6103) を解決
- Windows の `WebviewWindow.SetMenu(nil)` が、macOS のように何もしないのではなくパニックする問題を修正し、[#6104](https://github.com/wailsapp/wails/issues/6104) を解決
- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/6094) で、競合する push が発生しても再試行によって changelog エントリを保持
- モジュールに一度も同梱されていないバイナリを参照する埋め込みを削除し、すべてのプラットフォームで `go mod vendor` が `pattern arm64/WebView2Loader.dll: no matching files found` により失敗する問題を修正。[#5782](https://github.com/wailsapp/wails/issues/5782) および [#5376](https://github.com/wailsapp/wails/issues/5376) を解決。@Grantmartin2002 による [PR](https://github.com/wailsapp/wails/pull/6031) で対応
- @charlesbalderson による [PR](https://github.com/wailsapp/wails/pull/6159) で、macOS におけるウィンドウの重複破棄を防止
- @wzk2025 による [PR](https://github.com/wailsapp/wails/pull/6162) で、右クリック時にシステムトレイメニューを開くための WM_CONTEXTMENU を処理
- ディスパッチされたコールバックがパニックした場合に、5 種類すべての `InvokeSync*` が呼び出し元を永遠にブロックする問題を修正し、[#6107](https://github.com/wailsapp/wails/issues/6107) を解決。@RALIST による [PR](https://github.com/wailsapp/wails/pull/6108) で対応
- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/6236) で、破棄された GTK3 サブメニューを安全に再構築
- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6261) で、ネストされた worker 向けの Windows worker 転送 shim を実装
- Windows の `OpenFile()` が `CanChooseDirectories(true)` の場合に `PromptForMultipleSelection()` を無視する問題を修正。macOS および Linux と同様に、複数のフォルダーを一度に選択できるようになりました（#6263）
- windows/arm64 における WebView2 COM メソッドへの値渡しの `double` 引数を修正。`ZoomFactor`、`RasterizationScale`、cookie の `Expires`、および印刷設定の倍率、用紙サイズ、余白が影響を受けていました（従来、印刷設定の setter は `ErrDoubleArgUnsupported` を返していました）。影響を受ける呼び出しは、Go の syscall 経路ではロードできない `d0` をロードする `runtime.cgocall` trampoline を経由するようになりました（golang.org/issue/62583）。@APshenkin による修正

## セキュリティ
- サーバーモード（`-tags server`）では、`/wails/runtime` 呼び出しを、サーバー自身のオリジンまたは `ServerOptions.WebSocketOriginPatterns` に指定されたオリジンからの `POST` リクエストとしてのみ受け入れるようになりました。その他のメソッドには HTTP 405、その他のブラウザーオリジンには HTTP 403 を返します。`Host` ヘッダーを書き換えるリバースプロキシの背後にデプロイする場合は、公開オリジンを `WebSocketOriginPatterns` に追加してください
- ループバックアドレス（デフォルトの `localhost`、`127.0.0.0/8`、または `::1`）にバインドされたサーバーモード（`-tags server`）では、`Host` ヘッダーが `localhost`、IP リテラル、バインド先ホスト、または `ServerOptions.WebSocketOriginPatterns` に一致するホストでない限り、HTTP 403 でリクエストを拒否するようになりました。これにより、`/wails/runtime` および WebSocket エンドポイントを標的とする DNS リバインディング攻撃を防ぎます
- Windows の通知では、WinRT toast API が失敗した場合に PowerShell へフォールバックしなくなり、代わりにエラーを返すようになりました

## v3.0.0-beta.30 - 2026-10-15

## 修正
- アップデーターが `.AppImage` ファイルではなく、AppImage の読み取り専用マウントを置き換えようとする問題を修正
- macOS での `WindowDidUpdate` および `WindowWillUpdate` の送信を停止。`NSWindowDidUpdateNotification` は run loop の各処理サイクルで発火するため、ページは毎秒約 120 回の空の `evaluateJavaScript:` 呼び出しを受信していました。また macOS の 27.0 では、そのたびにユーザーのクリックによる一時的なアクティベーションが消費され、実際のクリックから呼び出した `navigator.clipboard.writeText()` やその他のアクティベーション必須 API が `NotAllowedError` で拒否されていました（[#6170](https://github.com/wailsapp/wails/issues/6170)）。@APshenkin による修正
- Explorer にバージョンの詳細が表示されるよう、Windows の info.json のデフォルト言語を en-US（0409）に修正（#6210）
- @john-okeefe による [PR](https://github.com/wailsapp/wails/pull/6135) で、Linux アップデーターのデバイス間 rename 失敗を、ステージングコピーへのフォールバックで処理
- Windows の `WebviewWindow.SetMenu` が、メニューバーのテーマを構成していない場合にメニューバーを再描画しない問題を修正し、[#6103](https://github.com/wailsapp/wails/issues/6103) を解決
- Windows の `WebviewWindow.SetMenu(nil)` が、macOS のように何もしないのではなくパニックする問題を修正し、[#6104](https://github.com/wailsapp/wails/issues/6104) を解決
- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/6094) で、競合する push が発生しても再試行によって changelog エントリを保持
- モジュールに一度も同梱されていないバイナリを参照する埋め込みを削除し、すべてのプラットフォームで `go mod vendor` が `pattern arm64/WebView2Loader.dll: no matching files found` により失敗する問題を修正。[#5782](https://github.com/wailsapp/wails/issues/5782) および [#5376](https://github.com/wailsapp/wails/issues/5376) を解決。@Grantmartin2002 による [PR](https://github.com/wailsapp/wails/pull/6031) で対応
- @charlesbalderson による [PR](https://github.com/wailsapp/wails/pull/6159) で、macOS におけるウィンドウの重複破棄を防止
- @wzk2025 による [PR](https://github.com/wailsapp/wails/pull/6162) で、右クリック時にシステムトレイメニューを開くための WM_CONTEXTMENU を処理
- ディスパッチされたコールバックがパニックした場合に、5 種類すべての `InvokeSync*` が呼び出し元を永遠にブロックする問題を修正し、[#6107](https://github.com/wailsapp/wails/issues/6107) を解決。@RALIST による [PR](https://github.com/wailsapp/wails/pull/6108) で対応

## v3.0.0-beta.29 - 2026-10-13

## 追加
- macOS のグローバルショートカットがキーボードレイアウトを認識するようになりました。非 QWERTY レイアウト（AZERTY、QWERTZ、Dvorak など）では、固定された QWERTY の物理位置ではなく、`UCKeyTranslate` を介して文字アクセラレーターがキーに表示された文字にバインドされます。

## 修正
- @abhicodes-007 による [PR](https://github.com/wailsapp/wails/pull/6233) で、Linux 上でネイティブメニューを設定する前に初期化
- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6237) で、出力処理およびセットアップウィザードにおける macOS の GA 対応上のリグレッションを解消
- Windows の WebView2 COM `AddRef`/`Release` ラッパーが、古いスレッドの last-error を失敗として報告する問題を修正。参照カウントを返すようになったため、release によってアプリケーションが中断されたり、アセットリクエストが失敗したりすることはなくなりました（#3931）
- Windows で `ICoreWebView2Deferral.Release` が `AddRef` スロットを呼び出し、遅延されたすべてのアセットリクエストの遅延処理が解放されなかった問題を修正（#3931）

## v3.0.0-beta.28 - 2026-10-05

## 修正
- macOS アプリの出力をファイルプロバイダー外に保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6225)）
- セットアップウィザードを完成させ、システムの詳細を表示（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6224)）
- Windows での WebView2 復旧時のエラー処理とリソースのクリーンアップを改善（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6222)）
- アプリのクラッシュとコントローラーの取り残しを防ぐため、Windows での WebView2 の復旧処理を改善（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6220)）
- macOS のメディア認可に関するガイダンスを修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6221)）
- Linux のシステムトレイアイコンを右クリックすると、メニューが開くだけでなくクリックハンドラーも実行され、右クリックのたびに関連付けられたウィンドウの表示状態が切り替わる問題を修正（#6018）
- Windows でネイティブ側から解放されるまで WebView2 コールバックハンドラーを存続させるよう修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6184)）

## v3.0.0-beta.27 - 2026-10-01

## 修正
- macOS 12 以降で `Permissions` オプションが無視される問題を修正。`WKUIDelegate` メディアキャプチャメソッドが実装され、Linux および Windows と同様に、`PermissionAllow` と `PermissionDeny` がカメラおよびマイクのリクエストに適用されるようになりました。`NSCameraUsageDescription` / `NSMicrophoneUsageDescription` は引き続き必要であり、サンドボックス環境では対応するデバイスエンタイトルメントも必要です

## v3.0.0-beta.26 - 2026-09-25

## 追加
- Watcher がフロントエンド開発サーバーに HTTP 準備完了チェックを適用するよう変更（@atterpac による [PR](https://github.com/wailsapp/wails/pull/6150)）

## 修正
- WebView2 プロセスの障害時に空白のウィンドウを残さず復旧するよう修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6002)）

## v3.0.0-beta.25 - 2026-09-22

## 修正
- メインアプリケーションの終了時に `wails3 dev` とそのバックグラウンドプロセスを停止するよう修正（#6048）

## v3.0.0-beta.24 - 2026-09-20

## セキュリティ
- 最小権限の原則に従い、ビルドおよび公開ワークフローにおける GITHUB_TOKEN の権限を制限
- screen サンプルのアセットミドルウェアにおけるパストラバーサルとシンボリックリンクによる脱出を修正
- セットアップウィザードによる依存関係のインストール機能を維持しつつ、承認済みのパッケージマネージャーコマンドのみに制限

## v3.0.0-beta.23 - 2026-09-16

## 追加

- Linux で GTK アプリケーション ID を上書きできるよう対応（@overlordtm による [PR](https://github.com/wailsapp/wails/pull/5972)）

## 修正

- Linux のシングルインスタンス機能が、アプリ固有の ID でバス名を取得するよう修正（@overlordtm による [PR](https://github.com/wailsapp/wails/pull/5973)）
- Linux で中断されたカスタムスキームのリクエストをキャンセルするよう修正（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6126)）

## v3.0.0-beta.22 - 2026-09-14

## 追加

- 翻訳を完成させ、コミュニティメディアを移行。[ドキュメントサイト](https://v3.wails.io)を参照（1812 ページを更新）（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6122)）
- さまざまなスタイルとサイズの埋め込み可能な Wails ボタンを文書化（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6120)）

## v3.0.0-beta.21 - 2026-09-13

## 追加

- Wails v3 ドキュメントを M-Press で配信（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6116)）

## 修正

- 変更履歴の生成時に MPD フロントマターの JSON slug 値を解析するよう修正（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6118)）
- バックアップ失敗後、Updater がヘルパー環境変数を消去し、元のターゲットを再起動するよう修正（@cnmax による [PR](https://github.com/wailsapp/wails/pull/6080)）
- App.Run の実行中にデフォルトのシグナルハンドラーを開始するよう修正（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6098)）
- Windows のメニューが nil メニューを処理し、置き換えられたリソースを解放してメニューバーを再描画するよう修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6112)）
- 共有 YAML 設定を使用する新規プロジェクトで MSIX パッケージングを復元（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6115)）
- ジェネリックモデルの生成関数が後続のヘルパー宣言を参照すると、生成された JavaScript および TypeScript バインディングの読み込みに失敗する問題を修正し、相互依存するジェネリックモデルの作成時にスタックオーバーフローが発生しないように修正（#6062）

## v3.0.0-beta.20 - 2026-09-10

## 変更

- Clave のショーケースリンクを現在の Web サイトとリポジトリに更新（@01xR4in による [PR](https://github.com/wailsapp/wails/pull/6082)）

## 修正

- ナビゲーションをまたいで keepalive ハンドラーを維持しながら、ワーカーからのリクエストを含む、中断された Windows のアセットリクエストをキャンセルするよう修正。Apple プラットフォームでは、ネイティブのリクエストコンテキストをアプリケーションラッパー経由で転送するよう修正。（#5963、#5969）
- 競合する push が発生しても、再試行によって変更履歴のエントリを保持するよう修正（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6094)）
- モジュールに一度も同梱されていないバイナリを参照する埋め込みを削除し、すべてのプラットフォームで `go mod vendor` が `pattern arm64/WebView2Loader.dll: no matching files found` により失敗する問題を修正。[#5782](https://github.com/wailsapp/wails/issues/5782) および [#5376](https://github.com/wailsapp/wails/issues/5376) を修正（@Grantmartin2002 による [PR](https://github.com/wailsapp/wails/pull/6031)）

## 削除

- pure Go ローダーに置き換えられたネイティブ WebView2 ローダーのサポートを削除。これにより、埋め込まれていた `WebView2Loader.dll` バイナリと `github.com/jchv/go-winloader` 依存関係が削除されます。`native_webview2loader` ビルドタグは引き続き受け付けられ、エラーは発生しなくなりましたが、v3 ビルドには影響しません（@Grantmartin2002 による [PR](https://github.com/wailsapp/wails/pull/6031)）
- macOS API ガイドから未使用のビルドタグと FPS オプションを削除（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6097)）

## v3.0.0-beta.19 - 2026-09-09

## 追加

- オプトインで使用できるよう、非公開の macOS API をビルドタグで制御。[ドキュメント](https://v3.wails.io/features/browser/integration)、[ドキュメント](https://v3.wails.io/features/environment/info)、[ドキュメント](https://v3.wails.io/features/windows/basics)、[ドキュメント](https://v3.wails.io/features/windows/frameless)、[ドキュメント](https://v3.wails.io/features/windows/notch-windows)、[ドキュメント](https://v3.wails.io/features/windows/options)、[ドキュメント](https://v3.wails.io/guides/build/macos)、[ドキュメント](https://v3.wails.io/guides/build/private-macos-apis)、[ドキュメント](https://v3.wails.io/reference/overview)を参照（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6087)）

## 修正

- 64 MiB を超えるランタイムリクエストを HTTP 413 で拒否するよう修正（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6091)）

## セキュリティ

- トークン認証により MCP のオリジンとリモートアクセスを強化（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6092)）

## v3.0.0-beta.18 - 2026-09-08

## 修正

- ポインタレシーバーを使用して、Linux および Darwin での Calloc のメモリリークを修正（@4RH1T3CT0R7、[PR](https://github.com/wailsapp/wails/pull/6083)）

## v3.0.0-beta.17 - 2026-09-06

## 修正

- Windows：WebResourceRequested ハンドラー内の `GetRequest` が失敗するか nil でも、プロセスが終了しないように修正（`log.Fatal`／nil 参照パニック）。代わりにリクエストを破棄してログに記録（@midagedev、[PR](https://github.com/wailsapp/wails/pull/6006)）

## v3.0.0-beta.16 - 2026-08-29

## 変更

- 新しいターミナルウィンドウで公証用パスワードを入力するように変更（@leaanthony、[PR](https://github.com/wailsapp/wails/pull/6029)）

## 修正

- macOS でシステムトレイのクリック種別を正しく処理（@ChewbaccaCookie、[PR](https://github.com/wailsapp/wails/pull/5919)）
- CI が更新前に未使用の Microsoft apt リポジトリを削除するように変更（@Grantmartin2002、[PR](https://github.com/wailsapp/wails/pull/6041)）

## v3.0.0-beta.15 - 2026-08-27

## 修正

- WebView2 の埋め込みタイムアウトを 60 秒に延長（@Grantmartin2002、[PR](https://github.com/wailsapp/wails/pull/6043)）

## v3.0.0-beta.14 - 2026-08-26

## 修正

- macOS で Control＋文字キーの押下に正しい名前を付けるように修正（@taliesin-ai、[PR](https://github.com/wailsapp/wails/pull/6032)）
- Windows で ICO 形式のトレイアイコンを修正し、タスクバーのテーマに追従するように変更（@nik9play、[PR](https://github.com/wailsapp/wails/pull/6016)）

## v3.0.0-beta.13 - 2026-08-25

## 修正

- macOS でモーダルループの実行中もメインスレッドの処理を継続（@leaanthony、[PR](https://github.com/wailsapp/wails/pull/6026)）
- モバイルのセキュアストレージで失敗を返せるようにし、障害時には安全側に閉じるように変更（@mortenolsrud、[PR](https://github.com/wailsapp/wails/pull/5923)）
- リスナーが登録されていない場合でも、アプリケーションのイベントフックを実行するように修正（@archy-rock3t-cloud、[PR](https://github.com/wailsapp/wails/pull/5999)）
- コメントおよびローカライズ済みドキュメントの誤字を修正（@haoku123、[PR](https://github.com/wailsapp/wails/pull/6023)）
- `v3/examples` 以下にコミットされていたコンパイル済み macOS バイナリを削除（@4RH1T3CT0R7、[PR](https://github.com/wailsapp/wails/pull/6025)）

## v3.0.0-beta.12 - 2026-08-21

## 追加

- ライフサイクルとテレメトリーのサンプルを備えた macOS ノッチ通知ウィンドウを追加。[ドキュメント](https://v3.wails.io/features/windows/notch-windows)を参照（@leaanthony、[PR](https://github.com/wailsapp/wails/pull/6010)）
- 新しいオプションとネイティブ統合を備えた macOS NSPanel ウィンドウのサポートを追加。[ドキュメント](https://v3.wails.io/features/windows/options)を参照（@leaanthony、[PR](https://github.com/wailsapp/wails/pull/6008)）

## 修正

- SQLite Prepare が並行呼び出し時にハングする問題を修正（@archy-rock3t-cloud、[PR](https://github.com/wailsapp/wails/pull/5998)）

## v3.0.0-beta.11 - 2026-08-20

## 削除

- 廃止された実装トラッカーをドキュメントから削除（@leaanthony、[PR](https://github.com/wailsapp/wails/pull/6005)）

## v3.0.0-beta.10 - 2026-08-19

## 修正

- GTK4 Linux ホストがカスタムプロトコルおよびファイル関連付けから起動された際の引数を破棄する問題を修正（@midagedev、[PR](https://github.com/wailsapp/wails/pull/6000)）
- 変更履歴の検証で、削除された行と同一ソース内の訂正を正しく処理するように修正（@taliesin-ai、[PR](https://github.com/wailsapp/wails/pull/5993)）

## v3.0.0-beta.9 - 2026-08-16

## 追加

- エージェント支援によるプロジェクト管理向けに、安全な wails3 mcp サーバーを追加（@leaanthony、[PR](https://github.com/wailsapp/wails/pull/5896)）
- バインディング内のモデルに関するドキュメントを追加。[ドキュメント](https://v3.wails.io/features/bindings/models)を参照（@taliesin-ai、[PR](https://github.com/wailsapp/wails/pull/5988)）
- Atomic Linux システムでの rpm-ostree インストールをサポート（@leaanthony、[PR](https://github.com/wailsapp/wails/pull/5987)）
- ネイティブな週次スター履歴チャートの生成と公開を追加。[ドキュメント](https://v3.wails.io/credits)、[ドキュメント](https://v3.wails.io/de/credits)、[ドキュメント](https://v3.wails.io/fr/credits)、[ドキュメント](https://v3.wails.io/id/credits)、[ドキュメント](https://v3.wails.io/ja/credits)、[ドキュメント](https://v3.wails.io/ko/credits)、[ドキュメント](https://v3.wails.io/pt/credits)、[ドキュメント](https://v3.wails.io/ru/credits)、[ドキュメント](https://v3.wails.io/zh-cn/credits)、[ドキュメント](https://v3.wails.io/zh-tw/credits)を参照（@leaanthony、[PR](https://github.com/wailsapp/wails/pull/5986)）
- アプリケーションバンドルのリソースを解決する Darwin 専用の mac パッケージを追加。[ドキュメント](https://v3.wails.io/guides/build/macos)を参照（@leaanthony、[PR](https://github.com/wailsapp/wails/pull/5965)）
- Condui のショーケースページと索引項目を追加。[ドキュメント](https://v3.wails.io/community/showcase/condui)および[ドキュメント](https://v3.wails.io/community/showcase)を参照（@mgueregath、[PR](https://github.com/wailsapp/wails/pull/5962)）
- スクリーンショットとプロジェクトへのリンクを備えた Redis Viewer のショーケースページを追加。[ドキュメント](https://v3.wails.io/community/showcase)および[ドキュメント](https://v3.wails.io/community/showcase/redisviewer)を参照（@redisviewer、[PR](https://github.com/wailsapp/wails/pull/5984)）

## 変更

- Linux の GTK アプリケーションフラグを G<em>APPLICATION</em>NON_UNIQUE に更新（@overlordtm、[PR](https://github.com/wailsapp/wails/pull/5971)）
- 見つからないウィンドウイベントを警告レベルではなくデバッグレベルでログに記録するように変更（@julianstorer、[PR](https://github.com/wailsapp/wails/pull/5914)）

## 修正

- 登録済みの macOS キーボードショートカットを webview より優先できるように修正（@julianstorer、[PR](https://github.com/wailsapp/wails/pull/5902)）
- ドキュメントのサイドバーにある壊れたリンクを修復（@northes、[PR](https://github.com/wailsapp/wails/pull/5937)）
- WebKit が対応するカスタムスキームのタスクを中止したとき、macOS および iOS のアセットリクエストコンテキストをキャンセル（#5963）
- WindowSetFullscreenButtonEnabled メッセージを処理（@archy-rock3t-cloud、[PR](https://github.com/wailsapp/wails/pull/5976)）
- ビルドエラーを解消するため、preact-ts テンプレートで Fragment をインポート（@haoku123、[PR](https://github.com/wailsapp/wails/pull/5979)）
- アクティブなウィンドウまたはディスプレイが利用可能になる前に画面検出が実行された場合でも、従来の GTK3 サービス専用アプリケーションがクラッシュしないように修正（#5966）
- 未リリースの変更履歴が空でも、バージョンを明示したリリース処理を続行できるように修正（#5977）

## セキュリティ

- セキュリティアドバイザリに対処するため、Web サイトの nanoid ロックファイルを修正済みの 3.3.18 に更新（@taliesin-ai、[PR](https://github.com/wailsapp/wails/pull/5985)）

## v3.0.0-beta.8 - 2026-08-12

## 追加

- 変更履歴の自動生成エントリにドキュメント URL の生成を追加（@taliesin-ai、[PR](https://github.com/wailsapp/wails/pull/5957)）
- Streams を追加：WebSocket プログラミングモデルを使用し、待受ソケットを必要としない、Go と JavaScript 間の双方向バイトストリームです。Go では `app.HandleStream(name, handler)` を使用してストリームを宣言し、フロントエンドから `Stream(name)` を使用して接続します。これは `WebSocket` 形式のオブジェクトを返します。Go→JS はアセットサーバーを介してウィンドウごとに保持される単一のポーリングで、JS→Go は通常の POST で転送されます。TCP ポートには何もバインドされず、`evaluateJavaScript` も経由しません。サーバービルド（`-tags server`）では、代わりに同じハンドラーが実際の WebSocket を介して提供されるため、ビルドが異なってもアプリケーションコードは同一です。by @leaanthony
- メールボックスの変更履歴エントリを Unreleased に移動 [PR](https://github.com/wailsapp/wails/pull/5935) by @leaanthony

## 変更

- ドキュメントのサイドバー自動生成とブログ著者型の導出を更新 [PR](https://github.com/wailsapp/wails/pull/5938) by @leaanthony

## 修正

- WebView2 の初期化で期限とメッセージポンプを使用 [PR](https://github.com/wailsapp/wails/pull/5952) by @leaanthony
- WebView2 の Cookie テストは、明示的に有効化されない限り CI ではスキップし、実行を現在の OS スレッドに固定 [PR](https://github.com/wailsapp/wails/pull/5951) by @leaanthony
- Windows のメニュービルダーで、サブメニューの親項目のコマンド ID を復元 [PR](https://github.com/wailsapp/wails/pull/5944) by @gilad-ch
- 公式クロスコンパイルイメージを GTK 4.14 以降の Linux サポート基準に合わせる（#5928）
- 継承されたリンカーフラグを保持し、-ObjC を追加するように iOS Xcode プロジェクトを構成 [PR](https://github.com/wailsapp/wails/pull/5915) by @mortenolsrud
- 大規模なフロントエンドで `wails3 dev` アセットプロキシの TCP 接続が過剰に入れ替わり、ホストのエフェメラルポートを枯渇させ、無関係なプロセスが `EADDRNOTAVAIL` で失敗する可能性がある問題を修正
- 順序どおりのディスパッチとバックプレッシャーのため、ウィンドウごとのイベント JavaScript をキューに格納 [PR](https://github.com/wailsapp/wails/pull/5934) by @leaanthony

## 削除

- デスクトップバイナリのリリースパイプラインを削除：v3 リリースはタグのみとなり、`wails3` CLI は `go install` でインストールされます。`release-v3.yml` と、それをディスパッチしていた nightly ステップを削除 [PR](https://github.com/wailsapp/wails/pull/5946) by @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## 追加

- メディア再生時のユーザー操作要件を無効にする macOS の自動再生設定を追加 [PR](https://github.com/wailsapp/wails/pull/5512) by @Eyalm321
- メールボックスの変更履歴エントリを Unreleased に移動 [PR](https://github.com/wailsapp/wails/pull/5935) by @leaanthony

## 変更

- macOS のズームアニメーションで CADisplayLink または NSTimer を使用し、動作を滑らかに改善 [PR](https://github.com/wailsapp/wails/pull/5945) by @savely-krasovsky

## 修正

- 継承されたリンカーフラグを保持し、-ObjC を追加するように iOS Xcode プロジェクトを構成 [PR](https://github.com/wailsapp/wails/pull/5915) by @mortenolsrud
- 大規模なフロントエンドで `wails3 dev` アセットプロキシの TCP 接続が過剰に入れ替わり、ホストのエフェメラルポートを枯渇させ、無関係なプロセスが `EADDRNOTAVAIL` で失敗する可能性がある問題を修正
- 順序どおりのディスパッチとバックプレッシャーのため、ウィンドウごとのイベント JavaScript をキューに格納 [PR](https://github.com/wailsapp/wails/pull/5934) by @leaanthony

### 追加

- イベントを順序どおりに配信するための汎用非同期 FIFO メールボックスを実装 [PR](https://github.com/wailsapp/wails/pull/5851) by @savely-krasovsky and @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## 追加

- サイズ超過イベント向けの容量制限付きホスト側ストレージと、順序どおりの JavaScript 配信を実装 [PR](https://github.com/wailsapp/wails/pull/5930) by @leaanthony
- ウィンドウの点滅表示に macOS Dock のバウンスを実装 [PR](https://github.com/wailsapp/wails/pull/5921) by @julianstorer

## 修正

- アセットサーバーがフラッシュ時に Content-Type スニッファーのエラーと未書き込みのプレフィックスを保持するように修正 [PR](https://github.com/wailsapp/wails/pull/5931) by @leaanthony
- Wails のコールバックからアプリケーションメニューを置き換えたときに macOS アプリケーションがクラッシュするのを防止
- Windows 10 1809 / Windows Server 2019（ビルド 17763）でネイティブメニューが読めなくなる問題を修正。ダークモード用の uxtheme エクスポートがビルド 18334 以降に制限されていたため、これらのホストではアプリレベルのダークモード有効化が実行されていませんでした。その結果、メニュー背景は暗く描画される一方、Windows はライトテーマのままメニューテキストを描画し続け、暗い背景に暗いテキストが表示されていました。序数によるエクスポートは 17763 から存在するため、判定条件をそれに合わせました。
- `w32.GetStockObject` が `GetStockObject` ではなく `GetDeviceCaps` を呼び出し、すべての標準オブジェクトに対して 0 を返していた問題を修正。
- WebView2 ブートストラッパーのダウンロードエラー処理と報告を改善 [PR](https://github.com/wailsapp/wails/pull/5924) by @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## 修正

- macOS アプリのアクティベーションで、通常のアプリに限りアクティベーションポリシーを尊重するように修正 [PR](https://github.com/wailsapp/wails/pull/5897) by @julianstorer
- Linux ビルドで未初期化の GTK ウィンドウを保護 [PR](https://github.com/wailsapp/wails/pull/5898) by @julianstorer
- URL の読み込み前に Linux WebKit ウィンドウへ明示的な不透明背景色を設定 [PR](https://github.com/wailsapp/wails/pull/5899) by @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## 変更

- Android のビルドタスクで arm64 をデフォルトとし、deploy-emulator でホストアーキテクチャを選択 [PR](https://github.com/wailsapp/wails/pull/5890) by @mortenolsrud

## 修正

- ドラッグ中も macOS ウィンドウのズーム状態を維持し、動きを抑制 [PR](https://github.com/wailsapp/wails/pull/5900) by @leaanthony
- `webview_window_windows_nonclient.go` のビルド制約に `!server` を追加し、Windows のサーバーモードビルドを修正

## v3.0.0-beta.3 - 2026-08-03

## 追加

- 実装詳細にフェーズ 10 のベータ検証完了を記載 [PR](https://github.com/wailsapp/wails/pull/5881) by @leaanthony

## 修正

- Windows のダークモード API にウィンドウハンドルを渡し、引数を検証 [PR](https://github.com/wailsapp/wails/pull/5877) by @leaanthony
- フレームレスウィンドウ向けの macOS タイトルバーボタン状態の決定処理を一元化 [PR](https://github.com/wailsapp/wails/pull/5870) by @taliesin-ai
- Windows のアプリテーマがライトであるにもかかわらず、Windows アプリケーションがダークモードを要求した場合に、ネイティブメニューのテキストが読めなくなるのを防止。Windows がダークテーマ用のメニューテキストを描画できるようになるまでは、メニューで対応するライトテーマのネイティブ背景を使用します。
- Windows 10 1809 / Windows Server 2019（ビルド 17763）でネイティブメニューが読めなくなる問題を修正。ダークモード用の uxtheme エクスポートがビルド 18334 以降に制限されていたため、これらのホストではアプリレベルのダークモード有効化が実行されていませんでした。その結果、メニュー背景は暗く描画される一方、Windows はライトテーマのままメニューテキストを描画し続け、暗い背景に暗いテキストが表示されていました。序数によるエクスポートは 17763 から存在するため、判定条件をそれに合わせました。

## v3.0.0-beta.2 - 2026-08-02

## 変更

- v3 をアルファからベータに昇格
- システムトレイのスマートデフォルトとポップアップの自動非表示動作を文書化し、クリックハンドラー選択の回帰テストを追加（#5840）。
- GitHub アップデーターで Windows インストーラーのアセットをデフォルトで除外 [PR](https://github.com/wailsapp/wails/pull/5861) by @leaanthony
- macOS のフレームレスウィンドウで、角丸、直角、およびカスタム半径のコーナーをサポート（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5866)）

## 修正

- GTK4 ウィンドウの現在のサイズを報告し、設定されたサーフェスからリサイズ、最大化、最小化、およびフルスクリーンの状態イベントを発行するようにしました（#5830）。
- fetch リクエストで Blob または FormData を送信した際に Linux の WebKit がクラッシュする問題を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5854)）
- Blob/FormData のヘッダーがない場合、fetch shim が undefined を渡すように変更（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5865)）

## v3.0.0-alpha2.122 - 2026-08-01

## 追加

## 変更

- macOS のフレームレスウィンドウで、角丸、直角、およびカスタム半径のコーナーをサポート（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5866)）

## 修正

- Blob/FormData のヘッダーがない場合、fetch shim が undefined を渡すように変更（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5865)）

## v3.0.0-alpha2.121 - 2026-07-31

## 追加

- 新しいオプションとビルドタスクを備えた macOS DMG パッケージングのサポートを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5857)）

## 変更

- GitHub updater がデフォルトで Windows インストーラーのアセットを除外するように変更（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5861)）

## 修正

- fetch リクエストで Blob または FormData を送信した際に Linux の WebKit がクラッシュする問題を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5854)）

## v3.0.0-alpha2.120 - 2026-07-31

## 追加

- macOS のタイトルバーをダブルクリックした際に最大化または最小化するアクションを実装（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5853)）

## 変更

- QR サービスのチュートリアルで NewServiceWithOptions を使用し、間隔を追加するように更新（@jeongkyu による [PR](https://github.com/wailsapp/wails/pull/5849)）

## 修正

- macOS でのズーム中も WKWebView の応答性を維持するように修正（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5856)）
- GTK4 のウィンドウサイズ取得を修正し、設定された `GdkSurface` からリサイズ、最大化、最小化、およびフルスクリーンの状態イベントを発行するようにしました。

## v3.0.0-alpha2.119 - 2026-07-27

## 修正

- 複数言語のドキュメントにアーキテクチャ図を追加（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5833)）

## v3.0.0-alpha2.118 - 2026-07-26

## 追加

- アイコン生成の入力および出力にデフォルトパスを設定（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5825)）
- runtime の package.json の sideEffects にソースエントリーモジュールを追加（@savely-krasovsky による [PR](https://github.com/wailsapp/wails/pull/5797)）
- コントリビューションガイドにライセンスと来歴のセクションを追加（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5816)）

## 修正

- GTK4 のフレームレスウィンドウにスコープを限定した CSS を適用し、角丸を削除（@savely-krasovsky による [PR](https://github.com/wailsapp/wails/pull/5800)）
- Windows のポップアップメニューと画面列挙で、カーソル位置の取得失敗を適切に処理するように修正（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5789)）
- macOS のファイルを開くダイアログで拡張子を正しくフィルタリングし、許可されたファイルを接尾辞で検証するように修正（@phergul による [PR](https://github.com/wailsapp/wails/pull/5678)）
- Windows のダークモード初期化で nil の API が呼び出されないように保護（@roachadam による [PR](https://github.com/wailsapp/wails/pull/5793)）
- updater の 32 ビットビルド失敗を修正しました。`GOARCH=386` 上で `fmt.Errorf` に渡した際、`maxArchiveTotalSize` 定数（2 GiB）がプラットフォームの `int` でオーバーフローしていました。現在は明示的に `int64` 型が指定されています。
- Windows 10 1809 / Windows Server 2019（ビルド 17763）など、ダークモード用 uxtheme API を読み込まない Windows ビルドで、ウィンドウが Dark（またはシステムのダークモード）のタイトルバーを使用している場合に、起動時に nil ポインターによるパニックが発生する問題を修正しました。ウィンドウテーマ設定内の `AllowDarkModeForWindow` 呼び出しに、`w32.SetMenuTheme` ですでに使用されているものと同じ nil ガードを追加しました。

## v3.0.0-alpha2.117 - 2026-07-08

## 追加

- Windows の非クライアント領域にカスタムヒットテストロジックを実装（@savely-krasovsky による [PR](https://github.com/wailsapp/wails/pull/5462)）

## 変更

- UseVisualHosting に基づいて WebView2 のモニタースケール検出を設定（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5761)）

## v3.0.0-alpha2.116 - 2026-07-07

## 追加

- Wails v3 の機能とガイダンスを中心とするよう FAQ ドキュメントを更新（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5763)）

## v3.0.0-alpha2.115 - 2026-07-06

## 修正

- GTK4 Linux で `Menu.Update()` がネイティブメニューを再構築しない問題を修正（#5659、#5539 で @puneetdixit200 が独自に診断および修正）
- ディスプレイ変更時に macOS の画面を列挙するとクラッシュする問題を、画面の ID/名前文字列をコピーして画面数のスナップショットを取得することで修正（#5565、#5584 で @x-haose が独自に診断および修正）
- 最小化からの復元中に `GetClientRect` が nil を返し、`WM_ERASEBKGND` が単色の背景を描画すると Windows でクラッシュする問題を修正（#5636 で @sinspired がガードを報告）
- フロントエンドバインディングのエラーが常にテキストとして解析される問題を修正（#5690、@mbaklor）
- `server` ビルドタグの使用時に Windows のビルドが失敗する問題を修正しました。原因は、Windows の GUI ファイルに、macOS および Linux の同等ファイルにはすでに存在する `!server` ビルド制約がなかったことでした（#5680）。

## v3.0.0-alpha2.114 - 2026-07-05

## 追加

- Update Manifest プロトコルとエンドポイントプロバイダーを実装（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5720)）

## 変更

- `webview2` バインディングを `v3/internal/webview2` として v3 モジュールに統合し、単独モジュール、その nightly リリース／同期ワークフロー、および go.mod のバージョン調整を削除しました（唯一の利用元は v3）（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5711)）。

## 修正

- WebView2 のモニタースケール検出と、DPI 変更時のホスト再同期に関する修正を Unreleased セクションに移動（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5750)）
- float64 および BOOL パラメーターに対する WebView2 COM マーシャリングを更新（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5741)）
- Windows のシステムトレイアイコンの更新および破棄で、パニックと nil 参照が発生しないように修正（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5703)）
- Windows で非表示ウィンドウが正しく再非表示にならない問題を修正（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5743)）
- ウィンドウの最小化、最大化、および復元に WebView2 コントローラーの表示状態を同期（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5742)）

### 修正

- WebView2 のモニタースケール検出を再度有効にし、DPI 変更時にのみホストを再同期するよう制御。@randalmurphal が検証した修正を基に、@eleclin が根本原因を確認し、@qq540491950 が実機でテスト。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5734)

## v3.0.0-alpha2.113 - 2026-07-04

## 追加

- `ANDROID_KEYSTORE_FILE` が設定されていない状態でリリース AAB をビルドすると警告するようにし（Google Play はデバッグ署名されたバンドルを拒否）、App Bundle のパッケージ化と署名について文書化。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5730)
- 「Why Wails」ドキュメントを複数言語で追加。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5739)
- バインディングで Go の time.Time を JS の Date または string にマッピングできるよう対応。@fbbdev による [PR](https://github.com/wailsapp/wails/pull/5398)
- Play Store への提出用に Android App Bundle（AAB）のパッケージ化タスク（`bundle`、`bundle:fat`、`assemble:aab`、`assemble:aab:release`）を追加。ローカル／エミュレーターでのテスト用 APK タスクは引き続き利用可能。@mortenolsrud による [PR](https://github.com/wailsapp/wails/pull/5728)（[#5726](https://github.com/wailsapp/wails/issues/5726) を修正）
- Android 実機向けのタスクターゲットを追加し、カメラ／位置情報の権限を復元。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5735)

## 変更

- `webview2` を v1.0.28 に更新（[リリースノート](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)）。
- Android テンプレートの `compileSdk`/`targetSdk` を 34 から、新規アプリの提出時に Google Play が要求する 35 へ更新。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5730)

## 修正

- sponsorkit で焼き込まれたアバターマスクを修正。@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5745)
- 辞書順でバージョンを並べ替えていたため、Android AVD の自動作成時に誤ったシステムイメージまたは cmdline-tools のバージョンが選択される問題を修正。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5730)
- セットアップウィザードが古い Android NDK バージョンを提案する問題を修正（現在は文書化された要件に一致する 26.3.11579264）。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5730)
- SvelteKit とオプションに関するフランス語ドキュメントを更新。@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5744)
- ディスプレイ変更中に macOS の画面列挙で発生する SIGSEGV を修正。@flofreud による [PR](https://github.com/wailsapp/wails/pull/5516)

## v3.0.0-alpha2.112 - 2026-07-03

## 追加

- Go ベースのコントリビューター SVG ジェネレーターを追加し、ドキュメント／ウェブサイトのクレジットページを更新。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5724)
- バインディングで Go の time.Time を JS の Date または string にマッピングできるよう対応。@fbbdev による [PR](https://github.com/wailsapp/wails/pull/5398)

## 変更

- Node ベースのスポンサー画像パイプラインを Go ジェネレーターに置き換え。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5719)

## 修正

- Android ビルドアセットの依存関係インストールスクリプトを修正。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5729)
- `ValidateAndSanitizeURL` で U+0085（NEXT LINE）制御文字を拒否するようにし、URL バリデーターによる空白文字の網羅を完了
- フレームレスウィンドウで DPI が変更された際に DWM フレームを再計算するよう修正。@leaanthony による [PR](https://github.com/wailsapp/wails/pull/4785)
- Windows でスケーリングが 100% 以外の場合に DnD ドロップゾーンの検出が失敗する問題を修正。@yulesxoxo による [PR](https://github.com/wailsapp/wails/pull/4632)
- Darwin のダイアログ、メニュー、トレイ、通知全体で Cocoa オブジェクトを明示的にメモリ管理するよう対応。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5714)
- Linux の CGO バックエンドのバグとシステムトレイの問題を修正。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5718)

## v3.0.0-alpha2.111 - 2026-07-01

## 追加

- コミュニティショーケースに HappyTools を追加。@Aliuyanfeng による [PR](https://github.com/wailsapp/wails/pull/5061)
- インドネシア語ロケールのサポートと包括的なドキュメントを追加。@triadmoko による [PR](https://github.com/wailsapp/wails/pull/5643)
- WindowsWindow に DisableMenu オプションを追加。@leaanthony による [PR](https://github.com/wailsapp/wails/pull/4813)

## 変更

- GOOS と ARCH を指定してビルド／パッケージタスクをディスパッチするよう、Taskfile テンプレートと CLI を更新。@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5617)

## 修正

- macOS のウィンドウタブ化に関する問題を修正。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5708)

## 削除

- コントリビューション、機能、ガイドからドイツ語翻訳済み MDX ファイルを削除。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5702)

## v3.0.0-alpha2.110 - 2026-06-30

## 追加

- macOS WebView の再読み込みと強制再読み込みを実装し、WebContent プロセス終了時の復旧処理を追加。@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5129)
- コントリビューション、機能、ガイドに関する包括的なドイツ語ドキュメントを追加。@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5396)
- 通知にサウンド、添付ファイル、スケジュール設定、更新 API を追加して機能を強化。@popaprozac による [PR](https://github.com/wailsapp/wails/pull/5333)

## 修正

- フレームレスウィンドウで DPI が変更された際に DWM フレームを再計算するよう修正。@leaanthony による [PR](https://github.com/wailsapp/wails/pull/4785)
- Windows でスケーリングが 100% 以外の場合に DnD ドロップゾーンの検出が失敗する問題を修正。@yulesxoxo による [PR](https://github.com/wailsapp/wails/pull/4632)

## v3.0.0-alpha2.109 - 2026-06-29

## 追加

- EventsEmit ドキュメントにコードサンプルを追加。@iamhabbeboy による [PR](https://github.com/wailsapp/wails/pull/5026)
- Windows WebView2 のビジュアルホスティングオプションを追加。@MerIijn による [PR](https://github.com/wailsapp/wails/pull/5380)
- コミュニティショーケースのドキュメントに Klustr を追加。@SametKUM による [PR](https://github.com/wailsapp/wails/pull/5536)
- 新しいページと変更履歴エントリを含め、コミュニティショーケースに Kira を追加。@thiennguyen93 による [PR](https://github.com/wailsapp/wails/pull/5685)
- MCP サービスガイドにフィードバックセクションを追加。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5694)

## 変更

- サーバーモードで、デスクトップのビルドタスクと一貫した正式なプロダクションビルドを利用できるようになりました（#5693）。`task build:server` はデフォルトでプロダクションバイナリ（`-tags server,production`、`-trimpath`、ストリップ済み）をビルドし、`DEV=true`（開発サーバー）、`OBFUSCATED=true`（garble）、`EXTRA_TAGS` を受け付けます。`task run:server` は開発サーバーを実行します。`Dockerfile.server` / `task build:docker` は、最初にプロダクションサーバー（`-tags server,production`）とプロダクションフロントエンドをビルドします。イメージはデフォルトで distroless/static 上の Pure Go 静的ビルドとなり、CGO アプリ向けに `CGO_ENABLED`、`GO_IMAGE`、`RUNTIME_IMAGE` が上書き可能なビルド引数として公開されます。

## 修正

- 非同期呼び出しの保留中にウィンドウを閉じた際のクラッシュを防止。@leaanthony による [PR](https://github.com/wailsapp/wails/pull/4435)
- Windows で非表示のアプリを開いた際にウィンドウがアクティブになるのを防止。@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5249)
- WebKit のリクエストメタデータ、レスポンス完了処理、ボディストリーム処理が GTK メインスレッドで実行されるよう修正。@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5668)
- GTK4 Linux で `Menu.Update()` がネイティブメニューを再構築しない問題を修正（#5659、@puneetdixit200 も #5539 で独自に診断・修正）
- ディスプレイ変更時に macOS の画面を列挙するとクラッシュする問題を、画面の ID／名前文字列をコピーして件数のスナップショットを取得することで修正（#5565、@x-haose も #5584 で独自に診断・修正）
- Windows で異なる DPI のモニター間をまたいでウィンドウをドラッグした後、WebView2 のコンテンツが縮小して消える問題を修正。最小化解除時の DPI 再同期と同様に、`WM_DPICHANGED` ハンドラーでコントローラーの境界を再設定するようにしました（#5677）

## v3.0.0-alpha2.108 - 2026-06-28

## 追加

- `app.GlobalShortcut` によるグローバル（システム全体）キーボードショートカットを追加（`Register`、`Unregister`、`UnregisterAll`、`IsRegistered`、`GetAll`）。アプリケーションにフォーカスがない場合でもショートカットが発動します。サードパーティ依存なしで、プラットフォームごとにネイティブ実装されています。macOS では Carbon ホットキー、Windows では `RegisterHotKey`、X11 では `XGrabKey`、Wayland では XDG Desktop Portal のグローバルショートカットインターフェースを使用します。
- 組み込み MCP サーバーを追加。`mcp` タグ付きでアプリケーションをビルドすると自動的に起動する Model Context Protocol サーバーで、LLM エージェントが実行中の Wails アプリケーションをテストおよび操作できます。ウィンドウ操作、DOM 検査、JavaScript の評価、バインド済みメソッドの呼び出し、イベント、画面上にアニメーション表示されるカーソルを使ったマウス／キーボード入力のシミュレーションに対応します。ユーザーコードは不要です。`WAILS_MCP=1` が設定されている場合、`wails3 build`／`wails3 dev` によって `mcp` タグが自動的に追加されます。設定はすべて環境変数（`WAILS_MCP_HOST`、`WAILS_MCP_PORT`、`WAILS_MCP_TIMEOUT`、`WAILS_MCP_HIDE_CURSOR`）で行います。

## 修正

- GTK4 Linux で `Menu.Update()` がネイティブメニューを再構築しない問題を修正（#5659、@puneetdixit200 も #5539 で独自に診断・修正）
- ディスプレイ変更時に macOS の画面を列挙するとクラッシュする問題を、画面の ID／名前文字列をコピーして件数のスナップショットを取得することで修正（#5565、@x-haose も #5584 で独自に診断・修正）

## v3.0.0-alpha2.107 - 2026-06-27

## 追加

- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5613) で、サイドバーナビゲーションを備えた実験的な Wake ドキュメントを追加

## v3.0.0-alpha2.106 - 2026-06-24

## 変更

- `webview2` を v1.0.27 に更新。
  - ci(webview2): リリースビルドを修正（Windows のクロスコンパイル＋完全な go.sum）（#5671）\

  **差分全体：** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5672) で、webview2 リリースワークフローのクロスコンパイルから go vet を削除
- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5670) で、自動変更履歴生成に使用する OpenRouter モデルを google/gemini-2.5-flash-lite に更新
- `webview2` を v1.0.26 に更新。

### 修正

- **一時的なランタイム COM エラーで終了せず、復旧するように変更**（#5658、#5580）。以前の `Chromium.errorCallback` は、*あらゆる* COM エラーに対して `os.Exit(1)` を呼び出していたため、起動後に復旧可能な一時的不具合が発生しただけでもアプリケーション全体が終了していました。ランタイムの処理経路（`Resize`／`GetClientRect`、`Navigate`／`NavigateToString`、`Init`、`MessageReceived`、`PutZoomFactor`、`OpenDevToolsWindow`）では、エラーをログに記録して復旧するようになりました。特に、`MessageReceived` 内の不正な形式または信頼できない Web メッセージは、プロセスを停止させずに破棄されるようになりました。これにより、異なる DPI のモニター間をまたぐ際に発生する一連のクラッシュ（#5544、#5650）に対処します。環境およびコントローラーの作成経路で発生するエラーは、引き続き致命的エラーとして扱われます。\

**差分全体：** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## 修正

- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5671) で、go.sum ファイルを正しく処理するよう release-webview2 ワークフローを修正
- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5659) で、ネイティブメニューをクリアして再構築することにより Linux GTK4 のメニュー更新を修正

## v3.0.0-alpha2.105 - 2026-06-21

## 追加

- 共有コードから実行時にプラットフォームを検出するための `application.System` を追加。`System.IsMobile()`（iOS／Android）、`System.IsDesktop()`（macOS／Windows／Linux）、`System.IsServer()`（`server` ビルドタグ）、および単一のターゲットを直接テストするための `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` を提供します。すべてのターゲットでコンパイルできるため、ビルドタグを使わずに分岐できます。対応するフロントエンド用ヘルパー（`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`）は `@wailsio/runtime` で利用できます
- 独自の Vite プロジェクトを `frontend/` に組み込む方法を示す「その他のフロントエンドフレームワークの使用」ガイドを追加（Solid、Preact、Lit、SvelteKit、Qwik、Angular などを対象）
- `wails3 setup` ウィザードで、モバイル（iOS／Android）ツールチェーン（Xcode と iOS Simulator ランタイム、JDK、Android SDK／NDK、エミュレーター）を確認するようになりました。該当する場合は、ワンクリックでのインストールと、コピー可能なシェル設定の修正も利用できます
- 生成されるプロジェクトには、新たに公開された（侵害されている可能性がある）パッケージへの露出を減らすため、7 日間の `minimum-release-age` を設定する `frontend/.npmrc` が同梱されます（pnpm と bun では適用され、npm では問題なく無視されます）

## 変更

- すべての組み込みスターターテンプレートを、新しいネオンマウンテンのヒーロービジュアル（Web、iOS、Android）で再設計
- **スターターテンプレートでは TypeScript がデフォルトになり、サフィックスなしのテンプレート名に割り当てられました。** `wails3 init`（`-t` なし）は TypeScript プロジェクトを生成します。`-t vanilla`、`-t react`、`-t vue`、`-t svelte` は TypeScript 版で、JavaScript 版はそれぞれ `-t vanilla-js`、`-t react-js`、`-t vue-js`、`-t svelte-js` です。組み込みテンプレートでは、`template.yaml` 内の `typescript:` で言語を宣言します。`-ts` サフィックスを使用するコミュニティテンプレートも、フォールバックとして引き続き動作します
- `wails3 setup` ウィザードを、ネオン調の「digital Wails」テーマ（山の背景に重ねたすりガラス風の鮮やかな表現）で再設計

## 修正

- WebView2 がサスペンドするか、そのレンダリング／GPU プロセスが再作成されるほど長時間最小化されていたアプリを復元すると、Windows でクラッシュする問題を修正。最小化／復元時の DPI 再同期（#5544）では、ウィンドウの DPI が実際に変更された場合にのみ WebView2 コントローラーを操作するようになり、一般的な同一 DPI での復元時に、サスペンド中のコントローラーに対して致命的な COM 呼び出しが行われることを回避します（#5605）
- アセット／メディアが頻繁に読み込まれる長時間稼働の Linux アプリで、ネイティブの `SIGABRT`／`SIGSEGV` クラッシュ（通常は GTK メインループ中の `g_object_unref` 内で発生）が繰り返し発生する問題を修正。アセットサーバーはワーカー goroutine から `WebKitURISchemeRequest` を完了させ、スレッドセーフでない WebKit2GTK 関数を GTK メインスレッド以外から呼び出していました。完了処理（`webkit_uri_scheme_request_finish_with_response`／`finish_error`）はメインスレッドで実行されるようになりました。#5566 の部分的な修正を完了するものです。GTK3 と GTK4／WebKitGTK の 6.0 ビルドの両方に影響します（#5631、#5557）
- Linux／GTK3 の `setupSignalHandlers` で断続的に発生する `fatal error: invalid pointer found on stack` を修正。シグナルの `user_data` として渡されるウィンドウ ID が Go のローカル変数 `unsafe.Pointer` に保持されていたため、スタックのコピー中にガベージコレクターがその（ポインターではない）値をスキャンすると異常終了していました。Go 側では ID を整数型（`uintptr_t`）のまま保持するようになりました。これは、#4958 で GTK4 の処理経路に適用されたものと同じ修正を従来の GTK3 の処理経路にバックポートするものです（GTK4 では、`-race`／checkptr エラーを解消するため、C のシグナル関数を `uintptr_t` に変更しました）（#5631）

## 削除

- `react-swc`、`preact`、`lit`、`solid`、`qwik`、`sveltekit` のスターターテンプレート（および各 `-ts` 版）を削除。サポートされる組み込みテンプレートは `vanilla`、`react`、`vue`、`svelte` になりました。いずれもデフォルトは TypeScript で、`-js` JavaScript 版も用意されています。その他のフレームワークも、[独自のフロントエンドを持ち込む](https://v3.wails.io/guides/dev/frontend-frameworks)か、カスタムテンプレートを使用することで引き続き利用できます

## v3.0.0-alpha2.104 - 2026-06-18

## 修正

- バインドされた Go サービスメソッドが空文字列を返すと iOS でクラッシュ（SIGABRT）する問題を修正。iOS のアセットレスポンスライターは、本文の長さではなく `buf != nil` で本文ポインターをガードしていたため、長さがゼロの本文では `&buf[0]` がパニックしていました。デスクトップのライターと同様に、長さでガードするようになりました

## v3.0.0-alpha2.103 - 2026-06-15

## 変更

- iOS と Android のネイティブ機能をプラットフォームマネージャーへ移行。従来の `application.IOS*`/`application.Android*` フリー関数ではなく、`application.IOS.*` と `application.Android.*`（例：`application.IOS.Haptic("medium")`、`application.Android.Share(payload)`）を介して呼び出すように変更（#5602）
- モバイルブリッジイベントの名前を変更。クロスプラットフォームイベントは `common:*` プレフィックス（例：`common:haptic`、`common:location`）を使用し、プラットフォーム固有イベントは `ios:*` / `android:*`（例：`ios:backgroundTask`、`android:foregroundService`）を使用するようになりました。`native:*` プレフィックスは使用されなくなりました（#5602）

## v3.0.0-alpha.102 - 2026-06-14

## 追加

- 対話形式のプロジェクトセットアップと依存関係チェックを行う実験的な `wails3 setup` ウィザードを追加
- 機械可読な出力を行う `--json` フラグを `wails3 doctor` に追加
- `wails3 doctor` コマンドに署名ステータスのセクションを追加

## 修正

- Linux で npm を検出する際、パッケージマネージャーに加えて PATH も確認するように修正

## v3.0.0-alpha.101 - 2026-06-13

## 追加

- iOS：ネイティブのメッセージダイアログ（UIAlertController）と、ファイル／複数ファイル／ディレクトリを開くダイアログ（UIDocumentPickerViewController）を追加。保存ダイアログは明示的なエラーを返します
- iOS：UIPasteboard によるクリップボード対応
- iOS：UIScreen による実際の画面メトリクス（ポイント、ピクセル、スケール、セーフエリア内の作業領域）
- iOS：実機向けビルド（`IOS_PLATFORM=device`）、コード署名 ID／プロビジョニングプロファイル／エンタイトルメントへの対応、`.ipa` パッケージング、および devicectl を介した `deploy-device`
- iOS：設定可能な最小 iOS バージョン（build/config.yml の `ios.minIOSVersion`）
- iOS：`wails3 doctor` が macOS 上の Xcode と iOS SDK の利用可否を報告
- iOS：バッテリー、ネットワーク、テーマ、画面ロック、メモリ不足のシステムイベントを、`events.IOS.*` およびプラットフォームに依存しない `events.Common.*` アプリケーションイベントとして公開
- iOS：ネイティブモバイル機能ブリッジ（公開された `application.IOS*`）— 共有シート、URL を開く、スリープ抑止、トーチ、セーフエリアのインセット、明るさ、アプリ情報、画面方向のロック、ステータスバー、生体認証（Face ID/Touch ID）、ローカル通知、Keychain のセキュアストレージ
- iOS：センサーとハードウェア — ハプティクス、単発の位置情報取得、加速度計、近接センサー、テキスト読み上げ、ストレージ情報、電源／バッテリー状態、ネットワーク状態、キーボードのインセット、画面キャプチャの検出
- iOS：ドキュメント（IOS.md およびドキュメントサイトのガイド）
- Android：ネイティブのメッセージダイアログ（AlertDialog）と、ファイル／複数ファイルを開くダイアログ（Storage Access Framework。キャッシュのコピーとしてインポート）を追加。ディレクトリを開くダイアログと保存ダイアログは明示的なエラーを返します
- Android：ClipboardManager によるクリップボード対応
- Android：WindowMetrics/DisplayMetrics による実際の画面メトリクス（dp、ピクセル、スケール、システムバーを除いた作業領域）
- Android：ハプティクス（`Android.Haptics.Vibrate`）、デバイス情報（`Android.Device.Info`）、トースト（`Android.Toast.Show`）のランタイムメソッド
- Android：型付きライフサイクルイベント（`events.Android.*`、events.txt から生成）を追加し、`ActivityCreated` を `Common.ApplicationStarted` にマッピング
- Android：ビルドパイプラインがインストール可能なデバッグ APK とリリース APK を生成（`android:run`、`android:package`、`android:package:fat`）。リリース署名には、デフォルトではデバッグ用キーストアを使用し、`ANDROID_KEYSTORE_*` 環境変数を介して実際のキーストアも使用できます
- Android：`wails3 doctor` が Android SDK、NDK、JDK を報告
- Android：バッテリー、ネットワーク、テーマ、画面ロック、メモリ不足のシステムイベントを、`events.Android.*` およびプラットフォームに依存しない `events.Common.*` アプリケーションイベントとして公開
- Android：ネイティブモバイル機能ブリッジ（公開された `application.Android*`）— 共有、URL を開く、スリープ抑止、トーチ、セーフエリアのインセット、明るさ、アプリ情報、画面方向のロック、ステータスバー、生体認証（BiometricPrompt）、ローカル通知、EncryptedSharedPreferences のセキュアストレージ
- Android：センサーとハードウェア — ハプティクス、単発の位置情報取得、加速度計、近接センサー、テキスト読み上げ、ストレージ情報、電源／バッテリー状態、ネットワーク状態、キーボードのインセット、FLAG_SECURE による画面キャプチャのブロック
- Android：ドキュメント（ANDROID.md およびドキュメントサイトのガイド）
- 例：`mobile` のキッチンシンクに、iOS と Android のネイティブ機能ブリッジを実演する Mobile タブと Hardware タブを追加（ピル型タブは複数行に折り返されます）
- モバイル：バッテリー — アプリがバックグラウンドになると、加速度計、近接センサー、トーチ、サンプルの定期クロックを一時停止し、復帰時に再開します（Android ではバックグラウンドでもプロセスが動作し続け、iOS ではトーチのハードウェア状態が維持されます）。また、Android のシステムイベントレシーバーは、アプリがフォアグラウンドにある間だけ登録されます
- iOS：カメラ撮影 — `application.IOSCapturePhoto`/`IOSCaptureVideo`（UIImagePickerController → base64 サムネイルを含む `native:capture` イベント）
- iOS：バックグラウンド実行 — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask`（UIApplication のバックグラウンドタスク時間枠）、および生成される Info.plist に `UIBackgroundModes` をテンプレート展開する、設定可能な `ios.backgroundModes`（build/config.yml）
- Android：カメラ撮影 — `application.AndroidCapturePhoto`/`AndroidCaptureVideo`（FileProvider を介したシステムカメラ → `native:capture` イベント）
- Android：フォアグラウンドサービス — `application.AndroidStartForegroundService`/`AndroidStopForegroundService`（継続中の通知を伴う `WailsForegroundService` により、長時間実行するバックグラウンド処理中もプロセスを維持）
- 例：写真／動画の撮影とバックグラウンド実行（Android ではフォアグラウンドサービス、iOS ではバックグラウンドタスク時間枠）を実演する Camera タブ

## 修正

- Linux で `getUserMedia` が常に `NotAllowedError` で失敗する問題を修正。WebKitGTK は処理するハンドラーがない権限要求を拒否しますが、`permission-request` シグナルが接続されていませんでした。カメラ／マイクは、新しいクロスプラットフォームの `WebviewWindowOptions.Permissions` マップ（`map[PermissionType]Permission`）に基づいて処理され、Linux（WebKitGTK）と Windows（WebView2）の両方で適用されるようになりました。ネイティブの確認プロンプトがない Linux では、カメラ／マイクをデフォルトで許可し（`getUserMedia` を復元）、`PermissionDeny` で無効にできます（#5552）
- iOS：`GOOS=ios` が再びコンパイル可能に（公開された `events.IOS`、モバイルのメソッド名スタブ）。また、production タグ付きビルドもコンパイル可能になりました（pkg/application および複数のサービスのビルドタグを修正）
- iOS：Go→JS イベントと ExecJS が動作するようになりました。起動時にページが二重に読み込まれなくなり、`wails:runtime:ready` ハンドシェイクが失われることもなくなりました
- iOS：`ApplicationDidFinishLaunching`/`ApplicationStarted` がアプリの起動処理と競合しなくなりました。固定の 2 秒の起動待機を削除しました
- iOS：Go→JS の JavaScript 実行ごとに発生していた C 文字列のリークを修正
- iOS：`hasListeners` が実際のリスナー登録状態を反映するようになりました
- iOS：本番ビルドではフレームワークのデバッグログをコンパイル対象から除外
- Android：`GOOS=android` が再びコンパイル可能になりました。`events.Android` を定義し、範囲外アクセスを起こす `events_android.go` リスナー配列を削除し、モバイルのメソッド名スタブを追加しました。また、デスクトップ Linux 用ファイル（`linux_cgo.*`、`events_linux.*`、`environment_linux.go`）が Android ビルドに混入しないようにしました
- Android：JS→Go バインディングが動作するようになりました。WebView は `fetch()` の POST 本文を `shouldInterceptRequest` に渡せないため、ランタイム呼び出しは JavascriptInterface トランスポート（`nativeHandleRuntimeCall`）経由になり、nil のリクエスト本文によるクラッシュを回避します
- Android：`Screens.*` ランタイム呼び出しが実際のデータを返すようになりました。起動時に ScreenManager が初期化されるようになりました（従来は接続されておらず、`GetAll` が nil を返していました）
- Android：本番ビルドではフレームワークのデバッグログをコンパイル対象から除外し、デバッグビルドでは `Wails` タグを使用して logcat に出力するようにしました
- Android：実際に機能する `hasListeners` レジストリ、JNI の参照および例外処理、単一ロードのページライフサイクル（二重ナビゲーションなし）を実装しました
- Windows で Vite 開発サーバーの実行中に `wails3 generate bindings` が「Access is denied」で失敗する問題を修正しました。生成ファイルを出力ディレクトリ上で置き換えるのではなく、出力ディレクトリへ同期するようにしました（#5515）
- ディスプレイの変更後に画面情報を読み取る際、macOS で断続的に致命的なクラッシュが発生する問題を修正しました。画面 ID と名前が、自動解放される `UTF8String` バッファへのポインターを保持しており、Go がコピーする前に解放される可能性がありました（解放後使用）。文字列を `strdup` し、変換後に解放するようにしました。また、画面の列挙を明示的な autorelease pool 内で実行することで、Go の goroutine から呼び出した場合のリークも解消しました（#5556）
- assetserver が `WebKitURISchemeRequest` を閉じる際、Linux で断続的に SIGSEGV が発生する問題を修正しました。最後の `g_object_unref` が assetserver の goroutine 上で実行され、GTK メインスレッド外で WebKit GObject をファイナライズしていました。unref は `g_main_context_invoke` を介して GTK のメインコンテキスト上にマーシャリングされるようになりました（#5557）

## v3.0.0-alpha.100 - 2026-06-13

## 追加

- `MacWebviewPreferences` に WKWebView の追加設定オプション `EnableAutoplayWithoutUserAction`、`AllowsAirPlayForMediaPlayback`、`AllowsMagnification`、`JavaScriptCanOpenWindowsAutomatically`、`MinimumFontSize`、`ApplicationNameForUserAgent` を追加しました（#5549）

## 修正

- Windows で Vite 開発サーバーの実行中に `wails3 generate bindings` が「Access is denied」で失敗する問題を修正しました。生成ファイルを出力ディレクトリ上で置き換えるのではなく、出力ディレクトリへ同期するようにしました（#5561）
- Linux のフレームレスウィンドウで JS のリサイズイベントが発生しない問題と、フレームレスウィンドウにおけるスクロールバーの端の検出を修正しました（#5368）
- Windows で一時ディレクトリとインストールディレクトリが異なるボリュームにある場合、アップデーターが「invalid cross-device link」で失敗する問題を修正しました（#5560）

## v3.0.0-alpha.99 - 2026-06-10

## 修正

- Windows で Vite 開発サーバーの実行中に `wails3 generate bindings` が「Access is denied」で失敗する問題を修正しました。生成ファイルを出力ディレクトリ上で置き換えるのではなく、出力ディレクトリへ同期するようにしました（#5515）

## v3.0.0-alpha.98 - 2026-06-03

## 修正

- アイドル時（インスペクターを開いている場合など）に Linux 上で WebKit UI がフリーズする問題を修正しました。JavaScriptCore の GC スレッド同期を壊していた、`SIGUSR1` への `SA_ONSTACK` の強制を廃止しました（#5527）

## v3.0.0-alpha.97 - 2026-05-31

## 追加

- デバッグページと `runtime/trace` の使用方法を追加しました

## 変更

- 不要な `_ "embed"` import をいくつか削除し、コードを整理しました

## 修正

- Windows でウィンドウの最大化を解除した後、最小幅および最小高さの制約が適用されない問題を修正しました（#4593）
- Frameless と Transparent のウィンドウオプションを使用したフルスクリーンモードで、マウスクリックがウィンドウを透過する問題を修正しました（#4408）

## v3.0.0-alpha.96 - 2026-05-25

## 追加

- Garble 難読化のサポート（[#4563](https://github.com/wailsapp/wails/issues/4563)）を追加しました。安定したバインディングメソッド ID、ビルドおよび Taskfile の仕組み（`build --obfuscated --garbleargs`、`generate bindings -obfuscated`）、さらに Garble によってエクスポート済みフィールドの名前が変更されても通信形式を維持できるよう、ランタイム向けのすべてのペイロードに JSON 構造体タグ（`EnvironmentInfo`、`OSInfo`、`Screen`、`Rect`、`Point`、`Size`、`Capabilities`）を追加しました。

## v3.0.0-alpha.95 - 2026-05-20

## 追加

- 不足していたプロジェクト構造のページを追加しました

## 変更

- ドキュメント：アーキテクチャページの図をより見やすくするため、一部をシーケンス図に変更しました
- ドキュメント：実行の前提条件として D2 のインストールが必要である旨の注記を追加しました

## 修正

- GTK4 がデフォルトの場合の `wails3 generate appimage` を修正しました。バンドラーがランタイムファイルを検索する前にバイナリから GTK スタックを検出するようになり、GTK4 ビルドでは `webkitgtk-6.0/` 配下の `libwebkitgtkinjectedbundle.so` を、`-tags gtk3` ビルドでは `webkit2gtk-4.1/` 配下の `libwebkit2gtkinjectedbundle.so` を選択するようになりました。`.relr.dyn` のプローブでも `libgtk-4.so.1` を確認するため、スタックにかかわらず、最新のツールチェーンではストリッピングが正しく無効化されます。（#5475）
- 相対 `-builddir` を指定して呼び出した場合に `wails3 generate appimage` が失敗する問題を修正しました。バンドラーが `-binary`、`-icon`、`-desktopfile`、`-builddir`、`-outputdir` を最初に絶対パスへ解決するようになったため、処理途中の `s.CD` によって AppRun のダウンロード用 goroutine やコピー後の `ldd` プローブが壊れることはありません。
- desktop の `Name=` フィールドがバイナリのベース名と一致しない場合に、`wails3 generate appimage` が最終的な AppImage を `-outputdir` へ移動できない問題を修正しました。バンドラーが `OUTPUT` 環境変数を介して linuxdeploy の appimage プラグインに、desktop ファイルから導出した名前ではなく `<binary>-<arch>.AppImage` へ AppImage を書き込ませるようになりました。
- alpha.93 で GTK4 + WebKitGTK 6.0 スタックがデフォルトに昇格した後、Linux で `events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`、`Common.SystemDidWake` が発生しない問題を修正しました。新しいデフォルトの `application_linux.go` `run()` が、`Linux.*` イベントを対応する `Common.*` へ転送する `setupCommonEvents()` と、`monitorPowerEvents()` を呼び出していませんでした。DBus の電源監視ヘルパーは、`application_linux_dbus.go` を介して GTK3 と GTK4 のビルドパス間で共有されるようになりました。（#5474）

## v3.0.0-alpha.94 - 2026-05-19

## 修正

- alpha.93 で GTK4 + WebKitGTK 6.0 スタックがデフォルトに昇格した後、Linux で `events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`、`Common.SystemDidWake` が発生しない問題を修正しました。新しいデフォルトの `application_linux.go` `run()` が、`Linux.*` イベントを対応する `Common.*` へ転送する `setupCommonEvents()` と、`monitorPowerEvents()` を呼び出していませんでした。DBus の電源監視ヘルパーは、`application_linux_dbus.go` を介して GTK3 と GTK4 のビルドパス間で共有されるようになりました。（#5474）

## v3.0.0-alpha.93 - 2026-05-17

## 追加

- @leaanthony により、Linux の `wails3 doctor` 出力に `XDG_SESSION_TYPE` を追加しました

## 修正

- appmenu-gtk-module が未実体化のウィンドウへアクセスすることで発生していた、Wayland 上のウィンドウメニューのクラッシュを @leaanthony が修正しました（#4769）
- アプリ名に無効な文字（スペース、丸括弧など）が含まれる場合に GTK アプリケーションがクラッシュする問題を @leaanthony が修正しました
- Windows でドラッグ＆ドロップを初期化する際に発生する「not enough memory」エラーを @overlordtm が修正しました（#4701）
- mainthread のコールバックストアで、map からの削除に誤って RLock を使用していたことによる競合状態を @leaanthony が修正しました（Linux、macOS、iOS）（#4424）
- タスクへコマンドライン引数を渡す際の変数処理を修正しました。KEY=VALUE のペアとして指定された CLI 変数が正しく初期化され、タスクの実行全体へ伝播されるようになりました。
- macOS における NSWindowZoomButton の競合を修正しました。`MaximiseButtonState` と `FullscreenButtonState` は、起動時と実行時の両方で、より制限の厳しい状態を適用するようになりました。どちらのセッターも、もう一方の設定を暗黙に上書きできなくなりました（#5319）
- CodeRabbit が #5463 で検出した、従来の GTK3 ビルドパス（`-tags gtk3`）に以前から存在していた一連のバグを修正しました。ファイル関連付けによる起動でスタートアップハンドラーがスキップされなくなりました。`getTheme` は境界および型の安全性を確保しました。`appName` は GLib が所有するメモリを解放しなくなりました。`clipboardGet` は GTK が返す `gchar*` をリークしなくなりました。`Calloc` がポインタレシーバーを使用するようになり（さらに `NewCalloc` が `*Calloc` を返すようになり）、プールが実際に割り当てを追跡するようになりました。`zoomOut` は、1.0 にクランプされる負の乗数ではなく、`zoomInFactor` の逆数を使用するようになりました。`execJS` は、呼び出しごとに `C.CString("")` をリークせず、事前割り当て済みの空のワールド名を再利用するようになりました。開発用の `fmt.Println` を `menuItem.setAccelerator` から削除しました。#5465 を解決しました。
- デフォルトの GTK4 ビルドパス（`linux_cgo.go`）でも、同じ `Calloc` の値レシーバーによるリークを修正しました。ポインタレシーバーと `NewCalloc() *Calloc` により、ウィンドウごとの `c.String(...)` の割り当てが実際に追跡され、解放されるようになりました。

## v3.0.0-alpha.92 - 2026-05-15

## 追加

- `PACKAGE_MANAGER` オプションで使用するフロントエンドパッケージマネージャーを制御できるように Taskfile を変更しました
- Taskfile テンプレートをより予測しやすく記述できるよう、テンプレートデータに `{{.Opn}}` と `{{.Cls}}` を追加しました

## 変更

- 既存の Taskfile の一部を、`{{.Opn}} and {{.Cls}}` を使用するように変更しました

## 修正

- パネルがトレイメニューを読み取っている間にメニューを更新すると、`linuxSystemTray` で発生する `concurrent map read and map write` ランタイムの致命的エラーを修正しました。
- Windows でコンソールを接続せずにアプリを実行した場合でもメッセージが失われないよう、WebView2 のエラーとスタックトレースの出力に `fmt` ではなく `log` を使用するようにしました。

## v3.0.0-alpha.91 - 2026-05-12

## 変更

- `@github-actions[bot]` による [PR](https://github.com/wailsapp/wails/pull/5414) でスポンサー SVG を更新しました
- <strong>破壊的変更（macOS）：</strong>macOS の座標系を正規化し、`GetScreens`、`Position`、`SetPosition` がすべて同じ座標空間を使用するようにしました。これは論理ポイント単位の Y 軸下向きの空間で、プライマリ画面の左上が `(0,0)` です。Windows、GTK、および Electron と Web の公開 API と同じ仕様です。プライマリ画面より物理的に上にある画面では、`Bounds.Y` が負の値として報告されるようになりました（従来は正の値）。また、`Position()`／`SetPosition()` の値は、`points × primaryScale` ではなく論理ポイント単位になりました。`Position()` → `SetPosition()` の往復変換は維持されます。以前のアルファビルドで記録した絶対値や、手動で計算した回避策（`primaryScale` を乗算する、画面の高さを基準に Y を反転するなど）は更新する必要があります。[#5117](https://github.com/wailsapp/wails/issues/5117) を解決しました。

## 修正

- パニックを防ぐため、DBus シグナル名と本体の長さを防御的に検証するようにしました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5416)）
- Linux の GTK メニュー処理におけるメモリ安全性の問題を修正しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5363)）
- NVIDIA GPU を検出し、Linux で DMA-BUF レンダラーを無効にするようにしました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5295)）
- macOS の `SetPosition` における画面間の Y 座標変換を修正しました。プライマリディスプレイから垂直方向にずれたモニターでもウィンドウが正しい位置に配置されるよう、プライマリ画面の高さをグローバルな基準として使用します（[#5117](https://github.com/wailsapp/wails/issues/5117)）
- Git の PR テンプレートが正しいフィードバック URL を参照するように修正しました（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5109)）
- 引数を1つではなく4つ渡していた壊れた `DestroyMenu` システムコールによって、すべての呼び出しが FALSE を返し、何も解放されなかったために発生していた、一連の Windows システムトレイの `SetMenu` クラッシュを修正しました。また、メニューの再構築時に HMENU と HBITMAP のハンドル（`MenuItem.SetBitmap` によって実行時に割り当てられたものを含む）を解放し、`Win32Menu.Update` に残っていた古いチェックボックス／ラジオボタンのマップをリセットするようにしました。さらに、割り当てを二重にしていた `systemtray.updateMenu` 内の冗長な `Update()` 呼び出しを削除しました。長時間稼働するシステムトレイアプリで、メニューを再構築するたびに GDI/USER オブジェクトがリークしなくなりました。

## v3.0.0-alpha.90 - 2026-05-11

## 追加

- macOS の WKWebView User-Agent に設定可能なアプリケーション名を追加しました（@vinhvoit225 による [PR](https://github.com/wailsapp/wails/pull/5261)）
- gin-service サンプルに間接依存関係 github.com/coder/websocket を追加しました（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5400)）
- ビルドアセットのテストに深い等価比較のサポートを追加しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5402)）

## 変更

- ビルド出力を assets ディレクトリに統合しました（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5401)）
- `@github-actions[bot]` による [PR](https://github.com/wailsapp/wails/pull/5399) でスポンサー SVG を更新しました

## 修正

- macOS の単一インスタンスメッセージに通知オブジェクトを使用するようにしました（@overlordtm による [PR](https://github.com/wailsapp/wails/pull/5289)）
- 高負荷時の Promise の消失を防ぐため、Windows のコールバックをバッチ処理するようにしました（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5383)）

## v3.0.0-alpha.89 - 2026-05-10

## 追加

- Go のテスト結果を集約する go<em>test</em>results ジョブを追加しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5316)）

## 変更

- 大きな RPC ペイロードを条件に応じて分割し、チャンク化した POST リクエストとして送信するようにしました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5369)）
- すべてのフロントエンドテンプレートで Vite を 5.x.x から 8.0.0 にアップグレードしました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5386)）
- Vite 開発サーバーのポート設定を環境変数へ移行しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5365)）
- すべてのテンプレートで Vite 開発サーバーが 127.0.0.1 にバインドするように設定しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5361)）
- `@github-actions[bot]` による [PR](https://github.com/wailsapp/wails/pull/5384) でスポンサー SVG を更新しました

## 修正

- ビルドアセットの更新時に Info.plist テンプレートのスタブをサニタイズするようにしました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5312)）
- メニュー項目の変更メソッド（`setMenuItemChecked()`、`setMenuItemLabel()`、`setMenuItemDisabled()`、`setMenuItemHidden()`、`setMenuItemTooltip()`）をメインスレッドで同期的に適用し、macOS メニューの古い状態が残る問題を修正しました。これにより、メニューを短時間で再度開いたときに以前の状態が描画されていた `dispatch_async` の競合状態を解消しました（#5002）
- 不要な再ビルドを防ぐため、開発モードで `*_test.go` ファイルを無視するようにしました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5203)）
- アプリが実行されていないときに Menu.Update() でセグメンテーション違反が発生しないようにしました（@wucm667 による [PR](https://github.com/wailsapp/wails/pull/5291)）
- Windows でメニューバーの再描画を制御するために lastSizeWParam を使用するようにしました（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5382)）

## v3.0.0-alpha.88 - 2026-05-09

## 変更

- HiddenOnTaskbar が WS<em>EX</em>TOOLWINDOW を使用するように変更しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5371)）
- 依存関係を並べ替え、go.mod から webview2 の replace ディレクティブを削除しました（@atterpac による [PR](https://github.com/wailsapp/wails/pull/5370)）
- `@github-actions[bot]` による [PR](https://github.com/wailsapp/wails/pull/5358) でスポンサー SVG を更新しました

## 修正

- 汎用的な間接参照エイリアスを削除し、マップのキー型を統合しました（@fbbdev による [PR](https://github.com/wailsapp/wails/pull/5331)）

## 削除

- PR-master ワークフローを削除し、@leaanthony が [PR](https://github.com/wailsapp/wails/pull/5377) でドキュメント、Go テスト、スキップテストを削除

## v3.0.0-alpha.87 - 2026-05-07

## 追加

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5352) で Wails v3 の韓国語ドキュメントを追加
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5354) でインストールとクイックスタートのフランス語ドキュメントを追加
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5355) でクイックスタート、概念、コミュニティのポルトガル語ドキュメントを追加

## v3.0.0-alpha.86 - 2026-05-06

## 追加

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5328) でフランス語のドキュメントローカライズを追加
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5343) でドキュメントサイトにドイツ語ロケールを追加

## 変更

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5347) で翻訳済みの全 8 ロケールをドキュメント設定に登録
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5317) で WebView2 に関連する各種 Windows ファイルを更新

## 修正

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5340) で Linux のダイアログディスパッチを GTK3 と GTK4 に分割
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5339) でダイアログのコールバックが GTK スレッド上で実行されるようにし、セグメンテーションフォルトを修正

## v3.0.0-alpha.85 - 2026-05-05

## 追加

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5179) で PR テンプレートの URL をリポジトリに追加
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5330) で Wails v3 のドイツ語ドキュメントを追加

## v3.0.0-alpha.84 - 2026-05-03

## 追加

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5307) で macOS において Escape キーによるフルスクリーン終了を無効化するオプションを追加
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5310) で macOS において Escape キーによるフルスクリーン終了を無効化するオプションを追加
- @yuseferi が [PR](https://github.com/wailsapp/wails/pull/5288) で Pausa のコミュニティショーケースに関するドキュメントを追加

## 変更

- `@github-actions[bot]` が [PR](https://github.com/wailsapp/wails/pull/5308) でスポンサーの SVG を更新
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5309) で未対応プラットフォームを処理できるようアイコン生成コマンドを更新
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5224) で真偽値のフルスクリーン API を3状態の ButtonState に置き換え、プラットフォームバインディングを実装

## 修正

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5315) でコントローラーの状態が nil の場合に WebView2 のフォーカス操作を実行しないよう保護
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5313) で PR のベースブランチを正しく参照するよう GitHub Actions ワークフローを更新
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5203) で不要な再ビルドを防ぐため、開発モードでは `*_test.go` ファイルを無視するように変更
- @wucm667 が [PR](https://github.com/wailsapp/wails/pull/5291) でアプリが実行されていない場合に Menu.Update() で発生するセグメンテーションフォルトを防止

## v3.0.0-alpha.83 - 2026-05-02

## 追加

- @symball が [PR](https://github.com/wailsapp/wails/pull/5094) でマシン単位／ユーザー単位のインストール向けに InstallScope フラグとビルドオプションを追加
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5294) で Window インターフェースを満たすため、何も行わない SetScreen メソッドを BrowserWindow に追加

## 修正

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5295) で NVIDIA GPU を検出し、Linux 上で DMA-BUF レンダラーを無効化
- @wayneforrest が [PR](https://github.com/wailsapp/wails/pull/5109) で正しいフィードバック URL を指すよう git PR テンプレートを修正
- 1個ではなく4個の引数を渡していた壊れた `DestroyMenu` システムコールにより、すべての呼び出しが FALSE を返し、何も解放されなかったことが原因の Windows システムトレイ `SetMenu` の一連のクラッシュを修正。さらに、メニューの再構築時に HMENU および HBITMAP ハンドル（`MenuItem.SetBitmap` を介して実行時に割り当てられたものを含む）を解放し、`Win32Menu.Update` 内の古いチェックボックス／ラジオマップをリセットし、割り当てを倍増させていた `systemtray.updateMenu` 内の冗長な `Update()` 呼び出しを削除。長時間稼働するシステムトレイアプリで、メニューを再構築するたびに GDI/USER オブジェクトがリークしなくなりました。

## v3.0.0-alpha.82 - 2026-05-01

## 修正

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5232) でデスクトップ名を正しく処理するようデスクトップファイル生成を修正

## v3.0.0-alpha.81 - 2026-04-30

## 変更

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5286) でナイトリーリリースのスケジュールを UTC 15:00 に調整

## 修正

- Retina Mac で Screen の Bounds、WorkArea、Size の値が半分になる問題を修正 -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## 変更

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5285) でドキュメントの依存関係とコンテンツコレクションローダーを更新

## v3.0.0-alpha.79 - 2026-04-29

## 追加

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5270) で trigger-release ジョブに actions: write 権限を付与

## 変更

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5283) でリリースタスクのデフォルトを master ブランチに設定し、変更履歴の文言を更新
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5282) で最新バージョンを使用するよう auto-changelog ワークフローを更新
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5280) でパスフィルターを追加し、使用されていないワークフローを削除して、ワークフローの効率を改善
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5274) でサンプルへのリンクが master ブランチを参照するようドキュメントを更新
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5272) で v3 のドキュメントとサンプルを更新

## 修正

- @AkagiYui が [PR](https://github.com/wailsapp/wails/pull/5265) で開発時のリトライロジックと IPv4 の強制使用をリバースプロキシに追加
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5281) で未リリース変更履歴のトリガーワークフローを書き直し

## 削除

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5267) でさまざまなテスト用途のシェルテストスクリプトを削除
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5266) で v3-alpha ドキュメントのデプロイワークフローと CNAME レコードを削除

### 追加

- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5196) でサイドバーナビゲーションに「フロントエンドルーティング」の項目を追加
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/5185) でフレームワーク別の推奨事項を含むフロントエンドルーティングガイドを追加
- モーダルシート（macOS）のサポートを追加
- Apple デバイスのサポートを改善するため、@leaanthony が ghw のバージョンを更新（#4977）
- Dock サービスに `GetBadge` メソッドを追加
- カスタム Go ビルドタグ（例：`wails3 build -tags gtk4`）を渡すため、`wails3 build` コマンドに `-tags` フラグを追加（#4957）
- バインディングジェネレーターによる enum の自動生成について、専用の Enums ページとサイドバーナビゲーションを含むドキュメントを追加（#4972）
- カスタム Go ビルドタグ（例：`wails3 build -tags gtk4`）を渡すため、`wails3 build` コマンドに `-tags` フラグを追加（#4957）
- Storage（localStorage、sessionStorage、IndexedDB、Cache API）、Network（Fetch、WebSocket、XMLHttpRequest、EventSource、Beacon）、Media（Canvas、WebGL、Web Audio、MediaDevices、MediaRecorder、Speech Synthesis）、Device（Geolocation、Clipboard、Fullscreen、Device Orientation、Vibration、Gamepad）、Performance（Performance API、Mutation Observer、Intersection/Resize Observer）、UI（Web Components、Pointer Events、Selection、Dialog、Drag and Drop）など、41 個以上のブラウザー API を紹介する Web API のサンプルを `v3/examples/web-apis/` に追加
- プラットフォーム間で 200 個以上のブラウザー API をテストする WebView API 互換性チェッカーのサンプル（`v3/examples/webview-api-check/`）を追加
- 並列検索とキャッシュに対応し、Flatpak/Snap/Nix をサポートする、Linux 上でネイティブライブラリのパスを検索するための `internal/libpath` パッケージを追加
- **作業中：** Linux 向けに実験的な WebKitGTK 6.0 / GTK4 サポートを追加。`-tags gtk4` で利用可能（GTK3/WebKit2GTK 4.1 が引き続きデフォルト）
- 注：タイル型ウィンドウマネージャー（Hyprland、Sway など）では、ウィンドウのジオメトリを WM が制御するため、最小化／最大化操作が期待どおりに動作しない場合があります
- @AbdelhadiSeddar が、ドキュメントの「**JavaScript でイベントをリッスンする**」に「**一度限りのハンドラー**」の実装方法を追加
- @leaanthony が、Windows/Linux 上のウィンドウで `app.Menu.Set()` により設定されたアプリケーションメニューを継承できるよう、`WebviewWindowOptions` に `UseApplicationMenu` オプションを追加
- @wimaha が、Liquid Glass アイコンとアセットカタログの生成に `.icon` ファイル（Apple Icon Composer 形式）を使用するためのサポートを追加（macOS）（#4934）
- ヘッドレス／Web デプロイ向けの実験的なサーバーモード（`-tags server`）を追加。ネイティブ GUI 依存関係なしで Wails アプリを HTTP サーバーとして実行できます。`wails3 task build:server` でビルドしてください。詳細は `examples/server` を参照してください。
- 並列検索とキャッシュに対応し、Flatpak/Snap/Nix をサポートする、Linux 上でネイティブライブラリのパスを検索するための `internal/libpath` パッケージを追加
- @leaanthony が、macOS の Spaces とフルスクリーンをまたぐウィンドウの動作を制御するため、`MacWindow` に `CollectionBehavior` オプションを追加（#4756）
- @leaanthony が pkg/application のユニットテストを追加
- @leaanthony が MSIX パッケージングにカスタムプロトコルのサポートを追加
- Linux にデスクトップ環境の検出機能を追加（[PR #4797](https://github.com/wailsapp/wails/pull/4797)）
- @leaanthony が、フロントエンドから印刷ダイアログを表示するための `Window.Print()` メソッドを JavaScript ランタイムに追加（#4290）
- @leaanthony が Linux 上の `wails3 doctor` 出力に `XDG_SESSION_TYPE` を追加
- @leaanthony が Linux 向けに WebKit2 の追加のロード変更イベント `WindowLoadStarted`、`WindowLoadRedirected`、`WindowLoadCommitted`、`WindowLoadFinished` を追加（#3896）
- @leaanthony が Linux 上の `wails3 doctor` 出力に `XDG_SESSION_TYPE` を追加
- パッケージング時だけでなく、Linux のビルド時にも `.desktop` ファイルを生成（#4575）
- @leaanthony が、ディストリビューション固有のパッケージ名と nfpm パッケージングの例を含む Linux ランタイム依存関係のドキュメントを追加（#4339）
- @leaanthony が Linux 上の `wails3 doctor` 出力に NVIDIA ドライバーのバージョン情報を追加
- @APshenkin が [PR](https://github.com/wailsapp/wails/pull/4710) で生メッセージハンドラーにオリジンを追加
- @APshenkin が [PR](https://github.com/wailsapp/wails/pull/4712) で macOS にユニバーサルリンクのサポートを追加
- @APshenkin が [PR](https://github.com/wailsapp/wails/pull/4702) でバインディングのトランスポート層をリファクタリング
- @chinenual が [PR](https://github.com/wailsapp/wails/pull/4760) で helloworld テンプレートに aria-label 識別子を追加し、Appium テストクライアントからサンプルアプリを容易にテストできるように変更
- @APshenkin が [PR](https://github.com/wailsapp/wails/pull/4710) で生メッセージハンドラーにオリジンを追加
- @APshenkin が [PR](https://github.com/wailsapp/wails/pull/4712) で macOS にユニバーサルリンクのサポートを追加
- @APshenkin が [PR](https://github.com/wailsapp/wails/pull/4702) でバインディングのトランスポート層をリファクタリング
- @fbbdev と @ianvs による型付きイベント（[#4633](https://github.com/wailsapp/wails/pull/4633)）
- ツールチップをリアルタイムで更新するヘッドレストレイを示す `systray-clock` サンプルを追加（#4653）。
- @Tolfx が #4510 で Windows 向け NSIS プロトコルテンプレートを追加
- @Tolfx が #4510 で build-assets のテストを追加
- macOS：@nidib により、メニューバーにネイティブのウィンドウコントロールを表示（[#4588](https://github.com/wailsapp/wails/pull/4588)）
- @popaprozac が [PR](https://github.com/wailsapp/wails/pull/4451) で、Dock にあるアプリアイコンを非表示／表示にする macOS Dock サービスを追加
- @popaprozac が [PR](https://github.com/wailsapp/wails/pull/4451) で、Dock にあるアプリアイコンを非表示／表示にする macOS Dock サービスを追加
- @leaanthony が [#4534](https://github.com/wailsapp/wails/pull/4534) で、NSGlassEffectView（macOS 15.0 以降）と NSVisualEffectView フォールバックを使用する macOS 向けネイティブ Liquid Glass エフェクトのサポートを追加。包括的なマテリアルカスタマイズオプションを含む
- @leaanthony によるブラウザー URL のサニタイズ（[#4500](https://github.dev/wailsapp/wails/pull/4500)）。@APShenkin による [#4484](https://github.com/wailsapp/wails/pull/4484) を基に実装。
- [@Taiterbase](https://github.com/Taiterbase) がこの [PR](https://github.com/wailsapp/wails/pull/4241) で行った当初の作業を基に、[@leaanthony](https://github.com/leaanthony) が Windows/Mac にコンテンツ保護を追加
- @leaanthony が [PR](https://github.com/wailsapp/wails/pull/4488) で、`wails3 build` および `wails3 package` エイリアスを介して CLI 変数を Task コマンドに渡すためのサポートを追加（#4422）
- [@atterpac](https://github.com/atterpac) により、イベントのソース要素データを使用するドロップゾーンをサポート（[#4318](https://github.com/wailsapp/wails/pull/4318)）
- WebView2 ブラウザーに追加のコマンドライン引数を渡せるよう、`WindowsWindow` オプションに `AdditionalLaunchArgs` を追加（[PR](https://github.com/wailsapp/wails/pull/4467)）
- [@triadmoko](https://github.com/triadmoko) により、wails init の後に go mod tidy を自動実行する機能を追加（[PR](https://github.com/wailsapp/wails/pull/4286)）
- @leaanthony による Windows Snap Assist 機能（[PR](https://github.dev/wailsapp/wails/pull/4463)）
- WebView2 ブラウザーに追加のコマンドライン引数を渡せるよう、`WindowsWindow` オプションに `AdditionalLaunchArgs` を追加（[PR](https://github.com/wailsapp/wails/pull/4467)）
- [@triadmoko](https://github.com/triadmoko) により、wails init の後に go mod tidy を自動実行する機能を追加（[PR](https://github.com/wailsapp/wails/pull/4286)）
- @leaanthony による Windows Snap Assist 機能（[PR](https://github.dev/wailsapp/wails/pull/4463)）
- [@almas-x](https://github.com/almas-x) が [PR](https://github.com/wailsapp/wails/pull/4427) で Windows 向け `getAccentColor` 実装を追加
- [@almas-x](https://github.com/almas-x) が [PR](https://github.com/wailsapp/wails/pull/4427) で Windows 向け `getAccentColor` 実装を追加
- Windows のダークテーマ対応メニューとメニューバー。@leaanthony による変更（[a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)）
- JS/TS バインディングを明確にするため、組み込みサービスの名前を変更。@popaprozac による変更（[PR](https://github.com/wailsapp/wails/pull/4405)）
- ユーザーのシステムのアクセントカラーを取得する `app.Env.GetAccentColor`。MacOS で動作。[@etesam913](https://github.com/etesam913) による変更
- `window.ToggleFrameless()` API を追加。[@atterpac](https://github.com/atterpac) による変更（[#4137](https://github.com/wailsapp/wails/pull/4137)）
- Linux 向けにディストリビューション固有のビルド依存関係を追加。@leaanthony による変更（[PR](https://github.com/wailsapp/wails/pull/4345)）
- バインディングガイドを追加。@atterpac による変更（[PR](https://github.com/wailsapp/wails/pull/4404)）
- **テスト基盤の整理**：Docker テストファイルを専用の `test/docker/` ディレクトリへ移動し、イメージを最適化してビルドの信頼性を向上。[@leaanthony](https://github.com/leaanthony) による変更（[#4359](https://github.com/wailsapp/wails/pull/4359)）
- **リソース管理パターンの改善**：サンプルに適切なイベントハンドラーのクリーンアップと、コンテキストを考慮した goroutine 管理を追加。[@leaanthony](https://github.com/leaanthony) による変更（[#4359](https://github.com/wailsapp/wails/pull/4359)）
- aarch64 AppImage のビルドをサポート。[@AkshayKalose](https://github.com/AkshayKalose) による変更（[#3981](https://github.com/wailsapp/wails/pull/3981)）
- `wails doctor` に診断セクションを追加。[@leaanthony](https://github.com/leaanthony) による変更
- サービスメソッドの呼び出し時に、ウィンドウをコンテキストへ追加。[@leaanthony](https://github.com/leaanthony) による変更
- どのウィンドウがサービスを呼び出しているかを識別する方法を示す `window-call` サンプルを追加。[@leaanthony](https://github.com/leaanthony) による変更
- 新しいメニューガイドを追加。[@leaanthony](https://github.com/leaanthony) による変更
- panic 処理を改善。[@leaanthony](https://github.com/leaanthony) による変更
- 新しいメニューガイドを追加。[@leaanthony](https://github.com/leaanthony) による変更
- Service API にドキュメントコメントを追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4024](https://github.com/wailsapp/wails/pull/4024)）
- 追加設定を使用してサービスを初期化する関数 `application.NewServiceWithOptions` を追加。[@leaanthony](https://github.com/leaanthony) による変更（[#4024](https://github.com/wailsapp/wails/pull/4024)）
- メニュー制御を改善。[@FalcoG](https://github.com/FalcoG) と [@leaanthony](https://github.com/leaanthony) による変更（[#4031](https://github.com/wailsapp/wails/pull/4031)）
- ドキュメントをさらに追加。[@leaanthony](https://github.com/leaanthony) による変更
- 標準イベントリスナーでイベントのキャンセルをサポート。[@leaanthony](https://github.com/leaanthony) による変更
- システムトレイで `Hide`、`Show`、`Destroy` をサポート。[@leaanthony](https://github.com/leaanthony) による変更
- システムトレイで `SetTooltip` をサポート。[@leaanthony](https://github.com/leaanthony) による変更。原案は [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- サポートされていない型に関するバインディングジェネレーターの警告で、パッケージパスを報告。[@fbbdev](https://github.com/fbbdev) による変更（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- バインディングジェネレーターでジェネリックエイリアスをサポート。[@fbbdev](https://github.com/fbbdev) による変更（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- バインディングジェネレーターで `omitzero` JSON フラグをサポート。[@fbbdev](https://github.com/fbbdev) による変更（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- 指定したサービスメソッドのバインディング生成を防止する `//wails:ignore` ディレクティブを追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- Go ではエクスポートするものの JS/TS ではエクスポートしない型を使用できるように、サービスとモデルへ `//wails:internal` ディレクティブを追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- 弱い型付けの列挙型を使用できるように、バインディングジェネレーターでエイリアス型の定数をサポート。[@fbbdev](https://github.com/fbbdev) による変更（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- Go の 1.24 機能に対するバインディングジェネレーターのテストを追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4068](https://github.com/wailsapp/wails/pull/4068)）
- OS バージョンの検出を改善するため、`OSInfo.Branding` に macOS 15「Sequoia」のサポートを追加（[#4065](https://github.com/wailsapp/wails/pull/4065)）
- シャットダウン処理の完了後にカスタムコードを実行する `PostShutdown` フックを追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- カスタムエラーハンドラーで致命的エラーを検出できるように、`FatalError` 構造体を追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- サービスの起動順序とシャットダウン順序を標準化し、文書化。[@fbbdev](https://github.com/fbbdev) による変更（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- アプリケーションの起動／シャットダウンシーケンス用テストハーネスと、サービスの起動／シャットダウンテストを追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- アプリケーションの作成後にサービスを登録する `RegisterService` メソッドを追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- バインディング呼び出しのエラーをカスタム処理するため、アプリケーションとサービスのオプションに `MarshalError` フィールドを追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- Promise チェーンを通じてキャンセル要求を伝播する、キャンセル可能な Promise ラッパーを追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4100](https://github.com/wailsapp/wails/pull/4100)）
- バインディング呼び出しのキャンセルを `AbortSignal` に関連付ける機能を追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4100](https://github.com/wailsapp/wails/pull/4100)）
- WML で、通常の `wml-*` 属性に加えて `data-wml-*` 属性をサポート。[@leaanthony](https://github.com/leaanthony) による変更
- 遅延設定と動的な再設定のため、すべてのサービスに `Configure` メソッドを追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- `fileserver` サービスは、未設定の場合に 503 Service Unavailable レスポンスを送信。[@fbbdev](https://github.com/fbbdev) による変更（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- `kvstore` サービスは、未設定の場合にデフォルトでインメモリのキーバリューストアを提供。[@fbbdev](https://github.com/fbbdev) による変更（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- 設定変更後にファイルからデータを再読み込みする `Load` メソッドを `kvstore` サービスに追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- すべてのキーを削除する `Clear` メソッドを `kvstore` サービスに追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- JS 側にログレベル定数を提供する型 `Level` を `log` サービスに追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- ログレベルを動的に指定する `Log` メソッドを `log` サービスに追加。[@fbbdev](https://github.com/fbbdev) による変更（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- `sqlite` サービスは、未設定の場合、デフォルトでインメモリ DB を提供するようになりました。[#4067](https://github.com/wailsapp/wails/pull/4067) で [@fbbdev](https://github.com/fbbdev) が実装しました
- DB を手動で閉じるための `Close` メソッドを `sqlite` サービスに追加しました。[#4067](https://github.com/wailsapp/wails/pull/4067) で [@fbbdev](https://github.com/fbbdev) が実装しました
- `sqlite` サービスのクエリメソッドにキャンセル対応を追加しました。[#4067](https://github.com/wailsapp/wails/pull/4067) で [@fbbdev](https://github.com/fbbdev) が実装しました
- `sqlite` サービスに、JS バインディングを備えたプリペアドステートメント対応を追加しました。[#4067](https://github.com/wailsapp/wails/pull/4067) で [@fbbdev](https://github.com/fbbdev) が実装しました
- [@AnalogJ](https://github.com/AnalogJ) によるこちらの [PR](https://github.com/wailsapp/wails/pull/3537) の元の実装を基に、[Lea Anthony](https://github.com/leaanthony) が [PR](https://github.com/wailsapp/wails/pull/3537) で Gin 対応を追加しました
- 自動保存とパスワードの自動保存が常に有効になる問題を修正しました。[#4134](https://github.com/wailsapp/wails/pull/4134) で [@oSethoum](https://github.com/osethoum) が実装しました
- ウィンドウにメニューを設定できるよう、ウィンドウに `SetMenu()` を追加しました。[@leaanthony](https://github.com/leaanthony) が実装しました
- 通知対応を追加しました。[#4098](https://github.com/wailsapp/wails/pull/4098) で [@popaprozac](https://github.com/popaprozac) が実装しました
-  mac にファイル関連付け対応を追加しました。[#4177](https://github.com/wailsapp/wails/pull/4177) で [@wimaha](https://github.com/wimaha) が実装しました
- セマンティックバージョンを上げるための `wails3 tool version` を追加しました。[@leaanthony](https://github.com/leaanthony) が実装しました
- macOS と Windows にバッジ対応を追加しました。[#](https://github.com/wailsapp/wails/pull/4234) で [@popaprozac](https://github.com/popaprozac) が実装しました
- 登録済み／厳密に型付けされたイベントへの対応を追加しました。[#4161](https://github.com/wailsapp/wails/pull/4161) で [@fbbdev](https://github.com/fbbdev) と [@IanVS](https://github.com/IanVS) が実装しました
- カスタムイベント用のフックを登録できるようにしました。[#4161](https://github.com/wailsapp/wails/pull/4161) で [@fbbdev](https://github.com/fbbdev) と [@IanVS](https://github.com/IanVS) が実装しました
- システムのファイルマネージャーでパス `path` を開くための `app.OpenFileManager(path string, selectFile bool)` を追加し、`selectFile` による任意の強調表示にも対応しました。[@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte) が実装しました
- `wails3 init` コマンドに新しい `-git` フラグを追加しました。[@leaanthony](https://github.com/leaanthony) が実装しました
- 新しい `wails3 generate webview2bootstrapper` コマンドを追加しました。[@leaanthony](https://github.com/leaanthony) が実装しました
- ランタイムを手動で初期化できるよう、ランタイムに `init()` メソッドを追加しました。[@leaanthony](https://github.com/leaanthony) が実装しました
- Window の WindowOptions に `WindowDidMoveDebounceMS` オプションを追加しました。[@leaanthony](https://github.com/leaanthony) が実装しました
- シングルインスタンス機能を追加しました。[@leaanthony](https://github.com/leaanthony) が実装しました。@APshenkin による [v2 PR](https://github.com/wailsapp/wails/pull/2951) を基にしています。
- `wails3 generate template` コマンド。[@leaanthony](https://github.com/leaanthony) が実装しました
- `wails3 releasenotes` コマンド。[@leaanthony](https://github.com/leaanthony) が実装しました
- `wails3 update cli` コマンド。[@leaanthony](https://github.com/leaanthony) が実装しました
- `wails3 generate bindings` コマンドの `-clean` オプション。[@leaanthony](https://github.com/leaanthony) が実装しました
- Linux で aarch64（arm64）AppImage をビルドできるようにしました。[#3981](https://github.com/wailsapp/wails/pull/3981) で [@AkshayKalose](https://github.com/AkshayKalose) が実装しました
- スポンサーへのハイパーリンクを追加しました。[#3958](https://github.com/wailsapp/wails/pull/3958) で @ansxuman が実装しました
- deb、rpm、Arch Linux パッケージャーによる Linux パッケージのビルドに対応。実装者：
- Darwin のユニバーサルビルドとパッケージに対応。実装者：
- Web サイトにイベントのドキュメントを追加。実装者：
- 非 SSR 開発用に設定された sveltekit および sveltekit-ts テンプレート
- 新しい `wails3 update build-assets` コマンドを使用してビルドアセットを更新。実装者：
- HTML Drag and Drop API をテストするためのサンプル。実装者：
- [leaanthony](https://github.com/leaanthony) によるファイル関連付け対応：
- 新しい `wails3 generate runtime` コマンド。実装者：
- ウィンドウを中央に配置するかどうかを指定する新しい `InitialPosition` オプション、または
- `application` パッケージに `Path` および `Paths` メソッドを追加。実装者：
- Windows の `GeneralAutofillEnabled` および `PasswordAutosaveEnabled` オプションを追加
- サービスメソッドを呼び出したウィンドウを取得できるようにしました。実装者：
- WebView2 に `EnabledFeatures` および `DisabledFeatures` オプションを追加。実装者：
- ⊞ 高度な高 DPI モニター対応のための新しい DIP システム。実装者：
- ⊞ ウィンドウクラス名のオプション。[windom](https://github.com/windom/) が以下で実装：
- プラグイン機能を提供できるようサービスを拡張しました。実装者：
- 🐧 WindowDidMove／WindowDidResize イベント：
- ⊞ WindowDidResize イベント：
-  Dock を処理できるよう Event ApplicationShouldHandleReopen を追加
-  実装に getPrimaryScreen/getScreens を追加。@tmclane が以下で実装：
-  macOS のフルスクリーンモードでツールバーを表示するためのオプションを追加。実装者：
- 🐧 Linux のキー入力をアクセラレーターに変換する onKeyPress ロジックを追加
- 🐧 タスク `run:linux` を追加。実装者：
- `SetIcon` メソッドをエクスポート。[@almas-x](https://github.com/almas-x) が以下で実装：
- `OnShutdown` を改善。[@almas-x](https://github.com/almas-x) が以下で実装：
- `Window` インターフェースに `ToggleMaximise` メソッドを復元。実装者：
- `Environment()` に詳細情報を追加。@leaanthony が以下で実装：
- `Window` インターフェースで `WebviewWindow.IsFocused` メソッドを公開。実装者：
- WML システムで、スペース区切りの複数のトリガーイベントに対応。実装者：
- バンドルされた JS ランタイムスクリプトからの ESM エクスポートを追加。実装者：
- 次の代わりにバンドルされた JS ランタイムスクリプトを使用するためのバインディングジェネレーターフラグを追加：
- Linux に `setIcon` を実装。[@abichinger](https://github.com/abichinger) が実装しました
- dev コマンドに `-port` フラグを追加し、環境変数に対応
- バインドされたメソッドの呼び出しに対するテストを追加。実装者：
- ⊞ 作成済みのウィンドウに `SetIgnoreMouseEvents` を追加。実装者：
-  ウィンドウのスタックレベル（順序）を設定できるようにしました。実装者：

### 修正

- Retina Mac で `Screen.Bounds`、`WorkArea`、`Size` が半分になる問題を修正しました。`Physical*` フィールドで NSScreen のポイント値をデバイスピクセルに変換し、トップレベルの `Screen.X`／`Y` に値を設定することで、マルチモニター間の接触検出と作業領域内への配置が正しく機能するようにしました。@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5168)
- ディスプレイ構成の変更時（例：スリープ／復帰中に外部モニターをホットプラグした場合）に WebKit DisplayLink のデッドロックを引き起こす ScreenManager のデータ競合を修正しました
- Assets.car が存在する場合、CFBundleIconName を appicon に直接設定するようにした（@symball による [PR](https://github.com/wailsapp/wails/pull/5154)）
- Fedora、openSUSE、Arch、NixOS で `wails3 doctor` が誤った WebKitGTK パッケージを報告する問題を修正。v3 はコンパイル時に 4.1 API を必要とするため、4.0 のフォールバックエントリを削除（#5071）
- openSUSE の webkit2gtk に対して doctor が示すパッケージ名を修正（`webkit2gtk4_1-devel` → 正しい openSUSE パッケージ名である `webkit2gtk3-devel`）（#5071）
- デスクトップ開発モードで `/wails/custom.js` が存在しない場合に発生する `Unexpected token '<'` エラーを修正。HTML SPA のフォールバックが JavaScript として挿入されるのを防ぐため、`/wails/custom.js` 用の明示的な 404 ハンドラーと、大文字・小文字を区別しない `Content-Type` 検証を `loadOptionalScript` に追加（[#5068](https://github.com/wailsapp/wails/issues/5068)）
- macOS のシステムトレイメニューのハイライト状態を修正。メニューを開いている間、アイコンが選択状態で表示されるようにした（#4910）
- macOS でシステムトレイに関連付けられたウィンドウがほかのウィンドウの背後に表示される問題を修正。適切なポップアップウィンドウレベルを使用するようにした（#4910）
- ドキュメント全体にある誤った `@wailsio/runtime` のインポート例を修正（#4989）
- darwin でフレームレスウィンドウを最小化できない問題を修正（#4294）
- `node_modules/` を go-task の最新状態チェックから除外し、`wails3 build` および `wails3 dev` の実行中に 20-30 分間停止する問題を修正。従来は `sources: "**/*"` glob により、go-task が `node_modules/` 内の全ファイル（MUI などの大規模な依存関係を含む 50000-100000+ 個）を列挙してチェックサムを計算していたため、特に Windows/NTFS で低速になっていた（#4939）
- C の `Screen` typedef と X11 の Xlib.h の衝突による GTK4 のビルド失敗を修正（#4957）
- macOS の Dock バッジメソッドの一貫性を修正
- `InvisibleTitleBarHeight` がフレームレスまたはタイトルバーが透明なウィンドウだけでなく、すべての macOS ウィンドウに適用される問題を修正（#4960）
- `InvisibleTitleBarHeight` が有効な場合、ウィンドウ端付近でドラッグを開始しないようにして、上隅からのサイズ変更時にウィンドウが揺れる問題を修正（#4960）
- JS/TS バインディングで enum キーを持つマップ型を生成する際の問題を修正（#4437、@fbbdev）
- Windows でディスプレイの拡大率が 100% 以外の場合に、ファイルのドラッグ＆ドロップが機能しない問題を修正
- Windows でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- Windows でファイルのドロップ座標が誤ったピクセル空間になる問題を修正（物理ピクセルと CSS ピクセル）
- Linux でホバー効果を使用すると、ファイルのドラッグ＆ドロップが安定して機能しない問題を修正
- Linux でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- `gtk_window_present()` を使用し、Linux/GTK4 でウィンドウの表示・非表示を切り替えた際に最小化状態へ戻ることがある問題を修正（#4957）
- `XTranslateCoordinates`/`XMoveWindow` を介して X11 条件付きサポートを追加し、Linux/GTK4 でウィンドウ位置の取得・設定が常に 0,0 を返す問題を修正（#4957）
- 削除された `gtk_window_set_geometry_hints` の代わりに、シグナルベースでサイズを上限に制限する処理を追加し、Linux/GTK4 でウィンドウの最大サイズが適用されない問題を修正（#4957）
- 適切な PhysicalBounds の計算と `gdk_monitor_get_scale` を介した小数スケーリングへの対応を実装し、Linux/GTK4 の DPI スケーリングを修正（GTK 4.14+）
- Linux/GTK4 で新しいウィンドウを作成するとメニュー項目が重複する問題を修正
- JS/TS バインディングで enum キーを持つマップ型を生成する際の問題を修正（#4437、@fbbdev）
- Windows でディスプレイの拡大率が 100% 以外の場合に、ファイルのドラッグ＆ドロップが機能しない問題を修正
- Windows でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- Windows でファイルのドロップ座標が誤ったピクセル空間になる問題を修正（物理ピクセルと CSS ピクセル）
- Linux でホバー効果を使用すると、ファイルのドラッグ＆ドロップが安定して機能しない問題を修正
- Linux でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- 適切な PhysicalBounds の計算と `gdk_monitor_get_scale` を介した小数スケーリングへの対応を実装し、Linux/GTK4 の DPI スケーリングを修正（GTK 4.14+）
- Linux/GTK4 で新しいウィンドウを作成するとメニュー項目が重複する問題を修正
- JS/TS バインディングで enum キーを持つマップ型を生成する際の問題を修正（#4437、@fbbdev）
- App.Window.Current() からメインスレッド上で AppKit API にアクセスしていなかったために macOS で発生する「ゴーストウィンドウ」問題を修正（#4947、@wimaha）
- WKUIDelegate runOpenPanelWithParameters を実装し、macOS で HTML の `<input type="file">` が機能しない問題を修正（#4862）
- macOS/Linux で `@wailsio/runtime` npm モジュールを使用すると、ネイティブのファイルドラッグ＆ドロップが機能しない問題を修正（#4953、@leaanthony）
- パッケージをまたぐ型エイリアスのバインディング生成を修正（#4578、@fbbdev）
- GTK のスレッドセーフティ違反により Linux で OpenFileDialog がクラッシュする問題を修正（#3683、@ddmoney420）
- 非表示または破棄済みのウィンドウで `Focus()` を呼び出すと SIGSEGV が発生する問題を修正（#4890、@ddmoney420）
- Linux で空のアイコンまたはビットマップを設定すると panic が発生する可能性がある問題を修正（#4923、@ddmoney420）
- macOS でサービスバインディングから ErrorDialog を呼び出すとクラッシュする問題を修正（#3631、@leaanthony）
- `v3\examples\dialogs` で Windows OS 上にメニューが表示されるようにした（@ndianabasi）
- ページの再読み込み中に TypeError を引き起こす競合状態を修正（#4872、@ddmoney420）
- `Collector.IsVoidAlias()` メソッドからグローバル状態を削除し、バインディングジェネレーターのテストが誤った出力を生成する問題を修正（#4941、@fbbdev）
- macOS で `<input type="file">` ファイルピッカーが機能しない問題を修正（#4862、@leaanthony）
- macOS で `Position()` と `SetPosition()` が異なる座標系を使用し、状態の保存・復元時にウィンドウ位置がずれる問題を修正（#4816、@leaanthony）
- アプリケーションマニフェストによって DPI 認識がすでに設定されている場合に発生する SetProcessDpiAwarenessContext の「Access is denied」エラーを修正（#4803）
- キーボードショートカットのドキュメントページを更新し、`KeyBinding.Add` のコールバック引数の型を修正（@ndianabasi）
- カスタムバインディングの生成に関するドキュメントを修正。`-o String` ではなく `-d String` を使用する必要がある
- `menu.Update()` でメニューの子項目がクリアされない問題を修正
- ドキュメント内の古い Manager API への参照を修正（31 ファイルを、`app.Window.New()`、`app.Event.Emit()` などの新しいパターンを使用するよう更新、@leaanthony）
- WebKit がシグナルハンドラーを上書きするため、JS にバインドされた Go メソッドで panic が発生すると Linux でクラッシュする問題を修正（#3965、@leaanthony）
- Linux で SaveFileDialog.SetFilename() が機能しない問題を修正（#4841、@samstanier）
- ドラッグ＆ドロップの例でドロップ座標が undefined と表示される問題を修正
- APP_NAME に空白が含まれる場合、macOS アプリバンドルの作成に失敗する問題を修正（ブレース展開の問題）
- Windows でサービスメソッドを呼び出すとインデックス範囲外の panic が発生する問題を修正（goccy/go-json を元に戻した）
- Windows でディスプレイの拡大率が 100% 以外の場合に、ファイルのドラッグ＆ドロップが機能しない問題を修正
- Windows でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- Windows でファイルのドロップ座標が誤ったピクセル空間になる問題を修正（物理ピクセルと CSS ピクセル）
- Linux でホバー効果を使用するとファイルのドラッグ＆ドロップが安定して動作しない問題を修正
- Linux でファイルドロップを有効にすると HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- @ndianabasi により、`APP_NAME` などの変数に含まれるスペースに対応するため、すべてのオペレーティングシステム向けの Taskfile.yml ファイル内の全コマンドを更新
- @wux1an により、Linux で 'build:universal:lipo:go' タスクを実行した際のコマンド引数エラーを修正
- @wux1an により、Linux で 'wails3 build GOOS=darwin GOARCH=arm64' を実行した際に発生する Docker エラー「undefined symbol: **<em>ubsan</em>handle_xxxxxxx」を修正
- @leaanthony により、カスタムプロトコルのドキュメントを統合し、Universal Links のセクションを追加
- @leaanthony により、TrackPopupMenuEx の同時呼び出しを防ぐガードを追加し、アイコンを繰り返しクリックした際に Windows のシステムトレイメニューがクラッシュする問題を修正（#4151）
- @leaanthony により、app.Run() より前に systray.Run() を呼び出した際にアプリがクラッシュすることを防止
- @leaanthony により、ApplicationShouldTerminateAfterLastWindowClosed が有効な状態で Hide()/Show() を使ってウィンドウの表示状態を切り替えると macOS でクラッシュする問題を修正（#4389）
- @leaanthony により、macOS および Windows でコンテキストメニューを繰り返し開いた際のメモリリークを修正（#4012）
- @leaanthony により、macOS でコンテキストメニューのネイティブリソースが再利用されず、表示のたびに新しいメニューが作成される問題を修正（#4012）
- @leaanthony により、`Hidden: true` を指定して起動したアプリで、macOS の Dock アイコンをクリックしても非表示のウィンドウが表示されない問題を修正（#4583）
- @leaanthony により、CGO 呼び出しでウィンドウポインターの型が誤っているため macOS の印刷ダイアログが開かない問題を修正（#4290）
- @leaanthony により、appmenu-gtk-module が未実体化のウィンドウにアクセスすることで発生する Wayland 上のウィンドウメニューのクラッシュを修正（#4769）
- @leaanthony により、アプリ名に無効な文字（スペース、括弧など）が含まれている場合に GTK アプリケーションがクラッシュする問題を修正
- @overlordtm により、Windows でドラッグ＆ドロップを初期化した際に発生する「not enough memory」エラーを修正（#4701）
- @leaanthony により、不正な URI エスケープが原因で Linux のファイルエクスプローラーが誤ったディレクトリを開く問題を修正（#4397）
- @leaanthony により、`.relr.dyn` ELF セクションを自動検出してストリッピングを無効にし、最新の Linux ディストリビューション（Arch、Fedora 39+、Ubuntu 24.04+）で AppImage のビルドが失敗する問題を修正（#4642）
- @leaanthony により、Fedora/DNF ベースのシステムで `wails doctor` が WebKit パッケージをインストール済みと誤って報告する問題を修正（#4457）
- @mbaklor により、デフォルトの `config.yml` が本番ビルドで `wails3 dev` を実行する問題を修正
- @leaanthony により、存在しないパッケージのインポートが原因で iOS サービススタブのビルドが失敗する問題を修正
- @leaanthony により、debug/info メソッドで構造化ログを使用すると「no formatting directives」エラーが発生する問題を修正
- @leaanthony により、モバイルプラットフォームのマージから誤って含まれた一時的なデバッグ出力文を削除
- @leaanthony により、DMA-BUF レンダラーを自動的に無効化し、NVIDIA GPU を使用する Wayland 上で WebKitGTK がクラッシュする問題（Error 71 Protocol error）を修正
- Linux の `application.WebviewWindowOptions.BackgroundColour` でアルファ値が無視される問題を解決（[#4722](https://github.com/wailsapp/wails/pull/4722)、@BradHacker）
- カスタムアイコンが指定されていない場合に、Windows のシステムトレイアイコンでアプリケーションアイコンがデフォルトとして使用されない問題を修正（#4704）
- ユーザーが作成したハンドルのみを破棄するよう `HICON` の所有権を追跡し、Explorer の再利用時のクラッシュを防止（#4653）。
- 破棄時に Windows のシステムテーマリスナーと保持中のトレイアイコンを解放し、goroutine とデバイスコンテキストのリークを解消（#4653）。
- サロゲートペアとマルチバイトグリフの破損を防ぐため、トレイのツールチップを 127 UTF-16 単位で切り詰め（#4653）。
- Windows の package タスクが失敗する問題を修正（#4667）
- Linux の Taskfile にある AppImage の appicon 変数を修正 [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- go-webview2 の v1.0.22 シグネチャ変更が原因で発生する Windows のビルドエラーを修正（#4513、#4645）
- Linux の Taskfile にある AppImage の appicon 変数を修正 [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- @Tolfx により、linux desktop.tmpl のプロトコル範囲から `<.Info.Protocol>`～`<.Protocol>` を削除して修正（#4510）
- @Etesam913 により、liquid glass デモの再定義エラーを修正（[#4542](https://github.com/wailsapp/wails/pull/4542)）
- [@JackDoan](https://github.com/JackDoan) により、Linux のシステムトレイメニュー更新を修正（[#4604](https://github.com/wailsapp/wails/issues/4604)）
- @leaanthony により、非表示ウィンドウの作成時に Windows で白いウィンドウが表示される問題を修正（[#4612](https://github.com/wailsapp/wails/pull/4612)）
- @rxliuli により、ドキュメント内の notifications パッケージのインポートパスを修正（[#4617](https://github.com/wailsapp/wails/pull/4617)）
- @leaanthony により、npm パッケージ @wailsio/runtime の使用時にドラッグ＆ドロップが機能しない問題（#4489）を #4616 で修正
- Windows：起動時にウィンドウがちらつく問題と、非表示ウィンドウが誤って表示される問題を @leaanthony が [PR](https://github.com/wailsapp/wails/pull/4600) で修正。
- [@samstanier](https://github.com/samstanier) により、Wayland でウィンドウを最大化した際のサイズ問題を修正（https://github.com/wailsapp/wails/issues/4429）
- [@samstanier](https://github.com/samstanier) により、Wayland でウィンドウを最大化した際のサイズ問題を修正（https://github.com/wailsapp/wails/issues/4429）
- @Etesam913 により、liquid glass デモの再定義エラーを修正（[#4542](https://github.com/wailsapp/wails/pull/4542)）
- @jghiloni により、MacOS で AssetServer がクラッシュする可能性がある問題を修正（[#4576](https://github.com/wailsapp/wails/pull/4576)）
- NextJs でのビルド時に発生するコンパイルの問題を修正。@rev42 により [#4585](https://github.com/wailsapp/wails/pull/4585) で修正
- @riadafridishibly により、ナイトリーリリース用パイプラインを修正（[#4597](https://github.com/wailsapp/wails/pull/4597)）
- @Etesam913 により、liquid glass デモの再定義エラーを修正（[#4542](https://github.com/wailsapp/wails/pull/4542)）
- @jghiloni により、MacOS で AssetServer がクラッシュする可能性がある問題を修正（[#4576](https://github.com/wailsapp/wails/pull/4576)）
- NextJs でのビルド時に発生するコンパイルの問題を修正。@rev42 により [#4585](https://github.com/wailsapp/wails/pull/4585) で修正
- @riadafridishibly により、ナイトリーリリース用パイプラインを修正（[#4597](https://github.com/wailsapp/wails/pull/4597)）
- @Etesam913 により、liquid glass デモの再定義エラーを修正（[#4542](https://github.com/wailsapp/wails/pull/4542)）
- @PPTGamer により、Windows の SetBackgroundColour を修正（[PR](https://github.com/wailsapp/wails/pull/4492)）
- @yulesxoxo により、Manager API のリファクタリングによる変更を反映するようドキュメントを更新（[PR #4476](https://github.com/wailsapp/wails/pull/4476)）
- Linux の Taskfile にある .desktop ファイルの appicon 変数を修正 [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- @yulesxoxo により、Manager API のリファクタリングによる変更を反映するようドキュメントを更新（[PR #4476](https://github.com/wailsapp/wails/pull/4476)）
- [#4456](https://github.com/wailsapp/wails/issues/4456) で報告された Windows の nil ポインター参照解除バグを、@leaanthony が [#4460](https://github.com/wailsapp/wails/pull/4460) で修正
- macOS の WKWebView で `allowsBackForwardNavigationGestures` をサポートし、2 本指スワイプによるナビゲーションジェスチャーを有効化（#1857）
- 初期状態で無効に設定されたメニュー項目で onClick が動作しない問題を、@leaanthony が [PR #4469](https://github.com/wailsapp/wails/pull/4469) で修正しました。初期調査を行った @IanVS に感謝します。
- ビルド失敗時に Vite サーバーがクリーンアップされない問題を修正（#4403）
- Windows で `SaveFileDialog` を閉じるかキャンセルしたときに発生するパニックを修正しました。@hkhere が [PR](https://github.com/wailsapp/wails/pull/4284) で修正しました。
- Windows における HTML レベルのドラッグ＆ドロップを、[@mbaklor](https://github.com/mbaklor) が [#4259](https://github.com/wailsapp/wails/pull/4259) で修正しました。
- 2 本指のスワイプによるナビゲーションジェスチャーを有効にするため、macOS WKWebView で `allowsBackForwardNavigationGestures` をサポート（#1857）
- 初期状態で無効に設定されたメニュー項目で onClick が動作しない問題を、@leaanthony が [PR #4469](https://github.com/wailsapp/wails/pull/4469) で修正しました。初期調査を行った @IanVS に感謝します。
- ビルド失敗時に Vite サーバーがクリーンアップされない問題を修正（#4403）
- Windows の通知解析を、@popaprozac が [PR](https://github.com/wailsapp/wails/pull/4450) で修正しました。
- doctor コマンドによる Windows SDK の依存関係チェックを、[@kodumulo](https://github.com/kodumulo) が [#4390](https://github.com/wailsapp/wails/issues/4390) で修正しました。
- Mac の processURLRequest における nil ポインターの参照解除を、[@etesam913](https://github.com/etesam913) が [#4366](https://github.com/wailsapp/wails/pull/4366) で修正しました。
- フィルター付きダイアログを使用できなくする Linux のバグを、[@bh90210](https://github.com/bh90210) が [#4287](https://github.com/wailsapp/wails/pull/4287) で修正しました。
- Windows と Linux の「編集」メニューの問題を、[@leaanthony](https://github.com/leaanthony) が [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062) で修正しました。
- macOS の .plist ファイルにある最小システムバージョンを、[@AkshayKalose](https://github.com/AkshayKalose) が [#3981](https://github.com/wailsapp/wails/pull/3981) で 10.13.0 から 10.15.0 に更新しました。
- ウィンドウ ID が飛ばされる問題を、[@leaanthony](https://github.com/leaanthony) が修正しました。
- RegisterContextMenu の呼び出し時に発生する nil メニューの問題を、[@leaanthony](https://github.com/leaanthony) が修正しました。
- バインディングジェネレーターの出力における依存関係の循環を、[@fbbdev](https://github.com/fbbdev) が [#4001](https://github.com/wailsapp/wails/pull/4001) で修正しました。
- バインディングジェネレーターの出力における定義前使用エラーを、[@fbbdev](https://github.com/fbbdev) が [#4001](https://github.com/wailsapp/wails/pull/4001) で修正しました。
- ビルドフラグをバインディングジェネレーターに渡すよう、[@fbbdev](https://github.com/fbbdev) が [#4023](https://github.com/wailsapp/wails/pull/4023) で変更しました。
- Windows 以外のプラットフォームでも動作するよう、Windows の Taskfile 内のパスをスラッシュ区切りに変更しました（[@leaanthony](https://github.com/leaanthony)）。
- Mac と Mac JS のイベントを、[@leaanthony](https://github.com/leaanthony) が修正しました。
- macOS のイベントデッドロックを、[@leaanthony](https://github.com/leaanthony) が修正しました。
- Windows で HTML は指定されているものの JS が指定されていない場合に、Window の初期化で発生する `Parameter incorrect` エラーを、[@leaanthony](https://github.com/leaanthony) が修正しました。
- アセットサーバーでコンテンツタイプのスニッフィングに使用するレスポンスプレフィックスのサイズを、[@fbbdev](https://github.com/fbbdev) が [#4049](https://github.com/wailsapp/wails/pull/4049) で修正しました。
- アセットサーバーのルートインデックスパスにおける、404 以外のレスポンスの処理を、[@fbbdev](https://github.com/fbbdev) が [#4049](https://github.com/wailsapp/wails/pull/4049) で修正しました。
- ジェネリック型のプロパティをテストする際にバインディングジェネレーターで発生する未定義動作を、[@fbbdev](https://github.com/fbbdev) が [#4045](https://github.com/wailsapp/wails/pull/4045) で修正しました。
- 基になる型のプロパティが名前付きラッパーと同一でない場合のモデル用バインディングジェネレーター出力を、[@fbbdev](https://github.com/fbbdev) が [#4045](https://github.com/wailsapp/wails/pull/4045) で修正しました。
- マップのキー型と前処理に関するバインディングジェネレーターの出力を、[@fbbdev](https://github.com/fbbdev) が [#4045](https://github.com/wailsapp/wails/pull/4045) で修正しました。
- マーシャラーインターフェースを実装する構造体に関するバインディングジェネレーターの出力を、[@fbbdev](https://github.com/fbbdev) が [#4045](https://github.com/wailsapp/wails/pull/4045) で修正しました。
- バインディングジェネレーターにおける、ジェネリック型を含む型循環の検出を、[@fbbdev](https://github.com/fbbdev) が [#4045](https://github.com/wailsapp/wails/pull/4045) で修正しました。
- バインディングジェネレーターの出力に含まれる、エクスポートされていないモデルへの無効な参照を、[@fbbdev](https://github.com/fbbdev) が [#4045](https://github.com/wailsapp/wails/pull/4045) で修正しました。
- 注入されるコードをサービスファイルの末尾へ移動しました（[@fbbdev](https://github.com/fbbdev)、[#4045](https://github.com/wailsapp/wails/pull/4045)）。
- バインディングジェネレーターにおけるファイルクローズ処理のエラー処理を、[@fbbdev](https://github.com/fbbdev) が [#4045](https://github.com/wailsapp/wails/pull/4045) で修正しました。
- ライフサイクルメソッドまたは http メソッドを定義しているものの、それ以外のバインドされたメソッドがないサービスに対する警告を、[@fbbdev](https://github.com/fbbdev) が [#4045](https://github.com/wailsapp/wails/pull/4045) で抑制しました。
- システムのカラースキームがライトの場合に、React 以外のテンプレートで Hello World フッターが表示されない問題を、[@marcus-crane](https://github.com/marcus-crane) が [#4056](https://github.com/wailsapp/wails/pull/4056) で修正しました。
- macOS で非表示になっていたメニュー項目を、[@leaanthony](https://github.com/leaanthony) が修正しました。
- メッセージプロセッサーにおけるエラーの処理と書式設定を、[@fbbdev](https://github.com/fbbdev) が [#4066](https://github.com/wailsapp/wails/pull/4066) で修正しました。
-  アプリケーション終了時にサービスのシャットダウンが省略される問題を、[@fbbdev](https://github.com/fbbdev) が [#4066](https://github.com/wailsapp/wails/pull/4066) で修正しました。
-  メニューの更新がメインスレッドで実行されるようにしました（[@leaanthony](https://github.com/leaanthony)）。
- ドラッグとサイズ変更の仕組みをより堅牢にし、プラットフォームで期待される動作にさらに近づけました（[@fbbdev](https://github.com/fbbdev)、[#4100](https://github.com/wailsapp/wails/pull/4100)）。
- Webpack/angular がランタイム初期化コードを破棄する問題 [#4097](https://github.com/wailsapp/wails/issues/4097) を、[@fbbdev](https://github.com/fbbdev) が [#4100](https://github.com/wailsapp/wails/pull/4100) で修正しました。
- 初期状態で非表示のメニュー項目を、[@IanVS](https://github.com/IanVS) が [#4116](https://github.com/wailsapp/wails/pull/4116) で修正しました。
- 拡張子のないリクエストで `[request]` が存在せず、`[request].html` が存在する場合に、assetFileServer が `.html` ファイルを配信しない問題を修正しました。
- アイコン生成パスを、[@robin-samuel](https://github.com/robin-samuel) が [#4125](https://github.com/wailsapp/wails/pull/4125) で修正しました。
- `fullscreen`、`unfullscreen`、`unminimise`、`unmaximise` の各イベントが発行されない問題を、[@oSethoum](https://github.com/osethoum) が [#4130](https://github.com/wailsapp/wails/pull/4130) で修正しました。
- 設定内のデフォルトバージョンに誤ったプレフィックスが付いていたために発生する NSIS Error を、[@robin-samuel](https://github.com/robin-samuel) が [#4126](https://github.com/wailsapp/wails/pull/4126) で修正しました。
- Windows で Dialogs ランタイム関数がエスケープされたパスを返す問題を、[TheGB0077](https://github.com/TheGB0077) が [#4188](https://github.com/wailsapp/wails/pull/4188) で修正しました。
- HKCU における Webview2 の検出パスを、[@leaanthony](https://github.com/leaanthony) が修正しました。
- macOS の入力に関する問題を、[@leaanthony](https://github.com/leaanthony) が修正しました。
- [@yulesxoxo](https://github.com/yulesxoxo) が [#4219](https://github.com/wailsapp/wails/pull/4219) で Windows アイコン生成タスクのファイル名を修正しました。
- @kron の作業を基に、[@leaanthony](https://github.com/leaanthony) がフレームレスウィンドウの透過に関する問題を修正しました。
- @kron の作業を基に、[@leaanthony](https://github.com/leaanthony) がウィンドウが無効または最小化されている場合のフォーカス呼び出しを修正しました。
- @kron の作業を基に、[@leaanthony](https://github.com/leaanthony) がタスクバーの再起動後にシステムトレイが表示されない問題を修正しました。
- [#4245](https://github.com/wailsapp/wails/pull/4245) で、fallbackResponseWriter が Flush() を実装していない問題を修正しました
- [@superDingda] が [#4236](https://github.com/wailsapp/wails/issues/4236) で、fallbackResponseWriter が Flush() を実装していない問題を修正しました
- [@joshhardy](https://github.com/joshhardy) が [#4354](https://github.com/wailsapp/wails/pull/4354) で、非同期の Go バインド関数呼び出しが保留中のときに macOS でウィンドウを閉じるとクラッシュする問題を修正しました
- [@leaanthony](https://github.com/leaanthony) が Windows の効率モード起動時の競合状態を修正しました
- [@leaanthony](https://github.com/leaanthony) が Windows のアイコンハンドルのクリーンアップを修正しました。
- [@PPTGamer](https://github.com/PPTGamer) が [#4375](https://github.com/wailsapp/wails/pull/4375) で、Windows 上の `OpenFileManager` を修正しました。
- @atterpac が [#3979](https://github.com/wailsapp/wails/pull/3979) で、Linux の最小幅／最大幅オプションを修正しました
- @atterpac が [#3966](https://github.com/wailsapp/wails/pull/3966) で、npm のバージョン更新により Typescript テンプレートの型定義を更新しました
- @atterpac が [#3945](https://github.com/wailsapp/wails/pull/3945) で、Sveltekit テンプレートの CSS 参照を修正しました
- [@leaanthony](https://github.com/leaanthony) が、window run() 内の主要なコールバックがメインスレッドで呼び出されるようにしました
- [@leaanthony](https://github.com/leaanthony) が、ダイアログのディレクトリ選択例を修正しました
- [@leaanthony](https://github.com/leaanthony) が、index.html が存在しない場合に表示する新しい中国語のエラーページを作成しました
-  [@leaanthony](https://github.com/leaanthony) が、`windowDidBecomeKey` コールバックがメインスレッドで実行されるようにしました
-  [@leaanthony](https://github.com/leaanthony) が、フレームレスウィンドウのフルスクリーン表示に対応しました
-  [@leaanthony](https://github.com/leaanthony) が、ウィンドウ破棄ロジックを改善しました
-  [@leaanthony](https://github.com/leaanthony) が、システムトレイに関連付けられたウィンドウの位置決定ロジックを修正しました
-  [@leaanthony](https://github.com/leaanthony) が、フレームレスウィンドウのフルスクリーン表示に対応しました
- [@leaanthony](https://github.com/leaanthony) がイベント処理を修正しました
- [@leaanthony](https://github.com/leaanthony) がウィンドウのシャットダウンロジックを修正しました
- [@leaanthony](https://github.com/leaanthony) により、共通 Taskfile では Typescript テンプレート用の Typescript バインディングがデフォルトで生成されるようになりました
- [@mmalcek](https://github.com/mmalcek) が [#3990](https://github.com/wailsapp/wails/pull/3990) で、ウィンドウが開いていない場合やシステムトレイのみの場合に、WM_CLOSE メッセージでアプリケーションを終了する処理を修正しました
- @5aaee9 が [#3192](https://github.com/wailsapp/wails/pull/3192) で、garble ビルドを修正しました
- [@leaanthony](https://github.com/leaanthony) が Windows の NSIS ビルドを修正しました
- 閉じられていないことが原因で Linux の複数選択ダイアログに発生するデッドロックを修正しました
- Windows ビルド時の .syso ファイルのクロスプラットフォームなクリーンアップを修正しました。担当：
- @atterpac が amd64 AppImage のコンパイルを修正しました：
- @ansxuman がビルドアセットの更新を修正しました：
- @atterpac が Linux のシステムトレイにおける `OnClick` と `OnRightClick` の実装を修正しました
- Mac で `AlwaysOnTop` が動作しない問題を修正しました。担当：
-  `application.NewEditMenu` に重複したものが含まれる問題を修正しました
- 🐧 aarch64 のコンパイルを修正しました
- ⊞ ラジオグループのメニュー項目を修正しました。担当：
- MacOS で実行可能な .app をビルドする際に、'name' と 'outputfilename' が原因で発生するエラーを修正しました
- ドラッグ＆ドロップの例で customEventProcessor を使用した際のバグを修正しました。担当：
- 🐧 IgnoreMouseEvents の追加によって発生した Linux のコンパイルエラーを修正しました。担当：
- ⊞ syso アイコンファイル生成のバグを修正しました。担当：
- 🐧 Wayland 上でネイティブ実行するための修正を取り込みました。出典：
- 次では内部サービスメソッドをバインドしないようにしました：
- ⊞ 次でシステムトレイ起動時のパニックを修正しました：
- 次では内部サービスメソッドをバインドしないようにしました：
- ⊞ 次でシステムトレイ起動時のパニックを修正しました：
- メニュー項目とイベント処理を大幅にリファクタリングしました。現時点では主に macOS が改善されています。担当：
- プラグインとイベントのリファクタリング後にテストを修正しました：
- ⊞ `Failed to unregister class Chrome_WidgetWin_0` の警告を修正しました。担当：
- モジュールの問題
- [atterpac](https://github.com/atterpac) がサイズ変更イベントのメッセージ処理を修正しました：
- 🐧 NixOS でのテーマ処理エラーを修正しました。担当：
- Windows でボリュームをまたぐプロジェクトのインストールを修正しました。担当：
- フッターが表示されるように React テンプレートの CSS を修正しました。担当：
- refresh を最新版に更新し、devmode で作業中に発生するゾンビプロセスを修正しました
- [Atterpac](https://github.com/atterpac) が AppImage の WebKit ファイル取得処理を修正しました
- [Atterpac](https://github.com/Atterpac) が Doctor の apt パッケージ検証を修正しました：
- @5aaee9 がアプリケーション終了時にフリーズする問題（Darwin）を修正しました：
- Windows 上のサンプルの背景色を修正しました。担当：
- [mmghv](https://github.com/mmghv) がデフォルトのコンテキストメニューを修正しました：
- Darwin の矢印キーの 16 進値を修正しました。担当：
- Windows のドラッグ＆ドロップを動作する状態にしました。追加者：
- ユーザーに適切なドライバーがない場合に Doctor で発生する Linux のバグを修正しました
- 起動時の DPI スケーリング（Windows）を修正しました。変更者：[@almas-x](https://github.com/almas-x)、変更箇所：
- 相対パスを使用するように `go.mod` の置換行を修正しました。次を含む Windows パスを修正します：
- ウィンドウが関連付けられていない場合の MacOS システムトレイのクリック処理を修正しました。担当：
- 不明なオプションが原因で失敗する Windows ビルドを修正しました。担当：
- 関連付けられたものがない状態で Windows のシステムトレイアイコンを左クリックするとクラッシュする問題を修正しました
- @5aaee9 が、ウィンドウを 2 回開いたときに baseURL が誤る問題を PR で修正しました
- `WebviewWindow.Restore` メソッド内の if 分岐の順序を修正しました。担当：
- 複数回の `GetStartURL` 呼び出しにわたって `startURL` を正しく計算するようにしました。条件：
- `Screen` 構造体の JS 型を対応する Go 型と一致するように修正しました。担当：
- 登録済みイベントが適切にクリーンアップされるよう、`WML.Reload` メソッドを修正しました
- Linux でカスタムコンテキストメニューがすぐに閉じる問題を修正しました。担当：
- バインディングによって生成されるモデルファイルの出力パスと拡張子を修正
- バインディングによって生成される JS コード内のモデルファイルのインポートパスを修正
- 一部の Linux ディストリビューションでドラッグ＆ドロップが動作しない問題を修正。対応者：
- `wails3 task dev` の使用時に macOS のタスクが欠落する問題を修正。対応者：
- イベントの登録時に nil マップへの代入が発生する問題を修正。対応者：
- バインドされたメソッドのパラメーターをアンマーシャリングできない問題を修正。対応者：
- バインドされたメソッドからの複数の戻り値の処理を修正。対応者：
- システムのパッケージマネージャーでインストールされていない npm を doctor が検出できない問題を修正
- MicrosoftEdgeWebview2Setup.exe が欠落する問題を修正。協力者：
- ウィンドウ ID の処理が原因で Linux 上でランダムにクラッシュする問題を @leaanthony が修正。以下を基にしています：
- Linux 上で systemTray.setIcon がクラッシュする問題を修正。対応者：
- `setFrameless` 関数の初回呼び出し時にウィンドウフレームが確実に適用されるよう修正。対象：

### 変更

- **破壊的変更**：Go のマップのセマンティクスを正確に反映するため、生成される JS/TS バインディングのマップキーを省略可能としてマークするようになりました。TypeScript でマップの値にアクセスした際、`T` ではなく `T | undefined` が返されるようになったため、null チェックまたはアサーションが必要です（#4943）。対応者：`@fbbdev`
- `@wailsio/runtime` の変更に合わせて、`Event` の使用を `Events` に変更し、`Features/Events/Event System` のドキュメント内の関数呼び出しも適切に変更。対応者：@AbdelhadiSeddar
- `EnabledFeatures`、`DisabledFeatures`、`AdditionalBrowserArgs` をウィンドウ単位のオプションからアプリケーションレベルの `Options.Windows` に移動（#4559）。対応者：@leaanthony
- `Drag N Drop` のサンプルの README を更新し、このサンプルで `Internal Drag and Drop` を実演していることを明記。対応者：@ndianabasi
- 各種デバッグログのレベルを Info から Debug に変更（対応者：@mbaklor）
- <strong>破壊的変更：</strong>ウィンドウオプションの `EnableDragAndDrop` を `EnableFileDrop` に名称変更
- <strong>破壊的変更：</strong>イベントコンテキストの `DropZoneDetails` を `DropTargetDetails` に名称変更
- **破壊的変更：**`WindowEventContext` の `DropZoneDetails()` メソッドを `DropTargetDetails()` に名称変更
- **破壊的変更：**`WindowDropZoneFilesDropped` イベントを削除し、代わりに `WindowFilesDropped` を使用
- <strong>破壊的変更：</strong>HTML 属性を `data-wails-dropzone` から `data-file-drop-target` に変更
- <strong>破壊的変更：</strong>CSS のホバークラスを `wails-dropzone-hover` から `file-drop-target-active` に変更
- <strong>破壊的変更：</strong>Windows から `DragEffect`、`OnEnterEffect`、`OnOverEffect` オプションを削除（削除された IDropTarget の一部だったため）
- 実行時のすべての JSON 処理（メソッドバインディング、イベント、WebView リクエスト、通知、kvstore）を goccy/go-json に切り替え、パフォーマンスを 21-63% 向上させ、メモリ割り当てを 40-60% 削減
- BoundMethod 構造体のレイアウトを最適化し、isVariadic フラグをキャッシュして呼び出しごとのオーバーヘッドを削減
- 引数が `<=8` 個のメソッドではスタックに割り当てた引数バッファを使用し、ヒープ割り当てを回避
- メソッド呼び出し時の結果収集を最適化し、戻り値が 1 つの場合のスライス割り当てを回避
- MIME タイプのキャッシュに sync.Map を使用し、並行処理性能を向上
- HTTP トランスポートのリクエストボディ読み取りにバッファプールを使用
- コンテンツタイプスニファーの CloseNotify チャネルを遅延割り当てし、リクエストごとの割り当てを削減
- アセットサーバーから CSS のデバッグログを削除
- MIME タイプの拡張子マップを拡充し、一般的な Web 形式（フォント、音声、動画など）を 50 種類以上カバー
- Window の `X/Y` オプションに関するドキュメントを更新。対応者：@ruhuang2001
- フロントエンドバインディング生成用のオプションを追加し、`Frontend Runtime` のドキュメントを更新。対応者：@ndianabasi
- Wails の v3 Asset Server に関するドキュメントページを更新。対応者：@ndianabasi
- **破壊的変更**：パッケージレベルのダイアログ関数（`application.InfoDialog()`、`application.QuestionDialog()` など）を削除。代わりに `app.Dialog` マネージャー（`app.Dialog.Info()`、`app.Dialog.Question()`、`app.Dialog.Warning()`、`app.Dialog.Error()`、`app.Dialog.OpenFile()`、`app.Dialog.SaveFile()`）を使用してください
- 実際の API に合わせてダイアログのドキュメントを更新：コールバック付きの `app.Dialog.*`、`AddButton()`（`SetButtons()` ではない）、`SetDefaultButton(*Button)`（文字列ではない）、`AddFilter()`（`SetFilters()` ではない）、`SetFilename()`（`SetDefaultFilename()` ではない）、およびフォルダー選択用の `app.Dialog.OpenFile().CanChooseDirectories(true)` を使用
- **破壊的変更**：本番ビルドがデフォルトになりました。開発ビルドを作成するには、Taskfile で `DEV=true` を設定してください。サンプル用に新しいプロジェクトを生成。対応者：@leaanthony
- データ引数が 0 個または 1 個のカスタムイベントを発行すると、データ値はスライスにラップされず、Data フィールドへ直接代入されます。対応者：[@fbbdev](https://github.com/fbbdev)、[#4633](https://github.com/wailsapp/wails/pull/4633)
- Windows のトレイが `NIS_HIDDEN` を切り替えることで `SystemTray.Show()`/`Hide()` に従うようになり、アプリを完全に非表示にして再表示できるようになりました（#4653）。
- トレイの登録で解決済みのアイコンを再利用し、`NOTIFYICON_VERSION_4` を一度だけ設定し、`NIF_SHOWTIP` を有効にすることで、Explorer の再起動後にツールチップが復元されるようになりました（#4653）。
- macOS：メニューバーと Dock の領域を除外してウィンドウを中央配置するため、`frame` ではなく `visibleFrame` を使用
- macOS：メニューバーと Dock の領域を除外してウィンドウを中央配置するため、`frame` ではなく `visibleFrame` を使用
- `-config` パラメーターを指定して `wails3 update build-assets` を実行すると、`-product*` パラメーターで設定した値は
- `window.NativeWindowHandle()` -> `window.NativeWindow()`。対応者：@leaanthony、[#4471](https://github.com/wailsapp/wails/pull/4471)
- 内部のウィンドウ処理をリファクタリング。対応者：@leaanthony、[#4471](https://github.com/wailsapp/wails/pull/4471)
- `application.WindowIDKey` と `application.WindowNameKey` を削除（`application.WindowKey` に置換）。対応者：[@leaanthony](https://github.com/leaanthony)
- ContextMenuData が any ではなく string を返すように変更。対応者：[@leaanthony](https://github.com/leaanthony)
- JS/TS バインディングで、固定長配列型のクラスフィールドが空ではなく、想定される長さで初期化されるようになりました。対応者：[@fbbdev](https://github.com/fbbdev)、[#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData が any ではなく string を返すように変更。対応者：[@leaanthony](https://github.com/leaanthony)
- `application.NewService` は、省略可能なパラメーターとして options を受け付けなくなりました（代わりに `application.NewServiceWithOptions` を使用してください）。対応者：[@leaanthony](https://github.com/leaanthony)、[#4024](https://github.com/wailsapp/wails/pull/4024)
- `nanoid` への依存を削除。対応者：[@leaanthony](https://github.com/leaanthony)
- mica/acrylic/tabbed ウィンドウスタイルに関する Window のサンプルを更新。対応者：[@leaanthony](https://github.com/leaanthony)
- JS/TS バインディングで `internal.js/ts` モデルファイルを削除しました。すべてのモデルは `models.js/ts` にあります。対応者：[@fbbdev](https://github.com/fbbdev)、[#4045](https://github.com/wailsapp/wails/pull/4045)
- JS/TS バインディングでは、名前付き型を別の名前付き型のエイリアスとしてレンダリングしなくなりました。従来の動作はエイリアスのみに制限されます。対応者：[@fbbdev](https://github.com/fbbdev)、[#4045](https://github.com/wailsapp/wails/pull/4045)
- JS/TS バインディングのクラスモードでは、型パラメーター型の構造体フィールドが省略可能としてマークされ、自動的に初期化されないようになりました（[@fbbdev](https://github.com/fbbdev)、[#4045](https://github.com/wailsapp/wails/pull/4045)）
- テンプレートから ESLint を削除しました（[@IanVS](https://github.com/IanVS)、[#4059](https://github.com/wailsapp/wails/pull/4059)）
- 著作権表示の年を2025に更新しました（[@IanVS](https://github.com/IanVS)、[#4037](https://github.com/wailsapp/wails/pull/4037)）
- event.Sender のドキュメントを追加しました（[@IanVS](https://github.com/IanVS)、[#4075](https://github.com/wailsapp/wails/pull/4075)）
- Go 1.24 をサポートしました（[@leaanthony](https://github.com/leaanthony)）
- `ServiceStartup` フックが、`application.New`内ではなく`App.Run`の呼び出し時に実行されるようになりました（[@fbbdev](https://github.com/fbbdev)、[#4066](https://github.com/wailsapp/wails/pull/4066)）
- `ServiceStartup` エラーでプロセスを終了せず、`App.Run`から返すようになりました（[@fbbdev](https://github.com/fbbdev)、[#4066](https://github.com/wailsapp/wails/pull/4066)）
- JS からのバインディングおよびダイアログ呼び出しが、文字列ではなくエラーオブジェクトで reject されるようになりました（[@fbbdev](https://github.com/fbbdev)、[#4066](https://github.com/wailsapp/wails/pull/4066)）
- Windows でのシステムトレイメニューの配置を改善しました（[@leaanthony](https://github.com/leaanthony)）
- JS ランタイムを TypeScript に移植しました（[@fbbdev](https://github.com/fbbdev)、[#4100](https://github.com/wailsapp/wails/pull/4100)）
- ランタイムはインポート直後に初期化されるため、ウィンドウの読み込みを待つ必要がなくなりました（[@fbbdev](https://github.com/fbbdev)、[#4100](https://github.com/wailsapp/wails/pull/4100)）
- ランタイムは init メソッドをエクスポートしなくなりました。副作用のみを目的としたインポートで初期化できます（[@fbbdev](https://github.com/fbbdev)、[#4100](https://github.com/wailsapp/wails/pull/4100)）
- バインドされたメソッドは、キャンセルされた場合に`CancelError`で reject される`CancellablePromise`を返すようになりました。実際の呼び出し結果は破棄されます（[@fbbdev](https://github.com/fbbdev)、[#4100](https://github.com/wailsapp/wails/pull/4100)）
- 組み込みサービス型の名前を`Service`に統一しました（[@fbbdev](https://github.com/fbbdev)、[#4067](https://github.com/wailsapp/wails/pull/4067)）
- オプションを取る組み込みサービス作成関数の名前を`NewWithConfig`に統一しました（[@fbbdev](https://github.com/fbbdev)、[#4067](https://github.com/wailsapp/wails/pull/4067)）
- Go API との一貫性を保つため、`sqlite`サービスの`Select`メソッドを`Query`に改名しました（[@fbbdev](https://github.com/fbbdev)、[#4067](https://github.com/wailsapp/wails/pull/4067)）
- テンプレート：ランタイムを「dependencies」に移し、package.json ファイルを整理しました（[@IanVS](https://github.com/IanVS)、[#4133](https://github.com/wailsapp/wails/pull/4133)）
- 特定の macOS API を利用できるように、開発時にアプリバンドルを作成してアドホック署名するようになりました（[@popaprozac](https://github.com/popaprozac)、[#4171](https://github.com/wailsapp/wails/pull/4171)）
- ビルドアセットをプラットフォーム別のディレクトリへ移動しました（[@leaanthony](https://github.com/leaanthony)）
- Taskfile をプラットフォーム別のディレクトリへ移動し、名前を変更しました（[@leaanthony](https://github.com/leaanthony)）
- `index.html`が見つからない場合のユーザー体験を大幅に改善しました（[@leaanthony](https://github.com/leaanthony)）
- [Windows] 最小化と復元のパフォーマンスを改善しました（[@leaanthony](https://github.com/leaanthony)）。[562589540](https://github.com/562589540)による元の[PR](https://github.com/wailsapp/wails/pull/3955)を基にしています
- `ShouldClose`オプションを削除しました（代わりに events.Common.WindowClosing のフックを登録してください）（[@leaanthony](https://github.com/leaanthony)）
- [Windows] ウィンドウを開く際のちらつきを軽減しました（[@leaanthony](https://github.com/leaanthony)）
- 内部関数として意図されていたため、`Window.Destroy`を削除しました（[@leaanthony](https://github.com/leaanthony)）
- `WindowClose`イベントを`WindowClosing`に改名しました（[@leaanthony](https://github.com/leaanthony)）
- フロントエンドのビルドで、ビルド種別に応じて Vite 環境「development」または「production」を使用するようになりました（[@leaanthony](https://github.com/leaanthony)）
- go-webview2 v1.19に更新しました（[@leaanthony](https://github.com/leaanthony)）
- taskfile のフォークが確実に使用されるようにしました（@leaanthony）
- 次を使用してインストールする際のバージョン問題を修正するため、Taskfile のフォークを更新
- 次を使用してインストールする際のバージョン問題を修正するため、Taskfile のフォークを使用
- `service.OnStartup`はエラー発生時にアプリケーションを終了し、次を実行するようになりました
- ユーザー操作との整合性を高めるため、システムトレイのクリックメッセージ処理をリファクタリング
- 次を生成するフレームワークをサポートするため、埋め込みアセットに`all:frontend/dist`を含めるように変更
- Taskfile をリファクタリング（[leaanthony](https://github.com/leaanthony)、次にて）
- `go-webview2` v1.0.16へアップグレード
- `Screen`型を修正し、`Id`ではなく`ID`を含めるように変更
- `application.ServiceOptions`をサポートするため、`go.mod.tmpl`の Wails バージョンを更新
- サービス名の判定を修正（[windom](https://github.com/windom/)、次にて）
- mkdocs serve が Docker を使用するようになりました（[leaanthony](https://github.com/leaanthony)）
- 開発設定を`config.yml`に統合
- システムトレイのダイアログで、利用可能な場合はアプリケーションアイコンをデフォルトで使用するように変更（Windows）
- macOS の GPU とメモリのレポートを改善
- `WebviewGpuIsDisabled`と`EnableFraudulentWebsiteWarnings`を削除
- Events API の変更：`On`/`Emit` -> ユーザーイベント、`OnApplicationEvent` ->
- Linux の Events API を修正（[TheGB0077](https://github.com/TheGB0077)、次にて）
- [CI] Actions を改善し、フォークでも Actions を実行できるようにして
- `AbsolutePosition()`を`Position()`に改名
- Linux の WebKit 依存関係を webkitgtk2-4.0から webkit2gtk-4.1へ更新し、次を実現
- 同梱の JS ランタイムスクリプトが ESM モジュールになりました。これをインポートする script タグは
- `@wailsio/runtime`パッケージは、その API を`window.wails`上に公開しません
- Window API モジュール`@wailsio/runtime/src/window`が、それを含む次のものを公開するようになりました
- JS の Window API を、現在の Go `WebviewWindow`に合わせて更新しました
- バインディングジェネレーターは、デフォルトで ID による呼び出しを使用するようになりました。`-id` CLI オプションは
- 新しいバインディングコードのレイアウト：以前、出力ファイルはフォルダー内に整理されていました
- 構造体フィールド`application.Options.Bind`の名前を次のように変更しました
- サービスをバインドするための新しい構文：サービスインスタンスを次でラップすることが必須になりました
- 非ターミナル環境または CI 環境ではスピナーを無効化

### 削除

- **破壊的変更**：ウィンドウごとの`WindowsWindow`オプションから`EnabledFeatures`、`DisabledFeatures`、`AdditionalLaunchArgs`を削除しました。代わりに、アプリケーションレベルの`Options.Windows.EnabledFeatures`、`Options.Windows.DisabledFeatures`、`Options.Windows.AdditionalBrowserArgs`を使用してください。これらのフラグは、共有 WebView2 環境全体に適用されます（#4559、@leaanthony）
- Windows のネイティブ`IDropTarget`実装を削除し、JavaScript ベースの方式に変更しました（v2の動作と同じ）
- github.com/wailsapp/mimetype 依存関係を削除し、拡張した拡張子マップと標準ライブラリの http.DetectContentType に置き換えて、バイナリサイズを約1.2MB 削減しました
- Linux ファイルエクスプローラー向けの最小限の .desktop ファイルパーサーを実装して gopkg.in/ini.v1 依存関係を削除し、約45KB 削減しました
- Go 1.21以降の標準ライブラリの slices パッケージと最小限の内部ヘルパーを使用して、ランタイムコードから samber/lo を削除し、約310KB 削減しました
- Darwin の URL スキームハンドラーからデバッグ用の printf 文を削除しました（#4834）
- **破壊的変更**：`linux:WindowLoadChanged`イベントを削除しました。WebView の読み込み完了を検出するには、代わりに`linux:WindowLoadFinished`を使用してください（#3896、@leaanthony）

### 破壊的変更

- **Manager API のリファクタリング**：コードの整理と API の見つけやすさを向上させるため、アプリケーション API をフラットな構造から体系化されたマネージャー構造へ再編しました。[#4359](https://github.com/wailsapp/wails/pull/4359)で[@leaanthony](https://github.com/leaanthony)が実施しました
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
- Service のメソッド名を変更しました：`Name` -> `ServiceName`、`OnStartup` -> `ServiceStartup`、`OnShutdown` -> `ServiceShutdown`（[@leaanthony](https://github.com/leaanthony)）
- `Path`メソッドと`Paths`メソッドを`application`パッケージへ移動しました（[@leaanthony](https://github.com/leaanthony)）
- アプリケーションメニューは macOS 専用になりました（[@leaanthony](https://github.com/leaanthony)）

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

+ 過度に広範な Linux パッケージ依存関係と、古くなった RPM 依存関係を修正しました。

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

`go install`（@leaanthony）

## v3.0.0-alpha.8.1 - 2024-12-07

### 変更

`go install`（@leaanthony）

## v3.0.0-alpha.8 - 2024-12-06

### 追加

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### 変更

`service.OnShutdown`以前に開始されたすべてのサービス（@atterpac、[#3920](https://github.com/wailsapp/wails/pull/3920)）   @atterpac（[#3907](https://github.com/wailsapp/wails/pull/3907)）   サブフォルダー（@atterpac、[#3887](https://github.com/wailsapp/wails/pull/3887)）   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913)（[#3778](https://github.com/wailsapp/wails/pull/3778)）   [northes](https://github.com/northes)（[#3836](https://github.com/wailsapp/wails/pull/3836)）   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   （`EnabledFeatures`および`DisabledFeatures`オプションにより置き換え済み）（[leaanthony](https://github.com/leaanthony)）

### 修正

チャネル変数（@michael-freling、[#3925](https://github.com/wailsapp/wails/pull/3925)）   [ansxuman](https://github.com/ansxuman)（[#3924](https://github.com/wailsapp/wails/pull/3924)）   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony)（[#3841](https://github.com/wailsapp/wails/pull/3841)）   Darwin の編集メニューにおける`PasteAndMatchStyle`ロール（[johnmccabe](https://github.com/johnmccabe)、[#3839](https://github.com/wailsapp/wails/pull/3839)）   [#3840](https://github.com/wailsapp/wails/issues/3840)（[#3854](https://github.com/wailsapp/wails/pull/3854)、[kodflow](https://github.com/kodflow)）   [@leaanthony](https://github.com/leaanthony)   が異なる問題（@nickisworking、[#3789](https://github.com/wailsapp/wails/pull/3789)）

## v3.0.0-alpha.7 - 2024-09-18

### 追加

[mmghv](https://github.com/mmghv)（[#3665](https://github.com/wailsapp/wails/pull/3665)）   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac)および[leaanthony](https://github.com/leaanthony)（[#3570](https://github.com/wailsapp/wails/pull/3570)）

### 変更

Application Events `OnWindowEvent` -> Window Events（[leaanthony](https://github.com/leaanthony)、[#3734](https://github.com/wailsapp/wails/pull/3734)）   `v3/`または`v3-`を接頭辞とするブランチ（[stendler](https://github.com/stendler)、[#3747](https://github.com/wailsapp/wails/pull/3747)）

### 修正

[etesam913](https://github.com/etesam913)（[#3742](https://github.com/wailsapp/wails/pull/3742)）   [atterpac](https://github.com/atterpac)（[#3721](https://github.com/wailsapp/wails/pull/3721)）   [atterpac](https://github.com/atterpac)（[#3675](https://github.com/wailsapp/wails/pull/3675)）   [#1811](https://github.com/wailsapp/wails/pull/1811)（[#3614](https://github.com/wailsapp/wails/pull/3614)、[@stendler](https://github.com/stendler)）   [#3720](https://github.com/wailsapp/wails/pull/3720)（[leaanthony](https://github.com/leaanthony)）   [#3693](https://github.com/wailsapp/wails/issues/3693)（[@DeltaLaboratory](https://github.com/DeltaLaboratory)）   [#3720](https://github.com/wailsapp/wails/pull/3720)（[leaanthony](https://github.com/leaanthony)）   [#3693](https://github.com/wailsapp/wails/issues/3693)（[@DeltaLaboratory](https://github.com/DeltaLaboratory)）   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746)（[@stendler](https://github.com/stendler)）   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### 修正

## v3.0.0-alpha.5 - 2024-07-30

### 追加

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   アイコンのクリック（@5aaee9、[#2991](https://github.com/wailsapp/wails/pull/2991)）   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev)（[#3282](https://github.com/wailsapp/wails/pull/3282)）   @[Atterpac](https://github.com/Atterpac)（[#3022](https://github.com/wailsapp/wails/pull/3022])）   [@marcus-crane](https://github.com/marcus-crane)（[#3146](https://github.com/wailsapp/wails/pull/3146)）   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev)（[#3281](https://github.com/wailsapp/wails/pull/3281)）   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   @Mai-Lapyst による[PR](https://github.com/wailsapp/wails/pull/2044)に基づく   [@fbbdev](https://github.com/fbbdev)（[#3295](https://github.com/wailsapp/wails/pull/3295)）   [@fbbdev](https://github.com/fbbdev)（[#3295](https://github.com/wailsapp/wails/pull/3295)）   [@fbbdev](https://github.com/fbbdev)（[#3295](https://github.com/wailsapp/wails/pull/3295)）   npm パッケージ（[@fbbdev](https://github.com/fbbdev)、[#3334](https://github.com/wailsapp/wails/pull/3334)）   [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT`（[@abichinger](https://github.com/abichinger)、[#3429](https://github.com/wailsapp/wails/pull/3429)）   [@abichinger](https://github.com/abichinger)（[#3431](https://github.com/wailsapp/wails/pull/3431)）   [@bruxaodev](https://github.com/bruxaodev)（[#3667](https://github.com/wailsapp/wails/pull/3667)）   [@OlegGulevskyy](https://github.com/OlegGulevskyy)（[#3674](https://github.com/wailsapp/wails/pull/3674)）

### 修正

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane)による変更（[#3515](https://github.com/wailsapp/wails/pull/3515)）   [atterpac](https://github.com/atterac)による変更（[#3512](https://github.com/wailsapp/wails/pull/3512)）   [atterpac](https://github.com/atterpac)による変更（[#3477](https://github.com/wailsapp/wails/pull/3477)）   [Atterpac](https://github.com/atterpac)による変更（[#3320](https://github.com/wailsapp/wails/pull/3320)）。   [#3306](https://github.com/wailsapp/wails/pull/3306)で変更。   [#2972](https://github.com/wailsapp/wails/pull/2972)。   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv)による変更（[#2750](https://github.com/wailsapp/wails/pull/2750)）。   [#2753](https://github.com/wailsapp/wails/pull/2753)。   [jaybeecave](https://github.com/jaybeecave)による変更（[#3052](https://github.com/wailsapp/wails/pull/3052)）。   [@pylotlight](https://github.com/pylotlight)による変更（[PR](https://github.com/wailsapp/wails/pull/3039)）   をインストール。[@pylotlight](https://github.com/pylotlight)が[PR](https://github.com/wailsapp/wails/pull/3032)で追加。   [PR](https://github.com/wailsapp/wails/pull/3145)   スペース — @leaanthony。   [thomas-senechal](https://github.com/thomas-senechal)によるPR [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal)によるPR [#3208](https://github.com/wailsapp/wails/pull/3208)   アタッチされたウィンドウ。[tw1nk](https://github.com/tw1nk)によるPR [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev)による変更（[#3279](https://github.com/wailsapp/wails/pull/3279)）   `FRONTEND_DEVSERVER_URL`が存在します。   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev)による変更（[#3295](https://github.com/wailsapp/wails/pull/3295)）   リスナー。[@fbbdev](https://github.com/fbbdev)による変更（[#3295](https://github.com/wailsapp/wails/pull/3295)）   [@abichinger](https://github.com/abichinger)による変更（[#3330](https://github.com/wailsapp/wails/pull/3330)）   ジェネレーター。[@fbbdev](https://github.com/fbbdev)による変更（[#3334](https://github.com/wailsapp/wails/pull/3334)）   ジェネレーター。[@fbbdev](https://github.com/fbbdev)による変更（[#3334](https://github.com/wailsapp/wails/pull/3334)）   [@abichinger](https://github.com/abichinger)による変更（[#3346](https://github.com/wailsapp/wails/pull/3346)）   [@hfoxy](https://github.com/hfoxy)による変更（[#3417](https://github.com/wailsapp/wails/pull/3417)）   [@hfoxy](https://github.com/hfoxy)による変更（[#3426](https://github.com/wailsapp/wails/pull/3426)）   [@fbbdev](https://github.com/fbbdev)による変更（[#3431](https://github.com/wailsapp/wails/pull/3431)）   [@fbbdev](https://github.com/fbbdev)による変更（[#3431](https://github.com/wailsapp/wails/pull/3431)）   [@pekim](https://github.com/pekim)による変更（[#3458](https://github.com/wailsapp/wails/pull/3458)）   [@robin-samuel](https://github.com/robin-samuel)。   [@5aaee9](https://github.com/5aaee9)によるPR [#3466](https://github.com/wailsapp/wails/pull/3622)。   [@windom](https://github.com/windom/)による変更（[#3636](https://github.com/wailsapp/wails/pull/3636)）。   Windows。[@bruxaodev](https://github.com/bruxaodev/)による変更（[#3691](https://github.com/wailsapp/wails/pull/3691)）。

### 変更

[mmghv](https://github.com/mmghv)による変更（[#3611](https://github.com/wailsapp/wails/pull/3611)）   Ubuntu 24.04 LTSをサポート。[atterpac](https://github.com/atterpac)による変更（[#3461](https://github.com/wailsapp/wails/pull/3461)）   には`type="module"`属性が必要です。[@fbbdev](https://github.com/fbbdev)による変更（[#3295](https://github.com/wailsapp/wails/pull/3295)）   オブジェクトを作成しますが、WMLシステムは起動しません。これはカプセル化を改善するための変更です。必要に応じて、新しい`WML.Enable`メソッドを呼び出すことでWMLシステムを手動で起動できます。バンドルされたJSランタイムスクリプトは、引き続き両方の操作を自動的に実行します。[@fbbdev](https://github.com/fbbdev)による変更（[#3295](https://github.com/wailsapp/wails/pull/3295)）   windowオブジェクトをデフォルトエクスポートとして公開します。ESMの名前付きインポートまたは名前空間インポート構文を使用して個別のメソッドをインポートすることは、今後できません。   API。一部のメソッドでは、名前またはシグネチャが変更されました。具体的には、`Screen`は`GetScreen`に、`GetZoomLevel`/`SetZoomLevel`は`GetZoom`/`SetZoom`に変更されました。また、`GetZoom`、`Width`、`Height`は、値をオブジェクトでラップせず直接返すようになりました。[@fbbdev](https://github.com/fbbdev)による変更（[#3295](https://github.com/wailsapp/wails/pull/3295)）   は削除されました。名前による呼び出しに戻すには、`-names` CLIオプションを使用してください。[@fbbdev](https://github.com/fbbdev)による変更（[#3468](https://github.com/wailsapp/wails/pull/3468)）   以前はそれらを含むパッケージにちなんで命名されていましたが、現在はモジュールパスを含む完全なGoインポートパスが使用されます。[@fbbdev](https://github.com/fbbdev)による変更（[#3468](https://github.com/wailsapp/wails/pull/3468)）   `application.Options.Services`。[@fbbdev](https://github.com/fbbdev)による変更（[#3468](https://github.com/wailsapp/wails/pull/3468)）   `application.NewService`の呼び出し。[@fbbdev](https://github.com/fbbdev)による変更（[#3468](https://github.com/wailsapp/wails/pull/3468)）   [@DeltaLaboratory](https://github.com/DeltaLaboratory)による変更（[#3574](https://github.com/wailsapp/wails/pull/3574)）
