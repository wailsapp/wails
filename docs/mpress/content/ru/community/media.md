---
title: "Медиа"
description: "Официальные кнопки Wails с предварительным просмотром и кодом для вставки в HTML и Markdown"
slug: "community/media"
sourcePath: "community/media.md"
---

Используйте официальные ресурсы Wails, когда рассказываете о проекте, созданном с помощью Wails. Кнопки ниже ведут на [wails.io](https://wails.io) и по-прежнему размещаются на серверах Wails, поэтому копировать изображение в свой проект не нужно.

## Создано с помощью Wails

Выберите вариант с достаточной контрастностью относительно фона. В каждом окне предварительного просмотра показано именно то изображение, которое увидят посетители по ссылке.

### Светлый фон

<div class="wails-media-preview wails-media-preview-light">  
<p>Стандартная (<code>180</code> × <code>48</code> пкс):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Создано с помощью Wails" /></a>  
<p>Компактная (<code>90</code> × <code>24</code> пкс):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="Создано с помощью Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

Для компактной версии используйте `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### Тёмный фон

<div class="wails-media-preview wails-media-preview-dark">  
<p>Стандартная (<code>180</code> × <code>48</code> пкс):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Создано с помощью Wails" /></a>  
<p>Компактная (<code>90</code> × <code>24</code> пкс):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="Создано с помощью Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

Для компактной версии используйте `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

Markdown отображает изображение с исходным размером `180` на `48` пикселей. Если требуется задать другой размер отображения, используйте версию HTML. Сохраняйте исходные пропорции кнопок `15:4`.
