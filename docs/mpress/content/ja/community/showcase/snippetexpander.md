---
title: "Snippet Expander"
description: "Wailsで構築されたデスクトップアプリケーション"
slug: "community/showcase/snippetexpander"
sourcePath: "community/showcase/snippetexpander.md"
---

![Snippet Expanderのスクリーンショット](/assets/showcase-images/snippetexpandergui-select-snippet.png)

Snippet Expanderの「スニペットを選択」ウィンドウのスクリーンショット

![Snippet Expanderのスクリーンショット](/assets/showcase-images/snippetexpandergui-add-snippet.png)

Snippet Expanderの「スニペットを追加」画面のスクリーンショット

![Snippet Expanderのスクリーンショット](/assets/showcase-images/snippetexpandergui-search-and-paste.png)

Snippet Expanderの「検索＆貼り付け」ウィンドウのスクリーンショット

[Snippet Expander](https://snippetexpander.org)は、Linux向けの「展開可能なテキストスニペットを支援する小さなツール」です。

Snippet Expanderは、スニペットと設定を管理するためにWailsで構築されたGUIアプリケーションです。スニペットをすばやく選択して貼り付けるための「検索＆貼り付け」ウィンドウモードも備えています。

WailsベースのGUI、go-lang CLI、vala-lang自動展開デーモンは、いずれもD-Busを介してgo-langデーモンと通信します。処理の大部分はこのデーモンが担い、スニペットのデータベースと共通設定を管理するほか、スニペットの展開や貼り付けなどのサービスを提供します。

[ソースコード](https://git.sr.ht/~ianmjones/snippetexpander/tree/trunk/item/cmd/snippetexpandergui/app.go#L38)を参照して、WailsアプリがUIからバックエンドへメッセージを送り、それがさらにデーモンへ送信される仕組みを確認してください。また、アプリの別インスタンスやCLIによるスニペットの変更を監視するためにD-Busイベントを購読し、その変更をWailsイベント経由で即座にUIへ反映する仕組みも確認できます。
