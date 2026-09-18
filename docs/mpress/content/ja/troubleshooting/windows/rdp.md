---
title: "リモート デスクトップ（RDP）経由で WebView2 がハングする"
description: "セッション中にモニターの DPI が変化する RDP セッション経由で Wails アプリを実行した際に、WebView2 の UI が数秒間停止する問題を解決します。"
slug: "troubleshooting/windows/rdp"
sourcePath: "troubleshooting/windows/rdp.md"
---

## 問題

Wails アプリケーションをリモート デスクトップ（RDP）セッション経由で使用すると、一般的な操作で UI が数秒間停止することがあります。

- ポップアップ ウィンドウを開くと、クリックしてから内容が表示されるまでに約 4～8 秒かかります。
- ウィンドウを閉じると、親ウィンドウが約 2 秒間ブロックされます。
- この低速な状態は再接続後も続き、ホスト マシンを再起動するまで解消されません。

この問題は、セッションの途中で Retina 向けに最適化された仮想モニターを構成する Microsoft Remote Desktop の iOS クライアントで特によく発生します。セッションの途中で異なる DPI コンテキストのモニターを追加する RDP クライアントであれば、どれでも同じ動作を引き起こす可能性があります。

## 発生する理由

WebView2 はデフォルトでウィンドウ ホスティングを使用し、コンポジター サーフェスは子ウィンドウ内に配置されます。RDP クライアントがセッションとは異なる DPI コンテキストのモニターを追加すると、WebView2 コントローラーの各呼び出し（`PutIsVisible`、`MoveFocus`、初回描画、サーフェスの解放）によって、同期的な DirectComposition の再マーシャリングが強制されます。再マーシャリングのたびに UI スレッドが約 2 秒間ブロックされるため、ポップアップを多用するアプリでは停止時間が積み重なります。

同じマシン上のネイティブ Win32 と WebView2 を組み合わせたアプリケーションは、ビジュアル ホスティングを使用するため影響を受けません。このことから、原因は WebView2 や Windows コンポジター全般の問題ではなく、ホスティング モードにあると判断できます。

## 解決策

Windows オプションで `UseVisualHosting` を設定し、ビジュアル ホスティングを有効にします。ビジュアル ホスティングでは、WebView2 のコンポジター サーフェスはホストが所有する DirectComposition ビジュアルを介して管理されるため、DPI コンテキストが変化しても同期的な再マーシャリングは発生しなくなります。

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Windows: application.WindowsOptions{
            UseVisualHosting: true,
        },
    })

    // ... create your windows, then:
    app.Run()
}
```

有効にすると、ポップアップは通常のナビゲーション時間内（約 150～500 ミリ秒）に開き、ウィンドウを閉じても親ウィンドウがブロックされなくなります。

@note{type="caution"}
`UseVisualHosting` は `app.Run()` より前に設定する必要があります。Wails はアプリケーションの起動時にこの値を読み取り、WebView2 環境を初期化する前に環境変数 `COREWEBVIEW2_FORCED_HOSTING_MODE` を `COREWEBVIEW2_HOSTING_MODE_WINDOW_TO_VISUAL` に設定します。それより後に設定しても効果はありません。

@end

このオプションのデフォルト値は `false` であるため、ウィンドウ ホスティングが引き続きデフォルトです。既存のアプリでは、明示的に有効にしない限り動作は変わりません。

## 有効にする場合

アプリが RDP 経由、特に Microsoft Remote Desktop の iOS クライアント経由で頻繁に使用され、ウィンドウを開閉する際に数秒間の停止が発生する場合は、`UseVisualHosting: true` を設定してください。アプリを RDP 経由で実行しない場合、このオプションは不要であり、デフォルトのままにできます。

## 参考資料

- [WebView2：ウィンドウ ホスティングとビジュアル ホスティング](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/windowed-vs-visual-hosting)
- [WebView2Feedback の課題 #5248](https://github.com/MicrosoftEdge/WebView2Feedback/issues/5248)
- [WebView2Feedback の課題 #4485](https://github.com/MicrosoftEdge/WebView2Feedback/issues/4485)
