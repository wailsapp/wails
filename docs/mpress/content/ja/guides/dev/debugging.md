---
title: "デバッグ"
description: "問題の調査とアプリのパフォーマンス分析"
slug: "guides/dev/debugging"
sourcePath: "guides/dev/debugging.md"
---

このガイドでは、次のものを使用して、Wails アプリで発生する可能性のあるパフォーマンス上の問題を検査・調査するためのさまざまなツールを紹介します。

- [`runtime/trace`](https://pkg.go.dev/runtime/trace) を使用してパフォーマンスグラフを作成し、ブラウザーで確認します

## パフォーマンストレースの作成

@steps
### トレース用にアプリを準備する
プログラムのエントリーポイント付近に、次のようなコードが存在することを確認します

```go
// Create the file to store our trace data within
traceFile, err := os.Create("trace.out")
if err != nil {
  log.Fatalf("trace.out could not be created: %v", err)
}

// Start the trace
if err := trace.Start(traceFile); err != nil {
  _ = traceFile.Close()
  log.Fatalf("trace.start could not start: %v", err)
}

// Trace cleanup on exit. Alternatively,
defer func() {
  trace.Stop()
  _ = traceFile.Close()
}()

...Start your wails app here...
```

出力は、アプリの実行時間全体を対象とするメトリクスとともに、作業パス内の `trace.out` に保存されます。デフォルトのトレースはかなり未加工です。より多くのコンテキスト情報を追加し、実際に記録する内容を制限する方法について、トレースに関する資料をさらに読むことを強く推奨します。例：`WithRegion, NewTask, Log`

### アプリを実行してトレースを作成する
アプリの実行中は終了時点まで継続的にトレースが生成されるため、アプリでいくつかの操作を行い、完了したら終了してください

### 可視化ヘルパーをインストールする
一部のビューには [Graphviz](https://graphviz.org/) が必要です。`dot -V` を実行して、インストールされていることを確認できます

```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz

# Arch
sudo pacman -S graphviz
```

### トレースを可視化する
トレースデータを用意できたら、Web インターフェースを起動できます

```bash
go tool trace trace.out
```

これにより、デフォルトのブラウザー（Chrome ベースを推奨）でトレースイベントビューアーのホームページが開くはずです。

初めて使用する場合は、`Syscall profile` 画面が最も役立つでしょう。この画面には、プログラムが具体的に何を実行しているかが、所要時間とともに詳しく表示されます

@end
