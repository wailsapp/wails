---
title: "Medien"
description: "Offizielle Wails-Schaltflächen mit Vorschauen sowie Einbettungscode für HTML und Markdown"
slug: "community/media"
sourcePath: "community/media.md"
---

Verwenden Sie beim Teilen eines mit Wails erstellten Projekts die offiziellen Wails-Ressourcen. Die folgenden Schaltflächen verweisen auf [wails.io](https://wails.io) und werden weiterhin von Wails gehostet. Sie müssen das Bild daher nicht in Ihr Projekt kopieren.

## Erstellt mit Wails

Wählen Sie die Version, die einen ausreichenden Kontrast zum Hintergrund bietet. Jede Vorschau zeigt exakt das verlinkte Bild, das Ihre Besucher sehen werden.

### Heller Hintergrund

<div class="wails-media-preview wails-media-preview-light">  
<p>Standard (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Erstellt mit Wails" /></a>  
<p>Kompakt (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="Erstellt mit Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

Verwenden Sie für die kompakte Version `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### Dunkler Hintergrund

<div class="wails-media-preview wails-media-preview-dark">  
<p>Standard (<code>180</code> × <code>48</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Erstellt mit Wails" /></a>  
<p>Kompakt (<code>90</code> × <code>24</code> px):</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="Erstellt mit Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

Verwenden Sie für die kompakte Version `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

Markdown stellt das Bild in seiner intrinsischen Größe von `180` × `48` Pixeln dar. Verwenden Sie die HTML-Version, wenn Sie eine andere Anzeigegröße festlegen müssen. Behalten Sie für die Schaltflächen das ursprüngliche Seitenverhältnis von `15:4` bei.
