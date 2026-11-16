---
title: "フィードバック"
description: "Wails v3 にフィードバックを提供し、問題を報告する方法"
slug: "feedback"
sourcePath: "feedback.md"
---

皆さんからのフィードバックを歓迎（そして推奨）します。新しい Issue やディスカッションを作成する前に、既存のものがないか検索してください。貢献する方法は次のとおりです。

@tabs
[バグ]
バグを見つけた場合は、バグ報告テンプレートを使用して GitHub で[ Issue を作成](https://github.com/wailsapp/wails/issues/new/choose)してください。

- 簡単に再現できる例を添えて、バグを明確に説明してください。何が起こる*べきか*がドキュメントで明確になっていない場合は、その点も報告に含めてください。
- `wails3 doctor` の出力を報告に含めてください。
- 現在のドキュメントに記載された動作と一致しないことがバグである場合は、以下も行ってください。
  - `v3/examples` ディレクトリにある既存のサンプルを更新するか、問題を明確に示す新しいサンプルを作成してください。
  - Issue を参照する[ PR](https://github.com/wailsapp/wails/pulls) を作成してください。


@note{type="caution"}
*覚えておいてください*。予期しない動作が必ずしもバグとは限りません。単に期待どおりに動作していないだけかもしれません。その場合は `Suggestions` を利用してください。

@end

Discord の [#v3](https://discord.gg/bdj28QNHmT) チャンネルでのバグに関する議論も歓迎します。

[修正]
バグの修正やドキュメントの改善がある場合は、次の手順を行ってください。

- [コントリビューションガイドライン](https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md)に従って、[Wails リポジトリ](https://github.com/wailsapp/wails)でプルリクエストを作成してください。
- 関連する Issue がある場合は、PR の説明で参照してください。

[機能強化]
新機能および公開動作の変更は、機能リクエストの Issue ではなく、<strong>WEP（Wails Enhancement Proposal）</strong>のドラフトプルリクエストを通じて提案します。

- [WEP のプロセス](https://github.com/wailsapp/wails/blob/v3.0.0/v3/wep/README.md)を確認してください。
- テンプレートをコピーし、WEP と補足資料のみを含む、タイトルが `[WEP] <title>` のドラフト PR を作成してください。
- 最初に[GitHub Discussions](https://github.com/wailsapp/wails/discussions)または Discord の [#v3](https://discord.gg/bdj28QNHmT) チャンネルでアイデアについて非公式に議論しても構いませんが、メンテナーが判断を下すには WEP の PR が必要です。

[賛成票]
- GitHub の :thumbsup: リアクションを使用して、バグ、WEP、ディスカッションへの支持を示してください。
- 「+1」や「私もです」のようなコメントを追加するだけの行為は*しないで*ください。
- 「このバグは ARM ビルドにも影響します」や「別の方法としては……」など、有意義な情報を提供できる場合はコメントを追加してください。

@end

既知の問題と進行中の作業は[こちら](https://github.com/orgs/wailsapp/projects/6)で確認できます。
