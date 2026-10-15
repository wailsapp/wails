---
title: "Créez des applications de bureau avec Go"
description: "Applications de bureau natives utilisant Go et les technologies web"
banner: {"content":"Wails v3 est actuellement en version bêta. \u003ca href=\"https://v2.wails.io\"\u003eVous recherchez la documentation de la v2 ?\u003c/a\u003e\n"}
template: "splash"
hero: {"actions":[{"icon":"right-arrow","link":"/quick-start/installation","text":"Bien démarrer","variant":"primary"},{"icon":"open-book","link":"/tutorials/03-notes-vanilla","text":"Voir le tutoriel","variant":"secondary"}],"image":{"alt":"Logo de Wails","dark":"/assets/wails-logo-dark.svg","light":"/assets/wails-logo-light.svg"},"tagline":"Créez de belles applications performantes avec Go et les technologies web modernes. Une seule base de code. Trois plateformes. Aucun navigateur.*"}
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

<span class="browser-footnote" data-browser-footnote>* sauf si vous y tenez vraiment</span>
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

## Démarrage rapide

@terminal{frame="macos" prompt="none"}
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard (experimental)
wails3 setup
@end

Vous êtes alors prêt à créer votre premier projet :

@terminal{frame="macos" prompt="none"}
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
@end

**Votre application est maintenant en cours d’exécution** avec le rechargement à chaud et des liaisons Go vers JS à typage sûr.

Vous rencontrez des problèmes avec `wails3 setup` ? Consultez le [guide d’installation manuelle](/quick-start/installation/).

## Votre application de bureau est déjà une application mobile

@terminal{frame="macos" prompt="none"}
wails3 task ios:run      # run your existing app on iOS Simulator
wails3 task android:run  # run your existing app on Android Emulator
@end

**Aucune modification de votre code Go.** Le même `main.go`, les mêmes services, le même frontend : Wails compile automatiquement le tout pour iOS et Android. Les fonctionnalités propres à chaque plateforme (retour haptique, boîtes de dialogue natives, zones sûres) sont disponibles lorsque vous en avez besoin, mais ne sont pas nécessaires pour publier votre application.

[Documentation mobile →](/guides/mobile/)

---

## Pourquoi Wails ?

@cards{cols="2"}
🚀 Des performances perceptibles par les utilisateurs
- Binaires d’environ 15 Mo contre 150 Mo pour Electron
- Environ 10 Mo de mémoire au repos contre plus de 100 Mo
- Démarrage en &lt;0.5 s contre 2-3 s
- Rendu natif à l’aide de la WebView du système d’exploitation
- Aucune surcharge due à un navigateur intégré

---
⚙ Expérience de développement
- Une seule base de code Go pour toutes les plateformes
- N’importe quel framework web — React, Vue, Svelte
- Rechargement à chaud pendant le développement
- Liaisons générées automatiquement pour appeler facilement Go depuis JavaScript
- IPC en mémoire. Aucun port réseau

---
✓ Prêt pour les applications de bureau
- Plusieurs fenêtres avec gestion de leur cycle de vie
- Menus natifs et zone de notification système
- Boîtes de dialogue de fichiers natives de la plateforme
- Intégration au système et raccourcis
- Outils de signature du code et de création de paquets

---
▣ Bureau et mobile
- Windows, macOS, Linux, iOS, Android
- Même base de code, aucune réécriture
- WebView native sur chaque plateforme
- Aucun port ouvert, aucun serveur localhost
- Fonctionnalités de la plateforme disponibles selon les besoins

@end

## Étapes suivantes

Ensuite : [créez une application complète](/tutorials/03-notes-vanilla/), parcourez les [exemples](https://github.com/wailsapp/wails/tree/master/v3/examples) ou consultez la [référence de l’API](/reference/application/). Vous migrez depuis la v2 ? Consultez le [guide de mise à niveau](/migration/v2-to-v3/).

@note{type="info" title="Wails v3 bêta"}
Wails v3 est un logiciel en version bêta doté d’une API de bureau stable. Des équipes l’utilisent déjà en production, mais devraient effectuer des tests approfondis avant le déploiement pendant que nous apportons les dernières finitions pour 3.0.

@end

@cards{cols="1"}
♥ Soutenez le développement de Wails
Wails est un logiciel libre et gratuit, créé par des développeurs pour des développeurs. Si Wails vous aide à créer des applications exceptionnelles, envisagez de soutenir la poursuite de son développement.

Votre parrainage contribue à maintenir le projet, à améliorer la documentation et à développer de nouvelles fonctionnalités qui profitent à toute la communauté.

[Devenir sponsor →](https://github.com/sponsors/leaanthony)

@end
