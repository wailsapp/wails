---
title: "macOS 上の Syso ファイル"
description: "macOS での Syso ファイルのビルドエラーを解決する"
slug: "troubleshooting/mac-syso"
sourcePath: "troubleshooting/mac-syso.md"
---

## 問題

macOS で Wails アプリケーションをビルドしようとすると、次のようなエラーが発生してビルドに失敗します。

```
Error: Users/runner/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.1.darwin-arm64/pkg/tool/darwin_arm64/link: running clang failed: exit status 1
ld: unknown file type in '/private/var/folders/ml/x_tvfgn50_s7p67dm1ypcqqm0000gn/T/go-link-774134794/000000.o'
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

## なぜこの問題が発生するのですか？

Windows 向けにビルドすると、Wails はアプリケーションアイコン、ウィンドウアイコン、メニューアイコン用の `.syso` ファイルを生成します。これらのファイルは、Windows でアプリケーションをビルドするために必要です。 これらの `.syso` ファイルはプロジェクトのルートディレクトリにあり、macOS 向けのビルド時に問題を引き起こす可能性があります。

## 解決方法

この問題を解決するには、プロジェクトのルートディレクトリから syso ファイルを削除します。または、自分で生成している場合は、`wails_windows_<arch>.syso` という形式で命名します（例：`wails_windows_arm64.syso`）。
