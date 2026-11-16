---
title: "Médias"
description: "Boutons Wails officiels avec aperçus et codes d’intégration HTML et Markdown"
slug: "community/media"
sourcePath: "community/media.md"
---

Utilisez les ressources Wails officielles lorsque vous partagez un projet créé avec Wails. Les boutons ci-dessous renvoient vers [wails.io](https://wails.io) et restent hébergés par Wails ; vous n’avez donc pas besoin de copier l’image dans votre projet.

## Créé avec Wails

Choisissez la version offrant un contraste suffisant avec son arrière-plan. Chaque aperçu correspond exactement à l’image liée que vos visiteurs verront.

### Arrière-plan clair

<div class="wails-media-preview wails-media-preview-light">  
<p>Standard (<code>180</code> × <code>48</code> px) :</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="180" height="48" alt="Créé avec Wails" /></a>  
<p>Compact (<code>90</code> × <code>24</code> px) :</p>  
<a href="https://wails.io"><img src="/img/wails-button-light.svg" width="90" height="24" alt="Créé avec Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-light.svg" width="180" height="48" alt="Built with Wails">
</a>
```

Pour la version compacte, utilisez `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-light.svg)](https://wails.io)
```

### Arrière-plan sombre

<div class="wails-media-preview wails-media-preview-dark">  
<p>Standard (<code>180</code> × <code>48</code> px) :</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="180" height="48" alt="Créé avec Wails" /></a>  
<p>Compact (<code>90</code> × <code>24</code> px) :</p>  
<a href="https://wails.io"><img src="/img/wails-button-dark.svg" width="90" height="24" alt="Créé avec Wails" /></a>  
</div>

#### HTML

```html
<a href="https://wails.io">
  <img src="https://wails.io/img/wails-button-dark.svg" width="180" height="48" alt="Built with Wails">
</a>
```

Pour la version compacte, utilisez `width="90" height="24"`.

#### Markdown

```md
[![Built with Wails](https://wails.io/img/wails-button-dark.svg)](https://wails.io)
```

Markdown affiche l’image à sa taille intrinsèque de `180` par `48` pixels. Utilisez la version HTML lorsque vous devez définir une taille d’affichage différente. Conservez les proportions d’origine `15:4` des boutons.
