---
title: "BulletinBoard"
description: "Wails で構築されたデスクトップアプリケーション"
slug: "community/showcase/bulletinboard"
sourcePath: "community/showcase/bulletinboard.md"
---

![BulletinBoard](/assets/showcase-images/bboard.webp)

[BulletinBoard](https://github.com/raguay/BulletinBoard) アプリケーションは、静的メッセージを表示したり、スクリプトでユーザーから情報を取得するためのダイアログを表示したりできる、多用途のメッセージボードです。新しいダイアログを作成するための TUI があり、作成したダイアログは後でユーザーから情報を取得する際に使用できます。このアプリケーションはシステム上で常駐し、必要に応じて情報を表示した後、非表示になるように設計されています。私の環境では、ファイルを監視し、変更時にその内容を BulletinBoard へ送信するプロセスを使用しています。私のワークフローとの相性は抜群です。また、プログラムへ情報を送信するための[Alfred ワークフロー](https://github.com/raguay/MyAlfred/blob/master/Alfred%205/EmailIt.alfredworkflow)もあります。このワークフローは[EmailIt](https://github.com/raguay/EmailIt)との連携にも使用できます。
