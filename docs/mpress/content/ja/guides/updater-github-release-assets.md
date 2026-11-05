---
title: "アップデーター向け GitHub Release アセット"
description: "Wails アップデーターが GitHub Releases からアプリケーションの成果物を選択し、インストーラーパッケージを除外する仕組み。"
slug: "guides/updater-github-release-assets"
sourcePath: "guides/updater-github-release-assets.md"
---

GitHub Releases プロバイダーは、設定された `AssetMatcher` を使用してリリースアセットを選択します。`AssetMatcher` が `nil` の場合は、`github.DefaultAssetMatcher` を使用します。

## デフォルトの照合

デフォルトのマッチャーは、各アセットのファイル名から現在のプラットフォームとアーキテクチャを検索します。次のような一般的なアーキテクチャの別名を認識します。

- `amd64`、`x86_64`、`x64`
- `arm64`、`aarch64`
- `386`、`i386`、`x86`、`ia32`

署名やチェックサムなどのサイドカーファイルは無視されます。

## インストーラーアセット

GitHub Release には、アップデーターが使用するアプリケーションバイナリと、初回インストール用の一般的なインストーラーの両方が含まれる場合があります。デフォルトのマッチャーは、小文字に変換したファイル名が次の条件に該当するアセットを無視します。

- `-installer.` を含む
- `_installer.` を含む
- `installer.exe` と完全に一致する

たとえば、次の Windows アセットがあるとします。

```text
myapp-windows-amd64.exe
myapp-windows-amd64-installer.exe
```

`DefaultAssetMatcher` は `myapp-windows-amd64.exe` を選択し、インストーラーを無視します。これにより、アップデーターが実行中のアプリケーションを NSIS または同様の形式でパッケージ化されたインストーラー実行ファイルに置き換えることを防ぎます。

このチェックは意図的に限定されています。`installer` という単語を単に含むだけのアプリケーション名は、次の例を含め、引き続き有効です。

```text
myinstaller.exe
installer-tool-windows-amd64.exe
myinstaller-windows-amd64.zip
```

## カスタム命名規則

リリースアセットがプラットフォームとアーキテクチャに基づく命名規則に従っていない場合、またはインストーラーを別の方法で除外する必要がある場合は、`AssetMatcher` を設定します。

```go
import (
    "strings"

    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

gh, err := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, asset := range assets {
            name := strings.ToLower(asset.Name)
            if strings.Contains(name, req.Platform) &&
                strings.Contains(name, req.Arch) &&
                !strings.Contains(name, "-setup.") {
                return i
            }
        }
        return -1
    },
})
```

カスタムマッチャーは `DefaultAssetMatcher` を完全に置き換えるため、署名、チェックサム、インストーラー、およびアプリケーションのアップデートとしてインストールすべきでないその他のアセットを除外する役割を担います。
