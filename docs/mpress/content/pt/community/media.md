---
title: "Mídia"
description: "Botões oficiais do Wails com pré-visualizações e códigos de incorporação em HTML e Markdown"
slug: "community/media"
sourcePath: "community/media.md"
---

Use os recursos oficiais do Wails ao compartilhar um projeto criado com o Wails. Os botões abaixo direcionam para [wails.io](https://wails.io) e continuam hospedados pelo Wails, portanto, não é necessário copiar a imagem para o seu projeto.

## Criado com Wails

Escolha a versão com contraste suficiente em relação ao plano de fundo. Cada pré-visualização corresponde exatamente à imagem vinculada que os visitantes verão.

### Plano de fundo claro

<div class="wails-media-preview wails-media-preview-light">  
<p>Padrão (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Criado com Wails" /></a>  
<p>Compacto (<code>90</code> × <code>24</code> px):</p>  
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

### Plano de fundo escuro

<div class="wails-media-preview wails-media-preview-dark">  
<p>Padrão (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Criado com Wails" /></a>  
<p>Compacto (<code>90</code> × <code>24</code> px):</p>  
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

O Markdown renderiza a imagem em seu tamanho intrínseco de `180` por `48` pixels. Use a versão em HTML quando precisar definir um tamanho de exibição diferente. Mantenha os botões em suas proporções originais de `15:4`.
