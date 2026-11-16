---
title: "Создавайте настольные приложения с помощью Go"
description: "Нативные настольные приложения на базе Go и веб-технологий"
banner: {"content": "Wails v3.0.0 — стабильный выпуск для настольных приложений под Windows, macOS и Linux. \u003ca href=\"https://v2.wails.io\"\u003eИщете документацию по v2?\u003c/a\u003e\n"}
template: "splash"
hero: {"actions":[{"icon":"right-arrow","link":"/quick-start/installation","text":"Начать работу","variant":"primary"},{"icon":"open-book","link":"/tutorials/03-notes-vanilla","text":"Открыть руководство","variant":"secondary"}],"image":{"alt":"Логотип Wails","dark":"/assets/wails-logo-dark.svg","light":"/assets/wails-logo-light.svg"},"tagline":"Создавайте красивые и производительные приложения с помощью Go и современных веб-технологий. Одна кодовая база. Три платформы. Никаких браузеров.*"}
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

<span class="browser-footnote" data-browser-footnote>* если только вам действительно этого не хочется</span>
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

## Быстрый старт

@terminal{frame="macos" prompt="none"}
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard (experimental)
wails3 setup
@end

Теперь всё готово для создания первого проекта:

@terminal{frame="macos" prompt="none"}
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
@end

**Ваше приложение уже запущено** с горячей перезагрузкой и типобезопасными привязками Go к JS.

Возникли проблемы с `wails3 setup`? См. [руководство по установке вручную](/quick-start/installation/).

## Ваше настольное приложение уже является мобильным

@terminal{frame="macos" prompt="none"}
wails3 task ios:run      # run your existing app on iOS Simulator
wails3 task android:run  # run your existing app on Android Emulator
@end

**Никаких изменений в коде Go.** Тот же `main.go`, те же сервисы, тот же фронтенд — Wails автоматически компилирует всё это для iOS и Android. Платформенные функции (тактильная обратная связь, нативные диалоговые окна, безопасные области) доступны, когда они вам нужны, но для выпуска приложения они не требуются.

[Документация для мобильных платформ →](/guides/mobile/)

---

## Почему Wails?

@cards{cols="2"}
🚀 Производительность, заметная пользователям
- Размер исполняемых файлов — около 15 МБ против 150 МБ у Electron
- Базовое потребление памяти — около 10 МБ против 100 МБ и более
- Время запуска — &lt;0.5 с против 2-3 с
- Нативный рендеринг с помощью системного WebView
- Никаких накладных расходов от встроенного браузера

---
⚙ Удобство разработки
- Единая кодовая база Go для всех платформ
- Любой веб-фреймворк — React, Vue, Svelte
- Горячая перезагрузка во время разработки
- Автоматически создаваемые привязки для удобного вызова Go из Javascript
- IPC в памяти. Никаких сетевых портов

---
✓ Готово для настольных приложений
- Несколько окон с управлением жизненным циклом
- Нативные меню и системный трей
- Платформенные нативные диалоговые окна выбора файлов
- Интеграция с системой и сочетания клавиш
- Инструменты для подписания кода и упаковки приложений

---
▣ Настольные и мобильные платформы
- Windows, macOS, Linux, iOS, Android
- Та же кодовая база, никаких переписываний
- Нативный WebView на каждой платформе
- Никаких открытых портов и сервера localhost
- Платформенные функции доступны при необходимости

@end

## Дальнейшие шаги

Далее: [создайте полноценное приложение](/tutorials/03-notes-vanilla/), изучите [примеры](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples) или ознакомьтесь со [справочником API](/reference/application/). Переходите с v2? См. [руководство по обновлению](/migration/v2-to-v3/).

@note{type="info" title="Wails v3.0.0"}
Wails v3.0.0 — стабильный выпуск для настольных приложений под Windows, macOS и Linux. Поддержка Android и iOS остаётся экспериментальной и не входит в гарантии совместимости стабильной версии для настольных платформ.

@end

@cards{cols="1"}
♥ Поддержите разработку Wails
Wails — бесплатный проект с открытым исходным кодом, созданный разработчиками для разработчиков. Если Wails помогает вам создавать замечательные приложения, поддержите его дальнейшую разработку.

Ваша спонсорская поддержка помогает сопровождать проект, улучшать документацию и разрабатывать новые функции на благо всего сообщества.

[Стать спонсором →](https://github.com/sponsors/leaanthony)

@end
