---
title: "難読化ビルド"
description: "Garble で Wails アプリケーションをビルドし、ソースコードをリバースエンジニアリングから保護します"
slug: "guides/build/obfuscation"
sourcePath: "guides/build/obfuscation.md"
---

[Garble](https://github.com/burrowers/garble) は、`go build` の代わりに使用し、シンボル名の変更、定数の難読化、生成されるバイナリからのデバッグ情報の削除を行う Go ビルドツールです。Wails v3 では、2 つの新しいコマンドによって Garble を標準でサポートしています。

## 前提条件

- **Go 1.26.2 以降** — Garble v0.16.0 に必要です
- **Garble v0.16.0**

```bash
go install mvdan.cc/garble@v0.16.0
```

@note{type="tip"}
Garble が必要とする Go の最小バージョンは、リリースごとに異なります。古い Go ツールチェーンを使用している場合は、インストール前に[Garble のリリースページ](https://github.com/burrowers/garble/releases)で、お使いのツールチェーンに対応するバージョンを確認してください。

@end

## 必須：サービス型に JSON タグを追加する

バインドされたサービスメソッドが返す、または受け取るすべての構造体では、エクスポートされた各フィールドに明示的な JSON タグを付ける必要があります：

```go
// Without tags — breaks under Garble
type OrderSummary struct {
    ID        int
    Total     float64
    LineItems []LineItem
}

// With tags — safe under Garble
type OrderSummary struct {
    ID        int       `json:"id"`
    Total     float64   `json:"total"`
    LineItems []LineItem `json:"lineItems"`
}
```

@note{type="caution"}
Garble はエクスポートされた構造体フィールドの名前を変更します。また、Wails はこれらの構造体を、Garble が静的に追跡できない `interface{}` パラメーターを通じて `json.Marshal` に渡します。JSON タグがなくても難読化ビルドは正常にコンパイルされますが、実行時にフロントエンドが受け取るフィールド名は難読化された名前になるか、空になります。難読化ビルドを実行する前にタグを追加してください。

@end

Wails 独自の型（`Screen`、`Rect`、`Point`、`Size`、`EnvironmentInfo`、`OSInfo`、`Capabilities`）には、すでにタグが付いています。タグを付ける必要があるのは、独自に定義した型だけです。

## 難読化してビルドする

@steps
### 安定 ID ファイルを生成する
バインドされたサービスメソッドを追加、名前変更、または削除するたびに、次のコマンドを実行してください：

```bash
wails3 generate bindings -obfuscated
```

これにより、メインパッケージのディレクトリに `wails_obfuscated.gen.go` が作成されます。このファイルをコミットしてください。

### Garble でビルドする
```bash
wails3 build --obfuscated
```

難読化されたバインディングを使用してアプリケーションをビルドします。

@end

## Garble に追加フラグを渡す

オプションを `garble` に直接渡すには、`--garbleargs` を使用します：

```bash
# Obfuscate string literals and reduce binary size
wails3 build --obfuscated --garbleargs "-literals -tiny"

# Reproducible output — same seed produces the same binary
wails3 build --obfuscated --garbleargs "-seed=deadbeef"
```

サポートされているフラグの全一覧については、[Garble のドキュメント](https://github.com/burrowers/garble#flags)を参照してください。

## 高度な設定：ID ファイルを別のパッケージに書き込む

デフォルトでは、`wails_obfuscated.gen.go` は `main` パッケージと同じ場所に書き込まれます。プロジェクトで、`main` からインポートされるサブパッケージにサービスを配置している場合は、`-obfuscated-output` を使用して、そのサブパッケージにファイルを書き込むこともできます：

```bash
wails3 generate bindings -obfuscated -obfuscated-output ./internal/services
```

@note{type="caution"}
起動時にそのパッケージの `init()` が実行されるように、出力先パッケージは `main` パッケージから直接または推移的にインポートされている必要があります。出力先パッケージに到達できない場合、安定 ID は登録されず、バインディング呼び出しは失敗します（たとえば、実行時に `binding not found` エラーが発生します）。

@end

## トラブルシューティング

### `garble: command not found`

Garble がインストールされていないか、`$(go env GOPATH)/bin` が `PATH` に含まれていません。

```bash
go install mvdan.cc/garble@v0.16.0
export PATH="$PATH:$(go env GOPATH)/bin"
```

### フロントエンドが誤ったフィールド値または空のフィールド値を受け取る

サービスの戻り値の型に `json:"..."` タグがありません。バインドされたメソッドが返す各構造体を確認し、エクスポートされたすべてのフィールドに明示的なタグを追加してください。

### ブラウザーコンソールに `binding not found` エラーが表示される

安定 ID ファイルが存在しないか、コンパイルに含まれていません。以下を確認してください：

- `wails_obfuscated.gen.go` がメインパッケージのディレクトリ（または `-obfuscated-output` に渡したディレクトリ）に存在する
- `wails_obfuscated` ビルドタグを追加する `wails3 build --obfuscated` を実行した
- `-obfuscated-output` を使用した場合、出力先パッケージが `main` からインポートされている

### Windows Defender がビルドをウイルスとして検出する

Garble で難読化された Go バイナリは、デバッグシンボルがなく、パックされた実行可能ファイルに類似しているため、ビルド中に Windows Defender のヒューリスティック検出によって脅威と判定されます。ビルドは次のエラーで失敗します：

```
open C:\Users\...\AppData\Local\Temp\go-build...\a.out.exe: The file contains a virus or potentially unwanted software.
```

一時ディレクトリ（Go が中間ビルド成果物を書き込む場所）とプロジェクトディレクトリを Defender の除外リストに追加してください：

```powershell
Add-MpPreference -ExclusionPath "$env:TEMP"
Add-MpPreference -ExclusionPath "C:\path\to\your\project"
```

これらの除外設定は指定したパスにのみ適用され、Defender をシステム全体で無効にするものではありません。

### `unsupported Go version` によりビルドが失敗する

Garble v0.16.0 には Go 1.26.2 以降が必要です。Go をアップグレードするか、[Garble のリリースページ](https://github.com/burrowers/garble/releases)で、お使いのツールチェーンと互換性のあるバージョンを確認してください。
