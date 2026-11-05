---
title: "Windows UAC の構成"
description: "Windows 向け Wails アプリケーションのユーザーアカウント制御（UAC）を構成する"
slug: "guides/windows-uac"
sourcePath: "guides/windows-uac.md"
---

対象プラットフォーム：<span class="mpress-badge mpress-badge-note">Windows</span>

<br/>

Windows のユーザーアカウント制御（UAC）は、Wails アプリケーションの実行権限を決定します。Wails v3 アプリケーションの Windows マニフェストには、デフォルトで明示的な UAC 構成が含まれているため、異なるマシンでも一貫して動作します。

## UAC の実行レベル

Windows アプリケーションは、マニフェストファイルを通じて異なる実行レベルを要求できます。Wails v3 にはデフォルトの実行レベルを指定した UAC 構成が自動的に含まれ、アプリケーションの要件に応じてカスタマイズできます。

### 利用可能な実行レベル

| レベル | 説明 | 用途 |
| --- | --- | --- |
| `asInvoker` | 親プロセスと同じ権限で実行される | ほとんどのアプリケーションのデフォルト |
| `highestAvailable` | ユーザーが利用できる最も高い権限で実行される | 昇格されたアクセス権が必要になる可能性があるアプリケーション |
| `requireAdministrator` | 常に管理者権限を必要とする | システムユーティリティ、インストーラー |

### デフォルト構成

Wails v3 アプリケーションの Windows マニフェストには、次のデフォルト UAC 構成が含まれています。

```xml
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

この構成により、アプリケーションは次のように動作します。

- 起動元プロセスと同じ権限で実行される
- デフォルトでは権限の昇格を必要としない
- 異なるマシンでも一貫して動作する
- 通常のユーザーに対して UAC プロンプトを表示しない

## UAC 構成のカスタマイズ

Wails v3 ではビルドアセットのカスタマイズが推奨されているため、Windows マニフェストテンプレートを直接編集して UAC 構成を変更できます。

### マニフェストテンプレートの場所

Windows マニフェストテンプレートは次の場所にあります。

```
build/windows/wails.exe.manifest
```

### 実行レベルの変更

実行レベルを変更するには、`requestedExecutionLevel` 要素の `level` 属性を編集します。

```xml {title="build/windows/wails.exe.manifest"}
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

### 例

#### 標準アプリケーション（デフォルト）

ほとんどのアプリケーションでは、デフォルトの `asInvoker` レベルを使用してください。

```xml
<requestedExecutionLevel level="asInvoker" uiAccess="false"/>
```

#### システムユーティリティ

利用可能な場合に昇格されたアクセス権を必要とするアプリケーション：

```xml
<requestedExecutionLevel level="highestAvailable" uiAccess="false"/>
```

#### 管理ツール

常に管理者権限を必要とするアプリケーション：

```xml
<requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
```

## UI アクセス

`uiAccess` 属性は、アプリケーションがより高い権限の UI 要素と対話できるかどうかを制御します。ほとんどの場合、`false` のままにしてください。

アプリケーションで次の操作が必要な場合に限り、`true` に設定してください。

- ほかのアプリケーションに入力を送信する
- ほかのアプリケーションの UI を操作する
- より高い権限で実行されるプロセスの UI 要素にアクセスする

@note{type="caution" title="UI アクセスの要件"}
`uiAccess="true"` に設定するには、アプリケーションが次の要件を満たす必要があります。

- 信頼された認証局が発行した証明書でデジタル署名されている
- 安全な場所（Program Files または Windows\System32）にインストールされている

@end

## カスタム UAC 設定でのビルド

マニフェストテンプレートを変更した後、通常どおりアプリケーションをビルドします。

```bash
wails3 build
```

ビルドプロセスにより、カスタム UAC 構成が実行可能ファイルへ自動的に埋め込まれます。

## UAC 構成の確認

`go-winres` ツールを使用して、UAC 設定が正しく埋め込まれていることを確認できます。

```bash
go-winres extract --in your-app.exe --out extracted-resources/
```

次に、抽出したマニフェストファイルを調べ、UAC 構成が含まれていることを確認します。

@note{type="tip" title="マニフェストの永続性"}
ほかの一部のフレームワークとは異なり、Wails v3 の UAC 構成はコンパイル時に実行可能ファイルへ直接埋め込まれるため、アプリケーションをほかのマシンにコピーしても保持されます。

@end

## トラブルシューティング

### UAC プロンプトが表示されない

`requireAdministrator` に設定しても UAC プロンプトが表示されない場合：

- マニフェストが実行可能ファイルに正しく埋め込まれていることを確認する
- すでに権限が昇格されたプロセスから実行していないことを確認する
- マニフェストの構文が有効な XML であることを確認する

### アプリケーションが起動しない

UAC の変更後にアプリケーションが起動しなくなった場合：

- マニフェストの構文に XML エラーがないか確認します
- 実行レベルの値が有効であることを確認します
- 問題を切り分けるため、`asInvoker` に戻してみます

### マシン間で動作が異なる場合

マシンによって UAC の動作が異なる場合：

- マニフェストが外部ファイルではなく実行可能ファイルに埋め込まれていることを確認します
- ビルド後に実行可能ファイルが変更されていないことを確認します
- 対象マシンで Windows の UAC 設定が有効になっていることを確認します
