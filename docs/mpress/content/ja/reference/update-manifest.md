---
title: "更新マニフェストプロトコル"
description: "Wails アプリが自己更新を検出して検証するために使用するオープンな JSON プロトコル。任意の静的ファイルホストまたは動的更新サーバーから配信できます。"
slug: "reference/update-manifest"
sourcePath: "reference/update-manifest.md"
---

Wails Update Manifest プロトコルは、Wails アプリケーションと更新元の間で使用する小規模でオープンな JSON コントラクトです。HTTPS 経由で JSON ファイルを配信できるものであれば、S3 バケット、GitHub Pages、CDN、ライセンスに基づいてリリースを制御する動的更新サーバーなど、どれでも Wails の更新を配信できます。

クライアント側は、フレームワークに `endpoint` プロバイダー（`github.com/wailsapp/wails/v3/pkg/updater/providers/endpoint`）として同梱されています。このページは、サーバー側を実装する開発者向けの通信形式リファレンスです。

## 設計目標

1. <strong>静的ホストに適しています。</strong>各プラットフォームの成果物を列挙した、チャネルごとに 1 つのマニフェストファイルだけで完全に実装できます。サーバーコードは不要です。
2. <strong>動的サーバーに適しています。</strong>クライアントは確認のたびに `platform`、`arch`、`version`、`channel` を送信します。そのため、サーバーは成果物を 1 つだけ返したり、ライセンスルールを適用したり、呼び出し元が最新版を使用している場合に `204 No Content` を返したりできます。
3. <strong>検証を最優先します。</strong>マニフェストには成果物ごとのチェックサムと署名が含まれ、Wails アップデーターは、ビルド時にアプリケーションバイナリへ固定された公開鍵を使用してそれらを検証します。更新元が独自の信頼の起点を選択することはありません。

## リクエスト

クライアントは、設定されたマニフェスト URL に対して、`Accept: application/json` と、アプリケーションで設定された任意のヘッダー（`Authorization: License <key>` など）を付けて `GET` を送信します。

URL にはプレースホルダーを埋め込むことができ、クライアントは確認のたびにそれらを置換します。

| プレースホルダー | 置換後の値 |
| --- | --- |
| `{{platform}}` | 実行中の OS を表す Go の `GOOS` 値（`darwin`、`windows`、`linux`） |
| `{{arch}}` | 実行中のアーキテクチャを表す Go の `GOARCH` 値（`amd64`、`arm64`、…） |
| `{{version}}` | 現在インストールされているバージョン |
| `{{channel}}` | 設定されている場合は、設定済みのリリースチャネル |

4 つの値のうち、プレースホルダーで使用されなかった値は、同名のクエリパラメーターとして追加されます（`channel` は設定されている場合のみ）。したがって、次の 2 つはいずれも有効で同等の設定です。

```text
# Dynamic server: reads query parameters
https://updates.example.com/check
  -> GET /check?platform=darwin&arch=arm64&version=1.0.0&channel=stable

# Static host: one manifest per platform/arch/channel path
https://cdn.example.com/updates/{{platform}}/{{arch}}/{{channel}}.json
  -> GET /updates/darwin/arm64/stable.json?version=1.0.0
```

静的ホストは、受信したクエリパラメーターを無視するだけです。

## レスポンス

| ステータス | 意味 |
| --- | --- |
| `200 OK` | 後続の内容はマニフェストです。アップグレードに該当するかどうかはクライアントが判断します。 |
| `204 No Content` | サーバーがバージョンを比較し、呼び出し元が最新版を使用していると判断しました。 |
| `404 Not Found` | 何も公開されていません（最新版を使用している場合と同様に扱われます）。 |
| その他すべて | エラーです。アップデーターは、次に設定されているプロバイダーへフォールスルーします。 |

`200` の本文はマニフェストドキュメントです。

```json
{
  "schemaVersion": 1,
  "version": "2.1.0",
  "channel": "stable",
  "name": "Summer Release",
  "notes": "## What's new\n\n- Faster startup\n- New themes",
  "publishedAt": "2026-07-03T10:00:00Z",
  "artifacts": [
    {
      "url": "MyApp-2.1.0-darwin-arm64.zip",
      "platform": "darwin",
      "arch": "arm64",
      "filetype": "zip",
      "size": 8388608,
      "digestAlgo": "sha512",
      "digest": "base64-encoded digest bytes",
      "signatureAlgo": "ed25519ph",
      "signature": "base64-encoded signature bytes"
    },
    {
      "url": "MyApp-2.1.0-windows-amd64.zip",
      "platform": "windows",
      "arch": "amd64",
      "filetype": "zip",
      "size": 9437184,
      "digestAlgo": "sha512",
      "digest": "...",
      "signatureAlgo": "ed25519ph",
      "signature": "..."
    }
  ]
}
```

### トップレベルのフィールド

| フィールド | 型 | 必須 | 注記 |
| --- | --- | --- | --- |
| `schemaVersion` | int | いいえ | プロトコルのバージョンです。省略した場合は `1` になります。クライアントは、自身が解釈できるものより新しい値を拒否します。 |
| `version` | string | **はい** | 先頭に `v` があってもなくてもよい SemVer 2.0.0。 |
| `channel` | string | いいえ | 情報提供用です。別のチャネルが設定されたクライアントは、そのマニフェストを更新なしとして扱います。 |
| `name` | string | いいえ | 更新ウィンドウに表示される、人が判読できるリリースタイトルです。 |
| `notes` | string | いいえ | 更新ウィンドウにレンダリングされる Markdown 形式のリリースノートです。 |
| `publishedAt` | string | いいえ | RFC 3339 形式のタイムスタンプです。 |
| `artifacts` | array | **はい** | ダウンロード可能な成果物ごとに 1 つのエントリを指定します。順序は公開者の優先順位を表します。 |
| `metadata` | object | いいえ | アプリケーションにそのまま渡される、自由形式のキーと値のデータです。 |

クライアントは未知のフィールドを無視するため、サーバーは互換性を損なうことなく独自のフィールドを追加できます。サーバー固有の追加情報は `metadata` に格納します。

### アーティファクトのフィールド

| フィールド | 型 | 必須 | 備考 |
| --- | --- | --- | --- |
| `url` | string | **はい** | 絶対 URL、またはマニフェスト URL からの相対 URL。`http(s)` のみ。 |
| `platform` | string | いいえ | Go の `GOOS` 値。一般的な別名（`macos`、`win` など）も使用できます。空の場合はすべてのプラットフォームに一致します。 |
| `arch` | string | いいえ | Go の `GOARCH` 値。一般的な別名（`x86_64`、`aarch64` など）も使用できます。空の場合はすべてのアーキテクチャに一致します。 |
| `filename` | string | いいえ | デフォルトは `url` のパスの最後のセグメントです。 |
| `filetype` | string | いいえ | デフォルトはファイル名の拡張子です。 |
| `size` | int | いいえ | バイト数。ダウンロードの進行状況に使用します。 |
| `digestAlgo` / `digest` | string / base64 | いいえ | `sha256` または `sha512`。 |
| `signatureAlgo` / `signature` | string / base64 | いいえ | `ed25519`、`ed25519ph`、または `ecdsa-p256`。`signature` が存在する場合は、必ず `signatureAlgo` が必要です。各アルゴリズムが署名する対象については、[アップデーターガイド](/guides/updater/#cryptographic-verification)を参照してください。 |

クライアントは、`platform` と `arch` が実行中のシステムに一致する<strong>最初の</strong>アーティファクトを選択します。Base64 値は、パディングの有無にかかわらず使用できます。

### バージョンの比較

マニフェストがアップグレードに該当するかどうかは、常にクライアント側で SemVer 2.0.0 の優先順位に従って判断されます。マニフェストの `version` は、インストール済みバージョンより厳密に新しくなければなりません。これにより、静的ホスティングでも容易に正しい動作を実現できます（マニフェストは常に最新リリースを記述し、すでに最新のクライアントは何もしません）。一方、動的サーバーでは帯域幅を節約するために `204` を引き続き利用できます。

## 検証と信頼

チェックサムと署名はマニフェストに含まれますが、信頼の起点は含まれません。署名は、ビルド時にアプリケーションが `updater.Config.PublicKey` を介して固定した公開鍵に対して検証されます。侵害された、または差し替えられた更新元が独自の鍵を提供することはできません。アプリケーションに固定された鍵がないにもかかわらずアーティファクトに署名が含まれている場合は、安全側に倒して処理に失敗します。宣言済みの `signatureAlgo` がない署名や、デコードできない署名も同様です。クライアントがダイジェストのみの検証へ暗黙にフォールバックすることはありません。

ダイジェストのみを持つアーティファクトは、ダイジェストの検査後にインストールされます。この検査は破損を防ぎますが、改ざん耐性については TLS とホスト自体の完全性に依存します。セキュリティ上重要なものには署名を付けて配布してください。

フレームワークの `ed25519ph` 方式でアーティファクトに署名する処理は、数行の Go コードで記述できます。

```go
digest := sha512.Sum512(artifactBytes)
sig, _ := privateKey.Sign(nil, digest[:], &ed25519.Options{Hash: crypto.SHA512})
manifest.Artifacts[i].DigestAlgo = "sha512"
manifest.Artifacts[i].Digest = base64.StdEncoding.EncodeToString(digest[:])
manifest.Artifacts[i].SignatureAlgo = "ed25519ph"
manifest.Artifacts[i].Signature = base64.StdEncoding.EncodeToString(sig)
```

実際には、このコードを自分で記述することはほとんどありません。CLI が代わりに処理します。

## wails3 CLI による公開

`wails3 updater` コマンドグループは、公開パイプライン全体を扱います。リリースに必要なコマンドは 3 つです。

```bash
# Once per application: create the signing keypair.
wails3 updater genkey
# updater.key      keep secret (CI secret store), signs every release
# updater.key.pub  embed in the app and pass as updater.Config.PublicKey

# Per release: digest, sign and describe every artifact in one manifest.
wails3 updater manifest -version 2.1.0 -channel stable \
    -key updater.key -notes-file notes.md \
    -url-prefix "https://cdn.example.com/myapp/2.1.0" \
    bin/updates/

# Before uploading: re-verify the files exactly as a shipped app would.
wails3 updater verify -manifest manifest.json -publickey updater.key.pub
```

`manifest` はファイルまたはディレクトリを受け付けます（鍵データ、`.json`、チェックサムおよびリリースノートのサイドカーファイルは自動的にスキップされます）。各アーティファクトをストリーミングしながら SHA-512 を計算し、`-key` が指定されている場合は Ed25519ph でダイジェストに署名します。また、`MyApp-2.1.0-darwin-arm64.zip` のような一般的なファイル名から `platform` と `arch` を推論します（`macOS`、`win64`、`x86_64`、`aarch64` などの一般的な別名も認識します。推論できないものについては警告が出力され、そのアーティファクトはすべてのプラットフォームに一致します）。`-url-prefix` を省略して相対 URL を出力し、マニフェストをアーティファクトと同じ場所にアップロードしてください。

`verify` は不一致が 1 つでもあると 0 以外の終了コードで終了するため、ビルドと公開の間に配置する CI ゲートとして自然に利用できます。マニフェストを独自に構築するサーバー向けには、`wails3 updater sign -key updater.key <files...>` が各ファイルの `digest`/`signature` フィールドを、独自のドキュメントへそのままマージできる JSON として出力します。

## 認証

認証はサーバー側の責務であり、プロトコルはヘッダーを伝送するだけです。クライアントは、マニフェストを要求するたびに設定済みのヘッダーを再送します。アーティファクトのダウンロードでは、アーティファクトの URL がマニフェストと同じホスト上にあり、かつ `https` から `http` へダウングレードしない場合に限り、`Authorization` ヘッダーを送信します。オリジンをまたぐリダイレクトやダウングレードを伴うリダイレクトでは、このヘッダーを必ず削除します。そのため、認証情報が CDN やオブジェクトストレージへ漏えいしたり、平文で送信されたりすることはありません。

ホスト型ライセンスサービスと自然に組み合わせられる、ライセンスによるアクセス制御の例：

```go
ep, _ := endpoint.New(endpoint.Config{
    URL:     "https://updates.example.com/check",
    Headers: map[string]string{"Authorization": "License " + licenseKey},
})
```

## クライアントの設定

`endpoint.Config` の完全なリファレンスと、このプロバイダーを GitHub、keygen.sh、AppCast の各プロバイダーとともにフォールバックチェーンへ組み込む方法については、[アップデーターガイド](/guides/updater/#providers)を参照してください。
