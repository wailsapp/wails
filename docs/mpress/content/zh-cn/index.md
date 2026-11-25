---
title: "使用 Go 构建桌面应用"
description: "使用 Go 和 Web 技术构建原生桌面应用"
banner: {"content": "Wails v3.0.0 是面向 Windows、macOS 和 Linux 桌面应用的稳定版本。\u003ca href=\"https://v2.wails.io\"\u003e在查找 v2 文档？\u003c/a\u003e\n"}
template: "splash"
hero: {"actions":[{"icon":"right-arrow","link":"/quick-start/installation","text":"开始使用","variant":"primary"},{"icon":"open-book","link":"/tutorials/03-notes-vanilla","text":"查看教程","variant":"secondary"}],"image":{"alt":"Wails 徽标","dark":"/assets/wails-logo-dark.svg","light":"/assets/wails-logo-light.svg"},"tagline":"使用 Go 和现代 Web 技术构建美观、高性能的应用。一个代码库，三个平台，无需浏览器。*"}
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

<span class="browser-footnote" data-browser-footnote>* 除非你确实想用浏览器</span>
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

## 快速开始

@terminal{frame="macos" prompt="none"}
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard (experimental)
wails3 setup
@end

接下来就可以创建你的第一个项目了：

@terminal{frame="macos" prompt="none"}
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
@end

**你的应用现已运行**，并支持热重载和类型安全的 Go 到 JS 绑定。

使用`wails3 setup`时遇到问题？请参阅[手动安装指南](/quick-start/installation/)。

## 你的桌面应用已经是移动应用

@terminal{frame="macos" prompt="none"}
wails3 task ios:run      # run your existing app on iOS Simulator
wails3 task android:run  # run your existing app on Android Emulator
@end

<strong>无需改动任何 Go 代码。</strong>使用相同的`main.go`、相同的服务和相同的前端，Wails 会自动将其编译到 iOS 和 Android。需要时可以使用平台特有功能（触觉反馈、原生对话框、安全区域），但发布应用并不依赖这些功能。

[移动端文档 →](/guides/mobile/)

---

## 为什么选择 Wails？

@cards{cols="2"}
🚀 用户切实可感的性能
- 二进制文件约为 15MB，而 Electron 为 150MB
- 基准内存占用约为 10MB，而对比方案为 100MB+
- 启动时间为&lt;0.5秒，而对比方案为 2-3秒
- 使用操作系统 WebView 进行原生渲染
- 没有捆绑浏览器带来的额外开销

---
⚙ 开发者体验
- 一个 Go 代码库适配所有平台
- 支持任意 Web 框架——React、Vue、Svelte
- 开发期间支持热重载
- 自动生成绑定，轻松从 Javascript 调用 Go
- 内存内 IPC，无需网络端口

---
✓ 桌面应用就绪
- 支持多个具有生命周期的窗口
- 原生菜单和系统托盘
- 平台原生文件对话框
- 系统集成和快捷键
- 代码签名和打包工具

---
▣ 桌面端与移动端
- Windows、macOS、Linux、iOS、Android
- 同一代码库，无需重写
- 每个平台均使用原生 WebView
- 不开放端口，也不运行 localhost 服务器
- 需要时可使用平台特有功能

@end

## 后续步骤

接下来：[构建一个完整应用](/tutorials/03-notes-vanilla/)、浏览[示例](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)，或查看[API 参考](/reference/application/)。要从v2迁移？请参阅[升级指南](/migration/v2-to-v3/)。

@note{type="info" title="Wails v3.0.0"}
Wails v3.0.0 是面向 Windows、macOS 和 Linux 桌面应用的稳定版本。Android 和 iOS 支持仍处于实验阶段，不在稳定桌面兼容性承诺范围内。

@end

@cards{cols="1"}
♥ 支持 Wails 的开发
Wails 是由开发者为开发者打造的免费开源项目。如果 Wails 帮助你构建了出色的应用，请考虑支持其持续开发。

你的赞助有助于维护项目、改进文档，并开发惠及整个社区的新功能。

[成为赞助者 →](https://github.com/sponsors/leaanthony)

@end
