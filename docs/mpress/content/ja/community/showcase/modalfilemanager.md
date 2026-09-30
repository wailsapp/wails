---
title: "Modal File Manager"
description: "Wailsで構築されたデスクトップアプリケーション"
slug: "community/showcase/modalfilemanager"
sourcePath: "community/showcase/modalfilemanager.md"
---

![Modal File Manager](/assets/showcase-images/modalfilemanager.webp)

[Modal File Manager](https://github.com/raguay/ModalFileManager)は、Web技術を使用したデュアルペイン型の ファイルマネージャーです。当初の設計はNW.jsをベースとしており、 [こちら](https://github.com/raguay/ModalFileManager-NWjs)で確認できます。この バージョンでも同じSvelteベースのフロントエンドコードを使用しています（ただし、NW.jsから 移行して以来、大幅に変更されています）が、バックエンドには [Wails 2](https://wails.io/)による実装を採用しています。この実装により、コマンドラインの `rm`や`cp`などのコマンドは使用しなくなりましたが、テーマや拡張機能を ダウンロードするには、システムにgitをインストールしておく必要があります。すべてGoで 実装されており、以前のバージョンよりもはるかに高速に動作します。

このファイルマネージャーは、Vimと同じ原則、すなわち状態によって 制御されるキーボード操作を中心に設計されています。状態の数は固定されておらず、非常に柔軟に プログラムできます。そのため、無限に多くのキーボード設定を 作成して使用できます。これが、ほかのファイルマネージャーとの主な違いです。GitHubから ダウンロードできるテーマや拡張機能も用意されています。
