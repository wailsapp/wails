---
title: "Mídia"
description: "Botões oficiais do Wails com prévias e códigos de incorporação HTML e Markdown"
slug: "community/media"
sourcePath: "community/media.md"
---

Use os recursos visuais oficiais do Wails ao compartilhar um projeto criado com Wails. Os botões abaixo levam a [wails.io](https://wails.io) e continuam hospedados pelo Wails, então você não precisa copiar a imagem para o seu projeto.

## Criado com Wails

Escolha a versão que ofereça contraste suficiente com o fundo. Cada prévia mostra exatamente a imagem com link que seus visitantes verão.

### Fundo claro

<div class="wails-media-preview wails-media-preview-light">  
<p>Tamanho padrão (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Criado com Wails" /></a>  
<p>Tamanho compacto (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="Criado com Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

Para a versão compacta, use `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### Fundo escuro

<div class="wails-media-preview wails-media-preview-dark">  
<p>Tamanho padrão (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Criado com Wails" /></a>  
<p>Tamanho compacto (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="Criado com Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

Para a versão compacta, use `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

O Markdown exibe a imagem em seu tamanho original de `180` por `48` pixels. Use a versão HTML quando precisar definir outro tamanho de exibição. Mantenha a proporção original dos botões de `15:4`.
