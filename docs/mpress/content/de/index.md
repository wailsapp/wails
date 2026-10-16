---
title: "Desktop-Anwendungen mit Go entwickeln"
description: "Native Desktop-Anwendungen mit Go und Webtechnologien"
banner: {"content":"Wails v3 befindet sich derzeit in der Betaphase. \u003ca href=\"https://v2.wails.io\"\u003eDu suchst die Dokumentation zu v2?\u003c/a\u003e\n"}
template: "splash"
hero: {"actions":[{"icon":"right-arrow","link":"/quick-start/installation","text":"Erste Schritte","variant":"primary"},{"icon":"open-book","link":"/tutorials/03-notes-vanilla","text":"Tutorial ansehen","variant":"secondary"}],"image":{"alt":"Wails-Logo","dark":"/assets/wails-logo-dark.svg","light":"/assets/wails-logo-light.svg"},"tagline":"Entwickle ansprechende, leistungsfähige Anwendungen mit Go und modernen Webtechnologien. Eine Codebasis. Drei Plattformen. Keine Browser.*"}
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

<span class="browser-footnote" data-browser-footnote>* außer du möchtest es wirklich</span>
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

## Schnellstart

@terminal{frame="macos" prompt="none"}
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard (experimental)
wails3 setup
@end

Jetzt bist du bereit für dein erstes Projekt:

@terminal{frame="macos" prompt="none"}
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
@end

**Deine Anwendung wird jetzt ausgeführt** – mit Hot Reload und typsicheren Go-zu-JS-Bindings.

Gibt es Probleme mit `wails3 setup`? Sieh dir die [Anleitung zur manuellen Installation](/quick-start/installation/) an.

## Deine Desktop-Anwendung ist bereits eine mobile Anwendung

@terminal{frame="macos" prompt="none"}
wails3 task ios:run      # run your existing app on iOS Simulator
wails3 task android:run  # run your existing app on Android Emulator
@end

**Keine Änderungen an deinem Go-Code.** Dieselbe `main.go`, dieselben Dienste, dasselbe Frontend – Wails kompiliert alles automatisch für iOS und Android. Plattformspezifische Funktionen wie Haptik, native Dialoge und Safe Areas stehen bei Bedarf zur Verfügung, sind für die Auslieferung aber nicht erforderlich.

[Dokumentation für Mobilgeräte →](/guides/mobile/)

---

## Warum Wails?

@cards{cols="2"}
🚀 Leistung, die Benutzer bemerken
- Binärdateien mit ~15 MB gegenüber 150 MB bei Electron
- ~10 MB Grundspeicherbedarf gegenüber mehr als 100 MB
- &lt;0.5 s Startzeit gegenüber 2-3 s
- Natives Rendering mit der WebView des Betriebssystems
- Kein Overhead durch einen gebündelten Browser

---
⚙ Entwicklungskomfort
- Eine Go-Codebasis für alle Plattformen
- Beliebiges Webframework – React, Vue oder Svelte
- Hot Reload während der Entwicklung
- Automatisch generierte Bindings für einfache Go-Aufrufe aus JavaScript
- IPC im Arbeitsspeicher. Keine Netzwerkports

---
✓ Bereit für den Desktop
- Mehrere Fenster mit eigenen Lebenszyklen
- Native Menüs und System-Tray
- Plattformnative Dateidialoge
- Systemintegration und Tastenkürzel
- Werkzeuge für Codesignierung und Paketierung

---
▣ Desktop und Mobilgeräte
- Windows, macOS, Linux, iOS, Android
- Dieselbe Codebasis, kein Neuschreiben
- Native WebView auf jeder Plattform
- Keine offenen Ports, kein localhost-Server
- Plattformfunktionen bei Bedarf verfügbar

@end

## Nächste Schritte

Als Nächstes kannst du [eine vollständige Anwendung entwickeln](/tutorials/03-notes-vanilla/), die [Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples) durchsuchen oder in der [API-Referenz](/reference/application/) nachschlagen. Du migrierst von v2? Sieh dir den [Upgrade-Leitfaden](/migration/v2-to-v3/) an.

@note{type="info" title="Wails v3 Beta"}
Wails v3 ist eine Betaversion mit einer stabilen Desktop-API. Teams setzen sie bereits produktiv ein, sollten sie vor der Bereitstellung jedoch gründlich testen, während wir 3.0 den letzten Feinschliff geben.

@end

@cards{cols="1"}
♥ Unterstütze die Entwicklung von Wails
Wails ist kostenlos, quelloffen und wird von Entwicklern für Entwickler entwickelt. Wenn Wails dir hilft, beeindruckende Anwendungen zu entwickeln, erwäge, die weitere Entwicklung zu unterstützen.

Dein Sponsoring hilft dabei, das Projekt zu pflegen, die Dokumentation zu verbessern und neue Funktionen zu entwickeln, von denen die gesamte Community profitiert.

[Sponsor werden →](https://github.com/sponsors/leaanthony)

@end
