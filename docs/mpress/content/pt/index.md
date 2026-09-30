---
title: "Crie aplicativos para desktop com Go"
description: "Aplicativos nativos para desktop usando Go e tecnologias web"
banner: {"content":"O Wails v3 está atualmente em versão beta. \u003ca href=\"https://v2.wails.io\"\u003eProcurando a documentação da v2?\u003c/a\u003e\n"}
template: "splash"
hero: {"actions":[{"icon":"right-arrow","link":"/quick-start/installation","text":"Comece agora","variant":"primary"},{"icon":"open-book","link":"/tutorials/03-notes-vanilla","text":"Ver tutorial","variant":"secondary"}],"image":{"alt":"Logotipo do Wails","dark":"/assets/wails-logo-dark.svg","light":"/assets/wails-logo-light.svg"},"tagline":"Crie aplicativos bonitos e de alto desempenho usando Go e tecnologias web modernas. Uma única base de código. Três plataformas. Nenhum navegador.*"}
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

<span class="browser-footnote" data-browser-footnote>* a menos que você realmente queira</span>
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

## Início rápido

@terminal{frame="macos" prompt="none"}
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard (experimental)
wails3 setup
@end

Agora você já pode criar seu primeiro projeto:

@terminal{frame="macos" prompt="none"}
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
@end

**Seu aplicativo agora está em execução** com recarregamento automático e bindings de Go para JS com segurança de tipos.

Está tendo problemas com `wails3 setup`? Consulte o [guia de instalação manual](/quick-start/installation/).

## Seu aplicativo para desktop já é um aplicativo móvel

@terminal{frame="macos" prompt="none"}
wails3 task ios:run      # run your existing app on iOS Simulator
wails3 task android:run  # run your existing app on Android Emulator
@end

**Não é necessário alterar o código Go.** O mesmo `main.go`, os mesmos serviços, o mesmo frontend — o Wails compila tudo automaticamente para iOS e Android. Os recursos específicos de cada plataforma (resposta tátil, caixas de diálogo nativas e áreas seguras) estão disponíveis quando você quiser usá-los, mas não são necessários para lançar o aplicativo.

[Documentação para dispositivos móveis →](/guides/mobile/)

---

## Por que usar o Wails?

@cards{cols="2"}
🚀 Desempenho que os usuários percebem
- Binários de ~15 MB, em comparação com 150 MB do Electron
- Consumo básico de memória de ~10 MB, em comparação com mais de 100 MB
- Tempo de inicialização de &lt;0.5 s, em comparação com 2-3 s
- Renderização nativa usando a WebView do sistema operacional
- Sem a sobrecarga de incluir um navegador no pacote

---
⚙ Experiência de desenvolvimento
- Uma única base de código Go para todas as plataformas
- Qualquer framework web — React, Vue ou Svelte
- Recarregamento automático durante o desenvolvimento
- Bindings gerados automaticamente para chamar Go facilmente a partir do JavaScript
- IPC na memória. Sem portas de rede

---
✓ Pronto para desktop
- Várias janelas com ciclos de vida próprios
- Menus nativos e área de notificação do sistema
- Caixas de diálogo de arquivos nativas da plataforma
- Integração com o sistema e atalhos
- Ferramentas de assinatura de código e empacotamento

---
▣ Desktop e dispositivos móveis
- Windows, macOS, Linux, iOS e Android
- A mesma base de código, sem reescritas
- WebView nativa em todas as plataformas
- Sem portas abertas nem servidor localhost
- Recursos da plataforma disponíveis quando necessário

@end

## Próximas etapas

A seguir: [crie um aplicativo completo](/tutorials/03-notes-vanilla/), explore os [exemplos](https://github.com/wailsapp/wails/tree/master/v3/examples) ou consulte a [referência da API](/reference/application/). Está migrando da v2? Consulte o [guia de atualização](/migration/v2-to-v3/).

@note{type="info" title="Wails v3 Beta"}
O Wails v3 é um software beta com uma API estável para desktop. Algumas equipes já o utilizam em produção, mas devem realizar testes completos antes da implantação enquanto concluímos os ajustes finais para 3.0.

@end

@cards{cols="1"}
♥ Apoie o desenvolvimento do Wails
O Wails é gratuito e de código aberto, criado por desenvolvedores para desenvolvedores. Se o Wails ajuda você a criar aplicativos incríveis, considere apoiar a continuidade de seu desenvolvimento.

Seu patrocínio ajuda a manter o projeto, aprimorar a documentação e desenvolver novos recursos que beneficiam toda a comunidade.

[Torne-se um patrocinador →](https://github.com/sponsors/leaanthony)

@end
