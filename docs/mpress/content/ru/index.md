---
title: "Создавайте настольные приложения на Go"
description: "Нативные настольные приложения на Go с использованием веб-технологий"
banner: {"content":"Wails v3 сейчас находится на стадии бета-тестирования. \u003ca href=\"https://v2.wails.io\"\u003eИщете документацию по v2?\u003c/a\u003e\n"}
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

<span class="browser-footnote" data-browser-footnote>* если только они вам действительно не нужны</span>
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
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard (experimental)
wails3 setup
@end

Теперь всё готово для создания первого проекта:

@terminal{frame="macos" prompt="none"}
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
@end

**Ваше приложение уже запущено** с горячей перезагрузкой и типобезопасными привязками Go к JS.

Возникли проблемы с `wails3 setup`? Обратитесь к [руководству по установке вручную](/quick-start/installation/).

## Ваше настольное приложение уже является мобильным

@terminal{frame="macos" prompt="none"}
wails3 task ios:run      # run your existing app on iOS Simulator
wails3 task android:run  # run your existing app on Android Emulator
@end

**Никаких изменений в коде Go.** Те же `main.go`, те же сервисы, тот же фронтенд — Wails автоматически компилирует всё это для iOS и Android. Возможности конкретных платформ (тактильная отдача, нативные диалоговые окна, безопасные области) доступны, когда они вам нужны, но для выпуска приложения они не обязательны.

[Документация по мобильным приложениям →](/guides/mobile/)

---

## Почему Wails?

@cards{cols="2"}
🚀 Производительность, которую замечают пользователи
- Бинарные файлы размером ~15 МБ против 150 МБ у Electron
- ~10 МБ памяти в состоянии покоя против 100 МБ и более
- Запуск за &lt;0.5 с против 2-3 с
- Нативный рендеринг с помощью системного WebView
- Нет накладных расходов на встроенный браузер

---
⚙ Удобство разработки
- Единая кодовая база на Go для всех платформ
- Любой веб-фреймворк — React, Vue, Svelte
- Горячая перезагрузка во время разработки
- Автоматически создаваемые привязки для удобного вызова Go из Javascript
- IPC в памяти. Никаких сетевых портов

---
✓ Готовность к настольным платформам
- Несколько окон с управлением жизненным циклом
- Нативные меню и системный трей
- Нативные для каждой платформы диалоговые окна выбора файлов
- Интеграция с системой и сочетания клавиш
- Инструменты для подписи кода и упаковки приложений

---
▣ Настольные и мобильные платформы
- Windows, macOS, Linux, iOS, Android
- Одна кодовая база — ничего не нужно переписывать
- Нативный WebView на каждой платформе
- Никаких открытых портов и серверов localhost
- Возможности платформ доступны по мере необходимости

@end

## Дальнейшие шаги

Далее: [создайте полноценное приложение](/tutorials/03-notes-vanilla/), просмотрите [примеры](https://github.com/wailsapp/wails/tree/master/v3/examples) или изучите [справочник по API](/reference/application/). Переходите с v2? Ознакомьтесь с [руководством по обновлению](/migration/v2-to-v3/).

@note{type="info" title="Бета-версия Wails v3"}
Wails v3 — бета-версия со стабильным API для настольных приложений. Команды уже используют её в рабочей среде, однако перед развёртыванием следует провести тщательное тестирование, пока мы завершаем финальную доработку для 3.0.

@end

@cards{cols="1"}
♥ Поддержите разработку Wails
Wails — бесплатный проект с открытым исходным кодом, созданный разработчиками для разработчиков. Если Wails помогает вам создавать замечательные приложения, подумайте о поддержке его дальнейшей разработки.

Ваша спонсорская поддержка помогает сопровождать проект, улучшать документацию и разрабатывать новые функции на благо всего сообщества.

[Стать спонсором →](https://github.com/sponsors/leaanthony)

@end
