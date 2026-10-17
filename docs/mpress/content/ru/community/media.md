---
title: "Медиаматериалы"
description: "Официальные кнопки Wails с предпросмотром и кодом для вставки в HTML и Markdown"
slug: "community/media"
sourcePath: "community/media.md"
---

Используйте официальные графические материалы Wails, когда рассказываете о проекте, созданном с помощью Wails. Кнопки ниже ведут на [wails.io](https://wails.io), а их изображения размещены на сайте Wails, поэтому копировать изображение в свой проект не нужно.

## Создано с помощью Wails

Выберите вариант с достаточным контрастом по отношению к фону. Каждый предпросмотр показывает именно то изображение со ссылкой, которое увидят ваши посетители.

### Светлый фон

<div class="wails-media-preview wails-media-preview-light">  
<p>Стандартный размер (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Создано с помощью Wails" /></a>  
<p>Компактный размер (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="Создано с помощью Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

Для компактного варианта используйте `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### Тёмный фон

<div class="wails-media-preview wails-media-preview-dark">  
<p>Стандартный размер (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Создано с помощью Wails" /></a>  
<p>Компактный размер (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="Создано с помощью Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

Для компактного варианта используйте `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

Markdown отображает изображение в его исходном размере — `180` на `48` пикселей. Чтобы задать другой размер отображения, используйте вариант HTML. Сохраняйте исходное соотношение сторон кнопок `15:4`.
