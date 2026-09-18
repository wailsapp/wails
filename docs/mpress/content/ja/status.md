---
title: "プロジェクトの状況"
description: "Wails v3 ベータ版の互換性、セキュリティサポート、アップグレード手順"
slug: "status"
sourcePath: "status.md"
---

## 現在の状況：ベータ版

最新の状況については、[変更履歴](/changelog/)を確認してください。

目標は安定版 v3.0 のリリースです。Wails v2 は引き続き現在の安定版であり、修正が提供されます。デプロイ前に、ご自身のアプリケーションでベータ版をテストしてください。

## ベータ版の互換性保証

v3ベータ版の互換性保証は、デスクトップアプリケーションを対象とします。

| プラットフォーム | サポート対象 | 要件と注記 |
| --- | --- | --- |
| Windows | amd64およびarm64 | WebView2ランタイム |
| macOS | IntelおよびApple Silicon | インストールガイドに記載されているmacOSおよびWebKitのバージョン |
| Linux | amd64およびarm64 | デフォルトではGTK4 + WebKitGTK 6.0。GTK3 + WebKit2GTK 4.1はv3.0.xまで`-tags gtk3`レガシーオプションとして残り、v3.1で削除されます |

すべてのターゲットで、開発にはGo 1.25以降が必要です。AndroidおよびiOSのサポートは実験的なものであり、デスクトップ版のベータ提供を妨げるものではありません。ベータ版APIは安定性を目指していますが、プレリリース版の不具合や明示的に告知された変更については、v3.0.0までに修正される場合があります。

## コントリビューション方法

- 最新のベータリリースをテストし、再現可能なバグを報告する
- ドキュメントやサンプルの作成に貢献する
- 議論に参加し、WEP草案にフィードバックを提供する
- バグ修正、ドキュメント、または承認済みのWEPに関するプルリクエストを提出する

コミュニティからのコントリビューションを歓迎します。これらの目標にご協力いただける場合は、コミュニティの議論に参加してください。新機能の提案は、機能リクエストのIssueではなく、WEPのPRとして提出してください。

## フィードバックと更新

再現可能な問題はIssueとして報告し、新機能は[WEP（Wails Enhancement Proposal）](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)のPRを通じて提案してください。

## ベータ版の利用

`latest` を追従するのではなく、CLI、Go モジュール、フロントエンドランタイムのバージョンを明示的に固定してください。既存のアルファ版プロジェクトについては、[アルファ版からベータ版へのアップグレードガイド](/migration/alpha-to-beta/)に従ってください。

[セキュリティポリシー](https://github.com/wailsapp/wails/blob/master/SECURITY.md)では、v3 ベータ版はサポート対象、アルファ版はサポート対象外とされています。脆弱性は公開 Issue ではなく、[非公開の脆弱性報告](https://github.com/wailsapp/wails/security/advisories/new)から報告してください。

## 追跡中の作業

- [v3 ラベルの付いた未解決のバグ](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3ABug+label%3Av3)
- [P0 または P1 ラベルの付いた未解決の v3 Issue](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3Av3+label%3AP0%2CP1)
- [リリースのマイルストーン](https://github.com/wailsapp/wails/milestones)

これらのリアルタイム検索は Issue のラベルに依存します。完全な一覧でも、リリース日や範囲を約束するものでもありません。Issue を読んで、ご自身のプロジェクトへの影響を評価してください。
