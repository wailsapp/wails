---
title: "Go でデスクトップアプリを構築"
description: "Go と Web 技術を使用したネイティブデスクトップアプリケーション"
banner: {"content":"Wails v3 は現在ベータ版です。\u003ca href=\"https://v2.wails.io\"\u003ev2 のドキュメントをお探しですか？\u003c/a\u003e\n"}
template: "splash"
hero: {"actions":[{"icon":"right-arrow","link":"/quick-start/installation","text":"はじめる","variant":"primary"},{"icon":"open-book","link":"/tutorials/03-notes-vanilla","text":"チュートリアルを見る","variant":"secondary"}],"image":{"alt":"Wails のロゴ","dark":"/assets/wails-logo-dark.svg","light":"/assets/wails-logo-light.svg"},"tagline":"Go と最新の Web 技術を使用して、美しく高性能なアプリケーションを構築できます。1 つのコードベース。3 つのプラットフォーム。ブラウザー不要。*"}
sourcePath: "index.md"
---

<style>
  /* Hero background — the neon "digital Wales" mountain, full-bleed and fixed. */
  body::after {
    content: '';
    position: fixed;
    inset: 0;
    z-index: -1;
    background: url('/digital_wales_master.webp') center center / cover no-repeat;
    opacity: 0.35;
    pointer-events: none;
  }
  
  /* Gradient overlay - highest z-index */
  html::before {
    content: '';
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: 
      linear-gradient(to bottom, var(--sl-color-bg) 0%, transparent 15%, transparent 85%, var(--sl-color-bg) 100%),
      linear-gradient(to right, var(--sl-color-bg) 0%, transparent 10%, transparent 90%, var(--sl-color-bg) 100%);
    z-index: 0;
    pointer-events: none;
  }
  
  
  /* Hero padding */
  .hero {
    padding-top: 3rem !important;
    padding-bottom: 2rem !important;
  }
  
  .small-buttons {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
    margin-top: 1rem;
  }
  
  .small-buttons a {
    display: inline-block;
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    border-radius: 6px;
    text-decoration: none;
    background: var(--sl-color-gray-6);
    color: var(--sl-color-white);
    border: 1px solid var(--sl-color-gray-5);
    transition: all 0.2s;
  }
  
  .small-buttons a:hover {
    background: var(--sl-color-gray-5);
    border-color: var(--sl-color-accent);
  }
  
  /* Round card corners and add translucent blur effect */
  .mpress-card {
    border-radius: 12px;
    background: rgba(var(--sl-color-gray-6-rgb), 0.6) !important;
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
  }

</style>

<span class="browser-footnote" data-browser-footnote>* どうしても使いたい場合を除く</span>
<script>
(() => {
  const words = ['Desktop', 'Server', 'Mobile'];
  let index = 0;
  const init = () => {
    const tagline = document.querySelector('.mpress-frontmatter-hero-tagline');
    const note = document.querySelector('[data-browser-footnote]');
    if (tagline && note && !tagline.contains(note)) tagline.appendChild(note);
    const title = document.querySelector('.mpress-frontmatter-hero-title');
    if (!title || title.querySelector('.morph-word-title')) return;
    if (title.textContent.trim() !== 'Build Desktop Apps with Go') return;
    const titleLine = document.createElement('span');
    titleLine.className = 'morph-title-line';
    const prefix = document.createElement('span');
    prefix.textContent = 'Build';
    const wordWrap = document.createElement('span');
    wordWrap.className = 'morph-word-title';
    const activeWord = document.createElement('span');
    activeWord.textContent = words[0];
    wordWrap.append(activeWord);
    titleLine.append(prefix, wordWrap);
    title.replaceChildren(titleLine, document.createElement('br'), document.createTextNode('Apps with Go'));
    const rotate = () => {
      const word = title.querySelector('.morph-word-title span');
      if (!word) return;
      word.classList.add('morph-word-out');
      setTimeout(() => {
        index = (index + 1) % words.length;
        word.textContent = words[index];
        word.classList.remove('morph-word-out');
        setTimeout(rotate, 3600);
      }, 650);
    };
    setTimeout(rotate, 3600);
  };
  document.readyState === 'loading' ? document.addEventListener('DOMContentLoaded', init) : init();
})();
</script>

## クイックスタート

@terminal{frame="macos" prompt="none"}
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard (experimental)
wails3 setup
@end

これで、最初のプロジェクトを作成する準備が整いました。

@terminal{frame="macos" prompt="none"}
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
@end

**アプリケーションが起動しました**。ホットリロードと型安全な Go-to-JS バインディングが有効になっています。

`wails3 setup`で問題が発生していますか？[手動インストールガイド](/quick-start/installation/)を参照してください。

## デスクトップアプリがそのままモバイルアプリに

@terminal{frame="macos" prompt="none"}
wails3 task ios:run      # run your existing app on iOS Simulator
wails3 task android:run  # run your existing app on Android Emulator
@end

**Go コードの変更は一切不要です。** 同じ`main.go`、同じサービス、同じフロントエンドを、Wails が iOS および Android 向けに自動的にコンパイルします。プラットフォーム固有の機能（触覚フィードバック、ネイティブダイアログ、セーフエリア）も必要に応じて利用できますが、リリースに必須ではありません。

[モバイル向けドキュメント →](/guides/mobile/)

---

## Wails を選ぶ理由

@cards{cols="2"}
🚀 ユーザーが実感できるパフォーマンス
- バイナリサイズは約15MB（Electron は150MB）
- ベースラインメモリ使用量は約10MB（比較対象は100MB 以上）
- 起動時間は&lt;0.5秒（比較対象は2-3秒）
- OS の WebView を使用したネイティブレンダリング
- ブラウザーを同梱するオーバーヘッドなし

---
⚙ 開発者体験
- すべてのプラットフォームに対応する単一の Go コードベース
- React、Vue、Svelte など、あらゆる Web フレームワークに対応
- 開発中のホットリロード
- 自動生成されたバインディングにより、Javascript から Go を簡単に呼び出し可能
- インメモリ IPC。ネットワークポートは使用しません

---
✓ デスクトップ対応
- ライフサイクルを備えた複数ウィンドウ
- ネイティブメニューとシステムトレイ
- 各プラットフォームのネイティブなファイルダイアログ
- システム連携とショートカット
- コード署名およびパッケージングツール

---
▣ デスクトップ＆モバイル
- Windows、macOS、Linux、iOS、Android
- 同じコードベースを使用し、書き直しは一切不要
- すべてのプラットフォームでネイティブ WebView を使用
- 開いているポートも localhost サーバーもありません
- 必要に応じてプラットフォーム固有の機能を利用可能

@end

## 次のステップ

次は、[完全なアプリケーションを構築する](/tutorials/03-notes-vanilla/)、[サンプル](https://github.com/wailsapp/wails/tree/master/v3/examples)を見る、または[API リファレンス](/reference/application/)を確認してください。v2 から移行する場合は、[アップグレードガイド](/migration/v2-to-v3/)を参照してください。

@note{type="info" title="Wails v3 ベータ版"}
Wails v3 は、安定したデスクトップ API を備えたベータ版ソフトウェアです。すでに本番環境で使用しているチームもありますが、私たちが3.0に向けた最終調整を進めている間は、デプロイ前に十分なテストを行うことをお勧めします。

@end

@cards{cols="1"}
♥ Wails の開発を支援
Wails は、開発者が開発者のために作る、無料のオープンソースソフトウェアです。Wails が優れたアプリケーションの構築に役立っている場合は、継続的な開発への支援をご検討ください。

スポンサーからの支援は、プロジェクトの維持、ドキュメントの改善、そしてコミュニティ全体に役立つ新機能の開発に活用されます。

[スポンサーになる →](https://github.com/sponsors/leaanthony)

@end
