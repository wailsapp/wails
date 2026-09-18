---
title: "クラッシュのデバッグ"
description: "Wails アプリケーションから診断レポートとミニダンプを収集する"
slug: "guides/debugging-crashes"
sourcePath: "guides/debugging-crashes.md"
---

本番環境で問題が発生した場合、ログだけでは十分でないことがよくあります。このガイドでは `github.com/wailsapp/wails/v3/pkg/debug` パッケージを説明します。このパッケージは構造化された診断レポートと、Windows では WinDbg や Visual Studio で開ける完全なミニダンプを生成します。

このパッケージは意図的に小さく保たれており、2 つのエントリーポイントがあります。

- `debug.Dump(...)` — 現在のプロセスの Windows ミニダンプを書き出します。
- `debug.Report(...)` — 詳細な診断スナップショットを収集し、任意でミニダンプも生成します。

どちらも呼び出したユーザーのトークンで動作します。権限の昇格も `SeDebugPrivilege` も `OpenProcess` も使用しません。現在のプロセスのダンプには疑似ハンドル `GetCurrentProcess` を使い、自分自身のプロセスへのアクセスでは ACL チェックを回避します。

## クイックリファレンス

```go
import "github.com/wailsapp/wails/v3/pkg/debug"

// 1. Just a minidump (Windows only).
path, err := debug.Dump()

// 2. Dump to a specific path.
path, err := debug.Dump(debug.WithPath("C:\\crashes\\my.dmp"))

// 3. Full-memory minidump (large, gigabytes).
path, err := debug.Dump(debug.WithFullMemory())

// 4. Full diagnostic report, no minidump.
r, err := debug.Report()

// 5. Report + minidump.
r, err := debug.Report(debug.WithDump())

// 6. Report + minidump at a specific path, full-memory.
r, err := debug.Report(
    debug.WithDumpPath("C:\\crashes\\my.dmp"),
    debug.WithDumpFullMemory(),
)
```

`debug.Report` は部分的な失敗でも必ず `*CrashReport` を返します。呼び出し側はエラーまでに収集できた情報を受け取れるため、`err != nil` の場合でもレポートが `nil` になることはありません。

## `CrashReport` 構造体

```go
type CrashReport struct {
    Timestamp   time.Time          // when the report was generated
    System      SystemInfo         // OS, hardware, CPU, GPU from doctor
    Build       BuildInfo          // Go version, buildmode, compiler, CGO flag
    Crash       *CrashInfo         // process, memory, modules, env vars
    Diagnostics []DiagnosticResult // doctor's health-check results
    DumpPath    string             // set only if WithDump() was passed
}
```

クラッシュ報告サービスに送る場合は `encoding/json` でシリアライズします。一部の情報だけが必要なら、フィールドを直接参照することもできます。

## `PanicHandler` との統合

最も実用的な方法は、Wails がパニックを捕捉した瞬間にレポートと必要に応じてミニダンプを収集し、それらをまとめて独自の報告処理に渡すことです。

```go
app := application.New(application.Options{
    PanicHandler: func(pd *application.PanicDetails) {
        // Always: capture a diagnostic snapshot plus a minidump.
        report, reportErr := debug.Report(debug.WithDump())
        if reportErr != nil {
            log.Printf("debug.Report: %v", reportErr)
        }

        // Now you have:
        //   pd.Error, pd.StackTrace, pd.FullStackTrace  — from wails
        //   report.DumpPath                              — minidump (Windows)
        //   report.System / report.Crash / report.Build — context
        //
        // Ship it off (Sentry, S3, support ticket, local log...).
        mycrashservice.Upload(pd, report)
    },
})
```

Wails は `debug.Report` を自動的には**呼び出しません**。クラッシュ報告にはユーザーの同意や個人情報の除去が必要になることが多いため、判断はハンドラーで行います。Wails が `PanicHandler` を呼び出すタイミングと、自分のコードで起動したゴルーチンも対象にする方法については、[パニック処理ガイド](/guides/panic-handling/)を参照してください。

## 手動でのデバッグ

`debug.Dump` と `debug.Report` はパニック以外にも役立ちます。

- ハングの検出：ウォッチドッグが作動したらプロセスをダンプします。WinDbg で開くと、各スレッドが何を待って停止していたかを詳しく調べられます。
- 異常状態の報告：例えばゴルーチン数が際限なく増えている場合は、レポートを取得してアプリケーションを動かし続け、後で調査できます。
- 要求に応じたサポート資料：「バグを報告」メニューを `debug.Report(debug.WithDump())` に接続し、結果をユーザーのサポートチケットに添付します。

## Windows ミニダンプの実際の使い方

ミニダンプの既定のフラグは、詳しい情報をコンパクトに記録する組み合わせです。

- `MiniDumpNormal`
- `MiniDumpWithThreadInfo`
- `MiniDumpWithHandleData`
- `MiniDumpWithUnloadedModules`

これにより、すべてのスレッドスタック、ハンドルテーブル、モジュール一覧を保持する約 5～50 MB のダンプが生成され、取得時に各ゴルーチンがどこにいたかを再構成できます。

`debug.WithFullMemory()`、または `Report` の `debug.WithDumpFullMemory()` を使うと、プロセスのアドレス空間全体を取得する `MiniDumpWithFullMemory` が追加されます。特定のポインターが破損した原因を調べる場合に有用ですが、ファイルは数百 MB から数 GB になります。詳細な調査用に限定してください。

生成された `.dmp` は `windbg -z C:\path\to\your.dmp` で WinDbg に読み込むか、Visual Studio の「ファイル → クラッシュダンプを開く」で開きます。リリースビルドでは通常シンボルが削除されるため、同じビルドの `.pdb`、またはデバッグ情報付きでビルドした Go バイナリ自体を併用すると有用なスタックトレースが得られます。

## Windows 以外のプラットフォーム

`debug.Dump` は Windows 以外では未実装を示すエラーを返します。Linux や macOS には `GOTRACEBACK=crash` によるコアダンプ、`rr`、プラットフォーム固有のクラッシュ報告機構など、別の事後解析ツールがあり、`MiniDumpWriteDump` API にそのまま対応しません。`debug.Report` は引き続き動作し、ダンプファイル以外の情報を取得できます。

`Dump` のエラーは確認できるように設計されているため、機能を制限して適切に処理を継続できます。

```go
path, err := debug.Dump()
if err != nil {
    log.Printf("minidump unavailable on %s: %v", runtime.GOOS, err)
    // Continue with just debug.Report() or pd.StackTrace from wails.
}
```

## セキュリティ上の注意

- **権限の昇格は不要です。** このパッケージは、認証情報を収集するツールで使われる `SeDebugPrivilege` や `OpenProcess` / 別プロセスの PID によるアクセスを意図的に避けています。ダンプ対象は呼び出し元プロセスだけです。
- **ダンプファイルにはメモリ上のすべてが含まれます。** 使用中の認証情報、認証トークン、ユーザーデータ、暗号鍵も含みます。完全なメモリイメージとして扱い、転送中は保護し、保存済みの機密情報を除去し、第三者に送信する前にも秘匿化を検討してください。
- **既定のパスは `os.TempDir()` です。** Windows では `%LOCALAPPDATA%\Temp` に対応し、ユーザーごとに分かれていて誰でも読めるわけではありませんが、ファイルは残ります。保持するダンプは安全な場所に移してください。
