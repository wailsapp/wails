---
title: "変更履歴"
description: "Wails v3 のバージョン履歴とリリースノート"
slug: "changelog"
sourcePath: "changelog.md"
---

凡例：

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- このプロジェクトの注目すべき変更はすべて、このファイルに記録されます。

この形式は[Keep a Changelog](https://keepachangelog.com/en/1.0.0/)に基づいており、このプロジェクトは[セマンティック バージョニング](https://semver.org/spec/v2.0.0.html)に準拠しています。

- 新機能には`Added`を使用します。
- 既存機能の変更には`Changed`を使用します。
- 間もなく削除される機能には`Deprecated`を使用します。
- 削除済みの機能には`Removed`を使用します。
- バグ修正には`Fixed`を使用します。
- 脆弱性がある場合は`Security`を使用します。

_/

/_   * 自動更新 *   リリースノートはマージ済み PR から生成されます。通常のコントリビューション PR では、このファイルを編集したり、削除済みの`v3/UNRELEASED_CHANGELOG.md`を再作成したりしないでください。変更履歴のメンテナンスについては`AGENTS.md`を参照してください。 _/

## [未リリース]

## v3.2.0-nightly.20261102 - 2026-11-02

## 変更
- @lethusang による[PR](https://github.com/wailsapp/wails/pull/6077)で、DisableBounce の代わりに DisableOverscroll を使用

## v3.2.0-nightly.20261101 - 2026-11-01

## 追加
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6521)で、CGO を使用しない実験的な macOS および Linux バックエンドを追加
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6531)で、ネイティブ YAML ビルド、ポータブルなパッケージ化および署名を追加

## 修正
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6519)で、各種v3ナイトリービルドの問題を修正し、依存関係を整理
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6533)で、v3リリース検証とフィクスチャ分類を修復
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6538)で、クロスプラットフォームのリリース検証エラーを解消

## v3.1.0 - 2026-11-01

## 追加
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6521)で、CGO を使用しない実験的な macOS および Linux バックエンドを追加
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6531)で、ネイティブ YAML ビルド、ポータブルなパッケージ化および署名を追加

## 変更
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6391)で、破棄されるデスクトップメタデータのコピーを回避
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6393)で、fmt を使わずにコンテンツ長をフォーマット
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6404)で、候補を収集せずにインデックスルートを選択
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6408)で、シークできないファイルのコピーバッファを再利用
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6410)で、プレフィックス判定用 reader の割り当てを回避
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6414)で、既知の拡張子に対する MIME 判定用の読み取りを省略
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6415)で、キューの増大と縮小の繰り返しを防止
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6419)で、破棄されるリクエストコンテキストの割り当てを回避
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6420)で、無効なデバッグログの準備を省略
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6428)で、MIME 判定用バッファを遅延割り当て
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6443)で、ハイパーリンクのフォーマットに伴う割り当てを削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6449)で、箇条書きのフォーマットに伴う一時割り当てを排除
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6453)で、UUID の文字列変換を 1 回の割り当てに削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6454)で、ステートメントハンドルのフォーマットに伴う割り当てを削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6455)で、名前ベース UUID のダイジェストに伴う一時割り当てを排除
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6460)で、未完了のコマンド出力行のコピーを回避
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6472)で、パッケージのステージング中に実行可能ファイルをストリーミングコピー
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6473)で、未完了の出力行の再スキャンを回避
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6475)で、コンパイル済みのコマンドルーティングパターンを再利用
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6482)で、ランタイム検出中のメタデータ検索を回避
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6485)で、ローカライズ済みフォールバックページの割り当てを削減
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6490)で、プロセッサごとのパターン再コンパイルを回避

## 修正
- @fan711 による[PR](https://github.com/wailsapp/wails/pull/6017)で、システムトレイの SetTooltip を実装
- @TechHutTV による[PR](https://github.com/wailsapp/wails/pull/6072)で、設定ポータルからダークモードを読み取るように修正
- @APshenkin による[PR](https://github.com/wailsapp/wails/pull/6166)で、Window.SetZoom が1.0未満の値を取れるように修正
- @DevLumuz による[PR](https://github.com/wailsapp/wails/pull/6275)で、各 Go リスナーへ emit 順に配信
- @DevLumuz による[PR](https://github.com/wailsapp/wails/pull/6276)で、Emit 内のイベントフックを hookLock の保持中に読み取るように修正
- @AlbinoGeek による[PR](https://github.com/wailsapp/wails/pull/6284)で、Run の開始後に作成したトレイが破棄されないよう、トレイの Menu プロパティを常に通知
- @AlbinoGeek による[PR](https://github.com/wailsapp/wails/pull/6287)で、Run の前に送信されたメインスレッド処理を Run の開始まで保留
- @gavidroselj による[PR](https://github.com/wailsapp/wails/pull/6304)で、GitHub updater provider に正しいアセット URL を使用
- @strraaa による[PR](https://github.com/wailsapp/wails/pull/6326)で、WebView2 が Content-Length を省略した場合もリクエストボディを保持
- @julianstorer による[PR](https://github.com/wailsapp/wails/pull/6346)で、macOS および GTK の画面を DIP 空間（混在 DPI の境界）で報告
- @bercedev による[PR](https://github.com/wailsapp/wails/pull/6367)で、閉じられた通知をクリックとして報告しないように修正
- @xmasdev による[PR](https://github.com/wailsapp/wails/pull/6371)で、ターミナルを閉じたときに開発プロセスを停止
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6373)で、エビデンスチェックにおいて変更履歴の実行を squash merge コミットに対応付け
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6386)で、タスクバー再起動時のデバウンスを 1 秒間維持
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6387)で、フルスクリーン状態から開始する場合もサイズ制約を維持
- @taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/6388)で、新しい制約に合わせるためにフルスクリーンまたは最大化されたウィンドウのサイズを変更しないように修正
- liquid-glass の例で Run より前にダイアログを表示しないようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6389)）
- Run より前に要求されたダイアログが永久に待機せず、何も行わないようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6390)）
- デスクトップエントリにコマンドがない場合のフォールバックを追加した（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6392)）
- ファイル名を保持せずに現在のコンテンツを判別するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6394)）
- 大文字と小文字が混在するアセット拡張子を認識するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6395)）
- 使用済みのコンテンツ判別用バッファーを解放するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6396)）
- ringqueue から削除された参照を解放するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6397)）
- フォールバックを選択する前にヘッダーの削除を反映するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6398)）
- フラッシュ前にフォールバックレスポンスのヘッダーを確定するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6399)）
- 名前が完全に一致する index ファイルのみを検出するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6400)）
- シーク不能なリソースへの HEAD で、ボディーの生成と全体の読み込みを回避するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6401)）
- 短い読み込みが続く場合でも MIME 判別用プレフィックスを満たすようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6402)）
- シーク不能で空の正常レスポンスを確定するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6403)）
- 展開した URI ターゲットの重複を回避するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6405)）
- フィールド展開時にパスを保持するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6406)）
- XDG のデスクトップファイル検索ディレクトリを尊重するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6407)）
- 空のフラッシュ時にコンテンツ判別用ヘッダーを確定するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6409)）
- デスクトップエントリを隠しているディレクトリをスキップするようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6411)）
- 実行ファイルのパス内で見つかったターゲットを保持するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6412)）
- タイマー発火後にコールバックがキャプチャーした値を解放するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6413)）
- ldconfig でライブラリが誤って一致することを回避した（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6416)）
- ファイルシステム上のライブラリを対象とするフォールバック検索の一致条件を限定した（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6417)）
- キャッシュ無効化時のデータ競合を解消した（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6418)）
- 任意指定のロガー設定を反映するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6421)）
- 共有された開始 URL の変更と古いポートの使用を回避した（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6422)）
- リクエストがキャンセルされたら開発プロキシの再試行を停止するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6423)）
- バンドル済みアセットがない場合に 404 を返すようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6424)）
- バンドル済みランタイムへの HEAD リクエストを処理するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6425)）
- サービスのルーティング時にエスケープ済みのパスを保持するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6426)）
- 非常に大きな数値のプレリリース識別子を正しく並べ替えるようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6427)）
- ツールのダウンロードを保存する前に HTTP エラーを検出して拒否するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6429)）
- デスクトップの検出中も起動期限を遵守するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6430)）
- 開発プロキシの再試行間でリクエストボディーを保持するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6432)）
- すべてのリターンパスで HTTP レスポンスボディーを閉じるようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6434)）
- ライブラリの並列検索でキャンセルエラーを保持するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6435)）
- 失敗した GitHub タグのレスポンスをデコード前に拒否するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6436)）
- 角括弧を含むテンプレートパスをリテラルのまま保持するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6437)）
- wails3 を一度だけインストールし、インストーラーのエラーを保持するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6438)）
- 非表示のエントリとコメントブロック内のヘッダーを無視するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6439)）
- 安定版リリースの検索エラーを通知するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6440)）
- プレーンテキストのメッセージをフォーマットとして解釈せずに出力するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6441)）
- JSON の null を読み込んだ後の書き込みでパニックが発生しないようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6442)）
- ストア全体の読み込みでスナップショットを返すようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6444)）
- 行を途中まで反復処理した後に発生したエラーを伝播するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6446)）
- 展開先のパスをリテラルのまま保持するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6447)）
- クローズ時と再オープン時にプリペアドステートメントのハンドルを解放するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6448)）
- ステートメントの準備とクローズまたは再オープンが競合する場合を処理するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6450)）
- フォールバック用 index ページへの HEAD リクエストを正しく処理するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6451)）
- ステートメントハンドルに完全な16進数表記を必須とした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6452)）
- 並行処理時もデフォルトのインメモリーデータベースを維持するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6456)）
- 解析対象の HEAD 出力から stderr を分離するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6457)）
- 埋め込みバンドルからランタイムアセットを生成するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6458)）
- TODO、FIXME、TBD のプレースホルダーを検出するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6459)）
- リリースノート API のレスポンスを検証するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6461)）
- ステータス確認時にタスクのディレクトリを反映するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6462)）
- 同一ファイルへのコピーによるデータ損失を防止した（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6463)）
- キャッシュの識別情報にタスクのディレクトリを含めるようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6464)）
- 1回のチェックにつきポートを一度だけ検査するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6465)）
- キャッシュエントリが null としてデコードされた場合のパニックを回避した（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6466)）
- 整数の上限に制約されずに数値のプレリリースをインクリメントするようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6467)）
- 従来のチェックによってキャッシュミスが見えなくなることを防止した（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6468)）
- executor のベースディレクトリでデフォルトタスクを実行するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6469)）
- エクスプローラーのコマンド内にあるリテラルのパーセント記号を保持するようにした（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6470)）
- [PR](https://github.com/wailsapp/wails/pull/6471) で事前条件とタスクコマンドのディレクトリを一致（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6474) でシンボリックリンクを生成された通常ファイルに置換（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6476) でビルド失敗時に既存の DMG を保持（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6477) で引用符で囲まれた正規表現とパスのバックスラッシュを保持（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6478) で画像を割り当てる前に生成されたアイコンのサイズを検証（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6479) でコマンドチェーンとリダイレクトをシェル経由で実行（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6480) でビルド失敗時に以前の Linux パッケージを保持（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6481) でプラットフォーム固有のシェルフォールバック実行を使用（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6483) で先頭にパス区切りのスラッシュが重複するアセットを配信（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6484) でビルドメタデータを含む安定版を認識（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6486) でリソース生成の失敗を正確に報告（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6487) でアクティブなレポーターの置換時に競合を防止（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6488) で同時実行されるキャッシュファイルの保存を直列化（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6489) で複数のテンプレート等価性候補をサポート（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6491) で保存済みカテゴリーマップの null を正規化（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6492) で失敗した変数シェルコマンドの診断情報を表示（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6493) でシェルコマンドチェーンを暗黙的な Go キャッシュの対象外に変更（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6494) で Linux のカテゴリー永続化を直列化（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6495) で引用符またはエスケープされた空白を含む Go 出力をキャッシュ（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6496) でエンコード失敗時に Mac アイコンを保持（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6497) で無効な Go テスト回数フラグを破棄しないように変更（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6498) で新規ファイルへのコピー時にコピー元の権限を保持（@taliesin-ai）
- v3 のさまざまなナイトリービルドの問題を修正し、依存関係を整理（[PR](https://github.com/wailsapp/wails/pull/6519)、@taliesin-ai）
- v3 のリリース検証とフィクスチャ分類を修復（[PR](https://github.com/wailsapp/wails/pull/6533)、@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6538) でクロスプラットフォームのリリース検証失敗を解消（@taliesin-ai）

## v3.1.0-nightly.20261028 - 2026-10-28

## 変更
- [PR](https://github.com/wailsapp/wails/pull/6490) でコンパイル済み正規表現パターンを再利用し、変更履歴の生成を高速化（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6485) でローカライズされたフォールバックページの配信時の割り当てを削減（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6482) で不要なメタデータ検索を避け、AppImage ランタイムの検出を高速化（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6475) でコンパイル済みコマンドルーティングパターンを再利用し、wake のパフォーマンスを向上（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6473) で未完了の出力行の再走査を避け、パフォーマンスを向上（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6472) で MSIX パッケージのステージング中に実行ファイルをストリーミングコピーし、メモリ使用量を削減してビルドを高速化（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6460) で未完了のコマンド出力行のコピーを避け、パフォーマンスを向上（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6455) で名前ベースの UUID 生成時の割り当てを削減（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6454) で SQLite ステートメントハンドルの書式設定時の割り当てを削減（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6453) で UUID の文字列変換に必要な割り当てを 1 回に削減（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6449) で変更履歴の箇条書きエントリーの書式設定時の割り当てを削減（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6443) でターミナルのハイパーリンクの書式設定時の割り当てを削減（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6428) で MIME バッファを遅延スニッフィングし、アセット配信時の割り当てを削減（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6420) でデバッグログが無効な場合、アセットサーバーでのデバッグログの準備を省略（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6419) で破棄されるリクエストコンテキストを作成せず、webview のリクエスト処理時の割り当てを削減（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6414) で既知のファイル拡張子に対する MIME スニッフィングを省略し、アセット配信を高速化（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6410) でコンテンツのスニッフィングにリーダーを使用せず、アセット配信時の割り当てを削減（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6408) でシークできないアセットファイルの配信時の割り当てを削減（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6404) で候補の収集を避け、アセットサーバーのインデックスルート選択を高速化（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6393) で fmt を使わずにコンテンツ長を書式設定し、アセットサーバーのパフォーマンスを向上（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6391) でデスクトップファイルエクスプローラーのメタデータ読み取り時のメモリコピーを削減（@taliesin-ai）

## 修正
- [PR](https://github.com/wailsapp/wails/pull/6498) で新しいコピー先ファイルへのコピー時にコピー元ファイルの権限を保持（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6388) で最小または最大サイズの制約が変更された際、全画面表示および最大化されたウィンドウのサイズが変わらないように修正（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6500) で wailsapp/wails リポジトリのスケジュールされたスター履歴更新を公開（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6497) で wake が無効な Go テスト回数フラグを破棄しないように修正（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6496) でエンコード失敗時に Mac アイコンを保持（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6495) で引用符またはエスケープされた空白を含む Go 出力のビルドキャッシュミスを修正（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6494) で Linux における通知カテゴリーの永続化を直列化し、同時書き込みの競合を防止（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6493) でシェルコマンドチェーンを暗黙的な Go ビルドキャッシュの対象外に変更（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6491) で保存済みデータが null の場合に通知カテゴリーマップを読み込めない問題を修正（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6492) で失敗した変数シェルコマンドの診断情報を表示（@taliesin-ai）
- [PR](https://github.com/wailsapp/wails/pull/6489) で複数の等価性候補が存在する場合の wake のテンプレートマッチングを修正（@taliesin-ai）
- キャッシュ書き込みの破損を防ぐため、同時実行されるキャッシュファイルの保存を直列化（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6488)）
- アクティブなレポーターを置き換える際の競合を防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6487)）
- CLI でリソース生成の失敗を正確に報告（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6486)）
- CLI でビルドメタデータを含む安定版を認識（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6484)）
- リクエストパスの先頭にスラッシュが重複しているアセットも配信（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6483)）
- 各プラットフォームで wake のシェルフォールバック実行を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6481)）
- ビルド失敗時に以前の Linux パッケージを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6480)）
- コマンドチェーンとリダイレクトをシェル経由で実行（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6479)）
- wake のコマンド解析で、引用符で囲まれた正規表現とパスのバックスラッシュを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6477)）
- パッケージ化ビルドの失敗時に既存の DMG ファイルを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6476)）
- バインディング出力内のシンボリックリンクを生成済みの通常ファイルに置換（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6474)）
- Linux のファイルエクスプローラーコマンドでリテラルのパーセント記号を保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6470)）
- 前提条件とタスクコマンドのディレクトリを揃え、一貫して解決されるように変更（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6471)）
- デフォルトタスクをエグゼキューターのベースディレクトリで実行（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6469)）
- 従来のチェックによってキャッシュミスが隠されるのを防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6468)）
- 整数オーバーフローの上限なく増分できるよう、CLI のプレリリース番号付けを修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6467)）
- キャッシュエントリーが null としてデコードされた場合のパニックを回避（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6466)）
- CLI でチェックごとにポートを一度だけプローブ（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6465)）
- 古い wake 結果を避けるため、キャッシュ ID にタスクディレクトリを追加（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6464)）
- パッケージ化処理でファイルをそれ自身へコピーする際のデータ損失を防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6463)）
- ステータスチェックでタスクディレクトリを考慮（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6462)）
- 変更履歴エントリー内の TODO、FIXME、TBD プレースホルダーを検出（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6459)）
- CLI で埋め込みバンドルからランタイムアセットを生成（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6458)）
- git コマンドで stderr を解析対象の HEAD 出力から分離（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6457)）
- 同時アクセス時もデフォルトのインメモリ SQLite データベースを維持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6456)）
- フォールバック用インデックスページへの HEAD リクエストに正しいレスポンスを返却（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6451)）
- SQLite ステートメントの準備とデータベースのクローズまたは再オープンとの競合に対応（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6450)）
- SQLite データベースのクローズまたは再オープン時に、準備済みステートメントのハンドルを解放（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6448)）
- gosod 使用時に展開先パスをリテラルのまま保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6447)）
- 行を途中まで反復処理した後に発生した SQLite クエリエラーを伝播（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6446)）
- kvstore でストア全体を読み取る際にスナップショットを返却（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6444)）
- JSON の null 値を読み込んだ後に kvstore で書き込み時のパニックが発生するのを防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6442)）
- 書式指定子として解釈せず、プレーンな端末メッセージを出力（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6441)）
- アップデーターで最新の安定版を検索する際のエラーを通知（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6440)）
- 変更履歴のコメントブロック内にある非表示エントリーとヘッダーを無視（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6439)）
- アップデーターで wails3 を一度だけインストールし、インストーラーのエラーを通知（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6438)）
- gosod で角括弧を含むリテラルのテンプレートパスを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6437)）
- アップデーターで失敗した GitHub タグレスポンスをデコード前に拒否（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6436)）
- ライブラリの並列検索中にキャンセルエラーを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6435)）
- リソースリークを防ぐため、すべてのリターンパスで HTTP レスポンスボディをクローズ（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6434)）
- 開発プロキシがリクエストを再試行する際にリクエストボディを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6432)）
- デスクトップのファイルエクスプローラー検出時に起動期限を遵守（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6430)）
- AppImage ツールのダウンロードを保存する前に HTTP エラーを拒否（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6429)）
- semver 比較で、桁数が非常に大きい数値のプレリリース識別子を正しく順序付け（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6427)）
- アセットサーバーでサービスをルーティングする際にエスケープ済みパスを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6426)）
- アセットサーバーでバンドル済みランタイムへの HEAD リクエストに対応（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6425)）
- バンドル済みアセットが見つからない場合、不正なレスポンスではなく 404 を返却（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6424)）
- リクエストがキャンセルされた場合、開発プロキシの再試行を停止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6423)）
- アセットサーバーで開始 URL の共有状態の変更と古いポートの使用を回避（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6422)）
- アセットサーバーでオプションのロガー設定を適用（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6421)）
- ライブラリパスのキャッシュ無効化におけるデータ競合を解消（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6418)）
- ldconfig の出力を検索する際のライブラリの誤一致を回避（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6416)）
- 負荷時にアセットサーバーのキューが増大し、縮小と拡大を繰り返すのを防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6415)）
- 古い参照が保持されないよう、デバウンスタイマーの発火後にコールバックのキャプチャを解放（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6413)）
- 実行可能ファイルのパス内で見つかったファイルエクスプローラーの対象を保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6412)）
- ファイルエクスプローラーでデスクトップエントリーを隠す同名ディレクトリをスキップ（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6411)）
- レスポンスボディが空の場合にアセットサーバーのヘッダーが送信されない問題を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6409)）
- アプリケーションアイコンの検索時に XDG のデスクトップファイル検索ディレクトリを考慮（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6407)）
- ファイルエクスプローラーでフィールドを展開する際にファイルパスを保持 [PR](https://github.com/wailsapp/wails/pull/6406) by @taliesin-ai
- ファイルエクスプローラーで展開済みの URI ターゲットが重複しないように修正 [PR](https://github.com/wailsapp/wails/pull/6405) by @taliesin-ai
- シーク不可の空のアセットレスポンスが成功した場合、保留のままにせずコミットするように修正 [PR](https://github.com/wailsapp/wails/pull/6403) by @taliesin-ai
- アセットサーバーで短い読み取りが複数回発生しても、MIME スニッフィング用プレフィックスを正しく埋めるように修正 [PR](https://github.com/wailsapp/wails/pull/6402) by @taliesin-ai
- アセットサーバーでシーク不可の HEAD リクエストに対してレスポンス本文全体を読み込まないように修正 [PR](https://github.com/wailsapp/wails/pull/6401) by @taliesin-ai
- アセットサーバーでディレクトリリクエストを解決する際、完全一致するインデックスファイルのみを配信するように修正 [PR](https://github.com/wailsapp/wails/pull/6400) by @taliesin-ai
- アクションをサポートしないサーバーで、閉じられた Linux 通知がクリックとして報告される問題を修正 [PR](https://github.com/wailsapp/wails/pull/6367) by @bercedev
- ターミナルが閉じられたときに開発プロセスを停止するように変更 [PR](https://github.com/wailsapp/wails/pull/6371) by @xmasdev
- Windows タスクバー再起動のデバウンスフラグをクリアする前に 1 秒間保持するように変更 [PR](https://github.com/wailsapp/wails/pull/6386) by @taliesin-ai
- アセットサーバーでフォールバックレスポンスのヘッダーが失われる問題を修正 [PR](https://github.com/wailsapp/wails/pull/6399) by @taliesin-ai
- アセットサーバーでフォールバックを選択する前に、ヘッダーの削除を反映するように修正 [PR](https://github.com/wailsapp/wails/pull/6398) by @taliesin-ai
- メモリリークを防ぐため、アセットサーバーでリングキューから削除された参照を解放 [PR](https://github.com/wailsapp/wails/pull/6397) by @taliesin-ai
- メモリ使用量を削減するため、アセットサーバーで使用済みのコンテンツスニッファーバッファを解放 [PR](https://github.com/wailsapp/wails/pull/6396) by @taliesin-ai
- アセット配信時に大文字と小文字が混在するファイル拡張子を認識するように修正 [PR](https://github.com/wailsapp/wails/pull/6395) by @taliesin-ai
- ファイル名を保持せずにアセットサーバーでコンテンツタイプを正しく検出するように修正 [PR](https://github.com/wailsapp/wails/pull/6394) by @taliesin-ai
- デスクトップエントリにコマンドがない場合はデフォルトにフォールバックするように変更 [PR](https://github.com/wailsapp/wails/pull/6392) by @taliesin-ai
- macOS 以外で、アプリの起動前に liquid-glass サンプルがダイアログを表示しないように修正 [PR](https://github.com/wailsapp/wails/pull/6389) by @taliesin-ai
- Windows で制約付きのサイズからフルスクリーンを開始した場合も、ウィンドウサイズの制約を維持するように修正 [PR](https://github.com/wailsapp/wails/pull/6387) by @taliesin-ai
- Run より前に呼び出されたダイアログが制御を返さず永久にブロックする問題を修正 [PR](https://github.com/wailsapp/wails/pull/6390) by @taliesin-ai

## セキュリティ
- 過剰なメモリ使用を防ぐため、画像を割り当てる前に生成アイコンのサイズを検証 [PR](https://github.com/wailsapp/wails/pull/6478) by @taliesin-ai
- 使用前にリリースノート API のレスポンスを検証 [PR](https://github.com/wailsapp/wails/pull/6461) by @taliesin-ai
- SQLite プラグインで完全な 16 進数形式のステートメントハンドルを必須化 [PR](https://github.com/wailsapp/wails/pull/6452) by @taliesin-ai
- 意図しないライブラリの読み込みを防ぐため、ファイルシステム上のライブラリに対するフォールバック検索の一致範囲を制限 [PR](https://github.com/wailsapp/wails/pull/6417) by @taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## 変更
- 404 ページから、v2.wails.io 上の対応する Wails v2 ドキュメントページへのリンクを追加 [PR](https://github.com/wailsapp/wails/pull/6366) by @taliesin-ai
- v2 ドキュメントサイトを v2.wails.io に移行 [PR](https://github.com/wailsapp/wails/pull/6365) by @taliesin-ai
- v3 ドキュメントを wails.io から、v2 ドキュメントを v2.wails.io から配信するように変更 [PR](https://github.com/wailsapp/wails/pull/6364) by @taliesin-ai
- v3.0.0 のリリースノートを更新し、セキュリティアドバイザリに関する記述を訂正 [PR](https://github.com/wailsapp/wails/pull/6363) by @taliesin-ai

## 修正
- Go のイベントリスナーコールバックを emit 順に直列で実行するように修正 [PR](https://github.com/wailsapp/wails/pull/6275) by @DevLumuz
- Run の後に作成されたトレイが Linux で破棄されないよう、トレイの Menu プロパティを常に通知するように修正 [PR](https://github.com/wailsapp/wails/pull/6284) by @AlbinoGeek
- Run より前に投入されたメインスレッドの処理を、アプリケーションが起動するまで保留するように修正 [PR](https://github.com/wailsapp/wails/pull/6287) by @AlbinoGeek
- GitHub アップデーターがリリースアセットをダウンロードする際に、正しいアセット API URL を使用するように修正 [PR](https://github.com/wailsapp/wails/pull/6304) by @gavidroselj
- DPI が混在する macOS および Linux 環境で、画面の境界と位置を修正 [PR](https://github.com/wailsapp/wails/pull/6346) by @julianstorer
- Windows で WebView2 が Content-Length ヘッダーを省略した場合も、リクエスト本文を保持するように修正 [PR](https://github.com/wailsapp/wails/pull/6326) by @strraaa
- ロック中にフックのスナップショットを取得し、イベントフックの発行時に発生するデータ競合を修正 [PR](https://github.com/wailsapp/wails/pull/6276) by @DevLumuz
- Linux で `SystemTray.SetTooltip()` が何も行わない問題を修正：StatusNotifierItem のツールチップ説明が設定されるようになりました（#6017）
- 週次リリースワークフローで変更履歴のエビデンスチェックを実行するように変更 [PR](https://github.com/wailsapp/wails/pull/6374) by @taliesin-ai
- 実行結果を squash merge コミットに対応付けるよう、リリース変更履歴のエビデンスチェックを修正 [PR](https://github.com/wailsapp/wails/pull/6373) by @taliesin-ai
- Linux のダークモード検出について、ダークデスクトップでもカラースキームが変更されるまでライトと判定される問題と、GTK4 バックエンドで `linux:SystemThemeChanged` が一度も発火しない問題を修正しました。両バックエンドで、`org.freedesktop.appearance` ポータル名前空間を基盤とする単一の実装を共有するようになりました（#6072）
- フレーム付きウィンドウ（Windows、macOS、Linux）では、`Window.SetZoom` が 1.0 未満の係数を 1.0 に引き上げなくなりました。WebView2 の `ZoomFactor`、WKWebView の `magnification`、WebKitGTK のズームレベルはいずれも任意の正の値を受け入れ、`Zoom` ウィンドウオプションもすでにそのような値をそのまま渡していました。フレームレスウィンドウでは、ドラッグ領域とリサイズ領域がズームに合わせて調整されないため、下限 1.0 を維持します（[#4590](https://github.com/wailsapp/wails/issues/4590)）。正でない値、NaN、Inf の係数は、1.0 に引き上げるのではなく無視するようになりました — by @APshenkin

## v3.0.0 - 2026-10-26

## 変更
- ナイトリープレリリースは次回の週次リリースのプレビューとして `vX.Y.0-nightly.YYYYMMDD` の形式で公開され、週次リリースではマイナーバージョンが上がります（`vX.Y.0`）。そのため、Go の `@latest` は常に週次リリースまたはホットフィックスリリースを解決します。`@wailsio/runtime` npm パッケージのバージョンは常に Wails のリリースバージョンと一致し、ナイトリーには `nightly` dist-tag が付きます
- プロジェクトテンプレートでは、`@wailsio/runtime` をテンプレートの生成に使用した CLI のバージョンに固定
- Wails v3 には Go 1.26 以降が必要です（Go 1.26.10、または 1.27.3 以降には `crypto/x509` のセキュリティ修正が含まれます）

## 修正
- `wails3 update cli` が v3 リリースを検出できない問題を修正：v2 タグのみを対象としており、誤ったパスからインストールしていました
- 一部の Windows ビルドでアプリケーションのメニューバーを描画する際、Windows アプリが `panic: SetTextColor failed` でクラッシュする問題を修正しました。メニューバーは有効なデバイスコンテキストに描画され、使用できない場合はデフォルトの描画にフォールバックします。`UAHMENUITEM.Position` は C の `int` と一致するようになり、メニュー描画で panic が発生することもなくなりました（#6311）
- 再起動の保留中もステージング済みの更新を安全に維持するように修正：シャットダウン中に開始されたダウンロードや 2 回目の再起動によって、スワップヘルパーがインストールを待っているペイロードが削除されたり、競合したりすることはありません by @julianstorer
- `pkg/w32` の初期化処理が `uxtheme.dll` をアンロードした後も、その序数関数ポインターが使用されていたため、Windows アプリが `applyDarkMode` でクラッシュする問題を修正しました。モジュールはプロセスの存続期間中、ロードされたままになりました（#6356）
- スワップ前にステージング済みペイロードを検証し、バックアップを破棄する前にインストール済みアプリケーションを検証することで、破損した更新によって正常に動作するアプリケーションが置き換えられず、ロールバックされるようにしました（@julianstorer）

## 削除
- 最終翻訳同期後、[PR](https://github.com/wailsapp/wails/pull/5736) で Crowdin 連携を削除（@github-actions[bot]）

## セキュリティ
- すべてのフレームにガードスクリプトを挿入し、Linux（GTK3 および GTK4）上のクロスオリジン子フレームからの予約済み `wails:` メッセージをブロックすることで、GHSA-2p5f-2gx8-mh8m を緩和
- `github.com/klauspost/compress` を 1.18.7（GHSA-259r-337f-4rfw）に、`golang.org/x/image` を 0.45.0 に更新し、`crypto/x509` の修正（GO-2026-6809）を取り込むため Go 1.26.10 でビルド
- ランタイムパッケージ、セットアップウィザード、サンプル、プロジェクトテンプレートの `vite` を 8.2.4 以降に更新（GHSA-9jrq-w75r-8gcw、GHSA-rq7h-c2jc-7f22、GHSA-vfpm-58rq-9qcg、GHSA-r825-w8mm-5pxw、GHSA-587h-7wmj-h5pm、GHSA-698v-4r3x-3vcv）

### 既知の制限事項
- **Linux：** WebKitGTK はハンドラーにフレーム情報を提供しないため、子フレームガードでは、*同一オリジン*の子フレーム（`srcdoc` または `about:blank` iframe など）が親を介してネイティブメッセージハンドラーに到達するのを防げません。クロスオリジンフレーム、および `allow-same-origin` のないサンドボックス化されたフレームはブロックされます。信頼できない同一オリジンのコンテンツを Wails ウィンドウに埋め込まないでください。
- **Android（実験的）：** `window.wails` JavaScript ブリッジは WebView 内のすべてのフレームに公開され、メインフレームやオリジンのチェックは行われません。クロスオリジンコンテンツを含め、埋め込まれた iframe 内のコンテンツは、バインドされたメソッドを呼び出し、予約済み `wails:` メッセージを送信できます。信頼できないサードパーティコンテンツをアプリの WebView に埋め込まないでください。これは、子フレームからの `wails:` メッセージを無視する macOS および iOS とは異なります。修正は 3.0.1 で予定されています。
- **すべてのプラットフォーム：** 子フレームからのバインド済みメソッド呼び出しは、ネイティブメッセージハンドラーではなく、HTTP 経由で `/wails/runtime` に送られるため、ガードによる制限を受けません。埋め込む iframe はすべて、バインドされた Go メソッドにアクセスできるものとして扱ってください。

## v3.0.0-beta.34 - 2026-10-20

## 修正
- Windows でウィンドウをアクティブ化しても webview がフォーカスを取り戻さない問題と、Tab キーによるフォーカス移動が webview 内に閉じ込められ、ネイティブウィンドウとの間を移動できない問題を修正（#3783）
- `internal/webview2/pkg/webview2` の COM vtable 構造体（`ICoreWebView2_2`..`_27` およびその他の派生インターフェース）で継承されたメソッドチェーンが欠落し、呼び出しが誤った vtable スロットに到達していた問題を修正（#6129）
- Windows の質問ダイアログでカスタムボタンラベルが無視され、ボタンのコールバックが一度も実行されない問題を修正。カスタムボタンは TaskDialog で表示されるようになり、comctl32 v6 が利用できない場合は、位置に基づくコールバックマッピングを使用する MessageBox にフォールバックします（#5794）
- Vite に対して 60 秒間 `wails3 dev` がハングする問題を修正。自動フロントエンド準備完了ゲートは、`Accept` ヘッダーを送信せず 404 を受信していた HTTP GET の代わりに TCP プローブを使用するようになりました（[#6209](https://github.com/wailsapp/wails/issues/6209)）
- 開いている間に更新してもシステムトレイメニューが動作し続けるように修正（[PR](https://github.com/wailsapp/wails/pull/6301)、@taliesin-ai）
- ランタイム操作で一時的な WebView2 COM エラーを適切に処理（[PR](https://github.com/wailsapp/wails/pull/6302)、@leaanthony）

## v3.0.0-beta.33 - 2026-10-19

## 追加
- バインディングジェネレーターが、コンパイル時定数の名前を使用して `application.RegisterEvent` に登録された検出可能な各イベントについて、型付きの名前定数をエクスポートする `bindings/events` モジュールを生成するようになりました（例：`user:login` に対する `UserLogin`）。Vite プラグインなしで `Events.On`/`Emit` とともに使用できます（#5645）

## 修正
- Windows の質問ダイアログでカスタムボタンラベルが無視され、ボタンのコールバックが一度も実行されない問題を修正。カスタムボタンは TaskDialog で表示されるようになり、comctl32 v6 が利用できない場合は、位置に基づくコールバックマッピングを使用する MessageBox にフォールバックします（#5794）
- Vite に対して 60 秒間 `wails3 dev` がハングする問題を修正。自動フロントエンド準備完了ゲートは、`Accept` ヘッダーを送信せず 404 を受信していた HTTP GET の代わりに TCP プローブを使用するようになりました（[#6209](https://github.com/wailsapp/wails/issues/6209)）
- 開いている間に更新してもシステムトレイメニューが動作し続けるように修正（[PR](https://github.com/wailsapp/wails/pull/6301)、@taliesin-ai）
- ランタイム操作で一時的な WebView2 COM エラーを適切に処理（[PR](https://github.com/wailsapp/wails/pull/6302)、@leaanthony）

## v3.0.0-beta.32 - 2026-10-18

## 追加
- macOS 14 以降向けのネイティブ macOS メニューセクション見出しを追加。[ドキュメントサイト](https://v3.wails.io)を参照（20 ページを更新）（[PR](https://github.com/wailsapp/wails/pull/6297)、@leaanthony）

## 修正
- macOS で SetRelativePosition が画面の作業領域を基準に動作するように修正（[PR](https://github.com/wailsapp/wails/pull/6299)、@leaanthony）
- 親プロセスのタイムアウトまたは置換失敗後に、アップデーターヘルパーが変更されていない対象を再起動するように修正（[PR](https://github.com/wailsapp/wails/pull/6298)、@taliesin-ai）
- バージョン別ドキュメントで pnpm をバージョン 10 に固定（[PR](https://github.com/wailsapp/wails/pull/6293)、@taliesin-ai）

## v3.0.0-beta.31 - 2026-10-18

## 追加
- ルートのランタイムエントリを評価せずに `CancellablePromise` をインポートできるよう、`@wailsio/runtime/cancellable` サブパスエクスポートを公開（#6153）

## 変更
- Go 1.26 以降を必須化。v3 モジュールは `go 1.26.0` を宣言するようになり、`wails3 init` および `wails3 service init` で生成されたプロジェクトも同じ最低バージョンを宣言します

## 修正
- Go 1.27 でビルドした際にバインディングジェネレーターの出力が変化する問題を修正。`Marshaler`、`Unmarshaler`、`RawMessage` などの `encoding/json` 型が `encoding/json/v2` 型および `encoding/json/jsontext` 型のエイリアスになっていましたが、再び `any` としてレンダリングされ、`encoding/json/v2` のモデルは生成されなくなりました（#6292）
- `wails3 init -t` に渡されたリモートテンプレート URL をクローン前に検証。`https://`、`http://`、`ssh://`、`git://`、`user@host:path` の URL のみを受け入れ、値は `--` の後に git へ渡し、クローン中に git が使用できるトランスポートをこれらに制限します
- HTTP ランタイムトランスポートで、範囲外の `x-wails-window-id` ヘッダー値を別のウィンドウ ID に暗黙的に切り詰めるのではなく拒否するように変更
- zip-slip に対するアップデーターのアーカイブ展開処理を強化。エントリパスを `filepath.IsLocal` でチェックしたうえで、結合後のパスが展開先に含まれることも確認するようになりました。また、zip および tar.gz アーカイブ内のネストされた `..`、バックスラッシュ、絶対パス、シンボリックリンク経由の書き込みを行うエントリに対するテストを追加しました
- `wails3 generate appimage` で作成した AppImage が Debian および Ubuntu 以外で起動時に「Unable to spawn a new child process: … WebKitNetworkProcess」とクラッシュする問題を修正。バンドルされた WebKit は、AppImage に同梱されたヘルパープロセスを実行するようになりました。WebKit の bubblewrap サンドボックスは引き続き有効です
- `-tags server` ビルドが Android および iOS 向けにコンパイルできない問題を修正し、[#6271](https://github.com/wailsapp/wails/issues/6271) を解決
- 一時ディレクトリが別のファイルシステム上にある場合（tmpfs の `/tmp` など）、Linux でアップデーターがスワップに失敗する問題を修正。インストール済みバイナリの隣にステージングし、ファイルシステム間ではコピーするようになりました
- アップデーターが `.AppImage` ファイルではなく、AppImage の読み取り専用マウントを置き換えようとする問題を修正
- macOS および iOS で、予約済みの `wails:` ウィンドウ制御メッセージが iframe から受け入れられていた問題を修正しました。今後はウィンドウのメインフレームからのみ受け入れられます。iframe からのその他のメッセージは、送信元情報とともに引き続き `RawMessageHandler` に届きます
- macOS での `WindowDidUpdate` および `WindowWillUpdate` の送信を停止しました。`NSWindowDidUpdateNotification` はランループを通過するたびに発火するため、ページでは毎秒約 120 回、空の `evaluateJavaScript:` 呼び出しが発生していました。また macOS の 27.0 では、そのたびにユーザーのクリックによる一時的なアクティベーションが消費され、実際にクリックしても `navigator.clipboard.writeText()` などのアクティベーションを必要とする API が `NotAllowedError` で拒否されていました（[#6170](https://github.com/wailsapp/wails/issues/6170)）。@APshenkin
- Explorer にバージョン情報が表示されるよう、Windows の info.json のデフォルト言語を en-US（0409）に修正しました（#6210）
- AMD または Intel の iGPU と NVIDIA dGPU を搭載したノート PC で、Linux の DMA-BUF レンダラーが無効になっていた問題を修正しました。Wayland では、WebKit が NVIDIA でレンダリングする場合（アプリが dGPU にオフロードされているか、コンポジターがそこで動作している場合）にのみ無効になり、iGPU でのレンダリングが極端に遅くなる原因となっていた nvidia モジュールのロードだけでは無効になりません。X11 セッションでは、Xorg のレンダリング GPU がその設定によって決まるため、従来の動作を維持します
- [PR](https://github.com/wailsapp/wails/pull/6135) で、Linux アップデーターにおけるデバイス間の名前変更失敗を、ステージングコピーへのフォールバックで処理するようにしました。@john-okeefe
- Windows の `WebviewWindow.SetMenu` で、メニューバーのテーマを設定していないとメニューバーが再描画されない問題を修正し、[#6103](https://github.com/wailsapp/wails/issues/6103) を解決しました
- Windows の `WebviewWindow.SetMenu(nil)` が macOS のように何もしないのではなくパニックする問題を修正し、[#6104](https://github.com/wailsapp/wails/issues/6104) を解決しました
- [PR](https://github.com/wailsapp/wails/pull/6094) で、競合するプッシュが発生しても再試行によって変更履歴のエントリを保持するようにしました。@leaanthony
- モジュールに一度も同梱されていないバイナリを参照していた埋め込みを削除し、すべてのプラットフォームで `go mod vendor` が `pattern arm64/WebView2Loader.dll: no matching files found` により失敗する問題を修正しました。[PR](https://github.com/wailsapp/wails/pull/6031) で [#5782](https://github.com/wailsapp/wails/issues/5782) および [#5376](https://github.com/wailsapp/wails/issues/5376) を解決しました。@Grantmartin2002
- [PR](https://github.com/wailsapp/wails/pull/6159) で、macOS におけるウィンドウの重複した破棄処理を防止しました。@charlesbalderson
- [PR](https://github.com/wailsapp/wails/pull/6162) で、右クリック時にシステムトレイメニューが開くよう WM_CONTEXTMENU を処理しました。@wzk2025
- ディスパッチされたコールバックがパニックすると、5 種類すべての `InvokeSync*` が呼び出し元を永久にブロックする問題を修正しました。[PR](https://github.com/wailsapp/wails/pull/6108) で [#6107](https://github.com/wailsapp/wails/issues/6107) を解決しました。@RALIST
- [PR](https://github.com/wailsapp/wails/pull/6236) で、破棄された GTK3 サブメニューを安全に再構築するようにしました。@leaanthony
- [PR](https://github.com/wailsapp/wails/pull/6261) で、ネストした Worker 用の Windows Worker 転送 shim を実装しました。@taliesin-ai
- Windows の `OpenFile()` で、`CanChooseDirectories(true)` の場合に `PromptForMultipleSelection()` が無視される問題を修正しました。macOS や Linux と同様に、複数のフォルダーを一度に選択できるようになりました（#6263）
- windows/arm64 で、WebView2 COM メソッドに値渡しされる `double` 引数の問題を修正しました。`ZoomFactor`、`RasterizationScale`、Cookie の `Expires`、および印刷設定の倍率、用紙サイズ、余白が影響を受けていました（以前は印刷設定の setter が `ErrDoubleArgUnsupported` を返していました）。影響を受ける呼び出しは、Go の syscall 経路ではロードできない `d0` をロードする `runtime.cgocall` トランポリンを経由するようになりました（golang.org/issue/62583）。@APshenkin

## セキュリティ
- サーバーモード（`-tags server`）では、`/wails/runtime` 呼び出しを、サーバー自身のオリジンまたは `ServerOptions.WebSocketOriginPatterns` に列挙されたオリジンからの `POST` リクエストとしてのみ受け付けるようになりました。その他のメソッドには HTTP 405、その他のブラウザーオリジンには HTTP 403 を返します。`Host` ヘッダーを書き換えるリバースプロキシの背後にデプロイする場合は、公開オリジンを `WebSocketOriginPatterns` に追加してください
- ループバックアドレス（デフォルトの `localhost`、`127.0.0.0/8`、または `::1`）にバインドされたサーバーモード（`-tags server`）では、`Host` ヘッダーが `localhost`、IP リテラル、バインド先ホスト、または `ServerOptions.WebSocketOriginPatterns` に一致するホストでない限り、リクエストを HTTP 403 で拒否するようになりました。これにより、`/wails/runtime` および WebSocket エンドポイントに対する DNS リバインディング攻撃を防ぎます
- Windows の通知で WinRT トースト API が失敗した場合に PowerShell へフォールバックしなくなり、代わりにエラーを返すようになりました

## v3.0.0-beta.30 - 2026-10-15

## 修正
- アップデーターが `.AppImage` ファイルではなく、AppImage の読み取り専用マウントを置き換えようとする問題を修正しました
- macOS での `WindowDidUpdate` および `WindowWillUpdate` の送信を停止しました。`NSWindowDidUpdateNotification` はランループを通過するたびに発火するため、ページでは毎秒約 120 回、空の `evaluateJavaScript:` 呼び出しが発生していました。また macOS の 27.0 では、そのたびにユーザーのクリックによる一時的なアクティベーションが消費され、実際にクリックしても `navigator.clipboard.writeText()` などのアクティベーションを必要とする API が `NotAllowedError` で拒否されていました（[#6170](https://github.com/wailsapp/wails/issues/6170)）。@APshenkin
- Explorer にバージョン情報が表示されるよう、Windows の info.json のデフォルト言語を en-US（0409）に修正しました（#6210）
- [PR](https://github.com/wailsapp/wails/pull/6135) で、Linux アップデーターにおけるデバイス間の名前変更失敗を、ステージングコピーへのフォールバックで処理するようにしました。@john-okeefe
- Windows の `WebviewWindow.SetMenu` で、メニューバーのテーマを設定していないとメニューバーが再描画されない問題を修正し、[#6103](https://github.com/wailsapp/wails/issues/6103) を解決しました
- Windows の `WebviewWindow.SetMenu(nil)` が macOS のように何もしないのではなくパニックする問題を修正し、[#6104](https://github.com/wailsapp/wails/issues/6104) を解決しました
- [PR](https://github.com/wailsapp/wails/pull/6094) で、競合するプッシュが発生しても再試行によって変更履歴のエントリを保持するようにしました。@leaanthony
- モジュールに一度も同梱されていないバイナリを参照していた埋め込みを削除し、すべてのプラットフォームで `go mod vendor` が `pattern arm64/WebView2Loader.dll: no matching files found` により失敗する問題を修正しました。[PR](https://github.com/wailsapp/wails/pull/6031) で [#5782](https://github.com/wailsapp/wails/issues/5782) および [#5376](https://github.com/wailsapp/wails/issues/5376) を解決しました。@Grantmartin2002
- [PR](https://github.com/wailsapp/wails/pull/6159) で、macOS におけるウィンドウの重複した破棄処理を防止しました。@charlesbalderson
- [PR](https://github.com/wailsapp/wails/pull/6162) で、右クリック時にシステムトレイメニューが開くよう WM_CONTEXTMENU を処理しました。@wzk2025
- ディスパッチされたコールバックがパニックすると、5 種類すべての `InvokeSync*` が呼び出し元を永久にブロックする問題を修正しました。[PR](https://github.com/wailsapp/wails/pull/6108) で [#6107](https://github.com/wailsapp/wails/issues/6107) を解決しました。@RALIST

## v3.0.0-beta.29 - 2026-10-13

## 追加
- macOS のグローバルショートカットがキーボードレイアウトを考慮するようになりました。非 QWERTY レイアウト（AZERTY、QWERTZ、Dvorak など）では、文字アクセラレーターが固定された QWERTY の物理位置ではなく、`UCKeyTranslate` を介してキーに表示された文字へバインドされます。

## 修正
- [PR](https://github.com/wailsapp/wails/pull/6233) で、Linux 上でネイティブメニューを設定する前に初期化するようにしました。@abhicodes-007
- [PR](https://github.com/wailsapp/wails/pull/6237) で、出力処理およびセットアップウィザードに関する macOS の GA 対応上のリグレッションを解消しました。@taliesin-ai
- Windows の WebView2 COM `AddRef`/`Release` ラッパーが、スレッドに残っていた古い最終エラーを失敗として報告する問題を修正しました。今後は参照カウントを返すため、release によってアプリケーションが中止されたり、アセットリクエストが失敗したりすることはありません（#3931）
- `ICoreWebView2Deferral.Release` が `AddRef` スロットを呼び出し、Windows 上で遅延されたすべてのアセットリクエストの deferral が解放されない問題を修正しました（#3931）

## v3.0.0-beta.28 - 2026-10-05

## 修正
- [PR](https://github.com/wailsapp/wails/pull/6225) で、macOS アプリの出力先をファイルプロバイダーの外部に維持するようにしました。@taliesin-ai
- [PR](https://github.com/wailsapp/wails/pull/6224) で、セットアップウィザードを完成させ、システムの詳細を表示するようにしました。@taliesin-ai
- [PR](https://github.com/wailsapp/wails/pull/6222) で、Windows における WebView2 の復旧時のエラー処理とリソースのクリーンアップを改善しました。@leaanthony
- Windows での WebView2 の復旧処理を改善し、アプリのクラッシュとコントローラーの取り残しを防止（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6220)）
- macOS のメディア認可に関するガイダンスを修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6221)）
- Linux のシステムトレイアイコンを右クリックした際、メニューが開くだけでなくクリックハンドラーも実行され、関連付けられたウィンドウが右クリックのたびに切り替わる問題を修正（#6018）
- Windows でネイティブ側から解放されるまで WebView2 コールバックハンドラーを保持（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6184)）

## v3.0.0-beta.27 - 2026-10-01

## 修正
- macOS 12 以降で `Permissions` オプションが無視される問題を修正しました。`WKUIDelegate` メディアキャプチャメソッドが実装され、Linux や Windows と同様に、カメラおよびマイクのリクエストに `PermissionAllow` と `PermissionDeny` が適用されるようになりました。`NSCameraUsageDescription` / `NSMicrophoneUsageDescription` は引き続き必要であり、サンドボックス化されている場合は、対応するデバイスエンタイトルメントも必要です

## v3.0.0-beta.26 - 2026-09-25

## 追加
- Watcher がフロントエンド開発サーバーに HTTP 準備状況チェックを適用（@atterpac による [PR](https://github.com/wailsapp/wails/pull/6150)）

## 修正
- WebView2 プロセスの障害発生時に空白のウィンドウを残さず復旧（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6002)）

## v3.0.0-beta.25 - 2026-09-22

## 修正
- メインアプリケーションの終了時に `wails3 dev` とそのバックグラウンドプロセスを停止（#6048）

## v3.0.0-beta.24 - 2026-09-20

## セキュリティ
- 最小権限の原則に従うよう、ビルドおよび公開ワークフローの GITHUB_TOKEN 権限を制限
- screen サンプルのアセットミドルウェアにおけるパストラバーサルとシンボリックリンクによる脱出を修正
- 承認済みのパッケージマネージャーコマンドのみに制限しながら、セットアップウィザードの依存関係インストール機能を維持

## v3.0.0-beta.23 - 2026-09-16

## 追加

- Linux で GTK アプリケーション ID を上書き可能に変更（@overlordtm による [PR](https://github.com/wailsapp/wails/pull/5972)）

## 修正

- Linux のシングルインスタンス機能が、アプリ固有の ID でバス名を取得するように修正（@overlordtm による [PR](https://github.com/wailsapp/wails/pull/5973)）
- Linux で中止されたカスタムスキームリクエストをキャンセル（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6126)）

## v3.0.0-beta.22 - 2026-09-14

## 追加

- 翻訳を完成させ、コミュニティメディアを移行。[ドキュメントサイト](https://v3.wails.io)を参照（1812 ページを更新）（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6122)）
- さまざまなスタイルとサイズの埋め込み可能な Wails ボタンを文書化（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6120)）

## v3.0.0-beta.21 - 2026-09-13

## 追加

- Wails v3 ドキュメントを M-Press で配信（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6116)）

## 修正

- 変更履歴生成のため、MPD フロントマターから JSON の slug 値を解析（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6118)）
- バックアップ失敗後、Updater がヘルパー環境変数を消去し、元のターゲットを再起動（@cnmax による [PR](https://github.com/wailsapp/wails/pull/6080)）
- App.Run の実行中にデフォルトのシグナルハンドラーを起動（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6098)）
- Windows のメニューで nil メニューを処理し、置き換えられたリソースを解放してメニューバーを再描画（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6112)）
- 共有 YAML 設定を使用する新規プロジェクトで MSIX パッケージングを復元（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6115)）
- ジェネリックモデルの生成関数が後方で宣言されたヘルパーを参照すると、生成された JavaScript および TypeScript バインディングの読み込みに失敗する問題を修正。また、相互依存するジェネリックモデルの作成時にスタックオーバーフローが発生しないように修正（#6062）

## v3.0.0-beta.20 - 2026-09-10

## 変更

- Clave のショーケースリンクを現在のウェブサイトとリポジトリに更新（@01xR4in による [PR](https://github.com/wailsapp/wails/pull/6082)）

## 修正

- ナビゲーション間で keepalive ハンドラーを維持しながら、ワーカーリクエストを含む、中止された Windows のアセットリクエストをキャンセル。Apple プラットフォームでは、ネイティブリクエストコンテキストをアプリケーションラッパー経由で転送（#5963、#5969）
- 競合するプッシュが発生しても、再試行によって変更履歴のエントリを保持（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6094)）
- モジュールに一度も同梱されなかったバイナリを参照する埋め込みを削除し、すべてのプラットフォームで `go mod vendor` が `pattern arm64/WebView2Loader.dll: no matching files found` により失敗する問題を修正。[#5782](https://github.com/wailsapp/wails/issues/5782) および [#5376](https://github.com/wailsapp/wails/issues/5376) を修正（@Grantmartin2002 による [PR](https://github.com/wailsapp/wails/pull/6031)）

## 削除

- 純粋な Go ローダーに置き換えられたため、ネイティブ WebView2 ローダーのサポートを削除しました。これにより、埋め込まれていた `WebView2Loader.dll` バイナリと `github.com/jchv/go-winloader` 依存関係が削除されます。`native_webview2loader` ビルドタグは引き続き受け入れられ、エラーは発生しませんが、v3 ビルドには影響しません（@Grantmartin2002 による [PR](https://github.com/wailsapp/wails/pull/6031)）
- macOS API ガイドから未使用のビルドタグと FPS オプションを削除（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6097)）

## v3.0.0-beta.19 - 2026-09-09

## 追加

- プライベートな macOS API をビルドタグで制限し、明示的に有効化した場合のみ使用可能に変更。[ドキュメント](https://v3.wails.io/features/browser/integration)、[ドキュメント](https://v3.wails.io/features/environment/info)、[ドキュメント](https://v3.wails.io/features/windows/basics)、[ドキュメント](https://v3.wails.io/features/windows/frameless)、[ドキュメント](https://v3.wails.io/features/windows/notch-windows)、[ドキュメント](https://v3.wails.io/features/windows/options)、[ドキュメント](https://v3.wails.io/guides/build/macos)、[ドキュメント](https://v3.wails.io/guides/build/private-macos-apis)、[ドキュメント](https://v3.wails.io/reference/overview)を参照（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6087)）

## 修正

- 64 MiB を超えるランタイムリクエストを HTTP 413 で拒否（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6091)）

## セキュリティ

- トークン認証により MCP のオリジンとリモートアクセスを強化（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/6092)）

## v3.0.0-beta.18 - 2026-09-08

## 修正

- ポインターレシーバーを使用し、Linux および Darwin での Calloc のメモリリークを修正（@4RH1T3CT0R7 による [PR](https://github.com/wailsapp/wails/pull/6083)）

## v3.0.0-beta.17 - 2026-09-06

## 修正

- Windows：WebResourceRequested ハンドラー内で `GetRequest` が失敗するか nil になっても、プロセスが終了しないように修正（`log.Fatal` / nil 参照による panic）。代わりにリクエストを破棄してログに記録（@midagedev による [PR](https://github.com/wailsapp/wails/pull/6006)）

## v3.0.0-beta.16 - 2026-08-29

## 変更

- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/6029) で、公証用パスワードを新しいターミナルウィンドウで入力するように変更

## 修正

- @ChewbaccaCookie による [PR](https://github.com/wailsapp/wails/pull/5919) で、macOS のシステムトレイのクリック種別を正しく処理
- @Grantmartin2002 による [PR](https://github.com/wailsapp/wails/pull/6041) で、CI が更新前に未使用の Microsoft apt リポジトリを削除するように変更

## v3.0.0-beta.15 - 2026-08-27

## 修正

- @Grantmartin2002 による [PR](https://github.com/wailsapp/wails/pull/6043) で、WebView2 の埋め込みタイムアウトを 60 秒に延長

## v3.0.0-beta.14 - 2026-08-26

## 修正

- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/6032) で、macOS における Control＋文字キーの押下名を正しく設定
- @nik9play による [PR](https://github.com/wailsapp/wails/pull/6016) で、Windows の ICO トレイアイコンを修正し、タスクバーのテーマに追従するように変更

## v3.0.0-beta.13 - 2026-08-25

## 修正

- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/6026) で、モーダルループの実行中も macOS のメインスレッド処理を継続するように修正
- @mortenolsrud による [PR](https://github.com/wailsapp/wails/pull/5923) で、モバイルのセキュアストレージ操作が失敗を返せるようにし、失敗時は安全側に閉じるように変更
- @archy-rock3t-cloud による [PR](https://github.com/wailsapp/wails/pull/5999) で、リスナーが登録されていない場合でもアプリケーションイベントフックを実行するように修正
- @haoku123 による [PR](https://github.com/wailsapp/wails/pull/6023) で、コメントとローカライズ済みドキュメントの誤字を修正
- @4RH1T3CT0R7 による [PR](https://github.com/wailsapp/wails/pull/6025) で、`v3/examples` 配下にコミットされていたコンパイル済み macOS バイナリを削除

## v3.0.0-beta.12 - 2026-08-21

## 追加

- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/6010) で、ライフサイクルとテレメトリの例を備えた macOS ノッチ通知ウィンドウを追加 — [ドキュメント](https://v3.wails.io/features/windows/notch-windows)を参照
- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/6008) で、新しいオプションとネイティブ統合を備えた macOS NSPanel ウィンドウのサポートを追加 — [ドキュメント](https://v3.wails.io/features/windows/options)を参照

## 修正

- @archy-rock3t-cloud による [PR](https://github.com/wailsapp/wails/pull/5998) で、SQLite の Prepare が同時呼び出し時にハングする問題を修正

## v3.0.0-beta.11 - 2026-08-20

## 削除

- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/6005) で、ドキュメントから廃止済みの実装トラッカーを削除

## v3.0.0-beta.10 - 2026-08-19

## 修正

- @midagedev による [PR](https://github.com/wailsapp/wails/pull/6000) で、GTK4 Linux ホストがカスタムプロトコルおよびファイル関連付けによる起動引数を破棄する問題を修正
- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5993) で、変更履歴の検証において削除された行と同一ソース内の修正を正しく処理

## v3.0.0-beta.9 - 2026-08-16

## 追加

- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5896) で、エージェント支援によるプロジェクト管理向けのセキュアな wails3 mcp サーバーを追加
- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5988) で、バインディング内のモデルに関するドキュメントを追加 — [ドキュメント](https://v3.wails.io/features/bindings/models)を参照
- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5987) で、アトミック Linux システムへの rpm-ostree を使用したインストールをサポート
- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5986) で、ネイティブな週次スター履歴チャートの生成と公開を追加 — [ドキュメント](https://v3.wails.io/credits)、[ドキュメント](https://v3.wails.io/de/credits)、[ドキュメント](https://v3.wails.io/fr/credits)、[ドキュメント](https://v3.wails.io/id/credits)、[ドキュメント](https://v3.wails.io/ja/credits)、[ドキュメント](https://v3.wails.io/ko/credits)、[ドキュメント](https://v3.wails.io/pt/credits)、[ドキュメント](https://v3.wails.io/ru/credits)、[ドキュメント](https://v3.wails.io/zh-cn/credits)、[ドキュメント](https://v3.wails.io/zh-tw/credits)を参照
- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5965) で、アプリケーションバンドルのリソースを解決する Darwin 専用の mac パッケージを追加 — [ドキュメント](https://v3.wails.io/guides/build/macos)を参照
- @mgueregath による [PR](https://github.com/wailsapp/wails/pull/5962) で、Condui のショーケースページとインデックス項目を追加 — [ドキュメント](https://v3.wails.io/community/showcase/condui)および[ドキュメント](https://v3.wails.io/community/showcase)を参照
- @redisviewer による [PR](https://github.com/wailsapp/wails/pull/5984) で、スクリーンショットとプロジェクトリンクを含む Redis Viewer のショーケースページを追加 — [ドキュメント](https://v3.wails.io/community/showcase)および[ドキュメント](https://v3.wails.io/community/showcase/redisviewer)を参照

## 変更

- @overlordtm による [PR](https://github.com/wailsapp/wails/pull/5971) で、Linux の GTK アプリケーションフラグを G<em>APPLICATION</em>NON_UNIQUE に更新
- @julianstorer による [PR](https://github.com/wailsapp/wails/pull/5914) で、見つからないウィンドウイベントを警告レベルではなくデバッグレベルでログに記録するように変更

## 修正

- @julianstorer による [PR](https://github.com/wailsapp/wails/pull/5902) で、登録済みの macOS キーボードショートカットが webview より優先されるように修正
- @northes による [PR](https://github.com/wailsapp/wails/pull/5937) で、ドキュメントのサイドバーにある壊れたリンクを修復
- WebKit が対応するカスタムスキームタスクを中止したとき、macOS および iOS のアセットリクエストコンテキストをキャンセルするように修正（#5963）
- @archy-rock3t-cloud による [PR](https://github.com/wailsapp/wails/pull/5976) で、WindowSetFullscreenButtonEnabled メッセージを処理
- @haoku123 による [PR](https://github.com/wailsapp/wails/pull/5979) で、ビルド失敗を解消するため preact-ts テンプレートに Fragment をインポート
- 従来の GTK3 サービス専用アプリケーションで、アクティブなウィンドウまたはディスプレイが利用可能になる前に画面検出が実行された場合にクラッシュする問題を修正（#5966）
- 未リリースの変更履歴が空の場合でも、明示的なバージョンを指定したリリース処理を続行できるように修正（#5977）

## セキュリティ

- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5985) で、セキュリティアドバイザリに対処するため、ウェブサイトの nanoid ロックファイルを修正済みの 3.3.18 に更新

## v3.0.0-beta.8 - 2026-08-12

## 追加

- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5957) で、自動生成される変更履歴の項目にドキュメント URL の生成を追加
- Streams を追加：WebSocket プログラミングモデルを使用し、リッスンソケットを必要としない、Go と JavaScript 間の双方向バイトストリーム。Go では `app.HandleStream(name, handler)` を使用してストリームを宣言し、フロントエンドでは `Stream(name)` で接続します。これは `WebSocket` 形式のオブジェクトを返します。Go→JS はアセットサーバーを介してウィンドウごとに保持される 1 つのポーリングで、JS→Go は通常の POST で転送されます。TCP ポートには何もバインドされず、`evaluateJavaScript` を経由するものもありません。サーバービルド（`-tags server`）では、代わりに同じハンドラーが実際の WebSocket 経由で提供されるため、アプリケーションコードはどのビルドでも同一です。@leaanthony
- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5935) で、mailbox の変更履歴項目を Unreleased に移動

## 変更

- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5938) で、ドキュメントのサイドバー自動生成とブログ著者型の導出を更新

## 修正

- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5952) で、WebView2 の初期化に期限とメッセージポンプを使用
- WebView2 の Cookie テストは、オプトインされていない限り CI ではスキップし、実行を現在の OS スレッドに固定する（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5951)）
- Windows のメニュービルダーで、サブメニューの親項目のコマンド ID を復元（@gilad-ch による[PR](https://github.com/wailsapp/wails/pull/5944)）
- 公式クロスコンパイルイメージを GTK 4.14以降の Linux サポート基準に合わせる（#5928）
- iOS Xcode プロジェクトで継承されたリンカーフラグを保持し、-ObjC を追加するように構成（@mortenolsrud による[PR](https://github.com/wailsapp/wails/pull/5915)）
- 大規模なフロントエンドで`wails3 dev`アセットプロキシの TCP 接続が過剰に入れ替わり、ホストのエフェメラルポートを使い果たして無関係なプロセスが`EADDRNOTAVAIL`で失敗する可能性がある問題を修正
- 順序どおりのディスパッチとバックプレッシャーのため、ウィンドウごとにイベント JavaScript をキューイング（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5934)）

## 削除

- デスクトップバイナリのリリースパイプラインを削除。v3リリースはタグのみとなり、`wails3` CLI は`go install`でインストールする。`release-v3.yml`と、それをディスパッチしていた nightly ステップを削除（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5946)）

## v3.0.0-beta.7 - 2026-08-11

## 追加

- メディア再生時のユーザー操作要件を無効にする macOS の自動再生設定を追加（@Eyalm321 による[PR](https://github.com/wailsapp/wails/pull/5512)）
- メールボックスの変更履歴エントリを未リリースへ移動（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5935)）

## 変更

- macOS のズームアニメーションで CADisplayLink または NSTimer を使用し、動作をより滑らかに改善（@savely-krasovsky による[PR](https://github.com/wailsapp/wails/pull/5945)）

## 修正

- iOS Xcode プロジェクトで継承されたリンカーフラグを保持し、-ObjC を追加するように構成（@mortenolsrud による[PR](https://github.com/wailsapp/wails/pull/5915)）
- 大規模なフロントエンドで`wails3 dev`アセットプロキシの TCP 接続が過剰に入れ替わり、ホストのエフェメラルポートを使い果たして無関係なプロセスが`EADDRNOTAVAIL`で失敗する可能性がある問題を修正
- 順序どおりのディスパッチとバックプレッシャーのため、ウィンドウごとにイベント JavaScript をキューイング（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5934)）

### 追加

- イベントを順序どおりに配信するための汎用非同期 FIFO メールボックスを実装（@savely-krasovsky と @DevLumuz による[PR](https://github.com/wailsapp/wails/pull/5851)）

## v3.0.0-beta.6 - 2026-08-09

## 追加

- サイズ超過イベント用の容量制限付きホスト側ストレージと、順序どおりの JavaScript 配信を実装（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5930)）
- ウィンドウの点滅時に macOS Dock アイコンをバウンドさせる機能を実装（@julianstorer による[PR](https://github.com/wailsapp/wails/pull/5921)）

## 修正

- アセットサーバーがフラッシュ時に、コンテンツタイプのスニッフィングエラーと未書き込みのプレフィックスを保持するように修正（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5931)）
- Wails のコールバックからアプリケーションメニューを置き換えた際に、macOS アプリケーションがクラッシュするのを防止
- Windows 10 1809 / Windows Server 2019（ビルド 17763）でネイティブメニューが判読不能になる問題を修正。ダークモードの uxtheme エクスポートがビルド 18334以降に限定されていたため、これらのホストではアプリレベルのダークモードへのオプトインが一度も実行されなかった。その結果、メニュー背景は暗く描画される一方で、Windows はライトテーマのメニューテキストを描画し続け、暗い背景に暗いテキストが表示されていた。序数によるエクスポートは17763以降に存在するため、条件をそれに合わせて修正。
- `w32.GetStockObject`が`GetStockObject`ではなく`GetDeviceCaps`を呼び出し、すべてのストックオブジェクトに対して0を返していた問題を修正。
- WebView2 ブートストラッパーのダウンロードエラー処理と報告を改善（@jannskiee による[PR](https://github.com/wailsapp/wails/pull/5924)）

## v3.0.0-beta.5 - 2026-08-07

## 修正

- macOS アプリのアクティベーションが、通常のアプリに限ってアクティベーションポリシーに従うように修正（@julianstorer による[PR](https://github.com/wailsapp/wails/pull/5897)）
- Linux ビルドで未初期化の GTK ウィンドウを保護（@julianstorer による[PR](https://github.com/wailsapp/wails/pull/5898)）
- Linux の WebKit ウィンドウで、URL の読み込み前に明示的な不透明背景色を設定（@julianstorer による[PR](https://github.com/wailsapp/wails/pull/5899)）

## v3.0.0-beta.4 - 2026-08-05

## 変更

- Android ビルドタスクのデフォルトを arm64 に変更し、deploy-emulator がホストアーキテクチャを選択するように変更（@mortenolsrud による[PR](https://github.com/wailsapp/wails/pull/5890)）

## 修正

- ドラッグ中も macOS ウィンドウのズーム状態を維持し、動きを抑制（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5900)）
- `webview_window_windows_nonclient.go`のビルド制約に`!server`を追加し、Windows のサーバーモードビルドを修正

## v3.0.0-beta.3 - 2026-08-03

## 追加

- 実装の詳細にフェーズ10のベータ検証完了を記載（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5881)）

## 修正

- Windows のダークモード API にウィンドウハンドルを渡し、引数を検証（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5877)）
- フレームレスウィンドウにおける macOS タイトルバーボタンの状態解決を一元化（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5870)）
- Windows のアプリテーマがライトであるにもかかわらず、Windows アプリケーションがダークモードを要求した際に、ネイティブメニューのテキストが判読不能になるのを防止。Windows が暗いメニューテキストを描画できるようになるまでは、メニューに対応する明るいネイティブ背景を使用する。
- Windows 10 1809 / Windows Server 2019（ビルド 17763）でネイティブメニューが判読不能になる問題を修正。ダークモードの uxtheme エクスポートがビルド 18334以降に限定されていたため、これらのホストではアプリレベルのダークモードへのオプトインが一度も実行されなかった。その結果、メニュー背景は暗く描画される一方で、Windows はライトテーマのメニューテキストを描画し続け、暗い背景に暗いテキストが表示されていた。序数によるエクスポートは17763以降に存在するため、条件をそれに合わせて修正。

## v3.0.0-beta.2 - 2026-08-02

## 変更

- v3をアルファからベータへ昇格
- システムトレイのスマートデフォルトとポップアップの自動非表示動作を文書化し、クリックハンドラー選択の回帰テストを追加（#5840）。
- GitHub アップデーターがデフォルトで Windows インストーラーのアセットを除外するように変更（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5861)）
- macOS のフレームレスウィンドウで、角丸、直角、およびカスタム半径の角をサポート（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5866)）

## 修正

- GTK4 ウィンドウの現在のサイズを報告し、構成済みサーフェスからサイズ変更、最大化、最小化、全画面表示の状態イベントを発行（#5830）。
- fetch リクエストで Blob または FormData を送信した際に Linux の WebKit がクラッシュする問題を修正（@taliesin-ai による[PR](https://github.com/wailsapp/wails/pull/5854)）
- fetch shim が、存在しない Blob/FormData ヘッダーに undefined を渡すように修正（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5865)）

## v3.0.0-alpha2.122 - 2026-08-01

## 追加

## 変更

- macOS のフレームレスウィンドウで、角丸、直角、およびカスタム半径の角をサポート（@leaanthony による[PR](https://github.com/wailsapp/wails/pull/5866)）

## 修正

- Fetch shim で Blob/FormData ヘッダーがない場合に undefined を渡すようにした（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5865)）

## v3.0.0-alpha2.121 - 2026-07-31

## 追加

- 新しいオプションとビルドタスクを備えた macOS DMG パッケージ作成サポートを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5857)）

## 変更

- GitHub updater で Windows インストーラーのアセットをデフォルトで除外するように変更（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5861)）

## 修正

- fetch リクエストで Blob または FormData を送信した際に Linux WebKit がクラッシュする問題を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5854)）

## v3.0.0-alpha2.120 - 2026-07-31

## 追加

- macOS のタイトルバーをダブルクリックした際に最大化または最小化するアクションを実装（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5853)）

## 変更

- QR サービスのチュートリアルを NewServiceWithOptions を使用するように更新し、余白を追加（@jeongkyu による [PR](https://github.com/wailsapp/wails/pull/5849)）

## 修正

- macOS でのズーム中も WKWebView が応答を維持するように修正（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5856)）
- GTK4 のウィンドウサイズ取得を修正し、設定済みの `GdkSurface` から、サイズ変更、最大化、最小化、全画面表示の状態イベントを発行するようにした。

## v3.0.0-alpha2.119 - 2026-07-27

## 修正

- 複数言語のドキュメントを更新し、アーキテクチャ図を追加（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5833)）

## v3.0.0-alpha2.118 - 2026-07-26

## 追加

- アイコン生成の入力および出力にデフォルトパスを提供（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5825)）
- runtime の package.json の sideEffects にソースエントリーモジュールを追加（@savely-krasovsky による [PR](https://github.com/wailsapp/wails/pull/5797)）
- コントリビューションガイドにライセンスと来歴のセクションを追加（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5816)）

## 修正

- 適用範囲を限定した GTK4 フレームレス用 CSS を適用し、角丸を削除（@savely-krasovsky による [PR](https://github.com/wailsapp/wails/pull/5800)）
- Windows でポップアップメニューや画面列挙のためのカーソル位置取得に失敗した場合を安全に処理（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5789)）
- macOS のファイルを開くダイアログで拡張子を正しくフィルタリングし、許可されたファイルを接尾辞で検証するように修正（@phergul による [PR](https://github.com/wailsapp/wails/pull/5678)）
- Windows のダークモード初期化で nil の API が呼び出されないように保護（@roachadam による [PR](https://github.com/wailsapp/wails/pull/5793)）
- updater の 32 ビットビルドが失敗する問題を修正。`GOARCH=386` 上で `fmt.Errorf` に渡した際、`maxArchiveTotalSize` 定数（2 GiB）がプラットフォームの `int` でオーバーフローしていた。現在は明示的に `int64` 型を指定している。
- ダークモード用の uxtheme API を読み込まない Windows ビルド（Windows 10 1809 / Windows Server 2019（ビルド 17763）など）で、ウィンドウが Dark（またはシステムのダークモード）のタイトルバーを使用している場合、起動時に nil ポインターによるパニックが発生する問題を修正。ウィンドウテーマ設定内の `AllowDarkModeForWindow` 呼び出しに、`w32.SetMenuTheme` ですでに使用されているものと同じ nil ガードを追加した。

## v3.0.0-alpha2.117 - 2026-07-08

## 追加

- Windows の非クライアント領域向けカスタムヒットテストロジックを実装（@savely-krasovsky による [PR](https://github.com/wailsapp/wails/pull/5462)）

## 変更

- UseVisualHosting に基づいて WebView2 のモニタースケール検出を設定（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5761)）

## v3.0.0-alpha2.116 - 2026-07-07

## 追加

- Wails v3 の機能とガイダンスに重点を置くよう FAQ ドキュメントを更新（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5763)）

## v3.0.0-alpha2.115 - 2026-07-06

## 修正

- GTK4 Linux で `Menu.Update()` がネイティブメニューを再構築しない問題を修正（#5659。@puneetdixit200 が #5539 で独自に診断および修正）
- ディスプレイ変更時に macOS の画面を列挙するとクラッシュする問題を、画面の ID/名前文字列をコピーし、画面数のスナップショットを取得することで修正（#5565。@x-haose が #5584 で独自に診断および修正）
- 最小化から復元への遷移中に `GetClientRect` が nil を返した際、`WM_ERASEBKGND` が単色の背景を描画すると Windows でクラッシュする問題を修正（ガードは #5636 で @sinspired が報告）
- フロントエンドバインディングのエラーが常にテキストとして解析される問題を修正（#5690、@mbaklor）
- `server` ビルドタグを使用すると Windows ビルドが失敗する問題を修正。原因は、macOS および Linux の同等ファイルにはすでに存在する `!server` ビルド制約が、Windows GUI ファイルにはなかったこと（#5680）

## v3.0.0-alpha2.114 - 2026-07-05

## 追加

- Update Manifest プロトコルとエンドポイントプロバイダーを実装（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5720)）

## 変更

- `webview2` バインディングを v3 モジュールに `v3/internal/webview2` として統合し、独立モジュール、その nightly リリース／同期ワークフロー、および go.mod のバージョン調整処理を削除（唯一の利用側は v3）（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5711)）

## 修正

- WebView2 のモニタースケール検出と、DPI 変更時のホスト再同期の修正を Unreleased セクションへ移動（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5750)）
- float64 および BOOL パラメーター向けに WebView2 の COM マーシャリングを更新（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5741)）
- Windows のシステムトレイアイコンの更新および破棄で、パニックと nil の参照解除を防止（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5703)）
- Windows で非表示ウィンドウが正しく再非表示にならない問題を修正（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5743)）
- WebView2 コントローラーの表示状態をウィンドウの最小化、最大化、復元と同期（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5742)）

### 修正

- WebView2 のモニタースケール検出を再有効化し、DPI 変更時のホスト再同期を条件付きで実行するように修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5734)）。@randalmurphal が検証した修正を基に、@eleclin が根本原因を確認し、@qq540491950 が実機テストを実施

## v3.0.0-alpha2.113 - 2026-07-04

## 追加

- `ANDROID_KEYSTORE_FILE` を設定せずにリリース用 AAB をビルドする場合に警告を追加し（Google Play はデバッグ署名されたバンドルを拒否する）、App Bundle のパッケージ作成と署名について文書化（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5730)）
- 複数言語で「Why Wails」ドキュメントを追加（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5739)）
- バインディングで Go の time.Time から JS の Date または string へのマッピングをサポート（@fbbdev による [PR](https://github.com/wailsapp/wails/pull/5398)）
- Play Store への提出用に Android App Bundle（AAB）のパッケージングタスク（`bundle`、`bundle:fat`、`assemble:aab`、`assemble:aab:release`）を追加。ローカル／エミュレーターでのテスト用 APK タスクは引き続き利用可能（@mortenolsrud による [PR](https://github.com/wailsapp/wails/pull/5728)、[#5726](https://github.com/wailsapp/wails/issues/5726) を修正）
- Android 実機用のタスクターゲットと、カメラ／位置情報権限の再開処理を追加（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5735)）

## 変更

- `webview2` を v1.0.28 に更新（[リリースノート](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)）。
- Android テンプレートの `compileSdk`／`targetSdk` を 34 から 35 に更新。Google Play では新規アプリの提出にこのバージョンが必要（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5730)）

## 修正

- sponsorkit で焼き込み済みのアバターマスクを修正（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5745)）
- 辞書順によるバージョンの並べ替えが原因で、Android AVD の自動作成時に誤ったシステムイメージまたは cmdline-tools のバージョンが選択される問題を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5730)）
- セットアップウィザードが古い Android NDK バージョンを提案する問題を修正（現在は文書化された要件と一致する 26.3.11579264）（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5730)）
- SvelteKit とオプションに関するフランス語ドキュメントを更新（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5744)）
- ディスプレイ変更中に macOS の画面を列挙すると SIGSEGV が発生する問題を修正（@flofreud による [PR](https://github.com/wailsapp/wails/pull/5516)）

## v3.0.0-alpha2.112 - 2026-07-03

## 追加

- Go ベースのコントリビューター SVG ジェネレーターを追加し、ドキュメント／Web サイトのクレジットページを更新（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5724)）
- バインディングで Go の time.Time を JS の Date または string にマッピングできるよう対応（@fbbdev による [PR](https://github.com/wailsapp/wails/pull/5398)）

## 変更

- Node ベースのスポンサー画像パイプラインを Go ジェネレーターに置き換え（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5719)）

## 修正

- Android ビルドアセットの依存関係インストールスクリプトを修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5729)）
- `ValidateAndSanitizeURL` で U+0085（NEXT LINE）制御文字を拒否するようにし、URL バリデーターの空白文字対応を完成
- フレームレスウィンドウで DPI が変更された際に DWM フレームを再計算（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/4785)）
- Windows でスケーリングが 100% 以外の場合に DnD ドロップゾーンの検出が失敗する問題を修正（@yulesxoxo による [PR](https://github.com/wailsapp/wails/pull/4632)）
- Darwin のダイアログ、メニュー、トレイ、通知全体の Cocoa オブジェクトに、明示的な Objective-C メモリ管理を追加（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5714)）
- Linux CGO バックエンドのバグとシステムトレイの問題を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5718)）

## v3.0.0-alpha2.111 - 2026-07-01

## 追加

- コミュニティショーケースに HappyTools を追加（@Aliuyanfeng による [PR](https://github.com/wailsapp/wails/pull/5061)）
- インドネシア語ロケールのサポートと包括的なドキュメントを追加（@triadmoko による [PR](https://github.com/wailsapp/wails/pull/5643)）
- WindowsWindow に DisableMenu オプションを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/4813)）

## 変更

- GOOS と ARCH を使用してビルド／パッケージタスクを振り分けるよう、Taskfile テンプレートと CLI を更新（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5617)）

## 修正

- macOS のウィンドウタブ化に関する問題を修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5708)）

## 削除

- contributing、features、guides からドイツ語翻訳済み MDX ファイルを削除（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5702)）

## v3.0.0-alpha2.110 - 2026-06-30

## 追加

- macOS の WebView の再読み込みと強制再読み込みを実装し、WebContent プロセス終了時の復旧処理を追加（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5129)）
- contributing、features、guides に関する包括的なドイツ語ドキュメントを追加（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5396)）
- 通知にサウンド、添付ファイル、スケジュール設定、更新 API を追加して機能を強化（@popaprozac による [PR](https://github.com/wailsapp/wails/pull/5333)）

## 修正

- フレームレスウィンドウで DPI が変更された際に DWM フレームを再計算（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/4785)）
- Windows でスケーリングが 100% 以外の場合に DnD ドロップゾーンの検出が失敗する問題を修正（@yulesxoxo による [PR](https://github.com/wailsapp/wails/pull/4632)）

## v3.0.0-alpha2.109 - 2026-06-29

## 追加

- EventsEmit ドキュメントにコードサンプルを追加（@iamhabbeboy による [PR](https://github.com/wailsapp/wails/pull/5026)）
- Windows WebView2 のビジュアルホスティングオプションを追加（@MerIijn による [PR](https://github.com/wailsapp/wails/pull/5380)）
- コミュニティショーケースのドキュメントに Klustr を追加（@SametKUM による [PR](https://github.com/wailsapp/wails/pull/5536)）
- 新しいページと変更履歴の項目を含め、コミュニティショーケースに Kira を追加（@thiennguyen93 による [PR](https://github.com/wailsapp/wails/pull/5685)）
- MCP サービスガイドにフィードバックセクションを追加（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5694)）

## 変更

- サーバーモードで、デスクトップのビルドタスクと一貫性のある正式な本番ビルドを利用可能にしました（#5693）。`task build:server` はデフォルトで本番用バイナリ（`-tags server,production`、`-trimpath`、シンボル削除済み）をビルドし、`DEV=true`（開発サーバー）、`OBFUSCATED=true`（garble）、`EXTRA_TAGS` を受け付けます。`task run:server` は開発サーバーを実行します。`Dockerfile.server`／`task build:docker` は、まず本番用サーバー（`-tags server,production`）と本番用フロントエンドをビルドします。イメージのデフォルトは distroless/static 上の Pure Go 静的ビルドで、CGO アプリ向けに `CGO_ENABLED`、`GO_IMAGE`、`RUNTIME_IMAGE` を上書き可能なビルド引数として公開します。

## 修正

- 保留中の非同期呼び出しがある状態でウィンドウを閉じた際のクラッシュを防止（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/4435)）
- Windows で非表示のアプリを開く際にウィンドウがアクティブ化されないよう修正（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5249)）
- WebKit のリクエストメタデータ、レスポンス完了処理、ボディストリーム処理が GTK メインスレッドで実行されるよう修正（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5668)）
- GTK4 Linux で `Menu.Update()` がネイティブメニューを再構築しない問題を修正（#5659、@puneetdixit200 も #5539 で独自に診断および修正）
- 画面 ID／名前の文字列をコピーして画面数のスナップショットを取得することで、ディスプレイ変更時に macOS の画面を列挙するとクラッシュする問題を修正（#5565、@x-haose も #5584 で独自に診断および修正）
- Windows で異なる DPI のモニター間をまたいでウィンドウをドラッグした後、WebView2 のコンテンツが縮小して消える問題を修正。最小化解除時の DPI 再同期と同様に、`WM_DPICHANGED` ハンドラーでコントローラーの境界を再設定（#5677）

## v3.0.0-alpha2.108 - 2026-06-28

## 追加

- `app.GlobalShortcut`（`Register`、`Unregister`、`UnregisterAll`、`IsRegistered`、`GetAll`）によるグローバル（システム全体）キーボードショートカットを追加。アプリケーションにフォーカスがない場合でもショートカットが発動します。サードパーティ依存関係を使用せず、プラットフォームごとにネイティブ実装されています。macOS では Carbon ホットキー、Windows では `RegisterHotKey`、X11 では `XGrabKey`、Wayland では XDG Desktop Portal のグローバルショートカットインターフェースを使用します。
- 組み込み MCP サーバーを追加。`mcp` タグ付きでアプリケーションをビルドすると自動的に起動する Model Context Protocol サーバーで、LLM エージェントが実行中の Wails アプリケーションをテストおよび制御できます。ウィンドウ制御、DOM 検査、JavaScript 評価、バインドされたメソッドの呼び出し、イベント、画面上にアニメーション表示されるカーソルによるマウス／キーボード入力のシミュレーションに対応します。ユーザーコードは不要です。`WAILS_MCP=1` が設定されている場合、`wails3 build`／`wails3 dev` によって `mcp` タグが自動的に追加されます。設定はすべて環境変数（`WAILS_MCP_HOST`、`WAILS_MCP_PORT`、`WAILS_MCP_TIMEOUT`、`WAILS_MCP_HIDE_CURSOR`）で行います。

## 修正

- GTK4 Linux で `Menu.Update()` がネイティブメニューを再構築しない問題を修正（#5659、#5539 では @puneetdixit200 が独自に診断して修正）
- 画面 ID／名前の文字列をコピーし、画面数のスナップショットを取得することで、ディスプレイ変更時に macOS の画面を列挙するとクラッシュする問題を修正（#5565、#5584 では @x-haose が独自に診断して修正）

## v3.0.0-alpha2.107 - 2026-06-27

## 追加

- @leaanthony による [PR](https://github.com/wailsapp/wails/pull/5613) で、サイドバーナビゲーションを備えた実験的な Wake ドキュメントを追加

## v3.0.0-alpha2.106 - 2026-06-24

## 変更

- `webview2` を v1.0.27 に更新。
  - ci(webview2): リリースビルドを修正（Windows のクロスコンパイル＋完全な go.sum）（#5671）\

  **差分全体：** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5672) で、webview2 リリースワークフローのクロスコンパイルから go vet を削除
- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5670) で、自動変更履歴生成用の OpenRouter モデルを google/gemini-2.5-flash-lite に更新
- `webview2` を v1.0.26 に更新。

### 修正

- **一時的なランタイム COM エラーで終了せず、復旧するように変更**（#5658、#5580）。従来、`Chromium.errorCallback` は *あらゆる* COM エラーに対して `os.Exit(1)` を呼び出していたため、起動後に復旧可能な一時的不具合が発生しただけでも、アプリケーション全体が終了していました。ランタイムの処理経路（`Resize`／`GetClientRect`、`Navigate`／`NavigateToString`、`Init`、`MessageReceived`、`PutZoomFactor`、`OpenDevToolsWindow`）では、エラーをログに記録して復旧するようになりました。特に、`MessageReceived` で不正な、または信頼できない Web メッセージを受信した場合、プロセスを停止させずにそのメッセージを破棄するようになりました。これにより、異なる DPI のモニター間を移動した際に発生する一連のクラッシュ（#5544、#5650）に対処します。環境／コントローラー作成の処理経路では、引き続き致命的エラーとして扱います。\

**差分全体：** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## 修正

- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5671) で、go.sum ファイルを正しく処理するよう release-webview2 ワークフローを修正
- @taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5659) で、ネイティブメニューをクリアして再構築することにより Linux GTK4 のメニュー更新を修正

## v3.0.0-alpha2.105 - 2026-06-21

## 追加

- 共有コードから実行時にプラットフォームを検出するための `application.System` を追加。`System.IsMobile()`（iOS／Android）、`System.IsDesktop()`（macOS／Windows／Linux）、`System.IsServer()`（`server` ビルドタグ）に加え、単一のターゲットを直接テストするための `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` を提供します。すべてのターゲットでコンパイルできるため、ビルドタグを使わずに処理を分岐できます。対応するフロントエンドヘルパー（`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`）は `@wailsio/runtime` で利用できます
- 独自の Vite プロジェクトを `frontend/` に組み込む方法を示す「その他のフロントエンドフレームワークの使用」ガイドを追加（Solid、Preact、Lit、SvelteKit、Qwik、Angular などを対象）
- `wails3 setup` ウィザードで、モバイル（iOS／Android）ツールチェーン（Xcode と iOS Simulator ランタイム、JDK、Android SDK／NDK、エミュレーター）を確認するようになりました。該当する場合は、ワンクリックインストールとコピー可能なシェル設定修正を利用できます
- 生成されるプロジェクトには、新たに公開された（侵害されている可能性がある）パッケージにさらされるリスクを軽減するため、7 日の `minimum-release-age` を設定する `frontend/.npmrc` が含まれるようになりました（pnpm と bun では適用され、npm では問題なく無視されます）

## 変更

- 組み込みのすべてのスターターテンプレートを、ネオンの山をヒーロービジュアルに採用した新しいデザインに刷新（Web、iOS、Android）
- **スターターテンプレートでは TypeScript がデフォルトになり、接尾辞のないテンプレート名が TypeScript 用になりました。** `wails3 init`（`-t` なし）では TypeScript プロジェクトが生成されます。`-t vanilla`、`-t react`、`-t vue`、`-t svelte` は TypeScript 版で、JavaScript 版はそれぞれ `-t vanilla-js`、`-t react-js`、`-t vue-js`、`-t svelte-js` です。組み込みテンプレートは、`template.yaml` 内の `typescript:` で言語を宣言します。`-ts` 接尾辞を使用するコミュニティテンプレートも、フォールバックとして引き続き機能します
- `wails3 setup` ウィザードを、ネオンの「digital Wails」テーマ（山の背景にフロストガラスの鮮やかさを重ねたデザイン）に刷新

## 修正

- WebView2 が一時停止するか、そのレンダー／GPU プロセスが再生成されるほど長時間最小化されていたアプリを復元すると、Windows でクラッシュする問題を修正。最小化／復元時の DPI 再同期（#5544）では、ウィンドウの DPI が実際に変更された場合にのみ WebView2 コントローラーを操作するようになり、一般的な同一 DPI での復元時に、一時停止中のコントローラーに対して致命的な COM 呼び出しを行わないようになりました（#5605）
- アセット／メディアを頻繁に読み込む長時間稼働中の Linux アプリで、ネイティブの `SIGABRT`／`SIGSEGV` クラッシュ（通常は GTK メインループ中の `g_object_unref` 内で発生）が繰り返される問題を修正。アセットサーバーがワーカー goroutine から `WebKitURISchemeRequest` を完了させ、スレッドセーフではない WebKit2GTK 関数を GTK メインスレッド外で呼び出していました。完了処理（`webkit_uri_scheme_request_finish_with_response`／`finish_error`）はメインスレッドで実行されるようになりました。#5566 の部分的な修正を完成させるものです。GTK3 と GTK4／WebKitGTK の 6.0 ビルドの両方に影響します（#5631、#5557）
- Linux／GTK3 の `setupSignalHandlers` で断続的に発生する `fatal error: invalid pointer found on stack` を修正。シグナルの `user_data` として渡されるウィンドウ ID が Go のローカル変数 `unsafe.Pointer` に保持されていたため、スタックコピー中にガベージコレクターがこの（非ポインターの）値をスキャンした際に異常終了していました。Go 側では ID を整数型（`uintptr_t`）のまま保持するようになりました。これは、#4958 で GTK4 の処理経路に適用したものと同じ修正を、従来の GTK3 の処理経路にもバックポートしたものです（GTK4 では、`-race`／checkptr エラーを解消するため、C のシグナル関数を `uintptr_t` に変更）（#5631）

## 削除

- `react-swc`、`preact`、`lit`、`solid`、`qwik`、`sveltekit` のスターターテンプレート（およびそれぞれの `-ts` 版）を削除。サポート対象の組み込みテンプレートは `vanilla`、`react`、`vue`、`svelte` となりました。いずれもデフォルトは TypeScript で、`-js` JavaScript 版もあります。その他のフレームワークも、[独自のフロントエンドを持ち込む](https://v3.wails.io/guides/dev/frontend-frameworks)か、カスタムテンプレートを使用することで引き続き利用できます

## v3.0.0-alpha2.104 - 2026-06-18

## 修正

- バインドされた Go サービスメソッドが空文字列を返すと iOS でクラッシュ（SIGABRT）する問題を修正。iOS のアセットレスポンスライターは、本文ポインターを長さではなく `buf != nil` で保護していたため、長さゼロの本文によって `&buf[0]` がパニックしていました。デスクトップ版のライターと同様に、長さで保護するようになりました

## v3.0.0-alpha2.103 - 2026-06-15

## 変更

- iOS と Android のネイティブ機能をプラットフォームマネージャーに移行。従来の `application.IOS*`／`application.Android*` フリー関数ではなく、`application.IOS.*` と `application.Android.*`（例：`application.IOS.Haptic("medium")`、`application.Android.Share(payload)`）経由で呼び出すようになりました（#5602）
- モバイルブリッジイベントの名前を変更：クロスプラットフォームイベントは `common:*` プレフィックス（例：`common:haptic`、`common:location`）を、プラットフォーム固有イベントは `ios:*` / `android:*`（例：`ios:backgroundTask`、`android:foregroundService`）を使用するようになりました。`native:*` プレフィックスは使用されなくなりました（#5602）

## v3.0.0-alpha.102 - 2026-06-14

## 追加

- 対話形式のプロジェクトセットアップと依存関係チェックを行う、実験的な `wails3 setup` ウィザードを追加
- 機械可読な出力を生成する `--json` フラグを `wails3 doctor` に追加
- `wails3 doctor` コマンドに署名ステータスのセクションを追加

## 修正

- Linuxでのnpm検出を修正し、パッケージマネージャーに加えてPATHも確認するように変更

## v3.0.0-alpha.101 - 2026-06-13

## 追加

- iOS：ネイティブメッセージダイアログ（UIAlertController）と、ファイル／複数ファイル／ディレクトリを開くダイアログ（UIDocumentPickerViewController）。保存ダイアログは明示的なエラーを返します
- iOS：UIPasteboardによるクリップボード対応
- iOS：UIScreenによる実際の画面メトリクス（ポイント、ピクセル、スケール、セーフエリアを考慮した作業領域）
- iOS：実機ビルド（`IOS_PLATFORM=device`）、コード署名ID／プロビジョニングプロファイル／エンタイトルメントへの対応、`.ipa` パッケージング、およびdevicectlによる `deploy-device`
- iOS：iOSの最小バージョンを設定可能（build/config.ymlの `ios.minIOSVersion`）
- iOS：`wails3 doctor` がmacOS上でXcodeとiOS SDKの利用可否を報告
- iOS：バッテリー、ネットワーク、テーマ、画面ロック、メモリ不足のシステムイベントを、`events.IOS.*` およびプラットフォーム非依存の `events.Common.*` アプリケーションイベントとして公開
- iOS：ネイティブモバイル機能ブリッジ（公開された `application.IOS*`）— 共有シート、URLを開く、スリープ防止、トーチ、セーフエリアのインセット、明るさ、アプリ情報、画面方向のロック、ステータスバー、生体認証（Face ID／Touch ID）、ローカル通知、Keychainによる安全なストレージ
- iOS：センサーとハードウェア — 触覚フィードバック、単発の位置情報取得、加速度センサー、近接センサー、テキスト読み上げ、ストレージ情報、電源／バッテリー状態、ネットワーク状態、キーボードのインセット、画面キャプチャの検出
- iOS：ドキュメント（IOS.mdおよびドキュメントサイトのガイド）
- Android：ネイティブメッセージダイアログ（AlertDialog）と、ファイル／複数ファイルを開くダイアログ（Storage Access Framework。キャッシュのコピーとしてインポート）。ディレクトリを開くダイアログと保存ダイアログは明示的なエラーを返します
- Android：ClipboardManagerによるクリップボード対応
- Android：WindowMetrics／DisplayMetricsによる実際の画面メトリクス（dp、ピクセル、スケール、システムバーを考慮した作業領域）
- Android：触覚フィードバック（`Android.Haptics.Vibrate`）、デバイス情報（`Android.Device.Info`）、トースト（`Android.Toast.Show`）のランタイムメソッド
- Android：型付きライフサイクルイベント（`events.Android.*`、events.txtから生成）。`ActivityCreated` は `Common.ApplicationStarted` にマッピングされます
- Android：ビルドパイプラインがインストール可能なデバッグAPKとリリースAPKを生成（`android:run`、`android:package`、`android:package:fat`）。リリース署名には、デフォルトではデバッグ用キーストアを使用し、`ANDROID_KEYSTORE_*` 環境変数を介して実運用のキーストアも使用可能
- Android：`wails3 doctor` がAndroid SDK、NDK、JDKを報告
- Android：バッテリー、ネットワーク、テーマ、画面ロック、メモリ不足のシステムイベントを、`events.Android.*` およびプラットフォーム非依存の `events.Common.*` アプリケーションイベントとして公開
- Android：ネイティブモバイル機能ブリッジ（公開された `application.Android*`）— 共有、URLを開く、スリープ防止、トーチ、セーフエリアのインセット、明るさ、アプリ情報、画面方向のロック、ステータスバー、生体認証（BiometricPrompt）、ローカル通知、EncryptedSharedPreferencesによる安全なストレージ
- Android：センサーとハードウェア — 触覚フィードバック、単発の位置情報取得、加速度センサー、近接センサー、テキスト読み上げ、ストレージ情報、電源／バッテリー状態、ネットワーク状態、キーボードのインセット、FLAG_SECUREによる画面キャプチャのブロック
- Android：ドキュメント（ANDROID.mdおよびドキュメントサイトのガイド）
- サンプル：`mobile` kitchen sinkに、iOSとAndroidのネイティブ機能ブリッジを実演するMobileタブとHardwareタブを追加（ピル型タブは複数行に折り返されます）
- モバイル：バッテリー — アプリがバックグラウンドに移行すると、加速度センサー、近接センサー、トーチ、サンプルの定期実行時計を一時停止し、復帰時に再開します（Androidではバックグラウンドでもプロセスが実行され続け、iOSではトーチが持続するハードウェア状態であるため）。また、Androidのシステムイベントレシーバーは、アプリがフォアグラウンドにある間だけ登録されます
- iOS：カメラ撮影 — `application.IOSCapturePhoto`／`IOSCaptureVideo`（UIImagePickerController → base64サムネイルを含む `native:capture` イベント）
- iOS：バックグラウンド実行 — `application.IOSBeginBackgroundTask`／`IOSEndBackgroundTask`（UIApplicationのバックグラウンドタスク実行枠）と、生成されるInfo.plistへ `UIBackgroundModes` をテンプレート展開する、設定可能な `ios.backgroundModes`（build/config.yml）
- Android：カメラ撮影 — `application.AndroidCapturePhoto`／`AndroidCaptureVideo`（FileProvider経由のシステムカメラ → `native:capture` イベント）
- Android：フォアグラウンドサービス — `application.AndroidStartForegroundService`／`AndroidStopForegroundService`（継続中の通知を伴う `WailsForegroundService` により、長時間のバックグラウンド処理中もプロセスを存続）
- サンプル：写真／動画の撮影とバックグラウンド実行を実演するCameraタブ（Androidではフォアグラウンドサービス、iOSではバックグラウンドタスク実行枠）

## 修正

- Linuxで `getUserMedia` が常に `NotAllowedError` で失敗する問題を修正：WebKitGTKは処理するものがない権限要求を拒否しますが、`permission-request` シグナルが接続されていませんでした。カメラ／マイクは、新しいクロスプラットフォームの `WebviewWindowOptions.Permissions` マップ（`map[PermissionType]Permission`）に基づいて処理され、Linux（WebKitGTK）とWindows（WebView2）の両方で適用されるようになりました。ネイティブの確認プロンプトがないLinuxでは、カメラ／マイクはデフォルトで許可され（`getUserMedia` を復元）、`PermissionDeny` で無効にできます（#5552）
- iOS：`GOOS=ios` が再びコンパイル可能になりました（`events.IOS` を公開し、モバイルのメソッド名スタブを追加）。また、productionタグ付きビルドもコンパイル可能になりました（pkg/applicationと複数のサービスでビルドタグを修正）
- iOS：Go→JSイベントとExecJSが動作するようになりました。起動時にページが二重に読み込まれることがなくなり、`wails:runtime:ready` ハンドシェイクも失われなくなりました
- iOS：`ApplicationDidFinishLaunching`／`ApplicationStarted` がアプリの起動と競合しなくなりました。固定の 2 秒の起動待機を削除しました
- iOS：Go→JSのJavaScript実行ごとに発生していたC文字列のリークを修正
- iOS：`hasListeners` が実際のリスナー登録状態を反映するようになりました
- iOS：本番ビルドではフレームワークのデバッグログをコンパイル時に除外
- Android：`GOOS=android` が再びコンパイル可能になりました。`events.Android` を定義し、範囲外アクセスを起こす `events_android.go` リスナー配列を削除し、モバイルのメソッド名スタブを追加するとともに、デスクトップLinux用ファイル（`linux_cgo.*`、`events_linux.*`、`environment_linux.go`）がAndroidビルドに混入しないようにしました
- Android：JS→Goバインディングが動作するようになりました。WebViewは `fetch()` のPOSTボディを `shouldInterceptRequest` に渡せないため、ランタイム呼び出しはJavascriptInterfaceトランスポート（`nativeHandleRuntimeCall`）を経由し、nilのリクエストボディでクラッシュしなくなりました
- Android：`Screens.*` ランタイム呼び出しが実際のデータを返すようになりました。起動時にScreenManagerを設定するようにしました（これまでは接続されていなかったため、`GetAll` がnilを返していました）
- Android：本番ビルドではフレームワークのデバッグログをコンパイル時に除外し、デバッグビルドでは `Wails` タグを使用してlogcatへ出力
- Android：実際に機能する `hasListeners` レジストリ、JNI参照／例外処理、および単一読み込みのページライフサイクル（二重ナビゲーションなし）
- WindowsでVite開発サーバーの実行中に `wails3 generate bindings` が「Access is denied」で失敗する問題を修正。出力ディレクトリを上書きするように名前変更する代わりに、生成ファイルを出力ディレクトリへ同期するようにしました（#5515）
- ディスプレイ変更後に画面情報を読み取る際、macOS で断続的に発生する致命的なクラッシュを修正しました。画面 ID と名前が、自動解放される `UTF8String` バッファへのポインターを保持しており、Go がコピーする前に解放される可能性がありました（解放後使用）。文字列は `strdup` され、変換後に解放されるようになりました。また、画面の列挙を明示的な自動解放プール内で実行することで、Go の goroutine から呼び出した場合にリークしないようにしました（#5556）
- assetserver が `WebKitURISchemeRequest` を閉じる際、Linux で断続的に発生する SIGSEGV を修正しました。最後の `g_object_unref` が assetserver の goroutine 上で実行され、GTK メインスレッド外で WebKit GObject をファイナライズしていました。unref は `g_main_context_invoke` を介して GTK メインコンテキスト上にマーシャリングされるようになりました（#5557）

## v3.0.0-alpha.100 - 2026-06-13

## 追加

- `MacWebviewPreferences` に WKWebView の追加設定オプション `EnableAutoplayWithoutUserAction`、`AllowsAirPlayForMediaPlayback`、`AllowsMagnification`、`JavaScriptCanOpenWindowsAutomatically`、`MinimumFontSize`、`ApplicationNameForUserAgent` を追加しました（#5549）

## 修正

- Vite 開発サーバーの実行中に Windows で `wails3 generate bindings` が「Access is denied」エラーで失敗する問題を修正しました。生成ファイルを出力ディレクトリ上のファイルと置き換えるために名前を変更するのではなく、出力ディレクトリへ同期するようにしました（#5561）
- Linux のフレームレスウィンドウで JS のリサイズイベントが発火しない問題と、フレームレスウィンドウにおけるスクロールバー端の検出を修正しました（#5368）
- 一時ディレクトリとインストールディレクトリが異なるボリュームにある場合、Windows のアップデーターが「invalid cross-device link」エラーで失敗する問題を修正しました（#5560）

## v3.0.0-alpha.99 - 2026-06-10

## 修正

- Vite 開発サーバーの実行中に Windows で `wails3 generate bindings` が「Access is denied」エラーで失敗する問題を修正しました。生成ファイルを出力ディレクトリ上のファイルと置き換えるために名前を変更するのではなく、出力ディレクトリへ同期するようにしました（#5515）

## v3.0.0-alpha.98 - 2026-06-03

## 修正

- アイドル時（インスペクターを開いている場合など）に Linux で WebKit UI がフリーズする問題を修正しました。JavaScriptCore の GC スレッド同期を破壊していた、`SIGUSR1` への `SA_ONSTACK` の強制を廃止しました（#5527）

## v3.0.0-alpha.97 - 2026-05-31

## 追加

- デバッグページと `runtime/trace` の操作方法を追加しました

## 変更

- 不要な `_ "embed"` インポートをいくつか削除し、コードを多少整理しました

## 修正

- Windows でウィンドウの最大化を解除した後、最小幅と最小高さの制約が適用されない問題を修正しました（#4593）
- Frameless と Transparent のウィンドウオプションを使用した全画面モードで、マウスクリックが背後へ抜ける問題を修正しました（#4408）

## v3.0.0-alpha.96 - 2026-05-25

## 追加

- Garble 難読化のサポート（[#4563](https://github.com/wailsapp/wails/issues/4563)）を追加しました。安定したバインディングメソッド ID、ビルドおよび Taskfile の連携（`build --obfuscated --garbleargs`、`generate bindings -obfuscated`）、ならびにランタイム向けのすべてのペイロードへの JSON 構造体タグ（`EnvironmentInfo`、`OSInfo`、`Screen`、`Rect`、`Point`、`Size`、`Capabilities`）により、Garble がエクスポートされたフィールドの名前を変更してもワイヤーフォーマットが維持されます。

## v3.0.0-alpha.95 - 2026-05-20

## 追加

- 欠落していたプロジェクト構造のページを追加しました

## 変更

- ドキュメント：アーキテクチャページの図をいくつかシーケンス図に変更し、表示を見やすくしました
- ドキュメント：実行の前提条件として D2 のインストールが必要である旨を追記しました

## 修正

- GTK4 がデフォルトの場合の `wails3 generate appimage` を修正しました。バンドラーはランタイムファイルを検索する前にバイナリから GTK スタックを検出するようになり、GTK4 ビルドでは `libwebkitgtkinjectedbundle.so`（`webkitgtk-6.0/` 配下）、`-tags gtk3` ビルドでは `libwebkit2gtkinjectedbundle.so`（`webkit2gtk-4.1/` 配下）を選択します。`.relr.dyn` のプローブでも `libgtk-4.so.1` を確認するようになり、スタックにかかわらず最新のツールチェーンではストリップが正しく無効化されます。（#5475）
- 相対 `-builddir` を指定して呼び出した場合に `wails3 generate appimage` が失敗する問題を修正しました。バンドラーは `-binary`、`-icon`、`-desktopfile`、`-builddir`、`-outputdir` を事前に絶対パスへ解決するようになり、処理途中の `s.CD` によって AppRun のダウンロード用 goroutine やコピー後の `ldd` プローブが壊れなくなりました。
- デスクトップの `Name=` フィールドがバイナリのベース名と一致しない場合、`wails3 generate appimage` が最終的な AppImage を `-outputdir` に移動できない問題を修正しました。バンドラーは linuxdeploy の appimage プラグインに対し、デスクトップファイルから導出した名前ではなく `<binary>-<arch>.AppImage` に AppImage を書き込むよう、`OUTPUT` 環境変数を介して強制するようになりました。
- alpha.93 で GTK4 + WebKitGTK 6.0 スタックがデフォルトに昇格した後、Linux で `events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`、`Common.SystemDidWake` が発火しない問題を修正しました。新しいデフォルトの `application_linux.go` `run()` は、`Linux.*` イベントを対応する `Common.*` へ転送する `setupCommonEvents()` も、`monitorPowerEvents()` も呼び出していませんでした。DBus 電源監視ヘルパーは、`application_linux_dbus.go` を介して GTK3 と GTK4 のビルドパス間で共有されるようになりました。（#5474）

## v3.0.0-alpha.94 - 2026-05-19

## 修正

- alpha.93 で GTK4 + WebKitGTK 6.0 スタックがデフォルトに昇格した後、Linux で `events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`、`Common.SystemDidWake` が発火しない問題を修正しました。新しいデフォルトの `application_linux.go` `run()` は、`Linux.*` イベントを対応する `Common.*` へ転送する `setupCommonEvents()` も、`monitorPowerEvents()` も呼び出していませんでした。DBus 電源監視ヘルパーは、`application_linux_dbus.go` を介して GTK3 と GTK4 のビルドパス間で共有されるようになりました。（#5474）

## v3.0.0-alpha.93 - 2026-05-17

## 追加

- @leaanthony により、Linux の `wails3 doctor` 出力に `XDG_SESSION_TYPE` を追加しました

## 修正

- appmenu-gtk-module が未実体化のウィンドウへアクセスすることで発生していた、Wayland 上のウィンドウメニューのクラッシュを @leaanthony が修正しました（#4769）
- アプリ名に無効な文字（空白、括弧など）が含まれる場合に GTK アプリケーションがクラッシュする問題を @leaanthony が修正しました
- Windows でドラッグ＆ドロップを初期化する際に発生する「not enough memory」エラーを @overlordtm が修正しました（#4701）
- mainthread のコールバックストアで、マップからの削除に誤って RLock を使用していたことによる競合状態（Linux、macOS、iOS）を @leaanthony が修正しました（#4424）
- タスクへコマンドライン引数を渡す際の変数処理を修正しました。KEY=VALUE のペアとして指定された CLI 変数が正しく初期化され、タスクの実行全体に伝播されるようになりました。
- macOS における NSWindowZoomButton の競合を修正しました。`MaximiseButtonState` と `FullscreenButtonState` は、起動時と実行時の両方で、より制限の厳しい状態を適用するようになりました。どちらのセッターも、もう一方の設定を暗黙に上書きできなくなりました（#5319）
- #5463 で CodeRabbit により明らかになった、従来の GTK3 ビルドパス（`-tags gtk3`）に以前から存在していた一連のバグを修正しました。ファイル関連付けから起動した場合でも起動ハンドラーをスキップしなくなりました。`getTheme` は境界と型の両面で安全になりました。`appName` は GLib が所有するメモリを解放しなくなりました。`clipboardGet` は GTK が返す `gchar*` をリークしなくなりました。`Calloc` はポインターレシーバーを使用するようになり、`NewCalloc` は `*Calloc` を返すため、プールが実際に割り当てを追跡するようになりました。`zoomOut` は、1.0 にクランプされる負の乗数ではなく、`zoomInFactor` の逆数を使用するようになりました。`execJS` は、呼び出すたびに `C.CString("")` をリークせず、事前割り当て済みの空のワールド名を再利用します。`menuItem.setAccelerator` から開発用の `fmt.Println` を削除しました。#5465 を解決しました。
- デフォルトの GTK4 ビルドパス（`linux_cgo.go`）でも、`Calloc` の値レシーバーによる同じリークを修正しました。ポインターレシーバーと `NewCalloc() *Calloc` を使用することで、ウィンドウごとの `c.String(...)` の割り当てが実際に追跡され、解放されるようになりました。

## v3.0.0-alpha.92 - 2026-05-15

## 追加

- `PACKAGE_MANAGER` オプションで使用するフロントエンドパッケージマネージャーを制御できるように Taskfile を変更しました
- Taskfile テンプレートをより予測可能に記述できるよう、テンプレートデータに `{{.Opn}}` と `{{.Cls}}` を追加しました

## 変更

- 既存の Taskfile の一部を、`{{.Opn}} and {{.Cls}}` を使用するように変更しました

## 修正

- パネルがトレイメニューを読み取っている間にメニューが更新された場合に、`linuxSystemTray` で発生する `concurrent map read and map write` ランタイム致命的エラーを修正しました。
- WebView2 のエラーとスタックトレースの出力に `fmt` ではなく `log` を使用し、Windows でコンソールを接続せずにアプリを実行した場合でもメッセージが失われないようにしました。

## v3.0.0-alpha.91 - 2026-05-12

## 変更

- `@github-actions[bot]` により、[PR](https://github.com/wailsapp/wails/pull/5414) でスポンサー SVG を更新しました
- <strong>破壊的変更（macOS）：</strong>macOS の座標系を正規化し、`GetScreens`、`Position`、`SetPosition` がすべて同じ空間を使用するようにしました。この空間は論理ポイント単位で Y 軸が下向きとなり、プライマリ画面の左上を `(0,0)` とします。これは Windows、GTK、および Electron と Web の公開 API に一致します。プライマリ画面より物理的に上にある画面では、`Bounds.Y` が負の値として報告されるようになりました（以前は正の値でした）。また、`Position()`／`SetPosition()` の値は `points × primaryScale` ではなく論理ポイント単位になりました。`Position()` → `SetPosition()` の往復変換は維持されますが、以前のアルファビルドで記録した絶対値や、手動で算出した回避策（たとえば `primaryScale` を乗算する、または画面の高さを基準に Y 軸を反転する処理）は更新が必要です。[#5117](https://github.com/wailsapp/wails/issues/5117) を解決しました。

## 修正

- パニックを防ぐため、DBus シグナル名と本文の長さを防御的に検証するようにしました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5416)）
- Linux の GTK メニュー処理におけるメモリ安全性の問題を修正しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5363)）
- NVIDIA GPU を検出し、Linux で DMA-BUF レンダラーを無効にするようにしました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5295)）
- macOS における `SetPosition` の画面間 Y 座標変換を修正しました。プライマリ画面の高さをグローバルな基準として使用することで、プライマリディスプレイから垂直方向にずれたモニターでもウィンドウが正しい位置に配置されるようになりました（[#5117](https://github.com/wailsapp/wails/issues/5117)）
- Git の PR テンプレートが正しいフィードバック URL を参照するように修正しました（@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5109)）
- 1 個ではなく 4 個の引数を渡していた壊れた `DestroyMenu` システムコールによって発生する、Windows システムトレイの一連の `SetMenu` クラッシュを修正しました。この不具合により、すべての呼び出しが FALSE を返し、何も解放されていませんでした。また、メニューの再構築時に HMENU と HBITMAP のハンドル（`MenuItem.SetBitmap` によって実行時に割り当てられたものを含む）を解放し、`Win32Menu.Update` に残った古いチェックボックス／ラジオボタンのマップをリセットするようにしました。さらに、割り当てを二重にしていた `systemtray.updateMenu` 内の冗長な `Update()` 呼び出しを削除しました。長時間稼働するシステムトレイアプリで、メニューを再構築するたびに GDI／USER オブジェクトがリークすることはなくなりました。

## v3.0.0-alpha.90 - 2026-05-11

## 追加

- macOS の WKWebView User-Agent に設定可能なアプリケーション名を追加しました（@vinhvoit225 による [PR](https://github.com/wailsapp/wails/pull/5261)）
- gin-service のサンプルに間接依存関係 github.com/coder/websocket を追加しました（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5400)）
- ビルドアセットのテストにディープイコール比較のサポートを追加しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5402)）

## 変更

- ビルド出力を assets ディレクトリに統合しました（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5401)）
- `@github-actions[bot]` により、[PR](https://github.com/wailsapp/wails/pull/5399) でスポンサー SVG を更新しました

## 修正

- macOS のシングルインスタンスメッセージに通知オブジェクトを使用するようにしました（@overlordtm による [PR](https://github.com/wailsapp/wails/pull/5289)）
- 高負荷時の Promise 消失を防ぐため、Windows のコールバックをバッチ処理するようにしました（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5383)）

## v3.0.0-alpha.89 - 2026-05-10

## 追加

- Go のテスト結果を集約する go<em>test</em>results ジョブを追加しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5316)）

## 変更

- 大きな RPC ペイロードを条件に応じて分割し、チャンク化した POST リクエストとして送信するようにしました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5369)）
- すべてのフロントエンドテンプレートで Vite を 5.x.x から 8.0.0 にアップグレードしました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5386)）
- Vite 開発サーバーのポート設定を環境変数に移行しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5365)）
- すべてのテンプレートで、Vite 開発サーバーが 127.0.0.1 にバインドするように設定しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5361)）
- `@github-actions[bot]` により、[PR](https://github.com/wailsapp/wails/pull/5384) でスポンサー SVG を更新しました

## 修正

- build-assets の更新時に Info.plist テンプレートのスタブをサニタイズするようにしました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5312)）
- メニュー項目のミューテーター（`setMenuItemChecked()`、`setMenuItemLabel()`、`setMenuItemDisabled()`、`setMenuItemHidden()`、`setMenuItemTooltip()`）をメインスレッドで同期的に適用することで、macOS メニューに古い状態が残る問題を修正しました。これにより、メニューをすぐに開き直すと以前の状態で描画される原因となっていた `dispatch_async` の競合を解消しました（#5002）
- 不要な再ビルドを防ぐため、開発モードで `*_test.go` ファイルを無視するようにしました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5203)）
- アプリが実行されていない場合に Menu.Update() でセグメンテーションフォールトが発生しないようにしました（@wucm667 による [PR](https://github.com/wailsapp/wails/pull/5291)）
- Windows でメニューバーを再描画するかどうかの判定に lastSizeWParam を使用するようにしました（@taliesin-ai による [PR](https://github.com/wailsapp/wails/pull/5382)）

## v3.0.0-alpha.88 - 2026-05-09

## 変更

- HiddenOnTaskbar で WS<em>EX</em>TOOLWINDOW を使用するように変更しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5371)）
- 依存関係の順序を変更し、go.mod から webview2 の replace ディレクティブを削除しました（@atterpac による [PR](https://github.com/wailsapp/wails/pull/5370)）
- `@github-actions[bot]` により、[PR](https://github.com/wailsapp/wails/pull/5358) でスポンサー SVG を更新しました

## 修正

- ジェネリックな間接参照エイリアスを削除し、マップのキー型を統合しました（@fbbdev による [PR](https://github.com/wailsapp/wails/pull/5331)）

## 削除

- PR-master ワークフローを削除し、ドキュメント、Go テスト、スキップテストを廃止しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5377)）

## v3.0.0-alpha.87 - 2026-05-07

## 追加

- Wails v3 の韓国語ドキュメントを追加しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5352)）
- インストールとクイックスタートのフランス語ドキュメントを追加しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5354)）
- クイックスタート、概念、コミュニティのポルトガル語ドキュメントを追加しました（@leaanthony による [PR](https://github.com/wailsapp/wails/pull/5355)）

## v3.0.0-alpha.86 - 2026-05-06

## 追加

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5328)でフランス語版ドキュメントを追加
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5343)でドキュメントサイトにドイツ語ロケールを追加

## 変更

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5347)で、翻訳済みの全8ロケールをドキュメント設定に登録
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5317)で、WebView2 用の各種 Windows 関連ファイルを更新

## 修正

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5340)で、Linux のダイアログ処理を GTK3 と GTK4 に分離
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5339)で、ダイアログのコールバックが GTK スレッド上で実行されるようにし、セグメンテーションフォルトを修正

## v3.0.0-alpha.85 - 2026-05-05

## 追加

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5179)で、PR テンプレートの URL をリポジトリに追加
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5330)で、Wails v3のドイツ語版ドキュメントを追加

## v3.0.0-alpha.84 - 2026-05-03

## 追加

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5307)で、macOS において Escape キーによるフルスクリーン終了を無効化するオプションを追加
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5310)で、macOS において Escape キーによるフルスクリーン終了を無効化するオプションを追加
- @yuseferi による[PR](https://github.com/wailsapp/wails/pull/5288)で、コミュニティショーケースに Pausa のドキュメントを追加

## 変更

- `@github-actions[bot]`による[PR](https://github.com/wailsapp/wails/pull/5308)で、スポンサーの SVG を更新
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5309)で、未対応プラットフォームを処理できるようアイコン生成コマンドを更新
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5224)で、真偽値のフルスクリーン API を3状態の ButtonState に置き換え、プラットフォームバインディングを実装

## 修正

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5315)で、WebView2 のフォーカス操作を nil のコントローラー状態から保護
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5313)で、PR のベースブランチを正しく参照するよう GitHub Actions ワークフローを更新
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5203)で、不要な再ビルドを防ぐため、開発モードでは`*_test.go`ファイルを無視
- @wucm667 による[PR](https://github.com/wailsapp/wails/pull/5291)で、アプリが実行されていない場合に Menu.Update() で発生するセグメンテーションフォルトを防止

## v3.0.0-alpha.83 - 2026-05-02

## 追加

- @symball による[PR](https://github.com/wailsapp/wails/pull/5094)で、マシン単位／ユーザー単位のインストール向けに InstallScope フラグとビルドオプションを追加
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5294)で、Window インターフェースを満たすため、何も実行しない SetScreen メソッドを BrowserWindow に追加

## 修正

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5295)で、NVIDIA GPU を検出し、Linux 上で DMA-BUF レンダラーを無効化
- @wayneforrest による[PR](https://github.com/wailsapp/wails/pull/5109)で、正しいフィードバック URL を参照するよう Git の PR テンプレートを修正
- 1個ではなく4個の引数を渡していた壊れた`DestroyMenu`システムコールにより、すべての呼び出しが FALSE を返し、何も解放されなかったことが原因の、Windows システムトレイにおける一連の`SetMenu`クラッシュを修正。さらに、メニューの再構築時に HMENU および HBITMAP ハンドル（実行時に`MenuItem.SetBitmap`を介して割り当てられたものを含む）を解放し、`Win32Menu.Update`内の古いチェックボックス／ラジオボタンのマップをリセットし、割り当てを二重にしていた`systemtray.updateMenu`内の冗長な`Update()`呼び出しを削除。長時間実行されるシステムトレイアプリで、メニューを再構築するたびに GDI/USER オブジェクトがリークしなくなりました。

## v3.0.0-alpha.82 - 2026-05-01

## 修正

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5232)で、デスクトップ名を正しく処理するようデスクトップファイル生成を修正

## v3.0.0-alpha.81 - 2026-04-30

## 変更

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5286)で、ナイトリーリリースのスケジュールを UTC 15:00に調整

## 修正

- Retina Mac で Screen Bounds、WorkArea、Size の値が半分になる問題を修正 -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## 変更

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5285)で、ドキュメントの依存関係とコンテンツコレクションローダーを更新

## v3.0.0-alpha.79 - 2026-04-29

## 追加

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5270)で、trigger-release ジョブに actions: write 権限を付与

## 変更

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5283)で、リリースタスクのデフォルトを master ブランチに変更し、変更履歴の文言を更新
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5282)で、最新バージョンを使用するよう auto-changelog ワークフローを更新
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5280)で、パスフィルターを追加し、使用されていないワークフローを削除して、ワークフローの効率を改善
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5274)で、サンプルへのリンクが master ブランチを参照するようドキュメントを更新
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5272)で、v3向けのドキュメントとサンプルを更新

## 修正

- @AkagiYui による[PR](https://github.com/wailsapp/wails/pull/5265)で、開発時のリトライロジックと IPv4 強制使用によりリバースプロキシを強化
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5281)で、未リリースの変更履歴を生成するトリガーワークフローを書き直し

## 削除

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5267)で、各種テスト用途のシェルテストスクリプトを削除
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5266)で、v3-alpha ドキュメントのデプロイワークフローと CNAME レコードを削除

### 追加

- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5196)で、サイドバーナビゲーションに「フロントエンドルーティング」の項目を追加
- @leaanthony による[PR](https://github.com/wailsapp/wails/pull/5185)で、フレームワーク別の推奨事項を含むフロントエンドルーティングガイドを追加
- モーダルシート（macOS）のサポートを追加
- Apple デバイスのサポートを改善するため、@leaanthony が ghw のバージョンを更新 (#4977)
- dock サービスに`GetBadge`メソッドを追加
- カスタム Go ビルドタグ（例：`wails3 build -tags gtk4`）を渡すため、`wails3 build`コマンドに`-tags`フラグを追加 (#4957)
- バインディングジェネレーターによる列挙型の自動生成について、専用の「列挙型」ページとサイドバーナビゲーションを含むドキュメントを追加 (#4972)
- カスタム Go ビルドタグ（例：`wails3 build -tags gtk4`）を渡すため、`wails3 build`コマンドに`-tags`フラグを追加 (#4957)
- Storage（localStorage、sessionStorage、IndexedDB、Cache API）、Network（Fetch、WebSocket、XMLHttpRequest、EventSource、Beacon）、Media（Canvas、WebGL、Web Audio、MediaDevices、MediaRecorder、Speech Synthesis）、Device（Geolocation、Clipboard、Fullscreen、Device Orientation、Vibration、Gamepad）、Performance（Performance API、Mutation Observer、Intersection/Resize Observer）、UI（Web Components、Pointer Events、Selection、Dialog、Drag and Drop）など、41種類のブラウザー API を実演する Web API サンプルを`v3/examples/web-apis/`に追加
- プラットフォーム間で200種類以上のブラウザー API をテストする WebView API 互換性チェッカーのサンプル（`v3/examples/webview-api-check/`）を追加
- 並列検索とキャッシュに対応し、Flatpak/Snap/Nix をサポートする、Linux 上のネイティブライブラリーパス検索用`internal/libpath`パッケージを追加
- <strong>作業中：</strong>Linux 向けに実験的な WebKitGTK 6.0 / GTK4 サポートを追加。`-tags gtk4`で利用可能（GTK3/WebKit2GTK 4.1は引き続きデフォルト）
- 注：タイル型ウィンドウマネージャー（Hyprland、Sway など）では、ウィンドウのジオメトリを WM が制御するため、最小化／最大化操作が期待どおりに動作しない場合があります
- @AbdelhadiSeddar により、ドキュメントの**JavaScript でイベントをリッスンする**に**一回限りのハンドラー**の使用方法を追加
- @leaanthony により、Windows/Linux 上のウィンドウが`app.Menu.Set()`で設定したアプリケーションメニューを継承できるようにする`UseApplicationMenu`オプションを`WebviewWindowOptions`に追加
- @wimaha により、Liquid Glass アイコンとアセットカタログ（macOS）の生成で`.icon`ファイル（Apple Icon Composer 形式）を使用するためのサポートを追加（#4934）
- ヘッドレス／Web デプロイ向けの実験的なサーバーモード（`-tags server`）を追加。ネイティブ GUI 依存関係なしで Wails アプリを HTTP サーバーとして実行できます。`wails3 task build:server`でビルドします。詳細は`examples/server`を参照してください。
- 並列検索とキャッシュに対応し、Flatpak/Snap/Nix をサポートする、Linux 上のネイティブライブラリーパス検索用`internal/libpath`パッケージを追加
- @leaanthony により、macOS の Spaces とフルスクリーンにまたがるウィンドウの動作を制御する`CollectionBehavior`オプションを`MacWindow`に追加（#4756）
- @leaanthony により pkg/application の単体テストを追加
- @leaanthony により MSIX パッケージングにカスタムプロトコルのサポートを追加
- Linux にデスクトップ環境の検出機能を追加（[PR #4797](https://github.com/wailsapp/wails/pull/4797)）
- @leaanthony により、フロントエンドから印刷ダイアログを開くための`Window.Print()`メソッドを JavaScript ランタイムに追加（#4290）
- @leaanthony により、Linux 上の`wails3 doctor`出力に`XDG_SESSION_TYPE`を追加
- @leaanthony により、Linux 向けに WebKit2 の追加の読み込み状態変更イベント`WindowLoadStarted`、`WindowLoadRedirected`、`WindowLoadCommitted`、`WindowLoadFinished`を追加（#3896）
- @leaanthony により、Linux 上の`wails3 doctor`出力に`XDG_SESSION_TYPE`を追加
- パッケージング時だけでなく、Linux のビルド時にも`.desktop`ファイルを生成（#4575）
- @leaanthony により、ディストリビューション固有のパッケージ名と nfpm パッケージングの例を含む Linux ランタイム依存関係のドキュメントを追加（#4339）
- @leaanthony により、Linux 上の`wails3 doctor`出力に NVIDIA ドライバーのバージョン情報を追加
- @APshenkin により、生メッセージハンドラーにオリジンを追加（[PR](https://github.com/wailsapp/wails/pull/4710)）
- @APshenkin により、macOS のユニバーサルリンクをサポート（[PR](https://github.com/wailsapp/wails/pull/4712)）
- @APshenkin により、バインディングのトランスポート層をリファクタリング（[PR](https://github.com/wailsapp/wails/pull/4702)）
- @chinenual により、サンプルアプリを Appium テストクライアントで容易にテストできるよう、helloworld テンプレートに aria-label 識別子を追加（[PR](https://github.com/wailsapp/wails/pull/4760)）
- @APshenkin により、生メッセージハンドラーにオリジンを追加（[PR](https://github.com/wailsapp/wails/pull/4710)）
- @APshenkin により、macOS のユニバーサルリンクをサポート（[PR](https://github.com/wailsapp/wails/pull/4712)）
- @APshenkin により、バインディングのトランスポート層をリファクタリング（[PR](https://github.com/wailsapp/wails/pull/4702)）
- @fbbdev と @ianvs による型付きイベント（[#4633](https://github.com/wailsapp/wails/pull/4633)）
- ツールチップをリアルタイムで更新するヘッドレストレイを示す`systray-clock`サンプルを追加（#4653）。
- @Tolfx により、Windows 向けの NSIS Protocol テンプレートを追加（#4510）
- @Tolfx により、build-assets のテストを追加（#4510）
- macOS：@nidib により、メニューバーにネイティブのウィンドウコントロールを表示（[#4588](https://github.com/wailsapp/wails/pull/4588)）
- @popaprozac により、Dock のアプリアイコンを表示／非表示にする macOS Dock サービスを追加（[PR](https://github.com/wailsapp/wails/pull/4451)）
- @popaprozac により、Dock のアプリアイコンを表示／非表示にする macOS Dock サービスを追加（[PR](https://github.com/wailsapp/wails/pull/4451)）
- @leaanthony により、NSGlassEffectView（macOS 15.0以降）と NSVisualEffectView フォールバックを使用した macOS ネイティブの Liquid Glass エフェクトをサポートし、包括的なマテリアルカスタマイズオプションを追加（[#4534](https://github.com/wailsapp/wails/pull/4534)）
- @leaanthony によるブラウザー URL のサニタイズ（[#4500](https://github.dev/wailsapp/wails/pull/4500)）。@APShenkin による[#4484](https://github.com/wailsapp/wails/pull/4484)を基に実装。
- [@Taiterbase](https://github.com/Taiterbase)による元の作業（この[PR](https://github.com/wailsapp/wails/pull/4241)）を基に、[@leaanthony](https://github.com/leaanthony)が Windows/Mac にコンテンツ保護を追加
- @leaanthony により、`wails3 build`および`wails3 package`エイリアスを介して CLI 変数を Task コマンドに渡す機能を追加（#4422、[PR](https://github.com/wailsapp/wails/pull/4488)）
- [@atterpac](https://github.com/atterpac)により、イベントの発生元となったドロップ要素のデータを取得できるドロップゾーンをサポート（[#4318](https://github.com/wailsapp/wails/pull/4318)）
- WebView2 ブラウザーに追加のコマンドライン引数を渡せるよう、`WindowsWindow`のオプションに`AdditionalLaunchArgs`を追加（[PR](https://github.com/wailsapp/wails/pull/4467)）
- [@triadmoko](https://github.com/triadmoko)により、wails init の後に go mod tidy を自動実行する機能を追加（[PR](https://github.com/wailsapp/wails/pull/4286)）
- @leaanthony による Windows Snapassist 機能（[PR](https://github.dev/wailsapp/wails/pull/4463)）
- WebView2 ブラウザーに追加のコマンドライン引数を渡せるよう、`WindowsWindow`のオプションに`AdditionalLaunchArgs`を追加（[PR](https://github.com/wailsapp/wails/pull/4467)）
- [@triadmoko](https://github.com/triadmoko)により、wails init の後に go mod tidy を自動実行する機能を追加（[PR](https://github.com/wailsapp/wails/pull/4286)）
- @leaanthony による Windows Snapassist 機能（[PR](https://github.dev/wailsapp/wails/pull/4463)）
- [@almas-x](https://github.com/almas-x)により、Windows 向け`getAccentColor`実装を追加（[PR](https://github.com/wailsapp/wails/pull/4427)）
- [@almas-x](https://github.com/almas-x)により、Windows 向け`getAccentColor`実装を追加（[PR](https://github.com/wailsapp/wails/pull/4427)）
- Windows のダークテーマ対応メニューとメニューバー。@leaanthony による変更（[a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)）
- @popaprozac により、JS/TS バインディングを明確にするため組み込みサービスの名前を変更（[PR](https://github.com/wailsapp/wails/pull/4405)）
- ユーザーのシステムのアクセントカラーを取得する`app.Env.GetAccentColor`。MacOS で動作。[@etesam913](https://github.com/etesam913)による実装
- [@atterpac](https://github.com/atterpac)により、`window.ToggleFrameless()` API を追加（[#4137](https://github.com/wailsapp/wails/pull/4137)）
- @leaanthony により、Linux 向けにディストリビューション固有のビルド依存関係を追加（[PR](https://github.com/wailsapp/wails/pull/4345)）
- @atterpac がバインディングガイドを追加（[PR](https://github.com/wailsapp/wails/pull/4404)）
- **テスト基盤の整理**：[@leaanthony](https://github.com/leaanthony) が Docker テストファイルを、最適化されたイメージを備える専用の `test/docker/` ディレクトリへ移動し、ビルドの信頼性を向上（[#4359](https://github.com/wailsapp/wails/pull/4359)）
- **リソース管理パターンの改善**：[@leaanthony](https://github.com/leaanthony) がサンプルに、イベントハンドラーの適切なクリーンアップとコンテキストを考慮した goroutine 管理を追加（[#4359](https://github.com/wailsapp/wails/pull/4359)）
- [@AkshayKalose](https://github.com/AkshayKalose) が aarch64 AppImage ビルドをサポート（[#3981](https://github.com/wailsapp/wails/pull/3981)）
- [@leaanthony](https://github.com/leaanthony) が `wails doctor` に診断セクションを追加
- [@leaanthony](https://github.com/leaanthony) がサービスメソッドの呼び出し時にウィンドウをコンテキストへ追加
- [@leaanthony](https://github.com/leaanthony) が、どのウィンドウからサービスが呼び出されたかを判別する方法を示す `window-call` サンプルを追加
- [@leaanthony](https://github.com/leaanthony) による新しいメニューガイド
- [@leaanthony](https://github.com/leaanthony) が panic 処理を改善
- [@leaanthony](https://github.com/leaanthony) による新しいメニューガイド
- [@fbbdev](https://github.com/fbbdev) が Service API のドキュメントコメントを追加（[#4024](https://github.com/wailsapp/wails/pull/4024)）
- [@leaanthony](https://github.com/leaanthony) が、追加設定を使用してサービスを初期化する関数 `application.NewServiceWithOptions` を追加（[#4024](https://github.com/wailsapp/wails/pull/4024)）
- [@FalcoG](https://github.com/FalcoG) と [@leaanthony](https://github.com/leaanthony) がメニュー制御を改善（[#4031](https://github.com/wailsapp/wails/pull/4031)）
- [@leaanthony](https://github.com/leaanthony) がドキュメントをさらに追加
- [@leaanthony](https://github.com/leaanthony) が標準イベントリスナーでのイベントキャンセルをサポート
- [@leaanthony](https://github.com/leaanthony) がシステムトレイの `Hide`、`Show`、`Destroy` をサポート
- [@leaanthony](https://github.com/leaanthony) がシステムトレイの `SetTooltip` をサポート。原案は [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- [@fbbdev](https://github.com/fbbdev) が、サポートされていない型に関するバインディングジェネレーターの警告にパッケージパスを表示（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- [@fbbdev](https://github.com/fbbdev) がバインディングジェネレーターでジェネリックエイリアスをサポート（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- [@fbbdev](https://github.com/fbbdev) がバインディングジェネレーターで `omitzero` JSON フラグをサポート（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- [@fbbdev](https://github.com/fbbdev) が、選択したサービスメソッドのバインディング生成を抑止する `//wails:ignore` ディレクティブを追加（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- [@fbbdev](https://github.com/fbbdev) が、Go ではエクスポートされるが JS/TS ではエクスポートされない型を使用できるよう、サービスとモデルに `//wails:internal` ディレクティブを追加（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- [@fbbdev](https://github.com/fbbdev) が、弱い型付けの列挙型を使用できるよう、バインディングジェネレーターでエイリアス型の定数をサポート（[#4045](https://github.com/wailsapp/wails/pull/4045)）
- [@fbbdev](https://github.com/fbbdev) が Go の 1.24 機能に関するバインディングジェネレーターのテストを追加（[#4068](https://github.com/wailsapp/wails/pull/4068)）
- [#4065](https://github.com/wailsapp/wails/pull/4065) で OS バージョンの検出を改善するため、`OSInfo.Branding` に macOS 15「Sequoia」のサポートを追加
- [@fbbdev](https://github.com/fbbdev) が、シャットダウン処理の完了後にカスタムコードを実行するための `PostShutdown` フックを追加（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- [@fbbdev](https://github.com/fbbdev) が、カスタムエラーハンドラーで致命的エラーを検出できるようにする `FatalError` 構造体を追加（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- [@fbbdev](https://github.com/fbbdev) がサービスの起動順序とシャットダウン順序を標準化し、文書化（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- [@fbbdev](https://github.com/fbbdev) が、アプリケーションの起動／シャットダウンシーケンス用テストハーネスと、サービスの起動／シャットダウンテストを追加（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- [@fbbdev](https://github.com/fbbdev) が、アプリケーション作成後にサービスを登録するための `RegisterService` メソッドを追加（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- [@fbbdev](https://github.com/fbbdev) が、バインディング呼び出しのカスタムエラー処理用に、アプリケーションとサービスのオプションへ `MarshalError` フィールドを追加（[#4066](https://github.com/wailsapp/wails/pull/4066)）
- [@fbbdev](https://github.com/fbbdev) が、キャンセル要求を promise チェーン全体へ伝播するキャンセル可能な promise ラッパーを追加（[#4100](https://github.com/wailsapp/wails/pull/4100)）
- [@fbbdev](https://github.com/fbbdev) が、バインディング呼び出しのキャンセルを `AbortSignal` に関連付ける機能を追加（[#4100](https://github.com/wailsapp/wails/pull/4100)）
- [@leaanthony](https://github.com/leaanthony) が、通常の `wml-*` 属性に加えて WML の `data-wml-*` 属性をサポート
- [@fbbdev](https://github.com/fbbdev) が、後からの設定や動的な再設定を行うための `Configure` メソッドをすべてのサービスに追加（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) により、`fileserver` サービスは未設定の場合に 503 Service Unavailable レスポンスを送信（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) により、`kvstore` サービスは未設定の場合、デフォルトでインメモリのキーバリューストアを提供（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) が、設定変更後にファイルからデータを再読み込みする `Load` メソッドを `kvstore` サービスに追加（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) が、すべてのキーを削除する `Clear` メソッドを `kvstore` サービスに追加（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) が、JS 側にログレベル定数を提供する型 `Level` を `log` サービスに追加（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) が、ログレベルを動的に指定する `Log` メソッドを `log` サービスに追加（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) により、`sqlite` サービスは未設定の場合、デフォルトでインメモリ DB を提供（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) が、DB を手動で閉じるメソッド `Close` を `sqlite` サービスに追加（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) が `sqlite` サービスのクエリメソッドでキャンセルをサポート（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@fbbdev](https://github.com/fbbdev) が、JS バインディングを備えたプリペアドステートメントのサポートを `sqlite` サービスに追加（[#4067](https://github.com/wailsapp/wails/pull/4067)）
- [@AnalogJ](https://github.com/AnalogJ) によるこの [PR](https://github.com/wailsapp/wails/pull/3537) での元の作業を基に、[Lea Anthony](https://github.com/leaanthony) が [PR](https://github.com/wailsapp/wails/pull/3537) で Gin サポートを追加
- 自動保存とパスワードの自動保存が常に有効になる問題を、[@oSethoum](https://github.com/osethoum) が [#4134](https://github.com/wailsapp/wails/pull/4134) で修正
- ウィンドウにメニューを設定できるよう、[@leaanthony](https://github.com/leaanthony) がウィンドウに `SetMenu()` を追加
- [@popaprozac](https://github.com/popaprozac) が [#4098](https://github.com/wailsapp/wails/pull/4098) で通知をサポート
-  [@wimaha](https://github.com/wimaha) が [#4177](https://github.com/wailsapp/wails/pull/4177) で mac のファイル関連付けをサポート
- [@leaanthony](https://github.com/leaanthony) がセマンティックバージョンの更新用に `wails3 tool version` を追加
- [@popaprozac](https://github.com/popaprozac) が [#](https://github.com/wailsapp/wails/pull/4234) で macOS と Windows のバッジ表示をサポート
- [@fbbdev](https://github.com/fbbdev) と [@IanVS](https://github.com/IanVS) が [#4161](https://github.com/wailsapp/wails/pull/4161) で登録済み／厳密に型付けされたイベントをサポート
- [@fbbdev](https://github.com/fbbdev) と [@IanVS](https://github.com/IanVS) が [#4161](https://github.com/wailsapp/wails/pull/4161) でカスタムイベント用フックを登録する機能を追加
- [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte) による、システムのファイルマネージャーでパス `path` を開き、必要に応じて `selectFile` で強調表示するための `app.OpenFileManager(path string, selectFile bool)`
- [@leaanthony](https://github.com/leaanthony) が `wails3 init` コマンドに新しい `-git` フラグを追加
- [@leaanthony](https://github.com/leaanthony) による新しい `wails3 generate webview2bootstrapper` コマンド
- [@leaanthony](https://github.com/leaanthony) がランタイムを手動で初期化できるよう、ランタイムに `init()` メソッドを追加
- [@leaanthony](https://github.com/leaanthony) がウィンドウの WindowOptions に `WindowDidMoveDebounceMS` オプションを追加
- [@leaanthony](https://github.com/leaanthony) がシングルインスタンス機能を追加。@APshenkin による [v2 PR](https://github.com/wailsapp/wails/pull/2951) が基になっています。
- [@leaanthony](https://github.com/leaanthony) による `wails3 generate template` コマンド
- [@leaanthony](https://github.com/leaanthony) による `wails3 releasenotes` コマンド
- [@leaanthony](https://github.com/leaanthony) による `wails3 update cli` コマンド
- [@leaanthony](https://github.com/leaanthony) による `wails3 generate bindings` コマンドの `-clean` オプション
- [@AkshayKalose](https://github.com/AkshayKalose) が [#3981](https://github.com/wailsapp/wails/pull/3981) で aarch64（arm64）向け AppImage Linux ビルドを可能に
- @ansxuman が [#3958](https://github.com/wailsapp/wails/pull/3958) でスポンサーへのハイパーリンクを追加
- による deb、rpm、Arch Linux パッケージャーの Linux パッケージビルドのサポート
- による Darwin ユニバーサルビルドおよびパッケージのサポートを追加
- による Web サイトへのイベントドキュメントの追加
- 非 SSR 開発向けに設定された sveltekit および sveltekit-ts のテンプレート
- による新しい `wails3 update build-assets` コマンドを使用したビルドアセットの更新
- による HTML Drag and Drop API をテストするためのサンプル
- [leaanthony](https://github.com/leaanthony) が次でファイル関連付けをサポート
- による新しい `wails3 generate runtime` コマンド
- ウィンドウを中央に配置するか、またはそれ以外にするかを指定する新しい `InitialPosition` オプション
- による `application` パッケージへの `Path` および `Paths` メソッドの追加
- Windows の `GeneralAutofillEnabled` および `PasswordAutosaveEnabled` オプションを追加
- によるサービスメソッドを呼び出したウィンドウを取得する機能の追加
- による WebView2 の `EnabledFeatures` および `DisabledFeatures` オプションの追加
- ⊞ による高 DPI モニターのサポートを強化する新しい DIP システム
- ⊞ [windom](https://github.com/windom/) が次でウィンドウクラス名オプションを追加
- サービスを拡張してプラグイン機能を提供。作成者：
- 🐧 次での WindowDidMove／WindowDidResize イベント
- ⊞ 次での WindowDidResize イベント
-  Dock を処理できるように Event ApplicationShouldHandleReopen を追加
-  @tmclane が次で実装に getPrimaryScreen/getScreens を追加
-  による macOS のフルスクリーンモードでツールバーを表示するオプションの追加
- 🐧 Linux のキー入力をアクセラレーターに変換する onKeyPress ロジックを追加
- 🐧 によるタスク `run:linux` の追加
- [@almas-x](https://github.com/almas-x) が次で `SetIcon` メソッドをエクスポート
- [@almas-x](https://github.com/almas-x) が次で `OnShutdown` を改善
- による `Window` インターフェースへの `ToggleMaximise` メソッドの復元
- @leaanthony が次で `Environment()` に詳細情報を追加
- による `Window` インターフェースでの `WebviewWindow.IsFocused` メソッドの公開
- による WML システムでのスペース区切りの複数トリガーイベントのサポート
- によるバンドル版 JS ランタイムスクリプトからの ESM エクスポートの追加
- 代わりにバンドル版 JS ランタイムスクリプトを使用するためのバインディングジェネレーターフラグを追加
- [@abichinger](https://github.com/abichinger) が Linux に `setIcon` を実装
- dev コマンドにフラグ `-port` を追加し、環境変数をサポート
- によるバインドされたメソッド呼び出しのテストの追加
- ⊞ による作成済みウィンドウへの `SetIgnoreMouseEvents` の追加
-  によるウィンドウの重なりレベル（順序）を設定する機能の追加

### 修正

- NSScreen のポイント値を `Physical*` フィールド内でデバイスピクセルに変換し、トップレベルの `Screen.X`／`Y` を設定することで、Retina Mac で `Screen.Bounds`、`WorkArea`、`Size` が半分になる問題を修正し、マルチモニター間の接触判定と作業領域への配置を正しくしました。@wayneforrest による [PR](https://github.com/wailsapp/wails/pull/5168)
- ディスプレイ構成の変更時（例：スリープ／復帰中の外部モニターのホットプラグ）に WebKit DisplayLink のデッドロックを引き起こす ScreenManager のデータ競合を修正
- Assets.car が存在する場合、CFBundleIconName を appicon に直接設定。@symball による [PR](https://github.com/wailsapp/wails/pull/5154)
- Fedora、openSUSE、Arch、NixOS で `wails3 doctor` が誤った WebKitGTK パッケージを報告する問題を修正。v3 はコンパイル時に 4.1 API を必要とするため、4.0 のフォールバック項目を削除しました（#5071）
- openSUSE の webkit2gtk に対して doctor が示すパッケージ名を修正（`webkit2gtk4_1-devel` → `webkit2gtk3-devel`。後者が正しい openSUSE パッケージ名）（#5071）
- デスクトップ開発モードで `/wails/custom.js` が見つからない場合に発生する `Unexpected token '<'` エラーを修正。HTML の SPA フォールバックが JavaScript として挿入されるのを防ぐため、`/wails/custom.js` 用の明示的な 404 ハンドラーと、大文字と小文字を区別しない `Content-Type` 検証を `loadOptionalScript` に追加しました。（[#5068](https://github.com/wailsapp/wails/issues/5068)）
- macOS でシステムトレイメニューのハイライト状態を修正。メニューを開いている間、アイコンが選択状態で表示されるようになりました（#4910）
- macOS でシステムトレイに関連付けられたウィンドウがほかのウィンドウの背後に表示される問題を修正。適切なポップアップウィンドウレベルを使用するようになりました（#4910）
- ドキュメント全体にある誤った `@wailsio/runtime` のインポート例を修正（#4989）
- darwin でフレームレスウィンドウを最小化できない問題を修正（#4294）
- `node_modules/` を go-task の最新状態チェックから除外し、`wails3 build` および `wails3 dev` の実行中に 20-30 分間停止する問題を修正。以前は `sources: "**/*"` glob により、go-task が `node_modules/` 内の全ファイル（MUI などの大規模な依存関係を含む 50000-100000 個超のファイル）を列挙してチェックサムを計算していたため、特に Windows/NTFS で処理が遅くなっていました（#4939）
- C の `Screen` typedef と X11 Xlib.h の衝突によって発生する GTK4 のビルド失敗を修正（#4957）
- macOS で Dock バッジメソッドの一貫性を修正
- `InvisibleTitleBarHeight` がフレームレスウィンドウまたはタイトルバーが透明なウィンドウだけでなく、すべての macOS ウィンドウに適用される問題を修正（#4960）
- `InvisibleTitleBarHeight` が有効な場合に上隅からサイズを変更するとウィンドウが揺れる問題を、ウィンドウ端の近くではドラッグを開始しないようにして修正（#4960）
- JS/TS バインディングで列挙型をキーとするマップ型の生成を修正（#4437、@fbbdev）
- Windows でディスプレイの拡大率が 100% 以外の場合に、ファイルのドラッグ＆ドロップが機能しない問題を修正
- Windows でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- Windows でファイルドロップの座標に誤ったピクセル空間（物理ピクセルと CSS ピクセル）が使用される問題を修正
- Linux でホバー効果を使用すると、ファイルのドラッグ＆ドロップが安定して機能しない問題を修正
- Linux でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- `gtk_window_present()` を使用し、Linux/GTK4 でウィンドウを表示または非表示にした際、最小化された状態に戻ることがある問題を修正（#4957）
- `XTranslateCoordinates`/`XMoveWindow` を介した X11 条件付きサポートを追加し、Linux/GTK4 でウィンドウ位置の取得／設定が常に 0,0 を返す問題を修正（#4957）
- 削除された `gtk_window_set_geometry_hints` の代わりにシグナルベースのサイズ制限を追加し、Linux/GTK4 でウィンドウの最大サイズが適用されない問題を修正（#4957）
- 適切な PhysicalBounds の計算と `gdk_monitor_get_scale`（GTK 4.14 以降）による小数スケーリングのサポートを実装し、Linux/GTK4 の DPI スケーリングを修正
- Linux/GTK4 で新しいウィンドウを作成するとメニュー項目が重複する問題を修正
- JS/TS バインディングで列挙型をキーとするマップ型の生成を修正（#4437、@fbbdev）
- Windows でディスプレイの拡大率が 100% 以外の場合に、ファイルのドラッグ＆ドロップが機能しない問題を修正
- Windows でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- Windows でファイルドロップの座標に誤ったピクセル空間（物理ピクセルと CSS ピクセル）が使用される問題を修正
- Linux でホバー効果を使用すると、ファイルのドラッグ＆ドロップが安定して機能しない問題を修正
- Linux でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- 適切な PhysicalBounds の計算と `gdk_monitor_get_scale`（GTK 4.14 以降）による小数スケーリングのサポートを実装し、Linux/GTK4 の DPI スケーリングを修正
- Linux/GTK4 で新しいウィンドウを作成するとメニュー項目が重複する問題を修正
- JS/TS バインディングで列挙型をキーとするマップ型の生成を修正（#4437、@fbbdev）
- App.Window.Current() でメインスレッドから AppKit API にアクセスしていなかったために macOS で発生する「ゴーストウィンドウ」の問題を修正（#4947、@wimaha）
- WKUIDelegate runOpenPanelWithParameters を実装し、macOS で HTML の `<input type="file">` が機能しない問題を修正（#4862）
- macOS/Linux で `@wailsio/runtime` npm モジュールを使用すると、ネイティブのファイルドラッグ＆ドロップが機能しない問題を修正（#4953、@leaanthony）
- パッケージをまたぐ型エイリアスのバインディング生成を修正（#4578、@fbbdev）
- GTK のスレッドセーフティ違反によって Linux で OpenFileDialog がクラッシュする問題を修正（#3683、@ddmoney420）
- 非表示または破棄済みのウィンドウで `Focus()` を呼び出すと SIGSEGV が発生する問題を修正（#4890、@ddmoney420）
- Linux で空のアイコンまたはビットマップを設定すると panic が発生する可能性がある問題を修正（#4923、@ddmoney420）
- macOS でサービスバインディングから ErrorDialog を呼び出すとクラッシュする問題を修正（#3631、@leaanthony）
- `v3\examples\dialogs` で Windows OS 上にメニューが表示されるように変更（@ndianabasi）
- ページの再読み込み中に TypeError を引き起こす競合状態を修正（#4872、@ddmoney420）
- `Collector.IsVoidAlias()` メソッドからグローバル状態を削除し、バインディングジェネレーターのテストで誤った出力が生成される問題を修正（#4941、@fbbdev）
- macOS で `<input type="file">` ファイルピッカーが機能しない問題を修正（#4862、@leaanthony）
- macOS で `Position()` と `SetPosition()` が異なる座標系を使用し、状態の保存／復元時にウィンドウ位置がずれる問題を修正（#4816、@leaanthony）
- アプリケーションマニフェストですでに DPI 認識が設定されている場合に発生する SetProcessDpiAwarenessContext の「Access is denied」エラーを修正（#4803）
- キーボードショートカットのドキュメントページを更新し、`KeyBinding.Add` のコールバック引数の型を修正（@ndianabasi）
- カスタムバインディングの生成に関するドキュメントを修正。`-o String` ではなく `-d String` を使用する必要があります
- `menu.Update()` でメニューの子項目がクリアされない問題を修正
- ドキュメント内の古い Manager API 参照を修正（31 ファイルを更新し、`app.Window.New()`、`app.Event.Emit()` などの新しいパターンを使用）（@leaanthony）
- WebKit がシグナルハンドラーを上書きするため、JS にバインドされた Go メソッド内で panic が発生した際に Linux がクラッシュする問題を修正（#3965、@leaanthony）
- Linux で SaveFileDialog.SetFilename() が機能しない問題を修正（#4841、@samstanier）
- ドラッグ＆ドロップの例でドロップ座標が undefined と表示される問題を修正
- APP_NAME に空白が含まれる場合に macOS アプリバンドルの作成が失敗する問題を修正（ブレース展開の問題）
- Windows でサービスメソッドを呼び出した際に発生するインデックス範囲外の panic を修正（goccy/go-json を元に戻しました）
- Windows でディスプレイの拡大率が 100% 以外の場合に、ファイルのドラッグ＆ドロップが機能しない問題を修正
- Windows でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- Windows でファイルドロップの座標に誤ったピクセル空間（物理ピクセルと CSS ピクセル）が使用される問題を修正
- Linux でホバー効果を使用すると、ファイルのドラッグ＆ドロップが安定して機能しない問題を修正
- Linux でファイルドロップを有効にすると、HTML5 の内部ドラッグ＆ドロップが機能しなくなる問題を修正
- `APP_NAME` などの変数に含まれる空白に対応するため、すべてのオペレーティングシステム向けの Taskfile.yml ファイル内にある全コマンドを更新（@ndianabasi）
- Linux で 'build:universal:lipo:go' タスクを実行した際に発生するコマンド引数エラーを修正（@wux1an）
- Linux で 'wails3 build GOOS=darwin GOARCH=arm64' を実行した際に発生する Docker エラー「undefined symbol: **<em>ubsan</em>handle_xxxxxxx」を修正（@wux1an）
- カスタムプロトコルのドキュメントを統合し、Universal Links のセクションを追加（@leaanthony）
- TrackPopupMenuEx の同時呼び出しを防ぐガードを追加し、アイコンを繰り返しクリックした際に Windows のシステムトレイメニューがクラッシュする問題を修正（#4151、@leaanthony）
- app.Run() より前に systray.Run() を呼び出した際にアプリがクラッシュする問題を修正（@leaanthony）
- ApplicationShouldTerminateAfterLastWindowClosed が有効な状態で Hide()/Show() によりウィンドウの表示状態を切り替えた際に、macOS でクラッシュする問題を修正（#4389、@leaanthony）
- macOS および Windows でメニューを繰り返し開いた際に発生する、コンテキストメニューのメモリリークを修正（#4012、@leaanthony）
- macOS でコンテキストメニューのネイティブリソースが再利用されず、表示のたびに新しいメニューが作成される問題を修正（#4012、@leaanthony）
- `Hidden: true` を指定して起動したアプリで、macOS の Dock アイコンをクリックしても非表示のウィンドウが表示されない問題を修正（#4583、@leaanthony）
- CGO 呼び出しでウィンドウポインターの型が誤っていたため、macOS で印刷ダイアログが開かない問題を修正（#4290、@leaanthony）
- appmenu-gtk-module が未実体化のウィンドウにアクセスすることで発生する、Wayland 上のウィンドウメニューのクラッシュを修正（#4769、@leaanthony）
- アプリ名に無効な文字（空白、丸括弧など）が含まれる場合に GTK アプリケーションがクラッシュする問題を修正（@leaanthony）
- Windows でドラッグ＆ドロップを初期化した際に発生する「not enough memory」エラーを修正（#4701、@overlordtm）
- URI のエスケープが誤っていたため、Linux でファイルエクスプローラーが別のディレクトリを開く問題を修正（#4397、@leaanthony）
- `.relr.dyn` ELF セクションを自動検出してストリップを無効化することで、最新の Linux ディストリビューション（Arch、Fedora 39 以降、Ubuntu 24.04 以降）で AppImage のビルドが失敗する問題を修正（#4642、@leaanthony）
- Fedora/DNF ベースのシステムで、`wails doctor` が webkit パッケージをインストール済みと誤って報告する問題を修正（#4457、@leaanthony）
- デフォルトの `config.yml` が本番ビルドで `wails3 dev` を実行する問題を修正（@mbaklor）
- 存在しないパッケージのインポートにより、iOS サービススタブでビルドが失敗する問題を修正（@leaanthony）
- debug/info メソッドの構造化ログで「no formatting directives」エラーが発生する問題を修正（@leaanthony）
- モバイルプラットフォームのマージから誤って含まれていた一時的なデバッグ出力文を削除（@leaanthony）
- NVIDIA GPU を搭載した Wayland 環境で WebKitGTK がクラッシュする問題（Error 71 Protocol error）を、DMA-BUF レンダラーの自動無効化により修正（@leaanthony）
- Linux の `application.WebviewWindowOptions.BackgroundColour` でアルファ値が無視される問題を解決（[#4722](https://github.com/wailsapp/wails/pull/4722)、@BradHacker）
- カスタムアイコンが指定されていない場合に、Windows のシステムトレイアイコンがデフォルトでアプリケーションアイコンにならない問題を修正（#4704）
- `HICON` の所有権を追跡し、ユーザーが作成したハンドルのみを破棄することで、Explorer による再利用時のクラッシュを防止（#4653）。
- 破棄時に Windows のシステムテーマリスナーと保持中のトレイアイコンを解放し、goroutine とデバイスコンテキストのリークを解消（#4653）。
- サロゲートペアやマルチバイトグリフの破損を防ぐため、トレイのツールチップを 127 UTF-16 単位で切り詰めるように変更（#4653）。
- Windows のパッケージタスクが失敗する問題を修正（#4667）
- Linux の Taskfile にある AppImage の appicon 変数を修正（[PR #4644](https://github.com/wailsapp/wails/pull/4644)）
- go-webview2 の v1.0.22 シグネチャ変更により発生する Windows のビルドエラーを修正（#4513、#4645）
- Linux の Taskfile にある AppImage の appicon 変数を修正（[PR #4644](https://github.com/wailsapp/wails/pull/4644)）
- `<.Info.Protocol>` から `<.Protocol>` までを削除し、Linux の desktop.tmpl にあるプロトコル範囲を修正（@Tolfx、#4510）
- liquid glass デモの再定義エラーを修正（[#4542](https://github.com/wailsapp/wails/pull/4542)、@Etesam913）
- Linux でのシステムトレイメニューの更新を修正（[#4604](https://github.com/wailsapp/wails/issues/4604)、[@JackDoan](https://github.com/JackDoan)）
- 非表示ウィンドウの作成時に Windows で白いウィンドウが表示される問題を修正（@leaanthony、[#4612](https://github.com/wailsapp/wails/pull/4612)）
- ドキュメント内の notifications パッケージのインポートパスを修正（@rxliuli、[#4617](https://github.com/wailsapp/wails/pull/4617)）
- npm パッケージ @wailsio/runtime の使用時にドラッグ＆ドロップが機能しない問題を修正（#4489、@leaanthony、#4616）
- Windows：起動時にウィンドウがちらつく問題と、非表示ウィンドウが誤って表示される問題を修正（[PR](https://github.com/wailsapp/wails/pull/4600)、@leaanthony）。
- Wayland でウィンドウを最大化した際のサイズ問題を修正（https://github.com/wailsapp/wails/issues/4429、[@samstanier](https://github.com/samstanier)）
- Wayland でウィンドウを最大化した際のサイズ問題を修正（https://github.com/wailsapp/wails/issues/4429、[@samstanier](https://github.com/samstanier)）
- liquid glass デモの再定義エラーを修正（[#4542](https://github.com/wailsapp/wails/pull/4542)、@Etesam913）
- MacOS で AssetServer がクラッシュする可能性がある問題を修正（[#4576](https://github.com/wailsapp/wails/pull/4576)、@jghiloni）
- NextJs でのビルド時に発生するコンパイルの問題を修正（[#4585](https://github.com/wailsapp/wails/pull/4585)、@rev42）
- ナイトリーリリース用のパイプラインを修正（[#4597](https://github.com/wailsapp/wails/pull/4597)、@riadafridishibly）
- liquid glass デモの再定義エラーを修正（[#4542](https://github.com/wailsapp/wails/pull/4542)、@Etesam913）
- MacOS で AssetServer がクラッシュする可能性がある問題を修正（[#4576](https://github.com/wailsapp/wails/pull/4576)、@jghiloni）
- NextJs でのビルド時に発生するコンパイルの問題を修正（[#4585](https://github.com/wailsapp/wails/pull/4585)、@rev42）
- ナイトリーリリース用のパイプラインを修正（[#4597](https://github.com/wailsapp/wails/pull/4597)、@riadafridishibly）
- liquid glass デモの再定義エラーを修正（[#4542](https://github.com/wailsapp/wails/pull/4542)、@Etesam913）
- Windows の SetBackgroundColour を修正（@PPTGamer、[PR](https://github.com/wailsapp/wails/pull/4492)）
- Manager API のリファクタリングによる変更を反映するようドキュメントを更新（@yulesxoxo、[PR #4476](https://github.com/wailsapp/wails/pull/4476)）
- Linux の Taskfile にある .desktop ファイルの appicon 変数を修正（[PR #4477](https://github.com/wailsapp/wails/pull/4477)）
- Manager API のリファクタリングによる変更を反映するようドキュメントを更新（@yulesxoxo、[PR #4476](https://github.com/wailsapp/wails/pull/4476)）
- [#4456](https://github.com/wailsapp/wails/issues/4456) で報告された Windows の nil ポインター参照解除バグを修正（@leaanthony、[#4460](https://github.com/wailsapp/wails/pull/4460)）
- 2 本指のスワイプによるナビゲーションジェスチャーを有効にするため、macOS WKWebView で `allowsBackForwardNavigationGestures` をサポート（#1857）
- 初期状態で無効に設定されたメニュー項目で onClick が機能しない問題を修正（@leaanthony、[PR #4469](https://github.com/wailsapp/wails/pull/4469)）。初期調査を行った @IanVS に感謝します。
- ビルド失敗時に Vite サーバーがクリーンアップされない問題を修正（#4403）
- Windows で `SaveFileDialog` を閉じるかキャンセルした際に発生するパニックを修正（[PR](https://github.com/wailsapp/wails/pull/4284)、@hkhere）
- WindowsでのHTMLレベルのドラッグ＆ドロップを、[#4259](https://github.com/wailsapp/wails/pull/4259)で[@mbaklor](https://github.com/mbaklor)が修正
- 2本指スワイプによるナビゲーションジェスチャーを有効にするため、macOSのWKWebViewで`allowsBackForwardNavigationGestures`をサポート（#1857）
- 最初から無効に設定されたメニュー項目でonClickが機能しない問題を、[PR #4469](https://github.com/wailsapp/wails/pull/4469)で@leaanthonyが修正。初期調査を行った@IanVSに感謝します。
- ビルド失敗時にViteサーバーがクリーンアップされない問題を修正（#4403）
- Windowsでの通知の解析を、[PR](https://github.com/wailsapp/wails/pull/4450)で@popaprozacが修正
- Windows SDKの依存関係を確認するようdoctorコマンドを、[#4390](https://github.com/wailsapp/wails/issues/4390)で[@kodumulo](https://github.com/kodumulo)が修正
- MacのprocessURLRequestにおけるnilポインター参照を、[#4366](https://github.com/wailsapp/wails/pull/4366)で[@etesam913](https://github.com/etesam913)が修正
- フィルター付きダイアログを使用できないLinuxのバグを、[#4287](https://github.com/wailsapp/wails/pull/4287)で[@bh90210](https://github.com/bh90210)が修正
- WindowsおよびLinuxの編集メニューの問題を、[#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)で[@leaanthony](https://github.com/leaanthony)が修正
- macOSの.plistファイルにある最小システムバージョンを10.13.0から10.15.0へ、[#3981](https://github.com/wailsapp/wails/pull/3981)で[@AkshayKalose](https://github.com/AkshayKalose)が更新
- ウィンドウIDが飛ぶ問題を[@leaanthony](https://github.com/leaanthony)が修正
- RegisterContextMenuの呼び出し時にメニューがnilになる問題を[@leaanthony](https://github.com/leaanthony)が修正
- バインディングジェネレーターの出力における依存関係の循環を、[#4001](https://github.com/wailsapp/wails/pull/4001)で[@fbbdev](https://github.com/fbbdev)が修正
- バインディングジェネレーターの出力における定義前使用エラーを、[#4001](https://github.com/wailsapp/wails/pull/4001)で[@fbbdev](https://github.com/fbbdev)が修正
- ビルドフラグをバインディングジェネレーターへ渡すよう、[#4023](https://github.com/wailsapp/wails/pull/4023)で[@fbbdev](https://github.com/fbbdev)が変更
- Windows以外のプラットフォームでも動作するよう、WindowsのTaskfile内のパスをスラッシュ区切りへ[@leaanthony](https://github.com/leaanthony)が変更
- Mac + Mac JSイベントを[@leaanthony](https://github.com/leaanthony)が修正
- macOSでのイベントのデッドロックを[@leaanthony](https://github.com/leaanthony)が修正
- WindowsでHTMLは指定されているもののJSが指定されていない場合に、Windowの初期化で発生する`Parameter incorrect`エラーを[@leaanthony](https://github.com/leaanthony)が修正
- アセットサーバーでコンテンツタイプの判別に使用するレスポンス接頭部のサイズを、[#4049](https://github.com/wailsapp/wails/pull/4049)で[@fbbdev](https://github.com/fbbdev)が修正
- アセットサーバーのルートインデックスパスで、404以外のレスポンスの処理を、[#4049](https://github.com/wailsapp/wails/pull/4049)で[@fbbdev](https://github.com/fbbdev)が修正
- バインディングジェネレーターでジェネリック型のプロパティを検査する際の未定義動作を、[#4045](https://github.com/wailsapp/wails/pull/4045)で[@fbbdev](https://github.com/fbbdev)が修正
- 基底型のプロパティが名前付きラッパーと一致しない場合の、モデルに対するバインディングジェネレーターの出力を、[#4045](https://github.com/wailsapp/wails/pull/4045)で[@fbbdev](https://github.com/fbbdev)が修正
- マップのキー型および前処理に対するバインディングジェネレーターの出力を、[#4045](https://github.com/wailsapp/wails/pull/4045)で[@fbbdev](https://github.com/fbbdev)が修正
- マーシャラーインターフェースを実装する構造体に対するバインディングジェネレーターの出力を、[#4045](https://github.com/wailsapp/wails/pull/4045)で[@fbbdev](https://github.com/fbbdev)が修正
- バインディングジェネレーターで、ジェネリック型を含む型循環の検出を、[#4045](https://github.com/wailsapp/wails/pull/4045)で[@fbbdev](https://github.com/fbbdev)が修正
- バインディングジェネレーターの出力における非公開モデルへの無効な参照を、[#4045](https://github.com/wailsapp/wails/pull/4045)で[@fbbdev](https://github.com/fbbdev)が修正
- 挿入コードをサービスファイルの末尾へ、[#4045](https://github.com/wailsapp/wails/pull/4045)で[@fbbdev](https://github.com/fbbdev)が移動
- バインディングジェネレーターでファイルを閉じる処理から返されるエラーの扱いを、[#4045](https://github.com/wailsapp/wails/pull/4045)で[@fbbdev](https://github.com/fbbdev)が修正
- ライフサイクルメソッドまたはhttpメソッドを定義しているものの、それ以外のバインドされたメソッドを持たないサービスに対する警告を、[#4045](https://github.com/wailsapp/wails/pull/4045)で[@fbbdev](https://github.com/fbbdev)が抑制
- システムの配色がライトの場合に、React以外のテンプレートでHello Worldフッターが表示されない問題を、[#4056](https://github.com/wailsapp/wails/pull/4056)で[@marcus-crane](https://github.com/marcus-crane)が修正
- macOSで非表示のメニュー項目に関する問題を[@leaanthony](https://github.com/leaanthony)が修正
- メッセージプロセッサーにおけるエラーの処理および書式設定を、[#4066](https://github.com/wailsapp/wails/pull/4066)で[@fbbdev](https://github.com/fbbdev)が修正
-  アプリケーション終了時にサービスのシャットダウンが実行されない問題を、[#4066](https://github.com/wailsapp/wails/pull/4066)で[@fbbdev](https://github.com/fbbdev)が修正
-  メニューの更新がメインスレッドで行われるよう[@leaanthony](https://github.com/leaanthony)が修正
- ドラッグおよびサイズ変更の仕組みをより堅牢にし、期待されるプラットフォームの動作へさらに近づける改善を、[#4100](https://github.com/wailsapp/wails/pull/4100)で[@fbbdev](https://github.com/fbbdev)が実施
- Webpack/angularがランタイム初期化コードを破棄する問題[#4097](https://github.com/wailsapp/wails/issues/4097)を、[#4100](https://github.com/wailsapp/wails/pull/4100)で[@fbbdev](https://github.com/fbbdev)が修正
- 最初から非表示に設定されたメニュー項目を、[#4116](https://github.com/wailsapp/wails/pull/4116)で[@IanVS](https://github.com/IanVS)が修正
- 拡張子のないリクエストで`[request]`が存在せず、`[request].html`が存在する場合に、assetFileServerが`.html`ファイルを配信しない問題を修正
- アイコン生成パスを、[#4125](https://github.com/wailsapp/wails/pull/4125)で[@robin-samuel](https://github.com/robin-samuel)が修正
- `fullscreen`、`unfullscreen`、`unminimise`、および`unmaximise`イベントが発行されない問題を、[#4130](https://github.com/wailsapp/wails/pull/4130)で[@oSethoum](https://github.com/osethoum)が修正
- 設定内のデフォルトバージョンに誤った接頭辞が付いていたために発生するNSIS Errorを、[#4126](https://github.com/wailsapp/wails/pull/4126)で[@robin-samuel](https://github.com/robin-samuel)が修正
- WindowsでDialogsランタイム関数がエスケープされたパスを返す問題を、[#4188](https://github.com/wailsapp/wails/pull/4188)で[TheGB0077](https://github.com/TheGB0077)が修正
- HKCU内のWebview2検出パスを[@leaanthony](https://github.com/leaanthony)が修正。
- macOSでの入力に関する問題を[@leaanthony](https://github.com/leaanthony)が修正。
- Windowsのアイコン生成タスクのファイル名を、[#4219](https://github.com/wailsapp/wails/pull/4219)で[@yulesxoxo](https://github.com/yulesxoxo)が修正。
- @kronの成果を基に、フレームレスウィンドウの透過性に関する問題を[@leaanthony](https://github.com/leaanthony)が修正。
- @kronの成果を基に、ウィンドウが無効または最小化されている場合のフォーカス呼び出しを[@leaanthony](https://github.com/leaanthony)が修正。
- @kronの成果を基に、タスクバーの再起動後にシステムトレイが表示されない問題を[@leaanthony](https://github.com/leaanthony)が修正。
- [#4245](https://github.com/wailsapp/wails/pull/4245)で、fallbackResponseWriterがFlush()を実装していない問題を修正
- fallbackResponseWriterがFlush()を実装していない問題を、[@superDingda]が[#4236](https://github.com/wailsapp/wails/issues/4236)で修正
- 非同期のGoバインド関数呼び出しが保留中にmacOSのウィンドウを閉じるとクラッシュする問題を、[@joshhardy](https://github.com/joshhardy)が[#4354](https://github.com/wailsapp/wails/pull/4354)で修正
- Windowsの効率モード起動時の競合状態を[@leaanthony](https://github.com/leaanthony)が修正
- Windowsのアイコンハンドルのクリーンアップを[@leaanthony](https://github.com/leaanthony)が修正。
- Windowsでの`OpenFileManager`を、[@PPTGamer](https://github.com/PPTGamer)が[#4375](https://github.com/wailsapp/wails/pull/4375)で修正。
- Linuxの最小／最大幅オプションを、@atterpacが[#3979](https://github.com/wailsapp/wails/pull/3979)で修正
- npmのバージョン更新により、TypeScriptテンプレートの型定義を追加（@atterpac、[#3966](https://github.com/wailsapp/wails/pull/3966)）
- SvelteKitテンプレートのCSS参照を、@atterpacが[#3945](https://github.com/wailsapp/wails/pull/3945)で修正
- window run()内の主要なコールバックがメインスレッドで呼び出されるように[@leaanthony](https://github.com/leaanthony)が修正
- ディレクトリ選択ダイアログの例を[@leaanthony](https://github.com/leaanthony)が修正
- index.htmlがない場合に表示する新しい中国語のエラーページを[@leaanthony](https://github.com/leaanthony)が作成
-  `windowDidBecomeKey`コールバックがメインスレッドで実行されるように[@leaanthony](https://github.com/leaanthony)が修正
-  フレームレスウィンドウのフルスクリーンを[@leaanthony](https://github.com/leaanthony)がサポート
-  ウィンドウ破棄ロジックを[@leaanthony](https://github.com/leaanthony)が改善
-  システムトレイに関連付けた場合のウィンドウ位置ロジックを[@leaanthony](https://github.com/leaanthony)が修正
-  フレームレスウィンドウのフルスクリーンを[@leaanthony](https://github.com/leaanthony)がサポート
- イベント処理を[@leaanthony](https://github.com/leaanthony)が修正
- ウィンドウのシャットダウンロジックを[@leaanthony](https://github.com/leaanthony)が修正
- 共通Taskfileで、TypeScriptテンプレートの場合にデフォルトでTypeScriptバインディングを生成するよう[@leaanthony](https://github.com/leaanthony)が変更
- ウィンドウが開いていない場合／システムトレイのみの場合に、WM_CLOSEメッセージでアプリケーションを終了する処理を、[@mmalcek](https://github.com/mmalcek)が[#3990](https://github.com/wailsapp/wails/pull/3990)で修正
- garbleビルドを、@5aaee9が[#3192](https://github.com/wailsapp/wails/pull/3192)で修正
- WindowsのNSISビルドを[@leaanthony](https://github.com/leaanthony)が修正
- 閉じられていないものが原因で、Linuxの複数選択ダイアログに発生するデッドロックを修正
- Windowsビルド中の.sysoファイルのクロスプラットフォームなクリーンアップを修正：
- amd64 AppImageのコンパイルを@atterpacが修正：
- ビルドアセットの更新を@ansxumanが修正：
- Linuxのシステムトレイにおける`OnClick`と`OnRightClick`の実装を@atterpacが修正
- Macで`AlwaysOnTop`が機能しない問題を修正：
-  `application.NewEditMenu`に重複項目が含まれる問題を修正
- 🐧 aarch64のコンパイルを修正
- ⊞ ラジオグループのメニュー項目を修正：
- MacOSで「name」と「outputfilename」が異なる場合に、実行可能な.appをビルドすると発生するエラーを修正
- ドラッグ＆ドロップの例でcustomEventProcessorを使用した際のバグを修正：
- 🐧 IgnoreMouseEventsの追加によって発生したLinuxのコンパイルエラーを修正：
- ⊞ sysoアイコンファイルの生成バグを修正：
- 🐧 Waylandでネイティブ実行するための修正を次から取り込み：
- 次では内部サービスメソッドをバインドしない：
- ⊞ 次でシステムトレイ起動時のパニックを修正：
- 次では内部サービスメソッドをバインドしない：
- ⊞ 次でシステムトレイ起動時のパニックを修正：
- メニュー項目とイベント処理を大幅にリファクタリング。現時点では主にmacOSを改善。担当：
- プラグインとイベントのリファクタリング後のテストを次で修正：
- ⊞ `Failed to unregister class Chrome_WidgetWin_0`の警告を修正。担当：
- モジュールの問題
- リサイズイベントのメッセージ処理を[atterpac](https://github.com/atterpac)が次で修正：
- 🐧 NixOSでのテーマ処理エラーを修正：
- Windowsで別ボリュームにまたがるプロジェクトのインストールを修正：
- フッターが表示されるようReactテンプレートのCSSを修正：
- refreshを最新版に更新し、devモードでの作業時に発生するゾンビプロセスを修正
- AppImageのWebKitファイル取得処理を[Atterpac](https://github.com/atterpac)が修正
- Doctorのaptパッケージ検証を[Atterpac](https://github.com/Atterpac)が次で修正：
- 終了時にアプリケーションがフリーズする問題（Darwin）を@5aaee9が次で修正：
- Windowsでのサンプルの背景色を修正：
- デフォルトのコンテキストメニューを[mmghv](https://github.com/mmghv)が次で修正：
- Darwinでの矢印キーの16進値を修正：
- Windowsのドラッグ＆ドロップを動作する状態に修正。追加者：
- ユーザーに適切なドライバーがない場合にDoctorで発生するLinuxのバグを修正
- 起動時のDPIスケーリング（Windows）を修正。[@almas-x](https://github.com/almas-x)が次で変更：
- 相対パスを使用するよう`go.mod`の置換行を修正。次を含むWindowsパスの問題を修正：
- ウィンドウが関連付けられていない場合のMacOSシステムトレイのクリック処理を修正：
- 不明なオプションによってWindowsビルドが失敗する問題を修正：
- 関連付けられたものがない状態でシステムトレイアイコンを左クリックするとWindowsでクラッシュする問題を修正
- ウィンドウを2回開いたときの誤ったbaseURLを、@5aaee9がPRで修正
- `WebviewWindow.Restore`メソッド内のif分岐の順序を修正：
- 複数回の`GetStartURL`呼び出しにまたがる`startURL`を正しく計算（次の場合）
- `Screen`構造体のJS型が対応するGo型と一致するよう修正：
- 登録済みイベントを適切にクリーンアップできるよう`WML.Reload`メソッドを修正
- Linuxでカスタムコンテキストメニューがすぐに閉じる問題を修正：
- バインディングによって生成されるモデルファイルの出力パスと拡張子を修正
- バインディングによって生成されるJSコード内のモデルファイルのインポートパスを修正
- 一部のLinuxディストリビューションでのドラッグ＆ドロップを修正：
- `wails3 task dev`使用時にmacOS用タスクが欠落する問題を修正：
- イベントの登録時にnilマップへの代入が発生する問題を修正：
- バインドされたメソッドのパラメーターのアンマーシャリングを修正：
- バインドされたメソッドから複数の戻り値を処理する際の問題を修正。修正者：
- システムのパッケージマネージャーでインストールされていない npm を doctor が検出できない問題を修正
- MicrosoftEdgeWebview2Setup.exe が欠落している問題を修正。協力者：
- ウィンドウ ID の処理が原因で Linux 上でランダムにクラッシュする問題を @leaanthony が修正。基になったもの：
- Linux 上で systemTray.setIcon がクラッシュする問題を修正。修正者：
- 上の `setFrameless` 関数で、最初の呼び出し時にウィンドウフレームが適用されるよう修正

### 変更

- **破壊的変更**：Go の map のセマンティクスを正確に反映するため、生成される JS/TS バインディングの map キーがオプショナルとしてマークされるようになりました。TypeScript で map の値にアクセスすると、`T` ではなく `T | undefined` が返されるようになったため、null チェックまたはアサーションが必要です（#4943）。変更者：`@fbbdev`
- `@wailsio/runtime` の変更に合わせて、`Event` の使用を `Events` に変更し、`Features/Events/Event System` のドキュメント内の関数呼び出しも適切に変更。変更者：@AbdelhadiSeddar
- `EnabledFeatures`、`DisabledFeatures`、`AdditionalBrowserArgs` をウィンドウごとのオプションからアプリケーションレベルの `Options.Windows` に移動（#4559）。変更者：@leaanthony
- `Drag N Drop` のサンプルの README を更新し、このサンプルで `Internal Drag and Drop` を実演していることを明記。変更者：@ndianabasi
- 各種デバッグログのレベルを Info から Debug に変更（@mbaklor）
- <strong>破壊的変更：</strong>ウィンドウオプションの `EnableDragAndDrop` を `EnableFileDrop` に名称変更
- <strong>破壊的変更：</strong>イベントコンテキストの `DropZoneDetails` を `DropTargetDetails` に名称変更
- **破壊的変更：**`WindowEventContext` の `DropZoneDetails()` メソッドを `DropTargetDetails()` に名称変更
- **破壊的変更：**`WindowDropZoneFilesDropped` イベントを削除。代わりに `WindowFilesDropped` を使用してください
- <strong>破壊的変更：</strong>HTML 属性を `data-wails-dropzone` から `data-file-drop-target` に変更
- <strong>破壊的変更：</strong>CSS のホバークラスを `wails-dropzone-hover` から `file-drop-target-active` に変更
- <strong>破壊的変更：</strong>Windows から `DragEffect`、`OnEnterEffect`、`OnOverEffect` オプションを削除（削除された IDropTarget の一部でした）
- ランタイムのすべての JSON 処理（メソッドバインディング、イベント、WebView リクエスト、通知、kvstore）を goccy/go-json に切り替え、パフォーマンスを 21-63% 向上し、メモリ割り当てを 40-60% 削減
- BoundMethod 構造体のレイアウトを最適化し、isVariadic フラグをキャッシュすることで、呼び出しごとのオーバーヘッドを削減
- 引数が `<=8` 個のメソッドではスタックに割り当てた引数バッファを使用し、ヒープ割り当てを回避
- メソッド呼び出しの結果収集を最適化し、戻り値が 1 つの場合のスライス割り当てを回避
- MIME タイプのキャッシュに sync.Map を使用し、並行処理性能を向上
- HTTP トランスポートのリクエストボディ読み取りにバッファプールを使用
- コンテンツタイプスニファーで CloseNotify チャネルを遅延割り当てし、リクエストごとの割り当てを削減
- アセットサーバーから CSS のデバッグログを削除
- MIME タイプの拡張子マップを拡張し、50 種類以上の一般的な Web 形式（フォント、音声、動画など）に対応
- Window の `X/Y` オプションに関するドキュメントを更新。変更者：@ruhuang2001
- フロントエンドバインディング生成用のオプションを追加し、`Frontend Runtime` のドキュメントを更新。変更者：@ndianabasi
- Wails の v3 Asset Server に関するドキュメントページを更新。変更者：@ndianabasi
- **破壊的変更**：パッケージレベルのダイアログ関数（`application.InfoDialog()`、`application.QuestionDialog()` など）を削除。代わりに `app.Dialog` マネージャー（`app.Dialog.Info()`、`app.Dialog.Question()`、`app.Dialog.Warning()`、`app.Dialog.Error()`、`app.Dialog.OpenFile()`、`app.Dialog.SaveFile()`）を使用してください
- 実際の API に合わせてダイアログのドキュメントを更新：`app.Dialog.*`、コールバックを伴う `AddButton()`（`SetButtons()` ではない）、`SetDefaultButton(*Button)`（文字列ではない）、`AddFilter()`（`SetFilters()` ではない）、`SetFilename()`（`SetDefaultFilename()` ではない）、およびフォルダー選択用の `app.Dialog.OpenFile().CanChooseDirectories(true)` を使用
- **破壊的変更**：プロダクションビルドがデフォルトになりました。開発ビルドを作成するには、Taskfile で `DEV=true` を設定してください。サンプル用に新しいプロジェクトを生成。変更者：@leaanthony
- データ引数が 0 個または 1 個のカスタムイベントを発行する場合、そのデータ値はスライスでラップされず、Data フィールドに直接代入されるようになりました。[#4633](https://github.com/wailsapp/wails/pull/4633) で [@fbbdev](https://github.com/fbbdev) が変更
- Windows のトレイで `NIS_HIDDEN` を切り替えることにより、`SystemTray.Show()`/`Hide()` が反映されるようになり、アプリを完全に非表示にして再表示できるようになりました（#4653）。
- トレイの登録で解決済みアイコンを再利用し、`NOTIFYICON_VERSION_4` を一度だけ設定するとともに、`NIF_SHOWTIP` を有効にすることで、Explorer の再起動後にツールチップが復元されるようになりました（#4653）。
- macOS：メニューバーと Dock の領域を除外してウィンドウを中央に配置するため、`frame` の代わりに `visibleFrame` を使用
- macOS：メニューバーと Dock の領域を除外してウィンドウを中央に配置するため、`frame` の代わりに `visibleFrame` を使用
- `-config` パラメーターを指定して `wails3 update build-assets` を実行すると、`-product*` パラメーターで設定した値は
- `window.NativeWindowHandle()` -> `window.NativeWindow()`。[#4471](https://github.com/wailsapp/wails/pull/4471) で @leaanthony が変更
- 内部のウィンドウ処理をリファクタリング。[#4471](https://github.com/wailsapp/wails/pull/4471) で @leaanthony が変更
- `application.WindowIDKey` と `application.WindowNameKey` を削除（`application.WindowKey` に置換）。変更者：[@leaanthony](https://github.com/leaanthony)
- ContextMenuData が any ではなく文字列を返すように変更。変更者：[@leaanthony](https://github.com/leaanthony)
- JS/TS バインディングで、固定長配列型のクラスフィールドが空ではなく、想定される長さで初期化されるようになりました。[#4001](https://github.com/wailsapp/wails/pull/4001) で [@fbbdev](https://github.com/fbbdev) が変更
- ContextMenuData が any ではなく文字列を返すように変更。変更者：[@leaanthony](https://github.com/leaanthony)
- `application.NewService` は、オプショナルパラメーターとしてオプションを受け取らなくなりました（代わりに `application.NewServiceWithOptions` を使用してください）。[#4024](https://github.com/wailsapp/wails/pull/4024) で [@leaanthony](https://github.com/leaanthony) が変更
- `nanoid` への依存を削除。変更者：[@leaanthony](https://github.com/leaanthony)
- mica/acrylic/tabbed ウィンドウスタイルに関する Window のサンプルを更新。変更者：[@leaanthony](https://github.com/leaanthony)
- JS/TS バインディングから `internal.js/ts` モデルファイルを削除し、すべてのモデルが `models.js/ts` に配置されるようになりました。[#4045](https://github.com/wailsapp/wails/pull/4045) で [@fbbdev](https://github.com/fbbdev) が変更
- JS/TS バインディングで、名前付き型が別の名前付き型のエイリアスとしてレンダリングされることはなくなりました。従来の動作はエイリアスのみに制限されます。[#4045](https://github.com/wailsapp/wails/pull/4045) で [@fbbdev](https://github.com/fbbdev) が変更
- JS/TS バインディングのクラスモードで、型が型パラメーターである構造体フィールドはオプショナルとしてマークされ、自動的には初期化されなくなりました。[#4045](https://github.com/wailsapp/wails/pull/4045) で [@fbbdev](https://github.com/fbbdev) が変更
- テンプレートから ESLint を削除。[#4059](https://github.com/wailsapp/wails/pull/4059) で [@IanVS](https://github.com/IanVS) が変更
- 著作権表記の年を 2025 に更新。[#4037](https://github.com/wailsapp/wails/pull/4037) で [@IanVS](https://github.com/IanVS) が変更
- [@IanVS](https://github.com/IanVS) が [#4075](https://github.com/wailsapp/wails/pull/4075) で event.Sender のドキュメントを追加
- [@leaanthony](https://github.com/leaanthony) が Go 1.24 のサポートを追加
- [@fbbdev](https://github.com/fbbdev) が [#4066](https://github.com/wailsapp/wails/pull/4066) で、`ServiceStartup` フックを `application.New` 内ではなく、`App.Run` の呼び出し時に実行するよう変更
- [@fbbdev](https://github.com/fbbdev) が [#4066](https://github.com/wailsapp/wails/pull/4066) で、`ServiceStartup` のエラー発生時にプロセスを終了せず、`App.Run` からエラーを返すよう変更
- [@fbbdev](https://github.com/fbbdev) が [#4066](https://github.com/wailsapp/wails/pull/4066) で、JS からのバインディングおよびダイアログの呼び出しが文字列ではなくエラーオブジェクトで reject されるよう変更
- [@leaanthony](https://github.com/leaanthony) が Windows でのシステムトレイメニューの配置を改善
- [@fbbdev](https://github.com/fbbdev) が [#4100](https://github.com/wailsapp/wails/pull/4100) で JS ランタイムを TypeScript に移植
- [@fbbdev](https://github.com/fbbdev) が [#4100](https://github.com/wailsapp/wails/pull/4100) で、ランタイムをインポート直後に初期化するよう変更。ウィンドウの読み込みを待つ必要はありません
- [@fbbdev](https://github.com/fbbdev) が [#4100](https://github.com/wailsapp/wails/pull/4100) で、ランタイムから init メソッドを削除。副作用のためのインポートを使用して初期化できます
- [@fbbdev](https://github.com/fbbdev) が [#4100](https://github.com/wailsapp/wails/pull/4100) で、バインドされたメソッドが、キャンセル時に `CancelError` で reject される `CancellablePromise` を返すよう変更。実際の呼び出し結果は破棄されます
- [@fbbdev](https://github.com/fbbdev) が [#4067](https://github.com/wailsapp/wails/pull/4067) で、組み込みサービス型の名前を `Service` に統一
- [@fbbdev](https://github.com/fbbdev) が [#4067](https://github.com/wailsapp/wails/pull/4067) で、オプションを受け取る組み込みサービス作成関数の名前を `NewWithConfig` に統一
- [@fbbdev](https://github.com/fbbdev) が [#4067](https://github.com/wailsapp/wails/pull/4067) で、Go API との一貫性を保つため、`sqlite` サービスの `Select` メソッドを `Query` に改名
- テンプレート：[@IanVS](https://github.com/IanVS) が [#4133](https://github.com/wailsapp/wails/pull/4133) でランタイムを "dependencies" に移動し、package.json ファイルを整理
- [@popaprozac](https://github.com/popaprozac) が [#4171](https://github.com/wailsapp/wails/pull/4171) で、特定の macOS API を有効にするため、開発時にアプリバンドルを作成してアドホック署名するよう変更
- [@leaanthony](https://github.com/leaanthony) がビルドアセットをプラットフォーム固有のディレクトリへ移動
- [@leaanthony](https://github.com/leaanthony) が Taskfile をプラットフォーム固有のディレクトリへ移動して改名
- [@leaanthony](https://github.com/leaanthony) が、`index.html` がない場合のユーザー体験を大幅に改善
- [Windows] [@leaanthony](https://github.com/leaanthony) が最小化と復元のパフォーマンスを改善。[562589540](https://github.com/562589540) による元の [PR](https://github.com/wailsapp/wails/pull/3955) に基づく変更
- [@leaanthony](https://github.com/leaanthony) が `ShouldClose` オプションを削除（代わりに events.Common.WindowClosing のフックを登録してください）
- [Windows] [@leaanthony](https://github.com/leaanthony) がウィンドウを開く際のちらつきを軽減
- [@leaanthony](https://github.com/leaanthony) が、内部関数として意図されていた `Window.Destroy` を削除
- [@leaanthony](https://github.com/leaanthony) が `WindowClose` イベントを `WindowClosing` に改名
- [@leaanthony](https://github.com/leaanthony) が、フロントエンドのビルドでビルド種別に応じて vite 環境 "development" または "production" を使用するよう変更
- [@leaanthony](https://github.com/leaanthony) が go-webview2 を v1.19 に更新
- taskfile のフォークが使用されることを保証（@leaanthony）
- 次を使用してインストールする際のバージョン問題を修正するため、Taskfile のフォークを更新
- 次を使用してインストールする際のバージョン問題を修正するため、Taskfile のフォークを使用
- `service.OnStartup` がエラー発生時にアプリケーションを終了し、次を実行するよう変更
- ユーザー操作との整合性を高めるため、システムトレイのクリックメッセージングをリファクタリング（担当：
- 次を生成するフレームワークをサポートするため、埋め込みアセットに `all:frontend/dist` を含めるよう変更
- [leaanthony](https://github.com/leaanthony) が次で Taskfile をリファクタリング
- `go-webview2` v1.0.16 へアップグレード（担当：
- `Screen` 型が `Id` ではなく `ID` を含むよう修正（担当：
- `application.ServiceOptions` をサポートするため、`go.mod.tmpl` の Wails バージョンを更新（担当：
- [windom](https://github.com/windom/) が次でサービス名の判定を修正
- [leaanthony](https://github.com/leaanthony) が mkdocs serve で docker を使用するよう変更
- 開発用設定を `config.yml` に統合（担当：
- システムトレイのダイアログで、利用可能な場合はアプリケーションアイコンをデフォルトで使用するよう変更（Windows、担当：
- macOS での GPU とメモリのレポートを改善（担当：
- `WebviewGpuIsDisabled` と `EnableFraudulentWebsiteWarnings` を削除
- Events API の変更：`On`/`Emit` -> ユーザーイベント、`OnApplicationEvent` ->
- [TheGB0077](https://github.com/TheGB0077) が次で Linux の Events API を修正
- [CI] Actions を改善し、フォークでも Actions を実行できるようにし、
- `AbsolutePosition()` を `Position()` に改名（担当：
- Linux の WebKit 依存関係を webkitgtk2-4.0 から webkit2gtk-4.1 に更新して
- 同梱の JS ランタイムスクリプトを ESM モジュールに変更：これをインポートする script タグは
- `@wailsio/runtime` パッケージは、その API を `window.wails` 上に公開しません
- Window API モジュール `@wailsio/runtime/src/window` が、包含するものを公開するようになりました
- JS Window API を現在の Go `WebviewWindow` と一致するよう更新
- バインディングジェネレーターは、デフォルトで ID による呼び出しを使用するようになりました。`-id` CLI オプションは
- 新しいバインディングコードのレイアウト：以前、出力ファイルはフォルダーに整理されていました
- 構造体フィールド `application.Options.Bind` の名前を次のように変更
- サービスをバインドするための新しい構文：サービスインスタンスを次でラップする必要があります
- 非ターミナル環境または CI 環境ではスピナーを無効化（担当：

### 削除

- **破壊的変更**：ウィンドウごとの `WindowsWindow` オプションから `EnabledFeatures`、`DisabledFeatures`、`AdditionalLaunchArgs` を削除。代わりに、アプリケーションレベルの `Options.Windows.EnabledFeatures`、`Options.Windows.DisabledFeatures`、`Options.Windows.AdditionalBrowserArgs` を使用してください。これらのフラグは、共有 WebView2 環境にグローバルに適用されます（#4559）（@leaanthony）
- Windows のネイティブ `IDropTarget` 実装を削除し、JavaScript ベースの方式に変更（v2 の動作と一致）
- github.com/wailsapp/mimetype 依存関係を削除し、拡張した拡張子マップと標準ライブラリの http.DetectContentType に変更。バイナリサイズを約 1.2MB 削減
- Linux ファイルエクスプローラー用の最小限の .desktop ファイルパーサーを実装し、gopkg.in/ini.v1 への依存を削除。約 45KB を削減
- Go 1.21+ の標準ライブラリの slices パッケージと最小限の内部ヘルパーを使用して、ランタイムコードから samber/lo を削除。約 310KB を削減
- Darwin の URL スキームハンドラーからデバッグ用 printf 文を削除（#4834）
- **破壊的変更**：`linux:WindowLoadChanged` イベントを削除。WebView の読み込み完了を検出するには、代わりに `linux:WindowLoadFinished` を使用（#3896、@leaanthony）

### 破壊的変更

- **Manager API のリファクタリング**：コード構成と見つけやすさを向上させるため、アプリケーション API をフラットな構造から整理されたマネージャー構造に再編。[#4359](https://github.com/wailsapp/wails/pull/4359) で [@leaanthony](https://github.com/leaanthony) が実施
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
- Service のメソッド名を変更：`Name` -> `ServiceName`、`OnStartup` -> `ServiceStartup`、`OnShutdown` -> `ServiceShutdown`。[@leaanthony](https://github.com/leaanthony) が実施
- `Path` メソッドと `Paths` メソッドを `application` パッケージへ移動。[@leaanthony](https://github.com/leaanthony) が実施
- アプリケーションメニューを macOS 専用に変更。[@leaanthony](https://github.com/leaanthony) が実施

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

無視されなくなり、設定値を上書きします。

## v3.0.0-alpha.24 - 2025-08-13

## 追加

## v3.0.0-alpha.23 - 2025-08-11

## 修正

## v3.0.0-alpha.22 - 2025-08-10

## 追加

## 変更

+ 過剰に広範な Linux パッケージ依存関係と、古くなった RPM 依存関係を修正。

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

`go install`（@leaanthony による）

## v3.0.0-alpha.8.1 - 2024-12-07

### 変更

`go install`（@leaanthony による）

## v3.0.0-alpha.8 - 2024-12-06

### 追加

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### 変更

`service.OnShutdown`以前に起動したすべてのサービスについて（@atterpac による、[#3920](https://github.com/wailsapp/wails/pull/3920)）   @atterpac（[#3907](https://github.com/wailsapp/wails/pull/3907)）   サブフォルダー（@atterpac による、[#3887](https://github.com/wailsapp/wails/pull/3887)）   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913)（[#3778](https://github.com/wailsapp/wails/pull/3778)）   [northes](https://github.com/northes)（[#3836](https://github.com/wailsapp/wails/pull/3836)）   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   （`EnabledFeatures`および`DisabledFeatures`オプションにより置き換え済み）（[leaanthony](https://github.com/leaanthony) による）

### 修正

channel 変数（@michael-freling による、[#3925](https://github.com/wailsapp/wails/pull/3925)）   [ansxuman](https://github.com/ansxuman)（[#3924](https://github.com/wailsapp/wails/pull/3924)）   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony)（[#3841](https://github.com/wailsapp/wails/pull/3841)）   Darwin の編集メニューにおける`PasteAndMatchStyle`ロール（[johnmccabe](https://github.com/johnmccabe) による、[#3839](https://github.com/wailsapp/wails/pull/3839)）   [#3840](https://github.com/wailsapp/wails/issues/3840)（[#3854](https://github.com/wailsapp/wails/pull/3854)、[kodflow](https://github.com/kodflow) による）   [@leaanthony](https://github.com/leaanthony)   が異なる問題（@nickisworking による、[#3789](https://github.com/wailsapp/wails/pull/3789)）

## v3.0.0-alpha.7 - 2024-09-18

### 追加

[mmghv](https://github.com/mmghv)（[#3665](https://github.com/wailsapp/wails/pull/3665)）   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac)および[leaanthony](https://github.com/leaanthony)（[#3570](https://github.com/wailsapp/wails/pull/3570)）

### 変更

Application Events `OnWindowEvent` -> Window Events（[leaanthony](https://github.com/leaanthony) による）   [#3734](https://github.com/wailsapp/wails/pull/3734)   `v3/`または`v3-`の接頭辞が付いたブランチ（[stendler](https://github.com/stendler) による、[#3747](https://github.com/wailsapp/wails/pull/3747)）

### 修正

[etesam913](https://github.com/etesam913)（[#3742](https://github.com/wailsapp/wails/pull/3742)）   [atterpac](https://github.com/atterpac)（[#3721](https://github.com/wailsapp/wails/pull/3721)）   [atterpac](https://github.com/atterpac)（[#3675](https://github.com/wailsapp/wails/pull/3675)）   [#1811](https://github.com/wailsapp/wails/pull/1811)（[#3614](https://github.com/wailsapp/wails/pull/3614)、[@stendler](https://github.com/stendler) による）   [#3720](https://github.com/wailsapp/wails/pull/3720)（[leaanthony](https://github.com/leaanthony) による）   [#3693](https://github.com/wailsapp/wails/issues/3693)（[@DeltaLaboratory](https://github.com/DeltaLaboratory) による）   [#3720](https://github.com/wailsapp/wails/pull/3720)（[leaanthony](https://github.com/leaanthony) による）   [#3693](https://github.com/wailsapp/wails/issues/3693)（[@DeltaLaboratory](https://github.com/DeltaLaboratory) による）   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746)（[@stendler](https://github.com/stendler) による）   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### 修正

## v3.0.0-alpha.5 - 2024-07-30

### 追加

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   アイコンのクリック（@5aaee9 による、[#2991](https://github.com/wailsapp/wails/pull/2991)）   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev)（[#3282](https://github.com/wailsapp/wails/pull/3282)）   @[Atterpac](https://github.com/Atterpac)（[#3022](https://github.com/wailsapp/wails/pull/3022])）   [@marcus-crane](https://github.com/marcus-crane)（[#3146](https://github.com/wailsapp/wails/pull/3146)）   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev)（[#3281](https://github.com/wailsapp/wails/pull/3281)）   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   @Mai-Lapyst による[PR](https://github.com/wailsapp/wails/pull/2044)を基に作成   [@fbbdev](https://github.com/fbbdev)（[#3295](https://github.com/wailsapp/wails/pull/3295)）   [@fbbdev](https://github.com/fbbdev)（[#3295](https://github.com/wailsapp/wails/pull/3295)）   [@fbbdev](https://github.com/fbbdev)（[#3295](https://github.com/wailsapp/wails/pull/3295)）   npm パッケージ（[@fbbdev](https://github.com/fbbdev) による、[#3334](https://github.com/wailsapp/wails/pull/3334)）   [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT`（[@abichinger](https://github.com/abichinger) による、[#3429](https://github.com/wailsapp/wails/pull/3429)）   [@abichinger](https://github.com/abichinger)（[#3431](https://github.com/wailsapp/wails/pull/3431)）   [@bruxaodev](https://github.com/bruxaodev)（[#3667](https://github.com/wailsapp/wails/pull/3667)）   [@OlegGulevskyy](https://github.com/OlegGulevskyy)（[#3674](https://github.com/wailsapp/wails/pull/3674)）

### 修正

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane)（[#3515](https://github.com/wailsapp/wails/pull/3515)）   [atterpac](https://github.com/atterac)（[#3512](https://github.com/wailsapp/wails/pull/3512)）   [atterpac](https://github.com/atterpac)（[#3477](https://github.com/wailsapp/wails/pull/3477)）   [Atterpac](https://github.com/atterpac)による変更（[#3320](https://github.com/wailsapp/wails/pull/3320)）。   [#3306](https://github.com/wailsapp/wails/pull/3306)で対応。   [#2972](https://github.com/wailsapp/wails/pull/2972)。   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv)（[#2750](https://github.com/wailsapp/wails/pull/2750)）。   [#2753](https://github.com/wailsapp/wails/pull/2753)。   [jaybeecave](https://github.com/jaybeecave)（[#3052](https://github.com/wailsapp/wails/pull/3052)）。   [@pylotlight](https://github.com/pylotlight)（[PR](https://github.com/wailsapp/wails/pull/3039)）   がインストールされるようになりました。[@pylotlight](https://github.com/pylotlight)が[PR](https://github.com/wailsapp/wails/pull/3032)で追加。   [PR](https://github.com/wailsapp/wails/pull/3145)   スペース — @leaanthony。   [thomas-senechal](https://github.com/thomas-senechal)によるPR [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal)によるPR [#3208](https://github.com/wailsapp/wails/pull/3208)   アタッチされたウィンドウ。[tw1nk](https://github.com/tw1nk)によるPR [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev)（[#3279](https://github.com/wailsapp/wails/pull/3279)）   `FRONTEND_DEVSERVER_URL`が存在します。   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev)（[#3295](https://github.com/wailsapp/wails/pull/3295)）   [@fbbdev](https://github.com/fbbdev)によるリスナー（[#3295](https://github.com/wailsapp/wails/pull/3295)）   [@abichinger](https://github.com/abichinger)（[#3330](https://github.com/wailsapp/wails/pull/3330)）   [@fbbdev](https://github.com/fbbdev)によるジェネレーター（[#3334](https://github.com/wailsapp/wails/pull/3334)）   [@fbbdev](https://github.com/fbbdev)によるジェネレーター（[#3334](https://github.com/wailsapp/wails/pull/3334)）   [@abichinger](https://github.com/abichinger)（[#3346](https://github.com/wailsapp/wails/pull/3346)）   [@hfoxy](https://github.com/hfoxy)（[#3417](https://github.com/wailsapp/wails/pull/3417)）   [@hfoxy](https://github.com/hfoxy)（[#3426](https://github.com/wailsapp/wails/pull/3426)）   [@fbbdev](https://github.com/fbbdev)（[#3431](https://github.com/wailsapp/wails/pull/3431)）   [@fbbdev](https://github.com/fbbdev)（[#3431](https://github.com/wailsapp/wails/pull/3431)）   [@pekim](https://github.com/pekim)による変更（[#3458](https://github.com/wailsapp/wails/pull/3458)）   [@robin-samuel](https://github.com/robin-samuel)。   [@5aaee9](https://github.com/5aaee9)によるPR [#3466](https://github.com/wailsapp/wails/pull/3622)。   [@windom](https://github.com/windom/)（[#3636](https://github.com/wailsapp/wails/pull/3636)）。   [@bruxaodev](https://github.com/bruxaodev/)によるWindows対応（[#3691](https://github.com/wailsapp/wails/pull/3691)）。

### 変更

[mmghv](https://github.com/mmghv)（[#3611](https://github.com/wailsapp/wails/pull/3611)）   [atterpac](https://github.com/atterpac)によるUbuntu 24.04 LTS対応（[#3461](https://github.com/wailsapp/wails/pull/3461)）   には`type="module"`属性が必要です。[@fbbdev](https://github.com/fbbdev)による変更（[#3295](https://github.com/wailsapp/wails/pull/3295)）   オブジェクトを作成しますが、WMLシステムは起動しません。これはカプセル化を改善するための変更です。必要に応じて、新しい`WML.Enable`メソッドを呼び出すことで、WMLシステムを手動で起動できます。バンドルされたJSランタイムスクリプトは、引き続き両方の操作を自動的に実行します。[@fbbdev](https://github.com/fbbdev)による変更（[#3295](https://github.com/wailsapp/wails/pull/3295)）   windowオブジェクトをデフォルトエクスポートとして提供します。ESMの名前付きインポートまたは名前空間インポート構文を使用して個々のメソッドをインポートすることは、今後できません。   API。一部のメソッドでは名前またはプロトタイプが変更されました。具体的には、`Screen`は`GetScreen`に、`GetZoomLevel`／`SetZoomLevel`は`GetZoom`／`SetZoom`になります。また、`GetZoom`、`Width`、`Height`は、値をオブジェクトにラップせず直接返すようになりました。[@fbbdev](https://github.com/fbbdev)による変更（[#3295](https://github.com/wailsapp/wails/pull/3295)）   は削除されました。名前による呼び出しに戻すには、`-names` CLIオプションを使用してください。[@fbbdev](https://github.com/fbbdev)による変更（[#3468](https://github.com/wailsapp/wails/pull/3468)）   以前は包含パッケージにちなんだ名前が付けられていましたが、現在はモジュールパスを含む完全なGoインポートパスが使用されます。[@fbbdev](https://github.com/fbbdev)による変更（[#3468](https://github.com/wailsapp/wails/pull/3468)）   `application.Options.Services`。[@fbbdev](https://github.com/fbbdev)による変更（[#3468](https://github.com/wailsapp/wails/pull/3468)）   `application.NewService`の呼び出し。[@fbbdev](https://github.com/fbbdev)による変更（[#3468](https://github.com/wailsapp/wails/pull/3468)）   [@DeltaLaboratory](https://github.com/DeltaLaboratory)（[#3574](https://github.com/wailsapp/wails/pull/3574)）
