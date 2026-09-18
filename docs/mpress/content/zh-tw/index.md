---
title: "使用 Go 建置桌面應用程式"
description: "使用 Go 與網頁技術打造原生桌面應用程式"
banner: {"content":"Wails v3 目前為測試版。\u003ca href=\"https://v2.wails.io\"\u003e正在尋找 v2 文件？\u003c/a\u003e\n"}
template: "splash"
hero: {"actions":[{"icon":"right-arrow","link":"/quick-start/installation","text":"開始使用","variant":"primary"},{"icon":"open-book","link":"/tutorials/03-notes-vanilla","text":"查看教學","variant":"secondary"}],"image":{"alt":"Wails 標誌","dark":"/assets/wails-logo-dark.svg","light":"/assets/wails-logo-light.svg"},"tagline":"使用 Go 與現代網頁技術建置美觀且高效能的應用程式。一套程式碼。三個平台。不含瀏覽器。*"}
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

<span class="browser-footnote" data-browser-footnote>* 除非您真的想要</span>
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

## 快速入門

@terminal{frame="macos" prompt="none"}
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard (experimental)
wails3 setup
@end

接著，您就可以建立第一個專案：

@terminal{frame="macos" prompt="none"}
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
@end

**您的應用程式現在已開始執行**，並支援熱重新載入和型別安全的 Go-to-JS 繫結。

使用`wails3 setup`時遇到問題？請參閱[手動安裝指南](/quick-start/installation/)。

## 您的桌面應用程式已經是行動應用程式

@terminal{frame="macos" prompt="none"}
wails3 task ios:run      # run your existing app on iOS Simulator
wails3 task android:run  # run your existing app on Android Emulator
@end

<strong>完全不必修改 Go 程式碼。</strong>同一個`main.go`、相同的服務、相同的前端——Wails 會自動將其編譯至 iOS 和 Android。需要時可使用平台特定功能（觸覺回饋、原生對話方塊、安全區域），但即使不用，也能發布應用程式。

[行動版文件 →](/guides/mobile/)

---

## 為什麼選擇 Wails？

@cards{cols="2"}
🚀 使用者能感受到的效能
- 約 15MB 的二進位檔，Electron 則為 150MB
- 基準記憶體用量：約10MB，對比100MB+
- 啟動時間：&lt;0.5秒，對比2-3秒
- 使用作業系統 WebView 進行原生繪製
- 沒有隨附瀏覽器所造成的額外負擔

---
⚙ 開發者體驗
- 所有平台共用一套 Go 程式碼
- 可使用任何網頁框架——React、Vue、Svelte
- 開發期間支援熱重新載入
- 自動產生繫結，讓 JavaScript 能輕鬆呼叫 Go
- 記憶體內 IPC，不使用網路連接埠

---
✓ 完整支援桌面應用程式
- 支援多個視窗及其生命週期
- 原生選單和系統匣
- 平台原生檔案對話方塊
- 系統整合和快速鍵
- 程式碼簽署和封裝工具

---
▣ 桌面與行動平台
- Windows、macOS、Linux、iOS、Android
- 共用同一套程式碼，完全不必重寫
- 每個平台皆使用原生 WebView
- 不開放任何連接埠，也不使用 localhost 伺服器
- 需要時即可使用平台功能

@end

## 後續步驟

接下來：[建置完整的應用程式](/tutorials/03-notes-vanilla/)、瀏覽[範例](https://github.com/wailsapp/wails/tree/master/v3/examples)，或查看[API 參考文件](/reference/application/)。要從 v2 遷移嗎？請參閱[升級指南](/migration/v2-to-v3/)。

@note{type="info" title="Wails v3 測試版"}
Wails v3 是測試版軟體，具備穩定的桌面 API。已有團隊將其用於 正式環境，但在我們完成3.0的最後潤飾期間，仍應在部署前 進行徹底測試。

@end

@cards{cols="1"}
♥ 支持 Wails 開發
Wails 是由開發者為開發者打造的免費開放原始碼軟體。如果 Wails 能協助您建置出色的應用程式，請考慮支持其持續開發。

您的贊助有助於維護專案、改善文件，並開發造福整個社群的新功能。

[成為贊助者 →](https://github.com/sponsors/leaanthony)

@end
