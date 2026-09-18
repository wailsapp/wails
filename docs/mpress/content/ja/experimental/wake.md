---
title: "Wake"
description: "既存の Taskfile を実行し、より高速なインクリメンタルビルド、構造化された出力、デフォルトでの並列実行を提供する、Wails 対応の実験的なビルドランナーです。"
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

@note{type="caution" title="実験的機能"}
Wake は `WAILS_USE_WAKE=true` によるオプトイン方式であり、デフォルトのランナーでは<strong>ありません</strong>。 この変数が未設定の場合、`wails3 build / package / sign / task` は従来とまったく同じように動作します。 対応する機能や動作は、リリース間で変更される可能性があります。

@end

Wake は、`wails3` 向けの<strong>実験的な代替ビルドランナー</strong>です。プロジェクトにすでに存在するものと同じ `Taskfile.yml`、つまり同じ task、dep、var、template、include、およびプラットフォーム名前空間の構文を読み取り、汎用の [Task](https://taskfile.dev) ランタイムではなく、Wails 対応のエグゼキューターで実行します。

目的は Task を置き換えることではありません。Wails プロジェクトの実際のビルド方法に特化し、`wails3` CLI の他の部分と一致するセマンティクス、出力、デフォルト設定を備えたランナーを提供することです。**Wake だけを使用する限り、Taskfile を変更する必要はありません。**

## 存在する理由

Wake と Task ランタイムはどちらも `wails3` に組み込まれているため、別途バイナリをインストールする必要はありません。違いは、Wake が<strong>対象領域を理解している</strong>ことです。汎用ランナーは、Taskfile に記載されたステップを指定された順序で実行します。Wake は、Wails のビルドが実際にどのようなもの<em>か</em>を理解しています。つまり、フロントエンドバンドルがバイナリに埋め込まれ、そのバイナリがプラットフォーム固有の成果物にパッケージ化され、アイコンやバインディングも併せて生成されることを認識しています。この知識を活用して、汎用ランナーにはできない方法でビルドを最適化します。

- **実際のビルドに必要な処理だけを実行します。** Wake は、各ステップの実際の入力と出力を自ら追跡します。Go ビルドの場合は、モジュールグラフに加えて、依存するステップの出力が対象になります。そのため、関連する変更が何もなければ、コンパイラーとリンカーの再実行を完全に省略します。汎用ランナーがステップを省略できるのは、監視対象のファイルが Taskfile にあらかじめすべて明記されている場合に限られます。一方、Wake はビルドについてすでに把握している情報から監視対象を特定します。何も変更されていない状態での再ビルドは、およそ **~20 ms（Wake）、Task では ~316 ms** です。コールドビルドの実時間は同じです。所要時間の大半を `npm install`、Vite、Go コンパイラーが占めるためです。

- **同時に実行できる処理を把握しています。** Wake はどのステップが互いに独立しているかを理解しているため、デフォルトで並列実行し、判定行には得られた高速化の効果が表示されます。兄弟ステップの出力が入り交じって調査の妨げになる場合は、`WAKE_SERIAL=true` で並列実行を無効にしてください。

- **wails3 が制御する構造化出力。** Wake は wails3 独自のレポーターを使用して表示します。計画されたステップごとに 1 行を表示し、ステータスをリアルタイムで更新し、最後にフェーズ別の内訳を色分けして示します。また、失敗パネル内の `file:line` リンクはクリックできます。`NO_COLOR` および非 TTY 環境（CI ログ）でも適切に簡略表示されます。

- **組み込みであるため、Wails とともに発展できます。** Wake はサードパーティ製ツールではなく `wails3` の一部なので、別プロジェクトによる実装を待たずに、新しいビルド機能を直接追加できます。これにより、現在の Taskfile が `wails3` バイナリをシェル経由で呼び出している箇所（呼び出しごとにプロセスを生成）についても、クロスプラットフォームのスクリプトやツールをネイティブに実行できる可能性が生まれます。こうした処理をプロセス内に取り込むことでオーバーヘッドが減り、今後さらに高速化できます。

## Wake の有効化

Wake は、`WAILS_USE_WAKE=true` 環境変数によって完全に制御されます。 この変数が未設定の場合（または `true` 以外に設定されている場合）、すべての `wails3` コマンドは従来どおり組み込みの Task ランタイムを使用します。

```bash
# Default: Task runtime, no Wake involvement
wails3 build

# Opt in: Wake drives build / package / sign / task <name>
WAILS_USE_WAKE=true wails3 build
WAILS_USE_WAKE=true wails3 package
WAILS_USE_WAKE=true wails3 task <some-task-name>
```

このフラグは `wails3 build`、`wails3 package`、`wails3 sign`、`wails3 task <name>` に適用されます。`wails3 dev` は現時点では影響を<strong>受けません</strong>。開発用ウォッチャーは引き続き独自のパイプラインを使用します。

@note{type="tip" title="Wake は安全に有効化できます"}
Wake が未実装の Taskfile 機能を検出すると、実行全体を同一プロセス内の組み込み Task ランタイムに引き渡します。外部の `task` バイナリをインストールする必要はありません。最悪の場合でも、フラグを指定しなかった場合とまったく同じ動作になります。

@end

## 階層化されたローカルオーバーライド

Wake は、<strong>基本 Taskfile とローカルオーバーライド</strong>をサポートしています。`Taskfile.yml` と同じ場所にファイルを配置すると、その定義が優先されます。

| ファイル | 用途 | 優先順位 |
| --- | --- | --- |
| `Taskfile.yml` | 基本設定、コミット対象 | 最低 |
| `Taskfile.override.yml` / `.yaml` | チーム全体で使用するオーバーライド、コミット対象 | 中間 |
| `Taskfile.local.yml` / `.yaml` | 個人用、通常は git の追跡対象外 | 最高 |

**マージのセマンティクス（ローカル側を優先）：**

- <strong>同じ名前</strong>のタスクは、基本タスクをオーバーライドします。リストフィールド（`cmds`、`deps`、`sources`、`generates`、`platforms`、`status`、`preconditions`、`aliases`）がオーバーライド側に指定されている場合、それらは基本側の値を<strong>置き換えます</strong>。オーバーライド側で省略されたフィールドは、基本側の値が保持されます。
- `env` と `vars` は<strong>キーごとにマージ</strong>され、キーが競合した場合はオーバーライド側が優先されます。
- オーバーライドファイルに<strong>のみ</strong>存在するタスクは<strong>追加</strong>されます。

たとえば、コミット済みの `Taskfile.yml` では開発用フラグを指定してビルドするものの、自分のマシンでは常に本番用ビルドを行う必要がある場合は、次のようにします。

```yaml
# Taskfile.local.yml (git-ignored, yours)
tasks:
  build:
    cmds:
      - go build -tags production -o bin/app .
  smoke:
    cmds:
      - ./bin/app --selftest
```

これで、コミット済みの Taskfile を変更することなく、`build` は本番用コマンドを実行し、`smoke` も利用できるようになります。

@note{type="note" title="信頼モデル"}
オーバーライドファイルは自動的に検出され、確認なしで適用されます。これによって新たな権限が付与されることはありません。Taskfile はもともと任意のシェルコマンドを実行できるため、オーバーライドで可能になるのは `Taskfile.yml` の編集でも実行できることに限られます。コミット済みの `Taskfile.override.*` は PR の差分に表示されます。`Taskfile.local.*` は各自のマシン上で作成されます。不正な形式のオーバーライドがある場合、黙って省略されるのではなく、実行が中止されます。再現可能な CI ビルドのためにオーバーライドの検出を完全に無効化するには、`WAILS_NO_OVERRIDES=true` を設定します。

@end

## 自動フォールバック

Wake が未実装の Taskfile 機能に遭遇すると、実行全体を組み込みの Task ランタイムに引き渡します。現在、次の機能がこのフォールバックを発生させます。

- Taskfile レベルの `dotenv`
- `interleaved` 以外の `output` モード
- `requires` ブロック
- `interval`（Taskfile レベルまたはタスクレベル）
- `always` 以外の `run` モード
- タスク内の `short`
- タスク内の`defer`

## 環境変数

| 変数 | 効果 |
| --- | --- |
| `WAILS_USE_WAKE` | `true`を指定すると、ルーティング可能な`wails3`動詞で Wake が有効になります。それ以外では Task ランタイムが使用されます |
| `WAILS_NO_OVERRIDES` | `true`を指定すると、`Taskfile.local.*` / `.override.*`の検出をスキップします（再現可能なビルド向け） |
| `WAKE_VERBOSE` | サブプロセスの stdout/stderr を取り込み、失敗時にのみ表示する代わりに、リアルタイムでストリーミングします |
| `WAKE_SILENT` | タスクの出力を完全に抑制します |
| `WAKE_SERIAL` | `true`を指定すると、`deps:`の並列ファンアウトを無効にします（デフォルトは並列です） |
| `WAKE_FORCE` | `true`を指定すると、すべてのキャッシュをバイパスし、完全なクリーンリビルドを行います |
| `WAKE_DEBUG` | リゾルバーの内部情報（DAG、依存関係、変数参照、実行ルーティング）をログに記録します |
| `WAKE_NOTICE` | 実行ごとの「wake (experimental)」通知を非表示にするには`off`を指定します |

ビルドキャッシュは`.wake/cache.json`に保存されます（Task は`.task/`を使用します）。

## フィードバック

Wake は実験的な機能であり、今後の方向性は皆様からのフィードバックによって決まります。お試しになった場合は、高速になったか、分かりやすくなったか、何か問題が発生したかをぜひお聞かせください。実行した内容、期待した結果、実際に起きたことが記載された報告が最も役立ちます。[Wake のフィードバックに関するディスカッション](https://github.com/wailsapp/wails/discussions/5679)でお知らせください。
