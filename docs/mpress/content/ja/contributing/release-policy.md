---
title: "リリースポリシー"
description: "Wails v3 のナイトリーおよびウィークリーリリースのバージョン付け、メンテナーによるリリース方法、ユーザーがバージョンを固定すべき理由について説明します。"
slug: "contributing/release-policy"
sourcePath: "contributing/release-policy.md"
---

## モデルの概要表

Wails v3 は**機能ではなく、スケジュールに基づいて**リリースされます。`v3.0.0`のベースライン以降は次のとおりです。

| リリース | タグ | 頻度 | GitHub リリース | npm `@wailsio/runtime` |
| --- | --- | --- | --- | --- |
| ナイトリー | `v3.<next>.0-nightly.<YYYYMMDD>`（UTC の日付。同じ日に再実行した場合は `.N` を付加） | 自動、毎日。`UNRELEASED_CHANGELOG.md`に内容がある場合のみ | プレリリース。Latest には**しない** | `v`を除いた同じバージョン、dist-tag は `nightly` |
| ウィークリー | `v3.<next>.0` | 自動、日曜日の 20:00（Australia/Sydney）**（曜日は未確定）** | 通常リリース、Latest | 同じバージョン、dist-tag は `latest` |
| ホットフィックス | `v3.<current>.<N>`, `N > 0` | まれ、手動 | 通常リリース、Latest | 同じバージョン、dist-tag は `latest` |
| メジャー | `v4.0.0` | 破壊的変更がある場合のみ | 通常リリース | `latest` |

`<next>` は、**リリース済み**の最上位の v3 マイナーバージョンに 1 を加えたものです。nightly は**次回の weekly**のプレビューであるため、そのベースには最上位のリリースタグを使用します（プレリリースは無視されます）。hotfix `v3.1.1` の後も、次の nightly は `v3.2.0-nightly.*` のままです。プレリリースは対応する正式リリースより前に並ぶため、より上位のバージョンを指定する必要があります。`v3.0.0-nightly.<date>` は `v3.0.0` より前に並びます。

**マイナーバージョンの更新は、機能のマイルストーンではありません。** `v3.4.0` は「`v3.3.0`の次のウィークリーリリース」を意味します。変更内容はバージョン番号ではなく、リリースノートで確認してください。

ポリシーはバージョンの**形式**によって決まります。`-nightly.`サフィックスはナイトリーです。それ以外のサフィックス（`-alpha`、`-beta`、`-rc`）は、GA 前のレガシービルドです。サフィックスがないものはリリースです。

### 連動するバージョン

CLI、Go モジュール、npm パッケージ `@wailsio/runtime` は、`v3/internal/version/version.txt`から取得した**同一のバージョン**を使用します（タグ `v3.2.0-nightly.20261102` は npm では `3.2.0-nightly.20261102`）。`wails3 init`で作成したプロジェクトは、`@wailsio/runtime`を CLI 自体のバージョンに固定します。CLI の開発ビルドでは、代わりに `latest` を使用します。

## 具体例

| 日付 | イベント | タグ | `go install ...@latest` の解決先 |
| --- | --- | --- | --- |
| 25 Oct（日） | GA、手動リリース | `v3.0.0` | `v3.0.0` |
| 26 Oct（月）～31 Oct（土） | ナイトリー | `v3.1.0-nightly.20261026` ... `v3.1.0-nightly.20261031` | `v3.0.0` |
| 1 Nov（日）20:00 Sydney | ウィークリー | `v3.1.0` | `v3.1.0` |
| 2 Nov（月） | ナイトリー（次のウィークリーをプレビュー） | `v3.2.0-nightly.20261102` | `v3.1.0` |
| 4 Nov（水） | リリースブランチからの手動ホットフィックス | `v3.1.1` | `v3.1.1` |
| 8 Nov（日） | ウィークリー | `v3.2.0` | `v3.2.0` |

ウィークリーは常に、同じ番号のすべてのナイトリーよりも上位にソートされます。また、ナイトリーが `@latest` になることはありません。

## ユーザーが実行するもの

| 目的 | コマンド |
| --- | --- |
| 最新リリース（ウィークリーまたはホットフィックス） | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` または `wails3 update cli` |
| **再現可能（推奨）** | `go.mod`、CI、ドキュメントで正確なタグ `v3.1.0` に固定 |
| ナイトリー | 正確なタグ（`@v3.2.0-nightly.20261102`）、`wails3 update cli -nightly` |
| 種類を問わない最新ビルド | `wails3 update cli -pre` |
| ランタイム npm パッケージ | `@wailsio/runtime`（`latest` = ウィークリー／ホットフィックス、`nightly` = ナイトリー） |

この仕組みが成り立つ理由は、Go が `@latest` を、**プレリリースサフィックスのない**最も大きいタグに解決し、リリースが存在しない場合にのみプレリリースを考慮するためです。GitHub の「pre-release」チェックボックスは Go プロキシから認識できないため、ナイトリーはこのチェックボックスに依存せず、semver のプレリリースサフィックスを付けます。本番環境では必ずバージョンを固定してください。

## ブランチ

| ブランチ | 用途 |
| --- | --- |
| `master` | 開発およびナイトリー。 |
| `releases/**` | 単発リリース用の短期ブランチ：`v3.0.0` GA およびホットフィックス。`master`と同じ CI が必須です。 |

ウィークリーは `master` から作成されますが、**その先端から作成されるわけではありません**。ワークフローは、前回のウィークリー以降で、必須チェック（`v3 Build Results`、`Run Go Tests`、`Bundled runtime matches source`）がすべて成功した最新コミットを選び、そのコミット自体にタグを付けます。成功したコミットがなければ、何もリリースされません。ナイトリーは `master` から作成され、HEAD が正常でない場合でも警告するだけです。

## リリースの作成方法

リリース権限はプロジェクト Lead に限定されています（ガバナンス文書を参照）。コントリビューターがタグ付けや公開を行うことはありません。コントリビューターは `v3/UNRELEASED_CHANGELOG.md` にエントリを追加します。

すべての処理は **Nightly + Weekly Release v3**（`.github/workflows/nightly-release-v3.yml`）から実行されます。

- <strong>ナイトリー（自動）</strong>は 15:00 UTC に実行されます。未リリースの変更履歴が空の場合はスキップされます。
- <strong>Weekly（自動）</strong>はシドニー時間の 20:00 に実行されます。GitHub の cron は UTC のみに対応しているため、2 つの cron（UTC の 09:00 と 10:00）をスケジュールし、ガードジョブはシドニーで 20:00 の場合にのみ処理を続行します（DST 対応）。weekly はリリースゲートを通過し、成功済みの最新コミットを選択してマイナーバージョンを上げ、changelog に手を加えずにそのコミットへタグを付けます。その後、`master` 上のバージョンフィールド（`version.txt` とランタイムパッケージのバージョン）のみを更新するため、選択したコミットより後にマージされた依存関係と lockfile の更新は維持されます。npm への公開元は `master` ではなくタグです。
- **手動。** `channel`、`dry_run`（デフォルトでオン）、必要に応じて `release_version` を指定してワークフローを実行します。ドライランでは GitHub や npm に変更を加えず、バージョン、種類、GitHub フラグ、npm dist-tag、リリースノートを出力します。
- **GA または hotfix。** `release_version`（`v3.0.0`、`v3.1.1`）を設定します。古いリリースライン向けの hotfix では `ref: releases/...` も設定できます。そのリリースコミットは `master` にチェリーピックで戻されます。

ローカルでのドライラン: `cd v3/tasks/release && go run release.go --dry-run --channel weekly`。`go run release.go --print-policy v3.2.0-nightly.20261102`を指定すると、導出されたフラグが表示されます。

### 安全策

- 自動生成されるバージョンは、既存のすべてのタグよりも厳密に大きい必要があります（semver 順。プレリリース < リリース）。そうでない場合、ツールは処理を拒否します。
- リリースジョブは一度に 1 つだけ実行されます（共有 concurrency group）。
- npm への公開は、リリースの**タグ**によってトリガーされ、タグ付けされたコミットをビルドします。

## リリースゲート

**weekly**（nightly では実行されません）は、次のいずれかが存在する間、タグ付けを拒否します。

1. `P0`または`release-blocker`のラベルが付いた、未解決の issue または pull request。
2. 前回の weekly 以降に作成され、`regression`のラベルが付いた未解決の issue。
3. 重大度が critical または high の未解決の Dependabot アラート、重大度が critical または high の未解決の code-scanning アラート、および未解決の secret-scanning アラート。ただし、**出荷されるコードのみ**が対象です。`v3/examples/`、`v2/examples/`、`website/`、`docs/`、`.github/`、または任意の`examples/`ディレクトリのみに存在するアラートは無視されます（`v3/tasks/releasegate/main.go`内の`SHIPPED_ALERT_IGNORE`。**オーナーによるレビューが必要**）。
4. 依然として`triage`の状態にあるリポジトリセキュリティアドバイザリ。

**ブロッカーを設定できるユーザー。** ラベルがブロッカーとして扱われるのは、そのラベルの最新の `labeled` イベント（後続の `unlabeled` がないもの）が、`wailsapp/maintainers` チームのメンバー、または `maintain`/`admin` 権限を持つコラボレーターによって実行された場合のみです。それ以外のユーザーが同じラベルを適用しても無視され、ログに記録されるため、外部ユーザーがリリースをブロックしたり、ブロックを解除したりすることはできません。自動適用される `high-priority` ラベルは、意図的に使用対象に**含めていません**。`taliesin-ai` などのボットアカウントは、`maintain` を保持していれば許可リストを通過します。これは意図した動作であり、トリアージエージェントは根拠を添えてブロッカーラベルを適用できます。

**失敗時の動作。**ゲートはフェイル**クローズ**方式です。アラート API または許可リストを読み取れない場合、weekly は停止します。失敗時にはタグを一切作成せず、リンク付きの一覧をジョブサマリーに書き込みます。また、リポジトリ変数`WEEKLY_STATUS_ISSUE`が設定されている場合は、単一の「Weekly release status」issue 上のコメントを更新します。

**オーナーが GitHub 上で行う必要がある設定**（自動的には何も作成されません）: `P0`、`release-blocker`、`regression`の各ラベルを作成します。「Weekly release status」issue を作成し、その番号を`WEEKLY_STATUS_ISSUE`に設定します。`WAILS_REPO_TOKEN`に、Issues、Pull requests、Administration（コラボレーター権限）、Dependabot alerts、Code scanning alerts、Secret scanning alerts、Repository security advisories、organization Members の読み取り権限と、ステータスコメント用の Issues 書き込み権限を付与します。デフォルトの`GITHUB_TOKEN`では、これらを読み取れません。

**オーバーライド。**`override_reason`（workflow-dispatch の入力）を指定すると、ゲートをスキップします。理由、実行者、およびスキップされたすべてのブロッカーが、サマリーとステータス issue に書き込まれます。これはリリース Lead のみが使用してください。GitHub では入力の使用者を 1 人に制限できないため、オーバーライドに明示的な承認が必要となるよう、必須レビュー担当者を設定した environment でリリースジョブを実行することを推奨します。

## 関連項目

- [テストと継続的インテグレーション](/contributing/testing-ci/)
- [メンテナーによるトリアージ](/contributing/maintainer-triage/)
